package vault

import (
	"github.com/armon/go-radix"
	"github.com/hashicorp/vault/logical"
)

// ACL is used to wrap a set of policies to provide
// an efficient interface for access control.
type ACL struct {
	// exactRules contains the path policies that are exact
	exactRules *radix.Tree

	// globRules contains the path policies that glob
	globRules *radix.Tree

	// root is enabled if the "root" named policy is present.
	root bool
}

// New is used to construct a policy based ACL from a set of policies.
func NewACL(policies []*Policy) (*ACL, error) {
	// Initialize
	a := &ACL{
		exactRules: radix.New(),
		globRules:  radix.New(),
		root:       false,
	}

	// Inject each policy
	for _, policy := range policies {
		// Ignore a nil policy object
		if policy == nil {
			continue
		}
		// Check if this is root
		if policy.Name == "root" {
			a.root = true
		}
		for _, pc := range policy.Paths {
			// Check which tree to use
			tree := a.exactRules
			if pc.Glob {
				tree = a.globRules
			}

			// Check for an existing policy
			raw, ok := tree.Get(pc.Prefix)
			if !ok {
				tree.Insert(pc.Prefix, pc.Permissions)
				continue
			}
			permissions := raw.(Permissions)
			existing := permissions.CapabilitiesBitmap

			switch {
			case existing&DenyCapabilityInt > 0:
				// If we are explicitly denied in the existing capability set,
				// don't save anything else

			case pc.Permissions.CapabilitiesBitmap&DenyCapabilityInt > 0:
				// If this new policy explicitly denies, only save the deny value
				pc.Permissions.CapabilitiesBitmap = DenyCapabilityInt
				tree.Insert(pc.Prefix, pc.Permissions)

			default:
				// Insert the capabilities in this new policy into the existing
				// value
				pc.Permissions.CapabilitiesBitmap = existing | pc.Permissions.CapabilitiesBitmap
				tree.Insert(pc.Prefix, pc.Permissions)
			}

			// Merge allowed parameters
			for key, value := range permissions.AllowedParameters {
				// Add new parameter
				if _, ok := pc.Permissions.AllowedParameters[key]; !ok {
					pc.Permissions.AllowedParameters[key] = permissions.AllowedParameters[key]
					continue
				}

				// Take more general allowed
				if (len(permissions.AllowedParameters[key]) == 0) || (len(pc.Permissions.AllowedParameters[key]) == 0) {
					pc.Permissions.AllowedParameters[key] = nil
					continue
				}

				// Merge allowed values for matching parameters
				for _, element := range value {
					pc.Permissions.AllowedParameters[key] = append(pc.Permissions.AllowedParameters[key], element)
				}
			}

			// Merge disallowed parameters
			for key, value := range permissions.DeniedParameters {
				// Add new parameter
				if _, ok := pc.Permissions.DeniedParameters[key]; !ok {
					pc.Permissions.DeniedParameters[key] = permissions.DeniedParameters[key]
					continue
				}

				// Take more general disallowed
				if (len(permissions.DeniedParameters[key]) == 0) || (len(pc.Permissions.DeniedParameters[key]) == 0) {
					pc.Permissions.DeniedParameters[key] = nil
					continue
				}

				// Merge disallowed values for matching parameters
				for _, element := range value {
					pc.Permissions.DeniedParameters[key] = append(pc.Permissions.DeniedParameters[key], element)
				}
			}

			tree.Insert(pc.Prefix, pc.Permissions)

		}
	}
	return a, nil
}

