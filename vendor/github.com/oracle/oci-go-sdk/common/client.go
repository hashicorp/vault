// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.

// Package common provides supporting functions and structs used by service packages
package common

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

const (
	// DefaultHostURLTemplate The default url template for service hosts
	DefaultHostURLTemplate = "%s.%s.oraclecloud.com"

	// requestHeaderAccept The key for passing a header to indicate Accept
	requestHeaderAccept = "Accept"

	// requestHeaderAuthorization The key for passing a header to indicate Authorization
	requestHeaderAuthorization = "Authorization"

	// requestHeaderContentLength The key for passing a header to indicate Content Length
	requestHeaderContentLength = "Content-Length"

	// requestHeaderContentType The key for passing a header to indicate Content Type
	requestHeaderContentType = "Content-Type"

	// requestHeaderDate The key for passing a header to indicate Date
	requestHeaderDate = "Date"

	// requestHeaderIfMatch The key for passing a header to indicate If Match
	requestHeaderIfMatch = "if-match"

	// requestHeaderOpcClientInfo The key for passing a header to indicate OPC Client Info
	requestHeaderOpcClientInfo = "opc-client-info"

	// requestHeaderOpcRetryToken The key for passing a header to indicate OPC Retry Token
	requestHeaderOpcRetryToken = "opc-retry-token"

	// requestHeaderOpcRequestID The key for unique Oracle-assigned identifier for the request.
	requestHeaderOpcRequestID = "opc-request-id"

	// requestHeaderOpcClientRequestID The key for unique Oracle-assigned identifier for the request.
	requestHeaderOpcClientRequestID = "opc-client-request-id"

	// requestHeaderUserAgent The key for passing a header to indicate User Agent
	requestHeaderUserAgent = "User-Agent"

	// requestHeaderXContentSHA256 The key for passing a header to indicate SHA256 hash
	requestHeaderXContentSHA256 = "X-Content-SHA256"

	// requestHeaderOpcOboToken The key for passing a header to use obo token
	requestHeaderOpcOboToken = "opc-obo-token"

	// private constants
	defaultScheme            = "https"
	defaultSDKMarker         = "Oracle-GoSDK"
	defaultUserAgentTemplate = "%s/%s (%s/%s; go/%s)" //SDK/SDKVersion (OS/OSVersion; Lang/LangVersion)
	defaultTimeout           = 60 * time.Second
	defaultConfigFileName    = "config"
	defaultConfigDirName     = ".oci"
	configFilePathEnvVarName = "OCI_CONFIG_FILE"
	secondaryConfigDirName   = ".oraclebmc"
	maxBodyLenForDebug       = 1024 * 1000
)

// RequestInterceptor function used to customize the request before calling the underlying service
type RequestInterceptor func(*http.Request) error

// HTTPRequestDispatcher wraps the execution of a http request, it is generally implemented by
// http.Client.Do, but can be customized for testing
type HTTPRequestDispatcher interface {
	Do(req *http.Request) (*http.Response, error)
}

// BaseClient struct implements all basic operations to call oci web services.
type BaseClient struct {
	//HTTPClient performs the http network operations
	HTTPClient HTTPRequestDispatcher

	//Signer performs auth operation
	Signer HTTPRequestSigner

	//A request interceptor can be used to customize the request before signing and dispatching
	Interceptor RequestInterceptor

	//The host of the service
	Host string

	//The user agent
	UserAgent string

	//Base path for all operations of this client
	BasePath string
}

func defaultUserAgent() string {
	userAgent := fmt.Sprintf(defaultUserAgentTemplate, defaultSDKMarker, Version(), runtime.GOOS, runtime.GOARCH, runtime.Version())
	return userAgent
}

var clientCounter int64

func getNextSeed() int64 {
	newCounterValue := atomic.AddInt64(&clientCounter, 1)
	return newCounterValue + time.Now().UnixNano()
}

func newBaseClient(signer HTTPRequestSigner, dispatcher HTTPRequestDispatcher) BaseClient {
	rand.Seed(getNextSeed())
	return BaseClient{
		UserAgent:   defaultUserAgent(),
		Interceptor: nil,
		Signer:      signer,
		HTTPClient:  dispatcher,
	}
}

func defaultHTTPDispatcher() http.Client {
	httpClient := http.Client{
		Timeout: defaultTimeout,
	}
	return httpClient
}

