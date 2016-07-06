package duo

import (
	"net/url"
	"strings"
	"testing"

	"github.com/duosecurity/duo_api_golang/authapi"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
)

type MockClientData struct {
	PreauthData  *authapi.PreauthResult
	PreauthError error
	AuthData     *authapi.AuthResult
	AuthError    error
}

type MockAuthClient struct {
	MockData *MockClientData
}

func (c *MockAuthClient) Preauth(options ...func(*url.Values)) (*authapi.PreauthResult, error) {
	return c.MockData.PreauthData, c.MockData.PreauthError
}

func (c *MockAuthClient) Auth(factor string, options ...func(*url.Values)) (*authapi.AuthResult, error) {
	return c.MockData.AuthData, c.MockData.AuthError
}

func MockGetDuoAuthClient(data *MockClientData) func(*logical.Request, *DuoConfig) (AuthClient, error) {
	return func(*logical.Request, *DuoConfig) (AuthClient, error) {
		return getDuoAuthClient(data), nil
	}
}

func getDuoAuthClient(data *MockClientData) AuthClient {
	var c MockAuthClient
	// set default response to be successful
	preauthSuccessJSON := `
	{
	  "Stat": "OK",
	  "Response": {
	    "Result": "auth",
	    "Status_Msg": "Needs authentication",
	    "Devices": []
	  }
	}`
	if data.PreauthData == nil {
		data.PreauthData = &authapi.PreauthResult{}
		jsonutil.DecodeJSON([]byte(preauthSuccessJSON), data.PreauthData)
	}

	authSuccessJSON := `
	{
	  "Stat": "OK",
	  "Response": {
	    "Result": "allow"
	  }
	}`
	if data.AuthData == nil {
		data.AuthData = &authapi.AuthResult{}
		jsonutil.DecodeJSON([]byte(authSuccessJSON), data.AuthData)
	}

	c.MockData = data
	return &c
}

func TestDuoHandlerSuccess(t *testing.T) {
	successResp := &logical.Response{
		Auth: &logical.Auth{},
	}
	duoConfig := &DuoConfig{
		UsernameFormat: "%s",
	}
	duoAuthClient := getDuoAuthClient(&MockClientData{})
	resp, err := duoHandler(duoConfig, duoAuthClient, &duoAuthRequest{
		successResp: successResp,
		username:    "",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if resp != successResp {
		t.Fatalf("Testing Duo authentication gave incorrect response (expected success, got: %v)", resp)
	}
}

func TestDuoHandlerReject(t *testing.T) {
	AuthData := &authapi.AuthResult{}
	authRejectJSON := `
	{
	  "Stat": "OK",
	  "Response": {
	    "Result": "deny",
	    "Status_Msg": "Invalid auth"
	  }
	}`
	jsonutil.DecodeJSON([]byte(authRejectJSON), AuthData)
	successResp := &logical.Response{
		Auth: &logical.Auth{},
	}
	expectedError := AuthData.Response.Status_Msg
	duoConfig := &DuoConfig{
		UsernameFormat: "%s",
	}
	duoAuthClient := getDuoAuthClient(&MockClientData{
		AuthData: AuthData,
	})
	resp, err := duoHandler(duoConfig, duoAuthClient, &duoAuthRequest{
		successResp: successResp,
		username:    "user",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	error, ok := resp.Data["error"].(string)
	if !ok || !strings.Contains(error, expectedError) {
		t.Fatalf("Testing Duo authentication gave incorrect response (expected deny, got: %v)", error)
	}
}
