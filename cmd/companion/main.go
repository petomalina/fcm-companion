package main

import (
	"context"
	"github.com/blendle/zapdriver"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/petomalina/fcm-companion/pkg/companion"
	"github.com/petomalina/fcm-companion/pkg/serverutil"
	"github.com/petomalina/xrpc/pkg/multiplexer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	ctx := context.Background()
	requestedBind := ":" + os.Getenv("PORT")

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

	svc, err := companion.New(ctx, os.Getenv("PROJECT_ID"), logger)
	if err != nil {
		logger.Fatal("An error occurred when initializing companion", zap.Error(err))
	}

	lis, err := net.Listen("tcp", requestedBind)
	if err != nil {
		logger.Fatal("Cannot start listener", zap.Error(err))
	}
	bind := lis.Addr().(*net.TCPAddr).String()

	// create and register the grpc server
	grpcServer := grpc.NewServer()
	gwmux := runtime.NewServeMux()

	svc.Register(grpcServer)
	err = svc.RegisterGateway(ctx, gwmux, bind, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Fatal("failed to register gateway endpoint handler", zap.Error(err))
	}

	serverutil.Serve(
		lis,
		logger,
		multiplexer.GRPCHandler(grpcServer),
		multiplexer.PubSubHandler(gwmux),
		multiplexer.HTTPHandler(gwmux),
	)
}
