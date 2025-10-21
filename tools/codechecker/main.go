// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"github.com/hashicorp/vault/tools/codechecker/pkg/godoctests"
	"github.com/hashicorp/vault/tools/codechecker/pkg/gonilnilfunctions"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(gonilnilfunctions.Analyzer, godoctests.Analyzer)
}
