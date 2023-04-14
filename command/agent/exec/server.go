package exec

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/consul-template/child"
	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/consul-template/manager"
	"github.com/hashicorp/go-hclog"
	"go.uber.org/atomic"

	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/agent/internal/ctmanager"
)

type ServerConfig struct {
	Logger      hclog.Logger
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
	runner *manager.Runner

	// lookupMap is a list of templates indexed by their consul-template ID. This
	// is used to ensure all Vault templates have been rendered before returning
	// from the runner in the event we're using exit after auth.
	lookupMap map[string][]*ctconfig.TemplateConfig

	DoneCh  chan struct{}
	stopped *atomic.Bool

	logger hclog.Logger

	child        *child.Child
	childInput   *child.NewInput
	childStarted *atomic.Bool
	childLock    sync.Mutex

	exitCh chan int
}

func NewServer(cfg *ServerConfig) *Server {
	server := Server{
		DoneCh:       make(chan struct{}),
		stopped:      atomic.NewBool(false),
		logger:       cfg.Logger,
		config:       cfg,
		childStarted: atomic.NewBool(false),
		exitCh:       make(chan int),
	}

	return &server
}

// ExitCh fires if the process dies on its own
func (s *Server) ExitCh() <-chan int {
	return s.exitCh
}

func (s *Server) Run(ctx context.Context, envTmpls map[string]*config.EnvTemplateConfig, execCfg *config.ExecConfig) error {
	s.logger.Info("starting exec server")
	defer func() {
		s.logger.Info("template server stopped")
	}()

	if len(envTmpls) == 0 || execCfg == nil {
		s.logger.Info("no env templates or exec config, exiting")
		<-ctx.Done()
		return nil
	}

	templates := make([]*ctconfig.TemplateConfig, len(envTmpls))

	for envName, envTmpl := range envTmpls {
		tmpl := envTmpl.TemplateConfig
		tmpl.EnvVar = &envName
		templates = append(templates, &tmpl)
	}

	managerConfig := ctmanager.ManagerConfig{
		AgentConfig: s.config.AgentConfig,
		Namespace:   s.config.Namespace,
		LogLevel:    s.config.LogLevel,
		LogWriter:   s.config.LogWriter,
	}

	runnerConfig, runnerConfigErr := ctmanager.NewManagerConfig(managerConfig, templates)
	if runnerConfigErr != nil {
		return fmt.Errorf("template server failed to runner generate config: %w", runnerConfigErr)
	}

	// we leave in "dry" mode, as there's no files
	// we will get the env var rendered contents from incoming events
	var err error
	s.runner, err = manager.NewRunner(runnerConfig, true)
	if err != nil {
		return fmt.Errorf("template server failed to create: %w", err)
	}

	go s.runner.Start()

	idMap := s.runner.TemplateConfigMapping()
	lookupMap := make(map[string][]*ctconfig.TemplateConfig, len(idMap))
	for id, ctmpls := range idMap {
		for _, ctmpl := range ctmpls {
			tl := lookupMap[id]
			tl = append(tl, ctmpl)
			lookupMap[id] = tl
		}
	}
	s.lookupMap = lookupMap

	s.childInput = &child.NewInput{
		Stdin:        os.Stdin,
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
		Command:      execCfg.Command,
		Args:         execCfg.Args,
		Timeout:      0,
		Env:          nil,
		ReloadSignal: nil,
		KillSignal:   os.Interrupt,
		KillTimeout:  30 * time.Second,
		Splay:        0,
		Setpgid:      true,
		Logger:       s.logger.StandardLogger(nil),
	}

	for {
		select {
		case <-ctx.Done():
			s.runner.Stop()
			s.childLock.Lock()
			if s.child != nil {
				// TODO: use kill to immediately kill or gracefully?
				s.child.Kill()
			}
			s.childLock.Unlock()
			return nil
		case err := <-s.runner.ErrCh:
			s.logger.Error("template server error", "error", err.Error())
			s.runner.StopImmediately()

			// Return after stopping the runner if exit on retry failure was specified
			if s.config.AgentConfig.TemplateConfig != nil && s.config.AgentConfig.TemplateConfig.ExitOnRetryFailure {
				return fmt.Errorf("template server: %w", err)
			}

			s.runner, err = manager.NewRunner(runnerConfig, false)
			if err != nil {
				return fmt.Errorf("template server failed to create: %w", err)
			}
			// TODO: why don't we start the runner above?
			go s.runner.Start()
		case <-s.runner.TemplateRenderedCh():
			// A template has been rendered, figure out what to do
			s.logger.Debug("template rendered")
			events := s.runner.RenderEvents()

			// events are keyed by template ID, and can be matched up to the id's from
			// the lookupMap
			if len(events) < len(s.lookupMap) {
				// Not all templates have been rendered yet
				continue
			}

			// assume the renders are finished, until we find otherwise
			doneRendering := true
			envVarToContents := make(map[string]string)
			for _, event := range events {
				// This template hasn't been rendered
				if event.LastWouldRender.IsZero() {
					doneRendering = false
				} else {
					// TODO: check for duplicates?
					for _, tcfg := range event.TemplateConfigs {
						envVarToContents[*tcfg.EnvVar] = string(event.Contents)
					}
				}
			}

			if doneRendering {
				s.logger.Info("done rendering templates, bouncing process")
				if err := s.bounceCmd(envVarToContents); err != nil {
					return fmt.Errorf("unable to bounce command: %w", err)
				}
			}
		}
	}

}

func (s *Server) bounceCmd(newEnvVars map[string]string) error {
	s.childLock.Lock()
	defer s.childLock.Unlock()
	if s.childStarted.Load() && s.child != nil {
		// process is running, need to kill it first
		s.child.Stop()
	}
	var err error
	s.childInput.Env = append(os.Environ(), envsToList(newEnvVars)...)
	s.child, err = child.New(s.childInput)
	if err != nil {
		return err
	}

	// TODO: would this leak?
	// forward process exits to server chan
	go func() {
		select {
		case exitCode := <-s.child.ExitCh():
			s.exitCh <- exitCode
		}
	}()

	if err := s.child.Start(); err != nil {
		return fmt.Errorf("error starting child process: %w", err)
	}
	s.childStarted.Store(true)
	return nil
}

func envsToList(envs map[string]string) []string {
	out := make([]string, len(envs))
	for key, value := range envs {
		e := fmt.Sprintf("%s=%s", key, value)
		out = append(out, e)
	}
	return out
}
