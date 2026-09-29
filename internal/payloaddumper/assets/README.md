# Bundled payload-dumper assets

This directory holds the official `payload-dumper-go` binaries that are embedded
into `ksuforge-cli` and `ksuforge` with `go:embed`.

The binaries are ignored by git. `make tools` (or a Wails pre-build hook) runs
`go run ./cmd/fetchtools`, which downloads and verifies the binary for the
current `GOOS/GOARCH` if it is missing:

- `darwin/arm64`  -> `payload-dumper-darwin-arm64`
- `darwin/amd64`  -> `payload-dumper-darwin-amd64-v3`
- `linux/arm64`   -> `payload-dumper-linux-arm64`
- `linux/amd64`   -> `payload-dumper-linux-amd64-v3`
- `windows/amd64` -> `payload-dumper-windows-amd64-v3.exe`
- `windows/arm64` -> `payload-dumper-windows-arm64.exe`

If no asset matches the build platform the application still compiles and falls
back to an external `payload-dumper` found on `PATH`.
