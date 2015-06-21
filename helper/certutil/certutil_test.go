package certutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/api"
)

// Tests converting back and forth between a CertBundle and a ParsedCertBundle.
//
// Also tests the GetSubjKeyID, GetOctalFormatted, and
// ParsedCertBundle.getSigner functions.
func TestCertBundleConversion(t *testing.T) {
	cbuts := []*CertBundle{
		refreshRSACertBundle(),
		refreshECCertBundle(),
	}

	for _, cbut := range cbuts {
		pcbut, err := cbut.ToParsedCertBundle()
		if err != nil {
			t.Fatalf("Error converting to parsed cert bundle: %s", err)
		}

		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
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
			if len(tlsConfig.ClientCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.ClientCAs.Subjects()[0], pcbut.IssuingCA.RawSubject) != 0 {
				t.Fatalf("CA certificate not in client cert pool as expected")
			}
			if len(tlsConfig.RootCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.RootCAs.Subjects()[0], pcbut.IssuingCA.RawSubject) != 0 {
				t.Fatalf("CA certificate not in root cert pool as expected")
			}
		case TLSServer:
			if len(tlsConfig.ClientCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.ClientCAs.Subjects()[0], pcbut.IssuingCA.RawSubject) != 0 {
				t.Fatalf("CA certificate not in client cert pool as expected")
			}
			if tlsConfig.RootCAs != nil {
				t.Fatalf("Found root pools in config object when not expected")
			}
		case TLSClient:
			if len(tlsConfig.RootCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.RootCAs.Subjects()[0], pcbut.IssuingCA.RawSubject) != 0 {
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

func TestCertBundleParsing(t *testing.T) {
	jsonBundle := refreshRSACertBundle()
	jsonString, err := json.Marshal(jsonBundle)
	if err != nil {
		t.Fatalf("Error marshaling testing certbundle to JSON: %s", err)
	}
	pcbut, err := ParsePKIJSON(jsonString)
	if err != nil {
		t.Fatalf("Error during JSON bundle handling: %s", err)
	}
	err = compareCertBundleToParsedCertBundle(jsonBundle, pcbut)
	if err != nil {
		t.Fatalf(err.Error())
	}

	secret := &api.Secret{
		Data: structs.New(jsonBundle).Map(),
	}
	pcbut, err = ParsePKIMap(secret.Data)
	if err != nil {
		t.Fatalf("Error during JSON bundle handling: %s", err)
	}
	err = compareCertBundleToParsedCertBundle(jsonBundle, pcbut)
	if err != nil {
		t.Fatalf(err.Error())
	}

	pemBundle := strings.Join([]string{
		jsonBundle.Certificate,
		jsonBundle.IssuingCA,
		jsonBundle.PrivateKey,
	}, "\n")
	pcbut, err = ParsePEMBundle(pemBundle)
	if err != nil {
		t.Fatalf("Error during JSON bundle handling: %s", err)
	}
	err = compareCertBundleToParsedCertBundle(jsonBundle, pcbut)
	if err != nil {
		t.Fatalf(err.Error())
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
	case pcbut.IssuingCA == nil:
		return fmt.Errorf("Parsed bundle has nil issuing CA")
	}

	switch cbut.PrivateKey {
	case privRSAKeyPem:
		if pcbut.PrivateKeyType != RSAPrivateKey {
			return fmt.Errorf("Parsed bundle has wrong private key type")
		}
	case privECKeyPem:
		if pcbut.PrivateKeyType != ECPrivateKey {
			return fmt.Errorf("Parsed bundle has wrong private key type")
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

	cb, err := pcbut.ToCertBundle()
	if err != nil {
		return fmt.Errorf("Thrown error during parsed bundle conversion: %s\n\nInput was: %#v", err, *pcbut)
	}

	switch {
	case len(cb.Certificate) == 0:
		return fmt.Errorf("Bundle has nil certificate")
	case len(cb.PrivateKey) == 0:
		return fmt.Errorf("Bundle has nil private key")
	case len(cb.IssuingCA) == 0:
		return fmt.Errorf("Bundle has nil issuing CA")
	}

	switch cb.PrivateKeyType {
	case "rsa":
		if pcbut.PrivateKeyType != RSAPrivateKey {
			return fmt.Errorf("Bundle has wrong private key type")
		}
		if cb.PrivateKey != privRSAKeyPem {
			return fmt.Errorf("Bundle private key does not match")
		}
	case "ec":
		if pcbut.PrivateKeyType != ECPrivateKey {
			return fmt.Errorf("Bundle has wrong private key type")
		}
		if cb.PrivateKey != privECKeyPem {
			return fmt.Errorf("Bundle private key does not match")
		}
	default:
		return fmt.Errorf("Bundle has unknown private key type")
	}

	if cb.SerialNumber != GetOctalFormatted(pcbut.Certificate.SerialNumber.Bytes(), ":") {
		return fmt.Errorf("Bundle serial number does not match")
	}

	return nil
}

func refreshRSACertBundle() *CertBundle {
	return &CertBundle{
		Certificate: certRSAPem,
		PrivateKey:  privRSAKeyPem,
		IssuingCA:   issuingCaPem,
	}
}

func refreshECCertBundle() *CertBundle {
	return &CertBundle{
		Certificate: certECPem,
		PrivateKey:  privECKeyPem,
		IssuingCA:   issuingCaPem,
	}
}

const (
	privRSAKeyPem = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAt3ZJUaztCRiVg87P0y8T7QMNFQi61BCSIKepxXXWc7zi5JJS
MfQAstXJEqBYiShsSpYm6soiT6hX074t7wQAHGS3+u7qNogWpmTAUTUnNIM+QCxH
2Nc/kzYxaWajupVzgGvLeiqU3d4tIUk/ZkftvWJryr2hZc8zEN3C4pGS/2F+RQ+z
Ov+BpAI1BdbQGhF7m92vn6KS/iWsqmwHG9oChgvWeBHjWUI8qGauBc+it4S5RxfN
8JJIBXUIZtbaqFZzgjv8kUDyqoQGvkY/4Ce1K0bFJsM7wmMPv+5QscIBF4KWgN0k
TSMPPXfnn/QIAfhaYQkT9MjwGr6+B3SODNGCTQIDAQABAoIBAHtSwhprGbNhmS/P
F5ioLsbFpEedZKkksnXM/qxDd/K45/Qp/6KgmM+eMdmZe6pHR/QjVunBEqtlSBSH
5KykjcaIVbwSWdJqTH9xfm2YQ1BjYLcWjP1QQ+YbKb/mRO0phUiwLUlj0koKDWAw
srN4anFB9Z+FNTcQvwz5ZQWUQbH0neQtWO1nDvLsScgu1kchoEzJEJaFOQ1+HfGe
WxD766fZyqZQi5+cLrhOqHOGSlO+IFVe0hguiEHFr9LEPTXXkZtOR4wTf7j1Us8s
1KQ/jv01sx9S7HEbZJurzIjS23OywEUdJd1EsIE2lJV2QUwSiAsPYZOSQZlgOGzP
VRKVkGkCgYEA1u+pVP2r+xSxYy8KcdcRCdGGBh00VLx1yJRHWZ5YjF56hp0R0cG+
xGLar5KCdBpr4jJnQGIrx8lw3SDCt4EXlxgJxitXlBtiKByM7/mYRRfURr9WMRr4
88GQlWDbo2Xalnuac0qlkFqVIg0BaW+Z15A/E1L69aUxaR0ozlA9Jl8CgYEA2oNA
5F2otqzo9eNYucNAjihVhATd11DECQvbIQp/0bEJe0Znnzq/QIGIOVapC0VKGBwB
P5DuLL1P/nTPjjE/ZhjFuhMNM5PzC6obAjBh+gCpc+c+21Qerv7RKUTi2sGTzRHu
lpccRDfuF8bhzD6lAo50FpSmPE/ovZzb9+IsXtMCgYBVnUdM9HKh47846870Q5+k
0pHZM57ZtewQxoeZOgq5dxTFNCGZ9NvBLENBtlFCYBfjFQKt0azwutu7KUaGg+Ra
qheSmUccVsAFjEHTgQ9XTkOfHq39h2ns5ohqCBfVAUhNstR14iEK3BoVYyrRzcNw
6yNE1kPivzdsUFIlxC5nbwKBgDUUjT7sQX+eoTiZ8YOumo/t3Fglln4ncHeCGcj8
8+/MQbFgeOuFKdBRpvXGx2mle0pAA02dtz3G/xeg6IpyDCSQ//cjiaFt3yyGNeli
N2qznnY5RluhI5L+83BC+5iITY8TPBH4wzUPIRdFiLREw3DLigeyNG+SOcdVw1mD
56NhAoGBALFh3sGkhvPiI/G/i/5tGZVA/dS/4DVXOoHW43+ZDHWEwqiN6vTf/VVi
cm+8kcfLY1E5fSf/4e7mIQq7o5qVn9Y3HWsajS1FFeznJjPj4Jaa1HvegNcycAzs
XOQ7xy23/8wUupgNeD1mFdSFCXQ3UedsJuVBHsElPc5W74q4F4+F
-----END RSA PRIVATE KEY-----`

	certRSAPem = `-----BEGIN CERTIFICATE-----
MIID+jCCAuSgAwIBAgIUcFCL9ESWTKLE6RqSYV7iZ78f1KcwCwYJKoZIhvcNAQEL
MBsxGTAXBgNVBAMMEFZhdWx0IFRlc3RpbmcgQ0EwHhcNMTUwNjE5MTcyMzA0WhcN
MTUwNzAzMTcyMzA0WjBPMRIwEAYDVQQDEwlsb2NhbGhvc3QxOTA3BgNVBAUTMDY0
MTIwMzIxNzY3NTk2MjQyMjU0OTg5MTUxMzAyMjg1NzQ0NTc0OTkzMjY3NjI2MzCC
ASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALd2SVGs7QkYlYPOz9MvE+0D
DRUIutQQkiCnqcV11nO84uSSUjH0ALLVyRKgWIkobEqWJurKIk+oV9O+Le8EABxk
t/ru6jaIFqZkwFE1JzSDPkAsR9jXP5M2MWlmo7qVc4Bry3oqlN3eLSFJP2ZH7b1i
a8q9oWXPMxDdwuKRkv9hfkUPszr/gaQCNQXW0BoRe5vdr5+ikv4lrKpsBxvaAoYL
1ngR41lCPKhmrgXPoreEuUcXzfCSSAV1CGbW2qhWc4I7/JFA8qqEBr5GP+AntStG
xSbDO8JjD7/uULHCAReCloDdJE0jDz1355/0CAH4WmEJE/TI8Bq+vgd0jgzRgk0C
AwEAAaOCAQQwggEAMA4GA1UdDwEB/wQEAwIAqDAdBgNVHSUEFjAUBggrBgEFBQcD
AQYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQUZHtkxSX5GVAYo3h8
B8TGJ36vTH4wHwYDVR0jBBgwFoAU5JzeXhccaOWk5X6vhuGV7NwLMVkwTgYDVR0R
BEcwRYIJbG9jYWxob3N0gg9mb28uZXhhbXBsZS5jb22CD2Jhci5leGFtcGxlLmNv
bYcEgAMFBocQ/gEAAAAAAAAAAAAAAAAAATAxBgNVHR8EKjAoMCagJKAihiBodHRw
Oi8vbG9jYWxob3N0OjgyMDAvdjEvcGtpL2NybDALBgkqhkiG9w0BAQsDggEBAAps
W2ZDOAfwWufclmGPHt+YRXXSTWvPfF/cBeg5Oq/F8qUCVMHqdE/+EDWzh+Kz8jp0
ggklnh76frROvHxygbVD2Hs9ACzgpnHPy8FYOdN+OblvAMtGlMyTq/5XheasmWdY
FFH/ft6tReG7BjGgfdyH8yL/R6b/RtU/qPlowfrZgAzOv7/ou6yRlfjIhsWbne/S
SQuGASRxRp3Txp7Cf3RcdCwVuiQhFLVeVHH+atTc8v2DO/CLfi9enQo96qUku8Bd
b5QPKIV0sQdtwGV5fo2JGd25rWpCo6TkAM9EeNkcVze8wgArSRk8zLkvM/5z+5sn
Qaka08px4wljGQ2Wc88=
-----END CERTIFICATE-----`

	privECKeyPem = `-----BEGIN EC PRIVATE KEY-----
MGgCAQEEHM3nuYLlrvawBN9hGVcu9mpaCEr7LMe44a7oQOygBwYFK4EEACGhPAM6
AATBZ3VXwBE9oeSREpM5b25PW6WiuLb4EXWpKZyjj552QYKYe7QBuGe9wvvgOeCB
ovN3tSuGKzTiUA==
-----END EC PRIVATE KEY-----`

	certECPem = `-----BEGIN CERTIFICATE-----
MIIDJDCCAg6gAwIBAgIUM3J02tw0ZvpHUVHv6t8kcoft2/MwCwYJKoZIhvcNAQEL
MBsxGTAXBgNVBAMMEFZhdWx0IFRlc3RpbmcgQ0EwHhcNMTUwNjE5MTcyODQyWhcN
MTUwNzAzMTcyODQyWjBPMRIwEAYDVQQDEwlsb2NhbGhvc3QxOTA3BgNVBAUTMDI5
MzcxMDk5Mzc2NDA3NDYyNjg3MTQzODcwMjc3Njg1OTkzMTkyMzkxNjM4MTE3MTBO
MBAGByqGSM49AgEGBSuBBAAhAzoABMFndVfAET2h5JESkzlvbk9bpaK4tvgRdakp
nKOPnnZBgph7tAG4Z73C++A54IGi83e1K4YrNOJQo4IBBDCCAQAwDgYDVR0PAQH/
BAQDAgCoMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8E
AjAAMB0GA1UdDgQWBBQiFoWDvInznUGjdJPjBAyoxIkQITAfBgNVHSMEGDAWgBTk
nN5eFxxo5aTlfq+G4ZXs3AsxWTBOBgNVHREERzBFgglsb2NhbGhvc3SCD2Zvby5l
eGFtcGxlLmNvbYIPYmFyLmV4YW1wbGUuY29thwSAAwUGhxD+AQAAAAAAAAAAAAAA
AAABMDEGA1UdHwQqMCgwJqAkoCKGIGh0dHA6Ly9sb2NhbGhvc3Q6ODIwMC92MS9w
a2kvY3JsMAsGCSqGSIb3DQEBCwOCAQEA0RU18OdSdt2k4FKWyUS7EhVFOybiUHof
1n9EeBoxd7fEP/IuQnJGr3CPV5LRFdHRxkihf4N5bRjsst7cqczaIZZLWkAj+P/2
JxBqv2Hm57dwaw2gtwt3GcYN/5j76fYaoZOgPMqas72vYgnBgdKQs8GYSoy7BVpC
x3nTYHwlOF+sM4wuVSi78lwkcgADF5GIWXrM3tYilmcT9fNbUgSvcVWdNTRJ0W+m
S2AF+4eby5PC9U8eIoCnZPRNmH0jZbNWzZyD0hDhBrDlaEbS2QXKRURPHzht/SqN
nWWcpQG3B8EI7p749dP5L+idi3ajHIH8vm/PK+o5TRrcHB585MlErQ==
-----END CERTIFICATE-----`

	issuingCaPem = `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`
)
