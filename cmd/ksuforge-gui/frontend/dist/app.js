const I18N = {
  zh: {
    tagline: "连接设备后，自动完成匹配、下载、提取、patch 与安全刷写。",
    "device.connected": "1 台设备已连接",
    "device.disconnected": "未检测到设备",
    "device.none": "未检测",
    "device.none.title": "尚未检测到设备",
    "device.none.hint": "连接手机并开启 USB 调试",
    "btn.detect": "检测设备",
    "log.idle": "等待检测设备…",
    "progress.step": "步骤",
    "progress.current": "当前",
    "step.usb": "USB / ADB",
    "step.detect": "检测设备",
    "step.info": "设备信息",
    "step.fetch": "获取原厂镜像",
    "step.patch": "KernelSU Patch",
    "step.ready": "准备完成",
    "step.fastboot": "Fastboot",
    "step.verify": "安全检查",
    "step.flash": "刷入",
    "step.reboot": "重启",
    "step.done": "完成",
    "config.title": "安装配置",
    "field.device": "ADB 设备",
    "btn.refresh": "刷新",
    "field.ota": "Recovery OTA 包",
    "hint.ota": "留空则自动匹配并下载",
    "hint.ota.auto": "自动匹配当前系统版本",
    "btn.clear": "清除",
    "ota.verify": "校验本地 OTA 的 MD5",
    "ota.verify.desc": "根据当前设备版本查询官方校验值，推荐保持开启。",
    "field.mode": "KernelSU 模式",
    "mode.recommended": "推荐",
    "mode.lkm": "LKM",
    "mode.gki": "GKI 替换内核",
    "mode.lkm.desc": "修补 init_boot / boot，保留原厂内核。",
    "mode.gki.desc": "替换 boot 中的内核，需自备匹配内核文件。",
    "field.partition": "目标分区",
    "partition.auto": "自动（依据 KMI 与分区）",
    "field.output": "输出目录",
    "hint.output": "默认 ~/Downloads/KSUForge",
    "btn.choose": "选择",
    "download.title": "ROM 下载源",
    "download.desc": "留空则自动测速小米官方 CDN；自定义源优先参与测速。",
    "download.resume": "支持断点续传",
    "field.mirrors": "镜像源基础地址",
    "hint.mirrors": "每行一个 HTTPS 地址",
    "field.connections": "并发连接数",
    "download.strategy": "并行分片 · 镜像故障切换 · MD5 校验",
    advanced: "高级选项",
    "field.ksud": "ksud 路径",
    "hint.ksud": "留空使用内置版本",
    "field.kmi": "KMI 覆盖",
    "hint.kmi": "通常留空",
    "field.kernel": "GKI 内核文件",
    "settings.title": "设置",
    "settings.short": "设置",
    "settings.general": "通用",
    "settings.install": "安装默认值",
    "settings.advanced": "下载与工具",
    "settings.language": "界面语言",
    "settings.logScroll": "日志自动滚动",
    "settings.logScroll.desc": "运行时自动滚动到最新日志。",
    "settings.autoFlash": "默认自动刷写",
    "settings.noReboot": "默认刷写后不重启",
    "settings.catalog": "备用 ROM 索引地址",
    "settings.catalog.desc":
      "优先查询小米官方；此地址用于回退。必须为 HTTPS，包含一个 %s 设备代号占位符。",
    "settings.ksud": "ksud 路径",
    "settings.chooseKsud": "选择 ksud",
    "btn.save": "保存",
    "btn.cancel": "取消",
    "toggle.flash": "下载并 patch 后自动刷写",
    "toggle.flash.desc": "会重启到 fastboot 并覆盖当前启动分区。",
    "toggle.noreboot": "刷写后不重启",
    "toggle.noreboot.desc": "停留在 fastboot，方便手动检查。",
    "btn.prepare": "仅准备镜像",
    "btn.start": "开始执行",
    "device.title": "设备信息",
    "kv.model": "机型",
    "kv.codename": "代号",
    "kv.android": "Android",
    "kv.version": "系统版本",
    "kv.slot": "当前 Slot",
    "kv.kmi": "KMI",
    "kv.kernel": "内核",
    "kv.initboot": "init_boot",
    "kv.present": "存在",
    "kv.absent": "不存在",
    "summary.title": "本次操作",
    "summary.mode": "模式",
    "summary.partition": "目标分区",
    "summary.rom": "匹配 ROM",
    "summary.flash": "刷写",
    "summary.flash.no": "否（仅准备）",
    "summary.flash.yes": "是（自动刷写）",
    "summary.flash.yes.noreboot": "是（刷写后停留在 fastboot）",
    "summary.download": "下载策略",
    "log.title": "运行日志",
    "log.autoScroll.on": "自动滚动",
    "log.autoScroll.off": "已暂停",
    "log.collapse": "折叠",
    "log.expand": "展开",
    "status.idle": "空闲",
    "status.running": "运行中",
    "status.succeeded": "完成",
    "status.failed": "失败",
    "foot.note": "真实刷写存在风险，请先解锁并备份。",
    "confirm.flash": "即将自动重启到 fastboot 并写入启动分区，确认继续？",
    "confirm.title": "确认操作",
    "confirm.ok": "确认",
    "confirm.cancel": "取消",
    "error.prefix": "错误: ",
  },
  en: {
    tagline:
      "Connect a device to auto-match, download, extract, patch and safely flash.",
    "device.connected": "1 device connected",
    "device.disconnected": "No device",
    "device.none": "N/A",
    "device.none.title": "No device detected",
    "device.none.hint": "Connect your phone and enable USB debugging",
    "btn.detect": "Detect Device",
    "log.idle": "Waiting for device…",
    "progress.step": "STEP",
    "progress.current": "CURRENT",
    "step.usb": "USB / ADB",
    "step.detect": "Detect",
    "step.info": "Device Info",
    "step.fetch": "Get Stock Image",
    "step.patch": "KernelSU Patch",
    "step.ready": "Ready",
    "step.fastboot": "Fastboot",
    "step.verify": "Safety Check",
    "step.flash": "Flash",
    "step.reboot": "Reboot",
    "step.done": "Done",
    "config.title": "Installation",
    "field.device": "ADB Device",
    "btn.refresh": "Refresh",
    "field.ota": "Recovery OTA",
    "hint.ota": "Leave empty to match and download",
    "hint.ota.auto": "Match the current system version automatically",
    "btn.clear": "Clear",
    "ota.verify": "Verify local OTA MD5",
    "ota.verify.desc":
      "Fetch the official checksum for the connected device version. Recommended.",
    "field.mode": "KernelSU Mode",
    "mode.recommended": "REC",
    "mode.lkm": "LKM",
    "mode.gki": "GKI Kernel",
    "mode.lkm.desc": "Patch init_boot / boot, keep the stock kernel.",
    "mode.gki.desc":
      "Replace the kernel in boot; a matching kernel file is required.",
    "field.partition": "Target Partition",
    "partition.auto": "Auto (by KMI & partition)",
    "field.output": "Output Directory",
    "hint.output": "Default ~/Downloads/KSUForge",
    "btn.choose": "Choose",
    "download.title": "ROM Sources",
    "download.desc":
      "Leave empty to benchmark Xiaomi CDNs; custom sources are tried first.",
    "download.resume": "Resume enabled",
    "field.mirrors": "Mirror base URLs",
    "hint.mirrors": "One HTTPS URL per line",
    "field.connections": "Connections",
    "download.strategy": "Parallel ranges · failover · MD5 verification",
    advanced: "Advanced",
    "field.ksud": "ksud Path",
    "hint.ksud": "Leave empty to use bundled",
    "field.kmi": "KMI Override",
    "hint.kmi": "Usually leave empty",
    "field.kernel": "GKI Kernel File",
    "settings.title": "Settings",
    "settings.short": "Settings",
    "settings.general": "General",
    "settings.install": "Install Defaults",
    "settings.advanced": "Download & Tools",
    "settings.language": "Language",
    "settings.logScroll": "Log auto-scroll",
    "settings.logScroll.desc": "Scroll to the newest log line while running.",
    "settings.autoFlash": "Default auto-flash",
    "settings.noReboot": "Default no-reboot",
    "settings.catalog": "Fallback ROM Index URL",
    "settings.catalog.desc":
      "Xiaomi official lookup runs first. Fallback URL must use HTTPS and a single %s device codename placeholder.",
    "settings.ksud": "ksud Path",
    "settings.chooseKsud": "Choose ksud",
    "btn.save": "Save",
    "btn.cancel": "Cancel",
    "toggle.flash": "Auto-flash after patching",
    "toggle.flash.desc":
      "Reboots to fastboot and overwrites the active boot partition.",
    "toggle.noreboot": "Don't reboot after flashing",
    "toggle.noreboot.desc": "Stay in fastboot for manual inspection.",
    "btn.prepare": "Prepare Only",
    "btn.start": "Start",
    "device.title": "Device",
    "kv.model": "Model",
    "kv.codename": "Codename",
    "kv.android": "Android",
    "kv.version": "Version",
    "kv.slot": "Active Slot",
    "kv.kmi": "KMI",
    "kv.kernel": "Kernel",
    "kv.initboot": "init_boot",
    "kv.present": "Present",
    "kv.absent": "Absent",
    "summary.title": "This Run",
    "summary.mode": "Mode",
    "summary.partition": "Partition",
    "summary.rom": "Matched ROM",
    "summary.flash": "Flash",
    "summary.flash.no": "No (prepare only)",
    "summary.flash.yes": "Yes (auto-flash)",
    "summary.flash.yes.noreboot": "Yes (stay in fastboot)",
    "summary.download": "Download",
    "log.title": "Log",
    "log.autoScroll.on": "Auto-scroll",
    "log.autoScroll.off": "Paused",
    "log.collapse": "Collapse",
    "log.expand": "Expand",
    "status.idle": "Idle",
    "status.running": "Running",
    "status.succeeded": "Done",
    "status.failed": "Failed",
    "foot.note": "Flashing is risky. Unlock the bootloader and back up first.",
    "confirm.flash":
      "This will reboot to fastboot and overwrite the active boot partition. Continue?",
    "confirm.title": "Confirm",
    "confirm.ok": "Confirm",
    "confirm.cancel": "Cancel",
    "error.prefix": "Error: ",
  },
};

