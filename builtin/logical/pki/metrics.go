// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"errors"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
)

type CertificateCounter struct {
	certCountEnabled                    *atomic.Bool
	publishCertCountMetrics             *atomic.Bool
	certCount                           *atomic.Uint32
	revokedCertCount                    *atomic.Uint32
	certsCounted                        *atomic.Bool
	certCountError                      error
	possibleDoubleCountedSerials        []string
	possibleDoubleCountedRevokedSerials []string
	backendUuid                         string
}

func (c *CertificateCounter) IsInitialized() bool {
	return c.certsCounted.Load()
}

func (c *CertificateCounter) IsEnabled() bool {
	return c.certCountEnabled.Load()
}

func (c *CertificateCounter) Error() error {
	return c.certCountError
}

func (c *CertificateCounter) SetError(err error) {
	c.certCountError = err
}

func (c *CertificateCounter) ReconfigureWithTidyConfig(config *tidyConfig) bool {
	if config.MaintainCount {
		c.enableCertCounting(config.PublishMetrics)
	} else {
		c.disableCertCounting()
	}

	return config.MaintainCount
}

func (c *CertificateCounter) disableCertCounting() {
	c.possibleDoubleCountedRevokedSerials = nil
	c.possibleDoubleCountedSerials = nil
	c.certsCounted.Store(false)
	c.certCount.Store(0)
	c.revokedCertCount.Store(0)
	c.certCountError = errors.New("Cert Count is Disabled: enable via Tidy Config maintain_stored_certificate_counts")
	c.certCountEnabled.Store(false)
	c.publishCertCountMetrics.Store(false)
}

func (c *CertificateCounter) enableCertCounting(publishMetrics bool) {
	c.publishCertCountMetrics.Store(publishMetrics)
	c.certCountEnabled.Store(true)

	if !c.certsCounted.Load() {
		c.certCountError = errors.New("Certificate Counting Has Not Been Initialized, re-initialize this mount")
	}
}

func (c *CertificateCounter) InitializeCountsFromStorage(certs, revoked []string) {
	c.certCount.Add(uint32(len(certs)))
	c.revokedCertCount.Add(uint32(len(revoked)))

	c.pruneDuplicates(certs, revoked)
	c.certCountError = nil
	c.certsCounted.Store(true)

	c.emitTotalCertCountMetric()
}

func (c *CertificateCounter) pruneDuplicates(entries, revokedEntries []string) {
	// Now that the metrics are set, we can switch from appending newly-stored certificates to the possible double-count
	// list, and instead have them update the counter directly.  We need to do this so that we are looking at a static
	// slice of possibly double counted serials.  Note that certsCounted is computed before the storage operation, so
	// there may be some delay here.

	// Sort the listed-entries first, to accommodate that delay.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i] < entries[j]
	})

	sort.Slice(revokedEntries, func(i, j int) bool {
		return revokedEntries[i] < revokedEntries[j]
	})

	// We assume here that these lists are now complete.
	sort.Slice(c.possibleDoubleCountedSerials, func(i, j int) bool {
		return c.possibleDoubleCountedSerials[i] < c.possibleDoubleCountedSerials[j]
	})

	listEntriesIndex := 0
	possibleDoubleCountIndex := 0
	for {
		if listEntriesIndex >= len(entries) {
			break
		}
		if possibleDoubleCountIndex >= len(c.possibleDoubleCountedSerials) {
			break
		}
		if entries[listEntriesIndex] == c.possibleDoubleCountedSerials[possibleDoubleCountIndex] {
			// This represents a double-counted entry
			c.decrementTotalCertificatesCountNoReport()
			listEntriesIndex = listEntriesIndex + 1
			possibleDoubleCountIndex = possibleDoubleCountIndex + 1
			continue
		}
		if entries[listEntriesIndex] < c.possibleDoubleCountedSerials[possibleDoubleCountIndex] {
			listEntriesIndex = listEntriesIndex + 1
			continue
		}
		if entries[listEntriesIndex] > c.possibleDoubleCountedSerials[possibleDoubleCountIndex] {
			possibleDoubleCountIndex = possibleDoubleCountIndex + 1
			continue
		}
	}

	sort.Slice(c.possibleDoubleCountedRevokedSerials, func(i, j int) bool {
		return c.possibleDoubleCountedRevokedSerials[i] < c.possibleDoubleCountedRevokedSerials[j]
	})

	listRevokedEntriesIndex := 0
	possibleRevokedDoubleCountIndex := 0
	for {
		if listRevokedEntriesIndex >= len(revokedEntries) {
			break
		}
		if possibleRevokedDoubleCountIndex >= len(c.possibleDoubleCountedRevokedSerials) {
			break
		}
		if revokedEntries[listRevokedEntriesIndex] == c.possibleDoubleCountedRevokedSerials[possibleRevokedDoubleCountIndex] {
			// This represents a double-counted revoked entry
			c.decrementTotalRevokedCertificatesCountNoReport()
			listRevokedEntriesIndex = listRevokedEntriesIndex + 1
			possibleRevokedDoubleCountIndex = possibleRevokedDoubleCountIndex + 1
			continue
		}
		if revokedEntries[listRevokedEntriesIndex] < c.possibleDoubleCountedRevokedSerials[possibleRevokedDoubleCountIndex] {
			listRevokedEntriesIndex = listRevokedEntriesIndex + 1
			continue
		}
		if revokedEntries[listRevokedEntriesIndex] > c.possibleDoubleCountedRevokedSerials[possibleRevokedDoubleCountIndex] {
			possibleRevokedDoubleCountIndex = possibleRevokedDoubleCountIndex + 1
			continue
		}
	}

	c.possibleDoubleCountedRevokedSerials = nil
	c.possibleDoubleCountedSerials = nil
}

