// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package audit

import "github.com/hashicorp/vault/internal/observability/event"

type backendEnt struct{}

func newBackendEnt(_ map[string]string) *backendEnt {
	return &backendEnt{}
}

func (b *backendEnt) IsFallback() bool {
	return false
}

// configureFilterNode is a no-op as filters are an Enterprise-only feature.
func (b *backend) configureFilterNode(_ string) error {
	return nil
}

func (b *backend) getMetricLabeler() event.Labeler {
	return &metricLabelerAuditSink{}
}