Object.assign(I18N.zh, {
  "nav.workspace": "工作空间",
  "nav.overview": "镜像工作台",
  "nav.docs": "使用文档",
  "nav.local": "在本机运行",
  "workspace.title": "为你的设备，准备好下一步。",
  "workspace.verified": "同版匹配 · 完整校验",
  "workspace.flow": "操作流程",
  "config.desc": "选择来源与操作方式，其余交给 KSUForge。",
  "source.local": "选择本地包",
  "advanced.hint": "KMI 与 GKI 内核",
  "device.ready": "已连接",
  "action.connect": "连接设备后即可开始。",
  "action.ready": "仅生成修补镜像，不会写入设备。",
  "action.flash": "开始前需确认，准备完成后将自动刷写。",
  "action.running": "任务执行中，请保持 USB 连接。",
  "result.ready": "镜像已准备完成",
  "result.flashed": "刷写已完成",
  "result.desc": "原厂镜像与修补镜像已保存在输出目录。",
  "result.failed": "任务未完成",
  "result.copy": "复制路径",
  "result.copied": "已复制",
  "result.copyFailed": "复制失败，请手动选择路径",
  "btn.flash": "准备并自动刷写",
  "btn.busy": "正在执行…",
  "summary.safe": "默认仅生成镜像，不写入设备。",
  "summary.caution": "将写入设备，请先解锁并完成备份。",
  tagline: "匹配原厂镜像，完成 KernelSU 修补。每一步，清晰可控。",
  "config.title": "准备镜像",
  "field.device": "目标设备",
  "field.ota": "原厂 Recovery 包",
  "device.title": "设备概览",
  "device.none.title": "连接你的 Android 设备",
  "device.none.hint": "通过 USB 连接手机，并开启 USB 调试。",
  "summary.title": "操作摘要",
  "btn.start": "开始准备镜像",
  "summary.flash.no": "仅准备镜像",
  "summary.flash.yes": "自动刷写",
  "summary.flash.yes.noreboot": "刷写后停留在 Fastboot",
  "status.idle": "等待开始",
  "log.idle": "任务开始后，详细进度将在这里显示。",
  "foot.note": "刷写前，请确认 Bootloader 已解锁并完成备份。",
});
Object.assign(I18N.en, {
  "nav.workspace": "WORKSPACE",
  "nav.overview": "Image workspace",
  "nav.docs": "Documentation",
  "nav.local": "Running locally",
  "workspace.title": "Your device. Ready for what’s next.",
  "workspace.verified": "Exact build · Verified images",
  "workspace.flow": "Workflow",
  "config.desc": "Choose your source and operation. We’ll handle the rest.",
  "source.local": "Local package",
  "advanced.hint": "KMI & GKI kernel",
  "device.ready": "Connected",
  "action.connect": "Connect a device to get started.",
  "action.ready": "Creates a patched image without writing to your device.",
  "action.flash":
    "Confirm before starting. Flashing will follow image preparation.",
  "action.running": "Task in progress. Keep your USB cable connected.",
  "result.ready": "Your image is ready",
  "result.flashed": "Flashing complete",
  "result.desc": "Stock and patched images are saved in your output directory.",
  "result.failed": "Task incomplete",
  "result.copy": "Copy path",
  "result.copied": "Copied",
  "result.copyFailed": "Copy failed. Select the path manually.",
  "btn.flash": "Prepare & auto-flash",
  "btn.busy": "Working…",
  "summary.safe": "Creates images without writing to your device.",
  "summary.caution": "Writes to your device. Unlock and back up first.",
  tagline: "Match the stock image. Patch with KernelSU. Stay in control.",
  "config.title": "Prepare an image",
  "field.device": "Target device",
  "field.ota": "Stock Recovery package",
  "device.title": "Device overview",
  "summary.title": "Run summary",
  "btn.start": "Prepare image",
  "summary.flash.no": "Prepare only",
  "summary.flash.yes": "Auto-flash",
  "summary.flash.yes.noreboot": "Stay in Fastboot",
  "status.idle": "Ready when you are",
  "log.idle": "Task details will appear here when you start.",
});

