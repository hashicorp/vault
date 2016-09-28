package certutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/api"
)

// Tests converting back and forth between a CertBundle and a ParsedCertBundle.
//
// Also tests the GetSubjKeyID, GetHexFormatted, and
// ParsedCertBundle.getSigner functions.
func TestCertBundleConversion(t *testing.T) {
	cbuts := []*CertBundle{
		refreshRSACertBundle(),
		refreshRSACertBundleWithChain(),
		refreshRSA8CertBundle(),
		refreshRSA8CertBundleWithChain(),
		refreshECCertBundle(),
		refreshECCertBundleWithChain(),
		refreshEC8CertBundle(),
		refreshEC8CertBundleWithChain(),
	}

	for i, cbut := range cbuts {
		pcbut, err := cbut.ToParsedCertBundle()
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Errorf("Error converting to parsed cert bundle: %s", err)
			continue
		}

		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Errorf(err.Error())
		}

		cbut, err := pcbut.ToCertBundle()
		if err != nil {
			t.Fatalf("Error converting to cert bundle: %s", err)
		}

		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func BenchmarkCertBundleParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cbuts := []*CertBundle{
			refreshRSACertBundle(),
			refreshRSACertBundleWithChain(),
			refreshRSA8CertBundle(),
			refreshRSA8CertBundleWithChain(),
			refreshECCertBundle(),
			refreshECCertBundleWithChain(),
			refreshEC8CertBundle(),
			refreshEC8CertBundleWithChain(),
		}

		for i, cbut := range cbuts {
			pcbut, err := cbut.ToParsedCertBundle()
			if err != nil {
				b.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
				b.Errorf("Error converting to parsed cert bundle: %s", err)
				continue
			}

			cbut, err = pcbut.ToCertBundle()
			if err != nil {
				b.Fatalf("Error converting to cert bundle: %s", err)
			}
		}
	}
}

func TestCertBundleParsing(t *testing.T) {
	cbuts := []*CertBundle{
		refreshRSACertBundle(),
		refreshRSACertBundleWithChain(),
		refreshRSA8CertBundle(),
		refreshRSA8CertBundleWithChain(),
		refreshECCertBundle(),
		refreshECCertBundleWithChain(),
		refreshEC8CertBundle(),
		refreshEC8CertBundleWithChain(),
	}

	for i, cbut := range cbuts {
		jsonString, err := json.Marshal(cbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf("Error marshaling testing certbundle to JSON: %s", err)
		}
		pcbut, err := ParsePKIJSON(jsonString)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf("Error during JSON bundle handling: %s", err)
		}
		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf(err.Error())
		}

		secret := &api.Secret{
			Data: structs.New(cbut).Map(),
		}
		pcbut, err = ParsePKIMap(secret.Data)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf("Error during JSON bundle handling: %s", err)
		}
		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf(err.Error())
		}

		pcbut, err = ParsePEMBundle(cbut.ToPEMBundle())
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf("Error during JSON bundle handling: %s", err)
		}
		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf(err.Error())
		}
	}
}

