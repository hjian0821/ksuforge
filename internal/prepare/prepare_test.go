package prepare

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hjian0821/ksuforge/internal/device"
)

type fakeDeviceRunner struct{}

func (fakeDeviceRunner) LookPath(name string) (string, error) {
	if name == "adb" {
		return name, nil
	}
	return "", os.ErrNotExist
}
func (fakeDeviceRunner) Run(_ context.Context, _ string, args ...string) (string, error) {
	if len(args) == 1 && args[0] == "devices" {
		return "List of devices attached\nphone\tdevice\n", nil
	}
	last := args[len(args)-1]
	props := map[string]string{"ro.product.manufacturer": "Xiaomi", "ro.product.brand": "Xiaomi", "ro.product.model": "Xiaomi 15", "ro.product.device": "dada", "ro.build.version.release": "17", "ro.build.version.incremental": "OS4.0.4.0.XOCCNXM", "ro.build.version.sdk": "37", "ro.boot.slot_suffix": "_b"}
	if v, ok := props[last]; ok {
		return v, nil
	}
	if last == "-r" {
		return "6.6.77-android15-8", nil
	}
	return "", nil
}

type noNetwork func(*http.Request) (*http.Response, error)

func (f noNetwork) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestVerifiedStockReuseWorksWithoutROMOrNetwork(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "stock", "dada", "OS4.0.4.0.XOCCNXM")
	os.MkdirAll(dir, 0o755)
	data := make([]byte, 4096)
	copy(data, "ANDROID!")
	sum := sha256.Sum256(data)
	image := filepath.Join(dir, "init_boot.img")
	os.WriteFile(image, data, 0o600)
	meta, _ := json.Marshal(map[string]string{"version": "OS4.0.4.0.XOCCNXM", "md5": "0123456789abcdef0123456789abcdef", "sha256": hex.EncodeToString(sum[:])})
	os.WriteFile(image+".json", meta, 0o600)
	calls := 0
	client := &http.Client{Transport: noNetwork(func(*http.Request) (*http.Response, error) {
		calls++
		return nil, errors.New("network must not be needed")
	})}
	service := Service{Devices: device.Client{Runner: fakeDeviceRunner{}}, HTTPClient: client}
	result, err := service.Run(context.Background(), Options{OutputDir: root, RequireKMI: true})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 0 || result.ImagePath != image || result.ROMPath != "" {
		t.Fatalf("result=%+v calls=%d", result, calls)
	}
	// An explicitly chosen local package still takes precedence over the cache.
	local := filepath.Join(root, "explicit.zip")
	os.WriteFile(local, []byte("invalid selected package"), 0o600)
	if _, err := service.Run(context.Background(), Options{OutputDir: root, ROMPath: local}); err == nil {
		t.Fatal("explicit local selection was silently bypassed")
	}
}
