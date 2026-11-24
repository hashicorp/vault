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

// mockSlackSendClient is a mock implementation of slackSendMsgClient
type mockSlackSendClient struct {
	postMessageFunc func(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error)
}

func (m *mockSlackSendClient) PostMessageContext(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error) {
	if m.postMessageFunc != nil {
		return m.postMessageFunc(ctx, channelID, options...)
	}
	return "", "", errors.New("mock not configured")
}

func TestSendSlackMsgReq_Run(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		req            *SendSlackMsgReq
		mockClient     *mockSlackSendClient
		expectedRes    *SendSlackMsgRes
		expectedErrMsg string
	}{
		{
			name:           "nil request",
			req:            nil,
			expectedErrMsg: "send slack message request is uninitialized",
		},
		{
			name: "empty message text",
			req: &SendSlackMsgReq{
				MessageText: "",
				ChannelID:   "C123456",
			},
			expectedErrMsg: "message text is required",
		},
		{
			name: "whitespace only message text",
			req: &SendSlackMsgReq{
				MessageText: "   ",
				ChannelID:   "C123456",
			},
			expectedErrMsg: "message text is required",
		},
		{
			name: "empty channel ID",
			req: &SendSlackMsgReq{
				MessageText: "Hello",
				ChannelID:   "",
			},
			expectedErrMsg: "channel ID is required",
		},
		{
			name: "whitespace only channel ID",
			req: &SendSlackMsgReq{
				MessageText: "Hello",
				ChannelID:   "   ",
			},
			expectedErrMsg: "channel ID is required",
		},
		{
			name: "successful message send with default emoji",
			req: &SendSlackMsgReq{
				MessageText: "Test message",
				ChannelID:   "C123456",
			},
			mockClient: &mockSlackSendClient{
				postMessageFunc: func(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error) {
					return "C123456", "1234567890.123456", nil
				},
			},
			expectedRes: &SendSlackMsgRes{
				ChannelID:       "C123456",
				ThreadTimestamp: "1234567890.123456",
			},
		},
		{
			name: "successful message send with custom emoji",
			req: &SendSlackMsgReq{
				MessageText: "Test message",
				ChannelID:   "C123456",
				Emoji:       ":rocket:",
			},
			mockClient: &mockSlackSendClient{
				postMessageFunc: func(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error) {
					return "C123456", "1234567890.123456", nil
				},
			},
			expectedRes: &SendSlackMsgRes{
				ChannelID:       "C123456",
				ThreadTimestamp: "1234567890.123456",
			},
		},
		{
			name: "successful message send with thread timestamp",
			req: &SendSlackMsgReq{
				MessageText:     "Reply message",
				ChannelID:       "C123456",
				ThreadTimestamp: "1234567890.000000",
			},
			mockClient: &mockSlackSendClient{
				postMessageFunc: func(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error) {
					return "C123456", "1234567890.123456", nil
				},
			},
			expectedRes: &SendSlackMsgRes{
				ChannelID:       "C123456",
				ThreadTimestamp: "1234567890.123456",
			},
		},
		{
			name: "whitespace in emoji is trimmed",
			req: &SendSlackMsgReq{
				MessageText: "Test message",
				ChannelID:   "C123456",
				Emoji:       "  :rocket:  ",
			},
			mockClient: &mockSlackSendClient{
				postMessageFunc: func(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error) {
					return "C123456", "1234567890.123456", nil
				},
			},
			expectedRes: &SendSlackMsgRes{
				ChannelID:       "C123456",
				ThreadTimestamp: "1234567890.123456",
			},
		},
		{
			name: "whitespace in thread timestamp is trimmed",
			req: &SendSlackMsgReq{
				MessageText:     "Reply message",
				ChannelID:       "C123456",
				ThreadTimestamp: "  1234567890.000000  ",
			},
			mockClient: &mockSlackSendClient{
				postMessageFunc: func(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error) {
					return "C123456", "1234567890.123456", nil
				},
			},
			expectedRes: &SendSlackMsgRes{
				ChannelID:       "C123456",
				ThreadTimestamp: "1234567890.123456",
			},
		},
		{
			name: "slack API error",
			req: &SendSlackMsgReq{
				MessageText: "Test message",
				ChannelID:   "C123456",
			},
			mockClient: &mockSlackSendClient{
				postMessageFunc: func(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error) {
					return "", "", errors.New("slack API error")
				},
			},
			expectedErrMsg: "failed to send slack message: slack API error",
		},
		{
			name: "empty timestamp returned",
			req: &SendSlackMsgReq{
				MessageText: "Test message",
				ChannelID:   "C123456",
			},
			mockClient: &mockSlackSendClient{
				postMessageFunc: func(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error) {
					return "C123456", "", nil
				},
			},
			expectedErrMsg: "slack API returned empty timestamp",
		},
		{
			name: "whitespace only timestamp returned",
			req: &SendSlackMsgReq{
				MessageText: "Test message",
				ChannelID:   "C123456",
			},
			mockClient: &mockSlackSendClient{
				postMessageFunc: func(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error) {
					return "C123456", "   ", nil
				},
			},
			expectedErrMsg: "slack API returned empty timestamp",
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
				assert.Equal(t, tt.expectedRes.ChannelID, res.ChannelID)
				assert.Equal(t, tt.expectedRes.ThreadTimestamp, res.ThreadTimestamp)
			}
		})
	}
}

func TestSendSlackMsgReq_Run_ContextCancellation(t *testing.T) {
	t.Parallel()

	req := &SendSlackMsgReq{
		MessageText: "Test message",
		ChannelID:   "C123456",
		Client: &mockSlackSendClient{
			postMessageFunc: func(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error) {
				return "C123456", "1234567890.123456", nil
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

func TestSendSlackMsgReq_client(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		req            *SendSlackMsgReq
		expectedErrMsg string
		shouldSucceed  bool
	}{
		{
			name: "client already set",
			req: &SendSlackMsgReq{
				Client: &mockSlackSendClient{},
			},
			shouldSucceed: true,
		},
		{
			name: "token in request",
			req: &SendSlackMsgReq{
				Token: "xoxb-test-token",
			},
			shouldSucceed: true,
		},
		{
			name: "token with whitespace",
			req: &SendSlackMsgReq{
				Token: "  xoxb-test-token  ",
			},
			shouldSucceed: true,
		},
		{
			name:           "no token provided",
			req:            &SendSlackMsgReq{},
			expectedErrMsg: "slack bot token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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

func TestSendSlackMsgRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &SendSlackMsgRes{
		ChannelID:       "C123456",
		ThreadTimestamp: "1234567890.123456",
	}

	jsonBytes, err := res.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Verify it's valid JSON
	var decoded SendSlackMsgRes
	err = json.Unmarshal(jsonBytes, &decoded)
	require.NoError(t, err)
	assert.Equal(t, res.ChannelID, decoded.ChannelID)
	assert.Equal(t, res.ThreadTimestamp, decoded.ThreadTimestamp)
}

func TestSendSlackMsgRes_ToGithubOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		res      *SendSlackMsgRes
		expected string
	}{
		{
			name: "normal timestamp",
			res: &SendSlackMsgRes{
				ChannelID:       "C123456",
				ThreadTimestamp: "1234567890.123456",
			},
			expected: "1234567890.123456",
		},
		{
			name: "empty timestamp",
			res: &SendSlackMsgRes{
				ChannelID:       "C123456",
				ThreadTimestamp: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := tt.res.ToGithubOutput()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(output))
		})
	}
}
