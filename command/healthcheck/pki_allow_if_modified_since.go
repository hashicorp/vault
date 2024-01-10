// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package healthcheck

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type AllowIfModifiedSince struct {
	Enabled            bool
	UnsupportedVersion bool

	TuneData map[string]interface{}
	Fetcher  *PathFetch
}

func NewAllowIfModifiedSinceCheck() Check {
	return &AllowIfModifiedSince{}
}

func (h *AllowIfModifiedSince) Name() string {
	return "allow_if_modified_since"
}

func (h *AllowIfModifiedSince) IsEnabled() bool {
	return h.Enabled
}

func (h *AllowIfModifiedSince) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func (h *AllowIfModifiedSince) LoadConfig(config map[string]interface{}) error {
	var err error

	h.Enabled, err = parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}

	return nil
}

func (h *AllowIfModifiedSince) FetchResources(e *Executor) error {
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

func (h *AllowIfModifiedSince) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/sys/mounts/{{mount}}/tune",
			Message:  "This health check requires Vault 1.12+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
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

	req, err := StringList(h.TuneData["passthrough_request_headers"])
	if err != nil {
		return nil, fmt.Errorf("unable to parse value from server for passthrough_request_headers: %w", err)
	}

	resp, err := StringList(h.TuneData["allowed_response_headers"])
	if err != nil {
		return nil, fmt.Errorf("unable to parse value from server for allowed_response_headers: %w", err)
	}

	foundIMS := false
	for _, param := range req {
		if strings.EqualFold(param, "If-Modified-Since") {
			foundIMS = true
			break
		}
	}

	foundLM := false
	for _, param := range resp {
		if strings.EqualFold(param, "Last-Modified") {
			foundLM = true
			break
		}
	}

	if !foundIMS || !foundLM {
		ret := Result{
			Status:   ResultInformational,
			Endpoint: "/sys/mounts/{{mount}}/tune",
			Message:  "Mount hasn't enabled If-Modified-Since Request or Last-Modified Response headers; consider enabling these headers to allow clients to fetch CAs and CRLs only when they've changed, reducing total bandwidth.",
		}
		results = append(results, &ret)
	} else {
		ret := Result{
			Status:   ResultOK,
			Endpoint: "/sys/mounts/{{mount}}/tune",
			Message:  "Mount allows the If-Modified-Since request header and Last-Modified response header.",
		}
		results = append(results, &ret)
	}

	return
}
