// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestACME_ValidateIdentifiersAgainstRole Verify the ACME order creation
// function verifies somewhat the identifiers that were provided have a
// decent chance of being allowed by the selected role.
func TestACME_ValidateIdentifiersAgainstRole(t *testing.T) {
	b, _ := CreateBackendWithStorage(t)

	tests := []struct {
		name        string
		role        *issuing.RoleEntry
		identifiers []*ACMEIdentifier
		expectErr   bool
	}{
		{
			name:        "verbatim-role-allows-dns-ip",
			role:        issuing.SignVerbatimRole(),
			identifiers: _buildACMEIdentifiers("test.com", "127.0.0.1"),
			expectErr:   false,
		},
		{
			name:        "default-role-does-not-allow-dns",
			role:        buildTestRole(t, nil),
			identifiers: _buildACMEIdentifiers("www.test.com"),
			expectErr:   true,
		},
		{
			name:        "default-role-allows-ip",
			role:        buildTestRole(t, nil),
			identifiers: _buildACMEIdentifiers("192.168.0.1"),
			expectErr:   false,
		},
		{
			name:        "disable-ip-sans-forbids-ip",
			role:        buildTestRole(t, map[string]interface{}{"allow_ip_sans": false}),
			identifiers: _buildACMEIdentifiers("192.168.0.1"),
			expectErr:   true,
		},
		{
			name: "role-no-wildcards-allowed-without",
			role: buildTestRole(t, map[string]interface{}{
				"allow_subdomains":            true,
				"allow_bare_domains":          true,
				"allowed_domains":             []string{"test.com"},
				"allow_wildcard_certificates": false,
			}),
			identifiers: _buildACMEIdentifiers("www.test.com", "test.com"),
			expectErr:   false,
		},
		{
			name: "role-no-wildcards-allowed-with-wildcard",
			role: buildTestRole(t, map[string]interface{}{
				"allow_subdomains":            true,
				"allowed_domains":             []string{"test.com"},
				"allow_wildcard_certificates": false,
			}),
			identifiers: _buildACMEIdentifiers("*.test.com"),
			expectErr:   true,
		},
		{
			name: "role-wildcards-allowed-with-wildcard",
			role: buildTestRole(t, map[string]interface{}{
				"allow_subdomains":            true,
				"allowed_domains":             []string{"test.com"},
				"allow_wildcard_certificates": true,
			}),
			identifiers: _buildACMEIdentifiers("*.test.com"),
			expectErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := b.validateIdentifiersAgainstRole(tt.role, tt.identifiers)

			if tt.expectErr {
				require.Error(t, err, "validateIdentifiersAgainstRole(%v, %v)", tt.role.ToResponseData(), tt.identifiers)
				// If we did return an error if should be classified as a ErrRejectedIdentifier
				require.ErrorIs(t, err, ErrRejectedIdentifier)
			} else {
				require.NoError(t, err, "validateIdentifiersAgainstRole(%v, %v)", tt.role.ToResponseData(), tt.identifiers)
			}
		})
	}
}

