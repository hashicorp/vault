// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package observe

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

type PkiCeObserver struct {
	logger   hclog.Logger
	observer PluginObserve
}

var _ PkiObserver = (*PkiCeObserver)(nil)

func (p PkiCeObserver) RecordPKIObservation(_ context.Context, _ *logical.Request, _ string, _ ...AdditionalPKIMetadata) {
	// No-op for Community Edition
}

func NewPkiCeObserver(logger hclog.Logger, observer PluginObserve) *PkiCeObserver {
	return &PkiCeObserver{
		logger:   logger,
		observer: observer,
	}
}
