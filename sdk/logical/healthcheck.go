// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

// HealthCheck represents a single health check performed by a plugin backend.
type HealthCheck struct {
	Type          string    `json:"type" mapstructure:"type"`
	Healthy       bool      `json:"healthy" mapstructure:"healthy"`
	Reason        string    `json:"reason" mapstructure:"reason"`
	ReasonDetails string    `json:"reason_details" mapstructure:"reason_details"`
	Timestamp     time.Time `json:"timestamp" mapstructure:"timestamp"`
	DurationMs    int       `json:"duration_ms" mapstructure:"duration_ms"`
}

// HealthCheckExecutionResult is the result of executing a plugin backend's health
// checks. A plugin reports one or more individual checks along with optional
// custom metadata describing the resource that was checked.
type HealthCheckExecutionResult struct {
	HealthChecks []HealthCheck `json:"health_checks" mapstructure:"health_checks"`
}

// ToLogicalResponse converts the typed response into a logical.Response for transport.
func (r *HealthCheckExecutionResult) ToLogicalResponse() *Response {
	if r == nil {
		return nil
	}

	checks := make([]map[string]interface{}, len(r.HealthChecks))
	for i, c := range r.HealthChecks {
		checks[i] = map[string]interface{}{
			"type":           c.Type,
			"healthy":        c.Healthy,
			"reason":         c.Reason,
			"reason_details": c.ReasonDetails,
			"timestamp":      c.Timestamp,
			"duration_ms":    c.DurationMs,
		}
	}

	return &Response{Data: map[string]interface{}{
		"health_checks": checks,
	}}
}

// HealthCheckExecutionResultFromLogicalResponse converts a logical.Response into a
// typed execution result. Decoding is strict so backends cannot return arbitrary
// fields, while weak typing and a time hook absorb the type changes introduced by
// the plugin transport (e.g. integers arriving as float64, timestamps as strings).
func HealthCheckExecutionResultFromLogicalResponse(resp *Response) (*HealthCheckExecutionResult, error) {
	if resp == nil {
		return nil, fmt.Errorf("nil health check response")
	}

	result := &HealthCheckExecutionResult{}

	// mapstructure.NewDecoder only returns an error when Result is nil or not a
	// pointer. Result is always a valid pointer here, so the error check below is
	// a guard against a programming error rather than a runtime condition.
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeHookFunc(time.RFC3339),
		ErrorUnused:      true,
		WeaklyTypedInput: true,
		Result:           result,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create health check response decoder: %w", err)
	}

	if err := decoder.Decode(resp.Data); err != nil {
		return nil, fmt.Errorf("decode health check response: %w", err)
	}

	return result, nil
}
