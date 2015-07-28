package duo

import (
	"fmt"
	"net/url"

	"github.com/duosecurity/duo_api_golang/authapi"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func DuoPaths() []*framework.Path {
	return []*framework.Path{
		pathDuoConfig(),
		pathDuoAccess(),
	}
}

func DuoPathsSpecial() []string {
	return []string {
		"duo/access",
		"duo/config",
	}
}

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

	method := d.Get("method").(string)
	passcode := d.Get("passcode").(string)

	return duoHandler(duoConfig, duoAuthClient, resp,
		username, method, passcode, req.Connection.RemoteAddr)
}

func duoHandler(
	duoConfig *DuoConfig, duoAuthClient AuthClient, successResp *logical.Response,
	username string, method string, passcode string, ipAddr string) (*logical.Response, error) {

	duoUser := fmt.Sprintf(duoConfig.UsernameFormat, username)

	preauth, err := duoAuthClient.Preauth(
		authapi.PreauthUsername(duoUser),
		authapi.PreauthIpAddr(ipAddr),
	)

	if preauth.StatResult.Stat != "OK" {
		return logical.ErrorResponse(fmt.Sprintf("Could not look up Duo user information: %v (%v)",
			preauth.StatResult.Message,
			preauth.StatResult.Message_Detail,
		)), nil
	}

	switch preauth.Response.Result {
	case "allow":
		return successResp, err
	case "deny":
		return logical.ErrorResponse(preauth.Response.Status_Msg), nil
	case "enroll":
		return logical.ErrorResponse(preauth.Response.Status_Msg), nil
	case "auth":
		break
	}

	options := []func(*url.Values){authapi.AuthUsername(duoUser)}
	if method == "" {
		method = "auto"
	}
	if passcode != "" {
		method = "passcode"
		options = append(options, authapi.AuthPasscode(passcode))
	} else {
		options = append(options, authapi.AuthDevice("auto"))
	}

	result, err := duoAuthClient.Auth(method, options...)

	if err != nil {
		return logical.ErrorResponse("Could not call Duo auth"), nil
	}

	if result.StatResult.Stat != "OK" {
		return logical.ErrorResponse(fmt.Sprintf("Could not authenticate Duo user: %v (%v)",
			preauth.StatResult.Message,
			preauth.StatResult.Message_Detail,
		)), nil
	}

	if result.Response.Result != "allow" {
		return logical.ErrorResponse(result.Response.Status_Msg), nil
	}

	return successResp, err
}
