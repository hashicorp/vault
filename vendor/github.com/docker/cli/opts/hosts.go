package opts

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

const (
	// defaultHTTPPort Default HTTP Port used if only the protocol is provided to -H flag e.g. dockerd -H tcp://
	// These are the IANA registered port numbers for use with Docker
	// see http://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml?search=docker
	defaultHTTPPort = "2375" // Default HTTP Port
	// defaultTLSHTTPPort Default HTTP Port used when TLS enabled
	defaultTLSHTTPPort = "2376" // Default TLS encrypted HTTP Port
	// defaultUnixSocket Path for the unix socket.
	// Docker daemon by default always listens on the default unix socket
	defaultUnixSocket = "/var/run/docker.sock"
	// defaultTCPHost constant defines the default host string used by docker on Windows
	defaultTCPHost = "tcp://" + defaultHTTPHost + ":" + defaultHTTPPort
	// DefaultTLSHost constant defines the default host string used by docker for TLS sockets
	defaultTLSHost = "tcp://" + defaultHTTPHost + ":" + defaultTLSHTTPPort
	// DefaultNamedPipe defines the default named pipe used by docker on Windows
	defaultNamedPipe = `//./pipe/docker_engine`
	// hostGatewayName defines a special string which users can append to --add-host
	// to add an extra entry in /etc/hosts that maps host.docker.internal to the host IP
	// TODO Consider moving the hostGatewayName constant defined in docker at
	// github.com/docker/docker/daemon/network/constants.go outside of the "daemon"
	// package, so that the CLI can consume it.
	hostGatewayName = "host-gateway"
)

// ValidateHost validates that the specified string is a valid host and returns it.
//
// TODO(thaJeztah): ValidateHost appears to be unused; deprecate it.
func ValidateHost(val string) (string, error) {
	host := strings.TrimSpace(val)
	// The empty string means default and is not handled by parseDockerDaemonHost
	if host != "" {
		_, err := parseDockerDaemonHost(host)
		if err != nil {
			return val, err
		}
	}
	// Note: unlike most flag validators, we don't return the mutated value here
	//       we need to know what the user entered later (using ParseHost) to adjust for TLS
	return val, nil
}

// ParseHost and set defaults for a Daemon host string
func ParseHost(defaultToTLS bool, val string) (string, error) {
	host := strings.TrimSpace(val)
	if host == "" {
		if defaultToTLS {
			host = defaultTLSHost
		} else {
			host = defaultHost
		}
	} else {
		var err error
		host, err = parseDockerDaemonHost(host)
		if err != nil {
			return val, err
		}
	}
	return host, nil
}

// parseDockerDaemonHost parses the specified address and returns an address that will be used as the host.
// Depending of the address specified, this may return one of the global Default* strings defined in hosts.go.
func parseDockerDaemonHost(addr string) (string, error) {
	proto, host, hasProto := strings.Cut(addr, "://")
	if !hasProto && proto != "" {
		host = proto
		proto = "tcp"
	}

	switch proto {
	case "tcp":
		return ParseTCPAddr(host, defaultTCPHost)
	case "unix":
		return parseSimpleProtoAddr(proto, host, defaultUnixSocket)
	case "npipe":
		return parseSimpleProtoAddr(proto, host, defaultNamedPipe)
	case "fd":
		return addr, nil
	case "ssh":
		return addr, nil
	default:
		return "", fmt.Errorf("invalid bind address format: %s", addr)
	}
}

// parseSimpleProtoAddr parses and validates that the specified address is a valid
// socket address for simple protocols like unix and npipe. It returns a formatted
// socket address, either using the address parsed from addr, or the contents of
// defaultAddr if addr is a blank string.
func parseSimpleProtoAddr(proto, addr, defaultAddr string) (string, error) {
	addr = strings.TrimPrefix(addr, proto+"://")
	if strings.Contains(addr, "://") {
		return "", fmt.Errorf("invalid proto, expected %s: %s", proto, addr)
	}
	if addr == "" {
		addr = defaultAddr
	}
	return fmt.Sprintf("%s://%s", proto, addr), nil
}

