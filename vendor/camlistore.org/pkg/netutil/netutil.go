/*
Copyright 2014 The Camlistore Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package netutil

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

// AwaitReachable tries to make a TCP connection to addr regularly.
// It returns an error if it's unable to make a connection before maxWait.
func AwaitReachable(addr string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		c, err := net.Dial("tcp", addr)
		if err == nil {
			c.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("%v unreachable for %v", addr, maxWait)
}

// HostPort takes a urlStr string URL, and returns a host:port string suitable
// to passing to net.Dial, with the port set as the scheme's default port if
// absent.
func HostPort(urlStr string) (string, error) {
	// TODO: rename this function to URLHostPort instead, like
	// ListenHostPort below.
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("could not parse %q as a url: %v", urlStr, err)
	}
	if u.Scheme == "" {
		return "", fmt.Errorf("url %q has no scheme", urlStr)
	}
	hostPort := u.Host
	if hostPort == "" || strings.HasPrefix(hostPort, ":") {
		return "", fmt.Errorf("url %q has no host", urlStr)
	}
	idx := strings.Index(hostPort, "]")
	if idx == -1 {
		idx = 0
	}
	if !strings.Contains(hostPort[idx:], ":") {
		if u.Scheme == "https" {
			hostPort += ":443"
		} else {
			hostPort += ":80"
		}
	}
	return hostPort, nil
}

// ListenHostPort maps a listen address into a host:port string.
// If the host part in listenAddr is empty or 0.0.0.0, localhost
// is used instead.
func ListenHostPort(listenAddr string) (string, error) {
	hp := listenAddr
	if strings.HasPrefix(hp, ":") {
		hp = "localhost" + hp
	} else if strings.HasPrefix(hp, "0.0.0.0:") {
		hp = "localhost:" + hp[len("0.0.0.0:"):]
	}
	if _, _, err := net.SplitHostPort(hp); err != nil {
		return "", err
	}
	return hp, nil
}

// ListenOnLocalRandomPort returns a TCP listener on a random
// localhost port.
func ListenOnLocalRandomPort() (net.Listener, error) {
	ip, err := Localhost()
	if err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", &net.TCPAddr{IP: ip, Port: 0})
}

// Localhost returns the first address found when
// doing a lookup of "localhost". If not successful,
// it looks for an ip on the loopback interfaces.
func Localhost() (net.IP, error) {
	if ip := localhostLookup(); ip != nil {
		return ip, nil
	}
	if ip := loopbackIP(); ip != nil {
		return ip, nil
	}
	return nil, errors.New("No loopback ip found.")
}

// localhostLookup looks for a loopback IP by resolving localhost.
func localhostLookup() net.IP {
	if ips, err := net.LookupIP("localhost"); err == nil && len(ips) > 0 {
		return ips[0]
	}
	return nil
}

// loopbackIP returns the first loopback IP address sniffing network
// interfaces or nil if none is found.
func loopbackIP() net.IP {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, inf := range interfaces {
		const flagUpLoopback = net.FlagUp | net.FlagLoopback
		if inf.Flags&flagUpLoopback == flagUpLoopback {
			addrs, _ := inf.Addrs()
			for _, addr := range addrs {
				ip, _, err := net.ParseCIDR(addr.String())
				if err == nil && ip.IsLoopback() {
					return ip
				}
			}
		}
	}
	return nil
}
