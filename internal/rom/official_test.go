package rom

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type otaTransport func(*http.Request) (*http.Response, error)

func (f otaTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

const testVersion = "OS4.0.4.0.XOCCNXM"
const testFilename = "dada-ota_full-OS4.0.4.0.XOCCNXM-user-17.0-92227559a4.zip"
const testMD5 = "92227559a419839436075a1bf25d7784"

func testOfficialROM() otaROM {
	return otaROM{Type: "rom", Device: "dada", Version: testVersion, Filename: testFilename, MD5: testMD5, Size: "8.2G"}
}

func otaResponse(t *testing.T, current, latest otaROM) string {
	t.Helper()
	data, err := json.Marshal(map[string]otaROM{"CurrentRom": current, "LatestRom": latest})
	if err != nil {
		t.Fatal(err)
	}
	return encryptOTA(data)
}

func textResponse(body string, status int) *http.Response {
	return &http.Response{StatusCode: status, Status: http.StatusText(status), Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}
}

func TestFindRecoveryOfficialFirst(t *testing.T) {
	client := &http.Client{Transport: otaTransport(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != officialOTAURL || req.Method != http.MethodPost {
			t.Fatalf("unexpected request (fallback must not be used): %s %s", req.Method, req.URL)
		}
		if err := req.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if req.Form.Get("s") != "1" || req.Form.Get("t") != "" {
			t.Fatalf("unexpected anonymous protocol: %v", req.Form)
		}
		plain, err := decryptOTA(req.Form.Get("q"))
		if err != nil {
			t.Fatal(err)
		}
		var query map[string]string
		if err := json.Unmarshal(plain, &query); err != nil {
			t.Fatal(err)
		}
		for key, want := range map[string]string{"d": "dada", "c": "17", "ov": testVersion, "v": "miui-" + testVersion, "l": "zh_CN", "r": "CN", "f": "1"} {
			if query[key] != want {
				t.Errorf("query[%s] = %q, want %q", key, query[key], want)
			}
		}
		return textResponse(otaResponse(t, testOfficialROM(), otaROM{}), http.StatusOK), nil
	})}
	var log bytes.Buffer
	got, err := FindRecovery(context.Background(), client, "https://fallback.example/%s.xml", "dada", testVersion, "17", &log)
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != testVersion || got.MD5 != testMD5 || got.URL != "https://bigota.d.miui.com/"+testVersion+"/"+testFilename {
		t.Fatalf("unexpected official release: %+v", got)
	}
	if !strings.Contains(log.String(), "rom.official_found") {
		t.Fatalf("missing source in log: %s", &log)
	}
}

func TestFindRecoveryOfficialFallback(t *testing.T) {
	for _, scenario := range []string{"missing", "different version", "different device", "incremental", "bad md5", "bad path", "invalid response", "HTTP failure", "network failure"} {
		t.Run(scenario, func(t *testing.T) {
			candidate := testOfficialROM()
			switch scenario {
			case "missing":
				candidate = otaROM{}
			case "different version":
				candidate.Version = "OS4.0.3.0.XOCCNXM"
			case "different device":
				candidate.Device = "haotian"
			case "incremental":
				candidate.Filename = "dada-ota_incremental-OS4.0.3.0.XOCCNXM-" + testVersion + "-user.zip"
			case "bad md5":
				candidate.MD5 = "invalid"
			case "bad path":
				candidate.Filename = "../" + testFilename
			}
			var calls []string
			client := &http.Client{Transport: otaTransport(func(req *http.Request) (*http.Response, error) {
				calls = append(calls, req.URL.Host)
				if req.URL.String() == officialOTAURL {
					switch scenario {
					case "invalid response":
						return textResponse("<html>unavailable</html>", http.StatusOK), nil
					case "HTTP failure":
						return textResponse("", http.StatusServiceUnavailable), nil
					case "network failure":
						return nil, errors.New("unavailable")
					}
					return textResponse(otaResponse(t, candidate, otaROM{}), http.StatusOK), nil
				}
				if req.URL.String() != "https://fallback.example/dada.xml" {
					t.Fatalf("unexpected fallback URL: %s", req.URL)
				}
				return textResponse(`<rss><channel><item><title>MIUI `+testVersion+` Recovery update</title><link>https://bigota.d.miui.com/fallback.zip</link><description>&lt;b&gt;MD5:&lt;/b&gt; `+testMD5+`</description></item></channel></rss>`, http.StatusOK), nil
			})}
			var log bytes.Buffer
			got, err := FindRecovery(context.Background(), client, "https://fallback.example/%s.xml", "dada", testVersion, "17", &log)
			if err != nil {
				t.Fatal(err)
			}
			if got.URL != "https://bigota.d.miui.com/fallback.zip" || strings.Join(calls, ",") != "update.miui.com,fallback.example" || !strings.Contains(log.String(), "rom.catalog_fallback") {
				t.Fatalf("fallback failed: %+v, calls=%v, log=%s", got, calls, &log)
			}
		})
	}
}

func TestOfficialLatestMustMatchAndBeFullROM(t *testing.T) {
	for _, match := range []bool{false, true} {
		latest := testOfficialROM()
		latest.Filename += "?t=123&s=abc"
		if !match {
			latest.Version = "OS4.0.5.0.XOCCNXM"
		}
		client := &http.Client{Transport: otaTransport(func(*http.Request) (*http.Response, error) {
			return textResponse(otaResponse(t, otaROM{}, latest), http.StatusOK), nil
		})}
		got, err := findOfficialRecovery(context.Background(), client, "dada", testVersion, "17")
		if match && (err != nil || !strings.HasSuffix(got.URL, "?t=123&s=abc")) {
			t.Fatalf("signed exact-version full ROM rejected: %+v, %v", got, err)
		}
		if !match && err == nil {
			t.Fatal("substituted latest version")
		}
	}
}

func TestOfficialCancellationDoesNotFallback(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	client := &http.Client{Transport: otaTransport(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != officialOTAURL {
			t.Fatal("fallback requested after cancellation")
		}
		cancel()
		return nil, ctx.Err()
	})}
	_, err := FindRecovery(ctx, client, "https://fallback.example/%s.xml", "dada", testVersion, "17")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
}

func TestOTADeviceRegion(t *testing.T) {
	for suffix, want := range map[string]string{"CNXM": "dada", "MIXM": "dada_global", "EUXM": "dada_eea_global", "INXM": "dada_in_global", "RUXM": "dada_ru_global", "TWXM": "dada_tw_global"} {
		got, err := otaDevice("dada", "OS3.0.1.0.WOC"+suffix)
		if err != nil || got != want {
			t.Fatalf("%s: %s, %v", suffix, got, err)
		}
	}
	if got, err := otaDevice("dada_global", "OS3.0.1.0.WOCMIXM"); err != nil || got != "dada_global" {
		t.Fatalf("duplicated region suffix: %s, %v", got, err)
	}
}

func TestDecryptOTARejectsMalformedData(t *testing.T) {
	for _, input := range []string{"", "invalid!", "YWJj", strings.Repeat("A", 24)} {
		if _, err := decryptOTA(input); err == nil {
			t.Fatalf("accepted malformed ciphertext %q", input)
		}
	}
}
