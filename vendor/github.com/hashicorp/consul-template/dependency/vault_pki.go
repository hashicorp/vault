// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dependency

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Ensure implements
var _ Dependency = (*VaultPKIQuery)(nil)

// Return type containing PEMs as strings
type PemEncoded struct {
	Cert, Key, CA string
	CAChain       []string
}

func (a PemEncoded) Equals(b PemEncoded) bool {
	if a.CA != b.CA || a.Cert != b.Cert || a.Key != b.Key {
		return false
	}

	if len(a.CAChain) != len(b.CAChain) {
		return false
	}

	for i, v := range a.CAChain {
		if v != b.CAChain[i] {
			return false
		}
	}
	return true
}

func (a PemEncoded) CaChainContains(item string) bool {
	for _, v := range a.CAChain {
		if v == item {
			return true
		}
	}
	return false
}

// a wrapper to mimic v2 secrets Data wrapper
func (a PemEncoded) Data() PemEncoded {
	return a
}

// VaultPKIQuery is the dependency to Vault for a secret
type VaultPKIQuery struct {
	stopCh  chan struct{}
	sleepCh chan time.Duration

	pkiPath  string
	data     map[string]interface{}
	filePath string
}

// NewVaultReadQuery creates a new datacenter dependency.
func NewVaultPKIQuery(urlpath, filepath string, data map[string]interface{}) (*VaultPKIQuery, error) {
	urlpath = strings.TrimSpace(urlpath)
	urlpath = strings.Trim(urlpath, "/")
	if urlpath == "" {
		return nil, fmt.Errorf("vault.read: invalid format: %q", urlpath)
	}

	secretURL, err := url.Parse(urlpath)
	if err != nil {
		return nil, err
	}

	return &VaultPKIQuery{
		stopCh:   make(chan struct{}, 1),
		sleepCh:  make(chan time.Duration, 1),
		pkiPath:  secretURL.Path,
		data:     data,
		filePath: filepath,
	}, nil
}

// Fetch queries the Vault API
func (d *VaultPKIQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}
	select {
	case dur := <-d.sleepCh:
		time.Sleep(dur)
	default:
	}

	needsRenewal := fmt.Errorf("needs renewal")
	getPEMs := func(renew bool) (PemEncoded, error) {
		rawPems, err := os.ReadFile(d.filePath)
		if renew || err != nil || len(rawPems) == 0 {
			rawPems, err = d.fetchPEMs(clients)
			// no need to write cert to file as it is the template dest
		}
		if err != nil {
			return PemEncoded{}, err
		}

		encPems, cert, err := pemsCert(rawPems)
		if err != nil {
			return encPems, err
		}

		if sleepFor, ok := goodFor(cert); ok {
			d.sleepCh <- sleepFor
			return encPems, nil
		}
		return encPems, needsRenewal
	}

	encPems, err := getPEMs(false)
	switch err {
	case nil:
	case needsRenewal:
		encPems, err = getPEMs(true)
		if err != nil {
			return PemEncoded{}, nil, err
		}
	default:
		return PemEncoded{}, nil, err
	}
	return respWithMetadata(encPems)
}

// returns time left in ~90% of the original lease and a boolean
// that returns false if cert needs renewing, true otherwise
func goodFor(cert *x509.Certificate) (time.Duration, bool) {
	// If we got called with a cert that doesn't exist, just say there's no
	// time left, and it needs to be renewed
	if cert == nil {
		return 0, false
	}
	start, end := cert.NotBefore.UTC(), cert.NotAfter.UTC()
	now := time.Now().UTC()
	if end.Before(now) || end.Equal(now) { // already expired
		return 0, false
	}

	lifespanDur := end.Sub(start)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	lifespanMilliseconds := lifespanDur.Milliseconds()
	// calculate the 'time the certificate should be rotated' by figuring out -3%
	// +3% + VaultLeaseRenewalThreshold of the lifespan and adding it to the
	// start
	rotationTime := start.Add(time.Millisecond * time.Duration(
		float64(lifespanMilliseconds)*VaultLeaseRenewalThreshold+float64(lifespanMilliseconds*(int64(r.Intn(6)-3)/100.0)),
	))

	// after we have the 'time the certificate should be rotated', figure out how
	// far it is from now to sleep
	sleepFor := rotationTime.Sub(now)
	if sleepFor <= 0 {
		return 0, false
	}

	return sleepFor, true
}

