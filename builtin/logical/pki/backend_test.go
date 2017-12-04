package pki

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	mathrand "math/rand"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/strutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
)

var (
	stepCount               = 0
	serialUnderTest         string
	parsedKeyUsageUnderTest int
)

// Performs basic tests on CA functionality
// Uses the RSA CA key
func TestBackend_RSAKey(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps:   []logicaltest.TestStep{},
	}

	stepCount = len(testCase.Steps)

	intdata := map[string]interface{}{}
	reqdata := map[string]interface{}{}
	testCase.Steps = append(testCase.Steps, generateCATestingSteps(t, rsaCACert, rsaCAKey, ecCACert, intdata, reqdata)...)

	logicaltest.Test(t, testCase)
}

// Performs basic tests on CA functionality
// Uses the EC CA key
func TestBackend_ECKey(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps:   []logicaltest.TestStep{},
	}

	stepCount = len(testCase.Steps)

	intdata := map[string]interface{}{}
	reqdata := map[string]interface{}{}
	testCase.Steps = append(testCase.Steps, generateCATestingSteps(t, ecCACert, ecCAKey, rsaCACert, intdata, reqdata)...)

	logicaltest.Test(t, testCase)
}

func TestBackend_CSRValues(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps:   []logicaltest.TestStep{},
	}

	stepCount = len(testCase.Steps)

	intdata := map[string]interface{}{}
	reqdata := map[string]interface{}{}
	testCase.Steps = append(testCase.Steps, generateCSRSteps(t, ecCACert, ecCAKey, intdata, reqdata)...)

	logicaltest.Test(t, testCase)
}

func TestBackend_URLsCRUD(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps:   []logicaltest.TestStep{},
	}

	stepCount = len(testCase.Steps)

	intdata := map[string]interface{}{}
	reqdata := map[string]interface{}{}
	testCase.Steps = append(testCase.Steps, generateURLSteps(t, ecCACert, ecCAKey, intdata, reqdata)...)

	logicaltest.Test(t, testCase)
}

// Generates and tests steps that walk through the various possibilities
// of role flags to ensure that they are properly restricted
// Uses the RSA CA key
func TestBackend_RSARoles(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "config/ca",
				Data: map[string]interface{}{
					"pem_bundle": rsaCAKey + rsaCACert,
				},
			},
		},
	}

	stepCount = len(testCase.Steps)

	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, false)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+stepCount, v)
		}
	}

	logicaltest.Test(t, testCase)
}

// Generates and tests steps that walk through the various possibilities
// of role flags to ensure that they are properly restricted
// Uses the RSA CA key
func TestBackend_RSARoles_CSR(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "config/ca",
				Data: map[string]interface{}{
					"pem_bundle": rsaCAKey + rsaCACert + rsaCAChain,
				},
			},
		},
	}

	stepCount = len(testCase.Steps)

	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, true)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+stepCount, v)
		}
	}

	logicaltest.Test(t, testCase)
}

// Generates and tests steps that walk through the various possibilities
// of role flags to ensure that they are properly restricted
// Uses the EC CA key
func TestBackend_ECRoles(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "config/ca",
				Data: map[string]interface{}{
					"pem_bundle": ecCAKey + ecCACert,
				},
			},
		},
	}

	stepCount = len(testCase.Steps)

	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, false)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+stepCount, v)
		}
	}

	logicaltest.Test(t, testCase)
}

// Generates and tests steps that walk through the various possibilities
// of role flags to ensure that they are properly restricted
// Uses the EC CA key
func TestBackend_ECRoles_CSR(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "config/ca",
				Data: map[string]interface{}{
					"pem_bundle": ecCAKey + ecCACert,
				},
			},
		},
	}

	stepCount = len(testCase.Steps)

	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, true)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+stepCount, v)
		}
	}

	logicaltest.Test(t, testCase)
}

// Performs some validity checking on the returned bundles
func checkCertsAndPrivateKey(keyType string, key crypto.Signer, usage x509.KeyUsage, extUsage x509.ExtKeyUsage, validity time.Duration, certBundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	parsedCertBundle, err := certBundle.ToParsedCertBundle()
	if err != nil {
		return nil, fmt.Errorf("Error parsing cert bundle: %s", err)
	}

	if key != nil {
		switch keyType {
		case "rsa":
			parsedCertBundle.PrivateKeyType = certutil.RSAPrivateKey
			parsedCertBundle.PrivateKey = key
			parsedCertBundle.PrivateKeyBytes = x509.MarshalPKCS1PrivateKey(key.(*rsa.PrivateKey))
		case "ec":
			parsedCertBundle.PrivateKeyType = certutil.ECPrivateKey
			parsedCertBundle.PrivateKey = key
			parsedCertBundle.PrivateKeyBytes, err = x509.MarshalECPrivateKey(key.(*ecdsa.PrivateKey))
			if err != nil {
				return nil, fmt.Errorf("Error parsing EC key: %s", err)
			}
		}
	}

	switch {
	case parsedCertBundle.Certificate == nil:
		return nil, fmt.Errorf("Did not find a certificate in the cert bundle")
	case len(parsedCertBundle.CAChain) == 0 || parsedCertBundle.CAChain[0].Certificate == nil:
		return nil, fmt.Errorf("Did not find a CA in the cert bundle")
	case parsedCertBundle.PrivateKey == nil:
		return nil, fmt.Errorf("Did not find a private key in the cert bundle")
	case parsedCertBundle.PrivateKeyType == certutil.UnknownPrivateKey:
		return nil, fmt.Errorf("Could not figure out type of private key")
	}

	switch {
	case parsedCertBundle.PrivateKeyType == certutil.RSAPrivateKey && keyType != "rsa":
		fallthrough
	case parsedCertBundle.PrivateKeyType == certutil.ECPrivateKey && keyType != "ec":
		return nil, fmt.Errorf("Given key type does not match type found in bundle")
	}

	cert := parsedCertBundle.Certificate

	if usage != cert.KeyUsage {
		return nil, fmt.Errorf("Expected usage of %#v, got %#v; ext usage is %#v", usage, cert.KeyUsage, cert.ExtKeyUsage)
	}

	// There should only be one ext usage type, because only one is requested
	// in the tests
	if len(cert.ExtKeyUsage) != 1 {
		return nil, fmt.Errorf("Got wrong size key usage in generated cert; expected 1, values are %#v", cert.ExtKeyUsage)
	}
	switch extUsage {
	case x509.ExtKeyUsageEmailProtection:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageEmailProtection {
			return nil, fmt.Errorf("Bad extended key usage")
		}
	case x509.ExtKeyUsageServerAuth:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageServerAuth {
			return nil, fmt.Errorf("Bad extended key usage")
		}
	case x509.ExtKeyUsageClientAuth:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth {
			return nil, fmt.Errorf("Bad extended key usage")
		}
	case x509.ExtKeyUsageCodeSigning:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageCodeSigning {
			return nil, fmt.Errorf("Bad extended key usage")
		}
	}

	// 40 seconds since we add 30 second slack for clock skew
	if math.Abs(float64(time.Now().Unix()-cert.NotBefore.Unix())) > 40 {
		return nil, fmt.Errorf("Validity period starts out of range")
	}
	if !cert.NotBefore.Before(time.Now().Add(-10 * time.Second)) {
		return nil, fmt.Errorf("Validity period not far enough in the past")
	}

	if math.Abs(float64(time.Now().Add(validity).Unix()-cert.NotAfter.Unix())) > 20 {
		return nil, fmt.Errorf("Certificate validity end: %s; expected within 20 seconds of %s", cert.NotAfter.Format(time.RFC3339), time.Now().Add(validity).Format(time.RFC3339))
	}

	return parsedCertBundle, nil
}

