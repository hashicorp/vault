package command

import (
	"archive/tar"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mholt/archiver"
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

	testDir, err := ioutil.TempDir("", "vault-debug")
	if err != nil {
		t.Fatal(err)
	}
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

func TestDebugCommand_Archive(t *testing.T) {
	t.Parallel()

	testDir, err := ioutil.TempDir("", "vault-debug")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	client, closer := testVaultServer(t)
	defer closer()

	_, cmd := testDebugCommand(t)
	cmd.client = client
	cmd.skipTimingChecks = true

	args := []string{
		"-duration=1s",
		fmt.Sprintf("-output=%s/archive", testDir),
		"-targets=server-status",
	}

	code := cmd.Run(args)
	if exp := 0; code != exp {
		t.Errorf("expected %d to be %d", code, exp)
	}

	basePath := "archive"
	bundlePath := filepath.Join(testDir, basePath+debugCompressionExt)
	_, err = os.Open(bundlePath)
	if err != nil {
		t.Fatalf("failed to open archive: %s", err)
	}

	tgz := archiver.NewTarGz()
	err = tgz.Walk(bundlePath, func(f archiver.File) error {
		fh, ok := f.Header.(*tar.Header)
		if !ok {
			t.Fatalf("invalid file header: %#v", f.Header)
		}

		// Ignore base directory
		if fh.Name == basePath+"/" {
			return nil
		}

		if fh.Name != filepath.Join(basePath, "index.json") && fh.Name != filepath.Join(basePath, "server_status.json") {
			t.Fatalf("unxexpected file: %s", fh.Name)
		}
		return nil
	})
}
