// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/cliconfig"
	"github.com/posener/complete"
)

type Predict struct {
	client     *api.Client
	clientOnce sync.Once
}

func NewPredict() *Predict {
	return &Predict{}
}

func (p *Predict) Client() *api.Client {
	p.clientOnce.Do(func() {
		if p.client == nil { // For tests
			client, _ := api.NewClient(nil)

			if client.Token() == "" {
				helper, err := cliconfig.DefaultTokenHelper()
				if err != nil {
					return
				}
				token, err := helper.Get()
				if err != nil {
					return
				}
				client.SetToken(token)
			}

			// Turn off retries for prediction
			if os.Getenv(api.EnvVaultMaxRetries) == "" {
				client.SetMaxRetries(0)
			}

			p.client = client
		}
	})
	return p.client
}

// defaultPredictVaultMounts is the default list of mounts to return to the
// user. This is a best-guess, given we haven't communicated with the Vault
// server. If the user has no token or if the token does not have the default
// policy attached, it won't be able to read cubbyhole/, but it's a better UX
// that returning nothing.
var defaultPredictVaultMounts = []string{"cubbyhole/"}

// predictClient is the API client to use for prediction. We create this at the
// beginning once, because completions are generated for each command (and this
// doesn't change), and the only way to configure the predict/autocomplete
// client is via environment variables. Even if the user specifies a flag, we
// can't parse that flag until after the command is submitted.
var (
	predictClient     *api.Client
	predictClientOnce sync.Once
)

// PredictClient returns the cached API client for the predictor.
func PredictClient() *api.Client {
	predictClientOnce.Do(func() {
		if predictClient == nil { // For tests
			predictClient, _ = api.NewClient(nil)
		}
	})
	return predictClient
}

// PredictVaultAvailableMounts returns a predictor for the available mounts in
// Vault. For now, there is no way to programmatically get this list. If, in the
// future, such a list exists, we can adapt it here. Until then, it's
// hard-coded.
func (b *BaseCommand) PredictVaultAvailableMounts() complete.Predictor {
	// This list does not contain deprecated backends. At present, there is no
	// API that lists all available secret backends, so this is hard-coded :(.
	return complete.PredictSet(
		"aws",
		"consul",
		"database",
		"generic",
		"pki",
		"plugin",
		"rabbitmq",
		"ssh",
		"totp",
		"transit",
	)
}

// PredictVaultAvailableAuths returns a predictor for the available auths in
// Vault. For now, there is no way to programmatically get this list. If, in the
// future, such a list exists, we can adapt it here. Until then, it's
// hard-coded.
func (b *BaseCommand) PredictVaultAvailableAuths() complete.Predictor {
	return complete.PredictSet(
		"app-id",
		"approle",
		"aws",
		"cert",
		"gcp",
		"github",
		"ldap",
		"okta",
		"plugin",
		"radius",
		"userpass",
	)
}

// PredictVaultFiles returns a predictor for Vault mounts and paths based on the
// configured client for the base command. Unfortunately this happens pre-flag
// parsing, so users must rely on environment variables for autocomplete if they
// are not using Vault at the default endpoints.
func (b *BaseCommand) PredictVaultFiles() complete.Predictor {
	return NewPredict().VaultFiles()
}

// PredictVaultFolders returns a predictor for "folders". See PredictVaultFiles
// for more information and restrictions.
func (b *BaseCommand) PredictVaultFolders() complete.Predictor {
	return NewPredict().VaultFolders()
}

// PredictVaultNamespaces returns a predictor for "namespaces". See PredictVaultFiles
// for more information an restrictions.
func (b *BaseCommand) PredictVaultNamespaces() complete.Predictor {
	return NewPredict().VaultNamespaces()
}

// PredictVaultMounts returns a predictor for "folders". See PredictVaultFiles
// for more information and restrictions.
func (b *BaseCommand) PredictVaultMounts() complete.Predictor {
	return NewPredict().VaultMounts()
}

// PredictVaultAudits returns a predictor for "folders". See PredictVaultFiles
// for more information and restrictions.
func (b *BaseCommand) PredictVaultAudits() complete.Predictor {
	return NewPredict().VaultAudits()
}

// PredictVaultAuths returns a predictor for "folders". See PredictVaultFiles
// for more information and restrictions.
func (b *BaseCommand) PredictVaultAuths() complete.Predictor {
	return NewPredict().VaultAuths()
}

// PredictVaultPlugins returns a predictor for installed plugins.
func (b *BaseCommand) PredictVaultPlugins(pluginTypes ...api.PluginType) complete.Predictor {
	return NewPredict().VaultPlugins(pluginTypes...)
}

// PredictVaultPolicies returns a predictor for "folders". See PredictVaultFiles
// for more information and restrictions.
func (b *BaseCommand) PredictVaultPolicies() complete.Predictor {
	return NewPredict().VaultPolicies()
}

func (b *BaseCommand) PredictVaultDebugTargets() complete.Predictor {
	return complete.PredictSet(
		"config",
		"host",
		"metrics",
		"pprof",
		"replication-status",
		"server-status",
	)
}

