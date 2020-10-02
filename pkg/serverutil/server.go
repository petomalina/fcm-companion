package serverutil

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/petomalina/xrpc/pkg/multiplexer"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
)

// GRPCServer is any structure that implements the Register method
// for grpc server registration.
type GRPCServer interface {
	Register(server *grpc.Server)
}

// GRPCGateway is any structure that implements the RegisterGateway method
// which registers a new gateway to the http mux.
type GRPCGateway interface {
	RegisterGateway(ctx context.Context, mux *runtime.ServeMux, bind string, opts []grpc.DialOption) error
}

func Serve(opts ...ServeContextOption) error {
	ctx := &ServeContext{}

	for _, o := range opts {
		o(ctx)
	}

	lis, err := net.Listen("tcp", ":"+ctx.port)
	if err != nil {
		return err
	}
	// get the correct bind if the getenv returned an empty string
	bind := lis.Addr().(*net.TCPAddr).String()

	if ctx.onListen != nil {
		ctx.onListen()
	}

	// create and register the grpc server
	grpcServer, gwmux, err := MakeGRPCServerWithGateway(ctx.ctx, bind, ctx.services...)
	if err != nil {
		return err
	}

	handler := multiplexer.Make(nil,
		multiplexer.GRPCHandler(grpcServer),
		multiplexer.PubSubHandler(gwmux),
		multiplexer.HTTPHandler(gwmux),
	)
	srv := http.Server{Handler: handler}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = srv.Close()
	}()

	err = srv.Serve(lis)

	if ctx.onExit != nil {
		ctx.onExit()
	}

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func MakeGRPCServerWithGateway(ctx context.Context, bind string, services ...interface{}) (*grpc.Server, *runtime.ServeMux, error) {
	grpcServer := grpc.NewServer()
	gwmux := runtime.NewServeMux()

	for _, svc := range services {
		if s, ok := svc.(GRPCServer); ok {
			s.Register(grpcServer)
		}

		if s, ok := svc.(GRPCGateway); ok {
			err := s.RegisterGateway(ctx, gwmux, bind, []grpc.DialOption{grpc.WithInsecure()})
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return grpcServer, gwmux, nil
}
