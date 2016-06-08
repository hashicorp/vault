package ssh

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
)

// Before the following tests are run, a username going by the name 'vaultssh' has
// to be created and its ~/.ssh/authorized_keys file should contain the below key.
//
// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC9i+hFxZHGo6KblVme4zrAcJstR6I0PTJozW286X4WyvPnkMYDQ5mnhEYC7UWCvjoTWbPEXPX7NjhRtwQTGD67bV+lrxgfyzK1JZbUXK4PwgKJvQD+XyyWYMzDgGSQY61KUSqCxymSm/9NZkPU3ElaQ9xQuTzPpztM4ROfb8f2Yv6/ZESZsTo0MTAkp8Pcy+WkioI/uJ1H7zqs0EA4OMY4aDJRu0UtP4rTVeYNEAuRXdX+eH4aW3KMvhzpFTjMbaJHJXlEeUm2SaX5TNQyTOvghCeQILfYIL/Ca2ij8iwCmulwdV6eQGfd4VDu40PvSnmfoaE38o6HaPnX0kUcnKiT

const (
	testIP               = "127.0.0.1"
	testUserName         = "vaultssh"
	testAdminUser        = "vaultssh"
	testOTPKeyType       = "otp"
	testDynamicKeyType   = "dynamic"
	testCIDRList         = "127.0.0.1/32"
	testDynamicRoleName  = "testDynamicRoleName"
	testOTPRoleName      = "testOTPRoleName"
	testKeyName          = "testKeyName"
	testSharedPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAvYvoRcWRxqOim5VZnuM6wHCbLUeiND0yaM1tvOl+Fsrz55DG
A0OZp4RGAu1Fgr46E1mzxFz1+zY4UbcEExg+u21fpa8YH8sytSWW1FyuD8ICib0A
/l8slmDMw4BkkGOtSlEqgscpkpv/TWZD1NxJWkPcULk8z6c7TOETn2/H9mL+v2RE
mbE6NDEwJKfD3MvlpIqCP7idR+86rNBAODjGOGgyUbtFLT+K01XmDRALkV3V/nh+
GltyjL4c6RU4zG2iRyV5RHlJtkml+UzUMkzr4IQnkCC32CC/wmtoo/IsAprpcHVe
nkBn3eFQ7uND70p5n6GhN/KOh2j519JFHJyokwIDAQABAoIBAHX7VOvBC3kCN9/x
+aPdup84OE7Z7MvpX6w+WlUhXVugnmsAAVDczhKoUc/WktLLx2huCGhsmKvyVuH+
MioUiE+vx75gm3qGx5xbtmOfALVMRLopjCnJYf6EaFA0ZeQ+NwowNW7Lu0PHmAU8
Z3JiX8IwxTz14DU82buDyewO7v+cEr97AnERe3PUcSTDoUXNaoNxjNpEJkKREY6h
4hAY676RT/GsRcQ8tqe/rnCqPHNd7JGqL+207FK4tJw7daoBjQyijWuB7K5chSal
oPInylM6b13ASXuOAOT/2uSUBWmFVCZPDCmnZxy2SdnJGbsJAMl7Ma3MUlaGvVI+
Tfh1aQkCgYEA4JlNOabTb3z42wz6mz+Nz3JRwbawD+PJXOk5JsSnV7DtPtfgkK9y
6FTQdhnozGWShAvJvc+C4QAihs9AlHXoaBY5bEU7R/8UK/pSqwzam+MmxmhVDV7G
IMQPV0FteoXTaJSikhZ88mETTegI2mik+zleBpVxvfdhE5TR+lq8Br0CgYEA2AwJ
CUD5CYUSj09PluR0HHqamWOrJkKPFPwa+5eiTTCzfBBxImYZh7nXnWuoviXC0sg2
AuvCW+uZ48ygv/D8gcz3j1JfbErKZJuV+TotK9rRtNIF5Ub7qysP7UjyI7zCssVM
kuDd9LfRXaB/qGAHNkcDA8NxmHW3gpln4CFdSY8CgYANs4xwfercHEWaJ1qKagAe
rZyrMpffAEhicJ/Z65lB0jtG4CiE6w8ZeUMWUVJQVcnwYD+4YpZbX4S7sJ0B8Ydy
AhkSr86D/92dKTIt2STk6aCN7gNyQ1vW198PtaAWH1/cO2UHgHOy3ZUt5X/Uwxl9
cex4flln+1Viumts2GgsCQKBgCJH7psgSyPekK5auFdKEr5+Gc/jB8I/Z3K9+g4X
5nH3G1PBTCJYLw7hRzw8W/8oALzvddqKzEFHphiGXK94Lqjt/A4q1OdbCrhiE68D
My21P/dAKB1UYRSs9Y8CNyHCjuZM9jSMJ8vv6vG/SOJPsnVDWVAckAbQDvlTHC9t
O98zAoGAcbW6uFDkrv0XMCpB9Su3KaNXOR0wzag+WIFQRXCcoTvxVi9iYfUReQPi
oOyBJU/HMVvBfv4g+OVFLVgSwwm6owwsouZ0+D/LasbuHqYyqYqdyPJQYzWA2Y+F
+B6f4RoPdSXj24JHPg/ioRxjaj094UXJxua2yfkcecGNEuBQHSs=
-----END RSA PRIVATE KEY-----
`
)

func TestBackend_allowed_users(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}

	roleData := map[string]interface{}{
		"key_type":      "otp",
		"default_user":  "ubuntu",
		"cidr_list":     "52.207.235.245/16",
		"allowed_users": "test",
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/role1",
		Storage:   config.StorageView,
		Data:      roleData,
	}

	resp, err := b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to create role: resp:%#v err:%s", resp, err)
	}

	credsData := map[string]interface{}{
		"ip":       "52.207.235.245",
		"username": "ubuntu",
	}
	credsReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Path:      "creds/role1",
		Data:      credsData,
	}

	resp, err = b.HandleRequest(credsReq)
	if err != nil || (resp != nil && resp.IsError()) || resp == nil {
		t.Fatalf("failed to create role: resp:%#v err:%s", resp, err)
	}
	if resp.Data["key"] == "" ||
		resp.Data["key_type"] != "otp" ||
		resp.Data["ip"] != "52.207.235.245" ||
		resp.Data["username"] != "ubuntu" {
		t.Fatalf("failed to create credential: resp:%#v", resp)
	}

	credsData["username"] = "test"
	resp, err = b.HandleRequest(credsReq)
	if err != nil || (resp != nil && resp.IsError()) || resp == nil {
		t.Fatalf("failed to create role: resp:%#v err:%s", resp, err)
	}
	if resp.Data["key"] == "" ||
		resp.Data["key_type"] != "otp" ||
		resp.Data["ip"] != "52.207.235.245" ||
		resp.Data["username"] != "test" {
		t.Fatalf("failed to create credential: resp:%#v", resp)
	}

	credsData["username"] = "random"
	resp, err = b.HandleRequest(credsReq)
	if err != nil || resp == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("expected failure: resp:%#v err:%s", resp, err)
	}

	delete(roleData, "allowed_users")
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to create role: resp:%#v err:%s", resp, err)
	}

	credsData["username"] = "ubuntu"
	resp, err = b.HandleRequest(credsReq)
	if err != nil || (resp != nil && resp.IsError()) || resp == nil {
		t.Fatalf("failed to create role: resp:%#v err:%s", resp, err)
	}
	if resp.Data["key"] == "" ||
		resp.Data["key_type"] != "otp" ||
		resp.Data["ip"] != "52.207.235.245" ||
		resp.Data["username"] != "ubuntu" {
		t.Fatalf("failed to create credential: resp:%#v", resp)
	}

	credsData["username"] = "test"
	resp, err = b.HandleRequest(credsReq)
	if err != nil || resp == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("expected failure: resp:%#v err:%s", resp, err)
	}

	roleData["allowed_users"] = "*"
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to create role: resp:%#v err:%s", resp, err)
	}

	resp, err = b.HandleRequest(credsReq)
	if err != nil || (resp != nil && resp.IsError()) || resp == nil {
		t.Fatalf("failed to create role: resp:%#v err:%s", resp, err)
	}
	if resp.Data["key"] == "" ||
		resp.Data["key_type"] != "otp" ||
		resp.Data["ip"] != "52.207.235.245" ||
		resp.Data["username"] != "test" {
		t.Fatalf("failed to create credential: resp:%#v", resp)
	}
}

func testingFactory(conf *logical.BackendConfig) (logical.Backend, error) {
	_, err := vault.StartSSHHostTestServer()
	if err != nil {
		panic(fmt.Sprintf("error starting mock server:%s", err))
	}
	defaultLeaseTTLVal := 2 * time.Minute
	maxLeaseTTLVal := 10 * time.Minute
	return Factory(&logical.BackendConfig{
		Logger:      nil,
		StorageView: &logical.InmemStorage{},
		System: &logical.StaticSystemView{
			DefaultLeaseTTLVal: defaultLeaseTTLVal,
			MaxLeaseTTLVal:     maxLeaseTTLVal,
		},
	})
}

func TestSSHBackend_Lookup(t *testing.T) {
	testOTPRoleData := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
		"cidr_list":    testCIDRList,
	}
	testDynamicRoleData := map[string]interface{}{
		"key_type":     testDynamicKeyType,
		"key":          testKeyName,
		"admin_user":   testAdminUser,
		"default_user": testAdminUser,
		"cidr_list":    testCIDRList,
	}
	data := map[string]interface{}{
		"ip": testIP,
	}
	resp1 := []string(nil)
	resp2 := []string{testOTPRoleName}
	resp3 := []string{testDynamicRoleName, testOTPRoleName}
	resp4 := []string{testDynamicRoleName}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Factory:        testingFactory,
		Steps: []logicaltest.TestStep{
			testLookupRead(t, data, resp1),
			testRoleWrite(t, testOTPRoleName, testOTPRoleData),
			testLookupRead(t, data, resp2),
			testNamedKeysWrite(t, testKeyName, testSharedPrivateKey),
			testRoleWrite(t, testDynamicRoleName, testDynamicRoleData),
			testLookupRead(t, data, resp3),
			testRoleDelete(t, testOTPRoleName),
			testLookupRead(t, data, resp4),
			testRoleDelete(t, testDynamicRoleName),
			testLookupRead(t, data, resp1),
		},
	})
}

func TestSSHBackend_DynamicKeyCreate(t *testing.T) {
	testDynamicRoleData := map[string]interface{}{
		"key_type":     testDynamicKeyType,
		"key":          testKeyName,
		"admin_user":   testAdminUser,
		"default_user": testAdminUser,
		"cidr_list":    testCIDRList,
	}
	data := map[string]interface{}{
		"username": testUserName,
		"ip":       testIP,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Factory:        testingFactory,
		Steps: []logicaltest.TestStep{
			testNamedKeysWrite(t, testKeyName, testSharedPrivateKey),
			testRoleWrite(t, testDynamicRoleName, testDynamicRoleData),
			testCredsWrite(t, testDynamicRoleName, data, false),
		},
	})
}

func TestSSHBackend_OTPRoleCrud(t *testing.T) {
	testOTPRoleData := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
		"cidr_list":    testCIDRList,
	}
	respOTPRoleData := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"port":         22,
		"default_user": testUserName,
		"cidr_list":    testCIDRList,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Factory:        testingFactory,
		Steps: []logicaltest.TestStep{
			testRoleWrite(t, testOTPRoleName, testOTPRoleData),
			testRoleRead(t, testOTPRoleName, respOTPRoleData),
			testRoleDelete(t, testOTPRoleName),
			testRoleRead(t, testOTPRoleName, nil),
		},
	})
}

func TestSSHBackend_DynamicRoleCrud(t *testing.T) {
	testDynamicRoleData := map[string]interface{}{
		"key_type":     testDynamicKeyType,
		"key":          testKeyName,
		"admin_user":   testAdminUser,
		"default_user": testAdminUser,
		"cidr_list":    testCIDRList,
	}
	respDynamicRoleData := map[string]interface{}{
		"cidr_list":      testCIDRList,
		"port":           22,
		"install_script": DefaultPublicKeyInstallScript,
		"key_bits":       1024,
		"key":            testKeyName,
		"admin_user":     testUserName,
		"default_user":   testUserName,
		"key_type":       testDynamicKeyType,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Factory:        testingFactory,
		Steps: []logicaltest.TestStep{
			testNamedKeysWrite(t, testKeyName, testSharedPrivateKey),
			testRoleWrite(t, testDynamicRoleName, testDynamicRoleData),
			testRoleRead(t, testDynamicRoleName, respDynamicRoleData),
			testRoleDelete(t, testDynamicRoleName),
			testRoleRead(t, testDynamicRoleName, nil),
		},
	})
}

func TestSSHBackend_NamedKeysCrud(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Factory:        testingFactory,
		Steps: []logicaltest.TestStep{
			testNamedKeysWrite(t, testKeyName, testSharedPrivateKey),
			testNamedKeysDelete(t),
		},
	})
}

func TestSSHBackend_OTPCreate(t *testing.T) {
	testOTPRoleData := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
		"cidr_list":    testCIDRList,
	}
	data := map[string]interface{}{
		"username": testUserName,
		"ip":       testIP,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Factory:        testingFactory,
		Steps: []logicaltest.TestStep{
			testRoleWrite(t, testOTPRoleName, testOTPRoleData),
			testCredsWrite(t, testOTPRoleName, data, false),
		},
	})
}

func TestSSHBackend_VerifyEcho(t *testing.T) {
	verifyData := map[string]interface{}{
		"otp": api.VerifyEchoRequest,
	}
	expectedData := map[string]interface{}{
		"message": api.VerifyEchoResponse,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Factory:        testingFactory,
		Steps: []logicaltest.TestStep{
			testVerifyWrite(t, verifyData, expectedData),
		},
	})
}

func TestSSHBackend_ConfigZeroAddressCRUD(t *testing.T) {
	testOTPRoleData := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
		"cidr_list":    testCIDRList,
	}
	testDynamicRoleData := map[string]interface{}{
		"key_type":     testDynamicKeyType,
		"key":          testKeyName,
		"admin_user":   testAdminUser,
		"default_user": testAdminUser,
		"cidr_list":    testCIDRList,
	}
	req1 := map[string]interface{}{
		"roles": testOTPRoleName,
	}
	resp1 := map[string]interface{}{
		"roles": []string{testOTPRoleName},
	}
	req2 := map[string]interface{}{
		"roles": fmt.Sprintf("%s,%s", testOTPRoleName, testDynamicRoleName),
	}
	resp2 := map[string]interface{}{
		"roles": []string{testOTPRoleName, testDynamicRoleName},
	}
	resp3 := map[string]interface{}{
		"roles": []string{},
	}

	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Factory:        testingFactory,
		Steps: []logicaltest.TestStep{
			testRoleWrite(t, testOTPRoleName, testOTPRoleData),
			testConfigZeroAddressWrite(t, req1),
			testConfigZeroAddressRead(t, resp1),
			testNamedKeysWrite(t, testKeyName, testSharedPrivateKey),
			testRoleWrite(t, testDynamicRoleName, testDynamicRoleData),
			testConfigZeroAddressWrite(t, req2),
			testConfigZeroAddressRead(t, resp2),
			testRoleDelete(t, testDynamicRoleName),
			testConfigZeroAddressRead(t, resp1),
			testRoleDelete(t, testOTPRoleName),
			testConfigZeroAddressRead(t, resp3),
			testConfigZeroAddressDelete(t),
		},
	})
}

func TestSSHBackend_CredsForZeroAddressRoles(t *testing.T) {
	dynamicRoleData := map[string]interface{}{
		"key_type":     testDynamicKeyType,
		"key":          testKeyName,
		"admin_user":   testAdminUser,
		"default_user": testAdminUser,
	}
	otpRoleData := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
	}
	data := map[string]interface{}{
		"username": testUserName,
		"ip":       testIP,
	}
	req1 := map[string]interface{}{
		"roles": testOTPRoleName,
	}
	req2 := map[string]interface{}{
		"roles": fmt.Sprintf("%s,%s", testOTPRoleName, testDynamicRoleName),
	}
	logicaltest.Test(t, logicaltest.TestCase{
		AcceptanceTest: true,
		Factory:        testingFactory,
		Steps: []logicaltest.TestStep{
			testRoleWrite(t, testOTPRoleName, otpRoleData),
			testCredsWrite(t, testOTPRoleName, data, true),
			testConfigZeroAddressWrite(t, req1),
			testCredsWrite(t, testOTPRoleName, data, false),
			testNamedKeysWrite(t, testKeyName, testSharedPrivateKey),
			testRoleWrite(t, testDynamicRoleName, dynamicRoleData),
			testCredsWrite(t, testDynamicRoleName, data, true),
			testConfigZeroAddressWrite(t, req2),
			testCredsWrite(t, testDynamicRoleName, data, false),
			testConfigZeroAddressDelete(t),
			testCredsWrite(t, testOTPRoleName, data, true),
			testCredsWrite(t, testDynamicRoleName, data, true),
		},
	})
}

func testConfigZeroAddressDelete(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "config/zeroaddress",
	}
}

func testConfigZeroAddressWrite(t *testing.T, data map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/zeroaddress",
		Data:      data,
	}
}

func testConfigZeroAddressRead(t *testing.T, expected map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "config/zeroaddress",
		Check: func(resp *logical.Response) error {
			var d zeroAddressRoles
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			var ex zeroAddressRoles
			if err := mapstructure.Decode(expected, &ex); err != nil {
				return err
			}

			if !reflect.DeepEqual(d, ex) {
				return fmt.Errorf("Response mismatch:\nActual:%#v\nExpected:%#v", d, ex)
			}

			return nil
		},
	}
}

func testVerifyWrite(t *testing.T, data map[string]interface{}, expected map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("verify"),
		Data:      data,
		Check: func(resp *logical.Response) error {
			var ac api.SSHVerifyResponse
			if err := mapstructure.Decode(resp.Data, &ac); err != nil {
				return err
			}
			var ex api.SSHVerifyResponse
			if err := mapstructure.Decode(expected, &ex); err != nil {
				return err
			}

			if !reflect.DeepEqual(ac, ex) {
				return fmt.Errorf("Invalid response")
			}
			return nil
		},
	}
}

func testNamedKeysWrite(t *testing.T, name, key string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("keys/%s", name),
		Data: map[string]interface{}{
			"key": key,
		},
	}
}

func testNamedKeysDelete(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      fmt.Sprintf("keys/%s", testKeyName),
	}
}

func testLookupRead(t *testing.T, data map[string]interface{}, expected []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "lookup",
		Data:      data,
		Check: func(resp *logical.Response) error {
			if resp.Data == nil || resp.Data["roles"] == nil {
				return fmt.Errorf("Missing roles information")
			}
			if !reflect.DeepEqual(resp.Data["roles"].([]string), expected) {
				return fmt.Errorf("Invalid response: \nactual:%#v\nexpected:%#v", resp.Data["roles"].([]string), expected)
			}
			return nil
		},
	}
}

func testRoleWrite(t *testing.T, name string, data map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data:      data,
	}
}

func testRoleRead(t *testing.T, roleName string, expected map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + roleName,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if expected == nil {
					return nil
				}
				return fmt.Errorf("bad: %#v", resp)
			}
			var d sshRole
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return fmt.Errorf("error decoding response:%s", err)
			}
			if roleName == testOTPRoleName {
				if d.KeyType != expected["key_type"] || d.DefaultUser != expected["default_user"] || d.CIDRList != expected["cidr_list"] {
					return fmt.Errorf("data mismatch. bad: %#v", resp)
				}
			} else {
				if d.AdminUser != expected["admin_user"] || d.CIDRList != expected["cidr_list"] || d.KeyName != expected["key"] || d.KeyType != expected["key_type"] {
					return fmt.Errorf("data mismatch. bad: %#v", resp)
				}
			}
			return nil
		},
	}
}

func testRoleDelete(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + name,
	}
}

func testCredsWrite(t *testing.T, roleName string, data map[string]interface{}, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("creds/%s", roleName),
		Data:      data,
		ErrorOk:   true,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				return fmt.Errorf("response is nil")
			}
			if resp.Data == nil {
				return fmt.Errorf("data is nil")
			}
			if expectError {
				var e struct {
					Error string `mapstructure:"error"`
				}
				if err := mapstructure.Decode(resp.Data, &e); err != nil {
					return err
				}
				if len(e.Error) == 0 {
					return fmt.Errorf("expected error, but write succeeded.")
				}
				return nil
			}
			if roleName == testDynamicRoleName {
				var d struct {
					Key string `mapstructure:"key"`
				}
				if err := mapstructure.Decode(resp.Data, &d); err != nil {
					return err
				}
				if d.Key == "" {
					return fmt.Errorf("Generated key is an empty string")
				}
				// Checking only for a parsable key
				_, err := ssh.ParsePrivateKey([]byte(d.Key))
				if err != nil {
					return fmt.Errorf("Generated key is invalid")
				}
			} else {
				if resp.Data["key_type"] != KeyTypeOTP {
					return fmt.Errorf("Incorrect key_type")
				}
				if resp.Data["key"] == nil {
					return fmt.Errorf("Invalid key")
				}
			}
			return nil
		},
	}
}
