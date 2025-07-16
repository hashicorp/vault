// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"flag"
	"os"
	"strconv"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/posener/complete"
)

// logFlags are the 'log' related flags that can be shared across commands.
type logFlags struct {
	flagCombineLogs       bool
	flagDisableGatedLogs  bool
	flagLogLevel          string
	flagLogFormat         string
	flagLogFile           string
	flagLogRotateBytes    int
	flagLogRotateDuration string
	flagLogRotateMaxFiles int
}

// valuesProvider has the intention of providing a way to supply a func with a
// way to retrieve values for flags and environment variables without having to
// directly call a specific implementation.
// The reasoning for its existence is to facilitate testing.
type valuesProvider struct {
	flagProvider   func(string) (flag.Value, bool)
	envVarProvider func(string) (string, bool)
}

// addLogFlags will add the set of 'log' related flags to a flag set.
func (f *FlagSet) addLogFlags(l *logFlags) {
	f.BoolVar(&BoolVar{
		Name:    flagNameCombineLogs,
		Target:  &l.flagCombineLogs,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    flagDisableGatedLogs,
		Target:  &l.flagDisableGatedLogs,
		Default: false,
		Hidden:  true,
	})

	f.StringVar(&StringVar{
		Name:       flagNameLogLevel,
		Target:     &l.flagLogLevel,
		Default:    notSetValue,
		EnvVar:     EnvVaultLogLevel,
		Completion: complete.PredictSet("trace", "debug", "info", "warn", "error"),
		Usage: "Log verbosity level. Supported values (in order of detail) are " +
			"\"trace\", \"debug\", \"info\", \"warn\", and \"error\".",
	})

	f.StringVar(&StringVar{
		Name:       flagNameLogFormat,
		Target:     &l.flagLogFormat,
		Default:    notSetValue,
		EnvVar:     EnvVaultLogFormat,
		Completion: complete.PredictSet("standard", "json"),
		Usage:      `Log format. Supported values are "standard" and "json".`,
	})

	f.StringVar(&StringVar{
		Name:   flagNameLogFile,
		Target: &l.flagLogFile,
		Usage:  "Path to the log file that Vault should use for logging",
	})

	f.IntVar(&IntVar{
		Name:   flagNameLogRotateBytes,
		Target: &l.flagLogRotateBytes,
		Usage: "Number of bytes that should be written to a log before it needs to be rotated. " +
			"Unless specified, there is no limit to the number of bytes that can be written to a log file",
	})

	f.StringVar(&StringVar{
		Name:   flagNameLogRotateDuration,
		Target: &l.flagLogRotateDuration,
		Usage: "The maximum duration a log should be written to before it needs to be rotated. " +
			"Must be a duration value such as 30s",
	})

	f.IntVar(&IntVar{
		Name:   flagNameLogRotateMaxFiles,
		Target: &l.flagLogRotateMaxFiles,
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

// applyLogConfigOverrides will accept a shared config and specifically attempt to update the 'log' related config keys.
// For each 'log' key, we aggregate file config, env vars and CLI flags to select the one with the highest precedence.
// This method mutates the config object passed into it.
func (f *FlagSets) applyLogConfigOverrides(config *configutil.SharedConfig) {
	p := &valuesProvider{
		flagProvider:   f.flagValue,
		envVarProvider: envVarValue,
	}

	// Update log level
	if val, found := p.overrideValue(flagNameLogLevel, EnvVaultLogLevel); found {
		config.LogLevel = val
	}

	// Update log format
	if val, found := p.overrideValue(flagNameLogFormat, EnvVaultLogFormat); found {
		config.LogFormat = val
	}

	// Update log file name
	if val, found := p.overrideValue(flagNameLogFile, ""); found {
		config.LogFile = val
	}

	// Update log rotation duration
	if val, found := p.overrideValue(flagNameLogRotateDuration, ""); found {
		config.LogRotateDuration = val
	}

	// Update log max files
	if val, found := p.overrideValue(flagNameLogRotateMaxFiles, ""); found {
		config.LogRotateMaxFiles, _ = strconv.Atoi(val)
	}

	// Update log rotation max bytes
	if val, found := p.overrideValue(flagNameLogRotateBytes, ""); found {
		config.LogRotateBytes, _ = strconv.Atoi(val)
	}
}
