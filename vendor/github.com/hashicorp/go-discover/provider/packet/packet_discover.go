package packet

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/packethost/packngo"
)

const baseURL = "https://api.packet.net/"

// Provider struct
type Provider struct {
	userAgent string
}

// SetUserAgent setter
func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

// Help function
func (p *Provider) Help() string {
	return `Packet:
	provider:       "packet"
	project:        UUID of packet project. Required
	auth_token:     Packet authentication token. Required
	url:            Packet REST URL. Optional
	address_type:   "private_v4", "public_v4" or "public_v6". Defaults to "private_v4". Optional
	facility:       Filter for specific facility (Examples: "ewr1,ams1")
	tag:            Filter by tag (Examples: "tag1,tag2")
	
	Variables can also be provided by environmental variables:
	export PACKET_PROJECT for project
	export PACKET_URL for url
	export PACKET_AUTH_TOKEN for auth_token
`
}

// Addrs function
func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	authToken := argsOrEnv(args, "auth_token", "PACKET_AUTH_TOKEN")
	projectID := argsOrEnv(args, "project", "PACKET_PROJECT")
	packetURL := argsOrEnv(args, "url", "PACKET_URL")
	addressType := args["address_type"]
	packetFacilities := args["facility"]
	packetTags := args["tag"]

	if addressType != "private_v4" && addressType != "public_v4" && addressType != "public_v6" {
		l.Printf("[INFO] discover-packet: Address type %s is not supported. Valid values are {private_v4,public_v4,public_v6}. Falling back to 'private_v4'", addressType)
		addressType = "private_v4"
	}

	includeFacilities := includeArgs(packetFacilities)
	includeTags := includeArgs(packetTags)

	c, err := client(p.userAgent, packetURL, authToken)
	if err != nil {
		return nil, fmt.Errorf("discover-packet: Initializing Packet client failed: %s", err)
	}

	var devices []packngo.Device

	if projectID == "" {
		return nil, fmt.Errorf("discover-packet: 'project' parameter must be provider")
	}

	devices, _, err = c.Devices.List(projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("discover-packet: Fetching Packet devices failed: %s", err)
	}

	var addrs []string
	for _, d := range devices {

		if len(includeFacilities) > 0 && !Include(includeFacilities, d.Facility.Code) {
			continue
		}

		if len(includeTags) > 0 && !Any(d.Tags, func(v string) bool { return Include(includeTags, v) }) {
			continue
		}

		addressFamily := 4
		if addressType == "public_v6" {
			addressFamily = 6
		}
		for _, n := range d.Network {

			if (n.Public == (addressType == "public_v4" || addressType == "public_v6")) && n.AddressFamily == addressFamily {
				addrs = append(addrs, n.Address)
			}
		}
	}
	return addrs, nil
}

func client(useragent, url, token string) (*packngo.Client, error) {
	if url == "" {
		url = baseURL
	}

	return packngo.NewClientWithBaseURL(useragent, token, nil, url)
}
func argsOrEnv(args map[string]string, key, env string) string {
	if value := args[key]; value != "" {
		return value
	}
	return os.Getenv(env)
}

func includeArgs(s string) []string {
	var include []string
	for _, localstring := range strings.Split(s, ",") {
		if len(localstring) > 0 {
			include = append(include, localstring)
		}
	}
	return include
}

// Index returns the first index of the target string t, or -1 if no match is found.
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Include returns true if the target string t is in the slice.
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

//Any returns true if one of the strings in the slice satisfies the predicate f.
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}
