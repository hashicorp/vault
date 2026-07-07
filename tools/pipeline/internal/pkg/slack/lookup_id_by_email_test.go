// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package slack

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	slackapi "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSlackClient is a mock implementation of slackLookupByEmailClient
type mockSlackClient struct {
	getUserByEmailFunc func(ctx context.Context, email string) (*slackapi.User, error)
}

func (m *mockSlackClient) GetUserByEmailContext(ctx context.Context, email string) (*slackapi.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, email)
	}
	return nil, errors.New("mock not configured")
}

func TestLookupIDByEmailReq_Run(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		req            *LookupIDByEmailReq
		mockClient     *mockSlackClient
		expectedRes    *LookupIDByEmailRes
		expectedErrMsg string
	}{
		{
			name:           "nil request",
			req:            nil,
			expectedErrMsg: "lookup id by email request is uninitialized",
		},
		{
			name: "empty emails list",
			req: &LookupIDByEmailReq{
				Emails: []string{},
			},
			expectedErrMsg: "at least one email is required",
		},
		{
			name: "single valid email with default domain",
			req: &LookupIDByEmailReq{
				Emails: []string{"user@ibm.com"},
			},
			mockClient: &mockSlackClient{
				getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
					return &slackapi.User{ID: "U12345"}, nil
				},
			},
			expectedRes: &LookupIDByEmailRes{
				EmailSlackIDMap: map[string]string{
					"user@ibm.com": "U12345",
				},
			},
		},
		{
			name: "multiple valid emails",
			req: &LookupIDByEmailReq{
				Emails: []string{"user1@ibm.com", "user2@ibm.com"},
			},
			mockClient: &mockSlackClient{
				getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
					if email == "user1@ibm.com" {
						return &slackapi.User{ID: "U11111"}, nil
					}
					return &slackapi.User{ID: "U22222"}, nil
				},
			},
			expectedRes: &LookupIDByEmailRes{
				EmailSlackIDMap: map[string]string{
					"user1@ibm.com": "U11111",
					"user2@ibm.com": "U22222",
				},
			},
		},
		{
			name: "custom allowed domain",
			req: &LookupIDByEmailReq{
				Emails:        []string{"user@example.com"},
				AllowedDomain: "example.com",
			},
			mockClient: &mockSlackClient{
				getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
					return &slackapi.User{ID: "U99999"}, nil
				},
			},
			expectedRes: &LookupIDByEmailRes{
				EmailSlackIDMap: map[string]string{
					"user@example.com": "U99999",
				},
			},
		},
		{
			name: "invalid domain",
			req: &LookupIDByEmailReq{
				Emails:        []string{"user@wrong.com"},
				AllowedDomain: "ibm.com",
			},
			mockClient: &mockSlackClient{
				getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
					return &slackapi.User{ID: "U12345"}, nil
				},
			},
			expectedErrMsg: `email "user@wrong.com" must use the @ibm.com domain`,
		},
		{
			name: "email with whitespace",
			req: &LookupIDByEmailReq{
				Emails: []string{"  user@ibm.com  "},
			},
			mockClient: &mockSlackClient{
				getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
					return &slackapi.User{ID: "U12345"}, nil
				},
			},
			expectedRes: &LookupIDByEmailRes{
				EmailSlackIDMap: map[string]string{
					"user@ibm.com": "U12345",
				},
			},
		},
		{
			name: "empty email strings filtered out",
			req: &LookupIDByEmailReq{
				Emails: []string{"", "  ", "user@ibm.com"},
			},
			mockClient: &mockSlackClient{
				getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
					return &slackapi.User{ID: "U12345"}, nil
				},
			},
			expectedRes: &LookupIDByEmailRes{
				EmailSlackIDMap: map[string]string{
					"user@ibm.com": "U12345",
				},
			},
		},
		{
			name: "all empty emails",
			req: &LookupIDByEmailReq{
				Emails: []string{"", "  ", "   "},
			},
			mockClient: &mockSlackClient{
				getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
					return &slackapi.User{ID: "U12345"}, nil
				},
			},
			expectedErrMsg: "at least one non-empty email is required",
		},
		{
			name: "slack API error",
			req: &LookupIDByEmailReq{
				Emails: []string{"user@ibm.com"},
			},
			mockClient: &mockSlackClient{
				getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
					return nil, errors.New("slack API error")
				},
			},
			expectedErrMsg: "looking up slack user id for user@ibm.com: slack API error",
		},
		{
			name: "empty user ID returned",
			req: &LookupIDByEmailReq{
				Emails: []string{"user@ibm.com"},
			},
			mockClient: &mockSlackClient{
				getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
					return &slackapi.User{ID: ""}, nil
				},
			},
			expectedErrMsg: "slack lookup returned empty user id for user@ibm.com",
		},
		{
			name: "case insensitive domain check",
			req: &LookupIDByEmailReq{
				Emails:        []string{"user@IBM.COM"},
				AllowedDomain: "ibm.com",
			},
			mockClient: &mockSlackClient{
				getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
					return &slackapi.User{ID: "U12345"}, nil
				},
			},
			expectedRes: &LookupIDByEmailRes{
				EmailSlackIDMap: map[string]string{
					"user@IBM.COM": "U12345",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.req != nil && tt.mockClient != nil {
				tt.req.Client = tt.mockClient
			}

			ctx := context.Background()
			res, err := tt.req.Run(ctx)

			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				assert.Equal(t, tt.expectedRes.EmailSlackIDMap, res.EmailSlackIDMap)
			}
		})
	}
}

