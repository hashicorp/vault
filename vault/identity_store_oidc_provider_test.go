package vault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2"
)

func TestOIDC_Path_OIDC_Authorize(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := new(logical.InmemStorage)

	// Create a key
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data:      map[string]interface{}{},
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Create an entity
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name": "test-entity",
		},
	})
	expectSuccess(t, resp, err)
	assert.NotNil(t, resp.Data["id"])
	entityID := resp.Data["id"].(string)

	// Create a group
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":              "test-group",
			"member_entity_ids": []string{entityID},
		},
	})
	expectSuccess(t, resp, err)
	assert.NotNil(t, resp.Data["id"])
	groupID := resp.Data["id"].(string)

	type args struct {
		entityID          string
		client            client
		provider          provider
		assignment        assignment
		authorizeRequest  *logical.Request
		tokenCreationTime func() time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "invalid authorize request with provider not found",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/non-existent-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with empty scopes",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with missing openid scope",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "groups email profile",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with missing response_type",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with unsupported response_type",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "id_token",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthUnsupportedResponseType,
		},
		{
			name: "invalid authorize request with client_id not found",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "non-existent-client-id",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthInvalidClientID,
		},
		{
			name: "invalid authorize request with client_id not allowed by provider",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				provider: provider{
					AllowedClientIDs: []string{"not-client-id"},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthUnauthorizedClient,
		},
		{
			name: "invalid authorize request with missing redirect_uri",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with redirect_uri not allowed by client",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://not.redirect.uri:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthInvalidRedirectURI,
		},
		{
			name: "invalid authorize request with missing state",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with missing nonce",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "",
					},
				},
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with request parameter provided",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
						"request":       "header.payload.signature",
					},
				},
			},
			wantErr: ErrAuthRequestNotSupported,
		},
		{
			name: "invalid authorize request with request_uri parameter provided",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
						"request_uri":   "https://client.example.org/request.jwt",
					},
				},
			},
			wantErr: ErrAuthRequestURINotSupported,
		},
		{
			name: "invalid authorize request with identity entity ID not found",
			args: args{
				entityID: "non-existent-entity",
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthAccessDenied,
		},
		{
			name: "invalid authorize request with entity not found in client assignment",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{"not-entity-id"},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthAccessDenied,
		},
		{
			name: "invalid authorize request with group not found in client assignment",
			args: args{
				entityID: entityID,
				assignment: assignment{
					GroupIDs: []string{"not-group-id"},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
			wantErr: ErrAuthAccessDenied,
		},
		{
			name: "invalid authorize request with negative max_age",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
						"max_age":       "-1",
					},
				},
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "active re-authentication required with token creation time exceeding max_age requirement",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
						"max_age":       "30",
					},
				},
				tokenCreationTime: func() time.Time {
					return time.Now().Add(-time.Minute)
				},
			},
			wantErr: ErrAuthMaxAgeReAuthenticate,
		},
		{
			name: "valid authorize request with token creation time within max_age requirement",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
						"max_age":       "30",
					},
				},
				tokenCreationTime: func() time.Time {
					return time.Now()
				},
			},
		},
		{
			name: "valid authorize request using update operation (HTTP PUT/POST)",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
		},
		{
			name: "valid authorize request using read operation (HTTP GET)",
			args: args{
				entityID: entityID,
				assignment: assignment{
					EntityIDs: []string{entityID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.ReadOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
		},
		{
			name: "valid authorize request using client assignment with group membership",
			args: args{
				entityID: entityID,
				assignment: assignment{
					GroupIDs: []string{groupID},
				},
				client: client{
					RedirectURIs: []string{"https://localhost:8251/callback"},
					Assignments:  []string{"test-assignment"},
					Key:          "test-key",
				},
				authorizeRequest: &logical.Request{
					Path:      "oidc/provider/test-provider/authorize",
					Operation: logical.UpdateOperation,
					Data: map[string]interface{}{
						"client_id":     "",
						"scope":         "openid",
						"redirect_uri":  "https://localhost:8251/callback",
						"response_type": "code",
						"state":         "abcdefg",
						"nonce":         "hijklmn",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a token entry and associate with the authorize request
			creationTime := time.Now()
			if tt.args.tokenCreationTime != nil {
				creationTime = tt.args.tokenCreationTime()
			}
			te := &logical.TokenEntry{
				Path:         "test",
				Policies:     []string{"default"},
				TTL:          time.Hour * 24,
				CreationTime: creationTime.Unix(),
			}
			testMakeTokenDirectly(t, c.tokenStore, te)
			assert.NotEmpty(t, te.ID)
			tt.args.authorizeRequest.ClientToken = te.ID

			// Create an assignment
			resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
				Path:      "oidc/assignment/test-assignment",
				Operation: logical.CreateOperation,
				Data: map[string]interface{}{
					"group_ids":  tt.args.assignment.GroupIDs,
					"entity_ids": tt.args.assignment.EntityIDs,
				},
				Storage: storage,
			})
			expectSuccess(t, resp, err)

			// Create a client
			resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
				Path:      "oidc/client/test-client",
				Operation: logical.CreateOperation,
				Storage:   storage,
				Data: map[string]interface{}{
					"key":              "test-key",
					"redirect_uris":    tt.args.client.RedirectURIs,
					"assignments":      tt.args.client.Assignments,
					"id_token_ttl":     tt.args.client.IDTokenTTL,
					"access_token_ttl": tt.args.client.AccessTokenTTL,
				},
			})
			expectSuccess(t, resp, err)

			// Read the client ID
			resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
				Path:      "oidc/client/test-client",
				Operation: logical.ReadOperation,
				Storage:   storage,
			})
			expectSuccess(t, resp, err)
			assert.NotNil(t, resp.Data["client_id"])
			clientID := resp.Data["client_id"].(string)

			// Use allowed client IDs if set by test args
			if len(tt.args.provider.AllowedClientIDs) == 0 {
				tt.args.provider.AllowedClientIDs = []string{clientID}
			}

			// Create a provider
			resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
				Path:      "oidc/provider/test-provider",
				Operation: logical.CreateOperation,
				Data: map[string]interface{}{
					"issuer":             tt.args.provider.Issuer,
					"allowed_client_ids": tt.args.provider.AllowedClientIDs,
					"scopes":             tt.args.provider.Scopes,
				},
				Storage: storage,
			})
			expectSuccess(t, resp, err)

			// Use the client ID if set by test args
			if len(tt.args.authorizeRequest.Data["client_id"].(string)) == 0 {
				tt.args.authorizeRequest.Data["client_id"] = clientID
			}

			// Send the request to the OIDC authorize endpoint
			tt.args.authorizeRequest.Storage = storage
			tt.args.authorizeRequest.EntityID = tt.args.entityID
			resp, err = c.identityStore.HandleRequest(ctx, tt.args.authorizeRequest)

			// Parse the response
			var res struct {
				Code             string `json:"code"`
				State            string `json:"state"`
				Error            string `json:"error"`
				ErrorDescription string `json:"error_description"`
			}
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Data[logical.HTTPRawBody])
			assert.NotNil(t, resp.Data[logical.HTTPContentType])
			assert.Equal(t, "application/json", resp.Data[logical.HTTPContentType].(string))
			assert.NoError(t, json.Unmarshal(resp.Data["http_raw_body"].([]byte), &res))

			if tt.wantErr != "" {
				// Assert that we receive the expected error code
				assert.Equal(t, tt.wantErr, res.Error)
				assert.NotEmpty(t, res.ErrorDescription)
				return
			}

			// Assert that we receive an authorization code (base62) and state
			expectSuccess(t, resp, err)
			assert.Regexp(t, "[a-zA-Z0-9]{32}", res.Code)
			assert.NotEmpty(t, res.State)
			assert.Empty(t, res.Error)
			assert.Empty(t, res.ErrorDescription)
		})
	}
}