func generateURLSteps(t *testing.T, caCert, caKey string, intdata, reqdata map[string]interface{}) []logicaltest.TestStep {
	expected := urlEntries{
		IssuingCertificates: []string{
			"http://example.com/ca1",
			"http://example.com/ca2",
		},
		CRLDistributionPoints: []string{
			"http://example.com/crl1",
			"http://example.com/crl2",
		},
		OCSPServers: []string{
			"http://example.com/ocsp1",
			"http://example.com/ocsp2",
		},
	}
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: "my@example.com",
		},
	}

	priv1024, _ := rsa.GenerateKey(rand.Reader, 1024)
	csr1024, _ := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, priv1024)
	csrPem1024 := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr1024,
	})

	priv2048, _ := rsa.GenerateKey(rand.Reader, 2048)
	csr2048, _ := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, priv2048)
	csrPem2048 := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr2048,
	})

	ret := []logicaltest.TestStep{
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/generate/exported",
			Data: map[string]interface{}{
				"common_name": "Root Cert",
				"ttl":         "180h",
			},
			Check: func(resp *logical.Response) error {
				if resp.Secret != nil && resp.Secret.LeaseID != "" {
					return fmt.Errorf("root returned with a lease")
				}
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "config/urls",
			Data: map[string]interface{}{
				"issuing_certificates":    strings.Join(expected.IssuingCertificates, ","),
				"crl_distribution_points": strings.Join(expected.CRLDistributionPoints, ","),
				"ocsp_servers":            strings.Join(expected.OCSPServers, ","),
			},
		},

		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			Path:      "config/urls",
			Check: func(resp *logical.Response) error {
				if resp.Data == nil {
					return fmt.Errorf("no data returned")
				}
				var entries urlEntries
				err := mapstructure.Decode(resp.Data, &entries)
				if err != nil {
					return err
				}

				if !reflect.DeepEqual(entries, expected) {
					return fmt.Errorf("expected urls\n%#v\ndoes not match provided\n%#v\n", expected, entries)
				}

				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data: map[string]interface{}{
				"common_name": "Intermediate Cert",
				"csr":         string(csrPem1024),
				"format":      "der",
			},
			ErrorOk: true,
			Check: func(resp *logical.Response) error {
				if !resp.IsError() {
					return fmt.Errorf("expected an error response but did not get one")
				}
				if !strings.Contains(resp.Data["error"].(string), "2048") {
					return fmt.Errorf("recieved an error but not about a 1024-bit key, error was: %s", resp.Data["error"].(string))
				}

				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data: map[string]interface{}{
				"common_name": "Intermediate Cert",
				"csr":         string(csrPem2048),
				"format":      "der",
			},
			Check: func(resp *logical.Response) error {
				certString := resp.Data["certificate"].(string)
				if certString == "" {
					return fmt.Errorf("no certificate returned")
				}
				if resp.Secret != nil && resp.Secret.LeaseID != "" {
					return fmt.Errorf("signed intermediate returned with a lease")
				}
				certBytes, _ := base64.StdEncoding.DecodeString(certString)
				certs, err := x509.ParseCertificates(certBytes)
				if err != nil {
					return fmt.Errorf("returned cert cannot be parsed: %v", err)
				}
				if len(certs) != 1 {
					return fmt.Errorf("unexpected returned length of certificates: %d", len(certs))
				}
				cert := certs[0]

				switch {
				case !reflect.DeepEqual(expected.IssuingCertificates, cert.IssuingCertificateURL):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", expected.IssuingCertificates, cert.IssuingCertificateURL)
				case !reflect.DeepEqual(expected.CRLDistributionPoints, cert.CRLDistributionPoints):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", expected.CRLDistributionPoints, cert.CRLDistributionPoints)
				case !reflect.DeepEqual(expected.OCSPServers, cert.OCSPServer):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", expected.OCSPServers, cert.OCSPServer)
				case !reflect.DeepEqual([]string{"Intermediate Cert"}, cert.DNSNames):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", []string{"Intermediate Cert"}, cert.DNSNames)
				}

				return nil
			},
		},

		// Same as above but exclude adding to sans
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data: map[string]interface{}{
				"common_name":          "Intermediate Cert",
				"csr":                  string(csrPem2048),
				"format":               "der",
				"exclude_cn_from_sans": true,
			},
			Check: func(resp *logical.Response) error {
				certString := resp.Data["certificate"].(string)
				if certString == "" {
					return fmt.Errorf("no certificate returned")
				}
				if resp.Secret != nil && resp.Secret.LeaseID != "" {
					return fmt.Errorf("signed intermediate returned with a lease")
				}
				certBytes, _ := base64.StdEncoding.DecodeString(certString)
				certs, err := x509.ParseCertificates(certBytes)
				if err != nil {
					return fmt.Errorf("returned cert cannot be parsed: %v", err)
				}
				if len(certs) != 1 {
					return fmt.Errorf("unexpected returned length of certificates: %d", len(certs))
				}
				cert := certs[0]

				switch {
				case !reflect.DeepEqual(expected.IssuingCertificates, cert.IssuingCertificateURL):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", expected.IssuingCertificates, cert.IssuingCertificateURL)
				case !reflect.DeepEqual(expected.CRLDistributionPoints, cert.CRLDistributionPoints):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", expected.CRLDistributionPoints, cert.CRLDistributionPoints)
				case !reflect.DeepEqual(expected.OCSPServers, cert.OCSPServer):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", expected.OCSPServers, cert.OCSPServer)
				case !reflect.DeepEqual([]string(nil), cert.DNSNames):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", []string(nil), cert.DNSNames)
				}

				return nil
			},
		},
	}
	return ret
}

func generateCSRSteps(t *testing.T, caCert, caKey string, intdata, reqdata map[string]interface{}) []logicaltest.TestStep {
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			Country:      []string{"MyCountry"},
			PostalCode:   []string{"MyPostalCode"},
			SerialNumber: "MySerialNumber",
			CommonName:   "my@example.com",
		},
		DNSNames: []string{
			"name1.example.com",
			"name2.example.com",
			"name3.example.com",
		},
		EmailAddresses: []string{
			"name1@example.com",
			"name2@example.com",
			"name3@example.com",
		},
		IPAddresses: []net.IP{
			net.ParseIP("::ff:1:2:3:4"),
			net.ParseIP("::ff:5:6:7:8"),
		},
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	csr, _ := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, priv)
	csrPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	})

	ret := []logicaltest.TestStep{
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/generate/exported",
			Data: map[string]interface{}{
				"common_name":     "Root Cert",
				"ttl":             "180h",
				"max_path_length": 0,
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data: map[string]interface{}{
				"use_csr_values": true,
				"csr":            string(csrPem),
				"format":         "der",
			},
			ErrorOk: true,
		},

		logicaltest.TestStep{
			Operation: logical.DeleteOperation,
			Path:      "root",
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/generate/exported",
			Data: map[string]interface{}{
				"common_name":     "Root Cert",
				"ttl":             "180h",
				"max_path_length": 1,
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data: map[string]interface{}{
				"use_csr_values": true,
				"csr":            string(csrPem),
				"format":         "der",
			},
			Check: func(resp *logical.Response) error {
				certString := resp.Data["certificate"].(string)
				if certString == "" {
					return fmt.Errorf("no certificate returned")
				}
				certBytes, _ := base64.StdEncoding.DecodeString(certString)
				certs, err := x509.ParseCertificates(certBytes)
				if err != nil {
					return fmt.Errorf("returned cert cannot be parsed: %v", err)
				}
				if len(certs) != 1 {
					return fmt.Errorf("unexpected returned length of certificates: %d", len(certs))
				}
				cert := certs[0]

				if cert.MaxPathLen != 0 {
					return fmt.Errorf("max path length of %d does not match the requested of 3", cert.MaxPathLen)
				}
				if !cert.MaxPathLenZero {
					return fmt.Errorf("max path length zero is not set")
				}

				// We need to set these as they are filled in with unparsed values in the final cert
				csrTemplate.Subject.Names = cert.Subject.Names
				csrTemplate.Subject.ExtraNames = cert.Subject.ExtraNames

				switch {
				case !reflect.DeepEqual(cert.Subject, csrTemplate.Subject):
					return fmt.Errorf("cert subject\n%#v\ndoes not match csr subject\n%#v\n", cert.Subject, csrTemplate.Subject)
				case !reflect.DeepEqual(cert.DNSNames, csrTemplate.DNSNames):
					return fmt.Errorf("cert dns names\n%#v\ndoes not match csr dns names\n%#v\n", cert.DNSNames, csrTemplate.DNSNames)
				case !reflect.DeepEqual(cert.EmailAddresses, csrTemplate.EmailAddresses):
					return fmt.Errorf("cert email addresses\n%#v\ndoes not match csr email addresses\n%#v\n", cert.EmailAddresses, csrTemplate.EmailAddresses)
				case !reflect.DeepEqual(cert.IPAddresses, csrTemplate.IPAddresses):
					return fmt.Errorf("cert ip addresses\n%#v\ndoes not match csr ip addresses\n%#v\n", cert.IPAddresses, csrTemplate.IPAddresses)
				}
				return nil
			},
		},
	}
	return ret
}

