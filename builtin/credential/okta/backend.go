package okta

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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
		return nil, logical.ErrorResponse("Okta backend not configured"), nil
	}

	client := cfg.OktaClient()
	auth, err := client.Authenticate(username, password)
	if err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("Okta auth failed: %v", err)), nil
	}
	if auth == nil {
		return nil, logical.ErrorResponse("okta auth backend unexpected failure"), nil
	}

	oktaGroups, err := b.getOktaGroups(cfg, auth.Embedded.User.ID)
	if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
	}
	if b.Logger().IsDebug() {
		b.Logger().Debug("auth/okta: Groups fetched from Okta", "num_groups", len(oktaGroups), "groups", oktaGroups)
	}

	oktaResponse := &logical.Response{
		Data: map[string]interface{}{},
	}
	if len(oktaGroups) == 0 {
		errString := fmt.Sprintf(
			"no Okta groups found; only policies from locally-defined groups available")
		oktaResponse.AddWarning(errString)
	}

	var allGroups []string
	// Import the custom added groups from okta backend
	user, err := b.User(req.Storage, username)
	if err == nil && user != nil && user.Groups != nil {
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: adding local groups", "num_local_groups", len(user.Groups), "local_groups", user.Groups)
		}
		allGroups = append(allGroups, user.Groups...)
	}
	// Merge local and Okta groups
	allGroups = append(allGroups, oktaGroups...)

	// Retrieve policies
	var policies []string
	for _, groupName := range allGroups {
		group, err := b.Group(req.Storage, groupName)
		if err == nil && group != nil && group.Policies != nil {
			policies = append(policies, group.Policies...)
		}
	}

	// Merge local Policies into Okta Policies
	if user != nil && user.Policies != nil {
		policies = append(policies, user.Policies...)
	}

	if len(policies) == 0 {
		errStr := "user is not a member of any authorized policy"
		if len(oktaResponse.Warnings) > 0 {
			errStr = fmt.Sprintf("%s; additionally, %s", errStr, oktaResponse.Warnings[0])
		}

		oktaResponse.Data["error"] = errStr
		return nil, oktaResponse, nil
	}

	return policies, oktaResponse, nil
}

func (b *backend) getOktaGroups(cfg *ConfigEntry, userID string) ([]string, error) {
	if cfg.Token != "" {
		client := cfg.OktaClient()
		groups, err := client.Groups(userID)
		if err != nil {
			return nil, err
		}

		oktaGroups := make([]string, 0, len(*groups))
		for _, group := range *groups {
			oktaGroups = append(oktaGroups, group.Profile.Name)
		}
		return oktaGroups, err
	}
	return nil, nil
}

const backendHelp = `
The Okta credential provider allows authentication querying,
checking username and password, and associating policies.  If an api token is configure
groups are pulled down from Okta.

Configuration of the connection is done through the "config" and "policies"
endpoints by a user with root access. Authentication is then done
by suppying the two fields for "login".
`