// TestOIDC_Path_OIDC_ProviderReadPublicKey_ProviderDoesNotExist tests that the
// path can handle the read operation when the provider does not exist
func TestOIDC_Path_OIDC_ProviderReadPublicKey_ProviderDoesNotExist(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Read "test-provider" .well-known keys
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider/.well-known/keys",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectedResp := &logical.Response{}
	if resp != expectedResp && err != nil {
		t.Fatalf("expected empty response but got success; error:\n%v\nresp: %#v", err, resp)
	}
}

// TestOIDC_Path_OIDC_ProviderReadPublicKey tests the provider .well-known
// keys endpoint read operations
func TestOIDC_Path_OIDC_ProviderReadPublicKey(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key-1"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key-1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test client "test-client-1"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client-1",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key-1",
		},
	})

	// get the clientID
	resp, _ := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client-1",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	clientID := resp.Data["client_id"].(string)

	// Create a test provider "test-provider" and allow all client IDs -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"issuer":             "https://example.com:8200",
			"allowed_client_ids": []string{"*"},
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-provider" .well-known keys
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider/.well-known/keys",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	responseJWKS := &jose.JSONWebKeySet{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), responseJWKS)
	if len(responseJWKS.Keys) != 2 {
		t.Fatalf("expected 2 public key but instead got %d", len(responseJWKS.Keys))
	}

	// Create a test key "test-key-2"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key-2",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test client "test-client-2"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client-2",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key-2",
		},
	})

	// Read "test-provider" .well-known keys
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider/.well-known/keys",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	responseJWKS = &jose.JSONWebKeySet{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), responseJWKS)
	if len(responseJWKS.Keys) != 4 {
		t.Fatalf("expected 4 public key but instead got %d", len(responseJWKS.Keys))
	}

	// Update the test provider "test-provider" to only allow test-client-1 -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"allowed_client_ids": []string{clientID},
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-provider" .well-known keys
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider/.well-known/keys",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	responseJWKS = &jose.JSONWebKeySet{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), responseJWKS)
	if len(responseJWKS.Keys) != 2 {
		t.Fatalf("expected 2 public key but instead got %d", len(responseJWKS.Keys))
	}
}

