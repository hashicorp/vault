package awsauth

import (
	"context"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
	vlttesting "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/helper/awsutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestBackend_pathRoleEc2(t *testing.T) {
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

	data := map[string]interface{}{
		"auth_type":    "ec2",
		"policies":     "p,q,r,s",
		"max_ttl":      "2h",
		"bound_ami_id": "ami-abcd123",
	}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
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

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/ami-abcd123",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected failure to create role with both allow_instance_migration true and disallow_reauthentication true")
	}
	data["disallow_reauthentication"] = false
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/ami-abcd123",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failure to update role: %v", resp.Data["error"])
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "role/ami-abcd123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Data["allow_instance_migration"].(bool) {
		t.Fatal("bad: expected allow_instance_migration:true got:false\n")
	}

	if resp.Data["disallow_reauthentication"].(bool) {
		t.Fatal("bad: expected disallow_reauthentication: false got:true\n")
	}

	// add another entry, to test listing of role entries
	data["bound_ami_id"] = "ami-abcd456"
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/ami-abcd456",
		Data:      data,
		Storage:   storage,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role: %s", resp.Data["error"])
	}
	if err != nil {
		t.Fatal(err)
	}

	data["bound_iam_principal_arn"] = ""
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "role/ami-abcd456",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to update role with empty bound_iam_principal_arn: %s", resp.Data["error"])
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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

	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "role/ami-abcd123",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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

func Test_enableIamIDResolution(t *testing.T) {
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
	roleName := "upgradable_role"

	b.resolveArnToUniqueIDFunc = resolveArnToFakeUniqueId

	boundIamRoleARNs := []string{"arn:aws:iam::123456789012:role/MyRole", "arn:aws:iam::123456789012:role/path/*"}
	data := map[string]interface{}{
		"auth_type":               iamAuthType,
		"policies":                "p,q",
		"bound_iam_principal_arn": boundIamRoleARNs,
		"resolve_aws_unique_ids":  false,
	}

	submitRequest := func(roleName string, op logical.Operation) (*logical.Response, error) {
		return b.HandleRequest(context.Background(), &logical.Request{
			Operation: op,
			Path:      "role/" + roleName,
			Data:      data,
			Storage:   storage,
		})
	}

	resp, err := submitRequest(roleName, logical.CreateOperation)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role: %#v", resp)
	}

	resp, err = submitRequest(roleName, logical.ReadOperation)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("failed to read role: resp:%#v,\nerr:%#v", resp, err)
	}
	if resp.Data["bound_iam_principal_id"] != nil && len(resp.Data["bound_iam_principal_id"].([]string)) > 0 {
		t.Fatalf("expected to get no unique ID in role, but got %q", resp.Data["bound_iam_principal_id"])
	}

	data = map[string]interface{}{
		"resolve_aws_unique_ids": true,
	}
	resp, err = submitRequest(roleName, logical.UpdateOperation)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("unable to upgrade role to resolve internal IDs: resp:%#v", resp)
	}

	resp, err = submitRequest(roleName, logical.ReadOperation)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("failed to read role: resp:%#v,\nerr:%#v", resp, err)
	}
	principalIDs := resp.Data["bound_iam_principal_id"].([]string)
	if len(principalIDs) != 1 || principalIDs[0] != "FakeUniqueId1" {
		t.Fatalf("bad: expected upgrade of role resolve principal ID to %q, but got %q instead", "FakeUniqueId1", resp.Data["bound_iam_principal_id"])
	}
	returnedARNs := resp.Data["bound_iam_principal_arn"].([]string)
	if !strutil.EquivalentSlices(returnedARNs, boundIamRoleARNs) {
		t.Fatalf("bad: expected to return bound_iam_principal_arn of %q, but got %q instead", boundIamRoleARNs, returnedARNs)
	}
}

