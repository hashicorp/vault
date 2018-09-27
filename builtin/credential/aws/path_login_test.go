package awsauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/vault/logical"
)

func TestBackend_pathLogin_getCallerIdentityResponse(t *testing.T) {
	responseFromUser := `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult>
    <Arn>arn:aws:iam::123456789012:user/MyUserName</Arn>
    <UserId>ASOMETHINGSOMETHINGSOMETHING</UserId>
    <Account>123456789012</Account>
  </GetCallerIdentityResult>
  <ResponseMetadata>
    <RequestId>7f4fc40c-853a-11e6-8848-8d035d01eb87</RequestId>
  </ResponseMetadata>
</GetCallerIdentityResponse>`
	expectedUserArn := "arn:aws:iam::123456789012:user/MyUserName"

	responseFromAssumedRole := `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult>
  <Arn>arn:aws:sts::123456789012:assumed-role/RoleName/RoleSessionName</Arn>
  <UserId>ASOMETHINGSOMETHINGELSE:RoleSessionName</UserId>
    <Account>123456789012</Account>
  </GetCallerIdentityResult>
  <ResponseMetadata>
    <RequestId>7f4fc40c-853a-11e6-8848-8d035d01eb87</RequestId>
  </ResponseMetadata>
</GetCallerIdentityResponse>`
	expectedRoleArn := "arn:aws:sts::123456789012:assumed-role/RoleName/RoleSessionName"

	parsedUserResponse, err := parseGetCallerIdentityResponse(responseFromUser)
	if err != nil {
		t.Fatal(err)
	}
	if parsedArn := parsedUserResponse.GetCallerIdentityResult[0].Arn; parsedArn != expectedUserArn {
		t.Errorf("expected to parse arn %#v, got %#v", expectedUserArn, parsedArn)
	}

	parsedRoleResponse, err := parseGetCallerIdentityResponse(responseFromAssumedRole)
	if err != nil {
		t.Fatal(err)
	}
	if parsedArn := parsedRoleResponse.GetCallerIdentityResult[0].Arn; parsedArn != expectedRoleArn {
		t.Errorf("expected to parn arn %#v; got %#v", expectedRoleArn, parsedArn)
	}

	_, err = parseGetCallerIdentityResponse("SomeRandomGibberish")
	if err == nil {
		t.Errorf("expected to NOT parse random giberish, but didn't get an error")
	}
}

func TestBackend_pathLogin_parseIamArn(t *testing.T) {
	testParser := func(inputArn, expectedCanonicalArn string, expectedEntity iamEntity) {
		entity, err := parseIamArn(inputArn)
		if err != nil {
			t.Fatal(err)
		}
		if expectedCanonicalArn != "" && entity.canonicalArn() != expectedCanonicalArn {
			t.Fatalf("expected to canonicalize ARN %q into %q but got %q instead", inputArn, expectedCanonicalArn, entity.canonicalArn())
		}
		if *entity != expectedEntity {
			t.Fatalf("expected to get iamEntity %#v from input ARN %q but instead got %#v", expectedEntity, inputArn, *entity)
		}
	}

	testParser("arn:aws:iam::123456789012:user/UserPath/MyUserName",
		"arn:aws:iam::123456789012:user/MyUserName",
		iamEntity{Partition: "aws", AccountNumber: "123456789012", Type: "user", Path: "UserPath", FriendlyName: "MyUserName"},
	)
	canonicalRoleArn := "arn:aws:iam::123456789012:role/RoleName"
	testParser("arn:aws:sts::123456789012:assumed-role/RoleName/RoleSessionName",
		canonicalRoleArn,
		iamEntity{Partition: "aws", AccountNumber: "123456789012", Type: "assumed-role", FriendlyName: "RoleName", SessionInfo: "RoleSessionName"},
	)
	testParser("arn:aws:iam::123456789012:role/RolePath/RoleName",
		canonicalRoleArn,
		iamEntity{Partition: "aws", AccountNumber: "123456789012", Type: "role", Path: "RolePath", FriendlyName: "RoleName"},
	)
	testParser("arn:aws:iam::123456789012:instance-profile/profilePath/InstanceProfileName",
		"",
		iamEntity{Partition: "aws", AccountNumber: "123456789012", Type: "instance-profile", Path: "profilePath", FriendlyName: "InstanceProfileName"},
	)

	// Test that it properly handles pathological inputs...
	_, err := parseIamArn("")
	if err == nil {
		t.Error("expected error from empty input string")
	}

	_, err = parseIamArn("arn:aws:iam::123456789012:role")
	if err == nil {
		t.Error("expected error from malformed ARN without a role name")
	}

	_, err = parseIamArn("arn:aws:iam")
	if err == nil {
		t.Error("expected error from incomplete ARN (arn:aws:iam)")
	}

	_, err = parseIamArn("arn:aws:iam::1234556789012:/")
	if err == nil {
		t.Error("expected error from empty principal type and no principal name (arn:aws:iam::1234556789012:/)")
	}
}

