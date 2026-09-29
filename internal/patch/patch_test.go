package patch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hjian0821/ksuforge/internal/ksud"
)

func TestPatchStateMatches(t *testing.T) {
	dir := t.TempDir()
	ksud := filepath.Join(dir, "ksud")
	if err := os.WriteFile(ksud, []byte("binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	state := filepath.Join(dir, "patched.img.json")
	if err := os.WriteFile(state, []byte(`{"stock_sha256":"aaa","kmi":"android14-6.1","mode":"lkm"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if !patchStateMatches(state, "aaa", "android14-6.1", "lkm", ksud) {
		t.Fatal("expected matching state to be reusable")
	}
	if patchStateMatches(state, "bbb", "android14-6.1", "lkm", ksud) {
		t.Fatal("different stock image must not be reused")
	}
	if patchStateMatches(state, "aaa", "android15-6.1", "lkm", ksud) {
		t.Fatal("different KMI must not be reused")
	}
	if patchStateMatches(filepath.Join(dir, "missing.json"), "aaa", "android14-6.1", "lkm", ksud) {
		t.Fatal("missing sidecar must not be reused")
	}
}

func TestResolveKSUdFallsBackToEmbedded(t *testing.T) {
	if !ksud.Embedded() {
		t.Skipf("no bundled ksud for this platform")
	}
	path, err := ResolveKSUd("")
	if err != nil {
		t.Fatalf("ResolveKSUd(empty): %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("resolved ksud not found: %v", err)
	}
}

func TestResolveKSUdRejectsMissingOverride(t *testing.T) {
	if _, err := ResolveKSUd("definitely-not-a-real-ksud-binary"); err == nil {
		t.Fatal("expected an error for a missing explicit path")
	}
}
