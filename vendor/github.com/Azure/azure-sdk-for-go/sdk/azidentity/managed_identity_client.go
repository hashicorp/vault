//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

const (
	arcIMDSEndpoint          = "IMDS_ENDPOINT"
	defaultIdentityClientID  = "DEFAULT_IDENTITY_CLIENT_ID"
	identityEndpoint         = "IDENTITY_ENDPOINT"
	identityHeader           = "IDENTITY_HEADER"
	identityServerThumbprint = "IDENTITY_SERVER_THUMBPRINT"
	headerMetadata           = "Metadata"
	imdsEndpoint             = "http://169.254.169.254/metadata/identity/oauth2/token"
	miResID                  = "mi_res_id"
	msiEndpoint              = "MSI_ENDPOINT"
	msiResID                 = "msi_res_id"
	msiSecret                = "MSI_SECRET"
	imdsAPIVersion           = "2018-02-01"
	azureArcAPIVersion       = "2019-08-15"
	qpClientID               = "client_id"
	serviceFabricAPIVersion  = "2019-07-01-preview"
)

var imdsProbeTimeout = time.Second

type msiType int

const (
	msiTypeAppService msiType = iota
	msiTypeAzureArc
	msiTypeAzureML
	msiTypeCloudShell
	msiTypeIMDS
	msiTypeServiceFabric
)

type managedIdentityClient struct {
	azClient  *azcore.Client
	endpoint  string
	id        ManagedIDKind
	msiType   msiType
	probeIMDS bool
}

// arcKeyDirectory returns the directory expected to contain Azure Arc keys
var arcKeyDirectory = func() (string, error) {
	switch runtime.GOOS {
	case "linux":
		return "/var/opt/azcmagent/tokens", nil
	case "windows":
		pd := os.Getenv("ProgramData")
		if pd == "" {
			return "", errors.New("environment variable ProgramData has no value")
		}
		return filepath.Join(pd, "AzureConnectedMachineAgent", "Tokens"), nil
	default:
		return "", fmt.Errorf("unsupported OS %q", runtime.GOOS)
	}
}

type wrappedNumber json.Number

func (n *wrappedNumber) UnmarshalJSON(b []byte) error {
	c := string(b)
	if c == "\"\"" {
		return nil
	}
	return json.Unmarshal(b, (*json.Number)(n))
}

// setIMDSRetryOptionDefaults sets zero-valued fields to default values appropriate for IMDS
func setIMDSRetryOptionDefaults(o *policy.RetryOptions) {
	if o.MaxRetries == 0 {
		o.MaxRetries = 5
	}
	if o.MaxRetryDelay == 0 {
		o.MaxRetryDelay = 1 * time.Minute
	}
	if o.RetryDelay == 0 {
		o.RetryDelay = 2 * time.Second
	}
	if o.StatusCodes == nil {
		o.StatusCodes = []int{
			// IMDS docs recommend retrying 404, 410, 429 and 5xx
			// https://learn.microsoft.com/entra/identity/managed-identities-azure-resources/how-to-use-vm-token#error-handling
			http.StatusNotFound,                      // 404
			http.StatusGone,                          // 410
			http.StatusTooManyRequests,               // 429
			http.StatusInternalServerError,           // 500
			http.StatusNotImplemented,                // 501
			http.StatusBadGateway,                    // 502
			http.StatusServiceUnavailable,            // 503
			http.StatusGatewayTimeout,                // 504
			http.StatusHTTPVersionNotSupported,       // 505
			http.StatusVariantAlsoNegotiates,         // 506
			http.StatusInsufficientStorage,           // 507
			http.StatusLoopDetected,                  // 508
			http.StatusNotExtended,                   // 510
			http.StatusNetworkAuthenticationRequired, // 511
		}
	}
	if o.TryTimeout == 0 {
		o.TryTimeout = 1 * time.Minute
	}
}

