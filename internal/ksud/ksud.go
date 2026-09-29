// Package ksud locates the KernelSU userspace tool. When no external ksud is
// found it falls back to the official binary embedded in the application,
// extracted to the per-user cache directory on first use.
package ksud

import (
	"embed"
	"runtime"

	"github.com/hjian0821/ksuforge/internal/bundled"
	"github.com/hjian0821/ksuforge/internal/i18n"
)

// Version is the bundled KernelSU release. Keep it in sync with KSUD_VERSION in
// the Makefile.
const Version = "v3.3.0"

//go:embed assets
var assets embed.FS

// assetNames maps GOOS/GOARCH to the official KernelSU release asset name.
var assetNames = map[string]string{
	"darwin/arm64":  "ksud-aarch64-apple-darwin",
	"darwin/amd64":  "ksud-x86_64-apple-darwin",
	"linux/arm64":   "ksud-aarch64-unknown-linux-musl",
	"linux/amd64":   "ksud-x86_64-unknown-linux-musl",
	"windows/amd64": "ksud-x86_64-pc-windows-gnu.exe",
}

// AssetName returns the bundled asset name for the current platform, or an
// empty string when the platform has no bundled binary.
func AssetName() string {
	return assetNames[runtime.GOOS+"/"+runtime.GOARCH]
}

// Embedded reports whether this build carries a ksud for the current platform.
func Embedded() bool {
	name := AssetName()
	if name == "" {
		return false
	}
	return bundled.Exists(assets, "assets/"+name)
}

// Extract writes the embedded ksud to the cache directory and returns it.
func Extract() (string, error) {
	name := AssetName()
	if name == "" {
		return "", i18n.NewError("error.bundled_ksud_unavailable", runtime.GOOS, runtime.GOARCH)
	}
	return bundled.Extract(assets, "assets/"+name, "ksud", Version)
}
