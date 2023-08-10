// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package etcd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	"go.etcd.io/etcd/client/v2"
)

var (
	EtcdMultipleBootstrapError = errors.New("client setup failed: multiple discovery or bootstrap flags specified, use either \"address\" or \"discovery_srv\"")
	EtcdAddressError           = errors.New("client setup failed: address must be valid URL (ex. 'scheme://host:port')")
	EtcdLockHeldError          = errors.New("lock already held")
	EtcdLockNotHeldError       = errors.New("lock not held")
	EtcdVersionUnknown         = errors.New("etcd: unknown API version")
)

// NewEtcdBackend constructs a etcd backend using a given machine address.
func NewEtcdBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	var (
		apiVersion string
		ok         bool
	)

	if apiVersion, ok = conf["etcd_api"]; !ok {
		apiVersion = os.Getenv("ETCD_API")
	}

	if apiVersion == "" {
		apiVersion = "v3"
	}

	switch apiVersion {
	case "3", "etcd3", "v3":
		return newEtcd3Backend(conf, logger)
	default:
		return nil, EtcdVersionUnknown
	}
}

// Retrieves the config option in order of priority:
//  1. The named environment variable if it exist
//  2. The key in the config map
func getEtcdOption(conf map[string]string, confKey, envVar string) (string, bool) {
	confVal, inConf := conf[confKey]
	envVal, inEnv := os.LookupEnv(envVar)
	if inEnv {
		return envVal, true
	}
	return confVal, inConf
}

func getEtcdEndpoints(conf map[string]string) ([]string, error) {
	address, staticBootstrap := getEtcdOption(conf, "address", "ETCD_ADDR")
	domain, useSrv := getEtcdOption(conf, "discovery_srv", "ETCD_DISCOVERY_SRV")
	if useSrv && staticBootstrap {
		return nil, EtcdMultipleBootstrapError
	}

	if staticBootstrap {
		endpoints := strings.Split(address, ",")
		// Verify that the machines are valid URLs
		for _, e := range endpoints {
			u, urlErr := url.Parse(e)
			if urlErr != nil || u.Scheme == "" {
				return nil, EtcdAddressError
			}
		}
		return endpoints, nil
	}

	if useSrv {
		srvName, _ := getEtcdOption(conf, "discovery_srv_name", "ETCD_DISCOVERY_SRV_NAME")
		discoverer := client.NewSRVDiscover()
		endpoints, err := discoverer.Discover(domain, srvName)
		if err != nil {
			return nil, fmt.Errorf("failed to discover etcd endpoints through SRV discovery: %w", err)
		}
		return endpoints, nil
	}

	// Set a default endpoints list if no option was set
	return []string{"http://127.0.0.1:2379"}, nil
}
