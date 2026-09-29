package cli

import (
	"context"
	"flag"
	"io"
	"os"
	"strings"

	"github.com/hjian0821/ksuforge/internal/command"
	"github.com/hjian0821/ksuforge/internal/config"
	"github.com/hjian0821/ksuforge/internal/device"
	"github.com/hjian0821/ksuforge/internal/flash"
	"github.com/hjian0821/ksuforge/internal/i18n"
	bootimage "github.com/hjian0821/ksuforge/internal/image"
	"github.com/hjian0821/ksuforge/internal/logging"
	ksupatch "github.com/hjian0821/ksuforge/internal/patch"
	appPaths "github.com/hjian0821/ksuforge/internal/paths"
	"github.com/hjian0821/ksuforge/internal/payloaddumper"
	"github.com/hjian0821/ksuforge/internal/prepare"
	"github.com/hjian0821/ksuforge/internal/rom"
	"github.com/hjian0821/ksuforge/internal/workflow"
)

// defaultOutputDir returns the remembered artifact directory, falling back to
// the platform default. Large artifacts stay in a user-visible location.
func defaultOutputDir() string {
	return config.Load().OutputDirOrDefault()
}

// defaultCatalog returns the remembered ROM index template.
func defaultCatalog() string {
	return config.Load().CatalogURLOr(rom.DefaultCatalogURL)
}

const usage = `ksuforge-cli - 安全、可扩展的 KernelSU 刷写工具

用法:
	ksuforge-cli [--lang zh|en] COMMAND
  ksuforge-cli doctor
  ksuforge-cli devices
  ksuforge-cli device-info [--serial ID]
  ksuforge-cli resolve-rom [--serial ID]
  ksuforge-cli inspect --image FILE
  ksuforge-cli prepare [--output-dir DIR] [选项]
  ksuforge-cli patch --image STOCK.img --kmi KMI [选项]
  ksuforge-cli flash --image PATCHED.img --stock-backup STOCK.img [选项]
  ksuforge-cli install [--output-dir DIR] [--yes]

prepare 选项:
  --serial ID             指定设备；仅有一台时可省略
  --mode lkm|gki          KernelSU 模式（默认 lkm）
  --partition auto|boot|init_boot
  --kmi KMI               覆盖自动识别的 KMI
  --output-dir DIR        产物根目录，按 rom|stock|patched/<代号>/<版本>/ 分层（默认记住的输出目录，首次 ~/Downloads/KSUForge）

install 选项:
  --serial ID             指定设备；仅有一台时可省略
  --mode lkm|gki          KernelSU 模式（默认 lkm）
  --partition auto|boot|init_boot
  --output-dir DIR        产物根目录，按 rom|stock|patched/<代号>/<版本>/ 分层
  --ksud FILE             覆盖随程序集成的 ksud
  --catalog URL           ROM 索引地址模板（含设备代号占位符）
  --kmi KMI               覆盖自动识别的 KMI
  --kernel FILE           GKI 模式使用的内核文件
  --yes                   自动重启 fastboot 并执行写入
  --no-reboot             写入后不重启

flash 选项:
  --serial ID             指定设备；仅有一台时可省略
  --mode lkm|gki          KernelSU 模式（默认 lkm）
  --partition auto|boot|init_boot
  --yes                   执行写入；省略时仅预演
  --no-reboot             写入后不重启
  --allow-other-oem       允许使用通用适配器处理非小米设备
`

func Run(ctx context.Context, args []string, out, errOut io.Writer) int {
	language := config.Load().LanguageOrDefault(i18n.Detect())
	var err error
	args, language, err = parseLanguage(args, language)
	ctx = logging.WithLanguage(ctx, language)
	if err != nil {
		logging.WithWriter(errOut, language).Error("app.command_failed", "error", err)
		return 1
	}
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" {
		_, _ = io.WriteString(out, usage)
		return 0
	}
	logger := logging.WithWriter(out, language)
	errLogger := logging.WithWriter(errOut, language)
	runner := command.ExecRunner{SearchPaths: command.DesktopSearchPaths()}
	client := device.Client{Runner: runner}
	switch args[0] {
	case "doctor":
		for _, tool := range []string{"adb", "fastboot"} {
			path, e := runner.LookPath(tool)
			if e != nil {
				err = i18n.NewError("error.tool_required", tool)
				break
			}
			logger.Info("tool.available", "tool", tool, "path", path)
		}
		if err == nil {
			if path, e := ksupatch.ResolveKSUd(""); e == nil {
				logger.Info("tool.available", "tool", "ksud", "path", path)
			} else {
				logger.Warn("tool.missing_optional", "tool", "ksud")
			}
			if path, e := payloaddumper.Resolve(); e == nil {
				logger.Info("tool.available", "tool", "payload-dumper", "path", path)
			} else {
				logger.Warn("tool.missing_optional", "tool", "payload-dumper")
			}
		}
	case "devices":
		var ds []device.Device
		ds, err = client.List(ctx)
		if err == nil && len(ds) == 0 {
			logger.Info("device.none")
		}
		for _, d := range ds {
			logger.Info("device.found", "serial", d.Serial, "mode", d.Mode)
		}
	case "device-info":
		err = runDeviceInfo(ctx, client, args[1:], out, errOut)
	case "resolve-rom":
		err = runResolveROM(ctx, client, args[1:], out, errOut)
	case "inspect":
		err = inspect(ctx, args[1:], out, errOut)
	case "prepare":
		err = runPrepare(ctx, client, args[1:], out, errOut)
	case "patch":
		err = runPatch(ctx, args[1:], out, errOut)
	case "flash":
		err = runFlash(ctx, client, args[1:], out, errOut)
	case "install":
		err = runInstall(ctx, client, args[1:], out, errOut)
	default:
		err = i18n.NewError("error.command_unknown", args[0])
	}
	if err != nil {
		errLogger.Error("app.command_failed", "error", err)
		return 1
	}
	return 0
}