func TestLookupIDByEmailReq_Run_ContextCancellation(t *testing.T) {
	t.Parallel()

	req := &LookupIDByEmailReq{
		Emails: []string{"user@ibm.com"},
		Client: &mockSlackClient{
			getUserByEmailFunc: func(ctx context.Context, email string) (*slackapi.User, error) {
				return &slackapi.User{ID: "U12345"}, nil
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	res, err := req.Run(ctx)
	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Nil(t, res)
}

func TestLookupIDByEmailReq_client(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		req            *LookupIDByEmailReq
		envToken       string
		expectedErrMsg string
		shouldSucceed  bool
	}{
		{
			name: "client already set",
			req: &LookupIDByEmailReq{
				Client: &mockSlackClient{},
			},
			shouldSucceed: true,
		},
		{
			name: "token in request",
			req: &LookupIDByEmailReq{
				Token: "xoxb-test-token",
			},
			shouldSucceed: true,
		},
		{
			name: "token with whitespace",
			req: &LookupIDByEmailReq{
				Token: "  xoxb-test-token  ",
			},
			shouldSucceed: true,
		},
		{
			name:           "no token provided",
			req:            &LookupIDByEmailReq{},
			expectedErrMsg: "slack token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Note: We can't easily test environment variable reading in parallel tests
			// without risking race conditions, so we skip that scenario
			client, err := tt.req.client()

			if tt.expectedErrMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
				assert.Nil(t, client)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestLookupIDByEmailRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &LookupIDByEmailRes{
		EmailSlackIDMap: map[string]string{
			"user1@ibm.com": "U11111",
			"user2@ibm.com": "U22222",
		},
	}

	jsonBytes, err := res.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Verify it's valid JSON
	var decoded LookupIDByEmailRes
	err = json.Unmarshal(jsonBytes, &decoded)
	require.NoError(t, err)
	assert.Equal(t, res.EmailSlackIDMap, decoded.EmailSlackIDMap)
}

func TestLookupIDByEmailRes_ToGithubOutput(t *testing.T) {
	t.Parallel()

	res := &LookupIDByEmailRes{
		EmailSlackIDMap: map[string]string{
			"user1@ibm.com": "U11111",
			"user2@ibm.com": "U22222",
		},
	}

	jsonBytes, err := res.ToGithubOutput()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Verify it's valid JSON and contains only the map
	var decoded map[string]string
	err = json.Unmarshal(jsonBytes, &decoded)
	require.NoError(t, err)
	assert.Equal(t, res.EmailSlackIDMap, decoded)
}

func TestLookupIDByEmailRes_Emails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		res           *LookupIDByEmailRes
		expectedCount int
	}{
		{
			name: "multiple emails",
			res: &LookupIDByEmailRes{
				EmailSlackIDMap: map[string]string{
					"user3@ibm.com": "U33333",
					"user1@ibm.com": "U11111",
					"user2@ibm.com": "U22222",
				},
			},
			expectedCount: 3,
		},
		{
			name: "single email",
			res: &LookupIDByEmailRes{
				EmailSlackIDMap: map[string]string{
					"user@ibm.com": "U12345",
				},
			},
			expectedCount: 1,
		},
		{
			name: "empty map",
			res: &LookupIDByEmailRes{
				EmailSlackIDMap: map[string]string{},
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			emails := tt.res.Emails()
			assert.Len(t, emails, tt.expectedCount)

			// Verify emails are sorted
			if len(emails) > 1 {
				for i := 1; i < len(emails); i++ {
					assert.True(t, emails[i-1] < emails[i], "emails should be sorted")
				}
			}

			// Verify all emails from map are present
			for email := range tt.res.EmailSlackIDMap {
				assert.Contains(t, emails, email)
			}
		})
	}
}

func TestLookupIDByEmailRes_Emails_Sorted(t *testing.T) {
	t.Parallel()

	res := &LookupIDByEmailRes{
		EmailSlackIDMap: map[string]string{
			"zebra@ibm.com": "U99999",
			"alpha@ibm.com": "U11111",
			"beta@ibm.com":  "U22222",
		},
	}

	emails := res.Emails()
	require.Len(t, emails, 3)
	assert.Equal(t, "alpha@ibm.com", emails[0])
	assert.Equal(t, "beta@ibm.com", emails[1])
	assert.Equal(t, "zebra@ibm.com", emails[2])
}
