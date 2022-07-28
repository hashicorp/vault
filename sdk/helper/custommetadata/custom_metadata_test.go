package custommetadata

import (
	"strconv"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name       string
		input      CustomMetadata
		shouldPass bool
	}{
		{
			"valid",
			CustomMetadata{
				"foo": "abc",
				"bar": "def",
				"baz": "ghi",
			},
			true,
		},
		{
			"too_many_keys",
			func() CustomMetadata {
				cm := make(CustomMetadata)

				for i := 0; i < maxKeyLength+1; i++ {
					s := strconv.Itoa(i)
					cm[s] = s
				}

				return cm
			}(),
			false,
		},
		{
			"key_too_long",
			CustomMetadata{
				strings.Repeat("a", maxKeyLength+1): "abc",
			},
			false,
		},
		{
			"value_too_long",
			CustomMetadata{
				"foo": strings.Repeat("a", maxValueLength+1),
			},
			false,
		},
		{
			"unprintable_key",
			CustomMetadata{
				"unprint\u200bable": "abc",
			},
			false,
		},
		{
			"unprintable_value",
			CustomMetadata{
				"foo": "unprint\u200bable",
			},
			false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := Validate(tc.input)

			if tc.shouldPass && err != nil {
				t.Fatalf("expected validation to pass, input: %#v, err: %v", tc.input, err)
			}

			if !tc.shouldPass && err == nil {
				t.Fatalf("expected validation to fail, input: %#v, err: %v", tc.input, err)
			}
		})
	}
}
