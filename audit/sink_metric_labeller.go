// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/internal/observability/event"
)

const (
	MetricLabelAuditSinkSuccess         = "vault.audit.sink.success"
	MetricLabelAuditSinkFailure         = "vault.audit.sink.failure"
	MetricLabelAuditSinkFallbackSuccess = "vault.audit.fallback.success"
	MetricLabelAuditSinkFallbackMiss    = "vault.audit.fallback.miss"
)

var (
	_ event.Labeler = (*MetricLabelerAuditSink)(nil)
	_ event.Labeler = (*MetricLabelerAuditFallback)(nil)
)

// MetricLabelerAuditSink can be used to provide labels for the success or failure
// of a sink node used for a normal audit device.
type MetricLabelerAuditSink struct{}

// MetricLabelerAuditFallback can be used to provide labels for the success or failure
// of a sink node used for an audit fallback device.
type MetricLabelerAuditFallback struct{}

// Label provides the success and failure labels for an audit sink, based on the error supplied.
// Success: 'vault.audit.sink.success'
// Failure: 'vault.audit.sink.failure'
func (m MetricLabelerAuditSink) Label(_ *eventlogger.Event, err error) string {
	if err != nil {
		return MetricLabelAuditSinkFailure
	}

	return MetricLabelAuditSinkSuccess
}

// Label provides the success and failures labels for an audit fallback sink, based on the error supplied.
// Success: 'vault.audit.fallback.success'
// Failure: 'vault.audit.sink.failure'
func (m MetricLabelerAuditFallback) Label(_ *eventlogger.Event, err error) string {
	if err != nil {
		return MetricLabelAuditSinkFailure
	}

	return MetricLabelAuditSinkFallbackSuccess
}