// Generates steps to test out CA configuration -- certificates + CRL expiry,
// and ensure that the certificates are readable after storing them
func generateCATestingSteps(t *testing.T, caCert, caKey, otherCaCert string, intdata, reqdata map[string]interface{}) []logicaltest.TestStep {
	setSerialUnderTest := func(req *logical.Request) error {
		req.Path = serialUnderTest
		return nil
	}

	ret := []logicaltest.TestStep{
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "config/ca",
			Data: map[string]interface{}{
				"pem_bundle": caKey + caCert,
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "config/crl",
			Data: map[string]interface{}{
				"expiry": "16h",
			},
		},

		// Ensure we can fetch it back via unauthenticated means, in various formats
		logicaltest.TestStep{
			Operation:       logical.ReadOperation,
			Path:            "cert/ca",
			Unauthenticated: true,
			Check: func(resp *logical.Response) error {
				if resp.Data["certificate"].(string) != caCert {
					return fmt.Errorf("CA certificate:\n%s\ndoes not match original:\n%s\n", resp.Data["certificate"].(string), caCert)
				}
				return nil
			},
		},

		logicaltest.TestStep{
			Operation:       logical.ReadOperation,
			Path:            "ca/pem",
			Unauthenticated: true,
			Check: func(resp *logical.Response) error {
				rawBytes := resp.Data["http_raw_body"].([]byte)
				if string(rawBytes) != caCert {
					return fmt.Errorf("CA certificate:\n%s\ndoes not match original:\n%s\n", string(rawBytes), caCert)
				}
				if resp.Data["http_content_type"].(string) != "application/pkix-cert" {
					return fmt.Errorf("Expected application/pkix-cert as content-type, but got %s", resp.Data["http_content_type"].(string))
				}
				return nil
			},
		},

		logicaltest.TestStep{
			Operation:       logical.ReadOperation,
			Path:            "ca",
			Unauthenticated: true,
			Check: func(resp *logical.Response) error {
				rawBytes := resp.Data["http_raw_body"].([]byte)
				pemBytes := pem.EncodeToMemory(&pem.Block{
					Type:  "CERTIFICATE",
					Bytes: rawBytes,
				})
				if string(pemBytes) != caCert {
					return fmt.Errorf("CA certificate:\n%s\ndoes not match original:\n%s\n", string(pemBytes), caCert)
				}
				if resp.Data["http_content_type"].(string) != "application/pkix-cert" {
					return fmt.Errorf("Expected application/pkix-cert as content-type, but got %s", resp.Data["http_content_type"].(string))
				}
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			Path:      "config/crl",
			Check: func(resp *logical.Response) error {
				if resp.Data["expiry"].(string) != "16h" {
					return fmt.Errorf("CRL lifetimes do not match (got %s)", resp.Data["expiry"].(string))
				}
				return nil
			},
		},

		// Ensure that both parts of the PEM bundle are required
		// Here, just the cert
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "config/ca",
			Data: map[string]interface{}{
				"pem_bundle": caCert,
			},
			ErrorOk: true,
		},

		// Here, just the key
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "config/ca",
			Data: map[string]interface{}{
				"pem_bundle": caKey,
			},
			ErrorOk: true,
		},

		// Ensure we can fetch it back via unauthenticated means, in various formats
		logicaltest.TestStep{
			Operation:       logical.ReadOperation,
			Path:            "cert/ca",
			Unauthenticated: true,
			Check: func(resp *logical.Response) error {
				if resp.Data["certificate"].(string) != caCert {
					return fmt.Errorf("CA certificate:\n%s\ndoes not match original:\n%s\n", resp.Data["certificate"].(string), caCert)
				}
				return nil
			},
		},

		logicaltest.TestStep{
			Operation:       logical.ReadOperation,
			Path:            "ca/pem",
			Unauthenticated: true,
			Check: func(resp *logical.Response) error {
				rawBytes := resp.Data["http_raw_body"].([]byte)
				if string(rawBytes) != caCert {
					return fmt.Errorf("CA certificate:\n%s\ndoes not match original:\n%s\n", string(rawBytes), caCert)
				}
				if resp.Data["http_content_type"].(string) != "application/pkix-cert" {
					return fmt.Errorf("Expected application/pkix-cert as content-type, but got %s", resp.Data["http_content_type"].(string))
				}
				return nil
			},
		},

		logicaltest.TestStep{
			Operation:       logical.ReadOperation,
			Path:            "ca",
			Unauthenticated: true,
			Check: func(resp *logical.Response) error {
				rawBytes := resp.Data["http_raw_body"].([]byte)
				pemBytes := pem.EncodeToMemory(&pem.Block{
					Type:  "CERTIFICATE",
					Bytes: rawBytes,
				})
				if string(pemBytes) != caCert {
					return fmt.Errorf("CA certificate:\n%s\ndoes not match original:\n%s\n", string(pemBytes), caCert)
				}
				if resp.Data["http_content_type"].(string) != "application/pkix-cert" {
					return fmt.Errorf("Expected application/pkix-cert as content-type, but got %s", resp.Data["http_content_type"].(string))
				}
				return nil
			},
		},

		// Test a bunch of generation stuff
		logicaltest.TestStep{
			Operation: logical.DeleteOperation,
			Path:      "root",
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/generate/exported",
			Data: map[string]interface{}{
				"common_name": "Root Cert",
				"ttl":         "180h",
			},
			Check: func(resp *logical.Response) error {
				intdata["root"] = resp.Data["certificate"].(string)
				intdata["rootkey"] = resp.Data["private_key"].(string)
				reqdata["pem_bundle"] = intdata["root"].(string) + "\n" + intdata["rootkey"].(string)
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "intermediate/generate/exported",
			Data: map[string]interface{}{
				"common_name": "Intermediate Cert",
			},
			Check: func(resp *logical.Response) error {
				intdata["intermediatecsr"] = resp.Data["csr"].(string)
				intdata["intermediatekey"] = resp.Data["private_key"].(string)
				return nil
			},
		},

		// Re-load the root key in so we can sign it
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "config/ca",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				delete(reqdata, "pem_bundle")
				delete(reqdata, "ttl")
				reqdata["csr"] = intdata["intermediatecsr"].(string)
				reqdata["common_name"] = "Intermediate Cert"
				reqdata["ttl"] = "10s"
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				delete(reqdata, "csr")
				delete(reqdata, "common_name")
				delete(reqdata, "ttl")
				intdata["intermediatecert"] = resp.Data["certificate"].(string)
				reqdata["serial_number"] = resp.Data["serial_number"].(string)
				reqdata["rsa_int_serial_number"] = resp.Data["serial_number"].(string)
				reqdata["certificate"] = resp.Data["certificate"].(string)
				reqdata["pem_bundle"] = intdata["intermediatekey"].(string) + "\n" + resp.Data["certificate"].(string)
				return nil
			},
		},

		// First load in this way to populate the private key
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "config/ca",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				delete(reqdata, "pem_bundle")
				return nil
			},
		},

		// Now test setting the intermediate, signed CA cert
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "intermediate/set-signed",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				delete(reqdata, "certificate")

				serialUnderTest = "cert/" + reqdata["rsa_int_serial_number"].(string)

				return nil
			},
		},

		// We expect to find a zero revocation time
		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			PreFlight: setSerialUnderTest,
			Check: func(resp *logical.Response) error {
				if resp.Data["error"] != nil && resp.Data["error"].(string) != "" {
					return fmt.Errorf("got an error: %s", resp.Data["error"].(string))
				}

				if resp.Data["revocation_time"].(int64) != 0 {
					return fmt.Errorf("expected a zero revocation time")
				}

				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "revoke",
			Data:      reqdata,
		},

		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			Path:      "crl",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				crlBytes := resp.Data["http_raw_body"].([]byte)
				certList, err := x509.ParseCRL(crlBytes)
				if err != nil {
					t.Fatalf("err: %s", err)
				}
				revokedList := certList.TBSCertList.RevokedCertificates
				if len(revokedList) != 1 {
					t.Fatalf("length of revoked list not 1; %d", len(revokedList))
				}
				revokedString := certutil.GetHexFormatted(revokedList[0].SerialNumber.Bytes(), ":")
				if revokedString != reqdata["serial_number"].(string) {
					t.Fatalf("got serial %s, expecting %s", revokedString, reqdata["serial_number"].(string))
				}
				delete(reqdata, "serial_number")
				return nil
			},
		},

		// Do it all again, with EC keys and DER format
		logicaltest.TestStep{
			Operation: logical.DeleteOperation,
			Path:      "root",
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/generate/exported",
			Data: map[string]interface{}{
				"common_name": "Root Cert",
				"ttl":         "180h",
				"key_type":    "ec",
				"key_bits":    384,
				"format":      "der",
			},
			Check: func(resp *logical.Response) error {
				certBytes, _ := base64.StdEncoding.DecodeString(resp.Data["certificate"].(string))
				certPem := pem.EncodeToMemory(&pem.Block{
					Type:  "CERTIFICATE",
					Bytes: certBytes,
				})
				keyBytes, _ := base64.StdEncoding.DecodeString(resp.Data["private_key"].(string))
				keyPem := pem.EncodeToMemory(&pem.Block{
					Type:  "EC PRIVATE KEY",
					Bytes: keyBytes,
				})
				intdata["root"] = string(certPem)
				intdata["rootkey"] = string(keyPem)
				reqdata["pem_bundle"] = string(certPem) + "\n" + string(keyPem)
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "intermediate/generate/exported",
			Data: map[string]interface{}{
				"format":      "der",
				"key_type":    "ec",
				"key_bits":    384,
				"common_name": "Intermediate Cert",
			},
			Check: func(resp *logical.Response) error {
				csrBytes, _ := base64.StdEncoding.DecodeString(resp.Data["csr"].(string))
				csrPem := pem.EncodeToMemory(&pem.Block{
					Type:  "CERTIFICATE REQUEST",
					Bytes: csrBytes,
				})
				keyBytes, _ := base64.StdEncoding.DecodeString(resp.Data["private_key"].(string))
				keyPem := pem.EncodeToMemory(&pem.Block{
					Type:  "EC PRIVATE KEY",
					Bytes: keyBytes,
				})
				intdata["intermediatecsr"] = string(csrPem)
				intdata["intermediatekey"] = string(keyPem)
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "config/ca",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				delete(reqdata, "pem_bundle")
				delete(reqdata, "ttl")
				reqdata["csr"] = intdata["intermediatecsr"].(string)
				reqdata["common_name"] = "Intermediate Cert"
				reqdata["ttl"] = "10s"
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				delete(reqdata, "csr")
				delete(reqdata, "common_name")
				delete(reqdata, "ttl")
				intdata["intermediatecert"] = resp.Data["certificate"].(string)
				reqdata["serial_number"] = resp.Data["serial_number"].(string)
				reqdata["ec_int_serial_number"] = resp.Data["serial_number"].(string)
				reqdata["certificate"] = resp.Data["certificate"].(string)
				reqdata["pem_bundle"] = intdata["intermediatekey"].(string) + "\n" + resp.Data["certificate"].(string)
				return nil
			},
		},

		// First load in this way to populate the private key
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "config/ca",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				delete(reqdata, "pem_bundle")
				return nil
			},
		},

		// Now test setting the intermediate, signed CA cert
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "intermediate/set-signed",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				delete(reqdata, "certificate")

				serialUnderTest = "cert/" + reqdata["ec_int_serial_number"].(string)

				return nil
			},
		},

		// We expect to find a zero revocation time
		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			PreFlight: setSerialUnderTest,
			Check: func(resp *logical.Response) error {
				if resp.Data["error"] != nil && resp.Data["error"].(string) != "" {
					return fmt.Errorf("got an error: %s", resp.Data["error"].(string))
				}

				if resp.Data["revocation_time"].(int64) != 0 {
					return fmt.Errorf("expected a zero revocation time")
				}

				return nil
			},
		},
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "revoke",
			Data:      reqdata,
		},

		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			Path:      "crl",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				crlBytes := resp.Data["http_raw_body"].([]byte)
				certList, err := x509.ParseCRL(crlBytes)
				if err != nil {
					t.Fatalf("err: %s", err)
				}
				revokedList := certList.TBSCertList.RevokedCertificates
				if len(revokedList) != 2 {
					t.Fatalf("length of revoked list not 2; %d", len(revokedList))
				}
				found := false
				for _, revEntry := range revokedList {
					revokedString := certutil.GetHexFormatted(revEntry.SerialNumber.Bytes(), ":")
					if revokedString == reqdata["serial_number"].(string) {
						found = true
					}
				}
				if !found {
					t.Fatalf("did not find %s in CRL", reqdata["serial_number"].(string))
				}
				delete(reqdata, "serial_number")

				serialUnderTest = "cert/" + reqdata["rsa_int_serial_number"].(string)

				return nil
			},
		},

		// Make sure both serial numbers we expect to find are found
		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			PreFlight: setSerialUnderTest,
			Check: func(resp *logical.Response) error {
				if resp.Data["error"] != nil && resp.Data["error"].(string) != "" {
					return fmt.Errorf("got an error: %s", resp.Data["error"].(string))
				}

				if resp.Data["revocation_time"].(int64) == 0 {
					return fmt.Errorf("expected a non-zero revocation time")
				}

				serialUnderTest = "cert/" + reqdata["ec_int_serial_number"].(string)

				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			PreFlight: setSerialUnderTest,
			Check: func(resp *logical.Response) error {
				if resp.Data["error"] != nil && resp.Data["error"].(string) != "" {
					return fmt.Errorf("got an error: %s", resp.Data["error"].(string))
				}

				if resp.Data["revocation_time"].(int64) == 0 {
					return fmt.Errorf("expected a non-zero revocation time")
				}

				// Give time for the certificates to pass the safety buffer
				t.Logf("Sleeping for 15 seconds to allow safety buffer time to pass before testing tidying")
				time.Sleep(15 * time.Second)

				serialUnderTest = "cert/" + reqdata["rsa_int_serial_number"].(string)

				return nil
			},
		},

		// This shouldn't do anything since the safety buffer is too long
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "tidy",
			Data: map[string]interface{}{
				"safety_buffer":        "3h",
				"tidy_cert_store":      true,
				"tidy_revocation_list": true,
			},
		},

		// We still expect to find these
		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			PreFlight: setSerialUnderTest,
			Check: func(resp *logical.Response) error {
				if resp != nil && resp.Data["error"] != nil && resp.Data["error"].(string) != "" {
					return fmt.Errorf("got an error: %s", resp.Data["error"].(string))
				}

				serialUnderTest = "cert/" + reqdata["ec_int_serial_number"].(string)

				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			PreFlight: setSerialUnderTest,
			Check: func(resp *logical.Response) error {
				if resp != nil && resp.Data["error"] != nil && resp.Data["error"].(string) != "" {
					return fmt.Errorf("got an error: %s", resp.Data["error"].(string))
				}

				serialUnderTest = "cert/" + reqdata["rsa_int_serial_number"].(string)

				return nil
			},
		},

		// Both should appear in the CRL
		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			Path:      "crl",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				crlBytes := resp.Data["http_raw_body"].([]byte)
				certList, err := x509.ParseCRL(crlBytes)
				if err != nil {
					t.Fatalf("err: %s", err)
				}
				revokedList := certList.TBSCertList.RevokedCertificates
				if len(revokedList) != 2 {
					t.Fatalf("length of revoked list not 2; %d", len(revokedList))
				}
				foundRsa := false
				foundEc := false
				for _, revEntry := range revokedList {
					revokedString := certutil.GetHexFormatted(revEntry.SerialNumber.Bytes(), ":")
					if revokedString == reqdata["rsa_int_serial_number"].(string) {
						foundRsa = true
					}
					if revokedString == reqdata["ec_int_serial_number"].(string) {
						foundEc = true
					}
				}
				if !foundRsa || !foundEc {
					t.Fatalf("did not find an expected entry in CRL")
				}

				return nil
			},
		},

		// This shouldn't do anything since the boolean values default to false
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "tidy",
			Data: map[string]interface{}{
				"safety_buffer": "1s",
			},
		},

		// We still expect to find these
		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			PreFlight: setSerialUnderTest,
			Check: func(resp *logical.Response) error {
				if resp != nil && resp.Data["error"] != nil && resp.Data["error"].(string) != "" {
					return fmt.Errorf("got an error: %s", resp.Data["error"].(string))
				}

				serialUnderTest = "cert/" + reqdata["ec_int_serial_number"].(string)

				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			PreFlight: setSerialUnderTest,
			Check: func(resp *logical.Response) error {
				if resp != nil && resp.Data["error"] != nil && resp.Data["error"].(string) != "" {
					return fmt.Errorf("got an error: %s", resp.Data["error"].(string))
				}

				serialUnderTest = "cert/" + reqdata["rsa_int_serial_number"].(string)

				return nil
			},
		},

		// This should remove the values since the safety buffer is short
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "tidy",
			Data: map[string]interface{}{
				"safety_buffer":        "1s",
				"tidy_cert_store":      true,
				"tidy_revocation_list": true,
			},
		},

		// We do *not* expect to find these
		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			PreFlight: setSerialUnderTest,
			Check: func(resp *logical.Response) error {
				if resp != nil {
					return fmt.Errorf("expected no response")
				}

				serialUnderTest = "cert/" + reqdata["ec_int_serial_number"].(string)

				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			PreFlight: setSerialUnderTest,
			Check: func(resp *logical.Response) error {
				if resp != nil {
					return fmt.Errorf("expected no response")
				}

				serialUnderTest = "cert/" + reqdata["rsa_int_serial_number"].(string)

				return nil
			},
		},

		// Both should be gone from the CRL
		logicaltest.TestStep{
			Operation: logical.ReadOperation,
			Path:      "crl",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				crlBytes := resp.Data["http_raw_body"].([]byte)
				certList, err := x509.ParseCRL(crlBytes)
				if err != nil {
					t.Fatalf("err: %s", err)
				}
				revokedList := certList.TBSCertList.RevokedCertificates
				if len(revokedList) != 0 {
					t.Fatalf("length of revoked list not 0; %d", len(revokedList))
				}

				return nil
			},
		},
	}

	return ret
}

