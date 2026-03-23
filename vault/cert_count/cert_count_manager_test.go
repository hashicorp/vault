// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cert_count

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// createTestCertificate creates a test certificate with the specified validity duration
func createTestCertificate(t *testing.T, validity time.Duration) *x509.Certificate {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	notBefore := time.Now()
	notAfter := notBefore.Add(validity)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "test-cert",
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(certBytes)
	require.NoError(t, err)

	return cert
}

// TestCertificateCountManager_IncrementAndConsume tests the behaviour of
// CertificateCountManager.
func TestCertificateCountManager_IncrementAndConsume(t *testing.T) {
	manager := newCertificateCountManager(hclog.NewNullLogger())
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

	// Create test certificates with different validity periods
	// 730 hours = 1 month = 1.0 billable unit
	cert1Month := createTestCertificate(t, 730*time.Hour)
	// 8760 hours = 1 year = 12.0 billable units
	cert1Year := createTestCertificate(t, 8760*time.Hour)

	manager.Increment().AddIssuedCertificate(true, cert1Month)
	manager.Increment().AddIssuedCertificate(false, cert1Year)

	time.Sleep(100 * time.Millisecond)

	// Calling stop again should not panic.
	manager.StopConsumerJob()

	jobCountLock.Lock()
	defer jobCountLock.Unlock()

	require.Equal(t, uint64(5), jobCount.IssuedCerts, "issued count mismatch")
	require.Equal(t, uint64(6), jobCount.StoredCerts, "stored count mismatch")
	// cert1Month: 730/730 = 1.0, cert1Year: 8760/730 = 12.0, total = 13.0
	require.InDelta(t, 13.0, jobCount.PkiDurationAdjustedCerts, 0.0001, "billable units mismatch")
	require.Zero(t, firstConsumerTotalCount.Load(), "first consumer should not have been called")
}