func TestBackend_validateVaultHeaderValue(t *testing.T) {
	const canaryHeaderValue = "Vault-Server"
	requestURL, err := url.Parse("https://sts.amazonaws.com/")
	if err != nil {
		t.Fatalf("error parsing test URL: %v", err)
	}
	postHeadersMissing := http.Header{
		"Host":          []string{"Foo"},
		"Authorization": []string{"AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/iam/aws4_request, SignedHeaders=content-type;host;x-amz-date;x-vault-aws-iam-server-id, Signature=5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7"},
	}
	postHeadersInvalid := http.Header{
		"Host":            []string{"Foo"},
		iamServerIdHeader: []string{"InvalidValue"},
		"Authorization":   []string{"AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/iam/aws4_request, SignedHeaders=content-type;host;x-amz-date;x-vault-aws-iam-server-id, Signature=5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7"},
	}
	postHeadersUnsigned := http.Header{
		"Host":            []string{"Foo"},
		iamServerIdHeader: []string{canaryHeaderValue},
		"Authorization":   []string{"AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/iam/aws4_request, SignedHeaders=content-type;host;x-amz-date, Signature=5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7"},
	}
	postHeadersValid := http.Header{
		"Host":            []string{"Foo"},
		iamServerIdHeader: []string{canaryHeaderValue},
		"Authorization":   []string{"AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/iam/aws4_request, SignedHeaders=content-type;host;x-amz-date;x-vault-aws-iam-server-id, Signature=5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7"},
	}

	postHeadersSplit := http.Header{
		"Host":            []string{"Foo"},
		iamServerIdHeader: []string{canaryHeaderValue},
		"Authorization":   []string{"AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/iam/aws4_request", "SignedHeaders=content-type;host;x-amz-date;x-vault-aws-iam-server-id, Signature=5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7"},
	}

	err = validateVaultHeaderValue(postHeadersMissing, requestURL, canaryHeaderValue)
	if err == nil {
		t.Error("validated POST request with missing Vault header")
	}

	err = validateVaultHeaderValue(postHeadersInvalid, requestURL, canaryHeaderValue)
	if err == nil {
		t.Error("validated POST request with invalid Vault header value")
	}

	err = validateVaultHeaderValue(postHeadersUnsigned, requestURL, canaryHeaderValue)
	if err == nil {
		t.Error("validated POST request with unsigned Vault header")
	}

	err = validateVaultHeaderValue(postHeadersValid, requestURL, canaryHeaderValue)
	if err != nil {
		t.Errorf("did NOT validate valid POST request: %v", err)
	}

	err = validateVaultHeaderValue(postHeadersSplit, requestURL, canaryHeaderValue)
	if err != nil {
		t.Errorf("did NOT validate valid POST request with split Authorization header: %v", err)
	}
}

