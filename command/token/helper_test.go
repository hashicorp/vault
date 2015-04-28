package token

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestHelperPath(t *testing.T) {
	cases := map[string]string{
		"/foo": "/foo",
		"foo":  exePath + " token-foo",
	}

	for k, v := range cases {
		actual := HelperPath(k)
		if actual != v {
			t.Fatalf(
				"input: %s, expected: %s, got: %s",
				k, v, actual)
		}
	}
}

func TestHelper(t *testing.T) {
	h := testHelper(t)
	Test(t, h.Path)
}

func testHelper(t *testing.T) *Helper {
	return &Helper{Path: helperPath("helper")}
}

func helperPath(s ...string) string {
	tf, err := ioutil.TempFile("", "vault")
	if err != nil {
		panic(err)
	}
	tf.Close()

	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	return fmt.Sprintf(
		"GO_HELPER_PATH=%s GO_WANT_HELPER_PROCESS=1 %s %s",
		tf.Name(),
		os.Args[0],
		strings.Join(cs, " "))
}

// This is not a real test. This is just a helper process kicked off by tests.
func TestHelperProcess(*testing.T) {
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
