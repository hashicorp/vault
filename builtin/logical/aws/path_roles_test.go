// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/logical"
)

const adminAccessPolicyARN = "arn:aws:iam::aws:policy/AdministratorAccess"

func TestBackend_PathListRoles(t *testing.T) {
	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	roleData := map[string]interface{}{
		"role_arns":       []string{"arn:aws:iam::123456789012:role/path/RoleName"},
		"credential_type": assumedRoleCred,
		"default_sts_ttl": 3600,
		"max_sts_ttl":     3600,
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Data:      roleData,
	}

	for i := 1; i <= 10; i++ {
		roleReq.Path = "roles/testrole" + strconv.Itoa(i)
		resp, err = b.HandleRequest(context.Background(), roleReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: role creation failed. resp:%#v\n err:%v", resp, err)
		}
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ListOperation,
		Path:      "roles",
		Storage:   config.StorageView,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: listing roles failed. resp:%#v\n err:%v", resp, err)
	}

	if len(resp.Data["keys"].([]string)) != 10 {
		t.Fatalf("failed to list all 10 roles")
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ListOperation,
		Path:      "roles/",
		Storage:   config.StorageView,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: listing roles failed. resp:%#v\n err:%v", resp, err)
	}

	if len(resp.Data["keys"].([]string)) != 10 {
		t.Fatalf("failed to list all 10 roles")
	}
}

func TestUpgradeLegacyPolicyEntry(t *testing.T) {
	var input string
	var expected awsRoleEntry
	var output *awsRoleEntry

	input = "arn:aws:iam::123456789012:role/path/RoleName"
	expected = awsRoleEntry{
		CredentialTypes:          []string{assumedRoleCred},
		RoleArns:                 []string{input},
		ProhibitFlexibleCredPath: true,
		Version:                  1,
	}
	output = upgradeLegacyPolicyEntry(input)
	if output.InvalidData != "" {
		t.Fatalf("bad: error processing upgrade of %q: got invalid data of %v", input, output.InvalidData)
	}
	if !reflect.DeepEqual(*output, expected) {
		t.Fatalf("bad: expected %#v; received %#v", expected, *output)
	}

	input = "arn:aws:iam::123456789012:policy/MyPolicy"
	expected = awsRoleEntry{
		CredentialTypes:          []string{iamUserCred},
		PolicyArns:               []string{input},
		ProhibitFlexibleCredPath: true,
		Version:                  1,
	}
	output = upgradeLegacyPolicyEntry(input)
	if output.InvalidData != "" {
		t.Fatalf("bad: error processing upgrade of %q: got invalid data of %v", input, output.InvalidData)
	}
	if !reflect.DeepEqual(*output, expected) {
		t.Fatalf("bad: expected %#v; received %#v", expected, *output)
	}

	input = "arn:aws:iam::aws:policy/AWSManagedPolicy"
	expected.PolicyArns = []string{input}
	output = upgradeLegacyPolicyEntry(input)
	if output.InvalidData != "" {
		t.Fatalf("bad: error processing upgrade of %q: got invalid data of %v", input, output.InvalidData)
	}
	if !reflect.DeepEqual(*output, expected) {
		t.Fatalf("bad: expected %#v; received %#v", expected, *output)
	}

	input = `
{
	"Version": "2012-10-07",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": "ec2:Describe*",
			"Resource": "*"
		}
	]
}`
	compacted, err := compactJSON(input)
	if err != nil {
		t.Fatalf("error parsing JSON: %v", err)
	}
	expected = awsRoleEntry{
		CredentialTypes:          []string{iamUserCred, federationTokenCred},
		PolicyDocument:           compacted,
		ProhibitFlexibleCredPath: true,
		Version:                  1,
	}
	output = upgradeLegacyPolicyEntry(input)
	if output.InvalidData != "" {
		t.Fatalf("bad: error processing upgrade of %q: got invalid data of %v", input, output.InvalidData)
	}
	if !reflect.DeepEqual(*output, expected) {
		t.Fatalf("bad: expected %#v; received %#v", expected, *output)
	}

	// Due to lack of prior input validation, this could exist in the storage, and we need
	// to be able to read it out in some fashion, so have to handle this in a poor fashion
	input = "arn:gobbledygook"
	expected = awsRoleEntry{
		InvalidData: input,
		Version:     1,
	}
	output = upgradeLegacyPolicyEntry(input)
	if !reflect.DeepEqual(*output, expected) {
		t.Fatalf("bad: expected %#v; received %#v", expected, *output)
	}
}

