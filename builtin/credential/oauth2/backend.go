package oauth2

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	goauth2 "golang.org/x/oauth2"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login/*",
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(&b),
			pathUsers(&b),
			pathGroups(&b),
			pathUsersList(&b),
			pathGroupsList(&b),
			pathLogin(&b),
		}),

		AuthRenew: b.pathLoginRenew,
	}

	return &b
}

type backend struct {
	*framework.Backend
}

func (b *backend) Login(req *logical.Request, username string, password string) ([]string, *logical.Response, error) {
	cfg, err := b.Config(req.Storage)
	if err != nil {
		return nil, nil, err
	}
	if cfg == nil {
		return nil, logical.ErrorResponse("Oauth2 backend not configured"), nil
	}

	oauthConfig := cfg.OauthConfig()
	auth, err := oauthConfig.PasswordCredentialsToken(goauth2.NoContext, username, password)
	if err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("Oauth2 auth failed: %v", err)), nil
	}
	if auth == nil {
		return nil, logical.ErrorResponse("Oauth2 auth backend unexpected failure"), nil
	}

	oauthGroups, err := b.getOauthGroups(cfg, auth)
	if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
	}
	if b.Logger().IsDebug() {
		b.Logger().Debug("auth/oauth2: Groups fetched", "num_groups", len(oauthGroups), "groups", oauthGroups)
	}

	oauthResponse := &logical.Response{
		Data: map[string]interface{}{},
	}
	if len(oauthGroups) == 0 {
		errString := fmt.Sprintf(
			"no groups found; only policies from locally-defined groups available")
		oauthResponse.AddWarning(errString)
	}

	var allGroups []string
	// Import the custom added groups from oauth backend
	user, err := b.User(req.Storage, username)
	if err == nil && user != nil && user.Groups != nil {
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/oauth2: adding local groups", "num_local_groups", len(user.Groups), "local_groups", user.Groups)
		}
		allGroups = append(allGroups, user.Groups...)
	}
	// Merge local and oauth groups
	allGroups = append(allGroups, oauthGroups...)

	// Retrieve policies
	var policies []string
	for _, groupName := range allGroups {
		group, err := b.Group(req.Storage, groupName)
		if err == nil && group != nil && group.Policies != nil {
			policies = append(policies, group.Policies...)
		}
	}

	// Merge local policies into oauth policies
	if user != nil && user.Policies != nil {
		policies = append(policies, user.Policies...)
	}

	if len(policies) == 0 {
		errStr := "user is not a member of any authorized policy"
		if len(oauthResponse.Warnings()) > 0 {
			errStr = fmt.Sprintf("%s; additionally, %s", errStr, oauthResponse.Warnings()[0])
		}

		oauthResponse.Data["error"] = errStr
		return nil, oauthResponse, nil
	}

	return policies, oauthResponse, nil
}

func (b *backend) getOauthGroups(cfg *ConfigEntry, token *goauth2.Token) ([]string, error) {
	/*  FIXME
	if cfg.Token != "" {
		client := cfg.OauthConfig()
		groups, err := client.Groups(userID)
		if err != nil {
			return nil, err
		}

		oauthGroups := make([]string, 0, len(*groups))
		for _, group := range *groups {
			oauthGroups = append(oauthGroups, group.Profile.Name)
		}
		return oauthGroups, err
	}
	return nil, nil
	*/
	oauthGroups := make([]string, 0, 1)
	oauthGroups = append(oauthGroups, "testOauthGroup")
	return oauthGroups, nil
}

const backendHelp = `
The Oauth2 backend allows for authenticating users using the 'Resource Owner
Password Flow'.  It associates policies using data provided along side the
returned bearer token or by querying a user info endpoint after successful
authentication.

Configuration of the connection is done through the "config" and "policies"
endpoints by a user with root access. Authentication is then done
by suppying the two fields for "login".
`
