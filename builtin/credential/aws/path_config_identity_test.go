package awsauth

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestBackend_pathConfigIdentity(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// Check if default values are returned before setting the configuration
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/identity",
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	if resp.Data["iam_alias"] == nil || resp.Data["iam_alias"] != identityAliasRoleID {
		t.Fatalf("bad: iam_alias; expected: %q, actual: %q", identityAliasIAMUniqueID, resp.Data["iam_alias"])
	}
	if resp.Data["ec2_alias"] == nil || resp.Data["ec2_alias"] != identityAliasRoleID {
		t.Fatalf("bad: ec2_alias; expected: %q, actual: %q", identityAliasIAMUniqueID, resp.Data["ec2_alias"])
	}

	// Invalid value for iam_alias
	data := map[string]interface{}{
		"iam_alias": "invalid",
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/identity",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("nil response from invalid config/identity request")
	}
	if !resp.IsError() {
		t.Fatalf("received non-error response from invalid config/identity request: %#v", resp)
	}

	// Valid value for iam_alias but invalid value for ec2_alias
	data["iam_alias"] = identityAliasIAMFullArn
	data["ec2_alias"] = "invalid"
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/identity",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("nil response from invalid config/identity request")
	}
	if !resp.IsError() {
		t.Fatalf("received non-error response from invalid config/identity request: %#v", resp)
	}

	// Valid value for both iam_alias and ec2_alias
	data["ec2_alias"] = identityAliasEC2ImageID
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/identity",
		Data:      data,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}

	// Check if both values are stored properly
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/identity",
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	if resp.Data["iam_alias"] != identityAliasIAMFullArn {
		t.Fatalf("bad: expected response with iam_alias value of %q; got %#v", identityAliasIAMFullArn, resp.Data["iam_alias"])
	}
	if resp.Data["ec2_alias"] != identityAliasEC2ImageID {
		t.Fatalf("bad: expected response with ec2_alias value of %q; got %#v", identityAliasEC2ImageID, resp.Data["ec2_alias"])
	}

	// Modify one field and ensure that the other one is unchanged
	data["ec2_alias"] = identityAliasEC2InstanceID
	delete(data, "iam_alias")
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/identity",
		Data:      data,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/identity",
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	if resp.Data["iam_alias"] != identityAliasIAMFullArn {
		t.Fatalf("bad: expected response with iam_alias value of %q; got %#v", identityAliasIAMFullArn, resp.Data["iam_alias"])
	}
	if resp.Data["ec2_alias"] != identityAliasEC2InstanceID {
		t.Fatalf("bad: expected response with ec2_alias value of %q; got %#v", identityAliasEC2ImageID, resp.Data["ec2_alias"])
	}

	// Update both iam_alias and ec2_alias
	data["iam_alias"] = identityAliasIAMUniqueID
	data["ec2_alias"] = identityAliasEC2InstanceID
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/identity",
		Data:      data,
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}

	// Check if updates were stored properly
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/identity",
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v\nresp: %#v", err, resp)
	}
	if resp.Data["iam_alias"] != identityAliasIAMUniqueID {
		t.Fatalf("bad: expected response with iam_alias value of %q; got %#v", identityAliasIAMFullArn, resp.Data["iam_alias"])
	}
	if resp.Data["ec2_alias"] != identityAliasEC2InstanceID {
		t.Fatalf("bad: expected response with ec2_alias value of %q; got %#v", identityAliasEC2ImageID, resp.Data["ec2_alias"])
	}
}
