# envguard — Complete Agent Context

## IDENTITY

You are a senior Go engineer.
Your only job is to build "envguard" completely.
You write production-ready, tested, working Go code.
Nothing else.

---

## HARD RULES — NEVER BREAK THESE

### Rule 1 — No Hallucination

Never invent packages, functions, or APIs.
If unsure → say exactly:
"I am not sure about [X]. My reasoning is [Y]. Please verify."
Never silently guess and present it as fact.

### Rule 2 — Complete Files Only

Never write partial code.
Never write snippets.
Never write pseudocode.
Every file you write must be 100% complete and runnable.

### Rule 3 — No Placeholders

Never write:
// TODO: implement this
// implement later
panic("not implemented")
Every function must be fully implemented.

### Rule 4 — Exact Packages Only

Only use these external packages:
github.com/spf13/cobra
github.com/spf13/viper
github.com/charmbracelet/lipgloss
github.com/joho/godotenv
github.com/stretchr/testify
No other external packages without explicit approval.

### Rule 5 — Privacy Is Non-Negotiable

Never read env var values for output.
Never log env var values.
Never print env var values.
Only process and display key names.
This is a critical security rule.

### Rule 6 — One File At A Time

Build one file per step.
Wait for "next" before moving forward.
Never skip steps.
Never reorder steps.

### Rule 7 — Error Handling Always

Every error must be handled.
Pattern:
if err != nil {
return fmt.Errorf("context here: %w", err)
}
Never ignore errors.
Never use panic() in production paths.

### Rule 8 — Test Everything

Every package needs tests.
Use table-driven tests always.
No file ships without its test file.

### Rule 9 — File Header Always

Every file starts with:
// path/to/filename.go
// What this file does in one sentence.

### Rule 10 — Ask Before Assuming

If a requirement is unclear → ask.
Do not assume and build wrong.

---

## PROJECT IDENTITY

Name: envguard
Tagline: "The ESLint for environment variables."
Language: Go 1.22+
Type: CLI tool
License: MIT

What it does:
Audits, validates, and protects environment variable
configuration in any software project.
Works with any language or framework.
100% offline. Zero config to start.

What problem it solves:
Most teams break production because of:

- Missing env vars
- Config drift between environments
- No validation of env var types or values
- Forgetting to update .env.example
- No CI check for env configuration
  envguard fixes all of this.

Target users:

- Backend developers
- DevOps engineers
- Full-stack developers
- Platform engineers
- Any team using .env files

---

## TECH STACK

### Language

Go 1.22+

### CLI Framework

github.com/spf13/cobra
Used for all command definitions.

### Config Parsing

github.com/spf13/viper
Used for .envguard.yaml parsing.

### Terminal Styling

github.com/charmbracelet/lipgloss
Used for pretty terminal output only.

### Env File Parsing

github.com/joho/godotenv
Used for .env file reading.

### Testing

github.com/stretchr/testify
Used for assertions in all test files.

### Release

goreleaser
Used for cross-platform binary builds.

### Allowed Standard Library

os
fmt
path/filepath
strings
regexp
encoding/json
io
bufio
errors
context
runtime
sync
bytes
sort
strconv
time

---

## PROJECT STRUCTURE

envguard/
├── main.go
├── go.mod
├── go.sum
├── .goreleaser.yaml
├── .envguard.yaml
├── README.md
├── CHANGELOG.md
├── CONTEXT.md
│
├── cmd/
│ ├── root.go
│ ├── audit.go
│ ├── sync.go
│ ├── diff.go
│ ├── doctor.go
│ ├── init.go
│ └── validate.go
│
├── internal/
│ ├── scanner/
│ │ ├── scanner.go
│ │ ├── scanner_test.go
│ │ ├── languages.go
│ │ └── patterns.go
│ ├── parser/
│ │ ├── dotenv.go
│ │ ├── dotenv_test.go
│ │ ├── config.go
│ │ └── config_test.go
│ ├── validator/
│ │ ├── validator.go
│ │ └── validator_test.go
│ ├── differ/
│ │ ├── differ.go
│ │ └── differ_test.go
│ ├── reporter/
│ │ ├── pretty.go
│ │ ├── json.go
│ │ └── github.go
│ └── doctor/
│ ├── doctor.go
│ └── doctor_test.go
│
├── pkg/
│ └── envfile/
│ ├── envfile.go
│ └── envfile_test.go
│
├── distributions/
│ ├── npm/
│ │ ├── bin/
│ │ │ └── envguard.js
│ │ ├── scripts/
│ │ │ └── install.js
│ │ └── package.json
│ └── pip/
│ ├── envguard/
│ │ ├── **init**.py
│ │ └── main.py
│ └── pyproject.toml
│
├── scripts/
│ └── bump-version.sh
│
├── testdata/
│ ├── js-project/
│ │ ├── .env
│ │ └── src/
│ │ └── index.js
│ ├── python-project/
│ │ ├── .env
│ │ └── app.py
│ └── go-project/
│ ├── .env
│ └── main.go
│
└── .github/
└── workflows/
├── ci.yml
└── release.yml

---

## COMMANDS

### audit

Command: envguard audit [path] [flags]

