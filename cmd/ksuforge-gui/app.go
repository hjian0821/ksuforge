package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hjian0821/ksuforge/internal/command"
	"github.com/hjian0821/ksuforge/internal/config"
	"github.com/hjian0821/ksuforge/internal/device"
	"github.com/hjian0821/ksuforge/internal/i18n"
	"github.com/hjian0821/ksuforge/internal/logging"
	appPaths "github.com/hjian0821/ksuforge/internal/paths"
	"github.com/hjian0821/ksuforge/internal/rom"
	"github.com/hjian0821/ksuforge/internal/workflow"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type InstallRequest struct {
	Serial, Mode, Partition, OutputDir, Catalog, KSUd, KMI, Kernel, Mirrors, ROMPath string
	Connections                                                                      int
	Flash, NoReboot, VerifyROM                                                       bool
}

type TaskState struct {
	Status   string `json:"status"`
	Language string `json:"language,omitempty"`
	Log      string `json:"log"`
	Error    string `json:"error,omitempty"`
	Result   any    `json:"result,omitempty"`
}

// Settings is the user-editable configuration shown on the settings page.
type Settings struct {
	OutputDir     string `json:"outputDir"`
	Language      string `json:"language"`
	Mode          string `json:"mode"`
	Partition     string `json:"partition"`
	AutoFlash     bool   `json:"autoFlash"`
	NoReboot      bool   `json:"noReboot"`
	CatalogURL    string `json:"catalogUrl"`
	Mirrors       string `json:"mirrors"`
	Connections   int    `json:"connections"`
	KSUd          string `json:"ksud"`
	LogAutoScroll bool   `json:"logAutoScroll"`
}

type App struct {
	ctx     context.Context
	devices device.Client
	mu      sync.Mutex
	state   TaskState
	cfgMu   sync.Mutex
	cfg     config.Config
}

func NewApp() *App {
	cfg := config.Load()
	return &App{
		devices: device.Client{Runner: command.ExecRunner{SearchPaths: command.DesktopSearchPaths()}},
		state:   TaskState{Status: "idle", Language: cfg.LanguageOrDefault("zh")},
		cfg:     cfg,
	}
}

func (a *App) startup(ctx context.Context) { a.ctx = ctx }

func (a *App) Devices() ([]device.Device, error) { return a.devices.List(a.ctx) }

func (a *App) Inspect(serial string) (device.Device, error) {
	if strings.TrimSpace(serial) == "" {
		return device.Device{}, i18n.NewError("error.device_required")
	}
	return a.devices.Inspect(a.ctx, serial)
}

func (a *App) GetSettings() Settings {
	a.cfgMu.Lock()
	defer a.cfgMu.Unlock()
	return Settings{
		OutputDir:     a.cfg.OutputDirOrDefault(),
		Language:      a.cfg.LanguageOrDefault("zh"),
		Mode:          a.cfg.ModeOrDefault(),
		Partition:     a.cfg.PartitionOrDefault(),
		AutoFlash:     a.cfg.AutoFlash,
		NoReboot:      a.cfg.NoReboot,
		CatalogURL:    a.cfg.CatalogURLOr(rom.DefaultCatalogURL),
		Mirrors:       a.cfg.Mirrors,
		Connections:   a.cfg.ConnectionsOrDefault(),
		KSUd:          a.cfg.KSUd,
		LogAutoScroll: a.cfg.LogAutoScrollDefault(),
	}
}

func (a *App) SaveSettings(s Settings) error {
	if s.Language != "zh" && s.Language != "en" {
		return i18n.NewError("error.language_unsupported")
	}
	if s.Mode != "lkm" && s.Mode != "gki" {
		return i18n.NewError("error.mode_invalid")
	}
	switch s.Partition {
	case "auto", "boot", "init_boot":
	default:
		return i18n.NewError("error.partition_invalid")
	}
	if strings.TrimSpace(s.OutputDir) == "" {
		return i18n.NewError("error.output_dir_required")
	}
	catalog := strings.TrimSpace(s.CatalogURL)
	if catalog == "" {
		catalog = rom.DefaultCatalogURL
	}
	if err := rom.ValidateCatalogURL(catalog); err != nil {
		return err
	}
	if s.Connections < 1 || s.Connections > 32 {
		return i18n.NewError("error.connections_invalid")
	}
	for _, mirror := range splitMirrors(s.Mirrors) {
		if err := rom.ValidateMirrorBaseURL(mirror); err != nil {
			return i18n.WrapError("error.mirror_invalid", err, mirror)
		}
	}
	autoScroll := s.LogAutoScroll
	a.cfgMu.Lock()
	defer a.cfgMu.Unlock()
	a.cfg = config.Config{
		OutputDir:     strings.TrimSpace(s.OutputDir),
		Language:      s.Language,
		Mode:          s.Mode,
		Partition:     s.Partition,
		AutoFlash:     s.AutoFlash,
		NoReboot:      s.NoReboot,
		CatalogURL:    catalog,
		Mirrors:       strings.TrimSpace(s.Mirrors),
		Connections:   s.Connections,
		KSUd:          strings.TrimSpace(s.KSUd),
		LogAutoScroll: &autoScroll,
	}
	if err := a.cfg.Save(); err != nil {
		return err
	}
	slog.Info("app.settings_saved", "output_dir", a.cfg.OutputDir, "language", a.cfg.Language, "mode", a.cfg.Mode, "partition", a.cfg.Partition)
	return nil
}

