package flash

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/hjian0821/ksuforge/internal/device"
)

type waitRunner struct {
	run func(context.Context, string, ...string) (string, error)
}

func (r waitRunner) LookPath(name string) (string, error) { return "/tools/" + name, nil }
func (r waitRunner) Run(ctx context.Context, name string, args ...string) (string, error) {
	return r.run(ctx, name, args...)
}

func TestWaitFastbootMatchesOnlyRequestedSerial(t *testing.T) {
	for _, tt := range []struct {
		name, out string
		connected bool
	}{
		{"target", "phone\tfastboot usb:1-2", true},
		{"another phone", "other\tfastboot", false},
		{"warning is not transport", "phone permission denied", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var log bytes.Buffer
			r := waitRunner{run: func(_ context.Context, name string, args ...string) (string, error) {
				if name == "adb" {
					return "List of devices attached\nphone device usb:1-2", nil
				}
				return tt.out, nil
			}}
			err := (Service{Devices: device.Client{Runner: r}, Out: &log}).waitFastboot(context.Background(), "phone", 10*time.Millisecond)
			if (err == nil) != tt.connected {
				t.Fatalf("connected=%v, err=%v", tt.connected, err)
			}
			if !tt.connected && (!strings.Contains(err.Error(), "phone device") || !strings.Contains(log.String(), "flash.fastboot_wait_failed")) {
				t.Fatalf("missing diagnostics: %v; %s", err, log.String())
			}
		})
	}
}
func TestWaitFastbootPreservesToolError(t *testing.T) {
	r := waitRunner{run: func(_ context.Context, name string, args ...string) (string, error) {
		if name == "fastboot" {
			return "", errors.New("USB access denied")
		}
		return "List of devices attached", nil
	}}
	err := (Service{Devices: device.Client{Runner: r}}).waitFastboot(context.Background(), "phone", 10*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "USB access denied") || !strings.Contains(err.Error(), "/tools/fastboot") {
		t.Fatalf("got %v", err)
	}
}
func TestWaitFastbootBoundsHungProbe(t *testing.T) {
	r := waitRunner{run: func(ctx context.Context, name string, args ...string) (string, error) {
		if name == "adb" {
			return "", nil
		}
		<-ctx.Done()
		return "", ctx.Err()
	}}
	started := time.Now()
	err := (Service{Devices: device.Client{Runner: r}}).waitFastboot(context.Background(), "phone", 10*time.Millisecond)
	if err == nil || time.Since(started) > time.Second {
		t.Fatalf("unbounded wait: %v", err)
	}
}
func TestWaitFastbootCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := waitRunner{run: func(ctx context.Context, name string, args ...string) (string, error) { return "", ctx.Err() }}
	err := (Service{Devices: device.Client{Runner: r}}).waitFastboot(ctx, "phone", time.Second)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("got %v", err)
	}
}
