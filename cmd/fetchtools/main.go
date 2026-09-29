// Command fetchtools downloads the official KernelSU ksud and payload-dumper-go
// binaries for a target platform, verifies their SHA-256 and stages them so the
// application can embed them. Existing files are reused, so the network is only
// touched when a tool is missing or has the wrong hash.
package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hjian0821/ksuforge/internal/i18n"
	"github.com/hjian0821/ksuforge/internal/logging"
)

const (
	ksudVersion    = "v3.3.0"
	payloadVersion = "v0.1.6"
	licenseURL     = "https://raw.githubusercontent.com/tiann/KernelSU/" + ksudVersion + "/LICENSE"
)

type artifact struct {
	name    string // display name for logs
	url     string
	sha256  string
	archive string // "", "zip" or "tar.gz"
	inner   string // file name inside the archive
	asset   string // destination under internal/<pkg>/assets
	bin     string // normalized copy under third_party/bin
}

type ksudEntry struct{ asset, sha string }

var ksudTable = map[string]ksudEntry{
	"darwin/arm64":  {"ksud-aarch64-apple-darwin", "40ca97a2fb61284129909abac5325dcae790736d9b88901f4f31cc7ec6d9a705"},
	"darwin/amd64":  {"ksud-x86_64-apple-darwin", "82a85db50e88d97a5f963c3ce29929d0ed7c8d817fbd1ad92d3a960e82bfbb34"},
	"linux/arm64":   {"ksud-aarch64-unknown-linux-musl", "2bb58f3bc15fd0058e248e808a4676da7ef293d03d495cde29db125aebffd01b"},
	"linux/amd64":   {"ksud-x86_64-unknown-linux-musl", "36c4f9350501e5aef3ef582a2bc581c1f48e2fea40c830d34319f0160a4bd790"},
	"windows/amd64": {"ksud-x86_64-pc-windows-gnu.exe", "f9880fbbb3751d847dd1fbe3ea5939f592c28887ff274cd5ddb090ed6c751ee0"},
}

type payloadEntry struct{ asset, sha, inner string }

var payloadTable = map[string]payloadEntry{
	"darwin/arm64":  {"payload-dumper-darwin-arm64.tar.gz", "ba88cb2f33b6e302eb5d87347aba4dfb12d875d167a8e657fe231b93251de89d", "payload-dumper-darwin-arm64"},
	"darwin/amd64":  {"payload-dumper-darwin-amd64-v3.tar.gz", "be93349c7782bd3995f6444d965b9c0a5489bb8a8d5bc7715de3e2838592a1e2", "payload-dumper-darwin-amd64-v3"},
	"linux/arm64":   {"payload-dumper-linux-arm64.tar.gz", "9b11cb88dff8b7e455469a52dddc798bf11b765fdd87772e29df611a61f0ded1", "payload-dumper-linux-arm64"},
	"linux/amd64":   {"payload-dumper-linux-amd64-v3.tar.gz", "c2960706e7f8d6e5a7f9b42ec55b7997828120de47a3b8a9de93c2cb7dc44503", "payload-dumper-linux-amd64-v3"},
	"windows/amd64": {"payload-dumper-windows-amd64-v3.zip", "14277c0ad855827e1058f36622b5ea0fdc6e44fcfd542f3293b80eff59e8465e", "payload-dumper-windows-amd64-v3.exe"},
	"windows/arm64": {"payload-dumper-windows-arm64.zip", "0fe11b0b88afa3a996278fd58b473d22116f4949bdaac5978d655fc0674b2b5d", "payload-dumper-windows-arm64.exe"},
}

func main() {
	logger, closeLog := logging.Setup(logging.Options{Level: slog.LevelInfo, Stderr: true, Prefix: "ksuforge-fetchtools"})
	defer closeLog()

	goos := flag.String("goos", envOr("GOOS", runtime.GOOS), "target GOOS")
	goarch := flag.String("goarch", envOr("GOARCH", runtime.GOARCH), "target GOARCH")
	flag.Parse()

	if err := run(*goos, *goarch); err != nil {
		logger.Error("fetch.failed", "error", err)
		os.Exit(1)
	}
}

func run(goos, goarch string) error {
	root, err := moduleRoot()
	if err != nil {
		return err
	}
	arts, err := manifest(goos, goarch)
	if err != nil {
		return err
	}
	for _, a := range arts {
		if err := fetch(root, a); err != nil {
			return err
		}
	}
	return fetchLicense(root)
}

