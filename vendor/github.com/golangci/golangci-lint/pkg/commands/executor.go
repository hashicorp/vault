package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goutil"
	"github.com/golangci/golangci-lint/pkg/lint"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
)

type Executor struct {
	rootCmd *cobra.Command
	runCmd  *cobra.Command

	exitCode              int
	version, commit, date string

	cfg               *config.Config
	log               logutils.Log
	reportData        report.Data
	DBManager         *lintersdb.Manager
	EnabledLintersSet *lintersdb.EnabledSet
	contextLoader     *lint.ContextLoader
	goenv             *goutil.Env
}

func NewExecutor(version, commit, date string) *Executor {
	e := &Executor{
		cfg:       config.NewDefault(),
		version:   version,
		commit:    commit,
		date:      date,
		DBManager: lintersdb.NewManager(),
	}

	e.log = report.NewLogWrapper(logutils.NewStderrLog(""), &e.reportData)

	// to setup log level early we need to parse config from command line extra time to
	// find `-v` option
	commandLineCfg, err := e.getConfigForCommandLine()
	if err != nil && err != pflag.ErrHelp {
		e.log.Fatalf("Can't get config for command line: %s", err)
	}
	if commandLineCfg != nil {
		logutils.SetupVerboseLog(e.log, commandLineCfg.Run.IsVerbose)
	}

	// init of commands must be done before config file reading because
	// init sets config with the default values of flags
	e.initRoot()
	e.initRun()
	e.initHelp()
	e.initLinters()
	e.initConfig()

	// init e.cfg by values from config: flags parse will see these values
	// like the default ones. It will overwrite them only if the same option
	// is found in command-line: it's ok, command-line has higher priority.

	r := config.NewFileReader(e.cfg, commandLineCfg, e.log.Child("config_reader"))
	if err := r.Read(); err != nil {
		e.log.Fatalf("Can't read config: %s", err)
	}

	e.cfg.LintersSettings.Gocritic.InferEnabledChecks(e.log)
	if err := e.cfg.LintersSettings.Gocritic.Validate(e.log); err != nil {
		e.log.Fatalf("Invalid gocritic settings: %s", err)
	}

	// Slice options must be explicitly set for proper merging of config and command-line options.
	fixSlicesFlags(e.runCmd.Flags())

	e.EnabledLintersSet = lintersdb.NewEnabledSet(e.DBManager,
		lintersdb.NewValidator(e.DBManager), e.log.Child("lintersdb"), e.cfg)
	e.goenv = goutil.NewEnv(e.log.Child("goenv"))
	e.contextLoader = lint.NewContextLoader(e.cfg, e.log.Child("loader"), e.goenv)

	return e
}

func (e *Executor) Execute() error {
	return e.rootCmd.Execute()
}