func defaultBaseClient(provider KeyProvider) BaseClient {
	dispatcher := defaultHTTPDispatcher()
	signer := DefaultRequestSigner(provider)
	return newBaseClient(signer, &dispatcher)
}

//DefaultBaseClientWithSigner creates a default base client with a given signer
func DefaultBaseClientWithSigner(signer HTTPRequestSigner) BaseClient {
	dispatcher := defaultHTTPDispatcher()
	return newBaseClient(signer, &dispatcher)
}

// NewClientWithConfig Create a new client with a configuration provider, the configuration provider
// will be used for the default signer as well as reading the region
// This function does not check for valid regions to implement forward compatibility
func NewClientWithConfig(configProvider ConfigurationProvider) (client BaseClient, err error) {
	var ok bool
	if ok, err = IsConfigurationProviderValid(configProvider); !ok {
		err = fmt.Errorf("can not create client, bad configuration: %s", err.Error())
		return
	}

	client = defaultBaseClient(configProvider)

	return
}

// NewClientWithOboToken Create a new client that will use oboToken for auth
func NewClientWithOboToken(configProvider ConfigurationProvider, oboToken string) (client BaseClient, err error) {
	client, err = NewClientWithConfig(configProvider)
	if err != nil {
		return
	}

	// Interceptor to add obo token header
	client.Interceptor = func(request *http.Request) error {
		request.Header.Add(requestHeaderOpcOboToken, oboToken)
		return nil
	}
	// Obo token will also be signed
	defaultHeaders := append(DefaultGenericHeaders(), requestHeaderOpcOboToken)
	client.Signer = RequestSigner(configProvider, defaultHeaders, DefaultBodyHeaders())

	return
}

