// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"fmt"
	"time"
)

const (
	// DefaultDedupPrefix is the default prefix used for deduplication mode.
	DefaultDedupPrefix = "consul-template/dedup/"

	// DefaultDedupTTL is the default TTL for deduplicate mode.
	DefaultDedupTTL = 15 * time.Second

	// DefaultDedupMaxStale is the default max staleness for the deduplication
	// manager.
	DefaultDedupMaxStale = DefaultMaxStale

	// DefaultDedupBlockQueryWaitTime is the default amount of time to do a blocking query for the deduplication
	DefaultDedupBlockQueryWaitTime = 60 * time.Second
)

// DedupConfig is used to enable the de-duplication mode, which depends
// on electing a leader per-template and watching of a key. This is used
// to reduce the cost of many instances of CT running the same template.
type DedupConfig struct {
	// Controls if deduplication mode is enabled
	Enabled *bool `mapstructure:"enabled"`

	// MaxStale is the maximum amount of time to allow for stale queries.
	MaxStale *time.Duration `mapstructure:"max_stale"`

	// Controls the KV prefix used. Defaults to defaultDedupPrefix
	Prefix *string `mapstructure:"prefix"`

	// TTL is the Session TTL used for lock acquisition, defaults to 15 seconds.
	TTL *time.Duration `mapstructure:"ttl"`

	// BlockQueryWaitTime is amount of time to do a blocking query for, defaults to 60 seconds.
	BlockQueryWaitTime *time.Duration `mapstructure:"block_query_wait"`
}

// DefaultDedupConfig returns a configuration that is populated with the
// default values.
func DefaultDedupConfig() *DedupConfig {
	return &DedupConfig{}
}

// Copy returns a deep copy of this configuration.
func (c *DedupConfig) Copy() *DedupConfig {
	if c == nil {
		return nil
	}

	var o DedupConfig
	o.Enabled = c.Enabled
	o.MaxStale = c.MaxStale
	o.Prefix = c.Prefix
	o.TTL = c.TTL
	o.BlockQueryWaitTime = c.BlockQueryWaitTime
	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *DedupConfig) Merge(o *DedupConfig) *DedupConfig {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	if o.Enabled != nil {
		r.Enabled = o.Enabled
	}

	if o.MaxStale != nil {
		r.MaxStale = o.MaxStale
	}

	if o.Prefix != nil {
		r.Prefix = o.Prefix
	}

	if o.TTL != nil {
		r.TTL = o.TTL
	}

	if o.BlockQueryWaitTime != nil {
		r.BlockQueryWaitTime = o.BlockQueryWaitTime
	}

	return r
}

// Finalize ensures there no nil pointers.
func (c *DedupConfig) Finalize() {
	if c.Enabled == nil {
		c.Enabled = Bool(false ||
			TimeDurationPresent(c.MaxStale) ||
			StringPresent(c.Prefix) ||
			TimeDurationPresent(c.TTL) ||
			TimeDurationPresent(c.BlockQueryWaitTime))
	}

	if c.MaxStale == nil {
		c.MaxStale = TimeDuration(DefaultDedupMaxStale)
	}

	if c.Prefix == nil {
		c.Prefix = String(DefaultDedupPrefix)
	}

	if c.TTL == nil {
		c.TTL = TimeDuration(DefaultDedupTTL)
	}

	if c.BlockQueryWaitTime == nil {
		c.BlockQueryWaitTime = TimeDuration(DefaultDedupBlockQueryWaitTime)
	}
}

// GoString defines the printable version of this struct.
func (c *DedupConfig) GoString() string {
	if c == nil {
		return "(*DedupConfig)(nil)"
	}
	return fmt.Sprintf("&DedupConfig{"+
		"Enabled:%s, "+
		"MaxStale:%s, "+
		"Prefix:%s, "+
		"TTL:%s, "+
		"BlockQueryWaitTime:%s"+
		"}",
		BoolGoString(c.Enabled),
		TimeDurationGoString(c.MaxStale),
		StringGoString(c.Prefix),
		TimeDurationGoString(c.TTL),
		TimeDurationGoString(c.BlockQueryWaitTime),
	)
}
