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

func TestBackend_CreateParseVerifyRoleTag(t *testing.T) {
	// create a backend
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// create an entry for ami
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

	// read the created image entry
	imageEntry, err := awsImage(storage, "abcd-123")
	if err != nil {
		t.Fatal(err)
	}

	// create a nonce for the role tag
	nonce, err := createRoleTagNonce()
	if err != nil {
		t.Fatal(err)
	}
	rTag1 := &roleTag{
		Version:  "v1",
		AmiID:    "abcd-123",
		Nonce:    nonce,
		Policies: []string{"p", "q", "r"},
		MaxTTL:   200,
	}

	// create a role tag against the image entry
	val, err := createRoleTagValue(rTag1, imageEntry)
	if err != nil {
		t.Fatal(err)
	}
	if val == "" {
		t.Fatalf("failed to create role tag")
	}

	// parse the created role tag
	rTag2, err := parseAndVerifyRoleTagValue(storage, val)
	if err != nil {
		t.Fatal(err)
	}

	// check the values in parsed role tag
	if rTag2.Version != "v1" ||
		rTag2.Nonce != nonce ||
		rTag2.AmiID != "abcd-123" ||
		rTag2.MaxTTL != 200 ||
		!policyutil.EquivalentPolicies(rTag2.Policies, []string{"p", "q", "r"}) ||
		len(rTag2.HMAC) == 0 {
		t.Fatalf("parsed role tag is invalid")
	}

	// verify the tag contents using image specific HMAC key
	verified, err := verifyRoleTagValue(rTag2, imageEntry)
	if err != nil {
		t.Fatal(err)
	}
	if !verified {
		t.Fatalf("failed to verify the role tag")
	}

	// register a different ami
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "image/ami-6789",
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}

	// entry for the newly created ami entry
	imageEntry2, err := awsImage(storage, "ami-6789")
	if err != nil {
		t.Fatal(err)
	}

	// try to verify the tag created with previous image's HMAC key
	// with the newly registered entry's HMAC key
	verified, err = verifyRoleTagValue(rTag2, imageEntry2)
	if err != nil {
		t.Fatal(err)
	}
	if verified {
		t.Fatalf("verification of role tag should have failed: invalid AMI ID")
	}

	// modify any value in role tag and try to verify it
	rTag2.Version = "v2"
	verified, err = verifyRoleTagValue(rTag2, imageEntry)
	if err != nil {
		t.Fatal(err)
	}
	if verified {
		t.Fatalf("verification of role tag should have failed: invalid Version")
	}
}

func TestBackend_prepareRoleTagPlaintextValue(t *testing.T) {
	// create a nonce for the role tag
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
	// try to create plaintext part of role tag
	// without specifying version
	val, err := prepareRoleTagPlaintextValue(rTag)
	if err == nil {
		t.Fatalf("expected error for missing version")
	}
	rTag.Version = "v1"

	rTag.Nonce = ""
	// try to create plaintext part of role tag
	// without specifying nonce
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err == nil {
		t.Fatalf("expected error for missing nonce")
	}
	rTag.Nonce = nonce

	rTag.AmiID = ""
	// try to create plaintext part of role tag
	// without specifying ami_id
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err == nil {
		t.Fatalf("expected error for missing ami_id")
	}
	rTag.AmiID = "abcd-123"

	// create the plaintext part of the tag
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		t.Fatal(err)
	}

	// verify if it contains known fields
	if !strings.Contains(val, "a=") ||
		!strings.Contains(val, "p=") ||
		!strings.Contains(val, "d=") ||
		!strings.HasPrefix(val, "v1") {
		t.Fatalf("incorrect information in role tag plaintext value")
	}

	rTag.InstanceID = "instance-123"
	// create the role tag with instance_id specified
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		t.Fatal(err)
	}
	// verify it
	if !strings.Contains(val, "i=") {
		t.Fatalf("missing instance ID in role tag plaintext value")
	}

	rTag.MaxTTL = 200
	// create the role tag with max_ttl specified
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		t.Fatal(err)
	}
	// verify it
	if !strings.Contains(val, "t=") {
		t.Fatalf("missing instance ID in role tag plaintext value")
	}
}