func getHomeFolder() string {
	current, e := user.Current()
	if e != nil {
		//Give up and try to return something sensible
		home := os.Getenv("HOME")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return current.HomeDir
}

// DefaultConfigProvider returns the default config provider. The default config provider
// will look for configurations in 3 places: file in $HOME/.oci/config, HOME/.obmcs/config and
// variables names starting with the string TF_VAR. If the same configuration is found in multiple
// places the provider will prefer the first one.
// If the config file is not placed in the default location, the environment variable
// OCI_CONFIG_FILE can provide the config file location.
func DefaultConfigProvider() ConfigurationProvider {
	defaultConfigFile := getDefaultConfigFilePath()
	homeFolder := getHomeFolder()
	secondaryConfigFile := path.Join(homeFolder, secondaryConfigDirName, defaultConfigFileName)

	defaultFileProvider, _ := ConfigurationProviderFromFile(defaultConfigFile, "")
	secondaryFileProvider, _ := ConfigurationProviderFromFile(secondaryConfigFile, "")
	environmentProvider := environmentConfigurationProvider{EnvironmentVariablePrefix: "TF_VAR"}

	provider, _ := ComposingConfigurationProvider([]ConfigurationProvider{defaultFileProvider, secondaryFileProvider, environmentProvider})
	Debugf("Configuration provided by: %s", provider)
	return provider
}

func getDefaultConfigFilePath() string {
	homeFolder := getHomeFolder()
	defaultConfigFile := path.Join(homeFolder, defaultConfigDirName, defaultConfigFileName)
	if _, err := os.Stat(defaultConfigFile); err == nil {
		return defaultConfigFile
	}
	Debugf("The %s does not exist, will check env var %s for file path.", defaultConfigFile, configFilePathEnvVarName)
	// Read configuration file path from OCI_CONFIG_FILE env var
	fallbackConfigFile, existed := os.LookupEnv(configFilePathEnvVarName)
	if !existed {
		Debugf("The env var %s does not exist...", configFilePathEnvVarName)
		return defaultConfigFile
	}
	if _, err := os.Stat(fallbackConfigFile); os.IsNotExist(err) {
		Debugf("The specified cfg file path in the env var %s does not exist: %s", configFilePathEnvVarName, fallbackConfigFile)
		return defaultConfigFile
	}
	return fallbackConfigFile
}

// CustomProfileConfigProvider returns the config provider of given profile. The custom profile config provider
// will look for configurations in 2 places: file in $HOME/.oci/config,  and variables names starting with the
// string TF_VAR. If the same configuration is found in multiple places the provider will prefer the first one.
func CustomProfileConfigProvider(customConfigPath string, profile string) ConfigurationProvider {
	homeFolder := getHomeFolder()
	if customConfigPath == "" {
		customConfigPath = path.Join(homeFolder, defaultConfigDirName, defaultConfigFileName)
	}
	customFileProvider, _ := ConfigurationProviderFromFileWithProfile(customConfigPath, profile, "")
	defaultFileProvider, _ := ConfigurationProviderFromFileWithProfile(customConfigPath, "DEFAULT", "")
	environmentProvider := environmentConfigurationProvider{EnvironmentVariablePrefix: "TF_VAR"}
	provider, _ := ComposingConfigurationProvider([]ConfigurationProvider{customFileProvider, defaultFileProvider, environmentProvider})
	Debugf("Configuration provided by: %s", provider)
	return provider
}

func (client *BaseClient) prepareRequest(request *http.Request) (err error) {
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

func (client BaseClient) intercept(request *http.Request) (err error) {
	if client.Interceptor != nil {
		err = client.Interceptor(request)
	}
	return
}

func checkForSuccessfulResponse(res *http.Response) error {
	familyStatusCode := res.StatusCode / 100
	if familyStatusCode == 4 || familyStatusCode == 5 {
		return newServiceFailureFromResponse(res)
	}
	return nil
}

// OCIRequest is any request made to an OCI service.
type OCIRequest interface {
	// HTTPRequest assembles an HTTP request.
	HTTPRequest(method, path string) (http.Request, error)
}

// RequestMetadata is metadata about an OCIRequest. This structure represents the behavior exhibited by the SDK when
// issuing (or reissuing) a request.
type RequestMetadata struct {
	// RetryPolicy is the policy for reissuing the request. If no retry policy is set on the request,
	// then the request will be issued exactly once.
	RetryPolicy *RetryPolicy
}

// OCIResponse is the response from issuing a request to an OCI service.
type OCIResponse interface {
	// HTTPResponse returns the raw HTTP response.
	HTTPResponse() *http.Response
}

// OCIOperation is the generalization of a request-response cycle undergone by an OCI service.
type OCIOperation func(context.Context, OCIRequest) (OCIResponse, error)

//ClientCallDetails a set of settings used by the a single Call operation of the http Client
type ClientCallDetails struct {
	Signer HTTPRequestSigner
}

// Call executes the http request with the given context
func (client BaseClient) Call(ctx context.Context, request *http.Request) (response *http.Response, err error) {
	return client.CallWithDetails(ctx, request, ClientCallDetails{Signer: client.Signer})
}

// CallWithDetails executes the http request, the given context using details specified in the paremeters, this function
// provides a way to override some settings present in the client
func (client BaseClient) CallWithDetails(ctx context.Context, request *http.Request, details ClientCallDetails) (response *http.Response, err error) {
	Debugln("Atempting to call downstream service")
	request = request.WithContext(ctx)

	err = client.prepareRequest(request)
	if err != nil {
		return
	}

	//Intercept
	err = client.intercept(request)
	if err != nil {
		return
	}

	//Sign the request
	err = details.Signer.Sign(request)
	if err != nil {
		return
	}

	IfDebug(func() {
		dumpBody := true
		if request.ContentLength > maxBodyLenForDebug {
			Debugf("not dumping body too big\n")
			dumpBody = false
		}
		dumpBody = dumpBody && defaultLogger.LogLevel() == verboseLogging
		if dump, e := httputil.DumpRequestOut(request, dumpBody); e == nil {
			Debugf("Dump Request %s", string(dump))
		} else {
			Debugf("%v\n", e)
		}
	})

	//Execute the http request
	response, err = client.HTTPClient.Do(request)

	IfDebug(func() {
		if err != nil {
			Debugf("%v\n", err)
			return
		}

		dumpBody := true
		if response.ContentLength > maxBodyLenForDebug {
			Debugf("not dumping body too big\n")
			dumpBody = false
		}

		dumpBody = dumpBody && defaultLogger.LogLevel() == verboseLogging
		if dump, e := httputil.DumpResponse(response, dumpBody); e == nil {
			Debugf("Dump Response %s", string(dump))
		} else {
			Debugf("%v\n", e)
		}
	})

	if err != nil {
		return
	}

	err = checkForSuccessfulResponse(response)
	return
}

//CloseBodyIfValid closes the body of an http response if the response and the body are valid
func CloseBodyIfValid(httpResponse *http.Response) {
	if httpResponse != nil && httpResponse.Body != nil {
		httpResponse.Body.Close()
	}
}
