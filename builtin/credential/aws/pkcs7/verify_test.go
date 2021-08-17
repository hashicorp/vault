package pkcs7

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestVerify(t *testing.T) {
	fixture := UnmarshalTestFixture(SignedTestFixture)
	p7, err := Parse(fixture.Input)
	if err != nil {
		t.Errorf("Parse encountered unexpected error: %v", err)
	}

	if err := p7.Verify(); err != nil {
		t.Errorf("Verify failed with error: %v", err)
	}
	expected := []byte("We the People")
	if !bytes.Equal(p7.Content, expected) {
		t.Errorf("Signed content does not match.\n\tExpected:%s\n\tActual:%s", expected, p7.Content)

	}
}

var SignedTestFixture = `
-----BEGIN PKCS7-----
MIIDVgYJKoZIhvcNAQcCoIIDRzCCA0MCAQExCTAHBgUrDgMCGjAcBgkqhkiG9w0B
BwGgDwQNV2UgdGhlIFBlb3BsZaCCAdkwggHVMIIBQKADAgECAgRpuDctMAsGCSqG
SIb3DQEBCzApMRAwDgYDVQQKEwdBY21lIENvMRUwEwYDVQQDEwxFZGRhcmQgU3Rh
cmswHhcNMTUwNTA2MDQyNDQ4WhcNMTYwNTA2MDQyNDQ4WjAlMRAwDgYDVQQKEwdB
Y21lIENvMREwDwYDVQQDEwhKb24gU25vdzCBnzANBgkqhkiG9w0BAQEFAAOBjQAw
gYkCgYEAqr+tTF4mZP5rMwlXp1y+crRtFpuLXF1zvBZiYMfIvAHwo1ta8E1IcyEP
J1jIiKMcwbzeo6kAmZzIJRCTezq9jwXUsKbQTvcfOH9HmjUmXBRWFXZYoQs/OaaF
a45deHmwEeMQkuSWEtYiVKKZXtJOtflKIT3MryJEDiiItMkdybUCAwEAAaMSMBAw
DgYDVR0PAQH/BAQDAgCgMAsGCSqGSIb3DQEBCwOBgQDK1EweZWRL+f7Z+J0kVzY8
zXptcBaV4Lf5wGZJLJVUgp33bpLNpT3yadS++XQJ+cvtW3wADQzBSTMduyOF8Zf+
L7TjjrQ2+F2HbNbKUhBQKudxTfv9dJHdKbD+ngCCdQJYkIy2YexsoNG0C8nQkggy
axZd/J69xDVx6pui3Sj8sDGCATYwggEyAgEBMDEwKTEQMA4GA1UEChMHQWNtZSBD
bzEVMBMGA1UEAxMMRWRkYXJkIFN0YXJrAgRpuDctMAcGBSsOAwIaoGEwGAYJKoZI
hvcNAQkDMQsGCSqGSIb3DQEHATAgBgkqhkiG9w0BCQUxExcRMTUwNTA2MDAyNDQ4
LTA0MDAwIwYJKoZIhvcNAQkEMRYEFG9D7gcTh9zfKiYNJ1lgB0yTh4sZMAsGCSqG
SIb3DQEBAQSBgFF3sGDU9PtXty/QMtpcFa35vvIOqmWQAIZt93XAskQOnBq4OloX
iL9Ct7t1m4pzjRm0o9nDkbaSLZe7HKASHdCqijroScGlI8M+alJ8drHSFv6ZIjnM
FIwIf0B2Lko6nh9/6mUXq7tbbIHa3Gd1JUVire/QFFtmgRXMbXYk8SIS
-----END PKCS7-----
-----BEGIN CERTIFICATE-----
MIIB1TCCAUCgAwIBAgIEabg3LTALBgkqhkiG9w0BAQswKTEQMA4GA1UEChMHQWNt
ZSBDbzEVMBMGA1UEAxMMRWRkYXJkIFN0YXJrMB4XDTE1MDUwNjA0MjQ0OFoXDTE2
MDUwNjA0MjQ0OFowJTEQMA4GA1UEChMHQWNtZSBDbzERMA8GA1UEAxMISm9uIFNu
b3cwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAKq/rUxeJmT+azMJV6dcvnK0
bRabi1xdc7wWYmDHyLwB8KNbWvBNSHMhDydYyIijHMG83qOpAJmcyCUQk3s6vY8F
1LCm0E73Hzh/R5o1JlwUVhV2WKELPzmmhWuOXXh5sBHjEJLklhLWIlSimV7STrX5
SiE9zK8iRA4oiLTJHcm1AgMBAAGjEjAQMA4GA1UdDwEB/wQEAwIAoDALBgkqhkiG
9w0BAQsDgYEAytRMHmVkS/n+2fidJFc2PM16bXAWleC3+cBmSSyVVIKd926SzaU9
8mnUvvl0CfnL7Vt8AA0MwUkzHbsjhfGX/i+04460Nvhdh2zWylIQUCrncU37/XSR
3Smw/p4AgnUCWJCMtmHsbKDRtAvJ0JIIMmsWXfyevcQ1ceqbot0o/LA=
-----END CERTIFICATE-----
-----BEGIN PRIVATE KEY-----
MIICXgIBAAKBgQCqv61MXiZk/mszCVenXL5ytG0Wm4tcXXO8FmJgx8i8AfCjW1rw
TUhzIQ8nWMiIoxzBvN6jqQCZnMglEJN7Or2PBdSwptBO9x84f0eaNSZcFFYVdlih
Cz85poVrjl14ebAR4xCS5JYS1iJUople0k61+UohPcyvIkQOKIi0yR3JtQIDAQAB
AoGBAIPLCR9N+IKxodq11lNXEaUFwMHXc1zqwP8no+2hpz3+nVfplqqubEJ4/PJY
5AgbJoIfnxVhyBXJXu7E+aD/OPneKZrgp58YvHKgGvvPyJg2gpC/1Fh0vQB0HNpI
1ZzIZUl8ZTUtVgtnCBUOh5JGI4bFokAqrT//Uvcfd+idgxqBAkEA1ZbP/Kseld14
qbWmgmU5GCVxsZRxgR1j4lG3UVjH36KXMtRTm1atAam1uw3OEGa6Y3ANjpU52FaB
Hep5rkk4FQJBAMynMo1L1uiN5GP+KYLEF5kKRxK+FLjXR0ywnMh+gpGcZDcOae+J
+t1gLoWBIESH/Xt639T7smuSfrZSA9V0EyECQA8cvZiWDvLxmaEAXkipmtGPjKzQ
4PsOtkuEFqFl07aKDYKmLUg3aMROWrJidqsIabWxbvQgsNgSvs38EiH3wkUCQQCg
ndxb7piVXb9RBwm3OoU2tE1BlXMX+sVXmAkEhd2dwDsaxrI3sHf1xGXem5AimQRF
JBOFyaCnMotGNioSHY5hAkEAxyXcNixQ2RpLXJTQZtwnbk0XDcbgB+fBgXnv/4f3
BCvcu85DqJeJyQv44Oe1qsXEX9BfcQIOVaoep35RPlKi9g==
-----END PRIVATE KEY-----`

