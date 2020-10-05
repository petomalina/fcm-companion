package main

import (
	"context"
	"github.com/blendle/zapdriver"
	"github.com/petomalina/fcm-companion/pkg/companion"
	"github.com/petomalina/fcm-companion/pkg/serverutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	ctx := context.Background()

	config := zapdriver.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger, err := config.Build(zapdriver.WrapCore(
		zapdriver.ReportAllErrors(true),
		zapdriver.ServiceName("fcm-companion"),
	))
	if err != nil {
		panic(err)
	}

	svc, err := companion.New(ctx, os.Getenv("PROJECT_ID"), logger, "")
	if err != nil {
		logger.Fatal("Cannot initialize companion", zap.Error(err))
	}

	logger.Info("Starting the server")
	if err := serverutil.Serve(
		serverutil.WithContext(ctx),
		serverutil.WithPort(os.Getenv("PORT")),
		serverutil.WithServices(svc),
		serverutil.WithGRPC(),
		serverutil.WithPubSub(),
	); err != nil {
		logger.Fatal("Serving crashed", zap.Error(err))
	}
}
