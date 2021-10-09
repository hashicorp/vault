package vault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"
)

func TestOIDC_Path_OIDC_Token(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	s := new(logical.InmemStorage)

	entityID, groupID, clientID, clientSecret := setupOIDCCommon(t, c, s)

	type args struct {
		clientReq              *logical.Request
		providerReq            *logical.Request
		assignmentReq          *logical.Request
		authorizeReq           *logical.Request
		tokenReq               *logical.Request
		vaultTokenCreationTime func() time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "invalid token request with provider not found",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq: func() *logical.Request {
					req := validTokenReq(s, "", clientID, clientSecret)
					req.Path = "oidc/provider/non-existent-provider/token"
					return req
				}(),
			},
			wantErr: ErrTokenInvalidRequest,
		},
		{
			name: "invalid token request with missing basic auth header",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq: func() *logical.Request {
					req := validTokenReq(s, "", clientID, clientSecret)
					req.Headers = nil
					return req
				}(),
			},
			wantErr: ErrTokenInvalidRequest,
		},
		{
			name: "invalid token request with client ID not found",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq:      validTokenReq(s, "", "non-existent-client-id", clientSecret),
			},
			wantErr: ErrTokenInvalidClient,
		},
		{
			name: "invalid token request with client secret mismatch",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq:      validTokenReq(s, "", clientID, "wrong-client-secret"),
			},
			wantErr: ErrTokenInvalidClient,
		},
		{
			name: "invalid token request with client_id not allowed by provider",
			args: args{
				clientReq: testClientReq(s),
				providerReq: func() *logical.Request {
					req := testProviderReq(s, clientID)
					req.Data["allowed_client_ids"] = []string{"not-client-id"}
					return req
				}(),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq:      validTokenReq(s, "", clientID, clientSecret),
			},
			wantErr: ErrTokenInvalidClient,
		},
		{
			name: "invalid token request with empty grant_type",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq: func() *logical.Request {
					req := validTokenReq(s, "", clientID, clientSecret)
					req.Data["grant_type"] = ""
					return req
				}(),
			},
			wantErr: ErrTokenInvalidRequest,
		},
		{
			name: "invalid token request with unsupported grant_type",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq: func() *logical.Request {
					req := validTokenReq(s, "", clientID, clientSecret)
					req.Data["grant_type"] = "not-supported-grant-type"
					return req
				}(),
			},
			wantErr: ErrTokenUnsupportedGrantType,
		},
		{
			name: "invalid token request with invalid code",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq: func() *logical.Request {
					req := validTokenReq(s, "", clientID, clientSecret)
					req.Data["code"] = "invalid-code"
					return req
				}(),
			},
			wantErr: ErrTokenInvalidGrant,
		},
		{
			name: "invalid token request with missing redirect_uri",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq: func() *logical.Request {
					req := validTokenReq(s, "", clientID, clientSecret)
					req.Data["redirect_uri"] = ""
					return req
				}(),
			},
			wantErr: ErrTokenInvalidRequest,
		},
		{
			name: "invalid token request with entity not found in client assignment",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, "not-entity-id", ""),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq:      validTokenReq(s, "", clientID, clientSecret),
			},
			wantErr: ErrTokenInvalidRequest,
		},
		{
			name: "invalid token request with redirect_uri mismatch",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq: func() *logical.Request {
					req := validTokenReq(s, "", clientID, clientSecret)
					req.Data["redirect_uri"] = "https://not.original.redirect.uri:8251/callback"
					return req
				}(),
			},
			wantErr: ErrTokenInvalidGrant,
		},
		{
			name: "invalid token request with group not found in client assignment",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, "", "not-group-id"),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq:      validTokenReq(s, "", clientID, clientSecret),
			},
			wantErr: ErrTokenInvalidRequest,
		},
		{
			name: "valid token request with max_age and auth_time claim",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["max_age"] = "30"
					return req
				}(),
				tokenReq: validTokenReq(s, "", clientID, clientSecret),
				vaultTokenCreationTime: func() time.Time {
					return time.Now()
				},
			},
		},
		{
			name: "valid token request",
			args: args{
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
				tokenReq:      validTokenReq(s, "", clientID, clientSecret),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a token entry to associate with the authorize request
			creationTime := time.Now()
			if tt.args.vaultTokenCreationTime != nil {
				creationTime = tt.args.vaultTokenCreationTime()
			}
			te := &logical.TokenEntry{
				Path:         "test",
				Policies:     []string{"default"},
				TTL:          time.Hour * 24,
				CreationTime: creationTime.Unix(),
			}
			testMakeTokenDirectly(t, c.tokenStore, te)
			require.NotEmpty(t, te.ID)

			// Reset any configuration modifications
			resetCommonOIDCConfig(t, s, c, entityID, groupID, clientID)

			// Send the request to the OIDC authorize endpoint
			tt.args.authorizeReq.EntityID = entityID
			tt.args.authorizeReq.ClientToken = te.ID
			resp, err := c.identityStore.HandleRequest(ctx, tt.args.authorizeReq)
			expectSuccess(t, resp, err)

			// Parse the authorize response
			var authRes struct {
				Code  string `json:"code"`
				State string `json:"state"`
			}
			require.NoError(t, json.Unmarshal(resp.Data["http_raw_body"].([]byte), &authRes))
			require.Regexp(t, "[a-zA-Z0-9]{32}", authRes.Code)
			require.NotEmpty(t, authRes.State)

			// Update the assignment
			tt.args.assignmentReq.Operation = logical.UpdateOperation
			resp, err = c.identityStore.HandleRequest(ctx, tt.args.assignmentReq)
			expectSuccess(t, resp, err)

			// Update the client
			tt.args.clientReq.Operation = logical.UpdateOperation
			resp, err = c.identityStore.HandleRequest(ctx, tt.args.clientReq)
			expectSuccess(t, resp, err)

			// Update the provider
			tt.args.providerReq.Operation = logical.UpdateOperation
			resp, err = c.identityStore.HandleRequest(ctx, tt.args.providerReq)
			expectSuccess(t, resp, err)

			// Update the code if provided by test arguments
			authCode := authRes.Code
			if tt.args.tokenReq.Data["code"] != "" {
				authCode = tt.args.tokenReq.Data["code"].(string)
			}

			// Send the request to the OIDC token endpoint
			tt.args.tokenReq.Data["code"] = authCode
			resp, err = c.identityStore.HandleRequest(ctx, tt.args.tokenReq)
			expectSuccess(t, resp, err)

			// Parse the token response
			var tokenRes struct {
				TokenType        string `json:"token_type"`
				AccessToken      string `json:"access_token"`
				IDToken          string `json:"id_token"`
				ExpiresIn        int64  `json:"expires_in"`
				Error            string `json:"error"`
				ErrorDescription string `json:"error_description"`
			}
			require.NotNil(t, resp)
			require.NotNil(t, resp.Data[logical.HTTPRawBody])
			require.NotNil(t, resp.Data[logical.HTTPStatusCode])
			require.NotNil(t, resp.Data[logical.HTTPContentType])
			require.NotNil(t, resp.Data[logical.HTTPPragmaHeader])
			require.NotNil(t, resp.Data[logical.HTTPCacheControlHeader])
			require.Equal(t, "no-cache", resp.Data[logical.HTTPPragmaHeader])
			require.Equal(t, "no-store", resp.Data[logical.HTTPCacheControlHeader])
			require.Equal(t, "application/json", resp.Data[logical.HTTPContentType].(string))
			require.NoError(t, json.Unmarshal(resp.Data["http_raw_body"].([]byte), &tokenRes))

			if tt.wantErr != "" {
				// Assert that we receive the expected error code and description
				require.Equal(t, tt.wantErr, tokenRes.Error)
				require.NotEmpty(t, tokenRes.ErrorDescription)

				// Assert that we receive the expected status code
				statusCode := resp.Data[logical.HTTPStatusCode].(int)
				switch tokenRes.Error {
				case ErrTokenInvalidClient:
					require.Equal(t, http.StatusUnauthorized, statusCode)
					require.Equal(t, "Basic", resp.Data[logical.HTTPWWWAuthenticateHeader])
				case ErrTokenServerError:
					require.Equal(t, http.StatusInternalServerError, statusCode)
				default:
					require.Equal(t, http.StatusBadRequest, statusCode)
				}
				return
			}

			// Assert that we receive the expected token response
			expectSuccess(t, resp, err)
			require.Equal(t, http.StatusOK, resp.Data[logical.HTTPStatusCode].(int))
			require.Equal(t, "Bearer", tokenRes.TokenType)
			require.NotEmpty(t, tokenRes.AccessToken)
			require.NotEmpty(t, tokenRes.IDToken)
			require.NotEmpty(t, tokenRes.ExpiresIn)
			require.Empty(t, tokenRes.Error)
			require.Empty(t, tokenRes.ErrorDescription)
		})
	}
}

