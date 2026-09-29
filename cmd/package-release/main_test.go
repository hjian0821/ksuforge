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

func TestPackageReleasesNamesAndContents(t *testing.T) {
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
			for _, spec := range releasePackages(tt.goos) {
				for _, input := range spec.inputs {
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
			}
			for _, helper := range []string{"ksud", "payload-dumper"} {
				if tt.goos == "windows" {
					helper += ".exe"
				}
				if err := os.WriteFile(filepath.Join(binDir, helper), []byte(helper), 0o755); err != nil {
					t.Fatal(err)
				}
			}

			archives, err := packageReleases(binDir, "v1.2.3", tt.goos, "amd64")
			if err != nil {
				t.Fatal(err)
			}
			if len(archives) != 2 {
				t.Fatalf("archive count = %d, want 2", len(archives))
			}

			guiRoot := "ksuforge-gui_1.2.3_" + tt.goos + "_amd64"
			cliRoot := "ksuforge-cli_1.2.3_" + tt.goos + "_amd64"
			cliName := "ksuforge-cli"
			if tt.goos == "windows" {
				cliName += ".exe"
			}
			assertArchive(t, tt.listFiles, archives[0], guiRoot+tt.ext,
				[]string{filepath.ToSlash(filepath.Join(guiRoot, tt.gui)), filepath.ToSlash(filepath.Join(guiRoot, "KernelSU-GPL-3.0.txt"))},
				[]string{cliName, "ksud", "payload-dumper"})
			assertArchive(t, tt.listFiles, archives[1], cliRoot+tt.ext,
				[]string{filepath.ToSlash(filepath.Join(cliRoot, cliName)), filepath.ToSlash(filepath.Join(cliRoot, "KernelSU-GPL-3.0.txt"))},
				[]string{tt.gui, "ksud", "payload-dumper"})
		})
	}
}

func TestPackageReleasesRejectsInvalidVersion(t *testing.T) {
	if _, err := packageReleases(t.TempDir(), "../release", "linux", "amd64"); err == nil {
		t.Fatal("expected invalid version to fail")
	}
}

func assertArchive(t *testing.T, listFiles func(*testing.T, string) []string, archive, wantName string, required, forbidden []string) {
	t.Helper()
	if filepath.Base(archive) != wantName {
		t.Fatalf("archive name = %q, want %q", filepath.Base(archive), wantName)
	}
	files := listFiles(t, archive)
	for _, want := range required {
		if !contains(files, want) && !contains(files, want+"/") {
			t.Errorf("archive %q does not contain %q: %v", wantName, want, files)
		}
	}
	for _, unwanted := range forbidden {
		root := strings.Split(required[0], "/")[0]
		if contains(files, filepath.ToSlash(filepath.Join(root, unwanted))) || contains(files, filepath.ToSlash(filepath.Join(root, unwanted+".exe"))) {
			t.Errorf("archive %q unexpectedly contains top-level %q: %v", wantName, unwanted, files)
		}
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