func TestBackend_CreateRoleTagNonce(t *testing.T) {
	// create a nonce for the role tag
	nonce, err := createRoleTagNonce()
	if err != nil {
		t.Fatal(err)
	}
	if nonce == "" {
		t.Fatalf("failed to create role tag nonce")
	}

	// verify that the value returned is base64 encoded
	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		t.Fatal(err)
	}
	if len(nonceBytes) == 0 {
		t.Fatalf("length of role tag nonce is zero")
	}
}

func TestBackend_ConfigTidyIdentities(t *testing.T) {
	// create a backend
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// test update operation
	tidyRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/tidy/identities",
		Storage:   storage,
	}
	data := map[string]interface{}{
		"safety_buffer":         "60",
		"disable_periodic_tidy": true,
	}
	tidyRequest.Data = data
	_, err = b.HandleRequest(tidyRequest)
	if err != nil {
		t.Fatal(err)
	}

	// test read operation
	tidyRequest.Operation = logical.ReadOperation
	resp, err := b.HandleRequest(tidyRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("failed to read config/tidy/identities endpoint")
	}
	if resp.Data["safety_buffer"].(int) != 60 || !resp.Data["disable_periodic_tidy"].(bool) {
		t.Fatalf("bad: expected: safety_buffer:60 disable_periodic_tidy:true actual: safety_buffer:%s disable_periodic_tidy:%t\n", resp.Data["safety_buffer"].(int), resp.Data["disable_periodic_tidy"].(bool))
	}

	// test delete operation
	tidyRequest.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(tidyRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("failed to delete config/tidy/identities")
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

	// test update operation
	tidyRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/tidy/roletags",
		Storage:   storage,
	}
	data := map[string]interface{}{
		"safety_buffer":         "60",
		"disable_periodic_tidy": true,
	}
	tidyRequest.Data = data
	_, err = b.HandleRequest(tidyRequest)
	if err != nil {
		t.Fatal(err)
	}

	// test read operation
	tidyRequest.Operation = logical.ReadOperation
	resp, err := b.HandleRequest(tidyRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("failed to read config/tidy/roletags endpoint")
	}
	if resp.Data["safety_buffer"].(int) != 60 || !resp.Data["disable_periodic_tidy"].(bool) {
		t.Fatalf("bad: expected: safety_buffer:60 disable_periodic_tidy:true actual: safety_buffer:%s disable_periodic_tidy:%t\n", resp.Data["safety_buffer"].(int), resp.Data["disable_periodic_tidy"].(bool))
	}

	// test delete operation
	tidyRequest.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(tidyRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("failed to delete config/tidy/roletags")
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

	// test update operation
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

	// test update operation
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

	// test existence check returning false
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

	// create an entry
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

	//test existence check returning true
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

	endpointData := map[string]interface{}{
		"secret_key": "secretkey",
		"access_key": "accesskey",
		"endpoint":   "endpointvalue",
	}

	endpointReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Storage:   storage,
		Data:      endpointData,
	}
	_, err = b.HandleRequest(endpointReq)
	if err != nil {
		t.Fatal(err)
	}

	endpointReq.Operation = logical.ReadOperation
	resp, err := b.HandleRequest(endpointReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil ||
		resp.IsError() {
		t.Fatalf("")
	}
	actual := resp.Data["endpoint"].(string)
	if actual != "endpointvalue" {
		t.Fatalf("bad: endpoint: expected:endpointvalue actual:%s\n", actual)
	}
}

func TestBackend_pathConfigCertificate(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	certReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   storage,
		Path:      "config/certificate/cert1",
	}
	checkFound, exists, err := b.HandleExistenceCheck(certReq)
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

	certReq.Data = data
	// test create operation
	_, err = b.HandleRequest(certReq)
	if err != nil {
		t.Fatal(err)
	}

	certReq.Data = nil
	// test existence check
	checkFound, exists, err = b.HandleExistenceCheck(certReq)
	if err != nil {
		t.Fatal(err)
	}
	if !checkFound {
		t.Fatal("existence check not found for path 'config/certificate/cert1'")
	}
	if !exists {
		t.Fatal("existence check should have returned 'true' for 'config/certificate/cert1'")
	}

	certReq.Operation = logical.ReadOperation
	// test read operation
	resp, err := b.HandleRequest(certReq)
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

	certReq.Operation = logical.CreateOperation
	certReq.Path = "config/certificate/cert2"
	certReq.Data = data
	// create another entry to test the list operation
	_, err = b.HandleRequest(certReq)
	if err != nil {
		t.Fatal(err)
	}

	certReq.Operation = logical.ListOperation
	certReq.Path = "config/certificates"
	// test list operation
	resp, err = b.HandleRequest(certReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("failed to list config/certificates")
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != 2 {
		t.Fatalf("invalid keys listed: %#v\n", keys)
	}

	certReq.Operation = logical.DeleteOperation
	certReq.Path = "config/certificate/cert1"
	_, err = b.HandleRequest(certReq)
	if err != nil {
		t.Fatal(err)
	}

	certReq.Path = "config/certificate/cert2"
	_, err = b.HandleRequest(certReq)
	if err != nil {
		t.Fatal(err)
	}

	certReq.Operation = logical.ListOperation
	certReq.Path = "config/certificates"
	// test list operation
	resp, err = b.HandleRequest(certReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("failed to list config/certificates")
	}
	if resp.Data["keys"] != nil {
		t.Fatalf("no entries should be present")
	}
}

