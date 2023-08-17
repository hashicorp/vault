package api_capability

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	scada "github.com/hashicorp/hcp-scada-provider"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/hcp_link/internal"
)

const (
	// HCP policy fetch retry limits
	policyRetryWaitMin = 500 * time.Millisecond
	policyRetryWaitMax = 30 * time.Second
	policyRetryMax     = 3

	// batchTokenDefaultTTL default TTL of a batch token
	batchTokenDefaultTTL = 5 * time.Minute

	// tokenExpiryOffset is deducted from token expiry to make sure a new token
	// is generated before an old one is expired.
	// This is used when Vault is unsealed.
	tokenExpiryOffset = 5 * time.Second
)

type HCPLinkTokenManager struct {
	lock            sync.RWMutex
	wrappedCore     internal.WrappedCoreHCPToken
	logger          hclog.Logger
	scadaConfig     *scada.Config
	policyUrl       string
	latestToken     string
	tokenTTL        time.Duration
	policy          string
	lastTokenExpiry time.Time
}

func (t *HCPLinkTokenManager) SetTokenTTL(ttl time.Duration) error {
	if ttl < tokenExpiryOffset {
		return fmt.Errorf("ttl cannot be less than 5 seconds")
	}

	t.lock.Lock()
	defer t.lock.Unlock()
	t.tokenTTL = ttl

	return nil
}

func (t *HCPLinkTokenManager) GetLatestToken() string {
	t.lock.RLock()
	latestToken := t.latestToken
	t.lock.RUnlock()

	return latestToken
}

func (t *HCPLinkTokenManager) GetLastTokenExpiry() time.Time {
	t.lock.RLock()
	te := t.lastTokenExpiry
	t.lock.RUnlock()

	return te
}

func (t *HCPLinkTokenManager) GetPolicy() string {
	t.lock.RLock()
	policy := t.policy
	t.lock.RUnlock()

	return policy
}

func (t *HCPLinkTokenManager) fetchPolicy() (string, error) {
	req, err := http.NewRequest(http.MethodGet, t.policyUrl, nil)
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}

	retryableReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return "", fmt.Errorf("error adding HTTP request retry wrapping: %w", err)
	}

	token, err := t.scadaConfig.HCPConfig.Token()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve HCP bearer token: %w", err)
	}

	retryableReq.Header.Add("Authorization", "Bearer "+token.AccessToken)

	client := &retryablehttp.Client{
		HTTPClient:   cleanhttp.DefaultClient(),
		RetryWaitMin: policyRetryWaitMin,
		RetryWaitMax: policyRetryWaitMax,
		RetryMax:     policyRetryMax,
		CheckRetry:   retryablehttp.DefaultRetryPolicy,
		Backoff:      retryablehttp.DefaultBackoff,
	}

	resp, err := client.Do(retryableReq)
	if err != nil {
		return "", fmt.Errorf("error retrieving policy from HCP: %w", err)
	}

	if resp.Body == nil {
		return "", fmt.Errorf("invalid HCP policy response: %w", err)
	}

	bodyBytes := bytes.NewBuffer(nil)
	_, err = bodyBytes.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading HCP policy response: %w", err)
	}

	if err := resp.Body.Close(); err != nil {
		return "", fmt.Errorf("error closing response body: %w", err)
	}

	body := make(map[string]interface{})
	err = json.Unmarshal(bodyBytes.Bytes(), &body)
	if err != nil {
		return "", fmt.Errorf("error parsing response body: %w", err)
	}

	policy, ok := body["policy"].(string)
	if !ok {
		return "", fmt.Errorf("formatting for policy fetched from HCP is invalid, expected string received %T", body["policy"])
	}

	return policy, nil
}

func (t *HCPLinkTokenManager) updateInLinePolicy() {
	policy, err := t.fetchPolicy()
	if err != nil {
		t.logger.Error("failed to fetch policy from HCP", "error", err)
		return
	}

	t.logger.Info("new policy fetched from HCP")

	t.lock.Lock()
	t.policy = policy
	t.lock.Unlock()

	return
}