func TestOIDC_Path_OIDC_Authorize(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	ctx := namespace.RootContext(nil)
	s := new(logical.InmemStorage)

	entityID, groupID, clientID, _ := setupOIDCCommon(t, c, s)

	type args struct {
		entityID               string
		clientReq              *logical.Request
		providerReq            *logical.Request
		assignmentReq          *logical.Request
		authorizeReq           *logical.Request
		vaultTokenCreationTime func() time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "invalid authorize request with provider not found",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Path = "oidc/provider/non-existent-provider/authorize"
					return req
				}(),
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with empty scope",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["scope"] = ""
					return req
				}(),
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with missing openid scope",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["scope"] = "groups email profile"
					return req
				}(),
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with missing response_type",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["response_type"] = ""
					return req
				}(),
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with unsupported response_type",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["response_type"] = "id_token"
					return req
				}(),
			},
			wantErr: ErrAuthUnsupportedResponseType,
		},
		{
			name: "invalid authorize request with client_id not found",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					return validAuthorizeReq(s, "non-existent-client-id")
				}(),
			},
			wantErr: ErrAuthInvalidClientID,
		},
		{
			name: "invalid authorize request with client_id not allowed by provider",
			args: args{
				entityID:  entityID,
				clientReq: testClientReq(s),
				providerReq: func() *logical.Request {
					req := testProviderReq(s, clientID)
					req.Data["allowed_client_ids"] = []string{"not-client-id"}
					return req
				}(),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
			},
			wantErr: ErrAuthUnauthorizedClient,
		},
		{
			name: "invalid authorize request with missing redirect_uri",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["redirect_uri"] = ""
					return req
				}(),
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with redirect_uri not allowed by client",
			args: args{
				entityID: entityID,
				clientReq: func() *logical.Request {
					req := testClientReq(s)
					req.Data["redirect_uris"] = []string{"https://not.redirect.uri:8251/callback"}
					return req
				}(),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
			},
			wantErr: ErrAuthInvalidRedirectURI,
		},
		{
			name: "invalid authorize request with missing state",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["state"] = ""
					return req
				}(),
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with missing nonce",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["nonce"] = ""
					return req
				}(),
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "invalid authorize request with request parameter provided",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["request"] = "header.payload.signature"
					return req
				}(),
			},
			wantErr: ErrAuthRequestNotSupported,
		},
		{
			name: "invalid authorize request with request_uri parameter provided",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["request_uri"] = "https://client.example.org/request.jwt"
					return req
				}(),
			},
			wantErr: ErrAuthRequestURINotSupported,
		},
		{
			name: "invalid authorize request with identity entity ID not found",
			args: args{
				entityID:      "non-existent-entity",
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
			},
			wantErr: ErrAuthAccessDenied,
		},
		{
			name: "invalid authorize request with entity not found in client assignment",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, "not-entity-id", ""),
				authorizeReq:  validAuthorizeReq(s, clientID),
			},
			wantErr: ErrAuthAccessDenied,
		},
		{
			name: "invalid authorize request with group not found in client assignment",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, "", "not-group-id"),
				authorizeReq:  validAuthorizeReq(s, clientID),
			},
			wantErr: ErrAuthAccessDenied,
		},
		{
			name: "invalid authorize request with negative max_age",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["max_age"] = "-1"
					return req
				}(),
			},
			wantErr: ErrAuthInvalidRequest,
		},
		{
			name: "active re-authentication required with token creation time exceeding max_age requirement",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["max_age"] = "30"
					return req
				}(),
				vaultTokenCreationTime: func() time.Time {
					return time.Now().Add(-time.Minute)
				},
			},
			wantErr: ErrAuthMaxAgeReAuthenticate,
		},
		{
			name: "valid authorize request with token creation time within max_age requirement",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Data["max_age"] = "30"
					return req
				}(),
				vaultTokenCreationTime: func() time.Time {
					return time.Now()
				},
			},
		},
		{
			name: "valid authorize request using update operation (HTTP POST)",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
			},
		},
		{
			name: "valid authorize request using read operation (HTTP GET)",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, groupID),
				authorizeReq: func() *logical.Request {
					req := validAuthorizeReq(s, clientID)
					req.Operation = logical.ReadOperation
					return req
				}(),
			},
		},
		{
			name: "valid authorize request using client assignment with only entity membership",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, entityID, ""),
				authorizeReq:  validAuthorizeReq(s, clientID),
			},
		},
		{
			name: "valid authorize request using client assignment with only group membership",
			args: args{
				entityID:      entityID,
				clientReq:     testClientReq(s),
				providerReq:   testProviderReq(s, clientID),
				assignmentReq: testAssignmentReq(s, "", groupID),
				authorizeReq:  validAuthorizeReq(s, clientID),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a token entry to associate with the authorize request
			creationTime := time.Now()
			if tt.args.vaultTokenCreationTime != nil {
				creationTime = tt.args.vaultTokenCreationTime()
			}
			te := &logical.TokenEntry{
				Path:         "test",
				Policies:     []string{"default"},
				TTL:          time.Hour * 24,
				CreationTime: creationTime.Unix(),
			}
			testMakeTokenDirectly(t, c.tokenStore, te)
			require.NotEmpty(t, te.ID)

			// Update the assignment
			tt.args.assignmentReq.Operation = logical.UpdateOperation
			resp, err := c.identityStore.HandleRequest(ctx, tt.args.assignmentReq)
			expectSuccess(t, resp, err)

			// Update the client
			tt.args.clientReq.Operation = logical.UpdateOperation
			resp, err = c.identityStore.HandleRequest(ctx, tt.args.clientReq)
			expectSuccess(t, resp, err)

			// Update the provider
			tt.args.providerReq.Operation = logical.UpdateOperation
			resp, err = c.identityStore.HandleRequest(ctx, tt.args.providerReq)
			expectSuccess(t, resp, err)

			// Send the request to the OIDC authorize endpoint
			tt.args.authorizeReq.EntityID = tt.args.entityID
			tt.args.authorizeReq.ClientToken = te.ID
			resp, err = c.identityStore.HandleRequest(ctx, tt.args.authorizeReq)

			// Parse the response
			var authRes struct {
				Code             string `json:"code"`
				State            string `json:"state"`
				Error            string `json:"error"`
				ErrorDescription string `json:"error_description"`
			}
			require.NotNil(t, resp)
			require.NotNil(t, resp.Data[logical.HTTPRawBody])
			require.NotNil(t, resp.Data[logical.HTTPStatusCode])
			require.NotNil(t, resp.Data[logical.HTTPContentType])
			require.Equal(t, "application/json", resp.Data[logical.HTTPContentType].(string))
			require.NoError(t, json.Unmarshal(resp.Data["http_raw_body"].([]byte), &authRes))

			if tt.wantErr != "" {
				// Assert that we receive the expected error code and description
				require.Equal(t, tt.wantErr, authRes.Error)
				require.NotEmpty(t, authRes.ErrorDescription)

				// Assert that we receive the expected status code
				statusCode := resp.Data[logical.HTTPStatusCode].(int)
				switch authRes.Error {
				case ErrAuthServerError:
					require.Equal(t, http.StatusInternalServerError, statusCode)
				default:
					require.Equal(t, http.StatusBadRequest, statusCode)
				}
				return
			}

			// Assert that we receive an authorization code (base62) and state
			expectSuccess(t, resp, err)
			require.Equal(t, http.StatusOK, resp.Data[logical.HTTPStatusCode].(int))
			require.Regexp(t, "[a-zA-Z0-9]{32}", authRes.Code)
			require.NotEmpty(t, authRes.State)
			require.Empty(t, authRes.Error)
			require.Empty(t, authRes.ErrorDescription)
		})
	}
}

