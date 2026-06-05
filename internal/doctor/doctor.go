// internal/doctor/doctor.go
// Performs project health checks for environment configuration.

package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Vamshavardhan50/envguard/pkg/envfile"
)

// DoctorResult holds the results of all health checks.
type DoctorResult struct {
	Checks []DoctorCheck
}

// DoctorCheck represents a single health check result.
type DoctorCheck struct {
	Name    string
	Status  string // "pass", "warn", "fail"
	Message string
}

// Engine implements environment health checks.
type Engine struct{}

// Check checks the health of the project located at basePath.
func (e *Engine) Check(basePath string) DoctorResult {
	res := DoctorResult{
		Checks: make([]DoctorCheck, 0, 6),
	}

	envExists := checkFileExists(filepath.Join(basePath, ".env"))
	res.Checks = append(res.Checks, DoctorCheck{
		Name:    ".env file",
		Status:  statusBool(envExists, "pass", "warn"),
		Message: messageBool(envExists, "Found .env file", "Missing .env file in project root"),
	})

	exampleExists := checkFileExists(filepath.Join(basePath, ".env.example"))
	res.Checks = append(res.Checks, DoctorCheck{
		Name:    ".env.example file",
		Status:  statusBool(exampleExists, "pass", "fail"),
		Message: messageBool(exampleExists, "Found .env.example file", "Missing .env.example template file"),
	})

	gitignore := filepath.Join(basePath, ".gitignore")
	if checkFileExists(gitignore) {
		content, err := os.ReadFile(gitignore)
		if err != nil {
			res.Checks = append(res.Checks, DoctorCheck{
				Name:    ".gitignore protects .env",
				Status:  "fail",
				Message: "Unable to read .gitignore to verify .env protection",
			})
			res.Checks = append(res.Checks, DoctorCheck{
				Name:    ".env.example not ignored",
				Status:  "warn",
				Message: "Unable to read .gitignore to verify .env.example tracking",
			})
		} else {
			lines := strings.Split(string(content), "\n")
			envIgnored := isIgnored(".env", lines)
			res.Checks = append(res.Checks, DoctorCheck{
				Name:    ".gitignore protects .env",
				Status:  statusBool(envIgnored, "pass", "fail"),
				Message: messageBool(envIgnored, ".env is in .gitignore", ".env must be added to .gitignore"),
			})

			exampleIgnored := isIgnored(".env.example", lines)
			res.Checks = append(res.Checks, DoctorCheck{
				Name:    ".env.example not ignored",
				Status:  statusBool(!exampleIgnored, "pass", "warn"),
				Message: messageBool(!exampleIgnored, ".env.example is not ignored", ".env.example should not be in .gitignore"),
			})
		}
	} else {
		res.Checks = append(res.Checks, DoctorCheck{
			Name:    ".gitignore protects .env",
			Status:  "fail",
			Message: "Missing .gitignore; add .env to .gitignore",
		})
		res.Checks = append(res.Checks, DoctorCheck{
			Name:    ".env.example not ignored",
			Status:  "warn",
			Message: "Missing .gitignore; ensure .env.example is tracked",
		})
	}

	configExists := checkFileExists(filepath.Join(basePath, ".envguard.yaml"))
	res.Checks = append(res.Checks, DoctorCheck{
		Name:    ".envguard.yaml config",
		Status:  statusBool(configExists, "pass", "warn"),
		Message: messageBool(configExists, "Found envguard configuration", "No .envguard.yaml found. Run `envguard init` to create one"),
	})

	if envExists {
		vars, err := envfile.ParseFile(filepath.Join(basePath, ".env"))
		if err != nil {
			res.Checks = append(res.Checks, DoctorCheck{
				Name:    "Secret detection",
				Status:  "warn",
				Message: "Unable to parse .env to inspect key names",
			})
		} else {
			secrets := countSuspectSecrets(vars)
			res.Checks = append(res.Checks, DoctorCheck{
				Name:    "Secret detection",
				Status:  statusBool(secrets == 0, "pass", "warn"),
				Message: messageBool(secrets == 0, "No obvious raw secrets in file", fmt.Sprintf("Found %d keys that might be raw production secrets (e.g. AWS_SECRET_ACCESS_KEY)", secrets)),
			})
		}
	}

	return res
}

func checkFileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isIgnored(target string, lines []string) bool {
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "#") {
			continue
		}
		if l == target || l == "/"+target {
			return true
		}
	}
	return false
}

func countSuspectSecrets(vars []envfile.EnvVar) int {
	matcher := regexp.MustCompile(`(?i)(secret|password|token|api_key|access_key)`)
	count := 0
	for _, v := range vars {
		if matcher.MatchString(v.Key) {
			count++
		}
	}
	return count
}

func statusBool(cond bool, t, f string) string {
	if cond {
		return t
	}
	return f
}

func messageBool(cond bool, t, f string) string {
	if cond {
		return t
	}
	return f
}
