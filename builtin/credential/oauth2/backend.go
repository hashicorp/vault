package oauth2

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	token, err := oauthConfig.PasswordCredentialsToken(goauth2.NoContext, username, password)
	if err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("Oauth2 auth failed: %v", err)), nil
	}
	if token == nil {
		return nil, logical.ErrorResponse("Oauth2 auth backend unexpected failure"), nil
	}

	// Get groups assigned by oauth provider from UserInfoURL
	oauthGroups, err := b.getOauthGroups(cfg, oauthConfig, token)
	if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
	}
	if b.Logger().IsDebug() {
		b.Logger().Debug("auth/oauth2: Groups fetched", "num_groups", len(oauthGroups),
			"groups", oauthGroups)
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
	// Get locally-assigned group memberships
	user, err := b.User(req.Storage, username)
	if err == nil && user != nil && user.Groups != nil {
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/oauth2: adding local groups", "num_local_groups",
				len(user.Groups), "local_groups", user.Groups)
		}
		allGroups = append(allGroups, user.Groups...)
	}
	// Merge local and oauth groups
	allGroups = append(allGroups, oauthGroups...)

	// Retrieve policies for given groups
	var policies []string
	for _, groupName := range allGroups {
		group, err := b.Group(req.Storage, groupName)
		if err == nil && group != nil && group.Policies != nil {
			policies = append(policies, group.Policies...)
		}
	}

	// Merge individually-assigned user policies into group policies
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

func (b *backend) getOauthGroups(cfg *ConfigEntry, oauthConfig *goauth2.Config, token *goauth2.Token) ([]string, error) {
	if len(cfg.UserInfoURL) != 0 {
		client := oauthConfig.Client(goauth2.NoContext, token)
		res, err := client.Get(cfg.UserInfoURL)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		var parsed map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&parsed)
		if err != nil {
			return nil, err
		}

		// Allow groups to be a JSON array or a CSV list
		switch parsedGroups := parsed[cfg.UserInfoGroupKey].(type) {
		case []interface{}:
			groups := make([]string, len(parsedGroups))
			for i, group := range parsedGroups {
				groups[i] = group.(string)
			}
			return groups, nil
		case string:
			groups := strings.Split(parsedGroups, ",")
			for i, group := range groups {
				groups[i] = strings.TrimSpace(group)
			}
			return groups, nil
		default:
			return nil, errors.New("Failed to parse groups")
		}
	}
	return []string{}, nil
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
