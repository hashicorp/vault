package vault

import (
	"fmt"
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
			fmt.Println("There is an existing policy")

			// these are the ones already in the tree
			permissions := raw.(*Permissions)
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

			// look for a * in allowed parameters for the node already in the tree
			if _, ok := permissions.AllowedParameters["*"]; ok {
				pc.Permissions.AllowedParameters = make(map[string]struct{})
				pc.Permissions.AllowedParameters["*"] = struct{}{}
				goto CHECK_DENIED
			}

			// look for a * in allowed parameters for the path capability we are merging
			if _, ok := pc.Permissions.AllowedParameters["*"]; ok {
				pc.Permissions.AllowedParameters = make(map[string]struct{})
				pc.Permissions.AllowedParameters["*"] = struct{}{}
				goto CHECK_DENIED
			}

			// Merge allowed parameters
			for key, _ := range permissions.AllowedParameters {
				// Add new parameter
				if _, ok := pc.Permissions.AllowedParameters[key]; !ok {
					pc.Permissions.AllowedParameters[key] = permissions.AllowedParameters[key]
				}
			}

		CHECK_DENIED:

			// look for a * in denied parameters for the node already in the tree
			if _, ok := permissions.DeniedParameters["*"]; ok {
				pc.Permissions.DeniedParameters = make(map[string]struct{})
				pc.Permissions.DeniedParameters["*"] = struct{}{}
				goto INSERT
			}

			// look for a * in denied parameters for the path capability we are merging
			if _, ok := pc.Permissions.DeniedParameters["*"]; ok {
				pc.Permissions.DeniedParameters = make(map[string]struct{})
				pc.Permissions.DeniedParameters["*"] = struct{}{}
				goto INSERT
			}

			fmt.Println("Entering Merge Denied")
			// Merge denied parameters
			for key, _ := range permissions.DeniedParameters {
				// Add new parameter
				fmt.Println("Checking if already in map")
				if _, ok := pc.Permissions.DeniedParameters[key]; !ok {
					fmt.Printf("DeniedParameter: %v\n", key)
					pc.Permissions.DeniedParameters[key] = permissions.DeniedParameters[key]
				}
			}

		INSERT:

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
		perm := raw.(*Permissions)
		capabilities = perm.CapabilitiesBitmap
		goto CHECK
	}

	// Find a glob rule, default deny if no match
	_, raw, ok = a.globRules.LongestPrefix(path)
	if !ok {
		return []string{DenyCapability}
	} else {
		perm := raw.(*Permissions)
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
func (a *ACL) AllowOperation(req *logical.Request) (allowed bool, sudo bool) {
	fmt.Printf("Operation: %v\nPath: %v\nData: %v\n", req.Operation, req.Path, req.Data)
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

	fmt.Printf("DeniedParameters: %v\n", permissions.DeniedParameters)
	// Only check parameter permissions for operations that can modify parameters.
	if op == logical.UpdateOperation || op == logical.DeleteOperation || op == logical.CreateOperation {
		// Check if all parameters have been denied
		if _, ok := permissions.DeniedParameters["*"]; ok {
			return false, sudo
		}

		for parameter, _ := range req.Data {
			// Check if parameter has explictly denied
			if _, ok := permissions.DeniedParameters[parameter]; ok {
				return false, sudo
			}
			// Specfic parameters have been allowed
			if len(permissions.AllowedParameters) > 0 {
				// Requested parameter is not in allowed list
				if _, ok := permissions.AllowedParameters[parameter]; !ok {
					return false, sudo
				}
			}
		}
		return true, sudo
	}

	return operationAllowed, sudo
}