// TestBackend_pathLogin_IAMHeaders tests login with iam_request_headers,
// supporting both base64 encoded string and JSON headers
func TestBackend_pathLogin_IAMHeaders(t *testing.T) {
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

	// sets up a test server to stand in for STS service
	ts := setupIAMTestServer()
	defer ts.Close()

	clientConfigData := map[string]interface{}{
		"iam_server_id_header_value": testVaultHeaderValue,
		"sts_endpoint":               ts.URL,
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

	// create a role entry
	roleEntry := &awsRoleEntry{
		Version:  currentRoleStorageVersion,
		AuthType: iamAuthType,
	}

	if err := b.nonLockedSetAWSRole(context.Background(), storage, testValidRoleName, roleEntry); err != nil {
		t.Fatalf("failed to set entry: %s", err)
	}

	// create a baseline loginData map structure, including iam_request_headers
	// already base64encoded. This is the "Default" loginData used for all tests.
	// Each sub test can override the map's iam_request_headers entry
	loginData, err := defaultLoginData()
	if err != nil {
		t.Fatal(err)
	}

	// expected errors for certain tests
	missingHeaderErr := errors.New("error validating X-Vault-AWS-IAM-Server-ID header: missing header \"X-Vault-AWS-IAM-Server-ID\"")
	parsingErr := errors.New("error making upstream request: error parsing STS response")

	testCases := []struct {
		Name      string
		Header    interface{}
		ExpectErr error
	}{
		{
			Name: "Default",
		},
		{
			Name: "Map-complete",
			Header: map[string]interface{}{
				"Content-Length":            "43",
				"Content-Type":              "application/x-www-form-urlencoded; charset=utf-8",
				"User-Agent":                "aws-sdk-go/1.14.24 (go1.11; darwin; amd64)",
				"X-Amz-Date":                "20180910T203328Z",
				"X-Vault-Aws-Iam-Server-Id": "VaultAcceptanceTesting",
				"Authorization":             "AWS4-HMAC-SHA256 Credential=AKIAJPQ466AIIQW4LPSQ/20180910/us-east-1/sts/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-vault-aws-iam-server-id, Signature=cdef5819b2e97f1ff0f3e898fd2621aa03af00a4ec3e019122c20e5482534bf4",
			},
		},
		{
			Name: "Map-incomplete",
			Header: map[string]interface{}{
				"Content-Length": "43",
				"Content-Type":   "application/x-www-form-urlencoded; charset=utf-8",
				"User-Agent":     "aws-sdk-go/1.14.24 (go1.11; darwin; amd64)",
				"X-Amz-Date":     "20180910T203328Z",
				"Authorization":  "AWS4-HMAC-SHA256 Credential=AKIAJPQ466AIIQW4LPSQ/20180910/us-east-1/sts/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-vault-aws-iam-server-id, Signature=cdef5819b2e97f1ff0f3e898fd2621aa03af00a4ec3e019122c20e5482534bf4",
			},
			ExpectErr: missingHeaderErr,
		},
		{
			Name: "JSON-complete",
			Header: `{
				"Content-Length":"43",
				"Content-Type":"application/x-www-form-urlencoded; charset=utf-8",
				"User-Agent":"aws-sdk-go/1.14.24 (go1.11; darwin; amd64)",
				"X-Amz-Date":"20180910T203328Z",
				"X-Vault-Aws-Iam-Server-Id": "VaultAcceptanceTesting",
				"Authorization":"AWS4-HMAC-SHA256 Credential=AKIAJPQ466AIIQW4LPSQ/20180910/us-east-1/sts/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-vault-aws-iam-server-id, Signature=cdef5819b2e97f1ff0f3e898fd2621aa03af00a4ec3e019122c20e5482534bf4"
			}`,
		},
		{
			Name: "JSON-incomplete",
			Header: `{
				"Content-Length":"43",
				"Content-Type":"application/x-www-form-urlencoded; charset=utf-8",
				"User-Agent":"aws-sdk-go/1.14.24 (go1.11; darwin; amd64)",
				"X-Amz-Date":"20180910T203328Z",
				"X-Vault-Aws-Iam-Server-Id": "VaultAcceptanceTesting",
				"Authorization":"AWS4-HMAC-SHA256 Credential=AKIAJPQ466AIIQW4LPSQ/20180910/us-east-1/sts/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-vault-aws-iam-server-id"
			}`,
			ExpectErr: parsingErr,
		},
		{
			Name:   "Base64-complete",
			Header: base64Complete(),
		},
		{
			Name:      "Base64-incomplete-missing-header",
			Header:    base64MissingVaultID(),
			ExpectErr: missingHeaderErr,
		},
		{
			Name:      "Base64-incomplete-missing-auth-sig",
			Header:    base64MissingAuthField(),
			ExpectErr: parsingErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Header != nil {
				loginData["iam_request_headers"] = tc.Header
			}

			loginRequest := &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "login",
				Storage:   storage,
				Data:      loginData,
			}

			resp, err := b.HandleRequest(context.Background(), loginRequest)
			if err != nil || resp == nil || resp.IsError() {
				if tc.ExpectErr != nil && tc.ExpectErr.Error() == resp.Error().Error() {
					return
				}
				t.Errorf("un expected failed login:\nresp: %#v\n\nerr: %v", resp, err)
			}
		})
	}
}