// setupOIDCCommon creates all of the resources needed to test a Vault OIDC provider.
// Returns the entity ID, group ID, and client ID to be used in tests.
func setupOIDCCommon(t *testing.T, c *Core, s logical.Storage) (string, string, string, string) {
	t.Helper()
	ctx := namespace.RootContext(nil)

	// Create a key
	resp, err := c.identityStore.HandleRequest(ctx, testKeyReq(s, []string{"*"}, "RS256"))
	expectSuccess(t, resp, err)

	// Create an entity
	resp, err = c.identityStore.HandleRequest(ctx, testEntityReq(s))
	expectSuccess(t, resp, err)
	require.NotNil(t, resp.Data["id"])
	entityID := resp.Data["id"].(string)

	// Create a group
	resp, err = c.identityStore.HandleRequest(ctx, testGroupReq(s, "test-group", []string{entityID}))
	expectSuccess(t, resp, err)
	require.NotNil(t, resp.Data["id"])
	groupID := resp.Data["id"].(string)

	// Create an assignment
	resp, err = c.identityStore.HandleRequest(ctx, testAssignmentReq(s, entityID, groupID))
	expectSuccess(t, resp, err)

	// Create a client
	resp, err = c.identityStore.HandleRequest(ctx, testClientReq(s))
	expectSuccess(t, resp, err)

	// Read the client ID and secret
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Storage:   s,
		Path:      "oidc/client/test-client",
		Operation: logical.ReadOperation,
	})
	expectSuccess(t, resp, err)
	require.NotNil(t, resp.Data["client_id"])
	require.NotNil(t, resp.Data["client_secret"])
	clientID := resp.Data["client_id"].(string)
	clientSecret := resp.Data["client_secret"].(string)

	// Create a custom scope
	template := `{
		"name": {{identity.entity.name}},
		"contact": {
			"email": {{identity.entity.metadata.email}},
			"phone_number": {{identity.entity.metadata.phone_number}}
		},
		"groups": {{identity.entity.groups.names}}
	}`
	resp, err = c.identityStore.HandleRequest(ctx, testScopeReq(s, template))
	expectSuccess(t, resp, err)

	// Create a provider
	resp, err = c.identityStore.HandleRequest(ctx, testProviderReq(s, clientID))
	expectSuccess(t, resp, err)

	return entityID, groupID, clientID, clientSecret
}

