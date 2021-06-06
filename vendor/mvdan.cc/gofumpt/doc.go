// Copyright (c) 2019, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package main

// First, sync the files with x/tools and GOROOT.
//go:generate go run gen.go

// Then, add the missing imports to our added code.
//go:generate goimports -w .

// Finally, ensure all code follows 'gofumpt -s'. Use the current source, to not
// need an extra 'go install' step.
//go:generate go run . -s -w .