Object.assign(I18N.zh, {
  "nav.overview": "工作台",
  "config.title": "准备修补镜像",
  "native.taskKind": "KernelSU 镜像修补",
  "advanced.hint": "覆盖 KMI / 自备内核",
  "btn.start": "准备镜像",
  "btn.flash": "准备并刷写",
  "toggle.flash": "准备完成后自动刷写",
  "toggle.flash.desc": "进入 Fastboot，写入当前槽位的启动分区。",
  "toggle.noreboot.desc": "保留在 Fastboot，方便检查刷写结果。",
  "mode.recommended": "推荐",
  "action.ready": "仅生成镜像，不写入设备。",
  "action.connect": "请连接并选择 ADB 设备。",
});
Object.assign(I18N.en, {
  "nav.overview": "Workspace",
  "config.title": "Prepare patched image",
  "native.taskKind": "KernelSU image patching",
  "advanced.hint": "Override KMI / custom kernel",
  "btn.start": "Prepare image",
  "btn.flash": "Prepare & flash",
  "toggle.flash": "Flash automatically after preparation",
  "action.ready": "Creates an image without writing to your device.",
  "action.connect": "Connect and select an ADB device.",
});

Object.assign(I18N.zh, {
  "logs.subtitle": "查看本次任务的完整执行记录",
  "logs.export": "导出日志",
  "logs.task": "本次任务",
  "logs.noTask": "尚未开始任务",
  "logs.currentTask": "当前任务的执行记录",
  "logs.prepareTask": "准备修补镜像",
  "logs.flashTask": "准备并刷写镜像",
  "logs.search": "搜索日志内容",
  "logs.all": "全部",
  "logs.warn": "警告",
  "logs.error": "错误",
  "logs.stages": "执行阶段",
  "logs.allRecords": "全部记录",
  "logs.device": "设备检查",
  "logs.rom": "原厂镜像",
  "logs.patch": "KernelSU 修补",
  "logs.flash": "设备刷写",
  "logs.time": "时间",
  "logs.level": "级别",
  "logs.message": "消息",
  "logs.detailHint": "点击记录可查看事件详情",
  "logs.copy": "复制全部",
  "logs.noMatches": "没有匹配的日志",
  "logs.copied": "日志已复制",
  "logs.copyFailed": "复制失败",
  "logs.exported": "日志已导出",
  "logs.exportFailed": "导出失败",
  "logs.phase.pending": "未执行",
  "logs.phase.running": "执行中",
  "logs.phase.passed": "通过",
  "logs.phase.failed": "失败",
});
Object.assign(I18N.en, {
  "logs.subtitle": "Execution details for the current task",
  "logs.export": "Export log",
  "logs.task": "Current task",
  "logs.noTask": "No task has started",
  "logs.currentTask": "Execution records for this task",
  "logs.prepareTask": "Prepare patched image",
  "logs.flashTask": "Prepare and flash image",
  "logs.search": "Search log content",
  "logs.all": "All",
  "logs.warn": "Warnings",
  "logs.error": "Errors",
  "logs.stages": "Stages",
  "logs.allRecords": "All records",
  "logs.device": "Device checks",
  "logs.rom": "Stock image",
  "logs.patch": "KernelSU patch",
  "logs.flash": "Device flash",
  "logs.time": "Time",
  "logs.level": "Level",
  "logs.message": "Message",
  "logs.detailHint": "Select a record to inspect its event",
  "logs.copy": "Copy all",
  "logs.noMatches": "No matching records",
  "logs.copied": "Log copied",
  "logs.copyFailed": "Copy failed",
  "logs.exported": "Log exported",
  "logs.exportFailed": "Export failed",
  "logs.phase.pending": "Not run",
  "logs.phase.running": "Running",
  "logs.phase.passed": "Passed",
  "logs.phase.failed": "Failed",
});

