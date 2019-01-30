package jwtauth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRoleList(b *jwtAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "role/?",
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},
		HelpSynopsis:    strings.TrimSpace(roleHelp["role-list"][0]),
		HelpDescription: strings.TrimSpace(roleHelp["role-list"][1]),
	}
}

// pathRole returns the path configurations for the CRUD operations on roles
func pathRole(b *jwtAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the role.",
			},
			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: "List of policies on the role.",
			},
			"num_uses": {
				Type:        framework.TypeInt,
				Description: `Number of times issued tokens can be used`,
			},
			"ttl": {
				Type: framework.TypeDurationSecond,
				Description: `Duration in seconds after which the issued token should expire. Defaults
to 0, in which case the value will fall back to the system/mount defaults.`,
			},
			"max_ttl": {
				Type: framework.TypeDurationSecond,
				Description: `Duration in seconds after which the issued token should not be allowed to
be renewed. Defaults to 0, in which case the value will fall back to the system/mount defaults.`,
			},
			"period": {
				Type: framework.TypeDurationSecond,
				Description: `If set, indicates that the token generated using this role
should never expire. The token should be renewed within the
duration specified by this value. At each renewal, the token's
TTL will be set to the value of this parameter.`,
			},
			"bound_subject": {
				Type:        framework.TypeString,
				Description: `The 'sub' claim that is valid for login. Optional.`,
			},
			"bound_audiences": {
				Type:        framework.TypeCommaStringSlice,
				Description: `Comma-separated list of 'aud' claims that are valid for login; any match is sufficient`,
			},
			"user_claim": {
				Type:        framework.TypeString,
				Description: `The claim to use for the Identity entity alias name`,
			},
			"groups_claim": {
				Type:        framework.TypeString,
				Description: `The claim to use for the Identity group alias names`,
			},
			"groups_claim_delimiter_pattern": {
				Type:        framework.TypeString,
				Description: `A pattern of delimiters used to allow the groups_claim to live outside of the top-level JWT structure. For instance, a "groups_claim" of "meta/user.name/groups" with this field set to "//" will expect nested structures named "meta", "user.name", and "groups". If this field was set to "/./" the groups information would expect to be via nested structures of "meta", "user", "name", and "groups".`,
			},
			"bound_cidrs": {
				Type: framework.TypeCommaStringSlice,
				Description: `Comma-separated list of IP CIDRS that are allowed to 
authenticate against this role`,
			},
		},
		ExistenceCheck: b.pathRoleExistenceCheck,
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathRoleCreateUpdate,
			logical.UpdateOperation: b.pathRoleCreateUpdate,
			logical.ReadOperation:   b.pathRoleRead,
			logical.DeleteOperation: b.pathRoleDelete,
		},
		HelpSynopsis:    strings.TrimSpace(roleHelp["role"][0]),
		HelpDescription: strings.TrimSpace(roleHelp["role"][1]),
	}
}

type jwtRole struct {
	// Policies that are to be required by the token to access this role
	Policies []string `json:"policies"`

	// TokenNumUses defines the number of allowed uses of the token issued
	NumUses int `json:"num_uses"`

	// Duration before which an issued token must be renewed
	TTL time.Duration `json:"ttl"`

	// Duration after which an issued token should not be allowed to be renewed
	MaxTTL time.Duration `json:"max_ttl"`

	// Period, if set, indicates that the token generated using this role
	// should never expire. The token should be renewed within the duration
	// specified by this value. The renewal duration will be fixed if the
	// value is not modified on the role. If the `Period` in the role is modified,
	// a token will pick up the new value during its next renewal.
	Period time.Duration `json:"period"`

	// Role binding properties
	BoundAudiences              []string                      `json:"bound_audiences"`
	BoundSubject                string                        `json:"bound_subject"`
	BoundCIDRs                  []*sockaddr.SockAddrMarshaler `json:"bound_cidrs"`
	UserClaim                   string                        `json:"user_claim"`
	GroupsClaim                 string                        `json:"groups_claim"`
	GroupsClaimDelimiterPattern string                        `json:"groups_claim_delimiter_pattern"`
}

// role takes a storage backend and the name and returns the role's storage
// entry
func (b *jwtAuthBackend) role(ctx context.Context, s logical.Storage, name string) (*jwtRole, error) {
	raw, err := s.Get(ctx, rolePrefix+name)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	role := new(jwtRole)
	if err := raw.DecodeJSON(role); err != nil {
		return nil, err
	}

	return role, nil
}

// pathRoleExistenceCheck returns whether the role with the given name exists or not.
func (b *jwtAuthBackend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	role, err := b.role(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}
	return role != nil, nil
}

// pathRoleList is used to list all the Roles registered with the backend.
func (b *jwtAuthBackend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, rolePrefix)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

// pathRoleRead grabs a read lock and reads the options set on the role from the storage
func (b *jwtAuthBackend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// Create a map of data to be returned
	resp := &logical.Response{
		Data: map[string]interface{}{
			"policies":                       role.Policies,
			"num_uses":                       role.NumUses,
			"period":                         int64(role.Period.Seconds()),
			"ttl":                            int64(role.TTL.Seconds()),
			"max_ttl":                        int64(role.MaxTTL.Seconds()),
			"bound_audiences":                role.BoundAudiences,
			"bound_subject":                  role.BoundSubject,
			"bound_cidrs":                    role.BoundCIDRs,
			"user_claim":                     role.UserClaim,
			"groups_claim":                   role.GroupsClaim,
			"groups_claim_delimiter_pattern": role.GroupsClaimDelimiterPattern,
		},
	}

	return resp, nil
}

