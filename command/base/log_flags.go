// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"flag"
	"os"
	"strconv"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/posener/complete"
)

// LogFlags are the 'log' related flags that can be shared across commands.
type LogFlags struct {
	FlagCombineLogs       bool
	FlagDisableGatedLogs  bool
	FlagLogLevel          string
	FlagLogFormat         string
	FlagLogFile           string
	FlagLogRotateBytes    int
	FlagLogRotateDuration string
	FlagLogRotateMaxFiles int
}

// valuesProvider has the intention of providing a way to supply a func with a
// way to retrieve values for flags and environment variables without having to
// directly call a specific implementation.
// The reasoning for its existence is to facilitate testing.
type valuesProvider struct {
	flagProvider   func(string) (flag.Value, bool)
	envVarProvider func(string) (string, bool)
}

// AddLogFlags will add the set of 'log' related flags to a flag set.
func (f *FlagSet) AddLogFlags(l *LogFlags) {
	f.BoolVar(&BoolVar{
		Name:    FlagNameCombineLogs,
		Target:  &l.FlagCombineLogs,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    FlagNameDisableGatedLogs,
		Target:  &l.FlagDisableGatedLogs,
		Default: false,
		Hidden:  true,
	})

	f.StringVar(&StringVar{
		Name:       FlagNameLogLevel,
		Target:     &l.FlagLogLevel,
		Default:    notSetValue,
		EnvVar:     EnvVaultLogLevel,
		Completion: complete.PredictSet("trace", "debug", "info", "warn", "error"),
		Usage: "Log verbosity level. Supported values (in order of detail) are " +
			"\"trace\", \"debug\", \"info\", \"warn\", and \"error\".",
	})

	f.StringVar(&StringVar{
		Name:       FlagNameLogFormat,
		Target:     &l.FlagLogFormat,
		Default:    notSetValue,
		EnvVar:     EnvVaultLogFormat,
		Completion: complete.PredictSet("standard", "json"),
		Usage:      `Log format. Supported values are "standard" and "json".`,
	})

	f.StringVar(&StringVar{
		Name:   FlagNameLogFile,
		Target: &l.FlagLogFile,
		Usage:  "Path to the log file that Vault should use for logging",
	})

	f.IntVar(&IntVar{
		Name:   FlagNameLogRotateBytes,
		Target: &l.FlagLogRotateBytes,
		Usage: "Number of bytes that should be written to a log before it needs to be rotated. " +
			"Unless specified, there is no limit to the number of bytes that can be written to a log file",
	})

	f.StringVar(&StringVar{
		Name:   FlagNameLogRotateDuration,
		Target: &l.FlagLogRotateDuration,
		Usage: "The maximum duration a log should be written to before it needs to be rotated. " +
			"Must be a duration value such as 30s",
	})

	f.IntVar(&IntVar{
		Name:   FlagNameLogRotateMaxFiles,
		Target: &l.FlagLogRotateMaxFiles,
		Usage:  "The maximum number of older log file archives to keep",
	})
}

// envVarValue attempts to get a named value from the environment variables.
// The value will be returned as a string along with a boolean value indiciating
// to the caller whether the named env var existed.
func envVarValue(key string) (string, bool) {
	if key == "" {
		return "", false
	}
	return os.LookupEnv(key)
}

// flagValue attempts to find the named flag in a set of FlagSets.
// The flag.Value is returned if it was specified, and the boolean value indicates
// to the caller if the flag was specified by the end user.
func (f *FlagSets) flagValue(flagName string) (flag.Value, bool) {
	var result flag.Value
	var isFlagSpecified bool

	if f != nil {
		f.Visit(func(fl *flag.Flag) {
			if fl.Name == flagName {
				result = fl.Value
				isFlagSpecified = true
			}
		})
	}

	return result, isFlagSpecified
}

// overrideValue uses the provided keys to check CLI flags and environment
// variables for values that may be used to override any specified configuration.
func (p *valuesProvider) overrideValue(flagKey, envVarKey string) (string, bool) {
	var result string
	found := true

	flg, flgFound := p.flagProvider(flagKey)
	env, envFound := p.envVarProvider(envVarKey)

	switch {
	case flgFound:
		result = flg.String()
	case envFound:
		result = env
	default:
		found = false
	}

	return result, found
}

// ApplyLogConfigOverrides will accept a shared config and specifically attempt to update the 'log' related config keys.
// For each 'log' key, we aggregate file config, env vars and CLI flags to select the one with the highest precedence.
// This method mutates the config object passed into it.
func (f *FlagSets) ApplyLogConfigOverrides(config *configutil.SharedConfig) {
	p := &valuesProvider{
		flagProvider:   f.flagValue,
		envVarProvider: envVarValue,
	}

	// Update log level
	if val, found := p.overrideValue(FlagNameLogLevel, EnvVaultLogLevel); found {
		config.LogLevel = val
	}

	// Update log format
	if val, found := p.overrideValue(FlagNameLogFormat, EnvVaultLogFormat); found {
		config.LogFormat = val
	}

	// Update log file name
	if val, found := p.overrideValue(FlagNameLogFile, ""); found {
		config.LogFile = val
	}

	// Update log rotation duration
	if val, found := p.overrideValue(FlagNameLogRotateDuration, ""); found {
		config.LogRotateDuration = val
	}

	// Update log max files
	if val, found := p.overrideValue(FlagNameLogRotateMaxFiles, ""); found {
		config.LogRotateMaxFiles, _ = strconv.Atoi(val)
	}

	// Update log rotation max bytes
	if val, found := p.overrideValue(FlagNameLogRotateBytes, ""); found {
		config.LogRotateBytes, _ = strconv.Atoi(val)
	}
}
