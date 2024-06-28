// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package builtinplugins

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"

	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/constants"
	dbMysql "github.com/hashicorp/vault/plugins/database/mysql"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"golang.org/x/exp/slices"
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
			builtin:    "userpass",
			pluginType: consts.PluginTypeCredential,
			want:       toFunc(credUserpass.Factory),
			wantOk:     true,
		},
		{
			name:       "removed builtin lookup",
			builtin:    "app-id",
			pluginType: consts.PluginTypeCredential,
			want:       nil,
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
		entWant    int
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
			want:       18,
			entWant:    1,
		},
		{
			name:       "number of database plugins",
			pluginType: consts.PluginTypeDatabase,
			want:       17,
		},
		{
			name:       "number of secrets plugins",
			pluginType: consts.PluginTypeSecrets,
			want:       19,
			entWant:    3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := Registry.Keys(tt.pluginType)
			want := tt.want
			if constants.IsEnterprise {
				want += tt.entWant
			}
			if len(keys) != want {
				t.Fatalf("got size: %d, want size: %d", len(keys), want)
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
			builtin:    "approle",
			pluginType: consts.PluginTypeCredential,
			want:       true,
		},
		{
			name:       "removed builtin lookup",
			builtin:    "app-id",
			pluginType: consts.PluginTypeCredential,
			want:       false,
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
			name:       "removed builtin lookup",
			builtin:    "app-id",
			pluginType: consts.PluginTypeCredential,
			want:       consts.Removed,
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

// Test_RegistryMatchesGenOpenapi ensures that the plugins mounted in gen_openapi.sh match registry.go
func Test_RegistryMatchesGenOpenapi(t *testing.T) {
	const scriptPath = "../../scripts/gen_openapi.sh"

	// parseScript fetches the contents of gen_openapi.sh script & extract the relevant lines
	parseScript := func(path string) ([]string, []string, error) {
		f, err := os.Open(scriptPath)
		if err != nil {
			return nil, nil, fmt.Errorf("could not open gen_openapi.sh script: %w", err)
		}
		defer f.Close()

		// This is a hack: the gen_openapi script contains a conditional block to
		// enable the enterprise plugins, whose lines are indented.  Tweak the
		// regexp to only include the indented lines on enterprise.
		leading := "^"
		if constants.IsEnterprise {
			leading = "^ *"
		}

		var (
			credentialBackends   []string
			credentialBackendsRe = regexp.MustCompile(leading + `vault auth enable (?:-.+ )*(?:"([a-zA-Z]+)"|([a-zA-Z]+))$`)

			secretsBackends   []string
			secretsBackendsRe = regexp.MustCompile(leading + `vault secrets enable (?:-.+ )*(?:"([a-zA-Z]+)"|([a-zA-Z]+))$`)
		)

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := scanner.Text()

			if m := credentialBackendsRe.FindStringSubmatch(line); m != nil {
				credentialBackends = append(credentialBackends, m[1])
			}
			if m := secretsBackendsRe.FindStringSubmatch(line); m != nil {
				secretsBackends = append(secretsBackends, m[1])
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, nil, fmt.Errorf("error scanning gen_openapi.sh: %v", err)
		}

		return credentialBackends, secretsBackends, nil
	}

	// ensureInRegistry ensures that the given plugin is in registry and marked as "supported"
	ensureInRegistry := func(t *testing.T, name string, pluginType consts.PluginType) {
		t.Helper()

		// "database" will not be present in registry, it is represented as
		// a list of database plugins instead
		if name == "database" && pluginType == consts.PluginTypeSecrets {
			return
		}

		deprecationStatus, ok := Registry.DeprecationStatus(name, pluginType)
		if !ok {
			t.Errorf("%q %s backend is missing from registry.go; please remove it from gen_openapi.sh", name, pluginType)
		}

		if deprecationStatus == consts.Removed {
			t.Errorf("%q %s backend is marked 'removed' in registry.go; please remove it from gen_openapi.sh", name, pluginType)
		}
	}

	// ensureInScript ensures that the given plugin name is in gen_openapi.sh script
	ensureInScript := func(t *testing.T, scriptBackends []string, name string) {
		t.Helper()

		for _, excluded := range []string{
			"oidc",     // alias for "jwt"
			"openldap", // alias for "ldap"
		} {
			if name == excluded {
				return
			}
		}

		if !slices.Contains(scriptBackends, name) {
			t.Errorf("%q backend could not be found in gen_openapi.sh, please add it there", name)
		}
	}

	// test starts here
	scriptCredentialBackends, scriptSecretsBackends, err := parseScript(scriptPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, name := range scriptCredentialBackends {
		ensureInRegistry(t, name, consts.PluginTypeCredential)
	}

	for _, name := range scriptSecretsBackends {
		ensureInRegistry(t, name, consts.PluginTypeSecrets)
	}

	for name, backend := range Registry.credentialBackends {
		if backend.DeprecationStatus == consts.Supported {
			ensureInScript(t, scriptCredentialBackends, name)
		}
	}

	for name, backend := range Registry.logicalBackends {
		if backend.DeprecationStatus == consts.Supported {
			ensureInScript(t, scriptSecretsBackends, name)
		}
	}
}
