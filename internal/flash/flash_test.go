package flash

import "testing"

func TestChoosePartition(t *testing.T) {
	tests := []struct {
		mode string
		sdk  int
		want string
	}{
		{"lkm", 32, "boot"}, {"lkm", 33, "init_boot"}, {"gki", 35, "boot"},
	}
	for _, tt := range tests {
		got, err := ChoosePartition("auto", tt.mode, tt.sdk)
		if err != nil || got != tt.want {
			t.Fatalf("mode=%s sdk=%d: got %q, %v", tt.mode, tt.sdk, got, err)
		}
	}
}

func TestExplicitPartitionWins(t *testing.T) {
	got, err := ChoosePartition("boot", "lkm", 35)
	if err != nil || got != "boot" {
		t.Fatalf("got %q, %v", got, err)
	}
}

func TestChooseDevicePartitionMatchesKernelSU(t *testing.T) {
	got, err := ChooseDevicePartition("auto", "lkm", "android14-6.1", true, 34)
	if err != nil || got != "init_boot" {
		t.Fatalf("got %q, %v", got, err)
	}
	got, err = ChooseDevicePartition("auto", "lkm", "android12-5.10", true, 32)
	if err != nil || got != "boot" {
		t.Fatalf("got %q, %v", got, err)
	}
}
