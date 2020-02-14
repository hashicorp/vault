package reload

import (
	"crypto/x509"
	"io/ioutil"
	"testing"

	"github.com/hashicorp/errwrap"
)

func TestReload_KeyWithPassphrase(t *testing.T) {
	password := "password"
	cert := []byte(`-----BEGIN CERTIFICATE-----
MIICLzCCAZgCCQCq27CeP4WhlDANBgkqhkiG9w0BAQUFADBcMQswCQYDVQQGEwJV
UzELMAkGA1UECAwCQ0ExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xEjAQBgNVBAoM
CUhhc2hpQ29ycDEUMBIGA1UEAwwLbXl2YXVsdC5jb20wHhcNMTcxMjEzMjEzNTM3
WhcNMTgxMjEzMjEzNTM3WjBcMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExFjAU
BgNVBAcMDVNhbiBGcmFuY2lzY28xEjAQBgNVBAoMCUhhc2hpQ29ycDEUMBIGA1UE
AwwLbXl2YXVsdC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMvsz/9l
EJIlRG6DOw4fXdB/aJgJk2rR8cU0D8+vECIzb+MdDK0cBHtLiVpZC/RnZMdMzjGn
Z++Fp3dEnT6CD0IjKdJcD+qSyZSjHIuYpHjnjrVlM/Le0xST7egoG+fXkSt4myzG
ec2WK1jcZefRRGPycvMqx1yUWU76jDdFZSL5AgMBAAEwDQYJKoZIhvcNAQEFBQAD
gYEAQfYE26FLZ9SPPU8bHNDxoxDmGrn8yJ78C490Qpix/w6gdLaBtILenrZbhpnB
3L3okraM8mplaN2KdAcpnsr4wPv9hbYkam0coxCQEKs8ltHSBaXT6uKRWb00nkGu
yAXDRpuPdFRqbXW3ZFC5broUrz4ujxTDKfVeIn0zpPZkv24=
-----END CERTIFICATE-----`)
	key := []byte(`-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: DES-EDE3-CBC,64B032D83BD6A6DC

qVJ+mXEBKMkUPrQ8odHunMpPgChQUny4CX73/dAcm7O9iXIv9eXQSxj2qfgCOloj
vthg7jYNwtRb0ydzCEnEud35zWw38K/l19/pe4ULfNXlOddlsk4XIHarBiz+KUaX
WTbNk0H+DwdcEwhprPgpTk8gp88lZBiHCnTG/s8v/JNt+wkdqjfAp0Xbm9m+OZ7s
hlNxZin1OuBdprBqfKWBltUALZYiIBhspMTmh+jGQSyEKNTAIBejIiRH5+xYWuOy
xKencq8UpQMOMPR2ZiSw42dU9j8HHMgldI7KszU2FDIEFXG7aSjcxNyyybeBT+Uz
YPoxGxSdUYWqaz50UszvHg/QWR8NlPlQc3nFAUVpGKUF9MEQCIAK8HjcpMP+IAVO
ertp4cTa2Rpm9YeoFrY6tabvmXApXlQPw6rBn6o5KpceWG3ceOsDOsT+e3edHu9g
SGO4hjggbRpO+dBOuwfw4rMn9X1BbqXKJcREAmrgVVSf9/s942E4YOQ+IGJPdtmY
WHAFk8hiJepsVCA2NpwVlAD+QbPPaR2RtvYOtq3IKlWRuVQ+6dpxDsz5FlJhs2L+
HsX6XqtwuQM8kk1hO8Gm3VeV7+b64r9kfbO8jCM18GexCYiCtig51mJW6IO42d1K
bS1axMx/KeDc/sy7LKEbHnjnYanpGz2Wa2EWhnWAeNXD1nUfUNFPp2SsIGbCMnat
mC4O4cO7YRl3+iJg3kHtTPGtgtCjrZcjlyBtxT2VC7SsTcTXZBWovczMIstyr4Ka
opM24uvQT3Bc0UM0WNh3tdRFuboxDeBDh7PX/2RIoiaMuCCiRZ3O0A==
-----END RSA PRIVATE KEY-----`)
	tempDir, err := ioutil.TempDir("", "vault-test")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %s", err)
	}
	keyFile := tempDir + "/server.key"
	certFile := tempDir + "/server.crt"

	err = ioutil.WriteFile(certFile, cert, 0755)
	if err != nil {
		t.Fatalf("Error writing to temp file: %s", err)
	}
	err = ioutil.WriteFile(keyFile, key, 0755)
	if err != nil {
		t.Fatalf("Error writing to temp file: %s", err)
	}

	cg := NewCertificateGetter(certFile, keyFile, "")
	err = cg.Reload(nil)
	if err == nil {
		t.Fatal("error expected")
	}
	if !errwrap.Contains(err, x509.IncorrectPasswordError.Error()) {
		t.Fatalf("expected incorrect password error, got %v", err)
	}

	cg = NewCertificateGetter(certFile, keyFile, password)
	if err := cg.Reload(nil); err != nil {
		t.Fatalf("err: %v", err)
	}
}
