package policies

import "sort"

// ComparePolicies checks whether the given policy sets are equivalent, as in,
// they contain the same values. The benefit of this method is that it leaves
// the "default" policy out of its comparisons as it may be added later by core
// after a set of policies has been saved by a backend.
func EquivalentPolicies(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
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
