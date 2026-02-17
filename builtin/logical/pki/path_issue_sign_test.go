// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPathIssueSign_KeyTypeAny is a regression test.  At one point, the signature bits
// were determined based on the keyType of the (CSR) key being signed, not by the (CA)
// signer key.  We also need to make sure we correctly set default bits for each key type.
func TestPathIssueSign_KeyTypeAny(t *testing.T) {
	// Try (Parent) Issuers of Several Different Types:
	// RSA2048; RSA3072; RSA4096; RSA8192; EC224; EC256; EC384; EC521; ED25519 (ignore keysize)
	// And Signature Bits of:
	// 224, 256, 384, 512 (+ Unset)
	// Try (Child) Issued Certificates of Each of Those Same (Key) Type (with a CSR)
	keyTypeOptions := []struct {
		name               string
		keyType            string
		keySize            int
		defaultSigBits     int
		successfulBitSizes []int
		generatedCaBundle  string
		generatedLeafCsr   string
	}{
		{"rsa-2048", "rsa", 2048, 256, []int{0, 256, 384, 512}, rsa2048CaAndKey, rsa2048leafCsr},
		{"rsa-3072", "rsa", 3072, 256, []int{0, 256, 384, 512}, rsa3072CaAndKey, rsa3072leafCsr},
		{"rsa-4096", "rsa", 4096, 256, []int{0, 256, 384, 512}, rsa4096CaAndKey, rsa4096leafCsr},
		{"rsa-8192", "rsa", 8192, 256, []int{0, 256, 384, 512}, rsa8192CaAndKey, rsa8192leafCsr},
		{"ec-224", "ec", 224, 256, []int{0, 256}, ec224CaAndKey, ec224leafCsr},
		{"ec-256", "ec", 256, 256, []int{0, 256}, ec256CaAndKey, ec256leafCsr},
		{"ec-384", "ec", 384, 384, []int{0, 384}, ec384CaAndKey, ec384leafCsr},
		{"ec-521", "ec", 521, 512, []int{0, 512}, ec512CaAndKey, ec512leafCsr},
		{"ed25519", "ed25519", 0, 512, []int{0, 512}, ed25519CaAndKey, ed25519leafCsr},
	}
	signatureBitOptions := []int{0, 224, 256, 384, 512, 1117}

	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	for _, parentKeyType := range keyTypeOptions {
		resp, err := CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
			"pem_bundle": parentKeyType.generatedCaBundle,
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected ca info")
		}
		if resp.IsError() {
			t.Fatal("expected successful import", resp.Error())
		}
		issuers := resp.Data["imported_issuers"].([]string)
		issuerId := issuers[0]
		resp, err = CBWrite(b, s, "issuer/"+issuerId, map[string]interface{}{
			"issuer_name": parentKeyType.name,
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected ca info updated")
		}
		if resp.IsError() {
			t.Fatal("expected successful update", resp.Error())
		}

		for _, signatureBitOption := range signatureBitOptions {
			roleName := "issuer_" + parentKeyType.name + "_sigBits_" + strconv.Itoa(signatureBitOption)

			// Now Make a Role with Key-Type Any, Set Signature-Bits:
			roleResp, err := CBWrite(b, s, "roles/"+roleName, map[string]interface{}{
				"allowed_domains":    "foobar.com",
				"allow_bare_domains": true,
				"max_ttl":            "2h",
				"key_type":           "any",
				"issuer_ref":         parentKeyType.name,
				"signature_bits":     signatureBitOption,
			})
			if err != nil { // We should never fail to create a role because of misconfigured Signature bits if the
				// configuration might be valid for a different issuer, since issuers are updated separately.
				t.Fatal(fmt.Errorf("test failed creating role %v: %v", roleName, err))
			}
			signingFailureExpected := false
			defaultSizeOverride := false
			if slices.Contains(parentKeyType.successfulBitSizes, signatureBitOption) {
				require.Equal(t, len(roleResp.Warnings), 0)
			} else {
				if parentKeyType.keyType == "rsa" {
					signingFailureExpected = true
					require.Contains(t, roleResp.Warnings, fmt.Sprintf("The Issuing Certificate %v for this role has a key algorithm, %v, incompatible with the set role signature bits, %d", parentKeyType.name, parentKeyType.keyType, signatureBitOption))
				} else {
					defaultSizeOverride = true
				}
			}

			// For Each Key-Type, Generate a CSR and try to have it signed by the Role
			for _, childKeyTypeOption := range keyTypeOptions {

				resp, err = CBWrite(b, s, "sign/"+roleName, map[string]interface{}{
					"common_name": "foobar.com",
					"csr":         childKeyTypeOption.generatedLeafCsr,
				})

				if signingFailureExpected {
					require.Error(t, err, "expected signing failure", roleName, childKeyTypeOption.keyType, childKeyTypeOption.keySize)
				} else {
					if err != nil {
						t.Fatal(fmt.Errorf("test failed signing csr with keyType %v size %d with role %v: %v", childKeyTypeOption.keyType, childKeyTypeOption.keySize, roleName, err))
					}
					if resp == nil {
						t.Fatal(fmt.Errorf("test role %v didn't give cert response to attemp to sign CSR with %v keyType", roleName, childKeyTypeOption))
					}

					rawCert := resp.Data["certificate"].(string)
					trimmedRawCert, _ := strings.CutPrefix(rawCert, "-----BEGIN CERTIFICATE-----\n")
					moreTrimmedCert, _ := strings.CutSuffix(trimmedRawCert, "-----END CERTIFICATE-----")
					cleanCert := strings.ReplaceAll(moreTrimmedCert, "\n", "")
					certBytes, err := base64.StdEncoding.DecodeString(cleanCert)
					if err != nil {
						require.NoError(t, err, "failed to decode certificate")
					}
					cert, err := x509.ParseCertificate(certBytes)
					if err != nil {
						require.NoError(t, err, "failed to parse certificate")
					}

					// Signature is Truncated, So We have to Look At Type To Get Digest Length
					resultingLength := signingBits(cert.SignatureAlgorithm)
					if signatureBitOption == 0 || defaultSizeOverride {
						require.Equal(t, parentKeyType.defaultSigBits, resultingLength, "signature length was not default size", "parent key type", parentKeyType.name, "signature bits set", signatureBitOption, "default signature bits", parentKeyType.defaultSigBits, "childKeyType", childKeyTypeOption.name)
					} else {
						require.Equal(t, signatureBitOption, resultingLength, "signature length was not what was set", "role", roleName, "parent key type", parentKeyType.name, "signature bits set", signatureBitOption, "childKeyType", childKeyTypeOption.name)
					}
				}
			}

		}
	}
}

