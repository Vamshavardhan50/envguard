# Contributing to envguard 🛡️

We love your input! We want to make contributing to `envguard` as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features (like scanner regex patterns for new languages!)

---

## 🛠️ Local Development Setup

To work on `envguard` locally, you will need:
- **Go 1.22** or higher.
- A terminal shell.

### Setup Instructions

1. **Fork and Clone the Repository:**
   ```bash
   git clone https://github.com/YOUR-USERNAME/envguard.git
   cd envguard
   ```

2. **Run Tests:**
   Ensure everything is working out-of-the-box by running the Go tests:
   ```bash
   go test ./... -v
   ```

3. **Build the Binary Locally:**
   ```bash
   go build -o envguard .
   ```

---

## 🔄 Contribution Workflow

We active use GitHub Pull Requests:

1. **Fork** the repo and create your branch from `main`.
2. If you've added code that should be tested, **add tests** (table-driven tests are preferred in this repository!).
3. Ensure your Go code is formatted and passes tests:
   ```bash
   go fmt ./...
   go test ./...
   ```
4. Place clear comments explaining non-trivial logic.
5. Submit a pull request!

---

## 📝 Coding Style & Rules

To keep the codebase robust and clean, please adhere to these core rules:
* **Privacy Engine:** Never read, log, or print raw environment variable values. Only process key names.
* **Go Style:** Follow standard Go linting rules. Use explicit errors (`%w` wrapping) and handle all errors.
* **No Snip/Placeholders:** Do not commit `TODO` comments or incomplete code templates.

---

## 🐞 Reporting Issues

* Use our GitHub Issue templates to report bugs or request new features.
* Provide clear, reproducible instructions.
* Do **not** post actual production credentials or environment values in issues. Only post key names (e.g. `DATABASE_URL` is fine, but `postgresql://admin:password@localhost:5432` is not).

Thank you for contributing to `envguard`! 🚀
