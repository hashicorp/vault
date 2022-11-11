package command

import "github.com/posener/complete"

type logFlags struct {
	flagLogLevel          string
	flagLogFormat         string
	flagLogFile           string
	flagLogRotateBytes    string
	flagLogRotateDuration string
	flagLogRotateMaxFiles string
	flagLogSyslog         bool
}

func (f *FlagSet) addLogFlags(l *logFlags) {
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
		EnvVar: EnvVaultLogFile,
		Usage:  "Path to the log file that Vault should use for logging",
	})

	f.StringVar(&StringVar{
		Name:   flagNameLogRotateBytes,
		Target: &l.flagLogRotateBytes,
		EnvVar: EnvVaultLogRotateBytes,
		Usage: "Number of bytes that should be written to a log before it needs to be rotated. " +
			"Unless specified, there is no limit to the number of bytes that can be written to a log file",
	})

	f.StringVar(&StringVar{
		Name:   flagNameLogRotateDuration,
		Target: &l.flagLogRotateDuration,
		EnvVar: EnvVaultLogRotateDuration,
		Usage: "The maximum duration a log should be written to before it needs to be rotated. " +
			"Must be a duration value such as 30s",
	})

	f.StringVar(&StringVar{
		Name:   flagNameLogRotateMaxFiles,
		Target: &l.flagLogRotateMaxFiles,
		EnvVar: EnvVaultLogRotateMaxFiles,
		Usage:  "The maximum number of older log file archives to keep",
	})

	f.BoolVar(&BoolVar{
		Name:   flagNameLogSyslog,
		Target: &l.flagLogSyslog,
		EnvVar: EnvVaultLogSyslog,
		Usage: "Enables logging to syslog. This is only supported on Linux and macOS. " +
			"It will result in an error if provided on Windows",
	})
}
