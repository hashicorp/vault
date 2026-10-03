// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package pki_cert_count

import (
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

// envVaultDisableCertCount is an environment variable that is used to disable
// all certificate counting.
const envVaultDisableCertCount = "VAULT_DISABLE_CERT_COUNT"

// consumerJobInterval is the interval the PkiCertificateCountManager uses
// for StartConsumerJob. It is a variable so that unit tests can override it.
var consumerJobInterval = 1 * time.Minute

// PkiCertificateCountConsumer is a callback for consumers of the PKI certificate counts.
type PkiCertificateCountConsumer func(issuedCount, storedCount uint64)

// PkiCertificateCountManager keeps track of issued and stored PKI certificate counts.
type PkiCertificateCountManager interface {
	logical.CertificateCounter
	// StartConsumerJob starts a background job that periodically reports the counts to the
	// given consumer. If a job is already running, it will be stopped and replaced.
	StartConsumerJob(consumer PkiCertificateCountConsumer)

	// StopConsumerJob stops the background job for the certificate count consumer, if one
	// is running.
	StopConsumerJob()

	// GetCounts returns the current counts of issued and stored certificates, without
	// consuming them. Meant to ease unit testing.
	GetCounts() (issuedCount, storedCount uint64)
}

// certCountManager is an implementation of PkiCertificateCountManager.
type certCountManager struct {
	issuedCount *atomic.Uint64
	storedCount *atomic.Uint64

	reportTimerStop     chan struct{}
	reportTimerStopLock sync.Mutex

	logger hclog.Logger
}

var _ PkiCertificateCountManager = (*certCountManager)(nil)

// InitPkiCertificateCountManager creates a new PkiCertificateCountManager, or a null
// implementation if certificate counting is disabled via the presence of the
// VAULT_DISABLE_CERT_COUNT environment variable (with any value).
func InitPkiCertificateCountManager(logger hclog.Logger) PkiCertificateCountManager {
	if os.Getenv(envVaultDisableCertCount) != "" {
		logger.Warn("PKI certificate counting disabled via environment variable")
		return newNullPkiCertificateCountManager()
	}
	return newPkiCertificateCountManager(logger)
}

func newPkiCertificateCountManager(logger hclog.Logger) PkiCertificateCountManager {
	ret := &certCountManager{
		issuedCount:     &atomic.Uint64{},
		storedCount:     &atomic.Uint64{},
		reportTimerStop: nil,
		logger:          logger,
	}
	return ret
}

func (m *certCountManager) StartConsumerJob(consumer PkiCertificateCountConsumer) {
	m.reportTimerStopLock.Lock()
	defer m.reportTimerStopLock.Unlock()

	m.stopConsumerJobWithLock()

	m.reportTimerStop = make(chan struct{})
	go m.reportLoop(m.reportTimerStop, consumer)
}

func (m *certCountManager) reportLoop(stop chan struct{}, consumer PkiCertificateCountConsumer) {
	reportTicker := time.NewTicker(consumerJobInterval)
	defer reportTicker.Stop()

	for {
		select {
		case <-reportTicker.C:
			m.consumeCount(consumer)

		case <-stop:
			reportTicker.Stop()
			m.consumeCount(consumer)
			return
		}
	}
}

func (m *certCountManager) consumeCount(consumer PkiCertificateCountConsumer) {
	issuedCount := m.issuedCount.Swap(0)
	storedCount := m.storedCount.Swap(0)
	consumer(issuedCount, storedCount)
}

func (m *certCountManager) StopConsumerJob() {
	m.reportTimerStopLock.Lock()
	defer m.reportTimerStopLock.Unlock()

	m.stopConsumerJobWithLock()
}

// stopConsumerJobWithLock must be called with reportTimerStopLock held.
func (m *certCountManager) stopConsumerJobWithLock() {
	var ch chan struct{}
	ch, m.reportTimerStop = m.reportTimerStop, nil

	if ch != nil {
		close(ch)
	}
}

func (m *certCountManager) AddIssuedCertificate(stored bool) {
	if stored {
		m.IncrementCount(1, 1)
	} else {
		m.IncrementCount(1, 0)
	}
}

func (m *certCountManager) IncrementCount(issuedCerts, storedCerts uint64) {
	issued := m.issuedCount.Add(issuedCerts)
	stored := m.storedCount.Add(storedCerts)
	m.logger.Trace("incremented in-memory PKI certificate counts", "issuedCerts", issued, "storedCerts", stored)
}

func (m *certCountManager) GetCounts() (issuedCount, storedCount uint64) {
	return m.issuedCount.Load(), m.storedCount.Load()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// nullPkiCertificateCountManager

type nullPkiCertificateCountManager struct{}

var _ PkiCertificateCountManager = (*nullPkiCertificateCountManager)(nil)

func newNullPkiCertificateCountManager() PkiCertificateCountManager {
	return &nullPkiCertificateCountManager{}
}

func (n *nullPkiCertificateCountManager) IncrementCount(_, _ uint64) {
	// nothing to do
}

func (n *nullPkiCertificateCountManager) AddIssuedCertificate(_ bool) {
	// nothing to do
}

func (n *nullPkiCertificateCountManager) StartConsumerJob(_ PkiCertificateCountConsumer) {
	// nothing to do
}

func (n *nullPkiCertificateCountManager) StopConsumerJob() {
	// nothing to do
}

func (n *nullPkiCertificateCountManager) GetCounts() (issuedCount, storedCount uint64) {
	return 0, 0
}
