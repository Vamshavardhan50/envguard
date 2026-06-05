// pkg/envfile/envfile_test.go
// Tests for envfile parsing utilities.

package envfile

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseReaderSuccess(t *testing.T) {
	content := strings.Join([]string{
		"# comment",
		"export PORT=3000",
		"DATABASE_URL=postgres://localhost/db",
		"QUOTED=\"hello world\"",
		"SINGLE='value'",
		"INLINE=foo # comment",
		"SEMICOLON=bar ; trailing",
		"ESCAPED=\"line\\nbreak\"",
	}, "\n")

	vars, err := ParseReader(strings.NewReader(content), "test.env")
	require.NoError(t, err)
	require.Len(t, vars, 7)

	expected := []EnvVar{
		{Key: "PORT", Value: "3000", File: "test.env", Line: 2},
		{Key: "DATABASE_URL", Value: "postgres://localhost/db", File: "test.env", Line: 3},
		{Key: "QUOTED", Value: "hello world", File: "test.env", Line: 4},
		{Key: "SINGLE", Value: "value", File: "test.env", Line: 5},
		{Key: "INLINE", Value: "foo", File: "test.env", Line: 6},
		{Key: "SEMICOLON", Value: "bar", File: "test.env", Line: 7},
		{Key: "ESCAPED", Value: "line\nbreak", File: "test.env", Line: 8},
	}

	assert.Equal(t, expected, vars)
}

func TestParseReaderErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		substring   string
		substring2  string
	}{
		{
			name:       "missing delimiter",
			input:      "MISSING",
			substring:  "missing '=' delimiter",
			substring2: "line 1",
		},
		{
			name:       "invalid key",
			input:      "1INVALID=bad",
			substring:  "invalid env var name",
			substring2: "line 1",
		},
		{
			name:       "unterminated double quote",
			input:      "BAD=\"no end",
			substring:  "unterminated double-quoted value",
			substring2: "line 1",
		},
		{
			name:       "unterminated single quote",
			input:      "BAD='no end",
			substring:  "unterminated single-quoted value",
			substring2: "line 1",
		},
		{
			name:       "dangling escape",
			input:      "BAD=\"value\\\"",
			substring:  "dangling escape sequence",
			substring2: "line 1",
		},
		{
			name:       "empty export",
			input:      "export ",
			substring:  "empty export directive",
			substring2: "line 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseReader(strings.NewReader(tt.input), "bad.env")
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.substring)
			assert.Contains(t, err.Error(), tt.substring2)
			assert.Contains(t, err.Error(), "bad.env")
		})
	}
}

func TestParseFile(t *testing.T) {
	tmpFile, err := os.CreateTemp(t.TempDir(), "envguard-*.env")
	require.NoError(t, err)

	content := "ONE=1\nTWO=2\n"
	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	vars, err := ParseFile(tmpFile.Name())
	require.NoError(t, err)
	require.Len(t, vars, 2)
	assert.Equal(t, "ONE", vars[0].Key)
	assert.Equal(t, "TWO", vars[1].Key)
}

func TestKeys(t *testing.T) {
	vars := []EnvVar{
		{Key: "B"},
		{Key: "A"},
		{Key: "A"},
		{Key: "C"},
	}

	keys := Keys(vars)
	assert.Equal(t, []string{"A", "B", "C"}, keys)
}

func TestIndexByKey(t *testing.T) {
	vars := []EnvVar{
		{Key: "PORT", Value: "3000"},
		{Key: "PORT", Value: "4000"},
		{Key: "DEBUG", Value: "true"},
	}

	indexed := IndexByKey(vars)
	assert.Equal(t, "4000", indexed["PORT"].Value)
	assert.Equal(t, "true", indexed["DEBUG"].Value)
}
