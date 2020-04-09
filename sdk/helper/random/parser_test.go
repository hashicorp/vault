package random

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type testCase struct {
		rawConfig string
		expected  StringGenerator
		expectErr bool
	}

	tests := map[string]testCase{
		"unrecognized rule": {
			rawConfig: `
				length = 20
				charset = "abcde"
				rule "testrule" {
					string = "teststring"
					int = 123
				}`,
			expected: StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
				Rules:   nil,
			},
			expectErr: true,
		},

		"charset restrictions": {
			rawConfig: `
				length = 20
				charset = "abcde"
				rule "CharsetRestriction" {
					charset = "abcde"
					min-chars = 2
				}`,
			expected: StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
				Rules: []Rule{
					&CharsetRestriction{
						Charset:  []rune("abcde"),
						MinChars: 2,
					},
				},
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := Parse(test.rawConfig)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("Actual: %#v\nExpected:%#v", actual, test.expected)
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	type testCase struct {
		registry map[string]ruleConstructor

		rawConfig string
		expected  StringGenerator
		expectErr bool
	}

	tests := map[string]testCase{
		"empty config": {
			registry:  defaultRuleNameMapping,
			rawConfig: "",
			expected:  StringGenerator{},
			expectErr: true,
		},
		"bogus config": {
			registry:  defaultRuleNameMapping,
			rawConfig: "asdf",
			expected:  StringGenerator{},
			expectErr: true,
		},
		"config with length and charset": {
			registry: defaultRuleNameMapping,
			rawConfig: `
				length = 20
				charset = "abcde"`,
			expected: StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
			},
			expectErr: false,
		},
		"config with zero length": {
			registry: defaultRuleNameMapping,
			rawConfig: `
				length = 0
				charset = "abcde"`,
			expected: StringGenerator{
				Length:  0,
				Charset: []rune("abcde"),
			},
			expectErr: true,
		},
		"config with negative length": {
			registry: defaultRuleNameMapping,
			rawConfig: `
				length = -2
				charset = "abcde"`,
			expected: StringGenerator{
				Length:  -2,
				Charset: []rune("abcde"),
			},
			expectErr: true,
		},
		"charset restrictions": {
			registry: defaultRuleNameMapping,
			rawConfig: `
				length = 20
				charset = "abcde"
				rule "CharsetRestriction" {
					charset = "abcde"
					min-chars = 2
				}`,
			expected: StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
				Rules: []Rule{
					&CharsetRestriction{
						Charset:  []rune("abcde"),
						MinChars: 2,
					},
				},
			},
			expectErr: false,
		},
		"test rule": {
			registry: map[string]ruleConstructor{
				"testrule": newTestRule,
			},
			rawConfig: `
				length = 20
				charset = "abcde"
				rule "testrule" {
					string = "teststring"
					int = 123
				}`,
			expected: StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
				Rules: []Rule{
					&testRule{
						String:  "teststring",
						Integer: 123,
					},
				},
			},
			expectErr: false,
		},
		"test rule and charset restrictions": {
			registry: map[string]ruleConstructor{
				"testrule":           newTestRule,
				"CharsetRestriction": ParseCharsetRestriction,
			},
			rawConfig: `
				length = 20
				charset = "abcde"
				rule "testrule" {
					string = "teststring"
					int = 123
				}
				rule "CharsetRestriction" {
					charset = "abcde"
					min-chars = 2
				}`,
			expected: StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
				Rules: []Rule{
					&testRule{
						String:  "teststring",
						Integer: 123,
					},
					&CharsetRestriction{
						Charset:  []rune("abcde"),
						MinChars: 2,
					},
				},
			},
			expectErr: false,
		},
		"unrecognized rule": {
			registry: defaultRuleNameMapping,
			rawConfig: `
				length = 20
				charset = "abcde"
				rule "testrule" {
					string = "teststring"
					int = 123
				}`,
			expected: StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
				Rules:   nil,
			},
			expectErr: true,
		},

		// /////////////////////////////////////////////////
		// JSON data
		"manually JSONified HCL": {
			registry: map[string]ruleConstructor{
				"testrule":           newTestRule,
				"CharsetRestriction": ParseCharsetRestriction,
			},
			rawConfig: `
				{
					"charset": "abcde",
					"length": 20,
					"rule": [
						{
							"testrule": [
								{
									"string": "teststring",
									"int": 123
								}
							]
						},
						{
							"CharsetRestriction": [
								{
									"charset": "abcde",
									"min-chars": 2
								}
							]
						}
					]
				}`,
			expected: StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
				Rules: []Rule{
					&testRule{
						String:  "teststring",
						Integer: 123,
					},
					&CharsetRestriction{
						Charset:  []rune("abcde"),
						MinChars: 2,
					},
				},
			},
			expectErr: false,
		},
		"JSONified HCL": {
			registry: map[string]ruleConstructor{
				"testrule":           newTestRule,
				"CharsetRestriction": ParseCharsetRestriction,
			},
			rawConfig: toJSON(t, StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
				Rules: []Rule{
					&testRule{
						String:  "teststring",
						Integer: 123,
					},
					&CharsetRestriction{
						Charset:  []rune("abcde"),
						MinChars: 2,
					},
				},
			}),
			expected: StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
				Rules: []Rule{
					&testRule{
						String:  "teststring",
						Integer: 123,
					},
					&CharsetRestriction{
						Charset:  []rune("abcde"),
						MinChars: 2,
					},
				},
			},
			expectErr: false,
		},
		"JSON unrecognized rule": {
			registry: defaultRuleNameMapping,
			rawConfig: `
				{
					"charset": "abcde",
					"length": 20,
					"rule": [
						{
							"testrule": [
								{
									"string": "teststring",
									"int": 123
								}
							],
						}
					]
				}`,
			expected: StringGenerator{
				Length:  20,
				Charset: []rune("abcde"),
				Rules:   nil,
			},
			expectErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			parser := Parser{
				RuleRegistry: Registry{
					Rules: test.registry,
				},
			}

			actual, err := parser.Parse(test.rawConfig)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("Actual: %#v\nExpected:%#v", actual, test.expected)
			}
		})
	}
}

