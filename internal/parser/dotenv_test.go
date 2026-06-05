// internal/parser/dotenv_test.go
// Tests for dotenv parsing helpers.

package parser

import (
	"os"
	"strings"
	"testing"

	"github.com/Vamshavardhan50/envguard/pkg/envfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDotEnvReader(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantKeys    []string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid content",
			input:    "PORT=3000\nDEBUG=true\n",
			wantKeys: []string{"PORT", "DEBUG"},
			wantErr:  false,
		},
		{
			name:        "invalid line",
			input:       "MISSING\n",
			wantErr:     true,
			errContains: "missing '=' delimiter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars, err := ParseDotEnvReader(strings.NewReader(tt.input), "test.env")
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantKeys, collectKeys(vars))
		})
	}
}

func TestParseDotEnvReaderNil(t *testing.T) {
	_, err := ParseDotEnvReader(nil, "test.env")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "reader is required")
}

func TestParseDotEnv(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "envguard-*.env")
	require.NoError(t, err)

	_, err = tmpFile.WriteString("ONE=1\nTWO=2\n")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	vars, err := ParseDotEnv(tmpFile.Name())
	require.NoError(t, err)
	assert.Equal(t, []string{"ONE", "TWO"}, collectKeys(vars))
}

func TestParseDotEnvEmptyPath(t *testing.T) {
	_, err := ParseDotEnv("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path is required")
}

func TestMapByKey(t *testing.T) {
	vars := []envfile.EnvVar{
		{Key: "PORT", Value: "3000"},
		{Key: "PORT", Value: "4000"},
		{Key: "DEBUG", Value: "true"},
	}

	mapped := MapByKey(vars)
	assert.Equal(t, "4000", mapped["PORT"])
	assert.Equal(t, "true", mapped["DEBUG"])
}

func collectKeys(vars []envfile.EnvVar) []string {
	keys := make([]string, 0, len(vars))
	for _, entry := range vars {
		keys = append(keys, entry.Key)
	}
	return keys
}
