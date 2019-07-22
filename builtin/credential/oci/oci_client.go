// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type OciClient struct {
	common.BaseClient
	config *common.ConfigurationProvider
}

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
		return
	}

	client = OciClient{BaseClient: baseClient}
	client.BasePath = ""
	err = client.setConfigurationProvider(configProvider)
	return
}

// SetRegion overrides the region of this client.
func (client *OciClient) SetHost(host string) {
	client.Host = host
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
	err = nil
	return
}

func (client *OciClient) prepareRequest(request *http.Request) (err error) {
	if client.UserAgent == "" {
		return fmt.Errorf("user agent can not be blank")
	}

	if request.Header == nil {
		request.Header = http.Header{}
	}
	request.Header.Set(requestHeaderUserAgent, client.UserAgent)
	request.Header.Set(requestHeaderDate, time.Now().UTC().Format(http.TimeFormat))

	if !strings.Contains(client.Host, "http") &&
		!strings.Contains(client.Host, "https") {
		client.Host = fmt.Sprintf("%s://%s", defaultScheme, client.Host)
	}

	clientURL, err := url.Parse(client.Host)
	if err != nil {
		return fmt.Errorf("host is invalid. %s", err.Error())
	}
	request.URL.Host = clientURL.Host
	request.URL.Scheme = clientURL.Scheme
	currentPath := request.URL.Path
	if !strings.Contains(currentPath, fmt.Sprintf("/%s", client.BasePath)) {
		request.URL.Path = path.Clean(fmt.Sprintf("/%s/%s", client.BasePath, currentPath))
	}
	return
}

// getRequestTarget returns the value of the special (request-target) header field name
// per https://tools.ietf.org/html/draft-cavage-http-signatures-06#section-2.3
func getRequestTarget(request *http.Request) string {
	lowercaseMethod := strings.ToLower(request.Method)
	return fmt.Sprintf("%s %s", lowercaseMethod, request.URL.RequestURI())
}
