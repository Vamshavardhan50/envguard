// internal/parser/config_test.go
// Tests for envguard config parsing.

package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfigMissingFileReturnsDefaults(t *testing.T) {
	cfg, err := ParseConfig(filepath.Join(t.TempDir(), "missing.yaml"))
	require.NoError(t, err)
	assert.Equal(t, DefaultConfig(), cfg)
}

func TestParseConfigEmptyPath(t *testing.T) {
	cfg, err := ParseConfig("")
	require.Error(t, err)
	assert.Equal(t, DefaultConfig(), cfg)
}

func TestParseConfigValidFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), ".envguard.yaml")
	content := []byte(`version: 1
scan:
  paths:
    - "src"
  ignore:
    - "node_modules"
  languages:
    - "auto"
rules:
  PORT:
    required: false
    type: number
    default: "3000"
    description: "Server port"
output:
  format: json
  exit_code: false
`)

	require.NoError(t, os.WriteFile(filePath, content, 0o600))

	cfg, err := ParseConfig(filePath)
	require.NoError(t, err)
	assert.Equal(t, 1, cfg.Version)
	assert.Equal(t, []string{"src"}, cfg.Scan.Paths)
	assert.Equal(t, []string{"node_modules"}, cfg.Scan.Ignore)
	assert.Equal(t, []string{"auto"}, cfg.Scan.Languages)
	assert.Equal(t, "json", cfg.Output.Format)
	assert.False(t, cfg.Output.ExitCode)
	assert.Contains(t, cfg.Rules, "PORT")
	assert.Equal(t, "number", cfg.Rules["PORT"].Type)
}

func TestParseConfigDefaultsApplied(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), ".envguard.yaml")
	content := []byte(`version: 1
output:
  format: pretty
`)
	require.NoError(t, os.WriteFile(filePath, content, 0o600))

	cfg, err := ParseConfig(filePath)
	require.NoError(t, err)
	assert.Equal(t, DefaultConfig().Scan.Paths, cfg.Scan.Paths)
	assert.Equal(t, DefaultConfig().Scan.Ignore, cfg.Scan.Ignore)
	assert.Equal(t, DefaultConfig().Scan.Languages, cfg.Scan.Languages)
	assert.Equal(t, true, cfg.Output.ExitCode)
}
