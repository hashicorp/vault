package oauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

type HttpClientFactory func() *http.Client

// TokenResponse represents successful token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// ErrorResponse repsents a failed response
type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

// OauthClient represents a stateful Oauth client
type OauthClient struct {
	Service         string
	Client          *http.Client
	Headers         map[string]string
	ClientID        string
	ClientSecret    string
	SourceHeader    string
	ResponseHeaders http.Header
}

// OauthConfig represents configuration used to create Oauth clients
type OauthConfig struct {
	ClientID     string   `json:"clientID" structs:"clientID" mapstructure:"clientID"`
	ClientSecret string   `json:"clientSecret" structs:"clientSecret" mapstructure:"clientSecret"`
	ServiceURL   string   `json:"serviceUrl" structs:"serviceUrl" mapstructure:"serviceUrl"`
	AppID        string   `json:"appID" structs:"appID" mapstructure:"appID"`
	Scope        string   `json:"scope,omitempty" structs:"scope" mapstructure:"scope"`
	Policies     []string `json:"policies,omitempty" structs:"policies" mapstructure:"policies"`
}

// GetConfig reads config from a file specified as arg or from the environment
func GetConfig(configPath string) (*OauthConfig, error) {
	var config OauthConfig

	// No path specified; use environment
	if configPath == "" {
		path, found := os.LookupEnv("CENTRIFY_OAUTH_CONFIGPATH")
		if found {
			configPath = path
		}
	}

	if configPath != "" {
		fileBytes, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("error reading configuration from path '%s': %s", configPath, err)
		}

		err = json.Unmarshal(fileBytes, &config)
		if err != nil {
			return nil, fmt.Errorf("error demarshalling configuration from path '%s': %s", configPath, err)
		}
	}

	return &config, nil
}

// GetNewClient creates a new client for the specified endpoint
func GetNewClient(service string, httpFactory HttpClientFactory) (*OauthClient, error) {
	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, err
	}

	// Munge on the service a little bit, force it to have no trailing / and always start with https://
	url, err := url.Parse(service)
	if err != nil {
		return nil, err
	}
	url.Scheme = "https"
	url.Path = ""

	client := &OauthClient{}
	client.Service = url.String()
	if httpFactory != nil {
		client.Client = httpFactory()
	} else {
		client.Client = &http.Client{}
	}
	client.Client.Jar = jar
	client.Headers = make(map[string]string)
	client.SourceHeader = "cloud-golang-sdk"
	return client, err
}

// GetNewConfidentialClient creates a new client for the specified endpoint
func GetNewConfidentialClient(service string, clientID string, clientSecret string, httpFactory HttpClientFactory) (*OauthClient, error) {
	client, err := GetNewClient(service, httpFactory)
	if err != nil {
		return nil, err
	}

	client.ClientID = clientID
	client.ClientSecret = clientSecret
	return client, nil
}

// ResourceOwner implements the ResourceOwner flow
func (c *OauthClient) ResourceOwner(appID string, scope string, owner string, ownerPassword string) (*TokenResponse, *ErrorResponse, error) {
	args := make(map[string]string)
	args["grant_type"] = "password"
	args["username"] = owner
	args["password"] = ownerPassword
	args["scope"] = scope
	return c.postAndGetResponse("/oauth2/token/"+appID, args)
}

func (c *OauthClient) ClientCredentials(appID string, scope string) (*TokenResponse, *ErrorResponse, error) {
	args := make(map[string]string)
	args["grant_type"] = "client_credentials"
	args["scope"] = scope
	return c.postAndGetResponse("/oauth2/token/"+appID, args)
}

func (c *OauthClient) RefreshToken(appID string, refreshToken string) (*TokenResponse, *ErrorResponse, error) {
	args := make(map[string]string)
	args["grant_type"] = "refresh_token"
	args["refresh_token"] = refreshToken
	return c.postAndGetResponse("/oauth2/token/"+appID, args)
}

func (c *OauthClient) postAndGetResponse(method string, args map[string]string) (*TokenResponse, *ErrorResponse, error) {
	body, status, err := c.postAndGetBody(method, args)
	if err != nil {
		return nil, nil, err
	}

	if status == 200 {
		response, err := bodyToTokenResponse(body)
		if err != nil {
			return nil, nil, err
		}
		return response, nil, nil
	}

	response, err := bodyToErrorResponse(body)
	if err != nil {
		return nil, nil, err
	}
	return nil, response, nil
}

func (c *OauthClient) postAndGetBody(method string, args map[string]string) ([]byte, int, error) {
	postdata := strings.NewReader(payloadFromMap(args))
	postreq, err := http.NewRequest("POST", c.Service+method, postdata)

	if err != nil {
		return nil, 0, err
	}

	if c.ClientID != "" && c.ClientSecret != "" {
		postreq.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.ClientID+":"+c.ClientSecret)))
	}

	postreq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	postreq.Header.Add("X-CENTRIFY-NATIVE-CLIENT", "Yes")
	postreq.Header.Add("X-CFY-SRC", c.SourceHeader)

	for k, v := range c.Headers {
		postreq.Header.Add(k, v)
	}

	httpresp, err := c.Client.Do(postreq)
	if err != nil {
		c.ResponseHeaders = nil
		return nil, 0, err
	}

	defer httpresp.Body.Close()
	c.ResponseHeaders = httpresp.Header

	body, err := ioutil.ReadAll(httpresp.Body)
	if err != nil {
		return nil, httpresp.StatusCode, err
	}
	return body, httpresp.StatusCode, nil
}

// GetLastResponseHeaders returns the response header for the previous REST call
func (c *OauthClient) GetLastResponseHeaders() http.Header {
	return c.ResponseHeaders
}

func payloadFromMap(input map[string]string) string {
	data := url.Values{}
	for i, v := range input {
		data.Add(i, v)
	}
	return data.Encode()
}

func bodyToTokenResponse(body []byte) (*TokenResponse, error) {
	reply := &TokenResponse{}
	err := json.Unmarshal(body, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func bodyToErrorResponse(body []byte) (*ErrorResponse, error) {
	reply := &ErrorResponse{}
	err := json.Unmarshal(body, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