func defaultLoginData() (map[string]interface{}, error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}

	stsService := sts.New(awsSession)
	stsInputParams := &sts.GetCallerIdentityInput{}
	stsRequestValid, _ := stsService.GetCallerIdentityRequest(stsInputParams)
	stsRequestValid.HTTPRequest.Header.Add(iamServerIdHeader, testVaultHeaderValue)
	stsRequestValid.HTTPRequest.Header.Add("Authorization", fmt.Sprintf("%s,%s,%s",
		"AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/iam/aws4_request",
		"SignedHeaders=content-type;host;x-amz-date;x-vault-aws-iam-server-id",
		"Signature=5d672d79c15b13162d9279b0855cfba6789a8edb4c82c400e06b5924a6f2b5d7"))
	stsRequestValid.Sign()

	return buildCallerIdentityLoginData(stsRequestValid.HTTPRequest, testValidRoleName)
}

// setupIAMTestServer configures httptest server to intercept and respond to the
// IAM login path's invocation of submitCallerIdentityRequest (which does not
// use the AWS SDK), which receieves the mocked response responseFromUser
// containing user information matching the role.
func setupIAMTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseString := `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetCallerIdentityResult>
    <Arn>arn:aws:iam::123456789012:user/valid-role</Arn>
    <UserId>ASOMETHINGSOMETHINGSOMETHING</UserId>
    <Account>123456789012</Account>
  </GetCallerIdentityResult>
  <ResponseMetadata>
    <RequestId>7f4fc40c-853a-11e6-8848-8d035d01eb87</RequestId>
  </ResponseMetadata>
</GetCallerIdentityResponse>`

		auth := r.Header.Get("Authorization")
		parts := strings.Split(auth, ",")
		for i, s := range parts {
			s = strings.TrimSpace(s)
			key := strings.Split(s, "=")
			parts[i] = key[0]
		}

		// verify the "Authorization" header contains all the expected parts
		expectedAuthParts := []string{"AWS4-HMAC-SHA256 Credential", "SignedHeaders", "Signature"}
		var matchingCount int
		for _, v := range parts {
			for _, z := range expectedAuthParts {
				if z == v {
					matchingCount++
				}
			}
		}
		if matchingCount != len(expectedAuthParts) {
			responseString = "missing auth parts"
		}
		fmt.Fprintln(w, responseString)
	}))
}

// base64Complete returns a base64 encoded auth header as expected
func base64Complete() string {
	min := `{"Authorization":["AWS4-HMAC-SHA256 Credential=AKIAJPQ466AIIQW4LPSQ/20180907/us-east-1/sts/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-vault-aws-iam-server-id, Signature=97086b0531854844099fc52733fa2c88a2bfb54b2689600c6e249358a8353b52"],"Content-Length":["43"],"Content-Type":["application/x-www-form-urlencoded; charset=utf-8"],"User-Agent":["aws-sdk-go/1.14.24 (go1.11; darwin; amd64)"],"X-Amz-Date":["20180907T222145Z"],"X-Vault-Aws-Iam-Server-Id":["VaultAcceptanceTesting"]}`
	return min
}

// base64MissingVaultID returns a base64 encoded auth header, that omits the
// Vault ID header
func base64MissingVaultID() string {
	min := `{"Authorization":["AWS4-HMAC-SHA256 Credential=AKIAJPQ466AIIQW4LPSQ/20180907/us-east-1/sts/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-vault-aws-iam-server-id, Signature=97086b0531854844099fc52733fa2c88a2bfb54b2689600c6e249358a8353b52"],"Content-Length":["43"],"Content-Type":["application/x-www-form-urlencoded; charset=utf-8"],"User-Agent":["aws-sdk-go/1.14.24 (go1.11; darwin; amd64)"],"X-Amz-Date":["20180907T222145Z"]}`
	return min
}

// base64MissingAuthField returns a base64 encoded Auth header, that omits the
// "Signature" part
func base64MissingAuthField() string {
	min := `{"Authorization":["AWS4-HMAC-SHA256 Credential=AKIAJPQ466AIIQW4LPSQ/20180907/us-east-1/sts/aws4_request, SignedHeaders=content-length;content-type;host;x-amz-date;x-vault-aws-iam-server-id"],"Content-Length":["43"],"Content-Type":["application/x-www-form-urlencoded; charset=utf-8"],"User-Agent":["aws-sdk-go/1.14.24 (go1.11; darwin; amd64)"],"X-Amz-Date":["20180907T222145Z"],"X-Vault-Aws-Iam-Server-Id":["VaultAcceptanceTesting"]}`
	return min
}
