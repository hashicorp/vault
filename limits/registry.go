// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1
package limits

import (
	"os"
	"strconv"
	"sync"

	"github.com/hashicorp/go-hclog"
)

const (
	WriteLimiter         = "write"
	SpecialPathLimiter   = "special-path"
	LimitsBadEnvVariable = "failed to process limiter environment variable, using default"
)

// NOTE: Great care should be taken when setting any of these variables to avoid
// adverse affects in optimal request servicing. It is strongly advised that
// these variables not be used unless there is a very good reason. These are
// intentionally undocumented environment variables that may be removed in
// future versions of Vault.
const (
	// EnvVaultDisableWriteLimiter is used to turn off the
	// RequestLimiter for write-based HTTP methods.
	EnvVaultDisableWriteLimiter = "VAULT_DISABLE_WRITE_LIMITER"

	// EnvVaultWriteLimiterMin is used to modify the minimum
	// concurrency limit for write-based HTTP methods.
	EnvVaultWriteLimiterMin = "VAULT_WRITE_LIMITER_MIN"

	// EnvVaultWriteLimiterMax is used to modify the maximum
	// concurrency limit for write-based HTTP methods.
	EnvVaultWriteLimiterMax = "VAULT_WRITE_LIMITER_MAX"

	// EnvVaultDisablePathBasedRequestLimiting is used to turn off the
	// RequestLimiter for special-cased paths, specified in
	// Backend.PathsSpecial.
	EnvVaultDisableSpecialPathLimiter = "VAULT_DISABLE_SPECIAL_PATH_LIMITER"

	// EnvVaultSpecialPathLimiterMin is used to modify the minimum
	// concurrency limit for write-based HTTP methods.
	EnvVaultSpecialPathLimiterMin = "VAULT_SPECIAL_PATH_LIMITER_MIN"

	// EnvVaultSpecialPathLimiterMax is used to modify the maximum
	// concurrency limit for write-based HTTP methods.
	EnvVaultSpecialPathLimiterMax = "VAULT_SPECIAL_PATH_LIMITER_MAX"
)

// LimiterRegistry holds the map of RequestLimiters mapped to keys.
type LimiterRegistry struct {
	Limiters map[string]*RequestLimiter
	Logger   hclog.Logger
	Enabled  bool
	sync.RWMutex
}

// NewLimiterRegistry is a basic LimiterRegistry constructor.
func NewLimiterRegistry(logger hclog.Logger) *LimiterRegistry {
	return &LimiterRegistry{
		Limiters: make(map[string]*RequestLimiter),
		Logger:   logger,
	}
}

// processEnvVars consults Limiter-specific environment variables and tells the
// caller if the Limiter should be disabled. If not, it adjusts the passed-in
// limiterFlags as appropriate.
func (r *LimiterRegistry) processEnvVars(name string, flags *LimiterFlags, envDisabled, envMin, envMax string) bool {
	envFlagsLogger := r.Logger.With("name", name)
	if disabledRaw := os.Getenv(envDisabled); disabledRaw != "" {
		disabled, err := strconv.ParseBool(disabledRaw)
		if err != nil {
			envFlagsLogger.Warn(LimitsBadEnvVariable,
				"env", envDisabled,
				"val", disabledRaw,
				"default", false,
				"error", err,
			)
		}

		if disabled {
			envFlagsLogger.Warn("limiter disabled by environment variable", "env", envDisabled, "val", disabledRaw)
			return true
		}
	}

	envFlags := &LimiterFlags{}
	if minRaw := os.Getenv(envMin); minRaw != "" {
		min, err := strconv.Atoi(minRaw)
		if err != nil {
			envFlagsLogger.Warn(LimitsBadEnvVariable,
				"env", envMin,
				"val", minRaw,
				"default", flags.MinLimit,
				"error", err,
			)
		} else {
			envFlags.MinLimit = min
		}
	}

	if maxRaw := os.Getenv(envMax); maxRaw != "" {
		max, err := strconv.Atoi(maxRaw)
		if err != nil {
			envFlagsLogger.Warn(LimitsBadEnvVariable,
				"env", envMax,
				"val", maxRaw,
				"default", flags.MaxLimit,
				"error", err,
			)
		} else {
			envFlags.MaxLimit = max
		}
	}

	switch {
	case envFlags.MinLimit == 0:
		// Assume no environment variable was provided.
	case envFlags.MinLimit > 0:
		flags.MinLimit = envFlags.MinLimit
	default:
		r.Logger.Warn("min limit must be greater than zero, falling back to defaults", "minLimit", flags.MinLimit)
	}

	switch {
	case envFlags.MaxLimit == 0:
		// Assume no environment variable was provided.
	case envFlags.MaxLimit > flags.MinLimit:
		flags.MaxLimit = envFlags.MaxLimit
	default:
		r.Logger.Warn("max limit must be greater than min, falling back to defaults", "maxLimit", flags.MaxLimit)
	}

	return false
}

// Enable sets up a new LimiterRegistry and marks it Enabled.
func (r *LimiterRegistry) Enable() {
	r.Lock()
	defer r.Unlock()

	if r.Enabled {
		return
	}

	r.Logger.Info("enabling request limiters")
	r.Limiters = map[string]*RequestLimiter{}

	for name, flags := range DefaultLimiterFlags {
		r.Register(name, flags)
	}

	r.Enabled = true
}

// Register creates a new request limiter and assigns it a slot in the
// LimiterRegistry. Locking should be done in the caller.
func (r *LimiterRegistry) Register(name string, flags LimiterFlags) {
	var disabled bool

	switch name {
	case WriteLimiter:
		disabled = r.processEnvVars(name, &flags,
			EnvVaultDisableWriteLimiter,
			EnvVaultWriteLimiterMin,
			EnvVaultWriteLimiterMax,
		)
		if disabled {
			return
		}
	case SpecialPathLimiter:
		disabled = r.processEnvVars(name, &flags,
			EnvVaultDisableSpecialPathLimiter,
			EnvVaultSpecialPathLimiterMin,
			EnvVaultSpecialPathLimiterMax,
		)
		if disabled {
			return
		}
	default:
		r.Logger.Warn("skipping invalid limiter type", "key", name)
		return
	}

	// Always set the initial limit to min so the system can find its own
	// equilibrium, since max might be too high.
	flags.InitialLimit = flags.MinLimit

	limiter, err := NewRequestLimiter(r.Logger.Named(name), name, flags)
	if err != nil {
		r.Logger.Error("failed to register limiter", "name", name, "error", err)
		return
	}

	r.Limiters[name] = limiter
}

// Disable drops its references to underlying limiters.
func (r *LimiterRegistry) Disable() {
	r.Lock()
	defer r.Unlock()

	if !r.Enabled {
		return
	}

	r.Logger.Info("disabling request limiters")
	// Any outstanding tokens will be flushed when their request completes, as
	// they've already acquired a listener. Just drop the limiter references
	// here and the garbage-collector should take care of the rest.
	r.Limiters = map[string]*RequestLimiter{}
	r.Enabled = false
}

// GetLimiter looks up a RequestLimiter by key in the LimiterRegistry.
func (r *LimiterRegistry) GetLimiter(key string) *RequestLimiter {
	r.RLock()
	defer r.RUnlock()
	return r.Limiters[key]
}
