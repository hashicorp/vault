// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains a certificate pool adapted from
// https://github.com/google/certificate-transparency-go/blob/v1.3.1/x509/cert_pool.go

package x509verify

import "crypto/x509"

// CertPool is a set of certificates.
type CertPool struct {
	bySubjectKeyID map[string][]int
	byName         map[string][]int
	certs          []*x509.Certificate
}

// NewCertPool returns a new, empty CertPool.
func NewCertPool() *CertPool {
	return &CertPool{
		bySubjectKeyID: make(map[string][]int),
		byName:         make(map[string][]int),
	}
}

// findPotentialParents returns the indexes of certificates in s which might
// have signed cert.
func (s *CertPool) findPotentialParents(cert *x509.Certificate) []int {
	if s == nil {
		return nil
	}
	var candidates []int
	if len(cert.AuthorityKeyId) > 0 {
		candidates = s.bySubjectKeyID[string(cert.AuthorityKeyId)]
	}
	if len(candidates) == 0 {
		candidates = s.byName[string(cert.RawIssuer)]
	}
	return candidates
}

func (s *CertPool) contains(cert *x509.Certificate) bool {
	if s == nil {
		return false
	}
	candidates := s.byName[string(cert.RawSubject)]
	for _, c := range candidates {
		if s.certs[c].Equal(cert) {
			return true
		}
	}
	return false
}

// AddCert adds a certificate to the pool.
func (s *CertPool) AddCert(cert *x509.Certificate) {
	if cert == nil {
		panic("adding nil Certificate to CertPool")
	}
	if s.contains(cert) {
		return
	}
	n := len(s.certs)
	s.certs = append(s.certs, cert)
	if len(cert.SubjectKeyId) > 0 {
		keyID := string(cert.SubjectKeyId)
		s.bySubjectKeyID[keyID] = append(s.bySubjectKeyID[keyID], n)
	}
	name := string(cert.RawSubject)
	s.byName[name] = append(s.byName[name], n)
}

// Len returns the number of certificates in the pool.
func (s *CertPool) Len() int {
	return len(s.certs)
}
