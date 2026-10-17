// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

// Package exactdeny contains regression tests for SECVULN-57761: an ACL
// exact-deny bypass via case-variant resource names. Several backends
// (AppRole, the ACL policy store, userpass, and identity) resolve
// resource names case-insensitively (lowercasing them before storage
// lookup/mutation), while vault/acl.go authorizes the raw, caller-supplied
// request path. Without normalizing the path before the ACL check, a token
// holding a wildcard allow and a more specific exact deny on the canonical
// (lowercase) resource name could bypass the deny by requesting a
// differently-cased variant: the variant misses the exact deny, matches the
// wildcard allow, and the backend then silently resolves it back to the
// protected canonical resource.
package exactdeny

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// newScopedClient creates a token with the given inline ACL policy and
// returns a client authenticated as that token, leaving the root client
// untouched for setup/verification.
func newScopedClient(t *testing.T, root *api.Client, policyName, policyHCL string) *api.Client {
	t.Helper()

	require.NoError(t, root.Sys().PutPolicy(policyName, policyHCL))

	secret, err := root.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{policyName},
		TTL:      "1h",
	})
	require.NoError(t, err)
	require.NotNil(t, secret.Auth)

	scoped, err := root.Clone()
	require.NoError(t, err)
	scoped.SetToken(secret.Auth.ClientToken)
	return scoped
}

func requirePermissionDenied(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	require.Contains(t, err.Error(), "permission denied")
}

// generateSelfSignedCertPEM creates a throwaway self-signed CA certificate
// (PEM-encoded) suitable for registering as a trusted cert auth entry in
// tests; its content is otherwise irrelevant to what's under test here.
func generateSelfSignedCertPEM() (string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "exactdeny-test"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// generateSelfSignedCRLPEM creates a throwaway, empty CRL (PEM-encoded),
// signed by its own throwaway CA, suitable for registering as a trusted CRL
// in cert auth tests; its content is otherwise irrelevant to what's under
// test here.
func generateSelfSignedCRLPEM() (string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	caTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "exactdeny-test-crl-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		SubjectKeyId:          []byte{1, 2, 3, 4},
	}

	caDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &key.PublicKey, key)
	if err != nil {
		return "", err
	}
	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		return "", err
	}

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: time.Now().Add(-time.Hour),
		NextUpdate: time.Now().Add(time.Hour),
	}, caCert, key)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, &pem.Block{Type: "X509 CRL", Bytes: crlDER}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// TestAppRole_ExactDenyNotBypassedByCaseVariant verifies that a delegated
