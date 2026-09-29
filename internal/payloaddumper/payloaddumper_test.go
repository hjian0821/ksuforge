package payloaddumper

import (
	"os"
	"testing"
)

func TestEmbeddedAndExtract(t *testing.T) {
	if !Embedded() {
		t.Skipf("no bundled payload-dumper for this platform")
	}
	path, err := Extract()
	if err != nil {
		t.Fatalf("Extract: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat extracted payload-dumper: %v", err)
	}
	if info.IsDir() || info.Size() == 0 {
		t.Fatalf("extracted payload-dumper is not a usable file: %+v", info)
	}
}
