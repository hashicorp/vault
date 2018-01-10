package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestParseArgsData(t *testing.T) {
	t.Parallel()

	t.Run("stdin_full", func(t *testing.T) {
		t.Parallel()

		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte(`{"foo":"bar"}`))
			stdinW.Close()
		}()

		m, err := parseArgsData(stdinR, []string{"-"})
		if err != nil {
			t.Fatal(err)
		}

		if v, ok := m["foo"]; !ok || v != "bar" {
			t.Errorf("expected %q to be %q", v, "bar")
		}
	})

	t.Run("stdin_value", func(t *testing.T) {
		t.Parallel()

		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte(`bar`))
			stdinW.Close()
		}()

		m, err := parseArgsData(stdinR, []string{"foo=-"})
		if err != nil {
			t.Fatal(err)
		}

		if v, ok := m["foo"]; !ok || v != "bar" {
			t.Errorf("expected %q to be %q", v, "bar")
		}
	})

	t.Run("file_full", func(t *testing.T) {
		t.Parallel()

		f, err := ioutil.TempFile("", "vault")
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte(`{"foo":"bar"}`))
		f.Close()
		defer os.Remove(f.Name())

		m, err := parseArgsData(os.Stdin, []string{"@" + f.Name()})
		if err != nil {
			t.Fatal(err)
		}

		if v, ok := m["foo"]; !ok || v != "bar" {
			t.Errorf("expected %q to be %q", v, "bar")
		}
	})

	t.Run("file_value", func(t *testing.T) {
		t.Parallel()

		f, err := ioutil.TempFile("", "vault")
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte(`bar`))
		f.Close()
		defer os.Remove(f.Name())

		m, err := parseArgsData(os.Stdin, []string{"foo=@" + f.Name()})
		if err != nil {
			t.Fatal(err)
		}

		if v, ok := m["foo"]; !ok || v != "bar" {
			t.Errorf("expected %q to be %q", v, "bar")
		}
	})

	t.Run("file_value_escaped", func(t *testing.T) {
		t.Parallel()

		m, err := parseArgsData(os.Stdin, []string{`foo=\@`})
		if err != nil {
			t.Fatal(err)
		}

		if v, ok := m["foo"]; !ok || v != "@" {
			t.Errorf("expected %q to be %q", v, "@")
		}
	})
}

func TestTruncateToSeconds(t *testing.T) {
	t.Parallel()

	cases := []struct {
		d   time.Duration
		exp int
	}{
		{
			10 * time.Nanosecond,
			0,
		},
		{
			10 * time.Microsecond,
			0,
		},
		{
			10 * time.Millisecond,
			0,
		},
		{
			1 * time.Second,
			1,
		},
		{
			10 * time.Second,
			10,
		},
		{
			100 * time.Second,
			100,
		},
		{
			3 * time.Minute,
			180,
		},
		{
			3 * time.Hour,
			10800,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(fmt.Sprintf("%s", tc.d), func(t *testing.T) {
			t.Parallel()

			act := truncateToSeconds(tc.d)
			if act != tc.exp {
				t.Errorf("expected %d to be %d", act, tc.exp)
			}
		})
	}
}
