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
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/okta/okta-sdk-golang/okta/cache"
)

type RequestExecutor struct {
	httpClient *http.Client
	config     *config
	BaseUrl    *url.URL
	cache      cache.Cache
}

func NewRequestExecutor(httpClient *http.Client, cache cache.Cache, config *config) *RequestExecutor {
	re := RequestExecutor{}
	re.httpClient = httpClient
	re.config = config
	re.cache = cache

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
	req.Header.Add("Authorization", "SSWS "+re.config.Okta.Client.Token)
	req.Header.Add("User-Agent", NewUserAgent(re.config).String())
	req.Header.Add("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (re *RequestExecutor) Do(req *http.Request, v interface{}) (*Response, error) {
	cacheKey := cache.CreateCacheKey(req)
	if req.Method != http.MethodGet {
		re.cache.Delete(cacheKey)
	}
	inCache := re.cache.Has(cacheKey)

	if !inCache {
		resp, err := re.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		origResp := ioutil.NopCloser(bytes.NewBuffer(respBody))
		resp.Body = origResp

		if req.Method == http.MethodGet && reflect.TypeOf(v).Kind() != reflect.Slice {
			re.cache.Set(cacheKey, resp)
		}

		return buildResponse(resp, &v)

	}

	resp := re.cache.Get(cacheKey)
	return buildResponse(resp, &v)

}

type Response struct {
	*http.Response
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
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
	response := newResponse(resp)

	err := CheckResponseForError(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		decodeError := json.NewDecoder(resp.Body).Decode(v)
		if decodeError == io.EOF {
			decodeError = nil
		}
		if decodeError != nil {
			err = decodeError
		}

	}
	return response, err
}