func TestBackend_pathIam(t *testing.T) {
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

	// make sure we start with empty roles, which gives us confidence that the read later
	// actually is the two roles we created
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
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
		"auth_type":               iamAuthType,
		"policies":                "p,q,r,s",
		"max_ttl":                 "2h",
		"bound_iam_principal_arn": "n:aws:iam::123456789012:user/MyUserName",
		"resolve_aws_unique_ids":  false,
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/MyRoleName",
		Data:      data,
		Storage:   storage,
	})

	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create the role entry; resp: %#v", resp)
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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

	data["inferred_entity_type"] = "invalid"
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/ShouldNeverExist",
		Data:      data,
		Storage:   storage,
	})
	if resp == nil || !resp.IsError() {
		t.Fatalf("Created role with invalid inferred_entity_type")
	}
	if err != nil {
		t.Fatal(err)
	}

	data["inferred_entity_type"] = ec2EntityType
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
	secondRole := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/MyOtherRoleName",
		Data:      data,
		Storage:   storage,
	}
	resp, err = b.HandleRequest(context.Background(), secondRole)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create additional role: %v", *secondRole)
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "role/MyOtherRoleName",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"policies":     "p,q,r,s",
		"bound_ami_id": "ami-abc1234",
		"auth_type":    "ec2,invalid",
	}

	submitRequest := func(roleName string, op logical.Operation) (*logical.Response, error) {
		return b.HandleRequest(context.Background(), &logical.Request{
			Operation: op,
			Path:      "role/" + roleName,
			Data:      data,
			Storage:   storage,
		})
	}

	resp, err := submitRequest("shouldNeverExist", logical.CreateOperation)
	if resp == nil || !resp.IsError() {
		t.Fatalf("created role with invalid auth_type; resp: %#v", resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	data["auth_type"] = "ec2,,iam"
	resp, err = submitRequest("shouldNeverExist", logical.CreateOperation)
	if resp == nil || !resp.IsError() {
		t.Fatalf("created role mixed auth types")
	}
	if err != nil {
		t.Fatal(err)
	}

	data["auth_type"] = ec2AuthType
	resp, err = submitRequest("ec2_to_iam", logical.CreateOperation)
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create valid role; resp: %#v", resp)
	}
	if err != nil {
		t.Fatal(err)
	}

	data["auth_type"] = iamAuthType
	delete(data, "bound_ami_id")
	boundIamPrincipalARNs := []string{"arn:aws:iam::123456789012:role/MyRole", "arn:aws:iam::123456789012:role/path/*"}
	data["bound_iam_principal_arn"] = boundIamPrincipalARNs
	resp, err = submitRequest("ec2_to_iam", logical.UpdateOperation)
	if resp == nil || !resp.IsError() {
		t.Fatalf("changed auth type on the role")
	}
	if err != nil {
		t.Fatal(err)
	}

	data["inferred_entity_type"] = ec2EntityType
	data["inferred_aws_region"] = "us-east-1"
	data["resolve_aws_unique_ids"] = false
	resp, err = submitRequest("multipleTypesInferred", logical.CreateOperation)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("didn't allow creation of roles with only inferred bindings")
	}

	b.resolveArnToUniqueIDFunc = resolveArnToFakeUniqueId
	data["resolve_aws_unique_ids"] = true
	resp, err = submitRequest("withInternalIdResolution", logical.CreateOperation)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsError() {
		t.Fatalf("didn't allow creation of role resolving unique IDs")
	}
	resp, err = submitRequest("withInternalIdResolution", logical.ReadOperation)
	if err != nil {
		t.Fatal(err)
	}
	principalIDs := resp.Data["bound_iam_principal_id"].([]string)
	if len(principalIDs) != 1 || principalIDs[0] != "FakeUniqueId1" {
		t.Fatalf("expected fake unique ID of FakeUniqueId1, got %q", resp.Data["bound_iam_principal_id"])
	}
	returnedARNs := resp.Data["bound_iam_principal_arn"].([]string)
	if !strutil.EquivalentSlices(returnedARNs, boundIamPrincipalARNs) {
		t.Fatalf("bad: expected to return bound_iam_principal_arn of %q, but got %q instead", boundIamPrincipalARNs, returnedARNs)
	}
	data["resolve_aws_unique_ids"] = false
	resp, err = submitRequest("withInternalIdResolution", logical.UpdateOperation)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("allowed changing resolve_aws_unique_ids from true to false")
	}
}

