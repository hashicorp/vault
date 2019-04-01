package pki

import (
	"bytes"
	"context"
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
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
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
	"golang.org/x/net/idna"
)

var (
	stepCount               = 0
	serialUnderTest         string
	parsedKeyUsageUnderTest int
)

func TestPKI_RequireCN(t *testing.T) {
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

	// Create a role which does require CN (default)
	_, err = client.Logical().Write("pki/roles/example", map[string]interface{}{
		"allowed_domains":    "foobar.com,zipzap.com,abc.com,xyz.com",
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"max_ttl":            "2h",
	})

	// Issue a cert with require_cn set to true and with common name supplied.
	// It should succeed.
	resp, err = client.Logical().Write("pki/issue/example", map[string]interface{}{
		"common_name": "foobar.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue a cert with require_cn set to true and with out supplying the
	// common name. It should error out.
	resp, err = client.Logical().Write("pki/issue/example", map[string]interface{}{})
	if err == nil {
		t.Fatalf("expected an error due to missing common_name")
	}

	// Modify the role to make the common name optional
	_, err = client.Logical().Write("pki/roles/example", map[string]interface{}{
		"allowed_domains":    "foobar.com,zipzap.com,abc.com,xyz.com",
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"max_ttl":            "2h",
		"require_cn":         false,
	})

	// Issue a cert with require_cn set to false and without supplying the
	// common name. It should succeed.
	resp, err = client.Logical().Write("pki/issue/example", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	if resp.Data["certificate"] == "" {
		t.Fatalf("expected a cert to be generated")
	}

	// Issue a cert with require_cn set to false and with a common name. It
	// should succeed.
	resp, err = client.Logical().Write("pki/issue/example", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	if resp.Data["certificate"] == "" {
		t.Fatalf("expected a cert to be generated")
	}
}

func TestBackend_CSRValues(t *testing.T) {
	initTest.Do(setCerts)
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		LogicalBackend: b,
		Steps:          []logicaltest.TestStep{},
	}

	intdata := map[string]interface{}{}
	reqdata := map[string]interface{}{}
	testCase.Steps = append(testCase.Steps, generateCSRSteps(t, ecCACert, ecCAKey, intdata, reqdata)...)

	logicaltest.Test(t, testCase)
}

func TestBackend_URLsCRUD(t *testing.T) {
	initTest.Do(setCerts)
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		LogicalBackend: b,
		Steps:          []logicaltest.TestStep{},
	}

	intdata := map[string]interface{}{}
	reqdata := map[string]interface{}{}
	testCase.Steps = append(testCase.Steps, generateURLSteps(t, ecCACert, ecCAKey, intdata, reqdata)...)

	logicaltest.Test(t, testCase)
}

// Generates and tests steps that walk through the various possibilities
// of role flags to ensure that they are properly restricted
// Uses the RSA CA key
func TestBackend_RSARoles(t *testing.T) {
	initTest.Do(setCerts)
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "config/ca",
				Data: map[string]interface{}{
					"pem_bundle": strings.Join([]string{rsaCAKey, rsaCACert}, "\n"),
				},
			},
		},
	}

	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, false)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+1, v)
		}
	}

	logicaltest.Test(t, testCase)
}

// Generates and tests steps that walk through the various possibilities
// of role flags to ensure that they are properly restricted
// Uses the RSA CA key
func TestBackend_RSARoles_CSR(t *testing.T) {
	initTest.Do(setCerts)
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "config/ca",
				Data: map[string]interface{}{
					"pem_bundle": strings.Join([]string{rsaCAKey, rsaCACert}, "\n"),
				},
			},
		},
	}

	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, true)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+1, v)
		}
	}

	logicaltest.Test(t, testCase)
}

// Generates and tests steps that walk through the various possibilities
// of role flags to ensure that they are properly restricted
// Uses the EC CA key
func TestBackend_ECRoles(t *testing.T) {
	initTest.Do(setCerts)
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "config/ca",
				Data: map[string]interface{}{
					"pem_bundle": strings.Join([]string{ecCAKey, ecCACert}, "\n"),
				},
			},
		},
	}

	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, false)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+1, v)
		}
	}

	logicaltest.Test(t, testCase)
}

