// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	metrics "github.com/hashicorp/go-metrics/compat"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/credential/approle"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestRequiresMaterializedTokenState verifies token materialization path
// requirements for enterprise token requests.
func TestRequiresMaterializedTokenState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "token lookup self", path: "auth/token/lookup-self", want: true},
		{name: "token lookup", path: "auth/token/lookup", want: true},
		{name: "capabilities self", path: "sys/capabilities-self", want: true},
		{name: "leases lookup", path: "sys/leases/lookup", want: true},
		{name: "leases lookup prefix", path: "sys/leases/lookup/secret/foo", want: true},
		{name: "leases count", path: "sys/leases/count", want: true},
		{name: "leases list", path: "sys/leases", want: true},
		{name: "cubbyhole", path: "cubbyhole/test", want: true},
		{name: "token renew self excluded", path: "auth/token/renew-self", want: false},
		{name: "leases renew excluded", path: "sys/leases/renew", want: false},
		{name: "unrelated", path: "secret/data/foo", want: false},
		{name: "ui mounts preflight", path: "sys/internal/ui/mounts/secret/data/foo", want: true},
		{name: "ui mounts with nested path", path: "sys/internal/ui/mounts/kv/data/nested/key", want: true},
		{name: "ui mounts prefix only", path: "sys/internal/ui/mounts/", want: true},
		{name: "ui mounts exact", path: "sys/internal/ui/mounts", want: true},
		{name: "ui namespaces exact", path: "sys/internal/ui/namespaces", want: true},
		{name: "ui namespaces with suffix excluded", path: "sys/internal/ui/namespaces/foo", want: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, requiresMaterializedTokenState(tc.path))
		})
	}
}

// TestNormalizeCaseInsensitiveResourcePath verifies that
// normalizeCaseInsensitiveResourcePath lowercases (and, for policy paths,
// trims whitespace from) the resource-name segment of paths handled by
// backends that resolve that name case-insensitively, while leaving
// identity paths untouched once the identity store has entered
// case-sensitive mode (IdentityStore.disableLowerCasedNames). Identity
// names are only guaranteed unique per-case once the store is in that mode
// (entered automatically when a duplicate name is detected across cases),
// so unconditionally lowercasing would incorrectly let a policy scoped to
// "admin" also authorize the distinct group "Admin".
func TestNormalizeCaseInsensitiveResourcePath(t *testing.T) {
	t.Parallel()

	userpassEntry := &MountEntry{Table: credentialTableType, Type: "userpass", Path: "userpass/"}
	approleEntry := &MountEntry{Table: credentialTableType, Type: "approle", Path: "approle/"}
	awsEC2Entry := &MountEntry{Table: credentialTableType, Type: "aws-ec2", Path: "aws-ec2/"}
	azureEntry := &MountEntry{Table: credentialTableType, Type: "azure", Path: "azure/"}
	gcpEntry := &MountEntry{Table: credentialTableType, Type: "gcp", Path: "gcp/"}
	githubEntry := &MountEntry{Table: credentialTableType, Type: "github", Path: "github/"}
	tpmEntry := &MountEntry{Table: credentialTableType, Type: "tpm", Path: "tpm/"}
	kubernetesEntry := &MountEntry{Table: credentialTableType, Type: "kubernetes", Path: "kubernetes/"}
	oktaEntry := &MountEntry{Table: credentialTableType, Type: "okta", Path: "okta/"}
	ldapEntry := &MountEntry{Table: credentialTableType, Type: "ldap", Path: "ldap/"}
	radiusEntry := &MountEntry{Table: credentialTableType, Type: "radius", Path: "radius/"}
	// A secrets mount whose type collides with an auth backend name; plugin
	// names are namespaced by plugin type, so this is legitimately possible.
	secretCertEntry := &MountEntry{Table: mountTableType, Type: "cert", Path: "cert/"}
	// An external auth plugin overriding the builtin of the same name, which
	// Vault permits for unversioned plugins. A non-empty RunningSha256 is
	// what marks a mount as running an external plugin.
	externalApproleEntry := &MountEntry{
		Table:         credentialTableType,
		Type:          "approle",
		Path:          "approle-ext/",
		RunningSha256: "abc123",
	}

	tests := []struct {
		name                  string
		path                  string
		entry                 *MountEntry
		disableLowerCaseNames bool
		want                  string
	}{
		{
			name: "identity group name lowercased by default",
			path: "identity/group/name/Admin",
			want: "identity/group/name/admin",
		},
		{
			name: "identity entity name lowercased by default",
			path: "identity/entity/name/Foo",
			want: "identity/entity/name/foo",
		},
		{
			name:                  "identity group name untouched in case-sensitive mode",
			path:                  "identity/group/name/Admin",
			disableLowerCaseNames: true,
			want:                  "identity/group/name/Admin",
		},
		{
			name: "acl policy name lowercased and trimmed",
			path: "sys/policies/acl/ Admin ",
			want: "sys/policies/acl/admin",
		},
		{
			name: "legacy policy name lowercased and trimmed",
			path: "sys/policy/ Admin ",
			want: "sys/policy/admin",
		},
		{
			name: "rgp policy name lowercased and trimmed",
			path: "sys/policies/rgp/ Admin ",
			want: "sys/policies/rgp/admin",
		},
		{
			name: "egp policy name lowercased and trimmed",
			path: "sys/policies/egp/ Admin ",
			want: "sys/policies/egp/admin",
		},
		{
			name:  "userpass user name lowercased",
			path:  "auth/userpass/users/Bob",
			entry: userpassEntry,
			want:  "auth/userpass/users/bob",
		},
		{
			name:  "userpass login name lowercased",
			path:  "auth/userpass/login/Bob",
			entry: userpassEntry,
			want:  "auth/userpass/login/bob",
		},
		{
			name:  "approle role name lowercased, sub-resource preserved",
			path:  "auth/approle/role/MyRole/role-id",
			entry: approleEntry,
			want:  "auth/approle/role/myrole/role-id",
		},
		{
			name:  "aws-ec2 alias mount resolves to aws role prefix",
			path:  "auth/aws-ec2/role/Admin",
			entry: awsEC2Entry,
			want:  "auth/aws-ec2/role/admin",
		},
		{
			name:  "azure role name lowercased",
			path:  "auth/azure/role/Admin",
			entry: azureEntry,
			want:  "auth/azure/role/admin",
		},
		{
			name:  "gcp role name lowercased, sub-resource preserved",
			path:  "auth/gcp/role/Admin/service-accounts",
			entry: gcpEntry,
			want:  "auth/gcp/role/admin/service-accounts",
		},
		{
			name:  "kubernetes role name lowercased",
			path:  "auth/kubernetes/role/Admin",
			entry: kubernetesEntry,
			want:  "auth/kubernetes/role/admin",
		},
		{
			name:  "github team map key lowercased",
			path:  "auth/github/map/teams/Admins",
			entry: githubEntry,
			want:  "auth/github/map/teams/admins",
		},
		{
			name:  "github user map key lowercased",
			path:  "auth/github/map/users/Bob",
			entry: githubEntry,
			want:  "auth/github/map/users/bob",
		},
		{
			name:  "tpm role name lowercased",
			path:  "auth/tpm/role/Admin",
			entry: tpmEntry,
			want:  "auth/tpm/role/admin",
		},
		{
			name:  "okta group name deliberately untouched",
			path:  "auth/okta/groups/Admin",
			entry: oktaEntry,
			want:  "auth/okta/groups/Admin",
		},
		// ldap, radius, and scim are known-uncovered rather than unaffected;
		// these pin that exclusion so a change in either direction is
		// deliberate. See normalizeCaseInsensitiveResourcePath's doc comment.
		{
			name:  "ldap user name deliberately untouched",
			path:  "auth/ldap/users/Bob",
			entry: ldapEntry,
			want:  "auth/ldap/users/Bob",
		},
		{
			name:  "radius user name deliberately untouched",
			path:  "auth/radius/users/Bob",
			entry: radiusEntry,
			want:  "auth/radius/users/Bob",
		},
		{
			name: "scim client name deliberately untouched",
			path: "identity/scim/client/Admin",
			want: "identity/scim/client/Admin",
		},
		{
			name:  "secrets mount with colliding type not normalized",
			path:  "cert/certs/Admin",
			entry: secretCertEntry,
			want:  "cert/certs/Admin",
		},
		{
			name:  "external plugin overriding builtin auth type not normalized",
			path:  "auth/approle-ext/role/MyRole",
			entry: externalApproleEntry,
			want:  "auth/approle-ext/role/MyRole",
		},
		{
			name:  "unrelated path untouched",
			path:  "secret/data/Foo",
			entry: approleEntry,
			want:  "secret/data/Foo",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := &Core{identityStore: &IdentityStore{}}
			if tc.disableLowerCaseNames {
				c.identityStore.SetDisableLowerCasedNames()
			}
			require.Equal(t, tc.want, c.normalizeCaseInsensitiveResourcePath(tc.path, tc.entry))
		})
	}
}