func TestUserPathValidity(t *testing.T) {
	testCases := []struct {
		description string
		userPath    string
		isValid     bool
	}{
		{
			description: "Default",
			userPath:    "/",
			isValid:     true,
		},
		{
			description: "Empty",
			userPath:    "",
			isValid:     false,
		},
		{
			description: "Valid",
			userPath:    "/path/",
			isValid:     true,
		},
		{
			description: "Missing leading slash",
			userPath:    "path/",
			isValid:     false,
		},
		{
			description: "Missing trailing slash",
			userPath:    "/path",
			isValid:     false,
		},
		{
			description: "Invalid character",
			userPath:    "/šiauliai/",
			isValid:     false,
		},
		{
			description: "Max length",
			userPath:    "/" + strings.Repeat("a", 510) + "/",
			isValid:     true,
		},
		{
			description: "Too long",
			userPath:    "/" + strings.Repeat("a", 511) + "/",
			isValid:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			if tc.isValid != userPathRegex.MatchString(tc.userPath) {
				t.Fatalf("bad: expected %s", strconv.FormatBool(tc.isValid))
			}
		})
	}
}

func TestRoleCRUDWithPermissionsBoundary(t *testing.T) {
	roleName := "test_perm_boundary"

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	permissionsBoundaryARN := "arn:aws:iam::aws:policy/EC2FullAccess"

	roleData := map[string]interface{}{
		"credential_type":          iamUserCred,
		"policy_arns":              []string{adminAccessPolicyARN},
		"permissions_boundary_arn": permissionsBoundaryARN,
	}
	request := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + roleName,
		Storage:   config.StorageView,
		Data:      roleData,
	}
	resp, err := b.HandleRequest(context.Background(), request)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: role creation failed. resp:%#v\nerr:%v", resp, err)
	}

	request = &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "roles/" + roleName,
		Storage:   config.StorageView,
	}
	resp, err = b.HandleRequest(context.Background(), request)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: reading role failed. resp:%#v\nerr:%v", resp, err)
	}
	if resp.Data["credential_type"] != iamUserCred {
		t.Errorf("bad: expected credential_type of %s, got %s instead", iamUserCred, resp.Data["credential_type"])
	}
	if resp.Data["permissions_boundary_arn"] != permissionsBoundaryARN {
		t.Errorf("bad: expected permissions_boundary_arn of %s, got %s instead", permissionsBoundaryARN, resp.Data["permissions_boundary_arn"])
	}
}

func TestRoleWithPermissionsBoundaryValidation(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	roleData := map[string]interface{}{
		"credential_type":          assumedRoleCred, // only iamUserCred supported with permissions_boundary_arn
		"role_arns":                []string{"arn:aws:iam::123456789012:role/VaultRole"},
		"permissions_boundary_arn": "arn:aws:iam::aws:policy/FooBar",
	}
	request := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test_perm_boundary",
		Storage:   config.StorageView,
		Data:      roleData,
	}
	resp, err := b.HandleRequest(context.Background(), request)
	if err == nil && (resp == nil || !resp.IsError()) {
		t.Fatalf("bad: expected role creation to fail due to bad credential_type, but it didn't. resp:%#v\nerr:%v", resp, err)
	}

	roleData = map[string]interface{}{
		"credential_type":          iamUserCred,
		"policy_arns":              []string{adminAccessPolicyARN},
		"permissions_boundary_arn": "arn:aws:notiam::aws:policy/FooBar",
	}
	request.Data = roleData
	resp, err = b.HandleRequest(context.Background(), request)
	if err == nil && (resp == nil || !resp.IsError()) {
		t.Fatalf("bad: expected role creation to fail due to malformed permissions_boundary_arn, but it didn't. resp:%#v\nerr:%v", resp, err)
	}
}

