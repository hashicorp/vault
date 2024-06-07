// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package healthcheck

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type AllowAcmeHeaders struct {
	Enabled            bool
	UnsupportedVersion bool

	TuneFetcher *PathFetch
	TuneData    map[string]interface{}

	AcmeConfigFetcher *PathFetch
}

func NewAllowAcmeHeaders() Check {
	return &AllowAcmeHeaders{}
}

func (h *AllowAcmeHeaders) Name() string {
	return "allow_acme_headers"
}

func (h *AllowAcmeHeaders) IsEnabled() bool {
	return h.Enabled
}

func (h *AllowAcmeHeaders) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func (h *AllowAcmeHeaders) LoadConfig(config map[string]interface{}) error {
	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *AllowAcmeHeaders) FetchResources(e *Executor) error {
	var err error
	h.AcmeConfigFetcher, err = e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/config/acme")
	if err != nil {
		return err
	}

	if h.AcmeConfigFetcher.IsUnsupportedPathError() {
		h.UnsupportedVersion = true
	}

	_, h.TuneFetcher, h.TuneData, err = fetchMountTune(e, func() {
		h.UnsupportedVersion = true
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *AllowAcmeHeaders) Evaluate(e *Executor) ([]*Result, error) {
	if h.UnsupportedVersion {
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: h.AcmeConfigFetcher.Path,
			Message:  "This health check requires Vault 1.14+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	if h.AcmeConfigFetcher.IsSecretPermissionsError() {
		msg := "Without read access to ACME configuration, this health check is unable to function."
		return craftInsufficientPermissionResult(e, h.AcmeConfigFetcher.Path, msg), nil
	}

	acmeEnabled, err := isAcmeEnabled(h.AcmeConfigFetcher)
	if err != nil {
		return nil, err
	}

	if !acmeEnabled {
		ret := Result{
			Status:   ResultNotApplicable,
			Endpoint: h.AcmeConfigFetcher.Path,
			Message:  "ACME is not enabled, no additional response headers required.",
		}
		return []*Result{&ret}, nil
	}

	if h.TuneFetcher.IsSecretPermissionsError() {
		msg := "Without access to mount tune information, this health check is unable to function."
		return craftInsufficientPermissionResult(e, h.TuneFetcher.Path, msg), nil
	}

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
			Status:   ResultWarning,
			Endpoint: "/sys/mounts/{{mount}}/tune",
			Message:  "Mount hasn't enabled 'Replay-Nonce', 'Link', 'Location' response headers, these are required for ACME to function.",
		}
		return []*Result{&ret}, nil
	}

	ret := Result{
		Status:   ResultOK,
		Endpoint: "/sys/mounts/{{mount}}/tune",
		Message:  "Mount has enabled 'Replay-Nonce', 'Link', 'Location' response headers.",
	}
	return []*Result{&ret}, nil
}

func craftInsufficientPermissionResult(e *Executor, path, errorMsg string) []*Result {
	ret := Result{
		Status:   ResultInsufficientPermissions,
		Endpoint: path,
		Message:  errorMsg,
	}

	if e.Client.Token() == "" {
		ret.Message = "No token available so unable read the tune endpoint for this mount. " + ret.Message
	} else {
		ret.Message = "This token lacks permission to read the tune endpoint for this mount. " + ret.Message
	}

	return []*Result{&ret}
}
