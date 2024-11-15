// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package ocsp is intended for internal use only. It is made available to
// facilitate use cases that require access to internal MongoDB driver
// functionality and state. The API of this package is not stable and there is
// no backward compatibility guarantee.
//
// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package ocsp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"golang.org/x/crypto/ocsp"
	"golang.org/x/sync/errgroup"
)

var (
	tlsFeatureExtensionOID = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 1, 24}
	mustStapleFeatureValue = big.NewInt(5)
)

// Error represents an OCSP verification error
type Error struct {
	wrapped error
}

// Error implements the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("OCSP verification failed: %v", e.wrapped)
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.wrapped
}

func newOCSPError(wrapped error) error {
	return &Error{wrapped: wrapped}
}

// ResponseDetails contains a subset of the details needed from an OCSP response after the original response has been
// validated.
type ResponseDetails struct {
	Status     int
	NextUpdate time.Time
}

func extractResponseDetails(res *ocsp.Response) *ResponseDetails {
	return &ResponseDetails{
		Status:     res.Status,
		NextUpdate: res.NextUpdate,
	}
}

// Verify performs OCSP verification for the provided ConnectionState instance.
func Verify(ctx context.Context, connState tls.ConnectionState, opts *VerifyOptions) error {
	if opts.Cache == nil {
		// There should always be an OCSP cache. Even if the user has specified the URI option to disable communication
		// with OCSP responders, the driver will cache any stapled responses. Requiring that the cache is non-nil
		// allows us to confirm that the cache is correctly being passed down from a higher level.
		return newOCSPError(errors.New("no OCSP cache provided"))
	}
	if len(connState.VerifiedChains) == 0 {
		return newOCSPError(errors.New("no verified certificate chains reported after TLS handshake"))
	}

	certChain := connState.VerifiedChains[0]
	if numCerts := len(certChain); numCerts == 0 {
		return newOCSPError(errors.New("verified chain contained no certificates"))
	}

	ocspCfg, err := newConfig(certChain, opts)
	if err != nil {
		return newOCSPError(err)
	}

	res, err := getParsedResponse(ctx, ocspCfg, connState)
	if err != nil {
		return err
	}
	if res == nil {
		// If no response was parsed from the staple and responders, the status of the certificate is unknown, so don't
		// error.
		return nil
	}

	if res.Status == ocsp.Revoked {
		return newOCSPError(errors.New("certificate is revoked"))
	}
	return nil
}

// getParsedResponse attempts to parse a response from the stapled OCSP data or by contacting OCSP responders if no
// staple is present.
func getParsedResponse(ctx context.Context, cfg config, connState tls.ConnectionState) (*ResponseDetails, error) {
	stapledResponse, err := processStaple(cfg, connState.OCSPResponse)
	if err != nil {
		return nil, err
	}

	if stapledResponse != nil {
		// If there is a staple, attempt to cache it. The cache.Update call will resolve conflicts with an existing
		// cache enry if necessary.
		return cfg.cache.Update(cfg.ocspRequest, stapledResponse), nil
	}
	if cachedResponse := cfg.cache.Get(cfg.ocspRequest); cachedResponse != nil {
		return cachedResponse, nil
	}

	// If there is no stapled or cached response, fall back to querying the responders if that functionality has not
	// been disabled.
	if cfg.disableEndpointChecking {
		return nil, nil
	}
	externalResponse := contactResponders(ctx, cfg)
	if externalResponse == nil {
		// None of the responders were available.
		return nil, nil
	}

	// Similar to the stapled response case above, unconditionally call Update and it will either cache the response
	// or resolve conflicts if a different connection has cached a response since the previous call to Get.
	return cfg.cache.Update(cfg.ocspRequest, externalResponse), nil
}

// processStaple returns the OCSP response from the provided staple. An error will be returned if any of the following
// are true:
//
// 1. cfg.serverCert has the Must-Staple extension but the staple is empty.
// 2. The staple is malformed.
// 3. The staple does not cover cfg.serverCert.
// 4. The OCSP response has an error status.
func processStaple(cfg config, staple []byte) (*ResponseDetails, error) {
	mustStaple, err := isMustStapleCertificate(cfg.serverCert)
	if err != nil {
		return nil, err
	}

	// If the server has a Must-Staple certificate and the server does not present a stapled OCSP response, error.
	if mustStaple && len(staple) == 0 {
		return nil, errors.New("server provided a certificate with the Must-Staple extension but did not " +
			"provide a stapled OCSP response")
	}

	if len(staple) == 0 {
		return nil, nil
	}

	parsedResponse, err := ocsp.ParseResponseForCert(staple, cfg.serverCert, cfg.issuer)
	if err != nil {
		// If the stapled response could not be parsed correctly, error. This can happen if the response is malformed,
		// the response does not cover the certificate presented by the server, or if the response contains an error
		// status.
		return nil, fmt.Errorf("error parsing stapled response: %w", err)
	}
	if err = verifyResponse(cfg, parsedResponse); err != nil {
		return nil, fmt.Errorf("error validating stapled response: %w", err)
	}

	return extractResponseDetails(parsedResponse), nil
}

