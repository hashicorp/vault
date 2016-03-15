package google

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/tebeka/selenium"
)

const GOOGLE_USERNAME_ENV_VAR_NAME = "GOOGLE_TESTING_ONLY_USERNAME"
const GOOGLE_PASSWORD_ENV_VAR_NAME = "GOOGLE_TESTING_ONLY_PASSWORD"
const GOOGLE_DOMAIN_ENV_VAR_NAME = "GOOGLE_DOMAIN"

func environmentVariable(name string) string {

	return os.Getenv(name)
}

func googleCode(t *testing.T) string {

	googleConfig := &oauth2.Config{
		ClientID:     "158113233735-figmusvbkf0ui8g8u58am2tkumf9cnl8.apps.googleusercontent.com",
		ClientSecret: "45UnnkbRwpUNkrCl9d8x3U48",
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{ "email" },
	}

	//TODO: nathang also read the test application data
	authUrl := googleConfig.AuthCodeURL("state", oauth2.AccessTypeOnline)

	user := environmentVariable(GOOGLE_USERNAME_ENV_VAR_NAME)
	pass := environmentVariable(GOOGLE_PASSWORD_ENV_VAR_NAME)

	caps := selenium.Capabilities {
		"browserName": "htmlunit",
		"javascriptEnabled": true,
	}

	wd, _ := selenium.NewRemote(caps, "")
	defer wd.Quit()

	wd.SetAsyncScriptTimeout(10000)
	wd.SetImplicitWaitTimeout(10000)
	wd.SetPageLoadTimeout(10000)
	wd.MaximizeWindow("")
	wd.Get(authUrl)
	fmt.Printf("authentication url: %s", authUrl)

	var currentUrl string

	currentUrl, _ = wd.CurrentURL()
	fmt.Printf("url before authentication %s\n", currentUrl)

	emailInput, err := wd.FindElement(selenium.ById, "Email")
	if err != nil {
		t.Error(err.Error())
	}
	err = emailInput.SendKeys(user)
	if err != nil {
		t.Error(err.Error())
	}

	passInput, err := wd.FindElement(selenium.ById, "Passwd")
	if err != nil {
		t.Error(err.Error())
	}
	err = passInput.SendKeys(pass)
	if err != nil {
		t.Error(err.Error())
	}

	authenticateButton, err := wd.FindElement(selenium.ById, "signIn")
	if err != nil {
		t.Error(err.Error())
	}
	err = authenticateButton.Click()
	if err != nil {
		t.Error(err.Error())
	}

	currentUrl, _ = wd.CurrentURL()
	fmt.Printf("url after authentication %s\n", currentUrl)

	time.Sleep(10 * time.Second)

	_, err = wd.ExecuteScript(`
	document.getElementById("submit_approve_access").click();
	`, []interface{}{})
	if err != nil {
		t.Error(err.Error())
	}

	currentUrl, _ = wd.CurrentURL()
	fmt.Printf("url after authorize %s\n", currentUrl)

	codeElement, err := wd.FindElement(selenium.ById, "code")
	if err != nil {
		t.Error(err.Error())
	}

	code, err := codeElement.GetAttribute("value")
	if err != nil {
		t.Error(err.Error())
	}

	return code
}

func loginData(t *testing.T) map[string]interface{} {
	return map[string]interface{}{
		"code": googleCode(t),
	}
}

func googleDomain() string {
	return environmentVariable(GOOGLE_DOMAIN_ENV_VAR_NAME)
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

	config_data1 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "",
		"max_ttl":      "",
	}
	expectedTTL1, _ := time.ParseDuration("24h0m0s")
	config_data2 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "1h",
		"max_ttl":      "2h",
	}
	expectedTTL2, _ := time.ParseDuration("1h0m0s")
	config_data3 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "50h",
		"max_ttl":      "50h",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, config_data1),
			testLoginWrite(t, expectedTTL1.Nanoseconds(), false),
			testConfigWrite(t, config_data2),
			testLoginWrite(t, expectedTTL2.Nanoseconds(), false),
			testConfigWrite(t, config_data3),
			testLoginWrite(t, 0, true),
		},
	})
}

func testLoginWrite(t *testing.T, expectedTTL int64, expectFail bool) logicaltest.TestStep {

	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		ErrorOk:   true,
		Data:      nil,
		PreFlight: func(r *logical.Request) error {
			r.Data = loginData(t)
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
		Path:      "config",
		Data:      d,
	}
}

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

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccMap(t, "default", "root"),
			testAccLogin(t, []string{"root"}),
		},
	})
}

func testAccPreCheck(t *testing.T) {

	//TODO: nathang - make sure selenium server is on? start selenium server? stop selenium server?

	requiredEnvVars := []string{ GOOGLE_USERNAME_ENV_VAR_NAME, GOOGLE_PASSWORD_ENV_VAR_NAME }
	for _, envVar := range requiredEnvVars {
		if value := environmentVariable(envVar); value == "" {
			t.Fatal(fmt.Sprintf("missing environment variable %s", envVar))
		}
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"domain": googleDomain(),
		},
	}
}

func testAccMap(t *testing.T, k string, v string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "map/teams/" + k,
		Data: map[string]interface{}{
			"value": v,
		},
	}
}

func testAccLogin(t *testing.T, keys []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"code": googleCode(t),
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth(keys),
	}
}
