// internal/scanner/languages.go
// Detects languages based on file path and known extensions.

package scanner

import (
	"path/filepath"
	"strings"
)

type languageMatcher struct {
	Language   string
	Extensions []string
	Filenames  []string
}

func languageMatchers() []languageMatcher {
	return []languageMatcher{
		{
			Language: "javascript",
			Extensions: []string{
				".js",
				".jsx",
				".ts",
				".tsx",
				".mjs",
				".cjs",
			},
		},
		{
			Language:   "python",
			Extensions: []string{".py"},
		},
		{
			Language:   "go",
			Extensions: []string{".go"},
		},
		{
			Language:   "ruby",
			Extensions: []string{".rb"},
		},
		{
			Language:   "rust",
			Extensions: []string{".rs"},
		},
		{
			Language:   "php",
			Extensions: []string{".php"},
		},
		{
			Language:   "java",
			Extensions: []string{".java"},
		},
		{
			Language: "shell",
			Extensions: []string{
				".sh",
				".bash",
				".zsh",
				".ksh",
			},
		},
		{
			Language:  "docker",
			Filenames: []string{"dockerfile"},
			Extensions: []string{
				".dockerfile",
			},
		},
	}
}

// SupportedLanguages returns the list of supported language identifiers.
func SupportedLanguages() []string {
	return []string{
		"javascript",
		"python",
		"go",
		"ruby",
		"rust",
		"php",
		"java",
		"shell",
		"docker",
	}
}

// DetectLanguage returns the language name for the given path, or empty if unknown.
func DetectLanguage(path string) string {
	base := strings.ToLower(filepath.Base(path))
	if base == "" {
		return ""
	}
	if strings.HasPrefix(base, "dockerfile.") {
		return "docker"
	}

	ext := strings.ToLower(filepath.Ext(base))
	for _, matcher := range languageMatchers() {
		for _, name := range matcher.Filenames {
			if base == name {
				return matcher.Language
			}
		}
		for _, extension := range matcher.Extensions {
			if ext == extension {
				return matcher.Language
			}
		}
	}
	return ""
}

// NormalizeLanguage normalizes user-provided language values.
func NormalizeLanguage(language string) string {
	return strings.ToLower(strings.TrimSpace(language))
}

// IsSupportedLanguage reports whether a language is supported.
func IsSupportedLanguage(language string) bool {
	normalized := NormalizeLanguage(language)
	for _, candidate := range SupportedLanguages() {
		if normalized == candidate {
			return true
		}
	}
	return false
}
