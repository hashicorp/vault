// Copyright (c) 2015-present Jeevanandam M (jeeva@myjeeva.com)
// 2024 Ahuigo (https://github.com/ahuigo)
// All rights reserved.
// resty source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

/*
Package shellescape provides the methods to escape arbitrary
strings for a safe use as command line arguments in the most common
POSIX shells.

The original Python package which this work was inspired by can be found
at https://pypi.python.org/pypi/shellescape.
*/
package shellescape

import (
	"regexp"
	"strings"
)

var pattern *regexp.Regexp

func init() {
	pattern = regexp.MustCompile(`[^\w@%+=:,./-]`)
}

// Quote method returns a shell-escaped version of the string. The returned value
// can safely be used as one token in a shell command line.
func Quote(s string) string {
	if len(s) == 0 {
		return "''"
	}

	if pattern.MatchString(s) {
		return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
	}

	return s
}
