package okta

import (
	"context"
	"fmt"
	"time"

	"encoding/json"
	"errors"
	"github.com/chrismalek/oktasdk-go/okta"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/mfa"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"net/url"
	"strings"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Root: mfa.MFARootPaths(),

			Unauthenticated: []string{
				"login/*",
			},
			SealWrapStorage: []string{
				"config",
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(&b),
			pathUsers(&b),
			pathGroups(&b),
			pathUsersList(&b),
			pathGroupsList(&b),
		},
			mfa.MFAPaths(b.Backend, pathLogin(&b))...,
		),

		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
	}

	return &b
}

type backend struct {
	*framework.Backend
}

// halResource is a Hypertext Application Language resource (used throughout the Okta API).
// See https://tools.ietf.org/html/draft-kelly-json-hal-06
// Embedding this Go struct helps deserialize dynamic embedded links and objects.
type halResource struct {
	Embedded embedded `mapstructure:"_embedded"`
	Links    links    `mapstructure:"_links"`
}

// embedded is a HAL "_embedded" object.
type embedded map[string]interface{}

// Decode decodes the named embedded object or object list.
// Dest should be a pointer to a struct, or a pointer to a struct slice for list-valued embeds.
func (m embedded) Decode(name string, dest interface{}) error {
	return mapstructure.Decode(m[name], dest)
}

// links is a HAL "_links" object. Annoyingly, each _links entry can be either a single link
// or a list.
type links map[string]interface{}

// Link returns a single link.
func (m links) Link(name string) (*link, error) {
	l := &link{}
	err := mapstructure.Decode(m[name], l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// link is a HAL link.
type link struct {
	Href string
}

type authResult struct {
	halResource  `mapstructure:",squash"`
	Status       string
	FactorResult string
	StateToken   string
}

func (r *authResult) UnmarshalJSON(data []byte) error {
	// Use mapstructure for JSON unmarshalling, since it plays better with HAL dynamic
	// embedded objects than static struct definitions.
	m := make(map[string]interface{})
	json.Unmarshal(data, &m)
	return mapstructure.Decode(m, r)
}

func (b *backend) Login(ctx context.Context, req *logical.Request, username string, password string) ([]string, *logical.Response, []string, error) {
	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, nil, nil, err
	}
	if cfg == nil {
		return nil, logical.ErrorResponse("Okta auth method not configured"), nil, nil
	}

	client := cfg.OktaClient()

	authReq, err := client.NewRequest("POST", "authn", map[string]interface{}{
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	result := &authResult{}
	_, err = client.Do(authReq, result)
	if err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("Okta auth failed: %v", err)), nil, nil
	}

	oktaResponse := &logical.Response{
		Data: map[string]interface{}{},
	}

	// More about Okta's Auth transation state here:
	// https://developer.okta.com/docs/api/resources/authn#transaction-state

	// If lockout failures are not configured to be hidden, the status needs to
	// be inspected for LOCKED_OUT status. Otherwise, it is handled above by an
	// error returned during the authentication request.
	switch result.Status {
	case "LOCKED_OUT":
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: user is locked out", "user", username)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil

	case "MFA_ENROLL", "MFA_ENROLL_ACTIVATE":
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: user must enroll or complete mfa enrollment", "user", username)
		}
		return nil, logical.ErrorResponse("okta authentication failed: you must complete MFA enrollment to continue"), nil, nil

	case "MFA_REQUIRED":
		if result, err = b.completeMfa(ctx, client, result); err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("Okta auth failed: %v", err)), nil, nil
		}

	case "PASSWORD_EXPIRED":
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: password is expired", "user", username)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil

	case "PASSWORD_WARN":
		oktaResponse.AddWarning("Your Okta password is in warning state and needs to be changed soon.")

	case "MFA_REQUIRED", "MFA_ENROLL":
		if !cfg.BypassOktaMFA {
			return nil, logical.ErrorResponse("okta mfa required for this account but mfa bypass not set in config"), nil, nil
		}

	case "SUCCESS":
		// Do nothing here

	default:
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: unhandled result status", "status", result.Status)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil
	}

	// Verify result status again in case a switch case above modifies result
	switch {
	case result.Status == "SUCCESS",
		result.Status == "PASSWORD_WARN",
		result.Status == "MFA_REQUIRED" && cfg.BypassOktaMFA,
		result.Status == "MFA_ENROLL" && cfg.BypassOktaMFA:
		// Allowed
	default:
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: authentication returned a non-success status", "status", result.Status)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil
	}

	var allGroups []string
	// Only query the Okta API for group membership if we have a token
	if cfg.Token != "" {
		user := &okta.User{}
		if err = result.Embedded.Decode("user", user); err != nil {
			return nil, nil, nil, err
		}
		oktaGroups, err := b.getOktaGroups(client, user)
		if err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("okta failure retrieving groups: %v", err)), nil, nil
		}
		if len(oktaGroups) == 0 {
			errString := fmt.Sprintf(
				"no Okta groups found; only policies from locally-defined groups available")
			oktaResponse.AddWarning(errString)
		}
		allGroups = append(allGroups, oktaGroups...)
	}

	// Import the custom added groups from okta backend
	user, err := b.User(ctx, req.Storage, username)
	if err != nil {
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: error looking up user", "error", err)
		}
	}
	if err == nil && user != nil && user.Groups != nil {
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: adding local groups", "num_local_groups", len(user.Groups), "local_groups", user.Groups)
		}
		allGroups = append(allGroups, user.Groups...)
	}

	// Retrieve policies
	var policies []string
	for _, groupName := range allGroups {
		entry, _, err := b.Group(ctx, req.Storage, groupName)
		if err != nil {
			if b.Logger().IsDebug() {
				b.Logger().Debug("auth/okta: error looking up group policies", "error", err)
			}
		}
		if err == nil && entry != nil && entry.Policies != nil {
			policies = append(policies, entry.Policies...)
		}
	}

	// Merge local Policies into Okta Policies
	if user != nil && user.Policies != nil {
		policies = append(policies, user.Policies...)
	}

	if len(policies) == 0 {
		errStr := "user is not a member of any authorized policy"
		if len(oktaResponse.Warnings) > 0 {
			errStr = fmt.Sprintf("%s; additionally, %s", errStr, oktaResponse.Warnings[0])
		}

		oktaResponse.Data["error"] = errStr
		return nil, oktaResponse, nil, nil
	}

	return policies, oktaResponse, allGroups, nil
}

