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
	baseUrl    *url.URL
	clusterUrl *url.URL
	sc         *storageContext
	role       *roleEntry
	issuer     *issuerEntry
	// acmeDirectory is a string that can distinguish the various acme directories we have configured
	// if something needs to remain locked into a directory path structure.
	acmeDirectory string
	eabPolicy     EabPolicy
}

type (
	acmeOperation                func(acmeCtx *acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error)
	acmeParsedOperation          func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}) (*logical.Response, error)
	acmeAccountRequiredOperation func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}, acct *acmeAccount) (*logical.Response, error)
)

// acmeErrorWrapper the lowest level wrapper that will translate errors into proper ACME error responses
func acmeErrorWrapper(op framework.OperationFunc) framework.OperationFunc {
	return func(ctx context.Context, r *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		resp, err := op(ctx, r, data)
		if err != nil {
			return TranslateError(err)
		}

		return resp, nil
	}
}

// acmeWrapper a basic wrapper that all ACME handlers should leverage as the basis.
// This will create a basic ACME context, validate basic ACME configuration is setup
// for operations. This pulls in acmeErrorWrapper to translate error messages for users,
// but does not enforce any sort of ACME authentication.
func (b *backend) acmeWrapper(op acmeOperation) framework.OperationFunc {
	return acmeErrorWrapper(func(ctx context.Context, r *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		sc := b.makeStorageContext(ctx, r.Storage)

		config, err := sc.Backend.acmeState.getConfigWithUpdate(sc)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch ACME configuration: %w", err)
		}

		// use string form in case someone messes up our config from raw storage.
		eabPolicy, err := getEabPolicyByString(string(config.EabPolicyName))
		if err != nil {
			return nil, err
		}

		if isAcmeDisabled(sc, config, eabPolicy) {
			return nil, ErrAcmeDisabled
		}

		if b.useLegacyBundleCaStorage() {
			return nil, fmt.Errorf("%w: Can not perform ACME operations until migration has completed", ErrServerInternal)
		}

		acmeBaseUrl, clusterBase, err := getAcmeBaseUrl(sc, r.Path)
		if err != nil {
			return nil, err
		}

		role, issuer, err := getAcmeRoleAndIssuer(sc, data, config)
		if err != nil {
			return nil, err
		}

		acmeDirectory := getAcmeDirectory(data)

		acmeCtx := &acmeContext{
			baseUrl:       acmeBaseUrl,
			clusterUrl:    clusterBase,
			sc:            sc,
			role:          role,
			issuer:        issuer,
			acmeDirectory: acmeDirectory,
			eabPolicy:     eabPolicy,
		}

		return op(acmeCtx, r, data)
	})
}

