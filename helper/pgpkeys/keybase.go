package pgpkeys

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/kr/pretty"
)

const (
	kbPrefix = "keybase:"
)

const leafCA = `
-----BEGIN CERTIFICATE-----
MIIGLTCCBRWgAwIBAgIQCwgUtNIo3MoF9C5vOob1czANBgkqhkiG9w0BAQsFADCB
jzELMAkGA1UEBhMCR0IxGzAZBgNVBAgTEkdyZWF0ZXIgTWFuY2hlc3RlcjEQMA4G
A1UEBxMHU2FsZm9yZDEYMBYGA1UEChMPU2VjdGlnbyBMaW1pdGVkMTcwNQYDVQQD
Ey5TZWN0aWdvIFJTQSBEb21haW4gVmFsaWRhdGlvbiBTZWN1cmUgU2VydmVyIENB
MB4XDTIxMDUyNDAwMDAwMFoXDTIyMDUyNDIzNTk1OVowFTETMBEGA1UEAxMKa2V5
YmFzZS5pbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAK4/kEUdphHV
9NyDUSKS98v195VJxk/w3xIfGrBKzxs1n1SOmZ7Uv7YyoNJga0rNfl7clk8stRuh
vSl+oWZ9fhJsPToR7nXJcrCzNMMNG2uXEXRDXKLe3/HFPTmWiaxIVijW9FXbPKO1
mnIdfWGUNqsBR8fQkyb1f7iZ/qdnyRuQRgdkzyJI2seFd7Ylwp0dbfKzU6CsOW4I
p0+gQNU1VLruSvi4DcqOb636IY38R3pG34KLrjZF0u3JMvyS7uHY1GRyZmtpqOZE
kBC5oQv0NNsNjDkT6pB26nCampOmEyA/lBZA0mBHLsa9UynmF4Aujb0T5gLKUVB+
pi6gGqewHhUCAwEAAaOCAvwwggL4MB8GA1UdIwQYMBaAFI2MXsRUrYrhd+mb+ZsF
4bgBjWHhMB0GA1UdDgQWBBQpHZYhXhc7fp5cXKukDZHHldI9BDAOBgNVHQ8BAf8E
BAMCBaAwDAYDVR0TAQH/BAIwADAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUH
AwIwSQYDVR0gBEIwQDA0BgsrBgEEAbIxAQICBzAlMCMGCCsGAQUFBwIBFhdodHRw
czovL3NlY3RpZ28uY29tL0NQUzAIBgZngQwBAgEwgYQGCCsGAQUFBwEBBHgwdjBP
BggrBgEFBQcwAoZDaHR0cDovL2NydC5zZWN0aWdvLmNvbS9TZWN0aWdvUlNBRG9t
YWluVmFsaWRhdGlvblNlY3VyZVNlcnZlckNBLmNydDAjBggrBgEFBQcwAYYXaHR0
cDovL29jc3Auc2VjdGlnby5jb20wJQYDVR0RBB4wHIIKa2V5YmFzZS5pb4IOd3d3
LmtleWJhc2UuaW8wggF+BgorBgEEAdZ5AgQCBIIBbgSCAWoBaAB2AEalVet1+pEg
MLWiiWn0830RLEF0vv1JuIWr8vxw/m1HAAABeZ8I474AAAQDAEcwRQIhAOU9ZeGf
OgVhjXHsqpkYNiB83ekpLDWpHl7lZlt6+bkaAiAbiePdir4Y1PoO4ME2vfzjugcq
B6Z1ezUrgsomm91vmgB2AN+lXqtogk8fbK3uuF9OPlrqzaISpGpejjsSwCBEXCpz
AAABeZ8I41AAAAQDAEcwRQIgItH+G72BmpUIAgit2gvJKJE68vzSYwIAquNVluIw
l2ICIQD37UsrggnKmrx0+ARHM7sIv2yHdzBtTaL3QELofNVeQgB2ACl5vvCeOTkh
8FZzn2Old+W+V32cYAr4+U1dJlwlXceEAAABeZ8I42EAAAQDAEcwRQIhANOsOxZ6
OASu/haq2scSAoGbNoBIbaKa5P9WCJMNQPLeAiAWErxoeb4UDJal3rjWuMcpP7do
SacTH4DHf7to7xRRqDANBgkqhkiG9w0BAQsFAAOCAQEACIoT1aqhCOnT2wzRdwUF
tpUtaxnBMpi0u4NWLVZkPwU8RWLZ4VAYYA0iuWxlUP7l/DfSbLvM/ECizDt26Pke
Qb9Wj8otgW4WoWHCfRBqyAIgg09SCAumYer0ni4zUhfASe9uSBcCXWhrkv8w1gQk
hgCX3fVwtPYdO7mCx4ZRveHIhlbVfWWf0F1HHP0PmVgJvQFc4IEU8FUjNreALHgG
bBADqMQD/FNb++ofS7C678ba1sajkfhirCa15CWJESW0Y6AWTTD0nttf7LZQziYM
ROrLMwGm75PhW2h7UXo9CB1LAvZNvvx0n7u2hjFrH4laBbDcbzKV6/ao+doIkjxb
1A==
-----END CERTIFICATE-----
`

