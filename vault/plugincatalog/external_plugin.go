// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"encoding/hex"
	"encoding/json"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// Only plugins running with identical PluginRunner config can be multiplexed,
// so we use the PluginRunner input as the key for the external plugins map.
//
// However, to be a map key, it must be comparable:
// https://go.dev/ref/spec#Comparison_operators.
// In particular, the PluginRunner struct has slices and a function which are not
// comparable, so we need to transform it into a struct which is.
type externalPluginsKey struct {
	name     string
	typ      consts.PluginType
	version  string
	command  string
	ociImage string
	runtime  string
	args     string
	env      string
	sha256   string
	builtin  bool
}

func makeExternalPluginsKey(p *pluginutil.PluginRunner) (externalPluginsKey, error) {
	args, err := json.Marshal(p.Args)
	if err != nil {
		return externalPluginsKey{}, err
	}

	env, err := json.Marshal(p.Env)
	if err != nil {
		return externalPluginsKey{}, err
	}

	return externalPluginsKey{
		name:     p.Name,
		typ:      p.Type,
		version:  p.Version,
		command:  p.Command,
		ociImage: p.OCIImage,
		runtime:  p.Runtime,
		args:     string(args),
		env:      string(env),
		sha256:   hex.EncodeToString(p.Sha256),
		builtin:  p.Builtin,
	}, nil
}

// externalPlugin holds client connections for multiplexed and
// non-multiplexed plugin processes
type externalPlugin struct {
	// connections holds client connections by ID
	connections map[string]*pluginClient

	multiplexingSupport bool
}
