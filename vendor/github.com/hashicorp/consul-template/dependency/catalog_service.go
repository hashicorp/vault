package dependency

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/pkg/errors"
)

var (
	// Ensure implements
	_ Dependency = (*CatalogServiceQuery)(nil)

	// CatalogServiceQueryRe is the regular expression to use.
	CatalogServiceQueryRe = regexp.MustCompile(`\A` + tagRe + serviceNameRe + dcRe + nearRe + `\z`)
)

func init() {
	gob.Register([]*CatalogSnippet{})
}

// CatalogService is a catalog entry in Consul.
type CatalogService struct {
	ID              string
	Node            string
	Address         string
	Datacenter      string
	TaggedAddresses map[string]string
	NodeMeta        map[string]string
	ServiceID       string
	ServiceName     string
	ServiceAddress  string
	ServiceTags     ServiceTags
	ServiceMeta     map[string]string
	ServicePort     int
}

// CatalogServiceQuery is the representation of a requested catalog services
// dependency from inside a template.
type CatalogServiceQuery struct {
	stopCh chan struct{}

	dc   string
	name string
	near string
	tag  string
}

// NewCatalogServiceQuery parses a string into a CatalogServiceQuery.
func NewCatalogServiceQuery(s string) (*CatalogServiceQuery, error) {
	if !CatalogServiceQueryRe.MatchString(s) {
		return nil, fmt.Errorf("catalog.service: invalid format: %q", s)
	}

	m := regexpMatch(CatalogServiceQueryRe, s)
	return &CatalogServiceQuery{
		stopCh: make(chan struct{}, 1),
		dc:     m["dc"],
		name:   m["name"],
		near:   m["near"],
		tag:    m["tag"],
	}, nil
}

// Fetch queries the Consul API defined by the given client and returns a slice
// of CatalogService objects.
func (d *CatalogServiceQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(&QueryOptions{
		Datacenter: d.dc,
		Near:       d.near,
	})

	u := &url.URL{
		Path:     "/v1/catalog/service/" + d.name,
		RawQuery: opts.String(),
	}
	if d.tag != "" {
		q := u.Query()
		q.Set("tag", d.tag)
		u.RawQuery = q.Encode()
	}
	log.Printf("[TRACE] %s: GET %s", d, u)

	entries, qm, err := clients.Consul().Catalog().Service(d.name, d.tag, opts.ToConsulOpts())
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	log.Printf("[TRACE] %s: returned %d results", d, len(entries))

	var list []*CatalogService
	for _, s := range entries {
		list = append(list, &CatalogService{
			ID:              s.ID,
			Node:            s.Node,
			Address:         s.Address,
			Datacenter:      s.Datacenter,
			TaggedAddresses: s.TaggedAddresses,
			NodeMeta:        s.NodeMeta,
			ServiceID:       s.ServiceID,
			ServiceName:     s.ServiceName,
			ServiceAddress:  s.ServiceAddress,
			ServiceTags:     ServiceTags(deepCopyAndSortTags(s.ServiceTags)),
			ServiceMeta:     s.ServiceMeta,
			ServicePort:     s.ServicePort,
		})
	}

	rm := &ResponseMetadata{
		LastIndex:   qm.LastIndex,
		LastContact: qm.LastContact,
	}

	return list, rm, nil
}

// CanShare returns a boolean if this dependency is shareable.
func (d *CatalogServiceQuery) CanShare() bool {
	return true
}

// String returns the human-friendly version of this dependency.
func (d *CatalogServiceQuery) String() string {
	name := d.name
	if d.tag != "" {
		name = d.tag + "." + name
	}
	if d.dc != "" {
		name = name + "@" + d.dc
	}
	if d.near != "" {
		name = name + "~" + d.near
	}
	return fmt.Sprintf("catalog.service(%s)", name)
}

// Stop halts the dependency's fetch function.
func (d *CatalogServiceQuery) Stop() {
	close(d.stopCh)
}

// Type returns the type of this dependency.
func (d *CatalogServiceQuery) Type() Type {
	return TypeConsul
}
