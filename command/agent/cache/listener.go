package cache

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"

	osuser "os/user"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/server"
)

func StartListener(lnConfig *config.Listener, unixSocketsConfig *config.UnixSockets) (net.Listener, *tls.Config, error) {
	addr, ok := lnConfig.Config["address"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("invalid address")
	}

	bindProto := "tcp"
	switch lnConfig.Type {
	case "tcp":
		if addr == "" {
			addr = "127.0.0.1:8007"
		}

		// If they've passed 0.0.0.0, we only want to bind on IPv4
		// rather than golang's dual stack default
		if strings.HasPrefix(addr, "0.0.0.0:") {
			bindProto = "tcp4"
		}

	case "unix":
		addr = "unix://" + addr
	default:
		return nil, nil, fmt.Errorf("invalid listener type: %q", lnConfig.Type)
	}

	var netAddr net.Addr
	switch {
	case strings.HasPrefix(addr, "unix://"):
		netAddr = &net.UnixAddr{
			Name: addr[len("unix://"):],
			Net:  "unix",
		}
	default:
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, nil, err
		}

		nPort, err := strconv.Atoi(port)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid tcp port %q", port)
		}

		ip := net.ParseIP(host)
		if ip == nil {
			return nil, nil, fmt.Errorf("invalid ip address %q", addr)
		}
		netAddr = &net.TCPAddr{
			IP:   ip,
			Port: nPort,
		}
	}

	var ln net.Listener
	var err error
	switch addrType := netAddr.(type) {
	case *net.UnixAddr:
		ln, err = listenSocket(addrType.Name, unixSocketsConfig)
		if err != nil {
			return nil, nil, err
		}

	case *net.TCPAddr:
		ln, err = net.Listen(bindProto, addrType.String())
		if err != nil {
			return nil, nil, err
		}
		ln = &server.TCPKeepAliveListener{ln.(*net.TCPListener)}
	}

	props := map[string]string{"addr": ln.Addr().String()}
	ln, props, _, tlsConf, err := server.ListenerWrapTLS(ln, props, lnConfig.Config, nil)
	if err != nil {
		return nil, nil, err
	}

	return ln, tlsConf, nil
}

func listenSocket(path string, unixSocketsConfig *config.UnixSockets) (net.Listener, error) {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to remove socket file: %v", err)
	}

	ln, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}

	if unixSocketsConfig != nil {
		err = setFilePermissions(path, unixSocketsConfig.User, unixSocketsConfig.Group, unixSocketsConfig.Mode)
		if err != nil {
			return nil, fmt.Errorf("failed to set file system permissions on the socket file: %s", err)
		}
	}

	// Wrap the listener in rmListener so that the Unix domain socket file is
	// removed on close.
	return &rmListener{
		Listener: ln,
		Path:     path,
	}, nil
}

// setFilePermissions handles configuring ownership and permissions
// settings on a given file. All permission/ownership settings are
// optional. If no user or group is specified, the current user/group
// will be used. Mode is optional, and has no default (the operation is
// not performed if absent). User may be specified by name or ID, but
// group may only be specified by ID.
func setFilePermissions(path string, user, group, mode string) error {
	var err error
	uid, gid := os.Getuid(), os.Getgid()

	if user != "" {
		if uid, err = strconv.Atoi(user); err == nil {
			goto GROUP
		}

		// Try looking up the user by name
		u, err := osuser.Lookup(user)
		if err != nil {
			return fmt.Errorf("failed to look up user %q: %v", user, err)
		}
		uid, _ = strconv.Atoi(u.Uid)
	}

GROUP:
	if group != "" {
		if gid, err = strconv.Atoi(group); err != nil {
			return fmt.Errorf("invalid group specified: %v", group)
		}
	}
	if err := os.Chown(path, uid, gid); err != nil {
		return fmt.Errorf("failed setting ownership to %d:%d on %q: %v",
			uid, gid, path, err)
	}

	if mode != "" {
		mode, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			return fmt.Errorf("invalid mode specified: %v", mode)
		}
		if err := os.Chmod(path, os.FileMode(mode)); err != nil {
			return fmt.Errorf("failed setting permissions to %d on %q: %v",
				mode, path, err)
		}
	}

	return nil
}

// rmListener is an implementation of net.Listener that forwards most
// calls to the listener but also removes a file as part of the close. We
// use this to cleanup the unix domain socket on close.
type rmListener struct {
	net.Listener
	Path string
}

func (l *rmListener) Close() error {
	// Close the listener itself
	if err := l.Listener.Close(); err != nil {
		return err
	}

	// Remove the file
	return os.Remove(l.Path)
}
