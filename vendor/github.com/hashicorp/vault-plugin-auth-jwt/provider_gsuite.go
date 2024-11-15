// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jwtauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

// GSuiteProvider provides G Suite-specific configuration and behavior.
type GSuiteProvider struct {
	// Configuration for the provider
	config GSuiteProviderConfig

	// Google admin service
	adminSvc *admin.Service
}

// GSuiteProviderConfig represents the configuration for a GSuiteProvider.
type GSuiteProviderConfig struct {
	// The path to or contents of a Google service account key file. Optional.
	// If left unspecified, Application Default Credentials will be used.
	ServiceAccount string `mapstructure:"gsuite_service_account"`

	// Email address of a Google Workspace user that has access to read users
	// and groups for the organization in the Google Workspace Directory API.
	// Required if accessing the Google Workspace Directory API through
	// domain-wide delegation of authority.
	AdminImpersonateEmail string `mapstructure:"gsuite_admin_impersonate"`

	// Service account email that has been granted domain-wide delegation of
	// authority in Google Workspace. Required if accessing the Google
	// Workspace Directory API through domain-wide delegation of authority,
	// without using a service account key. The service account vault is
	// running under must be granted the `iam.serviceAccounts.signJwt`
	// permission on this service account. If AdminImpersonateEmail is
	// specifed, that Workspace user will be impersonated.
	ImpersonatePrincipal string `mapstructure:"impersonate_principal"`

	// If set to true, groups will be fetched from the Google Workspace
	// Directory API.
	FetchGroups bool `mapstructure:"fetch_groups"`

	// If set to true, user info will be fetched from the Google Workspace
	// Directory API using UserCustomSchemas.
	FetchUserInfo bool `mapstructure:"fetch_user_info"`

	// Group membership recursion max depth (0 = do not recurse).
	GroupsRecurseMaxDepth int `mapstructure:"groups_recurse_max_depth"`

	// Comma-separated list of G Suite custom schemas to fetch as claims.
	UserCustomSchemas string `mapstructure:"user_custom_schemas"`

	// The domain to get groups from. Set this if your workspace is
	// configured with more than one domain.
	Domain string `mapstructure:"domain"`
}

// Initialize initializes the GSuiteProvider by validating and creating configuration.
func (g *GSuiteProvider) Initialize(ctx context.Context, jc *jwtConfig) error {
	// Decode the provider config
	var config GSuiteProviderConfig
	if err := mapstructure.Decode(jc.ProviderConfig, &config); err != nil {
		return err
	}

	// Validate configuration
	if config.GroupsRecurseMaxDepth < 0 {
		return errors.New("'gsuite_recurse_max_depth' must be a positive integer")
	}

	// Set the requested scopes
	scopes := []string{
		admin.AdminDirectoryGroupReadonlyScope,
		admin.AdminDirectoryUserReadonlyScope,
	}
	g.config = config

	var ts oauth2.TokenSource
	switch {
	// A file path or JSON string may be provided for the service account parameter.
	// Check to see if a file exists at the given path, and if so, read its contents.
	// Otherwise, assume the service account has been provided as a JSON string.
	case config.ServiceAccount != "":
		var err error
		keyJSON := []byte(config.ServiceAccount)
		if fileExists(config.ServiceAccount) {
			keyJSON, err = ioutil.ReadFile(config.ServiceAccount)
			if err != nil {
				return err
			}
		}

		// Create the google JWT config from the service account
		jwtConfig, err := google.JWTConfigFromJSON(keyJSON, scopes...)
		if err != nil {
			return fmt.Errorf("error parsing service account JSON: %w", err)
		}

		// Set the subject to impersonate
		jwtConfig.Subject = config.AdminImpersonateEmail

		ts = jwtConfig.TokenSource(ctx)
	// We are performing impersonation of a Workspace user through domain-wide
	// delegation of authority.
	case config.ImpersonatePrincipal != "":
		its, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: config.ImpersonatePrincipal,
			Scopes:          scopes,
			Subject:         config.AdminImpersonateEmail,
		})
		if err != nil {
			return fmt.Errorf("failed to impersonate principal: %q: %w", config.ImpersonatePrincipal, err)
		}

		ts = its
	// Assume Application Default Credentials and no impersonation.
	default:
		creds, err := google.FindDefaultCredentials(ctx, scopes...)
		if err != nil {
			return fmt.Errorf("failed to find application default credentials: %w", err)
		}

		ts = creds.TokenSource
	}

	svc, err := admin.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return err
	}
	g.adminSvc = svc

	return nil
}

