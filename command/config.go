package command

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
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
		return nil, fmt.Errorf("Error expanding config path %s: %s", path, err)
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return ParseConfig(string(contents))
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
		return nil, fmt.Errorf("Failed to parse config: does not contain a root object")
	}

	valid := []string{
		"token_helper",
	}
	if err := checkHCLKeys(list, valid); err != nil {
		return nil, err
	}

	var c DefaultConfig
	if err := hcl.DecodeObject(&c, list); err != nil {
		return nil, err
	}
	return &c, nil
}

func checkHCLKeys(node ast.Node, valid []string) error {
	var list *ast.ObjectList
	switch n := node.(type) {
	case *ast.ObjectList:
		list = n
	case *ast.ObjectType:
		list = n.List
	default:
		return fmt.Errorf("cannot check HCL keys of type %T", n)
	}

	validMap := make(map[string]struct{}, len(valid))
	for _, v := range valid {
		validMap[v] = struct{}{}
	}

	var result error
	for _, item := range list.Items {
		key := item.Keys[0].Token.Value().(string)
		if _, ok := validMap[key]; !ok {
			result = multierror.Append(result, fmt.Errorf(
				"invalid key '%s' on line %d", key, item.Assign.Line))
		}
	}

	return result
}
