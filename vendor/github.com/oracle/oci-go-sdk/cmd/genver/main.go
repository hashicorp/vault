// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

// Package main The following code is used to generate the version of the go sdk
// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

// Environment variables for version information
const (
	majorVer = "VER_MAJOR"
	minorVer = "VER_MINOR"
	patchVer = "VER_PATCH"
	tag      = "VER_TAG"
)

var output = flag.String("output", "", "output for the file")

func getEnvOrDefault(envkey, defaultValue string) string {
	if val := os.Getenv(envkey); val != "" {
		return val
	}
	return defaultValue
}

// Reads the output file as a flag and version information from environment variables
func main() {
	flag.Parse()
	genTemplate := template.Must(template.New("version").Parse(versionTemplate))

	versions := struct {
		Major, Minor, Patch string
		Tag                 string
	}{
		getEnvOrDefault(majorVer, "0"),
		getEnvOrDefault(minorVer, "0"),
		getEnvOrDefault(patchVer, "0"),
		getEnvOrDefault(tag, ""),
	}

	var buf bytes.Buffer

	if err := genTemplate.Execute(&buf, versions); err != nil {
		log.Printf("error while generation version: %s", err)
		return
	}

	if *output == "" {
		fmt.Print(buf.String())
		return
	}

	if err := ioutil.WriteFile(*output, buf.Bytes(), 0644); err != nil {
		log.Printf("could not write output file: %s", err)
		return
	}
}
