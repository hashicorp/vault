// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package slug

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-slug/internal/ignorefiles"
)

func parseIgnoreFile(rootPath string) *ignorefiles.Ruleset {
	// Look for .terraformignore at our root path/src
	file, err := os.Open(filepath.Join(rootPath, ".terraformignore"))
	defer file.Close()

	// If there's any kind of file error, punt and use the default ignore patterns
	if err != nil {
		// Only show the error debug if an error *other* than IsNotExist
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error reading .terraformignore, default exclusions will apply: %v \n", err)
		}
		return ignorefiles.DefaultRuleset
	}

	ret, err := ignorefiles.ParseIgnoreFileContent(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading .terraformignore, default exclusions will apply: %v \n", err)
		return ignorefiles.DefaultRuleset
	}

	return ret
}

func matchIgnoreRules(path string, ruleset *ignorefiles.Ruleset) ignorefiles.ExcludesResult {
	// Ruleset.Excludes explicitly allows ignoring its error, in which
	// case we are ignoring any individual invalid rules in the set
	// but still taking all of the others into account.
	ret, _ := ruleset.Excludes(path)
	return ret
}
