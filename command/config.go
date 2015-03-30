package command

// Config is the CLI configuration for Vault that can be specified via
// a `$HOME/.vault` file which is HCL-formatted (therefore HCL or JSON).
type Config struct {
	// TokenHelper is the executable/command that is executed for storing
	// and retrieving the authentication token for the Vault CLI. If this
	// is not specified, then vault token-disk will be used, which stores
	// the token on disk unencrypted.
	TokenHelper string `hcl:"token_helper"`
}
