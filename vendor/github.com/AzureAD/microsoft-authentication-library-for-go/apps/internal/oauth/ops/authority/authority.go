// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authority

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	authorizationEndpoint             = "https://%v/%v/oauth2/v2.0/authorize"
	aadInstanceDiscoveryEndpoint      = "https://%v/common/discovery/instance"
	tenantDiscoveryEndpointWithRegion = "https://%s.%s/%s/v2.0/.well-known/openid-configuration"
	regionName                        = "REGION_NAME"
	defaultAPIVersion                 = "2021-10-01"
	imdsEndpoint                      = "http://169.254.169.254/metadata/instance/compute/location?format=text&api-version=" + defaultAPIVersion
	autoDetectRegion                  = "TryAutoDetect"
	AccessTokenTypeBearer             = "Bearer"
)

// These are various hosts that host AAD Instance discovery endpoints.
const (
	defaultHost          = "login.microsoftonline.com"
	loginMicrosoft       = "login.microsoft.com"
	loginWindows         = "login.windows.net"
	loginSTSWindows      = "sts.windows.net"
	loginMicrosoftOnline = defaultHost
)

// jsonCaller is an interface that allows us to mock the JSONCall method.
type jsonCaller interface {
	JSONCall(ctx context.Context, endpoint string, headers http.Header, qv url.Values, body, resp interface{}) error
}

var aadTrustedHostList = map[string]bool{
	"login.windows.net":                true, // Microsoft Azure Worldwide - Used in validation scenarios where host is not this list
	"login.partner.microsoftonline.cn": true, // Microsoft Azure China
	"login.microsoftonline.de":         true, // Microsoft Azure Blackforest
	"login-us.microsoftonline.com":     true, // Microsoft Azure US Government - Legacy
	"login.microsoftonline.us":         true, // Microsoft Azure US Government
	"login.microsoftonline.com":        true, // Microsoft Azure Worldwide
}

// TrustedHost checks if an AAD host is trusted/valid.
func TrustedHost(host string) bool {
	if _, ok := aadTrustedHostList[host]; ok {
		return true
	}
	return false
}

// OAuthResponseBase is the base JSON return message for an OAuth call.
// This is embedded in other calls to get the base fields from every response.
type OAuthResponseBase struct {
	Error            string `json:"error"`
	SubError         string `json:"suberror"`
	ErrorDescription string `json:"error_description"`
	ErrorCodes       []int  `json:"error_codes"`
	CorrelationID    string `json:"correlation_id"`
	Claims           string `json:"claims"`
}

// TenantDiscoveryResponse is the tenant endpoints from the OpenID configuration endpoint.
type TenantDiscoveryResponse struct {
	OAuthResponseBase

	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	Issuer                string `json:"issuer"`

	AdditionalFields map[string]interface{}
}

// Validate validates that the response had the correct values required.
func (r *TenantDiscoveryResponse) Validate() error {
	switch "" {
	case r.AuthorizationEndpoint:
		return errors.New("TenantDiscoveryResponse: authorize endpoint was not found in the openid configuration")
	case r.TokenEndpoint:
		return errors.New("TenantDiscoveryResponse: token endpoint was not found in the openid configuration")
	case r.Issuer:
		return errors.New("TenantDiscoveryResponse: issuer was not found in the openid configuration")
	}
	return nil
}

type InstanceDiscoveryMetadata struct {
	PreferredNetwork string   `json:"preferred_network"`
	PreferredCache   string   `json:"preferred_cache"`
	Aliases          []string `json:"aliases"`

	AdditionalFields map[string]interface{}
}

type InstanceDiscoveryResponse struct {
	TenantDiscoveryEndpoint string                      `json:"tenant_discovery_endpoint"`
	Metadata                []InstanceDiscoveryMetadata `json:"metadata"`

	AdditionalFields map[string]interface{}
}

//go:generate stringer -type=AuthorizeType

// AuthorizeType represents the type of token flow.
type AuthorizeType int

// These are all the types of token flows.
const (
	ATUnknown AuthorizeType = iota
	ATUsernamePassword
	ATWindowsIntegrated
	ATAuthCode
	ATInteractive
	ATClientCredentials
	ATDeviceCode
	ATRefreshToken
	AccountByID
	ATOnBehalfOf
)

// These are all authority types
const (
	AAD  = "MSSTS"
	ADFS = "ADFS"
	DSTS = "DSTS"
)