// TestOIDC_Path_OIDC_ProviderClient_NoKeyParameter tests that a client cannot
// be created without a key parameter
func TestOIDC_Path_OIDC_ProviderClient_NoKeyParameter(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test client "test-client1" without a key param -- should fail
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client1",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"the key parameter is required": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderClient_NilKeyEntry tests that a client cannot be
// created when a key parameter is provided but the key does not exist
func TestOIDC_Path_OIDC_ProviderClient_NilKeyEntry(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test client "test-client1" with a non-existent key -- should fail
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key": "test-key",
		},
		Storage: storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"key \"test-key\" does not exist": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderClient_InvalidTokenTTL tests the TokenTTL validation
func TestOIDC_Path_OIDC_ProviderClient_InvalidTokenTTL(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": int64(60),
		},
		Storage: storage,
	})

	// Create a test role "test-client" with an id_token_ttl longer than the
	// verification_ttl -- should fail with error
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key":          "test-key",
			"id_token_ttl": int64(3600),
		},
		Storage: storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"a client's id_token_ttl cannot be greater than the verification_ttl of the key it references": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)

	// Read "test-client"
	respReadTestClient, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	// Ensure that "test-client" was not created
	expectSuccess(t, respReadTestClient, err)
	if respReadTestClient != nil {
		t.Fatalf("Expected a nil response but instead got:\n%#v", respReadTestClient)
	}
}

// TestOIDC_Path_OIDC_ProviderClient_UpdateKey tests that a client
// does not allow key modification on Update operations
func TestOIDC_Path_OIDC_ProviderClient_UpdateKey(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key1"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test key "test-key2"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key2",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test client "test-client" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key1",
		},
	})
	expectSuccess(t, resp, err)

	// Update the test client "test-client" -- should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key2",
		},
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"key modification is not allowed": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderClient_AssignmentDoesNotExist tests that a client
// cannot be created with assignments that do not exist
func TestOIDC_Path_OIDC_ProviderClient_AssignmentDoesNotExist(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test client "test-client" -- should fail
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key":         "test-key",
			"assignments": "my-assignment",
		},
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"assignment \"my-assignment\" does not exist": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderClient tests CRUD operations for clients
func TestOIDC_Path_OIDC_ProviderClient(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test client "test-client" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key",
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-client" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"redirect_uris":    []string{},
		"assignments":      []string{},
		"key":              "test-key",
		"id_token_ttl":     int64(86400),
		"access_token_ttl": int64(86400),
		"client_id":        resp.Data["client_id"],
		"client_secret":    resp.Data["client_secret"],
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Create a test assignment "my-assignment" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/my-assignment",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Update "test-client" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"redirect_uris":    "http://localhost:3456/callback",
			"assignments":      "my-assignment",
			"key":              "test-key",
			"id_token_ttl":     "90s",
			"access_token_ttl": "1m",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-client" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"redirect_uris":    []string{"http://localhost:3456/callback"},
		"assignments":      []string{"my-assignment"},
		"key":              "test-key",
		"id_token_ttl":     int64(90),
		"access_token_ttl": int64(60),
		"client_id":        resp.Data["client_id"],
		"client_secret":    resp.Data["client_secret"],
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Delete test-client -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-client" again and validate
	resp, _ = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if resp != nil {
		t.Fatalf("expected nil but got resp: %#v", resp)
	}
}

