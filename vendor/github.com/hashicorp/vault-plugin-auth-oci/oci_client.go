// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/oracle/oci-go-sdk/v59/common"
)

// OciClient stores the client and configuration details for making API requests to OCI Identity Service
type OciClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

// These constants store information related to signing the http request
const (
	// requestHeaderDate The key for passing a header to indicate Date
	requestHeaderDate = "Date"

	// requestHeaderUserAgent The key for passing a header to indicate User Agent
	requestHeaderUserAgent = "User-Agent"

	defaultScheme = "https"
)

// NewIdentityClientWithConfigurationProvider Creates a new default Identity client with the given configuration provider.
// the configuration provider will be used for the default signer as well as reading the region
func NewOciClientWithConfigurationProvider(configProvider common.ConfigurationProvider) (client OciClient, err error) {
	baseClient, err := common.NewClientWithConfig(configProvider)
	if err != nil {
		return client, err
	}

	client = OciClient{BaseClient: baseClient}
	err = client.setConfigurationProvider(configProvider)
	return client, err
}

// SetConfigurationProvider sets the configuration provider including the region, returns an error if is not valid
func (client *OciClient) setConfigurationProvider(configProvider common.ConfigurationProvider) error {
	if ok, err := common.IsConfigurationProviderValid(configProvider); !ok {
		return err
	}

	// Error has been checked already
	client.config = &configProvider
	return nil
}

// ConstructLoginRequest takes in a path and returns a signed http request
func (client OciClient) ConstructLoginRequest(path string) (request http.Request, err error) {
	httpRequest, err := common.MakeDefaultHTTPRequestWithTaggedStruct(http.MethodGet, path, request)
	if err != nil {
		return
	}

	err = client.prepareRequest(&httpRequest)
	if err != nil {
		return
	}

	err = client.Signer.Sign(&httpRequest)
	if err != nil {
		return
	}

	request = httpRequest
	return
}

// prepareRequest takes in a http request and adds the required information for signing it
func (client *OciClient) prepareRequest(request *http.Request) (err error) {
	if client.UserAgent == "" {
		return errors.New("user agent can not be blank")
	}

	if request.Header == nil {
		request.Header = http.Header{}
	}
	request.Header.Set(requestHeaderUserAgent, client.UserAgent)
	request.Header.Set(requestHeaderDate, time.Now().UTC().Format(http.TimeFormat))

	if !strings.HasPrefix(client.Host, "http://") &&
		!strings.HasPrefix(client.Host, "https://") {
		client.Host = fmt.Sprintf("%s://%s", defaultScheme, client.Host)
	}

	clientURL, err := url.Parse(client.Host)
	if err != nil {
		return errwrap.Wrapf("host is invalid. {{err}}", err)
	}
	request.URL.Host = clientURL.Host
	request.URL.Scheme = clientURL.Scheme
	currentPath := request.URL.Path
	if !strings.Contains(currentPath, fmt.Sprintf("/%s", client.BasePath)) {
		request.URL.Path = path.Clean(fmt.Sprintf("/%s/%s", client.BasePath, currentPath))
	}
	return
}
