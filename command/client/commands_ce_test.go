// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package command

// entClientCommands is empty in CE: enterprise-only commands are not built.
var entClientCommands []string

// entLoginHandlers is empty in CE: enterprise-only login handlers are not built.
var entLoginHandlers []string
