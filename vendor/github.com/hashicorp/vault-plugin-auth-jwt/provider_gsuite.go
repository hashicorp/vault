package jwtauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

// GSuiteProvider provides G Suite-specific configuration and behavior.
type GSuiteProvider struct {
	config    GSuiteProviderConfig // Configuration for the provider
	jwtConfig *jwt.Config          // Google JWT configuration
	adminSvc  *admin.Service       // Google admin service
}

// GSuiteProviderConfig represents the configuration for a GSuiteProvider.
type GSuiteProviderConfig struct {
	// Path to a Google service account key file. Required.
	ServiceAccountFilePath string `mapstructure:"gsuite_service_account"`

	// Email address of a G Suite admin to impersonate. Required.
	AdminImpersonateEmail string `mapstructure:"gsuite_admin_impersonate"`

	// If set to true, groups will be fetched from G Suite.
	FetchGroups bool `mapstructure:"fetch_groups"`

	// If set to true, user info will be fetched from G Suite using UserCustomSchemas.
	FetchUserInfo bool `mapstructure:"fetch_user_info"`

	// Group membership recursion max depth (0 = do not recurse).
	GroupsRecurseMaxDepth int `mapstructure:"groups_recurse_max_depth"`

	// Comma-separated list of G Suite custom schemas to fetch as claims.
	UserCustomSchemas string `mapstructure:"user_custom_schemas"`

	// JSON contents of a Google service account key file.
	serviceAccountKeyJSON []byte
}

// Initialize initializes the GSuiteProvider by validating and creating configuration.
func (g *GSuiteProvider) Initialize(jc *jwtConfig) error {
	// Decode the provider config
	var config GSuiteProviderConfig
	if err := mapstructure.Decode(jc.ProviderConfig, &config); err != nil {
		return err
	}

	// Read the Google service account key file
	keyJSON, err := ioutil.ReadFile(config.ServiceAccountFilePath)
	if err != nil {
		return err
	}
	config.serviceAccountKeyJSON = keyJSON

	return g.initialize(config)
}

func (g *GSuiteProvider) initialize(config GSuiteProviderConfig) error {
	var err error

	// Validate configuration
	if config.ServiceAccountFilePath == "" {
		return errors.New("'gsuite_service_account' must be set to the file path for a " +
			"service account key")
	}
	if config.AdminImpersonateEmail == "" {
		return errors.New("'gsuite_admin_impersonate' must be set to an email address of a " +
			"G Suite user with 'Read' permission to access the G Suite Admin User and Group APIs")
	}
	if config.GroupsRecurseMaxDepth < 0 {
		return errors.New("'gsuite_recurse_max_depth' must be a positive integer")
	}

	// Create the google JWT config from the service account key file
	if g.jwtConfig, err = google.JWTConfigFromJSON(config.serviceAccountKeyJSON,
		admin.AdminDirectoryGroupReadonlyScope, admin.AdminDirectoryUserReadonlyScope); err != nil {
		return err
	}

	// Set the subject to impersonate and config
	g.jwtConfig.Subject = config.AdminImpersonateEmail
	g.config = config
	return nil
}

// SensitiveKeys returns keys that should be redacted when reading the config of this provider
func (g *GSuiteProvider) SensitiveKeys() []string {
	return []string{}
}

// FetchGroups fetches and returns groups from G Suite.
func (g *GSuiteProvider) FetchGroups(b *jwtAuthBackend, allClaims map[string]interface{}, role *jwtRole) (interface{}, error) {
	if !g.config.FetchGroups {
		return nil, nil
	}

	userName, err := g.getUserClaim(b, allClaims, role)
	if err != nil {
		return nil, err
	}

	// Set context and create a new admin service for requests to Google admin APIs
	g.adminSvc, err = admin.NewService(b.providerCtx, option.WithHTTPClient(g.jwtConfig.Client(b.providerCtx)))
	if err != nil {
		return nil, err
	}

	// Get the G Suite groups
	userGroupsMap := make(map[string]bool)
	if err := g.search(b.providerCtx, userGroupsMap, userName, g.config.GroupsRecurseMaxDepth); err != nil {
		return nil, err
	}

	// Convert set of groups to list
	var userGroups = make([]interface{}, 0, len(userGroupsMap))
	for email := range userGroupsMap {
		userGroups = append(userGroups, email)
	}

	b.Logger().Debug("fetched G Suite groups", "groups", userGroups)
	return userGroups, nil
}

// search recursively searches for G Suite groups based on a configured depth for this provider.
func (g *GSuiteProvider) search(ctx context.Context, visited map[string]bool, userName string, depth int) error {
	call := g.adminSvc.Groups.List().UserKey(userName).Fields("nextPageToken", "groups(email)")
	if err := call.Pages(ctx, func(groups *admin.Groups) error {
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
func (g *GSuiteProvider) FetchUserInfo(b *jwtAuthBackend, allClaims map[string]interface{}, role *jwtRole) error {
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

	// Set context and create a new admin service for requests to Google admin APIs
	g.adminSvc, err = admin.NewService(b.providerCtx, option.WithHTTPClient(g.jwtConfig.Client(b.providerCtx)))
	if err != nil {
		return err
	}

	return g.fillCustomSchemas(b.providerCtx, userName, allClaims)
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
