package awsauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

const testVaultHeaderValue = "VaultAcceptanceTesting"
const testValidRoleName = "valid-role"
const testInvalidRoleName = "invalid-role"

func TestBackend_CreateParseVerifyRoleTag(t *testing.T) {
	// create a backend
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// create a role entry
	data := map[string]interface{}{
		"auth_type":    "ec2",
		"policies":     "p,q,r,s",
		"bound_ami_id": "abcd-123",
	}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/abcd-123",
		Storage:   storage,
		Data:      data,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role")
	}
	if err != nil {
		t.Fatal(err)
	}

	// read the created role entry
	roleEntry, err := b.lockedAWSRole(context.Background(), storage, "abcd-123")
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
		Role:     "abcd-123",
		Nonce:    nonce,
		Policies: []string{"p", "q", "r"},
		MaxTTL:   200000000000, // 200s
	}

	// create a role tag against the role entry
	val, err := createRoleTagValue(rTag1, roleEntry)
	if err != nil {
		t.Fatal(err)
	}
	if val == "" {
		t.Fatalf("failed to create role tag")
	}

	// parse the created role tag
	rTag2, err := b.parseAndVerifyRoleTagValue(context.Background(), storage, val)
	if err != nil {
		t.Fatal(err)
	}

	// check the values in parsed role tag
	if rTag2.Version != "v1" ||
		rTag2.Nonce != nonce ||
		rTag2.Role != "abcd-123" ||
		rTag2.MaxTTL != 200000000000 || // 200s
		!policyutil.EquivalentPolicies(rTag2.Policies, []string{"p", "q", "r"}) ||
		len(rTag2.HMAC) == 0 {
		t.Fatalf("parsed role tag is invalid")
	}

	// verify the tag contents using role specific HMAC key
	verified, err := verifyRoleTagValue(rTag2, roleEntry)
	if err != nil {
		t.Fatal(err)
	}
	if !verified {
		t.Fatalf("failed to verify the role tag")
	}

	// register a different role
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/ami-6789",
		Storage:   storage,
		Data:      data,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role")
	}
	if err != nil {
		t.Fatal(err)
	}

	// get the entry of the newly created role entry
	roleEntry2, err := b.lockedAWSRole(context.Background(), storage, "ami-6789")
	if err != nil {
		t.Fatal(err)
	}

	// try to verify the tag created with previous role's HMAC key
	// with the newly registered entry's HMAC key
	verified, err = verifyRoleTagValue(rTag2, roleEntry2)
	if err != nil {
		t.Fatal(err)
	}
	if verified {
		t.Fatalf("verification of role tag should have failed")
	}

	// modify any value in role tag and try to verify it
	rTag2.Version = "v2"
	verified, err = verifyRoleTagValue(rTag2, roleEntry)
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
		Role:    "abcd-123",
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

	rTag.Role = ""
	// try to create plaintext part of role tag
	// without specifying role
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err == nil {
		t.Fatalf("expected error for missing role")
	}
	rTag.Role = "abcd-123"

	// create the plaintext part of the tag
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		t.Fatal(err)
	}

	// verify if it contains known fields
	if !strings.Contains(val, "r=") ||
		!strings.Contains(val, "d=") ||
		!strings.Contains(val, "m=") ||
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

	rTag.MaxTTL = 200000000000
	// create the role tag with max_ttl specified
	val, err = prepareRoleTagPlaintextValue(rTag)
	if err != nil {
		t.Fatal(err)
	}
	// verify it
	if !strings.Contains(val, "t=") {
		t.Fatalf("missing max_ttl field in role tag plaintext value")
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

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// test update operation
	tidyRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/tidy/identity-whitelist",
		Storage:   storage,
	}
	data := map[string]interface{}{
		"safety_buffer":         "60",
		"disable_periodic_tidy": true,
	}
	tidyRequest.Data = data
	_, err = b.HandleRequest(context.Background(), tidyRequest)
	if err != nil {
		t.Fatal(err)
	}

	// test read operation
	tidyRequest.Operation = logical.ReadOperation
	resp, err := b.HandleRequest(context.Background(), tidyRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("failed to read config/tidy/identity-whitelist endpoint")
	}
	if resp.Data["safety_buffer"].(int) != 60 || !resp.Data["disable_periodic_tidy"].(bool) {
		t.Fatalf("bad: expected: safety_buffer:60 disable_periodic_tidy:true actual: safety_buffer:%d disable_periodic_tidy:%t\n", resp.Data["safety_buffer"].(int), resp.Data["disable_periodic_tidy"].(bool))
	}

	// test delete operation
	tidyRequest.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), tidyRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("failed to delete config/tidy/identity-whitelist")
	}
}