// isMustStapleCertificate determines whether or not an X509 certificate is a must-staple certificate.
func isMustStapleCertificate(cert *x509.Certificate) (bool, error) {
	var featureExtension pkix.Extension
	var foundExtension bool
	for _, ext := range cert.Extensions {
		if ext.Id.Equal(tlsFeatureExtensionOID) {
			featureExtension = ext
			foundExtension = true
			break
		}
	}
	if !foundExtension {
		return false, nil
	}

	// The value for the TLS feature extension is a sequence of integers. Per the asn1.Unmarshal documentation, an
	// integer can be unmarshalled into an int, int32, int64, or *big.Int and unmarshalling will error if the integer
	// cannot be encoded into the target type.
	//
	// Use []*big.Int to ensure that all values in the sequence can be successfully unmarshalled.
	var featureValues []*big.Int
	if _, err := asn1.Unmarshal(featureExtension.Value, &featureValues); err != nil {
		return false, fmt.Errorf("error unmarshalling TLS feature extension values: %w", err)
	}

	for _, value := range featureValues {
		if value.Cmp(mustStapleFeatureValue) == 0 {
			return true, nil
		}
	}
	return false, nil
}

// contactResponders will send a request to all OCSP responders reported by cfg.serverCert. The
// first response that conclusively identifies cfg.serverCert as good or revoked will be returned.
// If all responders are unavailable or no responder returns a conclusive status, it returns nil.
// contactResponders will wait for up to 5 seconds to get a certificate status response.
func contactResponders(ctx context.Context, cfg config) *ResponseDetails {
	if len(cfg.serverCert.OCSPServer) == 0 {
		return nil
	}

	// Limit all OCSP responder calls to a maximum of 5 seconds or when the passed-in context expires,
	// whichever happens first.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)
	ocspResponses := make(chan *ocsp.Response, len(cfg.serverCert.OCSPServer))
	defer close(ocspResponses)

	for _, endpoint := range cfg.serverCert.OCSPServer {
		// Re-assign endpoint so it gets re-scoped rather than using the iteration variable in the goroutine. See
		// https://golang.org/doc/faq#closures_and_goroutines.
		endpoint := endpoint

		// Start a group of goroutines that each attempt to request the certificate status from one
		// of the OCSP endpoints listed in the server certificate. We want to "soft fail" on all
		// errors, so this function never returns actual errors. Only a "done" error is returned
		// when a response is received so the errgroup cancels any other in-progress requests.
		group.Go(func() error {
			// Use bytes.NewReader instead of bytes.NewBuffer because a bytes.Buffer is an owning representation and the
			// docs recommend not using the underlying []byte after creating the buffer, so a new copy of the request
			// bytes would be needed for each request.
			request, err := http.NewRequest("POST", endpoint, bytes.NewReader(cfg.ocspRequestBytes))
			if err != nil {
				return nil
			}
			request = request.WithContext(ctx)

			httpResponse, err := cfg.httpClient.Do(request)
			if err != nil {
				return nil
			}
			defer func() {
				_ = httpResponse.Body.Close()
			}()

			if httpResponse.StatusCode != 200 {
				return nil
			}

			httpBytes, err := ioutil.ReadAll(httpResponse.Body)
			if err != nil {
				return nil
			}

			ocspResponse, err := ocsp.ParseResponseForCert(httpBytes, cfg.serverCert, cfg.issuer)
			if err != nil || verifyResponse(cfg, ocspResponse) != nil || ocspResponse.Status == ocsp.Unknown {
				// If there was an error parsing/validating the response or the response was
				// inconclusive, suppress the error because we want to ignore this responder.
				return nil
			}

			// Send the conclusive response on the response channel and return a "done" error that
			// will cause the errgroup to cancel all other in-progress requests.
			ocspResponses <- ocspResponse
			return errors.New("done")
		})
	}

	_ = group.Wait()
	select {
	case res := <-ocspResponses:
		return extractResponseDetails(res)
	default:
		// If there is no OCSP response on the response channel, all OCSP calls either failed or
		// were inconclusive. Return nil.
		return nil
	}
}

// verifyResponse checks that the provided OCSP response is valid.
func verifyResponse(cfg config, res *ocsp.Response) error {
	if err := verifyExtendedKeyUsage(cfg, res); err != nil {
		return err
	}

	currTime := time.Now().UTC()
	if res.ThisUpdate.After(currTime) {
		return fmt.Errorf("reported thisUpdate time %s is after current time %s", res.ThisUpdate, currTime)
	}
	if !res.NextUpdate.IsZero() && res.NextUpdate.Before(currTime) {
		return fmt.Errorf("reported nextUpdate time %s is before current time %s", res.NextUpdate, currTime)
	}
	return nil
}

func verifyExtendedKeyUsage(cfg config, res *ocsp.Response) error {
	if res.Certificate == nil {
		return nil
	}

	namesMatch := res.RawResponderName != nil && bytes.Equal(res.RawResponderName, cfg.issuer.RawSubject)
	keyHashesMatch := res.ResponderKeyHash != nil && bytes.Equal(res.ResponderKeyHash, cfg.ocspRequest.IssuerKeyHash)
	if namesMatch || keyHashesMatch {
		// The responder certificate is the same as the issuer certificate.
		return nil
	}

	// There is a delegate.
	for _, extKeyUsage := range res.Certificate.ExtKeyUsage {
		if extKeyUsage == x509.ExtKeyUsageOCSPSigning {
			return nil
		}
	}

	return errors.New("delegate responder certificate is missing the OCSP signing extended key usage")
}
