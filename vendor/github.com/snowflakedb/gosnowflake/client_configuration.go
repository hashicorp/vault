// Copyright (c) 2023 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// log levels for easy logging
const (
	levelOff   string = "OFF"   // log level for logging switched off
	levelError string = "ERROR" // error log level
	levelWarn  string = "WARN"  // warn log level
	levelInfo  string = "INFO"  // info log level
	levelDebug string = "DEBUG" // debug log level
	levelTrace string = "TRACE" // trace log level
)

const (
	defaultConfigName = "sf_client_config.json"
	clientConfEnvName = "SF_CLIENT_CONFIG_FILE"
)

func getClientConfig(filePathFromConnectionString string) (*ClientConfig, string, error) {
	configPredefinedDirPaths := clientConfigPredefinedDirs()
	filePath, err := findClientConfigFilePath(filePathFromConnectionString, configPredefinedDirPaths)
	if err != nil {
		return nil, "", err
	}
	if filePath == "" { // we did not find a config file
		return nil, "", nil
	}
	config, err := parseClientConfiguration(filePath)
	return config, filePath, err
}

func findClientConfigFilePath(filePathFromConnectionString string, configPredefinedDirs []string) (string, error) {
	if filePathFromConnectionString != "" {
		logger.Infof("Using client configuration path from a connection string: %s", filePathFromConnectionString)
		return filePathFromConnectionString, nil
	}
	envConfigFilePath := os.Getenv(clientConfEnvName)
	if envConfigFilePath != "" {
		logger.Infof("Using client configuration path from an environment variable: %s", envConfigFilePath)
		return envConfigFilePath, nil
	}
	return searchForConfigFile(configPredefinedDirs)
}

func searchForConfigFile(directories []string) (string, error) {
	for _, dir := range directories {
		filePath := path.Join(dir, defaultConfigName)
		exists, err := existsFile(filePath)
		if err != nil {
			return "", fmt.Errorf("error while searching for client config in directory: %s, err: %s", dir, err)
		}
		if exists {
			logger.Infof("Using client configuration from a default directory: %s", filePath)
			return filePath, nil
		}
		logger.Debugf("No client config found in directory: %s", dir)
	}
	logger.Info("No client config file found in default directories")
	return "", nil
}

func existsFile(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func clientConfigPredefinedDirs() []string {
	var predefinedDirs []string
	exeFile, err := os.Executable()
	if err != nil {
		logger.Warnf("Unable to access the application directory for client configuration search, err: %v", err)
	} else {
		predefinedDirs = append(predefinedDirs, filepath.Dir(exeFile))
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Warnf("Unable to access Home directory for client configuration search, err: %v", err)
	} else {
		predefinedDirs = append(predefinedDirs, homeDir)
	}
	if predefinedDirs == nil {
		return []string{}
	}
	return predefinedDirs
}

// ClientConfig config root
type ClientConfig struct {
	Common *ClientConfigCommonProps `json:"common"`
}

// ClientConfigCommonProps properties from "common" section
type ClientConfigCommonProps struct {
	LogLevel string `json:"log_level,omitempty"`
	LogPath  string `json:"log_path,omitempty"`
}

func parseClientConfiguration(filePath string) (*ClientConfig, error) {
	if filePath == "" {
		return nil, nil
	}
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, parsingClientConfigError(err)
	}
	err = validateCfgPerm(filePath)
	if err != nil {
		return nil, parsingClientConfigError(err)
	}
	var clientConfig ClientConfig
	err = json.Unmarshal(fileContents, &clientConfig)
	if err != nil {
		return nil, parsingClientConfigError(err)
	}
	unknownValues := getUnknownValues(fileContents)
	if len(unknownValues) > 0 {
		for val := range unknownValues {
			logger.Warnf("Unknown configuration entry: %s with value: %s", val, unknownValues[val])
		}
	}
	err = validateClientConfiguration(&clientConfig)
	if err != nil {
		return nil, parsingClientConfigError(err)
	}
	return &clientConfig, nil
}

func getUnknownValues(fileContents []byte) map[string]interface{} {
	var values map[string]interface{}
	err := json.Unmarshal(fileContents, &values)
	if err != nil {
		return nil
	}
	if values["common"] == nil {
		return nil
	}
	commonValues := values["common"].(map[string]interface{})
	lowercaseCommonValues := make(map[string]interface{}, len(commonValues))
	for k, v := range commonValues {
		lowercaseCommonValues[strings.ToLower(k)] = v
	}
	delete(lowercaseCommonValues, "log_level")
	delete(lowercaseCommonValues, "log_path")
	return lowercaseCommonValues
}

func parsingClientConfigError(err error) error {
	return fmt.Errorf("parsing client config failed: %w", err)
}

func validateClientConfiguration(clientConfig *ClientConfig) error {
	if clientConfig == nil {
		return errors.New("client config not found")
	}
	if clientConfig.Common == nil {
		return errors.New("common section in client config not found")
	}
	return validateLogLevel(*clientConfig)
}

func validateLogLevel(clientConfig ClientConfig) error {
	var logLevel = clientConfig.Common.LogLevel
	if logLevel != "" {
		_, err := toLogLevel(logLevel)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateCfgPerm(filePath string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	stat, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	perm := stat.Mode()
	// Check if group (5th LSB) or others (2nd LSB) have a write permission to the file
	if perm&(1<<4) != 0 || perm&(1<<1) != 0 {
		return fmt.Errorf("configuration file: %s can be modified by group or others", filePath)
	}
	return nil
}

func toLogLevel(logLevelString string) (string, error) {
	var logLevel = strings.ToUpper(logLevelString)
	switch logLevel {
	case levelOff, levelError, levelWarn, levelInfo, levelDebug, levelTrace:
		return logLevel, nil
	default:
		return "", errors.New("unknown log level: " + logLevelString)
	}
}
