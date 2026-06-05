// internal/validator/validator_test.go
// Tests for the validator engine.

package validator

import (
	"testing"

	"github.com/Vamshavardhan50/envguard/internal/parser"
	"github.com/Vamshavardhan50/envguard/pkg/envfile"
	"github.com/stretchr/testify/assert"
)

func TestEngineValidate(t *testing.T) {
	vars := []envfile.EnvVar{
		{Key: "PORT", Value: "8080"},
		{Key: "NODE_ENV", Value: "staging"},
		{Key: "DB_URL", Value: "postgres://host"},
		{Key: "DEBUG", Value: "yes"}, // invalid bool
	}

	schema := map[string]parser.RuleConfig{
		"PORT": {
			Type: "number",
		},
		"NODE_ENV": {
			Type:   "enum",
			Values: []string{"development", "production"},
		},
		"DB_URL": {
			Required: true,
			Type:     "url",
		},
		"DEBUG": {
			Type: "boolean",
		},
		"MISSING": {
			Required: true,
		},
	}

	engine := &Engine{}
	res := engine.Validate(vars, schema)

	assert.ElementsMatch(t, []string{"PORT", "DB_URL"}, res.Valid)
	requireInvalidKey(t, res.Invalid, "NODE_ENV", "enum")
	requireInvalidKey(t, res.Invalid, "DEBUG", "type")
	requireInvalidKey(t, res.Invalid, "MISSING", "required")
}

func requireInvalidKey(t *testing.T, fails []ValidationError, key, rule string) {
	for _, f := range fails {
		if f.Key == key && f.Rule == rule {
			return
		}
	}
	t.Fatalf("expected validation failure for key=%s rule=%s", key, rule)
}
