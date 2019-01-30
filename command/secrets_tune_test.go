package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testSecretsTuneCommand(tb testing.TB) (*cli.MockUi, *SecretsTuneCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &SecretsTuneCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestSecretsTuneCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			[]string{},
			"Not enough arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"foo", "bar"},
			"Too many arguments",
			1,
		},
	}

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client, closer := testVaultServer(t)
				defer closer()

				ui, cmd := testSecretsTuneCommand(t)
				cmd.client = client

				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("expected %d to be %d", code, tc.code)
				}

				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, tc.out) {
					t.Errorf("expected %q to contain %q", combined, tc.out)
				}
			})
		}
	})

	t.Run("protect_downgrade", func(t *testing.T) {
		t.Parallel()
		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testSecretsTuneCommand(t)
		cmd.client = client

		// Mount
		if err := client.Sys().Mount("kv", &api.MountInput{
			Type: "kv",
			Options: map[string]string{
				"version": "2",
			},
		}); err != nil {
			t.Fatal(err)
		}

		// confirm default max_versions
		mounts, err := client.Sys().ListMounts()
		if err != nil {
			t.Fatal(err)
		}

		mountInfo, ok := mounts["kv/"]
		if !ok {
			t.Fatalf("expected mount to exist")
		}
		if exp := "kv"; mountInfo.Type != exp {
			t.Errorf("expected %q to be %q", mountInfo.Type, exp)
		}
		if exp := "2"; mountInfo.Options["version"] != exp {
			t.Errorf("expected %q to be %q", mountInfo.Options["version"], exp)
		}

		if exp := ""; mountInfo.Options["max_versions"] != exp {
			t.Errorf("expected %s to be empty", mountInfo.Options["max_versions"])
		}

		// omitting the version should not cause a downgrade
		code := cmd.Run([]string{
			"-options", "max_versions=2",
			"kv/",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Tuned the secrets engine at: kv/"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		mounts, err = client.Sys().ListMounts()
		if err != nil {
			t.Fatal(err)
		}

		mountInfo, ok = mounts["kv/"]
		if !ok {
			t.Fatalf("expected mount to exist")
		}
		if exp := "2"; mountInfo.Options["version"] != exp {
			t.Errorf("expected %q to be %q", mountInfo.Options["version"], exp)
		}
		if exp := "kv"; mountInfo.Type != exp {
			t.Errorf("expected %q to be %q", mountInfo.Type, exp)
		}
		if exp := "2"; mountInfo.Options["max_versions"] != exp {
			t.Errorf("expected %s to be %s", mountInfo.Options["max_versions"], exp)
		}
	})

	t.Run("integration", func(t *testing.T) {
		t.Run("flags_all", func(t *testing.T) {
			t.Parallel()
			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testSecretsTuneCommand(t)
			cmd.client = client

			// Mount
			if err := client.Sys().Mount("mount_tune_integration", &api.MountInput{
				Type: "pki",
			}); err != nil {
				t.Fatal(err)
			}

			code := cmd.Run([]string{
				"-description", "new description",
				"-default-lease-ttl", "30m",
				"-max-lease-ttl", "1h",
				"-audit-non-hmac-request-keys", "foo,bar",
				"-audit-non-hmac-response-keys", "foo,bar",
				"-listing-visibility", "unauth",
				"mount_tune_integration/",
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d", code, exp)
			}

			expected := "Success! Tuned the secrets engine at: mount_tune_integration/"
			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, expected) {
				t.Errorf("expected %q to contain %q", combined, expected)
			}

			mounts, err := client.Sys().ListMounts()
			if err != nil {
				t.Fatal(err)
			}

			mountInfo, ok := mounts["mount_tune_integration/"]
			if !ok {
				t.Fatalf("expected mount to exist")
			}
			if exp := "new description"; mountInfo.Description != exp {
				t.Errorf("expected %q to be %q", mountInfo.Description, exp)
			}
			if exp := "pki"; mountInfo.Type != exp {
				t.Errorf("expected %q to be %q", mountInfo.Type, exp)
			}
			if exp := 1800; mountInfo.Config.DefaultLeaseTTL != exp {
				t.Errorf("expected %d to be %d", mountInfo.Config.DefaultLeaseTTL, exp)
			}
			if exp := 3600; mountInfo.Config.MaxLeaseTTL != exp {
				t.Errorf("expected %d to be %d", mountInfo.Config.MaxLeaseTTL, exp)
			}
		})

		t.Run("flags_description", func(t *testing.T) {
			t.Parallel()
			t.Run("not_provided", func(t *testing.T) {
				client, closer := testVaultServer(t)
				defer closer()

				ui, cmd := testSecretsTuneCommand(t)
				cmd.client = client

				// Mount
				if err := client.Sys().Mount("mount_tune_integration", &api.MountInput{
					Type:        "pki",
					Description: "initial description",
				}); err != nil {
					t.Fatal(err)
				}

				code := cmd.Run([]string{
					"-default-lease-ttl", "30m",
					"mount_tune_integration/",
				})
				if exp := 0; code != exp {
					t.Errorf("expected %d to be %d", code, exp)
				}

				expected := "Success! Tuned the secrets engine at: mount_tune_integration/"
				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, expected) {
					t.Errorf("expected %q to contain %q", combined, expected)
				}

				mounts, err := client.Sys().ListMounts()
				if err != nil {
					t.Fatal(err)
				}

				mountInfo, ok := mounts["mount_tune_integration/"]
				if !ok {
					t.Fatalf("expected mount to exist")
				}
				if exp := "initial description"; mountInfo.Description != exp {
					t.Errorf("expected %q to be %q", mountInfo.Description, exp)
				}
			})

			t.Run("provided_empty", func(t *testing.T) {
				client, closer := testVaultServer(t)
				defer closer()

				ui, cmd := testSecretsTuneCommand(t)
				cmd.client = client

				// Mount
				if err := client.Sys().Mount("mount_tune_integration", &api.MountInput{
					Type:        "pki",
					Description: "initial description",
				}); err != nil {
					t.Fatal(err)
				}

				code := cmd.Run([]string{
					"-description", "",
					"mount_tune_integration/",
				})
				if exp := 0; code != exp {
					t.Errorf("expected %d to be %d", code, exp)
				}

				expected := "Success! Tuned the secrets engine at: mount_tune_integration/"
				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, expected) {
					t.Errorf("expected %q to contain %q", combined, expected)
				}

				mounts, err := client.Sys().ListMounts()
				if err != nil {
					t.Fatal(err)
				}

				mountInfo, ok := mounts["mount_tune_integration/"]
				if !ok {
					t.Fatalf("expected mount to exist")
				}
				if exp := ""; mountInfo.Description != exp {
					t.Errorf("expected %q to be %q", mountInfo.Description, exp)
				}
			})
		})
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testSecretsTuneCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"pki/",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error tuning secrets engine pki/: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testSecretsTuneCommand(t)
		assertNoTabs(t, cmd)
	})
}
