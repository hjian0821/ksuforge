# KSUForge

[English](README.md) | **简体中文**

KSUForge 是一款使用 Go 和 Wails 构建的 KernelSU 镜像准备与安装工具，同时提供桌面 GUI 和命令行界面。它可以识别 Android 设备、下载与当前系统版本完全匹配的官方 Recovery ROM、提取并修补启动镜像，并在完成安全检查后通过 Fastboot 刷入设备。

目前默认支持 Xiaomi、Redmi 和 POCO 设备，项目内部保留了面向其他厂商的 profile 扩展能力。

> [!WARNING]
> 刷写 `boot` 或 `init_boot` 分区存在无法开机、数据丢失甚至设备变砖的风险。使用前请解锁 Bootloader、备份重要数据，并保留与当前 ROM **完全匹配**的原厂镜像。KSUForge 不会绕过 Bootloader 锁或设备安全校验，所有操作风险由使用者自行承担。

## 功能特性

- 提供跨平台桌面 GUI 和独立 CLI
- 检测 `adb`、`fastboot` 以及已连接设备
- 读取设备型号、代号、Android 版本、系统版本、当前 slot、内核版本和 KMI
- 自动查找与设备当前完整系统版本一致的 Xiaomi Recovery ROM
- 自动测速多个官方 OTA 镜像，支持并行分片下载、断点续传与失败切换
- 校验 ROM 的 MD5，以及启动镜像的 magic、大小和 SHA-256
- 从 ZIP 或 `payload.bin` 提取 `boot.img` / `init_boot.img`
- 支持 KernelSU LKM 和 GKI 模式，并按设备信息选择目标分区
- 自动复用已经校验通过的 ROM、原厂镜像和 patched 镜像
- 默认只生成 patched 镜像；仅在显式确认后执行刷写
- 刷写前复核 Bootloader 解锁状态、Fastboot slot 和原厂镜像备份
- 支持中文和英文界面及日志

## 工作流程

```text
连接 Android 设备
        ↓
读取设备、系统版本与 KMI
        ↓
下载并校验同版本 Recovery ROM
        ↓
提取 boot.img / init_boot.img
        ↓
使用 ksud 生成 KernelSU patched 镜像
        ↓
预览刷写命令，或确认后自动刷写
```

## 系统要求

运行 KSUForge 前，请准备：

