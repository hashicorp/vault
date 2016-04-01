package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	logicalssh "github.com/hashicorp/vault/builtin/logical/ssh"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
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
func initTest() {
	addr, err := vault.StartSSHHostTestServer()
	if err != nil {
		panic(fmt.Sprintf("Error starting mock server:%s", err))
	}
	input := strings.Split(addr, ":")
	testIP = input[0]
	testPort = input[1]

	testUserName := os.Getenv("VAULT_SSHTEST_USER")
	if len(testUserName) == 0 {
		panic("VAULT_SSHTEST_USER must be set to the desired user")
	}
	testAdminUser = testUserName
}

// This test is broken. Hence temporarily disabling it.
func testSSH(t *testing.T) {
	initTest()
	// Add the SSH backend to the unsealed test core.
	// This should be done before the unsealed core is created.
	err := vault.AddTestLogicalBackend("ssh", logicalssh.Factory)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	mountCmd := &MountCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{"-address", addr, "ssh"}

	// Mount the SSH backend
	if code := mountCmd.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	client, err := mountCmd.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mounts, err := client.Sys().ListMounts()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Check if SSH backend is mounted or not
	mount, ok := mounts["ssh/"]
	if !ok {
		t.Fatal("should have ssh mount")
	}
	if mount.Type != "ssh" {
		t.Fatal("should have ssh type")
	}

	writeCmd := &WriteCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	// Create a 'named' key in vault
	args = []string{
		"-address", addr,
		"ssh/keys/" + testKey,
		"key=" + testSharedPrivateKey,
	}
	if code := writeCmd.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Create a role using the named key along with cidr, username and port
	args = []string{
		"-address", addr,
		"ssh/roles/" + testRoleName,
		"key=" + testKey,
		"admin_user=" + testUserName,
		"cidr=" + testCidr,
		"port=" + testPort,
	}
	if code := writeCmd.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	sshCmd := &SSHCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	// Get the dynamic key and establish an SSH connection with target.
	// Inline command when supplied, runs on target and terminates the
	// connection. Use whoami as the inline command in target and get
	// the result. Compare the result with the username used to connect
	// to target. Test succeeds if they match.
	args = []string{
		"-address", addr,
		"-role=" + testRoleName,
		testUserName + "@" + testIP,
		"/usr/bin/whoami",
	}

	// Creating pipe to get the result of the inline command run in target machine.
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	os.Stdout = w
	if code := sshCmd.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
	bufChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		bufChan <- buf.String()
	}()
	w.Close()
	os.Stdout = stdout
	userName := <-bufChan
	userName = strings.TrimSpace(userName)

	// Comparing the username used to connect to target and
	// the username on the target, thereby verifying successful
	// execution
	if userName != testUserName {
		t.Fatalf("err: username mismatch")
	}
}
