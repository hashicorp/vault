package consul

import (
	"reflect"
	"testing"

	"github.com/hashicorp/consul/api"
)

func TestToken_parseServiceIdentities(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []*api.ACLServiceIdentity
	}{
		{
			name: "No datacenters",
			args: []string{"myservice-1"},
			want: []*api.ACLServiceIdentity{{ServiceName: "myservice-1", Datacenters: nil}},
		},
		{
			name: "One datacenter",
			args: []string{"myservice-1:dc1"},
			want: []*api.ACLServiceIdentity{{ServiceName: "myservice-1", Datacenters: []string{"dc1"}}},
		},
		{
			name: "Multiple datacenters",
			args: []string{"myservice-1:dc1,dc2,dc3"},
			want: []*api.ACLServiceIdentity{{ServiceName: "myservice-1", Datacenters: []string{"dc1", "dc2", "dc3"}}},
		},
		{
			name: "Missing service name with datacenter",
			args: []string{":dc1"},
			want: []*api.ACLServiceIdentity{{ServiceName: "", Datacenters: []string{"dc1"}}},
		},
		{
			name: "Missing service name and missing datacenter",
			args: []string{""},
			want: []*api.ACLServiceIdentity{{ServiceName: "", Datacenters: nil}},
		},
		{
			name: "Multiple service identities",
			args: []string{"myservice-1:dc1", "myservice-2:dc1", "myservice-3:dc1,dc2"},
			want: []*api.ACLServiceIdentity{
				{ServiceName: "myservice-1", Datacenters: []string{"dc1"}},
				{ServiceName: "myservice-2", Datacenters: []string{"dc1"}},
				{ServiceName: "myservice-3", Datacenters: []string{"dc1", "dc2"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseServiceIdentities(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseServiceIdentities() = {%s:%v}, want {%s:%v}", got[0].ServiceName, got[0].Datacenters, tt.want[0].ServiceName, tt.want[0].Datacenters)
			}
		})
	}
}

func TestToken_parseNodeIdentities(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []*api.ACLNodeIdentity
	}{
		{
			name: "No datacenter",
			args: []string{"server-1"},
			want: []*api.ACLNodeIdentity{{NodeName: "server-1", Datacenter: ""}},
		},
		{
			name: "One datacenter",
			args: []string{"server-1:dc1"},
			want: []*api.ACLNodeIdentity{{NodeName: "server-1", Datacenter: "dc1"}},
		},
		{
			name: "Missing node name with datacenter",
			args: []string{":dc1"},
			want: []*api.ACLNodeIdentity{{NodeName: "", Datacenter: "dc1"}},
		},
		{
			name: "Missing node name and missing datacenter",
			args: []string{""},
			want: []*api.ACLNodeIdentity{{NodeName: "", Datacenter: ""}},
		},
		{
			name: "Multiple node identities",
			args: []string{"server-1:dc1", "server-2:dc1", "server-3:dc1"},
			want: []*api.ACLNodeIdentity{
				{NodeName: "server-1", Datacenter: "dc1"},
				{NodeName: "server-2", Datacenter: "dc1"},
				{NodeName: "server-3", Datacenter: "dc1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNodeIdentities(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNodeIdentities() = {%s:%s}, want {%s:%s}", got[0].NodeName, got[0].Datacenter, tt.want[0].NodeName, tt.want[0].Datacenter)
			}
		})
	}
}
