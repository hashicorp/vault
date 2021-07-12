package token

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

// Test is a public function that can be used in other tests to
// test that a helper is functioning properly.
func Test(t *testing.T, h TokenHelper) {
	if err := h.Store("foo"); err != nil {
		t.Fatalf("err: %s", err)
	}

	v, err := h.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if v != "foo" {
		t.Fatalf("bad: %#v", v)
	}

	if err := h.Erase(); err != nil {
		t.Fatalf("err: %s", err)
	}

	v, err = h.Get()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if v != "" {
		t.Fatalf("bad: %#v", v)
	}
}

// TestProcess is used to re-execute this test in order to use it as the
// helper process. For this to work, the TestExternalTokenHelperProcess function must
// exist.
func TestProcess(t *testing.T, s ...string) {
	h := &ExternalTokenHelper{BinaryPath: TestProcessPath(t, s...)}
	Test(t, h)
}

// TestProcessPath returns the path to the test process.
func TestProcessPath(t *testing.T, s ...string) string {
	cs := []string{"-test.run=TestExternalTokenHelperProcess", "--", "GO_WANT_HELPER_PROCESS"}
	cs = append(cs, s...)
	return fmt.Sprintf(
		"%s %s",
		os.Args[0],
		strings.Join(cs, " "))
}

// TestExternalTokenHelperProcessCLI can be called to implement TestExternalTokenHelperProcess
// for TestProcess that just executes a CLI command.
func TestExternalTokenHelperProcessCLI(t *testing.T, cmd cli.Command) {
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}

		args = args[1:]
	}
	if len(args) == 0 || args[0] != "GO_WANT_HELPER_PROCESS" {
		return
	}
	args = args[1:]

	os.Exit(cmd.Run(args))
}