// VaultFiles returns a predictor for Vault "files". This is a public API for
// consumers, but you probably want BaseCommand.PredictVaultFiles instead.
func (p *Predict) VaultFiles() complete.Predictor {
	return p.vaultPaths(true)
}

// VaultFolders returns a predictor for Vault "folders". This is a public
// API for consumers, but you probably want BaseCommand.PredictVaultFolders
// instead.
func (p *Predict) VaultFolders() complete.Predictor {
	return p.vaultPaths(false)
}

// VaultNamespaces returns a predictor for Vault "namespaces". This is a public
// API for consumers, but you probably want BaseCommand.PredictVaultNamespaces
// instead.
func (p *Predict) VaultNamespaces() complete.Predictor {
	return p.filterFunc(p.namespaces)
}

// VaultMounts returns a predictor for Vault "folders". This is a public
// API for consumers, but you probably want BaseCommand.PredictVaultMounts
// instead.
func (p *Predict) VaultMounts() complete.Predictor {
	return p.filterFunc(p.mounts)
}

// VaultAudits returns a predictor for Vault "folders". This is a public API for
// consumers, but you probably want BaseCommand.PredictVaultAudits instead.
func (p *Predict) VaultAudits() complete.Predictor {
	return p.filterFunc(p.audits)
}

// VaultAuths returns a predictor for Vault "folders". This is a public API for
// consumers, but you probably want BaseCommand.PredictVaultAuths instead.
func (p *Predict) VaultAuths() complete.Predictor {
	return p.filterFunc(p.auths)
}

// VaultPlugins returns a predictor for Vault's plugin catalog. This is a public
// API for consumers, but you probably want BaseCommand.PredictVaultPlugins
// instead.
func (p *Predict) VaultPlugins(pluginTypes ...api.PluginType) complete.Predictor {
	filterFunc := func() []string {
		return p.plugins(pluginTypes...)
	}
	return p.filterFunc(filterFunc)
}

// VaultPolicies returns a predictor for Vault "folders". This is a public API for
// consumers, but you probably want BaseCommand.PredictVaultPolicies instead.
func (p *Predict) VaultPolicies() complete.Predictor {
	return p.filterFunc(p.policies)
}

// vaultPaths parses the CLI options and returns the "best" list of possible
// paths. If there are any errors, this function returns an empty result. All
// errors are suppressed since this is a prediction function.
func (p *Predict) vaultPaths(includeFiles bool) complete.PredictFunc {
	return func(args complete.Args) []string {
		// Do not predict more than one paths
		if p.hasPathArg(args.All) {
			return nil
		}

		client := p.Client()
		if client == nil {
			return nil
		}

		path := args.Last

		// Trim path with potential mount
		var relativePath string
		mountInfos, err := p.mountInfos()
		if err != nil {
			return nil
		}

		var mountType, mountVersion string
		for mount, mountInfo := range mountInfos {
			if strings.HasPrefix(path, mount) {
				relativePath = strings.TrimPrefix(path, mount+"/")
				mountType = mountInfo.Type
				if mountInfo.Options != nil {
					mountVersion = mountInfo.Options["version"]
				}
				break
			}
		}

		// Predict path or mount depending on path separator
		var predictions []string
		if strings.Contains(relativePath, "/") {
			predictions = p.paths(mountType, mountVersion, path, includeFiles)
		} else {
			predictions = p.filter(p.mounts(), path)
		}

		// Either no results or many results, so return.
		if len(predictions) != 1 {
			return predictions
		}

		// If this is not a "folder", do not try to recurse.
		if !strings.HasSuffix(predictions[0], "/") {
			return predictions
		}

		// If the prediction is the same as the last guess, return it (we have no
		// new information and we won't get anymore).
		if predictions[0] == args.Last {
			return predictions
		}

		// Re-predict with the remaining path
		args.Last = predictions[0]
		return p.vaultPaths(includeFiles).Predict(args)
	}
}

// paths predicts all paths which start with the given path.
func (p *Predict) paths(mountType, mountVersion, path string, includeFiles bool) []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	// Vault does not support listing based on a sub-key, so we have to back-pedal
	// to the last "/" and return all paths on that "folder". Then we perform
	// client-side filtering.
	root := path
	idx := strings.LastIndex(root, "/")
	if idx > 0 && idx < len(root) {
		root = root[:idx+1]
	}

	paths := p.listPaths(buildAPIListPath(root, mountType, mountVersion))

	var predictions []string
	for _, p := range paths {
		// Calculate the absolute "path" for matching.
		p = root + p

		if strings.HasPrefix(p, path) {
			// Ensure this is a directory or we've asked to include files.
			if includeFiles || strings.HasSuffix(p, "/") {
				predictions = append(predictions, p)
			}
		}
	}

	// Add root to the path
	if len(predictions) == 0 {
		predictions = append(predictions, path)
	}

	return predictions
}

func buildAPIListPath(path, mountType, mountVersion string) string {
	if mountType == "kv" && mountVersion == "2" {
		return toKVv2ListPath(path)
	}
	return path
}

