// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/vault/api"

	"github.com/oracle/oci-go-sdk/v59/common"
	"github.com/oracle/oci-go-sdk/v59/common/auth"
)

type CLIHandler struct{}

func (h *CLIHandler) Help() string {
	help := `
Usage: vault login -method=oci auth_type=apikey 
       vault login -method=oci auth_type=instance 

  The OCI auth method allows users to authenticate with OCI
  credentials. The OCI credentials may be specified in a number of ways,
  listed below:

    1. API Key

    2. Instance Principal

  Authenticate using API key:

		First create a configuration file as explained in https://docs.us-phoenix-1.oraclecloud.com/Content/API/Concepts/sdkconfig.htm
		Then login using the following command:

		$ vault login -method=oci auth_type=apikey role=<RoleName>

  Authenticate using Instance Principal:
		https://docs.cloud.oracle.com/iaas/Content/Identity/Tasks/callingservicesfrominstances.htm
		
		$ vault login -method=oci auth_type=instance role=<RoleName>

Configuration:
  auth_type=<string>
      Enter one of following: 
		apikey (or) ak		
		instance (or) ip
`
	return strings.TrimSpace(help)
}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (*api.Secret, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "oci"
	}
	mount = strings.TrimSuffix(mount, "/")

	role, ok := m["role"]
	if !ok {
		return nil, fmt.Errorf("'role' is required")
	}
	role = strings.ToLower(role)

	path := fmt.Sprintf(PathBaseFormat, mount, role)
	signingPath := PathVersionBase + path

	data, err := CreateLoginData(c.Address(), m, signingPath)
	if err != nil {
		return nil, err
	}

	// Now try to login
	secret, err := c.Logical().Write(path, data)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// CreateLoginData creates the interface required for a login request, signed using the corresponding OCI Identity Principal
func CreateLoginData(addr string, m map[string]string, path string) (map[string]interface{}, error) {
	authType, ok := m["auth_type"]
	if !ok {
		return nil, fmt.Errorf("'auth_type' is required")
	}

	var headerFunc func(string, string) (http.Header, error)
	switch strings.ToLower(authType) {
	case "ip", "instance":
		headerFunc = GetSignedInstanceRequestHeaders
	case "ak", "apikey":
		headerFunc = GetSignedAPIRequestHeaders
	default:
		return nil, fmt.Errorf("unsupported auth_type %q", authType)
	}

	headers, err := headerFunc(addr, path)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"request_headers": headers,
	}, nil
}

func GetSignedInstanceRequestHeaders(addr, path string) (http.Header, error) {
	ip, err := auth.InstancePrincipalConfigurationProvider()
	if err != nil {
		return nil, err
	}

	c, err := NewOciClientWithConfigurationProvider(ip)
	if err != nil {
		return nil, err
	}
	return getSignedRequestHeaders(addr, &c, path)
}

func GetSignedAPIRequestHeaders(addr, path string) (http.Header, error) {
	c, err := NewOciClientWithConfigurationProvider(common.DefaultConfigProvider())
	if err != nil {
		return nil, err
	}

	return getSignedRequestHeaders(addr, &c, path)
}

func getSignedRequestHeaders(addr string, client *OciClient, path string) (http.Header, error) {
	clientURL, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	client.Host = addr
	request, err := client.ConstructLoginRequest(path)
	if err != nil {
		return nil, err
	}

	request.Host = clientURL.Host
	request.Header.Set("host", request.Host)

	// ref: https://tools.ietf.org/html/draft-cavage-http-signatures-06#section-2.3
	request.Header.Set("(request-target)",
		fmt.Sprintf("%s %s", strings.ToLower(request.Method), request.URL.RequestURI()))

	return request.Header, nil
}
