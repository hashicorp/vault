// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import (
	"net/http"
	"time"

	graphql "github.com/hasura/go-graphql-client"
)

const apiURL = "https://api.monday.com/v2"

// defaultRetryDelay is the fixed delay between attempts when a mutation fails
// with a retryable error.
const defaultRetryDelay = 5 * time.Second

// retryDelayOrDefault returns d, or defaultRetryDelay when d is zero. Requests
// carry an unexported delay override so tests can shorten the wait.
func retryDelayOrDefault(d time.Duration) time.Duration {
	if d == 0 {
		return defaultRetryDelay
	}
	return d
}

// DuplicateBoardType is the GraphQL enum for the duplicate_board mutation's
// duplicate_type argument. Implements graphql.GraphQLType so that hasura emits
// "DuplicateBoardType!" in the operation variable declaration rather than
// inferring "String!" from the underlying Go string value.
type DuplicateBoardType string

// GetGraphQLType implements graphql.GraphQLType.
func (DuplicateBoardType) GetGraphQLType() string { return "DuplicateBoardType" }

const (
	// DuplicateBoardWithStructure copies structure only (no items).
	DuplicateBoardWithStructure DuplicateBoardType = "duplicate_board_with_structure"
	// DuplicateBoardWithPulses copies structure and items.
	DuplicateBoardWithPulses DuplicateBoardType = "duplicate_board_with_pulses"
	// DuplicateBoardWithPulsesAndUpdates copies structure, items, and updates.
	DuplicateBoardWithPulsesAndUpdates DuplicateBoardType = "duplicate_board_with_pulses_and_updates"
)

// mondayAuthTransport injects the Monday.com API token into every request.
type mondayAuthTransport struct {
	token string
	base  http.RoundTripper
}

func (t *mondayAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+t.token)
	req.Header.Set("API-Version", "2024-04")
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

// NewClient returns a Monday.com GraphQL client authenticated with token.
func NewClient(token string) *graphql.Client {
	httpClient := &http.Client{
		Transport: &mondayAuthTransport{token: token},
	}
	return graphql.NewClient(apiURL, httpClient)
}
