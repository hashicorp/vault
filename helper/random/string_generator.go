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
	DefaultStringGenerator = &StringGenerator{
		Length: 20,
		Rules: []Rule{
			CharsetRule{
				Charset:  LowercaseRuneset,
				MinChars: 1,
			},
			CharsetRule{
				Charset:  UppercaseRuneset,
				MinChars: 1,
			},
			CharsetRule{
				Charset:  NumericRuneset,
				MinChars: 1,
			},
			CharsetRule{
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

// StringGenerator generats random strings from the provided charset & adhering to a set of rules. The set of rules
// are things like CharsetRule which requires a certain number of characters from a sub-charset.
type StringGenerator struct {
	// Length of the string to generate.
	Length int `mapstructure:"length" json:"length"`

	// Rules the generated strings must adhere to.
	Rules serializableRules `mapstructure:"-" json:"rule"` // This is "rule" in JSON so it matches the HCL property type

	// CharsetRule to choose runes from. This is computed from the rules, not directly configurable
	charset runes
}

// Generate a random string from the charset and adhering to the provided rules.
// The io.Reader is optional. If not provided, it will default to the reader from crypto/rand
func (g *StringGenerator) Generate(ctx context.Context, rng io.Reader) (str string, err error) {
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
			str, err = g.generate(rng)
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

func (g *StringGenerator) generate(rng io.Reader) (str string, err error) {
	// If performance improvements need to be made, this can be changed to read a batch of
	// potential strings at once rather than one at a time. This will significantly
	// improve performance, but at the cost of added complexity.
	candidate, err := randomRunes(rng, g.charset, g.Length)
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

const (
	// maxCharsetLen is the maximum length a charset is allowed to be when generating a candidate string.
	// This is the total number of numbers available for selecting an index out of the charset slice.
	maxCharsetLen = 256
)

// randomRunes creates a random string based on the provided charset. The charset is limited to 255 characters, but
// could be expanded if needed. Expanding the maximum charset size will decrease performance because it will need to
// combine bytes into a larger integer using binary.BigEndian.Uint16() function.
func randomRunes(rng io.Reader, charset []rune, length int) (candidate []rune, err error) {
	if len(charset) == 0 {
		return nil, fmt.Errorf("no charset specified")
	}
	if len(charset) > maxCharsetLen {
		return nil, fmt.Errorf("charset is too long: limited to %d characters", math.MaxUint8)
	}
	if length <= 0 {
		return nil, fmt.Errorf("unable to generate a zero or negative length runeset")
	}

	// This can't always select indexes from [0-maxCharsetLen) because it could introduce bias to the character selection.
	// For instance, if the length of the charset is [a-zA-Z0-9-] (length of 63):
	// RNG ranges: [0-62][63-125][126-188][189-251] will equally select from the entirety of the charset. However,
	// the RNG values [252-255] will select the first 4 characters of the charset while ignoring the remaining 59.
	// This results in a bias towards the front of the charset.
	//
	// To avoid this, we determine the largest integer multiplier of the charset length that is <= maxCharsetLen
	// For instance, if the maxCharsetLen is 256 (the size of one byte) and the charset is length 63, the multiplier
	// equals 4:
	//   256/63 => 4.06
	//   Trunc(4.06) => 4
	// Multiply by the charset length
	// Subtract 1 to account for 0-based counting and you get the max index value: 251
	maxAllowedRNGValue := (maxCharsetLen/len(charset))*len(charset) - 1

	// rngBufferMultiplier increases the size of the RNG buffer to account for lost
	// indexes due to the maxAllowedRNGValue
	rngBufferMultiplier := 1.0

	// Don't set a multiplier if we are able to use the entire range of indexes
	if maxAllowedRNGValue < maxCharsetLen {
		// Anything more complicated than an arbitrary percentage appears to have little practical performance benefit
		rngBufferMultiplier = 1.5
	}

	// Default to the standard crypto reader if one isn't provided
	if rng == nil {
		rng = rand.Reader
	}

	charsetLen := byte(len(charset))

	runes := make([]rune, 0, length)

	for len(runes) < length {
		// Generate a bunch of indexes
		data := make([]byte, int(float64(length)*rngBufferMultiplier))
		numBytes, err := rng.Read(data)
		if err != nil {
			return nil, err
		}

		// Append characters until either we're out of indexes or the length is long enough
		for i := 0; i < numBytes; i++ {
			// Be careful to ensure that maxAllowedRNGValue isn't >= 256 as it will overflow and this
			// comparison will prevent characters from being selected from the charset
			if data[i] > byte(maxAllowedRNGValue) {
				continue
			}

			index := data[i]
			if len(charset) != maxCharsetLen {
				index = index % charsetLen
			}
			r := charset[index]
			runes = append(runes, r)

			if len(runes) == length {
				break
			}
		}
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
