// Copyright (C) 2015 . All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Interact with  API

// Package api contains client and functions to interact with  API
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

// https://cp-par1.scaleway.com/products/servers
// https://cp-ams1.scaleway.com/products/servers
// Default values
var (
	AccountAPI     = "https://account.scaleway.com/"
	MetadataAPI    = "http://169.254.42.42/"
	MarketplaceAPI = "https://api-marketplace.scaleway.com"
	ComputeAPIPar1 = "https://cp-par1.scaleway.com/"
	ComputeAPIAms1 = "https://cp-ams1.scaleway.com/"

	URLPublicDNS  = ".pub.cloud.scaleway.com"
	URLPrivateDNS = ".priv.cloud.scaleway.com"
)

func init() {
	if url := os.Getenv("SCW_ACCOUNT_API"); url != "" {
		AccountAPI = url
	}
	if url := os.Getenv("SCW_METADATA_API"); url != "" {
		MetadataAPI = url
	}
	if url := os.Getenv("SCW_MARKETPLACE_API"); url != "" {
		MarketplaceAPI = url
	}
}

const (
	perPage = 50
)

// HTTPClient wraps the net/http Client Do method
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// API is the interface used to communicate with the  API
type API struct {
	Organization string     // Organization is the identifier of the  organization
	Token        string     // Token is the authentication token for the  organization
	Client       HTTPClient // Client is used for all HTTP interactions

	password   string // Password is the authentication password
	userAgent  string
	computeAPI string

	Region string
}

// APIError represents a  API Error
type APIError struct {
	// Message is a human-friendly error message
	APIMessage string `json:"message,omitempty"`

	// Type is a string code that defines the kind of error
	Type string `json:"type,omitempty"`

	// Fields contains detail about validation error
	Fields map[string][]string `json:"fields,omitempty"`

	// StatusCode is the HTTP status code received
	StatusCode int `json:"-"`

	// Message
	Message string `json:"-"`
}

// Error returns a string representing the error
func (e APIError) Error() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "StatusCode: %v, ", e.StatusCode)
	fmt.Fprintf(&b, "Type: %v, ", e.Type)
	fmt.Fprintf(&b, "APIMessage: \x1b[31m%v\x1b[0m", e.APIMessage)
	if len(e.Fields) > 0 {
		fmt.Fprintf(&b, ", Details: %v", e.Fields)
	}
	return b.String()
}

// New creates a ready-to-use  SDK client
func New(organization, token, region string, options ...func(*API)) (*API, error) {
	s := &API{
		// exposed
		Organization: organization,
		Token:        token,
		Client:       &http.Client{},

		// internal
		password:  "",
		userAgent: "scaleway-sdk",
	}
	for _, option := range options {
		option(s)
	}
	switch region {
	case "par1", "":
		s.computeAPI = ComputeAPIPar1
	case "ams1":
		s.computeAPI = ComputeAPIAms1
	default:
		return nil, fmt.Errorf("%s isn't a valid region", region)
	}
	s.Region = region
	if url := os.Getenv("SCW_COMPUTE_API"); url != "" {
		s.computeAPI = url
	}
	return s, nil
}

func (s *API) response(method, uri string, content io.Reader) (resp *http.Response, err error) {
	var (
		req *http.Request
	)

	req, err = http.NewRequest(method, uri, content)
	if err != nil {
		err = fmt.Errorf("response %s %s", method, uri)
		return
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", s.userAgent)
	resp, err = s.Client.Do(req)
	return
}

// GetResponsePaginate fetchs all resources and returns an http.Response object for the requested resource
func (s *API) GetResponsePaginate(apiURL, resource string, values url.Values) (*http.Response, error) {
	resp, err := s.response("HEAD", fmt.Sprintf("%s/%s?%s", strings.TrimRight(apiURL, "/"), resource, values.Encode()), nil)
	if err != nil {
		return nil, err
	}

	count := resp.Header.Get("X-Total-Count")
	var maxElem int
	if count == "" {
		maxElem = 0
	} else {
		maxElem, err = strconv.Atoi(count)
		if err != nil {
			return nil, err
		}
	}

	get := maxElem / perPage
	if (float32(maxElem) / perPage) > float32(get) {
		get++
	}

	if get <= 1 { // If there is 0 or 1 page of result, the response is not paginated
		if len(values) == 0 {
			return s.response("GET", fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource), nil)
		}
		return s.response("GET", fmt.Sprintf("%s/%s?%s", strings.TrimRight(apiURL, "/"), resource, values.Encode()), nil)
	}

	fetchAll := !(values.Get("per_page") != "" || values.Get("page") != "")
	if fetchAll {
		var g errgroup.Group

		ch := make(chan *http.Response, get)
		for i := 1; i <= get; i++ {
			i := i // closure tricks
			g.Go(func() (err error) {
				var resp *http.Response

				val := url.Values{}
				val.Set("per_page", fmt.Sprintf("%v", perPage))
				val.Set("page", fmt.Sprintf("%v", i))
				resp, err = s.response("GET", fmt.Sprintf("%s/%s?%s", strings.TrimRight(apiURL, "/"), resource, val.Encode()), nil)
				ch <- resp
				return
			})
		}
		if err = g.Wait(); err != nil {
			return nil, err
		}
		newBody := make(map[string][]json.RawMessage)
		body := make(map[string][]json.RawMessage)
		key := ""
		for i := 0; i < get; i++ {
			res := <-ch
			if res.StatusCode != http.StatusOK {
				return res, nil
			}
			content, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(content, &body); err != nil {
				return nil, err
			}

			if i == 0 {
				resp = res
				for k := range body {
					key = k
					break
				}
			}
			newBody[key] = append(newBody[key], body[key]...)
		}
		payload := new(bytes.Buffer)
		if err := json.NewEncoder(payload).Encode(newBody); err != nil {
			return nil, err
		}
		resp.Body = ioutil.NopCloser(payload)
	} else {
		resp, err = s.response("GET", fmt.Sprintf("%s/%s?%s", strings.TrimRight(apiURL, "/"), resource, values.Encode()), nil)
	}
	return resp, err
}

// PostResponse returns an http.Response object for the updated resource
func (s *API) PostResponse(apiURL, resource string, data interface{}) (*http.Response, error) {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(data); err != nil {
		return nil, err
	}
	return s.response("POST", fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource), payload)
}

// PatchResponse returns an http.Response object for the updated resource
func (s *API) PatchResponse(apiURL, resource string, data interface{}) (*http.Response, error) {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(data); err != nil {
		return nil, err
	}
	return s.response("PATCH", fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource), payload)
}

// PutResponse returns an http.Response object for the updated resource
func (s *API) PutResponse(apiURL, resource string, data interface{}) (*http.Response, error) {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(data); err != nil {
		return nil, err
	}
	return s.response("PUT", fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource), payload)
}

// DeleteResponse returns an http.Response object for the deleted resource
func (s *API) DeleteResponse(apiURL, resource string) (*http.Response, error) {
	return s.response("DELETE", fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource), nil)
}

// handleHTTPError checks the statusCode and displays the error
func (s *API) handleHTTPError(goodStatusCode []int, resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		return nil, errors.New(string(body))
	}

	for _, code := range goodStatusCode {
		if code == resp.StatusCode {
			return body, nil
		}
	}

	var scwError APIError
	scwError.StatusCode = resp.StatusCode
	if len(body) > 0 {
		if err := json.Unmarshal(body, &scwError); err != nil {
			return nil, err
		}
	}
	return nil, scwError
}

// SetPassword register the password
func (s *API) SetPassword(password string) {
	s.password = password
}
