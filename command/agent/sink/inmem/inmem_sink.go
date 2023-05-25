// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package inmem

import (
	"errors"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/cache"
	"github.com/hashicorp/vault/command/agent/sink"
	"go.uber.org/atomic"
)

// inmemSink retains the auto-auth token in memory and exposes it via
// sink.SinkReader interface.
type inmemSink struct {
	logger     hclog.Logger
	token      *atomic.String
	leaseCache *cache.LeaseCache
}

// New creates a new instance of inmemSink.
func New(conf *sink.SinkConfig, leaseCache *cache.LeaseCache) (sink.Sink, error) {
	if conf.Logger == nil {
		return nil, errors.New("nil logger provided")
	}

	return &inmemSink{
		logger:     conf.Logger,
		leaseCache: leaseCache,
		token:      atomic.NewString(""),
	}, nil
}

func (s *inmemSink) WriteToken(token string) error {
	s.token.Store(token)

	if s.leaseCache != nil {
		s.leaseCache.RegisterAutoAuthToken(token)
	}

	return nil
}

func (s *inmemSink) Token() string {
	return s.token.Load()
}