func TestVerifyAppStore(t *testing.T) {
	fixture := UnmarshalTestFixture(AppStoreReceiptFixture)
	p7, err := Parse(fixture.Input)
	if err != nil {
		t.Errorf("Parse encountered unexpected error: %v", err)
	}
	if err := p7.Verify(); err != nil {
		t.Errorf("Verify failed with error: %v", err)
	}
}

var AppStoreReceiptFixture = `
-----BEGIN PKCS7-----
MIITtgYJKoZIhvcNAQcCoIITpzCCE6MCAQExCzAJBgUrDgMCGgUAMIIDVwYJKoZI
hvcNAQcBoIIDSASCA0QxggNAMAoCAQgCAQEEAhYAMAoCARQCAQEEAgwAMAsCAQEC
AQEEAwIBADALAgEDAgEBBAMMATEwCwIBCwIBAQQDAgEAMAsCAQ8CAQEEAwIBADAL
AgEQAgEBBAMCAQAwCwIBGQIBAQQDAgEDMAwCAQoCAQEEBBYCNCswDAIBDgIBAQQE
AgIAjTANAgENAgEBBAUCAwFgvTANAgETAgEBBAUMAzEuMDAOAgEJAgEBBAYCBFAy
NDcwGAIBAgIBAQQQDA5jb20uemhpaHUudGVzdDAYAgEEAgECBBCS+ZODNMHwT1Nz
gWYDXyWZMBsCAQACAQEEEwwRUHJvZHVjdGlvblNhbmRib3gwHAIBBQIBAQQU4nRh
YCEZx70Flzv7hvJRjJZckYIwHgIBDAIBAQQWFhQyMDE2LTA3LTIzVDA2OjIxOjEx
WjAeAgESAgEBBBYWFDIwMTMtMDgtMDFUMDc6MDA6MDBaMD0CAQYCAQEENbR21I+a
8+byMXo3NPRoDWQmSXQF2EcCeBoD4GaL//ZCRETp9rGFPSg1KekCP7Kr9HAqw09m
MEICAQcCAQEEOlVJozYYBdugybShbiiMsejDMNeCbZq6CrzGBwW6GBy+DGWxJI91
Y3ouXN4TZUhuVvLvN1b0m5T3ggQwggFaAgERAgEBBIIBUDGCAUwwCwICBqwCAQEE
AhYAMAsCAgatAgEBBAIMADALAgIGsAIBAQQCFgAwCwICBrICAQEEAgwAMAsCAgaz
AgEBBAIMADALAgIGtAIBAQQCDAAwCwICBrUCAQEEAgwAMAsCAga2AgEBBAIMADAM
AgIGpQIBAQQDAgEBMAwCAgarAgEBBAMCAQEwDAICBq4CAQEEAwIBADAMAgIGrwIB
AQQDAgEAMAwCAgaxAgEBBAMCAQAwGwICBqcCAQEEEgwQMTAwMDAwMDIyNTMyNTkw
MTAbAgIGqQIBAQQSDBAxMDAwMDAwMjI1MzI1OTAxMB8CAgaoAgEBBBYWFDIwMTYt
MDctMjNUMDY6MjE6MTFaMB8CAgaqAgEBBBYWFDIwMTYtMDctMjNUMDY6MjE6MTFa
MCACAgamAgEBBBcMFWNvbS56aGlodS50ZXN0LnRlc3RfMaCCDmUwggV8MIIEZKAD
AgECAggO61eH554JjTANBgkqhkiG9w0BAQUFADCBljELMAkGA1UEBhMCVVMxEzAR
BgNVBAoMCkFwcGxlIEluYy4xLDAqBgNVBAsMI0FwcGxlIFdvcmxkd2lkZSBEZXZl
bG9wZXIgUmVsYXRpb25zMUQwQgYDVQQDDDtBcHBsZSBXb3JsZHdpZGUgRGV2ZWxv
cGVyIFJlbGF0aW9ucyBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTAeFw0xNTExMTMw
MjE1MDlaFw0yMzAyMDcyMTQ4NDdaMIGJMTcwNQYDVQQDDC5NYWMgQXBwIFN0b3Jl
IGFuZCBpVHVuZXMgU3RvcmUgUmVjZWlwdCBTaWduaW5nMSwwKgYDVQQLDCNBcHBs
ZSBXb3JsZHdpZGUgRGV2ZWxvcGVyIFJlbGF0aW9uczETMBEGA1UECgwKQXBwbGUg
SW5jLjELMAkGA1UEBhMCVVMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQClz4H9JaKBW9aH7SPaMxyO4iPApcQmyz3Gn+xKDVWG/6QC15fKOVRtfX+yVBid
xCxScY5ke4LOibpJ1gjltIhxzz9bRi7GxB24A6lYogQ+IXjV27fQjhKNg0xbKmg3
k8LyvR7E0qEMSlhSqxLj7d0fmBWQNS3CzBLKjUiB91h4VGvojDE2H0oGDEdU8zeQ
uLKSiX1fpIVK4cCc4Lqku4KXY/Qrk8H9Pm/KwfU8qY9SGsAlCnYO3v6Z/v/Ca/Vb
XqxzUUkIVonMQ5DMjoEC0KCXtlyxoWlph5AQaCYmObgdEHOwCl3Fc9DfdjvYLdmI
HuPsB8/ijtDT+iZVge/iA0kjAgMBAAGjggHXMIIB0zA/BggrBgEFBQcBAQQzMDEw
LwYIKwYBBQUHMAGGI2h0dHA6Ly9vY3NwLmFwcGxlLmNvbS9vY3NwMDMtd3dkcjA0
MB0GA1UdDgQWBBSRpJz8xHa3n6CK9E31jzZd7SsEhTAMBgNVHRMBAf8EAjAAMB8G
A1UdIwQYMBaAFIgnFwmpthhgi+zruvZHWcVSVKO3MIIBHgYDVR0gBIIBFTCCAREw
ggENBgoqhkiG92NkBQYBMIH+MIHDBggrBgEFBQcCAjCBtgyBs1JlbGlhbmNlIG9u
IHRoaXMgY2VydGlmaWNhdGUgYnkgYW55IHBhcnR5IGFzc3VtZXMgYWNjZXB0YW5j
ZSBvZiB0aGUgdGhlbiBhcHBsaWNhYmxlIHN0YW5kYXJkIHRlcm1zIGFuZCBjb25k
aXRpb25zIG9mIHVzZSwgY2VydGlmaWNhdGUgcG9saWN5IGFuZCBjZXJ0aWZpY2F0
aW9uIHByYWN0aWNlIHN0YXRlbWVudHMuMDYGCCsGAQUFBwIBFipodHRwOi8vd3d3
LmFwcGxlLmNvbS9jZXJ0aWZpY2F0ZWF1dGhvcml0eS8wDgYDVR0PAQH/BAQDAgeA
MBAGCiqGSIb3Y2QGCwEEAgUAMA0GCSqGSIb3DQEBBQUAA4IBAQANphvTLj3jWysH
bkKWbNPojEMwgl/gXNGNvr0PvRr8JZLbjIXDgFnf4+LXLgUUrA3btrj+/DUufMut
F2uOfx/kd7mxZ5W0E16mGYZ2+FogledjjA9z/Ojtxh+umfhlSFyg4Cg6wBA3Lbmg
BDkfc7nIBf3y3n8aKipuKwH8oCBc2et9J6Yz+PWY4L5E27FMZ/xuCk/J4gao0pfz
p45rUaJahHVl0RYEYuPBX/UIqc9o2ZIAycGMs/iNAGS6WGDAfK+PdcppuVsq1h1o
bphC9UynNxmbzDscehlD86Ntv0hgBgw2kivs3hi1EdotI9CO/KBpnBcbnoB7OUdF
MGEvxxOoMIIEIjCCAwqgAwIBAgIIAd68xDltoBAwDQYJKoZIhvcNAQEFBQAwYjEL
MAkGA1UEBhMCVVMxEzARBgNVBAoTCkFwcGxlIEluYy4xJjAkBgNVBAsTHUFwcGxl
IENlcnRpZmljYXRpb24gQXV0aG9yaXR5MRYwFAYDVQQDEw1BcHBsZSBSb290IENB
MB4XDTEzMDIwNzIxNDg0N1oXDTIzMDIwNzIxNDg0N1owgZYxCzAJBgNVBAYTAlVT
MRMwEQYDVQQKDApBcHBsZSBJbmMuMSwwKgYDVQQLDCNBcHBsZSBXb3JsZHdpZGUg
RGV2ZWxvcGVyIFJlbGF0aW9uczFEMEIGA1UEAww7QXBwbGUgV29ybGR3aWRlIERl
dmVsb3BlciBSZWxhdGlvbnMgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwggEiMA0G
CSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDKOFSmy1aqyCQ5SOmM7uxfuH8mkbw0
U3rOfGOAYXdkXqUHI7Y5/lAtFVZYcC1+xG7BSoU+L/DehBqhV8mvexj/avoVEkkV
CBmsqtsqMu2WY2hSFT2Miuy/axiV4AOsAX2XBWfODoWVN2rtCbauZ81RZJ/GXNG8
V25nNYB2NqSHgW44j9grFU57Jdhav06DwY3Sk9UacbVgnJ0zTlX5ElgMhrgWDcHl
d0WNUEi6Ky3klIXh6MSdxmilsKP8Z35wugJZS3dCkTm59c3hTO/AO0iMpuUhXf1q
arunFjVg0uat80YpyejDi+l5wGphZxWy8P3laLxiX27Pmd3vG2P+kmWrAgMBAAGj
gaYwgaMwHQYDVR0OBBYEFIgnFwmpthhgi+zruvZHWcVSVKO3MA8GA1UdEwEB/wQF
MAMBAf8wHwYDVR0jBBgwFoAUK9BpR5R2Cf70a40uQKb3R01/CF4wLgYDVR0fBCcw
JTAjoCGgH4YdaHR0cDovL2NybC5hcHBsZS5jb20vcm9vdC5jcmwwDgYDVR0PAQH/
BAQDAgGGMBAGCiqGSIb3Y2QGAgEEAgUAMA0GCSqGSIb3DQEBBQUAA4IBAQBPz+9Z
viz1smwvj+4ThzLoBTWobot9yWkMudkXvHcs1Gfi/ZptOllc34MBvbKuKmFysa/N
w0Uwj6ODDc4dR7Txk4qjdJukw5hyhzs+r0ULklS5MruQGFNrCk4QttkdUGwhgAqJ
TleMa1s8Pab93vcNIx0LSiaHP7qRkkykGRIZbVf1eliHe2iK5IaMSuviSRSqpd1V
AKmuu0swruGgsbwpgOYJd+W+NKIByn/c4grmO7i77LpilfMFY0GCzQ87HUyVpNur
+cmV6U/kTecmmYHpvPm0KdIBembhLoz2IYrF+Hjhga6/05Cdqa3zr/04GpZnMBxR
pVzscYqCtGwPDBUfMIIEuzCCA6OgAwIBAgIBAjANBgkqhkiG9w0BAQUFADBiMQsw
CQYDVQQGEwJVUzETMBEGA1UEChMKQXBwbGUgSW5jLjEmMCQGA1UECxMdQXBwbGUg
Q2VydGlmaWNhdGlvbiBBdXRob3JpdHkxFjAUBgNVBAMTDUFwcGxlIFJvb3QgQ0Ew
HhcNMDYwNDI1MjE0MDM2WhcNMzUwMjA5MjE0MDM2WjBiMQswCQYDVQQGEwJVUzET
MBEGA1UEChMKQXBwbGUgSW5jLjEmMCQGA1UECxMdQXBwbGUgQ2VydGlmaWNhdGlv
biBBdXRob3JpdHkxFjAUBgNVBAMTDUFwcGxlIFJvb3QgQ0EwggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQDkkakJH5HbHkdQ6wXtXnmELes2oldMVeyLGYne
+Uts9QerIjAC6Bg++FAJ039BqJj50cpmnCRrEdCju+QbKsMflZ56DKRHi1vUFjcz
y8QPTc4UadHJGXL1XQ7Vf1+b8iUDulWPTV0N8WQ1IxVLFVkds5T39pyez1C6wVhQ
Z48ItCD3y6wsIG9wtj8BMIy3Q88PnT3zK0koGsj+zrW5DtleHNbLPbU6rfQPDgCS
C7EhFi501TwN22IWq6NxkkdTVcGvL0Gz+PvjcM3mo0xFfh9Ma1CWQYnEdGILEINB
hzOKgbEwWOxaBDKMaLOPHd5lc/9nXmW8Sdh2nzMUZaF3lMktAgMBAAGjggF6MIIB
djAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUK9Bp
R5R2Cf70a40uQKb3R01/CF4wHwYDVR0jBBgwFoAUK9BpR5R2Cf70a40uQKb3R01/
CF4wggERBgNVHSAEggEIMIIBBDCCAQAGCSqGSIb3Y2QFATCB8jAqBggrBgEFBQcC
ARYeaHR0cHM6Ly93d3cuYXBwbGUuY29tL2FwcGxlY2EvMIHDBggrBgEFBQcCAjCB
thqBs1JlbGlhbmNlIG9uIHRoaXMgY2VydGlmaWNhdGUgYnkgYW55IHBhcnR5IGFz
c3VtZXMgYWNjZXB0YW5jZSBvZiB0aGUgdGhlbiBhcHBsaWNhYmxlIHN0YW5kYXJk
IHRlcm1zIGFuZCBjb25kaXRpb25zIG9mIHVzZSwgY2VydGlmaWNhdGUgcG9saWN5
IGFuZCBjZXJ0aWZpY2F0aW9uIHByYWN0aWNlIHN0YXRlbWVudHMuMA0GCSqGSIb3
DQEBBQUAA4IBAQBcNplMLXi37Yyb3PN3m/J20ncwT8EfhYOFG5k9RzfyqZtAjizU
sZAS2L70c5vu0mQPy3lPNNiiPvl4/2vIB+x9OYOLUyDTOMSxv5pPCmv/K/xZpwUJ
fBdAVhEedNO3iyM7R6PVbyTi69G3cN8PReEnyvFteO3ntRcXqNx+IjXKJdXZD9Zr
1KIkIxH3oayPc4FgxhtbCS+SsvhESPBgOJ4V9T0mZyCKM2r3DYLP3uujL/lTaltk
wGMzd/c6ByxW69oPIQ7aunMZT7XZNn/Bh1XZp5m5MkL72NVxnn6hUrcbvZNCJBIq
xw8dtk2cXmPIS4AXUKqK1drk/NAJBzewdXUhMYIByzCCAccCAQEwgaMwgZYxCzAJ
BgNVBAYTAlVTMRMwEQYDVQQKDApBcHBsZSBJbmMuMSwwKgYDVQQLDCNBcHBsZSBX
b3JsZHdpZGUgRGV2ZWxvcGVyIFJlbGF0aW9uczFEMEIGA1UEAww7QXBwbGUgV29y
bGR3aWRlIERldmVsb3BlciBSZWxhdGlvbnMgQ2VydGlmaWNhdGlvbiBBdXRob3Jp
dHkCCA7rV4fnngmNMAkGBSsOAwIaBQAwDQYJKoZIhvcNAQEBBQAEggEAasPtnide
NWyfUtewW9OSgcQA8pW+5tWMR0469cBPZR84uJa0gyfmPspySvbNOAwnrwzZHYLa
ujOxZLip4DUw4F5s3QwUa3y4BXpF4J+NSn9XNvxNtnT/GcEQtCuFwgJ0o3F0ilhv
MTHrwiwyx/vr+uNDqlORK8lfK+1qNp+A/kzh8eszMrn4JSeTh9ZYxLHE56WkTQGD
VZXl0gKgxSOmDrcp1eQxdlymzrPv9U60wUJ0bkPfrU9qZj3mJrmrkQk61JTe3j6/
QfjfFBG9JG2mUmYQP1KQ3SypGHzDW8vngvsGu//tNU0NFfOqQu4bYU4VpQl0nPtD
4B85NkrgvQsWAQ==
-----END PKCS7-----`

