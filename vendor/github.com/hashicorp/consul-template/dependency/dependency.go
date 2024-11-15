// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dependency

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	nomadapi "github.com/hashicorp/nomad/api"
)

const (
	dcRe           = `(@(?P<dc>[[:word:]\.\-\_]+))?`
	keyRe          = `/?(?P<key>[^@\?]+)`
	filterRe       = `(\|(?P<filter>[[:word:]\,]+))?`
	serviceNameRe  = `(?P<name>[[:word:]\-\_]+)`
	queryRe        = `(\?(?P<query>[[:word:]\-\_\=\&]+))?`
	nodeNameRe     = `(?P<name>[[:word:]\.\-\_]+)`
	nearRe         = `(~(?P<near>[[:word:]\.\-\_]+))?`
	prefixRe       = `/?(?P<prefix>[^@\?]+)`
	tagRe          = `((?P<tag>[[:word:]=:\.\-\_]+)\.)?`
	regionRe       = `(@(?P<region>[[:word:]\.\-\_]+))?`
	nvPathRe       = `/?(?P<path>[^@]+)`
	nvNamespaceRe  = `(@(?P<namespace>[[:word:]\-\_]+))?`
	nvListPrefixRe = `/?(?P<prefix>[^@]*)`
	nvListNSRe     = `(@(?P<namespace>([[:word:]\-\_]+|\*)))?`
	nvRegionRe     = `(\.(?P<region>[[:word:]\-\_]+))?`

	DefaultNonBlockingQuerySleepTime = 15 * time.Second
)

type Type int

const (
	TypeConsul Type = iota
	TypeVault
	TypeLocal
	TypeNomad
)

const (
	// DefaultContextTimeout context wait timeout for blocking queries.
	DefaultContextTimeout = 60 * time.Second
)

// Dependency is an interface for a dependency that Consul Template is capable
// of watching.
type Dependency interface {
	Fetch(*ClientSet, *QueryOptions) (interface{}, *ResponseMetadata, error)
	CanShare() bool

	// String is used to uniquely identify a dependency. If a dependency receives arguments and
	// may be used multiple times w/ different args, then it should make sure to differentiate
	// itself based on those args; otherwise, only one instance of the dependency will be created.
	//
	// For example: Calling the exportedServices func multiple times in a template for different
	// partitions will result in the first partition provided being used repeatedly - regardless of the
	// partition provided in later calls - if the partition name is not included in the output of String.
	String() string

	Stop()
	Type() Type
}

// ServiceTags is a slice of tags assigned to a Service
type ServiceTags []string

// QueryOptions is a list of options to send with the query. These options are
// client-agnostic, and the dependency determines which, if any, of the options
// to use.
type QueryOptions struct {
	AllowStale          bool
	Datacenter          string
	Region              string
	Near                string
	Choose              string
	RequireConsistent   bool
	VaultGrace          time.Duration
	WaitIndex           uint64
	WaitTime            time.Duration
	ConsulPeer          string
	ConsulPartition     string
	ConsulNamespace     string
	ConsulSamenessGroup string
}

func (q *QueryOptions) Merge(o *QueryOptions) *QueryOptions {
	var r QueryOptions

	if q == nil {
		if o == nil {
			return &QueryOptions{}
		}
		r = *o
		return &r
	}

	r = *q

	if o == nil {
		return &r
	}

	if o.AllowStale {
		r.AllowStale = o.AllowStale
	}

	if o.Datacenter != "" {
		r.Datacenter = o.Datacenter
	}

	if o.Region != "" {
		r.Region = o.Region
	}

	if o.Near != "" {
		r.Near = o.Near
	}

	if o.Choose != "" {
		r.Choose = o.Choose
	}

	if o.RequireConsistent {
		r.RequireConsistent = o.RequireConsistent
	}

	if o.WaitIndex != 0 {
		r.WaitIndex = o.WaitIndex
	}

	if o.WaitTime != 0 {
		r.WaitTime = o.WaitTime
	}

	if o.ConsulNamespace != "" {
		r.ConsulNamespace = o.ConsulNamespace
	}

	if o.ConsulPartition != "" {
		r.ConsulPartition = o.ConsulPartition
	}

	if o.ConsulPeer != "" {
		r.ConsulPeer = o.ConsulPeer
	}

	if o.ConsulSamenessGroup != "" {
		r.ConsulSamenessGroup = o.ConsulSamenessGroup
	}

	return &r
}

