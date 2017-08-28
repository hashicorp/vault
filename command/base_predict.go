package command

import (
	"sort"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

// defaultPredictVaultMounts is the default list of mounts to return to the
// user. This is a best-guess, given we haven't communicated with the Vault
// server. If the user has no token or if the token does not have the default
// policy attached, it won't be able to read cubbyhole/, but it's a better UX
// that returning nothing.
var defaultPredictVaultMounts = []string{"cubbyhole/"}

// PredictVaultFiles returns a predictor for Vault mounts and paths based on the
// configured client for the base command. Unfortunately this happens pre-flag
// parsing, so users must rely on environment variables for autocomplete if they
// are not using Vault at the default endpoints.
func (b *BaseCommand) PredictVaultFiles() complete.Predictor {
	client, err := b.Client()
	if err != nil {
		return nil
	}
	return PredictVaultFiles(client)
}

// PredictVaultFolders returns a predictor for "folders". See PredictVaultFiles
// for more information and restrictions.
func (b *BaseCommand) PredictVaultFolders() complete.Predictor {
	client, err := b.Client()
	if err != nil {
		return nil
	}
	return PredictVaultFolders(client)
}

// PredictVaultFiles returns a predictor for Vault "files". This is a public API
// for consumers, but you probably want BaseCommand.PredictVaultFiles instead.
func PredictVaultFiles(client *api.Client) complete.Predictor {
	return predictVaultPaths(client, true)
}

// PredictVaultFolders returns a predictor for Vault "folders". This is a public
// API for consumers, but you probably want BaseCommand.PredictVaultFolders
// instead.
func PredictVaultFolders(client *api.Client) complete.Predictor {
	return predictVaultPaths(client, false)
}

// predictVaultPaths parses the CLI options and returns the "best" list of
// possible paths. If there are any errors, this function returns an empty
// result. All errors are suppressed since this is a prediction function.
func predictVaultPaths(client *api.Client, includeFiles bool) complete.PredictFunc {
	return func(args complete.Args) []string {
		// Do not predict more than one paths
		if predictHasPathArg(args.All) {
			return nil
		}

		path := args.Last

		var predictions []string
		if strings.Contains(path, "/") {
			predictions = predictPaths(client, path, includeFiles)
		} else {
			predictions = predictMounts(client, path)
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
		return predictVaultPaths(client, includeFiles).Predict(args)
	}
}

// predictMounts predicts all mounts which start with the given prefix. These
// are predicted on mount path, not "type".
func predictMounts(client *api.Client, path string) []string {
	mounts := predictListMounts(client)

	var predictions []string
	for _, m := range mounts {
		if strings.HasPrefix(m, path) {
			predictions = append(predictions, m)
		}
	}

	return predictions
}

// predictPaths predicts all paths which start with the given path.
func predictPaths(client *api.Client, path string, includeFiles bool) []string {
	// Vault does not support listing based on a sub-key, so we have to back-pedal
	// to the last "/" and return all paths on that "folder". Then we perform
	// client-side filtering.
	root := path
	idx := strings.LastIndex(root, "/")
	if idx > 0 && idx < len(root) {
		root = root[:idx+1]
	}

	paths := predictListPaths(client, root)

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

// predictListMounts returns a sorted list of the mount paths for Vault server
// for which the client is configured to communicate with. This function returns
// the default list of mounts if an error occurs.
func predictListMounts(c *api.Client) []string {
	mounts, err := c.Sys().ListMounts()
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

// predictListPaths returns a list of paths (HTTP LIST) for the given path. This
// function returns an empty list of any errors occur.
func predictListPaths(c *api.Client, path string) []string {
	secret, err := c.Logical().List(path)
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

// predictHasPathArg determines if the args have already accepted a path.
func predictHasPathArg(args []string) bool {
	var nonFlags []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			nonFlags = append(nonFlags, a)
		}
	}

	return len(nonFlags) > 2
}
