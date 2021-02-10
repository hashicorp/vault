// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.

package auth

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
	"strings"
)

const (
	metadataBaseURL             = `http://169.254.169.254/opc/v2`
	metadataFallbackURL         = `http://169.254.169.254/opc/v1`
	regionPath                  = `/instance/region`
	leafCertificatePath         = `/identity/cert.pem`
	leafCertificateKeyPath      = `/identity/key.pem`
	intermediateCertificatePath = `/identity/intermediate.pem`

	leafCertificateKeyPassphrase         = `` // No passphrase for the private key for Compute instances
	intermediateCertificateKeyURL        = ``
	intermediateCertificateKeyPassphrase = `` // No passphrase for the private key for Compute instances
)

var (
	regionURL, leafCertificateURL, leafCertificateKeyURL, intermediateCertificateURL string
)

// instancePrincipalKeyProvider implements KeyProvider to provide a key ID and its corresponding private key
// for an instance principal by getting a security token via x509FederationClient.
//
// The region name of the endpoint for x509FederationClient is obtained from the metadata service on the compute
// instance.
type instancePrincipalKeyProvider struct {
	Region           common.Region
	FederationClient federationClient
	TenancyID        string
}

// newInstancePrincipalKeyProvider creates and returns an instancePrincipalKeyProvider instance based on
// x509FederationClient.
//
// NOTE: There is a race condition between PrivateRSAKey() and KeyID().  These two pieces are tightly coupled; KeyID
// includes a security token obtained from Auth service by giving a public key which is paired with PrivateRSAKey.
// The x509FederationClient caches the security token in memory until it is expired.  Thus, even if a client obtains a
// KeyID that is not expired at the moment, the PrivateRSAKey that the client acquires at a next moment could be
// invalid because the KeyID could be already expired.
func newInstancePrincipalKeyProvider(modifier func(common.HTTPRequestDispatcher) (common.HTTPRequestDispatcher, error)) (provider *instancePrincipalKeyProvider, err error) {
	updateX509CertRetrieverURLParas(metadataBaseURL)
	clientModifier := newDispatcherModifier(modifier)

	client, err := clientModifier.Modify(&http.Client{})
	if err != nil {
		err = fmt.Errorf("failed to modify client: %s", err.Error())
		return nil, err
	}

	var region common.Region

	if region, err = getRegionForFederationClient(client, regionURL); err != nil {
		err = fmt.Errorf("failed to get the region name from %s: %s", regionURL, err.Error())
		common.Logf("%v\n", err)
		return nil, err
	}

	leafCertificateRetriever := newURLBasedX509CertificateRetriever(client,
		leafCertificateURL, leafCertificateKeyURL, leafCertificateKeyPassphrase)
	intermediateCertificateRetrievers := []x509CertificateRetriever{
		newURLBasedX509CertificateRetriever(
			client, intermediateCertificateURL, intermediateCertificateKeyURL,
			intermediateCertificateKeyPassphrase),
	}

	if err = leafCertificateRetriever.Refresh(); err != nil {
		err = fmt.Errorf("failed to refresh the leaf certificate: %s", err.Error())
		return nil, err
	}
	tenancyID := extractTenancyIDFromCertificate(leafCertificateRetriever.Certificate())

	federationClient, err := newX509FederationClient(region, tenancyID, leafCertificateRetriever, intermediateCertificateRetrievers, *clientModifier)

	if err != nil {
		err = fmt.Errorf("failed to create federation client: %s", err.Error())
		return nil, err
	}

	provider = &instancePrincipalKeyProvider{FederationClient: federationClient, TenancyID: tenancyID, Region: region}
	return
}

func getRegionForFederationClient(dispatcher common.HTTPRequestDispatcher, url string) (r common.Region, err error) {
	var body bytes.Buffer
	var statusCode int
	if body, statusCode, err = httpGet(dispatcher, url); err != nil {
		if statusCode == 404 && strings.Compare(url, metadataBaseURL+regionPath) == 0 {
			common.Logf("Falling back to http://169.254.169.254/opc/v1 to try again...\n")
			updateX509CertRetrieverURLParas(metadataFallbackURL)
			return getRegionForFederationClient(dispatcher, regionURL)
		}
		return
	}
	return common.StringToRegion(body.String()), nil
}

func updateX509CertRetrieverURLParas(baseURL string) {
	regionURL = baseURL + regionPath
	leafCertificateURL = baseURL + leafCertificatePath
	leafCertificateKeyURL = baseURL + leafCertificateKeyPath
	intermediateCertificateURL = baseURL + intermediateCertificatePath
}

func (p *instancePrincipalKeyProvider) RegionForFederationClient() common.Region {
	return p.Region
}

func (p *instancePrincipalKeyProvider) PrivateRSAKey() (privateKey *rsa.PrivateKey, err error) {
	if privateKey, err = p.FederationClient.PrivateKey(); err != nil {
		err = fmt.Errorf("failed to get private key: %s", err.Error())
		return nil, err
	}
	return privateKey, nil
}

func (p *instancePrincipalKeyProvider) KeyID() (string, error) {
	var securityToken string
	var err error
	if securityToken, err = p.FederationClient.SecurityToken(); err != nil {
		return "", fmt.Errorf("failed to get security token: %s", err.Error())
	}
	return fmt.Sprintf("ST$%s", securityToken), nil
}

func (p *instancePrincipalKeyProvider) TenancyOCID() (string, error) {
	return p.TenancyID, nil
}
