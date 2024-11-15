package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
)

type assumedRoleUser struct {
}

type credentials struct {
	SecurityToken   *string `json:"SecurityToken"`
	Expiration      *string `json:"Expiration"`
	AccessKeySecret *string `json:"AccessKeySecret"`
	AccessKeyId     *string `json:"AccessKeyId"`
}

type ecsRAMRoleCredentials struct {
	SecurityToken   *string `json:"SecurityToken"`
	Expiration      *string `json:"Expiration"`
	AccessKeySecret *string `json:"AccessKeySecret"`
	AccessKeyId     *string `json:"AccessKeyId"`
	LastUpdated     *string `json:"LastUpdated"`
	Code            *string `json:"Code"`
}

type assumeRoleResponse struct {
	RequestID       *string          `json:"RequestId"`
	AssumedRoleUser *assumedRoleUser `json:"AssumedRoleUser"`
	Credentials     *credentials     `json:"Credentials"`
}

type generateSessionAccessKeyResponse struct {
	RequestID        *string           `json:"RequestId"`
	SessionAccessKey *sessionAccessKey `json:"SessionAccessKey"`
}

type sessionAccessKey struct {
	SessionAccessKeyId     *string `json:"SessionAccessKeyId"`
	SessionAccessKeySecret *string `json:"SessionAccessKeySecret"`
	Expiration             *string `json:"Expiration"`
}

type SessionCredentials struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      string
}

type Credentials struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	BearerToken     string
	ProviderName    string
}

type do func(req *http.Request) (*http.Response, error)

var hookDo = func(fn do) do {
	return fn
}

type newReuqest func(method, url string, body io.Reader) (*http.Request, error)

var hookNewRequest = func(fn newReuqest) newReuqest {
	return fn
}

type HttpOptions struct {
	// Connection timeout
	ConnectTimeout time.Duration
	// Read timeout
	ReadTimeout time.Duration
}

type CredentialsProvider interface {
	GetCredentials() (cc *Credentials, err error)
	GetProviderName() string
}

type StaticAKCredentialsProvider struct {
	accessKeyId     string
	accessKeySecret string
}

type StaticAKCredentialsProviderBuilder struct {
	provider *StaticAKCredentialsProvider
}

func NewStaticAKCredentialsProviderBuilder() *StaticAKCredentialsProviderBuilder {
	return &StaticAKCredentialsProviderBuilder{
		provider: &StaticAKCredentialsProvider{},
	}
}

func (builder *StaticAKCredentialsProviderBuilder) WithAccessKeyId(accessKeyId string) *StaticAKCredentialsProviderBuilder {
	builder.provider.accessKeyId = accessKeyId
	return builder
}

func (builder *StaticAKCredentialsProviderBuilder) WithAccessKeySecret(accessKeySecret string) *StaticAKCredentialsProviderBuilder {
	builder.provider.accessKeySecret = accessKeySecret
	return builder
}

func (builder *StaticAKCredentialsProviderBuilder) Build() (provider *StaticAKCredentialsProvider, err error) {
	if builder.provider.accessKeyId == "" {
		builder.provider.accessKeyId = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	}

	if builder.provider.accessKeyId == "" {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "The access key id is empty", nil)
		return
	}

	if builder.provider.accessKeySecret == "" {
		builder.provider.accessKeySecret = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	}

	if builder.provider.accessKeySecret == "" {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "The access key secret is empty", nil)
		return
	}

	provider = builder.provider
	return
}

func NewStaticAKCredentialsProvider(accessKeyId, accessKeySecret string) *StaticAKCredentialsProvider {
	return &StaticAKCredentialsProvider{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}
}

func (provider *StaticAKCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
	cc = &Credentials{
		AccessKeyId:     provider.accessKeyId,
		AccessKeySecret: provider.accessKeySecret,
		ProviderName:    provider.GetProviderName(),
	}
	return
}

func (provider *StaticAKCredentialsProvider) GetProviderName() string {
	return "static_ak"
}

type StaticSTSCredentialsProvider struct {
	accessKeyId     string
	accessKeySecret string
	securityToken   string
}

type StaticSTSCredentialsProviderBuilder struct {
	provider *StaticSTSCredentialsProvider
}

func NewStaticSTSCredentialsProviderBuilder() *StaticSTSCredentialsProviderBuilder {
	return &StaticSTSCredentialsProviderBuilder{
		provider: &StaticSTSCredentialsProvider{},
	}
}

func (builder *StaticSTSCredentialsProviderBuilder) WithAccessKeyId(accessKeyId string) *StaticSTSCredentialsProviderBuilder {
	builder.provider.accessKeyId = accessKeyId
	return builder
}

func (builder *StaticSTSCredentialsProviderBuilder) WithAccessKeySecret(accessKeySecret string) *StaticSTSCredentialsProviderBuilder {
	builder.provider.accessKeySecret = accessKeySecret
	return builder
}

func (builder *StaticSTSCredentialsProviderBuilder) WithSecurityToken(securityToken string) *StaticSTSCredentialsProviderBuilder {
	builder.provider.securityToken = securityToken
	return builder
}

