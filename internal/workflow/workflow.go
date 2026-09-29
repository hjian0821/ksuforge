package workflow

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"

	"github.com/hjian0821/ksuforge/internal/device"
	"github.com/hjian0821/ksuforge/internal/flash"
	"github.com/hjian0821/ksuforge/internal/logging"
	ksupatch "github.com/hjian0821/ksuforge/internal/patch"
	appPaths "github.com/hjian0821/ksuforge/internal/paths"
	"github.com/hjian0821/ksuforge/internal/prepare"
)

type Options struct {
	Serial, Mode, Partition, OutputDir, Catalog, KSUd, KMI, Kernel, ROMPath string
	Yes, NoReboot, VerifyROM                                                bool
	Mirrors                                                                 []string
	Connections                                                             int
}

type Result struct {
	StockImage, PatchedImage, ROMPath, Partition string
}

type Service struct {
	Devices device.Client
	Out     io.Writer
}

func (s Service) Run(ctx context.Context, o Options) (Result, error) {
	if s.Out == nil {
		s.Out = io.Discard
	}
	if _, err := ksupatch.ResolveKSUd(o.KSUd); err != nil {
		return Result{}, err
	}
	logger := logging.WithWriter(s.Out, logging.Language(ctx))
	slog.Info("workflow.started", "serial", o.Serial, "mode", o.Mode, "partition", o.Partition, "flash", o.Yes, "output_dir", o.OutputDir)
	if o.Yes {
		logger.Info("workflow.step.check_usb", "step", 1, "total", 8)
		logger.Info("workflow.step.detect_device", "step", 2, "total", 8)
		logger.Info("workflow.step.acquire_stock", "step", 3, "total", 8)
	} else {
		logger.Info("workflow.step.check_usb", "step", 1, "total", 5)
		logger.Info("workflow.step.detect_device", "step", 2, "total", 5)
		logger.Info("workflow.step.acquire_stock", "step", 3, "total", 5)
	}
	prepared, err := (prepare.Service{Devices: s.Devices, HTTPClient: prepare.NewHTTPClient(), Out: s.Out}).Run(ctx, prepare.Options{
		Serial: o.Serial, Mode: o.Mode, Partition: o.Partition, OutputDir: o.OutputDir, Catalog: o.Catalog, KMI: o.KMI, RequireKMI: true,
		ROMPath: o.ROMPath, VerifyROM: o.VerifyROM, Mirrors: o.Mirrors, Connections: o.Connections,
	})
	if err != nil {
		return Result{}, err
	}
	if o.Yes {
		logger.Info("workflow.step.stock_verified", "step", 4, "total", 8)
		logger.Info("workflow.step.patch", "step", 5, "total", 8)
	} else {
		logger.Info("workflow.stock_acquired", "path", prepared.ImagePath)
		logger.Info("workflow.step.patch", "step", 4, "total", 5)
	}
	kmi := o.KMI
	if kmi == "" {
		kmi = prepared.Device.KMI
	}
	patchedDir := appPaths.PatchedDir(o.OutputDir, prepared.Device.Codename, prepared.Release.Version)
	patched, err := (ksupatch.Service{Out: s.Out}).Run(ctx, ksupatch.Options{
		Boot: prepared.ImagePath, KMI: kmi, Mode: o.Mode, Kernel: o.Kernel, OutputDir: patchedDir, KSUd: o.KSUd,
	})
	if err != nil {
		return Result{}, err
	}
	result := Result{StockImage: prepared.ImagePath, PatchedImage: patched, ROMPath: prepared.ROMPath, Partition: prepared.Partition}
	if !o.Yes {
		logger.Info("workflow.step.manual_ready", "step", 5, "total", 5, "path", patched, "partition", prepared.Partition)
		serialArg := quoteCommandArg(prepared.Device.Serial)
		logger.Info("workflow.manual_flash_command", "command", fmt.Sprintf("adb -s %s reboot bootloader", serialArg))
		logger.Info("workflow.manual_flash_command", "command", fmt.Sprintf("fastboot -s %s flash %s %s", serialArg, prepared.Partition, quoteCommandArg(patched)))
		logger.Info("workflow.manual_flash_command", "command", fmt.Sprintf("fastboot -s %s reboot", serialArg))
		slog.Info("workflow.completed", "partition", result.Partition, "rom", result.ROMPath, "patched", result.PatchedImage)
		return result, nil
	}
	logger.Info("workflow.step.reboot_fastboot", "step", 6, "total", 8)
	logger.Info("workflow.step.validate_flash", "step", 7, "total", 8)
	err = (flash.Service{Devices: s.Devices, Out: s.Out}).Run(ctx, flash.Options{
		Serial: o.Serial, Image: patched, StockBackup: prepared.ImagePath, Partition: prepared.Partition,
		Mode: o.Mode, Yes: o.Yes, NoReboot: o.NoReboot,
	})
	if err != nil {
		return Result{}, err
	}
	if o.NoReboot {
		logger.Info("workflow.step.flash_stay_fastboot", "step", 8, "total", 8)
	} else {
		logger.Info("workflow.step.flash_booting", "step", 8, "total", 8)
	}
	slog.Info("workflow.completed", "partition", result.Partition, "rom", result.ROMPath, "patched", result.PatchedImage)
	return result, nil
}

func quoteCommandArg(value string) string {
	if runtime.GOOS == "windows" {
		return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
	}
	return `'` + strings.ReplaceAll(value, `'`, `'"'"'`) + `'`
}
