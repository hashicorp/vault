// Package mdns provides node discovery via mDNS.
package mdns

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"

	m "github.com/hashicorp/mdns"
)

// Provider implements the Provider interface.
type Provider struct{}

// Help returns help information for the mDNS package.
func (p *Provider) Help() string {
	return `mDNS:

    provider:          "mdns"
    service:           The mDNS service name.
    domain:            The mDNS discovery domain.  Default "local".
    timeout:           The mDNS lookup timeout.  Default "5s" (five seconds).
    v6:                IPv6 will be allowed and preferred when set to "true"
                       and disabled when set to "false".  Default "true".
    v4:                IPv4 will be allowed when set to "true" and disabled
                       when set to "false".  Default "true".
`
}

// Addrs returns discovered addresses for the mDNS package.
func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	var params *m.QueryParam
	var ch chan *m.ServiceEntry
	var v6, v4 bool
	var addrs []string
	var err error

	// default to null logger
	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	// init params
	params = new(m.QueryParam)

	// validate and set service record
	if args["service"] == "" {
		return nil, fmt.Errorf("discover-mdns: Service record not provided." +
			"  Please specify a service record for the mDNS lookup.")
	}
	params.Service = args["service"]

	// validate and set domain
	if args["domain"] != "" {
		params.Domain = args["domain"]
	} else {
		params.Domain = "local"
	}

	// validate and set timeout
	if args["timeout"] != "" {
		if params.Timeout, err = time.ParseDuration(args["timeout"]); err != nil {
			return nil, fmt.Errorf("discover-mdns: Failed to parse timeout: %s", err)
		}
	} else {
		params.Timeout = 5 * time.Second
	}

	// validate and set v6 toggle
	if args["v6"] != "" {
		if v6, err = strconv.ParseBool(args["v6"]); err != nil {
			return nil, fmt.Errorf("discover-mdns: Failed to parse v6: %s", err)
		}
	} else {
		v6 = true
	}

	// validate and set v4 toggle
	if args["v4"] != "" {
		if v4, err = strconv.ParseBool(args["v4"]); err != nil {
			return nil, fmt.Errorf("discover-mdns: Failed to parse v4: %s", err)
		}
	} else {
		v4 = true
	}

	// init entries channel
	ch = make(chan *m.ServiceEntry)
	defer close(ch)
	params.Entries = ch

	// build addresses
	go func() {
		var addr string
		for e := range ch {
			addr = "" // reset addr each loop
			if v6 && e.AddrV6 != nil {
				addr = net.JoinHostPort(e.AddrV6.String(),
					strconv.Itoa(e.Port))
			}
			if addr == "" && v4 && e.AddrV4 != nil {
				addr = net.JoinHostPort(e.AddrV4.String(),
					strconv.Itoa(e.Port))
			}
			if addr != "" {
				l.Printf("[DEBUG] discover-mdns: %s -> %s",
					e.Host, addr)
				// build address list
				addrs = append(addrs, addr)
			}
		}
	}()

	// lookup and return
	return addrs, m.Query(params)
}
