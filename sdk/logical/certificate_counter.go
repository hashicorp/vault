// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

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
	IssuedCerts uint64
	StoredCerts uint64
}

func (i *CertCount) Add(other CertCount) {
	i.IssuedCerts += other.IssuedCerts
	i.StoredCerts += other.StoredCerts
}

func (i *CertCount) IsZero() bool {
	return i.IssuedCerts == 0 && i.StoredCerts == 0
}

type CertCountIncrementer interface {
	AddIssuedCertificate(stored bool) CertCountIncrementer
}

type certCountIncrementer struct {
	counter CertificateCounter
}

var _ CertCountIncrementer = (*certCountIncrementer)(nil)

// NewCertCountIncrementer creates a new CertCountIncrementer for the given counter.
func NewCertCountIncrementer(counter CertificateCounter) CertCountIncrementer {
	return &certCountIncrementer{counter: counter}
}

// AddIssuedCertificate increments the issued certificate count by 1, and also the
// stored certificate count if stored is true.
func (c *certCountIncrementer) AddIssuedCertificate(stored bool) CertCountIncrementer {
	count := CertCount{IssuedCerts: 1}
	if stored {
		count.StoredCerts = 1
	}
	c.counter.AddCount(count)

	return c
}
