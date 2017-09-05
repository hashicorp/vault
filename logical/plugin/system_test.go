package plugin

import (
	"testing"

	"reflect"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/logical"
)

func Test_impl(t *testing.T) {
	var _ logical.SystemView = new(SystemViewClient)
}

func TestSystem_defaultLeaseTTL(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	sys := logical.TestSystemView()

	server.RegisterName("Plugin", &SystemViewServer{
		impl: sys,
	})

	testSystemView := &SystemViewClient{client: client}

	expected := sys.DefaultLeaseTTL()
	actual := testSystemView.DefaultLeaseTTL()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_maxLeaseTTL(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	sys := logical.TestSystemView()

	server.RegisterName("Plugin", &SystemViewServer{
		impl: sys,
	})

	testSystemView := &SystemViewClient{client: client}

	expected := sys.MaxLeaseTTL()
	actual := testSystemView.MaxLeaseTTL()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_sudoPrivilege(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	sys := logical.TestSystemView()
	sys.SudoPrivilegeVal = true

	server.RegisterName("Plugin", &SystemViewServer{
		impl: sys,
	})

	testSystemView := &SystemViewClient{client: client}

	expected := sys.SudoPrivilege("foo", "bar")
	actual := testSystemView.SudoPrivilege("foo", "bar")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_tainted(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	sys := logical.TestSystemView()
	sys.TaintedVal = true

	server.RegisterName("Plugin", &SystemViewServer{
		impl: sys,
	})

	testSystemView := &SystemViewClient{client: client}

	expected := sys.Tainted()
	actual := testSystemView.Tainted()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_cachingDisabled(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	sys := logical.TestSystemView()
	sys.CachingDisabledVal = true

	server.RegisterName("Plugin", &SystemViewServer{
		impl: sys,
	})

	testSystemView := &SystemViewClient{client: client}

	expected := sys.CachingDisabled()
	actual := testSystemView.CachingDisabled()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_replicationState(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	sys := logical.TestSystemView()
	sys.ReplicationStateVal = consts.ReplicationPerformancePrimary

	server.RegisterName("Plugin", &SystemViewServer{
		impl: sys,
	})

	testSystemView := &SystemViewClient{client: client}

	expected := sys.ReplicationState()
	actual := testSystemView.ReplicationState()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_responseWrapData(t *testing.T) {
	t.SkipNow()
}

func TestSystem_lookupPlugin(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	sys := logical.TestSystemView()

	server.RegisterName("Plugin", &SystemViewServer{
		impl: sys,
	})

	testSystemView := &SystemViewClient{client: client}

	if _, err := testSystemView.LookupPlugin("foo"); err == nil {
		t.Fatal("LookPlugin(): expected error on due to unsupported call from plugin")
	}
}

func TestSystem_mlockEnabled(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	sys := logical.TestSystemView()
	sys.EnableMlock = true

	server.RegisterName("Plugin", &SystemViewServer{
		impl: sys,
	})

	testSystemView := &SystemViewClient{client: client}

	expected := sys.MlockEnabled()
	actual := testSystemView.MlockEnabled()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}
