package device

import "testing"

func TestSelectADB(t *testing.T) {
	devices := []Device{{Serial: "one", Mode: ADB}, {Serial: "boot", Mode: Fastboot}}
	got, err := SelectADB(devices, "")
	if err != nil || got.Serial != "one" {
		t.Fatalf("got %#v, %v", got, err)
	}
}

func TestSelectADBRequiresSerialForMultiple(t *testing.T) {
	devices := []Device{{Serial: "one", Mode: ADB}, {Serial: "two", Mode: ADB}}
	if _, err := SelectADB(devices, ""); err == nil {
		t.Fatal("expected multiple-device error")
	}
}
