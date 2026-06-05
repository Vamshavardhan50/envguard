// internal/differ/differ.go
// Compares two sets of environment variables.

package differ

import (
	"github.com/Vamshavardhan50/envguard/pkg/envfile"
)

// DiffResult captures the differences between two environment files.
type DiffResult struct {
	OnlyInFirst  []string
	OnlyInSecond []string
	InBoth       []string
}

// Differ defines the interface for environment variable comparison.
type Differ interface {
	Diff(first, second []envfile.EnvVar) DiffResult
}

// Engine implements the Differ interface.
type Engine struct{}

// Diff compares two sets of environment variables.
func (e *Engine) Diff(first, second []envfile.EnvVar) DiffResult {
	keys1 := envfile.Keys(first)
	keys2 := envfile.Keys(second)

	map1 := make(map[string]struct{}, len(keys1))
	for _, k := range keys1 {
		map1[k] = struct{}{}
	}

	map2 := make(map[string]struct{}, len(keys2))
	for _, k := range keys2 {
		map2[k] = struct{}{}
	}

	result := DiffResult{
		OnlyInFirst:  make([]string, 0, len(keys1)),
		OnlyInSecond: make([]string, 0, len(keys2)),
		InBoth:       make([]string, 0, len(keys1)/2),
	}

	for _, k := range keys1 {
		if _, ok := map2[k]; ok {
			result.InBoth = append(result.InBoth, k)
		} else {
			result.OnlyInFirst = append(result.OnlyInFirst, k)
		}
	}

	for _, k := range keys2 {
		if _, ok := map1[k]; !ok {
			result.OnlyInSecond = append(result.OnlyInSecond, k)
		}
	}

	return result
}
