// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

// Package awsutil provides helpers for generating the AWS IAM login data used to
// authenticate against Vault's AWS auth method with the AWS SDK for Go v2.
//
// The AWS SDK v2 migration removed the GenerateLoginData helper that previously
// lived in github.com/hashicorp/go-secure-stdlib/awsutil (see
// https://github.com/hashicorp/go-secure-stdlib/pull/83). This package re-homes
// that logic, along with the STS endpoint resolution it depends on, so it can be
// shared by the AWS auth backend, the login CLI, and other callers without
// duplication. It is intended to be straightforward to contribute back to
// go-secure-stdlib (or a similar shared library) in the future.
package awsutil

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/go-hclog"
	sdkawsutil "github.com/hashicorp/go-secure-stdlib/awsutil/v2"
)

// iamServerIdHeader is the name of the header that the client signs and the
// server validates so that a signed GetCallerIdentity request cannot be replayed
// against an unintended Vault server. The AWS auth backend defines its own copy
// of this constant for server-side request validation.
const iamServerIdHeader = "X-Vault-AWS-IAM-Server-ID"

// GenerateLoginDataV2 builds the login payload for the AWS IAM auth method using
// the AWS SDK for Go v2. It signs a GetCallerIdentity STS request with the
// credentials resolved from cfg and returns the base64-encoded request
// components expected by the auth/aws login endpoint.
//
// This replaces the v1 awsutil.GenerateLoginData helper that was dropped during
// the SDK v2 upgrade.
func GenerateLoginDataV2(ctx context.Context, cfg *awsv2.Config, headerValue, configuredRegion string, logger hclog.Logger) (map[string]interface{}, error) {
	if cfg == nil {
		return nil, fmt.Errorf("aws config must not be nil")
	}
	if cfg.Credentials == nil {
		return nil, fmt.Errorf("aws config credentials must not be nil")
	}
	loginData := make(map[string]interface{})

	region, err := sdkawsutil.GetRegion(ctx, configuredRegion)
	if err != nil {
		logger.Warn(fmt.Sprintf("defaulting region to %q due to %s", sdkawsutil.DefaultRegion, err.Error()))
		region = sdkawsutil.DefaultRegion
	}

	body := []byte("Action=GetCallerIdentity&Version=2011-06-15")

	stsURL, err := STSLoginEndpoint(ctx, region)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, stsURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	if headerValue != "" {
		req.Header.Set(iamServerIdHeader, headerValue)
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	payloadHash := sha256.Sum256(body)
	signer := v4.NewSigner()
	if err := signer.SignHTTP(ctx, creds, req, hex.EncodeToString(payloadHash[:]), "sts", region, time.Now()); err != nil {
		return nil, fmt.Errorf("failed to sign STS request: %w", err)
	}

	headersJSON, err := json.Marshal(req.Header)
	if err != nil {
		return nil, err
	}

	loginData["iam_http_request_method"] = req.Method
	loginData["iam_request_url"] = base64.StdEncoding.EncodeToString([]byte(req.URL.String()))
	loginData["iam_request_headers"] = base64.StdEncoding.EncodeToString(headersJSON)
	loginData["iam_request_body"] = base64.StdEncoding.EncodeToString(body)

	return loginData, nil
}

// STSLoginEndpoint returns the STS endpoint URL whose host matches the region the
// login request is signed for. The CLI has historically signed against the global
// endpoint, so that is retained for the default region; any other region resolves
// to a regional endpoint so the signed Host matches the request region (this also
// handles non-default partitions such as AWS China and GovCloud).
func STSLoginEndpoint(ctx context.Context, region string) (string, error) {
	if region == sdkawsutil.DefaultRegion {
		return "https://sts.amazonaws.com/", nil
	}
	regional, err := STSRegionalEndpoint(ctx, region)
	if err != nil {
		return "", err
	}
	// The SDK resolver may or may not return a trailing slash; trim it before
	// re-appending so the endpoint never ends up with a double slash.
	return strings.TrimRight(regional, "/") + "/", nil
}

// STSRegionalEndpoint resolves the regional STS endpoint URL for the given region
// using the AWS SDK v2 endpoint resolver, accounting for non-default partitions.
func STSRegionalEndpoint(ctx context.Context, region string) (string, error) {
	resolver := sts.NewDefaultEndpointResolverV2()
	resolvedEndpoint, err := resolver.ResolveEndpoint(ctx, sts.EndpointParameters{
		Region:            awsv2.String(region),
		UseDualStack:      awsv2.Bool(false),
		UseFIPS:           awsv2.Bool(false),
		UseGlobalEndpoint: awsv2.Bool(false),
	})
	if err != nil {
		return "", fmt.Errorf("unable to get regional STS endpoint for region: %v", region)
	}
	return resolvedEndpoint.URI.String(), nil
}