// token with wildcard access to a custom-mounted AppRole backend's role/*
// cannot bypass an exact deny protecting a specific role's
// role-id/secret-id endpoints by requesting a differently-cased variant of
// the protected role name. The backend is mounted at a non-default path so
// this exercises the APIPathNoNamespace-based mount-relative matching that
// normalizeCaseInsensitiveResourcePath needs for backends that can be
// mounted anywhere, as opposed to the fixed-prefix identity/ACL paths.
func TestAppRole_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	const mount = "custom-approle"
	require.NoError(t, root.Sys().EnableAuthWithOptions(mount, &api.EnableAuthOptions{Type: "approle"}))

	// The protected role, plus an unrelated role to confirm wildcard access
	// still works for non-protected resources.
	_, err := root.Logical().Write(fmt.Sprintf("auth/%s/role/admin", mount), map[string]interface{}{})
	require.NoError(t, err)
	_, err = root.Logical().Write(fmt.Sprintf("auth/%s/role/other", mount), map[string]interface{}{})
	require.NoError(t, err)

	policy := fmt.Sprintf(`
path "auth/%[1]s/role/*" {
  capabilities = ["read", "create", "update", "list"]
}
path "auth/%[1]s/role/admin/role-id" {
  capabilities = ["deny"]
}
path "auth/%[1]s/role/admin/secret-id" {
  capabilities = ["deny"]
}
`, mount)
	scoped := newScopedClient(t, root, "approle-wildcard-deny-admin", policy)

	// Baseline: the exact deny protects the canonical (lowercase) role name.
	_, err = scoped.Logical().Read(fmt.Sprintf("auth/%s/role/admin/role-id", mount))
	requirePermissionDenied(t, err)
	_, err = scoped.Logical().Write(fmt.Sprintf("auth/%s/role/admin/secret-id", mount), map[string]interface{}{})
	requirePermissionDenied(t, err)

	// The bypass: a case variant of the protected role name must also be
	// denied, even though AppRole resolves it to the same canonical role.
	for _, variant := range []string{"ADMIN", "Admin", "aDmIn"} {
		_, err = scoped.Logical().Read(fmt.Sprintf("auth/%s/role/%s/role-id", mount, variant))
		requirePermissionDenied(t, err)

		_, err = scoped.Logical().Write(fmt.Sprintf("auth/%s/role/%s/secret-id", mount, variant), map[string]interface{}{})
		requirePermissionDenied(t, err)
	}

	// Wildcard access to non-protected roles must continue to work.
	resp, err := scoped.Logical().Read(fmt.Sprintf("auth/%s/role/other/role-id", mount))
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestAWS_ExactDenyNotBypassedByCaseVariant verifies that a delegated token
// with wildcard access to a mounted AWS auth backend's role/* cannot bypass
// an exact deny protecting a specific role by requesting a differently-cased
// variant of that role name, which the backend would otherwise silently
// canonicalize back onto the protected role (builtin/credential/aws/path_role.go
// lowercases the role name before every storage lookup/write).
func TestAWS_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	require.NoError(t, root.Sys().EnableAuthWithOptions("aws", &api.EnableAuthOptions{Type: "aws"}))

	_, err := root.Logical().Write("auth/aws/role/admin", map[string]interface{}{
		"auth_type":               "iam",
		"bound_iam_principal_arn": "arn:aws:iam::123456789012:role/admin*",
	})
	require.NoError(t, err)
	_, err = root.Logical().Write("auth/aws/role/other", map[string]interface{}{
		"auth_type":               "iam",
		"bound_iam_principal_arn": "arn:aws:iam::123456789012:role/other*",
	})
	require.NoError(t, err)

	policy := `
path "auth/aws/role/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/aws/role/admin" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "aws-wildcard-deny-admin", policy)

	// Baseline: the exact deny protects the canonical (lowercase) role name.
	_, err = scoped.Logical().Read("auth/aws/role/admin")
	requirePermissionDenied(t, err)

	// The bypass: a case variant of the protected role name must also be
	// denied, even though AWS auth resolves it to the same canonical role.
	for _, variant := range []string{"ADMIN", "Admin"} {
		_, err = scoped.Logical().Read(fmt.Sprintf("auth/aws/role/%s", variant))
		requirePermissionDenied(t, err)
	}

	// Wildcard access to non-protected roles must continue to work.
	resp, err := scoped.Logical().Read("auth/aws/role/other")
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestAWSEC2Alias_ExactDenyNotBypassedByCaseVariant verifies that the
// case-variant exact-deny protection also covers a mount created with the
// legacy "aws-ec2" auth type. credentialAliases (vault/auth.go) maps
// "aws-ec2" onto the "aws" backend when instantiating it, but the mount
// keeps "aws-ec2" as its MountEntry.Type, so an authorization check that
// keys off the raw mount type would skip normalization for these mounts
// even though they run the same role-name-lowercasing AWS code.
func TestAWSEC2Alias_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	require.NoError(t, root.Sys().EnableAuthWithOptions("aws-ec2", &api.EnableAuthOptions{Type: "aws-ec2"}))

	_, err := root.Logical().Write("auth/aws-ec2/role/admin", map[string]interface{}{
		"auth_type":               "iam",
		"bound_iam_principal_arn": "arn:aws:iam::123456789012:role/admin*",
	})
	require.NoError(t, err)
	_, err = root.Logical().Write("auth/aws-ec2/role/other", map[string]interface{}{
		"auth_type":               "iam",
		"bound_iam_principal_arn": "arn:aws:iam::123456789012:role/other*",
	})
	require.NoError(t, err)

	policy := `
