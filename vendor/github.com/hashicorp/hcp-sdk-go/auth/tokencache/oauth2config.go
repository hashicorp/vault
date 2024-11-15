// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tokencache

import (
	"context"

	"golang.org/x/oauth2"
)

// oAuth2Config is an interface that is used to allow mocking of oauth2.Config values
type oAuth2Config interface {
	TokenSource(ctx context.Context, t *oauth2.Token) oauth2.TokenSource
}
