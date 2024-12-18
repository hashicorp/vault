package configutil

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// NormalizeAddr takes an address as a string and returns a normalized copy.
// If the addr is a URL, IP Address, or host:port address that includes an IPv6
// address, the normalized copy will be conformant with RFC-5942 §4
// See: https://rfc-editor.org/rfc/rfc5952.html
func NormalizeAddr(u string) string {
	if u == "" {
		return ""
	}

	var ip net.IP
	var port string
	bracketedIPv6 := false

	// Try parsing it as a URL
	pu, err := url.Parse(u)
	if err == nil {
		// We've been given something that appears to be a URL. See if the hostname
		// is an IP address
		ip = net.ParseIP(pu.Hostname())
	} else {
		// We haven't been given a URL. Try and parse it as an IP address
		ip = net.ParseIP(u)
		if ip == nil {
			// We haven't been given a URL or IP address, try parsing an IP:Port
			// combination.
			idx := strings.LastIndex(u, ":")
			if idx > 0 {
				// We've perhaps received an IP:Port address
				addr := u[:idx]
				port = u[idx+1:]
				if strings.HasPrefix(addr, "[") && strings.HasSuffix(addr, "]") {
					addr = strings.TrimPrefix(strings.TrimSuffix(addr, "]"), "[")
					bracketedIPv6 = true
				}
				ip = net.ParseIP(addr)
			}
		}
	}

	// If our IP is nil whatever was passed in does not contain an IP address.
	if ip == nil {
		return u
	}

	if v4 := ip.To4(); v4 != nil {
		// We don't need to normalize IPv4 addresses.
		if pu != nil && port == "" {
			return pu.String()
		}

		if port != "" {
			// Return the address:port
			return fmt.Sprintf("%s:%s", v4.String(), port)
		}

		// Return the ip addres
		return v4.String()
	}

	if v6 := ip.To16(); v6 != nil {
		// net.IP String() will return IPv6 RFC-5952 conformant addresses.

		if pu != nil {
			// Return the URL in conformant fashion
			if port := pu.Port(); port != "" {
				pu.Host = fmt.Sprintf("[%s]:%s", v6.String(), port)
			} else {
				pu.Host = fmt.Sprintf("[%s]", v6.String())
			}
			return pu.String()
		}

		// Handle IP:Port addresses
		if port != "" {
			// Return the address:port or [address]:port
			if bracketedIPv6 {
				return fmt.Sprintf("[%s]:%s", v6.String(), port)
			} else {
				return fmt.Sprintf("%s:%s", v6.String(), port)
			}
		}

		// Handle just an IP address
		return v6.String()
	}

	// It shouldn't be possible to get to this point. If we somehow we manage
	// to, return the string unchanged.
	return u
}