Object.assign(I18N.zh, {
  "storage.generalTab": "通用设置",
  "storage.tab": "存储与清理",
  "storage.heading": "管理本机产物",
  "storage.scope": "仅扫描当前输出目录中可识别的 ROM 和镜像文件。",
  "storage.refresh": "重新扫描",
  "storage.total": "识别文件占用",
  "storage.selected": "预计清理占用",
  "storage.preserve": "默认保留原厂备份和修补镜像；未完成的下载也会保留。",
  "storage.none": "请选择需要清理的文件",
  "storage.clean": "清理所选",
  "storage.processing": "正在处理…",
  "storage.selectGroup": "选择本组文件",
  "storage.kind.rom": "ROM 完整包",
  "storage.kind.partial": "未完成下载",
  "storage.kind.stock": "原厂备份镜像",
  "storage.kind.patched": "修补镜像",
  "storage.romReady": "已验证对应原厂备份，可重新下载此 ROM。",
  "storage.romUnknown": "尚未确认对应备份，默认保留。",
  "storage.note.partial": "删除后不能继续断点下载。",
  "storage.note.stock": "用于恢复设备，镜像与校验记录一起清理。",
  "storage.note.patched": "删除后需要重新修补，校验记录会一起清理。",
  "storage.files": "查看文件",
  "storage.scanning": "正在扫描文件，并验证原厂备份…",
  "storage.empty": "没有可清理的文件。",
  "storage.cancelled": "已取消清理，文件保持不变。",
  "storage.failed": "未能删除的文件",
  "storage.desktopOnly": "请在桌面应用中使用存储清理。",
  "storage.warning":
    "所选包含原厂备份或修补镜像，删除后可能需要重新下载或修补。",
});
Object.assign(I18N.en, {
  "storage.generalTab": "General",
  "storage.tab": "Storage & cleanup",
  "storage.heading": "Manage local artifacts",
  "storage.scope":
    "Only recognized ROM and image files in the configured output directory are scanned.",
  "storage.refresh": "Scan again",
  "storage.total": "Recognized file usage",
  "storage.selected": "Selected allocation",
  "storage.preserve":
    "Stock backups, patched images and incomplete downloads are kept by default.",
  "storage.none": "Select files to clean up",
  "storage.clean": "Clean selected",
  "storage.processing": "Processing…",
  "storage.selectGroup": "Select files in this group",
  "storage.kind.rom": "Full Recovery ROM",
  "storage.kind.partial": "Incomplete download",
  "storage.kind.stock": "Stock backup image",
  "storage.kind.patched": "Patched image",
  "storage.romReady":
    "A matching stock backup is verified. This ROM can be downloaded again.",
  "storage.romUnknown":
    "A matching backup has not been confirmed. Kept by default.",
  "storage.note.partial": "Deleting this prevents resuming the download.",
  "storage.note.stock":
    "Used for device recovery. The image and verification record are cleaned together.",
  "storage.note.patched":
    "Requires patching again. Includes the verification record.",
  "storage.files": "View files",
  "storage.scanning": "Scanning files and verifying stock backups…",
  "storage.empty": "No cleanable files found.",
  "storage.cancelled": "Cleanup cancelled. Files were kept.",
  "storage.failed": "Files that could not be removed",
  "storage.desktopOnly": "Use storage cleanup in the desktop application.",
  "storage.warning":
    "The selection includes stock backups or patched images. You may need to download or patch again.",
});

const PREPARE_STEP_KEYS = [
  "step.usb",
  "step.info",
  "step.fetch",
  "step.patch",
  "step.ready",
];
const FLASH_STEP_KEYS = [
  "step.usb",
  "step.info",
  "step.fetch",
  "step.patch",
  "step.fastboot",
  "step.verify",
  "step.flash",
  "step.reboot",
];
const CHECK_MARK = "✓";

const $ = (id) => document.getElementById(id);
const mockDeviceEnabled =
  new URLSearchParams(window.location.search).get("mock") === "device";
const mockApp = {
  Devices: async () =>
    mockDeviceEnabled ? [{ Serial: "d098a3f", Mode: "adb" }] : [],
  Inspect: async () => ({
    Serial: "d098a3f",
    Mode: "adb",
    Model: "Xiaomi 15",
    Codename: "dada",
    Android: "17",
    Version: "OS4.0.4.0.XOCCNXM",
    Slot: "b",
    KMI: "android15-6.6",
    Kernel: "6.6.77-android15-8",
    HasInitBoot: true,
  }),
  GetSettings: async () => ({
    outputDir: "~/Downloads/KSUForge",
    language: "zh",
    mode: "lkm",
    partition: "auto",
    autoFlash: false,
    noReboot: false,
    catalogUrl: "",
    mirrors: "",
    connections: 16,
    ksud: "",
    logAutoScroll: true,
  }),
  SaveSettings: async () => {},
  State: async () => ({ status: "idle", language: "zh", log: "" }),
  Start: async () => {},
  Confirm: async () => false,
  ChooseOutputDirectory: async () => "",
  ChooseFile: async () => "",
  ChooseOTAFile: async () => "",
};
const go = () => (window.go && window.go.main && window.go.main.App) || mockApp;
let lang = "zh";
let autoScroll = true;
let currentPage = "workspace";
let taskContext = null;
let isRunning = false;
let lastState = null;
let settingsTrigger = null;
let detected = false;
let lastDevice = null;
let deviceList = [];
let deviceCount = 0;
let settings = {
  outputDir: "",
  language: "zh",
  mode: "lkm",
  partition: "auto",
  autoFlash: false,
  noReboot: false,
  catalogUrl: "",
  mirrors: "",
  connections: 16,
  ksud: "",
  logAutoScroll: true,
};
let saveTimer = null;
let flowFlash = false;
let flowNoReboot = false;
let selectedROMPath = "";

function t(key) {
  return (I18N[lang] && I18N[lang][key]) || key;
}
function esc(v) {
  return String(v).replace(
    /[&<>"']/g,
    (c) =>
      ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" })[
        c
      ],
  );
}
function baseName(p) {
  return String(p).split(/[\\/]/).pop();
}

function formatTransfer(downloadedMiB, totalMiB) {
  if (totalMiB >= 1024)
    return `${(downloadedMiB / 1024).toFixed(2)} / ${(totalMiB / 1024).toFixed(2)} GiB`;
  return `${downloadedMiB.toFixed(1)} / ${totalMiB.toFixed(1)} MiB`;
}

function formatSpeed(mibPerSecond) {
  if (mibPerSecond >= 1) return `${mibPerSecond.toFixed(1)} MiB/s`;
  if (mibPerSecond > 0) return `${Math.round(mibPerSecond * 1024)} KiB/s`;
  return "0 KiB/s";
}

function readLogValue(line, key) {
  const match = line.match(
    new RegExp(`(?:^|\\s)${key}=("(?:\\\\.|[^"])*"|[^\\s]+)`),
  );
  if (!match) return "";
  if (match[1].startsWith('"')) {
    try {
      return JSON.parse(match[1]);
    } catch (e) {
      return match[1].slice(1, -1);
    }
  }
  return match[1];
}

