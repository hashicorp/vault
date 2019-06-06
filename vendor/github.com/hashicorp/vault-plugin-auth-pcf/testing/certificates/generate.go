package certificates

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
)

// Generate is a convenience method for testing. It creates a group of test certificates with the
// client certificate reflecting the given values. Close() should be called when done to immediately
// delete the three temporary files it has created.
//
// Usage:
//
// testCerts, err := certificates.Generate(...)
// if err != nil {
// 		...
// }
// defer func(){
// 		if err := testCerts.Close(); err != nil {
//			...
// 		}
// }()
//
func Generate(instanceID, orgID, spaceID, appID, ipAddress string) (*TestCertificates, error) {
	caCert, instanceCert, instanceKey, err := generate(instanceID, orgID, spaceID, appID, ipAddress)
	if err != nil {
		return nil, err
	}

	// Keep a list of paths we've created so that if we fail along the way,
	// we can attempt to clean them up.
	var paths []string
	pathToCACertificate, err := makePathTo(caCert)
	if err != nil {
		// No path was successfully created, so we don't need to cleanup here.
		return nil, err
	}
	paths = append(paths, pathToCACertificate)

	pathToInstanceCertificate, err := makePathTo(instanceCert)
	if err != nil {
		if cleanupErr := cleanup(paths); cleanupErr != nil {
			return nil, multierror.Append(err, cleanupErr)
		}
		return nil, err
	}
	paths = append(paths, pathToInstanceCertificate)

	pathToInstanceKey, err := makePathTo(instanceKey)
	if err != nil {
		if cleanupErr := cleanup(paths); cleanupErr != nil {
			return nil, multierror.Append(err, cleanupErr)
		}
		return nil, err
	}
	paths = append(paths, pathToInstanceKey)

	// Provide a function to be called at the end cleaning up our temporary files.
	cleanup := func() error {
		return cleanup(paths)
	}

	return &TestCertificates{
		CACertificate:             caCert,
		InstanceCertificate:       instanceCert,
		InstanceKey:               instanceKey,
		PathToCACertificate:       pathToCACertificate,
		PathToInstanceCertificate: pathToInstanceCertificate,
		PathToInstanceKey:         pathToInstanceKey,
		cleanup:                   cleanup,
	}, nil
}

type TestCertificates struct {
	CACertificate       string
	InstanceCertificate string
	InstanceKey         string

	PathToCACertificate       string
	PathToInstanceCertificate string
	PathToInstanceKey         string

	// cleanup contains a function that has a path to all the temporary files we made,
	// and deletes them. They're all in the /tmp folder so they'll disappear on the next
	// system restart anyways, but in case of repeated tests, it's best to leave nothing
	// behind if possible.
	cleanup func() error
}

func (e *TestCertificates) Close() error {
	return e.cleanup()
}

func generate(instanceID, orgID, spaceID, appID, ipAddress string) (caCert, instanceCert, instanceKey string, err error) {
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", "", err
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Province:     []string{"CA"},
			Organization: []string{"Testing, Inc."},
			CommonName:   "test-CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(caPrivateKey), caPrivateKey)
	if err != nil {
		return "", "", "", err
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	caCert = out.String()
	out.Reset()

	block, certBytes := pem.Decode([]byte(caCert))
	if block == nil {
		return "", "", "", errors.New("block shouldn't be nil")
	}
	if len(certBytes) > 0 {
		return "", "", "", errors.New("there shouldn't be more bytes")
	}
	ca509cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", "", "", err
	}

	template = x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{"US"},
			Province:     []string{"CA"},
			Organization: []string{"Cloud Foundry"},
			OrganizationalUnit: []string{
				fmt.Sprintf("organization:%s", orgID),
				fmt.Sprintf("space:%s", spaceID),
				fmt.Sprintf("app:%s", appID),
			},
			CommonName: instanceID,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		IPAddresses:           []net.IP{net.ParseIP(ipAddress)},
	}

	clientPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", "", err
	}

	derBytes, err = x509.CreateCertificate(rand.Reader, &template, ca509cert, publicKey(clientPrivateKey), caPrivateKey)
	if err != nil {
		return "", "", "", err
	}

	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	instanceCert = out.String()
	out.Reset()

	pem.Encode(out, pemBlockForKey(clientPrivateKey))
	instanceKey = out.String()
	out.Reset()

	return caCert, instanceCert, instanceKey, nil
}

func makePathTo(certOrKey string) (string, error) {
	u, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	tmpFile, err := ioutil.TempFile("", u)
	if err != nil {
		return "", err
	}
	if _, err := tmpFile.Write([]byte(certOrKey)); err != nil {
		return "", err
	}
	if err := tmpFile.Close(); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	default:
		return nil
	}
}

func cleanup(paths []string) error {
	var result error
	for i := 0; i < len(paths); i++ {
		if err := os.Remove(paths[i]); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}
