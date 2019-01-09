package vault

import (
	"encoding/base64"
	"testing"

	"github.com/hashicorp/vault/helper/base62"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/xor"
)

func TestCore_GenerateRoot_Lifecycle(t *testing.T) {
	c, masterKeys, _ := TestCoreUnsealed(t)
	testCore_GenerateRoot_Lifecycle_Common(t, c, masterKeys)
}

func testCore_GenerateRoot_Lifecycle_Common(t *testing.T, c *Core, keys [][]byte) {
	// Verify update not allowed
	if _, err := c.GenerateRootUpdate(namespace.RootContext(nil), keys[0], "", GenerateStandardRootTokenStrategy); err == nil {
		t.Fatalf("no root generation in progress")
	}

	// Should be no progress
	num, err := c.GenerateRootProgress()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if num != 0 {
		t.Fatalf("bad: %d", num)
	}

	// Should be no config
	conf, err := c.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}

	// Cancel should be idempotent
	err = c.GenerateRootCancel()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	otp, err := base62.Random(26)
	if err != nil {
		t.Fatal(err)
	}

	// Start a root generation
	err = c.GenerateRootInit(otp, "", GenerateStandardRootTokenStrategy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should get config
	conf, err = c.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Cancel should be clear
	err = c.GenerateRootCancel()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be no config
	conf, err = c.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}
}

func TestCore_GenerateRoot_Init(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	testCore_GenerateRoot_Init_Common(t, c)

	bc, rc := TestSealDefConfigs()
	c, _, _, _ = TestCoreUnsealedWithConfigs(t, bc, rc)
	testCore_GenerateRoot_Init_Common(t, c)
}

func testCore_GenerateRoot_Init_Common(t *testing.T, c *Core) {
	otp, err := base62.Random(26)
	if err != nil {
		t.Fatal(err)
	}

	err = c.GenerateRootInit(otp, "", GenerateStandardRootTokenStrategy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Second should fail
	err = c.GenerateRootInit("", pgpkeys.TestPubKey1, GenerateStandardRootTokenStrategy)
	if err == nil {
		t.Fatalf("should fail")
	}
}

func TestCore_GenerateRoot_InvalidMasterNonce(t *testing.T) {
	c, masterKeys, _ := TestCoreUnsealed(t)
	// Pass in master keys as they'll be invalid
	masterKeys[0][0]++
	testCore_GenerateRoot_InvalidMasterNonce_Common(t, c, masterKeys)
}

func testCore_GenerateRoot_InvalidMasterNonce_Common(t *testing.T, c *Core, keys [][]byte) {
	otp, err := base62.Random(26)
	if err != nil {
		t.Fatal(err)
	}

	err = c.GenerateRootInit(otp, "", GenerateStandardRootTokenStrategy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Fetch new config with generated nonce
	rgconf, err := c.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rgconf == nil {
		t.Fatalf("bad: no rekey config received")
	}

	// Provide the nonce (invalid)
	_, err = c.GenerateRootUpdate(namespace.RootContext(nil), keys[0], "abcd", GenerateStandardRootTokenStrategy)
	if err == nil {
		t.Fatalf("expected error")
	}

	// Provide the master (invalid)
	for _, key := range keys {
		_, err = c.GenerateRootUpdate(namespace.RootContext(nil), key, rgconf.Nonce, GenerateStandardRootTokenStrategy)
	}
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCore_GenerateRoot_Update_OTP(t *testing.T) {
	c, masterKeys, _ := TestCoreUnsealed(t)
	testCore_GenerateRoot_Update_OTP_Common(t, c, masterKeys)
}

func testCore_GenerateRoot_Update_OTP_Common(t *testing.T, c *Core, keys [][]byte) {
	otp, err := base62.Random(26)
	if err != nil {
		t.Fatal(err)
	}

	// Start a root generation
	err = c.GenerateRootInit(otp, "", GenerateStandardRootTokenStrategy)
	if err != nil {
		t.Fatal(err)
	}

	// Fetch new config with generated nonce
	rkconf, err := c.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rkconf == nil {
		t.Fatalf("bad: no root generation config received")
	}

	// Provide the keys
	var result *GenerateRootResult
	for _, key := range keys {
		result, err = c.GenerateRootUpdate(namespace.RootContext(nil), key, rkconf.Nonce, GenerateStandardRootTokenStrategy)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result.EncodedToken != "" {
			break
		}
	}
	if result == nil {
		t.Fatalf("Bad, result is nil")
	}

	encodedToken := result.EncodedToken

	// Should be no progress
	num, err := c.GenerateRootProgress()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if num != 0 {
		t.Fatalf("bad: %d", num)
	}

	// Should be no config
	conf, err := c.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}

	tokenBytes, err := base64.RawStdEncoding.DecodeString(encodedToken)
	if err != nil {
		t.Fatal(err)
	}

	tokenBytes, err = xor.XORBytes(tokenBytes, []byte(otp))
	if err != nil {
		t.Fatal(err)
	}

	token := string(tokenBytes)

	// Ensure that the token is a root token
	te, err := c.tokenStore.Lookup(namespace.RootContext(nil), token)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te == nil {
		t.Fatalf("token was nil")
	}
	if te.ID != token || te.Parent != "" ||
		len(te.Policies) != 1 || te.Policies[0] != "root" {
		t.Fatalf("bad: %#v", *te)
	}
}

func TestCore_GenerateRoot_Update_PGP(t *testing.T) {
	c, masterKeys, _ := TestCoreUnsealed(t)
	testCore_GenerateRoot_Update_PGP_Common(t, c, masterKeys)
}

func testCore_GenerateRoot_Update_PGP_Common(t *testing.T, c *Core, keys [][]byte) {
	// Start a root generation
	err := c.GenerateRootInit("", pgpkeys.TestPubKey1, GenerateStandardRootTokenStrategy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Fetch new config with generated nonce
	rkconf, err := c.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rkconf == nil {
		t.Fatalf("bad: no root generation config received")
	}

	// Provide the keys
	var result *GenerateRootResult
	for _, key := range keys {
		result, err = c.GenerateRootUpdate(namespace.RootContext(nil), key, rkconf.Nonce, GenerateStandardRootTokenStrategy)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result.EncodedToken != "" {
			break
		}
	}
	if result == nil {
		t.Fatalf("Bad, result is nil")
	}

	encodedToken := result.EncodedToken

	// Should be no progress
	num, err := c.GenerateRootProgress()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if num != 0 {
		t.Fatalf("bad: %d", num)
	}

	// Should be no config
	conf, err := c.GenerateRootConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}

	ptBuf, err := pgpkeys.DecryptBytes(encodedToken, pgpkeys.TestPrivKey1)
	if err != nil {
		t.Fatal(err)
	}
	if ptBuf == nil {
		t.Fatal("Got nil plaintext key")
	}

	token := ptBuf.String()

	// Ensure that the token is a root token
	te, err := c.tokenStore.Lookup(namespace.RootContext(nil), token)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te == nil {
		t.Fatalf("token was nil")
	}
	if te.ID != token || te.Parent != "" ||
		len(te.Policies) != 1 || te.Policies[0] != "root" {
		t.Fatalf("bad: %#v", *te)
	}
}