function formatLogForDisplay(log, language = lang) {
  return String(log || "")
    .split("\n")
    .map((line) => {
      if (line.includes("event=workflow.step.manual_ready")) {
        const path = readLogValue(line, "path");
        const partition = readLogValue(line, "partition");
        return language === "zh"
          ? `[5/5] Patch 完成 · 目标分区 ${partition}\nPatched 镜像：${path}`
          : `[5/5] Patch complete · Target ${partition}\nPatched image: ${path}`;
      }
      if (line.includes("event=workflow.manual_flash_command")) {
        const command = readLogValue(line, "command");
        return language === "zh"
          ? `手动刷写：${command}`
          : `Manual flash: ${command}`;
      }
      if (!line.includes("event=rom.download_progress")) return line;
      const readNumber = (key) => {
        const value = readLogValue(line, key);
        return value ? Number(value) : NaN;
      };
      const downloaded = readNumber("downloaded_mib");
      const total = readNumber("total_mib");
      const percent = readNumber("percent");
      const speed = readNumber("mib_per_second");
      if (![downloaded, total, percent, speed].every(Number.isFinite))
        return line;

      const safePercent = Math.max(0, Math.min(100, percent));
      const width = 20;
      const filled = Math.round((safePercent / 100) * width);
      const bar = `${"█".repeat(filled)}${"░".repeat(width - filled)}`;
      const label = language === "zh" ? "下载 ROM" : "Downloading ROM";
      return `${label} [${bar}] ${safePercent.toFixed(1)}% · ${formatTransfer(downloaded, total)} · ${formatSpeed(speed)}`;
    })
    .join("\n");
}

function applyLang(l) {
  lang = l;
  document.documentElement.lang = l === "zh" ? "zh-CN" : "en";
  document.querySelectorAll("[data-i18n]").forEach((el) => {
    const k = el.getAttribute("data-i18n");
    if (I18N[l][k] != null) el.textContent = I18N[l][k];
  });
  document.querySelectorAll("[data-i18n-title]").forEach((el) => {
    const k = el.getAttribute("data-i18n-title");
    if (I18N[l][k] != null) {
      el.title = I18N[l][k];
      if (el.hasAttribute("aria-label"))
        el.setAttribute("aria-label", I18N[l][k]);
    }
  });
  document.querySelectorAll("[data-i18n-placeholder]").forEach((el) => {
    const k = el.getAttribute("data-i18n-placeholder");
    if (I18N[l][k] != null) el.placeholder = I18N[l][k];
  });
  const ak = autoScroll ? "log.autoScroll.on" : "log.autoScroll.off";
  $("autoScrollText").setAttribute("data-i18n", ak);
  $("autoScrollText").textContent = t(ak);
  LogView.render();
  renderDeviceOptions();
  renderDevicePill();
  if (detected && lastDevice) fillDevice(lastDevice);
  updateSummary();
  renderROMSelection();
  if (lastState) renderResult(lastState);
  StorageView.render();
}

function currentStepKeys() {
  const keys = flowFlash ? [...FLASH_STEP_KEYS] : [...PREPARE_STEP_KEYS];
  if (flowFlash && flowNoReboot) keys[keys.length - 1] = "step.done";
  return keys;
}

function renderStepper(n = 1) {
  const keys = currentStepKeys();
  const total = keys.length;
  const stepper = $("stepper");
  stepper.className = `stepper ${flowFlash ? "flash-flow" : "prepare-flow"}`;
  stepper.innerHTML =
    `<div class="track" style="left:${50 / total}%;right:${50 / total}%"><i id="trackFill"></i></div>` +
    keys
      .map(
        (key, index) =>
          `<div class="step" data-step="${index + 1}"><div class="node">${index + 1}</div><div class="name" data-i18n="${key}">${esc(t(key))}</div></div>`,
      )
      .join("");
  $("stepTotal").textContent = total;
  setStep(Math.min(n, total + 1));
}

function setStep(n) {
  const keys = currentStepKeys();
  const total = keys.length;
  const done = n > total;
  document.querySelectorAll("#stepper .step").forEach((el) => {
    const k = Number(el.dataset.step);
    el.classList.toggle("done", done || k < n);
    el.classList.toggle("active", !done && k === n);
    el.classList.toggle(
      "failed",
      lastState?.status === "failed" && !done && k === n,
    );
    el.querySelector(".now")?.remove();
    const node = el.querySelector(".node");
    if (done || k < n) node.textContent = CHECK_MARK;
    else node.textContent = String(k);
    if (!done && k === n) {
      const now = document.createElement("span");
      now.className = "now";
      now.setAttribute("data-i18n", "progress.current");
      now.textContent = t("progress.current");
      el.appendChild(now);
    }
  });
  const shown = Math.min(n, total);
  $("stepNow").textContent = shown;
  const key = keys[shown - 1];
  $("stepName").setAttribute("data-i18n", key);
  $("stepName").textContent = t(key);
  $("trackFill").style.width =
    (done ? 100 : ((n - 1) / Math.max(1, total - 1)) * 100) + "%";
}

function renderDevicePill() {
  const pill = $("devicePill"),
    txt = $("deviceText");
  if (deviceCount > 0) {
    pill.className = "pill on";
    txt.removeAttribute("data-i18n");
    txt.textContent =
      lang === "zh"
        ? `${deviceCount} 台设备已连接`
        : `${deviceCount} device${deviceCount > 1 ? "s" : ""} connected`;
  } else {
    pill.className = "pill off";
    txt.setAttribute("data-i18n", "device.disconnected");
    txt.textContent = t("device.disconnected");
  }
}

function renderDeviceOptions() {
  const sel = $("deviceSelect"),
    prev = sel.value;
  sel.innerHTML = deviceList.length
    ? deviceList
        .map(
          (d) =>
            `<option value="${esc(d.Serial)}">${esc(d.Serial)} · ${esc(d.Mode)}</option>`,
        )
        .join("")
    : `<option value="">${esc(t("device.none"))}</option>`;
  if (prev) sel.value = prev;
}

function fillDevice(d) {
  $("deviceEmpty").style.display = "none";
  $("deviceData").style.display = "";
  $("deviceTag").removeAttribute("data-i18n");
  $("deviceTag").textContent = (d.Mode || "adb").toUpperCase();
  $("devModel").textContent = d.Model || "—";
  $("devCodename").textContent = d.Codename || "—";
  $("devAndroid").textContent = d.Android || "—";
  $("devVersion").textContent = d.Version || "—";
  $("devVersion").title = d.Version || "";
  $("devSlot").textContent = d.Slot || "—";
  $("devKMI").textContent = d.KMI || "—";
  $("devKernel").textContent = d.Kernel || "—";
  $("devKernel").title = d.Kernel || "";
  const ib = $("devInitBoot");
  ib.textContent = d.HasInitBoot ? t("kv.present") : t("kv.absent");
  ib.className = "chip" + (d.HasInitBoot ? "" : " warn");
}

function showEmpty() {
  detected = false;
  lastDevice = null;
  $("deviceEmpty").style.display = "";
  $("deviceData").style.display = "none";
  $("deviceTag").setAttribute("data-i18n", "device.none");
  $("deviceTag").textContent = t("device.none");
  $("sumRom").removeAttribute("data-i18n");
  $("sumRom").textContent = "—";
  updateSummary();
}

