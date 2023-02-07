package command

import (
	"strings"
	"sync/atomic"
	"testing"
	"time"

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
		code int64
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
			"haha is an unknown log level",
			1,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, closer := testVaultServer(t)
			defer closer()

			var code int64
			shutdownCh := make(chan struct{})

			ui, cmd := testMonitorCommand(t)
			cmd.client = client
			cmd.ShutdownCh = shutdownCh

			go func() {
				atomic.StoreInt64(&code, int64(cmd.Run(tc.args)))
			}()

			<-time.After(3 * time.Second)
			close(shutdownCh)

			if atomic.LoadInt64(&code) != tc.code {
				t.Errorf("expected %d to be %d", code, tc.code)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, tc.out) {
				t.Fatalf("expected %q to contain %q", combined, tc.out)
			}
		})
	}
}
