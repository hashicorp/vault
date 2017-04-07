package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
)

var (
	pluginCatalogPrefix = "plugin-catalog/"
)

type PluginCatalog struct {
	catalogView  *BarrierView
	directory    string
	vaultCommand string
	vaultSHA256  []byte

	lock sync.RWMutex
}

func (c *Core) setupPluginCatalog() error {
	c.pluginCatalog = &PluginCatalog{
		catalogView:  c.systemBarrierView.SubView(pluginCatalogPrefix),
		directory:    c.pluginDirectory,
		vaultCommand: c.vaultBinaryLocation,
		vaultSHA256:  c.vaultBinarySHA256,
	}

	return nil
}

func (c *PluginCatalog) Get(name string) (*pluginutil.PluginRunner, error) {
	// Look for external plugins in the barrier
	out, err := c.catalogView.Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve plugin \"%s\": %v", name, err)
	}
	if out != nil {
		entry := new(pluginutil.PluginRunner)
		if err := jsonutil.DecodeJSON(out.Value, entry); err != nil {
			return nil, fmt.Errorf("failed to decode plugin entry: %v", err)
		}

		return entry, nil
	}

	// Look for builtin plugins
	if _, ok := builtinplugins.BuiltinPlugins.Get(name); !ok {
		return nil, fmt.Errorf("no plugin found with name: %s", name)
	}

	return &pluginutil.PluginRunner{
		Name:    name,
		Command: c.vaultCommand,
		Args:    []string{"plugin-exec", name},
		Sha256:  c.vaultSHA256,
	}, nil
}

func (c *PluginCatalog) Set(name, command string, sha256 []byte) error {
	parts := strings.Split(command, " ")
	command = parts[0]
	args := parts[1:]

	command = filepath.Join(c.directory, command)

	// Best effort check to make sure the command isn't breaking out of the
	// configured plugin directory.
	sym, err := filepath.EvalSymlinks(command)
	if err != nil {
		return fmt.Errorf("error while validating the command path: %v", err)
	}
	symAbs, err := filepath.Abs(filepath.Dir(sym))
	if err != nil {
		return fmt.Errorf("error while validating the command path: %v", err)
	}

	if symAbs != c.directory {
		return errors.New("can not execute files outside of configured plugin directory")
	}

	entry := &pluginutil.PluginRunner{
		Name:    name,
		Command: command,
		Args:    args,
		Sha256:  sha256,
	}

	buf, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to encode plugin entry: %v", err)
	}

	logicalEntry := logical.StorageEntry{
		Key:   name,
		Value: buf,
	}
	if err := c.catalogView.Put(&logicalEntry); err != nil {
		return fmt.Errorf("failed to persist plugin entry: %v", err)
	}
	return nil
}
