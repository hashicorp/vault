package sink

import (
	"errors"
	"math/rand"
	"time"

	"github.com/hashicorp/errwrap"
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
	DHType  string
	DHPath  string
	AAD     string
}

// SinkServer is responsible for pushing tokens to sinks
type SinkServer struct {
	DoneCh     chan struct{}
	ShutdownCh chan struct{}
	logger     hclog.Logger
	client     *api.Client
	random     *rand.Rand
}

func NewSinkServer(conf *SinkConfig) *SinkServer {
	ss := &SinkServer{
		ShutdownCh: make(chan struct{}),
		DoneCh:     make(chan struct{}),
		logger:     conf.Logger,
		client:     conf.Client,
		random:     rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
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

	latestToken := new(string)
	sinkCh := make(chan func() error, len(sinks))
	for {
		select {
		case <-ss.ShutdownCh:
			return

		case token := <-incoming:
			if token != *latestToken {

				// Drain the existing funcs
			drainLoop:
				for {
					select {
					case <-sinkCh:
					default:
						break drainLoop
					}
				}

				*latestToken = token

				for _, s := range sinks {
					sinkFunc := func(currSink Sink, currToken string) func() error {
						return func() error {
							if currToken != *latestToken {
								return nil
							}
							if wrapTTL := currSink.WrapTTL(); wrapTTL != 0 {
								wrapClient, err := ss.client.Clone()
								if err != nil {
									return errwrap.Wrapf("error deriving client for wrapping, not writing out to sink: {{err}})", err)
								}
								wrapClient.SetToken(currToken)
								wrapClient.SetWrappingLookupFunc(func(string, string) string {
									return wrapTTL.String()
								})
								secret, err := wrapClient.Logical().Write("sys/wrapping/wrap", map[string]interface{}{
									"token": currToken,
								})
								if err != nil {
									return errwrap.Wrapf("error wrapping token, not writing out to sink: {{err}})", err)
								}
								if secret == nil {
									return errors.New("nil secret returned, not writing out to sink")
								}
								if secret.WrapInfo == nil {
									return errors.New("nil wrap info returned, not writing out to sink")
								}
								if m, err := jsonutil.EncodeJSON(secret.WrapInfo); err != nil {
									return errwrap.Wrapf("error marshaling token, not writing out to sink: {{err}})", err)
								} else {
									currToken = string(m)
								}
							}
							return currSink.WriteToken(currToken)
						}
					}
					sinkCh <- sinkFunc(s, token)
				}
			}

		case sinkFunc := <-sinkCh:
			select {
			case <-ss.ShutdownCh:
				return
			default:
			}

			if err := sinkFunc(); err != nil {
				ss.logger.Error("error returned by sink write function, retrying", "error", err)
				backoff := 2*time.Second + time.Duration(ss.random.Int63()%int64(time.Second*2)-int64(time.Second))
				time.AfterFunc(backoff, func() {
					sinkCh <- sinkFunc
				})
			}
		}
	}
}
