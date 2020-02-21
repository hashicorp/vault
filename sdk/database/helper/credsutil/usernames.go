package credsutil

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
)

const (
	DefaultCharset = LowerCharset + UpperCharset + NumericCharset
	UpperCharset   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LowerCharset   = "abcdefghijklmnopqrstuvwxyz"
	NumericCharset = "0123456789"
	SymbolsCharset = "!\"#$%&'()*+,-./:;<=>?`\\~[]^_@|" // no whitespace characters for safety

	DefaultTemplate = "v_{{.RoleName}}_{{.DisplayName}}_{{now_seconds}}_{{rand 10}}"
)

type UsernameOpt func(*UsernameProducer) error

// UsernameTemplate specifies the template to use when generating a username.
func UsernameTemplate(pattern string) UsernameOpt {
	return func(up *UsernameProducer) error {
		up.rawTemplate = pattern
		return nil
	}
}

// UsernameFuncMap allows the user to specify functions for use in the template. This can be used to override the
// behavior of any of the built-in functions (such as rand). Overriding the built-in functions is not recommended.
func UsernameFuncMap(name string, f interface{}) UsernameOpt {
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

// UsernameMaxLength is the maximum allowed length of the username. This will take effect after applying the template.
// A length of 0 indicates no maximum length
func UsernameMaxLength(maxLen int) UsernameOpt {
	return func(up *UsernameProducer) error {
		if maxLen < 0 {
			return fmt.Errorf("max username length must be >= 0")
		}
		up.maxLen = maxLen
		return nil
	}
}

// UsernameRandomCharset specifies the charset used when using the `rand` template function.
// This will remove duplicate characters from the charset to prevent bias.
func UsernameRandomCharset(charset string) UsernameOpt {
	if charset == "" {
		charset = DefaultCharset
	}
	charset = removeDuplicateChars(charset)
	f := func(length int) (randStr string, err error) {
		return randomChars(charset, length)
	}
	return UsernameFuncMap("rand", f)
}

// UsernameProducer creates usernames based on the provided template.
// This uses the go templating language, so anything that adheres to that language will function in this struct.
// There are several custom functions available for use in the template:
// - rand
//   - Randomly generated characters. This uses the charset specified in UsernameRandomCharset. Must include a length.
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
}

// NewUsernameProducer creates a UsernameProducer. No arguments are required
// as this has reasonable defaults for all values.
// The default template is specified in the DefaultTemplate constant.
func NewUsernameProducer(opts ...UsernameOpt) (up UsernameProducer, err error) {
	up = UsernameProducer{
		rawTemplate: DefaultTemplate,
		funcMap: map[string]interface{}{
			"rand":        randCharset(DefaultCharset),
			"truncate":    truncate,
			"now_seconds": nowSeconds,
			"now_nano":    nowNano,
			"uppercase":   uppercase,
			"lowercase":   lowercase,
		},
	}

	merr := &multierror.Error{}
	for _, opt := range opts {
		merr = multierror.Append(merr, opt(&up))
	}

	err = merr.ErrorOrNil()
	if err != nil {
		return up, err
	}

	tmpl := template.New("usernames")
	tmpl.Funcs(up.funcMap)
	tmpl, err = tmpl.Parse(up.rawTemplate)
	if err != nil {
		return up, fmt.Errorf("unable to parse template: %w", err)
	}
	up.tmpl = tmpl

	return up, nil
}

// GenerateUsername based on the provided template. This adheres to the CredentialsProducer interface.
func (up UsernameProducer) GenerateUsername(config dbplugin.UsernameConfig) (username string, err error) {
	str := &strings.Builder{}
	err = up.tmpl.Execute(str, config)
	if err != nil {
		return "", fmt.Errorf("unable to apply template: %w", err)
	}
	username = str.String()
	if up.maxLen > 0 && len(username) > up.maxLen {
		username = username[:up.maxLen]
	}
	return username, nil
}

// randCharset returns a function for use in UsernameFuncMap. This allows the user to specify custom charsets for the
// random string generation
func randCharset(charset string) func(length int) (randStr string, err error) {
	return func(length int) (randStr string, err error) {
		return randomChars(charset, length)
	}
}

// randomChars creates a string of the provided length
func randomChars(charset string, length int) (randStr string, err error) {
	if length < 0 {
		return "", fmt.Errorf("length must be >= 0")
	}
	reader := rand.Reader

	output := make([]byte, 0, length)

	// Request a bit more than length to reduce the chance
	// of needing more than one batch of random bytes
	batchSize := length + length/4

	for {
		buf, err := uuid.GenerateRandomBytesWithReader(batchSize, reader)
		if err != nil {
			return "", fmt.Errorf("unable to generate random bytes: %w", err)
		}

		csLen := byte(len(charset))

		for _, b := range buf {
			// Avoid bias by using a value range that's a multiple of the charset size
			if b < (csLen * 4) {
				output = append(output, charset[b%csLen])

				if len(output) == length {
					return string(output), nil
				}
			}
		}
	}
}

func randIntBetween(min, max int) int {
	if min == max {
		return min
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1))) // +1 so it's inclusive on both ends
	if err != nil {
		panic("unable to get random integer")
	}
	return int(nBig.Int64()) + min
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

// removeDuplicateChars from the provided string. Does not care about order. Will keep the first unique character.
func removeDuplicateChars(str string) string {
	seen := map[string]bool{}
	newStr := &strings.Builder{}
	for i := range str {
		v := str[i : i+1]
		if seen[v] {
			continue
		}
		newStr.WriteString(v)
		seen[v] = true
	}
	return newStr.String()
}
