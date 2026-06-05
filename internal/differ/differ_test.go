// internal/differ/differ_test.go
// Tests for the differ engine.

package differ

import (
	"testing"

	"github.com/Vamshavardhan50/envguard/pkg/envfile"
	"github.com/stretchr/testify/assert"
)

func TestEngineDiff(t *testing.T) {
	first := []envfile.EnvVar{
		{Key: "COMMON"},
		{Key: "FIRST_ONLY"},
	}

	second := []envfile.EnvVar{
		{Key: "COMMON"},
		{Key: "SECOND_ONLY"},
	}

	engine := &Engine{}
	res := engine.Diff(first, second)

	assert.Equal(t, []string{"FIRST_ONLY"}, res.OnlyInFirst)
	assert.Equal(t, []string{"SECOND_ONLY"}, res.OnlyInSecond)
	assert.Equal(t, []string{"COMMON"}, res.InBoth)
}
