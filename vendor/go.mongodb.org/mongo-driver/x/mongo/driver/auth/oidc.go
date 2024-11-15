// Copyright (C) MongoDB, Inc. 2024-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

// MongoDBOIDC is the string constant for the MONGODB-OIDC authentication mechanism.
const MongoDBOIDC = "MONGODB-OIDC"

// EnvironmentProp is the property key name that specifies the environment for the OIDC authenticator.
const EnvironmentProp = "ENVIRONMENT"

// ResourceProp is the property key name that specifies the token resource for GCP and AZURE OIDC auth.
const ResourceProp = "TOKEN_RESOURCE"

// AllowedHostsProp is the property key name that specifies the allowed hosts for the OIDC authenticator.
const AllowedHostsProp = "ALLOWED_HOSTS"

// AzureEnvironmentValue is the value for the Azure environment.
const AzureEnvironmentValue = "azure"

// GCPEnvironmentValue is the value for the GCP environment.
const GCPEnvironmentValue = "gcp"

// TestEnvironmentValue is the value for the test environment.
const TestEnvironmentValue = "test"

const apiVersion = 1
const invalidateSleepTimeout = 100 * time.Millisecond

// The CSOT specification says to apply a 1-minute timeout if "CSOT is not applied". That's
// ambiguous for the v1.x Go Driver because it could mean either "no timeout provided" or "CSOT not
// enabled". Always use a maximum timeout duration of 1 minute, allowing us to ignore the ambiguity.
// Contexts with a shorter timeout are unaffected.
const machineCallbackTimeout = time.Minute
const humanCallbackTimeout = 5 * time.Minute

var defaultAllowedHosts = []*regexp.Regexp{
	regexp.MustCompile(`^.*[.]mongodb[.]net(:\d+)?$`),
	regexp.MustCompile(`^.*[.]mongodb-qa[.]net(:\d+)?$`),
	regexp.MustCompile(`^.*[.]mongodb-dev[.]net(:\d+)?$`),
	regexp.MustCompile(`^.*[.]mongodbgov[.]net(:\d+)?$`),
	regexp.MustCompile(`^localhost(:\d+)?$`),
	regexp.MustCompile(`^127[.]0[.]0[.]1(:\d+)?$`),
	regexp.MustCompile(`^::1(:\d+)?$`),
}

// OIDCCallback is a function that takes a context and OIDCArgs and returns an OIDCCredential.
type OIDCCallback = driver.OIDCCallback

// OIDCArgs contains the arguments for the OIDC callback.
type OIDCArgs = driver.OIDCArgs

// OIDCCredential contains the access token and refresh token.
type OIDCCredential = driver.OIDCCredential

// IDPInfo contains the information needed to perform OIDC authentication with an Identity Provider.
type IDPInfo = driver.IDPInfo

var _ driver.Authenticator = (*OIDCAuthenticator)(nil)
var _ SpeculativeAuthenticator = (*OIDCAuthenticator)(nil)
var _ SaslClient = (*oidcOneStep)(nil)
var _ SaslClient = (*oidcTwoStep)(nil)

// OIDCAuthenticator is synchronized and handles caching of the access token, refreshToken,
// and IDPInfo. It also provides a mechanism to refresh the access token, but this functionality
// is only for the OIDC Human flow.
type OIDCAuthenticator struct {
	mu sync.Mutex // Guards all of the info in the OIDCAuthenticator struct.

	AuthMechanismProperties map[string]string
	OIDCMachineCallback     OIDCCallback
	OIDCHumanCallback       OIDCCallback

	allowedHosts *[]*regexp.Regexp
	userName     string
	httpClient   *http.Client
	accessToken  string
	refreshToken *string
	idpInfo      *IDPInfo
	tokenGenID   uint64
}

// SetAccessToken allows for manually setting the access token for the OIDCAuthenticator, this is
// only for testing purposes.
func (oa *OIDCAuthenticator) SetAccessToken(accessToken string) {
	oa.mu.Lock()
	defer oa.mu.Unlock()
	oa.accessToken = accessToken
}

