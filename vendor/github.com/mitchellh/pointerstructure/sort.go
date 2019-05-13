package pointerstructure

import (
	"sort"
)

// Sort does an in-place sort of the pointers so that they are in order
// of least specific to most specific alphabetized. For example:
// "/foo", "/foo/0", "/qux"
//
// This ordering is ideal for applying the changes in a way that ensures
// that parents are set first.
func Sort(p []*Pointer) { sort.Sort(PointerSlice(p)) }

// PointerSlice is a slice of pointers that adheres to sort.Interface
type PointerSlice []*Pointer

func (p PointerSlice) Len() int      { return len(p) }
func (p PointerSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PointerSlice) Less(i, j int) bool {
	// Equal number of parts, do a string compare per part
	for idx, ival := range p[i].Parts {
		// If we're passed the length of p[j] parts, then we're done
		if idx >= len(p[j].Parts) {
			break
		}

		// Compare the values if they're not equal
		jval := p[j].Parts[idx]
		if ival != jval {
			return ival < jval
		}
	}

	// Equal prefix, take the shorter
	if len(p[i].Parts) != len(p[j].Parts) {
		return len(p[i].Parts) < len(p[j].Parts)
	}

	// Equal, it doesn't matter
	return false
}
