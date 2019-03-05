package inmem

import (
	"errors"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/cache"
	"github.com/hashicorp/vault/command/agent/sink"
)

// inmemSink retains the auto-auth token in memory and exposes it via
// sink.SinkReader interface.
type inmemSink struct {
	logger     hclog.Logger
	token      string
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
	}, nil
}

func (s *inmemSink) WriteToken(token string) error {
	s.token = token

	if s.leaseCache != nil {
		s.leaseCache.RegisterAutoAuthToken(token)
	}

	return nil
}

func (s *inmemSink) Token() string {
	return s.token
}
