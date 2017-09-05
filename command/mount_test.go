package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testMountCommand(tb testing.TB) (*cli.MockUi, *MountCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &MountCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}
func TestMountCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"empty",
			nil,
			"Missing TYPE!",
			1,
		},
		{
			"too_many_args",
			[]string{"foo", "bar"},
			"Too many arguments",
			1,
		},
		{
			"not_a_valid_mount",
			[]string{"nope_definitely_not_a_valid_mount_like_ever"},
			"",
			2,
		},
		{
			"mount",
			[]string{"transit"},
			"Success! Mounted the transit secret backend at: transit/",
			0,
		},
		{
			"mount_path",
			[]string{
				"-path", "transit_mount_point",
				"transit",
			},
			"Success! Mounted the transit secret backend at: transit_mount_point/",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testMountCommand(t)
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

	t.Run("integration", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testMountCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-path", "mount_integration/",
			"-description", "The best kind of test",
			"-default-lease-ttl", "30m",
			"-max-lease-ttl", "1h",
			"-force-no-cache",
			"pki",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Mounted the pki secret backend at: mount_integration/"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		mounts, err := client.Sys().ListMounts()
		if err != nil {
			t.Fatal(err)
		}

		mountInfo, ok := mounts["mount_integration/"]
		if !ok {
			t.Fatalf("expected mount to exist")
		}
		if exp := "pki"; mountInfo.Type != exp {
			t.Errorf("expected %q to be %q", mountInfo.Type, exp)
		}
		if exp := "The best kind of test"; mountInfo.Description != exp {
			t.Errorf("expected %q to be %q", mountInfo.Description, exp)
		}
		if exp := 1800; mountInfo.Config.DefaultLeaseTTL != exp {
			t.Errorf("expected %d to be %d", mountInfo.Config.DefaultLeaseTTL, exp)
		}
		if exp := 3600; mountInfo.Config.MaxLeaseTTL != exp {
			t.Errorf("expected %d to be %d", mountInfo.Config.MaxLeaseTTL, exp)
		}
		if exp := true; mountInfo.Config.ForceNoCache != exp {
			t.Errorf("expected %t to be %t", mountInfo.Config.ForceNoCache, exp)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testMountCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"pki",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error mounting: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testMountCommand(t)
		assertNoTabs(t, cmd)
	})
}
