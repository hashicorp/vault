// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestGeneratePkiBillingMetric tests the PKI billing metric generation
func TestGeneratePkiBillingMetric(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := context.Background()
	backend := core.systemBackend

	currentMonth := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	t.Run("returns zero count when no data exists", func(t *testing.T) {
		overview, err := backend.buildPkiBillingMetric(ctx, currentMonth)
		require.NoError(t, err)
		require.NotNil(t, overview)

		// Check metric_name
		require.Equal(t, "pki_units", overview["metric_name"])

		// Check metric_data structure
		metricData, ok := overview["metric_data"].(map[string]interface{})
		require.True(t, ok)

		// Check total
		require.Equal(t, float64(0), metricData["total"])
	})

	t.Run("returns stored count when data exists", func(t *testing.T) {
		month := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
		expectedCount := 100.1234

		err := core.UpdatePkiDurationAdjustedCount(ctx, expectedCount, month)
		require.NoError(t, err)

		// Generate overview
		overview, err := backend.buildPkiBillingMetric(ctx, month)
		require.NoError(t, err)
		require.NotNil(t, overview)

		// Check metric_name
		require.Equal(t, "pki_units", overview["metric_name"])

		// Check metric_data structure
		metricData, ok := overview["metric_data"].(map[string]interface{})
		require.True(t, ok)

		// Check total
		require.Equal(t, expectedCount, metricData["total"])
	})

	t.Run("handles different months independently", func(t *testing.T) {
		month1 := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
		month2 := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

		// Store different counts for different months
		count1 := 50.5
		count2 := 75.25

		err := core.UpdatePkiDurationAdjustedCount(ctx, count1, month1)
		require.NoError(t, err)
		err = core.UpdatePkiDurationAdjustedCount(ctx, count2, month2)
		require.NoError(t, err)

		// Generate overview for month1
		overview1, err := backend.buildPkiBillingMetric(ctx, month1)
		require.NoError(t, err)
		metricData1 := overview1["metric_data"].(map[string]interface{})
		require.Equal(t, count1, metricData1["total"])

		// Generate overview for month2
		overview2, err := backend.buildPkiBillingMetric(ctx, month2)
		require.NoError(t, err)
		metricData2 := overview2["metric_data"].(map[string]interface{})
		require.Equal(t, count2, metricData2["total"])
	})

	t.Run("uses constant for metric name", func(t *testing.T) {
		month := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)

		overview, err := backend.buildPkiBillingMetric(ctx, month)
		require.NoError(t, err)

		// Verify it uses the right metric name
		require.Equal(t, "pki_units", overview["metric_name"])
	})
}

// Made with Bob