// pathRoleDelete removes the role from storage
func (b *jwtAuthBackend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("role name required"), nil
	}

	// Delete the role itself
	if err := req.Storage.Delete(ctx, rolePrefix+roleName); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathRoleCreateUpdate registers a new role with the backend or updates the options
// of an existing role
func (b *jwtAuthBackend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	// Check if the role already exists
	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	// Create a new entry object if this is a CreateOperation
	if role == nil {
		if req.Operation == logical.UpdateOperation {
			return nil, errors.New("role entry not found during update operation")
		}
		role = new(jwtRole)
	}

	if policiesRaw, ok := data.GetOk("policies"); ok {
		role.Policies = policyutil.ParsePolicies(policiesRaw)
	}

	periodRaw, ok := data.GetOk("period")
	if ok {
		role.Period = time.Duration(periodRaw.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		role.Period = time.Duration(data.Get("period").(int)) * time.Second
	}
	if role.Period > b.System().MaxLeaseTTL() {
		return logical.ErrorResponse(fmt.Sprintf("'period' of '%q' is greater than the backend's maximum lease TTL of '%q'", role.Period.String(), b.System().MaxLeaseTTL().String())), nil
	}

	if tokenNumUsesRaw, ok := data.GetOk("num_uses"); ok {
		role.NumUses = tokenNumUsesRaw.(int)
	} else if req.Operation == logical.CreateOperation {
		role.NumUses = data.Get("num_uses").(int)
	}
	if role.NumUses < 0 {
		return logical.ErrorResponse("num_uses cannot be negative"), nil
	}

	if tokenTTLRaw, ok := data.GetOk("ttl"); ok {
		role.TTL = time.Duration(tokenTTLRaw.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		role.TTL = time.Duration(data.Get("ttl").(int)) * time.Second
	}

	if tokenMaxTTLRaw, ok := data.GetOk("max_ttl"); ok {
		role.MaxTTL = time.Duration(tokenMaxTTLRaw.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		role.MaxTTL = time.Duration(data.Get("max_ttl").(int)) * time.Second
	}

	if boundAudiences, ok := data.GetOk("bound_audiences"); ok {
		role.BoundAudiences = boundAudiences.([]string)
	}

	if boundSubject, ok := data.GetOk("bound_subject"); ok {
		role.BoundSubject = boundSubject.(string)
	}

	if boundCIDRs, ok := data.GetOk("bound_cidrs"); ok {
		parsedCIDRs, err := parseutil.ParseAddrs(boundCIDRs)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
		role.BoundCIDRs = parsedCIDRs
	}

	if userClaim, ok := data.GetOk("user_claim"); ok {
		role.UserClaim = userClaim.(string)
	}
	if role.UserClaim == "" {
		return logical.ErrorResponse("a user claim must be defined on the role"), nil
	}

	if groupsClaim, ok := data.GetOk("groups_claim"); ok {
		role.GroupsClaim = groupsClaim.(string)
	}

	if groupsClaimDelimiterPattern, ok := data.GetOk("groups_claim_delimiter_pattern"); ok {
		role.GroupsClaimDelimiterPattern = groupsClaimDelimiterPattern.(string)
	}

	// Validate claim/delims
	if role.GroupsClaim != "" {
		if _, err := parseClaimWithDelimiters(role.GroupsClaim, role.GroupsClaimDelimiterPattern); err != nil {
			return logical.ErrorResponse(errwrap.Wrapf("error validating delimiters for groups claim: {{err}}", err).Error()), nil
		}
	}

	if len(role.BoundAudiences) == 0 &&
		len(role.BoundCIDRs) == 0 &&
		role.BoundSubject == "" {
		return logical.ErrorResponse("must have at least one bound constraint when creating/updating a role"), nil
	}

	// Check that the TTL value provided is less than the MaxTTL.
	// Sanitizing the TTL and MaxTTL is not required now and can be performed
	// at credential issue time.
	if role.MaxTTL > 0 && role.TTL > role.MaxTTL {
		return logical.ErrorResponse("ttl should not be greater than max_ttl"), nil
	}

	var resp *logical.Response
	if role.MaxTTL > b.System().MaxLeaseTTL() {
		resp = &logical.Response{}
		resp.AddWarning("max_ttl is greater than the system or backend mount's maximum TTL value; issued tokens' max TTL value will be truncated")
	}

	// Store the entry.
	entry, err := logical.StorageEntryJSON(rolePrefix+roleName, role)
	if err != nil {
		return nil, err
	}
	if err = req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return resp, nil
}

// parseClaimWithDelimiters parses a given claim string and ensures that we can
// separate it out into a "map path"
func parseClaimWithDelimiters(claim, delimiters string) ([]string, error) {
	if delimiters == "" {
		return []string{claim}, nil
	}
	var ret []string
	for _, runeVal := range delimiters {
		idx := strings.IndexRune(claim, runeVal)
		switch idx {
		case -1:
			return nil, fmt.Errorf("could not find instance of %q delimiter in claim", string(runeVal))
		case 0:
			return nil, fmt.Errorf("instance of %q delimiter in claim is at beginning of claim string", string(runeVal))
		case len(claim) - 1:
			return nil, fmt.Errorf("instance of %q delimiter in claim is at end of claim string", string(runeVal))
		default:
			ret = append(ret, claim[:idx])
			claim = claim[idx+1:]
		}
	}
	ret = append(ret, claim)

	return ret, nil
}

// roleStorageEntry stores all the options that are set on an role
var roleHelp = map[string][2]string{
	"role-list": {
		"Lists all the roles registered with the backend.",
		"The list will contain the names of the roles.",
	},
	"role": {
		"Register an role with the backend.",
		`A role is required to authenticate with this backend. The role binds
		JWT token information with token policies and settings.
		The bindings, token polices and token settings can all be configured
		using this endpoint`,
	},
}
