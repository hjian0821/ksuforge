package prepare

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hjian0821/ksuforge/internal/device"
	"github.com/hjian0821/ksuforge/internal/flash"
	"github.com/hjian0821/ksuforge/internal/i18n"
	"github.com/hjian0821/ksuforge/internal/logging"
	appPaths "github.com/hjian0821/ksuforge/internal/paths"
	"github.com/hjian0821/ksuforge/internal/profile"
	"github.com/hjian0821/ksuforge/internal/rom"
)

type Options struct {
	Serial, Mode, Partition, OutputDir, Catalog, KMI, ROMPath string
	RequireKMI, VerifyROM                                     bool
	Mirrors                                                   []string
	Connections                                               int
}

type Result struct {
	Device    device.Device
	Release   rom.Release
	Partition string
	ROMPath   string
	ImagePath string
}

type Service struct {
	Devices    device.Client
	HTTPClient *http.Client
	Out        io.Writer
}

func NewHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{ResponseHeaderTimeout: 30 * time.Second},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return i18n.NewError("error.rom_redirect_limit")
			}
			if len(via) > 0 && rom.ValidateOfficialURL(via[0].URL.String()) == nil {
				return rom.ValidateOfficialURL(req.URL.String())
			}
			if req.URL.Scheme != "https" || req.URL.Hostname() == "" || req.URL.User != nil {
				return i18n.NewError("error.rom_redirect_insecure")
			}
			return nil
		},
	}
}

func (s Service) Run(ctx context.Context, o Options) (Result, error) {
	if o.Mode == "" {
		o.Mode = "lkm"
	}
	if o.Partition == "" {
		o.Partition = "auto"
	}
	if o.OutputDir == "" {
		o.OutputDir = appPaths.DefaultOutputDir()
	}
	if o.Catalog == "" {
		o.Catalog = rom.DefaultCatalogURL
	}
	if s.HTTPClient == nil {
		s.HTTPClient = NewHTTPClient()
	}
	if s.Out == nil {
		s.Out = io.Discard
	}
	logger := logging.WithWriter(s.Out, logging.Language(ctx))
	devices, err := s.Devices.List(ctx)
	if err != nil {
		return Result{}, err
	}
	selected, err := device.SelectADB(devices, o.Serial)
	if err != nil {
		return Result{}, err
	}
	d, err := s.Devices.Inspect(ctx, selected.Serial)
	if err != nil {
		return Result{}, err
	}
	if _, supported := profile.Match(d.Manufacturer, d.Brand); !supported {
		return Result{}, i18n.NewError("error.rom_oem_unsupported")
	}
	if o.KMI != "" {
		d.KMI = o.KMI
	}
	if o.RequireKMI && d.KMI == "" {
		return Result{}, i18n.NewError("error.kmi_from_kernel", d.Kernel)
	}
	partition, err := flash.ChooseDevicePartition(o.Partition, o.Mode, d.KMI, d.HasInitBoot, d.SDK)
	if err != nil {
		return Result{}, err
	}
	logger.Info("device.target_image", "model", d.Model, "codename", d.Codename, "version", d.Version, "slot", d.Slot, "kmi", d.KMI, "image", partition+".img")
	release := rom.Release{Version: d.Version}
	romPath := strings.TrimSpace(o.ROMPath)
	if romPath != "" {
		romPath, err = filepath.Abs(romPath)
		if err != nil {
			return Result{}, i18n.WrapError("error.rom_local_path", err)
		}
		info, statErr := os.Stat(romPath)
		if statErr != nil {
			return Result{}, i18n.WrapError("error.rom_local_read", statErr)
		}
		if !info.Mode().IsRegular() || !strings.EqualFold(filepath.Ext(romPath), ".zip") {
			return Result{}, i18n.NewError("error.rom_local_not_zip")
		}
		logger.Info("rom.local_selected", "path", romPath)
		if o.VerifyROM {
			release, err = rom.FindRecovery(ctx, s.HTTPClient, o.Catalog, d.Codename, d.Version)
			if err != nil {
				return Result{}, err
			}
			if err := rom.VerifyFileMD5(romPath, release.MD5); err != nil {
				return Result{}, err
			}
			logger.Info("rom.local_md5_passed", "path", romPath)
		} else {
			logger.Warn("rom.local_md5_skipped", "path", romPath)
			// Keep extraction caching tied to this exact local file without
			// calculating a checksum the user explicitly chose to skip.
			release.MD5 = fmt.Sprintf("unverified:%d:%d:%s", info.Size(), info.ModTime().UnixNano(), romPath)
		}
	} else {
		release, err = rom.FindRecovery(ctx, s.HTTPClient, o.Catalog, d.Codename, d.Version)
		if err != nil {
			return Result{}, err
		}
		romDir := appPaths.ROMDir(o.OutputDir, d.Codename, release.Version)
		stockDir := appPaths.StockDir(o.OutputDir, d.Codename, release.Version)
		logger.Info("rom.match_found", "url", release.URL, "md5", release.MD5, "rom_dir", romDir, "stock_dir", stockDir)
		romPath, err = rom.DownloadWithOptions(ctx, s.HTTPClient, release, romDir, rom.DownloadOptions{
			Mirrors: o.Mirrors, Connections: o.Connections,
		}, s.Out)
		if err != nil {
			return Result{}, err
		}
		logger.Info("rom.verification_passed", "path", romPath)
	}
	stockDir := appPaths.StockDir(o.OutputDir, d.Codename, release.Version)
	imagePath, err := rom.ExtractImage(ctx, romPath, partition, stockDir, release, s.Out)
	if err != nil {
		return Result{}, err
	}
	logger.Info("rom.stock_extracted", "path", imagePath)
	slog.Info("rom.stock_preparation_completed", "serial", d.Serial, "partition", partition, "rom", romPath, "image", imagePath)
	return Result{Device: d, Release: release, Partition: partition, ROMPath: romPath, ImagePath: imagePath}, nil
}
