package consul

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

func TestBackend_Config_Access(t *testing.T) {
	t.Run("config_access", func(t *testing.T) {
		t.Parallel()
		t.Run("pre-1.4.0", func(t *testing.T) {
			t.Parallel()
			testBackendConfigAccess(t, "1.3.1", true)
		})
		t.Run("post-1.4.0", func(t *testing.T) {
			t.Parallel()
			testBackendConfigAccess(t, "", true)
		})
		t.Run("pre-1.4.0 automatic-bootstrap", func(t *testing.T) {
			t.Parallel()
			testBackendConfigAccess(t, "1.3.1", false)
		})
		t.Run("post-1.4.0 automatic-bootstrap", func(t *testing.T) {
			t.Parallel()
			testBackendConfigAccess(t, "", false)
		})
	})
}

func testBackendConfigAccess(t *testing.T, version string, bootstrap bool) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, consulConfig := consul.PrepareTestContainer(t, version, false, bootstrap)
	defer cleanup()

	connData := map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	}

	confReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Storage:   config.StorageView,
		Data:      connData,
	}

	resp, err := b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	confReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), confReq)
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

func TestBackend_Renew_Revoke(t *testing.T) {
	t.Run("renew_revoke", func(t *testing.T) {
		t.Parallel()
		t.Run("pre-1.4.0", func(t *testing.T) {
			t.Parallel()
			testBackendRenewRevoke(t, "1.3.1")
		})
		t.Run("post-1.4.0", func(t *testing.T) {
			t.Parallel()
			t.Run("legacy", func(t *testing.T) {
				t.Parallel()
				testBackendRenewRevoke(t, "1.4.4")
			})

			t.Run("param-policies", func(t *testing.T) {
				t.Parallel()
				testBackendRenewRevoke14(t, "", "policies")
			})
			t.Run("param-consul_policies", func(t *testing.T) {
				t.Parallel()
				testBackendRenewRevoke14(t, "", "consul_policies")
			})
			t.Run("both-params", func(t *testing.T) {
				t.Parallel()
				testBackendRenewRevoke14(t, "", "both")
			})
		})
	})
}

