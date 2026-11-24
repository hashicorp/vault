// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/vault/helper/testhelpers/testimages"
)

func main() {
	var source, target, binary string
	var hsm bool

	flag.StringVar(&source, "source", "", "Source image name")
	flag.StringVar(&target, "target", "", "Target image name")
	flag.StringVar(&binary, "binary", "", "Binary path")
	flag.BoolVar(&hsm, "hsm", false, "HSM style image")
	flag.Parse()

	if source == "" || target == "" || binary == "" {
		fmt.Fprintf(os.Stderr, "Error: all of the flags -source, -target, and -binary are required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	var output []byte
	var err error
	if hsm {
		output, err = testimages.CreateHSMDockerImage(source, target, binary)
	} else {
		output, err = testimages.CreateNonHSMDockerImage(source, target, binary)
	}
	fmt.Println(string(output))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
