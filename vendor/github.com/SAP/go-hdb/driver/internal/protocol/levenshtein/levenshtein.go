// Package levenshtein includes the levenshtein distance algorithm plus additional helper functions.
// The algorithm is taken from https://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Levenshtein_distance#Go.
package levenshtein

import (
	"math"
	"strings"
	"unicode/utf8"
)

// Distance returns the Lewenshtein distance.
func Distance(a, b string, caseSensitive bool) int {
	if caseSensitive {
		return distance(a, b)
	}
	return distance(strings.ToLower(a), strings.ToLower(b))
}

// MinString returns the string attribute determined by fn out of x with the minimal Lewenshtein distance to s.
func MinString[S ~[]E, E any](x S, fn func(a E) string, s string, caseSensitive bool) (rv string) {
	minInt := math.MaxInt
	for _, e := range x {
		xs := fn(e)
		if d := Distance(xs, s, caseSensitive); d < minInt {
			rv = xs
			minInt = d
		}
	}
	return
}

func distance(a, b string) int {
	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1) // delete & insert
			if cb != ca {
				mn = min(mn, fj1+1) // change
			} else {
				mn = min(mn, fj1) // matched
			}

			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
		}
	}

	return f[len(f)-1]
}
