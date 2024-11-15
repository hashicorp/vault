// Copyright (c) 2015-2024 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package resty

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Request struct and methods
//_______________________________________________________________________

// Request struct is used to compose and fire individual requests from
// Resty client. The [Request] provides an option to override client-level
// settings and also an option for the request composition.
type Request struct {
	URL           string
	Method        string
	Token         string
	AuthScheme    string
	QueryParam    url.Values
	FormData      url.Values
	PathParams    map[string]string
	RawPathParams map[string]string
	Header        http.Header
	Time          time.Time
	Body          interface{}
	Result        interface{}
	resultCurlCmd *string
	Error         interface{}
	RawRequest    *http.Request
	SRV           *SRVRecord
	UserInfo      *User
	Cookies       []*http.Cookie
	Debug         bool

	// Attempt is to represent the request attempt made during a Resty
	// request execution flow, including retry count.
	Attempt int

	isMultiPart         bool
	isFormData          bool
	setContentLength    bool
	isSaveResponse      bool
	notParseResponse    bool
	jsonEscapeHTML      bool
	trace               bool
	outputFile          string
	fallbackContentType string
	forceContentType    string
	ctx                 context.Context
	values              map[string]interface{}
	client              *Client
	bodyBuf             *bytes.Buffer
	clientTrace         *clientTrace
	log                 Logger
	multipartBoundary   string
	multipartFiles      []*File
	multipartFields     []*MultipartField
	retryConditions     []RetryConditionFunc
	responseBodyLimit   int
	generateCurlOnDebug bool
	unescapeQueryParams bool
}

// GenerateCurlCommand method generates the CURL command for the request.
func (r *Request) GenerateCurlCommand() string {
	if !(r.Debug && r.generateCurlOnDebug) {
		return ""
	}

	if r.resultCurlCmd != nil {
		return *r.resultCurlCmd
	}

	if r.RawRequest == nil {
		r.client.executeBefore(r) // mock with r.Get("/")
	}
	if r.resultCurlCmd == nil {
		r.resultCurlCmd = new(string)
	}
	*r.resultCurlCmd = buildCurlRequest(r.RawRequest, r.client.httpClient.Jar)
	return *r.resultCurlCmd
}

// Context method returns the Context if it is already set in the [Request]
// otherwise, it creates a new one using [context.Background].
func (r *Request) Context() context.Context {
	if r.ctx == nil {
		return context.Background()
	}
	return r.ctx
}

// SetContext method sets the [context.Context] for current [Request]. It allows
// to interrupt the request execution if `ctx.Done()` channel is closed.
// See https://blog.golang.org/context article and the package [context]
// documentation.
func (r *Request) SetContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// SetHeader method sets a single header field and its value in the current request.
//
// For Example: To set `Content-Type` and `Accept` as `application/json`.
//
//	client.R().
//		SetHeader("Content-Type", "application/json").
//		SetHeader("Accept", "application/json")
//
// It overrides the header value set at the client instance level.
func (r *Request) SetHeader(header, value string) *Request {
	r.Header.Set(header, value)
	return r
}

// SetHeaders method sets multiple header fields and their values at one go in the current request.
//
// For Example: To set `Content-Type` and `Accept` as `application/json`
//
//	client.R().
//		SetHeaders(map[string]string{
//			"Content-Type": "application/json",
//			"Accept": "application/json",
//		})
//
// It overrides the header value set at the client instance level.
func (r *Request) SetHeaders(headers map[string]string) *Request {
	for h, v := range headers {
		r.SetHeader(h, v)
	}
	return r
}

// SetHeaderMultiValues sets multiple header fields and their values as a list of strings in the current request.
//
// For Example: To set `Accept` as `text/html, application/xhtml+xml, application/xml;q=0.9, image/webp, */*;q=0.8`
//
//	client.R().
//		SetHeaderMultiValues(map[string][]string{
//			"Accept": []string{"text/html", "application/xhtml+xml", "application/xml;q=0.9", "image/webp", "*/*;q=0.8"},
//		})
//
// It overrides the header value set at the client instance level.
func (r *Request) SetHeaderMultiValues(headers map[string][]string) *Request {
	for key, values := range headers {
		r.SetHeader(key, strings.Join(values, ", "))
	}
	return r
}

// SetHeaderVerbatim method sets a single header field and its value verbatim in the current request.
//
// For Example: To set `all_lowercase` and `UPPERCASE` as `available`.
//
//	client.R().
//		SetHeaderVerbatim("all_lowercase", "available").
//		SetHeaderVerbatim("UPPERCASE", "available")
//
// It overrides the header value set at the client instance level.
func (r *Request) SetHeaderVerbatim(header, value string) *Request {
	r.Header[header] = []string{value}
	return r
}