// loops through all pem encoded blocks in the byte stream
// returning the CA, Certificate and Private Key PEM strings
// also returns the cert for the Certificate as we have it and need it
func pemsCert(encoded []byte) (PemEncoded, *x509.Certificate, error) {
	var block *pem.Block
	var cert *x509.Certificate
	var encPems PemEncoded
	var aPem []byte

	for {
		aPem, encoded = nextPem(encoded)
		// scan, find and parse PEM blocks
		block, _ = pem.Decode(aPem)

		switch {
		case block == nil: // end of scan, no more PEMs found
			return encPems, cert, nil
		case strings.HasSuffix(block.Type, "PRIVATE KEY"):
			// PKCS#1 and PKCS#8 matching to find private key
			encPems.Key = string(pem.EncodeToMemory(block))
			continue
		}
		// CERTIFICATE PEM blocks (Cert and CA) are left
		maybeCert, err := x509.ParseCertificate(block.Bytes)
		switch {
		case err != nil:
			return PemEncoded{}, nil, err
		case maybeCert.IsCA:
			if encPems.CA == "" {
				// set the first CA found to CA to be backward compatible
				encPems.CA = string(pem.EncodeToMemory(block))
			}
			if !encPems.CaChainContains(string(pem.EncodeToMemory(block))) {
				encPems.CAChain = append(encPems.CAChain, string(pem.EncodeToMemory(block)))
			}
		default: // the certificate
			cert = maybeCert
			encPems.Cert = string(pem.EncodeToMemory(block))
		}
	}
}

// find the next PEM in the byte stream
func nextPem(encoded []byte) (aPem []byte, theRest []byte) {
	start := bytes.Index(encoded, []byte("-----BEGIN"))
	if start >= 0 { // finds the PEM and pulls it to decode
		encoded = encoded[start:] // clip pre-pem junk
		// find the end
		end := bytes.Index(encoded, []byte("-----END")) + 8
		end = end + bytes.Index(encoded[end:], []byte("-----")) + 5
		// the PEM padded with newlines (what pem.Decode likes)
		aPem = append([]byte("\n"), encoded[:end]...)
		aPem = append(aPem, []byte("\n")...)
		theRest = encoded[end:] // the rest
	}
	return aPem, theRest
}

// Vault call to fetch the PKI Cert PEM data
func (d *VaultPKIQuery) fetchPEMs(clients *ClientSet) ([]byte, error) {
	vaultSecret, err := clients.Vault().Logical().Write(d.pkiPath, d.data)
	switch {
	case err != nil:
		return nil, errors.Wrap(err, d.String())
	case vaultSecret == nil:
		return nil, fmt.Errorf("no secret exists at %s", d.pkiPath)
	}
	printVaultWarnings(d, vaultSecret.Warnings)
	pems := bytes.Buffer{}

	for k, v := range vaultSecret.Data {
		switch v := v.(type) {
		case string:
			pems.WriteString(v + "\n")
		case []interface{}:
			if k == "ca_chain" {
				for _, item := range v {
					switch item := item.(type) {
					case string:
						pems.WriteString(item + "\n")
					}
				}
			}
		}
	}

	return pems.Bytes(), nil
}

// CanShare returns if this dependency is shareable.
func (d *VaultPKIQuery) CanShare() bool {
	return false
}

// Stop halts the given dependency's fetch.
func (d *VaultPKIQuery) Stop() {
	close(d.stopCh)
}

// String returns the human-friendly version of this dependency.
func (d *VaultPKIQuery) String() string {
	return fmt.Sprintf("vault.pki(%s->%s)", d.pkiPath, d.filePath)
}

// Type returns the type of this dependency.
func (d *VaultPKIQuery) Type() Type {
	return TypeVault
}
