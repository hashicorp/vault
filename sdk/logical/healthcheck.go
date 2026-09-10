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

// Reason strings for the aggregate health check summary produced after
// executing a backend's health checks.
const (
	HealthCheckReasonNoChecks       = "No checks reported"
	HealthCheckReasonAllPassed      = "All checks passed"
	HealthCheckReasonAllFailed      = "All checks failed"
	HealthCheckReasonPartialFailure = "%d out of %d checks failed"
)

// Check types for the Type field of a health check, naming what was verified.
// They are shared across plugins so operators and automation see the same names
// regardless of which backend reported the check.
const (
	// HealthCheckTypeRootConfigConnectivity verifies that the configured root
	// credentials can reach the backing service.
	HealthCheckTypeRootConfigConnectivity = "root_config_connectivity"

	// HealthCheckTypeStaticRoleCredentials verifies that the stored credentials
	// for a static role are still valid.
	HealthCheckTypeStaticRoleCredentials = "static_role_credentials"
)

// Tokens for the ReasonDetails field of a failed health check. They give
// callers a stable, machine-readable cause to branch on, since the accompanying
// Reason string is free-form prose meant for humans.
const (
	// HealthCheckReasonDetailsConfigReadError indicates the storage read of
	// the plugin config failed.
	HealthCheckReasonDetailsConfigReadError = "config_read_error"

	// HealthCheckReasonDetailsConfigNotFound indicates no config exists in
	// storage yet.
	HealthCheckReasonDetailsConfigNotFound = "config_not_found"

	// HealthCheckReasonDetailsRoleReadError indicates the storage read of a
	// static role failed.
	HealthCheckReasonDetailsRoleReadError = "role_read_error"

	// HealthCheckReasonDetailsRoleNotFound indicates the named static role
	// does not exist in storage.
	HealthCheckReasonDetailsRoleNotFound = "role_not_found"

	// HealthCheckReasonDetailsAccountLocked indicates the backing service
	// rejected the credentials because the account is locked.
	HealthCheckReasonDetailsAccountLocked = "account_locked"

	// HealthCheckReasonDetailsBindError indicates a generic authentication or
	// connectivity failure against the backing service.
	HealthCheckReasonDetailsBindError = "bind_error"
)

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

// HealthCheckExecutionResultBuilder accumulates one or more health checks and
// emits a HealthCheckExecutionResult via Build().
//
// Each individual check is created via NewCheck, which returns a
// HealthCheckHandle with its type bound and its timer already started.
// Calling Fail or Pass on the handle marks the outcome and returns the parent
// builder so the call site can either chain .Build() for an early exit or
// continue adding more checks:
//
//	hb := NewHealthCheckExecutionResultBuilder()
//
//	check := hb.NewCheck(HealthCheckTypeRootConfigConnectivity)
//	if config == nil {
//	    return check.Fail("no config found", HealthCheckReasonDetailsConfigNotFound).Build(), nil
//	}
//	return check.Pass("connected to " + url).Build(), nil
//
// Multi-check example:
//
//	hb := NewHealthCheckExecutionResultBuilder()
//
//	c1 := hb.NewCheck("check_one")
//	if err != nil {
//	    c1.Fail("check one failed", err.Error())  // mark failed, continue
//	} else {
//	    c1.Pass("check one ok")
//	}
//
//	c2 := hb.NewCheck("check_two")
//	if err2 != nil {
//	    return c2.Fail("check two failed", err2.Error()).Build(), nil
//	}
//	c2.Pass("check two ok")
//
//	return hb.Build(), nil
//
// A builder is not safe for concurrent use. Checks that run in parallel should
// be collected by the caller and recorded from a single goroutine.
type HealthCheckExecutionResultBuilder struct {
	checks []HealthCheck
}

// NewHealthCheckExecutionResultBuilder returns an empty builder ready to
// accumulate checks.
func NewHealthCheckExecutionResultBuilder() *HealthCheckExecutionResultBuilder {
	return &HealthCheckExecutionResultBuilder{}
}

// NewCheck creates a HealthCheckHandle with checkType bound and its start
// timer running. The handle's Fail/Pass methods append into this builder.
func (hb *HealthCheckExecutionResultBuilder) NewCheck(checkType string) *HealthCheckHandle {
	return &HealthCheckHandle{
		parent:    hb,
		checkType: checkType,
		// Deliberately not .UTC() here: that would strip the monotonic clock
		// reading time.Since needs, leaving DurationMs computed from the wall
		// clock and vulnerable to NTP steps. Fail/Pass convert on emit.
		start: time.Now(),
	}
}

// Build emits all accumulated checks as a HealthCheckExecutionResult.
func (hb *HealthCheckExecutionResultBuilder) Build() *HealthCheckExecutionResult {
	return &HealthCheckExecutionResult{HealthChecks: hb.checks}
}

func (hb *HealthCheckExecutionResultBuilder) append(c HealthCheck) {
	hb.checks = append(hb.checks, c)
}

// HealthCheckHandle is a scoped handle for a single named check. Its type and
// start time are fixed at creation. Calling Fail or Pass finalises the check,
// appends it to the parent builder, and returns the parent so the caller can
// chain .Build() or continue adding more checks.
//
// A handle records at most one result: the first Fail or Pass wins and any
// later call on the same handle is a no-op, so a check cannot be counted twice
// in the aggregate summary.
type HealthCheckHandle struct {
	parent    *HealthCheckExecutionResultBuilder
	checkType string
	start     time.Time
	done      bool
}

// Fail marks this check as failed and appends it to the parent builder.
// It returns the parent builder so the caller can chain .Build() or continue.
// It is a no-op if this handle has already been finalised.
func (h *HealthCheckHandle) Fail(reason, details string) *HealthCheckExecutionResultBuilder {
	if h.done {
		return h.parent
	}
	h.done = true

	h.parent.append(HealthCheck{
		Type:          h.checkType,
		Healthy:       false,
		Reason:        reason,
		ReasonDetails: details,
		Timestamp:     h.start.UTC(),
		DurationMs:    int(time.Since(h.start).Milliseconds()),
	})
	return h.parent
}

// Pass marks this check as passed and appends it to the parent builder.
// It returns the parent builder so the caller can chain .Build() or continue.
// It is a no-op if this handle has already been finalised.
func (h *HealthCheckHandle) Pass(reason string) *HealthCheckExecutionResultBuilder {
	if h.done {
		return h.parent
	}
	h.done = true

	h.parent.append(HealthCheck{
		Type:       h.checkType,
		Healthy:    true,
		Reason:     reason,
		Timestamp:  h.start.UTC(),
		DurationMs: int(time.Since(h.start).Milliseconds()),
	})
	return h.parent
}