// SetQueryParam method sets a single parameter and its value in the current request.
// It will be formed as a query string for the request.
//
// For Example: `search=kitchen%20papers&size=large` in the URL after the `?` mark.
//
//	client.R().
//		SetQueryParam("search", "kitchen papers").
//		SetQueryParam("size", "large")
//
// It overrides the query parameter value set at the client instance level.
func (r *Request) SetQueryParam(param, value string) *Request {
	r.QueryParam.Set(param, value)
	return r
}

// SetQueryParams method sets multiple parameters and their values at one go in the current request.
// It will be formed as a query string for the request.
//
// For Example: `search=kitchen%20papers&size=large` in the URL after the `?` mark.
//
//	client.R().
//		SetQueryParams(map[string]string{
//			"search": "kitchen papers",
//			"size": "large",
//		})
//
// It overrides the query parameter value set at the client instance level.
func (r *Request) SetQueryParams(params map[string]string) *Request {
	for p, v := range params {
		r.SetQueryParam(p, v)
	}
	return r
}

// SetUnescapeQueryParams method sets the unescape query parameters choice for request URL.
// To prevent broken URL, resty replaces space (" ") with "+" in the query parameters.
//
// This method overrides the value set by [Client.SetUnescapeQueryParams]
//
// NOTE: Request failure is possible due to non-standard usage of Unescaped Query Parameters.
func (r *Request) SetUnescapeQueryParams(unescape bool) *Request {
	r.unescapeQueryParams = unescape
	return r
}

// SetQueryParamsFromValues method appends multiple parameters with multi-value
// ([url.Values]) at one go in the current request. It will be formed as
// query string for the request.
//
// For Example: `status=pending&status=approved&status=open` in the URL after the `?` mark.
//
//	client.R().
//		SetQueryParamsFromValues(url.Values{
//			"status": []string{"pending", "approved", "open"},
//		})
//
// It overrides the query parameter value set at the client instance level.
func (r *Request) SetQueryParamsFromValues(params url.Values) *Request {
	for p, v := range params {
		for _, pv := range v {
			r.QueryParam.Add(p, pv)
		}
	}
	return r
}

// SetQueryString method provides the ability to use string as an input to set URL query string for the request.
//
//	client.R().
//		SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more")
//
// It overrides the query parameter value set at the client instance level.
func (r *Request) SetQueryString(query string) *Request {
	params, err := url.ParseQuery(strings.TrimSpace(query))
	if err == nil {
		for p, v := range params {
			for _, pv := range v {
				r.QueryParam.Add(p, pv)
			}
		}
	} else {
		r.log.Errorf("%v", err)
	}
	return r
}

// SetFormData method sets Form parameters and their values for the current request.
// It applies only to HTTP methods `POST` and `PUT`, and by default requests
// content type would be set as `application/x-www-form-urlencoded`.
//
//	client.R().
//		SetFormData(map[string]string{
//			"access_token": "BC594900-518B-4F7E-AC75-BD37F019E08F",
//			"user_id": "3455454545",
//		})
//
// It overrides the form data value set at the client instance level.
func (r *Request) SetFormData(data map[string]string) *Request {
	for k, v := range data {
		r.FormData.Set(k, v)
	}
	return r
}

// SetFormDataFromValues method appends multiple form parameters with multi-value
// ([url.Values]) at one go in the current request.
//
//	client.R().
//		SetFormDataFromValues(url.Values{
//			"search_criteria": []string{"book", "glass", "pencil"},
//		})
//
// It overrides the form data value set at the client instance level.
func (r *Request) SetFormDataFromValues(data url.Values) *Request {
	for k, v := range data {
		for _, kv := range v {
			r.FormData.Add(k, kv)
		}
	}
	return r
}

// SetBody method sets the request body for the request. It supports various practical needs as easy.
// It's quite handy and powerful. Supported request body data types are `string`,
// `[]byte`, `struct`, `map`, `slice` and [io.Reader].
//
// Body value can be pointer or non-pointer. Automatic marshalling for JSON and XML content type, if it is `struct`, `map`, or `slice`.
//
// NOTE: [io.Reader] is processed in bufferless mode while sending a request.
//
// For Example:
//
// `struct` gets marshaled based on the request header `Content-Type`.
//
//	client.R().
//		SetBody(User{
//			Username: "jeeva@myjeeva.com",
//			Password: "welcome2resty",
//		})
//
// 'map` gets marshaled based on the request header `Content-Type`.
//
//	client.R().
//		SetBody(map[string]interface{}{
//			"username": "jeeva@myjeeva.com",
//			"password": "welcome2resty",
//			"address": &Address{
//				Address1: "1111 This is my street",
//				Address2: "Apt 201",
//				City: "My City",
//				State: "My State",
//				ZipCode: 00000,
//			},
//		})
//
// `string` as a body input. Suitable for any need as a string input.
//
//	client.R().
//		SetBody(`{
//			"username": "jeeva@getrightcare.com",
//			"password": "admin"
//		}`)
//
// `[]byte` as a body input. Suitable for raw requests such as file upload, serialize & deserialize, etc.
//
//	client.R().
//		SetBody([]byte("This is my raw request, sent as-is"))
//
// and so on.
func (r *Request) SetBody(body interface{}) *Request {
	r.Body = body
	return r
}

