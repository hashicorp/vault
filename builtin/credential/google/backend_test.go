package google

import (
	"fmt"
	"testing"
	"time"
	"net/http"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/hashicorp/vault/logical/framework"
)

const GoogleApplicationIDEnvVarName = "GOOGLE_TESTING_ONLY_APPLICATION_ID"
const GoogleApplicationSecretEnvVarName = "GOOGLE_TESTING_ONLY_APPLICATION_SECRET"
const GoogleDomainEnvVarName = "GOOGLE_DOMAIN"

func googleClientID() string {
	return environmentVariable(GoogleApplicationIDEnvVarName)
}

func googleClientSecret() string {
	return environmentVariable(GoogleApplicationSecretEnvVarName)
}

func loginData(t *testing.T, authCodeURL string) map[string]interface{} {
	return map[string]interface{}{
		googleAuthCodeParameterName: googleCode(t, authCodeURL),
	}
}

func googleDomain() string {
	return environmentVariable(GoogleDomainEnvVarName)
}

func googleUsername() string {

	user := googleUser()
	name := localPartFromEmail(user)

	return name
}


func TestBackend_Config(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 2
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	stepsSharedState := &sharedTestState{}

	configData1 := map[string]interface{}{
		domainConfigPropertyName: googleDomain(),
		TTLConfigPropertyName:          "",
		maxTTLConfigPropertyName:      "",
		applicationIdConfigPropertyName: googleClientID(),
		applicationSecretConfigPropertyName: googleClientSecret(),
	}
	expectedTTL1, _ := time.ParseDuration("24h0m0s")
	configData2 := map[string]interface{}{
		domainConfigPropertyName: googleDomain(),
		TTLConfigPropertyName:          "1h",
		maxTTLConfigPropertyName:      "2h",
		applicationIdConfigPropertyName: googleClientID(),
		applicationSecretConfigPropertyName: googleClientSecret(),
	}
	expectedTTL2, _ := time.ParseDuration("1h0m0s")
	configData3 := map[string]interface{}{
		domainConfigPropertyName: googleDomain(),
		TTLConfigPropertyName:          "50h",
		maxTTLConfigPropertyName:      "50h",
		applicationIdConfigPropertyName: googleClientID(),
		applicationSecretConfigPropertyName: googleClientSecret(),
	}

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, configData1),
			testAuthCodeURL(t, stepsSharedState),
			testLoginWrite(t, stepsSharedState, expectedTTL1.Nanoseconds(), false),
			testConfigWrite(t, configData2),
			testLoginWrite(t, stepsSharedState, expectedTTL2.Nanoseconds(), false),
			testConfigWrite(t, configData3),
			testLoginWrite(t, stepsSharedState, 0, true),
		},
	})
}

type sharedTestState struct {
	authCodeUrl string
	auth *logical.Auth
}

func testLoginWrite(t *testing.T, stepsSharedState *sharedTestState, expectedTTL int64, expectFail bool) logicaltest.TestStep {

	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      loginPath,
		ErrorOk:   true,
		Data:      nil,
		PreFlight: func(r *logical.Request) error {
			r.Data = loginData(t, stepsSharedState.authCodeUrl)
			return nil
		},
		Check: func(resp *logical.Response) error {
			if resp.IsError() && expectFail {
				return nil
			}
			var actualTTL int64
			actualTTL = resp.Auth.LeaseOptions.TTL.Nanoseconds()
			if actualTTL != expectedTTL {
				return fmt.Errorf("TTL mismatched. Expected: %d Actual: %d", expectedTTL, resp.Auth.LeaseOptions.TTL.Nanoseconds())
			}
			return nil
		},
	}
}


func testConfigWrite(t *testing.T, d map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      configPath,
		Data:      d,
	}
}

const testPolicy = "myVeryOwnTestPolicy"

func TestBackend_basic(t *testing.T) {
	defaultLeaseTTLVal := time.Hour * 24
	maxLeaseTTLVal := time.Hour * 24 * 30
	b, err := Factory(&logical.BackendConfig{
		Logger: nil,
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	stepsSharedState := &sharedTestState{}

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAuthCodeURL(t, stepsSharedState),
			testUserPolicyMap(t, googleUsername(), testPolicy),
			testAccLogin(t, stepsSharedState, []string{ testPolicy, "default" }),
		},
	})
}

