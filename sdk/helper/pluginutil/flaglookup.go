// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package pluginutil

import "flag"

// FlagLookup is a minimal interface over a parsed CLI flag set. It allows
// LoginHandler implementations to read flags by name without importing the
// command package, avoiding an import cycle between vault-enterprise and
// auth plugins.
//
// *command.FlagSets satisfies this interface via its Lookup method.
type FlagLookup interface {
	Lookup(name string) *flag.Flag
}
