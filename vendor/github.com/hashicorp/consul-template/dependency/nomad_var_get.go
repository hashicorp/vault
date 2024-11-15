// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dependency

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	// Ensure implements
	_ Dependency = (*NVGetQuery)(nil)

	// NVGetQueryRe is the regular expression to use.
	NVGetQueryRe = regexp.MustCompile(`\A` + nvPathRe + nvNamespaceRe + nvRegionRe + `\z`)
)

// NVGetQuery queries the KV store for a single key.
type NVGetQuery struct {
	stopCh chan struct{}

	path      string
	namespace string
	region    string

	blockOnNil bool
}

// NewNVGetQuery parses a string into a dependency.
func NewNVGetQuery(ns, s string) (*NVGetQuery, error) {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "/")

	if s != "" && !NVGetQueryRe.MatchString(s) {
		return nil, fmt.Errorf("nomad.var.get: invalid format: %q", s)
	}

	m := regexpMatch(NVGetQueryRe, s)
	out := &NVGetQuery{
		stopCh:    make(chan struct{}, 1),
		path:      m["path"],
		namespace: m["namespace"],
		region:    m["region"],
	}
	if out.namespace == "" && ns != "" {
		out.namespace = ns
	}
	return out, nil
}

// Fetch queries the Nomad API defined by the given client.
func (d *NVGetQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-d.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(&QueryOptions{})

	log.Printf("[TRACE] %s: GET %s", d, &url.URL{
		Path:     "/v1/var/" + d.path,
		RawQuery: opts.String(),
	})

	nOpts := opts.ToNomadOpts()
	nOpts.Namespace = d.namespace
	nOpts.Region = d.region
	// NOTE: The Peek method of the Nomad Variables API will check a value,
	// return it if it exists, but return a nil value and NO error if it is
	// not found.
	nVar, qm, err := clients.Nomad().Variables().Peek(d.path, nOpts)
	if err != nil {
		return nil, nil, errors.Wrap(err, d.String())
	}

	rm := &ResponseMetadata{
		LastIndex:   qm.LastIndex,
		LastContact: qm.LastContact,
		BlockOnNil:  d.blockOnNil,
	}

	if nVar == nil {
		log.Printf("[TRACE] %s: returned nil", d)
		return nil, rm, nil
	}

	items := &NewNomadVariable(nVar).Items
	log.Printf("[TRACE] %s: returned %q", d, nVar.Path)
	return items, rm, nil
}

// EnableBlocking turns this into a blocking KV query.
func (d *NVGetQuery) EnableBlocking() {
	d.blockOnNil = true
}

// CanShare returns a boolean if this dependency is shareable.
func (d *NVGetQuery) CanShare() bool {
	return true
}

// String returns the human-friendly version of this dependency.
// This value is also used to disambiguate multiple instances in the Brain
func (d *NVGetQuery) String() string {
	ns := d.namespace
	if ns == "" {
		ns = "default"
	}
	region := d.region
	if region == "" {
		region = "global"
	}
	path := d.path
	key := fmt.Sprintf("%s@%s.%s", path, ns, region)
	if d.blockOnNil {
		return fmt.Sprintf("nomad.var.block(%s)", key)
	}
	return fmt.Sprintf("nomad.var.get(%s)", key)
}

// Stop halts the dependency's fetch function.
func (d *NVGetQuery) Stop() {
	close(d.stopCh)
}

// Type returns the type of this dependency.
func (d *NVGetQuery) Type() Type {
	return TypeNomad
}
