package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
)

var securityCredURL = "http://100.100.100.200/latest/meta-data/ram/security-credentials/"

func NewInstanceMetadataProvider() Provider {
	return &InstanceMetadataProvider{}
}

type InstanceMetadataProvider struct {
	RoleName string
}

func (p *InstanceMetadataProvider) Retrieve() (auth.Credential, error) {
	if p.RoleName == "" {
		// Instances can have only one role name that never changes,
		// so attempt to populate it.
		// If this call is executed in an environment that doesn't support instance metadata,
		// it will time out after 30 seconds and return an err.
		resp, err := http.Get(securityCredURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("received %d getting role name: %s", resp.StatusCode, bodyBytes)
		}
		roleName := string(bodyBytes)
		if roleName == "" {
			return nil, errors.New("unable to retrieve role name, it may be unset")
		}
		p.RoleName = roleName
	}

	resp, err := http.Get(securityCredURL + p.RoleName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received %d getting security credentials for %s", resp.StatusCode, p.RoleName)
	}
	body := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	accessKeyID, err := extractString(body, "AccessKeyId")
	if err != nil {
		return nil, err
	}
	accessKeySecret, err := extractString(body, "AccessKeySecret")
	if err != nil {
		return nil, err
	}
	securityToken, err := extractString(body, "SecurityToken")
	if err != nil {
		return nil, err
	}
	return credentials.NewStsTokenCredential(accessKeyID, accessKeySecret, securityToken), nil
}

func extractString(m map[string]interface{}, key string) (string, error) {
	raw, ok := m[key]
	if !ok {
		return "", fmt.Errorf("%s not in %+v", key, m)
	}
	str, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("%s is not a string in %+v", key, m)
	}
	return str, nil
}
