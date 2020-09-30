package main

import (
	"context"
	"github.com/blendle/zapdriver"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/petomalina/fcm-companion/apis/go-sdk/notification/v1"
	"github.com/petomalina/fcm-companion/pkg/notification"
	"github.com/petomalina/xrpc/examples/api"
	"github.com/petomalina/xrpc/pkg/multiplexer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()

	// create the zap logger for future use
	config := zapdriver.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger, err := config.Build(zapdriver.WrapCore(
		zapdriver.ReportAllErrors(true),
		zapdriver.ServiceName("Echo"),
	))
	if err != nil {
		panic(err)
	}

	// create and register the grpc server
	grpcServer := grpc.NewServer()
	notificationSvc := &notification.Service{Logger: logger}
	v1.RegisterNotificationServiceServer(grpcServer, notificationSvc)
	reflection.Register(grpcServer)

	// create the grpc-gateway server and register to grpc server
	gwmux := runtime.NewServeMux()
	err = api.RegisterEchoServiceHandlerFromEndpoint(ctx, gwmux, ":"+os.Getenv("PORT"), []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Fatal("gw: failed to register: %v", zap.Error(err))
	}

	// make multiplexer
	handler := multiplexer.Make(nil,
		// filters all application/grpc messages into the grpc server
		multiplexer.GRPCHandler(grpcServer),
		// filters all messages with Google Agent into the gwmux and
		// unpacks the PubSub message
		multiplexer.PubSubHandler(gwmux),
		// defaults all other messages into the http multiplexer
		multiplexer.HTTPHandler(gwmux),
	)
	srv := http.Server{Addr: ":" + os.Getenv("PORT"), Handler: handler}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Info("shutting down server")
		grpcServer.GracefulStop()
		_ = srv.Close()
	}()

	logger.Info("starting grpcServer")
	if err = srv.ListenAndServe(); err != nil {
		logger.Info("grpcServer exit", zap.Error(err))
	}
}
