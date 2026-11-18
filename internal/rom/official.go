package rom

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hjian0821/ksuforge/internal/i18n"
	"github.com/hjian0821/ksuforge/internal/logging"
)

const officialOTAURL = "https://update.miui.com/updates/miotaV3.php"

// These public protocol constants are used by Xiaomi's anonymous OTA v1 API.
// They encode the wire format; HTTPS authenticates the server.
const otaKey = "miuiotavalided11"
const otaIV = "0102030405060708"

var otaFilename = regexp.MustCompile(`^[A-Za-z0-9_.-]+\.zip$`)

// FindRecovery queries Xiaomi first, then uses the configured RSS index as a
// fallback. Both sources must supply an exact-version full ROM and an MD5.
func FindRecovery(ctx context.Context, client *http.Client, catalogTemplate, codename, version, android string, output ...io.Writer) (Release, error) {
	if !validCodename(codename) {
		return Release{}, i18n.NewError("error.codename_invalid", codename)
	}
	version = strings.TrimSpace(version)
	if version == "" {
		return Release{}, i18n.NewError("error.version_missing")
	}
	out := io.Discard
	if len(output) > 0 && output[0] != nil {
		out = output[0]
	}
	logger := logging.WithWriter(out, logging.Language(ctx))
	logger.Info("rom.official_query", "codename", codename, "version", version)
	// Bound the entire query, including reading its body, so an unavailable
	// official service cannot prevent the fallback from being tried.
	queryCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	release, err := findOfficialRecovery(queryCtx, client, codename, version, android)
	cancel()
	if ctx.Err() != nil {
		return Release{}, ctx.Err()
	}
	if err == nil {
		logger.Info("rom.official_found", "version", release.Version, "md5", release.MD5, "url", release.URL)
		return release, nil
	}
	logger.Warn("rom.catalog_fallback", "reason", err)
	return findCatalogRecovery(ctx, client, catalogTemplate, codename, version)
}

type otaROM struct {
	Type     string `json:"type"`
	Device   string `json:"device"`
	Version  string `json:"version"`
	Filename string `json:"filename"`
	MD5      string `json:"md5"`
	Size     string `json:"filesize"`
}

func findOfficialRecovery(ctx context.Context, client *http.Client, codename, version, android string) (Release, error) {
	if strings.TrimSpace(android) == "" {
		return Release{}, i18n.NewError("error.official_android_missing")
	}
	device, err := otaDevice(codename, version)
	if err != nil {
		return Release{}, err
	}
	language, region := "en_US", "GL"
	if strings.HasSuffix(version, "CNXM") {
		language, region = "zh_CN", "CN"
	}
	payload, err := json.Marshal(map[string]string{
		"id": "", "c": android, "d": device, "f": "1", "ov": version,
		"l": language, "r": region, "v": "miui-" + strings.Replace(version, "OS1", "V816", 1),
	})
	if err != nil {
		return Release{}, err
	}
	form := url.Values{"q": {encryptOTA(payload)}, "t": {""}, "s": {"1"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, officialOTAURL, strings.NewReader(form.Encode()))
	if err != nil {
		return Release{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return Release{}, i18n.WrapError("error.official_query", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Release{}, i18n.NewError("error.catalog_http", resp.Status)
	}
	const limit = 2 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return Release{}, i18n.WrapError("error.official_query", err)
	}
	if len(body) > limit {
		return Release{}, i18n.NewError("error.official_response")
	}
	plain, err := decryptOTA(strings.TrimSpace(string(body)))
	if err != nil {
		return Release{}, i18n.WrapError("error.official_response", err)
	}
	var response struct {
		Current otaROM `json:"CurrentRom"`
		Latest  otaROM `json:"LatestRom"`
	}
	if err := json.Unmarshal(plain, &response); err != nil {
		return Release{}, i18n.WrapError("error.official_response", err)
	}
	for _, candidate := range []otaROM{response.Current, response.Latest} {
		if candidate.Version != version || candidate.Device != device || candidate.Type != "rom" {
			continue
		}
		// LatestRom can include signed query parameters. Keep them on the
		// download URL, while validating the filename separately. Never use
		// incremental, fastboot, arbitrary paths or another build's package.
		file, err := url.Parse(candidate.Filename)
		if err != nil || file.IsAbs() || file.Host != "" || file.Fragment != "" || !otaFilename.MatchString(file.Path) {
			continue
		}
		if !strings.Contains(file.Path, "-"+version+"-") && !strings.Contains(file.Path, "_"+version+"_") {
			continue
		}
		legacyPrefix := "miui_" + strings.ReplaceAll(device, "_", "") + "_"
		if !strings.HasPrefix(file.Path, device+"-ota_full-") && !strings.HasPrefix(strings.ToLower(file.Path), strings.ToLower(legacyPrefix)) {
			continue
		}
		if strings.Contains(file.Path, "ota_incremental") {
			continue
		}
		checksum, err := hex.DecodeString(candidate.MD5)
		if err != nil || len(checksum) != 16 {
			continue
		}
		downloadURL := "https://bigota.d.miui.com/" + version + "/" + file.String()
		if err := ValidateOfficialURL(downloadURL); err != nil {
			return Release{}, err
		}
		return Release{Version: version, URL: downloadURL, MD5: strings.ToLower(candidate.MD5), Size: candidate.Size}, nil
	}
	return Release{}, i18n.NewError("error.official_no_match", codename, version)
}

func otaDevice(codename, version string) (string, error) {
	regions := map[string]string{
		"CNXM": "", "MIXM": "_global", "EUXM": "_eea_global",
		"INXM": "_in_global", "IDXM": "_id_global", "RUXM": "_ru_global",
		"TWXM": "_tw_global", "TRXM": "_tr_global", "JPXM": "_jp_global",
	}
	for region, suffix := range regions {
		if strings.HasSuffix(version, region) {
			// Some devices already expose the OTA region in ro.product.device.
			if suffix != "" && !strings.HasSuffix(codename, suffix) {
				codename += suffix
			}
			return codename, nil
		}
	}
	return "", i18n.NewError("error.official_region", version)
}

func encryptOTA(plain []byte) string {
	block, _ := aes.NewCipher([]byte(otaKey))
	padding := aes.BlockSize - len(plain)%aes.BlockSize
	padded := append(append([]byte(nil), plain...), bytes.Repeat([]byte{byte(padding)}, padding)...)
	cipher.NewCBCEncrypter(block, []byte(otaIV)).CryptBlocks(padded, padded)
	return base64.URLEncoding.EncodeToString(padded)
}

func decryptOTA(encoded string) ([]byte, error) {
	// Xiaomi responses use standard Base64; requests use URL-safe Base64.
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		data, err = base64.URLEncoding.DecodeString(encoded)
	}
	if err != nil || len(data) == 0 || len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("invalid OTA ciphertext")
	}
	block, _ := aes.NewCipher([]byte(otaKey))
	cipher.NewCBCDecrypter(block, []byte(otaIV)).CryptBlocks(data, data)
	padding := int(data[len(data)-1])
	if padding == 0 || padding > aes.BlockSize || !bytes.Equal(data[len(data)-padding:], bytes.Repeat([]byte{byte(padding)}, padding)) {
		return nil, fmt.Errorf("invalid OTA padding")
	}
	return data[:len(data)-padding], nil
}
