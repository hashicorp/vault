// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-2.0

package ldap

import (
	"bytes"
	"context"
	"encoding/binary"
	"os"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/hashicorp/vault/helper/testhelpers/ldap"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
)

// TestRotateRoot_DefaultSchema tests that the default schema type is used when not explicitly set and that rotation works in that case.
// This test relies on a docker ldap server with a suitable person object (cn=admin,dc=planetexpress,dc=com)
// with bindpassword "admin". `PrepareTestContainer` does this for us - see the backend_test for more details.
func TestRotateRoot_DefaultSchema(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip("skipping rotate root tests because VAULT_ACC is unset")
	}
	ctx := context.Background()

	b, store := createBackendWithStorage(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, ldap.DefaultVersion)
	defer cleanup()
	// set up auth config
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   store,
		Data: map[string]interface{}{
			"url":      cfg.Url,
			"binddn":   cfg.BindDN,
			"bindpass": cfg.BindPassword,
			"userdn":   cfg.UserDN,
		},
	}

	resp, err := b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to initialize ldap auth config: %s", err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to initialize ldap auth config: %s", resp.Data["error"])
	}

	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/rotate-root",
		Storage:   store,
	}

	_, err = b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to rotate password: %s", err)
	}

	newCFG, err := b.Config(ctx, req)
	if err != nil {
		t.Fatalf("failed to get config after rotation: %s", err)
	}
	if newCFG == nil {
		t.Fatal("config is nil after rotation")
	}
	if newCFG.BindDN != cfg.BindDN {
		t.Fatalf("a value in config that should have stayed the same changed: %s", cfg.BindDN)
	}
	if newCFG.BindPassword == cfg.BindPassword {
		t.Fatalf("the password should have changed, but it didn't")
	}
}

// TestRotateRootWithRotationUrl relies on a docker ldap server with a suitable person object (cn=admin,dc=planetexpress,dc=com)
// with bindpassword "admin". `PrepareTestContainer` does this for us. - see the backend_test for more details
// It checks that rotation url is being used instead of the main URL and assures that setting rotation url does't
// replace main URL
func TestRotateRootWithRotationUrl(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip("skipping rotate root tests because VAULT_ACC is unset")
	}
	ctx := context.Background()

	b, store := createBackendWithStorage(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, ldap.DefaultVersion)
	defer cleanup()
	const mainDummyUrl = "ldap://example.com:389"
	// set up auth config
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   store,
		Data: map[string]interface{}{
			"url":          mainDummyUrl,
			"binddn":       cfg.BindDN,
			"bindpass":     cfg.BindPassword,
			"userdn":       cfg.UserDN,
			"rotation_url": cfg.Url,
		},
	}

	resp, err := b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to initialize ldap auth config: %s", err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to initialize ldap auth config: %s", resp.Data["error"])
	}

	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/rotate-root",
		Storage:   store,
	}

	_, err = b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to rotate password: %s", err)
	}

	newCFG, err := b.Config(ctx, req)
	if err != nil {
		t.Fatalf("failed to get config after rotation: %s", err)
	}
	if newCFG == nil {
		t.Fatal("config is nil after rotation")
	}
	if newCFG.BindDN != cfg.BindDN {
		t.Fatalf("BindDN %q changed unexpectedly, found new value %q", cfg.BindDN, newCFG.BindDN)
	}
	if newCFG.BindPassword == cfg.BindPassword {
		t.Fatalf("the password should have changed, but it didn't")
	}
	// expecting the newCFG url to be "ldap://example.com:389"
	if newCFG.Url != mainDummyUrl {
		t.Fatalf("URL %q changed unexpectedly, found new value %q", mainDummyUrl, newCFG.Url)
	}
}

// TestRotateRoot_Schema_OpenLDAP tests that rotation for OpenLDAP schema
func TestRotateRoot_Schema_OpenLDAP(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip("skipping rotate root tests because VAULT_ACC is unset")
	}
	ctx := context.Background()
	b, store := createBackendWithStorage(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, ldap.DefaultVersion)
	defer cleanup()
	// set up auth config
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   store,
		Data: map[string]interface{}{
			"url":      cfg.Url,
			"binddn":   cfg.BindDN,
			"bindpass": cfg.BindPassword,
			"userdn":   cfg.UserDN,
			"schema":   ldaputil.SchemaOpenLDAP,
		},
	}
	resp, err := b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to initialize ldap auth config: %s", err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to initialize ldap auth config: %s", resp.Data["error"])
	}
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/rotate-root",
		Storage:   store,
	}
	_, err = b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to rotate password: %s", err)
	}
	newCFG, err := b.Config(ctx, req)
	if err != nil {
		t.Fatalf("failed to get config after rotation: %s", err)
	}
	if newCFG == nil {
		t.Fatal("config is nil after rotation")
	}
	if newCFG.BindDN != cfg.BindDN {
		t.Fatalf("a value in config that should have stayed the same changed: %s", cfg.BindDN)
	}
	if newCFG.BindPassword == cfg.BindPassword {
		t.Fatalf("the password should have changed, but it didn't")
	}
}

