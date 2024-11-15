// Copyright (c) 2015-2024 Jeevanandam M (jeeva@myjeeva.com), All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package resty

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const debugRequestLogKey = "__restyDebugRequestLog"

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Request Middleware(s)
//_______________________________________________________________________

func parseRequestURL(c *Client, r *Request) error {
	if l := len(c.PathParams) + len(c.RawPathParams) + len(r.PathParams) + len(r.RawPathParams); l > 0 {
		params := make(map[string]string, l)

		// GitHub #103 Path Params
		for p, v := range r.PathParams {
			params[p] = url.PathEscape(v)
		}
		for p, v := range c.PathParams {
			if _, ok := params[p]; !ok {
				params[p] = url.PathEscape(v)
			}
		}

		// GitHub #663 Raw Path Params
		for p, v := range r.RawPathParams {
			if _, ok := params[p]; !ok {
				params[p] = v
			}
		}
		for p, v := range c.RawPathParams {
			if _, ok := params[p]; !ok {
				params[p] = v
			}
		}

		if len(params) > 0 {
			var prev int
			buf := acquireBuffer()
			defer releaseBuffer(buf)
			// search for the next or first opened curly bracket
			for curr := strings.Index(r.URL, "{"); curr == 0 || curr > prev; curr = prev + strings.Index(r.URL[prev:], "{") {
				// write everything from the previous position up to the current
				if curr > prev {
					buf.WriteString(r.URL[prev:curr])
				}
				// search for the closed curly bracket from current position
				next := curr + strings.Index(r.URL[curr:], "}")
				// if not found, then write the remainder and exit
				if next < curr {
					buf.WriteString(r.URL[curr:])
					prev = len(r.URL)
					break
				}
				// special case for {}, without parameter's name
				if next == curr+1 {
					buf.WriteString("{}")
				} else {
					// check for the replacement
					key := r.URL[curr+1 : next]
					value, ok := params[key]
					/// keep the original string if the replacement not found
					if !ok {
						value = r.URL[curr : next+1]
					}
					buf.WriteString(value)
				}

				// set the previous position after the closed curly bracket
				prev = next + 1
				if prev >= len(r.URL) {
					break
				}
			}
			if buf.Len() > 0 {
				// write remainder
				if prev < len(r.URL) {
					buf.WriteString(r.URL[prev:])
				}
				r.URL = buf.String()
			}
		}
	}

	// Parsing request URL
	reqURL, err := url.Parse(r.URL)
	if err != nil {
		return err
	}

	// If Request.URL is relative path then added c.HostURL into
	// the request URL otherwise Request.URL will be used as-is
	if !reqURL.IsAbs() {
		r.URL = reqURL.String()
		if len(r.URL) > 0 && r.URL[0] != '/' {
			r.URL = "/" + r.URL
		}

		// TODO: change to use c.BaseURL only in v3.0.0
		baseURL := c.BaseURL
		if len(baseURL) == 0 {
			baseURL = c.HostURL
		}
		reqURL, err = url.Parse(baseURL + r.URL)
		if err != nil {
			return err
		}
	}

	// GH #407 && #318
	if reqURL.Scheme == "" && len(c.scheme) > 0 {
		reqURL.Scheme = c.scheme
	}

	// Adding Query Param
	if len(c.QueryParam)+len(r.QueryParam) > 0 {
		for k, v := range c.QueryParam {
			// skip query parameter if it was set in request
			if _, ok := r.QueryParam[k]; ok {
				continue
			}

			r.QueryParam[k] = v[:]
		}

		// GitHub #123 Preserve query string order partially.
		// Since not feasible in `SetQuery*` resty methods, because
		// standard package `url.Encode(...)` sorts the query params
		// alphabetically
		if len(r.QueryParam) > 0 {
			if IsStringEmpty(reqURL.RawQuery) {
				reqURL.RawQuery = r.QueryParam.Encode()
			} else {
				reqURL.RawQuery = reqURL.RawQuery + "&" + r.QueryParam.Encode()
			}
		}
	}

	// GH#797 Unescape query parameters
	if r.unescapeQueryParams && len(reqURL.RawQuery) > 0 {
		// at this point, all errors caught up in the above operations
		// so ignore the return error on query unescape; I realized
		// while writing the unit test
		unescapedQuery, _ := url.QueryUnescape(reqURL.RawQuery)
		reqURL.RawQuery = strings.ReplaceAll(unescapedQuery, " ", "+") // otherwise request becomes bad request
	}

	r.URL = reqURL.String()

	return nil
}

