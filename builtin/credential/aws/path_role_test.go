package awsauth

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
)

func TestBackend_pathRoleInstanceIdentityDocument(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"policies":     "p,q,r,s",
		"max_ttl":      "2h",
		"bound_ami_id": "ami-abcd123",
	}
	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/ami-abcd123",
		Data:      data,
		Storage:   storage,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role")
	}
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/ami-abcd123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("failed to read the role entry")
	}
	if !policyutil.EquivalentPolicies(strings.Split(data["policies"].(string), ","), resp.Data["policies"].([]string)) {
		t.Fatalf("bad: policies: expected: %#v\ngot: %#v\n", data, resp.Data)
	}

	data["allow_instance_migration"] = true
	data["disallow_reauthentication"] = true
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/ami-abcd123",
		Data:      data,
		Storage:   storage,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role")
	}
	if err != nil {
		t.Fatal(err)
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/ami-abcd123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Data["allow_instance_migration"].(bool) || !resp.Data["disallow_reauthentication"].(bool) {
		t.Fatal("bad: expected:true got:false\n")
	}

	// add another entry, to test listing of role entries
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/ami-abcd456",
		Data:      data,
		Storage:   storage,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role")
	}
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "roles",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil || resp.IsError() {
		t.Fatalf("failed to list the role entries")
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != 2 {
		t.Fatalf("bad: keys: %#v\n", keys)
	}

	_, err = b.HandleRequest(&logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "role/ami-abcd123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/ami-abcd123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("bad: response: expected:nil actual:%#v\n", resp)
	}

}

func TestBackend_pathRoleSignedCallerIdentity(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	// make sure we start with empty roles, which gives us confidence that the read later
	// actually is the two roles we created
	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "roles",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil || resp.IsError() {
		t.Fatalf("failed to list role entries")
	}
	if resp.Data["keys"] != nil {
		t.Fatalf("Received roles when expected none")
	}

	data := map[string]interface{}{
		"allowed_auth_types":      "signed_caller_identity_request",
		"policies":                "p,q,r,s",
		"max_ttl":                 "2h",
		"bound_iam_principal_arn": "n:aws:iam::123456789012:user/MyUserName",
	}
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/MyRoleName",
		Data:      data,
		Storage:   storage,
	})

	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal("failed to create the role entry")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/MyRoleName",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("failed to read the role entry")
	}
	if !policyutil.EquivalentPolicies(strings.Split(data["policies"].(string), ","), resp.Data["policies"].([]string)) {
		t.Fatalf("bad: policies: expected %#v\ngot: %#v\n", data, resp.Data)
	}

	data["infer_role_as_type"] = "invalid"
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/ShouldNeverExist",
		Data:      data,
		Storage:   storage,
	})
	if resp == nil || !resp.IsError() {
		t.Fatalf("Created role with invalid infer_role_as_type")
	}
	if err != nil {
		t.Fatal(err)
	}

	data["infer_role_as_type"] = "ec2Instance"
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/ShouldNeverExist",
		Data:      data,
		Storage:   storage,
	})
	if resp == nil || !resp.IsError() {
		t.Fatalf("Created role without necessary inferred_aws_region")
	}
	if err != nil {
		t.Fatal(err)
	}

	delete(data, "bound_iam_principal_arn")
	data["inferred_aws_region"] = "us-east-1"
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/ShouldNeverExist",
		Data:      data,
		Storage:   storage,
	})
	if resp == nil || !resp.IsError() {
		t.Fatalf("Created role without anything bound")
	}
	if err != nil {
		t.Fatal(err)
	}

	// generate a second role, ensure we're able to list both
	data["bound_ami_id"] = "ami-abcd123"
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/MyOtherRoleName",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create additional role: %s")
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "roles",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.Data == nil || resp.IsError() {
		t.Fatalf("failed to list role entries")
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != 2 {
		t.Fatalf("bad: keys %#v\n", keys)
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "role/MyOtherRoleName",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/MyOtherRoleName",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("bad: response: expected: nil actual:%3v\n", resp)
	}
}

func TestBackend_pathRoleMixedTypes(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"policies":           "p,q,r,s",
		"bound_ami_id":       "ami-abc1234",
		"allowed_auth_types": "instance_identity_document,invalid",
	}

	submitCreateRequest := func(roleName string) (*logical.Response, error) {
		return b.HandleRequest(&logical.Request{
			Operation: logical.CreateOperation,
			Path:      "role/" + roleName,
			Data:      data,
			Storage:   storage,
		})
	}

	resp, err := submitCreateRequest("shouldNeverExist")
	if resp != nil && !resp.IsError() {
		t.Fatalf("created role with invalid allowed_auth_type")
	}
	if err != nil {
		t.Fatal(err)
	}

	data["allowed_auth_types"] = "instance_identity_document,,signed_caller_identity_request"
	resp, err = submitCreateRequest("shouldNeverExist")
	if resp != nil && !resp.IsError() {
		t.Fatalf("created role without required bound_iam_principal_arn")
	}
	if err != nil {
		t.Fatal(err)
	}

	data["bound_iam_principal_arn"] = "arn:aws:iam::123456789012:role/MyRole"
	delete(data, "bound_ami_id")
	resp, err = submitCreateRequest("shouldNeverExist")
	if resp != nil && !resp.IsError() {
		t.Fatalf("created role without required instance_identity_document binding")
	}
	if err != nil {
		t.Fatal(err)
	}

	data["bound_ami_id"] = "ami-1234567"
	resp, err = submitCreateRequest("multipleTypes")
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("didn't allow creation of valid role with multiple bindings of different types")
	}

	delete(data, "bound_iam_principal_arn")
	data["infer_role_as_type"] = "ec2Instance"
	data["inferred_aws_region"] = "us-east-1"
	resp, err = submitCreateRequest("multipleTypesInferred")
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("didn't allow creation of roles with only inferred bindings")
	}
}
