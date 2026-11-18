const StorageView = (() => {
  let tab = "general",
    report = null,
    selected = new Set(),
    busy = "",
    result = null;
  let scanGeneration = 0;
  const key = (name) => t("storage." + name);
  const bytes = (n) => {
    const value = Math.max(0, Number(n) || 0);
    const units = ["B", "KiB", "MiB", "GiB", "TiB"];
    const i = value
      ? Math.min(4, Math.floor(Math.log(value) / Math.log(1024)))
      : 0;
    return `${(value / 1024 ** i).toFixed(i > 1 ? 2 : 0)} ${units[i]}`;
  };
  function setTab(next) {
    tab = next === "storage" ? "storage" : "general";
    const storage = tab === "storage";
    $("settingsModal")
      .querySelector(".modal")
      .classList.toggle("storage-modal", storage);
    $("generalSettings").classList.toggle("hidden", storage);
    $("storageSettings").classList.toggle("hidden", !storage);
    $("generalSettingsFooter").classList.toggle("hidden", storage);
    $("storageFooter").classList.toggle("hidden", !storage);
    $("settingsGeneralTab").setAttribute("aria-selected", String(!storage));
    $("settingsStorageTab").setAttribute("aria-selected", String(storage));
    if (storage) scan();
  }
  function chosen() {
    return (report?.items || []).filter((item) => selected.has(item.id));
  }
  function sync() {
    const items = chosen(),
      total = items.reduce((sum, item) => sum + item.bytes, 0);
    $("storageSelected").textContent = bytes(total);
    $("storageSelectionCount").textContent = items.length
      ? lang === "zh"
        ? `已选 ${items.length} 项 · 预计 ${bytes(total)}`
        : `${items.length} selected · approx. ${bytes(total)}`
      : key("none");
    $("storageClean").disabled = !items.length || !!busy || isRunning;
    $("storageRefresh").disabled = !!busy || isRunning;
    $("storageClean").textContent = key(
      busy === "clean" ? "processing" : "clean",
    );
    $("settingsStorageTab").disabled = isRunning || !!busy;
    document.querySelectorAll("[data-storage-item]").forEach((input) => {
      input.checked = selected.has(input.dataset.storageItem);
      input.disabled = !!busy || isRunning;
    });
    document.querySelectorAll("[data-storage-group]").forEach((input) => {
      const group = (report?.items || []).filter(
        (item) =>
          item.device + "/" + item.version === input.dataset.storageGroup,
      );
      const count = group.filter((item) => selected.has(item.id)).length;
      input.checked = !!count && count === group.length;
      input.indeterminate = count > 0 && count < group.length;
      input.disabled = !!busy || isRunning;
    });
    $("storageWarning").classList.toggle(
      "hidden",
      !items.some((item) => item.kind === "stock" || item.kind === "patched"),
    );
  }
  function render() {
    $("storageRoot").textContent = report?.root || "";
    $("storageRoot").title = report?.root || "";
    $("storageTotal").textContent = report ? bytes(report.totalBytes) : "—";
    const groupMap = new Map();
    for (const item of report?.items || []) {
      const groupKey = item.device + "/" + item.version;
      if (!groupMap.has(groupKey)) groupMap.set(groupKey, []);
      groupMap.get(groupKey).push(item);
    }
    const open = new Set(
      [...document.querySelectorAll(".storage-group[open]")].map(
        (el) => el.dataset.group,
      ),
    );
    $("storageGroups").innerHTML = [...groupMap]
      .map(
        ([name, items], index) =>
          `<details class="storage-group" data-group="${esc(name)}" ${open.has(name) || index === 0 ? "open" : ""}><summary><input type="checkbox" data-storage-group="${esc(name)}" aria-label="${esc(key("selectGroup"))}"><span><strong>${esc(items[0].device)}</strong><code>${esc(items[0].version)}</code></span><span class="storage-group-size">${bytes(items.reduce((sum, item) => sum + item.bytes, 0))}</span></summary><div>${items.map((item) => `<div class="storage-file"><label><input type="checkbox" data-storage-item="${esc(item.id)}"><span class="storage-file-info"><strong>${esc(key("kind." + item.kind))}</strong><span class="storage-filename" title="${esc(item.files[0].path)}">${esc(item.name)}</span><small>${esc(key(item.kind === "rom" ? (item.recommended ? "romReady" : "romUnknown") : "note." + item.kind))}</small></span><span class="storage-file-size">${bytes(item.bytes)}</span></label><details class="storage-files"><summary>${esc(key("files"))} · ${item.files.length}</summary>${item.files.map((file) => `<div><code>${esc(file.path)}</code><span>${bytes(file.bytes)} · ${esc(new Date(file.modified).toLocaleString(lang === "zh" ? "zh-CN" : "en-US"))}</span></div>`).join("")}</details></div>`).join("")}</div></details>`,
      )
      .join("");
    $("storageState").textContent =
      busy === "scan"
        ? key("scanning")
        : report && !report.items.length
          ? key("empty")
          : "";
    document.querySelectorAll("[data-storage-item]").forEach(
      (input) =>
        (input.onchange = () => {
          if (input.checked) selected.add(input.dataset.storageItem);
          else selected.delete(input.dataset.storageItem);
          sync();
        }),
    );
    document.querySelectorAll("[data-storage-group]").forEach(
      (input) =>
        (input.onchange = () => {
          for (const item of report.items.filter(
            (item) =>
              item.device + "/" + item.version === input.dataset.storageGroup,
          )) {
            if (input.checked) selected.add(item.id);
            else selected.delete(item.id);
          }
          sync();
        }),
    );
    const resultBox = $("storageResult");
    resultBox.classList.toggle("hidden", !result);
    if (result) {
      const message = result.error
        ? result.error
        : result.cancelled
          ? key("cancelled")
          : lang === "zh"
            ? `已删除 ${result.deletedFiles} 个文件（占用约 ${bytes(result.deletedBytes)}）。`
            : `Deleted ${result.deletedFiles} files (approx. ${bytes(result.deletedBytes)} allocated).`;
      resultBox.textContent = message;
      if (result.failed?.length) {
        const details = document.createElement("details");
        const summary = document.createElement("summary");
        summary.textContent = key("failed") + ` (${result.failed.length})`;
        details.append(summary);
        for (const failure of result.failed) {
          const p = document.createElement("p");
          p.textContent = failure.path + ": " + failure.error;
          details.append(p);
        }
        resultBox.append(details);
      }
      resultBox.classList.toggle(
        "warning",
        !!result.error || !!result.failed?.length,
      );
    }
    sync();
  }
  async function scan(preserveResult = false) {
    if (busy || isRunning) return;
    const generation = ++scanGeneration;
    busy = "scan";
    report = null;
    selected.clear();
    if (!preserveResult) result = null;
    render();
    try {
      if (!go().ScanStorage) throw new Error(key("desktopOnly"));
      const data = await go().ScanStorage();
      if (generation !== scanGeneration) return;
      report = data;
      selected = new Set(
        data.items.filter((item) => item.recommended).map((item) => item.id),
      );
    } catch (e) {
      result = { error: String(e) };
    } finally {
      if (generation === scanGeneration) {
        busy = "";
        render();
      }
    }
  }
  async function clean() {
    if (busy || isRunning || !report || !selected.size) return;
    busy = "clean";
    result = null;
    render();
    try {
      // The backend revalidates this scan and shows one native confirmation.
      result = await go().CleanStorage(report.token, [...selected]);
      busy = "";
      if (!result.cancelled) await scan(true);
    } catch (e) {
      result = { error: String(e) };
      report = null;
      selected.clear();
    } finally {
      busy = "";
      render();
    }
  }
  function init() {
    $("settingsGeneralTab").onclick = () => setTab("general");
    $("settingsStorageTab").onclick = () => setTab("storage");
    $("storageRefresh").onclick = () => scan();
    $("storageClean").onclick = clean;
    $("settingsStorageTab").onkeydown = $("settingsGeneralTab").onkeydown = (
      event,
    ) => {
      if (event.key === "ArrowLeft" || event.key === "ArrowRight") {
        event.preventDefault();
        const next = tab === "storage" ? "general" : "storage";
        if (next === "storage" && (isRunning || busy)) return;
        setTab(next);
        $(
          next === "storage" ? "settingsStorageTab" : "settingsGeneralTab",
        ).focus();
      }
    };
  }
  return { init, setTab, sync, render };
})();
