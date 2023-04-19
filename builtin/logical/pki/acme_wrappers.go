// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type acmeContext struct {
	// baseUrl is the combination of the configured cluster local URL and the acmePath up to /acme/
	baseUrl *url.URL
	sc      *storageContext
}

type (
	acmeOperation                func(acmeCtx *acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error)
	acmeParsedOperation          func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}) (*logical.Response, error)
	acmeAccountRequiredOperation func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}, acct *acmeAccount) (*logical.Response, error)
)

func acmeErrorWrapper(op framework.OperationFunc) framework.OperationFunc {
	return func(ctx context.Context, r *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		resp, err := op(ctx, r, data)
		if err != nil {
			return TranslateError(err)
		}

		return resp, nil
	}
}

func (b *backend) acmeWrapper(op acmeOperation) framework.OperationFunc {
	return acmeErrorWrapper(func(ctx context.Context, r *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		sc := b.makeStorageContext(ctx, r.Storage)

		if false {
			// TODO sclark: Check if ACME is enable here
			return nil, fmt.Errorf("ACME is disabled in configuration: %w", ErrServerInternal)
		}

		acmeBaseUrl, err := getAcmeBaseUrl(sc, r.Path)
		if err != nil {
			return nil, err
		}

		acmeCtx := &acmeContext{
			baseUrl: acmeBaseUrl,
			sc:      sc,
		}

		return op(acmeCtx, r, data)
	})
}

func (b *backend) acmeParsedWrapper(op acmeParsedOperation) framework.OperationFunc {
	return b.acmeWrapper(func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData) (*logical.Response, error) {
		user, data, err := b.acmeState.ParseRequestParams(acmeCtx, fields)
		if err != nil {
			return nil, err
		}

		resp, err := op(acmeCtx, r, fields, user, data)

		// Our response handlers might not add the necessary headers.
		if resp != nil {
			if resp.Headers == nil {
				resp.Headers = map[string][]string{}
			}

			if _, ok := resp.Headers["Replay-Nonce"]; !ok {
				nonce, _, err := b.acmeState.GetNonce()
				if err != nil {
					return nil, err
				}

				resp.Headers["Replay-Nonce"] = []string{nonce}
			}

			if _, ok := resp.Headers["Link"]; !ok {
				resp.Headers["Link"] = genAcmeLinkHeader(acmeCtx)
			} else {
				directory := genAcmeLinkHeader(acmeCtx)[0]
				addDirectory := true
				for _, item := range resp.Headers["Link"] {
					if item == directory {
						addDirectory = false
						break
					}
				}
				if addDirectory {
					resp.Headers["Link"] = append(resp.Headers["Link"], directory)
				}
			}

			// ACME responses don't understand Vault's default encoding
			// format. Rather than expecting everything to handle creating
			// ACME-formatted responses, do the marshaling in one place.
			if _, ok := resp.Data[logical.HTTPRawBody]; !ok {
				ignored_values := map[string]bool{logical.HTTPContentType: true, logical.HTTPStatusCode: true}
				fields := map[string]interface{}{}
				body := map[string]interface{}{
					logical.HTTPContentType: "application/json",
					logical.HTTPStatusCode:  http.StatusOK,
				}

				for key, value := range resp.Data {
					if _, present := ignored_values[key]; !present {
						fields[key] = value
					} else {
						body[key] = value
					}
				}

				rawBody, err := json.Marshal(fields)
				if err != nil {
					return nil, fmt.Errorf("Error marshaling JSON body: %w", err)
				}

				body[logical.HTTPRawBody] = rawBody
				resp.Data = body
			}
		}

		return resp, err
	})
}

func (b *backend) acmeAccountRequiredWrapper(op acmeAccountRequiredOperation) framework.OperationFunc {
	return b.acmeParsedWrapper(func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, uc *jwsCtx, data map[string]interface{}) (*logical.Response, error) {
		if !uc.Existing {
			return nil, fmt.Errorf("cannot process request without a 'kid': %w", ErrMalformed)
		}

		account, err := b.acmeState.LoadAccount(acmeCtx, uc.Kid)
		if err != nil {
			return nil, fmt.Errorf("error loading account: %w", err)
		}

		if account.Status != StatusValid {
			// Treating "revoked" and "deactivated" as the same here.
			return nil, fmt.Errorf("%w: account in status: %s", ErrUnauthorized, account.Status)
		}

		return op(acmeCtx, r, fields, uc, data, account)
	})
}

// A helper function that will build up the various path patterns we want for ACME APIs.
func buildAcmeFrameworkPaths(b *backend, patternFunc func(b *backend, pattern string) *framework.Path, acmeApi string) []*framework.Path {
	var patterns []*framework.Path
	for _, baseUrl := range []string{
		"acme",
		"roles/" + framework.GenericNameRegex("role") + "/acme",
		"issuer/" + framework.GenericNameRegex(issuerRefParam) + "/acme",
		"issuer/" + framework.GenericNameRegex(issuerRefParam) + "/roles/" + framework.GenericNameRegex("role") + "/acme",
	} {

		if !strings.HasPrefix(acmeApi, "/") {
			acmeApi = "/" + acmeApi
		}

		path := patternFunc(b, baseUrl+acmeApi)
		patterns = append(patterns, path)
	}

	return patterns
}

func getAcmeBaseUrl(sc *storageContext, path string) (*url.URL, error) {
	cfg, err := sc.getClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed loading cluster config: %w", err)
	}

	if cfg.Path == "" {
		return nil, fmt.Errorf("ACME feature requires local cluster path configuration to be set: %w", ErrServerInternal)
	}

	baseUrl, err := url.Parse(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("ACME feature a proper URL configured in local cluster path: %w", ErrServerInternal)
	}

	directoryPrefix := ""
	lastIndex := strings.LastIndex(path, "/acme/")
	if lastIndex != -1 {
		directoryPrefix = path[0:lastIndex]
	}

	return baseUrl.JoinPath(directoryPrefix, "/acme/"), nil
}
