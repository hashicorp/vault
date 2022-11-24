package command

import (
	"flag"
	"os"
	"strings"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/posener/complete"
)

// logFlags are the 'log' related flags that can be shared across commands.
type logFlags struct {
	flagCombineLogs       bool
	flagLogLevel          string
	flagLogFormat         string
	flagLogFile           string
	flagLogRotateBytes    string
	flagLogRotateDuration string
	flagLogRotateMaxFiles string
}

type provider = func(key string) (string, bool)

// valuesProvider has the intention of providing a way to supply a func with a
// way to retrieve values for flags and environment variables without having to
// directly call a specific implementation. The reasoning for its existence is
// to facilitate testing.
type valuesProvider struct {
	flagProvider   provider
	envVarProvider provider
}

// addLogFlags will add the set of 'log' related flags to a flag set.
func (f *FlagSet) addLogFlags(l *logFlags) {
	f.BoolVar(&BoolVar{
		Name:    flagNameCombineLogs,
		Target:  &l.flagCombineLogs,
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

	f.StringVar(&StringVar{
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

	f.StringVar(&StringVar{
		Name:   flagNameLogRotateMaxFiles,
		Target: &l.flagLogRotateMaxFiles,
		Usage:  "The maximum number of older log file archives to keep",
	})
}

// getValue will attempt to find the flag with the corresponding flag name (key)
// and return the value along with a bool representing whether of not the flag had been found/set.
func (f *FlagSets) getValue(flagName string) (string, bool) {
	var result string
	var isFlagSpecified bool

	if f != nil {
		f.Visit(func(fl *flag.Flag) {
			if fl.Name == flagName {
				result = fl.Value.String()
				isFlagSpecified = true
			}
		})
	}

	return result, isFlagSpecified
}

// getAggregatedConfigValue uses the provided keys to check CLI flags and environment
// variables for values that may be used to override any specified configuration.
// If nothing can be found in flags/env vars or config, the 'fallback' (default) value will be provided.
func (p *valuesProvider) getAggregatedConfigValue(flagKey, envVarKey, current, fallback string) string {
	var result string
	current = strings.TrimSpace(current)

	flg, flgFound := p.flagProvider(flagKey)
	env, envFound := p.envVarProvider(envVarKey)

	switch {
	case flgFound:
		result = flg
	case envFound:
		// Use value from env var
		result = env
	case current != "":
		// Use value from config
		result = current
	default:
		// Use the default value
		result = fallback
	}

	return result
}

// updateLogConfig will accept a shared config and specifically attempt to update the 'log' related config keys.
// For each 'log' key we aggregate file config/env vars and CLI flags to select the one with the highest precedence.
// This method mutates the config object passed into it.
func (f *FlagSets) updateLogConfig(config *configutil.SharedConfig) {
	p := &valuesProvider{
		flagProvider: func(key string) (string, bool) { return f.getValue(key) },
		envVarProvider: func(key string) (string, bool) {
			if key == "" {
				return "", false
			}
			return os.LookupEnv(key)
		},
	}

	config.LogLevel = p.getAggregatedConfigValue(flagNameLogLevel, EnvVaultLogLevel, config.LogLevel, "info")
	config.LogFormat = p.getAggregatedConfigValue(flagNameLogFormat, EnvVaultLogFormat, config.LogFormat, "")
	config.LogFile = p.getAggregatedConfigValue(flagNameLogFile, "", config.LogFile, "")
	config.LogRotateDuration = p.getAggregatedConfigValue(flagNameLogRotateDuration, "", config.LogRotateDuration, "")
	config.LogRotateBytes = p.getAggregatedConfigValue(flagNameLogRotateBytes, "", config.LogRotateBytes, "")
	config.LogRotateMaxFiles = p.getAggregatedConfigValue(flagNameLogRotateMaxFiles, "", config.LogRotateMaxFiles, "")
}
