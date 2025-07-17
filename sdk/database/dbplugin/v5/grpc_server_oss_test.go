// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package dbplugin

import (
	"github.com/hashicorp/vault/sdk/logical"
)

var _ Database = fakeDatabase{}

type fakeDatabase struct {
	initResp InitializeResponse
	initErr  error

	newUserResp NewUserResponse
	newUserErr  error

	updateUserResp UpdateUserResponse
	updateUserErr  error

	deleteUserResp DeleteUserResponse
	deleteUserErr  error

	typeResp string
	typeErr  error

	closeErr error
}

var _ Database = &recordingDatabase{}

type recordingDatabase struct {
	initializeCalls int
	newUserCalls    int
	updateUserCalls int
	deleteUserCalls int
	typeCalls       int
	closeCalls      int

	// recordingDatabase can act as middleware so we can record the calls to other test Database implementations
	next Database
}

type fakeDatabaseWithVersion struct {
	version string
}

var (
	_ Database                = (*fakeDatabaseWithVersion)(nil)
	_ logical.PluginVersioner = (*fakeDatabaseWithVersion)(nil)
)
