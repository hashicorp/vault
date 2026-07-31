// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHealthCheckExecutionResult_ToLogicalResponse verifies the typed result is
// flattened into the expected logical.Response data shape.
func TestHealthCheckExecutionResult_ToLogicalResponse(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	result := &HealthCheckExecutionResult{
		HealthChecks: []HealthCheck{
			{
				Type:          "connection",
				Healthy:       true,
				Reason:        "reachable",
				ReasonDetails: "ping ok",
				Timestamp:     ts,
				DurationMs:    42,
			},
		},
	}

	resp := result.ToLogicalResponse()
	require.NotNil(t, resp)

	checks, ok := resp.Data["health_checks"].([]map[string]interface{})
	require.True(t, ok, "health_checks should be a slice of maps")
	require.Len(t, checks, 1)

	assert.Equal(t, "connection", checks[0]["type"])
	assert.Equal(t, true, checks[0]["healthy"])
	assert.Equal(t, "reachable", checks[0]["reason"])
	assert.Equal(t, "ping ok", checks[0]["reason_details"])
	assert.Equal(t, ts, checks[0]["timestamp"])
	assert.Equal(t, 42, checks[0]["duration_ms"])
}

// TestHealthCheckExecutionResult_ToLogicalResponse_Nil verifies a nil receiver
// produces a nil response rather than panicking.
func TestHealthCheckExecutionResult_ToLogicalResponse_Nil(t *testing.T) {
	t.Parallel()

	var result *HealthCheckExecutionResult
	assert.Nil(t, result.ToLogicalResponse())
}

// TestHealthCheckExecutionResult_ToLogicalResponse_Empty verifies an empty set of
// checks still yields a valid, empty (non-nil) slice.
func TestHealthCheckExecutionResult_ToLogicalResponse_Empty(t *testing.T) {
	t.Parallel()

	result := &HealthCheckExecutionResult{}
	resp := result.ToLogicalResponse()
	require.NotNil(t, resp)

	checks, ok := resp.Data["health_checks"].([]map[string]interface{})
	require.True(t, ok)
	assert.Empty(t, checks)
}

// TestHealthCheckExecutionResultFromLogicalResponse_RoundTrip verifies that a
// result survives conversion to a logical.Response and back unchanged.
func TestHealthCheckExecutionResultFromLogicalResponse_RoundTrip(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	original := &HealthCheckExecutionResult{
		HealthChecks: []HealthCheck{
			{
				Type:          "connection",
				Healthy:       true,
				Reason:        "reachable",
				ReasonDetails: "ping ok",
				Timestamp:     ts,
				DurationMs:    42,
			},
			{
				Type:       "replication",
				Healthy:    false,
				Reason:     "lagging",
				Timestamp:  ts.Add(time.Second),
				DurationMs: 7,
			},
		},
	}

	decoded, err := HealthCheckExecutionResultFromLogicalResponse(original.ToLogicalResponse())
	require.NoError(t, err)
	require.NotNil(t, decoded)
	assert.Equal(t, original, decoded)
}

// TestHealthCheckExecutionResultFromLogicalResponse_TransportCoercion verifies the
// decoder absorbs the type changes introduced by the plugin transport: integers
// arriving as float64 and timestamps arriving as RFC3339 strings.
func TestHealthCheckExecutionResultFromLogicalResponse_TransportCoercion(t *testing.T) {
	t.Parallel()

	ts := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	resp := &Response{Data: map[string]interface{}{
		"health_checks": []map[string]interface{}{
			{
				"type":           "connection",
				"healthy":        true,
				"reason":         "reachable",
				"reason_details": "ping ok",
				"timestamp":      ts.Format(time.RFC3339),
				"duration_ms":    float64(42),
			},
		},
	}}

	decoded, err := HealthCheckExecutionResultFromLogicalResponse(resp)
	require.NoError(t, err)
	require.Len(t, decoded.HealthChecks, 1)

	check := decoded.HealthChecks[0]
	assert.Equal(t, "connection", check.Type)
	assert.True(t, check.Healthy)
	assert.Equal(t, ts, check.Timestamp)
	assert.Equal(t, 42, check.DurationMs)
}

// TestHealthCheckExecutionResultFromLogicalResponse_NilResponse verifies a nil
// response is rejected with an error.
func TestHealthCheckExecutionResultFromLogicalResponse_NilResponse(t *testing.T) {
	t.Parallel()

	decoded, err := HealthCheckExecutionResultFromLogicalResponse(nil)
	assert.Nil(t, decoded)
	assert.Error(t, err)
}

// TestHealthCheckExecutionResultFromLogicalResponse_UnknownField verifies strict
// decoding rejects fields the backend is not permitted to return.
func TestHealthCheckExecutionResultFromLogicalResponse_UnknownField(t *testing.T) {
	t.Parallel()

	resp := &Response{Data: map[string]interface{}{
		"health_checks":  []map[string]interface{}{},
		"unexpected_key": "value",
	}}

	decoded, err := HealthCheckExecutionResultFromLogicalResponse(resp)
	assert.Nil(t, decoded)
	assert.Error(t, err)
}

// TestHealthCheckExecutionResultFromLogicalResponse_UnknownNestedField verifies
// strict decoding also rejects unknown fields nested within an individual check.
func TestHealthCheckExecutionResultFromLogicalResponse_UnknownNestedField(t *testing.T) {
	t.Parallel()

	resp := &Response{Data: map[string]interface{}{
		"health_checks": []map[string]interface{}{
			{
				"type":    "connection",
				"healthy": true,
				"bogus":   "field",
			},
		},
	}}

	decoded, err := HealthCheckExecutionResultFromLogicalResponse(resp)
	assert.Nil(t, decoded)
	assert.Error(t, err)
}
