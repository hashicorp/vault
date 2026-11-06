// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// crlPEMBlockType is the PEM preamble used for certificate revocation lists.
const crlPEMBlockType = "X509 CRL"

// oidDeltaCRLIndicator is the extension marking a delta CRL. It
// lists only the changes since a base CRL rather than every revocation.
var oidDeltaCRLIndicator = asn1.ObjectIdentifier{2, 5, 29, 27}

// DefaultCRLRefreshInterval is how often CRLs from disk are checked for
// changes when TLSConfig.CRLRefreshInterval is not set.
const DefaultCRLRefreshInterval = 5 * time.Minute

// crlReloadRetryInterval is how long a failed reload is backed off for. It avoids
// a persistently unreadable CRL file from being read, turning every handshake into a disk
// read.
const crlReloadRetryInterval = 5 * time.Second

// CRLProvider supplies the certificate revocation lists used to check whether
// the Vault server's certificate chain has been revoked.
type CRLProvider interface {
	CRLs() ([]*x509.RevocationList, error)
}

// ParseCRLs parses one or more certificate revocation lists. The input may be
// a single DER-encoded CRL, or PEM data containing any number of "X509 CRL"
// blocks.
func ParseCRLs(data []byte) ([]*x509.RevocationList, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, errors.New("no CRL data provided")
	}

	var crls []*x509.RevocationList
	rest := data
	sawPEM := false
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		sawPEM = true
		if block.Type != crlPEMBlockType {
			continue
		}
		crl, err := x509.ParseRevocationList(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse PEM-encoded CRL: %w", err)
		}
		crls = append(crls, crl)
	}

	if sawPEM {
		if len(crls) == 0 {
			return nil, fmt.Errorf("no %q blocks found in PEM data", crlPEMBlockType)
		}
		return crls, nil
	}

	// No PEM blocks at all, so treat the input as raw DER.
	crl, err := x509.ParseRevocationList(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CRL as PEM or DER: %w", err)
	}
	return []*x509.RevocationList{crl}, nil
}

// staticCRLProvider is a CRLProvider backed by a fixed set of CRLs.
type staticCRLProvider struct {
	crls []*x509.RevocationList
}

func (s *staticCRLProvider) CRLs() ([]*x509.RevocationList, error) {
	return s.crls, nil
}

// NewStaticCRLProvider returns a CRLProvider that always returns the given
// CRLs. It is intended for callers that manage CRL retrieval and refresh
// themselves.
func NewStaticCRLProvider(crls []*x509.RevocationList) CRLProvider {
	return &staticCRLProvider{crls: crls}
}

// fileCRLProvider loads CRLs from a file or from every file in a directory,
// re-reading them when they change on disk.
// Lazy refresh at refresh interval
//
// Disk I/O is never performed while p.mu is held
//
// Design Decision:
// If a refresh fails, the previously loaded CRLs keep being served and the
// refresh is retried after crlReloadRetryInterval (or sooner, if the configured
// interval is shorter).
//
// Decision Rationale:
// Serving slightly stale revocation data beats failing
// every connection because a CRL file was mid-rewrite.
type fileCRLProvider struct {
	path             string
	isDir            bool
	interval         time.Duration
	mu               sync.Mutex
	crls             []*x509.RevocationList
	fingerprint      string
	checkedAt        time.Time
	reloading        bool
	lastReloadFailed bool
}