// TestOIDC_Path_OIDC_ProviderClient_DeDuplication tests that a
// client doesn't have duplicate redirect URIs or Assignments
func TestOIDC_Path_OIDC_ProviderClient_Deduplication(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test assignment "test-assignment1" -- should succeed
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment1",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// Create a test client "test-client" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key":           "test-key",
			"assignments":   []string{"test-assignment1", "test-assignment1"},
			"redirect_uris": []string{"http://example.com", "http://notduplicate.com", "http://example.com"},
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-client" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"redirect_uris":    []string{"http://example.com", "http://notduplicate.com"},
		"assignments":      []string{"test-assignment1"},
		"key":              "test-key",
		"id_token_ttl":     int64(86400),
		"access_token_ttl": int64(86400),
		"client_id":        resp.Data["client_id"],
		"client_secret":    resp.Data["client_secret"],
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDC_ProviderClient_Update tests Update operations for clients
func TestOIDC_Path_OIDC_ProviderClient_Update(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test assignment "my-assignment" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/my-assignment",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Create a test client "test-client" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"redirect_uris":    "http://localhost:3456/callback",
			"assignments":      "my-assignment",
			"key":              "test-key",
			"id_token_ttl":     "2m",
			"access_token_ttl": "1h",
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-client" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"redirect_uris":    []string{"http://localhost:3456/callback"},
		"assignments":      []string{"my-assignment"},
		"key":              "test-key",
		"id_token_ttl":     int64(120),
		"access_token_ttl": int64(3600),
		"client_id":        resp.Data["client_id"],
		"client_secret":    resp.Data["client_secret"],
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-client" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"redirect_uris":    "http://localhost:3456/callback2",
			"id_token_ttl":     "30",
			"access_token_ttl": "1m",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-client" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"redirect_uris":    []string{"http://localhost:3456/callback2"},
		"assignments":      []string{"my-assignment"},
		"key":              "test-key",
		"id_token_ttl":     int64(30),
		"access_token_ttl": int64(60),
		"client_id":        resp.Data["client_id"],
		"client_secret":    resp.Data["client_secret"],
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDC_ProviderClient_List tests the List operation for clients
func TestOIDC_Path_OIDC_ProviderClient_List(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Prepare two clients, test-client1 and test-client2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client1",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key",
		},
	})

	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client2",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key": "test-key",
		},
	})

	// list clients
	respListClients, listErr := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListClients, listErr)

	// validate list response
	expectedStrings := map[string]interface{}{"test-client1": true, "test-client2": true}
	expectStrings(t, respListClients.Data["keys"].([]string), expectedStrings)

	// delete test-client2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client2",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// list clients again and validate response
	respListClientAfterDelete, listErrAfterDelete := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListClientAfterDelete, listErrAfterDelete)

	// validate list response
	delete(expectedStrings, "test-client2")
	expectStrings(t, respListClientAfterDelete.Data["keys"].([]string), expectedStrings)
}

// TestOIDC_pathOIDCClientExistenceCheck tests pathOIDCClientExistenceCheck
func TestOIDC_pathOIDCClientExistenceCheck(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	clientName := "test"

	// Expect nil with empty storage
	exists, err := c.identityStore.pathOIDCClientExistenceCheck(
		ctx,
		&logical.Request{
			Storage: storage,
		},
		&framework.FieldData{
			Raw: map[string]interface{}{"name": clientName},
			Schema: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error during existence check on an expected nil entry, err:\n%#v", err)
	}
	if exists {
		t.Fatalf("Expected existence check to return false but instead returned: %t", exists)
	}

	// Populte storage with a client
	client := &client{}
	entry, _ := logical.StorageEntryJSON(clientPath+clientName, client)
	if err := storage.Put(ctx, entry); err != nil {
		t.Fatalf("writing to in mem storage failed")
	}

	// Expect true with a populated storage
	exists, err = c.identityStore.pathOIDCClientExistenceCheck(
		ctx,
		&logical.Request{
			Storage: storage,
		},
		&framework.FieldData{
			Raw: map[string]interface{}{"name": clientName},
			Schema: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error during existence check on an expected nil entry, err:\n%#v", err)
	}
	if !exists {
		t.Fatalf("Expected existence check to return true but instead returned: %t", exists)
	}
}

// TestOIDC_Path_OIDC_ProviderScope_ReservedName tests that the reserved name
// "openid" cannot be used when creating a scope
func TestOIDC_Path_OIDC_ProviderScope_ReservedName(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test scope "test-scope" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/openid",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"the \"openid\" scope name is reserved": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
}

// TestOIDC_Path_OIDC_ProviderScope_TemplateValidation tests that the template
// validation does not allow restricted claims
func TestOIDC_Path_OIDC_ProviderScope_TemplateValidation(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	testCases := []struct {
		templ         string
		restrictedKey string
	}{
		{
			templ:         `{"aud": "client-12345", "other": "test"}`,
			restrictedKey: "aud",
		},
		{
			templ:         `{"exp": 1311280970, "other": "test"}`,
			restrictedKey: "exp",
		},
		{
			templ:         `{"iat": 1311280970, "other": "test"}`,
			restrictedKey: "iat",
		},
		{
			templ:         `{"iss": "https://openid.c2id.com", "other": "test"}`,
			restrictedKey: "iss",
		},
		{
			templ:         `{"namespace": "n-0S6_WzA2Mj", "other": "test"}`,
			restrictedKey: "namespace",
		},
		{
			templ:         `{"sub": "alice", "other": "test"}`,
			restrictedKey: "sub",
		},
	}
	for _, tc := range testCases {
		encodedTempl := base64.StdEncoding.EncodeToString([]byte(tc.templ))
		// Create a test scope "test-scope" -- should fail
		resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
			Path:      "oidc/scope/test-scope",
			Operation: logical.CreateOperation,
			Storage:   storage,
			Data: map[string]interface{}{
				"template":    encodedTempl,
				"description": "my-description",
			},
		})
		expectError(t, resp, err)
		errString := fmt.Sprintf(
			"top level key %q not allowed. Restricted keys: iat, aud, exp, iss, sub, namespace",
			tc.restrictedKey,
		)
		// validate error message
		expectedStrings := map[string]interface{}{
			errString: true,
		}
		expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)
	}
}

