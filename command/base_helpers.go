// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	kvbuilder "github.com/hashicorp/go-secure-stdlib/kv-builder"
	"github.com/hashicorp/vault/api"
	"github.com/kr/text"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/ryanuber/columnize"
)

// extractListData reads the secret and returns a typed list of data and a
// boolean indicating whether the extraction was successful.
func extractListData(secret *api.Secret) ([]interface{}, bool) {
	if secret == nil || secret.Data == nil {
		return nil, false
	}

	k, ok := secret.Data["keys"]
	if !ok || k == nil {
		return nil, false
	}

	i, ok := k.([]interface{})
	return i, ok
}

// sanitizePath removes any leading or trailing things from a "path".
func sanitizePath(s string) string {
	return ensureNoTrailingSlash(ensureNoLeadingSlash(s))
}

// ensureTrailingSlash ensures the given string has a trailing slash.
func ensureTrailingSlash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	for len(s) > 0 && s[len(s)-1] != '/' {
		s = s + "/"
	}
	return s
}

// ensureNoTrailingSlash ensures the given string does not have a trailing slash.
func ensureNoTrailingSlash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	for len(s) > 0 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

// ensureNoLeadingSlash ensures the given string does not have a leading slash.
func ensureNoLeadingSlash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	for len(s) > 0 && s[0] == '/' {
		s = s[1:]
	}
	return s
}

// columnOutput prints the list of items as a table with no headers.
func columnOutput(list []string, c *columnize.Config) string {
	if len(list) == 0 {
		return ""
	}

	if c == nil {
		c = &columnize.Config{}
	}
	if c.Glue == "" {
		c.Glue = "    "
	}
	if c.Empty == "" {
		c.Empty = "n/a"
	}

	return columnize.Format(list, c)
}

// tableOutput prints the list of items as columns, where the first row is
// the list of headers.
func tableOutput(list []string, c *columnize.Config) string {
	if len(list) == 0 {
		return ""
	}

	delim := "|"
	if c != nil && c.Delim != "" {
		delim = c.Delim
	}

	underline := ""
	headers := strings.Split(list[0], delim)
	for i, h := range headers {
		h = strings.TrimSpace(h)
		u := strings.Repeat("-", len(h))

		underline = underline + u
		if i != len(headers)-1 {
			underline = underline + delim
		}
	}

	list = append(list, "")
	copy(list[2:], list[1:])
	list[1] = underline

	return columnOutput(list, c)
}

// parseArgsData parses the given args in the format key=value into a map of
// the provided arguments. The given reader can also supply key=value pairs.
func parseArgsData(stdin io.Reader, args []string) (map[string]interface{}, error) {
	builder := &kvbuilder.Builder{Stdin: stdin}
	if err := builder.Add(args...); err != nil {
		return nil, err
	}

	return builder.Map(), nil
}

// parseArgsDataString parses the args data and returns the values as strings.
// If the values cannot be represented as strings, an error is returned.
func parseArgsDataString(stdin io.Reader, args []string) (map[string]string, error) {
	raw, err := parseArgsData(stdin, args)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	if err := mapstructure.WeakDecode(raw, &result); err != nil {
		return nil, errors.Wrap(err, "failed to convert values to strings")
	}
	if result == nil {
		result = make(map[string]string)
	}
	return result, nil
}

// parseArgsDataStringLists parses the args data and returns the values as
// string lists. If the values cannot be represented as strings, an error is
// returned.
func parseArgsDataStringLists(stdin io.Reader, args []string) (map[string][]string, error) {
	raw, err := parseArgsData(stdin, args)
	if err != nil {
		return nil, err
	}

	var result map[string][]string
	if err := mapstructure.WeakDecode(raw, &result); err != nil {
		return nil, errors.Wrap(err, "failed to convert values to strings")
	}
	return result, nil
}

// truncateToSeconds truncates the given duration to the number of seconds. If
// the duration is less than 1s, it is returned as 0. The integer represents
// the whole number unit of seconds for the duration.
func truncateToSeconds(d time.Duration) int {
	d = d.Truncate(1 * time.Second)

	// Handle the case where someone requested a ridiculously short increment -
	// increments must be larger than a second.
	if d < 1*time.Second {
		return 0
	}

	return int(d.Seconds())
}

