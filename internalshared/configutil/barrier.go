package configutil

import (
	"time"
)

// 10% shy of the NIST recommended maximum, leaving a buffer to account for
// tracking losses.
const AbsoluteOperationMaximum = int64(3865470566)

var DefaultRotationConfig = KeyRotationConfig{
	MaxOperations: AbsoluteOperationMaximum,
}

type KeyRotationConfig struct {
	MaxOperations int64
	Interval      time.Duration
	nextRotation             time.Time
}

func (c *KeyRotationConfig) Sanitize() {
	if c.MaxOperations == 0 || c.MaxOperations > AbsoluteOperationMaximum {
		c.MaxOperations = AbsoluteOperationMaximum
	}
}

func (c *KeyRotationConfig) Equals(config KeyRotationConfig) bool {
	return c.MaxOperations == config.MaxOperations && c.Interval == config.Interval
}