func newOIDCAuthenticator(cred *Cred, httpClient *http.Client) (Authenticator, error) {
	if cred.Source != "" && cred.Source != sourceExternal {
		return nil, newAuthError("MONGODB-OIDC source must be empty or $external", nil)
	}
	if cred.Password != "" {
		return nil, fmt.Errorf("password cannot be specified for %q", MongoDBOIDC)
	}
	if cred.Props != nil {
		if env, ok := cred.Props[EnvironmentProp]; ok {
			switch strings.ToLower(env) {
			case AzureEnvironmentValue:
				fallthrough
			case GCPEnvironmentValue:
				if _, ok := cred.Props[ResourceProp]; !ok {
					return nil, fmt.Errorf("%q must be specified for %q %q", ResourceProp, env, EnvironmentProp)
				}
				fallthrough
			case TestEnvironmentValue:
				if cred.OIDCMachineCallback != nil || cred.OIDCHumanCallback != nil {
					return nil, fmt.Errorf("OIDC callbacks are not allowed for %q %q", env, EnvironmentProp)
				}
			}
		}
	}
	oa := &OIDCAuthenticator{
		userName:                cred.Username,
		httpClient:              httpClient,
		AuthMechanismProperties: cred.Props,
		OIDCMachineCallback:     cred.OIDCMachineCallback,
		OIDCHumanCallback:       cred.OIDCHumanCallback,
	}
	err := oa.setAllowedHosts()
	return oa, err
}

