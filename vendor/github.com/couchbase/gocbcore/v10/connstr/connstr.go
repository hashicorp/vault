package connstr

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
)

const (
	couchbaseScheme = iota + 1
	httpScheme
	nsServerScheme
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
	var onlyAllowSingleHost bool

	if parts[2] != "" {
		out.Scheme = parts[2]

		switch out.Scheme {
		case "couchbase":
		case "couchbases":
		case "http":
		case "ns_server":
			onlyAllowSingleHost = true
		case "ns_servers":
			onlyAllowSingleHost = true
		default:
			err = errors.New("bad scheme")
			return
		}
	}

	if parts[7] != "" {
		hosts := hostMatcher.FindAllStringSubmatch(parts[7], -1)
		if len(hosts) > 1 && onlyAllowSingleHost {
			err = errors.New("ns_server scheme can only be used with a single host")
			return
		}
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
	UseSsl       bool
	MemdHosts    []Address
	HttpHosts    []Address
	NSServerHost *Address
	Bucket       string
	Options      map[string][]string
	SrvRecord    *SrvRecord
}

// SrvRecord contains the information about the srv record used to extract addresses.
type SrvRecord struct {
	Proto  string
	Scheme string
	Host   string
}

// Resolve parses a ConnSpec into a ResolvedConnSpec.
func Resolve(connSpec ConnSpec) (out ResolvedConnSpec, err error) {
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
		defaultPort = DefaultHttpPort
		hasExplicitScheme = true
		scheme = nsServerScheme
		useSsl = true
	case "":
		defaultPort = DefaultHttpPort
		hasExplicitScheme = false
		scheme = httpScheme
		useSsl = false
	default:
		err = errors.New("bad scheme")
		return
	}

	var srvRecords []*net.SRV
	srvScheme, srvProto, srvHost, srvIsValid := connSpec.srvRecord()
	if srvIsValid {
		_, addrs, err := net.LookupSRV(srvScheme, srvProto, srvHost)
		if err == nil && len(addrs) > 0 {
			srvRecords = addrs
		}
	}

	if srvRecords != nil {
		for _, srv := range srvRecords {
			out.MemdHosts = append(out.MemdHosts, Address{
				Host: strings.TrimSuffix(srv.Target, "."),
				Port: int(srv.Port),
			})
		}
		out.SrvRecord = &SrvRecord{
			Host:   srvHost,
			Proto:  srvProto,
			Scheme: srvScheme,
		}
	} else if len(connSpec.Addresses) == 0 {
		if scheme == nsServerScheme {
			out.NSServerHost = &Address{
				Host: "127.0.0.1",
				Port: DefaultHttpPort,
			}
		} else {
			if useSsl {
				out.MemdHosts = append(out.MemdHosts, Address{
					Host: "127.0.0.1",
					Port: DefaultSslMemdPort,
				})
				out.HttpHosts = append(out.HttpHosts, Address{
					Host: "127.0.0.1",
					Port: DefaultSslHttpPort,
				})
			} else {
				out.MemdHosts = append(out.MemdHosts, Address{
					Host: "127.0.0.1",
					Port: DefaultMemdPort,
				})
				out.HttpHosts = append(out.HttpHosts, Address{
					Host: "127.0.0.1",
					Port: DefaultHttpPort,
				})
			}
		}
	} else {
		for _, address := range connSpec.Addresses {
			hasExplicitPort := address.Port > 0

			if !hasExplicitScheme && hasExplicitPort && address.Port != defaultPort {
				err = errors.New("ambiguous port without scheme")
				return
			}

			if hasExplicitScheme && scheme == couchbaseScheme && address.Port == DefaultHttpPort {
				err = errors.New("couchbase://host:8091 not supported for couchbase:// scheme. Use couchbase://host")
				return
			}

			if address.Port <= 0 || address.Port == defaultPort || address.Port == DefaultHttpPort {
				if scheme == nsServerScheme {
					out.NSServerHost = &Address{
						Host: address.Host,
						Port: DefaultHttpPort,
					}
				} else {
					if useSsl {
						out.MemdHosts = append(out.MemdHosts, Address{
							Host: address.Host,
							Port: DefaultSslMemdPort,
						})
						out.HttpHosts = append(out.HttpHosts, Address{
							Host: address.Host,
							Port: DefaultSslHttpPort,
						})
					} else {
						out.MemdHosts = append(out.MemdHosts, Address{
							Host: address.Host,
							Port: DefaultMemdPort,
						})
						out.HttpHosts = append(out.HttpHosts, Address{
							Host: address.Host,
							Port: DefaultHttpPort,
						})
					}
				}
			} else {
				switch scheme {
				case couchbaseScheme:
					out.MemdHosts = append(out.MemdHosts, Address{
						Host: address.Host,
						Port: address.Port,
					})
				case httpScheme:
					out.HttpHosts = append(out.HttpHosts, Address{
						Host: address.Host,
						Port: address.Port,
					})
				case nsServerScheme:
					out.NSServerHost = &Address{
						Host: address.Host,
						Port: address.Port,
					}
				}
			}
		}
	}

	out.UseSsl = useSsl
	out.Bucket = connSpec.Bucket
	out.Options = connSpec.Options
	return
}
