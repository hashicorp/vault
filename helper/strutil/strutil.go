package strutil

import (
	"sort"
	"strings"
)

// StrListContains looks for a string in a list of strings.
func StrListContains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// StrListSubset checks if a given list is a subset
// of another set
func StrListSubset(super, sub []string) bool {
	for _, item := range sub {
		if !StrListContains(super, item) {
			return false
		}
	}
	return true
}

// Parses a comma separated list of strings into a slice of strings.
// The return slice will be sorted and will not contain duplicate or
// empty items. The values will be converted to lower case.
func ParseStrings(input string) []string {
	input = strings.TrimSpace(input)
	var parsed []string
	if input == "" {
		// Don't return nil
		return parsed
	}
	return RemoveDuplicates(strings.Split(input, ","))
}

// Removes duplicate and empty elements from a slice of strings.
// This also converts the items in the slice to lower case and
// returns a sorted slice.
func RemoveDuplicates(items []string) []string {
	itemsMap := map[string]bool{}
	for _, item := range items {
		item = strings.ToLower(strings.TrimSpace(item))
		if item == "" {
			continue
		}
		itemsMap[item] = true
	}
	items = []string{}
	for item, _ := range itemsMap {
		items = append(items, item)
	}
	sort.Strings(items)
	return items
}

// EquivalentSlices checks whether the given string sets are equivalent, as in,
// they contain the same values.
func EquivalentSlices(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	// First we'll build maps to ensure unique values
	mapA := map[string]bool{}
	mapB := map[string]bool{}
	for _, keyA := range a {
		mapA[keyA] = true
	}
	for _, keyB := range b {
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
