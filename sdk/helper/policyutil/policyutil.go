// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package policyutil

import (
	"errors"
	"sort"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
)

const (
	AddDefaultPolicy      = true
	DoNotAddDefaultPolicy = false
)

var (
	errPolicyNameEmpty         = errors.New("policy name cannot be empty")
	errPolicyNameDotSegment    = errors.New("policy name cannot contain '.' path segments")
	errPolicyNameParentSegment = errors.New("policy name cannot contain '..' path segments")
)

// ValidatePolicyName validates a policy name and returns its canonical form.
func ValidatePolicyName(name string) (string, error) {
	canonicalName := CanonicalizePolicyName(name)

	if canonicalName == "" {
		return "", errPolicyNameEmpty
	}

	for _, segment := range strings.Split(canonicalName, "/") {
		switch segment {
		case ".":
			return "", errPolicyNameDotSegment
		case "..":
			return "", errPolicyNameParentSegment
		}
	}

	return canonicalName, nil
}

// ParsePolicies parses a comma-delimited list of policies.
// The resulting collection will have no duplicate elements.
// If 'root' policy was present in the list of policies, then
// all other policies will be ignored, the result will contain
// just the 'root'. In cases where 'root' is not present, if
// 'default' policy is not already present, it will be added.
func ParsePolicies(policiesRaw interface{}) []string {
	if policiesRaw == nil {
		return []string{"default"}
	}

	var policies []string
	switch policiesRaw.(type) {
	case string:
		if policiesRaw.(string) == "" {
			return []string{}
		}
		policies = strings.Split(policiesRaw.(string), ",")
	case []string:
		policies = policiesRaw.([]string)
	}

	return SanitizePolicies(policies, false)
}

// CanonicalizePolicyName returns the canonical form of a single policy name.
// Policy names are compared and stored case-insensitively and without
// surrounding whitespace, so canonicalization trims leading and trailing
// whitespace (including tabs, newlines, and Unicode space characters, per
// unicode.IsSpace) and lowercases the result.
func CanonicalizePolicyName(policy string) string {
	return strings.ToLower(strings.TrimSpace(policy))
}

// SanitizePolicies performs the common input validation tasks
// which are performed on the list of policies across Vault.
// The resulting collection will have no duplicate elements.
// If 'root' policy was present in the list of policies, then
// all other policies will be ignored, the result will contain
// just the 'root'. In cases where 'root' is not present, if
// 'default' policy is not already present, it will be added
// if addDefault is set to true.
func SanitizePolicies(policies []string, addDefault bool) []string {
	defaultFound := false
	for i, p := range policies {
		policies[i] = CanonicalizePolicyName(p)
		// Eliminate unnamed policies.
		if policies[i] == "" {
			continue
		}

		// If 'root' policy is present, ignore all other policies.
		if policies[i] == "root" {
			policies = []string{"root"}
			defaultFound = true
			break
		}
		if policies[i] == "default" {
			defaultFound = true
		}
	}

	// Always add 'default' except only if the policies contain 'root'.
	if addDefault && (len(policies) == 0 || !defaultFound) {
		policies = append(policies, "default")
	}

	return strutil.RemoveDuplicates(policies, true)
}

// EquivalentPolicies checks whether the given policy sets are equivalent, as in,
// they contain the same values. The benefit of this method is that it leaves
// the "default" policy out of its comparisons as it may be added later by core
// after a set of policies has been saved by a backend and performs policy name
// normalization.
func EquivalentPolicies(a, b []string) bool {
	// First we'll build maps to ensure unique values and filter default
	mapA := map[string]struct{}{}
	mapB := map[string]struct{}{}
	for _, keyA := range a {
		keyA := strings.ToLower(keyA)
		if keyA == "default" {
			continue
		}
		mapA[keyA] = struct{}{}
	}
	for _, keyB := range b {
		keyB := strings.ToLower(keyB)
		if keyB == "default" {
			continue
		}
		mapB[keyB] = struct{}{}
	}

	// Now we'll build our checking slices
	var sortedA, sortedB []string
	for keyA := range mapA {
		sortedA = append(sortedA, keyA)
	}
	for keyB := range mapB {
		sortedB = append(sortedB, keyB)
	}
	sort.Strings(sortedA)
	sort.Strings(sortedB)

	// Finally, compare
	if len(sortedA) != len(sortedB) {
		return false
	}

	for i := range sortedA {
		if sortedA[i] != sortedB[i] {
			return false
		}
	}

	return true
}
