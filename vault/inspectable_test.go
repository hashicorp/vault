package vault

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestInspectRouter(t *testing.T) {
	// Verify that all the expected tables exist when we inspect the router
	c, _, root := TestCoreUnsealed(t)

	rootCtx := namespace.RootContext(nil)
	subTrees := map[string][]string{
		"routeEntry": {"root", "storage"},
		"mountEntry": {"uuid", "accessor"},
	}
	subTreeFields := map[string][]string{
		"routeEntry": {"tainted", "storage_prefix", "accessor", "mount_namespace", "mount_path", "mount_type", "uuid"},
		"mountEntry": {"accessor", "mount_namespace", "mount_path", "mount_type", "uuid"},
	}
	for subTreeType, subTreeArray := range subTrees {
		for _, tag := range subTreeArray {
			resp, err := c.HandleRequest(rootCtx, &logical.Request{
				ClientToken: root,
				Operation:   logical.ReadOperation,
				Path:        "sys/internal/inspect/router/" + tag,
			})
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("bad: resp: %#v\n err: %v", resp, err)
			}
			// Verify that data exists
			data, ok := resp.Data[tag].([]map[string]interface{})
			if !ok {
				t.Fatalf("Router data is malformed")
			}
			for _, entry := range data {
				for _, field := range subTreeFields[subTreeType] {
					if _, ok := entry[field]; !ok {
						t.Fatalf("Field %s not found in %s", field, tag)
					}
				}
			}

		}
	}
}

func TestInvalidInspectRouterPath(t *testing.T) {
	// Verify that attempting to inspect an invalid tree in the router fails
	core, _, rootToken := testCoreSystemBackend(t)
	rootCtx := namespace.RootContext(nil)
	_, err := core.HandleRequest(rootCtx, &logical.Request{
		ClientToken: rootToken,
		Operation:   logical.ReadOperation,
		Path:        "sys/internal/inspect/router/random",
	})
	if !strings.Contains(err.Error(), logical.ErrUnsupportedPath.Error()) {
		t.Fatal("expected unsupported path error")
	}
}

func TestInspectAPISudoProtect(t *testing.T) {
	// Verify that the Inspect API path is sudo protected
	core, _, rootToken := testCoreSystemBackend(t)
	testMakeServiceTokenViaBackend(t, core.tokenStore, rootToken, "tokenid", "", []string{"secret"})
	rootCtx := namespace.RootContext(nil)
	_, err := core.HandleRequest(rootCtx, &logical.Request{
		ClientToken: "tokenid",
		Operation:   logical.ReadOperation,
		Path:        "sys/internal/inspect/router/root",
	})
	if !strings.Contains(err.Error(), logical.ErrPermissionDenied.Error()) {
		t.Fatal("expected permission denied error")
	}
}
