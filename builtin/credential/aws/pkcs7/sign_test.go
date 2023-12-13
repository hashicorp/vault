package pkcs7

import (
	"bytes"
	"crypto/dsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"testing"
)

func TestSign(t *testing.T) {
	content := []byte("Hello World")
	sigalgs := []x509.SignatureAlgorithm{
		x509.SHA256WithRSA,
		x509.SHA512WithRSA,
		x509.ECDSAWithSHA256,
		x509.ECDSAWithSHA384,
		x509.ECDSAWithSHA512,
	}
	for _, sigalgroot := range sigalgs {
		rootCert, err := createTestCertificateByIssuer("PKCS7 Test Root CA", nil, sigalgroot, true)
		if err != nil {
			t.Fatalf("test %s: cannot generate root cert: %s", sigalgroot, err)
		}
		truststore := x509.NewCertPool()
		truststore.AddCert(rootCert.Certificate)
		for _, sigalginter := range sigalgs {
			interCert, err := createTestCertificateByIssuer("PKCS7 Test Intermediate Cert", rootCert, sigalginter, true)
			if err != nil {
				t.Fatalf("test %s/%s: cannot generate intermediate cert: %s", sigalgroot, sigalginter, err)
			}
			var parents []*x509.Certificate
			parents = append(parents, interCert.Certificate)
			for _, sigalgsigner := range sigalgs {
				signerCert, err := createTestCertificateByIssuer("PKCS7 Test Signer Cert", interCert, sigalgsigner, false)
				if err != nil {
					t.Fatalf("test %s/%s/%s: cannot generate signer cert: %s", sigalgroot, sigalginter, sigalgsigner, err)
				}
				for _, testDetach := range []bool{false, true} {
					log.Printf("test %s/%s/%s detached %t\n", sigalgroot, sigalginter, sigalgsigner, testDetach)
					toBeSigned, err := NewSignedData(content)
					if err != nil {
						t.Fatalf("test %s/%s/%s: cannot initialize signed data: %s", sigalgroot, sigalginter, sigalgsigner, err)
					}

					// Set the digest to match the end entity cert
					signerDigest, _ := getDigestOIDForSignatureAlgorithm(signerCert.Certificate.SignatureAlgorithm)
					toBeSigned.SetDigestAlgorithm(signerDigest)

					if err := toBeSigned.AddSignerChain(signerCert.Certificate, *signerCert.PrivateKey, parents, SignerInfoConfig{}); err != nil {
						t.Fatalf("test %s/%s/%s: cannot add signer: %s", sigalgroot, sigalginter, sigalgsigner, err)
					}
					if testDetach {
						toBeSigned.Detach()
					}
					signed, err := toBeSigned.Finish()
					if err != nil {
						t.Fatalf("test %s/%s/%s: cannot finish signing data: %s", sigalgroot, sigalginter, sigalgsigner, err)
					}
					pem.Encode(os.Stdout, &pem.Block{Type: "PKCS7", Bytes: signed})
					p7, err := Parse(signed)
					if err != nil {
						t.Fatalf("test %s/%s/%s: cannot parse signed data: %s", sigalgroot, sigalginter, sigalgsigner, err)
					}
					if testDetach {
						p7.Content = content
					}
					if !bytes.Equal(content, p7.Content) {
						t.Errorf("test %s/%s/%s: content was not found in the parsed data:\n\tExpected: %s\n\tActual: %s", sigalgroot, sigalginter, sigalgsigner, content, p7.Content)
					}
					if err := p7.VerifyWithChain(truststore); err != nil {
						t.Errorf("test %s/%s/%s: cannot verify signed data: %s", sigalgroot, sigalginter, sigalgsigner, err)
					}
					if !signerDigest.Equal(p7.Signers[0].DigestAlgorithm.Algorithm) {
						t.Errorf("test %s/%s/%s: expected digest algorithm %q but got %q",
							sigalgroot, sigalginter, sigalgsigner, signerDigest, p7.Signers[0].DigestAlgorithm.Algorithm)
					}
				}
			}
		}
	}
}

