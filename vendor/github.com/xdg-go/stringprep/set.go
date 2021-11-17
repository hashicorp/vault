// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package stringprep

import "sort"

// RuneRange represents a close-ended range of runes: [N,M].  For a range
// consisting of a single rune, N and M will be equal.
type RuneRange [2]rune

// Contains returns true if a rune is within the bounds of the RuneRange.
func (rr RuneRange) Contains(r rune) bool {
	return rr[0] <= r && r <= rr[1]
}

func (rr RuneRange) isAbove(r rune) bool {
	return r <= rr[0]
}

// Set represents a stringprep data table used to identify runes of a
// particular type.
type Set []RuneRange

// Contains returns true if a rune is within any of the RuneRanges in the
// Set.
func (s Set) Contains(r rune) bool {
	i := sort.Search(len(s), func(i int) bool { return s[i].Contains(r) || s[i].isAbove(r) })
	if i < len(s) && s[i].Contains(r) {
		return true
	}
	return false
}
