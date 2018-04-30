package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	start := "2015-07-01T07:00:00Z"
	end := "2015-07-01T16:00:00Z"
	duration, err := ParseDuration(start, end)
	if err != nil {
		fmt.Printf(err.Error())
	}

	expectedHours := 9 * time.Hour
	assert.Equal(t, time.Duration(expectedHours), duration.Value)
}
