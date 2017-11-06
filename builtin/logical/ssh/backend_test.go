package ssh

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"encoding/base64"
	"errors"
	"strings"

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
	// Public half of `privateKey`, identical to how it would be fed in from a file
	publicKey = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDArgK0ilRRfk8E7HIsjz5l3BuxmwpDd8DHRCVfOhbZ4gOSVxjEOOqBwWGjygdboBIZwFXmwDlU6sWX0hBJAgpQz0Cjvbjxtq/NjkvATrYPgnrXUhTaEn2eQO0PsqRNSFH46SK/oJfTp0q8/WgojxWJ2L7FUV8PO8uIk49DzqAqPV7WXU63vFsjx+3WQOX/ILeQvHCvaqs3dWjjzEoDudRWCOdUqcHEOshV9azIzPrXlQVzRV3QAKl6u7pC+/Secorpwt6IHpMKoVPGiR0tMMuNOVH8zrAKzIxPGfy2WmNDpJopbXMTvSOGAqNcp49O4SKOQl9Fzfq2HEevJamKLrMB dummy@example.com
`
	publicKey2 = `AAAAB3NzaC1yc2EAAAADAQABAAABAQDArgK0ilRRfk8E7HIsjz5l3BuxmwpDd8DHRCVfOhbZ4gOSVxjEOOqBwWGjygdboBIZwFXmwDlU6sWX0hBJAgpQz0Cjvbjxtq/NjkvATrYPgnrXUhTaEn2eQO0PsqRNSFH46SK/oJfTp0q8/WgojxWJ2L7FUV8PO8uIk49DzqAqPV7WXU63vFsjx+3WQOX/ILeQvHCvaqs3dWjjzEoDudRWCOdUqcHEOshV9azIzPrXlQVzRV3QAKl6u7pC+/Secorpwt6IHpMKoVPGiR0tMMuNOVH8zrAKzIxPGfy2WmNDpJopbXMTvSOGAqNcp49O4SKOQl9Fzfq2HEevJamKLrMB
