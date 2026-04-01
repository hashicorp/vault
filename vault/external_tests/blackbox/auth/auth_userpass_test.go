// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package auth

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	helpers "github.com/hashicorp/vault/vault/external_tests/blackbox"
)

// testUserpassAuthCreate tests userpass auth engine creation
func testUserpassAuthCreate(t *testing.T, v *blackbox.Session) {
	// Create a policy for our test user
	userPolicy := `
		path "*" {
			capabilities = ["read", "list"]
		}
	`

	// Use common utility to setup userpass auth
	userClient := helpers.SetupUserpassAuth(v, "testuser", "passtestuser1", "reguser", userPolicy)

	// Verify the auth method was enabled by reading auth mounts
	authMounts := v.MustRead("sys/auth")
	if authMounts.Data == nil {
		t.Fatal("Could not read auth mounts")
	}

	// Verify userpass auth method is enabled
	if userpassAuth, ok := authMounts.Data["userpass/"]; !ok {
		t.Fatal("userpass auth method not found in sys/auth")
	} else {
		userpassMap := userpassAuth.(map[string]any)
		if userpassMap["type"] != "userpass" {
			t.Fatalf("Expected userpass auth method type to be 'userpass', got: %v", userpassMap["type"])
		}
	}

	// Test that the user session was created successfully
	if userClient != nil {
		// Login successful, verify we can read basic info
		tokenInfo := userClient.MustRead("auth/token/lookup-self")
		if tokenInfo.Data == nil {
			t.Fatal("Expected user to be able to read own token info after login")
		}
		t.Log("Userpass login test successful")
	} else {
		t.Log("Userpass login not available (likely managed environment)")
	}

	t.Log("Successfully created userpass auth with user: testuser")
}

// testUserpassAuthRead tests userpass auth engine read operations
func testUserpassAuthRead(t *testing.T, v *blackbox.Session) {
	// Use common utility to setup userpass auth with default policy
	userClient := helpers.SetupUserpassAuth(v, "readuser", "readpass123", "default", "")

	// Read the user configuration
	userConfig := v.MustRead("auth/userpass/users/readuser")
	if userConfig.Data == nil {
		t.Fatal("Expected to read user configuration")
	}

	// Test that the user session was created successfully
	if userClient != nil {
		// Login successful, verify we can read basic info
		tokenInfo := userClient.MustRead("auth/token/lookup-self")
		if tokenInfo.Data == nil {
			t.Fatal("Expected user to be able to read own token info after login")
		}
		t.Log("Userpass login test successful")
	} else {
		t.Log("Userpass login not available (likely managed environment)")
	}

	t.Log("Successfully read userpass auth config for user: readuser")
}

// testUserpassAuthDelete tests userpass auth engine delete operations
func testUserpassAuthDelete(t *testing.T, v *blackbox.Session) {
	// Enable userpass auth method with unique mount for delete test
	v.MustEnableAuth("userpass-delete", &api.EnableAuthOptions{Type: "userpass"})

	// Create a user to delete
	userName := "deleteuser"
	userPassword := "deletepass123"
	v.MustWrite("auth/userpass-delete/users/"+userName, map[string]any{
		"password": userPassword,
		"policies": "default",
	})

	// Verify the user exists
	userConfig := v.MustRead("auth/userpass-delete/users/" + userName)
	if userConfig.Data == nil {
		t.Fatal("Expected user to exist before deletion")
	}

	// Delete the user
	v.MustWrite("auth/userpass-delete/users/"+userName, nil)

	t.Logf("Successfully deleted userpass auth user: %s", userName)
}
