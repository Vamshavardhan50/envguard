// internal/parser/dotenv.go
// Parses .env files for use across envguard commands.

package parser

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/Vamshavardhan50/envguard/pkg/envfile"
)

// ParseDotEnv reads a dotenv file from disk and returns parsed variables.
func ParseDotEnv(path string) ([]envfile.EnvVar, error) {
	if path == "" {
		return nil, fmt.Errorf("parse dotenv file: path is required")
	}
	cleanPath := filepath.Clean(path)
	vars, err := envfile.ParseFile(cleanPath)
	if err != nil {
		return nil, err
	}
	return vars, nil
}

// ParseDotEnvReader reads dotenv content from a reader and returns parsed variables.
func ParseDotEnvReader(reader io.Reader, filename string) ([]envfile.EnvVar, error) {
	if reader == nil {
		return nil, fmt.Errorf("parse dotenv reader: reader is required")
	}
	if filename == "" {
		filename = "<reader>"
	}
	vars, err := envfile.ParseReader(reader, filename)
	if err != nil {
		return nil, err
	}
	return vars, nil
}

// MapByKey converts env vars into a key-value map.
func MapByKey(vars []envfile.EnvVar) map[string]string {
	mapped := make(map[string]string, len(vars))
	for _, entry := range vars {
		if entry.Key == "" {
			continue
		}
		mapped[entry.Key] = entry.Value
	}
	return mapped
}
