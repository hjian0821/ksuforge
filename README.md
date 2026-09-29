# KSUForge

**English** | [简体中文](README_zh-CN.md)

KSUForge is a KernelSU image preparation and installation tool built with Go and Wails. It provides both a desktop GUI and a command-line interface. KSUForge can identify an Android device, download the official Recovery ROM that exactly matches its current system version, extract and patch the boot image, and flash it through Fastboot after performing safety checks.

Xiaomi, Redmi, and POCO devices are supported by default. The project also provides a profile-based extension point for other manufacturers.

> [!WARNING]
> Flashing the `boot` or `init_boot` partition can cause boot loops, data loss, or a bricked device. Unlock the bootloader, back up important data, and keep a stock image that **exactly matches** the installed ROM before proceeding. KSUForge does not bypass bootloader locks or device safety checks. You use this software at your own risk.

## Features

- Cross-platform desktop GUI and standalone CLI
- Detection of `adb`, `fastboot`, and connected devices
- Device model, codename, Android version, build version, active slot, kernel version, and KMI detection
- Automatic lookup of a Xiaomi Recovery ROM that exactly matches the installed build
- Automatic benchmarking of official OTA mirrors, with parallel downloads, resume support, and mirror failover
- ROM MD5 verification and boot image magic, size, and SHA-256 validation
- Extraction of `boot.img` or `init_boot.img` from ZIP archives and `payload.bin`
- KernelSU LKM and GKI modes with automatic target partition selection
- Reuse of previously verified ROMs, stock images, and patched images
- Dry-run behavior by default; flashing requires explicit confirmation
- Bootloader state, Fastboot slot, and stock backup checks before flashing
- Chinese and English UI and logs

## Workflow

```text
Connect an Android device
        ↓
Read device information, build version, and KMI
        ↓
Download and verify the matching Recovery ROM
        ↓
Extract boot.img or init_boot.img
        ↓
Create a KernelSU-patched image with ksud
        ↓
Preview the flash command or confirm automatic flashing
```

## Requirements

Before running KSUForge, make sure you have:

