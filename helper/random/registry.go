package random

import (
	"fmt"
)

type ruleConstructor func(map[string]interface{}) (Rule, error)

var (
	// defaultRuleNameMapping is the default mapping of HCL rule names to the appropriate rule constructor.
	// Add to this map when adding a new Rule type to be recognized in HCL.
	defaultRuleNameMapping = map[string]ruleConstructor{
		"charset": ParseCharset,
	}

	defaultRegistry = Registry{
		Rules: defaultRuleNameMapping,
	}
)

// Registry of HCL rule names to rule constructors.
type Registry struct {
	// Rules maps names of rules to a constructor for the rule
	Rules map[string]ruleConstructor
}

func (r Registry) parseRule(ruleType string, ruleData map[string]interface{}) (rule Rule, err error) {
	constructor, exists := r.Rules[ruleType]
	if !exists {
		return nil, fmt.Errorf("unrecognized rule type %s", ruleType)
	}

	rule, err = constructor(ruleData)
	return rule, err
}
