// +build !go1.10

package topology

import (
	"crypto/x509"
)

// We don't support version less then 1.10, but Evergreen needs to be able to compile the driver
// using version 1.8.
func x509CertSubject(cert *x509.Certificate) string {
	return ""
}
