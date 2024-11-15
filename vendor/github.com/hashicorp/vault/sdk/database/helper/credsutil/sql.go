// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package credsutil

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
)

const (
	NoneLength int = -1
)

// SQLCredentialsProducer implements CredentialsProducer and provides a generic credentials producer for most sql database types.
type SQLCredentialsProducer struct {
	DisplayNameLen    int
	RoleNameLen       int
	UsernameLen       int
	Separator         string
	LowercaseUsername bool
}

func (scp *SQLCredentialsProducer) GenerateCredentials(ctx context.Context) (string, error) {
	password, err := scp.GeneratePassword()
	if err != nil {
		return "", err
	}
	return password, nil
}

func (scp *SQLCredentialsProducer) GenerateUsername(config dbplugin.UsernameConfig) (string, error) {
	caseOp := KeepCase
	if scp.LowercaseUsername {
		caseOp = Lowercase
	}
	return GenerateUsername(
		DisplayName(config.DisplayName, scp.DisplayNameLen),
		RoleName(config.RoleName, scp.RoleNameLen),
		Case(caseOp),
		Separator(scp.Separator),
		MaxLength(scp.UsernameLen),
	)
}

func (scp *SQLCredentialsProducer) GeneratePassword() (string, error) {
	password, err := RandomAlphaNumeric(20, true)
	if err != nil {
		return "", err
	}

	return password, nil
}

func (scp *SQLCredentialsProducer) GenerateExpiration(ttl time.Time) (string, error) {
	return ttl.Format("2006-01-02 15:04:05-0700"), nil
}
