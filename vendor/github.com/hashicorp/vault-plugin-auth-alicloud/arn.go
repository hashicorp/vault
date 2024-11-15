// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloud

import (
	"errors"
	"fmt"
	"strings"
)

func parseARN(a string) (*arn, error) {
	if a == "" {
		return nil, errors.New("no arn provided")
	}

	// Example: "acs:ram::5138828231865461:assumed-role/elk/vm-ram-i-rj978rorvlg76urhqh7q"
	parsed := &arn{
		Full: a,
	}
	outerFields := strings.Split(a, ":")
	if len(outerFields) != 5 {
		return nil, fmt.Errorf("unrecognized arn: contains %d colon-separated fields, expected 5", len(outerFields))
	}
	if outerFields[0] != "acs" {
		return nil, errors.New(`unrecognized arn: does not begin with "acs:"`)
	}
	if outerFields[1] != "ram" {
		return nil, fmt.Errorf("unrecognized service: %v, not ram", outerFields[1])
	}
	parsed.AccountNumber = outerFields[3]

	roleFields := strings.Split(outerFields[4], "/")
	if len(roleFields) < 2 {
		return nil, fmt.Errorf("unrecognized arn: %q contains fewer than 2 slash-separated roleFields", outerFields[4])
	}

	entityType := roleFields[0]
	switch entityType {
	case "assumed-role":
		parsed.Type = arnTypeAssumedRole
	case "role":
		parsed.Type = arnTypeRole
	default:
		return nil, fmt.Errorf("unsupported parsed type: %s", entityType)
	}

	parsed.RoleName = roleFields[1]
	if len(roleFields) > 2 {
		parsed.RoleAssumerName = roleFields[2]
	}
	return parsed, nil
}

type arnType int

func (t arnType) String() string {
	switch t {
	case arnTypeRole:
		return "role"
	case arnTypeAssumedRole:
		return "assumed-role"
	default:
		return ""
	}
}

const (
	arnTypeRole arnType = iota
	arnTypeAssumedRole
)

type arn struct {
	AccountNumber   string
	Type            arnType
	RoleName        string
	RoleAssumerName string
	Full            string
}

func (a *arn) IsMemberOf(possibleParent *arn) bool {
	if possibleParent.Type != arnTypeRole || a.Type != arnTypeAssumedRole {
		// We currently only support the relationship between roles and assumed roles.
		return false
	}
	if possibleParent.AccountNumber != a.AccountNumber {
		return false
	}
	if possibleParent.RoleName != a.RoleName {
		return false
	}
	return true
}

func (a *arn) String() string {
	return a.Full
}
