package linodego

import (
	"context"
)

// ChildAccount represents an account under the current account.
// NOTE: This is an alias to prevent any future breaking changes.
type ChildAccount = Account

// ChildAccountToken represents a short-lived token created using
// the CreateChildAccountToken(...) function.
// NOTE: This is an alias to prevent any future breaking changes.
type ChildAccountToken = Token

// ListChildAccounts lists child accounts under the current account.
// NOTE: Parent/Child related features may not be generally available.
func (c *Client) ListChildAccounts(ctx context.Context, opts *ListOptions) ([]ChildAccount, error) {
	return getPaginatedResults[ChildAccount](
		ctx,
		c,
		"account/child-accounts",
		opts,
	)
}

// GetChildAccount gets a single child accounts under the current account.
// NOTE: Parent/Child related features may not be generally available.
func (c *Client) GetChildAccount(ctx context.Context, euuid string) (*ChildAccount, error) {
	return doGETRequest[ChildAccount](
		ctx,
		c,
		formatAPIPath("account/child-accounts/%s", euuid),
	)
}

// CreateChildAccountToken creates a short-lived token that can be used to
// access the Linode API under a child account.
// The attributes of this token are not currently configurable.
// NOTE: Parent/Child related features may not be generally available.
func (c *Client) CreateChildAccountToken(ctx context.Context, euuid string) (*ChildAccountToken, error) {
	return doPOSTRequest[ChildAccountToken, any](
		ctx,
		c,
		formatAPIPath("account/child-accounts/%s/token", euuid),
	)
}