func (builder *StaticSTSCredentialsProviderBuilder) Build() (provider *StaticSTSCredentialsProvider, err error) {
	if builder.provider.accessKeyId == "" {
		builder.provider.accessKeyId = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")
	}

	if builder.provider.accessKeyId == "" {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "The access key id is empty", nil)
		return
	}

	if builder.provider.accessKeySecret == "" {
		builder.provider.accessKeySecret = os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
	}

	if builder.provider.accessKeySecret == "" {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "The access key secret is empty", nil)
		return
	}

	if builder.provider.securityToken == "" {
		builder.provider.securityToken = os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN")
	}

	if builder.provider.securityToken == "" {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "The security token is empty", nil)
		return
	}

	provider = builder.provider
	return
}

func NewStaticSTSCredentialsProvider(accessKeyId, accessKeySecret, securityToken string) *StaticSTSCredentialsProvider {
	return &StaticSTSCredentialsProvider{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		securityToken:   securityToken,
	}
}

func (provider *StaticSTSCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
	cc = &Credentials{
		AccessKeyId:     provider.accessKeyId,
		AccessKeySecret: provider.accessKeySecret,
		SecurityToken:   provider.securityToken,
		ProviderName:    provider.GetProviderName(),
	}
	return
}

func (provider *StaticSTSCredentialsProvider) GetProviderName() string {
	return "static_sts"
}

type BearerTokenCredentialsProvider struct {
	bearerToken string
}

func NewBearerTokenCredentialsProvider(bearerToken string) *BearerTokenCredentialsProvider {
	return &BearerTokenCredentialsProvider{
		bearerToken: bearerToken,
	}
}

func (provider *BearerTokenCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
	cc = &Credentials{
		BearerToken:  provider.bearerToken,
		ProviderName: provider.GetProviderName(),
	}
	return
}

func (provider *BearerTokenCredentialsProvider) GetProviderName() string {
	return "bearer_token"
}

// Deprecated: the RSA key pair credentials is deprecated
type RSAKeyPairCredentialsProvider struct {
	PublicKeyId         string
	PrivateKeyId        string
	durationSeconds     int
	sessionAccessKey    *sessionAccessKey
	lastUpdateTimestamp int64
	expirationTimestamp int64
}

// Deprecated: the RSA key pair credentials is deprecated
func NewRSAKeyPairCredentialsProvider(publicKeyId, privateKeyId string, durationSeconds int) (provider *RSAKeyPairCredentialsProvider, err error) {
	provider = &RSAKeyPairCredentialsProvider{
		PublicKeyId:  publicKeyId,
		PrivateKeyId: privateKeyId,
	}

	if durationSeconds > 0 {
		if durationSeconds >= 900 && durationSeconds <= 3600 {
			provider.durationSeconds = durationSeconds
		} else {
			err = errors.NewClientError(errors.InvalidParamErrorCode, "Key Pair session duration should be in the range of 15min - 1hr", nil)
		}
	} else {
		// set to default value
		provider.durationSeconds = 3600
	}
	return
}

// Deprecated: the RSA key pair credentials is deprecated
func (provider *RSAKeyPairCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
	if provider.sessionAccessKey == nil || provider.needUpdateCredential() {
		sessionAccessKey, err := provider.getCredentials()
		if err != nil {
			return nil, err
		}

		expirationTime, err := time.Parse("2006-01-02T15:04:05Z", *sessionAccessKey.Expiration)
		if err != nil {
			return nil, err
		}

		provider.sessionAccessKey = sessionAccessKey
		provider.lastUpdateTimestamp = time.Now().Unix()
		provider.expirationTimestamp = expirationTime.Unix()
	}

	cc = &Credentials{
		AccessKeyId:     *provider.sessionAccessKey.SessionAccessKeyId,
		AccessKeySecret: *provider.sessionAccessKey.SessionAccessKeySecret,
		ProviderName:    provider.GetProviderName(),
	}
	return
}

func (provider *RSAKeyPairCredentialsProvider) needUpdateCredential() bool {
	if provider.expirationTimestamp == 0 {
		return true
	}

	return provider.expirationTimestamp-time.Now().Unix() <= 180
}