func signingBits(alg x509.SignatureAlgorithm) int {
	switch alg {
	case x509.SHA256WithRSA, x509.SHA256WithRSAPSS, x509.ECDSAWithSHA256:
		return 256
	case x509.SHA384WithRSA, x509.SHA384WithRSAPSS, x509.ECDSAWithSHA384:
		return 384
	case x509.SHA512WithRSA, x509.SHA512WithRSAPSS, x509.ECDSAWithSHA512, x509.PureEd25519:
		return 512
	default:
		return 0
	}
}

// RSA keys (and certificates to a lesser degree) are very slow to generate which is
// inconvenient for testing; so these have been pre-generated.  They are valid for
// 100 years; so shouldn't expire.  They were generated with the following code (and
// tuning the mount for a similarly long ttl):
//
//	resp, err := CBWrite(b, s, "root/generate/exported", map[string]interface{}{
//		"common_name": "foobar.com",
//		"issuer_name": parentKeyType.name,
//		"key_name":    parentKeyType.name,
//		"key_type":    parentKeyType.keyType,
//		"key_bits":    parentKeyType.keySize,
//		"ttl":         "1000000h",
//	})
const rsa2048CaAndKey = `-----BEGIN CERTIFICATE-----
MIIDNDCCAhygAwIBAgIURfqw7VetXbOAIRomrgKZvUPwvCUwDQYJKoZIhvcNAQEL
BQAwFTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTQyMDI3MjVaGA8yMTM5
MTIxNDEyMjc1NVowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBALpBu3/YskXGW2WElPbXfUwibv9+vqMPlAud3iTh
an/RObRBMWK0heow4J6YY63z7PtF10XFxlwiXlv3HOrfNttvtuq4Ma1aC1JJDjZp
OAqniyoojtjOgH77TMDWEM1OIuX6vDwlgfoFySUjuSsrr1eK8yltFvTErqZJZ7RB
/RoHmzvEXeQdYBLV9gY+ku8kqDs7IL6lViKxnqQd33f5+2w4WvpOgVCZSTTyDEMZ
5kpLr/W9UaVx3Lp9WtWaEv2it36Bw+M3zjUxxidjoeHrQDzT1PCPA9mOVWT6s9uI
2WYfmbGAv79gWY3ziX/0XAED3txt/AKcYahBOy1AdrFa0zcCAwEAAaN6MHgwDgYD
VR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFNFSnZqM/NdZ
i+vmKgm7I858m1YgMB8GA1UdIwQYMBaAFNFSnZqM/NdZi+vmKgm7I858m1YgMBUG
A1UdEQQOMAyCCmZvb2Jhci5jb20wDQYJKoZIhvcNAQELBQADggEBAEJU2FLj94US
FpYek0ShL/M38YW0XM6NvAkA03YN1mzIfiGu122XWJwaG3BHGu/FOe3DNxQR65Cr
wT1ZHwzTuV6l54gUVbidGa3f5iscOvt+0PH+XXllVDfJymzUt+i67bC78bK3QEXs
bTIihjo1FH9XgTKTO5aFUgHiciBZqg5Df5REzoZM6s6HA5HWgjbyDTJICausUi+v
jQE1w2h8MrhMwmT+XopUWi2WAqqyTtenfETdA3xQciAZx9OU/4oy8uxjcH3q71/u
7DlylPrvCwOmO/hGUDR0pCZprwZ6fQ/OVYV9p+mtP5NIaRRJyg9bsJ3SyH1kzTGs
US17QT3YFNs=
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAukG7f9iyRcZbZYSU9td9TCJu/36+ow+UC53eJOFqf9E5tEEx
YrSF6jDgnphjrfPs+0XXRcXGXCJeW/cc6t8222+26rgxrVoLUkkONmk4CqeLKiiO
2M6AfvtMwNYQzU4i5fq8PCWB+gXJJSO5KyuvV4rzKW0W9MSupklntEH9GgebO8Rd
5B1gEtX2Bj6S7ySoOzsgvqVWIrGepB3fd/n7bDha+k6BUJlJNPIMQxnmSkuv9b1R
pXHcun1a1ZoS/aK3foHD4zfONTHGJ2Oh4etAPNPU8I8D2Y5VZPqz24jZZh+ZsYC/
v2BZjfOJf/RcAQPe3G38ApxhqEE7LUB2sVrTNwIDAQABAoIBACybsZxc+dVcPGeD
6Wl1Er05QfxPDrle8cYWeS28DxWttnRFaN6K/cepDSLuvHDdCtTjVTuQsoE+efrs
pDBcZXcIunZcxwkNl8iNVqoRaSqkFeBy9kNWsc+3wBovKrcBD7qk4pBFK2wGFrae
Z6q/O69rx/ET/3t/35RT4FJ7u3KQFtuvMujQqYyJvsmZ+GMqqqOS2Xfa9SJS+UDU
cUBNhseLtS4NeA/6gJaBzyFpBd/PTHxQUA/oct/xJr3mTD6whJQZ3vtDdx+GtwSe
lEvZdz12sCcMofbx74PWMefsiMRSFa6nRyrQdRNwjAqNj15bZtFOhlivtuiAu06i
dqkZ7gECgYEA7PMEwnbh9U2mwL2R3dlStUMZEsYWYUwh5QWlgu1pJ90VjPKgrZfJ
lu6l21fl3U23qlc+F19c6tSqHRf4euUugjXIp0h8E+SGc+MhadMgqGPLskUrPuqN
Wr4WV7MAlQqrfx7cFrU9PatugMiNRjVqp54cFU1otFrSoUPByWPFopECgYEAyTtX
Py2AArGPDli2EcWn+HpkYt5j8oip5J12wwU8a+BA3k+pB7G5PqdVwS0Xu5lUaJgG
do0vQZ02jqc8suSTeEZmVnQSU7Wx/+FDo6nWSS+0JVDjwLMuR80o5fmV01kwWjwN
/hPt4A/VrbjP1Bm11jh94NlgTrOaYUtj+Q2jbUcCgYAiNqjuR2ozGGZGmFjSlsm5
gJnDOzUKEYsnXZxbflpbtjGha3tF9Y/XKlhqhpObU9h8USKXD18ETXbOwqJPZH5F
sOxrMy0vViUP4LD3bdPeXKKR+CjZadbFToM9YIxp+ONwdI1E/iB8oh9PmyXDCH2A
/HSDouzGdgLJ5FW79ZsY8QKBgFZqlV0cPQzrE3QlxIp9R1T9un564pEU/2Cd/pJh
fUEWXMUbkIstV1AArGL46mg1wHnqT1w55UFYMkWwq/BnGK1eDjSyQ+yO6pHoOxPd
q5hiVApyYlwuloFfKWEZfa31bz5Q6/FgvZarNigUZavAHsaQG/6jWyhxGKsPpS8f
HD+hAoGAd5E+A9bHwdC1nIB48Ir7uv92jzqD6RTwuxYJteH9WAOSDiVZCGwRmJJT
ZkFDFj3Dr9T/45bJ8DUKv+CH5urX08jGRqCE6fAkaKKpJxcPsZxYtYEqdyoD3Vfl
2aql9ov9ut6BgP3yZhNGtJnW7HqYV2BV8xXPuPeEmSaaTAaioX4=
-----END RSA PRIVATE KEY-----`

