package aws

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestBackend_PathListRoles(t *testing.T) {
	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b := Backend()
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
