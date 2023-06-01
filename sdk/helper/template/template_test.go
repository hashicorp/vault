// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package template

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const templateCharSet = "abcdefghijklmnopqrstuvwxyz012346789"

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

func TestTemplateInputLength(t *testing.T) {
	type testCase struct {
		name        string
		length      int
		wantErr     bool
		expectedErr string
	}

	tests := []testCase{
		{
			name:    "below length limit",
			length:  100,
			wantErr: false,
		},
		{
			name:   "at length limit",
			length: maxTemplateInputLength,

			wantErr: false,
		},
		{
			name:        "exceeds length limit",
			length:      maxTemplateInputLength + 1,
			wantErr:     true,
			expectedErr: "exceeds the desired length limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := randString(tt.length)

			_, err := NewTemplate(Template(input))
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.expectedErr) {
				t.Fatalf("expected error %s, got %s", tt.expectedErr, err.Error())
			}
		})
	}
}

func randString(strlen int) string {
	return randStringFromCharSet(strlen, templateCharSet)
}

// RandStringFromCharSet generates a random string by selecting characters from
// the charset provided
func randStringFromCharSet(strlen int, charSet string) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}

// Fuzz test to test different inputs to NewTemplate
func FuzzNewTemplate(f *testing.F) {
	template := `{{ if (eq .Type "STS") }}{{ printf "vault-%s-%s"  (unix_time) (random 20) | truncate 32 }}{{ else }}{{ printf "vault-%s-%s-%s" (printf "%s-%s" (.DisplayName) (.PolicyName) | truncate 42) (unix_time) (random 20) | truncate 64 }}{{ end }}`
	f.Add(template)
	f.Fuzz(func(t *testing.T, input string) {
		st, _ := NewTemplate(Template(input))
		st.Generate(nil)
	})
}
