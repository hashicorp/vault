// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"time"

	"github.com/hashicorp/jsonapi"
)

// Access returns a pointer to the given team access type.
func Access(v AccessType) *AccessType {
	return &v
}

// ProjectAccess returns a pointer to the given team access project type.
func ProjectAccess(v TeamProjectAccessType) *TeamProjectAccessType {
	return &v
}

// ProjectSettingsPermission returns a pointer to the given team access project type.
func ProjectSettingsPermission(v ProjectSettingsPermissionType) *ProjectSettingsPermissionType {
	return &v
}

// ProjectTeamsPermission returns a pointer to the given team access project type.
func ProjectTeamsPermission(v ProjectTeamsPermissionType) *ProjectTeamsPermissionType {
	return &v
}

// WorkspaceRunsPermission returns a pointer to the given team access project type.
func WorkspaceRunsPermission(v WorkspaceRunsPermissionType) *WorkspaceRunsPermissionType {
	return &v
}

// WorkspaceSentinelMocksPermission returns a pointer to the given team access project type.
func WorkspaceSentinelMocksPermission(v WorkspaceSentinelMocksPermissionType) *WorkspaceSentinelMocksPermissionType {
	return &v
}

// WorkspaceStateVersionsPermission returns a pointer to the given team access project type.
func WorkspaceStateVersionsPermission(v WorkspaceStateVersionsPermissionType) *WorkspaceStateVersionsPermissionType {
	return &v
}

// WorkspaceStateVersionsPermission returns a pointer to the given team access project type.
func WorkspaceVariablesPermission(v WorkspaceVariablesPermissionType) *WorkspaceVariablesPermissionType {
	return &v
}

// RunsPermission returns a pointer to the given team runs permission type.
func RunsPermission(v RunsPermissionType) *RunsPermissionType {
	return &v
}

// VariablesPermission returns a pointer to the given team variables permission type.
func VariablesPermission(v VariablesPermissionType) *VariablesPermissionType {
	return &v
}

// StateVersionsPermission returns a pointer to the given team state versions permission type.
func StateVersionsPermission(v StateVersionsPermissionType) *StateVersionsPermissionType {
	return &v
}

// SentinelMocksPermission returns a pointer to the given team Sentinel mocks permission type.
func SentinelMocksPermission(v SentinelMocksPermissionType) *SentinelMocksPermissionType {
	return &v
}

// AuthPolicy returns a pointer to the given authentication poliy.
func AuthPolicy(v AuthPolicyType) *AuthPolicyType {
	return &v
}

// Bool returns a pointer to the given bool
func Bool(v bool) *bool {
	return &v
}

// Category returns a pointer to the given category type.
func Category(v CategoryType) *CategoryType {
	return &v
}

// EnforcementMode returns a pointer to the given enforcement level.
func EnforcementMode(v EnforcementLevel) *EnforcementLevel {
	return &v
}

// Int returns a pointer to the given int.
func Int(v int) *int {
	return &v
}

// Int64 returns a pointer to the given int64.
func Int64(v int64) *int64 {
	return &v
}

// NotificationDestination returns a pointer to the given notification configuration destination type
func NotificationDestination(v NotificationDestinationType) *NotificationDestinationType {
	return &v
}

// PlanExportType returns a pointer to the given plan export data type.
func PlanExportType(v PlanExportDataType) *PlanExportDataType {
	return &v
}

// ServiceProvider returns a pointer to the given service provider type.
func ServiceProvider(v ServiceProviderType) *ServiceProviderType {
	return &v
}

// SMTPAuthValue returns a pointer to a given smtp auth type.
func SMTPAuthValue(v SMTPAuthType) *SMTPAuthType {
	return &v
}

// String returns a pointer to the given string.
func String(v string) *string {
	return &v
}

func NullableBool(v bool) jsonapi.NullableAttr[bool] {
	return jsonapi.NewNullableAttrWithValue[bool](v)
}

func NullBool() jsonapi.NullableAttr[bool] {
	return jsonapi.NewNullNullableAttr[bool]()
}

func NullableTime(v time.Time) jsonapi.NullableAttr[time.Time] {
	return jsonapi.NewNullableAttrWithValue[time.Time](v)
}

func NullTime() jsonapi.NullableAttr[time.Time] {
	return jsonapi.NewNullNullableAttr[time.Time]()
}
