package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
)

type URLCredentialsProvider struct {
	url string
	// for sts
	sessionCredentials *SessionCredentials
	// for http options
	httpOptions *HttpOptions
	// inner
	expirationTimestamp int64
}

type URLCredentialsProviderBuilder struct {
	provider *URLCredentialsProvider
}

func NewURLCredentialsProviderBuilderBuilder() *URLCredentialsProviderBuilder {
	return &URLCredentialsProviderBuilder{
		provider: &URLCredentialsProvider{},
	}
}

func (builder *URLCredentialsProviderBuilder) WithUrl(url string) *URLCredentialsProviderBuilder {
	builder.provider.url = url
	return builder
}

func (builder *URLCredentialsProviderBuilder) WithHttpOptions(httpOptions *HttpOptions) *URLCredentialsProviderBuilder {
	builder.provider.httpOptions = httpOptions
	return builder
}

func (builder *URLCredentialsProviderBuilder) Build() (provider *URLCredentialsProvider, err error) {

	if builder.provider.url == "" {
		builder.provider.url = os.Getenv("ALIBABA_CLOUD_CREDENTIALS_URI")
	}

	if builder.provider.url == "" {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "The url is empty", nil)
		return
	}

	provider = builder.provider
	return
}

type urlResponse struct {
	AccessKeyId     *string `json:"AccessKeyId"`
	AccessKeySecret *string `json:"AccessKeySecret"`
	SecurityToken   *string `json:"SecurityToken"`
	Expiration      *string `json:"Expiration"`
}

func (provider *URLCredentialsProvider) getCredentials() (session *SessionCredentials, err error) {
	connectTimeout := 5 * time.Second
	readTimeout := 10 * time.Second
	if provider.httpOptions != nil && provider.httpOptions.ConnectTimeout > 0 {
		connectTimeout = provider.httpOptions.ConnectTimeout
	}
	if provider.httpOptions != nil && provider.httpOptions.ReadTimeout > 0 {
		readTimeout = provider.httpOptions.ReadTimeout
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		return (&net.Dialer{
			Timeout:   connectTimeout,
			DualStack: true,
		}).DialContext(ctx, network, address)
	}

	httpClient := &http.Client{
		Timeout:   connectTimeout + readTimeout,
		Transport: transport,
	}

	httpRequest, err := hookNewRequest(http.NewRequest)("GET", provider.url, strings.NewReader(""))
	if err != nil {
		return
	}

	httpResponse, err := hookDo(httpClient.Do)(httpRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer httpResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = fmt.Errorf("get credentials from %s failed: %s", provider.url, string(responseBody))
		return
	}

	var resp urlResponse
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		err = fmt.Errorf("get credentials from %s failed with error, json unmarshal fail: %s", provider.url, err.Error())
		return
	}

	if resp.AccessKeyId == nil || resp.AccessKeySecret == nil || resp.SecurityToken == nil || resp.Expiration == nil {
		err = fmt.Errorf("refresh credentials from %s failed: %s", provider.url, string(responseBody))
		return
	}

	session = &SessionCredentials{
		AccessKeyId:     *resp.AccessKeyId,
		AccessKeySecret: *resp.AccessKeySecret,
		SecurityToken:   *resp.SecurityToken,
		Expiration:      *resp.Expiration,
	}
	return
}

func (provider *URLCredentialsProvider) needUpdateCredential() (result bool) {
	if provider.expirationTimestamp == 0 {
		return true
	}

	return provider.expirationTimestamp-time.Now().Unix() <= 180
}

func (provider *URLCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
	if provider.sessionCredentials == nil || provider.needUpdateCredential() {
		sessionCredentials, err1 := provider.getCredentials()
		if err1 != nil {
			return nil, err1
		}

		provider.sessionCredentials = sessionCredentials
		expirationTime, err2 := time.Parse("2006-01-02T15:04:05Z", sessionCredentials.Expiration)
		if err2 != nil {
			return nil, err2
		}
		provider.expirationTimestamp = expirationTime.Unix()
	}

	cc = &Credentials{
		AccessKeyId:     provider.sessionCredentials.AccessKeyId,
		AccessKeySecret: provider.sessionCredentials.AccessKeySecret,
		SecurityToken:   provider.sessionCredentials.SecurityToken,
		ProviderName:    provider.GetProviderName(),
	}
	return
}

func (provider *URLCredentialsProvider) GetProviderName() string {
	return "credential_uri"
}
