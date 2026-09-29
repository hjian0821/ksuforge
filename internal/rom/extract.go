package rom

import (
	"archive/zip"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hjian0821/ksuforge/internal/command"
	"github.com/hjian0821/ksuforge/internal/i18n"
	bootimage "github.com/hjian0821/ksuforge/internal/image"
	"github.com/hjian0821/ksuforge/internal/logging"
	"github.com/hjian0821/ksuforge/internal/payloaddumper"
)

type extractState struct {
	Version string `json:"version"`
	MD5     string `json:"md5"`
	SHA256  string `json:"sha256"`
}

func ExtractImage(ctx context.Context, romPath, partition, outputDir string, release Release, output ...io.Writer) (string, error) {
	if partition != "boot" && partition != "init_boot" {
		return "", i18n.NewError("error.extract_partition")
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}
	out := io.Discard
	if len(output) > 0 && output[0] != nil {
		out = output[0]
	}
	logger := logging.WithWriter(out, logging.Language(ctx))
	destination := filepath.Join(outputDir, partition+".img")
	statePath := destination + ".json"
	if _, err := os.Stat(destination); err == nil {
		if _, err := bootimage.Inspect(destination); err == nil && extractStateMatches(statePath, release) {
			logger.Info("rom.image_reused", "path", destination)
			return destination, nil
		}
		logger.Info("rom.image_reextract", "image", partition+".img", "path", destination)
		if err := os.Remove(destination); err != nil {
			return "", err
		}
		_ = os.Remove(statePath)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	zr, err := zip.OpenReader(romPath)
	if err != nil {
		return "", i18n.WrapError("error.rom_open", err)
	}
	hasPayload := false
	for _, entry := range zr.File {
		base := filepath.Base(strings.ReplaceAll(entry.Name, "\\", "/"))
		if base == "payload.bin" {
			hasPayload = true
		}
		if base != partition+".img" {
			continue
		}
		err := copyZipEntry(entry, destination)
		zr.Close()
		if err != nil {
			return "", err
		}
		if _, err := bootimage.Inspect(destination); err != nil {
			os.Remove(destination)
			return "", i18n.WrapError("error.extracted_image_check", err)
		}
		writeExtractState(statePath, destination, release)
		slog.Info("rom.image_extracted_direct", "partition", partition, "path", destination)
		return destination, nil
	}
	zr.Close()
	if !hasPayload {
		return "", i18n.NewError("error.rom_image_missing", partition)
	}
	tool, err := payloaddumper.Resolve()
	if err != nil {
		return "", i18n.WrapError("error.payload_required", err)
	}
	cmd := command.CommandContext(ctx, tool, "extract", romPath, "--partitions", partition, "--out", outputDir)
	cmd.Stdout, cmd.Stderr = out, out
	if err := cmd.Run(); err != nil {
		return "", i18n.WrapError("error.payload_extract_partition", err, partition)
	}
	if _, err := bootimage.Inspect(destination); err != nil {
		return "", i18n.WrapError("error.extracted_image_check", err)
	}
	writeExtractState(statePath, destination, release)
	slog.Info("rom.image_extracted_payload", "partition", partition, "path", destination)
	return destination, nil
}

// extractStateMatches reports whether an existing image can be reused for the
// requested release. A recorded version or MD5 that differs forces a re-extract;
// an unreadable or absent sidecar falls back to trusting a valid image.
func extractStateMatches(statePath string, release Release) bool {
	if release.Version == "" && release.MD5 == "" {
		return true
	}
	data, err := os.ReadFile(statePath)
	if err != nil {
		return true
	}
	var state extractState
	if err := json.Unmarshal(data, &state); err != nil {
		return true
	}
	if release.Version != "" && state.Version != "" && state.Version != release.Version {
		return false
	}
	if release.MD5 != "" && state.MD5 != "" && state.MD5 != release.MD5 {
		return false
	}
	return true
}

func writeExtractState(statePath, imagePath string, release Release) {
	state := extractState{Version: release.Version, MD5: release.MD5}
	if info, err := bootimage.Inspect(imagePath); err == nil {
		state.SHA256 = info.SHA256
	}
	data, err := json.Marshal(state)
	if err != nil {
		return
	}
	_ = os.WriteFile(statePath, data, 0o600)
}

func copyZipEntry(entry *zip.File, destination string) error {
	in, err := entry.Open()
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		os.Remove(destination)
		return copyErr
	}
	return closeErr
}