func TestAwsEc2_RoleCrud(t *testing.T) {
	var err error
	var resp *logical.Response
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

	role1Data := map[string]interface{}{
		"auth_type":                "ec2",
		"bound_vpc_id":             "testvpcid",
		"allow_instance_migration": true,
		"policies":                 "testpolicy1,testpolicy2",
	}
	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Path:      "role/role1",
		Data:      role1Data,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleData := map[string]interface{}{
		"auth_type":                      "ec2",
		"bound_ami_id":                   "testamiid",
		"bound_account_id":               "testaccountid",
		"bound_region":                   "testregion",
		"bound_iam_role_arn":             "arn:aws:iam::123456789012:role/MyRole",
		"bound_iam_instance_profile_arn": "arn:aws:iam::123456789012:instance-profile/MyInstancePro*",
		"bound_subnet_id":                "testsubnetid",
		"bound_vpc_id":                   "testvpcid",
		"bound_ec2_instance_id":          "i-12345678901234567,i-76543210987654321",
		"role_tag":                       "testtag",
		"resolve_aws_unique_ids":         false,
		"allow_instance_migration":       true,
		"ttl":                            "10m",
		"max_ttl":                        "20m",
		"policies":                       "testpolicy1,testpolicy2",
		"disallow_reauthentication":      false,
		"hmac_key":                       "testhmackey",
		"period":                         "1m",
	}

	roleReq.Path = "role/testrole"
	roleReq.Data = roleData
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	expected := map[string]interface{}{
		"auth_type":                      ec2AuthType,
		"bound_ami_id":                   []string{"testamiid"},
		"bound_account_id":               []string{"testaccountid"},
		"bound_region":                   []string{"testregion"},
		"bound_ec2_instance_id":          []string{"i-12345678901234567", "i-76543210987654321"},
		"bound_iam_principal_arn":        []string{},
		"bound_iam_principal_id":         []string{},
		"bound_iam_role_arn":             []string{"arn:aws:iam::123456789012:role/MyRole"},
		"bound_iam_instance_profile_arn": []string{"arn:aws:iam::123456789012:instance-profile/MyInstancePro*"},
		"bound_subnet_id":                []string{"testsubnetid"},
		"bound_vpc_id":                   []string{"testvpcid"},
		"inferred_entity_type":           "",
		"inferred_aws_region":            "",
		"resolve_aws_unique_ids":         false,
		"role_tag":                       "testtag",
		"allow_instance_migration":       true,
		"ttl":                            int64(600),
		"token_ttl":                      int64(600),
		"max_ttl":                        int64(1200),
		"token_max_ttl":                  int64(1200),
		"token_explicit_max_ttl":         int64(0),
		"policies":                       []string{"testpolicy1", "testpolicy2"},
		"token_policies":                 []string{"testpolicy1", "testpolicy2"},
		"disallow_reauthentication":      false,
		"period":                         int64(60),
		"token_period":                   int64(60),
		"token_bound_cidrs":              []string{},
		"token_no_default_policy":        false,
		"token_num_uses":                 0,
		"token_type":                     "default",
	}

	if resp.Data["role_id"] == nil {
		t.Fatal("role_id not found in repsonse")
	}
	expected["role_id"] = resp.Data["role_id"]
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	roleData["bound_vpc_id"] = "newvpcid"
	roleReq.Operation = logical.UpdateOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}
	expected["bound_vpc_id"] = []string{"newvpcid"}
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: role data: expected: %#v\n actual: %#v", expected, resp.Data)
	}

	// Create a new backend so we have a new cache (thus populating from disk).
	// Then test reading (reading from disk + lock), writing, reading,
	// deleting, reading.
	b, err = Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// Read again, make sure things are what we expect
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: role data: expected: %#v\n actual: %#v", expected, resp.Data)
	}

	roleReq.Operation = logical.UpdateOperation
	roleData["bound_ami_id"] = "testamiid2"
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	expected["bound_ami_id"] = []string{"testamiid2"}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Delete which should remove from disk and also cache
	roleReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}
	if resp != nil {
		t.Fatalf("failed to delete role entry")
	}

	// Verify it was deleted, e.g. it isn't found in the role cache
	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}
	if resp != nil {
		t.Fatal("expected nil")
	}
}

