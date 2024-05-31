// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package random

import (
	"fmt"
	"reflect"
	"unicode/utf8"

	"github.com/hashicorp/hcl"
	"github.com/mitchellh/mapstructure"
)

// ParsePolicy is a convenience function for parsing HCL into a StringGenerator.
// See PolicyParser.ParsePolicy for details.
func ParsePolicy(raw string) (gen StringGenerator, err error) {
	parser := PolicyParser{
		RuleRegistry: Registry{
			Rules: defaultRuleNameMapping,
		},
	}
	return parser.ParsePolicy(raw)
}

// ParsePolicyBytes is a convenience function for parsing HCL into a StringGenerator.
// See PolicyParser.ParsePolicy for details.
func ParsePolicyBytes(raw []byte) (gen StringGenerator, err error) {
	return ParsePolicy(string(raw))
}

// PolicyParser parses string generator configuration from HCL.
type PolicyParser struct {
	// RuleRegistry maps rule names in HCL to Rule constructors.
	RuleRegistry Registry
}

// ParsePolicy parses the provided HCL into a StringGenerator.
func (p PolicyParser) ParsePolicy(raw string) (sg StringGenerator, err error) {
	rawData := map[string]interface{}{}
	err = hcl.Decode(&rawData, raw)
	if err != nil {
		return sg, fmt.Errorf("unable to decode: %w", err)
	}

	// Decode the top level items
	gen := StringGenerator{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     &gen,
		DecodeHook: stringToRunesFunc,
	})
	if err != nil {
		return sg, fmt.Errorf("unable to decode configuration: %w", err)
	}

	err = decoder.Decode(rawData)
	if err != nil {
		return sg, fmt.Errorf("failed to decode configuration: %w", err)
	}

	// Decode & parse rules
	rawRules, err := getMapSlice(rawData, "rule")
	if err != nil {
		return sg, fmt.Errorf("unable to retrieve rules: %w", err)
	}

	rules, err := parseRules(p.RuleRegistry, rawRules)
	if err != nil {
		return sg, fmt.Errorf("unable to parse rules: %w", err)
	}

	gen = StringGenerator{
		Length: gen.Length,
		Rules:  rules,
	}

	err = gen.validateConfig()
	if err != nil {
		return sg, err
	}

	return gen, nil
}

func parseRules(registry Registry, rawRules []map[string]interface{}) (rules []Rule, err error) {
	for _, rawRule := range rawRules {
		info, err := getRuleInfo(rawRule)
		if err != nil {
			return nil, fmt.Errorf("unable to get rule info: %w", err)
		}

		rule, err := registry.parseRule(info.ruleType, info.data)
		if err != nil {
			return nil, fmt.Errorf("unable to parse rule %s: %w", info.ruleType, err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// getMapSlice from the provided map. This will retrieve and type-assert a []map[string]interface{} from the map
// This will not error if the key does not exist
// This will return an error if the value at the provided key is not of type []map[string]interface{}
func getMapSlice(m map[string]interface{}, key string) (mapSlice []map[string]interface{}, err error) {
	rawSlice, exists := m[key]
	if !exists {
		return nil, nil
	}

	mapSlice = []map[string]interface{}{}
	err = mapstructure.Decode(rawSlice, &mapSlice)
	if err != nil {
		return nil, err
	}
	return mapSlice, nil
}

type ruleInfo struct {
	ruleType string
	data     map[string]interface{}
}

// getRuleInfo splits the provided HCL-decoded rule into its rule type along with the data associated with it
func getRuleInfo(rule map[string]interface{}) (data ruleInfo, err error) {
	// There should only be one key, but it's a dynamic key yay!
	for key := range rule {
		slice, err := getMapSlice(rule, key)
		if err != nil {
			return data, fmt.Errorf("unable to get rule data: %w", err)
		}

		if len(slice) == 0 {
			return data, fmt.Errorf("rule info cannot be empty")
		}

		data = ruleInfo{
			ruleType: key,
			data:     slice[0],
		}
		return data, nil
	}
	return data, fmt.Errorf("rule is empty")
}

// stringToRunesFunc converts a string to a []rune for use in the mapstructure library
func stringToRunesFunc(from reflect.Kind, to reflect.Kind, data interface{}) (interface{}, error) {
	if from != reflect.String || to != reflect.Slice {
		return data, nil
	}

	raw := data.(string)

	if !utf8.ValidString(raw) {
		return nil, fmt.Errorf("invalid UTF8 string")
	}
	return []rune(raw), nil
}
