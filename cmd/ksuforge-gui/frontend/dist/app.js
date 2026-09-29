const I18N = {
  zh: {
    "tagline":"连接设备后，自动完成匹配、下载、提取、patch 与安全刷写。",
    "device.connected":"1 台设备已连接",
    "device.disconnected":"未检测到设备","device.none":"未检测",
    "device.none.title":"尚未检测到设备","device.none.hint":"连接手机并开启 USB 调试",
    "btn.detect":"检测设备",
    "log.idle":"等待检测设备…",
    "progress.step":"步骤","progress.current":"当前",
    "step.usb":"USB / ADB","step.detect":"检测设备","step.info":"设备信息","step.fetch":"获取原厂镜像",
    "step.patch":"KernelSU Patch","step.ready":"准备完成","step.fastboot":"Fastboot","step.verify":"安全检查","step.flash":"刷入","step.reboot":"重启","step.done":"完成",
    "config.title":"安装配置",
    "field.device":"ADB 设备","btn.refresh":"刷新",
    "field.ota":"Recovery OTA 包","hint.ota":"留空则自动匹配并下载","hint.ota.auto":"自动匹配当前系统版本","btn.clear":"清除",
    "ota.verify":"校验本地 OTA 的 MD5","ota.verify.desc":"根据当前设备版本查询官方校验值，推荐保持开启。",
    "field.mode":"KernelSU 模式","mode.recommended":"推荐","mode.lkm":"LKM","mode.gki":"GKI 替换内核",
    "mode.lkm.desc":"修补 init_boot / boot，保留原厂内核。",
    "mode.gki.desc":"替换 boot 中的内核，需自备匹配内核文件。",
    "field.partition":"目标分区","partition.auto":"自动（依据 KMI 与分区）",
    "field.output":"输出目录","hint.output":"默认 ~/Downloads/KSUForge","btn.choose":"选择",
    "download.title":"ROM 下载源","download.desc":"留空则自动测速小米官方 CDN；自定义源优先参与测速。","download.resume":"支持断点续传",
    "field.mirrors":"镜像源基础地址","hint.mirrors":"每行一个 HTTPS 地址","field.connections":"并发连接数",
    "download.strategy":"并行分片 · 镜像故障切换 · MD5 校验",
    "advanced":"高级选项","field.ksud":"ksud 路径","hint.ksud":"留空使用内置版本",
    "field.kmi":"KMI 覆盖","hint.kmi":"通常留空","field.kernel":"GKI 内核文件",
    "settings.title":"设置","settings.short":"设置","settings.general":"通用","settings.install":"安装默认值","settings.advanced":"下载与工具",
    "settings.language":"界面语言","settings.logScroll":"日志自动滚动","settings.logScroll.desc":"运行时自动滚动到最新日志。",
    "settings.autoFlash":"默认自动刷写","settings.noReboot":"默认刷写后不重启",
    "settings.catalog":"ROM 索引地址","settings.catalog.desc":"必须为 HTTPS，且包含一个 %s 作为设备代号占位符。",
    "settings.ksud":"ksud 路径","settings.chooseKsud":"选择 ksud",
    "btn.save":"保存","btn.cancel":"取消",
    "toggle.flash":"下载并 patch 后自动刷写","toggle.flash.desc":"会重启到 fastboot 并覆盖当前启动分区。",
    "toggle.noreboot":"刷写后不重启","toggle.noreboot.desc":"停留在 fastboot，方便手动检查。",
    "btn.prepare":"仅准备镜像","btn.start":"开始执行",
    "device.title":"设备信息","kv.model":"机型","kv.codename":"代号","kv.android":"Android",
    "kv.version":"系统版本","kv.slot":"当前 Slot","kv.kmi":"KMI","kv.kernel":"内核","kv.initboot":"init_boot",
    "kv.present":"存在","kv.absent":"不存在",
    "summary.title":"本次操作","summary.mode":"模式","summary.partition":"目标分区","summary.rom":"匹配 ROM",
    "summary.flash":"刷写","summary.flash.no":"否（仅准备）","summary.flash.yes":"是（自动刷写）","summary.flash.yes.noreboot":"是（刷写后停留在 fastboot）","summary.download":"下载策略",
    "log.title":"运行日志","log.autoScroll.on":"自动滚动","log.autoScroll.off":"已暂停","log.collapse":"折叠","log.expand":"展开",
    "status.idle":"空闲","status.running":"运行中","status.succeeded":"完成","status.failed":"失败",
    "foot.note":"真实刷写存在风险，请先解锁并备份。",
    "confirm.flash":"即将自动重启到 fastboot 并写入启动分区，确认继续？",
    "confirm.title":"确认操作","confirm.ok":"确认","confirm.cancel":"取消",
    "error.prefix":"错误: "
  },
  en: {
    "tagline":"Connect a device to auto-match, download, extract, patch and safely flash.",
    "device.connected":"1 device connected",
    "device.disconnected":"No device","device.none":"N/A",
    "device.none.title":"No device detected","device.none.hint":"Connect your phone and enable USB debugging",
    "btn.detect":"Detect Device",
    "log.idle":"Waiting for device…",
    "progress.step":"STEP","progress.current":"CURRENT",
    "step.usb":"USB / ADB","step.detect":"Detect","step.info":"Device Info","step.fetch":"Get Stock Image",
    "step.patch":"KernelSU Patch","step.ready":"Ready","step.fastboot":"Fastboot","step.verify":"Safety Check","step.flash":"Flash","step.reboot":"Reboot","step.done":"Done",
    "config.title":"Installation",
    "field.device":"ADB Device","btn.refresh":"Refresh",
    "field.ota":"Recovery OTA","hint.ota":"Leave empty to match and download","hint.ota.auto":"Match the current system version automatically","btn.clear":"Clear",
    "ota.verify":"Verify local OTA MD5","ota.verify.desc":"Fetch the official checksum for the connected device version. Recommended.",
    "field.mode":"KernelSU Mode","mode.recommended":"REC","mode.lkm":"LKM","mode.gki":"GKI Kernel",
    "mode.lkm.desc":"Patch init_boot / boot, keep the stock kernel.",
    "mode.gki.desc":"Replace the kernel in boot; a matching kernel file is required.",
    "field.partition":"Target Partition","partition.auto":"Auto (by KMI & partition)",
    "field.output":"Output Directory","hint.output":"Default ~/Downloads/KSUForge","btn.choose":"Choose",
    "download.title":"ROM Sources","download.desc":"Leave empty to benchmark Xiaomi CDNs; custom sources are tried first.","download.resume":"Resume enabled",
    "field.mirrors":"Mirror base URLs","hint.mirrors":"One HTTPS URL per line","field.connections":"Connections",
    "download.strategy":"Parallel ranges · failover · MD5 verification",
    "advanced":"Advanced","field.ksud":"ksud Path","hint.ksud":"Leave empty to use bundled",
    "field.kmi":"KMI Override","hint.kmi":"Usually leave empty","field.kernel":"GKI Kernel File",
    "settings.title":"Settings","settings.short":"Settings","settings.general":"General","settings.install":"Install Defaults","settings.advanced":"Download & Tools",
    "settings.language":"Language","settings.logScroll":"Log auto-scroll","settings.logScroll.desc":"Scroll to the newest log line while running.",
    "settings.autoFlash":"Default auto-flash","settings.noReboot":"Default no-reboot",
    "settings.catalog":"ROM Index URL","settings.catalog.desc":"Must be HTTPS and contain a single %s placeholder for the device codename.",
    "settings.ksud":"ksud Path","settings.chooseKsud":"Choose ksud",
    "btn.save":"Save","btn.cancel":"Cancel",
    "toggle.flash":"Auto-flash after patching","toggle.flash.desc":"Reboots to fastboot and overwrites the active boot partition.",
    "toggle.noreboot":"Don't reboot after flashing","toggle.noreboot.desc":"Stay in fastboot for manual inspection.",
    "btn.prepare":"Prepare Only","btn.start":"Start",
    "device.title":"Device","kv.model":"Model","kv.codename":"Codename","kv.android":"Android",
    "kv.version":"Version","kv.slot":"Active Slot","kv.kmi":"KMI","kv.kernel":"Kernel","kv.initboot":"init_boot",
    "kv.present":"Present","kv.absent":"Absent",
    "summary.title":"This Run","summary.mode":"Mode","summary.partition":"Partition","summary.rom":"Matched ROM",
    "summary.flash":"Flash","summary.flash.no":"No (prepare only)","summary.flash.yes":"Yes (auto-flash)","summary.flash.yes.noreboot":"Yes (stay in fastboot)","summary.download":"Download",
    "log.title":"Log","log.autoScroll.on":"Auto-scroll","log.autoScroll.off":"Paused","log.collapse":"Collapse","log.expand":"Expand",
    "status.idle":"Idle","status.running":"Running","status.succeeded":"Done","status.failed":"Failed",
    "foot.note":"Flashing is risky. Unlock the bootloader and back up first.",
    "confirm.flash":"This will reboot to fastboot and overwrite the active boot partition. Continue?",
    "confirm.title":"Confirm","confirm.ok":"Confirm","confirm.cancel":"Cancel",
    "error.prefix":"Error: "
  }
};

