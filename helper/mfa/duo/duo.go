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
	duo_config, err := GetDuoConfig(req)
	if err != nil || duo_config == nil {
		return logical.ErrorResponse("Could not load Duo configuration"), nil
	}

	duo_auth_client, err := GetDuoAuthClient(req, duo_config)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	username, ok := resp.Auth.Metadata["username"]
	if !ok {
		return logical.ErrorResponse("Could not read username for MFA"), nil
	}

	duo_user := fmt.Sprintf(duo_config.UsernameFormat, username)

	preauth, err := duo_auth_client.Preauth(
		authapi.PreauthUsername(duo_user),
		authapi.PreauthIpAddr(req.Connection.RemoteAddr),
	)

	if preauth.StatResult.Stat != "OK" {
		return logical.ErrorResponse(fmt.Sprintf("Could not look up Duo user information: %s (%s)",
			*preauth.StatResult.Message,
			*preauth.StatResult.Message_Detail,
		)), nil
	}

	switch preauth.Response.Result {
	case "allow":
		return resp, err
	case "deny":
		return logical.ErrorResponse(preauth.Response.Status_Msg), nil
	case "enroll":
		return logical.ErrorResponse(preauth.Response.Status_Msg), nil
	case "auth":
		break
	}

	options := []func(*url.Values){authapi.AuthUsername(duo_user)}

	method := d.Get("method").(string)
	if method == "" {
		method = "auto"
	}

	passcode := d.Get("passcode").(string)
	if passcode != "" {
		method = "passcode"
		options = append(options, authapi.AuthPasscode(passcode))
	} else {
		options = append(options, authapi.AuthDevice("auto"))
	}

	result, err := duo_auth_client.Auth(method, options...)

	if err != nil {
		return logical.ErrorResponse("Could not call Duo auth"), nil
	}

	if result.StatResult.Stat != "OK" {
		return logical.ErrorResponse(fmt.Sprintf("Could not authenticate Duo user: %s (%s)",
			*preauth.StatResult.Message,
			*preauth.StatResult.Message_Detail,
		)), nil
	}

	if result.Response.Result != "allow" {
		return logical.ErrorResponse(result.Response.Status_Msg), nil
	}

	return resp, err
}