func NewHCPLinkTokenManager(scadaConfig *scada.Config, core *vault.Core, logger hclog.Logger) (*HCPLinkTokenManager, error) {
	tokenLogger := logger.Named("token_manager")

	policyURL := fmt.Sprintf("https://%s/vault-link/2022-11-07/organizations/%s/projects/%s/link/policy/%s",
		scadaConfig.HCPConfig.APIAddress(),
		scadaConfig.Resource.Location.OrganizationID,
		scadaConfig.Resource.Location.ProjectID,
		scadaConfig.Resource.ID,
	)

	m := &HCPLinkTokenManager{
		wrappedCore:     core,
		logger:          tokenLogger,
		scadaConfig:     scadaConfig,
		policyUrl:       policyURL,
		tokenTTL:        batchTokenDefaultTTL,
		lastTokenExpiry: time.Time{},
	}

	m.logger.Info("initialized HCP Link token manager")

	return m, nil
}

func (m *HCPLinkTokenManager) Shutdown() {
	m.ForgetTokenPolicy()
}

// HandleTokenPolicy checks if Vault is sealed or not an active node,
// then it removes both token and policy. And, if Vault is not sealed,
// and token needs to be refreshed it refreshes both policy and token
func (m *HCPLinkTokenManager) HandleTokenPolicy(ctx context.Context, activeNode bool) string {
	switch {
	case m.wrappedCore.Sealed(), !activeNode:
		m.logger.Debug("failed to create a token as Vault is either sealed or a non-active node. Setting the token to an empty string")
		m.ForgetTokenPolicy()
	case m.GetLatestToken() == "", m.GetLastTokenExpiry().Before(time.Now()):
		m.createToken(ctx)
	}

	return m.GetLatestToken()
}

// ForgetTokenPolicy Forgets the current Batch token, its associated policy,
// and sets the lastTokenExpiry to Now such that the policy is forced to be
// refreshed by the next valid request.
func (m *HCPLinkTokenManager) ForgetTokenPolicy() {
	m.lock.Lock()
	defer m.lock.Unlock()

	// purging the latest token and policy
	m.latestToken = ""
	m.policy = ""

	// setting the token expiry to now as a cleanup step
	m.lastTokenExpiry = time.Now()
	m.logger.Info("purged token and inline policy")
}

func (m *HCPLinkTokenManager) createToken(ctx context.Context) {
	// updating the policy first
	m.updateInLinePolicy()
	policy := m.GetPolicy()

	m.lock.Lock()
	defer m.lock.Unlock()

	// an orphan batch token is required.
	// For an orphan token, we need to not set the parent
	// Also setting the time of creation, and metadata for auditing
	te := &logical.TokenEntry{
		Type:               logical.TokenTypeBatch,
		TTL:                m.tokenTTL,
		CreationTime:       time.Now().Unix(),
		NamespaceID:        namespace.RootNamespaceID,
		NoIdentityPolicies: true,
		Meta: map[string]string{
			"hcp_link_token": "HCP Link Access Token",
		},
		InlinePolicy: policy,
		InternalMeta: map[string]string{vault.IgnoreForBilling: "true"},
	}
	// try creating the token, if the trial fails, let's set the latestToken
	// to an empty string, and reset token expiry with a backoff strategy
	err := m.wrappedCore.CreateToken(ctx, te)
	if err != nil {
		m.logger.Error("failed to create a token, setting the token to an empty string", "error", err)
		m.latestToken = ""
		m.lastTokenExpiry = time.Now()
		return
	}

	if te == nil || len(te.ID) == 0 {
		m.logger.Error("token creation returned an empty token entry")
		m.latestToken = ""
		m.lastTokenExpiry = time.Now()
		return
	}

	m.latestToken = te.ID
	// storing the new token and setting lastTokenExpiry to be 5 seconds before
	// the TTL
	m.lastTokenExpiry = time.Now().Add(m.tokenTTL - tokenExpiryOffset)

	m.logger.Info("successfully generated a new token for HCP link")
}