const PREPARE_STEP_KEYS = ["step.usb","step.info","step.fetch","step.patch","step.ready"];
const FLASH_STEP_KEYS = ["step.usb","step.info","step.fetch","step.patch","step.fastboot","step.verify","step.flash","step.reboot"];
const CHECK_MARK = '✓';

const $ = id => document.getElementById(id);
const mockDeviceEnabled = new URLSearchParams(window.location.search).get("mock") === "device";
const mockApp = {
  Devices: async () => mockDeviceEnabled ? [{Serial:"d098a3f", Mode:"adb"}] : [],
  Inspect: async () => ({Serial:"d098a3f",Mode:"adb",Model:"Xiaomi 15",Codename:"dada",Android:"16",Version:"OS3.0.305.0.WOCCNXM",Slot:"b",KMI:"android15-6.6",Kernel:"6.6.77-android15-8",HasInitBoot:true}),
  GetSettings: async () => ({outputDir:"~/Downloads/KSUForge",language:"zh",mode:"lkm",partition:"auto",autoFlash:false,noReboot:false,catalogUrl:"",mirrors:"",connections:16,ksud:"",logAutoScroll:true}),
  SaveSettings: async () => {}, State: async () => ({status:"idle",language:"zh",log:""}), Start: async () => {}, Confirm: async () => false,
  ChooseOutputDirectory: async () => "", ChooseFile: async () => "", ChooseOTAFile: async () => ""
};
const go = () => (window.go && window.go.main && window.go.main.App) || mockApp;
let lang = "zh";
let autoScroll = true;
let logExpanded = true;
let detected = false;
let lastDevice = null;
let deviceList = [];
let deviceCount = 0;
let settings = {
  outputDir: "", language: "zh", mode: "lkm", partition: "auto",
  autoFlash: false, noReboot: false, catalogUrl: "", mirrors: "", connections: 16,
  ksud: "", logAutoScroll: true
};
let saveTimer = null;
let flowFlash = false;
let flowNoReboot = false;
let selectedROMPath = "";

