package config

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

const (
	// DefaultExecKillSignal is the default signal to send to the process to
	// tell it to gracefully terminate.
	DefaultExecKillSignal = syscall.SIGINT

	// DefaultExecKillTimeout is the maximum amount of time to wait for the
	// process to gracefully terminate before force-killing it.
	DefaultExecKillTimeout = 30 * time.Second

	// DefaultExecTimeout is the default amount of time to wait for a
	// command to exit. By default, this is disabled, which means the command
	// is allowed to run for an infinite amount of time.
	DefaultExecTimeout = 0 * time.Second
)

var (
	// DefaultExecReloadSignal is the default signal to send to the process to
	// tell it to reload its configuration.
	DefaultExecReloadSignal = (os.Signal)(nil)
)

// ExecConfig is used to configure the application when it runs in
// exec/supervise mode.
type ExecConfig struct {
	// Command is the command to execute and watch as a child process.
	Command *string `mapstructure:"command"`

	// Enabled controls if this exec is enabled.
	Enabled *bool `mapstructure:"enabled"`

	// EnvConfig is the environmental customizations.
	Env *EnvConfig `mapstructure:"env"`

	// KillSignal is the signal to send to the command to kill it gracefully.
	KillSignal *os.Signal `mapstructure:"kill_signal"`

	// KillTimeout is the amount of time to give the process to cleanup before
	// hard-killing it.
	KillTimeout *time.Duration `mapstructure:"kill_timeout"`

	// ReloadSignal is the signal to send to the child process when a template
	// changes. This tells the child process that templates have
	ReloadSignal *os.Signal `mapstructure:"reload_signal"`

	// Splay is the maximum amount of random time to wait to signal or kill the
	// process. By default this is disabled, but it can be set to low values to
	// reduce the "thundering herd" problem where all tasks are restarted at once.
	Splay *time.Duration `mapstructure:"splay"`

	// Timeout is the maximum amount of time to wait for a command to complete.
	// By default, this is 0, which means "wait forever".
	Timeout *time.Duration `mapstructure:"timeout"`
}

// DefaultExecConfig returns a configuration that is populated with the
// default values.
func DefaultExecConfig() *ExecConfig {
	return &ExecConfig{
		Env: DefaultEnvConfig(),
	}
}

// Copy returns a deep copy of this configuration.
func (c *ExecConfig) Copy() *ExecConfig {
	if c == nil {
		return nil
	}

	var o ExecConfig

	o.Command = c.Command

	o.Enabled = c.Enabled

	if c.Env != nil {
		o.Env = c.Env.Copy()
	}

	o.KillSignal = c.KillSignal

	o.KillTimeout = c.KillTimeout

	o.ReloadSignal = c.ReloadSignal

	o.Splay = c.Splay

	o.Timeout = c.Timeout

	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *ExecConfig) Merge(o *ExecConfig) *ExecConfig {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	if o.Command != nil {
		r.Command = o.Command
	}

	if o.Enabled != nil {
		r.Enabled = o.Enabled
	}

	if o.Env != nil {
		r.Env = r.Env.Merge(o.Env)
	}

	if o.KillSignal != nil {
		r.KillSignal = o.KillSignal
	}

	if o.KillTimeout != nil {
		r.KillTimeout = o.KillTimeout
	}

	if o.ReloadSignal != nil {
		r.ReloadSignal = o.ReloadSignal
	}

	if o.Splay != nil {
		r.Splay = o.Splay
	}

	if o.Timeout != nil {
		r.Timeout = o.Timeout
	}

	return r
}

// Finalize ensures there no nil pointers.
func (c *ExecConfig) Finalize() {
	if c.Enabled == nil {
		c.Enabled = Bool(StringPresent(c.Command))
	}

	if c.Command == nil {
		c.Command = String("")
	}

	if c.Env == nil {
		c.Env = DefaultEnvConfig()
	}
	c.Env.Finalize()

	if c.KillSignal == nil {
		c.KillSignal = Signal(DefaultExecKillSignal)
	}

	if c.KillTimeout == nil {
		c.KillTimeout = TimeDuration(DefaultExecKillTimeout)
	}

	if c.ReloadSignal == nil {
		c.ReloadSignal = Signal(DefaultExecReloadSignal)
	}

	if c.Splay == nil {
		c.Splay = TimeDuration(0 * time.Second)
	}

	if c.Timeout == nil {
		c.Timeout = TimeDuration(DefaultExecTimeout)
	}
}

// GoString defines the printable version of this struct.
func (c *ExecConfig) GoString() string {
	if c == nil {
		return "(*ExecConfig)(nil)"
	}

	return fmt.Sprintf("&ExecConfig{"+
		"Command:%s, "+
		"Enabled:%s, "+
		"Env:%#v, "+
		"KillSignal:%s, "+
		"KillTimeout:%s, "+
		"ReloadSignal:%s, "+
		"Splay:%s, "+
		"Timeout:%s"+
		"}",
		StringGoString(c.Command),
		BoolGoString(c.Enabled),
		c.Env,
		SignalGoString(c.KillSignal),
		TimeDurationGoString(c.KillTimeout),
		SignalGoString(c.ReloadSignal),
		TimeDurationGoString(c.Splay),
		TimeDurationGoString(c.Timeout),
	)
}
