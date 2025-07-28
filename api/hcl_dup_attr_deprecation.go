// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	hclParser "github.com/hashicorp/hcl/hcl/parser"
)

// allowHclDuplicatesEnvVar is an environment variable that allows Vault to revert back to accepting HCL files with
// duplicate attributes. It's temporary until we finish the deprecation process, at which point this will be removed
const allowHclDuplicatesEnvVar = "VAULT_ALLOW_PENDING_REMOVAL_DUPLICATE_HCL_ATTRIBUTES"

// parseAndCheckForDuplicateHclAttributes parses the input JSON/HCL file and if it is HCL it also checks
// for duplicate keys in the HCL file, allowing callers to handle the issue accordingly. It now only accepts duplicate
// keys if the environment variable VAULT_ALLOW_PENDING_REMOVAL_DUPLICATE_HCL_ATTRIBUTES is set to true. In a future
// release we'll remove this function entirely and there will be no way to parse HCL files with duplicate keys.
// TODO (HCL_DUP_KEYS_DEPRECATION): remove once not used anymore
func parseAndCheckForDuplicateHclAttributes(input string) (res *ast.File, duplicate bool, err error) {
	res, err = hcl.Parse(input)
	if err != nil && strings.Contains(err.Error(), "Each argument can only be defined once") {
		allowHclDuplicatesRaw := os.Getenv(allowHclDuplicatesEnvVar)
		if allowHclDuplicatesRaw == "" {
			// default is to not allow duplicates
			return nil, false, err
		}
		allowHclDuplicates, envParseErr := strconv.ParseBool(allowHclDuplicatesRaw)
		if envParseErr != nil {
			return nil, false, fmt.Errorf("error parsing %q environment variable: %w", allowHclDuplicatesEnvVar, err)
		}
		if !allowHclDuplicates {
			return nil, false, err
		}

		// if allowed by the environment variable, parse again without failing on duplicate attributes
		duplicate = true
		res, err = hclParser.ParseDontErrorOnDuplicateKeys([]byte(input))
	}
	return res, duplicate, err
}
