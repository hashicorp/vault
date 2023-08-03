// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package pki

func getEntProperAuthingPaths(_ string) map[string]pathAuthChecker {
	return map[string]pathAuthChecker{}
}

func getEntAcmePrefixes() []string {
	return []string{}
}

func entProperAuthingPathReplacer(rawPath string) string {
	return rawPath
}