func createPatternsForGlobs(hosts []string) ([]*regexp.Regexp, error) {
	var err error
	ret := make([]*regexp.Regexp, len(hosts))
	for i := range hosts {
		hosts[i] = strings.ReplaceAll(hosts[i], ".", "[.]")
		hosts[i] = strings.ReplaceAll(hosts[i], "*", ".*")
		hosts[i] = "^" + hosts[i] + "(:\\d+)?$"
		ret[i], err = regexp.Compile(hosts[i])
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (oa *OIDCAuthenticator) setAllowedHosts() error {
	if oa.AuthMechanismProperties == nil {
		oa.allowedHosts = &defaultAllowedHosts
		return nil
	}

	allowedHosts, ok := oa.AuthMechanismProperties[AllowedHostsProp]
	if !ok {
		oa.allowedHosts = &defaultAllowedHosts
		return nil
	}
	globs := strings.Split(allowedHosts, ",")
	ret, err := createPatternsForGlobs(globs)
	if err != nil {
		return err
	}
	oa.allowedHosts = &ret
	return nil
}

func (oa *OIDCAuthenticator) validateConnectionAddressWithAllowedHosts(conn driver.Connection) error {
	if oa.allowedHosts == nil {
		// should be unreachable, but this is a safety check.
		return newAuthError(fmt.Sprintf("%q missing", AllowedHostsProp), nil)
	}
	allowedHosts := *oa.allowedHosts
	if len(allowedHosts) == 0 {
		return newAuthError(fmt.Sprintf("empty %q specified", AllowedHostsProp), nil)
	}
	for _, pattern := range allowedHosts {
		if pattern.MatchString(string(conn.Address())) {
			return nil
		}
	}
	return newAuthError(fmt.Sprintf("address %q not allowed by %q: %v", conn.Address(), AllowedHostsProp, allowedHosts), nil)
}

type oidcOneStep struct {
	userName    string
	accessToken string
}

type oidcTwoStep struct {
	conn driver.Connection
	oa   *OIDCAuthenticator
}

func jwtStepRequest(accessToken string) []byte {
	return bsoncore.NewDocumentBuilder().
		AppendString("jwt", accessToken).
		Build()
}

func principalStepRequest(principal string) []byte {
	doc := bsoncore.NewDocumentBuilder()
	if principal != "" {
		doc.AppendString("n", principal)
	}
	return doc.Build()
}

func (oos *oidcOneStep) Start() (string, []byte, error) {
	return MongoDBOIDC, jwtStepRequest(oos.accessToken), nil
}

func (oos *oidcOneStep) Next(context.Context, []byte) ([]byte, error) {
	return nil, newAuthError("unexpected step in OIDC authentication", nil)
}

func (*oidcOneStep) Completed() bool {
	return true
}

func (ots *oidcTwoStep) Start() (string, []byte, error) {
	return MongoDBOIDC, principalStepRequest(ots.oa.userName), nil
}

func (ots *oidcTwoStep) Next(ctx context.Context, msg []byte) ([]byte, error) {
	var idpInfo IDPInfo
	err := bson.Unmarshal(msg, &idpInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling BSON document: %w", err)
	}

	accessToken, err := ots.oa.getAccessToken(ctx,
		ots.conn,
		&OIDCArgs{
			Version: apiVersion,
			// idpInfo is nil for machine callbacks in the current spec.
			IDPInfo: &idpInfo,
			// there is no way there could be a refresh token when there is no IDPInfo.
			RefreshToken: nil,
		},
		// two-step callbacks are always human callbacks.
		ots.oa.OIDCHumanCallback)

	return jwtStepRequest(accessToken), err
}

func (*oidcTwoStep) Completed() bool {
	return true
}

func (oa *OIDCAuthenticator) providerCallback() (OIDCCallback, error) {
	env, ok := oa.AuthMechanismProperties[EnvironmentProp]
	if !ok {
		return nil, nil
	}

	switch env {
	case AzureEnvironmentValue:
		resource, ok := oa.AuthMechanismProperties[ResourceProp]
		if !ok {
			return nil, newAuthError(fmt.Sprintf("%q must be specified for Azure OIDC", ResourceProp), nil)
		}
		return getAzureOIDCCallback(oa.userName, resource, oa.httpClient), nil
	case GCPEnvironmentValue:
		resource, ok := oa.AuthMechanismProperties[ResourceProp]
		if !ok {
			return nil, newAuthError(fmt.Sprintf("%q must be specified for GCP OIDC", ResourceProp), nil)
		}
		return getGCPOIDCCallback(resource, oa.httpClient), nil
	}

	return nil, fmt.Errorf("%q %q not supported for MONGODB-OIDC", EnvironmentProp, env)
}

// getAzureOIDCCallback returns the callback for the Azure Identity Provider.
func getAzureOIDCCallback(clientID string, resource string, httpClient *http.Client) OIDCCallback {
	// return the callback parameterized by the clientID and resource, also passing in the user
	// configured httpClient.
	return func(ctx context.Context, _ *OIDCArgs) (*OIDCCredential, error) {
		resource = url.QueryEscape(resource)
		var uri string
		if clientID != "" {
			uri = fmt.Sprintf("http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=%s&client_id=%s", resource, clientID)
		} else {
			uri = fmt.Sprintf("http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=%s", resource)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
		if err != nil {
			return nil, newAuthError("error creating http request to Azure Identity Provider", err)
		}
		req.Header.Add("Metadata", "true")
		req.Header.Add("Accept", "application/json")
		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, newAuthError("error getting access token from Azure Identity Provider", err)
		}
		defer resp.Body.Close()
		var azureResp struct {
			AccessToken string `json:"access_token"`
			ExpiresOn   int64  `json:"expires_on,string"`
		}

		if resp.StatusCode != http.StatusOK {
			return nil, newAuthError(fmt.Sprintf("failed to get a valid response from Azure Identity Provider, http code: %d", resp.StatusCode), nil)
		}
		err = json.NewDecoder(resp.Body).Decode(&azureResp)
		if err != nil {
			return nil, newAuthError("failed parsing result from Azure Identity Provider", err)
		}
		expireTime := time.Unix(azureResp.ExpiresOn, 0)
		return &OIDCCredential{
			AccessToken: azureResp.AccessToken,
			ExpiresAt:   &expireTime,
		}, nil
	}
}

// getGCPOIDCCallback returns the callback for the GCP Identity Provider.
func getGCPOIDCCallback(resource string, httpClient *http.Client) OIDCCallback {
	// return the callback parameterized by the clientID and resource, also passing in the user
	// configured httpClient.
	return func(ctx context.Context, _ *OIDCArgs) (*OIDCCredential, error) {
		resource = url.QueryEscape(resource)
		uri := fmt.Sprintf("http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience=%s", resource)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
		if err != nil {
			return nil, newAuthError("error creating http request to GCP Identity Provider", err)
		}
		req.Header.Add("Metadata-Flavor", "Google")
		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, newAuthError("error getting access token from GCP Identity Provider", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, newAuthError(fmt.Sprintf("failed to get a valid response from GCP Identity Provider, http code: %d", resp.StatusCode), nil)
		}
		accessToken, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, newAuthError("failed parsing reading response from GCP Identity Provider", err)
		}
		return &OIDCCredential{
			AccessToken: string(accessToken),
			ExpiresAt:   nil,
		}, nil
	}
}

func (oa *OIDCAuthenticator) getAccessToken(
	ctx context.Context,
	conn driver.Connection,
	args *OIDCArgs,
	callback OIDCCallback,
) (string, error) {
	oa.mu.Lock()
	defer oa.mu.Unlock()

	if oa.accessToken != "" {
		return oa.accessToken, nil
	}

	// Attempt to refresh the access token if a refresh token is available.
	if args.RefreshToken != nil {
		cred, err := callback(ctx, args)
		if err == nil && cred != nil {
			oa.accessToken = cred.AccessToken
			oa.tokenGenID++
			conn.SetOIDCTokenGenID(oa.tokenGenID)
			oa.refreshToken = cred.RefreshToken
			return cred.AccessToken, nil
		}
		oa.refreshToken = nil
		args.RefreshToken = nil
	}
	// If we get here this means there either was no refresh token or the refresh token failed.
	cred, err := callback(ctx, args)
	if err != nil {
		return "", err
	}
	// This line should never occur, if go conventions are followed, but it is a safety check such
	// that we do not throw nil pointer errors to our users if they abuse the API.
	if cred == nil {
		return "", newAuthError("OIDC callback returned nil credential with no specified error", nil)
	}

	oa.accessToken = cred.AccessToken
	oa.tokenGenID++
	conn.SetOIDCTokenGenID(oa.tokenGenID)
	oa.refreshToken = cred.RefreshToken
	// always set the IdPInfo, in most cases, this should just be recopying the same pointer, or nil
	// in the machine flow.
	oa.idpInfo = args.IDPInfo

	return cred.AccessToken, nil
}

// invalidateAccessToken invalidates the access token, if the force flag is set to true (which is
// only on a Reauth call) or if the tokenGenID of the connection is greater than or equal to the
// tokenGenID of the OIDCAuthenticator. It should never actually be greater than, but only equal,
// but this is a safety check, since extra invalidation is only a performance impact, not a
// correctness impact.
func (oa *OIDCAuthenticator) invalidateAccessToken(conn driver.Connection) {
	oa.mu.Lock()
	defer oa.mu.Unlock()
	tokenGenID := conn.OIDCTokenGenID()
	// If the connection used in a Reauth is a new connection it will not have a correct tokenGenID,
	// it will instead be set to 0. In the absence of information, the only safe thing to do is to
	// invalidate the cached accessToken.
	if tokenGenID == 0 || tokenGenID >= oa.tokenGenID {
		oa.accessToken = ""
		conn.SetOIDCTokenGenID(0)
	}
}

// Reauth reauthenticates the connection when the server returns a 391 code. Reauth is part of the
// driver.Authenticator interface.
func (oa *OIDCAuthenticator) Reauth(ctx context.Context, cfg *Config) error {
	oa.invalidateAccessToken(cfg.Connection)
	return oa.Auth(ctx, cfg)
}

// Auth authenticates the connection.
func (oa *OIDCAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	var err error

	if cfg == nil {
		return newAuthError(fmt.Sprintf("config must be set for %q authentication", MongoDBOIDC), nil)
	}
	conn := cfg.Connection

	oa.mu.Lock()
	cachedAccessToken := oa.accessToken
	cachedRefreshToken := oa.refreshToken
	cachedIDPInfo := oa.idpInfo
	oa.mu.Unlock()

	if cachedAccessToken != "" {
		err = ConductSaslConversation(ctx, cfg, sourceExternal, &oidcOneStep{
			userName:    oa.userName,
			accessToken: cachedAccessToken,
		})
		if err == nil {
			return nil
		}
		// this seems like it could be incorrect since we could be inavlidating an access token that
		// has already been replaced by a different auth attempt, but the TokenGenID will prevernt
		// that from happening.
		oa.invalidateAccessToken(conn)
		time.Sleep(invalidateSleepTimeout)
	}

	if oa.OIDCHumanCallback != nil {
		return oa.doAuthHuman(ctx, cfg, oa.OIDCHumanCallback, cachedIDPInfo, cachedRefreshToken)
	}

	// Handle user provided or automatic provider machine callback.
	var machineCallback OIDCCallback
	if oa.OIDCMachineCallback != nil {
		machineCallback = oa.OIDCMachineCallback
	} else {
		machineCallback, err = oa.providerCallback()
		if err != nil {
			return fmt.Errorf("error getting built-in OIDC provider: %w", err)
		}
	}

	if machineCallback != nil {
		return oa.doAuthMachine(ctx, cfg, machineCallback)
	}
	return newAuthError("no OIDC callback provided", nil)
}

func (oa *OIDCAuthenticator) doAuthHuman(ctx context.Context, cfg *Config, humanCallback OIDCCallback, idpInfo *IDPInfo, refreshToken *string) error {
	// Ensure that the connection address is allowed by the allowed hosts.
	err := oa.validateConnectionAddressWithAllowedHosts(cfg.Connection)
	if err != nil {
		return err
	}
	subCtx, cancel := context.WithTimeout(ctx, humanCallbackTimeout)
	defer cancel()
	// If the idpInfo exists, we can just do one step
	if idpInfo != nil {
		accessToken, err := oa.getAccessToken(subCtx,
			cfg.Connection,
			&OIDCArgs{
				Version: apiVersion,
				// idpInfo is nil for machine callbacks in the current spec.
				IDPInfo:      idpInfo,
				RefreshToken: refreshToken,
			},
			humanCallback)
		if err != nil {
			return err
		}
		return ConductSaslConversation(
			subCtx,
			cfg,
			sourceExternal,
			&oidcOneStep{accessToken: accessToken},
		)
	}
	// otherwise, we need the two step where we ask the server for the IdPInfo first.
	ots := &oidcTwoStep{
		conn: cfg.Connection,
		oa:   oa,
	}
	return ConductSaslConversation(subCtx, cfg, sourceExternal, ots)
}

func (oa *OIDCAuthenticator) doAuthMachine(ctx context.Context, cfg *Config, machineCallback OIDCCallback) error {
	subCtx, cancel := context.WithTimeout(ctx, machineCallbackTimeout)
	accessToken, err := oa.getAccessToken(subCtx,
		cfg.Connection,
		&OIDCArgs{
			Version: apiVersion,
			// idpInfo is nil for machine callbacks in the current spec.
			IDPInfo:      nil,
			RefreshToken: nil,
		},
		machineCallback)
	cancel()
	if err != nil {
		return err
	}
	return ConductSaslConversation(
		ctx,
		cfg,
		sourceExternal,
		&oidcOneStep{accessToken: accessToken},
	)
}

// CreateSpeculativeConversation creates a speculative conversation for OIDC authentication.
func (oa *OIDCAuthenticator) CreateSpeculativeConversation() (SpeculativeConversation, error) {
	oa.mu.Lock()
	defer oa.mu.Unlock()
	accessToken := oa.accessToken
	if accessToken == "" {
		return nil, nil // Skip speculative auth.
	}

	return newSaslConversation(&oidcOneStep{accessToken: accessToken}, sourceExternal, true), nil
}
