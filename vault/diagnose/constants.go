// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package diagnose

const (
	AutoLoadedLicenseValidatorError    = "Autoloaded license could not be validated: "
	AutoloadedLicenseValidationError   = "Autoloaded license validation failed due to error: "
	LicenseAutoloadingError            = "License could not be autoloaded: "
	StoredLicenseNoAutoloadingWarning  = "Vault is using a stored license, which is deprecated! Vault should use autoloaded licenses instead."
	NoStoredOrAutoloadedLicenseWarning = "No autoloaded or stored license could be detected."
	LicenseExpiredError                = "Autoloaded license is expired."
	LicenseExpiryThresholdWarning      = "Autoloaded license will expire "
	LicenseTerminatedError             = "Autoloaded license is terminated."
	LicenseTerminationThresholdWarning = "Autoloaded license will be terminated "
)
