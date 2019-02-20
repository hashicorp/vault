package inmem

import (
	"errors"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/sink"
)

// inmemSink retains the auto-auth token in memory and exposes it via
// sink.SinkReader interface.
type inmemSink struct {
	logger hclog.Logger
	token  string
}

func New(conf *sink.SinkConfig) (sink.Sink, error) {
	if conf.Logger == nil {
		return nil, errors.New("nil logger provided")
	}
	return &inmemSink{
		logger: conf.Logger,
	}, nil
}

func (s *inmemSink) WriteToken(token string) error {
	s.token = token
	return nil
}

func (s *inmemSink) Token() string {
	return s.token
}
