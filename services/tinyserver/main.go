package main

import (
	"log/slog"
	"tinyserver/app"
)

func main() {
	// Create and run the tinyserver application
	a := app.NewApp()
	if err := a.Run(); err != nil {
		slog.Error("tinyserver exited", slog.Any("error", err))
	}
}