// printKeyStatus prints the KeyStatus response from the API.
func printKeyStatus(ks *api.KeyStatus) string {
	return columnOutput([]string{
		fmt.Sprintf("Key Term | %d", ks.Term),
		fmt.Sprintf("Install Time | %s", ks.InstallTime.UTC().Format(time.RFC822)),
		fmt.Sprintf("Encryption Count | %d", ks.Encryptions),
	}, nil)
}

// expandPath takes a filepath and returns the full expanded path, accounting
// for user-relative things like ~/.
func expandPath(s string) string {
	if s == "" {
		return ""
	}

	e, err := homedir.Expand(s)
	if err != nil {
		return s
	}
	return e
}

// wrapAtLengthWithPadding wraps the given text at the maxLineLength, taking
// into account any provided left padding.
func wrapAtLengthWithPadding(s string, pad int) string {
	wrapped := text.Wrap(s, maxLineLength-pad)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		lines[i] = strings.Repeat(" ", pad) + line
	}
	return strings.Join(lines, "\n")
}

// wrapAtLength wraps the given text to maxLineLength.
func wrapAtLength(s string) string {
	return wrapAtLengthWithPadding(s, 0)
}

// ttlToAPI converts a user-supplied ttl into an API-compatible string. If
// the TTL is 0, this returns the empty string. If the TTL is negative, this
// returns "system" to indicate to use the system values. Otherwise, the
// time.Duration ttl is used.
func ttlToAPI(d time.Duration) string {
	if d == 0 {
		return ""
	}

	if d < 0 {
		return "system"
	}

	return d.String()
}

// humanDuration prints the time duration without those pesky zeros.
func humanDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}

	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if idx := strings.Index(s, "h0m"); idx > 0 {
		s = s[:idx+1] + s[idx+3:]
	}
	return s
}

// humanDurationInt prints the given int as if it were a time.Duration  number
// of seconds.
func humanDurationInt(i interface{}) interface{} {
	switch i := i.(type) {
	case int:
		return humanDuration(time.Duration(i) * time.Second)
	case int64:
		return humanDuration(time.Duration(i) * time.Second)
	case json.Number:
		if i, err := i.Int64(); err == nil {
			return humanDuration(time.Duration(i) * time.Second)
		}
	}

	// If we don't know what type it is, just return the original value
	return i
}

// parseFlagFile accepts a flag value returns the contets of that value. If the
// value starts with '@', that indicates the value is a file and its content
// should be read and returned. Otherwise, the raw value is returned.
func parseFlagFile(raw string) (string, error) {
	// check if the provided argument should be read from file
	if len(raw) > 0 && raw[0] == '@' {
		contents, err := ioutil.ReadFile(raw[1:])
		if err != nil {
			return "", fmt.Errorf("error reading file: %w", err)
		}

		return string(contents), nil
	}

	return raw, nil
}

func generateFlagWarnings(args []string) string {
	var trailingFlags []string
	for _, arg := range args {
		// "-" can be used where a file is expected to denote stdin.
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			continue
		}

		isGlobalFlag := false
		trimmedArg, _, _ := strings.Cut(strings.TrimLeft(arg, "-"), "=")
		for _, flag := range globalFlags {
			if trimmedArg == flag {
				isGlobalFlag = true
			}
		}
		if isGlobalFlag {
			continue
		}

		trailingFlags = append(trailingFlags, arg)
	}

	if len(trailingFlags) > 0 {
		return fmt.Sprintf("Command flags must be provided before positional arguments. "+
			"The following arguments will not be parsed as flags: [%s]", strings.Join(trailingFlags, ","))
	} else {
		return ""
	}
}

func generateFlagErrors(f *FlagSets, opts ...ParseOptions) error {
	if Format(f.ui) == "raw" {
		canUseRaw := false
		for _, opt := range opts {
			if value, ok := opt.(ParseOptionAllowRawFormat); ok {
				canUseRaw = bool(value)
			}
		}

		if !canUseRaw {
			return fmt.Errorf("This command does not support the -format=raw option.")
		}
	}

	return nil
}