// resetCommonOIDCConfig resets the state of common configuration resources
// (i.e., created by setupOIDCCommon) that are modified during tests. This
// enables the tests to continue operating using the same underlying storage
// throughout many test cases that modify the configuration resources.
func resetCommonOIDCConfig(t *testing.T, s logical.Storage, c *Core, entityID, groupID, clientID string) {
	ctx := namespace.RootContext(nil)

	req := testAssignmentReq(s, entityID, groupID)
	req.Operation = logical.UpdateOperation
	resp, err := c.identityStore.HandleRequest(ctx, req)
	expectSuccess(t, resp, err)

	req = testClientReq(s)
	req.Operation = logical.UpdateOperation
	resp, err = c.identityStore.HandleRequest(ctx, req)
	expectSuccess(t, resp, err)

	req = testProviderReq(s, clientID)
	req.Operation = logical.UpdateOperation
	resp, err = c.identityStore.HandleRequest(ctx, req)
	expectSuccess(t, resp, err)
}

func validTokenReq(s logical.Storage, code, clientID, clientSecret string) *logical.Request {
	return &logical.Request{
		Storage:   s,
		Path:      "oidc/provider/test-provider/token",
		Operation: logical.UpdateOperation,
		Headers: map[string][]string{
			"Authorization": {basicAuthHeader(clientID, clientSecret)},
		},
		Data: map[string]interface{}{
			// The code is unknown until returned from the authorization endpoint
			"code":         code,
			"grant_type":   "authorization_code",
			"redirect_uri": "https://localhost:8251/callback",
		},
	}
}

