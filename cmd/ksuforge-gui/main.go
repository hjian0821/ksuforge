package main

import (
	"embed"
	"log/slog"
	"os"
	"runtime"

	"github.com/hjian0821/ksuforge/internal/logging"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	logger, closeLog := logging.Setup(logging.Options{Level: slog.LevelInfo, Prefix: "ksuforge-gui"})
	logger.Info("app.gui_started")

	app := NewApp()
	var appMenu *menu.Menu
	if runtime.GOOS == "darwin" {
		appMenu = menu.NewMenu()
		appMenu.Append(menu.AppMenu())
		appMenu.Append(menu.EditMenu())
	}
	err := wails.Run(&options.App{
		Title:     "KSUForge",
		Width:     1180,
		Height:    840,
		MinWidth:  720,
		MinHeight: 560,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 16, G: 21, B: 27, A: 1},
		Menu:             appMenu,
		Mac:              &mac.Options{Appearance: mac.NSAppearanceNameDarkAqua},
		Windows:          &windows.Options{Theme: windows.Dark},
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
