package command

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/command/config"
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
//
// Config just calls into config.Config for backwards compatibility purposes.
// Use config.Config instead.
func Config() (*DefaultConfig, error) {
	conf, err := config.Config()
	return (*DefaultConfig)(conf), err
}

// LoadConfig reads the configuration from the given path. If path is
// empty, then the default path will be used, or the environment variable
// if set.
//
// LoadConfig just calls into config.LoadConfig for backwards compatibility
// purposes. Use config.LoadConfig instead.
func LoadConfig(path string) (*DefaultConfig, error) {
	conf, err := config.LoadConfig(path)
	return (*DefaultConfig)(conf), err
}

// ParseConfig parses the given configuration as a string.
//
// ParseConfig just calls into config.ParseConfig for backwards compatibility
// purposes. Use config.ParseConfig instead.
func ParseConfig(contents string) (*DefaultConfig, error) {
	conf, err := config.ParseConfig(contents)
	return (*DefaultConfig)(conf), err
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
