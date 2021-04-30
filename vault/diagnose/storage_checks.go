package diagnose

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/physical"
)

var success string = "success"
var secretKey string = "diagnose"
var secretVal string = "diagnoseSecret"

const timeOutErr string = "storage call timed out after 20 seconds: "
const wrongRWValsPrefix string = "Storage get and put gave wrong values: "

// StorageEndToEndLatencyCheck calls Write, Read, and Delete on a secret in the root
// directory of the backend.
// Note: Just checking read, write, and delete for root. It's a very basic check,
// but I don't think we can necessarily do any more than that. We could check list,
// but I don't think List is ever going to break in isolation.
func StorageEndToEndLatencyCheck(ctx context.Context, b physical.Backend) error {

	c2 := make(chan string, 1)
	go func() {
		err := b.Put(context.Background(), &physical.Entry{Key: secretKey, Value: []byte(secretVal)})
		if err != nil {
			c2 <- err.Error()
		} else {
			c2 <- success
		}
	}()
	select {
	case errString := <-c2:
		if errString != success {
			return fmt.Errorf(errString)
		}
	case <-time.After(20 * time.Second):
		return fmt.Errorf(timeOutErr + "operation: Put")
	}

	c3 := make(chan *physical.Entry)
	c4 := make(chan error)
	go func() {
		val, err := b.Get(context.Background(), "diagnose")
		if err != nil {
			c4 <- err
		} else {
			c3 <- val
		}
	}()
	select {
	case err := <-c4:
		return err
	case val := <-c3:
		if val.Key != "diagnose" && string(val.Value) != "diagnose" {
			return fmt.Errorf(wrongRWValsPrefix+"expecting diagnose, but got %s, %s", val.Key, val.Value)
		}
	case <-time.After(20 * time.Second):
		return fmt.Errorf(timeOutErr + "operation: Get")
	}

	c5 := make(chan string, 1)
	go func() {
		err := b.Delete(context.Background(), "diagnose")
		if err != nil {
			c5 <- err.Error()
		} else {
			c5 <- success
		}
	}()
	select {
	case errString := <-c5:
		if errString != success {
			return fmt.Errorf(errString)
		}
	case <-time.After(20 * time.Second):
		return fmt.Errorf(timeOutErr + "operation: Delete")
	}
	return nil
}

// ConsulDirectAccess verifies that consul is connecting to local agent,
// versus directly to a remote server. We can only assume that the local address
// is a server, not a client.
func ConsulDirectAccess(config map[string]string) string {
	serviceRegistrationAddr := config["address"]
	if !strings.Contains(serviceRegistrationAddr, "localhost") && !strings.Contains(serviceRegistrationAddr, "127.0.0.1") {
		return "consul storage does not connect to local agent, but directly to server"
	}
	return ""
}
