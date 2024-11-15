// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azureauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/pluginidentityutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"

	"github.com/hashicorp/vault-plugin-auth-azure/client"
)

// https://learn.microsoft.com/en-us/graph/sdks/national-clouds
const (
	azurePublicCloudBaseURI = "https://graph.microsoft.com"
	azureChinaCloudBaseURI  = "https://microsoftgraph.chinacloudapi.cn"
	azureUSGovCloudBaseURI  = "https://graph.microsoft.us"
	azurePublicCloudEnvName = "AZUREPUBLICCLOUD"
	azureChinaCloudEnvName  = "AZURECHINACLOUD"
	azureUSGovCloudEnvName  = "AZUREUSGOVERNMENTCLOUD"
)

type provider interface {
	TokenVerifier() client.TokenVerifier
	ComputeClient(subscriptionID string) (client.ComputeClient, error)
	VMSSClient(subscriptionID string) (client.VMSSClient, error)
	MSIClient(subscriptionID string) (client.MSIClient, error)
	MSGraphClient() (client.MSGraphClient, error)
	ResourceClient(subscriptionID string) (client.ResourceClient, error)
	ProvidersClient(subscriptionID string) (client.ProvidersClient, error)
}

type azureProvider struct {
	oidcVerifier *oidc.IDTokenVerifier
	settings     *azureSettings
	httpClient   *http.Client
	logger       hclog.Logger
	systemView   logical.SystemView
}

type oidcDiscoveryInfo struct {
	Issuer  string `json:"issuer"`
	JWKSURL string `json:"jwks_uri"`
}

// transporter implements the azure exported.Transporter interface to send HTTP
// requests. This allows us to set our custom http client and user agent.
type transporter struct {
	pluginEnv *logical.PluginEnvironment
	sender    *http.Client
}

func (tp transporter) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", useragent.PluginString(tp.pluginEnv,
		userAgentPluginName))

	client := tp.sender

	// don't attempt redirects so we aren't acting as an unintended network proxy
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *azureAuthBackend) newAzureProvider(ctx context.Context, config *azureConfig) (*azureProvider, error) {
	httpClient := cleanhttp.DefaultClient()
	settings, err := b.getAzureSettings(ctx, config)
	if err != nil {
		return nil, err
	}

	// In many OIDC providers, the discovery endpoint matches the issuer. For Azure AD, the discovery
	// endpoint is the AD endpoint which does not match the issuer defined in the discovery payload. This
	// makes a request to the discovery URL to determine the issuer and key set information to configure
	// the OIDC verifier
	discoveryURL := fmt.Sprintf("%s%s/.well-known/openid-configuration", settings.CloudConfig.ActiveDirectoryAuthorityHost, settings.TenantID)
	req, err := http.NewRequestWithContext(ctx, "GET", discoveryURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", useragent.PluginString(settings.PluginEnv,
		userAgentPluginName))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}
	var discoveryInfo oidcDiscoveryInfo
	if err := json.Unmarshal(body, &discoveryInfo); err != nil {
		return nil, fmt.Errorf("unable to unmarshal discovery url: %w", err)
	}

	// Create a remote key set from the discovery endpoint
	keySetCtx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	remoteKeySet := oidc.NewRemoteKeySet(keySetCtx, discoveryInfo.JWKSURL)

	verifierConfig := &oidc.Config{
		ClientID:             settings.Resource,
		SupportedSigningAlgs: []string{oidc.RS256},
	}
	oidcVerifier := oidc.NewVerifier(discoveryInfo.Issuer, remoteKeySet, verifierConfig)

	return &azureProvider{
		settings:     settings,
		oidcVerifier: oidcVerifier,
		httpClient:   httpClient,
		logger:       b.Logger(),
		systemView:   b.System(),
	}, nil
}

func (p *azureProvider) TokenVerifier() client.TokenVerifier {
	return p.oidcVerifier
}

func (p *azureProvider) MSGraphClient() (client.MSGraphClient, error) {
	cred, err := p.getTokenCredential()
	if err != nil {
		return nil, err
	}

	msGraphAppClient, err := client.NewMSGraphApplicationClient(p.settings.GraphURI, cred)
	if err != nil {
		return nil, fmt.Errorf("failed to create MS graph client: %w", err)
	}

	return msGraphAppClient, nil
}

