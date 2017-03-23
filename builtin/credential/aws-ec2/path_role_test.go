package awsec2

import (
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestAwsEc2_RoleCrud(t *testing.T) {
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

	roleData := map[string]interface{}{
		"bound_ami_id":                   "testamiid",
		"bound_account_id":               "testaccountid",
		"bound_region":                   "testregion",
		"bound_iam_role_arn":             "testiamrolearn",
		"bound_iam_instance_profile_arn": "testiaminstanceprofilearn",
		"bound_subnet_id":                "testsubnetid",
		"bound_vpc_id":                   "testvpcid",
		"role_tag":                       "testtag",
		"allow_instance_migration":       true,
		"ttl":                       "10m",
		"max_ttl":                   "20m",
		"policies":                  "testpolicy1,testpolicy2",
		"disallow_reauthentication": true,
		"hmac_key":                  "testhmackey",
		"period":                    "1m",
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
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

	expected := map[string]interface{}{
		"bound_ami_id":                   "testamiid",
		"bound_account_id":               "testaccountid",
		"bound_region":                   "testregion",
		"bound_iam_role_arn":             "testiamrolearn",
		"bound_iam_instance_profile_arn": "testiaminstanceprofilearn",
		"bound_subnet_id":                "testsubnetid",
		"bound_vpc_id":                   "testvpcid",
		"role_tag":                       "testtag",
		"allow_instance_migration":       true,
		"ttl":                       time.Duration(600),
		"max_ttl":                   time.Duration(1200),
		"policies":                  []string{"default", "testpolicy1", "testpolicy2"},
		"disallow_reauthentication": true,
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
	_, err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	roleData := map[string]interface{}{
		"bound_iam_instance_profile_arn": "testarn",
		"ttl":     "10s",
		"max_ttl": "20s",
		"period":  "30s",
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
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
