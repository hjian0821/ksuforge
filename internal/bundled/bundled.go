// Package bundled extracts executables that are embedded in the application
// with go:embed into the per-user cache directory so they can be run.
package bundled

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hjian0821/ksuforge/internal/i18n"
	appPaths "github.com/hjian0821/ksuforge/internal/paths"
)

// Exists reports whether the embedded filesystem contains assetPath.
func Exists(assets embed.FS, assetPath string) bool {
	_, err := fs.Stat(assets, assetPath)
	return err == nil
}

// Extract writes a content-addressed copy of an embedded executable to the
// per-user cache directory and returns its path. baseName is the file name to
// use on disk (without the platform .exe suffix); version and a content hash
// are appended so upgraded builds never reuse a stale binary.
func Extract(assets embed.FS, assetPath, baseName, version string) (string, error) {
	data, err := assets.ReadFile(assetPath)
	if err != nil {
		return "", i18n.WrapError("error.bundled_tool_read", err)
	}
	name := baseName
	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(name), ".exe") {
		name += ".exe"
	}
	sum := sha256.Sum256(data)
	dir := filepath.Join(appPaths.CacheDir(), "tools")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	dest := filepath.Join(dir, fmt.Sprintf("%s-%s-%s", name, version, hex.EncodeToString(sum[:])[:12]))
	if info, err := os.Stat(dest); err == nil && info.Size() == int64(len(data)) {
		return dest, nil
	}
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, data, 0o755); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, dest); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}
	_ = os.Chmod(dest, 0o755)
	return dest, nil
}
