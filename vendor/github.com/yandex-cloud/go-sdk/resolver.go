// Copyright (c) 2018 Yandex LLC. All rights reserved.
// Author: Dmitry Novikov <novikoff@yandex-team.ru>

package ycsdk

import (
	"context"

	"google.golang.org/grpc"
)

type Resolver interface {
	ID() string
	Err() error

	Run(context.Context, *SDK, ...grpc.CallOption) error
}
