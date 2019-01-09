package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/helper/hclutil"
	homedir "github.com/mitchellh/go-homedir"
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
type DefaultConfig struct {
	// TokenHelper is the executable/command that is executed for storing
	// and retrieving the authentication token for the Vault CLI. If this
	// is not specified, then vault's internal token store will be used, which
	// stores the token on disk unencrypted.
	TokenHelper string `hcl:"token_helper"`
}

// Config loads the configuration and returns it. If the configuration
// is already loaded, it is returned.
func Config() (*DefaultConfig, error) {
	var err error
	config, err := LoadConfig("")
	if err != nil {
		return nil, err
	}

	return config, nil
}

// LoadConfig reads the configuration from the given path. If path is
// empty, then the default path will be used, or the environment variable
// if set.
func LoadConfig(path string) (*DefaultConfig, error) {
	if path == "" {
		path = DefaultConfigPath
	}
	if v := os.Getenv(ConfigPathEnv); v != "" {
		path = v
	}

	// NOTE: requires HOME env var to be set
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("error expanding config path %q: {{err}}", path), err)
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	conf, err := ParseConfig(string(contents))
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("error parsing config file at %q: {{err}}; ensure that the file is valid; Ansible Vault is known to conflict with it.", path), err)
	}

	return conf, nil
}

// ParseConfig parses the given configuration as a string.
func ParseConfig(contents string) (*DefaultConfig, error) {
	root, err := hcl.Parse(contents)
	if err != nil {
		return nil, err
	}

	// Top-level item should be the object list
	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("failed to parse config; does not contain a root object")
	}

	valid := []string{
		"token_helper",
	}
	if err := hclutil.CheckHCLKeys(list, valid); err != nil {
		return nil, err
	}

	var c DefaultConfig
	if err := hcl.DecodeObject(&c, list); err != nil {
		return nil, err
	}
	return &c, nil
}
