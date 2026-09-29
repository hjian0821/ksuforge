package command

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hjian0821/ksuforge/internal/i18n"
)

type Runner interface {
	Run(ctx context.Context, name string, args ...string) (string, error)
	LookPath(name string) (string, error)
}

type ExecRunner struct{ SearchPaths []string }

func (r ExecRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	if !strings.ContainsRune(name, filepath.Separator) {
		if resolved, err := r.LookPath(name); err == nil {
			name = resolved
		}
	}
	cmd := CommandContext(ctx, name, args...)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	err := cmd.Run()
	text := strings.TrimSpace(output.String())
	if err != nil {
		if text == "" {
			return "", fmt.Errorf("%s: %w", name, err)
		}
		return text, fmt.Errorf("%s: %w: %s", name, err, text)
	}
	return text, nil
}

func (r ExecRunner) LookPath(name string) (string, error) {
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}
	for _, dir := range r.SearchPaths {
		candidate := filepath.Join(dir, name)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() && info.Mode()&0o111 != 0 {
			return candidate, nil
		}
		if info, err := os.Stat(candidate + ".exe"); err == nil && !info.IsDir() {
			return candidate + ".exe", nil
		}
	}
	return "", i18n.NewError("error.command_not_found", name)
}

func DesktopSearchPaths() []string {
	home, _ := os.UserHomeDir()
	paths := []string{}
	if executable, err := os.Executable(); err == nil {
		dir := filepath.Dir(executable)
		paths = append(paths, dir, filepath.Join(dir, "tools"), filepath.Clean(filepath.Join(dir, "..", "Resources", "tools")))
	}
	paths = append(paths, "/opt/homebrew/bin", "/usr/local/bin")
	if home != "" {
		paths = append(paths, filepath.Join(home, "Library", "Android", "sdk", "platform-tools"), filepath.Join(home, "Android", "Sdk", "platform-tools"), filepath.Join(home, "go", "bin"))
	}
	if sdk := os.Getenv("ANDROID_SDK_ROOT"); sdk != "" {
		paths = append(paths, filepath.Join(sdk, "platform-tools"))
	}
	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		paths = append(paths, filepath.Join(local, "Android", "Sdk", "platform-tools"))
	}
	return paths
}
