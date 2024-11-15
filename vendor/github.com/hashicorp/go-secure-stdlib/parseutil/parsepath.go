// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parseutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

var (
	ErrNotAUrl   = errors.New("not a url")
	ErrNotParsed = errors.New("not a parsed value")
)

// ParsePath parses a URL with schemes file://, env://, or any other. Depending
// on the scheme it will return specific types of data:
//
// * file:// will return a string with the file's contents
//
// * env:// will return a string with the env var's contents
//
// * Anything else will return the string as it was. Functionally this means
// anything for which Go's `url.Parse` function does not throw an error. If you
// want to ensure that this function errors if a known scheme is not found, use
// MustParsePath.
//
// On error, we return the original string along with the error. The caller can
// switch on errors.Is(err, ErrNotAUrl) to understand whether it was the parsing
// step that errored or something else (such as a file not found). This is
// useful to attempt to read a non-URL string from some resource, but where the
// original input may simply be a valid string of that type.
func ParsePath(path string) (string, error) {
	return parsePath(path, false)
}

// MustParsePath behaves like ParsePath but will return ErrNotAUrl if the value
// is not a URL with a scheme that can be parsed by this function.
func MustParsePath(path string) (string, error) {
	return parsePath(path, true)
}

func parsePath(path string, mustParse bool) (string, error) {
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
	default:
		if mustParse {
			return "", ErrNotParsed
		}
		return path, nil
	}
}
