package vault

import (
	"fmt"
	"sort"
	"strings"
)

// CapabilitiesResult holds the result of fetching the capabilities of token on a path
type CapabilitiesResult struct {
	Capabilities []string
}

// Capabilities is used to fetch the capabilities of the given token on the given path
func (c *Core) Capabilities(token, path string) (*CapabilitiesResult, error) {
	if path == "" {
		return nil, fmt.Errorf("missing path")
	}

	if token == "" {
		return nil, fmt.Errorf("missing token")
	}

	te, err := c.tokenStore.Lookup(token)
	if err != nil {
		return nil, err
	}
	if te == nil {
		return nil, fmt.Errorf("invalid token")
	}

	if te.Policies == nil {
		return nil, nil
	}

	maps := make(map[string]bool)
	for _, tePolicy := range te.Policies {
		if tePolicy == "root" {
			//TODO: check if the path is actually a valid path. Otherwise, there is no
			// meaning in returning the capabilities
			// Add all the capabilities
			maps["create"] = true
			maps["read"] = true
			maps["update"] = true
			maps["delete"] = true
			maps["list"] = true
			maps["sudo"] = true
			break
		}
		policy, err := c.policyStore.GetPolicy(tePolicy)
		if err != nil {
			return nil, err
		}
		if policy == nil {
			return nil, fmt.Errorf("policy '%s' not found", tePolicy)
		}

		if policy.Paths == nil {
			return nil, fmt.Errorf("policy '%s' does not contain any paths", tePolicy)
		}
		for _, pathCapability := range policy.Paths {
			switch pathCapability.Glob {
			case true:
				if strings.HasPrefix(path, pathCapability.Prefix) {
					for _, capability := range pathCapability.Capabilities {
						if _, ok := maps[capability]; !ok {
							maps[capability] = true
						}
					}
				}
			case false:
				if path == pathCapability.Prefix {
					for _, capability := range pathCapability.Capabilities {
						if _, ok := maps[capability]; !ok {
							maps[capability] = true
						}
					}
				}
			}
		}
	}

	var capabilities []string
	for capability, _ := range maps {
		capabilities = append(capabilities, capability)
	}
	sort.Strings(capabilities)
	return &CapabilitiesResult{
		Capabilities: capabilities,
	}, nil
}