// DSTSTenant is referenced throughout multiple files, let us use a const in case we ever need to change it.
const DSTSTenant = "7a433bfc-2514-4697-b467-e0933190487f"

// AuthenticationScheme is an extensibility mechanism designed to be used only by Azure Arc for proof of possession access tokens.
type AuthenticationScheme interface {
	// Extra parameters that are added to the request to the /token endpoint.
	TokenRequestParams() map[string]string
	// Key ID of the public / private key pair used by the encryption algorithm, if any.
	// Tokens obtained by authentication schemes that use this are bound to the KeyId, i.e.
	// if a different kid is presented, the access token cannot be used.
	KeyID() string
	// Creates the access token that goes into an Authorization HTTP header.
	FormatAccessToken(accessToken string) (string, error)
	//Expected to match the token_type parameter returned by ESTS. Used to disambiguate
	// between ATs of different types (e.g. Bearer and PoP) when loading from cache etc.
	AccessTokenType() string
}

// default authn scheme realizing AuthenticationScheme for "Bearer" tokens
type BearerAuthenticationScheme struct{}

var bearerAuthnScheme BearerAuthenticationScheme

func (ba *BearerAuthenticationScheme) TokenRequestParams() map[string]string {
	return nil
}
func (ba *BearerAuthenticationScheme) KeyID() string {
	return ""
}
func (ba *BearerAuthenticationScheme) FormatAccessToken(accessToken string) (string, error) {
	return accessToken, nil
}
func (ba *BearerAuthenticationScheme) AccessTokenType() string {
	return AccessTokenTypeBearer
}

// AuthParams represents the parameters used for authorization for token acquisition.
type AuthParams struct {
	AuthorityInfo Info
	CorrelationID string
	Endpoints     Endpoints
	ClientID      string
	// Redirecturi is used for auth flows that specify a redirect URI (e.g. local server for interactive auth flow).
	Redirecturi   string
	HomeAccountID string
	// Username is the user-name portion for username/password auth flow.
	Username string
	// Password is the password portion for username/password auth flow.
	Password string
	// Scopes is the list of scopes the user consents to.
	Scopes []string
	// AuthorizationType specifies the auth flow being used.
	AuthorizationType AuthorizeType
	// State is a random value used to prevent cross-site request forgery attacks.
	State string
	// CodeChallenge is derived from a code verifier and is sent in the auth request.
	CodeChallenge string
	// CodeChallengeMethod describes the method used to create the CodeChallenge.
	CodeChallengeMethod string
	// Prompt specifies the user prompt type during interactive auth.
	Prompt string
	// IsConfidentialClient specifies if it is a confidential client.
	IsConfidentialClient bool
	// SendX5C specifies if x5c claim(public key of the certificate) should be sent to STS.
	SendX5C bool
	// UserAssertion is the access token used to acquire token on behalf of user
	UserAssertion string
	// Capabilities the client will include with each token request, for example "CP1".
	// Call [NewClientCapabilities] to construct a value for this field.
	Capabilities ClientCapabilities
	// Claims required for an access token to satisfy a conditional access policy
	Claims string
	// KnownAuthorityHosts don't require metadata discovery because they're known to the user
	KnownAuthorityHosts []string
	// LoginHint is a username with which to pre-populate account selection during interactive auth
	LoginHint string
	// DomainHint is a directive that can be used to accelerate the user to their federated IdP sign-in page
	DomainHint string
	// AuthnScheme is an optional scheme for formatting access tokens
	AuthnScheme AuthenticationScheme
}

// NewAuthParams creates an authorization parameters object.
func NewAuthParams(clientID string, authorityInfo Info) AuthParams {
	return AuthParams{
		ClientID:      clientID,
		AuthorityInfo: authorityInfo,
		CorrelationID: uuid.New().String(),
		AuthnScheme:   &bearerAuthnScheme,
	}
}

