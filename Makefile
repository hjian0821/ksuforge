.PHONY: build cli gui tools release release-all test clean

WAILS ?= $(shell go env GOPATH)/bin/wails
WAILS_FLAGS ?=
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
HOST_GOOS := $(shell go env GOHOSTOS)
HOST_GOARCH := $(shell go env GOHOSTARCH)
VERSION ?= $(shell git describe --tags --always --dirty)
EXE := $(if $(filter windows,$(GOOS)),.exe,)
PROJECT_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BIN_DIR := $(PROJECT_ROOT)/bin
GUI_DIR := $(PROJECT_ROOT)/cmd/ksuforge-gui
GUI_BUILD_BIN := $(GUI_DIR)/build/bin
# Wails resolves -o relative to its internal build/bin directory.
WAILS_OUTPUT := ../../../../bin/ksuforge$(EXE)

build: cli gui

# Download, verify and stage the bundled tools for the target platform. Existing
# files are reused, so this is a no-op once the assets are present.
tools:
	cd "$(PROJECT_ROOT)" && CGO_ENABLED=0 GOOS="$(HOST_GOOS)" GOARCH="$(HOST_GOARCH)" go run ./cmd/fetchtools -goos $(GOOS) -goarch $(GOARCH)

cli: tools
	mkdir -p "$(BIN_DIR)"
	cd "$(PROJECT_ROOT)" && go build -o "$(BIN_DIR)/ksuforge-cli$(EXE)" ./cmd/ksuforge
	cp "$(PROJECT_ROOT)/third_party/bin/ksud$(EXE)" "$(BIN_DIR)/"
	cp "$(PROJECT_ROOT)/third_party/bin/payload-dumper$(EXE)" "$(BIN_DIR)/"
	cp "$(PROJECT_ROOT)/third_party/licenses/KernelSU-GPL-3.0.txt" "$(BIN_DIR)/KernelSU-GPL-3.0.txt"

gui: tools
	mkdir -p "$(BIN_DIR)"
	cd "$(GUI_DIR)" && "$(WAILS)" build $(WAILS_FLAGS) -o "$(WAILS_OUTPUT)"
	@if [ "$(GOOS)" = "darwin" ]; then \
		mkdir -p "$(GUI_BUILD_BIN)/ksuforge.app/Contents/Resources/tools" "$(GUI_BUILD_BIN)/ksuforge.app/Contents/Resources/licenses"; \
		cp "$(PROJECT_ROOT)/third_party/bin/ksud" "$(GUI_BUILD_BIN)/ksuforge.app/Contents/Resources/tools/ksud"; \
		cp "$(PROJECT_ROOT)/third_party/bin/payload-dumper" "$(GUI_BUILD_BIN)/ksuforge.app/Contents/Resources/tools/payload-dumper"; \
		cp "$(PROJECT_ROOT)/third_party/licenses/KernelSU-GPL-3.0.txt" "$(GUI_BUILD_BIN)/ksuforge.app/Contents/Resources/licenses/KernelSU-GPL-3.0.txt"; \
		codesign --force --deep --sign - "$(GUI_BUILD_BIN)/ksuforge.app"; \
		rm -rf "$(BIN_DIR)/ksuforge.app"; \
		cp -R "$(GUI_BUILD_BIN)/ksuforge.app" "$(BIN_DIR)/"; \
	else \
		cp "$(PROJECT_ROOT)/third_party/bin/ksud$(EXE)" "$(BIN_DIR)/"; \
		cp "$(PROJECT_ROOT)/third_party/bin/payload-dumper$(EXE)" "$(BIN_DIR)/"; \
		cp "$(PROJECT_ROOT)/third_party/licenses/KernelSU-GPL-3.0.txt" "$(BIN_DIR)/KernelSU-GPL-3.0.txt"; \
	fi

# Build and package the current target for distribution.
release: build
	cd "$(PROJECT_ROOT)" && CGO_ENABLED=0 GOOS="$(HOST_GOOS)" GOARCH="$(HOST_GOARCH)" go run ./cmd/package-release -bin "$(BIN_DIR)" -version "$(VERSION)" -goos "$(GOOS)" -goarch "$(GOARCH)"

# Build macOS natively and build Linux/Windows with the local Docker toolchain.
# This only writes archives to bin/; it never uploads them.
release-all:
	@printf '%s\n' "$(VERSION)" | grep -Eq '^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z]+([.-][0-9A-Za-z]+)*)?$$' || (echo "VERSION must be semantic and start with v (for example: v1.2.3 or v1.2.3-rc.1)" && exit 1)
	bash "$(PROJECT_ROOT)/scripts/release-all.sh" "$(VERSION)"

test:
	cd "$(PROJECT_ROOT)" && go test ./...
	cd "$(PROJECT_ROOT)" && go vet ./...

clean:
	cd "$(PROJECT_ROOT)" && go clean
