// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki_cert_count

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

// TestPkiCertificateCountManager_IncrementAndConsume tests the behaviour of
// PkiCertificateCountManager.
func TestPkiCertificateCountManager_IncrementAndConsume(t *testing.T) {
	manager := newPkiCertificateCountManager(hclog.NewNullLogger())
	consumerJobInterval = 10 * time.Millisecond

	firstConsumerTotalCount := &atomic.Uint64{}
	manager.StartConsumerJob(func(i, s uint64) {
		firstConsumerTotalCount.Add(i + s)
	})

	issued := &atomic.Uint64{}
	stored := &atomic.Uint64{}

	consumer := func(i, s uint64) {
		issued.Add(i)
		stored.Add(s)
	}

	manager.StartConsumerJob(consumer)
	// StartConsumerJob calls StopConsumerJob, which will make one last call to the consumer,
	// so lets wait a bit not to lose any increments.
	time.Sleep(20 * time.Millisecond)
	firstConsumerTotalCount.Store(0)

	manager.IncrementCount(3, 0)
	manager.IncrementCount(0, 5)
	manager.AddIssuedCertificate(true)
	manager.AddIssuedCertificate(false)

	time.Sleep(100 * time.Millisecond)

	// Calling stop again should not panic.
	manager.StopConsumerJob()

	require.Equal(t, uint64(5), issued.Load(), "issued count mismatch")
	require.Equal(t, uint64(6), stored.Load(), "stored count mismatch")
	require.Zero(t, firstConsumerTotalCount.Load(), "first consumer should not have been called")
}
