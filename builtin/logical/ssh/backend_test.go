package ssh

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"os/user"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

const (
	testCidr      = "127.0.0.1/32"
	testRoleName  = "testRoleName"
	testKey       = "testKey"
	testPublicKey = `
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCaKEIkyRuzYdWPABDoLSPJY3eMCEOXIE0kRI5jqCwJtbkLFydSPvF7swN3r3v/StSBUP+8jmCD8zbXOxmfZHF1XMYGLVJdqfZDT1VCy0HI7PkJbuTIFhdJo3RyOyOlSzj4JV4I3iN7BFbx8RBckEYegKykOps82hZwJYMdykq2iynVJEw+FEg2Y+Zte4DHcy75kR61HE3PM3BK7R5nIPNcuDXTXQZbmFq57LONi8EjAiVWIZitCGdQJg+8aDAceaHdb8xu3GiZUGWQVO8M3OUYbSqWgPIp7R9JI9XZBfby2twJsgJs4PKIH0kjYRW+0Q3iDZH51RTOX3F8yN8Zk7mv
`
	testPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAmihCJMkbs2HVjwAQ6C0jyWN3jAhDlyBNJESOY6gsCbW5Cxcn
Uj7xe7MDd697/0rUgVD/vI5gg/M21zsZn2RxdVzGBi1SXan2Q09VQstByOz5CW7k
yBYXSaN0cjsjpUs4+CVeCN4jewRW8fEQXJBGHoCspDqbPNoWcCWDHcpKtosp1SRM
PhRINmPmbXuAx3Mu+ZEetRxNzzNwSu0eZyDzXLg1010GW5haueyzjYvBIwIlViGY
rQhnUCYPvGgwHHmh3W/MbtxomVBlkFTvDNzlGG0qloDyKe0fSSPV2QX28trcCbIC
bODyiB9JI2EVvtEN4g2R+dUUzl9xfMjfGZO5rwIDAQABAoIBAGHMUpIVx+4YjiyH
hTJWmNKFuOzsvTyeMHJmz9KneTC7yeYgTUDfT8IDQprmiIrghUp5AZU02kQ7wznu
c4XsahJjxflbPVrQnbv8E4IpgtWeiSuT366UXTfJa/GgVS/jNgQvaKXFj8rWaPZa
0d93ZBSr21rhF2UWko+ZLMJ0eMuvJ6yc+BsNjSXq5tGAeT+0vkMBcP+ltZWoEibq
d3YvxAzDmb4CwG4AqcSF1UMnuF6GEdRc/NLlq6YB72pPWaOi2oVEkIQPeMdSfTj/
fFI61JB/MlnkQbAAPq/R/5pGhjiCqHds2uSinAAQuaE/cMdhfFBMYNfvadQIEZzm
U6F7O7ECgYEAzS7o+lm+W/1bAXmOiddwLAF4olXs3q0Am+sbZF6zMsq67ZT3txU2
V3c3vBiXy4MOkOp5CcN9m1hai5CwMxEYoNE77+kwuxFV5pzGnHseHSbu2hWinLOg
j0+NQwKqy7U55amwz+Y41Wwn9obzU6AXQ38I9Kf+YWDiVIDVEBxVRbcCgYEAwFYu
+fEPAioSg3sn0S+z0TbEFp9p0meZWuqct3Lyn83lOpbfVNL6GSYBFwy92jxhQCMu
vGPzkK6ITRe4rapOjMLWosT6wzfgjubeHlhjt3Ccf4zm9OJQ7ghfqR5lKkxoKwZw
eB/iB/Li+ZCn2HpkrLQ6V4HAuJD2Fj+T7LFn68kCgYEAyPNNd4sXNU6vp4UehX96
u46BUDPpNbin5Qxgmm9o/7CvXGnOJf/fZdA7xLstR0LGrEUHX/mW9eKVYyTEfG8c
+LuTAQcYE84JnD8lATJPLuvnd61CwkfmUxTtW5isH7AQ0Q3dPe/S76rqhLZsbxVW
U2OCKOKy7zoM0AgRI6MsHIcCgYAMd4mj+dQXN9LrYtg53vWw4fPj44FgegaetgZi
fbjsUtRA7/aZ8PL1HlmDvPexZaiIF7+3xmLLRgTfumHmH9vnk9mFw27dqImNubk8
Dk6oXUxHmEKALQtB4pkQxT+ZdkpqP4iawLZN/ZhoxM+cYJKV/zio42gyjnLlDknw
Va9+wQKBgQDE7aUItIquTwNtcOsar7aMAYup7wHprEDSb7Y2PclUamKyLfjvJrX3
7ZyXgH4PxDXeezwd+XdE2qdCwlW+3vMnveA9qFz+jyJ3hcxG+hcHMrTLM0A3NBH1
eWhDYXIMZdnt2TojESQHBZhImgPL0nVfynj+I1uMbb84xGHVkACSHw==
-----END RSA PRIVATE KEY-----
`
)

var testIP string
var testPort string
var testUserName string
var testAdminUser string

func init() {
	addr, err := startTestServer()
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

func TestSSHBackend(t *testing.T) {
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testNamedKeys(t),
			testNewRole(t),
			testRoleCreate(t),
		},
	})
}

func startTestServer() (string, error) {
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(testPublicKey))
	if err != nil {
		return "", fmt.Errorf("Error parsing public key")
	}
	serverConfig := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if bytes.Compare(pubKey.Marshal(), key.Marshal()) == 0 {
				return &ssh.Permissions{}, nil
			} else {
				return nil, fmt.Errorf("Key does not match")
			}
		},
	}
	signer, err := ssh.ParsePrivateKey([]byte(testPrivateKey))
	if err != nil {
		panic("Error parsing private key")
	}
	serverConfig.AddHostKey(signer)

	soc, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", fmt.Errorf("Error listening to connection")
	}

	go func() {
		for {
			conn, err := soc.Accept()
			if err != nil {
				panic(fmt.Sprintf("Error accepting incoming connection: %s", err))
			}
			defer conn.Close()
			sshConn, chanReqs, _, err := ssh.NewServerConn(conn, serverConfig)
			if err != nil {
				panic(fmt.Sprintf("Handshaking error: %v", err))
			}

			go func() {
				for chanReq := range chanReqs {
					go func(chanReq ssh.NewChannel) {
						if chanReq.ChannelType() != "session" {
							chanReq.Reject(ssh.UnknownChannelType, "unknown channel type")
							return
						}

						ch, requests, err := chanReq.Accept()
						if err != nil {
							panic(fmt.Sprintf("Error accepting channel: %s", err))
						}

						go func(ch ssh.Channel, in <-chan *ssh.Request) {
							for req := range in {
								executeCommand(ch, req)
							}
						}(ch, requests)
					}(chanReq)
				}
				sshConn.Close()
			}()
		}
	}()
	return soc.Addr().String(), nil
}

func executeCommand(ch ssh.Channel, req *ssh.Request) {
	command := string(req.Payload[4:])
	cmd := exec.Command("/bin/bash", []string{"-c", command}...)
	req.Reply(true, nil)

	cmd.Stdout = ch
	cmd.Stderr = ch
	cmd.Stdin = ch

	err := cmd.Start()
	if err != nil {
		panic(fmt.Sprintf("Error starting the command: '%s'", err))
	}

	go func() {
		_, err := cmd.Process.Wait()
		if err != nil {
			panic(fmt.Sprintf("Error while waiting for command to finish:'%s'", err))
		}
		ch.Close()
	}()
}

func testRoleCreate(t *testing.T) logicaltest.TestStep {
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

func testNewRole(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      fmt.Sprintf("roles/%s", testRoleName),
		Data: map[string]interface{}{
			"key":        testKey,
			"admin_user": testAdminUser,
			"cidr":       testCidr,
			"port":       testPort,
		},
	}
}

func testNamedKeys(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      fmt.Sprintf("keys/%s", testKey),
		Data: map[string]interface{}{
			"key": testPrivateKey,
		},
	}
}
