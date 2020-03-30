package random

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
)

type testRule struct {
	String  string `mapstructure:"string" json:"string"`
	Integer int    `mapstructure:"int"    json:"int"`

	// Default to passing
	fail bool
}

func newTestRule(data map[string]interface{}) (rule Rule, err error) {
	tr := &testRule{}
	err = mapstructure.Decode(data, tr)
	if err != nil {
		return nil, fmt.Errorf("unable to decode test rule")
	}
	return tr, nil
}

func (tr testRule) Pass([]rune) bool { return !tr.fail }

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
				"string": "omgwtfbbq",
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
			expectedRule: &testRule{},
			expectErr:    false,
		},
		"good rule": {
			rules: map[string]ruleConstructor{
				"testrule": newTestRule,
			},
			ruleType: "testrule",
			ruleData: map[string]interface{}{
				"string": "omgwtfbbq",
				"int":    123,
			},
			expectedRule: &testRule{
				String:  "omgwtfbbq",
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
