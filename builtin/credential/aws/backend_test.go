package aws

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func TestBackend_createRoleTagValue(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"policies": "p,q,r,s",
	}
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "image/abcd-123",
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}

	imageEntry, err := awsImage(storage, "abcd-123")
	if err != nil {
		t.Fatal(err)
	}

	nonce, err := createRoleTagNonce()
	if err != nil {
		t.Fatal(err)
	}
	rTag := &roleTag{
		Version:  "v1",
		AmiID:    "abcd-123",
		Nonce:    nonce,
		Policies: []string{"p", "q", "r"},
		MaxTTL:   200,
	}
	val, err := createRoleTagValue(rTag, imageEntry)
	if err != nil {
		t.Fatal(err)
	}
	if val == "" {
		t.Fatalf("failed to create role tag")
	}
}

func TestBackend_prepareRoleTagPlaintextValue(t *testing.T) {
	nonce, err := createRoleTagNonce()
	if err != nil {
		t.Fatal(err)
	}
	rTag := &roleTag{
		Version: "v1",
		Nonce:   nonce,
		AmiID:   "abcd-123",
	}

	rTag.Version = ""
	val, err := prepareRoleTagPlaintextValue(rTag)
	if err == nil {
		t.Fatalf("expected error for missing version")
	}
	rTag.Version = "v1"

	rTag.Nonce = ""
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err == nil {
		t.Fatalf("expected error for missing nonce")
	}
	rTag.Nonce = nonce

	rTag.AmiID = ""
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err == nil {
		t.Fatalf("expected error for missing ami_id")
	}
	rTag.AmiID = "abcd-123"

	val, err = prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(val, "a=") ||
		!strings.Contains(val, "p=") ||
		!strings.Contains(val, "d=") ||
		!strings.HasPrefix(val, "v1") {
		t.Fatalf("incorrect information in role tag plaintext value")
	}

	rTag.InstanceID = "instance-123"
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(val, "i=") {
		t.Fatalf("missing instance ID in role tag plaintext value")
	}

	rTag.MaxTTL = 200
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(val, "t=") {
		t.Fatalf("missing instance ID in role tag plaintext value")
	}
}

func TestBackend_CreateRoleTagNonce(t *testing.T) {
	nonce, err := createRoleTagNonce()
	if err != nil {
		t.Fatal(err)
	}
	if nonce == "" {
		t.Fatalf("failed to create role tag nonce")
	}
	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		t.Fatal(err)
	}
	if len(nonceBytes) == 0 {
		t.Fatalf("length of role tag nonce is zero")
	}
}

