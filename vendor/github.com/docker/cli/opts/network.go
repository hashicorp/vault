package opts

import (
	"encoding/csv"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	networkOptName        = "name"
	networkOptAlias       = "alias"
	networkOptIPv4Address = "ip"
	networkOptIPv6Address = "ip6"
	networkOptMacAddress  = "mac-address"
	networkOptLinkLocalIP = "link-local-ip"
	driverOpt             = "driver-opt"
)

// NetworkAttachmentOpts represents the network options for endpoint creation
type NetworkAttachmentOpts struct {
	Target       string
	Aliases      []string
	DriverOpts   map[string]string
	Links        []string // TODO add support for links in the csv notation of `--network`
	IPv4Address  string
	IPv6Address  string
	LinkLocalIPs []string
	MacAddress   string
}

// NetworkOpt represents a network config in swarm mode.
type NetworkOpt struct {
	options []NetworkAttachmentOpts
}

// Set networkopts value
func (n *NetworkOpt) Set(value string) error { //nolint:gocyclo
	longSyntax, err := regexp.MatchString(`\w+=\w+(,\w+=\w+)*`, value)
	if err != nil {
		return err
	}

	var netOpt NetworkAttachmentOpts
	if longSyntax {
		csvReader := csv.NewReader(strings.NewReader(value))
		fields, err := csvReader.Read()
		if err != nil {
			return err
		}

		netOpt.Aliases = []string{}
		for _, field := range fields {
			// TODO(thaJeztah): these options should not be case-insensitive.
			key, val, ok := strings.Cut(strings.ToLower(field), "=")
			if !ok || key == "" {
				return fmt.Errorf("invalid field %s", field)
			}

			key = strings.TrimSpace(key)
			val = strings.TrimSpace(val)

			switch key {
			case networkOptName:
				netOpt.Target = val
			case networkOptAlias:
				netOpt.Aliases = append(netOpt.Aliases, val)
			case networkOptIPv4Address:
				netOpt.IPv4Address = val
			case networkOptIPv6Address:
				netOpt.IPv6Address = val
			case networkOptMacAddress:
				netOpt.MacAddress = val
			case networkOptLinkLocalIP:
				netOpt.LinkLocalIPs = append(netOpt.LinkLocalIPs, val)
			case driverOpt:
				key, val, err = parseDriverOpt(val)
				if err != nil {
					return err
				}
				if netOpt.DriverOpts == nil {
					netOpt.DriverOpts = make(map[string]string)
				}
				netOpt.DriverOpts[key] = val
			default:
				return errors.New("invalid field key " + key)
			}
		}
		if len(netOpt.Target) == 0 {
			return errors.New("network name/id is not specified")
		}
	} else {
		netOpt.Target = value
	}
	n.options = append(n.options, netOpt)
	return nil
}

// Type returns the type of this option
func (n *NetworkOpt) Type() string {
	return "network"
}

// Value returns the networkopts
func (n *NetworkOpt) Value() []NetworkAttachmentOpts {
	return n.options
}

// String returns the network opts as a string
func (n *NetworkOpt) String() string {
	return ""
}

// NetworkMode return the network mode for the network option
func (n *NetworkOpt) NetworkMode() string {
	networkIDOrName := "default"
	netOptVal := n.Value()
	if len(netOptVal) > 0 {
		networkIDOrName = netOptVal[0].Target
	}
	return networkIDOrName
}

func parseDriverOpt(driverOpt string) (string, string, error) {
	// TODO(thaJeztah): these options should not be case-insensitive.
	// TODO(thaJeztah): should value be converted to lowercase as well, or only the key?
	key, value, ok := strings.Cut(strings.ToLower(driverOpt), "=")
	if !ok || key == "" {
		return "", "", errors.New("invalid key value pair format in driver options")
	}
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	return key, value, nil
}
