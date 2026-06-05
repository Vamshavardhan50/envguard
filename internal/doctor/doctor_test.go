// internal/doctor/doctor_test.go
// Tests for the doctor engine.

package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngineCheckAllGood(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".env", "API_URL=xyz")
	writeFile(t, root, ".env.example", "API_URL=")
	writeFile(t, root, ".gitignore", ".env\n")
	writeFile(t, root, ".envguard.yaml", "version: 1")

	engine := &Engine{}
	res := engine.Check(root)

	assert.Len(t, res.Checks, 6)
	for _, c := range res.Checks {
		assert.Equal(t, "pass", c.Status, c.Name)
	}
}

func TestEngineCheckMissingAndSecrets(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".env", "AWS_SECRET_KEY=123\nDB_PASSWORD=abc")
	writeFile(t, root, ".gitignore", ".env.example\n") // purposefully bad

	engine := &Engine{}
	res := engine.Check(root)

	assert.Equal(t, "pass", findCheck(res.Checks, ".env file"))
	assert.Equal(t, "fail", findCheck(res.Checks, ".env.example file"))
	assert.Equal(t, "fail", findCheck(res.Checks, ".gitignore protects .env"))
	assert.Equal(t, "warn", findCheck(res.Checks, ".env.example not ignored"))
	assert.Equal(t, "warn", findCheck(res.Checks, "Secret detection"))
}

func TestEngineCheckMissingGitignore(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".env", "API_KEY=abc")
	writeFile(t, root, ".env.example", "API_KEY=")

	engine := &Engine{}
	res := engine.Check(root)

	assert.Equal(t, "fail", findCheck(res.Checks, ".gitignore protects .env"))
	assert.Equal(t, "warn", findCheck(res.Checks, ".env.example not ignored"))
}

func writeFile(t *testing.T, dir, name, content string) {
	err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600)
	require.NoError(t, err)
}

func findCheck(checks []DoctorCheck, name string) string {
	for _, c := range checks {
		if c.Name == name {
			return c.Status
		}
	}
	return ""
}