path "auth/aws-ec2/role/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/aws-ec2/role/admin" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "aws-ec2-wildcard-deny-admin", policy)

	// Baseline: the exact deny protects the canonical (lowercase) role name.
	_, err = scoped.Logical().Read("auth/aws-ec2/role/admin")
	requirePermissionDenied(t, err)

	// The bypass: a case variant of the protected role name must also be
	// denied, even though the aliased backend resolves it to the same
	// canonical role.
	for _, variant := range []string{"ADMIN", "Admin"} {
		_, err = scoped.Logical().Read(fmt.Sprintf("auth/aws-ec2/role/%s", variant))
		requirePermissionDenied(t, err)
	}

	// Wildcard access to non-protected roles must continue to work.
	resp, err := scoped.Logical().Read("auth/aws-ec2/role/other")
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestCert_ExactDenyNotBypassedByCaseVariant verifies that a delegated token
// with wildcard access to a mounted cert auth backend's certs/* and crls/*
// cannot bypass an exact deny protecting a specific trusted certificate or
// CRL by requesting a differently-cased variant of its name, which the
// backend would otherwise silently canonicalize back onto the protected
// resource (builtin/credential/cert/path_certs.go,
// builtin/credential/cert/path_crls.go lowercase the name before every
// storage lookup/write).
func TestCert_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	require.NoError(t, root.Sys().EnableAuthWithOptions("cert", &api.EnableAuthOptions{Type: "cert"}))

	certPEM, err := generateSelfSignedCertPEM()
	require.NoError(t, err)

	_, err = root.Logical().Write("auth/cert/certs/admin", map[string]interface{}{"certificate": certPEM})
	require.NoError(t, err)
	_, err = root.Logical().Write("auth/cert/certs/other", map[string]interface{}{"certificate": certPEM})
	require.NoError(t, err)

	policy := `
path "auth/cert/certs/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/cert/certs/admin" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "cert-wildcard-deny-admin", policy)

	// Baseline: the exact deny protects the canonical (lowercase) cert name.
	_, err = scoped.Logical().Read("auth/cert/certs/admin")
	requirePermissionDenied(t, err)

	// The bypass: a case variant of the protected cert name must also be
	// denied, even though cert auth resolves it to the same canonical entry.
	for _, variant := range []string{"ADMIN", "Admin"} {
		_, err = scoped.Logical().Read(fmt.Sprintf("auth/cert/certs/%s", variant))
		requirePermissionDenied(t, err)
	}

	// Wildcard access to non-protected certs must continue to work.
	resp, err := scoped.Logical().Read("auth/cert/certs/other")
	require.NoError(t, err)
	require.NotNil(t, resp)

	// The same exact-deny bypass applies independently to crls/*: register
	// two distinct CRLs and verify a case variant of the protected CRL name
	// is also denied.
	crlPEM, err := generateSelfSignedCRLPEM()
	require.NoError(t, err)

	_, err = root.Logical().Write("auth/cert/crls/admin", map[string]interface{}{"crl": crlPEM})
	require.NoError(t, err)
	_, err = root.Logical().Write("auth/cert/crls/other", map[string]interface{}{"crl": crlPEM})
	require.NoError(t, err)

	crlPolicy := `