// TestRotateRoot_UnsupportedSchema tests unsupported schema handling
func TestRotateRoot_UnsupportedSchema(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip("skipping rotate root tests because VAULT_ACC is unset")
	}
	ctx := context.Background()
	b, store := createBackendWithStorage(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, ldap.DefaultVersion)
	defer cleanup()
	// set up auth config
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   store,
		Data: map[string]interface{}{
			"url":      cfg.Url,
			"binddn":   cfg.BindDN,
			"bindpass": cfg.BindPassword,
			"userdn":   cfg.UserDN,
			"schema":   "unsupported_schema", // unsupported schema type
		},
	}

	resp, err := b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to initialize ldap auth config: %s", err)
	}
	// Config creation should fail with logical error
	if resp == nil || !resp.IsError() {
		t.Fatal("expected config creation to fail with unsupported schema type, but it succeeded")
	}
	// Verify error message contains expected text
	errMsg := resp.Error().Error()
	if !strings.Contains(errMsg, "unsupported schema type") {
		t.Fatalf("expected error containing 'unsupported schema type', got: %s", errMsg)
	}
}

// TestRotateRoot_EncodeUTF16LEBytes tests the encoding of UTF-16LE bytes for AD password modification.
func TestRotateRoot_EncodeUTF16LEBytes(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty_string", input: ""},
		{name: "single_char", input: "A"},
		{name: "quoted_password", input: "\"password\""},
		{name: "alphanumeric_with_special", input: "Pass123!"},
		{name: "unicode_chars", input: "Pāsswörd✓"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeUTF16LEBytes(tt.input)
			r := utf16.Encode([]rune(tt.input))
			expected := make([]byte, len(r)*2)
			for i, v := range r {
				binary.LittleEndian.PutUint16(expected[i*2:], v)
			}

			if !bytes.Equal(got, expected) {
				t.Fatalf("encodeUTF16LEBytes(%q) = %v, want %v", tt.input, got, expected)
			}
		})
	}
}

// adTLSTestCase defines a test case for AD TLS validation
type adTLSTestCase struct {
	name           string
	url            string
	insecureTLS    bool
	startTLS       bool
	certificate    string
	passwordPolicy string
	expectError    bool
	errorContains  string
}

// TestRotateRoot_AD_TLS_Validation tests TLS validation for AD schema
func TestRotateRoot_AD_TLS_Validation(t *testing.T) {
	tests := []adTLSTestCase{
		// Valid: ldaps - TLS validation passes, connection refused immediately (no real server).
		{
			name:        "ldaps_valid",
			url:         "ldaps://127.0.0.1:1",
			insecureTLS: true,
		},
		// Valid: ldap + starttls - with optional cert and password policy
		{
			name:           "ldap_with_starttls_valid",
			url:            "ldap://127.0.0.1:1",
			startTLS:       true,
			certificate:    testCert,
			passwordPolicy: "test-policy",
		},
		// Invalid: ldap without starttls - must surface as LogicalErrorResponse
		{
			name:          "ldap_without_starttls_invalid",
			url:           "ldap://example.com",
			expectError:   true,
			errorContains: "AD password rotation with ldap:// requires starttls=true for encrypted connection",
		},
		// Invalid: unsupported protocol - must surface as LogicalErrorResponse
		{
			name:          "invalid_protocol",
			url:           "http://example.com",
			expectError:   true,
			errorContains: "AD password rotation requires ldap:// or ldaps:// protocol, got: http://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testADTLSValidation(t, tt)
		})
	}
}

// Valid self-signed certificate for testing (shared across tests)
const testCert = `-----BEGIN CERTIFICATE-----
MIIB1jCCAUGgAwIBAgIFAMv4K9YwCwYJKoZIhvcNAQELMCkxEDAOBgNVBAoTB0Fj
bWUgQ28xFTATBgNVBAMTDEVkZGFyZCBTdGFyazAeFw0xNTA1MDYwMzU2NDBaFw0x
NjA1MDYwMzU2NDBaMCUxEDAOBgNVBAoTB0FjbWUgQ28xETAPBgNVBAMTCEpvbiBT
bm93MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDK6NU0R0eiCYVquU4RcjKc
LzGfx0aa1lMr2TnLQUSeLFZHFxsyyMXXuMPig3HK4A7SGFHupO+/1H/sL4xpH5zg
8+Zg2r8xnnney7abxcuv0uATWSIeKlNnb1ZO1BAxFnESc3GtyOCr2dUwZHX5mRVP
+Zxp2ni5qHNraf3wE2VPIQIDAQABoxIwEDAOBgNVHQ8BAf8EBAMCAKAwCwYJKoZI
hvcNAQELA4GBAIr2F7wsqmEU/J/kLyrCgEVXgaV/sKZq4pPNnzS0tBYk8fkV3V18
sBJyHKRLL/wFZASvzDcVGCplXyMdAOCyfd8jO3F9Ac/xdlz10RrHJT75hNu3a7/n
9KNwKhfN4A1CQv2x372oGjRhCW5bHNCWx4PIVeNzCyq/KZhyY9sxHE1g
-----END CERTIFICATE-----`