func TestBackend_ConfigTidyRoleTags(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// test update operation
	tidyRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/tidy/roletag-blacklist",
		Storage:   storage,
	}
	data := map[string]interface{}{
		"safety_buffer":         "60",
		"disable_periodic_tidy": true,
	}
	tidyRequest.Data = data
	_, err = b.HandleRequest(context.Background(), tidyRequest)
	if err != nil {
		t.Fatal(err)
	}

	// test read operation
	tidyRequest.Operation = logical.ReadOperation
	resp, err := b.HandleRequest(context.Background(), tidyRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("failed to read config/tidy/roletag-blacklist endpoint")
	}
	if resp.Data["safety_buffer"].(int) != 60 || !resp.Data["disable_periodic_tidy"].(bool) {
		t.Fatalf("bad: expected: safety_buffer:60 disable_periodic_tidy:true actual: safety_buffer:%d disable_periodic_tidy:%t\n", resp.Data["safety_buffer"].(int), resp.Data["disable_periodic_tidy"].(bool))
	}

	// test delete operation
	tidyRequest.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), tidyRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("failed to delete config/tidy/roletag-blacklist")
	}
}

func TestBackend_TidyIdentities(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	expiredIdentityWhitelist := &whitelistIdentity{
		ExpirationTime: time.Now().Add(-1 * 24 * 365 * time.Hour),
	}
	entry, err := logical.StorageEntryJSON("whitelist/identity/id1", expiredIdentityWhitelist)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	// test update operation
	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "tidy/identity-whitelist",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// let tidy finish in the background
	time.Sleep(1 * time.Second)

	entry, err = storage.Get(context.Background(), "whitelist/identity/id1")
	if err != nil {
		t.Fatal(err)
	}
	if entry != nil {
		t.Fatal("wl tidy did not remove expired entry")
	}
}

func TestBackend_TidyRoleTags(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	expiredIdentityWhitelist := &roleTagBlacklistEntry{
		ExpirationTime: time.Now().Add(-1 * 24 * 365 * time.Hour),
	}
	entry, err := logical.StorageEntryJSON("blacklist/roletag/id1", expiredIdentityWhitelist)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	// test update operation
	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "tidy/roletag-blacklist",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// let tidy finish in the background
	time.Sleep(1 * time.Second)

	entry, err = storage.Get(context.Background(), "blacklist/roletag/id1")
	if err != nil {
		t.Fatal(err)
	}
	if entry != nil {
		t.Fatal("bl tidy did not remove expired entry")
	}
}

