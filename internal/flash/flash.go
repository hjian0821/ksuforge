package flash

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/hjian0821/ksuforge/internal/device"
	"github.com/hjian0821/ksuforge/internal/i18n"
	bootimage "github.com/hjian0821/ksuforge/internal/image"
	"github.com/hjian0821/ksuforge/internal/logging"
	"github.com/hjian0821/ksuforge/internal/profile"
)

type Options struct {
	Serial, Image, StockBackup, Partition, Mode string
	Yes, NoReboot, AllowOtherOEM                bool
}

type Service struct {
	Devices device.Client
	Out     io.Writer
}

func (s Service) Run(ctx context.Context, o Options) error {
	if s.Out == nil {
		s.Out = io.Discard
	}
	logger := logging.WithWriter(s.Out, logging.Language(ctx))
	patched, err := bootimage.Inspect(o.Image)
	if err != nil {
		return i18n.WrapError("error.flash_image_check", err)
	}
	stock, err := bootimage.Inspect(o.StockBackup)
	if err != nil {
		return i18n.WrapError("error.stock_backup_check", err)
	}
	if patched.SHA256 == stock.SHA256 {
		return i18n.NewError("error.flash_images_identical")
	}

	devices, err := s.Devices.List(ctx)
	if err != nil {
		return err
	}
	d, err := selectDevice(devices, o.Serial)
	if err != nil {
		return err
	}
	if d.Mode != device.ADB {
		return i18n.NewError("error.fastboot_preexisting")
	}
	d, err = s.Devices.Inspect(ctx, d.Serial)
	if err != nil {
		return err
	}
	p, supported := profile.Match(d.Manufacturer, d.Brand)
	if !supported && !o.AllowOtherOEM {
		return i18n.NewError("error.oem_unsupported")
	}
	profileName := "generic"
	if supported {
		profileName = p.Name()
		if err := p.ValidateCodename(d.Codename); err != nil {
			return err
		}
	}
	partition, err := ChooseDevicePartition(o.Partition, o.Mode, d.KMI, d.HasInitBoot, d.SDK)
	if err != nil {
		return err
	}

	logger.Info("flash.device_information", "model", d.Model, "codename", d.Codename, "profile", profileName, "android", d.Android, "sdk", d.SDK)
	logger.Info("flash.target", "partition", partition, "image", patched.Path, "sha256", patched.SHA256)
	logger.Info("flash.stock_backup", "image", stock.Path, "sha256", stock.SHA256)
	if !o.Yes {
		logger.Info("flash.preview_completed")
		return nil
	}

	if err := s.Devices.RebootBootloader(ctx, d.Serial); err != nil {
		return i18n.WrapError("error.reboot_bootloader", err)
	}
	deadline := time.Now().Add(45 * time.Second)
	for !s.Devices.IsFastboot(ctx, d.Serial) {
		if time.Now().After(deadline) {
			return i18n.NewError("error.fastboot_timeout")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
	unlocked, detail := s.Devices.BootloaderUnlocked(ctx, d.Serial)
	if !unlocked {
		return i18n.NewError("error.bootloader_locked", detail)
	}
	if currentSlot := s.Devices.CurrentSlot(ctx, d.Serial); d.Slot != "" && currentSlot != "" && currentSlot != d.Slot {
		return i18n.NewError("error.slot_mismatch", d.Slot, currentSlot)
	}
	if err := s.Devices.Flash(ctx, d.Serial, partition, patched.Path); err != nil {
		return i18n.WrapError("error.flash_partition", err, partition)
	}
	logger.Info("flash.completed")
	if !o.NoReboot {
		if err := s.Devices.Reboot(ctx, d.Serial); err != nil {
			return i18n.WrapError("error.reboot_after_flash", err)
		}
		logger.Info("flash.rebooting")
	}
	slog.Info("flash.workflow_completed", "serial", d.Serial, "partition", partition, "reboot", !o.NoReboot)
	return nil
}

func selectDevice(ds []device.Device, serial string) (device.Device, error) {
	if serial != "" {
		for _, d := range ds {
			if d.Serial == serial {
				return d, nil
			}
		}
		return device.Device{}, i18n.NewError("error.device_not_found", serial)
	}
	if len(ds) == 0 {
		return device.Device{}, i18n.NewError("error.device_no_transport")
	}
	if len(ds) > 1 {
		return device.Device{}, i18n.NewError("error.device_multiple")
	}
	return ds[0], nil
}

func ChoosePartition(partition, mode string, sdk int) (string, error) {
	if partition == "boot" || partition == "init_boot" {
		return partition, nil
	}
	if partition != "auto" {
		return "", i18n.NewError("error.partition_invalid")
	}
	switch mode {
	case "gki":
		return "boot", nil
	case "lkm":
		if sdk >= 33 {
			return "init_boot", nil
		}
		return "boot", nil
	default:
		return "", i18n.NewError("error.mode_invalid")
	}
}

func ChooseDevicePartition(partition, mode, kmi string, hasInitBoot bool, sdk int) (string, error) {
	if partition != "auto" {
		return ChoosePartition(partition, mode, sdk)
	}
	if mode == "gki" {
		return "boot", nil
	}
	if mode != "lkm" {
		return "", i18n.NewError("error.mode_invalid")
	}
	if hasInitBoot && !strings.HasPrefix(kmi, "android12-") {
		return "init_boot", nil
	}
	return "boot", nil
}
