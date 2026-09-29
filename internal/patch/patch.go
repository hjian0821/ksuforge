package patch

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hjian0821/ksuforge/internal/command"
	"github.com/hjian0821/ksuforge/internal/i18n"
	bootimage "github.com/hjian0821/ksuforge/internal/image"
	"github.com/hjian0821/ksuforge/internal/ksud"
	"github.com/hjian0821/ksuforge/internal/logging"
	appPaths "github.com/hjian0821/ksuforge/internal/paths"
)

type Options struct {
	Boot, KMI, Mode, Kernel, OutputDir, KSUd string
}

type Service struct{ Out io.Writer }

type patchState struct {
	StockSHA256   string `json:"stock_sha256"`
	KMI           string `json:"kmi"`
	Mode          string `json:"mode"`
	KSUd          string `json:"ksud"`
	KSUdMod       int64  `json:"ksud_mod"`
	PatchedSHA256 string `json:"patched_sha256"`
}

// ResolveKSUd returns the ksud executable to use. An explicit override is
// honoured as-is; otherwise ksud is looked up on PATH (which also covers the
// copy shipped next to the application) and finally falls back to the binary
// embedded in the application.
func ResolveKSUd(path string) (string, error) {
	runner := command.ExecRunner{SearchPaths: command.DesktopSearchPaths()}
	if strings.TrimSpace(path) != "" {
		tool, err := runner.LookPath(path)
		if err != nil {
			return "", i18n.WrapError("error.ksud_override_not_found", err)
		}
		return tool, nil
	}
	if tool, err := runner.LookPath("ksud"); err == nil {
		return tool, nil
	}
	if ksud.Embedded() {
		tool, err := ksud.Extract()
		if err != nil {
			return "", i18n.WrapError("error.ksud_extract", err)
		}
		return tool, nil
	}
	return "", i18n.NewError("error.ksud_missing")
}

func (s Service) Run(ctx context.Context, o Options) (string, error) {
	if s.Out == nil {
		s.Out = io.Discard
	}
	logger := logging.WithWriter(s.Out, logging.Language(ctx))
	stockInfo, err := bootimage.Inspect(o.Boot)
	if err != nil {
		return "", i18n.WrapError("error.stock_image_check", err)
	}
	if strings.TrimSpace(o.KMI) == "" {
		return "", i18n.NewError("error.kmi_required")
	}
	if o.OutputDir == "" {
		o.OutputDir = appPaths.DefaultOutputDir()
	}
	if err := os.MkdirAll(o.OutputDir, 0o755); err != nil {
		return "", err
	}
	tool, err := ResolveKSUd(o.KSUd)
	if err != nil {
		return "", err
	}
	name := "kernelsu-patched-" + filepath.Base(o.Boot)
	destination := filepath.Join(o.OutputDir, name)
	statePath := destination + ".json"
	if _, err := os.Stat(destination); err == nil {
		if _, err := bootimage.Inspect(destination); err == nil && patchStateMatches(statePath, stockInfo.SHA256, o.KMI, o.Mode, tool) {
			logger.Info("patch.image_reused", "path", destination)
			return destination, nil
		}
		logger.Info("patch.image_stale", "path", destination)
		if err := os.Remove(destination); err != nil {
			return "", err
		}
		_ = os.Remove(statePath)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	args := []string{"boot-patch", "--boot", o.Boot, "--kmi", o.KMI, "--out", o.OutputDir, "--out-name", name}
	if o.Mode == "gki" {
		if strings.TrimSpace(o.Kernel) == "" {
			return "", i18n.NewError("error.gki_kernel_required")
		}
		args = append(args, "--kernel", o.Kernel)
	} else if o.Mode != "" && o.Mode != "lkm" {
		return "", i18n.NewError("error.mode_invalid")
	}
	logger.Info("patch.started", "image", o.Boot, "kmi", o.KMI, "mode", o.Mode)
	cmd := command.CommandContext(ctx, tool, args...)
	cmd.Stdout, cmd.Stderr = s.Out, s.Out
	if err := cmd.Run(); err != nil {
		return "", i18n.WrapError("error.patch_failed", err)
	}
	if _, err := bootimage.Inspect(destination); err != nil {
		return "", i18n.WrapError("error.patched_image_check", err)
	}
	writePatchState(statePath, stockInfo.SHA256, o.KMI, o.Mode, tool, destination)
	logger.Info("patch.image_created", "path", destination)
	slog.Info("patch.completed", "image", o.Boot, "patched", destination, "kmi", o.KMI, "mode", o.Mode)
	return destination, nil
}

// patchStateMatches reports whether an existing patched image was produced from
// the same stock image, KMI, mode and ksud binary. A missing or unreadable
// sidecar forces a re-patch so a stale or foreign image is never reused.
func patchStateMatches(statePath, stockSHA, kmi, mode, ksud string) bool {
	data, err := os.ReadFile(statePath)
	if err != nil {
		return false
	}
	var state patchState
	if err := json.Unmarshal(data, &state); err != nil {
		return false
	}
	if state.StockSHA256 != stockSHA || state.KMI != kmi || state.Mode != mode {
		return false
	}
	if state.KSUd != "" && state.KSUd != ksud {
		return false
	}
	if state.KSUdMod != 0 && state.KSUdMod != fileModTime(ksud) {
		return false
	}
	return true
}

func writePatchState(statePath, stockSHA, kmi, mode, ksud, patched string) {
	state := patchState{StockSHA256: stockSHA, KMI: kmi, Mode: mode, KSUd: ksud, KSUdMod: fileModTime(ksud)}
	if info, err := bootimage.Inspect(patched); err == nil {
		state.PatchedSHA256 = info.SHA256
	}
	data, err := json.Marshal(state)
	if err != nil {
		return
	}
	_ = os.WriteFile(statePath, data, 0o600)
}

func fileModTime(path string) int64 {
	if info, err := os.Stat(path); err == nil {
		return info.ModTime().Unix()
	}
	return 0
}
