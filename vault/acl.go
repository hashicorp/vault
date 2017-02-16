package vault

import (
	"reflect"

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

			// these are the ones already in the tree
			existingPerms := raw.(*Permissions)

			switch {
			case existingPerms.CapabilitiesBitmap&DenyCapabilityInt > 0:
				// If we are explicitly denied in the existing capability set,
				// don't save anything else
				continue

			case pc.Permissions.CapabilitiesBitmap&DenyCapabilityInt > 0:
				// If this new policy explicitly denies, only save the deny value
				pc.Permissions.CapabilitiesBitmap = DenyCapabilityInt
				pc.Permissions.AllowedParameters = nil
				pc.Permissions.DeniedParameters = nil
				goto INSERT

			default:
				// Insert the capabilities in this new policy into the existing
				// value
				pc.Permissions.CapabilitiesBitmap = existingPerms.CapabilitiesBitmap | pc.Permissions.CapabilitiesBitmap
			}

			if len(existingPerms.AllowedParameters) > 0 {
				if pc.Permissions.AllowedParameters == nil {
					pc.Permissions.AllowedParameters = make(map[string][]interface{}, len(existingPerms.AllowedParameters))
				}

				// If this policy allows everything skip to checking denied
				if _, ok := pc.Permissions.AllowedParameters["*"]; ok {
					goto CHECK_DENIED
				}

				// If the exising policy allows everything set this policy to
				// allow everything and skip to check denied
				if _, ok = existingPerms.AllowedParameters["*"]; ok {
					pc.Permissions.AllowedParameters = map[string][]interface{}{
						"*": []interface{}{},
					}
					goto CHECK_DENIED
				}

				// Merge the two maps, appending values on key conflict.
				for key, value := range existingPerms.AllowedParameters {
					pc.Permissions.AllowedParameters[key] = append(value, pc.Permissions.AllowedParameters[key]...)
				}
			}

		CHECK_DENIED:

			if len(existingPerms.DeniedParameters) > 0 {
				if pc.Permissions.DeniedParameters == nil {
					pc.Permissions.DeniedParameters = make(map[string][]interface{}, len(existingPerms.DeniedParameters))
				}

				// If this policy denies everything go to insert
				if _, ok := pc.Permissions.DeniedParameters["*"]; ok {
					goto INSERT
				}

				// If exising policy denies everything set this policy to
				// deny everything and go to insert
				if _, ok = existingPerms.DeniedParameters["*"]; ok {
					pc.Permissions.DeniedParameters = map[string][]interface{}{
						"*": []interface{}{},
					}
					goto INSERT
				}

				// Merge the two maps, appending values on key conflict.
				for key, value := range existingPerms.DeniedParameters {
					pc.Permissions.DeniedParameters[key] = append(value, pc.Permissions.DeniedParameters[key]...)
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

	// These three re-use UpdateCapabilityInt since that's the most appropriate
	// capability/operation mapping
	case logical.RevokeOperation, logical.RenewOperation, logical.RollbackOperation:
		operationAllowed = capabilities&UpdateCapabilityInt > 0

	default:
		return false, false
	}

	if !operationAllowed {
		return false, sudo
	}

	// Only check parameter permissions for operations that can modify
	// parameters.
	if op == logical.UpdateOperation || op == logical.DeleteOperation || op == logical.CreateOperation {
		// Check if all parameters have been denied
		if _, ok := permissions.DeniedParameters["*"]; ok {
			return false, sudo
		}

		_, allowedAll := permissions.AllowedParameters["*"]
		if len(permissions.DeniedParameters) == 0 && allowedAll {
			return true, sudo
		}

		for parameter, value := range req.Data {
			// Check if parameter has explictly denied
			if valueSlice, ok := permissions.DeniedParameters[parameter]; ok {
				// If the value exists in denied values slice, deny
				if valueInParameterList(value, valueSlice) {
					return false, sudo
				}
				// If the value doesn't exist in the denied values slice,
				// continue
				continue
			}

			// Specfic parameters have been allowed
			if len(permissions.AllowedParameters) > 0 && !allowedAll {
				// Requested parameter is not in allowed list
				valueSlice, ok := permissions.AllowedParameters[parameter]
				if !ok {
					return false, sudo
				}

				// If the value doesn't exists in the allowed values slice,
				// deny
				if !valueInParameterList(value, valueSlice) {
					return false, sudo
				}
			}
		}
	}

	return true, sudo
}

func valueInParameterList(v interface{}, list []interface{}) bool {
	if len(list) == 0 || valueInSlice("*", list) {
		return true
	}

	return valueInSlice(v, list)
}

func valueInSlice(v interface{}, list []interface{}) bool {
	for _, el := range list {
		if reflect.DeepEqual(el, v) {
			return true
		}
	}

	return false
}
