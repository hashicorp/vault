// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Based on github.com/aws/aws-sdk-go by Amazon.com, Inc. with code from:
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/request/request.go
// See THIRD-PARTY-NOTICES for original license terms

package awsv4

import (
	"net/http"
	"strings"
)

// Returns host from request
func getHost(r *http.Request) string {
	if r.Host != "" {
		return r.Host
	}

	if r.URL == nil {
		return ""
	}

	return r.URL.Host
}

// Hostname returns u.Host, without any port number.
//
// If Host is an IPv6 literal with a port number, Hostname returns the
// IPv6 literal without the square brackets. IPv6 literals may include
// a zone identifier.
//
// Copied from the Go 1.8 standard library (net/url)
func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

// Port returns the port part of u.Host, without the leading colon.
// If u.Host doesn't contain a port, Port returns an empty string.
//
// Copied from the Go 1.8 standard library (net/url)
func portOnly(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return ""
	}
	if i := strings.Index(hostport, "]:"); i != -1 {
		return hostport[i+len("]:"):]
	}
	if strings.Contains(hostport, "]") {
		return ""
	}
	return hostport[colon+len(":"):]
}

// Returns true if the specified URI is using the standard port
// (i.e. port 80 for HTTP URIs or 443 for HTTPS URIs)
func isDefaultPort(scheme, port string) bool {
	if port == "" {
		return true
	}

	lowerCaseScheme := strings.ToLower(scheme)
	if (lowerCaseScheme == "http" && port == "80") || (lowerCaseScheme == "https" && port == "443") {
		return true
	}

	return false
}
