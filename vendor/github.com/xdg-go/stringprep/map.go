// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package stringprep

// Mapping represents a stringprep mapping, from a single rune to zero or more
// runes.
type Mapping map[rune][]rune

// Map maps a rune to a (possibly empty) rune slice via a stringprep Mapping.
// The ok return value is false if the rune was not found.
func (m Mapping) Map(r rune) (replacement []rune, ok bool) {
	rs, ok := m[r]
	if !ok {
		return nil, false
	}
	return rs, true
}
