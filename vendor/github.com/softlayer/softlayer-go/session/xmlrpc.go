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
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/renier/xmlrpc"
	"github.com/softlayer/softlayer-go/sl"
)

// Debugging RoundTripper
type debugRoundTripper struct{}

func (mrt debugRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	log := Logger
	log.Println("->>>Request:")
	dumpedReq, _ := httputil.DumpRequestOut(request, true)
	log.Println(string(dumpedReq))

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		log.Println("Error:", err)
		return response, err
	}

	log.Println("\n\n<<<-Response:")
	dumpedResp, _ := httputil.DumpResponse(response, true)
	log.Println(string(dumpedResp))

	return response, err
}

// XML-RPC Transport
type XmlRpcTransport struct{}

func (x *XmlRpcTransport) DoRequest(
	sess *Session,
	service string,
	method string,
	args []interface{},
	options *sl.Options,
	pResult interface{},
) error {

	var err error
	serviceUrl := fmt.Sprintf("%s/%s", strings.TrimRight(sess.Endpoint, "/"), service)

	timeout := DefaultTimeout
	if sess.Timeout != 0 {
		timeout = sess.Timeout
	}

	// Declaring client outside of the if /else. So we can set the correct http transport based if it is TLS or not
	var client *xmlrpc.Client
	if sess.HTTPClient != nil && sess.HTTPClient.Transport != nil {
		client, err = xmlrpc.NewClient(serviceUrl, sess.HTTPClient.Transport, timeout)
	} else {
		var roundTripper http.RoundTripper
		if sess.Debug {
			roundTripper = debugRoundTripper{}
		}

		client, err = xmlrpc.NewClient(serviceUrl, roundTripper, timeout)
	}
	//Verify no errors happened in creating the xmlrpc client
	if err != nil {
		return fmt.Errorf("Could not create an xmlrpc client for %s: %s", service, err)
	}

	authenticate := map[string]interface{}{}
	if sess.UserName != "" {
		authenticate["username"] = sess.UserName
	}

	if sess.APIKey != "" {
		authenticate["apiKey"] = sess.APIKey
	}

	if sess.UserId != 0 {
		authenticate["userId"] = sess.UserId
		authenticate["complexType"] = "PortalLoginToken"
	}

	if sess.AuthToken != "" {
		authenticate["authToken"] = sess.AuthToken
		authenticate["complexType"] = "PortalLoginToken"
	}

	// For cases where session is built from the raw structure and not using New() , the UserAgent would be empty
	if sess.userAgent == "" {
		sess.userAgent = getDefaultUserAgent()
	}

	headers := map[string]interface{}{}
	headers["User-Agent"] = sess.userAgent

	if len(authenticate) > 0 {
		headers["authenticate"] = authenticate
	}

	if options.Id != nil {
		headers[fmt.Sprintf("%sInitParameters", service)] = map[string]int{
			"id": *options.Id,
		}
	}

	mask := options.Mask
	if mask != "" {
		if !strings.HasPrefix(mask, "mask[") && !strings.Contains(mask, ";") && strings.Contains(mask, ",") {
			mask = fmt.Sprintf("mask[%s]", mask)
			headers["SoftLayer_ObjectMask"] = map[string]string{"mask": mask}
		} else {
			headers[fmt.Sprintf("%sObjectMask", service)] =
				map[string]interface{}{"mask": genXMLMask(mask)}
		}
	}

	if options.Filter != "" {
		// FIXME: This json unmarshaling presents a performance problem,
		// since the filter builder marshals a data structure to json.
		// This then undoes that step to pass it to the xmlrpc request.
		// It would be better to get the umarshaled data structure
		// from the filter builder, but that will require changes to the
		// public API in Options.
		objFilter := map[string]interface{}{}
		err := json.Unmarshal([]byte(options.Filter), &objFilter)
		if err != nil {
			return fmt.Errorf("Error encoding object filter: %s", err)
		}
		headers[fmt.Sprintf("%sObjectFilter", service)] = objFilter
	}

	if options.Limit != nil {
		offset := 0
		if options.Offset != nil {
			offset = *options.Offset
		}

		headers["resultLimit"] = map[string]int{
			"limit":  *options.Limit,
			"offset": offset,
		}
	}

	// Add incoming arguments to xmlrpc parameter array
	params := []interface{}{}

	if len(headers) > 0 {
		params = append(params, map[string]interface{}{"headers": headers})
	}

	for _, arg := range args {
		params = append(params, arg)
	}

	retries := sess.Retries
	if retries < 2 {
		err = client.Call(method, params, pResult)
	} else {
		wait := sess.RetryWait
		if wait == 0 {
			wait = DefaultRetryWait
		}

		err = makeXmlRequest(retries, wait, client, method, params, pResult)
	}

	if xmlRpcError, ok := err.(*xmlrpc.XmlRpcError); ok {
		err = sl.Error{
			StatusCode: xmlRpcError.HttpStatusCode,
			Exception:  xmlRpcError.Code.(string),
			Message:    xmlRpcError.Err,
		}
	}
	return err
}

func makeXmlRequest(
	retries int, wait time.Duration, client *xmlrpc.Client,
	method string, params []interface{}, pResult interface{}) error {

	err := client.Call(method, params, pResult)

	if xmlRpcError, ok := err.(*xmlrpc.XmlRpcError); ok {
		err = sl.Error{
			StatusCode: xmlRpcError.HttpStatusCode,
			Exception:  xmlRpcError.Code.(string),
			Message:    xmlRpcError.Err,
		}
	}

	if err != nil {
		if !isRetryable(err) {
			return err
		}

		if retries--; retries > 0 {
			jitter := time.Duration(rand.Int63n(int64(wait)))
			wait = wait + jitter/2
			time.Sleep(wait)
			return makeXmlRequest(
				retries, wait, client, method, params, pResult)
		}
	}

	return err
}

func genXMLMask(mask string) interface{} {
	objectMask := map[string]interface{}{}
	for _, item := range strings.Split(mask, ";") {
		if !strings.Contains(item, ".") {
			objectMask[item] = []string{}
			continue
		}

		level := objectMask
		names := strings.Split(item, ".")
		totalNames := len(names)
		for i, name := range names {
			if i == totalNames-1 {
				level[name] = []string{}
				continue
			}

			level[name] = map[string]interface{}{}
			level = level[name].(map[string]interface{})
		}
	}

	return objectMask
}
