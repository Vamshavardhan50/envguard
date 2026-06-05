// internal/parser/config.go
// Loads and parses envguard configuration files.

package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config defines the envguard configuration schema.
type Config struct {
	Version int                  `mapstructure:"version"`
	Scan    ScanConfig           `mapstructure:"scan"`
	Rules   map[string]RuleConfig `mapstructure:"rules"`
	Output  OutputConfig         `mapstructure:"output"`
}

// CustomPattern defines a user-specified regex pattern for a language.
type CustomPattern struct {
	Language string `mapstructure:"language"`
	Regex    string `mapstructure:"regex"`
}

// ScanConfig defines file scanning configuration.
type ScanConfig struct {
	Paths          []string        `mapstructure:"paths"`
	Ignore         []string        `mapstructure:"ignore"`
	Languages      []string        `mapstructure:"languages"`
	CustomPatterns []CustomPattern `mapstructure:"patterns"`
}

// RuleConfig defines validation rules for a single env var.
type RuleConfig struct {
	Required    bool     `mapstructure:"required"`
	Type        string   `mapstructure:"type"`
	Values      []string `mapstructure:"values"`
	Default     string   `mapstructure:"default"`
	Description string   `mapstructure:"description"`
}

// OutputConfig defines output formatting settings.
type OutputConfig struct {
	Format   string `mapstructure:"format"`
	ExitCode bool   `mapstructure:"exit_code"`
}

// DefaultConfig returns the baseline configuration used when no config file is present.
func DefaultConfig() Config {
	return Config{
		Version: 1,
		Scan: ScanConfig{
			Paths: []string{"."},
			Ignore: []string{
				"node_modules",
				".git",
				"dist",
				"build",
				"vendor",
				".goreleaser.yaml",
			},
			Languages: []string{"auto"},
		},
		Rules: map[string]RuleConfig{},
		Output: OutputConfig{
			Format:   "pretty",
			ExitCode: true,
		},
	}
}

// ParseConfig reads and parses the envguard config file.
// If the file does not exist, defaults are returned without error.
func ParseConfig(path string) (Config, error) {
	defaults := DefaultConfig()
	if path == "" {
		return defaults, fmt.Errorf("parse config: path is required")
	}

	cleanPath := filepath.Clean(path)
	if _, err := os.Stat(cleanPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaults, nil
		}
		return defaults, fmt.Errorf("stat config file: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(cleanPath)
	v.SetConfigType("yaml")
	v.SetDefault("version", defaults.Version)
	v.SetDefault("scan.paths", defaults.Scan.Paths)
	v.SetDefault("scan.ignore", defaults.Scan.Ignore)
	v.SetDefault("scan.languages", defaults.Scan.Languages)
	v.SetDefault("output.format", defaults.Output.Format)
	v.SetDefault("output.exit_code", defaults.Output.ExitCode)

	if err := v.ReadInConfig(); err != nil {
		return defaults, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return defaults, fmt.Errorf("parse config file: %w", err)
	}

	applyDefaults(v, defaults, &cfg)
	return cfg, nil
}

func applyDefaults(v *viper.Viper, defaults Config, cfg *Config) {
	if !v.IsSet("version") {
		cfg.Version = defaults.Version
	}
	if !v.IsSet("scan.paths") || len(cfg.Scan.Paths) == 0 {
		cfg.Scan.Paths = append([]string{}, defaults.Scan.Paths...)
	}
	if !v.IsSet("scan.ignore") || len(cfg.Scan.Ignore) == 0 {
		cfg.Scan.Ignore = append([]string{}, defaults.Scan.Ignore...)
	}
	if !v.IsSet("scan.languages") || len(cfg.Scan.Languages) == 0 {
		cfg.Scan.Languages = append([]string{}, defaults.Scan.Languages...)
	}
	if cfg.Rules == nil {
		cfg.Rules = map[string]RuleConfig{}
	} else {
		upperRules := make(map[string]RuleConfig, len(cfg.Rules))
		for k, val := range cfg.Rules {
			upperRules[strings.ToUpper(k)] = val
		}
		cfg.Rules = upperRules
	}
	if !v.IsSet("output.format") {
		cfg.Output.Format = defaults.Output.Format
	}
	if !v.IsSet("output.exit_code") {
		cfg.Output.ExitCode = defaults.Output.ExitCode
	}
}
