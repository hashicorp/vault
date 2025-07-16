// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package healthcheck

import (
	"bytes"
	"crypto/x509"
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type RootIssuedLeaves struct {
	Enabled            bool
	UnsupportedVersion bool

	CertsToFetch int

	FetchIssues map[string]*PathFetch
	RootCertMap map[string]*x509.Certificate
	LeafCertMap map[string]*x509.Certificate
}

func NewRootIssuedLeavesCheck() Check {
	return &RootIssuedLeaves{
		FetchIssues: make(map[string]*PathFetch),
		RootCertMap: make(map[string]*x509.Certificate),
		LeafCertMap: make(map[string]*x509.Certificate),
	}
}

func (h *RootIssuedLeaves) Name() string {
	return "root_issued_leaves"
}

func (h *RootIssuedLeaves) IsEnabled() bool {
	return h.Enabled
}

func (h *RootIssuedLeaves) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"certs_to_fetch": 100,
	}
}

func (h *RootIssuedLeaves) LoadConfig(config map[string]interface{}) error {
	count, err := parseutil.SafeParseIntRange(config["certs_to_fetch"], 1, 100000)
	if err != nil {
		return fmt.Errorf("error parsing %v.certs_to_fetch: %w", h.Name(), err)
	}
	h.CertsToFetch = int(count)

	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *RootIssuedLeaves) FetchResources(e *Executor) error {
	exit, _, issuers, err := pkiFetchIssuersList(e, func() {
		h.UnsupportedVersion = true
	})
	if exit || err != nil {
		return err
	}

	for _, issuer := range issuers {
		skip, pathFetch, cert, err := pkiFetchIssuer(e, issuer, func() {
			h.UnsupportedVersion = true
		})
		h.FetchIssues[issuer] = pathFetch
		if skip || err != nil {
			if err != nil {
				return err
			}
			continue
		}

		// Ensure we only check Root CAs.
		if !bytes.Equal(cert.RawSubject, cert.RawIssuer) {
			continue
		}
		if err := cert.CheckSignatureFrom(cert); err != nil {
			continue
		}

		h.RootCertMap[issuer] = cert
	}

	exit, f, leaves, err := pkiFetchLeavesList(e, func() {
		h.UnsupportedVersion = true
	})
	if exit || err != nil {
		if f != nil && f.IsSecretPermissionsError() {
			for _, issuer := range issuers {
				h.FetchIssues[issuer] = f
			}
		}
		return err
	}

	var leafCount int
	for _, serial := range leaves {
		if leafCount >= h.CertsToFetch {
			break
		}

		skip, _, cert, err := pkiFetchLeaf(e, serial, func() {
			h.UnsupportedVersion = true
		})
		if skip || err != nil {
			if err != nil {
				return err
			}
			continue
		}

		// Ignore other CAs.
		if cert.BasicConstraintsValid && cert.IsCA {
			continue
		}

		leafCount += 1
		h.LeafCertMap[serial] = cert
	}

	return nil
}

func (h *RootIssuedLeaves) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/{{mount}}/issuers",
			Message:  "This health check requires Vault 1.11+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	for issuer, fetchPath := range h.FetchIssues {
		if fetchPath != nil && fetchPath.IsSecretPermissionsError() {
			delete(h.RootCertMap, issuer)
			ret := Result{
				Status:   ResultInsufficientPermissions,
				Endpoint: fetchPath.Path,
				Message:  "Without this information, this health check is unable to function.",
			}

			if e.Client.Token() == "" {
				ret.Message = "No token available so unable for the endpoint for this mount. " + ret.Message
			} else {
				ret.Message = "This token lacks permission for the endpoint for this mount. " + ret.Message
			}

			results = append(results, &ret)
		}
	}

	issuerHasLeaf := make(map[string]bool)
	for serial, leaf := range h.LeafCertMap {
		if len(issuerHasLeaf) == len(h.RootCertMap) {
			break
		}

		for issuer, root := range h.RootCertMap {
			if issuerHasLeaf[issuer] {
				continue
			}

			if !bytes.Equal(leaf.RawIssuer, root.RawSubject) {
				continue
			}

			if err := leaf.CheckSignatureFrom(root); err != nil {
				continue
			}

			ret := Result{
				Status:   ResultWarning,
				Endpoint: "/{{mount}}/issuer/" + issuer,
				Message:  fmt.Sprintf("Root issuer has directly issued non-CA leaf certificates (%v) instead of via an intermediate CA. This can make rotating the root CA harder as direct cross-signing of the roots must be used, rather than cross-signing of the intermediates. It is encouraged to set up and use an intermediate CA and tidy the mount when all directly issued leaves have expired.", serial),
			}

			issuerHasLeaf[issuer] = true

			results = append(results, &ret)
		}
	}

	if len(results) == 0 && len(h.RootCertMap) > 0 {
		ret := Result{
			Status:   ResultOK,
			Endpoint: "/{{mount}}/certs",
			Message:  "Root certificate(s) in this mount have not directly issued non-CA leaf certificates.",
		}

		results = append(results, &ret)
	}

	return
}