func (p *azureProvider) ComputeClient(subscriptionID string) (client.ComputeClient, error) {
	cred, err := p.getTokenCredential()
	if err != nil {
		return nil, err
	}

	clientOptions := p.getClientOptions()
	client, err := armcompute.NewVirtualMachinesClient(subscriptionID, cred, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual machines client: %w", err)
	}

	return client, nil
}

func (p *azureProvider) VMSSClient(subscriptionID string) (client.VMSSClient, error) {
	cred, err := p.getTokenCredential()
	if err != nil {
		return nil, err
	}

	clientOptions := p.getClientOptions()
	client, err := armcompute.NewVirtualMachineScaleSetsClient(subscriptionID, cred, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create virtual machine scale sets client: %w", err)
	}

	return client, nil
}

func (p *azureProvider) MSIClient(subscriptionID string) (client.MSIClient, error) {
	cred, err := p.getTokenCredential()
	if err != nil {
		return nil, err
	}

	clientOptions := p.getClientOptions()
	client, err := armmsi.NewUserAssignedIdentitiesClient(subscriptionID, cred, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create user assigned identity client: %w", err)
	}

	return client, nil
}

func (p *azureProvider) ProvidersClient(subscriptionID string) (client.ProvidersClient, error) {
	cred, err := p.getTokenCredential()
	if err != nil {
		return nil, err
	}

	clientOptions := p.getClientOptions()
	client, err := armresources.NewProvidersClient(subscriptionID, cred, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create providers client: %w", err)
	}

	return client, nil
}

func (p *azureProvider) ResourceClient(subscriptionID string) (client.ResourceClient, error) {
	cred, err := p.getTokenCredential()
	if err != nil {
		return nil, err
	}

	clientOptions := p.getClientOptions()
	client, err := armresources.NewClient(subscriptionID, cred, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource client: %w", err)
	}

	return client, nil
}

func (p *azureProvider) getClientOptions() *arm.ClientOptions {
	return &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: p.settings.CloudConfig,
			Transport: transporter{
				pluginEnv: p.settings.PluginEnv,
				sender:    p.httpClient,
			},
			Retry: policy.RetryOptions{
				MaxRetries:    p.settings.MaxRetries,
				MaxRetryDelay: p.settings.MaxRetryDelay,
				RetryDelay:    p.settings.RetryDelay,
			},
		},
	}
}

func (p *azureProvider) getTokenCredential() (azcore.TokenCredential, error) {
	clientCloudOpts := azcore.ClientOptions{Cloud: p.settings.CloudConfig}

	if p.settings.ClientSecret != "" {
		options := &azidentity.ClientSecretCredentialOptions{
			ClientOptions: clientCloudOpts,
		}

		cred, err := azidentity.NewClientSecretCredential(p.settings.TenantID, p.settings.ClientID,
			p.settings.ClientSecret, options)
		if err != nil {
			return nil, fmt.Errorf("failed to create client secret token credential: %w", err)
		}

		return cred, nil
	}

	if p.settings.IdentityTokenAudience != "" {
		options := &azidentity.ClientAssertionCredentialOptions{
			ClientOptions: clientCloudOpts,
		}
		getAssertion := getAssertionFunc(p.logger, p.systemView, p.settings)
		cred, err := azidentity.NewClientAssertionCredential(
			p.settings.TenantID,
			p.settings.ClientID,
			getAssertion,
			options,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create client assertion credential: %w", err)
		}

		return cred, nil
	}

	// Fall back to using managed service identity
	options := &azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: clientCloudOpts,
	}
	cred, err := azidentity.NewManagedIdentityCredential(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create managed identity token credential: %w", err)
	}

	return cred, nil
}

type getAssertion func(context.Context) (string, error)

func getAssertionFunc(logger hclog.Logger, sys logical.SystemView, s *azureSettings) getAssertion {
	return func(ctx context.Context) (string, error) {
		req := &pluginutil.IdentityTokenRequest{
			Audience: s.IdentityTokenAudience,
			TTL:      s.IdentityTokenTTL * time.Second,
		}
		resp, err := sys.GenerateIdentityToken(ctx, req)
		if err != nil {
			return "", fmt.Errorf("failed to generate plugin identity token: %w", err)
		}
		logger.Info("fetched new plugin identity token")

		if resp.TTL < req.TTL {
			logger.Debug("generated plugin identity token has shorter TTL than requested",
				"requested", req.TTL, "actual", resp.TTL)
		}

		return resp.Token.Token(), nil
	}
}

