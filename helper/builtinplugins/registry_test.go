package builtinplugins

import (
	"reflect"
	"testing"

	credAppId "github.com/hashicorp/vault/builtin/credential/app-id"
	dbMysql "github.com/hashicorp/vault/plugins/database/mysql"
	"github.com/hashicorp/vault/sdk/helper/consts"
)

// Test_RegistryGet exercises the (registry).Get functionality by comparing
// factory types and ok response.
func Test_RegistryGet(t *testing.T) {
	tests := []struct {
		name       string
		builtin    string
		pluginType consts.PluginType
		want       BuiltinFactory
		wantOk     bool
	}{
		{
			name:       "non-existent builtin",
			builtin:    "foo",
			pluginType: consts.PluginTypeCredential,
			want:       nil,
			wantOk:     false,
		},
		{
			name:       "bad plugin type",
			builtin:    "app-id",
			pluginType: 9000,
			want:       nil,
			wantOk:     false,
		},
		{
			name:       "known builtin lookup",
			builtin:    "app-id",
			pluginType: consts.PluginTypeCredential,
			want:       toFunc(credAppId.Factory),
			wantOk:     true,
		},
		{
			name:       "known builtin lookup",
			builtin:    "mysql-database-plugin",
			pluginType: consts.PluginTypeDatabase,
			want:       dbMysql.New(dbMysql.DefaultUserNameTemplate),
			wantOk:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got BuiltinFactory
			got, ok := Registry.Get(tt.builtin, tt.pluginType)
			if ok {
				if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
					t.Fatalf("got type: %T, want type: %T", got, tt.want)
				}
			}
			if tt.wantOk != ok {
				t.Fatalf("error: got %v, want %v", ok, tt.wantOk)
			}
		})
	}
}

// Test_RegistryKeyCounts is a light unit test used to check the builtin
// registry lists for each plugin type and make sure they match in length.
func Test_RegistryKeyCounts(t *testing.T) {
	tests := []struct {
		name       string
		pluginType consts.PluginType
		want       int // use slice length as test condition
		wantOk     bool
	}{
		{
			name:       "bad plugin type",
			pluginType: 9001,
			want:       0,
		},
		{
			name:       "number of auth plugins",
			pluginType: consts.PluginTypeCredential,
			want:       20,
		},
		{
			name:       "number of database plugins",
			pluginType: consts.PluginTypeDatabase,
			want:       17,
		},
		{
			name:       "number of secrets plugins",
			pluginType: consts.PluginTypeSecrets,
			want:       18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := Registry.Keys(tt.pluginType)
			if len(keys) != tt.want {
				t.Fatalf("got size: %d, want size: %d", len(keys), tt.want)
			}
		})
	}
}

// Test_RegistryContains exercises the (registry).Contains functionality.
func Test_RegistryContains(t *testing.T) {
	tests := []struct {
		name       string
		builtin    string
		pluginType consts.PluginType
		want       bool
	}{
		{
			name:       "non-existent builtin",
			builtin:    "foo",
			pluginType: consts.PluginTypeCredential,
			want:       false,
		},
		{
			name:       "bad plugin type",
			builtin:    "app-id",
			pluginType: 9001,
			want:       false,
		},
		{
			name:       "known builtin lookup",
			builtin:    "app-id",
			pluginType: consts.PluginTypeCredential,
			want:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Registry.Contains(tt.builtin, tt.pluginType)
			if got != tt.want {
				t.Fatalf("error: got %v, wanted %v", got, tt.want)
			}
		})
	}
}

// Test_RegistryStatus exercises the (registry).Status functionality.
func Test_RegistryStatus(t *testing.T) {
	tests := []struct {
		name       string
		builtin    string
		pluginType consts.PluginType
		want       consts.DeprecationStatus
		wantOk     bool
	}{
		{
			name:       "non-existent builtin and valid type",
			builtin:    "foo",
			pluginType: consts.PluginTypeCredential,
			want:       consts.Unknown,
			wantOk:     false,
		},
		{
			name:       "mismatch builtin and plugin type",
			builtin:    "app-id",
			pluginType: consts.PluginTypeSecrets,
			want:       consts.Unknown,
			wantOk:     false,
		},
		{
			name:       "existing builtin and invalid plugin type",
			builtin:    "app-id",
			pluginType: 9000,
			want:       consts.Unknown,
			wantOk:     false,
		},
		{
			name:       "supported builtin lookup",
			builtin:    "approle",
			pluginType: consts.PluginTypeCredential,
			want:       consts.Supported,
			wantOk:     true,
		},
		{
			name:       "deprecated builtin lookup",
			builtin:    "pcf",
			pluginType: consts.PluginTypeCredential,
			want:       consts.Deprecated,
			wantOk:     true,
		},
		{
			name:       "pending removal builtin lookup",
			builtin:    "app-id",
			pluginType: consts.PluginTypeCredential,
			want:       consts.PendingRemoval,
			wantOk:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Registry.DeprecationStatus(tt.builtin, tt.pluginType)
			if got != tt.want {
				t.Fatalf("got %+v, wanted %+v", got, tt.want)
			}
			if ok != tt.wantOk {
				t.Fatalf("got ok: %t, want ok: %t", ok, tt.wantOk)
			}
		})
	}
}