const interOne = `-----BEGIN CERTIFICATE-----
MIIGEzCCA/ugAwIBAgIQfVtRJrR2uhHbdBYLvFMNpzANBgkqhkiG9w0BAQwFADCB
iDELMAkGA1UEBhMCVVMxEzARBgNVBAgTCk5ldyBKZXJzZXkxFDASBgNVBAcTC0pl
cnNleSBDaXR5MR4wHAYDVQQKExVUaGUgVVNFUlRSVVNUIE5ldHdvcmsxLjAsBgNV
BAMTJVVTRVJUcnVzdCBSU0EgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwHhcNMTgx
MTAyMDAwMDAwWhcNMzAxMjMxMjM1OTU5WjCBjzELMAkGA1UEBhMCR0IxGzAZBgNV
BAgTEkdyZWF0ZXIgTWFuY2hlc3RlcjEQMA4GA1UEBxMHU2FsZm9yZDEYMBYGA1UE
ChMPU2VjdGlnbyBMaW1pdGVkMTcwNQYDVQQDEy5TZWN0aWdvIFJTQSBEb21haW4g
VmFsaWRhdGlvbiBTZWN1cmUgU2VydmVyIENBMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEA1nMz1tc8INAA0hdFuNY+B6I/x0HuMjDJsGz99J/LEpgPLT+N
TQEMgg8Xf2Iu6bhIefsWg06t1zIlk7cHv7lQP6lMw0Aq6Tn/2YHKHxYyQdqAJrkj
eocgHuP/IJo8lURvh3UGkEC0MpMWCRAIIz7S3YcPb11RFGoKacVPAXJpz9OTTG0E
oKMbgn6xmrntxZ7FN3ifmgg0+1YuWMQJDgZkW7w33PGfKGioVrCSo1yfu4iYCBsk
Haswha6vsC6eep3BwEIc4gLw6uBK0u+QDrTBQBbwb4VCSmT3pDCg/r8uoydajotY
uK3DGReEY+1vVv2Dy2A0xHS+5p3b4eTlygxfFQIDAQABo4IBbjCCAWowHwYDVR0j
BBgwFoAUU3m/WqorSs9UgOHYm8Cd8rIDZsswHQYDVR0OBBYEFI2MXsRUrYrhd+mb
+ZsF4bgBjWHhMA4GA1UdDwEB/wQEAwIBhjASBgNVHRMBAf8ECDAGAQH/AgEAMB0G
A1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAbBgNVHSAEFDASMAYGBFUdIAAw
CAYGZ4EMAQIBMFAGA1UdHwRJMEcwRaBDoEGGP2h0dHA6Ly9jcmwudXNlcnRydXN0
LmNvbS9VU0VSVHJ1c3RSU0FDZXJ0aWZpY2F0aW9uQXV0aG9yaXR5LmNybDB2Bggr
BgEFBQcBAQRqMGgwPwYIKwYBBQUHMAKGM2h0dHA6Ly9jcnQudXNlcnRydXN0LmNv
bS9VU0VSVHJ1c3RSU0FBZGRUcnVzdENBLmNydDAlBggrBgEFBQcwAYYZaHR0cDov
L29jc3AudXNlcnRydXN0LmNvbTANBgkqhkiG9w0BAQwFAAOCAgEAMr9hvQ5Iw0/H
ukdN+Jx4GQHcEx2Ab/zDcLRSmjEzmldS+zGea6TvVKqJjUAXaPgREHzSyrHxVYbH
7rM2kYb2OVG/Rr8PoLq0935JxCo2F57kaDl6r5ROVm+yezu/Coa9zcV3HAO4OLGi
H19+24rcRki2aArPsrW04jTkZ6k4Zgle0rj8nSg6F0AnwnJOKf0hPHzPE/uWLMUx
RP0T7dWbqWlod3zu4f+k+TY4CFM5ooQ0nBnzvg6s1SQ36yOoeNDT5++SR2RiOSLv
xvcRviKFxmZEJCaOEDKNyJOuB56DPi/Z+fVGjmO+wea03KbNIaiGCpXZLoUmGv38
sbZXQm2V0TP2ORQGgkE49Y9Y3IBbpNV9lXj9p5v//cWoaasm56ekBYdbqbe4oyAL
l6lFhd2zi+WJN44pDfwGF/Y4QA5C5BIG+3vzxhFoYt/jmPQT2BVPi7Fp2RBgvGQq
6jG35LWjOhSbJuMLe/0CjraZwTiXWTb2qHSihrZe68Zk6s+go/lunrotEbaGmAhY
LcmsJWTyXnW0OMGuf1pGg+pRyrbxmRE1a6Vqe8YAsOf4vmSyrcjC8azjUeqkk+B5
yOGBQMkKW+ESPMFgKuOXwIlCypTPRpgSabuY0MLTDXJLR27lk8QyKGOHQ+SwMj4K
00u/I5sUKUErmgQfky3xxzlIPK1aEn8=
-----END CERTIFICATE-----`

