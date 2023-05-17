// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcheck

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type VerifyAcmeBasics struct {
	Enabled            bool
	UnsupportedVersion bool

	TuneFetcher *PathFetch
	TuneData    map[string]interface{}

	ClusterConfigFetcher *PathFetch
}

func NewVerifyAcmeBasics() Check {
	return &VerifyAcmeBasics{}
}

func (h *VerifyAcmeBasics) Name() string {
	return "verify_acme_basics"
}

func (h *VerifyAcmeBasics) IsEnabled() bool {
	return h.Enabled
}

func (h *VerifyAcmeBasics) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func (h *VerifyAcmeBasics) LoadConfig(config map[string]interface{}) error {
	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *VerifyAcmeBasics) FetchResources(e *Executor) error {
	var err error

	// We ignore errors from fetchMountTune as we want both fetches to occur no matter what.
	_, h.TuneFetcher, h.TuneData, _ = fetchMountTune(e, func() {
		h.UnsupportedVersion = true
	})

	h.ClusterConfigFetcher, err = e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/config/cluster")
	if err != nil {
		return fmt.Errorf("failed to fetch mount's config/cluster value: %v", err)
	}

	if h.ClusterConfigFetcher.IsUnsupportedPathError() {
		h.UnsupportedVersion = true
	}

	return nil
}

func verifyLocalPathUrl(h *VerifyAcmeBasics) error {
	localPath, ok := h.ClusterConfigFetcher.Secret.Data["path"]
	if !ok {
		return fmt.Errorf("path variable missing from config")
	}

	if localPath == "" {
		return fmt.Errorf("path variable not configured within /{{mount}}/config/cluster")
	}

	parsedUrl, err := url.Parse(localPath.(string))
	if err != nil {
		return fmt.Errorf("failed to parse URL from path config: %v: %w", localPath, err)
	}

	if parsedUrl.Scheme != "https" {
		return fmt.Errorf("the configured path value in /{{mount}}/config/cluster was not using an https scheme")
	}
	return nil
}

func (h *VerifyAcmeBasics) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion || h.TuneFetcher == nil || h.ClusterConfigFetcher == nil {
		endpoint := "Unknown all fetchers failed"
		if h.TuneFetcher != nil && h.ClusterConfigFetcher != nil {
			endpoint = h.TuneFetcher.Path + " or " + h.ClusterConfigFetcher.Path
		} else if h.TuneFetcher != nil {
			endpoint = h.TuneFetcher.Path
		} else if h.ClusterConfigFetcher != nil {
			endpoint = h.ClusterConfigFetcher.Path
		}
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: endpoint,
			Message:  "This health check requires Vault 1.13+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	// Tune information
	if h.TuneFetcher.IsSecretPermissionsError() {
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
	} else {
		resp, err := StringList(h.TuneData["allowed_response_headers"])
		if err != nil {
			return nil, fmt.Errorf("unable to parse value from server for allowed_response_headers: %w", err)
		}

		requiredResponseHeaders := []string{"Replay-Nonce", "Link", "Location"}
		foundResponseHeaders := []string{}
		for _, param := range resp {
			for _, reqHeader := range requiredResponseHeaders {
				if strings.EqualFold(param, reqHeader) {
					foundResponseHeaders = append(foundResponseHeaders, reqHeader)
					break
				}
			}
		}

		foundAllHeaders := strutil.EquivalentSlices(requiredResponseHeaders, foundResponseHeaders)

		if !foundAllHeaders {
			ret := Result{
				Status:   ResultInformational,
				Endpoint: "/sys/mounts/{{mount}}/tune",
				Message:  "Mount hasn't enabled 'Replay-Nonce', 'Link', 'Location' response headers, these are required for ACME to function.",
			}
			results = append(results, &ret)
		} else {
			ret := Result{
				Status:   ResultOK,
				Endpoint: "/sys/mounts/{{mount}}/tune",
				Message:  "Mount has enabled 'Replay-Nonce', 'Link', 'Location' response headers.",
			}
			results = append(results, &ret)
		}
	}

	// Check for local cluster config
	if h.ClusterConfigFetcher.IsSecretPermissionsError() {
		ret := Result{
			Status:   ResultInsufficientPermissions,
			Endpoint: h.ClusterConfigFetcher.Path,
			Message:  "Without this information, this health check is unable to function.",
		}

		if e.Client.Token() == "" {
			ret.Message = "No token available so unable to read the local cluster config for this mount. " + ret.Message
		} else {
			ret.Message = "This token lacks permission to read the local cluster config for this mount. " + ret.Message
		}

		results = append(results, &ret)
	} else {
		localPathIssue := verifyLocalPathUrl(h)

		if localPathIssue != nil {
			ret := Result{
				Status:   ResultInformational,
				Endpoint: h.ClusterConfigFetcher.Path,
				Message:  localPathIssue.Error(),
			}
			results = append(results, &ret)
		} else {
			ret := Result{
				Status:   ResultOK,
				Endpoint: h.ClusterConfigFetcher.Path,
				Message:  "Local cluster config 'path' argument has an https URL.",
			}
			results = append(results, &ret)
		}
	}

	return
}
