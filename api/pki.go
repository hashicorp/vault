package api

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

type SerialNumber string

// Pki is used to perform PKI-related operations on Vault.
type Pki struct {
	c    *Client
	path string
}

// Pki is used to return the client for PKI-related API calls.
func (c *Client) Pki(mountPoint string) *Pki {
	return &Pki{c: c, path: mountPoint}
}

func (pki *Pki) Ca() (*x509.Certificate, error) {
	r := pki.c.NewRequest("GET", fmt.Sprintf("/v1/%s/ca", pki.path))
	resp, err := pki.c.RawRequest(r)
	if err != nil {
		return &x509.Certificate{}, err
	}
	defer resp.Body.Close()

	certBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &x509.Certificate{}, err
	}

	cert, err := x509.ParseCertificate(certBody)
	if err != nil {
		return &x509.Certificate{}, err
	}
	return cert, nil
}

func (pki *Pki) Cert(serial SerialNumber) (*x509.Certificate, error) {
	path := fmt.Sprintf("/v1/%s/cert/%s", pki.path, serial)
	r := pki.c.NewRequest("GET", path)
	resp, err := pki.c.RawRequest(r)
	if err != nil {
		return &x509.Certificate{}, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return &x509.Certificate{}, err
	}
	certPem := secret.Data["certificate"]
	certString := certPem.(string)
	certBytes := []byte(certString)

	// We don't need "rest" here because the response contains only one Cert.
	// Also there is no need to check for block.Type, because it will always
	// be a pem "CERTIFICATE" type
	block, _ := pem.Decode(certBytes)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return &x509.Certificate{}, err
	}
	return cert, nil
}
