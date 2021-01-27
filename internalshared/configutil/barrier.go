package configutil

import (
	"time"
)

// 10% shy of the NIST recommended maximum, leaving a buffer to account for
// tracking losses.
const AbsoluteOperationMaximum = int64(3865470566)

var DefaultRotationConfig = KeyRotationConfig{
	KeyRotationMaxOperations: AbsoluteOperationMaximum,
}

type KeyRotationConfig struct {
	KeyRotationMaxOperations int64 `hcl:"key_rotation_max_operations"`
	KeyRotationInterval      time.Duration
	KeyRotationIntervalRaw   interface{} `hcl:"key_rotation_interval"`
}

func (c *KeyRotationConfig) Sanitize() {
	if c.KeyRotationMaxOperations == 0 || c.KeyRotationMaxOperations > AbsoluteOperationMaximum {
		c.KeyRotationMaxOperations = AbsoluteOperationMaximum
	}
}

func (c *KeyRotationConfig) Equals(config KeyRotationConfig) bool {
	return c.KeyRotationMaxOperations == config.KeyRotationMaxOperations && c.KeyRotationInterval == c.KeyRotationInterval
}