func TestValidateAWSManagedPolicy(t *testing.T) {
	expectErr := func(arn string) {
		err := validateAWSManagedPolicy(arn)
		if err == nil {
			t.Errorf("bad: expected arn of %s to return an error but it didn't", arn)
		}
	}

	expectErr("not_an_arn")
	expectErr("notarn:aws:iam::aws:policy/FooBar")
	expectErr("arn:aws:notiam::aws:policy/FooBar")
	expectErr("arn:aws:iam::aws:notpolicy/FooBar")
	expectErr("arn:aws:iam::aws:policynot/FooBar")

	arn := "arn:aws:iam::aws:policy/FooBar"
	err := validateAWSManagedPolicy(arn)
	if err != nil {
		t.Errorf("bad: expected arn of %s to not return an error but it did: %#v", arn, err)
	}
}

func TestRoleEntryValidationCredTypes(t *testing.T) {
	roleEntry := awsRoleEntry{
		CredentialTypes: []string{},
		PolicyArns:      []string{adminAccessPolicyARN},
	}
	if roleEntry.validate() == nil {
		t.Errorf("bad: invalid roleEntry with no CredentialTypes %#v passed validation", roleEntry)
	}
	roleEntry.CredentialTypes = []string{"invalid_type"}
	if roleEntry.validate() == nil {
		t.Errorf("bad: invalid roleEntry with invalid CredentialTypes %#v passed validation", roleEntry)
	}
	roleEntry.CredentialTypes = []string{iamUserCred, "invalid_type"}
	if roleEntry.validate() == nil {
		t.Errorf("bad: invalid roleEntry with invalid CredentialTypes %#v passed validation", roleEntry)
	}
}

func TestRoleEntryValidationIamUserCred(t *testing.T) {
	allowAllPolicyDocument := `{"Version": "2012-10-17", "Statement": [{"Sid": "AllowAll", "Effect": "Allow", "Action": "*", "Resource": "*"}]}`
	roleEntry := awsRoleEntry{
		CredentialTypes:        []string{iamUserCred},
		PolicyArns:             []string{adminAccessPolicyARN},
		PermissionsBoundaryARN: adminAccessPolicyARN,
	}
	err := roleEntry.validate()
	if err != nil {
		t.Errorf("bad: valid roleEntry %#v failed validation: %v", roleEntry, err)
	}
	roleEntry.PolicyDocument = allowAllPolicyDocument
	err = roleEntry.validate()
	if err != nil {
		t.Errorf("bad: valid roleEntry %#v failed validation: %v", roleEntry, err)
	}
	roleEntry.PolicyArns = []string{}
	err = roleEntry.validate()
	if err != nil {
		t.Errorf("bad: valid roleEntry %#v failed validation: %v", roleEntry, err)
	}

	roleEntry = awsRoleEntry{
		CredentialTypes: []string{iamUserCred},
		RoleArns:        []string{"arn:aws:iam::123456789012:role/SomeRole"},
	}
	assertMultiError(t, roleEntry.validate(),
		[]error{
			errors.New(
				"cannot supply role_arns when credential_type isn't assumed_role",
			),
		})

	roleEntry = awsRoleEntry{
		CredentialTypes: []string{iamUserCred},
		PolicyArns:      []string{adminAccessPolicyARN},
		DefaultSTSTTL:   1,
	}
	assertMultiError(t, roleEntry.validate(),
		[]error{
			errors.New(
				"default_sts_ttl parameter only valid for assumed_role, federation_token, and session_token credential types",
			),
		})
	roleEntry.DefaultSTSTTL = 0

	roleEntry.MaxSTSTTL = 1
	assertMultiError(t, roleEntry.validate(),
		[]error{
			errors.New(
				"max_sts_ttl parameter only valid for assumed_role, federation_token, and session_token credential types",
			),
		})
	roleEntry.MaxSTSTTL = 0

	roleEntry.SessionTags = map[string]string{
		"Key1": "Value1",
		"Key2": "Value2",
	}
	assertMultiError(t, roleEntry.validate(),
		[]error{
			errors.New(
				"cannot supply session_tags when credential_type isn't assumed_role",
			),
		})
	roleEntry.SessionTags = nil

	roleEntry.ExternalID = "my-ext-id"
	assertMultiError(t, roleEntry.validate(),
		[]error{
			errors.New(
				"cannot supply external_id when credential_type isn't assumed_role"),
		})
}

