package sink

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
)

type Sink interface {
	WriteToken(string) error
}

type SinkConfig struct {
	Logger hclog.Logger
	Config map[string]interface{}
}

// SinkServer is responsible for pushing tokens to sinks
type SinkServer struct {
	DoneCh     chan struct{}
	ShutdownCh chan struct{}
	logger     hclog.Logger
	client     *api.Client
}

type SinkConfig struct {
	Logger hclog.Logger
	Client *api.Client
}

func NewSinkServer(conf *SinkConfig) *SinkServer {
	ss := &SinkServer{
		ShutdownCh: make(chan struct{}),
		DoneCh:     make(chan struct{}),
		logger:     conf.Logger,
		client:     conf.Client,
	}

	return ss
}

// Run executes the server's run loop, which is responsible for reading
// in new tokens and pushing them out to the various sinks.
func (ss *SinkServer) Run(incoming chan string, sinks []Sink) {
	if incoming == nil {
		panic("incoming or shutdown channel are nil")
	}

	ss.logger.Info("starting sink server")
	defer func() {
		ss.logger.Info("sink server stopped")
		close(ss.DoneCh)
	}()

	var prevToken string
	for {
		select {
		case <-ss.ShutdownCh:
			return
		case token := <-incoming:
			if prevToken != token {
				for _, sink := range sinks {
					if err := sink.WriteToken(token); err != nil {
						ss.logger.Error("error writing token", "error", err)
					}
				}
			}
			prevToken = token
		}
	}
}
