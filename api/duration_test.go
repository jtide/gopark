package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	start := "2015-07-01T07:00:00Z"
	end := "2015-07-01T16:00:00Z"
	expectedHours := 9 * time.Hour

	duration, err := ParseDuration(start, end)
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(expectedHours), duration.Value)
}
