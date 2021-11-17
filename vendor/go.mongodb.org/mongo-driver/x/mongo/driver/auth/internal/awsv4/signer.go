// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Based on github.com/aws/aws-sdk-go by Amazon.com, Inc. with code from:
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/request/request.go
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/signer/v4/v4.go
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/signer/v4/uri_path.go
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/types.go
// See THIRD-PARTY-NOTICES for original license terms

package awsv4

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	authorizationHeader     = "Authorization"
	authHeaderSignatureElem = "Signature="

	authHeaderPrefix = "AWS4-HMAC-SHA256"
	timeFormat       = "20060102T150405Z"
	shortTimeFormat  = "20060102"
	awsV4Request     = "aws4_request"

	// emptyStringSHA256 is a SHA256 of an empty string
	emptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
)

var ignoredHeaders = rules{
	blacklist{
		mapRule{
			authorizationHeader: struct{}{},
			"User-Agent":        struct{}{},
			"X-Amzn-Trace-Id":   struct{}{},
		},
	},
}

// Signer applies AWS v4 signing to given request. Use this to sign requests
// that need to be signed with AWS V4 Signatures.
type Signer struct {
	Credentials *StaticProvider
}

// NewSigner returns a Signer pointer configured with the credentials and optional
// option values provided. If not options are provided the Signer will use its
// default configuration.
func NewSigner(credentials *StaticProvider) *Signer {
	v4 := &Signer{
		Credentials: credentials,
	}

	return v4
}

type signingCtx struct {
	ServiceName      string
	Region           string
	Request          *http.Request
	Body             io.ReadSeeker
	Query            url.Values
	Time             time.Time
	SignedHeaderVals http.Header

	credValues Value

	bodyDigest       string
	signedHeaders    string
	canonicalHeaders string
	canonicalString  string
	credentialString string
	stringToSign     string
	signature        string
	authorization    string
}

// Sign signs AWS v4 requests with the provided body, service name, region the
// request is made to, and time the request is signed at. The signTime allows
// you to specify that a request is signed for the future, and cannot be
// used until then.
//
// Returns a list of HTTP headers that were included in the signature or an
// error if signing the request failed. Generally for signed requests this value
// is not needed as the full request context will be captured by the http.Request
// value. It is included for reference though.
//
// Sign will set the request's Body to be the `body` parameter passed in. If
// the body is not already an io.ReadCloser, it will be wrapped within one. If
// a `nil` body parameter passed to Sign, the request's Body field will be
// also set to nil. Its important to note that this functionality will not
// change the request's ContentLength of the request.
//
// Sign differs from Presign in that it will sign the request using HTTP
// header values. This type of signing is intended for http.Request values that
// will not be shared, or are shared in a way the header values on the request
// will not be lost.
//
// The requests body is an io.ReadSeeker so the SHA256 of the body can be
// generated. To bypass the signer computing the hash you can set the
// "X-Amz-Content-Sha256" header with a precomputed value. The signer will
// only compute the hash if the request header value is empty.
func (v4 Signer) Sign(r *http.Request, body io.ReadSeeker, service, region string, signTime time.Time) (http.Header, error) {
	return v4.signWithBody(r, body, service, region, signTime)
}

func (v4 Signer) signWithBody(r *http.Request, body io.ReadSeeker, service, region string, signTime time.Time) (http.Header, error) {
	ctx := &signingCtx{
		Request:     r,
		Body:        body,
		Query:       r.URL.Query(),
		Time:        signTime,
		ServiceName: service,
		Region:      region,
	}

	for key := range ctx.Query {
		sort.Strings(ctx.Query[key])
	}

	if ctx.isRequestSigned() {
		ctx.Time = time.Now()
	}

	var err error
	ctx.credValues, err = v4.Credentials.Retrieve()
	if err != nil {
		return http.Header{}, err
	}

	ctx.sanitizeHostForHeader()
	ctx.assignAmzQueryValues()
	if err := ctx.build(); err != nil {
		return nil, err
	}

	var reader io.ReadCloser
	if body != nil {
		var ok bool
		if reader, ok = body.(io.ReadCloser); !ok {
			reader = ioutil.NopCloser(body)
		}
	}
	r.Body = reader

	return ctx.SignedHeaderVals, nil
}

// sanitizeHostForHeader removes default port from host and updates request.Host
func (ctx *signingCtx) sanitizeHostForHeader() {
	r := ctx.Request
	host := getHost(r)
	port := portOnly(host)
	if port != "" && isDefaultPort(r.URL.Scheme, port) {
		r.Host = stripPort(host)
	}
}

