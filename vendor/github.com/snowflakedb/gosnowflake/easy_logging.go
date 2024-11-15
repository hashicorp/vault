package gosnowflake

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
)

type initTrials struct {
	everTriedToInitialize bool
	clientConfigFileInput string
	configureCounter      int
	mu                    sync.Mutex
}

var easyLoggingInitTrials = initTrials{
	everTriedToInitialize: false,
	clientConfigFileInput: "",
	configureCounter:      0,
	mu:                    sync.Mutex{},
}

func (i *initTrials) setInitTrial(clientConfigFileInput string) {
	i.everTriedToInitialize = true
	i.clientConfigFileInput = clientConfigFileInput
}

func (i *initTrials) increaseReconfigureCounter() {
	i.configureCounter++
}

func (i *initTrials) reset() {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.everTriedToInitialize = false
	i.clientConfigFileInput = ""
	i.configureCounter = 0
}

func initEasyLogging(clientConfigFileInput string) error {
	easyLoggingInitTrials.mu.Lock()
	defer easyLoggingInitTrials.mu.Unlock()

	if !allowedToInitialize(clientConfigFileInput) {
		logger.Info("Skipping Easy Logging initialization as it is not allowed to initialize")
		return nil
	}
	logger.Infof("Trying to initialize Easy Logging")
	config, configPath, err := getClientConfig(clientConfigFileInput)
	if err != nil {
		logger.Errorf("Failed to initialize Easy Logging, err: %s", err)
		return easyLoggingInitError(err)
	}
	if config == nil {
		logger.Info("Easy Logging is disabled as no config has been found")
		easyLoggingInitTrials.setInitTrial(clientConfigFileInput)
		return nil
	}
	var logLevel string
	logLevel, err = getLogLevel(config.Common.LogLevel)
	if err != nil {
		logger.Errorf("Failed to initialize Easy Logging, err: %s", err)
		return easyLoggingInitError(err)
	}
	var logPath string
	logPath, err = getLogPath(config.Common.LogPath)
	if err != nil {
		logger.Errorf("Failed to initialize Easy Logging, err: %s", err)
		return easyLoggingInitError(err)
	}
	logger.Infof("Initializing Easy Logging with logPath=%s and logLevel=%s from file: %s", logPath, logLevel, configPath)
	err = reconfigureEasyLogging(logLevel, logPath)
	if err != nil {
		logger.Errorf("Failed to initialize Easy Logging, err: %s", err)
	}
	easyLoggingInitTrials.setInitTrial(clientConfigFileInput)
	easyLoggingInitTrials.increaseReconfigureCounter()
	return err
}

func easyLoggingInitError(err error) error {
	return &SnowflakeError{
		Number:      ErrCodeClientConfigFailed,
		Message:     errMsgClientConfigFailed,
		MessageArgs: []interface{}{err.Error()},
	}
}

func reconfigureEasyLogging(logLevel string, logPath string) error {
	newLogger := CreateDefaultLogger()
	err := newLogger.SetLogLevel(logLevel)
	if err != nil {
		return err
	}
	var output io.Writer
	var file *os.File
	output, file, err = createLogWriter(logPath)
	if err != nil {
		return err
	}
	newLogger.SetOutput(output)
	err = newLogger.CloseFileOnLoggerReplace(file)
	if err != nil {
		logger.Errorf("%s", err)
	}
	logger.Replace(&newLogger)
	return nil
}

func createLogWriter(logPath string) (io.Writer, *os.File, error) {
	if strings.EqualFold(logPath, "STDOUT") {
		return os.Stdout, nil, nil
	}
	logFileName := path.Join(logPath, "snowflake.log")
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return nil, nil, err
	}
	return file, file, nil
}

func allowedToInitialize(clientConfigFileInput string) bool {
	triedToInitializeWithoutConfigFile := easyLoggingInitTrials.everTriedToInitialize && easyLoggingInitTrials.clientConfigFileInput == ""
	isAllowedToInitialize := !easyLoggingInitTrials.everTriedToInitialize || (triedToInitializeWithoutConfigFile && clientConfigFileInput != "")
	if !isAllowedToInitialize && easyLoggingInitTrials.clientConfigFileInput != clientConfigFileInput {
		logger.Warnf("Easy logging will not be configured for CLIENT_CONFIG_FILE=%s because it was previously configured for a different client config", clientConfigFileInput)
	}
	return isAllowedToInitialize
}

func getLogLevel(logLevel string) (string, error) {
	if logLevel == "" {
		logger.Warn("LogLevel in client config not found. Using default value: OFF")
		return levelOff, nil
	}
	return toLogLevel(logLevel)
}

func getLogPath(logPath string) (string, error) {
	logPathOrDefault := logPath
	if logPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("user home directory is not accessible, err: %w", err)
		}
		logPathOrDefault = homeDir
		logger.Warnf("LogPath in client config not found. Using user home directory as a default value: %s", logPathOrDefault)
	}
	pathWithGoSubdir := path.Join(logPathOrDefault, "go")
	exists, err := dirExists(pathWithGoSubdir)
	if err != nil {
		return "", err
	}
	if !exists {
		err = os.MkdirAll(pathWithGoSubdir, 0700)
		if err != nil {
			return "", err
		}
	}
	logDirPermValid, perm, err := isDirAccessCorrect(pathWithGoSubdir)
	if err != nil {
		return "", err
	}
	if !logDirPermValid {
		logger.Warnf("Log directory: %s could potentially be accessed by others. Directory chmod: 0%o", pathWithGoSubdir, *perm)
	}
	return pathWithGoSubdir, nil
}

func isDirAccessCorrect(dirPath string) (bool, *os.FileMode, error) {
	if runtime.GOOS == "windows" {
		return true, nil, nil
	}
	dirStat, err := os.Stat(dirPath)
	if err != nil {
		return false, nil, err
	}
	perm := dirStat.Mode().Perm()
	if perm != 0700 {
		return false, &perm, nil
	}
	return true, &perm, nil
}

func dirExists(dirPath string) (bool, error) {
	stat, err := os.Stat(dirPath)
	if err == nil {
		return stat.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