// Generates steps to test out various role permutations
func generateRoleSteps(t *testing.T, useCSRs bool) []logicaltest.TestStep {
	roleVals := roleEntry{
		MaxTTL:  "12h",
		KeyType: "rsa",
		KeyBits: 2048,
	}
	issueVals := certutil.IssueData{}
	ret := []logicaltest.TestStep{}

	roleTestStep := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/test",
	}
	var issueTestStep logicaltest.TestStep
	if useCSRs {
		issueTestStep = logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "sign/test",
		}
	} else {
		issueTestStep = logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "issue/test",
		}
	}

	generatedRSAKeys := map[int]crypto.Signer{}
	generatedECKeys := map[int]crypto.Signer{}

	/*
		// For the number of tests being run, a seed of 1 has been tested
		// to hit all of the various values below. However, for normal
		// testing we use a randomized time for maximum fuzziness.
	*/
	var seed int64 = 1
	fixedSeed := os.Getenv("VAULT_PKITESTS_FIXED_SEED")
	if len(fixedSeed) == 0 {
		seed = time.Now().UnixNano()
	} else {
		var err error
		seed, err = strconv.ParseInt(fixedSeed, 10, 64)
		if err != nil {
			t.Fatalf("error parsing fixed seed of %s: %v", fixedSeed, err)
		}
	}
	mathRand := mathrand.New(mathrand.NewSource(seed))
	t.Logf("seed under test: %v", seed)

	// Used by tests not toggling common names to turn off the behavior of random key bit fuzziness
	keybitSizeRandOff := false

	genericErrorOkCheck := func(resp *logical.Response) error {
		if resp.IsError() {
			return nil
		}
		return fmt.Errorf("Expected an error, but did not seem to get one")
	}

	// Adds tests with the currently configured issue/role information
	addTests := func(testCheck logicaltest.TestCheckFunc) {
		stepCount++
		//t.Logf("test step %d\nrole vals: %#v\n", stepCount, roleVals)
		stepCount++
		//t.Logf("test step %d\nissue vals: %#v\n", stepCount, issueTestStep)
		roleTestStep.Data = structs.New(roleVals).Map()
		roleTestStep.Data["generate_lease"] = false
		ret = append(ret, roleTestStep)
		issueTestStep.Data = structs.New(issueVals).Map()
		switch {
		case issueTestStep.ErrorOk:
			issueTestStep.Check = genericErrorOkCheck
		case testCheck != nil:
			issueTestStep.Check = testCheck
		default:
			issueTestStep.Check = nil
		}
		ret = append(ret, issueTestStep)
	}

	getOuCheck := func(role roleEntry) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return fmt.Errorf("Error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate

			expected := strutil.ParseDedupLowercaseAndSortStrings(role.OU, ",")
			if !reflect.DeepEqual(cert.Subject.OrganizationalUnit, expected) {
				return fmt.Errorf("Error: returned certificate has OU of %s but %s was specified in the role.", cert.Subject.OrganizationalUnit, expected)
			}
			return nil
		}
	}

	getOrganizationCheck := func(role roleEntry) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return fmt.Errorf("Error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate

			expected := strutil.ParseDedupLowercaseAndSortStrings(role.Organization, ",")
			if !reflect.DeepEqual(cert.Subject.Organization, expected) {
				return fmt.Errorf("Error: returned certificate has Organization of %s but %s was specified in the role.", cert.Subject.Organization, expected)
			}
			return nil
		}
	}

	// Returns a TestCheckFunc that performs various validity checks on the
	// returned certificate information, mostly within checkCertsAndPrivateKey
	getCnCheck := func(name string, role roleEntry, key crypto.Signer, usage x509.KeyUsage, extUsage x509.ExtKeyUsage, validity time.Duration) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			parsedCertBundle, err := checkCertsAndPrivateKey(role.KeyType, key, usage, extUsage, validity, &certBundle)
			if err != nil {
				return fmt.Errorf("Error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate
			if cert.Subject.CommonName != name {
				return fmt.Errorf("Error: returned certificate has CN of %s but %s was requested", cert.Subject.CommonName, name)
			}
			if strings.Contains(cert.Subject.CommonName, "@") {
				if len(cert.DNSNames) != 0 || len(cert.EmailAddresses) != 1 {
					return fmt.Errorf("Error: found more than one DNS SAN or not one Email SAN but only one was requested, cert.DNSNames = %#v, cert.EmailAddresses = %#v", cert.DNSNames, cert.EmailAddresses)
				}
			} else {
				if len(cert.DNSNames) != 1 || len(cert.EmailAddresses) != 0 {
					return fmt.Errorf("Error: found more than one Email SAN or not one DNS SAN but only one was requested, cert.DNSNames = %#v, cert.EmailAddresses = %#v", cert.DNSNames, cert.EmailAddresses)
				}
			}
			var retName string
			if len(cert.DNSNames) > 0 {
				retName = cert.DNSNames[0]
			}
			if len(cert.EmailAddresses) > 0 {
				retName = cert.EmailAddresses[0]
			}
			if retName != name {
				return fmt.Errorf("Error: returned certificate has a DNS SAN of %s but %s was requested", retName, name)
			}
			return nil
		}
	}

	// Common names to test with the various role flags toggled
	var commonNames struct {
		Localhost            bool `structs:"localhost"`
		BareDomain           bool `structs:"example.com"`
		SecondDomain         bool `structs:"foobar.com"`
		SubDomain            bool `structs:"foo.example.com"`
		Wildcard             bool `structs:"*.example.com"`
		SubSubdomain         bool `structs:"foo.bar.example.com"`
		SubSubdomainWildcard bool `structs:"*.bar.example.com"`
		GlobDomain           bool `structs:"fooexample.com"`
		NonHostname          bool `structs:"daɪˈɛrɨsɨs"`
		AnyHost              bool `structs:"porkslap.beer"`
	}

	// Adds a series of tests based on the current selection of
	// allowed common names; contains some (seeded) randomness
	//
	// This allows for a variety of common names to be tested in various
	// combinations with allowed toggles of the role
	addCnTests := func() {
		cnMap := structs.New(commonNames).Map()
		for name, allowedInt := range cnMap {
			roleVals.KeyType = "rsa"
			roleVals.KeyBits = 2048
			if mathRand.Int()%2 == 1 {
				roleVals.KeyType = "ec"
				roleVals.KeyBits = 224
			}

			roleVals.ServerFlag = false
			roleVals.ClientFlag = false
			roleVals.CodeSigningFlag = false
			roleVals.EmailProtectionFlag = false

			var usage string
			if mathRand.Int()%2 == 1 {
				usage = usage + ",DigitalSignature"
			}
			if mathRand.Int()%2 == 1 {
				usage = usage + ",ContentCoMmitment"
			}
			if mathRand.Int()%2 == 1 {
				usage = usage + ",KeyEncipherment"
			}
			if mathRand.Int()%2 == 1 {
				usage = usage + ",DataEncipherment"
			}
			if mathRand.Int()%2 == 1 {
				usage = usage + ",KeyAgreemEnt"
			}
			if mathRand.Int()%2 == 1 {
				usage = usage + ",CertSign"
			}
			if mathRand.Int()%2 == 1 {
				usage = usage + ",CRLSign"
			}
			if mathRand.Int()%2 == 1 {
				usage = usage + ",EncipherOnly"
			}
			if mathRand.Int()%2 == 1 {
				usage = usage + ",DecipherOnly"
			}

			roleVals.KeyUsage = usage
			parsedKeyUsage := parseKeyUsages(roleVals.KeyUsage)
			if parsedKeyUsage == 0 && usage != "" {
				panic("parsed key usages was zero")
			}
			parsedKeyUsageUnderTest = parsedKeyUsage

			var extUsage x509.ExtKeyUsage
			i := mathRand.Int() % 4
			switch {
			case i == 0:
				extUsage = x509.ExtKeyUsageEmailProtection
				roleVals.EmailProtectionFlag = true
			case i == 1:
				extUsage = x509.ExtKeyUsageServerAuth
				roleVals.ServerFlag = true
			case i == 2:
				extUsage = x509.ExtKeyUsageClientAuth
				roleVals.ClientFlag = true
			default:
				extUsage = x509.ExtKeyUsageCodeSigning
				roleVals.CodeSigningFlag = true
			}

			allowed := allowedInt.(bool)
			issueVals.CommonName = name
			if roleVals.EmailProtectionFlag {
				if !strings.HasPrefix(name, "*") {
					issueVals.CommonName = "user@" + issueVals.CommonName
				}
			}

			issueTestStep.ErrorOk = !allowed

			validity, _ := time.ParseDuration(roleVals.MaxTTL)

			var testBitSize int

			if useCSRs {
				rsaKeyBits := []int{2048, 4096}
				ecKeyBits := []int{224, 256, 384, 521}

				var privKey crypto.Signer
				var ok bool
				switch roleVals.KeyType {
				case "rsa":
					roleVals.KeyBits = rsaKeyBits[mathRand.Int()%2]

					// If we don't expect an error already, randomly choose a
					// key size and expect an error if it's less than the role
					// setting
					testBitSize = roleVals.KeyBits
					if !keybitSizeRandOff && !issueTestStep.ErrorOk {
						testBitSize = rsaKeyBits[mathRand.Int()%2]
					}

					if testBitSize < roleVals.KeyBits {
						issueTestStep.ErrorOk = true
					}

					privKey, ok = generatedRSAKeys[testBitSize]
					if !ok {
						privKey, _ = rsa.GenerateKey(rand.Reader, testBitSize)
						generatedRSAKeys[testBitSize] = privKey
					}

				case "ec":
					roleVals.KeyBits = ecKeyBits[mathRand.Int()%4]

					var curve elliptic.Curve

					// If we don't expect an error already, randomly choose a
					// key size and expect an error if it's less than the role
					// setting
					testBitSize = roleVals.KeyBits
					if !keybitSizeRandOff && !issueTestStep.ErrorOk {
						testBitSize = ecKeyBits[mathRand.Int()%4]
					}

					switch testBitSize {
					case 224:
						curve = elliptic.P224()
					case 256:
						curve = elliptic.P256()
					case 384:
						curve = elliptic.P384()
					case 521:
						curve = elliptic.P521()
					}

					if curve.Params().BitSize < roleVals.KeyBits {
						issueTestStep.ErrorOk = true
					}

					privKey, ok = generatedECKeys[testBitSize]
					if !ok {
						privKey, _ = ecdsa.GenerateKey(curve, rand.Reader)
						generatedECKeys[testBitSize] = privKey
					}
				}
				templ := &x509.CertificateRequest{
					Subject: pkix.Name{
						CommonName: issueVals.CommonName,
					},
				}
				csr, err := x509.CreateCertificateRequest(rand.Reader, templ, privKey)
				if err != nil {
					t.Fatalf("Error creating certificate request: %s", err)
				}
				block := pem.Block{
					Type:  "CERTIFICATE REQUEST",
					Bytes: csr,
				}
				issueVals.CSR = strings.TrimSpace(string(pem.EncodeToMemory(&block)))

				addTests(getCnCheck(issueVals.CommonName, roleVals, privKey, x509.KeyUsage(parsedKeyUsage), extUsage, validity))
			} else {
				addTests(getCnCheck(issueVals.CommonName, roleVals, nil, x509.KeyUsage(parsedKeyUsage), extUsage, validity))
			}
		}
	}

	// Common Name tests
	{
		// common_name not provided
		issueVals.CommonName = ""
		issueTestStep.ErrorOk = true
		addTests(nil)

		// Nothing is allowed
		addCnTests()

		roleVals.AllowLocalhost = true
		commonNames.Localhost = true
		addCnTests()

		roleVals.AllowedDomains = "foobar.com"
		addCnTests()

		roleVals.AllowedDomains = "example.com"
		roleVals.AllowSubdomains = true
		commonNames.SubDomain = true
		commonNames.Wildcard = true
		commonNames.SubSubdomain = true
		commonNames.SubSubdomainWildcard = true
		addCnTests()

		roleVals.AllowedDomains = "foobar.com,example.com"
		commonNames.SecondDomain = true
		roleVals.AllowBareDomains = true
		commonNames.BareDomain = true
		addCnTests()

		roleVals.AllowedDomains = "foobar.com,*example.com"
		roleVals.AllowGlobDomains = true
		commonNames.GlobDomain = true
		addCnTests()

		roleVals.AllowAnyName = true
		roleVals.EnforceHostnames = true
		commonNames.AnyHost = true
		addCnTests()

		roleVals.EnforceHostnames = false
		commonNames.NonHostname = true
		addCnTests()

		// Ensure that we end up with acceptable key sizes since they won't be
		// toggled any longer
		keybitSizeRandOff = true
		addCnTests()
	}
	// OU tests
	{
		roleVals.OU = "foo"
		addTests(getOuCheck(roleVals))

		roleVals.OU = "foo,bar"
		addTests(getOuCheck(roleVals))
	}
	// Organization tests
	{
		roleVals.Organization = "system:masters"
		addTests(getOrganizationCheck(roleVals))

		roleVals.Organization = "foo,bar"
		addTests(getOrganizationCheck(roleVals))
	}
	// IP SAN tests
	{
		roleVals.UseCSRSANs = true
		roleVals.AllowIPSANs = false
		issueTestStep.ErrorOk = false
		addTests(nil)

		roleVals.UseCSRSANs = false
		issueVals.IPSANs = "127.0.0.1,::1"
		issueTestStep.ErrorOk = true
		addTests(nil)

		roleVals.AllowIPSANs = true
		issueTestStep.ErrorOk = false
		addTests(nil)

		issueVals.IPSANs = "foobar"
		issueTestStep.ErrorOk = true
		addTests(nil)

		issueTestStep.ErrorOk = false
		issueVals.IPSANs = ""
	}

	// Lease tests
	{
		roleTestStep.ErrorOk = true
		roleVals.Lease = ""
		roleVals.MaxTTL = ""
		addTests(nil)

		roleVals.Lease = "12h"
		roleVals.MaxTTL = "6h"
		addTests(nil)

		roleTestStep.ErrorOk = false
		roleVals.TTL = ""
		roleVals.MaxTTL = "12h"
	}

	// Listing test
	ret = append(ret, logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "roles/",
		Check: func(resp *logical.Response) error {
			if resp.Data == nil {
				return fmt.Errorf("nil data")
			}

			keysRaw, ok := resp.Data["keys"]
			if !ok {
				return fmt.Errorf("no keys found")
			}

			keys, ok := keysRaw.([]string)
			if !ok {
				return fmt.Errorf("could not convert keys to a string list")
			}

			if len(keys) != 1 {
				return fmt.Errorf("unexpected keys length of %d", len(keys))
			}

			if keys[0] != "test" {
				return fmt.Errorf("unexpected key value of %s", keys[0])
			}

			return nil
		},
	})

	return ret
}

