package credsutil

import (
	"regexp"
	"testing"
)

func TestGeneratePassword_Constructor(t *testing.T) {
	type testCase struct {
		opts        []PasswordOpt
		expectRegex string
	}

	tests := map[string]testCase{
		"default": {
			opts:        []PasswordOpt{},
			expectRegex: "^A1a-[a-zA-Z0-9]{11}$",
		},
		"charset lowercase only": {
			opts: []PasswordOpt{
				PasswordCharset(LowerCharset),
			},
			expectRegex: "^A1a-[a-z]{11}$",
		},
		"charset uppercase only": {
			opts: []PasswordOpt{
				PasswordCharset(UpperCharset),
			},
			expectRegex: "^A1a-[A-Z]{11}$",
		},
		"charset numbers only": {
			opts: []PasswordOpt{
				PasswordCharset(NumericCharset),
			},
			expectRegex: "^A1a-[0-9]{11}$",
		},
		"long length": {
			opts: []PasswordOpt{
				PasswordLength(100),
			},
			expectRegex: "^A1a-[a-zA-Z0-9]{96}$",
		},
		"empty prefix": {
			opts: []PasswordOpt{
				PasswordPrefix(""),
			},
			expectRegex: "^[a-zA-Z0-9]{15}$",
		},
		"non-default prefix": {
			opts: []PasswordOpt{
				PasswordPrefix("abcde"),
			},
			expectRegex: "^abcde[a-zA-Z0-9]{10}$",
		},
		"mixed opts": {
			opts: []PasswordOpt{
				PasswordCharset(UpperCharset),
				PasswordLength(20),
				PasswordPrefix("abcde"),
			},
			expectRegex: "^abcde[A-Z]{15}$",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p, err := NewPasswordProducer(test.opts...)
			if err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
			re := regexp.MustCompile(test.expectRegex)

			// Because this is a randomly generated value, let's run this a bunch of times to detect any flakiness
			for i := 0; i < 1000; i++ {
				password, err := p.GeneratePassword()
				if err != nil {
					t.Fatalf("no error expected, got: %s", err)
				}

				if !re.MatchString(password) {
					t.Fatalf("Password [%s] did not match regex [%s]", password, test.expectRegex)
				}
			}
		})
	}
}

func TestGeneratePassword_NoConstructor(t *testing.T) {
	type testCase struct {
		charset     string
		length      int
		prefix      string
		expectRegex string
	}

	tests := map[string]testCase{
		"default": {
			expectRegex: "^[a-zA-Z0-9]{15}$",
		},
		"charset lowercase only": {
			charset:     LowerCharset,
			expectRegex: "^[a-z]{15}$",
		},
		"charset uppercase only": {
			charset:     UpperCharset,
			expectRegex: "^[A-Z]{15}$",
		},
		"charset numbers only": {
			charset:     NumericCharset,
			expectRegex: "^[0-9]{15}$",
		},
		"long length": {
			length:      100,
			expectRegex: "^[a-zA-Z0-9]{100}$",
		},
		"empty prefix": {
			prefix:      "",
			expectRegex: "^[a-zA-Z0-9]{15}$",
		},
		"non-default prefix": {
			prefix:      "abcde",
			expectRegex: "^abcde[a-zA-Z0-9]{10}$",
		},
		"mixed opts": {
			charset:     UpperCharset,
			length:      20,
			prefix:      "abcde",
			expectRegex: "^abcde[A-Z]{15}$",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			p := PasswordProducer{
				Charset: test.charset,
				Length:  test.length,
				Prefix:  test.prefix,
			}
			re := regexp.MustCompile(test.expectRegex)

			// Because this is a randomly generated value, let's run this a bunch of times to detect any flakiness
			for i := 0; i < 1000; i++ {
				password, err := p.GeneratePassword()
				if err != nil {
					t.Fatalf("no error expected, got: %s", err)
				}

				if !re.MatchString(password) {
					t.Fatalf("Password [%s] did not match regex [%s]", password, test.expectRegex)
				}
			}
		})
	}
}

func TestPasswordProducer_Constructor(t *testing.T) {
	type testCase struct {
		opts      []PasswordOpt
		expectErr bool
	}

	tests := map[string]testCase{
		"password length too short": {
			opts: []PasswordOpt{
				PasswordLength(MinPasswordLength - 1),
			},
			expectErr: true,
		},
		"password length at minimum": {
			opts: []PasswordOpt{
				PasswordLength(MinPasswordLength),
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewPasswordProducer(test.opts...)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}