func TestBackend_ConfigTidyIdentities(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"safety_buffer": "60",
	}

	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/tidy/identities",
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_ConfigTidyRoleTags(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"safety_buffer": "60",
	}

	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/tidy/roletags",
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_TidyIdentities(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "tidy/identities",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_TidyRoleTags(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "tidy/roletags",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
}

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

	data3 := map[string]interface{}{"access_key": "",
		"secret_key": "mCtSM8ZUEQ3mOFVZYPBQkf2sO6F/W7a5TVzrl3Oj",
	}
	stepInvalidAccessKey := logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data3,
		ErrorOk:   true,
	}

	data4 := map[string]interface{}{"access_key": "accesskey",
		"secret_key": "",
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
		clientConfig.SecretKey != data["secret_key"] {
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
		Path:      "config/certificate/cert1",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checkFound {
		t.Fatal("existence check not found for path 'config/certificate/cert1'")
	}
	if exists {
		t.Fatal("existence check should have returned 'false' for 'config/certificate/cert1'")
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
		Path:      "config/certificate/cert1",
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFound, exists, err = b.HandleExistenceCheck(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config/certificate/cert1",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !checkFound {
		t.Fatal("existence check not found for path 'config/certificate/cert1'")
	}
	if !exists {
		t.Fatal("existence check should have returned 'true' for 'config/certificate/cert1'")
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/certificate/cert1",
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
	// create a backend
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// create an entry for an AMI
	data := map[string]interface{}{
		"policies": "p,q,r,s",
		"max_ttl":  "120s",
		"role_tag": "VaultRole",
	}
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "image/abcd-123",
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}

	// verify that the entry is created
	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "image/abcd-123",
		Storage:   storage,
	})
	if resp == nil {
		t.Fatalf("expected an image entry for abcd-123")
	}

	// create a role tag
	data2 := map[string]interface{}{
		"policies": "p,q,r,s",
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "image/abcd-123/roletag",
		Storage:   storage,
		Data:      data2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data["tag_key"].(string) == "" ||
		resp.Data["tag_value"].(string) == "" {
		t.Fatalf("invalid tag response: %#v\n", resp)
	}
	tagValue := resp.Data["tag_value"].(string)

	// parse the value and check if the verifiable values match
	rTag, err := parseRoleTagValue(storage, tagValue)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if rTag == nil {
		t.Fatalf("failed to parse role tag")
	}
	if rTag.Version != "v1" ||
		!policyutil.EquivalentPolicies(rTag.Policies, []string{"p", "q", "r", "s"}) ||
		rTag.AmiID != "abcd-123" {
		t.Fatalf("bad: parsed role tag contains incorrect values. Got: %#v\n", rTag)
	}
}

func TestBackend_PathImageTag(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"policies": "p,q,r,s",
		"max_ttl":  "120s",
		"role_tag": "VaultRole",
	}
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "image/abcd-123",
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "image/abcd-123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("failed to find an entry for ami_id: abcd-123")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "image/abcd-123/roletag",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("failed to create a tag on ami_id: abcd-123")
	}
	if resp.IsError() {
		t.Fatalf("failed to create a tag on ami_id: abcd-123: %s\n", resp.Data["error"])
	}
	if resp.Data["tag_value"].(string) == "" {
		t.Fatalf("role tag not present in the response data: %#v\n", resp.Data)
	}
}

func TestBackend_PathBlacklistRoleTag(t *testing.T) {
	// create the backend
	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// create an image entry
	data := map[string]interface{}{
		"ami_id":   "abcd-123",
		"policies": "p,q,r,s",
		"role_tag": "VaultRole",
	}
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "image/abcd-123",
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}

	// create a role tag against an image registered before
	data2 := map[string]interface{}{
		"policies": "p,q,r,s",
	}
	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "image/abcd-123/roletag",
		Storage:   storage,
		Data:      data2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("failed to create a tag on ami_id: abcd-123")
	}
	if resp.IsError() {
		t.Fatalf("failed to create a tag on ami_id: abcd-123: %s\n", resp.Data["error"])
	}
	tag := resp.Data["tag_value"].(string)
	if tag == "" {
		t.Fatalf("role tag not present in the response data: %#v\n", resp.Data)
	}

	// blacklist that role tag
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "blacklist/roletag/" + tag,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("failed to blacklist the roletag: %s\n", tag)
	}

	// read the blacklist entry
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "blacklist/roletag/" + tag,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("failed to read the blacklisted role tag: %s\n", tag)
	}
	if resp.IsError() {
		t.Fatalf("failed to read the blacklisted role tag:%s. Err: %s\n", tag, resp.Data["error"])
	}

	// delete the blacklisted entry
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "blacklist/roletag/" + tag,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// try to read the deleted entry
	tagEntry, err := blacklistRoleTagEntry(storage, tag)
	if err != nil {
		t.Fatal(err)
	}
	if tagEntry != nil {
		t.Fatalf("role tag should not have been present: %s\n", tag)
	}
}

