// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	slackapi "github.com/slack-go/slack"
)

type LookupIDByEmailReq struct {
	Emails              []string
	Token               string
	WriteToGithubOutput bool
	AllowedDomain       string
	Client              slackLookupByEmailClient
}

type LookupIDByEmailRes struct {
	EmailSlackIDMap map[string]string `json:"email_slack_id_map,omitempty"`
}

type slackLookupByEmailClient interface {
	GetUserByEmailContext(ctx context.Context, email string) (*slackapi.User, error)
}

func (r *LookupIDByEmailReq) Run(ctx context.Context) (*LookupIDByEmailRes, error) {
	if r == nil {
		return nil, fmt.Errorf("lookup id by email request is uninitialized")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if len(r.Emails) == 0 {
		return nil, fmt.Errorf("at least one email is required")
	}

	allowedDomain := strings.TrimSpace(r.AllowedDomain)
	if allowedDomain == "" {
		allowedDomain = "ibm.com"
	}

	client, err := r.client()
	if err != nil {
		return nil, err
	}

	res := &LookupIDByEmailRes{
		EmailSlackIDMap: map[string]string{},
	}

	for _, email := range r.Emails {
		email = strings.TrimSpace(email)
		if email == "" {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(email), "@"+strings.ToLower(allowedDomain)) {
			return nil, fmt.Errorf("email %q must use the @%s domain", email, allowedDomain)
		}

		slog.Default().DebugContext(ctx, "looking up slack user id by email", "email", email)

		user, err := client.GetUserByEmailContext(ctx, email)
		if err != nil {
			return nil, fmt.Errorf("looking up slack user id for %s: %w", email, err)
		}
		if strings.TrimSpace(user.ID) == "" {
			return nil, fmt.Errorf("slack lookup returned empty user id for %s", email)
		}

		res.EmailSlackIDMap[email] = user.ID
	}

	if len(res.EmailSlackIDMap) == 0 {
		return nil, fmt.Errorf("at least one non-empty email is required")
	}

	return res, nil
}

func (r *LookupIDByEmailReq) client() (slackLookupByEmailClient, error) {
	if r.Client != nil {
		return r.Client, nil
	}

	token := strings.TrimSpace(r.Token)
	if token == "" {
		token = strings.TrimSpace(os.Getenv("SLACK_TOKEN"))
	}
	if token == "" {
		return nil, fmt.Errorf("slack token is required")
	}

	return slackapi.New(token), nil
}

func (r *LookupIDByEmailRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling slack lookup id by email response to JSON: %w", err)
	}

	return b, nil
}

func (r *LookupIDByEmailRes) ToGithubOutput() ([]byte, error) {
	b, err := json.Marshal(r.EmailSlackIDMap)
	if err != nil {
		return nil, fmt.Errorf("marshaling slack lookup id by email github output to JSON: %w", err)
	}

	return b, nil
}

func (r *LookupIDByEmailRes) Emails() []string {
	emails := make([]string, 0, len(r.EmailSlackIDMap))
	for email := range r.EmailSlackIDMap {
		emails = append(emails, email)
	}
	slices.Sort(emails)

	return emails
}
