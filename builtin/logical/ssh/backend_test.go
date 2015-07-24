package ssh

import (
	"fmt"
	"log"
	"os/user"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
)

const (
	testCidr             = "127.0.0.1/32"
	testRoleName         = "testRoleName"
	testKey              = "testKey"
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
var testPort string
var testUserName string
var testAdminUser string

// Starts the server and initializes the servers IP address,
// port and usernames to be used by the test cases.
func init() {
	addr, err := vault.StartTestServer()
	if err != nil {
		panic(fmt.Sprintf("Error starting mock server:%s", err))
	}
	input := strings.Split(addr, ":")
	testIP = input[0]
	testPort = input[1]

	u, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("Error getting current username: '%s'", err))
	}
	testUserName = u.Username
	testAdminUser = u.Username
}

func TestSSHDynamicKeyBackend(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testNamedKeys(t),
			testNewDynamicKeyRole(t),
			testDynamicKeyCredsCreate(t),
		},
	})
}

func testNamedKeys(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      fmt.Sprintf("keys/%s", testKey),
		Data: map[string]interface{}{
			"key": testSharedPrivateKey,
		},
	}
}

func testNewDynamicKeyRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      fmt.Sprintf("roles/%s", testRoleName),
		Data: map[string]interface{}{
			"key_type":   "dynamic",
			"key":        testKey,
			"admin_user": testAdminUser,
			"cidr":       testCidr,
			"port":       testPort,
		},
	}
}

func testDynamicKeyCredsCreate(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      fmt.Sprintf("creds/%s", testRoleName),
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
			log.Printf("[WARN] Generated Key:%s\n", d.Key)
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
