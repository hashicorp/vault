package random

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"sort"
	"time"
	"unicode"

	"github.com/hashicorp/go-multierror"
)

var (
	LowercaseCharset   = sortCharset("abcdefghijklmnopqrstuvwxyz")
	UppercaseCharset   = sortCharset("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	NumericCharset     = sortCharset("0123456789")
	FullSymbolCharset  = sortCharset("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~")
	ShortSymbolCharset = sortCharset("-")

	AlphabeticCharset              = sortCharset(UppercaseCharset + LowercaseCharset)
	AlphaNumericCharset            = sortCharset(AlphabeticCharset + NumericCharset)
	AlphaNumericShortSymbolCharset = sortCharset(AlphaNumericCharset + ShortSymbolCharset)
	AlphaNumericFullSymbolCharset  = sortCharset(AlphaNumericCharset + FullSymbolCharset)

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
		Length: 20,
		Rules: []Rule{
			Charset{
				Charset:  LowercaseRuneset,
				MinChars: 1,
			},
			Charset{
				Charset:  UppercaseRuneset,
				MinChars: 1,
			},
			Charset{
				Charset:  NumericRuneset,
				MinChars: 1,
			},
			Charset{
				Charset:  ShortSymbolRuneset,
				MinChars: 1,
			},
		},
	}
)

func sortCharset(chars string) string {
	r := runes(chars)
	sort.Sort(r)
	return string(r)
}

// Rule to assert on string values.
type Rule interface {
	// Pass should return true if the provided value passes any assertions this Rule is making.
	Pass(value []rune) bool

	// Type returns the name of the rule as associated in the registry
	Type() string
}

// StringGenerator generats random strings from the provided charset & adhering to a set of rules. The set of rules
// are things like Charset which requires a certain number of characters from a sub-charset.
type StringGenerator struct {
	// Length of the string to generate.
	Length int `mapstructure:"length" json:"length"`

	// Rules the generated strings must adhere to.
	Rules serializableRules `mapstructure:"-" json:"rule"` // This is "rule" in JSON so it matches the HCL property type

	// Charset to choose runes from. This is computed from the rules, not directly configurable
	charset runes

	// rng for testing purposes to ensure error handling from the crypto/rand package is working properly.
	rng io.Reader
}

// Generate a random string from the charset and adhering to the provided rules.
func (g *StringGenerator) Generate(ctx context.Context) (str string, err error) {
	if _, hasTimeout := ctx.Deadline(); !hasTimeout {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, 1*time.Second) // Ensure there's a timeout on the context
		defer cancel()
	}

	// Ensure the generator is configured well since it may be manually created rather than parsed from HCL
	err = g.validateConfig()
	if err != nil {
		return "", err
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

func (g *StringGenerator) generate() (str string, err error) {
	// If performance improvements need to be made, this can be changed to read a batch of
	// potential strings at once rather than one at a time. This will significantly
	// improve performance, but at the cost of added complexity.
	candidate, err := randomRunes(g.rng, g.charset, g.Length)
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

// validateConfig of the generator to ensure that we can successfully generate a string.
func (g *StringGenerator) validateConfig() (err error) {
	merr := &multierror.Error{}

	// Ensure the sum of minimum lengths in the rules doesn't exceed the length specified
	minLen := getMinLength(g.Rules)
	if g.Length <= 0 {
		merr = multierror.Append(merr, fmt.Errorf("length must be > 0"))
	} else if g.Length < minLen {
		merr = multierror.Append(merr, fmt.Errorf("specified rules require at least %d characters but %d is specified", minLen, g.Length))
	}

	// Ensure we have a charset & all characters are printable
	if len(g.charset) == 0 {
		// Yes this is mutating the generator but this is done so we don't have to compute this on every generation
		g.charset = getChars(g.Rules)
	}
	if len(g.charset) == 0 {
		merr = multierror.Append(merr, fmt.Errorf("no charset specified"))
	} else {
		for _, r := range g.charset {
			if !unicode.IsPrint(r) {
				merr = multierror.Append(merr, fmt.Errorf("non-printable character in charset"))
				break
			}
		}
	}
	return merr.ErrorOrNil()
}

// getMinLength from the rules using the optional interface: `MinLength() int`
func getMinLength(rules []Rule) (minLen int) {
	type minLengthProvider interface {
		MinLength() int
	}

	for _, rule := range rules {
		mlp, ok := rule.(minLengthProvider)
		if !ok {
			continue
		}
		minLen += mlp.MinLength()
	}
	return minLen
}

// getChars from the rules using the optional interface: `Chars() []rune`
func getChars(rules []Rule) (chars []rune) {
	type charsetProvider interface {
		Chars() []rune
	}

	for _, rule := range rules {
		cp, ok := rule.(charsetProvider)
		if !ok {
			continue
		}
		chars = append(chars, cp.Chars()...)
	}
	return deduplicateRunes(chars)
}

// deduplicateRunes returns a new slice of sorted & de-duplicated runes
func deduplicateRunes(original []rune) (deduped []rune) {
	if len(original) == 0 {
		return nil
	}

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