async function refreshDevices() {
  try {
    deviceList = (await go().Devices()) || [];
  } catch (e) {
    deviceList = [];
  }
  deviceCount =
    deviceList.filter((d) => d.Mode === "adb").length || deviceList.length;
  renderDeviceOptions();
  renderDevicePill();
  if (deviceList.length) {
    const sel = $("deviceSelect");
    const adb = deviceList.find((d) => d.Mode === "adb");
    const pick =
      deviceList.find((d) => d.Serial === sel.value) || adb || deviceList[0];
    sel.value = pick.Serial;
    await inspect(pick.Serial);
  } else {
    showEmpty();
    setStep(1);
  }
}

async function inspect(serial) {
  if (!serial) {
    showEmpty();
    setStep(1);
    return;
  }
  try {
    const d = await go().Inspect(serial);
    lastDevice = d;
    detected = true;
    fillDevice(d);
    updateSummary();
    setStep(2);
  } catch (e) {
    showEmpty();
    setStep(2);
  }
}

function currentMode() {
  return settings.mode || "lkm";
}

function updateSummary() {
  const mode = currentMode();
  $("sumMode").textContent = mode === "gki" ? "GKI" : "LKM";
  const part = settings.partition || "auto";
  if (part !== "auto") {
    $("sumPartition").textContent = part;
  } else if (lastDevice && lastDevice.HasInitBoot) {
    $("sumPartition").textContent =
      mode === "gki" || String(lastDevice.KMI).startsWith("android12-")
        ? "boot"
        : "init_boot";
  } else if (lastDevice) {
    $("sumPartition").textContent = "boot";
  } else {
    $("sumPartition").textContent = "—";
  }
  const flash = !!settings.autoFlash;
  const el = $("sumFlash");
  const key = flash
    ? settings.noReboot
      ? "summary.flash.yes.noreboot"
      : "summary.flash.yes"
    : "summary.flash.no";
  el.setAttribute("data-i18n", key);
  el.textContent = t(key);
  const connections = Math.max(
    1,
    Math.min(32, Number(settings.connections) || 16),
  );
  const custom = String(settings.mirrors || "").trim() !== "";
  $("sumDownload").textContent =
    lang === "zh"
      ? `${connections} 连接 · ${custom ? "自定义优先" : "自动镜像"}`
      : `${connections} connections · ${custom ? "custom first" : "auto mirrors"}`;
  if (!selectedROMPath)
    $("sumRom").textContent = lastDevice ? lastDevice.Version : "—";
  updateActions();
}

function updateActions() {
  const ready = detected && lastDevice && lastDevice.Mode === "adb";
  $("btnStart").disabled = isRunning || !ready;
  $("btnPrepare").disabled = isRunning || !ready;
  $("btnPrepare").classList.toggle("hidden", !settings.autoFlash);
  $("setNoReboot").disabled = isRunning || !settings.autoFlash;
  $("noRebootOption").classList.toggle("hidden", !settings.autoFlash);
  const key = isRunning
    ? "btn.busy"
    : settings.autoFlash
      ? "btn.flash"
      : "btn.start";
  $("startLabel").textContent = t(key);
  $("startLabel").setAttribute("data-i18n", key);
  const hint = isRunning
    ? "action.running"
    : !ready
      ? "action.connect"
      : settings.autoFlash
        ? "action.flash"
        : "action.ready";
  $("actionHint").textContent = t(hint);
  $("actionHint").setAttribute("data-i18n", hint);
  const safety = document.querySelector(".safety-note");
  safety.classList.toggle("warning", !!settings.autoFlash);
  const note = settings.autoFlash ? "summary.caution" : "summary.safe";
  safety.querySelector("span").textContent = t(note);
  safety.querySelector("span").setAttribute("data-i18n", note);
}

function renderROMSelection() {
  $("otaPath").value = selectedROMPath;
  $("clearOTA").classList.toggle("hidden", !selectedROMPath);
  $("verifyROMOption").classList.toggle("hidden", !selectedROMPath);
  if (selectedROMPath) {
    $("sumRom").removeAttribute("data-i18n");
    $("sumRom").textContent = baseName(selectedROMPath);
  } else if (!lastDevice) {
    $("sumRom").textContent = "—";
  }
}

function setRunning(running) {
  isRunning = running;
  document.body.dataset.running = String(running);
  [
    "btnStart",
    "btnPrepare",
    "refresh",
    "deviceSelect",
    "chooseKernel",
    "chooseOTA",
    "clearOTA",
    "verifyROM",
    "setAutoFlash",
    "setNoReboot",
    "setMirrors",
    "setConnections",
    "openSettings",
    "modeLKM",
    "modeGKI",
    "kmi",
    "kernel",
  ].forEach((id) => {
    const el = $(id);
    if (el) el.disabled = running;
  });
  $("runBar").style.display = running ? "" : "none";
  updateActions();
  StorageView.sync();
}

/* ---- settings ---- */
function setMode(mode) {
  $("kernelField").style.display = mode === "gki" ? "" : "none";
  $("modeLKM").setAttribute("aria-pressed", String(mode !== "gki"));
  $("modeGKI").setAttribute("aria-pressed", String(mode === "gki"));
  const key = mode === "gki" ? "mode.gki.desc" : "mode.lkm.desc";
  $("modeDescription").textContent = t(key);
  $("modeDescription").setAttribute("data-i18n", key);
  if (mode === "gki") document.querySelector("details.adv").open = true;
}

function renderAutoScroll() {
  $("autoScroll").classList.toggle("on", autoScroll);
  $("autoScroll").setAttribute("aria-pressed", String(autoScroll));
  const k = autoScroll ? "log.autoScroll.on" : "log.autoScroll.off";
  $("autoScrollText").setAttribute("data-i18n", k);
  $("autoScrollText").textContent = t(k);
}

function setPage(page) {
  currentPage = page === "logs" ? "logs" : "workspace";
  $("workspace").classList.toggle("logs-active", currentPage === "logs");
  $("taskLog").classList.toggle("hidden", currentPage !== "logs");
  document.querySelectorAll("nav [data-page]").forEach((link) => {
    const active = link.dataset.page === currentPage;
    link.classList.toggle("selected", active);
    if (active) link.setAttribute("aria-current", "page");
    else link.removeAttribute("aria-current");
  });
  if (currentPage === "logs") {
    LogView.render();
    $("taskLog").focus({ preventScroll: true });
    if (autoScroll) requestAnimationFrame(scrollLog);
  }
}

