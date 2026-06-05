# envguard

[![CI](https://github.com/Vamshavardhan50/envguard/actions/workflows/ci.yml/badge.svg)](https://github.com/Vamshavardhan50/envguard/actions/workflows/ci.yml)
[![Release](https://github.com/Vamshavardhan50/envguard/actions/workflows/release.yml/badge.svg)](https://github.com/Vamshavardhan50/envguard/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **The ESLint for environment variables.** Audits, validates, and protects environment variable configurations in any software project. Fast, offline, and zero-config to start.

---

## Why envguard?

Most development teams break production due to:
- Missing environment variables in deployment.
- Configuration drift between development, staging, and production.
- No validation of environment variable types (e.g. string vs number) or formats.
- Forgetting to keep `.env.example` templates updated.
- Lack of continuous integration checks for environment state.

**envguard** solves all of these with a single static binary. It acts as a static analysis tool for your configuration, verifying that what your code expects matches what is configured.

---

## Core Promises

- **Offline-First:** Never makes network requests. Runs 100% locally.
- **Privacy-Engineered:** Never reads or logs actual environment values; it only parses and displays **key names**.
- **Polyglot Scanner:** Auto-detects variable usages in JavaScript, TypeScript, Python, Go, Ruby, Rust, PHP, Java, Shell, and Dockerfiles.
- **Zero Dependencies:** Compiles to a single static binary.

---

## Installation

### Via NPM (Global)
```bash
npm install -g envguard
```

### Via PIP
```bash
pip install envguard
```

### Via Go
```bash
go install github.com/Vamshavardhan50/envguard@latest
```

---

## Command Reference

### 1. `envguard audit`
Scans your codebase for environment variable usage and compares it against a `.env` file to detect missing or unused variables.

```bash
envguard audit [path] [flags]
```

#### Flags
- `--path string` - Path to scan (default `"."`)
- `--env string` - Path to `.env` file (default `".env"`)
- `--config string` - Path to `.envguard.yaml` (default `".envguard.yaml"`)
- `--format string` - Output format: `pretty|json|github` (default `"pretty"`)
- `--ci` - Exits with status `1` if any missing or unused variables are found.
- `--ignore string` - Comma-separated list of paths to ignore.

---

### 2. `envguard sync`
Reads your local `.env` file and outputs an updated `.env.example` containing only the keys and stripping all secret values. It optionally appends key descriptions defined in `.envguard.yaml`.

```bash
envguard sync [flags]
```

#### Flags
- `--env string` - Source `.env` file (default `".env"`)
- `--output string` - Output template file (default `".env.example"`)
- `--force` - Overwrite `.env.example` without asking for confirmation.

---

### 3. `envguard diff`
Compares key presence between two environment files without exposing their values. Useful for checking differences between `.env` and `.env.production`.

```bash
envguard diff <file1> <file2> [flags]
```

#### Flags
- `--format string` - Output format: `pretty|json` (default `"pretty"`)
- `--ci` - Exits with status `1` if there are any differences between the two files.

---

### 4. `envguard doctor`
Performs a structural health check on your project's environment variables. It verifies file existence, gitignore rules, and flags potential raw production secrets in `.env` keys.

```bash
envguard doctor [path] [flags]
```

#### Flags
- `--path string` - Path to check (default `"."`)
- `--format string` - Output format: `pretty|json` (default `"pretty"`)

---

### 5. `envguard init`
Initializes a new `.envguard.yaml` configuration file for your project. It automatically scans your primary programming language and pre-populates config rules based on your existing `.env` file keys.

```bash
envguard init [flags]
```

#### Flags
- `--yes` - Skip interactive prompts and write default configuration values.
- `--path string` - Path to initialize (default `"."`)

---

### 6. `envguard validate`
Validates `.env` variables against a schema defined in `.envguard.yaml`. It checks for missing required variables, incorrect types (string, number, boolean, url), and enum constraints.

```bash
envguard validate [flags]
```

#### Flags
- `--env string` - Path to `.env` file (default `".env"`)
- `--config string` - Path to `.envguard.yaml` (default `".envguard.yaml"`)
- `--format string` - Output format: `pretty|json` (default `"pretty"`)
- `--ci` - Exits with status `1` if any validation checks fail.

---

## Configuration Schema (`.envguard.yaml`)

Initialize a configuration file with `envguard init`. Below is an example `.envguard.yaml`:

```yaml
version: 1

scan:
  paths:
    - "."
  ignore:
    - "node_modules"
    - ".git"
    - "dist"
    - "build"
    - "vendor"
    - ".goreleaser.yaml"
  languages:
    - auto # Auto-detect project languages or specify ["javascript", "go", "python"]

rules:
  DATABASE_URL:
    required: true
    type: url
    description: "Primary database connection string"
  PORT:
    required: false
    type: number
    default: "3000"
    description: "Server port"
  NODE_ENV:
    required: true
    type: enum
    values:
      - development
      - production
      - test
    description: "Application environment"

output:
  format: pretty
  exit_code: true
```

---

## Verification & CI Integration

### GitHub Actions Workflow Example

You can run `envguard` directly in your GitHub Actions workflows to block pull requests containing invalid configuration:

```yaml
name: Guard Environment

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - name: Install envguard
        run: go install github.com/Vamshavardhan50/envguard@latest
      - name: Run audit in CI Mode
        run: envguard audit --ci --format github
```

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