// FetchKeybasePubkeys fetches public keys from Keybase given a set of
// usernames, which are derived from correctly formatted input entries. It
// doesn't use their client code due to both the API and the fact that it is
// considered alpha and probably best not to rely on it.  The keys are returned
// as base64-encoded strings.
func FetchKeybasePubkeys(input []string) (map[string]string, error) {

	client := &http.Client{}
	// client := cleanhttp.DefaultClient()
	if client == nil {
		return nil, errors.New("unable to create an http client")
	}
	usernames := make([]string, 0, len(input))
	u := fmt.Sprintf("https://keybase.io/_/api/1.0/user/lookup.json?usernames=%s&fields=public_keys", strings.Join(usernames, ","))
	resp, err := client.Get(u)
	if err != nil {
		if ue, ok := err.(*url.Error); ok {
			pretty.Print(ue.Err)
		}
		return nil, err
	}
	defer resp.Body.Close()
	// if len(input) == 0 {
	// 	return nil, nil
	// }

	// usernames := make([]string, 0, len(input))
	// for _, v := range input {
	// 	if strings.HasPrefix(v, kbPrefix) {
	// 		usernames = append(usernames, strings.TrimPrefix(v, kbPrefix))
	// 	}
	// }

	// if len(usernames) == 0 {
	// 	return nil, nil
	// }

	// certs := x509.NewCertPool()
	// mTLSConfig := &tls.Config{}
	// certs.AppendCertsFromPEM([]byte(APICA))
	// mTLSConfig.RootCAs = certs

	// tr := &http.Transport{
	// 	TLSClientConfig: mTLSConfig,
	// }

	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	ret := make(map[string]string, len(usernames))
	// u := fmt.Sprintf("https://keybase.io/_/api/1.0/user/lookup.json?usernames=%s&fields=public_keys", strings.Join(usernames, ","))

	// client := &http.Client{Transport: tr}
	// client := &http.Client{}

	// _, err = client.Do(req)
	// if err != nil {
	// 	return nil, err
	// }

	// resp, err := client.Get(u)
	// if err != nil {
	// 	if ue, ok := err.(*url.Error); ok {
	// 		pretty.Print(ue.Err)
	// 	}
	// 	return nil, err
	// }
	// defer resp.Body.Close()

	// **** THIS WORKS ****

	// cmd := exec.Command("curl", u)
	// err := cmd.Run()
	// if err != nil {
	// 	return nil, err
	// }

	// type PublicKeys struct {
	// 	Primary struct {
	// 		Bundle string
	// 	}
	// }

	// type LThem struct {
	// 	PublicKeys `json:"public_keys"`
	// }

	// type KbResp struct {
	// 	Status struct {
	// 		Name string
	// 	}
	// 	Them []LThem
	// }

	// out := &KbResp{
	// 	Them: []LThem{},
	// }

	// if err := jsonutil.DecodeJSONFromReader(resp.Body, out); err != nil {
	// 	return nil, err
	// }

	// if out.Status.Name != "OK" {
	// 	return nil, fmt.Errorf("got non-OK response: %q", out.Status.Name)
	// }

	// missingNames := make([]string, 0, len(usernames))
	// var keyReader *bytes.Reader
	// serializedEntity := bytes.NewBuffer(nil)
	// for i, themVal := range out.Them {
	// 	if themVal.Primary.Bundle == "" {
	// 		missingNames = append(missingNames, usernames[i])
	// 		continue
	// 	}
	// 	keyReader = bytes.NewReader([]byte(themVal.Primary.Bundle))
	// 	entityList, err := openpgp.ReadArmoredKeyRing(keyReader)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if len(entityList) != 1 {
	// 		return nil, fmt.Errorf("primary key could not be parsed for user %q", usernames[i])
	// 	}
	// 	if entityList[0] == nil {
	// 		return nil, fmt.Errorf("primary key was nil for user %q", usernames[i])
	// 	}

	// 	serializedEntity.Reset()
	// 	err = entityList[0].Serialize(serializedEntity)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error serializing entity for user %q: %w", usernames[i], err)
	// 	}

	// 	// The API returns values in the same ordering requested, so this should properly match
	// 	ret[kbPrefix+usernames[i]] = base64.StdEncoding.EncodeToString(serializedEntity.Bytes())
	// }

	// if len(missingNames) > 0 {
	// 	return nil, fmt.Errorf("unable to fetch keys for user(s) %q from keybase", strings.Join(missingNames, ","))
	// }

	return ret, nil
}
