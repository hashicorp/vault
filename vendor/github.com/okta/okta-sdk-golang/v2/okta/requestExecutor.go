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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	nUrl "net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

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
		re.httpClient = &http.Client{Transport: tr}
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

			_, err = buildResponse(tokenResponse, &accessToken)
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
	requestStarted := time.Now().Unix()
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

		resp, err := re.doWithRetries(ctx, req, 0, requestStarted, nil)

		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 && req.Method == http.MethodGet && v != nil && reflect.TypeOf(v).Kind() != reflect.Slice {
			re.cache.Set(cacheKey, resp)
		}
		return buildResponse(resp, &v)
	}

	resp := re.cache.Get(cacheKey)
	return buildResponse(resp, &v)

}

func (re *RequestExecutor) doWithRetries(ctx context.Context, req *http.Request, retryCount int32, requestStarted int64, lastResponse *http.Response) (*http.Response, error) {
	iterationStart := time.Now().Unix()
	maxRetries := re.config.Okta.Client.RateLimit.MaxRetries
	requestTimeout := int64(re.config.Okta.Client.RequestTimeout)

	if req.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if requestTimeout > 0 && (iterationStart-requestStarted) >= requestTimeout {
		return lastResponse, errors.New("reached the max request time")
	}

	req = req.WithContext(ctx)
	resp, err := re.httpClient.Do(req)

	if (err != nil || tooManyRequests(resp)) && retryCount < maxRetries {
		if resp != nil {
			err := tryDrainBody(resp.Body)
			if err != nil {
				return nil, err
			}

			retryLimitReset := resp.Header.Get("X-Rate-Limit-Reset")
			date := resp.Header.Get("Date")
			if retryLimitReset == "" || date == "" {
				return resp, errors.New("a 429 response must include the x-retry-limit-reset and date headers")
			}

			if tooManyRequests(resp) {
				err := backoffPause(ctx, retryCount, resp)
				if err != nil {
					return nil, err
				}
			}
			retryCount++

			req.Header.Add("X-Okta-Retry-For", resp.Header.Get("X-Okta-Request-Id"))
			req.Header.Add("X-Okta-Retry-Count", fmt.Sprint(retryCount))

			resp, err = re.doWithRetries(ctx, req, retryCount, requestStarted, resp)
		}
	}

	return resp, err
}

func tooManyRequests(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusTooManyRequests
}

func tryDrainBody(body io.ReadCloser) error {
	defer body.Close()
	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, 4096))
	if err != nil {
		return err
	}
	return nil
}

func backoffPause(ctx context.Context, retryCount int32, response *http.Response) error {
	if response.StatusCode == http.StatusTooManyRequests {
		backoffSeconds := Get429BackoffTime(ctx, response)
		time.Sleep(time.Duration(backoffSeconds) * time.Second)

		return nil
	}

	return nil
}

func Get429BackoffTime(ctx context.Context, response *http.Response) int64 {
	var limitResetMap []int

	for _, time := range response.Header["X-Rate-Limit-Reset"] {
		timestamp, _ := strconv.Atoi(time)
		limitResetMap = append(limitResetMap, timestamp)
	}

	sort.Ints(limitResetMap)

	requestDate, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 Z", response.Header.Get("Date"))
	requestDateUnix := requestDate.Unix()
	backoffSeconds := int64(limitResetMap[0]) - requestDateUnix + 1
	return backoffSeconds
}

type Response struct {
	*http.Response
	Self     string
	NextPage string
}

func (r *Response) Next(ctx context.Context, v interface{}) (*Response, error) {
	client, _ := ClientFromContext(ctx)

	req, err := client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", r.NextPage, nil)
	if err != nil {
		return nil, err
	}

	return client.requestExecutor.Do(ctx, req, v)

}

func (r *Response) HasNextPage() bool {
	return r.NextPage != ""
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
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

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	e := new(Error)
	json.Unmarshal(bodyBytes, &e)
	return e

}

func buildResponse(resp *http.Response, v interface{}) (*Response, error) {
	ct := resp.Header.Get("Content-Type")

	if strings.Contains(ct, "application/xml") {
		return buildXmlResponse(resp, v)
	} else if strings.Contains(ct, "application/json") {
		return buildJsonResponse(resp, v)
	} else if ct == "" {
		return buildJsonResponse(resp, v)
	} else {
		return nil, errors.New("could not build a response for type: " + ct)
	}

}

func buildJsonResponse(resp *http.Response, v interface{}) (*Response, error) {
	response := newResponse(resp)

	err := CheckResponseForError(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		origResp := ioutil.NopCloser(bytes.NewBuffer(respBody))
		response.Body = origResp

		decodeError := json.NewDecoder(resp.Body).Decode(v)
		if decodeError == io.EOF {
			decodeError = nil
		}
		if decodeError != nil {
			return nil, decodeError
		}

	}
	return response, nil
}

func buildXmlResponse(resp *http.Response, v interface{}) (*Response, error) {
	response := newResponse(resp)

	err := CheckResponseForError(resp)
	if err != nil {
		return response, err
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	v = string(out)

	return response, nil
}
