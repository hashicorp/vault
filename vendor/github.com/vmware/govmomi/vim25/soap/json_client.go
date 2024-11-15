/*
Copyright (c) 2023-2023 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package soap

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"reflect"
	"strings"

	"github.com/vmware/govmomi/vim25/xml"

	"github.com/vmware/govmomi/vim25/types"
)

const (
	sessionHeader = "vmware-api-session-id"
)

var (
	// errInvalidResponse is used during unmarshaling when the response content
	// does not match expectations e.g. unexpected HTTP status code or MIME
	// type.
	errInvalidResponse error = errors.New("Invalid response")
	// errInputError is used as root error when the request is malformed.
	errInputError error = errors.New("Invalid input error")
)

// Handles round trip using json HTTP
func (c *Client) jsonRoundTrip(ctx context.Context, req, res HasFault) error {
	this, method, params, err := unpackSOAPRequest(req)
	if err != nil {
		return fmt.Errorf("Cannot unpack the request. %w", err)
	}

	return c.invoke(ctx, this, method, params, res)
}

// Invoke calls a managed object method
func (c *Client) invoke(ctx context.Context, this types.ManagedObjectReference, method string, params interface{}, res HasFault) error {
	buffer := bytes.Buffer{}
	if params != nil {
		marshaller := types.NewJSONEncoder(&buffer)
		err := marshaller.Encode(params)
		if err != nil {
			return fmt.Errorf("Encoding request to JSON failed. %w", err)
		}
	}
	uri := c.getPathForName(this, method)
	req, err := http.NewRequest(http.MethodPost, uri, &buffer)
	if err != nil {
		return err
	}

	if len(c.cookie) != 0 {
		req.Header.Add(sessionHeader, c.cookie)
	}

	result, err := getSOAPResultPtr(res)
	if err != nil {
		return fmt.Errorf("Cannot get pointer to the result structure. %w", err)
	}

	return c.Do(ctx, req, c.responseUnmarshaler(&result))
}

// responseUnmarshaler create unmarshaler function for VMOMI JSON request. The
// unmarshaler checks for errors and tries to load the response body in the
// result structure. It is assumed that result is pointer to a data structure or
// interface{}.
func (c *Client) responseUnmarshaler(result interface{}) func(resp *http.Response) error {
	return func(resp *http.Response) error {
		if resp.StatusCode == http.StatusNoContent ||
			(!isError(resp.StatusCode) && resp.ContentLength == 0) {
			return nil
		}

		if e := checkJSONContentType(resp); e != nil {
			return e
		}

		if resp.StatusCode == 500 {
			bodyBytes, e := io.ReadAll(resp.Body)
			if e != nil {
				return e
			}
			var serverErr interface{}
			dec := types.NewJSONDecoder(bytes.NewReader(bodyBytes))
			e = dec.Decode(&serverErr)
			if e != nil {
				return e
			}
			var faultStringStruct struct {
				FaultString string `json:"faultstring,omitempty"`
			}
			dec = types.NewJSONDecoder(bytes.NewReader(bodyBytes))
			e = dec.Decode(&faultStringStruct)
			if e != nil {
				return e
			}

			f := &Fault{
				XMLName: xml.Name{
					Space: c.Namespace,
					Local: reflect.TypeOf(serverErr).Name() + "Fault",
				},
				String: faultStringStruct.FaultString,
				Code:   "ServerFaultCode",
			}
			f.Detail.Fault = serverErr
			return WrapSoapFault(f)
		}

		if isError(resp.StatusCode) {
			return fmt.Errorf("Unexpected HTTP error code: %v. %w", resp.StatusCode, errInvalidResponse)
		}

		dec := types.NewJSONDecoder(resp.Body)
		e := dec.Decode(result)
		if e != nil {
			return e
		}

		c.checkForSessionHeader(resp)

		return nil
	}
}

func isError(statusCode int) bool {
	return statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices
}

// checkForSessionHeader checks if we have new session id.
// This is a hack that intercepts the session id header and then repeats it.
// It is very similar to cookie store but only for the special vCenter
// session header.
func (c *Client) checkForSessionHeader(resp *http.Response) {
	sessionKey := resp.Header.Get(sessionHeader)
	if len(sessionKey) > 0 {
		c.cookie = sessionKey
	}
}

// Checks if the payload of an HTTP response has the JSON MIME type.
func checkJSONContentType(resp *http.Response) error {
	contentType := resp.Header.Get("content-type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("error parsing content-type: %v, error %w", contentType, err)
	}
	if mediaType != "application/json" {
		return fmt.Errorf("content-type is not application/json: %v. %w", contentType, errInvalidResponse)
	}
	return nil
}

func (c *Client) getPathForName(this types.ManagedObjectReference, name string) string {
	const urnPrefix = "urn:"
	ns := c.Namespace
	if strings.HasPrefix(ns, urnPrefix) {
		ns = ns[len(urnPrefix):]
	}
	return fmt.Sprintf("%v/%v/%v/%v/%v/%v", c.u, ns, c.Version, this.Type, this.Value, name)
}

// unpackSOAPRequest converts SOAP request into this value, method nam and
// parameters using reflection. The input is a one of the *Body structures
// defined in methods.go. It is expected to have "Req" field that is a non-null
// pointer to a struct. The struct simple type name is the method name. The
// struct "This" member is the this MoRef value.
func unpackSOAPRequest(req HasFault) (this types.ManagedObjectReference, method string, params interface{}, err error) {
	reqBodyPtr := reflect.ValueOf(req)
	if reqBodyPtr.Kind() != reflect.Ptr {
		err = fmt.Errorf("Expected pointer to request body as input. %w", errInputError)
		return
	}
	reqBody := reqBodyPtr.Elem()
	if reqBody.Kind() != reflect.Struct {
		err = fmt.Errorf("Expected Request body to be structure. %w", errInputError)
		return
	}
	methodRequestPtr := reqBody.FieldByName("Req")
	if methodRequestPtr.Kind() != reflect.Ptr {
		err = fmt.Errorf("Expected method request body field to be pointer to struct. %w", errInputError)
		return
	}
	methodRequest := methodRequestPtr.Elem()
	if methodRequest.Kind() != reflect.Struct {
		err = fmt.Errorf("Expected method request body to be structure. %w", errInputError)
		return
	}
	thisValue := methodRequest.FieldByName("This")
	if thisValue.Kind() != reflect.Struct {
		err = fmt.Errorf("Expected This field in the method request body to be structure. %w", errInputError)
		return
	}
	var ok bool
	if this, ok = thisValue.Interface().(types.ManagedObjectReference); !ok {
		err = fmt.Errorf("Expected This field to be MoRef. %w", errInputError)
		return
	}
	method = methodRequest.Type().Name()
	params = methodRequestPtr.Interface()

	return

}

// getSOAPResultPtr extract a pointer to the result data structure using go
// reflection from a SOAP data structure used for marshalling.
func getSOAPResultPtr(result HasFault) (res interface{}, err error) {
	resBodyPtr := reflect.ValueOf(result)
	if resBodyPtr.Kind() != reflect.Ptr {
		err = fmt.Errorf("Expected pointer to result body as input. %w", errInputError)
		return
	}
	resBody := resBodyPtr.Elem()
	if resBody.Kind() != reflect.Struct {
		err = fmt.Errorf("Expected result body to be structure. %w", errInputError)
		return
	}
	methodResponsePtr := resBody.FieldByName("Res")
	if methodResponsePtr.Kind() != reflect.Ptr {
		err = fmt.Errorf("Expected method response body field to be pointer to struct. %w", errInputError)
		return
	}
	if methodResponsePtr.IsNil() {
		methodResponsePtr.Set(reflect.New(methodResponsePtr.Type().Elem()))
	}
	methodResponse := methodResponsePtr.Elem()
	if methodResponse.Kind() != reflect.Struct {
		err = fmt.Errorf("Expected method response body to be structure. %w", errInputError)
		return
	}
	returnval := methodResponse.FieldByName("Returnval")
	if !returnval.IsValid() {
		// void method and we return nil, nil
		return
	}
	res = returnval.Addr().Interface()
	return
}
