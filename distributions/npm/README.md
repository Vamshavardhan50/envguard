# envguard 🛡️

[![CI](https://github.com/Vamshavardhan50/envguard/actions/workflows/ci.yml/badge.svg)](https://github.com/Vamshavardhan50/envguard/actions/workflows/ci.yml)
[![Release](https://github.com/Vamshavardhan50/envguard/actions/workflows/release.yml/badge.svg)](https://github.com/Vamshavardhan50/envguard/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Think of `envguard` as a spell-checker/linter for your environment variables.** 
> It scans your code, checks your `.env` files, and makes sure you never break your app in production due to a missing or misconfigured setting. It is fast, works 100% offline, and requires zero configuration to start.

---

## What is envguard?
Have you ever deployed an app, only for it to immediately crash because you forgot to copy a new API key to the server? Or because someone configured `PORT` as a word instead of a number?

`envguard` solves this by automatically finding all environment variables your code uses (like `process.env.DATABASE_URL` or `os.environ.get('PORT')`) and checking them against your `.env` configuration file. It warns you immediately about:
- ❌ **Missing variables** that your code expects but are not configured.
- ⚠️ **Unused variables** in `.env` that your code doesn't actually use.
- 🚫 **Invalid formats** (e.g. database URLs that are not valid URLs, or ports that are not numbers).

---

## ⚡ Quick 1-Minute Start

Get up and running in three simple steps:

### 1. Install it
Run the install command for your favorite package manager:

```bash
# Using Node (NPM)
npm install -g envguard-bin

# Using Python (PIP)
pip install envguard

# Using Go
go install github.com/Vamshavardhan50/envguard@latest
```

### 2. Set it up
Run the setup wizard in your project folder. It will scan your project, find your environment variables, and create a safe configuration file (`.envguard.yaml`) for you:

```bash
envguard init
```

### 3. Guard your project
Scan your project to find any discrepancies:

```bash
envguard audit
```

---

## 🔒 Security & Privacy First

- **100% Offline:** `envguard` never makes network requests and never uploads anything.
- **Privacy-Engineered:** It only reads and displays the **names** of the keys (e.g. `STRIPE_API_KEY`). It **never** reads, logs, or prints the actual secret values.

---

## 🛠️ How to Use: Common Tasks

Here is how to run `envguard`'s most common commands in your project:

### 🔍 Find Missing or Unused Keys
To check if your local `.env` file matches what your code actually uses, run:
```bash
envguard audit
```
*Tip: In your CI/CD pipelines, run `envguard audit --ci` to automatically fail build pipelines if keys are missing.*

---

### 🛡️ Check If Values Are Correct (Validation)
You can define rules in `.envguard.yaml` (e.g., checking if `DATABASE_URL` is a valid URL, or `PORT` is a number). To validate your current configuration against these rules, run:
```bash
envguard validate
```

---

### 📝 Auto-Update Your `.env.example` Template
Stop updating `.env.example` templates manually! Run this command to automatically read your `.env` file and generate a clean `.env.example` containing only the keys (values are safely stripped):
```bash
envguard sync --force
```

---

### 🩺 Run a Project Health Check
To perform a structure audit of your environment file setup (making sure `.env` is inside `.gitignore` so you don't accidentally publish secrets, checking file integrity, etc.), run:
```bash
envguard doctor
```

---

## ⚙️ Advanced Configuration (`.envguard.yaml`)

You can customize how `envguard` behaves by editing your `.envguard.yaml` file. Here is an example:

```yaml
version: 1

scan:
  paths:
    - "."
  ignore:
    - "node_modules"
    - ".git"
    - "dist"
  languages:
    - auto # Auto-detects JavaScript, TypeScript, Python, Go, Rust, Ruby, Dockerfiles, etc.

# Define rules for validating values
rules:
  DATABASE_URL:
    required: true
    type: url
    description: "Primary PostgreSQL database URL"
  PORT:
    required: false
    type: number
    default: "3000"
    description: "The port the web server runs on"
  NODE_ENV:
    required: true
    type: enum
    values:
      - development
      - production
      - test
```

---

## 🤖 Integrate with GitHub Actions

Add `envguard` to your GitHub Actions workflow to block pull requests containing invalid or incomplete configurations:

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
      - name: Run audit
        run: envguard audit --ci --format github
```

---

## 🙋 Frequently Asked Questions (FAQ)

### Does `envguard` send my secrets to a third-party server?
**No.** `envguard` runs entirely on your local machine. It does not send any telemetry, analytics, or credentials over the internet.

### What languages does the code scanner support?
`envguard` scans JavaScript (`process.env.VAR`), TypeScript, React/Vue (`import.meta.env.VAR` or `process.env.VAR`), Python (`os.environ`), Go (`os.Getenv`), Ruby, Rust, PHP, Java, Shell scripts, and Dockerfiles out-of-the-box.

### How is this different from other dotenv validators?
Unlike most tools, `envguard` does not just check if a `.env` file exists. It **statically scans your source code files** to find what keys your code actually references, highlighting code references that are completely missing from your config.

---

## 📄 License
This project is open-source software licensed under the [MIT License](LICENSE).
