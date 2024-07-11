// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki_backend

import (
	"context"
	"fmt"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type SystemViewGetter interface {
	System() logical.SystemView
}

type MountInfo interface {
	BackendUUID() string
}

type Logger interface {
	Logger() log.Logger
}

//go:generate enumer -type=RolePathPolicy -text -json -transform=kebab-case
type RolePathPolicy int

const (
	RPPUnknown RolePathPolicy = iota
	RPPSignVerbatim
	RPPRole
)

var (
	pathPolicyRolePrefix       = "role:"
	pathPolicyRolePrefixLength = len(pathPolicyRolePrefix)
)

// GetRoleByPathOrPathPolicy loads an existing role based on if the data field data contains a 'role' parameter
// or by the values within the pathPolicy
func GetRoleByPathOrPathPolicy(ctx context.Context, s logical.Storage, data *framework.FieldData, pathPolicy string) (*issuing.RoleEntry, error) {
	var role *issuing.RoleEntry

	// The role name from the path is the highest priority
	if roleName, ok := getRoleNameFromPath(data); ok {
		var err error
		role, err = issuing.GetRole(ctx, s, roleName)
		if err != nil {
			return nil, err
		}
	} else {
		policyType, policyVal, err := GetPathPolicyType(pathPolicy)
		if err != nil {
			return nil, err
		}

		switch policyType {
		case RPPRole:
			role, err = issuing.GetRole(ctx, s, policyVal)
			if err != nil {
				return nil, err
			}
		case RPPSignVerbatim:
			role = issuing.SignVerbatimRole()
		default:
			return nil, fmt.Errorf("unsupported policy type returned: %s from policy path: %s", policyType, pathPolicy)
		}
	}

	return role, nil
}

func GetPathPolicyType(pathPolicy string) (RolePathPolicy, string, error) {
	policy := strings.TrimSpace(pathPolicy)

	switch {
	case policy == "sign-verbatim":
		return RPPSignVerbatim, "", nil
	case strings.HasPrefix(policy, pathPolicyRolePrefix):
		if policy == pathPolicyRolePrefix {
			return RPPUnknown, "", fmt.Errorf("no role specified by policy %v", pathPolicy)
		}
		roleName := pathPolicy[pathPolicyRolePrefixLength:]
		return RPPRole, roleName, nil
	default:
		return RPPUnknown, "", fmt.Errorf("string %v was not a valid default path policy", pathPolicy)
	}
}

func getRoleNameFromPath(data *framework.FieldData) (string, bool) {
	// If our schema doesn't include the parameter bail
	if _, ok := data.Schema["role"]; !ok {
		return "", false
	}

	if roleName, ok := data.GetOk("role"); ok {
		return roleName.(string), true
	}

	return "", false
}
