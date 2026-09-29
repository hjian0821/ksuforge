package main

import (
	"embed"
	"log/slog"
	"os"

	"github.com/hjian0821/ksuforge/internal/logging"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	logger, closeLog := logging.Setup(logging.Options{Level: slog.LevelInfo, Prefix: "ksuforge-gui"})
	logger.Info("app.gui_started")

	app := NewApp()
	err := wails.Run(&options.App{
		Title:     "ksuforge",
		Width:     1180,
		Height:    840,
		MinWidth:  720,
		MinHeight: 560,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 9, G: 11, B: 16, A: 1},
		OnStartup:        app.startup,
		Bind:             []interface{}{app},
	})
	if err != nil {
		logger.Error("app.exit_error", "error", err)
		closeLog()
		os.Exit(1)
	}
	closeLog()
}