Purpose:
Scan codebase for all env var usages.
Compare against .env file.
Report missing, unused, and valid vars.

Flags:
--path string Path to scan (default ".")
--env string Path to .env file (default ".env")
--config string Path to config (default ".envguard.yaml")
--format string Output format: pretty|json|github (default "pretty")
--ci Exit 1 on any findings
--ignore string Comma-separated paths to ignore

Behavior:

1. Load .env file → get list of defined keys
2. Walk codebase → find all env var usages
3. Compare defined vs used
4. Missing = used in code but not in .env
5. Unused = in .env but never used in code
6. Valid = in .env and used in code
7. Output results in chosen format
8. If --ci flag and errors exist → exit 1

### sync

Command: envguard sync [flags]

Purpose:
Generate .env.example from .env.
Strip values, keep only key names.
Add descriptions from config if available.

Flags:
--env string Source .env file (default ".env")
--output string Output file (default ".env.example")
--force Overwrite without confirmation

Behavior:

1. Read .env file
2. For each line: keep key, strip value
3. Add description comment if in .envguard.yaml
4. If .env.example exists and no --force → ask confirm
5. Write .env.example
6. Never write actual values

### diff

Command: envguard diff <file1> <file2> [flags]

Purpose:
Compare two env files.
Show what is different between them.

Flags:
--format string Output format: pretty|json (default "pretty")
--ci Exit 1 if any differences found

Behavior:

1. Read both files → get key lists only
2. Find keys only in file1
3. Find keys only in file2
4. Find keys in both
5. Output comparison
6. Never compare or show values

### doctor

Command: envguard doctor [path] [flags]

Purpose:
Full project health check.
Checks file presence, gitignore rules, common issues.

Flags:
--path string Path to check (default ".")
--format string Output format: pretty|json (default "pretty")

Checks performed:

1. .env file exists → fail if missing
2. .env.example file exists → warn if missing
3. .env is in .gitignore → fail if not
4. .env.example is NOT in .gitignore → warn if it is
5. .envguard.yaml exists → warn if missing
6. .env has at least one var → warn if empty
7. No obvious secret key names exposed → warn if found

### init

Command: envguard init [flags]

Purpose:
Interactive setup wizard.
Creates .envguard.yaml for the project.

Flags:
--yes Skip prompts, use all defaults
--path string Path to initialize (default ".")

Behavior:

1. Detect project language automatically
2. Read existing .env if present
3. For each key in .env → ask for type and description
4. Write .envguard.yaml with results
5. If --yes → use defaults for everything

### validate

Command: envguard validate [flags]

Purpose:
Validate .env against schema in .envguard.yaml.

Flags:
--env string Path to .env file (default ".env")
--config string Path to config (default ".envguard.yaml")
--format string Output format: pretty|json (default "pretty")
--ci Exit 1 on any validation errors

Validation rules:
required: true → key must exist in .env
type: string → any non-empty string
type: number → must be valid integer or float
type: url → must be valid URL format
type: boolean → must be true|false|1|0
type: enum → must be one of allowed values
default: "value" → used in output docs only

---

## CORE DATA TYPES

These exact types must be used across all packages.
Do not modify or add fields without updating all usages.

```go
// EnvVar represents a single environment variable.
// Value field must NEVER be logged or printed.
type EnvVar struct {
    Key         string
    Value       string
    File        string
    Line        int
    Description string
}

// VarUsage represents a detected env var usage in code.
type VarUsage struct {
    Key      string
    File     string
    Line     int
    Language string
}

// AuditResult holds the complete result of an audit command.
type AuditResult struct {
    ScannedFiles int
    DurationMs   int64
    Missing      []MissingVar
    Unused       []UnusedVar
    Valid         []string
    ErrorCount   int
    WarnCount    int
}

// MissingVar is an env var used in code but missing from .env.
type MissingVar struct {
    Key    string
    UsedIn []VarUsage
}

// UnusedVar is an env var in .env but never used in code.
type UnusedVar struct {
    Key string
}

// DiffResult holds the result of comparing two env files.
type DiffResult struct {
    OnlyInFirst  []string
    OnlyInSecond []string
    InBoth       []string
}

// DoctorResult holds all health check results.
type DoctorResult struct {
    Checks []DoctorCheck
}

// DoctorCheck is a single health check result.
type DoctorCheck struct {
    Name    string
    Status  string
    Message string
}

// ValidationResult holds the result of schema validation.
type ValidationResult struct {
    Valid   []string
    Invalid []ValidationError
}

// ValidationError describes a single validation failure.
type ValidationError struct {
    Key     string
    Rule    string
    Message string
}

// Config is the parsed .envguard.yaml config file.
type Config struct {
    Version int
    Scan    ScanConfig
    Rules   map[string]RuleConfig
    Output  OutputConfig
}

// ScanConfig holds scanning configuration.
type ScanConfig struct {
    Paths     []string
    Ignore    []string
    Languages []string
}

// RuleConfig holds validation rules for one env var.
type RuleConfig struct {
    Required    bool
    Type        string
    Values      []string
    Default     string
    Description string
}

// OutputConfig holds output configuration.
type OutputConfig struct {
    Format   string
    ExitCode bool
}
```
