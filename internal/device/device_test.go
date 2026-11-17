package device

import (
	"context"
	"testing"
)

func TestParseGetVar(t *testing.T) {
	tests := []struct {
		out, key, want string
		ok             bool
	}{
		{"(bootloader) unlocked: yes\nFinished.", "unlocked", "yes", true},
		{"secure: no", "secure", "no", true},
		{"FAILED (remote: 'GetVar Variable Not found')", "secure", "", false},
	}
	for _, tt := range tests {
		got, ok := parseGetVar(tt.out, tt.key)
		if got != tt.want || ok != tt.ok {
			t.Fatalf("parseGetVar(%q): got %q/%v, want %q/%v", tt.out, got, ok, tt.want, tt.ok)
		}
	}
}

func TestParseKMI(t *testing.T) {
	tests := map[string]string{
		"5.10.198-android13-9-g123": "android13-5.10",
		"6.1.75-android14-11-g123":  "android14-6.1",
		"5.15.0-custom":             "",
	}
	for kernel, want := range tests {
		if got := ParseKMI(kernel); got != want {
			t.Fatalf("ParseKMI(%q)=%q, want %q", kernel, got, want)
		}
	}
}

type captureRunner struct {
	args []string
}

func (r *captureRunner) LookPath(name string) (string, error) { return name, nil }
func (r *captureRunner) Run(_ context.Context, name string, args ...string) (string, error) {
	r.args = append([]string{name}, args...)
	return "", nil
}

func TestRebootFastboot(t *testing.T) {
	var r captureRunner
	c := Client{Runner: &r}
	if err := c.RebootFastboot(context.Background(), "abc123"); err != nil {
		t.Fatal(err)
	}
	want := []string{"adb", "-s", "abc123", "reboot", "fastboot"}
	if len(r.args) != len(want) {
		t.Fatalf("args = %v, want %v", r.args, want)
	}
	for i := range want {
		if r.args[i] != want[i] {
			t.Fatalf("args = %v, want %v", r.args, want)
		}
	}
}

type getvarRunner struct{ out string }

func (r getvarRunner) LookPath(name string) (string, error) { return name, nil }
func (r getvarRunner) Run(_ context.Context, name string, args ...string) (string, error) {
	return r.out, nil
}

func TestIsUserspaceFastboot(t *testing.T) {
	c := Client{Runner: getvarRunner{out: "is-userspace: yes\nFinished."}}
	if !c.IsUserspaceFastboot(context.Background(), "abc") {
		t.Fatal("expected userspace fastboot")
	}
	c = Client{Runner: getvarRunner{out: "FAILED (remote: 'GetVar Variable Not found')"}}
	if c.IsUserspaceFastboot(context.Background(), "abc") {
		t.Fatal("expected non-userspace fastboot")
	}
}
