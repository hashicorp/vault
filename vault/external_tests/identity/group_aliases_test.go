// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package identity

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
)

func TestIdentityStore_GroupAliasLocalMount(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// Create a local auth mount
	err := client.Sys().EnableAuthWithOptions("ldap", &api.EnableAuthOptions{
		Type:  "ldap",
		Local: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Extract out the mount accessor for LDAP auth
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	ldapMountAccessor := auths["ldap/"].Accessor

	// Create an external group
	secret, err := client.Logical().Write("identity/group", map[string]interface{}{
		"type": "external",
	})
	if err != nil {
		t.Fatal(err)
	}
	groupID := secret.Data["id"].(string)

	// Attempt to create a group alias against a local mount should fail
	secret, err = client.Logical().Write("identity/group-alias", map[string]interface{}{
		"name":           "testuser",
		"mount_accessor": ldapMountAccessor,
		"canonical_id":   groupID,
	})
	if err == nil {
		t.Fatalf("expected error since mount is local")
	}
}