// ParseTCPAddr parses and validates that the specified address is a valid TCP
// address. It returns a formatted TCP address, either using the address parsed
// from tryAddr, or the contents of defaultAddr if tryAddr is a blank string.
// tryAddr is expected to have already been Trim()'d
// defaultAddr must be in the full `tcp://host:port` form
func ParseTCPAddr(tryAddr string, defaultAddr string) (string, error) {
	if tryAddr == "" || tryAddr == "tcp://" {
		return defaultAddr, nil
	}
	addr := strings.TrimPrefix(tryAddr, "tcp://")
	if strings.Contains(addr, "://") || addr == "" {
		return "", fmt.Errorf("invalid proto, expected tcp: %s", tryAddr)
	}

	defaultAddr = strings.TrimPrefix(defaultAddr, "tcp://")
	defaultHost, defaultPort, err := net.SplitHostPort(defaultAddr)
	if err != nil {
		return "", err
	}
	// url.Parse fails for trailing colon on IPv6 brackets on Go 1.5, but
	// not 1.4. See https://github.com/golang/go/issues/12200 and
	// https://github.com/golang/go/issues/6530.
	if strings.HasSuffix(addr, "]:") {
		addr += defaultPort
	}

	u, err := url.Parse("tcp://" + addr)
	if err != nil {
		return "", err
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		// try port addition once
		host, port, err = net.SplitHostPort(net.JoinHostPort(u.Host, defaultPort))
	}
	if err != nil {
		return "", fmt.Errorf("invalid bind address format: %s", tryAddr)
	}

	if host == "" {
		host = defaultHost
	}
	if port == "" {
		port = defaultPort
	}
	p, err := strconv.Atoi(port)
	if err != nil && p == 0 {
		return "", fmt.Errorf("invalid bind address format: %s", tryAddr)
	}

	return fmt.Sprintf("tcp://%s%s", net.JoinHostPort(host, port), u.Path), nil
}

// ValidateExtraHost validates that the specified string is a valid extrahost and
// returns it. ExtraHost is in the form of name:ip or name=ip, where the ip has
// to be a valid ip (IPv4 or IPv6). The address may be enclosed in square
// brackets.
//
// For example:
//
//	my-hostname:127.0.0.1
//	my-hostname:::1
//	my-hostname=::1
//	my-hostname:[::1]
//
// For compatibility with the API server, this function normalises the given
// argument to use the ':' separator and strip square brackets enclosing the
// address.
func ValidateExtraHost(val string) (string, error) {
	k, v, ok := strings.Cut(val, "=")
	if !ok {
		// allow for IPv6 addresses in extra hosts by only splitting on first ":"
		k, v, ok = strings.Cut(val, ":")
	}
	// Check that a hostname was given, and that it doesn't contain a ":". (Colon
	// isn't allowed in a hostname, along with many other characters. It's
	// special-cased here because the API server doesn't know about '=' separators in
	// '--add-host'. So, it'll split at the first colon and generate a strange error
	// message.)
	if !ok || k == "" || strings.Contains(k, ":") {
		return "", fmt.Errorf("bad format for add-host: %q", val)
	}
	// Skip IPaddr validation for "host-gateway" string
	if v != hostGatewayName {
		// If the address is enclosed in square brackets, extract it (for IPv6, but
		// permit it for IPv4 as well; we don't know the address family here, but it's
		// unambiguous).
		if len(v) > 2 && v[0] == '[' && v[len(v)-1] == ']' {
			v = v[1 : len(v)-1]
		}
		// ValidateIPAddress returns the address in canonical form (for example,
		// 0:0:0:0:0:0:0:1 -> ::1). But, stick with the original form, to avoid
		// surprising a user who's expecting to see the address they supplied in the
		// output of 'docker inspect' or '/etc/hosts'.
		if _, err := ValidateIPAddress(v); err != nil {
			return "", fmt.Errorf("invalid IP address in add-host: %q", v)
		}
	}
	// This result is passed directly to the API, the daemon doesn't accept the '='
	// separator or an address enclosed in brackets. So, construct something it can
	// understand.
	return k + ":" + v, nil
}
