// Package scaleway provides node discovery for Scaleway.
package scaleway

import (
	"fmt"
	"io/ioutil"
	"log"

	api "github.com/nicolai86/scaleway-sdk"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Scaleway:

    provider:     "scaleway"
    organization: The Scaleway organization access key
    tag_name:     The tag name to filter on
    token:        The Scaleway API access token
    region:       The Scalway region
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "scaleway" {
		return nil, fmt.Errorf("discover-scaleway: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	organization := args["organization"]
	tagName := args["tag_name"]
	token := args["token"]
	region := args["region"]

	l.Printf("[INFO] discover-scaleway: Organization is %q", organization)
	l.Printf("[INFO] discover-scaleway: Region is %q", region)

	// Create a new API client
	api, err := api.New(organization, token, region)
	if err != nil {
		return nil, fmt.Errorf("discover-scaleway: %s", err)
	}

	// Currently fetching all servers since the API doesn't support
	// filter options.
	// api.GetServers() takes the following two arguments:
	// * all (bool) - lets you list all servers in any state (stopped, running etc)
	// * limit (int) - limits the results to a certain number. In this case we are listing
	servers, err := api.GetServers(true, 0)
	if err != nil {
		return nil, fmt.Errorf("discover-scaleway: %s", err)
	}

	// Filter servers by tag
	var addrs []string
	if servers != nil {
		for _, server := range servers {
			if stringInSlice(tagName, server.Tags) {
				l.Printf("[DEBUG] discover-scaleway: Found server (%s) - %s with private IP: %s",
					server.Name, server.Hostname, server.PrivateIP)
				addrs = append(addrs, server.PrivateIP)
			}
		}
	}

	l.Printf("[DEBUG] discover-scaleway: Found ip addresses: %v", addrs)
	return addrs, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
