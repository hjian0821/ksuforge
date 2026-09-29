package device

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hjian0821/ksuforge/internal/command"
	"github.com/hjian0821/ksuforge/internal/i18n"
)

type Mode string

const (
	ADB      Mode = "adb"
	Fastboot Mode = "fastboot"
)

type Device struct {
	Serial       string
	Mode         Mode
	Manufacturer string
	Brand        string
	Model        string
	Codename     string
	Android      string
	Version      string
	Slot         string
	Kernel       string
	KMI          string
	HasInitBoot  bool
	SDK          int
}

type Client struct{ Runner command.Runner }

func SelectADB(devices []Device, serial string) (Device, error) {
	var selected *Device
	for i := range devices {
		if serial != "" && devices[i].Serial != serial {
			continue
		}
		if devices[i].Mode != ADB {
			continue
		}
		if selected != nil {
			return Device{}, i18n.NewError("error.device_multiple_adb")
		}
		selected = &devices[i]
	}
	if selected == nil {
		return Device{}, i18n.NewError("error.device_no_adb")
	}
	return *selected, nil
}

func (c Client) List(ctx context.Context) ([]Device, error) {
	var found []Device
	if _, err := c.Runner.LookPath("adb"); err == nil {
		out, err := c.Runner.Run(ctx, "adb", "devices")
		if err == nil {
			s := bufio.NewScanner(strings.NewReader(out))
			for s.Scan() {
				f := strings.Fields(s.Text())
				if len(f) >= 2 && f[1] == "device" {
					found = append(found, Device{Serial: f[0], Mode: ADB})
				}
			}
		}
	}
	if _, err := c.Runner.LookPath("fastboot"); err == nil {
		out, err := c.Runner.Run(ctx, "fastboot", "devices")
		if err == nil {
			for _, line := range strings.Split(out, "\n") {
				f := strings.Fields(line)
				if len(f) > 0 && !has(found, f[0]) {
					found = append(found, Device{Serial: f[0], Mode: Fastboot})
				}
			}
		}
	}
	return found, nil
}

func has(devices []Device, serial string) bool {
	for _, d := range devices {
		if d.Serial == serial {
			return true
		}
	}
	return false
}

func (c Client) Inspect(ctx context.Context, serial string) (Device, error) {
	d := Device{Serial: serial, Mode: ADB}
	props := []struct {
		key  string
		dest *string
	}{
		{"ro.product.manufacturer", &d.Manufacturer}, {"ro.product.brand", &d.Brand},
		{"ro.product.model", &d.Model}, {"ro.product.device", &d.Codename},
		{"ro.build.version.release", &d.Android},
		{"ro.build.version.incremental", &d.Version},
	}
	for _, p := range props {
		v, err := c.Runner.Run(ctx, "adb", "-s", serial, "shell", "getprop", p.key)
		if err != nil {
			return Device{}, i18n.WrapError("error.device_property", err, p.key)
		}
		*p.dest = strings.TrimSpace(v)
	}
	sdk, err := c.Runner.Run(ctx, "adb", "-s", serial, "shell", "getprop", "ro.build.version.sdk")
	if err != nil {
		return Device{}, i18n.WrapError("error.device_sdk", err)
	}
	d.SDK, _ = strconv.Atoi(strings.TrimSpace(sdk))
	slot, _ := c.Runner.Run(ctx, "adb", "-s", serial, "shell", "getprop", "ro.boot.slot_suffix")
	d.Slot = strings.TrimPrefix(strings.TrimSpace(slot), "_")
	d.Kernel, err = c.Runner.Run(ctx, "adb", "-s", serial, "shell", "uname", "-r")
	if err != nil {
		return Device{}, i18n.WrapError("error.device_kernel", err)
	}
	d.Kernel = strings.TrimSpace(d.Kernel)
	d.KMI = ParseKMI(d.Kernel)
	initBoot := "/dev/block/by-name/init_boot"
	if d.Slot != "" {
		initBoot += "_" + d.Slot
	}
	if _, err := c.Runner.Run(ctx, "adb", "-s", serial, "shell", "ls", initBoot); err == nil {
		d.HasInitBoot = true
	}
	return d, nil
}

func ParseKMI(kernel string) string {
	versionEnd := strings.IndexByte(kernel, '-')
	version := kernel
	if versionEnd >= 0 {
		version = kernel[:versionEnd]
	}
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return ""
	}
	androidAt := strings.Index(kernel, "android")
	if androidAt < 0 {
		return ""
	}
	android := kernel[androidAt+len("android"):]
	end := 0
	for end < len(android) && android[end] >= '0' && android[end] <= '9' {
		end++
	}
	if end == 0 {
		return ""
	}
	return "android" + android[:end] + "-" + parts[0] + "." + parts[1]
}

func (c Client) RebootBootloader(ctx context.Context, serial string) error {
	_, err := c.Runner.Run(ctx, "adb", "-s", serial, "reboot", "bootloader")
	return err
}

func (c Client) IsFastboot(ctx context.Context, serial string) bool {
	out, err := c.Runner.Run(ctx, "fastboot", "devices")
	return err == nil && hasSerialLine(out, serial)
}

func hasSerialLine(out, serial string) bool {
	for _, line := range strings.Split(out, "\n") {
		f := strings.Fields(line)
		if len(f) > 0 && f[0] == serial {
			return true
		}
	}
	return false
}

func (c Client) GetVar(ctx context.Context, serial, key string) string {
	out, _ := c.Runner.Run(ctx, "fastboot", "-s", serial, "getvar", key)
	return out
}

// BootloaderUnlocked requires an explicit positive signal. fastboot often writes
// getvar values to stderr, and unsupported variables can include misleading text
// such as "not found", so substring matching is intentionally avoided.
func (c Client) BootloaderUnlocked(ctx context.Context, serial string) (bool, string) {
	unlockedRaw := c.GetVar(ctx, serial, "unlocked")
	secureRaw := c.GetVar(ctx, serial, "secure")
	unlocked, unlockedOK := parseGetVar(unlockedRaw, "unlocked")
	secure, secureOK := parseGetVar(secureRaw, "secure")
	ok := (unlockedOK && (unlocked == "yes" || unlocked == "true" || unlocked == "1")) ||
		(secureOK && (secure == "no" || secure == "false" || secure == "0"))
	return ok, fmt.Sprintf("unlocked=%q, secure=%q", unlocked, secure)
}

func (c Client) CurrentSlot(ctx context.Context, serial string) string {
	raw := c.GetVar(ctx, serial, "current-slot")
	slot, ok := parseGetVar(raw, "current-slot")
	if !ok {
		return ""
	}
	return strings.TrimPrefix(slot, "_")
}

func parseGetVar(output, key string) (string, bool) {
	key = strings.ToLower(strings.TrimSpace(key))
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(strings.ToLower(line))
		line = strings.TrimPrefix(line, "(bootloader)")
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			value := strings.Fields(strings.TrimSpace(parts[1]))
			if len(value) > 0 {
				return value[0], true
			}
		}
	}
	return "", false
}

func (c Client) Flash(ctx context.Context, serial, partition, image string) error {
	_, err := c.Runner.Run(ctx, "fastboot", "-s", serial, "flash", partition, image)
	return err
}

func (c Client) Reboot(ctx context.Context, serial string) error {
	_, err := c.Runner.Run(ctx, "fastboot", "-s", serial, "reboot")
	return err
}
