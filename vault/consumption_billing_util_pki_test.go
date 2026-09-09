// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestGetStoredPkiDurationAdjustedCount tests reading PKI duration-adjusted counts from storage
func TestGetStoredPkiDurationAdjustedCount(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := context.Background()
	currentMonth := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	t.Run("returns zero when no data exists", func(t *testing.T) {
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, currentMonth)
		require.NoError(t, err)
		require.Equal(t, float64(0), count)
	})

	t.Run("returns stored count when data exists", func(t *testing.T) {
		// Store a count first
		expectedCount := 42.5
		err := core.UpdatePkiDurationAdjustedCount(ctx, expectedCount, currentMonth)
		require.NoError(t, err)

		// Retrieve and verify
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, currentMonth)
		require.NoError(t, err)
		require.Equal(t, expectedCount, count)
	})

	t.Run("returns different counts for different months", func(t *testing.T) {
		month1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		month2 := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

		// Store different counts for different months
		count1 := 10.5
		count2 := 20.5
		err := core.UpdatePkiDurationAdjustedCount(ctx, count1, month1)
		require.NoError(t, err)
		err = core.UpdatePkiDurationAdjustedCount(ctx, count2, month2)
		require.NoError(t, err)

		// Verify each month has its own count
		retrievedCount1, err := core.GetStoredPkiDurationAdjustedCount(ctx, month1)
		require.NoError(t, err)
		require.Equal(t, count1, retrievedCount1)

		retrievedCount2, err := core.GetStoredPkiDurationAdjustedCount(ctx, month2)
		require.NoError(t, err)
		require.Equal(t, count2, retrievedCount2)
	})

	t.Run("returns zero for non-existent month", func(t *testing.T) {
		futureMonth := time.Date(2027, 12, 1, 0, 0, 0, 0, time.UTC)
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, futureMonth)
		require.NoError(t, err)
		require.Equal(t, float64(0), count)
	})
}

// TestUpdatePkiDurationAdjustedCount tests storing and incrementing PKI duration-adjusted counts
func TestUpdatePkiDurationAdjustedCount(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := context.Background()
	currentMonth := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	t.Run("stores initial count", func(t *testing.T) {
		initialCount := 15.5
		err := core.UpdatePkiDurationAdjustedCount(ctx, initialCount, currentMonth)
		require.NoError(t, err)

		// Verify the count was stored
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, currentMonth)
		require.NoError(t, err)
		require.Equal(t, initialCount, count)
	})

	t.Run("increments existing count", func(t *testing.T) {
		month := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)

		// Store initial count
		initialCount := 10.0
		err := core.UpdatePkiDurationAdjustedCount(ctx, initialCount, month)
		require.NoError(t, err)

		// Increment the count
		increment := 5.5
		err = core.UpdatePkiDurationAdjustedCount(ctx, increment, month)
		require.NoError(t, err)

		// Verify the count was incremented
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, month)
		require.NoError(t, err)
		require.Equal(t, initialCount+increment, count)
	})

	t.Run("handles multiple increments", func(t *testing.T) {
		// Start with zero
		month := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
		increments := []float64{1.5, 2.5, 3.0, 4.5}
		expectedTotal := 0.0

		for _, inc := range increments {
			err := core.UpdatePkiDurationAdjustedCount(ctx, inc, month)
			require.NoError(t, err)
			expectedTotal += inc
		}

		// Verify the total
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, month)
		require.NoError(t, err)
		require.Equal(t, expectedTotal, count)
	})

	t.Run("rejects negative increments", func(t *testing.T) {
		month := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

		// Store initial count
		initialCount := 100.0
		err := core.UpdatePkiDurationAdjustedCount(ctx, initialCount, month)
		require.NoError(t, err)

		// Attempt to apply negative increment
		decrement := -25.5
		err = core.UpdatePkiDurationAdjustedCount(ctx, decrement, month)
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be non-negative")

		// Verify the count was not changed
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, month)
		require.NoError(t, err)
		require.Equal(t, initialCount, count)
	})

	t.Run("handles zero increment", func(t *testing.T) {
		month := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
		initialCount := 30.0

		// Store initial count
		err := core.UpdatePkiDurationAdjustedCount(ctx, initialCount, month)
		require.NoError(t, err)

		// Apply zero increment
		err = core.UpdatePkiDurationAdjustedCount(ctx, 0.0, month)
		require.NoError(t, err)

		// Verify the count remains unchanged
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, month)
		require.NoError(t, err)
		require.Equal(t, initialCount, count)
	})

	t.Run("handles fractional increments accurately", func(t *testing.T) {
		month := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

		// Store fractional counts
		err := core.UpdatePkiDurationAdjustedCount(ctx, 0.1, month)
		require.NoError(t, err)
		err = core.UpdatePkiDurationAdjustedCount(ctx, 0.2, month)
		require.NoError(t, err)
		err = core.UpdatePkiDurationAdjustedCount(ctx, 0.3, month)
		require.NoError(t, err)

		// Verify the total
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, month)
		require.NoError(t, err)
		require.InDelta(t, 0.6, count, 0.0001) // Use InDelta for floating point comparison
	})

	t.Run("updates count through public method", func(t *testing.T) {
		month := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
		increment := 25.5
		err := core.UpdatePkiDurationAdjustedCount(ctx, increment, month)
		require.NoError(t, err)

		// Verify the count was stored
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, month)
		require.NoError(t, err)
		require.Equal(t, increment, count)
	})

	t.Run("handles minimum increment of 0.0001", func(t *testing.T) {
		month := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)

		err := core.UpdatePkiDurationAdjustedCount(ctx, 1.0, month)
		require.NoError(t, err)

		// Apply multiple minimum increment
		minIncrement := 0.0001

		err = core.UpdatePkiDurationAdjustedCount(ctx, minIncrement, month)
		require.NoError(t, err)

		// Verify the total count
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, month)
		require.NoError(t, err)
		expectedTotal := 1.0001
		require.InDelta(t, expectedTotal, count, 0.00001) // Use InDelta for floating point comparison
	})
}

// TestConcurrentPkiDurationAdjustedCount tests concurrent updates
func TestConcurrentPkiDurationAdjustedCount(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := context.Background()

	t.Run("handles concurrent updates", func(t *testing.T) {
		month := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
		numGoroutines := 10
		incrementPerGoroutine := 1.0

		// Launch concurrent updates
		done := make(chan bool, numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				err := core.UpdatePkiDurationAdjustedCount(ctx, incrementPerGoroutine, month)
				require.NoError(t, err)
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Verify the total count
		count, err := core.GetStoredPkiDurationAdjustedCount(ctx, month)
		require.NoError(t, err)
		require.Equal(t, float64(numGoroutines)*incrementPerGoroutine, count)
	})
}

// Made with Bob