// TestOIDC_Path_OIDC_ProviderScope tests CRUD operations for scopes
func TestOIDC_Path_OIDC_ProviderScope(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test scope "test-scope" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-scope" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"template":    "",
		"description": "",
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	templ := `{ "groups": {{identity.entity.groups.names}} }`
	encodedTempl := base64.StdEncoding.EncodeToString([]byte(templ))
	// Update "test-scope" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"template":    encodedTempl,
			"description": "my-description",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-scope" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"template":    templ,
		"description": "my-description",
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Delete test-scope -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-scope" again and validate
	resp, _ = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if resp != nil {
		t.Fatalf("expected nil but got resp: %#v", resp)
	}
}

// TestOIDC_Path_OIDC_ProviderScope_Update tests Update operations for scopes
func TestOIDC_Path_OIDC_ProviderScope_Update(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	templ := `{ "groups": {{identity.entity.groups.names}} }`
	encodedTempl := base64.StdEncoding.EncodeToString([]byte(templ))
	// Create a test scope "test-scope" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"template":    encodedTempl,
			"description": "my-description",
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-scope" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"template":    templ,
		"description": "my-description",
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-scope" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"template":    encodedTempl,
			"description": "my-description-2",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-scope" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"template":    "{ \"groups\": {{identity.entity.groups.names}} }",
		"description": "my-description-2",
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDC_ProviderScope_List tests the List operation for scopes
func TestOIDC_Path_OIDC_ProviderScope_List(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Prepare two scopes, test-scope1 and test-scope2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope1",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope2",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// list scopes
	respListScopes, listErr := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListScopes, listErr)

	// validate list response
	expectedStrings := map[string]interface{}{"test-scope1": true, "test-scope2": true}
	expectStrings(t, respListScopes.Data["keys"].([]string), expectedStrings)

	// delete test-scope2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope2",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// list scopes again and validate response
	respListScopeAfterDelete, listErrAfterDelete := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListScopeAfterDelete, listErrAfterDelete)

	// validate list response
	delete(expectedStrings, "test-scope2")
	expectStrings(t, respListScopeAfterDelete.Data["keys"].([]string), expectedStrings)
}

