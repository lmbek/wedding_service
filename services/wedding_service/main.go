package main

import (
	"log/slog"
	"wedding_service/app"
)

func main() {

	a, err := app.NewApp()
	if err != nil {
		slog.Error("failed to init app", slog.Any("error", err))
		return
	}
	_ = a.Start()
}