path "auth/cert/crls/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/cert/crls/admin" {
  capabilities = ["deny"]
}
`
	crlScoped := newScopedClient(t, root, "cert-crl-wildcard-deny-admin", crlPolicy)

	_, err = crlScoped.Logical().Read("auth/cert/crls/admin")
	requirePermissionDenied(t, err)

	for _, variant := range []string{"ADMIN", "Admin"} {
		_, err = crlScoped.Logical().Read(fmt.Sprintf("auth/cert/crls/%s", variant))
		requirePermissionDenied(t, err)
	}

	crlResp, err := crlScoped.Logical().Read("auth/cert/crls/other")
	require.NoError(t, err)
	require.NotNil(t, crlResp)
}

// TestACLPolicyName_ExactDenyNotBypassedByCaseVariant verifies that a
// delegated token with wildcard access to sys/policies/acl/* cannot bypass
// an exact deny protecting a specific policy name by writing a
// differently-cased variant of that name, which the policy store would
// otherwise silently canonicalize back onto the protected policy.
func TestACLPolicyName_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	protectedRules := `path "secret/protected" { capabilities = ["read"] }`
	require.NoError(t, root.Sys().PutPolicy("admin", protectedRules))

	policy := `
path "sys/policies/acl/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "sys/policies/acl/admin" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "policy-wildcard-deny-admin", policy)

	// Baseline: writing the canonical protected policy name is denied.
	err := scoped.Sys().PutPolicy("admin", `path "secret/attacker" { capabilities = ["read"] }`)
	requirePermissionDenied(t, err)

	// The bypass: a case variant of the protected policy name must also be
	// denied, and must not silently overwrite the protected policy. This
	// also covers whitespace-padded variants, since PolicyStore.sanitizeName
	// trims surrounding whitespace in addition to lowercasing.
	for _, variant := range []string{"ADMIN", "Admin", " admin", "admin ", " Admin "} {
		err = scoped.Sys().PutPolicy(variant, `path "secret/attacker" { capabilities = ["read"] }`)
		requirePermissionDenied(t, err)
	}

	// The protected policy must be unchanged.
	rules, err := root.Sys().GetPolicy("admin")
	require.NoError(t, err)
	require.Contains(t, rules, "secret/protected")
	require.NotContains(t, rules, "secret/attacker")

	// Wildcard access to non-protected policy names must continue to work.
	require.NoError(t, scoped.Sys().PutPolicy("other", `path "secret/other" { capabilities = ["read"] }`))
}

// TestLegacyACLPolicyName_ExactDenyNotBypassedByCaseVariant is the
// equivalent of TestACLPolicyName_ExactDenyNotBypassedByCaseVariant for the
// legacy sys/policy/<name> endpoint. sys/policy/ and sys/policies/acl/ are
// normalized by separate prefixes in normalizeCaseInsensitiveResourcePath,
// so each needs its own regression coverage to catch a typo or regression
// in either branch. api.Sys().PutPolicy/GetPolicy always target
// sys/policies/acl/, so this test writes/reads sys/policy/<name> directly
// through Logical() to exercise the legacy endpoint.
func TestLegacyACLPolicyName_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	protectedRules := `path "secret/protected" { capabilities = ["read"] }`
	require.NoError(t, root.Sys().PutPolicy("admin", protectedRules))

	policy := `
path "sys/policy/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "sys/policy/admin" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "legacy-policy-wildcard-deny-admin", policy)

	// Baseline: writing the canonical protected policy name is denied.
	_, err := scoped.Logical().Write("sys/policy/admin", map[string]interface{}{"rules": `path "secret/attacker" { capabilities = ["read"] }`})
	requirePermissionDenied(t, err)

	// The bypass: a case variant of the protected policy name must also be
	// denied, and must not silently overwrite the protected policy. This
	// also covers whitespace-padded variants, since PolicyStore.sanitizeName
	// trims surrounding whitespace in addition to lowercasing.
	for _, variant := range []string{"ADMIN", "Admin", " admin", "admin ", " Admin "} {
		_, err = scoped.Logical().Write(fmt.Sprintf("sys/policy/%s", variant), map[string]interface{}{"rules": `path "secret/attacker" { capabilities = ["read"] }`})
		requirePermissionDenied(t, err)
	}

	// The protected policy must be unchanged.
	rules, err := root.Sys().GetPolicy("admin")
	require.NoError(t, err)
	require.Contains(t, rules, "secret/protected")
	require.NotContains(t, rules, "secret/attacker")

	// Wildcard access to non-protected policy names must continue to work.
	_, err = scoped.Logical().Write("sys/policy/other", map[string]interface{}{"rules": `path "secret/other" { capabilities = ["read"] }`})
	require.NoError(t, err)
}

// TestUserpass_ExactDenyNotBypassedByCaseVariant verifies that a delegated
// token with wildcard access to auth/userpass/users/* cannot bypass an
// exact deny protecting a specific user's password endpoint by requesting a
// differently-cased variant of the protected username.
func TestUserpass_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	require.NoError(t, root.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{Type: "userpass"}))

	_, err := root.Logical().Write("auth/userpass/users/admin", map[string]interface{}{"password": "orig-password"})
	require.NoError(t, err)
	_, err = root.Logical().Write("auth/userpass/users/other", map[string]interface{}{"password": "orig-password"})
	require.NoError(t, err)

	policy := `