func TestBackend_pathImage(t *testing.T) {
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
	data["disallow_reauthentication"] = true
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
	if !resp.Data["allow_instance_migration"].(bool) || !resp.Data["disallow_reauthentication"].(bool) {
		t.Fatal("bad: expected:true got:false\n")
	}

	// add another entry, to test listing of image entries
	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "image/ami-abcd456",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "images",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil || resp.IsError() {
		t.Fatalf("failed to list the image entries")
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != 2 {
		t.Fatalf("bad: keys: %#v\n", keys)
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

func TestBackend_parseAndVerifyRoleTagValue(t *testing.T) {
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
	rTag, err := parseAndVerifyRoleTagValue(storage, tagValue)
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
// Requires TEST_AWS_EC2_PKCS7, TEST_AWS_EC2_AMI_ID to be set.
// If the test is not being run on an EC2 instance that has access to credentials using EC2RoleProvider,
// then TEST_AWS_SECRET_KEY and TEST_AWS_ACCESS_KEY env vars are also required.
func TestBackendAcc_LoginAndWhitelistIdentity(t *testing.T) {
	// This test case should be run only when certain env vars are set and
	// executed as an acceptance test.
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	pkcs7 := os.Getenv("TEST_AWS_EC2_PKCS7")
	if pkcs7 == "" {
		t.Fatalf("env var TEST_AWS_EC2_PKCS7 not set")
	}

	amiID := os.Getenv("TEST_AWS_EC2_AMI_ID")
	if amiID == "" {
		t.Fatalf("env var TEST_AWS_EC2_AMI_ID not set")
	}

	// create the backend
	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	accessKey := os.Getenv("TEST_AWS_ACCESS_KEY")
	secretKey := os.Getenv("TEST_AWS_SECRET_KEY")

	// In case of problems with making API calls using the credentials (2FA enabled,
	// for instance), the keys need not be set if the test is running on an EC2
	// instance with permissions to get the credentials using EC2RoleProvider.
	if accessKey != "" && secretKey != "" {
		// get the API credentials from env vars
		clientConfig := map[string]interface{}{
			"access_key": accessKey,
			"secret_key": secretKey,
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
	}

	// create an entry for the AMI. This is required for login to work.
	data := map[string]interface{}{
		"policies": "root",
		"max_ttl":  "120s",
	}

	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "image/" + amiID,
		Storage:   storage,
		Data:      data,
	})
	if err != nil {
		t.Fatal(err)
	}

	loginInput := map[string]interface{}{
		"pkcs7": pkcs7,
		"nonce": "vault-client-nonce",
	}

	// perform the login operation.
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
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("first login attempt failed")
	}

	// Attempt to login again and see if it succeeds
	resp, err = b.HandleRequest(loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("second login attempt failed")
	}

	// verify the presence of instance_id in the response object.
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

	// Check if a whitelist identity entry is created after the login.
	wlRequest := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "whitelist/identity/" + instanceID,
		Storage:   storage,
	}
	resp, err = b.HandleRequest(wlRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil || resp.Data["ami_id"] != amiID {
		t.Fatalf("failed to read whitelist identity")
	}

	// Delete the whitelist identity entry.
	wlRequest.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(wlRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("failed to delete whitelist identity")
	}

	// Allow a fresh login.
	resp, err = b.HandleRequest(loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("login attempt failed")
	}
}