// WithTenant returns a copy of the AuthParams having the specified tenant ID. If the given
// ID is empty, the copy is identical to the original. This function returns an error in
// several cases:
//   - ID isn't specific (for example, it's "common")
//   - ID is non-empty and the authority doesn't support tenants (for example, it's an ADFS authority)
//   - the client is configured to authenticate only Microsoft accounts via the "consumers" endpoint
//   - the resulting authority URL is invalid
func (p AuthParams) WithTenant(ID string) (AuthParams, error) {
	if ID == "" || ID == p.AuthorityInfo.Tenant {
		return p, nil
	}

	var authority string
	switch p.AuthorityInfo.AuthorityType {
	case AAD:
		if ID == "common" || ID == "consumers" || ID == "organizations" {
			return p, fmt.Errorf(`tenant ID must be a specific tenant, not "%s"`, ID)
		}
		if p.AuthorityInfo.Tenant == "consumers" {
			return p, errors.New(`client is configured to authenticate only personal Microsoft accounts, via the "consumers" endpoint`)
		}
		authority = "https://" + path.Join(p.AuthorityInfo.Host, ID)
	case ADFS:
		return p, errors.New("ADFS authority doesn't support tenants")
	case DSTS:
		return p, errors.New("dSTS authority doesn't support tenants")
	}

	info, err := NewInfoFromAuthorityURI(authority, p.AuthorityInfo.ValidateAuthority, p.AuthorityInfo.InstanceDiscoveryDisabled)
	if err == nil {
		info.Region = p.AuthorityInfo.Region
		p.AuthorityInfo = info
	}
	return p, err
}

// MergeCapabilitiesAndClaims combines client capabilities and challenge claims into a value suitable for an authentication request's "claims" parameter.
func (p AuthParams) MergeCapabilitiesAndClaims() (string, error) {
	claims := p.Claims
	if len(p.Capabilities.asMap) > 0 {
		if claims == "" {
			// without claims the result is simply the capabilities
			return p.Capabilities.asJSON, nil
		}
		// Otherwise, merge claims and capabilties into a single JSON object.
		// We handle the claims challenge as a map because we don't know its structure.
		var challenge map[string]any
		if err := json.Unmarshal([]byte(claims), &challenge); err != nil {
			return "", fmt.Errorf(`claims must be JSON. Are they base64 encoded? json.Unmarshal returned "%v"`, err)
		}
		if err := merge(p.Capabilities.asMap, challenge); err != nil {
			return "", err
		}
		b, err := json.Marshal(challenge)
		if err != nil {
			return "", err
		}
		claims = string(b)
	}
	return claims, nil
}

// merges a into b without overwriting b's values. Returns an error when a and b share a key for which either has a non-object value.
func merge(a, b map[string]any) error {
	for k, av := range a {
		if bv, ok := b[k]; !ok {
			// b doesn't contain this key => simply set it to a's value
			b[k] = av
		} else {
			// b does contain this key => recursively merge a[k] into b[k], provided both are maps. If a[k] or b[k] isn't
			// a map, return an error because merging would overwrite some value in b. Errors shouldn't occur in practice
			// because the challenge will be from AAD, which knows the capabilities format.
			if A, ok := av.(map[string]any); ok {
				if B, ok := bv.(map[string]any); ok {
					return merge(A, B)
				} else {
					// b[k] isn't a map
					return errors.New("challenge claims conflict with client capabilities")
				}
			} else {
				// a[k] isn't a map
				return errors.New("challenge claims conflict with client capabilities")
			}
		}
	}
	return nil
}

// ClientCapabilities stores capabilities in the formats used by AuthParams.MergeCapabilitiesAndClaims.
// [NewClientCapabilities] precomputes these representations because capabilities are static for the
// lifetime of a client and are included with every authentication request i.e., these computations
// always have the same result and would otherwise have to be repeated for every request.
type ClientCapabilities struct {
	// asJSON is for the common case: adding the capabilities to an auth request with no challenge claims
	asJSON string
	// asMap is for merging the capabilities with challenge claims
	asMap map[string]any
}

func NewClientCapabilities(capabilities []string) (ClientCapabilities, error) {
	c := ClientCapabilities{}
	var err error
	if len(capabilities) > 0 {
		cpbs := make([]string, len(capabilities))
		for i := 0; i < len(cpbs); i++ {
			cpbs[i] = fmt.Sprintf(`"%s"`, capabilities[i])
		}
		c.asJSON = fmt.Sprintf(`{"access_token":{"xms_cc":{"values":[%s]}}}`, strings.Join(cpbs, ","))
		// note our JSON is valid but we can't stop users breaking it with garbage like "}"
		err = json.Unmarshal([]byte(c.asJSON), &c.asMap)
	}
	return c, err
}

// Info consists of information about the authority.
type Info struct {
	Host                      string
	CanonicalAuthorityURI     string
	AuthorityType             string
	ValidateAuthority         bool
	Tenant                    string
	Region                    string
	InstanceDiscoveryDisabled bool
}

