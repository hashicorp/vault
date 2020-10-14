package main

import (
	"crypto/rand"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/vault/helper/random"

	"github.com/hashicorp/go-multierror"
)

const (
	DefaultCharset = LowerCharset + UpperCharset + NumericCharset
	UpperCharset   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerCharset   = "abcdefghijklmnopqrstuvwxyz"
	NumericCharset = "0123456789"
	SymbolsCharset = "!\"#$%&'()*+,-./:;<=>?`\\~[]^_@|" // no whitespace characters for safety
)

type UsernameOpt func(*UsernameProducer) error

func Template(rawTemplate string) UsernameOpt {
	return func(up *UsernameProducer) error {
		up.rawTemplate = rawTemplate
		return nil
	}
}

// SetFunction allows the user to specify functions for use in the template. This can be used to override the
// behavior of any of the built-in functions (such as rand). Overriding the built-in functions is not recommended.
func SetFunction(name string, f interface{}) UsernameOpt {
	return func(up *UsernameProducer) error {
		if name == "" {
			return fmt.Errorf("missing function name")
		}
		if f == nil {
			return fmt.Errorf("missing function")
		}
		up.funcMap[name] = f
		return nil
	}
}

// MaxLength is the maximum allowed length of the username. This will take effect after applying the template.
// A length of 0 indicates no maximum length
func MaxLength(maxLen int) UsernameOpt {
	return func(up *UsernameProducer) error {
		if maxLen < 0 {
			return fmt.Errorf("max username length must be >= 0")
		}
		up.maxLen = maxLen
		return nil
	}
}

// RandomCharset specifies the charset used when using the `rand` template function.
// This will remove duplicate characters from the charset to prevent bias.
func RandomCharset(charset string) UsernameOpt {
	if charset == "" {
		charset = DefaultCharset
	}
	dedupedRunes := deduplicateRunes([]rune(charset))
	return SetFunction("rand", randCharset([]rune(dedupedRunes)))
}

func ToUppercase() UsernameOpt {
	return func(up *UsernameProducer) error {
		up.casing = upperCase
		return nil
	}
}

func ToLowercase() UsernameOpt {
	return func(up *UsernameProducer) error {
		up.casing = lowerCase
		return nil
	}
}

type casing uint8

const (
	ignoreCase casing = iota
	upperCase
	lowerCase
)

// UsernameProducer creates usernames based on the provided template.
// This uses the go templating language, so anything that adheres to that language will function in this struct.
// There are several custom functions available for use in the template:
// - rand
//   - Randomly generated characters. This uses the charset specified in RandomCharset. Must include a length.
//     Example: {{rand 20}}
// - truncate
//   - Truncates the previous value to the specified length. Must include a maximum length.
//     Example: {{.DisplayName | truncate 10}}
// - now_seconds
//   - Provides the current unix time in seconds.
//     Example: {{now_seconds}}
// - now_nano
//   - Provides the current unix time in nanoseconds.
//     Example: {{now_nano}}
// - uppercase
//   - Uppercases the previous value.
//     Example: {{.RoleName | uppercase}}
// - lowercase
//   - Lowercases the previous value.
//     Example: {{.DisplayName | lowercase}}
type UsernameProducer struct {
	rawTemplate string
	tmpl        *template.Template
	funcMap     template.FuncMap
	maxLen      int
	casing      casing
}

// NewUsernameProducer creates a UsernameProducer. No arguments are required
// as this has reasonable defaults for all values.
// The default template is specified in the DefaultTemplate constant.
func NewUsernameProducer(opts ...UsernameOpt) (up UsernameProducer, err error) {
	up = UsernameProducer{
		funcMap: map[string]interface{}{
			"rand":        randCharset([]rune(DefaultCharset)),
			"truncate":    truncate,
			"now_seconds": nowSeconds,
			"now_nano":    nowNano,
			"uppercase":   uppercase,
			"lowercase":   lowercase,
		},
		casing: ignoreCase,
	}

	merr := &multierror.Error{}
	for _, opt := range opts {
		merr = multierror.Append(merr, opt(&up))
	}

	err = merr.ErrorOrNil()
	if err != nil {
		return up, err
	}

	tmpl := template.New("usernames").
		Funcs(up.funcMap)
	up.tmpl = tmpl

	if up.rawTemplate != "" {
		_, err = tmpl.Parse(up.rawTemplate)
		if err != nil {
			return up, fmt.Errorf("unable to parse template: %w", err)
		}
	}

	return up, nil
}

func (up UsernameProducer) SetTemplate(rawTemplate string) error {
	_, err := up.tmpl.Parse(rawTemplate)
	if err != nil {
		return err
	}
	return nil
}

// GenerateUsername based on the provided template. This adheres to the CredentialsProducer interface.
func (up UsernameProducer) GenerateUsername(data interface{}) (username string, err error) {
	str := &strings.Builder{}
	err = up.tmpl.Execute(str, data)
	if err != nil {
		return "", fmt.Errorf("unable to apply template: %w", err)
	}
	username = str.String()
	if up.maxLen > 0 && len(username) > up.maxLen {
		username = username[:up.maxLen]
	}

	switch up.casing {
	case upperCase:
		username = strings.ToUpper(username)
	case lowerCase:
		username = strings.ToLower(username)
	}

	return username, nil
}

// randCharset returns a function for use in SetFunction. This allows the user to specify custom charsets for the
// random string generation
func randCharset(charset []rune) func(length int) (randStr string, err error) {
	return func(length int) (randStr string, err error) {
		randRunes, err := random.RandomRunes(rand.Reader, charset, length)
		if err != nil {
			return "", err
		}
		return string(randRunes), nil
	}
}

func nowSeconds() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func nowNano() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func truncate(maxLen int, str string) (string, error) {
	if maxLen <= 0 {
		return str, fmt.Errorf("max length must be > 0 but was %d", maxLen)
	}
	if len(str) > maxLen {
		return str[:maxLen], nil
	}
	return str, nil
}

func uppercase(str string) string {
	return strings.ToUpper(str)
}

func lowercase(str string) string {
	return strings.ToLower(str)
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

type runes []rune

func (r runes) Len() int           { return len(r) }
func (r runes) Less(i, j int) bool { return r[i] < r[j] }
func (r runes) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
