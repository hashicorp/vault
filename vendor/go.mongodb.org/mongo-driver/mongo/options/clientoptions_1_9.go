// +build !go1.10

package options

import (
	"crypto/x509"
	"fmt"
)

// We don't support Go versions less than 1.10, but Evergreen needs to be able to compile the driver
// using version 1.9 and cert.Subject
func x509CertSubject(cert *x509.Certificate) string {
	return ""
}

// We don't support Go versions less than 1.10, but Evergreen needs to be able to compile the driver
// using version 1.9 and x509.MarshalPKCS8PrivateKey()
func x509MarshalPKCS8PrivateKey(pkcs8 interface{}) ([]byte, error) {
	return nil, fmt.Errorf("PKCS8-encrypted client private keys are only supported with go1.10+")
}
