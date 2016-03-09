package vault

// Struct to identify user input errors.
// This is helpful in responding the appropriate status codes to clients
// from the HTTP endpoints.
type ErrUserInput struct {
	Message string
}

// Implementing error interface
func (e *ErrUserInput) Error() string {
	return e.Message
}

// CapabilitiesAccessor is used to fetch the capabilities of the token associated with
// the token accessor an the given path
func (c *Core) CapabilitiesAccessor(accessorID, path string) ([]string, error) {
	if path == "" {
		return nil, &ErrUserInput{
			Message: "missing path",
		}
	}

	if accessorID == "" {
		return nil, &ErrUserInput{
			Message: "missing accessor_id",
		}
	}

	token, err := c.tokenStore.lookupByAccessorID(accessorID)
	if err != nil {
		return nil, err
	}

	return c.Capabilities(token, path)
}

// Capabilities is used to fetch the capabilities of the given token on the given path
func (c *Core) Capabilities(token, path string) ([]string, error) {
	if path == "" {
		return nil, &ErrUserInput{
			Message: "missing path",
		}
	}

	if token == "" {
		return nil, &ErrUserInput{
			Message: "missing token",
		}
	}

	te, err := c.tokenStore.Lookup(token)
	if err != nil {
		return nil, err
	}
	if te == nil {
		return nil, &ErrUserInput{
			Message: "invalid token",
		}
	}

	if te.Policies == nil {
		return []string{DenyCapability}, nil
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
		return []string{DenyCapability}, nil
	}

	acl, err := NewACL(policies)
	if err != nil {
		return nil, err
	}

	return acl.Capabilities(path), nil
}