func TestVerifyApkEcdsa(t *testing.T) {
	fixture := UnmarshalTestFixture(ApkEcdsaFixture)
	p7, err := Parse(fixture.Input)
	if err != nil {
		t.Errorf("Parse encountered unexpected error: %v", err)
	}
	p7.Content, err = base64.StdEncoding.DecodeString(ApkEcdsaContent)
	if err != nil {
		t.Errorf("Failed to decode base64 signature file: %v", err)
	}
	if err := p7.Verify(); err != nil {
		t.Errorf("Verify failed with error: %v", err)
	}
}

var ApkEcdsaFixture = `-----BEGIN PKCS7-----
MIIDAgYJKoZIhvcNAQcCoIIC8zCCAu8CAQExDzANBglghkgBZQMEAgMFADALBgkq
hkiG9w0BBwGgggH3MIIB8zCCAVSgAwIBAgIJAOxXdFsvm3YiMAoGCCqGSM49BAME
MBIxEDAOBgNVBAMMB2VjLXA1MjEwHhcNMTYwMzMxMTUzMTIyWhcNNDMwODE3MTUz
MTIyWjASMRAwDgYDVQQDDAdlYy1wNTIxMIGbMBAGByqGSM49AgEGBSuBBAAjA4GG
AAQAYX95sSjPEQqgyLD04tNUyq9y/w8seblOpfqa/Amx6H4GFdrjGXX0YTfXKr9G
hAyIyQSnNrIg0zDlWQUbBPRW4CYBLFOg1pUn1NBhKFD4NtO1KWvYtNOYDegFjRCP
B0p+fEXDbq8QFDYvlh+NZUJ16+ih8XNIf1C29xuLEqN6oKOnAvajUDBOMB0GA1Ud
DgQWBBT/Ra3kz60gQ7tYk3byZckcLabt8TAfBgNVHSMEGDAWgBT/Ra3kz60gQ7tY
k3byZckcLabt8TAMBgNVHRMEBTADAQH/MAoGCCqGSM49BAMEA4GMADCBiAJCAP39
hYLsWk2H84oEw+HJqGGjexhqeD3vSO1mWhopripE/81oy3yV10puYoJe11xDSfcD
j2VfNCHazuXO3kSxGA/1AkIBLUJxp/WYbYzhBGKr6lcxczKI/wuMfkZ6vL+0EMJV
A/2uEoeqvnl7BsdkicyaOBNEADijuVdaPPIWzKClt9OaVxExgdAwgc0CAQEwHzAS
MRAwDgYDVQQDDAdlYy1wNTIxAgkA7Fd0Wy+bdiIwDQYJYIZIAWUDBAIDBQAwCgYI
KoZIzj0EAwQEgYswgYgCQgD1pVSNo7qTm9A6tpt3SU2yRa+xpJAnUbpZ+Gu36B71
JnQBUzRgTGevniqHpyagi7b2zjWh1uvfz9FfrITUwGMddgJCAPjiBRcl7rKpxmZn
V1MvcJOX41xRSJu1wmBiYXqaJarL+gQ/Wl7RYsMtqLjmNColvLaHNxCaWOO/8nAE
Hg0OMA60
-----END PKCS7-----`

