package device

import "testing"

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
