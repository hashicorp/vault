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
	addrParts := strings.SplitN(addr, "://", 2)
	if len(addrParts) == 1 && addrParts[0] != "" {
		addrParts = []string{"tcp", addrParts[0]}
	}

	switch addrParts[0] {
	case "tcp":
		return ParseTCPAddr(addrParts[1], defaultTCPHost)
	case "unix":
		return parseSimpleProtoAddr("unix", addrParts[1], defaultUnixSocket)
	case "npipe":
		return parseSimpleProtoAddr("npipe", addrParts[1], defaultNamedPipe)
	case "fd":
		return addr, nil
	case "ssh":
		return addr, nil
	default:
		return "", fmt.Errorf("Invalid bind address format: %s", addr)
	}
}

// parseSimpleProtoAddr parses and validates that the specified address is a valid
// socket address for simple protocols like unix and npipe. It returns a formatted
// socket address, either using the address parsed from addr, or the contents of
// defaultAddr if addr is a blank string.
func parseSimpleProtoAddr(proto, addr, defaultAddr string) (string, error) {
	addr = strings.TrimPrefix(addr, proto+"://")
	if strings.Contains(addr, "://") {
		return "", fmt.Errorf("Invalid proto, expected %s: %s", proto, addr)
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
		return "", fmt.Errorf("Invalid proto, expected tcp: %s", tryAddr)
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
		return "", fmt.Errorf("Invalid bind address format: %s", tryAddr)
	}

	if host == "" {
		host = defaultHost
	}
	if port == "" {
		port = defaultPort
	}
	p, err := strconv.Atoi(port)
	if err != nil && p == 0 {
		return "", fmt.Errorf("Invalid bind address format: %s", tryAddr)
	}

	return fmt.Sprintf("tcp://%s%s", net.JoinHostPort(host, port), u.Path), nil
}

// ValidateExtraHost validates that the specified string is a valid extrahost and returns it.
// ExtraHost is in the form of name:ip where the ip has to be a valid ip (IPv4 or IPv6).
func ValidateExtraHost(val string) (string, error) {
	// allow for IPv6 addresses in extra hosts by only splitting on first ":"
	arr := strings.SplitN(val, ":", 2)
	if len(arr) != 2 || len(arr[0]) == 0 {
		return "", fmt.Errorf("bad format for add-host: %q", val)
	}
	// Skip IPaddr validation for "host-gateway" string
	if arr[1] != hostGatewayName {
		if _, err := ValidateIPAddress(arr[1]); err != nil {
			return "", fmt.Errorf("invalid IP address in add-host: %q", arr[1])
		}
	}
	return val, nil
}
