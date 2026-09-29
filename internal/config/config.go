package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/hjian0821/ksuforge/internal/paths"
)

// Config is the small amount of state KSUForge remembers between runs. Large
// artifacts never live here; only user preferences.
type Config struct {
	OutputDir     string `json:"output_dir,omitempty"`
	Language      string `json:"language,omitempty"`
	Mode          string `json:"mode,omitempty"`
	Partition     string `json:"partition,omitempty"`
	AutoFlash     bool   `json:"auto_flash,omitempty"`
	NoReboot      bool   `json:"no_reboot,omitempty"`
	CatalogURL    string `json:"catalog_url,omitempty"`
	Mirrors       string `json:"mirrors,omitempty"`
	Connections   int    `json:"connections,omitempty"`
	KSUd          string `json:"ksud,omitempty"`
	LogAutoScroll *bool  `json:"log_auto_scroll,omitempty"`
}

func (c Config) ConnectionsOrDefault() int {
	if c.Connections >= 1 && c.Connections <= 32 {
		return c.Connections
	}
	return 16
}

// Path returns the config file location under the OS config directory.
func Path() string {
	return filepath.Join(paths.ConfigDir(), "config.json")
}

// Load reads the config file, returning an empty config when it is missing or
// unreadable. A corrupt file never blocks startup.
func Load() Config {
	return LoadFrom(Path())
}

// LoadFrom reads a config from an explicit path.
func LoadFrom(file string) Config {
	data, err := os.ReadFile(file)
	if err != nil {
		return Config{}
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}
	}
	return cfg
}

// Save writes the config atomically to the OS config directory.
func (c Config) Save() error {
	return c.SaveTo(Path())
}

// SaveTo writes the config atomically to an explicit path.
func (c Config) SaveTo(file string) error {
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	tmp := file + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, file)
}

// OutputDirOrDefault returns the remembered artifact directory or the platform
// default when none has been chosen yet.
func (c Config) OutputDirOrDefault() string {
	if c.OutputDir != "" {
		return c.OutputDir
	}
	return paths.DefaultOutputDir()
}

// LanguageOrDefault returns the remembered UI language when valid.
func (c Config) LanguageOrDefault(fallback string) string {
	if c.Language == "zh" || c.Language == "en" {
		return c.Language
	}
	return fallback
}

// ModeOrDefault returns the remembered KernelSU mode, defaulting to lkm.
func (c Config) ModeOrDefault() string {
	if c.Mode == "lkm" || c.Mode == "gki" {
		return c.Mode
	}
	return "lkm"
}

// PartitionOrDefault returns the remembered target partition, defaulting to auto.
func (c Config) PartitionOrDefault() string {
	switch c.Partition {
	case "auto", "boot", "init_boot":
		return c.Partition
	default:
		return "auto"
	}
}

// LogAutoScrollDefault reports whether the log should auto-scroll. It defaults
// to true so first-run behavior is unchanged.
func (c Config) LogAutoScrollDefault() bool {
	if c.LogAutoScroll == nil {
		return true
	}
	return *c.LogAutoScroll
}

// CatalogURLOr returns the remembered ROM index template, or def when unset.
func (c Config) CatalogURLOr(def string) string {
	if strings.TrimSpace(c.CatalogURL) != "" {
		return c.CatalogURL
	}
	return def
}
