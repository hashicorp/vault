// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package parseip

import (
	"strings"

	"k8s.io/utils/net"
)

// In Go 1.17 the behaviour of net.ParseIP and net.ParseCIDR changed
// (https://golang.org/doc/go1.17#net) so that leading zeros in the input results
// in an error.  This package contains helpers that strip leading zeroes so as
// to avoid those errors.

// You should probably not be using anything here unless you've found a new place
// where IPs/CIDRs are read from storage and re-parsed.

// trimLeadingZeroes returns its input trimmed of any leading zeroes.
func trimLeadingZeroes(s string) string {
	for i, r := range s {
		if r == '0' {
			continue
		}
		return s[i:]
	}
	return ""
}

// trimLeadingZeroesIPv4 takes an IPv4 string and returns the input
// trimmed of any excess leading zeroes in each octet.
func trimLeadingZeroesIPv4(s string) string {
	if len(s) == 0 {
		return s
	}

	pieces := strings.Split(s, ".")
	var sb strings.Builder
	for i, piece := range pieces {
		trimmed := trimLeadingZeroes(piece)
		if trimmed == "" && len(piece) > 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteString(trimmed)
		}
		if i != len(pieces)-1 {
			sb.WriteByte('.')
		}
	}
	return sb.String()
}

// trimLeadingZeroesIP does the same work as trimLeadingZeroesIPv4 but also accepts
// an IPv6 address that may contain an IPv4 address representation. Only decimal
// IPv4 addresses get zero-stripped.
func trimLeadingZeroesIP(s string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ':' && net.ParseIPSloppy(s[i+1:]) != nil {
			return s[:i+1] + trimLeadingZeroesIPv4(s[i+1:])
		}
	}
	return trimLeadingZeroesIPv4(s)
}

// TrimLeadingZeroesCIDR does the same thing as trimLeadingZeroesIP but expects
// a CIDR address as input.  If the input isn't a valid CIDR address, it is returned
// unchanged.
func TrimLeadingZeroesCIDR(s string) string {
	pieces := strings.Split(s, "/")
	if len(pieces) != 2 {
		return s
	}
	pieces[0] = trimLeadingZeroesIP(pieces[0])
	return strings.Join(pieces, "/")
}