var ApkEcdsaContent = `U2lnbmF0dXJlLVZlcnNpb246IDEuMA0KU0hBLTUxMi1EaWdlc3QtTWFuaWZlc3Q6IFAvVDRqSWtTMjQvNzFxeFE2WW1MeEtNdkRPUUF0WjUxR090dFRzUU9yemhHRQ0KIEpaUGVpWUtyUzZYY090bStYaWlFVC9uS2tYdWVtUVBwZ2RBRzFKUzFnPT0NCkNyZWF0ZWQtQnk6IDEuMCAoQW5kcm9pZCBTaWduQXBrKQ0KDQpOYW1lOiBBbmRyb2lkTWFuaWZlc3QueG1sDQpTSEEtNTEyLURpZ2VzdDogcm9NbWVQZmllYUNQSjFJK2VzMVpsYis0anB2UXowNHZqRWVpL2U0dkN1ald0VVVWSHEzMkNXDQogMUxsOHZiZGMzMCtRc1FlN29ibld4dmhLdXN2K3c1a2c9PQ0KDQpOYW1lOiByZXNvdXJjZXMuYXJzYw0KU0hBLTUxMi1EaWdlc3Q6IG5aYW1aUzlPZTRBRW41cEZaaCtoQ1JFT3krb1N6a3hHdU5YZU0wUFF6WGVBVlVQV3hSVzFPYQ0KIGVLbThRbXdmTmhhaS9HOEcwRUhIbHZEQWdlcy9HUGtBPT0NCg0KTmFtZTogY2xhc3Nlcy5kZXgNClNIQS01MTItRGlnZXN0OiBlbWlDQld2bkVSb0g2N2lCa3EwcUgrdm5tMkpaZDlMWUNEV051N3RNYzJ3bTRtV0dYSUVpWmcNCiBWZkVPV083MFRlZnFjUVhldkNtN2hQMnRpT0U3Y0w5UT09DQoNCg==`

