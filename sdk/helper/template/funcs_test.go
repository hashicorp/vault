package template

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnixTimestamp(t *testing.T) {
	now := time.Now().Unix()
	for i := 0; i < 100; i++ {
		str := unixTime()
		actual, err := strconv.Atoi(str)
		require.NoError(t, err)
		// Make sure the value generated is from now (or later if the clock ticked over)
		require.GreaterOrEqual(t, int64(actual), now)
	}
}

func TestNowNano(t *testing.T) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	for i := 0; i < 100; i++ {
		str := unixTimeMillis()
		actual, err := strconv.ParseUint(str, 10, 64)
		require.NoError(t, err)
		// Make sure the value generated is from now (or later if the clock ticked over)
		require.GreaterOrEqual(t, int64(actual), now)
	}
}

func TestTruncate(t *testing.T) {
	type testCase struct {
		maxLen    int
		input     string
		expected  string
		expectErr bool
	}

	tests := map[string]testCase{
		"negative max length": {
			maxLen:    -1,
			input:     "foobarbaz",
			expected:  "",
			expectErr: true,
		},
		"zero max length": {
			maxLen:    0,
			input:     "foobarbaz",
			expected:  "",
			expectErr: true,
		},
		"one max length": {
			maxLen:    1,
			input:     "foobarbaz",
			expected:  "f",
			expectErr: false,
		},
		"half max length": {
			maxLen:    5,
			input:     "foobarbaz",
			expected:  "fooba",
			expectErr: false,
		},
		"max length one less than length": {
			maxLen:    8,
			input:     "foobarbaz",
			expected:  "foobarba",
			expectErr: false,
		},
		"max length equals string length": {
			maxLen:    9,
			input:     "foobarbaz",
			expected:  "foobarbaz",
			expectErr: false,
		},
		"max length greater than string length": {
			maxLen:    10,
			input:     "foobarbaz",
			expected:  "foobarbaz",
			expectErr: false,
		},
		"max length significantly greater than string length": {
			maxLen:    100,
			input:     "foobarbaz",
			expected:  "foobarbaz",
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := truncate(test.maxLen, test.input)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			require.Equal(t, test.expected, actual)
		})
	}
}

func TestTruncateSHA256(t *testing.T) {
	type testCase struct {
		maxLen    int
		input     string
		expected  string
		expectErr bool
	}

	tests := map[string]testCase{
		"negative max length": {
			maxLen:    -1,
			input:     "thisisareallylongstring",
			expected:  "",
			expectErr: true,
		},
		"zero max length": {
			maxLen:    0,
			input:     "thisisareallylongstring",
			expected:  "",
			expectErr: true,
		},
		"8 max length": {
			maxLen:    8,
			input:     "thisisareallylongstring",
			expected:  "",
			expectErr: true,
		},
		"nine max length": {
			maxLen:    9,
			input:     "thisisareallylongstring",
			expected:  "t4bb25641",
			expectErr: false,
		},
		"half max length": {
			maxLen:    12,
			input:     "thisisareallylongstring",
			expected:  "this704cd12b",
			expectErr: false,
		},
		"max length one less than length": {
			maxLen:    22,
			input:     "thisisareallylongstring",
			expected:  "thisisareallyl7f978be6",
			expectErr: false,
		},
		"max length equals string length": {
			maxLen:    23,
			input:     "thisisareallylongstring",
			expected:  "thisisareallylongstring",
			expectErr: false,
		},
		"max length greater than string length": {
			maxLen:    24,
			input:     "thisisareallylongstring",
			expected:  "thisisareallylongstring",
			expectErr: false,
		},
		"max length significantly greater than string length": {
			maxLen:    100,
			input:     "thisisareallylongstring",
			expected:  "thisisareallylongstring",
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := truncateSHA256(test.maxLen, test.input)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			require.Equal(t, test.expected, actual)
		})
	}
}

func TestSHA256(t *testing.T) {
	type testCase struct {
		input    string
		expected string
	}

	tests := map[string]testCase{
		"empty string": {
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		"foobar": {
			input:    "foobar",
			expected: "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2",
		},
		"mystring": {
			input:    "mystring",
			expected: "bd3ff47540b31e62d4ca6b07794e5a886b0f655fc322730f26ecd65cc7dd5c90",
		},
		"very long string": {
			input: "Nullam pharetra mattis laoreet. Mauris feugiat, tortor in malesuada convallis, " +
				"eros nunc dapibus erat, eget malesuada purus leo id lorem. Morbi pharetra, libero at malesuada bibendum, " +
				"dui quam tristique libero, bibendum cursus diam quam at sem. Vivamus vestibulum orci vel odio posuere, " +
				"quis tincidunt ipsum lacinia. Donec elementum a orci quis lobortis. Etiam bibendum ullamcorper varius. " +
				"Mauris tempor eros est, at porta erat rutrum ac. Aliquam erat volutpat. Sed sagittis leo non bibendum " +
				"lacinia. Praesent id justo iaculis, mattis libero vel, feugiat dui. Morbi id diam non magna imperdiet " +
				"imperdiet. Ut tortor arcu, mollis ac maximus ac, sagittis commodo augue. Ut semper, diam pulvinar porta " +
				"dignissim, massa ex condimentum enim, sed euismod urna quam vitae ex. Sed id neque vitae magna sagittis " +
				"pretium. Suspendisse potenti.",
			expected: "3e2a996c20b7a02378204f0843507d335e1ba203df2c4ded8d839d44af24482f",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := hashSHA256(test.input)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestUppercase(t *testing.T) {
	type testCase struct {
		input    string
		expected string
	}

	tests := map[string]testCase{
		"empty string": {
			input:    "",
			expected: "",
		},
		"lowercase": {
			input:    "foobar",
			expected: "FOOBAR",
		},
		"uppercase": {
			input:    "FOOBAR",
			expected: "FOOBAR",
		},
		"mixed case": {
			input:    "fOoBaR",
			expected: "FOOBAR",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := uppercase(test.input)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestLowercase(t *testing.T) {
	type testCase struct {
		input    string
		expected string
	}

	tests := map[string]testCase{
		"empty string": {
			input:    "",
			expected: "",
		},
		"lowercase": {
			input:    "foobar",
			expected: "foobar",
		},
		"uppercase": {
			input:    "FOOBAR",
			expected: "foobar",
		},
		"mixed case": {
			input:    "fOoBaR",
			expected: "foobar",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := lowercase(test.input)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestReplace(t *testing.T) {
	type testCase struct {
		input    string
		find     string
		replace  string
		expected string
	}

	tests := map[string]testCase{
		"empty string": {
			input:    "",
			find:     "",
			replace:  "",
			expected: "",
		},
		"search not found": {
			input:    "foobar",
			find:     ".",
			replace:  "_",
			expected: "foobar",
		},
		"single character found": {
			input:    "foo.bar",
			find:     ".",
			replace:  "_",
			expected: "foo_bar",
		},
		"multiple characters found": {
			input:    "foo.bar.baz",
			find:     ".",
			replace:  "_",
			expected: "foo_bar_baz",
		},
		"find and remove": {
			input:    "foo.bar",
			find:     ".",
			replace:  "",
			expected: "foobar",
		},
		"find full string": {
			input:    "foobarbaz",
			find:     "bar",
			replace:  "_",
			expected: "foo_baz",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := replace(test.find, test.replace, test.input)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestUUID(t *testing.T) {
	re := "^[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}$"
	for i := 0; i < 100; i++ {
		id, err := uuid()
		require.NoError(t, err)
		require.Regexp(t, re, id)
	}
}
