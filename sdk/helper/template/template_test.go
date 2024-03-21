// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package template

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	type testCase struct {
		template       string
		additionalOpts []Opt
		data           interface{}

		expected  string
		expectErr bool
	}

	tests := map[string]testCase{
		"template without arguments": {
			template:  "this is a template",
			data:      nil,
			expected:  "this is a template",
			expectErr: false,
		},
		"template with arguments but no data": {
			template:  "this is a {{.String}}",
			data:      nil,
			expected:  "this is a <no value>",
			expectErr: false,
		},
		"template with arguments": {
			template: "this is a {{.String}}",
			data: struct {
				String string
			}{
				String: "foobar",
			},
			expected:  "this is a foobar",
			expectErr: false,
		},
		"template with builtin functions": {
			template: `{{.String | truncate 10}}
{{.String | uppercase}}
{{.String | lowercase}}
{{.String | replace " " "."}}
{{.String | sha256}}
{{.String | base64}}
{{.String | truncate_sha256 20}}`,
			data: struct {
				String string
			}{
				String: "Some string with Multiple Capitals LETTERS",
			},
			expected: `Some strin
SOME STRING WITH MULTIPLE CAPITALS LETTERS
some string with multiple capitals letters
Some.string.with.Multiple.Capitals.LETTERS
da9872dd96609c72897defa11fe81017a62c3f44339d9d3b43fe37540ede3601
U29tZSBzdHJpbmcgd2l0aCBNdWx0aXBsZSBDYXBpdGFscyBMRVRURVJT
Some string 6841cf80`,
			expectErr: false,
		},
		"custom function": {
			template: "{{foo}}",
			additionalOpts: []Opt{
				Function("foo", func() string {
					return "custom-foo"
				}),
			},
			expected:  "custom-foo",
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			opts := append(test.additionalOpts, Template(test.template))
			st, err := NewTemplate(opts...)
			require.NoError(t, err)

			actual, err := st.Generate(test.data)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			require.Equal(t, test.expected, actual)
		})
	}

	t.Run("random", func(t *testing.T) {
		for i := 1; i < 100; i++ {
			st, err := NewTemplate(
				Template(fmt.Sprintf("{{random %d}}", i)),
			)
			require.NoError(t, err)

			actual, err := st.Generate(nil)
			require.NoError(t, err)

			require.Regexp(t, fmt.Sprintf("^[a-zA-Z0-9]{%d}$", i), actual)
		}
	})

	t.Run("unix_time", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			st, err := NewTemplate(
				Template("{{unix_time}}"),
			)
			require.NoError(t, err)

			actual, err := st.Generate(nil)
			require.NoError(t, err)

			require.Regexp(t, "^[0-9]+$", actual)
		}
	})

	t.Run("unix_time_millis", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			st, err := NewTemplate(
				Template("{{unix_time_millis}}"),
			)
			require.NoError(t, err)

			actual, err := st.Generate(nil)
			require.NoError(t, err)

			require.Regexp(t, "^[0-9]+$", actual)
		}
	})

	t.Run("timestamp", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			st, err := NewTemplate(
				Template(`{{timestamp "2006-01-02T15:04:05.000Z"}}`),
			)
			require.NoError(t, err)

			actual, err := st.Generate(nil)
			require.NoError(t, err)

			require.Regexp(t, `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`, actual)
		}
	})

	t.Run("too-large-overflow", func(t *testing.T) {
		data := "{{" + strings.Repeat("(", 1000000)
		_, err := NewTemplate(
			Template(data),
		)
		// We expect an error due it being too large,
		// this test should not fail with an overflow
		require.Error(t, err)
	})
}

func TestBadConstructorArguments(t *testing.T) {
	type testCase struct {
		opts []Opt
	}

	tests := map[string]testCase{
		"missing template": {
			opts: nil,
		},
		"missing custom function name": {
			opts: []Opt{
				Template("foo bar"),
				Function("", func() string {
					return "foo"
				}),
			},
		},
		"missing custom function": {
			opts: []Opt{
				Template("foo bar"),
				Function("foo", nil),
			},
		},
		"bad template": {
			opts: []Opt{
				Template("{{.String"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			st, err := NewTemplate(test.opts...)
			require.Error(t, err)

			str, err := st.Generate(nil)
			require.Error(t, err)
			require.Equal(t, "", str)
		})
	}

	t.Run("erroring custom function", func(t *testing.T) {
		st, err := NewTemplate(
			Template("{{foo}}"),
			Function("foo", func() (string, error) {
				return "", fmt.Errorf("an error!")
			}),
		)
		require.NoError(t, err)

		str, err := st.Generate(nil)
		require.Error(t, err)
		require.Equal(t, "", str)
	})
}
