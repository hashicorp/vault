package vault

import (
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/credential"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
)

var (
	// invalidKey is used to test Unseal
	invalidKey = []byte("abcdefghijklmnopqrstuvwxyz")[:17]
)

func TestCore_Init(t *testing.T) {
	inm := physical.NewInmem()
	conf := &CoreConfig{Physical: inm}
	c, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	init, err := c.Initialized()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if init {
		t.Fatalf("should not be init")
	}

	// Check the seal configuration
	outConf, err := c.SealConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if outConf != nil {
		t.Fatalf("bad: %v", outConf)
	}

	sealConf := &SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}
	res, err := c.Initialize(sealConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(res.SecretShares) != 1 {
		t.Fatalf("Bad: %v", res)
	}
	if res.RootToken == "" {
		t.Fatalf("Bad: %v", res)
	}

	_, err = c.Initialize(sealConf)
	if err != ErrAlreadyInit {
		t.Fatalf("err: %v", err)
	}

	init, err = c.Initialized()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !init {
		t.Fatalf("should be init")
	}

	// Check the seal configuration
	outConf, err = c.SealConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(outConf, sealConf) {
		t.Fatalf("bad: %v expect: %v", outConf, sealConf)
	}

	// New Core, same backend
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	_, err = c2.Initialize(sealConf)
	if err != ErrAlreadyInit {
		t.Fatalf("err: %v", err)
	}

	init, err = c2.Initialized()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !init {
		t.Fatalf("should be init")
	}

	// Check the seal configuration
	outConf, err = c2.SealConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(outConf, sealConf) {
		t.Fatalf("bad: %v expect: %v", outConf, sealConf)
	}
}