const rsa3072CaAndKey = `-----BEGIN CERTIFICATE-----
MIIENDCCApygAwIBAgIUTQvhmoa/czFL4Y00RVeUxWWa7c4wDQYJKoZIhvcNAQEL
BQAwFTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTQyMjI1MDBaGA8yMTM5
MTIxNDE0MjUzMFowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTCCAaIwDQYJKoZIhvcN
AQEBBQADggGPADCCAYoCggGBAKcaQt+AGRbBTkd5zwEeCeNBf72A74kBhs3R8BFK
0nBbbtjySZH05hDd5RNofuvW+zlr8m6sWBvyEZ3adNKhmGEiJyZtBqWMePI5p1yj
NyBZRd0M24uAPDtJYRnlKWkqBu9DIcKyHSGnqlqsQgcbRehs4Vn533WL5AhYd+Vx
EQF8fZjgZfMfJhhXE3Tnib3bS038njmxVMCMCT8+6V4tEGRZzK61FT0NqpDuXxYa
Zm57uEem8MY3TAtENtNrbNyyblQRdWZVmwpmI2oOJ9EMwTUU6Y8RdmuhNpsz4OO+
EWaFwGzg994H/kHda6vpyryxWDrAmwOFg5BxGiz+5KpUgtwzzq17rhEfVZhrdTMT
+NYFX52u3kFqqXXaTjLWGah96mNSC9FuORHGjwYVf7e5rcP2pQBHGfNmXfEVwnBX
tt6iJjp0sSr1ifrvdmcGBVR2mp4fYeqc8BW5EoUibRw+BPiGepikqz2U48rDTAPe
vQCsSDYU422qOA3FmmN8CvZ8IwIDAQABo3oweDAOBgNVHQ8BAf8EBAMCAQYwDwYD
VR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUANpvM+sId1IkzF1eSdBfyl19HFAwHwYD
VR0jBBgwFoAUANpvM+sId1IkzF1eSdBfyl19HFAwFQYDVR0RBA4wDIIKZm9vYmFy
LmNvbTANBgkqhkiG9w0BAQsFAAOCAYEAgUlvaCDvVqj6654dtFuMzoHcr6j8/GVb
xWloujxOmnVq9oCy1Pkhiv3g3iseBEU4ty19qbpTqhThJgbOYWagwz+gHQ9RPgaB
fFd4HzmjfQXhwW0+jJvLRGl5Qn3ydNasGPcxdhtEMVQtXUMYzX02JW8swqeSkk3D
xNnCAucxlcsK/NhmdFoNCebNMGegEkDOTJjVNbbPiuBSEzuySGbzvl+DGmUlshw4
O0UIyxE3+gq3BYT7HK1gQxyf6j3w/pxgO6yYoOWxMERruwTB/6omPXXB80rkQlsS
DzAIWqYOTGSKVLis+D80wL971eMd9fET7riHfXQVjLSvvVfM1If9jzrxs8Hc52TO
OZPgq3VWQLf0I/MhPGeoKfMVhqrX31TJzwJeI+UR3fsqzAZKTrnKeZ4m7WYntgHq
yV/3BMXxCWYljYv5k1mumbzk2uw17wC/KHB2ecTSIAAjMFKP5jTVf1atVy+aLyhO
SLBYd9gGLrRiV2zL61/g5HnYy+TlMoHe
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

const rsa4096CaAndKey = `-----BEGIN CERTIFICATE-----
MIIFNDCCAxygAwIBAgIUArsz+X8+6yiAov3a80qzEa7pF1EwDQYJKoZIhvcNAQEL
BQAwFTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTQyMjQyMDJaGA8yMTM5
MTIxNDE0NDIzMVowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTCCAiIwDQYJKoZIhvcN
AQEBBQADggIPADCCAgoCggIBAM/ltvEBy4eEYhrzjhQd4k7xeRz86p3pdA3xcrba
B7jmxDsqw+UxQExEM6CFnTjJsPABSq3rPDiAOVwLM9NP5H356LZItA2ERaTvgCPg
vYt78bEiAF9Tis676NIUNh0u62ytxJnpD43jl6zSJugjIWEirNFuIiwMWUcXVkGK
Mt9QW7PtkNl7+VtVsR9VnIzBx/C/M4HzIRDIVa8m60vUbADA3u/PnI2MzIjh2ptA
dXFSMWxBK99o0E3DpzvQjjUh2B1d/UTupN6BwShWe+mVH466VuZKVtrN30Efm8xe
Tv4oZ6smbykQ/0bdVi68YiTBSfWzY1ROZ67SX9rA1hTNle6+7iXE51cJkhcGE8vi
QfcQXAe+2uttkxwrnl2qgkfL65FsWceYpL834Qb8UM3jlBI+EttDnansMP7M/+Y9
gu5GLK8+RwXIBqEb+GefsRHIeyUPEUZsHUUanrNXquYjMMSCeEfRe2941GYKVkPI
CKoo7MtTlJCymNpSuqCcBUcAcTje+j0JJyNYaFT2ZXC03Rf8bcrK1wryb7urOnfx
Fg1cGTJOasDFr8CH4G5KvwLcLwZfq8A/yp4jfyRbsQESDi8U7a8RFFXBrXkCLxgN
lO0YxKIJ0wDHpR0IYoyeYxVtVOP774ZGXd1s5VJSf0lke53+Ehzhjlf4P00xaNmb
IPT3AgMBAAGjejB4MA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0G
A1UdDgQWBBT3rlpZOrZu6SsSanJIDCDuaK36WzAfBgNVHSMEGDAWgBT3rlpZOrZu
6SsSanJIDCDuaK36WzAVBgNVHREEDjAMggpmb29iYXIuY29tMA0GCSqGSIb3DQEB
CwUAA4ICAQB0DGuWEBWDX5Kg7nayiKE7QvmiS7fDJstaXcnRbiPs1OitrCRKdDp+
QjtGzLQt/9DYLG+yPsxxLsRJxmTixq1qAFyG2Wka28F1QBAXfUiPBCE9cbHYN2Mk
M+Z9CMeuuZkXmY3wUG8uYOeNkXlrvshiCs6PZRMS6M7C8GL+4VeYld3kVuFqGeJT
p5kFwHTDy5kdqwrW0uLqL5VADPsxCvD1AeAio52J9zEMQoNSpOzb9kk4R3JM0Mse
eXIF6xol1F9IB12/XqjODe1IEFEZv+FVSM1JSyHpaO+vlQHw+OF7N7P7Dw/2YD0l
O1rqeLoLUIKNUuzzOTtoZjqNZ1irQHuN2+TwDJt72dGJDRh7XiRVxpZGRsH1DeH6
nig+VL2iCTPpv68Ge//eLRshRp/4dj0L7nleWw1ve830IOR7Sw4fpoILPFG7NECN
jLEU/hDb+8tsFJQF05m0X43SfxLq3ghexq/iIs09HvQ9AHsX1/Zz+mgyyvtI8e9/
wevNIplA5z3lVNjDUr91S2xZrGHk8zyAkQ925evvIeyiTOaW/zckMi4xyTVHkl0w
oFgquUU6saGjYcT9322ewvdoxsbYSlJR2Kzx0V96duSfY7DfxFgFoOh7oqzM/IJm
dUpagZCAn3I72DHt0uogRa5y3SpJwTlzCQhH1iH3bUFo6y0c/l3jIQ==
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

