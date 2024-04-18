// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/internal/observability/event"
)

var (
	_ event.Labeler = (*metricLabelerAuditSink)(nil)
	_ event.Labeler = (*metricLabelerAuditFallback)(nil)
)

var (
	metricLabelAuditSinkSuccess     = []string{"audit", "sink", "success"}
	metricLabelAuditSinkFailure     = []string{"audit", "sink", "failure"}
	metricLabelAuditFallbackSuccess = []string{"audit", "fallback", "success"}
	metricLabelAuditFallbackMiss    = []string{"audit", "fallback", "miss"}
)

// metricLabelerAuditSink can be used to provide labels for the success or failure
// of a sink node used for a normal audit device.
type metricLabelerAuditSink struct{}

// metricLabelerAuditFallback can be used to provide labels for the success or failure
// of a sink node used for an audit fallback device.
type metricLabelerAuditFallback struct{}

// Labels provides the success and failure labels for an audit sink, based on the error supplied.
// Success: 'vault.audit.sink.success'
// Failure: 'vault.audit.sink.failure'
func (m metricLabelerAuditSink) Labels(_ *eventlogger.Event, err error) []string {
	if err != nil {
		return metricLabelAuditSinkFailure
	}

	return metricLabelAuditSinkSuccess
}

// Labels provides the success and failures labels for an audit fallback sink, based on the error supplied.
// Success: 'vault.audit.fallback.success'
// Failure: 'vault.audit.sink.failure'
func (m metricLabelerAuditFallback) Labels(_ *eventlogger.Event, err error) []string {
	if err != nil {
		return metricLabelAuditSinkFailure
	}

	return metricLabelAuditFallbackSuccess
}

// MetricLabelsFallbackMiss returns the labels which indicate an audit entry was missed.
// 'vault.audit.fallback.miss'
func MetricLabelsFallbackMiss() []string {
	return metricLabelAuditFallbackMiss
}
