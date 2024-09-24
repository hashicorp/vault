// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type acmeContext struct {
	issuing.IssuerRoleContext

	// baseUrl is the combination of the configured cluster local URL and the acmePath up to /acme/
	baseUrl    *url.URL
	clusterUrl *url.URL
	sc         *storageContext
	acmeState  *acmeState
	// acmeDirectory is a string that can distinguish the various acme directories we have configured
	// if something needs to remain locked into a directory path structure.
	acmeDirectory string
	eabPolicy     EabPolicy
	ciepsPolicy   string
	runtimeOpts   acmeWrapperOpts
}

func (c acmeContext) getAcmeState() *acmeState {
	return c.acmeState
}

type (
	acmeOperation                func(acmeCtx *acmeContext, r *logical.Request, _ *framework.FieldData) (*logical.Response, error)
	acmeParsedOperation          func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}) (*logical.Response, error)
	acmeAccountRequiredOperation func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}, acct *acmeAccount) (*logical.Response, error)
)

// setupAcmeDirectory will populate a prefix'd URL with all the paths required
// for a given ACME directory.
func setupAcmeDirectory(b *backend, acmePrefix string, unauthPrefix string, opts acmeWrapperOpts) {
	acmePrefix = strings.TrimRight(acmePrefix, "/")
	unauthPrefix = strings.TrimRight(unauthPrefix, "/")

	b.Backend.Paths = append(b.Backend.Paths, pathAcmeDirectory(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeNonce(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeNewAccount(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeUpdateAccount(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeGetOrder(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeListOrders(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeNewOrder(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeFinalizeOrder(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeFetchOrderCert(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeChallenge(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeAuthorization(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeRevoke(b, acmePrefix, opts))
	b.Backend.Paths = append(b.Backend.Paths, pathAcmeNewEab(b, acmePrefix)) // auth'd API that lives underneath the various /acme paths

	// Add specific un-auth'd paths for ACME APIs
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/directory")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/new-nonce")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/new-account")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/new-order")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/revoke-cert")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/key-change")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/account/+")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/authorization/+")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/challenge/+/+")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/orders")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/order/+")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/order/+/finalize")
	b.PathsSpecial.Unauthenticated = append(b.PathsSpecial.Unauthenticated, unauthPrefix+"/order/+/cert")
	// We specifically do NOT add acme/new-eab to this as it should be auth'd
}

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

type acmeWrapperOpts struct {
	isDefault      bool
	isCiepsEnabled bool
}

func (o acmeWrapperOpts) Clone() acmeWrapperOpts {
	return acmeWrapperOpts{
		isDefault:      o.isDefault,
		isCiepsEnabled: o.isCiepsEnabled,
	}
}

// acmeWrapper a basic wrapper that all ACME handlers should leverage as the basis.
// This will create a basic ACME context, validate basic ACME configuration is setup
// for operations. This pulls in acmeErrorWrapper to translate error messages for users,
// but does not enforce any sort of ACME authentication.
func (b *backend) acmeWrapper(opts acmeWrapperOpts, op acmeOperation) framework.OperationFunc {
	return acmeErrorWrapper(func(ctx context.Context, r *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		sc := b.makeStorageContext(ctx, r.Storage)

		config, err := b.GetAcmeState().getConfigWithUpdate(sc)
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

		if b.UseLegacyBundleCaStorage() {
			return nil, fmt.Errorf("%w: Can not perform ACME operations until migration has completed", ErrServerInternal)
		}

		acmeBaseUrl, clusterBase, err := getAcmeBaseUrl(sc, r)
		if err != nil {
			return nil, err
		}

		role, issuer, err := getAcmeRoleAndIssuer(sc, data, config)
		if err != nil {
			return nil, err
		}

		acmeDirectory, err := getAcmeDirectory(r)
		if err != nil {
			return nil, err
		}

		isCiepsEnabled, ciepsPolicy, err := getCiepsAcmeSettings(b, sc, opts, config, data)
		if err != nil {
			return nil, err
		}
		runtimeOpts := opts.Clone()

		if isCiepsEnabled {
			// We need to possibly reset the isCiepsEnabled option to true if we are in
			// the default folder with the external-policy set as it would have been
			// normally disabled.
			if runtimeOpts.isDefault {
				runtimeOpts.isCiepsEnabled = true
			}
		}

		acmeCtx := &acmeContext{
			IssuerRoleContext: issuing.NewIssuerRoleContext(ctx, issuer, role),
			baseUrl:           acmeBaseUrl,
			clusterUrl:        clusterBase,
			sc:                sc,
			acmeState:         b.acmeState,
			acmeDirectory:     acmeDirectory,
			eabPolicy:         eabPolicy,
			ciepsPolicy:       ciepsPolicy,
			runtimeOpts:       runtimeOpts,
		}

		return op(acmeCtx, r, data)
	})
}

// acmeParsedWrapper is an ACME wrapper that will parse out the ACME request parameters, validate
// that we have a proper signature and pass to the operation a decoded map of arguments received.
// This wrapper builds on top of acmeWrapper. Note that this does perform signature verification
// it does not enforce the account being in a valid state nor existing.
func (b *backend) acmeParsedWrapper(opt acmeWrapperOpts, op acmeParsedOperation) framework.OperationFunc {
	return b.acmeWrapper(opt, func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData) (*logical.Response, error) {
		user, data, err := b.GetAcmeState().ParseRequestParams(acmeCtx, r, fields)
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
				nonce, _, err := b.GetAcmeState().GetNonce()
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
func (b *backend) acmeAccountRequiredWrapper(opt acmeWrapperOpts, op acmeAccountRequiredOperation) framework.OperationFunc {
	return b.acmeParsedWrapper(opt, func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, uc *jwsCtx, data map[string]interface{}) (*logical.Response, error) {
		if !uc.Existing {
			return nil, fmt.Errorf("cannot process request without a 'kid': %w", ErrMalformed)
		}

		account, err := requireValidAcmeAccount(acmeCtx, uc)
		if err != nil {
			return nil, err
		}

		return op(acmeCtx, r, fields, uc, data, account)
	})
}

func requireValidAcmeAccount(acmeCtx *acmeContext, uc *jwsCtx) (*acmeAccount, error) {
	account, err := acmeCtx.getAcmeState().LoadAccount(acmeCtx, uc.Kid)
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
	return account, nil
}

func getAcmeBaseUrl(sc *storageContext, r *logical.Request) (*url.URL, *url.URL, error) {
	baseUrl, err := getBasePathFromClusterConfig(sc)
	if err != nil {
		return nil, nil, err
	}

	directoryPrefix, err := getAcmeDirectory(r)
	if err != nil {
		return nil, nil, err
	}

	return baseUrl.JoinPath(directoryPrefix), baseUrl, nil
}

func getBasePathFromClusterConfig(sc *storageContext) (*url.URL, error) {
	cfg, err := sc.getClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed loading cluster config: %w", err)
	}

	if cfg.Path == "" {
		return nil, fmt.Errorf("ACME feature requires local cluster 'path' field configuration to be set")
	}

	baseUrl, err := url.Parse(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed parsing URL configured in local cluster 'path' configuration: %s: %s",
			cfg.Path, err.Error())
	}
	return baseUrl, nil
}

func getAcmeIssuer(sc *storageContext, issuerName string) (*issuing.IssuerEntry, error) {
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

	if issuer.Usage.HasUsage(issuing.IssuanceUsage) && len(issuer.KeyID) > 0 {
		return issuer, nil
	}

	return nil, fmt.Errorf("%w: issuer missing proper issuance usage or key", ErrServerInternal)
}

// getAcmeDirectory return the base acme directory path, without a leading '/' and including
// the trailing /acme/ folder which is the root of all our various directories
func getAcmeDirectory(r *logical.Request) (string, error) {
	acmePath := r.Path
	if !strings.HasPrefix(acmePath, "/") {
		acmePath = "/" + acmePath
	}

	lastIndex := strings.LastIndex(acmePath, "/acme/")
	if lastIndex == -1 {
		return "", fmt.Errorf("%w: unable to determine acme base folder path: %s", ErrServerInternal, acmePath)
	}

	// Skip the leading '/' and return our base path with the /acme/
	return strings.TrimLeft(acmePath[0:lastIndex]+"/acme/", "/"), nil
}

func getAcmeRoleAndIssuer(sc *storageContext, data *framework.FieldData, config *acmeConfigEntry) (*issuing.RoleEntry, *issuing.IssuerEntry, error) {
	requestedIssuer := getRequestedAcmeIssuerFromPath(data)
	requestedRole := getRequestedAcmeRoleFromPath(data)
	issuerToLoad := requestedIssuer

	var role *issuing.RoleEntry
	var err error

	if len(requestedRole) == 0 { // Default Directory
		policyType, extraInfo, err := getDefaultDirectoryPolicyType(config.DefaultDirectoryPolicy)
		if err != nil {
			return nil, nil, err
		}
		switch policyType {
		case Forbid:
			return nil, nil, fmt.Errorf("%w: default directory not allowed by ACME policy", ErrServerInternal)
		case SignVerbatim, ExternalPolicy:
			role = issuing.SignVerbatimRoleWithOpts(
				issuing.WithIssuer(requestedIssuer),
				issuing.WithNoStore(false))
		case Role:
			role, err = getAndValidateAcmeRole(sc, extraInfo)
			if err != nil {
				return nil, nil, err
			}
		}
	} else { // Requested Role
		role, err = getAndValidateAcmeRole(sc, requestedRole)
		if err != nil {
			return nil, nil, err
		}

		// Check the Requested Role is Allowed
		allowAnyRole := len(config.AllowedRoles) == 1 && config.AllowedRoles[0] == "*"
		if !allowAnyRole {

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
	}

	// If we haven't loaded an issuer directly from our path and the specified (or default)
	// role does specify an issuer prefer the role's issuer rather than the default issuer.
	if len(role.Issuer) > 0 && len(requestedIssuer) == 0 {
		issuerToLoad = role.Issuer
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

	// If not allowed in configuration, override ExtKeyUsage behavior to force it to only be
	// ServerAuth within ACME issued certs
	if !config.AllowRoleExtKeyUsage {
		role.ExtKeyUsage = []string{"serverauth"}
		role.ExtKeyUsageOIDs = []string{}
		role.ServerFlag = true
		role.ClientFlag = false
		role.CodeSigningFlag = false
		role.EmailProtectionFlag = false
	}

	return role, issuer, nil
}

func getAndValidateAcmeRole(sc *storageContext, requestedRole string) (*issuing.RoleEntry, error) {
	var err error
	role, err := sc.GetRole(requestedRole)
	if err != nil {
		return nil, fmt.Errorf("%w: err loading role", ErrServerInternal)
	}

	if role == nil {
		return nil, fmt.Errorf("%w: role does not exist", ErrMalformed)
	}

	if role.NoStore {
		return nil, fmt.Errorf("%w: role can not be used as NoStore is set to true", ErrServerInternal)
	}

	return role, nil
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
		sc.Logger().Warn(fmt.Sprintf("could not parse env var '%s'", disableAcmeEnvVar), "error", nonFatalErr)
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
