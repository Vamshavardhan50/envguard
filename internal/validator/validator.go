// internal/validator/validator.go
// Validates environment variables against a predefined schema.

package validator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Vamshavardhan50/envguard/internal/parser"
	"github.com/Vamshavardhan50/envguard/pkg/envfile"
)

// ValidationError represents a single validation failure.
type ValidationError struct {
	Key     string
	Rule    string
	Message string
}

// ValidationResult captures the outcome of a validation run.
type ValidationResult struct {
	Valid   []string
	Invalid []ValidationError
}

// Validator defines the interface for environment variable validation.
type Validator interface {
	Validate(vars []envfile.EnvVar, schema map[string]parser.RuleConfig) ValidationResult
}

// Engine implements the Validator interface.
type Engine struct{}

// Validate checks a set of environment variables against defined rules.
func (e *Engine) Validate(vars []envfile.EnvVar, schema map[string]parser.RuleConfig) ValidationResult {
	result := ValidationResult{
		Valid:   make([]string, 0, len(vars)),
		Invalid: make([]ValidationError, 0, len(vars)),
	}

	provided := envfile.IndexByKey(vars)

	for key, rule := range schema {
		entry, exists := provided[key]

		// 1. Required check
		if rule.Required && (!exists || entry.Value == "") {
			result.Invalid = append(result.Invalid, ValidationError{
				Key:     key,
				Rule:    "required",
				Message: "missing required environment variable",
			})
			continue
		}

		// Skip type/enum checks if not present (and not required/empty)
		if !exists || entry.Value == "" {
			continue
		}

		// 2. Enum check (if type is enum)
		if rule.Type == "enum" {
			if !isValidEnum(entry.Value, rule.Values) {
				result.Invalid = append(result.Invalid, ValidationError{
					Key:     key,
					Rule:    "enum",
					Message: fmt.Sprintf("invalid value, expected one of: %s", strings.Join(rule.Values, ", ")),
				})
				continue
			}
		}

		// 3. Type check
		if rule.Type != "" && rule.Type != "enum" {
			if err := checkType(entry.Value, rule.Type); err != nil {
				result.Invalid = append(result.Invalid, ValidationError{
					Key:     key,
					Rule:    "type",
					Message: err.Error(),
				})
				continue
			}
		}

		result.Valid = append(result.Valid, key)
	}

	return result
}

func isValidEnum(val string, allowed []string) bool {
	for _, a := range allowed {
		if val == a {
			return true
		}
	}
	return false
}

func checkType(val, typ string) error {
	switch typ {
	case "string":
		return nil
	case "number":
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return fmt.Errorf("expected numeric value")
		}
	case "boolean":
		v := strings.ToLower(val)
		if v != "true" && v != "false" && v != "1" && v != "0" {
			return fmt.Errorf("expected boolean value (true/false/1/0)")
		}
	case "url":
		if !looksLikeURL(val) {
			return fmt.Errorf("expected valid URL")
		}
	default:
		return fmt.Errorf("unknown type: %s", typ)
	}
	return nil
}

func looksLikeURL(val string) bool {
	if val == "" {
		return false
	}
	sep := strings.Index(val, "://")
	if sep <= 0 {
		return false
	}
	scheme := val[:sep]
	if !isValidScheme(scheme) {
		return false
	}
	remain := val[sep+3:]
	if remain == "" {
		return false
	}
	host := remain
	if slash := strings.IndexByte(remain, '/'); slash >= 0 {
		host = remain[:slash]
	}
	return host != ""
}

func isValidScheme(scheme string) bool {
	if scheme == "" {
		return false
	}
	first := scheme[0]
	if !isAlpha(first) {
		return false
	}
	for i := 1; i < len(scheme); i++ {
		ch := scheme[i]
		if isAlphaNumeric(ch) || ch == '+' || ch == '-' || ch == '.' {
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
