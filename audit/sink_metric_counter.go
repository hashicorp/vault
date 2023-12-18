// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/internal/observability/event"
)

const (
	auditSinkLabelSuccess     = "vault.audit.sink.success"
	auditSinkLabelFailure     = "vault.audit.sink.failure"
	auditFallbackLabelSuccess = "vault.audit.fallback.success"
	auditFallbackLabelMiss    = "vault.audit.fallback.miss"
)

var _ event.Labeler = (*MetricCounterAuditSink)(nil)

// MetricCounterAuditSink can be used to provide labels for the successor failure
// of a sink node for audit.
type MetricCounterAuditSink struct{}
type MetricCounterAuditFallback struct{}

// Label provides the success and failures labels based on the error supplied.
// Success: 'vault.audit.sink.success'
// Failure: 'vault.audit.sink.failure'
func (m MetricCounterAuditSink) Label(_ *eventlogger.Event, err error) string {
	if err != nil {
		return auditSinkLabelFailure
	}

	return auditSinkLabelSuccess
}

// Label provides the success and failures labels based on the error supplied.
// Success: 'vault.audit.fallback.success'
// Failure: 'vault.audit.sink.failure'
func (m MetricCounterAuditFallback) Label(_ *eventlogger.Event, err error) string {
	if err != nil {
		return auditFallbackLabelMiss
	}

	return auditFallbackLabelSuccess
}
