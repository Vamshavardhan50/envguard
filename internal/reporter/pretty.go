// internal/reporter/pretty.go
// Pretty terminal output rendering for envguard.

package reporter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/Vamshavardhan50/envguard/internal/differ"
	"github.com/Vamshavardhan50/envguard/internal/doctor"
	"github.com/Vamshavardhan50/envguard/internal/scanner"
	"github.com/Vamshavardhan50/envguard/internal/validator"
)

// MissingVar captures a missing environment variable and its usage locations.
type MissingVar struct {
	Key    string
	UsedIn []scanner.VarUsage
}

// UnusedVar captures an environment variable defined but not used.
type UnusedVar struct {
	Key string
}

// AuditResult represents the output of an audit run.
type AuditResult struct {
	ScannedFiles int
	DurationMs   int64
	Missing      []MissingVar
	Unused       []UnusedVar
	Valid        []string
	ErrorCount   int
	WarnCount    int
}

// PrettyReporter renders human-friendly output to stdout.
type PrettyReporter struct{}

// RenderAudit renders audit results in pretty format.
func (p *PrettyReporter) RenderAudit(result AuditResult) string {
	header := lipgloss.NewStyle().Bold(true)
	border := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	fail := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	warn := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	ok := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)

	var b strings.Builder

	box := border.Render("╭─────────────────────────────────────────╮") + "\n" +
		border.Render("│  envguard audit") + strings.Repeat(" ", 25) + border.Render("│") + "\n" +
		border.Render(fmt.Sprintf("│  Scanned %d files in %dms", result.ScannedFiles, result.DurationMs)) +
		strings.Repeat(" ", max(0, 23-len(fmt.Sprintf("%d", result.ScannedFiles))-len(fmt.Sprintf("%d", result.DurationMs)))) + border.Render("│") + "\n" +
		border.Render("╰─────────────────────────────────────────╯")

	b.WriteString(box)
	b.WriteString("\n\n")

	if len(result.Missing) > 0 {
		b.WriteString(fail.Render("❌ Missing (") + fmt.Sprintf("%d", len(result.Missing)) + fail.Render(")") + "\n")
		for _, m := range result.Missing {
			if len(m.UsedIn) > 0 {
				usage := m.UsedIn[0]
				b.WriteString(fmt.Sprintf("  %s   → used in %s:%d\n", m.Key, usage.File, usage.Line))
			} else {
				b.WriteString(fmt.Sprintf("  %s\n", m.Key))
			}
		}
		b.WriteString("\n")
	}

	if len(result.Unused) > 0 {
		b.WriteString(warn.Render("⚠️  Unused (") + fmt.Sprintf("%d", len(result.Unused)) + warn.Render(")") + "\n")
		for _, u := range result.Unused {
			b.WriteString(fmt.Sprintf("  %s   → defined in .env but never used\n", u.Key))
		}
		b.WriteString("\n")
	}

	if len(result.Valid) > 0 {
		b.WriteString(ok.Render("✅ Valid (") + fmt.Sprintf("%d", len(result.Valid)) + ok.Render(")") + "\n")
		if len(result.Valid) > 5 {
			b.WriteString(fmt.Sprintf("  %s ... and %d more\n", strings.Join(result.Valid[:5], ", "), len(result.Valid)-5))
		} else {
			b.WriteString(fmt.Sprintf("  %s\n", strings.Join(result.Valid, ", ")))
		}
		b.WriteString("\n")
	}

	b.WriteString(border.Render("──────────────────────────────────────────"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Result: %d errors, %d warnings\n", result.ErrorCount, result.WarnCount))
	b.WriteString("Run envguard doctor for full health check\n")

	return header.Render(b.String())
}

// RenderDiff renders diff results in pretty format.
func (p *PrettyReporter) RenderDiff(result differ.DiffResult) string {
	var b strings.Builder
	b.WriteString("envguard diff\n\n")

	writeDiffList(&b, "Only in first", result.OnlyInFirst)
	writeDiffList(&b, "Only in second", result.OnlyInSecond)
	writeDiffList(&b, "In both", result.InBoth)

	return b.String()
}

// RenderDoctor renders doctor results in pretty format.
func (p *PrettyReporter) RenderDoctor(result doctor.DoctorResult) string {
	var b strings.Builder
	b.WriteString("envguard doctor\n\n")

	for _, check := range result.Checks {
		status := strings.ToUpper(check.Status)
		b.WriteString(fmt.Sprintf("[%s] %s - %s\n", status, check.Name, check.Message))
	}

	return b.String()
}

// RenderValidate renders validation results in pretty format.
func (p *PrettyReporter) RenderValidate(result validator.ValidationResult) string {
	var b strings.Builder
	b.WriteString("envguard validate\n\n")
	if len(result.Invalid) == 0 {
		b.WriteString("All variables passed validation.\n")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("Invalid (%d)\n", len(result.Invalid)))
	for _, inv := range result.Invalid {
		b.WriteString(fmt.Sprintf("  %s - %s\n", inv.Key, inv.Message))
	}
	return b.String()
}

// RenderAuditSummary returns a compact summary suitable for CI output.
func (p *PrettyReporter) RenderAuditSummary(result scanner.ScanResult, report AuditResult) string {
	return fmt.Sprintf("scanned=%d duration_ms=%d errors=%d warnings=%d", result.ScannedFiles, result.DurationMs, report.ErrorCount, report.WarnCount)
}

func writeDiffList(b *strings.Builder, title string, items []string) {
	b.WriteString(title + ":\n")
	if len(items) == 0 {
		b.WriteString("  (none)\n\n")
		return
	}

	sorted := append([]string{}, items...)
	sort.Strings(sorted)
	for _, item := range sorted {
		b.WriteString("  " + item + "\n")
	}
	b.WriteString("\n")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