func (provider *RSAKeyPairCredentialsProvider) getCredentials() (sessionAK *sessionAccessKey, err error) {
	method := "POST"
	host := "sts.ap-northeast-1.aliyuncs.com"

	queries := make(map[string]string)
	queries["Version"] = "2015-04-01"
	queries["Action"] = "GenerateSessionAccessKey"
	queries["Format"] = "JSON"
	queries["Timestamp"] = utils.GetTimeInFormatISO8601()
	queries["SignatureMethod"] = "SHA256withRSA"
	queries["SignatureVersion"] = "1.0"
	queries["SignatureNonce"] = utils.GetNonce()
	queries["PublicKeyId"] = provider.PublicKeyId
	queries["SignatureType"] = "PRIVATEKEY"

	bodyForm := make(map[string]string)
	bodyForm["DurationSeconds"] = strconv.Itoa(provider.durationSeconds)

	// caculate signature
	signParams := make(map[string]string)
	for key, value := range queries {
		signParams[key] = value
	}
	for key, value := range bodyForm {
		signParams[key] = value
	}

	stringToSign := utils.GetUrlFormedMap(signParams)
	stringToSign = strings.Replace(stringToSign, "+", "%20", -1)
	stringToSign = strings.Replace(stringToSign, "*", "%2A", -1)
	stringToSign = strings.Replace(stringToSign, "%7E", "~", -1)
	stringToSign = url.QueryEscape(stringToSign)
	stringToSign = method + "&%2F&" + stringToSign

	queries["Signature"] = utils.Sha256WithRsa(stringToSign, provider.PrivateKeyId)

	querystring := utils.GetUrlFormedMap(queries)
	// do request
	httpUrl := fmt.Sprintf("https://%s/?%s", host, querystring)

	body := utils.GetUrlFormedMap(bodyForm)

	httpRequest, err := hookNewRequest(http.NewRequest)(method, httpUrl, strings.NewReader(body))
	if err != nil {
		return
	}

	// set headers
	httpRequest.Header["Accept-Encoding"] = []string{"identity"}
	httpRequest.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	httpClient := &http.Client{}

	httpResponse, err := hookDo(httpClient.Do)(httpRequest)
	if err != nil {
		return
	}

	defer httpResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		message := "refresh temp ak failed"
		err = errors.NewServerError(httpResponse.StatusCode, string(responseBody), message)
		return
	}

	var data generateSessionAccessKeyResponse
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		err = fmt.Errorf("refresh temp ak err, json.Unmarshal fail: %s", err.Error())
		return
	}

	if data.SessionAccessKey == nil {
		err = fmt.Errorf("refresh temp ak token err, fail to get credentials")
		return
	}

	if data.SessionAccessKey.SessionAccessKeyId == nil || data.SessionAccessKey.SessionAccessKeySecret == nil {
		err = fmt.Errorf("refresh temp ak token err, fail to get credentials")
		return
	}

	sessionAK = data.SessionAccessKey
	return
}

func (provider *RSAKeyPairCredentialsProvider) GetProviderName() string {
	return "rsa_key_pair"
}

type RAMRoleARNCredentialsProvider struct {
	// for previous credentials
	accessKeyId         string
	accessKeySecret     string
	securityToken       string
	credentialsProvider CredentialsProvider

	roleArn         string
	roleSessionName string
	durationSeconds int
	policy          string
	externalId      string

	// for sts endpoint
	stsRegion   string
	enableVpc   bool
	stsEndpoint string

	// for http options
	httpOptions *HttpOptions

	// inner
	expirationTimestamp  int64
	lastUpdateTimestamp  int64
	previousProviderName string
	sessionCredentials   *SessionCredentials
}

type RAMRoleARNCredentialsProviderBuilder struct {
	provider *RAMRoleARNCredentialsProvider
}

