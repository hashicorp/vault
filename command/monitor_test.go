package command

import (
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/mitchellh/cli"
)

func testMonitorCommand(tb testing.TB) (*cli.MockUi, *MonitorCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &MonitorCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestMonitorCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"valid",
			[]string{
				"-log-level=debug",
			},
			"",
			0,
		},
		{
			"too_many_args",
			[]string{
				"-log-level=debug",
				"foo",
			},
			"Too many arguments",
			1,
		},
		{
			"unknown_log_level",
			[]string{
				"-log-level=haha",
			},
			"HAHA is an unknown log level",
			1,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testMonitorCommand(t)
			cmd.client = client

			var code int
			stopCh := make(chan struct{})

			testhelpers.GenerateDebugLogs(t, stopCh, client)

			go func() {
				code = cmd.Run(tc.args)
			}()

			select {
			case <-time.After(3 * time.Second):
				close(stopCh)
				syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
			}

			if code != tc.code {
				t.Errorf("expected %d to be %d", code, tc.code)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, tc.out) {
				t.Fatalf("expected %q to contain %q", combined, tc.out)
			}
		})
	}
}

func TestMonitorCommand_NoTabs(t *testing.T) {
	t.Parallel()

	_, cmd := testMonitorCommand(t)
	assertNoTabs(t, cmd)
}
