package awsauth

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
)

func TestBackend_pathRoleEc2(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"auth_type":    "ec2",
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
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || !resp.IsError() {
		t.Fatalf("expected failure to create role with both allow_instance_migration true and disallow_reauthentication true")
	}
	data["disallow_reauthentication"] = false
	resp, err = b.HandleRequest(&logical.Request{
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
	resp, err = b.HandleRequest(&logical.Request{
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
	resp, err = b.HandleRequest(&logical.Request{
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

func Test_enableIamIDResolution(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}
	roleName := "upgradable_role"

	b.resolveArnToUniqueIDFunc = resolveArnToFakeUniqueId

	data := map[string]interface{}{
		"auth_type":               iamAuthType,
		"policies":                "p,q",
		"bound_iam_principal_arn": "arn:aws:iam::123456789012:role/MyRole",
		"resolve_aws_unique_ids":  false,
	}

	submitRequest := func(roleName string, op logical.Operation) (*logical.Response, error) {
		return b.HandleRequest(&logical.Request{
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
	if resp.Data["bound_iam_principal_id"] != "" {
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
	if resp.Data["bound_iam_principal_id"] != "FakeUniqueId1" {
		t.Fatalf("bad: expected upgrade of role resolve principal ID to %q, but got %q instead", "FakeUniqueId1", resp.Data["bound_iam_principal_id"])
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
	err = b.Setup(config)
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
		"auth_type":               iamAuthType,
		"policies":                "p,q,r,s",
		"max_ttl":                 "2h",
		"bound_iam_principal_arn": "n:aws:iam::123456789012:user/MyUserName",
		"resolve_aws_unique_ids":  false,
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
		t.Fatalf("failed to create the role entry; resp: %#v", resp)
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

	data["inferred_entity_type"] = "invalid"
	resp, err = b.HandleRequest(&logical.Request{
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
	secondRole := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "role/MyOtherRoleName",
		Data:      data,
		Storage:   storage,
	}
	resp, err = b.HandleRequest(secondRole)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create additional role: %v", *secondRole)
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
	err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"policies":     "p,q,r,s",
		"bound_ami_id": "ami-abc1234",
		"auth_type":    "ec2,invalid",
	}

	submitRequest := func(roleName string, op logical.Operation) (*logical.Response, error) {
		return b.HandleRequest(&logical.Request{
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
	data["bound_iam_principal_arn"] = "arn:aws:iam::123456789012:role/MyRole"
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
	if resp.Data["bound_iam_principal_id"] != "FakeUniqueId1" {
		t.Fatalf("expected fake unique ID of FakeUniqueId1, got %q", resp.Data["bound_iam_principal_id"])
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
	err = b.Setup(config)
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

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleData := map[string]interface{}{
		"auth_type":                      "ec2",
		"bound_ami_id":                   "testamiid",
		"bound_account_id":               "testaccountid",
		"bound_region":                   "testregion",
		"bound_iam_role_arn":             "arn:aws:iam::123456789012:role/MyRole",
		"bound_iam_instance_profile_arn": "arn:aws:iam::123456789012:instance-profile/MyInstanceProfile",
		"bound_subnet_id":                "testsubnetid",
		"bound_vpc_id":                   "testvpcid",
		"role_tag":                       "testtag",
		"resolve_aws_unique_ids":         false,
		"allow_instance_migration":       true,
		"ttl":                       "10m",
		"max_ttl":                   "20m",
		"policies":                  "testpolicy1,testpolicy2",
		"disallow_reauthentication": false,
		"hmac_key":                  "testhmackey",
		"period":                    "1m",
	}

	roleReq.Path = "role/testrole"
	roleReq.Data = roleData
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	expected := map[string]interface{}{
		"auth_type":                      ec2AuthType,
		"bound_ami_id":                   "testamiid",
		"bound_account_id":               "testaccountid",
		"bound_region":                   "testregion",
		"bound_iam_principal_arn":        "",
		"bound_iam_principal_id":         "",
		"bound_iam_role_arn":             "arn:aws:iam::123456789012:role/MyRole",
		"bound_iam_instance_profile_arn": "arn:aws:iam::123456789012:instance-profile/MyInstanceProfile",
		"bound_subnet_id":                "testsubnetid",
		"bound_vpc_id":                   "testvpcid",
		"inferred_entity_type":           "",
		"inferred_aws_region":            "",
		"resolve_aws_unique_ids":         false,
		"role_tag":                       "testtag",
		"allow_instance_migration":       true,
		"ttl":                       time.Duration(600),
		"max_ttl":                   time.Duration(1200),
		"policies":                  []string{"testpolicy1", "testpolicy2"},
		"disallow_reauthentication": false,
		"period":                    time.Duration(60),
	}

	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: role data: expected: %#v\n actual: %#v", expected, resp.Data)
	}

	roleData["bound_vpc_id"] = "newvpcid"
	roleReq.Operation = logical.UpdateOperation

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	expected["bound_vpc_id"] = "newvpcid"

	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: role data: expected: %#v\n actual: %#v", expected, resp.Data)
	}

	roleReq.Operation = logical.DeleteOperation

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	if resp != nil {
		t.Fatalf("failed to delete role entry")
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
	err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	roleData := map[string]interface{}{
		"auth_type":                      "ec2",
		"bound_iam_instance_profile_arn": "arn:aws:iam::123456789012:instance-profile/test-profile-name",
		"resolve_aws_unique_ids":         false,
		"ttl":     "10s",
		"max_ttl": "20s",
		"period":  "30s",
	}

	roleReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   storage,
		Path:      "role/testrole",
		Data:      roleData,
	}

	resp, err := b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v, err: %v", resp, err)
	}

	if int64(resp.Data["ttl"].(time.Duration)) != 10 {
		t.Fatalf("bad: period; expected: 10, actual: %d", resp.Data["ttl"])
	}
	if int64(resp.Data["max_ttl"].(time.Duration)) != 20 {
		t.Fatalf("bad: period; expected: 20, actual: %d", resp.Data["max_ttl"])
	}
	if int64(resp.Data["period"].(time.Duration)) != 30 {
		t.Fatalf("bad: period; expected: 30, actual: %d", resp.Data["period"])
	}
}

func resolveArnToFakeUniqueId(s logical.Storage, arn string) (string, error) {
	return "FakeUniqueId1", nil
}