function t(key){ return (I18N[lang] && I18N[lang][key]) || key; }
function esc(v){ return String(v).replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c])); }
function baseName(p){ return String(p).split(/[\\/]/).pop(); }

function formatTransfer(downloadedMiB, totalMiB){
  if (totalMiB >= 1024) return `${(downloadedMiB / 1024).toFixed(2)} / ${(totalMiB / 1024).toFixed(2)} GiB`;
  return `${downloadedMiB.toFixed(1)} / ${totalMiB.toFixed(1)} MiB`;
}

function formatSpeed(mibPerSecond){
  if (mibPerSecond >= 1) return `${mibPerSecond.toFixed(1)} MiB/s`;
  if (mibPerSecond > 0) return `${Math.round(mibPerSecond * 1024)} KiB/s`;
  return "0 KiB/s";
}

function readLogValue(line, key){
  const match = line.match(new RegExp(`(?:^|\\s)${key}=("(?:\\\\.|[^"])*"|[^\\s]+)`));
  if (!match) return "";
  if (match[1].startsWith('"')) {
    try { return JSON.parse(match[1]); } catch (e) { return match[1].slice(1, -1); }
  }
  return match[1];
}

function formatLogForDisplay(log, language = lang){
  return String(log || "").split("\n").map(line => {
    if (line.includes("event=workflow.step.manual_ready")) {
      const path = readLogValue(line, "path");
      const partition = readLogValue(line, "partition");
      return language === "zh"
        ? `[5/5] Patch 完成 · 目标分区 ${partition}\nPatched 镜像：${path}`
        : `[5/5] Patch complete · Target ${partition}\nPatched image: ${path}`;
    }
    if (line.includes("event=workflow.manual_flash_command")) {
      const command = readLogValue(line, "command");
      return language === "zh" ? `手动刷写：${command}` : `Manual flash: ${command}`;
    }
    if (!line.includes("event=rom.download_progress")) return line;
    const readNumber = key => {
      const value = readLogValue(line, key);
      return value ? Number(value) : NaN;
    };
    const downloaded = readNumber("downloaded_mib");
    const total = readNumber("total_mib");
    const percent = readNumber("percent");
    const speed = readNumber("mib_per_second");
    if (![downloaded, total, percent, speed].every(Number.isFinite)) return line;

    const safePercent = Math.max(0, Math.min(100, percent));
    const width = 20;
    const filled = Math.round(safePercent / 100 * width);
    const bar = `${"█".repeat(filled)}${"░".repeat(width - filled)}`;
    const label = language === "zh" ? "下载 ROM" : "Downloading ROM";
    return `${label} [${bar}] ${safePercent.toFixed(1)}% · ${formatTransfer(downloaded, total)} · ${formatSpeed(speed)}`;
  }).join("\n");
}

