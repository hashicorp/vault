package sink

import (
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/jsonutil"
)

type Sink interface {
	WriteToken(string) error
	WrapTTL() time.Duration
}

type SinkConfig struct {
	Logger  hclog.Logger
	Config  map[string]interface{}
	Client  *api.Client
	WrapTTL time.Duration
}

// SinkServer is responsible for pushing tokens to sinks
type SinkServer struct {
	DoneCh     chan struct{}
	ShutdownCh chan struct{}
	logger     hclog.Logger
	client     *api.Client
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
					sinkToken := token
					if wrapTTL := sink.WrapTTL(); wrapTTL != 0 {
						wrapClient, err := ss.client.Clone()
						if err != nil {
							ss.logger.Error("error deriving client for wrapping, not writing out to sink", "error", err)
							continue
						}
						wrapClient.SetToken(token)
						wrapClient.SetWrappingLookupFunc(func(string, string) string {
							return wrapTTL.String()
						})
						secret, err := wrapClient.Logical().Write("sys/wrapping/wrap", map[string]interface{}{
							"token": token,
						})
						if err != nil {
							ss.logger.Error("error wrapping token, not writing out to sink", "error", err)
							continue
						}
						if secret == nil {
							ss.logger.Error("nil secret returned, not writing out to sink", "error", err)
							continue
						}
						if secret.WrapInfo == nil {
							ss.logger.Error("nil wrap info returned, not writing out to sink", "error", err)
							continue
						}
						if m, err := jsonutil.EncodeJSON(secret.WrapInfo); err != nil {
							ss.logger.Error("error marshaling response, not writing out to sink", "error", err)
							continue
						} else {
							sinkToken = string(m)
						}
					}
					if err := sink.WriteToken(sinkToken); err != nil {
						ss.logger.Error("error writing token", "error", err)
					}
				}
			}
			prevToken = token
		}
	}
}
