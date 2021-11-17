// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Based on github.com/aws/aws-sdk-go by Amazon.com, Inc. with code from:
// - github.com/aws/aws-sdk-go/blob/v1.34.28/aws/signer/v4/header_rules.go
// - github.com/aws/aws-sdk-go/blob/v1.34.28/internal/strings/strings.go
// See THIRD-PARTY-NOTICES for original license terms

package awsv4

import (
	"strings"
)

// validator houses a set of rule needed for validation of a
// string value
type rules []rule

// rule interface allows for more flexible rules and just simply
// checks whether or not a value adheres to that rule
type rule interface {
	IsValid(value string) bool
}

// IsValid will iterate through all rules and see if any rules
// apply to the value and supports nested rules
func (r rules) IsValid(value string) bool {
	for _, rule := range r {
		if rule.IsValid(value) {
			return true
		}
	}
	return false
}

// mapRule generic rule for maps
type mapRule map[string]struct{}

// IsValid for the map rule satisfies whether it exists in the map
func (m mapRule) IsValid(value string) bool {
	_, ok := m[value]
	return ok
}

// whitelist is a generic rule for whitelisting
type whitelist struct {
	rule
}

// IsValid for whitelist checks if the value is within the whitelist
func (w whitelist) IsValid(value string) bool {
	return w.rule.IsValid(value)
}

// blacklist is a generic rule for blacklisting
type blacklist struct {
	rule
}

// IsValid for whitelist checks if the value is within the whitelist
func (b blacklist) IsValid(value string) bool {
	return !b.rule.IsValid(value)
}

type patterns []string

// hasPrefixFold tests whether the string s begins with prefix, interpreted as UTF-8 strings,
// under Unicode case-folding.
func hasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}

// IsValid for patterns checks each pattern and returns if a match has
// been found
func (p patterns) IsValid(value string) bool {
	for _, pattern := range p {
		if hasPrefixFold(value, pattern) {
			return true
		}
	}
	return false
}

// inclusiveRules rules allow for rules to depend on one another
type inclusiveRules []rule

// IsValid will return true if all rules are true
func (r inclusiveRules) IsValid(value string) bool {
	for _, rule := range r {
		if !rule.IsValid(value) {
			return false
		}
	}
	return true
}