const rsa8192CaAndKey = `-----BEGIN CERTIFICATE-----
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
lb/Tm0GFPS0=
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

const ec224CaAndKey = `-----BEGIN CERTIFICATE-----
MIIBlDCCAUOgAwIBAgIUDIdGUQAa8UMKCxuIdBDQqMSccOMwCgYIKoZIzj0EAwIw
FTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTcxNDUwMDlaGA8yMTM5MTIx
NzA2NTAzOVowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTBOMBAGByqGSM49AgEGBSuB
BAAhAzoABK/6iwAVzsJYvitjJXJ78YJrRXWExJD/kCmQehI+UrGPAPzHfp81RZuS
c6rHveBdJ3nXlzZCUHvXo3oweDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUw
AwEB/zAdBgNVHQ4EFgQU3jejWn+SRK8KMDn89eXMfiRgiWYwHwYDVR0jBBgwFoAU
3jejWn+SRK8KMDn89eXMfiRgiWYwFQYDVR0RBA4wDIIKZm9vYmFyLmNvbTAKBggq
hkjOPQQDAgM/ADA8AhwTua19D2CHI5Riv5WmBKz3uVLuSG8ZwYSX4ufwAhxDNrLX
2bMvnb8hnlEPI7gVhgB9bq0RGfpYE3v/
-----END CERTIFICATE-----
-----BEGIN EC PRIVATE KEY-----
MGgCAQEEHPybY8Z4tkbQV7HnbJ4Eq6AQpcccLxeA5fXoxbigBwYFK4EEACGhPAM6
AASv+osAFc7CWL4rYyVye/GCa0V1hMSQ/5ApkHoSPlKxjwD8x36fNUWbknOqx73g
XSd515c2QlB71w==
-----END EC PRIVATE KEY-----`

const ec256CaAndKey = `-----BEGIN CERTIFICATE-----
MIIBqTCCAU6gAwIBAgIUaRWkG+4WeYb1bQASXwBDq1hu0eMwCgYIKoZIzj0EAwIw
FTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTcxNTA0NTlaGA8yMTM5MTIx
NzA3MDUyOVowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTBZMBMGByqGSM49AgEGCCqG
SM49AwEHA0IABPQf8311uaA/7ROV2vyjUGcyaJcc5YKshMg2VjDWvG8RrpUtvXLF
/VXW11zngxT97V8dpA1Lj1B1aaKIzx0CTfqjejB4MA4GA1UdDwEB/wQEAwIBBjAP
BgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTd+NrxlY9Q4LisDjWcB574b1pn4jAf
BgNVHSMEGDAWgBTd+NrxlY9Q4LisDjWcB574b1pn4jAVBgNVHREEDjAMggpmb29i
YXIuY29tMAoGCCqGSM49BAMCA0kAMEYCIQDooDJ1p45cwYAUIwNYfU3HO1l1exor
pmNzZu+H0+d2tAIhAPalt/8lIeReeeaDcg2m0bsEBKpjm6VYCe4wWucQ7htW
-----END CERTIFICATE-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIDoSaLJmKY4+qeeMMk7IjkvXoz9NOndpYa0NKlSGdJYoAoGCCqGSM49
AwEHoUQDQgAE9B/zfXW5oD/tE5Xa/KNQZzJolxzlgqyEyDZWMNa8bxGulS29csX9
VdbXXOeDFP3tXx2kDUuPUHVpoojPHQJN+g==
-----END EC PRIVATE KEY-----`

const ec384CaAndKey = `-----BEGIN CERTIFICATE-----
MIIB5jCCAWugAwIBAgIUKLfWtF+8L5myt+fvO9zmpRn6dxowCgYIKoZIzj0EAwMw
FTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTcxNzM1MzBaGA8yMTM5MTIx
NzA5MzYwMFowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTB2MBAGByqGSM49AgEGBSuB
BAAiA2IABDd28Nf40cXJREW4BtC8Ig0dzingZ2jtU1pXS2edHrSQDclwBa8UYZe6
kgpkLNKeEYHhUlYPj98kxKl9E4ekjCn+CxJ5Zx8HPX3iPl01ORlrkI1ZcugNFSjP
aQ0kxB9JV6N6MHgwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYD
VR0OBBYEFE61WPmDUVk+ClC/VZkarqvXKOPsMB8GA1UdIwQYMBaAFE61WPmDUVk+
ClC/VZkarqvXKOPsMBUGA1UdEQQOMAyCCmZvb2Jhci5jb20wCgYIKoZIzj0EAwMD
aQAwZgIxAKju3Se4UrkpUV59tcINouv/l24An3cHyeP2EPk8anRFyAe3P146lPqt
n1+B0Uzj+QIxALQMejPd/Mpps+CLMN13+fOigXLy6jsUXXTt3bUBvrfAC5udGOSC
crZno5EVVS04Wg==
-----END CERTIFICATE-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDBUvahM6eA7z19A3p8Cny55ML83K0W+eIOF7nHLrlwaAlq6toNi3DoZ
SpcLHUOXS/CgBwYFK4EEACKhZANiAAQ3dvDX+NHFyURFuAbQvCINHc4p4Gdo7VNa
V0tnnR60kA3JcAWvFGGXupIKZCzSnhGB4VJWD4/fJMSpfROHpIwp/gsSeWcfBz19
4j5dNTkZa5CNWXLoDRUoz2kNJMQfSVc=
-----END EC PRIVATE KEY-----`

const ec512CaAndKey = `-----BEGIN CERTIFICATE-----
MIICLjCCAZGgAwIBAgIUeLTj4f4z04rgaEJSxBwXjdXxMlUwCgYIKoZIzj0EAwQw
FTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTcxNzM2NDRaGA8yMTM5MTIx
NzA5MzcxNFowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTCBmzAQBgcqhkjOPQIBBgUr
gQQAIwOBhgAEAUr6WhqfUw43B0f6oWKkY3q83mYV2EE0zDhGlg5c1eIDDctD6Dby
jdWPwVjb9yyZ+1en2jveMgJqieN1vA0+Ov5WAb4f4JzIMsaCbEJl68riIh0zx+dg
O5hbisHR7+FWTg8HEQC0/EgSRlgcdBKQnv9tlozDxXHUBX1kQNsMdQoLidtco3ow
eDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUP794
iA+PSFLUTNYbzrB7o9ErmCswHwYDVR0jBBgwFoAUP794iA+PSFLUTNYbzrB7o9Er
mCswFQYDVR0RBA4wDIIKZm9vYmFyLmNvbTAKBggqhkjOPQQDBAOBigAwgYYCQSSi
X/IeqFFKhhyUEg757qkRnT1czdF6HmjaUrpCo5yblkVzHLTO8gxM4e8wDmxL0EFr
+Dxl/3OitC3Ro9D9BghzAkF2t3z3L353jnap1JEAX9gtXQVXoQFGaJ0gZ00iyx2W
pqa9U6pR8qMtDf9uGx01gINsqRruwDkprRXYBYcTU9NoWQ==
-----END CERTIFICATE-----
-----BEGIN EC PRIVATE KEY-----
MIHcAgEBBEIAvjYG3qoX52+P4HWKXNJlLVPwWysnmh/ABpZ+rA+VIyPujIYNfLem
QVaN9QMLKoaHudVh8IRFPOPGxlwOV4tKNeagBwYFK4EEACOhgYkDgYYABAFK+loa
n1MONwdH+qFipGN6vN5mFdhBNMw4RpYOXNXiAw3LQ+g28o3Vj8FY2/csmftXp9o7
3jICaonjdbwNPjr+VgG+H+CcyDLGgmxCZevK4iIdM8fnYDuYW4rB0e/hVk4PBxEA
tPxIEkZYHHQSkJ7/bZaMw8Vx1AV9ZEDbDHUKC4nbXA==
-----END EC PRIVATE KEY-----`

const ed25519CaAndKey = `-----BEGIN CERTIFICATE-----
MIIBaDCCARqgAwIBAgIUKsd3rEygnaf5fiWZUwyh2WlRQxAwBQYDK2VwMBUxEzAR
BgNVBAMTCmZvb2Jhci5jb20wIBcNMjUxMTE3MTc0MDE0WhgPMjEzOTEyMTcwOTQw
NDRaMBUxEzARBgNVBAMTCmZvb2Jhci5jb20wKjAFBgMrZXADIQBpSFBQV9cYAZQT
S0xslNFjdbpgH5rC5yEgUYdgVFWHoqN6MHgwDgYDVR0PAQH/BAQDAgEGMA8GA1Ud
EwEB/wQFMAMBAf8wHQYDVR0OBBYEFHNSUSvGUPM7WKBVhRvojQBnoU35MB8GA1Ud
IwQYMBaAFHNSUSvGUPM7WKBVhRvojQBnoU35MBUGA1UdEQQOMAyCCmZvb2Jhci5j
b20wBQYDK2VwA0EAuNsyYppYWemshxpBinChoxdahOUx8G6uozUSY0VN8HTHcVRP
KhhRZO3a2AzTKjpC/GxD3lrZbPjGpb8LxvEpCA==
-----END CERTIFICATE-----
-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIDALK7hTzLj3lAbHDGIBVGMRI7UnrXnUZaX/wQJABEWj
-----END PRIVATE KEY-----`

// Generating CSRs can take a great deal of time, particularly when
// it requires generating long RSA keys.  Therefore, the following
// CSR were generated with this code, but have been cached here for
// testing.
//
//	goodCr := &x509.CertificateRequest{}
//	var csrKey any
//	var csrPem string
//	for _, childKeyTypeOption := range keyTypeOptions {
//		switch childKeyTypeOption.keyType {
//		case "rsa":
//			csrKey, err = rsa.GenerateKey(rand.Reader, childKeyTypeOption.keySize)
//		case "ec":
//			switch childKeyTypeOption.keySize {
//			case 224:
//				csrKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
//			case 256:
//				csrKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
//			case 384:
//				csrKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
//			case 521:
//				csrKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
//			}
//		case "ed25519":
//			_, csrKey, err = ed25519.GenerateKey(rand.Reader)
//		}
//		require.NoError(t, err, "failed generated key for CSR")
//		csr, err := x509.CreateCertificateRequest(rand.Reader, goodCr, csrKey)
//		require.NoError(t, err, "failed generating csr")
//
//		csrPem = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
//			Type:  "CERTIFICATE REQUEST",
//			Bytes: csr,
//		})))
const rsa2048leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIICRTCCAS0CAQAwADCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANKQ
24KKkOQ+JBhO+s322v4nbioOGVXV7h9YSwM8EIQnZg6N9QxfLTtnx78rfR73yI34
2v+eucxhWmixoFs2/JEye7BiJ5INV9j3eUXQfzmmM9OuLm7RELD+vXcwGKS3so9x
wz57vhttaNM73VylNQCJvoftbhLcKHdkDzqb8fTcDhP+m1v7VKGXtMFF/4wgiJDp
kV3EZX2LHeZ7GcHweIW0m+JgCJBmHH6+mCnhIEiZXemzpJ6MtVLCt/hXdBjUzK/e
LrGkKG9QJEO7HOZKc/KAtVfrROL3PJFIrqJWaD8UWqmuTF5JQOLGqG038QRE+Zcw
0woON+Ssq/Ncl5X5wXsCAwEAAaAAMA0GCSqGSIb3DQEBCwUAA4IBAQAUfcsyYePl
8LdE6V0HVEYLVCfN4n+MVAjWBUjAWJZdimfxTGqMaAUBcvuhyKI6kYRcIXTgptQm
DaXK0BX6EjEbindmA1BpOEWgZyb+heE2OV5VtMKyX/+fnJ2ZJ3VY7aAACsx91VNJ
9pMdGKrcKewP6E2pbS4yryx8NbbZRLR46LBFwfFIcoPfFgrSGmv6mjnNNzLXdF+w
nKJyPhjP5c3Aq54kkDwzu68JEZxapwL3Zm49sWoNb0gnkm/XtgkkRrCYamsUnqqq
LSjy4Ov7yvK31HAULfkyxXgVxNNdxAsyUGXT/pVc6F9dtqq8Ni0s/KmvKk+ghUi5
w49o1WDkKU1p
-----END CERTIFICATE REQUEST-----`

