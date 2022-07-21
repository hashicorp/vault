// DEPRECATED: this has been moved to go-secure-stdlib and will be removed
package strutil

import (
	extstrutil "github.com/hashicorp/go-secure-stdlib/strutil"
)

func StrListContainsGlob(haystack []string, needle string) bool {
	return extstrutil.StrListContainsGlob(haystack, needle)
}

func StrListContains(haystack []string, needle string) bool {
	return extstrutil.StrListContains(haystack, needle)
}

func StrListContainsCaseInsensitive(haystack []string, needle string) bool {
	return extstrutil.StrListContainsCaseInsensitive(haystack, needle)
}

func StrListSubset(super, sub []string) bool {
	return extstrutil.StrListSubset(super, sub)
}

func ParseDedupAndSortStrings(input string, sep string) []string {
	return extstrutil.ParseDedupAndSortStrings(input, sep)
}

func ParseDedupLowercaseAndSortStrings(input string, sep string) []string {
	return extstrutil.ParseDedupLowercaseAndSortStrings(input, sep)
}

func ParseKeyValues(input string, out map[string]string, sep string) error {
	return extstrutil.ParseKeyValues(input, out, sep)
}

func ParseArbitraryKeyValues(input string, out map[string]string, sep string) error {
	return extstrutil.ParseArbitraryKeyValues(input, out, sep)
}

func ParseStringSlice(input string, sep string) []string {
	return extstrutil.ParseStringSlice(input, sep)
}

func ParseArbitraryStringSlice(input string, sep string) []string {
	return extstrutil.ParseArbitraryStringSlice(input, sep)
}

func TrimStrings(items []string) []string {
	return extstrutil.TrimStrings(items)
}

func RemoveDuplicates(items []string, lowercase bool) []string {
	return extstrutil.RemoveDuplicates(items, lowercase)
}

func RemoveDuplicatesStable(items []string, caseInsensitive bool) []string {
	return extstrutil.RemoveDuplicatesStable(items, caseInsensitive)
}

func RemoveEmpty(items []string) []string {
	return extstrutil.RemoveEmpty(items)
}

func EquivalentSlices(a, b []string) bool {
	return extstrutil.EquivalentSlices(a, b)
}

func EqualStringMaps(a, b map[string]string) bool {
	return extstrutil.EqualStringMaps(a, b)
}

func StrListDelete(s []string, d string) []string {
	return extstrutil.StrListDelete(s, d)
}

func GlobbedStringsMatch(item, val string) bool {
	return extstrutil.GlobbedStringsMatch(item, val)
}

func AppendIfMissing(slice []string, i string) []string {
	return extstrutil.AppendIfMissing(slice, i)
}

func MergeSlices(args ...[]string) []string {
	return extstrutil.MergeSlices(args...)
}

func Difference(a, b []string, lowercase bool) []string {
	return extstrutil.Difference(a, b, lowercase)
}

func GetString(m map[string]interface{}, key string) (string, error) {
	return extstrutil.GetString(m, key)
}
