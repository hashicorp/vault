package reload

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/hashicorp/errwrap"
)

// ReloadFunc are functions that are called when a reload is requested
type ReloadFunc func(map[string]interface{}) error

// CertificateGetter satisfies ReloadFunc and its GetCertificate method
// satisfies the tls.GetCertificate function signature.  Currently it does not
// allow changing paths after the fact.
type CertificateGetter struct {
	sync.RWMutex

	cert *tls.Certificate

	certFile   string
	keyFile    string
	passphrase string
}

func NewCertificateGetter(certFile, keyFile, passphrase string) *CertificateGetter {
	return &CertificateGetter{
		certFile:   certFile,
		keyFile:    keyFile,
		passphrase: passphrase,
	}
}

func (cg *CertificateGetter) Reload(_ map[string]interface{}) error {
	certPEMBlock, err := ioutil.ReadFile(cg.certFile)
	if err != nil {
		return err
	}
	keyPEMBlock, err := ioutil.ReadFile(cg.keyFile)
	if err != nil {
		return err
	}

	// Check for encrypted pem block
	keyBlock, _ := pem.Decode(keyPEMBlock)
	if keyBlock == nil {
		return errors.New("decoded PEM is blank")
	}

	if x509.IsEncryptedPEMBlock(keyBlock) {
		keyBlock.Bytes, err = x509.DecryptPEMBlock(keyBlock, []byte(cg.passphrase))
		if err != nil {
			return errwrap.Wrapf("Decrypting PEM block failed {{err}}", err)
		}
		keyPEMBlock = pem.EncodeToMemory(keyBlock)
	}

	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
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
