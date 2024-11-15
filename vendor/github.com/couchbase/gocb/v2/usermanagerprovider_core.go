package gocb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/url"
	"strings"
	"time"
)

type userManagerProviderCore struct {
	provider mgmtProvider

	tracer RequestTracer
	meter  *meterWrapper
}

func (um *userManagerProviderCore) tryParseErrorMessage(req *mgmtRequest, resp *mgmtResponse) error {
	b, err := io.ReadAll(resp.Body)
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
		if err := checkForRateLimitError(resp.StatusCode, string(b)); err != nil {
			return makeGenericMgmtError(err, req, resp, string(b))
		}

		bodyErr = errors.New(string(b))
	}

	return makeGenericMgmtError(bodyErr, req, resp, string(b))
}

func (um *userManagerProviderCore) GetAllUsers(opts *GetAllUsersOptions) ([]UserAndMetadata, error) {
	if opts == nil {
		opts = &GetAllUsersOptions{}
	}

	if opts.DomainName == "" {
		opts.DomainName = string(LocalDomain)
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_get_all_users", start)

	path := fmt.Sprintf("/settings/rbac/users/%s", url.PathEscape(opts.DomainName))
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
		return nil, makeGenericMgmtError(err, &req, resp, "")
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

func (um *userManagerProviderCore) GetUser(name string, opts *GetUserOptions) (*UserAndMetadata, error) {
	if opts == nil {
		opts = &GetUserOptions{}
	}

	if opts.DomainName == "" {
		opts.DomainName = string(LocalDomain)
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_get_user", start)

	path := fmt.Sprintf("/settings/rbac/users/%s/%s", url.PathEscape(opts.DomainName), url.PathEscape(name))
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
		return nil, makeGenericMgmtError(err, &req, resp, "")
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

func (um *userManagerProviderCore) UpsertUser(user User, opts *UpsertUserOptions) error {
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

	path := fmt.Sprintf("/settings/rbac/users/%s/%s", url.PathEscape(opts.DomainName), url.PathEscape(user.Username))
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
		return makeGenericMgmtError(err, &req, resp, "")
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

func (um *userManagerProviderCore) DropUser(name string, opts *DropUserOptions) error {
	if opts == nil {
		opts = &DropUserOptions{}
	}

	if opts.DomainName == "" {
		opts.DomainName = string(LocalDomain)
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_drop_user", start)

	path := fmt.Sprintf("/settings/rbac/users/%s/%s", url.PathEscape(opts.DomainName), url.PathEscape(name))
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
		return makeGenericMgmtError(err, &req, resp, "")
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

func (um *userManagerProviderCore) GetRoles(opts *GetRolesOptions) ([]RoleAndDescription, error) {
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
		return nil, makeGenericMgmtError(err, &req, resp, "")
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

func (um *userManagerProviderCore) GetGroup(groupName string, opts *GetGroupOptions) (*Group, error) {
	if groupName == "" {
		return nil, makeInvalidArgumentsError("groupName cannot be empty")
	}
	if opts == nil {
		opts = &GetGroupOptions{}
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_get_group", start)

	path := fmt.Sprintf("/settings/rbac/groups/%s", url.PathEscape(groupName))
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
		return nil, makeGenericMgmtError(err, &req, resp, "")
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

func (um *userManagerProviderCore) GetAllGroups(opts *GetAllGroupsOptions) ([]Group, error) {
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
		return nil, makeGenericMgmtError(err, &req, resp, "")
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

func (um *userManagerProviderCore) UpsertGroup(group Group, opts *UpsertGroupOptions) error {
	if group.Name == "" {
		return makeInvalidArgumentsError("group name cannot be empty")
	}
	if opts == nil {
		opts = &UpsertGroupOptions{}
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_upsert_group", start)

	path := fmt.Sprintf("/settings/rbac/groups/%s", url.PathEscape(group.Name))
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
		return makeGenericMgmtError(err, &req, resp, "")
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

func (um *userManagerProviderCore) DropGroup(groupName string, opts *DropGroupOptions) error {
	if groupName == "" {
		return makeInvalidArgumentsError("groupName cannot be empty")
	}

	if opts == nil {
		opts = &DropGroupOptions{}
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_drop_group", start)

	path := fmt.Sprintf("/settings/rbac/groups/%s", url.PathEscape(groupName))
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
		return makeGenericMgmtError(err, &req, resp, "")
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

func (um *userManagerProviderCore) ChangePassword(newPassword string, opts *ChangePasswordOptions) error {
	if newPassword == "" {
		return makeInvalidArgumentsError("new password cannot be empty")
	}

	if opts == nil {
		opts = &ChangePasswordOptions{}
	}

	start := time.Now()
	defer um.meter.ValueRecord(meterValueServiceManagement, "manager_users_change_password", start)

	path := "/controller/changePassword"
	span := createSpan(um.tracer, opts.ParentSpan, "manager_users_change_password", "management")
	span.SetAttribute("db.operation", "POST "+path)
	defer span.End()

	reqForm := make(url.Values)
	reqForm.Add("password", newPassword)

	req := mgmtRequest{
		Service:       ServiceTypeManagement,
		Method:        "POST",
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
		return makeGenericMgmtError(err, &req, resp, "")
	}
	defer ensureBodyClosed(resp.Body)

	if resp.StatusCode != 200 {
		usrErr := um.tryParseErrorMessage(&req, resp)
		if usrErr != nil {
			return usrErr
		}
		return makeMgmtBadStatusError("failed to change password", &req, resp)
	}

	return nil
}
