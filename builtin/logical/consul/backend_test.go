package consul

import (
	"bufio"
	"encoding/base64"
	"fmt"
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
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	config, process := testStartConsulServer(t)
	defer testStopConsulServer(t, process)

	b, _ := Factory(logical.TestBackendConfig())
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, config),
			testAccStepWritePolicy(t, "test", testPolicy, ""),
			testAccStepReadToken(t, "test", config),
		},
	})
}

func TestBackend_management(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	config, process := testStartConsulServer(t)
	defer testStopConsulServer(t, process)

	b, _ := Factory(logical.TestBackendConfig())
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, config),
			testAccStepWriteManagementPolicy(t, "test", ""),
			testAccStepReadManagementToken(t, "test", config),
		},
	})
}

func TestBackend_crud(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	_, process := testStartConsulServer(t)
	defer testStopConsulServer(t, process)

	b, _ := Factory(logical.TestBackendConfig())
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", testPolicy, ""),
			testAccStepReadPolicy(t, "test", testPolicy, 0),
			testAccStepDeletePolicy(t, "test"),
		},
	})
}

func TestBackend_role_lease(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", logicaltest.TestEnvVar))
		return
	}

	_, process := testStartConsulServer(t)
	defer testStopConsulServer(t, process)

	b, _ := Factory(logical.TestBackendConfig())
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		PreCheck:       func() { testAccPreCheck(t) },
		Backend:        b,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", testPolicy, "6h"),
			testAccStepReadPolicy(t, "test", testPolicy, 6*time.Hour),
			testAccStepDeletePolicy(t, "test"),
		},
	})
}

func testStartConsulServer(t *testing.T) (map[string]interface{}, *os.Process) {
	if _, err := exec.LookPath("consul"); err != nil {
		t.Errorf("consul not found: %s", err)
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
		"-advertise", "127.0.0.1",
		"-config-file", tf.Name(),
		"-data-dir", td)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)
	stdoutScanFunc := func() {
		for stdoutScanner.Scan() {
			t.Logf("Consul stdout: %s\n", stdoutScanner.Text())
		}
	}
	stderrScanFunc := func() {
		for stderrScanner.Scan() {
			t.Logf("Consul stderr: %s\n", stderrScanner.Text())
		}
	}
	if os.Getenv("VAULT_VERBOSE_ACC_TESTS") != "" {
		go stdoutScanFunc()
		go stderrScanFunc()
	}

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
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      config,
	}
}

func testAccStepReadToken(
	t *testing.T, name string, conf map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "creds/" + name,
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

func testAccStepReadManagementToken(
	t *testing.T, name string, conf map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "creds/" + name,
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
			_, _, err = client.ACL().Create(&api.ACLEntry{
				Type: "management",
				Name: "test2",
			}, nil)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func testAccStepWritePolicy(t *testing.T, name string, policy string, lease string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"policy": base64.StdEncoding.EncodeToString([]byte(policy)),
			"lease":  lease,
		},
	}
}

func testAccStepWriteManagementPolicy(t *testing.T, name string, lease string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"token_type": "management",
			"lease":      lease,
		},
	}
}

func testAccStepReadPolicy(t *testing.T, name string, policy string, lease time.Duration) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			policyRaw := resp.Data["policy"].(string)
			out, err := base64.StdEncoding.DecodeString(policyRaw)
			if err != nil {
				return err
			}
			if string(out) != policy {
				return fmt.Errorf("mismatch: %s %s", out, policy)
			}

			leaseRaw := resp.Data["lease"].(string)
			l, err := time.ParseDuration(leaseRaw)
			if err != nil {
				return err
			}
			if l != lease {
				return fmt.Errorf("mismatch: %v %v", l, lease)
			}
			return nil
		},
	}
}

func testAccStepDeletePolicy(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + name,
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