`
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
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
)

func TestBackend_allowed_users(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}
	err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}
	err = b.Initialize()
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
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: testingFactory,
		Steps: []logicaltest.TestStep{
			testRoleList(t, resp1),
			testRoleWrite(t, testOTPRoleName, testOTPRoleData),
			testRoleList(t, resp2),
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

func TestBackend_AbleToRetrievePublicKey(t *testing.T) {

	config := logical.TestBackendConfig()

	b, err := Factory(config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(),

			logicaltest.TestStep{
				Operation:       logical.ReadOperation,
				Path:            "public_key",
				Unauthenticated: true,

				Check: func(resp *logical.Response) error {

					key := string(resp.Data["http_raw_body"].([]byte))

					if key != publicKey {
						return fmt.Errorf("public_key incorrect. Expected %v, actual %v", publicKey, key)
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

	b, err := Factory(config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "config/ca",
			},

			logicaltest.TestStep{
				Operation:       logical.ReadOperation,
				Path:            "public_key",
				Unauthenticated: true,

				Check: func(resp *logical.Response) error {

					key := string(resp.Data["http_raw_body"].([]byte))

					if key == "" {
						return fmt.Errorf("public_key empty. Expected not empty, actual %s", key)
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

	b, err := Factory(config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(),

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

	b, err := Factory(config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(),

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

func TestBackend_CustomKeyIDFormat(t *testing.T) {
	config := logical.TestBackendConfig()

	b, err := Factory(config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(),

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

	b, err := Factory(config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	testCase := logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			configCaStep(),

			createRoleStep("testing", map[string]interface{}{
				"key_type":                "ca",
				"allow_user_key_ids":      false,
				"allow_user_certificates": true,
			}),
			logicaltest.TestStep{
				Operation: logical.UpdateOperation,
				Path:      "sign/testing",
				Data: map[string]interface{}{
					"public_key": publicKey2,
					"key_id":     "override",
				},
				ErrorOk: true,
				Check: func(resp *logical.Response) error {
					if resp.Data["error"] != "setting key_id is not allowed by role" {
						return errors.New("Custom user key id was allowed even when 'allow_user_key_ids' is false.")
					}
					return nil
				},
			},
		},
	}

	logicaltest.Test(t, testCase)
}

func configCaStep() logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/ca",
		Data: map[string]interface{}{
			"public_key":  publicKey,
			"private_key": privateKey,
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
	role, keyId string, certType int, validPrincipals []string,
	criticalOptionPermissions, extensionPermissions map[string]string,
	ttl time.Duration,
	requestParameters map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "sign/" + role,
		Data:      requestParameters,

		Check: func(resp *logical.Response) error {

			serialNumber := resp.Data["serial_number"].(string)
			if serialNumber == "" {
				return errors.New("No serial number in response")
			}

			signedKey := strings.TrimSpace(resp.Data["signed_key"].(string))
			if signedKey == "" {
				return errors.New("No signed key in response")
			}

			key, _ := base64.StdEncoding.DecodeString(strings.Split(signedKey, " ")[1])

			parsedKey, err := ssh.ParsePublicKey(key)
			if err != nil {
				return err
			}

			return validateSSHCertificate(parsedKey.(*ssh.Certificate), keyId, certType, validPrincipals, criticalOptionPermissions, extensionPermissions, ttl)
		},
	}
}

func validateSSHCertificate(cert *ssh.Certificate, keyId string, certType int, validPrincipals []string, criticalOptionPermissions, extensionPermissions map[string]string,
	ttl time.Duration) error {

	if cert.KeyId != keyId {
		return fmt.Errorf("Incorrect KeyId: %v, wanted %v", cert.KeyId, keyId)
	}

	if cert.CertType != uint32(certType) {
		return fmt.Errorf("Incorrect CertType: %v", cert.CertType)
	}

	if time.Unix(int64(cert.ValidAfter), 0).After(time.Now()) {
		return fmt.Errorf("Incorrect ValidAfter: %v", cert.ValidAfter)
	}

	if time.Unix(int64(cert.ValidBefore), 0).Before(time.Now()) {
		return fmt.Errorf("Incorrect ValidBefore: %v", cert.ValidBefore)
	}

	actualTtl := time.Unix(int64(cert.ValidBefore), 0).Add(-30 * time.Second).Sub(time.Unix(int64(cert.ValidAfter), 0))
	if actualTtl != ttl {
		return fmt.Errorf("Incorrect ttl: expected: %v, actualL %v", ttl, actualTtl)
	}

	if !reflect.DeepEqual(cert.ValidPrincipals, validPrincipals) {
		return fmt.Errorf("Incorrect ValidPrincipals: expected: %#v actual: %#v", validPrincipals, cert.ValidPrincipals)
	}

	publicSigningKey, err := getSigningPublicKey()
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(cert.SignatureKey, publicSigningKey) {
		return fmt.Errorf("Incorrect SignatureKey: %v", cert.SignatureKey)
	}

	if cert.Signature == nil {
		return fmt.Errorf("Incorrect Signature: %v", cert.Signature)
	}

	if !reflect.DeepEqual(cert.Permissions.Extensions, extensionPermissions) {
		return fmt.Errorf("Incorrect Permissions.Extensions: Expected: %v, Actual: %v", extensionPermissions, cert.Permissions.Extensions)
	}

	if !reflect.DeepEqual(cert.Permissions.CriticalOptions, criticalOptionPermissions) {
		return fmt.Errorf("Incorrect Permissions.CriticalOptions: %v", cert.Permissions.CriticalOptions)
	}

	return nil
}

func getSigningPublicKey() (ssh.PublicKey, error) {
	key, err := base64.StdEncoding.DecodeString(strings.Split(publicKey, " ")[1])
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
