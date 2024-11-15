// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dependency

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/nomad/api"
	"github.com/pkg/errors"
)

var (
	// Ensure implements
	_ Dependency = (*NVListQuery)(nil)

	// NVListQueryRe is the regular expression to use.
	NVListQueryRe = regexp.MustCompile(`\A` + nvListPrefixRe + nvListNSRe + nvRegionRe + `\z`)
)

func init() {
	gob.Register([]*NomadVarMeta{})
}

// NVListQuery queries the SV store for the metadata for keys matching the given
// prefix.
type NVListQuery struct {
	stopCh    chan struct{}
	namespace string
	region    string
	prefix    string
}

// NewNVListQuery parses a string into a dependency.
func NewNVListQuery(ns, s string) (*NVListQuery, error) {
	if s != "" && !NVListQueryRe.MatchString(s) {
		return nil, fmt.Errorf("nomad.var.list: invalid format: %q", s)
	}

	m := regexpMatch(NVListQueryRe, s)
	out := &NVListQuery{
		stopCh:    make(chan struct{}, 1),
		namespace: m["namespace"],
		region:    m["region"],
		prefix:    m["prefix"],
	}

	// Handle paths that are only slashes and discard them
	if strings.Trim(out.prefix, "/") == "" {
		out.prefix = ""
	}

	if out.namespace == "" && ns != "" {
		out.namespace = ns
	}

	return out, nil
}

// Fetch queries the Nomad API defined by the given client.
func (d *NVListQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(&QueryOptions{})

	log.Printf("[TRACE] %s: GET %s", d, &url.URL{
		Path:     "/v1/vars/",
		RawQuery: opts.String(),
	})

	nOpts := opts.ToNomadOpts()
	nOpts.Namespace = d.namespace
	nOpts.Region = d.region
	list, qm, err := clients.Nomad().Variables().PrefixList(d.prefix, nOpts)
	if err != nil && !strings.Contains(err.Error(), "Permission denied") {
		return nil, nil, errors.Wrap(err, d.String())
	}

	log.Printf("[TRACE] %s: returned %d paths", d, len(list))

	vars := make([]*NomadVarMeta, 0, len(list))
	for _, nVar := range list {
		vars = append(vars, NewNomadVarMeta(nVar))
	}

	// 404's don't return QueryMeta.
	if qm == nil {
		qm = &api.QueryMeta{
			LastIndex: 1,
		}
	}

	rm := &ResponseMetadata{
		LastIndex:   qm.LastIndex,
		LastContact: qm.LastContact,
	}

	return vars, rm, nil
}

// CanShare returns a boolean if this dependency is shareable.
func (d *NVListQuery) CanShare() bool {
	return true
}

// String returns the human-friendly version of this dependency.
func (d *NVListQuery) String() string {
	ns := d.namespace
	if ns == "" {
		ns = "default"
	}
	region := d.region
	if region == "" {
		region = "global"
	}
	prefix := d.prefix
	key := fmt.Sprintf("%s@%s.%s", prefix, ns, region)

	return fmt.Sprintf("nomad.var.list(%s)", key)
}

// Stop halts the dependency's fetch function.
func (d *NVListQuery) Stop() {
	close(d.stopCh)
}

// Type returns the type of this dependency.
func (d *NVListQuery) Type() Type {
	return TypeNomad
}