func TestBackend_ConfigClient(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
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
		AcceptanceTest:    false,
		CredentialBackend: b,
		Steps: []logicaltest.TestStep{
			stepCreate,
			stepInvalidAccessKey,
			stepInvalidSecretKey,
			stepUpdate,
		},
	})

	// test existence check returning false
	checkFound, exists, err := b.HandleExistenceCheck(context.Background(), &logical.Request{
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
	_, err = b.HandleRequest(context.Background(), configClientCreateRequest)
	if err != nil {
		t.Fatal(err)
	}

	//test existence check returning true
	checkFound, exists, err = b.HandleExistenceCheck(context.Background(), &logical.Request{
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
	_, err = b.HandleRequest(context.Background(), endpointReq)
	if err != nil {
		t.Fatal(err)
	}

	endpointReq.Operation = logical.ReadOperation
	resp, err := b.HandleRequest(context.Background(), endpointReq)
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

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	certReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   storage,
		Path:      "config/certificate/cert1",
	}
	checkFound, exists, err := b.HandleExistenceCheck(context.Background(), certReq)
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
		"type": "pkcs7",
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
	resp, err := b.HandleRequest(context.Background(), certReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	certReq.Data = nil
	// test existence check
	checkFound, exists, err = b.HandleExistenceCheck(context.Background(), certReq)
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
	resp, err = b.HandleRequest(context.Background(), certReq)
	if err != nil {
		t.Fatal(err)
	}
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
		t.Fatalf("bad: expected:%s\n got:%s\n", expectedCert, resp.Data["aws_public_cert"].(string))
	}

	certReq.Operation = logical.CreateOperation
	certReq.Path = "config/certificate/cert2"
	certReq.Data = data
	// create another entry to test the list operation
	_, err = b.HandleRequest(context.Background(), certReq)
	if err != nil {
		t.Fatal(err)
	}

	certReq.Operation = logical.ListOperation
	certReq.Path = "config/certificates"
	// test list operation
	resp, err = b.HandleRequest(context.Background(), certReq)
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
	_, err = b.HandleRequest(context.Background(), certReq)
	if err != nil {
		t.Fatal(err)
	}

	certReq.Path = "config/certificate/cert2"
	_, err = b.HandleRequest(context.Background(), certReq)
	if err != nil {
		t.Fatal(err)
	}

	certReq.Operation = logical.ListOperation
	certReq.Path = "config/certificates"
	// test list operation
	resp, err = b.HandleRequest(context.Background(), certReq)
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

func TestBackend_parseAndVerifyRoleTagValue(t *testing.T) {
	// create a backend
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// create a role
	data := map[string]interface{}{
		"auth_type":    "ec2",
		"policies":     "p,q,r,s",
		"max_ttl":      "120s",
		"role_tag":     "VaultRole",
		"bound_ami_id": "abcd-123",
	}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/abcd-123",
		Storage:   storage,
		Data:      data,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role")
	}
	if err != nil {
		t.Fatal(err)
	}

	// verify that the entry is created
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/abcd-123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("expected an role entry for abcd-123")
	}

	// create a role tag
	data2 := map[string]interface{}{
		"policies": "p,q,r,s",
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/abcd-123/tag",
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
	rTag, err := b.parseAndVerifyRoleTagValue(context.Background(), storage, tagValue)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if rTag == nil {
		t.Fatalf("failed to parse role tag")
	}
	if rTag.Version != "v1" ||
		!policyutil.EquivalentPolicies(rTag.Policies, []string{"p", "q", "r", "s"}) ||
		rTag.Role != "abcd-123" {
		t.Fatalf("bad: parsed role tag contains incorrect values. Got: %#v\n", rTag)
	}
}

func TestBackend_PathRoleTag(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"auth_type":    "ec2",
		"policies":     "p,q,r,s",
		"max_ttl":      "120s",
		"role_tag":     "VaultRole",
		"bound_ami_id": "abcd-123",
	}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/abcd-123",
		Storage:   storage,
		Data:      data,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role")
	}
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/abcd-123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("failed to find a role entry for abcd-123")
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/abcd-123/tag",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("failed to create a tag on role: abcd-123")
	}
	if resp.IsError() {
		t.Fatalf("failed to create a tag on role: abcd-123: %s\n", resp.Data["error"])
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
	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// create an role entry
	data := map[string]interface{}{
		"auth_type":    "ec2",
		"policies":     "p,q,r,s",
		"role_tag":     "VaultRole",
		"bound_ami_id": "abcd-123",
	}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/abcd-123",
		Storage:   storage,
		Data:      data,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role")
	}
	if err != nil {
		t.Fatal(err)
	}

	// create a role tag against an role registered before
	data2 := map[string]interface{}{
		"policies": "p,q,r,s",
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/abcd-123/tag",
		Storage:   storage,
		Data:      data2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("failed to create a tag on role: abcd-123")
	}
	if resp.IsError() {
		t.Fatalf("failed to create a tag on role: abcd-123: %s\n", resp.Data["error"])
	}
	tag := resp.Data["tag_value"].(string)
	if tag == "" {
		t.Fatalf("role tag not present in the response data: %#v\n", resp.Data)
	}

	// blacklist that role tag
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roletag-blacklist/" + tag,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("failed to blacklist the roletag: %s\n", tag)
	}

	// read the blacklist entry
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "roletag-blacklist/" + tag,
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
	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "roletag-blacklist/" + tag,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// try to read the deleted entry
	tagEntry, err := b.lockedBlacklistRoleTagEntry(context.Background(), storage, tag)
	if err != nil {
		t.Fatal(err)
	}
	if tagEntry != nil {
		t.Fatalf("role tag should not have been present: %s\n", tag)
	}
}