func compareCertBundleToParsedCertBundle(cbut *CertBundle, pcbut *ParsedCertBundle) error {
	if cbut == nil {
		return fmt.Errorf("Got nil bundle")
	}
	if pcbut == nil {
		return fmt.Errorf("Got nil parsed bundle")
	}

	switch {
	case pcbut.Certificate == nil:
		return fmt.Errorf("Parsed bundle has nil certificate")
	case pcbut.PrivateKey == nil:
		return fmt.Errorf("Parsed bundle has nil private key")
	}

	switch cbut.PrivateKey {
	case privRSAKeyPem:
		if pcbut.PrivateKeyType != RSAPrivateKey {
			return fmt.Errorf("Parsed bundle has wrong private key type: %v, should be 'rsa' (%v)", pcbut.PrivateKeyType, RSAPrivateKey)
		}
	case privRSA8KeyPem:
		if pcbut.PrivateKeyType != RSAPrivateKey {
			return fmt.Errorf("Parsed bundle has wrong pkcs8 private key type: %v, should be 'rsa' (%v)", pcbut.PrivateKeyType, RSAPrivateKey)
		}
	case privECKeyPem:
		if pcbut.PrivateKeyType != ECPrivateKey {
			return fmt.Errorf("Parsed bundle has wrong private key type: %v, should be 'ec' (%v)", pcbut.PrivateKeyType, ECPrivateKey)
		}
	case privEC8KeyPem:
		if pcbut.PrivateKeyType != ECPrivateKey {
			return fmt.Errorf("Parsed bundle has wrong pkcs8 private key type: %v, should be 'ec' (%v)", pcbut.PrivateKeyType, ECPrivateKey)
		}
	default:
		return fmt.Errorf("Parsed bundle has unknown private key type")
	}

	subjKeyID, err := GetSubjKeyID(pcbut.PrivateKey)
	if err != nil {
		return fmt.Errorf("Error when getting subject key id: %s", err)
	}
	if bytes.Compare(subjKeyID, pcbut.Certificate.SubjectKeyId) != 0 {
		return fmt.Errorf("Parsed bundle private key does not match subject key id")
	}

	switch {
	case len(pcbut.CAChain) > 0 && len(cbut.CAChain) == 0:
		return fmt.Errorf("Parsed bundle ca chain has certs when cert bundle does not")
	case len(pcbut.CAChain) == 0 && len(cbut.CAChain) > 0:
		return fmt.Errorf("Cert bundle ca chain has certs when parsed cert bundle does not")
	}

	cb, err := pcbut.ToCertBundle()
	if err != nil {
		return fmt.Errorf("Thrown error during parsed bundle conversion: %s\n\nInput was: %#v", err, *pcbut)
	}

	switch {
	case len(cb.Certificate) == 0:
		return fmt.Errorf("Bundle has nil certificate")
	case len(cb.PrivateKey) == 0:
		return fmt.Errorf("Bundle has nil private key")
	case len(cb.CAChain[0]) == 0:
		return fmt.Errorf("Bundle has nil issuing CA")
	}

	switch pcbut.PrivateKeyType {
	case RSAPrivateKey:
		if cb.PrivateKey != privRSAKeyPem && cb.PrivateKey != privRSA8KeyPem {
			return fmt.Errorf("Bundle private key does not match")
		}
	case ECPrivateKey:
		if cb.PrivateKey != privECKeyPem && cb.PrivateKey != privEC8KeyPem {
			return fmt.Errorf("Bundle private key does not match")
		}
	default:
		return fmt.Errorf("CertBundle has unknown private key type")
	}

	if cb.SerialNumber != GetHexFormatted(pcbut.Certificate.SerialNumber.Bytes(), ":") {
		return fmt.Errorf("Bundle serial number does not match")
	}

	switch {
	case len(pcbut.CAChain) > 0 && len(cb.CAChain) == 0:
		return fmt.Errorf("Parsed bundle ca chain has certs when cert bundle does not")
	case len(pcbut.CAChain) == 0 && len(cb.CAChain) > 0:
		return fmt.Errorf("Cert bundle ca chain has certs when parsed cert bundle does not")
	case !reflect.DeepEqual(cbut.CAChain, cb.CAChain):
		return fmt.Errorf("Cert bundle ca chain does not match: %#v\n\n%#v", cbut.CAChain, cb.CAChain)
	}

	return nil
}