func parseRequestHeader(c *Client, r *Request) error {
	for k, v := range c.Header {
		if _, ok := r.Header[k]; ok {
			continue
		}
		r.Header[k] = v[:]
	}

	if IsStringEmpty(r.Header.Get(hdrUserAgentKey)) {
		r.Header.Set(hdrUserAgentKey, hdrUserAgentValue)
	}

	if ct := r.Header.Get(hdrContentTypeKey); IsStringEmpty(r.Header.Get(hdrAcceptKey)) && !IsStringEmpty(ct) && (IsJSONType(ct) || IsXMLType(ct)) {
		r.Header.Set(hdrAcceptKey, r.Header.Get(hdrContentTypeKey))
	}

	return nil
}

func parseRequestBody(c *Client, r *Request) error {
	if isPayloadSupported(r.Method, c.AllowGetMethodPayload) {
		switch {
		case r.isMultiPart: // Handling Multipart
			if err := handleMultipart(c, r); err != nil {
				return err
			}
		case len(c.FormData) > 0 || len(r.FormData) > 0: // Handling Form Data
			handleFormData(c, r)
		case r.Body != nil: // Handling Request body
			handleContentType(c, r)

			if err := handleRequestBody(c, r); err != nil {
				return err
			}
		}
	}

	// by default resty won't set content length, you can if you want to :)
	if c.setContentLength || r.setContentLength {
		if r.bodyBuf == nil {
			r.Header.Set(hdrContentLengthKey, "0")
		} else {
			r.Header.Set(hdrContentLengthKey, strconv.Itoa(r.bodyBuf.Len()))
		}
	}

	return nil
}

func createHTTPRequest(c *Client, r *Request) (err error) {
	if r.bodyBuf == nil {
		if reader, ok := r.Body.(io.Reader); ok && isPayloadSupported(r.Method, c.AllowGetMethodPayload) {
			r.RawRequest, err = http.NewRequest(r.Method, r.URL, reader)
		} else if c.setContentLength || r.setContentLength {
			r.RawRequest, err = http.NewRequest(r.Method, r.URL, http.NoBody)
		} else {
			r.RawRequest, err = http.NewRequest(r.Method, r.URL, nil)
		}
	} else {
		// fix data race: must deep copy.
		bodyBuf := bytes.NewBuffer(append([]byte{}, r.bodyBuf.Bytes()...))
		r.RawRequest, err = http.NewRequest(r.Method, r.URL, bodyBuf)
	}

	if err != nil {
		return
	}

	// Assign close connection option
	r.RawRequest.Close = c.closeConnection

	// Add headers into http request
	r.RawRequest.Header = r.Header

	// Add cookies from client instance into http request
	for _, cookie := range c.Cookies {
		r.RawRequest.AddCookie(cookie)
	}

	// Add cookies from request instance into http request
	for _, cookie := range r.Cookies {
		r.RawRequest.AddCookie(cookie)
	}

	// Enable trace
	if c.trace || r.trace {
		r.clientTrace = &clientTrace{}
		r.ctx = r.clientTrace.createContext(r.Context())
	}

	// Use context if it was specified
	if r.ctx != nil {
		r.RawRequest = r.RawRequest.WithContext(r.ctx)
	}

	// assign get body func for the underlying raw request instance
	if r.RawRequest.GetBody == nil {
		bodyCopy, err := getBodyCopy(r)
		if err != nil {
			return err
		}
		if bodyCopy != nil {
			buf := bodyCopy.Bytes()
			r.RawRequest.GetBody = func() (io.ReadCloser, error) {
				b := bytes.NewReader(buf)
				return io.NopCloser(b), nil
			}
		}
	}

	return
}

func addCredentials(c *Client, r *Request) error {
	var isBasicAuth bool
	// Basic Auth
	if r.UserInfo != nil { // takes precedence
		r.RawRequest.SetBasicAuth(r.UserInfo.Username, r.UserInfo.Password)
		isBasicAuth = true
	} else if c.UserInfo != nil {
		r.RawRequest.SetBasicAuth(c.UserInfo.Username, c.UserInfo.Password)
		isBasicAuth = true
	}

	if !c.DisableWarn {
		if isBasicAuth && !strings.HasPrefix(r.URL, "https") {
			r.log.Warnf("Using Basic Auth in HTTP mode is not secure, use HTTPS")
		}
	}

	// Set the Authorization Header Scheme
	var authScheme string
	if !IsStringEmpty(r.AuthScheme) {
		authScheme = r.AuthScheme
	} else if !IsStringEmpty(c.AuthScheme) {
		authScheme = c.AuthScheme
	} else {
		authScheme = "Bearer"
	}

	// Build the Token Auth header
	if !IsStringEmpty(r.Token) { // takes precedence
		r.RawRequest.Header.Set(c.HeaderAuthorizationKey, authScheme+" "+r.Token)
	} else if !IsStringEmpty(c.Token) {
		r.RawRequest.Header.Set(c.HeaderAuthorizationKey, authScheme+" "+c.Token)
	}

	return nil
}

