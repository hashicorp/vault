// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package random

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type testCharsetRule struct {
	String  string `mapstructure:"string" json:"string"`
	Integer int    `mapstructure:"int"    json:"int"`

	// Default to passing
	fail bool
}

func newTestRule(data map[string]interface{}) (rule Rule, err error) {
	tr := &testCharsetRule{}
	err = mapstructure.Decode(data, tr)
	if err != nil {
		return nil, fmt.Errorf("unable to decode test rule")
	}
	return *tr, nil
}

func (tr testCharsetRule) Pass([]rune) bool { return !tr.fail }
func (tr testCharsetRule) Type() string     { return "testrule" }
func (tr testCharsetRule) Chars() []rune    { return []rune(tr.String) }

func TestParseRule(t *testing.T) {
	type testCase struct {
		rules map[string]ruleConstructor

		ruleType string
		ruleData map[string]interface{}

		expectedRule Rule
		expectErr    bool
	}

	tests := map[string]testCase{
		"missing rule": {
			rules:    map[string]ruleConstructor{},
			ruleType: "testrule",
			ruleData: map[string]interface{}{
				"string": "teststring",
				"int":    123,
			},
			expectedRule: nil,
			expectErr:    true,
		},
		"nil data": {
			rules: map[string]ruleConstructor{
				"testrule": newTestRule,
			},
			ruleType:     "testrule",
			ruleData:     nil,
			expectedRule: testCharsetRule{},
			expectErr:    false,
		},
		"good rule": {
			rules: map[string]ruleConstructor{
				"testrule": newTestRule,
			},
			ruleType: "testrule",
			ruleData: map[string]interface{}{
				"string": "teststring",
				"int":    123,
			},
			expectedRule: testCharsetRule{
				String:  "teststring",
				Integer: 123,
			},
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reg := Registry{
				Rules: test.rules,
			}

			actualRule, err := reg.parseRule(test.ruleType, test.ruleData)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}

			if !reflect.DeepEqual(actualRule, test.expectedRule) {
				t.Fatalf("Actual: %#v\nExpected:%#v", actualRule, test.expectedRule)
			}
		})
	}
}

// Ensure the mappings in the defaultRuleNameMapping are consistent between the keys
// in the map and the Type() calls on the Rule values
func TestDefaultRuleNameMapping(t *testing.T) {
	for expectedType, constructor := range defaultRuleNameMapping {
		// In this case, we don't care about the error since we're checking the types, not the contents
		instance, _ := constructor(map[string]interface{}{})
		actualType := instance.Type()
		if actualType != expectedType {
			t.Fatalf("Default registry mismatched types: Actual: %s Expected: %s", actualType, expectedType)
		}
	}
}
