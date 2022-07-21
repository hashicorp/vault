package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/vault/api"
)

type AzureAuth struct {
	roleName  string
	mountPath string
	resource  string
}

var _ api.AuthMethod = (*AzureAuth)(nil)

type LoginOption func(a *AzureAuth) error

type responseJSON struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
}

type errorJSON struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type metadataJSON struct {
	Compute computeJSON `json:"compute"`
}

type computeJSON struct {
	VMName            string `json:"name"`
	VMScaleSetName    string `json:"vmScaleSetName"`
	SubscriptionID    string `json:"subscriptionId"`
	ResourceGroupName string `json:"resourceGroupName"`
}

const (
	defaultMountPath     = "azure"
	defaultResourceURL   = "https://management.azure.com/"
	metadataEndpoint     = "http://169.254.169.254"
	metadataAPIVersion   = "2021-05-01"
	apiVersionQueryParam = "api-version"
	resourceQueryParam   = "resource"
	clientTimeout        = 10 * time.Second
)

// NewAzureAuth initializes a new Azure auth method interface to be
// passed as a parameter to the client.Auth().Login method.
//
// Supported options: WithMountPath, WithResource
func NewAzureAuth(roleName string, opts ...LoginOption) (*AzureAuth, error) {
	if roleName == "" {
		return nil, fmt.Errorf("no role name provided for login")
	}

	a := &AzureAuth{
		roleName:  roleName,
		mountPath: defaultMountPath,
		resource:  defaultResourceURL,
	}

	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		// *AzureAuth as the argument
		err := opt(a)
		if err != nil {
			return nil, fmt.Errorf("error with login option: %w", err)
		}
	}

	// return the modified auth struct instance
	return a, nil
}

// Login sets up the required request body for the Azure auth method's /login
// endpoint, and performs a write to it.
func (a *AzureAuth) Login(ctx context.Context, client *api.Client) (*api.Secret, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	jwtResp, err := a.getJWT()
	if err != nil {
		return nil, fmt.Errorf("unable to get access token: %w", err)
	}

	metadataRespJSON, err := getMetadata()
	if err != nil {
		return nil, fmt.Errorf("unable to get instance metadata: %w", err)
	}

	loginData := map[string]interface{}{
		"role":                a.roleName,
		"jwt":                 jwtResp,
		"vm_name":             metadataRespJSON.Compute.VMName,
		"vmss_name":           metadataRespJSON.Compute.VMScaleSetName,
		"subscription_id":     metadataRespJSON.Compute.SubscriptionID,
		"resource_group_name": metadataRespJSON.Compute.ResourceGroupName,
	}

	path := fmt.Sprintf("auth/%s/login", a.mountPath)
	resp, err := client.Logical().WriteWithContext(ctx, path, loginData)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Azure auth: %w", err)
	}

	return resp, nil
}

func WithMountPath(mountPath string) LoginOption {
	return func(a *AzureAuth) error {
		a.mountPath = mountPath
		return nil
	}
}

// WithResource allows you to specify a different resource URL to use as the aud value
// on the JWT token than the default of Azure Public Cloud's ARM URL.
// This should match the resource URI that an administrator configured your
// Vault server to use.
//
// See https://github.com/Azure/go-autorest/blob/master/autorest/azure/environments.go
// for a list of valid environments.
func WithResource(url string) LoginOption {
	return func(a *AzureAuth) error {
		a.resource = url
		return nil
	}
}

// Retrieves an access token from Managed Identities for Azure Resources
//
// Learn more here: https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/how-to-use-vm-token
func (a *AzureAuth) getJWT() (string, error) {
	identityEndpoint, err := url.Parse(fmt.Sprintf("%s/metadata/identity/oauth2/token", metadataEndpoint))
	if err != nil {
		return "", fmt.Errorf("error creating metadata URL: %w", err)
	}

	identityParameters := identityEndpoint.Query()
	identityParameters.Add(apiVersionQueryParam, metadataAPIVersion)
	identityParameters.Add(resourceQueryParam, a.resource)
	identityEndpoint.RawQuery = identityParameters.Encode()

	req, err := http.NewRequest(http.MethodGet, identityEndpoint.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Add("Metadata", "true")

	client := &http.Client{
		Timeout: clientTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error calling Azure token endpoint: %w", err)
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body from Azure token endpoint: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp errorJSON
		err = json.Unmarshal(responseBytes, &errResp)
		if err != nil {
			return "", fmt.Errorf("received error message but was unable to unmarshal its contents")
		}
		return "", fmt.Errorf("%s error from Azure token endpoint: %s", errResp.Error, errResp.ErrorDescription)
	}

	var r responseJSON
	err = json.Unmarshal(responseBytes, &r)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response from Azure token endpoint: %w", err)
	}

	return r.AccessToken, nil
}

func getMetadata() (metadataJSON, error) {
	metadataEndpoint, err := url.Parse(fmt.Sprintf("%s/metadata/instance", metadataEndpoint))
	if err != nil {
		return metadataJSON{}, err
	}

	metadataParameters := metadataEndpoint.Query()
	metadataParameters.Add(apiVersionQueryParam, metadataAPIVersion)
	metadataEndpoint.RawQuery = metadataParameters.Encode()
	req, err := http.NewRequest(http.MethodGet, metadataEndpoint.String(), nil)
	if err != nil {
		return metadataJSON{}, fmt.Errorf("error creating HTTP Request for metadata endpoint: %w", err)
	}
	req.Header.Add("Metadata", "true")

	client := &http.Client{
		Timeout: clientTimeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return metadataJSON{}, fmt.Errorf("error calling metadata endpoint: %w", err)
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return metadataJSON{}, fmt.Errorf("error reading response body from metadata endpoint: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp errorJSON
		_ = json.Unmarshal(responseBytes, &errResp)
		if err != nil {
			return metadataJSON{}, fmt.Errorf("received error message but was unable to unmarshal its contents")
		}
		return metadataJSON{}, fmt.Errorf("%s error from metadata endpoint: %s", errResp.Error, errResp.ErrorDescription)
	}

	var r metadataJSON
	err = json.Unmarshal(responseBytes, &r)
	if err != nil {
		return metadataJSON{}, fmt.Errorf("error unmarshaling the response from metadata endpoint: %w", err)
	}

	return r, nil
}
