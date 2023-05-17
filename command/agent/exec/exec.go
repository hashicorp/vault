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

	stopped *atomic.Bool

	logger hclog.Logger

	proc        *child.Child
	procStarted bool
	procLock    sync.Mutex

	// exit channel of the child process
	exitCh <-chan int
}

type ProcessExitError struct {
	ExitCode int
}

func (e *ProcessExitError) Error() string {
	return fmt.Sprintf("process exited with %d", e.ExitCode)
}

func NewServer(cfg *ServerConfig) *Server {
	server := Server{
		stopped:     atomic.NewBool(false),
		logger:      cfg.Logger,
		config:      cfg,
		procStarted: false,
		// exitCh: make(<-chan int),
	}

	return &server
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("starting exec server")
	defer func() {
		s.logger.Info("exec server stopped")
	}()

	if len(s.config.AgentConfig.EnvTemplates) == 0 || s.config.AgentConfig.Exec == nil {
		s.logger.Info("no env templates or exec config, exiting")
		<-ctx.Done()
		return nil
	}

	templates := make([]*ctconfig.TemplateConfig, 0, len(s.config.AgentConfig.EnvTemplates))

	for _, envTmpl := range s.config.AgentConfig.EnvTemplates {
		templates = append(templates, envTmpl)
	}

	managerConfig := ctmanager.ManagerConfig{
		AgentConfig: s.config.AgentConfig,
		Namespace:   s.config.Namespace,
		LogLevel:    s.config.LogLevel,
		LogWriter:   s.config.LogWriter,
	}

	runnerConfig, runnerConfigErr := ctmanager.NewConfig(managerConfig, templates)
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

	for {
		select {
		case <-ctx.Done():
			s.runner.Stop()
			s.procLock.Lock()
			if s.proc != nil {
				s.proc.Stop()
			}
			s.procStarted = false
			s.proc.Unlock()
			return nil
		case err := <-s.runner.ErrCh:
			s.logger.Error("template server error", "error", err.Error())
			s.runner.StopImmediately()

			// Return after stopping the runner if exit on retry failure was specified
			if s.config.AgentConfig.TemplateConfig != nil && s.config.AgentConfig.TemplateConfig.ExitOnRetryFailure {
				return fmt.Errorf("template server: %w", err)
			}

			s.runner, err = manager.NewRunner(runnerConfig, true)
			if err != nil {
				return fmt.Errorf("template server failed to create: %w", err)
			}
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
			envVarToContents := map[string]string{}
			for _, event := range events {
				// This template hasn't been rendered
				if event.LastWouldRender.IsZero() {
					doneRendering = false
				} else {
					// TODO: check for duplicates?
					for _, tcfg := range event.TemplateConfigs {
						envVarToContents[*tcfg.MapToEnvironmentVariable] = string(event.Contents)
					}
				}
			}

			if doneRendering {
				s.logger.Debug("done rendering templates/detected change, bouncing process")
				if err := s.bounceCmd(envVarToContents); err != nil {
					return fmt.Errorf("unable to bounce command: %w", err)
				}
			}
		case exitCode := <-s.exitCh:
			// process exited on its own
			return &ProcessExitError{ExitCode: exitCode}
		}
	}
}

func (s *Server) bounceCmd(newEnvVars map[string]string) error {
	s.procLock.Lock()
	defer s.procLock.Unlock()

	// TODO: replace w/enum?
	switch s.config.AgentConfig.Exec.RestartOnNewSecret {
	case "always":
		if s.procStarted {
			// process is running, need to kill it first
			s.logger.Info("stopping process", "process_id", s.proc.Pid())
			s.proc.Stop()
		}
	case "never":
		if s.procStarted {
			s.logger.Info("detected update, but not restarting process", "process_id", s.proc.Pid())
			return nil
		}
	}

	childInput := &child.NewInput{
		Stdin:        os.Stdin,
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
		Command:      s.config.AgentConfig.Exec.Command,
		Args:         s.config.AgentConfig.Exec.Args,
		Timeout:      0,
		Env:          append(os.Environ(), envsToList(newEnvVars)...),
		ReloadSignal: nil,
		KillSignal:   s.config.AgentConfig.Exec.RestartKillSignal,
		KillTimeout:  30 * time.Second,
		Splay:        0,
		Setpgid:      true,
		Logger:       s.logger.StandardLogger(nil),
	}

	proc, err := child.New(childInput)
	if err != nil {
		return err
	}
	s.proc = proc
	s.exitCh = s.proc.ExitCh()

	if err := s.proc.Start(); err != nil {
		return fmt.Errorf("error starting child process: %w", err)
	}
	s.procStarted = true

	return nil
}

func (s *Server) Stop() {
	if s.stopped.CompareAndSwap(false, true) {
		s.procLock.Lock()
		defer s.procLock.Unlock()
		if s.procStarted {
			s.proc.Stop()
		}
		s.procStarted = false
	}
}

func envsToList(envs map[string]string) []string {
	out := make([]string, len(envs))
	for key, value := range envs {
		e := fmt.Sprintf("%s=%s", key, value)
		out = append(out, e)
	}
	return out
}
