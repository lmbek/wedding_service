package logger

import (
	"context"
	"log/slog"
	"os"
)

type Logger interface {
	Slog() *slog.Logger
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
}

type logger struct {
	slogLogger *slog.Logger
}

func NewLogger(appName string, debugLevel int) Logger {
	// Map our integer debug level to slog.Level
	var slogLevel slog.Level
	switch debugLevel {
	case 4: // All/Debug
		slogLevel = slog.LevelDebug
	case 3: // Info
		slogLevel = slog.LevelInfo
	case 2: // Warning
		slogLevel = slog.LevelWarn
	case 1: // Error
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})
	slogLogger := slog.New(textHandler).With(slog.String("app", appName))
	slog.SetDefault(slogLogger)
	return &logger{slogLogger}
}

func (l *logger) Slog() *slog.Logger {
	return l.slogLogger
}

func (l *logger) Info(msg string, args ...any) {
	l.slogLogger.Info(msg, args)
}

func (l *logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.slogLogger.InfoContext(ctx, msg, args)
}

func (l *logger) Debug(msg string, args ...any) {
	l.slogLogger.Debug(msg, args)
}

func (l *logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.slogLogger.DebugContext(ctx, msg, args)
}

func (l *logger) Warn(msg string, args ...any) {
	l.slogLogger.Warn(msg, args)
}

func (l *logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.slogLogger.WarnContext(ctx, msg, args)
}

func (l *logger) Error(msg string, args ...any) {
	l.slogLogger.Error(msg, args)
}

func (l *logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.slogLogger.ErrorContext(ctx, msg, args)
}
