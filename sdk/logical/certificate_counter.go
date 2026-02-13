// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"crypto/x509"
	"math"
)

// CertificateCounter is an interface for incrementing the count of issued and stored
// certificates.
type CertificateCounter interface {
	// IncrementCount increments the count of issued and stored certificates.
	AddCount(params CertCount)

	// Increment returns a CertCountIncrementer that can be used to add
	// to the count.
	Increment() CertCountIncrementer
}

// CertCount represents the parameters for incrementing certificate counts.
type CertCount struct {
	IssuedCerts uint64 // TODO(victorr): Rename to PkiIssuedCerts
	StoredCerts uint64

	// PkiDurationAdjustedCerts tracks the normalized certificate duration units for billing
	// purposes. Each certificate's billable units = (Validity Hours ÷ 730), rounded to 4 decimal
	// places.
	PkiDurationAdjustedCerts float64
}

func (i *CertCount) Add(other CertCount) {
	i.IssuedCerts += other.IssuedCerts
	i.StoredCerts += other.StoredCerts
	i.PkiDurationAdjustedCerts += other.PkiDurationAdjustedCerts
}

func (i *CertCount) IsZero() bool {
	return i.IssuedCerts == 0 && i.StoredCerts == 0 && i.PkiDurationAdjustedCerts == 0
}

// durationAdjustedCertificateCount calculates the billable units for a certificate based on its
// validity duration. WARNING: Beware the maximum value for time.Duration (approximately 290 years).
//
// The calculation follows the billing specification:
// - Standard duration is 730 hours (1 month)
// - Units = (Validity Hours ÷ 730), rounded to 4 decimal places
// - Example: 1-year cert (8760 hours) = 12.0000 units
// - Example: 1-day cert (24 hours) = 0.0329 units
func durationAdjustedCertificateCount(validitySeconds int64) float64 {
	const standardDuration = 730.0
	validityHours := float64(validitySeconds) / 3600.0
	units := validityHours / standardDuration
	// Round to 4 decimal places
	return math.Round(units*10000) / 10000
}

type CertCountIncrementer interface {
	AddIssuedCertificate(stored bool, cert *x509.Certificate) CertCountIncrementer
}

type certCountIncrementer struct {
	counter CertificateCounter
}

var _ CertCountIncrementer = (*certCountIncrementer)(nil)

// NewCertCountIncrementer creates a new CertCountIncrementer for the given counter.
func NewCertCountIncrementer(counter CertificateCounter) CertCountIncrementer {
	return &certCountIncrementer{counter: counter}
}

// AddIssuedCertificate increments the issued certificate count by 1, the stored certificate
// count if stored is true, and adds the calculated billable units based on the certificate's
// validity duration.
// cert: The X.509 certificate to extract validity duration from.
func (c *certCountIncrementer) AddIssuedCertificate(stored bool, cert *x509.Certificate) CertCountIncrementer {
	validity := int64(cert.NotAfter.Unix() - cert.NotBefore.Unix())
	count := CertCount{
		IssuedCerts:              1,
		PkiDurationAdjustedCerts: durationAdjustedCertificateCount(validity),
	}
	if stored {
		count.StoredCerts = 1
	}
	c.counter.AddCount(count)

	return c
}
