package azure

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	"github.com/hashicorp/vault/helper/jsonutil"
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

	switch {
	case a.role == "":
		return nil, errors.New("'role' value is empty")
	case a.resource == "":
		return nil, errors.New("'resource' value is empty")
	}

	return a, nil
}

func (a *azureMethod) Authenticate(client *api.Client) (*api.Secret, error) {
	a.logger.Trace("beginning authentication")

	// Fetch instance data
	var instance struct {
		Compute struct {
			Name              string
			ResourceGroupName string
			SubscriptionID    string
		}
	}

	body, err := getMetadataInfo(instanceEndpoint, "")
	if err != nil {
		return nil, err
	}

	err = jsonutil.DecodeJSON(body, &instance)
	if err != nil {
		return nil, errwrap.Wrapf("error parsing instance metadata response: {{err}}", err)
	}

	// Fetch JWT
	var identity struct {
		AccessToken string `json:"access_token"`
	}

	body, err = getMetadataInfo(identityEndpoint, a.resource)
	if err != nil {
		return nil, err
	}

	err = jsonutil.DecodeJSON(body, &identity)
	if err != nil {
		return nil, errwrap.Wrapf("error parsing identity metadata response: {{err}}", err)
	}

	// Attempt login
	data := map[string]interface{}{
		"role":                a.role,
		"vm_name":             instance.Compute.Name,
		"resource_group_name": instance.Compute.ResourceGroupName,
		"subscription_id":     instance.Compute.SubscriptionID,
		"jwt":                 identity.AccessToken,
	}

	secret, err := client.Logical().Write(fmt.Sprintf("%s/login", a.mountPath), data)
	if err != nil {
		return nil, errwrap.Wrapf("error logging in: {{err}}", err)
	}

	return secret, nil
}

func (a *azureMethod) NewCreds() chan struct{} {
	return nil
}

func (a *azureMethod) Shutdown() {
}

func getMetadataInfo(endpoint, resource string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("api-version", apiVersion)
	if resource != "" {
		q.Add("resource", resource)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Metadata", "true")

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("error fetching metadata from %s: {{err}}", endpoint), err)
	}

	if resp == nil {
		return nil, fmt.Errorf("empty response fetching metadata from %s", endpoint)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("error reading metadata from %s: {{err}}", endpoint), err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response in metadata from %s: %s", endpoint, body)
	}

	return body, nil
}
