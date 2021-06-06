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
	_ Dependency = (*CatalogNodeQuery)(nil)

	// CatalogNodeQueryRe is the regular expression to use.
	CatalogNodeQueryRe = regexp.MustCompile(`\A` + nodeNameRe + dcRe + `\z`)
)

func init() {
	gob.Register([]*CatalogNode{})
	gob.Register([]*CatalogNodeService{})
}

// CatalogNodeQuery represents a single node from the Consul catalog.
type CatalogNodeQuery struct {
	stopCh chan struct{}

	dc   string
	name string
}

// CatalogNode is a wrapper around the node and its services.
type CatalogNode struct {
	Node     *Node
	Services []*CatalogNodeService
}

// CatalogNodeService is a service on a single node.
type CatalogNodeService struct {
	ID                string
	Service           string
	Tags              ServiceTags
	Meta              map[string]string
	Port              int
	Address           string
	EnableTagOverride bool
}

// NewCatalogNodeQuery parses the given string into a dependency. If the name is
// empty then the name of the local agent is used.
func NewCatalogNodeQuery(s string) (*CatalogNodeQuery, error) {
	if s != "" && !CatalogNodeQueryRe.MatchString(s) {
		return nil, fmt.Errorf("catalog.node: invalid format: %q", s)
	}

	m := regexpMatch(CatalogNodeQueryRe, s)
	return &CatalogNodeQuery{
		dc:     m["dc"],
		name:   m["name"],
		stopCh: make(chan struct{}, 1),
	}, nil
}

// Fetch queries the Consul API defined by the given client and returns a
// of CatalogNode object.
func (d *CatalogNodeQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(&QueryOptions{
		Datacenter: d.dc,
	})

	// Grab the name
	name := d.name

	if name == "" {
		log.Printf("[TRACE] %s: getting local agent name", d)
		var err error
		name, err = clients.Consul().Agent().NodeName()
		if err != nil {
			return nil, nil, errors.Wrapf(err, d.String())
		}
	}

	log.Printf("[TRACE] %s: GET %s", d, &url.URL{
		Path:     "/v1/catalog/node/" + name,
		RawQuery: opts.String(),
	})
	node, qm, err := clients.Consul().Catalog().Node(name, opts.ToConsulOpts())
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	log.Printf("[TRACE] %s: returned response", d)

	rm := &ResponseMetadata{
		LastIndex:   qm.LastIndex,
		LastContact: qm.LastContact,
	}

	if node == nil {
		log.Printf("[WARN] %s: no node exists with the name %q", d, name)
		var node CatalogNode
		return &node, rm, nil
	}

	services := make([]*CatalogNodeService, 0, len(node.Services))
	for _, v := range node.Services {
		services = append(services, &CatalogNodeService{
			ID:                v.ID,
			Service:           v.Service,
			Tags:              ServiceTags(deepCopyAndSortTags(v.Tags)),
			Meta:              v.Meta,
			Port:              v.Port,
			Address:           v.Address,
			EnableTagOverride: v.EnableTagOverride,
		})
	}
	sort.Stable(ByService(services))

	detail := &CatalogNode{
		Node: &Node{
			ID:              node.Node.ID,
			Node:            node.Node.Node,
			Address:         node.Node.Address,
			Datacenter:      node.Node.Datacenter,
			TaggedAddresses: node.Node.TaggedAddresses,
			Meta:            node.Node.Meta,
		},
		Services: services,
	}

	return detail, rm, nil
}

// CanShare returns a boolean if this dependency is shareable.
func (d *CatalogNodeQuery) CanShare() bool {
	return false
}

// String returns the human-friendly version of this dependency.
func (d *CatalogNodeQuery) String() string {
	name := d.name
	if d.dc != "" {
		name = name + "@" + d.dc
	}

	if name == "" {
		return "catalog.node"
	}
	return fmt.Sprintf("catalog.node(%s)", name)
}

// Stop halts the dependency's fetch function.
func (d *CatalogNodeQuery) Stop() {
	close(d.stopCh)
}

// Type returns the type of this dependency.
func (d *CatalogNodeQuery) Type() Type {
	return TypeConsul
}

// ByService is a sorter of node services by their service name and then ID.
type ByService []*CatalogNodeService

func (s ByService) Len() int      { return len(s) }
func (s ByService) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByService) Less(i, j int) bool {
	if s[i].Service == s[j].Service {
		return s[i].ID <= s[j].ID
	}
	return s[i].Service <= s[j].Service
}