func (ctx *signingCtx) assignAmzQueryValues() {
	if ctx.credValues.SessionToken != "" {
		ctx.Request.Header.Set("X-Amz-Security-Token", ctx.credValues.SessionToken)
	}
}

func (ctx *signingCtx) build() error {
	ctx.buildTime()             // no depends
	ctx.buildCredentialString() // no depends

	if err := ctx.buildBodyDigest(); err != nil {
		return err
	}

	unsignedHeaders := ctx.Request.Header

	ctx.buildCanonicalHeaders(ignoredHeaders, unsignedHeaders)
	ctx.buildCanonicalString() // depends on canon headers / signed headers
	ctx.buildStringToSign()    // depends on canon string
	ctx.buildSignature()       // depends on string to sign

	parts := []string{
		authHeaderPrefix + " Credential=" + ctx.credValues.AccessKeyID + "/" + ctx.credentialString,
		"SignedHeaders=" + ctx.signedHeaders,
		authHeaderSignatureElem + ctx.signature,
	}
	ctx.Request.Header.Set(authorizationHeader, strings.Join(parts, ", "))

	return nil
}

// GetSignedRequestSignature attempts to extract the signature of the request.
// Returning an error if the request is unsigned, or unable to extract the
// signature.
func GetSignedRequestSignature(r *http.Request) ([]byte, error) {

	if auth := r.Header.Get(authorizationHeader); len(auth) != 0 {
		ps := strings.Split(auth, ", ")
		for _, p := range ps {
			if idx := strings.Index(p, authHeaderSignatureElem); idx >= 0 {
				sig := p[len(authHeaderSignatureElem):]
				if len(sig) == 0 {
					return nil, fmt.Errorf("invalid request signature authorization header")
				}
				return hex.DecodeString(sig)
			}
		}
	}

	if sig := r.URL.Query().Get("X-Amz-Signature"); len(sig) != 0 {
		return hex.DecodeString(sig)
	}

	return nil, fmt.Errorf("request not signed")
}

func (ctx *signingCtx) buildTime() {
	ctx.Request.Header.Set("X-Amz-Date", formatTime(ctx.Time))
}

func (ctx *signingCtx) buildCredentialString() {
	ctx.credentialString = buildSigningScope(ctx.Region, ctx.ServiceName, ctx.Time)
}

func (ctx *signingCtx) buildCanonicalHeaders(r rule, header http.Header) {
	var headers []string
	headers = append(headers, "host")
	for k, v := range header {
		if !r.IsValid(k) {
			continue // ignored header
		}
		if ctx.SignedHeaderVals == nil {
			ctx.SignedHeaderVals = make(http.Header)
		}

		lowerCaseKey := strings.ToLower(k)
		if _, ok := ctx.SignedHeaderVals[lowerCaseKey]; ok {
			// include additional values
			ctx.SignedHeaderVals[lowerCaseKey] = append(ctx.SignedHeaderVals[lowerCaseKey], v...)
			continue
		}

		headers = append(headers, lowerCaseKey)
		ctx.SignedHeaderVals[lowerCaseKey] = v
	}
	sort.Strings(headers)

	ctx.signedHeaders = strings.Join(headers, ";")

	headerValues := make([]string, len(headers))
	for i, k := range headers {
		if k == "host" {
			if ctx.Request.Host != "" {
				headerValues[i] = "host:" + ctx.Request.Host
			} else {
				headerValues[i] = "host:" + ctx.Request.URL.Host
			}
		} else {
			headerValues[i] = k + ":" +
				strings.Join(ctx.SignedHeaderVals[k], ",")
		}
	}
	stripExcessSpaces(headerValues)
	ctx.canonicalHeaders = strings.Join(headerValues, "\n")
}

func getURIPath(u *url.URL) string {
	var uri string

	if len(u.Opaque) > 0 {
		uri = "/" + strings.Join(strings.Split(u.Opaque, "/")[3:], "/")
	} else {
		uri = u.EscapedPath()
	}

	if len(uri) == 0 {
		uri = "/"
	}

	return uri
}