// SensitiveKeys returns keys that should be redacted when reading the config of this provider
func (g *GSuiteProvider) SensitiveKeys() []string {
	return []string{"gsuite_service_account"}
}

// FetchGroups fetches and returns groups from G Suite.
func (g *GSuiteProvider) FetchGroups(ctx context.Context, b *jwtAuthBackend, allClaims map[string]interface{}, role *jwtRole, _ oauth2.TokenSource) (interface{}, error) {
	if !g.config.FetchGroups {
		return nil, nil
	}

	userName, err := g.getUserClaim(b, allClaims, role)
	if err != nil {
		return nil, err
	}

	// Get the G Suite groups
	userGroupsMap := make(map[string]bool)
	if err := g.search(ctx, userGroupsMap, userName, g.config.GroupsRecurseMaxDepth); err != nil {
		return nil, err
	}

	// Convert set of groups to list
	userGroups := make([]interface{}, 0, len(userGroupsMap))
	for email := range userGroupsMap {
		userGroups = append(userGroups, email)
	}

	b.Logger().Debug("fetched G Suite groups", "groups", userGroups)
	return userGroups, nil
}

// search recursively searches for G Suite groups based on a configured depth for this provider.
func (g *GSuiteProvider) search(ctx context.Context, visited map[string]bool, userName string, depth int) error {
	req := g.adminSvc.Groups.List().UserKey(userName)

	// Request for a specific domain if one is set
	if g.config.Domain != "" {
		req = req.Domain(g.config.Domain)
	}

	req = req.Fields("nextPageToken", "groups(email)")
	if err := req.Pages(ctx, func(groups *admin.Groups) error {
		var newGroups []string
		for _, group := range groups.Groups {
			if _, ok := visited[group.Email]; ok {
				continue
			}
			visited[group.Email] = true
			newGroups = append(newGroups, group.Email)
		}
		// Only recursively search for new groups that haven't been seen
		if depth > 0 {
			for _, email := range newGroups {
				if err := g.search(ctx, visited, email, depth-1); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// FetchUserInfo fetches additional user information from G Suite using custom schemas.
func (g *GSuiteProvider) FetchUserInfo(ctx context.Context, b *jwtAuthBackend, allClaims map[string]interface{}, role *jwtRole) error {
	if !g.config.FetchUserInfo || g.config.UserCustomSchemas == "" {
		if g.config.UserCustomSchemas != "" {
			b.Logger().Warn(fmt.Sprintf("must set 'fetch_user_info=true' to fetch 'user_custom_schemas': %s", g.config.UserCustomSchemas))
		}

		return nil
	}

	userName, err := g.getUserClaim(b, allClaims, role)
	if err != nil {
		return err
	}

	return g.fillCustomSchemas(ctx, userName, allClaims)
}

// fillCustomSchemas fetches G Suite user information associated with the custom schemas
// configured for this provider. It inserts the schema -> value pairs into the passed
// allClaims so that the values can be used for claim mapping to token and identity metadata.
func (g *GSuiteProvider) fillCustomSchemas(ctx context.Context, userName string, allClaims map[string]interface{}) error {
	userResponse, err := g.adminSvc.Users.Get(userName).Context(ctx).Projection("custom").
		CustomFieldMask(g.config.UserCustomSchemas).Fields("customSchemas").Do()
	if err != nil {
		return err
	}

	for schema, rawValue := range userResponse.CustomSchemas {
		// note: metadata extraction via claim_mappings only supports strings
		// as values, but filtering happens later so we must use interface{}
		var value map[string]interface{}
		if err := json.Unmarshal(rawValue, &value); err != nil {
			return err
		}

		allClaims[schema] = value
	}

	return nil
}

// getUserClaim returns the user claim value configured in the passed role.
// If the user claim is not found or is not a string, an error is returned.
func (g *GSuiteProvider) getUserClaim(b *jwtAuthBackend, allClaims map[string]interface{}, role *jwtRole) (string, error) {
	userClaimRaw := getClaim(b.Logger(), allClaims, role.UserClaim)
	if userClaimRaw == nil {
		return "", fmt.Errorf("unable to locate %q in claims", role.UserClaim)
	}
	userClaim, ok := userClaimRaw.(string)
	if !ok {
		return "", fmt.Errorf("claim %q could not be converted to string", role.UserClaim)
	}

	return userClaim, nil
}

// fileExists returns true if a file exists at the given path.
func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi != nil && !fi.IsDir()
}
