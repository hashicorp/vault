package vault

import "fmt"

// CapabilitiesResponse holds the result of fetching the capabilities of token on a path
type CapabilitiesResponse struct {
	Root         bool
	Capabilities []string
}

// Capabilities is used to fetch the capabilities of the given token on the given path
func (c *Core) Capabilities(token, path string) (*CapabilitiesResponse, error) {
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

	var policies []*Policy
	for _, tePolicy := range te.Policies {
		policy, err := c.policyStore.GetPolicy(tePolicy)
		if err != nil {
			return nil, err
		}
		policies = append(policies, policy)
	}

	if len(policies) == 0 {
		return nil, nil
	}

	acl, err := NewACL(policies)
	if err != nil {
		return nil, err
	}

	caps := acl.Capabilities(path)
	/*
		log.Printf("vishal: caps:%#v\n", caps)

		var result CapabilitiesResponse
		capabilities := make(map[string]bool)
		for _, tePolicy := range te.Policies {
			if tePolicy == "root" {
				capabilities = map[string]bool{
					"root": true,
				}
				break
			}
			policy, err := c.policyStore.GetPolicy(tePolicy)
			if err != nil {
				return nil, err
			}
			if policy == nil || policy.Paths == nil {
				continue
			}
			for _, pathCapability := range policy.Paths {
				switch {
				case pathCapability.Glob:
					if strings.HasPrefix(path, pathCapability.Prefix) {
						for _, capability := range pathCapability.Capabilities {
							if _, ok := capabilities[capability]; !ok {
								capabilities[capability] = true
							}
						}
					}
				default:
					if path == pathCapability.Prefix {
						for _, capability := range pathCapability.Capabilities {
							if _, ok := capabilities[capability]; !ok {
								capabilities[capability] = true
							}
						}
					}
				}
			}
		}

		if len(capabilities) == 0 {
			result.Capabilities = []string{"deny"}
			return &result, nil
		}

		for capability, _ := range capabilities {
			result.Capabilities = append(result.Capabilities, capability)
		}
		sort.Strings(result.Capabilities)
	*/
	return &result, nil
}
