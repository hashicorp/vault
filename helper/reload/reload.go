package reload

import (
	"crypto/tls"
	"fmt"
	"sync"
)

// ReloadFunc are functions that are called when a reload is requested
type ReloadFunc func(map[string]interface{}) error

// CertificateGetter satisfies ReloadFunc and its GetCertificate method
// satisfies the tls.GetCertificate function signature.  Currently it does not
// allow changing paths after the fact.
type CertificateGetter struct {
	sync.RWMutex

	cert *tls.Certificate

	certFile string
	keyFile  string
}

func NewCertificateGetter(certFile, keyFile string) *CertificateGetter {
	return &CertificateGetter{
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (cg *CertificateGetter) Reload(_ map[string]interface{}) error {
	cert, err := tls.LoadX509KeyPair(cg.certFile, cg.keyFile)
	if err != nil {
		return err
	}

	cg.Lock()
	defer cg.Unlock()

	cg.cert = &cert

	return nil
}

func (cg *CertificateGetter) GetCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cg.RLock()
	defer cg.RUnlock()

	if cg.cert == nil {
		return nil, fmt.Errorf("nil certificate")
	}

	return cg.cert, nil
}
