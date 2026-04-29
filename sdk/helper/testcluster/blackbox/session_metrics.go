// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"encoding/json"
	"fmt"
	"time"
)

// MetricsResponse represents the response from sys/metrics endpoint
type MetricsResponse struct {
	Data struct {
		Gauges []struct {
			Name   string            `json:"Name"`
			Value  float64           `json:"Value"`
			Labels map[string]string `json:"Labels"`
		} `json:"Gauges"`
		Counters []struct {
			Name   string            `json:"Name"`
			Count  int               `json:"Count"`
			Sum    float64           `json:"Sum"`
			Labels map[string]string `json:"Labels"`
		} `json:"Counters"`
		Samples []struct {
			Name   string            `json:"Name"`
			Count  int               `json:"Count"`
			Sum    float64           `json:"Sum"`
			Labels map[string]string `json:"Labels"`
		} `json:"Samples"`
	} `json:"data"`
}

// AssertMetricGaugeValue verifies that a specific gauge metric has the expected value
// This method includes retry logic with configurable timeout
// Note: retryInterval parameter is ignored as the SDK uses a fixed 200ms interval
func (s *Session) AssertMetricGaugeValue(gaugeName string, expectedValue float64, timeout time.Duration, retryInterval time.Duration) {
	s.t.Helper()

	s.EventuallyWithTimeout(func() error {
		// Read sys/metrics endpoint
		secret, err := s.Client.Logical().Read("sys/metrics")
		if err != nil {
			return fmt.Errorf("failed to read sys/metrics: %w", err)
		}

		if secret == nil || secret.Data == nil {
			return fmt.Errorf("sys/metrics returned nil data")
		}

		// Marshal and unmarshal to get proper structure
		dataBytes, err := json.Marshal(secret.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal metrics data: %w", err)
		}

		var metricsData struct {
			Gauges []struct {
				Name   string            `json:"Name"`
				Value  float64           `json:"Value"`
				Labels map[string]string `json:"Labels"`
			} `json:"Gauges"`
		}

		if err := json.Unmarshal(dataBytes, &metricsData); err != nil {
			return fmt.Errorf("failed to unmarshal metrics data: %w", err)
		}

		// Find the gauge by name
		var found bool
		var actualValue float64
		for _, gauge := range metricsData.Gauges {
			if gauge.Name == gaugeName {
				found = true
				actualValue = gauge.Value
				break
			}
		}

		if !found {
			return fmt.Errorf("gauge metric %q not found in sys/metrics response", gaugeName)
		}

		if actualValue != expectedValue {
			return fmt.Errorf("gauge %q has value %.0f, expected %.0f", gaugeName, actualValue, expectedValue)
		}

		s.t.Logf("Gauge metric %q has expected value: %.0f", gaugeName, expectedValue)
		return nil
	}, timeout)
}
