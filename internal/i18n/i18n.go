package i18n

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	Chinese = "zh"
	English = "en"
)

type Message struct {
	Zh string
	En string
}

// Normalize returns a supported language. Chinese remains the application
// fallback so existing first-run behaviour is preserved.
func Normalize(language string) string {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "en", "en-us", "en_us", "english":
		return English
	case "zh", "zh-cn", "zh_cn", "zh-hans", "chinese":
		return Chinese
	default:
		return Chinese
	}
}

// Detect returns the language requested by the process environment.
func Detect() string {
	for _, name := range []string{"KSUFORGE_LANG", "LC_ALL", "LC_MESSAGES", "LANG"} {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			return Normalize(strings.Split(value, ".")[0])
		}
	}
	return Chinese
}

func Lookup(key, language string) (string, bool) {
	message, ok := messages[key]
	if !ok {
		return key, false
	}
	if Normalize(language) == English {
		return message.En, true
	}
	return message.Zh, true
}

func Text(key, language string, args ...any) string {
	format, ok := Lookup(key, language)
	if !ok {
		return key
	}
	return fmt.Sprintf(format, args...)
}

// AppError keeps a stable error code and formatting arguments while preserving
// the underlying cause for errors.Is/errors.As and diagnostic file logs.
type AppError struct {
	Code  string
	Args  []any
	Cause error
}

func (e *AppError) Error() string { return ErrorText(e, Chinese) }
func (e *AppError) Unwrap() error { return e.Cause }

func NewError(code string, args ...any) error {
	return &AppError{Code: code, Args: args}
}

func WrapError(code string, cause error, args ...any) error {
	return &AppError{Code: code, Args: args, Cause: cause}
}

func ErrorText(err error, language string) string {
	if err == nil {
		return ""
	}
	if appErr, ok := err.(*AppError); ok {
		text := Text(appErr.Code, language, appErr.Args...)
		if appErr.Cause == nil {
			return text
		}
		return text + ": " + ErrorText(appErr.Cause, language)
	}
	if cause := errors.Unwrap(err); cause != nil {
		raw, rawCause := err.Error(), cause.Error()
		if strings.HasSuffix(raw, rawCause) {
			return strings.TrimSuffix(raw, rawCause) + ErrorText(cause, language)
		}
	}
	return err.Error()
}

