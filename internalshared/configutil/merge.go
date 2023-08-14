// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

func (c *SharedConfig) Merge(c2 *SharedConfig) *SharedConfig {
	if c2 == nil {
		return c
	}

	result := new(SharedConfig)

	for _, l := range c.Listeners {
		result.Listeners = append(result.Listeners, l)
	}
	for _, l := range c2.Listeners {
		result.Listeners = append(result.Listeners, l)
	}

	for _, userlockout := range c.UserLockouts {
		result.UserLockouts = append(result.UserLockouts, userlockout)
	}
	for _, userlockout := range c2.UserLockouts {
		result.UserLockouts = append(result.UserLockouts, userlockout)
	}

	result.HCPLinkConf = c.HCPLinkConf
	if c2.HCPLinkConf != nil {
		result.HCPLinkConf = c2.HCPLinkConf
	}

	result.Entropy = c.Entropy
	if c2.Entropy != nil {
		result.Entropy = c2.Entropy
	}

	for _, s := range c.Seals {
		result.Seals = append(result.Seals, s)
	}
	for _, s := range c2.Seals {
		result.Seals = append(result.Seals, s)
	}

	result.Telemetry = c.Telemetry
	if c2.Telemetry != nil {
		result.Telemetry = c2.Telemetry
	}

	result.DisableMlock = c.DisableMlock
	if c2.DisableMlock {
		result.DisableMlock = c2.DisableMlock
	}

	result.DefaultMaxRequestDuration = c.DefaultMaxRequestDuration
	if c2.DefaultMaxRequestDuration > result.DefaultMaxRequestDuration {
		result.DefaultMaxRequestDuration = c2.DefaultMaxRequestDuration
	}

	result.LogLevel = c.LogLevel
	if c2.LogLevel != "" {
		result.LogLevel = c2.LogLevel
	}

	result.LogFormat = c.LogFormat
	if c2.LogFormat != "" {
		result.LogFormat = c2.LogFormat
	}

	result.LogFile = c.LogFile
	if c2.LogFile != "" {
		result.LogFile = c2.LogFile
	}

	result.LogRotateBytes = c.LogRotateBytes
	if c2.LogRotateBytesRaw != nil {
		result.LogRotateBytes = c2.LogRotateBytes
		result.LogRotateBytesRaw = c2.LogRotateBytesRaw
	}

	result.LogRotateMaxFiles = c.LogRotateMaxFiles
	if c2.LogRotateMaxFilesRaw != nil {
		result.LogRotateMaxFiles = c2.LogRotateMaxFiles
		result.LogRotateMaxFilesRaw = c2.LogRotateMaxFilesRaw
	}

	result.LogRotateDuration = c.LogRotateDuration
	if c2.LogRotateDuration != "" {
		result.LogRotateDuration = c2.LogRotateDuration
	}

	result.PidFile = c.PidFile
	if c2.PidFile != "" {
		result.PidFile = c2.PidFile
	}

	result.ClusterName = c.ClusterName
	if c2.ClusterName != "" {
		result.ClusterName = c2.ClusterName
	}

	return result
}
