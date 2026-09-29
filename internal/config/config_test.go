package config

import (
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	file := filepath.Join(t.TempDir(), "nested", "config.json")
	want := Config{OutputDir: "/tmp/artifacts", Language: "en"}
	if err := want.SaveTo(file); err != nil {
		t.Fatal(err)
	}
	got := LoadFrom(file)
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
func TestLoadMissingReturnsEmpty(t *testing.T) {
	if got := LoadFrom(filepath.Join(t.TempDir(), "absent.json")); got != (Config{}) {
		t.Fatalf("got %+v", got)
	}
}

func TestOutputDirOrDefault(t *testing.T) {
	if got := (Config{}).OutputDirOrDefault(); got == "" {
		t.Fatal("expected non-empty default")
	}
	if got := (Config{OutputDir: "/custom"}).OutputDirOrDefault(); got != "/custom" {
		t.Fatalf("got %q", got)
	}
}

func TestLanguageOrDefault(t *testing.T) {
	if got := (Config{Language: "en"}).LanguageOrDefault("zh"); got != "en" {
		t.Fatalf("got %q", got)
	}
	if got := (Config{Language: "fr"}).LanguageOrDefault("zh"); got != "zh" {
		t.Fatalf("got %q", got)
	}
}

func TestDefaults(t *testing.T) {
	empty := Config{}
	if empty.ModeOrDefault() != "lkm" {
		t.Fatalf("mode default = %q", empty.ModeOrDefault())
	}
	if empty.PartitionOrDefault() != "auto" {
		t.Fatalf("partition default = %q", empty.PartitionOrDefault())
	}
	if !empty.LogAutoScrollDefault() {
		t.Fatal("log auto-scroll should default to true")
	}
	if got := empty.CatalogURLOr("default"); got != "default" {
		t.Fatalf("catalog default = %q", got)
	}
	off := false
	if (Config{LogAutoScroll: &off}).LogAutoScrollDefault() {
		t.Fatal("explicit false should win")
	}
	if (Config{Mode: "bogus"}).ModeOrDefault() != "lkm" {
		t.Fatal("invalid mode should fall back")
	}
	if (Config{Partition: "bogus"}).PartitionOrDefault() != "auto" {
		t.Fatal("invalid partition should fall back")
	}
}
