// Copyright (c) 2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/99designs/keyring"
)

const (
	driverName        = "SNOWFLAKE-GO-DRIVER"
	credCacheDirEnv   = "SF_TEMPORARY_CREDENTIAL_CACHE_DIR"
	credCacheFileName = "temporary_credential.json"
)

var (
	credCacheDir   = ""
	credCache      = ""
	localCredCache = map[string]string{}
)

var (
	credCacheLock sync.RWMutex
)

func createCredentialCacheDir() {
	credCacheDir = os.Getenv(credCacheDirEnv)
	if credCacheDir == "" {
		switch runtime.GOOS {
		case "windows":
			credCacheDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Snowflake", "Caches")
		case "darwin":
			home := os.Getenv("HOME")
			if home == "" {
				logger.Info("HOME is blank.")
			}
			credCacheDir = filepath.Join(home, "Library", "Caches", "Snowflake")
		default:
			home := os.Getenv("HOME")
			if home == "" {
				logger.Info("HOME is blank")
			}
			credCacheDir = filepath.Join(home, ".cache", "snowflake")
		}
	}

	if _, err := os.Stat(credCacheDir); os.IsNotExist(err) {
		if err = os.MkdirAll(credCacheDir, os.ModePerm); err != nil {
			logger.Debugf("Failed to create cache directory. %v, err: %v. ignored\n", credCacheDir, err)
		}
	}
	credCache = filepath.Join(credCacheDir, credCacheFileName)
	logger.Infof("Cache directory: %v", credCache)
}

func setCredential(sc *snowflakeConn, credType, token string) {
	if token == "" {
		logger.Debug("no token provided")
	} else {
		var target string
		if runtime.GOOS == "windows" {
			target = driverName + ":" + credType
			ring, _ := keyring.Open(keyring.Config{
				WinCredPrefix: strings.ToUpper(sc.cfg.Host),
				ServiceName:   strings.ToUpper(sc.cfg.User),
			})
			item := keyring.Item{
				Key:  target,
				Data: []byte(token),
			}
			if err := ring.Set(item); err != nil {
				logger.Debugf("Failed to write to Windows credential manager. Err: %v", err)
			}
		} else if runtime.GOOS == "darwin" {
			target = convertTarget(sc.cfg.Host, sc.cfg.User, credType)
			ring, _ := keyring.Open(keyring.Config{
				ServiceName: target,
			})
			account := strings.ToUpper(sc.cfg.User)
			item := keyring.Item{
				Key:  account,
				Data: []byte(token),
			}
			if err := ring.Set(item); err != nil {
				logger.Debugf("Failed to write to keychain. Err: %v", err)
			}
		} else if runtime.GOOS == "linux" {
			createCredentialCacheDir()
			writeTemporaryCredential(sc, credType, token)
		} else {
			logger.Debug("OS not supported for Local Secure Storage")
		}
	}
}

func getCredential(sc *snowflakeConn, credType string) {
	var target string
	cred := ""
	if runtime.GOOS == "windows" {
		target = driverName + ":" + credType
		ring, _ := keyring.Open(keyring.Config{
			WinCredPrefix: strings.ToUpper(sc.cfg.Host),
			ServiceName:   strings.ToUpper(sc.cfg.User),
		})
		i, err := ring.Get(target)
		if err != nil {
			logger.Debugf("Failed to read target or could not find it in Windows Credential Manager. Error: %v", err)
		}
		cred = string(i.Data)
	} else if runtime.GOOS == "darwin" {
		target = convertTarget(sc.cfg.Host, sc.cfg.User, credType)
		ring, _ := keyring.Open(keyring.Config{
			ServiceName: target,
		})
		account := strings.ToUpper(sc.cfg.User)
		i, err := ring.Get(account)
		if err != nil {
			logger.Debugf("Failed to find the item in keychain or item does not exist. Error: %v", err)
		}
		cred = string(i.Data)
		if cred == "" {
			logger.Debug("Returned credential is empty")
		} else {
			logger.Debug("Successfully read token. Returning as string")
		}
	} else if runtime.GOOS == "linux" {
		createCredentialCacheDir()
		cred = readTemporaryCredential(sc, credType)
	} else {
		logger.Debug("OS not supported for Local Secure Storage")
	}

	if credType == idToken {
		sc.cfg.IDToken = cred
	} else if credType == mfaToken {
		sc.cfg.MfaToken = cred
	} else {
		logger.Debugf("Unrecognized type %v for local cached credential", credType)
	}
}

