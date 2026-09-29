//go:build windows

package command

import (
	"context"
	"os/exec"
	"syscall"
)

// createNoWindow keeps console child processes (adb, fastboot, ksud, ...) from
// flashing a console window when launched by the windowed GUI application.
const createNoWindow = 0x08000000

// CommandContext is exec.CommandContext with a hidden console on Windows.
func CommandContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: createNoWindow,
	}
	return cmd
}
