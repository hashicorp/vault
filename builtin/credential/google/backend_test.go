package google

import (
	"fmt"
	"os"
	"testing"
	"time"
	"net/http"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/tebeka/selenium"
)

const GoogleApplicationIDEnvVarName = "GOOGLE_TESTING_ONLY_APPLICATION_ID"
const GoogleApplicationSecretEnvVarName = "GOOGLE_TESTING_ONLY_APPLICATION_SECRET"
const GoogleUsernameEnvVarName = "GOOGLE_TESTING_ONLY_USERNAME"
const GooglePasswordEnvVarName = "GOOGLE_TESTING_ONLY_PASSWORD"
const GoogleDomainEnvVarName = "GOOGLE_DOMAIN"

func environmentVariable(name string) string {

	return os.Getenv(name)
}

func googleClientID() string {
	return environmentVariable(GoogleApplicationIDEnvVarName)
}

func googleClientSecret() string {
	return environmentVariable(GoogleApplicationSecretEnvVarName)
}

func handleError(msg string, wd selenium.WebDriver, t *testing.T, err error) {
	if err != nil {
		currentURL, getURLErr := wd.CurrentURL()
		var errorURL string
		if getURLErr != nil {
			errorURL = "unknown url (url retrieval failed)"
		} else {
			errorURL = currentURL
		}
		t.Errorf("error while %s at url %s, error details: %s\n", msg, errorURL, err.Error())
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

func googleCode(t *testing.T, authURL string) string {

	user := environmentVariable(GoogleUsernameEnvVarName)
	pass := environmentVariable(GooglePasswordEnvVarName)

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
	wd.Get(authURL)

	var err error

	emailInput, _ := element("Email", wd, t)
	err = emailInput.SendKeys(user)
	handleError("filling out user text box", wd, t, err)

	//two flows here, one fill out user + pass, the other fill user, click next, enter pass...
	passwordTextInputID := "Passwd"
	passInput, err := wd.FindElement(selenium.ById, passwordTextInputID)
	if err != nil {
		nextButton, _ := element("next", wd, t)
		err = nextButton.Click()
		handleError("clicking next after inserting email", wd, t, err)
		passInput, _ = element(passwordTextInputID, wd, t)
	}
	err = passInput.SendKeys(pass)
	handleError("filling out password text box", wd, t, err)

	authenticateButton, err := element("signIn", wd, t)
	err = authenticateButton.Click()
	handleError("clicking sign in after filling password", wd, t, err)

	authorizeButtonID := "submit_approve_access"
	authorizeButton, _ := element(authorizeButtonID, wd, t)
	authorizationButtonEnabled, _ :=  authorizeButton.IsEnabled()
	for i := 0 ; (!authorizationButtonEnabled) && (i < 100) ; i++ {
		time.Sleep(100 * time.Millisecond)
		authorizationButtonEnabled, _ =  authorizeButton.IsEnabled()
	}
	_, err = wd.ExecuteScript(fmt.Sprintf(`document.getElementById("%s").click();`, authorizeButtonID), []interface{}{})
	handleError("authorizing application with required permissions", wd, t, err)

	codeElement, err := element("code", wd, t)
	code, err := codeElement.GetAttribute("value")
	handleError("retrieving value of code", wd, t, err)

	return code
}

func loginData(t *testing.T, authCodeURL string) map[string]interface{} {
	return map[string]interface{}{
		"code": googleCode(t, authCodeURL),
	}
}

func googleDomain() string {
	return environmentVariable(GoogleDomainEnvVarName)
}

const sharedAuthCodeURL = "auth_code"

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

	configData1 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "",
		"max_ttl":      "",
		"applicationId": googleClientID(),
		"applicationSecret": googleClientSecret(),
	}
	expectedTTL1, _ := time.ParseDuration("24h0m0s")
	configData2 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "1h",
		"max_ttl":      "2h",
		"applicationId": googleClientID(),
		"applicationSecret": googleClientSecret(),
	}
	expectedTTL2, _ := time.ParseDuration("1h0m0s")
	configData3 := map[string]interface{}{
		"domain": googleDomain(),
		"ttl":          "50h",
		"max_ttl":      "50h",
		"applicationId": googleClientID(),
		"applicationSecret": googleClientSecret(),
	}

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, configData1),
			testAuthCodeURL(t, &stepsSharedState),
			testLoginWrite(t, &stepsSharedState, expectedTTL1.Nanoseconds(), false),
			testConfigWrite(t, configData2),
			testLoginWrite(t, &stepsSharedState, expectedTTL2.Nanoseconds(), false),
			testConfigWrite(t, configData3),
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
			r.Data = loginData(t, (*stepsSharedState)[sharedAuthCodeURL])
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
			testAuthCodeURL(t, &stepsSharedState),
			testAccMap(t, "default", "root"),
			testAccLogin(t, &stepsSharedState, []string{"root"}),
		},
	})
}

func testAuthCodeURL(t *testing.T, stepsSharedState *map[string]string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path: codeURLPath,
		Check: func(resp *logical.Response) error {
			(*stepsSharedState)[sharedAuthCodeURL] = resp.Data["url"].(string)
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
		Path:      "config",
		Data: map[string]interface{}{
			"domain": googleDomain(),
			"applicationId": googleClientID(),
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
			r.Data = loginData(t, (*stepsSharedState)[sharedAuthCodeURL])
			return nil
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckAuth(keys),
	}
}
