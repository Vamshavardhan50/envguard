// internal/reporter/github.go
// GitHub Actions output rendering.

package reporter

import (
	"fmt"
	"strings"
)

// GitHubReporter renders GitHub Actions annotations.
type GitHubReporter struct{}

// RenderAudit emits GitHub Actions annotations for audit results.
func (g *GitHubReporter) RenderAudit(report AuditResult) string {
	var b strings.Builder
	for _, missing := range report.Missing {
		if len(missing.UsedIn) == 0 {
			b.WriteString(fmt.Sprintf("::error::Missing env var: %s\n", missing.Key))
			continue
		}
		usage := missing.UsedIn[0]
		b.WriteString(fmt.Sprintf("::error file=%s,line=%d::Missing env var: %s\n", usage.File, usage.Line, missing.Key))
	}

	for _, unused := range report.Unused {
		b.WriteString(fmt.Sprintf("::warning::Unused env var: %s\n", unused.Key))
	}

	return b.String()
}