const rsa3072leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIIDRTCCAa0CAQAwADCCAaIwDQYJKoZIhvcNAQEBBQADggGPADCCAYoCggGBALXT
IIejmww8yIh3/ly7d+/B5xowUhk9ij4sg5LAyxQtQa3ZWYNjipdGBG14G0nspG0e
Z1pevHDuhmODH7xi6shwo+ng8phTDYoW5/9aJn10lyHBsvuCic6kRuSZTqpS3NFk
8l2y4vHKYS2ZhcGlBujO4m/cbafkUN7v/a4L0h2fYq7EG/rGRud79OMnFaVhFqaR
macqukHqzhxo9Wj2y4iqdDiQpL0mm84ymKE27fSdfBDdkATKN3mUwfkkqBDM4bhM
ulX0saqjqxLB1IA2NBmt3jKO0HH36e7UlwDmBnEi2dpdMwvzLntCcAECUqSEDrZr
Ht+QZtBOinQ/1XWpYeCSV2XtO4dFfooM2/dSOFm3IK1nohVsIDyzFlRGAPHbT8xC
GbkxY53RWcM/sAKQsnx5kCdefrtQEX1DSOy1w1xEsHvJi1ptuCU0wLZNyq8f9nnr
RRoAFUJmCnjBm+KUOT4p6JVovwbr9I6/5DIq/iVqtXMVzx6NFxDCbv9ok8PXhwID
AQABoAAwDQYJKoZIhvcNAQELBQADggGBAC8raimxwX2lkOxS2Wo3vBYd+VOYD1z3
XXSfCdgqix0tscL09Sz1w0SOttKIoLpgZxMa7/p0fYtjpj1Y4SR2ee54kgRxdHHY
6VGPd6eVDfpYhb+6TEIBdZS2gIC150gOAA7vDrsIVXnLwjTOg6JwQbFgMvthVZ5l
ne+Y2MmSdO1EmKNwYxM3O80HDRIrcnU84OtLxxWGZRFGOaIwC8vFtArnHdHQDw/V
hsgFMov+00mcPUYrLHUy3kLZRm/PEwASgoy6CRhI5UM0LRARV8u+BUMX6B8pjri2
/eRm4euRg4t9jig+jU9oNM+coFKl8ZaHgwoDpmHiUjhGDAprZRwgYB2IhGqbG6TS
RmWkuzQD31uONdTrSIoVcDeVzREmzVDfERQ8KoQQ2nrGqlXgk2QUUmfHt6L+kpZ+
aG/m3iQLttdhotO9DMeZe0qkQ7bgwOnKdOhU9LdjeH1wBHqtKh2IoG+uZ8UOQp+f
uQeKKhf1NSEOyv81mrUvhdVD//6tspRzwg==
-----END CERTIFICATE REQUEST-----`

const rsa4096leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIIERTCCAi0CAQAwADCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAL93
Fn32HImbEe9aD6UqfLIM3VB5pM9LgehMbZUo1NBdkM5urpEe9fG6mP4Or2b3hvn/
HJmV0si4V+IU9DTfc3rmVRUm8weQ5HTyF5pqmbQWbiPpibV6t3BNCuWfLJEfhCW2
pPhyhyTriKIffuo994N3VRW+TXdGFsPvlPAt91miR/DH5GfNZAmDqlxJhgrT70dP
H1eE3PBwPhH8xaAWiccsF5/4R4idDn5DusoB4CWZXSx+FDQdMNTV7XSDOtgZax0R
RKo0nwipBJiIulgp0SsSu4zy8fN2G9EjKkROjiX1YWomqwiK7nrsjjuR5z63Qnwk
pzSGIQteTHSTFH8XoWQQa91pUUy22j3DKwcT9gJEptzt7NiDsbLdqR52dQRyfuYJ
UNDrbB0aC+k2rVBeHPEzvvm3rOJpS6dq+5IU/CCroi0KMt2+ZumMHDe2vGsuNvwS
NvWd4ibk/EZWtiOUTISuEVBOJvtaPGtNmd25CdC/7EN+PpM8vCkW8my0O/hN7frB
TOTR8KTc6LmgBPVXRMUyB6Gss2fI86m38ydZioLGtNCT2ZkAmZzp1w8YNsqdnGvy
kwDIgtjvFR1XtvGv2XO20hIDd7BjJAgDBDoRfU2gpA3RUpQZcyAK4JNexG//gQAV
1cwvInjPSeGsLEhsz9V07mt+OHKEyDJABUecnBWnAgMBAAGgADANBgkqhkiG9w0B
AQsFAAOCAgEAH7PJi7vGU7eghKMeM6Y4r5yYUcQC665AUsgKPax/pyX0E/blPZyL
RhTupB6y4G6He+l4EuoAUcVBx3k/uFHy4jsIElSo7HyqZi9UtUzOn448ZQbGg5fU
zC2qJp+KKEBYnpuAul2e0YDklklOk4oSsuXAb/gYEhiN5foqtn5hUme9AoXBRwyS
hGmJLwkZaJIFDzIlvG1NNfhNnzfTR4o1t/7YMdcp7sJ+Qa/IShc++Xr3vgE6dUmY
3aGQIlWeIkO6qxsuKcfKGt4wwkCRHPR981wT6xHEl3lgeAJg4qalFcCmxOXaWhz+
tppB/Hj1e6FYqfnCBkPe0jjLrenGDYeoSKQkXH4XK2I2VX+kGUPnYN5yaNP9N6hL
31dhrPHknB13mx3UQoGCD+tRglelHceLbgTZteiVC/PjKRN9ZJit98KalHOKQhOv
AFPtQXbBKn3jbIKwYxJE37e3HqQ69eBudFzo0bUwDvEgAsbZXmh6D5hKpEIMRDOd
usCQjur8/Mjwn/c+VI89wAxnBevlFvj6Sy3EDtnSUmnEJsQ8Qhrg2LwISkcodUqg
HUJkn7Yn06drnCTa4WkBBTJOisg39ZvQMKPnjMCxG80CiqgOOZztFH6YeDWfux6b
Tw/AULSSqwxWnEzEwr9CWpeZbgNbXrgGW/iQnCpEwH/X08PsYqSyYIs=
-----END CERTIFICATE REQUEST-----`

