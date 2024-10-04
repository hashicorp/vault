// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package azure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	policy "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	az "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	"github.com/hashicorp/vault/helper/useragent"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
)

const (
	instanceEndpoint = "http://169.254.169.254/metadata/instance"
	identityEndpoint = "http://169.254.169.254/metadata/identity/oauth2/token"

	// minimum version 2018-02-01 needed for identity metadata
	// regional availability: https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service
	apiVersion = "2018-02-01"
)

type azureMethod struct {
	logger    hclog.Logger
	mountPath string

	authenticateFromEnvironment bool
	role                        string
	scope                       string
	resource                    string
	objectID                    string
	clientID                    string
}

func NewAzureAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	a := &azureMethod{
		logger:    conf.Logger,
		mountPath: conf.MountPath,
	}

	roleRaw, ok := conf.Config["role"]
	if !ok {
		return nil, errors.New("missing 'role' value")
	}
	a.role, ok = roleRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role' config value to string")
	}

	resourceRaw, ok := conf.Config["resource"]
	if !ok {
		return nil, errors.New("missing 'resource' value")
	}
	a.resource, ok = resourceRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'resource' config value to string")
	}

	objectIDRaw, ok := conf.Config["object_id"]
	if ok {
		a.objectID, ok = objectIDRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'object_id' config value to string")
		}
	}

	clientIDRaw, ok := conf.Config["client_id"]
	if ok {
		a.clientID, ok = clientIDRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'client_id' config value to string")
		}
	}

	scopeRaw, ok := conf.Config["scope"]
	if ok {
		a.scope, ok = scopeRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'scope' config value to string")
		}
	}
	if a.scope == "" {
		a.scope = fmt.Sprintf("%s/.default", a.resource)
	}

	authenticateFromEnvironmentRaw, ok := conf.Config["authenticate_from_environment"]
	if ok {
		authenticateFromEnvironment, err := parseutil.ParseBool(authenticateFromEnvironmentRaw)
		if err != nil {
			return nil, fmt.Errorf("could not convert 'authenticate_from_environment' config value to bool: %w", err)
		}
		a.authenticateFromEnvironment = authenticateFromEnvironment
	}

	switch {
	case a.role == "":
		return nil, errors.New("'role' value is empty")
	case a.resource == "":
		return nil, errors.New("'resource' value is empty")
	case a.objectID != "" && a.clientID != "":
		return nil, errors.New("only one of 'object_id' or 'client_id' may be provided")
	}

	return a, nil
}

func (a *azureMethod) Authenticate(ctx context.Context, client *api.Client) (retPath string, header http.Header, retData map[string]interface{}, retErr error) {
	a.logger.Trace("beginning authentication")

	// Fetch instance data
	var instance struct {
		Compute struct {
			Name              string
			ResourceGroupName string
			SubscriptionID    string
			VMScaleSetName    string
			ResourceID        string
		}
	}

	body, err := getInstanceMetadataInfo(ctx)
	if err != nil {
		retErr = err
		return
	}

	err = jsonutil.DecodeJSON(body, &instance)
	if err != nil {
		retErr = fmt.Errorf("error parsing instance metadata response: %w", err)
		return
	}

	token := ""
	if a.authenticateFromEnvironment {
		token, err = getAzureTokenFromEnvironment(ctx, a.scope)
		if err != nil {
			retErr = err
			return
		}
	} else {
		token, err = getTokenFromIdentityEndpoint(ctx, a.resource, a.objectID, a.clientID)
		if err != nil {
			retErr = err
			return
		}
	}

	// Attempt login
	data := map[string]interface{}{
		"role":                a.role,
		"vm_name":             instance.Compute.Name,
		"vmss_name":           instance.Compute.VMScaleSetName,
		"resource_group_name": instance.Compute.ResourceGroupName,
		"subscription_id":     instance.Compute.SubscriptionID,
		"jwt":                 token,
	}

	return fmt.Sprintf("%s/login", a.mountPath), nil, data, nil
}

func (a *azureMethod) NewCreds() chan struct{} {
	return nil
}

func (a *azureMethod) CredSuccess() {
}

func (a *azureMethod) Shutdown() {
}

// getAzureTokenFromEnvironment Is Azure's preferred way for authentication, and takes values
// from environment variables to form a credential.
// It uses a DefaultAzureCredential:
// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-defaultazurecredential
// Environment variables are taken into account in the following order:
// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-environment-variables
func getAzureTokenFromEnvironment(ctx context.Context, scope string) (string, error) {
	cred, err := az.NewDefaultAzureCredential(nil)
	if err != nil {
		return "", err
	}

	tokenOpts := policy.TokenRequestOptions{Scopes: []string{scope}}
	tk, err := cred.GetToken(ctx, tokenOpts)
	if err != nil {
		return "", err
	}
	return tk.Token, nil
}

// getInstanceMetadataInfo calls the Azure Instance Metadata endpoint to get
// information about the Azure environment it's running in.
func getInstanceMetadataInfo(ctx context.Context) ([]byte, error) {
	return getMetadataInfo(ctx, instanceEndpoint, "", "", "")
}

// getTokenFromIdentityEndpoint is kept for backwards compatibility purposes. Using the
// newer APIs and the Azure SDK should be preferred over this mechanism.
func getTokenFromIdentityEndpoint(ctx context.Context, resource, objectID, clientID string) (string, error) {
	var identity struct {
		AccessToken string `json:"access_token"`
	}

	body, err := getMetadataInfo(ctx, identityEndpoint, resource, objectID, clientID)
	if err != nil {
		return "", err
	}

	err = jsonutil.DecodeJSON(body, &identity)
	if err != nil {
		return "", fmt.Errorf("error parsing identity metadata response: %w", err)
	}

	return identity.AccessToken, nil
}

// getMetadataInfo calls the Azure metadata endpoint with the given parameters.
// An empty resource, objectID and clientID will return metadata information.
func getMetadataInfo(ctx context.Context, endpoint, resource, objectID, clientID string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("api-version", apiVersion)
	if resource != "" {
		q.Add("resource", resource)
	}
	if objectID != "" {
		q.Add("object_id", objectID)
	}
	if clientID != "" {
		q.Add("client_id", clientID)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Metadata", "true")
	req.Header.Set("User-Agent", useragent.String())
	req = req.WithContext(ctx)

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching metadata from %s: %w", endpoint, err)
	}

	if resp == nil {
		return nil, fmt.Errorf("empty response fetching metadata from %s", endpoint)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading metadata from %s: %w", endpoint, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response in metadata from %s: %s", endpoint, body)
	}

	return body, nil
}
