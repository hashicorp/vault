// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

// PkiCertificateCounter is an interface for incrementing the count of issued and stored
// PKI certificates.
type PkiCertificateCounter interface {
	// IncrementCount increments the count of issued and stored certificates.
	IncrementCount(issuedCerts, storedCerts uint64)

	// AddIssuedCertificate increments the issued certificate count by 1, and also the
	// stored certificate count if stored is true.
	AddIssuedCertificate(stored bool)
}

type PkiCertificateCountSystemView interface {
	GetPkiCertificateCounter() PkiCertificateCounter
}

type nullPkiCertificateCounter struct{}

func (n *nullPkiCertificateCounter) IncrementCount(_, _ uint64) {
}

func (n *nullPkiCertificateCounter) AddIssuedCertificate(_ bool) {
}

var _ PkiCertificateCounter = (*nullPkiCertificateCounter)(nil)

func NewNullPkiCertificateCounter() PkiCertificateCounter {
	return &nullPkiCertificateCounter{}
}
