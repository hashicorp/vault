package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testAuthTuneCommand(tb testing.TB) (*cli.MockUi, *AuthTuneCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AuthTuneCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestAuthTuneCommand_Run(t *testing.T) {
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

				ui, cmd := testAuthTuneCommand(t)
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

	t.Run("integration", func(t *testing.T) {
		t.Run("flags_all", func(t *testing.T) {
			t.Parallel()
			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testAuthTuneCommand(t)
			cmd.client = client

			// Mount
			if err := client.Sys().EnableAuthWithOptions("my-auth", &api.EnableAuthOptions{
				Type: "userpass",
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
				"my-auth/",
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d", code, exp)
			}

			expected := "Success! Tuned the auth method at: my-auth/"
			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, expected) {
				t.Errorf("expected %q to contain %q", combined, expected)
			}

			auths, err := client.Sys().ListAuth()
			if err != nil {
				t.Fatal(err)
			}

			mountInfo, ok := auths["my-auth/"]
			if !ok {
				t.Fatalf("expected auth to exist")
			}
			if exp := "new description"; mountInfo.Description != exp {
				t.Errorf("expected %q to be %q", mountInfo.Description, exp)
			}
			if exp := "userpass"; mountInfo.Type != exp {
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

				ui, cmd := testAuthTuneCommand(t)
				cmd.client = client

				// Mount
				if err := client.Sys().EnableAuthWithOptions("my-auth", &api.EnableAuthOptions{
					Type:        "userpass",
					Description: "initial description",
				}); err != nil {
					t.Fatal(err)
				}

				code := cmd.Run([]string{
					"-default-lease-ttl", "30m",
					"my-auth/",
				})
				if exp := 0; code != exp {
					t.Errorf("expected %d to be %d", code, exp)
				}

				expected := "Success! Tuned the auth method at: my-auth/"
				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, expected) {
					t.Errorf("expected %q to contain %q", combined, expected)
				}

				auths, err := client.Sys().ListAuth()
				if err != nil {
					t.Fatal(err)
				}

				mountInfo, ok := auths["my-auth/"]
				if !ok {
					t.Fatalf("expected auth to exist")
				}
				if exp := "initial description"; mountInfo.Description != exp {
					t.Errorf("expected %q to be %q", mountInfo.Description, exp)
				}
			})

			t.Run("provided_empty", func(t *testing.T) {
				client, closer := testVaultServer(t)
				defer closer()

				ui, cmd := testAuthTuneCommand(t)
				cmd.client = client

				// Mount
				if err := client.Sys().EnableAuthWithOptions("my-auth", &api.EnableAuthOptions{
					Type:        "userpass",
					Description: "initial description",
				}); err != nil {
					t.Fatal(err)
				}

				code := cmd.Run([]string{
					"-description", "",
					"my-auth/",
				})
				if exp := 0; code != exp {
					t.Errorf("expected %d to be %d", code, exp)
				}

				expected := "Success! Tuned the auth method at: my-auth/"
				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, expected) {
					t.Errorf("expected %q to contain %q", combined, expected)
				}

				auths, err := client.Sys().ListAuth()
				if err != nil {
					t.Fatal(err)
				}

				mountInfo, ok := auths["my-auth/"]
				if !ok {
					t.Fatalf("expected auth to exist")
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

		ui, cmd := testAuthTuneCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"userpass/",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error tuning auth method userpass/: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuthTuneCommand(t)
		assertNoTabs(t, cmd)
	})
}
