// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/internal/observability/event"
)

var _ event.Labeler = (*metricLabelerAuditSink)(nil)

var (
	metricLabelAuditSinkSuccess = []string{"audit", "sink", "success"}
	metricLabelAuditSinkFailure = []string{"audit", "sink", "failure"}
)

// metricLabelerAuditSink can be used to provide labels for the success or failure
// of a sink node used for a normal audit device.
type metricLabelerAuditSink struct{}

// Labels provides the success and failure labels for an audit sink, based on the error supplied.
// Success: 'vault.audit.sink.success'
// Failure: 'vault.audit.sink.failure'
func (m metricLabelerAuditSink) Labels(_ *eventlogger.Event, err error) []string {
	if err != nil {
		// NOTE: a cancelled context would still result in an error.
		return metricLabelAuditSinkFailure
	}

	return metricLabelAuditSinkSuccess
}
