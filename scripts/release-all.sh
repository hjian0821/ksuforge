#!/usr/bin/env bash
set -euo pipefail

version="${1:-}"
project_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
builder_image="ksuforge-release-builder:local"

if [[ "$(uname -s)" != "Darwin" ]]; then
  echo "release-all must run on macOS because Wails macOS applications require the native toolchain." >&2
  exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker is required to build the Linux and Windows packages." >&2
  exit 1
fi
if ! docker info >/dev/null 2>&1; then
  echo "Docker Desktop is not running. Start it and retry." >&2
  exit 1
fi

host_arch="$(go env GOARCH)"
case "$host_arch" in
  arm64|amd64) ;;
  *)
    echo "Unsupported macOS architecture: $host_arch" >&2
    exit 1
    ;;
esac

echo "==> Building macOS/$host_arch"
make -C "$project_root" release \
  VERSION="$version" \
  GOOS=darwin \
  GOARCH="$host_arch" \
  WAILS_FLAGS="-platform darwin/$host_arch"

echo "==> Preparing Linux/Windows build container"
docker build \
  --platform linux/amd64 \
  --tag "$builder_image" \
  --file "$project_root/scripts/release/Dockerfile" \
  "$project_root/scripts/release"

docker_make() {
  docker run --rm \
    --platform linux/amd64 \
    --user "$(id -u):$(id -g)" \
    --env HOME=/tmp/ksuforge-build-home \
    --env GOCACHE=/tmp/ksuforge-build-home/go-cache \
    --env GOMODCACHE=/tmp/ksuforge-build-home/go-mod \
    --volume "$project_root:/src" \
    --workdir /src \
    "$builder_image" \
    "$@"
}

echo "==> Building linux/amd64"
docker_make make release \
  VERSION="$version" \
  GOOS=linux \
  GOARCH=amd64 \
  WAILS_FLAGS="-platform linux/amd64 -tags webkit2_41"

echo "==> Building windows/amd64"
docker_make env CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc make release \
  VERSION="$version" \
  GOOS=windows \
  GOARCH=amd64 \
  WAILS_FLAGS="-platform windows/amd64"

echo "==> Release archives"
find "$project_root/bin" -maxdepth 1 -type f \
  \( -name "ksuforge_*.tar.gz" -o -name "ksuforge_*.zip" \) \
  -print | sort
