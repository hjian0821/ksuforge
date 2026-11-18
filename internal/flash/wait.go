package flash

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hjian0821/ksuforge/internal/i18n"
	"github.com/hjian0821/ksuforge/internal/logging"
)

const fastbootWaitTimeout = 45 * time.Second

// Never substitute another attached phone when the requested serial disappears.
// Both the overall wait and each tool invocation have deadlines.
func (s Service) waitFastboot(ctx context.Context, serial string, timeout time.Duration) error {
	logger := logging.WithWriter(s.Out, logging.Language(ctx))
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	tool, _ := s.Devices.Runner.LookPath("fastboot")
	last, previous := "not detected", ""
	for {
		probeCtx, stop := context.WithTimeout(waitCtx, 3*time.Second)
		out, err := s.Devices.Runner.Run(probeCtx, "fastboot", "devices", "-l")
		stop()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err == nil {
			for _, line := range strings.Split(out, "\n") {
				fields := strings.Fields(line)
				if len(fields) >= 2 && fields[0] == serial && fields[1] == "fastboot" {
					logger.Info("flash.fastboot_connected", "serial", serial, "tool", tool)
					return nil
				}
			}
			last = "no devices"
			if strings.TrimSpace(out) != "" {
				last = diagnosticText(out)
			}
		} else {
			last = diagnosticText(err.Error())
		}
		if last != previous {
			logger.Info("flash.fastboot_wait", "serial", serial, "fastboot", last, "tool", tool)
			previous = last
		}
		select {
		case <-waitCtx.Done():
			if ctx.Err() != nil {
				return ctx.Err()
			}
			// A fresh bounded read explains whether Android has returned after reboot.
			probeCtx, stop := context.WithTimeout(ctx, 3*time.Second)
			adb, err := s.Devices.Runner.Run(probeCtx, "adb", "devices", "-l")
			stop()
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if err != nil {
				adb = err.Error()
			}
			adb = diagnosticText(adb)
			logger.Warn("flash.fastboot_wait_failed", "serial", serial, "fastboot", last, "adb", adb, "tool", tool)
			return i18n.WrapError("error.fastboot_timeout", fmt.Errorf("target=%s; fastboot=%s; adb=%s; tool=%s", serial, last, adb, tool))
		case <-time.After(time.Second):
		}
	}
}
func diagnosticText(text string) string {
	text = strings.Join(strings.Fields(text), " ")
	if text == "" {
		return "no devices"
	}
	runes := []rune(text)
	if len(runes) > 512 {
		return string(runes[:512]) + "…"
	}
	return text
}
