package azureauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	oidc "github.com/coreos/go-oidc"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"golang.org/x/oauth2"
)

type computeClient interface {
	Get(ctx context.Context, resourceGroup, vmName string, instanceView compute.InstanceViewTypes) (compute.VirtualMachine, error)
}

type vmssClient interface {
	Get(ctx context.Context, resourceGroup, vmssName string) (compute.VirtualMachineScaleSet, error)
}

type tokenVerifier interface {
	Verify(ctx context.Context, token string) (*oidc.IDToken, error)
}

type provider interface {
	Verifier() tokenVerifier
	ComputeClient(subscriptionID string) computeClient
	VMSSClient(subscriptionID string) vmssClient
}

type azureProvider struct {
	oidcVerifier *oidc.IDTokenVerifier
	authorizer   autorest.Authorizer
	httpClient   *http.Client
}

type oidcDiscoveryInfo struct {
	Issuer  string `json:"issuer"`
	JWKSURL string `json:"jwks_uri"`
}

func newAzureProvider(config *azureConfig) (*azureProvider, error) {
	httpClient := cleanhttp.DefaultClient()
	settings, err := getAzureSettings(config)
	if err != nil {
		return nil, err
	}

	// In many OIDC providers, the discovery endpoint matches the issuer. For Azure AD, the discovery
	// endpoint is the AD endpoint which does not match the issuer defined in the discovery payload. This
	// makes a request to the discovery URL to determine the issuer and key set information to configure
	// the OIDC verifier
	discoveryURL := fmt.Sprintf("%s%s/.well-known/openid-configuration", settings.Environment.ActiveDirectoryEndpoint, settings.TenantID)
	req, err := http.NewRequest("GET", discoveryURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent())

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errwrap.Wrapf("unable to read response body: {{err}}", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}
	var discoveryInfo oidcDiscoveryInfo
	if err := json.Unmarshal(body, &discoveryInfo); err != nil {
		return nil, errwrap.Wrapf("unable to unmarshal discovery url: {{err}}", err)
	}

	// Create a remote key set from the discovery endpoint
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	remoteKeySet := oidc.NewRemoteKeySet(ctx, discoveryInfo.JWKSURL)

	verifierConfig := &oidc.Config{
		ClientID:             settings.Resource,
		SupportedSigningAlgs: []string{oidc.RS256},
	}
	oidcVerifier := oidc.NewVerifier(discoveryInfo.Issuer, remoteKeySet, verifierConfig)

	// Create an OAuth2 client for retrieving VM data
	var authorizer autorest.Authorizer
	switch {
	// Use environment/config first
	case settings.ClientSecret != "":
		config := auth.NewClientCredentialsConfig(settings.ClientID, settings.ClientSecret, settings.TenantID)
		config.AADEndpoint = settings.Environment.ActiveDirectoryEndpoint
		config.Resource = settings.Environment.ResourceManagerEndpoint
		authorizer, err = config.Authorizer()
		if err != nil {
			return nil, err
		}
	// By default use MSI
	default:
		config := auth.NewMSIConfig()
		config.Resource = settings.Environment.ResourceManagerEndpoint
		authorizer, err = config.Authorizer()
		if err != nil {
			return nil, err
		}
	}

	// Ping the metadata service (if available)
	go pingMetadataService()

	return &azureProvider{
		authorizer:   authorizer,
		oidcVerifier: oidcVerifier,
		httpClient:   httpClient,
	}, nil
}

func (p *azureProvider) Verifier() tokenVerifier {
	return p.oidcVerifier
}

func (p *azureProvider) ComputeClient(subscriptionID string) computeClient {
	client := compute.NewVirtualMachinesClient(subscriptionID)
	client.Authorizer = p.authorizer
	client.Sender = p.httpClient
	client.AddToUserAgent(userAgent())
	return client
}

func (p *azureProvider) VMSSClient(subscriptionID string) vmssClient {
	client := compute.NewVirtualMachineScaleSetsClient(subscriptionID)
	client.Authorizer = p.authorizer
	client.Sender = p.httpClient
	client.AddToUserAgent(userAgent())
	return client
}

type azureSettings struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	Environment  azure.Environment
	Resource     string
}

func getAzureSettings(config *azureConfig) (*azureSettings, error) {
	settings := new(azureSettings)

	envTenantID := os.Getenv("AZURE_TENANT_ID")
	switch {
	case envTenantID != "":
		settings.TenantID = envTenantID
	case config.TenantID != "":
		settings.TenantID = config.TenantID
	default:
		return nil, errors.New("tenant_id is required")
	}

	envResource := os.Getenv("AZURE_AD_RESOURCE")
	switch {
	case envResource != "":
		settings.Resource = envResource
	case config.Resource != "":
		settings.Resource = config.Resource
	default:
		return nil, errors.New("resource is required")
	}

	clientID := os.Getenv("AZURE_CLIENT_ID")
	if clientID == "" {
		clientID = config.ClientID
	}
	settings.ClientID = clientID

	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if clientSecret == "" {
		clientSecret = config.ClientSecret
	}
	settings.ClientSecret = clientSecret

	envName := os.Getenv("AZURE_ENVIRONMENT")
	if envName == "" {
		envName = config.Environment
	}
	if envName == "" {
		settings.Environment = azure.PublicCloud
	} else {
		var err error
		settings.Environment, err = azure.EnvironmentFromName(envName)
		if err != nil {
			return nil, err
		}
	}

	return settings, nil
}

// This is simply to ping the Azure metadata service, if it is running
// in Azure
func pingMetadataService() {
	client := cleanhttp.DefaultClient()
	client.Timeout = 5 * time.Second
	req, _ := http.NewRequest("GET", "http://169.254.169.254/metadata/instance", nil)
	req.Header.Add("Metadata", "True")
	req.Header.Set("User-Agent", userAgent())

	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("api-version", "2017-04-02")
	req.URL.RawQuery = q.Encode()

	client.Do(req)
}