func TestAwsEc2_RoleDurationSeconds(t *testing.T) {
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

	roleData := map[string]interface{}{
		"auth_type":                      "ec2",
		"bound_iam_instance_profile_arn": "arn:aws:iam::123456789012:instance-profile/test-profile-name",
		"resolve_aws_unique_ids":         false,
		"ttl":                            "10s",
		"max_ttl":                        "20s",
		"period":                         "30s",
	}

	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   storage,
		Path:      "role/testrole",
		Data:      roleData,
	}

	resp, err := b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	if resp.Data["ttl"].(int64) != 10 {
		t.Fatalf("bad: ttl; expected: 10, actual: %d", resp.Data["ttl"])
	}
	if resp.Data["max_ttl"].(int64) != 20 {
		t.Fatalf("bad: max_ttl; expected: 20, actual: %d", resp.Data["max_ttl"])
	}
	if resp.Data["period"].(int64) != 30 {
		t.Fatalf("bad: period; expected: 30, actual: %d", resp.Data["period"])
	}
}

func TestAwsIam_RoleDurationSeconds(t *testing.T) {
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

	roleData := map[string]interface{}{
		"auth_type":               "iam",
		"bound_iam_principal_arn": "arn:aws:iam::123456789012:role/testrole",
		"resolve_aws_unique_ids":  false,
		"ttl":                     "20m",
		"max_ttl":                 "30m",
	}

	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   storage,
		Path:      "role/testrole",
		Data:      roleData,
	}

	resp, err := b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	// since default lease TTL for system is 24hr, Token TTL should not be capped
	// since max lease TTL for system is 48hr, Token Max TTL should not be capped
	if resp.Data["token_ttl"].(int64) != 1200 {
		t.Fatalf("bad: token_ttl; expected: 1200, actual: %d", resp.Data["ttl"])
	}
	if resp.Data["max_ttl"].(int64) != 1800 {
		t.Fatalf("bad: max_ttl; expected: 1800, actual: %d", resp.Data["max_ttl"])
	}

	// set default lease TTL to 10m; Token TTL should be capped at 10m
	// set max lease TTL to 20m; Token Max TTL should be capped at 20m
	config = &logical.BackendConfig{
		Logger: logging.NewVaultLogger(hclog.Trace),

		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Minute * 10,
			MaxLeaseTTLVal:     time.Minute * 20,
		},
		StorageView: &logical.InmemStorage{},
	}

	b, err = Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	roleReq.Operation = logical.CreateOperation

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	if resp.Data["token_ttl"].(int64) != 600 {
		t.Fatalf("bad: token_ttl; expected: 600, actual: %d", resp.Data["ttl"])
	}

	if resp.Data["token_max_ttl"].(int64) != 1200 {
		t.Fatalf("bad: token_max_ttl; expected: 1200, actual: %d", resp.Data["ttl"])
	}
}

