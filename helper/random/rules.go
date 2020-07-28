package random

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Rule to assert on string values.
type Rule interface {
	// Pass should return true if the provided value passes any assertions this Rule is making.
	Pass(value []rune) bool

	// Type returns the name of the rule as associated in the registry
	Type() string
}

// CharsetRule requires a certain number of characters from the specified charset.
type CharsetRule struct {
	// CharsetRule is the list of rules that candidate strings must contain a minimum number of.
	Charset runes `mapstructure:"charset" json:"charset"`

	// MinChars indicates the minimum (inclusive) number of characters from the charset that should appear in the string.
	MinChars int `mapstructure:"min-chars" json:"min-chars"`
}

// ParseCharset from the provided data map. The data map is expected to be parsed from HCL.
func ParseCharset(data map[string]interface{}) (rule Rule, err error) {
	cr := &CharsetRule{}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		Result:     cr,
		DecodeHook: stringToRunesFunc,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to decode charset restriction: %w", err)
	}

	err = decoder.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse charset restriction: %w", err)
	}

	return *cr, nil
}

func (c CharsetRule) Type() string {
	return "charset"
}

// Chars returns the charset that this rule is looking for.
func (c CharsetRule) Chars() []rune {
	return c.Charset
}

func (c CharsetRule) MinLength() int {
	return c.MinChars
}

// Pass returns true if the provided candidate string has a minimum number of chars in it.
// This adheres to the Rule interface
func (c CharsetRule) Pass(value []rune) bool {
	if c.MinChars <= 0 {
		return true
	}

	count := 0
	for _, r := range value {
		// charIn is sometimes faster than a map lookup because the data is so small
		// This is being kept rather than converted to a map to keep the code cleaner,
		// otherwise there would need to be additional parsing logic.
		if charIn(r, c.Charset) {
			count++
			if count >= c.MinChars {
				return true
			}
		}
	}

	return false
}

func charIn(search rune, charset []rune) bool {
	for _, r := range charset {
		if search == r {
			return true
		}
	}
	return false
}