- A supported Android device with an unlocked bootloader
- A USB cable that supports data transfer
- [Android SDK Platform-Tools](https://developer.android.com/tools/releases/platform-tools), with `adb` and `fastboot` available in `PATH`
- USB debugging enabled and the computer authorized on the device

Both the GUI and CLI packages contain `ksud` and `payload-dumper` for supported platforms, so they normally do not need to be installed separately. To use another `ksud` version, specify it with the CLI `--ksud` option or in the GUI advanced settings.

## Installation

Download the latest package for your operating system and architecture from [GitHub Releases](https://github.com/hjian0821/ksuforge/releases):

- Download a `ksuforge-gui_*` archive for the desktop application.
- Download a `ksuforge-cli_*` archive only if you want the command-line tool.
- Do not download both unless you intend to use both interfaces.
- Apple Silicon Macs use `darwin_arm64`; Intel Macs use `darwin_amd64`.
- Most Windows and Linux PCs use `amd64`.

GUI and CLI are packaged separately:

| Operating system | Architecture | GUI archive | CLI archive |
| --- | --- | --- | --- |
| macOS | ARM64 | `ksuforge-gui_<version>_darwin_arm64.tar.gz` | `ksuforge-cli_<version>_darwin_arm64.tar.gz` |
| macOS | AMD64 | `ksuforge-gui_<version>_darwin_amd64.tar.gz` | `ksuforge-cli_<version>_darwin_amd64.tar.gz` |
| Linux | ARM64 | `ksuforge-gui_<version>_linux_arm64.tar.gz` | `ksuforge-cli_<version>_linux_arm64.tar.gz` |
| Linux | AMD64 | `ksuforge-gui_<version>_linux_amd64.tar.gz` | `ksuforge-cli_<version>_linux_amd64.tar.gz` |
| Windows | AMD64 | `ksuforge-gui_<version>_windows_amd64.zip` | `ksuforge-cli_<version>_windows_amd64.zip` |

## Usage

### Desktop application

1. Connect your device and confirm that it appears in `adb devices`.
2. Start `ksuforge`.
3. Select the KernelSU mode and output directory. The default LKM mode and automatic partition selection are recommended for most users.
4. Review the detected model, system version, KMI, and target partition.
5. Start the task and wait for the ROM download, verification, image extraction, and patching to finish.
6. By default, KSUForge only creates the patched image and displays the manual flash command.
7. To flash automatically, enable automatic flashing and accept the second confirmation prompt.

The default artifact directory is `~/Downloads/KSUForge`. The GUI also lets you configure download concurrency, custom HTTPS mirrors, the ROM index URL, default mode and partition, output directory, interface language, and log scrolling.

### Command line

Check the environment and connected devices first:

```bash
ksuforge-cli doctor
ksuforge-cli devices
ksuforge-cli device-info
```

Run the complete workflow:

```bash
# Download the ROM, extract the image, and patch it without flashing
ksuforge-cli install

# Reboot into Fastboot and flash after all checks pass
ksuforge-cli install --yes
```

Select a device, mode, partition, or output directory:

```bash
ksuforge-cli install \
  --serial DEVICE_SERIAL \
  --mode lkm \
  --partition auto \
  --output-dir ./ksuforge-output
```

The workflow can also be run one step at a time:

```bash
# 1. Download the matching ROM and extract the stock image
ksuforge-cli prepare --output-dir ./ksuforge-output --mode lkm

# 2. Patch the stock image
ksuforge-cli patch \
  --image ./ksuforge-output/stock/DEVICE/VERSION/init_boot.img \
  --kmi android14-6.1 \
  --output-dir ./ksuforge-output

# 3. Preview the flash operation without writing to the device
ksuforge-cli flash \
  --image ./ksuforge-output/patched/DEVICE/VERSION/kernelsu-patched-init_boot.img \
  --stock-backup ./ksuforge-output/stock/DEVICE/VERSION/init_boot.img \
  --mode lkm

# 4. Flash after reviewing the preview
ksuforge-cli flash \
  --image ./ksuforge-output/patched/DEVICE/VERSION/kernelsu-patched-init_boot.img \
  --stock-backup ./ksuforge-output/stock/DEVICE/VERSION/init_boot.img \
  --mode lkm \
  --yes
```

Display the complete CLI help:

```bash
ksuforge-cli --help
```

Use `--lang zh` or `--lang en` before the command to override the configured CLI language:

```bash
ksuforge-cli --lang en device-info
```

> [!NOTE]
> Automatic ROM downloads currently support Xiaomi, Redmi, and POCO devices only. KSUForge stops if it cannot find an exact build match; it never substitutes the latest available ROM. The generic flashing path for other manufacturers is intended for adapter development and validation and requires the explicit `--allow-other-oem` option.

## Artifact Layout

KSUForge stores files by type, device codename, and build version so that artifacts from different devices or releases do not overwrite one another:

```text
~/Downloads/KSUForge/
├── rom/<codename>/<build-version>/<recovery-rom>.zip
├── stock/<codename>/<build-version>/boot.img
├── stock/<codename>/<build-version>/init_boot.img
└── patched/<codename>/<build-version>/kernelsu-patched-<partition>.img
```

Downloads in progress use `.part` and `.part.json` files to preserve chunk state and can resume during the next run. Configuration and logs are stored in the standard application data directories for each operating system. The CLI and GUI share the same configuration.

## Building from Source

### Prerequisites

- Go 1.25 or later
- GNU Make
- [Wails v2.15.0](https://wails.io/docs/gettingstarted/installation/)
- The Wails build dependencies for your operating system

Install the Wails CLI:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.15.0
```

On the first build, the build scripts download the appropriate KernelSU `ksud` and `payload-dumper-go` binaries from their official GitHub Releases, verify their SHA-256 checksums, and embed them in the applications. Existing verified files are reused.

### Build all applications

```bash
make build
```

Build outputs are written to the repository's `bin/` directory:

- macOS: `bin/ksuforge.app` and `bin/ksuforge-cli`
- Linux: `bin/ksuforge` and `bin/ksuforge-cli`
- Windows: `bin/ksuforge.exe` and `bin/ksuforge-cli.exe`

### Build the CLI or GUI separately

```bash
make cli
make gui
```

You can also build the GUI directly:

```bash
cd cmd/ksuforge-gui
wails build -o ../../../../bin/ksuforge
```

On Windows, use `../../../../bin/ksuforge.exe` as the output path.

### Run tests

```bash
make test
```

This command runs all Go tests followed by `go vet`.

## Release Builds

Create separate GUI and CLI release archives for the current target:

```bash
make release VERSION=v1.2.3
```

Archives are written to `bin/` using the following naming scheme:

```text
ksuforge-gui_<version>_<os>_<architecture>.tar.gz
ksuforge-cli_<version>_<os>_<architecture>.tar.gz

# Windows uses .zip instead of .tar.gz
```

On macOS, with Docker Desktop installed and running, build packages for the current macOS architecture, Linux AMD64, and Windows AMD64 with:

```bash
make release-all VERSION=v1.2.3
```

Every push or merged pull request to `main` updates the rolling `nightly` pre-release with macOS, Linux, and Windows packages plus SHA-256 checksums. The `nightly` assets always represent the latest commit on `main`.

Pushing a semantic `v*` tag creates a separate stable GitHub Release with generated release notes:

```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

## Technical Notes

- LKM mode uses the official `ksud boot-patch` command and its bundled module.
- GKI mode requires a kernel that matches the device KMI, security patch level, and compression format.
- Standard Fastboot cannot normally read partitions, so the stock backup comes from the matching official ROM rather than being dumped from the device.
- On A/B devices, KSUForge flashes only the current Fastboot slot and cross-checks it against the slot detected through ADB.
- To retain root after an OTA update, use KernelSU Manager to install to the inactive slot before rebooting.

For more information, see the [official KernelSU installation guide](https://kernelsu.org/guide/installation.html).

## Extending Device Support

Manufacturer adapters live in `internal/profile`. To add support for another device family, implement and register a `Profile`, then add any required codename validation, partition rules, bootloader status parsing, or manufacturer-specific Fastboot behavior. The shared workflow does not need to be duplicated.

## License

KSUForge is released under the [Apache License 2.0](LICENSE).

KernelSU and payload-dumper-go remain subject to their respective open-source licenses. Fetching or packaging their binaries does not change their original license terms.
