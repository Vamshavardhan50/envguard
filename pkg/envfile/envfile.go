// pkg/envfile/envfile.go
// Provides parsing utilities for .env files.

package envfile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// EnvVar represents a single environment variable entry parsed from a file.
type EnvVar struct {
	Key         string
	Value       string
	File        string
	Line        int
	Description string
}

// ParseFile reads a .env file from disk and returns the parsed entries.
func ParseFile(path string) ([]EnvVar, error) {
	cleanPath := filepath.Clean(path)
	file, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("open env file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	vars, err := ParseReader(file, cleanPath)
	if err != nil {
		return nil, err
	}

	return vars, nil
}

// ParseReader reads .env content from a reader and returns the parsed entries.
func ParseReader(reader io.Reader, filename string) ([]EnvVar, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	vars := make([]EnvVar, 0, 64)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if line == "export" || strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export"))
			if line == "" {
				return nil, parseLineError(filename, lineNum, errors.New("empty export directive"))
			}
		}

		key, rawValue, ok := splitKeyValue(line)
		if !ok {
			return nil, parseLineError(filename, lineNum, errors.New("missing '=' delimiter"))
		}
		if !isValidKey(key) {
			return nil, parseLineError(filename, lineNum, errors.New("invalid env var name"))
		}

		value, err := parseValue(rawValue)
		if err != nil {
			return nil, parseLineError(filename, lineNum, err)
		}

		vars = append(vars, EnvVar{
			Key:  key,
			Value: value,
			File: filename,
			Line: lineNum,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}

	return vars, nil
}

// Keys returns unique, sorted keys from the parsed variables.
func Keys(vars []EnvVar) []string {
	unique := make(map[string]struct{}, len(vars))
	for _, entry := range vars {
		if entry.Key == "" {
			continue
		}
		unique[entry.Key] = struct{}{}
	}

	keys := make([]string, 0, len(unique))
	for key := range unique {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// IndexByKey returns a map keyed by env var name, with the last entry winning.
func IndexByKey(vars []EnvVar) map[string]EnvVar {
	indexed := make(map[string]EnvVar, len(vars))
	for _, entry := range vars {
		if entry.Key == "" {
			continue
		}
		indexed[entry.Key] = entry
	}
	return indexed
}

func splitKeyValue(line string) (string, string, bool) {
	idx := strings.IndexRune(line, '=')
	if idx == -1 {
		return "", "", false
	}
	key := strings.TrimSpace(line[:idx])
	value := strings.TrimSpace(line[idx+1:])
	return key, value, true
}

func parseValue(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}
	if raw[0] == '"' {
		if raw[len(raw)-1] != '"' {
			return "", errors.New("unterminated double-quoted value")
		}
		return unescapeDoubleQuoted(raw[1 : len(raw)-1])
	}
	if raw[0] == '\'' {
		if raw[len(raw)-1] != '\'' {
			return "", errors.New("unterminated single-quoted value")
		}
		return raw[1 : len(raw)-1], nil
	}

	return stripInlineComment(raw), nil
}

func stripInlineComment(value string) string {
	for i := 0; i < len(value); i++ {
		if value[i] == '#' || value[i] == ';' {
			if i > 0 && isWhitespace(value[i-1]) {
				return strings.TrimSpace(value[:i])
			}
		}
	}
	return strings.TrimSpace(value)
}

func unescapeDoubleQuoted(value string) (string, error) {
	var builder strings.Builder
	builder.Grow(len(value))

	for i := 0; i < len(value); i++ {
		ch := value[i]
		if ch != '\\' {
			builder.WriteByte(ch)
			continue
		}
		if i == len(value)-1 {
			return "", errors.New("dangling escape sequence")
		}
		i++
		switch value[i] {
		case 'n':
			builder.WriteByte('\n')
		case 'r':
			builder.WriteByte('\r')
		case 't':
			builder.WriteByte('\t')
		case '\\':
			builder.WriteByte('\\')
		case '"':
			builder.WriteByte('"')
		default:
			builder.WriteByte(value[i])
		}
	}

	return builder.String(), nil
}

func isValidKey(key string) bool {
	if key == "" {
		return false
	}
	first := key[0]
	if !isAlpha(first) && first != '_' {
		return false
	}
	for i := 1; i < len(key); i++ {
		ch := key[i]
		if isAlphaNumeric(ch) || ch == '_' || ch == '-' || ch == '.' {
			continue
		}
		return false
	}
	return true
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || (ch >= '0' && ch <= '9')
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

func parseLineError(filename string, line int, err error) error {
	return fmt.Errorf("parse env file %s line %d: %w", filename, line, err)
}