func TestCSRBundleConversion(t *testing.T) {
	csrbuts := []*CSRBundle{
		refreshRSACSRBundle(),
		refreshECCSRBundle(),
	}

	for _, csrbut := range csrbuts {
		pcsrbut, err := csrbut.ToParsedCSRBundle()
		if err != nil {
			t.Fatalf("Error converting to parsed CSR bundle: %v", err)
		}

		err = compareCSRBundleToParsedCSRBundle(csrbut, pcsrbut)
		if err != nil {
			t.Fatalf(err.Error())
		}

		csrbut, err = pcsrbut.ToCSRBundle()
		if err != nil {
			t.Fatalf("Error converting to CSR bundle: %v", err)
		}

		err = compareCSRBundleToParsedCSRBundle(csrbut, pcsrbut)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func compareCSRBundleToParsedCSRBundle(csrbut *CSRBundle, pcsrbut *ParsedCSRBundle) error {
	if csrbut == nil {
		return fmt.Errorf("Got nil bundle")
	}
	if pcsrbut == nil {
		return fmt.Errorf("Got nil parsed bundle")
	}

	switch {
	case pcsrbut.CSR == nil:
		return fmt.Errorf("Parsed bundle has nil csr")
	case pcsrbut.PrivateKey == nil:
		return fmt.Errorf("Parsed bundle has nil private key")
	}

	switch csrbut.PrivateKey {
	case privRSAKeyPem:
		if pcsrbut.PrivateKeyType != RSAPrivateKey {
			return fmt.Errorf("Parsed bundle has wrong private key type")
		}
	case privECKeyPem:
		if pcsrbut.PrivateKeyType != ECPrivateKey {
			return fmt.Errorf("Parsed bundle has wrong private key type")
		}
	default:
		return fmt.Errorf("Parsed bundle has unknown private key type")
	}

	csrb, err := pcsrbut.ToCSRBundle()
	if err != nil {
		return fmt.Errorf("Thrown error during parsed bundle conversion: %s\n\nInput was: %#v", err, *pcsrbut)
	}

	switch {
	case len(csrb.CSR) == 0:
		return fmt.Errorf("Bundle has nil certificate")
	case len(csrb.PrivateKey) == 0:
		return fmt.Errorf("Bundle has nil private key")
	}

	switch csrb.PrivateKeyType {
	case "rsa":
		if pcsrbut.PrivateKeyType != RSAPrivateKey {
			return fmt.Errorf("Bundle has wrong private key type")
		}
		if csrb.PrivateKey != privRSAKeyPem {
			return fmt.Errorf("Bundle private key does not match")
		}
	case "ec":
		if pcsrbut.PrivateKeyType != ECPrivateKey {
			return fmt.Errorf("Bundle has wrong private key type")
		}
		if csrb.PrivateKey != privECKeyPem {
			return fmt.Errorf("Bundle private key does not match")
		}
	default:
		return fmt.Errorf("Bundle has unknown private key type")
	}

	return nil
}

func TestTLSConfig(t *testing.T) {
	cbut := refreshRSACertBundle()

	pcbut, err := cbut.ToParsedCertBundle()
	if err != nil {
		t.Fatalf("Error getting parsed cert bundle: %s", err)
	}

	usages := []TLSUsage{
		TLSUnknown,
		TLSClient,
		TLSServer,
		TLSClient | TLSServer,
	}

	for _, usage := range usages {
		tlsConfig, err := pcbut.GetTLSConfig(usage)
		if err != nil {
			t.Fatalf("Error getting tls config: %s", err)
		}
		if tlsConfig == nil {
			t.Fatalf("Got nil tls.Config")
		}

		if len(tlsConfig.Certificates) != 1 {
			t.Fatalf("Unexpected length in config.Certificates")
		}

		// Length should be 2, since we passed in a CA
		if len(tlsConfig.Certificates[0].Certificate) != 2 {
			t.Fatalf("Did not find both certificates in config.Certificates.Certificate")
		}

		if tlsConfig.Certificates[0].Leaf != pcbut.Certificate {
			t.Fatalf("Leaf certificate does not match parsed bundle's certificate")
		}

		if tlsConfig.Certificates[0].PrivateKey != pcbut.PrivateKey {
			t.Fatalf("Config's private key does not match parsed bundle's private key")
		}

		switch usage {
		case TLSServer | TLSClient:
			if len(tlsConfig.ClientCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.ClientCAs.Subjects()[0], pcbut.CAChain[0].Certificate.RawSubject) != 0 {
				t.Fatalf("CA certificate not in client cert pool as expected")
			}
			if len(tlsConfig.RootCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.RootCAs.Subjects()[0], pcbut.CAChain[0].Certificate.RawSubject) != 0 {
				t.Fatalf("CA certificate not in root cert pool as expected")
			}
		case TLSServer:
			if len(tlsConfig.ClientCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.ClientCAs.Subjects()[0], pcbut.CAChain[0].Certificate.RawSubject) != 0 {
				t.Fatalf("CA certificate not in client cert pool as expected")
			}
			if tlsConfig.RootCAs != nil {
				t.Fatalf("Found root pools in config object when not expected")
			}
		case TLSClient:
			if len(tlsConfig.RootCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.RootCAs.Subjects()[0], pcbut.CAChain[0].Certificate.RawSubject) != 0 {
				t.Fatalf("CA certificate not in root cert pool as expected")
			}
			if tlsConfig.ClientCAs != nil {
				t.Fatalf("Found root pools in config object when not expected")
			}
		default:
			if tlsConfig.RootCAs != nil || tlsConfig.ClientCAs != nil {
				t.Fatalf("Found root pools in config object when not expected")
			}
		}
	}
}

func refreshRSA8CertBundle() *CertBundle {
	return &CertBundle{
		Certificate: certRSAPem,
		PrivateKey:  privRSA8KeyPem,
		CAChain:     []string{issuingCaChainPem[0]},
	}
}

func refreshRSA8CertBundleWithChain() *CertBundle {
	ret := refreshRSA8CertBundle()
	ret.CAChain = issuingCaChainPem
	return ret
}

func refreshRSACertBundle() *CertBundle {
	return &CertBundle{
		Certificate: certRSAPem,
		CAChain:     []string{issuingCaChainPem[0]},
		PrivateKey:  privRSAKeyPem,
	}
}

func refreshRSACertBundleWithChain() *CertBundle {
	ret := refreshRSACertBundle()
	ret.CAChain = issuingCaChainPem
	return ret
}

func refreshECCertBundle() *CertBundle {
	return &CertBundle{
		Certificate: certECPem,
		CAChain:     []string{issuingCaChainPem[0]},
		PrivateKey:  privECKeyPem,
	}
}

func refreshECCertBundleWithChain() *CertBundle {
	ret := refreshECCertBundle()
	ret.CAChain = issuingCaChainPem
	return ret
}

func refreshRSACSRBundle() *CSRBundle {
	return &CSRBundle{
		CSR:        csrRSAPem,
		PrivateKey: privRSAKeyPem,
	}
}

func refreshECCSRBundle() *CSRBundle {
	return &CSRBundle{
		CSR:        csrECPem,
		PrivateKey: privECKeyPem,
	}
}

func refreshEC8CertBundle() *CertBundle {
	return &CertBundle{
		Certificate: certECPem,
		PrivateKey:  privEC8KeyPem,
		CAChain:     []string{issuingCaChainPem[0]},
	}
}

func refreshEC8CertBundleWithChain() *CertBundle {
	ret := refreshEC8CertBundle()
	ret.CAChain = issuingCaChainPem
	return ret
}

var (
	privRSA8KeyPem = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC92mr7+D/tGkW5
nvDH/fYkOLywbxsU9wU7lKVPCdj+zNzQYHiixTZtmZPYVTBj27lZgaUDUXuiw9Ru
BWHTuAb/Cpn1I+71qJbh8FgWot+MRDFKuV0PLkgHz5eRVC4JmKy9hbcgo1q0FfGf
qxL+VQmI0GcQ4IYK/ppVMrKbn4ndIg70uR46vPiU11GqRIz5wkiPeLklrhoWa5qE
IHINnR83eHUbijaCuqPEcz0QTz0iLM8dfubVaJ+Gn/DseUtku+qBdcYcUK2hQyCc
NRKtu953gr5hhEX0N9x5JBb10WaI1UL5HGa2wa6ndZ7yVb42B2WTHzNQRR5Cr4N4
ve31//gRAgMBAAECggEAfEvyvTLDz5zix2tK4vTfYMmQp8amKWysjVx9eijNW8yO
SRLQCGkrgEgLJphnjQk+6V3axjhjxKWHf9ygNrgGRJYRRBCZk1YkKpprYa6Sw0em
KfD//z9iw1JjPi+p0HiXp6FSytiIOt0fC1U6oy7ThjJDOCZ3O92C94KwsviZjx9r
DZbTLDm7Ya2LF4jGCq0dQ+AVqZ65QJ3yjdxm87PSE6q2eiV9wdMUx9RDOmFy+Meq
Mm3L9TW1QzyFtFMXeIF5QYGpmxWP/iii5V0CP573apXMIqQ+wTNpwK3WU5iURypZ
kJ1Iaxbzjfok6wpwLj7SJytF+fOVcygUxud7GPH8UQKBgQDPhQhB3+o+y+bwkUTx
Qdj/YNKcA/bjo/b9KMq+3rufwN9u/DK5z7vVfVklidbh5DVqhlLREsdSuZvb/IHc
OdCYwNeDxk1rLr+1W/iPYSBJod4eWDteIH1U9lts+/mH+u+iSsWVuikbeA8/MUJ3
nnAYu4FR1nz8I/CrvGbQL/KCdQKBgQDqNNI562Ch+4dJ407F3MC4gNPwPgksfLXn
ZRcPVVwGagil9oIIte0BIT0rAG/jVACfghGxfrj719uwjcFFxnUaSHGQcATseSf6
SgoruIVF15lI4e8lEcWrOypsW8Id2/amwUiIWYCgwlYG2Q7dggpXfgjmKfjSlvJ8
+yKR/Y6zrQKBgQCkx2aqICm5mWEUbtWGmJm9Ft3FQqSdV4n8tZJgAy6KiLUiRKHm
x1vIBtNtqkj1b6c2odhK6ZVaS8XF5XgcLdBEKwQ2P5Uj4agaUyBIgYAI174u7DKf
6D58423vWRln70qu3J6N6JdRl4DL1cqIf0dVbDYgjKcL82HcjCo7b4cqLQKBgFGU
TJX4MxS5NIq8LrglCMw7s5c/RJrGZeZQBBRHO2LQlGqazviRxhhap5/O6ypYHE9z
Uw5sgarXqaJ5/hR76FZbXZNeMZjdKtu35osMHwAQ9Ue5yz8yTZQza7eKzrbv4556
PPWhl3hnuOdxvAfUQB3xvM/PVuijw5tdLtGDbK2RAoGBAKB7OsTgF7wVEkzccJTE
hrbVKD+KBZz8WKnEgNoyyTIT0Kaugk15MCXkGrXIY8bW0IzYAw69qhTOgaWkcu4E
JbTK5UerP8V+X7XPBiw72StPVM4bxaXx2/B+78IuMOI/GR0tHQCF8S6DwTHeBXnl
ke8GFExnXHTPqG6Bku0r/G47
-----END PRIVATE KEY-----`

	privRSAKeyPem = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAvdpq+/g/7RpFuZ7wx/32JDi8sG8bFPcFO5SlTwnY/szc0GB4
osU2bZmT2FUwY9u5WYGlA1F7osPUbgVh07gG/wqZ9SPu9aiW4fBYFqLfjEQxSrld
Dy5IB8+XkVQuCZisvYW3IKNatBXxn6sS/lUJiNBnEOCGCv6aVTKym5+J3SIO9Lke
Orz4lNdRqkSM+cJIj3i5Ja4aFmuahCByDZ0fN3h1G4o2grqjxHM9EE89IizPHX7m
1Wifhp/w7HlLZLvqgXXGHFCtoUMgnDUSrbved4K+YYRF9DfceSQW9dFmiNVC+Rxm
tsGup3We8lW+Ngdlkx8zUEUeQq+DeL3t9f/4EQIDAQABAoIBAHxL8r0yw8+c4sdr
SuL032DJkKfGpilsrI1cfXoozVvMjkkS0AhpK4BICyaYZ40JPuld2sY4Y8Slh3/c
oDa4BkSWEUQQmZNWJCqaa2GuksNHpinw//8/YsNSYz4vqdB4l6ehUsrYiDrdHwtV
OqMu04YyQzgmdzvdgveCsLL4mY8faw2W0yw5u2GtixeIxgqtHUPgFameuUCd8o3c
ZvOz0hOqtnolfcHTFMfUQzphcvjHqjJty/U1tUM8hbRTF3iBeUGBqZsVj/4oouVd
Aj+e92qVzCKkPsEzacCt1lOYlEcqWZCdSGsW8436JOsKcC4+0icrRfnzlXMoFMbn
exjx/FECgYEAz4UIQd/qPsvm8JFE8UHY/2DSnAP246P2/SjKvt67n8Dfbvwyuc+7
1X1ZJYnW4eQ1aoZS0RLHUrmb2/yB3DnQmMDXg8ZNay6/tVv4j2EgSaHeHlg7XiB9
VPZbbPv5h/rvokrFlbopG3gPPzFCd55wGLuBUdZ8/CPwq7xm0C/ygnUCgYEA6jTS
OetgofuHSeNOxdzAuIDT8D4JLHy152UXD1VcBmoIpfaCCLXtASE9KwBv41QAn4IR
sX64+9fbsI3BRcZ1GkhxkHAE7Hkn+koKK7iFRdeZSOHvJRHFqzsqbFvCHdv2psFI
iFmAoMJWBtkO3YIKV34I5in40pbyfPsikf2Os60CgYEApMdmqiApuZlhFG7VhpiZ
vRbdxUKknVeJ/LWSYAMuioi1IkSh5sdbyAbTbapI9W+nNqHYSumVWkvFxeV4HC3Q
RCsENj+VI+GoGlMgSIGACNe+Luwyn+g+fONt71kZZ+9KrtyejeiXUZeAy9XKiH9H
VWw2IIynC/Nh3IwqO2+HKi0CgYBRlEyV+DMUuTSKvC64JQjMO7OXP0SaxmXmUAQU
Rzti0JRqms74kcYYWqefzusqWBxPc1MObIGq16mief4Ue+hWW12TXjGY3Srbt+aL
DB8AEPVHucs/Mk2UM2u3is627+Oeejz1oZd4Z7jncbwH1EAd8bzPz1boo8ObXS7R
g2ytkQKBgQCgezrE4Be8FRJM3HCUxIa21Sg/igWc/FipxIDaMskyE9CmroJNeTAl
5Bq1yGPG1tCM2AMOvaoUzoGlpHLuBCW0yuVHqz/Ffl+1zwYsO9krT1TOG8Wl8dvw
fu/CLjDiPxkdLR0AhfEug8Ex3gV55ZHvBhRMZ1x0z6hugZLtK/xuOw==
-----END RSA PRIVATE KEY-----`

	csrRSAPem = `-----BEGIN CERTIFICATE REQUEST-----
MIICijCCAXICAQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUx
ITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBALd2SVGs7QkYlYPOz9MvE+0DDRUIutQQkiCnqcV1
1nO84uSSUjH0ALLVyRKgWIkobEqWJurKIk+oV9O+Le8EABxkt/ru6jaIFqZkwFE1
JzSDPkAsR9jXP5M2MWlmo7qVc4Bry3oqlN3eLSFJP2ZH7b1ia8q9oWXPMxDdwuKR
kv9hfkUPszr/gaQCNQXW0BoRe5vdr5+ikv4lrKpsBxvaAoYL1ngR41lCPKhmrgXP
oreEuUcXzfCSSAV1CGbW2qhWc4I7/JFA8qqEBr5GP+AntStGxSbDO8JjD7/uULHC
AReCloDdJE0jDz1355/0CAH4WmEJE/TI8Bq+vgd0jgzRgk0CAwEAAaAAMA0GCSqG
SIb3DQEBCwUAA4IBAQAR8U1vZMJf7YFvGU69QvoWPTDe/o8SwYy1j+++AAO9Y7H2
C7nb+9tnEMtXm+3pkY0aJIecAnq8H4QWimOrJa/ZsoZLzz9LKW2nzARdWo63j4nB
jKld/EDBzQ/nQSTyoX7s9JiDiSC9yqTXBrPHSXruPbh7sE0yXROar+6atjNdCpDp
uLw86gwewDJrMaB1aFAmDvwaRQQDONwRy0zG1UdMxLQxsxpKOHaGM/ZvV3FPir2B
7mKupki/dvap5UW0lTMJBlKf3qhoeHKMHFo9i5vGCIkWUIv+XgTF0NjbYv9i7bfq
WdW905v4wiuWRlddNwqFtLx9Pf1/fRJVT5mBbjIx
-----END CERTIFICATE REQUEST-----`

	certRSAPem = `-----BEGIN CERTIFICATE-----
MIIDfDCCAmSgAwIBAgIUad4Q9EhVvqc06H7fCfKaLGcyDw0wDQYJKoZIhvcNAQEL
BQAwNzE1MDMGA1UEAxMsVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3ViIFN1
YiBBdXRob3JpdHkwHhcNMTYwODA0MTkyMjAyWhcNMTYwODA0MjAyMjMyWjAhMR8w
HQYDVQQDExZWYXVsdCBUZXN0IENlcnRpZmljYXRlMIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAvdpq+/g/7RpFuZ7wx/32JDi8sG8bFPcFO5SlTwnY/szc
0GB4osU2bZmT2FUwY9u5WYGlA1F7osPUbgVh07gG/wqZ9SPu9aiW4fBYFqLfjEQx
SrldDy5IB8+XkVQuCZisvYW3IKNatBXxn6sS/lUJiNBnEOCGCv6aVTKym5+J3SIO
9LkeOrz4lNdRqkSM+cJIj3i5Ja4aFmuahCByDZ0fN3h1G4o2grqjxHM9EE89IizP
HX7m1Wifhp/w7HlLZLvqgXXGHFCtoUMgnDUSrbved4K+YYRF9DfceSQW9dFmiNVC
+RxmtsGup3We8lW+Ngdlkx8zUEUeQq+DeL3t9f/4EQIDAQABo4GVMIGSMA4GA1Ud
DwEB/wQEAwIDqDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwHQYDVR0O
BBYEFMKLVTrdDyyF0kTxkxcAMcGNcGWlMB8GA1UdIwQYMBaAFNmGqFL215GlYyD0
mWVIoMB71s+NMCEGA1UdEQQaMBiCFlZhdWx0IFRlc3QgQ2VydGlmaWNhdGUwDQYJ
KoZIhvcNAQELBQADggEBAJJP9OWG3W5uUluKdeFYCzKMIY+rsCUb86QrKRqQ5xYR
w4pKC3yuryEfreBs3iQA4NNw2mMWxuI8t/i+km2H7NzQytTRn6L0sxTa8ThNZ3e7
xCdWaZZzd1O6Xwq/pDbE1MZ/4z5nvsKaKJVVIvVFL5algi4A8njiFMVSww035c1e
waLww4AOHydlLky/RJBJPOkQNoDBToC9ojDqPtNJVWWaQL2TsUCu+Q+L5QL5djgj
LxPwqGOiM4SLSUrXSXMpHNLX1rhBH1/sNb3Kn1FDBaZ+M9kZglCDwuQyQuH8xKwB
qukeKfgFUp7rH0yoQTZa0eaXAYTFoRLjnTQ+fS7e19s=
-----END CERTIFICATE-----`

	privECKeyPem = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEICC2XihYLxEYEseFesZEXjV1va6rMAdtkpkaxT4hGu5boAoGCCqGSM49
AwEHoUQDQgAEti0uWkq7MAkQevNNrBpYY0FLni8OAZroHXkij2x6Vo0xIvClftbC
L33BU/520t23TcewtQYsNqv86Bvhx9PeAw==
-----END EC PRIVATE KEY-----`

	csrECPem = `-----BEGIN CERTIFICATE REQUEST-----
MIHsMIGcAgEAMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwTjAQBgcqhkjOPQIBBgUr
gQQAIQM6AATBZ3VXwBE9oeSREpM5b25PW6WiuLb4EXWpKZyjj552QYKYe7QBuGe9
wvvgOeCBovN3tSuGKzTiUKAAMAoGCCqGSM49BAMCAz8AMDwCHFap/5XDuqtXCG1g
ljbYH5OWGBqGYCfL2k2+/6cCHAuk1bmOkGx7JAq/fSPd09i0DQIqUu7WHQHms48=
-----END CERTIFICATE REQUEST-----`

	privEC8KeyPem = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgILZeKFgvERgSx4V6
xkReNXW9rqswB22SmRrFPiEa7luhRANCAAS2LS5aSrswCRB6802sGlhjQUueLw4B
mugdeSKPbHpWjTEi8KV+1sIvfcFT/nbS3bdNx7C1Biw2q/zoG+HH094D
-----END PRIVATE KEY-----`

	certECPem = `-----BEGIN CERTIFICATE-----
MIICtzCCAZ+gAwIBAgIUNDYMWd9SOGVMs4I1hezvRnGDMyUwDQYJKoZIhvcNAQEL
BQAwNzE1MDMGA1UEAxMsVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3ViIFN1
YiBBdXRob3JpdHkwHhcNMTYwODA0MTkyOTM0WhcNMTYwODA0MjAzMDA0WjAkMSIw
IAYDVQQDExlWYXVsdCBUZXN0IEVDIENlcnRpZmljYXRlMFkwEwYHKoZIzj0CAQYI
KoZIzj0DAQcDQgAEti0uWkq7MAkQevNNrBpYY0FLni8OAZroHXkij2x6Vo0xIvCl
ftbCL33BU/520t23TcewtQYsNqv86Bvhx9PeA6OBmDCBlTAOBgNVHQ8BAf8EBAMC
A6gwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB0GA1UdDgQWBBStnbW/
ga2/dz4FyRNafwhTzM1UbzAfBgNVHSMEGDAWgBTZhqhS9teRpWMg9JllSKDAe9bP
jTAkBgNVHREEHTAbghlWYXVsdCBUZXN0IEVDIENlcnRpZmljYXRlMA0GCSqGSIb3
DQEBCwUAA4IBAQBsPhwRB51de3sGBMnjDiOMViYpRH7kKhUWAY1W2W/1hqk5HgZw
4c3r0LmdIQ94gShaXng8ojYRDW/5D7LeXJdbtLy9U29xfeCb+vqKDc2oN7Ps3/HB
4YLnseqDiZFKPEAdOE4rtwyFWJI7JR9sOSG1B5El6duN0i9FWOLSklQ4EbV5R45r
cy/fJq0DOYje7MXsFuNl5iQ92gfDjPD2P98DK9lCIquSzB3WkpjE41UtKJ0IKPeD
wYoyl0J33Alxq2eC2subR7xISR3MzZFcdkzNNrBddeaSviYlR4SgTUiqOldAcdR4
QZxtxazcUqQDZ+wZFOpBOnp94bzVeXT9BF+L
-----END CERTIFICATE-----`

	issuingCaChainPem = []string{`-----BEGIN CERTIFICATE-----
MIIDljCCAn6gAwIBAgIUHjciEzUzeNVqI9mwFJeduNtXWzMwDQYJKoZIhvcNAQEL
BQAwMzExMC8GA1UEAxMoVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3ViIEF1
dGhvcml0eTAeFw0xNjA4MDQxOTEyNTdaFw0xNjA4MDUyMDEzMjdaMDcxNTAzBgNV
BAMTLFZhdWx0IFRlc3RpbmcgSW50ZXJtZWRpYXRlIFN1YiBTdWIgQXV0aG9yaXR5
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy2pAH1U8KzhjO+MLRPTb
ic7Iyk57d0TFnj6CAJWqZaKNGXoTkwD8wRCirY8mQv8YrfBy3hwGqSLYj6oxwA0R
8FxsiWdf4gFTX2cJpxThFnIllGbzqIXnEZLvCIMydp44Ls9eYxoXfZQ9X24u/Wmf
kWEQFGUzrpyklkIOx2Yo5g7OHbFLl3OfPz89/TDM8VeymlGzCTJZ+Y+iNGDBPT0L
X9aE65lL76dUx/bcKnfQEgAcH4nkE4K/Kgjnj5umZKQUH4+6wKFwDCQT2RwaBkve
WyAiz0LY9a1WFXt7RYCPs+QWLJAhv7wJL8l4gnxYA1k+ovLXDjUqYweU+WHV6/lR
7wIDAQABo4GdMIGaMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0G
A1UdDgQWBBTZhqhS9teRpWMg9JllSKDAe9bPjTAfBgNVHSMEGDAWgBRTY6430DXg
cIDAnEnA+fostjcbqDA3BgNVHREEMDAugixWYXVsdCBUZXN0aW5nIEludGVybWVk
aWF0ZSBTdWIgU3ViIEF1dGhvcml0eTANBgkqhkiG9w0BAQsFAAOCAQEAZp3VwrUw
jw6TzjZJEnXkfRxLWZNmqMPAXYOtBl/+5FjAfehifTTzIdIxR4mfdgH5YZnSQpzY
m/w17BXElao8uOX6CUaX+sLTVzwsl2csswpcGlNwHooVREoMq9X187qxSr1HS7zF
O550XgDVIf5e7sXrVuV1rd1XUo3xZLaSLUhU70y/343mcN2TRUslXO4QrIE5lo2v
awyQl0NW0hSO0F9VZYzOvPPVwu7mf1ijTzbkPtUbAXDnmlvOCrlx2JZd/BqXb75e
UgYDq7hIyQ109FBOjv0weAM5tZCdesyvro4/43Krd8pa74zHdZMjfQAsTr66WOi4
yedj8LnWl66JOA==
-----END CERTIFICATE-----`,
		`-----BEGIN CERTIFICATE-----
MIIDijCCAnKgAwIBAgIUBNDYCUsOT2Wth8Fz3layfjEVbcIwDQYJKoZIhvcNAQEL
BQAwLzEtMCsGA1UEAxMkVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgQXV0aG9y
aXR5MB4XDTE2MDgwNDE5MTI1NloXDTE2MDgwNjIxMTMyNlowMzExMC8GA1UEAxMo
VmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3ViIEF1dGhvcml0eTCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBALHoD7g5YYu2akO8hkFlUCF45Bxjckq4
WTyDIcDwv/wr7vhZnngCClnP+7Rc30XTmkq2RnH6N7iuqowGM5RNcBV/C9R1weVx
9esXtWr/AUMyuNb3HSjwDwQGuiAVEgk67fXYy08Ii78+ap3uY3CKC1AFDkHdgDZt
e946rJ3Nps00TcH0KwyP5voitLgt6dMBR9ttuUdSoQ4uLQDdDf0HRw/IAQswO4Av
lgUgQObBecnLGhh7e3PM5VVz5f0IqG2ZYnDs3ncl2UYOrj0/JqOMDIMvSQMc2pzS
Hjty0d1wKWWPC9waguL/24oQR4VG5b7TL62elc2kcEg7r8u5L/sCi/8CAwEAAaOB
mTCBljAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU
U2OuN9A14HCAwJxJwPn6LLY3G6gwHwYDVR0jBBgwFoAUgAz80p6Pkzk6Cb7lYmTI
T1jc7iYwMwYDVR0RBCwwKoIoVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgU3Vi
IEF1dGhvcml0eTANBgkqhkiG9w0BAQsFAAOCAQEACXjzGVqRZi6suHwBqvXwVxlS
nf5YwudfBDJ4LfNd5nTCypsvHkfXnaki6LZMCS1rwPvxzssZ5Wp/7zO5fu6lpSTx
yjuiH5fBUGHy+f1Ygu6tlAZtUnxAi6pU4eoCDNZpqunJMM4IdaahHeICdjPhx/bH
AlmwaN0FsNvOlgUuPTjQ3z6jMZn3p2lXI3HiRlcz+nR7gQizPb2L7u8mQ+5EZFmC
AmXMj40g3bTJVmKoGeAR7cb0pYG/GUELmERjEjCfP7W15eYfuu1j7EYTUAVuPAlJ
34HDxCuM8cPJwCGMDKfb3Q39AYRmLT6sE3/sq2CZ5xlj8wfwDpVfpXikRDpI0A==
-----END CERTIFICATE-----`,
		`-----BEGIN CERTIFICATE-----
MIIDejCCAmKgAwIBAgIUEtjlbdzIth3U71TELA0PVW7HvaEwDQYJKoZIhvcNAQEL
BQAwJzElMCMGA1UEAxMcVmF1bHQgVGVzdGluZyBSb290IEF1dGhvcml0eTAeFw0x
NjA4MDQxOTEyNTVaFw0xNjA4MDgyMzEzMjVaMC8xLTArBgNVBAMTJFZhdWx0IFRl
c3RpbmcgSW50ZXJtZWRpYXRlIEF1dGhvcml0eTCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBAMYAQAHCm9V9062NF/UuAa6z6aYqsS5g2YGkd9DvgYxfU5JI
yIdSz7rkp9QprlQYl2abptZocq+1C9yRVmRJWKjZYDckSwXdmQam/sOfNuiw6Gbd
3OJGdQ82jhx3v3mIQp+3u9E43wXX0StaJ44+9DgkgwG8iybiv4fh0LzuHPSeKsXe
/IvJZ0YAInWuzFNegYxU32UT2CEvLtZdru8+sLr4NFWRu/nYIMPJDeZ2JEQVi9IF
lcB3dP63c6vMBrn4Wn2xBo12JPsQp+ezf5Z5zmtAe68PwRmIXZVAUa2q+CfEuJ36
66756Ypa0Z3brhPWfX2ahhxSg8DjqFGmZZ5Gfl8CAwEAAaOBlTCBkjAOBgNVHQ8B
Af8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUgAz80p6Pkzk6Cb7l
YmTIT1jc7iYwHwYDVR0jBBgwFoAU6dC1U32HZp7iq97KSu2i+g8+rf4wLwYDVR0R
BCgwJoIkVmF1bHQgVGVzdGluZyBJbnRlcm1lZGlhdGUgQXV0aG9yaXR5MA0GCSqG
SIb3DQEBCwUAA4IBAQA6xVMyuZgpJhIQnG2LwwD5Zcbmm4+rHkGVNSpwkUH8ga8X
b4Owog+MvBw8R7ADVwAh/aOh1/qsRfHv8KkMWW+SAQ84ICVXJiPBzEUJaMWujpyr
SDkbaD05avRtfvSrPCagaUGVRt+wK24g8hpJqQ+trkufzjq9ySU018+NNX9yGRyA
VjwZAqALlNEAkdcvd4adEBpZqum2x1Fl9EXnjp6NEWQ7nuGkp3X2DP4gDtQPxgmn
omOo4GHhO0U57exEIl0d4kiy9WU0qcIISOr6I+gzesMooX6aI43CaqJoZKsHXYY6
1uxFLss+/wDtvIcyXdTdjPrgD38YIgk1/iKNIgKO
-----END CERTIFICATE-----`}
)
