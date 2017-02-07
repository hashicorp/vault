package awsiam

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

func buildLoginData(request *http.Request, roleName string) (map[string]interface{}, error) {
	headersJson, err := json.Marshal(request.Header)
	if err != nil {
		return nil, err
	}
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"method":  request.Method,
		"url":     request.URL.String(),
		"headers": base64.StdEncoding.EncodeToString(headersJson),
		"body":    base64.StdEncoding.EncodeToString(requestBody),
		"role":    roleName,
	}, nil
}

// This is an acceptance test.
// If the test is NOT being run on an AWS EC2 instance in an instance profile,
// it requires the following environment variables to be set:
// TEST_AWS_ACCESS_KEY_ID
// TEST_AWS_SECRET_ACCESS_KEY
// TEST_AWS_SECURITY_TOKEN (optional, if you are using short-lived creds)
// These are intentionally NOT the "standard" variables to prevent accidentally
// using prod creds in acceptance tests
func TestBackendAcc_Login(t *testing.T) {
	// This test case should be run only when certain env vars are set and
	// executed as an acceptance test.
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	storage := &logical.InmemStorage{}
	config := logical.TestBackendConfig()
	config.StorageView = storage
	b := Backend()
	_, err := b.Setup(config)
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
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SECURITY_TOKEN"} {
		os.Setenv("TEST_"+envvar, os.Getenv(envvar))
	}
	awsSession, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	stsService := sts.New(awsSession)
	var stsInputParams *sts.GetCallerIdentityInput

	testIdentity, err := stsService.GetCallerIdentity(stsInputParams)
	if err != nil {
		t.Fatalf("Received error retrieving identity: %s", err)
	}
	testIdentityArn, _, err := parseIamArn(*testIdentity.Arn)
	if err != nil {
		t.Fatal(err)
	}

	// Test setup largely done
	// At this point, we're going to:
	// 1. Configure the client to require our test header value
	// 2. Configure two different roles:
	//    a. One bound to our test user
	//    b. One bound to a garbage ARN
	// 3. Pass in a request that doesn't have the signed header, ensure
	//    we're not allowed to login
	// 4. Passin a request that has a validly signed header, but the wrong
	//    value, ensure it doesn't allow login
	// 5. Pass in a request that has a validly signed request, ensure
	//    it allows us to login to our role
	// 6. Pass in a request that has a validly signed request, asking for
	//    the other role, ensure it fails
	const testVaultHeaderValue = "VaultAcceptanceTesting"
	const testValidRoleName = "valid-role"
	const testInvalidRoleName = "invalid-role"

	clientConfigData := map[string]interface{}{
		"vault_header_value": testVaultHeaderValue,
	}
	clientRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Storage:   storage,
		Data:      clientConfigData,
	}
	_, err = b.HandleRequest(clientRequest)
	if err != nil {
		t.Fatal(err)
	}

	// configuring the valid role we'll be able to login to
	roleData := map[string]interface{}{
		"bound_iam_principal": testIdentityArn,
		"policies":            "root",
	}
	roleRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/" + testValidRoleName,
		Storage:   storage,
		Data:      roleData,
	}
	resp, err := b.HandleRequest(roleRequest)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: failed to create role: resp:%#v\nerr:%v", resp, err)
	}

	// now we're creating the invalid role we won't be able to login to
	roleData["bound_iam_principal"] = "arn:aws:iam::123456789012:role/FakeRole"
	roleRequest.Path = "role/" + testInvalidRoleName
	resp, err = b.HandleRequest(roleRequest)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: failed to create role: resp:%#v\nerr:%v", resp, err)
	}

	// now, create the request without the signed header
	stsRequestNoHeader, _ := stsService.GetCallerIdentityRequest(stsInputParams)
	stsRequestNoHeader.Sign()
	loginData, err := buildLoginData(stsRequestNoHeader.HTTPRequest, testValidRoleName)
	if err != nil {
		t.Fatal(err)
	}
	loginRequest := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginData,
	}
	resp, err = b.HandleRequest(loginRequest)
	if err != nil || resp == nil || !resp.IsError() {
		t.Errorf("bad: expected failed login due to missing header: resp:%#v\nerr:%v", resp, err)
	}

	// create the request with the invalid header value

	// Not reusing stsRequestNoHeader because the process of signing the request
	// and reading the body modifies the underlying request, so it's just cleaner
	// to get new requests.
	stsRequestInvalidHeader, _ := stsService.GetCallerIdentityRequest(stsInputParams)
	stsRequestInvalidHeader.HTTPRequest.Header.Add(magicVaultHeader, "InvalidValue")
	stsRequestInvalidHeader.Sign()
	loginData, err = buildLoginData(stsRequestInvalidHeader.HTTPRequest, testValidRoleName)
	if err != nil {
		t.Fatal(err)
	}
	loginRequest = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginData,
	}
	resp, err = b.HandleRequest(loginRequest)
	if err != nil || resp == nil || !resp.IsError() {
		t.Errorf("bad: expected failed login due to invalid header: resp:%#v\nerr:%v", resp, err)
	}

	// Now, valid request against invalid role
	stsRequestValid, _ := stsService.GetCallerIdentityRequest(stsInputParams)
	stsRequestValid.HTTPRequest.Header.Add(magicVaultHeader, testVaultHeaderValue)
	stsRequestValid.Sign()
	loginData, err = buildLoginData(stsRequestValid.HTTPRequest, testInvalidRoleName)
	if err != nil {
		t.Fatal(err)
	}
	loginRequest = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Storage:   storage,
		Data:      loginData,
	}
	resp, err = b.HandleRequest(loginRequest)
	if err != nil || resp == nil || !resp.IsError() {
		t.Errorf("bad: expected failed login due to invalid role: resp:%#v\nerr:%v", resp, err)
	}

	loginData["role"] = testValidRoleName
	resp, err = b.HandleRequest(loginRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Auth == nil || resp.IsError() {
		t.Errorf("bad: expected valid login: resp:%#v", resp)
	}
}
