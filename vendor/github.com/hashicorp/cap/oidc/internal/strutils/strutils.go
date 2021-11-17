package strutils

import "strings"

// StrListContains looks for a string in a list of strings.
func StrListContains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// RemoveDuplicatesStable removes duplicate and empty elements from a slice of
// strings, preserving order (and case) of the original slice.
// In all cases, strings are compared after trimming whitespace
// If caseInsensitive, strings will be compared after ToLower()
func RemoveDuplicatesStable(items []string, caseInsensitive bool) []string {
	itemsMap := make(map[string]bool, len(items))
	deduplicated := make([]string, 0, len(items))

	for _, item := range items {
		key := strings.TrimSpace(item)
		if caseInsensitive {
			key = strings.ToLower(key)
		}
		if key == "" || itemsMap[key] {
			continue
		}
		itemsMap[key] = true
		deduplicated = append(deduplicated, item)
	}
	return deduplicated
}
