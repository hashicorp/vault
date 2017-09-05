package command

import (
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/vault/api"
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

		var predictions []string
		if strings.Contains(path, "/") {
			predictions = p.paths(path, includeFiles)
		} else {
			predictions = p.mounts(path)
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

// mounts predicts all mounts which start with the given prefix. These are
// predicted on mount path, not "type".
func (p *Predict) mounts(path string) []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	mounts := p.listMounts()

	var predictions []string
	for _, m := range mounts {
		if strings.HasPrefix(m, path) {
			predictions = append(predictions, m)
		}
	}

	return predictions
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

// listMounts returns a sorted list of the mount paths for Vault server for
// which the client is configured to communicate with. This function returns the
// default list of mounts if an error occurs.
func (p *Predict) listMounts() []string {
	client := p.Client()
	if client == nil {
		return nil
	}

	mounts, err := client.Sys().ListMounts()
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
