// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package random

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParsePolicy(t *testing.T) {
	type testCase struct {
		rawConfig string
		expected  StringGenerator
		expectErr bool
	}

	tests := map[string]testCase{
		"unrecognized rule": {
			rawConfig: `
				length = 20
				rule "testrule" {
					string = "teststring"
					int = 123
				}`,
			expected:  StringGenerator{},
			expectErr: true,
		},

		"charset restrictions": {
			rawConfig: `
				length = 20
				rule "charset" {
					charset = "abcde"
					min-chars = 2
				}`,
			expected: StringGenerator{
				Length:  20,
				charset: []rune("abcde"),
				Rules: []Rule{
					CharsetRule{
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
			actual, err := ParsePolicy(test.rawConfig)
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

func TestParser_ParsePolicy(t *testing.T) {
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
		"config with only length": {
			registry: defaultRuleNameMapping,
			rawConfig: `
				length = 20`,
			expected:  StringGenerator{},
			expectErr: true,
		},
		"config with zero length": {
			registry: defaultRuleNameMapping,
			rawConfig: `
				length = 0
				rule "charset" {
					charset = "abcde"
				}`,
			expected:  StringGenerator{},
			expectErr: true,
		},
		"config with negative length": {
			registry: defaultRuleNameMapping,
			rawConfig: `
				length = -2
				rule "charset" {
					charset = "abcde"
				}`,
			expected:  StringGenerator{},
			expectErr: true,
		},
		"charset restrictions": {
			registry: defaultRuleNameMapping,
			rawConfig: `
				length = 20
				rule "charset" {
					charset = "abcde"
					min-chars = 2
				}`,
			expected: StringGenerator{
				Length:  20,
				charset: []rune("abcde"),
				Rules: []Rule{
					CharsetRule{
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
				rule "testrule" {
					string = "teststring"
					int = 123
				}`,
			expected: StringGenerator{
				Length:  20,
				charset: deduplicateRunes([]rune("teststring")),
				Rules: []Rule{
					testCharsetRule{
						String:  "teststring",
						Integer: 123,
					},
				},
			},
			expectErr: false,
		},
		"test rule and charset restrictions": {
			registry: map[string]ruleConstructor{
				"testrule": newTestRule,
				"charset":  ParseCharset,
			},
			rawConfig: `
				length = 20
				rule "testrule" {
					string = "teststring"
					int = 123
				}
				rule "charset" {
					charset = "abcde"
					min-chars = 2
				}`,
			expected: StringGenerator{
				Length:  20,
				charset: deduplicateRunes([]rune("abcdeteststring")),
				Rules: []Rule{
					testCharsetRule{
						String:  "teststring",
						Integer: 123,
					},
					CharsetRule{
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
				rule "testrule" {
					string = "teststring"
					int = 123
				}`,
			expected:  StringGenerator{},
			expectErr: true,
		},

		// /////////////////////////////////////////////////
		// JSON data
		"manually JSONified HCL": {
			registry: map[string]ruleConstructor{
				"testrule": newTestRule,
				"charset":  ParseCharset,
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
							"charset": [
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
				charset: deduplicateRunes([]rune("abcdeteststring")),
				Rules: []Rule{
					testCharsetRule{
						String:  "teststring",
						Integer: 123,
					},
					CharsetRule{
						Charset:  []rune("abcde"),
						MinChars: 2,
					},
				},
			},
			expectErr: false,
		},
		"JSONified HCL": {
			registry: map[string]ruleConstructor{
				"testrule": newTestRule,
				"charset":  ParseCharset,
			},
			rawConfig: toJSON(t, StringGenerator{
				Length: 20,
				Rules: []Rule{
					testCharsetRule{
						String:  "teststring",
						Integer: 123,
					},
					CharsetRule{
						Charset:  []rune("abcde"),
						MinChars: 2,
					},
				},
			}),
			expected: StringGenerator{
				Length:  20,
				charset: deduplicateRunes([]rune("abcdeteststring")),
				Rules: []Rule{
					testCharsetRule{
						String:  "teststring",
						Integer: 123,
					},
					CharsetRule{
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
			expected:  StringGenerator{},
			expectErr: true,
		},
		"config value with empty slice": {
			registry: defaultRuleNameMapping,
			rawConfig: `
                rule {
                    n = []
                }`,
			expected:  StringGenerator{},
			expectErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			parser := PolicyParser{
				RuleRegistry: Registry{
					Rules: test.registry,
				},
			}

			actual, err := parser.ParsePolicy(test.rawConfig)
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
				testCharsetRule{
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

func BenchmarkParser_Parse(b *testing.B) {
	config := `length = 20
               rule "charset" {
                   charset = "abcde"
                   min-chars = 2
               }`

	for i := 0; i < b.N; i++ {
		parser := PolicyParser{
			RuleRegistry: Registry{
				Rules: defaultRuleNameMapping,
			},
		}
		_, err := parser.ParsePolicy(config)
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