const rsa8192leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIIIRTCCBC0CAQAwADCCBCIwDQYJKoZIhvcNAQEBBQADggQPADCCBAoCggQBAMGb
OhJMjmaC239/T3ihtNdRQw+z3iRxR+Uo/jY/FENlBCFaVZpdDYUlvg6ywLVreFPt
frtumnocb2zmkGqovEfZbwpXvPOrohVFI6CtJePkx5ax6CIUuGCT8EpZg02ysQuM
dX8sbojP/3zfrgYqr9+jV0br73xLma9psOVjrSL1L1E+YFh7+ypSZ99Q0QbPj8jW
WwIWXK8nIywHr1hpvlmDsA9QxDxOgmjxUrC7cCUQsvOQWLUxGuvJW+ec7tXGigb+
T1mG9AxadU8Hf+vJaPmRTsDEdK9o76VEec45g4PwNPd51gs7NZq3APq/bOzxPWwj
dyX7kFkP4hSVPYAiQ6IR2MHp8Jr4eRIQrxFGWi/vTCLuixWkXcDzk7U9XeJZWfQz
dfhDhNhk2GjwfdpXCH+fEe9UTJXK3HATqjIH3EaDPk1ba7I1PMFMMEmm4Le/oCXk
fDjZwBxY/CExb6iO1ZIi8Wq9gdsCfLcJTESSpJ4i2ebRrLNk2FI0+cImDt8HyWSR
1wm0FVhJnPxJrR/yRRF/ldjNuF+OAIOAk/k2Emh0b2ioEVybWo1s0vo6+Mo9nIp7
qVKswSAZhkvS5t40/Syrh3845CMSN64C3KFVkQiCN6kRhZXGVCn6OdLBhgCELUn+
95wxUInGBMzN6wRs5MLrHDHyrr99BymL/+jRWbgV+4e8SMPybLPByUbVN1wwvweE
LOcSdHQea+MdFkdbmgeb0/cNLtqx2DpLisUTlvQFLVnX4MFuIJoxbyh+PXYO0oEf
lFTxmw2srf2kNmxYNRDLTu0WDYlmziFiGuppQM8r591+EEM9qdKEanxFD1oKFQ/e
Re0FgxeFAvE6/XuCJMxt1IXYKWzoannHVOlypaH23McPQ9Ci847d5eCaPadvwLO/
fjcyU4MtPdUgNOophSyykJb7+jgK+z+lLlSaJ/EB5ONhidRYrbXrwJiwzcpb66Mb
kJCjpKBfR93wE/LuV1hckAcG2kk5Lehs9KWzs0IKgHsFAWnNd8laIvD46jvejiNP
hnAR2bV4F0+nMLtqvQC60XdLtdvwTolUMYWcpelvdSw9xCaLsA4tLcS+mh7XPTK4
qZNK0cihCDjittgcqe+NU5vdFdmhrL6Fc3SEwZGbUYrDsoSPuVbIXP/cpAjEqW6V
LvahTz76041k5u+TREW972kxVQq0LLcJOAJgaWO3aD74lkQc85XIGI5EUMKpH/2/
R1XFTQ1QJiFtEj1cD4tPFaTZOOvtX7WN5R7pWWS/1UnSddAdz/+8BaFxG+bhtRrX
Vh3GpmucxN7SVTLXGFFSbX2QYmCtXQGuDkGiBYRKJVqUo/VTlrbhJAP5OlSE++XU
/gQ/RMo9wjukZ+GgcnUCAwEAAaAAMA0GCSqGSIb3DQEBCwUAA4IEAQBn5IjJgD9t
xfH2VC6ycEJ9AeWr+ZdC+uayCDbOajnbMNsgE0ssVGcPpq/zeIr1auD3Eme1fG7c
TYapz/tAPQsHVKMXZRB/xB3YdA2VDcGAolp9XWPGpptn3CRIsJfVFWZ5JgJcjoKj
xd72aILT38NqRQ5GXjTQE9ifiJPJ55zVbTOC9TzXonOM0WkjEiCZF5Doc82T2ZKY
Z6NXfENBdGkNydyK/5ZjV7UVrnEHxZI0YkxsfksR0MET7E2uYJCAWLY5mRIiJWAM
QEUWIUA5/JXmv7PcN7eHujDsTRlJcjybENud0PcbOsYJ2EjDeuQvTVMCDKj5oAcb
5dRd0c4fAKgGQ/YnmDUX2fFB9nyYJvWL047RbQQb5ZIyD4ABW3NvCdLbvASOMV7u
4UkdVn3NCQ7drJkY48z1i2o6NnYL2FE93XhHLDiDiKMoEi1su4AMF1S5bSftBwpR
nI2fnkGGMgYs120qYx+/vT4u1QqHeWqg4Qrg6mxKpAGjsqAmaSy7117PnBVHOLVl
krbGmANBTUV/z+w9hkOnnv9iMXd+yG5z0y+tZ9Z58iuje2CgmbDZJPfuja5OYg8f
hfxfmKI803F2bYgDV0MWx1vn6+PJqQwm1utXneTzgHPOFR/SoQXIE+6oxAF0b6w2
02OHehnYmeIvMKlmmzw5i/0I+0/wfOOqunK4mZtiBRW0d43L20j/KjNHSiFoMkfy
gD1j06hmCbvS6VJxrry4OA2naIuDY8ZSxYlbAb5ichpMDx66UjPa592io0OktXte
F+VV5irlRUaceeylfEykngE6EtqtgeiHnRct0bOm7CWuivhltVSXL7dwBsH5i/mc
KnkSKRPGsFlXDS6zrcUYsu3N4AIFBTffCuVn+rJIBM33BS4ncLmMAgfDrPxt1YNU
xw2d91GV3f0eu1XyqQSJJ26lj79aa6buHeEGromaPnfKaPqTAEKL5OqJ76HpXi6I
xdEHu6NrbXzNxWzRE7Hyqkw/AcTg4SJOjxkYJSuelpPaEUfNti0gaSV6ZBAUkJnv
+RZvtTnDm9clIh/6vnjxIjpO3yof2PrKlkOqfWCcdmhDbBe1xyVjCeUZY8KyjHP6
YMRDISJvwcPW8zjBGBzNyLWSkLs8wf2wm+chSsTymQzJSahQsosXWvgezuWbHIME
rek1nw6LtQ9RfQmN3bIrR9p1jJaHCC3cIMPVgyA1TJB6/a3dUksaNzkZe2Jb5Jpw
RTNiiaL37XVQG4Dl3EXtEFRZu++4E/AJY6gXTpRPA3vaBQsIMTvmAmPAs0WdX0wb
9bFk7QsRbkb9tdLrPmx71gdM5iBOPyoSC06IrpnCiJIJOj1DuFEjgUV/oQjgWC2/
THW7sEerAulu
-----END CERTIFICATE REQUEST-----`

const ec224leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIGnMFcCAQAwADBOMBAGByqGSM49AgEGBSuBBAAhAzoABDPEiaV7GSzlRPZoY1IE
gbhxZHhP9vqReTNRfAZSCzIjrIAb46j8HvpfruQUCiXnc9F+RDE90LRuoAAwCgYI
KoZIzj0EAwIDQAAwPQIdAOCbImu4+9pUb8BgbHXa13eY3HaKsNdc8m3bizQCHGB0
xwWb7zUBGOrVVyqs/vFhBLq16OTzw3b/wS0=
-----END CERTIFICATE REQUEST-----`

