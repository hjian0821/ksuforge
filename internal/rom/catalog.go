package rom

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/hjian0821/ksuforge/internal/i18n"
)

const DefaultCatalogURL = "https://raw.githubusercontent.com/XiaomiFirmwareUpdater/miui-updates-tracker/master/rss/%s.xml"

var md5Pattern = regexp.MustCompile(`(?i)MD5:</b>\s*([0-9a-f]{32})`)
var sizePattern = regexp.MustCompile(`(?i)Size:</b>\s*([^<]+)`)

type Release struct {
	Version string
	URL     string
	MD5     string
	Size    string
}

type rss struct {
	Channel struct {
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func FindRecovery(ctx context.Context, client *http.Client, catalogTemplate, codename, version string) (Release, error) {
	if !validCodename(codename) {
		return Release{}, i18n.NewError("error.codename_invalid", codename)
	}
	if strings.TrimSpace(version) == "" {
		return Release{}, i18n.NewError("error.version_missing")
	}
	if catalogTemplate == "" {
		catalogTemplate = DefaultCatalogURL
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(catalogTemplate, codename), nil)
	if err != nil {
		return Release{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return Release{}, i18n.WrapError("error.catalog_read", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Release{}, i18n.NewError("error.catalog_http", resp.Status)
	}
	var feed rss
	if err := xml.NewDecoder(io.LimitReader(resp.Body, 32<<20)).Decode(&feed); err != nil {
		return Release{}, i18n.WrapError("error.catalog_parse", err)
	}
	version = strings.TrimSpace(version)
	for _, item := range feed.Channel.Items {
		fields := strings.Fields(item.Title)
		if len(fields) < 4 || fields[1] != version || !strings.EqualFold(fields[2], "Recovery") {
			continue
		}
		if err := ValidateOfficialURL(item.Link); err != nil {
			return Release{}, err
		}
		description := html.UnescapeString(item.Description)
		match := md5Pattern.FindStringSubmatch(description)
		if len(match) != 2 {
			return Release{}, i18n.NewError("error.catalog_md5_missing")
		}
		size := ""
		if sizeMatch := sizePattern.FindStringSubmatch(description); len(sizeMatch) == 2 {
			size = strings.TrimSpace(sizeMatch[1])
		}
		return Release{Version: version, URL: item.Link, MD5: strings.ToLower(match[1]), Size: size}, nil
	}
	return Release{}, i18n.NewError("error.catalog_no_match", codename, version)
}

// ValidateCatalogURL checks a user-supplied ROM index template. It must be an
// HTTPS URL with exactly one %s placeholder for the device codename. The links
// the index returns are still validated against the Xiaomi CDN, so a custom
// index host cannot redirect downloads to an untrusted source.
func ValidateCatalogURL(template string) error {
	template = strings.TrimSpace(template)
	if template == "" {
		return i18n.NewError("error.catalog_url_required")
	}
	if strings.Count(template, "%s") != 1 {
		return i18n.NewError("error.catalog_placeholder")
	}
	probe := strings.ReplaceAll(template, "%s", "codename")
	u, err := url.Parse(probe)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
		return i18n.NewError("error.catalog_https")
	}
	return nil
}

func ValidateOfficialURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil {
		return i18n.NewError("error.rom_url_untrusted")
	}
	host := strings.ToLower(u.Hostname())
	trusted := host == "bigota.d.miui.com" || host == "hugeota.d.miui.com" ||
		strings.HasSuffix(host, ".miui.com") || strings.HasSuffix(host, ".mi-img.com") ||
		host == "bkt-sgp-miui-ota-update-alisgp.oss-ap-southeast-1.aliyuncs.com"
	if !trusted {
		return i18n.NewError("error.rom_cdn_untrusted", host)
	}
	return nil
}

func validCodename(v string) bool {
	if v == "" {
		return false
	}
	for _, r := range v {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '_' && r != '-' {
			return false
		}
	}
	return true
}
