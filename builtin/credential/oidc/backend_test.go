package oidc

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/coreos/go-oidc/jose"
	oidcutil "github.com/coreos/go-oidc/oidc"
	"github.com/hashicorp/vault/builtin/credential/oidc/oidctesting"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	log "github.com/mgutz/logxi/v1"
)

var (
	goodClientId       = "myGoodClientId"
	goodSubjectId      = "DEADBEEF"
	usernameClaim      = "test_user_claim"
	groupClaim         = "test_group_claim"
	defaultLeaseTTLVal = time.Hour * 24
	maxLeaseTTLVal     = time.Hour * 24 * 32
)

// oidcProviderHarness starts a small httptest-based OIDC provider for the purpose of tests.
type oidcProviderHarness struct {
	tempDir  string
	provider *oidctesting.OIDCProvider
	server   *httptest.Server
}

func NewOidcProviderHarness(t *testing.T) *oidcProviderHarness {
	var err error
	h := &oidcProviderHarness{}
	h.tempDir = os.TempDir()
	oidctesting.GenerateSelfSignedCert(t, "127.0.0.1", path.Join(h.tempDir, "cert-1"), path.Join(h.tempDir, h.tempDir, "key-1"))
	h.provider = oidctesting.NewOIDCProvider(t, "")
	h.server, err = h.provider.ServeTLSWithKeyPair(path.Join(h.tempDir, "cert-1"), path.Join(h.tempDir, h.tempDir, "key-1"))
	if err != nil {
		t.Fatalf("failed starting oidc test harness: %v", err)
	}
	return h
}

func (h *oidcProviderHarness) issuerUrl(t *testing.T) string {
	return h.server.URL
}

func (h *oidcProviderHarness) genIdentityToken(t *testing.T, username string, groups []string, iat time.Time, exp time.Time) string {
	claims := oidcutil.NewClaims(h.issuerUrl(t), goodSubjectId, goodClientId, iat, exp)
	claims.Add(usernameClaim, username)
	if len(groups) > 0 {
		claims.Add(groupClaim, groups)
	}
	signer := h.provider.PrivKey.Signer()
	jwt, err := jose.NewSignedJWT(claims, signer)
	if err != nil {
		t.Fatalf("Cannot generate token: %v", err)
	}
	return jwt.Encode()
}

func (h *oidcProviderHarness) getCaCertBundle(t *testing.T) string {
	data, err := ioutil.ReadFile(path.Join(h.tempDir, "cert-1"))
	if err != nil {
		t.Fatalf("failed opening CA cert from testdir of harness %v", err)
	}
	return string(data)
}

func (h *oidcProviderHarness) Close() error {
	h.server.Close()
	os.RemoveAll(h.tempDir)
	return nil
}

func TestBackend_HappyPath_WithGroupName(t *testing.T) {
	h := NewOidcProviderHarness(t)
	defer h.Close()

	b, err := Factory(&logical.BackendConfig{
		Logger: logformat.NewVaultLogger(log.LevelTrace),
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
	if err != nil {
		t.Fatalf("Unable to create backend: %s", err)
	}

	config_data_good := map[string]interface{}{
		"issuer_url":       h.issuerUrl(t),
		"client_ids":       []string{goodClientId},
		"username_claim":   usernameClaim,
		"groups_claim":     groupClaim,
		"issuer_verify_ca": h.getCaCertBundle(t),
	}

	goodUserName := "myuser@example.com"
	goodGroupName := "myGroup"
	goodToken := h.genIdentityToken(t, goodUserName, []string{goodGroupName}, time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour))
	t.Logf("show token: %v", goodToken)
	//expiredToken := h.genIdentityToken(t, goodUserName, []string{goodGroupName}, time.Now().Add(-1 * time.Hour), time.Now().Add(1 * time.Hour))

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: false,
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testConfigWrite(t, config_data_good),
			testAccGroups(t, goodGroupName, "somepolicy"),
			//testLoginWrite(t, "dummytoken", 0, true),
			testLoginWrite(t, goodToken, 2*time.Hour, false),
			//testLoginWrite(t, login_data, expectedTTL1.Nanoseconds(), false),
			//testConfigWrite(t, config_data2),
			//testLoginWrite(t, login_data, expectedTTL2.Nanoseconds(), false),
			//testConfigWrite(t, config_data3),
			//testLoginWrite(t, login_data, 0, true),
		},
	})
}

func testLoginWrite(t *testing.T, token string, expectedTTL time.Duration, expectFail bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "login",
		ErrorOk:   true,
		Data:      map[string]interface{}{"token": token},
		Check: func(resp *logical.Response) error {
			if resp.IsError() && expectFail {
				return nil
			} else if expectFail {
				return fmt.Errorf("expected failiure")
			} else if resp.IsError() {
				return fmt.Errorf("expected to succeed but got %v", resp.Error())
			}
			t.Logf("resp.Auth %v", resp.Auth)

			actualTTL := resp.Auth.LeaseOptions.TTL
			ttlDiff := expectedTTL - actualTTL
			t.Logf("actualTtl %v expectedTtl %v ttldiff %v", actualTTL, expectedTTL, ttlDiff)
			if -5*time.Second < ttlDiff && ttlDiff < 5*time.Second { // 5s grace period
				return fmt.Errorf("TTL mismatched. Expected: %d Actual: %d", expectedTTL, resp.Auth.LeaseOptions.TTL.Nanoseconds())
			}
			return nil
		},
	}
}

func testConfigWrite(t *testing.T, d map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      d,
	}
}

func testAccGroups(t *testing.T, group string, policies string) logicaltest.TestStep {
	t.Logf("[testAccGroups] - Registering group %s, policy %s", group, policies)
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "groups/" + group,
		Data: map[string]interface{}{
			"policies": policies,
		},
	}
}