/* This is an acceptance test.
   Requires the following env vars:
   TEST_AWS_EC2_PKCS7
   TEST_AWS_EC2_IDENTITY_DOCUMENT
   TEST_AWS_EC2_IDENTITY_DOCUMENT_SIG
   TEST_AWS_EC2_AMI_ID
   TEST_AWS_EC2_ACCOUNT_ID
   TEST_AWS_EC2_IAM_ROLE_ARN

   If this is being run on an EC2 instance, you can set the environment vars using this bash snippet:

   export TEST_AWS_EC2_PKCS7=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/pkcs7)
   export TEST_AWS_EC2_IDENTITY_DOCUMENT=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document | base64 -w 0)
   export TEST_AWS_EC2_IDENTITY_DOCUMENT_SIG=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/signature | tr -d '\n')
   export TEST_AWS_EC2_AMI_ID=$(curl -s http://169.254.169.254/latest/meta-data/ami-id)
   export TEST_AWS_EC2_IAM_ROLE_ARN=$(aws iam get-role --role-name $(curl -q http://169.254.169.254/latest/meta-data/iam/security-credentials/ -S -s) --query Role.Arn --output text)
   export TEST_AWS_EC2_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

   If the test is not being run on an EC2 instance that has access to
   credentials using EC2RoleProvider, on top of the above vars, following
   needs to be set:
   TEST_AWS_SECRET_KEY
   TEST_AWS_ACCESS_KEY
*/
func TestBackendAcc_LoginWithInstanceIdentityDocAndWhitelistIdentity(t *testing.T) {
	// This test case should be run only when certain env vars are set and
	// executed as an acceptance test.
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	pkcs7 := os.Getenv("TEST_AWS_EC2_PKCS7")
	if pkcs7 == "" {
		t.Skipf("env var TEST_AWS_EC2_PKCS7 not set, skipping test")
	}

	identityDoc := os.Getenv("TEST_AWS_EC2_IDENTITY_DOCUMENT")
	if identityDoc == "" {
		t.Skipf("env var TEST_AWS_EC2_IDENTITY_DOCUMENT not set, skipping test")
	}

	identityDocSig := os.Getenv("TEST_AWS_EC2_IDENTITY_DOCUMENT_SIG")
	if identityDocSig == "" {
		t.Skipf("env var TEST_AWS_EC2_IDENTITY_DOCUMENT_SIG not set, skipping test")
	}

	amiID := os.Getenv("TEST_AWS_EC2_AMI_ID")
	if amiID == "" {
		t.Skipf("env var TEST_AWS_EC2_AMI_ID not set, skipping test")
	}

	iamARN := os.Getenv("TEST_AWS_EC2_IAM_ROLE_ARN")
	if iamARN == "" {
		t.Skipf("env var TEST_AWS_EC2_IAM_ROLE_ARN not set, skipping test")
	}

	accountID := os.Getenv("TEST_AWS_EC2_ACCOUNT_ID")
	if accountID == "" {
		t.Skipf("env var TEST_AWS_EC2_ACCOUNT_ID not set, skipping test")
	}

	roleName := amiID

	// create the backend
	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage
	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
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
		_, err = b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.UpdateOperation,
			Storage:   storage,
			Path:      "config/client",
			Data:      clientConfig,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	loginInput := map[string]interface{}{
		"pkcs7": pkcs7,
		"nonce": "vault-client-nonce",
	}

	parsedIdentityDoc, err := b.parseIdentityDocument(context.Background(), storage, pkcs7)
	if err != nil {
		t.Fatal(err)
	}

	// Perform the login operation with a AMI ID that is not matching
	// the bound on the role.
	loginRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginInput,
	}

	// Baseline role data that should succeed permit login
	data := map[string]interface{}{
		"auth_type":             "ec2",
		"policies":              "root",
		"max_ttl":               "120s",
		"bound_ami_id":          []string{"wrong_ami_id", amiID, "wrong_ami_id2"},
		"bound_account_id":      accountID,
		"bound_iam_role_arn":    iamARN,
		"bound_ec2_instance_id": []string{parsedIdentityDoc.InstanceID, "i-1234567"},
	}

	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/" + roleName,
		Storage:   storage,
		Data:      data,
	}

	updateRoleExpectLoginFail := func(roleRequest, loginRequest *logical.Request) error {
		resp, err := b.HandleRequest(context.Background(), roleRequest)
		if err != nil || (resp != nil && resp.IsError()) {
			return fmt.Errorf("bad: failed to create role: resp:%#v\nerr:%v", resp, err)
		}
		resp, err = b.HandleRequest(context.Background(), loginRequest)
		if err != nil || resp == nil || (resp != nil && !resp.IsError()) {
			return fmt.Errorf("bad: expected login failure: resp:%#v\nerr:%v", resp, err)
		}
		return nil
	}

	// Test a role with the wrong AMI ID
	data["bound_ami_id"] = []string{"ami-1234567", "ami-7654321"}
	if err := updateRoleExpectLoginFail(roleReq, loginRequest); err != nil {
		t.Fatal(err)
	}

	roleReq.Operation = logical.UpdateOperation
	// Place the correct AMI ID in one of the values, but make the AccountID wrong
	data["bound_ami_id"] = []string{"wrong_ami_id_1", amiID, "wrong_ami_id_2"}
	data["bound_account_id"] = []string{"wrong-account-id", "wrong-account-id-2"}
	if err := updateRoleExpectLoginFail(roleReq, loginRequest); err != nil {
		t.Fatal(err)
	}

	// Place the correct AccountID in one of the values, but make the wrong IAMRoleARN
	data["bound_account_id"] = []string{"wrong-account-id-1", accountID, "wrong-account-id-2"}
	data["bound_iam_role_arn"] = []string{"wrong_iam_role_arn", "wrong_iam_role_arn_2"}
	if err := updateRoleExpectLoginFail(roleReq, loginRequest); err != nil {
		t.Fatal(err)
	}

	// Place correct IAM role ARN, but incorrect instance ID
	data["bound_iam_role_arn"] = []string{"wrong_iam_role_arn_1", iamARN, "wrong_iam_role_arn_2"}
	data["bound_ec2_instance_id"] = "i-1234567"
	if err := updateRoleExpectLoginFail(roleReq, loginRequest); err != nil {
		t.Fatal(err)
	}

	// Place correct instance ID, but substring of the IAM role ARN
	data["bound_ec2_instance_id"] = []string{parsedIdentityDoc.InstanceID, "i-1234567"}
	data["bound_iam_role_arn"] = []string{"wrong_iam_role_arn", iamARN[:len(iamARN)-2], "wrong_iam_role_arn_2"}
	if err := updateRoleExpectLoginFail(roleReq, loginRequest); err != nil {
		t.Fatal(err)
	}

	// place a wildcard in the middle of the role ARN
	// The :31 gets arn:aws:iam::123456789012:role/
	// This test relies on the role name having at least two characters
	data["bound_iam_role_arn"] = []string{"wrong_iam_role_arn", fmt.Sprintf("%s*%s", iamARN[:31], iamARN[32:])}
	if err := updateRoleExpectLoginFail(roleReq, loginRequest); err != nil {
		t.Fatal(err)
	}

	// globbed IAM role ARN
	data["bound_iam_role_arn"] = []string{"wrong_iam_role_arn_1", fmt.Sprintf("%s*", iamARN[:len(iamARN)-2]), "wrong_iam_role_arn_2"}
	resp, err := b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: failed to create role: resp:%#v\nerr:%v", resp, err)
	}

	// Now, the login attempt should succeed
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("bad: failed to login: resp:%#v\nerr:%v", resp, err)
	}

	// Attempt to re-login with the identity signature
	delete(loginInput, "pkcs7")
	loginInput["identity"] = identityDoc
	loginInput["signature"] = identityDocSig
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("bad: failed to login: resp:%#v\nerr:%v", resp, err)
	}

	// verify the presence of instance_id in the response object.
	instanceID := resp.Auth.Metadata["instance_id"]
	if instanceID == "" {
		t.Fatalf("instance ID not present in the response object")
	}
	if instanceID != parsedIdentityDoc.InstanceID {
		t.Fatalf("instance ID in response (%q) did not match instance ID from identity document (%q)", instanceID, parsedIdentityDoc.InstanceID)
	}

	_, ok := resp.Auth.Metadata["nonce"]
	if ok {
		t.Fatalf("client nonce should not have been returned")
	}

	loginInput["nonce"] = "changed-vault-client-nonce"
	// try to login again with changed nonce
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("login attempt should have failed due to client nonce mismatch")
	}

	// Check if a whitelist identity entry is created after the login.
	wlRequest := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "identity-whitelist/" + instanceID,
		Storage:   storage,
	}
	resp, err = b.HandleRequest(context.Background(), wlRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil || resp.Data["role"] != roleName {
		t.Fatalf("failed to read whitelist identity")
	}

	// Delete the whitelist identity entry.
	wlRequest.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), wlRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("failed to delete whitelist identity")
	}

	// Allow a fresh login without supplying the nonce
	delete(loginInput, "nonce")

	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("login attempt failed")
	}

	_, ok = resp.Auth.Metadata["nonce"]
	if !ok {
		t.Fatalf("expected nonce to be returned")
	}
}