function applySettings() {
  setMode(settings.mode || "lkm");
  $("setAutoFlash").checked = !!settings.autoFlash;
  $("setNoReboot").checked = !!settings.noReboot;
  autoScroll = settings.logAutoScroll !== false;
  renderAutoScroll();
  if (settings.language && settings.language !== lang)
    applyLang(settings.language);
  flowFlash = !!settings.autoFlash;
  flowNoReboot = !!settings.noReboot;
  renderStepper(detected ? 2 : 1);
  updateSummary();
}

function syncFromUI() {
  settings.language = lang;
}

function persistSettings() {
  syncFromUI();
  clearTimeout(saveTimer);
  saveTimer = setTimeout(() => {
    go()
      .SaveSettings(settings)
      .catch(() => {});
  }, 300);
}

async function loadSettings() {
  try {
    const s = await go().GetSettings();
    if (s) settings = Object.assign(settings, s);
  } catch (e) {}
  applySettings();
}

function openSettings() {
  StorageView.setTab("general");
  $("setOutputDir").value = settings.outputDir || "";
  $("setLanguage").value = settings.language || "zh";
  $("setMode").value = settings.mode || "lkm";
  $("setPartition").value = settings.partition || "auto";
  $("setAutoFlash").checked = !!settings.autoFlash;
  $("setNoReboot").checked = !!settings.noReboot;
  $("setCatalog").value = settings.catalogUrl || "";
  $("setMirrors").value = settings.mirrors || "";
  $("setConnections").value = settings.connections || 16;
  $("setKsud").value = settings.ksud || "";
  $("setLogScroll").checked = settings.logAutoScroll !== false;
  $("settingsError").textContent = "";
  settingsTrigger = document.activeElement;
  $("settingsModal").classList.remove("hidden");
  document.querySelector(".workspace").inert = true;
  document.querySelector(".sidebar").inert = true;
  $("setLanguage").focus();
}

function closeSettings() {
  $("settingsModal").classList.add("hidden");
  document.querySelector(".workspace").inert = false;
  document.querySelector(".sidebar").inert = false;
  settingsTrigger?.focus();
}

async function saveSettingsForm() {
  const next = {
    outputDir: $("setOutputDir").value.trim(),
    language: $("setLanguage").value,
    mode: $("setMode").value,
    partition: $("setPartition").value,
    autoFlash: $("setAutoFlash").checked,
    noReboot: $("setNoReboot").checked,
    catalogUrl: $("setCatalog").value.trim(),
    mirrors: $("setMirrors").value,
    connections: Math.max(
      1,
      Math.min(32, Number($("setConnections").value) || 16),
    ),
    ksud: $("setKsud").value.trim(),
    logAutoScroll: $("setLogScroll").checked,
  };
  try {
    await go().SaveSettings(next);
  } catch (e) {
    $("settingsError").textContent = String(e);
    return;
  }
  settings = next;
  applySettings();
  closeSettings();
}

function scrollLog() {
  const el = $("log");
  el.scrollTop = el.scrollHeight;
}

function updateFlowFromLog(log) {
  let backendStage = 1,
    backendTotal = flowFlash ? 8 : 5,
    m;
  const re = /(?:^|\s)step=(\d+)(?:\s|.*?\s)total=(\d+)/g;
  while ((m = re.exec(log)) !== null) {
    backendStage = Number(m[1]);
    backendTotal = Number(m[2]);
  }
  let stage;
  if (flowFlash) {
    stage =
      backendStage <= 1
        ? 1
        : backendStage <= 3
          ? 2
          : backendStage === 4
            ? 3
            : backendStage === 5
              ? 4
              : backendStage === 6
                ? 5
                : backendStage === 7
                  ? 6
                  : 8;
    if (log.includes("event=flash.completed")) stage = Math.max(stage, 7);
    if (log.includes("event=flash.rebooting")) stage = 8;
  } else {
    stage =
      backendTotal === 5
        ? backendStage
        : backendStage <= 1
          ? 1
          : backendStage <= 3
            ? 2
            : backendStage === 4
              ? 3
              : backendStage === 5
                ? 4
                : 5;
  }
  setStep(stage);
}

function renderState(s) {
  lastState = s;
  const st = $("logStatus");
  st.className =
    "status" + (s.status && s.status !== "idle" ? " " + s.status : "");
  const key = "status." + (s.status || "idle");
  st.setAttribute("data-i18n", I18N[lang][key] ? key : "status.idle");
  st.textContent = t(I18N[lang][key] ? key : "status.idle");

  LogView.render(s);
  if (s.log) updateFlowFromLog(s.log);

  const running = s.status === "running";
  setRunning(running);
  renderResult(s);
  document
    .querySelectorAll("#stepper .step")
    .forEach((el) =>
      el.classList.toggle(
        "failed",
        s.status === "failed" && el.classList.contains("active"),
      ),
    );

  if (s.status === "succeeded" && s.result) {
    if (s.result.Partition) $("sumPartition").textContent = s.result.Partition;
    if (s.result.ROMPath) {
      $("sumRom").removeAttribute("data-i18n");
      $("sumRom").textContent = baseName(s.result.ROMPath);
    }
    setStep(9);
  }
}

function renderResult(s) {
  const panel = $("resultPanel");
  const done = s.status === "succeeded",
    failed = s.status === "failed";
  panel.classList.toggle("hidden", !done && !failed);
  panel.classList.toggle("failed", failed);
  if (!done && !failed) return;
  $("resultTitle").textContent = t(
    failed ? "result.failed" : flowFlash ? "result.flashed" : "result.ready",
  );
  $("resultDescription").textContent = failed
    ? String(s.error || "")
    : t("result.desc");
  const path = done ? s.result?.PatchedImage || s.result?.StockImage || "" : "";
  $("resultPath").textContent = path;
  $("copyResult").classList.toggle("hidden", !path);
  $("copyResult").textContent = t("result.copy");
}

async function poll() {
  try {
    const s = await go().State();
    renderState(s);
    if (s.status === "running") setTimeout(poll, 600);
  } catch (e) {
    setTimeout(poll, 1200);
  }
}

