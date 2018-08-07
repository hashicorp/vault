package structs

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strings"

	"golang.org/x/crypto/blake2b"

	multierror "github.com/hashicorp/go-multierror"
	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/nomad/acl"
)

// MergeMultierrorWarnings takes job warnings and canonicalize warnings and
// merges them into a returnable string. Both the errors may be nil.
func MergeMultierrorWarnings(warnings ...error) string {
	var warningMsg multierror.Error
	for _, warn := range warnings {
		if warn != nil {
			multierror.Append(&warningMsg, warn)
		}
	}

	if len(warningMsg.Errors) == 0 {
		return ""
	}

	// Set the formatter
	warningMsg.ErrorFormat = warningsFormatter
	return warningMsg.Error()
}

// warningsFormatter is used to format job warnings
func warningsFormatter(es []error) string {
	points := make([]string, len(es))
	for i, err := range es {
		points[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d warning(s):\n\n%s",
		len(es), strings.Join(points, "\n"))
}

// RemoveAllocs is used to remove any allocs with the given IDs
// from the list of allocations
func RemoveAllocs(alloc []*Allocation, remove []*Allocation) []*Allocation {
	// Convert remove into a set
	removeSet := make(map[string]struct{})
	for _, remove := range remove {
		removeSet[remove.ID] = struct{}{}
	}

	n := len(alloc)
	for i := 0; i < n; i++ {
		if _, ok := removeSet[alloc[i].ID]; ok {
			alloc[i], alloc[n-1] = alloc[n-1], nil
			i--
			n--
		}
	}

	alloc = alloc[:n]
	return alloc
}

// FilterTerminalAllocs filters out all allocations in a terminal state and
// returns the latest terminal allocations
func FilterTerminalAllocs(allocs []*Allocation) ([]*Allocation, map[string]*Allocation) {
	terminalAllocsByName := make(map[string]*Allocation)
	n := len(allocs)
	for i := 0; i < n; i++ {
		if allocs[i].TerminalStatus() {

			// Add the allocation to the terminal allocs map if it's not already
			// added or has a higher create index than the one which is
			// currently present.
			alloc, ok := terminalAllocsByName[allocs[i].Name]
			if !ok || alloc.CreateIndex < allocs[i].CreateIndex {
				terminalAllocsByName[allocs[i].Name] = allocs[i]
			}

			// Remove the allocation
			allocs[i], allocs[n-1] = allocs[n-1], nil
			i--
			n--
		}
	}
	return allocs[:n], terminalAllocsByName
}

// AllocsFit checks if a given set of allocations will fit on a node.
// The netIdx can optionally be provided if its already been computed.
// If the netIdx is provided, it is assumed that the client has already
// ensured there are no collisions.
func AllocsFit(node *Node, allocs []*Allocation, netIdx *NetworkIndex) (bool, string, *Resources, error) {
	// Compute the utilization from zero
	used := new(Resources)

	// Add the reserved resources of the node
	if node.Reserved != nil {
		if err := used.Add(node.Reserved); err != nil {
			return false, "", nil, err
		}
	}

	// For each alloc, add the resources
	for _, alloc := range allocs {
		if alloc.Resources != nil {
			if err := used.Add(alloc.Resources); err != nil {
				return false, "", nil, err
			}
		} else if alloc.TaskResources != nil {

			// Adding the shared resource asks for the allocation to the used
			// resources
			if err := used.Add(alloc.SharedResources); err != nil {
				return false, "", nil, err
			}
			// Allocations within the plan have the combined resources stripped
			// to save space, so sum up the individual task resources.
			for _, taskResource := range alloc.TaskResources {
				if err := used.Add(taskResource); err != nil {
					return false, "", nil, err
				}
			}
		} else {
			return false, "", nil, fmt.Errorf("allocation %q has no resources set", alloc.ID)
		}
	}

	// Check that the node resources are a super set of those
	// that are being allocated
	if superset, dimension := node.Resources.Superset(used); !superset {
		return false, dimension, used, nil
	}

	// Create the network index if missing
	if netIdx == nil {
		netIdx = NewNetworkIndex()
		defer netIdx.Release()
		if netIdx.SetNode(node) || netIdx.AddAllocs(allocs) {
			return false, "reserved port collision", used, nil
		}
	}

	// Check if the network is overcommitted
	if netIdx.Overcommitted() {
		return false, "bandwidth exceeded", used, nil
	}

	// Allocations fit!
	return true, "", used, nil
}

// ScoreFit is used to score the fit based on the Google work published here:
// http://www.columbia.edu/~cs2035/courses/ieor4405.S13/datacenter_scheduling.ppt
// This is equivalent to their BestFit v3
func ScoreFit(node *Node, util *Resources) float64 {
	// Determine the node availability
	nodeCpu := float64(node.Resources.CPU)
	if node.Reserved != nil {
		nodeCpu -= float64(node.Reserved.CPU)
	}
	nodeMem := float64(node.Resources.MemoryMB)
	if node.Reserved != nil {
		nodeMem -= float64(node.Reserved.MemoryMB)
	}

	// Compute the free percentage
	freePctCpu := 1 - (float64(util.CPU) / nodeCpu)
	freePctRam := 1 - (float64(util.MemoryMB) / nodeMem)

	// Total will be "maximized" the smaller the value is.
	// At 100% utilization, the total is 2, while at 0% util it is 20.
	total := math.Pow(10, freePctCpu) + math.Pow(10, freePctRam)

	// Invert so that the "maximized" total represents a high-value
	// score. Because the floor is 20, we simply use that as an anchor.
	// This means at a perfect fit, we return 18 as the score.
	score := 20.0 - total

	// Bound the score, just in case
	// If the score is over 18, that means we've overfit the node.
	if score > 18.0 {
		score = 18.0
	} else if score < 0 {
		score = 0
	}
	return score
}

func CopySliceConstraints(s []*Constraint) []*Constraint {
	l := len(s)
	if l == 0 {
		return nil
	}

	c := make([]*Constraint, l)
	for i, v := range s {
		c[i] = v.Copy()
	}
	return c
}

// VaultPoliciesSet takes the structure returned by VaultPolicies and returns
// the set of required policies
func VaultPoliciesSet(policies map[string]map[string]*Vault) []string {
	set := make(map[string]struct{})

	for _, tgp := range policies {
		for _, tp := range tgp {
			for _, p := range tp.Policies {
				set[p] = struct{}{}
			}
		}
	}

	flattened := make([]string, 0, len(set))
	for p := range set {
		flattened = append(flattened, p)
	}
	return flattened
}

// DenormalizeAllocationJobs is used to attach a job to all allocations that are
// non-terminal and do not have a job already. This is useful in cases where the
// job is normalized.
func DenormalizeAllocationJobs(job *Job, allocs []*Allocation) {
	if job != nil {
		for _, alloc := range allocs {
			if alloc.Job == nil && !alloc.TerminalStatus() {
				alloc.Job = job
			}
		}
	}
}

// AllocName returns the name of the allocation given the input.
func AllocName(job, group string, idx uint) string {
	return fmt.Sprintf("%s.%s[%d]", job, group, idx)
}

// ACLPolicyListHash returns a consistent hash for a set of policies.
func ACLPolicyListHash(policies []*ACLPolicy) string {
	cacheKeyHash, err := blake2b.New256(nil)
	if err != nil {
		panic(err)
	}
	for _, policy := range policies {
		cacheKeyHash.Write([]byte(policy.Name))
		binary.Write(cacheKeyHash, binary.BigEndian, policy.ModifyIndex)
	}
	cacheKey := string(cacheKeyHash.Sum(nil))
	return cacheKey
}

// CompileACLObject compiles a set of ACL policies into an ACL object with a cache
func CompileACLObject(cache *lru.TwoQueueCache, policies []*ACLPolicy) (*acl.ACL, error) {
	// Sort the policies to ensure consistent ordering
	sort.Slice(policies, func(i, j int) bool {
		return policies[i].Name < policies[j].Name
	})

	// Determine the cache key
	cacheKey := ACLPolicyListHash(policies)
	aclRaw, ok := cache.Get(cacheKey)
	if ok {
		return aclRaw.(*acl.ACL), nil
	}

	// Parse the policies
	parsed := make([]*acl.Policy, 0, len(policies))
	for _, policy := range policies {
		p, err := acl.Parse(policy.Rules)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q: %v", policy.Name, err)
		}
		parsed = append(parsed, p)
	}

	// Create the ACL object
	aclObj, err := acl.NewACL(false, parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to construct ACL: %v", err)
	}

	// Update the cache
	cache.Add(cacheKey, aclObj)
	return aclObj, nil
}

// GenerateMigrateToken will create a token for a client to access an
// authenticated volume of another client to migrate data for sticky volumes.
func GenerateMigrateToken(allocID, nodeSecretID string) (string, error) {
	h, err := blake2b.New512([]byte(nodeSecretID))
	if err != nil {
		return "", err
	}
	h.Write([]byte(allocID))
	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}

// CompareMigrateToken returns true if two migration tokens can be computed and
// are equal.
func CompareMigrateToken(allocID, nodeSecretID, otherMigrateToken string) bool {
	h, err := blake2b.New512([]byte(nodeSecretID))
	if err != nil {
		return false
	}
	h.Write([]byte(allocID))

	otherBytes, err := base64.URLEncoding.DecodeString(otherMigrateToken)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(h.Sum(nil), otherBytes) == 1
}