// newFileCRLProvider builds a provider for path, which may be a single CRL
// file or a directory of them. The initial load happens immediately so that
// configuration errors surface from ConfigureTLS rather than from an unrelated
// request later on.
func newFileCRLProvider(path string, dir bool, interval time.Duration) (*fileCRLProvider, error) {
	if interval <= 0 {
		interval = DefaultCRLRefreshInterval
	}
	p := &fileCRLProvider{
		path:     path,
		isDir:    dir,
		interval: interval,
	}
	if err := p.reload(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *fileCRLProvider) CRLs() ([]*x509.RevocationList, error) {
	p.mu.Lock()
	if !p.refreshDueLocked() {
		crls := p.crls
		p.mu.Unlock()
		return crls, nil
	}
	// Claim the refresh so that concurrent callers are served the cached set
	// instead of piling up behind the read below.
	p.reloading = true
	cached := p.crls
	p.mu.Unlock()

	if err := p.reload(); err != nil {
		// Keep serving the last known good set.
		if len(cached) > 0 {
			return cached, nil
		}
		return nil, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	return p.crls, nil
}

// refreshDueLocked reports whether the source should be re-examined. p.mu must
// be held.
func (p *fileCRLProvider) refreshDueLocked() bool {
	if p.reloading {
		return false
	}
	interval := p.interval
	if p.lastReloadFailed && crlReloadRetryInterval < interval {
		interval = crlReloadRetryInterval
	}
	return time.Since(p.checkedAt) >= interval
}

// reload re-reads the CRL source if its fingerprint changed
func (p *fileCRLProvider) reload() error {
	files, fingerprint, err := p.stat()
	if err != nil {
		p.finishReload(err)
		return err
	}

	p.mu.Lock()
	unchanged := fingerprint == p.fingerprint && p.fingerprint != ""
	p.mu.Unlock()
	if unchanged {
		p.finishReload(nil)
		return nil
	}

	var crls []*x509.RevocationList
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			err = fmt.Errorf("error reading CRL file %q: %w", file, err)
			p.finishReload(err)
			return err
		}
		parsed, err := ParseCRLs(data)
		if err != nil {
			err = fmt.Errorf("error parsing CRL file %q: %w", file, err)
			p.finishReload(err)
			return err
		}
		crls = append(crls, parsed...)
	}
	if len(crls) == 0 {
		err := fmt.Errorf("no CRLs found in %q", p.path)
		p.finishReload(err)
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.crls = crls
	p.fingerprint = fingerprint
	p.checkedAt = time.Now()
	p.reloading = false
	p.lastReloadFailed = false
	return nil
}

// finishReload records the outcome of a reload attempt and releases the claim
// taken by CRLs.
func (p *fileCRLProvider) finishReload(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.checkedAt = time.Now()
	p.reloading = false
	p.lastReloadFailed = err != nil
}

// stat returns the files making up the CRL source along with their names,
// sizes and modification times. CRLPath is walked recursively
func (p *fileCRLProvider) stat() ([]string, string, error) {
	var files []string
	if p.isDir {
		err := filepath.Walk(p.path, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			return nil, "", fmt.Errorf("error reading CRL directory %q: %w", p.path, err)
		}
		if len(files) == 0 {
			return nil, "", fmt.Errorf("no files found in CRL directory %q", p.path)
		}
		sort.Strings(files)
	} else {
		files = []string{p.path}
	}

	var fingerprint strings.Builder
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return nil, "", fmt.Errorf("error reading CRL file %q: %w", file, err)
		}
		fmt.Fprintf(&fingerprint, "%s:%d:%d;", file, info.Size(), info.ModTime().UnixNano())
	}

	return files, fingerprint.String(), nil
}

// crlChecker checks a server's certificate chain against a set of CRLs.
type crlChecker struct {
	provider CRLProvider
}

