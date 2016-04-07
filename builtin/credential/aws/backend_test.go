package aws

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_ConfigClient(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{"access_key": "AKIAJBRHKV6EVTTNXDHA",
		"secret_key": "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",
		"region":     "us-east-1",
	}

	stepCreate := logicaltest.TestStep{
		Operation: logical.CreateOperation,
		Path:      "config/client",
		Data:      data,
	}

	stepUpdate := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data,
	}

	data2 := map[string]interface{}{"access_key": "AKIAJBRHKV6EVTTNXDHA",
		"secret_key": "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",
		"region":     "",
	}
	stepEmptyRegion := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data2,
		ErrorOk:   true,
	}

	data3 := map[string]interface{}{"access_key": "",
		"secret_key": "mCtSM8ZUEQ3mOFVZYPBQkf2sO6F/W7a5TVzrl3Oj",
		"region":     "us-east-1",
	}
	stepInvalidAccessKey := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data3,
		ErrorOk:   true,
	}

	data4 := map[string]interface{}{"access_key": "accesskey",
		"secret_key": "",
		"region":     "us-east-1",
	}
	stepInvalidSecretKey := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data4,
		ErrorOk:   true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: false,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			stepCreate,
			stepEmptyRegion,
			stepInvalidAccessKey,
			stepInvalidSecretKey,
			stepUpdate,
		},
	})

	checkFound, exists, err := b.HandleExistenceCheck(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config/client",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checkFound {
		t.Fatal("existence check not found for path 'config/client'")
	}
	if exists {
		t.Fatal("existence check should have returned 'false' for 'config/client'")
	}

	configClientCreateRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data,
		Storage:   storage,
	}
	_, err = b.HandleRequest(configClientCreateRequest)
	if err != nil {
		t.Fatal(err)
	}

	checkFound, exists, err = b.HandleExistenceCheck(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config/client",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checkFound {
		t.Fatal("existence check not found for path 'config/client'")
	}
	if !exists {
		t.Fatal("existence check should have returned 'true' for 'config/client'")
	}

	clientConfig, err := clientConfigEntry(storage)
	if err != nil {
		t.Fatal(err)
	}
	if clientConfig.AccessKey != data["access_key"] ||
		clientConfig.SecretKey != data["secret_key"] ||
		clientConfig.Region != data["region"] {
		t.Fatalf("bad: expected: %#v\ngot: %#v\n", data, clientConfig)
	}
}