function applyLang(l){
  lang = l;
  document.documentElement.lang = l === "zh" ? "zh-CN" : "en";
  document.querySelectorAll("[data-i18n]").forEach(el => {
    const k = el.getAttribute("data-i18n");
    if (I18N[l][k] != null) el.textContent = I18N[l][k];
  });
  document.querySelectorAll("[data-i18n-title]").forEach(el => {
    const k = el.getAttribute("data-i18n-title");
    if (I18N[l][k] != null) el.title = I18N[l][k];
  });
  document.querySelectorAll("[data-i18n-placeholder]").forEach(el => {
    const k = el.getAttribute("data-i18n-placeholder");
    if (I18N[l][k] != null) el.placeholder = I18N[l][k];
  });
  const ak = autoScroll ? "log.autoScroll.on" : "log.autoScroll.off";
  $("autoScrollText").setAttribute("data-i18n", ak);
  $("autoScrollText").textContent = t(ak);
  renderLogExpanded();
  renderDeviceOptions();
  renderDevicePill();
  if (detected && lastDevice) fillDevice(lastDevice);
  updateSummary();
}

function currentStepKeys(){
  const keys = flowFlash ? [...FLASH_STEP_KEYS] : [...PREPARE_STEP_KEYS];
  if (flowFlash && flowNoReboot) keys[keys.length - 1] = "step.done";
  return keys;
}

function renderStepper(n = 1){
  const keys = currentStepKeys();
  const total = keys.length;
  const stepper = $("stepper");
  stepper.className = `stepper ${flowFlash ? "flash-flow" : "prepare-flow"}`;
  stepper.innerHTML = `<div class="track" style="left:${50 / total}%;right:${50 / total}%"><i id="trackFill"></i></div>` +
    keys.map((key, index) => `<div class="step" data-step="${index + 1}"><div class="node">${index + 1}</div><div class="name" data-i18n="${key}">${esc(t(key))}</div></div>`).join("");
  $("stepTotal").textContent = total;
  setStep(Math.min(n, total + 1));
}