func TestVerifyFirefoxAddon(t *testing.T) {
	fixture := UnmarshalTestFixture(FirefoxAddonFixture)
	p7, err := Parse(fixture.Input)
	if err != nil {
		t.Errorf("Parse encountered unexpected error: %v", err)
	}
	p7.Content = FirefoxAddonContent
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(FirefoxAddonRootCert)
	// verifies at the signingTime authenticated attr
	if err := p7.VerifyWithChain(certPool); err != nil {
		t.Errorf("Verify failed with error: %v", err)
	}

	// TODO: update to check for an expiration error when the EE
	// expires on 2021-08-16 20:04:58 +0000 UTC
	//
	// The chain has validity:
	//
	// EE:           2016-08-17 20:04:58 +0000 UTC 2021-08-16 20:04:58 +0000 UTC
	// Intermediate: 2015-03-17 23:52:42 +0000 UTC 2025-03-14 23:52:42 +0000 UTC
	// Root:         2015-03-17 22:53:57 +0000 UTC 2025-03-14 22:53:57 +0000 UTC
	if err = p7.VerifyWithChainAtTime(certPool, time.Now().UTC()); err != nil {
		t.Errorf("Verify at UTC now failed with error: %v", err)
	}

	expiredTime := time.Date(2030, time.January, 1, 0, 0, 0, 0, time.UTC)
	if err = p7.VerifyWithChainAtTime(certPool, expiredTime); err == nil {
		t.Errorf("Verify at expired time %s did not error", expiredTime)
	}
	notYetValidTime := time.Date(1999, time.July, 5, 0, 13, 0, 0, time.UTC)
	if err = p7.VerifyWithChainAtTime(certPool, notYetValidTime); err == nil {
		t.Errorf("Verify at not yet valid time %s did not error", notYetValidTime)
	}

	// Verify the certificate chain to make sure the identified root
	// is the one we expect
	ee := getCertFromCertsByIssuerAndSerial(p7.Certificates, p7.Signers[0].IssuerAndSerialNumber)
	if ee == nil {
		t.Errorf("No end-entity certificate found for signer")
	}
	signingTime, _ := time.Parse(time.RFC3339, "2017-02-23 09:06:16-05:00")
	chains, err := verifyCertChain(ee, p7.Certificates, certPool, signingTime)
	if err != nil {
		t.Error(err)
	}
	if len(chains) != 1 {
		t.Errorf("Expected to find one chain, but found %d", len(chains))
	}
	if len(chains[0]) != 3 {
		t.Errorf("Expected to find three certificates in chain, but found %d", len(chains[0]))
	}
	if chains[0][0].Subject.CommonName != "tabscope@xuldev.org" {
		t.Errorf("Expected to find EE certificate with subject 'tabscope@xuldev.org', but found '%s'", chains[0][0].Subject.CommonName)
	}
	if chains[0][1].Subject.CommonName != "production-signing-ca.addons.mozilla.org" {
		t.Errorf("Expected to find intermediate certificate with subject 'production-signing-ca.addons.mozilla.org', but found '%s'", chains[0][1].Subject.CommonName)
	}
	if chains[0][2].Subject.CommonName != "root-ca-production-amo" {
		t.Errorf("Expected to find root certificate with subject 'root-ca-production-amo', but found '%s'", chains[0][2].Subject.CommonName)
	}
}

var FirefoxAddonContent = []byte(`Signature-Version: 1.0
MD5-Digest-Manifest: KjRavc6/KNpuT1QLcB/Gsg==
SHA1-Digest-Manifest: 5Md5nUg+U7hQ/UfzV+xGKWOruVI=

`)