func createCurlCmd(c *Client, r *Request) (err error) {
	if r.Debug && r.generateCurlOnDebug {
		if r.resultCurlCmd == nil {
			r.resultCurlCmd = new(string)
		}
		*r.resultCurlCmd = buildCurlRequest(r.RawRequest, c.httpClient.Jar)
	}
	return nil
}

func requestLogger(c *Client, r *Request) error {
	if r.Debug {
		rr := r.RawRequest
		rh := copyHeaders(rr.Header)
		if c.GetClient().Jar != nil {
			for _, cookie := range c.GetClient().Jar.Cookies(r.RawRequest.URL) {
				s := fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
				if c := rh.Get("Cookie"); c != "" {
					rh.Set("Cookie", c+"; "+s)
				} else {
					rh.Set("Cookie", s)
				}
			}
		}
		rl := &RequestLog{Header: rh, Body: r.fmtBodyString(c.debugBodySizeLimit)}
		if c.requestLog != nil {
			if err := c.requestLog(rl); err != nil {
				return err
			}
		}

		reqLog := "\n==============================================================================\n"

		if r.Debug && r.generateCurlOnDebug {
			reqLog += "~~~ REQUEST(CURL) ~~~\n" +
				fmt.Sprintf("	%v\n", *r.resultCurlCmd)
		}

		reqLog += "~~~ REQUEST ~~~\n" +
			fmt.Sprintf("%s  %s  %s\n", r.Method, rr.URL.RequestURI(), rr.Proto) +
			fmt.Sprintf("HOST   : %s\n", rr.URL.Host) +
			fmt.Sprintf("HEADERS:\n%s\n", composeHeaders(c, r, rl.Header)) +
			fmt.Sprintf("BODY   :\n%v\n", rl.Body) +
			"------------------------------------------------------------------------------\n"

		r.initValuesMap()
		r.values[debugRequestLogKey] = reqLog
	}

	return nil
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Response Middleware(s)
//_______________________________________________________________________

func responseLogger(c *Client, res *Response) error {
	if res.Request.Debug {
		rl := &ResponseLog{Header: copyHeaders(res.Header()), Body: res.fmtBodyString(c.debugBodySizeLimit)}
		if c.responseLog != nil {
			if err := c.responseLog(rl); err != nil {
				return err
			}
		}

		debugLog := res.Request.values[debugRequestLogKey].(string)
		debugLog += "~~~ RESPONSE ~~~\n" +
			fmt.Sprintf("STATUS       : %s\n", res.Status()) +
			fmt.Sprintf("PROTO        : %s\n", res.Proto()) +
			fmt.Sprintf("RECEIVED AT  : %v\n", res.ReceivedAt().Format(time.RFC3339Nano)) +
			fmt.Sprintf("TIME DURATION: %v\n", res.Time()) +
			"HEADERS      :\n" +
			composeHeaders(c, res.Request, rl.Header) + "\n"
		if res.Request.isSaveResponse {
			debugLog += "BODY         :\n***** RESPONSE WRITTEN INTO FILE *****\n"
		} else {
			debugLog += fmt.Sprintf("BODY         :\n%v\n", rl.Body)
		}
		debugLog += "==============================================================================\n"

		res.Request.log.Debugf("%s", debugLog)
	}

	return nil
}

func parseResponseBody(c *Client, res *Response) (err error) {
	if res.StatusCode() == http.StatusNoContent {
		res.Request.Error = nil
		return
	}
	// Handles only JSON or XML content type
	ct := firstNonEmpty(res.Request.forceContentType, res.Header().Get(hdrContentTypeKey), res.Request.fallbackContentType)
	if IsJSONType(ct) || IsXMLType(ct) {
		// HTTP status code > 199 and < 300, considered as Result
		if res.IsSuccess() {
			res.Request.Error = nil
			if res.Request.Result != nil {
				err = Unmarshalc(c, ct, res.body, res.Request.Result)
				return
			}
		}

		// HTTP status code > 399, considered as Error
		if res.IsError() {
			// global error interface
			if res.Request.Error == nil && c.Error != nil {
				res.Request.Error = reflect.New(c.Error).Interface()
			}

			if res.Request.Error != nil {
				unmarshalErr := Unmarshalc(c, ct, res.body, res.Request.Error)
				if unmarshalErr != nil {
					c.log.Warnf("Cannot unmarshal response body: %s", unmarshalErr)
				}
			}
		}
	}

	return
}

func handleMultipart(c *Client, r *Request) error {
	r.bodyBuf = acquireBuffer()
	w := multipart.NewWriter(r.bodyBuf)

	// Set boundary if not set by user
	if r.multipartBoundary != "" {
		if err := w.SetBoundary(r.multipartBoundary); err != nil {
			return err
		}
	}

	for k, v := range c.FormData {
		for _, iv := range v {
			if err := w.WriteField(k, iv); err != nil {
				return err
			}
		}
	}

	for k, v := range r.FormData {
		for _, iv := range v {
			if strings.HasPrefix(k, "@") { // file
				if err := addFile(w, k[1:], iv); err != nil {
					return err
				}
			} else { // form value
				if err := w.WriteField(k, iv); err != nil {
					return err
				}
			}
		}
	}

	// #21 - adding io.Reader support
	for _, f := range r.multipartFiles {
		if err := addFileReader(w, f); err != nil {
			return err
		}
	}

	// GitHub #130 adding multipart field support with content type
	for _, mf := range r.multipartFields {
		if err := addMultipartFormField(w, mf); err != nil {
			return err
		}
	}

	r.Header.Set(hdrContentTypeKey, w.FormDataContentType())
	return w.Close()
}

func handleFormData(c *Client, r *Request) {
	for k, v := range c.FormData {
		if _, ok := r.FormData[k]; ok {
			continue
		}
		r.FormData[k] = v[:]
	}

	r.bodyBuf = acquireBuffer()
	r.bodyBuf.WriteString(r.FormData.Encode())
	r.Header.Set(hdrContentTypeKey, formContentType)
	r.isFormData = true
}

func handleContentType(c *Client, r *Request) {
	contentType := r.Header.Get(hdrContentTypeKey)
	if IsStringEmpty(contentType) {
		contentType = DetectContentType(r.Body)
		r.Header.Set(hdrContentTypeKey, contentType)
	}
}

func handleRequestBody(c *Client, r *Request) error {
	var bodyBytes []byte
	r.bodyBuf = nil

	switch body := r.Body.(type) {
	case io.Reader:
		if c.setContentLength || r.setContentLength { // keep backward compatibility
			r.bodyBuf = acquireBuffer()
			if _, err := r.bodyBuf.ReadFrom(body); err != nil {
				return err
			}
			r.Body = nil
		} else {
			// Otherwise buffer less processing for `io.Reader`, sounds good.
			return nil
		}
	case []byte:
		bodyBytes = body
	case string:
		bodyBytes = []byte(body)
	default:
		contentType := r.Header.Get(hdrContentTypeKey)
		kind := kindOf(r.Body)
		var err error
		if IsJSONType(contentType) && (kind == reflect.Struct || kind == reflect.Map || kind == reflect.Slice) {
			r.bodyBuf, err = jsonMarshal(c, r, r.Body)
		} else if IsXMLType(contentType) && (kind == reflect.Struct) {
			bodyBytes, err = c.XMLMarshal(r.Body)
		}
		if err != nil {
			return err
		}
	}

	if bodyBytes == nil && r.bodyBuf == nil {
		return errors.New("unsupported 'Body' type/value")
	}

	// []byte into Buffer
	if bodyBytes != nil && r.bodyBuf == nil {
		r.bodyBuf = acquireBuffer()
		_, _ = r.bodyBuf.Write(bodyBytes)
	}

	return nil
}

func saveResponseIntoFile(c *Client, res *Response) error {
	if res.Request.isSaveResponse {
		file := ""

		if len(c.outputDirectory) > 0 && !filepath.IsAbs(res.Request.outputFile) {
			file += c.outputDirectory + string(filepath.Separator)
		}

		file = filepath.Clean(file + res.Request.outputFile)
		if err := createDirectory(filepath.Dir(file)); err != nil {
			return err
		}

		outFile, err := os.Create(file)
		if err != nil {
			return err
		}
		defer closeq(outFile)

		// io.Copy reads maximum 32kb size, it is perfect for large file download too
		defer closeq(res.RawResponse.Body)

		written, err := io.Copy(outFile, res.RawResponse.Body)
		if err != nil {
			return err
		}

		res.size = written
	}

	return nil
}

func getBodyCopy(r *Request) (*bytes.Buffer, error) {
	// If r.bodyBuf present, return the copy
	if r.bodyBuf != nil {
		bodyCopy := acquireBuffer()
		if _, err := io.Copy(bodyCopy, bytes.NewReader(r.bodyBuf.Bytes())); err != nil {
			// cannot use io.Copy(bodyCopy, r.bodyBuf) because io.Copy reset r.bodyBuf
			return nil, err
		}
		return bodyCopy, nil
	}

	// Maybe body is `io.Reader`.
	// Note: Resty user have to watchout for large body size of `io.Reader`
	if r.RawRequest.Body != nil {
		b, err := io.ReadAll(r.RawRequest.Body)
		if err != nil {
			return nil, err
		}

		// Restore the Body
		closeq(r.RawRequest.Body)
		r.RawRequest.Body = io.NopCloser(bytes.NewBuffer(b))

		// Return the Body bytes
		return bytes.NewBuffer(b), nil
	}
	return nil, nil
}