// SetResult method is to register the response `Result` object for automatic
// unmarshalling of the HTTP response if the response status code is
// between 200 and 299, and the content type is JSON or XML.
//
// Note: [Request.SetResult] input can be a pointer or non-pointer.
//
// The pointer with handle
//
//	authToken := &AuthToken{}
//	client.R().SetResult(authToken)
//
//	// Can be accessed via -
//	fmt.Println(authToken) OR fmt.Println(response.Result().(*AuthToken))
//
// OR -
//
// The pointer without handle or non-pointer
//
//	client.R().SetResult(&AuthToken{})
//	// OR
//	client.R().SetResult(AuthToken{})
//
//	// Can be accessed via -
//	fmt.Println(response.Result().(*AuthToken))
func (r *Request) SetResult(res interface{}) *Request {
	if res != nil {
		r.Result = getPointer(res)
	}
	return r
}

// SetError method is to register the request `Error` object for automatic unmarshalling for the request,
// if the response status code is greater than 399 and the content type is either JSON or XML.
//
// NOTE: [Request.SetError] input can be a pointer or non-pointer.
//
//	client.R().SetError(&AuthError{})
//	// OR
//	client.R().SetError(AuthError{})
//
// Accessing an error value from response instance.
//
//	response.Error().(*AuthError)
//
// If this request Error object is nil, Resty will use the client-level error object Type if it is set.
func (r *Request) SetError(err interface{}) *Request {
	r.Error = getPointer(err)
	return r
}

// SetFile method sets a single file field name and its path for multipart upload.
//
//	client.R().
//		SetFile("my_file", "/Users/jeeva/Gas Bill - Sep.pdf")
func (r *Request) SetFile(param, filePath string) *Request {
	r.isMultiPart = true
	r.FormData.Set("@"+param, filePath)
	return r
}

// SetFiles method sets multiple file field names and their paths for multipart uploads.
//
//	client.R().
//		SetFiles(map[string]string{
//				"my_file1": "/Users/jeeva/Gas Bill - Sep.pdf",
//				"my_file2": "/Users/jeeva/Electricity Bill - Sep.pdf",
//				"my_file3": "/Users/jeeva/Water Bill - Sep.pdf",
//			})
func (r *Request) SetFiles(files map[string]string) *Request {
	r.isMultiPart = true
	for f, fp := range files {
		r.FormData.Set("@"+f, fp)
	}
	return r
}

// SetFileReader method is to set a file using [io.Reader] for multipart upload.
//
//	client.R().
//		SetFileReader("profile_img", "my-profile-img.png", bytes.NewReader(profileImgBytes)).
//		SetFileReader("notes", "user-notes.txt", bytes.NewReader(notesBytes))
func (r *Request) SetFileReader(param, fileName string, reader io.Reader) *Request {
	r.isMultiPart = true
	r.multipartFiles = append(r.multipartFiles, &File{
		Name:      fileName,
		ParamName: param,
		Reader:    reader,
	})
	return r
}

// SetMultipartFormData method allows simple form data to be attached to the request
// as `multipart:form-data`
func (r *Request) SetMultipartFormData(data map[string]string) *Request {
	for k, v := range data {
		r = r.SetMultipartField(k, "", "", strings.NewReader(v))
	}

	return r
}

// SetMultipartField method sets custom data with Content-Type using [io.Reader] for multipart upload.
func (r *Request) SetMultipartField(param, fileName, contentType string, reader io.Reader) *Request {
	r.isMultiPart = true
	r.multipartFields = append(r.multipartFields, &MultipartField{
		Param:       param,
		FileName:    fileName,
		ContentType: contentType,
		Reader:      reader,
	})
	return r
}

