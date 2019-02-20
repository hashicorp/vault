package passthroughsink

import (
	"errors"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/sink"
)

// passthroughSink is a simple sink that will send the token to an output channel
type passthroughSink struct {
	logger   hclog.Logger
	outputCh chan<- string
}

func New(conf *sink.SinkConfig, outputCh chan<- string) (sink.Sink, error) {
	if conf.Logger == nil {
		return nil, errors.New("nil logger provided")
	}

	conf.Logger.Info("creating passthrough sink")

	s := &passthroughSink{
		logger:   conf.Logger,
		outputCh: outputCh,
	}

	s.logger.Info("passthrough sink configured")

	return s, nil
}

func (s *passthroughSink) WriteToken(token string) error {
	s.logger.Trace("enter write_token")
	defer s.logger.Trace("exit write_token")

	s.outputCh <- token

	s.logger.Info("token written")
	return nil
}
