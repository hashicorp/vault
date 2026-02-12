// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

type PkiCertificateCountSystemView interface {
	GetPkiCertificateCounter() CertificateCounter
}

type nullPkiCertificateCounter struct{}

func (n *nullPkiCertificateCounter) AddCount(_ CertCount) {
}

func (n *nullPkiCertificateCounter) Increment() CertCountIncrementer {
	return NewCertCountIncrementer(n)
}

var _ CertificateCounter = (*nullPkiCertificateCounter)(nil)

func NewNullPkiCertificateCounter() CertificateCounter {
	return &nullPkiCertificateCounter{}
}