async function start(flash) {
  if (isRunning || !detected || lastDevice?.Mode !== "adb") return;
  if (flash) {
    let ok = false;
    try {
      ok = await go().Confirm(
        t("confirm.title"),
        `${t("confirm.flash")}

${t("field.device")}: ${lastDevice.Model} (${lastDevice.Codename})
${t("kv.version")}: ${lastDevice.Version}
${t("summary.partition")}: ${$("sumPartition").textContent} · Slot ${lastDevice.Slot || "—"}`,
        t("confirm.ok"),
        t("confirm.cancel"),
      );
    } catch (e) {
      ok = false;
    }
    if (!ok) return;
  }
  const req = {
    Serial: $("deviceSelect").value,
    Mode: settings.mode || "lkm",
    Partition: settings.partition || "auto",
    OutputDir: settings.outputDir || "",
    Catalog: settings.catalogUrl || "",
    KSUd: settings.ksud || "",
    KMI: $("kmi").value,
    Kernel: $("kernel").value,
    Mirrors: settings.mirrors || "",
    Connections: Math.max(1, Math.min(32, Number(settings.connections) || 16)),
    ROMPath: selectedROMPath,
    VerifyROM: selectedROMPath ? $("verifyROM").checked : true,
    Flash: flash,
    NoReboot: !!settings.noReboot,
  };
  flowFlash = flash;
  flowNoReboot = !!settings.noReboot;
  renderStepper(detected ? 2 : 1);
  setRunning(true);
  $("resultPanel").classList.add("hidden");
  taskContext = {
    ...lastDevice,
    partition: $("sumPartition").textContent,
    flash,
  };
  lastState = { status: "running", language: lang, log: "" };
  LogView.reset();
  clearTimeout(saveTimer);
  syncFromUI();
  go()
    .SaveSettings(settings)
    .catch(() => {});
  try {
    await go().Start(req);
    poll();
  } catch (e) {
    setRunning(false);
    renderState({ status: "failed", language: lang, error: String(e) });
  }
}

/* ---- events ---- */
["lkm", "gki"].forEach((mode) => {
  $(mode === "lkm" ? "modeLKM" : "modeGKI").onclick = () => {
    if (isRunning) return;
    settings.mode = mode;
    setMode(mode);
    updateSummary();
    persistSettings();
  };
});

document.querySelectorAll("nav [data-page]").forEach((link) => {
  link.onclick = (event) => {
    event.preventDefault();
    setPage(link.dataset.page);
  };
});
$("deviceSelect").onchange = () => inspect($("deviceSelect").value);
$("refresh").onclick = refreshDevices;
$("btnDetect").onclick = refreshDevices;
$("chooseKernel").onclick = async () => {
  const v = await go().ChooseFile(t("field.kernel"));
  if (v) $("kernel").value = v;
};
$("chooseOTA").onclick = async () => {
  const v = await go().ChooseOTAFile();
  if (v) {
    selectedROMPath = v;
    $("verifyROM").checked = true;
    renderROMSelection();
  }
};
$("clearOTA").onclick = () => {
  selectedROMPath = "";
  renderROMSelection();
  updateSummary();
};
$("btnPrepare").onclick = () => start(false);
$("btnStart").onclick = () => start(!!settings.autoFlash);
$("setAutoFlash").onchange = () => {
  settings.autoFlash = $("setAutoFlash").checked;
  flowFlash = settings.autoFlash;
  renderStepper(detected ? 2 : 1);
  updateSummary();
  persistSettings();
};
$("setNoReboot").onchange = () => {
  settings.noReboot = $("setNoReboot").checked;
  flowNoReboot = settings.noReboot;
  renderStepper(detected ? 2 : 1);
  updateSummary();
  persistSettings();
};
$("autoScroll").onclick = () => {
  autoScroll = !autoScroll;
  settings.logAutoScroll = autoScroll;
  renderAutoScroll();
  persistSettings();
};
$("docsLink").onclick = (e) => {
  if (window.runtime && window.runtime.BrowserOpenURL) {
    e.preventDefault();
    window.runtime.BrowserOpenURL(
      "https://kernelsu.org/guide/installation.html",
    );
  }
};
$("copyResult").onclick = async () => {
  try {
    const path = $("resultPath").textContent;
    if (window.runtime?.ClipboardSetText)
      await window.runtime.ClipboardSetText(path);
    else await navigator.clipboard.writeText(path);
    $("copyResult").textContent = t("result.copied");
  } catch (e) {
    $("copyResult").textContent = t("result.copyFailed");
  }
};

$("openSettings").onclick = openSettings;
$("closeSettings").onclick = closeSettings;
$("cancelSettings").onclick = closeSettings;
$("saveSettings").onclick = saveSettingsForm;
$("setChooseOutput").onclick = async () => {
  const v = await go().ChooseOutputDirectory();
  if (v) $("setOutputDir").value = v;
};
$("setChooseKsud").onclick = async () => {
  const v = await go().ChooseFile(t("settings.chooseKsud"));
  if (v) $("setKsud").value = v;
};
$("settingsModal").onclick = (e) => {
  if (e.target === $("settingsModal")) closeSettings();
};
document.addEventListener("keydown", (e) => {
  const preferencesOpen = !$("settingsModal").classList.contains("hidden");
  if ((e.metaKey || e.ctrlKey) && e.key === "," && !isRunning) {
    e.preventDefault();
    if (!preferencesOpen) openSettings();
    return;
  }
  if (
    (e.metaKey || e.ctrlKey) &&
    e.key === "Enter" &&
    !preferencesOpen &&
    currentPage === "workspace"
  ) {
    e.preventDefault();
    start(!!settings.autoFlash);
    return;
  }

  if (
    (e.metaKey || e.ctrlKey) &&
    e.key.toLowerCase() === "f" &&
    currentPage === "logs" &&
    !preferencesOpen
  ) {
    e.preventDefault();
    $("logSearch").focus();
    return;
  }
  if (e.key === "Tab" && !$("settingsModal").classList.contains("hidden")) {
    const items = [
      ...$("settingsModal").querySelectorAll(
        "button:not(:disabled), input:not(:disabled), select:not(:disabled), textarea:not(:disabled)",
      ),
    ].filter((el) => el.getClientRects().length);
    const first = items[0],
      last = items[items.length - 1];
    if (e.shiftKey && document.activeElement === first) {
      e.preventDefault();
      last.focus();
    } else if (!e.shiftKey && document.activeElement === last) {
      e.preventDefault();
      first.focus();
    }
  }
  if (e.key === "Escape" && !$("settingsModal").classList.contains("hidden"))
    closeSettings();
});

/* ---- init ---- */
StorageView.init();
LogView.init();
setPage("workspace");
applyLang("zh");
renderStepper(1);
renderROMSelection();
loadSettings();
refreshDevices();
poll();
