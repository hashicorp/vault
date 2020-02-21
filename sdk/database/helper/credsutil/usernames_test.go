package credsutil

import (
	"regexp"
	"testing"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
)

func TestRemoveDuplicateChars(t *testing.T) {
	type testCase struct {
		input    string
		expected string
	}

	tests := map[string]testCase{
		"DefaultCharset": {
			input:    DefaultCharset,
			expected: DefaultCharset,
		},
		"UpperCharset": {
			input:    UpperCharset,
			expected: UpperCharset,
		},
		"LowerCharset": {
			input:    LowerCharset,
			expected: LowerCharset,
		},
		"NumericCharset": {
			input:    NumericCharset,
			expected: NumericCharset,
		},
		"SymbolsCharset": {
			input:    SymbolsCharset,
			expected: SymbolsCharset,
		},
		"Duplicate default charset": {
			input:    DefaultCharset + DefaultCharset,
			expected: DefaultCharset,
		},
		"many duplicates": {
			input:    "aaaaaaaaaaaaaaaaaaa",
			expected: "a",
		},
		"order is kept": {
			input:    "abcdefabcdefabcdefabcdef",
			expected: "abcdef",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := removeDuplicateChars(test.input)
			if actual != test.expected {
				t.Fatalf("Actual: [%s]\nExpected: [%s]", actual, test.expected)
			}
		})
	}
}

func BenchmarkRemoveDuplicateChars(b *testing.B) {
	for i := 0; i < b.N; i++ {
		removeDuplicateChars(DefaultCharset + SymbolsCharset)
	}
}