func TestCore_Init_MultiShare(t *testing.T) {
	c := TestCore(t)
	sealConf := &SealConfig{
		SecretShares:    5,
		SecretThreshold: 3,
	}
	res, err := c.Initialize(sealConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(res.SecretShares) != 5 {
		t.Fatalf("Bad: %v", res)
	}
	if res.RootToken == "" {
		t.Fatalf("Bad: %v", res)
	}

	// Check the seal configuration
	outConf, err := c.SealConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(outConf, sealConf) {
		t.Fatalf("bad: %v expect: %v", outConf, sealConf)
	}
}

func TestCore_Unseal_MultiShare(t *testing.T) {
	c := TestCore(t)

	_, err := c.Unseal(invalidKey)
	if err != ErrNotInit {
		t.Fatalf("err: %v", err)
	}

	sealConf := &SealConfig{
		SecretShares:    5,
		SecretThreshold: 3,
	}
	res, err := c.Initialize(sealConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	sealed, err := c.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !sealed {
		t.Fatalf("should be sealed")
	}

	if prog := c.SecretProgress(); prog != 0 {
		t.Fatalf("bad progress: %d", prog)
	}

	for i := 0; i < 5; i++ {
		unseal, err := c.Unseal(res.SecretShares[i])
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Ignore redundant
		_, err = c.Unseal(res.SecretShares[i])
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i >= 2 {
			if !unseal {
				t.Fatalf("should be unsealed")
			}
			if prog := c.SecretProgress(); prog != 0 {
				t.Fatalf("bad progress: %d", prog)
			}
		} else {
			if unseal {
				t.Fatalf("should not be unsealed")
			}
			if prog := c.SecretProgress(); prog != i+1 {
				t.Fatalf("bad progress: %d", prog)
			}
		}
	}

	sealed, err = c.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if sealed {
		t.Fatalf("should not be sealed")
	}

	err = c.Seal()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ignore redundant
	err = c.Seal()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	sealed, err = c.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !sealed {
		t.Fatalf("should be sealed")
	}
}

func TestCore_Unseal_Single(t *testing.T) {
	c := TestCore(t)

	_, err := c.Unseal(invalidKey)
	if err != ErrNotInit {
		t.Fatalf("err: %v", err)
	}

	sealConf := &SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}
	res, err := c.Initialize(sealConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	sealed, err := c.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !sealed {
		t.Fatalf("should be sealed")
	}

	if prog := c.SecretProgress(); prog != 0 {
		t.Fatalf("bad progress: %d", prog)
	}

	unseal, err := c.Unseal(res.SecretShares[0])
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !unseal {
		t.Fatalf("should be unsealed")
	}
	if prog := c.SecretProgress(); prog != 0 {
		t.Fatalf("bad progress: %d", prog)
	}

	sealed, err = c.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if sealed {
		t.Fatalf("should not be sealed")
	}
}

func TestCore_Route_Sealed(t *testing.T) {
	c := TestCore(t)
	sealConf := &SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}

	// Should not route anything
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "sys/mounts",
	}
	_, err := c.HandleRequest(req)
	if err != ErrSealed {
		t.Fatalf("err: %v", err)
	}

	res, err := c.Initialize(sealConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	unseal, err := c.Unseal(res.SecretShares[0])
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Should not error after unseal
	req.ClientToken = res.RootToken
	_, err = c.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

// Attempt to unseal after doing a first seal
func TestCore_SealUnseal(t *testing.T) {
	c, key := TestCoreUnsealed(t)
	if err := c.Seal(); err != nil {
		t.Fatalf("err: %v", err)
	}
	if unseal, err := c.Unseal(key); err != nil || !unseal {
		t.Fatalf("err: %v", err)
	}
}

// Ensure we get a VaultID
func TestCore_HandleRequest_Lease(t *testing.T) {
	c, _, root := TestCoreUnsealedToken(t)

	req := &logical.Request{
		Operation: logical.WriteOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
		ClientToken: root,
	}
	resp, err := c.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read the key
	req.Operation = logical.ReadOperation
	req.Data = nil
	resp, err = c.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	if resp.Secret.Lease != time.Hour {
		t.Fatalf("bad: %#v", resp.Secret)
	}
	if resp.Secret.VaultID == "" {
		t.Fatalf("bad: %#v", resp.Secret)
	}
	if resp.Data["foo"] != "bar" {
		t.Fatalf("bad: %#v", resp.Data)
	}
}

func TestCore_HandleRequest_MissingToken(t *testing.T) {
	c, _, _ := TestCoreUnsealedToken(t)

	req := &logical.Request{
		Operation: logical.WriteOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
	}
	resp, err := c.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "missing client token" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestCore_HandleRequest_InvalidToken(t *testing.T) {
	c, _, _ := TestCoreUnsealedToken(t)

	req := &logical.Request{
		Operation: logical.WriteOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
		ClientToken: "foobarbaz",
	}
	resp, err := c.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "invalid client token" {
		t.Fatalf("bad: %#v", resp)
	}
}

// TODO: Test a root path is denied if non-root
//func TestCore_HandleRequest_RootPath(t *testing.T) {
//    c, _, _ := TestCoreUnsealedToken(t)
//    req := &logical.Request{
//        Operation: logical.WriteOperation,
//        Path:      "secret/test",
//        Data: map[string]interface{}{
//            "foo":   "bar",
//            "lease": "1h",
//        },
//        ClientToken: "foobarbaz",
//    }
//    resp, err := c.HandleRequest(req)
//    if err != logical.ErrInvalidRequest {
//        t.Fatalf("err: %v", err)
//    }
//    if resp.Data["error"] != "invalid client token" {
//        t.Fatalf("bad: %#v", resp)
//    }
//}

// TODO: Check that standard permissions work
//func TestCore_HandleRequest_PermissionDenied(t *testing.T) {
//    c, _, _ := TestCoreUnsealedToken(t)
//    req := &logical.Request{
//        Operation: logical.WriteOperation,
//        Path:      "secret/test",
//        Data: map[string]interface{}{
//            "foo":   "bar",
//            "lease": "1h",
//        },
//        ClientToken: "foobarbaz",
//    }
//    resp, err := c.HandleRequest(req)
//    if err != logical.ErrInvalidRequest {
//        t.Fatalf("err: %v", err)
//    }
//    if resp.Data["error"] != "invalid client token" {
//        t.Fatalf("bad: %#v", resp)
//    }
//}

// Ensure we get a client token
func TestCore_HandleLogin_Token(t *testing.T) {
	// Create a badass credential backend that always logs in as armon
	noop := &NoopCred{
		Login: []string{"login"},
		LoginResponse: &credential.Response{
			Secret: &logical.Secret{
				InternalData: map[string]interface{}{
					credential.PolicyKey:            []string{"foo", "bar"},
					credential.MetadataKey + "user": "armon",
				},
				Lease: time.Hour,
			},
		},
	}
	c, _, root := TestCoreUnsealedToken(t)
	c.credentialBackends["noop"] = func(map[string]string) (credential.Backend, error) {
		return noop, nil
	}

	// Enable the credential backend
	req := logical.TestRequest(t, logical.WriteOperation, "sys/auth/foo")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err := c.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to login
	lreq := &credential.Request{
		Path: "auth/foo/login",
	}
	lresp, err := c.HandleLogin(lreq)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ensure we got a client token back
	clientToken, ok := lresp.Data[clientTokenKey].(string)
	if !ok || clientToken == "" {
		t.Fatalf("bad: %#v", lresp)
	}

	// Check the policy and metadata
	te, err := c.tokenStore.Lookup(clientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expect := &TokenEntry{
		ID:       clientToken,
		Parent:   "",
		Policies: []string{"foo", "bar"},
		Path:     "auth/foo/login",
		Meta: map[string]interface{}{
			"user": "armon",
		},
	}
	if !reflect.DeepEqual(te, expect) {
		t.Fatalf("Bad: %#v expect: %#v", te, expect)
	}

	// Check that we have a lease with a VaultID
	if lresp.Secret.Lease != time.Hour {
		t.Fatalf("bad: %#v", lresp.Secret)
	}
	if lresp.Secret.VaultID == "" {
		t.Fatalf("bad: %#v", lresp.Secret)
	}
}