func manifest(goos, goarch string) ([]artifact, error) {
	key := goos + "/" + goarch
	exe := ""
	if goos == "windows" {
		exe = ".exe"
	}
	var arts []artifact
	if e, ok := ksudTable[key]; ok {
		arts = append(arts, artifact{
			name:   "ksud " + ksudVersion,
			url:    fmt.Sprintf("https://github.com/tiann/KernelSU/releases/download/%s/%s", ksudVersion, e.asset),
			sha256: e.sha,
			asset:  filepath.ToSlash(filepath.Join("internal", "ksud", "assets", e.asset)),
			bin:    filepath.ToSlash(filepath.Join("third_party", "bin", "ksud"+exe)),
		})
	}
	if e, ok := payloadTable[key]; ok {
		arts = append(arts, artifact{
			name:    "payload-dumper " + payloadVersion,
			url:     fmt.Sprintf("https://github.com/xishang0128/payload-dumper-go/releases/download/%s/%s", payloadVersion, e.asset),
			sha256:  e.sha,
			archive: archiveKind(e.asset),
			inner:   e.inner,
			asset:   filepath.ToSlash(filepath.Join("internal", "payloaddumper", "assets", e.inner)),
			bin:     filepath.ToSlash(filepath.Join("third_party", "bin", "payload-dumper"+exe)),
		})
	}
	if len(arts) == 0 {
		return nil, i18n.NewError("error.bundled_tool_unavailable", goos, goarch)
	}
	return arts, nil
}

func archiveKind(name string) string {
	switch {
	case strings.HasSuffix(name, ".zip"):
		return "zip"
	case strings.HasSuffix(name, ".tar.gz"):
		return "tar.gz"
	default:
		return ""
	}
}

func fetch(root string, a artifact) error {
	dest := filepath.Join(root, filepath.FromSlash(a.asset))
	reuse := false
	if info, err := os.Stat(dest); err == nil && !info.IsDir() && info.Size() > 0 {
		// For archives the manifest hash covers the archive, not the extracted
		// binary, so an existing file is trusted as-is.
		reuse = a.archive != "" || matchSHA256(dest, a.sha256)
	}
	if reuse {
		slog.Info("fetch.exists", "tool", a.name)
	} else {
		slog.Info("fetch.downloading", "tool", a.name, "url", a.url)
		data, err := download(a.url)
		if err != nil {
			return fmt.Errorf("%s: %w", a.name, err)
		}
		sum := sha256.Sum256(data)
		if got := hex.EncodeToString(sum[:]); got != a.sha256 {
			return i18n.NewError("error.sha256_mismatch", a.name, a.sha256, got)
		}
		if err := writeArtifact(dest, a, data); err != nil {
			return fmt.Errorf("%s: %w", a.name, err)
		}
		slog.Info("fetch.saved", "tool", a.name, "path", dest)
	}
	binDest := filepath.Join(root, filepath.FromSlash(a.bin))
	if err := copyFile(dest, binDest); err != nil {
		return fmt.Errorf("%s: %w", a.name, err)
	}
	return nil
}

func writeArtifact(dest string, a artifact, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	switch a.archive {
	case "zip":
		return extractZip(data, a.inner, dest)
	case "tar.gz":
		return extractTarGz(data, a.inner, dest)
	default:
		return writeFileAtomic(dest, data)
	}
}

func extractZip(data []byte, inner, dest string) error {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	for _, f := range zr.File {
		if filepath.Base(f.Name) != inner {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		return writeFileFrom(dest, rc)
	}
	return i18n.NewError("error.archive_entry_missing", inner)
}

func extractTarGz(data []byte, inner, dest string) error {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag != tar.TypeReg || filepath.Base(hdr.Name) != inner {
			continue
		}
		return writeFileFrom(dest, tr)
	}
	return i18n.NewError("error.archive_entry_missing", inner)
}

func fetchLicense(root string) error {
	dest := filepath.Join(root, "third_party", "licenses", "KernelSU-GPL-3.0.txt")
	if info, err := os.Stat(dest); err == nil && info.Size() > 0 {
		return nil
	}
	data, err := download(licenseURL)
	if err != nil {
		// The license is only needed when packaging; never block an embedded build.
		slog.Warn("fetch.license_failed", "error", err)
		return nil
	}
	return writeFileAtomic(dest, data)
}

func download(url string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func matchSHA256(path, want string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false
	}
	return hex.EncodeToString(h.Sum(nil)) == want
}

func writeFileAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o755); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Chmod(path, 0o755)
}

func writeFileFrom(path string, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Chmod(path, 0o755)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	return writeFileFrom(dst, in)
}

// moduleRoot walks up from the working directory to the module root (go.mod).
func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", i18n.NewError("error.module_root_missing")
		}
		dir = parent
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