func parseLanguage(args []string, fallback string) ([]string, string, error) {
	if len(args) == 0 {
		return args, fallback, nil
	}
	if strings.HasPrefix(args[0], "--lang=") {
		value := strings.TrimPrefix(args[0], "--lang=")
		if value != "zh" && value != "en" {
			return nil, fallback, i18n.NewError("error.language_unsupported")
		}
		return args[1:], value, nil
	}
	if args[0] == "--lang" {
		if len(args) < 2 || (args[1] != "zh" && args[1] != "en") {
			return nil, fallback, i18n.NewError("error.language_unsupported")
		}
		return args[2:], args[1], nil
	}
	return args, fallback, nil
}

func runResolveROM(ctx context.Context, client device.Client, args []string, out, errOut io.Writer) error {
	fs := flag.NewFlagSet("resolve-rom", flag.ContinueOnError)
	fs.SetOutput(errOut)
	serial := fs.String("serial", "", "device serial")
	catalog := fs.String("catalog", defaultCatalog(), "ROM catalog URL template")
	mode := fs.String("mode", "lkm", "KernelSU mode")
	partitionFlag := fs.String("partition", "auto", "target partition")
	if err := fs.Parse(args); err != nil {
		return err
	}
	devices, err := client.List(ctx)
	if err != nil {
		return err
	}
	selected, err := device.SelectADB(devices, *serial)
	if err != nil {
		return err
	}
	d, err := client.Inspect(ctx, selected.Serial)
	if err != nil {
		return err
	}
	partition, err := flash.ChooseDevicePartition(*partitionFlag, *mode, d.KMI, d.HasInitBoot, d.SDK)
	if err != nil {
		return err
	}
	release, err := rom.FindRecovery(ctx, prepare.NewHTTPClient(), *catalog, d.Codename, d.Version)
	if err != nil {
		return err
	}
	logging.WithWriter(out, logging.Language(ctx)).Info("rom.match_result",
		"codename", d.Codename, "version", release.Version, "kmi", d.KMI,
		"mode", *mode, "partition", partition, "type", "Recovery",
		"size", release.Size, "md5", release.MD5, "url", release.URL,
	)
	return nil
}

func runDeviceInfo(ctx context.Context, client device.Client, args []string, out, errOut io.Writer) error {
	fs := flag.NewFlagSet("device-info", flag.ContinueOnError)
	fs.SetOutput(errOut)
	serial := fs.String("serial", "", "device serial")
	if err := fs.Parse(args); err != nil {
		return err
	}
	devices, err := client.List(ctx)
	if err != nil {
		return err
	}
	selected, err := device.SelectADB(devices, *serial)
	if err != nil {
		return err
	}
	d, err := client.Inspect(ctx, selected.Serial)
	if err != nil {
		return err
	}
	logging.WithWriter(out, logging.Language(ctx)).Info("device.information",
		"serial", d.Serial, "manufacturer", d.Manufacturer, "brand", d.Brand,
		"model", d.Model, "codename", d.Codename, "android", d.Android,
		"sdk", d.SDK, "version", d.Version, "slot", d.Slot, "kernel", d.Kernel,
		"kmi", d.KMI, "has_init_boot", d.HasInitBoot,
	)
	return nil
}

func runPatch(ctx context.Context, args []string, out, errOut io.Writer) error {
	fs := flag.NewFlagSet("patch", flag.ContinueOnError)
	fs.SetOutput(errOut)
	var o ksupatch.Options
	fs.StringVar(&o.Boot, "image", "", "stock boot image")
	fs.StringVar(&o.KMI, "kmi", "", "KMI, e.g. android14-6.1")
	fs.StringVar(&o.Mode, "mode", "lkm", "KernelSU mode")
	fs.StringVar(&o.Kernel, "kernel", "", "replacement GKI kernel")
	fs.StringVar(&o.OutputDir, "output-dir", defaultOutputDir(), "output directory")
	fs.StringVar(&o.KSUd, "ksud", "", "ksud path override")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(o.Boot) == "" {
		return i18n.NewError("error.image_required")
	}
	o.OutputDir = appPaths.PatchedDirFor(o.OutputDir, o.Boot)
	_, err := (ksupatch.Service{Out: out}).Run(ctx, o)
	return err
}

