package api

import (
	"crypto/x509"
	"fmt"
)

type certList struct {
	pki           *Pki
	serialNumbers []SerialNumber
}

type CertLister interface {
	Count() int
	Next() (*x509.Certificate, error)
	All(threads int) ([]*x509.Certificate, error)
	SerialNumbers() []SerialNumber
}

func (pki *Pki) List() (CertLister, error) {
	path := fmt.Sprintf("/v1/%s/certs", pki.path)
	r := pki.c.NewRequest("LIST", path)
	resp, err := pki.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}

	keyInterfaces := secret.Data["keys"].([]interface{})
	serials := make([]SerialNumber, len(keyInterfaces))
	for i, v := range keyInterfaces {
		s := v.(string)
		serials[i] = SerialNumber(s)
	}

	return &certList{pki, serials}, nil
}

func (cl *certList) Count() int {
	return len(cl.serialNumbers)
}

func (cl *certList) SerialNumbers() []SerialNumber {
	return cl.serialNumbers
}
