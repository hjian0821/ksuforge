package rom

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"testing"
)

func TestDownloadSendsXiaomiReferer(t *testing.T) {
	payload := []byte("test zip payload")
	sum := md5.Sum(payload)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Referer(); got != "https://www.mi.com/" {
			t.Fatalf("Referer = %q", got)
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Range", "bytes 0-"+strconv.Itoa(len(payload)-1)+"/"+strconv.Itoa(len(payload)))
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write(payload)
	}))
	defer server.Close()

	// The production URL validator intentionally rejects test hosts, so replace
	// the server host with an allowed one and route requests back to the server.
	client := server.Client()
	transport := client.Transport.(*http.Transport).Clone()
	transport.TLSClientConfig = transport.TLSClientConfig.Clone()
	transport.TLSClientConfig.InsecureSkipVerify = true // Test server certificate cannot name the allowed production host.
	dialer := &net.Dialer{}
	transport.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, server.Listener.Addr().String())
	}
	client.Transport = transport

	dir := t.TempDir()
	path, err := Download(context.Background(), client, Release{
		URL: "https://bigota.d.miui.com/test.zip",
		MD5: hex.EncodeToString(sum[:]),
	}, dir)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(payload) {
		t.Fatalf("payload = %q", got)
	}
}

func TestDownloadResumesCompletedChunks(t *testing.T) {
	payload := bytes.Repeat([]byte{0x5a}, int(defaultChunkSize)+17)
	sum := md5.Sum(payload)
	var repeatedFirstChunk atomic.Bool
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var start, end int64
		if _, err := fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end); err != nil {
			http.Error(w, "bad range", http.StatusBadRequest)
			return
		}
		if start == 0 && end == defaultChunkSize-1 {
			repeatedFirstChunk.Store(true)
		}
		if end >= int64(len(payload)) {
			end = int64(len(payload)) - 1
		}
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(payload)))
		w.WriteHeader(http.StatusPartialContent)
		_, _ = w.Write(payload[start : end+1])
	}))
	defer server.Close()
	client := routeAllowedHostsToServer(server)

	dir := t.TempDir()
	destination := filepath.Join(dir, "resume.zip")
	part := destination + ".part"
	f, err := os.OpenFile(part, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(int64(len(payload))); err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteAt(payload[:defaultChunkSize], 0); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	if err := saveDownloadState(part+".json", downloadState{Version: 1, Size: int64(len(payload)), ChunkSize: defaultChunkSize, MD5: hex.EncodeToString(sum[:]), Done: []bool{true, false}}); err != nil {
		t.Fatal(err)
	}

	path, err := Download(context.Background(), client, Release{URL: "https://bigota.d.miui.com/resume.zip", MD5: hex.EncodeToString(sum[:])}, dir)
	if err != nil {
		t.Fatal(err)
	}
	if repeatedFirstChunk.Load() {
		t.Fatal("completed first chunk was downloaded again")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatal("resumed payload mismatch")
	}
}

func TestDownloadReusesLocalROMByVersion(t *testing.T) {
	payload := []byte("versioned ROM")
	sum := md5.Sum(payload)
	dir := t.TempDir()
	name := "houji_OS2.0.9.0.VNCCNXM_abc.zip"
	if err := os.WriteFile(filepath.Join(dir, name), payload, 0o600); err != nil {
		t.Fatal(err)
	}
	path, err := Download(context.Background(), http.DefaultClient, Release{
		URL:     "https://bigota.d.miui.com/other-name.zip",
		Version: "OS2.0.9.0.VNCCNXM",
		MD5:     hex.EncodeToString(sum[:]),
	}, dir)
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(dir, name) {
		t.Fatalf("path = %q", path)
	}
}

func routeAllowedHostsToServer(server *httptest.Server) *http.Client {
	client := server.Client()
	transport := client.Transport.(*http.Transport).Clone()
	transport.TLSClientConfig = transport.TLSClientConfig.Clone()
	transport.TLSClientConfig.InsecureSkipVerify = true // Test server certificate cannot name the allowed production host.
	dialer := &net.Dialer{}
	transport.DialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, server.Listener.Addr().String())
	}
	client.Transport = transport
	return client
}

func TestDownloadReusesMatchingExistingFile(t *testing.T) {
	payload := []byte("existing ROM")
	sum := md5.Sum(payload)
	dir := t.TempDir()
	name := "existing.zip"
	if err := os.WriteFile(filepath.Join(dir, name), payload, 0o600); err != nil {
		t.Fatal(err)
	}
	path, err := Download(context.Background(), http.DefaultClient, Release{
		URL: "https://bigota.d.miui.com/" + name,
		MD5: hex.EncodeToString(sum[:]),
	}, dir)
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(dir, name) {
		t.Fatalf("path = %q", path)
	}
}

func TestConfiguredMirrorURLKeepsOTAPath(t *testing.T) {
	original, err := url.Parse("https://bigota.d.miui.com/OS3.0.1/dada%20ota.zip")
	if err != nil {
		t.Fatal(err)
	}
	got, err := configuredMirrorURL("https://mirror.example.com/xiaomi", original)
	if err != nil {
		t.Fatal(err)
	}
	want := "https://mirror.example.com/xiaomi/OS3.0.1/dada%20ota.zip"
	if got != want {
		t.Fatalf("mirror URL = %q, want %q", got, want)
	}
}

func TestValidateMirrorBaseURL(t *testing.T) {
	for _, raw := range []string{"http://mirror.example.com", "https://user@mirror.example.com", "https://mirror.example.com?q=1"} {
		if err := ValidateMirrorBaseURL(raw); err == nil {
			t.Fatalf("expected %q to be rejected", raw)
		}
	}
	if err := ValidateMirrorBaseURL("https://mirror.example.com/xiaomi"); err != nil {
		t.Fatal(err)
	}
}

func TestDownloadRejectsInvalidConnectionsBeforeNetwork(t *testing.T) {
	_, err := DownloadWithOptions(context.Background(), http.DefaultClient, Release{
		URL: "https://bigota.d.miui.com/test.zip", MD5: "00",
	}, t.TempDir(), DownloadOptions{Connections: 33})
	if err == nil {
		t.Fatal("expected invalid connection count to fail")
	}
}

func TestVerifyFileMD5(t *testing.T) {
	payload := []byte("local recovery ROM")
	path := filepath.Join(t.TempDir(), "ota.zip")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	sum := md5.Sum(payload)
	if err := VerifyFileMD5(path, hex.EncodeToString(sum[:])); err != nil {
		t.Fatal(err)
	}
	if err := VerifyFileMD5(path, "00000000000000000000000000000000"); err == nil {
		t.Fatal("expected mismatched checksum to fail")
	}
}
