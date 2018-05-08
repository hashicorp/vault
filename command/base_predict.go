package command

import (
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

type Predict struct {
	client     *api.Client
	clientOnce sync.Once

	kv2s     []string
	kv2sOnce sync.Once
}

func NewPredict() *Predict {
	return &Predict{}
}

func (p *Predict) Client() *api.Client {
	p.clientOnce.Do(func() {
		if p.client == nil { // For tests
			client, _ := api.NewClient(nil)

			if client.Token() == "" {
				helper, err := DefaultTokenHelper()
				if err != nil {
					return
				}
				token, err := helper.Get()
				if err != nil {
					return
				}
				client.SetToken(token)
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
var predictClient *api.Client
var predictClientOnce sync.Once

// PredictClient returns the cached API client for the predictor.
func PredictClient() *api.Client {
	predictClientOnce.Do(func() {
		if predictClient == nil { // For tests
			predictClient, _ = api.NewClient(nil)
		}
	})
	return predictClient
}

// mountFilter is a function that filters mounts.
type mountFilter func(map[string]*api.MountOutput)

// mountFilterOnlyKV is a mount filter that returns only mount points of type
// "kv".
func mountFilterOnlyKV() mountFilter {
	return func(m map[string]*api.MountOutput) {
		for k, v := range m {
			if v.Type != "kv" {
				delete(m, k)
			}
		}
	}
}

// mountFilterOnlyKV is a mount filter that returns all mount points except type
// "kv".
func mountFilterExceptKV() mountFilter {
	return func(m map[string]*api.MountOutput) {
		for k, v := range m {
			if v.Type == "kv" {
				delete(m, k)
			}
		}
	}
}

// mountFilterOnlyKV2 is a mount filter that returns only KV v2 mounts.
func mountFilterOnlyKV2() mountFilter {
	return func(m map[string]*api.MountOutput) {
		for k, v := range m {
			if !mountFilterIsKV2(v) {
				delete(m, k)
			}
		}
	}
}

// mountFilterExceptKV2 is a mount filter that returns all mounts except those
// that are KV v2.
func mountFilterExceptKV2() mountFilter {
	return func(m map[string]*api.MountOutput) {
		for k, v := range m {
			if mountFilterIsKV2(v) {
				delete(m, k)
			}
		}
	}
}

// mountFilterIsKV2 is a helper to check if a given mount is type kv v2.
func mountFilterIsKV2(o *api.MountOutput) bool {
	if o.Type != "kv" {
		return false
	}

	if v, ok := o.Options["version"]; ok && v == "2" {
		return true
	}

	return false
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
func (b *BaseCommand) PredictVaultFiles(fs ...mountFilter) complete.Predictor {
	return NewPredict().VaultFiles(fs...)
}

// PredictVaultFolders returns a predictor for "folders". See PredictVaultFiles
// for more information and restrictions.
func (b *BaseCommand) PredictVaultFolders(fs ...mountFilter) complete.Predictor {
	return NewPredict().VaultFolders(fs...)
}

// PredictVaultMounts returns a predictor for "folders". See PredictVaultFiles
// for more information and restrictions.
func (b *BaseCommand) PredictVaultMounts(fs ...mountFilter) complete.Predictor {
	return NewPredict().VaultMounts(fs...)
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
func (b *BaseCommand) PredictVaultPlugins() complete.Predictor {
	return NewPredict().VaultPlugins()
}

// PredictVaultPolicies returns a predictor for "folders". See PredictVaultFiles
// for more information and restrictions.
func (b *BaseCommand) PredictVaultPolicies() complete.Predictor {
	return NewPredict().VaultPolicies()
}

// VaultFiles returns a predictor for Vault "files". This is a public API for
// consumers, but you probably want BaseCommand.PredictVaultFiles instead.
func (p *Predict) VaultFiles(fs ...mountFilter) complete.Predictor {
	return p.vaultPaths(p.mounts(fs...), true)
}

// VaultFolders returns a predictor for Vault "folders". This is a public
// API for consumers, but you probably want BaseCommand.PredictVaultFolders
// instead.
func (p *Predict) VaultFolders(fs ...mountFilter) complete.Predictor {
	return p.vaultPaths(p.mounts(fs...), false)
}

// VaultMounts returns a predictor for Vault "folders". This is a public
// API for consumers, but you probably want BaseCommand.PredictVaultMounts
// instead.
func (p *Predict) VaultMounts(fs ...mountFilter) complete.Predictor {
	return p.filterFunc(p.mountsFn(fs...))
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
func (p *Predict) VaultPlugins() complete.Predictor {
	return p.filterFunc(p.plugins)
}

// VaultPolicies returns a predictor for Vault "folders". This is a public API for
// consumers, but you probably want BaseCommand.PredictVaultPolicies instead.
func (p *Predict) VaultPolicies() complete.Predictor {
	return p.filterFunc(p.policies)
}

// vaultPaths parses the CLI options and returns the "best" list of possible
// paths. If there are any errors, this function returns an empty result. All
// errors are suppressed since this is a prediction function.
func (p *Predict) vaultPaths(mounts []string, includeFiles bool) complete.PredictFunc {
	return func(args complete.Args) []string {
		// Do not predict more than one paths
		if p.hasPathArg(args.All) {
			return nil
		}

		path := args.Last

		var predictions []string
		if strings.Contains(path, "/") {
			predictions = p.paths(path, includeFiles)
		} else {
			predictions = p.filter(mounts, path)
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
		return p.vaultPaths(mounts, includeFiles).Predict(args)
	}
}

// paths predicts all paths which start with the given path.
func (p *Predict) paths(path string, includeFiles bool) []string {
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

	paths := p.listPaths(root)

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
func (p *Predict) plugins() []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	result, err := client.Sys().ListPlugins(nil)
	if err != nil {
		return nil
	}
	plugins := result.Names
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

// mounts returns a sorted list of the mount paths for Vault server for
// which the client is configured to communicate with. This function returns the
// default list of mounts if an error occurs.
//
// If any filters are given, they are executed in order.
func (p *Predict) mounts(fs ...mountFilter) []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return defaultPredictVaultMounts
	}

	for _, f := range fs {
		f(mounts)
	}

	list := make([]string, 0, len(mounts))
	for m := range mounts {
		list = append(list, m)
	}
	sort.Strings(list)
	return list
}

// mountFn is a wrapper around the mount func to return a slice of string
// with no args for filterFunc
func (p *Predict) mountsFn(fs ...mountFilter) func() []string {
	return func() []string {
		return p.mounts(fs...)
	}
}

// listPaths returns a list of paths (HTTP LIST) for the given path. This
// function returns an empty list of any errors occur.
func (p *Predict) listPaths(pth string) []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	// Handle listing for kv v2
	if mountPath, ok := p.kv2Path(pth); ok {
		pth = strings.TrimPrefix(pth, mountPath)
		pth = path.Join(mountPath, "metadata", pth)
	}

	secret, err := client.Logical().List(pth)
	if err != nil || secret == nil || secret.Data == nil {
		return nil
	}

	pths, ok := secret.Data["keys"].([]interface{})
	if !ok {
		return nil
	}

	list := make([]string, 0, len(pths))
	for _, p := range pths {
		if str, ok := p.(string); ok {
			list = append(list, str)
		}
	}
	sort.Strings(list)
	return list
}

// kv2Path returns the mount point and "true" if the given path is a KV v2
// mount. Otherwise it returns the empty string and false.
func (p *Predict) kv2Path(pth string) (string, bool) {
	for _, v := range p.kv2Paths() {
		if !strings.HasSuffix(v, "/") {
			v = v + "/"
		}

		if strings.HasPrefix(pth, v) {
			return v, true
		}
	}

	return "", false
}

// kv2Paths is the list of mounts that are of type "kv" with version "2". This
// is cached for performance.
func (p *Predict) kv2Paths() []string {
	p.kv2sOnce.Do(func() {
		p.kv2s = p.mounts(mountFilterOnlyKV2())
		sort.Strings(p.kv2s)
	})
	return p.kv2s
}

// hasPathArg determines if the args have already accepted a path.
func (p *Predict) hasPathArg(args []string) bool {
	var nonFlags []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			nonFlags = append(nonFlags, a)
		}
	}

	includesKV := false
	for _, v := range nonFlags {
		if v == "kv" {
			includesKV = true
		}
	}

	if includesKV {
		return len(nonFlags) > 3
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
