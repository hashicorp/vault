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

const GOOGLE_APPLICATION_ID_ENV_VAR_NAME = "GOOGLE_TESTING_ONLY_APPLICATION_ID"
const GOOGLE_APPLICATION_SECRET_ENV_VAR_NAME = "GOOGLE_TESTING_ONLY_APPLICATION_SECRET"
const GOOGLE_USERNAME_ENV_VAR_NAME = "GOOGLE_TESTING_ONLY_USERNAME"
const GOOGLE_PASSWORD_ENV_VAR_NAME = "GOOGLE_TESTING_ONLY_PASSWORD"
const GOOGLE_DOMAIN_ENV_VAR_NAME = "GOOGLE_DOMAIN"

func environmentVariable(name string) string {

	return os.Getenv(name)
}

func googleClientId() string {
	return environmentVariable(GOOGLE_APPLICATION_ID_ENV_VAR_NAME)
}

func googleClientSecret() string {
	return environmentVariable(GOOGLE_APPLICATION_SECRET_ENV_VAR_NAME)
}

func handleError(msg string, wd selenium.WebDriver, t *testing.T, err error) {
	if err != nil {
		currentUrl, getUrlErr := wd.CurrentURL()
		var errorUrl string
		if getUrlErr != nil {
			errorUrl = "unknown url (url retrieval failed)"
		} else {
			errorUrl = currentUrl
		}
		t.Error(fmt.Sprintf("error while %s at url %s, error details: %s\n", msg, errorUrl, err.Error()))
	}
	return
}

func handleFindElementError(id string, wd selenium.WebDriver, t *testing.T, err error) {
	msg := fmt.Sprintf("retrieving element by id %s", id)
	handleError(msg, wd, t, err)
	return
}

func element(id string, wd selenium.WebDriver, t *testing.T) (selenium.WebElement, error) {

	element, err := wd.FindElement(selenium.ById, id)

	handleFindElementError(id, wd, t, err)

	return element, err
}

func googleCode(t *testing.T, authCodeUrl string) string {

	googleConfig := &oauth2.Config{
		ClientID:     googleClientId(),
		ClientSecret: googleClientSecret(),
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{ "email" },
	}

	authUrl := googleConfig.AuthCodeURL("state", oauth2.AccessTypeOnline, oauth2.ApprovalForce)

	user := environmentVariable(GOOGLE_USERNAME_ENV_VAR_NAME)
	pass := environmentVariable(GOOGLE_PASSWORD_ENV_VAR_NAME)

	caps := selenium.Capabilities {
		"browserName": "htmlunit",
		"javascriptEnabled": true,
	}

	wd, _ := selenium.NewRemote(caps, "")
	defer wd.Quit()

	wd.SetAsyncScriptTimeout(30000)
	wd.SetImplicitWaitTimeout(30000)
	wd.SetPageLoadTimeout(30000)
	wd.MaximizeWindow("")
	wd.Get(authUrl)

	var err error

	emailInput, _ := element("Email", wd, t)
	err = emailInput.SendKeys(user)
	handleError("filling out user text box", wd, t, err)

	//two flows here, one fill out user + pass, the other fill user, click next, enter pass...
	passwordTextInputId := "Passwd"
	passInput, err := wd.FindElement(selenium.ById, passwordTextInputId)
	if err != nil {
		nextButton, _ := element("next", wd, t)
		err = nextButton.Click()
		handleError("clicking next after inserting email", wd, t, err)
		passInput, _ = element(passwordTextInputId, wd, t)
	}
	err = passInput.SendKeys(pass)
	handleError("filling out password text box", wd, t, err)

	authenticateButton, err := element("signIn", wd, t)
	err = authenticateButton.Click()
	handleError("clicking sign in after filling password", wd, t, err)

	authorizeButtonId := "submit_approve_access"
	authorizeButton, _ := element(authorizeButtonId, wd, t)
	authorizationButtonEnabled, _ :=  authorizeButton.IsEnabled()
	for i := 0 ; (!authorizationButtonEnabled) && (i < 100) ; i++ {
		time.Sleep(100 * time.Millisecond)
		authorizationButtonEnabled, _ =  authorizeButton.IsEnabled()
	}
	_, err = wd.ExecuteScript(fmt.Sprintf(`document.getElementById("%s").click();`, authorizeButtonId), []interface{}{})
	handleError("authorizing application with required permissions", wd, t, err)

	codeElement, err := element("code", wd, t)
	code, err := codeElement.GetAttribute("value")
	handleError("retrieving value of code", wd, t, err)

	return code
}

