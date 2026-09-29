package rom

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractDirectBootImage(t *testing.T) {
	dir := t.TempDir()
	romPath := filepath.Join(dir, "rom.zip")
	f, err := os.Create(romPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	entry, err := zw.Create("firmware-update/boot.img")
	if err != nil {
		t.Fatal(err)
	}
	image := make([]byte, 4096)
	copy(image, []byte("ANDROID!"))
	if _, err := entry.Write(image); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	got, err := ExtractImage(context.Background(), romPath, "boot", out, Release{})
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.Join(out, "boot.img") {
		t.Fatalf("got %q", got)
	}
}

func TestExtractReusesMatchingImage(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out")
	if err := os.MkdirAll(out, 0o755); err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(out, "init_boot.img")
	image := make([]byte, 4096)
	copy(image, []byte("ANDROID!"))
	if err := os.WriteFile(destination, image, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination+".json", []byte(`{"version":"OS1","md5":"abc"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ExtractImage(context.Background(), filepath.Join(dir, "missing.zip"), "init_boot", out, Release{Version: "OS1", MD5: "abc"})
	if err != nil {
		t.Fatal(err)
	}
	if got != destination {
		t.Fatalf("got %q", got)
	}
}

func TestExtractReextractsOnVersionMismatch(t *testing.T) {
	dir := t.TempDir()
	romPath := filepath.Join(dir, "rom.zip")
	f, err := os.Create(romPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	entry, err := zw.Create("boot.img")
	if err != nil {
		t.Fatal(err)
	}
	image := make([]byte, 4096)
	copy(image, []byte("ANDROID!"))
	if _, err := entry.Write(image); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	if err := os.MkdirAll(out, 0o755); err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(out, "boot.img")
	if err := os.WriteFile(destination, image, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination+".json", []byte(`{"version":"OS1","md5":"old"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ExtractImage(context.Background(), romPath, "boot", out, Release{Version: "OS2", MD5: "new"})
	if err != nil {
		t.Fatal(err)
	}
	if got != destination {
		t.Fatalf("got %q", got)
	}
	if data, err := os.ReadFile(destination + ".json"); err != nil || !strings.Contains(string(data), "OS2") {
		t.Fatalf("sidecar not refreshed: %s, %v", data, err)
	}
}