func TestRandCharset(t *testing.T) {
	type testCase struct {
		charset     string
		length      int
		expectRegex string
		expectErr   bool
	}

	tests := map[string]testCase{
		"default charset": {
			charset:     DefaultCharset,
			length:      10,
			expectRegex: "^[a-zA-Z0-9]{10}$",
			expectErr:   false,
		},
		"unusual charset": {
			charset:     "abcdefg01234",
			length:      20,
			expectRegex: "^[a-g0-4]{20}$",
			expectErr:   false,
		},
		"very long symbols": {
			charset:     SymbolsCharset,
			length:      200,
			expectRegex: "^[!\"#$%&'\\(\\)\\*\\+,-./:;<=>?`\\~\\[\\]^_@|\\\\]{200}$",
			expectErr:   false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Because the value is random, run this a bunch of times to ensure that each run has the expected behavior
			for i := 0; i < 100; i++ {
				actual, err := randCharset(test.charset)(test.length)
				if test.expectErr && err == nil {
					t.Fatalf("err expected, got nil")
				}
				if !test.expectErr && err != nil {
					t.Fatalf("no error expected, got: %s", err)
				}
				expectedRegex := regexp.MustCompile(test.expectRegex)
				if !expectedRegex.MatchString(actual) {
					t.Fatalf("Random characters [%s] did not match regexp [%s]", actual, test.expectRegex)
				}
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	type testCase struct {
		input     string
		maxLen    int
		expected  string
		expectErr bool
	}

	tests := map[string]testCase{
		"no truncate": {
			input:     "foobar",
			maxLen:    10,
			expected:  "foobar",
			expectErr: false,
		},
		"max len equals length": {
			input:     "foobar",
			maxLen:    6,
			expected:  "foobar",
			expectErr: false,
		},
		"max len equals length-1": {
			input:     "foobar",
			maxLen:    5,
			expected:  "fooba",
			expectErr: false,
		},
		"max len 1": {
			input:     "foobar",
			maxLen:    1,
			expected:  "f",
			expectErr: false,
		},
		"max len is zero": {
			input:     "foobar",
			maxLen:    0,
			expected:  "foobar",
			expectErr: true,
		},
		"max len is negative": {
			input:     "foobar",
			maxLen:    -10,
			expected:  "foobar",
			expectErr: true,
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
			if actual != test.expected {
				t.Fatalf("Actual [%s] Expected [%s]", actual, test.expected)
			}
		})
	}
}

func TestGenerateUsername(t *testing.T) {
	type testCase struct {
		opts        []UsernameOpt
		config      dbplugin.UsernameConfig
		expectRegex string
	}

	tests := map[string]testCase{
		"default values": {
			config: dbplugin.UsernameConfig{
				DisplayName: "dispname",
				RoleName:    "rolename",
			},
			expectRegex: `^v_rolename_dispname_[0-9]+_[a-zA-Z0-9]{10}$`,
		},
		"multiple roles": {
			opts: []UsernameOpt{
				UsernameTemplate("{{.RoleName}}_{{.RoleName}}"),
			},
			config: dbplugin.UsernameConfig{
				DisplayName: "dispname",
				RoleName:    "rolename",
			},
			expectRegex: `^rolename_rolename$`,
		},
		"truncate": {
			opts: []UsernameOpt{
				UsernameTemplate("{{.RoleName | truncate 5}}"),
			},
			config: dbplugin.UsernameConfig{
				DisplayName: "dispname",
				RoleName:    "rolename",
			},
			expectRegex: `^rolen$`,
		},
		"multiple displays": {
			opts: []UsernameOpt{
				UsernameTemplate("{{.DisplayName}}_{{.DisplayName}}"),
			},
			config: dbplugin.UsernameConfig{
				DisplayName: "dispname",
				RoleName:    "rolename",
			},
			expectRegex: `^dispname_dispname$`,
		},
		"random": {
			opts: []UsernameOpt{
				UsernameTemplate("{{rand 10}}"),
			},
			config: dbplugin.UsernameConfig{
				DisplayName: "dispname",
				RoleName:    "rolename",
			},
			expectRegex: `^[a-zA-Z0-9]{10}$`,
		},
		"multiple randoms": {
			opts: []UsernameOpt{
				UsernameTemplate("{{rand 10}}_{{rand 20}}"),
			},
			config: dbplugin.UsernameConfig{
				DisplayName: "dispname",
				RoleName:    "rolename",
			},
			expectRegex: `^[a-zA-Z0-9]{10}_[a-zA-Z0-9]{20}$`,
		},
		"random with custom charset": {
			opts: []UsernameOpt{
				UsernameTemplate("{{rand 30}}"),
				UsernameFuncMap("rand", randCharset("abcdefg012345")),
			},
			config: dbplugin.UsernameConfig{
				DisplayName: "dispname",
				RoleName:    "rolename",
			},
			expectRegex: `^[a-g0-5]{30}$`,
		},
		"mix and match": {
			opts: []UsernameOpt{
				UsernameTemplate("Prefix_{{.DisplayName | truncate 5}}-{{rand 10}}={{.RoleName}}_suffix"),
			},
			config: dbplugin.UsernameConfig{
				DisplayName: "dispname",
				RoleName:    "rolename",
			},
			expectRegex: `^Prefix_dispn-[a-zA-Z0-9]{10}=rolename_suffix$`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			up, err := NewUsernameProducer(test.opts...)
			if err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			username, err := up.GenerateUsername(test.config)
			if err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			expectedRegex := regexp.MustCompile(test.expectRegex)
			if !expectedRegex.MatchString(username) {
				t.Fatalf("Username [%s] did not match regexp [%s]", username, test.expectRegex)
			}
		})
	}
}

func TestUsernameBackwardsCompatibility(t *testing.T) {
	type testCase struct {
		dispLen   int
		roleLen   int
		maxLen    int
		separator string
		lowercase bool

		template string

		expectRegex string
	}

	tests := map[string]testCase{
		"cassandra & influxdb": {
			dispLen:     15,
			roleLen:     15,
			maxLen:      100,
			separator:   "_",
			lowercase:   false,
			template:    "v_{{.DisplayName | truncate 15}}_{{.RoleName | truncate 15}}_{{rand 20}}_{{now_seconds}}",
			expectRegex: "^v_aBcDeFgHiJkLmNo_012345678901234_[a-zA-Z0-9]{20}_[0-9]+$",
		},
		"hana": {
			dispLen:     32,
			roleLen:     20,
			maxLen:      128,
			separator:   "_",
			lowercase:   false,
			template:    "v_{{.DisplayName | truncate 32}}_{{.RoleName | truncate 20}}_{{rand 20}}_{{now_seconds}}",
			expectRegex: "^v_aBcDeFgHiJkLmNoPqRsTuVwXyZaBcDeF_01234567890123456789_[a-zA-Z0-9]{20}_[0-9]+$",
		},
		"mongodb": {
			dispLen:     15,
			roleLen:     15,
			maxLen:      100,
			separator:   "-",
			lowercase:   false,
			template:    "v-{{.DisplayName | truncate 15}}-{{.RoleName | truncate 15}}-{{rand 20}}-{{now_seconds}}",
			expectRegex: "^v-aBcDeFgHiJkLmNo-012345678901234-[a-zA-Z0-9]{20}-[0-9]+$",
		},
		"mssql": {
			dispLen:     20,
			roleLen:     20,
			maxLen:      128,
			separator:   "-",
			lowercase:   false,
			template:    "v-{{.DisplayName | truncate 20}}-{{.RoleName | truncate 20}}-{{rand 20}}-{{now_seconds}}",
			expectRegex: "^v-aBcDeFgHiJkLmNoPqRsT-01234567890123456789-[a-zA-Z0-9]{20}-[0-9]+$",
		},
		"mysql": {
			dispLen:     10,
			roleLen:     10,
			maxLen:      32,
			separator:   "-",
			lowercase:   false,
			template:    "v-{{.DisplayName | truncate 10}}-{{.RoleName | truncate 10}}-{{rand 20}}-{{now_seconds}}",
			expectRegex: "^v-aBcDeFgHiJ-0123456789-[a-zA-Z0-9]{8}$",
		},
		"mysql alternatives": {
			dispLen:     -1,
			roleLen:     4,
			maxLen:      16,
			separator:   "-",
			lowercase:   false,
			template:    "v-{{.RoleName | truncate 4}}-{{rand 20}}-{{now_seconds}}",
			expectRegex: "^v-0123-[a-zA-Z0-9]{9}$",
		},
		"postgresql": {
			dispLen:     8,
			roleLen:     8,
			maxLen:      63,
			separator:   "-",
			lowercase:   false,
			template:    "v-{{.DisplayName | truncate 8}}-{{.RoleName | truncate 8}}-{{rand 20}}-{{now_seconds}}",
			expectRegex: "^v-aBcDeFgH-01234567-[a-zA-Z0-9]{20}-[0-9]+$",
		},
		"redshift": {
			dispLen:     8,
			roleLen:     8,
			maxLen:      63,
			separator:   "-",
			lowercase:   false,
			template:    "v-{{.DisplayName | truncate 8}}-{{.RoleName | truncate 8}}-{{rand 20}}-{{now_seconds}}",
			expectRegex: "^v-aBcDeFgH-01234567-[a-zA-Z0-9]{20}-[0-9]+$",
		},
		"redshift lowercase": {
			dispLen:     8,
			roleLen:     8,
			maxLen:      63,
			separator:   "-",
			lowercase:   true,
			template:    "v-{{.DisplayName | truncate 8 | lowercase}}-{{.RoleName | truncate 8 | lowercase}}-{{rand 20 | lowercase}}-{{now_seconds}}",
			expectRegex: "^v-abcdefgh-01234567-[a-z0-9]{20}-[0-9]+$",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sqlCredProducer := &SQLCredentialsProducer{
				DisplayNameLen:    test.dispLen,
				RoleNameLen:       test.roleLen,
				UsernameLen:       test.maxLen,
				Separator:         test.separator,
				LowercaseUsername: test.lowercase,
			}

			t.Logf("Template: %s", test.template)

			up, err := NewUsernamePasswordProducer(
				UsernameOpts(
					UsernameTemplate(test.template),
					UsernameMaxLength(test.maxLen),
				),
			)
			if err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			config := dbplugin.UsernameConfig{
				DisplayName: "aBcDeFgHiJkLmNoPqRsTuVwXyZaBcDeFgHiJkLmNoPqRsTuVwXyZ",
				RoleName:    "0123456789012345678901234567890123456789012345678901",
			}
			sqlUsername, err := sqlCredProducer.GenerateUsername(config)
			if err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			upUsername, err := up.GenerateUsername(config)
			if err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			re := regexp.MustCompile(test.expectRegex)
			if !re.MatchString(sqlUsername) {
				t.Fatalf("SQL username [%s] did not match regex [%s]", sqlUsername, test.expectRegex)
			}
			if !re.MatchString(upUsername) {
				t.Fatalf("UsernameProducer username [%s] did not match regex [%s]", upUsername, test.expectRegex)
			}
		})
	}
}
