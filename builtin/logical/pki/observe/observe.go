// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package observe

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
)

type PluginObserve interface {
	// RecordObservation is used to record observations through the plugin's observation system.
	// It returns ErrNoObservations if the observation system has not been configured or enabled.
	RecordObservation(ctx context.Context, observationType string, data map[string]interface{}) error
}

type PkiObserver interface {
	RecordPKIObservation(ctx context.Context, req *logical.Request, observationType string, additionalMetadata ...AdditionalPKIMetadata)
}

type AdditionalPKIMetadata struct {
	key   string
	value any
}

func NewAdditionalPKIMetadata(key string, value any) AdditionalPKIMetadata {
	return AdditionalPKIMetadata{key: key, value: value}
}
