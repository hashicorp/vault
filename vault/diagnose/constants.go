package diagnose

const (
	AutoLoadedLicenseValidatorError    = "Autoloaded license could not be validated: "
	AutoloadedLicenseValidationError   = "Autoloaded license validation failed due to error: "
	LicenseAutoloadingError            = "license could not be autoloaded: "
	StoredLicenseNoAutoloadingWarning  = "Vault is using a stored license, which is deprecated! Vault should use autoloaded licenses instead."
	NoStoredOrAutoloadedLicenseWarning = "No autoloaded or stored license could be detected. If the binary is not a pro/prem binary, this means Vault does not have access to a license at all."
	LicenseExpiredError                = "Autoloaded license is expired."
	LicenseExpiryThresholdWarning      = "Autoloaded license will expire "
	LicenseTerminatedError             = "Autoloaded license is terminated."
	LicenseTerminationThresholdWarning = "Autoloaded license will be terminated "
)
