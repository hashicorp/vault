// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	"github.com/go-test/deep"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/credential/approle"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/namespace"
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
		{name: "leases lookup", path: "sys/leases/lookup", want: true},
		{name: "leases lookup prefix", path: "sys/leases/lookup/secret/foo", want: true},
		{name: "leases count", path: "sys/leases/count", want: true},
		{name: "leases list", path: "sys/leases", want: true},
		{name: "cubbyhole", path: "cubbyhole/test", want: true},
		{name: "token renew self excluded", path: "auth/token/renew-self", want: false},
		{name: "leases renew excluded", path: "sys/leases/renew", want: false},
		{name: "unrelated", path: "secret/data/foo", want: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, requiresMaterializedTokenState(tc.path))
		})
	}
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

// Test_allPoliciesAllowOnly tests a helper function that checks if all policies in
// a given set have only "allow" capabilities, and not "deny" or "sudo"
func Test_allPoliciesAllowOnly(t *testing.T) {
	t.Parallel()

	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(context.Background())

	allowPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/*" {
	capabilities = ["read", "list"]
}
`)
	require.NoError(t, err)
	allowPolicy.Name = "allow-only"
	require.NoError(t, c.policyStore.SetPolicy(ctx, allowPolicy))

	denyPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/*" {
	capabilities = ["deny"]
}
`)
	require.NoError(t, err)
	denyPolicy.Name = "deny-policy"
	require.NoError(t, c.policyStore.SetPolicy(ctx, denyPolicy))

	sudoPolicy, err := ParseACLPolicy(namespace.RootNamespace, `
path "secret/data/*" {
	capabilities = ["read", "sudo"]
}
`)
	require.NoError(t, err)
	sudoPolicy.Name = "sudo-policy"
	require.NoError(t, c.policyStore.SetPolicy(ctx, sudoPolicy))

	tests := map[string]struct {
		policyNamesByNamespace map[string][]string
		expected               bool
		wantErr                string
	}{
		"all allow only": {
			policyNamesByNamespace: map[string][]string{
				namespace.RootNamespaceID: {"allow-only"},
			},
			expected: true,
		},
		"deny policy": {
			policyNamesByNamespace: map[string][]string{
				namespace.RootNamespaceID: {"deny-policy"},
			},
			expected: false,
		},
		"sudo policy": {
			policyNamesByNamespace: map[string][]string{
				namespace.RootNamespaceID: {"sudo-policy"},
			},
			expected: false,
		},
		"root policy": {
			policyNamesByNamespace: map[string][]string{
				namespace.RootNamespaceID: {"root"},
			},
			expected: false,
		},
		"missing policy": {
			policyNamesByNamespace: map[string][]string{
				namespace.RootNamespaceID: {"missing-policy"},
			},
			expected: false,
			wantErr:  "policy \"missing-policy\" not found",
		},
		"missing namespace": {
			policyNamesByNamespace: map[string][]string{
				"missing-namespace": {"allow-only"},
			},
			expected: false,
			wantErr:  namespace.ErrNoNamespace.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := c.allPoliciesAllowOnly(ctx, tc.policyNamesByNamespace)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expected, actual)
		})
	}
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
		EnterpriseTokenAuthorizationDetails: details,
	}

	// Simulate the assignment performed in CheckToken.
	auth.AuthorizationDetails = req.EnterpriseTokenAuthorizationDetails

	require.Equal(t, details, auth.AuthorizationDetails, "auth.AuthorizationDetails must equal req.EnterpriseTokenAuthorizationDetails")
}

// TestAuth_AuthorizationDetails_NilWhenAbsent verifies that auth.AuthorizationDetails is nil
// when the request does not carry authorization details.
func TestAuth_AuthorizationDetails_NilWhenAbsent(t *testing.T) {
	t.Parallel()

	auth := &logical.Auth{}
	req := &logical.Request{}

	auth.AuthorizationDetails = req.EnterpriseTokenAuthorizationDetails

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