func validAuthorizeReq(s logical.Storage, clientID string) *logical.Request {
	return &logical.Request{
		Storage:   s,
		Path:      "oidc/provider/test-provider/authorize",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"client_id":     clientID,
			"scope":         "openid",
			"redirect_uri":  "https://localhost:8251/callback",
			"response_type": "code",
			"state":         "abcdefg",
			"nonce":         "hijklmn",
		},
	}
}

func testAssignmentReq(s logical.Storage, entityID, groupID string) *logical.Request {
	return &logical.Request{
		Storage:   s,
		Path:      "oidc/assignment/test-assignment",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"entity_ids": []string{entityID},
			"group_ids":  []string{groupID},
		},
	}
}

func testClientReq(s logical.Storage) *logical.Request {
	return &logical.Request{
		Storage:   s,
		Path:      "oidc/client/test-client",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"key":              "test-key",
			"redirect_uris":    []string{"https://localhost:8251/callback"},
			"assignments":      []string{"test-assignment"},
			"id_token_ttl":     "24h",
			"access_token_ttl": "24h",
		},
	}
}

func testProviderReq(s logical.Storage, clientID string) *logical.Request {
	return &logical.Request{
		Storage:   s,
		Path:      "oidc/provider/test-provider",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"allowed_client_ids": []string{clientID},
			"scopes":             []string{"test-scope"},
		},
	}
}

