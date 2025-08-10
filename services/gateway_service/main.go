package main

import (
	"gateway_service/app"
	"gateway_service/certificate"
	"gateway_service/config"
	"gateway_service/logger"
	"gateway_service/tracer"
	"log/slog"
	"os"
)

func main() {
	startGatewayService()
}

func startGatewayService() {
	var err error

	// config
	c, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	// logger, tracer, cert manager
	newLogger := logger.NewLogger(c.ServiceName(), c.DebugLevel())
	newTracer := tracer.NewTracer("X-Trace-Id")
	newACME, err := certificate.NewAutoCertManager(c)
	if err != nil {
		slog.Error("failed to init ACME", slog.Any("error", err))
		return
	}

	// application
	a, err := app.NewApp(c, newLogger, newTracer, newACME)
	if err != nil {
		slog.Error("application new app failed", slog.Any("error", err))
		return
	}

	err = a.Run()
	if err != nil {
		slog.Error("application run failed", slog.Any("error", err))
		return
	}
}
