package parseutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

var ErrNotAUrl = errors.New("not a url")

// ParsePath parses a URL with schemes file://, env://, or any other. Depending
// on the scheme it will return specific types of data:
//
// * file:// will return a string with the file's contents
//
// * env:// will return a string with the env var's contents
//
// * Anything else will return the string as it was
//
// On error, we return the original string along with the error. The caller can
// switch on errors.Is(err, ErrNotAUrl) to understand whether it was the parsing
// step that errored or something else (such as a file not found). This is
// useful to attempt to read a non-URL string from some resource, but where the
// original input may simply be a valid string of that type.
func ParsePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	parsed, err := url.Parse(path)
	if err != nil {
		return path, fmt.Errorf("error parsing url (%q): %w", err.Error(), ErrNotAUrl)
	}
	switch parsed.Scheme {
	case "file":
		contents, err := ioutil.ReadFile(strings.TrimPrefix(path, "file://"))
		if err != nil {
			return path, fmt.Errorf("error reading file at %s: %w", path, err)
		}
		return strings.TrimSpace(string(contents)), nil
	case "env":
		return strings.TrimSpace(os.Getenv(strings.TrimPrefix(path, "env://"))), nil
	}

	return path, nil
}
