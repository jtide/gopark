package api_test

import (
	"fmt"
	"github.com/jtide/gopark/api"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMinutesSinceMidnightFromTime(t *testing.T) {
	startTime, err := time.Parse(time.RFC3339, "2015-07-01T01:00:00Z")
	assert.NoError(t, err)

	startMinutes := api.MinutesSinceMidnightFromTime(startTime)
	assert.Equal(t, uint64(60), startMinutes)

	fmt.Printf("startTime %v: startMinutes: %v", startTime, startMinutes)
}
