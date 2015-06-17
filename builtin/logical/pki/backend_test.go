package pki

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math"
	"math/rand"
	"os"
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

func checkCertsAndPrivateKey(keyType string, usage certUsage, validity time.Duration, certBundle *certutil.CertBundle) (cert, ca *x509.Certificate, privKey crypto.Signer, err error) {
	var pemBlock *pem.Block

	pemBlock, _ = pem.Decode([]byte(certBundle.Certificate))
	if pemBlock == nil {
		return nil, nil, nil, fmt.Errorf("No PEM data found for cert")
	}
	cert, err = x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, nil, nil, err
	}

	pemBlock, _ = pem.Decode([]byte(certBundle.IssuingCA))
	if pemBlock == nil {
		return nil, nil, nil, fmt.Errorf("No PEM data found for issuing CA")
	}
	ca, err = x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, nil, nil, err
	}

	pemBlock, _ = pem.Decode([]byte(certBundle.PrivateKey))
	if pemBlock == nil {
		return nil, nil, nil, fmt.Errorf("No PEM data found for private key")
	}
	switch keyType {
	case "rsa":
		privKey, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
		if err != nil {
			return nil, nil, nil, err
		}
	case "ec":
		privKey, err = x509.ParseECPrivateKey(pemBlock.Bytes)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Unable to decode EC private key: %s; value was %s", err, certBundle.PrivateKey)
		}
	default:
		return nil, nil, nil, fmt.Errorf("Unknown private key type %s", keyType)
	}

	// There should only be one usage type, because only one is requested
	// in the tests
	if len(cert.ExtKeyUsage) != 1 {
		return cert, nil, nil, fmt.Errorf("Got wrong size key usage in generated cert")
	}
	switch usage {
	case serverUsage:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageServerAuth {
			return cert, nil, nil, fmt.Errorf("Bad key usage")
		}
	case clientUsage:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth {
			return cert, nil, nil, fmt.Errorf("Bad key usage")
		}
	case codeSigningUsage:
		if cert.ExtKeyUsage[0] != x509.ExtKeyUsageCodeSigning {
			return cert, nil, nil, fmt.Errorf("Bad key usage")
		}
	}

	if math.Abs(float64(time.Now().Unix()-cert.NotBefore.Unix())) > 10 {
		return cert, nil, nil, fmt.Errorf("Validity period starts out of range")
	}

	if math.Abs(float64(time.Now().Add(validity).Unix()-cert.NotAfter.Unix())) > 10 {
		return cert, nil, nil, fmt.Errorf("Validity period too large")
	}

	return
}

func TestBackend_basic(t *testing.T) {
	b := Backend()

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps:   []logicaltest.TestStep{},
	}

	stepCount += len(testCase.Steps)

	testCase.Steps = append(testCase.Steps, generateCASteps(t)...)

	logicaltest.Test(t, testCase)
}

func TestBackend_roles(t *testing.T) {
	b := Backend()

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps:   []logicaltest.TestStep{},
	}

	testCase.Steps = append(testCase.Steps, generateCASteps(t)...)
	testCase.Steps = append(testCase.Steps, generateRoleSteps(t)...)
	if len(os.Getenv("VAULT_VERBOSE_PKITESTS")) > 0 {
		for i, v := range testCase.Steps {
			fmt.Printf("Step %d:\n%+v\n\n", i+stepCount, v)
		}
	}

	stepCount += len(testCase.Steps)

	logicaltest.Test(t, testCase)
}

func generateCASteps(t *testing.T) []logicaltest.TestStep {
	ret := []logicaltest.TestStep{
		logicaltest.TestStep{
			Operation: logical.WriteOperation,
			Path:      "config/ca",
			Data: map[string]interface{}{
				"pem_bundle": caKey + caCert,
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
	}

	return ret
}

func generateRoleSteps(t *testing.T) []logicaltest.TestStep {
	roleVals := roleEntry{
		LeaseMax: "12h",
	}
	issueVals := certutil.IssueData{}
	ret := []logicaltest.TestStep{}

	roleTestStep := logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "roles/test",
	}
	issueTestStep := logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "issue/test",
	}

	genericErrorOkCheck := func(resp *logical.Response) error {
		if resp.IsError() {
			return nil
		}
		return fmt.Errorf("Expected an error, but did not seem to get one")
	}

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

	getCnCheck := func(name, keyType string, usage certUsage, validity time.Duration) logicaltest.TestCheckFunc {
		var certBundle certutil.CertBundle
		return func(resp *logical.Response) error {
			err := mapstructure.Decode(resp.Data, &certBundle)
			if err != nil {
				return err
			}
			cert, _, _, err := checkCertsAndPrivateKey(keyType, usage, validity, &certBundle)
			if err != nil {
				return fmt.Errorf("Error checking generated certificate: %s", err)
			}
			if cert.Subject.CommonName != name {
				return fmt.Errorf("Error: returned certificate has CN of %s but %s was requested", cert.Subject.CommonName, name)
			}
			if len(cert.DNSNames) != 1 {
				return fmt.Errorf("Error: found more than one DNS SAN but only one was requested")
			}
			if cert.DNSNames[0] != name {
				return fmt.Errorf("Error: returned certificate has a DNS SAN of %s but %s was requested", cert.DNSNames[0], name)
			}
			return nil
		}
	}

	var commonNames struct {
		Localhost         bool `structs:"localhost"`
		BaseDomain        bool `structs:"foo.example.com"`
		Wildcard          bool `structs:"*.example.com"`
		Subdomain         bool `structs:"foo.bar.example.com"`
		SubdomainWildcard bool `structs:"*.bar.example.com"`
		AnyHost           bool `structs:"porkslap.beer"`
	}

	addCnTests := func() {
		cnMap := structs.New(commonNames).Map()
		// For the number of tests being run, this is known to hit all
		// of the various values below
		mathRand := rand.New(rand.NewSource(1))
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
			var usage certUsage
			i := mathRand.Int()
			switch {
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
			if allowed {
				issueTestStep.ErrorOk = false
			} else {
				issueTestStep.ErrorOk = true
			}

			validity, _ := time.ParseDuration(roleVals.LeaseMax)
			addTests(getCnCheck(name, roleVals.KeyType, usage, validity))
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
		commonNames.BaseDomain = true
		commonNames.Wildcard = true
		addCnTests()

		roleVals.AllowSubdomains = true
		commonNames.Subdomain = true
		commonNames.SubdomainWildcard = true
		addCnTests()

		roleVals.AllowAnyName = true
		commonNames.AnyHost = true
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
		roleVals.LeaseMax = ""
		addTests(nil)

		roleVals.Lease = "12h"
		roleVals.LeaseMax = "6h"
		addTests(nil)

		roleTestStep.ErrorOk = false
	}

	return ret
}

const (
	caKey string = `-----BEGIN RSA PRIVATE KEY-----
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
	caCert string = `-----BEGIN CERTIFICATE-----
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
)
