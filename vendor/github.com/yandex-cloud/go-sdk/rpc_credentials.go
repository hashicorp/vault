// Copyright (c) 2018 Yandex LLC. All rights reserved.
// Author: Maxim Kolganov <manykey@yandex-team.ru>

package ycsdk

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
	"github.com/yandex-cloud/go-sdk/pkg/sdkerrors"
)

type rpcCredentials struct {
	plaintext bool

	createToken createIAMTokenFunc // Injected on Init
	// now may be replaced in tests
	now func() time.Time

	// mutex guards conn and currentState, and excludes multiple simultaneous token updates
	mutex        sync.RWMutex
	currentState rpcCredentialsState
}

var _ credentials.PerRPCCredentials = &rpcCredentials{}

type rpcCredentialsState struct {
	token     string
	expiresAt time.Time
	version   int64
}

func newRPCCredentials(plaintext bool) *rpcCredentials {
	return &rpcCredentials{
		plaintext: plaintext,
		now:       time.Now,
	}
}

type createIAMTokenFunc func(ctx context.Context) (*iam.CreateIamTokenResponse, error)

func (c *rpcCredentials) Init(createToken createIAMTokenFunc) {
	c.createToken = createToken
}

func (c *rpcCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	audienceURL, err := url.Parse(uri[0])
	if err != nil {
		return nil, err
	}
	if audienceURL.Path == "/yandex.cloud.iam.v1.IamTokenService" ||
		audienceURL.Path == "/yandex.cloud.endpoint.ApiEndpointService" {
		return nil, nil
	}

	c.mutex.RLock()
	state := c.currentState
	c.mutex.RUnlock()

	token := state.token
	expired := c.now().After(state.expiresAt)
	if expired {
		token, err = c.updateToken(ctx, state)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unauthenticated {
				return nil, err
			}
			return nil, status.Errorf(codes.Unauthenticated, "%v", err)
		}
	}

	return map[string]string{
		"authorization": "Bearer " + token,
	}, nil
}

func (c *rpcCredentials) RequireTransportSecurity() bool {
	return !c.plaintext
}

func (c *rpcCredentials) updateToken(ctx context.Context, currentState rpcCredentialsState) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.currentState.version != currentState.version {
		// someone have already updated it
		return c.currentState.token, nil
	}

	resp, err := c.createToken(ctx)
	if err != nil {
		return "", sdkerrors.WithMessage(err, "iam token create failed")
	}
	expiresAt, expiresAtErr := ptypes.Timestamp(resp.ExpiresAt)
	if expiresAtErr != nil {
		grpclog.Warningf("invalid IAM Token expires_at: %s", expiresAtErr)
		// Fallback to short term caching.
		expiresAt = time.Now().Add(time.Minute)
	}
	c.currentState = rpcCredentialsState{
		token:     resp.IamToken,
		expiresAt: expiresAt,
		version:   currentState.version + 1,
	}
	return c.currentState.token, nil
}