func (a *App) DefaultOutputDirectory() string {
	a.cfgMu.Lock()
	defer a.cfgMu.Unlock()
	return a.cfg.OutputDirOrDefault()
}

func (a *App) currentLanguage() string {
	a.cfgMu.Lock()
	defer a.cfgMu.Unlock()
	return a.cfg.LanguageOrDefault("zh")
}

func (a *App) ChooseOutputDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: "选择输出目录"})
}

func (a *App) ChooseFile(title string) (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: title})
}

func (a *App) ChooseOTAFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择 Recovery OTA 卡刷包",
		Filters: []runtime.FileFilter{{
			DisplayName: "Recovery OTA (*.zip)",
			Pattern:     "*.zip",
		}},
	})
}

// Confirm shows a native question dialog. WKWebView does not implement the
// JavaScript dialog delegates, so window.confirm() always returns false there.
func (a *App) Confirm(title, message, confirmLabel, cancelLabel string) (bool, error) {
	if confirmLabel == "" {
		confirmLabel = "Confirm"
	}
	if cancelLabel == "" {
		cancelLabel = "Cancel"
	}
	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         title,
		Message:       message,
		Buttons:       []string{cancelLabel, confirmLabel},
		DefaultButton: confirmLabel,
		CancelButton:  cancelLabel,
	})
	if err != nil {
		return false, err
	}
	return result == confirmLabel, nil
}

func (a *App) State() TaskState {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state
}

func (a *App) Start(request InstallRequest) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.state.Status == "running" {
		return i18n.NewError("error.task_running")
	}
	if strings.TrimSpace(request.OutputDir) == "" {
		request.OutputDir = a.DefaultOutputDirectory()
	}
	language := a.currentLanguage()
	slog.Info("app.task_started", "serial", request.Serial, "mode", request.Mode, "partition", request.Partition,
		"flash", request.Flash, "no_reboot", request.NoReboot, "output_dir", request.OutputDir)
	a.state = TaskState{Status: "running", Language: language}
	go a.run(request, language)
	return nil
}

func (a *App) run(request InstallRequest, language string) {
	ctx := logging.WithLanguage(a.ctx, language)
	result, err := (workflow.Service{Devices: a.devices, Out: &appLog{app: a}}).Run(ctx, workflow.Options{
		Serial: request.Serial, Mode: request.Mode, Partition: request.Partition,
		OutputDir: request.OutputDir, Catalog: request.Catalog, KSUd: request.KSUd, KMI: request.KMI,
		Kernel: request.Kernel, ROMPath: request.ROMPath, VerifyROM: request.VerifyROM, Yes: request.Flash, NoReboot: request.NoReboot,
		Mirrors: splitMirrors(request.Mirrors), Connections: request.Connections,
	})
	a.mu.Lock()
	if err != nil {
		a.state.Status, a.state.Error = "failed", i18n.ErrorText(err, language)
	} else {
		a.state.Status, a.state.Result = "succeeded", result
	}
	logText := a.state.Log
	if a.state.Error != "" {
		logText += "\n" + a.state.Error
	}
	a.mu.Unlock()
	if err != nil {
		slog.Error("app.task_failed", "error", err)
	} else {
		slog.Info("app.task_completed", "partition", result.Partition, "patched", result.PatchedImage)
	}
	writeRunLog(logText)
}

func splitMirrors(value string) []string {
	return strings.FieldsFunc(value, func(r rune) bool { return r == '\n' || r == '\r' || r == ',' })
}

// writeRunLog persists the run transcript under the OS log directory so failed
// runs can be inspected after the app closes.
func writeRunLog(text string) {
	dir := appPaths.LogDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}
	name := "ksuforge-" + time.Now().Format("20060102-150405") + ".log"
	_ = os.WriteFile(filepath.Join(dir, name), []byte(text), 0o600)
}

type appLog struct{ app *App }

func (w *appLog) Write(p []byte) (int, error) {
	w.app.mu.Lock()
	if len(w.app.state.Log) < 2<<20 {
		w.app.state.Log += string(p)
	}
	w.app.mu.Unlock()
	return len(p), nil
}
