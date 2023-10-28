// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/sdk/helper/hclutil"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	// DefaultConfigPath is the default path to the configuration file
	DefaultConfigPath = "~/.vault"

	// ConfigPathEnv is the environment variable that can be used to
	// override where the Vault configuration is.
	ConfigPathEnv = "VAULT_CONFIG_PATH"

	// DefaultClientContextConfig is the default path to the client context configuration file
	DefaultClientContextConfig = "~/.vault-client-context"

	// ClientContextConfigPathEnv is the environment variable that can be used to
	// override where the client context configuration is.
	ClientContextConfigPathEnv = "VAULT_CLIENT_CONTEXT_CONFIG_PATH"
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

type ClientContextConfig struct {
	ClientContexts []ContextInfo `hcl:"client_context,block"`
	CurrentContext ContextInfo   `hcl:"current_context,block"`
}

type ContextInfo struct {
	Name          string `hcl:"name"`
	ClusterToken  string `hcl:"cluster_token"`
	VaultAddr     string `hcl:"cluster_addr"`
	NamespacePath string `hcl:"namespace_path"`
}

func NewContextConfig() ClientContextConfig {
	return ClientContextConfig{}
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
		return nil, fmt.Errorf("error expanding config path %q: %w", path, err)
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	conf, err := ParseConfig(string(contents))
	if err != nil {
		return nil, fmt.Errorf("error parsing config file at %q: %w; ensure that the file is valid; Ansible Vault is known to conflict with it.", path, err)
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

func WriteClientContextConfig(path string, config ClientContextConfig) (err error) {
	if path == "" {
		path = DefaultClientContextConfig
	}
	if v := os.Getenv(ClientContextConfigPathEnv); v != "" {
		path = v
	}

	contents := hclwrite.NewEmptyFile()

	gohcl.EncodeIntoBody(&config, contents.Body())

	// NOTE: requires HOME env var to be set
	path, err = homedir.Expand(path)
	if err != nil {
		return fmt.Errorf("error expanding client context config path %q: %w", path, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create configuration file %q: %v", path, err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			err = fmt.Errorf("failed to close configuration file %q: %v", path, err)
		}
	}()

	if _, err := contents.WriteTo(f); err != nil {
		return fmt.Errorf("failed to write to configuration file %q: %v", path, err)
	}

	return err
}

func LoadClientContextConfig(path string) (ClientContextConfig, error) {
	result := NewContextConfig()

	if path == "" {
		path = DefaultClientContextConfig
	}
	if v := os.Getenv(ClientContextConfigPathEnv); v != "" {
		path = v
	}

	// NOTE: requires HOME env var to be set
	path, err := homedir.Expand(path)
	if err != nil {
		return result, fmt.Errorf("error expanding client context config path %q: %w", path, err)
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return result, err
	}

	conf, err := ParseClientContextConfig(string(contents))
	if err != nil {
		return result, fmt.Errorf("error parsing client context config file at %q: %w", path, err)
	}

	return conf, nil
}

// ParseClientContextConfig parses the given configuration as a string.
func ParseClientContextConfig(contents string) (ClientContextConfig, error) {
	result := NewContextConfig()

	root, err := hcl.Parse(contents)
	if err != nil {
		return result, err
	}

	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return result, fmt.Errorf("failed to parse config; does not contain a root object")
	}

	if o := list.Filter("current_context"); len(o.Items) > 0 {
		if err := parseCurrentContext(&result.CurrentContext, o, "current_context"); err != nil {
			return result, fmt.Errorf("error parsing 'current_context', %w", err)
		}
	}

	if o := list.Filter("client_context"); len(o.Items) > 0 {
		if err := parseClientContexts(&result, o, "client_contexts"); err != nil {
			return result, fmt.Errorf("error parsing 'client_contexts', %w", err)
		}
	}

	return result, nil
}

func parseCurrentContext(result *ContextInfo, list *ast.ObjectList, name string) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one %q is allowed at a time", name)
	}
	item := list.Items[0]

	key := name
	if len(item.Keys) > 0 {
		key = item.Keys[0].Token.Value().(string)
	}

	if err := hcl.DecodeObject(result, item.Val); err != nil {
		return fmt.Errorf("failed to decode object for the current context, name.key: %s.%s, error:%w", name, key, err)
	}

	return nil
}

func parseClientContexts(result *ClientContextConfig, list *ast.ObjectList, name string) error {
	if result.ClientContexts == nil {
		result.ClientContexts = make([]ContextInfo, 0, len(list.Items))
	}

	for i, item := range list.Items {
		var ci ContextInfo
		if err := hcl.DecodeObject(&ci, item.Val); err != nil {
			return fmt.Errorf("failed to decode %q at location %d, error: %w", name, i, err)
		}
		switch {
		case ci.Name != "":
		case len(item.Keys) == 1:
			ci.Name = strings.ToLower(item.Keys[0].Token.Value().(string))
		default:
			return fmt.Errorf("failed to parse client context name for location %d", i)
		}

		result.ClientContexts = append(result.ClientContexts, ci)
	}

	return nil
}

func FindContextInfoIndexByName(infoSlice []ContextInfo, name string) (int, bool) {
	var index int
	var found bool
	for i, ctx := range infoSlice {
		if ctx.Name == name {
			found = true
			index = i
			break
		}
	}

	return index, found
}