func runInstall(ctx context.Context, client device.Client, args []string, out, errOut io.Writer) error {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(errOut)
	var o workflow.Options
	var mirrors string
	fs.StringVar(&o.Serial, "serial", "", "device serial")
	fs.StringVar(&o.Mode, "mode", "lkm", "KernelSU mode")
	fs.StringVar(&o.Partition, "partition", "auto", "target partition")
	fs.StringVar(&o.OutputDir, "output-dir", defaultOutputDir(), "output directory")
	fs.StringVar(&o.KMI, "kmi", "", "override detected KMI")
	fs.StringVar(&o.Kernel, "kernel", "", "replacement GKI kernel")
	fs.StringVar(&o.KSUd, "ksud", "", "ksud path override")
	fs.StringVar(&o.Catalog, "catalog", defaultCatalog(), "ROM catalog URL template")
	fs.StringVar(&mirrors, "mirrors", "", "comma-separated HTTPS ROM mirror base URLs")
	fs.IntVar(&o.Connections, "connections", 16, "parallel ROM download connections (1-32)")
	fs.BoolVar(&o.Yes, "yes", false, "perform flash")
	fs.BoolVar(&o.NoReboot, "no-reboot", false, "do not reboot")
	if err := fs.Parse(args); err != nil {
		return err
	}
	o.Mirrors = splitList(mirrors)
	_, err := (workflow.Service{Devices: client, Out: out}).Run(ctx, o)
	return err
}

func runPrepare(ctx context.Context, client device.Client, args []string, out, errOut io.Writer) error {
	fs := flag.NewFlagSet("prepare", flag.ContinueOnError)
	fs.SetOutput(errOut)
	serial := fs.String("serial", "", "device serial")
	mode := fs.String("mode", "lkm", "KernelSU mode")
	partitionFlag := fs.String("partition", "auto", "target partition")
	outputDir := fs.String("output-dir", defaultOutputDir(), "output directory")
	kmi := fs.String("kmi", "", "override detected KMI")
	catalog := fs.String("catalog", defaultCatalog(), "ROM catalog URL template")
	mirrors := fs.String("mirrors", "", "comma-separated HTTPS ROM mirror base URLs")
	connections := fs.Int("connections", 16, "parallel ROM download connections (1-32)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, err := (prepare.Service{Devices: client, HTTPClient: prepare.NewHTTPClient(), Out: out}).Run(ctx, prepare.Options{
		Serial: *serial, Mode: *mode, Partition: *partitionFlag, OutputDir: *outputDir, Catalog: *catalog, KMI: *kmi,
		Mirrors: splitList(*mirrors), Connections: *connections,
	})
	return err
}

func splitList(value string) []string {
	return strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == '\n' || r == '\r' })
}

func inspect(ctx context.Context, args []string, out, errOut io.Writer) error {
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	fs.SetOutput(errOut)
	path := fs.String("image", "", "image path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return i18n.NewError("error.image_required")
	}
	i, err := bootimage.Inspect(*path)
	if err != nil {
		return err
	}
	logging.WithWriter(out, logging.Language(ctx)).Info("image.information", "path", i.Path, "size", i.Size, "sha256", i.SHA256)
	return nil
}

func runFlash(ctx context.Context, client device.Client, args []string, out, errOut io.Writer) error {
	fs := flag.NewFlagSet("flash", flag.ContinueOnError)
	fs.SetOutput(errOut)
	var o flash.Options
	fs.StringVar(&o.Serial, "serial", "", "device serial")
	fs.StringVar(&o.Image, "image", "", "patched image")
	fs.StringVar(&o.StockBackup, "stock-backup", "", "stock image backup")
	fs.StringVar(&o.Partition, "partition", "auto", "target partition")
	fs.StringVar(&o.Mode, "mode", "lkm", "KernelSU mode")
	fs.BoolVar(&o.Yes, "yes", false, "perform flash")
	fs.BoolVar(&o.NoReboot, "no-reboot", false, "do not reboot")
	fs.BoolVar(&o.AllowOtherOEM, "allow-other-oem", false, "allow generic OEM")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(o.Image) == "" || strings.TrimSpace(o.StockBackup) == "" {
		return i18n.NewError("error.image_and_backup_required")
	}
	if _, err := os.Stat(o.Image); err != nil {
		return err
	}
	return (flash.Service{Devices: client, Out: out}).Run(ctx, o)
}