func testBackendRenewRevoke(t *testing.T, version string) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, consulConfig := consul.PrepareTestContainer(t, version, false, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "roles/test"
	req.Data = map[string]interface{}{
		"policy": base64.StdEncoding.EncodeToString([]byte(testPolicy)),
		"lease":  "6h",
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test"
	resp, err := b.HandleRequest(context.Background(), req)
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
	generatedSecret.TTL = 6 * time.Hour

	var d struct {
		Token string `mapstructure:"token"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.KV().Put(&consulapi.KVPair{
		Key:   "foo",
		Value: []byte("bar"),
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.RenewOperation
	req.Secret = generatedSecret
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}

	req.Operation = logical.RevokeOperation
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.KV().Put(&consulapi.KVPair{
		Key:   "foo",
		Value: []byte("bar"),
	}, nil)
	if err == nil {
		t.Fatal("err: expected error")
	}
}

func testBackendRenewRevoke14(t *testing.T, version string, policiesParam string) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, consulConfig := consul.PrepareTestContainer(t, version, false, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "roles/test"
	req.Data = map[string]interface{}{
		"lease": "6h",
	}
	if policiesParam == "both" {
		req.Data["policies"] = []string{"wrong-name"}
		req.Data["consul_policies"] = []string{"test"}
	} else {
		req.Data[policiesParam] = []string{"test"}
	}

	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	read := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.ReadOperation,
		Path:      "roles/test",
		Data:      connData,
	}
	roleResp, err := b.HandleRequest(context.Background(), read)

	expectExtract := roleResp.Data["consul_policies"]
	respExtract := roleResp.Data[policiesParam]
	if respExtract != nil {
		if expectExtract.([]string)[0] != respExtract.([]string)[0] {
			t.Errorf("mismatch: response consul_policies '%s' does not match '[test]'", roleResp.Data["consul_policies"])
		}
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test"
	resp, err := b.HandleRequest(context.Background(), req)
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
	generatedSecret.TTL = 6 * time.Hour

	var d struct {
		Token    string `mapstructure:"token"`
		Accessor string `mapstructure:"accessor"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultNonPooledConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Catalog(), nil
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.RenewOperation
	req.Secret = generatedSecret
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}

	req.Operation = logical.RevokeOperation
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Build a management client and verify that the token does not exist anymore
	consulmgmtConfig := consulapi.DefaultNonPooledConfig()
	consulmgmtConfig.Address = connData["address"].(string)
	consulmgmtConfig.Token = connData["token"].(string)
	mgmtclient, err := consulapi.NewClient(consulmgmtConfig)
	if err != nil {
		t.Fatal(err)
	}
	q := &consulapi.QueryOptions{
		Datacenter: "DC1",
	}

	_, _, err = mgmtclient.ACL().TokenRead(d.Accessor, q)
	if err == nil {
		t.Fatal("err: expected error")
	}
}

func TestBackend_LocalToken(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, consulConfig := consul.PrepareTestContainer(t, "", false, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "roles/test"
	req.Data = map[string]interface{}{
		"consul_policies": []string{"test"},
		"ttl":             "6h",
		"local":           false,
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "roles/test_local"
	req.Data = map[string]interface{}{
		"consul_policies": []string{"test"},
		"ttl":             "6h",
		"local":           true,
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test"
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
	if resp.IsError() {
		t.Fatalf("resp is error: %v", resp.Error())
	}

	var d struct {
		Token    string `mapstructure:"token"`
		Accessor string `mapstructure:"accessor"`
		Local    bool   `mapstructure:"local"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}

	if d.Local {
		t.Fatalf("requested global token, got local one")
	}

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultNonPooledConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Catalog(), nil
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test_local"
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
	if resp.IsError() {
		t.Fatalf("resp is error: %v", resp.Error())
	}

	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}

	if !d.Local {
		t.Fatalf("requested local token, got global one")
	}

	// Build a client and verify that the credentials work
	consulapiConfig = consulapi.DefaultNonPooledConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err = consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Catalog(), nil
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_Management(t *testing.T) {
	t.Run("management", func(t *testing.T) {
		t.Parallel()
		t.Run("pre-1.4.0", func(t *testing.T) {
			t.Parallel()
			testBackendManagement(t, "1.3.1")
		})
		t.Run("post-1.4.0", func(t *testing.T) {
			t.Parallel()
			testBackendManagement(t, "1.4.4")
		})

		testBackendManagement(t, "1.10.8")
	})
}

func testBackendManagement(t *testing.T, version string) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, consulConfig := consul.PrepareTestContainer(t, version, false, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData),
			testAccStepWriteManagementPolicy(t, "test", ""),
			testAccStepReadManagementToken(t, "test", connData),
		},
	})
}

func TestBackend_Basic(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		t.Run("pre-1.4.0", func(t *testing.T) {
			t.Parallel()
			testBackendBasic(t, "1.3.1")
		})
		t.Run("post-1.4.0", func(t *testing.T) {
			t.Parallel()
			t.Run("legacy", func(t *testing.T) {
				t.Parallel()
				testBackendBasic(t, "1.4.4")
			})

			testBackendBasic(t, "1.10.8")
		})
	})
}

func testBackendBasic(t *testing.T, version string) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, consulConfig := consul.PrepareTestContainer(t, version, false, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData),
			testAccStepWritePolicy(t, "test", testPolicy, ""),
			testAccStepReadToken(t, "test", connData),
		},
	})
}

func TestBackend_crud(t *testing.T) {
	b, _ := Factory(context.Background(), logical.TestBackendConfig())
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
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
	b, _ := Factory(context.Background(), logical.TestBackendConfig())
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", testPolicy, "6h"),
			testAccStepReadPolicy(t, "test", testPolicy, 6*time.Hour),
			testAccStepDeletePolicy(t, "test"),
		},
	})
}

func testAccStepConfig(t *testing.T, config map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      config,
	}
}

func testAccStepReadToken(t *testing.T, name string, conf map[string]interface{}) logicaltest.TestStep {
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

func testAccStepReadManagementToken(t *testing.T, name string, conf map[string]interface{}) logicaltest.TestStep {
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

			l := resp.Data["lease"].(int64)
			if lease != time.Second*time.Duration(l) {
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

func TestBackend_Roles(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, consulConfig := consul.PrepareTestContainer(t, "", false, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Create the consul_roles role
	req.Path = "roles/test-consul-roles"
	req.Data = map[string]interface{}{
		"consul_roles": []string{"role-test"},
		"lease":        "6h",
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test-consul-roles"
	resp, err := b.HandleRequest(context.Background(), req)
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
	generatedSecret.TTL = 6 * time.Hour

	var d struct {
		Token    string `mapstructure:"token"`
		Accessor string `mapstructure:"accessor"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultNonPooledConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Catalog(), nil
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.RenewOperation
	req.Secret = generatedSecret
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}

	req.Operation = logical.RevokeOperation
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Build a management client and verify that the token does not exist anymore
	consulmgmtConfig := consulapi.DefaultNonPooledConfig()
	consulmgmtConfig.Address = connData["address"].(string)
	consulmgmtConfig.Token = connData["token"].(string)
	mgmtclient, err := consulapi.NewClient(consulmgmtConfig)
	if err != nil {
		t.Fatal(err)
	}
	q := &consulapi.QueryOptions{
		Datacenter: "DC1",
	}

	_, _, err = mgmtclient.ACL().TokenRead(d.Accessor, q)
	if err == nil {
		t.Fatal("err: expected error")
	}
}

func TestBackend_Enterprise_Namespace(t *testing.T) {
	if _, hasLicense := os.LookupEnv("CONSUL_LICENSE"); !hasLicense {
		t.Skip("Skipping: No enterprise license found")
	}

	testBackendEntNamespace(t)
}

func TestBackend_Enterprise_Partition(t *testing.T) {
	if _, hasLicense := os.LookupEnv("CONSUL_LICENSE"); !hasLicense {
		t.Skip("Skipping: No enterprise license found")
	}

	testBackendEntPartition(t)
}

func testBackendEntNamespace(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, consulConfig := consul.PrepareTestContainer(t, "", true, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Create the role in namespace "ns1"
	req.Path = "roles/test-ns"
	req.Data = map[string]interface{}{
		"consul_policies":  []string{"ns-test"},
		"lease":            "6h",
		"consul_namespace": "ns1",
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test-ns"
	resp, err := b.HandleRequest(context.Background(), req)
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
	generatedSecret.TTL = 6 * time.Hour

	var d struct {
		Token           string `mapstructure:"token"`
		Accessor        string `mapstructure:"accessor"`
		ConsulNamespace string `mapstructure:"consul_namespace"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}

	if d.ConsulNamespace != "ns1" {
		t.Fatalf("Failed to access namespace")
	}

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultNonPooledConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Catalog(), nil
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.RenewOperation
	req.Secret = generatedSecret
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}

	req.Operation = logical.RevokeOperation
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Build a management client and verify that the token does not exist anymore
	consulmgmtConfig := consulapi.DefaultNonPooledConfig()
	consulmgmtConfig.Address = connData["address"].(string)
	consulmgmtConfig.Token = connData["token"].(string)
	mgmtclient, err := consulapi.NewClient(consulmgmtConfig)
	if err != nil {
		t.Fatal(err)
	}
	q := &consulapi.QueryOptions{
		Datacenter: "DC1",
		Namespace:  "ns1",
	}

	_, _, err = mgmtclient.ACL().TokenRead(d.Accessor, q)
	if err == nil {
		t.Fatal("err: expected error")
	}
}

func testBackendEntPartition(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, consulConfig := consul.PrepareTestContainer(t, "", true, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Create the role in partition "part1"
	req.Path = "roles/test-part"
	req.Data = map[string]interface{}{
		"consul_policies": []string{"part-test"},
		"lease":           "6h",
		"partition":       "part1",
	}
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test-part"
	resp, err := b.HandleRequest(context.Background(), req)
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
	generatedSecret.TTL = 6 * time.Hour

	var d struct {
		Token     string `mapstructure:"token"`
		Accessor  string `mapstructure:"accessor"`
		Partition string `mapstructure:"partition"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}

	if d.Partition != "part1" {
		t.Fatalf("Failed to access partition")
	}

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultNonPooledConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Catalog(), nil
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.RenewOperation
	req.Secret = generatedSecret
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}

	req.Operation = logical.RevokeOperation
	_, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Build a management client and verify that the token does not exist anymore
	consulmgmtConfig := consulapi.DefaultNonPooledConfig()
	consulmgmtConfig.Address = connData["address"].(string)
	consulmgmtConfig.Token = connData["token"].(string)
	mgmtclient, err := consulapi.NewClient(consulmgmtConfig)
	if err != nil {
		t.Fatal(err)
	}
	q := &consulapi.QueryOptions{
		Datacenter: "DC1",
		Partition:  "test1",
	}

	_, _, err = mgmtclient.ACL().TokenRead(d.Accessor, q)
	if err == nil {
		t.Fatal("err: expected error")
	}
}

func TestBackendRenewRevokeRolesAndIdentities(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, consulConfig := consul.PrepareTestContainer(t, "", false, true)
	defer cleanup()

	connData := map[string]interface{}{
		"address": consulConfig.Address(),
		"token":   consulConfig.Token,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	cases := map[string]struct {
		RoleName string
		RoleData map[string]interface{}
	}{
		"just role": {
			"r",
			map[string]interface{}{
				"consul_roles": []string{"role-test"},
				"lease":        "6h",
			},
		},
		"role and policies": {
			"rp",
			map[string]interface{}{
				"consul_policies": []string{"test"},
				"consul_roles":    []string{"role-test"},
				"lease":           "6h",
			},
		},
		"service identity": {
			"si",
			map[string]interface{}{
				"service_identities": "service1",
				"lease":              "6h",
			},
		},
		"service identity and policies": {
			"sip",
			map[string]interface{}{
				"consul_policies":    []string{"test"},
				"service_identities": "service1",
				"lease":              "6h",
			},
		},
		"service identity and role": {
			"sir",
			map[string]interface{}{
				"consul_roles":       []string{"role-test"},
				"service_identities": "service1",
				"lease":              "6h",
			},
		},
		"service identity and role and policies": {
			"sirp",
			map[string]interface{}{
				"consul_policies":    []string{"test"},
				"consul_roles":       []string{"role-test"},
				"service_identities": "service1",
				"lease":              "6h",
			},
		},
		"node identity": {
			"ni",
			map[string]interface{}{
				"node_identities": []string{"node1:dc1"},
				"lease":           "6h",
			},
		},
		"node identity and policies": {
			"nip",
			map[string]interface{}{
				"consul_policies": []string{"test"},
				"node_identities": []string{"node1:dc1"},
				"lease":           "6h",
			},
		},
		"node identity and role": {
			"nir",
			map[string]interface{}{
				"consul_roles":    []string{"role-test"},
				"node_identities": []string{"node1:dc1"},
				"lease":           "6h",
			},
		},
		"node identity and role and policies": {
			"nirp",
			map[string]interface{}{
				"consul_policies": []string{"test"},
				"consul_roles":    []string{"role-test"},
				"node_identities": []string{"node1:dc1"},
				"lease":           "6h",
			},
		},
		"node identity and service identity": {
			"nisi",
			map[string]interface{}{
				"service_identities": "service1",
				"node_identities":    []string{"node1:dc1"},
				"lease":              "6h",
			},
		},
		"node identity and service identity and policies": {
			"nisip",
			map[string]interface{}{
				"consul_policies":    []string{"test"},
				"service_identities": "service1",
				"node_identities":    []string{"node1:dc1"},
				"lease":              "6h",
			},
		},
		"node identity and service identity and role": {
			"nisir",
			map[string]interface{}{
				"consul_roles":       []string{"role-test"},
				"service_identities": "service1",
				"node_identities":    []string{"node1:dc1"},
				"lease":              "6h",
			},
		},
		"node identity and service identity and role and policies": {
			"nisirp",
			map[string]interface{}{
				"consul_policies":    []string{"test"},
				"consul_roles":       []string{"role-test"},
				"service_identities": "service1",
				"node_identities":    []string{"node1:dc1"},
				"lease":              "6h",
			},
		},
	}

	for description, tc := range cases {
		t.Logf("Testing: %s", description)

		req.Operation = logical.UpdateOperation
		req.Path = fmt.Sprintf("roles/%s", tc.RoleName)
		req.Data = tc.RoleData
		resp, err = b.HandleRequest(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}

		req.Operation = logical.ReadOperation
		req.Path = fmt.Sprintf("creds/%s", tc.RoleName)
		resp, err = b.HandleRequest(context.Background(), req)
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
		generatedSecret.TTL = 6 * time.Hour

		var d struct {
			Token    string `mapstructure:"token"`
			Accessor string `mapstructure:"accessor"`
		}
		if err := mapstructure.Decode(resp.Data, &d); err != nil {
			t.Fatal(err)
		}

		// Build a client and verify that the credentials work
		consulapiConfig := consulapi.DefaultNonPooledConfig()
		consulapiConfig.Address = connData["address"].(string)
		consulapiConfig.Token = d.Token
		client, err := consulapi.NewClient(consulapiConfig)
		if err != nil {
			t.Fatal(err)
		}

		_, err = client.Catalog(), nil
		if err != nil {
			t.Fatal(err)
		}

		req.Operation = logical.RenewOperation
		req.Secret = generatedSecret
		resp, err = b.HandleRequest(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("got nil response from renew")
		}

		req.Operation = logical.RevokeOperation
		resp, err = b.HandleRequest(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}

		// Build a management client and verify that the token does not exist anymore
		consulmgmtConfig := consulapi.DefaultNonPooledConfig()
		consulmgmtConfig.Address = connData["address"].(string)
		consulmgmtConfig.Token = connData["token"].(string)
		mgmtclient, err := consulapi.NewClient(consulmgmtConfig)

		q := &consulapi.QueryOptions{
			Datacenter: "DC1",
		}

		_, _, err = mgmtclient.ACL().TokenRead(d.Accessor, q)
		if err == nil {
			t.Fatal("err: expected error")
		}
	}
}

const testPolicy = `
key "" {
	policy = "write"
}`
