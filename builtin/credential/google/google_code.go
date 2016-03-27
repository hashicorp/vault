package google

import (
	"github.com/tebeka/selenium"
	"fmt"
	"testing"
	"time"
)

func handleError(msg string, wd selenium.WebDriver, t *testing.T, err error) {
	if err != nil {
		currentURL, getURLErr := wd.CurrentURL()
		var errorURL string
		if getURLErr != nil {
			errorURL = "unknown url (url retrieval failed)"
		} else {
			errorURL = currentURL
		}
		t.Fatalf("error while %s at url %s, error details: %s\n", msg, errorURL, err.Error())
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