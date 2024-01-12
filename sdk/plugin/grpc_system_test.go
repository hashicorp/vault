// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func TestSystem_GRPC_ReturnsErrIfSystemViewNil(t *testing.T) {
	_, err := new(gRPCSystemViewServer).ReplicationState(context.Background(), nil)
	if err == nil {
		t.Error("Expected error when using server with no impl")
	}
}

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

func TestSystem_GRPC_lookupPlugin(t *testing.T) {
	sys := logical.TestSystemView()
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()

	testSystemView := newGRPCSystemView(client)

	if _, err := testSystemView.LookupPlugin(context.Background(), "foo", consts.PluginTypeDatabase); err == nil {
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

func TestSystem_GRPC_entityInfo(t *testing.T) {
	sys := logical.TestSystemView()
	sys.EntityVal = &logical.Entity{
		ID:   "id",
		Name: "name",
		Metadata: map[string]string{
			"foo": "bar",
		},
		Aliases: []*logical.Alias{
			{
				MountType:     "logical",
				MountAccessor: "accessor",
				Name:          "name",
				Metadata: map[string]string{
					"zip": "zap",
				},
			},
		},
		Disabled: true,
	}
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()
	testSystemView := newGRPCSystemView(client)

	actual, err := testSystemView.EntityInfo("")
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(sys.EntityVal, actual) {
		t.Fatalf("expected: %v, got: %v", sys.EntityVal, actual)
	}
}

func TestSystem_GRPC_groupsForEntity(t *testing.T) {
	sys := logical.TestSystemView()
	sys.GroupsVal = []*logical.Group{
		{
			ID:   "group1-id",
			Name: "group1",
			Metadata: map[string]string{
				"group-metadata": "metadata-value",
			},
		},
	}
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()
	testSystemView := newGRPCSystemView(client)

	actual, err := testSystemView.GroupsForEntity("")
	if err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(sys.GroupsVal[0], actual[0]) {
		t.Fatalf("expected: %v, got: %v", sys.GroupsVal, actual)
	}
}

func TestSystem_GRPC_pluginEnv(t *testing.T) {
	sys := logical.TestSystemView()
	sys.PluginEnvironment = &logical.PluginEnvironment{
		VaultVersion:           "0.10.42",
		VaultVersionPrerelease: "dev",
		VaultVersionMetadata:   "ent",
	}
	client, _ := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer client.Close()

	testSystemView := newGRPCSystemView(client)

	expected, err := sys.PluginEnv(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	actual, err := testSystemView.PluginEnv(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func TestSystem_GRPC_GeneratePasswordFromPolicy(t *testing.T) {
	policyName := "testpolicy"
	expectedPassword := "87354qtnjgrehiogd9u0t43"
	passGen := func() (password string, err error) {
		return expectedPassword, nil
	}
	sys := &logical.StaticSystemView{
		PasswordPolicies: map[string]logical.PasswordGenerator{
			policyName: passGen,
		},
	}

	client, server := plugin.TestGRPCConn(t, func(s *grpc.Server) {
		pb.RegisterSystemViewServer(s, &gRPCSystemViewServer{
			impl: sys,
		})
	})
	defer server.Stop()
	defer client.Close()

	testSystemView := newGRPCSystemView(client)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	password, err := testSystemView.GeneratePasswordFromPolicy(ctx, policyName)
	if err != nil {
		t.Fatalf("no error expected, got: %s", err)
	}

	if password != expectedPassword {
		t.Fatalf("Actual password: %s\nExpected password: %s", password, expectedPassword)
	}
}
