// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestSystemBackend_PluginCatalog_Update_Download_Should_Fail tests the update failure
// cases when download is true
func TestSystemBackend_PluginCatalog_Update_Download_Should_Fail(t *testing.T) {
	const expectedErrStr = "download is an enterprise only feature"
	sym, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	c, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{
		PluginDirectory: sym,
	})
	b := c.systemBackend

	tests := []struct {
		pluginType    consts.PluginType
		pluginVersion string
		pluginName    string
	}{
		{
			pluginName:    "vault-plugin-database-redis",
			pluginVersion: "v0.6.0",
			pluginType:    consts.PluginTypeDatabase,
		},
		{
			pluginName:    "vault-plugin-secrets-kv",
			pluginVersion: "v0.24.0",
			pluginType:    consts.PluginTypeSecrets,
		},
		{
			pluginName:    "vault-plugin-auth-jwt",
			pluginVersion: "v0.24.1",
			pluginType:    consts.PluginTypeCredential,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.pluginName, tt.pluginVersion), func(t *testing.T) {
			req := logical.TestRequest(t, logical.UpdateOperation,
				"plugins/catalog/"+tt.pluginType.String()+"/"+tt.pluginName)
			req.Data["version"] = tt.pluginVersion
			req.Data["download"] = true
			resp, err := b.HandleRequest(namespace.RootContext(nil), req)
			if err != nil || resp.Error() == nil {
				t.Fatalf("expected error when download is true, got resp: %v, err: %v", resp, err)
			} else if !strings.Contains(resp.Error().Error(), expectedErrStr) {
				t.Fatalf("expected error %q, got resp: %v", expectedErrStr, resp)
			}
		})
	}
}
