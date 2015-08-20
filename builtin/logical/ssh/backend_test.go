package ssh

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
)

const (
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

var testIP string
var testOTP string
var testPort int
var testUserName string
var testAdminUser string
var testInstallScript string

// Starts the server and initializes the servers IP address,
// port and usernames to be used by the test cases.
func init() {
	addr, err := vault.StartSSHHostTestServer()
	if err != nil {
		panic(fmt.Sprintf("error starting mock server:%s", err))
	}
	input := strings.Split(addr, ":")
	testIP = input[0]
	testPort, err = strconv.Atoi(input[1])
	if err != nil {
		panic(fmt.Sprintf("error parsing port number:%s", err))
	}

	u, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("error getting current username: '%s'", err))
	}
	testUserName = u.Username
	testAdminUser = u.Username
	testInstallScript = DefaultPublicKeyInstallScript
}

func TestSSHBackend_Lookup(t *testing.T) {
	data := map[string]interface{}{
		"ip": testIP,
	}
	otpData := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
		"cidr_list":    testCIDRList,
	}
	dynamicData := map[string]interface{}{
		"key_type":       testDynamicKeyType,
		"key":            testKeyName,
		"admin_user":     testAdminUser,
		"default_user":   testAdminUser,
		"cidr_list":      testCIDRList,
		"install_script": testInstallScript,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testLookupRead(t, data, 0),
			testRoleWrite(t, testOTPRoleName, otpData),
			testLookupRead(t, data, 1),
			testNamedKeysWrite(t),
			testRoleWrite(t, testDynamicRoleName, dynamicData),
			testLookupRead(t, data, 2),
			testRoleDelete(t, testOTPRoleName),
			testLookupRead(t, data, 1),
			testRoleDelete(t, testDynamicRoleName),
			testLookupRead(t, data, 0),
		},
	})
}

func TestSSHBackend_DynamicKeyCreate(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testNamedKeysWrite(t),
			testNewDynamicKeyRole(t),
			testDynamicKeyCredsCreate(t),
		},
	})
}

func TestSSHBackend_OTPRoleCrud(t *testing.T) {
	data := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
		"cidr_list":    testCIDRList,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testRoleWrite(t, testOTPRoleName, data),
			testRoleRead(t, testOTPRoleName, data),
			testRoleDelete(t, testOTPRoleName),
			testRoleRead(t, testOTPRoleName, nil),
		},
	})
}

func TestSSHBackend_DynamicRoleCrud(t *testing.T) {
	data := map[string]interface{}{
		"key_type":       testDynamicKeyType,
		"key":            testKeyName,
		"admin_user":     testAdminUser,
		"default_user":   testAdminUser,
		"cidr_list":      testCIDRList,
		"install_script": testInstallScript,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testNamedKeysWrite(t),
			testRoleWrite(t, testDynamicRoleName, data),
			testRoleRead(t, testDynamicRoleName, data),
			testRoleDelete(t, testDynamicRoleName),
			testRoleRead(t, testDynamicRoleName, nil),
		},
	})
}

func TestSSHBackend_NamedKeysCrud(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testNamedKeysRead(t, ""),
			testNamedKeysWrite(t),
			testNamedKeysRead(t, testSharedPrivateKey),
			testNamedKeysDelete(t),
		},
	})
}

func TestSSHBackend_OTPCreate(t *testing.T) {
	data := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
		"cidr_list":    testCIDRList,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testRoleWrite(t, testOTPRoleName, data),
			testCredsWrite(t, testOTPRoleName),
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
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testVerifyWrite(t, verifyData, expectedData),
		},
	})
}

func testVerifyWrite(t *testing.T, d map[string]interface{}, expected map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      fmt.Sprintf("verify"),
		Data:      d,
		Check: func(resp *logical.Response) error {
			var ac api.SSHVerifyResponse
			if err := mapstructure.Decode(resp.Data, &ac); err != nil {
				return err
			}
			var ex api.SSHVerifyResponse
			if err := mapstructure.Decode(expected, &ex); err != nil {
				return err
			}

			if ac.Message != ex.Message || ac.IP != ex.IP || ac.Username != ex.Username {
				return fmt.Errorf("Invalid response")
			}
			return nil
		},
	}
}