func TestParseRules(t *testing.T) {
	type testCase struct {
		registry map[string]ruleConstructor

		rawRules      []map[string]interface{}
		expectedRules []Rule
		expectErr     bool
	}

	tests := map[string]testCase{
		"nil rule data": {
			registry:      defaultRuleNameMapping,
			rawRules:      nil,
			expectedRules: nil,
			expectErr:     false,
		},
		"empty rule data": {
			registry:      defaultRuleNameMapping,
			rawRules:      []map[string]interface{}{},
			expectedRules: nil,
			expectErr:     false,
		},
		"invalid rule data": {
			registry: defaultRuleNameMapping,
			rawRules: []map[string]interface{}{
				{
					"testrule": map[string]interface{}{
						"string": "teststring",
					},
				},
			},
			expectedRules: nil,
			expectErr:     true,
		},
		"unrecognized rule data": {
			registry: defaultRuleNameMapping,
			rawRules: []map[string]interface{}{
				{
					"testrule": []map[string]interface{}{
						{
							"string": "teststring",
							"int":    123,
						},
					},
				},
			},
			expectedRules: nil,
			expectErr:     true,
		},
		"recognized rule": {
			registry: map[string]ruleConstructor{
				"testrule": newTestRule,
			},
			rawRules: []map[string]interface{}{
				{
					"testrule": []map[string]interface{}{
						{
							"string": "teststring",
							"int":    123,
						},
					},
				},
			},
			expectedRules: []Rule{
				&testRule{
					String:  "teststring",
					Integer: 123,
				},
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			registry := Registry{
				Rules: test.registry,
			}

			actualRules, err := parseRules(registry, test.rawRules)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(actualRules, test.expectedRules) {
				t.Fatalf("Actual: %#v\nExpected:%#v", actualRules, test.expectedRules)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	type testCase struct {
		generator StringGenerator
		expectErr bool
	}

	tests := map[string]testCase{
		"default generator": {
			generator: DefaultStringGenerator,
			expectErr: false,
		},
		"length is 0": {
			generator: StringGenerator{
				Length:  0,
				Charset: []rune("abcde"),
			},
			expectErr: true,
		},
		"length is negative": {
			generator: StringGenerator{
				Length:  -2,
				Charset: []rune("abcde"),
			},
			expectErr: true,
		},
		"nil charset": {
			generator: StringGenerator{
				Length:  5,
				Charset: nil,
			},
			expectErr: true,
		},
		"zero length charset": {
			generator: StringGenerator{
				Length:  5,
				Charset: []rune{},
			},
			expectErr: true,
		},
		"rules require password longer than length": {
			generator: StringGenerator{
				Length:  5,
				Charset: []rune("abcde"),
				Rules: []Rule{
					CharsetRestriction{
						Charset:  []rune("abcde"),
						MinChars: 6,
					},
				},
			},
			expectErr: true,
		},
		"charset has non-printable characters": {
			generator: StringGenerator{
				Length: 0,
				Charset: []rune{
					'a',
					'b',
					0, // Null character
					'd',
					'e',
				},
			},
			expectErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := validate(test.generator)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}

func TestGetChars(t *testing.T) {
	type testCase struct {
		rules    []Rule
		expected []rune
	}

	tests := map[string]testCase{
		"nil rules": {
			rules:    nil,
			expected: []rune(nil),
		},
		"empty rules": {
			rules:    []Rule{},
			expected: []rune(nil),
		},
		"rule without chars": {
			rules: []Rule{
				testRule{
					String:  "teststring",
					Integer: 123,
				},
			},
			expected: []rune(nil),
		},
		"rule with chars": {
			rules: []Rule{
				CharsetRestriction{
					Charset:  []rune("abcdefghij"),
					MinChars: 1,
				},
			},
			expected: []rune("abcdefghij"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := getChars(test.rules)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("Actual: %v\nExpected: %v", actual, test.expected)
			}
		})
	}
}

func TestGetMapSlice(t *testing.T) {
	type testCase struct {
		input         map[string]interface{}
		key           string
		expectedSlice []map[string]interface{}
		expectErr     bool
	}

	tests := map[string]testCase{
		"nil map": {
			input:         nil,
			key:           "testkey",
			expectedSlice: nil,
			expectErr:     false,
		},
		"empty map": {
			input:         map[string]interface{}{},
			key:           "testkey",
			expectedSlice: nil,
			expectErr:     false,
		},
		"ignored keys": {
			input: map[string]interface{}{
				"foo": "bar",
			},
			key:           "testkey",
			expectedSlice: nil,
			expectErr:     false,
		},
		"key has wrong type": {
			input: map[string]interface{}{
				"foo": "bar",
			},
			key:           "foo",
			expectedSlice: nil,
			expectErr:     true,
		},
		"good data": {
			input: map[string]interface{}{
				"foo": []map[string]interface{}{
					{
						"sub-foo": "bar",
					},
				},
			},
			key: "foo",
			expectedSlice: []map[string]interface{}{
				{
					"sub-foo": "bar",
				},
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualSlice, err := getMapSlice(test.input, test.key)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(actualSlice, test.expectedSlice) {
				t.Fatalf("Actual: %#v\nExpected:%#v", actualSlice, test.expectedSlice)
			}
		})
	}
}

func TestGetRuleInfo(t *testing.T) {
	type testCase struct {
		rule         map[string]interface{}
		expectedInfo ruleInfo
		expectErr    bool
	}

	tests := map[string]testCase{
		"nil rule": {
			rule:         nil,
			expectedInfo: ruleInfo{},
			expectErr:    true,
		},
		"empty rule": {
			rule:         map[string]interface{}{},
			expectedInfo: ruleInfo{},
			expectErr:    true,
		},
		"rule with invalid type": {
			rule: map[string]interface{}{
				"TestRuleType": "wrong type",
			},
			expectedInfo: ruleInfo{},
			expectErr:    true,
		},
		"rule with good data": {
			rule: map[string]interface{}{
				"TestRuleType": []map[string]interface{}{
					{
						"foo": "bar",
					},
				},
			},
			expectedInfo: ruleInfo{
				ruleType: "TestRuleType",
				data: map[string]interface{}{
					"foo": "bar",
				},
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualInfo, err := getRuleInfo(test.rule)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(actualInfo, test.expectedInfo) {
				t.Fatalf("Actual: %#v\nExpected:%#v", actualInfo, test.expectedInfo)
			}
		})
	}
}

func TestApplyShortcuts(t *testing.T) {
	type testCase struct {
		input    map[string]interface{}
		expected map[string]interface{}
	}

	tests := map[string]testCase{
		"nil map": {
			input:    nil,
			expected: nil,
		},
		"empty map": {
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		"non-matching key": {
			input: map[string]interface{}{
				"foo": "teststring",
			},
			expected: map[string]interface{}{
				"foo": "teststring",
			},
		},
		"matching key": {
			input: map[string]interface{}{
				"charset": "lower-alpha",
			},
			expected: map[string]interface{}{
				"charset": LowercaseCharset,
			},
		},
		"matching and non-matching keys": {
			input: map[string]interface{}{
				"charset": "lower-alpha",
				"foo":     "teststring",
			},
			expected: map[string]interface{}{
				"charset": LowercaseCharset,
				"foo":     "teststring",
			},
		},
		"invalid value type": {
			input: map[string]interface{}{
				"charset": 123,
			},
			expected: map[string]interface{}{
				"charset": 123,
			},
		},
		"unrecognized shortcut": {
			input: map[string]interface{}{
				"charset": LowercaseCharset,
			},
			expected: map[string]interface{}{
				"charset": LowercaseCharset,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			applyShortcuts(test.input)
			if !reflect.DeepEqual(test.input, test.expected) {
				t.Fatalf("Actual: %#v\nExpected:%#v", test.input, test.expected)
			}
		})
	}
}

func TestDeduplicateRunes(t *testing.T) {
	type testCase struct {
		input    []rune
		expected []rune
	}

	tests := map[string]testCase{
		"empty string": {
			input:    []rune(""),
			expected: []rune(nil),
		},
		"no duplicates": {
			input:    []rune("abcde"),
			expected: []rune("abcde"),
		},
		"in order duplicates": {
			input:    []rune("aaaabbbbcccccccddddeeeee"),
			expected: []rune("abcde"),
		},
		"out of order duplicates": {
			input:    []rune("abcdeabcdeabcdeabcde"),
			expected: []rune("abcde"),
		},
		"unicode no duplicates": {
			input:    []rune("日本語"),
			expected: []rune("日本語"),
		},
		"unicode in order duplicates": {
			input:    []rune("日日日日本本本語語語語語"),
			expected: []rune("日本語"),
		},
		"unicode out of order duplicates": {
			input:    []rune("日本語日本語日本語日本語"),
			expected: []rune("日本語"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := deduplicateRunes(test.input)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("Actual: %#v\nExpected:%#v", actual, test.expected)
			}
		})
	}
}

func BenchmarkParser_Parse(b *testing.B) {
	config := `length = 20
               charset = "abcde"
               rule "CharsetRestriction" {
                   charset = "abcde"
                   min-chars = 2
               }`

	for i := 0; i < b.N; i++ {
		parser := Parser{
			RuleRegistry: Registry{
				Rules: defaultRuleNameMapping,
			},
		}
		_, err := parser.Parse(config)
		if err != nil {
			b.Fatalf("Failed to parse: %s", err)
		}
	}
}

func toJSON(t *testing.T, val interface{}) string {
	t.Helper()
	b, err := json.Marshal(val)
	if err != nil {
		t.Fatalf("unable to marshal to JSON: %s", err)
	}
	return string(b)
}
