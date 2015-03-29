package consul

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

func TestBackend_basic(t *testing.T) {
	config, process := testStartConsulServer(t)
	defer testStopConsulServer(t, process)

	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Backend:  Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, config),
			testAccStepWritePolicy(t, "test", testPolicy),
			testAccStepReadToken(t, "test", config),
		},
	})
}

func testStartConsulServer(t *testing.T) (map[string]interface{}, *os.Process) {
	if _, err := exec.LookPath("consul"); err != nil {
		t.Skipf("consul not found: %s", err)
	}

	td, err := ioutil.TempDir("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	tf, err := ioutil.TempFile("", "vault")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if _, err := tf.Write([]byte(strings.TrimSpace(testConsulConfig))); err != nil {
		t.Fatalf("err: %s", err)
	}
	tf.Close()

	cmd := exec.Command(
		"consul", "agent",
		"-server",
		"-bootstrap",
		"-config-file", tf.Name(),
		"-data-dir", td)
	if err := cmd.Start(); err != nil {
		t.Fatalf("error starting Consul: %s", err)
	}

	// Give Consul time to startup
	time.Sleep(2 * time.Second)

	config := map[string]interface{}{
		"address": "127.0.0.1:8500",
		"token":   "test",
	}
	return config, cmd.Process
}

func testStopConsulServer(t *testing.T, p *os.Process) {
	p.Kill()
}

func testAccPreCheck(t *testing.T) {
	if _, err := exec.LookPath("consul"); err != nil {
		t.Fatal("consul must be on PATH")
	}
}

func testAccStepConfig(
	t *testing.T, config map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "config",
		Data:      config,
	}
}

func testAccStepReadToken(
	t *testing.T, name string, conf map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      name,
		Check: func(resp *logical.Response) error {
			var d struct {
				Token string `mapstructure:"token"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[WARN] Generated token: %s", d.Token)

			// Build a client and verify that the credentials work
			config := api.DefaultConfig()
			config.Address = conf["address"].(string)
			config.Token = d.Token
			client, err := api.NewClient(config)
			if err != nil {
				return err
			}

			log.Printf("[WARN] Verifying that the generated token works...")
			_, err = client.KV().Put(&api.KVPair{
				Key:   "foo",
				Value: []byte("bar"),
			}, nil)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func testAccStepWritePolicy(t *testing.T, name string, policy string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "policy/" + name,
		Data: map[string]interface{}{
			"policy": base64.StdEncoding.EncodeToString([]byte(policy)),
		},
	}
}

const testPolicy = `
key "" {
	policy = "write"
}
`

const testConsulConfig = `
{
	"datacenter": "test",
	"acl_datacenter": "test",
	"acl_master_token": "test"
}
`
