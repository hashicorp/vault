package gocbconnstr

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	// DefaultHttpPort is the default HTTP port to use to connect to Couchbase Server.
	DefaultHttpPort = 8091

	// DefaultSslHttpPort is the default HTTPS port to use to connect to Couchbase Server.
	DefaultSslHttpPort = 18091

	// DefaultMemdPort is the default memd port to use to connect to Couchbase Server.
	DefaultMemdPort = 11210

	// DefaultSslMemdPort is the default memd SSL port to use to connect to Couchbase Server.
	DefaultSslMemdPort = 11207

	// DefaultCouchbase2Port is the default (SSL) port to use to connect with Couchbase2.
	DefaultCouchbase2Port = 18098
)

const (
	couchbaseScheme = iota + 1
	httpScheme
	nsServerScheme
	couchbase2Scheme
)

func hostIsIpAddress(host string) bool {
	if strings.HasPrefix(host, "[") {
		// This is an IPv6 address
		return true
	}
	if net.ParseIP(host) != nil {
		// This is an IPv4 address
		return true
	}
	return false
}

// Address represents a host:port pair.
type Address struct {
	Host string
	Port int
}

// ConnSpec describes a connection specification.
type ConnSpec struct {
	Scheme    string
	Addresses []Address
	Bucket    string
	Options   map[string][]string
}

func (spec ConnSpec) srvRecord() (string, string, string, bool) {
	// Only `couchbase`-type schemes allow SRV records
	if spec.Scheme != "couchbase" && spec.Scheme != "couchbases" {
		return "", "", "", false
	}

	// Must have only a single host, with no port specified
	if len(spec.Addresses) != 1 || spec.Addresses[0].Port != -1 {
		return "", "", "", false
	}

	if hostIsIpAddress(spec.Addresses[0].Host) {
		return "", "", "", false
	}

	return spec.Scheme, "tcp", spec.Addresses[0].Host, true
}

// SrvRecordName returns the record name for the ConnSpec.
func (spec ConnSpec) SrvRecordName() (recordName string) {
	scheme, proto, host, isValid := spec.srvRecord()
	if !isValid {
		return ""
	}

	return fmt.Sprintf("_%s._%s.%s", scheme, proto, host)
}

// GetOption returns the specified option value for the ConnSpec.
func (spec ConnSpec) GetOption(name string) []string {
	if opt, ok := spec.Options[name]; ok {
		return opt
	}
	return nil
}

// GetOptionString returns the specified option value for the ConnSpec.
func (spec ConnSpec) GetOptionString(name string) string {
	opts := spec.GetOption(name)
	if len(opts) > 0 {
		return opts[0]
	}
	return ""
}

// Parse parses the connection string into a ConnSpec.
func Parse(connStr string) (out ConnSpec, err error) {
	partMatcher := regexp.MustCompile(`((.*):\/\/)?(([^\/?:]*)(:([^\/?:@]*))?@)?([^\/?]*)(\/([^\?]*))?(\?(.*))?`)
	hostMatcher := regexp.MustCompile(`((\[[^\]]+\]+)|([^;\,\:]+))(:([0-9]*))?(;\,)?`)
	parts := partMatcher.FindStringSubmatch(connStr)

	if parts[2] != "" {
		out.Scheme = parts[2]
	}

	if parts[7] != "" {
		hosts := hostMatcher.FindAllStringSubmatch(parts[7], -1)
		for _, hostInfo := range hosts {
			address := Address{
				Host: hostInfo[1],
				Port: -1,
			}

			if hostInfo[5] != "" {
				address.Port, err = strconv.Atoi(hostInfo[5])
				if err != nil {
					return
				}
			}

			out.Addresses = append(out.Addresses, address)
		}
	}

	if parts[9] != "" {
		out.Bucket, err = url.QueryUnescape(parts[9])
		if err != nil {
			return
		}
	}

	if parts[11] != "" {
		out.Options, err = url.ParseQuery(parts[11])
		if err != nil {
			return
		}
	}

	return
}

func (spec ConnSpec) String() string {
	var out string

	if spec.Scheme != "" {
		out += fmt.Sprintf("%s://", spec.Scheme)
	}

	for i, address := range spec.Addresses {
		if i > 0 {
			out += ","
		}

		if address.Port >= 0 {
			out += fmt.Sprintf("%s:%d", address.Host, address.Port)
		} else {
			out += address.Host
		}
	}

	if spec.Bucket != "" {
		out += "/"
		out += spec.Bucket
	}

	urlOptions := url.Values(spec.Options)
	if len(urlOptions) > 0 {
		out += "?" + urlOptions.Encode()
	}

	return out
}

// ResolvedConnSpec is the result of resolving a ConnSpec.
type ResolvedConnSpec struct {
	UseSsl         bool
	MemdHosts      []Address
	HttpHosts      []Address
	NSServerHost   *Address
	Couchbase2Host *Address
	Bucket         string
	SrvRecord      *SrvRecord
	Options        map[string][]string
}

// SrvRecord contains the information about the srv record used to extract addresses.
type SrvRecord struct {
	Proto  string
	Scheme string
	Host   string
}