function setStep(n){
  const keys = currentStepKeys();
  const total = keys.length;
  const done = n > total;
  document.querySelectorAll("#stepper .step").forEach(el => {
    const k = Number(el.dataset.step);
    el.classList.toggle("done", done || k < n);
    el.classList.toggle("active", !done && k === n);
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
  $("trackFill").style.width = (done ? 100 : (n - 1) / 7 * 100) + "%";
}

function renderDevicePill(){
  const pill = $("devicePill"), txt = $("deviceText");
  if (deviceCount > 0) {
    pill.className = "pill on";
    txt.removeAttribute("data-i18n");
    txt.textContent = lang === "zh"
      ? `${deviceCount} 台设备已连接`
      : `${deviceCount} device${deviceCount > 1 ? "s" : ""} connected`;
  } else {
    pill.className = "pill off";
    txt.setAttribute("data-i18n", "device.disconnected");
    txt.textContent = t("device.disconnected");
  }
}

function renderDeviceOptions(){
  const sel = $("deviceSelect"), prev = sel.value;
  sel.innerHTML = deviceList.length
    ? deviceList.map(d => `<option value="${esc(d.Serial)}">${esc(d.Serial)} · ${esc(d.Mode)}</option>`).join("")
    : `<option value="">${esc(t("device.none"))}</option>`;
  if (prev) sel.value = prev;
}

function fillDevice(d){
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

function showEmpty(){
  detected = false; lastDevice = null;
  $("deviceEmpty").style.display = "";
  $("deviceData").style.display = "none";
  $("deviceTag").setAttribute("data-i18n", "device.none");
  $("deviceTag").textContent = t("device.none");
  $("sumRom").removeAttribute("data-i18n");
  $("sumRom").textContent = "—";
  updateSummary();
}

async function refreshDevices(){
  try {
    deviceList = await go().Devices() || [];
  } catch (e) {
    deviceList = [];
  }
  deviceCount = deviceList.filter(d => d.Mode === "adb").length || deviceList.length;
  renderDeviceOptions();
  renderDevicePill();
  if (deviceList.length) {
    const sel = $("deviceSelect");
    const adb = deviceList.find(d => d.Mode === "adb");
    const pick = deviceList.find(d => d.Serial === sel.value) || adb || deviceList[0];
    sel.value = pick.Serial;
    await inspect(pick.Serial);
  } else {
    showEmpty();
    setStep(1);
  }
}

async function inspect(serial){
  if (!serial) { showEmpty(); setStep(1); return; }
  try {
    const d = await go().Inspect(serial);
    lastDevice = d; detected = true;
    fillDevice(d);
    updateSummary();
    setStep(2);
  } catch (e) {
    showEmpty();
    setStep(2);
  }
}

function currentMode(){ return settings.mode || "lkm"; }

function updateSummary(){
  const mode = currentMode();
  $("sumMode").textContent = mode === "gki" ? "GKI" : "LKM";
  const part = settings.partition || "auto";
  if (part !== "auto") {
    $("sumPartition").textContent = part;
  } else if (lastDevice && lastDevice.HasInitBoot) {
    $("sumPartition").textContent = (mode === "gki" || String(lastDevice.KMI).startsWith("android12-")) ? "boot" : "init_boot";
  } else if (lastDevice) {
    $("sumPartition").textContent = "boot";
  } else {
    $("sumPartition").textContent = "—";
  }
  const flash = !!settings.autoFlash;
  const el = $("sumFlash");
  const key = flash
    ? (settings.noReboot ? "summary.flash.yes.noreboot" : "summary.flash.yes")
    : "summary.flash.no";
  el.setAttribute("data-i18n", key);
  el.textContent = t(key);
  const connections = Math.max(1, Math.min(32, Number(settings.connections) || 16));
  const custom = String(settings.mirrors || "").trim() !== "";
  $("sumDownload").textContent = lang === "zh"
    ? `${connections} 连接 · ${custom ? "自定义优先" : "自动镜像"}`
    : `${connections} connections · ${custom ? "custom first" : "auto mirrors"}`;
}

function renderROMSelection(){
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

function setRunning(running){
  ["btnStart","btnPrepare","refresh","deviceSelect","chooseKernel","chooseOTA","clearOTA","verifyROM","setAutoFlash","setNoReboot","setMirrors","setConnections"].forEach(id => {
    const el = $(id); if (el) el.disabled = running;
  });
  $("runBar").style.display = running ? "" : "none";
}

/* ---- settings ---- */
function setMode(mode){
  $("kernelField").style.display = mode === "gki" ? "" : "none";
}

function renderAutoScroll(){
  $("autoScroll").classList.toggle("on", autoScroll);
  const k = autoScroll ? "log.autoScroll.on" : "log.autoScroll.off";
  $("autoScrollText").setAttribute("data-i18n", k);
  $("autoScrollText").textContent = t(k);
}

function renderLogExpanded(){
  const key = logExpanded ? "log.collapse" : "log.expand";
  $("logBody").hidden = !logExpanded;
  $("logToggle").setAttribute("aria-expanded", String(logExpanded));
  $("logToggleText").setAttribute("data-i18n", key);
  $("logToggleText").textContent = t(key);
  document.querySelector(".console").classList.toggle("collapsed", !logExpanded);
  if (logExpanded && autoScroll) requestAnimationFrame(scrollLog);
}

function applySettings(){
  setMode(settings.mode || "lkm");
  $("setAutoFlash").checked = !!settings.autoFlash;
  $("setNoReboot").checked = !!settings.noReboot;
  autoScroll = settings.logAutoScroll !== false;
  renderAutoScroll();
  if (settings.language && settings.language !== lang) applyLang(settings.language);
  flowFlash = !!settings.autoFlash;
  flowNoReboot = !!settings.noReboot;
  renderStepper(detected ? 2 : 1);
  updateSummary();
}

function syncFromUI(){
  settings.language = lang;
}

function persistSettings(){
  syncFromUI();
  clearTimeout(saveTimer);
  saveTimer = setTimeout(() => { go().SaveSettings(settings).catch(() => {}); }, 300);
}

async function loadSettings(){
  try {
    const s = await go().GetSettings();
    if (s) settings = Object.assign(settings, s);
  } catch (e) {}
  applySettings();
}

function openSettings(){
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
  $("settingsModal").classList.remove("hidden");
}

function closeSettings(){ $("settingsModal").classList.add("hidden"); }

async function saveSettingsForm(){
  const next = {
    outputDir: $("setOutputDir").value.trim(),
    language: $("setLanguage").value,
    mode: $("setMode").value,
    partition: $("setPartition").value,
    autoFlash: $("setAutoFlash").checked,
    noReboot: $("setNoReboot").checked,
    catalogUrl: $("setCatalog").value.trim(),
    mirrors: $("setMirrors").value,
    connections: Math.max(1, Math.min(32, Number($("setConnections").value) || 16)),
    ksud: $("setKsud").value.trim(),
    logAutoScroll: $("setLogScroll").checked
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

function scrollLog(){ const el = $("log"); el.scrollTop = el.scrollHeight; }

function updateFlowFromLog(log){
  let backendStage = 1, backendTotal = flowFlash ? 8 : 5, m;
  const re = /(?:^|\s)step=(\d+)(?:\s|.*?\s)total=(\d+)/g;
  while ((m = re.exec(log)) !== null) {
    backendStage = Number(m[1]);
    backendTotal = Number(m[2]);
  }
  let stage;
  if (flowFlash) {
    stage = backendStage <= 1 ? 1
      : backendStage <= 3 ? 2
      : backendStage === 4 ? 3
      : backendStage === 5 ? 4
      : backendStage === 6 ? 5
      : backendStage === 7 ? 6
      : 8;
    if (log.includes("event=flash.completed")) stage = Math.max(stage, 7);
    if (log.includes("event=flash.rebooting")) stage = 8;
  } else {
    stage = backendTotal === 5 ? backendStage
      : backendStage <= 1 ? 1
        : backendStage <= 3 ? 2
        : backendStage === 4 ? 3
        : backendStage === 5 ? 4
        : 5;
  }
  setStep(stage);
}

function renderState(s){
  const st = $("logStatus");
  st.className = "status" + (s.status && s.status !== "idle" ? " " + s.status : "");
  const key = "status." + (s.status || "idle");
  st.setAttribute("data-i18n", I18N[lang][key] ? key : "status.idle");
  st.textContent = t(I18N[lang][key] ? key : "status.idle");

  const logLanguage = s.language === "en" ? "en" : "zh";
  const errorPrefix = I18N[logLanguage]["error.prefix"];
  const text = formatLogForDisplay(s.log || "", logLanguage) + (s.error ? "\n" + errorPrefix + s.error : "");
  $("log").textContent = text;
  updateFlowFromLog(s.log || "");
  if (autoScroll) scrollLog();

  const running = s.status === "running";
  setRunning(running);

  if (s.status === "succeeded" && s.result){
    if (s.result.Partition) $("sumPartition").textContent = s.result.Partition;
    if (s.result.ROMPath) {
      $("sumRom").removeAttribute("data-i18n");
      $("sumRom").textContent = baseName(s.result.ROMPath);
    }
    setStep(9);
  }
}

async function poll(){
  try {
    const s = await go().State();
    renderState(s);
    if (s.status === "running") setTimeout(poll, 600);
  } catch (e) {
    setTimeout(poll, 1200);
  }
}

async function start(flash){
  if (flash) {
    let ok = false;
    try {
      ok = await go().Confirm(t("confirm.title"), t("confirm.flash"), t("confirm.ok"), t("confirm.cancel"));
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
    NoReboot: !!settings.noReboot
  };
  flowFlash = flash;
  flowNoReboot = !!settings.noReboot;
  renderStepper(detected ? 2 : 1);
  setRunning(true);
  $("log").textContent = "";
  clearTimeout(saveTimer);
  syncFromUI();
  go().SaveSettings(settings).catch(() => {});
  try {
    await go().Start(req);
    poll();
  } catch (e) {
    setRunning(false);
    $("log").textContent = t("error.prefix") + e;
  }
}

/* ---- events ---- */
$("deviceSelect").onchange = () => inspect($("deviceSelect").value);
$("refresh").onclick = refreshDevices;
$("btnDetect").onclick = refreshDevices;
$("chooseKernel").onclick = async () => { const v = await go().ChooseFile(t("field.kernel")); if (v) $("kernel").value = v; };
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
  $("sumRom").textContent = "—";
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
$("logToggle").onclick = () => {
  logExpanded = !logExpanded;
  renderLogExpanded();
};
$("docsLink").onclick = e => {
  e.preventDefault();
  if (window.runtime && window.runtime.BrowserOpenURL) window.runtime.BrowserOpenURL("https://kernelsu.org/guide/installation.html");
};

$("openSettings").onclick = openSettings;
$("closeSettings").onclick = closeSettings;
$("cancelSettings").onclick = closeSettings;
$("saveSettings").onclick = saveSettingsForm;
$("setChooseOutput").onclick = async () => { const v = await go().ChooseOutputDirectory(); if (v) $("setOutputDir").value = v; };
$("setChooseKsud").onclick = async () => { const v = await go().ChooseFile(t("settings.chooseKsud")); if (v) $("setKsud").value = v; };
$("settingsModal").onclick = e => { if (e.target === $("settingsModal")) closeSettings(); };
document.addEventListener("keydown", e => {
  if (e.key === "Escape" && !$("settingsModal").classList.contains("hidden")) closeSettings();
});

/* ---- init ---- */
applyLang("zh");
renderStepper(1);
renderROMSelection();
loadSettings();
refreshDevices();
poll();