// TestOIDC_pathOIDCScopeExistenceCheck tests pathOIDCScopeExistenceCheck
func TestOIDC_pathOIDCScopeExistenceCheck(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	scopeName := "test"

	// Expect nil with empty storage
	exists, err := c.identityStore.pathOIDCScopeExistenceCheck(
		ctx,
		&logical.Request{
			Storage: storage,
		},
		&framework.FieldData{
			Raw: map[string]interface{}{"name": scopeName},
			Schema: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error during existence check on an expected nil entry, err:\n%#v", err)
	}
	if exists {
		t.Fatalf("Expected existence check to return false but instead returned: %t", exists)
	}

	// Populte storage with a scope
	scope := &scope{}
	entry, _ := logical.StorageEntryJSON(scopePath+scopeName, scope)
	if err := storage.Put(ctx, entry); err != nil {
		t.Fatalf("writing to in mem storage failed")
	}

	// Expect true with a populated storage
	exists, err = c.identityStore.pathOIDCScopeExistenceCheck(
		ctx,
		&logical.Request{
			Storage: storage,
		},
		&framework.FieldData{
			Raw: map[string]interface{}{"name": scopeName},
			Schema: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error during existence check on an expected nil entry, err:\n%#v", err)
	}
	if !exists {
		t.Fatalf("Expected existence check to return true but instead returned: %t", exists)
	}
}

// TestOIDC_Path_OIDC_ProviderScope_DeleteWithExistingProvider tests that a
// Scope cannot be deleted when it is referenced by a provider
func TestOIDC_Path_OIDC_ProviderScope_DeleteWithExistingProvider(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test scope "test-scope" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"template":    `{"groups": "{{identity.entity.groups.names}}"}`,
			"description": "my-description",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Create a test provider "test-provider"
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"scopes": []string{"test-scope"},
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Delete test-scope -- should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"unable to delete scope \"test-scope\" because it is currently referenced by these providers: test-provider": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)

	// Read "test-scope" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
}

// TestOIDC_Path_OIDC_ProviderAssignment tests CRUD operations for assignments
func TestOIDC_Path_OIDC_ProviderAssignment(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test assignment "test-assignment" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-assignment" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"group_ids":  []string{},
		"entity_ids": []string{},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-assignment" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"group_ids":  "my-group",
			"entity_ids": "my-entity",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-assignment" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"group_ids":  []string{"my-group"},
		"entity_ids": []string{"my-entity"},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Delete test-assignment -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-assignment" again and validate
	resp, _ = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if resp != nil {
		t.Fatalf("expected nil but got resp: %#v", resp)
	}
}

// TestOIDC_Path_OIDC_ProviderAssignment_DeleteWithExistingClient tests that an
// assignment cannot be deleted when it is referenced by a client
func TestOIDC_Path_OIDC_ProviderAssignment_DeleteWithExistingClient(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test assignment "test-assignment" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Create a test key "test-key"
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"verification_ttl": "2m",
			"rotation_period":  "2m",
		},
		Storage: storage,
	})

	// Create a test client "test-client" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key":         "test-key",
			"assignments": []string{"test-assignment"},
		},
	})
	expectSuccess(t, resp, err)

	// Delete test-assignment -- should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"unable to delete assignment \"test-assignment\" because it is currently referenced by these clients: test-client": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)

	// Read "test-assignment" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"group_ids":  []string{},
		"entity_ids": []string{},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDC_ProviderAssignment_Update tests Update operations for assignments
func TestOIDC_Path_OIDC_ProviderAssignment_Update(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test assignment "test-assignment" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"group_ids":  "my-group",
			"entity_ids": "my-entity",
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-assignment" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"group_ids":  []string{"my-group"},
		"entity_ids": []string{"my-entity"},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-assignment" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"group_ids": "my-group2",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-assignment" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"group_ids":  []string{"my-group2"},
		"entity_ids": []string{"my-entity"},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDC_ProviderAssignment_List tests the List operation for assignments
func TestOIDC_Path_OIDC_ProviderAssignment_List(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Prepare two assignments, test-assignment1 and test-assignment2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment1",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment2",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// list assignments
	respListAssignments, listErr := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListAssignments, listErr)

	// validate list response
	expectedStrings := map[string]interface{}{"test-assignment1": true, "test-assignment2": true}
	expectStrings(t, respListAssignments.Data["keys"].([]string), expectedStrings)

	// delete test-assignment2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment/test-assignment2",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// list assignments again and validate response
	respListAssignmentAfterDelete, listErrAfterDelete := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/assignment",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListAssignmentAfterDelete, listErrAfterDelete)

	// validate list response
	delete(expectedStrings, "test-assignment2")
	expectStrings(t, respListAssignmentAfterDelete.Data["keys"].([]string), expectedStrings)
}