func TestDSASignAndVerifyWithOpenSSL(t *testing.T) {
	dsaPublicCert := []byte(`-----BEGIN CERTIFICATE-----
MIIDOjCCAvWgAwIBAgIEPCY/UDANBglghkgBZQMEAwIFADBsMRAwDgYDVQQGEwdV
bmtub3duMRAwDgYDVQQIEwdVbmtub3duMRAwDgYDVQQHEwdVbmtub3duMRAwDgYD
VQQKEwdVbmtub3duMRAwDgYDVQQLEwdVbmtub3duMRAwDgYDVQQDEwdVbmtub3du
MB4XDTE4MTAyMjEzNDMwN1oXDTQ2MDMwOTEzNDMwN1owbDEQMA4GA1UEBhMHVW5r
bm93bjEQMA4GA1UECBMHVW5rbm93bjEQMA4GA1UEBxMHVW5rbm93bjEQMA4GA1UE
ChMHVW5rbm93bjEQMA4GA1UECxMHVW5rbm93bjEQMA4GA1UEAxMHVW5rbm93bjCC
AbgwggEsBgcqhkjOOAQBMIIBHwKBgQD9f1OBHXUSKVLfSpwu7OTn9hG3UjzvRADD
Hj+AtlEmaUVdQCJR+1k9jVj6v8X1ujD2y5tVbNeBO4AdNG/yZmC3a5lQpaSfn+gE
exAiwk+7qdf+t8Yb+DtX58aophUPBPuD9tPFHsMCNVQTWhaRMvZ1864rYdcq7/Ii
Axmd0UgBxwIVAJdgUI8VIwvMspK5gqLrhAvwWBz1AoGBAPfhoIXWmz3ey7yrXDa4
V7l5lK+7+jrqgvlXTAs9B4JnUVlXjrrUWU/mcQcQgYC0SRZxI+hMKBYTt88JMozI
puE8FnqLVHyNKOCjrh4rs6Z1kW6jfwv6ITVi8ftiegEkO8yk8b6oUZCJqIPf4Vrl
nwaSi2ZegHtVJWQBTDv+z0kqA4GFAAKBgQDCriMPbEVBoRK4SOUeFwg7+VRf4TTp
rcOQC9IVVoCjXzuWEGrp3ZI7YWJSpFnSch4lk29RH8O0HpI/NOzKnOBtnKr782pt
1k/bJVMH9EaLd6MKnAVjrCDMYBB0MhebZ8QHY2elZZCWoqDYAcIDOsEx+m4NLErT
ypPnjS5M0jm1PKMhMB8wHQYDVR0OBBYEFC0Yt5XdM0Kc95IX8NQ8XRssGPx7MA0G
CWCGSAFlAwQDAgUAAzAAMC0CFQCIgQtrZZ9hdZG1ROhR5hc8nYEmbgIUAIlgC688
qzy/7yePTlhlpj+ahMM=
-----END CERTIFICATE-----`)

	content := []byte("Hello World")
	// write the content to a temp file
	tmpContentFile, err := ioutil.TempFile("", "TestDSASignAndVerifyWithOpenSSL_content")
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile(tmpContentFile.Name(), content, 0o755)

	block, _ := pem.Decode(dsaPublicCert)
	if block == nil {
		t.Fatal("failed to parse certificate PEM")
	}
	signerCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal("failed to parse certificate: " + err.Error())
	}

	// write the signer cert to a temp file
	tmpSignerCertFile, err := ioutil.TempFile("", "TestDSASignAndVerifyWithOpenSSL_signer")
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile(tmpSignerCertFile.Name(), dsaPublicCert, 0o755)

	priv := dsa.PrivateKey{
		PublicKey: dsa.PublicKey{
			Parameters: dsa.Parameters{
				P: fromHex("fd7f53811d75122952df4a9c2eece4e7f611b7523cef4400c31e3f80b6512669455d402251fb593d8d58fabfc5f5ba30f6cb9b556cd7813b801d346ff26660b76b9950a5a49f9fe8047b1022c24fbba9d7feb7c61bf83b57e7c6a8a6150f04fb83f6d3c51ec3023554135a169132f675f3ae2b61d72aeff22203199dd14801c7"),
				Q: fromHex("9760508F15230BCCB292B982A2EB840BF0581CF5"),
				G: fromHex("F7E1A085D69B3DDECBBCAB5C36B857B97994AFBBFA3AEA82F9574C0B3D0782675159578EBAD4594FE67107108180B449167123E84C281613B7CF09328CC8A6E13C167A8B547C8D28E0A3AE1E2BB3A675916EA37F0BFA213562F1FB627A01243BCCA4F1BEA8519089A883DFE15AE59F06928B665E807B552564014C3BFECF492A"),
			},
		},
		X: fromHex("7D6E1A3DD4019FD809669D8AB8DA73807CEF7EC1"),
	}
	toBeSigned, err := NewSignedData(content)
	if err != nil {
		t.Fatalf("test case: cannot initialize signed data: %s", err)
	}
	// openssl DSA only supports SHA1 for our 1024-bit DSA key, since that is all the standard officially supports
	toBeSigned.digestOid = OIDDigestAlgorithmSHA1
	if err := toBeSigned.SignWithoutAttr(signerCert, &priv, SignerInfoConfig{}); err != nil {
		t.Fatalf("Cannot add signer: %s", err)
	}
	toBeSigned.Detach()
	signed, err := toBeSigned.Finish()
	if err != nil {
		t.Fatalf("test case: cannot finish signing data: %s", err)
	}

	// write the signature to a temp file
	tmpSignatureFile, err := ioutil.TempFile("", "TestDSASignAndVerifyWithOpenSSL_signature")
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile(tmpSignatureFile.Name(), pem.EncodeToMemory(&pem.Block{Type: "PKCS7", Bytes: signed}), 0o755)

	// call openssl to verify the signature on the content using the root
	opensslCMD := exec.Command("openssl", "smime", "-verify", "-noverify",
		"-in", tmpSignatureFile.Name(), "-inform", "PEM",
		"-content", tmpContentFile.Name())
	out, err := opensslCMD.CombinedOutput()
	if err != nil {
		t.Errorf("Command: %s", opensslCMD.Args)
		t.Fatalf("test case: openssl command failed with %s: %s", err, out)
	}
	os.Remove(tmpSignatureFile.Name())  // clean up
	os.Remove(tmpContentFile.Name())    // clean up
	os.Remove(tmpSignerCertFile.Name()) // clean up
}

