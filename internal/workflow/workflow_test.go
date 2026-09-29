package workflow

import (
	"runtime"
	"testing"
)

func TestQuoteCommandArg(t *testing.T) {
	got := quoteCommandArg("/tmp/My ROM/image.img")
	if runtime.GOOS == "windows" {
		if got != `"/tmp/My ROM/image.img"` {
			t.Fatalf("quoteCommandArg() = %q", got)
		}
		return
	}
	if got != `'/tmp/My ROM/image.img'` {
		t.Fatalf("quoteCommandArg() = %q", got)
	}
}
