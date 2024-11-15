package proxyproto

import (
	"fmt"
	"net"
	"strings"
)

// PolicyFunc can be used to decide whether to trust the PROXY info from
// upstream. If set, the connecting address is passed in as an argument.
//
// See below for the different policies.
//
// In case an error is returned the connection is denied.
type PolicyFunc func(upstream net.Addr) (Policy, error)

// ConnPolicyFunc can be used to decide whether to trust the PROXY info
// based on connection policy options. If set, the connecting addresses
// (remote and local) are passed in as argument.
//
// See below for the different policies.
//
// In case an error is returned the connection is denied.
type ConnPolicyFunc func(connPolicyOptions ConnPolicyOptions) (Policy, error)

// ConnPolicyOptions contains the remote and local addresses of a connection.
type ConnPolicyOptions struct {
	Upstream   net.Addr
	Downstream net.Addr
}

// Policy defines how a connection with a PROXY header address is treated.
type Policy int

const (
	// USE address from PROXY header
	USE Policy = iota
	// IGNORE address from PROXY header, but accept connection
	IGNORE
	// REJECT connection when PROXY header is sent
	// Note: even though the first read on the connection returns an error if
	// a PROXY header is present, subsequent reads do not. It is the task of
	// the code using the connection to handle that case properly.
	REJECT
	// REQUIRE connection to send PROXY header, reject if not present
	// Note: even though the first read on the connection returns an error if
	// a PROXY header is not present, subsequent reads do not. It is the task
	// of the code using the connection to handle that case properly.
	REQUIRE
	// SKIP accepts a connection without requiring the PROXY header
	// Note: an example usage can be found in the SkipProxyHeaderForCIDR
	// function.
	SKIP
)

// SkipProxyHeaderForCIDR returns a PolicyFunc which can be used to accept a
// connection from a skipHeaderCIDR without requiring a PROXY header, e.g.
// Kubernetes pods local traffic. The def is a policy to use when an upstream
// address doesn't match the skipHeaderCIDR.
func SkipProxyHeaderForCIDR(skipHeaderCIDR *net.IPNet, def Policy) PolicyFunc {
	return func(upstream net.Addr) (Policy, error) {
		ip, err := ipFromAddr(upstream)
		if err != nil {
			return def, err
		}

		if skipHeaderCIDR != nil && skipHeaderCIDR.Contains(ip) {
			return SKIP, nil
		}

		return def, nil
	}
}

// WithPolicy adds given policy to a connection when passed as option to NewConn()
func WithPolicy(p Policy) func(*Conn) {
	return func(c *Conn) {
		c.ProxyHeaderPolicy = p
	}
}

// LaxWhiteListPolicy returns a PolicyFunc which decides whether the
// upstream ip is allowed to send a proxy header based on a list of allowed
// IP addresses and IP ranges. In case upstream IP is not in list the proxy
// header will be ignored. If one of the provided IP addresses or IP ranges
// is invalid it will return an error instead of a PolicyFunc.
func LaxWhiteListPolicy(allowed []string) (PolicyFunc, error) {
	allowFrom, err := parse(allowed)
	if err != nil {
		return nil, err
	}

	return whitelistPolicy(allowFrom, IGNORE), nil
}

// MustLaxWhiteListPolicy returns a LaxWhiteListPolicy but will panic if one
// of the provided IP addresses or IP ranges is invalid.
func MustLaxWhiteListPolicy(allowed []string) PolicyFunc {
	pfunc, err := LaxWhiteListPolicy(allowed)
	if err != nil {
		panic(err)
	}

	return pfunc
}

// StrictWhiteListPolicy returns a PolicyFunc which decides whether the
// upstream ip is allowed to send a proxy header based on a list of allowed
// IP addresses and IP ranges. In case upstream IP is not in list reading on
// the connection will be refused on the first read. Please note: subsequent
// reads do not error. It is the task of the code using the connection to
// handle that case properly. If one of the provided IP addresses or IP
// ranges is invalid it will return an error instead of a PolicyFunc.
func StrictWhiteListPolicy(allowed []string) (PolicyFunc, error) {
	allowFrom, err := parse(allowed)
	if err != nil {
		return nil, err
	}

	return whitelistPolicy(allowFrom, REJECT), nil
}

// MustStrictWhiteListPolicy returns a StrictWhiteListPolicy but will panic
// if one of the provided IP addresses or IP ranges is invalid.
func MustStrictWhiteListPolicy(allowed []string) PolicyFunc {
	pfunc, err := StrictWhiteListPolicy(allowed)
	if err != nil {
		panic(err)
	}

	return pfunc
}

func whitelistPolicy(allowed []func(net.IP) bool, def Policy) PolicyFunc {
	return func(upstream net.Addr) (Policy, error) {
		upstreamIP, err := ipFromAddr(upstream)
		if err != nil {
			// something is wrong with the source IP, better reject the connection
			return REJECT, err
		}

		for _, allowFrom := range allowed {
			if allowFrom(upstreamIP) {
				return USE, nil
			}
		}

		return def, nil
	}
}

func parse(allowed []string) ([]func(net.IP) bool, error) {
	a := make([]func(net.IP) bool, len(allowed))
	for i, allowFrom := range allowed {
		if strings.LastIndex(allowFrom, "/") > 0 {
			_, ipRange, err := net.ParseCIDR(allowFrom)
			if err != nil {
				return nil, fmt.Errorf("proxyproto: given string %q is not a valid IP range: %v", allowFrom, err)
			}

			a[i] = ipRange.Contains
		} else {
			allowed := net.ParseIP(allowFrom)
			if allowed == nil {
				return nil, fmt.Errorf("proxyproto: given string %q is not a valid IP address", allowFrom)
			}

			a[i] = allowed.Equal
		}
	}

	return a, nil
}

func ipFromAddr(upstream net.Addr) (net.IP, error) {
	upstreamString, _, err := net.SplitHostPort(upstream.String())
	if err != nil {
		return nil, err
	}

	upstreamIP := net.ParseIP(upstreamString)
	if nil == upstreamIP {
		return nil, fmt.Errorf("proxyproto: invalid IP address")
	}

	return upstreamIP, nil
}

// IgnoreProxyHeaderNotOnInterface retuns a ConnPolicyFunc which can be used to
// decide whether to use or ignore PROXY headers depending on the connection
// being made on a specific interface. This policy can be used when the server
// is bound to multiple interfaces but wants to allow on only one interface.
func IgnoreProxyHeaderNotOnInterface(allowedIP net.IP) ConnPolicyFunc {
	return func(connOpts ConnPolicyOptions) (Policy, error) {
		ip, err := ipFromAddr(connOpts.Downstream)
		if err != nil {
			return REJECT, err
		}

		if allowedIP.Equal(ip) {
			return USE, nil
		}

		return IGNORE, nil
	}
}
