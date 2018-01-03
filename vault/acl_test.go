package vault

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestACL_Capabilities(t *testing.T) {
	// Create the root policy ACL
	policy := []*Policy{&Policy{Name: "root"}}
	acl, err := NewACL(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	actual := acl.Capabilities("any/path")
	expected := []string{"root"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	policies, err := ParseACLPolicy(aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err = NewACL([]*Policy{policies})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	actual = acl.Capabilities("dev")
	expected = []string{"deny"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: path:%s\ngot\n%#v\nexpected\n%#v\n", "deny", actual, expected)
	}

	actual = acl.Capabilities("dev/")
	expected = []string{"sudo", "read", "list", "update", "delete", "create"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: path:%s\ngot\n%#v\nexpected\n%#v\n", "dev/", actual, expected)
	}

	actual = acl.Capabilities("stage/aws/test")
	expected = []string{"sudo", "read", "list", "update"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: path:%s\ngot\n%#v\nexpected\n%#v\n", "stage/aws/test", actual, expected)
	}

}

func TestACL_Root(t *testing.T) {
	// Create the root policy ACL
	policy := []*Policy{&Policy{Name: "root"}}
	acl, err := NewACL(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	request := new(logical.Request)
	request.Operation = logical.UpdateOperation
	request.Path = "sys/mount/foo"
	authResults := acl.AllowOperation(request)
	if !authResults.RootPrivs {
		t.Fatalf("expected root")
	}
	if !authResults.Allowed {
		t.Fatalf("expected permissions")
	}
}

func TestACL_Single(t *testing.T) {
	policy, err := ParseACLPolicy(aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := NewACL([]*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Type of operation is not important here as we only care about checking
	// sudo/root
	request := new(logical.Request)
	request.Operation = logical.ReadOperation
	request.Path = "sys/mount/foo"
	authResults := acl.AllowOperation(request)
	if authResults.RootPrivs {
		t.Fatalf("unexpected root")
	}

	type tcase struct {
		op        logical.Operation
		path      string
		allowed   bool
		rootPrivs bool
	}
	tcases := []tcase{
		{logical.ReadOperation, "root", false, false},
		{logical.HelpOperation, "root", true, false},

		{logical.ReadOperation, "dev/foo", true, true},
		{logical.UpdateOperation, "dev/foo", true, true},

		{logical.DeleteOperation, "stage/foo", true, false},
		{logical.ListOperation, "stage/aws/foo", true, true},
		{logical.UpdateOperation, "stage/aws/foo", true, true},
		{logical.UpdateOperation, "stage/aws/policy/foo", true, true},

		{logical.DeleteOperation, "prod/foo", false, false},
		{logical.UpdateOperation, "prod/foo", false, false},
		{logical.ReadOperation, "prod/foo", true, false},
		{logical.ListOperation, "prod/foo", true, false},
		{logical.ReadOperation, "prod/aws/foo", false, false},

		{logical.ReadOperation, "foo/bar", true, true},
		{logical.ListOperation, "foo/bar", false, true},
		{logical.UpdateOperation, "foo/bar", false, true},
		{logical.CreateOperation, "foo/bar", true, true},
	}

	for _, tc := range tcases {
		request := new(logical.Request)
		request.Operation = tc.op
		request.Path = tc.path
		authResults := acl.AllowOperation(request)
		if authResults.Allowed != tc.allowed {
			t.Fatalf("bad: case %#v: %v, %v", tc, authResults.Allowed, authResults.RootPrivs)
		}
		if authResults.RootPrivs != tc.rootPrivs {
			t.Fatalf("bad: case %#v: %v, %v", tc, authResults.Allowed, authResults.RootPrivs)
		}
	}
}

func TestACL_Layered(t *testing.T) {
	policy1, err := ParseACLPolicy(aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy2, err := ParseACLPolicy(aclPolicy2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := NewACL([]*Policy{policy1, policy2})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	testLayeredACL(t, acl)
}

func testLayeredACL(t *testing.T, acl *ACL) {
	// Type of operation is not important here as we only care about checking
	// sudo/root
	request := new(logical.Request)
	request.Operation = logical.ReadOperation
	request.Path = "sys/mount/foo"
	authResults := acl.AllowOperation(request)
	if authResults.RootPrivs {
		t.Fatalf("unexpected root")
	}

	type tcase struct {
		op        logical.Operation
		path      string
		allowed   bool
		rootPrivs bool
	}
	tcases := []tcase{
		{logical.ReadOperation, "root", false, false},
		{logical.HelpOperation, "root", true, false},

		{logical.ReadOperation, "dev/foo", true, true},
		{logical.UpdateOperation, "dev/foo", true, true},
		{logical.ReadOperation, "dev/hide/foo", false, false},
		{logical.UpdateOperation, "dev/hide/foo", false, false},

		{logical.DeleteOperation, "stage/foo", true, false},
		{logical.ListOperation, "stage/aws/foo", true, true},
		{logical.UpdateOperation, "stage/aws/foo", true, true},
		{logical.UpdateOperation, "stage/aws/policy/foo", false, false},

		{logical.DeleteOperation, "prod/foo", true, false},
		{logical.UpdateOperation, "prod/foo", true, false},
		{logical.ReadOperation, "prod/foo", true, false},
		{logical.ListOperation, "prod/foo", true, false},
		{logical.ReadOperation, "prod/aws/foo", false, false},

		{logical.ReadOperation, "sys/status", false, false},
		{logical.UpdateOperation, "sys/seal", true, true},

		{logical.ReadOperation, "foo/bar", false, false},
		{logical.ListOperation, "foo/bar", false, false},
		{logical.UpdateOperation, "foo/bar", false, false},
		{logical.CreateOperation, "foo/bar", false, false},
	}

	for _, tc := range tcases {
		request := new(logical.Request)
		request.Operation = tc.op
		request.Path = tc.path
		authResults := acl.AllowOperation(request)
		if authResults.Allowed != tc.allowed {
			t.Fatalf("bad: case %#v: %v, %v", tc, authResults.Allowed, authResults.RootPrivs)
		}
		if authResults.RootPrivs != tc.rootPrivs {
			t.Fatalf("bad: case %#v: %v, %v", tc, authResults.Allowed, authResults.RootPrivs)
		}
	}
}

func TestACL_PolicyMerge(t *testing.T) {
	policy, err := ParseACLPolicy(mergingPolicies)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	acl, err := NewACL([]*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	type tcase struct {
		path           string
		minWrappingTTL *time.Duration
		maxWrappingTTL *time.Duration
		allowed        map[string][]interface{}
		denied         map[string][]interface{}
		required       []string
	}

	createDuration := func(seconds int) *time.Duration {
		ret := time.Duration(seconds) * time.Second
		return &ret
	}

	tcases := []tcase{
		{"foo/bar", nil, nil, nil, map[string][]interface{}{"zip": []interface{}{}, "baz": []interface{}{}}, []string{"baz"}},
		{"hello/universe", createDuration(50), createDuration(200), map[string][]interface{}{"foo": []interface{}{}, "bar": []interface{}{}}, nil, []string{"foo", "bar"}},
		{"allow/all", nil, nil, map[string][]interface{}{"*": []interface{}{}, "test": []interface{}{}, "test1": []interface{}{"foo"}}, nil, nil},
		{"allow/all1", nil, nil, map[string][]interface{}{"*": []interface{}{}, "test": []interface{}{}, "test1": []interface{}{"foo"}}, nil, nil},
		{"deny/all", nil, nil, nil, map[string][]interface{}{"*": []interface{}{}, "test": []interface{}{}}, nil},
		{"deny/all1", nil, nil, nil, map[string][]interface{}{"*": []interface{}{}, "test": []interface{}{}}, nil},
		{"value/merge", nil, nil, map[string][]interface{}{"test": []interface{}{3, 4, 1, 2}}, map[string][]interface{}{"test": []interface{}{3, 4, 1, 2}}, nil},
		{"value/empty", nil, nil, map[string][]interface{}{"empty": []interface{}{}}, map[string][]interface{}{"empty": []interface{}{}}, nil},
	}

	for _, tc := range tcases {
		raw, ok := acl.exactRules.Get(tc.path)
		if !ok {
			t.Fatalf("Could not find acl entry for path %s", tc.path)
		}

		p := raw.(*ACLPermissions)
		if !reflect.DeepEqual(tc.allowed, p.AllowedParameters) {
			t.Fatalf("Allowed paramaters did not match, Expected: %#v, Got: %#v", tc.allowed, p.AllowedParameters)
		}
		if !reflect.DeepEqual(tc.denied, p.DeniedParameters) {
			t.Fatalf("Denied paramaters did not match, Expected: %#v, Got: %#v", tc.denied, p.DeniedParameters)
		}
		if !reflect.DeepEqual(tc.required, p.RequiredParameters) {
			t.Fatalf("Required paramaters did not match, Expected: %#v, Got: %#v", tc.required, p.RequiredParameters)
		}
		if tc.minWrappingTTL != nil && *tc.minWrappingTTL != p.MinWrappingTTL {
			t.Fatalf("Min wrapping TTL did not match, Expected: %#v, Got: %#v", tc.minWrappingTTL, p.MinWrappingTTL)
		}
		if tc.minWrappingTTL != nil && *tc.maxWrappingTTL != p.MaxWrappingTTL {
			t.Fatalf("Max wrapping TTL did not match, Expected: %#v, Got: %#v", tc.maxWrappingTTL, p.MaxWrappingTTL)
		}
	}
}

func TestACL_AllowOperation(t *testing.T) {
	policy, err := ParseACLPolicy(permissionsPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	acl, err := NewACL([]*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	toperations := []logical.Operation{
		logical.UpdateOperation,
		logical.CreateOperation,
	}
	type tcase struct {
		path        string
		wrappingTTL *time.Duration
		parameters  []string
		allowed     bool
	}

	createDuration := func(seconds int) *time.Duration {
		ret := time.Duration(seconds) * time.Second
		return &ret
	}

	tcases := []tcase{
		{"dev/ops", nil, []string{"zip"}, true},
		{"foo/bar", nil, []string{"zap"}, false},
		{"foo/bar", nil, []string{"zip"}, false},
		{"foo/bar", createDuration(50), []string{"zip"}, false},
		{"foo/bar", createDuration(450), []string{"zip"}, false},
		{"foo/bar", createDuration(350), []string{"zip"}, true},
		{"foo/baz", nil, []string{"hello"}, false},
		{"foo/baz", createDuration(50), []string{"hello"}, false},
		{"foo/baz", createDuration(450), []string{"hello"}, true},
		{"foo/baz", nil, []string{"zap"}, false},
		{"broken/phone", nil, []string{"steve"}, false},
		{"working/phone", nil, []string{""}, false},
		{"working/phone", createDuration(450), []string{""}, false},
		{"working/phone", createDuration(350), []string{""}, true},
		{"hello/world", nil, []string{"one"}, false},
		{"tree/fort", nil, []string{"one"}, true},
		{"tree/fort", nil, []string{"foo"}, false},
		{"fruit/apple", nil, []string{"pear"}, false},
		{"fruit/apple", nil, []string{"one"}, false},
		{"cold/weather", nil, []string{"four"}, true},
		{"var/aws", nil, []string{"cold", "warm", "kitty"}, false},
		{"var/req", nil, []string{"cold", "warm", "kitty"}, false},
		{"var/req", nil, []string{"cold", "warm", "kitty", "foo"}, true},
	}

	for _, tc := range tcases {
		request := logical.Request{Path: tc.path, Data: make(map[string]interface{})}
		for _, parameter := range tc.parameters {
			request.Data[parameter] = ""
		}
		if tc.wrappingTTL != nil {
			request.WrapInfo = &logical.RequestWrapInfo{
				TTL: *tc.wrappingTTL,
			}
		}
		for _, op := range toperations {
			request.Operation = op
			authResults := acl.AllowOperation(&request)
			if authResults.Allowed != tc.allowed {
				t.Fatalf("bad: case %#v: %v", tc, authResults.Allowed)
			}
		}
	}
}

func TestACL_ValuePermissions(t *testing.T) {
	policy, err := ParseACLPolicy(valuePermissionsPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	acl, err := NewACL([]*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	toperations := []logical.Operation{
		logical.UpdateOperation,
		logical.CreateOperation,
	}
	type tcase struct {
		path       string
		parameters []string
		values     []interface{}
		allowed    bool
	}

	tcases := []tcase{
		{"dev/ops", []string{"allow"}, []interface{}{"good"}, true},
		{"dev/ops", []string{"allow"}, []interface{}{"bad"}, false},
		{"foo/bar", []string{"deny"}, []interface{}{"bad"}, false},
		{"foo/bar", []string{"deny"}, []interface{}{"bad glob"}, false},
		{"foo/bar", []string{"deny"}, []interface{}{"good"}, true},
		{"foo/bar", []string{"allow"}, []interface{}{"good"}, true},
		{"foo/baz", []string{"aLLow"}, []interface{}{"good"}, true},
		{"foo/baz", []string{"deny"}, []interface{}{"bad"}, false},
		{"foo/baz", []string{"deny"}, []interface{}{"good"}, false},
		{"foo/baz", []string{"allow", "deny"}, []interface{}{"good", "bad"}, false},
		{"foo/baz", []string{"deny", "allow"}, []interface{}{"good", "bad"}, false},
		{"foo/baz", []string{"deNy", "allow"}, []interface{}{"bad", "good"}, false},
		{"foo/baz", []string{"aLLow"}, []interface{}{"bad"}, false},
		{"foo/baz", []string{"Neither"}, []interface{}{"bad"}, false},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"good"}, true},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"good1"}, true},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"good2"}, true},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"glob good2"}, false},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"glob good3"}, true},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"bad"}, false},
		{"fizz/buzz", []string{"allow_multi"}, []interface{}{"bad"}, false},
		{"fizz/buzz", []string{"allow_multi", "allow"}, []interface{}{"good1", "good"}, true},
		{"fizz/buzz", []string{"deny_multi"}, []interface{}{"bad2"}, false},
		{"fizz/buzz", []string{"deny_multi", "allow_multi"}, []interface{}{"good", "good2"}, false},
		//	{"test/types", []string{"array"}, []interface{}{[1]string{"good"}}, true},
		{"test/types", []string{"map"}, []interface{}{map[string]interface{}{"good": "one"}}, true},
		{"test/types", []string{"map"}, []interface{}{map[string]interface{}{"bad": "one"}}, false},
		{"test/types", []string{"int"}, []interface{}{1}, true},
		{"test/types", []string{"int"}, []interface{}{3}, false},
		{"test/types", []string{"bool"}, []interface{}{false}, true},
		{"test/types", []string{"bool"}, []interface{}{true}, false},
		{"test/star", []string{"anything"}, []interface{}{true}, true},
		{"test/star", []string{"foo"}, []interface{}{true}, true},
		{"test/star", []string{"bar"}, []interface{}{false}, true},
		{"test/star", []string{"bar"}, []interface{}{true}, false},
	}

	for _, tc := range tcases {
		request := logical.Request{Path: tc.path, Data: make(map[string]interface{})}
		for i, parameter := range tc.parameters {
			request.Data[parameter] = tc.values[i]
		}
		for _, op := range toperations {
			request.Operation = op
			authResults := acl.AllowOperation(&request)
			if authResults.Allowed != tc.allowed {
				t.Fatalf("bad: case %#v: %v", tc, authResults.Allowed)
			}
		}
	}
}

// NOTE: this test doesn't catch any races ATM
func TestACL_CreationRace(t *testing.T) {
	policy, err := ParseACLPolicy(valuePermissionsPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	var wg sync.WaitGroup
	stopTime := time.Now().Add(20 * time.Second)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if time.Now().After(stopTime) {
					return
				}
				_, err := NewACL([]*Policy{policy})
				if err != nil {
					t.Fatalf("err: %v", err)
				}
			}
		}()
	}

	wg.Wait()
}

