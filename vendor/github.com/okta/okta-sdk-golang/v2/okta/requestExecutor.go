/*
 * Copyright 2018 - Present Okta, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package okta

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/okta/okta-sdk-golang/v2/okta/cache"
	goCache "github.com/patrickmn/go-cache"
)

const AccessTokenCacheKey = "OKTA_ACCESS_TOKEN"

type RequestExecutor struct {
	httpClient        *http.Client
	config            *config
	BaseUrl           *urlpkg.URL
	cache             cache.Cache
	tokenCache        *goCache.Cache
	binary            bool
	headerAccept      string
	headerContentType string
	freshCache        bool
}

type ClientAssertionClaims struct {
	Issuer   string           `json:"iss,omitempty"`
	Subject  string           `json:"sub,omitempty"`
	Audience string           `json:"aud,omitempty"`
	Expiry   *jwt.NumericDate `json:"exp,omitempty"`
	IssuedAt *jwt.NumericDate `json:"iat,omitempty"`
	ID       string           `json:"jti,omitempty"`
}

type RequestAccessToken struct {
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

type Authorization interface {
	Authorize() error
}

type SSWSAuth struct {
	token string
	req   *http.Request
}

func NewSSWSAuth(token string, req *http.Request) *SSWSAuth {
	return &SSWSAuth{token: token, req: req}
}

func (a *SSWSAuth) Authorize() error {
	a.req.Header.Add("Authorization", "SSWS "+a.token)
	return nil
}

type BearerAuth struct {
	token string
	req   *http.Request
}

func NewBearerAuth(token string, req *http.Request) *BearerAuth {
	return &BearerAuth{token: token, req: req}
}

func (a *BearerAuth) Authorize() error {
	a.req.Header.Add("Authorization", "Bearer "+a.token)
	return nil
}

type PrivateKeyAuth struct {
	tokenCache       *goCache.Cache
	httpClient       *http.Client
	privateKeySigner jose.Signer
	privateKey       string
	privateKeyId     string
	clientId         string
	orgURL           string
	scopes           []string
	maxRetries       int32
	maxBackoff       int64
	req              *http.Request
}

type PrivateKeyAuthConfig struct {
	TokenCache       *goCache.Cache
	HttpClient       *http.Client
	PrivateKeySigner jose.Signer
	PrivateKey       string
	PrivateKeyId     string
	ClientId         string
	OrgURL           string
	Scopes           []string
	MaxRetries       int32
	MaxBackoff       int64
	Req              *http.Request
}

func NewPrivateKeyAuth(config PrivateKeyAuthConfig) *PrivateKeyAuth {
	return &PrivateKeyAuth{
		tokenCache:       config.TokenCache,
		httpClient:       config.HttpClient,
		privateKeySigner: config.PrivateKeySigner,
		privateKey:       config.PrivateKey,
		privateKeyId:     config.PrivateKeyId,
		clientId:         config.ClientId,
		orgURL:           config.OrgURL,
		scopes:           config.Scopes,
		maxRetries:       config.MaxRetries,
		maxBackoff:       config.MaxBackoff,
		req:              config.Req,
	}
}

func (a *PrivateKeyAuth) Authorize() error {
	accessToken, hasToken := a.tokenCache.Get(AccessTokenCacheKey)
	if hasToken {
		a.req.Header.Add("Authorization", "Bearer "+accessToken.(string))
	} else {
		if a.privateKeySigner == nil {
			var err error
			a.privateKeySigner, err = CreateKeySigner(a.privateKey, a.privateKeyId)
			if err != nil {
				return err
			}
		}

		clientAssertion, err := CreateClientAssertion(a.orgURL, a.clientId, a.privateKeySigner)
		if err != nil {
			return err
		}

		accessToken, err := getAccessTokenForPrivateKey(a.httpClient, a.orgURL, clientAssertion, a.scopes, a.maxRetries, a.maxBackoff)
		if err != nil {
			return err
		}

		a.req.Header.Add("Authorization", "Bearer "+accessToken.AccessToken)

		// Trim a couple of seconds off calculated expiry so cache expiry
		// occures before Okta server side expiry.
		expiration := accessToken.ExpiresIn - 2
		a.tokenCache.Set(AccessTokenCacheKey, accessToken.AccessToken, time.Second*time.Duration(expiration))
	}
	return nil
}

type JWTAuth struct {
	tokenCache      *goCache.Cache
	httpClient      *http.Client
	orgURL          string
	scopes          []string
	clientAssertion string
	maxRetries      int32
	maxBackoff      int64
	req             *http.Request
}

type JWTAuthConfig struct {
	TokenCache      *goCache.Cache
	HttpClient      *http.Client
	OrgURL          string
	Scopes          []string
	ClientAssertion string
	MaxRetries      int32
	MaxBackoff      int64
	Req             *http.Request
}

func NewJWTAuth(config JWTAuthConfig) *JWTAuth {
	return &JWTAuth{
		tokenCache:      config.TokenCache,
		httpClient:      config.HttpClient,
		orgURL:          config.OrgURL,
		scopes:          config.Scopes,
		clientAssertion: config.ClientAssertion,
		maxRetries:      config.MaxRetries,
		maxBackoff:      config.MaxBackoff,
		req:             config.Req,
	}
}

func (a *JWTAuth) Authorize() error {
	accessToken, hasToken := a.tokenCache.Get(AccessTokenCacheKey)
	if hasToken {
		a.req.Header.Add("Authorization", "Bearer "+accessToken.(string))
	} else {
		accessToken, err := getAccessTokenForPrivateKey(a.httpClient, a.orgURL, a.clientAssertion, a.scopes, a.maxRetries, a.maxBackoff)
		if err != nil {
			return err
		}
		a.req.Header.Add("Authorization", "Bearer "+accessToken.AccessToken)

		// Trim a couple of seconds off calculated expiry so cache expiry
		// occures before Okta server side expiry.
		expiration := accessToken.ExpiresIn - 2
		a.tokenCache.Set(AccessTokenCacheKey, accessToken.AccessToken, time.Second*time.Duration(expiration))
	}
	return nil
}

func CreateKeySigner(privateKey, privateKeyID string) (jose.Signer, error) {
	priv := []byte(strings.ReplaceAll(privateKey, `\n`, "\n"))

	privPem, _ := pem.Decode(priv)
	if privPem == nil {
		return nil, errors.New("invalid private key")
	}
	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("RSA private key is of the wrong type")
	}

	parsedKey, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		return nil, err
	}

	var signerOptions *jose.SignerOptions
	if privateKeyID != "" {
		signerOptions = (&jose.SignerOptions{}).WithHeader("kid", privateKeyID)
	}

	return jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: parsedKey}, signerOptions)
}

func CreateClientAssertion(orgURL, clientID string, privateKeySinger jose.Signer) (clientAssertion string, err error) {
	claims := ClientAssertionClaims{
		Subject:  clientID,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(1))),
		Issuer:   clientID,
		Audience: orgURL + "/oauth2/v1/token",
	}
	jwtBuilder := jwt.Signed(privateKeySinger).Claims(claims)
	return jwtBuilder.CompactSerialize()
}

func getAccessTokenForPrivateKey(httpClient *http.Client, orgURL, clientAssertion string, scopes []string, maxRetries int32, maxBackoff int64) (*RequestAccessToken, error) {
	var tokenRequestBuff io.ReadWriter
	query := urlpkg.Values{}
	tokenRequestURL := orgURL + "/oauth2/v1/token"

	query.Add("grant_type", "client_credentials")
	query.Add("scope", strings.Join(scopes, " "))
	query.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	query.Add("client_assertion", clientAssertion)
	tokenRequestURL += "?" + query.Encode()
	tokenRequest, err := http.NewRequest("POST", tokenRequestURL, tokenRequestBuff)
	if err != nil {
		return nil, err
	}

	tokenRequest.Header.Add("Accept", "application/json")
	tokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	bOff := &oktaBackoff{
		ctx:             context.TODO(),
		maxRetries:      maxRetries,
		backoffDuration: time.Duration(maxBackoff),
	}
	var tokenResponse *http.Response
	operation := func() error {
		tokenResponse, err = httpClient.Do(tokenRequest)
		bOff.retryCount++
		return err
	}
	err = backoff.Retry(operation, bOff)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(tokenResponse.Body)
	if err != nil {
		return nil, err
	}
	origResp := io.NopCloser(bytes.NewBuffer(respBody))
	tokenResponse.Body = origResp
	var accessToken *RequestAccessToken

	_, err = buildResponse(tokenResponse, nil, &accessToken)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func NewRequestExecutor(httpClient *http.Client, cache cache.Cache, config *config) *RequestExecutor {
	re := RequestExecutor{
		tokenCache: goCache.New(5*time.Minute, 10*time.Minute),
	}

	re.httpClient = httpClient
	re.config = config
	re.cache = cache
	re.binary = false
	re.headerAccept = "application/json"
	re.headerContentType = "application/json"

	if httpClient == nil {
		tr := &http.Transport{
			IdleConnTimeout: 30 * time.Second,
		}
		re.httpClient = &http.Client{
			Transport: tr,
			Timeout:   time.Second * time.Duration(re.config.Okta.Client.ConnectionTimeout),
		}
	}

	return &re
}

func (re *RequestExecutor) NewRequest(method string, url string, body interface{}) (*http.Request, error) {
	var buff io.ReadWriter
	if body != nil {
		switch v := body.(type) {
		case []byte:
			buff = bytes.NewBuffer(v)
		case *bytes.Buffer:
			buff = v
		default:
			buff = new(bytes.Buffer)
			encoder := json.NewEncoder(buff)
			encoder.SetEscapeHTML(false)
			err := encoder.Encode(body)
			if err != nil {
				return nil, err
			}
		}
	}
	url = re.config.Okta.Client.OrgUrl + url

	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return nil, err
	}

	var auth Authorization

	switch re.config.Okta.Client.AuthorizationMode {
	case "SSWS":
		auth = NewSSWSAuth(re.config.Okta.Client.Token, req)
	case "Bearer":
		auth = NewBearerAuth(re.config.Okta.Client.Token, req)
	case "PrivateKey":
		auth = NewPrivateKeyAuth(PrivateKeyAuthConfig{
			TokenCache:       re.tokenCache,
			HttpClient:       re.httpClient,
			PrivateKeySigner: re.config.PrivateKeySigner,
			PrivateKey:       re.config.Okta.Client.PrivateKey,
			PrivateKeyId:     re.config.Okta.Client.PrivateKeyId,
			ClientId:         re.config.Okta.Client.ClientId,
			OrgURL:           re.config.Okta.Client.OrgUrl,
			Scopes:           re.config.Okta.Client.Scopes,
			MaxRetries:       re.config.Okta.Client.RateLimit.MaxRetries,
			MaxBackoff:       re.config.Okta.Client.RateLimit.MaxBackoff,
			Req:              req,
		})
	case "JWT":
		auth = NewJWTAuth(JWTAuthConfig{
			TokenCache:      re.tokenCache,
			HttpClient:      re.httpClient,
			OrgURL:          re.config.Okta.Client.OrgUrl,
			Scopes:          re.config.Okta.Client.Scopes,
			ClientAssertion: re.config.Okta.Client.ClientAssertion,
			MaxRetries:      re.config.Okta.Client.RateLimit.MaxRetries,
			MaxBackoff:      re.config.Okta.Client.RateLimit.MaxBackoff,
			Req:             req,
		})
	default:
		return nil, fmt.Errorf("unknown authorization mode %v", re.config.Okta.Client.AuthorizationMode)
	}

	err = auth.Authorize()
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", NewUserAgent(re.config).String())
	req.Header.Add("Accept", re.headerAccept)

	if body != nil {
		req.Header.Set("Content-Type", re.headerContentType)
	}

	// Force reset defaults
	re.binary = false
	re.headerAccept = "application/json"
	re.headerContentType = "application/json"
	return req, nil
}

func (re *RequestExecutor) AsBinary() *RequestExecutor {
	re.binary = true
	return re
}

func (re *RequestExecutor) WithAccept(acceptHeader string) *RequestExecutor {
	re.headerAccept = acceptHeader
	return re
}

func (re *RequestExecutor) WithContentType(contentTypeHeader string) *RequestExecutor {
	re.headerContentType = contentTypeHeader
	return re
}

func (re *RequestExecutor) RefreshNext() *RequestExecutor {
	re.freshCache = true
	return re
}

func (re *RequestExecutor) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	cacheKey := cache.CreateCacheKey(req)
	if req.Method != http.MethodGet {
		re.cache.Delete(cacheKey)
	}
	inCache := re.cache.Has(cacheKey)
	if re.freshCache {
		re.cache.Delete(cacheKey)
		inCache = false
		re.freshCache = false
	}
	if !inCache {
		resp, done, err := re.doWithRetries(ctx, req)
		defer done()
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 && req.Method == http.MethodGet && v != nil && reflect.TypeOf(v).Kind() != reflect.Slice {
			re.cache.Set(cacheKey, resp)
		}
		return buildResponse(resp, re, &v)
	}
	resp := re.cache.Get(cacheKey)
	return buildResponse(resp, re, &v)
}

type oktaBackoff struct {
	retryCount, maxRetries int32
	backoffDuration        time.Duration
	ctx                    context.Context
}

// NextBackOff returns the duration to wait before retrying the operation,
// or backoff. Stop to indicate that no more retries should be made.
func (o *oktaBackoff) NextBackOff() time.Duration {
	// stop retrying if operation reached retry limit
	if o.retryCount > o.maxRetries {
		return backoff.Stop
	}
	return o.backoffDuration
}

// Reset to initial state.
func (o *oktaBackoff) Reset() {}

func (o *oktaBackoff) Context() context.Context {
	return o.ctx
}

func (re *RequestExecutor) doWithRetries(ctx context.Context, req *http.Request) (*http.Response, func(), error) {
	var bodyReader func() io.ReadCloser
	done := func() {}
	if req.Body != nil {
		buf, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, done, err
		}
		bodyReader = func() io.ReadCloser {
			return io.NopCloser(bytes.NewReader(buf))
		}
	}
	var (
		resp *http.Response
		err  error
	)
	if re.config.Okta.Client.RequestTimeout > 0 {
		ctx, done = context.WithTimeout(ctx, time.Second*time.Duration(re.config.Okta.Client.RequestTimeout))
	}
	bOff := &oktaBackoff{
		ctx:        ctx,
		maxRetries: re.config.Okta.Client.RateLimit.MaxRetries,
	}
	operation := func() error {
		// Always rewind the request body when non-nil.
		if bodyReader != nil {
			req.Body = bodyReader()
		}
		resp, err = re.httpClient.Do(req.WithContext(ctx))
		if errors.Is(err, io.EOF) {
			// retry on EOF errors, which might be caused by network connectivity issues
			return fmt.Errorf("network error: %w", err)
		} else if err != nil {
			// this is error is considered to be permanent and should not be retried
			return backoff.Permanent(err)
		}
		if !tooManyRequests(resp) {
			return nil
		}
		if err = tryDrainBody(resp.Body); err != nil {
			return err
		}
		backoffDuration, err := Get429BackoffTime(resp)
		if err != nil {
			return err
		}
		if re.config.Okta.Client.RateLimit.MaxBackoff < backoffDuration {
			backoffDuration = re.config.Okta.Client.RateLimit.MaxBackoff
		}
		bOff.backoffDuration = time.Second * time.Duration(backoffDuration)
		bOff.retryCount++
		req.Header.Add("X-Okta-Retry-For", resp.Header.Get("X-Okta-Request-Id"))
		req.Header.Add("X-Okta-Retry-Count", fmt.Sprint(bOff.retryCount))
		return errors.New("too many requests")
	}
	err = backoff.Retry(operation, bOff)
	return resp, done, err
}

func tooManyRequests(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusTooManyRequests
}

func tryDrainBody(body io.ReadCloser) error {
	defer body.Close()
	_, err := io.Copy(io.Discard, io.LimitReader(body, 4096))
	return err
}

func Get429BackoffTime(resp *http.Response) (int64, error) {
	requestDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", resp.Header.Get("Date"))
	if err != nil {
		// this is error is considered to be permanent and should not be retried
		return 0, backoff.Permanent(fmt.Errorf("date header is missing or invalid: %w", err))
	}
	rateLimitReset, err := strconv.Atoi(resp.Header.Get("X-Rate-Limit-Reset"))
	if err != nil {
		// this is error is considered to be permanent and should not be retried
		return 0, backoff.Permanent(fmt.Errorf("X-Rate-Limit-Reset header is missing or invalid: %w", err))
	}
	return int64(rateLimitReset) - requestDate.Unix() + 1, nil
}

type Response struct {
	*http.Response
	re       *RequestExecutor
	Self     string
	NextPage string
}

func (r *Response) Next(ctx context.Context, v interface{}) (*Response, error) {
	if r.re == nil {
		return nil, errors.New("no initial response provided from previous request")
	}
	req, err := r.re.NewRequest("GET", r.NextPage, nil)
	if err != nil {
		return nil, err
	}
	return r.re.Do(ctx, req, v)
}

func (r *Response) HasNextPage() bool {
	return r.NextPage != ""
}

func newResponse(r *http.Response, re *RequestExecutor) *Response {
	response := &Response{Response: r, re: re}
	links := r.Header["Link"]

	if len(links) > 0 {
		for _, link := range links {
			splitLinkHeader := strings.Split(link, ";")
			if len(splitLinkHeader) < 2 {
				continue
			}
			rawLink := strings.TrimRight(strings.TrimLeft(splitLinkHeader[0], "<"), ">")
			rawURL, _ := urlpkg.Parse(rawLink)
			rawURL.Scheme = ""
			rawURL.Host = ""
			if r.Request != nil {
				q := r.Request.URL.Query()
				for k, v := range rawURL.Query() {
					q.Set(k, v[0])
				}
				rawURL.RawQuery = q.Encode()
			}
			if strings.Contains(link, `rel="self"`) {
				response.Self = rawURL.String()
			}
			if strings.Contains(link, `rel="next"`) {
				response.NextPage = rawURL.String()
			}
		}
	}

	return response
}

func CheckResponseForError(resp *http.Response) error {
	statusCode := resp.StatusCode
	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		return nil
	}
	e := Error{}
	if (statusCode == http.StatusUnauthorized || statusCode == http.StatusForbidden) &&
		strings.Contains(resp.Header.Get("Www-Authenticate"), "Bearer") {
		for _, v := range strings.Split(resp.Header.Get("Www-Authenticate"), ", ") {
			if strings.Contains(v, "error_description") {
				_, err := toml.Decode(v, &e)
				if err != nil {
					e.ErrorSummary = "unauthorized"
				}
				return &e
			}
		}
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	copyBodyBytes := make([]byte, len(bodyBytes))
	copy(copyBodyBytes, bodyBytes)
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	_ = json.NewDecoder(bytes.NewReader(copyBodyBytes)).Decode(&e)
	if statusCode == http.StatusInternalServerError {
		e.ErrorSummary += fmt.Sprintf(", x-okta-request-id=%s", resp.Header.Get("x-okta-request-id"))
	}
	return &e
}

func buildResponse(resp *http.Response, re *RequestExecutor, v interface{}) (*Response, error) {
	ct := resp.Header.Get("Content-Type")
	response := newResponse(resp, re)
	err := CheckResponseForError(resp)
	if err != nil {
		return response, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	copyBodyBytes := make([]byte, len(bodyBytes))
	copy(copyBodyBytes, bodyBytes)
	_ = resp.Body.Close()                                // close it to avoid memory leaks
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // restore the original response body
	if len(copyBodyBytes) == 0 {
		return response, nil
	}
	switch {
	case strings.Contains(ct, "application/xml"):
		err = xml.NewDecoder(bytes.NewReader(copyBodyBytes)).Decode(v)
	case strings.Contains(ct, "application/json"):
		err = json.NewDecoder(bytes.NewReader(copyBodyBytes)).Decode(v)
	case strings.Contains(ct, "application/octet-stream"):
		// since the response is arbitrary binary data, we leave it to the user to decode it
		return response, nil
	default:
		return nil, errors.New("could not build a response for type: " + ct)
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}