var FirefoxAddonFixture = `
-----BEGIN PKCS7-----
MIIQTAYJKoZIhvcNAQcCoIIQPTCCEDkCAQExCzAJBgUrDgMCGgUAMAsGCSqGSIb3
DQEHAaCCDL0wggW6MIIDoqADAgECAgYBVpobWVwwDQYJKoZIhvcNAQELBQAwgcUx
CzAJBgNVBAYTAlVTMRwwGgYDVQQKExNNb3ppbGxhIENvcnBvcmF0aW9uMS8wLQYD
VQQLEyZNb3ppbGxhIEFNTyBQcm9kdWN0aW9uIFNpZ25pbmcgU2VydmljZTExMC8G
A1UEAxMocHJvZHVjdGlvbi1zaWduaW5nLWNhLmFkZG9ucy5tb3ppbGxhLm9yZzE0
MDIGCSqGSIb3DQEJARYlc2VydmljZXMtb3BzK2FkZG9uc2lnbmluZ0Btb3ppbGxh
LmNvbTAeFw0xNjA4MTcyMDA0NThaFw0yMTA4MTYyMDA0NThaMHYxEzARBgNVBAsT
ClByb2R1Y3Rpb24xCzAJBgNVBAYTAlVTMRYwFAYDVQQHEw1Nb3VudGFpbiBWaWV3
MQ8wDQYDVQQKEwZBZGRvbnMxCzAJBgNVBAgTAkNBMRwwGgYDVQQDFBN0YWJzY29w
ZUB4dWxkZXYub3JnMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAv6e0
mPD8dt4J8HTNNq4ODns2DV6Weh1hllCIFvOeu1u3UrR03st0BMY8OXYwr/NvRVjg
bA8gRySWAL+XqLzbhtXNeNegAoxrF+3mYY5rJjsLj/FGI6P6OXjngqwgm9VTBl7m
jh/KXBSwYoUcavJo6cmk8sCFwoblyQiv+tsWaUCOI6zMzubNtIS+GFvET9y/VZMP
j6mk8O10wBgJF5MMtA19va3qXy7aCZ7DnZp1l3equd/L6t324TtXoqx6xWQKo6TM
I0mcTlKvm6TKegTGBCyGn3JRARoIJv4AW1qqgyaHXf9EoY2pKT8Avkri5++NuSJ6
jtO4k/diBA2MZU20U0KGffYZNTxKDqd6XtI6y1tJPd/OWRFyU+mHntkcm9sar7L3
nPKujHRox2re10ec1WBnJE3PjlAoesNjxzp+xs2mGGc8DX9NuWn+1uK9xmgGIIMl
OFfyQ4s0G6hKp5goFcrFZxmexu0ZahOs8vZf8xDBW7yR1zToQElOXHvrscM386os
kOF9IxQZfcCoPuNQVg1haCONNkx0oau3RQQlOSAZtC79b+rBjQ5JYfjRLYAworf2
xQaprCh33TD1dTBrvzEbCGszgkN53Vqh5TFBjbU/NyldOkGvK8Xf6WhT5u+aftnV
lbuE2McAg6x1AlloUZq6PNTBpz7zypcIISnQ+y8CAwEAATANBgkqhkiG9w0BAQsF
AAOCAgEAIBoo2+OEYNCgP/IbUj9azaf/lde1q4AK/uTMoUeS5WcrXd8aqA0Y1qV7
xUALgDQAExXgqcOMGu4mPMaoZDgwGI4Tj7XPJQq5Z5zYxpRf/Wtzae33T9BF6QPW
v5xiRYuol+FbEtqRHZqxDWtIrd1MWBy3wjO3pLPdzDM9jWh+HLxdGWThJszaZp3T
CqsOx+l9W0Q7qM5ioZpHStgXDfhw38Lg++kLnzcX9MqsjYyezdwE4krqW6hK3+4S
0LZE4dTgsy8JULkyAF3HrPWEXESnD7c4mx6owZe+BNDK5hsVM/obAqH7sJq/igbM
5N1l832p/ws8l5xKOr3qBWSzWn6u7ExvqG6Ckh0foJOVXvzGqvrXcoiBGV8S9Z7c
DghUvMt6b0pZ0ildRCHfTUz7eG3g4MhfbjupR7b+L9FWEJhcd/H0dxpw7SKYha/n
ePuRL7MXmbW8WLMqO/ImxzL8TPOB3pUg3nITfubV6gpPBmn+0nwbqYUmggJuwgvK
I2GpN2Ny6EErZy17EEgyhJygJZMj+UzQjC781xxsl3ljpYEqqwgRLIZBSBUD5dXj
XBuU24w162SeSyHZzkBbuv6lr52pqoZyFrG29DCHECgO9ZmNWgSpiWSkh+vExAG7
wNs0y61t2HUG+BCMGPQ9sOzouyTfrnLVAWwzswGftFYQfoIBeJIwggb7MIIE46AD
AgECAgMQAAIwDQYJKoZIhvcNAQEMBQAwfTELMAkGA1UEBhMCVVMxHDAaBgNVBAoT
E01vemlsbGEgQ29ycG9yYXRpb24xLzAtBgNVBAsTJk1vemlsbGEgQU1PIFByb2R1
Y3Rpb24gU2lnbmluZyBTZXJ2aWNlMR8wHQYDVQQDExZyb290LWNhLXByb2R1Y3Rp
b24tYW1vMB4XDTE1MDMxNzIzNTI0MloXDTI1MDMxNDIzNTI0MlowgcUxCzAJBgNV
BAYTAlVTMRwwGgYDVQQKExNNb3ppbGxhIENvcnBvcmF0aW9uMS8wLQYDVQQLEyZN
b3ppbGxhIEFNTyBQcm9kdWN0aW9uIFNpZ25pbmcgU2VydmljZTExMC8GA1UEAxMo
cHJvZHVjdGlvbi1zaWduaW5nLWNhLmFkZG9ucy5tb3ppbGxhLm9yZzE0MDIGCSqG
SIb3DQEJARYlc2VydmljZXMtb3BzK2FkZG9uc2lnbmluZ0Btb3ppbGxhLmNvbTCC
AiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAMLMM9m2HBLhCiO9mhljpehT
hxpzlCnxluzDZ51I/H7MvBbIvZBm9zSpHdffubSsak2qYE69d+ebTa/CK83WIosM
24/2Qp7n/GGaPJcCC4Y3JkrCsgA8+wV2MbFlKSv+qMdvI/sE3BPYDMCjVPMhHmIP
XaPWd42OoHpI8R3GGUtVnR3Hm76pa2+v6TwgeMiO8om+ogGufiyv6FNMZ5NuY1Z9
aLNEvehnAzSfddQyki+6FJd7XkgZbP7pb1Kl8yYgiy4piBerJ9H09uPehffE3Ell
3cApQL3+0kjaUX4scMjuNQDMKziRZkYgJAM+qA9WA5Jn77AjerQBWQeEev1PWHYh
0IDlgS/a0bjKmVjNZYG6adrY/R5/whzWGFCIE1UfhPm6PdN0557qvF838C2RFHsI
KzV6KQf0chMjpa02tPaIctjVhnDQZZNKm2ZfLOt9kQ57Is/e6KxH7pYMit46+s99
lYM7ZquvWbK19b1Ili/6S1BxSzd3wztgfN5jGsc+jCCYLm+AcVtfNKc8cFZHXKrB
CwhGmdbWDSBCicZNA7FKJpO3oIx26VPF2XUldA/T5Mh/POGLilK3t9m9qbjEyDp1
EwoBToOR/aMrdnNYvSWp0g/GHMzSfJjjXyAqrZY2itam/IJd8r9FoRAzevPt/zTX
BET3INoiCDGRH0XrxUYtAgMGVTejggE5MIIBNTAMBgNVHRMEBTADAQH/MA4GA1Ud
DwEB/wQEAwIBBjAWBgNVHSUBAf8EDDAKBggrBgEFBQcDAzAdBgNVHQ4EFgQUdHxf
FKXipZLjs20GqIUdNkQXH4gwgagGA1UdIwSBoDCBnYAUs7zqWHSr4W54KrKrnCMe
qGMsl7ehgYGkfzB9MQswCQYDVQQGEwJVUzEcMBoGA1UEChMTTW96aWxsYSBDb3Jw
b3JhdGlvbjEvMC0GA1UECxMmTW96aWxsYSBBTU8gUHJvZHVjdGlvbiBTaWduaW5n
IFNlcnZpY2UxHzAdBgNVBAMTFnJvb3QtY2EtcHJvZHVjdGlvbi1hbW+CAQEwMwYJ
YIZIAYb4QgEEBCYWJGh0dHA6Ly9hZGRvbnMubW96aWxsYS5vcmcvY2EvY3JsLnBl
bTANBgkqhkiG9w0BAQwFAAOCAgEArde/fdjb7TE0eH7Ij7xU4JbcSyhY3cQhVYCw
Fg+Q/2pj+NAfazcjUuLWA0Y/YZs9HOx6j+ZAqO4C/xfMP4RDs9IypxvzHDU6SXgD
RK6uOKtS07HXLcXgFUBvJEQhbT/h5+IQOA4/GcpCshfD6iyiBBi+IocR+tnKPCuZ
T3m1t60Eja/MkPKG/Gx8vSodHvlTTsJ2GzjUEANveCZOnlAdp9fjTvFZny9qqnbg
sfVbuTqKndbCFW5QLXfkna6jBqMrY0+CpMYY2oJ5gwpHbE/7hhukjxGCTcpv7r/O
M53bb/DZnybDlLLepacljvz7DBA1O1FFtEhf9MR+vyvmBpniAyKQhqG2hsVGurE1
nBcE+oteZWar2lMp6+etDAb9DRC+jZv0aEQs2o/qQwyD8AGquLgBsJq5Jz3gGxzn
4r3vGu2lV8VdzIm0C8sOFSWTmTZxQmJbF8xSsQBojnsvEah4DPER+eAt6qKolaWe
s4drJQjzFyC7HJn2VqalpCwbe9CdMB7eRqzeP6GujJBi80/gx0pAysUtuKKpH5IJ
WbXAOszfrjb3CaHafYZDnwPoOfj74ogFzjt2f54jwnU+ET/byfjZ7J8SLH316C1V
HrvFXcTzyMV4aRluVPjPg9x1G58hMIbeuT4GpwQUNdJ9uL8t65v0XwG2t6Y7jpRO
sFVxBtgxggNXMIIDUwIBATCB0DCBxTELMAkGA1UEBhMCVVMxHDAaBgNVBAoTE01v
emlsbGEgQ29ycG9yYXRpb24xLzAtBgNVBAsTJk1vemlsbGEgQU1PIFByb2R1Y3Rp
b24gU2lnbmluZyBTZXJ2aWNlMTEwLwYDVQQDEyhwcm9kdWN0aW9uLXNpZ25pbmct
Y2EuYWRkb25zLm1vemlsbGEub3JnMTQwMgYJKoZIhvcNAQkBFiVzZXJ2aWNlcy1v
cHMrYWRkb25zaWduaW5nQG1vemlsbGEuY29tAgYBVpobWVwwCQYFKw4DAhoFAKBd
MBgGCSqGSIb3DQEJAzELBgkqhkiG9w0BBwEwHAYJKoZIhvcNAQkFMQ8XDTE2MDgx
NzIwMDQ1OFowIwYJKoZIhvcNAQkEMRYEFAxlGvNFSx+Jqj70haE8b7UZk+2GMA0G
CSqGSIb3DQEBAQUABIICADsDlrucYRgwq9o2QSsO6X6cRa5Zu6w+1n07PTIyc1zn
Pi1cgkkWZ0kZBHDrJ5CY33yRQPl6I1tHXaq7SkOSdOppKhpUmBiKZxQRAZR21QHk
R3v1XS+st/o0N+0btv3YoplUifLIwtH89oolxqlStChELu7FuOBretdhx/z12ytA
EhIIS53o/XjDL7XKJbQA02vzOtOC/Eq6p8BI7F3y6pvtmJIRkeGv+u6ssJa6g5q8
74w8hHXaH94Z9+hDPqjNWlsXJHgPdAKiEjzDz9oLkvDyX4Pd8JMK5ILskirpG+hj
Q8jkTc5oYwyuSlBAUTGxW6ZbuOrtfVZvOVtRL/ixuiFiVlJ+JOQOxrtK19ukamsI
iacFlbLgiA7w0HCtm2DsT9aL67/1e4rJ0lv0MjnQYUMmKQy7g0Gd3+nQPU9pn+Lf
Z/UmSNWiJ8Csc/seDMyzT6jrzcGPfoSVaUowH0wGrI9If1snwcr+mMg7dWRGf1fm
y/dcVSzed0ax4LqDmike1EshU+51cKWWlnhyNHK4KH+0fNsBQ0c6clrFpGx9MPmV
YXie6C+LWkh5x12RU0sJt/SmSZV6q9VliIkX+yY3jBrC/pKgRahtcIyq46Da1E6K
lc15Euur3NfGow+nott0Z8XutpYdK/2vBKcIh9JOdkd+oe6pcIP6hnhHRp53wqmG
-----END PKCS7-----`

