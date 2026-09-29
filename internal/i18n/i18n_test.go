package i18n

import (
	"errors"
	"testing"
)

func TestEventTranslations(t *testing.T) {
	if got := Text("flash.completed", Chinese); got != "刷写成功" {
		t.Fatalf("Chinese translation = %q", got)
	}
	if got := Text("flash.completed", English); got != "Flash completed successfully" {
		t.Fatalf("English translation = %q", got)
	}
}

func TestLocalizedWrappedError(t *testing.T) {
	err := WrapError("error.flash_partition", errors.New("transport closed"), "init_boot")
	if got := ErrorText(err, Chinese); got != "刷写 init_boot 失败: transport closed" {
		t.Fatalf("Chinese error = %q", got)
	}
	if got := ErrorText(err, English); got != "Failed to flash init_boot: transport closed" {
		t.Fatalf("English error = %q", got)
	}
}

func TestUnknownMessageFallsBackToKey(t *testing.T) {
	if got := Text("unknown.event", English); got != "unknown.event" {
		t.Fatalf("unknown message = %q", got)
	}
}
