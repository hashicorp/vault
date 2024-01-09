// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"strings"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/internal/observability/event"
)

const (
	metricLabelAuditSinkSuccess         = "audit.sink.success"
	metricLabelAuditSinkFailure         = "audit.sink.failure"
	metricLabelAuditSinkFallbackSuccess = "audit.fallback.success"
	metricLabelAuditSinkFallbackMiss    = "audit.fallback.miss"
	metricLabelSeparator                = "."
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

// Labels provides the success and failure labels for an audit sink, based on the error supplied.
// Success: 'vault.audit.sink.success'
// Failure: 'vault.audit.sink.failure'
func (m MetricLabelerAuditSink) Labels(_ *eventlogger.Event, err error) []string {
	if err != nil {
		return splitLabel(metricLabelAuditSinkFailure)
	}

	return splitLabel(metricLabelAuditSinkSuccess)
}

// Labels provides the success and failures labels for an audit fallback sink, based on the error supplied.
// Success: 'vault.audit.fallback.success'
// Failure: 'vault.audit.sink.failure'
func (m MetricLabelerAuditFallback) Labels(_ *eventlogger.Event, err error) []string {
	if err != nil {
		return splitLabel(metricLabelAuditSinkFailure)
	}

	return splitLabel(metricLabelAuditSinkFallbackSuccess)
}

// MetricLabelsFallbackMiss returns the labels which indicate an audit entry was missed.
// 'vault.audit.fallback.miss'
func MetricLabelsFallbackMiss() []string {
	return splitLabel(metricLabelAuditSinkFallbackMiss)
}

// splitLabel takes a label and splits it using the metricLabelSeparator.
func splitLabel(metricLabel string) []string {
	return strings.Split(metricLabel, metricLabelSeparator)
}