// Test_parseOrderIdentifiers validates we convert ACME requests into proper ACMEIdentifiers
func Test_parseOrderIdentifiers(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		want    *ACMEIdentifier
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ipv4",
			data: map[string]interface{}{"type": "ip", "value": "192.168.1.1"},
			want: &ACMEIdentifier{
				Type:          ACMEIPIdentifier,
				Value:         "192.168.1.1",
				OriginalValue: "192.168.1.1",
				IsWildcard:    false,
				IsV6IP:        false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "ipv6",
			data: map[string]interface{}{"type": "ip", "value": "2001:0:130F::9C0:876A:130B"},
			want: &ACMEIdentifier{
				Type:          ACMEIPIdentifier,
				Value:         "2001:0:130F::9C0:876A:130B",
				OriginalValue: "2001:0:130F::9C0:876A:130B",
				IsWildcard:    false,
				IsV6IP:        true,
			},
			wantErr: assert.NoError,
		},
		{
			name: "ipv4-in-ipv6",
			data: map[string]interface{}{"type": "ip", "value": "::ffff:192.168.1.1"},
			want: &ACMEIdentifier{
				Type:          ACMEIPIdentifier,
				Value:         "::ffff:192.168.1.1",
				OriginalValue: "::ffff:192.168.1.1",
				IsWildcard:    false,
				IsV6IP:        true,
			},
			wantErr: assert.NoError,
		},
		{
			name: "dns",
			data: map[string]interface{}{"type": "dns", "value": "dadgarcorp.com"},
			want: &ACMEIdentifier{
				Type:          ACMEDNSIdentifier,
				Value:         "dadgarcorp.com",
				OriginalValue: "dadgarcorp.com",
				IsWildcard:    false,
				IsV6IP:        false,
			},
			wantErr: assert.NoError,
		},
		{
			name: "wildcard-dns",
			data: map[string]interface{}{"type": "dns", "value": "*.dadgarcorp.com"},
			want: &ACMEIdentifier{
				Type:          ACMEDNSIdentifier,
				Value:         "dadgarcorp.com",
				OriginalValue: "*.dadgarcorp.com",
				IsWildcard:    true,
				IsV6IP:        false,
			},
			wantErr: assert.NoError,
		},
		{
			name:    "ipv6-with-zone", // This is debatable if we should strip or fail
			data:    map[string]interface{}{"type": "ip", "value": "fe80::1cc0:3e8c:119f:c2e1%ens18"},
			wantErr: ErrorContains("IPv6 identifiers with zone information are not allowed"),
		},
		{
			name:    "bad-dns-wildcard",
			data:    map[string]interface{}{"type": "dns", "value": "*192.168.1.1"},
			wantErr: ErrorContains("invalid wildcard"),
		},
		{
			name:    "ip-in-dns",
			data:    map[string]interface{}{"type": "dns", "value": "192.168.1.1"},
			wantErr: ErrorContains("parsed OK as IP address"),
		},
		{
			name:    "empty-identifiers",
			data:    nil,
			wantErr: ErrorContains("no parsed identifiers were found"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			identifiers := map[string]interface{}{"identifiers": []interface{}{}}
			if tt.data != nil {
				identifiers["identifiers"] = append(identifiers["identifiers"].([]interface{}), tt.data)
			}
			got, err := parseOrderIdentifiers(identifiers)
			if !tt.wantErr(t, err, fmt.Sprintf("parseOrderIdentifiers(%v)", tt.data)) {
				return
			} else if err != nil {
				// If we passed the test above and an error was set no point in testing below
				return
			}

			require.Len(t, got, 1, "expected a single return value")
			acmeId := got[0]
			require.Equal(t, tt.want.Type, acmeId.Type)
			require.Equal(t, tt.want.Value, acmeId.Value)
			require.Equal(t, tt.want.OriginalValue, acmeId.OriginalValue)
			require.Equal(t, tt.want.IsWildcard, acmeId.IsWildcard)
			require.Equal(t, tt.want.IsV6IP, acmeId.IsV6IP)
		})
	}
}

func ErrorContains(errMsg string) assert.ErrorAssertionFunc {
	return func(t assert.TestingT, err error, i ...interface{}) bool {
		if err == nil {
			return assert.Fail(t, "expected error got none", i...)
		}

		if !strings.Contains(err.Error(), errMsg) {
			return assert.Fail(t, fmt.Sprintf("error did not contain '%s':\n%+v", errMsg, err), i...)
		}

		return true
	}
}

func _buildACMEIdentifiers(values ...string) []*ACMEIdentifier {
	var identifiers []*ACMEIdentifier

	for _, value := range values {
		identifiers = append(identifiers, _buildACMEIdentifier(value))
	}

	return identifiers
}

func _buildACMEIdentifier(val string) *ACMEIdentifier {
	ip := net.ParseIP(val)
	if ip == nil {
		identifier := &ACMEIdentifier{Type: "dns", Value: val, OriginalValue: val, IsWildcard: false}
		_, _, _ = identifier.MaybeParseWildcard()
		return identifier
	}

	return &ACMEIdentifier{Type: "ip", Value: val, OriginalValue: val, IsWildcard: false}
}

// Easily allow tests to create valid roles with proper defaults, since we don't have an easy
// way to generate roles with proper defaults, go through the createRole handler with the handlers
// field data so we pickup all the defaults specified there.
func buildTestRole(t *testing.T, config map[string]interface{}) *issuing.RoleEntry {
	b, s := CreateBackendWithStorage(t)

	path := pathRoles(b)
	fields := path.Fields
	if config == nil {
		config = map[string]interface{}{}
	}

	if _, exists := config["name"]; !exists {
		config["name"] = genUuid()
	}

	_, err := b.pathRoleCreate(ctx, &logical.Request{Storage: s}, &framework.FieldData{Raw: config, Schema: fields})
	require.NoError(t, err, "failed generating role with config %v", config)

	role, err := b.GetRole(ctx, s, config["name"].(string))
	require.NoError(t, err, "failed loading stored role")

	return role
}

