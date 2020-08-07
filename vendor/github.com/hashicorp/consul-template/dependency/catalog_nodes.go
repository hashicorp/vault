package dependency

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"sort"

	"github.com/pkg/errors"
)

var (
	// Ensure implements
	_ Dependency = (*CatalogNodesQuery)(nil)

	// CatalogNodesQueryRe is the regular expression to use.
	CatalogNodesQueryRe = regexp.MustCompile(`\A` + dcRe + nearRe + `\z`)
)

func init() {
	gob.Register([]*Node{})
}

// Node is a node entry in Consul
type Node struct {
	ID              string
	Node            string
	Address         string
	Datacenter      string
	TaggedAddresses map[string]string
	Meta            map[string]string
}

// CatalogNodesQuery is the representation of all registered nodes in Consul.
type CatalogNodesQuery struct {
	stopCh chan struct{}

	dc   string
	near string
}

// NewCatalogNodesQuery parses the given string into a dependency. If the name is
// empty then the name of the local agent is used.
func NewCatalogNodesQuery(s string) (*CatalogNodesQuery, error) {
	if !CatalogNodesQueryRe.MatchString(s) {
		return nil, fmt.Errorf("catalog.nodes: invalid format: %q", s)
	}

	m := regexpMatch(CatalogNodesQueryRe, s)
	return &CatalogNodesQuery{
		dc:     m["dc"],
		near:   m["near"],
		stopCh: make(chan struct{}, 1),
	}, nil
}

// Fetch queries the Consul API defined by the given client and returns a slice
// of Node objects
func (d *CatalogNodesQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(&QueryOptions{
		Datacenter: d.dc,
		Near:       d.near,
	})

	log.Printf("[TRACE] %s: GET %s", d, &url.URL{
		Path:     "/v1/catalog/nodes",
		RawQuery: opts.String(),
	})
	n, qm, err := clients.Consul().Catalog().Nodes(opts.ToConsulOpts())
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	log.Printf("[TRACE] %s: returned %d results", d, len(n))

	nodes := make([]*Node, 0, len(n))
	for _, node := range n {
		nodes = append(nodes, &Node{
			ID:              node.ID,
			Node:            node.Node,
			Address:         node.Address,
			Datacenter:      node.Datacenter,
			TaggedAddresses: node.TaggedAddresses,
			Meta:            node.Meta,
		})
	}

	// Sort unless the user explicitly asked for nearness
	if d.near == "" {
		sort.Stable(ByNode(nodes))
	}

	rm := &ResponseMetadata{
		LastIndex:   qm.LastIndex,
		LastContact: qm.LastContact,
	}

	return nodes, rm, nil
}

// CanShare returns a boolean if this dependency is shareable.
func (d *CatalogNodesQuery) CanShare() bool {
	return true
}

// String returns the human-friendly version of this dependency.
func (d *CatalogNodesQuery) String() string {
	name := ""
	if d.dc != "" {
		name = name + "@" + d.dc
	}
	if d.near != "" {
		name = name + "~" + d.near
	}

	if name == "" {
		return "catalog.nodes"
	}
	return fmt.Sprintf("catalog.nodes(%s)", name)
}

// Stop halts the dependency's fetch function.
func (d *CatalogNodesQuery) Stop() {
	close(d.stopCh)
}

// Type returns the type of this dependency.
func (d *CatalogNodesQuery) Type() Type {
	return TypeConsul
}

// ByNode is a sortable list of nodes by name and then IP address.
type ByNode []*Node

func (s ByNode) Len() int      { return len(s) }
func (s ByNode) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByNode) Less(i, j int) bool {
	if s[i].Node == s[j].Node {
		return s[i].Address <= s[j].Address
	}
	return s[i].Node <= s[j].Node
}