var messages = map[string]Message{
	"app.started":              {Zh: "KSUForge 已启动", En: "KSUForge started"},
	"app.gui_started":          {Zh: "KSUForge GUI 已启动", En: "KSUForge GUI started"},
	"app.exit_error":           {Zh: "应用异常退出", En: "Application exited with an error"},
	"app.command_failed":       {Zh: "命令执行失败", En: "Command failed"},
	"app.settings_saved":       {Zh: "设置已保存", En: "Settings saved"},
	"app.task_started":         {Zh: "任务已开始", En: "Task started"},
	"app.task_failed":          {Zh: "任务失败", En: "Task failed"},
	"app.task_completed":       {Zh: "任务完成", En: "Task completed"},
	"logging.file_unavailable": {Zh: "无法写入日志文件，仅使用标准错误输出", En: "Unable to write log file; using stderr only"},

	"workflow.started":                  {Zh: "工作流已开始", En: "Workflow started"},
	"workflow.completed":                {Zh: "工作流已完成", En: "Workflow completed"},
	"workflow.step.check_usb":           {Zh: "检查 USB / ADB 连接", En: "Checking USB/ADB connection"},
	"workflow.step.detect_device":       {Zh: "检测设备", En: "Detecting device"},
	"workflow.step.acquire_stock":       {Zh: "读取设备信息并获取原厂镜像", En: "Reading device information and acquiring the stock image"},
	"workflow.step.stock_verified":      {Zh: "原厂启动镜像已获取并通过校验", En: "Stock boot image acquired and verified"},
	"workflow.step.patch":               {Zh: "应用 KernelSU Patch", En: "Applying KernelSU patch"},
	"workflow.stock_acquired":           {Zh: "已获取原厂启动镜像", En: "Stock boot image acquired"},
	"workflow.step.manual_ready":        {Zh: "Patch 完成，可以手动刷写", En: "Patched image ready for manual flashing"},
	"workflow.manual_flash_command":     {Zh: "手动刷写命令", En: "Manual flash command"},
	"workflow.step.reboot_fastboot":     {Zh: "准备重启到 Fastboot", En: "Preparing to reboot into fastboot"},
	"workflow.step.validate_flash":      {Zh: "刷写前校验 Slot 和 Bootloader 状态", En: "Validating slot and bootloader state before flashing"},
	"workflow.step.flash_stay_fastboot": {Zh: "刷写完成，设备停留在 Fastboot", En: "Flash completed; device remains in fastboot"},
	"workflow.step.flash_booting":       {Zh: "刷写完成，设备正在启动 Android", En: "Flash completed; device is booting Android"},

	"device.information":    {Zh: "设备信息", En: "Device information"},
	"device.target_image":   {Zh: "设备与目标镜像", En: "Device and target image"},
	"device.found":          {Zh: "发现设备", En: "Device found"},
	"device.none":           {Zh: "未发现设备", En: "No device found"},
	"tool.available":        {Zh: "工具可用", En: "Tool available"},
	"tool.missing_optional": {Zh: "未找到工具；完整自动化工作流需要该工具", En: "Tool not found; required for the complete automated workflow"},

	"rom.match_result":                {Zh: "ROM 匹配结果", En: "ROM match result"},
	"rom.local_selected":              {Zh: "使用选定的本地 Recovery ROM", En: "Using selected local recovery ROM"},
	"rom.local_md5_passed":            {Zh: "选定的本地 Recovery ROM 已通过 MD5 校验", En: "Selected local recovery ROM passed MD5 verification"},
	"rom.local_md5_skipped":           {Zh: "已跳过选定本地 Recovery ROM 的 MD5 校验", En: "Selected local recovery ROM MD5 verification skipped"},
	"rom.match_found":                 {Zh: "已找到匹配的 Recovery ROM", En: "Matching recovery ROM found"},
	"rom.verification_passed":         {Zh: "Recovery ROM 校验通过", En: "Recovery ROM verification passed"},
	"rom.local_reused":                {Zh: "本地 Recovery ROM 已通过 MD5 校验，跳过下载", En: "Matching local recovery ROM passed MD5 verification; skipping download"},
	"rom.local_md5_mismatch":          {Zh: "本地文件 MD5 与目录不一致，将重新下载", En: "Local file MD5 does not match the catalog; ignoring and downloading again"},
	"rom.mirror_probe_failed":         {Zh: "镜像源探测失败", En: "Mirror probe failed"},
	"rom.mirror_benchmark":            {Zh: "镜像源测速", En: "Mirror benchmark"},
	"rom.mirror_benchmark_completed":  {Zh: "镜像源测速完成", En: "Mirror benchmark completed"},
	"rom.download_resumed":            {Zh: "发现有效分片状态，继续下载", En: "Valid chunk state found; resuming download"},
	"rom.download_verifying":          {Zh: "下载完成，正在校验 MD5", En: "Download completed; verifying MD5"},
	"rom.download_progress":           {Zh: "下载 ROM", En: "Downloading ROM"},
	"rom.image_reused":                {Zh: "复用已校验的原厂镜像", En: "Reusing verified stock image"},
	"rom.image_reextract":             {Zh: "本地原厂镜像无效或版本不匹配，重新提取", En: "Local stock image is invalid or does not match the current version; extracting again"},
	"rom.image_extracted_direct":      {Zh: "已从 Recovery ROM 直接提取镜像", En: "Image extracted directly from recovery ROM"},
	"rom.image_extracted_payload":     {Zh: "已从 payload.bin 提取镜像", En: "Image extracted from payload.bin"},
	"rom.stock_extracted":             {Zh: "原厂镜像提取完成", En: "Stock image extracted"},
	"rom.stock_preparation_completed": {Zh: "原厂镜像准备完成", En: "Stock image preparation completed"},

	"patch.image_reused":  {Zh: "复用已校验的 Patched 镜像", En: "Reusing verified patched image"},
	"patch.image_stale":   {Zh: "本地 Patched 镜像与当前输入不匹配，重新 Patch", En: "Local patched image does not match the current input; patching again"},
	"patch.started":       {Zh: "开始执行 KernelSU Patch", En: "Starting KernelSU patch"},
	"patch.image_created": {Zh: "KernelSU Patched 镜像已生成", En: "KernelSU patched image created"},
	"patch.completed":     {Zh: "KernelSU Patch 完成", En: "KernelSU patch completed"},

	"flash.device_information": {Zh: "刷写设备信息", En: "Device information"},
	"flash.target":             {Zh: "刷写目标", En: "Flash target"},
	"flash.stock_backup":       {Zh: "原厂镜像备份", En: "Stock backup"},
	"flash.preview_completed":  {Zh: "预演完成，未写入设备；添加 --yes 执行刷写", En: "Preview completed without writing to the device; add --yes to flash"},
	"flash.completed":          {Zh: "刷写成功", En: "Flash completed successfully"},
	"flash.rebooting":          {Zh: "设备正在重启", En: "Device is rebooting"},
	"flash.workflow_completed": {Zh: "刷写工作流已完成", En: "Flash workflow completed"},
	"image.information":        {Zh: "启动镜像信息", En: "Boot image information"},

	"fetch.exists":         {Zh: "内置工具已存在，跳过", En: "Embedded tool already exists; skipping"},
	"fetch.downloading":    {Zh: "正在下载内置工具", En: "Downloading embedded tool"},
	"fetch.saved":          {Zh: "内置工具已保存", En: "Embedded tool saved"},
	"fetch.failed":         {Zh: "内置工具获取失败", En: "Embedded tool acquisition failed"},
	"fetch.license_failed": {Zh: "KernelSU 许可证下载失败", En: "Failed to download KernelSU license"},

	"error.device_required":                    {Zh: "未指定设备", En: "No device was specified"},
	"error.language_unsupported":               {Zh: "不支持的语言", En: "Unsupported language"},
	"error.mode_invalid":                       {Zh: "模式只能是 lkm 或 gki", En: "Mode must be lkm or gki"},
	"error.partition_invalid":                  {Zh: "分区只能是 auto、boot 或 init_boot", En: "Partition must be auto, boot, or init_boot"},
	"error.output_dir_required":                {Zh: "输出目录不能为空", En: "Output directory cannot be empty"},
	"error.connections_invalid":                {Zh: "并发连接数必须在 1 到 32 之间", En: "Connections must be between 1 and 32"},
	"error.mirror_invalid":                     {Zh: "镜像源 %q 无效", En: "Mirror %q is invalid"},
	"error.task_running":                       {Zh: "已有任务正在运行", En: "A task is already running"},
	"error.command_unknown":                    {Zh: "未知命令 %q", En: "Unknown command %q"},
	"error.image_required":                     {Zh: "必须提供 --image", En: "--image is required"},
	"error.image_and_backup_required":          {Zh: "必须提供 --image 和 --stock-backup", En: "--image and --stock-backup are required"},
	"error.tool_required":                      {Zh: "缺少 %s（请安装 Android SDK Platform-Tools）", En: "%s is missing (install Android SDK Platform-Tools)"},
	"error.command_not_found":                  {Zh: "找不到命令 %s", En: "Command %s was not found"},
	"error.bundled_tool_read":                  {Zh: "读取内置工具失败", En: "Failed to read bundled tool"},
	"error.bundled_tool_unavailable":           {Zh: "没有适用于 %s/%s 的内置工具", En: "No bundled tool is available for %s/%s"},
	"error.bundled_ksud_unavailable":           {Zh: "没有适用于 %s/%s 的内置 ksud", En: "No bundled ksud is available for %s/%s"},
	"error.bundled_payload_dumper_unavailable": {Zh: "没有适用于 %s/%s 的内置 payload-dumper", En: "No bundled payload-dumper is available for %s/%s"},
	"error.ksud_override_not_found":            {Zh: "找不到指定的 ksud", En: "The specified ksud was not found"},
	"error.ksud_extract":                       {Zh: "释放内置 ksud 失败", En: "Failed to extract bundled ksud"},
	"error.ksud_missing":                       {Zh: "缺少 ksud；请从 KernelSU 官方 Release 下载并加入 PATH", En: "ksud is missing; download it from the official KernelSU release and add it to PATH"},
	"error.payload_dumper_extract":             {Zh: "释放内置 payload-dumper 失败", En: "Failed to extract bundled payload-dumper"},
	"error.payload_dumper_missing":             {Zh: "缺少 payload-dumper；请从 payload-dumper-go 官方 Release 下载并加入 PATH", En: "payload-dumper is missing; download it from the official payload-dumper-go release and add it to PATH"},
	"error.archive_entry_missing":              {Zh: "压缩包中找不到 %s", En: "%s was not found in the archive"},
	"error.module_root_missing":                {Zh: "找不到模块根目录（go.mod）", En: "Module root (go.mod) was not found"},
	"error.sha256_mismatch":                    {Zh: "%s: SHA-256 校验失败，期望 %s，实际 %s", En: "%s: SHA-256 mismatch, expected %s, got %s"},

	"error.device_multiple_adb": {Zh: "发现多个 ADB 设备，请用 --serial 指定", En: "Multiple ADB devices were found; specify one with --serial"},
	"error.device_no_adb":       {Zh: "未找到可用的 ADB 设备", En: "No available ADB device was found"},
	"error.device_property":     {Zh: "读取设备属性 %s 失败", En: "Failed to read device property %s"},
	"error.device_sdk":          {Zh: "读取 Android SDK 失败", En: "Failed to read Android SDK"},
	"error.device_kernel":       {Zh: "读取内核版本失败", En: "Failed to read the kernel version"},
	"error.device_not_found":    {Zh: "未找到设备 %s", En: "Device %s was not found"},
	"error.device_no_transport": {Zh: "未找到 adb/fastboot 设备", En: "No adb/fastboot device was found"},
	"error.device_multiple":     {Zh: "发现多个设备，请用 --serial 指定", En: "Multiple devices were found; specify one with --serial"},

	"error.image_invalid":          {Zh: "不是有效的 Android 启动镜像（文件过小或非普通文件）", En: "Not a valid Android boot image (the file is too small or is not a regular file)"},
	"error.image_magic":            {Zh: "无法识别启动镜像 magic；拒绝刷写", En: "Boot image magic is not recognized; flashing was refused"},
	"error.stock_image_check":      {Zh: "检查原厂镜像失败", En: "Failed to inspect the stock image"},
	"error.patched_image_check":    {Zh: "检查 Patched 镜像失败", En: "Failed to inspect the patched image"},
	"error.flash_image_check":      {Zh: "检查待刷镜像失败", En: "Failed to inspect the image to flash"},
	"error.stock_backup_check":     {Zh: "检查原厂备份镜像失败", En: "Failed to inspect the stock backup image"},
	"error.flash_images_identical": {Zh: "待刷镜像与原厂备份完全相同", En: "The image to flash is identical to the stock backup"},
	"error.kmi_required":           {Zh: "无法识别 KMI，请用 --kmi 明确指定，例如 android14-6.1", En: "KMI could not be detected; specify --kmi explicitly, for example android14-6.1"},
	"error.kmi_from_kernel":        {Zh: "无法从内核版本 %q 识别 KMI，请用 --kmi 明确指定", En: "KMI could not be detected from kernel version %q; specify --kmi explicitly"},
	"error.gki_kernel_required":    {Zh: "GKI 替换内核模式必须提供与 KMI 和压缩格式匹配的 --kernel", En: "GKI kernel replacement mode requires --kernel matching the KMI and compression format"},
	"error.patch_failed":           {Zh: "KernelSU Patch 失败", En: "KernelSU patch failed"},

	"error.fastboot_preexisting": {Zh: "设备已在 Fastboot；为避免失去厂商/Android 信息，请先启动系统再运行", En: "The device is already in fastboot; boot Android first so vendor and Android information can be verified"},
	"error.oem_unsupported":      {Zh: "当前仅默认支持小米/Redmi/POCO；其他厂商需显式传入 --allow-other-oem", En: "Only Xiaomi/Redmi/POCO is supported by default; pass --allow-other-oem explicitly for other vendors"},
	"error.reboot_bootloader":    {Zh: "重启到 Bootloader 失败", En: "Failed to reboot to the bootloader"},
	"error.fastboot_timeout":     {Zh: "等待 Fastboot 设备超时", En: "Timed out waiting for the fastboot device"},
	"error.bootloader_locked":    {Zh: "无法确认 Bootloader 已解锁（%s），拒绝刷写", En: "Unable to confirm that the bootloader is unlocked (%s); flashing was refused"},
	"error.slot_mismatch":        {Zh: "ADB 检测到 Slot %s，但 Fastboot 当前 Slot 为 %s，拒绝刷写", En: "ADB detected slot %s but the current fastboot slot is %s; flashing was refused"},
	"error.flash_partition":      {Zh: "刷写 %s 失败", En: "Failed to flash %s"},
	"error.reboot_after_flash":   {Zh: "刷写成功，但重启失败", En: "Flashing succeeded, but the reboot failed"},

	"error.rom_redirect_limit":        {Zh: "ROM 下载重定向次数过多", En: "Too many ROM download redirects"},
	"error.rom_redirect_insecure":     {Zh: "ROM 镜像重定向到不安全地址", En: "The ROM mirror redirected to an insecure address"},
	"error.rom_oem_unsupported":       {Zh: "自动下载目前仅支持小米/Redmi/POCO 设备", En: "Automatic download currently supports only Xiaomi/Redmi/POCO devices"},
	"error.rom_local_path":            {Zh: "解析本地卡刷包路径失败", En: "Failed to resolve the local Recovery ROM path"},
	"error.rom_local_read":            {Zh: "读取本地卡刷包失败", En: "Failed to read the local Recovery ROM"},
	"error.rom_local_not_zip":         {Zh: "本地卡刷包必须是 ZIP 文件", En: "The local Recovery ROM must be a ZIP file"},
	"error.rom_url_not_zip":           {Zh: "下载地址不是 ZIP 卡刷包", En: "The download URL does not point to a ZIP Recovery ROM"},
	"error.mirror_https":              {Zh: "自定义镜像源必须是有效的 HTTPS 地址", En: "A custom mirror must be a valid HTTPS address"},
	"error.mirror_base_url":           {Zh: "镜像源必须是无账号、查询参数或片段的 HTTPS 基础地址", En: "A mirror must be an HTTPS base URL without credentials, query parameters, or fragments"},
	"error.mirrors_unavailable":       {Zh: "所有小米 OTA 镜像均不可用或不支持断点续传", En: "All Xiaomi OTA mirrors are unavailable or do not support resuming"},
	"error.rom_md5_parts":             {Zh: "卡刷包 MD5 不匹配：实际 %s，期望 %s；分片文件已保留", En: "Recovery ROM MD5 mismatch: got %s, expected %s; chunk files were retained"},
	"error.chunk_download":            {Zh: "分片 %d 下载失败", En: "Chunk %d download failed"},
	"error.content_range_missing":     {Zh: "响应缺少有效 Content-Range", En: "The response is missing a valid Content-Range"},
	"error.content_range_mismatch":    {Zh: "Content-Range 不匹配：%s", En: "Content-Range mismatch: %s"},
	"error.md5_invalid":               {Zh: "无效的 MD5 校验值", En: "Invalid MD5 checksum"},
	"error.rom_md5":                   {Zh: "卡刷包 MD5 不匹配：实际 %s，期望 %s", En: "Recovery ROM MD5 mismatch: got %s, expected %s"},
	"error.codename_invalid":          {Zh: "设备代号 %q 含有非法字符", En: "Device codename %q contains invalid characters"},
	"error.version_missing":           {Zh: "无法读取当前系统版本", En: "The current system version could not be read"},
	"error.catalog_read":              {Zh: "读取 ROM 索引失败", En: "Failed to read the ROM catalog"},
	"error.catalog_http":              {Zh: "读取 ROM 索引失败：HTTP %s", En: "Failed to read the ROM catalog: HTTP %s"},
	"error.catalog_parse":             {Zh: "解析 ROM 索引失败", En: "Failed to parse the ROM catalog"},
	"error.catalog_md5_missing":       {Zh: "找到同版卡刷包，但索引没有可验证的 MD5", En: "A matching Recovery ROM was found, but the catalog has no verifiable MD5"},
	"error.catalog_no_match":          {Zh: "索引中没有找到设备 %s 的同版 Recovery ROM %s；拒绝用其他版本替代", En: "No matching Recovery ROM was found in the catalog for device %s and version %s; another version will not be substituted"},
	"error.catalog_url_required":      {Zh: "ROM 索引地址不能为空", En: "The ROM catalog URL cannot be empty"},
	"error.catalog_placeholder":       {Zh: "ROM 索引地址必须且只能包含一个 %%s 占位符用于设备代号", En: "The ROM catalog URL must contain exactly one %%s placeholder for the device codename"},
	"error.catalog_https":             {Zh: "ROM 索引地址必须是有效的 HTTPS URL", En: "The ROM catalog URL must be a valid HTTPS URL"},
	"error.rom_url_untrusted":         {Zh: "ROM 下载地址不是可信 HTTPS URL", En: "The ROM download URL is not a trusted HTTPS URL"},
	"error.rom_cdn_untrusted":         {Zh: "ROM 下载地址不属于允许的小米 CDN：%s", En: "The ROM download URL does not belong to an allowed Xiaomi CDN: %s"},
	"error.extract_partition":         {Zh: "只能提取 boot 或 init_boot", En: "Only boot or init_boot can be extracted"},
	"error.rom_open":                  {Zh: "打开卡刷包失败", En: "Failed to open the Recovery ROM"},
	"error.extracted_image_check":     {Zh: "检查提取镜像失败", En: "Failed to inspect the extracted image"},
	"error.rom_image_missing":         {Zh: "卡刷包中既没有 %s.img，也没有 payload.bin", En: "The Recovery ROM contains neither %s.img nor payload.bin"},
	"error.payload_required":          {Zh: "该卡刷包使用 payload.bin", En: "This Recovery ROM uses payload.bin"},
	"error.payload_extract_partition": {Zh: "从 payload.bin 提取 %s 失败", En: "Failed to extract %s from payload.bin"},
}
