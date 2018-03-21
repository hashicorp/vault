package plugin

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	oidc "github.com/coreos/go-oidc"
)

const (
	issuerBaseURI = "https://sts.windows.net"
)

type computeClient interface {
	Get(ctx context.Context, resourceGroup, vmName string, instanceView compute.InstanceViewTypes) (compute.VirtualMachine, error)
}

type tokenVerifier interface {
	Verify(ctx context.Context, token string) (*oidc.IDToken, error)
}

type provider interface {
	Verifier() tokenVerifier
	ComputeClient(subscriptionID string) computeClient
}

var _ provider = &azureProvider{}

type azureProvider struct {
	settings     *azureSettings
	oidcProvider *oidc.Provider
	authorizer   autorest.Authorizer
}

func NewAzureProvider(config *azureConfig) (*azureProvider, error) {
	settings, err := getAzureSettings(config)
	if err != nil {
		return nil, err
	}

	issuer := fmt.Sprintf("%s/%s/", issuerBaseURI, settings.TenantID)
	oidcProvider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return nil, err
	}

	provider := &azureProvider{
		settings:     settings,
		oidcProvider: oidcProvider,
	}

	// OAuth2 client for querying VM data
	switch {
	// Use environment/config first
	case settings.ClientSecret != "":
		config := auth.NewClientCredentialsConfig(settings.ClientID, settings.ClientSecret, settings.TenantID)
		config.AADEndpoint = settings.Environment.ActiveDirectoryEndpoint
		config.Resource = settings.Environment.ResourceManagerEndpoint
		provider.authorizer, err = config.Authorizer()
		if err != nil {
			return nil, err
		}
	// By default use MSI
	default:
		config := auth.NewMSIConfig()
		config.Resource = settings.Environment.ResourceManagerEndpoint
		provider.authorizer, err = config.Authorizer()
		if err != nil {
			return nil, err
		}
	}
	return provider, nil
}

func (p *azureProvider) Verifier() tokenVerifier {
	verifierConfig := &oidc.Config{
		ClientID:             p.settings.Resource,
		SupportedSigningAlgs: []string{oidc.RS256},
	}
	return p.oidcProvider.Verifier(verifierConfig)
}

func (p *azureProvider) ComputeClient(subscriptionID string) computeClient {
	client := compute.NewVirtualMachinesClient(subscriptionID)
	client.Authorizer = p.authorizer
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
