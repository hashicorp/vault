package gocb

import (
	"context"
	"time"
)

// AuthDomain specifies the user domain of a specific user
type AuthDomain string

const (
	// LocalDomain specifies users that are locally stored in Couchbase.
	LocalDomain AuthDomain = "local"

	// ExternalDomain specifies users that are externally stored
	// (in LDAP for instance).
	ExternalDomain AuthDomain = "external"
)

type jsonOrigin struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type jsonRole struct {
	RoleName       string `json:"role"`
	BucketName     string `json:"bucket_name"`
	ScopeName      string `json:"scope_name"`
	CollectionName string `json:"collection_name"`
}

type jsonRoleDescription struct {
	jsonRole

	Name        string `json:"name"`
	Description string `json:"desc"`
}

type jsonRoleOrigins struct {
	jsonRole

	Origins []jsonOrigin
}

type jsonUserMetadata struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Roles           []jsonRoleOrigins `json:"roles"`
	Groups          []string          `json:"groups"`
	Domain          AuthDomain        `json:"domain"`
	ExternalGroups  []string          `json:"external_groups"`
	PasswordChanged time.Time         `json:"password_change_date"`
}

type jsonGroup struct {
	Name               string     `json:"id"`
	Description        string     `json:"description"`
	Roles              []jsonRole `json:"roles"`
	LDAPGroupReference string     `json:"ldap_group_ref"`
}

// Role represents a specific permission.
type Role struct {
	Name       string `json:"role"`
	Bucket     string `json:"bucket_name"`
	Scope      string `json:"scope_name"`
	Collection string `json:"collection_name"`
}

func (ro *Role) fromData(data jsonRole) error {
	ro.Name = data.RoleName
	ro.Bucket = data.BucketName
	ro.Scope = data.ScopeName
	ro.Collection = data.CollectionName

	if ro.Scope == "*" {
		ro.Scope = ""
	}
	if ro.Collection == "*" {
		ro.Collection = ""
	}

	return nil
}

// RoleAndDescription represents a role with its display name and description.
type RoleAndDescription struct {
	Role

	DisplayName string
	Description string
}

func (rd *RoleAndDescription) fromData(data jsonRoleDescription) error {
	err := rd.Role.fromData(data.jsonRole)
	if err != nil {
		return err
	}

	rd.DisplayName = data.Name
	rd.Description = data.Description

	return nil
}

// Origin indicates why a user has a specific role. Is the Origin Type is "user" then the role is assigned
// directly to the user. If the type is "group" then it means that the role has been inherited from the group
// identified by the Name field.
type Origin struct {
	Type string
	Name string
}

func (o *Origin) fromData(data jsonOrigin) error {
	o.Type = data.Type
	o.Name = data.Name

	return nil
}

// RoleAndOrigins associates a role with its origins.
type RoleAndOrigins struct {
	Role

	Origins []Origin
}

func (ro *RoleAndOrigins) fromData(data jsonRoleOrigins) error {
	err := ro.Role.fromData(data.jsonRole)
	if err != nil {
		return err
	}

	origins := make([]Origin, len(data.Origins))
	for i, originData := range data.Origins {
		var origin Origin
		err := origin.fromData(originData)
		if err != nil {
			return err
		}

		origins[i] = origin
	}
	ro.Origins = origins

	return nil
}

// User represents a user which was retrieved from the server.
type User struct {
	Username    string
	DisplayName string
	// Roles are the roles assigned to the user that are of type "user".
	Roles    []Role
	Groups   []string
	Password string
}

// UserAndMetadata represents a user and user meta-data from the server.
type UserAndMetadata struct {
	User

	Domain AuthDomain
	// EffectiveRoles are all of the user's roles and the origins.
	EffectiveRoles  []RoleAndOrigins
	ExternalGroups  []string
	PasswordChanged time.Time
}

func (um *UserAndMetadata) fromData(data jsonUserMetadata) error {
	um.User.Username = data.ID
	um.User.DisplayName = data.Name
	um.User.Groups = data.Groups

	um.ExternalGroups = data.ExternalGroups
	um.Domain = data.Domain
	um.PasswordChanged = data.PasswordChanged

	var roles []Role
	var effectiveRoles []RoleAndOrigins
	for _, roleData := range data.Roles {
		var effectiveRole RoleAndOrigins
		err := effectiveRole.fromData(roleData)
		if err != nil {
			return err
		}

		effectiveRoles = append(effectiveRoles, effectiveRole)

		role := effectiveRole.Role
		if roleData.Origins == nil {
			roles = append(roles, role)
		} else {
			for _, origin := range effectiveRole.Origins {
				if origin.Type == "user" {
					roles = append(roles, role)
					break
				}
			}
		}
	}
	um.EffectiveRoles = effectiveRoles
	um.User.Roles = roles

	return nil
}

// Group represents a user group on the server.
type Group struct {
	Name               string
	Description        string
	Roles              []Role
	LDAPGroupReference string
}

func (g *Group) fromData(data jsonGroup) error {
	g.Name = data.Name
	g.Description = data.Description
	g.LDAPGroupReference = data.LDAPGroupReference

	roles := make([]Role, len(data.Roles))
	for roleIdx, roleData := range data.Roles {
		err := roles[roleIdx].fromData(roleData)
		if err != nil {
			return err
		}
	}
	g.Roles = roles

	return nil
}

