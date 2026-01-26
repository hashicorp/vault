// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

// CertificateCounter is an interface for incrementing the count of issued and stored
// certificates.
type CertificateCounter interface {
	// IncrementCount increments the count of issued and stored certificates.
	IncrementCount(issuedCerts, storedCerts uint64)

	// AddIssuedCertificate increments the issued certificate count by 1, and also the
	// stored certificate count if stored is true.
	AddIssuedCertificate(stored bool)
}
