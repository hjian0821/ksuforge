package logging

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hjian0821/ksuforge/internal/i18n"
)

func TestWithWriterMirrorsStructuredRecord(t *testing.T) {
	var output bytes.Buffer
	logger := WithWriter(&output)
	logger.Info("download progress", "percent", 50)
	got := output.String()
	if !strings.Contains(got, "download progress") || !strings.Contains(got, "percent=50") {
		t.Fatalf("stream log content = %q", got)
	}
	if strings.Contains(got, "level=INFO") || strings.Contains(got, "time=") {
		t.Fatalf("stream log should omit level and time: %q", got)
	}
}

func TestWithWriterLocalizesEventsAndErrors(t *testing.T) {
	var chinese bytes.Buffer
	logger := WithWriter(&chinese, "zh")
	logger.Error("app.command_failed", "error", i18n.WrapError("error.stock_image_check", errors.New("disk failure")))
	got := chinese.String()
	for _, expected := range []string{"命令执行失败", "event=app.command_failed", "检查原厂镜像失败: disk failure"} {
		if !strings.Contains(got, expected) {
			t.Fatalf("Chinese stream log %q does not contain %q", got, expected)
		}
	}

	var english bytes.Buffer
	WithWriter(&english, "en").Info("workflow.step.check_usb", "step", 1, "total", 5)
	got = english.String()
	for _, expected := range []string{"Checking USB/ADB connection", "event=workflow.step.check_usb", "step=1", "total=5"} {
		if !strings.Contains(got, expected) {
			t.Fatalf("English stream log %q does not contain %q", got, expected)
		}
	}
}

func TestPersistentLogRemainsEnglish(t *testing.T) {
	dir := t.TempDir()
	_, closeLog := Setup(Options{Level: slog.LevelInfo, Dir: dir, Prefix: "localized"})
	var task bytes.Buffer
	WithWriter(&task, "zh").Info("workflow.step.patch", "step", 4, "total", 5)
	closeLog()

	if !strings.Contains(task.String(), "应用 KernelSU Patch") {
		t.Fatalf("task stream was not Chinese: %q", task.String())
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "Applying KernelSU patch") || strings.Contains(got, "应用 KernelSU Patch") {
		t.Fatalf("persistent log should remain English: %q", got)
	}
}

func TestSetupWritesToDailyFile(t *testing.T) {
	dir := t.TempDir()
	logger, closeLog := Setup(Options{Level: slog.LevelInfo, Dir: dir, Prefix: "test"})
	logger.Info("hello", "key", "value")
	closeLog()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one log file, got %d", len(entries))
	}
	if !strings.HasPrefix(entries[0].Name(), "test-") || filepath.Ext(entries[0].Name()) != ".log" {
		t.Fatalf("unexpected log name %q", entries[0].Name())
	}
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "hello") || !strings.Contains(string(data), "key=value") {
		t.Fatalf("log content = %q", data)
	}
}

func TestSetupFallsBackToStderr(t *testing.T) {
	// An unwritable directory must not panic or lose the logger.
	logger, closeLog := Setup(Options{Level: slog.LevelInfo, Dir: "/proc/does-not-exist/nope"})
	defer closeLog()
	if logger == nil {
		t.Fatal("expected a usable logger")
	}
	logger.Info("still works")
}
