package vault

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsTimeNowBetween(t *testing.T) {
	testcases := []struct {
		name        string
		startTime   time.Time
		endTime     time.Time
		expectation bool
	}{
		{
			name:        "is between start and end times",
			startTime:   time.Now().Add(-1 * time.Hour),
			endTime:     time.Now().Add(time.Hour),
			expectation: true,
		},
		{
			name:        "both start and end times before",
			startTime:   time.Now().Add(-2 * time.Hour),
			endTime:     time.Now().Add(-1 * time.Hour),
			expectation: false,
		},
		{
			name:        "both start and end times after",
			startTime:   time.Now().Add(time.Hour),
			endTime:     time.Now().Add(2 * time.Hour),
			expectation: false,
		},
		{
			name:        "is between start and end times, reversed",
			startTime:   time.Now().Add(time.Hour),
			endTime:     time.Now().Add(-1 * time.Hour),
			expectation: false,
		},
		{
			name:        "both start and end times before, reversed",
			startTime:   time.Now().Add(-1 * time.Hour),
			endTime:     time.Now().Add(-2 * time.Hour),
			expectation: false,
		},
		{
			name:        "both start and end times after, reversed",
			startTime:   time.Now().Add(2 * time.Hour),
			endTime:     time.Now().Add(time.Hour),
			expectation: false,
		},
	}

	for _, testcase := range testcases {
		result := isTimeNowBetween(testcase.startTime, testcase.endTime)
		assert.Equal(t, testcase.expectation, result, testcase.name)
	}
}
