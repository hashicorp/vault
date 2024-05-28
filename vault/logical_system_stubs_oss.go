// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

type entSystemBackend struct{}

func entUnauthenticatedPaths() []string {
	return []string{}
}

func (s *SystemBackend) entInit() {}