var FirefoxAddonRootCert = []byte(`
-----BEGIN CERTIFICATE-----
MIIGYTCCBEmgAwIBAgIBATANBgkqhkiG9w0BAQwFADB9MQswCQYDVQQGEwJVUzEc
MBoGA1UEChMTTW96aWxsYSBDb3Jwb3JhdGlvbjEvMC0GA1UECxMmTW96aWxsYSBB
TU8gUHJvZHVjdGlvbiBTaWduaW5nIFNlcnZpY2UxHzAdBgNVBAMTFnJvb3QtY2Et
cHJvZHVjdGlvbi1hbW8wHhcNMTUwMzE3MjI1MzU3WhcNMjUwMzE0MjI1MzU3WjB9
MQswCQYDVQQGEwJVUzEcMBoGA1UEChMTTW96aWxsYSBDb3Jwb3JhdGlvbjEvMC0G
A1UECxMmTW96aWxsYSBBTU8gUHJvZHVjdGlvbiBTaWduaW5nIFNlcnZpY2UxHzAd
BgNVBAMTFnJvb3QtY2EtcHJvZHVjdGlvbi1hbW8wggIgMA0GCSqGSIb3DQEBAQUA
A4ICDQAwggIIAoICAQC0u2HXXbrwy36+MPeKf5jgoASMfMNz7mJWBecJgvlTf4hH
JbLzMPsIUauzI9GEpLfHdZ6wzSyFOb4AM+D1mxAWhuZJ3MDAJOf3B1Rs6QorHrl8
qqlNtPGqepnpNJcLo7JsSqqE3NUm72MgqIHRgTRsqUs+7LIPGe7262U+N/T0LPYV
Le4rZ2RDHoaZhYY7a9+49mHOI/g2YFB+9yZjE+XdplT2kBgA4P8db7i7I0tIi4b0
B0N6y9MhL+CRZJyxdFe2wBykJX14LsheKsM1azHjZO56SKNrW8VAJTLkpRxCmsiT
r08fnPyDKmaeZ0BtsugicdipcZpXriIGmsZbI12q5yuwjSELdkDV6Uajo2n+2ws5
uXrP342X71WiWhC/dF5dz1LKtjBdmUkxaQMOP/uhtXEKBrZo1ounDRQx1j7+SkQ4
BEwjB3SEtr7XDWGOcOIkoJZWPACfBLC3PJCBWjTAyBlud0C5n3Cy9regAAnOIqI1
t16GU2laRh7elJ7gPRNgQgwLXeZcFxw6wvyiEcmCjOEQ6PM8UQjthOsKlszMhlKw
vjyOGDoztkqSBy/v+Asx7OW2Q7rlVfKarL0mREZdSMfoy3zTgtMVCM0vhNl6zcvf
5HNNopoEdg5yuXo2chZ1p1J+q86b0G5yJRMeT2+iOVY2EQ37tHrqUURncCy4uwIB
A6OB7TCB6jAMBgNVHRMEBTADAQH/MA4GA1UdDwEB/wQEAwIBBjAWBgNVHSUBAf8E
DDAKBggrBgEFBQcDAzCBkgYDVR0jBIGKMIGHoYGBpH8wfTELMAkGA1UEBhMCVVMx
HDAaBgNVBAoTE01vemlsbGEgQ29ycG9yYXRpb24xLzAtBgNVBAsTJk1vemlsbGEg
QU1PIFByb2R1Y3Rpb24gU2lnbmluZyBTZXJ2aWNlMR8wHQYDVQQDExZyb290LWNh
LXByb2R1Y3Rpb24tYW1vggEBMB0GA1UdDgQWBBSzvOpYdKvhbngqsqucIx6oYyyX
tzANBgkqhkiG9w0BAQwFAAOCAgEAaNSRYAaECAePQFyfk12kl8UPLh8hBNidP2H6
KT6O0vCVBjxmMrwr8Aqz6NL+TgdPmGRPDDLPDpDJTdWzdj7khAjxqWYhutACTew5
eWEaAzyErbKQl+duKvtThhV2p6F6YHJ2vutu4KIciOMKB8dslIqIQr90IX2Usljq
8Ttdyf+GhUmazqLtoB0GOuESEqT4unX6X7vSGu1oLV20t7t5eCnMMYD67ZBn0YIU
/cm/+pan66hHrja+NeDGF8wabJxdqKItCS3p3GN1zUGuJKrLykxqbOp/21byAGog
Z1amhz6NHUcfE6jki7sM7LHjPostU5ZWs3PEfVVgha9fZUhOrIDsyXEpCWVa3481
LlAq3GiUMKZ5DVRh9/Nvm4NwrTfB3QkQQJCwfXvO9pwnPKtISYkZUqhEqvXk5nBg
QCkDSLDjXTx39naBBGIVIqBtKKuVTla9enngdq692xX/CgO6QJVrwpqdGjebj5P8
5fNZPABzTezG3Uls5Vp+4iIWVAEDkK23cUj3c/HhE+Oo7kxfUeu5Y1ZV3qr61+6t
ZARKjbu1TuYQHf0fs+GwID8zeLc2zJL7UzcHFwwQ6Nda9OJN4uPAuC/BKaIpxCLL
26b24/tRam4SJjqpiq20lynhUrmTtt6hbG3E1Hpy3bmkt2DYnuMFwEx2gfXNcnbT
wNuvFqc=
-----END CERTIFICATE-----`)