// SetMultipartFields method sets multiple data fields using [io.Reader] for multipart upload.
//
// For Example:
//
//	client.R().SetMultipartFields(
//		&resty.MultipartField{
//			Param:       "uploadManifest1",
//			FileName:    "upload-file-1.json",
//			ContentType: "application/json",
//			Reader:      strings.NewReader(`{"input": {"name": "Uploaded document 1", "_filename" : ["file1.txt"]}}`),
//		},
//		&resty.MultipartField{
//			Param:       "uploadManifest2",
//			FileName:    "upload-file-2.json",
//			ContentType: "application/json",
//			Reader:      strings.NewReader(`{"input": {"name": "Uploaded document 2", "_filename" : ["file2.txt"]}}`),
//		})
//
// If you have a `slice` of fields already, then call-
//
//	client.R().SetMultipartFields(fields...)
func (r *Request) SetMultipartFields(fields ...*MultipartField) *Request {
	r.isMultiPart = true
	r.multipartFields = append(r.multipartFields, fields...)
	return r
}

// SetMultipartBoundary method sets the custom multipart boundary for the multipart request.
// Typically, the `mime/multipart` package generates a random multipart boundary if not provided.
func (r *Request) SetMultipartBoundary(boundary string) *Request {
	r.multipartBoundary = boundary
	return r
}

// SetContentLength method sets the current request's HTTP header `Content-Length` value.
// By default, Resty won't set `Content-Length`.
//
// See [Client.SetContentLength]
//
//	client.R().SetContentLength(true)
//
// It overrides the value set at the client instance level.
func (r *Request) SetContentLength(l bool) *Request {
	r.setContentLength = l
	return r
}

// SetBasicAuth method sets the basic authentication header in the current HTTP request.
//
// For Example:
//
//	Authorization: Basic <base64-encoded-value>
//
// To set the header for username "go-resty" and password "welcome"
//
//	client.R().SetBasicAuth("go-resty", "welcome")
//
// It overrides the credentials set by method [Client.SetBasicAuth].
func (r *Request) SetBasicAuth(username, password string) *Request {
	r.UserInfo = &User{Username: username, Password: password}
	return r
}

// SetAuthToken method sets the auth token header(Default Scheme: Bearer) in the current HTTP request. Header example:
//
//	Authorization: Bearer <auth-token-value-comes-here>
//
// For Example: To set auth token BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F
//
//	client.R().SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")
//
// It overrides the Auth token set by method [Client.SetAuthToken].
func (r *Request) SetAuthToken(token string) *Request {
	r.Token = token
	return r
}

// SetAuthScheme method sets the auth token scheme type in the HTTP request.
//
// Example Header value structure:
//
//	Authorization: <auth-scheme-value-set-here> <auth-token-value>
//
// For Example: To set the scheme to use OAuth
//
//	client.R().SetAuthScheme("OAuth")
//
//	// The outcome will be -
//	Authorization: OAuth <auth-token-value>
//
// Information about Auth schemes can be found in [RFC 7235], IANA [HTTP Auth schemes]
//
// It overrides the `Authorization` scheme set by method [Client.SetAuthScheme].
//
// [RFC 7235]: https://tools.ietf.org/html/rfc7235
// [HTTP Auth schemes]: https://www.iana.org/assignments/http-authschemes/http-authschemes.xhtml#authschemes
func (r *Request) SetAuthScheme(scheme string) *Request {
	r.AuthScheme = scheme
	return r
}

// SetDigestAuth method sets the Digest Access auth scheme for the HTTP request.
// If a server responds with 401 and sends a Digest challenge in the WWW-Authenticate Header,
// the request will be resent with the appropriate Authorization Header.
//
// For Example: To set the Digest scheme with username "Mufasa" and password "Circle Of Life"
//
//	client.R().SetDigestAuth("Mufasa", "Circle Of Life")
//
// Information about Digest Access Authentication can be found in [RFC 7616]
//
// It overrides the digest username and password set by method [Client.SetDigestAuth].
//
// [RFC 7616]: https://datatracker.ietf.org/doc/html/rfc7616
func (r *Request) SetDigestAuth(username, password string) *Request {
	oldTransport := r.client.httpClient.Transport
	r.client.OnBeforeRequest(func(c *Client, _ *Request) error {
		c.httpClient.Transport = &digestTransport{
			digestCredentials: digestCredentials{username, password},
			transport:         oldTransport,
		}
		return nil
	})
	r.client.OnAfterResponse(func(c *Client, _ *Response) error {
		c.httpClient.Transport = oldTransport
		return nil
	})

	return r
}

// SetOutput method sets the output file for the current HTTP request. The current
// HTTP response will be saved in the given file. It is similar to the `curl -o` flag.
//
// Absolute path or relative path can be used.
//
// If it is a relative path, then the output file goes under the output directory, as mentioned
// in the [Client.SetOutputDirectory].
//
//	client.R().
//		SetOutput("/Users/jeeva/Downloads/ReplyWithHeader-v5.1-beta.zip").
//		Get("http://bit.ly/1LouEKr")
//
// NOTE: In this scenario [Response.Body] might be nil.
func (r *Request) SetOutput(file string) *Request {
	r.outputFile = file
	r.isSaveResponse = true
	return r
}