func TestBackend_pathStsConfig(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	stsReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   storage,
		Path:      "config/sts/account1",
	}
	checkFound, exists, err := b.HandleExistenceCheck(context.Background(), stsReq)
	if err != nil {
		t.Fatal(err)
	}
	if !checkFound {
		t.Fatal("existence check not found for path 'config/sts/account1'")
	}
	if exists {
		t.Fatal("existence check should have returned 'false' for 'config/sts/account1'")
	}

	data := map[string]interface{}{
		"sts_role": "arn:aws:iam:account1:role/myRole",
	}

	stsReq.Data = data
	// test create operation
	resp, err := b.HandleRequest(context.Background(), stsReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	stsReq.Data = nil
	// test existence check
	checkFound, exists, err = b.HandleExistenceCheck(context.Background(), stsReq)
	if err != nil {
		t.Fatal(err)
	}
	if !checkFound {
		t.Fatal("existence check not found for path 'config/sts/account1'")
	}
	if !exists {
		t.Fatal("existence check should have returned 'true' for 'config/sts/account1'")
	}

	stsReq.Operation = logical.ReadOperation
	// test read operation
	resp, err = b.HandleRequest(context.Background(), stsReq)
	if err != nil {
		t.Fatal(err)
	}
	expectedStsRole := "arn:aws:iam:account1:role/myRole"
	if resp.Data["sts_role"].(string) != expectedStsRole {
		t.Fatalf("bad: expected:%s\n got:%s\n", expectedStsRole, resp.Data["sts_role"].(string))
	}

	stsReq.Operation = logical.CreateOperation
	stsReq.Path = "config/sts/account2"
	stsReq.Data = data
	// create another entry to test the list operation
	resp, err = b.HandleRequest(context.Background(), stsReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatal(err)
	}

	stsReq.Operation = logical.ListOperation
	stsReq.Path = "config/sts"
	// test list operation
	resp, err = b.HandleRequest(context.Background(), stsReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("failed to list config/sts")
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != 2 {
		t.Fatalf("invalid keys listed: %#v\n", keys)
	}

	stsReq.Operation = logical.DeleteOperation
	stsReq.Path = "config/sts/account1"
	resp, err = b.HandleRequest(context.Background(), stsReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatal(err)
	}

	stsReq.Path = "config/sts/account2"
	resp, err = b.HandleRequest(context.Background(), stsReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatal(err)
	}

	stsReq.Operation = logical.ListOperation
	stsReq.Path = "config/sts"
	// test list operation
	resp, err = b.HandleRequest(context.Background(), stsReq)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("failed to list config/sts")
	}
	if resp.Data["keys"] != nil {
		t.Fatalf("no entries should be present")
	}
}

