package identity

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestIdentityStore_UserLockout(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}
	// tune auth mount
	userlockoutConfig := &api.UserLockoutConfigInput{
		LockoutThreshold:            "3",
		LockoutDuration:             "600",
		LockoutCounterResetDuration: "600",
	}
	err = client.Sys().TuneMount("auth/userpass", api.MountConfigInput{
		UserLockoutConfig: userlockoutConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// read auth tune
	resp, err := client.Logical().Read("sys/auth/userpass/tune")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for reading auth tune")
	}
	_, err = client.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}
	// login failure 1
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 2
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 3
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login : permission denied as user locked out
	_, err = client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err == nil {
		t.Fatal("expected login to fail as user locked out")
	}
	if !strings.Contains(err.Error(), logical.ErrPermissionDenied.Error()) {
		t.Fatalf("expected to see permission denied error as user locked out, got %v", err)
	}
}

func TestIdentityStore_UserFailedLoginMapResetOnSuccess(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}
	// tune auth mount
	userlockoutConfig := &api.UserLockoutConfigInput{
		LockoutThreshold:            "3",
		LockoutDuration:             "600",
		LockoutCounterResetDuration: "600",
	}
	err = client.Sys().TuneMount("auth/userpass", api.MountConfigInput{
		UserLockoutConfig: userlockoutConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// read auth tune
	resp, err := client.Logical().Read("sys/auth/userpass/tune")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for reading auth tune")
	}
	_, err = client.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}
	// login failure 1
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 2
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login with right credentials - successful login
	// userFailedLogin map that contains failed user login info reset
	_, err = client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}
	// login failure 1
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})

	// login failure 2 after successful login
	// error should not be permission denied as user not locked out
	_, err = client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	if err == nil {
		t.Fatal("expected login to fail due to wrong credentials")
	}
	if !strings.Contains(err.Error(), "invalid username or password") {
		t.Fatalf("expected to see invalid username or password error as user is not locked out, got %v", err)
	}
}

func TestIdentityStore_LockoutDurationTest(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatal(err)
	}
	// tune auth mount
	userlockoutConfig := &api.UserLockoutConfigInput{
		LockoutThreshold:            "3",
		LockoutDuration:             "600",
		LockoutCounterResetDuration: "600",
	}
	err = client.Sys().TuneMount("auth/userpass", api.MountConfigInput{
		UserLockoutConfig: userlockoutConfig,
	})
	if err != nil {
		t.Fatal(err)
	}

	// read auth tune
	resp, err := client.Logical().Read("sys/auth/userpass/tune")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected a response for reading auth tune")
	}
	_, err = client.Logical().Write("auth/userpass/users/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err != nil {
		t.Fatal(err)
	}
	// login failure 1
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 2
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login failure 3
	client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "wrongPassword",
	})
	// login : permission denied as user locked out
	_, err = client.Logical().Write("auth/userpass/login/bsmith", map[string]interface{}{
		"password": "training",
	})
	if err == nil {
		t.Fatal("expected login to fail as user locked out")
	}
	if !strings.Contains(err.Error(), logical.ErrPermissionDenied.Error()) {
		t.Fatalf("expected to see permission denied error as user locked out, got %v", err)
	}
}