- 已解锁 Bootloader 的受支持 Android 设备
- 可正常传输数据的 USB 线
- [Android SDK Platform-Tools](https://developer.android.com/tools/releases/platform-tools)，并确保 `adb` 和 `fastboot` 位于 `PATH`
- 设备已开启 USB 调试并授权当前电脑

GUI 和 CLI 发布包都已经内置受支持平台的 `ksud` 和 `payload-dumper`，通常不需要单独安装。使用其他版本时，可以通过 CLI 的 `--ksud` 参数或 GUI 高级设置指定外部 `ksud`。

## 安装

请从 [GitHub Releases](https://github.com/hjian0821/ksuforge/releases) 下载与系统和架构对应的压缩包：

- 普通用户下载 `ksuforge-gui_*` 桌面版。
- 只在需要终端或脚本调用时下载 `ksuforge-cli_*` 命令行版。
- 除非需要同时使用两种界面，否则不必下载两份。
- Apple Silicon Mac 选择 `darwin_arm64`，Intel Mac 选择 `darwin_amd64`。
- 大多数 Windows 和 Linux 电脑选择 `amd64`。

GUI 与 CLI 分开打包：

| 系统 | 架构 | GUI 压缩包 | CLI 压缩包 |
| --- | --- | --- | --- |
| macOS | ARM64 | `ksuforge-gui_<版本>_darwin_arm64.tar.gz` | `ksuforge-cli_<版本>_darwin_arm64.tar.gz` |
| macOS | AMD64 | `ksuforge-gui_<版本>_darwin_amd64.tar.gz` | `ksuforge-cli_<版本>_darwin_amd64.tar.gz` |
| Linux | ARM64 | `ksuforge-gui_<版本>_linux_arm64.tar.gz` | `ksuforge-cli_<版本>_linux_arm64.tar.gz` |
| Linux | AMD64 | `ksuforge-gui_<版本>_linux_amd64.tar.gz` | `ksuforge-cli_<版本>_linux_amd64.tar.gz` |
| Windows | AMD64 | `ksuforge-gui_<版本>_windows_amd64.zip` | `ksuforge-cli_<版本>_windows_amd64.zip` |

## 使用方法

### 桌面应用

1. 连接手机，并确认 `adb devices` 能够识别设备。
2. 启动 `ksuforge`。
3. 选择 KernelSU 模式和输出目录；一般保持默认的 LKM 和自动分区即可。
4. 检查识别到的机型、系统版本、KMI 和目标分区。
5. 开始任务，等待 ROM 下载、校验、镜像提取和 patch 完成。
6. 默认情况下，应用只生成 patched 镜像并显示手动刷写命令。
7. 如需自动刷写，启用“自动刷写”并在二次确认后继续。

默认产物目录为 `~/Downloads/KSUForge`。GUI 还可以设置下载并发数、自定义 HTTPS 镜像、ROM 索引地址、默认模式、目标分区、输出目录、界面语言和日志滚动行为。

### 命令行

先检查运行环境和设备连接：

```bash
ksuforge-cli doctor
ksuforge-cli devices
ksuforge-cli device-info
```

执行完整流程：

```bash
# 下载 ROM、提取并 patch；不会执行刷写
ksuforge-cli install

# 确认信息无误后，自动进入 Fastboot 并执行刷写
ksuforge-cli install --yes
```

指定设备、模式或输出目录：

```bash
ksuforge-cli install \
  --serial DEVICE_SERIAL \
  --mode lkm \
  --partition auto \
  --output-dir ./ksuforge-output
```

也可以分步运行：

```bash
# 1. 下载同版本 ROM 并提取原厂镜像
ksuforge-cli prepare --output-dir ./ksuforge-output --mode lkm

# 2. 修补原厂镜像
ksuforge-cli patch \
  --image ./ksuforge-output/stock/DEVICE/VERSION/init_boot.img \
  --kmi android14-6.1 \
  --output-dir ./ksuforge-output

# 3. 预演刷写，不会写入设备
ksuforge-cli flash \
  --image ./ksuforge-output/patched/DEVICE/VERSION/kernelsu-patched-init_boot.img \
  --stock-backup ./ksuforge-output/stock/DEVICE/VERSION/init_boot.img \
  --mode lkm

# 4. 核对输出后执行真实刷写
ksuforge-cli flash \
  --image ./ksuforge-output/patched/DEVICE/VERSION/kernelsu-patched-init_boot.img \
  --stock-backup ./ksuforge-output/stock/DEVICE/VERSION/init_boot.img \
  --mode lkm \
  --yes
```

查看完整命令帮助：

```bash
ksuforge-cli --help
```

使用 `--lang zh` 或 `--lang en` 可以临时指定 CLI 语言：

```bash
ksuforge-cli --lang en device-info
```

> [!NOTE]
> 自动 ROM 下载目前仅支持 Xiaomi、Redmi 和 POCO。找不到完全一致的系统版本时，程序会停止，不会使用“最新版本”替代。非 Xiaomi 设备的通用刷写路径仅用于适配开发和验证，需要显式传入 `--allow-other-oem`。

## 产物目录

KSUForge 按类型、设备代号和系统版本保存文件，避免不同设备或版本相互覆盖：

```text
~/Downloads/KSUForge/
├── rom/<设备代号>/<系统版本>/<recovery-rom>.zip
├── stock/<设备代号>/<系统版本>/boot.img
├── stock/<设备代号>/<系统版本>/init_boot.img
└── patched/<设备代号>/<系统版本>/kernelsu-patched-<分区>.img
```

下载中的文件会使用 `.part` 和 `.part.json` 保存分片状态，下次运行时可以继续下载。配置和日志保存在各系统的标准应用数据目录中，CLI 与 GUI 共用配置。

## 从源码构建

### 准备环境

- Go 1.25 或更高版本
- GNU Make
- [Wails v2.15.0](https://wails.io/docs/gettingstarted/installation/)
- 对应平台所需的 Wails 编译依赖

安装 Wails CLI：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.15.0
```

首次构建时，构建脚本会从官方 GitHub Release 下载匹配平台的 KernelSU `ksud` 和 `payload-dumper-go`，校验 SHA-256 后嵌入应用；已经存在且有效的文件会直接复用。

### 构建全部应用

```bash
make build
```

构建结果写入仓库根目录的 `bin/`：

- macOS：`bin/ksuforge.app`、`bin/ksuforge-cli`
- Linux：`bin/ksuforge`、`bin/ksuforge-cli`
- Windows：`bin/ksuforge.exe`、`bin/ksuforge-cli.exe`

### 分别构建 CLI 或 GUI

```bash
make cli
make gui
```

也可以直接构建 GUI：

```bash
cd cmd/ksuforge-gui
wails build -o ../../../../bin/ksuforge
```

Windows 下请将输出文件名改为 `../../../../bin/ksuforge.exe`。

### 测试

```bash
make test
```

该命令会运行全部 Go 测试和 `go vet`。

## 发布构建

在当前系统上分别生成 GUI 和 CLI 发布压缩包：

```bash
make release VERSION=v1.2.3
```

产物位于 `bin/`，文件名格式为：

```text
ksuforge-gui_<版本>_<系统>_<架构>.tar.gz
ksuforge-cli_<版本>_<系统>_<架构>.tar.gz

# Windows 使用 .zip，不使用 .tar.gz
```

在 macOS 上安装并启动 Docker Desktop 后，可以构建当前 macOS 架构、Linux AMD64 和 Windows AMD64 的发布包：

```bash
make release-all VERSION=v1.2.3
```

每次推送或合并 Pull Request 到 `main` 后，GitHub Actions 都会更新滚动的 `nightly` 预发布版本，其中包含 macOS、Linux 和 Windows 构建包以及 SHA-256 校验文件。`nightly` 中的产物始终对应 `main` 的最新提交。

推送符合语义化版本规范的 `v*` tag 后，则会创建独立的正式 GitHub Release，并自动生成发布说明：

```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

## 技术说明

- LKM 模式由官方 `ksud boot-patch` 使用内置模块完成修补。
- GKI 模式需要提供与目标设备 KMI、补丁级别和压缩格式匹配的内核文件。
- 标准 Fastboot 通常不能读取分区，因此原厂备份来自同版本官方 ROM，而不是从设备导出。
- A/B 设备只刷写 Fastboot 当前 slot，并会与 ADB 阶段检测到的 slot 交叉校验。
- OTA 更新后保留 Root 时，建议在重启前通过 KernelSU Manager 安装到未使用槽位。

更多 KernelSU 安装说明请参考 [KernelSU 官方文档](https://kernelsu.org/guide/installation.html)。

## 扩展设备支持

厂商适配位于 `internal/profile`。新增设备支持时，实现并注册 `Profile`，再根据需要补充设备代号校验、分区策略、Bootloader 状态解析和厂商特有的 Fastboot 流程即可；通用工作流无需复制。

## 许可证

KSUForge 使用 [Apache License 2.0](LICENSE) 发布。

KernelSU 和 payload-dumper-go 分别遵循各自的开源许可证，本项目构建过程中获取的第三方二进制文件不改变其原有许可条款。
