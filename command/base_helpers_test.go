// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
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
		f.WriteString(`{"foo":"bar"}`)
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
		f.WriteString(`bar`)
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

		t.Run(tc.d.String(), func(t *testing.T) {
			t.Parallel()

			act := truncateToSeconds(tc.d)
			if act != tc.exp {
				t.Errorf("expected %d to be %d", act, tc.exp)
			}
		})
	}
}

func TestParseFlagFile(t *testing.T) {
	t.Parallel()

	content := "some raw content"
	tmpFile, err := ioutil.TempFile(os.TempDir(), "TestParseFlagFile")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}

	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	cases := []struct {
		value string
		exp   string
	}{
		{
			"",
			"",
		},
		{
			content,
			content,
		},
		{
			fmt.Sprintf("@%s", tmpFile.Name()),
			content,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.value, func(t *testing.T) {
			content, err := ParseFlagFile(tc.value)
			if err != nil {
				t.Fatalf("unexpected error parsing flag value: %v", err)
			}

			if content != tc.exp {
				t.Fatalf("expected %s to be %s", content, tc.exp)
			}
		})
	}
}

func TestArgWarnings(t *testing.T) {
	t.Parallel()

	cases := []struct {
		args     []string
		expected string
	}{
		{
			[]string{"a", "b", "c"},
			"",
		},
		{
			[]string{"a", "-b"},
			"-b",
		},
		{
			[]string{"a", "--b"},
			"--b",
		},
		{
			[]string{"a-b", "-c"},
			"-c",
		},
		{
			[]string{"a", "-b-c"},
			"-b-c",
		},
		{
			[]string{"-a", "b"},
			"-a",
		},
		{
			[]string{globalFlagDetailed},
			"",
		},
		{
			[]string{"-" + globalFlagOutputCurlString + "=true"},
			"",
		},
		{
			[]string{"--" + globalFlagFormat + "=false"},
			"",
		},
		{
			[]string{"-x" + globalFlagDetailed},
			"-x" + globalFlagDetailed,
		},
		{
			[]string{"--x=" + globalFlagDetailed},
			"--x=" + globalFlagDetailed,
		},
		{
			[]string{"policy", "write", "my-policy", "-"},
			"",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.expected, func(t *testing.T) {
			warnings := generateFlagWarnings(tc.args)
			if !strings.Contains(warnings, tc.expected) {
				t.Fatalf("expected %s to contain %s", warnings, tc.expected)
			}
		})
	}
}