// TestOIDC_pathOIDCAssignmentExistenceCheck tests pathOIDCAssignmentExistenceCheck
func TestOIDC_pathOIDCAssignmentExistenceCheck(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	assignmentName := "test"

	// Expect nil with empty storage
	exists, err := c.identityStore.pathOIDCAssignmentExistenceCheck(
		ctx,
		&logical.Request{
			Storage: storage,
		},
		&framework.FieldData{
			Raw: map[string]interface{}{"name": assignmentName},
			Schema: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error during existence check on an expected nil entry, err:\n%#v", err)
	}
	if exists {
		t.Fatalf("Expected existence check to return false but instead returned: %t", exists)
	}

	// Populte storage with a assignment
	assignment := &assignment{}
	entry, _ := logical.StorageEntryJSON(assignmentPath+assignmentName, assignment)
	if err := storage.Put(ctx, entry); err != nil {
		t.Fatalf("writing to in mem storage failed")
	}

	// Expect true with a populated storage
	exists, err = c.identityStore.pathOIDCAssignmentExistenceCheck(
		ctx,
		&logical.Request{
			Storage: storage,
		},
		&framework.FieldData{
			Raw: map[string]interface{}{"name": assignmentName},
			Schema: map[string]*framework.FieldSchema{
				"name": {
					Type: framework.TypeString,
				},
			},
		},
	)
	if err != nil {
		t.Fatalf("Error during existence check on an expected nil entry, err:\n%#v", err)
	}
	if !exists {
		t.Fatalf("Expected existence check to return true but instead returned: %t", exists)
	}
}

// TestOIDC_Path_OIDCProvider tests CRUD operations for providers
func TestOIDC_Path_OIDCProvider(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test provider "test-provider" with non-existing scope
	// Should fail
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"scopes": []string{"test-scope"},
		},
		Storage: storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings := map[string]interface{}{
		"scope \"test-scope\" does not exist": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)

	// Create a test provider "test-provider" with no scopes -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-provider" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"issuer":             "",
		"allowed_client_ids": []string{},
		"scopes":             []string{},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Create a test scope "test-scope" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"template":    `{"groups": {{identity.entity.groups.names}} }`,
			"description": "my-description",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Update "test-provider" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"allowed_client_ids": []string{"test-client-id"},
			"scopes":             []string{"test-scope"},
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-provider" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"issuer":             "",
		"allowed_client_ids": []string{"test-client-id"},
		"scopes":             []string{"test-scope"},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-provider" -- should fail issuer validation
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"issuer": "test-issuer",
		},
		Storage: storage,
	})
	expectError(t, resp, err)
	// validate error message
	expectedStrings = map[string]interface{}{
		"invalid issuer, which must include only a scheme, host, and optional port (e.g. https://example.com:8200)": true,
	}
	expectStrings(t, []string{resp.Data["error"].(string)}, expectedStrings)

	// Update "test-provider" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"issuer": "https://example.com:8200",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-provider" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"issuer":             "https://example.com:8200",
		"allowed_client_ids": []string{"test-client-id"},
		"scopes":             []string{"test-scope"},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Delete test-provider -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-provider" again and validate
	resp, _ = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	if resp != nil {
		t.Fatalf("expected nil but got resp: %#v", resp)
	}
}

// TestOIDC_Path_OIDCProvider_DuplicateTempalteKeys tests that no two
// scopes have the same top-level keys when creating a provider
func TestOIDC_Path_OIDCProvider_DuplicateTemplateKeys(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test scope "test-scope1" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"template":    `{"groups": {{identity.entity.groups.names}} }`,
			"description": "desc1",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Create another test scope "test-scope2" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope2",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"template":    `{"groups": {{identity.entity.groups.names}} }`,
			"description": "desc2",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Create a test provider "test-provider" with scopes that have same top-level keys
	// Should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"scopes": []string{"test-scope1", "test-scope2"},
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)
	if resp.Warnings[0] != "Found scope templates with conflicting top-level keys: conflict \"groups\" in scopes \"test-scope2\", \"test-scope1\". This may result in an error if the scopes are requested in an OIDC Authentication Request." {
		t.Fatalf("expected a warning for conflicting keys, got %s", resp.Warnings[0])
	}

	// // Update "test-scope1" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope1",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"template": `{"roles": {{identity.entity.groups.names}} }`,
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Create a test provider "test-provider" with updated scopes
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"scopes": []string{"test-scope1", "test-scope2"},
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)
}

// TestOIDC_Path_OIDCProvider_DeDuplication tests that a
// provider doensn't have duplicate scopes or client IDs
func TestOIDC_Path_OIDCProvider_Deduplication(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test scope "test-scope1" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"template":    `{"groups": {{identity.entity.groups.names}} }`,
			"description": "desc1",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Create a test provider "test-provider" with duplicates
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"scopes":             []string{"test-scope1", "test-scope1"},
			"allowed_client_ids": []string{"test-id1", "test-id2", "test-id1"},
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-provider" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"issuer":             "",
		"allowed_client_ids": []string{"test-id1", "test-id2"},
		"scopes":             []string{"test-scope1"},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDCProvider_Update tests Update operations for providers
func TestOIDC_Path_OIDCProvider_Update(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test provider "test-provider" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"issuer":             "https://example.com:8200",
			"allowed_client_ids": []string{"test-client-id"},
		},
	})
	expectSuccess(t, resp, err)

	// Read "test-provider" and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected := map[string]interface{}{
		"issuer":             "https://example.com:8200",
		"allowed_client_ids": []string{"test-client-id"},
		"scopes":             []string{},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}

	// Update "test-provider" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"issuer": "https://changedurl.com",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Read "test-provider" again and validate
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	expected = map[string]interface{}{
		"issuer":             "https://changedurl.com",
		"allowed_client_ids": []string{"test-client-id"},
		"scopes":             []string{},
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OIDC_ProviderList tests the List operation for providers
func TestOIDC_Path_OIDC_Provider_List(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Prepare two providers, test-provider1 and test-provider2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider1",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider2",
		Operation: logical.CreateOperation,
		Storage:   storage,
	})

	// list providers
	respListProviders, listErr := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListProviders, listErr)

	// validate list response
	expectedStrings := map[string]interface{}{"test-provider1": true, "test-provider2": true}
	expectStrings(t, respListProviders.Data["keys"].([]string), expectedStrings)

	// delete test-provider2
	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider2",
		Operation: logical.DeleteOperation,
		Storage:   storage,
	})

	// list providers again and validate response
	respListProvidersAfterDelete, listErrAfterDelete := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider",
		Operation: logical.ListOperation,
		Storage:   storage,
	})
	expectSuccess(t, respListProvidersAfterDelete, listErrAfterDelete)

	// validate list response
	delete(expectedStrings, "test-provider2")
	expectStrings(t, respListProvidersAfterDelete.Data["keys"].([]string), expectedStrings)
}

