package pki

import (
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
	mathrand "math/rand"
	"net"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

var (
	stepCount = 0
)

// Performs basic tests on CA functionality
// Uses the RSA CA key
func TestBackend_RSAKey(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 30
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

	stepCount += len(testCase.Steps)

	intdata := map[string]interface{}{}
	reqdata := map[string]interface{}{}
	testCase.Steps = append(testCase.Steps, generateCATestingSteps(t, rsaCACert, rsaCAKey, intdata, reqdata)...)

	logicaltest.Test(t, testCase)
}

// Performs basic tests on CA functionality
// Uses the EC CA key
func TestBackend_ECKey(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 30
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

	stepCount += len(testCase.Steps)

	intdata := map[string]interface{}{}
	reqdata := map[string]interface{}{}
	testCase.Steps = append(testCase.Steps, generateCATestingSteps(t, ecCACert, ecCAKey, intdata, reqdata)...)

	logicaltest.Test(t, testCase)
}

func TestBackend_CSRValues(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 30
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

	stepCount += len(testCase.Steps)

	intdata := map[string]interface{}{}
	reqdata := map[string]interface{}{}
	testCase.Steps = append(testCase.Steps, generateCSRSteps(t, ecCACert, ecCAKey, intdata, reqdata)...)

	logicaltest.Test(t, testCase)
}

func TestBackend_URLsCRUD(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 30
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

	stepCount += len(testCase.Steps)

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
	maxLeaseTTLVal := time.Hour * 24 * 30
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
				Operation: logical.WriteOperation,
				Path:      "config/ca/set",
				Data: map[string]interface{}{
					"pem_bundle": rsaCAKey + rsaCACert,
				},
			},
		},
	}

	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, false)...)
	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, true)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+stepCount, v)
		}
	}

	stepCount += len(testCase.Steps)

	logicaltest.Test(t, testCase)
}

// Generates and tests steps that walk through the various possibilities
// of role flags to ensure that they are properly restricted
// Uses the EC CA key
func TestBackend_ECRoles(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 30
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
				Operation: logical.WriteOperation,
				Path:      "config/ca/set",
				Data: map[string]interface{}{
					"pem_bundle": ecCAKey + ecCACert,
				},
			},
		},
	}

	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, false)...)
	testCase.Steps = append(testCase.Steps, generateRoleSteps(t, true)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+stepCount, v)
		}
	}

	stepCount += len(testCase.Steps)

	logicaltest.Test(t, testCase)
}

