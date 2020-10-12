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
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/go-test/deep"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
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
func TestBackend_Roles(t *testing.T) {
	cases := []struct {
		name      string
		key, cert *string
		useCSR    bool
	}{
		{"RSA", &rsaCAKey, &rsaCACert, false},
		{"RSACSR", &rsaCAKey, &rsaCACert, true},
		{"EC", &ecCAKey, &ecCACert, false},
		{"ECCSR", &ecCAKey, &ecCACert, true},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
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
							"pem_bundle": *tc.key + "\n" + *tc.cert,
						},
					},
				},
			}

			testCase.Steps = append(testCase.Steps, generateRoleSteps(t, tc.useCSR)...)
			if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
				for i, v := range testCase.Steps {
					data := map[string]interface{}{}
					var keys []string
					for k := range v.Data {
						keys = append(keys, k)
					}
					sort.Strings(keys)
					for _, k := range keys {
						interf := v.Data[k]
						switch v := interf.(type) {
						case bool:
							if !v {
								continue
							}
						case int:
							if v == 0 {
								continue
							}
						case []string:
							if len(v) == 0 {
								continue
							}
						case string:
							if v == "" {
								continue
							}
							lines := strings.Split(v, "\n")
							if len(lines) > 1 {
								data[k] = lines[0] + " ... (truncated)"
								continue
							}
						}
						data[k] = interf

					}
					t.Logf("Step %d:\n%s %s err=%v %+v\n\n", i+1, v.Operation, v.Path, v.ErrorOk, data)
				}
			}

			logicaltest.Test(t, testCase)
		})
	}
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
	expected := certutil.URLEntries{
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
				var entries certutil.URLEntries
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
	// t.Logf("seed under test: %v", seed)

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

			expected := strutil.RemoveDuplicatesStable(role.OU, true)
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

	type csrPlan struct {
		errorOk     bool
		roleKeyBits int
		cert        string
		privKey     crypto.Signer
	}

	getCsr := func(keyType string, keyBits int, csrTemplate *x509.CertificateRequest) (*pem.Block, crypto.Signer) {
		var privKey crypto.Signer
		var ok bool
		switch keyType {
		case "rsa":
			privKey, ok = generatedRSAKeys[keyBits]
			if !ok {
				privKey, _ = rsa.GenerateKey(rand.Reader, keyBits)
				generatedRSAKeys[keyBits] = privKey
			}

		case "ec":
			var curve elliptic.Curve

			switch keyBits {
			case 224:
				curve = elliptic.P224()
			case 256:
				curve = elliptic.P256()
			case 384:
				curve = elliptic.P384()
			case 521:
				curve = elliptic.P521()
			}

			privKey, ok = generatedECKeys[keyBits]
			if !ok {
				privKey, _ = ecdsa.GenerateKey(curve, rand.Reader)
				generatedECKeys[keyBits] = privKey
			}

		default:
			panic("invalid key type: " + keyType)
		}

		csr, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, privKey)
		if err != nil {
			t.Fatalf("Error creating certificate request: %s", err)
		}
		block := pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: csr,
		}
		return &block, privKey
	}

	getRandCsr := func(keyType string, errorOk bool, csrTemplate *x509.CertificateRequest) csrPlan {
		rsaKeyBits := []int{2048, 4096}
		ecKeyBits := []int{224, 256, 384, 521}
		var plan = csrPlan{errorOk: errorOk}

		var testBitSize int
		switch keyType {
		case "rsa":
			plan.roleKeyBits = rsaKeyBits[mathRand.Int()%2]
			testBitSize = plan.roleKeyBits

			// If we don't expect an error already, randomly choose a
			// key size and expect an error if it's less than the role
			// setting
			if !keybitSizeRandOff && !errorOk {
				testBitSize = rsaKeyBits[mathRand.Int()%2]
			}

			if testBitSize < plan.roleKeyBits {
				plan.errorOk = true
			}

		case "ec":
			plan.roleKeyBits = ecKeyBits[mathRand.Int()%4]
			testBitSize = plan.roleKeyBits

			// If we don't expect an error already, randomly choose a
			// key size and expect an error if it's less than the role
			// setting
			if !keybitSizeRandOff && !errorOk {
				testBitSize = ecKeyBits[mathRand.Int()%4]
			}

			if testBitSize < plan.roleKeyBits {
				plan.errorOk = true
			}

		default:
			panic("invalid key type: " + keyType)
		}
		if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
			t.Logf("roleKeyBits=%d testBitSize=%d errorOk=%v", plan.roleKeyBits, testBitSize, plan.errorOk)
		}

		block, privKey := getCsr(keyType, testBitSize, csrTemplate)
		plan.cert = strings.TrimSpace(string(pem.EncodeToMemory(block)))
		plan.privKey = privKey
		return plan
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

			if useCSRs {
				templ := &x509.CertificateRequest{
					Subject: pkix.Name{
						CommonName: issueVals.CommonName,
					},
				}
				plan := getRandCsr(roleVals.KeyType, issueTestStep.ErrorOk, templ)
				issueVals.CSR = plan.cert
				roleVals.KeyBits = plan.roleKeyBits
				issueTestStep.ErrorOk = plan.errorOk

				addTests(getCnCheck(issueVals.CommonName, roleVals, plan.privKey, x509.KeyUsage(parsedKeyUsage), extUsage, validity))
			} else {
				addTests(getCnCheck(issueVals.CommonName, roleVals, nil, x509.KeyUsage(parsedKeyUsage), extUsage, validity))
			}
		}
	}

	funcs := []interface{}{addCnTests, getCnCheck, getCountryCheck, getLocalityCheck, getNotBeforeCheck,
		getOrganizationCheck, getOuCheck, getPostalCodeCheck, getRandCsr, getStreetAddressCheck,
		getProvinceCheck}
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		t.Logf("funcs=%d", len(funcs))
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
		getIpCheck := func(expectedIp ...net.IP) logicaltest.TestCheckFunc {
			return func(resp *logical.Response) error {
				var certBundle certutil.CertBundle
				err := mapstructure.Decode(resp.Data, &certBundle)
				if err != nil {
					return err
				}
				parsedCertBundle, err := certBundle.ToParsedCertBundle()
				if err != nil {
					return fmt.Errorf("error parsing cert bundle: %s", err)
				}
				cert := parsedCertBundle.Certificate
				var emptyIPs []net.IP
				var expected []net.IP = append(emptyIPs, expectedIp...)
				if diff := deep.Equal(cert.IPAddresses, expected); len(diff) > 0 {
					return fmt.Errorf("wrong SAN IPs, diff: %v", diff)
				}
				return nil
			}
		}
		addIPSANTests := func(useCSRs, useCSRSANs, allowIPSANs, errorOk bool, ipSANs string, csrIPSANs []net.IP, check logicaltest.TestCheckFunc) {
			if useCSRs {
				csrTemplate := &x509.CertificateRequest{
					Subject: pkix.Name{
						CommonName: issueVals.CommonName,
					},
					IPAddresses: csrIPSANs,
				}
				block, _ := getCsr(roleVals.KeyType, roleVals.KeyBits, csrTemplate)
				issueVals.CSR = strings.TrimSpace(string(pem.EncodeToMemory(block)))
			}
			oldRoleVals, oldIssueVals, oldIssueTestStep := roleVals, issueVals, issueTestStep
			roleVals.UseCSRSANs = useCSRSANs
			roleVals.AllowIPSANs = allowIPSANs
			issueVals.CommonName = "someone@example.com"
			issueVals.IPSANs = ipSANs
			issueTestStep.ErrorOk = errorOk
			addTests(check)
			roleVals, issueVals, issueTestStep = oldRoleVals, oldIssueVals, oldIssueTestStep
		}
		roleVals.AllowAnyName = true
		roleVals.EnforceHostnames = true
		roleVals.AllowLocalhost = true
		roleVals.UseCSRCommonName = true
		commonNames.Localhost = true

		netip1, netip2 := net.IP{127, 0, 0, 1}, net.IP{170, 171, 172, 173}
		textip1, textip3 := "127.0.0.1", "::1"

		// IPSANs not allowed and not provided, should not be an error.
		addIPSANTests(useCSRs, false, false, false, "", nil, getIpCheck())

		// IPSANs not allowed, valid IPSANs provided, should be an error.
		addIPSANTests(useCSRs, false, false, true, textip1+","+textip3, nil, nil)

		// IPSANs allowed, bogus IPSANs provided, should be an error.
		addIPSANTests(useCSRs, false, true, true, "foobar", nil, nil)

		// Given IPSANs as API argument and useCSRSANs false, CSR arg ignored.
		addIPSANTests(useCSRs, false, true, false, textip1,
			[]net.IP{netip2}, getIpCheck(netip1))

		if useCSRs {
			// IPSANs not allowed, valid IPSANs provided via CSR, should be an error.
			addIPSANTests(useCSRs, true, false, true, "", []net.IP{netip1}, nil)

			// Given IPSANs as both API and CSR arguments and useCSRSANs=true, API arg ignored.
			addIPSANTests(useCSRs, true, true, false, textip3,
				[]net.IP{netip1, netip2}, getIpCheck(netip1, netip2))
		}
	}

	{
		getOtherCheck := func(expectedOthers ...otherNameUtf8) logicaltest.TestCheckFunc {
			return func(resp *logical.Response) error {
				var certBundle certutil.CertBundle
				err := mapstructure.Decode(resp.Data, &certBundle)
				if err != nil {
					return err
				}
				parsedCertBundle, err := certBundle.ToParsedCertBundle()
				if err != nil {
					return fmt.Errorf("error parsing cert bundle: %s", err)
				}
				cert := parsedCertBundle.Certificate
				foundOthers, err := getOtherSANsFromX509Extensions(cert.Extensions)
				var emptyOthers []otherNameUtf8
				var expected []otherNameUtf8 = append(emptyOthers, expectedOthers...)
				if diff := deep.Equal(foundOthers, expected); len(diff) > 0 {
					return fmt.Errorf("wrong SAN IPs, diff: %v", diff)
				}
				return nil
			}
		}

		addOtherSANTests := func(useCSRs, useCSRSANs bool, allowedOtherSANs []string, errorOk bool, otherSANs []string, csrOtherSANs []otherNameUtf8, check logicaltest.TestCheckFunc) {
			otherSansMap := func(os []otherNameUtf8) map[string][]string {
				ret := make(map[string][]string)
				for _, o := range os {
					ret[o.oid] = append(ret[o.oid], o.value)
				}
				return ret
			}
			if useCSRs {
				csrTemplate := &x509.CertificateRequest{
					Subject: pkix.Name{
						CommonName: issueVals.CommonName,
					},
				}
				if err := handleOtherCSRSANs(csrTemplate, otherSansMap(csrOtherSANs)); err != nil {
					t.Fatal(err)
				}
				block, _ := getCsr(roleVals.KeyType, roleVals.KeyBits, csrTemplate)
				issueVals.CSR = strings.TrimSpace(string(pem.EncodeToMemory(block)))
			}
			oldRoleVals, oldIssueVals, oldIssueTestStep := roleVals, issueVals, issueTestStep
			roleVals.UseCSRSANs = useCSRSANs
			roleVals.AllowedOtherSANs = allowedOtherSANs
			issueVals.CommonName = "someone@example.com"
			issueVals.OtherSANs = strings.Join(otherSANs, ",")
			issueTestStep.ErrorOk = errorOk
			addTests(check)
			roleVals, issueVals, issueTestStep = oldRoleVals, oldIssueVals, oldIssueTestStep
		}
		roleVals.AllowAnyName = true
		roleVals.EnforceHostnames = true
		roleVals.AllowLocalhost = true
		roleVals.UseCSRCommonName = true
		commonNames.Localhost = true

		newOtherNameUtf8 := func(s string) (ret otherNameUtf8) {
			pieces := strings.Split(s, ";")
			if len(pieces) == 2 {
				piecesRest := strings.Split(pieces[1], ":")
				if len(piecesRest) == 2 {
					switch strings.ToUpper(piecesRest[0]) {
					case "UTF-8", "UTF8":
						return otherNameUtf8{oid: pieces[0], value: piecesRest[1]}
					}
				}
			}
			t.Fatalf("error parsing otherName: %q", s)
			return
		}
		oid1 := "1.3.6.1.4.1.311.20.2.3"
		oth1str := oid1 + ";utf8:devops@nope.com"
		oth1 := newOtherNameUtf8(oth1str)
		oth2 := otherNameUtf8{oid1, "me@example.com"}
		// allowNone, allowAll := []string{}, []string{oid1 + ";UTF-8:*"}
		allowNone, allowAll := []string{}, []string{"*"}

		// OtherSANs not allowed and not provided, should not be an error.
		addOtherSANTests(useCSRs, false, allowNone, false, nil, nil, getOtherCheck())

		// OtherSANs not allowed, valid OtherSANs provided, should be an error.
		addOtherSANTests(useCSRs, false, allowNone, true, []string{oth1str}, nil, nil)

		// OtherSANs allowed, bogus OtherSANs provided, should be an error.
		addOtherSANTests(useCSRs, false, allowAll, true, []string{"foobar"}, nil, nil)

		// Given OtherSANs as API argument and useCSRSANs false, CSR arg ignored.
		addOtherSANTests(useCSRs, false, allowAll, false, []string{oth1str},
			[]otherNameUtf8{oth2}, getOtherCheck(oth1))

		if useCSRs {
			// OtherSANs not allowed, valid OtherSANs provided via CSR, should be an error.
			addOtherSANTests(useCSRs, true, allowNone, true, nil, []otherNameUtf8{oth1}, nil)

			// Given OtherSANs as both API and CSR arguments and useCSRSANs=true, API arg ignored.
			addOtherSANTests(useCSRs, false, allowAll, false, []string{oth2.String()},
				[]otherNameUtf8{oth1}, getOtherCheck(oth2))
		}
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
		"alt_names":   "foobar.com,foo.foobar.com,bar.foobar.com",
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
	if len(cert.DNSNames) != 3 ||
		cert.DNSNames[0] != "bar.foobar.com" ||
		cert.DNSNames[1] != "foo.foobar.com" ||
		cert.DNSNames[2] != "foobar.com" {
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
	if len(cert.DNSNames) != 3 ||
		cert.DNSNames[0] != "bar.foobar.com" ||
		cert.DNSNames[1] != "foo.foobar.com" ||
		cert.DNSNames[2] != "foobar.com" {
		t.Fatalf("unexpected DNS SANs %v", cert.DNSNames)
	}
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		t.Logf("certificate 1 to check:\n%s", certStr)
	}

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
	if len(cert.DNSNames) != 3 ||
		cert.DNSNames[0] != "bar.foobar.com" ||
		cert.DNSNames[1] != "foo.foobar.com" ||
		cert.DNSNames[2] != "foobar.com" {
		t.Fatalf("unexpected DNS SANs %v", cert.DNSNames)
	}
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		t.Logf("certificate 2 to check:\n%s", certStr)
	}

	// Valid for both
	oid1, type1, val1 := "1.3.6.1.4.1.311.20.2.3", "utf8", "devops@nope.com"
	oid2, type2, val2 := "1.3.6.1.4.1.311.20.2.4", "utf-8", "d234e@foobar.com"
	otherNames := []string{fmt.Sprintf("%s;%s:%s", oid1, type1, val1),
		fmt.Sprintf("%s;%s:%s", oid2, type2, val2)}
	resp, err = client.Logical().Write("root/issue/test", map[string]interface{}{
		"common_name": "foobar.com",
		"ip_sans":     "1.2.3.4",
		"alt_names":   "foo.foobar.com,bar.foobar.com",
		"ttl":         "1h",
		"other_sans":  strings.Join(otherNames, ","),
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
	if len(cert.DNSNames) != 3 ||
		cert.DNSNames[0] != "bar.foobar.com" ||
		cert.DNSNames[1] != "foo.foobar.com" ||
		cert.DNSNames[2] != "foobar.com" {
		t.Fatalf("unexpected DNS SANs %v", cert.DNSNames)
	}
	expectedOtherNames := []otherNameUtf8{{oid1, val1}, {oid2, val2}}
	foundOtherNames, err := getOtherSANsFromX509Extensions(cert.Extensions)
	if err != nil {
		t.Fatal(err)
	}
	if diff := deep.Equal(expectedOtherNames, foundOtherNames); len(diff) != 0 {
		t.Errorf("unexpected otherNames: %v", diff)
	}
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		t.Logf("certificate 3 to check:\n%s", certStr)
	}
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
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		t.Logf("certificate 1 to check:\n%s", certStr)
	}

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
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		t.Logf("certificate 2 to check:\n%s", certStr)
	}
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