func buildCallerIdentityLoginData(request *http.Request, roleName string) (map[string]interface{}, error) {
	headersJson, err := json.Marshal(request.Header)
	if err != nil {
		return nil, err
	}
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"iam_http_request_method": request.Method,
		"iam_request_url":         base64.StdEncoding.EncodeToString([]byte(request.URL.String())),
		"iam_request_headers":     base64.StdEncoding.EncodeToString(headersJson),
		"iam_request_body":        base64.StdEncoding.EncodeToString(requestBody),
		"request_role":            roleName,
	}, nil
}

// This is an acceptance test.
// If the test is NOT being run on an AWS EC2 instance in an instance profile,
// it requires the following environment variables to be set:
// TEST_AWS_ACCESS_KEY_ID
// TEST_AWS_SECRET_ACCESS_KEY
// TEST_AWS_SECURITY_TOKEN or TEST_AWS_SESSION_TOKEN (optional, if you are using short-lived creds)
// These are intentionally NOT the "standard" variables to prevent accidentally
// using prod creds in acceptance tests
func TestBackendAcc_LoginWithCallerIdentity(t *testing.T) {
	// This test case should be run only when certain env vars are set and
	// executed as an acceptance test.
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage
	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// Override the default AWS env vars (if set) with our test creds
	// so that the credential provider chain will pick them up
	// NOTE that I'm not bothing to override the shared config file location,
	// so if creds are specified there, they will be used before IAM
	// instance profile creds
	// This doesn't provide perfect leakage protection (e.g., it will still
	// potentially pick up credentials from the ~/.config files), but probably
	// good enough rather than having to muck around in the low-level details
	for _, envvar := range []string{
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SECURITY_TOKEN", "AWS_SESSION_TOKEN"} {
		// Skip test if any of the required env vars are missing
		testEnvVar := os.Getenv("TEST_" + envvar)
		if testEnvVar == "" {
			t.Skipf("env var %s not set, skipping test", "TEST_"+envvar)
		}

		// restore existing environment variables (in case future tests need them)
		defer os.Setenv(envvar, os.Getenv(envvar))

		os.Setenv(envvar, testEnvVar)
	}
	awsSession, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	stsService := sts.New(awsSession)
	stsInputParams := &sts.GetCallerIdentityInput{}

	testIdentity, err := stsService.GetCallerIdentity(stsInputParams)
	if err != nil {
		t.Fatalf("Received error retrieving identity: %s", err)
	}
	entity, err := parseIamArn(*testIdentity.Arn)
	if err != nil {
		t.Fatal(err)
	}

	// Test setup largely done
	// At this point, we're going to:
	// 1. Configure the client to require our test header value
	// 2. Configure identity to use the ARN for the alias
	// 3. Configure two different roles:
	//    a. One bound to our test user
	//    b. One bound to a garbage ARN
	// 4. Pass in a request that doesn't have the signed header, ensure
	//    we're not allowed to login
	// 5. Passin a request that has a validly signed header, but the wrong
	//    value, ensure it doesn't allow login
	// 6. Pass in a request that has a validly signed request, ensure
	//    it allows us to login to our role
	// 7. Pass in a request that has a validly signed request, asking for
	//    the other role, ensure it fails

	clientConfigData := map[string]interface{}{
		"iam_server_id_header_value": testVaultHeaderValue,
	}
	clientRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Storage:   storage,
		Data:      clientConfigData,
	}
	_, err = b.HandleRequest(context.Background(), clientRequest)
	if err != nil {
		t.Fatal(err)
	}

	configIdentityData := map[string]interface{}{
		"iam_alias": identityAliasIAMFullArn,
	}
	configIdentityRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/identity",
		Storage:   storage,
		Data:      configIdentityData,
	}
	resp, err := b.HandleRequest(context.Background(), configIdentityRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("received error response when configuring identity: %#v", resp)
	}

	// configuring the valid role we'll be able to login to
	roleData := map[string]interface{}{
		"bound_iam_principal_arn": []string{entity.canonicalArn(), "arn:aws:iam::123456789012:role/FakeRoleArn1*"}, // Fake ARN MUST be wildcard terminated because we're resolving unique IDs, and the wildcard termination prevents unique ID resolution
		"policies":                "root",
		"auth_type":               iamAuthType,
	}
	roleRequest := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/" + testValidRoleName,
		Storage:   storage,
		Data:      roleData,
	}
	resp, err = b.HandleRequest(context.Background(), roleRequest)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: failed to create role: resp:%#v\nerr:%v", resp, err)
	}

	// configuring a valid role we won't be able to login to
	roleDataEc2 := map[string]interface{}{
		"auth_type":    "ec2",
		"policies":     "root",
		"bound_ami_id": "ami-1234567",
	}
	roleRequestEc2 := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/ec2only",
		Storage:   storage,
		Data:      roleDataEc2,
	}
	resp, err = b.HandleRequest(context.Background(), roleRequestEc2)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: failed to create role; resp:%#v\nerr:%v", resp, err)
	}

	fakeArn := "arn:aws:iam::123456789012:role/somePath/FakeRole"
	fakeArn2 := "arn:aws:iam::123456789012:role/somePath/FakeRole2"
	fakeArnResolverCount := 0
	fakeArnResolver := func(ctx context.Context, s logical.Storage, arn string) (string, error) {
		if strings.HasPrefix(arn, fakeArn) {
			fakeArnResolverCount++
			return fmt.Sprintf("FakeUniqueIdFor%s%d", arn, fakeArnResolverCount), nil
		}
		return b.resolveArnToRealUniqueId(context.Background(), s, arn)
	}
	b.resolveArnToUniqueIDFunc = fakeArnResolver

	// now we're creating the invalid role we won't be able to login to
	roleData["bound_iam_principal_arn"] = []string{fakeArn, fakeArn2}
	roleRequest.Path = "role/" + testInvalidRoleName
	resp, err = b.HandleRequest(context.Background(), roleRequest)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: didn't fail to create role: resp:%#v\nerr:%v", resp, err)
	}

	// now, create the request without the signed header
	stsRequestNoHeader, _ := stsService.GetCallerIdentityRequest(stsInputParams)
	stsRequestNoHeader.Sign()
	loginData, err := buildCallerIdentityLoginData(stsRequestNoHeader.HTTPRequest, testValidRoleName)
	if err != nil {
		t.Fatal(err)
	}
	loginRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginData,
	}
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil || resp == nil || !resp.IsError() {
		t.Errorf("bad: expected failed login due to missing header: resp:%#v\nerr:%v", resp, err)
	}

	// create the request with the invalid header value

	// Not reusing stsRequestNoHeader because the process of signing the request
	// and reading the body modifies the underlying request, so it's just cleaner
	// to get new requests.
	stsRequestInvalidHeader, _ := stsService.GetCallerIdentityRequest(stsInputParams)
	stsRequestInvalidHeader.HTTPRequest.Header.Add(iamServerIdHeader, "InvalidValue")
	stsRequestInvalidHeader.Sign()
	loginData, err = buildCallerIdentityLoginData(stsRequestInvalidHeader.HTTPRequest, testValidRoleName)
	if err != nil {
		t.Fatal(err)
	}
	loginRequest = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginData,
	}
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil || resp == nil || !resp.IsError() {
		t.Errorf("bad: expected failed login due to invalid header: resp:%#v\nerr:%v", resp, err)
	}

	// Now, valid request against invalid role
	stsRequestValid, _ := stsService.GetCallerIdentityRequest(stsInputParams)
	stsRequestValid.HTTPRequest.Header.Add(iamServerIdHeader, testVaultHeaderValue)
	stsRequestValid.Sign()
	loginData, err = buildCallerIdentityLoginData(stsRequestValid.HTTPRequest, testInvalidRoleName)
	if err != nil {
		t.Fatal(err)
	}
	loginRequest = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginData,
	}
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil || resp == nil || !resp.IsError() {
		t.Errorf("bad: expected failed login due to invalid role: resp:%#v\nerr:%v", resp, err)
	}

	loginData["role"] = "ec2only"
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil || resp == nil || !resp.IsError() {
		t.Errorf("bad: expected failed login due to bad auth type: resp:%#v\nerr:%v", resp, err)
	}

	// finally, the happy path test :)

	loginData["role"] = testValidRoleName
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("bad: expected valid login: resp:%#v", resp)
	}
	if resp.Auth.Alias == nil {
		t.Fatalf("bad: nil auth Alias")
	}
	if resp.Auth.Alias.Name != *testIdentity.Arn {
		t.Fatalf("bad: expected identity alias of %q, got %q instead", *testIdentity.Arn, resp.Auth.Alias.Name)
	}

	renewReq := generateRenewRequest(storage, resp.Auth)
	// dump a fake ARN into the metadata to ensure that we ONLY look
	// at the unique ID that has been generated
	renewReq.Auth.Metadata["canonical_arn"] = "fake_arn"
	empty_login_fd := &framework.FieldData{
		Raw:    map[string]interface{}{},
		Schema: pathLogin(b).Fields,
	}
	// ensure we can renew
	resp, err = b.pathLoginRenew(context.Background(), renewReq, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}
	if resp.IsError() {
		t.Fatalf("got error when renewing: %#v", *resp)
	}

	// Now, fake out the unique ID resolver to ensure we fail login if the unique ID
	// changes from under us
	b.resolveArnToUniqueIDFunc = resolveArnToFakeUniqueId
	// First, we need to update the role to force Vault to use our fake resolver to
	// pick up the fake user ID
	roleData["bound_iam_principal_arn"] = entity.canonicalArn()
	roleRequest.Path = "role/" + testValidRoleName
	resp, err = b.HandleRequest(context.Background(), roleRequest)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: failed to recreate role: resp:%#v\nerr:%v", resp, err)
	}
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil || resp == nil || !resp.IsError() {
		t.Errorf("bad: expected failed login due to changed AWS role ID: resp: %#v\nerr:%v", resp, err)
	}

	// and ensure a renew no longer works
	resp, err = b.pathLoginRenew(context.Background(), renewReq, empty_login_fd)
	if err == nil || (resp != nil && !resp.IsError()) {
		t.Errorf("bad: expected failed renew due to changed AWS role ID: resp: %#v", resp)
	}
	// Undo the fake resolver...
	b.resolveArnToUniqueIDFunc = b.resolveArnToRealUniqueId

	// Now test that wildcard matching works
	wildcardRoleName := "valid_wildcard"
	wildcardEntity := *entity
	wildcardEntity.FriendlyName = "*"
	roleData["bound_iam_principal_arn"] = []string{wildcardEntity.canonicalArn(), "arn:aws:iam::123456789012:role/DoesNotExist/Vault_Fake_Role*"}
	roleRequest.Path = "role/" + wildcardRoleName
	resp, err = b.HandleRequest(context.Background(), roleRequest)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: failed to create wildcard roles: resp:%#v\nerr:%v", resp, err)
	}

	loginData["role"] = wildcardRoleName
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("bad: expected valid login: resp:%#v", resp)
	}
	// and ensure we can renew
	renewReq = generateRenewRequest(storage, resp.Auth)
	resp, err = b.pathLoginRenew(context.Background(), renewReq, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}
	if resp.IsError() {
		t.Fatalf("got error when renewing: %#v", *resp)
	}
	// ensure the cache is populated
	cachedArn := b.getCachedUserId(resp.Auth.Metadata["client_user_id"])
	if cachedArn == "" {
		t.Errorf("got empty ARN back from user ID cache; expected full arn")
	}

	// Test for renewal with period
	period := 600 * time.Second
	roleData["period"] = period.String()
	roleRequest.Path = "role/" + testValidRoleName
	resp, err = b.HandleRequest(context.Background(), roleRequest)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: failed to create wildcard role: resp:%#v\nerr:%v", resp, err)
	}

	loginData["role"] = testValidRoleName
	resp, err = b.HandleRequest(context.Background(), loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Fatalf("bad: expected valid login: resp:%#v", resp)
	}

	renewReq = generateRenewRequest(storage, resp.Auth)
	resp, err = b.pathLoginRenew(context.Background(), renewReq, empty_login_fd)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}
	if resp.IsError() {
		t.Fatalf("got error when renewing: %#v", *resp)
	}

	if resp.Auth.Period != period {
		t.Fatalf("expected a period value of %s in the response, got: %s", period, resp.Auth.Period)
	}
}

func generateRenewRequest(s logical.Storage, auth *logical.Auth) *logical.Request {
	renewReq := &logical.Request{
		Storage: s,
		Auth:    &logical.Auth{},
	}
	renewReq.Auth.InternalData = auth.InternalData
	renewReq.Auth.Metadata = auth.Metadata
	renewReq.Auth.LeaseOptions = auth.LeaseOptions
	renewReq.Auth.Policies = auth.Policies
	renewReq.Auth.Period = auth.Period

	return renewReq
}