path "auth/userpass/users/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/userpass/users/admin/password" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "userpass-wildcard-deny-admin", policy)

	// Baseline: the exact deny protects the canonical (lowercase) username.
	_, err = scoped.Logical().Write("auth/userpass/users/admin/password", map[string]interface{}{"password": "attacker-password"})
	requirePermissionDenied(t, err)

	// The bypass: a case variant of the protected username must also be
	// denied, even though userpass resolves it to the same canonical user.
	for _, variant := range []string{"ADMIN", "Admin"} {
		_, err = scoped.Logical().Write(fmt.Sprintf("auth/userpass/users/%s/password", variant), map[string]interface{}{"password": "attacker-password"})
		requirePermissionDenied(t, err)
	}

	// Wildcard access to non-protected users must continue to work.
	_, err = scoped.Logical().Write("auth/userpass/users/other/password", map[string]interface{}{"password": "new-password"})
	require.NoError(t, err)
}

// TestIdentityGroup_ExactDenyNotBypassedByCaseVariant verifies that a
// delegated token with wildcard access to identity/group/name/* cannot
// bypass an exact deny protecting a specific group by requesting a
// differently-cased variant of the protected group name. This is a
// regression test for VAULT-34957, which fixed the identity case but had
// no direct test coverage.
func TestIdentityGroup_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	_, err := root.Logical().Write("identity/group/name/admin", map[string]interface{}{
		"metadata": map[string]interface{}{"role": "protected"},
	})
	require.NoError(t, err)
	_, err = root.Logical().Write("identity/group/name/other", map[string]interface{}{})
	require.NoError(t, err)

	policy := `
path "identity/group/name/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "identity/group/name/admin" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "identity-wildcard-deny-admin", policy)

	// Baseline: the exact deny protects the canonical (lowercase) group name.
	_, err = scoped.Logical().Write("identity/group/name/admin", map[string]interface{}{
		"metadata": map[string]interface{}{"role": "compromised"},
	})
	requirePermissionDenied(t, err)

	// The bypass: a case variant of the protected group name must also be
	// denied, even though identity resolves it to the same canonical group.
	for _, variant := range []string{"ADMIN", "Admin"} {
		_, err = scoped.Logical().Write(fmt.Sprintf("identity/group/name/%s", variant), map[string]interface{}{
			"metadata": map[string]interface{}{"role": "compromised"},
		})
		requirePermissionDenied(t, err)
	}

	// The protected group's metadata must be unchanged.
	resp, err := root.Logical().Read("identity/group/name/admin")
	require.NoError(t, err)
	require.NotNil(t, resp)
	metadata, ok := resp.Data["metadata"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "protected", metadata["role"])

	// Wildcard access to non-protected groups must continue to work.
	_, err = scoped.Logical().Write("identity/group/name/other", map[string]interface{}{
		"metadata": map[string]interface{}{"role": "unprotected"},
	})
	require.NoError(t, err)
}

// TestAzure_ExactDenyNotBypassedByCaseVariant verifies that an exact deny on
// an Azure auth role name cannot be bypassed by requesting a case variant.
// The Azure backend lowercases the role name before every storage read,
// write, and delete (vault-plugin-auth-azure path_role.go), so "Admin" and
// "admin" address the same stored role.
func TestAzure_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	require.NoError(t, root.Sys().EnableAuthWithOptions("azure", &api.EnableAuthOptions{Type: "azure"}))

	for _, role := range []string{"admin", "other"} {
		_, err := root.Logical().Write("auth/azure/role/"+role, map[string]interface{}{
			"bound_subscription_ids": []string{"11111111-2222-3333-4444-555555555555"},
		})
		require.NoError(t, err)
	}

	policy := `
