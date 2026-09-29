package ksud

import (
	"os"
	"testing"
)

func TestEmbeddedAndExtract(t *testing.T) {
	if !Embedded() {
		t.Skipf("no bundled ksud for this platform")
	}
	path, err := Extract()
	if err != nil {
		t.Fatalf("Extract: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat extracted ksud: %v", err)
	}
	if info.IsDir() || info.Size() == 0 {
		t.Fatalf("extracted ksud is not a usable file: %+v", info)
	}
	again, err := Extract()
	if err != nil {
		t.Fatalf("Extract (second call): %v", err)
	}
	if again != path {
		t.Fatalf("Extract is not stable: %q vs %q", path, again)
	}
}
