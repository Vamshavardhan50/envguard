// internal/scanner/scanner_test.go
// Tests for the scanner implementation.

package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{name: "js", path: "src/index.js", expected: "javascript"},
		{name: "ts", path: "src/index.ts", expected: "javascript"},
		{name: "py", path: "app.py", expected: "python"},
		{name: "go", path: "main.go", expected: "go"},
		{name: "docker", path: "Dockerfile", expected: "docker"},
		{name: "docker alt", path: "Dockerfile.dev", expected: "docker"},
		{name: "unknown", path: "README.md", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, DetectLanguage(tt.path))
		})
	}
}

func TestNormalizeLanguage(t *testing.T) {
	assert.Equal(t, "javascript", NormalizeLanguage("JavaScript"))
	assert.Equal(t, "", NormalizeLanguage(""))
}

func TestIsSupportedLanguage(t *testing.T) {
	assert.True(t, IsSupportedLanguage("python"))
	assert.False(t, IsSupportedLanguage("haskell"))
}

func TestPatternsExtractKeys(t *testing.T) {
	patterns := Patterns()
	line := "const port = process.env.PORT;"

	var matched bool
	for _, pattern := range patterns {
		if pattern.Language != "javascript" {
			continue
		}
		if pattern.Regex.MatchString(line) {
			matches := pattern.Regex.FindStringSubmatch(line)
			require.Len(t, matches, 2)
			assert.Equal(t, "PORT", matches[1])
			matched = true
		}
	}
	assert.True(t, matched)
}

func TestFileScannerScan(t *testing.T) {
	root := t.TempDir()
	filePath := filepath.Join(root, "src", "index.js")
	require.NoError(t, os.MkdirAll(filepath.Dir(filePath), 0o755))

	content := strings.Join([]string{
		"const port = process.env.PORT;",
		"const db = process.env['DATABASE_URL'];",
		"const secret = process.env[\"JWT_SECRET\"];",
	}, "\n")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0o600))

	scanner := &FileScanner{}
	result, err := scanner.Scan(root, []string{"node_modules"}, []string{"auto"})
	require.NoError(t, err)
	assert.Equal(t, 1, result.ScannedFiles)
	assert.Len(t, result.Usages, 3)

	keys := []string{result.Usages[0].Key, result.Usages[1].Key, result.Usages[2].Key}
	assert.ElementsMatch(t, []string{"PORT", "DATABASE_URL", "JWT_SECRET"}, keys)
}

func TestShouldIgnore(t *testing.T) {
	assert.True(t, shouldIgnore("/root/node_modules/file.js", []string{"node_modules"}))
	assert.False(t, shouldIgnore("/root/src/index.js", []string{"node_modules"}))
}

func TestShouldScanFile(t *testing.T) {
	root := t.TempDir()
	filePath := filepath.Join(root, "main.go")
	require.NoError(t, os.WriteFile(filePath, []byte("package main\n"), 0o600))

	allowed := map[string]struct{}{ "go": {} }
	assert.True(t, shouldScanFile(filePath, allowed))
}