// SetSRV method sets the details to query the service SRV record and execute the
// request.
//
//	client.R().
//		SetSRV(SRVRecord{"web", "testservice.com"}).
//		Get("/get")
func (r *Request) SetSRV(srv *SRVRecord) *Request {
	r.SRV = srv
	return r
}

// SetDoNotParseResponse method instructs Resty not to parse the response body automatically.
// Resty exposes the raw response body as [io.ReadCloser]. If you use it, do not
// forget to close the body, otherwise, you might get into connection leaks, and connection
// reuse may not happen.
//
// NOTE: [Response] middlewares are not executed using this option. You have
// taken over the control of response parsing from Resty.
func (r *Request) SetDoNotParseResponse(parse bool) *Request {
	r.notParseResponse = parse
	return r
}

// SetResponseBodyLimit method sets a maximum body size limit in bytes on response,
// avoid reading too much data to memory.
//
// Client will return [resty.ErrResponseBodyTooLarge] if the body size of the body
// in the uncompressed response is larger than the limit.
// Body size limit will not be enforced in the following cases:
//   - ResponseBodyLimit <= 0, which is the default behavior.
//   - [Request.SetOutput] is called to save response data to the file.
//   - "DoNotParseResponse" is set for client or request.
//
// It overrides the value set at the client instance level. see [Client.SetResponseBodyLimit]
func (r *Request) SetResponseBodyLimit(v int) *Request {
	r.responseBodyLimit = v
	return r
}

// SetPathParam method sets a single URL path key-value pair in the
// Resty current request instance.
//
//	client.R().SetPathParam("userId", "sample@sample.com")
//
//	Result:
//	   URL - /v1/users/{userId}/details
//	   Composed URL - /v1/users/sample@sample.com/details
//
//	client.R().SetPathParam("path", "groups/developers")
//
//	Result:
//	   URL - /v1/users/{userId}/details
//	   Composed URL - /v1/users/groups%2Fdevelopers/details
//
// It replaces the value of the key while composing the request URL.
// The values will be escaped using function [url.PathEscape].
//
// It overrides the path parameter set at the client instance level.
func (r *Request) SetPathParam(param, value string) *Request {
	r.PathParams[param] = value
	return r
}

// SetPathParams method sets multiple URL path key-value pairs at one go in the
// Resty current request instance.
//
//	client.R().SetPathParams(map[string]string{
//		"userId":       "sample@sample.com",
//		"subAccountId": "100002",
//		"path":         "groups/developers",
//	})
//
//	Result:
//	   URL - /v1/users/{userId}/{subAccountId}/{path}/details
//	   Composed URL - /v1/users/sample@sample.com/100002/groups%2Fdevelopers/details
//
// It replaces the value of the key while composing the request URL.
// The values will be escaped using function [url.PathEscape].
//
// It overrides the path parameter set at the client instance level.
func (r *Request) SetPathParams(params map[string]string) *Request {
	for p, v := range params {
		r.SetPathParam(p, v)
	}
	return r
}

// SetRawPathParam method sets a single URL path key-value pair in the
// Resty current request instance.
//
//	client.R().SetPathParam("userId", "sample@sample.com")
//
//	Result:
//	   URL - /v1/users/{userId}/details
//	   Composed URL - /v1/users/sample@sample.com/details
//
//	client.R().SetPathParam("path", "groups/developers")
//
//	Result:
//	   URL - /v1/users/{userId}/details
//	   Composed URL - /v1/users/groups/developers/details
//
// It replaces the value of the key while composing the request URL.
// The value will be used as-is and has not been escaped.
//
// It overrides the raw path parameter set at the client instance level.
func (r *Request) SetRawPathParam(param, value string) *Request {
	r.RawPathParams[param] = value
	return r
}

// SetRawPathParams method sets multiple URL path key-value pairs at one go in the
// Resty current request instance.
//
//	client.R().SetPathParams(map[string]string{
//		"userId": "sample@sample.com",
//		"subAccountId": "100002",
//		"path":         "groups/developers",
//	})
//
//	Result:
//	   URL - /v1/users/{userId}/{subAccountId}/{path}/details
//	   Composed URL - /v1/users/sample@sample.com/100002/groups/developers/details
//
// It replaces the value of the key while composing the request URL.
// The value will be used as-is and has not been escaped.
//
// It overrides the raw path parameter set at the client instance level.
func (r *Request) SetRawPathParams(params map[string]string) *Request {
	for p, v := range params {
		r.SetRawPathParam(p, v)
	}
	return r
}