// acmeParsedWrapper is an ACME wrapper that will parse out the ACME request parameters, validate
// that we have a proper signature and pass to the operation a decoded map of arguments received.
// This wrapper builds on top of acmeWrapper. Note that this does perform signature verification
// it does not enforce the account being in a valid state nor existing.
func (b *backend) acmeParsedWrapper(op acmeParsedOperation) framework.OperationFunc {
	return b.acmeWrapper(func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData) (*logical.Response, error) {
		user, data, err := b.acmeState.ParseRequestParams(acmeCtx, r, fields)
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

// acmeAccountRequiredWrapper builds on top of acmeParsedWrapper, enforcing the
// request has a proper signature for an existing account, and that account is
// in a valid status. It passes to the operation a decoded form of the request
// parameters as well as the ACME account the request is for.
func (b *backend) acmeAccountRequiredWrapper(op acmeAccountRequiredOperation) framework.OperationFunc {
	return b.acmeParsedWrapper(func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, uc *jwsCtx, data map[string]interface{}) (*logical.Response, error) {
		if !uc.Existing {
			return nil, fmt.Errorf("cannot process request without a 'kid': %w", ErrMalformed)
		}

		account, err := b.acmeState.LoadAccount(acmeCtx, uc.Kid)
		if err != nil {
			return nil, fmt.Errorf("error loading account: %w", err)
		}

		if err = acmeCtx.eabPolicy.EnforceForExistingAccount(account); err != nil {
			return nil, err
		}

		if account.Status != AccountStatusValid {
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

func getAcmeBaseUrl(sc *storageContext, path string) (*url.URL, *url.URL, error) {
	cfg, err := sc.getClusterConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed loading cluster config: %w", err)
	}

	if cfg.Path == "" {
		return nil, nil, fmt.Errorf("ACME feature requires local cluster path configuration to be set: %w", ErrServerInternal)
	}

	baseUrl, err := url.Parse(cfg.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("ACME feature a proper URL configured in local cluster path: %w", ErrServerInternal)
	}

	directoryPrefix := ""
	lastIndex := strings.LastIndex(path, "/acme/")
	if lastIndex != -1 {
		directoryPrefix = path[0:lastIndex]
	}

	return baseUrl.JoinPath(directoryPrefix, "/acme/"), baseUrl, nil
}

func getAcmeIssuer(sc *storageContext, issuerName string) (*issuerEntry, error) {
	if issuerName == "" {
		issuerName = defaultRef
	}
	issuerId, err := sc.resolveIssuerReference(issuerName)
	if err != nil {
		return nil, fmt.Errorf("%w: issuer does not exist", ErrMalformed)
	}

	issuer, err := sc.fetchIssuerById(issuerId)
	if err != nil {
		return nil, fmt.Errorf("issuer failed to load: %w", err)
	}

	if issuer.Usage.HasUsage(IssuanceUsage) && len(issuer.KeyID) > 0 {
		return issuer, nil
	}

	return nil, fmt.Errorf("%w: issuer missing proper issuance usage or key", ErrServerInternal)
}

func getAcmeDirectory(data *framework.FieldData) string {
	requestedIssuer := getRequestedAcmeIssuerFromPath(data)
	requestedRole := getRequestedAcmeRoleFromPath(data)

	return fmt.Sprintf("issuer-%s::role-%s", requestedIssuer, requestedRole)
}

func getAcmeRoleAndIssuer(sc *storageContext, data *framework.FieldData, config *acmeConfigEntry) (*roleEntry, *issuerEntry, error) {
	requestedIssuer := getRequestedAcmeIssuerFromPath(data)
	requestedRole := getRequestedAcmeRoleFromPath(data)
	issuerToLoad := requestedIssuer

	var wasVerbatim bool
	var role *roleEntry

	if len(requestedRole) > 0 || len(config.DefaultRole) > 0 {
		if len(requestedRole) == 0 {
			requestedRole = config.DefaultRole
		}

		var err error
		role, err = sc.Backend.getRole(sc.Context, sc.Storage, requestedRole)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: err loading role", ErrServerInternal)
		}

		if role == nil {
			return nil, nil, fmt.Errorf("%w: role does not exist", ErrMalformed)
		}

		if role.NoStore {
			return nil, nil, fmt.Errorf("%w: role can not be used as NoStore is set to true", ErrServerInternal)
		}

		// If we haven't loaded an issuer directly from our path and the specified
		// role does specify an issuer prefer the role's issuer rather than the default issuer.
		if len(role.Issuer) > 0 && len(requestedIssuer) == 0 {
			issuerToLoad = role.Issuer
		}
	} else {
		role = buildSignVerbatimRoleWithNoData(&roleEntry{
			Issuer:  requestedIssuer,
			NoStore: false,
			Name:    requestedRole,
		})
		wasVerbatim = true
	}

	allowAnyRole := len(config.AllowedRoles) == 1 && config.AllowedRoles[0] == "*"
	if !allowAnyRole {
		if wasVerbatim {
			return nil, nil, fmt.Errorf("%w: using the default directory without specifying a role is not supported by this configuration; specify 'default_role' in the acme config to the default directories", ErrServerInternal)
		}

		var foundRole bool
		for _, name := range config.AllowedRoles {
			if name == role.Name {
				foundRole = true
				break
			}
		}

		if !foundRole {
			return nil, nil, fmt.Errorf("%w: specified role not allowed by ACME policy", ErrServerInternal)
		}
	}

	issuer, err := getAcmeIssuer(sc, issuerToLoad)
	if err != nil {
		return nil, nil, err
	}

	allowAnyIssuer := len(config.AllowedIssuers) == 1 && config.AllowedIssuers[0] == "*"
	if !allowAnyIssuer {
		var foundIssuer bool
		for index, name := range config.AllowedIssuers {
			candidateId, err := sc.resolveIssuerReference(name)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to resolve reference for allowed_issuer entry %d: %w", index, err)
			}

			if candidateId == issuer.ID {
				foundIssuer = true
				break
			}
		}

		if !foundIssuer {
			return nil, nil, fmt.Errorf("%w: specified issuer not allowed by ACME policy", ErrServerInternal)
		}
	}

	return role, issuer, nil
}

func getRequestedAcmeRoleFromPath(data *framework.FieldData) string {
	requestedRole := ""
	roleNameRaw, present := data.GetOk("role")
	if present {
		requestedRole = roleNameRaw.(string)
	}
	return requestedRole
}

func getRequestedAcmeIssuerFromPath(data *framework.FieldData) string {
	requestedIssuer := ""
	requestedIssuerRaw, present := data.GetOk(issuerRefParam)
	if present {
		requestedIssuer = requestedIssuerRaw.(string)
	}
	return requestedIssuer
}

func isAcmeDisabled(sc *storageContext, config *acmeConfigEntry, policy EabPolicy) bool {
	if !config.Enabled {
		return true
	}

	disableAcme, nonFatalErr := isPublicACMEDisabledByEnv()
	if nonFatalErr != nil {
		sc.Backend.Logger().Warn(fmt.Sprintf("could not parse env var '%s'", disableAcmeEnvVar), "error", nonFatalErr)
	}

	// The OS environment if true will override any configuration option.
	if disableAcme {
		if policy.OverrideEnvDisablingPublicAcme() {
			return false
		}
		return true
	}

	return false
}
