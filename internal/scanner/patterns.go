// internal/scanner/patterns.go
// Defines regex patterns for env var usage across languages.

package scanner

import "regexp"

// Pattern describes a regex pattern and language for env usage detection.
type Pattern struct {
	Language string
	Regex    *regexp.Regexp
}

// Patterns returns the compiled patterns for supported languages.
func Patterns() []Pattern {
	return []Pattern{
		{
			Language: "javascript",
			Regex:    regexp.MustCompile(`process\.env\.([A-Z_][A-Z0-9_]*)`),
		},
		{
			Language: "javascript",
			Regex:    regexp.MustCompile(`process\.env\[['\"]([A-Z_][A-Z0-9_]*)['\"]\]`),
		},
		{
			Language: "javascript",
			Regex:    regexp.MustCompile(`import\.meta\.env\.([A-Z_][A-Z0-9_]*)`),
		},
		{
			Language: "javascript",
			Regex:    regexp.MustCompile(`import\.meta\.env\[['\"]([A-Z_][A-Z0-9_]*)['\"]\]`),
		},
		{
			Language: "python",
			Regex:    regexp.MustCompile(`os\.environ\.get\(['\"]([A-Z_][A-Z0-9_]*)['\"]\)`),
		},
		{
			Language: "python",
			Regex:    regexp.MustCompile(`os\.environ\[['\"]([A-Z_][A-Z0-9_]*)['\"]\]`),
		},
		{
			Language: "python",
			Regex:    regexp.MustCompile(`os\.getenv\(['\"]([A-Z_][A-Z0-9_]*)['\"]\)`),
		},
		{
			Language: "go",
			Regex:    regexp.MustCompile(`os\.Getenv\(['\"]([A-Z_][A-Z0-9_]*)['\"]\)`),
		},
		{
			Language: "go",
			Regex:    regexp.MustCompile(`os\.LookupEnv\(['\"]([A-Z_][A-Z0-9_]*)['\"]\)`),
		},
		{
			Language: "ruby",
			Regex:    regexp.MustCompile(`ENV\[['\"]([A-Z_][A-Z0-9_]*)['\"]\]`),
		},
		{
			Language: "ruby",
			Regex:    regexp.MustCompile(`ENV\.fetch\(['\"]([A-Z_][A-Z0-9_]*)['\"]\)`),
		},
		{
			Language: "rust",
			Regex:    regexp.MustCompile(`std::env::var\(['\"]([A-Z_][A-Z0-9_]*)['\"]\)`),
		},
		{
			Language: "rust",
			Regex:    regexp.MustCompile(`env!\(['\"]([A-Z_][A-Z0-9_]*)['\"]\)`),
		},
		{
			Language: "php",
			Regex:    regexp.MustCompile(`getenv\(['\"]([A-Z_][A-Z0-9_]*)['\"]\)`),
		},
		{
			Language: "php",
			Regex:    regexp.MustCompile(`\$_ENV\[['\"]([A-Z_][A-Z0-9_]*)['\"]\]`),
		},
		{
			Language: "php",
			Regex:    regexp.MustCompile(`\$_SERVER\[['\"]([A-Z_][A-Z0-9_]*)['\"]\]`),
		},
		{
			Language: "java",
			Regex:    regexp.MustCompile(`System\.getenv\(['\"]([A-Z_][A-Z0-9_]*)['\"]\)`),
		},
		{
			Language: "shell",
			Regex:    regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`),
		},
		{
			Language: "shell",
			Regex:    regexp.MustCompile(`\$([A-Z_][A-Z0-9_]*)`),
		},
		{
			Language: "docker",
			Regex:    regexp.MustCompile(`^\s*ARG\s+([A-Z_][A-Z0-9_]*)`),
		},
		{
			Language: "docker",
			Regex:    regexp.MustCompile(`^\s*ENV\s+([A-Z_][A-Z0-9_]*)`),
		},
	}
}