func (q *QueryOptions) ToConsulOpts() *consulapi.QueryOptions {
	return &consulapi.QueryOptions{
		AllowStale:        q.AllowStale,
		Datacenter:        q.Datacenter,
		Namespace:         q.ConsulNamespace,
		Partition:         q.ConsulPartition,
		SamenessGroup:     q.ConsulSamenessGroup,
		Peer:              q.ConsulPeer,
		Near:              q.Near,
		RequireConsistent: q.RequireConsistent,
		WaitIndex:         q.WaitIndex,
		WaitTime:          q.WaitTime,
	}
}

// GetConsulQueryOpts parses optional consul query params into key pairs.
// supports namespace, peer and partition params
func GetConsulQueryOpts(queryMap map[string]string, endpointLabel string) (url.Values, error) {
	queryParams := url.Values{}

	if queryRaw := queryMap["query"]; queryRaw != "" {
		var err error
		queryParams, err = url.ParseQuery(queryRaw)
		if err != nil {
			return nil, fmt.Errorf(
				"%s: invalid query: %q: %s", endpointLabel, queryRaw, err)
		}
		// Validate keys.
		for key := range queryParams {
			switch key {
			case QueryNamespace,
				QueryPeer,
				QueryPartition,
				QuerySamenessGroup:
			default:
				return nil,
					fmt.Errorf("%s: invalid query parameter key %q in query %q: supported keys: %s,%s,%s,%s", endpointLabel, key, queryRaw, QueryNamespace, QueryPeer, QueryPartition, QuerySamenessGroup)
			}
		}
	}

	return queryParams, nil
}

func (q *QueryOptions) ToNomadOpts() *nomadapi.QueryOptions {
	var params map[string]string
	if q.Choose != "" {
		params = map[string]string{
			"choose": q.Choose,
		}
	}
	return &nomadapi.QueryOptions{
		AllowStale: q.AllowStale,
		Region:     q.Region,
		Params:     params,
		WaitIndex:  q.WaitIndex,
		WaitTime:   q.WaitTime,
	}
}

func (q *QueryOptions) String() string {
	u := &url.Values{}

	if q.AllowStale {
		u.Add("stale", strconv.FormatBool(q.AllowStale))
	}

	if q.Datacenter != "" {
		u.Add("dc", q.Datacenter)
	}

	if q.Region != "" {
		u.Add("region", q.Region)
	}

	if q.ConsulNamespace != "" {
		u.Add(QueryNamespace, q.ConsulNamespace)
	}

	if q.ConsulPeer != "" {
		u.Add(QueryPeer, q.ConsulPeer)
	}

	if q.ConsulPartition != "" {
		u.Add(QueryPartition, q.ConsulPartition)
	}

	if q.Near != "" {
		u.Add("near", q.Near)
	}

	if q.Choose != "" {
		u.Add("choose", q.Choose)
	}

	if q.RequireConsistent {
		u.Add("consistent", strconv.FormatBool(q.RequireConsistent))
	}

	if q.WaitIndex != 0 {
		u.Add("index", strconv.FormatUint(q.WaitIndex, 10))
	}

	if q.WaitTime != 0 {
		u.Add("wait", q.WaitTime.String())
	}

	return u.Encode()
}

// ResponseMetadata is a struct that contains metadata about the response. This
// is returned from a Fetch function call.
type ResponseMetadata struct {
	LastIndex   uint64
	LastContact time.Duration
	BlockOnNil  bool // keep blocking on `nil` data returns
}

// deepCopyAndSortTags deep copies the tags in the given string slice and then
// sorts and returns the copied result.
func deepCopyAndSortTags(tags []string) []string {
	newTags := make([]string, 0, len(tags))
	newTags = append(newTags, tags...)
	sort.Strings(newTags)
	return newTags
}

// respWithMetadata is a short wrapper to return the given interface with fake
// response metadata for non-Consul dependencies.
func respWithMetadata(i interface{}) (interface{}, *ResponseMetadata, error) {
	return i, &ResponseMetadata{
		LastContact: 0,
		LastIndex:   uint64(time.Now().Unix()),
	}, nil
}

// regexpMatch matches the given regexp and extracts the match groups into a
// named map.
func regexpMatch(re *regexp.Regexp, q string) map[string]string {
	names := re.SubexpNames()
	match := re.FindAllStringSubmatch(q, -1)

	if len(match) == 0 {
		return map[string]string{}
	}

	m := map[string]string{}
	for i, n := range match[0] {
		if names[i] != "" {
			m[names[i]] = n
		}
	}

	return m
}
