package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunLocalizesUnknownCommand(t *testing.T) {
	for _, test := range []struct {
		language string
		want     string
	}{
		{language: "zh", want: "未知命令"},
		{language: "en", want: "Unknown command"},
	} {
		t.Run(test.language, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run(context.Background(), []string{"--lang", test.language, "does-not-exist"}, &stdout, &stderr)
			if code != 1 {
				t.Fatalf("exit code = %d", code)
			}
			if !strings.Contains(stderr.String(), test.want) {
				t.Fatalf("stderr %q does not contain %q", stderr.String(), test.want)
			}
		})
	}
}

func TestParseLanguageEqualsSyntax(t *testing.T) {
	args, language, err := parseLanguage([]string{"--lang=en", "devices"}, "zh")
	if err != nil || language != "en" || len(args) != 1 || args[0] != "devices" {
		t.Fatalf("args=%v language=%q err=%v", args, language, err)
	}
}