const ec256leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIG7MGICAQAwADBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABNBo3srjDRR9pIKp
Xkg4A4cYLZIdVzY/atAevIWb8tNtIwPu14VOFs1C01YwxQS2TY0iSwgUyAo+BIws
01nGKc6gADAKBggqhkjOPQQDAgNJADBGAiEA76gHPJbAOX85DHwWsgzn+9GpV0t7
qX98i/SIM6z2yIACIQCZfb8cCcx3PL34RzOw4WGDero7QOUvQ9DQp+w036KJRg==
-----END CERTIFICATE REQUEST-----`

const ec384leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIH2MH8CAQAwADB2MBAGByqGSM49AgEGBSuBBAAiA2IABN0yemr+Ij23X7Y1SYbJ
Naf5y2gYlfi9lDG4JG9dcGJ+rc1XQKJ8zesfODjnf0uvJ+tM4ujhZsmz7kF4MKu0
pCdmfF8TfWdvE7ys0wuRlxOgngDsx8sQSbweoFEk1hrRfaAAMAoGCCqGSM49BAMD
A2cAMGQCMFQlHDMjwsKYpfxPmJQDeErzjJynPy+8KAWiJ3zMj39Otdrx0md+SyBG
ibjgw3uXTQIwMdYhMKTEFtF36nqtNTtPODc5MMySHRfxB9Sy0TnjdTyGvdvrHzSr
AWEW03wZUpbW
-----END CERTIFICATE REQUEST-----`

