// +build go1.10

package options

import "crypto/x509"

func x509CertSubject(cert *x509.Certificate) string {
	return cert.Subject.String()
}

func x509MarshalPKCS8PrivateKey(pkcs8 interface{}) ([]byte, error) {
	return x509.MarshalPKCS8PrivateKey(pkcs8)
}
