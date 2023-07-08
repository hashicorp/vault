// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcheck

import (
	"bytes"
	"crypto/x509"
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type HardwareBackedRoot struct {
	Enabled bool

	UnsupportedVersion bool

	FetchIssues  map[string]*PathFetch
	IssuerKeyMap map[string]string
	KeyIsManaged map[string]string
}

func NewHardwareBackedRootCheck() Check {
	return &HardwareBackedRoot{
		FetchIssues:  make(map[string]*PathFetch),
		IssuerKeyMap: make(map[string]string),
		KeyIsManaged: make(map[string]string),
	}
}

func (h *HardwareBackedRoot) Name() string {
	return "hardware_backed_root"
}

func (h *HardwareBackedRoot) IsEnabled() bool {
	return h.Enabled
}

func (h *HardwareBackedRoot) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"enabled": false,
	}
}

func (h *HardwareBackedRoot) LoadConfig(config map[string]interface{}) error {
	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *HardwareBackedRoot) FetchResources(e *Executor) error {
	exit, _, issuers, err := pkiFetchIssuersList(e, func() {
		h.UnsupportedVersion = true
	})
	if exit || err != nil {
		return err
	}

	for _, issuer := range issuers {
		skip, ret, entry, err := pkiFetchIssuerEntry(e, issuer, func() {
			h.UnsupportedVersion = true
		})
		if skip || err != nil || entry == nil {
			if err != nil {
				return err
			}
			h.FetchIssues[issuer] = ret
			continue
		}

		// Ensure we only check Root CAs.
		cert := ret.ParsedCache["certificate"].(*x509.Certificate)
		if !bytes.Equal(cert.RawSubject, cert.RawIssuer) {
			continue
		}
		if err := cert.CheckSignatureFrom(cert); err != nil {
			continue
		}

		// Ensure we only check issuers with keys.
		keyId, present := entry["key_id"].(string)
		if !present || len(keyId) == 0 {
			continue
		}

		h.IssuerKeyMap[issuer] = keyId
		skip, ret, keyEntry, err := pkiFetchKeyEntry(e, keyId, func() {
			h.UnsupportedVersion = true
		})
		if skip || err != nil || keyEntry == nil {
			if err != nil {
				return err
			}

			h.FetchIssues[issuer] = ret
			continue
		}

		uuid, present := keyEntry["managed_key_id"].(string)
		if present {
			h.KeyIsManaged[keyId] = uuid
		}
	}

	return nil
}

func (h *HardwareBackedRoot) Evaluate(e *Executor) (results []*Result, err error) {
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
			delete(h.IssuerKeyMap, issuer)
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

	for name, keyId := range h.IssuerKeyMap {
		var ret Result
		ret.Status = ResultInformational
		ret.Endpoint = "/{{mount}}/issuer/" + name
		ret.Message = "Root issuer was created using Vault-backed software keys; for added safety of long-lived, important root CAs, you may wish to consider using a HSM or KSM Managed Key to store key material for this issuer."

		uuid, present := h.KeyIsManaged[keyId]
		if present {
			ret.Status = ResultOK
			ret.Message = fmt.Sprintf("Root issuer was backed by a HSM or KMS Managed Key: %v.", uuid)
		}

		results = append(results, &ret)
	}

	return
}