// TestRestoreForwardingTokenHeaders_UsesInboundToken verifies Authorization
// forwarding prefers the original inbound token when present.
func TestRestoreForwardingTokenHeaders_UsesInboundToken(t *testing.T) {
	t.Parallel()

	req := &logical.Request{
		ClientToken:       "jwt.internal-id",
		InboundSSCToken:   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.payload.sig",
		ClientTokenSource: logical.ClientTokenFromAuthzHeader,
		Headers: map[string][]string{
			"Authorization": {"Basic abc123"},
		},
	}

	restoreForwardingTokenHeaders(req)

	require.Equal(t, []string{"Basic abc123", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.payload.sig"}, req.Headers["Authorization"])
}

// TestRestoreForwardingTokenHeaders_FallsBackToClientToken verifies fallback to
// req.ClientToken when no inbound token is present.
func TestRestoreForwardingTokenHeaders_FallsBackToClientToken(t *testing.T) {
	t.Parallel()

	req := &logical.Request{
		ClientToken:       "jwt.jti-value",
		ClientTokenSource: logical.ClientTokenFromVaultHeader,
	}

	restoreForwardingTokenHeaders(req)

	require.Equal(t, []string{"jwt.jti-value"}, req.Headers["X-Vault-Token"])
}

// TestRestoreForwardingTokenHeaders_UsesInboundTokenForVaultHeader verifies
// X-Vault-Token forwarding prefers the original inbound token.
func TestRestoreForwardingTokenHeaders_UsesInboundTokenForVaultHeader(t *testing.T) {
	t.Parallel()

	req := &logical.Request{
		ClientToken:       "jwt.jti-value",
		InboundSSCToken:   "jwt.raw.value",
		ClientTokenSource: logical.ClientTokenFromVaultHeader,
	}

	restoreForwardingTokenHeaders(req)

	require.Equal(t, []string{"jwt.raw.value"}, req.Headers["X-Vault-Token"])
}

func TestRequestHandling_Wrapping(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	core.logicalBackends["kv"] = PassthroughBackendFactory

	meUUID, _ := uuid.GenerateUUID()
	err := core.mount(namespace.RootContext(nil), &MountEntry{
		Table: mountTableType,
		UUID:  meUUID,
		Path:  "wraptest",
		Type:  "kv",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// No duration specified
	req := &logical.Request{
		Path:        "wraptest/foo",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"zip": "zap",
		},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = &logical.Request{
		Path:        "wraptest/foo",
		ClientToken: root,
		Operation:   logical.ReadOperation,
		WrapInfo: &logical.RequestWrapInfo{
			TTL: time.Duration(15 * time.Second),
		},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.WrapInfo == nil || resp.WrapInfo.TTL != time.Duration(15*time.Second) {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestRequestHandling_LoginWrapping(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	if err := core.loadMounts(namespace.RootContext(nil)); err != nil {
		t.Fatalf("err: %v", err)
	}

	core.credentialBackends["userpass"] = credUserpass.Factory

	// No duration specified
	req := &logical.Request{
		Path:        "sys/auth/userpass",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "userpass",
		},
		Connection: &logical.Connection{},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req.Path = "auth/userpass/users/test"
	req.Data = map[string]interface{}{
		"password": "foo",
		"policies": "default",
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = &logical.Request{
		Path:      "auth/userpass/login/test",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"password": "foo",
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.WrapInfo != nil {
		t.Fatalf("bad: %#v", resp)
	}

	req = &logical.Request{
		Path:      "auth/userpass/login/test",
		Operation: logical.UpdateOperation,
		WrapInfo: &logical.RequestWrapInfo{
			TTL: time.Duration(15 * time.Second),
		},
		Data: map[string]interface{}{
			"password": "foo",
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.WrapInfo == nil || resp.WrapInfo.TTL != time.Duration(15*time.Second) {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestRequestHandling_Login_PeriodicToken(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	if err := core.loadMounts(namespace.RootContext(nil)); err != nil {
		t.Fatalf("err: %v", err)
	}

	core.credentialBackends["approle"] = approle.Factory

	// Enable approle
	req := &logical.Request{
		Path:        "sys/auth/approle",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "approle",
		},
		Connection: &logical.Connection{},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Create role
	req.Path = "auth/approle/role/role-period"
	req.Data = map[string]interface{}{
		"period": "5s",
	}
	_, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Get role ID
	req.Path = "auth/approle/role/role-period/role-id"
	req.Operation = logical.ReadOperation
	req.Data = nil
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	roleID := resp.Data["role_id"]

	// Get secret ID
	req.Path = "auth/approle/role/role-period/secret-id"
	req.Operation = logical.UpdateOperation
	req.Data = nil
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	secretID := resp.Data["secret_id"]

	// Perform login
	req = &logical.Request{
		Path:      "auth/approle/login",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"role_id":   roleID,
			"secret_id": secretID,
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatalf("bad: %v", resp)
	}
	loginToken := resp.Auth.ClientToken
	entityID := resp.Auth.EntityID
	accessor := resp.Auth.Accessor

	// Perform token lookup on the generated token
	req = &logical.Request{
		Path:        "auth/token/lookup",
		Operation:   logical.UpdateOperation,
		ClientToken: root,
		Data: map[string]interface{}{
			"token": loginToken,
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}

	if resp.Data["creation_time"].(int64) == 0 {
		t.Fatal("creation time was zero")
	}

	// Depending on timing of the test this may have ticked down, so reset it
	// back to the original value as long as it's not expired.
	if resp.Data["ttl"].(int64) > 0 && resp.Data["ttl"].(int64) < 5 {
		resp.Data["ttl"] = int64(5)
	}

	exp := map[string]interface{}{
		"accessor":         accessor,
		"creation_time":    resp.Data["creation_time"].(int64),
		"creation_ttl":     int64(5),
		"display_name":     "approle",
		"entity_id":        entityID,
		"expire_time":      resp.Data["expire_time"].(time.Time),
		"explicit_max_ttl": int64(0),
		"id":               loginToken,
		"issue_time":       resp.Data["issue_time"].(time.Time),
		"meta":             map[string]string{"role_name": "role-period"},
		"num_uses":         0,
		"orphan":           true,
		"path":             "auth/approle/login",
		"period":           int64(5),
		"policies":         []string{"default"},
		"renewable":        true,
		"ttl":              int64(5),
		"type":             "service",
	}

	if diff := deep.Equal(resp.Data, exp); diff != nil {
		t.Fatal(diff)
	}
}

// TestLoginCreateToken_SkipsLegacyInvalidTokenPolicyNames verifies token
// creation skips malformed legacy token policy names instead of failing login.
//
// This test intentionally calls LoginCreateToken directly so it can inject a
// legacy malformed policy value into logical.Response.Auth. The public API
// validates policy names during normal writes and login flows, so it cannot
// create this legacy input state.
func TestLoginCreateToken_SkipsLegacyInvalidTokenPolicyNames(t *testing.T) {
	t.Parallel()

	cluster := NewTestCluster(t, &CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}, nil)
	core := cluster.Cores[0].Core
	ctx := namespace.RootContext(t.Context())

	resp, err := core.HandleRequest(ctx, &logical.Request{
		Path:        "sys/auth/userpass",
		ClientToken: cluster.RootToken,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "userpass",
		},
		Connection: &logical.Connection{},
	})
	require.NoError(t, err)
	require.Nil(t, resp)

	loginResp := &logical.Response{
		Auth: &logical.Auth{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
			DisplayName: "legacy-token-policy-user",
			Policies:    []string{"default", "team/../read"},
			TokenType:   logical.TokenTypeService,
		},
	}

	leaseGenerated, resp, err := core.LoginCreateToken(
		ctx,
		namespace.RootNamespace,
		"auth/userpass/login/legacy-token-policy-user",
		"auth/userpass/",
		"",
		loginResp,
	)
	require.NoError(t, err)
	require.True(t, leaseGenerated)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Auth)
	require.ElementsMatch(t, []string{"default"}, resp.Auth.TokenPolicies)
	require.ElementsMatch(t, []string{"default"}, resp.Auth.Policies)
}

// TestLoginCreateToken_SkipsLegacyInvalidIdentityPolicyNames verifies token
// creation skips malformed legacy identity policy names instead of failing login.
//
// This test must inject legacy malformed identity policy values directly into
// identity storage before LoginCreateToken runs. The public API rejects these
// malformed policy names, so it cannot construct this legacy policy state.
func TestLoginCreateToken_SkipsLegacyInvalidIdentityPolicyNames(t *testing.T) {
	t.Parallel()

	cluster := NewTestCluster(t, &CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}, nil)
	core := cluster.Cores[0].Core
	ctx := namespace.RootContext(t.Context())

	resp, err := core.HandleRequest(ctx, &logical.Request{
		Path:        "sys/auth/userpass",
		ClientToken: cluster.RootToken,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "userpass",
		},
		Connection: &logical.Connection{},
	})
	require.NoError(t, err)
	require.Nil(t, resp)

	entityID, err := uuid.GenerateUUID()
	require.NoError(t, err)
	require.NoError(t, core.identityStore.upsertEntity(ctx, &identity.Entity{
		ID:          entityID,
		Name:        "legacy-identity-policy-user",
		NamespaceID: namespace.RootNamespaceID,
		BucketKey:   core.identityStore.entityPacker.BucketKey(entityID),
		Policies:    []string{" Team-Read ", "team/../read"},
	}, nil, true))

	loginResp := &logical.Response{
		Auth: &logical.Auth{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
			DisplayName: "legacy-identity-policy-user",
			EntityID:    entityID,
			Policies:    []string{"default"},
			TokenType:   logical.TokenTypeService,
		},
	}

	leaseGenerated, resp, err := core.LoginCreateToken(
		ctx,
		namespace.RootNamespace,
		"auth/userpass/login/legacy-identity-policy-user",
		"auth/userpass/",
		"",
		loginResp,
	)
	require.NoError(t, err)
	require.True(t, leaseGenerated)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Auth)
	require.ElementsMatch(t, []string{"team-read"}, resp.Auth.IdentityPolicies)
	require.ElementsMatch(t, []string{"default", "team-read"}, resp.Auth.Policies)
}

func labelsMatch(actual, expected map[string]string) bool {
	for expected_label, expected_val := range expected {
		if v, ok := actual[expected_label]; ok {
			if v != expected_val {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func checkCounter(t *testing.T, inmemSink *metrics.InmemSink, keyPrefix string, expectedLabels map[string]string) {
	t.Helper()

	intervals := inmemSink.Data()
	if len(intervals) > 1 {
		t.Skip("Detected interval crossing.")
	}

	var counter *metrics.SampledValue = nil
	var labels map[string]string
	for _, c := range intervals[0].Counters {
		if !strings.HasPrefix(c.Name, keyPrefix) {
			continue
		}
		counter = &c

		labels = make(map[string]string)
		for _, l := range counter.Labels {
			labels[l.Name] = l.Value
		}

		// Distinguish between different label sets
		if labelsMatch(labels, expectedLabels) {
			break
		}
	}
	if counter == nil {
		t.Fatalf("No %q counter found with matching labels", keyPrefix)
	}

	if !labelsMatch(labels, expectedLabels) {
		t.Errorf("No matching label set, found %v", labels)
	}

	if counter.Count != 1 {
		t.Errorf("Counter number of samples %v is not 1.", counter.Count)
	}

	if counter.Sum != 1.0 {
		t.Errorf("Counter sum %v is not 1.", counter.Sum)
	}
}

func TestRequestHandling_LoginMetric(t *testing.T) {
	core, _, root, sink := TestCoreUnsealedWithMetrics(t)

	if err := core.loadMounts(namespace.RootContext(nil)); err != nil {
		t.Fatalf("err: %v", err)
	}

	core.credentialBackends["userpass"] = credUserpass.Factory

	// Setup mount
	req := &logical.Request{
		Path:        "sys/auth/userpass",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"type": "userpass",
		},
		Connection: &logical.Connection{},
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Create user
	req.Path = "auth/userpass/users/test"
	req.Data = map[string]interface{}{
		"password": "foo",
		"policies": "default",
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Login with response wrapping
	req = &logical.Request{
		Path:      "auth/userpass/login/test",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"password": "foo",
		},
		WrapInfo: &logical.RequestWrapInfo{
			TTL: time.Duration(15 * time.Second),
		},
		Connection: &logical.Connection{},
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}

	// There should be two counters
	checkCounter(t, sink, "token.creation",
		map[string]string{
			"cluster":      "test-cluster",
			"namespace":    "root",
			"auth_method":  "userpass",
			"mount_point":  "auth/userpass/",
			"creation_ttl": "+Inf",
			"token_type":   "service",
		},
	)
	checkCounter(t, sink, "token.creation",
		map[string]string{
			"cluster":      "test-cluster",
			"namespace":    "root",
			"auth_method":  "response_wrapping",
			"mount_point":  "auth/userpass/",
			"creation_ttl": "1m",
			"token_type":   "service",
		},
	)
}

func TestRequestHandling_SecretLeaseMetric(t *testing.T) {
	coreConfig := &CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": LeasedPassthroughBackendFactory,
		},
	}
	core, _, root, sink := TestCoreUnsealedWithMetricsAndConfig(t, coreConfig)

	// Create a key with a lease
	req := logical.TestRequest(t, logical.UpdateOperation, "secret/foo")
	req.Data["foo"] = "bar"
	req.ClientToken = root
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read a key with a LeaseID
	req = logical.TestRequest(t, logical.ReadOperation, "secret/foo")
	req.ClientToken = root
	err = core.PopulateTokenEntry(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp)
	}

	checkCounter(t, sink, "secret.lease.creation",
		map[string]string{
			"cluster":       "test-cluster",
			"namespace":     "root",
			"secret_engine": "kv",
			"mount_point":   "secret/",
			"creation_ttl":  "+Inf",
		},
	)
}

// TestRequestHandling_isRetryableRPCError tests that a retryable RPC error
// can be distinguished from a normal error
func TestRequestHandling_isRetryableRPCError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	deadlineCtx, deadlineCancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
	defer deadlineCancel()
	testCases := []struct {
		name string
		ctx  context.Context
		err  error
		want bool
	}{
		{
			name: "req context canceled, not deadline",
			ctx:  ctx,
			err:  status.Error(codes.Canceled, "context canceled"),
			want: true,
		},
		{
			name: "req context deadline exceeded",
			ctx:  deadlineCtx,
			err:  status.Error(codes.Canceled, "context canceled"),
			want: false,
		},
		{
			name: "server context canceled",
			err:  status.Error(codes.Canceled, "context canceled"),
			want: true,
		},
		{
			name: "unavailable",
			err:  status.Error(codes.Unavailable, "unavailable"),
			want: true,
		},
		{
			name: "other status",
			err:  status.Error(codes.FailedPrecondition, "failed"),
			want: false,
		},
		{
			name: "other unknown",
			err:  status.Error(codes.Unknown, "unknown"),
			want: false,
		},
		{
			name: "malformed header unknown",
			err:  status.Error(codes.Unknown, "malformed header: missing HTTP content-type"),
			want: true,
		},
		{
			name: "other error",
			err:  errors.New("other type of error"),
			want: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			useCtx := tc.ctx
			if tc.ctx == nil {
				useCtx = context.Background()
			}
			require.Equal(t, tc.want, isRetryableRPCError(useCtx, tc.err))
		})
	}
}

// TestRequestHandling_TokenRenewal tests that a renewable token can be renewed
// and that an error is returned when lease_id is not a string
func TestRequestHandling_TokenRenewal(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	// First, create a renewable token with a short TTL
	req := &logical.Request{
		Path:        "auth/token/create",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"ttl":       "1h",
			"renewable": true,
			"policies":  []string{"default"},
		},
	}

	resp, err := core.HandleRequest(namespace.RootContext(context.TODO()), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatalf("bad: %v", resp)
	}

	newToken := resp.Auth.ClientToken
	if newToken == "" {
		t.Fatal("expected non-empty token")
	}
	if !resp.Auth.Renewable {
		t.Fatal("expected renewable token")
	}

	// Test token renewal
	req = &logical.Request{
		Path:        "auth/token/renew-self",
		ClientToken: newToken,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"increment": "2h", // Extend by 2 hours
		},
	}

	resp, err = core.HandleRequest(namespace.RootContext(context.TODO()), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Auth == nil {
		t.Fatalf("bad: %v", resp)
	}

	// Verify the token was renewed
	if resp.Auth.ClientToken != newToken {
		t.Fatalf("expected same token, got %s", resp.Auth.ClientToken)
	}
	if !resp.Auth.Renewable {
		t.Fatal("expected renewable token after renewal")
	}

	req = &logical.Request{
		Path:        "sys/leases/renew",
		ClientToken: root,
		Operation:   logical.UpdateOperation,
		Data: map[string]interface{}{
			"lease_id": 12345, // Non-string value
		},
	}

	resp, err = core.HandleRequest(namespace.RootContext(context.TODO()), req)
	if err == nil {
		t.Fatal("expected error when lease_id is not a string")
	}
	if !strings.Contains(err.Error(), "invalid request") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_NilRequest tests that
// fetchACLTokenEntryAndEntity returns an error when called with a nil request
func TestRequestHandling_fetchACLTokenEntryAndEntity_NilRequest(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	// Call with nil request - should return ErrInternalError
	_, _, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, nil)

	require.Error(t, err)
	require.Equal(t, ErrInternalError, err)
}

// TestAuth_AuthorizationDetails_CopiedFromRequest verifies that logical.Auth.AuthorizationDetails
// matches the authorization details already carried on the request.
func TestAuth_AuthorizationDetails_CopiedFromRequest(t *testing.T) {
	t.Parallel()

	details := []logical.AuthorizationDetail{
		{"type": "account_information", "scope": "read"},
		{"type": "payment_initiation", "amount": "100"},
	}

	auth := &logical.Auth{}
	req := &logical.Request{
		JwtAuthorizationDetails: details,
	}

	// Simulate the assignment performed in CheckToken.
	auth.AuthorizationDetails = req.JwtAuthorizationDetails

	require.Equal(t, details, auth.AuthorizationDetails, "auth.AuthorizationDetails must equal req.JwtAuthorizationDetails")
}

// TestAuth_AuthorizationDetails_NilWhenAbsent verifies that auth.AuthorizationDetails is nil
// when the request does not carry authorization details.
func TestAuth_AuthorizationDetails_NilWhenAbsent(t *testing.T) {
	t.Parallel()

	auth := &logical.Auth{}
	req := &logical.Request{}

	auth.AuthorizationDetails = req.JwtAuthorizationDetails

	require.Nil(t, auth.AuthorizationDetails)
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_EmptyToken verifies that a
// request with an empty ClientToken is rejected with ErrPermissionDenied.
func TestRequestHandling_fetchACLTokenEntryAndEntity_EmptyToken(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	req := &logical.Request{ClientToken: ""}
	_, _, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, req)

	require.Error(t, err)
	require.Equal(t, logical.ErrPermissionDenied, err)
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_UnknownToken verifies that a
// non-existent token returns ErrPermissionDenied combined with ErrInvalidToken.
func TestRequestHandling_fetchACLTokenEntryAndEntity_UnknownToken(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	req := &logical.Request{ClientToken: "hvs.nonexistent-token-id"}
	_, _, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, req)

	require.Error(t, err)
	require.ErrorIs(t, err, logical.ErrPermissionDenied)
	require.ErrorIs(t, err, logical.ErrInvalidToken)
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_ValidRootToken verifies the
// happy path: a valid root token returns an ACL, the token entry, and no error.
func TestRequestHandling_fetchACLTokenEntryAndEntity_ValidRootToken(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	req := &logical.Request{ClientToken: root}
	acl, te, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, acl)
	require.NotNil(t, te)
	require.Equal(t, root, te.ID)
	require.Contains(t, te.Policies, "root")
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_CachedTokenEntry verifies
// that when a token entry is already cached on the request, the function uses
// the cached entry instead of performing a lookup.
func TestRequestHandling_fetchACLTokenEntryAndEntity_CachedTokenEntry(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	// Look up the root token to get a valid entry
	te, err := core.tokenStore.Lookup(ctx, root)
	require.NoError(t, err)
	require.NotNil(t, te)

	// Pre-cache the entry on the request
	req := &logical.Request{ClientToken: root}
	req.SetTokenEntry(te)

	acl, returnedTE, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, acl)
	// The returned token entry should be the same object we cached
	require.Same(t, te, returnedTE)
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_BoundCIDR_NoConnection
// verifies that a token with BoundCIDRs is rejected when the request has no
// connection information. The token entry is pre-cached on the request to
// isolate the CIDR check logic from token storage concerns.
func TestRequestHandling_fetchACLTokenEntryAndEntity_BoundCIDR_NoConnection(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	boundCIDRs, err := parseutil.ParseAddrs([]string{"10.0.0.0/8"})
	require.NoError(t, err)

	// Look up the root token to get a valid base entry, then add BoundCIDRs
	te, err := core.tokenStore.Lookup(ctx, root)
	require.NoError(t, err)
	te.TTL = time.Hour
	te.BoundCIDRs = boundCIDRs

	req := &logical.Request{
		ClientToken: root,
		// No Connection field set
	}
	req.SetTokenEntry(te)

	_, _, _, _, err = core.fetchACLTokenEntryAndEntity(ctx, req)

	require.Error(t, err)
	require.ErrorIs(t, err, logical.ErrPermissionDenied)
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_BoundCIDR_OutOfRange
// verifies that a token with BoundCIDRs is rejected when the request comes
// from an IP outside the allowed range. The token entry is pre-cached on the
// request to isolate the CIDR check logic from token storage concerns.
func TestRequestHandling_fetchACLTokenEntryAndEntity_BoundCIDR_OutOfRange(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	boundCIDRs, err := parseutil.ParseAddrs([]string{"10.0.0.0/8"})
	require.NoError(t, err)

	te, err := core.tokenStore.Lookup(ctx, root)
	require.NoError(t, err)
	te.TTL = time.Hour
	te.BoundCIDRs = boundCIDRs

	req := &logical.Request{
		ClientToken: root,
		Connection:  &logical.Connection{RemoteAddr: "192.168.1.1"},
	}
	req.SetTokenEntry(te)

	_, _, _, _, err = core.fetchACLTokenEntryAndEntity(ctx, req)

	require.Error(t, err)
	require.ErrorIs(t, err, logical.ErrPermissionDenied)
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_BoundCIDR_InRange verifies
// that a token with BoundCIDRs succeeds when the request comes from an IP
// within the allowed CIDR range. The token entry is pre-cached on the request
// to isolate the CIDR check logic from token storage concerns.
func TestRequestHandling_fetchACLTokenEntryAndEntity_BoundCIDR_InRange(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	boundCIDRs, err := parseutil.ParseAddrs([]string{"10.0.0.0/8"})
	require.NoError(t, err)

	te, err := core.tokenStore.Lookup(ctx, root)
	require.NoError(t, err)
	te.TTL = time.Hour
	te.BoundCIDRs = boundCIDRs

	req := &logical.Request{
		ClientToken: root,
		Connection:  &logical.Connection{RemoteAddr: "10.1.2.3"},
	}
	req.SetTokenEntry(te)

	acl, returnedTE, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, acl)
	require.NotNil(t, returnedTE)
	require.Equal(t, root, returnedTE.ID)
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_NonExpiring_RootIgnoresCIDR
// verifies that a non-expiring root token (TTL == 0) bypasses CIDR checks even
// if BoundCIDRs is set, since the CIDR check is gated on TTL != 0.
func TestRequestHandling_fetchACLTokenEntryAndEntity_NonExpiring_RootIgnoresCIDR(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	// Root token has TTL == 0, so CIDR checks should not apply
	req := &logical.Request{
		ClientToken: root,
		// No Connection, which would fail if CIDR checks ran
	}
	acl, te, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, acl)
	require.NotNil(t, te)
	require.Equal(t, time.Duration(0), te.TTL)
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_LegacyInvalidTokenPoliciesNotApplied
// verifies invalid legacy token policy names are ignored at request time even
// when they exist in storage, policyTypeMap, and cache.
func TestRequestHandling_fetchACLTokenEntryAndEntity_LegacyInvalidTokenPoliciesNotApplied(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		invalidPolicyName string
		primeCache        bool
	}{
		{
			name:              "cached sibling target",
			invalidPolicyName: "sibling/./admin",
			primeCache:        true,
		},
		{
			name:              "uncached root target",
			invalidPolicyName: "root/./admin",
			primeCache:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cluster := NewTestCluster(t, nil, nil)
			core := cluster.Cores[0].Core
			ctx := namespace.RootContext(t.Context())

			validPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/allowed" {
	capabilities = ["read"]
}
`, WithDenySlashInTemplatedPaths(core.denySlashInTemplatedPolicyPaths))
			require.NoError(t, err)
			validPolicy.Name = "legacy-token-valid"
			validPolicy.namespace = namespace.RootNamespace
			require.NoError(t, core.policyStore.SetPolicy(ctx, validPolicy))

			legacyInvalidPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/forbidden" {
	capabilities = ["read"]
}
`, WithDenySlashInTemplatedPaths(core.denySlashInTemplatedPolicyPaths))
			require.NoError(t, err)
			legacyInvalidPolicy.Name = tc.invalidPolicyName
			legacyInvalidPolicy.namespace = namespace.RootNamespace
			require.NoError(t, core.policyStore.setPolicyInternal(ctx, legacyInvalidPolicy, nil))

			invalidKey := policyCacheKey(namespace.RootNamespaceID, core.policyStore.sanitizeName(tc.invalidPolicyName))
			_, found := core.policyStore.policyTypeMap.Load(invalidKey)
			require.True(t, found)
			storedEntry, err := core.policyStore.getACLView(namespace.RootNamespace).Get(ctx, tc.invalidPolicyName)
			require.NoError(t, err)
			require.NotNil(t, storedEntry)
			if !tc.primeCache {
				core.policyStore.tokenPoliciesLRU.Remove(invalidKey)
			}

			tokenID, err := uuid.GenerateUUID()
			require.NoError(t, err)
			legacyToken := &logical.TokenEntry{
				ID:       tokenID,
				Path:     "auth/token/create",
				Policies: []string{"legacy-token-valid"},
				TTL:      time.Hour,
			}
			testMakeTokenDirectly(t, core.tokenStore, legacyToken)

			legacyToken.Policies = []string{"legacy-token-valid", tc.invalidPolicyName}
			require.NoError(t, core.tokenStore.store(ctx, legacyToken))

			acl, returnedTE, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, &logical.Request{
				ClientToken: tokenID,
			})
			require.NoError(t, err)
			require.NotNil(t, acl)
			require.NotNil(t, returnedTE)

			allowed := acl.AllowOperation(ctx, &logical.Request{
				Path:      "secret/data/allowed",
				Operation: logical.ReadOperation,
			}, false)
			require.True(t, allowed.Allowed)

			forbidden := acl.AllowOperation(ctx, &logical.Request{
				Path:      "secret/data/forbidden",
				Operation: logical.ReadOperation,
			}, false)
			require.False(t, forbidden.Allowed)
		})
	}
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_BackslashAndSlashPolicyNamesRemainDistinct
// verifies request-time ACL construction treats slash and backslash policy names
// as distinct, independently resolvable entries.
func TestRequestHandling_fetchACLTokenEntryAndEntity_BackslashAndSlashPolicyNamesRemainDistinct(t *testing.T) {
	t.Parallel()

	cluster := NewTestCluster(t, nil, nil)
	core := cluster.Cores[0].Core
	ctx := namespace.RootContext(t.Context())

	slashPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/slash" {
	capabilities = ["read"]
}
`, WithDenySlashInTemplatedPaths(core.denySlashInTemplatedPolicyPaths))
	require.NoError(t, err)
	slashPolicy.Name = "team/read"
	slashPolicy.namespace = namespace.RootNamespace
	require.NoError(t, core.policyStore.SetPolicy(ctx, slashPolicy))

	backslashPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/backslash" {
	capabilities = ["read"]
}
`, WithDenySlashInTemplatedPaths(core.denySlashInTemplatedPolicyPaths))
	require.NoError(t, err)
	backslashPolicy.Name = `team\read`
	backslashPolicy.namespace = namespace.RootNamespace
	require.NoError(t, core.policyStore.SetPolicy(ctx, backslashPolicy))

	tokenBackslashID, err := uuid.GenerateUUID()
	require.NoError(t, err)
	testMakeTokenDirectly(t, core.tokenStore, &logical.TokenEntry{
		ID:       tokenBackslashID,
		Path:     "auth/token/create",
		Policies: []string{`team\read`},
		TTL:      time.Hour,
	})

	backslashACL, returnedBackslashTE, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, &logical.Request{
		ClientToken: tokenBackslashID,
	})
	require.NoError(t, err)
	require.NotNil(t, backslashACL)
	require.NotNil(t, returnedBackslashTE)

	backslashAllowed := backslashACL.AllowOperation(ctx, &logical.Request{
		Path:      "secret/data/backslash",
		Operation: logical.ReadOperation,
	}, false)
	require.True(t, backslashAllowed.Allowed)

	slashDeniedFromBackslashToken := backslashACL.AllowOperation(ctx, &logical.Request{
		Path:      "secret/data/slash",
		Operation: logical.ReadOperation,
	}, false)
	require.False(t, slashDeniedFromBackslashToken.Allowed)

	tokenSlashID, err := uuid.GenerateUUID()
	require.NoError(t, err)
	testMakeTokenDirectly(t, core.tokenStore, &logical.TokenEntry{
		ID:       tokenSlashID,
		Path:     "auth/token/create",
		Policies: []string{"team/read"},
		TTL:      time.Hour,
	})

	slashACL, returnedSlashTE, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, &logical.Request{
		ClientToken: tokenSlashID,
	})
	require.NoError(t, err)
	require.NotNil(t, slashACL)
	require.NotNil(t, returnedSlashTE)

	slashAllowed := slashACL.AllowOperation(ctx, &logical.Request{
		Path:      "secret/data/slash",
		Operation: logical.ReadOperation,
	}, false)
	require.True(t, slashAllowed.Allowed)

	backslashDeniedFromSlashToken := slashACL.AllowOperation(ctx, &logical.Request{
		Path:      "secret/data/backslash",
		Operation: logical.ReadOperation,
	}, false)
	require.False(t, backslashDeniedFromSlashToken.Allowed)
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_LegacyParentTraversalTokenPoliciesNotApplied
// verifies parent-traversal token policy names are ignored during ACL
// construction, even when cache and type-map entries exist.
func TestRequestHandling_fetchACLTokenEntryAndEntity_LegacyParentTraversalTokenPoliciesNotApplied(t *testing.T) {
	t.Parallel()

	cluster := NewTestCluster(t, nil, nil)
	core := cluster.Cores[0].Core
	ctx := namespace.RootContext(t.Context())

	validPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/allowed" {
	capabilities = ["read"]
}
`, WithDenySlashInTemplatedPaths(core.denySlashInTemplatedPolicyPaths))
	require.NoError(t, err)
	validPolicy.Name = "legacy-token-valid"
	validPolicy.namespace = namespace.RootNamespace
	require.NoError(t, core.policyStore.SetPolicy(ctx, validPolicy))

	const invalidPolicyName = "team/../admin"
	invalidKey := policyCacheKey(namespace.RootNamespaceID, core.policyStore.sanitizeName(invalidPolicyName))
	core.policyStore.policyTypeMap.Store(invalidKey, PolicyTypeACL)
	require.NotNil(t, core.policyStore.tokenPoliciesLRU)
	core.policyStore.tokenPoliciesLRU.Add(invalidKey, &Policy{
		Name:      invalidPolicyName,
		Type:      PolicyTypeACL,
		Raw:       `path "secret/data/forbidden" { capabilities = ["read"] }`,
		namespace: namespace.RootNamespace,
	})

	tokenID, err := uuid.GenerateUUID()
	require.NoError(t, err)
	legacyToken := &logical.TokenEntry{
		ID:       tokenID,
		Path:     "auth/token/create",
		Policies: []string{"legacy-token-valid"},
		TTL:      time.Hour,
	}
	testMakeTokenDirectly(t, core.tokenStore, legacyToken)

	legacyToken.Policies = []string{"legacy-token-valid", invalidPolicyName}
	require.NoError(t, core.tokenStore.store(ctx, legacyToken))

	acl, returnedTE, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, &logical.Request{
		ClientToken: tokenID,
	})
	require.NoError(t, err)
	require.NotNil(t, acl)
	require.NotNil(t, returnedTE)

	allowed := acl.AllowOperation(ctx, &logical.Request{
		Path:      "secret/data/allowed",
		Operation: logical.ReadOperation,
	}, false)
	require.True(t, allowed.Allowed)

	forbidden := acl.AllowOperation(ctx, &logical.Request{
		Path:      "secret/data/forbidden",
		Operation: logical.ReadOperation,
	}, false)
	require.False(t, forbidden.Allowed)
}

// TestRequestHandling_fetchACLTokenEntryAndEntity_LegacyInvalidIdentityPoliciesNotApplied
// verifies invalid legacy identity policy names are ignored at request time
// even when they exist in storage, policyTypeMap, and cache.
//
// This test intentionally uses core.HandleRequest for auth enable/user creation
// and login setup. In this in-package harness, using client.Logical().Write for
// these calls can fail JSON decoding into api.Secret, while HandleRequest
// exercises the same logical paths without HTTP client decoding behavior.
func TestRequestHandling_fetchACLTokenEntryAndEntity_LegacyInvalidIdentityPoliciesNotApplied(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		invalidPolicyName string
		primeCache        bool
	}{
		{
			name:              "cached sibling target",
			invalidPolicyName: "sibling/./admin",
			primeCache:        true,
		},
		{
			name:              "uncached root target",
			invalidPolicyName: "root/./admin",
			primeCache:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cluster := NewTestCluster(t, &CoreConfig{
				CredentialBackends: map[string]logical.Factory{
					"userpass": credUserpass.Factory,
				},
			}, nil)
			core := cluster.Cores[0].Core
			ctx := namespace.RootContext(t.Context())

			req := &logical.Request{
				Path:        "sys/auth/userpass",
				ClientToken: cluster.RootToken,
				Operation:   logical.UpdateOperation,
				Data: map[string]interface{}{
					"type": "userpass",
				},
				Connection: &logical.Connection{},
			}
			resp, err := core.HandleRequest(ctx, req)
			require.NoError(t, err)
			require.Nil(t, resp)

			req.Path = "auth/userpass/users/legacy-identity-user"
			req.Data = map[string]interface{}{
				"password": "testpassword",
			}
			resp, err = core.HandleRequest(ctx, req)
			require.NoError(t, err)
			require.Nil(t, resp)

			validPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/allowed" {
	capabilities = ["read"]
}
`, WithDenySlashInTemplatedPaths(core.denySlashInTemplatedPolicyPaths))
			require.NoError(t, err)
			validPolicy.Name = "legacy-identity-valid"
			validPolicy.namespace = namespace.RootNamespace
			require.NoError(t, core.policyStore.SetPolicy(ctx, validPolicy))

			legacyInvalidPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/forbidden" {
	capabilities = ["read"]
}
`, WithDenySlashInTemplatedPaths(core.denySlashInTemplatedPolicyPaths))
			require.NoError(t, err)
			legacyInvalidPolicy.Name = tc.invalidPolicyName
			legacyInvalidPolicy.namespace = namespace.RootNamespace
			require.NoError(t, core.policyStore.setPolicyInternal(ctx, legacyInvalidPolicy, nil))

			invalidKey := policyCacheKey(namespace.RootNamespaceID, core.policyStore.sanitizeName(tc.invalidPolicyName))
			_, found := core.policyStore.policyTypeMap.Load(invalidKey)
			require.True(t, found)
			if !tc.primeCache {
				core.policyStore.tokenPoliciesLRU.Remove(invalidKey)
			}

			loginResp, err := core.HandleRequest(ctx, &logical.Request{
				Path:      "auth/userpass/login/legacy-identity-user",
				Operation: logical.UpdateOperation,
				Data: map[string]interface{}{
					"password": "testpassword",
				},
				Connection: &logical.Connection{},
			})
			require.NoError(t, err)
			require.NotNil(t, loginResp)
			require.NotNil(t, loginResp.Auth)
			tokenEntry, err := core.tokenStore.Lookup(ctx, loginResp.Auth.ClientToken)
			require.NoError(t, err)
			require.NotNil(t, tokenEntry)
			require.NotEmpty(t, tokenEntry.EntityID)

			var getErr error
			core.identityStore.lock.Lock()
			entity, getErr := core.identityStore.MemDBEntityByID(tokenEntry.EntityID, true)
			if getErr == nil {
				if entity == nil {
					getErr = errors.New("entity not found")
				}
			}
			if getErr == nil {
				entity.Policies = []string{"legacy-identity-valid", tc.invalidPolicyName}
				getErr = core.identityStore.upsertEntity(ctx, entity, nil, true)
			}
			core.identityStore.lock.Unlock()
			require.NoError(t, getErr)

			acl, returnedTE, _, _, err := core.fetchACLTokenEntryAndEntity(ctx, &logical.Request{
				ClientToken: loginResp.Auth.ClientToken,
			})
			require.NoError(t, err)
			require.NotNil(t, acl)
			require.NotNil(t, returnedTE)

			allowed := acl.AllowOperation(ctx, &logical.Request{
				Path:      "secret/data/allowed",
				Operation: logical.ReadOperation,
			}, false)
			require.True(t, allowed.Allowed)

			forbidden := acl.AllowOperation(ctx, &logical.Request{
				Path:      "secret/data/forbidden",
				Operation: logical.ReadOperation,
			}, false)
			require.False(t, forbidden.Allowed)
		})
	}
}

// TestRequestHandling_handleCancelableTestNumericToken tests that if a token
// that is passed in, is somehow a number rather than a string (not currently
// possible), then the handling will error, not panic.
func TestRequestHandling_handleCancelableTestNumericToken(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	data := map[string]interface{}{"token": 5}
	req := &logical.Request{Data: data, Path: "auth/token/lookup"}

	resp, err := core.handleCancelableRequest(ctx, req)
	require.True(t, resp != nil && err != nil)
	require.ErrorContains(t, err, logical.ErrPermissionDenied.Error())
	require.ErrorContains(t, resp.Error(), "invalid token")
}