// NewInfoFromAuthorityURI creates an AuthorityInfo instance from the authority URL provided.
func NewInfoFromAuthorityURI(authority string, validateAuthority bool, instanceDiscoveryDisabled bool) (Info, error) {

	cannonicalAuthority := authority

	// suffix authority with / if it doesn't have one
	if !strings.HasSuffix(cannonicalAuthority, "/") {
		cannonicalAuthority += "/"
	}

	u, err := url.Parse(strings.ToLower(cannonicalAuthority))

	if err != nil {
		return Info{}, fmt.Errorf("couldn't parse authority url: %w", err)
	}
	if u.Scheme != "https" {
		return Info{}, errors.New("authority url scheme must be https")
	}

	pathParts := strings.Split(u.EscapedPath(), "/")
	if len(pathParts) < 3 {
		return Info{}, errors.New(`authority must be an URL such as "https://login.microsoftonline.com/<your tenant>"`)
	}

	var authorityType, tenant string
	switch pathParts[1] {
	case "adfs":
		authorityType = ADFS
	case "dstsv2":
		if len(pathParts) != 4 {
			return Info{}, fmt.Errorf("dSTS authority must be an https URL such as https://<authority>/dstsv2/%s", DSTSTenant)
		}
		if pathParts[2] != DSTSTenant {
			return Info{}, fmt.Errorf("dSTS authority only accepts a single tenant %q", DSTSTenant)
		}
		authorityType = DSTS
		tenant = DSTSTenant
	default:
		authorityType = AAD
		tenant = pathParts[1]
	}

	// u.Host includes the port, if any, which is required for private cloud deployments
	return Info{
		Host:                      u.Host,
		CanonicalAuthorityURI:     cannonicalAuthority,
		AuthorityType:             authorityType,
		ValidateAuthority:         validateAuthority,
		Tenant:                    tenant,
		InstanceDiscoveryDisabled: instanceDiscoveryDisabled,
	}, nil
}

// Endpoints consists of the endpoints from the tenant discovery response.
type Endpoints struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
	selfSignedJwtAudience string
	authorityHost         string
}

// NewEndpoints creates an Endpoints object.
func NewEndpoints(authorizationEndpoint string, tokenEndpoint string, selfSignedJwtAudience string, authorityHost string) Endpoints {
	return Endpoints{authorizationEndpoint, tokenEndpoint, selfSignedJwtAudience, authorityHost}
}

// UserRealmAccountType refers to the type of user realm.
type UserRealmAccountType string

// These are the different types of user realms.
const (
	Unknown   UserRealmAccountType = ""
	Federated UserRealmAccountType = "Federated"
	Managed   UserRealmAccountType = "Managed"
)

// UserRealm is used for the username password request to determine user type
type UserRealm struct {
	AccountType       UserRealmAccountType `json:"account_type"`
	DomainName        string               `json:"domain_name"`
	CloudInstanceName string               `json:"cloud_instance_name"`
	CloudAudienceURN  string               `json:"cloud_audience_urn"`

	// required if accountType is Federated
	FederationProtocol    string `json:"federation_protocol"`
	FederationMetadataURL string `json:"federation_metadata_url"`

	AdditionalFields map[string]interface{}
}

func (u UserRealm) validate() error {
	switch "" {
	case string(u.AccountType):
		return errors.New("the account type (Federated or Managed) is missing")
	case u.DomainName:
		return errors.New("domain name of user realm is missing")
	case u.CloudInstanceName:
		return errors.New("cloud instance name of user realm is missing")
	case u.CloudAudienceURN:
		return errors.New("cloud Instance URN is missing")
	}

	if u.AccountType == Federated {
		switch "" {
		case u.FederationProtocol:
			return errors.New("federation protocol of user realm is missing")
		case u.FederationMetadataURL:
			return errors.New("federation metadata URL of user realm is missing")
		}
	}
	return nil
}

// Client represents the REST calls to authority backends.
type Client struct {
	// Comm provides the HTTP transport client.
	Comm jsonCaller // *comm.Client
}

func (c Client) UserRealm(ctx context.Context, authParams AuthParams) (UserRealm, error) {
	endpoint := fmt.Sprintf("https://%s/common/UserRealm/%s", authParams.Endpoints.authorityHost, url.PathEscape(authParams.Username))
	qv := url.Values{
		"api-version": []string{"1.0"},
	}

	resp := UserRealm{}
	err := c.Comm.JSONCall(
		ctx,
		endpoint,
		http.Header{"client-request-id": []string{authParams.CorrelationID}},
		qv,
		nil,
		&resp,
	)
	if err != nil {
		return resp, err
	}

	return resp, resp.validate()
}

