package template

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/helper/base62"
)

type Opt func(*StringTemplate) error

func Template(rawTemplate string) Opt {
	return func(up *StringTemplate) error {
		up.rawTemplate = rawTemplate
		return nil
	}
}

// Function allows the user to specify functions for use in the template. If the name provided is a function that
// already exists in the function map, this will override the previously specified function.
func Function(name string, f interface{}) Opt {
	return func(up *StringTemplate) error {
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

// StringTemplate creates strings based on the provided template.
// This uses the go templating language, so anything that adheres to that language will function in this struct.
// There are several custom functions available for use in the template:
// - random
//   - Randomly generated characters. This uses the charset specified in RandomCharset. Must include a length.
//     Example: {{ rand 20 }}
// - truncate
//   - Truncates the previous value to the specified length. Must include a maximum length.
//     Example: {{ .DisplayName | truncate 10 }}
// - truncate_sha256
//   - Truncates the previous value to the specified length. If the original length is greater than the length
//     specified, the remaining characters will be sha256 hashed and appended to the end. The hash will be only the first 8 characters The maximum length will
//     be no longer than the length specified.
//     Example: {{ .DisplayName | truncate_sha256 30 }}
// - uppercase
//   - Uppercases the previous value.
//     Example: {{ .RoleName | uppercase }}
// - lowercase
//   - Lowercases the previous value.
//     Example: {{ .DisplayName | lowercase }}
// - replace
//   - Performs a string find & replace
//     Example: {{ .DisplayName | replace - _ }}
// - sha256
//   - SHA256 hashes the previous value.
//     Example: {{ .DisplayName | sha256 }}
// - base64
//   - base64 encodes the previous value.
//     Example: {{ .DisplayName | base64 }}
// - unix_time
//   - Provides the current unix time in seconds.
//     Example: {{ unix_time }}
// - unix_time_millis
//   - Provides the current unix time in milliseconds.
//     Example: {{ unix_time_millis }}
// - timestamp
//   - Provides the current time. Must include a standard Go format string
// - uuid
//   - Generates a UUID
//     Example: {{ uuid }}
type StringTemplate struct {
	rawTemplate string
	tmpl        *template.Template
	funcMap     template.FuncMap
}

// NewTemplate creates a StringTemplate. No arguments are required
// as this has reasonable defaults for all values.
// The default template is specified in the DefaultTemplate constant.
func NewTemplate(opts ...Opt) (up StringTemplate, err error) {
	up = StringTemplate{
		funcMap: map[string]interface{}{
			"random":          base62.Random,
			"truncate":        truncate,
			"truncate_sha256": truncateSHA256,
			"uppercase":       uppercase,
			"lowercase":       lowercase,
			"replace":         replace,
			"sha256":          hashSHA256,
			"base64":          encodeBase64,

			"unix_time":        unixTime,
			"unix_time_millis": unixTimeMillis,
			"timestamp":        timestamp,
			"uuid":             uuid,
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

	if up.rawTemplate == "" {
		return StringTemplate{}, fmt.Errorf("missing template")
	}

	tmpl, err := template.New("template").
		Funcs(up.funcMap).
		Parse(up.rawTemplate)
	if err != nil {
		return StringTemplate{}, fmt.Errorf("unable to parse template: %w", err)
	}
	up.tmpl = tmpl

	return up, nil
}

// Generate based on the provided template
func (up StringTemplate) Generate(data interface{}) (string, error) {
	if up.tmpl == nil || up.rawTemplate == "" {
		return "", fmt.Errorf("failed to generate: template not initialized")
	}
	str := &strings.Builder{}
	err := up.tmpl.Execute(str, data)
	if err != nil {
		return "", fmt.Errorf("unable to apply template: %w", err)
	}

	return str.String(), nil
}
