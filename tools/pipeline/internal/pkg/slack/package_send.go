// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	slackapi "github.com/slack-go/slack"
)

type SendSlackMsgReq struct {
	MessageText         string
	ChannelID           string
	Token               string
	Emoji               string
	ThreadTimestamp     string
	WriteToGithubOutput bool
	Client              slackSendMsgClient
}

type SendSlackMsgRes struct {
	ChannelID       string `json:"channel_id,omitempty"`
	ThreadTimestamp string `json:"thread_timestamp,omitempty"`
}

type slackSendMsgClient interface {
	PostMessageContext(ctx context.Context, channelID string, options ...slackapi.MsgOption) (string, string, error)
}

func (r *SendSlackMsgReq) Run(ctx context.Context) (*SendSlackMsgRes, error) {
	if r == nil {
		return nil, fmt.Errorf("send slack message request is uninitialized")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if strings.TrimSpace(r.MessageText) == "" {
		return nil, fmt.Errorf("message text is required")
	}
	if strings.TrimSpace(r.ChannelID) == "" {
		return nil, fmt.Errorf("channel ID is required")
	}

	client, err := r.client()
	if err != nil {
		return nil, err
	}

	emoji := strings.TrimSpace(r.Emoji)
	if emoji == "" {
		emoji = ":vault:"
	}

	options := []slackapi.MsgOption{
		slackapi.MsgOptionText(r.MessageText, false),
		slackapi.MsgOptionIconEmoji(emoji),
	}

	threadTS := strings.TrimSpace(r.ThreadTimestamp)
	if threadTS != "" {
		options = append(options, slackapi.MsgOptionTS(threadTS))
	}

	slog.Default().DebugContext(ctx, "sending slack message", "channel", r.ChannelID, "thread_ts", threadTS)

	channelID, ts, err := client.PostMessageContext(ctx, r.ChannelID, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to send slack message: %w", err)
	}

	if strings.TrimSpace(ts) == "" {
		return nil, fmt.Errorf("slack API returned empty timestamp")
	}

	return &SendSlackMsgRes{
		ChannelID:       channelID,
		ThreadTimestamp: ts,
	}, nil
}

func (r *SendSlackMsgReq) client() (slackSendMsgClient, error) {
	if r.Client != nil {
		return r.Client, nil
	}

	token := strings.TrimSpace(r.Token)
	if token == "" {
		token = strings.TrimSpace(os.Getenv("SLACK_BOT_TOKEN"))
	}
	if token == "" {
		return nil, fmt.Errorf("slack bot token is required")
	}

	return slackapi.New(token), nil
}

func (r *SendSlackMsgRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling slack send message response to JSON: %w", err)
	}

	return b, nil
}

func (r *SendSlackMsgRes) ToGithubOutput() ([]byte, error) {
	return []byte(r.ThreadTimestamp), nil
}
