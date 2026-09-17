// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

type HealthCheckManager struct{}

func NewHealthCheckManager(core *Core) *HealthCheckManager {
	return nil
}
