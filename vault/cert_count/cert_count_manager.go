// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package cert_count

import (
	"crypto/x509"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

// envVaultDisableCertCount is an environment variable that is used to disable
// all certificate counting.
const envVaultDisableCertCount = "VAULT_DISABLE_CERT_COUNT"

// consumerJobInterval is the interval the CertificateCountManager uses
// for StartConsumerJob. It is a variable so that unit tests can override it.
var consumerJobInterval = 1 * time.Minute

// CertificateCountConsumer is a callback for consumers of the certificate counts.
type CertificateCountConsumer func(logical.CertCount)

// CertificateCountManager keeps track of issued and stored certificate counts.
type CertificateCountManager interface {
	logical.CertificateCounter
	// StartConsumerJob starts a background job that periodically reports the counts to the
	// given consumer. If a job is already running, it will be stopped and replaced.
	StartConsumerJob(consumer CertificateCountConsumer)

	// StopConsumerJob stops the background job for the certificate count consumer, if one
	// is running.
	StopConsumerJob()

	// GetCounts returns the current counts of issued and stored certificates, without
	// consuming them. Meant to ease unit testing.
	GetCounts() logical.CertCount
}

// certCountManager is an implementation of CertificateCountManager.
type certCountManager struct {
	count     logical.CertCount
	countLock sync.RWMutex

	reportTimerStop     chan struct{}
	reportTimerStopLock sync.Mutex

	logger hclog.Logger
}

var _ CertificateCountManager = (*certCountManager)(nil)

// InitCertificateCountManager creates a new CertificateCountManager, or a null
// implementation if certificate counting is disabled via the presence of the
// VAULT_DISABLE_CERT_COUNT environment variable (with any value).
func InitCertificateCountManager(logger hclog.Logger) CertificateCountManager {
	if os.Getenv(envVaultDisableCertCount) != "" {
		logger.Warn("certificate counting disabled via environment variable")
		return newNullCertificateCountManager()
	}
	return newCertificateCountManager(logger)
}

func newCertificateCountManager(logger hclog.Logger) CertificateCountManager {
	ret := &certCountManager{
		count:           logical.CertCount{},
		reportTimerStop: nil,
		logger:          logger,
	}
	return ret
}

func (m *certCountManager) StartConsumerJob(consumer CertificateCountConsumer) {
	m.reportTimerStopLock.Lock()
	defer m.reportTimerStopLock.Unlock()

	m.stopConsumerJobWithLock()

	m.reportTimerStop = make(chan struct{})
	go m.reportLoop(m.reportTimerStop, consumer)
}

func (m *certCountManager) reportLoop(stop chan struct{}, consumer CertificateCountConsumer) {
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

func (m *certCountManager) consumeCount(consumer CertificateCountConsumer) {
	m.countLock.Lock()
	increment := m.count
	m.count = logical.CertCount{}
	m.countLock.Unlock()

	consumer(increment)
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

func (m *certCountManager) AddCount(params logical.CertCount) {
	m.countLock.Lock()
	defer m.countLock.Unlock()

	m.count.Add(params)

	m.logger.Trace("incremented in-memory certificate counts", "issuedCerts", m.count.IssuedCerts,
		"storedCerts", m.count.StoredCerts, "pkiDurationAdjustedCerts", m.count.PkiDurationAdjustedCerts)
}

func (m *certCountManager) Increment() logical.CertCountIncrementer {
	return logical.NewCertCountIncrementer(m)
}

func (m *certCountManager) GetCounts() (issuedCount logical.CertCount) {
	m.countLock.RLock()
	defer m.countLock.RUnlock()
	ret := logical.CertCount{}
	ret.Add(m.count)
	return ret
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// nullCertificateCountManager

type nullCertificateCountManager struct{}

var _ CertificateCountManager = (*nullCertificateCountManager)(nil)

func newNullCertificateCountManager() CertificateCountManager {
	return &nullCertificateCountManager{}
}

func (n *nullCertificateCountManager) AddCount(_ logical.CertCount) {
	// nothing to do
}

func (n *nullCertificateCountManager) Increment() logical.CertCountIncrementer {
	return logical.NewCertCountIncrementer(n)
}

func (n *nullCertificateCountManager) AddIssuedCertificate(_ bool, _ *x509.Certificate) {
	// nothing to do
}

func (n *nullCertificateCountManager) StartConsumerJob(_ CertificateCountConsumer) {
	// nothing to do
}

func (n *nullCertificateCountManager) StopConsumerJob() {
	// nothing to do
}

func (n *nullCertificateCountManager) GetCounts() (issuedCount logical.CertCount) {
	return logical.CertCount{}
}