// Resolve parses a ConnSpec into a ResolvedConnSpec.
func Resolve(connSpec ConnSpec) (ResolvedConnSpec, error) {
	defaultPort := 0
	hasExplicitScheme := false
	var scheme int
	useSsl := false

	switch connSpec.Scheme {
	case "couchbase":
		defaultPort = DefaultMemdPort
		hasExplicitScheme = true
		scheme = couchbaseScheme
		useSsl = false
	case "couchbases":
		defaultPort = DefaultSslMemdPort
		hasExplicitScheme = true
		scheme = couchbaseScheme
		useSsl = true
	case "http":
		defaultPort = DefaultHttpPort
		hasExplicitScheme = true
		scheme = httpScheme
		useSsl = false
	case "ns_server":
		return handleNsServerScheme(connSpec)
	case "ns_servers":
		return handleNsServerScheme(connSpec)
	case "couchbase2":
		return handleCouchbase2Scheme(connSpec)
	case "":
		defaultPort = DefaultHttpPort
		hasExplicitScheme = false
		scheme = httpScheme
		useSsl = false
	default:
		return ResolvedConnSpec{}, errors.New("bad scheme")
	}

	var srvRecords []*net.SRV
	srvScheme, srvProto, srvHost, srvIsValid := connSpec.srvRecord()
	if srvIsValid {
		_, addrs, err := net.LookupSRV(srvScheme, srvProto, srvHost)
		if err == nil && len(addrs) > 0 {
			srvRecords = addrs
		}
	}

	out := ResolvedConnSpec{}
	out.UseSsl = useSsl
	out.Bucket = connSpec.Bucket
	out.Options = connSpec.Options

	if srvRecords != nil {
		for _, srv := range srvRecords {
			out.MemdHosts = append(out.MemdHosts, Address{
				Host: strings.TrimSuffix(srv.Target, "."),
				Port: int(srv.Port),
			})
		}
		return out, nil
	}

	if len(connSpec.Addresses) == 0 {
		appendMemdAndHttpHostsWithDefaultPorts(&out, useSsl, "127.0.0.1")
		return out, nil
	}

	for _, address := range connSpec.Addresses {
		hasExplicitPort := address.Port > 0

		if !hasExplicitScheme && hasExplicitPort && address.Port != defaultPort {
			return ResolvedConnSpec{}, errors.New("ambiguous port without scheme")
		}

		if hasExplicitScheme && scheme == couchbaseScheme && address.Port == DefaultHttpPort {
			return ResolvedConnSpec{}, errors.New("couchbase://host:8091 not supported for couchbase:// scheme. Use couchbase://host")
		}

		if isDefaultOrNoPortOrDefaultHttpPort(address.Port, defaultPort) {
			appendMemdAndHttpHostsWithDefaultPorts(&out, useSsl, address.Host)
		} else {
			switch scheme {
			case couchbaseScheme:
				out.MemdHosts = append(out.MemdHosts, makeAddress(address.Host, address.Port))
			case httpScheme:
				out.HttpHosts = append(out.HttpHosts, makeAddress(address.Host, address.Port))
			}
		}
	}

	return out, nil
}

func handleCouchbase2Scheme(connSpec ConnSpec) (ResolvedConnSpec, error) {
	out := ResolvedConnSpec{}
	out.UseSsl = true
	out.Bucket = connSpec.Bucket
	out.Options = connSpec.Options

	if connSpec.Bucket != "" {
		return ResolvedConnSpec{}, errors.New("couchbase2 scheme cannot only be used with bucket option")
	}

	if len(connSpec.Addresses) > 1 {
		return ResolvedConnSpec{}, errors.New("couchbase2 scheme can only be used with a single host")
	}

	if len(connSpec.Addresses) == 0 {
		return populateCouchbase2Host(&out, "127.0.0.1"), nil
	}

	address := connSpec.Addresses[0]
	if address.Port <= 0 || address.Port == DefaultCouchbase2Port {
		return populateCouchbase2Host(&out, address.Host), nil
	}

	a := makeAddress(address.Host, address.Port)
	out.Couchbase2Host = &a
	return out, nil
}

func populateCouchbase2Host(out *ResolvedConnSpec, address string) ResolvedConnSpec {
	addr := makeAddress(address, DefaultCouchbase2Port)
	out.Couchbase2Host = &addr

	return *out
}

func handleNsServerScheme(connSpec ConnSpec) (ResolvedConnSpec, error) {
	out := ResolvedConnSpec{}
	out.UseSsl = true
	out.Bucket = connSpec.Bucket
	out.Options = connSpec.Options

	if len(connSpec.Addresses) > 1 {
		return ResolvedConnSpec{}, errors.New("ns_server schemes can only be used with a single host")
	}

	if len(connSpec.Addresses) == 0 {
		return populateNsServerHost(&out, "127.0.0.1"), nil
	}

	address := connSpec.Addresses[0]
	if isDefaultOrNoPortOrDefaultHttpPort(address.Port, DefaultHttpPort) {
		return populateNsServerHost(&out, address.Host), nil
	}

	a := makeAddress(address.Host, address.Port)
	out.NSServerHost = &a
	return out, nil
}

func populateNsServerHost(out *ResolvedConnSpec, address string) ResolvedConnSpec {
	addr := makeAddress(address, DefaultHttpPort)
	out.NSServerHost = &addr

	return *out
}

func appendMemdAndHttpHostsWithDefaultPorts(out *ResolvedConnSpec, useSsl bool, addr string) {
	if useSsl {
		out.MemdHosts = append(out.MemdHosts, makeAddress(addr, DefaultSslMemdPort))
		out.HttpHosts = append(out.HttpHosts, makeAddress(addr, DefaultSslHttpPort))
	} else {
		out.MemdHosts = append(out.MemdHosts, makeAddress(addr, DefaultMemdPort))
		out.HttpHosts = append(out.HttpHosts, makeAddress(addr, DefaultHttpPort))
	}
}

func makeAddress(host string, port int) Address {
	return Address{
		Host: host,
		Port: port,
	}
}

func isDefaultOrNoPortOrDefaultHttpPort(port int, defaultPort int) bool {
	return port <= 0 || port == defaultPort || port == DefaultHttpPort
}