path "auth/azure/role/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/azure/role/admin" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "azure-wildcard-deny-admin", policy)

	_, err := scoped.Logical().Read("auth/azure/role/admin")
	requirePermissionDenied(t, err)

	for _, variant := range []string{"ADMIN", "Admin"} {
		_, err = scoped.Logical().Read(fmt.Sprintf("auth/azure/role/%s", variant))
		requirePermissionDenied(t, err)
	}

	resp, err := scoped.Logical().Read("auth/azure/role/other")
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestGCP_ExactDenyNotBypassedByCaseVariant verifies that an exact deny on a
// GCP auth role name cannot be bypassed by requesting a case variant. The
// GCP backend lowercases the role name in both pathRoleCreateUpdate and
// role() (vault-plugin-auth-gcp plugin/path_role.go), so "Admin" and "admin"
// address the same stored role. The sub-resource case additionally covers
// role/<name>/service-accounts, where only the name segment may be rewritten.
func TestGCP_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	require.NoError(t, root.Sys().EnableAuthWithOptions("gcp", &api.EnableAuthOptions{Type: "gcp"}))

	for _, role := range []string{"admin", "other"} {
		_, err := root.Logical().Write("auth/gcp/role/"+role, map[string]interface{}{
			"type":                   "iam",
			"bound_service_accounts": []string{"test@project.iam.gserviceaccount.com"},
		})
		require.NoError(t, err)
	}

	policy := `
path "auth/gcp/role/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/gcp/role/admin" {
  capabilities = ["deny"]
}
path "auth/gcp/role/admin/service-accounts" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "gcp-wildcard-deny-admin", policy)

	_, err := scoped.Logical().Read("auth/gcp/role/admin")
	requirePermissionDenied(t, err)

	for _, variant := range []string{"ADMIN", "Admin"} {
		_, err = scoped.Logical().Read(fmt.Sprintf("auth/gcp/role/%s", variant))
		requirePermissionDenied(t, err)

		// The name segment must be normalized even when a sub-resource
		// follows it.
		_, err = scoped.Logical().Write(fmt.Sprintf("auth/gcp/role/%s/service-accounts", variant), map[string]interface{}{
			"add": []string{"added@project.iam.gserviceaccount.com"},
		})
		requirePermissionDenied(t, err)
	}

	resp, err := scoped.Logical().Read("auth/gcp/role/other")
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestKubernetes_ExactDenyNotBypassedByCaseVariant verifies that an exact
// deny on a Kubernetes auth role name cannot be bypassed by requesting a
// case variant. The Kubernetes backend lowercases the role name before every
// storage read, write, and delete (vault-plugin-auth-kubernetes backend.go
// and path_role.go), so "Admin" and "admin" address the same stored role.
func TestKubernetes_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	require.NoError(t, root.Sys().EnableAuthWithOptions("kubernetes", &api.EnableAuthOptions{Type: "kubernetes"}))

	for _, role := range []string{"admin", "other"} {
		_, err := root.Logical().Write("auth/kubernetes/role/"+role, map[string]interface{}{
			"bound_service_account_names":      []string{"vault-auth"},
			"bound_service_account_namespaces": []string{"default"},
		})
		require.NoError(t, err)
	}

	policy := `
