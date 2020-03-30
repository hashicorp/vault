package random

import (
	"fmt"
	"reflect"
	"sort"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/mitchellh/mapstructure"
)

func Parse(raw string) (strs StringGenerator, err error) {
	parser := Parser{
		RuleRegistry: Registry{
			Rules: defaultRuleNameMapping,
		},
	}
	return parser.Parse(raw)
}

type Parser struct {
	RuleRegistry Registry
}

func (p Parser) Parse(raw string) (strs StringGenerator, err error) {
	rawData := map[string]interface{}{}
	err = hcl.Decode(&rawData, raw)
	if err != nil {
		return strs, fmt.Errorf("unable to decode: %w", err)
	}

	// Decode the top level items
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     &strs,
		DecodeHook: stringToRunesFunc,
	})
	if err != nil {
		return strs, fmt.Errorf("unable to decode configuration: %w", err)
	}

	err = decoder.Decode(rawData)
	if err != nil {
		return strs, fmt.Errorf("failed to decode configuration: %w", err)
	}

	// Decode & parse rules
	rawRules, err := getMapSlice(rawData, "rule")
	if err != nil {
		return strs, fmt.Errorf("unable to retrieve rules: %w", err)
	}

	rules, err := p.parseRules(rawRules)
	if err != nil {
		return strs, fmt.Errorf("unable to parse rules: %w", err)
	}

	// Add any charsets found in rules to the overall charset & deduplicate
	cs := append(strs.Charset, getChars(rules)...)
	cs = deduplicateRunes(cs)

	strs = StringGenerator{
		Length:  strs.Length,
		Charset: cs,
		Rules:   rules,
	}

	err = validate(strs)
	if err != nil {
		return strs, err
	}

	return strs, nil
}

func (p Parser) parseRules(rawRules []map[string]interface{}) (rules []Rule, err error) {
	for _, rawRule := range rawRules {
		info, err := getRuleInfo(rawRule)
		if err != nil {
			return nil, fmt.Errorf("unable to get rule info: %w", err)
		}

		// Map names like "lower-alpha" to lowercase alphabetical characters
		applyShortcuts(info.data)

		rule, err := p.RuleRegistry.parseRule(info.ruleType, info.data)
		if err != nil {
			return nil, fmt.Errorf("unable to parse rule %s: %w", info.ruleType, err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func validate(strs StringGenerator) (err error) {
	merr := &multierror.Error{}
	if strs.Length < 1 {
		merr = multierror.Append(merr, fmt.Errorf("length must be >= 1"))
	}
	if len(strs.Charset) == 0 {
		merr = multierror.Append(merr, fmt.Errorf("no charset specified"))
	}

	for _, r := range strs.Charset {
		if !unicode.IsPrint(r) {
			merr = multierror.Append(merr, fmt.Errorf("non-printable character in charset"))
			break
		}
	}

	return merr.ErrorOrNil()
}

func getChars(rules []Rule) (chars []rune) {
	type charsetProvider interface {
		Chars() []rune
	}

	for _, rule := range rules {
		cv, ok := rule.(charsetProvider)
		if !ok {
			continue
		}
		chars = append(chars, cv.Chars()...)
	}
	return chars
}

// getMapSlice from the provided map. This will retrieve and type-assert a []map[string]interface{} from the map
// This will not error if the key does not exist
// This will return an error if the value at the provided key is not of type []map[string]interface{}
func getMapSlice(m map[string]interface{}, key string) (mapSlice []map[string]interface{}, err error) {
	rawSlice, exists := m[key]
	if !exists {
		return nil, nil
	}

	slice, ok := rawSlice.([]map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("key %s is not a []map[string]interface{}", key)
	}

	return slice, nil
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
		data = ruleInfo{
			ruleType: key,
			data:     slice[0],
		}
		return data, nil
	}
	return data, fmt.Errorf("rule is empty")
}

var (
	charsetShortcuts = map[string]string{
		// Base
		"lower-alpha": "abcdefghijklmnopqrstuvwxyz",
		"upper-alpha": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"numeric":     "0123456789",

		// Combinations
		"lower-upper-alpha":        "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"lower-upper-alphanumeric": "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	}
)

// applyShortcuts to the provided map. This will look for a "charset" key. If it exists and equals one of the keys
// in `charsetShortcuts`, it replaces the value with the value found in the `charsetShortcuts` map. For instance:
//
// Input map:
// map[string]interface{}{
//   "charset": "upper-alpha",
// }
//
// This will convert it to:
// map[string]interface{}{
//   "charset": "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
// }
func applyShortcuts(m map[string]interface{}) {
	rawCharset, exists := m["charset"]
	if !exists {
		return
	}
	charset, ok := rawCharset.(string)
	if !ok {
		return
	}
	newCharset, shortcutExists := charsetShortcuts[charset]
	if !shortcutExists {
		return
	}
	m["charset"] = newCharset
}

func deduplicateRunes(original []rune) (deduped []rune) {
	m := map[rune]bool{}
	dedupedRunes := []rune(nil)

	for _, r := range original {
		if m[r] {
			continue
		}
		m[r] = true
		dedupedRunes = append(dedupedRunes, r)
	}

	// They don't have to be sorted, but this is being done to make the charset easier to visualize
	sort.Sort(runes(dedupedRunes))
	return dedupedRunes
}

type runes []rune

func (r runes) Len() int           { return len(r) }
func (r runes) Less(i, j int) bool { return r[i] < r[j] }
func (r runes) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

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