func TestRoleEntryUpgradeV(t *testing.T) {
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

	roleEntryToUpgrade := &awsRoleEntry{
		BoundIamRoleARNs:            []string{"arn:aws:iam::123456789012:role/my_role_prefix"},
		BoundIamInstanceProfileARNs: []string{"arn:aws:iam::123456789012:instance-profile/my_profile-prefix"},
		Version:                     1,
	}
	expected := &awsRoleEntry{
		BoundIamRoleARNs:            []string{"arn:aws:iam::123456789012:role/my_role_prefix*"},
		BoundIamInstanceProfileARNs: []string{"arn:aws:iam::123456789012:instance-profile/my_profile-prefix*"},
		Version:                     currentRoleStorageVersion,
	}

	upgraded, err := b.upgradeRole(context.Background(), storage, roleEntryToUpgrade)
	if err != nil {
		t.Fatalf("error upgrading role entry: %#v", err)
	}
	if !upgraded {
		t.Fatalf("expected to upgrade role entry %#v but got no upgrade", roleEntryToUpgrade)
	}
	if roleEntryToUpgrade.RoleID == "" {
		t.Fatal("expected role ID to be populated")
	}
	expected.RoleID = roleEntryToUpgrade.RoleID
	if diff := deep.Equal(*roleEntryToUpgrade, *expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestRoleInitialize(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	err = b.Setup(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	// create some role entries, some of which will need to be upgraded
	type testData struct {
		name  string
		entry *awsRoleEntry
	}

	before := []testData{
		{
			name: "role1",
			entry: &awsRoleEntry{
				BoundIamRoleARNs:            []string{"arn:aws:iam::000000000001:role/my_role_prefix"},
				BoundIamInstanceProfileARNs: []string{"arn:aws:iam::000000000001:instance-profile/my_profile-prefix"},
				Version:                     1,
			},
		},
		{
			name: "role2",
			entry: &awsRoleEntry{
				BoundIamRoleARNs:            []string{"arn:aws:iam::000000000002:role/my_role_prefix"},
				BoundIamInstanceProfileARNs: []string{"arn:aws:iam::000000000002:instance-profile/my_profile-prefix"},
				Version:                     2,
			},
		},
		{
			name: "role3",
			entry: &awsRoleEntry{
				BoundIamRoleARNs:            []string{"arn:aws:iam::000000000003:role/my_role_prefix"},
				BoundIamInstanceProfileARNs: []string{"arn:aws:iam::000000000003:instance-profile/my_profile-prefix"},
				Version:                     currentRoleStorageVersion,
			},
		},
	}

	// put the entries in storage
	for _, role := range before {
		err = b.setRole(ctx, storage, role.name, role.entry)
		if err != nil {
			t.Fatal(err)
		}
	}

	// upgrade all the entries
	upgraded, err := b.upgrade(ctx, storage)
	if err != nil {
		t.Fatal(err)
	}
	if !upgraded {
		t.Fatalf("expected upgrade")
	}

	// read the entries from storage
	after := make([]testData, 0)
	names, err := storage.List(ctx, "role/")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		entry, err := b.role(ctx, storage, name)
		if err != nil {
			t.Fatal(err)
		}
		after = append(after, testData{name: name, entry: entry})
	}

	// make sure each entry is at the current version
	expected := []testData{
		{
			name: "role1",
			entry: &awsRoleEntry{
				BoundIamRoleARNs:            []string{"arn:aws:iam::000000000001:role/my_role_prefix"},
				BoundIamInstanceProfileARNs: []string{"arn:aws:iam::000000000001:instance-profile/my_profile-prefix"},
				Version:                     currentRoleStorageVersion,
			},
		},
		{
			name: "role2",
			entry: &awsRoleEntry{
				BoundIamRoleARNs:            []string{"arn:aws:iam::000000000002:role/my_role_prefix"},
				BoundIamInstanceProfileARNs: []string{"arn:aws:iam::000000000002:instance-profile/my_profile-prefix"},
				Version:                     currentRoleStorageVersion,
			},
		},
		{
			name: "role3",
			entry: &awsRoleEntry{
				BoundIamRoleARNs:            []string{"arn:aws:iam::000000000003:role/my_role_prefix"},
				BoundIamInstanceProfileARNs: []string{"arn:aws:iam::000000000003:instance-profile/my_profile-prefix"},
				Version:                     currentRoleStorageVersion,
			},
		},
	}
	if diff := deep.Equal(expected, after); diff != nil {
		t.Fatal(diff)
	}

	// run it again -- nothing will happen
	upgraded, err = b.upgrade(ctx, storage)
	if err != nil {
		t.Fatal(err)
	}
	if upgraded {
		t.Fatalf("expected no upgrade")
	}

	// make sure saved role version is correct
	entry, err := storage.Get(ctx, "config/version")
	if err != nil {
		t.Fatal(err)
	}
	var version awsVersion
	err = entry.DecodeJSON(&version)
	if err != nil {
		t.Fatal(err)
	}
	if version.Version != currentAwsVersion {
		t.Fatalf("expected version %d, got  %d", currentAwsVersion, version.Version)
	}

	// stomp on the saved version
	version.Version = 0
	e2, err := logical.StorageEntryJSON("config/version", version)
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Put(ctx, e2)
	if err != nil {
		t.Fatal(err)
	}

	// run it again -- now an upgrade will happen
	upgraded, err = b.upgrade(ctx, storage)
	if err != nil {
		t.Fatal(err)
	}
	if !upgraded {
		t.Fatalf("expected upgrade")
	}
}

func TestAwsVersion(t *testing.T) {
	before := awsVersion{
		Version: 42,
	}

	entry, err := logical.StorageEntryJSON("config/version", &before)
	if err != nil {
		t.Fatal(err)
	}

	var after awsVersion
	err = entry.DecodeJSON(&after)
	if err != nil {
		t.Fatal(err)
	}

	if diff := deep.Equal(before, after); diff != nil {
		t.Fatal(diff)
	}
}

// This test was used to reproduce https://github.com/hashicorp/vault/issues/7418
// and verify its fix.
// Please run it at least 3 times to ensure that passing tests are due to actually
// passing, rather than the region being randomly chosen tying to the one in the
// test through luck.
func TestRoleResolutionWithSTSEndpointConfigured(t *testing.T) {
	if enabled := os.Getenv(vlttesting.TestEnvVar); enabled == "" {
		t.Skip()
	}

	/* ARN of an AWS role that Vault can query during testing.
	   This role should exist in your current AWS account and your credentials
	   should have iam:GetRole permissions to query it.
	*/
	assumableRoleArn := os.Getenv("AWS_ASSUMABLE_ROLE_ARN")
	if assumableRoleArn == "" {
		t.Skip("skipping because AWS_ASSUMABLE_ROLE_ARN is unset")
	}

	// Ensure aws credentials are available locally for testing.
	logger := logging.NewVaultLogger(hclog.Debug)
	credsConfig := &awsutil.CredentialsConfig{Logger: logger}
	credsChain, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		t.Fatal(err)
	}
	_, err = credsChain.Get()
	if err != nil {
		t.SkipNow()
	}

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

	// configure the client with an sts endpoint that should be used in creating the role
	data := map[string]interface{}{
		"sts_endpoint": "https://sts.eu-west-1.amazonaws.com",
		// Note - if you comment this out, you can reproduce the error shown
		// in the linked GH issue above. This essentially reproduces the problem
		// we had when we didn't have an sts_region field.
		"sts_region": "eu-west-1",
	}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config/client",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create the role entry; resp: %#v", resp)
	}

	data = map[string]interface{}{
		"auth_type":               iamAuthType,
		"bound_iam_principal_arn": assumableRoleArn,
		"resolve_aws_unique_ids":  true,
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/MyRoleName",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create the role entry; resp: %#v", resp)
	}
}

func resolveArnToFakeUniqueId(_ context.Context, _ logical.Storage, _ string) (string, error) {
	return "FakeUniqueId1", nil
}
