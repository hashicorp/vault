// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// vendor tool to copy external source code from GOPATH or remote location to the
// local vendor folder. See README.md for usage.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kardianos/govendor/cliprompt"
	"github.com/kardianos/govendor/help"
	"github.com/kardianos/govendor/run"
)

func main() {
	prompt := &cliprompt.Prompt{}

	allArgs := os.Args

	if allArgs[len(allArgs)-1] == "-" {
		stdin := &bytes.Buffer{}
		if _, err := io.Copy(stdin, os.Stdin); err == nil {
			stdinArgs := strings.Fields(stdin.String())
			allArgs = append(allArgs[:len(allArgs)-1], stdinArgs...)
		}
	}

	msg, err := run.Run(os.Stdout, allArgs, prompt)
	if err == flag.ErrHelp {
		err = nil
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	msgText := msg.String()
	if len(msgText) > 0 {
		fmt.Fprint(os.Stderr, msgText)
	}
	if err != nil {
		os.Exit(2)
	}
	if msg != help.MsgNone {
		os.Exit(1)
	}
}
