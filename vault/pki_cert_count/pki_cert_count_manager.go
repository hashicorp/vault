// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package pki_cert_count

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

// consumerJobInterval is the interval the PkiCertificateCountManager uses
// for StartConsumerJob. It is a variable so that unit tests can override it.
var consumerJobInterval = 1 * time.Minute

// PkiCertificateCountConsumer is a callback for consumers of the PKI certificate counts.
type PkiCertificateCountConsumer func(logical.CertCount)

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
	GetCounts() logical.CertCount
}

// certCountManager is an implementation of PkiCertificateCountManager.
type certCountManager struct {
	count     logical.CertCount
	countLock sync.RWMutex

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
		count:           logical.CertCount{},
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
	m.countLock.Lock()
	defer m.countLock.Unlock()

	increment := m.count
	m.count = logical.CertCount{}

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

	m.logger.Trace("incremented in-memory PKI certificate counts", "issuedCerts", m.count.IssuedCerts, "storedCerts", m.count.StoredCerts)
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
// nullPkiCertificateCountManager

type nullPkiCertificateCountManager struct{}

var _ PkiCertificateCountManager = (*nullPkiCertificateCountManager)(nil)

func newNullPkiCertificateCountManager() PkiCertificateCountManager {
	return &nullPkiCertificateCountManager{}
}

func (n *nullPkiCertificateCountManager) AddCount(_ logical.CertCount) {
	// nothing to do
}

func (n *nullPkiCertificateCountManager) Increment() logical.CertCountIncrementer {
	return logical.NewCertCountIncrementer(n)
}

func (n *nullPkiCertificateCountManager) AddIssuedCertificate(_ bool, _ *x509.Certificate) {
	// nothing to do
}

func (n *nullPkiCertificateCountManager) StartConsumerJob(_ PkiCertificateCountConsumer) {
	// nothing to do
}

func (n *nullPkiCertificateCountManager) StopConsumerJob() {
	// nothing to do
}

func (n *nullPkiCertificateCountManager) GetCounts() (issuedCount logical.CertCount) {
	return logical.CertCount{}
}