// Generates and tests steps that walk through the various possibilities
// of role flags to ensure that they are properly restricted
// Uses the EC CA key
func TestBackend_ECRoles_CSR(t *testing.T) {
	initTest.Do(setCerts)
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 32
	b, err := Factory(context.Background(), &logical.BackendConfig{
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
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "config/ca",
				Data: map[string]interface{}{
					"pem_bundle": strings.Join([]string{ecCAKey, ecCACert}, "\n"),
				},
			},
		},
	}

	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, true)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+1, v)
		}
	}

	logicaltest.Test(t, testCase)
}

// Performs some validity checking on the returned bundles
func checkCertsAndPrivateKey(keyType string, key crypto.Signer, usage x509.KeyUsage, extUsage x509.ExtKeyUsage, validity time.Duration, certBundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	parsedCertBundle, err := certBundle.ToParsedCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error parsing cert bundle: %s", err)
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
				return nil, fmt.Errorf("error parsing EC key: %s", err)
			}
		}
	}

	switch {
	case parsedCertBundle.Certificate == nil:
		return nil, fmt.Errorf("did not find a certificate in the cert bundle")
	case len(parsedCertBundle.CAChain) == 0 || parsedCertBundle.CAChain[0].Certificate == nil:
		return nil, fmt.Errorf("did not find a CA in the cert bundle")
	case parsedCertBundle.PrivateKey == nil:
		return nil, fmt.Errorf("did not find a private key in the cert bundle")
	case parsedCertBundle.PrivateKeyType == certutil.UnknownPrivateKey:
		return nil, fmt.Errorf("could not figure out type of private key")
	}

	switch {
	case parsedCertBundle.PrivateKeyType == certutil.RSAPrivateKey && keyType != "rsa":
		fallthrough
	case parsedCertBundle.PrivateKeyType == certutil.ECPrivateKey && keyType != "ec":
		return nil, fmt.Errorf("given key type does not match type found in bundle")
	}

	cert := parsedCertBundle.Certificate

	if usage != cert.KeyUsage {
		return nil, fmt.Errorf("expected usage of %#v, got %#v; ext usage is %#v", usage, cert.KeyUsage, cert.ExtKeyUsage)
	}

	// There should only be one ext usage type, because only one is requested
	// in the tests
	if len(cert.ExtKeyUsage) != 1 {
		return nil, fmt.Errorf("got wrong size key usage in generated cert; expected 1, values are %#v", cert.ExtKeyUsage)
	}
	switch extUsage {
	case x509.ExtKeyUsageEmailProtection:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageEmailProtection {
			return nil, fmt.Errorf("bad extended key usage")
		}
	case x509.ExtKeyUsageServerAuth:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageServerAuth {
			return nil, fmt.Errorf("bad extended key usage")
		}
	case x509.ExtKeyUsageClientAuth:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth {
			return nil, fmt.Errorf("bad extended key usage")
		}
	case x509.ExtKeyUsageCodeSigning:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageCodeSigning {
			return nil, fmt.Errorf("bad extended key usage")
		}
	}

	if math.Abs(float64(time.Now().Add(validity).Unix()-cert.NotAfter.Unix())) > 20 {
		return nil, fmt.Errorf("certificate validity end: %s; expected within 20 seconds of %s", cert.NotAfter.Format(time.RFC3339), time.Now().Add(validity).Format(time.RFC3339))
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
	csrPem1024 := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr1024,
	})))

	priv2048, _ := rsa.GenerateKey(rand.Reader, 2048)
	csr2048, _ := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, priv2048)
	csrPem2048 := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr2048,
	})))

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
				"common_name": "intermediate.cert.com",
				"csr":         csrPem1024,
				"format":      "der",
			},
			ErrorOk: true,
			Check: func(resp *logical.Response) error {
				if !resp.IsError() {
					return fmt.Errorf("expected an error response but did not get one")
				}
				if !strings.Contains(resp.Data["error"].(string), "2048") {
					return fmt.Errorf("received an error but not about a 1024-bit key, error was: %s", resp.Data["error"].(string))
				}

				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data: map[string]interface{}{
				"common_name": "intermediate.cert.com",
				"csr":         csrPem2048,
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
				case !reflect.DeepEqual([]string{"intermediate.cert.com"}, cert.DNSNames):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", []string{"intermediate.cert.com"}, cert.DNSNames)
				}

				return nil
			},
		},

		// Same as above but exclude adding to sans
		logicaltest.TestStep{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data: map[string]interface{}{
				"common_name":          "intermediate.cert.com",
				"csr":                  csrPem2048,
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
	csrPem := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	})))

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
				"csr":            csrPem,
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
				"csr":            csrPem,
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

