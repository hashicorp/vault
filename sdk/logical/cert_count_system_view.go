// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

type CertificateCountSystemView interface {
	GetCertificateCounter() CertificateCounter
}

type nullCertificateCounter struct{}

func (n *nullCertificateCounter) AddCount(_ CertCount) {
}

func (n *nullCertificateCounter) Increment() CertCountIncrementer {
	return NewCertCountIncrementer(n)
}

var _ CertificateCounter = (*nullCertificateCounter)(nil)

func NewNullCertificateCounter() CertificateCounter {
	return &nullCertificateCounter{}
}