func testAuthCodeURL(t *testing.T, stepsSharedState *sharedTestState) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path: codeURLPath,
		Check: func(resp *logical.Response) error {
			stepsSharedState.authCodeUrl = resp.Data[codeURLResponsePropertyName].(string)
			return nil
		},
	}
}

func testAccPreCheck(t *testing.T) {

	requiredEnvVars := []string{
		GoogleUsernameEnvVarName,
		GooglePasswordEnvVarName,
		GoogleApplicationIDEnvVarName,
		GoogleApplicationSecretEnvVarName,
	}

	for _, envVar := range requiredEnvVars {
		if value := environmentVariable(envVar); value == "" {
			t.Fatal(fmt.Sprintf("missing environment variable %s", envVar))
		}
	}

	_, err := http.Get("http://127.0.0.1:4444/wd/hub")
	if (err != nil) {
		t.Fatal("google integration tests require selenium server. use: 'make selenium testacc' instead of 'make testacc' to have one provided for you.")
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      configPath,
		Data: map[string]interface{}{
			domainConfigPropertyName: googleDomain(),
			applicationIdConfigPropertyName: googleClientID(),
			applicationSecretConfigPropertyName: googleClientSecret(),
		},
	}
}

func testUserPolicyMap(t *testing.T, k string, v string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "map/" + usersToPoliciesMapPath + "/" + k,
		Data: map[string]interface{}{
			"value": v,
		},
	}
}

func testAccLogin(t *testing.T, stepsSharedState *sharedTestState, keys []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      loginPath,
		Data:	nil,
		PreFlight: func(r *logical.Request) error {
			r.Data = loginData(t, stepsSharedState.authCodeUrl)
			return nil
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth(keys),
	}
}

type AuthContainer struct {
	Auth *logical.Auth
}

func Test_Renew(t *testing.T) {
	storage := &logical.InmemStorage{}
	ttl, _ := time.ParseDuration("1m")
	conf, err := logical.StorageEntryJSON("config", config{
		Domain:     googleDomain(),
		TTL:        ttl,
		MaxTTL:     ttl,
		ApplicationID: googleClientID(),
		ApplicationSecret: googleClientSecret(),
	})
	if err != nil {
		t.Fatal(err)
	}
	storage.Put(conf)

	lb, err := Factory(&logical.BackendConfig{
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: 300 * time.Second,
			MaxLeaseTTLVal:     1800 * time.Second,
		},
		StorageView: storage,
	})
	if err != nil {
		t.Fatal("error: %s", err)
	}

	b := lb.(*backend)

	if err != nil {
		t.Fatal(err)
	}

	code := codeUrl(googleClientID(), googleClientSecret())

	fd := &framework.FieldData{
		Raw: map[string]interface{}{
			googleAuthCodeParameterName: googleCode(t, code),
		},
		Schema: pathLogin(b).Fields,
	}

	req := &logical.Request{
		Storage: storage,
		Auth:    &logical.Auth{},
	}

	resp, err := b.pathLogin(req, fd)
	if err != nil {
		t.Fatal(err)
	}

	//serialization and deserialization happens in real flow and effects the structure of the token
	authContainer := &AuthContainer{
		Auth: resp.Auth,
	}
	storageEntry, err := logical.StorageEntryJSON("key", authContainer)
	if (err != nil) {
		t.Fatal(err)
	}
	storageEntry.DecodeJSON(&authContainer)

	req.Auth.InternalData = authContainer.Auth.InternalData
	req.Auth.Metadata = resp.Auth.Metadata
	req.Auth.LeaseOptions = resp.Auth.LeaseOptions
	req.Auth.IssueTime = time.Now()
	req.Auth.Policies = append(resp.Auth.Policies, "default") //this happens in core

	resp, err = b.pathLoginRenew(req, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}
	if resp.IsError() {
		t.Fatalf("got error: %#v", *resp)
	}
}

