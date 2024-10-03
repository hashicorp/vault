// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package audit

import (
	"context"
	"fmt"
)

// brokerEnt provides extensions to the broker behavior, but not in the community edition.
type brokerEnt struct{}

func newBrokerEnt() (*brokerEnt, error) {
	return &brokerEnt{}, nil
}

func (b *Broker) validateRegistrationRequest(_ Backend) error {
	return nil
}

func (b *Broker) handlePipelineRegistration(backend Backend) error {
	err := b.register(backend)
	if err != nil {
		return fmt.Errorf("unable to register device for %q: %w", backend.Name(), err)
	}

	return nil
}

func (b *Broker) handlePipelineDeregistration(ctx context.Context, name string) error {
	return b.deregister(ctx, name)
}

// requiredSuccessThresholdSinks is the value that should be used as the success
// threshold in the eventlogger broker.
func (b *Broker) requiredSuccessThresholdSinks() int {
	if len(b.backends) > 0 {
		return 1
	}

	return 0
}

func (b *brokerEnt) handleAdditionalAudit(_ context.Context, _ *Event) error {
	return nil
}
