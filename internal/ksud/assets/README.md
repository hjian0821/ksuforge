# Bundled ksud assets

This directory holds the official KernelSU `ksud` binaries that are embedded
into `ksuforge-cli` and `ksuforge` with `go:embed`.

The binaries are ignored by git. `make tools` (or a Wails pre-build hook) runs
`go run ./cmd/fetchtools`, which downloads and verifies the binary for the
current `GOOS/GOARCH` if it is missing:

- `darwin/arm64`  -> `ksud-aarch64-apple-darwin`
- `darwin/amd64`  -> `ksud-x86_64-apple-darwin`
- `linux/arm64`   -> `ksud-aarch64-unknown-linux-musl`
- `linux/amd64`   -> `ksud-x86_64-unknown-linux-musl`
- `windows/amd64` -> `ksud-x86_64-pc-windows-gnu.exe`

If no asset matches the build platform the application still compiles and falls
back to an external `ksud` found on `PATH`.