func deleteCredential(sc *snowflakeConn, credType string) {
	target := driverName + ":" + credType
	if runtime.GOOS == "windows" {
		ring, _ := keyring.Open(keyring.Config{
			WinCredPrefix: strings.ToUpper(sc.cfg.Host),
			ServiceName:   strings.ToUpper(sc.cfg.User),
		})
		err := ring.Remove(target)
		if err != nil {
			logger.Debugf("Failed to delete target in Windows Credential Manager. Error: %v", err)
		}
	} else if runtime.GOOS == "darwin" {
		target = convertTarget(sc.cfg.Host, sc.cfg.User, credType)
		ring, _ := keyring.Open(keyring.Config{
			ServiceName: target,
		})
		account := strings.ToUpper(sc.cfg.User)
		err := ring.Remove(account)
		if err != nil {
			logger.Debugf("Failed to delete target in keychain. Error: %v", err)
		}
	} else if runtime.GOOS == "linux" {
		deleteTemporaryCredential(sc, credType)
	}
}

// Reads temporary credential file when OS is Linux.
func readTemporaryCredential(sc *snowflakeConn, credType string) string {
	target := convertTarget(sc.cfg.Host, sc.cfg.User, credType)
	credCacheLock.Lock()
	defer credCacheLock.Unlock()
	localCredCache := readTemporaryCacheFile()
	cred := localCredCache[target]
	if cred != "" {
		logger.Debug("Successfully read token. Returning as string")
	} else {
		logger.Debug("Returned credential is empty")
	}
	return cred
}

// Writes to temporary credential file when OS is Linux.
func writeTemporaryCredential(sc *snowflakeConn, credType, token string) {
	target := convertTarget(sc.cfg.Host, sc.cfg.User, credType)
	credCacheLock.Lock()
	defer credCacheLock.Unlock()
	localCredCache[target] = token

	j, err := json.Marshal(localCredCache)
	if err != nil {
		logger.Warnf("failed to convert credential to JSON.")
		return
	}
	writeTemporaryCacheFile(j)
}

func deleteTemporaryCredential(sc *snowflakeConn, credType string) {
	if credCacheDir == "" {
		logger.Debug("Cache file doesn't exist. Skipping deleting credential file.")
	} else {
		credCacheLock.Lock()
		defer credCacheLock.Unlock()
		target := convertTarget(sc.cfg.Host, sc.cfg.User, credType)
		delete(localCredCache, target)
		j, err := json.Marshal(localCredCache)
		if err != nil {
			logger.Warnf("failed to convert credential to JSON.")
			return
		}
		writeTemporaryCacheFile(j)
	}
}

func readTemporaryCacheFile() map[string]string {
	if credCache == "" {
		logger.Debug("Cache file doesn't exist. Skipping reading credential file.")
		return nil
	}
	jsonData, err := os.ReadFile(credCache)
	if err != nil {
		logger.Debugf("Failed to read credential file: %v", err)
		return nil
	}
	err = json.Unmarshal([]byte(jsonData), &localCredCache)
	if err != nil {
		logger.Debugf("failed to read JSON. Err: %v", err)
		return nil
	}

	return localCredCache
}

func writeTemporaryCacheFile(input []byte) {
	if credCache == "" {
		logger.Debug("Cache file doesn't exist. Skipping writing temporary credential file.")
	} else {
		logger.Debugf("writing credential cache file. %v\n", credCache)
		credCacheLockFileName := credCache + ".lck"
		err := os.Mkdir(credCacheLockFileName, 0600)
		logger.Debugf("Creating lock file. %v", credCacheLockFileName)

		switch {
		case os.IsExist(err):
			statinfo, err := os.Stat(credCacheLockFileName)
			if err != nil {
				logger.Debugf("failed to write credential cache file. file: %v, err: %v. ignored.\n", credCache, err)
				return
			}
			if time.Since(statinfo.ModTime()) < 15*time.Minute {
				logger.Debugf("other process locks the cache file. %v. ignored.\n", credCache)
				return
			}
			if err = os.Remove(credCacheLockFileName); err != nil {
				logger.Debugf("failed to delete lock file. file: %v, err: %v. ignored.\n", credCacheLockFileName, err)
				return
			}
			if err = os.Mkdir(credCacheLockFileName, 0600); err != nil {
				logger.Debugf("failed to delete lock file. file: %v, err: %v. ignored.\n", credCacheLockFileName, err)
				return
			}
		}
		defer os.RemoveAll(credCacheLockFileName)

		if err = os.WriteFile(credCache, input, 0644); err != nil {
			logger.Debugf("Failed to write the cache file. File: %v err: %v.", credCache, err)
		}
	}
}

func convertTarget(host, user, credType string) string {
	host = strings.ToUpper(host)
	user = strings.ToUpper(user)
	credType = strings.ToUpper(credType)
	target := host + ":" + user + ":" + driverName + ":" + credType
	return target
}
