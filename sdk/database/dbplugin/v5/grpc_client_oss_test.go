// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package dbplugin

import (
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
)

var _ proto.DatabaseClient = fakeClient{}

type fakeClient struct {
	initResp *proto.InitializeResponse
	initErr  error

	newUserResp *proto.NewUserResponse
	newUserErr  error

	updateUserResp *proto.UpdateUserResponse
	updateUserErr  error

	deleteUserResp *proto.DeleteUserResponse
	deleteUserErr  error

	typeResp *proto.TypeResponse
	typeErr  error

	closeErr error
}