// newManagedIdentityClient creates a new instance of the ManagedIdentityClient with the ManagedIdentityCredentialOptions
// that are passed into it along with a default pipeline.
// options: ManagedIdentityCredentialOptions configure policies for the pipeline and the authority host that
// will be used to retrieve tokens and authenticate
func newManagedIdentityClient(options *ManagedIdentityCredentialOptions) (*managedIdentityClient, error) {
	if options == nil {
		options = &ManagedIdentityCredentialOptions{}
	}
	cp := options.ClientOptions
	c := managedIdentityClient{id: options.ID, endpoint: imdsEndpoint, msiType: msiTypeIMDS}
	env := "IMDS"
	if endpoint, ok := os.LookupEnv(identityEndpoint); ok {
		if _, ok := os.LookupEnv(identityHeader); ok {
			if _, ok := os.LookupEnv(identityServerThumbprint); ok {
				if options.ID != nil {
					return nil, errors.New("the Service Fabric API doesn't support specifying a user-assigned managed identity at runtime")
				}
				env = "Service Fabric"
				c.endpoint = endpoint
				c.msiType = msiTypeServiceFabric
			} else {
				env = "App Service"
				c.endpoint = endpoint
				c.msiType = msiTypeAppService
			}
		} else if _, ok := os.LookupEnv(arcIMDSEndpoint); ok {
			if options.ID != nil {
				return nil, errors.New("the Azure Arc API doesn't support specifying a user-assigned managed identity at runtime")
			}
			env = "Azure Arc"
			c.endpoint = endpoint
			c.msiType = msiTypeAzureArc
		}
	} else if endpoint, ok := os.LookupEnv(msiEndpoint); ok {
		c.endpoint = endpoint
		if _, ok := os.LookupEnv(msiSecret); ok {
			if options.ID != nil && options.ID.idKind() != miClientID {
				return nil, errors.New("the Azure ML API supports specifying a user-assigned managed identity by client ID only")
			}
			env = "Azure ML"
			c.msiType = msiTypeAzureML
		} else {
			if options.ID != nil {
				return nil, errors.New("the Cloud Shell API doesn't support user-assigned managed identities")
			}
			env = "Cloud Shell"
			c.msiType = msiTypeCloudShell
		}
	} else {
		c.probeIMDS = options.dac
		setIMDSRetryOptionDefaults(&cp.Retry)
	}

	client, err := azcore.NewClient(module, version, azruntime.PipelineOptions{
		Tracing: azruntime.TracingOptions{
			Namespace: traceNamespace,
		},
	}, &cp)
	if err != nil {
		return nil, err
	}
	c.azClient = client

	if log.Should(EventAuthentication) {
		log.Writef(EventAuthentication, "Managed Identity Credential will use %s managed identity", env)
	}

	return &c, nil
}

// provideToken acquires a token for MSAL's confidential.Client, which caches the token
func (c *managedIdentityClient) provideToken(ctx context.Context, params confidential.TokenProviderParameters) (confidential.TokenProviderResult, error) {
	result := confidential.TokenProviderResult{}
	tk, err := c.authenticate(ctx, c.id, params.Scopes)
	if err == nil {
		result.AccessToken = tk.Token
		result.ExpiresInSeconds = int(time.Until(tk.ExpiresOn).Seconds())
	}
	return result, err
}

// authenticate acquires an access token
func (c *managedIdentityClient) authenticate(ctx context.Context, id ManagedIDKind, scopes []string) (azcore.AccessToken, error) {
	// no need to synchronize around this value because it's true only when DefaultAzureCredential constructed the client,
	// and in that case ChainedTokenCredential.GetToken synchronizes goroutines that would execute this block
	if c.probeIMDS {
		cx, cancel := context.WithTimeout(ctx, imdsProbeTimeout)
		defer cancel()
		cx = policy.WithRetryOptions(cx, policy.RetryOptions{MaxRetries: -1})
		req, err := azruntime.NewRequest(cx, http.MethodGet, c.endpoint)
		if err != nil {
			return azcore.AccessToken{}, fmt.Errorf("failed to create IMDS probe request: %s", err)
		}
		res, err := c.azClient.Pipeline().Do(req)
		if err != nil {
			msg := err.Error()
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				msg = "managed identity timed out. See https://aka.ms/azsdk/go/identity/troubleshoot#dac for more information"
			}
			return azcore.AccessToken{}, newCredentialUnavailableError(credNameManagedIdentity, msg)
		}
		// because IMDS always responds with JSON, assume a non-JSON response is from something else, such
		// as a proxy, and return credentialUnavailableError so DefaultAzureCredential continues iterating
		b, err := azruntime.Payload(res)
		if err != nil {
			return azcore.AccessToken{}, newCredentialUnavailableError(credNameManagedIdentity, fmt.Sprintf("failed to read IMDS probe response: %s", err))
		}
		if !json.Valid(b) {
			return azcore.AccessToken{}, newCredentialUnavailableError(credNameManagedIdentity, "unexpected response to IMDS probe")
		}
		// send normal token requests from now on because IMDS responded
		c.probeIMDS = false
	}

	msg, err := c.createAuthRequest(ctx, id, scopes)
	if err != nil {
		return azcore.AccessToken{}, err
	}

	resp, err := c.azClient.Pipeline().Do(msg)
	if err != nil {
		return azcore.AccessToken{}, newAuthenticationFailedError(credNameManagedIdentity, err.Error(), nil)
	}

	if azruntime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return c.createAccessToken(resp)
	}

	if c.msiType == msiTypeIMDS {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			if id != nil {
				return azcore.AccessToken{}, newAuthenticationFailedError(credNameManagedIdentity, "the requested identity isn't assigned to this resource", resp)
			}
			msg := "failed to authenticate a system assigned identity"
			if body, err := azruntime.Payload(resp); err == nil && len(body) > 0 {
				msg += fmt.Sprintf(". The endpoint responded with %s", body)
			}
			return azcore.AccessToken{}, newCredentialUnavailableError(credNameManagedIdentity, msg)
		case http.StatusForbidden:
			// Docker Desktop runs a proxy that responds 403 to IMDS token requests. If we get that response,
			// we return credentialUnavailableError so credential chains continue to their next credential
			body, err := azruntime.Payload(resp)
			if err == nil && strings.Contains(string(body), "unreachable") {
				return azcore.AccessToken{}, newCredentialUnavailableError(credNameManagedIdentity, fmt.Sprintf("unexpected response %q", string(body)))
			}
		}
	}

	return azcore.AccessToken{}, newAuthenticationFailedError(credNameManagedIdentity, "", resp)
}

