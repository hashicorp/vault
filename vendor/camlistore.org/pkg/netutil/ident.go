/*
Copyright 2011 Google Inc.

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

// Package netutil identifies the system userid responsible for
// localhost TCP connections.
package netutil // import "camlistore.org/pkg/netutil"

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	ErrNotFound      = errors.New("netutil: connection not found")
	ErrUnsupportedOS = errors.New("netutil: not implemented on this operating system")
)

// ConnUserid returns the uid that owns the given localhost connection.
// The returned error is ErrNotFound if the connection wasn't found.
func ConnUserid(conn net.Conn) (uid int, err error) {
	return AddrPairUserid(conn.LocalAddr(), conn.RemoteAddr())
}

// HostPortToIP parses a host:port to a TCPAddr without resolving names.
// If given a context IP, it will resolve localhost to match the context's IP family.
func HostPortToIP(hostport string, ctx *net.TCPAddr) (hostaddr *net.TCPAddr, err error) {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return nil, err
	}
	iport, err := strconv.Atoi(port)
	if err != nil || iport < 0 || iport > 0xFFFF {
		return nil, fmt.Errorf("invalid port %d", iport)
	}
	var addr net.IP
	if ctx != nil && host == "localhost" {
		if ctx.IP.To4() != nil {
			addr = net.IPv4(127, 0, 0, 1)
		} else {
			addr = net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
		}
	} else if addr = net.ParseIP(host); addr == nil {
		return nil, fmt.Errorf("could not parse IP %s", host)
	}

	return &net.TCPAddr{IP: addr, Port: iport}, nil
}

// AddrPairUserid returns the local userid who owns the TCP connection
// given by the local and remote ip:port (lipport and ripport,
// respectively).  Returns ErrNotFound for the error if the TCP connection
// isn't found.
func AddrPairUserid(local, remote net.Addr) (uid int, err error) {
	lAddr, lOk := local.(*net.TCPAddr)
	rAddr, rOk := remote.(*net.TCPAddr)
	if !(lOk && rOk) {
		return -1, fmt.Errorf("netutil: Could not convert Addr to TCPAddr.")
	}

	localv4 := (lAddr.IP.To4() != nil)
	remotev4 := (rAddr.IP.To4() != nil)
	if localv4 != remotev4 {
		return -1, fmt.Errorf("netutil: address pairs of different families; localv4=%v, remotev4=%v",
			localv4, remotev4)
	}

	switch runtime.GOOS {
	case "darwin":
		return uidFromLsof(lAddr.IP, lAddr.Port, rAddr.IP, rAddr.Port)
	case "freebsd":
		return uidFromSockstat(lAddr.IP, lAddr.Port, rAddr.IP, rAddr.Port)
	case "linux":
		file := "/proc/net/tcp"
		if !localv4 {
			file = "/proc/net/tcp6"
		}
		f, err := os.Open(file)
		if err != nil {
			return -1, fmt.Errorf("Error opening %s: %v", file, err)
		}
		defer f.Close()
		return uidFromProcReader(lAddr.IP, lAddr.Port, rAddr.IP, rAddr.Port, f)
	}
	return 0, ErrUnsupportedOS
}

func toLinuxIPv4Order(b []byte) []byte {
	binary.BigEndian.PutUint32(b, binary.LittleEndian.Uint32(b))
	return b
}

func toLinuxIPv6Order(b []byte) []byte {
	for i := 0; i < 16; i += 4 {
		sb := b[i : i+4]
		binary.BigEndian.PutUint32(sb, binary.LittleEndian.Uint32(sb))
	}
	return b
}

type maybeBrackets net.IP

func (p maybeBrackets) String() string {
	s := net.IP(p).String()
	if strings.Contains(s, ":") {
		return "[" + s + "]"
	}
	return s
}

// Changed by tests.
var uidFromUsername = uidFromUsernameFn

func uidFromUsernameFn(username string) (uid int, err error) {
	if uid := os.Getuid(); uid != 0 && username == os.Getenv("USER") {
		return uid, nil
	}
	u, err := user.Lookup(username)
	if err == nil {
		uid, err := strconv.Atoi(u.Uid)
		return uid, err
	}
	return 0, err
}

func uidFromLsof(lip net.IP, lport int, rip net.IP, rport int) (uid int, err error) {
	seek := fmt.Sprintf("%s:%d->%s:%d", maybeBrackets(lip), lport, maybeBrackets(rip), rport)
	seekb := []byte(seek)
	if _, err = exec.LookPath("lsof"); err != nil {
		return
	}
	cmd := exec.Command("lsof",
		"-b", // avoid system calls that could block
		"-w", // and don't warn about cases where -b fails
		"-n", // don't resolve network names
		"-P", // don't resolve network ports,
		// TODO(bradfitz): pass down the uid we care about, then do: ?
		//"-a",  // AND the following together:
		// "-u", strconv.Itoa(uid)  // just this uid
		"-itcp") // we only care about TCP connections
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	defer cmd.Wait()
	defer stdout.Close()
	err = cmd.Start()
	if err != nil {
		return
	}
	defer cmd.Process.Kill()
	br := bufio.NewReader(stdout)
	for {
		line, err := br.ReadSlice('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return -1, err
		}
		if !bytes.Contains(line, seekb) {
			continue
		}
		// SystemUIS   276 bradfitz   15u  IPv4 0xffffff801a7c74e0      0t0  TCP 127.0.0.1:56718->127.0.0.1:5204 (ESTABLISHED)
		f := bytes.Fields(line)
		if len(f) < 8 {
			continue
		}
		username := string(f[2])
		return uidFromUsername(username)
	}
	return -1, ErrNotFound

}

func uidFromSockstat(lip net.IP, lport int, rip net.IP, rport int) (int, error) {
	cmd := exec.Command("sockstat", "-Ptcp")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}
	defer cmd.Wait()
	defer stdout.Close()
	err = cmd.Start()
	if err != nil {
		return -1, err
	}
	defer cmd.Process.Kill()

	return uidFromSockstatReader(lip, lport, rip, rport, stdout)
}

func uidFromSockstatReader(lip net.IP, lport int, rip net.IP, rport int, r io.Reader) (int, error) {
	pat, err := regexp.Compile(fmt.Sprintf(`^([^ ]+).*%s:%d *%s:%d$`,
		lip.String(), lport, rip.String(), rport))
	if err != nil {
		return -1, err
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l := scanner.Text()
		m := pat.FindStringSubmatch(l)
		if len(m) == 2 {
			return uidFromUsername(m[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return -1, err
	}

	return -1, ErrNotFound
}

func uidFromProcReader(lip net.IP, lport int, rip net.IP, rport int, r io.Reader) (uid int, err error) {
	buf := bufio.NewReader(r)

	localHex := ""
	remoteHex := ""
	ipv4 := lip.To4() != nil
	if ipv4 {
		// In the kernel, the port is run through ntohs(), and
		// the inet_request_socket in
		// include/net/inet_socket.h says the "loc_addr" and
		// "rmt_addr" fields are __be32, but get_openreq4's
		// printf of them is raw, without byte order
		// converstion.
		localHex = fmt.Sprintf("%08X:%04X", toLinuxIPv4Order([]byte(lip.To4())), lport)
		remoteHex = fmt.Sprintf("%08X:%04X", toLinuxIPv4Order([]byte(rip.To4())), rport)
	} else {
		localHex = fmt.Sprintf("%032X:%04X", toLinuxIPv6Order([]byte(lip.To16())), lport)
		remoteHex = fmt.Sprintf("%032X:%04X", toLinuxIPv6Order([]byte(rip.To16())), rport)
	}

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return -1, ErrNotFound
		}
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) < 8 {
			continue
		}
		// log.Printf("parts[1] = %q; localHex = %q", parts[1], localHex)
		if parts[1] == localHex && parts[2] == remoteHex {
			uid, err = strconv.Atoi(parts[7])
			return uid, err
		}
	}
	panic("unreachable")
}
