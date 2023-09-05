package server

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/hashicorp/go-secure-stdlib/strutil"
)

// TestGenerateCertExtraSans ensures the implementation backing the flag
// -dev-tls-san populates alternate DNS and IP address names in the generated
// certificate as expected.
func TestGenerateCertExtraSans(t *testing.T) {
	ca, err := GenerateCA()
	if err != nil {
		t.Fatal(err)
	}

	for name, tc := range map[string]struct {
		extraSans           []string
		expectedDNSNames    []string
		expectedIPAddresses []string
	}{
		"empty": {},
		"DNS names": {
			extraSans:        []string{"foo", "foo.bar"},
			expectedDNSNames: []string{"foo", "foo.bar"},
		},
		"IP addresses": {
			extraSans:           []string{"0.0.0.0", "::1"},
			expectedIPAddresses: []string{"0.0.0.0", "::1"},
		},
		"mixed": {
			extraSans:           []string{"bar", "0.0.0.0", "::1"},
			expectedDNSNames:    []string{"bar"},
			expectedIPAddresses: []string{"0.0.0.0", "::1"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			certStr, _, err := generateCert(ca.Template, ca.Signer, tc.extraSans)
			if err != nil {
				t.Fatal(err)
			}

			block, _ := pem.Decode([]byte(certStr))
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				t.Fatal(err)
			}

			expectedDNSNamesLen := len(tc.expectedDNSNames) + 5
			if len(cert.DNSNames) != expectedDNSNamesLen {
				t.Errorf("Wrong number of DNS names, expected %d but got %v", expectedDNSNamesLen, cert.DNSNames)
			}
			expectedIPAddrLen := len(tc.expectedIPAddresses) + 1
			if len(cert.IPAddresses) != expectedIPAddrLen {
				t.Errorf("Wrong number of IP addresses, expected %d but got %v", expectedIPAddrLen, cert.IPAddresses)
			}

			for _, expected := range tc.expectedDNSNames {
				if !strutil.StrListContains(cert.DNSNames, expected) {
					t.Errorf("Missing DNS name %s", expected)
				}
			}
			for _, expected := range tc.expectedIPAddresses {
				var found bool
				for _, ip := range cert.IPAddresses {
					if ip.String() == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Missing IP address %s", expected)
				}
			}
		})
	}
}
