// internal/scanner/scanner.go
// Scans source files to detect environment variable usage.

package scanner

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// VarUsage captures a single environment variable usage in code.
type VarUsage struct {
	Key      string
	File     string
	Line     int
	Language string
}

// ScanResult captures the result of a scan operation.
type ScanResult struct {
	Usages       []VarUsage
	ScannedFiles int
	DurationMs   int64
}

// Scanner defines the scanning interface for dependency injection.
type Scanner interface {
	Scan(path string, ignore []string, languages []string) (ScanResult, error)
}

// FileScanner implements scanning across a filesystem.
type FileScanner struct {
	Patterns []Pattern
}

// Scan walks the provided path and finds env var usages.
func (s *FileScanner) Scan(path string, ignore []string, languages []string) (ScanResult, error) {
	start := time.Now()
	if path == "" {
		return ScanResult{}, fmt.Errorf("scan path is required")
	}

	cleanPath := filepath.Clean(path)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return ScanResult{}, fmt.Errorf("stat scan path: %w", err)
	}

	allowed := normalizeLanguages(languages)
	patterns := s.Patterns
	if len(patterns) == 0 {
		patterns = Patterns()
	}

	jobs := make(chan string, runtime.NumCPU()*2)
	errs := make(chan error, 1)
	var mu sync.Mutex
	allUsages := make([]VarUsage, 0, 128)

	var wg sync.WaitGroup
	workerCount := runtime.NumCPU()
	if workerCount < 2 {
		workerCount = 2
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range jobs {
				if len(errs) > 0 {
					continue
				}
				usages, err := scanFile(filePath, patterns, allowed)
				if err != nil {
					select {
					case errs <- err:
					default:
					}
					continue
				}
				if len(usages) > 0 {
					mu.Lock()
					allUsages = append(allUsages, usages...)
					mu.Unlock()
				}
			}
		}()
	}

	var scanErr error
	var scannedFiles int
	walkErr := filepath.WalkDir(cleanPath, func(current string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if current == cleanPath && info.IsDir() {
			return nil
		}
		if shouldIgnore(current, ignore) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		if !shouldScanFile(current, allowed) {
			return nil
		}
		scannedFiles++
		jobs <- current
		return nil
	})

	close(jobs)
	wg.Wait()

	select {
	case scanErr = <-errs:
	default:
	}

	if walkErr != nil {
		return ScanResult{}, fmt.Errorf("walk scan path: %w", walkErr)
	}
	if scanErr != nil {
		return ScanResult{}, scanErr
	}

	sort.Slice(allUsages, func(i, j int) bool {
		if allUsages[i].File == allUsages[j].File {
			return allUsages[i].Line < allUsages[j].Line
		}
		return allUsages[i].File < allUsages[j].File
	})

	duration := time.Since(start).Milliseconds()
	return ScanResult{
		Usages:       allUsages,
		ScannedFiles: scannedFiles,
		DurationMs:   duration,
	}, nil
}

func scanFile(path string, patterns []Pattern, allowed map[string]struct{}) ([]VarUsage, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", path, err)
	}
	defer func() {
		_ = file.Close()
	}()

	language := DetectLanguage(path)
	if language == "" {
		return nil, nil
	}
	if _, ok := allowed[language]; !ok {
		return nil, nil
	}

	matching := filterPatterns(patterns, language)
	if len(matching) == 0 {
		return nil, nil
	}

	usages := make([]VarUsage, 0, 8)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		for _, pattern := range matching {
			matches := pattern.Regex.FindAllStringSubmatchIndex(line, -1)
			for _, match := range matches {
				if len(match) < 4 {
					continue
				}
				key := line[match[2]:match[3]]
				usages = append(usages, VarUsage{
					Key:      key,
					File:     path,
					Line:     lineNum,
					Language: language,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan file %s: %w", path, err)
	}

	return usages, nil
}

func filterPatterns(patterns []Pattern, language string) []Pattern {
	matching := make([]Pattern, 0, len(patterns))
	for _, pattern := range patterns {
		if pattern.Language == language {
			matching = append(matching, pattern)
		}
	}
	return matching
}

func normalizeLanguages(languages []string) map[string]struct{} {
	allowed := make(map[string]struct{})
	if len(languages) == 0 {
		allowed["auto"] = struct{}{}
		return allowed
	}
	for _, lang := range languages {
		lang = NormalizeLanguage(lang)
		if lang == "" {
			continue
		}
		allowed[lang] = struct{}{}
	}
	if _, ok := allowed["auto"]; ok {
		allowed = make(map[string]struct{})
		for _, lang := range SupportedLanguages() {
			allowed[lang] = struct{}{}
		}
	}
	return allowed
}

func shouldIgnore(path string, ignore []string) bool {
	if len(ignore) == 0 {
		return false
	}
	normalizedPath := filepath.ToSlash(path)
	segments := strings.Split(normalizedPath, "/")
	for _, entry := range ignore {
		entry = filepath.ToSlash(strings.TrimSpace(entry))
		if entry == "" {
			continue
		}
		if strings.Contains(entry, "/") {
			if strings.Contains(normalizedPath, entry) {
				return true
			}
		} else {
			for _, segment := range segments {
				if segment == entry {
					return true
				}
			}
		}
	}
	return false
}

func shouldScanFile(path string, allowed map[string]struct{}) bool {
	language := DetectLanguage(path)
	if language == "" {
		return false
	}
	if _, ok := allowed[language]; !ok {
		return false
	}
	if fileSizeZero(path) {
		return false
	}
	return true
}

func fileSizeZero(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return true
	}
	return info.Size() == 0
}
