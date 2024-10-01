// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package healthcheck

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/acme"
)

type EnableAcmeIssuance struct {
	Enabled            bool
	UnsupportedVersion bool

	AcmeConfigFetcher    *PathFetch
	ClusterConfigFetcher *PathFetch
	TotalIssuers         int
	RootIssuers          int
}

func NewEnableAcmeIssuance() Check {
	return &EnableAcmeIssuance{}
}

func (h *EnableAcmeIssuance) Name() string {
	return "enable_acme_issuance"
}

func (h *EnableAcmeIssuance) IsEnabled() bool {
	return h.Enabled
}

func (h *EnableAcmeIssuance) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func (h *EnableAcmeIssuance) LoadConfig(config map[string]interface{}) error {
	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *EnableAcmeIssuance) FetchResources(e *Executor) error {
	var err error
	h.AcmeConfigFetcher, err = e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/config/acme")
	if err != nil {
		return err
	}

	if h.AcmeConfigFetcher.IsUnsupportedPathError() {
		h.UnsupportedVersion = true
	}

	h.ClusterConfigFetcher, err = e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/config/cluster")
	if err != nil {
		return err
	}

	if h.ClusterConfigFetcher.IsUnsupportedPathError() {
		h.UnsupportedVersion = true
	}

	h.TotalIssuers, h.RootIssuers, err = doesMountContainOnlyRootIssuers(e)
	return err
}

func doesMountContainOnlyRootIssuers(e *Executor) (int, int, error) {
	exit, _, issuers, err := pkiFetchIssuersList(e, func() {})
	if exit || err != nil {
		return 0, 0, err
	}

	totalIssuers := 0
	rootIssuers := 0

	for _, issuer := range issuers {
		skip, _, cert, err := pkiFetchIssuer(e, issuer, func() {})

		if skip || err != nil {
			if err != nil {
				return 0, 0, err
			}
			continue
		}
		totalIssuers++

		if !bytes.Equal(cert.RawSubject, cert.RawIssuer) {
			continue
		}
		if err := cert.CheckSignatureFrom(cert); err != nil {
			continue
		}
		rootIssuers++
	}

	return totalIssuers, rootIssuers, nil
}

func isAcmeEnabled(fetcher *PathFetch) (bool, error) {
	isEnabledRaw, ok := fetcher.Secret.Data["enabled"]
	if !ok {
		return false, fmt.Errorf("enabled configuration field missing from acme config")
	}

	parseBool, err := parseutil.ParseBool(isEnabledRaw)
	if err != nil {
		return false, fmt.Errorf("failed parsing 'enabled' field from ACME config: %w", err)
	}

	return parseBool, nil
}

func verifyLocalPathUrl(h *EnableAcmeIssuance) error {
	localPathRaw, ok := h.ClusterConfigFetcher.Secret.Data["path"]
	if !ok {
		return fmt.Errorf("'path' field missing from config")
	}

	localPath, err := parseutil.ParseString(localPathRaw)
	if err != nil {
		return fmt.Errorf("failed converting 'path' field from local config: %w", err)
	}

	if localPath == "" {
		return fmt.Errorf("'path' field not configured within /{{mount}}/config/cluster")
	}

	parsedUrl, err := url.Parse(localPath)
	if err != nil {
		return fmt.Errorf("failed to parse URL from path config: %v: %w", localPathRaw, err)
	}

	if parsedUrl.Scheme != "https" {
		return fmt.Errorf("the configured 'path' field in /{{mount}}/config/cluster was not using an https scheme")
	}

	// Avoid issues with SSL certificates for this check, we just want to validate that we would
	// hit an ACME server with the path they specified in configuration
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	acmeDirectoryUrl := parsedUrl.JoinPath("/acme/", "directory")
	acmeClient := acme.Client{HTTPClient: client, DirectoryURL: acmeDirectoryUrl.String()}
	_, err = acmeClient.Discover(context.Background())
	if err != nil {
		return fmt.Errorf("using configured 'path' field ('%s') in /{{mount}}/config/cluster failed to reach the ACME"+
			" directory: %s: %w", parsedUrl.String(), acmeDirectoryUrl.String(), err)
	}

	return nil
}

func (h *EnableAcmeIssuance) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: h.AcmeConfigFetcher.Path,
			Message:  "This health check requires Vault 1.14+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	if h.AcmeConfigFetcher.IsSecretPermissionsError() {
		msg := "Without this information, this health check is unable to function."
		return craftInsufficientPermissionResult(e, h.AcmeConfigFetcher.Path, msg), nil
	}

	acmeEnabled, err := isAcmeEnabled(h.AcmeConfigFetcher)
	if err != nil {
		return nil, err
	}

	if !acmeEnabled {
		if h.TotalIssuers == 0 {
			ret := Result{
				Status:   ResultNotApplicable,
				Endpoint: h.AcmeConfigFetcher.Path,
				Message:  "No issuers in mount, ACME is not required.",
			}
			return []*Result{&ret}, nil
		}

		if h.TotalIssuers == h.RootIssuers {
			ret := Result{
				Status:   ResultNotApplicable,
				Endpoint: h.AcmeConfigFetcher.Path,
				Message:  "Mount contains only root issuers, ACME is not required.",
			}
			return []*Result{&ret}, nil
		}

		ret := Result{
			Status:   ResultInformational,
			Endpoint: h.AcmeConfigFetcher.Path,
			Message:  "Consider enabling ACME support to support a self-rotating PKI infrastructure.",
		}
		return []*Result{&ret}, nil
	}

	if h.ClusterConfigFetcher.IsSecretPermissionsError() {
		msg := "Without this information, this health check is unable to function."
		return craftInsufficientPermissionResult(e, h.ClusterConfigFetcher.Path, msg), nil
	}

	localPathIssue := verifyLocalPathUrl(h)

	if localPathIssue != nil {
		ret := Result{
			Status:   ResultWarning,
			Endpoint: h.ClusterConfigFetcher.Path,
			Message:  "ACME enabled in config but not functional: " + localPathIssue.Error(),
		}
		return []*Result{&ret}, nil
	}

	ret := Result{
		Status:   ResultOK,
		Endpoint: h.ClusterConfigFetcher.Path,
		Message:  "ACME enabled and successfully connected to the ACME directory.",
	}
	return []*Result{&ret}, nil
}