func TestBackend_PathFetchCertList(t *testing.T) {
	// create the backend
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b := Backend()
	err := b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// generate root
	rootData := map[string]interface{}{
		"common_name": "test.com",
		"ttl":         "6h",
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/generate/internal",
		Storage:   storage,
		Data:      rootData,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to generate root, %#v", resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	// config urls
	urlsData := map[string]interface{}{
		"issuing_certificates":    "http://127.0.0.1:8200/v1/pki/ca",
		"crl_distribution_points": "http://127.0.0.1:8200/v1/pki/crl",
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/urls",
		Storage:   storage,
		Data:      urlsData,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to config urls, %#v", resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	// create a role entry
	roleData := map[string]interface{}{
		"allowed_domains":  "test.com",
		"allow_subdomains": "true",
		"max_ttl":          "4h",
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-example",
		Storage:   storage,
		Data:      roleData,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create a role, %#v", resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	// issue some certs
	i := 1
	for i < 10 {
		certData := map[string]interface{}{
			"common_name": "example.test.com",
		}
		resp, err = b.HandleRequest(&logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "issue/test-example",
			Storage:   storage,
			Data:      certData,
		})
		if resp != nil && resp.IsError() {
			t.Fatalf("failed to issue a cert, %#v", resp)
		}
		if err != nil {
			t.Fatal(err)
		}

		i = i + 1
	}

	// list certs
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "certs",
		Storage:   storage,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to list certs, %#v", resp)
	}
	if err != nil {
		t.Fatal(err)
	}
	// check that the root and 9 additional certs are all listed
	if len(resp.Data["keys"].([]string)) != 10 {
		t.Fatalf("failed to list all 10 certs")
	}

	// list certs/
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "certs/",
		Storage:   storage,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to list certs, %#v", resp)
	}
	if err != nil {
		t.Fatal(err)
	}
	// check that the root and 9 additional certs are all listed
	if len(resp.Data["keys"].([]string)) != 10 {
		t.Fatalf("failed to list all 10 certs")
	}
}

