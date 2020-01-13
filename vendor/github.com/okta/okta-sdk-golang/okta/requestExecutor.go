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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
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
	requestStarted := time.Now().Unix()
	cacheKey := cache.CreateCacheKey(req)
	if req.Method != http.MethodGet {
		re.cache.Delete(cacheKey)
	}
	inCache := re.cache.Has(cacheKey)

	if !inCache {

		resp, err := re.doWithRetries(req, 0, requestStarted, nil)

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
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 && req.Method == http.MethodGet && reflect.TypeOf(v).Kind() != reflect.Slice {
			re.cache.Set(cacheKey, resp)
		}
		return buildResponse(resp, &v)
	}

	resp := re.cache.Get(cacheKey)
	return buildResponse(resp, &v)

}

func (re *RequestExecutor) doWithRetries(req *http.Request, retryCount int32, requestStarted int64, lastResponse *http.Response) (*http.Response, error) {
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

	resp, err := re.httpClient.Do(req)

	if (err != nil || tooManyRequests(resp)) && retryCount < maxRetries {
		if resp != nil {
			err := tryDrainBody(resp.Body)
			if err != nil {
				return nil, err
			}
		}

		retryLimitReset := resp.Header.Get("X-Rate-Limit-Reset")
		date := resp.Header.Get("Date")
		if retryLimitReset == "" || date == "" {
			return resp, errors.New("a 429 response must include the x-retry-limit-reset and date headers")
		}

		if tooManyRequests(resp) {
			err := backoffPause(retryCount, resp)
			if err != nil {
				return nil, err
			}
		}
		retryCount++

		req.Header.Add("X-Okta-Retry-For", resp.Header.Get("X-Okta-Request-Id"))
		req.Header.Add("X-Okta-Retry-Count", fmt.Sprint(retryCount))

		resp, err = re.doWithRetries(req, retryCount, requestStarted, resp)
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

func backoffPause(retryCount int32, response *http.Response) error {
	if response.StatusCode == http.StatusTooManyRequests {
		backoffSeconds := Get429BackoffTime(response)
		time.Sleep(time.Duration(backoffSeconds) * time.Second)

		return nil
	}

	return nil
}

func Get429BackoffTime(response *http.Response) int64 {
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
