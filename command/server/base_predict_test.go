// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/posener/complete"
	"github.com/stretchr/testify/require"
)

// TestPredictVaultPaths verifies that path prediction distinguishes files from
// folders across partial, nested, and mount-level inputs.
func TestPredictVaultPaths(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().Mount("secret", &api.MountInput{Type: "kv"})
	require.NoError(t, err)

	data := map[string]interface{}{"a": "b"}
	if _, err := client.Logical().Write("secret/bar", data); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Logical().Write("secret/foo", data); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Logical().Write("secret/zip/zap", data); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Logical().Write("secret/zip/zonk", data); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Logical().Write("secret/zip/twoot", data); err != nil {
		t.Fatal(err)
	}
	if err := client.Sys().Mount("level1a/level2a/level3a", &api.MountInput{Type: "kv"}); err != nil {
		t.Fatal(err)
	}
	if err := client.Sys().Mount("level1a/level2a/level3b", &api.MountInput{Type: "kv"}); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name         string
		args         complete.Args
		includeFiles bool
		exp          []string
	}{
		{
			"has_args",
			complete.Args{
				All:  []string{"read", "secret/foo", "a=b"},
				Last: "a=b",
			},
			true,
			nil,
		},
		{
			"has_args_no_files",
			complete.Args{
				All:  []string{"read", "secret/foo", "a=b"},
				Last: "a=b",
			},
			false,
			nil,
		},
		{
			"part_mount",
			complete.Args{
				All:  []string{"read", "s"},
				Last: "s",
			},
			true,
			[]string{"secret/", "sys/"},
		},
		{
			"part_mount_no_files",
			complete.Args{
				All:  []string{"read", "s"},
				Last: "s",
			},
			false,
			[]string{"secret/", "sys/"},
		},
		{
			"only_mount",
			complete.Args{
				All:  []string{"read", "sec"},
				Last: "sec",
			},
			true,
			[]string{"secret/bar", "secret/foo", "secret/zip/"},
		},
		{
			"only_mount_no_files",
			complete.Args{
				All:  []string{"read", "sec"},
				Last: "sec",
			},
			false,
			[]string{"secret/zip/"},
		},
		{
			"full_mount",
			complete.Args{
				All:  []string{"read", "secret"},
				Last: "secret",
			},
			true,
			[]string{"secret/bar", "secret/foo", "secret/zip/"},
		},
		{
			"full_mount_no_files",
			complete.Args{
				All:  []string{"read", "secret"},
				Last: "secret",
			},
			false,
			[]string{"secret/zip/"},
		},
		{
			"full_mount_slash",
			complete.Args{
				All:  []string{"read", "secret/"},
				Last: "secret/",
			},
			true,
			[]string{"secret/bar", "secret/foo", "secret/zip/"},
		},
		{
			"full_mount_slash_no_files",
			complete.Args{
				All:  []string{"read", "secret/"},
				Last: "secret/",
			},
			false,
			[]string{"secret/zip/"},
		},
		{
			"path_partial",
			complete.Args{
				All:  []string{"read", "secret/z"},
				Last: "secret/z",
			},
			true,
			[]string{"secret/zip/twoot", "secret/zip/zap", "secret/zip/zonk"},
		},
		{
			"path_partial_no_files",
			complete.Args{
				All:  []string{"read", "secret/z"},
				Last: "secret/z",
			},
			false,
			[]string{"secret/zip/"},
		},
		{
			"subpath_partial_z",
			complete.Args{
				All:  []string{"read", "secret/zip/z"},
				Last: "secret/zip/z",
			},
			true,
			[]string{"secret/zip/zap", "secret/zip/zonk"},
		},
		{
			"subpath_partial_z_no_files",
			complete.Args{
				All:  []string{"read", "secret/zip/z"},
				Last: "secret/zip/z",
			},
			false,
			[]string{"secret/zip/z"},
		},
		{
			"subpath_partial_t",
			complete.Args{
				All:  []string{"read", "secret/zip/t"},
				Last: "secret/zip/t",
			},
			true,
			[]string{"secret/zip/twoot"},
		},
		{
			"subpath_partial_t_no_files",
			complete.Args{
				All:  []string{"read", "secret/zip/t"},
				Last: "secret/zip/t",
			},
			false,
			[]string{"secret/zip/t"},
		},
		{
			"multi_nested",
			complete.Args{
				All:  []string{"read", "level1a/level2a"},
				Last: "level1a/level2a",
			},
			false,
			[]string{
				"level1a/level2a/level3a/",
				"level1a/level2a/level3b/",
			},
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := base.NewPredict()
				p.SetClient(client)

				f := p.VaultFolders()
				if tc.includeFiles {
					f = p.VaultFiles()
				}
				act := f.Predict(tc.args)
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

// TestPredict_Audits verifies that audit predictions list configured devices
// and return no results when Vault is unavailable.
func TestPredict_Audits(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	badClient, badCloser := testVaultServerBad(t)
	defer badCloser()

	if err := client.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": "discard",
		},
	}); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name   string
		client *api.Client
		exp    []string
	}{
		{
			"not_connected_client",
			badClient,
			nil,
		},
		{
			"good_path",
			client,
			[]string{"file/"},
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := base.NewPredict()
				p.SetClient(tc.client)

				act := p.VaultAudits().Predict(complete.Args{})
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

// TestPredict_Mounts verifies that mount predictions use live mounts when
// available and retain the safe built-in fallback when Vault is unavailable.
func TestPredict_Mounts(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	badClient, badCloser := testVaultServerBad(t)
	defer badCloser()

	cases := []struct {
		name   string
		client *api.Client
		exp    []string
	}{
		{
			"not_connected_client",
			badClient,
			[]string{"cubbyhole/"},
		},
		{
			"good_path",
			client,
			[]string{"agent-registry/", "cubbyhole/", "identity/", "sys/"},
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := base.NewPredict()
				p.SetClient(tc.client)

				act := p.VaultMounts().Predict(complete.Args{})
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

// TestPredict_Plugins verifies that plugin prediction exposes the server's
// available built-in and external plugin names.
func TestPredict_Plugins(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	badClient, badCloser := testVaultServerBad(t)
	defer badCloser()

	cases := []struct {
		name   string
		client *api.Client
		exp    []string
	}{
		{
			"not_connected_client",
			badClient,
			nil,
		},
		{
			"good_path",
			client,
			[]string{
				"ad",
				"alicloud",
				"approle",
				"aws",
				"azure",
				"cassandra-database-plugin",
				"cert",
				"cf",
				"consul",
				"couchbase-database-plugin",
				"elasticsearch-database-plugin",
				"gcp",
				"gcpkms",
				"github",
				"hana-database-plugin",
				"influxdb-database-plugin",
				"jwt",
				"kerberos",
				"keymgmt",
				"kmip",
				"kubernetes",
				"kv",
				"ldap",
				"mongodb-database-plugin",
				"mongodbatlas",
				"mongodbatlas-database-plugin",
				"mssql-database-plugin",
				"mysql-aurora-database-plugin",
				"mysql-database-plugin",
				"mysql-legacy-database-plugin",
				"mysql-rds-database-plugin",
				"nomad",
				"oci",
				"oidc",
				"okta",
				"openldap",
				"pcf", // Deprecated.
				"pki",
				"pki-external-ca",
				"postgresql-database-plugin",
				"rabbitmq",
				"radius",
				"redis-database-plugin",
				"redis-elasticache-database-plugin",
				"redshift-database-plugin",
				"saml",
				"scep",
				"snowflake-database-plugin",
				"spiffe",
				"ssh",
				"terraform",
				"totp",
				"transform",
				"transit",
				"userpass",
			},
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := base.NewPredict()
				p.SetClient(tc.client)

				act := p.VaultPlugins().Predict(complete.Args{})

				for _, pluginName := range []string{"keymgmt", "kmip", "transform", "saml", "scep", "spiffe", "pki-external-ca"} {
					if !strutil.StrListContains(act, pluginName) {
						for i, v := range tc.exp {
							if v == pluginName {
								tc.exp = append(tc.exp[:i], tc.exp[i+1:]...)
								break
							}
						}
					}
				}
				if d := cmp.Diff(act, tc.exp); len(d) > 0 {
					t.Errorf("expected: %q, got: %q, diff: %v", tc.exp, act, d)
				}
			})
		}
	})
}

// TestPredict_Policies verifies that policy prediction lists server policies
// and returns no results when Vault is unavailable.
func TestPredict_Policies(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	badClient, badCloser := testVaultServerBad(t)
	defer badCloser()

	cases := []struct {
		name   string
		client *api.Client
		exp    []string
	}{
		{
			"not_connected_client",
			badClient,
			nil,
		},
		{
			"good_path",
			client,
			[]string{"default", "default-ceiling", "root"},
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := base.NewPredict()
				p.SetClient(tc.client)

				act := p.VaultPolicies().Predict(complete.Args{})
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

// TestPredict_Paths verifies KV v1 path prediction for complete, partial, and
// nonexistent paths with and without leaf secrets.
func TestPredict_Paths(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	err := client.Sys().Mount("secret", &api.MountInput{Type: "kv"})
	require.NoError(t, err)

	data := map[string]interface{}{"a": "b"}
	if _, err := client.Logical().Write("secret/bar", data); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Logical().Write("secret/foo", data); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Logical().Write("secret/zip/zap", data); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name         string
		path         string
		includeFiles bool
		exp          []string
	}{
		{
			"bad_path",
			"nope/not/a/real/path/ever",
			true,
			nil,
		},
		{
			"good_path",
			"secret/",
			true,
			[]string{"secret/bar", "secret/foo", "secret/zip/"},
		},
		{
			"good_path_no_files",
			"secret/",
			false,
			[]string{"secret/zip/"},
		},
		{
			"partial_match",
			"secret/z",
			true,
			[]string{"secret/zip/zap"},
		},
		{
			"partial_match_no_files",
			"secret/z",
			false,
			[]string{"secret/zip/"},
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := base.NewPredict()
				p.SetClient(client)
				predictor := p.VaultFolders()
				if tc.includeFiles {
					predictor = p.VaultFiles()
				}

				act := predictor.Predict(
					complete.Args{All: []string{"read", tc.path}, Last: tc.path},
				)
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

// TestPredict_PathsKVv2 verifies that KV v2 metadata paths are translated into
// user-facing predictions with the same file and folder behavior as KV v1.
func TestPredict_PathsKVv2(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().
		Mount("secret", &api.MountInput{Type: "kv", Options: map[string]string{"version": "2"}})
	require.NoError(t, err)

	data := map[string]interface{}{"data": map[string]interface{}{"a": "b"}}
	_, err = client.Logical().Write("secret/data/bar", data)
	require.NoError(t, err)
	_, err = client.Logical().Write("secret/data/foo", data)
	require.NoError(t, err)
	_, err = client.Logical().Write("secret/data/zip/zap", data)
	require.NoError(t, err)

	cases := []struct {
		name         string
		path         string
		includeFiles bool
		exp          []string
	}{
		{
			"bad_path",
			"nope/not/a/real/path/ever",
			true,
			nil,
		},
		{
			"good_path",
			"secret/",
			true,
			[]string{"secret/bar", "secret/foo", "secret/zip/"},
		},
		{
			"good_path_no_files",
			"secret/",
			false,
			[]string{"secret/zip/"},
		},
		{
			"partial_match",
			"secret/z",
			true,
			[]string{"secret/zip/zap"},
		},
		{
			"partial_match_no_files",
			"secret/z",
			false,
			[]string{"secret/zip/"},
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := base.NewPredict()
				p.SetClient(client)
				predictor := p.VaultFolders()
				if tc.includeFiles {
					predictor = p.VaultFiles()
				}

				act := predictor.Predict(
					complete.Args{All: []string{"read", tc.path}, Last: tc.path},
				)
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

// TestPredict_VaultFiles verifies that file prediction returns fully qualified
// secret paths and handles missing paths or unavailable Vault servers.
func TestPredict_VaultFiles(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	badClient, badCloser := testVaultServerBad(t)
	defer badCloser()

	err := client.Sys().Mount("secret", &api.MountInput{Type: "kv"})
	require.NoError(t, err)

	data := map[string]interface{}{"a": "b"}
	if _, err := client.Logical().Write("secret/bar", data); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Logical().Write("secret/foo", data); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name   string
		client *api.Client
		path   string
		exp    []string
	}{
		{
			"bad_path",
			client,
			"nope/not/a/real/path/ever",
			nil,
		},
		{
			"good_path",
			client,
			"secret/",
			[]string{"secret/bar", "secret/foo"},
		},
		{
			"not_connected_client",
			badClient,
			"secret/",
			nil,
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := base.NewPredict()
				p.SetClient(tc.client)

				act := p.VaultFiles().
					Predict(complete.Args{All: []string{"read", tc.path}, Last: tc.path})
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}
