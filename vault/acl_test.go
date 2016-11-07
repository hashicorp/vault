package vault

import (
	"reflect"
	"testing"

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

	policies, err := Parse(aclPolicy)
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
	allowed, rootPrivs := acl.AllowOperation(request)
	if !rootPrivs {
		t.Fatalf("expected root")
	}
	if !allowed {
		t.Fatalf("expected permissions")
	}
}

func TestACL_Single(t *testing.T) {
	policy, err := Parse(aclPolicy)
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
	_, rootPrivs := acl.AllowOperation(request)
	if rootPrivs {
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
		allowed, rootPrivs := acl.AllowOperation(request)
		if allowed != tc.allowed {
			t.Fatalf("bad: case %#v: %v, %v", tc, allowed, rootPrivs)
		}
		if rootPrivs != tc.rootPrivs {
			t.Fatalf("bad: case %#v: %v, %v", tc, allowed, rootPrivs)
		}
	}
}

func TestACL_Layered(t *testing.T) {
	policy1, err := Parse(aclPolicy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	policy2, err := Parse(aclPolicy2)
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
	_, rootPrivs := acl.AllowOperation(request)
	if rootPrivs {
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
		allowed, rootPrivs := acl.AllowOperation(request)
		if allowed != tc.allowed {
			t.Fatalf("bad: case %#v: %v, %v", tc, allowed, rootPrivs)
		}
		if rootPrivs != tc.rootPrivs {
			t.Fatalf("bad: case %#v: %v, %v", tc, allowed, rootPrivs)
		}
	}
}

func TestPolicyMerge(t *testing.T) {
	policy, err := Parse(permissionsPolicy2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	acl, err := NewACL([]*Policy{policy})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	type tcase struct {
		op        logical.Operation
		pathName  string
		parameter string
		allowed   bool
		rootPrivs bool
	}

	tcases := []tcase{
		{logical.UpdateOperation, "foo/bar", "baz", false, false},
		{logical.UpdateOperation, "foo/bar", "zip", false, false},
		{logical.CreateOperation, "foo/bar", "baz", false, false},
		{logical.CreateOperation, "foo/bar", "zip", false, false},
		{logical.DeleteOperation, "foo/bar", "baz", false, false},
		{logical.DeleteOperation, "foo/bar", "zip", false, false},
		{logical.UpdateOperation, "hello/universe", "bob", true, false},
		{logical.UpdateOperation, "hello/universe", "tom", true, false},
		{logical.CreateOperation, "hello/universe", "bob", true, false},
		{logical.CreateOperation, "hello/universe", "tom", true, false},
		{logical.DeleteOperation, "hello/universe", "bob", true, false},
		{logical.DeleteOperation, "hello/universe", "tom", true, false},
		{logical.UpdateOperation, "rainy/day", "bob", true, false},
		{logical.UpdateOperation, "rainy/day", "tom", true, false},
		{logical.CreateOperation, "rainy/day", "bob", true, false},
		{logical.CreateOperation, "rainy/day", "tom", true, false},
		{logical.DeleteOperation, "rainy/day", "bob", true, false},
		{logical.DeleteOperation, "rainy/day", "tom", true, false},
		{logical.UpdateOperation, "cool/bike", "frank", false, false},
		{logical.UpdateOperation, "cool/bike", "two", false, false},
		{logical.CreateOperation, "cool/bike", "frank", false, false},
		{logical.CreateOperation, "cool/bike", "four", false, false},
		{logical.DeleteOperation, "cool/bike", "frank", false, false},
		{logical.DeleteOperation, "cool/bike", "six", false, false},
		{logical.UpdateOperation, "clean/bed", "one", false, false},
		{logical.UpdateOperation, "clean/bed", "two", false, false},
		{logical.CreateOperation, "clean/bed", "three", false, false},
		{logical.CreateOperation, "clean/bed", "four", false, false},
		{logical.DeleteOperation, "clean/bed", "five", false, false},
		{logical.DeleteOperation, "clean/bed", "six", false, false},
		{logical.UpdateOperation, "coca/cola", "john", false, false},
		{logical.UpdateOperation, "coca/cola", "two", false, false},
		{logical.CreateOperation, "coca/cola", "john", false, false},
		{logical.CreateOperation, "coca/cola", "four", false, false},
		{logical.DeleteOperation, "coca/cola", "john", false, false},
		{logical.DeleteOperation, "coca/cola", "six", false, false},
	}

	for _, tc := range tcases {
		request := new(logical.Request)
		request.Operation = tc.op
		request.Path = tc.pathName
		paramMap := make(map[string]interface{})
		paramMap[tc.parameter] = ""
		request.Data = paramMap
		allowed, rootPrivs := acl.AllowOperation(request)
		if allowed != tc.allowed {
			t.Fatalf("bad: case %#v: %v, %v", tc, allowed, rootPrivs)
		}
		if rootPrivs != tc.rootPrivs {
			t.Fatalf("bad: case %#v: %v, %v", tc, allowed, rootPrivs)
		}
	}

}

//test merging

var permissionsPolicy2 = `
name = "ops"
path "foo/bar" {
	policy = "write"
	permissions = {
		deniedparameters = {
			"baz" = {}
		}
	}
}
path "foo/bar" {
	policy = "write"
	permissions = {
		deniedparameters = {
			"zip" = {}
		}
	}
}
path "hello/universe" {
	policy = "write"
	permissions = {
		allowedparameters = {
			"bob" = {}
		}
	}
}
path "hello/universe" {
	policy = "write"
	permissions = {
		allowedparameters = {
			"tom" = {}
		}
  }
}
path "rainy/day" {
	policy = "write"
	permissions = {
		allowedparameters = {
			"bob" = {}
		}
	}
}
path "rainy/day" {
	policy = "write"
	permissions = {
		allowedparameters = {
			"*" = {}
		}
  }
}
path "cool/bike" {
	policy = "write"
	permissions = {
		deniedparameters = {
			"frank" = {}
		}
	}
}
path "cool/bike" {
	policy = "write"
	permissions = {
		deniedparameters = {
			"*" = {}
		}
  }
}
path "clean/bed" {
	policy = "write"
	permissions = {
		deniedparameters = {
			"*" = {}
		}
	}
}
path "clean/bed" {
	policy = "write"
	permissions = {
		allowedparameters = {
			"*" = {}
		}
  }
}
path "coca/cola" {
	policy = "write"
	permissions = {
		deniedparameters = {
			"john" = {}
		}
	}
}
path "coca/cola" {
	policy = "write"
	permissions = {
		allowedparameters = {
			"john" = {}
		}
  }
}
`
var tokenCreationPolicy = `
name = "tokenCreation"
path "auth/token/create*" {
	capabilities = ["update", "create", "sudo"]
}
`

var aclPolicy = `
name = "dev"
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
name = "ops"
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

//allow operation testing
var permissionsPolicy = `
name = "dev"
path "dev/*" {
	policy = "write"
	
  permissions = {
  	allowedparameters = {
  		"zip" = {}
  	}
  }
}
path "foo/bar" {
	policy = "write"
	permissions = {
		deniedparameters = {
			"zap" = {}
		}
  }
}
path "foo/baz" {
	policy = "write"
	permissions = {
		allowedparameters = {
			"hello" = {}
		}
		deniedparameters = {
			"zap" = {}
		}
  }
}
path "broken/phone" {
	policy = "write"
	permissions = {
		allowedparameters = {
		  "steve" = {}
		}
		deniedparameters = {
		  "steve" = {}
		}
	}
}
path "hello/world" {
	policy = "write"
	permissions = {
		allowedparameters = {
			"*" = {}
		}
		deniedparameters = {
			"*" = {}
		}
  }
}
path "tree/fort" {
	policy = "write"
	permissions = {
		allowedparameters = {
			"*" = {}
		}
		deniedparameters = {
			"beer" = {}
		}
  }
}
path "fruit/apple" {
	policy = "write"
	permissions = {
		allowedparameters = {
			"pear" = {}
		}
		deniedparameters = {
			"*" = {}
		}
  }
}
path "cold/weather" {
	policy = "write"
	permissions = {
		allowedparameters = {}
		deniedparameters = {}
	}
}
`
