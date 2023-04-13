// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azure

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/vault/helper/useragent"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	az "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
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

	role     string
	resource string
	objectID string
	clientID string
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

	token, err := getManagedIdentityCredentialToken(ctx, a.resource, a.objectID, a.clientID)
	if err != nil {
		retErr = err
		return
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

func getManagedIdentityCredentialToken(ctx context.Context, resource, objectID, clientID string) (string, error) {
	opts := &az.ManagedIdentityCredentialOptions{}
	if objectID != "" {
		opts.ID = az.ResourceID(objectID)
	}
	if clientID != "" {
		opts.ID = az.ClientID(clientID)
	}

	cred, err := az.NewManagedIdentityCredential(opts)
	if err != nil {
		return "", err
	}
	tokenOpts := policy.TokenRequestOptions{Scopes: []string{resource}}
	tk, err := cred.GetToken(ctx, tokenOpts)
	if err != nil {
		return "", err
	}
	return tk.Token, nil
}

func getInstanceMetadataInfo(ctx context.Context) ([]byte, error) {
	endpoint := instanceEndpoint
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("api-version", apiVersion)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Metadata", "true")
	req.Header.Set("User-Agent", useragent.AgentString())
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
