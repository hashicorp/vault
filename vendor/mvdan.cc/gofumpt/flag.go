// Copyright (c) 2019, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package main

import "flag"

var (
	langVersion = flag.String("lang", "", "target Go version in the form 1.X (default from go.mod)")
	extraRules  = flag.Bool("extra", false, "enable extra rules which should be vetted by a human")
)

func init() {
	// make -s default to true
	*simplifyAST = true
}