// testADTLSValidation is a helper function that runs AD TLS validation tests
func testADTLSValidation(t *testing.T, tt adTLSTestCase) {
	t.Helper()
	ctx := context.Background()
	// Create fresh backend and storage for each subtest to avoid state bleeding
	b, store := createBackendWithStorage(t)

	data := map[string]interface{}{
		"url":             tt.url,
		"binddn":          "cn=admin,dc=example,dc=com",
		"bindpass":        "password",
		"userdn":          "ou=users,dc=example,dc=com",
		"schema":          ldaputil.SchemaAD,
		"insecure_tls":    tt.insecureTLS,
		"starttls":        tt.startTLS,
		"certificate":     tt.certificate,
		"password_policy": tt.passwordPolicy,
	}

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   store,
		Data:      data,
	}

	resp, err := b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error during config: %v", err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("unexpected response error during config: %v", resp.Error())
	}

	// Try to rotate - this is where TLS validation happens
	rotateReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/rotate-root",
		Storage:   store,
	}

	rotateResp, rotateErr := b.HandleRequest(ctx, rotateReq)

	if tt.expectError {
		// responseError is returned as logical.ErrorResponse with nil error
		if rotateResp == nil || !rotateResp.IsError() {
			t.Fatalf("expected error response, got resp=%v", rotateResp)
		}
		errMsg := rotateResp.Error().Error()
		if !strings.Contains(errMsg, tt.errorContains) {
			t.Fatalf("expected error containing %q, got: %q", tt.errorContains, errMsg)
		}
	} else {
		// For valid configs, TLS validation passes and we expect a connection error
		// (no real LDAP server). The backend logs these at ERROR level,
		// which is expected and does not indicate a test failure.
		// Assert neither the Go error nor the response error is a TLS validation failure,
		// proving execution got past the TLS checks.
		if rotateErr != nil {
			if isTLSValidationError(rotateErr.Error()) {
				t.Fatalf("unexpected TLS validation error: %s", rotateErr)
			}
			// connection/bind error is expected — TLS validation passed
		}
		if rotateResp != nil && rotateResp.IsError() {
			errMsg := rotateResp.Error().Error()
			if isTLSValidationError(errMsg) {
				t.Fatalf("unexpected TLS validation error in response: %s", errMsg)
			}
			// connection/bind error in response is expected — TLS validation passed
		}
	}
}

// isTLSValidationError checks if an error message is a TLS validation error
func isTLSValidationError(msg string) bool {
	return strings.Contains(msg, "AD password rotation with ldap:// requires starttls=true for encrypted connection") ||
		strings.Contains(msg, "AD password rotation requires ldap:// or ldaps:// protocol, got:")
}

// TestValidateADRotationURLs tests the URL validation logic for AD password rotation
func TestValidateADRotationURLs(t *testing.T) {
	tests := []struct {
		name      string
		urlString string
		startTLS  bool
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "single ldaps URL - valid",
			urlString: "ldaps://secure.example.com",
			startTLS:  false,
			wantErr:   false,
		},
		{
			name:      "single ldap URL with StartTLS - valid",
			urlString: "ldap://example.com",
			startTLS:  true,
			wantErr:   false,
		},
		{
			name:      "single ldap URL without StartTLS - invalid",
			urlString: "ldap://example.com",
			startTLS:  false,
			wantErr:   true,
			errMsg:    "starttls=true",
		},
		{
			name:      "multiple ldaps URLs - valid",
			urlString: "ldaps://server1.example.com,ldaps://server2.example.com",
			startTLS:  false,
			wantErr:   false,
		},
		{
			name:      "multiple ldap URLs with StartTLS - valid",
			urlString: "ldap://server1.example.com,ldap://server2.example.com",
			startTLS:  true,
			wantErr:   false,
		},
		{
			name:      "mixed ldaps and ldap with StartTLS - valid",
			urlString: "ldaps://server1.example.com,ldap://server2.example.com",
			startTLS:  true,
			wantErr:   false,
		},
		{
			name:      "mixed ldaps and ldap without StartTLS - invalid",
			urlString: "ldaps://server1.example.com,ldap://server2.example.com",
			startTLS:  false,
			wantErr:   true,
			errMsg:    "starttls=true",
		},
		{
			name:      "multiple ldap URLs without StartTLS - invalid",
			urlString: "ldap://server1.example.com,ldap://server2.example.com",
			startTLS:  false,
			wantErr:   true,
			errMsg:    "starttls=true",
		},
		{
			name:      "invalid protocol - invalid",
			urlString: "http://example.com",
			startTLS:  false,
			wantErr:   true,
			errMsg:    "ldap:// or ldaps://",
		},
		{
			name:      "empty URL - invalid",
			urlString: "",
			startTLS:  false,
			wantErr:   true,
			errMsg:    "requires a configured URL",
		},
		{
			name:      "URL with spaces - valid",
			urlString: "ldaps://server1.example.com, ldaps://server2.example.com",
			startTLS:  false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateADRotationURLs(tt.urlString, tt.startTLS)
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateADRotationURLs() expected error but got none")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateADRotationURLs() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateADRotationURLs() unexpected error = %v", err)
				}
			}
		})
	}
}