// Generates steps to test out various role permutations
func generateRoleSteps(t *testing.T, useCSRs bool) []logicaltest.TestStep {
	roleVals := roleEntry{
		MaxTTL:    12 * time.Hour,
		KeyType:   "rsa",
		KeyBits:   2048,
		RequireCN: true,
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
		return fmt.Errorf("expected an error, but did not seem to get one")
	}

	// Adds tests with the currently configured issue/role information
	addTests := func(testCheck logicaltest.TestCheckFunc) {
		stepCount++
		//t.Logf("test step %d\nrole vals: %#v\n", stepCount, roleVals)
		stepCount++
		//t.Logf("test step %d\nissue vals: %#v\n", stepCount, issueTestStep)
		roleTestStep.Data = roleVals.ToResponseData()
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

	getCountryCheck := func(role roleEntry) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return fmt.Errorf("error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate

			expected := strutil.RemoveDuplicates(role.Country, true)
			if !reflect.DeepEqual(cert.Subject.Country, expected) {
				return fmt.Errorf("error: returned certificate has Country of %s but %s was specified in the role", cert.Subject.Country, expected)
			}
			return nil
		}
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
				return fmt.Errorf("error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate

			expected := strutil.RemoveDuplicates(role.OU, true)
			if !reflect.DeepEqual(cert.Subject.OrganizationalUnit, expected) {
				return fmt.Errorf("error: returned certificate has OU of %s but %s was specified in the role", cert.Subject.OrganizationalUnit, expected)
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
				return fmt.Errorf("error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate

			expected := strutil.RemoveDuplicates(role.Organization, true)
			if !reflect.DeepEqual(cert.Subject.Organization, expected) {
				return fmt.Errorf("error: returned certificate has Organization of %s but %s was specified in the role", cert.Subject.Organization, expected)
			}
			return nil
		}
	}

	getLocalityCheck := func(role roleEntry) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return fmt.Errorf("error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate

			expected := strutil.RemoveDuplicates(role.Locality, true)
			if !reflect.DeepEqual(cert.Subject.Locality, expected) {
				return fmt.Errorf("error: returned certificate has Locality of %s but %s was specified in the role", cert.Subject.Locality, expected)
			}
			return nil
		}
	}

	getProvinceCheck := func(role roleEntry) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return fmt.Errorf("error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate

			expected := strutil.RemoveDuplicates(role.Province, true)
			if !reflect.DeepEqual(cert.Subject.Province, expected) {
				return fmt.Errorf("error: returned certificate has Province of %s but %s was specified in the role", cert.Subject.Province, expected)
			}
			return nil
		}
	}

	getStreetAddressCheck := func(role roleEntry) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return fmt.Errorf("error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate

			expected := strutil.RemoveDuplicates(role.StreetAddress, true)
			if !reflect.DeepEqual(cert.Subject.StreetAddress, expected) {
				return fmt.Errorf("error: returned certificate has StreetAddress of %s but %s was specified in the role", cert.Subject.StreetAddress, expected)
			}
			return nil
		}
	}

	getPostalCodeCheck := func(role roleEntry) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return fmt.Errorf("error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate

			expected := strutil.RemoveDuplicates(role.PostalCode, true)
			if !reflect.DeepEqual(cert.Subject.PostalCode, expected) {
				return fmt.Errorf("error: returned certificate has PostalCode of %s but %s was specified in the role", cert.Subject.PostalCode, expected)
			}
			return nil
		}
	}

	getNotBeforeCheck := func(role roleEntry) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return fmt.Errorf("error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate

			actualDiff := time.Now().Sub(cert.NotBefore)
			certRoleDiff := (role.NotBeforeDuration - actualDiff).Truncate(time.Second)
			// These times get truncated, so give a 1 second buffer on each side
			if certRoleDiff >= -1*time.Second && certRoleDiff <= 1*time.Second {
				return nil
			}
			return fmt.Errorf("validity period out of range diff: %v", certRoleDiff)
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
				return fmt.Errorf("error checking generated certificate: %s", err)
			}
			cert := parsedCertBundle.Certificate
			if cert.Subject.CommonName != name {
				return fmt.Errorf("error: returned certificate has CN of %s but %s was requested", cert.Subject.CommonName, name)
			}
			if strings.Contains(cert.Subject.CommonName, "@") {
				if len(cert.DNSNames) != 0 || len(cert.EmailAddresses) != 1 {
					return fmt.Errorf("error: found more than one DNS SAN or not one Email SAN but only one was requested, cert.DNSNames = %#v, cert.EmailAddresses = %#v", cert.DNSNames, cert.EmailAddresses)
				}
			} else {
				if len(cert.DNSNames) != 1 || len(cert.EmailAddresses) != 0 {
					return fmt.Errorf("error: found more than one Email SAN or not one DNS SAN but only one was requested, cert.DNSNames = %#v, cert.EmailAddresses = %#v", cert.DNSNames, cert.EmailAddresses)
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
				// Check IDNA
				p := idna.New(
					idna.StrictDomainName(true),
					idna.VerifyDNSLength(true),
				)
				converted, err := p.ToUnicode(retName)
				if err != nil {
					t.Fatal(err)
				}
				if converted != name {
					return fmt.Errorf("error: returned certificate has a DNS SAN of %s (from idna: %s) but %s was requested", retName, converted, name)
				}
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
		IDN                  bool `structs:"daɪˈɛrɨsɨs"`
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

			var usage []string
			if mathRand.Int()%2 == 1 {
				usage = append(usage, "DigitalSignature")
			}
			if mathRand.Int()%2 == 1 {
				usage = append(usage, "ContentCoMmitment")
			}
			if mathRand.Int()%2 == 1 {
				usage = append(usage, "KeyEncipherment")
			}
			if mathRand.Int()%2 == 1 {
				usage = append(usage, "DataEncipherment")
			}
			if mathRand.Int()%2 == 1 {
				usage = append(usage, "KeyAgreemEnt")
			}
			if mathRand.Int()%2 == 1 {
				usage = append(usage, "CertSign")
			}
			if mathRand.Int()%2 == 1 {
				usage = append(usage, "CRLSign")
			}
			if mathRand.Int()%2 == 1 {
				usage = append(usage, "EncipherOnly")
			}
			if mathRand.Int()%2 == 1 {
				usage = append(usage, "DecipherOnly")
			}

			roleVals.KeyUsage = usage
			parsedKeyUsage := parseKeyUsages(roleVals.KeyUsage)
			if parsedKeyUsage == 0 && len(usage) != 0 {
				panic("parsed key usages was zero")
			}
			parsedKeyUsageUnderTest = parsedKeyUsage

			var extUsage x509.ExtKeyUsage
			i := mathRand.Int() % 4
			switch {
			case i == 0:
				// Punt on this for now since I'm not clear the actual proper
				// way to format these
				if name != "daɪˈɛrɨsɨs" {
					extUsage = x509.ExtKeyUsageEmailProtection
					roleVals.EmailProtectionFlag = true
					break
				}
				fallthrough
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

			validity := roleVals.MaxTTL

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

		roleVals.AllowedDomains = []string{"foobar.com"}
		addCnTests()

		roleVals.AllowedDomains = []string{"example.com"}
		roleVals.AllowSubdomains = true
		commonNames.SubDomain = true
		commonNames.Wildcard = true
		commonNames.SubSubdomain = true
		commonNames.SubSubdomainWildcard = true
		addCnTests()

		roleVals.AllowedDomains = []string{"foobar.com", "example.com"}
		commonNames.SecondDomain = true
		roleVals.AllowBareDomains = true
		commonNames.BareDomain = true
		addCnTests()

		roleVals.AllowedDomains = []string{"foobar.com", "*example.com"}
		roleVals.AllowGlobDomains = true
		commonNames.GlobDomain = true
		addCnTests()

		roleVals.AllowAnyName = true
		roleVals.EnforceHostnames = true
		commonNames.AnyHost = true
		commonNames.IDN = true
		addCnTests()

		roleVals.EnforceHostnames = false
		addCnTests()

		// Ensure that we end up with acceptable key sizes since they won't be
		// toggled any longer
		keybitSizeRandOff = true
		addCnTests()
	}
	// Country tests
	{
		roleVals.Country = []string{"foo"}
		addTests(getCountryCheck(roleVals))

		roleVals.Country = []string{"foo", "bar"}
		addTests(getCountryCheck(roleVals))
	}
	// OU tests
	{
		roleVals.OU = []string{"foo"}
		addTests(getOuCheck(roleVals))

		roleVals.OU = []string{"foo", "bar"}
		addTests(getOuCheck(roleVals))
	}
	// Organization tests
	{
		roleVals.Organization = []string{"system:masters"}
		addTests(getOrganizationCheck(roleVals))

		roleVals.Organization = []string{"foo", "bar"}
		addTests(getOrganizationCheck(roleVals))
	}
	// Locality tests
	{
		roleVals.Locality = []string{"foo"}
		addTests(getLocalityCheck(roleVals))

		roleVals.Locality = []string{"foo", "bar"}
		addTests(getLocalityCheck(roleVals))
	}
	// Province tests
	{
		roleVals.Province = []string{"foo"}
		addTests(getProvinceCheck(roleVals))

		roleVals.Province = []string{"foo", "bar"}
		addTests(getProvinceCheck(roleVals))
	}
	// StreetAddress tests
	{
		roleVals.StreetAddress = []string{"123 foo street"}
		addTests(getStreetAddressCheck(roleVals))

		roleVals.StreetAddress = []string{"123 foo street", "456 bar avenue"}
		addTests(getStreetAddressCheck(roleVals))
	}
	// PostalCode tests
	{
		roleVals.PostalCode = []string{"f00"}
		addTests(getPostalCodeCheck(roleVals))

		roleVals.PostalCode = []string{"f00", "b4r"}
		addTests(getPostalCodeCheck(roleVals))
	}
	// NotBefore tests
	{
		roleVals.NotBeforeDuration = 10 * time.Second
		addTests(getNotBeforeCheck(roleVals))

		roleVals.NotBeforeDuration = 30 * time.Second
		addTests(getNotBeforeCheck(roleVals))

		roleVals.NotBeforeDuration = 0
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
		roleVals.MaxTTL = 0
		addTests(nil)

		roleVals.Lease = "12h"
		roleVals.MaxTTL = 6 * time.Hour
		addTests(nil)

		roleTestStep.ErrorOk = false
		roleVals.TTL = 0
		roleVals.MaxTTL = 12 * time.Hour
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

	b := Backend(config)
	err := b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// generate root
	rootData := map[string]interface{}{
		"common_name": "test.com",
		"ttl":         "6h",
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
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

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
		resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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

	b := Backend(config)
	err := b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// generate root
	rootData := map[string]interface{}{
		"common_name": "test.com",
		"ttl":         "172800",
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
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
	pemCSR := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	})))
	if len(pemCSR) == 0 {
		t.Fatal("pem csr is empty")
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sign-verbatim",
		Storage:   storage,
		Data: map[string]interface{}{
			"csr": pemCSR,
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
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sign-verbatim/test",
		Storage:   storage,
		Data: map[string]interface{}{
			"csr": pemCSR,
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
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sign-verbatim/test",
		Storage:   storage,
		Data: map[string]interface{}{
			"csr": pemCSR,
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
		t.Fatalf("sign-verbatim did not properly cap validity period on signed CSR")
	}

	// now check that if we set generate-lease it takes it from the role and the TTLs match
	roleData = map[string]interface{}{
		"ttl":            "4h",
		"max_ttl":        "8h",
		"generate_lease": true,
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sign-verbatim/test",
		Storage:   storage,
		Data: map[string]interface{}{
			"csr": pemCSR,
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

func TestBackend_Root_Idempotency(t *testing.T) {
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
	if resp == nil {
		t.Fatal("expected a warning")
	}
	if resp.Data != nil || len(resp.Warnings) == 0 {
		t.Fatalf("bad response: %#v", *resp)
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

	b := Backend(config)
	err := b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// generate root
	rootData := map[string]interface{}{
		"common_name": "test.com",
		"ttl":         "172800",
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
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
		pemSS := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: selfSigned,
		})))
		return pemSS, cert
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "foo.bar.com",
		},
		SerialNumber:          big.NewInt(1234),
		IsCA:                  false,
		BasicConstraintsValid: true,
	}

	ss, _ := getSelfSigned(template, template)
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
		SerialNumber:          big.NewInt(2345),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	ss, ssCert := getSelfSigned(template, issuer)
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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

	signingBundle, err := fetchCAInfo(context.Background(), &logical.Request{Storage: storage})
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

// This is a really tricky test because the Go stdlib asn1 package is incapable
// of doing the right thing with custom OID SANs (see comments in the package,
// it's readily admitted that it's too magic) but that means that any
// validation logic written for this test isn't being independently verified,
// as in, if cryptobytes is used to decode it to make the test work, that
// doesn't mean we're encoding and decoding correctly, only that we made the
// test pass. Instead, when run verbosely it will first perform a bunch of
// checks to verify that the OID SAN logic doesn't screw up other SANs, then
// will spit out the PEM. This can be validated independently.
//
// You want the hex dump of the octet string corresponding to the X509v3
// Subject Alternative Name. There's a nice online utility at
// https://lapo.it/asn1js that can be used to view the structure of an
// openssl-generated other SAN at
// https://lapo.it/asn1js/#3022A020060A2B060104018237140203A0120C106465766F7073406C6F63616C686F7374
// (openssl asn1parse can also be used with -strparse using an offset of the
// hex blob for the subject alternative names extension).
//
// The structure output from here should match that precisely (even if the OID
// itself doesn't) in the second test.
//
// The test that encodes two should have them be in separate elements in the
// top-level sequence; see
// https://lapo.it/asn1js/#3046A020060A2B060104018237140203A0120C106465766F7073406C6F63616C686F7374A022060A2B060104018237140204A0140C12322D6465766F7073406C6F63616C686F7374 for an openssl-generated example.
//
// The good news is that it's valid to simply copy and paste the PEM output from
// here into the form at that site as it will do the right thing so it's pretty
// easy to validate.
func TestBackend_OID_SANs(t *testing.T) {
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

	var resp *api.Secret
	var certStr string
	var block *pem.Block
	var cert *x509.Certificate

	_, err = client.Logical().Write("root/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("root/roles/test", map[string]interface{}{
		"allowed_domains":    []string{"foobar.com", "zipzap.com"},
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"allow_ip_sans":      true,
		"allowed_other_sans": "1.3.6.1.4.1.311.20.2.3;UTF8:devops@*,1.3.6.1.4.1.311.20.2.4;utf8:d*e@foobar.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get a baseline before adding OID SANs. In the next sections we'll verify
	// that the SANs are all added even as the OID SAN inclusion forces other
	// adding logic (custom rather than built-in Golang logic)
	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
	})
	if err != nil {
		t.Fatal(err)
	}
	certStr = resp.Data["certificate"].(string)
	block, _ = pem.Decode([]byte(certStr))
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if cert.IPAddresses[0].String() != "1.2.3.4" {
		t.Fatalf("unexpected IP SAN %q", cert.IPAddresses[0].String())
	}
	if cert.DNSNames[0] != "foobar.com" ||
		cert.DNSNames[1] != "bar.foobar.com" ||
		cert.DNSNames[2] != "foo.foobar.com" {
		t.Fatalf("unexpected DNS SANs %v", cert.DNSNames)
	}

	// First test some bad stuff that shouldn't work
	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		// Not a valid value for the first possibility
		"other_sans": "1.3.6.1.4.1.311.20.2.3;UTF8:devop@nope.com",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		// Not a valid OID for the first possibility
		"other_sans": "1.3.6.1.4.1.311.20.2.5;UTF8:devops@nope.com",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		// Not a valid name for the second possibility
		"other_sans": "1.3.6.1.4.1.311.20.2.4;UTF8:d34g@foobar.com",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		// Not a valid OID for the second possibility
		"other_sans": "1.3.6.1.4.1.311.20.2.5;UTF8:d34e@foobar.com",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		// Not a valid type
		"other_sans": "1.3.6.1.4.1.311.20.2.5;UTF2:d34e@foobar.com",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// Valid for first possibility
	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		"other_sans":  "1.3.6.1.4.1.311.20.2.3;utf8:devops@nope.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	certStr = resp.Data["certificate"].(string)
	block, _ = pem.Decode([]byte(certStr))
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if cert.IPAddresses[0].String() != "1.2.3.4" {
		t.Fatalf("unexpected IP SAN %q", cert.IPAddresses[0].String())
	}
	if cert.DNSNames[0] != "foobar.com" ||
		cert.DNSNames[1] != "bar.foobar.com" ||
		cert.DNSNames[2] != "foo.foobar.com" {
		t.Fatalf("unexpected DNS SANs %v", cert.DNSNames)
	}
	t.Logf("certificate 1 to check:\n%s", certStr)

	// Valid for second possibility
	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		"other_sans":  "1.3.6.1.4.1.311.20.2.4;UTF8:d234e@foobar.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	certStr = resp.Data["certificate"].(string)
	block, _ = pem.Decode([]byte(certStr))
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if cert.IPAddresses[0].String() != "1.2.3.4" {
		t.Fatalf("unexpected IP SAN %q", cert.IPAddresses[0].String())
	}
	if cert.DNSNames[0] != "foobar.com" ||
		cert.DNSNames[1] != "bar.foobar.com" ||
		cert.DNSNames[2] != "foo.foobar.com" {
		t.Fatalf("unexpected DNS SANs %v", cert.DNSNames)
	}
	t.Logf("certificate 2 to check:\n%s", certStr)

	// Valid for both
	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		"other_sans":  "1.3.6.1.4.1.311.20.2.3;utf8:devops@nope.com,1.3.6.1.4.1.311.20.2.4;utf8:d234e@foobar.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	certStr = resp.Data["certificate"].(string)
	block, _ = pem.Decode([]byte(certStr))
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if cert.IPAddresses[0].String() != "1.2.3.4" {
		t.Fatalf("unexpected IP SAN %q", cert.IPAddresses[0].String())
	}
	if cert.DNSNames[0] != "foobar.com" ||
		cert.DNSNames[1] != "bar.foobar.com" ||
		cert.DNSNames[2] != "foo.foobar.com" {
		t.Fatalf("unexpected DNS SANs %v", cert.DNSNames)
	}
	t.Logf("certificate 3 to check:\n%s", certStr)
}

func TestBackend_AllowedSerialNumbers(t *testing.T) {
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

	var resp *api.Secret
	var certStr string
	var block *pem.Block
	var cert *x509.Certificate

	_, err = client.Logical().Write("root/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// First test that Serial Numbers are not allowed
	_, err = client.Logical().Write("root/roles/test", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar",
		"ttl":         "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name":   "foobar",
		"ttl":           "1h",
		"serial_number": "foobar",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// Update the role to allow serial numbers
	_, err = client.Logical().Write("root/roles/test", map[string]interface{}{
		"allow_any_name":         true,
		"enforce_hostnames":      false,
		"allowed_serial_numbers": "f00*,b4r*",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar",
		"ttl":         "1h",
		// Not a valid serial number
		"serial_number": "foobar",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// Valid for first possibility
	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name":   "foobar",
		"serial_number": "f00bar",
	})
	if err != nil {
		t.Fatal(err)
	}
	certStr = resp.Data["certificate"].(string)
	block, _ = pem.Decode([]byte(certStr))
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if cert.Subject.SerialNumber != "f00bar" {
		t.Fatalf("unexpected Subject SerialNumber %s", cert.Subject.SerialNumber)
	}
	t.Logf("certificate 1 to check:\n%s", certStr)

	// Valid for second possibility
	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name":   "foobar",
		"serial_number": "b4rf00",
	})
	if err != nil {
		t.Fatal(err)
	}
	certStr = resp.Data["certificate"].(string)
	block, _ = pem.Decode([]byte(certStr))
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if cert.Subject.SerialNumber != "b4rf00" {
		t.Fatalf("unexpected Subject SerialNumber %s", cert.Subject.SerialNumber)
	}
	t.Logf("certificate 2 to check:\n%s", certStr)
}

func TestBackend_URI_SANs(t *testing.T) {
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

	_, err = client.Logical().Write("root/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("root/roles/test", map[string]interface{}{
		"allowed_domains":    []string{"foobar.com", "zipzap.com"},
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"allow_ip_sans":      true,
		"allowed_uri_sans":   []string{"http://someuri/abc", "spiffe://host.com/*"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// First test some bad stuff that shouldn't work
	_, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		"uri_sans":    "http://www.mydomain.com/zxf",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// Test valid single entry
	_, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		"uri_sans":    "http://someuri/abc",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test globed entry
	_, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		"uri_sans":    "spiffe://host.com/something",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test multiple entries
	resp, err := client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		"uri_sans":    "spiffe://host.com/something,http://someuri/abc",
	})
	if err != nil {
		t.Fatal(err)
	}

	certStr := resp.Data["certificate"].(string)
	block, _ := pem.Decode([]byte(certStr))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	URI0, _ := url.Parse("spiffe://host.com/something")
	URI1, _ := url.Parse("http://someuri/abc")

	if len(cert.URIs) != 2 {
		t.Fatalf("expected 2 valid URIs SANs %v", cert.URIs)
	}

	if cert.URIs[0].String() != URI0.String() || cert.URIs[1].String() != URI1.String() {
		t.Fatalf(
			"expected URIs SANs %v to equal provided values spiffe://host.com/something, http://someuri/abc",
			cert.URIs)
	}
}
func setCerts() {
	cak, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	marshaledKey, err := x509.MarshalECPrivateKey(cak)
	if err != nil {
		panic(err)
	}
	keyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledKey,
	}
	ecCAKey = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
	if err != nil {
		panic(err)
	}
	subjKeyID, err := certutil.GetSubjKeyID(cak)
	if err != nil {
		panic(err)
	}
	caCertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "root.localhost",
		},
		SubjectKeyId:          subjKeyID,
		DNSNames:              []string{"root.localhost"},
		KeyUsage:              x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign),
		SerialNumber:          big.NewInt(mathrand.Int63()),
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, cak.Public(), cak)
	if err != nil {
		panic(err)
	}
	caCertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	ecCACert = strings.TrimSpace(string(pem.EncodeToMemory(caCertPEMBlock)))

	rak, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	marshaledKey = x509.MarshalPKCS1PrivateKey(rak)
	keyPEMBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: marshaledKey,
	}
	rsaCAKey = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
	if err != nil {
		panic(err)
	}
	subjKeyID, err = certutil.GetSubjKeyID(rak)
	if err != nil {
		panic(err)
	}
	caBytes, err = x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, rak.Public(), rak)
	if err != nil {
		panic(err)
	}
	caCertPEMBlock = &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	rsaCACert = strings.TrimSpace(string(pem.EncodeToMemory(caCertPEMBlock)))
}

var (
	initTest  sync.Once
	rsaCAKey  string
	rsaCACert string
	ecCAKey   string
	ecCACert  string
)