func (c *managedIdentityClient) createAccessToken(res *http.Response) (azcore.AccessToken, error) {
	value := struct {
		// these are the only fields that we use
		Token        string        `json:"access_token,omitempty"`
		RefreshToken string        `json:"refresh_token,omitempty"`
		ExpiresIn    wrappedNumber `json:"expires_in,omitempty"` // this field should always return the number of seconds for which a token is valid
		ExpiresOn    interface{}   `json:"expires_on,omitempty"` // the value returned in this field varies between a number and a date string
	}{}
	if err := azruntime.UnmarshalAsJSON(res, &value); err != nil {
		return azcore.AccessToken{}, fmt.Errorf("internal AccessToken: %v", err)
	}
	if value.ExpiresIn != "" {
		expiresIn, err := json.Number(value.ExpiresIn).Int64()
		if err != nil {
			return azcore.AccessToken{}, err
		}
		return azcore.AccessToken{Token: value.Token, ExpiresOn: time.Now().Add(time.Second * time.Duration(expiresIn)).UTC()}, nil
	}
	switch v := value.ExpiresOn.(type) {
	case float64:
		return azcore.AccessToken{Token: value.Token, ExpiresOn: time.Unix(int64(v), 0).UTC()}, nil
	case string:
		if expiresOn, err := strconv.Atoi(v); err == nil {
			return azcore.AccessToken{Token: value.Token, ExpiresOn: time.Unix(int64(expiresOn), 0).UTC()}, nil
		}
		return azcore.AccessToken{}, newAuthenticationFailedError(credNameManagedIdentity, "unexpected expires_on value: "+v, res)
	default:
		msg := fmt.Sprintf("unsupported type received in expires_on: %T, %v", v, v)
		return azcore.AccessToken{}, newAuthenticationFailedError(credNameManagedIdentity, msg, res)
	}
}

func (c *managedIdentityClient) createAuthRequest(ctx context.Context, id ManagedIDKind, scopes []string) (*policy.Request, error) {
	switch c.msiType {
	case msiTypeIMDS:
		return c.createIMDSAuthRequest(ctx, id, scopes)
	case msiTypeAppService:
		return c.createAppServiceAuthRequest(ctx, id, scopes)
	case msiTypeAzureArc:
		// need to perform preliminary request to retreive the secret key challenge provided by the HIMDS service
		key, err := c.getAzureArcSecretKey(ctx, scopes)
		if err != nil {
			msg := fmt.Sprintf("failed to retreive secret key from the identity endpoint: %v", err)
			return nil, newAuthenticationFailedError(credNameManagedIdentity, msg, nil)
		}
		return c.createAzureArcAuthRequest(ctx, scopes, key)
	case msiTypeAzureML:
		return c.createAzureMLAuthRequest(ctx, id, scopes)
	case msiTypeServiceFabric:
		return c.createServiceFabricAuthRequest(ctx, scopes)
	case msiTypeCloudShell:
		return c.createCloudShellAuthRequest(ctx, scopes)
	default:
		return nil, newCredentialUnavailableError(credNameManagedIdentity, "managed identity isn't supported in this environment")
	}
}