func NewRAMRoleARNCredentialsProviderBuilder() *RAMRoleARNCredentialsProviderBuilder {
	return &RAMRoleARNCredentialsProviderBuilder{
		provider: &RAMRoleARNCredentialsProvider{},
	}
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithAccessKeyId(accessKeyId string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.accessKeyId = accessKeyId
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithAccessKeySecret(accessKeySecret string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.accessKeySecret = accessKeySecret
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithSecurityToken(securityToken string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.securityToken = securityToken
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithCredentialsProvider(credentialsProvider CredentialsProvider) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.credentialsProvider = credentialsProvider
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithRoleArn(roleArn string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.roleArn = roleArn
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithStsRegion(region string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.stsRegion = region
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithEnableVpc(enableVpc bool) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.enableVpc = enableVpc
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithStsEndpoint(endpoint string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.stsEndpoint = endpoint
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithRoleSessionName(roleSessionName string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.roleSessionName = roleSessionName
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithPolicy(policy string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.policy = policy
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithExternalId(externalId string) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.externalId = externalId
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithDurationSeconds(durationSeconds int) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.durationSeconds = durationSeconds
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) WithHttpOptions(httpOptions *HttpOptions) *RAMRoleARNCredentialsProviderBuilder {
	builder.provider.httpOptions = httpOptions
	return builder
}

func (builder *RAMRoleARNCredentialsProviderBuilder) Build() (provider *RAMRoleARNCredentialsProvider, err error) {
	if builder.provider.credentialsProvider == nil {
		if builder.provider.accessKeyId != "" && builder.provider.accessKeySecret != "" && builder.provider.securityToken != "" {
			builder.provider.credentialsProvider, err = NewStaticSTSCredentialsProviderBuilder().
				WithAccessKeyId(builder.provider.accessKeyId).
				WithAccessKeySecret(builder.provider.accessKeySecret).
				WithSecurityToken(builder.provider.securityToken).
				Build()
			if err != nil {
				return
			}
		} else if builder.provider.accessKeyId != "" && builder.provider.accessKeySecret != "" {
			builder.provider.credentialsProvider, err = NewStaticAKCredentialsProviderBuilder().
				WithAccessKeyId(builder.provider.accessKeyId).
				WithAccessKeySecret(builder.provider.accessKeySecret).
				Build()
			if err != nil {
				return
			}
		} else {
			err = errors.NewClientError(errors.InvalidParamErrorCode, "Must specify a previous credentials provider to assume role", nil)
			return
		}
	}

	if builder.provider.roleArn == "" {
		if roleArn := os.Getenv("ALIBABA_CLOUD_ROLE_ARN"); roleArn != "" {
			builder.provider.roleArn = roleArn
		} else {
			err = errors.NewClientError(errors.InvalidParamErrorCode, "The RoleArn is empty", nil)
			return
		}
	}

	if builder.provider.roleSessionName == "" {
		if roleSessionName := os.Getenv("ALIBABA_CLOUD_ROLE_SESSION_NAME"); roleSessionName != "" {
			builder.provider.roleSessionName = roleSessionName
		} else {
			builder.provider.roleSessionName = "aliyun-go-sdk-" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
		}
	}

	// duration seconds
	if builder.provider.durationSeconds == 0 {
		// default to 3600
		builder.provider.durationSeconds = 3600
	}

	if builder.provider.durationSeconds < 900 {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "Session duration should be in the range of 900s - max session duration", nil)
		return
	}

	// sts endpoint
	if builder.provider.stsEndpoint == "" {
		if !builder.provider.enableVpc {
			builder.provider.enableVpc = strings.ToLower(os.Getenv("ALIBABA_CLOUD_VPC_ENDPOINT_ENABLED")) == "true"
		}
		prefix := "sts"
		if builder.provider.enableVpc {
			prefix = "sts-vpc"
		}
		if builder.provider.stsRegion != "" {
			builder.provider.stsEndpoint = fmt.Sprintf("%s.%s.aliyuncs.com", prefix, builder.provider.stsRegion)
		} else if region := os.Getenv("ALIBABA_CLOUD_STS_REGION"); region != "" {
			builder.provider.stsEndpoint = fmt.Sprintf("%s.%s.aliyuncs.com", prefix, region)
		} else {
			builder.provider.stsEndpoint = "sts.aliyuncs.com"
		}
	}

	provider = builder.provider
	return
}

func NewRAMRoleARNCredentialsProvider(credentialsProvider CredentialsProvider, roleArn, roleSessionName string, durationSeconds int, policy, stsRegion, externalId string) (provider *RAMRoleARNCredentialsProvider, err error) {
	provider = &RAMRoleARNCredentialsProvider{
		credentialsProvider: credentialsProvider,
		roleArn:             roleArn,
		durationSeconds:     durationSeconds,
		policy:              policy,
		stsRegion:           stsRegion,
		externalId:          externalId,
	}

	if len(roleSessionName) > 0 {
		provider.roleSessionName = roleSessionName
	} else {
		provider.roleSessionName = "aliyun-go-sdk-" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
	}

	if durationSeconds > 0 {
		if durationSeconds >= 900 && durationSeconds <= 3600 {
			provider.durationSeconds = durationSeconds
		} else {
			err = errors.NewClientError(errors.InvalidParamErrorCode, "Assume Role session duration should be in the range of 15min - 1hr", nil)
		}
	} else {
		// default to 3600
		provider.durationSeconds = 3600
	}

	return
}

func (provider *RAMRoleARNCredentialsProvider) getCredentials(cc *Credentials) (sessionCredentials *SessionCredentials, err error) {
	method := "POST"
	var host string
	if provider.stsEndpoint != "" {
		host = provider.stsEndpoint
	} else if provider.stsRegion != "" {
		host = fmt.Sprintf("sts.%s.aliyuncs.com", provider.stsRegion)
	} else {
		host = "sts.aliyuncs.com"
	}

	queries := make(map[string]string)
	queries["Version"] = "2015-04-01"
	queries["Action"] = "AssumeRole"
	queries["Format"] = "JSON"
	queries["Timestamp"] = utils.GetTimeInFormatISO8601()
	queries["SignatureMethod"] = "HMAC-SHA1"
	queries["SignatureVersion"] = "1.0"
	queries["SignatureNonce"] = utils.GetNonce()
	queries["AccessKeyId"] = cc.AccessKeyId
	if cc.SecurityToken != "" {
		queries["SecurityToken"] = cc.SecurityToken
	}

	bodyForm := make(map[string]string)
	bodyForm["RoleArn"] = provider.roleArn
	if provider.policy != "" {
		bodyForm["Policy"] = provider.policy
	}
	if provider.externalId != "" {
		bodyForm["ExternalId"] = provider.externalId
	}
	bodyForm["RoleSessionName"] = provider.roleSessionName
	bodyForm["DurationSeconds"] = strconv.Itoa(provider.durationSeconds)

	// caculate signature
	signParams := make(map[string]string)
	for key, value := range queries {
		signParams[key] = value
	}
	for key, value := range bodyForm {
		signParams[key] = value
	}

	stringToSign := utils.GetUrlFormedMap(signParams)
	stringToSign = strings.Replace(stringToSign, "+", "%20", -1)
	stringToSign = strings.Replace(stringToSign, "*", "%2A", -1)
	stringToSign = strings.Replace(stringToSign, "%7E", "~", -1)
	stringToSign = url.QueryEscape(stringToSign)
	stringToSign = method + "&%2F&" + stringToSign
	secret := cc.AccessKeySecret + "&"
	queries["Signature"] = utils.ShaHmac1(stringToSign, secret)

	querystring := utils.GetUrlFormedMap(queries)
	// do request
	httpUrl := fmt.Sprintf("https://%s/?%s", host, querystring)

	body := utils.GetUrlFormedMap(bodyForm)

	httpRequest, err := hookNewRequest(http.NewRequest)(method, httpUrl, strings.NewReader(body))
	if err != nil {
		return
	}

	// set headers
	httpRequest.Header["Accept-Encoding"] = []string{"identity"}
	httpRequest.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}

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

	httpResponse, err := hookDo(httpClient.Do)(httpRequest)
	if err != nil {
		return
	}

	defer httpResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		message := "refresh session token failed"
		err = errors.NewServerError(httpResponse.StatusCode, string(responseBody), message)
		return
	}
	var data assumeRoleResponse
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		err = fmt.Errorf("refresh RoleArn sts token err, json.Unmarshal fail: %s", err.Error())
		return
	}
	if data.Credentials == nil {
		err = fmt.Errorf("refresh RoleArn sts token err, fail to get credentials")
		return
	}

	if data.Credentials.AccessKeyId == nil || data.Credentials.AccessKeySecret == nil || data.Credentials.SecurityToken == nil {
		err = fmt.Errorf("refresh RoleArn sts token err, fail to get credentials")
		return
	}

	sessionCredentials = &SessionCredentials{
		AccessKeyId:     *data.Credentials.AccessKeyId,
		AccessKeySecret: *data.Credentials.AccessKeySecret,
		SecurityToken:   *data.Credentials.SecurityToken,
		Expiration:      *data.Credentials.Expiration,
	}
	return
}

func (provider *RAMRoleARNCredentialsProvider) needUpdateCredential() (result bool) {
	if provider.expirationTimestamp == 0 {
		return true
	}

	return provider.expirationTimestamp-time.Now().Unix() <= 180
}

func (provider *RAMRoleARNCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
	if provider.sessionCredentials == nil || provider.needUpdateCredential() {
		// 获取前置凭证
		previousCredentials, err1 := provider.credentialsProvider.GetCredentials()
		if err1 != nil {
			return nil, err1
		}
		sessionCredentials, err2 := provider.getCredentials(previousCredentials)
		if err2 != nil {
			return nil, err2
		}

		expirationTime, err := time.Parse("2006-01-02T15:04:05Z", sessionCredentials.Expiration)
		if err != nil {
			return nil, err
		}

		provider.expirationTimestamp = expirationTime.Unix()
		provider.lastUpdateTimestamp = time.Now().Unix()
		provider.previousProviderName = previousCredentials.ProviderName
		provider.sessionCredentials = sessionCredentials
	}
	if provider.previousProviderName == "" {
		provider.previousProviderName = provider.credentialsProvider.GetProviderName()
	}

	cc = &Credentials{
		AccessKeyId:     provider.sessionCredentials.AccessKeyId,
		AccessKeySecret: provider.sessionCredentials.AccessKeySecret,
		SecurityToken:   provider.sessionCredentials.SecurityToken,
		ProviderName:    fmt.Sprintf("%s/%s", provider.GetProviderName(), provider.previousProviderName),
	}
	return
}

func (provider *RAMRoleARNCredentialsProvider) GetProviderName() string {
	return "ram_role_arn"
}

type ECSRAMRoleCredentialsProvider struct {
	roleName      string
	disableIMDSv1 bool

	// for http options
	httpOptions *HttpOptions

	sessionCredentials  *SessionCredentials
	expirationTimestamp int64
}

type ECSRAMRoleCredentialsProviderBuilder struct {
	provider *ECSRAMRoleCredentialsProvider
}

func NewECSRAMRoleCredentialsProviderBuilder() *ECSRAMRoleCredentialsProviderBuilder {
	return &ECSRAMRoleCredentialsProviderBuilder{
		provider: &ECSRAMRoleCredentialsProvider{},
	}
}

func (builder *ECSRAMRoleCredentialsProviderBuilder) WithRoleName(roleName string) *ECSRAMRoleCredentialsProviderBuilder {
	builder.provider.roleName = roleName
	return builder
}

func (builder *ECSRAMRoleCredentialsProviderBuilder) WithDisableIMDSv1(disableIMDSv1 bool) *ECSRAMRoleCredentialsProviderBuilder {
	builder.provider.disableIMDSv1 = disableIMDSv1
	return builder
}

func (builder *ECSRAMRoleCredentialsProviderBuilder) WithHttpOptions(httpOptions *HttpOptions) *ECSRAMRoleCredentialsProviderBuilder {
	builder.provider.httpOptions = httpOptions
	return builder
}

const defaultMetadataTokenDuration = 21600 // 6 hours

func (builder *ECSRAMRoleCredentialsProviderBuilder) Build() (provider *ECSRAMRoleCredentialsProvider, err error) {

	if strings.ToLower(os.Getenv("ALIBABA_CLOUD_ECS_METADATA_DISABLED")) == "true" {
		err = fmt.Errorf("IMDS credentials is disabled")
		return
	}

	// 设置 roleName 默认值
	if builder.provider.roleName == "" {
		builder.provider.roleName = os.Getenv("ALIBABA_CLOUD_ECS_METADATA")
	}

	if !builder.provider.disableIMDSv1 {
		builder.provider.disableIMDSv1 = strings.ToLower(os.Getenv("ALIBABA_CLOUD_IMDSV1_DISABLED")) == "true"
	}

	provider = builder.provider
	return
}

func NewECSRAMRoleCredentialsProvider(roleName string) *ECSRAMRoleCredentialsProvider {
	return &ECSRAMRoleCredentialsProvider{
		roleName: roleName,
	}
}

func (provider *ECSRAMRoleCredentialsProvider) needUpdateCredential() bool {
	if provider.expirationTimestamp == 0 {
		return true
	}

	return provider.expirationTimestamp-time.Now().Unix() <= 180
}

func (provider *ECSRAMRoleCredentialsProvider) getMetadataToken() (metadataToken string, err error) {
	// PUT http://100.100.100.200/latest/api/token
	var requestUrl = "http://100.100.100.200/latest/api/token"
	httpRequest, _err := hookNewRequest(http.NewRequest)("PUT", requestUrl, strings.NewReader(""))
	if _err != nil {
		if provider.disableIMDSv1 {
			err = fmt.Errorf("get metadata token failed: %s", _err.Error())
		}
		return
	}
	httpRequest.Header.Set("X-aliyun-ecs-metadata-token-ttl-seconds", strconv.Itoa(defaultMetadataTokenDuration))

	connectTimeout := 1 * time.Second
	readTimeout := 1 * time.Second

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

	httpResponse, _err := hookDo(httpClient.Do)(httpRequest)
	if _err != nil {
		if provider.disableIMDSv1 {
			err = fmt.Errorf("get metadata token failed: %s", _err.Error())
		}
		return
	}

	defer httpResponse.Body.Close()

	responseBody, _err := ioutil.ReadAll(httpResponse.Body)
	if _err != nil {
		if provider.disableIMDSv1 {
			err = fmt.Errorf("get metadata token failed: %s", _err.Error())
		}
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		if provider.disableIMDSv1 {
			err = errors.NewServerError(httpResponse.StatusCode, string(responseBody), "refresh Ecs sts token err")
		}
		return
	}

	metadataToken = strings.TrimSpace(string(responseBody))
	return
}

func (provider *ECSRAMRoleCredentialsProvider) getRoleName() (roleName string, err error) {
	connectTimeout := 1 * time.Second
	readTimeout := 1 * time.Second

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

	var securityCredURL = "http://100.100.100.200/latest/meta-data/ram/security-credentials/"
	httpRequest, err := hookNewRequest(http.NewRequest)("GET", securityCredURL, strings.NewReader(""))
	if err != nil {
		err = fmt.Errorf("get role name failed: %s", err.Error())
		return
	}

	metadataToken, err := provider.getMetadataToken()
	if err != nil {
		return
	}
	if metadataToken != "" {
		httpRequest.Header.Set("X-aliyun-ecs-metadata-token", metadataToken)
	}

	httpResponse, err := hookDo(httpClient.Do)(httpRequest)
	if err != nil {
		err = fmt.Errorf("get role name failed: %s", err.Error())
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = fmt.Errorf("get role name failed: request %s %d", securityCredURL, httpResponse.StatusCode)
		return
	}

	defer httpResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	roleName = strings.TrimSpace(string(responseBody))
	return
}

func (provider *ECSRAMRoleCredentialsProvider) getCredentials() (sessionCredentials *SessionCredentials, err error) {
	roleName := provider.roleName
	if roleName == "" {
		roleName, err = provider.getRoleName()
		if err != nil {
			return
		}
	}

	connectTimeout := 1 * time.Second
	readTimeout := 1 * time.Second

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

	var requestUrl = "http://100.100.100.200/latest/meta-data/ram/security-credentials/" + roleName
	httpRequest, err := hookNewRequest(http.NewRequest)("GET", requestUrl, strings.NewReader(""))
	if err != nil {
		err = fmt.Errorf("refresh Ecs sts token err: %s", err.Error())
		return
	}

	metadataToken, err := provider.getMetadataToken()
	if err != nil {
		return
	}
	if metadataToken != "" {
		httpRequest.Header.Set("X-aliyun-ecs-metadata-token", metadataToken)
	}

	httpResponse, err := hookDo(httpClient.Do)(httpRequest)
	if err != nil {
		err = fmt.Errorf("refresh Ecs sts token err: %s", err.Error())
		return
	}

	defer httpResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		err = fmt.Errorf("refresh Ecs sts token err, httpStatus: %d, message = %s", httpResponse.StatusCode, string(responseBody))
		return
	}

	var data ecsRAMRoleCredentials
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		err = fmt.Errorf("refresh Ecs sts token err, json.Unmarshal fail: %s", err.Error())
		return
	}

	if data.AccessKeyId == nil || data.AccessKeySecret == nil || data.SecurityToken == nil {
		err = fmt.Errorf("refresh Ecs sts token err, fail to get credentials")
		return
	}

	if *data.Code != "Success" {
		err = fmt.Errorf("refresh Ecs sts token err, Code is not Success")
		return
	}

	sessionCredentials = &SessionCredentials{
		AccessKeyId:     *data.AccessKeyId,
		AccessKeySecret: *data.AccessKeySecret,
		SecurityToken:   *data.SecurityToken,
		Expiration:      *data.Expiration,
	}
	return
}

func (provider *ECSRAMRoleCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
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

func (provider *ECSRAMRoleCredentialsProvider) GetProviderName() string {
	return "ecs_ram_role"
}

type OIDCCredentialsProvider struct {
	oidcProviderARN   string
	oidcTokenFilePath string
	roleArn           string
	roleSessionName   string
	durationSeconds   int
	policy            string
	// for sts endpoint
	stsRegion           string
	enableVpc           bool
	stsEndpoint         string
	lastUpdateTimestamp int64
	expirationTimestamp int64
	sessionCredentials  *SessionCredentials
	// for http options
	httpOptions *HttpOptions
}

type OIDCCredentialsProviderBuilder struct {
	provider *OIDCCredentialsProvider
}

func NewOIDCCredentialsProviderBuilder() *OIDCCredentialsProviderBuilder {
	return &OIDCCredentialsProviderBuilder{
		provider: &OIDCCredentialsProvider{},
	}
}

func (b *OIDCCredentialsProviderBuilder) WithOIDCProviderARN(oidcProviderArn string) *OIDCCredentialsProviderBuilder {
	b.provider.oidcProviderARN = oidcProviderArn
	return b
}

func (b *OIDCCredentialsProviderBuilder) WithOIDCTokenFilePath(oidcTokenFilePath string) *OIDCCredentialsProviderBuilder {
	b.provider.oidcTokenFilePath = oidcTokenFilePath
	return b
}

func (b *OIDCCredentialsProviderBuilder) WithRoleArn(roleArn string) *OIDCCredentialsProviderBuilder {
	b.provider.roleArn = roleArn
	return b
}

func (b *OIDCCredentialsProviderBuilder) WithRoleSessionName(roleSessionName string) *OIDCCredentialsProviderBuilder {
	b.provider.roleSessionName = roleSessionName
	return b
}

func (b *OIDCCredentialsProviderBuilder) WithDurationSeconds(durationSeconds int) *OIDCCredentialsProviderBuilder {
	b.provider.durationSeconds = durationSeconds
	return b
}

func (b *OIDCCredentialsProviderBuilder) WithStsRegion(region string) *OIDCCredentialsProviderBuilder {
	b.provider.stsRegion = region
	return b
}

func (b *OIDCCredentialsProviderBuilder) WithEnableVpc(enableVpc bool) *OIDCCredentialsProviderBuilder {
	b.provider.enableVpc = enableVpc
	return b
}

func (b *OIDCCredentialsProviderBuilder) WithSTSEndpoint(stsEndpoint string) *OIDCCredentialsProviderBuilder {
	b.provider.stsEndpoint = stsEndpoint
	return b
}

func (b *OIDCCredentialsProviderBuilder) WithPolicy(policy string) *OIDCCredentialsProviderBuilder {
	b.provider.policy = policy
	return b
}

func (b *OIDCCredentialsProviderBuilder) WithHttpOptions(httpOptions *HttpOptions) *OIDCCredentialsProviderBuilder {
	b.provider.httpOptions = httpOptions
	return b
}

func (b *OIDCCredentialsProviderBuilder) Build() (provider *OIDCCredentialsProvider, err error) {
	provider = b.provider

	if provider.roleSessionName == "" {
		provider.roleSessionName = "aliyun-go-sdk-" + strconv.FormatInt(time.Now().UnixNano()/1000, 10)
	}

	if provider.oidcTokenFilePath == "" {
		provider.oidcTokenFilePath = os.Getenv("ALIBABA_CLOUD_OIDC_TOKEN_FILE")
	}

	if provider.oidcTokenFilePath == "" {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "OIDCTokenFilePath can not be empty", nil)
		return
	}

	if provider.oidcProviderARN == "" {
		provider.oidcProviderARN = os.Getenv("ALIBABA_CLOUD_OIDC_PROVIDER_ARN")
	}

	if provider.oidcProviderARN == "" {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "OIDCProviderARN can not be empty", nil)
		return
	}

	if provider.roleArn == "" {
		provider.roleArn = os.Getenv("ALIBABA_CLOUD_ROLE_ARN")
	}

	if provider.roleArn == "" {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "RoleArn can not be empty", nil)
		return
	}

	if provider.durationSeconds == 0 {
		provider.durationSeconds = 3600
	}

	if provider.durationSeconds < 900 || provider.durationSeconds > 3600 {
		err = errors.NewClientError(errors.InvalidParamErrorCode, "Assume Role session duration should be in the range of 15min - 1hr", nil)
	}

	// sts endpoint
	if provider.stsEndpoint == "" {
		if !provider.enableVpc {
			provider.enableVpc = strings.ToLower(os.Getenv("ALIBABA_CLOUD_VPC_ENDPOINT_ENABLED")) == "true"
		}
		prefix := "sts"
		if provider.enableVpc {
			prefix = "sts-vpc"
		}
		if provider.stsRegion != "" {
			provider.stsEndpoint = fmt.Sprintf("%s.%s.aliyuncs.com", prefix, provider.stsRegion)
		} else if region := os.Getenv("ALIBABA_CLOUD_STS_REGION"); region != "" {
			provider.stsEndpoint = fmt.Sprintf("%s.%s.aliyuncs.com", prefix, region)
		} else {
			provider.stsEndpoint = "sts.aliyuncs.com"
		}
	}

	return
}

