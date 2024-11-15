// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

// A set of strings
type StringSet map[string]struct{}

func ToSet(values []string) StringSet {
	s := make(StringSet)
	for _, v := range values {
		s[v] = struct{}{}
	}
	return s
}

func (ss StringSet) Add(v string) {
	ss[v] = struct{}{}
}

func (ss StringSet) ToSlice() []string {
	ls := make([]string, len(ss))
	i := 0
	for r := range ss {
		ls[i] = r
		i++
	}
	return ls
}

func (ss StringSet) Includes(v string) bool {
	_, ok := ss[v]
	return ok
}

func (ss StringSet) Update(members ...string) {
	for _, v := range members {
		ss[v] = struct{}{}
	}
}

func (ss StringSet) Union(other StringSet) StringSet {
	un := make(StringSet)
	for v := range ss {
		un[v] = struct{}{}
	}
	for v := range other {
		un[v] = struct{}{}
	}
	return un
}

func (ss StringSet) Intersection(other StringSet) StringSet {
	inter := make(StringSet)

	var s StringSet
	if len(ss) > len(other) {
		s = other
	} else {
		s = ss
	}

	for v := range s {
		if other.Includes(v) {
			inter[v] = struct{}{}
		}
	}
	return inter
}

func (ss StringSet) Sub(other StringSet) StringSet {
	sub := make(StringSet)
	for v := range ss {
		if !other.Includes(v) {
			sub[v] = struct{}{}
		}
	}
	return sub
}

func (ss StringSet) Equals(other StringSet) bool {
	return len(ss.Intersection(other)) == len(ss)
}
