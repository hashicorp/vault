// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/docker"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

const (
	testIP            = "127.0.0.1"
	testUserName      = "vaultssh"
	testMultiUserName = "vaultssh,otherssh"
	testAdminUser     = "vaultssh"
	testCaKeyType     = "ca"
	testOTPKeyType    = "otp"
	testCIDRList      = "127.0.0.1/32"
	testAtRoleName    = "test@RoleName"
	testOTPRoleName   = "testOTPRoleName"
	// testKeyName is the name of the entry that will be written to SSHMOUNTPOINT/ssh/keys
	testKeyName = "testKeyName"
	// testSharedPrivateKey is the value of the entry that will be written to SSHMOUNTPOINT/ssh/keys
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
	// Public half of `testCAPrivateKey`, identical to how it would be fed in from a file
	testCAPublicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDArgK0ilRRfk8E7HIsjz5l3BuxmwpDd8DHRCVfOhbZ4gOSVxjEOOqBwWGjygdboBIZwFXmwDlU6sWX0hBJAgpQz0Cjvbjxtq/NjkvATrYPgnrXUhTaEn2eQO0PsqRNSFH46SK/oJfTp0q8/WgojxWJ2L7FUV8PO8uIk49DzqAqPV7WXU63vFsjx+3WQOX/ILeQvHCvaqs3dWjjzEoDudRWCOdUqcHEOshV9azIzPrXlQVzRV3QAKl6u7pC+/Secorpwt6IHpMKoVPGiR0tMMuNOVH8zrAKzIxPGfy2WmNDpJopbXMTvSOGAqNcp49O4SKOQl9Fzfq2HEevJamKLrMB dummy@example.com
`
	publicKey2 = `AAAAB3NzaC1yc2EAAAADAQABAAABAQDArgK0ilRRfk8E7HIsjz5l3BuxmwpDd8DHRCVfOhbZ4gOSVxjEOOqBwWGjygdboBIZwFXmwDlU6sWX0hBJAgpQz0Cjvbjxtq/NjkvATrYPgnrXUhTaEn2eQO0PsqRNSFH46SK/oJfTp0q8/WgojxWJ2L7FUV8PO8uIk49DzqAqPV7WXU63vFsjx+3WQOX/ILeQvHCvaqs3dWjjzEoDudRWCOdUqcHEOshV9azIzPrXlQVzRV3QAKl6u7pC+/Secorpwt6IHpMKoVPGiR0tMMuNOVH8zrAKzIxPGfy2WmNDpJopbXMTvSOGAqNcp49O4SKOQl9Fzfq2HEevJamKLrMB
`

	publicKey3072 = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDlsMr3K1d0nzE1TjUULPRuVjEGETmOqHtWq4gVPq3HiuNVHE/e/BJnkXc40BoClQ2Z5ZZPJZ6izF9PnlzNDjpq8DrILUrn/6KrzCHvRwnkYMAXbfM/Br09z5QGptbOe1EMLeVe0b/udmUicbYAGPxMruZk+ljyr4vXkO+gOAIrxeSIQSdMVLU4g0pCPQuDCOx5IQpDYSlOB3091frpN8npfMueKPflNYzxnqqYgAVeDKAIqMCGOMOHUeIZJ7A7HuynEAVOsOkJwC9nesy9D6ppdWNduGl42IkzlwVdDMZtUAEznMUT/dnHNG1Krx9SuNZ/S9fGjxGVsT+jzUmizrWB9/6XIEHDxPBzcqlWFuwYTGz1OL8bfZ+HldOGPcnqZn9hKntWwjUc3whcvWt+NCmXpHSVLSxf+WN8pdmfEsCqn8mpvo2MXa+iJrtAVPX4i0u8AQUuqC3NuXHv4Cn0LNwtziBT544UjgbWkAZqzFZJREYA09OHscc3akEIrTnPehk= demo@example.com`

	publicKey4096 = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC54Oj4YCFDYxYv69Q9KfU6rWYtUB1eByQdUW0nXFi/vr98QUIV77sEeUVhaQzZcuCojAi/GrloW7ta0Z2DaEv5jOQMAnGpXBcqLJsz3KdrHbpvl93MPNdmNaGPU0GnUEsjBVuDVn9HdIUa8CNrxShvPu7/VqoaRHKLqphGgzFb37vi4qvnQ+5VYAO/TzyVYMD6qJX6I/9Pw8d74jCfEdOh2yGKkP7rXWOghreyIl8H2zTJKg9KoZuPq9F5M8nNt7Oi3rf+DwQiYvamzIqlDP4s5oFVTZW0E9lwWvYDpyiJnUrkQqksebBK/rcyfiFG3onb4qLo2WVWXeK3si8IhGik/TEzprScyAWIf9RviT8O+l5hTA2/c+ctn3MVCLRNfez2lKpdxCoprv1MbIcySGWblTJEcY6RA+aauVJpu7FMtRxHHtZKtMpep8cLu8GKbiP6Ifq2JXBtXtNxDeIgo2MkNoMh/NHAsACJniE/dqV/+u9HvhvgrTbJ69ell0nE4ivzA7O4kZgbR/4MHlLgLFvaqC8RrWRLY6BdFagPIMxghWha7Qw16zqoIjRnolvRzUWvSXanJVg8Z6ua1VxwgirNaAH1ivmJhUh2+4lNxCX6jmZyR3zjJsWY03gjJTairvI762opjjalF8fH6Xrs15mB14JiAlNbk6+5REQcvXlGqw== dummy@example.com`

	testCAPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwK4CtIpUUX5PBOxyLI8+ZdwbsZsKQ3fAx0QlXzoW2eIDklcY
xDjqgcFho8oHW6ASGcBV5sA5VOrFl9IQSQIKUM9Ao7248bavzY5LwE62D4J611IU
2hJ9nkDtD7KkTUhR+Okiv6CX06dKvP1oKI8Vidi+xVFfDzvLiJOPQ86gKj1e1l1O
t7xbI8ft1kDl/yC3kLxwr2qrN3Vo48xKA7nUVgjnVKnBxDrIVfWsyMz615UFc0Vd
0ACperu6Qvv0nnKK6cLeiB6TCqFTxokdLTDLjTlR/M6wCsyMTxn8tlpjQ6SaKW1z
E70jhgKjXKePTuEijkJfRc36thxHryWpii6zAQIDAQABAoIBAA/DrPD8iF2KigiL
F+RRa/eFhLaJStOuTpV/G9eotwnolgY5Hguf5H/tRIHUG7oBZLm6pMyWWZp7AuOj
CjYO9q0Z5939vc349nVI+SWoyviF4msPiik1bhWulja8lPjFu/8zg+ZNy15Dx7ei
vAzleAupMiKOv8pNSB/KguQ3WZ9a9bcQcoFQ2Foru6mXpLJ03kghVRlkqvQ7t5cA
n11d2Hiipq9mleESr0c+MUPKLBX/neaWfGA4xgJTjIYjZi6avmYc/Ox3sQ9aLq2J
tH0D4HVUZvaU28hn+jhbs64rRFbu++qQMe3vNvi/Q/iqcYU4b6tgDNzm/JFRTS/W
njiz4mkCgYEA44CnQVmonN6qQ0AgNNlBY5+RX3wwBJZ1AaxpzwDRylAt2vlVUA0n
YY4RW4J4+RMRKwHwjxK5RRmHjsIJx+nrpqihW3fte3ev5F2A9Wha4dzzEHxBY6IL
362T/x2f+vYk6tV+uTZSUPHsuELH26mitbBVFNB/00nbMNdEc2bO5FMCgYEA2NCw
ubt+g2bRkkT/Qf8gIM8ZDpZbARt6onqxVcWkQFT16ZjbsBWUrH1Xi7alv9+lwYLJ
ckY/XDX4KeU19HabeAbpyy6G9Q2uBSWZlJbjl7QNhdLeuzV82U1/r8fy6Uu3gQnU
WSFx2GesRpSmZpqNKMs5ksqteZ9Yjg1EIgXdINsCgYBIn9REt1NtKGOf7kOZu1T1
cYXdvm4xuLoHW7u3OiK+e9P3mCqU0G4m5UxDMyZdFKohWZAqjCaamWi9uNGYgOMa
I7DG20TzaiS7OOIm9TY17eul8pSJMrypnealxRZB7fug/6Bhjaa/cktIEwFr7P4l
E/JFH73+fBA9yipu0H3xQwKBgHmiwrLAZF6VrVcxDD9bQQwHA5iyc4Wwg+Fpkdl7
0wUgZQHTdtRXlxwaCaZhJqX5c4WXuSo6DMvPn1TpuZZXgCsbPch2ZtJOBWXvzTSW
XkK6iaedQMWoYU2L8+mK9FU73EwxVodWgwcUSosiVCRV6oGLWdZnjGEiK00uVh38
Si1nAoGBAL47wWinv1cDTnh5mm0mybz3oI2a6V9aIYCloQ/EFcvtahyR/gyB8qNF
lObH9Faf0WGdnACZvTz22U9gWhw79S0SpDV31tC5Kl8dXHFiZ09vYUKkYmSd/kms
SeKWrUkryx46LVf6NMhkyYmRqCEjBwfOozzezi5WbiJy6nn54GQt
-----END RSA PRIVATE KEY-----
`

	testCAPublicKeyEd25519 = `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIO1S6g5Bib7vT8eoFnvTl3dZSjOQL/GkH1nkRcDS9++a ca
`

	testCAPrivateKeyEd25519 = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDtUuoOQYm+70/HqBZ705d3WUozkC/xpB9Z5EXA0vfvmgAAAIhfRuszX0br
MwAAAAtzc2gtZWQyNTUxOQAAACDtUuoOQYm+70/HqBZ705d3WUozkC/xpB9Z5EXA0vfvmg
AAAEBQYa029SP/7AGPFQLmzwOc9eCoOZuwCq3iIf2C6fj9j+1S6g5Bib7vT8eoFnvTl3dZ
SjOQL/GkH1nkRcDS9++aAAAAAmNhAQID
-----END OPENSSH PRIVATE KEY-----
`

	publicKeyECDSA256 = `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBJsfOouYIjJNI23QJqaDsFTGukm21fRAMeGvKZDB59i5jnX1EubMH1AEjjzz4fgySUlyWKo+TS31rxU8kX3DDM4= demo@example.com`
	publicKeyECDSA521 = `ecdsa-sha2-nistp521 AAAAE2VjZHNhLXNoYTItbmlzdHA1MjEAAAAIbmlzdHA1MjEAAACFBAEg73ORD4J3FV2CrL01gLSKREO2EHrZPlJCOeDL5OKD3M1GCHv3q8O452RW49Aw+8zFFFU5u6d1Ys3Qsj05zdaQwQDt/D3ceWLGVkWiKyLPQStfn0GGOZh3YFKEw5XmeW9jh6xudEHlKs4Pfv2FrroaUKZvM2SlxR/feOK0tCQyq3MN/g== demo@example.com`

	// testPublicKeyInstall is the public key that is installed in the
	// admin account's authorized_keys
	testPublicKeyInstall = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC9i+hFxZHGo6KblVme4zrAcJstR6I0PTJozW286X4WyvPnkMYDQ5mnhEYC7UWCvjoTWbPEXPX7NjhRtwQTGD67bV+lrxgfyzK1JZbUXK4PwgKJvQD+XyyWYMzDgGSQY61KUSqCxymSm/9NZkPU3ElaQ9xQuTzPpztM4ROfb8f2Yv6/ZESZsTo0MTAkp8Pcy+WkioI/uJ1H7zqs0EA4OMY4aDJRu0UtP4rTVeYNEAuRXdX+eH4aW3KMvhzpFTjMbaJHJXlEeUm2SaX5TNQyTOvghCeQILfYIL/Ca2ij8iwCmulwdV6eQGfd4VDu40PvSnmfoaE38o6HaPnX0kUcnKiT"

	dockerImageTagSupportsRSA1   = "8.1_p1-r0-ls20"
	dockerImageTagSupportsNoRSA1 = "8.4_p1-r3-ls48"
)

var ctx = context.Background()

func prepareTestContainer(t *testing.T, tag, caPublicKeyPEM string) (func(), string) {
	if tag == "" {
		tag = dockerImageTagSupportsNoRSA1
	}
	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ContainerName: "openssh",
		ImageRepo:     "docker.mirror.hashicorp.services/linuxserver/openssh-server",
		ImageTag:      tag,
		Env: []string{
			"DOCKER_MODS=linuxserver/mods:openssh-server-openssh-client",
			"PUBLIC_KEY=" + testPublicKeyInstall,
			"SUDO_ACCESS=true",
			"USER_NAME=vaultssh",
		},
		Ports: []string{"2222/tcp"},
	})
	if err != nil {
		t.Fatalf("Could not start local ssh docker container: %s", err)
	}

	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		ipaddr, err := net.ResolveIPAddr("ip", host)
		if err != nil {
			return nil, err
		}
		sshAddress := fmt.Sprintf("%s:%d", ipaddr.String(), port)

		signer, err := ssh.ParsePrivateKey([]byte(testSharedPrivateKey))
		if err != nil {
			return nil, err
		}

		// Install util-linux for non-busybox flock that supports timeout option
		err = testSSH("vaultssh", sshAddress, ssh.PublicKeys(signer), fmt.Sprintf(`
			set -e;
			sudo ln -s /config /home/vaultssh
			sudo apk add util-linux;
			echo "LogLevel DEBUG" | sudo tee -a /config/ssh_host_keys/sshd_config;
			echo "TrustedUserCAKeys /config/ssh_host_keys/trusted-user-ca-keys.pem" | sudo tee -a /config/ssh_host_keys/sshd_config;
			kill -HUP $(cat /config/sshd.pid)
			echo "%s" | sudo tee /config/ssh_host_keys/trusted-user-ca-keys.pem
		`, caPublicKeyPEM))
		if err != nil {
			return nil, err
		}

		return docker.NewServiceHostPort(ipaddr.String(), port), nil
	})
	if err != nil {
		t.Fatalf("Could not start docker ssh server: %s", err)
	}
	return svc.Cleanup, svc.Config.Address()
}

func testSSH(user, host string, auth ssh.AuthMethod, command string) error {
	client, err := ssh.Dial("tcp", host, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("unable to dial sshd to host %q: %v", host, err)
	}
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("unable to create sshd session to host %q: %v", host, err)
	}
	var stderr bytes.Buffer
	session.Stderr = &stderr
	defer session.Close()
	err = session.Run(command)
	if err != nil {
		return fmt.Errorf("command %v failed, error: %v, stderr: %v", command, err, stderr.String())
	}
	return nil
}

func TestBackend_AllowedUsers(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	err = b.Setup(context.Background(), config)
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

	resp, err := b.HandleRequest(context.Background(), roleReq)
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

	resp, err = b.HandleRequest(context.Background(), credsReq)
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
	resp, err = b.HandleRequest(context.Background(), credsReq)
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
	resp, err = b.HandleRequest(context.Background(), credsReq)
	if err != nil || resp == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("expected failure: resp:%#v err:%s", resp, err)
	}

	delete(roleData, "allowed_users")
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to create role: resp:%#v err:%s", resp, err)
	}

	credsData["username"] = "ubuntu"
	resp, err = b.HandleRequest(context.Background(), credsReq)
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
	resp, err = b.HandleRequest(context.Background(), credsReq)
	if err != nil || resp == nil || (resp != nil && !resp.IsError()) {
		t.Fatalf("expected failure: resp:%#v err:%s", resp, err)
	}

	roleData["allowed_users"] = "*"
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to create role: resp:%#v err:%s", resp, err)
	}

	resp, err = b.HandleRequest(context.Background(), credsReq)
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

func TestBackend_AllowedDomainsTemplate(t *testing.T) {
	testAllowedDomainsTemplate := "{{ identity.entity.metadata.ssh_username }}.example.com"
	expectedValidPrincipal := "foo." + testUserName + ".example.com"
	testAllowedPrincipalsTemplate(
		t, testAllowedDomainsTemplate,
		expectedValidPrincipal,
		map[string]string{
			"ssh_username": testUserName,
		},
		map[string]interface{}{
			"key_type":                 testCaKeyType,
			"algorithm_signer":         "rsa-sha2-256",
			"allow_host_certificates":  true,
			"allow_subdomains":         true,
			"allowed_domains":          testAllowedDomainsTemplate,
			"allowed_domains_template": true,
		},
		map[string]interface{}{
			"cert_type":        "host",
			"public_key":       testCAPublicKey,
			"valid_principals": expectedValidPrincipal,
		},
	)
}

func TestBackend_AllowedUsersTemplate(t *testing.T) {
	testAllowedUsersTemplate(t,
		"{{ identity.entity.metadata.ssh_username }}",
		testUserName, map[string]string{
			"ssh_username": testUserName,
		},
	)
}

func TestBackend_MultipleAllowedUsersTemplate(t *testing.T) {
	testAllowedUsersTemplate(t,
		"{{ identity.entity.metadata.ssh_username }}",
		testUserName, map[string]string{
			"ssh_username": testMultiUserName,
		},
	)
}

func TestBackend_AllowedUsersTemplate_WithStaticPrefix(t *testing.T) {
	testAllowedUsersTemplate(t,
		"ssh-{{ identity.entity.metadata.ssh_username }}",
		"ssh-"+testUserName, map[string]string{
			"ssh_username": testUserName,
		},
	)
}

func TestBackend_DefaultUserTemplate(t *testing.T) {
	testDefaultUserTemplate(t,
		"{{ identity.entity.metadata.ssh_username }}",
		testUserName,
		map[string]string{
			"ssh_username": testUserName,
		},
	)
}

func TestBackend_DefaultUserTemplate_WithStaticPrefix(t *testing.T) {
	testDefaultUserTemplate(t,
		"user-{{ identity.entity.metadata.ssh_username }}",
		"user-"+testUserName,
		map[string]string{
			"ssh_username": testUserName,
		},
	)
}

func TestBackend_DefaultUserTemplateFalse_AllowedUsersTemplateTrue(t *testing.T) {
	cluster, userpassToken := getSshCaTestCluster(t, testUserName)
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// set metadata "ssh_username" to userpass username
	tokenLookupResponse, err := client.Logical().Write("/auth/token/lookup", map[string]interface{}{
		"token": userpassToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	entityID := tokenLookupResponse.Data["entity_id"].(string)
	_, err = client.Logical().Write("/identity/entity/id/"+entityID, map[string]interface{}{
		"metadata": map[string]string{
			"ssh_username": testUserName,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("ssh/roles/my-role", map[string]interface{}{
		"key_type":                testCaKeyType,
		"allow_user_certificates": true,
		"default_user":            "{{identity.entity.metadata.ssh_username}}",
		// disable user templating but not allowed_user_template and the request should fail
		"default_user_template":  false,
		"allowed_users":          "{{identity.entity.metadata.ssh_username}}",
		"allowed_users_template": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// sign SSH key as userpass user
	client.SetToken(userpassToken)
	_, err = client.Logical().Write("ssh/sign/my-role", map[string]interface{}{
		"public_key": testCAPublicKey,
	})
	if err == nil {
		t.Errorf("signing request should fail when default_user is not in the allowed_users list, because allowed_users_template is true and default_user_template is not")
	}

	expectedErrStr := "{{identity.entity.metadata.ssh_username}} is not a valid value for valid_principals"
	if !strings.Contains(err.Error(), expectedErrStr) {
		t.Errorf("expected error to include %q but it was: %q", expectedErrStr, err.Error())
	}
}

func TestBackend_DefaultUserTemplateFalse_AllowedUsersTemplateFalse(t *testing.T) {
	cluster, userpassToken := getSshCaTestCluster(t, testUserName)
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// set metadata "ssh_username" to userpass username
	tokenLookupResponse, err := client.Logical().Write("/auth/token/lookup", map[string]interface{}{
		"token": userpassToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	entityID := tokenLookupResponse.Data["entity_id"].(string)
	_, err = client.Logical().Write("/identity/entity/id/"+entityID, map[string]interface{}{
		"metadata": map[string]string{
			"ssh_username": testUserName,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("ssh/roles/my-role", map[string]interface{}{
		"key_type":                testCaKeyType,
		"allow_user_certificates": true,
		"default_user":            "{{identity.entity.metadata.ssh_username}}",
		"default_user_template":   false,
		"allowed_users":           "{{identity.entity.metadata.ssh_username}}",
		"allowed_users_template":  false,
	})
	if err != nil {
		t.Fatal(err)
	}

	// sign SSH key as userpass user
	client.SetToken(userpassToken)
	signResponse, err := client.Logical().Write("ssh/sign/my-role", map[string]interface{}{
		"public_key": testCAPublicKey,
	})
	if err != nil {
		t.Fatal(err)
	}

	// check for the expected valid principals of certificate
	signedKey := signResponse.Data["signed_key"].(string)
	key, _ := base64.StdEncoding.DecodeString(strings.Split(signedKey, " ")[1])
	parsedKey, err := ssh.ParsePublicKey(key)
	if err != nil {
		t.Fatal(err)
	}
	actualPrincipals := parsedKey.(*ssh.Certificate).ValidPrincipals
	if len(actualPrincipals) < 1 {
		t.Fatal(
			fmt.Sprintf("No ValidPrincipals returned: should have been %v",
				[]string{"{{identity.entity.metadata.ssh_username}}"}),
		)
	}
	if len(actualPrincipals) > 1 {
		t.Error(
			fmt.Sprintf("incorrect number ValidPrincipals, expected only 1: %v should be %v",
				actualPrincipals, []string{"{{identity.entity.metadata.ssh_username}}"}),
		)
	}
	if actualPrincipals[0] != "{{identity.entity.metadata.ssh_username}}" {
		t.Fatal(
			fmt.Sprintf("incorrect ValidPrincipals: %v should be %v",
				actualPrincipals, []string{"{{identity.entity.metadata.ssh_username}}"}),
		)
	}
}

func newTestingFactory(t *testing.T) func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		defaultLeaseTTLVal := 2 * time.Minute
		maxLeaseTTLVal := 10 * time.Minute
		return Factory(context.Background(), &logical.BackendConfig{
			Logger:      corehelpers.NewTestLogger(t),
			StorageView: &logical.InmemStorage{},
			System: &logical.StaticSystemView{
				DefaultLeaseTTLVal: defaultLeaseTTLVal,
				MaxLeaseTTLVal:     maxLeaseTTLVal,
			},
		})
	}
}

func TestSSHBackend_Lookup(t *testing.T) {
	testOTPRoleData := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
		"cidr_list":    testCIDRList,
	}
	data := map[string]interface{}{
		"ip": testIP,
	}
	resp1 := []string(nil)
	resp2 := []string{testOTPRoleName}
	resp3 := []string{testAtRoleName}
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalFactory: newTestingFactory(t),
		Steps: []logicaltest.TestStep{
			testLookupRead(t, data, resp1),
			testRoleWrite(t, testOTPRoleName, testOTPRoleData),
			testLookupRead(t, data, resp2),
			testRoleDelete(t, testOTPRoleName),
			testLookupRead(t, data, resp1),
			testRoleWrite(t, testAtRoleName, testOTPRoleData),
			testLookupRead(t, data, resp3),
			testRoleDelete(t, testAtRoleName),
			testLookupRead(t, data, resp1),
		},
	})
}

func TestSSHBackend_RoleList(t *testing.T) {
	testOTPRoleData := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
		"cidr_list":    testCIDRList,
	}
	resp1 := map[string]interface{}{}
	resp2 := map[string]interface{}{
		"keys": []string{testOTPRoleName},
		"key_info": map[string]interface{}{
			testOTPRoleName: map[string]interface{}{
				"key_type": testOTPKeyType,
			},
		},
	}
	resp3 := map[string]interface{}{
		"keys": []string{testAtRoleName, testOTPRoleName},
		"key_info": map[string]interface{}{
			testOTPRoleName: map[string]interface{}{
				"key_type": testOTPKeyType,
			},
			testAtRoleName: map[string]interface{}{
				"key_type": testOTPKeyType,
			},
		},
	}
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalFactory: newTestingFactory(t),
		Steps: []logicaltest.TestStep{
			testRoleList(t, resp1),
			testRoleWrite(t, testOTPRoleName, testOTPRoleData),
			testRoleList(t, resp2),
			testRoleWrite(t, testAtRoleName, testOTPRoleData),
			testRoleList(t, resp3),
			testRoleDelete(t, testAtRoleName),
			testRoleList(t, resp2),
			testRoleDelete(t, testOTPRoleName),
			testRoleList(t, resp1),
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
		LogicalFactory: newTestingFactory(t),
		Steps: []logicaltest.TestStep{
			testRoleWrite(t, testOTPRoleName, testOTPRoleData),
			testRoleRead(t, testOTPRoleName, respOTPRoleData),
			testRoleDelete(t, testOTPRoleName),
			testRoleRead(t, testOTPRoleName, nil),
			testRoleWrite(t, testAtRoleName, testOTPRoleData),
			testRoleRead(t, testAtRoleName, respOTPRoleData),
			testRoleDelete(t, testAtRoleName),
			testRoleRead(t, testAtRoleName, nil),
		},
	})
}

func TestSSHBackend_OTPCreate(t *testing.T) {
	cleanup, sshAddress := prepareTestContainer(t, "", "")
	defer func() {
		if !t.Failed() {
			cleanup()
		}
	}()

	host, port, err := net.SplitHostPort(sshAddress)
	if err != nil {
		t.Fatal(err)
	}

	testOTPRoleData := map[string]interface{}{
		"key_type":     testOTPKeyType,
		"default_user": testUserName,
		"cidr_list":    host + "/32",
		"port":         port,
	}
	data := map[string]interface{}{
		"username": testUserName,
		"ip":       host,
	}
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalFactory: newTestingFactory(t),
		Steps: []logicaltest.TestStep{
			testRoleWrite(t, testOTPRoleName, testOTPRoleData),
			testCredsWrite(t, testOTPRoleName, data, false, sshAddress),
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
		LogicalFactory: newTestingFactory(t),
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
	req1 := map[string]interface{}{
		"roles": testOTPRoleName,
	}
	resp1 := map[string]interface{}{
		"roles": []string{testOTPRoleName},
	}
	resp2 := map[string]interface{}{
		"roles": []string{testOTPRoleName},
	}
	resp3 := map[string]interface{}{
		"roles": []string{},
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalFactory: newTestingFactory(t),
		Steps: []logicaltest.TestStep{
			testRoleWrite(t, testOTPRoleName, testOTPRoleData),
			testConfigZeroAddressWrite(t, req1),
			testConfigZeroAddressRead(t, resp1),
			testConfigZeroAddressRead(t, resp2),
			testConfigZeroAddressRead(t, resp1),
			testRoleDelete(t, testOTPRoleName),
			testConfigZeroAddressRead(t, resp3),
			testConfigZeroAddressDelete(t),
		},
	})
}

func TestSSHBackend_CredsForZeroAddressRoles_otp(t *testing.T) {
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
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalFactory: newTestingFactory(t),
		Steps: []logicaltest.TestStep{
			testRoleWrite(t, testOTPRoleName, otpRoleData),
			testCredsWrite(t, testOTPRoleName, data, true, ""),
			testConfigZeroAddressWrite(t, req1),
			testCredsWrite(t, testOTPRoleName, data, false, ""),
			testConfigZeroAddressDelete(t),
			testCredsWrite(t, testOTPRoleName, data, true, ""),
		},
	})
}

func TestSSHBackend_CA(t *testing.T) {
	testCases := []struct {
		name         string
		tag          string
		caPublicKey  string
		caPrivateKey string
		algoSigner   string
		expectError  bool
	}{
		{
			"RSAKey_EmptyAlgoSigner_ImageSupportsRSA1",
			dockerImageTagSupportsRSA1,
			testCAPublicKey,
			testCAPrivateKey,
			"",
			false,
		},
		{
			"RSAKey_EmptyAlgoSigner_ImageSupportsNoRSA1",
			dockerImageTagSupportsNoRSA1,
			testCAPublicKey,
			testCAPrivateKey,
			"",
			false,
		},
		{
			"RSAKey_DefaultAlgoSigner_ImageSupportsRSA1",
			dockerImageTagSupportsRSA1,
			testCAPublicKey,
			testCAPrivateKey,
			"default",
			false,
		},
		{
			"RSAKey_DefaultAlgoSigner_ImageSupportsNoRSA1",
			dockerImageTagSupportsNoRSA1,
			testCAPublicKey,
			testCAPrivateKey,
			"default",
			false,
		},
		{
			"RSAKey_RSA1AlgoSigner_ImageSupportsRSA1",
			dockerImageTagSupportsRSA1,
			testCAPublicKey,
			testCAPrivateKey,
			ssh.SigAlgoRSA,
			false,
		},
		{
			"RSAKey_RSA1AlgoSigner_ImageSupportsNoRSA1",
			dockerImageTagSupportsNoRSA1,
			testCAPublicKey,
			testCAPrivateKey,
			ssh.SigAlgoRSA,
			true,
		},
		{
			"RSAKey_RSASHA2256AlgoSigner_ImageSupportsRSA1",
			dockerImageTagSupportsRSA1,
			testCAPublicKey,
			testCAPrivateKey,
			ssh.SigAlgoRSASHA2256,
			false,
		},
		{
			"RSAKey_RSASHA2256AlgoSigner_ImageSupportsNoRSA1",
			dockerImageTagSupportsNoRSA1,
			testCAPublicKey,
			testCAPrivateKey,
			ssh.SigAlgoRSASHA2256,
			false,
		},
		{
			"ed25519Key_EmptyAlgoSigner_ImageSupportsRSA1",
			dockerImageTagSupportsRSA1,
			testCAPublicKeyEd25519,
			testCAPrivateKeyEd25519,
			"",
			false,
		},
		{
			"ed25519Key_EmptyAlgoSigner_ImageSupportsNoRSA1",
			dockerImageTagSupportsNoRSA1,
			testCAPublicKeyEd25519,
			testCAPrivateKeyEd25519,
			"",
			false,
		},
		{
			"ed25519Key_RSA1AlgoSigner_ImageSupportsRSA1",
			dockerImageTagSupportsRSA1,
			testCAPublicKeyEd25519,
			testCAPrivateKeyEd25519,
			ssh.SigAlgoRSA,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testSSHBackend_CA(t, tc.tag, tc.caPublicKey, tc.caPrivateKey, tc.algoSigner, tc.expectError)
		})
	}
}

func testSSHBackend_CA(t *testing.T, dockerImageTag, caPublicKey, caPrivateKey, algorithmSigner string, expectError bool) {
	cleanup, sshAddress := prepareTestContainer(t, dockerImageTag, caPublicKey)
	defer cleanup()
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testKeyToSignPrivate := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAwn1V2xd/EgJXIY53fBTtc20k/ajekqQngvkpFSwNHW63XNEQK8Ll
FOCyGXoje9DUGxnYs3F/ohfsBBWkLNfU7fiENdSJL1pbkAgJ+2uhV9sLZjvYhikrXWoyJX
LDKfY12LjpcBS2HeLMT04laZ/xSJrOBEJHGzHyr2wUO0NUQUQPUODAFhnHKgvvA4Uu79UY
gcdThF4w83+EAnE4JzBZMKPMjzy4u1C0R/LoD8DuapHwX6NGWdEUvUZZ+XRcIWeCOvR0ne
qGBRH35k1Mv7k65d7kkE0uvM5Z36erw3tdoszxPYf7AKnO1DpeU2uwMcym6xNwfwynKjhL
qL/Mgi4uRwAAA8iAsY0zgLGNMwAAAAdzc2gtcnNhAAABAQDCfVXbF38SAlchjnd8FO1zbS
T9qN6SpCeC+SkVLA0dbrdc0RArwuUU4LIZeiN70NQbGdizcX+iF+wEFaQs19Tt+IQ11Ikv
WluQCAn7a6FX2wtmO9iGKStdajIlcsMp9jXYuOlwFLYd4sxPTiVpn/FIms4EQkcbMfKvbB
Q7Q1RBRA9Q4MAWGccqC+8DhS7v1RiBx1OEXjDzf4QCcTgnMFkwo8yPPLi7ULRH8ugPwO5q
kfBfo0ZZ0RS9Rln5dFwhZ4I69HSd6oYFEffmTUy/uTrl3uSQTS68zlnfp6vDe12izPE9h/
sAqc7UOl5Ta7AxzKbrE3B/DKcqOEuov8yCLi5HAAAAAwEAAQAAAQABns2yT5XNbpuPOgKg
1APObGBchKWmDxwNKUpAVOefEScR7OP3mV4TOHQDZlMZWvoJZ8O4av+nOA/NUOjXPs0VVn
azhBvIezY8EvUSVSk49Cg6J9F7/KfR1WqpiTU7CkQUlCXNuz5xLUyKdJo3MQ/vjOqeenbh
MR9Wes4IWF1BVe4VOD6lxRsjwuIieIgmScW28FFh2rgsEfO2spzZ3AWOGExw+ih757hFz5
4A2fhsQXP8m3r8m7iiqcjTLWXdxTUk4zot2kZEjbI4Avk0BL+wVeFq6f/y+G+g5edqSo7j
uuSgzbUQtA9PMnGxhrhU2Ob7n3VGdya7WbGZkaKP8zJhAAAAgQC3bJurmOSLIi3KVhp7lD
/FfxwXHwVBFALCgq7EyNlkTz6RDoMFM4eOTRMDvsgWxT+bSB8R8eg1sfgY8rkHOuvTAVI5
3oEYco3H7NWE9X8Zt0lyhO1uaE49EENNSQ8hY7R3UIw5becyI+7ZZxs9HkBgCQCZzSjzA+
SIyAoMKM261AAAAIEA+PCkcDRp3J0PaoiuetXSlWZ5WjP3CtwT2xrvEX9x+ZsDgXCDYQ5T
osxvEKOGSfIrHUUhzZbFGvqWyfrziPe9ypJrtCM7RJT/fApBXnbWFcDZzWamkQvohst+0w
XHYCmNoJ6/Y+roLv3pzyFUmqRNcrQaohex7TZmsvHJT513UakAAACBAMgBXxH8DyNYdniX
mIXEto4GqMh4rXdNwCghfpyWdJE6vCyDt7g7bYMq7AQ2ynSKRtQDT/ZgQNfSbilUq3iXz7
xNZn5U9ndwFs90VmEpBup/PmhfX+Gwt5hQZLbkKZcgQ9XrhSKdMxVm1yy/fk0U457enlz5
cKumubUxOfFdy1ZvAAAAEm5jY0BtYnAudWJudC5sb2NhbA==
-----END OPENSSH PRIVATE KEY-----
`
	testKeyToSignPublic := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDCfVXbF38SAlchjnd8FO1zbST9qN6SpCeC+SkVLA0dbrdc0RArwuUU4LIZeiN70NQbGdizcX+iF+wEFaQs19Tt+IQ11IkvWluQCAn7a6FX2wtmO9iGKStdajIlcsMp9jXYuOlwFLYd4sxPTiVpn/FIms4EQkcbMfKvbBQ7Q1RBRA9Q4MAWGccqC+8DhS7v1RiBx1OEXjDzf4QCcTgnMFkwo8yPPLi7ULRH8ugPwO5qkfBfo0ZZ0RS9Rln5dFwhZ4I69HSd6oYFEffmTUy/uTrl3uSQTS68zlnfp6vDe12izPE9h/sAqc7UOl5Ta7AxzKbrE3B/DKcqOEuov8yCLi5H `

	roleOptions := map[string]interface{}{
		"allow_user_certificates": true,
		"allowed_users":           "*",
		"default_extensions": []map[string]string{
			{
				"permit-pty": "",
			},
		},
		"key_type":     "ca",
		"default_user": testUserName,
		"ttl":          "30m0s",
	}
	if algorithmSigner != "" {
		roleOptions["algorithm_signer"] = algorithmSigner
	}
	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(caPublicKey, caPrivateKey),
			testRoleWrite(t, "testcarole", roleOptions),
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/testcarole",
				ErrorOk:   expectError,
				Data: map[string]interface{}{
					"public_key":       testKeyToSignPublic,
					"valid_principals": testUserName,
				},

				Check: func(resp *logical.Response) error {
					// Tolerate nil response if an error was expected
					if expectError && resp == nil {
						return nil
					}

					signedKey := strings.TrimSpace(resp.Data["signed_key"].(string))
					if signedKey == "" {
						return errors.New("no signed key in response")
					}

					privKey, err := ssh.ParsePrivateKey([]byte(testKeyToSignPrivate))
					if err != nil {
						return fmt.Errorf("error parsing private key: %v", err)
					}

					parsedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(signedKey))
					if err != nil {
						return fmt.Errorf("error parsing signed key: %v", err)
					}
					certSigner, err := ssh.NewCertSigner(parsedKey.(*ssh.Certificate), privKey)
					if err != nil {
						return err
					}

					err = testSSH(testUserName, sshAddress, ssh.PublicKeys(certSigner), "date")
					if expectError && err == nil {
						return fmt.Errorf("expected error but got none")
					}
					if !expectError && err != nil {
						return err
					}

					return nil
				},
			},
			testIssueCert("testcarole", "ec", testUserName, sshAddress, expectError),
			testIssueCert("testcarole", "ed25519", testUserName, sshAddress, expectError),
			testIssueCert("testcarole", "rsa", testUserName, sshAddress, expectError),
		},
	}

	logicaltest.Test(t, testCase)
}

func testIssueCert(role string, keyType string, testUserName string, sshAddress string, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "issue/" + role,
		ErrorOk:   expectError,
		Data: map[string]interface{}{
			"key_type":         keyType,
			"valid_principals": testUserName,
		},

		Check: func(resp *logical.Response) error {
			// Tolerate nil response if an error was expected
			if expectError && resp == nil {
				return nil
			}

			signedKey := strings.TrimSpace(resp.Data["signed_key"].(string))
			if signedKey == "" {
				return errors.New("no signed key in response")
			}

			privKey, err := ssh.ParsePrivateKey([]byte(resp.Data["private_key"].(string)))
			if err != nil {
				return fmt.Errorf("error parsing private key: %v", err)
			}

			parsedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(signedKey))
			if err != nil {
				return fmt.Errorf("error parsing signed key: %v", err)
			}
			certSigner, err := ssh.NewCertSigner(parsedKey.(*ssh.Certificate), privKey)
			if err != nil {
				return err
			}

			err = testSSH(testUserName, sshAddress, ssh.PublicKeys(certSigner), "date")
			if expectError && err == nil {
				return fmt.Errorf("expected error but got none")
			}
			if !expectError && err != nil {
				return err
			}

			return nil
		},
	}
}

func TestSSHBackend_CAUpgradeAlgorithmSigner(t *testing.T) {
	cleanup, sshAddress := prepareTestContainer(t, dockerImageTagSupportsRSA1, testCAPublicKey)
	defer cleanup()
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testKeyToSignPrivate := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAwn1V2xd/EgJXIY53fBTtc20k/ajekqQngvkpFSwNHW63XNEQK8Ll
FOCyGXoje9DUGxnYs3F/ohfsBBWkLNfU7fiENdSJL1pbkAgJ+2uhV9sLZjvYhikrXWoyJX
LDKfY12LjpcBS2HeLMT04laZ/xSJrOBEJHGzHyr2wUO0NUQUQPUODAFhnHKgvvA4Uu79UY
gcdThF4w83+EAnE4JzBZMKPMjzy4u1C0R/LoD8DuapHwX6NGWdEUvUZZ+XRcIWeCOvR0ne
qGBRH35k1Mv7k65d7kkE0uvM5Z36erw3tdoszxPYf7AKnO1DpeU2uwMcym6xNwfwynKjhL
qL/Mgi4uRwAAA8iAsY0zgLGNMwAAAAdzc2gtcnNhAAABAQDCfVXbF38SAlchjnd8FO1zbS
T9qN6SpCeC+SkVLA0dbrdc0RArwuUU4LIZeiN70NQbGdizcX+iF+wEFaQs19Tt+IQ11Ikv
WluQCAn7a6FX2wtmO9iGKStdajIlcsMp9jXYuOlwFLYd4sxPTiVpn/FIms4EQkcbMfKvbB
Q7Q1RBRA9Q4MAWGccqC+8DhS7v1RiBx1OEXjDzf4QCcTgnMFkwo8yPPLi7ULRH8ugPwO5q
kfBfo0ZZ0RS9Rln5dFwhZ4I69HSd6oYFEffmTUy/uTrl3uSQTS68zlnfp6vDe12izPE9h/
sAqc7UOl5Ta7AxzKbrE3B/DKcqOEuov8yCLi5HAAAAAwEAAQAAAQABns2yT5XNbpuPOgKg
1APObGBchKWmDxwNKUpAVOefEScR7OP3mV4TOHQDZlMZWvoJZ8O4av+nOA/NUOjXPs0VVn
azhBvIezY8EvUSVSk49Cg6J9F7/KfR1WqpiTU7CkQUlCXNuz5xLUyKdJo3MQ/vjOqeenbh
MR9Wes4IWF1BVe4VOD6lxRsjwuIieIgmScW28FFh2rgsEfO2spzZ3AWOGExw+ih757hFz5
4A2fhsQXP8m3r8m7iiqcjTLWXdxTUk4zot2kZEjbI4Avk0BL+wVeFq6f/y+G+g5edqSo7j
uuSgzbUQtA9PMnGxhrhU2Ob7n3VGdya7WbGZkaKP8zJhAAAAgQC3bJurmOSLIi3KVhp7lD
/FfxwXHwVBFALCgq7EyNlkTz6RDoMFM4eOTRMDvsgWxT+bSB8R8eg1sfgY8rkHOuvTAVI5
3oEYco3H7NWE9X8Zt0lyhO1uaE49EENNSQ8hY7R3UIw5becyI+7ZZxs9HkBgCQCZzSjzA+
SIyAoMKM261AAAAIEA+PCkcDRp3J0PaoiuetXSlWZ5WjP3CtwT2xrvEX9x+ZsDgXCDYQ5T
osxvEKOGSfIrHUUhzZbFGvqWyfrziPe9ypJrtCM7RJT/fApBXnbWFcDZzWamkQvohst+0w
XHYCmNoJ6/Y+roLv3pzyFUmqRNcrQaohex7TZmsvHJT513UakAAACBAMgBXxH8DyNYdniX
mIXEto4GqMh4rXdNwCghfpyWdJE6vCyDt7g7bYMq7AQ2ynSKRtQDT/ZgQNfSbilUq3iXz7
xNZn5U9ndwFs90VmEpBup/PmhfX+Gwt5hQZLbkKZcgQ9XrhSKdMxVm1yy/fk0U457enlz5
cKumubUxOfFdy1ZvAAAAEm5jY0BtYnAudWJudC5sb2NhbA==
-----END OPENSSH PRIVATE KEY-----
`
	testKeyToSignPublic := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDCfVXbF38SAlchjnd8FO1zbST9qN6SpCeC+SkVLA0dbrdc0RArwuUU4LIZeiN70NQbGdizcX+iF+wEFaQs19Tt+IQ11IkvWluQCAn7a6FX2wtmO9iGKStdajIlcsMp9jXYuOlwFLYd4sxPTiVpn/FIms4EQkcbMfKvbBQ7Q1RBRA9Q4MAWGccqC+8DhS7v1RiBx1OEXjDzf4QCcTgnMFkwo8yPPLi7ULRH8ugPwO5qkfBfo0ZZ0RS9Rln5dFwhZ4I69HSd6oYFEffmTUy/uTrl3uSQTS68zlnfp6vDe12izPE9h/sAqc7UOl5Ta7AxzKbrE3B/DKcqOEuov8yCLi5H `

	// Old role entries between 1.4.3 and 1.5.2 had algorithm_signer default to
	// ssh-rsa if not provided.
	roleOptionsOldEntry := map[string]interface{}{
		"allow_user_certificates": true,
		"allowed_users":           "*",
		"default_extensions": []map[string]string{
			{
				"permit-pty": "",
			},
		},
		"key_type":         "ca",
		"default_user":     testUserName,
		"ttl":              "30m0s",
		"algorithm_signer": ssh.SigAlgoRSA,
	}

	// Upgrade entry by overwriting algorithm_signer with an empty value
	roleOptionsUpgradedEntry := map[string]interface{}{
		"allow_user_certificates": true,
		"allowed_users":           "*",
		"default_extensions": []map[string]string{
			{
				"permit-pty": "",
			},
		},
		"key_type":         "ca",
		"default_user":     testUserName,
		"ttl":              "30m0s",
		"algorithm_signer": "",
	}

	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(testCAPublicKey, testCAPrivateKey),
			testRoleWrite(t, "testcarole", roleOptionsOldEntry),
			testRoleWrite(t, "testcarole", roleOptionsUpgradedEntry),
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/testcarole",
				ErrorOk:   false,
				Data: map[string]interface{}{
					"public_key":       testKeyToSignPublic,
					"valid_principals": testUserName,
				},

				Check: func(resp *logical.Response) error {
					signedKey := strings.TrimSpace(resp.Data["signed_key"].(string))
					if signedKey == "" {
						return errors.New("no signed key in response")
					}

					privKey, err := ssh.ParsePrivateKey([]byte(testKeyToSignPrivate))
					if err != nil {
						return fmt.Errorf("error parsing private key: %v", err)
					}

					parsedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(signedKey))
					if err != nil {
						return fmt.Errorf("error parsing signed key: %v", err)
					}
					certSigner, err := ssh.NewCertSigner(parsedKey.(*ssh.Certificate), privKey)
					if err != nil {
						return err
					}

					err = testSSH(testUserName, sshAddress, ssh.PublicKeys(certSigner), "date")
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	logicaltest.Test(t, testCase)
}

func TestBackend_AbleToRetrievePublicKey(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(testCAPublicKey, testCAPrivateKey),

			{
				Operation:       logical.ReadOperation,
				Path:            "public_key",
				Unauthenticated: true,

				Check: func(resp *logical.Response) error {
					key := string(resp.Data["http_raw_body"].([]byte))

					if key != testCAPublicKey {
						return fmt.Errorf("public_key incorrect. Expected %v, actual %v", testCAPublicKey, key)
					}

					return nil
				},
			},
		},
	}

	logicaltest.Test(t, testCase)
}

func TestBackend_AbleToAutoGenerateSigningKeys(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	var expectedPublicKey string
	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			{
				Operation: logical.UpdateOperation,
				Path:      "config/ca",
				Check: func(resp *logical.Response) error {
					if resp.Data["public_key"].(string) == "" {
						return fmt.Errorf("public_key empty")
					}
					expectedPublicKey = resp.Data["public_key"].(string)
					return nil
				},
			},

			{
				Operation:       logical.ReadOperation,
				Path:            "public_key",
				Unauthenticated: true,

				Check: func(resp *logical.Response) error {
					key := string(resp.Data["http_raw_body"].([]byte))

					if key == "" {
						return fmt.Errorf("public_key empty. Expected not empty, actual %s", key)
					}
					if key != expectedPublicKey {
						return fmt.Errorf("public_key mismatch. Expected %s, actual %s", expectedPublicKey, key)
					}

					return nil
				},
			},
		},
	}

	logicaltest.Test(t, testCase)
}

func TestBackend_ValidPrincipalsValidatedForHostCertificates(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(testCAPublicKey, testCAPrivateKey),

			createRoleStep("testing", map[string]interface{}{
				"key_type":                "ca",
				"allow_host_certificates": true,
				"allowed_domains":         "example.com,example.org",
				"allow_subdomains":        true,
				"default_critical_options": map[string]interface{}{
					"option": "value",
				},
				"default_extensions": map[string]interface{}{
					"extension": "extended",
				},
			}),

			signCertificateStep("testing", "vault-root-22608f5ef173aabf700797cb95c5641e792698ec6380e8e1eb55523e39aa5e51", ssh.HostCert, []string{"dummy.example.org", "second.example.com"}, map[string]string{
				"option": "value",
			}, map[string]string{
				"extension": "extended",
			},
				2*time.Hour, map[string]interface{}{
					"public_key":       publicKey2,
					"ttl":              "2h",
					"cert_type":        "host",
					"valid_principals": "dummy.example.org,second.example.com",
				}),
		},
	}

	logicaltest.Test(t, testCase)
}

func TestBackend_OptionsOverrideDefaults(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(testCAPublicKey, testCAPrivateKey),

			createRoleStep("testing", map[string]interface{}{
				"key_type":                 "ca",
				"allowed_users":            "tuber",
				"default_user":             "tuber",
				"allow_user_certificates":  true,
				"allowed_critical_options": "option,secondary",
				"allowed_extensions":       "extension,additional",
				"default_critical_options": map[string]interface{}{
					"option": "value",
				},
				"default_extensions": map[string]interface{}{
					"extension": "extended",
				},
			}),

			signCertificateStep("testing", "vault-root-22608f5ef173aabf700797cb95c5641e792698ec6380e8e1eb55523e39aa5e51", ssh.UserCert, []string{"tuber"}, map[string]string{
				"secondary": "value",
			}, map[string]string{
				"additional": "value",
			}, 2*time.Hour, map[string]interface{}{
				"public_key": publicKey2,
				"ttl":        "2h",
				"critical_options": map[string]interface{}{
					"secondary": "value",
				},
				"extensions": map[string]interface{}{
					"additional": "value",
				},
			}),
		},
	}

	logicaltest.Test(t, testCase)
}

func TestBackend_EmptyPrincipals(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}
	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(testCAPublicKey, testCAPrivateKey),
			createRoleStep("no_user_principals", map[string]interface{}{
				"key_type":                "ca",
				"allow_user_certificates": true,
				"allowed_user_key_lengths": map[string]interface{}{
					"rsa": 2048,
				},
				"allowed_users": "no_principals",
			}),
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/no_user_principals",
				Data: map[string]interface{}{
					"public_key": testCAPublicKey,
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp.Data["error"] != "empty valid principals not allowed by role" {
						return errors.New("expected empty valid principals not allowed by role")
					}
					return nil
				},
			},
			createRoleStep("no_host_principals", map[string]interface{}{
				"key_type":                "ca",
				"allow_host_certificates": true,
				"allowed_domains":         "*",
			}),
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/no_host_principals",
				Data: map[string]interface{}{
					"cert_type":  "host",
					"public_key": testCAPublicKeyEd25519,
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp.Data["error"] != "empty valid principals not allowed by role" {
						return errors.New("expected empty valid principals not allowed by role")
					}
					return nil
				},
			},
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/no_host_principals",
				Data: map[string]interface{}{
					"cert_type":        "host",
					"public_key":       testCAPublicKeyEd25519,
					"valid_principals": "example.com",
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp.Data["error"] != nil {
						return errors.New("expected no error")
					}
					return nil
				},
			},
		},
	}
	logicaltest.Test(t, testCase)
}

func TestBackend_AllowedUserKeyLengths(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}
	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(testCAPublicKey, testCAPrivateKey),
			createRoleStep("weakkey", map[string]interface{}{
				"key_type":                "ca",
				"allow_user_certificates": true,
				"allowed_user_key_lengths": map[string]interface{}{
					"rsa": 4096,
				},
				"allowed_users": "guest",
			}),
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/weakkey",
				Data: map[string]interface{}{
					"public_key": testCAPublicKey,
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp.Data["error"] != "public_key failed to meet the key requirements: key is of an invalid size: 2048" {
						return errors.New("a smaller key (2048) was allowed, when the minimum was set for 4096")
					}
					return nil
				},
			},
			createRoleStep("stdkey", map[string]interface{}{
				"key_type":                "ca",
				"allow_user_certificates": true,
				"allowed_user_key_lengths": map[string]interface{}{
					"rsa": 2048,
				},
				"allowed_users": "guest",
			}),
			// Pass with 2048 key
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/stdkey",
				Data: map[string]interface{}{
					"public_key":       testCAPublicKey,
					"valid_principals": "guest",
				},
			},
			// Fail with 4096 key
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/stdkey",
				Data: map[string]interface{}{
					"public_key":       publicKey4096,
					"valid_principals": "guest",
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp.Data["error"] != "public_key failed to meet the key requirements: key is of an invalid size: 4096" {
						return errors.New("a larger key (4096) was allowed, when the size was set for 2048")
					}
					return nil
				},
			},
			createRoleStep("multikey", map[string]interface{}{
				"key_type":                "ca",
				"allow_user_certificates": true,
				"allowed_user_key_lengths": map[string]interface{}{
					"rsa": []int{2048, 4096},
				},
				"allowed_users": "guest",
			}),
			// Pass with 2048-bit key
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/multikey",
				Data: map[string]interface{}{
					"public_key":       testCAPublicKey,
					"valid_principals": "guest",
				},
			},
			// Pass with 4096-bit key
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/multikey",
				Data: map[string]interface{}{
					"public_key":       publicKey4096,
					"valid_principals": "guest",
				},
			},
			// Fail with 3072-bit key
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/multikey",
				Data: map[string]interface{}{
					"public_key": publicKey3072,
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp.Data["error"] != "public_key failed to meet the key requirements: key is of an invalid size: 3072" {
						return errors.New("a larger key (3072) was allowed, when the size was set for 2048")
					}
					return nil
				},
			},
			// Fail with ECDSA key
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/multikey",
				Data: map[string]interface{}{
					"public_key":       publicKeyECDSA256,
					"valid_principals": "guest",
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp.Data["error"] != "public_key failed to meet the key requirements: key of type ecdsa is not allowed" {
						return errors.New("an ECDSA key was allowed under RSA-only policy")
					}
					return nil
				},
			},
			createRoleStep("ectypes", map[string]interface{}{
				"key_type":                "ca",
				"allow_user_certificates": true,
				"allowed_user_key_lengths": map[string]interface{}{
					"ec":                  []int{256},
					"ecdsa-sha2-nistp521": 0,
				},
				"allowed_users": "guest",
			}),
			// Pass with ECDSA P-256
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/ectypes",
				Data: map[string]interface{}{
					"public_key":       publicKeyECDSA256,
					"valid_principals": "guest",
				},
			},
			// Pass with ECDSA P-521
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/ectypes",
				Data: map[string]interface{}{
					"public_key":       publicKeyECDSA521,
					"valid_principals": "guest",
				},
			},
			// Fail with RSA key
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/ectypes",
				Data: map[string]interface{}{
					"public_key":       publicKey3072,
					"valid_principals": "guest",
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp.Data["error"] != "public_key failed to meet the key requirements: key of type rsa is not allowed" {
						return errors.New("an RSA key was allowed under ECDSA-only policy")
					}
					return nil
				},
			},
		},
	}

	logicaltest.Test(t, testCase)
}

func TestBackend_CustomKeyIDFormat(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(testCAPublicKey, testCAPrivateKey),

			createRoleStep("customrole", map[string]interface{}{
				"key_type":                 "ca",
				"key_id_format":            "{{role_name}}-{{token_display_name}}-{{public_key_hash}}",
				"allowed_users":            "tuber",
				"default_user":             "tuber",
				"allow_user_certificates":  true,
				"allowed_critical_options": "option,secondary",
				"allowed_extensions":       "extension,additional",
				"default_critical_options": map[string]interface{}{
					"option": "value",
				},
				"default_extensions": map[string]interface{}{
					"extension": "extended",
				},
			}),

			signCertificateStep("customrole", "customrole-root-22608f5ef173aabf700797cb95c5641e792698ec6380e8e1eb55523e39aa5e51", ssh.UserCert, []string{"tuber"}, map[string]string{
				"secondary": "value",
			}, map[string]string{
				"additional": "value",
			}, 2*time.Hour, map[string]interface{}{
				"public_key": publicKey2,
				"ttl":        "2h",
				"critical_options": map[string]interface{}{
					"secondary": "value",
				},
				"extensions": map[string]interface{}{
					"additional": "value",
				},
			}),
		},
	}

	logicaltest.Test(t, testCase)
}

func TestBackend_DisallowUserProvidedKeyIDs(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(testCAPublicKey, testCAPrivateKey),

			createRoleStep("testing", map[string]interface{}{
				"key_type":                "ca",
				"allow_user_key_ids":      false,
				"allow_user_certificates": true,
			}),
			{
				Operation: logical.UpdateOperation,
				Path:      "sign/testing",
				Data: map[string]interface{}{
					"public_key": publicKey2,
					"key_id":     "override",
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp.Data["error"] != "setting key_id is not allowed by role" {
						return errors.New("custom user key id was allowed even when 'allow_user_key_ids' is false")
					}
					return nil
				},
			},
		},
	}

	logicaltest.Test(t, testCase)
}

func TestBackend_DefExtTemplatingEnabled(t *testing.T) {
	cluster, userpassToken := getSshCaTestCluster(t, testUserName)
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// Get auth accessor for identity template.
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	userpassAccessor := auths["userpass/"].Accessor

	// Write SSH role.
	_, err = client.Logical().Write("ssh/roles/test", map[string]interface{}{
		"key_type":                    "ca",
		"allowed_extensions":          "login@zipzap.com",
		"allow_user_certificates":     true,
		"allowed_users":               "tuber",
		"default_user":                "tuber",
		"default_extensions_template": true,
		"default_extensions": map[string]interface{}{
			"login@foobar.com": "{{identity.entity.aliases." + userpassAccessor + ".name}}",
			"login@foobar2.com": "{{identity.entity.aliases." + userpassAccessor + ".name}}, " +
				"{{identity.entity.aliases." + userpassAccessor + ".name}}_foobar",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	sshKeyID := "vault-userpass-" + testUserName + "-9bd0f01b7dfc50a13aa5e5cd11aea19276968755c8f1f9c98965d04147f30ed0"

	// Issue SSH certificate with default extensions templating enabled, and no user-provided extensions
	client.SetToken(userpassToken)
	resp, err := client.Logical().Write("ssh/sign/test", map[string]interface{}{
		"public_key": publicKey4096,
	})
	if err != nil {
		t.Fatal(err)
	}
	signedKey := resp.Data["signed_key"].(string)
	key, _ := base64.StdEncoding.DecodeString(strings.Split(signedKey, " ")[1])

	parsedKey, err := ssh.ParsePublicKey(key)
	if err != nil {
		t.Fatal(err)
	}

	defaultExtensionPermissions := map[string]string{
		"login@foobar.com":  testUserName,
		"login@foobar2.com": fmt.Sprintf("%s, %s_foobar", testUserName, testUserName),
	}

	err = validateSSHCertificate(parsedKey.(*ssh.Certificate), sshKeyID, ssh.UserCert, []string{"tuber"}, map[string]string{}, defaultExtensionPermissions, 16*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// Issue SSH certificate with default extensions templating enabled, and user-provided extensions
	// The certificate should only have the user-provided extensions, and no templated extensions
	userProvidedExtensionPermissions := map[string]string{
		"login@zipzap.com": "some_other_user_name",
	}
	resp, err = client.Logical().Write("ssh/sign/test", map[string]interface{}{
		"public_key": publicKey4096,
		"extensions": userProvidedExtensionPermissions,
	})
	if err != nil {
		t.Fatal(err)
	}
	signedKey = resp.Data["signed_key"].(string)
	key, _ = base64.StdEncoding.DecodeString(strings.Split(signedKey, " ")[1])

	parsedKey, err = ssh.ParsePublicKey(key)
	if err != nil {
		t.Fatal(err)
	}

	err = validateSSHCertificate(parsedKey.(*ssh.Certificate), sshKeyID, ssh.UserCert, []string{"tuber"}, map[string]string{}, userProvidedExtensionPermissions, 16*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// Issue SSH certificate with default extensions templating enabled, and invalid user-provided extensions - it should fail
	invalidUserProvidedExtensionPermissions := map[string]string{
		"login@foobar.com": "{{identity.entity.metadata}}",
	}
	resp, err = client.Logical().Write("ssh/sign/test", map[string]interface{}{
		"public_key": publicKey4096,
		"extensions": invalidUserProvidedExtensionPermissions,
	})
	if err == nil {
		t.Fatal("expected an error while attempting to sign a key with invalid permissions")
	}
}

func TestBackend_EmptyAllowedExtensionFailsClosed(t *testing.T) {
	cluster, userpassToken := getSshCaTestCluster(t, testUserName)
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// Get auth accessor for identity template.
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	userpassAccessor := auths["userpass/"].Accessor

	// Write SSH role to test with no allowed extension. We also provide a templated default extension,
	// to verify that it's not actually being evaluated
	_, err = client.Logical().Write("ssh/roles/test_allow_all_extensions", map[string]interface{}{
		"key_type":                    "ca",
		"allow_user_certificates":     true,
		"allowed_users":               "tuber",
		"default_user":                "tuber",
		"allowed_extensions":          "",
		"default_extensions_template": false,
		"default_extensions": map[string]interface{}{
			"login@foobar.com": "{{identity.entity.aliases." + userpassAccessor + ".name}}",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Issue SSH certificate with default extensions templating disabled, and user-provided extensions
	client.SetToken(userpassToken)
	userProvidedAnyExtensionPermissions := map[string]string{
		"login@foobar.com": "not_userpassname",
	}
	_, err = client.Logical().Write("ssh/sign/test_allow_all_extensions", map[string]interface{}{
		"public_key": publicKey4096,
		"extensions": userProvidedAnyExtensionPermissions,
	})
	if err == nil {
		t.Fatal("Expected failure we should not have allowed specifying custom extensions")
	}

	if !strings.Contains(err.Error(), "are not on allowed list") {
		t.Fatalf("Expected failure to contain 'are not on allowed list' but was %s", err)
	}
}

func TestBackend_DefExtTemplatingDisabled(t *testing.T) {
	cluster, userpassToken := getSshCaTestCluster(t, testUserName)
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// Get auth accessor for identity template.
	auths, err := client.Sys().ListAuth()
	if err != nil {
		t.Fatal(err)
	}
	userpassAccessor := auths["userpass/"].Accessor

	// Write SSH role to test with any extension. We also provide a templated default extension,
	// to verify that it's not actually being evaluated
	_, err = client.Logical().Write("ssh/roles/test_allow_all_extensions", map[string]interface{}{
		"key_type":                    "ca",
		"allow_user_certificates":     true,
		"allowed_users":               "tuber",
		"default_user":                "tuber",
		"allowed_extensions":          "*",
		"default_extensions_template": false,
		"default_extensions": map[string]interface{}{
			"login@foobar.com": "{{identity.entity.aliases." + userpassAccessor + ".name}}",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	sshKeyID := "vault-userpass-" + testUserName + "-9bd0f01b7dfc50a13aa5e5cd11aea19276968755c8f1f9c98965d04147f30ed0"

	// Issue SSH certificate with default extensions templating disabled, and no user-provided extensions
	client.SetToken(userpassToken)
	defaultExtensionPermissions := map[string]string{
		"login@foobar.com": "{{identity.entity.aliases." + userpassAccessor + ".name}}",
		"login@zipzap.com": "some_other_user_name",
	}
	resp, err := client.Logical().Write("ssh/sign/test_allow_all_extensions", map[string]interface{}{
		"public_key": publicKey4096,
		"extensions": defaultExtensionPermissions,
	})
	if err != nil {
		t.Fatal(err)
	}
	signedKey := resp.Data["signed_key"].(string)
	key, _ := base64.StdEncoding.DecodeString(strings.Split(signedKey, " ")[1])

	parsedKey, err := ssh.ParsePublicKey(key)
	if err != nil {
		t.Fatal(err)
	}

	err = validateSSHCertificate(parsedKey.(*ssh.Certificate), sshKeyID, ssh.UserCert, []string{"tuber"}, map[string]string{}, defaultExtensionPermissions, 16*time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// Issue SSH certificate with default extensions templating disabled, and user-provided extensions
	client.SetToken(userpassToken)
	userProvidedAnyExtensionPermissions := map[string]string{
		"login@foobar.com": "not_userpassname",
		"login@zipzap.com": "some_other_user_name",
	}
	resp, err = client.Logical().Write("ssh/sign/test_allow_all_extensions", map[string]interface{}{
		"public_key": publicKey4096,
		"extensions": userProvidedAnyExtensionPermissions,
	})
	if err != nil {
		t.Fatal(err)
	}
	signedKey = resp.Data["signed_key"].(string)
	key, _ = base64.StdEncoding.DecodeString(strings.Split(signedKey, " ")[1])

	parsedKey, err = ssh.ParsePublicKey(key)
	if err != nil {
		t.Fatal(err)
	}

	err = validateSSHCertificate(parsedKey.(*ssh.Certificate), sshKeyID, ssh.UserCert, []string{"tuber"}, map[string]string{}, userProvidedAnyExtensionPermissions, 16*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSSHBackend_ValidateNotBeforeDuration(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}
	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(testCAPublicKey, testCAPrivateKey),

			createRoleStep("testing", map[string]interface{}{
				"key_type":                "ca",
				"allow_host_certificates": true,
				"allowed_domains":         "example.com,example.org",
				"allow_subdomains":        true,
				"default_critical_options": map[string]interface{}{
					"option": "value",
				},
				"default_extensions": map[string]interface{}{
					"extension": "extended",
				},
				"not_before_duration": "300s",
			}),

			signCertificateStep("testing", "vault-root-22608f5ef173aabf700797cb95c5641e792698ec6380e8e1eb55523e39aa5e51", ssh.HostCert, []string{"dummy.example.org", "second.example.com"}, map[string]string{
				"option": "value",
			}, map[string]string{
				"extension": "extended",
			},
				2*time.Hour+5*time.Minute-30*time.Second, map[string]interface{}{
					"public_key":       publicKey2,
					"ttl":              "2h",
					"cert_type":        "host",
					"valid_principals": "dummy.example.org,second.example.com",
				}),

			createRoleStep("testing", map[string]interface{}{
				"key_type":                "ca",
				"allow_host_certificates": true,
				"allowed_domains":         "example.com,example.org",
				"allow_subdomains":        true,
				"default_critical_options": map[string]interface{}{
					"option": "value",
				},
				"default_extensions": map[string]interface{}{
					"extension": "extended",
				},
				"not_before_duration": "2h",
			}),

			signCertificateStep("testing", "vault-root-22608f5ef173aabf700797cb95c5641e792698ec6380e8e1eb55523e39aa5e51", ssh.HostCert, []string{"dummy.example.org", "second.example.com"}, map[string]string{
				"option": "value",
			}, map[string]string{
				"extension": "extended",
			},
				4*time.Hour-30*time.Second, map[string]interface{}{
					"public_key":       publicKey2,
					"ttl":              "2h",
					"cert_type":        "host",
					"valid_principals": "dummy.example.org,second.example.com",
				}),
			createRoleStep("testing", map[string]interface{}{
				"key_type":                "ca",
				"allow_host_certificates": true,
				"allowed_domains":         "example.com,example.org",
				"allow_subdomains":        true,
				"default_critical_options": map[string]interface{}{
					"option": "value",
				},
				"default_extensions": map[string]interface{}{
					"extension": "extended",
				},
				"not_before_duration": "30s",
			}),

			signCertificateStep("testing", "vault-root-22608f5ef173aabf700797cb95c5641e792698ec6380e8e1eb55523e39aa5e51", ssh.HostCert, []string{"dummy.example.org", "second.example.com"}, map[string]string{
				"option": "value",
			}, map[string]string{
				"extension": "extended",
			},
				2*time.Hour, map[string]interface{}{
					"public_key":       publicKey2,
					"ttl":              "2h",
					"cert_type":        "host",
					"valid_principals": "dummy.example.org,second.example.com",
				}),
		},
	}

	logicaltest.Test(t, testCase)
}

func TestSSHBackend_IssueSign(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(testCAPublicKey, testCAPrivateKey),

			createRoleStep("testing", map[string]interface{}{
				"key_type":     "otp",
				"default_user": "user",
			}),
			// Key pair not issued with invalid role key type
			issueSSHKeyPairStep("testing", "rsa", 0, true, "role key type 'otp' not allowed to issue key pairs"),

			createRoleStep("testing", map[string]interface{}{
				"key_type":                "ca",
				"allow_user_key_ids":      false,
				"allow_user_certificates": true,
				"allowed_user_key_lengths": map[string]interface{}{
					"ssh-rsa":             []int{2048, 3072, 4096},
					"ecdsa-sha2-nistp521": 0,
					"ed25519":             0,
				},
				"allow_empty_principals": true,
			}),
			// Key_type not in allowed_user_key_types_lengths
			issueSSHKeyPairStep("testing", "ec", 256, true, "provided key_type value not in allowed_user_key_types"),
			// Key_bits not in allowed_user_key_types_lengths for provided key_type
			issueSSHKeyPairStep("testing", "rsa", 2560, true, "provided key_bits value not in list of role's allowed_user_key_types"),
			// key_type `rsa` and key_bits `2048` successfully created
			issueSSHKeyPairStep("testing", "rsa", 2048, false, ""),
			// key_type `ed22519` and key_bits `0` successfully created
			issueSSHKeyPairStep("testing", "ed25519", 0, false, ""),
		},
	}

	logicaltest.Test(t, testCase)
}

func getSshCaTestCluster(t *testing.T, userIdentity string) (*vault.TestCluster, string) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"ssh": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	client := cluster.Cores[0].Client

	// Write test policy for userpass auth method.
	err := client.Sys().PutPolicy("test", `
   path "ssh/*" {
     capabilities = ["update"]
   }`)
	if err != nil {
		t.Fatal(err)
	}

	// Enable userpass auth method.
	if err := client.Sys().EnableAuth("userpass", "userpass", ""); err != nil {
		t.Fatal(err)
	}

	// Configure test role for userpass.
	if _, err := client.Logical().Write("auth/userpass/users/"+userIdentity, map[string]interface{}{
		"password": "test",
		"policies": "test",
	}); err != nil {
		t.Fatal(err)
	}

	// Login userpass for test role and keep client token.
	secret, err := client.Logical().Write("auth/userpass/login/"+userIdentity, map[string]interface{}{
		"password": "test",
	})
	if err != nil || secret == nil {
		t.Fatal(err)
	}
	userpassToken := secret.Auth.ClientToken

	// Mount SSH.
	err = client.Sys().Mount("ssh", &api.MountInput{
		Type: "ssh",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "60h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Configure SSH CA.
	_, err = client.Logical().Write("ssh/config/ca", map[string]interface{}{
		"public_key":  testCAPublicKey,
		"private_key": testCAPrivateKey,
	})
	if err != nil {
		t.Fatal(err)
	}

	return cluster, userpassToken
}

func testDefaultUserTemplate(t *testing.T, testDefaultUserTemplate string,
	expectedValidPrincipal string, testEntityMetadata map[string]string,
) {
	cluster, userpassToken := getSshCaTestCluster(t, testUserName)
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// set metadata "ssh_username" to userpass username
	tokenLookupResponse, err := client.Logical().Write("/auth/token/lookup", map[string]interface{}{
		"token": userpassToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	entityID := tokenLookupResponse.Data["entity_id"].(string)
	_, err = client.Logical().Write("/identity/entity/id/"+entityID, map[string]interface{}{
		"metadata": testEntityMetadata,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("ssh/roles/my-role", map[string]interface{}{
		"key_type":                testCaKeyType,
		"allow_user_certificates": true,
		"default_user":            testDefaultUserTemplate,
		"default_user_template":   true,
		"allowed_users":           testDefaultUserTemplate,
		"allowed_users_template":  true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// sign SSH key as userpass user
	client.SetToken(userpassToken)
	signResponse, err := client.Logical().Write("ssh/sign/my-role", map[string]interface{}{
		"public_key": testCAPublicKey,
	})
	if err != nil {
		t.Fatal(err)
	}

	// check for the expected valid principals of certificate
	signedKey := signResponse.Data["signed_key"].(string)
	key, _ := base64.StdEncoding.DecodeString(strings.Split(signedKey, " ")[1])
	parsedKey, err := ssh.ParsePublicKey(key)
	if err != nil {
		t.Fatal(err)
	}
	actualPrincipals := parsedKey.(*ssh.Certificate).ValidPrincipals
	if actualPrincipals[0] != expectedValidPrincipal {
		t.Fatal(
			fmt.Sprintf("incorrect ValidPrincipals: %v should be %v",
				actualPrincipals, []string{expectedValidPrincipal}),
		)
	}
}

func testAllowedPrincipalsTemplate(t *testing.T, testAllowedDomainsTemplate string,
	expectedValidPrincipal string, testEntityMetadata map[string]string,
	roleConfigPayload map[string]interface{}, signingPayload map[string]interface{},
) {
	cluster, userpassToken := getSshCaTestCluster(t, testUserName)
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	// set metadata "ssh_username" to userpass username
	tokenLookupResponse, err := client.Logical().Write("/auth/token/lookup", map[string]interface{}{
		"token": userpassToken,
	})
	if err != nil {
		t.Fatal(err)
	}
	entityID := tokenLookupResponse.Data["entity_id"].(string)
	_, err = client.Logical().Write("/identity/entity/id/"+entityID, map[string]interface{}{
		"metadata": testEntityMetadata,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("ssh/roles/my-role", roleConfigPayload)
	if err != nil {
		t.Fatal(err)
	}

	// sign SSH key as userpass user
	client.SetToken(userpassToken)
	signResponse, err := client.Logical().Write("ssh/sign/my-role", signingPayload)
	if err != nil {
		t.Fatal(err)
	}

	// check for the expected valid principals of certificate
	signedKey := signResponse.Data["signed_key"].(string)
	key, _ := base64.StdEncoding.DecodeString(strings.Split(signedKey, " ")[1])
	parsedKey, err := ssh.ParsePublicKey(key)
	if err != nil {
		t.Fatal(err)
	}
	actualPrincipals := parsedKey.(*ssh.Certificate).ValidPrincipals
	if actualPrincipals[0] != expectedValidPrincipal {
		t.Fatal(
			fmt.Sprintf("incorrect ValidPrincipals: %v should be %v",
				actualPrincipals, []string{expectedValidPrincipal}),
		)
	}
}

func testAllowedUsersTemplate(t *testing.T, testAllowedUsersTemplate string,
	expectedValidPrincipal string, testEntityMetadata map[string]string,
) {
	testAllowedPrincipalsTemplate(
		t, testAllowedUsersTemplate,
		expectedValidPrincipal, testEntityMetadata,
		map[string]interface{}{
			"key_type":                testCaKeyType,
			"allow_user_certificates": true,
			"allowed_users":           testAllowedUsersTemplate,
			"allowed_users_template":  true,
		},
		map[string]interface{}{
			"public_key":       testCAPublicKey,
			"valid_principals": expectedValidPrincipal,
		},
	)
}

func configCaStep(caPublicKey, caPrivateKey string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/ca",
		Data: map[string]interface{}{
			"public_key":  caPublicKey,
			"private_key": caPrivateKey,
		},
	}
}

func createRoleStep(name string, parameters map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.CreateOperation,
		Path:      "roles/" + name,
		Data:      parameters,
	}
}

func signCertificateStep(
	role, keyID string, certType int, validPrincipals []string,
	criticalOptionPermissions, extensionPermissions map[string]string,
	ttl time.Duration,
	requestParameters map[string]interface{},
) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "sign/" + role,
		Data:      requestParameters,

		Check: func(resp *logical.Response) error {
			serialNumber := resp.Data["serial_number"].(string)
			if serialNumber == "" {
				return errors.New("no serial number in response")
			}

			signedKey := strings.TrimSpace(resp.Data["signed_key"].(string))
			if signedKey == "" {
				return errors.New("no signed key in response")
			}

			key, _ := base64.StdEncoding.DecodeString(strings.Split(signedKey, " ")[1])

			parsedKey, err := ssh.ParsePublicKey(key)
			if err != nil {
				return err
			}

			return validateSSHCertificate(parsedKey.(*ssh.Certificate), keyID, certType, validPrincipals, criticalOptionPermissions, extensionPermissions, ttl)
		},
	}
}

func issueSSHKeyPairStep(role, keyType string, keyBits int, expectError bool, errorMsg string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "issue/" + role,
		Data: map[string]interface{}{
			"key_type": keyType,
			"key_bits": keyBits,
		},
		ErrorOk: true,
		Check: func(resp *logical.Response) error {
			if expectError {
				var err error
				if resp.Data["error"] != errorMsg {
					err = fmt.Errorf("actual error message \"%s\" different from expected error message \"%s\"", resp.Data["error"], errorMsg)
				}

				return err
			}

			if resp.IsError() {
				return fmt.Errorf("unexpected error response returned: %v", resp.Error())
			}

			if resp.Data["private_key_type"] != keyType {
				return fmt.Errorf("response private_key_type (%s) does not match the provided key_type (%s)", resp.Data["private_key_type"], keyType)
			}

			if resp.Data["signed_key"] == "" {
				return errors.New("certificate/signed_key should not be empty")
			}

			return nil
		},
	}
}

func validateSSHCertificate(cert *ssh.Certificate, keyID string, certType int, validPrincipals []string, criticalOptionPermissions, extensionPermissions map[string]string,
	ttl time.Duration,
) error {
	if cert.KeyId != keyID {
		return fmt.Errorf("incorrect KeyId: %v, wanted %v", cert.KeyId, keyID)
	}

	if cert.CertType != uint32(certType) {
		return fmt.Errorf("incorrect CertType: %v", cert.CertType)
	}

	if time.Unix(int64(cert.ValidAfter), 0).After(time.Now()) {
		return fmt.Errorf("incorrect ValidAfter: %v", cert.ValidAfter)
	}

	if time.Unix(int64(cert.ValidBefore), 0).Before(time.Now()) {
		return fmt.Errorf("incorrect ValidBefore: %v", cert.ValidBefore)
	}

	actualTTL := time.Unix(int64(cert.ValidBefore), 0).Add(-30 * time.Second).Sub(time.Unix(int64(cert.ValidAfter), 0))
	if actualTTL != ttl {
		return fmt.Errorf("incorrect ttl: expected: %v, actual %v", ttl, actualTTL)
	}

	if !reflect.DeepEqual(cert.ValidPrincipals, validPrincipals) {
		return fmt.Errorf("incorrect ValidPrincipals: expected: %#v actual: %#v", validPrincipals, cert.ValidPrincipals)
	}

	publicSigningKey, err := getSigningPublicKey()
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(cert.SignatureKey, publicSigningKey) {
		return fmt.Errorf("incorrect SignatureKey: %v", cert.SignatureKey)
	}

	if cert.Signature == nil {
		return fmt.Errorf("incorrect Signature: %v", cert.Signature)
	}

	if !reflect.DeepEqual(cert.Permissions.Extensions, extensionPermissions) {
		return fmt.Errorf("incorrect Permissions.Extensions: Expected: %v, Actual: %v", extensionPermissions, cert.Permissions.Extensions)
	}

	if !reflect.DeepEqual(cert.Permissions.CriticalOptions, criticalOptionPermissions) {
		return fmt.Errorf("incorrect Permissions.CriticalOptions: %v", cert.Permissions.CriticalOptions)
	}

	return nil
}

func getSigningPublicKey() (ssh.PublicKey, error) {
	key, err := base64.StdEncoding.DecodeString(strings.Split(testCAPublicKey, " ")[1])
	if err != nil {
		return nil, err
	}

	parsedKey, err := ssh.ParsePublicKey(key)
	if err != nil {
		return nil, err
	}

	return parsedKey, nil
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
				return fmt.Errorf("invalid response")
			}
			return nil
		},
	}
}

func testLookupRead(t *testing.T, data map[string]interface{}, expected []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "lookup",
		Data:      data,
		Check: func(resp *logical.Response) error {
			if resp.Data == nil || resp.Data["roles"] == nil {
				return fmt.Errorf("missing roles information")
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

func testRoleList(t *testing.T, expected map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "roles",
		Check: func(resp *logical.Response) error {
			if resp == nil {
				return fmt.Errorf("nil response")
			}
			if resp.Data == nil {
				return fmt.Errorf("nil data")
			}
			if !reflect.DeepEqual(resp.Data, expected) {
				return fmt.Errorf("Invalid response:\nactual:%#v\nexpected is %#v", resp.Data, expected)
			}
			return nil
		},
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
			switch d.KeyType {
			case "otp":
				if d.KeyType != expected["key_type"] || d.DefaultUser != expected["default_user"] || d.CIDRList != expected["cidr_list"] {
					return fmt.Errorf("data mismatch. bad: %#v", resp)
				}
			default:
				return fmt.Errorf("unknown key type. bad: %#v", resp)
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

func testCredsWrite(t *testing.T, roleName string, data map[string]interface{}, expectError bool, address string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("creds/%s", roleName),
		Data:      data,
		ErrorOk:   expectError,
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
					return fmt.Errorf("expected error, but write succeeded")
				}
				return nil
			}
			if roleName == testAtRoleName {
				var d struct {
					Key string `mapstructure:"key"`
				}
				if err := mapstructure.Decode(resp.Data, &d); err != nil {
					return err
				}
				if d.Key == "" {
					return fmt.Errorf("generated key is an empty string")
				}
				// Checking only for a parsable key
				privKey, err := ssh.ParsePrivateKey([]byte(d.Key))
				if err != nil {
					return fmt.Errorf("generated key is invalid")
				}
				if err := testSSH(data["username"].(string), address, ssh.PublicKeys(privKey), "date"); err != nil {
					return fmt.Errorf("unable to SSH with new key (%s): %w", d.Key, err)
				}
			} else {
				if resp.Data["key_type"] != KeyTypeOTP {
					return fmt.Errorf("incorrect key_type")
				}
				if resp.Data["key"] == nil {
					return fmt.Errorf("invalid key")
				}
			}
			return nil
		},
	}
}

func TestBackend_CleanupDynamicHostKeys(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// Running on a clean mount shouldn't do anything.
	cleanRequest := &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "tidy/dynamic-keys",
		Storage:   config.StorageView,
	}

	resp, err := b.HandleRequest(context.Background(), cleanRequest)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotNil(t, resp.Data["message"])
	require.Contains(t, resp.Data["message"], "0 of 0")

	// Write a bunch of bogus entries.
	for i := 0; i < 15; i++ {
		data := map[string]interface{}{
			"host": "localhost",
			"key":  "nothing-to-see-here",
		}
		entry, err := logical.StorageEntryJSON(fmt.Sprintf("%vexample-%v", keysStoragePrefix, i), &data)
		require.NoError(t, err)
		err = config.StorageView.Put(context.Background(), entry)
		require.NoError(t, err)
	}

	// Should now have 15
	resp, err = b.HandleRequest(context.Background(), cleanRequest)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotNil(t, resp.Data["message"])
	require.Contains(t, resp.Data["message"], "15 of 15")

	// Should have none left.
	resp, err = b.HandleRequest(context.Background(), cleanRequest)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotNil(t, resp.Data["message"])
	require.Contains(t, resp.Data["message"], "0 of 0")
}

type pathAuthCheckerFunc func(t *testing.T, client *api.Client, path string, token string)

func isPermDenied(err error) bool {
	return strings.Contains(err.Error(), "permission denied")
}

func isUnsupportedPathOperation(err error) bool {
	return strings.Contains(err.Error(), "unsupported path") || strings.Contains(err.Error(), "unsupported operation")
}

func isDeniedOp(err error) bool {
	return isPermDenied(err) || isUnsupportedPathOperation(err)
}

func pathShouldBeAuthed(t *testing.T, client *api.Client, path string, token string) {
	client.SetToken("")
	resp, err := client.Logical().ReadWithContext(ctx, path)
	if err == nil || !isPermDenied(err) {
		t.Fatalf("expected failure to read %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().ListWithContext(ctx, path)
	if err == nil || !isPermDenied(err) {
		t.Fatalf("expected failure to list %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().WriteWithContext(ctx, path, map[string]interface{}{})
	if err == nil || !isPermDenied(err) {
		t.Fatalf("expected failure to write %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().DeleteWithContext(ctx, path)
	if err == nil || !isPermDenied(err) {
		t.Fatalf("expected failure to delete %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().JSONMergePatch(ctx, path, map[string]interface{}{})
	if err == nil || !isPermDenied(err) {
		t.Fatalf("expected failure to patch %v while unauthed: %v / %v", path, err, resp)
	}
}

func pathShouldBeUnauthedReadList(t *testing.T, client *api.Client, path string, token string) {
	// Should be able to read both with and without a token.
	client.SetToken("")
	resp, err := client.Logical().ReadWithContext(ctx, path)
	if err != nil && isPermDenied(err) {
		// Read will sometimes return permission denied, when the handler
		// does not support the given operation. Retry with the token.
		client.SetToken(token)
		resp2, err2 := client.Logical().ReadWithContext(ctx, path)
		if err2 != nil && !isUnsupportedPathOperation(err2) {
			t.Fatalf("unexpected failure to read %v while unauthed: %v / %v\nWhile authed: %v / %v", path, err, resp, err2, resp2)
		}
		client.SetToken("")
	}
	resp, err = client.Logical().ListWithContext(ctx, path)
	if err != nil && isPermDenied(err) {
		// List will sometimes return permission denied, when the handler
		// does not support the given operation. Retry with the token.
		client.SetToken(token)
		resp2, err2 := client.Logical().ListWithContext(ctx, path)
		if err2 != nil && !isUnsupportedPathOperation(err2) {
			t.Fatalf("unexpected failure to list %v while unauthed: %v / %v\nWhile authed: %v / %v", path, err, resp, err2, resp2)
		}
		client.SetToken("")
	}

	// These should all be denied.
	resp, err = client.Logical().WriteWithContext(ctx, path, map[string]interface{}{})
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during write on read-only path %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().DeleteWithContext(ctx, path)
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during delete on read-only path %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().JSONMergePatch(ctx, path, map[string]interface{}{})
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during patch on read-only path %v while unauthed: %v / %v", path, err, resp)
	}

	// Retrying with token should allow read/list, but not modification still.
	client.SetToken(token)
	resp, err = client.Logical().ReadWithContext(ctx, path)
	if err != nil && isPermDenied(err) {
		t.Fatalf("unexpected failure to read %v while authed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().ListWithContext(ctx, path)
	if err != nil && isPermDenied(err) {
		t.Fatalf("unexpected failure to list %v while authed: %v / %v", path, err, resp)
	}

	// Should all be denied.
	resp, err = client.Logical().WriteWithContext(ctx, path, map[string]interface{}{})
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during write on read-only path %v while authed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().DeleteWithContext(ctx, path)
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during delete on read-only path %v while authed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().JSONMergePatch(ctx, path, map[string]interface{}{})
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during patch on read-only path %v while authed: %v / %v", path, err, resp)
	}
}

func pathShouldBeUnauthedWriteOnly(t *testing.T, client *api.Client, path string, token string) {
	client.SetToken("")
	resp, err := client.Logical().WriteWithContext(ctx, path, map[string]interface{}{})
	if err != nil && isPermDenied(err) {
		t.Fatalf("unexpected failure to write %v while unauthed: %v / %v", path, err, resp)
	}

	// These should all be denied.
	resp, err = client.Logical().ReadWithContext(ctx, path)
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during read on write-only path %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().ListWithContext(ctx, path)
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during list on write-only path %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().DeleteWithContext(ctx, path)
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during delete on write-only path %v while unauthed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().JSONMergePatch(ctx, path, map[string]interface{}{})
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during patch on write-only path %v while unauthed: %v / %v", path, err, resp)
	}

	// Retrying with token should allow writing, but nothing else.
	client.SetToken(token)
	resp, err = client.Logical().WriteWithContext(ctx, path, map[string]interface{}{})
	if err != nil && isPermDenied(err) {
		t.Fatalf("unexpected failure to write %v while unauthed: %v / %v", path, err, resp)
	}

	// These should all be denied.
	resp, err = client.Logical().ReadWithContext(ctx, path)
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during read on write-only path %v while authed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().ListWithContext(ctx, path)
	if err == nil || !isDeniedOp(err) {
		if resp != nil || err != nil {
			t.Fatalf("unexpected failure during list on write-only path %v while authed: %v / %v", path, err, resp)
		}
	}
	resp, err = client.Logical().DeleteWithContext(ctx, path)
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during delete on write-only path %v while authed: %v / %v", path, err, resp)
	}
	resp, err = client.Logical().JSONMergePatch(ctx, path, map[string]interface{}{})
	if err == nil || !isDeniedOp(err) {
		t.Fatalf("unexpected failure during patch on write-only path %v while authed: %v / %v", path, err, resp)
	}
}

type pathAuthChecker int

const (
	shouldBeAuthed pathAuthChecker = iota
	shouldBeUnauthedReadList
	shouldBeUnauthedWriteOnly
)

var pathAuthChckerMap = map[pathAuthChecker]pathAuthCheckerFunc{
	shouldBeAuthed:            pathShouldBeAuthed,
	shouldBeUnauthedReadList:  pathShouldBeUnauthedReadList,
	shouldBeUnauthedWriteOnly: pathShouldBeUnauthedWriteOnly,
}

func TestProperAuthing(t *testing.T) {
	t.Parallel()
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"ssh": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client
	token := client.Token()

	// Mount SSH.
	err := client.Sys().MountWithContext(ctx, "ssh", &api.MountInput{
		Type: "ssh",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "60h",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Setup basic configuration.
	_, err = client.Logical().WriteWithContext(ctx, "ssh/config/ca", map[string]interface{}{
		"generate_signing_key": true,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().WriteWithContext(ctx, "ssh/roles/test-ca", map[string]interface{}{
		"key_type":                "ca",
		"allow_user_certificates": true,
		"allowed_users":           "toor",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().WriteWithContext(ctx, "ssh/issue/test-ca", map[string]interface{}{
		"valid_principals": "toor",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().WriteWithContext(ctx, "ssh/roles/test-otp", map[string]interface{}{
		"key_type":     "otp",
		"default_user": "toor",
		"cidr_list":    "127.0.0.0/24",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().WriteWithContext(ctx, "ssh/creds/test-otp", map[string]interface{}{
		"username": "toor",
		"ip":       "127.0.0.1",
	})
	if err != nil || resp == nil {
		t.Fatal(err)
	}
	// key := resp.Data["key"].(string)

	paths := map[string]pathAuthChecker{
		"config/ca":          shouldBeAuthed,
		"config/zeroaddress": shouldBeAuthed,
		"creds/test-otp":     shouldBeAuthed,
		"issue/test-ca":      shouldBeAuthed,
		"lookup":             shouldBeAuthed,
		"public_key":         shouldBeUnauthedReadList,
		"roles/test-ca":      shouldBeAuthed,
		"roles/test-otp":     shouldBeAuthed,
		"roles/":             shouldBeAuthed,
		"sign/test-ca":       shouldBeAuthed,
		"tidy/dynamic-keys":  shouldBeAuthed,
		"verify":             shouldBeUnauthedWriteOnly,
	}
	for path, checkerType := range paths {
		checker := pathAuthChckerMap[checkerType]
		checker(t, client, "ssh/"+path, token)
	}

	client.SetToken(token)
	openAPIResp, err := client.Logical().ReadWithContext(ctx, "sys/internal/specs/openapi")
	if err != nil {
		t.Fatalf("failed to get openapi data: %v", err)
	}

	if len(openAPIResp.Data["paths"].(map[string]interface{})) == 0 {
		t.Fatalf("expected to get response from OpenAPI; got empty path list")
	}

	validatedPath := false
	for openapi_path, raw_data := range openAPIResp.Data["paths"].(map[string]interface{}) {
		if !strings.HasPrefix(openapi_path, "/ssh/") {
			t.Logf("Skipping path: %v", openapi_path)
			continue
		}

		t.Logf("Validating path: %v", openapi_path)
		validatedPath = true

		// Substitute values in from our testing map.
		raw_path := openapi_path[5:]
		if strings.Contains(raw_path, "{role}") && strings.Contains(raw_path, "roles/") {
			raw_path = strings.ReplaceAll(raw_path, "{role}", "test-ca")
		}
		if strings.Contains(raw_path, "{role}") && (strings.Contains(raw_path, "sign/") || strings.Contains(raw_path, "issue/")) {
			raw_path = strings.ReplaceAll(raw_path, "{role}", "test-ca")
		}
		if strings.Contains(raw_path, "{role}") && strings.Contains(raw_path, "creds") {
			raw_path = strings.ReplaceAll(raw_path, "{role}", "test-otp")
		}

		handler, present := paths[raw_path]
		if !present {
			t.Fatalf("OpenAPI reports SSH mount contains %v -> %v but was not tested to be authed or not authed.",
				openapi_path, raw_path)
		}

		openapi_data := raw_data.(map[string]interface{})
		hasList := false
		rawGetData, hasGet := openapi_data["get"]
		if hasGet {
			getData := rawGetData.(map[string]interface{})
			getParams, paramsPresent := getData["parameters"].(map[string]interface{})
			if getParams != nil && paramsPresent {
				if _, hasList = getParams["list"]; hasList {
					// LIST is exclusive from GET on the same endpoint usually.
					hasGet = false
				}
			}
		}
		_, hasPost := openapi_data["post"]
		_, hasDelete := openapi_data["delete"]

		if handler == shouldBeUnauthedReadList {
			if hasPost || hasDelete {
				t.Fatalf("Unauthed read-only endpoints should not have POST/DELETE capabilities")
			}
		}
	}

	if !validatedPath {
		t.Fatalf("Expected to have validated at least one path.")
	}
}