func testCredsWrite(t *testing.T, name string) logicaltest.TestStep {
	data := map[string]interface{}{
		"ip": testIP,
	}
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      fmt.Sprintf("creds/%s", name),
		Data:      data,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				return fmt.Errorf("response is nil")
			}
			if resp.Data == nil {
				return fmt.Errorf("data is nil")
			}
			if resp.Data["key_type"] != KeyTypeOTP {
				return fmt.Errorf("Incorrect key_type")
			}
			if resp.Data["key"] == nil {
				return fmt.Errorf("Invalid key")
			}
			testOTP = resp.Data["key"].(string)
			return nil
		},
	}
}

func testNamedKeysRead(t *testing.T, key string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      fmt.Sprintf("keys/%s", testKeyName),
		Check: func(resp *logical.Response) error {
			if key != "" {
				if resp == nil || resp.Data == nil {
					return fmt.Errorf("Key missing in response")
				}
				var d struct {
					Key string `mapstructure:"key"`
				}
				if err := mapstructure.Decode(resp.Data, &d); err != nil {
					return err
				}

				if d.Key != key {
					return fmt.Errorf("Key mismatch")
				}
			}
			return nil
		},
	}
}

func testNamedKeysWrite(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      fmt.Sprintf("keys/%s", testKeyName),
		Data: map[string]interface{}{
			"key": testSharedPrivateKey,
		},
	}
}

func testNamedKeysDelete(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      fmt.Sprintf("keys/%s", testKeyName),
	}
}

func testLookupRead(t *testing.T, data map[string]interface{}, length int) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "lookup",
		Data:      data,
		Check: func(resp *logical.Response) error {
			if resp.Data == nil || resp.Data["roles"] == nil {
				return fmt.Errorf("Missing roles information")
			}
			if len(resp.Data["roles"].([]string)) != length {
				return fmt.Errorf("Role information incorrect")
			}
			return nil
		},
	}
}

func testRoleWrite(t *testing.T, name string, data map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "roles/" + name,
		Data:      data,
	}
}

func testRoleRead(t *testing.T, name string, data map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if data == nil {
					return nil
				}
				return fmt.Errorf("bad: %#v", resp)
			}
			var d sshRole
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return fmt.Errorf("error decoding response:%s", err)
			}
			if name == testOTPRoleName {
				if d.KeyType != data["key_type"] || d.DefaultUser != data["default_user"] || d.CIDRList != data["cidr_list"] {
					return fmt.Errorf("data mismatch. bad: %#v", resp)
				}
			} else {
				if d.AdminUser != data["admin_user"] || d.CIDRList != data["cidr_list"] || d.KeyName != data["key"] || d.KeyType != data["key_type"] {
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

func testNewDynamicKeyRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      fmt.Sprintf("roles/%s", testDynamicRoleName),
		Data: map[string]interface{}{
			"key_type":       "dynamic",
			"key":            testKeyName,
			"admin_user":     testAdminUser,
			"default_user":   testAdminUser,
			"cidr_list":      testCIDRList,
			"port":           testPort,
			"install_script": testInstallScript,
		},
	}
}

func testDynamicKeyCredsCreate(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      fmt.Sprintf("creds/%s", testDynamicRoleName),
		Data: map[string]interface{}{
			"username": testUserName,
			"ip":       testIP,
		},
		Check: func(resp *logical.Response) error {
			var d struct {
				Key string `mapstructure:"key"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if d.Key == "" {
				return fmt.Errorf("Generated key is an empty string")
			}
			_, err := ssh.ParsePrivateKey([]byte(d.Key))
			if err != nil {
				return fmt.Errorf("Generated key is invalid")
			}
			return nil
		},
	}
}
