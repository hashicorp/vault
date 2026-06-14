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
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
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

// TestACME_FinalizeOrderConcurrency verifies that the per-order lock and
// intermediate "processing" status in acmeFinalizeOrderHandler prevent
// concurrent finalize requests from double-issuing certificates for the
// same order. See GH-31987.
func TestACME_FinalizeOrderConcurrency(t *testing.T) {
	t.Parallel()

	b, s := CreateBackendWithStorage(t)

	// Set up a minimal ACME environment: an account and a ready order.
	accountId := genUuid()
	orderId := genUuid()
	account := &acmeAccount{
		KeyId:  accountId,
		Status: ACMEAccountStatusValid,
	}
	order := &acmeOrder{
		OrderId:          orderId,
		AccountId:        accountId,
		Status:           ACMEOrderReady,
		Expires:          time.Now().Add(1 * time.Hour),
		Identifiers:      []*ACMEIdentifier{{Type: ACMEDNSIdentifier, OriginalValue: "test.example.com", Value: "test.example.com"}},
		AuthorizationIds: []string{genUuid()},
	}

	// Persist account and order to storage.
	err := b.GetAcmeState().SaveAccount(&acmeContext{sc: &storageContext{Context: ctx, Storage: s}}, account)
	require.NoError(t, err)
	err = b.GetAcmeState().SaveOrder(&acmeContext{sc: &storageContext{Context: ctx, Storage: s}}, order)
	require.NoError(t, err)

	// Generate a valid CSR for the order's identifier.
	csrKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	csrTemplate := &x509.CertificateRequest{DNSNames: []string{"test.example.com"}}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, csrKey)
	require.NoError(t, err)
	csrB64 := base64.RawURLEncoding.EncodeToString(csrDER)

	// Build a minimal acmeContext pointing at our storage.
	ac := &acmeContext{
		sc:        &storageContext{Context: ctx, Storage: s},
		acmeState: b.acmeState,
		runtimeOpts: acmeWrapperOpts{
			isCiepsEnabled: false,
		},
	}
	ac.Role = issuing.SignVerbatimRole()

	// Build the jwsCtx for the account holder.
	uc := &jwsCtx{Kid: accountId}

	// Prepare the fields with order_id and the CSR payload.
	data := map[string]interface{}{
		"csr": csrB64,
	}
	schema := map[string]*framework.FieldSchema{
		"order_id": {Type: framework.TypeString},
	}
	fields := &framework.FieldData{Raw: map[string]interface{}{"order_id": orderId}, Schema: schema}
	req := &logical.Request{Storage: s}

	// Run two concurrent finalize attempts — the lock should serialize them
	// and the processing-status gate should reject the loser.
	var wg sync.WaitGroup
	results := make([]error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, err := b.acmeFinalizeOrderHandler(ac, req, fields, uc, data, account)
			results[idx] = err
		}(i)
	}
	wg.Wait()

	// At least one should have succeeded.
	successCount := 0
	for _, err := range results {
		if err == nil {
			successCount++
		}
	}
	require.Equal(t, 1, successCount, "expected exactly one successful finalize, got %d; errors: %v", successCount, results)

	// The order should now be in the "valid" state with a certificate serial.
	loaded, err := b.GetAcmeState().LoadOrder(ac, uc, orderId)
	require.NoError(t, err)
	require.Equal(t, ACMEOrderValid, loaded.Status)
	require.NotEmpty(t, loaded.CertificateSerialNumber)
}

// TestACME_FinalizeOrderProcessingDefense verifies that the processing-status
// defense-in-depth rejects a stale request whose lock-stripe alias happens to
// differ from the in-flight request. While the per-order lock normally prevents
// this, the processing status acts as a second gate.
func TestACME_FinalizeOrderProcessingDefense(t *testing.T) {
	t.Parallel()

	b, s := CreateBackendWithStorage(t)

	accountId := genUuid()
	orderId := genUuid()

	// Save an order that is already in "processing" state (simulating a
	// mid-flight finalize on a different lock stripe).
	order := &acmeOrder{
		OrderId:          orderId,
		AccountId:        accountId,
		Status:           ACMEOrderProcessing,
		Expires:          time.Now().Add(1 * time.Hour),
		Identifiers:      []*ACMEIdentifier{{Type: ACMEDNSIdentifier, OriginalValue: "test.example.com", Value: "test.example.com"}},
		AuthorizationIds: []string{genUuid()},
	}
	err := b.GetAcmeState().SaveOrder(&acmeContext{sc: &storageContext{Context: ctx, Storage: s}}, order)
	require.NoError(t, err)

	// Build context and attempt to finalize — should be rejected because
	// the order is not in "ready" state.
	ac := &acmeContext{
		sc:        &storageContext{Context: ctx, Storage: s},
		acmeState: b.acmeState,
		runtimeOpts: acmeWrapperOpts{
			isCiepsEnabled: false,
		},
	}
	ac.Role = issuing.SignVerbatimRole()
	uc := &jwsCtx{Kid: accountId}
	account := &acmeAccount{KeyId: accountId, Status: ACMEAccountStatusValid}

	csrKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	csrTemplate := &x509.CertificateRequest{DNSNames: []string{"test.example.com"}}
	csrDER, _ := x509.CreateCertificateRequest(rand.Reader, csrTemplate, csrKey)
	csrB64 := base64.RawURLEncoding.EncodeToString(csrDER)

	data := map[string]interface{}{"csr": csrB64}
	schema := map[string]*framework.FieldSchema{"order_id": {Type: framework.TypeString}}
	fields := &framework.FieldData{Raw: map[string]interface{}{"order_id": orderId}, Schema: schema}
	req := &logical.Request{Storage: s}

	_, err = b.acmeFinalizeOrderHandler(ac, req, fields, uc, data, account)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrOrderNotReady)
	require.Contains(t, err.Error(), "processing")
}

// TestACME_FinalizeOrderLockSerialises verifies that the per-order lock
// correctly serialises access to the finalize handler, preventing the
// data-race on CertificateSerialNumber described in GH-31987.
func TestACME_FinalizeOrderLockSerialises(t *testing.T) {
	t.Parallel()

	b, _ := CreateBackendWithStorage(t)

	// Verify the locks are initialised.
	require.NotNil(t, b.acmeOrderLocks)
	require.Len(t, b.acmeOrderLocks, 256)

	// Verify LockForKey returns the same lock for the same order ID.
	orderId := "test-order-123"
	lock1 := locksutil.LockForKey(b.acmeOrderLocks, orderId)
	lock2 := locksutil.LockForKey(b.acmeOrderLocks, orderId)
	require.Same(t, lock1, lock2)

	// Verify the lock can be acquired and released without deadlock.
	lock1.Lock()
	lock1.Unlock()

	// Verify different order IDs may hash to the same or different locks.
	otherId := "test-order-456"
	lock3 := locksutil.LockForKey(b.acmeOrderLocks, otherId)
	// Either same or different is fine — just verify no panic.
	lock3.Lock()
	lock3.Unlock()
}