// Performs some validity checking on the returned bundles
func checkCertsAndPrivateKey(keyType string, key crypto.Signer, usage certUsage, validity time.Duration, certBundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
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
	case parsedCertBundle.IssuingCA == nil:
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
	// There should only be one usage type, because only one is requested
	// in the tests
	if len(cert.ExtKeyUsage) != 1 {
		return nil, fmt.Errorf("Got wrong size key usage in generated cert; values are %#v", cert.ExtKeyUsage)
	}
	switch usage {
	case emailProtectionUsage:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageEmailProtection {
			return nil, fmt.Errorf("Bad key usage")
		}
	case serverUsage:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageServerAuth {
			return nil, fmt.Errorf("Bad key usage")
		}
	case clientUsage:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth {
			return nil, fmt.Errorf("Bad key usage")
		}
	case codeSigningUsage:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageCodeSigning {
			return nil, fmt.Errorf("Bad key usage")
		}
	}

	if math.Abs(float64(time.Now().Unix()-cert.NotBefore.Unix())) > 10 {
		return nil, fmt.Errorf("Validity period starts out of range")
	}

	if math.Abs(float64(time.Now().Add(validity).Unix()-cert.NotAfter.Unix())) > 10 {
		return nil, fmt.Errorf("Validity period of %d too large vs max of 10", cert.NotAfter.Unix())
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

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	csr, _ := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, priv)
	csrPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	})

	ret := []logicaltest.TestStep{
		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/generate/root/exported",
			Data: map[string]interface{}{
				"common_name": "Root Cert",
				"ttl":         "180h",
			},
		},

		logicaltest.TestStep{
			Operation: logical.WriteOperation,
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
			Operation: logical.WriteOperation,
			Path:      "config/ca/sign",
			Data: map[string]interface{}{
				"common_name": "Intermediate Cert",
				"csr":         string(csrPem),
				"format":      "der",
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

				switch {
				case !reflect.DeepEqual(expected.IssuingCertificates, cert.IssuingCertificateURL):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", expected.IssuingCertificates, cert.IssuingCertificateURL)
				case !reflect.DeepEqual(expected.CRLDistributionPoints, cert.CRLDistributionPoints):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", expected.CRLDistributionPoints, cert.CRLDistributionPoints)
				case !reflect.DeepEqual(expected.OCSPServers, cert.OCSPServer):
					return fmt.Errorf("expected\n%#v\ngot\n%#v\n", expected.OCSPServers, cert.OCSPServer)
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
			Operation: logical.WriteOperation,
			Path:      "config/ca/generate/root/exported",
			Data: map[string]interface{}{
				"common_name":     "Root Cert",
				"ttl":             "180h",
				"max_path_length": 0,
			},
		},

		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/sign",
			Data: map[string]interface{}{
				"use_csr_values": true,
				"csr":            string(csrPem),
				"format":         "der",
			},
			ErrorOk: true,
		},

		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/generate/root/exported",
			Data: map[string]interface{}{
				"common_name":     "Root Cert",
				"ttl":             "180h",
				"max_path_length": 1,
			},
		},
		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/sign",
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
func generateCATestingSteps(t *testing.T, caCert, caKey string, intdata, reqdata map[string]interface{}) []logicaltest.TestStep {
	ret := []logicaltest.TestStep{
		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/set",
			Data: map[string]interface{}{
				"pem_bundle": caKey + caCert,
			},
		},

		logicaltest.TestStep{
			Operation: logical.WriteOperation,
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

		// Now test uploading when the private key is already stored, such
		// as when uploading a CA signed as the result of a generated CSR
		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/set",
			Data: map[string]interface{}{
				"pem_bundle": caCert,
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

		// Test a bunch of generation stuff
		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/generate/root/exported",
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
			Operation: logical.WriteOperation,
			Path:      "config/ca/generate/intermediate/exported",
			Data: map[string]interface{}{
				"common_name": "Intermediate Cert",
			},
			Check: func(resp *logical.Response) error {
				intdata["intermediatecsr"] = resp.Data["csr"].(string)
				intdata["intermediatekey"] = resp.Data["private_key"].(string)
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/set",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				delete(reqdata, "pem_bundle")
				delete(reqdata, "ttl")
				reqdata["csr"] = intdata["intermediatecsr"].(string)
				reqdata["common_name"] = "Intermediate Cert"
				reqdata["ttl"] = "90h"
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/sign",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				intdata["intermediatecert"] = resp.Data["certificate"].(string)
				delete(reqdata, "csr")
				delete(reqdata, "common_name")
				delete(reqdata, "ttl")
				reqdata["serial_number"] = resp.Data["serial_number"].(string)
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.WriteOperation,
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
				revokedString := certutil.GetOctalFormatted(revokedList[0].SerialNumber.Bytes(), ":")
				if revokedString != reqdata["serial_number"].(string) {
					t.Fatalf("got serial %s, expecting %s", revokedString, reqdata["serial_number"].(string))
				}
				delete(reqdata, "serial_number")
				return nil
			},
		},

		// Do it all again, with EC keys and DER format
		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/generate/root/exported",
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
			Operation: logical.WriteOperation,
			Path:      "config/ca/generate/intermediate/exported",
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
			Operation: logical.WriteOperation,
			Path:      "config/ca/set",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				delete(reqdata, "pem_bundle")
				delete(reqdata, "ttl")
				reqdata["csr"] = intdata["intermediatecsr"].(string)
				reqdata["common_name"] = "Intermediate Cert"
				reqdata["ttl"] = "90h"
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca/sign",
			Data:      reqdata,
			Check: func(resp *logical.Response) error {
				intdata["intermediatecert"] = resp.Data["certificate"].(string)
				delete(reqdata, "csr")
				delete(reqdata, "common_name")
				delete(reqdata, "ttl")
				reqdata["serial_number"] = resp.Data["serial_number"].(string)
				return nil
			},
		},

		logicaltest.TestStep{
			Operation: logical.WriteOperation,
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
					revokedString := certutil.GetOctalFormatted(revEntry.SerialNumber.Bytes(), ":")
					if revokedString == reqdata["serial_number"].(string) {
						found = true

					}
				}
				if !found {
					t.Fatalf("did not find %s in CRL", reqdata["serial_number"].(string))
				}
				delete(reqdata, "serial_number")
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
		Operation: logical.WriteOperation,
		Path:      "roles/test",
	}
	var issueTestStep logicaltest.TestStep
	if useCSRs {
		issueTestStep = logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "sign/test",
		}
	} else {
		issueTestStep = logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "issue/test",
		}
	}

	genericErrorOkCheck := func(resp *logical.Response) error {
		if resp.IsError() {
			return nil
		}
		return fmt.Errorf("Expected an error, but did not seem to get one")
	}

	// Adds tests with the currently configured issue/role information
	addTests := func(testCheck logicaltest.TestCheckFunc) {
		//fmt.Printf("role vals: %#v\n", roleVals)
		//fmt.Printf("issue vals: %#v\n", issueTestStep)
		roleTestStep.Data = structs.New(roleVals).Map()
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

	// Returns a TestCheckFunc that performs various validity checks on the
	// returned certificate information, mostly within checkCertsAndPrivateKey
	getCnCheck := func(name, keyType string, key crypto.Signer, usage certUsage, validity time.Duration) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			parsedCertBundle, err := checkCertsAndPrivateKey(keyType, key, usage, validity, &certBundle)
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
		BaseDomain           bool `structs:"example.com"`
		SubDomain            bool `structs:"foo.example.com"`
		Wildcard             bool `structs:"*.example.com"`
		SubSubdomain         bool `structs:"foo.bar.example.com"`
		SubSubdomainWildcard bool `structs:"*.bar.example.com"`
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
		// For the number of tests being run, this is known to hit all
		// of the various values below
		mathRand := mathrand.New(mathrand.NewSource(1))
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
			var usage certUsage
			i := mathRand.Int()
			switch {
			case i%5 == 0:
				usage = emailProtectionUsage
				roleVals.EmailProtectionFlag = true
			case i%3 == 0:
				usage = serverUsage
				roleVals.ServerFlag = true
			case i%2 == 0:
				usage = clientUsage
				roleVals.ClientFlag = true
			default:
				usage = codeSigningUsage
				roleVals.CodeSigningFlag = true
			}

			allowed := allowedInt.(bool)
			issueVals.CommonName = name
			if roleVals.EmailProtectionFlag {
				issueVals.CommonName = "user@" + issueVals.CommonName
			}
			if allowed {
				issueTestStep.ErrorOk = false
			} else {
				issueTestStep.ErrorOk = true
			}

			validity, _ := time.ParseDuration(roleVals.MaxTTL)
			if useCSRs {
				var privKey crypto.Signer
				switch roleVals.KeyType {
				case "rsa":
					privKey, _ = rsa.GenerateKey(rand.Reader, roleVals.KeyBits)
				case "ec":
					var curve elliptic.Curve
					switch roleVals.KeyBits {
					case 224:
						curve = elliptic.P224()
					case 256:
						curve = elliptic.P256()
					case 384:
						curve = elliptic.P384()
					case 521:
						curve = elliptic.P521()
					}
					privKey, _ = ecdsa.GenerateKey(curve, rand.Reader)
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
				addTests(getCnCheck(issueVals.CommonName, roleVals.KeyType, privKey, usage, validity))
			} else {
				addTests(getCnCheck(issueVals.CommonName, roleVals.KeyType, nil, usage, validity))
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

		roleVals.AllowedBaseDomain = "foobar.com"
		addCnTests()

		roleVals.AllowedBaseDomain = "example.com"
		roleVals.AllowSubdomains = true
		commonNames.SubDomain = true
		commonNames.Wildcard = true
		commonNames.SubSubdomain = true
		commonNames.SubSubdomainWildcard = true
		addCnTests()

		roleVals.AllowBaseDomain = true
		commonNames.BaseDomain = true
		addCnTests()

		roleVals.AllowAnyName = true
		roleVals.EnforceHostnames = true
		commonNames.AnyHost = true
		addCnTests()

		roleVals.EnforceHostnames = false
		commonNames.NonHostname = true
		addCnTests()
	}
	// IP SAN tests
	{
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
	}

	return ret
}

const (
	rsaCAKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1eKB2nFbRqTFs7KyZjbzB5VRCBbnLZfEXVP1c3bHe+YGjlfl
34cy52dmancUzOf1/Jfo+VglocjTLVy5wHSGJwQYs8b6pEuuvAVo/6wUL5Z7ZlQD
R4kDe5Q+xgoRT6Bi/Bs57E+fNYgyUq/YAUY5WLuC+ZliCbJkLnb15ItuP1yVUTDX
TYORRE3qJS5RRol8D3QvteG9LyPEc7C+jsm5iBCagyxluzU0dnEOib5q7xwZncoM
bQz+rZH3QnwOij41FOGRPazrD5Mv6xLBkFnE5VAJ+GIgvd4bpOwvYMuofvF4PS7S
FzxkGssMLlICap6PFpKz86DpAoDxPuoZeOhU4QIDAQABAoIBAQCp6VIdFdZcDYPd
WIVuvBJfINiJo6AtURa2yX8BJggdPkRRCjTcWUwwFq1+wHDuwwtgidGTW9oxZxeU
Psh1wlvcXN2+28C7ikAar/WUvsAeed44EV+1kXwJzV/89XyBFDnuazadqzcgUL0h
gP4JLR9bhULsRFRkvanmW6zFzZpcjBzi/UoFuWkFRRqZ0euM2Lpz8L75PFfW9s9M
kNglZpcV6ZmvR9c1JkEMUs/mrB8ZgCd1uvmcVosQ+u7sE8Yk/xAurHXuNJQlGXx4
azrLW0XY1CLO2Tm4l4MwPjmhH0WytXNjOSKycBCXVnBIfZsI128DsP5YyA/fW9qA
BAqFSzABAoGBAPcBNk9sf3cnZ5w6qwlE2ysDwGIGR+I1fb09YjRI6vjwwdWZgGR0
EE4UB1Pp+KIehXaTJHcEgvBBErM2NLS4qKzh25O30C2EwK6o//3jEAribuYutBhJ
ihu1qKzqcPbKClG+34kjX6nmtux2wlYM05f5v3ALki5Is7W/RrfceBuBAoGBAN2s
hdt4TcgIcZymPG2931qCBGF3E8AaA8bUl9TKaZHuFikOMFKA/KM5O5mznPGnQP2d
kXYKXuqdYhVLwp32FTbIbozGZZ8XliO5oS7J3vIID+sLWQhrvyFO7d0lpSjv41HH
yJ2DrykHRg8hxsbh2D4By7olBx6Q2m+B8lPzHmlhAoGACHUeKvIIG0haH9tSZ+rX
pk1mlPSqGXDDcWtcpXWptgRoXqv23Xmr5UCCT7k/Li3lW/4FzZ117kwMG97LRzTb
ca/6GMC+fBCDmHdo7ISN1BGUwoTu3bYG6JP7xo/wdkLMv6fNd6CicerYcJhQZynh
RN7kUy3SP4t1u89k2H7QDgECgYBpU0bKr8+tQq3Qs3+02OmeFHbGZJDCztmKiIqX
tZERoGFxIme9W8IuP8xczGW+wCx2FH7/6g+NRDhNTBDtgvYzcGpugvnX7JoO4W1/
ULWYpFID6QFlqeRHjDwivndKCykkO1vL07zPLsCQAglzh+16ENpe2KcYU9Ul9EVS
tAp4IQKBgQDrb/NpiVx7NI6PyTCm6ctuUAYm3ihAiQNV4Bmr0liPDp9PozbqkhcF
udNtivO4LlRb/PJ+DK6afDyH8aJQdDqe3NpDvyrmKiMSYOY3iVFvan4tbIiofxdQ
flwiZUzox814fzXbxheO9Cs6pXz7PUBVU4fN0Y/hXJCfRO4Ns9152A==
-----END RSA PRIVATE KEY-----
`
	rsaCACert string = `-----BEGIN CERTIFICATE-----
MIIDUTCCAjmgAwIBAgIJAKM+z4MSfw2mMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV
BAMMEFZhdWx0IFRlc3RpbmcgQ0EwHhcNMTUwNjAxMjA1MTUzWhcNMjUwNTI5MjA1
MTUzWjAbMRkwFwYDVQQDDBBWYXVsdCBUZXN0aW5nIENBMIIBIjANBgkqhkiG9w0B
AQEFAAOCAQ8AMIIBCgKCAQEA1eKB2nFbRqTFs7KyZjbzB5VRCBbnLZfEXVP1c3bH
e+YGjlfl34cy52dmancUzOf1/Jfo+VglocjTLVy5wHSGJwQYs8b6pEuuvAVo/6wU
L5Z7ZlQDR4kDe5Q+xgoRT6Bi/Bs57E+fNYgyUq/YAUY5WLuC+ZliCbJkLnb15Itu
P1yVUTDXTYORRE3qJS5RRol8D3QvteG9LyPEc7C+jsm5iBCagyxluzU0dnEOib5q
7xwZncoMbQz+rZH3QnwOij41FOGRPazrD5Mv6xLBkFnE5VAJ+GIgvd4bpOwvYMuo
fvF4PS7SFzxkGssMLlICap6PFpKz86DpAoDxPuoZeOhU4QIDAQABo4GXMIGUMB0G
A1UdDgQWBBTknN5eFxxo5aTlfq+G4ZXs3AsxWTAfBgNVHSMEGDAWgBTknN5eFxxo
5aTlfq+G4ZXs3AsxWTAxBgNVHR8EKjAoMCagJKAihiBodHRwOi8vbG9jYWxob3N0
OjgyMDAvdjEvcGtpL2NybDAPBgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQEAwIB
BjANBgkqhkiG9w0BAQsFAAOCAQEAsINcA4PZm+OyldgNrwRVgxoSrhV1I9zszhc9
VV340ZWlpTTxFKVb/K5Hg+jMF9tv70X1HwlYdlutE6KdrsA3gks5zanh4/3zlrYk
ABNBmSD6SSU2HKX1bFCBAAS3YHONE5o1K5tzwLsMl5uilNf+Wid3NjFnQ4KfuYI5
loN/opnM6+a/O3Zua8RAuMMAv9wyqwn88aVuLvVzDNSMe5qC5kkuLGmRkNgY06rI
S/fXIHIOldeQxgYCqhdVmcDWJ1PtVaDfBsKVpRg1GRU8LUGw2E4AY+twd+J2FBfa
G/7g4koczXLoUM3OQXd5Aq2cs4SS1vODrYmgbioFsQ3eDHd1fg==
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
