package gocb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	provider mgmtProvider
	tracer   RequestTracer
	meter    *meterWrapper
}

func (um *UserManager) tryParseErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logDebugf("Failed to read search index response body: %s", err)
		return nil
	}

	var bodyErr error
	if resp.StatusCode == 404 {
		if strings.Contains(strings.ToLower(string(b)), "unknown user") {
			bodyErr = ErrUserNotFound
		} else if strings.Contains(strings.ToLower(string(b)), "user was not found") {
			bodyErr = ErrUserNotFound
		} else if strings.Contains(strings.ToLower(string(b)), "group was not found") {
			bodyErr = ErrGroupNotFound
		} else if strings.Contains(strings.ToLower(string(b)), "unknown group") {
			bodyErr = ErrGroupNotFound
		} else {
			bodyErr = errors.New(string(b))
		}
	} else {
		bodyErr = errors.New(string(b))
	}

	return makeGenericMgmtError(bodyErr, req, resp)
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
	if opts == nil {
		opts = &GetAllUsersOptions{}
	}

	if opts.DomainName == "" {
		opts.DomainName = string(LocalDomain)
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_get_all_users", start)

	path := fmt.Sprintf("/settings/rbac/users/%s", opts.DomainName)
	span := createSpan(um.tracer, opts.ParentSpan, "manager_users_get_all_users", "management")
	span.SetAttribute("db.operation", "GET "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Method:        "GET",
		Path:          path,
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := um.provider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		usrErr := um.tryParseErrorMessage(&req, resp)
		if usrErr != nil {
			return nil, usrErr
		}
		return nil, makeMgmtBadStatusError("failed to get users", &req, resp)
	}

	var usersData []jsonUserMetadata
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&usersData)
	if err != nil {
		return nil, err
	}

	users := make([]UserAndMetadata, len(usersData))
	for userIdx, userData := range usersData {
		err := users[userIdx].fromData(userData)
		if err != nil {
			return nil, err
		}
	}

	return users, nil
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
	if opts == nil {
		opts = &GetUserOptions{}
	}

	if opts.DomainName == "" {
		opts.DomainName = string(LocalDomain)
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_get_user", start)

	path := fmt.Sprintf("/settings/rbac/users/%s/%s", opts.DomainName, name)
	span := createSpan(um.tracer, opts.ParentSpan, "manager_users_get_user", "management")
	span.SetAttribute("db.operation", "GET "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Method:        "GET",
		Path:          path,
		IsIdempotent:  true,
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := um.provider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		usrErr := um.tryParseErrorMessage(&req, resp)
		if usrErr != nil {
			return nil, usrErr
		}
		return nil, makeMgmtBadStatusError("failed to get user", &req, resp)
	}

	var userData jsonUserMetadata
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&userData)
	if err != nil {
		return nil, err
	}

	var user UserAndMetadata
	err = user.fromData(userData)
	if err != nil {
		return nil, err
	}

	return &user, nil
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
	if opts == nil {
		opts = &UpsertUserOptions{}
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_upsert_user", start)

	if opts.DomainName == "" {
		opts.DomainName = string(LocalDomain)
	}

	parseWildcard := func(str string) string {
		if str == "*" {
			return ""
		}

		return str
	}

	isNullOrWildcard := func(str string) bool {
		if str == "*" || str == "" {
			return true
		}

		return false
	}

	path := fmt.Sprintf("/settings/rbac/users/%s/%s", opts.DomainName, user.Username)
	span := createSpan(um.tracer, opts.ParentSpan, "manager_users_upsert_user", "management")
	span.SetAttribute("db.operation", "PUT "+path)
	defer span.End()

	var reqRoleStrs []string
	for _, roleData := range user.Roles {
		if roleData.Bucket == "" {
			reqRoleStrs = append(reqRoleStrs, roleData.Name)
		} else {
			scope := parseWildcard(roleData.Scope)
			collection := parseWildcard(roleData.Collection)

			if scope != "" && isNullOrWildcard(roleData.Bucket) {
				return makeInvalidArgumentsError("when a scope is specified, the bucket cannot be null or wildcard")
			}
			if collection != "" && isNullOrWildcard(scope) {
				return makeInvalidArgumentsError("when a collection is specified, the scope cannot be null or wildcard")
			}

			roleStr := fmt.Sprintf("%s[%s", roleData.Name, roleData.Bucket)
			if scope != "" {
				roleStr += ":" + roleData.Scope
			}
			if collection != "" {
				roleStr += ":" + roleData.Collection
			}
			roleStr += "]"

			reqRoleStrs = append(reqRoleStrs, roleStr)

		}
	}

	reqForm := make(url.Values)
	reqForm.Add("name", user.DisplayName)
	if user.Password != "" {
		reqForm.Add("password", user.Password)
	}
	if len(user.Groups) > 0 {
		reqForm.Add("groups", strings.Join(user.Groups, ","))
	}
	reqForm.Add("roles", strings.Join(reqRoleStrs, ","))

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Method:        "PUT",
		Path:          path,
		Body:          []byte(reqForm.Encode()),
		ContentType:   "application/x-www-form-urlencoded",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := um.provider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		usrErr := um.tryParseErrorMessage(&req, resp)
		if usrErr != nil {
			return usrErr
		}
		return makeMgmtBadStatusError("failed to upsert user", &req, resp)
	}

	return nil
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
	if opts == nil {
		opts = &DropUserOptions{}
	}

	if opts.DomainName == "" {
		opts.DomainName = string(LocalDomain)
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_drop_user", start)

	path := fmt.Sprintf("/settings/rbac/users/%s/%s", opts.DomainName, name)
	span := createSpan(um.tracer, opts.ParentSpan, "manager_users_drop_user", "management")
	span.SetAttribute("db.operation", "DELETE "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Method:        "DELETE",
		Path:          path,
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := um.provider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		usrErr := um.tryParseErrorMessage(&req, resp)
		if usrErr != nil {
			return usrErr
		}
		return makeMgmtBadStatusError("failed to drop user", &req, resp)
	}

	return nil
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
	if opts == nil {
		opts = &GetRolesOptions{}
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_get_roles", start)

	span := createSpan(um.tracer, opts.ParentSpan, "manager_users_get_roles", "management")
	span.SetAttribute("db.operation", "GET /settings/rbac/roles")
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Method:        "GET",
		Path:          "/settings/rbac/roles",
		RetryStrategy: opts.RetryStrategy,
		IsIdempotent:  true,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := um.provider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		usrErr := um.tryParseErrorMessage(&req, resp)
		if usrErr != nil {
			return nil, usrErr
		}
		return nil, makeMgmtBadStatusError("failed to get roles", &req, resp)
	}

	var roleDatas []jsonRoleDescription
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&roleDatas)
	if err != nil {
		return nil, err
	}

	roles := make([]RoleAndDescription, len(roleDatas))
	for roleIdx, roleData := range roleDatas {
		err := roles[roleIdx].fromData(roleData)
		if err != nil {
			return nil, err
		}
	}

	return roles, nil
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
	if groupName == "" {
		return nil, makeInvalidArgumentsError("groupName cannot be empty")
	}
	if opts == nil {
		opts = &GetGroupOptions{}
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_get_group", start)

	path := fmt.Sprintf("/settings/rbac/groups/%s", groupName)
	span := createSpan(um.tracer, opts.ParentSpan, "manager_users_get_group", "management")
	span.SetAttribute("db.operation", "GET "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Method:        "GET",
		Path:          path,
		RetryStrategy: opts.RetryStrategy,
		IsIdempotent:  true,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := um.provider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		usrErr := um.tryParseErrorMessage(&req, resp)
		if usrErr != nil {
			return nil, usrErr
		}
		return nil, makeMgmtBadStatusError("failed to get group", &req, resp)
	}

	var groupData jsonGroup
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&groupData)
	if err != nil {
		return nil, err
	}

	var group Group
	err = group.fromData(groupData)
	if err != nil {
		return nil, err
	}

	return &group, nil
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
	if opts == nil {
		opts = &GetAllGroupsOptions{}
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_get_all_groups", start)

	path := "/settings/rbac/groups"
	span := createSpan(um.tracer, opts.ParentSpan, "manager_users_get_all_groups", "management")
	span.SetAttribute("db.operation", "GET "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Method:        "GET",
		Path:          path,
		RetryStrategy: opts.RetryStrategy,
		IsIdempotent:  true,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := um.provider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return nil, makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		usrErr := um.tryParseErrorMessage(&req, resp)
		if usrErr != nil {
			return nil, usrErr
		}
		return nil, makeMgmtBadStatusError("failed to get all groups", &req, resp)
	}

	var groupDatas []jsonGroup
	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&groupDatas)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, len(groupDatas))
	for groupIdx, groupData := range groupDatas {
		err = groups[groupIdx].fromData(groupData)
		if err != nil {
			return nil, err
		}
	}

	return groups, nil
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
	if group.Name == "" {
		return makeInvalidArgumentsError("group name cannot be empty")
	}
	if opts == nil {
		opts = &UpsertGroupOptions{}
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_upsert_group", start)

	path := fmt.Sprintf("/settings/rbac/groups/%s", group.Name)
	span := createSpan(um.tracer, opts.ParentSpan, "manager_users_upsert_group", "management")
	span.SetAttribute("db.operation", "PUT "+path)
	defer span.End()

	var reqRoleStrs []string
	for _, roleData := range group.Roles {
		if roleData.Bucket == "" {
			reqRoleStrs = append(reqRoleStrs, roleData.Name)
		} else {
			reqRoleStrs = append(reqRoleStrs, fmt.Sprintf("%s[%s]", roleData.Name, roleData.Bucket))
		}
	}

	reqForm := make(url.Values)
	reqForm.Add("description", group.Description)
	reqForm.Add("ldap_group_ref", group.LDAPGroupReference)
	reqForm.Add("roles", strings.Join(reqRoleStrs, ","))

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Method:        "PUT",
		Path:          path,
		Body:          []byte(reqForm.Encode()),
		ContentType:   "application/x-www-form-urlencoded",
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := um.provider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		usrErr := um.tryParseErrorMessage(&req, resp)
		if usrErr != nil {
			return usrErr
		}
		return makeMgmtBadStatusError("failed to upsert group", &req, resp)
	}

	return nil
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
	if groupName == "" {
		return makeInvalidArgumentsError("groupName cannot be empty")
	}

	if opts == nil {
		opts = &DropGroupOptions{}
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_drop_group", start)

	path := fmt.Sprintf("/settings/rbac/groups/%s", groupName)
	span := createSpan(um.tracer, opts.ParentSpan, "manager_users_drop_group", "management")
	span.SetAttribute("db.operation", "DELETE "+path)
	defer span.End()

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Method:        "DELETE",
		Path:          path,
		RetryStrategy: opts.RetryStrategy,
		UniqueID:      uuid.New().String(),
		Timeout:       opts.Timeout,
		parentSpanCtx: span.Context(),
	}

	resp, err := um.provider.executeMgmtRequest(opts.Context, req)
	if err != nil {
		return makeGenericMgmtError(err, &req, resp)
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		usrErr := um.tryParseErrorMessage(&req, resp)
		if usrErr != nil {
			return usrErr
		}
		return makeMgmtBadStatusError("failed to drop group", &req, resp)
	}

	return nil
}