// This is an acceptance test.
func TestBackendAcc_LoginAndWhitelistIdentity(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	// create the backend
	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// get the API credentials from env vars
	clientConfig := map[string]interface{}{
		"access_key": os.Getenv("AWS_ACCESS_KEY"),
		"secret_key": os.Getenv("AWS_SECRET_KEY"),
	}
	if clientConfig["access_key"] == "" ||
		clientConfig["secret_key"] == "" {
		t.Fatalf("credentials not configured")
	}

	// store the credentials
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "config/client",
		Data:      clientConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// write an entry for an ami
	data := map[string]interface{}{
		"policies": "root",
		"max_ttl":  "120s",
	}
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "image/ami-fce3c696",
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}

	loginInput := map[string]interface{}{"pkcs7": `MIAGCSqGSIb3DQEHAqCAMIACAQExCzAJBgUrDgMCGgUAMIAGCSqGSIb3DQEHAaCAJIAEggGmewog
ICJkZXZwYXlQcm9kdWN0Q29kZXMiIDogbnVsbCwKICAicHJpdmF0ZUlwIiA6ICIxNzIuMzEuNjMu
NjAiLAogICJhdmFpbGFiaWxpdHlab25lIiA6ICJ1cy1lYXN0LTFjIiwKICAidmVyc2lvbiIgOiAi
MjAxMC0wOC0zMSIsCiAgImluc3RhbmNlSWQiIDogImktZGUwZjEzNDQiLAogICJiaWxsaW5nUHJv
ZHVjdHMiIDogbnVsbCwKICAiaW5zdGFuY2VUeXBlIiA6ICJ0Mi5taWNybyIsCiAgImFjY291bnRJ
ZCIgOiAiMjQxNjU2NjE1ODU5IiwKICAiaW1hZ2VJZCIgOiAiYW1pLWZjZTNjNjk2IiwKICAicGVu
ZGluZ1RpbWUiIDogIjIwMTYtMDQtMDVUMTY6MjY6NTVaIiwKICAiYXJjaGl0ZWN0dXJlIiA6ICJ4
ODZfNjQiLAogICJrZXJuZWxJZCIgOiBudWxsLAogICJyYW1kaXNrSWQiIDogbnVsbCwKICAicmVn
aW9uIiA6ICJ1cy1lYXN0LTEiCn0AAAAAAAAxggEXMIIBEwIBATBpMFwxCzAJBgNVBAYTAlVTMRkw
FwYDVQQIExBXYXNoaW5ndG9uIFN0YXRlMRAwDgYDVQQHEwdTZWF0dGxlMSAwHgYDVQQKExdBbWF6
b24gV2ViIFNlcnZpY2VzIExMQwIJAJa6SNnlXhpnMAkGBSsOAwIaBQCgXTAYBgkqhkiG9w0BCQMx
CwYJKoZIhvcNAQcBMBwGCSqGSIb3DQEJBTEPFw0xNjA0MDUxNjI3MDBaMCMGCSqGSIb3DQEJBDEW
BBRtiynzMTNfTw1TV/d8NvfgVw+XfTAJBgcqhkjOOAQDBC4wLAIUVfpVcNYoOKzN1c+h1Vsm/c5U
0tQCFAK/K72idWrONIqMOVJ8Uen0wYg4AAAAAAAA`,
		"nonce": "vault-client-nonce",
	}

	loginRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginInput,
	}
	resp, err := b.HandleRequest(loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatalf("login attempt failed")
	}

	// try to login again and see if it succeeds
	resp, err = b.HandleRequest(loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatalf("login attempt failed")
	}

	//instanceID := resp.Auth.Metadata.(map[string]string)["instance_id"]
	instanceID := resp.Auth.Metadata["instance_id"]
	if instanceID == "" {
		t.Fatalf("instance ID not present in the response object")
	}

	loginInput["nonce"] = "changed-vault-client-nonce"
	// try to login again with changed nonce
	resp, err = b.HandleRequest(loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("login attempt should have failed due to client nonce mismatch")
	}

	wlRequest := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "whitelist/identity/" + instanceID,
		Storage:   storage,
	}
	resp, err = b.HandleRequest(wlRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil || resp.Data["ami_id"] != "ami-fce3c696" {
		t.Fatalf("failed to read whitelist identity")
	}

	wlRequest.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(wlRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("failed to delete whitelist identity")
	}

	// try to login again and see if it succeeds
	resp, err = b.HandleRequest(loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatalf("login attempt failed")
	}
}
