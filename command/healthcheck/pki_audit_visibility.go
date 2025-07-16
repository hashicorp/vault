// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package healthcheck

import (
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

var VisibleReqParams = []string{
	"csr",
	"certificate",
	"issuer_ref",
	"common_name",
	"alt_names",
	"other_sans",
	"ip_sans",
	"uri_sans",
	"ttl",
	"not_after",
	"serial_number",
	"key_type",
	"private_key_format",
	"managed_key_name",
	"managed_key_id",
	"ou",
	"organization",
	"country",
	"locality",
	"province",
	"street_address",
	"postal_code",
	"permitted_dns_domains",
	"policy_identifiers",
	"ext_key_usage_oids",
}

var VisibleRespParams = []string{
	"certificate",
	"issuing_ca",
	"serial_number",
	"error",
	"ca_chain",
}

var HiddenReqParams = []string{
	"private_key",
	"pem_bundle",
}

var HiddenRespParams = []string{
	"private_key",
	"pem_bundle",
}

type AuditVisibility struct {
	Enabled            bool
	UnsupportedVersion bool

	IgnoredParameters map[string]bool
	TuneData          map[string]interface{}
	Fetcher           *PathFetch
}

func NewAuditVisibilityCheck() Check {
	return &AuditVisibility{
		IgnoredParameters: make(map[string]bool),
	}
}

func (h *AuditVisibility) Name() string {
	return "audit_visibility"
}

func (h *AuditVisibility) IsEnabled() bool {
	return h.Enabled
}

func (h *AuditVisibility) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"ignored_parameters": []string{},
	}
}

func (h *AuditVisibility) LoadConfig(config map[string]interface{}) error {
	var err error

	coerced, err := StringList(config["ignored_parameters"])
	if err != nil {
		return fmt.Errorf("error parsing %v.ignored_parameters: %v", h.Name(), err)
	}
	for _, ignored := range coerced {
		h.IgnoredParameters[ignored] = true
	}

	h.Enabled, err = parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}

	return nil
}

func (h *AuditVisibility) FetchResources(e *Executor) error {
	var exit bool
	var err error

	exit, h.Fetcher, h.TuneData, err = fetchMountTune(e, func() {
		h.UnsupportedVersion = true
	})

	if exit || err != nil {
		return err
	}
	return nil
}

func (h *AuditVisibility) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/sys/mounts/{{mount}}/tune",
			Message:  "This health check requires Vault 1.9+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	if h.Fetcher.IsSecretPermissionsError() {
		ret := Result{
			Status:   ResultInsufficientPermissions,
			Endpoint: "/sys/mounts/{{mount}}/tune",
			Message:  "Without this information, this health check is unable to function.",
		}

		if e.Client.Token() == "" {
			ret.Message = "No token available so unable read the tune endpoint for this mount. " + ret.Message
		} else {
			ret.Message = "This token lacks permission to read the tune endpoint for this mount. " + ret.Message
		}

		results = append(results, &ret)
		return
	}

	sourceMap := map[string][]string{
		"audit_non_hmac_request_keys":  VisibleReqParams,
		"audit_non_hmac_response_keys": VisibleRespParams,
	}
	for source, visibleList := range sourceMap {
		actual, err := StringList(h.TuneData[source])
		if err != nil {
			return nil, fmt.Errorf("error parsing %v from server: %v", source, err)
		}

		for _, param := range visibleList {
			found := false
			for _, tuned := range actual {
				if param == tuned {
					found = true
					break
				}
			}

			if !found {
				ret := Result{
					Status:   ResultInformational,
					Endpoint: "/sys/mounts/{{mount}}/tune",
					Message:  fmt.Sprintf("Mount currently HMACs %v because it is not in %v; as this is not a sensitive security parameter, it is encouraged to disable HMACing to allow better auditing of the PKI engine.", param, source),
				}
				results = append(results, &ret)
			}
		}
	}

	sourceMap = map[string][]string{
		"audit_non_hmac_request_keys":  HiddenReqParams,
		"audit_non_hmac_response_keys": HiddenRespParams,
	}
	for source, hiddenList := range sourceMap {
		actual, err := StringList(h.TuneData[source])
		if err != nil {
			return nil, fmt.Errorf("error parsing %v from server: %v", source, err)
		}
		for _, param := range hiddenList {
			found := false
			for _, tuned := range actual {
				if param == tuned {
					found = true
					break
				}
			}

			if found {
				ret := Result{
					Status:   ResultWarning,
					Endpoint: "/sys/mounts/{{mount}}/tune",
					Message:  fmt.Sprintf("Mount currently doesn't HMAC %v because it is in %v; as this is a sensitive security parameter it is encouraged to HMAC it in the Audit logs.", param, source),
				}
				results = append(results, &ret)
			}
		}
	}

	if len(results) == 0 {
		ret := Result{
			Status:   ResultOK,
			Endpoint: "/sys/mounts/{{mount}}/tune",
			Message:  "Mount audit information is configured appropriately.",
		}
		results = append(results, &ret)
	}

	return
}