path "auth/kubernetes/role/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/kubernetes/role/admin" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "kubernetes-wildcard-deny-admin", policy)

	_, err := scoped.Logical().Read("auth/kubernetes/role/admin")
	requirePermissionDenied(t, err)

	for _, variant := range []string{"ADMIN", "Admin"} {
		_, err = scoped.Logical().Read(fmt.Sprintf("auth/kubernetes/role/%s", variant))
		requirePermissionDenied(t, err)
	}

	resp, err := scoped.Logical().Read("auth/kubernetes/role/other")
	require.NoError(t, err)
	require.NotNil(t, resp)
}

// TestGitHub_ExactDenyNotBypassedByCaseVariant verifies that an exact deny on
// a GitHub auth policy mapping cannot be bypassed by requesting a case
// variant. The github backend exposes its team and user mappings through
// framework.PathMap without setting CaseSensitive, and PathMap.pathStruct
// (sdk/framework/path_map.go) lowercases the key for every read, write, and
// delete, so "Admins" and "admins" address the same stored mapping. These
// paths assign Vault policies, so a bypass here grants the ability to
// rewrite a protected team's policy set.
func TestGitHub_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	require.NoError(t, root.Sys().EnableAuthWithOptions("github", &api.EnableAuthOptions{Type: "github"}))

	for _, team := range []string{"admins", "others"} {
		_, err := root.Logical().Write("auth/github/map/teams/"+team, map[string]interface{}{
			"value": "default",
		})
		require.NoError(t, err)
	}
	_, err := root.Logical().Write("auth/github/map/users/bob", map[string]interface{}{
		"value": "default",
	})
	require.NoError(t, err)

	policy := `
path "auth/github/map/teams/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/github/map/teams/admins" {
  capabilities = ["deny"]
}
path "auth/github/map/users/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "auth/github/map/users/bob" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "github-wildcard-deny-admins", policy)

	// Baseline: the exact denies protect the canonical (lowercase) keys.
	_, err = scoped.Logical().Read("auth/github/map/teams/admins")
	requirePermissionDenied(t, err)
	_, err = scoped.Logical().Read("auth/github/map/users/bob")
	requirePermissionDenied(t, err)

	// The bypass: case variants resolve to the same stored mapping, so they
	// must be denied too. Writing is the dangerous operation here, since it
	// would let the caller assign policies to a protected team.
	for _, variant := range []string{"ADMINS", "Admins"} {
		_, err = scoped.Logical().Read(fmt.Sprintf("auth/github/map/teams/%s", variant))
		requirePermissionDenied(t, err)

		_, err = scoped.Logical().Write(fmt.Sprintf("auth/github/map/teams/%s", variant), map[string]interface{}{
			"value": "root",
		})
		requirePermissionDenied(t, err)
	}
	for _, variant := range []string{"BOB", "Bob"} {
		_, err = scoped.Logical().Read(fmt.Sprintf("auth/github/map/users/%s", variant))
		requirePermissionDenied(t, err)
	}

	// Wildcard access to non-protected mappings must continue to work.
	resp, err := scoped.Logical().Read("auth/github/map/teams/others")
	require.NoError(t, err)
	require.NotNil(t, resp)

	// The protected mapping must still hold its original value.
	resp, err = root.Logical().Read("auth/github/map/teams/admins")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "default", resp.Data["value"])
}

// TestIdentityEntity_ExactDenyNotBypassedByCaseVariant verifies that a
// delegated token with wildcard access to identity/entity/name/* cannot
// bypass an exact deny protecting a specific entity by requesting a
// differently-cased variant of the protected entity name. This mirrors
// TestIdentityGroup_ExactDenyNotBypassedByCaseVariant but for identity
// entities, ensuring both identity resource types have matching regression
// coverage.
func TestIdentityEntity_ExactDenyNotBypassedByCaseVariant(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	_, err := root.Logical().Write("identity/entity/name/admin", map[string]interface{}{
		"metadata": map[string]interface{}{"role": "protected"},
	})
	require.NoError(t, err)
	_, err = root.Logical().Write("identity/entity/name/other", map[string]interface{}{})
	require.NoError(t, err)

	policy := `
