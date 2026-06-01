// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package consul

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/hashicorp/vault/sdk/helper/testcluster"
)

type ClusterStorage struct {
	// Set these after calling `NewConsulClusterStorage` but before `Start` (or
	// passing in to NewDockerCluster) to control Consul version specifically in
	// your test. Leave empty for latest OSS (defined in consulhelper.go).
	ConsulVersion    string
	ConsulEnterprise bool
	ConsulConfig     string
	cleanup          func()
	config           *Config
	started          atomic.Bool
}

var _ testcluster.ClusterStorage = &ClusterStorage{}

func NewClusterStorage() *ClusterStorage {
	return &ClusterStorage{}
}

func (s *ClusterStorage) Start(ctx context.Context, opts *testcluster.ClusterOptions) error {
	if s.started.Load() {
		return nil
	}
	prefix := ""
	if opts != nil && opts.ClusterName != "" {
		prefix = fmt.Sprintf("%s-", opts.ClusterName)
	}
	cleanup, config, err := RunContainerConfig(ctx, prefix, s.ConsulVersion, s.ConsulEnterprise, true, s.ConsulConfig)
	if err != nil {
		return err
	}
	s.cleanup = cleanup
	s.config = config
	s.started.Store(true)

	return nil
}

func (s *ClusterStorage) Cleanup() error {
	if s.cleanup != nil {
		s.cleanup()
		s.cleanup = nil
	}
	return nil
}

func (s *ClusterStorage) Opts() map[string]interface{} {
	if s.config == nil {
		return nil
	}
	return map[string]interface{}{
		"address":      s.config.ContainerHTTPAddr,
		"token":        s.config.Token,
		"max_parallel": "32",
	}
}

func (s *ClusterStorage) Type() string {
	return "consul"
}

func (s *ClusterStorage) Config() *Config {
	return s.config
}