func (b *backend) getOktaGroups(client *okta.Client, user *okta.User) ([]string, error) {
	rsp, err := client.Users.PopulateGroups(user)
	if err != nil {
		return nil, err
	}
	if rsp == nil {
		return nil, fmt.Errorf("okta auth method unexpected failure")
	}
	oktaGroups := make([]string, 0, len(user.Groups))
	for _, group := range user.Groups {
		oktaGroups = append(oktaGroups, group.Profile.Name)
	}
	if b.Logger().IsDebug() {
		b.Logger().Debug("auth/okta: Groups fetched from Okta", "num_groups", len(oktaGroups), "groups", oktaGroups)
	}
	return oktaGroups, nil
}

func (b *backend) completeMfa(ctx context.Context, client *okta.Client, mfaRequiredResult *authResult) (*authResult, error) {
	var factors []mfaFactor
	if err := mfaRequiredResult.Embedded.Decode("factors", &factors); err != nil {
		return nil, err
	}

	var handler verifyMfaHandler
	var factorId string
	for _, v := range factors {
		if v.Provider == "OKTA" && v.FactorType == "push" {
			// For Okta push, the only thing we need to do is hit the verify endpoint.
			handler = doOktaVerify
			factorId = v.Id
			break
		} else if v.Provider == "DUO" && v.FactorType == "web" {
			handler = b.verifyDuoPush
			factorId = v.Id
			break
		}
	}
	if handler == nil || factorId == "" {
		return nil, errors.New("Only Okta Push and Duo Push are supported for Okta MFA.")
	}

	stateToken := mfaRequiredResult.StateToken
	// Perform our part of the MFA flow.
	result, err := handler(client, stateToken, factorId)
	if err != nil {
		return nil, err
	}

	if result.FactorResult == "WAITING" {
		b.Logger().Debug("auth/okta: waiting for Okta MFA result")
	}
	for result.FactorResult == "WAITING" {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// okta has a documented rate limit of ~1tps, so we try to respect that here
			time.Sleep(500 * time.Millisecond)
			if result, err = doOktaVerify(client, stateToken, factorId); err != nil {
				return nil, err
			}
		}
	}

	switch result.FactorResult {
	case "REJECTED", "TIMEOUT", "CANCELLED", "ERROR":
		// These are known MFA failure cases and we want to communicate them individually.
		return nil, fmt.Errorf("multi-factor authentication failed: %s", result.FactorResult)
	case "":
		// Empty is success.
		b.Logger().Debug("auth/okta: multi-factor authentication succeeded")
		return result, nil
	default:
		// In the unknown case, let the caller decide based on status.
		b.Logger().Debug("auth/okta: unhandled factor result", "status", result.Status, "factor_result", result.FactorResult)
		return result, nil
	}
}

type mfaFactor struct {
	halResource `mapstructure:",squash"`
	Id          string
	FactorType  string
	Provider    string
}

// verifyMfaHandler verifies a single type of Okta MFA factor.
type verifyMfaHandler func(oktaClient *okta.Client, stateToken string, factorId string) (*authResult, error)