func loginData(t *testing.T, authCodeUrl string) map[string]interface{} {
	return map[string]interface{}{
		"code": googleCode(t, authCodeUrl),
	}
}

func googleDomain() string {
	return environmentVariable(GOOGLE_DOMAIN_ENV_VAR_NAME)
}

const SHARED_AUTH_CODE_URL = "auth_code"

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

	stepsSharedState := map[string]string{}

	config_data1 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "",
		"max_ttl":      "",
		"applicationId": googleClientId(),
		"applicationSecret": googleClientSecret(),
	}
	expectedTTL1, _ := time.ParseDuration("24h0m0s")
	config_data2 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "1h",
		"max_ttl":      "2h",
		"applicationId": googleClientId(),
		"applicationSecret": googleClientSecret(),
	}
	expectedTTL2, _ := time.ParseDuration("1h0m0s")
	config_data3 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "50h",
		"max_ttl":      "50h",
		"applicationId": googleClientId(),
		"applicationSecret": googleClientSecret(),
	}

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, config_data1),
			testAuthCodeUrl(t, &stepsSharedState),
			testLoginWrite(t, &stepsSharedState, expectedTTL1.Nanoseconds(), false),
			testConfigWrite(t, config_data2),
			testLoginWrite(t, &stepsSharedState, expectedTTL2.Nanoseconds(), false),
			testConfigWrite(t, config_data3),
			testLoginWrite(t, &stepsSharedState, 0, true),
		},
	})
}

func testLoginWrite(t *testing.T, stepsSharedState *map[string]string, expectedTTL int64, expectFail bool) logicaltest.TestStep {

	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		ErrorOk:   true,
		Data:      nil,
		PreFlight: func(r *logical.Request) error {
			r.Data = loginData(t, (*stepsSharedState)[SHARED_AUTH_CODE_URL])
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

	stepsSharedState := map[string]string{}

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAuthCodeUrl(t, &stepsSharedState),
			testAccMap(t, "default", "root"),
			testAccLogin(t, &stepsSharedState, []string{"root"}),
		},
	})
}

func testAuthCodeUrl(t *testing.T, stepsSharedState *map[string]string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path: PATH_CODE_URL,
		Check: func(resp *logical.Response) error {
			(*stepsSharedState)[SHARED_AUTH_CODE_URL] = resp.Data["url"].(string)
			return nil
		},
	}
}

func testAccPreCheck(t *testing.T) {

	//TODO: nathang - make sure selenium server is on? start selenium server? stop selenium server?

	requiredEnvVars := []string{
		GOOGLE_USERNAME_ENV_VAR_NAME,
		GOOGLE_PASSWORD_ENV_VAR_NAME,
		GOOGLE_APPLICATION_ID_ENV_VAR_NAME,
		GOOGLE_APPLICATION_SECRET_ENV_VAR_NAME,
	}

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
			"applicationId": googleClientId(),
			"applicationSecret": googleClientSecret(),
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

func testAccLogin(t *testing.T, stepsSharedState *map[string]string, keys []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		Data:	nil,
		PreFlight: func(r *logical.Request) error {
			r.Data = loginData(t, (*stepsSharedState)[SHARED_AUTH_CODE_URL])
			return nil
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth(keys),
	}
}
