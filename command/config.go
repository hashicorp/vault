package command

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl"
	"github.com/mitchellh/go-homedir"
)

const (
	// DefaultConfigPath is the default path to the configuration file
	DefaultConfigPath = "~/.vault"

	// ConfigPathEnv is the environment variable that can be used to
	// override where the Vault configuration is.
	ConfigPathEnv = "VAULT_CONFIG_PATH"
)

// Config is the CLI configuration for Vault that can be specified via
// a `$HOME/.vault` file which is HCL-formatted (therefore HCL or JSON).
type Config struct {
	// TokenHelper is the executable/command that is executed for storing
	// and retrieving the authentication token for the Vault CLI. If this
	// is not specified, then vault's internal token store will be used, which
	// stores the token on disk unencrypted.
	TokenHelper string `hcl:"token_helper"`
}

// LoadConfig reads the configuration from the given path. If path is
// empty, then the default path will be used, or the environment variable
// if set.
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigPath
	}
	if v := os.Getenv(ConfigPathEnv); v != "" {
		path = v
	}

	path, err := homedir.Expand(path)
	if err != nil {
		return nil, fmt.Errorf("Error expanding config path: %s", err)
	}

	var config Config
	contents, err := ioutil.ReadFile(path)
	if !os.IsNotExist(err) {
		if err != nil {
			return nil, err
		}

		obj, err := hcl.Parse(string(contents))
		if err != nil {
			return nil, err
		}

		if err := hcl.DecodeObject(&config, obj); err != nil {
			return nil, err
		}
	}

	return &config, nil
}