// doOktaVerify POSTs to the given factor's verification URL. The first call starts verification,
// and subsequent calls with the same state token get the current status of the verification.
func doOktaVerify(oktaClient *okta.Client, stateToken string, factorId string) (*authResult, error) {
	verifyUrl := fmt.Sprintf("authn/factors/%s/verify", factorId)
	payload := map[string]interface{}{
		"stateToken": stateToken,
	}
	verifyReq, err := oktaClient.NewRequest("POST", verifyUrl, payload)
	if err != nil {
		return nil, err
	}
	result := &authResult{}
	_, err = oktaClient.Do(verifyReq, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *backend) verifyDuoPush(oktaClient *okta.Client, stateToken string, factorId string) (*authResult, error) {
	result, err := doOktaVerify(oktaClient, stateToken, factorId)
	if err != nil {
		return nil, err
	}
	if result.Status != "MFA_CHALLENGE" {
		return nil, fmt.Errorf("expected auth status MFA_CHALLENGE, got %s", result.Status)
	}

	// Extract the challenge details.
	factor := mfaFactor{}
	if err = result.Embedded.Decode("factor", &factor); err != nil {
		return nil, err
	}
	var verification struct {
		halResource  `mapstructure:",squash"`
		Host         string
		Signature    string
		FactorResult string
	}
	if err = factor.Embedded.Decode("verification", &verification); err != nil {
		return nil, err
	}
	sigParts := strings.Split(verification.Signature, ":")
	txSig := sigParts[0]
	appSig := sigParts[1]

	// Kick off Duo verification exchange.
	parent := &url.URL{
		Scheme: oktaClient.BaseURL.Scheme,
		Host:   oktaClient.BaseURL.Host,
		Path:   "/signin/verify/duo/web",
	}
	queryVals := make(url.Values)
	queryVals.Add("parent", parent.String())
	queryVals.Add("tx", txSig)
	queryVals.Add("v", "2.6")
	reqUrl := &url.URL{
		Scheme:   "https",
		Host:     verification.Host,
		Path:     "/frame/web/v1/auth",
		RawQuery: queryVals.Encode(),
	}
	var promptUrl *url.URL
	// We don't actually need to load the UI, just to capture the prompt URL.
	captureRedirectClient := &http.Client{
		Transport: cleanhttp.DefaultTransport(),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			promptUrl = req.URL
			return http.ErrUseLastResponse
		},
	}
	resp, err := captureRedirectClient.PostForm(reqUrl.String(), queryVals)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// Send Duo verification push.
	sid := promptUrl.Query().Get("sid")
	promptUrl.RawQuery = ""
	promptVals := map[string]string{
		"sid": sid,
		// phone1 is the user's default device. phoneN is the Nth phone in the user's device list.
		"device": "phone1",
		"factor": "Duo Push",
	}
	var duoResponse struct {
		Response map[string]string `json:"response"`
		Stat     string            `json:"stat"`
	}
	if err = postForm(promptUrl.String(), promptVals, &duoResponse); err != nil {
		return nil, err
	}
	txid := duoResponse.Response["txid"]

	// Long-poll the status endpoint for status updates.
	statusUrl := &(*promptUrl)
	statusUrl.Path = "/frame/status"
	statusVals := map[string]string{
		"sid":  sid,
		"txid": txid,
	}
	for statusCode := ""; statusCode != "allow"; {
		if err = postForm(statusUrl.String(), statusVals, &duoResponse); err != nil {
			return nil, err
		}
		statusCode = duoResponse.Response["status_code"]
		status := duoResponse.Response["status"]
		switch statusCode {
		case "allow", "pushed":
			b.Logger().Debug(
				"auth/okta: Duo verification", "status_code", statusCode, "status", status)
		case "timeout":
			return nil, errors.New("Duo verification push timed out")
		default:
			return nil, fmt.Errorf("unknown Duo verification status %s: %s", statusCode, status)
		}
	}

	// Report the success back to Okta.
	completeVals := map[string]string{
		"id":           factorId,
		"stateToken":   stateToken,
		"sig_response": fmt.Sprintf("%s:%s", duoResponse.Response["cookie"], appSig),
	}
	completeLink, err := verification.Links.Link("complete")
	if err != nil {
		return nil, err
	}
	if err = postForm(completeLink.Href, completeVals, nil); err != nil {
		return nil, err
	}
	b.Logger().Debug("auth/okta: Duo verification: result posted back to Okta")

	// Verify success with Okta.
	return doOktaVerify(oktaClient, stateToken, factorId)
}

func postForm(uri string, values map[string]string, dest interface{}) error {
	urlValues := make(url.Values)
	for k, v := range values {
		urlValues.Set(k, v)
	}
	resp, err := http.PostForm(uri, urlValues)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if dest == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(dest)
}

const backendHelp = `
The Okta credential provider allows authentication querying,
checking username and password, and associating policies.  If an api token is configure
groups are pulled down from Okta.

Configuration of the connection is done through the "config" and "policies"
endpoints by a user with root access. Authentication is then done
by suppying the two fields for "login".
`
