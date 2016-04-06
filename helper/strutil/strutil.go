package strutil

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