func (c *CertificateCounter) decrementTotalCertificatesCountNoReport() uint32 {
	newCount := c.certCount.Add(^uint32(0))
	return newCount
}

func (c *CertificateCounter) decrementTotalRevokedCertificatesCountNoReport() uint32 {
	newRevokedCertCount := c.revokedCertCount.Add(^uint32(0))
	return newRevokedCertCount
}

func (c *CertificateCounter) CertificateCount() uint32 {
	return c.certCount.Load()
}

func (c *CertificateCounter) RevokedCount() uint32 {
	return c.revokedCertCount.Load()
}

func (c *CertificateCounter) IncrementTotalCertificatesCount(certsCounted bool, newSerial string) {
	if c.certCountEnabled.Load() {
		c.certCount.Add(1)
		switch {
		case !certsCounted:
			// This is unsafe, but a good best-attempt
			if strings.HasPrefix(newSerial, issuing.PathCerts) {
				newSerial = newSerial[6:]
			}
			c.possibleDoubleCountedSerials = append(c.possibleDoubleCountedSerials, newSerial)
		default:
			c.emitTotalCertCountMetric()
		}
	}
}

// The "certsCounted" boolean here should be loaded from the backend certsCounted before the corresponding storage call:
// eg. certsCounted := certCounter.IsInitialized()
func (c *CertificateCounter) IncrementTotalRevokedCertificatesCount(certsCounted bool, newSerial string) {
	if c.certCountEnabled.Load() {
		c.revokedCertCount.Add(1)
		switch {
		case !certsCounted:
			// This is unsafe, but a good best-attempt
			if strings.HasPrefix(newSerial, "revoked/") { // allow passing in the path (revoked/serial) OR the serial
				newSerial = newSerial[8:]
			}
			c.possibleDoubleCountedRevokedSerials = append(c.possibleDoubleCountedRevokedSerials, newSerial)
		default:
			c.emitTotalRevokedCountMetric()
		}
	}
}

func (c *CertificateCounter) DecrementTotalCertificatesCountReport() {
	if c.certCountEnabled.Load() {
		c.decrementTotalCertificatesCountNoReport()
		c.emitTotalCertCountMetric()
	}
}

func (c *CertificateCounter) DecrementTotalRevokedCertificatesCountReport() {
	if c.certCountEnabled.Load() {
		c.decrementTotalRevokedCertificatesCountNoReport()
		c.emitTotalRevokedCountMetric()
	}
}

func (c *CertificateCounter) EmitCertStoreMetrics() {
	c.emitTotalCertCountMetric()
	c.emitTotalRevokedCountMetric()
}

func (c *CertificateCounter) emitTotalCertCountMetric() {
	if c.publishCertCountMetrics.Load() {
		certCount := float32(c.CertificateCount())
		metrics.SetGauge([]string{"secrets", "pki", c.backendUuid, "total_certificates_stored"}, certCount)
	}
}

func (c *CertificateCounter) emitTotalRevokedCountMetric() {
	if c.publishCertCountMetrics.Load() {
		revokedCount := float32(c.RevokedCount())
		metrics.SetGauge([]string{"secrets", "pki", c.backendUuid, "total_revoked_certificates_stored"}, revokedCount)
	}
}

func NewCertificateCounter(backendUuid string) *CertificateCounter {
	counter := &CertificateCounter{
		backendUuid:                         backendUuid,
		certCountEnabled:                    &atomic.Bool{},
		publishCertCountMetrics:             &atomic.Bool{},
		certCount:                           &atomic.Uint32{},
		revokedCertCount:                    &atomic.Uint32{},
		certsCounted:                        &atomic.Bool{},
		certCountError:                      errors.New("Initialize Not Yet Run, Cert Counts Unavailable"),
		possibleDoubleCountedSerials:        make([]string, 0, 250),
		possibleDoubleCountedRevokedSerials: make([]string, 0, 250),
	}

	return counter
}