func toKVv2ListPath(path string) string {
	firstSlashIdx := strings.Index(path, "/")
	if firstSlashIdx < 0 {
		return path
	}

	return path[:firstSlashIdx] + "/metadata" + path[firstSlashIdx:]
}

// audits returns a sorted list of the audit backends for Vault server for
// which the client is configured to communicate with.
func (p *Predict) audits() []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	audits, err := client.Sys().ListAudit()
	if err != nil {
		return nil
	}

	list := make([]string, 0, len(audits))
	for m := range audits {
		list = append(list, m)
	}
	sort.Strings(list)
	return list
}

// auths returns a sorted list of the enabled auth provides for Vault server for
// which the client is configured to communicate with.
func (p *Predict) auths() []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	auths, err := client.Sys().ListAuth()
	if err != nil {
		return nil
	}

	list := make([]string, 0, len(auths))
	for m := range auths {
		list = append(list, m)
	}
	sort.Strings(list)
	return list
}

// plugins returns a sorted list of the plugins in the catalog.
func (p *Predict) plugins(pluginTypes ...api.PluginType) []string {
	// This method's signature doesn't enforce that a pluginType must be passed in.
	// If it's not, it's likely the caller's intent is go get a list of all of them,
	// so let's help them out.
	if len(pluginTypes) == 0 {
		pluginTypes = append(pluginTypes, api.PluginTypeUnknown)
	}

	client := p.Client()
	if client == nil {
		return nil
	}

	var plugins []string
	pluginsAdded := make(map[string]bool)
	for _, pluginType := range pluginTypes {
		result, err := client.Sys().ListPlugins(&api.ListPluginsInput{Type: api.PluginType(pluginType)})
		if err != nil {
			return nil
		}
		if result == nil {
			return nil
		}
		for _, names := range result.PluginsByType {
			for _, name := range names {
				if _, ok := pluginsAdded[name]; !ok {
					plugins = append(plugins, name)
					pluginsAdded[name] = true
				}
			}
		}
	}
	sort.Strings(plugins)
	return plugins
}

// policies returns a sorted list of the policies stored in this Vault
// server.
func (p *Predict) policies() []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	policies, err := client.Sys().ListPolicies()
	if err != nil {
		return nil
	}
	sort.Strings(policies)
	return policies
}

// mountInfos returns a map with mount paths as keys and MountOutputs as values
// for the Vault server which the client is configured to communicate with.
// Returns error if server communication fails.
func (p *Predict) mountInfos() (map[string]*api.MountOutput, error) {
	client := p.Client()
	if client == nil {
		return nil, nil
	}

	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	return mounts, nil
}

// mounts returns a sorted list of the mount paths for Vault server for
// which the client is configured to communicate with. This function returns the
// default list of mounts if an error occurs.
func (p *Predict) mounts() []string {
	mounts, err := p.mountInfos()
	if err != nil {
		return defaultPredictVaultMounts
	}

	list := make([]string, 0, len(mounts))
	for m := range mounts {
		list = append(list, m)
	}
	sort.Strings(list)
	return list
}

// namespaces returns a sorted list of the namespace paths for Vault server for
// which the client is configured to communicate with. This function returns
// an empty list in any error occurs.
func (p *Predict) namespaces() []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	secret, err := client.Logical().List("sys/namespaces")
	if err != nil {
		return nil
	}
	namespaces, ok := extractListData(secret)
	if !ok {
		return nil
	}

	list := make([]string, 0, len(namespaces))
	for _, n := range namespaces {
		s, ok := n.(string)
		if !ok {
			continue
		}
		list = append(list, s)
	}
	sort.Strings(list)
	return list
}

// listPaths returns a list of paths (HTTP LIST) for the given path. This
// function returns an empty list of any errors occur.
func (p *Predict) listPaths(path string) []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	secret, err := client.Logical().List(path)
	if err != nil || secret == nil || secret.Data == nil {
		return nil
	}

	paths, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil
	}

	list := make([]string, 0, len(paths))
	for _, p := range paths {
		if str, ok := p.(string); ok {
			list = append(list, str)
		}
	}
	sort.Strings(list)
	return list
}

// hasPathArg determines if the args have already accepted a path.
func (p *Predict) hasPathArg(args []string) bool {
	var nonFlags []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			nonFlags = append(nonFlags, a)
		}
	}

	return len(nonFlags) > 2
}

// filterFunc is used to compose a complete predictor that filters an array
// of strings as per the filter function.
func (p *Predict) filterFunc(f func() []string) complete.Predictor {
	return complete.PredictFunc(func(args complete.Args) []string {
		if p.hasPathArg(args.All) {
			return nil
		}

		client := p.Client()
		if client == nil {
			return nil
		}

		return p.filter(f(), args.Last)
	})
}

// filter filters the given list for items that start with the prefix.
func (p *Predict) filter(list []string, prefix string) []string {
	var predictions []string
	for _, item := range list {
		if strings.HasPrefix(item, prefix) {
			predictions = append(predictions, item)
		}
	}
	return predictions
}
