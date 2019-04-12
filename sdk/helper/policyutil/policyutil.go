package policyutil

import (
	"sort"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/strutil"
)

const (
	AddDefaultPolicy      = true
	DoNotAddDefaultPolicy = false
)

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
		policies[i] = strings.ToLower(strings.TrimSpace(p))
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
// after a set of policies has been saved by a backend.
func EquivalentPolicies(a, b []string) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil && len(b) == 1 && b[0] == "default":
		return true
	case b == nil && len(a) == 1 && a[0] == "default":
		return true
	case a == nil || b == nil:
		return false
	}

	// First we'll build maps to ensure unique values and filter default
	mapA := map[string]bool{}
	mapB := map[string]bool{}
	for _, keyA := range a {
		if keyA == "default" {
			continue
		}
		mapA[keyA] = true
	}
	for _, keyB := range b {
		if keyB == "default" {
			continue
		}
		mapB[keyB] = true
	}

	// Now we'll build our checking slices
	var sortedA, sortedB []string
	for keyA, _ := range mapA {
		sortedA = append(sortedA, keyA)
	}
	for keyB, _ := range mapB {
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
