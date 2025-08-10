package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"wedding_service/config"
	"wedding_service/logging"
	"wedding_service/webserver"
	"wedding_service/webserver/website/frontend"
)

// App represents the application with all its components
type App interface {
	// Start starts the application and blocks until it's stopped
	Start() error

	// Shutdown gracefully shuts down the application
	Shutdown() error
}

// application is the implementation of the App interface
type application struct {
	ctx       context.Context
	cancel    context.CancelFunc
	config    config.Config
	webserver webserver.Webserver
	frontend  frontend.Frontend
}

// NewApp creates a new application instance
func NewApp(ctx context.Context) (App, error) {
	// Initialize the environment
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	// Create a context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(ctx)

	// Set up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		slog.Info("Received shutdown signal")
		cancel()
	}()

	app := &application{
		ctx:    ctx,
		cancel: cancel,
		config: cfg,
	}

	if err := app.initialize(); err != nil {
		return nil, err
	}
	return app, nil
}

// Initialize initializes all components of the application
func (a *application) initialize() error {
	var err error
	// Initialize logging (slog)
	if err = logging.Init("wedding_service", a.config.DebugLevel()); err != nil {
		return err
	}

	// Initialize the frontend
	a.frontend, err = frontend.NewFrontend(a.config.FrontendPath(), a.config.HotReloadEnabled())
	if err != nil {
		return err
	}

	// Initialize the webserver
	a.webserver, err = webserver.NewWebserver(a.config, a.frontend)
	if err != nil {
		return err
	}

	return nil
}

// Start starts the application and blocks until it's stopped
func (a *application) Start() error {
	// Start the webserver in a goroutine
	go func() {
		if err := a.webserver.ListenAndServe(); err != nil {
			slog.Error("Error starting webserver", slog.Any("error", err))
			a.cancel()
		}
	}()

	slog.Info("Application started",
		slog.String("http_port", a.config.HttpPort()))

	// Wait for context cancellation (from signal or error)
	<-a.ctx.Done()

	return a.Shutdown()
}

// Shutdown gracefully shuts down the application
func (a *application) Shutdown() error {
	slog.Info("Shutting down application")

	// Close the webserver
	if a.webserver != nil {
		if err := a.webserver.Close(); err != nil {
			slog.Error("Error closing webserver", slog.Any("error", err))
		}
	}

	// Shutdown logging
	if err := logging.Shutdown(a.ctx); err != nil {
		slog.Error("Error shutting down logging", slog.Any("error", err))
	}

	return nil
}
