// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"bytes"
	"cmp"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/url"
	"os"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/fatih/structs"
	"github.com/go-test/deep"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/userpass"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/parsing"
	"github.com/hashicorp/vault/builtin/logical/pki/pki_backend"
	"github.com/hashicorp/vault/helper/testhelpers"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/helper/testhelpers/teststorage"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
	"golang.org/x/net/idna"
)

var stepCount = 0

// From builtin/credential/cert/test-fixtures/root/rootcacert.pem
const (
	rootCACertPEM = `-----BEGIN CERTIFICATE-----
MIIDPDCCAiSgAwIBAgIUb5id+GcaMeMnYBv3MvdTGWigyJ0wDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLZXhhbXBsZS5jb20wHhcNMTYwMjI5MDIyNzI5WhcNMjYw
MjI2MDIyNzU5WjAWMRQwEgYDVQQDEwtleGFtcGxlLmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAOxTMvhTuIRc2YhxZpmPwegP86cgnqfT1mXxi1A7
Q7qax24Nqbf00I3oDMQtAJlj2RB3hvRSCb0/lkF7i1Bub+TGxuM7NtZqp2F8FgG0
z2md+W6adwW26rlxbQKjmRvMn66G9YPTkoJmPmxt2Tccb9+apmwW7lslL5j8H48x
AHJTMb+PMP9kbOHV5Abr3PT4jXUPUr/mWBvBiKiHG0Xd/HEmlyOEPeAThxK+I5tb
6m+eB+7cL9BsvQpy135+2bRAxUphvFi5NhryJ2vlAvoJ8UqigsNK3E28ut60FAoH
SWRfFUFFYtfPgTDS1yOKU/z/XMU2giQv2HrleWt0mp4jqBUCAwEAAaOBgTB/MA4G
A1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSdxLNP/ocx
7HK6JT3/sSAe76iTmzAfBgNVHSMEGDAWgBSdxLNP/ocx7HK6JT3/sSAe76iTmzAc
BgNVHREEFTATggtleGFtcGxlLmNvbYcEfwAAATANBgkqhkiG9w0BAQsFAAOCAQEA
wHThDRsXJunKbAapxmQ6bDxSvTvkLA6m97TXlsFgL+Q3Jrg9HoJCNowJ0pUTwhP2
U946dCnSCkZck0fqkwVi4vJ5EQnkvyEbfN4W5qVsQKOFaFVzep6Qid4rZT6owWPa
cNNzNcXAee3/j6hgr6OQ/i3J6fYR4YouYxYkjojYyg+CMdn6q8BoV0BTsHdnw1/N
ScbnBHQIvIZMBDAmQueQZolgJcdOuBLYHe/kRy167z8nGg+PUFKIYOL8NaOU1+CJ
t2YaEibVq5MRqCbRgnd9a2vG0jr5a3Mn4CUUYv+5qIjP3hUusYenW1/EWtn1s/gk
zehNe5dFTjFpylg1o6b8Ow==
-----END CERTIFICATE-----`
	rootCAKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA7FMy+FO4hFzZiHFmmY/B6A/zpyCep9PWZfGLUDtDuprHbg2p
t/TQjegMxC0AmWPZEHeG9FIJvT+WQXuLUG5v5MbG4zs21mqnYXwWAbTPaZ35bpp3
BbbquXFtAqOZG8yfrob1g9OSgmY+bG3ZNxxv35qmbBbuWyUvmPwfjzEAclMxv48w
/2Rs4dXkBuvc9PiNdQ9Sv+ZYG8GIqIcbRd38cSaXI4Q94BOHEr4jm1vqb54H7twv
0Gy9CnLXfn7ZtEDFSmG8WLk2GvIna+UC+gnxSqKCw0rcTby63rQUCgdJZF8VQUVi
18+BMNLXI4pT/P9cxTaCJC/YeuV5a3SaniOoFQIDAQABAoIBAQCoGZJC84JnnIgb
ttZNWuWKBXbCJcDVDikOQJ9hBZbqsFg1X0CfGmQS3MHf9Ubc1Ro8zVjQh15oIEfn
8lIpdzTeXcpxLdiW8ix3ekVJF20F6pnXY8ZP6UnTeOwamXY6QPZAtb0D9UXcvY+f
nw+IVRD6082XS0Rmzu+peYWVXDy+FDN+HJRANBcdJZz8gOmNBIe0qDWx1b85d/s8
2Kk1Wwdss1IwAGeSddTSwzBNaaHdItZaMZOqPW1gRyBfVSkcUQIE6zn2RKw2b70t
grkIvyRcTdfmiKbqkkJ+eR+ITOUt0cBZSH4cDjlQA+r7hulvoBpQBRj068Toxkcc
bTagHaPBAoGBAPWPGVkHqhTbJ/DjmqDIStxby2M1fhhHt4xUGHinhUYjQjGOtDQ9
0mfaB7HObudRiSLydRAVGAHGyNJdQcTeFxeQbovwGiYKfZSA1IGpea7dTxPpGEdN
ksA0pzSp9MfKzX/MdLuAkEtO58aAg5YzsgX9hDNxo4MhH/gremZhEGZlAoGBAPZf
lqdYvAL0fjHGJ1FUEalhzGCGE9PH2iOqsxqLCXK7bDbzYSjvuiHkhYJHAOgVdiW1
lB34UHHYAqZ1VVoFqJ05gax6DE2+r7K5VV3FUCaC0Zm3pavxchU9R/TKP82xRrBj
AFWwdgDTxUyvQEmgPR9sqorftO71Iz2tiwyTpIfxAoGBAIhEMLzHFAse0rtKkrRG
ccR27BbRyHeQ1Lp6sFnEHKEfT8xQdI/I/snCpCJ3e/PBu2g5Q9z416mktiyGs8ib
thTNgYsGYnxZtfaCx2pssanoBcn2wBJRae5fSapf5gY49HDG9MBYR7qCvvvYtSzU
4yWP2ZzyotpRt3vwJKxLkN5BAoGAORHpZvhiDNkvxj3da7Rqpu7VleJZA2y+9hYb
iOF+HcqWhaAY+I+XcTRrTMM/zYLzLEcEeXDEyao86uwxCjpXVZw1kotvAC9UqbTO
tnr3VwRkoxPsV4kFYTAh0+1pnC8dbcxxDmhi3Uww3tOVs7hfkEDuvF6XnebA9A+Y
LyCgMzECgYEA6cCU8QODOivIKWFRXucvWckgE6MYDBaAwe6qcLsd1Q/gpE2e3yQc
4RB3bcyiPROLzMLlXFxf1vSNJQdIaVfrRv+zJeGIiivLPU8+Eq4Lrb+tl1LepcOX
OzQeADTSCn5VidOfjDkIst9UXjMlrFfV9/oJEw5Eiqa6lkNPCGDhfA8=
-----END RSA PRIVATE KEY-----`
)

func TestPKI_RequireCN(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected ca info")
	}

	// Create a role which does require CN (default)
	_, err = CBWrite(b, s, "roles/example", map[string]interface{}{
		"allowed_domains":    "foobar.com,zipzap.com,abc.com,xyz.com",
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"max_ttl":            "2h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue a cert with require_cn set to true and with common name supplied.
	// It should succeed.
	resp, err = CBWrite(b, s, "issue/example", map[string]interface{}{
		"common_name": "foobar.com",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("issue/example"), logical.UpdateOperation), resp, true)
	if err != nil {
		t.Fatal(err)
	}

	// Issue a cert with require_cn set to true and with out supplying the
	// common name. It should error out.
	_, err = CBWrite(b, s, "issue/example", map[string]interface{}{})
	if err == nil {
		t.Fatalf("expected an error due to missing common_name")
	}

	// Modify the role to make the common name optional
	_, err = CBWrite(b, s, "roles/example", map[string]interface{}{
		"allowed_domains":    "foobar.com,zipzap.com,abc.com,xyz.com",
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"max_ttl":            "2h",
		"require_cn":         false,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue a cert with require_cn set to false and without supplying the
	// common name. It should succeed.
	resp, err = CBWrite(b, s, "issue/example", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	if resp.Data["certificate"] == "" {
		t.Fatalf("expected a cert to be generated")
	}

	// Issue a cert with require_cn set to false and with a common name. It
	// should succeed.
	resp, err = CBWrite(b, s, "issue/example", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	if resp.Data["certificate"] == "" {
		t.Fatalf("expected a cert to be generated")
	}
}

func TestPKI_DeviceCert(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name":         "myvault.com",
		"not_after":           "9999-12-31T23:59:59Z",
		"not_before_duration": "2h",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected ca info")
	}
	var certBundle certutil.CertBundle
	err = mapstructure.Decode(resp.Data, &certBundle)
	if err != nil {
		t.Fatal(err)
	}

	parsedCertBundle, err := certBundle.ToParsedCertBundle()
	if err != nil {
		t.Fatal(err)
	}
	cert := parsedCertBundle.Certificate
	notAfter := cert.NotAfter.Format(time.RFC3339)
	if notAfter != "9999-12-31T23:59:59Z" {
		t.Fatalf("not after from certificate: %v is not matching with input parameter: %v", cert.NotAfter, "9999-12-31T23:59:59Z")
	}
	if math.Abs(float64(time.Now().Add(-2*time.Hour).Unix()-cert.NotBefore.Unix())) > 10 {
		t.Fatalf("root/generate/internal did not properly set validity period (notBefore): was %v vs expected %v", cert.NotBefore, time.Now().Add(-2*time.Hour))
	}

	// Create a role which does require CN (default)
	_, err = CBWrite(b, s, "roles/example", map[string]interface{}{
		"allowed_domains":    "foobar.com,zipzap.com,abc.com,xyz.com",
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"not_after":          "9999-12-31T23:59:59Z",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue a cert with require_cn set to true and with common name supplied.
	// It should succeed.
	resp, err = CBWrite(b, s, "issue/example", map[string]interface{}{
		"common_name": "foobar.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = mapstructure.Decode(resp.Data, &certBundle)
	if err != nil {
		t.Fatal(err)
	}

	parsedCertBundle, err = certBundle.ToParsedCertBundle()
	if err != nil {
		t.Fatal(err)
	}
	cert = parsedCertBundle.Certificate
	notAfter = cert.NotAfter.Format(time.RFC3339)
	if notAfter != "9999-12-31T23:59:59Z" {
		t.Fatal(fmt.Errorf("not after from certificate  is not matching with input parameter"))
	}
}

func TestBackend_InvalidParameter(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	_, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "myvault.com",
		"not_after":   "9999-12-31T23:59:59Z",
		"ttl":         "25h",
	})
	if err == nil {
		t.Fatal(err)
	}

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "myvault.com",
		"not_after":   "9999-12-31T23:59:59",
	})
	if err == nil {
		t.Fatal(err)
	}
}

func TestBackend_CSRValues(t *testing.T) {
	t.Parallel()
	initTest.Do(setCerts)
	b, _ := CreateBackendWithStorage(t)

	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps:          []logicaltest.TestStep{},
	}

	intdata := map[string]interface{}{}
	reqdata := map[string]interface{}{}
	testCase.Steps = append(testCase.Steps, generateCSRSteps(t, ecCACert, ecCAKey, intdata, reqdata)...)

	logicaltest.Test(t, testCase)
}

func TestBackend_SerialNumberSource(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	var err error

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = CBWrite(b, s, "roles/json-csr", map[string]interface{}{
		"allow_any_name":         true,
		"enforce_hostnames":      false,
		"allowed_serial_numbers": "foo*",
		"serial_number_source":   "json-csr",
		"key_type":               "any",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = CBWrite(b, s, "roles/json", map[string]interface{}{
		"allow_any_name":         true,
		"enforce_hostnames":      false,
		"allowed_serial_numbers": "foo*",
		"serial_number_source":   "json",
		"key_type":               "any",
	})

	// Create a CSR with a serial number not allowed by the role.
	tmpl := &x509.CertificateRequest{
		Subject: pkix.Name{SerialNumber: "bar"},
	}
	_, _, csrPem := generateCSR(t, tmpl, "ec", 256)

	// Signing a csr with a disallowed subject serial number in the CSR
	// with serial_number_source=json-csr should fail.
	_, err = CBWrite(b, s, "sign/json-csr", map[string]interface{}{
		"common_name": "localhost",
		"csr":         csrPem,
	})
	if err == nil {
		t.Fatal("expected an error")
	}

	// The serial number in the request should take precedence.
	_, err = CBWrite(b, s, "sign/json-csr", map[string]interface{}{
		"common_name":   "localhost",
		"csr":           csrPem,
		"serial_number": "foobar",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Try signing the cert with serial_number_source=json.
	// The serial in the CSR should be ignored.
	_, err = CBWrite(b, s, "sign/json", map[string]interface{}{
		"common_name": "localhost",
		"csr":         csrPem,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Try signing the cert with serial_number_source=json
	// and a serial number in the request
	_, err = CBWrite(b, s, "sign/json", map[string]interface{}{
		"common_name":   "localhost",
		"csr":           csrPem,
		"serial_number": "foobar2",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_URLsCRUD(t *testing.T) {
	t.Parallel()
	initTest.Do(setCerts)
	b, _ := CreateBackendWithStorage(t)

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
	t.Parallel()
	cases := []struct {
		name      string
		key, cert *string
		useCSR    bool
	}{
		{"RSA", &rsaCAKey, &rsaCACert, false},
		{"RSACSR", &rsaCAKey, &rsaCACert, true},
		{"EC", &ecCAKey, &ecCACert, false},
		{"ECCSR", &ecCAKey, &ecCACert, true},
		{"ED", &edCAKey, &edCACert, false},
		{"EDCSR", &edCAKey, &edCACert, true},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			initTest.Do(setCerts)
			b, _ := CreateBackendWithStorage(t)

			testCase := logicaltest.TestCase{
				LogicalBackend: b,
				Steps: []logicaltest.TestStep{
					{
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
		case "ed25519":
			parsedCertBundle.PrivateKeyType = certutil.Ed25519PrivateKey
			parsedCertBundle.PrivateKey = key
			parsedCertBundle.PrivateKeyBytes, err = x509.MarshalPKCS8PrivateKey(key.(ed25519.PrivateKey))
			if err != nil {
				return nil, fmt.Errorf("error parsing Ed25519 key: %s", err)
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
	case parsedCertBundle.PrivateKeyType == certutil.Ed25519PrivateKey && keyType != "ed25519":
		fallthrough
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

	// TODO: We incremented 20->25 due to CircleCI execution
	// being slow and pausing this test. We might consider recording the
	// actual issuance time of the cert and calculating the expected
	// validity period +/- fuzz, but that'd require recording and passing
	// through more information.
	if math.Abs(float64(time.Now().Add(validity).Unix()-cert.NotAfter.Unix())) > 25 {
		return nil, fmt.Errorf("certificate validity end: %s; expected within 25 seconds of %s", cert.NotAfter.Format(time.RFC3339), time.Now().Add(validity).Format(time.RFC3339))
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

	priv1024, _ := cryptoutil.GenerateRSAKey(rand.Reader, 1024)
	csr1024, _ := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, priv1024)
	csrPem1024 := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr1024,
	})))

	priv2048, _ := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	csr2048, _ := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, priv2048)
	csrPem2048 := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr2048,
	})))

	ret := []logicaltest.TestStep{
		{
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

		{
			Operation: logical.UpdateOperation,
			Path:      "config/urls",
			Data: map[string]interface{}{
				"issuing_certificates":    strings.Join(expected.IssuingCertificates, ","),
				"crl_distribution_points": strings.Join(expected.CRLDistributionPoints, ","),
				"ocsp_servers":            strings.Join(expected.OCSPServers, ","),
			},
		},

		{
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

		{
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

		{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data: map[string]interface{}{
				"common_name":         "intermediate.cert.com",
				"csr":                 csrPem2048,
				"signature_bits":      512,
				"format":              "der",
				"not_before_duration": "2h",
				// Let's Encrypt -- R3 SKID
				"skid": "14:2E:B3:17:B7:58:56:CB:AE:50:09:40:E6:1F:AF:9D:8B:14:C2:C6",
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
					return fmt.Errorf("returned cert cannot be parsed: %w", err)
				}
				if len(certs) != 1 {
					return fmt.Errorf("unexpected returned length of certificates: %d", len(certs))
				}
				cert := certs[0]

				skid, _ := hex.DecodeString("142EB317B75856CBAE500940E61FAF9D8B14C2C6")

				switch {
				case !reflect.DeepEqual(expected.IssuingCertificates, cert.IssuingCertificateURL):
					return fmt.Errorf("IssuingCertificateURL:\nexpected\n%#v\ngot\n%#v\n", expected.IssuingCertificates, cert.IssuingCertificateURL)
				case !reflect.DeepEqual(expected.CRLDistributionPoints, cert.CRLDistributionPoints):
					return fmt.Errorf("CRLDistributionPoints:\nexpected\n%#v\ngot\n%#v\n", expected.CRLDistributionPoints, cert.CRLDistributionPoints)
				case !reflect.DeepEqual(expected.OCSPServers, cert.OCSPServer):
					return fmt.Errorf("OCSPServer:\nexpected\n%#v\ngot\n%#v\n", expected.OCSPServers, cert.OCSPServer)
				case !reflect.DeepEqual([]string{"intermediate.cert.com"}, cert.DNSNames):
					return fmt.Errorf("DNSNames\nexpected\n%#v\ngot\n%#v\n", []string{"intermediate.cert.com"}, cert.DNSNames)
				case !reflect.DeepEqual(x509.SHA512WithRSA, cert.SignatureAlgorithm):
					return fmt.Errorf("Signature Algorithm:\nexpected\n%#v\ngot\n%#v\n", x509.SHA512WithRSA, cert.SignatureAlgorithm)
				case !reflect.DeepEqual(skid, cert.SubjectKeyId):
					return fmt.Errorf("SKID:\nexpected\n%#v\ngot\n%#v\n", skid, cert.SubjectKeyId)
				}

				if math.Abs(float64(time.Now().Add(-2*time.Hour).Unix()-cert.NotBefore.Unix())) > 10 {
					t.Fatalf("root/sign-intermediate did not properly set validity period (notBefore): was %v vs expected %v", cert.NotBefore, time.Now().Add(-2*time.Hour))
				}

				return nil
			},
		},

		// Same as above but exclude adding to sans
		{
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
					return fmt.Errorf("returned cert cannot be parsed: %w", err)
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

func generateCSR(t *testing.T, csrTemplate *x509.CertificateRequest, keyType string, keyBits int) (interface{}, []byte, string) {
	t.Helper()

	var priv interface{}
	var err error
	switch keyType {
	case "rsa":
		priv, err = cryptoutil.GenerateRSAKey(rand.Reader, keyBits)
	case "ec":
		switch keyBits {
		case 224:
			priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
		case 256:
			priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		case 384:
			priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		case 521:
			priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		default:
			t.Fatalf("Got unknown ec< key bits: %v", keyBits)
		}
	case "ed25519":
		_, priv, err = ed25519.GenerateKey(rand.Reader)
	}

	if err != nil {
		t.Fatalf("Got error generating private key for CSR: %v", err)
	}

	return generateCSRWithKey(t, csrTemplate, priv)
}

func generateCSRWithKey(t *testing.T, csrTemplate *x509.CertificateRequest, priv interface{}) (interface{}, []byte, string) {
	csr, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, priv)
	if err != nil {
		t.Fatalf("Got error generating CSR: %v", err)
	}

	csrPem := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	})))

	return priv, csr, csrPem
}

func generateCSRSteps(t *testing.T, caCert, caKey string, intdata, reqdata map[string]interface{}) []logicaltest.TestStep {
	csrTemplate, csrPem := generateTestCsr(t, certutil.RSAPrivateKey, 2048)

	ret := []logicaltest.TestStep{
		{
			Operation: logical.UpdateOperation,
			Path:      "root/generate/exported",
			Data: map[string]interface{}{
				"common_name":     "Root Cert",
				"ttl":             "180h",
				"max_path_length": 0,
			},
		},

		{
			Operation: logical.UpdateOperation,
			Path:      "root/sign-intermediate",
			Data: map[string]interface{}{
				"use_csr_values": true,
				"csr":            csrPem,
				"format":         "der",
			},
			ErrorOk: true,
		},

		{
			Operation: logical.DeleteOperation,
			Path:      "root",
		},

		{
			Operation: logical.UpdateOperation,
			Path:      "root/generate/exported",
			Data: map[string]interface{}{
				"common_name":     "Root Cert",
				"ttl":             "180h",
				"max_path_length": 1,
			},
		},

		{
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
					return fmt.Errorf("returned cert cannot be parsed: %w", err)
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

func generateTestCsr(t *testing.T, keyType certutil.PrivateKeyType, keyBits int) (x509.CertificateRequest, string) {
	t.Helper()

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

	_, _, csrPem := generateCSR(t, &csrTemplate, string(keyType), keyBits)
	return csrTemplate, csrPem
}

// Generates steps to test out various role permutations
func generateRoleSteps(t *testing.T, useCSRs bool) []logicaltest.TestStep {
	roleVals := issuing.RoleEntry{
		MaxTTL:                    12 * time.Hour,
		KeyType:                   "rsa",
		KeyBits:                   2048,
		RequireCN:                 true,
		AllowWildcardCertificates: new(bool),
	}
	*roleVals.AllowWildcardCertificates = true

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
	generatedEdKeys := map[int]crypto.Signer{}
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
		// t.Logf("test step %d\nrole vals: %#v\n", stepCount, roleVals)
		stepCount++
		// t.Logf("test step %d\nissue vals: %#v\n", stepCount, issueTestStep)
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

	getCountryCheck := func(role issuing.RoleEntry) logicaltest.TestCheckFunc {
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

	getOuCheck := func(role issuing.RoleEntry) logicaltest.TestCheckFunc {
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

	getOrganizationCheck := func(role issuing.RoleEntry) logicaltest.TestCheckFunc {
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

	getLocalityCheck := func(role issuing.RoleEntry) logicaltest.TestCheckFunc {
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

	getProvinceCheck := func(role issuing.RoleEntry) logicaltest.TestCheckFunc {
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

	getStreetAddressCheck := func(role issuing.RoleEntry) logicaltest.TestCheckFunc {
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

	getPostalCodeCheck := func(role issuing.RoleEntry) logicaltest.TestCheckFunc {
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

	getNotBeforeCheck := func(role issuing.RoleEntry) logicaltest.TestCheckFunc {
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

			actualDiff := time.Since(cert.NotBefore)
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
	getCnCheck := func(name string, role issuing.RoleEntry, key crypto.Signer, usage x509.KeyUsage,
		extUsage x509.ExtKeyUsage, validity time.Duration,
	) logicaltest.TestCheckFunc {
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
				privKey, _ = cryptoutil.GenerateRSAKey(rand.Reader, keyBits)
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

		case "ed25519":
			privKey, ok = generatedEdKeys[keyBits]
			if !ok {
				_, privKey, _ = ed25519.GenerateKey(rand.Reader)
				generatedEdKeys[keyBits] = privKey
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
		rsaKeyBits := []int{2048, 3072, 4096, 8192}
		ecKeyBits := []int{224, 256, 384, 521}
		plan := csrPlan{errorOk: errorOk}

		var testBitSize int
		switch keyType {
		case "rsa":
			plan.roleKeyBits = rsaKeyBits[mathRand.Int()%len(rsaKeyBits)]
			testBitSize = plan.roleKeyBits

			// If we don't expect an error already, randomly choose a
			// key size and expect an error if it's less than the role
			// setting
			if !keybitSizeRandOff && !errorOk {
				testBitSize = rsaKeyBits[mathRand.Int()%len(rsaKeyBits)]
			}

			if testBitSize < plan.roleKeyBits {
				plan.errorOk = true
			}

		case "ec":
			plan.roleKeyBits = ecKeyBits[mathRand.Int()%len(ecKeyBits)]
			testBitSize = plan.roleKeyBits

			// If we don't expect an error already, randomly choose a
			// key size and expect an error if it's less than the role
			// setting
			if !keybitSizeRandOff && !errorOk {
				testBitSize = ecKeyBits[mathRand.Int()%len(ecKeyBits)]
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
			if mathRand.Int()%3 == 1 {
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
			parsedKeyUsage := parsing.ParseKeyUsages(roleVals.KeyUsage)
			if parsedKeyUsage == 0 && len(usage) != 0 {
				panic("parsed key usages was zero")
			}

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

	funcs := []interface{}{
		addCnTests, getCnCheck, getCountryCheck, getLocalityCheck, getNotBeforeCheck,
		getOrganizationCheck, getOuCheck, getPostalCodeCheck, getRandCsr, getStreetAddressCheck,
		getProvinceCheck,
	}
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

		roleVals.OU = []string{"bar", "foo"}
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
				var expected []net.IP
				expected = append(expected, expectedIp...)
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
		getOtherCheck := func(expectedOthers ...certutil.OtherNameUtf8) logicaltest.TestCheckFunc {
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
				if err != nil {
					return err
				}
				var expected []certutil.OtherNameUtf8
				expected = append(expected, expectedOthers...)
				if diff := deep.Equal(foundOthers, expected); len(diff) > 0 {
					return fmt.Errorf("wrong SAN IPs, diff: %v", diff)
				}
				return nil
			}
		}

		addOtherSANTests := func(useCSRs, useCSRSANs bool, allowedOtherSANs []string, errorOk bool, otherSANs []string, csrOtherSANs []certutil.OtherNameUtf8, check logicaltest.TestCheckFunc) {
			otherSansMap := func(os []certutil.OtherNameUtf8) map[string][]string {
				ret := make(map[string][]string)
				for _, o := range os {
					ret[o.Oid] = append(ret[o.Oid], o.Value)
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

		newOtherNameUtf8 := func(s string) (ret certutil.OtherNameUtf8) {
			pieces := strings.Split(s, ";")
			if len(pieces) == 2 {
				piecesRest := strings.Split(pieces[1], ":")
				if len(piecesRest) == 2 {
					switch strings.ToUpper(piecesRest[0]) {
					case "UTF-8", "UTF8":
						return certutil.OtherNameUtf8{Oid: pieces[0], Value: piecesRest[1]}
					}
				}
			}
			t.Fatalf("error parsing otherName: %q", s)
			return
		}
		oid1 := "1.3.6.1.4.1.311.20.2.3"
		oth1str := oid1 + ";utf8:devops@nope.com"
		oth1 := newOtherNameUtf8(oth1str)
		oth2 := certutil.OtherNameUtf8{oid1, "me@example.com"}
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
			[]certutil.OtherNameUtf8{oth2}, getOtherCheck(oth1))

		if useCSRs {
			// OtherSANs not allowed, valid OtherSANs provided via CSR, should be an error.
			addOtherSANTests(useCSRs, true, allowNone, true, nil, []certutil.OtherNameUtf8{oth1}, nil)

			// Given OtherSANs as both API and CSR arguments and useCSRSANs=true, API arg ignored.
			addOtherSANTests(useCSRs, false, allowAll, false, []string{oth2.String()},
				[]certutil.OtherNameUtf8{oth1}, getOtherCheck(oth2))
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

func TestRolesAltIssuer(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Create two issuers.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root a - example.com",
		"issuer_name": "root-a",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	rootAPem := resp.Data["certificate"].(string)
	rootACert := parseCert(t, rootAPem)

	resp, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root b - example.com",
		"issuer_name": "root-b",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	rootBPem := resp.Data["certificate"].(string)
	rootBCert := parseCert(t, rootBPem)

	// Create three roles: one with no assignment, one with explicit root-a,
	// one with explicit root-b.
	_, err = CBWrite(b, s, "roles/use-default", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
	})
	require.NoError(t, err)

	_, err = CBWrite(b, s, "roles/use-root-a", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
		"issuer_ref":        "root-a",
	})
	require.NoError(t, err)

	_, err = CBWrite(b, s, "roles/use-root-b", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"issuer_ref":        "root-b",
	})
	require.NoError(t, err)

	// Now issue certs against these roles.
	resp, err = CBWrite(b, s, "issue/use-default", map[string]interface{}{
		"common_name": "testing",
		"ttl":         "5s",
	})
	require.NoError(t, err)
	leafPem := resp.Data["certificate"].(string)
	leafCert := parseCert(t, leafPem)
	err = leafCert.CheckSignatureFrom(rootACert)
	require.NoError(t, err, "should be signed by root-a but wasn't")

	resp, err = CBWrite(b, s, "issue/use-root-a", map[string]interface{}{
		"common_name": "testing",
		"ttl":         "5s",
	})
	require.NoError(t, err)
	leafPem = resp.Data["certificate"].(string)
	leafCert = parseCert(t, leafPem)
	err = leafCert.CheckSignatureFrom(rootACert)
	require.NoError(t, err, "should be signed by root-a but wasn't")

	resp, err = CBWrite(b, s, "issue/use-root-b", map[string]interface{}{
		"common_name": "testing",
		"ttl":         "5s",
	})
	require.NoError(t, err)
	leafPem = resp.Data["certificate"].(string)
	leafCert = parseCert(t, leafPem)
	err = leafCert.CheckSignatureFrom(rootBCert)
	require.NoError(t, err, "should be signed by root-b but wasn't")

	// Update the default issuer to be root B and make sure that the
	// use-default role updates.
	_, err = CBWrite(b, s, "config/issuers", map[string]interface{}{
		"default": "root-b",
	})
	require.NoError(t, err)

	resp, err = CBWrite(b, s, "issue/use-default", map[string]interface{}{
		"common_name": "testing",
		"ttl":         "5s",
	})
	require.NoError(t, err)
	leafPem = resp.Data["certificate"].(string)
	leafCert = parseCert(t, leafPem)
	err = leafCert.CheckSignatureFrom(rootBCert)
	require.NoError(t, err, "should be signed by root-b but wasn't")
}

func TestBackend_PathFetchValidRaw(t *testing.T) {
	t.Parallel()
	b, storage := CreateBackendWithStorage(t)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/generate/internal",
		Storage:   storage,
		Data: map[string]interface{}{
			"common_name": "test.com",
			"ttl":         "6h",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to generate root, %#v", resp)
	}
	rootCaAsPem := resp.Data["certificate"].(string)

	// Chain should contain the root.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "ca_chain",
		Storage:    storage,
		Data:       map[string]interface{}{},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	if resp != nil && resp.IsError() {
		t.Fatalf("failed read ca_chain, %#v", resp)
	}
	if strings.Count(string(resp.Data[logical.HTTPRawBody].([]byte)), rootCaAsPem) != 1 {
		t.Fatalf("expected raw chain to contain the root cert")
	}

	// The ca/pem should return us the actual CA...
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "ca/pem",
		Storage:    storage,
		Data:       map[string]interface{}{},
		MountPoint: "pki/",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("ca/pem"), logical.ReadOperation), resp, true)
	require.NoError(t, err)
	if resp != nil && resp.IsError() {
		t.Fatalf("failed read ca/pem, %#v", resp)
	}
	// check the raw cert matches the response body
	if !bytes.Equal(resp.Data[logical.HTTPRawBody].([]byte), []byte(rootCaAsPem)) {
		t.Fatalf("failed to get raw cert")
	}

	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/example",
		Storage:   storage,
		Data: map[string]interface{}{
			"allowed_domains":  "example.com",
			"allow_subdomains": "true",
			"max_ttl":          "1h",
			"no_store":         "false",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error setting up pki role: %v", err)

	// Now issue a short-lived certificate from our pki-external.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issue/example",
		Storage:   storage,
		Data: map[string]interface{}{
			"common_name": "test.example.com",
			"ttl":         "5m",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "error issuing certificate: %v", err)
	require.NotNil(t, resp, "got nil response from issuing request")

	issueCrtAsPem := resp.Data["certificate"].(string)
	issuedCrt := parseCert(t, issueCrtAsPem)
	expectedSerial := serialFromCert(issuedCrt)
	expectedCert := []byte(issueCrtAsPem)

	// get der cert
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      fmt.Sprintf("cert/%s/raw", expectedSerial),
		Storage:   storage,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to get raw cert, %#v", resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	// check the raw cert matches the response body
	rawBody := resp.Data[logical.HTTPRawBody].([]byte)
	bodyAsPem := []byte(strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rawBody}))))
	if !bytes.Equal(bodyAsPem, expectedCert) {
		t.Fatalf("failed to get raw cert for serial number: %s", expectedSerial)
	}
	if resp.Data[logical.HTTPContentType] != "application/pkix-cert" {
		t.Fatalf("failed to get raw cert content-type")
	}

	// get pem
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      fmt.Sprintf("cert/%s/raw/pem", expectedSerial),
		Storage:   storage,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to get raw, %#v", resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	// check the pem cert matches the response body
	if !bytes.Equal(resp.Data[logical.HTTPRawBody].([]byte), expectedCert) {
		t.Fatalf("failed to get pem cert")
	}
	if resp.Data[logical.HTTPContentType] != "application/pem-certificate-chain" {
		t.Fatalf("failed to get raw cert content-type")
	}
}

func TestBackend_PathFetchCertList(t *testing.T) {
	t.Parallel()
	// create the backend
	b, storage := CreateBackendWithStorage(t)

	// generate root
	rootData := map[string]interface{}{
		"common_name": "test.com",
		"ttl":         "6h",
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "root/generate/internal",
		Storage:    storage,
		Data:       rootData,
		MountPoint: "pki/",
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
		Operation:  logical.UpdateOperation,
		Path:       "config/urls",
		Storage:    storage,
		Data:       urlsData,
		MountPoint: "pki/",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("config/urls"), logical.UpdateOperation), resp, true)

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "config/urls",
		Storage:    storage,
		MountPoint: "pki/",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("config/urls"), logical.ReadOperation), resp, true)

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
		Operation:  logical.UpdateOperation,
		Path:       "roles/test-example",
		Storage:    storage,
		Data:       roleData,
		MountPoint: "pki/",
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
			Operation:  logical.UpdateOperation,
			Path:       "issue/test-example",
			Storage:    storage,
			Data:       certData,
			MountPoint: "pki/",
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
		Operation:  logical.ListOperation,
		Path:       "certs",
		Storage:    storage,
		MountPoint: "pki/",
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
		Operation:  logical.ListOperation,
		Path:       issuing.PathCerts,
		Storage:    storage,
		MountPoint: "pki/",
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
	t.Parallel()
	testCases := []struct {
		testName string
		keyType  string
	}{
		{testName: "RSA", keyType: "rsa"},
		{testName: "ED25519", keyType: "ed25519"},
		{testName: "EC", keyType: "ec"},
		{testName: "Any", keyType: "any"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			runTestSignVerbatim(t, tc.keyType)
		})
	}
}

func runTestSignVerbatim(t *testing.T, keyType string) {
	// create the backend
	b, storage := CreateBackendWithStorage(t)

	// generate root
	rootData := map[string]interface{}{
		"common_name": "test.com",
		"not_after":   "9999-12-31T23:59:59Z",
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "root/generate/internal",
		Storage:    storage,
		Data:       rootData,
		MountPoint: "pki/",
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to generate root, %#v", *resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	// create a CSR and key
	key, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	csrReq := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: "foo.bar.com",
		},
		// Check that otherName extensions are not duplicated (see hashicorp/vault#16700).
		// If these extensions are duplicated, sign-verbatim will fail when parsing the signed certificate on Go 1.19+ (see golang/go#50988).
		// On older versions of Go this test will fail due to an explicit check for duplicate otherNames later in this test.
		ExtraExtensions: []pkix.Extension{
			{
				Id:       certutil.OidExtensionSubjectAltName,
				Critical: false,
				Value:    []byte{0x30, 0x26, 0xA0, 0x24, 0x06, 0x0A, 0x2B, 0x06, 0x01, 0x04, 0x01, 0x82, 0x37, 0x14, 0x02, 0x03, 0xA0, 0x16, 0x0C, 0x14, 0x75, 0x73, 0x65, 0x72, 0x6E, 0x61, 0x6D, 0x65, 0x40, 0x65, 0x78, 0x61, 0x6D, 0x70, 0x6C, 0x65, 0x2E, 0x63, 0x6F, 0x6D},
			},
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

	signVerbatimData := map[string]interface{}{
		"csr": pemCSR,
	}
	if keyType == "rsa" {
		signVerbatimData["signature_bits"] = 512
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "sign-verbatim",
		Storage:    storage,
		Data:       signVerbatimData,
		MountPoint: "pki/",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("sign-verbatim"), logical.UpdateOperation), resp, true)

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
		"ttl":                 "4h",
		"max_ttl":             "8h",
		"key_type":            keyType,
		"not_before_duration": "2h",
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "roles/test",
		Storage:    storage,
		Data:       roleData,
		MountPoint: "pki/",
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
		MountPoint: "pki/",
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
		MountPoint: "pki/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal(resp.Error().Error())
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
		t.Fatalf("sign-verbatim did not properly cap validity period (notAfter) on signed CSR: was %v vs requested %v but should've been %v", cert.NotAfter, time.Now().Add(12*time.Hour), time.Now().Add(8*time.Hour))
	}
	if math.Abs(float64(time.Now().Add(-2*time.Hour).Unix()-cert.NotBefore.Unix())) > 10 {
		t.Fatalf("sign-verbatim did not properly cap validity period (notBefore) on signed CSR: was %v vs expected %v", cert.NotBefore, time.Now().Add(-2*time.Hour))
	}

	// Now check signing a certificate using the not_after input using the Y10K value
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sign-verbatim/test",
		Storage:   storage,
		Data: map[string]interface{}{
			"csr":       pemCSR,
			"not_after": "9999-12-31T23:59:59Z",
		},
		MountPoint: "pki/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal(resp.Error().Error())
	}
	if resp.Data == nil || resp.Data["certificate"] == nil {
		t.Fatal("did not get expected data")
	}
	certString = resp.Data["certificate"].(string)
	block, _ = pem.Decode([]byte(certString))
	if block == nil {
		t.Fatal("nil pem block")
	}
	certs, err = x509.ParseCertificates(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(certs) != 1 {
		t.Fatalf("expected a single cert, got %d", len(certs))
	}
	cert = certs[0]

	// Fallback check for duplicate otherName, necessary on Go versions before 1.19.
	// We assume that there is only one SAN in the original CSR and that it is an otherName.
	san_count := 0
	for _, ext := range cert.Extensions {
		if ext.Id.Equal(certutil.OidExtensionSubjectAltName) {
			san_count += 1
		}
	}
	if san_count != 1 {
		t.Fatalf("expected one SAN extension, got %d", san_count)
	}

	notAfter := cert.NotAfter.Format(time.RFC3339)
	if notAfter != "9999-12-31T23:59:59Z" {
		t.Fatal(fmt.Errorf("not after from certificate is not matching with input parameter"))
	}

	// now check that if we set generate-lease it takes it from the role and the TTLs match
	roleData = map[string]interface{}{
		"ttl":            "4h",
		"max_ttl":        "8h",
		"generate_lease": true,
		"key_type":       keyType,
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "roles/test",
		Storage:    storage,
		Data:       roleData,
		MountPoint: "pki/",
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
		MountPoint: "pki/",
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
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// This is a change within 1.11, we are no longer idempotent across generate/internal calls.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "myvault.com",
	})
	require.NoError(t, err)
	require.NotNil(t, resp, "expected ca info")
	keyId1 := resp.Data["key_id"]
	issuerId1 := resp.Data["issuer_id"]
	cert := parseCert(t, resp.Data["certificate"].(string))
	certSkid := certutil.GetHexFormatted(cert.SubjectKeyId, ":")

	//  -> Validate the SKID matches between the root cert and the key
	resp, err = CBRead(b, s, "key/"+keyId1.(issuing.KeyID).String())
	require.NoError(t, err)
	require.NotNil(t, resp, "expected a response")
	require.Equal(t, resp.Data["subject_key_id"], certSkid)

	resp, err = CBRead(b, s, "cert/ca_chain")
	require.NoError(t, err, "error reading ca_chain: %v", err)

	r1Data := resp.Data

	// Calling generate/internal should generate a new CA as well.
	resp, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "myvault.com",
	})
	require.NoError(t, err)
	require.NotNil(t, resp, "expected ca info")
	keyId2 := resp.Data["key_id"]
	issuerId2 := resp.Data["issuer_id"]
	cert = parseCert(t, resp.Data["certificate"].(string))
	certSkid = certutil.GetHexFormatted(cert.SubjectKeyId, ":")

	//  -> Validate the SKID matches between the root cert and the key
	resp, err = CBRead(b, s, "key/"+keyId2.(issuing.KeyID).String())
	require.NoError(t, err)
	require.NotNil(t, resp, "expected a response")
	require.Equal(t, resp.Data["subject_key_id"], certSkid)

	// Make sure that we actually generated different issuer and key values
	require.NotEqual(t, keyId1, keyId2)
	require.NotEqual(t, issuerId1, issuerId2)

	// Now because the issued CA's have no links, the call to ca_chain should return the same data (ca chain from default)
	resp, err = CBRead(b, s, "cert/ca_chain")
	require.NoError(t, err, "error reading ca_chain: %v", err)
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("cert/ca_chain"), logical.ReadOperation), resp, true)

	r2Data := resp.Data
	if !reflect.DeepEqual(r1Data, r2Data) {
		t.Fatal("got different ca certs")
	}

	// Now let's validate that the import bundle is idempotent.
	pemBundleRootCA := rootCACertPEM + "\n" + rootCAKeyPEM
	resp, err = CBWrite(b, s, "config/ca", map[string]interface{}{
		"pem_bundle": pemBundleRootCA,
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("config/ca"), logical.UpdateOperation), resp, true)

	require.NoError(t, err)
	require.NotNil(t, resp, "expected ca info")
	firstMapping := resp.Data["mapping"].(map[string]string)
	firstImportedKeys := resp.Data["imported_keys"].([]string)
	firstImportedIssuers := resp.Data["imported_issuers"].([]string)
	firstExistingKeys := resp.Data["existing_keys"].([]string)
	firstExistingIssuers := resp.Data["existing_issuers"].([]string)

	require.NotContains(t, firstImportedKeys, keyId1)
	require.NotContains(t, firstImportedKeys, keyId2)
	require.NotContains(t, firstImportedIssuers, issuerId1)
	require.NotContains(t, firstImportedIssuers, issuerId2)
	require.Empty(t, firstExistingKeys)
	require.Empty(t, firstExistingIssuers)
	require.NotEmpty(t, firstMapping)
	require.Equal(t, 1, len(firstMapping))

	var issuerId3 string
	var keyId3 string
	for i, k := range firstMapping {
		issuerId3 = i
		keyId3 = k
	}

	// Performing this again should result in no key/issuer ids being imported/generated.
	resp, err = CBWrite(b, s, "config/ca", map[string]interface{}{
		"pem_bundle": pemBundleRootCA,
	})
	require.NoError(t, err)
	require.NotNil(t, resp, "expected ca info")
	secondMapping := resp.Data["mapping"].(map[string]string)
	secondImportedKeys := resp.Data["imported_keys"]
	secondImportedIssuers := resp.Data["imported_issuers"]
	secondExistingKeys := resp.Data["existing_keys"]
	secondExistingIssuers := resp.Data["existing_issuers"]

	require.Empty(t, secondImportedKeys)
	require.Empty(t, secondImportedIssuers)
	require.Contains(t, secondExistingKeys, keyId3)
	require.Contains(t, secondExistingIssuers, issuerId3)
	require.Equal(t, 1, len(secondMapping))

	resp, err = CBDelete(b, s, "root")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 1, len(resp.Warnings))

	// Make sure we can delete twice...
	resp, err = CBDelete(b, s, "root")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 1, len(resp.Warnings))

	_, err = CBRead(b, s, "cert/ca_chain")
	require.Error(t, err, "expected an error fetching deleted ca_chain")

	// We should be able to import the same ca bundle as before and get a different key/issuer ids
	resp, err = CBWrite(b, s, "config/ca", map[string]interface{}{
		"pem_bundle": pemBundleRootCA,
	})
	require.NoError(t, err)
	require.NotNil(t, resp, "expected ca info")
	postDeleteImportedKeys := resp.Data["imported_keys"]
	postDeleteImportedIssuers := resp.Data["imported_issuers"]

	// Make sure that we actually generated different issuer and key values, then the previous import
	require.NotNil(t, postDeleteImportedKeys)
	require.NotNil(t, postDeleteImportedIssuers)
	require.NotEqual(t, postDeleteImportedKeys, firstImportedKeys)
	require.NotEqual(t, postDeleteImportedIssuers, firstImportedIssuers)

	resp, err = CBRead(b, s, "cert/ca_chain")
	require.NoError(t, err)

	caChainPostDelete := resp.Data
	if reflect.DeepEqual(r1Data, caChainPostDelete) {
		t.Fatal("ca certs from ca_chain were the same post delete, should have changed.")
	}
}

// TestBackend_SignIntermediate_EnforceLeafFlag verifies if the flag is true
// that we will leverage the issuer's configured behavior
func TestBackend_SignIntermediate_EnforceLeafFlag(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	require.NoError(t, err, "failed generating root cert")
	rootCert := parseCert(t, resp.Data["certificate"].(string))

	_, err = CBWrite(b, s, "issuer/default", map[string]interface{}{
		"leaf_not_after_behavior": "err",
	})
	require.NoError(t, err, "failed updating root issuer cert behavior")

	resp, err = CBWrite(b, s, "intermediate/generate/internal", map[string]interface{}{
		"common_name": "myint.com",
	})
	require.NoError(t, err, "failed generating intermediary CSR")
	csr := resp.Data["csr"]

	_, err = CBWrite(b, s, "root/sign-intermediate", map[string]interface{}{
		"common_name":                     "myint.com",
		"other_sans":                      "1.3.6.1.4.1.311.20.2.3;utf8:caadmin@example.com",
		"csr":                             csr,
		"ttl":                             "60h",
		"enforce_leaf_not_after_behavior": true,
	})
	require.Error(t, err, "sign-intermediate should have failed as root issuer leaf behavior is set to err")

	// Now test with permit, the old default behavior
	_, err = CBWrite(b, s, "issuer/default", map[string]interface{}{
		"leaf_not_after_behavior": "permit",
	})
	require.NoError(t, err, "failed updating root issuer cert behavior to permit")

	resp, err = CBWrite(b, s, "root/sign-intermediate", map[string]interface{}{
		"common_name":                     "myint.com",
		"other_sans":                      "1.3.6.1.4.1.311.20.2.3;utf8:caadmin@example.com",
		"csr":                             csr,
		"ttl":                             "60h",
		"enforce_leaf_not_after_behavior": true,
	})
	require.NoError(t, err, "failed to sign intermediary CA with permit as issuer")
	intCert := parseCert(t, resp.Data["certificate"].(string))

	require.Truef(t, rootCert.NotAfter.Before(intCert.NotAfter),
		"root cert notAfter %v was not before ca cert's notAfter %v", rootCert.NotAfter, intCert.NotAfter)
}

func TestBackend_SignIntermediate_AllowedPastCAValidity(t *testing.T) {
	t.Parallel()
	b_root, s_root := CreateBackendWithStorage(t)
	b_int, s_int := CreateBackendWithStorage(t)
	var err error

	// Direct issuing from root
	resp, err := CBWrite(b_root, s_root, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	rootCert := parseCert(t, resp.Data["certificate"].(string))

	_, err = CBWrite(b_root, s_root, "roles/test", map[string]interface{}{
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"allow_any_name":     true,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = CBWrite(b_int, s_int, "intermediate/generate/internal", map[string]interface{}{
		"common_name": "myint.com",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b_root.Route("intermediate/generate/internal"), logical.UpdateOperation), resp, true)
	require.Contains(t, resp.Data, "key_id")
	intKeyId := resp.Data["key_id"].(issuing.KeyID)
	csr := resp.Data["csr"]

	resp, err = CBRead(b_int, s_int, "key/"+intKeyId.String())
	require.NoError(t, err)
	require.NotNil(t, resp, "expected a response")
	intSkid := resp.Data["subject_key_id"].(string)

	if err != nil {
		t.Fatal(err)
	}

	_, err = CBWrite(b_root, s_root, "sign/test", map[string]interface{}{
		"common_name": "myint.com",
		"csr":         csr,
		"ttl":         "60h",
	})
	require.ErrorContains(t, err, "that is beyond the expiration of the CA certificate")

	_, err = CBWrite(b_root, s_root, "sign-verbatim/test", map[string]interface{}{
		"common_name": "myint.com",
		"other_sans":  "1.3.6.1.4.1.311.20.2.3;utf8:caadmin@example.com",
		"csr":         csr,
		"ttl":         "60h",
	})
	require.ErrorContains(t, err, "that is beyond the expiration of the CA certificate")

	resp, err = CBWrite(b_root, s_root, "root/sign-intermediate", map[string]interface{}{
		"common_name": "myint.com",
		"other_sans":  "1.3.6.1.4.1.311.20.2.3;utf8:caadmin@example.com",
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

	cert := parseCert(t, resp.Data["certificate"].(string))
	certSkid := certutil.GetHexFormatted(cert.SubjectKeyId, ":")
	require.Equal(t, intSkid, certSkid)

	require.Equal(t, rootCert.NotAfter, cert.NotAfter, "intermediary cert's NotAfter did not match root cert's NotAfter")
	require.Contains(t, resp.Warnings, intCaTruncatationWarning, "missing warning about intermediary CA notAfter truncation")
}

func TestBackend_ConsulSignLeafWithLegacyRole(t *testing.T) {
	t.Parallel()
	// create the backend
	b, s := CreateBackendWithStorage(t)

	// generate root
	data, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	require.NoError(t, err, "failed generating internal root cert")
	rootCaPem := data.Data["certificate"].(string)

	// Create a signing role like Consul did with the default args prior to Vault 1.10
	_, err = CBWrite(b, s, "roles/test", map[string]interface{}{
		"allow_any_name":         true,
		"allowed_serial_numbers": []string{"MySerialNumber"},
		"key_type":               "any",
		"key_bits":               "2048",
		"signature_bits":         "256",
	})
	require.NoError(t, err, "failed creating legacy role")

	_, csrPem := generateTestCsr(t, certutil.ECPrivateKey, 256)
	data, err = CBWrite(b, s, "sign/test", map[string]interface{}{
		"csr": csrPem,
	})
	require.NoError(t, err, "failed signing csr")
	certAsPem := data.Data["certificate"].(string)

	signedCert := parseCert(t, certAsPem)
	rootCert := parseCert(t, rootCaPem)
	requireSignedBy(t, signedCert, rootCert)
}

func TestBackend_SignSelfIssued(t *testing.T) {
	t.Parallel()
	// create the backend
	b, storage := CreateBackendWithStorage(t)

	// generate root
	rootData := map[string]interface{}{
		"common_name": "test.com",
		"ttl":         "172800",
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "root/generate/internal",
		Storage:    storage,
		Data:       rootData,
		MountPoint: "pki/",
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to generate root, %#v", *resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	key, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "foo.bar.com",
		},
		SerialNumber:          big.NewInt(1234),
		IsCA:                  false,
		BasicConstraintsValid: true,
	}

	ss, _ := getSelfSigned(t, template, template, key)
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-self-issued",
		Storage:   storage,
		Data: map[string]interface{}{
			"certificate": ss,
		},
		MountPoint: "pki/",
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
	ss, ssCert := getSelfSigned(t, template, issuer, key)
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-self-issued",
		Storage:   storage,
		Data: map[string]interface{}{
			"certificate": ss,
		},
		MountPoint: "pki/",
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

	ss, _ = getSelfSigned(t, template, template, key)
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-self-issued",
		Storage:   storage,
		Data: map[string]interface{}{
			"certificate": ss,
		},
		MountPoint: "pki/",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("root/sign-self-issued"), logical.UpdateOperation), resp, true)
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

	sc := b.makeStorageContext(context.Background(), storage)
	signingBundle, err := sc.fetchCAInfo(defaultRef, issuing.ReadOnlyUsage)
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

// TestBackend_SignSelfIssued_DifferentTypes tests the functionality of the
// require_matching_certificate_algorithms flag.
func TestBackend_SignSelfIssued_DifferentTypes(t *testing.T) {
	t.Parallel()
	// create the backend
	b, storage := CreateBackendWithStorage(t)

	// generate root
	rootData := map[string]interface{}{
		"common_name": "test.com",
		"ttl":         "172800",
		"key_type":    "ec",
		"key_bits":    "521",
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "root/generate/internal",
		Storage:    storage,
		Data:       rootData,
		MountPoint: "pki/",
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to generate root, %#v", *resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	key, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "foo.bar.com",
		},
		SerialNumber:          big.NewInt(1234),
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	// Tests absent the flag
	ss, _ := getSelfSigned(t, template, template, key)
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-self-issued",
		Storage:   storage,
		Data: map[string]interface{}{
			"certificate": ss,
		},
		MountPoint: "pki/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response")
	}

	// Set CA to true, but leave issuer alone
	template.IsCA = true

	// Tests with flag present but false
	ss, _ = getSelfSigned(t, template, template, key)
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-self-issued",
		Storage:   storage,
		Data: map[string]interface{}{
			"certificate": ss,
			"require_matching_certificate_algorithms": false,
		},
		MountPoint: "pki/",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response")
	}

	// Test with flag present and true
	ss, _ = getSelfSigned(t, template, template, key)
	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-self-issued",
		Storage:   storage,
		Data: map[string]interface{}{
			"certificate": ss,
			"require_matching_certificate_algorithms": true,
		},
		MountPoint: "pki/",
	})
	if err == nil {
		t.Fatal("expected error due to mismatched algorithms")
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
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	var err error
	var resp *logical.Response
	var certStr string
	var block *pem.Block
	var cert *x509.Certificate

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = CBWrite(b, s, "roles/test", map[string]interface{}{
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
	resp, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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
	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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

	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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

	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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

	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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

	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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
	resp, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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
	resp, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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
	otherNames := []string{
		fmt.Sprintf("%s;%s:%s", oid1, type1, val1),
		fmt.Sprintf("%s;%s:%s", oid2, type2, val2),
	}
	resp, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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
	expectedOtherNames := []certutil.OtherNameUtf8{{oid1, val1}, {oid2, val2}}
	foundOtherNames, err := getOtherSANsFromX509Extensions(cert.Extensions)
	if err != nil {
		t.Fatal(err)
	}
	// Sort our returned list as SANS are built internally with a map so ordering can be inconsistent
	slices.SortFunc(foundOtherNames, func(a, b certutil.OtherNameUtf8) int { return cmp.Compare(a.Oid, b.Oid) })

	if diff := deep.Equal(expectedOtherNames, foundOtherNames); len(diff) != 0 {
		t.Errorf("unexpected otherNames: %v", diff)
	}
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		t.Logf("certificate 3 to check:\n%s", certStr)
	}
}

func TestBackend_AllowedSerialNumbers(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	var err error
	var resp *logical.Response
	var certStr string
	var block *pem.Block
	var cert *x509.Certificate

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// First test that Serial Numbers are not allowed
	_, err = CBWrite(b, s, "roles/test", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
		"common_name": "foobar",
		"ttl":         "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
		"common_name":   "foobar",
		"ttl":           "1h",
		"serial_number": "foobar",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// Update the role to allow serial numbers
	_, err = CBWrite(b, s, "roles/test", map[string]interface{}{
		"allow_any_name":         true,
		"enforce_hostnames":      false,
		"allowed_serial_numbers": "f00*,b4r*",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
		"common_name": "foobar",
		"ttl":         "1h",
		// Not a valid serial number
		"serial_number": "foobar",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	// Valid for first possibility
	resp, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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
	resp, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	var err error

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = CBWrite(b, s, "roles/test", map[string]interface{}{
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
	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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
	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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
	_, err = CBWrite(b, s, "issue/test", map[string]interface{}{
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
	resp, err := CBWrite(b, s, "issue/test", map[string]interface{}{
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

func TestBackend_AllowedURISANsTemplate(t *testing.T) {
	t.Parallel()
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
		"allowed_uri_sans": []string{
			"spiffe://domain/{{identity.entity.aliases." + userpassAccessor + ".name}}",
			"spiffe://domain/{{identity.entity.aliases." + userpassAccessor + ".name}}/*", "spiffe://domain/foo",
		},
		"allowed_uri_sans_template": true,
		"require_cn":                false,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue certificate with identity templating
	client.SetToken(userpassToken)
	_, err = client.Logical().Write("pki/issue/test", map[string]interface{}{"uri_sans": "spiffe://domain/userpassname, spiffe://domain/foo"})
	if err != nil {
		t.Fatal(err)
	}

	// Issue certificate with identity templating and glob
	client.SetToken(userpassToken)
	_, err = client.Logical().Write("pki/issue/test", map[string]interface{}{"uri_sans": "spiffe://domain/userpassname/bar"})
	if err != nil {
		t.Fatal(err)
	}

	// Issue certificate with non-matching identity template parameter
	client.SetToken(userpassToken)
	_, err = client.Logical().Write("pki/issue/test", map[string]interface{}{"uri_sans": "spiffe://domain/unknownuser"})
	if err == nil {
		t.Fatal(err)
	}

	// Set allowed_uri_sans_template to false.
	_, err = client.Logical().Write("pki/roles/test", map[string]interface{}{
		"allowed_uri_sans_template": false,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue certificate with userpassToken.
	_, err = client.Logical().Write("pki/issue/test", map[string]interface{}{"uri_sans": "spiffe://domain/users/userpassname"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBackend_AllowedDomainsTemplate(t *testing.T) {
	t.Parallel()
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

	// Login userpass for test role and set client token
	userpassAuth, err := auth.NewUserpassAuth("userpassname", &auth.Password{FromString: "test"})
	if err != nil {
		t.Fatal(err)
	}

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
		"allowed_domains": []string{
			"foobar.com", "zipzap.com", "{{identity.entity.aliases." + userpassAccessor + ".name}}",
			"foo.{{identity.entity.aliases." + userpassAccessor + ".name}}.example.com",
		},
		"allowed_domains_template": true,
		"allow_bare_domains":       true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue certificate with userpassToken.
	secret, err := client.Auth().Login(context.TODO(), userpassAuth)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil || secret == nil {
		t.Fatal(err)
	}
	_, err = client.Logical().Write("pki/issue/test", map[string]interface{}{"common_name": "userpassname"})
	if err != nil {
		t.Fatal(err)
	}

	// Issue certificate for foobar.com to verify allowed_domain_template doesn't break plain domains.
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

func TestReadWriteDeleteRoles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
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

	// Mount PKI.
	err := client.Sys().MountWithContext(ctx, "pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "60h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().ReadWithContext(ctx, "pki/roles/test")
	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		t.Fatalf("response should have been emtpy but was:\n%#v", resp)
	}

	// Write role PKI.
	_, err = client.Logical().WriteWithContext(ctx, "pki/roles/test", map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	// Read the role.
	resp, err = client.Logical().ReadWithContext(ctx, "pki/roles/test")
	if err != nil {
		t.Fatal(err)
	}

	if resp.Data == nil {
		t.Fatal("default data within response was nil when it should have contained data")
	}

	// Validate that we have not changed any defaults unknowingly
	expectedData := map[string]interface{}{
		"key_type":                           "rsa",
		"use_csr_sans":                       true,
		"serial_number_source":               "json-csr",
		"client_flag":                        true,
		"allowed_serial_numbers":             []interface{}{},
		"generate_lease":                     false,
		"signature_bits":                     json.Number("256"),
		"use_pss":                            false,
		"allowed_domains":                    []interface{}{},
		"allowed_uri_sans_template":          false,
		"enforce_hostnames":                  true,
		"policy_identifiers":                 []interface{}{},
		"require_cn":                         true,
		"allowed_domains_template":           false,
		"allow_token_displayname":            false,
		"country":                            []interface{}{},
		"not_after":                          "",
		"postal_code":                        []interface{}{},
		"use_csr_common_name":                true,
		"allow_localhost":                    true,
		"allow_subdomains":                   false,
		"allow_wildcard_certificates":        true,
		"allowed_other_sans":                 []interface{}{},
		"allowed_uri_sans":                   []interface{}{},
		"basic_constraints_valid_for_non_ca": false,
		"key_usage":                          []interface{}{"DigitalSignature", "KeyAgreement", "KeyEncipherment"},
		"not_before_duration":                json.Number("30"),
		"allow_glob_domains":                 false,
		"ttl":                                json.Number("0"),
		"ou":                                 []interface{}{},
		"email_protection_flag":              false,
		"locality":                           []interface{}{},
		"server_flag":                        true,
		"allow_bare_domains":                 false,
		"allow_ip_sans":                      true,
		"ext_key_usage_oids":                 []interface{}{},
		"allow_any_name":                     false,
		"ext_key_usage":                      []interface{}{},
		"key_bits":                           json.Number("2048"),
		"max_ttl":                            json.Number("0"),
		"no_store":                           false,
		"organization":                       []interface{}{},
		"province":                           []interface{}{},
		"street_address":                     []interface{}{},
		"code_signing_flag":                  false,
		"issuer_ref":                         "default",
		"cn_validations":                     []interface{}{"email", "hostname"},
		"allowed_user_ids":                   []interface{}{},
	}

	if issuing.MetadataPermitted {
		expectedData["no_store_metadata"] = false
	}

	if diff := deep.Equal(expectedData, resp.Data); len(diff) > 0 {
		t.Fatalf("pki role default values have changed, diff: %v", diff)
	}

	_, err = client.Logical().DeleteWithContext(ctx, "pki/roles/test")
	if err != nil {
		t.Fatal(err)
	}

	resp, err = client.Logical().ReadWithContext(ctx, "pki/roles/test")
	if err != nil {
		t.Fatal(err)
	}

	if resp != nil {
		t.Fatalf("response should have been empty but was:\n%#v", resp)
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

	rak, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
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
	_, err = certutil.GetSubjKeyID(rak)
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

	_, edk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	marshaledKey, err = x509.MarshalPKCS8PrivateKey(edk)
	if err != nil {
		panic(err)
	}
	keyPEMBlock = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: marshaledKey,
	}
	edCAKey = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
	if err != nil {
		panic(err)
	}
	_, err = certutil.GetSubjKeyID(edk)
	if err != nil {
		panic(err)
	}
	caBytes, err = x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, edk.Public(), edk)
	if err != nil {
		panic(err)
	}
	caCertPEMBlock = &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	edCACert = strings.TrimSpace(string(pem.EncodeToMemory(caCertPEMBlock)))
}

func TestBackend_RevokePlusTidy_Intermediate(t *testing.T) {
	// Use a ridiculously long time to minimize the chance
	// that we have to deal with more than one interval.
	// InMemSink rounds down to an interval boundary rather than
	// starting one at the time of initialization.
	//
	// This test is not parallelizable.
	inmemSink := metrics.NewInmemSink(
		1000000*time.Hour,
		2000000*time.Hour)

	metricsConf := metrics.DefaultConfig("")
	metricsConf.EnableHostname = false
	metricsConf.EnableHostnameLabel = false
	metricsConf.EnableServiceLabel = false
	metricsConf.EnableTypePrefix = false

	_, err := metrics.NewGlobal(metricsConf, inmemSink)
	if err != nil {
		t.Fatal(err)
	}

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

	// Set up Metric Configuration, then restart to enable it
	_, err = client.Logical().Write("pki/config/auto-tidy", map[string]interface{}{
		"maintain_stored_certificate_counts":       true,
		"publish_stored_certificate_count_metrics": true,
	})
	require.NoError(t, err, "failed calling auto-tidy")
	_, err = client.Logical().Write("/sys/plugins/reload/backend", map[string]interface{}{
		"mounts": "pki/",
	})
	require.NoError(t, err, "failed calling backend reload")

	// Check the metrics initialized in order to calculate backendUUID for /pki
	// BackendUUID not consistent during tests with UUID from /sys/mounts/pki
	metricsSuffix := "total_certificates_stored"
	backendUUID := ""
	mostRecentInterval := inmemSink.Data()[len(inmemSink.Data())-1]
	for _, existingGauge := range mostRecentInterval.Gauges {
		if strings.HasSuffix(existingGauge.Name, metricsSuffix) {
			expandedGaugeName := existingGauge.Name
			backendUUID = strings.Split(expandedGaugeName, ".")[2]
			break
		}
	}
	if backendUUID == "" {
		t.Fatalf("No Gauge Found ending with %s", metricsSuffix)
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
	// Set up Metric Configuration, then restart to enable it
	_, err = client.Logical().Write("pki2/config/auto-tidy", map[string]interface{}{
		"maintain_stored_certificate_counts":       true,
		"publish_stored_certificate_count_metrics": true,
	})
	_, err = client.Logical().Write("/sys/plugins/reload/backend", map[string]interface{}{
		"mounts": "pki2/",
	})

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
		"ttl":                   "10s",
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
	if err != nil {
		t.Fatal(err)
	}

	// Check the cert-count metrics
	expectedCertCountGaugeMetrics := map[string]float32{
		"secrets.pki." + backendUUID + ".total_revoked_certificates_stored": 1,
		"secrets.pki." + backendUUID + ".total_certificates_stored":         1,
	}
	mostRecentInterval = inmemSink.Data()[len(inmemSink.Data())-1]
	for gauge, value := range expectedCertCountGaugeMetrics {
		if _, ok := mostRecentInterval.Gauges[gauge]; !ok {
			t.Fatalf("Expected metrics to include a value for gauge %s", gauge)
		}
		if value != mostRecentInterval.Gauges[gauge].Value {
			t.Fatalf("Expected value metric %s to be %f but got %f", gauge, value, mostRecentInterval.Gauges[gauge].Value)
		}
	}

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

	// Get CRL and ensure the tidied cert is still in the list after the tidy
	// operation since it's not past the NotAfter (ttl) value yet.
	crl := getParsedCrl(t, client, "pki")

	revokedCerts := crl.TBSCertList.RevokedCertificates
	if len(revokedCerts) == 0 {
		t.Fatal("expected CRL to be non-empty")
	}

	sn := certutil.GetHexFormatted(revokedCerts[0].SerialNumber.Bytes(), ":")
	if sn != intermediateCertSerial {
		t.Fatalf("expected: %v, got: %v", intermediateCertSerial, sn)
	}

	// Wait for cert to expire
	time.Sleep(10 * time.Second)

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

	// Issue a tidy-status on /pki
	{
		tidyStatus, err := client.Logical().Read("pki/tidy-status")
		if err != nil {
			t.Fatal(err)
		}
		expectedData := map[string]interface{}{
			"safety_buffer":                         json.Number("1"),
			"issuer_safety_buffer":                  json.Number("31536000"),
			"revocation_queue_safety_buffer":        json.Number("172800"),
			"tidy_cert_store":                       true,
			"tidy_revoked_certs":                    true,
			"tidy_revoked_cert_issuer_associations": false,
			"tidy_expired_issuers":                  false,
			"tidy_move_legacy_ca_bundle":            false,
			"tidy_revocation_queue":                 false,
			"tidy_cross_cluster_revoked_certs":      false,
			"tidy_cert_metadata":                    false,
			"tidy_cmpv2_nonce_store":                false,
			"pause_duration":                        "0s",
			"state":                                 "Finished",
			"error":                                 nil,
			"time_started":                          nil,
			"time_finished":                         nil,
			"last_auto_tidy_finished":               nil,
			"message":                               nil,
			"cert_store_deleted_count":              json.Number("1"),
			"revoked_cert_deleted_count":            json.Number("1"),
			"missing_issuer_cert_count":             json.Number("0"),
			"current_cert_store_count":              json.Number("0"),
			"current_revoked_cert_count":            json.Number("0"),
			"revocation_queue_deleted_count":        json.Number("0"),
			"cross_revoked_cert_deleted_count":      json.Number("0"),
			"internal_backend_uuid":                 backendUUID,
			"tidy_acme":                             false,
			"acme_account_safety_buffer":            json.Number("2592000"),
			"acme_orders_deleted_count":             json.Number("0"),
			"acme_account_revoked_count":            json.Number("0"),
			"acme_account_deleted_count":            json.Number("0"),
			"total_acme_account_count":              json.Number("0"),
			"cert_metadata_deleted_count":           json.Number("0"),
			"cmpv2_nonce_deleted_count":             json.Number("0"),
		}
		// Let's copy the times from the response so that we can use deep.Equal()
		timeStarted, ok := tidyStatus.Data["time_started"]
		if !ok || timeStarted == "" {
			t.Fatal("Expected tidy status response to include a value for time_started")
		}
		expectedData["time_started"] = timeStarted
		timeFinished, ok := tidyStatus.Data["time_finished"]
		if !ok || timeFinished == "" {
			t.Fatal("Expected tidy status response to include a value for time_finished")
		}
		expectedData["time_finished"] = timeFinished
		expectedData["last_auto_tidy_finished"] = tidyStatus.Data["last_auto_tidy_finished"]

		if diff := deep.Equal(expectedData, tidyStatus.Data); diff != nil {
			t.Fatal(diff)
		}
	}
	// Check the tidy metrics
	{
		// Map of gauges to expected value
		expectedGauges := map[string]float32{
			"secrets.pki.tidy.cert_store_current_entry":                         0,
			"secrets.pki.tidy.cert_store_total_entries":                         1,
			"secrets.pki.tidy.revoked_cert_current_entry":                       0,
			"secrets.pki.tidy.revoked_cert_total_entries":                       1,
			"secrets.pki.tidy.start_time_epoch":                                 0,
			"secrets.pki." + backendUUID + ".total_certificates_stored":         0,
			"secrets.pki." + backendUUID + ".total_revoked_certificates_stored": 0,
			"secrets.pki.tidy.cert_store_total_entries_remaining":               0,
			"secrets.pki.tidy.revoked_cert_total_entries_remaining":             0,
		}
		// Map of counters to the sum of the metrics for that counter
		expectedCounters := map[string]float64{
			"secrets.pki.tidy.cert_store_deleted_count":   1,
			"secrets.pki.tidy.revoked_cert_deleted_count": 1,
			"secrets.pki.tidy.success":                    2,
			// Note that "secrets.pki.tidy.failure" won't be in the captured metrics
		}

		// If the metrics span more than one interval, skip the checks
		intervals := inmemSink.Data()
		if len(intervals) == 1 {
			interval := inmemSink.Data()[0]

			for gauge, value := range expectedGauges {
				if _, ok := interval.Gauges[gauge]; !ok {
					t.Fatalf("Expected metrics to include a value for gauge %s", gauge)
				}
				if value != interval.Gauges[gauge].Value {
					t.Fatalf("Expected value metric %s to be %f but got %f", gauge, value, interval.Gauges[gauge].Value)
				}

			}
			for counter, value := range expectedCounters {
				if _, ok := interval.Counters[counter]; !ok {
					t.Fatalf("Expected metrics to include a value for couter %s", counter)
				}
				if value != interval.Counters[counter].Sum {
					t.Fatalf("Expected the sum of metric %s to be %f but got %f", counter, value, interval.Counters[counter].Sum)
				}
			}

			tidyDuration, ok := interval.Samples["secrets.pki.tidy.duration"]
			if !ok {
				t.Fatal("Expected metrics to include a value for sample secrets.pki.tidy.duration")
			}
			if tidyDuration.Count <= 0 {
				t.Fatalf("Expected metrics to have count > 0 for sample secrets.pki.tidy.duration, but got %d", tidyDuration.Count)
			}
		}
	}

	crl = getParsedCrl(t, client, "pki")

	revokedCerts = crl.TBSCertList.RevokedCertificates
	if len(revokedCerts) != 0 {
		t.Fatal("expected CRL to be empty")
	}
}

func TestBackend_Root_FullCAChain(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		testName string
		keyType  string
	}{
		{testName: "RSA", keyType: "rsa"},
		{testName: "ED25519", keyType: "ed25519"},
		{testName: "EC", keyType: "ec"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testName, func(t *testing.T) {
			runFullCAChainTest(t, tc.keyType)
		})
	}
}

func runFullCAChainTest(t *testing.T, keyType string) {
	// Generate a root CA at /pki-root
	b_root, s_root := CreateBackendWithStorage(t)

	var err error

	resp, err := CBWrite(b_root, s_root, "root/generate/exported", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_type":    keyType,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected ca info")
	}
	rootData := resp.Data
	rootCert := rootData["certificate"].(string)

	// Validate that root's /cert/ca-chain now contains the certificate.
	resp, err = CBRead(b_root, s_root, "cert/ca_chain")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected intermediate chain information")
	}

	fullChain := resp.Data["ca_chain"].(string)
	requireCertInCaChainString(t, fullChain, rootCert, "expected root cert within root cert/ca_chain")

	// Make sure when we issue a leaf certificate we get the full chain back.
	_, err = CBWrite(b_root, s_root, "roles/example", map[string]interface{}{
		"allowed_domains":  "example.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
	})
	require.NoError(t, err, "error setting up pki root role: %v", err)

	resp, err = CBWrite(b_root, s_root, "issue/example", map[string]interface{}{
		"common_name": "test.example.com",
		"ttl":         "5m",
	})
	require.NoError(t, err, "error issuing certificate from pki root: %v", err)
	fullChainArray := resp.Data["ca_chain"].([]string)
	requireCertInCaChainArray(t, fullChainArray, rootCert, "expected root cert within root issuance pki-root/issue/example")

	// Now generate an intermediate at /pki-intermediate, signed by the root.
	b_int, s_int := CreateBackendWithStorage(t)

	resp, err = CBWrite(b_int, s_int, "intermediate/generate/exported", map[string]interface{}{
		"common_name": "intermediate myvault.com",
		"key_type":    keyType,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected intermediate CSR info")
	}
	intermediateData := resp.Data
	intermediateKey := intermediateData["private_key"].(string)

	resp, err = CBWrite(b_root, s_root, "root/sign-intermediate", map[string]interface{}{
		"csr":    intermediateData["csr"],
		"format": "pem",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected signed intermediate info")
	}
	intermediateSignedData := resp.Data
	intermediateCert := intermediateSignedData["certificate"].(string)

	rootCaCert := parseCert(t, rootCert)
	intermediaryCaCert := parseCert(t, intermediateCert)
	requireSignedBy(t, intermediaryCaCert, rootCaCert)
	intermediateCaChain := intermediateSignedData["ca_chain"].([]string)

	require.Equal(t, parseCert(t, intermediateCaChain[0]), intermediaryCaCert, "intermediate signed cert should have been part of ca_chain")
	require.Equal(t, parseCert(t, intermediateCaChain[1]), rootCaCert, "root cert should have been part of ca_chain")

	_, err = CBWrite(b_int, s_int, "intermediate/set-signed", map[string]interface{}{
		"certificate": intermediateCert + "\n" + rootCert + "\n",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Validate that intermediate's ca_chain field now includes the full
	// chain.
	resp, err = CBRead(b_int, s_int, "cert/ca_chain")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected intermediate chain information")
	}

	// Verify we have a proper CRL now
	crl := getParsedCrlFromBackend(t, b_int, s_int, "crl")
	require.Equal(t, 0, len(crl.TBSCertList.RevokedCertificates))

	fullChain = resp.Data["ca_chain"].(string)
	requireCertInCaChainString(t, fullChain, intermediateCert, "expected full chain to contain intermediate certificate from pki-intermediate/cert/ca_chain")
	requireCertInCaChainString(t, fullChain, rootCert, "expected full chain to contain root certificate from pki-intermediate/cert/ca_chain")

	// Make sure when we issue a leaf certificate we get the full chain back.
	_, err = CBWrite(b_int, s_int, "roles/example", map[string]interface{}{
		"allowed_domains":  "example.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
	})
	require.NoError(t, err, "error setting up pki intermediate role: %v", err)

	resp, err = CBWrite(b_int, s_int, "issue/example", map[string]interface{}{
		"common_name": "test.example.com",
		"ttl":         "5m",
	})
	require.NoError(t, err, "error issuing certificate from pki intermediate: %v", err)
	fullChainArray = resp.Data["ca_chain"].([]string)
	requireCertInCaChainArray(t, fullChainArray, intermediateCert, "expected full chain to contain intermediate certificate from pki-intermediate/issue/example")
	requireCertInCaChainArray(t, fullChainArray, rootCert, "expected full chain to contain root certificate from pki-intermediate/issue/example")

	// Finally, import this signing cert chain into a new mount to ensure
	// "external" CAs behave as expected.
	b_ext, s_ext := CreateBackendWithStorage(t)

	_, err = CBWrite(b_ext, s_ext, "config/ca", map[string]interface{}{
		"pem_bundle": intermediateKey + "\n" + intermediateCert + "\n" + rootCert + "\n",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Validate the external chain information was loaded correctly.
	resp, err = CBRead(b_ext, s_ext, "cert/ca_chain")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected intermediate chain information")
	}

	fullChain = resp.Data["ca_chain"].(string)
	if strings.Count(fullChain, intermediateCert) != 1 {
		t.Fatalf("expected full chain to contain intermediate certificate; got %v occurrences", strings.Count(fullChain, intermediateCert))
	}
	if strings.Count(fullChain, rootCert) != 1 {
		t.Fatalf("expected full chain to contain root certificate; got %v occurrences", strings.Count(fullChain, rootCert))
	}

	// Now issue a short-lived certificate from our pki-external.
	_, err = CBWrite(b_ext, s_ext, "roles/example", map[string]interface{}{
		"allowed_domains":  "example.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
	})
	require.NoError(t, err, "error setting up pki role: %v", err)

	resp, err = CBWrite(b_ext, s_ext, "issue/example", map[string]interface{}{
		"common_name": "test.example.com",
		"ttl":         "5m",
	})
	require.NoError(t, err, "error issuing certificate: %v", err)
	require.NotNil(t, resp, "got nil response from issuing request")
	issueCrtAsPem := resp.Data["certificate"].(string)
	issuedCrt := parseCert(t, issueCrtAsPem)

	// Verify that the certificates are signed by the intermediary CA key...
	requireSignedBy(t, issuedCrt, intermediaryCaCert)

	// Test that we can request that the root ca certificate not appear in the ca_chain field
	resp, err = CBWrite(b_ext, s_ext, "issue/example", map[string]interface{}{
		"common_name":             "test.example.com",
		"ttl":                     "5m",
		"remove_roots_from_chain": "true",
	})
	requireSuccessNonNilResponse(t, resp, err, "error issuing certificate when removing self signed")
	fullChain = strings.Join(resp.Data["ca_chain"].([]string), "\n")
	if strings.Count(fullChain, intermediateCert) != 1 {
		t.Fatalf("expected full chain to contain intermediate certificate; got %v occurrences", strings.Count(fullChain, intermediateCert))
	}
	if strings.Count(fullChain, rootCert) != 0 {
		t.Fatalf("expected full chain to NOT contain root certificate; got %v occurrences", strings.Count(fullChain, rootCert))
	}
}

func requireCertInCaChainArray(t *testing.T, chain []string, cert string, msgAndArgs ...interface{}) {
	var fullChain string
	for _, caCert := range chain {
		fullChain = fullChain + "\n" + caCert
	}

	requireCertInCaChainString(t, fullChain, cert, msgAndArgs)
}

func requireCertInCaChainString(t *testing.T, chain string, cert string, msgAndArgs ...interface{}) {
	count := strings.Count(chain, cert)
	if count != 1 {
		failMsg := fmt.Sprintf("Found %d occurrances of the cert in the provided chain", count)
		require.FailNow(t, failMsg, msgAndArgs...)
	}
}

type MultiBool int

const (
	MFalse MultiBool = iota
	MTrue  MultiBool = iota
	MAny   MultiBool = iota
)

func (o MultiBool) ToValues() []bool {
	if o == MTrue {
		return []bool{true}
	}

	if o == MFalse {
		return []bool{false}
	}

	if o == MAny {
		return []bool{true, false}
	}

	return []bool{}
}

type IssuanceRegression struct {
	AllowedDomains            []string
	AllowBareDomains          MultiBool
	AllowGlobDomains          MultiBool
	AllowSubdomains           MultiBool
	AllowLocalhost            MultiBool
	AllowWildcardCertificates MultiBool
	CNValidations             []string
	CommonName                string
	Issued                    bool
}

func RoleIssuanceRegressionHelper(t *testing.T, b *backend, s logical.Storage, index int, test IssuanceRegression) int {
	tested := 0
	for _, AllowBareDomains := range test.AllowBareDomains.ToValues() {
		for _, AllowGlobDomains := range test.AllowGlobDomains.ToValues() {
			for _, AllowSubdomains := range test.AllowSubdomains.ToValues() {
				for _, AllowLocalhost := range test.AllowLocalhost.ToValues() {
					for _, AllowWildcardCertificates := range test.AllowWildcardCertificates.ToValues() {
						role := fmt.Sprintf("issuance-regression-%d-bare-%v-glob-%v-subdomains-%v-localhost-%v-wildcard-%v", index, AllowBareDomains, AllowGlobDomains, AllowSubdomains, AllowLocalhost, AllowWildcardCertificates)
						_, err := CBWrite(b, s, "roles/"+role, map[string]interface{}{
							"allowed_domains":             test.AllowedDomains,
							"allow_bare_domains":          AllowBareDomains,
							"allow_glob_domains":          AllowGlobDomains,
							"allow_subdomains":            AllowSubdomains,
							"allow_localhost":             AllowLocalhost,
							"allow_wildcard_certificates": AllowWildcardCertificates,
							"cn_validations":              test.CNValidations,
							// TODO: test across this vector as well. Currently certain wildcard
							// matching is broken with it enabled (such as x*x.foo).
							"enforce_hostnames": false,
							"key_type":          "ec",
							"key_bits":          256,
							"no_store":          true,
							// With the CN Validations field, ensure we prevent CN from appearing
							// in SANs.
						})
						if err != nil {
							t.Fatal(err)
						}

						resp, err := CBWrite(b, s, "issue/"+role, map[string]interface{}{
							"common_name":          test.CommonName,
							"exclude_cn_from_sans": true,
						})

						haveErr := err != nil || resp == nil
						expectErr := !test.Issued

						if haveErr != expectErr {
							t.Fatalf("issuance regression test [%d] failed: haveErr: %v, expectErr: %v, err: %v, resp: %v, test case: %v, role: %v", index, haveErr, expectErr, err, resp, test, role)
						}

						tested += 1
					}
				}
			}
		}
	}

	return tested
}

func TestBackend_Roles_IssuanceRegression(t *testing.T) {
	t.Parallel()
	// Regression testing of role's issuance policy.
	testCases := []IssuanceRegression{
		// allowed, bare, glob, subdomains, localhost, wildcards, cn, issued

		// === Globs not allowed but used === //
		// Allowed contains globs, but globbing not allowed, resulting in all
		// issuances failing. Note that tests against issuing a wildcard with
		// a bare domain will be covered later.
		/*  0 */ {[]string{"*.*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "baz.fud.bar.foo", false},
		/*  1 */ {[]string{"*.*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "*.fud.bar.foo", false},
		/*  2 */ {[]string{"*.*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "fud.bar.foo", false},
		/*  3 */ {[]string{"*.*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "*.bar.foo", false},
		/*  4 */ {[]string{"*.*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "bar.foo", false},
		/*  5 */ {[]string{"*.*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "*.foo", false},
		/*  6 */ {[]string{"*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "foo", false},
		/*  7 */ {[]string{"*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "baz.fud.bar.foo", false},
		/*  8 */ {[]string{"*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "*.fud.bar.foo", false},
		/*  9 */ {[]string{"*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "fud.bar.foo", false},
		/* 10 */ {[]string{"*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "*.bar.foo", false},
		/* 11 */ {[]string{"*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "bar.foo", false},
		/* 12 */ {[]string{"*.foo"}, MAny, MFalse, MAny, MAny, MAny, nil, "foo", false},

		// === Localhost sanity === //
		// Localhost forbidden, not matching allowed domains -> not issued
		/* 13 */ {[]string{"*.*.foo"}, MAny, MAny, MAny, MFalse, MAny, nil, "localhost", false},
		// Localhost allowed, not matching allowed domains -> issued
		/* 14 */ {[]string{"*.*.foo"}, MAny, MAny, MAny, MTrue, MAny, nil, "localhost", true},
		// Localhost allowed via allowed domains (and bare allowed), not by AllowLocalhost -> issued
		/* 15 */ {[]string{"localhost"}, MTrue, MAny, MAny, MFalse, MAny, nil, "localhost", true},
		// Localhost allowed via allowed domains (and bare not allowed), not by AllowLocalhost -> not issued
		/* 16 */ {[]string{"localhost"}, MFalse, MAny, MAny, MFalse, MAny, nil, "localhost", false},
		// Localhost allowed via allowed domains (but bare not allowed), and by AllowLocalhost -> issued
		/* 17 */ {[]string{"localhost"}, MFalse, MAny, MAny, MTrue, MAny, nil, "localhost", true},

		// === Bare wildcard issuance == //
		// allowed_domains contains one or more wildcards and bare domains allowed,
		// resulting in the cert being issued.
		/* 18 */ {[]string{"*.foo"}, MTrue, MAny, MAny, MAny, MTrue, nil, "*.foo", true},
		/* 19 */ {[]string{"*.*.foo"}, MTrue, MAny, MAny, MAny, MAny, nil, "*.*.foo", false}, // Does not conform to RFC 6125

		// === Double Leading Glob Testing === //
		// Allowed contains globs, but glob allowed so certain matches work.
		// The value of bare and localhost does not impact these results.
		/* 20 */ {[]string{"*.*.foo"}, MAny, MTrue, MFalse, MAny, MAny, nil, "baz.fud.bar.foo", true}, // glob domains allow infinite subdomains
		/* 21 */ {[]string{"*.*.foo"}, MAny, MTrue, MFalse, MAny, MTrue, nil, "*.fud.bar.foo", true}, // glob domain allows wildcard of subdomains
		/* 22 */ {[]string{"*.*.foo"}, MAny, MTrue, MFalse, MAny, MAny, nil, "fud.bar.foo", true},
		/* 23 */ {[]string{"*.*.foo"}, MAny, MTrue, MFalse, MAny, MTrue, nil, "*.bar.foo", true}, // Regression fix: Vault#13530
		/* 24 */ {[]string{"*.*.foo"}, MAny, MTrue, MFalse, MAny, MAny, nil, "bar.foo", false},
		/* 25 */ {[]string{"*.*.foo"}, MAny, MTrue, MFalse, MAny, MAny, nil, "*.foo", false},
		/* 26 */ {[]string{"*.*.foo"}, MAny, MTrue, MFalse, MAny, MAny, nil, "foo", false},

		// Allowed contains globs, but glob and subdomain both work, so we expect
		// wildcard issuance to work as well. The value of bare and localhost does
		// not impact these results.
		/* 27 */ {[]string{"*.*.foo"}, MAny, MTrue, MTrue, MAny, MAny, nil, "baz.fud.bar.foo", true},
		/* 28 */ {[]string{"*.*.foo"}, MAny, MTrue, MTrue, MAny, MTrue, nil, "*.fud.bar.foo", true},
		/* 29 */ {[]string{"*.*.foo"}, MAny, MTrue, MTrue, MAny, MAny, nil, "fud.bar.foo", true},
		/* 30 */ {[]string{"*.*.foo"}, MAny, MTrue, MTrue, MAny, MTrue, nil, "*.bar.foo", true}, // Regression fix: Vault#13530
		/* 31 */ {[]string{"*.*.foo"}, MAny, MTrue, MTrue, MAny, MAny, nil, "bar.foo", false},
		/* 32 */ {[]string{"*.*.foo"}, MAny, MTrue, MTrue, MAny, MAny, nil, "*.foo", false},
		/* 33 */ {[]string{"*.*.foo"}, MAny, MTrue, MTrue, MAny, MAny, nil, "foo", false},

		// === Single Leading Glob Testing === //
		// Allowed contains globs, but glob allowed so certain matches work.
		// The value of bare and localhost does not impact these results.
		/* 34 */ {[]string{"*.foo"}, MAny, MTrue, MFalse, MAny, MAny, nil, "baz.fud.bar.foo", true}, // glob domains allow infinite subdomains
		/* 35 */ {[]string{"*.foo"}, MAny, MTrue, MFalse, MAny, MTrue, nil, "*.fud.bar.foo", true}, // glob domain allows wildcard of subdomains
		/* 36 */ {[]string{"*.foo"}, MAny, MTrue, MFalse, MAny, MAny, nil, "fud.bar.foo", true}, // glob domains allow infinite subdomains
		/* 37 */ {[]string{"*.foo"}, MAny, MTrue, MFalse, MAny, MTrue, nil, "*.bar.foo", true}, // glob domain allows wildcards of subdomains
		/* 38 */ {[]string{"*.foo"}, MAny, MTrue, MFalse, MAny, MAny, nil, "bar.foo", true},
		/* 39 */ {[]string{"*.foo"}, MAny, MTrue, MFalse, MAny, MAny, nil, "foo", false},

		// Allowed contains globs, but glob and subdomain both work, so we expect
		// wildcard issuance to work as well. The value of bare and localhost does
		// not impact these results.
		/* 40 */ {[]string{"*.foo"}, MAny, MTrue, MTrue, MAny, MAny, nil, "baz.fud.bar.foo", true},
		/* 41 */ {[]string{"*.foo"}, MAny, MTrue, MTrue, MAny, MTrue, nil, "*.fud.bar.foo", true},
		/* 42 */ {[]string{"*.foo"}, MAny, MTrue, MTrue, MAny, MAny, nil, "fud.bar.foo", true},
		/* 43 */ {[]string{"*.foo"}, MAny, MTrue, MTrue, MAny, MTrue, nil, "*.bar.foo", true},
		/* 44 */ {[]string{"*.foo"}, MAny, MTrue, MTrue, MAny, MAny, nil, "bar.foo", true},
		/* 45 */ {[]string{"*.foo"}, MAny, MTrue, MTrue, MAny, MAny, nil, "foo", false},

		// === Only base domain name === //
		// Allowed contains only domain components, but subdomains not allowed. This
		// results in most issuances failing unless we allow bare domains, in which
		// case only the final issuance for "foo" will succeed.
		/* 46 */ {[]string{"foo"}, MAny, MAny, MFalse, MAny, MAny, nil, "baz.fud.bar.foo", false},
		/* 47 */ {[]string{"foo"}, MAny, MAny, MFalse, MAny, MAny, nil, "*.fud.bar.foo", false},
		/* 48 */ {[]string{"foo"}, MAny, MAny, MFalse, MAny, MAny, nil, "fud.bar.foo", false},
		/* 49 */ {[]string{"foo"}, MAny, MAny, MFalse, MAny, MAny, nil, "*.bar.foo", false},
		/* 50 */ {[]string{"foo"}, MAny, MAny, MFalse, MAny, MAny, nil, "bar.foo", false},
		/* 51 */ {[]string{"foo"}, MAny, MAny, MFalse, MAny, MAny, nil, "*.foo", false},
		/* 52 */ {[]string{"foo"}, MFalse, MAny, MFalse, MAny, MAny, nil, "foo", false},
		/* 53 */ {[]string{"foo"}, MTrue, MAny, MFalse, MAny, MAny, nil, "foo", true},

		// Allowed contains only domain components, and subdomains are now allowed.
		// This results in most issuances succeeding, with the exception of the
		// base foo, which is still governed by base's value.
		/* 54 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MAny, nil, "baz.fud.bar.foo", true},
		/* 55 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MTrue, nil, "*.fud.bar.foo", true},
		/* 56 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MAny, nil, "fud.bar.foo", true},
		/* 57 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MTrue, nil, "*.bar.foo", true},
		/* 58 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MAny, nil, "bar.foo", true},
		/* 59 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MTrue, nil, "*.foo", true},
		/* 60 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MTrue, nil, "x*x.foo", true}, // internal wildcards should be allowed per RFC 6125/6.4.3
		/* 61 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MTrue, nil, "*x.foo", true}, // prefix wildcards should be allowed per RFC 6125/6.4.3
		/* 62 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MTrue, nil, "x*.foo", true}, // suffix wildcards should be allowed per RFC 6125/6.4.3
		/* 63 */ {[]string{"foo"}, MFalse, MAny, MTrue, MAny, MAny, nil, "foo", false},
		/* 64 */ {[]string{"foo"}, MTrue, MAny, MTrue, MAny, MAny, nil, "foo", true},

		// === Internal Glob Matching === //
		// Basic glob matching requirements
		/* 65 */ {[]string{"x*x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "xerox.foo", true},
		/* 66 */ {[]string{"x*x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "xylophone.files.pyrex.foo", true}, // globs can match across subdomains
		/* 67 */ {[]string{"x*x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "xercex.bar.foo", false}, // x.foo isn't matched
		/* 68 */ {[]string{"x*x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "bar.foo", false}, // x*x isn't matched.
		/* 69 */ {[]string{"x*x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "*.foo", false}, // unrelated wildcard
		/* 70 */ {[]string{"x*x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "*.x*x.foo", false}, // Does not conform to RFC 6125
		/* 71 */ {[]string{"x*x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "*.xyx.foo", false}, // Globs and Subdomains do not layer per docs.

		// Various requirements around x*x.foo wildcard matching.
		/* 72 */ {[]string{"x*x.foo"}, MFalse, MFalse, MAny, MAny, MAny, nil, "x*x.foo", false}, // base disabled, shouldn't match wildcard
		/* 73 */ {[]string{"x*x.foo"}, MFalse, MTrue, MAny, MAny, MTrue, nil, "x*x.foo", true}, // base disallowed, but globbing allowed and should match
		/* 74 */ {[]string{"x*x.foo"}, MTrue, MAny, MAny, MAny, MTrue, nil, "x*x.foo", true}, // base allowed, should match wildcard

		// Basic glob matching requirements with internal dots.
		/* 75 */ {[]string{"x.*.x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "xerox.foo", false}, // missing dots
		/* 76 */ {[]string{"x.*.x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "x.ero.x.foo", true},
		/* 77 */ {[]string{"x.*.x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "xylophone.files.pyrex.foo", false}, // missing dots
		/* 78 */ {[]string{"x.*.x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "x.ylophone.files.pyre.x.foo", true}, // globs can match across subdomains
		/* 79 */ {[]string{"x.*.x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "xercex.bar.foo", false}, // x.foo isn't matched
		/* 80 */ {[]string{"x.*.x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "bar.foo", false}, // x.*.x isn't matched.
		/* 81 */ {[]string{"x.*.x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "*.foo", false}, // unrelated wildcard
		/* 82 */ {[]string{"x.*.x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "*.x.*.x.foo", false}, // Does not conform to RFC 6125
		/* 83 */ {[]string{"x.*.x.foo"}, MAny, MTrue, MAny, MAny, MAny, nil, "*.x.y.x.foo", false}, // Globs and Subdomains do not layer per docs.

		// === Wildcard restriction testing === //
		/* 84 */ {[]string{"*.foo"}, MAny, MTrue, MFalse, MAny, MFalse, nil, "*.fud.bar.foo", false}, // glob domain allows wildcard of subdomains
		/* 85 */ {[]string{"*.foo"}, MAny, MTrue, MFalse, MAny, MFalse, nil, "*.bar.foo", false}, // glob domain allows wildcards of subdomains
		/* 86 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MFalse, nil, "*.fud.bar.foo", false},
		/* 87 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MFalse, nil, "*.bar.foo", false},
		/* 88 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MFalse, nil, "*.foo", false},
		/* 89 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MFalse, nil, "x*x.foo", false},
		/* 90 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MFalse, nil, "*x.foo", false},
		/* 91 */ {[]string{"foo"}, MAny, MAny, MTrue, MAny, MFalse, nil, "x*.foo", false},
		/* 92 */ {[]string{"x*x.foo"}, MTrue, MAny, MAny, MAny, MFalse, nil, "x*x.foo", false},
		/* 93 */ {[]string{"*.foo"}, MFalse, MFalse, MAny, MAny, MAny, nil, "*.foo", false}, // Bare and globs forbidden despite (potentially) allowing wildcards.
		/* 94 */ {[]string{"x.*.x.foo"}, MAny, MAny, MAny, MAny, MAny, nil, "x.*.x.foo", false}, // Does not conform to RFC 6125

		// === CN validation allowances === //
		/*  95 */ {[]string{"foo"}, MAny, MAny, MAny, MAny, MAny, []string{"disabled"}, "*.fud.bar.foo", true},
		/*  96 */ {[]string{"foo"}, MAny, MAny, MAny, MAny, MAny, []string{"disabled"}, "*.fud.*.foo", true},
		/*  97 */ {[]string{"foo"}, MAny, MAny, MAny, MAny, MAny, []string{"disabled"}, "*.bar.*.bar", true},
		/*  98 */ {[]string{"foo"}, MAny, MAny, MAny, MAny, MAny, []string{"disabled"}, "foo@foo", true},
		/*  99 */ {[]string{"foo"}, MAny, MAny, MAny, MAny, MAny, []string{"disabled"}, "foo@foo@foo", true},
		/* 100 */ {[]string{"foo"}, MAny, MAny, MAny, MAny, MAny, []string{"disabled"}, "bar@bar@bar", true},
		/* 101 */ {[]string{"foo"}, MTrue, MTrue, MTrue, MTrue, MTrue, []string{"email"}, "bar@bar@bar", false},
		/* 102 */ {[]string{"foo"}, MTrue, MTrue, MTrue, MTrue, MTrue, []string{"email"}, "bar@bar", false},
		/* 103 */ {[]string{"foo"}, MTrue, MTrue, MTrue, MTrue, MTrue, []string{"email"}, "bar@foo", true},
		/* 104 */ {[]string{"foo"}, MTrue, MTrue, MTrue, MTrue, MTrue, []string{"hostname"}, "bar@foo", false},
		/* 105 */ {[]string{"foo"}, MTrue, MTrue, MTrue, MTrue, MTrue, []string{"hostname"}, "bar@bar", false},
		/* 106 */ {[]string{"foo"}, MTrue, MTrue, MTrue, MTrue, MTrue, []string{"hostname"}, "bar.foo", true},
		/* 107 */ {[]string{"foo"}, MTrue, MTrue, MTrue, MTrue, MTrue, []string{"hostname"}, "bar.bar", false},
		/* 108 */ {[]string{"foo"}, MTrue, MTrue, MTrue, MTrue, MTrue, []string{"email"}, "bar.foo", false},
		/* 109 */ {[]string{"foo"}, MTrue, MTrue, MTrue, MTrue, MTrue, []string{"email"}, "bar.bar", false},
	}

	if len(testCases) != 110 {
		t.Fatalf("misnumbered test case entries will make it hard to find bugs: %v", len(testCases))
	}

	b, s := CreateBackendWithStorage(t)

	// We need a RSA key so all signature sizes are valid with it.
	resp, err := CBWrite(b, s, "root/generate/exported", map[string]interface{}{
		"common_name": "myvault.com",
		"ttl":         "128h",
		"key_type":    "rsa",
		"key_bits":    2048,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected ca info")
	}

	tested := 0
	for index, test := range testCases {
		tested += RoleIssuanceRegressionHelper(t, b, s, index, test)
	}

	t.Logf("Issuance regression expanded matrix test scenarios: %d", tested)
}

type KeySizeRegression struct {
	// Values reused for both Role and CA configuration.
	RoleKeyType string
	RoleKeyBits []int

	// Signature Bits presently is only specified on the role.
	RoleSignatureBits []int
	RoleUsePSS        bool

	// These are tuples; must be of the same length.
	TestKeyTypes []string
	TestKeyBits  []int

	// All of the above key types/sizes must pass or fail together.
	ExpectError bool
}

func (k KeySizeRegression) KeyTypeValues() []string {
	if k.RoleKeyType == "any" {
		return []string{"rsa", "ec", "ed25519"}
	}

	return []string{k.RoleKeyType}
}

func RoleKeySizeRegressionHelper(t *testing.T, b *backend, s logical.Storage, index int, test KeySizeRegression) int {
	tested := 0

	for _, caKeyType := range test.KeyTypeValues() {
		for _, caKeyBits := range test.RoleKeyBits {
			// Generate a new CA key.
			resp, err := CBWrite(b, s, "root/generate/exported", map[string]interface{}{
				"common_name": "myvault.com",
				"ttl":         "128h",
				"key_type":    caKeyType,
				"key_bits":    caKeyBits,
			})
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected ca info")
			}

			for _, roleKeyBits := range test.RoleKeyBits {
				for _, roleSignatureBits := range test.RoleSignatureBits {
					role := fmt.Sprintf("key-size-regression-%d-keytype-%v-keybits-%d-signature-bits-%d", index, test.RoleKeyType, roleKeyBits, roleSignatureBits)
					_, err := CBWrite(b, s, "roles/"+role, map[string]interface{}{
						"key_type":       test.RoleKeyType,
						"key_bits":       roleKeyBits,
						"signature_bits": roleSignatureBits,
						"use_pss":        test.RoleUsePSS,
					})
					if err != nil {
						t.Fatal(err)
					}

					for index, keyType := range test.TestKeyTypes {
						keyBits := test.TestKeyBits[index]

						_, _, csrPem := generateCSR(t, &x509.CertificateRequest{
							Subject: pkix.Name{
								CommonName: "localhost",
							},
						}, keyType, keyBits)

						resp, err = CBWrite(b, s, "sign/"+role, map[string]interface{}{
							"common_name": "localhost",
							"csr":         csrPem,
						})

						haveErr := err != nil || resp == nil

						if haveErr != test.ExpectError {
							t.Fatalf("key size regression test [%d] failed: haveErr: %v, expectErr: %v, err: %v, resp: %v, test case: %v, caKeyType: %v, caKeyBits: %v, role: %v, keyType: %v, keyBits: %v", index, haveErr, test.ExpectError, err, resp, test, caKeyType, caKeyBits, role, keyType, keyBits)
						}

						if resp != nil && test.RoleUsePSS && caKeyType == "rsa" {
							leafCert := parseCert(t, resp.Data["certificate"].(string))
							switch leafCert.SignatureAlgorithm {
							case x509.SHA256WithRSAPSS, x509.SHA384WithRSAPSS, x509.SHA512WithRSAPSS:
							default:
								t.Fatalf("key size regression test [%d] failed on role %v: unexpected signature algorithm; expected RSA-type CA to sign a leaf cert with PSS algorithm; got %v", index, role, leafCert.SignatureAlgorithm.String())
							}
						}

						tested += 1
					}
				}
			}

			_, err = CBDelete(b, s, "root")
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	return tested
}

func TestBackend_Roles_KeySizeRegression(t *testing.T) {
	t.Parallel()
	// Regression testing of role's issuance policy.
	testCases := []KeySizeRegression{
		// RSA with default parameters should fail to issue smaller RSA keys
		// and any size ECDSA/Ed25519 keys.
		/*  0 */ {"rsa", []int{0, 2048}, []int{0, 256, 384, 512}, false, []string{"rsa", "ec", "ec", "ec", "ec", "ed25519"}, []int{1024, 224, 256, 384, 521, 0}, true},
		// But it should work to issue larger RSA keys.
		/*  1 */ {"rsa", []int{0, 2048}, []int{0, 256, 384, 512}, false, []string{"rsa", "rsa"}, []int{2048, 3072}, false},

		// EC with default parameters should fail to issue smaller EC keys
		// and any size RSA/Ed25519 keys.
		/*  2 */ {"ec", []int{0}, []int{0}, false, []string{"rsa", "ec", "ed25519"}, []int{2048, 224, 0}, true},
		// But it should work to issue larger EC keys. Note that we should be
		// independent of signature bits as that's computed from the issuer
		// type (for EC based issuers).
		/*  3 */ {"ec", []int{224}, []int{0, 256, 384, 521}, false, []string{"ec", "ec", "ec", "ec"}, []int{224, 256, 384, 521}, false},
		/*  4 */ {"ec", []int{0, 256}, []int{0, 256, 384, 521}, false, []string{"ec", "ec", "ec"}, []int{256, 384, 521}, false},
		/*  5 */ {"ec", []int{384}, []int{0, 256, 384, 521}, false, []string{"ec", "ec"}, []int{384, 521}, false},
		/*  6 */ {"ec", []int{521}, []int{0, 256, 384, 512}, false, []string{"ec"}, []int{521}, false},

		// Ed25519 should reject RSA and EC keys.
		/*  7 */ {"ed25519", []int{0}, []int{0}, false, []string{"rsa", "ec", "ec"}, []int{2048, 256, 521}, true},
		// But it should work to issue Ed25519 keys.
		/*  8 */ {"ed25519", []int{0}, []int{0}, false, []string{"ed25519"}, []int{0}, false},

		// Any key type should reject insecure RSA key sizes.
		/*  9 */ {"any", []int{0}, []int{0, 256, 384, 512}, false, []string{"rsa", "rsa"}, []int{512, 1024}, true},
		// But work for everything else.
		/* 10 */ {"any", []int{0}, []int{0, 256, 384, 512}, false, []string{"rsa", "rsa", "ec", "ec", "ec", "ec", "ed25519"}, []int{2048, 3072, 224, 256, 384, 521, 0}, false},

		// RSA with larger than default key size should reject smaller ones.
		/* 11 */ {"rsa", []int{3072}, []int{0, 256, 384, 512}, false, []string{"rsa"}, []int{2048}, true},

		// We should be able to sign with PSS with any CA key type.
		/* 12 */ {"rsa", []int{0}, []int{0, 256, 384, 512}, true, []string{"rsa"}, []int{2048}, false},
		/* 13 */ {"ec", []int{0}, []int{0}, true, []string{"ec"}, []int{256}, false},
		/* 14 */ {"ed25519", []int{0}, []int{0}, true, []string{"ed25519"}, []int{0}, false},
	}

	if len(testCases) != 15 {
		t.Fatalf("misnumbered test case entries will make it hard to find bugs: %v", len(testCases))
	}

	b, s := CreateBackendWithStorage(t)

	tested := 0
	for index, test := range testCases {
		tested += RoleKeySizeRegressionHelper(t, b, s, index, test)
	}

	t.Logf("Key size regression expanded matrix test scenarios: %d", tested)
}

func TestRootWithExistingKey(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	var err error

	// Fail requests if type is existing, and we specify the key_type param
	_, err = CBWrite(b, s, "root/generate/existing", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_type":    "rsa",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "key_type nor key_bits arguments can be set in this mode")

	// Fail requests if type is existing, and we specify the key_bits param
	_, err = CBWrite(b, s, "root/generate/existing", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_bits":    "2048",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "key_type nor key_bits arguments can be set in this mode")

	// Fail if the specified key does not exist.
	_, err = CBWrite(b, s, "issuers/generate/root/existing", map[string]interface{}{
		"common_name": "root myvault.com",
		"issuer_name": "my-issuer1",
		"key_ref":     "my-key1",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to find PKI key for reference: my-key1")

	// Fail if the specified key name is default.
	_, err = CBWrite(b, s, "issuers/generate/root/internal", map[string]interface{}{
		"common_name": "root myvault.com",
		"issuer_name": "my-issuer1",
		"key_name":    "Default",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "reserved keyword 'default' can not be used as key name")

	// Fail if the specified issuer name is default.
	_, err = CBWrite(b, s, "issuers/generate/root/internal", map[string]interface{}{
		"common_name": "root myvault.com",
		"issuer_name": "DEFAULT",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "reserved keyword 'default' can not be used as issuer name")

	// Create the first CA
	resp, err := CBWrite(b, s, "issuers/generate/root/internal", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_type":    "rsa",
		"issuer_name": "my-issuer1",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("issuers/generate/root/internal"), logical.UpdateOperation), resp, true)
	require.NoError(t, err)
	require.NotNil(t, resp.Data["certificate"])
	myIssuerId1 := resp.Data["issuer_id"]
	myKeyId1 := resp.Data["key_id"]
	require.NotEmpty(t, myIssuerId1)
	require.NotEmpty(t, myKeyId1)

	// Fetch the parsed CRL; it should be empty as we've not revoked anything
	parsedCrl := getParsedCrlFromBackend(t, b, s, "issuer/my-issuer1/crl/der")
	require.Equal(t, len(parsedCrl.TBSCertList.RevokedCertificates), 0, "should have no revoked certificates")

	// Fail if the specified issuer name is re-used.
	_, err = CBWrite(b, s, "issuers/generate/root/internal", map[string]interface{}{
		"common_name": "root myvault.com",
		"issuer_name": "my-issuer1",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "issuer name already in use")

	// Create the second CA
	resp, err = CBWrite(b, s, "issuers/generate/root/internal", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_type":    "rsa",
		"issuer_name": "my-issuer2",
		"key_name":    "root-key2",
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Data["certificate"])
	myIssuerId2 := resp.Data["issuer_id"]
	myKeyId2 := resp.Data["key_id"]
	require.NotEmpty(t, myIssuerId2)
	require.NotEmpty(t, myKeyId2)

	// Fetch the parsed CRL; it should be empty as we've not revoked anything
	parsedCrl = getParsedCrlFromBackend(t, b, s, "issuer/my-issuer2/crl/der")
	require.Equal(t, len(parsedCrl.TBSCertList.RevokedCertificates), 0, "should have no revoked certificates")

	// Fail if the specified key name is re-used.
	_, err = CBWrite(b, s, "issuers/generate/root/internal", map[string]interface{}{
		"common_name": "root myvault.com",
		"issuer_name": "my-issuer3",
		"key_name":    "root-key2",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "key name already in use")

	// Create a third CA re-using key from CA 1
	resp, err = CBWrite(b, s, "issuers/generate/root/existing", map[string]interface{}{
		"common_name": "root myvault.com",
		"issuer_name": "my-issuer3",
		"key_ref":     myKeyId1,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Data["certificate"])
	myIssuerId3 := resp.Data["issuer_id"]
	myKeyId3 := resp.Data["key_id"]
	require.NotEmpty(t, myIssuerId3)
	require.NotEmpty(t, myKeyId3)

	// Fetch the parsed CRL; it should be empty as we've not revoking anything.
	parsedCrl = getParsedCrlFromBackend(t, b, s, "issuer/my-issuer3/crl/der")
	require.Equal(t, len(parsedCrl.TBSCertList.RevokedCertificates), 0, "should have no revoked certificates")
	// Signatures should be the same since this is just a reissued cert. We
	// use signature as a proxy for "these two CRLs are equal".
	firstCrl := getParsedCrlFromBackend(t, b, s, "issuer/my-issuer1/crl/der")
	require.Equal(t, parsedCrl.SignatureValue, firstCrl.SignatureValue)

	require.NotEqual(t, myIssuerId1, myIssuerId2)
	require.NotEqual(t, myIssuerId1, myIssuerId3)
	require.NotEqual(t, myKeyId1, myKeyId2)
	require.Equal(t, myKeyId1, myKeyId3)

	resp, err = CBList(b, s, "issuers")
	require.NoError(t, err)
	require.Equal(t, 3, len(resp.Data["keys"].([]string)))
	require.Contains(t, resp.Data["keys"], string(myIssuerId1.(issuing.IssuerID)))
	require.Contains(t, resp.Data["keys"], string(myIssuerId2.(issuing.IssuerID)))
	require.Contains(t, resp.Data["keys"], string(myIssuerId3.(issuing.IssuerID)))
}

func TestIntermediateWithExistingKey(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	var err error

	// Fail requests if type is existing, and we specify the key_type param
	_, err = CBWrite(b, s, "intermediate/generate/existing", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_type":    "rsa",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "key_type nor key_bits arguments can be set in this mode")

	// Fail requests if type is existing, and we specify the key_bits param
	_, err = CBWrite(b, s, "intermediate/generate/existing", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_bits":    "2048",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "key_type nor key_bits arguments can be set in this mode")

	// Fail if the specified key does not exist.
	_, err = CBWrite(b, s, "issuers/generate/intermediate/existing", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_ref":     "my-key1",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to find PKI key for reference: my-key1")

	// Create the first intermediate CA
	resp, err := CBWrite(b, s, "issuers/generate/intermediate/internal", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_type":    "rsa",
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("issuers/generate/intermediate/internal"), logical.UpdateOperation), resp, true)
	require.NoError(t, err)
	// csr1 := resp.Data["csr"]
	myKeyId1 := resp.Data["key_id"]
	require.NotEmpty(t, myKeyId1)

	// Create the second intermediate CA
	resp, err = CBWrite(b, s, "issuers/generate/intermediate/internal", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_type":    "rsa",
		"key_name":    "interkey1",
	})
	require.NoError(t, err)
	// csr2 := resp.Data["csr"]
	myKeyId2 := resp.Data["key_id"]
	require.NotEmpty(t, myKeyId2)

	// Create a third intermediate CA re-using key from intermediate CA 1
	resp, err = CBWrite(b, s, "issuers/generate/intermediate/existing", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_ref":     myKeyId1,
	})
	require.NoError(t, err)
	// csr3 := resp.Data["csr"]
	myKeyId3 := resp.Data["key_id"]
	require.NotEmpty(t, myKeyId3)

	require.NotEqual(t, myKeyId1, myKeyId2)
	require.Equal(t, myKeyId1, myKeyId3, "our new ca did not seem to reuse the key as we expected.")
}

func TestIssuanceTTLs(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root example.com",
		"issuer_name": "root",
		"ttl":         "10s",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	rootCert := parseCert(t, resp.Data["certificate"].(string))

	_, err = CBWrite(b, s, "roles/local-testing", map[string]interface{}{
		"allow_any_name":    true,
		"enforce_hostnames": false,
		"key_type":          "ec",
	})
	require.NoError(t, err)

	_, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing",
		"ttl":         "1s",
	})
	require.NoError(t, err, "expected issuance to succeed due to shorter ttl than cert ttl")

	_, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing",
	})
	require.Error(t, err, "expected issuance to fail due to longer default ttl than cert ttl")

	resp, err = CBPatch(b, s, "issuer/root", map[string]interface{}{
		"leaf_not_after_behavior": "permit",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.Equal(t, resp.Data["leaf_not_after_behavior"], "permit")

	_, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing",
	})
	require.NoError(t, err, "expected issuance to succeed due to permitted longer TTL")

	resp, err = CBWrite(b, s, "issuer/root", map[string]interface{}{
		"issuer_name":             "root",
		"leaf_not_after_behavior": "truncate",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.Equal(t, resp.Data["leaf_not_after_behavior"], "truncate")

	_, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing",
	})
	require.NoError(t, err, "expected issuance to succeed due to truncated ttl")

	// Sleep until the parent cert expires and the clock rolls over
	// to the next second.
	time.Sleep(time.Until(rootCert.NotAfter) + (1500 * time.Millisecond))

	resp, err = CBWrite(b, s, "issuer/root", map[string]interface{}{
		"issuer_name":             "root",
		"leaf_not_after_behavior": "err",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Even 1s ttl should now fail.
	_, err = CBWrite(b, s, "issue/local-testing", map[string]interface{}{
		"common_name": "testing",
		"ttl":         "1s",
	})
	require.Error(t, err, "expected issuance to fail due to longer default ttl than cert ttl")
}

func TestSealWrappedStorageConfigured(t *testing.T) {
	t.Parallel()
	b, _ := CreateBackendWithStorage(t)
	wrappedEntries := b.Backend.PathsSpecial.SealWrapStorage

	// Make sure our legacy bundle is within the list
	// NOTE: do not convert these test values to constants, we should always have these paths within seal wrap config
	require.Contains(t, wrappedEntries, "config/ca_bundle", "Legacy bundle missing from seal wrap")
	// The trailing / is important as it treats the entire folder requiring seal wrapping, not just config/key
	require.Contains(t, wrappedEntries, "config/key/", "key prefix with trailing / missing from seal wrap.")
}

func TestBackend_ConfigCA_WithECParams(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Generated key with OpenSSL:
	// $ openssl ecparam -out p256.key -name prime256v1 -genkey
	//
	// Regression test for https://github.com/hashicorp/vault/issues/16667
	resp, err := CBWrite(b, s, "config/ca", map[string]interface{}{
		"pem_bundle": `
-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEINzXthCZdhyV7+wIEBl/ty+ctNsUS99ykTeax6EbYZtvoAoGCCqGSM49
AwEHoUQDQgAE57NX8bR/nDoW8yRgLswoXBQcjHrdyfuHS0gPwki6BNnfunUzryVb
8f22/JWj6fsEF6AOADZlrswKIbR2Es9e/w==
-----END EC PRIVATE KEY-----
		`,
	})
	require.NoError(t, err)
	require.NotNil(t, resp, "expected ca info")
	importedKeys := resp.Data["imported_keys"].([]string)
	importedIssuers := resp.Data["imported_issuers"].([]string)

	require.Equal(t, len(importedKeys), 1)
	require.Equal(t, len(importedIssuers), 0)
}

func TestPerIssuerAIA(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Generating a root without anything should not have AIAs.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root example.com",
		"issuer_name": "root",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	rootCert := parseCert(t, resp.Data["certificate"].(string))
	require.Empty(t, rootCert.OCSPServer)
	require.Empty(t, rootCert.IssuingCertificateURL)
	require.Empty(t, rootCert.CRLDistributionPoints)

	// Set some local URLs on the issuer.
	resp, err = CBWrite(b, s, "issuer/default", map[string]interface{}{
		"issuing_certificates": []string{"https://google.com"},
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("issuer/default"), logical.UpdateOperation), resp, true)

	require.NoError(t, err)

	_, err = CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allow_any_name": true,
		"ttl":            "85s",
		"key_type":       "ec",
	})
	require.NoError(t, err)

	// Issue something with this re-configured issuer.
	resp, err = CBWrite(b, s, "issuer/default/issue/testing", map[string]interface{}{
		"common_name": "localhost.com",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	leafCert := parseCert(t, resp.Data["certificate"].(string))
	require.Empty(t, leafCert.OCSPServer)
	require.Equal(t, leafCert.IssuingCertificateURL, []string{"https://google.com"})
	require.Empty(t, leafCert.CRLDistributionPoints)

	// Set global URLs and ensure they don't appear on this issuer's leaf.
	_, err = CBWrite(b, s, "config/urls", map[string]interface{}{
		"issuing_certificates":    []string{"https://example.com/ca", "https://backup.example.com/ca"},
		"crl_distribution_points": []string{"https://example.com/crl", "https://backup.example.com/crl"},
		"ocsp_servers":            []string{"https://example.com/ocsp", "https://backup.example.com/ocsp"},
	})
	require.NoError(t, err)
	resp, err = CBWrite(b, s, "issuer/default/issue/testing", map[string]interface{}{
		"common_name": "localhost.com",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	leafCert = parseCert(t, resp.Data["certificate"].(string))
	require.Empty(t, leafCert.OCSPServer)
	require.Equal(t, leafCert.IssuingCertificateURL, []string{"https://google.com"})
	require.Empty(t, leafCert.CRLDistributionPoints)

	// Now come back and remove the local modifications and ensure we get
	// the defaults again.
	_, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
		"issuing_certificates": []string{},
	})
	require.NoError(t, err)
	resp, err = CBWrite(b, s, "issuer/default/issue/testing", map[string]interface{}{
		"common_name": "localhost.com",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	leafCert = parseCert(t, resp.Data["certificate"].(string))
	require.Equal(t, leafCert.IssuingCertificateURL, []string{"https://example.com/ca", "https://backup.example.com/ca"})
	require.Equal(t, leafCert.OCSPServer, []string{"https://example.com/ocsp", "https://backup.example.com/ocsp"})
	require.Equal(t, leafCert.CRLDistributionPoints, []string{"https://example.com/crl", "https://backup.example.com/crl"})

	// Validate that we can set an issuer name and remove it.
	_, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
		"issuer_name": "my-issuer",
	})
	require.NoError(t, err)
	_, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
		"issuer_name": "",
	})
	require.NoError(t, err)
}

func TestIssuersWithoutCRLBits(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Importing a root without CRL signing bits should work fine.
	customBundleWithoutCRLBits := `
-----BEGIN CERTIFICATE-----
MIIDGTCCAgGgAwIBAgIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDDAhyb290
LW5ldzAeFw0yMjA4MjQxMjEzNTVaFw0yMzA5MDMxMjEzNTVaMBMxETAPBgNVBAMM
CHJvb3QtbmV3MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAojTA/Mx7
LVW/Zgn/N4BqZbaF82MrTIBFug3ob7mqycNRlWp4/PH8v37+jYn8e691HUsKjden
rDTrO06kiQKiJinAzmlLJvgcazE3aXoh7wSzVG9lFHYvljEmVj+yDbkeaqaCktup
skuNjxCoN9BLmKzZIwVCHn92ZHlhN6LI7CNaU3SDJdu7VftWF9Ugzt9FIvI+6Gcn
/WNE9FWvZ9o7035rZ+1vvTn7/tgxrj2k3XvD51Kq4tsSbqjnSf3QieXT6E6uvtUE
TbPp3xjBElgBCKmeogR1l28rs1aujqqwzZ0B/zOeF8ptaH0aZOIBsVDJR8yTwHzq
s34hNdNfKLHzOwIDAQABo3gwdjAdBgNVHQ4EFgQUF4djNmx+1+uJINhZ82pN+7jz
H8EwHwYDVR0jBBgwFoAUF4djNmx+1+uJINhZ82pN+7jzH8EwDwYDVR0TAQH/BAUw
AwEB/zAOBgNVHQ8BAf8EBAMCAoQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDQYJKoZI
hvcNAQELBQADggEBAICQovBz4KLWlLmXeZ2Vf6WfQYyGNgGyJa10XNXtWQ5dM2NU
OLAit4x1c2dz+aFocc8ZsX/ikYi/bruT2rsGWqMAGC4at3U4GuaYGO5a6XzMKIDC
nxIlbiO+Pn6Xum7fAqUri7+ZNf/Cygmc5sByi3MAAIkszeObUDZFTJL7gEOuXIMT
rKIXCINq/U+qc7m9AQ8vKhF1Ddj+dLGLzNQ5j3cKfilPs/wRaYqbMQvnmarX+5Cs
k1UL6kWSQsiP3+UWaBlcWkmD6oZ3fIG7c0aMxf7RISq1eTAM9XjH3vMxWQJlS5q3
2weJ2LYoPe/DwX5CijR0IezapBCrin1BscJMLFQ=
-----END CERTIFICATE-----
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCiNMD8zHstVb9m
Cf83gGpltoXzYytMgEW6DehvuarJw1GVanj88fy/fv6Nifx7r3UdSwqN16esNOs7
TqSJAqImKcDOaUsm+BxrMTdpeiHvBLNUb2UUdi+WMSZWP7INuR5qpoKS26myS42P
EKg30EuYrNkjBUIef3ZkeWE3osjsI1pTdIMl27tV+1YX1SDO30Ui8j7oZyf9Y0T0
Va9n2jvTfmtn7W+9Ofv+2DGuPaTde8PnUqri2xJuqOdJ/dCJ5dPoTq6+1QRNs+nf
GMESWAEIqZ6iBHWXbyuzVq6OqrDNnQH/M54Xym1ofRpk4gGxUMlHzJPAfOqzfiE1
018osfM7AgMBAAECggEAAVd6kZZaN69IZITIc1vHRYa2rlZpKS2JP7c8Vd3Z/4Fz
ZZvnJ7LgVAmUYg5WPZ2sOqBNLfKVN/oke5Q0dALgdxYl7dWQIhPjHeRFbZFtjqEV
OXZGBniamMO/HSKGWGrqFf7BM/H7AhClUwQgjnzVSz+B+LJJidM+SVys3n1xuDmC
EP+iOda+bAHqHv/7oCELQKhLmCvPc9v2fDy+180ttdo8EHuxwVnKiyR/ryKFhSyx
K1wgAPQ9jO+V+GESL90rqpX/r501REsIOOpm4orueelHTD4+dnHxvUPqJ++9aYGX
79qBNPPUhxrQI1yoHxwW0cTxW5EqkZ9bT2lSd5rjcQKBgQDNyPBpidkHPrYemQDT
RldtS6FiW/jc1It/CRbjU4A6Gi7s3Cda43pEUObKNLeXMyLQaMf4GbDPDX+eh7B8
RkUq0Q/N0H4bn1hbxYSUdgv0j/6czpMo6rLcJHGwOTSpHGsNsxSLL7xlpgzuzqrG
FzEgjMA1aD3w8B9+/77AoSLoMQKBgQDJyYMw82+euLYRbR5Wc/SbrWfh2n1Mr2BG
pp1ZNYorXE5CL4ScdLcgH1q/b8r5XGwmhMcpeA+geAAaKmk1CGG+gPLoq20c9Q1Y
Ykq9tUVJasIkelvbb/SPxyjkJdBwylzcPP14IJBsqQM0be+yVqLJJVHSaoKhXZcl
IW2xgCpjKwKBgFpeX5U5P+F6nKebMU2WmlYY3GpBUWxIummzKCX0SV86mFjT5UR4
mPzfOjqaI/V2M1eqbAZ74bVLjDumAs7QXReMb5BGetrOgxLqDmrT3DQt9/YMkXtq
ddlO984XkRSisjB18BOfhvBsl0lX4I7VKHHO3amWeX0RNgOjc7VMDfRBAoGAWAQH
r1BfvZHACLXZ58fISCdJCqCsysgsbGS8eW77B5LJp+DmLQBT6DUE9j+i/0Wq/ton
rRTrbAkrsj4RicpQKDJCwe4UN+9DlOu6wijRQgbJC/Q7IOoieJxcX7eGxcve2UnZ
HY7GsD7AYRwa02UquCYJHIjM1enmxZFhMW1AD+UCgYEAm4jdNz5e4QjA4AkNF+cB
ZenrAZ0q3NbTyiSsJEAtRe/c5fNFpmXo3mqgCannarREQYYDF0+jpSoTUY8XAc4q
wL7EZNzwxITLqBnnHQbdLdAvYxB43kvWTy+JRK8qY9LAMCCFeDoYwXkWV4Wkx/b0
TgM7RZnmEjNdeaa4M52o7VY=
-----END PRIVATE KEY-----
	`
	resp, err := CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": customBundleWithoutCRLBits,
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("issuers/import/bundle"), logical.UpdateOperation), resp, true)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data)
	require.NotEmpty(t, resp.Data["imported_issuers"])
	require.NotEmpty(t, resp.Data["imported_keys"])
	require.NotEmpty(t, resp.Data["mapping"])

	// Shouldn't have crl-signing on the newly imported issuer's usage.
	resp, err = CBRead(b, s, "issuer/default")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data)
	require.NotEmpty(t, resp.Data["usage"])
	require.NotContains(t, resp.Data["usage"], "crl-signing")

	// Modifying to set CRL should fail.
	resp, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
		"usage": "issuing-certificates,crl-signing",
	})
	require.Error(t, err)
	require.True(t, resp.IsError())

	// Modifying to set issuing-certificates and ocsp-signing should succeed.
	resp, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
		"usage": "issuing-certificates,ocsp-signing",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Data)
	require.NotEmpty(t, resp.Data["usage"])
	require.NotContains(t, resp.Data["usage"], "crl-signing")
}

func TestBackend_IfModifiedSinceHeaders(t *testing.T) {
	t.Parallel()
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc:             vaulthttp.Handler,
		RequestResponseCallback: schema.ResponseValidatingCallback(t),
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// Mount PKI.
	err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "60h",
			// Required to allow the header to be passed through.
			PassthroughRequestHeaders: []string{"if-modified-since"},
			AllowedResponseHeaders:    []string{"Last-Modified"},
		},
	})
	require.NoError(t, err)

	// Get a time before CA generation. Subtract two seconds to ensure
	// the value in the seconds field is different than the time the CA
	// is actually generated at.
	beforeOldCAGeneration := time.Now().Add(-2 * time.Second)

	// Generate an internal CA. This one is the default.
	resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "Root X1",
		"key_type":    "ec",
		"issuer_name": "old-root",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["certificate"])

	// CA is generated, but give a grace window.
	afterOldCAGeneration := time.Now().Add(2 * time.Second)

	// When you _save_ headers, client returns a copy. But when you go to
	// reset them, it doesn't create a new copy (and instead directly
	// assigns). This means we have to continually refresh our view of the
	// last headers, otherwise the headers added after the last set operation
	// leak into this copy... Yuck!
	lastHeaders := client.Headers()
	for _, path := range []string{"pki/cert/ca", "pki/cert/crl", "pki/issuer/default/json", "pki/issuer/old-root/json", "pki/issuer/old-root/crl", "pki/cert/delta-crl", "pki/issuer/old-root/crl/delta"} {
		t.Logf("path: %v", path)
		field := "certificate"
		if strings.HasPrefix(path, "pki/issuer") && strings.Contains(path, "/crl") {
			field = "crl"
		}

		// Reading the CA should work, without a header.
		resp, err := client.Logical().Read(path)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data[field])

		// Ensure that the CA is returned correctly if we give it the old time.
		client.AddHeader("If-Modified-Since", beforeOldCAGeneration.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data[field])
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()

		// Ensure that the CA is elided if we give it the present time (plus a
		// grace window).
		client.AddHeader("If-Modified-Since", afterOldCAGeneration.Format(time.RFC1123))
		t.Logf("headers: %v", client.Headers())
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.Nil(t, resp)
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()
	}

	// Wait three seconds. This ensures we have adequate grace period
	// to distinguish the two cases, even with grace periods.
	time.Sleep(3 * time.Second)

	// Generating a second root. This one isn't the default.
	beforeNewCAGeneration := time.Now().Add(-2 * time.Second)

	// Generate an internal CA. This one is the default.
	_, err = client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "Root X1",
		"key_type":    "ec",
		"issuer_name": "new-root",
	})
	require.NoError(t, err)

	// As above.
	afterNewCAGeneration := time.Now().Add(2 * time.Second)

	// New root isn't the default, so it has fewer paths.
	for _, path := range []string{"pki/issuer/new-root/json", "pki/issuer/new-root/crl", "pki/issuer/new-root/crl/delta"} {
		t.Logf("path: %v", path)
		field := "certificate"
		if strings.HasPrefix(path, "pki/issuer") && strings.Contains(path, "/crl") {
			field = "crl"
		}

		// Reading the CA should work, without a header.
		resp, err := client.Logical().Read(path)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data[field])

		// Ensure that the CA is returned correctly if we give it the old time.
		client.AddHeader("If-Modified-Since", beforeNewCAGeneration.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data[field])
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()

		// Ensure that the CA is elided if we give it the present time (plus a
		// grace window).
		client.AddHeader("If-Modified-Since", afterNewCAGeneration.Format(time.RFC1123))
		t.Logf("headers: %v", client.Headers())
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.Nil(t, resp)
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()
	}

	// Wait three seconds. This ensures we have adequate grace period
	// to distinguish the two cases, even with grace periods.
	time.Sleep(3 * time.Second)

	// Now swap the default issuers around.
	_, err = client.Logical().Write("pki/config/issuers", map[string]interface{}{
		"default": "new-root",
	})
	require.NoError(t, err)

	// Reading both with the last modified date should return new values.
	for _, path := range []string{"pki/cert/ca", "pki/cert/crl", "pki/issuer/default/json", "pki/issuer/old-root/json", "pki/issuer/new-root/json", "pki/issuer/old-root/crl", "pki/issuer/new-root/crl", "pki/cert/delta-crl", "pki/issuer/old-root/crl/delta", "pki/issuer/new-root/crl/delta"} {
		t.Logf("path: %v", path)
		field := "certificate"
		if strings.HasPrefix(path, "pki/issuer") && strings.Contains(path, "/crl") {
			field = "crl"
		}

		// Ensure that the CA is returned correctly if we give it the old time.
		client.AddHeader("If-Modified-Since", afterOldCAGeneration.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data[field])
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()

		// Ensure that the CA is returned correctly if we give it the old time.
		client.AddHeader("If-Modified-Since", afterNewCAGeneration.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data[field])
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()
	}

	// Wait for things to settle, record the present time, and wait for the
	// clock to definitely tick over again.
	time.Sleep(2 * time.Second)
	preRevocationTimestamp := time.Now()
	time.Sleep(2 * time.Second)

	// The above tests should say everything is cached.
	for _, path := range []string{"pki/cert/ca", "pki/cert/crl", "pki/issuer/default/json", "pki/issuer/old-root/json", "pki/issuer/new-root/json", "pki/issuer/old-root/crl", "pki/issuer/new-root/crl", "pki/cert/delta-crl", "pki/issuer/old-root/crl/delta", "pki/issuer/new-root/crl/delta"} {
		t.Logf("path: %v", path)

		// Ensure that the CA is returned correctly if we give it the new time.
		client.AddHeader("If-Modified-Since", preRevocationTimestamp.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.Nil(t, resp)
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()
	}

	// We could generate some leaves and verify the revocation updates the
	// CRL. But, revoking the issuer behaves the same, so let's do that
	// instead.
	_, err = client.Logical().Write("pki/issuer/old-root/revoke", map[string]interface{}{})
	require.NoError(t, err)

	// CA should still be valid.
	for _, path := range []string{"pki/cert/ca", "pki/issuer/default/json", "pki/issuer/old-root/json", "pki/issuer/new-root/json"} {
		t.Logf("path: %v", path)

		// Ensure that the CA is returned correctly if we give it the old time.
		client.AddHeader("If-Modified-Since", preRevocationTimestamp.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.Nil(t, resp)
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()
	}

	// CRL should be invalidated
	for _, path := range []string{"pki/cert/crl", "pki/issuer/old-root/crl", "pki/issuer/new-root/crl", "pki/cert/delta-crl", "pki/issuer/old-root/crl/delta", "pki/issuer/new-root/crl/delta"} {
		t.Logf("path: %v", path)
		field := "certificate"
		if strings.HasPrefix(path, "pki/issuer") && strings.Contains(path, "/crl") {
			field = "crl"
		}

		client.AddHeader("If-Modified-Since", preRevocationTimestamp.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data[field])
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()
	}

	// If we send some time in the future, everything should be cached again!
	futureTime := time.Now().Add(30 * time.Second)
	for _, path := range []string{"pki/cert/ca", "pki/cert/crl", "pki/issuer/default/json", "pki/issuer/old-root/json", "pki/issuer/new-root/json", "pki/issuer/old-root/crl", "pki/issuer/new-root/crl", "pki/cert/delta-crl", "pki/issuer/old-root/crl/delta", "pki/issuer/new-root/crl/delta"} {
		t.Logf("path: %v", path)

		// Ensure that the CA is returned correctly if we give it the new time.
		client.AddHeader("If-Modified-Since", futureTime.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.Nil(t, resp)
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()
	}

	beforeThreeWaySwap := time.Now().Add(-2 * time.Second)

	// Now, do a three-way swap of names (old->tmp; new->old; tmp->new). This
	// should result in all names/CRLs being invalidated.
	_, err = client.Logical().JSONMergePatch(ctx, "pki/issuer/old-root", map[string]interface{}{
		"issuer_name": "tmp-root",
	})
	require.NoError(t, err)
	_, err = client.Logical().JSONMergePatch(ctx, "pki/issuer/new-root", map[string]interface{}{
		"issuer_name": "old-root",
	})
	require.NoError(t, err)
	_, err = client.Logical().JSONMergePatch(ctx, "pki/issuer/tmp-root", map[string]interface{}{
		"issuer_name": "new-root",
	})
	require.NoError(t, err)

	afterThreeWaySwap := time.Now().Add(2 * time.Second)

	for _, path := range []string{"pki/cert/ca", "pki/cert/crl", "pki/issuer/default/json", "pki/issuer/old-root/json", "pki/issuer/new-root/json", "pki/issuer/old-root/crl", "pki/issuer/new-root/crl", "pki/cert/delta-crl", "pki/issuer/old-root/crl/delta", "pki/issuer/new-root/crl/delta"} {
		t.Logf("path: %v", path)
		field := "certificate"
		if strings.HasPrefix(path, "pki/issuer") && strings.Contains(path, "/crl") {
			field = "crl"
		}

		// Ensure that the CA is returned if we give it the pre-update time.
		client.AddHeader("If-Modified-Since", beforeThreeWaySwap.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data[field])
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()

		// Ensure that the CA is elided correctly if we give it the after time.
		client.AddHeader("If-Modified-Since", afterThreeWaySwap.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.Nil(t, resp)
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()
	}

	// Finally, rebuild the delta CRL and ensure that only that is
	// invalidated. We first need to enable it though, and wait for
	// all CRLs to rebuild.
	_, err = client.Logical().Write("pki/config/crl", map[string]interface{}{
		"auto_rebuild": true,
		"enable_delta": true,
	})
	require.NoError(t, err)
	time.Sleep(4 * time.Second)
	beforeDeltaRotation := time.Now().Add(-2 * time.Second)

	resp, err = client.Logical().Read("pki/crl/rotate-delta")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.Equal(t, resp.Data["success"], true)

	afterDeltaRotation := time.Now().Add(2 * time.Second)

	for _, path := range []string{"pki/cert/ca", "pki/cert/crl", "pki/issuer/default/json", "pki/issuer/old-root/json", "pki/issuer/new-root/json", "pki/issuer/old-root/crl", "pki/issuer/new-root/crl"} {
		t.Logf("path: %v", path)

		for _, when := range []time.Time{beforeDeltaRotation, afterDeltaRotation} {
			client.AddHeader("If-Modified-Since", when.Format(time.RFC1123))
			resp, err = client.Logical().Read(path)
			require.NoError(t, err)
			require.Nil(t, resp)
			client.SetHeaders(lastHeaders)
			lastHeaders = client.Headers()
		}
	}

	for _, path := range []string{"pki/cert/delta-crl", "pki/issuer/old-root/crl/delta", "pki/issuer/new-root/crl/delta"} {
		t.Logf("path: %v", path)
		field := "certificate"
		if strings.HasPrefix(path, "pki/issuer") && strings.Contains(path, "/crl") {
			field = "crl"
		}

		// Ensure that the CRL is present if we give it the pre-update time.
		client.AddHeader("If-Modified-Since", beforeDeltaRotation.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data[field])
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()

		client.AddHeader("If-Modified-Since", afterDeltaRotation.Format(time.RFC1123))
		resp, err = client.Logical().Read(path)
		require.NoError(t, err)
		require.Nil(t, resp)
		client.SetHeaders(lastHeaders)
		lastHeaders = client.Headers()
	}
}

func TestBackend_InitializeCertificateCounts(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	ctx := context.Background()

	// Set up an Issuer and Role
	// We need a root certificate to write/revoke certificates with
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected ca info")
	}

	// Create a role
	_, err = CBWrite(b, s, "roles/example", map[string]interface{}{
		"allowed_domains":    "myvault.com",
		"allow_bare_domains": true,
		"allow_subdomains":   true,
		"max_ttl":            "2h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Put certificates A, B, C, D, E in backend
	var certificates []string = []string{"a", "b", "c", "d", "e"}
	serials := make([]string, 5)
	for i, cn := range certificates {
		resp, err = CBWrite(b, s, "issue/example", map[string]interface{}{
			"common_name": cn + ".myvault.com",
		})
		if err != nil {
			t.Fatal(err)
		}
		serials[i] = resp.Data["serial_number"].(string)
	}

	// Turn on certificate counting:
	CBWrite(b, s, "config/auto-tidy", map[string]interface{}{
		"maintain_stored_certificate_counts":       true,
		"publish_stored_certificate_count_metrics": false,
	})
	// Assert initialize from clean is correct:
	b.initializeStoredCertificateCounts(ctx)

	// Revoke certificates A + B
	revocations := serials[0:2]
	for _, key := range revocations {
		resp, err = CBWrite(b, s, "revoke", map[string]interface{}{
			"serial_number": key,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	certCounter := b.GetCertificateCounter()
	if certCounter.CertificateCount() != 6 {
		t.Fatalf("Failed to count six certificates root,A,B,C,D,E, instead counted %d certs", certCounter.CertificateCount())
	}
	if certCounter.RevokedCount() != 2 {
		t.Fatalf("Failed to count two revoked certificates A+B, instead counted %d certs", certCounter.RevokedCount())
	}

	// Simulates listing while initialize in progress, by "restarting it"
	certCounter.certCount.Store(0)
	certCounter.revokedCertCount.Store(0)
	certCounter.certsCounted.Store(false)

	// Revoke certificates C, D
	dirtyRevocations := serials[2:4]
	for _, key := range dirtyRevocations {
		resp, err = CBWrite(b, s, "revoke", map[string]interface{}{
			"serial_number": key,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Put certificates F, G in the backend
	dirtyCertificates := []string{"f", "g"}
	for _, cn := range dirtyCertificates {
		resp, err = CBWrite(b, s, "issue/example", map[string]interface{}{
			"common_name": cn + ".myvault.com",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// Run initialize
	err = b.initializeStoredCertificateCounts(ctx)
	require.NoError(t, err, "failed initializing certificate counts")

	// Test certificate count
	if certCounter.CertificateCount() != 8 {
		t.Fatalf("Failed to initialize count of certificates root, A,B,C,D,E,F,G counted %d certs", certCounter.CertificateCount())
	}

	if certCounter.RevokedCount() != 4 {
		t.Fatalf("Failed to count revoked certificates A,B,C,D counted %d certs", certCounter.RevokedCount())
	}

	return
}

// Verify that our default values are consistent when creating an issuer and when we do an
// empty POST update to it. This will hopefully identify if we have different default values
// for fields across the two APIs.
func TestBackend_VerifyIssuerUpdateDefaultsMatchCreation(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "myvault.com",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed generating root issuer")

	resp, err = CBRead(b, s, "issuer/default")
	requireSuccessNonNilResponse(t, resp, err, "failed reading default issuer")
	preUpdateValues := resp.Data

	// This field gets reset during issuer update to the empty string
	// (meaning Go will auto-detect the rev-sig-algo).
	preUpdateValues["revocation_signature_algorithm"] = ""

	resp, err = CBWrite(b, s, "issuer/default", map[string]interface{}{})
	requireSuccessNonNilResponse(t, resp, err, "failed updating default issuer with no values")

	resp, err = CBRead(b, s, "issuer/default")
	requireSuccessNonNilResponse(t, resp, err, "failed reading default issuer")
	postUpdateValues := resp.Data

	require.Equal(t, preUpdateValues, postUpdateValues,
		"A value was updated based on the empty update of an issuer, "+
			"most likely we have a different set of field parameters across create and update of issuers.")
}

func TestBackend_VerifyPSSKeysIssuersFailImport(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// PKCS8 parsing fails on this key due to rsaPSS OID
	rsaOIDKey := `
-----BEGIN PRIVATE KEY-----
MIIEugIBADALBgkqhkiG9w0BAQoEggSmMIIEogIBAAKCAQEAtN0/NPuJHLuyEdBr
tUikXoXOV741XZcNvLAIVBIqDA0ege2gXt9A15FGUI4X3u6kT16Fl6MRdtUZ/qNS
Vs15nK9A1PI/AVekMgTVFTnoCzs550CKN8iRk9Om+lwHimpyXxKkFW69v8fsXwKE
Bsz69jjT7HV9VZQ7fQhmE79brAMuwKP1fUQKdHq5OBKtQ7Cl3Gmipp0izCsVuQIE
kBHvT3UUgyaSp2n+FONpOiyuBoYUH5tVEv9sZzBqSsrYBJYF+GvfnFy9AcTdqRe2
VX2SjjWjDF84T30OBA798gIFIPwu9R4OjWOlPeh2bo2kGeo3AITjwFZ28m7kS7kc
OtvHpwIDAQABAoIBAFQxmjbj0RQbG+3HBBzD0CBgUYnu9ZC3vKFVoMriGci6YrVB
FSKU8u5mpkDhpKMWnE6GRdItCvgyg4NSLAZUaIRT4O5ARqwtTDYsobTb2/U+gNnx
5WXKbFpQcK6jIK+ClfNEDjYb8yDPxG0GEsfHrBvqoFy25L1t37N4sWwH7HjJyZIe
Hbqx4NVDur9qgqaUwkfSeufn4ycHqFtkzKNzCUarDkST9cxE6/1AKfhl09PPuMEa
lAY2JLiEplQL5sh9cxG5FObJbutJo5EIhR2OdM0VcPf0MTD9LXKRoGR3SNlG7IlS
llJzBjlh4J1ByMX32btKMHzEvlhyrMI90E1SEGECgYEAx1yDQWe4/b1MBqCxA3d0
20dDmUHSRQFhkd/Mzkl5dPzRkG42W3ryNbMKdeuL0ZgK9AhfaLCjcj1i+44O7dHb
qBTVwfRrer2uoQVCqqJ6z8PGxPJJxTaqh9QuJxkoQ0i43ZNPcjc2M2sWLn+lkkdE
MaGMiyrmjIQEC6tmgCtZ1VUCgYEA6D9xoT9VuAnQjDvW2tO5N2U2H/8ZyRd1pC3z
H1CzjwShhxsP4YOUaVdw59K95JL4SMxSmpRrhthlW3cRaiT/exBcXLEvz0Qu0OhW
a6155ZFjK3UaLDKlwvmtuoAsuAFqX084LO0B1oxvUJESgyPncQ36fv2lZGV7A66z
Uo+BKQsCgYB2yGBMMAjA5nDN4iCV+C7gF+3m+pjWFKSVzcqxfoWndptGeuRYTUDT
TgIFkHqWPwkHrZVrQxOflYPMbi/m8wr1crSKA5+mWi4aMpAuKvERqYxc/B+IKbIh
jAKTuSGMNWAwZP0JCGx65mso+VUleuDe0Wpz4PPM9TuT2GQSKcI0oQKBgHAHcouC
npmo+lU65DgoWzaydrpWdpy+2Tt6AsW/Su4ZIMWoMy/oJaXuzQK2cG0ay/NpxArW
v0uLhNDrDZZzBF3blYIM4nALhr205UMJqjwntnuXACoDwFvdzoShIXEdFa+l6gYZ
yYIxudxWLmTd491wDb5GIgrcvMsY8V1I5dfjAoGAM9g2LtdqgPgK33dCDtZpBm8m
y4ri9PqHxnpps9WJ1dO6MW/YbW+a7vbsmNczdJ6XNLEfy2NWho1dw3xe7ztFVDjF
cWNUzs1+/6aFsi41UX7EFn3zAFhQUPxT59hXspuWuKbRAWc5fMnxbCfI/Cr8wTLJ
E/0kiZ4swUMyI4tYSbM=
-----END PRIVATE KEY-----
`
	_, err := CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": rsaOIDKey,
	})
	require.Error(t, err, "expected error importing PKCS8 rsaPSS OID key")

	_, err = CBWrite(b, s, "keys/import", map[string]interface{}{
		"key": rsaOIDKey,
	})
	require.Error(t, err, "expected error importing PKCS8 rsaPSS OID key")

	// Importing a cert with rsaPSS OID should also fail
	rsaOIDCert := `
-----BEGIN CERTIFICATE-----
MIIDfjCCAjGgAwIBAgIBATBCBgkqhkiG9w0BAQowNaAPMA0GCWCGSAFlAwQCAQUA
oRwwGgYJKoZIhvcNAQEIMA0GCWCGSAFlAwQCAQUAogQCAgDeMBMxETAPBgNVBAMM
CHJvb3Qtb2xkMB4XDTIyMDkxNjE0MDEwM1oXDTIzMDkyNjE0MDEwM1owEzERMA8G
A1UEAwwIcm9vdC1vbGQwggEgMAsGCSqGSIb3DQEBCgOCAQ8AMIIBCgKCAQEAtN0/
NPuJHLuyEdBrtUikXoXOV741XZcNvLAIVBIqDA0ege2gXt9A15FGUI4X3u6kT16F
l6MRdtUZ/qNSVs15nK9A1PI/AVekMgTVFTnoCzs550CKN8iRk9Om+lwHimpyXxKk
FW69v8fsXwKEBsz69jjT7HV9VZQ7fQhmE79brAMuwKP1fUQKdHq5OBKtQ7Cl3Gmi
pp0izCsVuQIEkBHvT3UUgyaSp2n+FONpOiyuBoYUH5tVEv9sZzBqSsrYBJYF+Gvf
nFy9AcTdqRe2VX2SjjWjDF84T30OBA798gIFIPwu9R4OjWOlPeh2bo2kGeo3AITj
wFZ28m7kS7kcOtvHpwIDAQABo3UwczAdBgNVHQ4EFgQUVGkTAUJ8inxIVGBlfxf4
cDhRSnowHwYDVR0jBBgwFoAUVGkTAUJ8inxIVGBlfxf4cDhRSnowDAYDVR0TBAUw
AwEB/zAOBgNVHQ8BAf8EBAMCAYYwEwYDVR0lBAwwCgYIKwYBBQUHAwEwQgYJKoZI
hvcNAQEKMDWgDzANBglghkgBZQMEAgEFAKEcMBoGCSqGSIb3DQEBCDANBglghkgB
ZQMEAgEFAKIEAgIA3gOCAQEAQZ3iQ3NjvS4FYJ5WG41huZI0dkvNFNan+ZYWlYHJ
MIQhbFogb/UQB0rlsuldG0+HF1RDXoYNuThfzt5hiBWYEtMBNurezvnOn4DF0hrl
Uk3sBVnvTalVXg+UVjqh9hBGB75JYJl6a5Oa2Zrq++4qGNwjd0FqgnoXzqS5UGuB
TJL8nlnXPuOIK3VHoXEy7l9GtvEzKcys0xa7g1PYpaJ5D2kpbBJmuQGmU6CDcbP+
m0hI4QDfVfHtnBp2VMCvhj0yzowtwF4BFIhv4EXZBU10mzxVj0zyKKft9++X8auH
nebuK22ZwzbPe4NhOvAdfNDElkrrtGvTnzkDB7ezPYjelA==
-----END CERTIFICATE-----
`
	_, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": rsaOIDCert,
	})
	require.Error(t, err, "expected error importing PKCS8 rsaPSS OID cert")

	_, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": rsaOIDKey + "\n" + rsaOIDCert,
	})
	require.Error(t, err, "expected error importing PKCS8 rsaPSS OID key+cert")

	_, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": rsaOIDCert + "\n" + rsaOIDKey,
	})
	require.Error(t, err, "expected error importing PKCS8 rsaPSS OID cert+key")

	// After all these errors, we should have zero issuers and keys.
	resp, err := CBList(b, s, "issuers")
	require.NoError(t, err)
	require.Equal(t, nil, resp.Data["keys"])

	resp, err = CBList(b, s, "keys")
	require.NoError(t, err)
	require.Equal(t, nil, resp.Data["keys"])

	// If we create a new PSS root, we should be able to issue an intermediate
	// under it.
	resp, err = CBWrite(b, s, "root/generate/exported", map[string]interface{}{
		"use_pss":     "true",
		"common_name": "root x1 - pss",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["certificate"])
	require.NotEmpty(t, resp.Data["private_key"])

	resp, err = CBWrite(b, s, "intermediate/generate/exported", map[string]interface{}{
		"use_pss":     "true",
		"common_name": "int x1 - pss",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["csr"])
	require.NotEmpty(t, resp.Data["private_key"])

	resp, err = CBWrite(b, s, "issuer/default/sign-intermediate", map[string]interface{}{
		"use_pss":     "true",
		"common_name": "int x1 - pss",
		"csr":         resp.Data["csr"].(string),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["certificate"])

	resp, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": resp.Data["certificate"].(string),
	})
	require.NoError(t, err)

	// Finally, if we were to take an rsaPSS OID'd CSR and use it against this
	// mount, it will fail.
	_, err = CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allow_any_name": true,
		"ttl":            "85s",
		"key_type":       "any",
	})
	require.NoError(t, err)

	// Issuing a leaf from a CSR with rsaPSS OID should fail...
	rsaOIDCSR := `-----BEGIN CERTIFICATE REQUEST-----
MIICkTCCAUQCAQAwGTEXMBUGA1UEAwwOcmFuY2hlci5teS5vcmcwggEgMAsGCSqG
SIb3DQEBCgOCAQ8AMIIBCgKCAQEAtzHuGEUK55lXI08yp9DXoye9yCZbkJZO+Hej
1TWGEkbX4hzauRJeNp2+wn8xU5y8ITjWSIXEVDHeezosLCSy0Y2QT7/V45zWPUYY
ld0oUnPiwsb9CPFlBRFnX3dO9SS5MONIrNCJGKXmLdF3lgSl8zPT6J/hWM+JBjHO
hBzK6L8IYwmcEujrQfnOnOztzgMEBJtWG8rnI8roz1adpczTddDKGymh2QevjhlL
X9CLeYSSQZInOMsgaDYl98Hn00K5x0CBp8ADzzXtaPSQ9nsnihN8VvZ/wHw6YbBS
BSHa6OD+MrYnw3Sao6/YgBRNT2glIX85uro4ARW9zGB9/748dwIDAQABoAAwQgYJ
KoZIhvcNAQEKMDWgDzANBglghkgBZQMEAgEFAKEcMBoGCSqGSIb3DQEBCDANBglg
hkgBZQMEAgEFAKIEAgIA3gOCAQEARGAa0HiwzWCpvAdLOVc4/srEyOYFZPLbtv+Y
ezZIaUBNaWhOvkunqpa48avmcbGlji7r6fxJ5sT28lHt7ODWcJfn1XPAnqesXErm
EBuOIhCv6WiwVyGeTVynuHYkHyw3rIL/zU7N8+zIFV2G2M1UAv5D/eyh/74cr9Of
+nvm9jAbkHix8UwOBCFY2LLNl6bXvbIeJEdDOEtA9UmDXs8QGBg4lngyqcE2Z7rz
+5N/x4guMk2FqblbFGiCc5fLB0Gp6lFFOqhX9Q8nLJ6HteV42xGJUUtsFpppNCRm
82dGIH2PTbXZ0k7iAAwLaPjzOv1v58Wq90o35d4iEsOfJ8v98Q==
-----END CERTIFICATE REQUEST-----`

	_, err = CBWrite(b, s, "issuer/default/sign/testing", map[string]interface{}{
		"common_name": "example.com",
		"csr":         rsaOIDCSR,
	})
	require.Error(t, err)

	_, err = CBWrite(b, s, "issuer/default/sign-verbatim", map[string]interface{}{
		"common_name": "example.com",
		"use_pss":     true,
		"csr":         rsaOIDCSR,
	})
	require.Error(t, err)

	_, err = CBWrite(b, s, "issuer/default/sign-intermediate", map[string]interface{}{
		"common_name": "faulty x1 - pss",
		"use_pss":     true,
		"csr":         rsaOIDCSR,
	})
	require.Error(t, err)

	// Vault has a weird API for signing self-signed certificates. Ensure
	// that doesn't accept rsaPSS OID'd certificates either.
	_, err = CBWrite(b, s, "issuer/default/sign-self-issued", map[string]interface{}{
		"use_pss":     true,
		"certificate": rsaOIDCert,
	})
	require.Error(t, err)

	// Issuing a regular leaf should succeed.
	_, err = CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allow_any_name": true,
		"ttl":            "85s",
		"key_type":       "rsa",
		"use_pss":        "true",
	})
	require.NoError(t, err)

	resp, err = CBWrite(b, s, "issuer/default/issue/testing", map[string]interface{}{
		"common_name": "example.com",
		"use_pss":     "true",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to issue PSS leaf")
}

func TestPKI_EmptyCRLConfigUpgraded(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Write an empty CRLConfig into storage.
	crlConfigEntry, err := logical.StorageEntryJSON("config/crl", &pki_backend.CrlConfig{})
	require.NoError(t, err)
	err = s.Put(ctx, crlConfigEntry)
	require.NoError(t, err)

	resp, err := CBRead(b, s, "config/crl")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.Equal(t, resp.Data["expiry"], pki_backend.DefaultCrlConfig.Expiry)
	require.Equal(t, resp.Data["disable"], pki_backend.DefaultCrlConfig.Disable)
	require.Equal(t, resp.Data["ocsp_disable"], pki_backend.DefaultCrlConfig.OcspDisable)
	require.Equal(t, resp.Data["auto_rebuild"], pki_backend.DefaultCrlConfig.AutoRebuild)
	require.Equal(t, resp.Data["auto_rebuild_grace_period"], pki_backend.DefaultCrlConfig.AutoRebuildGracePeriod)
	require.Equal(t, resp.Data["enable_delta"], pki_backend.DefaultCrlConfig.EnableDelta)
	require.Equal(t, resp.Data["delta_rebuild_interval"], pki_backend.DefaultCrlConfig.DeltaRebuildInterval)
	require.Equal(t, resp.Data["max_crl_entries"], pki_backend.DefaultCrlConfig.MaxCRLEntries)
}

func TestPKI_ListRevokedCerts(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Test empty cluster
	resp, err := CBList(b, s, "certs/revoked")
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("certs/revoked"), logical.ListOperation), resp, true)
	requireSuccessNonNilResponse(t, resp, err, "failed listing empty cluster")
	require.Empty(t, resp.Data, "response map contained data that we did not expect")

	// Set up a mount that we can revoke under (We will create 3 leaf certs, 2 of which will be revoked)
	resp, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "test.com",
		"key_type":    "ec",
	})
	requireSuccessNonNilResponse(t, resp, err, "error generating root CA")
	requireFieldsSetInResp(t, resp, "serial_number")
	issuerSerial := resp.Data["serial_number"]

	resp, err = CBWrite(b, s, "roles/test", map[string]interface{}{
		"allowed_domains":  "test.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
	})
	requireSuccessNonNilResponse(t, resp, err, "error setting up pki role")

	resp, err = CBWrite(b, s, "issue/test", map[string]interface{}{
		"common_name": "test1.test.com",
	})
	requireSuccessNonNilResponse(t, resp, err, "error issuing cert 1")
	requireFieldsSetInResp(t, resp, "serial_number")
	serial1 := resp.Data["serial_number"]

	resp, err = CBWrite(b, s, "issue/test", map[string]interface{}{
		"common_name": "test2.test.com",
	})
	requireSuccessNonNilResponse(t, resp, err, "error issuing cert 2")
	requireFieldsSetInResp(t, resp, "serial_number")
	serial2 := resp.Data["serial_number"]

	resp, err = CBWrite(b, s, "issue/test", map[string]interface{}{
		"common_name": "test3.test.com",
	})
	requireSuccessNonNilResponse(t, resp, err, "error issuing cert 2")
	requireFieldsSetInResp(t, resp, "serial_number")
	serial3 := resp.Data["serial_number"]

	resp, err = CBWrite(b, s, "revoke", map[string]interface{}{"serial_number": serial1})
	requireSuccessNonNilResponse(t, resp, err, "error revoking cert 1")

	resp, err = CBWrite(b, s, "revoke", map[string]interface{}{"serial_number": serial2})
	requireSuccessNonNilResponse(t, resp, err, "error revoking cert 2")

	// Test that we get back the expected revoked serial numbers.
	resp, err = CBList(b, s, "certs/revoked")
	requireSuccessNonNilResponse(t, resp, err, "failed listing revoked certs")
	requireFieldsSetInResp(t, resp, "keys")
	revokedKeys := resp.Data["keys"].([]string)

	require.Contains(t, revokedKeys, serial1)
	require.Contains(t, revokedKeys, serial2)
	require.Equal(t, 2, len(revokedKeys), "Expected 2 revoked entries got %d: %v", len(revokedKeys), revokedKeys)

	// Test that listing our certs returns a different response
	resp, err = CBList(b, s, "certs")
	requireSuccessNonNilResponse(t, resp, err, "failed listing written certs")
	requireFieldsSetInResp(t, resp, "keys")
	certKeys := resp.Data["keys"].([]string)

	require.Contains(t, certKeys, serial1)
	require.Contains(t, certKeys, serial2)
	require.Contains(t, certKeys, serial3)
	require.Contains(t, certKeys, issuerSerial)
	require.Equal(t, 4, len(certKeys), "Expected 4 cert entries got %d: %v", len(certKeys), certKeys)
}

func TestPKI_TemplatedAIAs(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Setting templated AIAs should succeed.
	resp, err := CBWrite(b, s, "config/cluster", map[string]interface{}{
		"path":     "http://localhost:8200/v1/pki",
		"aia_path": "http://localhost:8200/cdn/pki",
	})
	require.NoError(t, err)
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("config/cluster"), logical.UpdateOperation), resp, true)

	resp, err = CBRead(b, s, "config/cluster")
	require.NoError(t, err)
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("config/cluster"), logical.ReadOperation), resp, true)

	aiaData := map[string]interface{}{
		"crl_distribution_points": "{{cluster_path}}/issuer/{{issuer_id}}/crl/der",
		"issuing_certificates":    "{{cluster_aia_path}}/issuer/{{issuer_id}}/der",
		"ocsp_servers":            "{{cluster_path}}/ocsp",
		"enable_templating":       true,
	}
	_, err = CBWrite(b, s, "config/urls", aiaData)
	require.NoError(t, err)

	// Root generation should succeed, but without AIA info.
	rootData := map[string]interface{}{
		"common_name": "Long-Lived Root X1",
		"issuer_name": "long-root-x1",
		"key_type":    "ec",
	}
	resp, err = CBWrite(b, s, "root/generate/internal", rootData)
	require.NoError(t, err)
	_, err = CBDelete(b, s, "root")
	require.NoError(t, err)

	// Clearing the config and regenerating the root should still succeed.
	_, err = CBWrite(b, s, "config/urls", map[string]interface{}{
		"crl_distribution_points": "{{cluster_path}}/issuer/my-root-id/crl/der",
		"issuing_certificates":    "{{cluster_aia_path}}/issuer/my-root-id/der",
		"ocsp_servers":            "{{cluster_path}}/ocsp",
		"enable_templating":       true,
	})
	require.NoError(t, err)
	resp, err = CBWrite(b, s, "root/generate/internal", rootData)
	requireSuccessNonNilResponse(t, resp, err)
	issuerId := string(resp.Data["issuer_id"].(issuing.IssuerID))

	// Now write the original AIA config and sign a leaf.
	_, err = CBWrite(b, s, "config/urls", aiaData)
	require.NoError(t, err)
	_, err = CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allow_any_name": "true",
		"key_type":       "ec",
		"ttl":            "50m",
	})
	require.NoError(t, err)
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "example.com",
	})
	requireSuccessNonNilResponse(t, resp, err)

	// Validate the AIA info is correctly templated.
	cert := parseCert(t, resp.Data["certificate"].(string))
	require.Equal(t, cert.OCSPServer, []string{"http://localhost:8200/v1/pki/ocsp"})
	require.Equal(t, cert.IssuingCertificateURL, []string{"http://localhost:8200/cdn/pki/issuer/" + issuerId + "/der"})
	require.Equal(t, cert.CRLDistributionPoints, []string{"http://localhost:8200/v1/pki/issuer/" + issuerId + "/crl/der"})

	// Modify our issuer to set custom AIAs: these URLs are bad.
	_, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
		"enable_aia_url_templating": "false",
		"crl_distribution_points":   "a",
		"issuing_certificates":      "b",
		"ocsp_servers":              "c",
	})
	require.Error(t, err)

	// These URLs are good.
	_, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
		"enable_aia_url_templating": "false",
		"crl_distribution_points":   "http://localhost/a",
		"issuing_certificates":      "http://localhost/b",
		"ocsp_servers":              "http://localhost/c",
	})

	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "example.com",
	})
	requireSuccessNonNilResponse(t, resp, err)

	// Validate the AIA info is correctly templated.
	cert = parseCert(t, resp.Data["certificate"].(string))
	require.Equal(t, cert.OCSPServer, []string{"http://localhost/c"})
	require.Equal(t, cert.IssuingCertificateURL, []string{"http://localhost/b"})
	require.Equal(t, cert.CRLDistributionPoints, []string{"http://localhost/a"})

	// These URLs are bad, but will fail at issuance time due to AIA templating.
	resp, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
		"enable_aia_url_templating": "true",
		"crl_distribution_points":   "a",
		"issuing_certificates":      "b",
		"ocsp_servers":              "c",
	})
	requireSuccessNonNilResponse(t, resp, err)
	require.NotEmpty(t, resp.Warnings)
	_, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "example.com",
	})
	require.Error(t, err)
}

func requireSubjectUserIDAttr(t *testing.T, cert string, target string) {
	xCert := parseCert(t, cert)

	for _, attr := range xCert.Subject.Names {
		var userID string
		if attr.Type.Equal(certutil.SubjectPilotUserIDAttributeOID) {
			if target == "" {
				t.Fatalf("expected no UserID (OID: %v) subject attributes in cert:\n%v", certutil.SubjectPilotUserIDAttributeOID, cert)
			}

			switch aValue := attr.Value.(type) {
			case string:
				userID = aValue
			case []byte:
				userID = string(aValue)
			default:
				t.Fatalf("unknown type for UserID attribute: %v\nCert: %v", attr, cert)
			}

			if userID == target {
				return
			}
		}
	}

	if target != "" {
		t.Fatalf("failed to find UserID (OID: %v) matching %v in cert:\n%v", certutil.SubjectPilotUserIDAttributeOID, target, cert)
	}
}

func TestUserIDsInLeafCerts(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// 1. Setup root issuer.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "Vault Root CA",
		"key_type":    "ec",
		"ttl":         "7200h",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed generating root issuer")

	// 2. Allow no user IDs.
	resp, err = CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allowed_user_ids": "",
		"key_type":         "ec",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed setting up role")

	// - Issue cert without user IDs should work.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "")

	// - Issue cert with user ID should fail.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "humanoid",
	})
	require.Error(t, err)
	require.True(t, resp.IsError())

	// 3. Allow any user IDs.
	resp, err = CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allowed_user_ids": "*",
		"key_type":         "ec",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed setting up role")

	// - Issue cert without user IDs.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "")

	// - Issue cert with one user ID.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "humanoid",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "humanoid")

	// - Issue cert with two user IDs.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "humanoid,robot",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "humanoid")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "robot")

	// 4. Allow one specific user ID.
	resp, err = CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allowed_user_ids": "humanoid",
		"key_type":         "ec",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed setting up role")

	// - Issue cert without user IDs.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "")

	// - Issue cert with approved ID.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "humanoid",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "humanoid")

	// - Issue cert with non-approved user ID should fail.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "robot",
	})
	require.Error(t, err)
	require.True(t, resp.IsError())

	// - Issue cert with one approved and one non-approved should also fail.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "humanoid,robot",
	})
	require.Error(t, err)
	require.True(t, resp.IsError())

	// 5. Allow two specific user IDs.
	resp, err = CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allowed_user_ids": "humanoid,robot",
		"key_type":         "ec",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed setting up role")

	// - Issue cert without user IDs.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "")

	// - Issue cert with one approved ID.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "humanoid",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "humanoid")

	// - Issue cert with other user ID.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "robot",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "robot")

	// - Issue cert with unknown user ID will fail.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "robot2",
	})
	require.Error(t, err)
	require.True(t, resp.IsError())

	// - Issue cert with both should succeed.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "humanoid,robot",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "humanoid")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "robot")

	// 6. Use a glob.
	resp, err = CBWrite(b, s, "roles/testing", map[string]interface{}{
		"allowed_user_ids": "human*",
		"key_type":         "ec",
		"use_csr_sans":     true, // setup for further testing.
	})
	requireSuccessNonNilResponse(t, resp, err, "failed setting up role")

	// - Issue cert without user IDs.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "")

	// - Issue cert with approved ID.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "humanoid",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "humanoid")

	// - Issue cert with another approved ID.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "human",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "human")

	// - Issue cert with literal glob.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "human*",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "human*")

	// - Still no robotic certs are allowed; will fail.
	resp, err = CBWrite(b, s, "issue/testing", map[string]interface{}{
		"common_name": "localhost",
		"user_ids":    "robot",
	})
	require.Error(t, err)
	require.True(t, resp.IsError())

	// Create a CSR and validate it works with both sign/ and sign-verbatim.
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: "localhost",
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type:  certutil.SubjectPilotUserIDAttributeOID,
					Value: "humanoid",
				},
			},
		},
	}
	_, _, csrPem := generateCSR(t, &csrTemplate, "ec", 256)

	// Should work with role-based signing.
	resp, err = CBWrite(b, s, "sign/testing", map[string]interface{}{
		"csr": csrPem,
	})
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("sign/testing"), logical.UpdateOperation), resp, true)
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "humanoid")

	// - Definitely will work with sign-verbatim.
	resp, err = CBWrite(b, s, "sign-verbatim", map[string]interface{}{
		"csr": csrPem,
	})
	requireSuccessNonNilResponse(t, resp, err, "failed issuing leaf cert")
	requireSubjectUserIDAttr(t, resp.Data["certificate"].(string), "humanoid")
}

// TestStandby_Operations test proper forwarding for PKI requests from a standby node to the
// active node within a cluster.
func TestStandby_Operations(t *testing.T) {
	conf, opts := teststorage.ClusterSetup(&vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
	}, nil, teststorage.InmemBackendSetup)
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	testhelpers.WaitForActiveNodeAndStandbys(t, cluster)
	standbyCores := testhelpers.DeriveStandbyCores(t, cluster)
	require.Greater(t, len(standbyCores), 0, "Need at least one standby core.")
	client := standbyCores[0].Client

	mountPKIEndpoint(t, client, "pki")

	_, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "root-ca.com",
		"ttl":         "600h",
	})
	require.NoError(t, err, "error setting up pki role: %v", err)

	_, err = client.Logical().Write("pki/roles/example", map[string]interface{}{
		"allowed_domains":  "example.com",
		"allow_subdomains": "true",
		"no_store":         "false", // make sure we store this cert
		"ttl":              "5h",
		"key_type":         "ec",
	})
	require.NoError(t, err, "error setting up pki role: %v", err)

	resp, err := client.Logical().Write("pki/issue/example", map[string]interface{}{
		"common_name": "test.example.com",
	})
	require.NoError(t, err, "error issuing certificate: %v", err)
	require.NotNil(t, resp, "got nil response from issuing request")
	serialOfCert := resp.Data["serial_number"].(string)

	resp, err = client.Logical().Write("pki/revoke", map[string]interface{}{
		"serial_number": serialOfCert,
	})
	require.NoError(t, err, "error revoking certificate: %v", err)
	require.NotNil(t, resp, "got nil response from revoke request")
}

type pathAuthCheckerFunc func(t *testing.T, client *api.Client, path string, token string)

func isPermDenied(err error) bool {
	return err != nil && strings.Contains(err.Error(), "permission denied")
}

func isUnsupportedPathOperation(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "unsupported path") || strings.Contains(err.Error(), "unsupported operation"))
}

func isDeniedOp(err error) bool {
	return isPermDenied(err) || isUnsupportedPathOperation(err)
}

func pathShouldBeAuthed(t *testing.T, client *api.Client, path string, token string) {
	client.SetToken("")
	resp, err := client.Logical().ReadWithContext(ctx, path)
	if err == nil || !isPermDenied(err) {
		t.Fatalf("expected failure to read %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().ListWithContext(ctx, path)
	if err == nil || !isPermDenied(err) {
		t.Fatalf("expected failure to list %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().WriteWithContext(ctx, path, map[string]interface{}{})
	if err == nil || !isPermDenied(err) {
		t.Fatalf("expected failure to write %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().DeleteWithContext(ctx, path)
	if err == nil || !isPermDenied(err) {
		t.Fatalf("expected failure to delete %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().JSONMergePatch(ctx, path, map[string]interface{}{})
	if err == nil || !isPermDenied(err) {
		t.Fatalf("expected failure to patch %v while unauthed: %v / %v", path, err, resp)
	}
}

func pathShouldBeUnauthedReadList(t *testing.T, client *api.Client, path string, token string) {
	// Should be able to read both with and without a token.
	client.SetToken("")
	resp, err := client.Logical().ReadWithContext(ctx, path)
	if err != nil && isPermDenied(err) {
		// Read will sometimes return permission denied, when the handler
		// does not support the given operation. Retry with the token.
		client.SetToken(token)
		resp2, err2 := client.Logical().ReadWithContext(ctx, path)
		if err2 != nil && !isUnsupportedPathOperation(err2) {
			t.Fatalf("unexpected failure to read %v while unauthed: %v / %v\nWhile authed: %v / %v", path, err, resp, err2, resp2)
		}
		client.SetToken("")
	}
	resp, err = client.Logical().ListWithContext(ctx, path)
	if err != nil && isPermDenied(err) {
		// List will sometimes return permission denied, when the handler
		// does not support the given operation. Retry with the token.
		client.SetToken(token)
		resp2, err2 := client.Logical().ListWithContext(ctx, path)
		if err2 != nil && !isUnsupportedPathOperation(err2) {
			t.Fatalf("unexpected failure to list %v while unauthed: %v / %v\nWhile authed: %v / %v", path, err, resp, err2, resp2)
		}
		client.SetToken("")
	}

	// These should all be denied.
	resp, err = client.Logical().WriteWithContext(ctx, path, map[string]interface{}{})
	if err == nil || !isDeniedOp(err) {
		if !strings.Contains(path, "ocsp") || !strings.Contains(err.Error(), "Code: 40") {
			t.Fatalf("unexpected failure during write on read-only path %v while unauthed: %v / %v", path, err, resp)
		}
	}
	resp, err = client.Logical().DeleteWithContext(ctx, path)
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during delete on read-only path %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().JSONMergePatch(ctx, path, map[string]interface{}{})
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during patch on read-only path %v while unauthed: %v / %v", path, err, resp)
	}

	// Retrying with token should allow read/list, but not modification still.
	client.SetToken(token)
	resp, err = client.Logical().ReadWithContext(ctx, path)
	if err != nil && isPermDenied(err) {
		t.Fatalf("unexpected failure to read %v while authed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().ListWithContext(ctx, path)
	if err != nil && isPermDenied(err) {
		t.Fatalf("unexpected failure to list %v while authed: %v / %v", path, err, resp)
	}

	// Should all be denied.
	resp, err = client.Logical().WriteWithContext(ctx, path, map[string]interface{}{})
	if err == nil || !isDeniedOp(err) {
		if !strings.Contains(path, "ocsp") || !strings.Contains(err.Error(), "Code: 40") {
			t.Fatalf("unexpected failure during write on read-only path %v while authed: %v / %v", path, err, resp)
		}
	}
	resp, err = client.Logical().DeleteWithContext(ctx, path)
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during delete on read-only path %v while authed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().JSONMergePatch(ctx, path, map[string]interface{}{})
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during patch on read-only path %v while authed: %v / %v", path, err, resp)
	}
}

func pathShouldBeUnauthedWriteOnly(t *testing.T, client *api.Client, path string, token string) {
	client.SetToken("")
	resp, err := client.Logical().WriteWithContext(ctx, path, map[string]interface{}{})
	if err != nil && isPermDenied(err) {
		t.Fatalf("unexpected failure to write %v while unauthed: %v / %v", path, err, resp)
	}

	// These should all be denied. However, on OSS, we might end up with
	// a regular 404, which looks like err == resp == nil; hence we only
	// fail when there's a non-nil response and/or a non-nil err.
	resp, err = client.Logical().ReadWithContext(ctx, path)
	if (err == nil && resp != nil) || (err != nil && !isDeniedOp(err)) {
		t.Fatalf("unexpected failure during read on write-only path %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().ListWithContext(ctx, path)
	if (err == nil && resp != nil) || (err != nil && !isDeniedOp(err)) {
		t.Fatalf("unexpected failure during list on write-only path %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().DeleteWithContext(ctx, path)
	if (err == nil && resp != nil) || (err != nil && !isDeniedOp(err)) {
		t.Fatalf("unexpected failure during delete on write-only path %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().JSONMergePatch(ctx, path, map[string]interface{}{})
	if (err == nil && resp != nil) || (err != nil && !isDeniedOp(err)) {
		t.Fatalf("unexpected failure during patch on write-only path %v while unauthed: %v / %v", path, err, resp)
	}

	// Retrying with token should allow writing, but nothing else.
	client.SetToken(token)
	resp, err = client.Logical().WriteWithContext(ctx, path, map[string]interface{}{})
	if err != nil && isPermDenied(err) {
		t.Fatalf("unexpected failure to write %v while unauthed: %v / %v", path, err, resp)
	}

	// These should all be denied.
	resp, err = client.Logical().ReadWithContext(ctx, path)
	if (err == nil && resp != nil) || (err != nil && !isDeniedOp(err)) {
		t.Fatalf("unexpected failure during read on write-only path %v while authed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().ListWithContext(ctx, path)
	if (err == nil && resp != nil) || (err != nil && !isDeniedOp(err)) {
		if resp != nil || err != nil {
			t.Fatalf("unexpected failure during list on write-only path %v while authed: %v / %v", path, err, resp)
		}
	}
	resp, err = client.Logical().DeleteWithContext(ctx, path)
	if (err == nil && resp != nil) || (err != nil && !isDeniedOp(err)) {
		t.Fatalf("unexpected failure during delete on write-only path %v while authed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().JSONMergePatch(ctx, path, map[string]interface{}{})
	if (err == nil && resp != nil) || (err != nil && !isDeniedOp(err)) {
		t.Fatalf("unexpected failure during patch on write-only path %v while authed: %v / %v", path, err, resp)
	}
}

type pathAuthChecker int

const (
	shouldBeAuthed pathAuthChecker = iota
	shouldBeUnauthedReadList
	shouldBeUnauthedWriteOnly
	shouldBeUnauthedReadWriteOnly
)

var pathAuthChckerMap = map[pathAuthChecker]pathAuthCheckerFunc{
	shouldBeAuthed:                pathShouldBeAuthed,
	shouldBeUnauthedReadList:      pathShouldBeUnauthedReadList,
	shouldBeUnauthedWriteOnly:     pathShouldBeUnauthedWriteOnly,
	shouldBeUnauthedReadWriteOnly: pathShouldBeUnauthedWriteOnly,
}

func TestProperAuthing(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
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
	token := client.Token()

	// Mount PKI.
	err := client.Sys().MountWithContext(ctx, "pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "60h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Setup basic configuration.
	_, err = client.Logical().WriteWithContext(ctx, "pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "myvault.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().WriteWithContext(ctx, "pki/roles/test", map[string]interface{}{
		"allow_localhost": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().WriteWithContext(ctx, "pki/issue/test", map[string]interface{}{
		"common_name": "localhost",
	})
	if err != nil || resp == nil {
		t.Fatal(err)
	}
	serial := resp.Data["serial_number"].(string)
	eabKid := "13b80844-e60d-42d2-b7e9-152a8e834b90"
	acmeKeyId := "hrKmDYTvicHoHGVN2-3uzZV_BPGdE0W_dNaqYTtYqeo="
	paths := map[string]pathAuthChecker{
		"ca_chain":                               shouldBeUnauthedReadList,
		"cert/ca_chain":                          shouldBeUnauthedReadList,
		"ca":                                     shouldBeUnauthedReadList,
		"ca/pem":                                 shouldBeUnauthedReadList,
		"cert/" + serial:                         shouldBeUnauthedReadList,
		"cert/" + serial + "/raw":                shouldBeUnauthedReadList,
		"cert/" + serial + "/raw/pem":            shouldBeUnauthedReadList,
		"cert/crl":                               shouldBeUnauthedReadList,
		"cert/crl/raw":                           shouldBeUnauthedReadList,
		"cert/crl/raw/pem":                       shouldBeUnauthedReadList,
		"cert/delta-crl":                         shouldBeUnauthedReadList,
		"cert/delta-crl/raw":                     shouldBeUnauthedReadList,
		"cert/delta-crl/raw/pem":                 shouldBeUnauthedReadList,
		"cert/unified-crl":                       shouldBeUnauthedReadList,
		"cert/unified-crl/raw":                   shouldBeUnauthedReadList,
		"cert/unified-crl/raw/pem":               shouldBeUnauthedReadList,
		"cert/unified-delta-crl":                 shouldBeUnauthedReadList,
		"cert/unified-delta-crl/raw":             shouldBeUnauthedReadList,
		"cert/unified-delta-crl/raw/pem":         shouldBeUnauthedReadList,
		issuing.PathCerts:                        shouldBeAuthed,
		"certs/revoked/":                         shouldBeAuthed,
		"certs/revocation-queue/":                shouldBeAuthed,
		"certs/unified-revoked/":                 shouldBeAuthed,
		"config/acme":                            shouldBeAuthed,
		"config/auto-tidy":                       shouldBeAuthed,
		"config/ca":                              shouldBeAuthed,
		"config/cluster":                         shouldBeAuthed,
		"config/crl":                             shouldBeAuthed,
		"config/issuers":                         shouldBeAuthed,
		"config/keys":                            shouldBeAuthed,
		"config/urls":                            shouldBeAuthed,
		"crl":                                    shouldBeUnauthedReadList,
		"crl/pem":                                shouldBeUnauthedReadList,
		"crl/delta":                              shouldBeUnauthedReadList,
		"crl/delta/pem":                          shouldBeUnauthedReadList,
		"crl/rotate":                             shouldBeAuthed,
		"crl/rotate-delta":                       shouldBeAuthed,
		"intermediate/cross-sign":                shouldBeAuthed,
		"intermediate/generate/exported":         shouldBeAuthed,
		"intermediate/generate/internal":         shouldBeAuthed,
		"intermediate/generate/existing":         shouldBeAuthed,
		"intermediate/generate/kms":              shouldBeAuthed,
		"intermediate/set-signed":                shouldBeAuthed,
		"issue/test":                             shouldBeAuthed,
		"issuer/default":                         shouldBeAuthed,
		"issuer/default/der":                     shouldBeUnauthedReadList,
		"issuer/default/json":                    shouldBeUnauthedReadList,
		"issuer/default/pem":                     shouldBeUnauthedReadList,
		"issuer/default/crl":                     shouldBeUnauthedReadList,
		"issuer/default/crl/pem":                 shouldBeUnauthedReadList,
		"issuer/default/crl/der":                 shouldBeUnauthedReadList,
		"issuer/default/crl/delta":               shouldBeUnauthedReadList,
		"issuer/default/crl/delta/der":           shouldBeUnauthedReadList,
		"issuer/default/crl/delta/pem":           shouldBeUnauthedReadList,
		"issuer/default/unified-crl":             shouldBeUnauthedReadList,
		"issuer/default/unified-crl/pem":         shouldBeUnauthedReadList,
		"issuer/default/unified-crl/der":         shouldBeUnauthedReadList,
		"issuer/default/unified-crl/delta":       shouldBeUnauthedReadList,
		"issuer/default/unified-crl/delta/der":   shouldBeUnauthedReadList,
		"issuer/default/unified-crl/delta/pem":   shouldBeUnauthedReadList,
		"issuer/default/issue/test":              shouldBeAuthed,
		"issuer/default/resign-crls":             shouldBeAuthed,
		"issuer/default/revoke":                  shouldBeAuthed,
		"issuer/default/sign-intermediate":       shouldBeAuthed,
		"issuer/default/sign-revocation-list":    shouldBeAuthed,
		"issuer/default/sign-self-issued":        shouldBeAuthed,
		"issuer/default/sign-verbatim":           shouldBeAuthed,
		"issuer/default/sign-verbatim/test":      shouldBeAuthed,
		"issuer/default/sign/test":               shouldBeAuthed,
		"issuers/":                               shouldBeUnauthedReadList,
		"issuers/generate/intermediate/exported": shouldBeAuthed,
		"issuers/generate/intermediate/internal": shouldBeAuthed,
		"issuers/generate/intermediate/existing": shouldBeAuthed,
		"issuers/generate/intermediate/kms":      shouldBeAuthed,
		"issuers/generate/root/exported":         shouldBeAuthed,
		"issuers/generate/root/internal":         shouldBeAuthed,
		"issuers/generate/root/existing":         shouldBeAuthed,
		"issuers/generate/root/kms":              shouldBeAuthed,
		"issuers/import/cert":                    shouldBeAuthed,
		"issuers/import/bundle":                  shouldBeAuthed,
		"key/default":                            shouldBeAuthed,
		"keys/":                                  shouldBeAuthed,
		"keys/generate/internal":                 shouldBeAuthed,
		"keys/generate/exported":                 shouldBeAuthed,
		"keys/generate/kms":                      shouldBeAuthed,
		"keys/import":                            shouldBeAuthed,
		"ocsp":                                   shouldBeUnauthedWriteOnly,
		"ocsp/dGVzdAo=":                          shouldBeUnauthedReadList,
		"revoke":                                 shouldBeAuthed,
		"revoke-with-key":                        shouldBeAuthed,
		"roles/test":                             shouldBeAuthed,
		"roles/":                                 shouldBeAuthed,
		"root":                                   shouldBeAuthed,
		"root/generate/exported":                 shouldBeAuthed,
		"root/generate/internal":                 shouldBeAuthed,
		"root/generate/existing":                 shouldBeAuthed,
		"root/generate/kms":                      shouldBeAuthed,
		"root/replace":                           shouldBeAuthed,
		"root/rotate/internal":                   shouldBeAuthed,
		"root/rotate/exported":                   shouldBeAuthed,
		"root/rotate/existing":                   shouldBeAuthed,
		"root/rotate/kms":                        shouldBeAuthed,
		"root/sign-intermediate":                 shouldBeAuthed,
		"root/sign-self-issued":                  shouldBeAuthed,
		"sign-verbatim":                          shouldBeAuthed,
		"sign-verbatim/test":                     shouldBeAuthed,
		"sign/test":                              shouldBeAuthed,
		"tidy":                                   shouldBeAuthed,
		"tidy-cancel":                            shouldBeAuthed,
		"tidy-status":                            shouldBeAuthed,
		"unified-crl":                            shouldBeUnauthedReadList,
		"unified-crl/pem":                        shouldBeUnauthedReadList,
		"unified-crl/delta":                      shouldBeUnauthedReadList,
		"unified-crl/delta/pem":                  shouldBeUnauthedReadList,
		"unified-ocsp":                           shouldBeUnauthedWriteOnly,
		"unified-ocsp/dGVzdAo=":                  shouldBeUnauthedReadList,
		"eab/":                                   shouldBeAuthed,
		"eab/" + eabKid:                          shouldBeAuthed,
		"acme/mgmt/account/keyid/":               shouldBeAuthed,
		"acme/mgmt/account/keyid/" + acmeKeyId:   shouldBeAuthed,
	}

	entPaths := getEntProperAuthingPaths(serial)
	maps.Copy(paths, entPaths)

	// Add ACME based paths to the test suite
	ossAcmePrefixes := []string{"acme/", "issuer/default/acme/", "roles/test/acme/", "issuer/default/roles/test/acme/"}
	entAcmePrefixes := getEntAcmePrefixes()
	for _, acmePrefix := range append(ossAcmePrefixes, entAcmePrefixes...) {
		paths[acmePrefix+"directory"] = shouldBeUnauthedReadList
		paths[acmePrefix+"new-nonce"] = shouldBeUnauthedReadList
		paths[acmePrefix+"new-account"] = shouldBeUnauthedWriteOnly
		paths[acmePrefix+"revoke-cert"] = shouldBeUnauthedWriteOnly
		paths[acmePrefix+"new-order"] = shouldBeUnauthedWriteOnly
		paths[acmePrefix+"orders"] = shouldBeUnauthedWriteOnly
		paths[acmePrefix+"account/hrKmDYTvicHoHGVN2-3uzZV_BPGdE0W_dNaqYTtYqeo="] = shouldBeUnauthedWriteOnly
		paths[acmePrefix+"authorization/29da8c38-7a09-465e-b9a6-3d76802b1afd"] = shouldBeUnauthedWriteOnly
		paths[acmePrefix+"challenge/29da8c38-7a09-465e-b9a6-3d76802b1afd/http-01"] = shouldBeUnauthedWriteOnly
		paths[acmePrefix+"order/13b80844-e60d-42d2-b7e9-152a8e834b90"] = shouldBeUnauthedWriteOnly
		paths[acmePrefix+"order/13b80844-e60d-42d2-b7e9-152a8e834b90/finalize"] = shouldBeUnauthedWriteOnly
		paths[acmePrefix+"order/13b80844-e60d-42d2-b7e9-152a8e834b90/cert"] = shouldBeUnauthedWriteOnly

		// Make sure this new-eab path is auth'd
		paths[acmePrefix+"new-eab"] = shouldBeAuthed
	}

	for path, checkerType := range paths {
		checker := pathAuthChckerMap[checkerType]
		checker(t, client, "pki/"+path, token)
	}

	client.SetToken(token)
	openAPIResp, err := client.Logical().ReadWithContext(ctx, "sys/internal/specs/openapi")
	if err != nil {
		t.Fatalf("failed to get openapi data: %v", err)
	}

	validatedPath := false
	for openapi_path, raw_data := range openAPIResp.Data["paths"].(map[string]interface{}) {
		if !strings.HasPrefix(openapi_path, "/pki/") {
			t.Logf("Skipping path: %v", openapi_path)
			continue
		}

		t.Logf("Validating path: %v", openapi_path)
		validatedPath = true
		// Substitute values in from our testing map.
		raw_path := openapi_path[5:]
		if strings.Contains(raw_path, "roles/") && strings.Contains(raw_path, "{name}") {
			raw_path = strings.ReplaceAll(raw_path, "{name}", "test")
		}
		if strings.Contains(raw_path, "{role}") {
			raw_path = strings.ReplaceAll(raw_path, "{role}", "test")
		}
		if strings.Contains(raw_path, "ocsp/") && strings.Contains(raw_path, "{req}") {
			raw_path = strings.ReplaceAll(raw_path, "{req}", "dGVzdAo=")
		}
		if strings.Contains(raw_path, "{issuer_ref}") {
			raw_path = strings.ReplaceAll(raw_path, "{issuer_ref}", "default")
		}
		if strings.Contains(raw_path, "{key_ref}") {
			raw_path = strings.ReplaceAll(raw_path, "{key_ref}", "default")
		}
		if strings.Contains(raw_path, "{exported}") {
			raw_path = strings.ReplaceAll(raw_path, "{exported}", "internal")
		}
		if strings.Contains(raw_path, "{serial}") {
			raw_path = strings.ReplaceAll(raw_path, "{serial}", serial)
		}
		if strings.Contains(raw_path, "acme/account/") && strings.Contains(raw_path, "{kid}") {
			raw_path = strings.ReplaceAll(raw_path, "{kid}", acmeKeyId)
		}
		if strings.Contains(raw_path, "acme/mgmt/account/") && strings.Contains(raw_path, "{keyid}") {
			raw_path = strings.ReplaceAll(raw_path, "{keyid}", acmeKeyId)
		}
		if strings.Contains(raw_path, "acme/") && strings.Contains(raw_path, "{auth_id}") {
			raw_path = strings.ReplaceAll(raw_path, "{auth_id}", "29da8c38-7a09-465e-b9a6-3d76802b1afd")
		}
		if strings.Contains(raw_path, "acme/") && strings.Contains(raw_path, "{challenge_type}") {
			raw_path = strings.ReplaceAll(raw_path, "{challenge_type}", "http-01")
		}
		if strings.Contains(raw_path, "acme/") && strings.Contains(raw_path, "{order_id}") {
			raw_path = strings.ReplaceAll(raw_path, "{order_id}", "13b80844-e60d-42d2-b7e9-152a8e834b90")
		}
		if strings.Contains(raw_path, "eab") && strings.Contains(raw_path, "{key_id}") {
			raw_path = strings.ReplaceAll(raw_path, "{key_id}", eabKid)
		}
		if strings.Contains(raw_path, "external-policy/") && strings.Contains(raw_path, "{policy}") {
			raw_path = strings.ReplaceAll(raw_path, "{policy}", "a-policy")
		}

		raw_path = entProperAuthingPathReplacer(raw_path)

		handler, present := paths[raw_path]
		if !present {
			t.Fatalf("OpenAPI reports PKI mount contains %v -> %v but was not tested to be authed or not authed.",
				openapi_path, raw_path)
		}

		openapi_data := raw_data.(map[string]interface{})
		hasList := false
		rawGetData, hasGet := openapi_data["get"]
		if hasGet {
			getData := rawGetData.(map[string]interface{})
			getParams, paramsPresent := getData["parameters"].(map[string]interface{})
			if getParams != nil && paramsPresent {
				if _, hasList = getParams["list"]; hasList {
					// LIST is exclusive from GET on the same endpoint usually.
					hasGet = false
				}
			}
		}
		_, hasPost := openapi_data["post"]
		_, hasDelete := openapi_data["delete"]

		if handler == shouldBeUnauthedReadList {
			if hasPost || hasDelete {
				t.Fatalf("Unauthed read-only endpoints should not have POST/DELETE capabilities: %v->%v", openapi_path, raw_path)
			}
		} else if handler == shouldBeUnauthedWriteOnly {
			if hasGet || hasList {
				t.Fatalf("Unauthed write-only endpoints should not have GET/LIST capabilities: %v->%v", openapi_path, raw_path)
			}
		} else if handler == shouldBeUnauthedReadWriteOnly {
			if hasDelete || hasList {
				t.Fatalf("Unauthed read-write-only endpoints should not have DELETE/LIST capabilities: %v->%v", openapi_path, raw_path)
			}
		}
	}

	if !validatedPath {
		t.Fatalf("Expected to have validated at least one path.")
	}
}

type patchIssuerTestCase struct {
	Field   string
	Before  interface{}
	Patched interface{}
}

func TestPatchIssuer(t *testing.T) {
	t.Parallel()

	testCases := []patchIssuerTestCase{
		{
			Field:   "issuer_name",
			Before:  "root",
			Patched: "root-new",
		},
		{
			Field:   "leaf_not_after_behavior",
			Before:  "err",
			Patched: "permit",
		},
		{
			Field:   "usage",
			Before:  "crl-signing,issuing-certificates,ocsp-signing,read-only",
			Patched: "issuing-certificates,read-only",
		},
		{
			Field:   "revocation_signature_algorithm",
			Before:  "ECDSAWithSHA256",
			Patched: "ECDSAWithSHA384",
		},
		{
			Field:   "issuing_certificates",
			Before:  []string{"http://localhost/v1/pki-1/ca"},
			Patched: []string{"http://localhost/v1/pki/ca"},
		},
		{
			Field:   "crl_distribution_points",
			Before:  []string{"http://localhost/v1/pki-1/crl"},
			Patched: []string{"http://localhost/v1/pki/crl"},
		},
		{
			Field:   "ocsp_servers",
			Before:  []string{"http://localhost/v1/pki-1/ocsp"},
			Patched: []string{"http://localhost/v1/pki/ocsp"},
		},
		{
			Field:   "enable_aia_url_templating",
			Before:  false,
			Patched: true,
		},
		{
			Field:   "manual_chain",
			Before:  []string(nil),
			Patched: []string{"self"},
		},
	}
	testPatchIssuer(t, testCases)
}

func testPatchIssuer(t *testing.T, testCases []patchIssuerTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.Field, func(t *testing.T) {
			b, s := CreateBackendWithStorage(t)

			// 1. Setup root issuer.
			resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
				"common_name": "Vault Root CA",
				"key_type":    "ec",
				"ttl":         "7200h",
				"issuer_name": "root",
			})
			requireSuccessNonNilResponse(t, resp, err, "failed generating root issuer")
			id := string(resp.Data["issuer_id"].(issuing.IssuerID))

			// 2. Enable Cluster paths
			resp, err = CBWrite(b, s, "config/urls", map[string]interface{}{
				"path":     "https://localhost/v1/pki",
				"aia_path": "http://localhost/v1/pki",
			})
			requireSuccessNonNilResponse(t, resp, err, "failed updating AIA config")

			// 3. Add AIA information
			resp, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
				"issuing_certificates":    "http://localhost/v1/pki-1/ca",
				"crl_distribution_points": "http://localhost/v1/pki-1/crl",
				"ocsp_servers":            "http://localhost/v1/pki-1/ocsp",
			})
			requireSuccessNonNilResponse(t, resp, err, "failed setting up issuer")

			// 4. Read the issuer before.
			resp, err = CBRead(b, s, "issuer/default")
			requireSuccessNonNilResponse(t, resp, err, "failed reading root issuer before")
			require.Equal(t, testCase.Before, resp.Data[testCase.Field], "bad expectations")

			// 5. Perform modification.
			resp, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
				testCase.Field: testCase.Patched,
			})
			requireSuccessNonNilResponse(t, resp, err, "failed patching root issuer")

			if testCase.Field != "manual_chain" {
				require.Equal(t, testCase.Patched, resp.Data[testCase.Field], "failed persisting value")
			} else {
				// self->id
				require.Equal(t, []string{id}, resp.Data[testCase.Field], "failed persisting value")
			}

			// 6. Ensure it stuck
			resp, err = CBRead(b, s, "issuer/default")
			requireSuccessNonNilResponse(t, resp, err, "failed reading root issuer after")

			if testCase.Field != "manual_chain" {
				require.Equal(t, testCase.Patched, resp.Data[testCase.Field])
			} else {
				// self->id
				require.Equal(t, []string{id}, resp.Data[testCase.Field], "failed persisting value")
			}

			// 7. Patch it back
			resp, err = CBPatch(b, s, "issuer/default", map[string]interface{}{
				testCase.Field: testCase.Before,
			})
			requireSuccessNonNilResponse(t, resp, err, "failed patching root issuer")

			require.Equal(t, testCase.Before, resp.Data[testCase.Field], "failed persisting value")

			// 8. Ensure it stuck
			resp, err = CBRead(b, s, "issuer/default")
			requireSuccessNonNilResponse(t, resp, err, "failed reading root issuer after")

			require.Equal(t, testCase.Before, resp.Data[testCase.Field])
		})
	}
}

func TestGenerateRootCAWithAIA(t *testing.T) {
	// Generate a root CA at /pki-root
	b_root, s_root := CreateBackendWithStorage(t)

	// Setup templated AIA information
	_, err := CBWrite(b_root, s_root, "config/cluster", map[string]interface{}{
		"path":     "https://localhost:8200",
		"aia_path": "https://localhost:8200",
	})
	require.NoError(t, err, "failed to write AIA settings")

	_, err = CBWrite(b_root, s_root, "config/urls", map[string]interface{}{
		"crl_distribution_points": "{{cluster_path}}/issuer/{{issuer_id}}/crl/der",
		"issuing_certificates":    "{{cluster_aia_path}}/issuer/{{issuer_id}}/der",
		"ocsp_servers":            "{{cluster_path}}/ocsp",
		"enable_templating":       true,
	})
	require.NoError(t, err, "failed to write AIA settings")

	// Write a root issuer, this should succeed.
	resp, err := CBWrite(b_root, s_root, "root/generate/exported", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_type":    "ec",
	})
	requireSuccessNonNilResponse(t, resp, err, "expected root generation to succeed")
}

// TestIssuance_AlwaysEnforceErr validates that we properly return an error in all request
// types that go beyond the issuer's NotAfter
func TestIssuance_AlwaysEnforceErr(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root myvault.com",
		"key_type":    "ec",
		"ttl":         "10h",
		"issuer_name": "root-ca",
		"key_name":    "root-key",
	})
	requireSuccessNonNilResponse(t, resp, err, "expected root generation to succeed")

	resp, err = CBPatch(b, s, "issuer/root-ca", map[string]interface{}{
		"leaf_not_after_behavior": "always_enforce_err",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed updating root issuer with always_enforce_err")

	resp, err = CBWrite(b, s, "roles/test-role", map[string]interface{}{
		"allow_any_name":         true,
		"key_type":               "ec",
		"allowed_serial_numbers": "*",
	})

	expectedErrContains := "cannot satisfy request, as TTL would result in notAfter"

	// Make sure we fail on CA issuance requests now
	t.Run("ca-issuance", func(t *testing.T) {
		resp, err = CBWrite(b, s, "intermediate/generate/internal", map[string]interface{}{
			"common_name": "myint.com",
		})
		requireSuccessNonNilResponse(t, resp, err, "failed generating intermediary CSR")
		requireFieldsSetInResp(t, resp, "csr")
		csr := resp.Data["csr"]

		_, err = CBWrite(b, s, "issuer/root-ca/sign-intermediate", map[string]interface{}{
			"csr":            csr,
			"use_csr_values": true,
			"ttl":            "60h",
		})
		require.ErrorContains(t, err, expectedErrContains, "sign-intermediate should have failed as root issuer leaf behavior is set to always_enforce_err")

		// Make sure it works if we are under
		resp, err = CBWrite(b, s, "issuer/root-ca/sign-intermediate", map[string]interface{}{
			"csr":            csr,
			"use_csr_values": true,
			"ttl":            "30m",
		})
		requireSuccessNonNilResponse(t, resp, err, "sign-intermediate should have passed with a lower TTL value and always_enforce_err")
	})

	// Make sure we fail on leaf csr signing leaf as we always did for 'err'
	t.Run("sign-leaf-csr", func(t *testing.T) {
		_, csrPem := generateTestCsr(t, certutil.ECPrivateKey, 256)

		resp, err = CBWrite(b, s, "issuer/root-ca/sign/test-role", map[string]interface{}{
			"ttl": "60h",
			"csr": csrPem,
		})
		require.ErrorContains(t, err, expectedErrContains, "expected error from sign csr got: %v", resp)

		// Make sure it works if we are under
		resp, err = CBWrite(b, s, "issuer/root-ca/sign/test-role", map[string]interface{}{
			"ttl": "30m",
			"csr": csrPem,
		})
		requireSuccessNonNilResponse(t, resp, err, "sign should have succeeded with a lower TTL and always_enforce_err")
	})

	// Make sure we fail on leaf csr signing leaf as we always did for 'err'
	t.Run("issue-leaf-csr", func(t *testing.T) {
		resp, err = CBWrite(b, s, "issuer/root-ca/issue/test-role", map[string]interface{}{
			"ttl":         "60h",
			"common_name": "leaf.example.com",
		})
		require.ErrorContains(t, err, expectedErrContains, "expected error from issue got: %v", resp)

		// Make sure it works if we are under
		resp, err = CBWrite(b, s, "issuer/root-ca/issue/test-role", map[string]interface{}{
			"ttl":         "30m",
			"common_name": "leaf.example.com",
		})
		requireSuccessNonNilResponse(t, resp, err, "issue should have worked with a lower TTL and always_enforce_err")
	})
}

var (
	initTest  sync.Once
	rsaCAKey  string
	rsaCACert string
	ecCAKey   string
	ecCACert  string
	edCAKey   string
	edCACert  string
)
