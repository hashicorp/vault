package vault

import (
	"time"

	"github.com/hashicorp/vault/command/server"
)

type DebugConfig struct {
	pprofDisabled           bool
	pprofProfileMaxDuration time.Duration
	pprofTraceMaxDuration   time.Duration
}

// NewDebugConfig takes the values from server.Debug and
// returns a populated DebugConfig.
func NewDebugConfig(debug *server.Debug) *DebugConfig {
	if debug == nil {
		return nil
	}

	return &DebugConfig{
		pprofDisabled:           debug.PprofDisable,
		pprofProfileMaxDuration: debug.PprofProfileMaxDuration,
		pprofTraceMaxDuration:   debug.PprofTraceMaxDuration,
	}
}

func (c *DebugConfig) PprofDisabled() bool {
	return c.pprofDisabled
}

func (c *DebugConfig) PprofProfileMaxDuration() time.Duration {
	return c.pprofProfileMaxDuration
}

func (c *DebugConfig) PprofTraceMaxDuration() time.Duration {
	return c.pprofTraceMaxDuration
}
