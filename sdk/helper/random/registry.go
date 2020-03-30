package random

import (
	"fmt"
)

type ruleConstructor func(map[string]interface{}) (Rule, error)

var (
	defaultRegistry = Registry{
		Rules: defaultRuleNameMapping,
	}

	defaultRuleNameMapping = map[string]ruleConstructor{
		"CharsetRestriction": NewCharsetRestriction,
	}
)

type Registry struct {
	Rules map[string]ruleConstructor
}

func (r Registry) parseRule(ruleType string, ruleData map[string]interface{}) (rule Rule, err error) {
	constructor, exists := r.Rules[ruleType]
	if !exists {
		return nil, fmt.Errorf("unrecognized rule type %s", ruleType)
	}

	return constructor(ruleData)
}