func assertMultiError(t *testing.T, err error, expected []error) {
	t.Helper()

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	var multiErr *multierror.Error
	if errors.As(err, &multiErr) {
		if multiErr.Len() != len(expected) {
			t.Errorf("expected %d error, got %d", len(expected), multiErr.Len())
		} else {
			if !reflect.DeepEqual(expected, multiErr.Errors) {
				t.Errorf("expected error %q, actual %q", expected, multiErr.Errors)
			}
		}
	} else {
		t.Errorf("expected multierror, got %T", err)
	}
}

func TestRoleEntryValidationAssumedRoleCred(t *testing.T) {
	allowAllPolicyDocument := `{"Version": "2012-10-17", "Statement": [{"Sid": "AllowAll", "Effect": "Allow", "Action": "*", "Resource": "*"}]}`
	roleEntry := awsRoleEntry{
		CredentialTypes: []string{assumedRoleCred},
		RoleArns:        []string{"arn:aws:iam::123456789012:role/SomeRole"},
		PolicyArns:      []string{adminAccessPolicyARN},
		PolicyDocument:  allowAllPolicyDocument,
		ExternalID:      "my-ext-id",
		SessionTags: map[string]string{
			"Key1": "Value1",
			"Key2": "Value2",
		},
		DefaultSTSTTL: 2,
		MaxSTSTTL:     3,
	}
	if err := roleEntry.validate(); err != nil {
		t.Errorf("bad: valid roleEntry %#v failed validation: %v", roleEntry, err)
	}

	roleEntry.MaxSTSTTL = 1
	if roleEntry.validate() == nil {
		t.Errorf("bad: invalid roleEntry with MaxSTSTTL < DefaultSTSTTL %#v passed validation", roleEntry)
	}
	roleEntry.MaxSTSTTL = 0
	roleEntry.UserPath = "/foobar/"
	if roleEntry.validate() == nil {
		t.Errorf("bad: invalid roleEntry with unrecognized UserPath %#v passed validation", roleEntry)
	}
	roleEntry.UserPath = ""
	roleEntry.PermissionsBoundaryARN = adminAccessPolicyARN
	if roleEntry.validate() == nil {
		t.Errorf("bad: invalid roleEntry with unrecognized PermissionsBoundary %#v passed validation", roleEntry)
	}
}

func TestRoleEntryValidationFederationTokenCred(t *testing.T) {
	allowAllPolicyDocument := `{"Version": "2012-10-17", "Statement": [{"Sid": "AllowAll", "Effect": "Allow", "Action": "*", "Resource": "*"}]}`
	roleEntry := awsRoleEntry{
		CredentialTypes: []string{federationTokenCred},
		PolicyDocument:  allowAllPolicyDocument,
		PolicyArns:      []string{adminAccessPolicyARN},
		DefaultSTSTTL:   2,
		MaxSTSTTL:       3,
	}
	if err := roleEntry.validate(); err != nil {
		t.Errorf("bad: valid roleEntry %#v failed validation: %v", roleEntry, err)
	}

	roleEntry.RoleArns = []string{"arn:aws:iam::123456789012:role/SomeRole"}
	if roleEntry.validate() == nil {
		t.Errorf("bad: invalid roleEntry with unrecognized RoleArns %#v passed validation", roleEntry)
	}
	roleEntry.RoleArns = []string{}
	roleEntry.UserPath = "/foobar/"
	if roleEntry.validate() == nil {
		t.Errorf("bad: invalid roleEntry with unrecognized UserPath %#v passed validation", roleEntry)
	}

	roleEntry.UserPath = ""
	roleEntry.MaxSTSTTL = 1
	if roleEntry.validate() == nil {
		t.Errorf("bad: invalid roleEntry with MaxSTSTTL < DefaultSTSTTL %#v passed validation", roleEntry)
	}
	roleEntry.MaxSTSTTL = 0
	roleEntry.PermissionsBoundaryARN = adminAccessPolicyARN
	if roleEntry.validate() == nil {
		t.Errorf("bad: invalid roleEntry with unrecognized PermissionsBoundary %#v passed validation", roleEntry)
	}
}