// acmeFinalizeTestEnv holds a PKI backend with a usable issuer plus a ready
// ACME order, so a test can drive acmeFinalizeOrderHandler directly. It exists
// to exercise the per-order locking and processing-state handling added for
// GH-31987.
type acmeFinalizeTestEnv struct {
	b       *backend
	storage logical.Storage
	issuer  *issuing.IssuerEntry
	role    *issuing.RoleEntry
	baseUrl *url.URL
	uc      *jwsCtx
	account *acmeAccount
	orderId string
	csrB64  string
}

// setupACMEFinalizeTest builds a backend with a real root issuer and persists a
// ready order for the given DNS identifier, along with a CSR that matches it.
func setupACMEFinalizeTest(t *testing.T, identifier string) *acmeFinalizeTestEnv {
	t.Helper()

	b, s := CreateBackendWithStorage(t)

	_, err := CBWrite(b, s, "config/cluster", map[string]interface{}{
		"path": "https://localhost:8200/v1/pki",
	})
	require.NoError(t, err)

	_, err = CBWrite(b, s, "config/acme", map[string]interface{}{"enabled": true})
	require.NoError(t, err)

	_, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "root.example.com",
		"issuer_name": "root",
		"key_type":    "ec",
	})
	require.NoError(t, err, "failed generating root issuer")

	sc := b.makeStorageContext(ctx, s)
	issuer, err := getAcmeIssuer(sc, "")
	require.NoError(t, err, "failed resolving default acme issuer")

	accountId := genUuid()
	orderId := genUuid()

	accountKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	order := &acmeOrder{
		OrderId:          orderId,
		AccountId:        accountId,
		Status:           ACMEOrderReady,
		Expires:          time.Now().Add(time.Hour),
		Identifiers:      _buildACMEIdentifiers(identifier),
		AuthorizationIds: []string{genUuid()},
	}
	require.NoError(t, b.GetAcmeState().SaveOrder(&acmeContext{sc: sc}, order))

	baseUrl, err := url.Parse("https://localhost:8200/v1/pki/acme/")
	require.NoError(t, err)

	// The CSR is signed with a key distinct from the account key (so
	// validateCsrNotUsingAccountKey passes) and carries the order's DNS
	// identifier (so validateCsrMatchesOrder passes).
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		DNSNames: []string{identifier},
	}, csrKey)
	require.NoError(t, err)

	return &acmeFinalizeTestEnv{
		b:       b,
		storage: s,
		issuer:  issuer,
		role:    issuing.SignVerbatimRole(),
		baseUrl: baseUrl,
		uc:      &jwsCtx{Kid: accountId, Key: jose.JSONWebKey{Key: accountKey}},
		account: &acmeAccount{KeyId: accountId, Status: AccountStatusValid},
		orderId: orderId,
		csrB64:  base64.RawURLEncoding.EncodeToString(csrDER),
	}
}

// runFinalize invokes the finalize handler with a fresh storage context, the
// way a real request would. It is safe to call concurrently because the only
// shared state is the backend storage and the per-order lock pool.
func (e *acmeFinalizeTestEnv) runFinalize() (*logical.Response, error) {
	sc := e.b.makeStorageContext(ctx, e.storage)
	ac := &acmeContext{
		baseUrl:     e.baseUrl,
		sc:          sc,
		acmeState:   e.b.GetAcmeState(),
		runtimeOpts: acmeWrapperOpts{},
	}
	ac.Issuer = e.issuer
	ac.Role = e.role

	fields := &framework.FieldData{
		Raw:    map[string]interface{}{"order_id": e.orderId},
		Schema: map[string]*framework.FieldSchema{"order_id": {Type: framework.TypeString}},
	}
	data := map[string]interface{}{"csr": e.csrB64}

	return e.b.acmeFinalizeOrderHandler(ac, &logical.Request{Storage: e.storage}, fields, e.uc, data, e.account)
}