// TestOIDC_Path_OpenIDProviderConfig tests read operations for the
// openid-configuration path
func TestOIDC_Path_OpenIDProviderConfig(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Create a test scope "test-scope-1" -- should succeed
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope-1",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"template":    `{"groups": "{{identity.entity.groups.names}}"}`,
			"description": "my-description",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Create a test provider "test-provider"
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"scopes": []string{"test-scope-1"},
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Expect defaults from .well-known/openid-configuration
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider/.well-known/openid-configuration",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)

	basePath := "/v1/identity/oidc/provider/test-provider"
	expected := &providerDiscovery{
		Issuer:                basePath,
		Keys:                  basePath + "/.well-known/keys",
		ResponseTypes:         []string{"code"},
		Scopes:                []string{"test-scope-1", "openid"},
		Subjects:              []string{"public"},
		IDTokenAlgs:           supportedAlgs,
		AuthorizationEndpoint: "/ui/vault/identity/oidc/provider/test-provider/authorize",
		TokenEndpoint:         basePath + "/token",
		UserinfoEndpoint:      basePath + "/userinfo",
		GrantTypes:            []string{"authorization_code"},
		AuthMethods:           []string{"client_secret_basic"},
		RequestURIParameter:   false,
	}
	discoveryResp := &providerDiscovery{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), discoveryResp)
	if diff := deep.Equal(expected, discoveryResp); diff != nil {
		t.Fatal(diff)
	}

	// Create a test scope "test-scope-2" -- should succeed
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/scope/test-scope-2",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"template":    `{"groups": "{{identity.entity.groups.names}}"}`,
			"description": "my-description",
		},
		Storage: storage,
	})
	expectSuccess(t, resp, err)

	// Update provider issuer config
	testIssuer := "https://example.com:1234"
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"issuer": testIssuer,
			"scopes": []string{"test-scope-2"},
		},
	})
	expectSuccess(t, resp, err)

	// Expect updates from .well-known/openid-configuration
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider/.well-known/openid-configuration",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectSuccess(t, resp, err)
	// Validate
	basePath = testIssuer + basePath
	expected = &providerDiscovery{
		Issuer:                basePath,
		Keys:                  basePath + "/.well-known/keys",
		ResponseTypes:         []string{"code"},
		Scopes:                []string{"test-scope-2", "openid"},
		Subjects:              []string{"public"},
		IDTokenAlgs:           supportedAlgs,
		AuthorizationEndpoint: testIssuer + "/ui/vault/identity/oidc/provider/test-provider/authorize",
		TokenEndpoint:         basePath + "/token",
		UserinfoEndpoint:      basePath + "/userinfo",
		GrantTypes:            []string{"authorization_code"},
		AuthMethods:           []string{"client_secret_basic"},
		RequestURIParameter:   false,
	}
	discoveryResp = &providerDiscovery{}
	json.Unmarshal(resp.Data["http_raw_body"].([]byte), discoveryResp)
	if diff := deep.Equal(expected, discoveryResp); diff != nil {
		t.Fatal(diff)
	}
}

// TestOIDC_Path_OpenIDProviderConfig_ProviderDoesNotExist tests read
// operations for the openid-configuration path when the provider does not
// exist
func TestOIDC_Path_OpenIDProviderConfig_ProviderDoesNotExist(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	storage := &logical.InmemStorage{}

	// Expect defaults from .well-known/openid-configuration
	// test-provider does not exist
	resp, err := c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/provider/test-provider/.well-known/openid-configuration",
		Operation: logical.ReadOperation,
		Storage:   storage,
	})
	expectedResp := &logical.Response{}
	if resp != expectedResp && err != nil {
		t.Fatalf("expected empty response but got success; error:\n%v\nresp: %#v", err, resp)
	}
}
