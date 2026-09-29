// Command package-release creates a distributable archive from native build
// outputs. The archive name follows <project>_<version>_<goos>_<goarch>.
package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const projectName = "ksuforge"

var safeVersion = regexp.MustCompile(`^[0-9A-Za-z][0-9A-Za-z._-]*$`)

func main() {
	binDir := flag.String("bin", "bin", "directory containing native build outputs")
	version := flag.String("version", "dev", "release version, optionally prefixed with v")
	goos := flag.String("goos", runtime.GOOS, "target operating system")
	goarch := flag.String("goarch", runtime.GOARCH, "target architecture")
	flag.Parse()

	archive, err := packageRelease(*binDir, *version, *goos, *goarch)
	if err != nil {
		fmt.Fprintln(os.Stderr, "package release:", err)
		os.Exit(1)
	}
	fmt.Println(archive)
}

func packageRelease(binDir, version, goos, goarch string) (string, error) {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	if !safeVersion.MatchString(version) {
		return "", fmt.Errorf("invalid version %q", version)
	}
	if goos != "darwin" && goos != "linux" && goos != "windows" {
		return "", fmt.Errorf("unsupported operating system %q", goos)
	}
	if !safeVersion.MatchString(goarch) {
		return "", fmt.Errorf("invalid architecture %q", goarch)
	}

	rootName := strings.Join([]string{projectName, version, goos, goarch}, "_")
	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}
	destination := filepath.Join(binDir, rootName+ext)

	inputs := releaseInputs(goos)
	for _, input := range inputs {
		if _, err := os.Stat(filepath.Join(binDir, input)); err != nil {
			return "", fmt.Errorf("required build output %s: %w", input, err)
		}
	}

	temp, err := os.CreateTemp(binDir, ".package-release-*")
	if err != nil {
		return "", err
	}
	tempName := temp.Name()
	defer os.Remove(tempName)

	if goos == "windows" {
		err = writeZip(temp, binDir, rootName, inputs)
	} else {
		err = writeTarGz(temp, binDir, rootName, inputs)
	}
	if closeErr := temp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return "", err
	}
	if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if err := os.Rename(tempName, destination); err != nil {
		return "", err
	}
	return destination, nil
}

func releaseInputs(goos string) []string {
	exe := ""
	gui := projectName
	if goos == "windows" {
		exe = ".exe"
		gui += exe
	} else if goos == "darwin" {
		gui += ".app"
	}
	return []string{
		gui,
		projectName + "-cli" + exe,
		"ksud" + exe,
		"payload-dumper" + exe,
		"KernelSU-GPL-3.0.txt",
	}
}

func writeZip(out io.Writer, binDir, rootName string, inputs []string) error {
	zw := zip.NewWriter(out)
	for _, input := range inputs {
		if err := walkInput(binDir, input, func(path, archivePath string, info fs.FileInfo) error {
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = filepath.ToSlash(filepath.Join(rootName, archivePath))
			header.Modified = time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}
			entry, err := zw.CreateHeader(header)
			if err != nil || info.IsDir() {
				return err
			}
			return copyFile(entry, path)
		}); err != nil {
			_ = zw.Close()
			return err
		}
	}
	return zw.Close()
}

func writeTarGz(out io.Writer, binDir, rootName string, inputs []string) error {
	gz := gzip.NewWriter(out)
	tw := tar.NewWriter(gz)
	for _, input := range inputs {
		if err := walkInput(binDir, input, func(path, archivePath string, info fs.FileInfo) error {
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			header.Name = filepath.ToSlash(filepath.Join(rootName, archivePath))
			header.ModTime = time.Unix(0, 0)
			header.AccessTime = time.Time{}
			header.ChangeTime = time.Time{}
			if err := tw.WriteHeader(header); err != nil || info.IsDir() {
				return err
			}
			return copyFile(tw, path)
		}); err != nil {
			_ = tw.Close()
			_ = gz.Close()
			return err
		}
	}
	if err := tw.Close(); err != nil {
		_ = gz.Close()
		return err
	}
	return gz.Close()
}

func walkInput(binDir, input string, visit func(string, string, fs.FileInfo) error) error {
	root := filepath.Join(binDir, input)
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(binDir, path)
		if err != nil {
			return err
		}
		return visit(path, rel, info)
	})
}

func copyFile(dst io.Writer, source string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = io.Copy(dst, src)
	return err
}
