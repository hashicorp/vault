package duo

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"github.com/duosecurity/duo_api_golang/authapi"
	"github.com/hashicorp/vault/logical"
)

type MockClientData struct {
	PreauthData *authapi.PreauthResult
	PreauthError error
	AuthData *authapi.AuthResult
	AuthError error
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

func MockGetDuoAuthClient(data *MockClientData) func (*logical.Request, *DuoConfig) (AuthClient, error) {
	return func (*logical.Request, *DuoConfig) (AuthClient, error) {
		return getDuoAuthClient(data), nil
	}
}

func getDuoAuthClient(data *MockClientData) AuthClient {
	var c MockAuthClient
	// set default response to auth user 
	if data.PreauthData == nil {
		data.PreauthData = &authapi.PreauthResult{}
		json.Unmarshal([]byte(`
{
  "Stat": "OK",
  "Response": {
    "Result": "auth",
    "Status_Msg": "Needs authentication",
    "Devices": []
  }
}`), data.PreauthData)
	}

	if data.AuthData == nil {
		data.AuthData = &authapi.AuthResult{}
		json.Unmarshal([]byte(`
{
  "Stat": "OK",
  "Response": {
    "Result": "allow"
  }
}`), data.AuthData)
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
	resp, err := duoHandler(duoConfig, duoAuthClient, successResp, "user", "", "", "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if resp != successResp {
		t.Fatalf("Testing Duo authentication gave incorrect response (expected success, got: %v)", resp)
	}
}

func TestDuoHandlerReject(t *testing.T) {
	AuthData := &authapi.AuthResult{}
		json.Unmarshal([]byte(`
{
  "Stat": "OK",
  "Response": {
    "Result": "deny",
    "Status_Msg": "Invalid auth"
  }
}`), AuthData)
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
	resp, err := duoHandler(duoConfig, duoAuthClient, successResp, "user", "", "", "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	error, ok := resp.Data["error"].(string)
	if !ok || !strings.Contains(error, expectedError) {
		t.Fatalf("Testing Duo authentication gave incorrect response (expected deny, got: %v)", error)
	}
}
