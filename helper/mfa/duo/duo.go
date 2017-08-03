// Package duo provides a Duo MFA handler to authenticate users
// with Duo. This handler is registered as the "duo" type in
// mfa_config.
package duo

import (
	"fmt"
	"net/url"

	"github.com/duosecurity/duo_api_golang/authapi"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// DuoPaths returns path functions to configure Duo.
func DuoPaths() []*framework.Path {
	return []*framework.Path{
		pathDuoConfig(),
		pathDuoAccess(),
	}
}

// DuoRootPaths returns the paths that are used to configure Duo.
func DuoRootPaths() []string {
	return []string{
		"duo/access",
		"duo/config",
	}
}

// DuoHandler interacts with the Duo Auth API to authenticate a user
// login request. If successful, the original response from the login
// backend is returned.
func DuoHandler(req *logical.Request, d *framework.FieldData, resp *logical.Response) (
	*logical.Response, error) {
	duoConfig, err := GetDuoConfig(req)
	if err != nil || duoConfig == nil {
		return logical.ErrorResponse("Could not load Duo configuration"), nil
	}

	duoAuthClient, err := GetDuoAuthClient(req, duoConfig)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	username, ok := resp.Auth.Metadata["username"]
	if !ok {
		return logical.ErrorResponse("Could not read username for MFA"), nil
	}

	var request *duoAuthRequest = &duoAuthRequest{}
	request.successResp = resp
	request.username = username
	request.method = d.Get("method").(string)
	request.passcode = d.Get("passcode").(string)
	request.ipAddr = req.Connection.RemoteAddr

	return duoHandler(duoConfig, duoAuthClient, request)
}

type duoAuthRequest struct {
	successResp *logical.Response
	username    string
	method      string
	passcode    string
	ipAddr      string
}

func duoHandler(duoConfig *DuoConfig, duoAuthClient AuthClient, request *duoAuthRequest) (
	*logical.Response, error) {

	duoUser := fmt.Sprintf(duoConfig.UsernameFormat, request.username)

	preauth, err := duoAuthClient.Preauth(
		authapi.PreauthUsername(duoUser),
		authapi.PreauthIpAddr(request.ipAddr),
	)

	if err != nil || preauth == nil {
		return logical.ErrorResponse("Could not call Duo preauth"), nil
	}

	if preauth.StatResult.Stat != "OK" {
		errorMsg := "Could not look up Duo user information"
		if preauth.StatResult.Message != nil {
			errorMsg = errorMsg + ": " + *preauth.StatResult.Message
		}
		if preauth.StatResult.Message_Detail != nil {
			errorMsg = errorMsg + " (" + *preauth.StatResult.Message_Detail + ")"
		}
		return logical.ErrorResponse(errorMsg), nil
	}

	switch preauth.Response.Result {
	case "allow":
		return request.successResp, err
	case "deny":
		return logical.ErrorResponse(preauth.Response.Status_Msg), nil
	case "enroll":
		return logical.ErrorResponse(fmt.Sprintf("%s (%s)",
			preauth.Response.Status_Msg,
			preauth.Response.Enroll_Portal_Url)), nil
	case "auth":
		break
	default:
		return logical.ErrorResponse(fmt.Sprintf("Invalid Duo preauth response: %s",
			preauth.Response.Result)), nil
	}

	options := []func(*url.Values){authapi.AuthUsername(duoUser)}
	if request.method == "" {
		request.method = "auto"
	}
	if request.method == "auto" || request.method == "push" {
		if duoConfig.PushInfo != "" {
			options = append(options, authapi.AuthPushinfo(duoConfig.PushInfo))
		}
	}
	if request.passcode != "" {
		request.method = "passcode"
		options = append(options, authapi.AuthPasscode(request.passcode))
	} else {
		options = append(options, authapi.AuthDevice("auto"))
	}

	result, err := duoAuthClient.Auth(request.method, options...)

	if err != nil || result == nil {
		return logical.ErrorResponse("Could not call Duo auth"), nil
	}

	if result.StatResult.Stat != "OK" {
		errorMsg := "Could not authenticate Duo user"
		if result.StatResult.Message != nil {
			errorMsg = errorMsg + ": " + *result.StatResult.Message
		}
		if result.StatResult.Message_Detail != nil {
			errorMsg = errorMsg + " (" + *result.StatResult.Message_Detail + ")"
		}
		return logical.ErrorResponse(errorMsg), nil
	}

	if result.Response.Result != "allow" {
		return logical.ErrorResponse(result.Response.Status_Msg), nil
	}

	return request.successResp, nil
}
