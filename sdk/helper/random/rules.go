package random

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// CharsetRestriction requires a certain number of characters from the specified charset
type CharsetRestriction struct {
	Charset []rune `mapstructure:"charset" json:"charset"`

	// MinChars indicates the minimum (inclusive) number of characters from the charset that should appear in the string
	MinChars int `mapstructure:"min-chars" json:"min-chars"`
}

func NewCharsetRestriction(data map[string]interface{}) (rule Rule, err error) {
	cr := &CharsetRestriction{}

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
	return cr, nil
}

func (c CharsetRestriction) Chars() []rune {
	return c.Charset
}

func (c CharsetRestriction) Pass(value []rune) bool {
	if c.MinChars <= 0 {
		return true
	}

	count := 0
	for _, r := range value {
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
