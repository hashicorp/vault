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
	"io/ioutil"
	"net/http"
	"net/url"
	nUrl "net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cenkalti/backoff/v4"
	"github.com/okta/okta-sdk-golang/v2/okta/cache"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type RequestExecutor struct {
	httpClient        *http.Client
	config            *config
	BaseUrl           *url.URL
	cache             cache.Cache
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
	ExpireIn    int    `json:"expire_in,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

func NewRequestExecutor(httpClient *http.Client, cache cache.Cache, config *config) *RequestExecutor {
	re := RequestExecutor{}
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
		buff = new(bytes.Buffer)
		encoder := json.NewEncoder(buff)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(body)
		if err != nil {
			return nil, err
		}
	}
	url = re.config.Okta.Client.OrgUrl + url

	req, err := http.NewRequest(method, url, buff)

	if err != nil {
		return nil, err
	}

	if re.config.Okta.Client.AuthorizationMode == "SSWS" {
		req.Header.Add("Authorization", "SSWS "+re.config.Okta.Client.Token)
	}

	if re.config.Okta.Client.AuthorizationMode == "PrivateKey" {
		if re.cache.Has("OKTA_ACCESS_TOKEN") {
			token := re.cache.GetString("OKTA_ACCESS_TOKEN")
			req.Header.Add("Authorization", "Bearer "+token)
		} else {
			priv := []byte(strings.ReplaceAll(re.config.Okta.Client.PrivateKey, `\n`, "\n"))

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

			signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: parsedKey}, nil)
			if err != nil {
				return nil, err
			}

			claims := ClientAssertionClaims{
				Subject:  re.config.Okta.Client.ClientId,
				IssuedAt: jwt.NewNumericDate(time.Now()),
				Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(1))),
				Issuer:   re.config.Okta.Client.ClientId,
				Audience: re.config.Okta.Client.OrgUrl + "/oauth2/v1/token",
			}
			jwtBuilder := jwt.Signed(signer).Claims(claims)
			clientAssertion, err := jwtBuilder.CompactSerialize()
			if err != nil {
				return nil, err
			}

			var tokenRequestBuff io.ReadWriter
			query := nUrl.Values{}
			tokenRequestUrl := re.config.Okta.Client.OrgUrl + "/oauth2/v1/token"

			query.Add("grant_type", "client_credentials")
			query.Add("scope", strings.Join(re.config.Okta.Client.Scopes, " "))
			query.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
			query.Add("client_assertion", clientAssertion)
			tokenRequestUrl += "?" + query.Encode()
			tokenRequest, err := http.NewRequest("POST", tokenRequestUrl, tokenRequestBuff)
			if err != nil {
				return nil, err
			}

			tokenRequest.Header.Add("Accept", "application/json")
			tokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			tokenResponse, err := re.httpClient.Do(tokenRequest)
			if err != nil {
				return nil, err
			}

			respBody, err := ioutil.ReadAll(tokenResponse.Body)
			if err != nil {
				return nil, err
			}
			origResp := ioutil.NopCloser(bytes.NewBuffer(respBody))
			tokenResponse.Body = origResp
			var accessToken *RequestAccessToken

			_, err = buildResponse(tokenResponse, nil, &accessToken)
			if err != nil {
				return nil, err
			}
			req.Header.Add("Authorization", "Bearer "+accessToken.AccessToken)

			re.cache.SetString("OKTA_ACCESS_TOKEN", accessToken.AccessToken)
		}

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

		resp, err := re.doWithRetries(ctx, req)

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
	err                    error
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

func (re *RequestExecutor) doWithRetries(ctx context.Context, req *http.Request) (*http.Response, error) {
	var bodyReader func() io.ReadCloser
	if req.Body != nil {
		buf, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = func() io.ReadCloser {
			return ioutil.NopCloser(bytes.NewReader(buf))
		}
	}
	var (
		resp *http.Response
		err  error
	)
	if re.config.Okta.Client.RequestTimeout > 0 {
		ctx, _ = context.WithTimeout(ctx, time.Second*time.Duration(re.config.Okta.Client.RequestTimeout))
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
			return fmt.Errorf("network error: %v", err)
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
		return errors.New("to many requests")
	}
	err = backoff.Retry(operation, bOff)
	return resp, err
}

func tooManyRequests(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusTooManyRequests
}

func tryDrainBody(body io.ReadCloser) error {
	defer body.Close()
	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, 4096))
	return err
}

func Get429BackoffTime(resp *http.Response) (int64, error) {
	requestDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", resp.Header.Get("Date"))
	if err != nil {
		// this is error is considered to be permanent and should not be retried
		return 0, backoff.Permanent(errors.New(fmt.Sprintf("Date header is missing or invalid: %v", err)))
	}
	rateLimitReset, err := strconv.Atoi(resp.Header.Get("X-Rate-Limit-Reset"))
	if err != nil {
		// this is error is considered to be permanent and should not be retried
		return 0, backoff.Permanent(errors.New(fmt.Sprintf("X-Rate-Limit-Reset header is missing or invalid: %v", err)))
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
			rawUrl, _ := url.Parse(rawLink)
			rawUrl.Scheme = ""
			rawUrl.Host = ""

			if strings.Contains(link, `rel="self"`) {
				response.Self = rawUrl.String()
			}

			if strings.Contains(link, `rel="next"`) {
				response.NextPage = rawUrl.String()
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
	if statusCode == http.StatusUnauthorized && strings.Contains(resp.Header.Get("WWW-Authenticate"), "Bearer") {
		for _, v := range strings.Split(resp.Header.Get("WWW-Authenticate"), ", ") {
			if strings.Contains(v, "error_description") {
				_, err := toml.Decode(v, &e)
				if err != nil {
					e.ErrorSummary = "unauthorized"
				}
				return &e
			}
		}
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	copyBodyBytes := make([]byte, len(bodyBytes))
	copy(copyBodyBytes, bodyBytes)
	_ = resp.Body.Close()
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	_ = json.NewDecoder(bytes.NewReader(copyBodyBytes)).Decode(&e)
	return &e
}

func buildResponse(resp *http.Response, re *RequestExecutor, v interface{}) (*Response, error) {
	ct := resp.Header.Get("Content-Type")
	response := newResponse(resp, re)
	err := CheckResponseForError(resp)
	if err != nil {
		return response, err
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	copyBodyBytes := make([]byte, len(bodyBytes))
	copy(copyBodyBytes, bodyBytes)
	_ = resp.Body.Close()                                    // close it to avoid memory leaks
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // restore the original response body
	if strings.Contains(ct, "application/xml") {
		err = xml.NewDecoder(bytes.NewReader(copyBodyBytes)).Decode(v)
	} else if strings.Contains(ct, "application/json") || ct == "" {
		err = json.NewDecoder(bytes.NewReader(copyBodyBytes)).Decode(v)
	} else if strings.Contains(ct, "application/octet-stream") {
		// since the response is arbitrary binary data, we leave it to the user to decode it
		return response, nil
	} else {
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