func ExampleSignedData() {
	// generate a signing cert or load a key pair
	cert, err := createTestCertificate(x509.SHA256WithRSA)
	if err != nil {
		fmt.Printf("Cannot create test certificates: %s", err)
	}

	// Initialize a SignedData struct with content to be signed
	signedData, err := NewSignedData([]byte("Example data to be signed"))
	if err != nil {
		fmt.Printf("Cannot initialize signed data: %s", err)
	}

	// Add the signing cert and private key
	if err := signedData.AddSigner(cert.Certificate, cert.PrivateKey, SignerInfoConfig{}); err != nil {
		fmt.Printf("Cannot add signer: %s", err)
	}

	// Call Detach() is you want to remove content from the signature
	// and generate an S/MIME detached signature
	signedData.Detach()

	// Finish() to obtain the signature bytes
	detachedSignature, err := signedData.Finish()
	if err != nil {
		fmt.Printf("Cannot finish signing data: %s", err)
	}
	pem.Encode(os.Stdout, &pem.Block{Type: "PKCS7", Bytes: detachedSignature})
}

func TestUnmarshalSignedAttribute(t *testing.T) {
	cert, err := createTestCertificate(x509.SHA512WithRSA)
	if err != nil {
		t.Fatal(err)
	}
	content := []byte("Hello World")
	toBeSigned, err := NewSignedData(content)
	if err != nil {
		t.Fatalf("Cannot initialize signed data: %s", err)
	}
	oidTest := asn1.ObjectIdentifier{2, 3, 4, 5, 6, 7}
	testValue := "TestValue"
	if err := toBeSigned.AddSigner(cert.Certificate, *cert.PrivateKey, SignerInfoConfig{
		ExtraSignedAttributes: []Attribute{{Type: oidTest, Value: testValue}},
	}); err != nil {
		t.Fatalf("Cannot add signer: %s", err)
	}
	signed, err := toBeSigned.Finish()
	if err != nil {
		t.Fatalf("Cannot finish signing data: %s", err)
	}
	p7, err := Parse(signed)
	if err != nil {
		t.Fatalf("Cannot parse signed data: %v", err)
	}
	var actual string
	err = p7.UnmarshalSignedAttribute(oidTest, &actual)
	if err != nil {
		t.Fatalf("Cannot unmarshal test value: %s", err)
	}
	if testValue != actual {
		t.Errorf("Attribute does not match test value\n\tExpected: %s\n\tActual: %s", testValue, actual)
	}
}

func TestDegenerateCertificate(t *testing.T) {
	cert, err := createTestCertificate(x509.SHA256WithRSA)
	if err != nil {
		t.Fatal(err)
	}
	deg, err := DegenerateCertificate(cert.Certificate.Raw)
	if err != nil {
		t.Fatal(err)
	}
	testOpenSSLParse(t, deg)
	pem.Encode(os.Stdout, &pem.Block{Type: "PKCS7", Bytes: deg})
}

// writes the cert to a temporary file and tests that openssl can read it.
func testOpenSSLParse(t *testing.T, certBytes []byte) {
	tmpCertFile, err := ioutil.TempFile("", "testCertificate")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpCertFile.Name()) // clean up

	if _, err := tmpCertFile.Write(certBytes); err != nil {
		t.Fatal(err)
	}

	opensslCMD := exec.Command("openssl", "pkcs7", "-inform", "der", "-in", tmpCertFile.Name())
	_, err = opensslCMD.Output()
	if err != nil {
		t.Fatal(err)
	}

	if err := tmpCertFile.Close(); err != nil {
		t.Fatal(err)
	}
}

func fromHex(s string) *big.Int {
	result, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(s)
	}
	return result
}
