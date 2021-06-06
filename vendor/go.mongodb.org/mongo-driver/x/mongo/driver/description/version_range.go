// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

import "fmt"

// VersionRange represents a range of versions.
type VersionRange struct {
	Min int32
	Max int32
}

// NewVersionRange creates a new VersionRange given a min and a max.
func NewVersionRange(min, max int32) VersionRange {
	return VersionRange{Min: min, Max: max}
}

// Includes returns a bool indicating whether the supplied integer is included
// in the range.
func (vr VersionRange) Includes(v int32) bool {
	return v >= vr.Min && v <= vr.Max
}

// String implements the fmt.Stringer interface.
func (vr VersionRange) String() string {
	return fmt.Sprintf("[%d, %d]", vr.Min, vr.Max)
}
