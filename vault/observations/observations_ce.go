// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package observations

import (
	"context"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// ObservationSystem contains the main logic of running the observation system.
type ObservationSystem struct{}

type pluginObservationSystem struct{}

func (observations *ObservationSystem) Start() {}

func (observations *ObservationSystem) WithPlugin(_ *namespace.Namespace, _ *logical.ObservationPluginInfo) (*pluginObservationSystem, error) {
	return &pluginObservationSystem{}, nil
}

func (observations *pluginObservationSystem) RecordObservationFromPlugin(_ context.Context, _ string, _ map[string]interface{}) error {
	return nil
}

func (observations *ObservationSystem) RecordObservationToLedger(_ context.Context, _ string, _ *namespace.Namespace, _ map[string]interface{}) error {
	return nil
}

func NewObservationSystem(_ *NewObservationSystemConfig) (*ObservationSystem, error) {
	return &ObservationSystem{}, nil
}