// ExpectContentType method allows to provide fallback `Content-Type` for automatic unmarshalling
// when the `Content-Type` response header is unavailable.
func (r *Request) ExpectContentType(contentType string) *Request {
	r.fallbackContentType = contentType
	return r
}

// ForceContentType method provides a strong sense of response `Content-Type` for
// automatic unmarshalling. Resty gives this a higher priority than the `Content-Type`
// response header.
//
// This means that if both [Request.ForceContentType] is set and
// the response `Content-Type` is available, `ForceContentType` will win.
func (r *Request) ForceContentType(contentType string) *Request {
	r.forceContentType = contentType
	return r
}

// SetJSONEscapeHTML method enables or disables the HTML escape on JSON marshal.
// By default, escape HTML is false.
//
// NOTE: This option only applies to the standard JSON Marshaller used by Resty.
//
// It overrides the value set at the client instance level, see [Client.SetJSONEscapeHTML]
func (r *Request) SetJSONEscapeHTML(b bool) *Request {
	r.jsonEscapeHTML = b
	return r
}

// SetCookie method appends a single cookie in the current request instance.
//
//	client.R().SetCookie(&http.Cookie{
//				Name:"go-resty",
//				Value:"This is cookie value",
//			})
//
// NOTE: Method appends the Cookie value into existing Cookie even if its already existing.
func (r *Request) SetCookie(hc *http.Cookie) *Request {
	r.Cookies = append(r.Cookies, hc)
	return r
}

// SetCookies method sets an array of cookies in the current request instance.
//
//	cookies := []*http.Cookie{
//		&http.Cookie{
//			Name:"go-resty-1",
//			Value:"This is cookie 1 value",
//		},
//		&http.Cookie{
//			Name:"go-resty-2",
//			Value:"This is cookie 2 value",
//		},
//	}
//
//	// Setting a cookies into resty's current request
//	client.R().SetCookies(cookies)
//
// NOTE: Method appends the Cookie value into existing Cookie even if its already existing.
func (r *Request) SetCookies(rs []*http.Cookie) *Request {
	r.Cookies = append(r.Cookies, rs...)
	return r
}

// SetLogger method sets given writer for logging Resty request and response details.
// By default, requests and responses inherit their logger from the client.
//
// Compliant to interface [resty.Logger].
//
// It overrides the logger value set at the client instance level.
func (r *Request) SetLogger(l Logger) *Request {
	r.log = l
	return r
}

// SetDebug method enables the debug mode on the current request. It logs
// the details current request and response.
//
//	client.SetDebug(true)
//
// Also, it can be enabled at the request level for a particular request; see [Request.SetDebug].
//   - For [Request], it logs information such as HTTP verb, Relative URL path,
//     Host, Headers, and Body if it has one.
//   - For [Response], it logs information such as Status, Response Time, Headers,
//     and Body if it has one.
func (r *Request) SetDebug(d bool) *Request {
	r.Debug = d
	return r
}

