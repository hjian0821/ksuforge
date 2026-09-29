package rom

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindRecoveryExactVersion(t *testing.T) {
	feed := `<?xml version="1.0"?><rss><channel>
<item><title>MIUI OS1.0.2.0.UNCCNXM Fastboot update for Xiaomi 14 China</title><link>https://bigota.d.miui.com/fastboot.tgz</link><description></description></item>
<item><title>MIUI OS1.0.2.0.UNCCNXM Recovery update for Xiaomi 14 China</title><link>https://bigota.d.miui.com/rom.zip</link><description>&amp;lt;b&amp;gt;MD5:&amp;lt;/b&amp;gt; 0123456789abcdef0123456789abcdef</description></item>
</channel></rss>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, feed)
	}))
	defer server.Close()
	release, err := FindRecovery(context.Background(), server.Client(), server.URL+"/%s.xml", "houji", "OS1.0.2.0.UNCCNXM")
	if err != nil {
		t.Fatal(err)
	}
	if release.URL != "https://bigota.d.miui.com/rom.zip" || release.MD5 != "0123456789abcdef0123456789abcdef" {
		t.Fatalf("unexpected release: %#v", release)
	}
}
func TestFindRecoveryDoesNotSubstituteVersion(t *testing.T) {
	feed := `<?xml version="1.0"?><rss><channel><item><title>MIUI OS1.0.3.0.UNCCNXM Recovery update</title><link>https://bigota.d.miui.com/rom.zip</link></item></channel></rss>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, feed) }))
	defer server.Close()
	_, err := FindRecovery(context.Background(), server.Client(), server.URL+"/%s.xml", "houji", "OS1.0.2.0.UNCCNXM")
	if err == nil {
		t.Fatal("expected exact-version lookup to fail")
	}
}

func TestValidateOfficialURL(t *testing.T) {
	if err := ValidateOfficialURL("https://bigota.d.miui.com/version/rom.zip"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateOfficialURL("https://example.com/rom.zip"); err == nil {
		t.Fatal("expected untrusted host to fail")
	}
}

func TestValidateCatalogURL(t *testing.T) {
	if err := ValidateCatalogURL("https://example.com/rss/%s.xml"); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{
		"",
		"http://example.com/%s.xml",
		"https://example.com/rss.xml",
		"https://example.com/%s/%s.xml",
		"https://example.com/%d.xml",
		"https://user@example.com/%s.xml",
	} {
		if err := ValidateCatalogURL(bad); err == nil {
			t.Fatalf("expected %q to fail", bad)
		}
	}
}