type azureSettings struct {
	pluginidentityutil.PluginIdentityTokenParams

	TenantID      string
	ClientID      string
	ClientSecret  string
	CloudConfig   cloud.Configuration
	GraphURI      string
	Resource      string
	PluginEnv     *logical.PluginEnvironment
	MaxRetries    int32
	MaxRetryDelay time.Duration
	RetryDelay    time.Duration
}

func (b *azureAuthBackend) getAzureSettings(ctx context.Context, config *azureConfig) (*azureSettings, error) {
	settings := &azureSettings{
		MaxRetries:    config.MaxRetries,
		MaxRetryDelay: config.MaxRetryDelay,
		RetryDelay:    config.RetryDelay,
	}

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

	settings.IdentityTokenAudience = config.IdentityTokenAudience
	settings.IdentityTokenTTL = config.IdentityTokenTTL

	environment := os.Getenv("AZURE_ENVIRONMENT")
	if environment == "" {
		// set environment from config
		environment = config.Environment
	}
	if environment == "" {
		// Default to Azure public cloud
		settings.CloudConfig = cloud.AzurePublic
		settings.GraphURI = azurePublicCloudBaseURI
	} else {
		var err error
		settings.CloudConfig, err = cloudConfigFromName(environment)
		if err != nil {
			return nil, err
		}

		settings.GraphURI, err = graphURIFromName(environment)
		if err != nil {
			return nil, err
		}
	}

	pluginEnv, err := b.System().PluginEnv(ctx)
	if err != nil {
		b.Logger().Warn("failed to read plugin environment, user-agent will not be set",
			"error", err)
	}
	settings.PluginEnv = pluginEnv

	return settings, nil
}

func cloudConfigFromName(name string) (cloud.Configuration, error) {
	configs := map[string]cloud.Configuration{
		azureChinaCloudEnvName:  cloud.AzureChina,
		azurePublicCloudEnvName: cloud.AzurePublic,
		azureUSGovCloudEnvName:  cloud.AzureGovernment,
	}

	name = strings.ToUpper(name)
	c, ok := configs[name]
	if !ok {
		return c, fmt.Errorf("err: no cloud configuration matching the name %q", name)
	}

	return c, nil
}

func graphURIFromName(name string) (string, error) {
	configs := map[string]string{
		azureChinaCloudEnvName:  azureChinaCloudBaseURI,
		azurePublicCloudEnvName: azurePublicCloudBaseURI,
		azureUSGovCloudEnvName:  azureUSGovCloudBaseURI,
	}

	name = strings.ToUpper(name)
	c, ok := configs[name]
	if !ok {
		return c, fmt.Errorf("err: no MS Graph URI matching the name %q", name)
	}

	return c, nil
}

// guidRx from https://learn.microsoft.com/en-us/rest/api/defenderforcloud/tasks/get-subscription-level-task
var guidRx = regexp.MustCompile(`^[0-9A-Fa-f]{8}-([0-9A-Fa-f]{4}-){3}[0-9A-Fa-f]{12}$`) // just a uuid
// nameRx based on https://azure.github.io/PSRule.Rules.Azure/en/rules/Azure.VM.Name/#description
var nameRx = regexp.MustCompile(`^[a-zA-Z]$|^[a-zA-Z][a-zA-Z0-9.\-_]*[a-zA-Z0-9_]$`) // alphanumeric, doesn't start with a number, at least 1 character, doesn't end with a . or -
// https://azure.github.io/PSRule.Rules.Azure/en/rules/Azure.ResourceGroup.Name/ and https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules
// The latter documentation specifically allows characters in unicode letter/digit categories, which is wider than a-zA-Z0-9.
var rgRx = regexp.MustCompile(`^[\-_.()\pL\pN]*[\-_()\pL\pN]$`)

// verify the field provided matches Azure's requirements
// (see: https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules).
func validateAzureField(regex *regexp.Regexp, value string) bool {
	return regex.MatchString(value)
}