// verifyConnection runs after normal chain verification
func (c *crlChecker) verifyConnection(cs tls.ConnectionState) error {
	chains := cs.VerifiedChains
	if len(chains) == 0 {
		return errors.New("tls: cannot check certificate revocation: no verified certificate chain, which means chain verification was disabled (InsecureSkipVerify)")
	}

	crls, err := c.provider.CRLs()
	if err != nil {
		return fmt.Errorf("tls: cannot check certificate revocation: %w", err)
	}

	var firstErr error
	for _, chain := range chains {
		err := checkChainRevocation(chain, crls)
		if err == nil {
			return nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// checkChainRevocation walks a verified chain from leaf to root
func checkChainRevocation(chain []*x509.Certificate, crls []*x509.RevocationList) error {
	if len(chain) == 0 {
		return errors.New("tls: cannot check certificate revocation: empty certificate chain")
	}

	for i, cert := range chain {
		issuer := cert
		if i+1 < len(chain) {
			issuer = chain[i+1]
		}

		// Skipping self-signed CA as it can revoke its own certificate.
		selfIssued := issuer == cert && isSelfSigned(cert)

		var covered, staleCRLs, unverifiableCRLs, deltaCRLs int
		for _, crl := range crls {
			if !crlIssuedBy(crl, issuer) {
				continue
			}
			if err := crl.CheckSignatureFrom(issuer); err != nil {
				// The issuer names match but the signature does not,
				// mark unverifiable
				unverifiableCRLs++
				continue
			}

			// Honour the entries even if the CRL
			// or is only a delta.
			for _, entry := range crl.RevokedCertificateEntries {
				if entry.SerialNumber.Cmp(cert.SerialNumber) == 0 {
					return fmt.Errorf("tls: certificate %q with serial number %s was revoked at %s by CRL from issuer %q",
						cert.Subject, cert.SerialNumber, entry.RevocationTime.Format(time.RFC3339), crl.Issuer)
				}
			}

			// A CRL that is past due tells us nothing about revocations that
			// happened after it was issued, so it does not count as coverage.
			if !crl.NextUpdate.IsZero() && time.Now().After(crl.NextUpdate) {
				staleCRLs++
				continue
			}

			// A delta CRL lists only the changes since some base CRL, so on its
			// own it does not establish that the certificate is unrevoked.
			if isDeltaCRL(crl) {
				deltaCRLs++
				continue
			}
			covered++
		}

		if i == 0 && covered == 0 && !selfIssued {
			return fmt.Errorf("tls: cannot check revocation of certificate %q: %s",
				cert.Subject, coverageFailureReason(issuer, staleCRLs, unverifiableCRLs, deltaCRLs))
		}
	}

	return nil
}

func coverageFailureReason(issuer *x509.Certificate, stale, unverifiable, delta int) string {
	var reasons []string
	if stale > 0 {
		reasons = append(reasons, fmt.Sprintf("%d are past their next update", stale))
	}
	if unverifiable > 0 {
		reasons = append(reasons, fmt.Sprintf("%d could not be verified against the issuer's key", unverifiable))
	}
	if delta > 0 {
		reasons = append(reasons, fmt.Sprintf("%d are delta CRLs, which do not list every revocation", delta))
	}
	if len(reasons) == 0 {
		return fmt.Sprintf("no CRL was supplied for issuer %q", issuer.Subject)
	}
	return fmt.Sprintf("of the CRLs from issuer %q, %s", issuer.Subject, strings.Join(reasons, " and "))
}

// isDeltaCRL reports whether the CRL carries the deltaCRLIndicator extension,
// marking it as listing only the changes since a base CRL.
func isDeltaCRL(crl *x509.RevocationList) bool {
	for _, ext := range crl.Extensions {
		if ext.Id.Equal(oidDeltaCRLIndicator) {
			return true
		}
	}
	return false
}

// crlIssuedBy reports whether the CRL claims to have been issued by the given
// certificate
func crlIssuedBy(crl *x509.RevocationList, issuer *x509.Certificate) bool {
	if len(crl.AuthorityKeyId) > 0 && len(issuer.SubjectKeyId) > 0 {
		return bytes.Equal(crl.AuthorityKeyId, issuer.SubjectKeyId)
	}
	return bytes.Equal(crl.RawIssuer, issuer.RawSubject)
}

// isSelfSigned reports whether the certificate signed itself.
func isSelfSigned(cert *x509.Certificate) bool {
	if !bytes.Equal(cert.RawSubject, cert.RawIssuer) {
		return false
	}
	return cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature) == nil
}