func (c *managedIdentityClient) createIMDSAuthRequest(ctx context.Context, id ManagedIDKind, scopes []string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Raw().Header.Set(headerMetadata, "true")
	q := request.Raw().URL.Query()
	q.Set("api-version", imdsAPIVersion)
	q.Set("resource", strings.Join(scopes, " "))
	if id != nil {
		switch id.idKind() {
		case miClientID:
			q.Set(qpClientID, id.String())
		case miObjectID:
			q.Set("object_id", id.String())
		case miResourceID:
			q.Set(msiResID, id.String())
		}
	}
	request.Raw().URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createAppServiceAuthRequest(ctx context.Context, id ManagedIDKind, scopes []string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Raw().Header.Set("X-IDENTITY-HEADER", os.Getenv(identityHeader))
	q := request.Raw().URL.Query()
	q.Set("api-version", "2019-08-01")
	q.Set("resource", scopes[0])
	if id != nil {
		switch id.idKind() {
		case miClientID:
			q.Set(qpClientID, id.String())
		case miObjectID:
			q.Set("principal_id", id.String())
		case miResourceID:
			q.Set(miResID, id.String())
		}
	}
	request.Raw().URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createAzureMLAuthRequest(ctx context.Context, id ManagedIDKind, scopes []string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Raw().Header.Set("secret", os.Getenv(msiSecret))
	q := request.Raw().URL.Query()
	q.Set("api-version", "2017-09-01")
	q.Set("resource", strings.Join(scopes, " "))
	q.Set("clientid", os.Getenv(defaultIdentityClientID))
	if id != nil {
		switch id.idKind() {
		case miClientID:
			q.Set("clientid", id.String())
		case miObjectID:
			return nil, newAuthenticationFailedError(credNameManagedIdentity, "Azure ML doesn't support specifying a managed identity by object ID", nil)
		case miResourceID:
			return nil, newAuthenticationFailedError(credNameManagedIdentity, "Azure ML doesn't support specifying a managed identity by resource ID", nil)
		}
	}
	request.Raw().URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createServiceFabricAuthRequest(ctx context.Context, scopes []string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	q := request.Raw().URL.Query()
	request.Raw().Header.Set("Accept", "application/json")
	request.Raw().Header.Set("Secret", os.Getenv(identityHeader))
	q.Set("api-version", serviceFabricAPIVersion)
	q.Set("resource", strings.Join(scopes, " "))
	request.Raw().URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) getAzureArcSecretKey(ctx context.Context, resources []string) (string, error) {
	// create the request to retreive the secret key challenge provided by the HIMDS service
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return "", err
	}
	request.Raw().Header.Set(headerMetadata, "true")
	q := request.Raw().URL.Query()
	q.Set("api-version", azureArcAPIVersion)
	q.Set("resource", strings.Join(resources, " "))
	request.Raw().URL.RawQuery = q.Encode()
	// send the initial request to get the short-lived secret key
	response, err := c.azClient.Pipeline().Do(request)
	if err != nil {
		return "", err
	}
	// the endpoint is expected to return a 401 with the WWW-Authenticate header set to the location
	// of the secret key file. Any other status code indicates an error in the request.
	if response.StatusCode != 401 {
		msg := fmt.Sprintf("expected a 401 response, received %d", response.StatusCode)
		return "", newAuthenticationFailedError(credNameManagedIdentity, msg, response)
	}
	header := response.Header.Get("WWW-Authenticate")
	if len(header) == 0 {
		return "", newAuthenticationFailedError(credNameManagedIdentity, "HIMDS response has no WWW-Authenticate header", nil)
	}
	// the WWW-Authenticate header is expected in the following format: Basic realm=/some/file/path.key
	_, p, found := strings.Cut(header, "=")
	if !found {
		return "", newAuthenticationFailedError(credNameManagedIdentity, "unexpected WWW-Authenticate header from HIMDS: "+header, nil)
	}
	expected, err := arcKeyDirectory()
	if err != nil {
		return "", err
	}
	if filepath.Dir(p) != expected || !strings.HasSuffix(p, ".key") {
		return "", newAuthenticationFailedError(credNameManagedIdentity, "unexpected file path from HIMDS service: "+p, nil)
	}
	f, err := os.Stat(p)
	if err != nil {
		return "", newAuthenticationFailedError(credNameManagedIdentity, fmt.Sprintf("could not stat %q: %v", p, err), nil)
	}
	if s := f.Size(); s > 4096 {
		return "", newAuthenticationFailedError(credNameManagedIdentity, fmt.Sprintf("key is too large (%d bytes)", s), nil)
	}
	key, err := os.ReadFile(p)
	if err != nil {
		return "", newAuthenticationFailedError(credNameManagedIdentity, fmt.Sprintf("could not read %q: %v", p, err), nil)
	}
	return string(key), nil
}

func (c *managedIdentityClient) createAzureArcAuthRequest(ctx context.Context, resources []string, key string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Raw().Header.Set(headerMetadata, "true")
	request.Raw().Header.Set("Authorization", fmt.Sprintf("Basic %s", key))
	q := request.Raw().URL.Query()
	q.Set("api-version", azureArcAPIVersion)
	q.Set("resource", strings.Join(resources, " "))
	request.Raw().URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createCloudShellAuthRequest(ctx context.Context, scopes []string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodPost, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Raw().Header.Set(headerMetadata, "true")
	data := url.Values{}
	data.Set("resource", strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := streaming.NopCloser(strings.NewReader(dataEncoded))
	if err := request.SetBody(body, "application/x-www-form-urlencoded"); err != nil {
		return nil, err
	}
	return request, nil
}