// AddRetryCondition method adds a retry condition function to the request's
// array of functions is checked to determine if the request can be retried.
// The request will retry if any functions return true and the error is nil.
//
// NOTE: The request level retry conditions are checked before all retry
// conditions from the client instance.
func (r *Request) AddRetryCondition(condition RetryConditionFunc) *Request {
	r.retryConditions = append(r.retryConditions, condition)
	return r
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// HTTP request tracing
//_______________________________________________________________________

// EnableTrace method enables trace for the current request
// using [httptrace.ClientTrace] and provides insights.
//
//	client := resty.New()
//
//	resp, err := client.R().EnableTrace().Get("https://httpbin.org/get")
//	fmt.Println("Error:", err)
//	fmt.Println("Trace Info:", resp.Request.TraceInfo())
//
// See [Client.EnableTrace] is also available to get trace info for all requests.
func (r *Request) EnableTrace() *Request {
	r.trace = true
	return r
}

// EnableGenerateCurlOnDebug method enables the generation of CURL commands in the debug log.
// It works in conjunction with debug mode. It overrides the options set by the [Client].
//
// NOTE: Use with care.
//   - Potential to leak sensitive data from [Request] and [Response] in the debug log.
//   - Beware of memory usage since the request body is reread.
func (r *Request) EnableGenerateCurlOnDebug() *Request {
	r.generateCurlOnDebug = true
	return r
}

// DisableGenerateCurlOnDebug method disables the option set by [Request.EnableGenerateCurlOnDebug].
// It overrides the options set by the [Client].
func (r *Request) DisableGenerateCurlOnDebug() *Request {
	r.generateCurlOnDebug = false
	return r
}

// TraceInfo method returns the trace info for the request.
// If either the [Client.EnableTrace] or [Request.EnableTrace] function has not been called
// before the request is made, an empty [resty.TraceInfo] object is returned.
func (r *Request) TraceInfo() TraceInfo {
	ct := r.clientTrace

	if ct == nil {
		return TraceInfo{}
	}

	ti := TraceInfo{
		DNSLookup:      ct.dnsDone.Sub(ct.dnsStart),
		TLSHandshake:   ct.tlsHandshakeDone.Sub(ct.tlsHandshakeStart),
		ServerTime:     ct.gotFirstResponseByte.Sub(ct.gotConn),
		IsConnReused:   ct.gotConnInfo.Reused,
		IsConnWasIdle:  ct.gotConnInfo.WasIdle,
		ConnIdleTime:   ct.gotConnInfo.IdleTime,
		RequestAttempt: r.Attempt,
	}

	// Calculate the total time accordingly,
	// when connection is reused
	if ct.gotConnInfo.Reused {
		ti.TotalTime = ct.endTime.Sub(ct.getConn)
	} else {
		ti.TotalTime = ct.endTime.Sub(ct.dnsStart)
	}

	// Only calculate on successful connections
	if !ct.connectDone.IsZero() {
		ti.TCPConnTime = ct.connectDone.Sub(ct.dnsDone)
	}

	// Only calculate on successful connections
	if !ct.gotConn.IsZero() {
		ti.ConnTime = ct.gotConn.Sub(ct.getConn)
	}

	// Only calculate on successful connections
	if !ct.gotFirstResponseByte.IsZero() {
		ti.ResponseTime = ct.endTime.Sub(ct.gotFirstResponseByte)
	}

	// Capture remote address info when connection is non-nil
	if ct.gotConnInfo.Conn != nil {
		ti.RemoteAddr = ct.gotConnInfo.Conn.RemoteAddr()
	}

	return ti
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// HTTP verb method starts here
//_______________________________________________________________________

// Get method does GET HTTP request. It's defined in section 4.3.1 of RFC7231.
func (r *Request) Get(url string) (*Response, error) {
	return r.Execute(MethodGet, url)
}

// Head method does HEAD HTTP request. It's defined in section 4.3.2 of RFC7231.
func (r *Request) Head(url string) (*Response, error) {
	return r.Execute(MethodHead, url)
}

// Post method does POST HTTP request. It's defined in section 4.3.3 of RFC7231.
func (r *Request) Post(url string) (*Response, error) {
	return r.Execute(MethodPost, url)
}

// Put method does PUT HTTP request. It's defined in section 4.3.4 of RFC7231.
func (r *Request) Put(url string) (*Response, error) {
	return r.Execute(MethodPut, url)
}

// Delete method does DELETE HTTP request. It's defined in section 4.3.5 of RFC7231.
func (r *Request) Delete(url string) (*Response, error) {
	return r.Execute(MethodDelete, url)
}

// Options method does OPTIONS HTTP request. It's defined in section 4.3.7 of RFC7231.
func (r *Request) Options(url string) (*Response, error) {
	return r.Execute(MethodOptions, url)
}

// Patch method does PATCH HTTP request. It's defined in section 2 of RFC5789.
func (r *Request) Patch(url string) (*Response, error) {
	return r.Execute(MethodPatch, url)
}

// Send method performs the HTTP request using the method and URL already defined
// for current [Request].
//
//	req := client.R()
//	req.Method = resty.MethodGet
//	req.URL = "http://httpbin.org/get"
//	resp, err := req.Send()
func (r *Request) Send() (*Response, error) {
	return r.Execute(r.Method, r.URL)
}

// Execute method performs the HTTP request with the given HTTP method and URL
// for current [Request].
//
//	resp, err := client.R().Execute(resty.MethodGet, "http://httpbin.org/get")
func (r *Request) Execute(method, url string) (*Response, error) {
	var addrs []*net.SRV
	var resp *Response
	var err error

	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				r.client.onPanicHooks(r, err)
			} else {
				r.client.onPanicHooks(r, fmt.Errorf("panic %v", rec))
			}
			panic(rec)
		}
	}()

	if r.isMultiPart && !(method == MethodPost || method == MethodPut || method == MethodPatch) {
		// No OnError hook here since this is a request validation error
		err := fmt.Errorf("multipart content is not allowed in HTTP verb [%v]", method)
		r.client.onInvalidHooks(r, err)
		return nil, err
	}

	if r.SRV != nil {
		_, addrs, err = net.LookupSRV(r.SRV.Service, "tcp", r.SRV.Domain)
		if err != nil {
			r.client.onErrorHooks(r, nil, err)
			return nil, err
		}
	}

	r.Method = method
	r.URL = r.selectAddr(addrs, url, 0)

	if r.client.RetryCount == 0 {
		r.Attempt = 1
		resp, err = r.client.execute(r)
		r.client.onErrorHooks(r, resp, unwrapNoRetryErr(err))
		return resp, unwrapNoRetryErr(err)
	}

	err = Backoff(
		func() (*Response, error) {
			r.Attempt++

			r.URL = r.selectAddr(addrs, url, r.Attempt)

			resp, err = r.client.execute(r)
			if err != nil {
				r.log.Warnf("%v, Attempt %v", err, r.Attempt)
			}

			return resp, err
		},
		Retries(r.client.RetryCount),
		WaitTime(r.client.RetryWaitTime),
		MaxWaitTime(r.client.RetryMaxWaitTime),
		RetryConditions(append(r.retryConditions, r.client.RetryConditions...)),
		RetryHooks(r.client.RetryHooks),
		ResetMultipartReaders(r.client.RetryResetReaders),
	)

	if err != nil {
		r.log.Errorf("%v", err)
	}

	r.client.onErrorHooks(r, resp, unwrapNoRetryErr(err))
	return resp, unwrapNoRetryErr(err)
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// SRVRecord struct
//_______________________________________________________________________

