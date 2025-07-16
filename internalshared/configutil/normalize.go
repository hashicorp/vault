// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"net"
	"net/url"
	"strings"
)

// NormalizeAddr takes a string of a Host, Host:Port, URL, or Destination
// Address and returns a copy where any IP addresses have been normalized to be
// conformant with RFC 5942 §4. If the input string does not match any of the
// supported syntaxes, or the "host" section is not an IP address, the input
// will be returned unchanged. Supported syntaxes are:
//
//	Host                 host                                or [host]
//	Host:Port            host:port                           or [host]:port
//	URL                  scheme://user@host/path?query#frag  or scheme://user@[host]/path?query#frag
//	Destination Address  user@host:port                      or user@[host]:port
//
// See:
//
//	https://rfc-editor.org/rfc/rfc3986.html
//	https://rfc-editor.org/rfc/rfc5942.html
//	https://rfc-editor.org/rfc/rfc5952.html
func NormalizeAddr(addr string) string {
	if addr == "" {
		return ""
	}

	// Host
	ip := net.ParseIP(addr)
	if ip != nil {
		// net.IP.String() is RFC 5942 §4 compliant
		return ip.String()
	}

	// [Host]
	if strings.HasPrefix(addr, "[") && strings.HasSuffix(addr, "]") {
		if len(addr) < 3 {
			return addr
		}

		// If we've been given a bracketed IP address, return the address
		// normalized without brackets.
		ip := net.ParseIP(addr[1 : len(addr)-1])
		if ip != nil {
			return ip.String()
		}

		// Our input is not a valid schema.
		return addr
	}

	// Host:Port
	host, port, err := net.SplitHostPort(addr)
	if err == nil {
		ip := net.ParseIP(host)
		if ip == nil {
			// Our host isn't an IP address so we can return it unchanged
			return addr
		}

		// net.JoinHostPort handles bracketing for RFC 5952 §6
		return net.JoinHostPort(ip.String(), port)
	}

	// URL
	u, err := url.Parse(addr)
	if err == nil {
		uhost := u.Hostname()
		ip := net.ParseIP(uhost)
		if ip == nil {
			// Our URL doesn't contain an IP address so we can return our input unchanged.
			return addr
		} else {
			uhost = ip.String()
		}

		if uport := u.Port(); uport != "" {
			uhost = net.JoinHostPort(uhost, uport)
		} else {
			if !strings.HasPrefix(uhost, "[") && !strings.HasSuffix(uhost, "]") {
				// Ensure the IPv6 URL host is bracketed post-normalization.
				// When*url.URL.String() reassembles the URL it will not consider
				// whether or not the  *url.URL.Host is RFC 5952 §6 and RFC 3986 §3.2.2
				// conformant.
				uhost = "[" + uhost + "]"
			}
		}
		u.Host = uhost

		return u.String()
	}

	// Destination Address
	if idx := strings.LastIndex(addr, "@"); idx > 0 {
		if idx+1 > len(addr) {
			return addr
		}

		return addr[:idx+1] + NormalizeAddr(addr[idx+1:])
	}

	// Our input did not match our supported schemas. Return it unchanged.
	return addr
}
