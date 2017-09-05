package command

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	kvbuilder "github.com/hashicorp/vault/helper/kv-builder"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/ryanuber/columnize"
)

var ErrMissingID = fmt.Errorf("Missing ID!")
var ErrMissingPath = fmt.Errorf("Missing PATH!")
var ErrMissingThing = fmt.Errorf("Missing THING!")

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

// extractPath extracts the path and list of arguments from the args. If there
// are no extra arguments, the remaining args will be nil.
func extractPath(args []string) (string, []string, error) {
	str, remaining, err := extractThings(args)
	if err == ErrMissingThing {
		err = ErrMissingPath
	}
	return str, remaining, err
}

// extractID extracts the path and list of arguments from the args. If there
// are no extra arguments, the remaining args will be nil.
func extractID(args []string) (string, []string, error) {
	str, remaining, err := extractThings(args)
	if err == ErrMissingThing {
		err = ErrMissingID
	}
	return str, remaining, err
}

func extractThings(args []string) (string, []string, error) {
	if len(args) < 1 {
		return "", nil, ErrMissingThing
	}

	// Path is always the first argument after all flags
	thing := args[0]

	// Strip leading and trailing slashes
	thing = sanitizePath(thing)

	// Verify we have a thing
	if thing == "" {
		return "", nil, ErrMissingThing
	}

	// Splice remaining args
	var remaining []string
	if len(args) > 1 {
		remaining = args[1:]
	}

	return thing, remaining, nil
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

// ensureNoTrailingSlash ensures the given string has a trailing slash.
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

// ensureNoLeadingSlash ensures the given string has a trailing slash.
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

// columnOuput prints the list of items as a table with no headers.
func columnOutput(list []string) string {
	if len(list) == 0 {
		return ""
	}

	return columnize.Format(list, &columnize.Config{
		Glue:  "    ",
		Empty: "n/a",
	})
}

// tableOutput prints the list of items as columns, where the first row is
// the list of headers.
func tableOutput(list []string) string {
	if len(list) == 0 {
		return ""
	}

	underline := ""
	headers := strings.Split(list[0], "|")
	for i, h := range headers {
		h = strings.TrimSpace(h)
		u := strings.Repeat("-", len(h))

		underline = underline + u
		if i != len(headers)-1 {
			underline = underline + " | "
		}
	}

	list = append(list, "")
	copy(list[2:], list[1:])
	list[1] = underline

	return columnOutput(list)
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
	return result, nil
}

// truncateToSeconds truncates the given duaration to the number of seconds. If
// the duration is less than 1s, it is returned as 0. The integer represents
// the whole number unit of seconds for the duration.
func truncateToSeconds(d time.Duration) int {
	d = d.Truncate(1 * time.Second)

	// Handle the case where someone requested a ridiculously short increment -
	// incremenents must be larger than a second.
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
	})
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
