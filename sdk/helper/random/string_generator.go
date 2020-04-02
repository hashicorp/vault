package random

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"time"
)

const (
	LowercaseCharset   = "abcdefghijklmnopqrstuvwxyz"
	UppercaseCharset   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumericCharset     = "0123456789"
	FullSymbolCharset  = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	ShortSymbolCharset = "-"

	AlphabeticCharset              = UppercaseCharset + LowercaseCharset
	AlphaNumericCharset            = AlphabeticCharset + NumericCharset
	AlphaNumericShortSymbolCharset = AlphaNumericCharset + ShortSymbolCharset
	AlphaNumericFullSymbolCharset  = AlphaNumericCharset + FullSymbolCharset
)

var (
	LowercaseRuneset   = []rune(LowercaseCharset)
	UppercaseRuneset   = []rune(UppercaseCharset)
	NumericRuneset     = []rune(NumericCharset)
	FullSymbolRuneset  = []rune(FullSymbolCharset)
	ShortSymbolRuneset = []rune(ShortSymbolCharset)

	AlphabeticRuneset              = []rune(AlphabeticCharset)
	AlphaNumericRuneset            = []rune(AlphaNumericCharset)
	AlphaNumericShortSymbolRuneset = []rune(AlphaNumericShortSymbolCharset)
	AlphaNumericFullSymbolRuneset  = []rune(AlphaNumericFullSymbolCharset)

	// DefaultStringGenerator has reasonable default rules for generating strings
	DefaultStringGenerator = StringGenerator{
		Length:  20,
		Charset: []rune(LowercaseCharset + UppercaseCharset + NumericCharset + ShortSymbolCharset),
		Rules: []Rule{
			CharsetRestriction{
				Charset:  LowercaseRuneset,
				MinChars: 1,
			},
			CharsetRestriction{
				Charset:  UppercaseRuneset,
				MinChars: 1,
			},
			CharsetRestriction{
				Charset:  NumericRuneset,
				MinChars: 1,
			},
			CharsetRestriction{
				Charset:  ShortSymbolRuneset,
				MinChars: 1,
			},
		},
	}
)

// Rule to assert on string values.
type Rule interface {
	// Pass should return true if the provided value passes any assertions this Rule is making.
	Pass(value []rune) bool
}

// StringGenerator generats random strings from the provided charset & adhering to a set of rules. The set of rules
// are things like CharsetRestriction which requires a certain number of characters from a sub-charset.
type StringGenerator struct {
	// Length of the string to generate.
	Length int `mapstructure:"length"`

	// Charset to choose runes from.
	Charset []rune `mapstructure:"charset"`

	// Rules the generated strings must adhere to.
	Rules []Rule `mapstructure:"-"`

	// rng for testing purposes to ensure error handling from the crypto/rand package is working properly.
	rng io.Reader
}

// Generate a random string from the charset and adhering to the provided rules.
func (g StringGenerator) Generate(ctx context.Context) (str string, err error) {
	if _, hasTimeout := ctx.Deadline(); !hasTimeout {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 1*time.Second) // Ensure there's a timeout on the context
		defer cancel()
	}

LOOP:
	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timed out generating string")
		default:
			str, err = g.generate()
			if err != nil {
				return "", err
			}
			if str == "" {
				continue LOOP
			}
			return str, err
		}
	}
}

func (g StringGenerator) generate() (str string, err error) {
	// If performance improvements need to be made, this can be changed to read a batch of
	// potential strings at once rather than one at a time. This will significantly
	// improve performance, but at the cost of added complexity.
	candidate, err := randomRunes(g.rng, g.Charset, g.Length)
	if err != nil {
		return "", fmt.Errorf("unable to generate random characters: %w", err)
	}

	for _, rule := range g.Rules {
		if !rule.Pass(candidate) {
			return "", nil
		}
	}

	// Passed all rules
	return string(candidate), nil
}

// randomRunes creates a random string based on the provided charset. The charset is limited to 255 characters, but
// could be expanded if needed. Expanding the maximum charset size will decrease performance because it will need to
// combine bytes into a larger integer using binary.BigEndian.Uint16() function.
func randomRunes(rng io.Reader, charset []rune, length int) (candidate []rune, err error) {
	if len(charset) == 0 {
		return nil, fmt.Errorf("no charset specified")
	}
	if len(charset) > math.MaxUint8 {
		return nil, fmt.Errorf("charset is too long: limited to %d characters", math.MaxUint8)
	}
	if length <= 0 {
		return nil, fmt.Errorf("unable to generate a zero or negative length runeset")
	}

	if rng == nil {
		rng = rand.Reader
	}

	charsetLen := byte(len(charset))
	data := make([]byte, length)
	_, err = rng.Read(data)
	if err != nil {
		return nil, err
	}

	runes := make([]rune, 0, length)
	for i := 0; i < len(data); i++ {
		r := charset[data[i]%charsetLen]
		runes = append(runes, r)
	}

	return runes, nil
}
