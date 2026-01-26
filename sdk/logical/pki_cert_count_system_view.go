// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

type PkiCertificateCountSystemView interface {
	GetPkiCertificateCounter() CertificateCounter
}

type nullPkiCertificateCounter struct{}

func (n *nullPkiCertificateCounter) IncrementCount(_, _ uint64) {
}

func (n *nullPkiCertificateCounter) AddIssuedCertificate(_ bool) {
}

var _ CertificateCounter = (*nullPkiCertificateCounter)(nil)

func NewNullPkiCertificateCounter() CertificateCounter {
	return &nullPkiCertificateCounter{}
}
