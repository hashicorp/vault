package command

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/mitchellh/cli"
)

func testDebugCommand(tb testing.TB) (*cli.MockUi, *DebugCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &DebugCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestDebugCommand_Run(t *testing.T) {
	t.Parallel()

	testDir := testutil.TempDir(t, "debug")
	defer os.RemoveAll(testDir)

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"valid",
			[]string{
				"-duration=1s",
				fmt.Sprintf("-output=%s/valid", testDir),
			},
			"",
			0,
		},
		{
			"too_many_args",
			[]string{
				"-duration=1s",
				fmt.Sprintf("-output=%s/too_many_args", testDir),
				"foo",
			},
			"Too many arguments",
			1,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testDebugCommand(t)
			cmd.client = client
			cmd.skipTimingChecks = true

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
}