func (c Client) GetTenantDiscoveryResponse(ctx context.Context, openIDConfigurationEndpoint string) (TenantDiscoveryResponse, error) {
	resp := TenantDiscoveryResponse{}
	err := c.Comm.JSONCall(
		ctx,
		openIDConfigurationEndpoint,
		http.Header{},
		nil,
		nil,
		&resp,
	)

	return resp, err
}

// AADInstanceDiscovery attempts to discover a tenant endpoint (used in OIDC auth with an authorization endpoint).
// This is done by AAD which allows for aliasing of tenants (windows.sts.net is the same as login.windows.com).
func (c Client) AADInstanceDiscovery(ctx context.Context, authorityInfo Info) (InstanceDiscoveryResponse, error) {
	region := ""
	var err error
	resp := InstanceDiscoveryResponse{}
	if authorityInfo.Region != "" && authorityInfo.Region != autoDetectRegion {
		region = authorityInfo.Region
	} else if authorityInfo.Region == autoDetectRegion {
		region = detectRegion(ctx)
	}
	if region != "" {
		environment := authorityInfo.Host
		switch environment {
		case loginMicrosoft, loginWindows, loginSTSWindows, defaultHost:
			environment = loginMicrosoft
		}

		resp.TenantDiscoveryEndpoint = fmt.Sprintf(tenantDiscoveryEndpointWithRegion, region, environment, authorityInfo.Tenant)
		metadata := InstanceDiscoveryMetadata{
			PreferredNetwork: fmt.Sprintf("%v.%v", region, authorityInfo.Host),
			PreferredCache:   authorityInfo.Host,
			Aliases:          []string{fmt.Sprintf("%v.%v", region, authorityInfo.Host), authorityInfo.Host},
		}
		resp.Metadata = []InstanceDiscoveryMetadata{metadata}
	} else {
		qv := url.Values{}
		qv.Set("api-version", "1.1")
		qv.Set("authorization_endpoint", fmt.Sprintf(authorizationEndpoint, authorityInfo.Host, authorityInfo.Tenant))

		discoveryHost := defaultHost
		if TrustedHost(authorityInfo.Host) {
			discoveryHost = authorityInfo.Host
		}

		endpoint := fmt.Sprintf(aadInstanceDiscoveryEndpoint, discoveryHost)
		err = c.Comm.JSONCall(ctx, endpoint, http.Header{}, qv, nil, &resp)
	}
	return resp, err
}

func detectRegion(ctx context.Context) string {
	region := os.Getenv(regionName)
	if region != "" {
		region = strings.ReplaceAll(region, " ", "")
		return strings.ToLower(region)
	}
	// HTTP call to IMDS endpoint to get region
	// Refer : https://identitydivision.visualstudio.com/DevEx/_git/AuthLibrariesApiReview?path=%2FPinAuthToRegion%2FAAD%20SDK%20Proposal%20to%20Pin%20Auth%20to%20region.md&_a=preview&version=GBdev
	// Set a 2 second timeout for this http client which only does calls to IMDS endpoint
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, imdsEndpoint, nil)
	req.Header.Set("Metadata", "true")
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
	// If the request times out or there is an error, it is retried once
	if err != nil || resp.StatusCode != http.StatusOK {
		resp, err = client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return ""
		}
	}
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(response)
}

func (a *AuthParams) CacheKey(isAppCache bool) string {
	if a.AuthorizationType == ATOnBehalfOf {
		return a.AssertionHash()
	}
	if a.AuthorizationType == ATClientCredentials || isAppCache {
		return a.AppKey()
	}
	if a.AuthorizationType == ATRefreshToken || a.AuthorizationType == AccountByID {
		return a.HomeAccountID
	}
	return ""
}
func (a *AuthParams) AssertionHash() string {
	hasher := sha256.New()
	// Per documentation this never returns an error : https://pkg.go.dev/hash#pkg-types
	_, _ = hasher.Write([]byte(a.UserAssertion))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func (a *AuthParams) AppKey() string {
	if a.AuthorityInfo.Tenant != "" {
		return fmt.Sprintf("%s_%s_AppTokenCache", a.ClientID, a.AuthorityInfo.Tenant)
	}
	return fmt.Sprintf("%s__AppTokenCache", a.ClientID)
}