var tokenCreationPolicy = `
name = "tokenCreation"
path "auth/token/create*" {
	capabilities = ["update", "create", "sudo"]
}
`

var aclPolicy = `
name = "DeV"
path "dev/*" {
	policy = "sudo"
}
path "stage/*" {
	policy = "write"
}
path "stage/aws/*" {
	policy = "read"
	capabilities = ["update", "sudo"]
}
path "stage/aws/policy/*" {
	policy = "sudo"
}
path "prod/*" {
	policy = "read"
}
path "prod/aws/*" {
	policy = "deny"
}
path "sys/*" {
	policy = "deny"
}
path "foo/bar" {
	capabilities = ["read", "create", "sudo"]
}
`

var aclPolicy2 = `
name = "OpS"
path "dev/hide/*" {
	policy = "deny"
}
path "stage/aws/policy/*" {
	policy = "deny"
	# This should have no effect
	capabilities = ["read", "update", "sudo"]
}
path "prod/*" {
	policy = "write"
}
path "sys/seal" {
	policy = "sudo"
}
path "foo/bar" {
	capabilities = ["deny"]
}
`

//test merging
var mergingPolicies = `
name = "ops"
path "foo/bar" {
	policy = "write"
	denied_parameters = {
		"baz" = []
	}
	required_parameters = ["baz"]
}
path "foo/bar" {
	policy = "write"
	denied_parameters = {
		"zip" = []
	}
}
path "hello/universe" {
	policy = "write"
	allowed_parameters = {
		"foo" = []
	}
	required_parameters = ["foo"]
	max_wrapping_ttl = 300
	min_wrapping_ttl = 100
}
path "hello/universe" {
	policy = "write"
	allowed_parameters = {
		"bar" = []
	}
	required_parameters = ["bar"]
	max_wrapping_ttl = 200
	min_wrapping_ttl = 50
}
path "allow/all" {
	policy = "write"
	allowed_parameters = {
		"test" = []
		"test1" = ["foo"]
	}
}
path "allow/all" {
	policy = "write"
	allowed_parameters = {
		"*" = []
	}
}
path "allow/all1" {
	policy = "write"
	allowed_parameters = {
		"*" = []
	}
}
path "allow/all1" {
	policy = "write"
	allowed_parameters = {
		"test" = []
		"test1" = ["foo"]
	}
}
path "deny/all" {
	policy = "write"
	denied_parameters = {
		"test" = []
	}
}
path "deny/all" {
	policy = "write"
	denied_parameters = {
		"*" = []
	}
}
path "deny/all1" {
	policy = "write"
	denied_parameters = {
		"*" = []
	}
}
path "deny/all1" {
	policy = "write"
	denied_parameters = {
		"test" = []
	}
}
path "value/merge" {
	policy = "write"
	allowed_parameters = {
		"test" = [1, 2]
	}
	denied_parameters = {
		"test" = [1, 2]
	}
}
path "value/merge" {
	policy = "write"
	allowed_parameters = {
		"test" = [3, 4]
	}
	denied_parameters = {
		"test" = [3, 4]
	}
}
path "value/empty" {
	policy = "write"
	allowed_parameters = {
		"empty" = []
	}
	denied_parameters = {
		"empty" = [1]
	}
}
path "value/empty" {
	policy = "write"
	allowed_parameters = {
		"empty" = [1]
	}
	denied_parameters = {
		"empty" = []
	}
}
`

