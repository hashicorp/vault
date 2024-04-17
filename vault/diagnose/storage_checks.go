// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package diagnose

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/physical"
)

const (
	success   string = "success"
	secretVal string = "diagnoseSecret"

	LatencyWarning    string        = "Latency above 100 ms: "
	DirAccessErr      string        = "Vault storage is directly connected to a Consul server."
	DirAccessAdvice   string        = "We recommend connecting to a local agent."
	AddrDNExistErr    string        = "Storage config address does not exist: 127.0.0.1:8500 will be used."
	wrongRWValsPrefix string        = "Storage get and put gave wrong values: "
	latencyThreshold  time.Duration = time.Millisecond * 100
)

func EndToEndLatencyCheckWrite(ctx context.Context, uuid string, b physical.Backend) (time.Duration, error) {
	start := time.Now()
	err := b.Put(context.Background(), &physical.Entry{Key: uuid, Value: []byte(secretVal)})
	duration := time.Since(start)
	if err != nil {
		return time.Duration(0), err
	}
	if duration > latencyThreshold {
		return duration, nil
	}
	return time.Duration(0), nil
}

func EndToEndLatencyCheckRead(ctx context.Context, uuid string, b physical.Backend) (time.Duration, error) {
	start := time.Now()
	val, err := b.Get(context.Background(), uuid)
	duration := time.Since(start)
	if err != nil {
		return time.Duration(0), err
	}
	if val == nil {
		return time.Duration(0), fmt.Errorf("No value found when reading generated data.")
	}
	if val.Key != uuid && string(val.Value) != secretVal {
		return time.Duration(0), fmt.Errorf(wrongRWValsPrefix+"expecting %s as key and diagnose for value, but got %s, %s.", uuid, val.Key, val.Value)
	}
	if duration > latencyThreshold {
		return duration, nil
	}
	return time.Duration(0), nil
}

func EndToEndLatencyCheckDelete(ctx context.Context, uuid string, b physical.Backend) (time.Duration, error) {
	start := time.Now()
	err := b.Delete(context.Background(), uuid)
	duration := time.Since(start)
	if err != nil {
		return time.Duration(0), err
	}
	if duration > latencyThreshold {
		return duration, nil
	}
	return time.Duration(0), nil
}

// ConsulDirectAccess verifies that consul is connecting to local agent,
// versus directly to a remote server. We can only assume that the local address
// is a server, not a client.
func ConsulDirectAccess(config map[string]string) string {
	configAddr, ok := config["address"]
	if !ok {
		return AddrDNExistErr
	}
	if !strings.Contains(configAddr, "localhost") && !strings.Contains(configAddr, "127.0.0.1") {
		return DirAccessErr
	}
	return ""
}