path "identity/entity/name/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "identity/entity/name/admin" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "identity-entity-wildcard-deny-admin", policy)

	// Baseline: the exact deny protects the canonical (lowercase) entity name.
	_, err = scoped.Logical().Write("identity/entity/name/admin", map[string]interface{}{
		"metadata": map[string]interface{}{"role": "compromised"},
	})
	requirePermissionDenied(t, err)

	// The bypass: a case variant of the protected entity name must also be
	// denied, even though identity resolves it to the same canonical entity.
	for _, variant := range []string{"ADMIN", "Admin"} {
		_, err = scoped.Logical().Write(fmt.Sprintf("identity/entity/name/%s", variant), map[string]interface{}{
			"metadata": map[string]interface{}{"role": "compromised"},
		})
		requirePermissionDenied(t, err)
	}

	// The protected entity's metadata must be unchanged.
	resp, err := root.Logical().Read("identity/entity/name/admin")
	require.NoError(t, err)
	require.NotNil(t, resp)
	metadata, ok := resp.Data["metadata"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "protected", metadata["role"])

	// Wildcard access to non-protected entities must continue to work.
	_, err = scoped.Logical().Write("identity/entity/name/other", map[string]interface{}{
		"metadata": map[string]interface{}{"role": "unprotected"},
	})
	require.NoError(t, err)
}

// TestIdentityEntity_ExactDenyWithNonCanonicalInitialName verifies that
// exact-deny protection works correctly when the entity is initially created
// with a non-canonical (uppercase) name. This tests the opposite direction:
// creating "Admin", setting a deny on "Admin", and verifying that both the
// canonical lowercase form and other case variants are denied.
func TestIdentityEntity_ExactDenyWithNonCanonicalInitialName(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	root := cluster.Cores[0].Client

	// Create entities with non-canonical uppercase names.
	_, err := root.Logical().Write("identity/entity/name/Admin", map[string]interface{}{
		"metadata": map[string]interface{}{"role": "protected"},
	})
	require.NoError(t, err)
	_, err = root.Logical().Write("identity/entity/name/Other", map[string]interface{}{})
	require.NoError(t, err)

	// Deny at the canonical lowercase name; policy paths are not normalized,
	// so the deny must be on lowercase "admin" even though the entity was
	// created as "Admin". Requests to "Admin" get normalized to "admin"
	// before matching against the policy.
	policy := `
path "identity/entity/name/*" {
  capabilities = ["read", "create", "update", "list", "delete"]
}
path "identity/entity/name/admin" {
  capabilities = ["deny"]
}
`
	scoped := newScopedClient(t, root, "identity-entity-noncanon-deny-admin", policy)

	// Baseline: the exact deny on the uppercase name is enforced.
	_, err = scoped.Logical().Write("identity/entity/name/Admin", map[string]interface{}{
		"metadata": map[string]interface{}{"role": "compromised"},
	})
	requirePermissionDenied(t, err)

	// Case variants must also be denied, even though the protected entity
	// was created with an uppercase name.
	for _, variant := range []string{"admin", "ADMIN", "aDmIn"} {
		_, err = scoped.Logical().Write(fmt.Sprintf("identity/entity/name/%s", variant), map[string]interface{}{
			"metadata": map[string]interface{}{"role": "compromised"},
		})
		requirePermissionDenied(t, err)
	}

	// The protected entity's metadata must be unchanged.
	resp, err := root.Logical().Read("identity/entity/name/Admin")
	require.NoError(t, err)
	require.NotNil(t, resp)
	metadata, ok := resp.Data["metadata"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "protected", metadata["role"])

	// Wildcard access to non-protected entities must continue to work.
	_, err = scoped.Logical().Write("identity/entity/name/Other", map[string]interface{}{
		"metadata": map[string]interface{}{"role": "unprotected"},
	})
	require.NoError(t, err)
}
