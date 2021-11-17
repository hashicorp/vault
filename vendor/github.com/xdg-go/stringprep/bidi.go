// Copyright 2018 by David A. Golden. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package stringprep

var errHasLCat = "BiDi string can't have runes from category L"
var errFirstRune = "BiDi string first rune must have category R or AL"
var errLastRune = "BiDi string last rune must have category R or AL"

// Check for prohibited characters from table C.8
func checkBiDiProhibitedRune(s string) error {
	for _, r := range s {
		if TableC8.Contains(r) {
			return Error{Msg: errProhibited, Rune: r}
		}
	}
	return nil
}

// Check for LCat characters from table D.2
func checkBiDiLCat(s string) error {
	for _, r := range s {
		if TableD2.Contains(r) {
			return Error{Msg: errHasLCat, Rune: r}
		}
	}
	return nil
}

// Check first and last characters are in table D.1; requires non-empty string
func checkBadFirstAndLastRandALCat(s string) error {
	rs := []rune(s)
	if !TableD1.Contains(rs[0]) {
		return Error{Msg: errFirstRune, Rune: rs[0]}
	}
	n := len(rs) - 1
	if !TableD1.Contains(rs[n]) {
		return Error{Msg: errLastRune, Rune: rs[n]}
	}
	return nil
}

// Look for RandALCat characters from table D.1
func hasBiDiRandALCat(s string) bool {
	for _, r := range s {
		if TableD1.Contains(r) {
			return true
		}
	}
	return false
}

// Check that BiDi rules are satisfied ; let empty string pass this rule
func passesBiDiRules(s string) error {
	if len(s) == 0 {
		return nil
	}
	if err := checkBiDiProhibitedRune(s); err != nil {
		return err
	}
	if hasBiDiRandALCat(s) {
		if err := checkBiDiLCat(s); err != nil {
			return err
		}
		if err := checkBadFirstAndLastRandALCat(s); err != nil {
			return err
		}
	}
	return nil
}