// sign a document with openssl and verify the signature with pkcs7.
// this uses a chain of root, intermediate and signer cert, where the
// intermediate is added to the certs but the root isn't.
func TestSignWithOpenSSLAndVerify(t *testing.T) {
	content := []byte(`
A ship in port is safe,
but that's not what ships are built for.
-- Grace Hopper`)
	// write the content to a temp file
	tmpContentFile, err := ioutil.TempFile("", "TestSignWithOpenSSLAndVerify_content")
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile(tmpContentFile.Name(), content, 0755)
	sigalgs := []x509.SignatureAlgorithm{
		x509.SHA1WithRSA,
		x509.SHA256WithRSA,
		x509.SHA512WithRSA,
		x509.ECDSAWithSHA1,
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
			// write the intermediate cert to a temp file
			tmpInterCertFile, err := ioutil.TempFile("", "TestSignWithOpenSSLAndVerify_intermediate")
			if err != nil {
				t.Fatal(err)
			}
			fd, err := os.OpenFile(tmpInterCertFile.Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				t.Fatal(err)
			}
			pem.Encode(fd, &pem.Block{Type: "CERTIFICATE", Bytes: interCert.Certificate.Raw})
			fd.Close()
			for _, sigalgsigner := range sigalgs {
				signerCert, err := createTestCertificateByIssuer("PKCS7 Test Signer Cert", interCert, sigalgsigner, false)
				if err != nil {
					t.Fatalf("test %s/%s/%s: cannot generate signer cert: %s", sigalgroot, sigalginter, sigalgsigner, err)
				}

				// write the signer cert to a temp file
				tmpSignerCertFile, err := ioutil.TempFile("", "TestSignWithOpenSSLAndVerify_signer")
				if err != nil {
					t.Fatal(err)
				}
				fd, err = os.OpenFile(tmpSignerCertFile.Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
				if err != nil {
					t.Fatal(err)
				}
				pem.Encode(fd, &pem.Block{Type: "CERTIFICATE", Bytes: signerCert.Certificate.Raw})
				fd.Close()

				// write the signer key to a temp file
				tmpSignerKeyFile, err := ioutil.TempFile("", "TestSignWithOpenSSLAndVerify_key")
				if err != nil {
					t.Fatal(err)
				}
				fd, err = os.OpenFile(tmpSignerKeyFile.Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
				if err != nil {
					t.Fatal(err)
				}
				var derKey []byte
				priv := *signerCert.PrivateKey
				switch priv := priv.(type) {
				case *rsa.PrivateKey:
					derKey = x509.MarshalPKCS1PrivateKey(priv)
					pem.Encode(fd, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: derKey})
				case *ecdsa.PrivateKey:
					derKey, err = x509.MarshalECPrivateKey(priv)
					if err != nil {
						t.Fatal(err)
					}
					pem.Encode(fd, &pem.Block{Type: "EC PRIVATE KEY", Bytes: derKey})
				}
				fd.Close()

				// write the root cert to a temp file
				tmpSignedFile, err := ioutil.TempFile("", "TestSignWithOpenSSLAndVerify_signature")
				if err != nil {
					t.Fatal(err)
				}
				// call openssl to sign the content
				opensslCMD := exec.Command("openssl", "smime", "-sign", "-nodetach",
					"-in", tmpContentFile.Name(), "-out", tmpSignedFile.Name(),
					"-signer", tmpSignerCertFile.Name(), "-inkey", tmpSignerKeyFile.Name(),
					"-certfile", tmpInterCertFile.Name(), "-outform", "PEM")
				out, err := opensslCMD.CombinedOutput()
				if err != nil {
					t.Fatalf("test %s/%s/%s: openssl command failed with %s: %s", sigalgroot, sigalginter, sigalgsigner, err, out)
				}

				// verify the signed content
				pemSignature, err := ioutil.ReadFile(tmpSignedFile.Name())
				if err != nil {
					t.Fatal(err)
				}
				derBlock, _ := pem.Decode(pemSignature)
				if derBlock == nil {
					break
				}
				p7, err := Parse(derBlock.Bytes)
				if err != nil {
					t.Fatalf("Parse encountered unexpected error: %v", err)
				}
				if err := p7.VerifyWithChain(truststore); err != nil {
					t.Fatalf("Verify failed with error: %v", err)
				}
				// Verify the certificate chain to make sure the identified root
				// is the one we expect
				ee := getCertFromCertsByIssuerAndSerial(p7.Certificates, p7.Signers[0].IssuerAndSerialNumber)
				if ee == nil {
					t.Fatalf("No end-entity certificate found for signer")
				}
				chains, err := verifyCertChain(ee, p7.Certificates, truststore, time.Now())
				if err != nil {
					t.Fatal(err)
				}
				if len(chains) != 1 {
					t.Fatalf("Expected to find one chain, but found %d", len(chains))
				}
				if len(chains[0]) != 3 {
					t.Fatalf("Expected to find three certificates in chain, but found %d", len(chains[0]))
				}
				if chains[0][0].Subject.CommonName != "PKCS7 Test Signer Cert" {
					t.Fatalf("Expected to find EE certificate with subject 'PKCS7 Test Signer Cert', but found '%s'", chains[0][0].Subject.CommonName)
				}
				if chains[0][1].Subject.CommonName != "PKCS7 Test Intermediate Cert" {
					t.Fatalf("Expected to find intermediate certificate with subject 'PKCS7 Test Intermediate Cert', but found '%s'", chains[0][1].Subject.CommonName)
				}
				if chains[0][2].Subject.CommonName != "PKCS7 Test Root CA" {
					t.Fatalf("Expected to find root certificate with subject 'PKCS7 Test Root CA', but found '%s'", chains[0][2].Subject.CommonName)
				}
				os.Remove(tmpSignerCertFile.Name()) // clean up
				os.Remove(tmpSignerKeyFile.Name())  // clean up
			}
			os.Remove(tmpInterCertFile.Name()) // clean up
		}
	}
	os.Remove(tmpContentFile.Name()) // clean up
}
