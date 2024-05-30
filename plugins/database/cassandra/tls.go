// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cassandra

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
)

func jsonBundleToTLSConfig(rawJSON string, tlsMinVersion uint16, serverName string, insecureSkipVerify bool) (*tls.Config, error) {
	var certBundle certutil.CertBundle
	err := json.Unmarshal([]byte(rawJSON), &certBundle)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if certBundle.IssuingCA != "" && len(certBundle.CAChain) > 0 {
		return nil, fmt.Errorf("issuing_ca and ca_chain cannot both be specified")
	}
	if certBundle.IssuingCA != "" {
		certBundle.CAChain = []string{certBundle.IssuingCA}
		certBundle.IssuingCA = ""
	}

	return toClientTLSConfig(certBundle.Certificate, certBundle.PrivateKey, certBundle.CAChain, tlsMinVersion, serverName, insecureSkipVerify)
}

func pemBundleToTLSConfig(pemBundle string, tlsMinVersion uint16, serverName string, insecureSkipVerify bool) (*tls.Config, error) {
	if len(pemBundle) == 0 {
		return nil, errutil.UserError{Err: "empty pem bundle"}
	}

	pemBytes := []byte(pemBundle)
	var pemBlock *pem.Block

	certificate := ""
	privateKey := ""
	caChain := []string{}

	for len(pemBytes) > 0 {
		pemBlock, pemBytes = pem.Decode(pemBytes)
		if pemBlock == nil {
			return nil, errutil.UserError{Err: "no data found in PEM block"}
		}
		blockBytes := pem.EncodeToMemory(pemBlock)

		switch pemBlock.Type {
		case "CERTIFICATE":
			// Parse the cert so we know if it's a CA or not
			cert, err := x509.ParseCertificate(pemBlock.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate: %w", err)
			}
			if cert.IsCA {
				caChain = append(caChain, string(blockBytes))
				continue
			}

			// Only one leaf certificate supported
			if certificate != "" {
				return nil, errutil.UserError{Err: "multiple leaf certificates not supported"}
			}
			certificate = string(blockBytes)

		case "RSA PRIVATE KEY", "EC PRIVATE KEY", "PRIVATE KEY":
			if privateKey != "" {
				return nil, errutil.UserError{Err: "multiple private keys not supported"}
			}
			privateKey = string(blockBytes)
		default:
			return nil, fmt.Errorf("unsupported PEM block type [%s]", pemBlock.Type)
		}
	}

	return toClientTLSConfig(certificate, privateKey, caChain, tlsMinVersion, serverName, insecureSkipVerify)
}

func toClientTLSConfig(certificatePEM string, privateKeyPEM string, caChainPEMs []string, tlsMinVersion uint16, serverName string, insecureSkipVerify bool) (*tls.Config, error) {
	if certificatePEM != "" && privateKeyPEM == "" {
		return nil, fmt.Errorf("found certificate for client-side TLS authentication but no private key")
	} else if certificatePEM == "" && privateKeyPEM != "" {
		return nil, fmt.Errorf("found private key for client-side TLS authentication but no certificate")
	}

	var certificates []tls.Certificate
	if certificatePEM != "" {
		certificate, err := tls.X509KeyPair([]byte(certificatePEM), []byte(privateKeyPEM))
		if err != nil {
			return nil, fmt.Errorf("failed to parse certificate and private key pair: %w", err)
		}
		certificates = append(certificates, certificate)
	}

	var rootCAs *x509.CertPool
	if len(caChainPEMs) > 0 {
		rootCAs = x509.NewCertPool()
		for _, caBlock := range caChainPEMs {
			ok := rootCAs.AppendCertsFromPEM([]byte(caBlock))
			if !ok {
				return nil, fmt.Errorf("failed to add CA certificate to certificate pool: it may be malformed or empty")
			}
		}
	}

	config := &tls.Config{
		Certificates:       certificates,
		RootCAs:            rootCAs,
		ServerName:         serverName,
		InsecureSkipVerify: insecureSkipVerify,
		MinVersion:         tlsMinVersion,
	}
	return config, nil
}
