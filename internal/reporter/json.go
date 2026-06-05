// internal/reporter/json.go
// JSON output rendering for envguard.

package reporter

import (
	"encoding/json"
	"fmt"

	"github.com/Vamshavardhan50/envguard/internal/differ"
	"github.com/Vamshavardhan50/envguard/internal/doctor"
	"github.com/Vamshavardhan50/envguard/internal/validator"
)

// JSONReporter renders machine-readable JSON output.
type JSONReporter struct{}

// RenderAudit renders audit results as JSON.
func (j *JSONReporter) RenderAudit(report AuditResult) (string, error) {
	payload := map[string]any{
		"command":       "audit",
		"scanned_files": report.ScannedFiles,
		"duration_ms":   report.DurationMs,
		"missing":       report.Missing,
		"unused":        report.Unused,
		"valid":         report.Valid,
		"errors":        report.ErrorCount,
		"warnings":      report.WarnCount,
	}

	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal audit json: %w", err)
	}
	return string(bytes), nil
}

// RenderDiff renders diff results as JSON.
func (j *JSONReporter) RenderDiff(result differ.DiffResult) (string, error) {
	payload := map[string]any{
		"command":        "diff",
		"only_in_first":  result.OnlyInFirst,
		"only_in_second": result.OnlyInSecond,
		"in_both":        result.InBoth,
	}

	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal diff json: %w", err)
	}
	return string(bytes), nil
}

// RenderDoctor renders doctor results as JSON.
func (j *JSONReporter) RenderDoctor(result doctor.DoctorResult) (string, error) {
	payload := map[string]any{
		"command": "doctor",
		"checks":  result.Checks,
	}

	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal doctor json: %w", err)
	}
	return string(bytes), nil
}

// RenderValidate renders validation results as JSON.
func (j *JSONReporter) RenderValidate(result validator.ValidationResult) (string, error) {
	payload := map[string]any{
		"command": "validate",
		"valid":   result.Valid,
		"invalid": result.Invalid,
	}

	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal validate json: %w", err)
	}
	return string(bytes), nil
}
