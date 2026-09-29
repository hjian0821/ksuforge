// Package payloaddumper locates the payload-dumper-go tool used to extract
// boot images from payload.bin OTA packages. When no external copy is found it
// falls back to the official binary embedded in the application.
package payloaddumper

import (
	"embed"
	"runtime"

	"github.com/hjian0821/ksuforge/internal/bundled"
	"github.com/hjian0821/ksuforge/internal/command"
	"github.com/hjian0821/ksuforge/internal/i18n"
)

// Version is the bundled payload-dumper-go release. Keep it in sync with
// PAYLOAD_VERSION in the Makefile.
const Version = "v0.1.6"

//go:embed assets
var assets embed.FS

// assetNames maps GOOS/GOARCH to the official release asset name (the archive
// stem, which is also the binary name inside the archive).
var assetNames = map[string]string{
	"darwin/arm64":  "payload-dumper-darwin-arm64",
	"darwin/amd64":  "payload-dumper-darwin-amd64-v3",
	"linux/arm64":   "payload-dumper-linux-arm64",
	"linux/amd64":   "payload-dumper-linux-amd64-v3",
	"windows/amd64": "payload-dumper-windows-amd64-v3.exe",
	"windows/arm64": "payload-dumper-windows-arm64.exe",
}

// AssetName returns the bundled asset name for the current platform, or an
// empty string when the platform has no bundled binary.
func AssetName() string {
	return assetNames[runtime.GOOS+"/"+runtime.GOARCH]
}

// Embedded reports whether this build carries payload-dumper for the platform.
func Embedded() bool {
	name := AssetName()
	if name == "" {
		return false
	}
	return bundled.Exists(assets, "assets/"+name)
}

// Extract writes the embedded payload-dumper to the cache directory and returns
// its path.
func Extract() (string, error) {
	name := AssetName()
	if name == "" {
		return "", i18n.NewError("error.bundled_payload_dumper_unavailable", runtime.GOOS, runtime.GOARCH)
	}
	return bundled.Extract(assets, "assets/"+name, "payload-dumper", Version)
}

// Resolve returns the payload-dumper executable to use: the copy on PATH
// (which includes the one shipped next to the application) and finally the
// binary embedded in the application.
func Resolve() (string, error) {
	runner := command.ExecRunner{SearchPaths: command.DesktopSearchPaths()}
	if tool, err := runner.LookPath("payload-dumper"); err == nil {
		return tool, nil
	}
	if Embedded() {
		tool, err := Extract()
		if err != nil {
			return "", i18n.WrapError("error.payload_dumper_extract", err)
		}
		return tool, nil
	}
	return "", i18n.NewError("error.payload_dumper_missing")
}
