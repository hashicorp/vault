// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"errors"
	"fmt"
	"io"
	paths "path"
	"sort"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func kvReadRequest(client *api.Client, path string, params map[string]string) (*api.Secret, error) {
	r := client.NewRequest("GET", "/v1/"+path)
	for k, v := range params {
		r.Params.Set(k, v)
	}
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := api.ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return api.ParseSecret(resp.Body)
}

func kvPreflightVersionRequest(client *api.Client, path string) (string, int, error) {
	// We don't want to use a wrapping call here so save any custom value and
	// restore after
	currentWrappingLookupFunc := client.CurrentWrappingLookupFunc()
	client.SetWrappingLookupFunc(nil)
	defer client.SetWrappingLookupFunc(currentWrappingLookupFunc)
	currentOutputCurlString := client.OutputCurlString()
	client.SetOutputCurlString(false)
	defer client.SetOutputCurlString(currentOutputCurlString)
	currentOutputPolicy := client.OutputPolicy()
	client.SetOutputPolicy(false)
	defer client.SetOutputPolicy(currentOutputPolicy)

	r := client.NewRequest("GET", "/v1/sys/internal/ui/mounts/"+path)
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// If we get a 404 we are using an older version of vault, default to
		// version 1
		if resp != nil {
			if resp.StatusCode == 404 {
				return "", 1, nil
			}

			// if the original request had the -output-curl-string or -output-policy flag,
			if (currentOutputCurlString || currentOutputPolicy) && resp.StatusCode == 403 {
				// we provide a more helpful error for the user,
				// who may not understand why the flag isn't working.
				err = fmt.Errorf(
					`This output flag requires the success of a preflight request 
to determine the version of a KV secrets engine. Please 
re-run this command with a token with read access to %s. 
Note that if the path you are trying to reach is a KV v2 path, your token's policy must 
allow read access to that path in the format 'mount-path/data/foo', not just 'mount-path/foo'.`, path)
			}
		}

		return "", 0, err
	}

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", 0, err
	}
	if secret == nil {
		return "", 0, errors.New("nil response from pre-flight request")
	}
	var mountPath string
	if mountPathRaw, ok := secret.Data["path"]; ok {
		mountPath = mountPathRaw.(string)
	}
	options := secret.Data["options"]
	if options == nil {
		return mountPath, 1, nil
	}
	versionRaw := options.(map[string]interface{})["version"]
	if versionRaw == nil {
		return mountPath, 1, nil
	}
	version := versionRaw.(string)
	switch version {
	case "", "1":
		return mountPath, 1, nil
	case "2":
		return mountPath, 2, nil
	}

	return mountPath, 1, nil
}

func IsKVv2(path string, client *api.Client) (string, bool, error) {
	mountPath, version, err := kvPreflightVersionRequest(client, path)
	if err != nil {
		return "", false, err
	}

	return mountPath, version == 2, nil
}

func AddPrefixToKVPath(path, mountPath, apiPrefix string, skipIfExists bool) string {
	if path == mountPath || path == strings.TrimSuffix(mountPath, "/") {
		return paths.Join(mountPath, apiPrefix)
	}

	pathSuffix := strings.TrimPrefix(path, mountPath)
	for {
		// If the entire mountPath is included in the path, we are done
		if pathSuffix != path {
			break
		}
		// Trim the parts of the mountPath that are not included in the
		// path, for example, in cases where the mountPath contains
		// namespaces which are not included in the path.
		partialMountPath := strings.SplitN(mountPath, "/", 2)
		if len(partialMountPath) <= 1 || partialMountPath[1] == "" {
			break
		}
		mountPath = strings.TrimSuffix(partialMountPath[1], "/")
		pathSuffix = strings.TrimPrefix(pathSuffix, mountPath)
	}

	if skipIfExists {
		if strings.HasPrefix(pathSuffix, apiPrefix) || strings.HasPrefix(pathSuffix, "/"+apiPrefix) {
			return paths.Join(mountPath, pathSuffix)
		}
	}

	return paths.Join(mountPath, apiPrefix, pathSuffix)
}

func getHeaderForMap(header string, data map[string]interface{}) string {
	maxKey := 0
	for k := range data {
		if len(k) > maxKey {
			maxKey = len(k)
		}
	}

	// 4 for the column spaces and 5 for the len("value")
	totalLen := maxKey + 4 + 5

	return padEqualSigns(header, totalLen)
}

func kvParseVersionsFlags(versions []string) []string {
	versionsOut := make([]string, 0, len(versions))
	for _, v := range versions {
		versionsOut = append(versionsOut, strutil.ParseStringSlice(v, ",")...)
	}

	return versionsOut
}

func outputPath(ui cli.Ui, path string, title string) {
	ui.Info(padEqualSigns(title, len(path)))
	ui.Info(path)
	ui.Info("")
}

// Pad the table header with equal signs on each side
func padEqualSigns(header string, totalLen int) string {
	equalSigns := totalLen - (len(header) + 2)

	// If we have zero or fewer equal signs bump it back up to two on either
	// side of the header.
	if equalSigns <= 0 {
		equalSigns = 4
	}

	// If the number of equal signs is not divisible by two add a sign.
	if equalSigns%2 != 0 {
		equalSigns = equalSigns + 1
	}

	return fmt.Sprintf("%s %s %s", strings.Repeat("=", equalSigns/2), header, strings.Repeat("=", equalSigns/2))
}

// WalkSecretsTree dfs-traverses the secrets tree rooted at the given path
// and calls the `visit` functor for each of the directory and leaf paths.
// Note: for kv-v2, a "metadata" path is expected and "metadata" paths will be
// returned in the visit functor.
func WalkSecretsTree(ctx context.Context, client *api.Client, path string, visit func(path string, directory bool) error) error {
	resp, err := client.Logical().ListWithContext(ctx, path)
	if err != nil {
		return fmt.Errorf("could not list %q path: %w", path, err)
	}

	if resp == nil || resp.Data == nil {
		return fmt.Errorf("no value found at %q: %w", path, err)
	}

	keysRaw, ok := resp.Data["keys"]
	if !ok {
		return fmt.Errorf("unexpected list response at %q", path)
	}

	keysRawSlice, ok := keysRaw.([]interface{})
	if !ok {
		return fmt.Errorf("unexpected list response type %T at %q", keysRaw, path)
	}

	keys := make([]string, 0, len(keysRawSlice))

	for _, keyRaw := range keysRawSlice {
		key, ok := keyRaw.(string)
		if !ok {
			return fmt.Errorf("unexpected key type %T at %q", keyRaw, path)
		}
		keys = append(keys, key)
	}

	// sort the keys for a deterministic output
	sort.Strings(keys)

	for _, key := range keys {
		// the keys are relative to the current path: combine them
		child := paths.Join(path, key)

		if strings.HasSuffix(key, "/") {
			// visit the directory
			if err := visit(child, true); err != nil {
				return err
			}

			// this is not a leaf node: we need to go deeper...
			if err := WalkSecretsTree(ctx, client, child, visit); err != nil {
				return err
			}
		} else {
			// this is a leaf node: add it to the list
			if err := visit(child, false); err != nil {
				return err
			}
		}
	}

	return nil
}
