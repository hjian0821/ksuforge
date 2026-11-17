package rom

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFindVerifiedStockRequiresExactProvenance(t *testing.T) {
	for _, scenario := range []string{"valid", "wrong version", "unverified source", "missing metadata", "modified image"} {
		t.Run(scenario, func(t *testing.T) {
			dir := t.TempDir()
			data := make([]byte, 4096)
			copy(data, "ANDROID!")
			sum := sha256.Sum256(data)
			state := extractState{Version: "OS4", MD5: "0123456789abcdef0123456789abcdef", SHA256: hex.EncodeToString(sum[:])}
			switch scenario {
			case "wrong version":
				state.Version = "OS3"
			case "unverified source":
				state.MD5 = "unverified:123"
			case "modified image":
				data[100] = 1
			}
			file := filepath.Join(dir, "init_boot.img")
			if err := os.WriteFile(file, data, 0o600); err != nil {
				t.Fatal(err)
			}
			if scenario != "missing metadata" {
				meta, _ := json.Marshal(state)
				if err := os.WriteFile(file+".json", meta, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			path, release, ok := FindVerifiedStock(dir, "init_boot", "OS4")
			if ok != (scenario == "valid") {
				t.Fatalf("reused=%v", ok)
			}
			if ok && (path != file || release.Version != "OS4") {
				t.Fatalf("%s %+v", path, release)
			}
		})
	}
}
