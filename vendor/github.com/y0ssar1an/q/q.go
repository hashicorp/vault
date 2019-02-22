// Copyright 2016 Ryan Boehning. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package q

import (
	"bytes"
	"fmt"
)

// nolint: gochecknoglobals
var (
	// std is the singleton logger.
	std = &logger{
		buf: &bytes.Buffer{},
	}
)

// Q pretty-prints the given arguments to the $TMPDIR/q log file.
func Q(v ...interface{}) {
	std.mu.Lock()
	defer std.mu.Unlock()

	// Flush the buffered writes to disk.
	defer func() {
		if err := std.flush(); err != nil {
			fmt.Println(err)
		}
	}()

	args := formatArgs(v...)
	funcName, file, line, err := getCallerInfo()
	if err != nil {
		std.output(args...) // no name=value printing
		return
	}

	// Print a header line if this q.Q() call is in a different file or
	// function than the previous q.Q() call, or if the 2s timer expired.
	// A header line looks like this: [14:00:36 main.go main.main:122].
	header := std.header(funcName, file, line)
	if header != "" {
		fmt.Fprint(std.buf, "\n", header, "\n")
	}

	// q.Q(foo, bar, baz) -> []string{"foo", "bar", "baz"}
	names, err := argNames(file, line)
	if err != nil {
		std.output(args...) // no name=value printing
		return
	}

	// Convert the arguments to name=value strings.
	args = prependArgName(names, args)
	std.output(args...)
}