func TestBackend_SignVerbatim(t *testing.T) {
	// create the backend
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b := Backend()
	err := b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// generate root
	rootData := map[string]interface{}{
		"common_name": "test.com",
		"ttl":         "172800",
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/generate/internal",
		Storage:   storage,
		Data:      rootData,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to generate root, %#v", *resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	// create a CSR and key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	csrReq := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: "foo.bar.com",
		},
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, csrReq, key)
	if err != nil {
		t.Fatal(err)
	}
	if len(csr) == 0 {
		t.Fatal("generated csr is empty")
	}
	pemCSR := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	})
	if len(pemCSR) == 0 {
		t.Fatal("pem csr is empty")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sign-verbatim",
		Storage:   storage,
		Data: map[string]interface{}{
			"csr": string(pemCSR),
		},
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to sign-verbatim basic CSR: %#v", *resp)
	}
	if err != nil {
		t.Fatal(err)
	}
	if resp.Secret != nil {
		t.Fatal("secret is not nil")
	}

	// create a role entry; we use this to check that sign-verbatim when used with a role is still honoring TTLs
	roleData := map[string]interface{}{
		"ttl":     "4h",
		"max_ttl": "8h",
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test",
		Storage:   storage,
		Data:      roleData,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create a role, %#v", *resp)
	}
	if err != nil {
		t.Fatal(err)
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sign-verbatim/test",
		Storage:   storage,
		Data: map[string]interface{}{
			"csr": string(pemCSR),
			"ttl": "5h",
		},
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to sign-verbatim ttl'd CSR: %#v", *resp)
	}
	if err != nil {
		t.Fatal(err)
	}
	if resp.Secret != nil {
		t.Fatal("got a lease when we should not have")
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sign-verbatim/test",
		Storage:   storage,
		Data: map[string]interface{}{
			"csr": string(pemCSR),
			"ttl": "12h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf(resp.Error().Error())
	}
	if resp.Data == nil || resp.Data["certificate"] == nil {
		t.Fatal("did not get expected data")
	}
	certString := resp.Data["certificate"].(string)
	block, _ := pem.Decode([]byte(certString))
	if block == nil {
		t.Fatal("nil pem block")
	}
	certs, err := x509.ParseCertificates(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(certs) != 1 {
		t.Fatalf("expected a single cert, got %d", len(certs))
	}
	cert := certs[0]
	if math.Abs(float64(time.Now().Add(12*time.Hour).Unix()-cert.NotAfter.Unix())) < 10 {
		t.Fatalf("sign-verbatim did not properly cap validiaty period on signed CSR")
	}

	// now check that if we set generate-lease it takes it from the role and the TTLs match
	roleData = map[string]interface{}{
		"ttl":            "4h",
		"max_ttl":        "8h",
		"generate_lease": true,
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test",
		Storage:   storage,
		Data:      roleData,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create a role, %#v", *resp)
	}
	if err != nil {
		t.Fatal(err)
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sign-verbatim/test",
		Storage:   storage,
		Data: map[string]interface{}{
			"csr": string(pemCSR),
			"ttl": "5h",
		},
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to sign-verbatim role-leased CSR: %#v", *resp)
	}
	if err != nil {
		t.Fatal(err)
	}
	if resp.Secret == nil {
		t.Fatalf("secret is nil, response is %#v", *resp)
	}
	if math.Abs(float64(resp.Secret.TTL-(5*time.Hour))) > float64(5*time.Hour) {
		t.Fatalf("ttl not default; wanted %v, got %v", b.System().DefaultLeaseTTL(), resp.Secret.TTL)
	}
}

func TestBackend_Root_Idempotentcy(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	var err error
	err = client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected ca info")
	}
	resp, err = client.Logical().Read("pki/cert/ca_chain")
	if err != nil {
		t.Fatalf("error reading ca_chain: %v", err)
	}

	r1Data := resp.Data

	// Try again, make sure it's a 204 and same CA
	resp, err = client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatal("expected no ca info")
	}
	resp, err = client.Logical().Read("pki/cert/ca_chain")
	if err != nil {
		t.Fatalf("error reading ca_chain: %v", err)
	}
	r2Data := resp.Data
	if !reflect.DeepEqual(r1Data, r2Data) {
		t.Fatal("got different ca certs")
	}

	resp, err = client.Logical().Delete("pki/root")
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}
	// Make sure it behaves the same
	resp, err = client.Logical().Delete("pki/root")
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatal("expected nil response")
	}

	_, err = client.Logical().Read("pki/cert/ca_chain")
	if err == nil {
		t.Fatal("expected error")
	}

	resp, err = client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected ca info")
	}

	_, err = client.Logical().Read("pki/cert/ca_chain")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_Permitted_DNS_Domains(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	var err error
	err = client.Sys().Mount("root", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.Sys().Mount("int", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "4h",
			MaxLeaseTTL:     "20h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("root/roles/example", map[string]interface{}{
		"allowed_domains":    "foobar.com,zipzap.com,abc.com,xyz.com",
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"max_ttl":            "2h",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("int/roles/example", map[string]interface{}{
		"allowed_domains":    "foobar.com,zipzap.com,abc.com,xyz.com",
		"allow_subdomains":   true,
		"allow_bare_domains": true,
		"max_ttl":            "2h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Direct issuing from root
	_, err = client.Logical().Write("root/root/generate/internal", map[string]interface{}{
		"ttl":                   "40h",
		"common_name":           "myvault.com",
		"permitted_dns_domains": []string{"foobar.com", ".zipzap.com"},
	})
	if err != nil {
		t.Fatal(err)
	}

	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	path := "root/"
	checkIssue := func(valid bool, args ...interface{}) {
		argMap := map[string]interface{}{}
		var currString string
		for i, arg := range args {
			if i%2 == 0 {
				currString = arg.(string)
			} else {
				argMap[currString] = arg
			}
		}
		_, err = client.Logical().Write(path+"issue/example", argMap)
		switch {
		case valid && err != nil:
			t.Fatal(err)
		case !valid && err == nil:
			t.Fatal("expected error")
		}

		csr, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName: argMap["common_name"].(string),
			},
		}, clientKey)
		if err != nil {
			t.Fatal(err)
		}
		delete(argMap, "common_name")
		argMap["csr"] = string(pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: csr,
		}))

		_, err = client.Logical().Write(path+"sign/example", argMap)
		switch {
		case valid && err != nil:
			t.Fatal(err)
		case !valid && err == nil:
			t.Fatal("expected error")
		}
	}

	// Check issuing and signing against root's permitted domains
	checkIssue(false, "common_name", "zipzap.com")
	checkIssue(false, "common_name", "host.foobar.com")
	checkIssue(true, "common_name", "host.zipzap.com")
	checkIssue(true, "common_name", "foobar.com")

	// Verify that root also won't issue an intermediate outside of its permitted domains
	resp, err := client.Logical().Write("int/intermediate/generate/internal", map[string]interface{}{
		"common_name": "issuer.abc.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("root/root/sign-intermediate", map[string]interface{}{
		"common_name": "issuer.abc.com",
		"csr":         resp.Data["csr"],
		"permitted_dns_domains": []string{"abc.com", ".xyz.com"},
		"ttl": "5h",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	_, err = client.Logical().Write("root/root/sign-intermediate", map[string]interface{}{
		"use_csr_values": true,
		"csr":            resp.Data["csr"],
		"permitted_dns_domains": []string{"abc.com", ".xyz.com"},
		"ttl": "5h",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// Sign a valid intermediate
	resp, err = client.Logical().Write("root/root/sign-intermediate", map[string]interface{}{
		"common_name": "issuer.zipzap.com",
		"csr":         resp.Data["csr"],
		"permitted_dns_domains": []string{"abc.com", ".xyz.com"},
		"ttl": "5h",
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err = client.Logical().Write("int/intermediate/set-signed", map[string]interface{}{
		"certificate": resp.Data["certificate"],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check enforcement with the intermediate's set values
	path = "int/"
	checkIssue(false, "common_name", "host.abc.com")
	checkIssue(false, "common_name", "xyz.com")
	checkIssue(true, "common_name", "abc.com")
	checkIssue(true, "common_name", "host.xyz.com")
}

func TestBackend_SignIntermediate_AllowedPastCA(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	var err error
	err = client.Sys().Mount("root", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "60h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = client.Sys().Mount("int", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "4h",
			MaxLeaseTTL:     "20h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Direct issuing from root
	_, err = client.Logical().Write("root/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("root/roles/test", map[string]interface{}{
		"allow_bare_domains": true,
		"allow_subdomains":   true,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().Write("int/intermediate/generate/internal", map[string]interface{}{
		"common_name": "myint.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	csr := resp.Data["csr"]

	_, err = client.Logical().Write("root/sign/test", map[string]interface{}{
		"common_name": "myint.com",
		"csr":         csr,
		"ttl":         "60h",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	_, err = client.Logical().Write("root/sign-verbatim/test", map[string]interface{}{
		"common_name": "myint.com",
		"csr":         csr,
		"ttl":         "60h",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	resp, err = client.Logical().Write("root/root/sign-intermediate", map[string]interface{}{
		"common_name": "myint.com",
		"csr":         csr,
		"ttl":         "60h",
	})
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	if resp == nil {
		t.Fatal("got nil response")
	}
	if len(resp.Warnings) == 0 {
		t.Fatalf("expected warnings, got %#v", *resp)
	}
}

func TestBackend_SignSelfIssued(t *testing.T) {
	// create the backend
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b := Backend()
	err := b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// generate root
	rootData := map[string]interface{}{
		"common_name": "test.com",
		"ttl":         "172800",
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/generate/internal",
		Storage:   storage,
		Data:      rootData,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to generate root, %#v", *resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	getSelfSigned := func(subject, issuer *x509.Certificate) (string, *x509.Certificate) {
		selfSigned, err := x509.CreateCertificate(rand.Reader, subject, issuer, key.Public(), key)
		if err != nil {
			t.Fatal(err)
		}
		cert, err := x509.ParseCertificate(selfSigned)
		if err != nil {
			t.Fatal(err)
		}
		pemSS := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: selfSigned,
		})
		return string(pemSS), cert
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "foo.bar.com",
		},
		SerialNumber: big.NewInt(1234),
		IsCA:         false,
		BasicConstraintsValid: true,
	}

	ss, _ := getSelfSigned(template, template)
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-self-issued",
		Storage:   storage,
		Data: map[string]interface{}{
			"certificate": ss,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response")
	}
	if !resp.IsError() {
		t.Fatalf("expected error due to non-CA; got: %#v", *resp)
	}

	// Set CA to true, but leave issuer alone
	template.IsCA = true

	issuer := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "bar.foo.com",
		},
		SerialNumber: big.NewInt(2345),
		IsCA:         true,
		BasicConstraintsValid: true,
	}
	ss, ssCert := getSelfSigned(template, issuer)
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-self-issued",
		Storage:   storage,
		Data: map[string]interface{}{
			"certificate": ss,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response")
	}
	if !resp.IsError() {
		t.Fatalf("expected error due to different issuer; cert info is\nIssuer\n%#v\nSubject\n%#v\n", ssCert.Issuer, ssCert.Subject)
	}

	ss, ssCert = getSelfSigned(template, template)
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-self-issued",
		Storage:   storage,
		Data: map[string]interface{}{
			"certificate": ss,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response")
	}
	if resp.IsError() {
		t.Fatalf("error in response: %s", resp.Error().Error())
	}

	newCertString := resp.Data["certificate"].(string)
	block, _ := pem.Decode([]byte(newCertString))
	newCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	signingBundle, err := fetchCAInfo(&logical.Request{Storage: storage})
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(newCert.Subject, newCert.Issuer) {
		t.Fatal("expected different subject/issuer")
	}
	if !reflect.DeepEqual(newCert.Issuer, signingBundle.Certificate.Subject) {
		t.Fatalf("expected matching issuer/CA subject\n\nIssuer:\n%#v\nSubject:\n%#v\n", newCert.Issuer, signingBundle.Certificate.Subject)
	}
	if bytes.Equal(newCert.AuthorityKeyId, newCert.SubjectKeyId) {
		t.Fatal("expected different authority/subject")
	}
	if !bytes.Equal(newCert.AuthorityKeyId, signingBundle.Certificate.SubjectKeyId) {
		t.Fatal("expected authority on new cert to be same as signing subject")
	}
	if newCert.Subject.CommonName != "foo.bar.com" {
		t.Fatalf("unexpected common name on new cert: %s", newCert.Subject.CommonName)
	}
}

const (
	rsaCAKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAmPQlK7xD5p+E8iLQ8XlVmll5uU2NKMxKY3UF5tbh+0vkc+Fy
XmutLxxXAyYRPoztZ1g7ocr8XBFYsQPK26TFc3TzrLL7bBEYHQArd8M+VUHjziB7
zwwpbV7tG8WPqIScDKMNncavDcT8sDg3DUqb8/zWkBD8WEYmsVr1VfKY5pFdxIZU
kHP3/MkDpGmfrED9K5qPu17dIHTL2VYi4KxKhtIryapZTk6vDwRNfIYJD23QbQnt
Si1j0X9MTRUf3BIcd0Ch60aGvv0VSL+1NTafsZQD+z1RY/zNp9IUHz5bNIiePZ6l
JrlddodAAXZ4sN1CMetf4bA2RXssxBEIb5FyiQIDAQABAoIBAGMScSk9DvZJCUIV
zyU6JHqPzkp6sx5kBSMa37HAKiwt4lI1C3GhaVIEl0/Qzoannfa8rhOEeaXhDoPK
IxHWTpcUf+mzHSvIfsf6Hi2655stzLLtU4SvKf5P6GF+vCi5jKKa0u0JjsXqfIpg
Pzh6xT1q3kf+2JUNC28Brbv4IZXmPmqWwu21VN+t3GsMGYgOnEOzBjXMhvNnm9kN
kznV9Y2y0UIcT4dhbe2VRs4Dp8dGEyrFM7/Ovb3hIJrTkPcxjBbL5eMqpXnIkiW2
7NyPMWFvX2lGnGdZ1Erh65SVtMjnHFwnSJ8jD+x9RAH9c1LQrYASws3MvMV8Bdzg
2iljNqECgYEAw3Ow0clLx2alj9qFXcS2ap1lUCJxXZ9UiIU5lOcPxpCpHPloua14
46rj2EJ9SD1L2kyB5gCq4nGK5uUIx37AJryy1SGzUmtmIVxQLnm6XK6zKnTBk0gx
gevS6D7fHLDiVGGl3oGw4evibUFCk7dFOb/I/uBRb1zyaJrqOIlDS7UCgYEAyFYi
RYQbYJJ0k18fUWDKy/P/Rl7uy9D67Qa9+wxoYN2Kh/aQwnNxYHAbwG7Pupd0oGcW
Yl4bgUliAX3IFGs/cCkPJAIHzwWBPjUDhsJ020TGxKfL4SWP9OaxOpN5TOAixvBY
ar9aSaKEl7QShmzc/Dknxu58LcoZUwI82pKIGAUCgYAxaHJ/ZcpxOsKJjez+2jZe
1zEAQ+SyjQ96f2sh+BMl1/XYLDhMD80qiE2WoqA2/b/KDGMd+Hc6TQeW/LjubV03
raXreNxy7lFgB40BYqY4vbTu+5rfl3VkaW/kY9hU0WY1fIXIrLJBOjb/9WpWGxM1
2QR/YcdURoPE67xf1FsdrQKBgE8KdNEakzah8e6nLBMOblTTutcH4410sVvNOi2P
sqrtHZgRNwIRTB0xfjGJRtomoXQb2CANYyq6SjmuZ79upQPan0ekqXILiPeDMRX9
KN/OHeI/FdiJ2mdUkX476zLih7YX47qSLsw4m7nC6UAyOWomHsSFGWdzglRW4K2X
/KwFAoGAYQUEWhXp5vpKzAly1ivSH9+sGC59Cujdy50oJSjaw9J+W1fM5WO9z+MH
CoEpRt8epIgvCBBP2IM7uJUu8i2jQgJ/rrn3NTJgZn2UEPzyxUxbuWnSyueyUsD6
uhTwBDf8LWOpvdZHMI4CPZ5WJwxAGkvde9xtlzuZUSAlyI2X8m0=
-----END RSA PRIVATE KEY-----
`
	rsaCACert string = `-----BEGIN CERTIFICATE-----
MIIDljCCAn6gAwIBAgIUQVapfgyAeDH9rAmpw3PQrhMcjRMwDQYJKoZIhvcNAQEL
BQAwMzExMC8GA1UEAxMoVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3ViIEF1
dGhvcml0eTAeFw0xNjA4MDcyMjUzNTRaFw0yNjA3MjQxMDU0MjRaMDcxNTAzBgNV
BAMTLFZhdWx0IFRlc3RpbmcgSW50ZXJtZWRpYXRlIFN1YiBTdWIgQXV0aG9yaXR5
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAmPQlK7xD5p+E8iLQ8XlV
mll5uU2NKMxKY3UF5tbh+0vkc+FyXmutLxxXAyYRPoztZ1g7ocr8XBFYsQPK26TF
c3TzrLL7bBEYHQArd8M+VUHjziB7zwwpbV7tG8WPqIScDKMNncavDcT8sDg3DUqb
8/zWkBD8WEYmsVr1VfKY5pFdxIZUkHP3/MkDpGmfrED9K5qPu17dIHTL2VYi4KxK
htIryapZTk6vDwRNfIYJD23QbQntSi1j0X9MTRUf3BIcd0Ch60aGvv0VSL+1NTaf
sZQD+z1RY/zNp9IUHz5bNIiePZ6lJrlddodAAXZ4sN1CMetf4bA2RXssxBEIb5Fy
iQIDAQABo4GdMIGaMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0G
A1UdDgQWBBRMeQTX9VkLqb1wzrvN/vFG09yhUTAfBgNVHSMEGDAWgBR0Oq2VTUBE
dOm6a1sKJTvdZMV5LjA3BgNVHREEMDAugixWYXVsdCBUZXN0aW5nIEludGVybWVk
aWF0ZSBTdWIgU3ViIEF1dGhvcml0eTANBgkqhkiG9w0BAQsFAAOCAQEAagYM7uFa
tUziraBkuU7cIyX83y7lYFsDhUse2hkpqmgO14oEOwFsDox1Jg2QGt4FEfJoCOXf
oCZZN8XmaWdSrfgs1nDmtE0xwXiX1z7JuJZ+Ygt3dcRHO1zs5tmuHLxrvMnKfIfG
bsGmES4mknt0qQ7tGhpyC+KgEmcVL1QQJXNjzCrw5iQ9sgvQt+oCqV28pxOUSYkq
FdrozmNdJwMgVADywiY/FqYJWgkixlFHQkPR7eiXwpahON+zRMk1JSgr/8N8fRDj
aqVBRppPzVU9joUME0vOc8cK3VozNe4iRkKNZFelHU2NPPJSDjRLVH9tJ7jPVOEA
/k6w2PwdoRom7Q==
-----END CERTIFICATE-----
`

	rsaCAChain string = `-----BEGIN CERTIFICATE-----
MIIDijCCAnKgAwIBAgIUOiGo/1EOhRhuupTRGDYnqdALk/swDQYJKoZIhvcNAQEL
BQAwLzEtMCsGA1UEAxMkVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgQXV0aG9y
aXR5MB4XDTE2MDgwNzIyNTA1MloXDTI2MDcyODE0NTEyMlowMzExMC8GA1UEAxMo
VmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3ViIEF1dGhvcml0eTCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBAMTPRQREwW3BEifNcm0XElMRB0GNTXHr
XCuNoFVsVBlIEsNVQkka+SHZcmNBdEcZLBXP/W3tBT82B48GVN8jyxAGfYZ5hoOQ
ed3GVft1A7lAnxcGvf5e9kfecKDcBB4G4rBhqdDNcAtklS2hV4uZUcVcEJKggpsQ
a1wZkCn8eg6sqEYG/SxPouwL52PblxIN+Dd57sBeqx4qdL297XR8LuLkxqftwUCZ
l2iFBnSDID/06ZmHDXA38I0n3jT2ZGjgPGFnIFKxRGq1vpVc3F5ga8qk+u66ybBu
xWHzINQrrryjELbl2YBTr6i0R9HnZle6OPcXMWp0JuGjtDC1xb5NmnkCAwEAAaOB
mTCBljAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU
dDqtlU1ARHTpumtbCiU73WTFeS4wHwYDVR0jBBgwFoAU+UO/nrlKr4COZCxLZSY/
ul+YMvMwMwYDVR0RBCwwKoIoVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3Vi
IEF1dGhvcml0eTANBgkqhkiG9w0BAQsFAAOCAQEAjgCuTXsLFkf0DVkfSsKReNwI
U/yBcP8Ttbx/ltanJGIVfD5TZoCnNTWm6RkML29ohfxI27sHTUhj+/6Ba0MRiLeI
FXdclXmHOU2dTHlrmUa0m/4cb5uYoiiEnpmyWL5k94fqPOZAvJcFHnP3db4vsaUW
47YcOvJbPSJqFXZHadqnsf3Fur5NCeTkIk6yZSvwTaZJT0JIWcqfE5LK3mYAMMC3
iPaIa1cYqOZhWx9ilQfW6u6WxWeOphGuDIusP7Q4qc2Dr9sekyD59dfIYsroK5TP
QVJb69nIYINpYdg3l3VNmmkY4G30N9QNs6acaH49rYzLcRX6tLBgPklO6d+TPA==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIDejCCAmKgAwIBAgIUULdIdrdK4Y8d+XM9fuOpDlNcJIYwDQYJKoZIhvcNAQEL
BQAwJzElMCMGA1UEAxMcVmF1bHQgVGVzdGluZyBSb290IEF1dGhvcml0eTAeFw0x
NjA4MDcyMjUwNTFaFw0yNjA4MDExODUxMjFaMC8xLTArBgNVBAMTJFZhdWx0IFRl
c3RpbmcgSW50ZXJtZWRpYXRlIEF1dGhvcml0eTCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBANXa6U+MDiUrryeZeGxgkmAZdrm9wCKz/6SmxYSebKr8aZwD
nfbsPLRFxU6BXp9Nc6pP7e8HLBv6PtFTQG389zxOBwAHxZQvUsFESumUd64oTLRG
J+AErTh7rtSWbLZsgDtQVvpx+6mKkvm53f/aKcq+DbqAFOg6slYOaQix0ZvP/qL0
iWGIPr1JZk9uBJOUuIUBJdbsgTk+KQqJL9M6up8bCnM0noCafwrNKwZWtsbkfOZE
OLSycdzCEBeHejpHTIU0vgAkdj63oEy2AbK3hMPxKzNthL3DX6W0tssoVgL//92i
oSfpDTxiXqqdr+J3accpsAvA+F+D2TqaxdAfjLcCAwEAAaOBlTCBkjAOBgNVHQ8B
Af8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU+UO/nrlKr4COZCxL
ZSY/ul+YMvMwHwYDVR0jBBgwFoAUA3jY4OUWi1Y7zQgM7S9QeXjNgIQwLwYDVR0R
BCgwJoIkVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgQXV0aG9yaXR5MA0GCSqG
SIb3DQEBCwUAA4IBAQA9VJt92LsOOegtAx35rr41LSSfPWB2SCKg0fphL2gMPO5y
fE2u8O5TF5dJWJ1fF3cg9/XK30Ohdl1/ujfHbVcX+O6Sb3cgwKTQQOhTdsAZKFRT
RPHaf/Ja8uqAXMITApxOp7YiYQwukwZr+OsKi66s+zhlfI790PoQbC3UJvgmNDmv
V5oP63mw4yhqRNPn4NOjzoC/hJSIM0AIdRB1nx2rsSUw0P354R1j9gO43L/Lj33S
NEaPmw+SC3Tbcx4yxeKnTvGdu3sw/ndmZkCjaq5jxgTy9FONqT45TPJOyk29o5gl
+AVQz5fD2M3C1L/sZIPH2OQbXxePHcsvUZVgaKyk
-----END CERTIFICATE-----
`

	ecCAKey string = `-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDBP/t89wrC0RFVs0N+jiRuGPptoxI1Iyu42/PzzZWMKYnO7yCWFG/Qv
zC8cRa8PDqegBwYFK4EEACKhZANiAAQI9e8n9RD6gOd5YpWpDi5AoPbskxQSogxx
dYFzzHwS0RYIucmlcJ2CuJQNc+9E4dUCMsYr2cAnCgA4iUHzGaje3Fa4O667LVH1
imAyAj5nbfSd89iNzg4XNPkFjuVNBlE=
-----END EC PRIVATE KEY-----
`

	ecCACert string = `-----BEGIN CERTIFICATE-----
MIIDHzCCAqSgAwIBAgIUEQ4L+8Xl9+/uxU3MMCrd3Bw0HMcwCgYIKoZIzj0EAwIw
XzEjMCEGA1UEAxMaVmF1bHQgRUMgdGVzdGluZyByb290IGNlcnQxODA2BgNVBAUT
Lzk3MzY2MDk3NDQ1ODU2MDI3MDY5MDQ0MTkxNjIxODI4NjI0NjM0NTI5MTkzMTU5
MB4XDTE1MTAwNTE2MzAwMFoXDTM1MDkzMDE2MzAwMFowXzEjMCEGA1UEAxMaVmF1
bHQgRUMgdGVzdGluZyByb290IGNlcnQxODA2BgNVBAUTLzk3MzY2MDk3NDQ1ODU2
MDI3MDY5MDQ0MTkxNjIxODI4NjI0NjM0NTI5MTkzMTU5MHYwEAYHKoZIzj0CAQYF
K4EEACIDYgAECPXvJ/UQ+oDneWKVqQ4uQKD27JMUEqIMcXWBc8x8EtEWCLnJpXCd
griUDXPvROHVAjLGK9nAJwoAOIlB8xmo3txWuDuuuy1R9YpgMgI+Z230nfPYjc4O
FzT5BY7lTQZRo4IBHzCCARswDgYDVR0PAQH/BAQDAgGuMBMGA1UdJQQMMAoGCCsG
AQUFBwMJMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFCIBqs15CiKuj7vqmIW5
L07WSeLhMB8GA1UdIwQYMBaAFCIBqs15CiKuj7vqmIW5L07WSeLhMEIGCCsGAQUF
BwEBBDYwNDAyBggrBgEFBQcwAoYmaHR0cDovL3ZhdWx0LmV4YW1wbGUuY29tL3Yx
L3Jvb3Rwa2kvY2EwJQYDVR0RBB4wHIIaVmF1bHQgRUMgdGVzdGluZyByb290IGNl
cnQwOAYDVR0fBDEwLzAtoCugKYYnaHR0cDovL3ZhdWx0LmV4YW1wbGUuY29tL3Yx
L3Jvb3Rwa2kvY3JsMAoGCCqGSM49BAMCA2kAMGYCMQDRrxXskBtXjuZ1tUTk+qae
3bNVE1oeTDJhe0m3KN7qTykSGslxfEjlv83GYXziiv0CMQDsqu1U9uXPn3ezSbgG
O30prQ/sanDzNAeJhftoGtNPJDspwx0fzclHvKIhgl3JVUc=
-----END CERTIFICATE-----
`
)