func (provider *OIDCCredentialsProvider) getCredentials() (sessionCredentials *SessionCredentials, err error) {
	method := "POST"
	var host string
	if provider.stsEndpoint != "" {
		host = provider.stsEndpoint
	} else if provider.stsRegion != "" {
		host = fmt.Sprintf("sts.%s.aliyuncs.com", provider.stsRegion)
	} else {
		host = "sts.aliyuncs.com"
	}

	queries := make(map[string]string)
	queries["Version"] = "2015-04-01"
	queries["Action"] = "AssumeRoleWithOIDC"
	queries["Format"] = "JSON"
	queries["Timestamp"] = utils.GetTimeInFormatISO8601()

	bodyForm := make(map[string]string)
	bodyForm["RoleArn"] = provider.roleArn
	bodyForm["OIDCProviderArn"] = provider.oidcProviderARN
	token, err := ioutil.ReadFile(provider.oidcTokenFilePath)
	if err != nil {
		return
	}

	bodyForm["OIDCToken"] = string(token)
	if provider.policy != "" {
		bodyForm["Policy"] = provider.policy
	}

	bodyForm["RoleSessionName"] = provider.roleSessionName
	bodyForm["DurationSeconds"] = strconv.Itoa(provider.durationSeconds)

	// caculate signature
	signParams := make(map[string]string)
	for key, value := range queries {
		signParams[key] = value
	}
	for key, value := range bodyForm {
		signParams[key] = value
	}

	querystring := utils.GetUrlFormedMap(queries)
	// do request
	httpUrl := fmt.Sprintf("https://%s/?%s", host, querystring)

	body := utils.GetUrlFormedMap(bodyForm)

	httpRequest, err := hookNewRequest(http.NewRequest)(method, httpUrl, strings.NewReader(body))
	if err != nil {
		return
	}

	// set headers
	httpRequest.Header["Accept-Encoding"] = []string{"identity"}
	httpRequest.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}

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

	httpResponse, err := hookDo(httpClient.Do)(httpRequest)
	if err != nil {
		return
	}

	defer httpResponse.Body.Close()

	responseBody, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return
	}

	if httpResponse.StatusCode != http.StatusOK {
		message := "get session token failed"
		err = errors.NewServerError(httpResponse.StatusCode, string(responseBody), message)
		return
	}
	var data assumeRoleResponse
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		err = fmt.Errorf("get oidc sts token err, json.Unmarshal fail: %s", err.Error())
		return
	}
	if data.Credentials == nil {
		err = fmt.Errorf("get oidc sts token err, fail to get credentials")
		return
	}

	if data.Credentials.AccessKeyId == nil || data.Credentials.AccessKeySecret == nil || data.Credentials.SecurityToken == nil {
		err = fmt.Errorf("refresh RoleArn sts token err, fail to get credentials")
		return
	}

	sessionCredentials = &SessionCredentials{
		AccessKeyId:     *data.Credentials.AccessKeyId,
		AccessKeySecret: *data.Credentials.AccessKeySecret,
		SecurityToken:   *data.Credentials.SecurityToken,
		Expiration:      *data.Credentials.Expiration,
	}
	return
}

func (provider *OIDCCredentialsProvider) needUpdateCredential() (result bool) {
	if provider.expirationTimestamp == 0 {
		return true
	}

	return provider.expirationTimestamp-time.Now().Unix() <= 180
}

func (provider *OIDCCredentialsProvider) GetCredentials() (cc *Credentials, err error) {
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

		provider.lastUpdateTimestamp = time.Now().Unix()
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

func (provider *OIDCCredentialsProvider) GetProviderName() string {
	return "oidc_role_arn"
}
