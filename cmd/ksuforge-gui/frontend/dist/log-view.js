// The viewer consumes the task transcript; navigation never recreates the task UI.
const LogView = (() => {
  let records = [],
    source = "",
    signature = "";
  let level = "all",
    stage = "all",
    query = "",
    activeStage = "other";
  let actionKey = "",
    actionDetail = "";
  const expanded = new Set();
  const stages = ["device", "rom", "patch", "flash"];

  function classify(event, previous) {
    if (
      /acquire_stock|stock_verified|stock_acquired/.test(event) ||
      event.startsWith("rom.")
    )
      return "rom";
    if (/patch|manual_ready/.test(event) || event.startsWith("ksud."))
      return "patch";
    if (/flash|fastboot|reboot|validate_flash/.test(event)) return "flash";
    if (/check_usb|detect_device/.test(event) || event.startsWith("device."))
      return "device";
    return previous;
  }

  function parse(log, error = "") {
    let previous = "other";
    const result = String(log || "")
      .split(/\r?\n/)
      .filter((line) => line.trim())
      .map((raw, index) => {
        const event = readLogValue(raw, "event");
        previous = classify(event, previous);
        const stamp = readLogValue(raw, "time");
        const time = stamp.match(/T(\d{2}:\d{2}:\d{2})/i)?.[1] || "—";
        let severity = readLogValue(raw, "level").toUpperCase();
        if (!["DEBUG", "INFO", "WARN", "ERROR"].includes(severity))
          severity = "INFO";
        let message = readLogValue(raw, "msg") || raw;
        if (event === "rom.download_progress")
          message = formatLogForDisplay(raw);
        const command = readLogValue(raw, "command");
        if (command) message += ": " + command;
        return {
          id: index,
          raw,
          event,
          stage: previous,
          level: severity,
          time,
          stamp,
          message,
        };
      });
    if (error)
      result.push({
        id: result.length,
        raw: `level=ERROR event=task.failed msg=${JSON.stringify(error)}`,
        event: "task.failed",
        stage: previous,
        level: "ERROR",
        time: "—",
        stamp: "",
        message: error,
      });
    return result;
  }

  function rawText() {
    return (
      (lastState?.log || "") + (lastState?.error ? "\n" + lastState.error : "")
    );
  }

  function action(key, detail = "") {
    actionKey = key;
    actionDetail = detail;
    $("logActionStatus").textContent = key
      ? t(key) + (detail ? ": " + detail : "")
      : "";
  }

  function render(state = lastState || { status: "idle" }) {
    action(actionKey, actionDetail);
    const next = (state.log || "") + "\u0000" + (state.error || "");
    const changed = source !== next;
    if (changed) {
      if (!(state.log || "").startsWith(source.split("\u0000")[0]))
        expanded.clear();
      source = next;
      records = parse(state.log, state.error);
    }
    activeStage =
      records.filter((r) => stages.includes(r.stage)).at(-1)?.stage || "other";
    $("logCountAll").textContent = records.length;
    $("logStageAll").textContent = records.length;
    $("logCountWarn").textContent = records.filter(
      (r) => r.level === "WARN",
    ).length;
    $("logCountError").textContent = records.filter(
      (r) => r.level === "ERROR",
    ).length;
    $("exportLog").disabled = $("copyLog").disabled = !rawText().trim();
    $("logTaskTitle").textContent = t(
      taskContext?.flash ? "logs.flashTask" : "logs.prepareTask",
    );
    const context = taskContext || (records.length ? lastDevice : null);
    $("logTaskMeta").textContent = context
      ? [
          context.Model,
          context.Version,
          context.partition || state.result?.Partition,
        ]
          .filter(Boolean)
          .join(" · ")
      : t(records.length ? "logs.currentTask" : "logs.noTask");
    document.querySelectorAll("[data-log-stage]").forEach((button) => {
      const id = button.dataset.logStage;
      if (id === "all") return;
      const seen = records.some((r) => r.stage === id);
      let phase = !seen ? "pending" : "passed";
      if (seen && id === activeStage && state.status === "running")
        phase = "running";
      if (records.some((r) => r.stage === id && r.level === "ERROR"))
        phase = "failed";
      button.dataset.phase = phase;
      button.querySelector(".stage-state").textContent = t(
        "logs.phase." + phase,
      );
    });
    const nextSignature = JSON.stringify([source, level, stage, query, lang]);
    if (signature === nextSignature) return;
    signature = nextSignature;
    const shown = records.filter(
      (r) =>
        (level === "all" || r.level === level) &&
        (stage === "all" || r.stage === stage) &&
        (!query ||
          (r.message + " " + r.raw)
            .toLowerCase()
            .includes(query.toLowerCase())),
    );
    const scroll = $("log").scrollTop;
    $("log").innerHTML = shown
      .map(
        (r) =>
          `<div class="log-item ${r.level.toLowerCase()}" role="listitem"><button class="log-row" type="button" data-log-record="${r.id}" aria-expanded="${expanded.has(r.id)}" aria-controls="log-detail-${r.id}"><time class="row-time" ${r.stamp ? `datetime="${esc(r.stamp)}"` : ""}>${esc(r.time)}</time><span class="row-level">${esc(r.level)}</span><span class="row-message">${esc(r.message)}</span></button><pre class="row-detail" id="log-detail-${r.id}" ${expanded.has(r.id) ? "" : "hidden"}>${esc(r.raw)}</pre></div>`,
      )
      .join("");
    $("logEmpty").hidden = shown.length > 0;
    $("logEmpty").textContent = t(
      records.length ? "logs.noMatches" : "log.idle",
    );
    $("logRecordCount").textContent =
      lang === "zh"
        ? `${shown.length} 条记录 / 共 ${records.length} 条`
        : `${shown.length} of ${records.length} records`;
    $("log").scrollTop = scroll;
    if (changed && autoScroll && currentPage === "logs")
      requestAnimationFrame(scrollLog);
  }

  function reset() {
    source = signature = "";
    records = [];
    expanded.clear();
    level = stage = "all";
    query = $("logSearch").value = "";
    action("");
    syncFilters();
    render({ status: "running", log: "" });
  }

  function syncFilters() {
    document
      .querySelectorAll("[data-log-level], [data-log-stage]")
      .forEach((button) => {
        const selected = button.dataset.logLevel
          ? button.dataset.logLevel === level
          : button.dataset.logStage === stage;
        button.classList.toggle("selected", selected);
        button.setAttribute("aria-pressed", String(selected));
      });
  }

  function init() {
    $("logSearch").addEventListener("input", (event) => {
      query = event.target.value;
      render();
    });
    document.querySelectorAll("[data-log-level]").forEach(
      (button) =>
        (button.onclick = () => {
          level = button.dataset.logLevel;
          syncFilters();
          render();
        }),
    );
    document.querySelectorAll("[data-log-stage]").forEach(
      (button) =>
        (button.onclick = () => {
          stage = button.dataset.logStage;
          syncFilters();
          render();
        }),
    );
    $("log").addEventListener("click", (event) => {
      const button = event.target.closest("[data-log-record]");
      if (!button) return;
      const id = Number(button.dataset.logRecord);
      if (expanded.has(id)) expanded.delete(id);
      else expanded.add(id);
      button.setAttribute("aria-expanded", String(expanded.has(id)));
      $("log-detail-" + id).hidden = !expanded.has(id);
    });
    $("copyLog").onclick = async () => {
      try {
        if (window.runtime?.ClipboardSetText)
          await window.runtime.ClipboardSetText(rawText());
        else await navigator.clipboard.writeText(rawText());
        action("logs.copied");
      } catch (e) {
        action("logs.copyFailed");
      }
    };
    $("exportLog").onclick = async () => {
      action("");
      try {
        if (window.go?.main?.App?.ExportRunLog) {
          const path = await go().ExportRunLog();
          if (path) action("logs.exported");
        } else {
          const url = URL.createObjectURL(
            new Blob([rawText()], { type: "text/plain;charset=utf-8" }),
          );
          const link = document.createElement("a");
          link.href = url;
          link.download = "ksuforge-run.log";
          link.click();
          setTimeout(() => URL.revokeObjectURL(url), 1000);
        }
      } catch (e) {
        action("logs.exportFailed", String(e));
      }
    };
  }
  return { init, render, reset, parse };
})();
