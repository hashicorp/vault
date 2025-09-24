// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

// PkiCertificateCounter is an interface for incrementing the count of issued and stored
// PKI certificates.
type PkiCertificateCounter interface {
	IncrementCount(issuedCerts, storedCerts uint64)
}

type PkiCertificateCountSystemView interface {
	GetPkiCertificateCounter() PkiCertificateCounter
}

type nullPkiCertificateCounter struct{}

func (n *nullPkiCertificateCounter) IncrementCount(_, _ uint64) {
}

var _ PkiCertificateCounter = (*nullPkiCertificateCounter)(nil)

func NewNullPkiCertificateCounter() PkiCertificateCounter {
	return &nullPkiCertificateCounter{}
}