//allow operation testing
var permissionsPolicy = `
name = "dev"
path "dev/*" {
	policy = "write"
	
	allowed_parameters = {
		"zip" = []
	}
}
path "foo/bar" {
	policy = "write"
	denied_parameters = {
		"zap" = []
	}
	min_wrapping_ttl = 300
	max_wrapping_ttl = 400
}
path "foo/baz" {
	policy = "write"
	allowed_parameters = {
		"hello" = []
	}
	denied_parameters = {
		"zap" = []
	}
	min_wrapping_ttl = 300
}
path "working/phone" {
	policy = "write"
	max_wrapping_ttl = 400
}
path "broken/phone" {
	policy = "write"
	allowed_parameters = {
	  "steve" = []
	}
	denied_parameters = {
	  "steve" = []
	}
}
path "hello/world" {
	policy = "write"
	allowed_parameters = {
		"*" = []
	}
	denied_parameters = {
		"*" = []
	}
}
path "tree/fort" {
	policy = "write"
	allowed_parameters = {
		"*" = []
	}
	denied_parameters = {
		"foo" = []
	}
}
path "fruit/apple" {
	policy = "write"
	allowed_parameters = {
		"pear" = []
	}
	denied_parameters = {
		"*" = []
	}
}
path "cold/weather" {
	policy = "write"
	allowed_parameters = {}
	denied_parameters = {}
}
path "var/aws" {
	policy = "write"
	allowed_parameters = {
		"*" = []
	}
	denied_parameters = {
		"soft" = []
		"warm" = []
		"kitty" = []
	}
}
path "var/req" {
	policy = "write"
	required_parameters = ["foo"]
}
`

//allow operation testing
var valuePermissionsPolicy = `
name = "op"
path "dev/*" {
	policy = "write"
	
	allowed_parameters = {
		"allow" = ["good"]
	}
}
path "foo/bar" {
	policy = "write"
	denied_parameters = {
		"deny" = ["bad*"]
	}
}
path "foo/baz" {
	policy = "write"
	allowed_parameters = {
		"ALLOW" = ["good"]
	}
	denied_parameters = {
		"dEny" = ["bad"]
	}
}
path "fizz/buzz" {
	policy = "write"
	allowed_parameters = {
		"allow_multi" = ["good", "good1", "good2", "*good3"]
		"allow" = ["good"]
	}
	denied_parameters = {
		"deny_multi" = ["bad", "bad1", "bad2"]
	}
}
path "test/types" {
	policy = "write"
	allowed_parameters = {
		"map" = [{"good" = "one"}]
		"int" = [1, 2]
		"bool" = [false]
	}
	denied_parameters = {
	}
}
path "test/star" {
	policy = "write"
	allowed_parameters = {
		"*" = []
		"foo" = []
		"bar" = [false]
	}
	denied_parameters = {
	}
}
`
