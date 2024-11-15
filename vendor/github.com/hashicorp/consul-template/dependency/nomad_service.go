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
	// Ensure NomadServiceQuery meets the Dependency interface.
	_ Dependency = (*NomadServiceQuery)(nil)

	// NomadServiceQueryRe is the regex that is used to understand a service
	// specific Nomad query.
	//
	// e.g. "<tag=value>.<name>@<region>"
	NomadServiceQueryRe = regexp.MustCompile(`\A` + tagRe + serviceNameRe + regionRe + `\z`)
)

func init() {
	gob.Register([]*NomadService{})
}

// NomadService is a fully hydrated service registration response from the
// mirroring the Nomad API response object.
type NomadService struct {
	ID         string
	Name       string
	Node       string
	Address    string
	Port       int
	Datacenter string
	Tags       ServiceTags
	JobID      string
	AllocID    string
}

// NomadServiceQuery is the representation of a requested Nomad services
// dependency from inside a template.
type NomadServiceQuery struct {
	stopCh chan struct{}

	region string
	name   string
	tag    string
	choose string
}

// NewNomadServiceQuery parses a string into a NomadServiceQuery which is
// used to list services registered within Nomad which match a particular name.
func NewNomadServiceQuery(s string) (*NomadServiceQuery, error) {
	if !NomadServiceQueryRe.MatchString(s) {
		return nil, fmt.Errorf("nomad.service: invalid format: %q", s)
	}

	m := regexpMatch(NomadServiceQueryRe, s)

	return &NomadServiceQuery{
		stopCh: make(chan struct{}, 1),
		region: m["region"],
		name:   m["name"],
		tag:    m["tag"],
	}, nil
}

// NewNomadServiceChooseQuery parses s using NewNomadServiceQuery, and then also
// configures the resulting query with the choose parameter set according to the
// count and key arguments.
func NewNomadServiceChooseQuery(count int, key, s string) (*NomadServiceQuery, error) {
	query, err := NewNomadServiceQuery(s)
	if err != nil {
		return nil, err
	}

	choose := fmt.Sprintf("%d|%s", count, key)
	query.choose = choose

	return query, nil
}

// Fetch queries the Nomad API defined by the given client and returns a slice
// of NomadService objects.
func (d *NomadServiceQuery) Fetch(client *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(&QueryOptions{
		Region: d.region,
		Choose: d.choose,
	})

	u := &url.URL{
		Path:     "/v1/service/" + d.name,
		RawQuery: opts.String(),
	}

	log.Printf("[TRACE] %s: GET %s", d, u)

	entries, qm, err := client.Nomad().Services().Get(d.name, opts.ToNomadOpts())
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	log.Printf("[TRACE] %s: returned %d results", d, len(entries))

	services := make([]*NomadService, 0, len(entries))
	for _, s := range entries {
		// Filter by tag
		if d.tag != "" {
			found := false
			for i := 0; !found && i < len(s.Tags); i++ {
				found = s.Tags[i] == d.tag
			}
			if !found {
				continue
			}
		}

		services = append(services, &NomadService{
			ID:         s.ID,
			Name:       s.ServiceName,
			Node:       s.NodeID,
			Address:    s.Address,
			Port:       s.Port,
			Datacenter: s.Datacenter,
			Tags:       deepCopyAndSortTags(s.Tags),
			JobID:      s.JobID,
			AllocID:    s.AllocID,
		})
	}

	sort.Stable(NomadServiceByName(services))

	rm := &ResponseMetadata{
		LastIndex: qm.LastIndex,
	}

	return services, rm, nil
}

func (d *NomadServiceQuery) CanShare() bool {
	return true
}

func (d *NomadServiceQuery) String() string {
	name := d.name
	if d.tag != "" {
		name = d.tag + "." + name
	}
	if d.region != "" {
		name = name + "@" + d.region
	}
	if d.choose != "" {
		name = name + ":" + d.choose
	}
	return fmt.Sprintf("nomad.service(%s)", name)
}

// Stop halts the dependency's fetch function.
func (d *NomadServiceQuery) Stop() {
	close(d.stopCh)
}

// Type returns the type of this dependency.
func (d *NomadServiceQuery) Type() Type {
	return TypeNomad
}

// NomadServiceByName is a sortable slice of NomadService structs.
type NomadServiceByName []*NomadService

func (s NomadServiceByName) Len() int           { return len(s) }
func (s NomadServiceByName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s NomadServiceByName) Less(i, j int) bool { return s[i].Name < s[j].Name }
