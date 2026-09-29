package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// AppName is the directory name used under the OS application data locations.
const AppName = "KSUForge"

// DefaultOutputDir is the user-visible directory for large artifacts such as the
// ROM, stock boot images, patched images and backups. Keeping them in a visible
// location makes them easy to find, reuse and delete.
func DefaultOutputDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "ksuforge-output"
	}
	return filepath.Join(home, "Downloads", AppName)
}

// ConfigDir is the per-user configuration directory:
// macOS ~/Library/Application Support/KSUForge, Windows %AppData%\KSUForge.
func ConfigDir() string {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, AppName)
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".config", AppName)
	}
	return AppName
}

// CacheDir is the per-user cache directory for regenerable data:
// macOS ~/Library/Caches/KSUForge, Windows %LocalAppData%\KSUForge.
func CacheDir() string {
	if dir, err := os.UserCacheDir(); err == nil && dir != "" {
		return filepath.Join(dir, AppName)
	}
	return filepath.Join(os.TempDir(), AppName)
}

// LogDir is the per-user log directory.
func LogDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(CacheDir(), "logs")
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Logs", AppName)
	case "windows":
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			return filepath.Join(local, AppName, "Logs")
		}
		return filepath.Join(home, AppName, "Logs")
	default:
		return filepath.Join(CacheDir(), "logs")
	}
}

// The artifact directory is split by type and then by device codename and system
// version so different devices and builds never overwrite each other:
//
//	<root>/rom/<codename>/<version>/
//	<root>/stock/<codename>/<version>/
//	<root>/patched/<codename>/<version>/
//
// Empty segments are skipped, so ROMDir(root, "", "") is <root>/rom.
func ROMDir(root, codename, version string) string {
	return joinLayout(root, "rom", codename, version)
}

// StockDir is where extracted stock boot images and their provenance live.
func StockDir(root, codename, version string) string {
	return joinLayout(root, "stock", codename, version)
}

// PatchedDir is where ksud writes the patched boot images.
func PatchedDir(root, codename, version string) string {
	return joinLayout(root, "patched", codename, version)
}

// PatchedDirFor mirrors a stock image's device/version under patched/. When the
// image does not follow the managed layout it falls back to <root>/patched.
func PatchedDirFor(root, stockImage string) string {
	if codename, version, ok := splitStockPath(stockImage); ok {
		return PatchedDir(root, codename, version)
	}
	return PatchedDir(root, "", "")
}

// splitStockPath recognises .../stock/<codename>/<version>/<file>.
func splitStockPath(path string) (codename, version string, ok bool) {
	dir := filepath.Dir(filepath.Clean(path))
	version = filepath.Base(dir)
	codename = filepath.Base(filepath.Dir(dir))
	parent := filepath.Base(filepath.Dir(filepath.Dir(dir)))
	if parent != "stock" || version == "" || version == "." || codename == "" || codename == "." {
		return "", "", false
	}
	return codename, version, true
}

func joinLayout(root string, parts ...string) string {
	segments := []string{root}
	for _, part := range parts {
		if part = sanitizeSegment(part); part != "" {
			segments = append(segments, part)
		}
	}
	return filepath.Join(segments...)
}

// sanitizeSegment keeps device and version strings usable as a single path
// segment on every platform.
func sanitizeSegment(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	v = strings.Map(func(r rune) rune {
		if r < 0x20 || strings.ContainsRune(`/\:*?"<>|`, r) {
			return '_'
		}
		return r
	}, v)
	v = strings.Trim(v, ".")
	if v == "" {
		return ""
	}
	return v
}
