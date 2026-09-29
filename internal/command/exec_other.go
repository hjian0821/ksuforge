//go:build !windows

package command

import (
	"context"
	"os/exec"
)

// CommandContext is exec.CommandContext. On non-Windows platforms there is no
// console window to hide.
func CommandContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, name, args...)
}