func testEntityReq(s logical.Storage) *logical.Request {
	return &logical.Request{
		Storage:   s,
		Path:      "entity",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name": "test-entity",
			"metadata": map[string]string{
				"email":        "test@hashicorp.com",
				"phone_number": "123-456-7890",
			},
		},
	}
}

func testKeyReq(s logical.Storage, allowedClientIDs []string, alg string) *logical.Request {
	return &logical.Request{
		Storage:   s,
		Path:      "oidc/key/test-key",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"allowed_client_ids": allowedClientIDs,
			"algorithm":          alg,
		},
	}
}

func testGroupReq(s logical.Storage, name string, entityIDs []string) *logical.Request {
	return &logical.Request{
		Storage:   s,
		Path:      "group",
		Operation: logical.UpdateOperation,
		Data: map[string]interface{}{
			"name":              name,
			"member_entity_ids": entityIDs,
		},
	}
}

func testScopeReq(s logical.Storage, template string) *logical.Request {
	return &logical.Request{
		Storage:   s,
		Path:      "oidc/scope/test-scope",
		Operation: logical.CreateOperation,
		Data: map[string]interface{}{
			"template": template,
		},
	}
}

func basicAuthHeader(username, password string) string {
	auth := fmt.Sprintf("%s:%s", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return fmt.Sprintf("Basic %s", encoded)
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
			"key":          "test-key-1",
			"id_token_ttl": "1m",
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
			"key":          "test-key-2",
			"id_token_ttl": "1m",
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

	// Create a test client "test-client" with an id_token_ttl longer than the
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
			"key":          "test-key1",
			"id_token_ttl": "1m",
		},
	})
	expectSuccess(t, resp, err)

	// Update the test client "test-client" -- should fail
	resp, err = c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client",
		Operation: logical.UpdateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key":          "test-key2",
			"id_token_ttl": "1m",
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
			"key":          "test-key",
			"id_token_ttl": "1m",
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
		"id_token_ttl":     int64(60),
		"access_token_ttl": int64(86400),
		"client_id":        resp.Data["client_id"],
		"client_secret":    resp.Data["client_secret"],
	}
	if diff := deep.Equal(expected, resp.Data); diff != nil {
		t.Fatal(diff)
	}
	clientID := resp.Data["client_id"].(string)
	if len(clientID) != clientIDLength {
		t.Fatalf("client_id format is incorrect: %#v", clientID)
	}
	clientSecret := resp.Data["client_secret"].(string)
	if !strings.HasPrefix(clientSecret, clientSecretPrefix) {
		t.Fatalf("client_secret format is incorrect: %#v", clientSecret)
	}
	if len(clientSecret) != clientSecretLength+len(clientSecretPrefix) {
		t.Fatalf("client_secret format is incorrect: %#v", clientSecret)
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
			"id_token_ttl":  "1m",
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
		"id_token_ttl":     int64(60),
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
			"key":          "test-key",
			"id_token_ttl": "1m",
		},
	})

	c.identityStore.HandleRequest(ctx, &logical.Request{
		Path:      "oidc/client/test-client2",
		Operation: logical.CreateOperation,
		Storage:   storage,
		Data: map[string]interface{}{
			"key":          "test-key",
			"id_token_ttl": "1m",
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
		{
			templ:         `{"auth_time": 123456, "other": "test"}`,
			restrictedKey: "auth_time",
		},
		{
			templ:         `{"at_hash": "abcdefg", "other": "test"}`,
			restrictedKey: "at_hash",
		},
		{
			templ:         `{"c_hash": "hijklmn", "other": "test"}`,
			restrictedKey: "c_hash",
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
			"top level key %q not allowed. Restricted keys: iat, aud, exp, iss, sub, namespace, nonce, auth_time, at_hash, c_hash",
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
			"key":          "test-key",
			"assignments":  []string{"test-assignment"},
			"id_token_ttl": "1m",
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

	// Populate storage with a assignment
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
