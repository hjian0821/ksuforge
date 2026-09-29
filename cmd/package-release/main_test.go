package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPackageReleaseNamesAndContents(t *testing.T) {
	tests := []struct {
		goos      string
		ext       string
		gui       string
		listFiles func(*testing.T, string) []string
	}{
		{goos: "linux", ext: ".tar.gz", gui: "ksuforge", listFiles: listTarGz},
		{goos: "windows", ext: ".zip", gui: "ksuforge.exe", listFiles: listZip},
		{goos: "darwin", ext: ".tar.gz", gui: "ksuforge.app", listFiles: listTarGz},
	}

	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			binDir := t.TempDir()
			for _, input := range releaseInputs(tt.goos) {
				path := filepath.Join(binDir, input)
				if strings.HasSuffix(input, ".app") {
					path = filepath.Join(path, "Contents", "MacOS", "ksuforge")
				}
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(input), 0o755); err != nil {
					t.Fatal(err)
				}
			}

			archive, err := packageRelease(binDir, "v1.2.3", tt.goos, "amd64")
			if err != nil {
				t.Fatal(err)
			}
			wantName := "ksuforge_1.2.3_" + tt.goos + "_amd64" + tt.ext
			if filepath.Base(archive) != wantName {
				t.Fatalf("archive name = %q, want %q", filepath.Base(archive), wantName)
			}
			files := tt.listFiles(t, archive)
			wantGUI := filepath.ToSlash(filepath.Join("ksuforge_1.2.3_"+tt.goos+"_amd64", tt.gui))
			if !contains(files, wantGUI) && !contains(files, wantGUI+"/") {
				t.Fatalf("archive does not contain %q: %v", wantGUI, files)
			}
		})
	}
}

func TestPackageReleaseRejectsInvalidVersion(t *testing.T) {
	if _, err := packageRelease(t.TempDir(), "../release", "linux", "amd64"); err == nil {
		t.Fatal("expected invalid version to fail")
	}
}

func listZip(t *testing.T, path string) []string {
	t.Helper()
	zr, err := zip.OpenReader(path)
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	files := make([]string, 0, len(zr.File))
	for _, file := range zr.File {
		files = append(files, file.Name)
	}
	return files
}

func listTarGz(t *testing.T, path string) []string {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	var files []string
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return files
		}
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, header.Name)
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
