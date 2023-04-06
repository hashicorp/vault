package exec

import (
	"io"

	"go.uber.org/atomic"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/consul-template/manager"
	ctconfig "github.com/hashicorp/consul-template/config"

	"github.com/hashicorp/vault/command/agent/config"
)

type ServerConfig struct {
	Logger hclog.Logger
	// Client        *api.Client
	AgentConfig *config.Config

	Namespace string

	// LogLevel is needed to set the internal Consul Template Runner's log level
	// to match the log level of Vault Agent. The internal Runner creates it's own
	// logger and can't be set externally or copied from the Template Server.
	//
	// LogWriter is needed to initialize Consul Template's internal logger to use
	// the same io.Writer that Vault Agent itself is using.
	LogLevel  hclog.Level
	LogWriter io.Writer
}

type Server struct {
	// config holds the ServerConfig used to create it. It's passed along in other
	// methods
	config *ServerConfig

	// runner is the consul-template runner
	runner        *manager.Runner
	runnerStarted *atomic.Bool

	// Templates holds the parsed Consul Templates
	Templates []*ctconfig.TemplateConfig
	// flags when all templates rendered (at least once) and we can start the command
	templatesRendered *atomic.Bool

	// lookupMap is a list of templates indexed by their consul-template ID. This
	// is used to ensure all Vault templates have been rendered before returning
	// from the runner in the event we're using exit after auth.
	lookupMap map[string][]*ctconfig.TemplateConfig

	DoneCh  chan struct{}
	stopped *atomic.Bool

	logger hclog.Logger
}

func NewServer(cfg *ServerConfig) *Server {
	server := Server{
		DoneCh:            make(chan struct{}),
		stopped:           atomic.NewBool(false),
		runnerStarted:     atomic.NewBool(false),
		templatesRendered: atomic.NewBool(false),
		logger:            cfg.Logger,
		config:            cfg,
	}

	return &server
}
