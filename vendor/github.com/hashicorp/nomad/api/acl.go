// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ACLPolicies is used to query the ACL Policy endpoints.
type ACLPolicies struct {
	client *Client
}

// ACLPolicies returns a new handle on the ACL policies.
func (c *Client) ACLPolicies() *ACLPolicies {
	return &ACLPolicies{client: c}
}

// List is used to dump all of the policies.
func (a *ACLPolicies) List(q *QueryOptions) ([]*ACLPolicyListStub, *QueryMeta, error) {
	var resp []*ACLPolicyListStub
	qm, err := a.client.query("/v1/acl/policies", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

// Upsert is used to create or update a policy
func (a *ACLPolicies) Upsert(policy *ACLPolicy, q *WriteOptions) (*WriteMeta, error) {
	if policy == nil || policy.Name == "" {
		return nil, errors.New("missing policy name")
	}
	wm, err := a.client.put("/v1/acl/policy/"+policy.Name, policy, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Delete is used to delete a policy
func (a *ACLPolicies) Delete(policyName string, q *WriteOptions) (*WriteMeta, error) {
	if policyName == "" {
		return nil, errors.New("missing policy name")
	}
	wm, err := a.client.delete("/v1/acl/policy/"+policyName, nil, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Info is used to query a specific policy
func (a *ACLPolicies) Info(policyName string, q *QueryOptions) (*ACLPolicy, *QueryMeta, error) {
	if policyName == "" {
		return nil, nil, errors.New("missing policy name")
	}
	var resp ACLPolicy
	wm, err := a.client.query("/v1/acl/policy/"+policyName, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// ACLTokens is used to query the ACL token endpoints.
type ACLTokens struct {
	client *Client
}

// ACLTokens returns a new handle on the ACL tokens.
func (c *Client) ACLTokens() *ACLTokens {
	return &ACLTokens{client: c}
}

// Bootstrap is used to get the initial bootstrap token
//
// See BootstrapOpts to set ACL bootstrapping options.
func (a *ACLTokens) Bootstrap(q *WriteOptions) (*ACLToken, *WriteMeta, error) {
	var resp ACLToken
	wm, err := a.client.put("/v1/acl/bootstrap", nil, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// BootstrapOpts is used to get the initial bootstrap token or pass in the one that was provided in the API
func (a *ACLTokens) BootstrapOpts(btoken string, q *WriteOptions) (*ACLToken, *WriteMeta, error) {
	if q == nil {
		q = &WriteOptions{}
	}
	req := &BootstrapRequest{
		BootstrapSecret: btoken,
	}

	var resp ACLToken
	wm, err := a.client.put("/v1/acl/bootstrap", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// List is used to dump all of the tokens.
func (a *ACLTokens) List(q *QueryOptions) ([]*ACLTokenListStub, *QueryMeta, error) {
	var resp []*ACLTokenListStub
	qm, err := a.client.query("/v1/acl/tokens", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

// Create is used to create a token
func (a *ACLTokens) Create(token *ACLToken, q *WriteOptions) (*ACLToken, *WriteMeta, error) {
	if token.AccessorID != "" {
		return nil, nil, errors.New("cannot specify Accessor ID")
	}
	var resp ACLToken
	wm, err := a.client.put("/v1/acl/token", token, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Update is used to update an existing token
func (a *ACLTokens) Update(token *ACLToken, q *WriteOptions) (*ACLToken, *WriteMeta, error) {
	if token.AccessorID == "" {
		return nil, nil, errors.New("missing accessor ID")
	}
	var resp ACLToken
	wm, err := a.client.put("/v1/acl/token/"+token.AccessorID,
		token, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Delete is used to delete a token
func (a *ACLTokens) Delete(accessorID string, q *WriteOptions) (*WriteMeta, error) {
	if accessorID == "" {
		return nil, errors.New("missing accessor ID")
	}
	wm, err := a.client.delete("/v1/acl/token/"+accessorID, nil, nil, q)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Info is used to query a token
func (a *ACLTokens) Info(accessorID string, q *QueryOptions) (*ACLToken, *QueryMeta, error) {
	if accessorID == "" {
		return nil, nil, errors.New("missing accessor ID")
	}
	var resp ACLToken
	wm, err := a.client.query("/v1/acl/token/"+accessorID, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Self is used to query our own token
func (a *ACLTokens) Self(q *QueryOptions) (*ACLToken, *QueryMeta, error) {
	var resp ACLToken
	wm, err := a.client.query("/v1/acl/token/self", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// UpsertOneTimeToken is used to create a one-time token
func (a *ACLTokens) UpsertOneTimeToken(q *WriteOptions) (*OneTimeToken, *WriteMeta, error) {
	var resp *OneTimeTokenUpsertResponse
	wm, err := a.client.put("/v1/acl/token/onetime", nil, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, errors.New("no one-time token returned")
	}
	return resp.OneTimeToken, wm, nil
}

// ExchangeOneTimeToken is used to create a one-time token
func (a *ACLTokens) ExchangeOneTimeToken(secret string, q *WriteOptions) (*ACLToken, *WriteMeta, error) {
	if secret == "" {
		return nil, nil, errors.New("missing secret ID")
	}
	req := &OneTimeTokenExchangeRequest{OneTimeSecretID: secret}
	var resp *OneTimeTokenExchangeResponse
	wm, err := a.client.put("/v1/acl/token/onetime/exchange", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, errors.New("no ACL token returned")
	}
	return resp.Token, wm, nil
}

var (
	// errMissingACLRoleID is the generic errors to use when a call is missing
	// the required ACL Role ID parameter.
	errMissingACLRoleID = errors.New("missing ACL role ID")

	// errMissingACLAuthMethodName is the generic error to use when a call is
	// missing the required ACL auth-method name parameter.
	errMissingACLAuthMethodName = errors.New("missing ACL auth-method name")

	// errMissingACLBindingRuleID is the generic error to use when a call is
	// missing the required ACL binding rule ID parameter.
	errMissingACLBindingRuleID = errors.New("missing ACL binding rule ID")
)

// ACLRoles is used to query the ACL Role endpoints.
type ACLRoles struct {
	client *Client
}

// ACLRoles returns a new handle on the ACL roles API client.
func (c *Client) ACLRoles() *ACLRoles {
	return &ACLRoles{client: c}
}

// List is used to detail all the ACL roles currently stored within state.
func (a *ACLRoles) List(q *QueryOptions) ([]*ACLRoleListStub, *QueryMeta, error) {
	var resp []*ACLRoleListStub
	qm, err := a.client.query("/v1/acl/roles", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

// Create is used to create an ACL role.
func (a *ACLRoles) Create(role *ACLRole, w *WriteOptions) (*ACLRole, *WriteMeta, error) {
	if role.ID != "" {
		return nil, nil, errors.New("cannot specify ACL role ID")
	}
	var resp ACLRole
	wm, err := a.client.put("/v1/acl/role", role, &resp, w)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Update is used to update an existing ACL role.
func (a *ACLRoles) Update(role *ACLRole, w *WriteOptions) (*ACLRole, *WriteMeta, error) {
	if role.ID == "" {
		return nil, nil, errMissingACLRoleID
	}
	var resp ACLRole
	wm, err := a.client.put("/v1/acl/role/"+role.ID, role, &resp, w)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Delete is used to delete an ACL role.
func (a *ACLRoles) Delete(roleID string, w *WriteOptions) (*WriteMeta, error) {
	if roleID == "" {
		return nil, errMissingACLRoleID
	}
	wm, err := a.client.delete("/v1/acl/role/"+roleID, nil, nil, w)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Get is used to look up an ACL role.
func (a *ACLRoles) Get(roleID string, q *QueryOptions) (*ACLRole, *QueryMeta, error) {
	if roleID == "" {
		return nil, nil, errMissingACLRoleID
	}
	var resp ACLRole
	qm, err := a.client.query("/v1/acl/role/"+roleID, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// GetByName is used to look up an ACL role using its name.
func (a *ACLRoles) GetByName(roleName string, q *QueryOptions) (*ACLRole, *QueryMeta, error) {
	if roleName == "" {
		return nil, nil, errors.New("missing ACL role name")
	}
	var resp ACLRole
	qm, err := a.client.query("/v1/acl/role/name/"+roleName, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// ACLAuthMethods is used to query the ACL auth-methods endpoints.
type ACLAuthMethods struct {
	client *Client
}

// ACLAuthMethods returns a new handle on the ACL auth-methods API client.
func (c *Client) ACLAuthMethods() *ACLAuthMethods {
	return &ACLAuthMethods{client: c}
}

// List is used to detail all the ACL auth-methods currently stored within
// state.
func (a *ACLAuthMethods) List(q *QueryOptions) ([]*ACLAuthMethodListStub, *QueryMeta, error) {
	var resp []*ACLAuthMethodListStub
	qm, err := a.client.query("/v1/acl/auth-methods", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

// Create is used to create an ACL auth-method.
func (a *ACLAuthMethods) Create(authMethod *ACLAuthMethod, w *WriteOptions) (*ACLAuthMethod, *WriteMeta, error) {
	if authMethod.Name == "" {
		return nil, nil, errMissingACLAuthMethodName
	}
	var resp ACLAuthMethod
	wm, err := a.client.put("/v1/acl/auth-method", authMethod, &resp, w)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Update is used to update an existing ACL auth-method.
func (a *ACLAuthMethods) Update(authMethod *ACLAuthMethod, w *WriteOptions) (*ACLAuthMethod, *WriteMeta, error) {
	if authMethod.Name == "" {
		return nil, nil, errMissingACLAuthMethodName
	}
	var resp ACLAuthMethod
	wm, err := a.client.put("/v1/acl/auth-method/"+authMethod.Name, authMethod, &resp, w)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Delete is used to delete an ACL auth-method.
func (a *ACLAuthMethods) Delete(authMethodName string, w *WriteOptions) (*WriteMeta, error) {
	if authMethodName == "" {
		return nil, errMissingACLAuthMethodName
	}
	wm, err := a.client.delete("/v1/acl/auth-method/"+authMethodName, nil, nil, w)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Get is used to look up an ACL auth-method.
func (a *ACLAuthMethods) Get(authMethodName string, q *QueryOptions) (*ACLAuthMethod, *QueryMeta, error) {
	if authMethodName == "" {
		return nil, nil, errMissingACLAuthMethodName
	}
	var resp ACLAuthMethod
	qm, err := a.client.query("/v1/acl/auth-method/"+authMethodName, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// ACLBindingRules is used to query the ACL auth-methods endpoints.
type ACLBindingRules struct {
	client *Client
}

// ACLBindingRules returns a new handle on the ACL auth-methods API client.
func (c *Client) ACLBindingRules() *ACLBindingRules {
	return &ACLBindingRules{client: c}
}

// List is used to detail all the ACL binding rules currently stored within
// state.
func (a *ACLBindingRules) List(q *QueryOptions) ([]*ACLBindingRuleListStub, *QueryMeta, error) {
	var resp []*ACLBindingRuleListStub
	qm, err := a.client.query("/v1/acl/binding-rules", &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return resp, qm, nil
}

// Create is used to create an ACL binding rule.
func (a *ACLBindingRules) Create(bindingRule *ACLBindingRule, w *WriteOptions) (*ACLBindingRule, *WriteMeta, error) {
	var resp ACLBindingRule
	wm, err := a.client.put("/v1/acl/binding-rule", bindingRule, &resp, w)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Update is used to update an existing ACL binding rule.
func (a *ACLBindingRules) Update(bindingRule *ACLBindingRule, w *WriteOptions) (*ACLBindingRule, *WriteMeta, error) {
	if bindingRule.ID == "" {
		return nil, nil, errMissingACLBindingRuleID
	}
	var resp ACLBindingRule
	wm, err := a.client.put("/v1/acl/binding-rule/"+bindingRule.ID, bindingRule, &resp, w)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Delete is used to delete an ACL binding rule.
func (a *ACLBindingRules) Delete(bindingRuleID string, w *WriteOptions) (*WriteMeta, error) {
	if bindingRuleID == "" {
		return nil, errMissingACLBindingRuleID
	}
	wm, err := a.client.delete("/v1/acl/binding-rule/"+bindingRuleID, nil, nil, w)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

// Get is used to look up an ACL binding rule.
func (a *ACLBindingRules) Get(bindingRuleID string, q *QueryOptions) (*ACLBindingRule, *QueryMeta, error) {
	if bindingRuleID == "" {
		return nil, nil, errMissingACLBindingRuleID
	}
	var resp ACLBindingRule
	qm, err := a.client.query("/v1/acl/binding-rule/"+bindingRuleID, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, qm, nil
}

// ACLOIDC is used to query the ACL OIDC endpoints.
//
// Deprecated: ACLOIDC is deprecated, use ACLAuth instead.
type ACLOIDC struct {
	client *Client
	ACLAuth
}

// ACLOIDC returns a new handle on the ACL auth-methods API client.
//
// Deprecated: c.ACLOIDC() is deprecated, use c.ACLAuth() instead.
func (c *Client) ACLOIDC() *ACLOIDC {
	return &ACLOIDC{client: c}
}

// ACLAuth is used to query the ACL auth endpoints.
type ACLAuth struct {
	client *Client
}

// ACLAuth returns a new handle on the ACL auth-methods API client.
func (c *Client) ACLAuth() *ACLAuth {
	return &ACLAuth{client: c}
}

// GetAuthURL generates the OIDC provider authentication URL. This URL should
// be visited in order to sign in to the provider.
func (a *ACLAuth) GetAuthURL(req *ACLOIDCAuthURLRequest, q *WriteOptions) (*ACLOIDCAuthURLResponse, *WriteMeta, error) {
	var resp ACLOIDCAuthURLResponse
	wm, err := a.client.put("/v1/acl/oidc/auth-url", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// CompleteAuth exchanges the OIDC provider token for a Nomad token with the
// appropriate claims attached.
func (a *ACLAuth) CompleteAuth(req *ACLOIDCCompleteAuthRequest, q *WriteOptions) (*ACLToken, *WriteMeta, error) {
	var resp ACLToken
	wm, err := a.client.put("/v1/acl/oidc/complete-auth", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// Login exchanges the third party token for a Nomad token with the appropriate
// claims attached.
func (a *ACLAuth) Login(req *ACLLoginRequest, q *WriteOptions) (*ACLToken, *WriteMeta, error) {
	var resp ACLToken
	wm, err := a.client.put("/v1/acl/login", req, &resp, q)
	if err != nil {
		return nil, nil, err
	}
	return &resp, wm, nil
}

// ACLPolicyListStub is used to for listing ACL policies
type ACLPolicyListStub struct {
	Name        string
	Description string
	CreateIndex uint64
	ModifyIndex uint64
}

// ACLPolicy is used to represent an ACL policy
type ACLPolicy struct {
	Name        string
	Description string
	Rules       string
	JobACL      *JobACL

	CreateIndex uint64
	ModifyIndex uint64
}

// JobACL represents an ACL policy's attachment to a job, group, or task.
type JobACL struct {
	Namespace string
	JobID     string
	Group     string
	Task      string
}

// ACLToken represents a client token which is used to Authenticate
type ACLToken struct {
	AccessorID string
	SecretID   string
	Name       string
	Type       string
	Policies   []string

	// Roles represents the ACL roles that this token is tied to. The token
	// will inherit the permissions of all policies detailed within the role.
	Roles []*ACLTokenRoleLink

	Global     bool
	CreateTime time.Time

	// ExpirationTime represents the point after which a token should be
	// considered revoked and is eligible for destruction. The zero value of
	// time.Time does not respect json omitempty directives, so we must use a
	// pointer.
	ExpirationTime *time.Time `json:",omitempty"`

	// ExpirationTTL is a convenience field for helping set ExpirationTime to a
	// value of CreateTime+ExpirationTTL. This can only be set during token
	// creation. This is a string version of a time.Duration like "2m".
	ExpirationTTL time.Duration `json:",omitempty"`

	CreateIndex uint64
	ModifyIndex uint64
}

// ACLTokenRoleLink is used to link an ACL token to an ACL role. The ACL token
// can therefore inherit all the ACL policy permissions that the ACL role
// contains.
type ACLTokenRoleLink struct {

	// ID is the ACLRole.ID UUID. This field is immutable and represents the
	// absolute truth for the link.
	ID string

	// Name is the human friendly identifier for the ACL role and is a
	// convenience field for operators.
	Name string
}

// MarshalJSON implements the json.Marshaler interface and allows
// ACLToken.ExpirationTTL to be marshaled correctly.
func (a *ACLToken) MarshalJSON() ([]byte, error) {
	type Alias ACLToken
	exported := &struct {
		ExpirationTTL string
		*Alias
	}{
		ExpirationTTL: a.ExpirationTTL.String(),
		Alias:         (*Alias)(a),
	}
	if a.ExpirationTTL == 0 {
		exported.ExpirationTTL = ""
	}
	return json.Marshal(exported)
}

// UnmarshalJSON implements the json.Unmarshaler interface and allows
// ACLToken.ExpirationTTL to be unmarshalled correctly.
func (a *ACLToken) UnmarshalJSON(data []byte) (err error) {
	type Alias ACLToken
	aux := &struct {
		ExpirationTTL any
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.ExpirationTTL != nil {
		switch v := aux.ExpirationTTL.(type) {
		case string:
			if v != "" {
				if a.ExpirationTTL, err = time.ParseDuration(v); err != nil {
					return err
				}
			}
		case float64:
			a.ExpirationTTL = time.Duration(v)
		}

	}
	return nil
}

type ACLTokenListStub struct {
	AccessorID string
	Name       string
	Type       string
	Policies   []string
	Roles      []*ACLTokenRoleLink
	Global     bool
	CreateTime time.Time

	// ExpirationTime represents the point after which a token should be
	// considered revoked and is eligible for destruction. A nil value
	// indicates no expiration has been set on the token.
	ExpirationTime *time.Time `json:",omitempty"`

	CreateIndex uint64
	ModifyIndex uint64
}

type OneTimeToken struct {
	OneTimeSecretID string
	AccessorID      string
	ExpiresAt       time.Time
	CreateIndex     uint64
	ModifyIndex     uint64
}

type OneTimeTokenUpsertResponse struct {
	OneTimeToken *OneTimeToken
}

type OneTimeTokenExchangeRequest struct {
	OneTimeSecretID string
}

type OneTimeTokenExchangeResponse struct {
	Token *ACLToken
}

// BootstrapRequest is used for when operators provide an ACL Bootstrap Token
type BootstrapRequest struct {
	BootstrapSecret string
}

// ACLRole is an abstraction for the ACL system which allows the grouping of
// ACL policies into a single object. ACL tokens can be created and linked to
// a role; the token then inherits all the permissions granted by the policies.
type ACLRole struct {

	// ID is an internally generated UUID for this role and is controlled by
	// Nomad. It can be used after role creation to update the existing role.
	ID string

	// Name is unique across the entire set of federated clusters and is
	// supplied by the operator on role creation. The name can be modified by
	// updating the role and including the Nomad generated ID. This update will
	// not affect tokens created and linked to this role. This is a required
	// field.
	Name string

	// Description is a human-readable, operator set description that can
	// provide additional context about the role. This is an optional field.
	Description string

	// Policies is an array of ACL policy links. Although currently policies
	// can only be linked using their name, in the future we will want to add
	// IDs also and thus allow operators to specify either a name, an ID, or
	// both. At least one entry is required.
	Policies []*ACLRolePolicyLink

	CreateIndex uint64
	ModifyIndex uint64
}

// ACLRolePolicyLink is used to link a policy to an ACL role. We use a struct
// rather than a list of strings as in the future we will want to add IDs to
// policies and then link via these.
type ACLRolePolicyLink struct {

	// Name is the ACLPolicy.Name value which will be linked to the ACL role.
	Name string
}

// ACLRoleListStub is the stub object returned when performing a listing of ACL
// roles. While it might not currently be different to the full response
// object, it allows us to future-proof the RPC in the event the ACLRole object
// grows over time.
type ACLRoleListStub struct {

	// ID is an internally generated UUID for this role and is controlled by
	// Nomad.
	ID string

	// Name is unique across the entire set of federated clusters and is
	// supplied by the operator on role creation. The name can be modified by
	// updating the role and including the Nomad generated ID. This update will
	// not affect tokens created and linked to this role. This is a required
	// field.
	Name string

	// Description is a human-readable, operator set description that can
	// provide additional context about the role. This is an operational field.
	Description string

	// Policies is an array of ACL policy links. Although currently policies
	// can only be linked using their name, in the future we will want to add
	// IDs also and thus allow operators to specify either a name, an ID, or
	// both.
	Policies []*ACLRolePolicyLink

	CreateIndex uint64
	ModifyIndex uint64
}

// ACLAuthMethod is used to capture the properties of an authentication method
// used for single sing-on.
type ACLAuthMethod struct {

	// Name is the identifier for this auth-method and is a required parameter.
	Name string

	// Type is the SSO identifier this auth-method is. Nomad currently only
	// supports "oidc" and the API contains ACLAuthMethodTypeOIDC for
	// convenience.
	Type string

	// Defines whether the auth-method creates a local or global token when
	// performing SSO login. This should be set to either "local" or "global"
	// and the API contains ACLAuthMethodTokenLocalityLocal and
	// ACLAuthMethodTokenLocalityGlobal for convenience.
	TokenLocality string

	// TokenNameFormat defines the HIL template to use when building the token name
	TokenNameFormat string

	// MaxTokenTTL is the maximum life of a token created by this method.
	MaxTokenTTL time.Duration

	// Default identifies whether this is the default auth-method to use when
	// attempting to login without specifying an auth-method name to use.
	Default bool

	// Config contains the detailed configuration which is specific to the
	// auth-method.
	Config *ACLAuthMethodConfig

	CreateTime  time.Time
	ModifyTime  time.Time
	CreateIndex uint64
	ModifyIndex uint64
}

// MarshalJSON implements the json.Marshaler interface and allows
// ACLAuthMethod.MaxTokenTTL to be marshaled correctly.
func (m *ACLAuthMethod) MarshalJSON() ([]byte, error) {
	type Alias ACLAuthMethod
	exported := &struct {
		MaxTokenTTL string
		*Alias
	}{
		MaxTokenTTL: m.MaxTokenTTL.String(),
		Alias:       (*Alias)(m),
	}
	if m.MaxTokenTTL == 0 {
		exported.MaxTokenTTL = ""
	}
	return json.Marshal(exported)
}

// UnmarshalJSON implements the json.Unmarshaler interface and allows
// ACLAuthMethod.MaxTokenTTL to be unmarshalled correctly.
func (m *ACLAuthMethod) UnmarshalJSON(data []byte) error {
	type Alias ACLAuthMethod
	aux := &struct {
		MaxTokenTTL string
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	if aux.MaxTokenTTL != "" {
		if m.MaxTokenTTL, err = time.ParseDuration(aux.MaxTokenTTL); err != nil {
			return err
		}
	}
	return nil
}

// ACLAuthMethodConfig is used to store configuration of an auth method.
type ACLAuthMethodConfig struct {
	// A list of PEM-encoded public keys to use to authenticate signatures
	// locally
	JWTValidationPubKeys []string
	// JSON Web Key Sets url for authenticating signatures
	JWKSURL string
	// The OIDC Discovery URL, without any .well-known component (base path)
	OIDCDiscoveryURL string
	// The OAuth Client ID configured with the OIDC provider
	OIDCClientID string
	// The OAuth Client Secret configured with the OIDC provider
	OIDCClientSecret string
	// Disable claims from the OIDC UserInfo endpoint
	OIDCDisableUserInfo bool
	// List of OIDC scopes
	OIDCScopes []string
	// List of auth claims that are valid for login
	BoundAudiences []string
	// The value against which to match the iss claim in a JWT
	BoundIssuer []string
	// A list of allowed values for redirect_uri
	AllowedRedirectURIs []string
	// PEM encoded CA certs for use by the TLS client used to talk with the
	// OIDC Discovery URL.
	DiscoveryCaPem []string
	// PEM encoded CA cert for use by the TLS client used to talk with the JWKS
	// URL
	JWKSCACert string
	// A list of supported signing algorithms
	SigningAlgs []string
	// Duration in seconds of leeway when validating expiration of a token to
	// account for clock skew
	ExpirationLeeway time.Duration
	// Duration in seconds of leeway when validating not before values of a
	// token to account for clock skew.
	NotBeforeLeeway time.Duration
	// Duration in seconds of leeway when validating all claims to account for
	// clock skew.
	ClockSkewLeeway time.Duration
	// Mappings of claims (key) that will be copied to a metadata field
	// (value).
	ClaimMappings     map[string]string
	ListClaimMappings map[string]string
}

// MarshalJSON implements the json.Marshaler interface and allows
// time.Duration fields to be marshaled correctly.
func (c *ACLAuthMethodConfig) MarshalJSON() ([]byte, error) {
	type Alias ACLAuthMethodConfig
	exported := &struct {
		ExpirationLeeway string
		NotBeforeLeeway  string
		ClockSkewLeeway  string
		*Alias
	}{
		ExpirationLeeway: c.ExpirationLeeway.String(),
		NotBeforeLeeway:  c.NotBeforeLeeway.String(),
		ClockSkewLeeway:  c.ClockSkewLeeway.String(),
		Alias:            (*Alias)(c),
	}
	if c.ExpirationLeeway == 0 {
		exported.ExpirationLeeway = ""
	}
	if c.NotBeforeLeeway == 0 {
		exported.NotBeforeLeeway = ""
	}
	if c.ClockSkewLeeway == 0 {
		exported.ClockSkewLeeway = ""
	}
	return json.Marshal(exported)
}

// UnmarshalJSON implements the json.Unmarshaler interface and allows
// time.Duration fields to be unmarshalled correctly.
func (c *ACLAuthMethodConfig) UnmarshalJSON(data []byte) error {
	type Alias ACLAuthMethodConfig
	aux := &struct {
		ExpirationLeeway any
		NotBeforeLeeway  any
		ClockSkewLeeway  any
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	if aux.ExpirationLeeway != nil {
		switch v := aux.ExpirationLeeway.(type) {
		case string:
			if v != "" {
				if c.ExpirationLeeway, err = time.ParseDuration(v); err != nil {
					return err
				}
			}
		case float64:
			c.ExpirationLeeway = time.Duration(v)
		default:
			return fmt.Errorf("unexpected ExpirationLeeway type: %v", v)
		}
	}
	if aux.NotBeforeLeeway != nil {
		switch v := aux.NotBeforeLeeway.(type) {
		case string:
			if v != "" {
				if c.NotBeforeLeeway, err = time.ParseDuration(v); err != nil {
					return err
				}
			}
		case float64:
			c.NotBeforeLeeway = time.Duration(v)
		default:
			return fmt.Errorf("unexpected NotBeforeLeeway type: %v", v)
		}
	}
	if aux.ClockSkewLeeway != nil {
		switch v := aux.ClockSkewLeeway.(type) {
		case string:
			if v != "" {
				if c.ClockSkewLeeway, err = time.ParseDuration(v); err != nil {
					return err
				}
			}
		case float64:
			c.ClockSkewLeeway = time.Duration(v)
		default:
			return fmt.Errorf("unexpected ClockSkewLeeway type: %v", v)
		}
	}
	return nil
}

// ACLAuthMethodListStub is the stub object returned when performing a listing
// of ACL auth-methods. It is intentionally minimal due to the unauthenticated
// nature of the list endpoint.
type ACLAuthMethodListStub struct {
	Name    string
	Type    string
	Default bool

	CreateIndex uint64
	ModifyIndex uint64
}

const (
	// ACLAuthMethodTokenLocalityLocal is the ACLAuthMethod.TokenLocality that
	// will generate ACL tokens which can only be used on the local cluster the
	// request was made.
	ACLAuthMethodTokenLocalityLocal = "local"

	// ACLAuthMethodTokenLocalityGlobal is the ACLAuthMethod.TokenLocality that
	// will generate ACL tokens which can be used on all federated clusters.
	ACLAuthMethodTokenLocalityGlobal = "global"

	// ACLAuthMethodTypeOIDC the ACLAuthMethod.Type and represents an
	// auth-method which uses the OIDC protocol.
	ACLAuthMethodTypeOIDC = "OIDC"

	// ACLAuthMethodTypeJWT the ACLAuthMethod.Type and represents an auth-method
	// which uses the JWT type.
	ACLAuthMethodTypeJWT = "JWT"
)

// ACLBindingRule contains a direct relation to an ACLAuthMethod and represents
// a rule to apply when logging in via the named AuthMethod. This allows the
// transformation of OIDC provider claims, to Nomad based ACL concepts such as
// ACL Roles and Policies.
type ACLBindingRule struct {

	// ID is an internally generated UUID for this rule and is controlled by
	// Nomad.
	ID string

	// Description is a human-readable, operator set description that can
	// provide additional context about the binding rule. This is an
	// operational field.
	Description string

	// AuthMethod is the name of the auth method for which this rule applies
	// to. This is required and the method must exist within state before the
	// cluster administrator can create the rule.
	AuthMethod string

	// Selector is an expression that matches against verified identity
	// attributes returned from the auth method during login. This is optional
	// and when not set, provides a catch-all rule.
	Selector string

	// BindType adjusts how this binding rule is applied at login time. The
	// valid values are ACLBindingRuleBindTypeRole,
	// ACLBindingRuleBindTypePolicy, and ACLBindingRuleBindTypeManagement.
	BindType string

	// BindName is the target of the binding. Can be lightly templated using
	// HIL ${foo} syntax from available field names. How it is used depends
	// upon the BindType.
	BindName string

	CreateTime  time.Time
	ModifyTime  time.Time
	CreateIndex uint64
	ModifyIndex uint64
}

const (
	// ACLBindingRuleBindTypeRole is the ACL binding rule bind type that only
	// allows the binding rule to function if a role exists at login-time. The
	// role will be specified within the ACLBindingRule.BindName parameter, and
	// will identify whether this is an ID or Name.
	ACLBindingRuleBindTypeRole = "role"

	// ACLBindingRuleBindTypePolicy is the ACL binding rule bind type that
	// assigns a policy to the generate ACL token. The role will be specified
	// within the ACLBindingRule.BindName parameter, and will be the policy
	// name.
	ACLBindingRuleBindTypePolicy = "policy"

	// ACLBindingRuleBindTypeManagement is the ACL binding rule bind type that
	// will generate management ACL tokens when matched.
	ACLBindingRuleBindTypeManagement = "management"
)

// ACLBindingRuleListStub is the stub object returned when performing a listing
// of ACL binding rules.
type ACLBindingRuleListStub struct {

	// ID is an internally generated UUID for this role and is controlled by
	// Nomad.
	ID string

	// Description is a human-readable, operator set description that can
	// provide additional context about the binding role. This is an
	// operational field.
	Description string

	// AuthMethod is the name of the auth method for which this rule applies
	// to. This is required and the method must exist within state before the
	// cluster administrator can create the rule.
	AuthMethod string

	CreateIndex uint64
	ModifyIndex uint64
}

// ACLOIDCAuthURLRequest is the request to make when starting the OIDC
// authentication login flow.
type ACLOIDCAuthURLRequest struct {

	// AuthMethodName is the OIDC auth-method to use. This is a required
	// parameter.
	AuthMethodName string

	// RedirectURI is the URL that authorization should redirect to. This is a
	// required parameter.
	RedirectURI string

	// ClientNonce is a randomly generated string to prevent replay attacks. It
	// is up to the client to generate this and Go integrations should use the
	// oidc.NewID function within the hashicorp/cap library.
	ClientNonce string
}

// ACLOIDCAuthURLResponse is the response when starting the OIDC authentication
// login flow.
type ACLOIDCAuthURLResponse struct {

	// AuthURL is URL to begin authorization and is where the user logging in
	// should go.
	AuthURL string
}

// ACLOIDCCompleteAuthRequest is the request object to begin completing the
// OIDC auth cycle after receiving the callback from the OIDC provider.
type ACLOIDCCompleteAuthRequest struct {

	// AuthMethodName is the name of the auth method being used to login via
	// OIDC. This will match AuthUrlArgs.AuthMethodName. This is a required
	// parameter.
	AuthMethodName string

	// ClientNonce, State, and Code are provided from the parameters given to
	// the redirect URL. These are all required parameters.
	ClientNonce string
	State       string
	Code        string

	// RedirectURI is the URL that authorization should redirect to. This is a
	// required parameter.
	RedirectURI string
}

// ACLLoginRequest is the request object to begin auth with an external bearer
// token provider.
type ACLLoginRequest struct {
	// AuthMethodName is the name of the auth method being used to login. This
	// is a required parameter.
	AuthMethodName string
	// LoginToken is the token used to login. This is a required parameter.
	LoginToken string
}