// UserManager provides methods for performing Couchbase user management.
type UserManager struct {
	controller *providerController[userManagerProvider]
}

// GetAllUsersOptions is the set of options available to the user manager GetAll operation.
type GetAllUsersOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy

	DomainName string
	ParentSpan RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllUsers returns a list of all the users from the cluster.
func (um *UserManager) GetAllUsers(opts *GetAllUsersOptions) ([]UserAndMetadata, error) {
	return autoOpControl(um.controller, func(provider userManagerProvider) ([]UserAndMetadata, error) {
		if opts == nil {
			opts = &GetAllUsersOptions{}
		}

		return provider.GetAllUsers(opts)
	})
}

// GetUserOptions is the set of options available to the user manager Get operation.
type GetUserOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy

	DomainName string
	ParentSpan RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetUser returns the data for a particular user
func (um *UserManager) GetUser(name string, opts *GetUserOptions) (*UserAndMetadata, error) {
	return autoOpControl(um.controller, func(provider userManagerProvider) (*UserAndMetadata, error) {
		if opts == nil {
			opts = &GetUserOptions{}
		}

		return provider.GetUser(name, opts)
	})
}

// UpsertUserOptions is the set of options available to the user manager Upsert operation.
type UpsertUserOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy

	DomainName string
	ParentSpan RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// UpsertUser updates a built-in RBAC user on the cluster.
func (um *UserManager) UpsertUser(user User, opts *UpsertUserOptions) error {
	return autoOpControlErrorOnly(um.controller, func(provider userManagerProvider) error {
		if opts == nil {
			opts = &UpsertUserOptions{}
		}

		return provider.UpsertUser(user, opts)
	})
}

// DropUserOptions is the set of options available to the user manager Drop operation.
type DropUserOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy

	DomainName string
	ParentSpan RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropUser removes a built-in RBAC user on the cluster.
func (um *UserManager) DropUser(name string, opts *DropUserOptions) error {
	return autoOpControlErrorOnly(um.controller, func(provider userManagerProvider) error {
		if opts == nil {
			opts = &DropUserOptions{}
		}

		return provider.DropUser(name, opts)
	})
}

// GetRolesOptions is the set of options available to the user manager GetRoles operation.
type GetRolesOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetRoles lists the roles supported by the cluster.
func (um *UserManager) GetRoles(opts *GetRolesOptions) ([]RoleAndDescription, error) {
	return autoOpControl(um.controller, func(provider userManagerProvider) ([]RoleAndDescription, error) {
		if opts == nil {
			opts = &GetRolesOptions{}
		}

		return provider.GetRoles(opts)
	})
}

// GetGroupOptions is the set of options available to the group manager Get operation.
type GetGroupOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetGroup fetches a single group from the server.
func (um *UserManager) GetGroup(groupName string, opts *GetGroupOptions) (*Group, error) {
	return autoOpControl(um.controller, func(provider userManagerProvider) (*Group, error) {
		if groupName == "" {
			return nil, makeInvalidArgumentsError("groupName cannot be empty")
		}
		if opts == nil {
			opts = &GetGroupOptions{}
		}

		return provider.GetGroup(groupName, opts)
	})
}

// GetAllGroupsOptions is the set of options available to the group manager GetAll operation.
type GetAllGroupsOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetAllGroups fetches all groups from the server.
func (um *UserManager) GetAllGroups(opts *GetAllGroupsOptions) ([]Group, error) {
	return autoOpControl(um.controller, func(provider userManagerProvider) ([]Group, error) {
		if opts == nil {
			opts = &GetAllGroupsOptions{}
		}

		return provider.GetAllGroups(opts)
	})
}

// UpsertGroupOptions is the set of options available to the group manager Upsert operation.
type UpsertGroupOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// UpsertGroup creates, or updates, a group on the server.
func (um *UserManager) UpsertGroup(group Group, opts *UpsertGroupOptions) error {
	return autoOpControlErrorOnly(um.controller, func(provider userManagerProvider) error {
		if group.Name == "" {
			return makeInvalidArgumentsError("group name cannot be empty")
		}
		if opts == nil {
			opts = &UpsertGroupOptions{}
		}

		return provider.UpsertGroup(group, opts)
	})
}

// DropGroupOptions is the set of options available to the group manager Drop operation.
type DropGroupOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// DropGroup removes a group from the server.
func (um *UserManager) DropGroup(groupName string, opts *DropGroupOptions) error {
	return autoOpControlErrorOnly(um.controller, func(provider userManagerProvider) error {
		if groupName == "" {
			return makeInvalidArgumentsError("groupName cannot be empty")
		}

		if opts == nil {
			opts = &DropGroupOptions{}
		}

		return provider.DropGroup(groupName, opts)
	})
}

// ChangePasswordOptions is the set of options available to the user manager ChangePassword operation.
type ChangePasswordOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// ChangePassword changes the password for the currently authenticated user.
// *Note*: Usage of this function will effectively invalidate the SDK instance and further requests will fail
// due to authentication errors. After using this function the SDK must be reinitialized.
func (um *UserManager) ChangePassword(newPassword string, opts *ChangePasswordOptions) error {
	return autoOpControlErrorOnly(um.controller, func(provider userManagerProvider) error {
		if newPassword == "" {
			return makeInvalidArgumentsError("new password cannot be empty")
		}

		if opts == nil {
			opts = &ChangePasswordOptions{}
		}

		return provider.ChangePassword(newPassword, opts)
	})
}