func TestBackend_PathConfigCertificate(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	checkFound, exists, err := b.HandleExistenceCheck(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config/certificate",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checkFound {
		t.Fatal("existence check not found for path 'config/certificate'")
	}
	if exists {
		t.Fatal("existence check should have returned 'false' for 'config/certificate'")
	}

	data := map[string]interface{}{
		"aws_public_cert": `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM3VENDQXEwQ0NRQ1d1a2paNVY0YVp6QUpC
Z2NxaGtqT09BUURNRnd4Q3pBSkJnTlZCQVlUQWxWVE1Sa3cKRndZRFZRUUlFeEJYWVhOb2FXNW5k
Rzl1SUZOMFlYUmxNUkF3RGdZRFZRUUhFd2RUWldGMGRHeGxNU0F3SGdZRApWUVFLRXhkQmJXRjZi
MjRnVjJWaUlGTmxjblpwWTJWeklFeE1RekFlRncweE1qQXhNRFV4TWpVMk1USmFGdzB6Ck9EQXhN
RFV4TWpVMk1USmFNRnd4Q3pBSkJnTlZCQVlUQWxWVE1Sa3dGd1lEVlFRSUV4QlhZWE5vYVc1bmRH
OXUKSUZOMFlYUmxNUkF3RGdZRFZRUUhFd2RUWldGMGRHeGxNU0F3SGdZRFZRUUtFeGRCYldGNmIy
NGdWMlZpSUZObApjblpwWTJWeklFeE1RekNDQWJjd2dnRXNCZ2NxaGtqT09BUUJNSUlCSHdLQmdR
Q2prdmNTMmJiMVZRNHl0LzVlCmloNU9PNmtLL24xTHpsbHI3RDhad3RRUDhmT0VwcDVFMm5nK0Q2
VWQxWjFnWWlwcjU4S2ozbnNzU05wSTZiWDMKVnlJUXpLN3dMY2xuZC9Zb3pxTk5tZ0l5WmVjTjdF
Z2xLOUlUSEpMUCt4OEZ0VXB0M1FieVlYSmRtVk1lZ042UApodmlZdDVKSC9uWWw0aGgzUGExSEpk
c2tnUUlWQUxWSjNFUjExK0tvNHRQNm53dkh3aDYrRVJZUkFvR0JBSTFqCmsrdGtxTVZIdUFGY3ZB
R0tvY1Rnc2pKZW02LzVxb216SnVLRG1iSk51OVF4dzNyQW90WGF1OFFlK01CY0psL1UKaGh5MUtI
VnBDR2w5ZnVlUTJzNklMMENhTy9idXljVTFDaVlRazQwS05IQ2NIZk5pWmJkbHgxRTlycFVwN2Ju
RgpsUmEydjFudE1YM2NhUlZEZGJ0UEVXbWR4U0NZc1lGRGs0bVpyT0xCQTRHRUFBS0JnRWJtZXZl
NWY4TElFL0dmCk1ObVA5Q001ZW92UU9HeDVobzhXcUQrYVRlYnMrazJ0bjkyQkJQcWVacXBXUmE1
UC8ranJkS21sMXF4NGxsSFcKTVhyczNJZ0liNitoVUlCK1M4ZHo4L21tTzBicHI3NlJvWlZDWFlh
YjJDWmVkRnV0N3FjM1dVSDkrRVVBSDVtdwp2U2VEQ09VTVlRUjdSOUxJTll3b3VISXppcVFZTUFr
R0J5cUdTTTQ0QkFNREx3QXdMQUlVV1hCbGs0MHhUd1N3CjdIWDMyTXhYWXJ1c2U5QUNGQk5HbWRY
MlpCclZOR3JOOU4yZjZST2swazlLCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
`,
	}
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config/certificate",
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFound, exists, err = b.HandleExistenceCheck(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config/certificate",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checkFound {
		t.Fatal("existence check not found for path 'config/certificate'")
	}
	if !exists {
		t.Fatal("existence check should have returned 'true' for 'config/certificate'")
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/certificate",
		Storage:   storage,
	})
	expectedCert := `-----BEGIN CERTIFICATE-----
MIIC7TCCAq0CCQCWukjZ5V4aZzAJBgcqhkjOOAQDMFwxCzAJBgNVBAYTAlVTMRkw
FwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYD
VQQKExdBbWF6b24gV2ViIFNlcnZpY2VzIExMQzAeFw0xMjAxMDUxMjU2MTJaFw0z
ODAxMDUxMjU2MTJaMFwxCzAJBgNVBAYTAlVTMRkwFwYDVQQIExBXYXNoaW5ndG9u
IFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6b24gV2ViIFNl
cnZpY2VzIExMQzCCAbcwggEsBgcqhkjOOAQBMIIBHwKBgQCjkvcS2bb1VQ4yt/5e
ih5OO6kK/n1Lzllr7D8ZwtQP8fOEpp5E2ng+D6Ud1Z1gYipr58Kj3nssSNpI6bX3
VyIQzK7wLclnd/YozqNNmgIyZecN7EglK9ITHJLP+x8FtUpt3QbyYXJdmVMegN6P
hviYt5JH/nYl4hh3Pa1HJdskgQIVALVJ3ER11+Ko4tP6nwvHwh6+ERYRAoGBAI1j
k+tkqMVHuAFcvAGKocTgsjJem6/5qomzJuKDmbJNu9Qxw3rAotXau8Qe+MBcJl/U
hhy1KHVpCGl9fueQ2s6IL0CaO/buycU1CiYQk40KNHCcHfNiZbdlx1E9rpUp7bnF
lRa2v1ntMX3caRVDdbtPEWmdxSCYsYFDk4mZrOLBA4GEAAKBgEbmeve5f8LIE/Gf
MNmP9CM5eovQOGx5ho8WqD+aTebs+k2tn92BBPqeZqpWRa5P/+jrdKml1qx4llHW
MXrs3IgIb6+hUIB+S8dz8/mmO0bpr76RoZVCXYab2CZedFut7qc3WUH9+EUAH5mw
vSeDCOUMYQR7R9LINYwouHIziqQYMAkGByqGSM44BAMDLwAwLAIUWXBlk40xTwSw
7HX32MxXYruse9ACFBNGmdX2ZBrVNGrN9N2f6ROk0k9K
-----END CERTIFICATE-----
`
	if resp.Data["aws_public_cert"].(string) != expectedCert {
		t.Fatal("bad: expected:%s\n got:%s\n", expectedCert, resp.Data["aws_public_cert"].(string))
	}
}

func TestBackend_PathImage(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"policies": "p,q,r,s",
		"max_ttl":  "2h",
	}
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "image/ami-abcd123",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "image/ami-abcd123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !policyutil.EquivalentPolicies(strings.Split(data["policies"].(string), ","), resp.Data["policies"].([]string)) {
		t.Fatalf("bad: policies: expected: %#v\ngot: %#v\n", data, resp.Data)
	}

	data["allow_instance_migration"] = true
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "image/ami-abcd123",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "image/ami-abcd123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Data["allow_instance_migration"].(bool) {
		t.Fatal("bad: allow_instance_migration: expected:true got:false\n")
	}

	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "image/ami-abcd123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "image/ami-abcd123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("bad: response: expected:nil actual:%#v\n", resp)
	}
}

func TestBackend_parseRoleTagValue(t *testing.T) {
	tag := "v1:XwuKhyyBNJc=:a=ami-fce3c696:p=root:t=3h0m0s:lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ="
	expected := roleTag{
		Version:  "v1",
		Nonce:    "XwuKhyyBNJc=",
		Policies: []string{"root"},
		MaxTTL:   10800000000000,
		ImageID:  "ami-fce3c696",
		HMAC:     "lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ=",
	}
	actual, err := parseRoleTagValue(tag)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !actual.Equal(&expected) {
		t.Fatalf("err: expected:%#v \ngot: %#v\n", expected, actual)
	}

	tag = "v2:XwuKhyyBNJc=:a=ami-fce3c696:p=root:t=3h0m0s:lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ="
	actual, err = parseRoleTagValue(tag)
	if err == nil {
		t.Fatalf("err: expected error due to invalid role tag version", err)
	}

	tag = "v1:XwuKhyyBNJc=:a=ami-fce3c696:lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ="
	expected = roleTag{
		Version: "v1",
		Nonce:   "XwuKhyyBNJc=",
		ImageID: "ami-fce3c696",
		HMAC:    "lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ=",
	}
	actual, err = parseRoleTagValue(tag)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	tag = "v1:XwuKhyyBNJc=:p=ami-fce3c696:lhvKJAZn8kxNwmPFnyXzmphQTtbXqQe6WG6sLiIf3dQ="
	actual, err = parseRoleTagValue(tag)
	if err == nil {
		t.Fatalf("err: expected error due to missing image ID", err)
	}
}
