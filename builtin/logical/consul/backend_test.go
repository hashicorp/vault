package consul

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
	dockertest "gopkg.in/ory-am/dockertest.v2"
)

var (
	testImagePull sync.Once
)

func prepareTestContainer(t *testing.T, s logical.Storage, b logical.Backend) (cid dockertest.ContainerID, retAddress string) {
	if os.Getenv("CONSUL_ADDR") != "" {
		return "", os.Getenv("CONSUL_ADDR")
	}

	// Without this the checks for whether the container has started seem to
	// never actually pass. There's really no reason to expose the test
	// containers, so don't.
	dockertest.BindDockerToLocalhost = "yep"

	testImagePull.Do(func() {
		dockertest.Pull(dockertest.ConsulImageName)
	})

	try := 0
	cid, connErr := dockertest.ConnectToConsul(60, 500*time.Millisecond, func(connAddress string) bool {
		try += 1
		// Build a client and verify that the credentials work
		config := consulapi.DefaultConfig()
		config.Address = connAddress
		config.Token = dockertest.ConsulACLMasterToken
		client, err := consulapi.NewClient(config)
		if err != nil {
			if try > 50 {
				panic(err)
			}
			return false
		}

		_, err = client.KV().Put(&consulapi.KVPair{
			Key:   "setuptest",
			Value: []byte("setuptest"),
		}, nil)
		if err != nil {
			if try > 50 {
				panic(err)
			}
			return false
		}

		retAddress = connAddress
		return true
	})

	if connErr != nil {
		t.Fatalf("could not connect to consul: %v", connErr)
	}

	return
}

func cleanupTestContainer(t *testing.T, cid dockertest.ContainerID) {
	err := cid.KillRemove()
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_config_access(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, connURL := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}
	connData := map[string]interface{}{
		"address": connURL,
		"token":   dockertest.ConsulACLMasterToken,
	}

	confReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Storage:   config.StorageView,
		Data:      connData,
	}

	resp, err := b.HandleRequest(confReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	confReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(confReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	expected := map[string]interface{}{
		"address": connData["address"].(string),
		"scheme":  "http",
	}
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
	}
	if resp.Data["token"] != nil {
		t.Fatalf("token should not be set in the response")
	}
}

func TestBackend_basic(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, connURL := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}
	connData := map[string]interface{}{
		"address": connURL,
		"token":   dockertest.ConsulACLMasterToken,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData),
			testAccStepWritePolicy(t, "test", testPolicy, ""),
			testAccStepReadToken(t, "test", connData),
		},
	})
}

func TestBackend_renew_revoke(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, connURL := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}
	connData := map[string]interface{}{
		"address": connURL,
		"token":   dockertest.ConsulACLMasterToken,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "roles/test"
	req.Data = map[string]interface{}{
		"policy": base64.StdEncoding.EncodeToString([]byte(testPolicy)),
		"lease":  "6h",
	}
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test"
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
	if resp.IsError() {
		t.Fatalf("resp is error: %v", resp.Error())
	}

	generatedSecret := resp.Secret
	generatedSecret.IssueTime = time.Now()
	generatedSecret.TTL = 6 * time.Hour

	var d struct {
		Token string `mapstructure:"token"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}
	log.Printf("[WARN] Generated token: %s", d.Token)

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("[WARN] Verifying that the generated token works...")
	_, err = client.KV().Put(&consulapi.KVPair{
		Key:   "foo",
		Value: []byte("bar"),
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.RenewOperation
	req.Secret = generatedSecret
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}

	req.Operation = logical.RevokeOperation
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("[WARN] Verifying that the generated token does not work...")
	_, err = client.KV().Put(&consulapi.KVPair{
		Key:   "foo",
		Value: []byte("bar"),
	}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBackend_management(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	cid, connURL := prepareTestContainer(t, config.StorageView, b)
	if cid != "" {
		defer cleanupTestContainer(t, cid)
	}
	connData := map[string]interface{}{
		"address": connURL,
		"token":   dockertest.ConsulACLMasterToken,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData),
			testAccStepWriteManagementPolicy(t, "test", ""),
			testAccStepReadManagementToken(t, "test", connData),
		},
	})
}

func TestBackend_crud(t *testing.T) {
	b, _ := Factory(logical.TestBackendConfig())
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", testPolicy, ""),
			testAccStepWritePolicy(t, "test2", testPolicy, ""),
			testAccStepWritePolicy(t, "test3", testPolicy, ""),
			testAccStepReadPolicy(t, "test", testPolicy, 0),
			testAccStepListPolicy(t, []string{"test", "test2", "test3"}),
			testAccStepDeletePolicy(t, "test"),
		},
	})
}

func TestBackend_role_lease(t *testing.T) {
	b, _ := Factory(logical.TestBackendConfig())
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", testPolicy, "6h"),
			testAccStepReadPolicy(t, "test", testPolicy, 6*time.Hour),
			testAccStepDeletePolicy(t, "test"),
		},
	})
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
			config := consulapi.DefaultConfig()
			config.Address = conf["address"].(string)
			config.Token = d.Token
			client, err := consulapi.NewClient(config)
			if err != nil {
				return err
			}

			log.Printf("[WARN] Verifying that the generated token works...")
			_, err = client.KV().Put(&consulapi.KVPair{
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
			config := consulapi.DefaultConfig()
			config.Address = conf["address"].(string)
			config.Token = d.Token
			client, err := consulapi.NewClient(config)
			if err != nil {
				return err
			}

			log.Printf("[WARN] Verifying that the generated token works...")
			_, _, err = client.ACL().Create(&consulapi.ACLEntry{
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

func testAccStepListPolicy(t *testing.T, names []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "roles/",
		Check: func(resp *logical.Response) error {
			respKeys := resp.Data["keys"].([]string)
			if !reflect.DeepEqual(respKeys, names) {
				return fmt.Errorf("mismatch: %#v %#v", respKeys, names)
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
