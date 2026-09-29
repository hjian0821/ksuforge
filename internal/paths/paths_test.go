package paths

import (
	"path/filepath"
	"testing"
)

func TestLayoutDirs(t *testing.T) {
	root := "/out"
	if got, want := ROMDir(root, "dada", "OS1"), filepath.Join(root, "rom", "dada", "OS1"); got != want {
		t.Fatalf("ROMDir = %q, want %q", got, want)
	}
	if got, want := StockDir(root, "dada", "OS1"), filepath.Join(root, "stock", "dada", "OS1"); got != want {
		t.Fatalf("StockDir = %q, want %q", got, want)
	}
	if got, want := PatchedDir(root, "dada", "OS1"), filepath.Join(root, "patched", "dada", "OS1"); got != want {
		t.Fatalf("PatchedDir = %q, want %q", got, want)
	}
}

func TestLayoutSkipsEmptySegments(t *testing.T) {
	root := "/out"
	if got, want := ROMDir(root, "", ""), filepath.Join(root, "rom"); got != want {
		t.Fatalf("ROMDir = %q, want %q", got, want)
	}
	if got, want := PatchedDir(root, "", ""), filepath.Join(root, "patched"); got != want {
		t.Fatalf("PatchedDir = %q, want %q", got, want)
	}
}

func TestLayoutSanitizesSegments(t *testing.T) {
	got := ROMDir("/out", "../evil", "a/b:c")
	want := filepath.Join("/out", "rom", "_evil", "a_b_c")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestPatchedDirForManagedStock(t *testing.T) {
	root := "/out"
	stock := filepath.Join(root, "stock", "dada", "OS1", "init_boot.img")
	if got, want := PatchedDirFor(root, stock), filepath.Join(root, "patched", "dada", "OS1"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestPatchedDirForForeignStock(t *testing.T) {
	if got, want := PatchedDirFor("/out", "/tmp/custom.img"), filepath.Join("/out", "patched"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
