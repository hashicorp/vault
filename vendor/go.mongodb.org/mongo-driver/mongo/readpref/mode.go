// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package readpref

import (
	"fmt"
	"strings"
)

// Mode indicates the user's preference on reads.
type Mode uint8

// Mode constants
const (
	_ Mode = iota
	// PrimaryMode indicates that only a primary is
	// considered for reading. This is the default
	// mode.
	PrimaryMode
	// PrimaryPreferredMode indicates that if a primary
	// is available, use it; otherwise, eligible
	// secondaries will be considered.
	PrimaryPreferredMode
	// SecondaryMode indicates that only secondaries
	// should be considered.
	SecondaryMode
	// SecondaryPreferredMode indicates that only secondaries
	// should be considered when one is available. If none
	// are available, then a primary will be considered.
	SecondaryPreferredMode
	// NearestMode indicates that all primaries and secondaries
	// will be considered.
	NearestMode
)

// ModeFromString returns a mode corresponding to
// mode.
func ModeFromString(mode string) (Mode, error) {
	switch strings.ToLower(mode) {
	case "primary":
		return PrimaryMode, nil
	case "primarypreferred":
		return PrimaryPreferredMode, nil
	case "secondary":
		return SecondaryMode, nil
	case "secondarypreferred":
		return SecondaryPreferredMode, nil
	case "nearest":
		return NearestMode, nil
	}
	return Mode(0), fmt.Errorf("unknown read preference %v", mode)
}

// String returns the string representation of mode.
func (mode Mode) String() string {
	switch mode {
	case PrimaryMode:
		return "primary"
	case PrimaryPreferredMode:
		return "primaryPreferred"
	case SecondaryMode:
		return "secondary"
	case SecondaryPreferredMode:
		return "secondaryPreferred"
	case NearestMode:
		return "nearest"
	default:
		return "unknown"
	}
}