func TestBackend_AllowedDomainsTemplate(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
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

	// Write test policy for userpass auth method.
	err := client.Sys().PutPolicy("test", `
   path "pki/*" {  
     capabilities = ["update"]
   }`)
	if err != nil {
		t.Fatal(err)
	}

	// Enable userpass auth method.
	if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
		t.Fatal(err)
	}

	// Configure test role for userpass.
	if _, err := client.Logical().Write("auth/userpass/users/userpassname", map[string]interface{}{
		"password": "test",
		"policies": "test",
	}); err != nil {
		t.Fatal(err)
	}

	// Login userpass for test role and keep client token.
	secret, err := client.Logical().Write("auth/userpass/login/userpassname", map[string]interface{}{
		"password": "test",
	})
	if err != nil || secret == nil {
		t.Fatal(err)
	}
	userpassToken := secret.Auth.ClientToken

	// Get auth accessor for identity template.
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	userpassAccessor := auths["userpass/"].Accessor

	// Mount PKI.
	err = client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "60h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Generate internal CA.
	_, err = client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Write role PKI.
	_, err = client.Logical().Write("pki/roles/test", map[string]interface{}{
		"allowed_domains": []string{"foobar.com", "zipzap.com", "{{identity.entity.aliases." + userpassAccessor + ".name}}",
			"foo.{{identity.entity.aliases." + userpassAccessor + ".name}}.example.com"},
		"allowed_domains_template": true,
		"allow_bare_domains":       true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue certificate with userpassToken.
	client.SetToken(userpassToken)
	_, err = client.Logical().Write("pki/issue/test", map[string]interface{}{"common_name": "userpassname"})
	if err != nil {
		t.Fatal(err)
	}

	// Issue certificate for foobar.com to verify allowed_domain_templae doesnt break plain domains.
	_, err = client.Logical().Write("pki/issue/test", map[string]interface{}{"common_name": "foobar.com"})
	if err != nil {
		t.Fatal(err)
	}

	// Issue certificate for unknown userpassname.
	_, err = client.Logical().Write("pki/issue/test", map[string]interface{}{"common_name": "unknownuserpassname"})
	if err == nil {
		t.Fatal("expected error")
	}

	// Issue certificate for foo.userpassname.domain.
	_, err = client.Logical().Write("pki/issue/test", map[string]interface{}{"common_name": "foo.userpassname.example.com"})
	if err != nil {
		t.Fatal("expected error")
	}

	// Set allowed_domains_template to false.
	_, err = client.Logical().Write("pki/roles/test", map[string]interface{}{
		"allowed_domains_template": false,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue certificate with userpassToken.
	_, err = client.Logical().Write("pki/issue/test", map[string]interface{}{"common_name": "userpassname"})
	if err == nil {
		t.Fatal("expected error")
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

func TestBackend_RevokePlusTidy_Intermediate(t *testing.T) {
	// Enable PKI secret engine
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
	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client

	var err error

	// Mount /pki as a root CA
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

	// Set the cluster's certificate as the root CA in /pki
	pemBundleRootCA := string(cluster.CACertPEM) + string(cluster.CAKeyPEM)
	_, err = client.Logical().Write("pki/config/ca", map[string]interface{}{
		"pem_bundle": pemBundleRootCA,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Mount /pki2 to operate as an intermediate CA
	err = client.Sys().Mount("pki2", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a CSR for the intermediate CA
	secret, err := client.Logical().Write("pki2/intermediate/generate/internal", nil)
	if err != nil {
		t.Fatal(err)
	}
	intermediateCSR := secret.Data["csr"].(string)

	// Sign the intermediate CSR using /pki
	secret, err = client.Logical().Write("pki/root/sign-intermediate", map[string]interface{}{
		"permitted_dns_domains": ".myvault.com",
		"csr":                   intermediateCSR,
	})
	if err != nil {
		t.Fatal(err)
	}
	intermediateCertSerial := secret.Data["serial_number"].(string)
	intermediateCASerialColon := strings.ReplaceAll(strings.ToLower(intermediateCertSerial), ":", "-")

	// Get the intermediate cert after signing
	secret, err = client.Logical().Read("pki/cert/" + intermediateCASerialColon)
	if err != nil {
		t.Fatal(err)
	}
	if secret == nil || len(secret.Data) == 0 || len(secret.Data["certificate"].(string)) == 0 {
		t.Fatal("expected certificate information from read operation")
	}

	// Issue a revoke on on /pki
	_, err = client.Logical().Write("pki/revoke", map[string]interface{}{
		"serial_number": intermediateCertSerial,
	})

	// Revoke adds a fixed 2s buffer, so we sleep for a bit longer to ensure
	// the revocation time is past the current time.
	time.Sleep(3 * time.Second)

	// Issue a tidy on /pki
	_, err = client.Logical().Write("pki/tidy", map[string]interface{}{
		"tidy_cert_store":    true,
		"tidy_revoked_certs": true,
		"safety_buffer":      "1s",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sleep a bit to make sure we're past the safety buffer
	time.Sleep(2 * time.Second)

	// Attempt to read the intermediate cert after revoke + tidy, and ensure
	// that it's no longer present
	secret, err = client.Logical().Read("pki/cert/" + intermediateCASerialColon)
	if err != nil {
		t.Fatal(err)
	}
	if secret != nil {
		t.Fatalf("expected empty response data, got: %#v", secret.Data)
	}

}

var (
	initTest  sync.Once
	rsaCAKey  string
	rsaCACert string
	ecCAKey   string
	ecCACert  string
)
