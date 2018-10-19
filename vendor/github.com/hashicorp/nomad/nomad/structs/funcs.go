package structs

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strconv"
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
func AllocsFit(node *Node, allocs []*Allocation, netIdx *NetworkIndex) (bool, string, *ComparableResources, error) {
	// Compute the utilization from zero
	used := new(ComparableResources)

	// Add the reserved resources of the node
	used.Add(node.ComparableReservedResources())

	// For each alloc, add the resources
	for _, alloc := range allocs {
		// Do not consider the resource impact of terminal allocations
		if alloc.TerminalStatus() {
			continue
		}

		used.Add(alloc.ComparableResources())
	}

	// Check that the node resources are a super set of those
	// that are being allocated
	if superset, dimension := node.ComparableResources().Superset(used); !superset {
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
func ScoreFit(node *Node, util *ComparableResources) float64 {
	// COMPAT(0.11): Remove in 0.11
	reserved := node.ComparableReservedResources()
	res := node.ComparableResources()

	// Determine the node availability
	nodeCpu := float64(res.Flattened.Cpu.CpuShares)
	nodeMem := float64(res.Flattened.Memory.MemoryMB)
	if reserved != nil {
		nodeCpu -= float64(reserved.Flattened.Cpu.CpuShares)
		nodeMem -= float64(reserved.Flattened.Memory.MemoryMB)
	}

	// Compute the free percentage
	freePctCpu := 1 - (float64(util.Flattened.Cpu.CpuShares) / nodeCpu)
	freePctRam := 1 - (float64(util.Flattened.Memory.MemoryMB) / nodeMem)

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

func CopySliceAffinities(s []*Affinity) []*Affinity {
	l := len(s)
	if l == 0 {
		return nil
	}

	c := make([]*Affinity, l)
	for i, v := range s {
		c[i] = v.Copy()
	}
	return c
}

func CopySliceSpreads(s []*Spread) []*Spread {
	l := len(s)
	if l == 0 {
		return nil
	}

	c := make([]*Spread, l)
	for i, v := range s {
		c[i] = v.Copy()
	}
	return c
}

func CopySliceSpreadTarget(s []*SpreadTarget) []*SpreadTarget {
	l := len(s)
	if l == 0 {
		return nil
	}

	c := make([]*SpreadTarget, l)
	for i, v := range s {
		c[i] = v.Copy()
	}
	return c
}

func CopySliceNodeScoreMeta(s []*NodeScoreMeta) []*NodeScoreMeta {
	l := len(s)
	if l == 0 {
		return nil
	}

	c := make([]*NodeScoreMeta, l)
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

// ParsePortRanges parses the passed port range string and returns a list of the
// ports. The specification is a comma separated list of either port numbers or
// port ranges. A port number is a single integer and a port range is two
// integers separated by a hyphen. As an example the following spec would
// convert to: ParsePortRanges("10,12-14,16") -> []uint64{10, 12, 13, 14, 16}
func ParsePortRanges(spec string) ([]uint64, error) {
	parts := strings.Split(spec, ",")

	// Hot path the empty case
	if len(parts) == 1 && parts[0] == "" {
		return nil, nil
	}

	ports := make(map[uint64]struct{})
	for _, part := range parts {
		part = strings.TrimSpace(part)
		rangeParts := strings.Split(part, "-")
		l := len(rangeParts)
		switch l {
		case 1:
			if val := rangeParts[0]; val == "" {
				return nil, fmt.Errorf("can't specify empty port")
			} else {
				port, err := strconv.ParseUint(val, 10, 0)
				if err != nil {
					return nil, err
				}
				ports[port] = struct{}{}
			}
		case 2:
			// We are parsing a range
			start, err := strconv.ParseUint(rangeParts[0], 10, 0)
			if err != nil {
				return nil, err
			}

			end, err := strconv.ParseUint(rangeParts[1], 10, 0)
			if err != nil {
				return nil, err
			}

			if end < start {
				return nil, fmt.Errorf("invalid range: starting value (%v) less than ending (%v) value", end, start)
			}

			for i := start; i <= end; i++ {
				ports[i] = struct{}{}
			}
		default:
			return nil, fmt.Errorf("can only parse single port numbers or port ranges (ex. 80,100-120,150)")
		}
	}

	var results []uint64
	for port := range ports {
		results = append(results, port)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i] < results[j]
	})
	return results, nil
}
