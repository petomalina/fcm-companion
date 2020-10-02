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

	var grpcServer *grpc.Server
	var gateway *runtime.ServeMux
	var handlers []multiplexer.Handler

	if ctx.grpcEnabled {
		grpcServer = grpc.NewServer()
		handlers = append(handlers, multiplexer.GRPCHandler(grpcServer))
	}

	if ctx.pubsubEnabled || ctx.gatewayEnabled {
		gateway = runtime.NewServeMux()
	}

	// pubsub must be registered before the gateway
	if ctx.pubsubEnabled {
		handlers = append(handlers, multiplexer.PubSubHandler(gateway))
	}

	if ctx.gatewayEnabled {
		handlers = append(handlers, multiplexer.HTTPHandler(grpcServer))
	}

	// register all services provided by the user
	for _, svc := range ctx.services {
		if s, ok := svc.(GRPCServer); ok {
			s.Register(grpcServer)
		}

		if s, ok := svc.(GRPCGateway); ok {
			err := s.RegisterGateway(ctx.ctx, gateway, bind, []grpc.DialOption{grpc.WithInsecure()})
			if err != nil {
				return err
			}
		}
	}

	handler := multiplexer.Make(nil,
		handlers...,
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

func makeGRPCServerWithGateway(ctx context.Context, bind string, services ...interface{}) (*grpc.Server, *runtime.ServeMux, error) {
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
