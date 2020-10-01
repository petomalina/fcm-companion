package serverutil

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/petomalina/xrpc/pkg/multiplexer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
)

type GRPCServer interface {
	Register(server *grpc.Server)
}

type GRPCGateway interface {
	RegisterGateway(ctx context.Context, mux *runtime.ServeMux, bind string, opts []grpc.DialOption) error
}

func Serve(lis net.Listener, logger *zap.Logger, handlers ...multiplexer.Handler) {
	handler := multiplexer.Make(nil,
		handlers...,
	)
	srv := http.Server{Handler: handler}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Info("OS Interrupt caught, shutting down")
		_ = srv.Close()
	}()

	logger.Info("Starting the FCM Companion")
	if err := srv.Serve(lis); err != nil {
		logger.Info("Exiting the FCM Companion", zap.Error(err))
	}
}
