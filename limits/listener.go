// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package limits

type RequestListener struct{}

func (l *RequestListener) OnSuccess() {}

func (l *RequestListener) OnDropped() {}

func (l *RequestListener) OnIgnore() {}