func (e *acmeFinalizeTestEnv) loadOrder(t *testing.T) *acmeOrder {
	t.Helper()
	order, err := e.b.GetAcmeState().LoadOrder(&acmeContext{sc: e.b.makeStorageContext(ctx, e.storage)}, e.uc, e.orderId)
	require.NoError(t, err)
	return order
}

func (e *acmeFinalizeTestEnv) countStoredCerts(t *testing.T) int {
	t.Helper()
	serials, err := e.storage.List(ctx, issuing.PathCerts)
	require.NoError(t, err)
	return len(serials)
}

// TestACMEFinalizeOrder_ConcurrentRequestsIssueSingleCert fires several finalize
// requests for the same order at once. Without the per-order lock, more than one
// would pass the readiness gate and issue its own certificate, leaving every
// certificate but the last one orphaned from the order. With the fix, exactly
// one request issues a certificate and the rest are rejected. See GH-31987.
func TestACMEFinalizeOrder_ConcurrentRequestsIssueSingleCert(t *testing.T) {
	t.Parallel()

	env := setupACMEFinalizeTest(t, "test.example.com")
	certsBefore := env.countStoredCerts(t)

	const workers = 8
	var wg sync.WaitGroup
	results := make([]error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, results[idx] = env.runFinalize()
		}(i)
	}
	wg.Wait()

	successes := 0
	for _, err := range results {
		if err == nil {
			successes++
			continue
		}
		require.ErrorIs(t, err, ErrOrderNotReady, "a losing finalize must be rejected by the readiness gate")
	}
	require.Equal(t, 1, successes, "exactly one finalize must succeed, got errors: %v", results)

	// Only one certificate may have been issued for the order.
	require.Equal(t, certsBefore+1, env.countStoredCerts(t),
		"concurrent finalize must not issue more than one certificate")

	order := env.loadOrder(t)
	require.Equal(t, ACMEOrderValid, order.Status)
	require.NotEmpty(t, order.CertificateSerialNumber)

	stored, err := env.storage.Get(ctx, issuing.PathCerts+order.CertificateSerialNumber)
	require.NoError(t, err)
	require.NotNil(t, stored, "the order's serial must reference a stored certificate, not an orphan")
}

// TestACMEFinalizeOrder_FailedIssuanceLeavesOrderRetryable checks that a
// finalize which fails during issuance does not alter the stored order, so the
// client can simply retry. The handler must not leave behind a half-finished
// certificate or any state that strands the order.
func TestACMEFinalizeOrder_FailedIssuanceLeavesOrderRetryable(t *testing.T) {
	t.Parallel()

	env := setupACMEFinalizeTest(t, "test.example.com")
	certsBefore := env.countStoredCerts(t)

	// A role that does not permit the order's identifier makes signCert reject
	// the CSR, so issuance fails partway through finalize.
	env.role = buildTestRole(t, map[string]interface{}{
		"allowed_domains":    []string{"unrelated.example.com"},
		"allow_bare_domains": true,
		"allow_subdomains":   false,
	})

	_, err := env.runFinalize()
	require.Error(t, err, "issuance with a non-matching role must fail")

	order := env.loadOrder(t)
	require.Equal(t, ACMEOrderReady, order.Status, "a failed issuance must leave the order finalizable")
	require.Empty(t, order.CertificateSerialNumber)
	require.Equal(t, certsBefore, env.countStoredCerts(t), "a failed issuance must not leave a certificate behind")

	// Retrying with a permissive role now succeeds, proving the order was not
	// stranded by the failed attempt.
	env.role = issuing.SignVerbatimRole()
	resp, err := env.runFinalize()
	require.NoError(t, err, "retry after a transient issuance failure must succeed")
	require.NotNil(t, resp)

	final := env.loadOrder(t)
	require.Equal(t, ACMEOrderValid, final.Status)
	require.NotEmpty(t, final.CertificateSerialNumber)
	require.Equal(t, certsBefore+1, env.countStoredCerts(t))
}

// TestACMEFinalizeOrder_OrderLockIsStablePerOrder confirms the striped lock pool
// is wired up and returns a stable lock per order id.
func TestACMEFinalizeOrder_OrderLockIsStablePerOrder(t *testing.T) {
	t.Parallel()

	b, _ := CreateBackendWithStorage(t)
	state := b.GetAcmeState()

	require.NotNil(t, state.orderLocks)
	require.Same(t, state.orderLockFor("order-abc"), state.orderLockFor("order-abc"),
		"the same order id must always map to the same lock")
}