func (ctx *signingCtx) buildCanonicalString() {
	ctx.Request.URL.RawQuery = strings.Replace(ctx.Query.Encode(), "+", "%20", -1)

	uri := getURIPath(ctx.Request.URL)

	uri = EscapePath(uri, false)

	ctx.canonicalString = strings.Join([]string{
		ctx.Request.Method,
		uri,
		ctx.Request.URL.RawQuery,
		ctx.canonicalHeaders + "\n",
		ctx.signedHeaders,
		ctx.bodyDigest,
	}, "\n")
}

func (ctx *signingCtx) buildStringToSign() {
	ctx.stringToSign = strings.Join([]string{
		authHeaderPrefix,
		formatTime(ctx.Time),
		ctx.credentialString,
		hex.EncodeToString(hashSHA256([]byte(ctx.canonicalString))),
	}, "\n")
}

func (ctx *signingCtx) buildSignature() {
	creds := deriveSigningKey(ctx.Region, ctx.ServiceName, ctx.credValues.SecretAccessKey, ctx.Time)
	signature := hmacSHA256(creds, []byte(ctx.stringToSign))
	ctx.signature = hex.EncodeToString(signature)
}

func (ctx *signingCtx) buildBodyDigest() error {
	hash := ctx.Request.Header.Get("X-Amz-Content-Sha256")
	if hash == "" {
		if ctx.Body == nil {
			hash = emptyStringSHA256
		} else {
			hashBytes, err := makeSha256Reader(ctx.Body)
			if err != nil {
				return err
			}
			hash = hex.EncodeToString(hashBytes)
		}
	}
	ctx.bodyDigest = hash

	return nil
}

// isRequestSigned returns if the request is currently signed or presigned
func (ctx *signingCtx) isRequestSigned() bool {
	return ctx.Request.Header.Get("Authorization") != ""
}

func hmacSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func hashSHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// seekerLen attempts to get the number of bytes remaining at the seeker's
// current position.  Returns the number of bytes remaining or error.
func seekerLen(s io.Seeker) (int64, error) {
	curOffset, err := s.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	endOffset, err := s.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = s.Seek(curOffset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return endOffset - curOffset, nil
}

func makeSha256Reader(reader io.ReadSeeker) (hashBytes []byte, err error) {
	hash := sha256.New()
	start, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	defer func() {
		// ensure error is return if unable to seek back to start of payload.
		_, err = reader.Seek(start, io.SeekStart)
	}()

	// Use CopyN to avoid allocating the 32KB buffer in io.Copy for bodies
	// smaller than 32KB. Fall back to io.Copy if we fail to determine the size.
	size, err := seekerLen(reader)
	if err != nil {
		_, _ = io.Copy(hash, reader)
	} else {
		_, _ = io.CopyN(hash, reader, size)
	}

	return hash.Sum(nil), nil
}

const doubleSpace = "  "

// stripExcessSpaces will rewrite the passed in slice's string values to not
// contain multiple side-by-side spaces.
func stripExcessSpaces(vals []string) {
	var j, k, l, m, spaces int
	for i, str := range vals {
		// Trim trailing spaces
		for j = len(str) - 1; j >= 0 && str[j] == ' '; j-- {
		}

		// Trim leading spaces
		for k = 0; k < j && str[k] == ' '; k++ {
		}
		str = str[k : j+1]

		// Strip multiple spaces.
		j = strings.Index(str, doubleSpace)
		if j < 0 {
			vals[i] = str
			continue
		}

		buf := []byte(str)
		for k, m, l = j, j, len(buf); k < l; k++ {
			if buf[k] == ' ' {
				if spaces == 0 {
					// First space.
					buf[m] = buf[k]
					m++
				}
				spaces++
			} else {
				// End of multiple spaces.
				spaces = 0
				buf[m] = buf[k]
				m++
			}
		}

		vals[i] = string(buf[:m])
	}
}

func buildSigningScope(region, service string, dt time.Time) string {
	return strings.Join([]string{
		formatShortTime(dt),
		region,
		service,
		awsV4Request,
	}, "/")
}

func deriveSigningKey(region, service, secretKey string, dt time.Time) []byte {
	keyDate := hmacSHA256([]byte("AWS4"+secretKey), []byte(formatShortTime(dt)))
	keyRegion := hmacSHA256(keyDate, []byte(region))
	keyService := hmacSHA256(keyRegion, []byte(service))
	signingKey := hmacSHA256(keyService, []byte(awsV4Request))
	return signingKey
}

func formatShortTime(dt time.Time) string {
	return dt.UTC().Format(shortTimeFormat)
}

func formatTime(dt time.Time) string {
	return dt.UTC().Format(timeFormat)
}
