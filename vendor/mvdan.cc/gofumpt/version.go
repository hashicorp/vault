// Copyright (c) 2020, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package main

import (
	"flag"
	"fmt"
	"runtime/debug"
)

var (
	showVersion = flag.Bool("version", false, "show version and exit")

	version = "(devel)" // to match the default from runtime/debug
)

func printVersion() {
	// don't overwrite the version if it was set by -ldflags=-X
	if info, ok := debug.ReadBuildInfo(); ok && version == "(devel)" {
		mod := &info.Main
		if mod.Replace != nil {
			mod = mod.Replace
		}
		version = mod.Version
	}
	fmt.Println(version)
}
