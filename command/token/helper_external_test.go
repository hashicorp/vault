package token

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestExternalTokenHelperPath(t *testing.T) {
	cases := map[string]string{}

	unixCases := map[string]string{
		"/foo": "/foo",
	}
	windowsCases := map[string]string{
		"C:/foo":           "C:/foo",
		`C:\Program Files`: `C:\Program Files`,
	}

	var runtimeCases map[string]string
	if runtime.GOOS == "windows" {
		runtimeCases = windowsCases
	} else {
		runtimeCases = unixCases
	}

	for k, v := range runtimeCases {
		cases[k] = v
	}

	// We don't expect those to actually exist, so we expect an error. For now,
	// I'm commenting out the rest of this code as we don't have real external
	// helpers to test with and the os.Stat will fail with our fake test cases.
	/*
		for k, v := range cases {
			actual, err := ExternalTokenHelperPath(k)
				if err != nil {
					t.Fatalf("error getting external helper path: %v", err)
				}
				if actual != v {
					t.Fatalf(
						"input: %s, expected: %s, got: %s",
						k, v, actual)
				}
		}
	*/
}

func TestExternalTokenHelper(t *testing.T) {
	Test(t, testExternalTokenHelper(t))
}

func testExternalTokenHelper(t *testing.T) *ExternalTokenHelper {
	return &ExternalTokenHelper{BinaryPath: helperPath("helper"), Env: helperEnv()}
}

func helperPath(s ...string) string {
	cs := []string{"-test.run=TestExternalTokenHelperProcess", "--"}
	cs = append(cs, s...)
	return fmt.Sprintf(
		"%s %s",
		os.Args[0],
		strings.Join(cs, " "))
}

func helperEnv() []string {
	var env []string

	tf, err := ioutil.TempFile("", "vault")
	if err != nil {
		panic(err)
	}
	tf.Close()

	env = append(env, "GO_HELPER_PATH="+tf.Name(), "GO_WANT_HELPER_PROCESS=1")
	return env
}

// This is not a real test. This is just a helper process kicked off by tests.
func TestExternalTokenHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}

		args = args[1:]
	}

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "helper":
		path := os.Getenv("GO_HELPER_PATH")

		switch args[0] {
		case "erase":
			os.Remove(path)
		case "get":
			f, err := os.Open(path)
			if os.IsNotExist(err) {
				return
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Err: %s\n", err)
				os.Exit(1)
			}
			defer f.Close()
			io.Copy(os.Stdout, f)
		case "store":
			f, err := os.Create(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Err: %s\n", err)
				os.Exit(1)
			}
			defer f.Close()
			io.Copy(f, os.Stdin)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n", cmd)
		os.Exit(2)
	}
}
