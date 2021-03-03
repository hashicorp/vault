/**
 * Copyright 2016 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
)

type RestTransport struct{}

// DoRequest - Implementation of the TransportHandler interface for handling
// calls to the REST endpoint.
func (r *RestTransport) DoRequest(sess *Session, service string, method string, args []interface{}, options *sl.Options, pResult interface{}) error {
	restMethod := httpMethod(method, args)

	// Parse any method parameters and determine the HTTP method
	var parameters []byte
	if len(args) > 0 {
		// parse the parameters
		parameters, _ = json.Marshal(
			map[string]interface{}{
				"parameters": args,
			})
	}

	path := buildPath(service, method, options)

	resp, code, err := sendHTTPRequest(
		sess,
		path,
		restMethod,
		bytes.NewBuffer(parameters),
		options)

	if err != nil {
		//Preserve the original sl error
		if _, ok := err.(sl.Error); ok {
			return err
		}
		return sl.Error{Wrapped: err, StatusCode: code}
	}

	err = findResponseError(code, resp)
	if err != nil {
		return err
	}

	// Some APIs that normally return a collection, omit the []'s when the API returns a single value
	returnType := reflect.TypeOf(pResult).String()
	if strings.Index(returnType, "[]") == 1 && strings.Index(string(resp), "[") != 0 {
		resp = []byte("[" + string(resp) + "]")
	}

	// At this point, all that's left to do is parse the return value to the appropriate type, and return
	// any parse errors (or nil if successful)

	err = nil
	switch pResult.(type) {
	case *[]uint8:
		// exclude quotes
		*pResult.(*[]uint8) = resp[1 : len(resp)-1]
	case *datatypes.Void:
	case *uint:
		var val uint64
		val, err = strconv.ParseUint(string(resp), 0, 64)
		if err == nil {
			*pResult.(*uint) = uint(val)
		}
	case *bool:
		*pResult.(*bool), err = strconv.ParseBool(string(resp))
	case *string:
		str := string(resp)
		strIdx := len(str) - 1
		if str == "null" {
			str = ""
		} else if str[0] == '"' && str[strIdx] == '"' {
			rawStr := rawString{str}
			err = json.Unmarshal([]byte(`{"val":`+str+`}`), &rawStr)
			if err == nil {
				str = rawStr.Val
			}
		}
		*pResult.(*string) = str
	default:
		// Must be a json representation of one of the many softlayer datatypes
		err = json.Unmarshal(resp, pResult)
	}

	if err != nil {
		err = sl.Error{Message: err.Error(), Wrapped: err}
	}

	return err
}

type rawString struct {
	Val string
}

func buildPath(service string, method string, options *sl.Options) string {
	path := service

	if options.Id != nil {
		path = path + "/" + strconv.Itoa(*options.Id)
	}

	// omit the API method name if the method represents one of the basic REST methods
	if method != "getObject" && method != "deleteObject" && method != "createObject" &&
		method != "editObject" && method != "editObjects" {
		path = path + "/" + method
	}

	return path + ".json"
}

func encodeQuery(opts *sl.Options) string {
	query := new(url.URL).Query()

	if opts.Mask != "" {
		query.Add("objectMask", opts.Mask)
	}

	if opts.Filter != "" {
		query.Add("objectFilter", opts.Filter)
	}

	// resultLimit=<offset>,<limit>
	// If offset unspecified, default to 0
	if opts.Limit != nil {
		startOffset := 0
		if opts.Offset != nil {
			startOffset = *opts.Offset
		}

		query.Add("resultLimit", fmt.Sprintf("%d,%d", startOffset, *opts.Limit))
	}

	return query.Encode()
}

func sendHTTPRequest(
	sess *Session, path string, requestType string,
	requestBody *bytes.Buffer, options *sl.Options) ([]byte, int, error) {

	retries := sess.Retries
	if retries < 2 {
		return makeHTTPRequest(sess, path, requestType, requestBody, options)
	}

	wait := sess.RetryWait
	if wait == 0 {
		wait = DefaultRetryWait
	}

	return tryHTTPRequest(retries, wait, sess, path, requestType, requestBody, options)
}

func tryHTTPRequest(
	retries int, wait time.Duration, sess *Session,
	path string, requestType string, requestBody *bytes.Buffer,
	options *sl.Options) ([]byte, int, error) {

	resp, code, err := makeHTTPRequest(sess, path, requestType, requestBody, options)
	if err != nil {
		if !isRetryable(err) {
			return resp, code, err
		}

		if retries--; retries > 0 {
			jitter := time.Duration(rand.Int63n(int64(wait)))
			wait = wait + jitter/2
			time.Sleep(wait)
			return tryHTTPRequest(
				retries, wait, sess, path, requestType, requestBody, options)
		}
	}

	return resp, code, err
}

func makeHTTPRequest(
	session *Session, path string, requestType string,
	requestBody *bytes.Buffer, options *sl.Options) ([]byte, int, error) {
	log := Logger

	client := session.HTTPClient
	if client == nil {
		client = &http.Client{}
	}

	client.Timeout = DefaultTimeout
	if session.Timeout != 0 {
		client.Timeout = session.Timeout
	}

	var url string
	if session.Endpoint == "" {
		url = url + DefaultEndpoint
	} else {
		url = url + session.Endpoint
	}
	url = fmt.Sprintf("%s/%s", strings.TrimRight(url, "/"), path)
	req, err := http.NewRequest(requestType, url, requestBody)
	if err != nil {
		return nil, 0, err
	}

	if session.APIKey != "" {
		req.SetBasicAuth(session.UserName, session.APIKey)
	} else if session.AuthToken != "" {
		req.SetBasicAuth(fmt.Sprintf("%d", session.UserId), session.AuthToken)
	}

	// For cases where session is built from the raw structure and not using New() , the UserAgent would be empty
	if session.userAgent == "" {
		session.userAgent = getDefaultUserAgent()
	}

	req.Header.Set("User-Agent", session.userAgent)

	if session.Headers != nil {
		for key, value := range session.Headers {
			req.Header.Set(key, value)
		}
	}

	req.URL.RawQuery = encodeQuery(options)

	if session.Debug {
		log.Println("[DEBUG] Request URL: ", requestType, req.URL)
		log.Println("[DEBUG] Parameters: ", requestBody.String())
	}

	resp, err := client.Do(req)
	if err != nil {
		statusCode := 520
		if resp != nil && resp.StatusCode != 0 {
			statusCode = resp.StatusCode
		}

		if isTimeout(err) {
			statusCode = 599
		}

		return nil, statusCode, err
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	if session.Debug {
		log.Println("[DEBUG] Status Code: ", resp.StatusCode)
		log.Println("[DEBUG] Response: ", string(responseBody))
	}
	err = findResponseError(resp.StatusCode, responseBody)
	return responseBody, resp.StatusCode, err
}

func httpMethod(name string, args []interface{}) string {
	if name == "deleteObject" {
		return "DELETE"
	} else if name == "editObject" || name == "editObjects" {
		return "PUT"
	} else if name == "createObject" || name == "createObjects" || len(args) > 0 {
		return "POST"
	}

	return "GET"
}

func findResponseError(code int, resp []byte) error {
	if code < 200 || code > 299 {
		e := sl.Error{StatusCode: code}
		err := json.Unmarshal(resp, &e)
		// If unparseable, wrap the json error
		if err != nil {
			e.Wrapped = err
			e.Message = err.Error()
		}
		return e
	}
	return nil
}
