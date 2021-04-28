package random

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// serializableRules is a slice of rules that can be marshalled to JSON in an HCL format
type serializableRules []Rule

// MarshalJSON in an HCL-friendly way
func (r serializableRules) MarshalJSON() (b []byte, err error) {
	// Example:
	// [
	//   {
	//     "testrule": [
	//       {
	//         "string": "teststring",
	//         "int": 123
	//       }
	//     ]
	//   },
	//   {
	//     "charset": [
	//       {
	//         "charset": "abcde",
	//         "min-chars": 2
	//       }
	//     ]
	//   }
	// ]
	data := []map[string][]map[string]interface{}{} // Totally not confusing at all
	for _, rule := range r {
		ruleData := map[string]interface{}{}
		err = mapstructure.Decode(rule, &ruleData)
		if err != nil {
			return nil, fmt.Errorf("unable to decode rule: %w", err)
		}

		ruleMap := map[string][]map[string]interface{}{
			rule.Type(): {
				ruleData,
			},
		}
		data = append(data, ruleMap)
	}

	b, err = json.Marshal(data)
	return b, err
}

func (r *serializableRules) UnmarshalJSON(data []byte) (err error) {
	mapData := []map[string]interface{}{}
	err = json.Unmarshal(data, &mapData)
	if err != nil {
		return err
	}
	rules, err := parseRules(defaultRegistry, mapData)
	if err != nil {
		return err
	}
	*r = rules
	return nil
}

type runes []rune

func (r runes) Len() int           { return len(r) }
func (r runes) Less(i, j int) bool { return r[i] < r[j] }
func (r runes) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// MarshalJSON converts the runes to a string for smaller JSON and easier readability
func (r runes) MarshalJSON() (b []byte, err error) {
	return json.Marshal(string(r))
}

// UnmarshalJSON converts a string to []rune
func (r *runes) UnmarshalJSON(data []byte) (err error) {
	var str string
	err = json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	*r = []rune(str)
	return nil
}
