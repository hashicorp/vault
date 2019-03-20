package command

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

func TestPredictVaultPaths(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

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
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := NewPredict()
				p.client = client

				f := p.vaultPaths(tc.includeFiles)
				act := f(tc.args)
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

func TestPredict_Audits(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

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

				p := NewPredict()
				p.client = tc.client

				act := p.audits()
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

func TestPredict_Mounts(t *testing.T) {
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
			defaultPredictVaultMounts,
		},
		{
			"good_path",
			client,
			[]string{"cubbyhole/", "identity/", "secret/", "sys/"},
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := NewPredict()
				p.client = tc.client

				act := p.mounts()
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

func TestPredict_Plugins(t *testing.T) {
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
			[]string{
				"ad",
				"alicloud",
				"app-id",
				"approle",
				"aws",
				"azure",
				"cassandra",
				"cassandra-database-plugin",
				"centrify",
				"cert",
				"consul",
				"gcp",
				"gcpkms",
				"github",
				"hana-database-plugin",
				"influxdb-database-plugin",
				"jwt",
				"kubernetes",
				"kv",
				"ldap",
				"mongodb",
				"mongodb-database-plugin",
				"mssql",
				"mssql-database-plugin",
				"mysql",
				"mysql-aurora-database-plugin",
				"mysql-database-plugin",
				"mysql-legacy-database-plugin",
				"mysql-rds-database-plugin",
				"nomad",
				"oidc",
				"okta",
				"pki",
				"postgresql",
				"postgresql-database-plugin",
				"rabbitmq",
				"radius",
				"ssh",
				"totp",
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

				p := NewPredict()
				p.client = tc.client

				act := p.plugins()
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

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
			[]string{"default", "root"},
		},
	}

	t.Run("group", func(t *testing.T) {
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				p := NewPredict()
				p.client = tc.client

				act := p.policies()
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

func TestPredict_Paths(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

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
			[]string{"nope/not/a/real/path/ever"},
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
			[]string{"secret/zip/"},
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

				p := NewPredict()
				p.client = client

				act := p.paths(tc.path, tc.includeFiles)
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

func TestPredict_ListPaths(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	badClient, badCloser := testVaultServerBad(t)
	defer badCloser()

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
			[]string{"bar", "foo"},
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

				p := NewPredict()
				p.client = tc.client

				act := p.listPaths(tc.path)
				if !reflect.DeepEqual(act, tc.exp) {
					t.Errorf("expected %q to be %q", act, tc.exp)
				}
			})
		}
	})
}

func TestPredict_HasPathArg(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		exp  bool
	}{
		{
			"nil",
			nil,
			false,
		},
		{
			"empty",
			[]string{},
			false,
		},
		{
			"empty_string",
			[]string{""},
			false,
		},
		{
			"single",
			[]string{"foo"},
			false,
		},
		{
			"multiple",
			[]string{"foo", "bar", "baz"},
			true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := NewPredict()
			if act := p.hasPathArg(tc.args); act != tc.exp {
				t.Errorf("expected %t to be %t", act, tc.exp)
			}
		})
	}
}
