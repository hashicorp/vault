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
func Test(t *testing.T, path string) {
	h := &Helper{Path: path}
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
// helper process. For this to work, the TestHelperProcess function must
// exist.
func TestProcess(t *testing.T, s ...string) {
	// Build the path to the CLI to execute
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	path := fmt.Sprintf(
		"GO_WANT_HELPER_PROCESS=1 %s %s",
		os.Args[0],
		strings.Join(cs, " "))

	// Run the tests
	Test(t, path)
}

// TestHelperProcessCLI can be called to implement TestHelperProcess
// for TestProcess that just executes a CLI command.
func TestHelperProcessCLI(t *testing.T, cmd cli.Command) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}

		args = args[1:]
	}

	os.Exit(cmd.Run(args))
}
