package parseip

import (
	"net"
)

// In Go 1.17 the behaviour of net.ParseIP and net.ParseCIDR changed
// (https://golang.org/doc/go1.17#net) so that leading zeros in the input results
// in an error.  This package contains helpers that strip leading zeroes so as
// to avoid those errors.

// TrimLeadingZeroes returns its input trimmed of any leading zeroes.
func TrimLeadingZeroes(s string) string {
	for i, r := range s {
		if r == '0' {
			continue
		}
		return s[i:]
	}
	return ""
}

// ParseIP strips any leading zeroes and then calls net.ParseIP.  This is not
// exactly the same as the net.ParseIP behaviour prior to Go 1.17, but it at least
// won't return any errors if there are leading zeroes.
func ParseIP(s string) net.IP {
	return net.ParseIP(TrimLeadingZeroes(s))
}

// ParseCIDR strips any leading zeroes and then calls net.ParseCIDR.  This is not
// exactly the same as the net.ParseCIDR behaviour prior to Go 1.17, but it at least
// won't return any errors if there are leading zeroes.
func ParseCIDR(s string) (net.IP, *net.IPNet, error) {
	return net.ParseCIDR(TrimLeadingZeroes(s))
}