func (a *ACL) Capabilities(path string) (pathCapabilities []string) {
	// Fast-path root
	if a.root {
		return []string{RootCapability}
	}

	// Find an exact matching rule, look for glob if no match
	var capabilities uint32
	raw, ok := a.exactRules.Get(path)

	if ok {
		perm := raw.(Permissions)
		capabilities = perm.CapabilitiesBitmap
		goto CHECK
	}

	// Find a glob rule, default deny if no match
	_, raw, ok = a.globRules.LongestPrefix(path)
	if !ok {
		return []string{DenyCapability}
	} else {
		perm := raw.(Permissions)
		capabilities = perm.CapabilitiesBitmap
	}

CHECK:
	if capabilities&SudoCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, SudoCapability)
	}
	if capabilities&ReadCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, ReadCapability)
	}
	if capabilities&ListCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, ListCapability)
	}
	if capabilities&UpdateCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, UpdateCapability)
	}
	if capabilities&DeleteCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, DeleteCapability)
	}
	if capabilities&CreateCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, CreateCapability)
	}

	// If "deny" is explicitly set or if the path has no capabilities at all,
	// set the path capabilities to "deny"
	if capabilities&DenyCapabilityInt > 0 || len(pathCapabilities) == 0 {
		pathCapabilities = []string{DenyCapability}
	}
	return
}

// AllowOperation is used to check if the given operation is permitted. The
// first bool indicates if an op is allowed, the second whether sudo priviliges
// exist for that op and path.

// change arguments to hold a full request that holds the operation, path, and parameter
// that is to be modified.
func (a *ACL) AllowOperation(req *logical.Request) (allowed bool, sudo bool) {
	// Fast-path root
	if a.root {
		return true, true
	}

	op := req.Operation
	path := req.Path

	// Help is always allowed
	if op == logical.HelpOperation {
		return true, false
	}

	var permissions *Permissions

	// Find an exact matching rule, look for glob if no match
	var capabilities uint32
	raw, ok := a.exactRules.Get(path)
	if ok {
		permissions = raw.(*Permissions)
		capabilities = permissions.CapabilitiesBitmap
		goto CHECK
	}

	// Find a glob rule, default deny if no match
	_, raw, ok = a.globRules.LongestPrefix(path)
	if !ok {
		return false, false
	} else {
		permissions = raw.(*Permissions)
		capabilities = permissions.CapabilitiesBitmap
	}

CHECK:
	// Check if the minimum permissions are met
	// If "deny" has been explicitly set, only deny will be in the map, so we
	// only need to check for the existence of other values
	sudo = capabilities&SudoCapabilityInt > 0
	operationAllowed := false
	switch op {
	case logical.ReadOperation:
		operationAllowed = capabilities&ReadCapabilityInt > 0
	case logical.ListOperation:
		operationAllowed = capabilities&ListCapabilityInt > 0
	case logical.UpdateOperation:
		operationAllowed = capabilities&UpdateCapabilityInt > 0
	case logical.DeleteOperation:
		operationAllowed = capabilities&DeleteCapabilityInt > 0
	case logical.CreateOperation:
		operationAllowed = capabilities&CreateCapabilityInt > 0

	// These three re-use UpdateCapabilityInt since that's the most appropriate capability/operation mapping
	case logical.RevokeOperation, logical.RenewOperation, logical.RollbackOperation:
		operationAllowed = capabilities&UpdateCapabilityInt > 0

	default:
		return false, false
	}

	if !operationAllowed {
		return false, sudo
	}

	var denied = permissions.DeniedParameters.(map[string]struct{})

	// Check parameter permissions for operations that can modify only.
	if op == logical.UpdateOperation || op == logical.DeleteOperation || op == logical.CreateOperation {
		// Check if all parameters have been denied
		if _, ok := denied["*"]; ok {
			return false, sudo
		}
		for _, value := range req.Data {
			// Check if parameter has explictly been denied
			if _, ok := denied[value]; ok {
				return false, sudo
			}
			// Specfic parameters have been allowed
			if len(permissions.AllowedParameters) > 0 {
				// Requested parameter is not in allowed list
				if _, ok := permissions.AllowedParameters[value]; !ok {
					return false, sudo
				}
			}
		}
		return true, sudo
	}

	return
}
