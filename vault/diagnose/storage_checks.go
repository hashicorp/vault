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
	secretKey string = "diagnose"
	secretVal string = "diagnoseSecret"

	LatencyWarning    string        = "latency above 100 ms for storage call: "
	DirAccessErr      string        = "consul storage does not connect to local agent, but directly to server"
	AddrDNExistErr    string        = "config address does not exist: 127.0.0.1:8500 will be used"
	wrongRWValsPrefix string        = "Storage get and put gave wrong values: "
	latencyThreshold  time.Duration = time.Millisecond * 100
)

func EndToEndLatencyCheckWrite(ctx context.Context, uuid string, b physical.Backend) error {
	start := time.Now()
	err := b.Put(context.Background(), &physical.Entry{Key: secretKey, Value: []byte(secretVal)})
	duration := time.Since(start)
	if err != nil {
		return err
	}
	if duration > latencyThreshold {
		return fmt.Errorf(LatencyWarning + "operation: put")
	}
	return nil
}

func EndToEndLatencyCheckRead(ctx context.Context, uuid string, b physical.Backend) error {

	start := time.Now()
	val, err := b.Get(context.Background(), "diagnose")
	duration := time.Since(start)
	if err != nil {
		return err
	}
	if val.Key != "diagnose" && string(val.Value) != "diagnose" {
		return fmt.Errorf(wrongRWValsPrefix+"expecting diagnose, but got %s, %s", val.Key, val.Value)
	}
	if duration > latencyThreshold {
		return fmt.Errorf(LatencyWarning + "operation: get")
	}
	return nil
}
func EndToEndLatencyCheckDelete(ctx context.Context, uuid string, b physical.Backend) error {

	start := time.Now()
	err := b.Delete(context.Background(), "diagnose")
	duration := time.Since(start)
	if err != nil {
		return err
	}
	if duration > latencyThreshold {
		return fmt.Errorf(LatencyWarning + "operation: get")
	}
	return nil
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
