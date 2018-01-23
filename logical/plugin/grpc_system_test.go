package plugin

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	"reflect"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/plugin/pb"
)

func TestSystem_GRPC_GRPC_impl(t *testing.T) {
	var _ logical.SystemView = new(gRPCSystemViewClient)
}

func TestSystem_GRPC_defaultLeaseTTL(t *testing.T) {
	sys := logical.TestSystemView()
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()
	testSystemView := newGRPCSystemView(client)

	expected := sys.DefaultLeaseTTL()
	actual := testSystemView.DefaultLeaseTTL()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_GRPC_maxLeaseTTL(t *testing.T) {
	sys := logical.TestSystemView()
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()
	testSystemView := newGRPCSystemView(client)

	expected := sys.MaxLeaseTTL()
	actual := testSystemView.MaxLeaseTTL()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_GRPC_sudoPrivilege(t *testing.T) {
	sys := logical.TestSystemView()
	sys.SudoPrivilegeVal = true
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()
	testSystemView := newGRPCSystemView(client)

	ctx := context.Background()

	expected := sys.SudoPrivilege(ctx, "foo", "bar")
	actual := testSystemView.SudoPrivilege(ctx, "foo", "bar")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_GRPC_tainted(t *testing.T) {
	sys := logical.TestSystemView()
	sys.TaintedVal = true
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()
	testSystemView := newGRPCSystemView(client)

	expected := sys.Tainted()
	actual := testSystemView.Tainted()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_GRPC_cachingDisabled(t *testing.T) {
	sys := logical.TestSystemView()
	sys.CachingDisabledVal = true
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()
	testSystemView := newGRPCSystemView(client)

	expected := sys.CachingDisabled()
	actual := testSystemView.CachingDisabled()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_GRPC_replicationState(t *testing.T) {
	sys := logical.TestSystemView()
	sys.ReplicationStateVal = consts.ReplicationPerformancePrimary
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()
	testSystemView := newGRPCSystemView(client)

	expected := sys.ReplicationState()
	actual := testSystemView.ReplicationState()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_GRPC_responseWrapData(t *testing.T) {
	t.SkipNow()
}

func TestSystem_GRPC_lookupPlugin(t *testing.T) {
	sys := logical.TestSystemView()
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()

	testSystemView := newGRPCSystemView(client)

	if _, err := testSystemView.LookupPlugin(context.Background(), "foo"); err == nil {
		t.Fatal("LookPlugin(): expected error on due to unsupported call from plugin")
	}
}

func TestSystem_GRPC_mlockEnabled(t *testing.T) {
	sys := logical.TestSystemView()
	sys.EnableMlock = true
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()

	testSystemView := newGRPCSystemView(client)

	expected := sys.MlockEnabled()
	actual := testSystemView.MlockEnabled()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}
