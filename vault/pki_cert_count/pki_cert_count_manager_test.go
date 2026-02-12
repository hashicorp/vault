// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki_cert_count

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestPkiCertificateCountManager_IncrementAndConsume tests the behaviour of
// PkiCertificateCountManager.
func TestPkiCertificateCountManager_IncrementAndConsume(t *testing.T) {
	manager := newPkiCertificateCountManager(hclog.NewNullLogger())
	consumerJobInterval = 10 * time.Millisecond

	firstConsumerTotalCount := &atomic.Uint64{}
	manager.StartConsumerJob(func(inc logical.CertCount) {
		firstConsumerTotalCount.Add(inc.IssuedCerts + inc.StoredCerts)
	})

	jobCountLock := sync.Mutex{}
	jobCount := logical.CertCount{}
	consumer := func(inc logical.CertCount) {
		jobCountLock.Lock()
		defer jobCountLock.Unlock()
		jobCount.Add(inc)
	}

	manager.StartConsumerJob(consumer)
	// StartConsumerJob calls StopConsumerJob, which will make one last call to the consumer,
	// so lets wait a bit not to lose any increments.
	time.Sleep(20 * time.Millisecond)
	firstConsumerTotalCount.Store(0)

	manager.AddCount(logical.CertCount{IssuedCerts: 3, StoredCerts: 0})
	manager.AddCount(logical.CertCount{IssuedCerts: 0, StoredCerts: 5})
	manager.Increment().AddIssuedCertificate(true)
	manager.Increment().AddIssuedCertificate(false)

	time.Sleep(100 * time.Millisecond)

	// Calling stop again should not panic.
	manager.StopConsumerJob()

	jobCountLock.Lock()
	defer jobCountLock.Unlock()

	require.Equal(t, uint64(5), jobCount.IssuedCerts, "issued count mismatch")
	require.Equal(t, uint64(6), jobCount.StoredCerts, "stored count mismatch")
	require.Zero(t, firstConsumerTotalCount.Load(), "first consumer should not have been called")
}
