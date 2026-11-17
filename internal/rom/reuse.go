package rom

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	bootimage "github.com/hjian0821/ksuforge/internal/image"
)

// FindVerifiedStock only reuses a previously verified, exact-build stock image.
// Unlike extraction's legacy cache, it requires complete provenance and SHA-256
// because its source ROM may already have been cleaned up.
func FindVerifiedStock(dir, partition, version string) (string, Release, bool) {
	if partition != "boot" && partition != "init_boot" {
		return "", Release{}, false
	}
	path := filepath.Join(dir, partition+".img")
	for _, file := range []string{path, path + ".json"} {
		info, err := os.Lstat(file)
		if err != nil || !info.Mode().IsRegular() {
			return "", Release{}, false
		}
	}
	data, err := os.ReadFile(path + ".json")
	if err != nil {
		return "", Release{}, false
	}
	var state extractState
	if json.Unmarshal(data, &state) != nil || version == "" || state.Version != version {
		return "", Release{}, false
	}
	checksum, err := hex.DecodeString(state.MD5)
	if err != nil || len(checksum) != 16 || len(state.SHA256) != 64 {
		return "", Release{}, false
	}
	info, err := bootimage.Inspect(path)
	if err != nil || !strings.EqualFold(info.SHA256, state.SHA256) {
		return "", Release{}, false
	}
	return path, Release{Version: version, MD5: strings.ToLower(state.MD5)}, true
}
