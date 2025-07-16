// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package namespaces

import (
	"context"
	"errors"
	"fmt"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
)

// RootNamespacePath is the path of the root namespace.
const RootNamespacePath = ""

// RootNamespaceID is the ID of the root namespace.
const RootNamespaceID = "root"

// ErrNotFound is returned by funcs in this package when something isn't found,
// instead of returning (nil, nil).
var ErrNotFound = errors.New("no namespaces found")

// folderPath transforms an input path that refers to a namespace or mount point,
// such that it adheres to the norms Vault prefers.  The result will have any
// leading "/" stripped, and, except for the root namespace which is always
// RootNamespacePath, will always end in a "/".
func folderPath(path string) string {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return strings.TrimPrefix(path, "/")
}

// joinPath concatenates its inputs using "/" as a delimiter.  The result will
// adhere to folderPath conventions.
func joinPath(s ...string) string {
	return folderPath(path.Join(s...))
}

// GetNamespaceIDPaths does a namespace list and extracts the resulting paths
// and namespace IDs, returning a map from namespace ID to path.  Returns
// ErrNotFound if no namespaces exist beneath the current namespace set on the
// client.
func GetNamespaceIDPaths(client *api.Client) (map[string]string, error) {
	secret, err := client.Logical().List("sys/namespaces")
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, ErrNotFound
	}
	if _, ok := secret.Data["key_info"]; !ok {
		return nil, ErrNotFound
	}

	ret := map[string]string{}
	for relNsPath, infoAny := range secret.Data["key_info"].(map[string]any) {
		info := infoAny.(map[string]any)
		id := info["id"].(string)
		ret[id] = relNsPath
	}
	return ret, err
}

// WalkNamespaces does recursive namespace list commands to discover the complete
// namespace hierarchy.  This may yield an error or inconsistent results if
// namespaces change while we're querying them.
// The callback f is invoked for every namespace discovered.  Namespace traversal
// is pre-order depth-first. If f returns an error, traversal is aborted and the
// error is returned.  Otherwise, an error is only returned if a request results
// in an error.
func WalkNamespaces(client *api.Client, f func(id, apiPath string) error) error {
	return walkNamespacesRecursive(client, RootNamespaceID, RootNamespacePath, f)
}

func walkNamespacesRecursive(client *api.Client, startID, startApiPath string, f func(id, apiPath string) error) error {
	if err := f(startID, startApiPath); err != nil {
		return err
	}

	idpaths, err := GetNamespaceIDPaths(client.WithNamespace(startApiPath))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil
		}
		return err
	}

	for id, path := range idpaths {
		fullPath := joinPath(startApiPath, path)

		if err = walkNamespacesRecursive(client, id, fullPath, f); err != nil {
			return err
		}
	}
	return nil
}

// PollDeleteNamespace issues a namespace delete request and waits for it
// to complete (since namespace deletes are asynchronous), at least until
// ctx expires.
func PollDeleteNamespace(ctx context.Context, client *api.Client, nsPath string) error {
	_, err := client.Logical().Delete("sys/namespaces/" + nsPath)
	if err != nil {
		return err
	}

LOOP:
	for ctx.Err() == nil {
		resp, err := client.Logical().Delete("sys/namespaces/" + nsPath)
		if err != nil {
			return err
		}
		for _, warn := range resp.Warnings {
			if strings.HasPrefix(warn, "Namespace is already being deleted") {
				time.Sleep(10 * time.Millisecond)
				continue LOOP
			}
		}
		break
	}

	return ctx.Err()
}

// DeleteAllNamespaces uses WalkNamespaces to delete all namespaces,
// waiting for deletion to complete before returning.  The same caveats about
// namespaces changing underneath us apply as in WalkNamespaces.
// Traversal is depth-first pre-order, but we must do the deletion in the reverse
// order, since a namespace containing namespaces cannot be deleted.
func DeleteAllNamespaces(ctx context.Context, client *api.Client) error {
	var nss []string
	err := WalkNamespaces(client, func(id, apiPath string) error {
		if apiPath != RootNamespacePath {
			nss = append(nss, apiPath)
		}
		return nil
	})
	if err != nil {
		return err
	}
	slices.Reverse(nss)
	for _, apiPath := range nss {
		if err := PollDeleteNamespace(ctx, client, apiPath); err != nil {
			return fmt.Errorf("error deleting namespace %q: %v", apiPath, err)
		}
	}

	// Do a final check to make sure that we got everything, and so that the
	// caller doesn't assume that all namespaces are deleted when a glitch
	// occurred due to namespaces changing while we were traversing or deleting
	// them.
	_, err = GetNamespaceIDPaths(client)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	return nil
}
