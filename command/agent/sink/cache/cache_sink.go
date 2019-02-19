package cachesink

import (
	"errors"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent"
	"github.com/hashicorp/vault/command/agent/sink"
)

type cacheSink struct {
	logger        hclog.Logger
	clientManager *agent.ClientManager
}

func NewCacheSink(conf *sink.SinkConfig, clientManager *agent.ClientManager) (sink.Sink, error) {
	if conf.Logger == nil {
		return nil, errors.New("nil logger provided")
	}

	conf.Logger.Info("creating cache sink")

	s := &cacheSink{
		logger:        conf.Logger,
		clientManager: clientManager,
	}

	s.logger.Info("cache sink configured")

	return s, nil
}

func (s *cacheSink) WriteToken(token string) error {
	s.logger.Trace("enter write_token")
	defer s.logger.Trace("exit write_token")

	s.clientManager.SetToken(token)

	s.logger.Info("token written")
	return nil
}
