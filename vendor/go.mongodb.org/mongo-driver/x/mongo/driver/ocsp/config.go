// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package ocsp

import (
	"crypto/x509"
	"errors"
	"fmt"

	"golang.org/x/crypto/ocsp"
)

type config struct {
	serverCert, issuer      *x509.Certificate
	cache                   Cache
	disableEndpointChecking bool
	ocspRequest             *ocsp.Request
	ocspRequestBytes        []byte
}

func newConfig(certChain []*x509.Certificate, opts *VerifyOptions) (config, error) {
	cfg := config{
		cache:                   opts.Cache,
		disableEndpointChecking: opts.DisableEndpointChecking,
	}

	if len(certChain) == 0 {
		return cfg, errors.New("verified certificate chain contained no certificates")
	}

	// In the case where the leaf certificate and CA are the same, the chain may only contain one certificate.
	cfg.serverCert = certChain[0]
	cfg.issuer = certChain[0]
	if len(certChain) > 1 {
		// If the chain has multiple certificates, the one directly after the leaf should be the issuer. Use
		// CheckSignatureFrom to verify that it is the issuer.
		cfg.issuer = certChain[1]

		if err := cfg.serverCert.CheckSignatureFrom(cfg.issuer); err != nil {
			errString := "error checking if server certificate is signed by the issuer in the verified chain: %v"
			return cfg, fmt.Errorf(errString, err)
		}
	}

	var err error
	cfg.ocspRequestBytes, err = ocsp.CreateRequest(cfg.serverCert, cfg.issuer, nil)
	if err != nil {
		return cfg, fmt.Errorf("error creating OCSP request: %v", err)
	}
	cfg.ocspRequest, err = ocsp.ParseRequest(cfg.ocspRequestBytes)
	if err != nil {
		return cfg, fmt.Errorf("error parsing OCSP request bytes: %v", err)
	}

	return cfg, nil
}
