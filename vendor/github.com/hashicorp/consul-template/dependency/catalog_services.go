// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	_ Dependency = (*CatalogServicesQuery)(nil)

	// CatalogServicesQueryRe is the regular expression to use for CatalogServicesQuery.
	CatalogServicesQueryRe = regexp.MustCompile(`\A` + queryRe + dcRe + `\z`)
)

func init() {
	gob.Register([]*CatalogSnippet{})
}

// CatalogSnippet is a catalog entry in Consul.
type CatalogSnippet struct {
	Name string
	Tags ServiceTags
}

// CatalogServicesQuery is the representation of a requested catalog service
// dependency from inside a template.
type CatalogServicesQuery struct {
	stopCh chan struct{}

	dc        string
	namespace string
	partition string
}

// NewCatalogServicesQuery parses a string of the format @dc.
func NewCatalogServicesQuery(s string) (*CatalogServicesQuery, error) {
	if !CatalogServicesQueryRe.MatchString(s) {
		return nil, fmt.Errorf("catalog.services: invalid format: %q", s)
	}

	m := regexpMatch(CatalogServicesQueryRe, s)
	queryParams, err := GetConsulQueryOpts(m, "catalog.services")
	if err != nil {
		return nil, err
	}

	return &CatalogServicesQuery{
		stopCh:    make(chan struct{}, 1),
		dc:        m["dc"],
		namespace: queryParams.Get(QueryNamespace),
		partition: queryParams.Get(QueryPartition),
	}, nil
}

// Fetch queries the Consul API defined by the given client and returns a slice
// of CatalogService objects.
func (d *CatalogServicesQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	// default to the query params present while creating NewCatalogServicesQuery
	// and then merge with the query params present in the query
	defaultOpts := &QueryOptions{
		Datacenter:      d.dc,
		ConsulPartition: d.partition,
		ConsulNamespace: d.namespace,
	}

	opts = defaultOpts.Merge(opts)

	log.Printf("[TRACE] %s: GET %s", d, &url.URL{
		Path:     "/v1/catalog/services",
		RawQuery: opts.String(),
	})

	entries, qm, err := clients.Consul().Catalog().Services(opts.ToConsulOpts())
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	log.Printf("[TRACE] %s: returned %d results", d, len(entries))

	var catalogServices []*CatalogSnippet
	for name, tags := range entries {
		catalogServices = append(catalogServices, &CatalogSnippet{
			Name: name,
			Tags: ServiceTags(deepCopyAndSortTags(tags)),
		})
	}

	sort.Stable(ByName(catalogServices))

	rm := &ResponseMetadata{
		LastIndex:   qm.LastIndex,
		LastContact: qm.LastContact,
	}

	return catalogServices, rm, nil
}

// CanShare returns a boolean if this dependency is shareable.
func (d *CatalogServicesQuery) CanShare() bool {
	return true
}

// String returns the human-friendly version of this dependency.
func (d *CatalogServicesQuery) String() string {
	var name string
	if d.dc != "" {
		name = name + "@" + d.dc
	}
	if d.partition != "" {
		name = name + "@partition=" + d.partition
	}
	if d.namespace != "" {
		name = name + "@ns=" + d.namespace
	}

	if len(name) == 0 {
		return "catalog.services"
	}

	return fmt.Sprintf("catalog.services(%s)", name)
}

// Stop halts the dependency's fetch function.
func (d *CatalogServicesQuery) Stop() {
	close(d.stopCh)
}

// Type returns the type of this dependency.
func (d *CatalogServicesQuery) Type() Type {
	return TypeConsul
}

// ByName is a sortable slice of CatalogService structs.
type ByName []*CatalogSnippet

func (s ByName) Len() int           { return len(s) }
func (s ByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByName) Less(i, j int) bool { return s[i].Name < s[j].Name }