const ec512leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIIBQzCBpQIBADAAMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBctfkWcBYsL0j
xDocFiSLCLujQoMokv+1wBc+J9oWfmYFatpqdd1OlS2A8UdaLc8HIZCIaeV6rUBy
7/LrqZI1Zx0BTZt7Yl3KimUgCLkrq11WLKlQdxuc59ejFUtQR4ci1sR63MGenPgp
/aWUETKoQ8O/Xvur3nkHtVoFD9lmmD2PEeCgADAKBggqhkjOPQQDBAOBjAAwgYgC
QgFpswcXJ+bpXDihwtExKPTTwRIVv7t0JMHgQjolIfSf6T20P9KuuTIYKuls87y9
DGQo5Ku/tj6PkZFwO6VRsbgn5QJCAI6feCOcmwttGXN6YzoajN4458oyD+/UuQ9/
RnoT+li4M/y1QETDjT0h7IIf5lHYulzI+rP7VCE9IU3PS3CeYrdx
-----END CERTIFICATE REQUEST-----`

const ed25519leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MH8wMwIBADAAMCowBQYDK2VwAyEAmX98EACbAIG9PNcwMG6Zo5rkLDrQ02n8ZEXs
A4AbUdKgADAFBgMrZXADQQCmxK2WBPw5S7Bs2xHmBgG3+yF5HxPjVrVloRWxY9/O
vaeNkZybMZSaATniE3IAG04xnddeDR3MZZwDuUamuZcA
-----END CERTIFICATE REQUEST-----`