// SRVRecord struct holds the data to query the SRV record for the
// following service.
type SRVRecord struct {
	Service string
	Domain  string
}

func (r *Request) fmtBodyString(sl int64) (body string) {
	body = "***** NO CONTENT *****"
	if !isPayloadSupported(r.Method, r.client.AllowGetMethodPayload) {
		return
	}

	if _, ok := r.Body.(io.Reader); ok {
		body = "***** BODY IS io.Reader *****"
		return
	}

	// multipart or form-data
	if r.isMultiPart || r.isFormData {
		bodySize := int64(r.bodyBuf.Len())
		if bodySize > sl {
			body = fmt.Sprintf("***** REQUEST TOO LARGE (size - %d) *****", bodySize)
			return
		}
		body = r.bodyBuf.String()
		return
	}

	// request body data
	if r.Body == nil {
		return
	}
	var prtBodyBytes []byte
	var err error

	contentType := r.Header.Get(hdrContentTypeKey)
	kind := kindOf(r.Body)
	if canJSONMarshal(contentType, kind) {
		var bodyBuf *bytes.Buffer
		bodyBuf, err = noescapeJSONMarshalIndent(&r.Body)
		if err == nil {
			prtBodyBytes = bodyBuf.Bytes()
			defer releaseBuffer(bodyBuf)
		}
	} else if IsXMLType(contentType) && (kind == reflect.Struct) {
		prtBodyBytes, err = xml.MarshalIndent(&r.Body, "", "   ")
	} else if b, ok := r.Body.(string); ok {
		if IsJSONType(contentType) {
			bodyBytes := []byte(b)
			out := acquireBuffer()
			defer releaseBuffer(out)
			if err = json.Indent(out, bodyBytes, "", "   "); err == nil {
				prtBodyBytes = out.Bytes()
			}
		} else {
			body = b
		}
	} else if b, ok := r.Body.([]byte); ok {
		body = fmt.Sprintf("***** BODY IS byte(s) (size - %d) *****", len(b))
		return
	}

	if prtBodyBytes != nil && err == nil {
		body = string(prtBodyBytes)
	}

	if len(body) > 0 {
		bodySize := int64(len([]byte(body)))
		if bodySize > sl {
			body = fmt.Sprintf("***** REQUEST TOO LARGE (size - %d) *****", bodySize)
		}
	}

	return
}

func (r *Request) selectAddr(addrs []*net.SRV, path string, attempt int) string {
	if addrs == nil {
		return path
	}

	idx := attempt % len(addrs)
	domain := strings.TrimRight(addrs[idx].Target, ".")
	path = strings.TrimLeft(path, "/")

	return fmt.Sprintf("%s://%s:%d/%s", r.client.scheme, domain, addrs[idx].Port, path)
}

func (r *Request) initValuesMap() {
	if r.values == nil {
		r.values = make(map[string]interface{})
	}
}

var noescapeJSONMarshal = func(v interface{}) (*bytes.Buffer, error) {
	buf := acquireBuffer()
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		releaseBuffer(buf)
		return nil, err
	}

	return buf, nil
}

var noescapeJSONMarshalIndent = func(v interface{}) (*bytes.Buffer, error) {
	buf := acquireBuffer()
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "   ")

	if err := encoder.Encode(v); err != nil {
		releaseBuffer(buf)
		return nil, err
	}

	return buf, nil
}
