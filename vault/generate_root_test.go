package vault

import (
	"encoding/base64"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/xor"
)

func TestCore_GenerateRoot_Lifecycle(t *testing.T) {
	bc, rc := TestSealDefConfigs()
	c, masterKeys, _, _ := TestCoreUnsealedWithConfigs(t, bc, rc)
	testCore_GenerateRoot_Lifecycle_Common(t, c, masterKeys)
}

func testCore_GenerateRoot_Lifecycle_Common(t *testing.T, c *Core, keys [][]byte) {
	// Verify update not allowed
	if _, err := c.GenerateRootUpdate(keys[0], "", GenerateStandardRootTokenStrategy); err == nil {
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

	otpBytes, err := GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}

	// Start a root generation
	err = c.GenerateRootInit(base64.StdEncoding.EncodeToString(otpBytes), "", GenerateStandardRootTokenStrategy)
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
	otpBytes, err := GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}

	err = c.GenerateRootInit(base64.StdEncoding.EncodeToString(otpBytes), "", GenerateStandardRootTokenStrategy)
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
	bc, _ := TestSealDefConfigs()
	bc.SecretShares = 3
	bc.SecretThreshold = 3
	c, masterKeys, _, _ := TestCoreUnsealedWithConfigs(t, bc, nil)
	// Pass in master keys as they'll be invalid
	masterKeys[0][0]++
	testCore_GenerateRoot_InvalidMasterNonce_Common(t, c, masterKeys)
}

func testCore_GenerateRoot_InvalidMasterNonce_Common(t *testing.T, c *Core, keys [][]byte) {
	otpBytes, err := GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}

	err = c.GenerateRootInit(base64.StdEncoding.EncodeToString(otpBytes), "", GenerateStandardRootTokenStrategy)
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
	_, err = c.GenerateRootUpdate(keys[0], "abcd", GenerateStandardRootTokenStrategy)
	if err == nil {
		t.Fatalf("expected error")
	}

	// Provide the master (invalid)
	for _, key := range keys {
		_, err = c.GenerateRootUpdate(key, rgconf.Nonce, GenerateStandardRootTokenStrategy)
	}
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCore_GenerateRoot_Update_OTP(t *testing.T) {
	bc, rc := TestSealDefConfigs()
	c, masterKeys, _, _ := TestCoreUnsealedWithConfigs(t, bc, rc)
	testCore_GenerateRoot_Update_OTP_Common(t, c, masterKeys[0:bc.SecretThreshold])
}

func testCore_GenerateRoot_Update_OTP_Common(t *testing.T, c *Core, keys [][]byte) {
	otpBytes, err := GenerateRandBytes(16)
	if err != nil {
		t.Fatal(err)
	}

	otp := base64.StdEncoding.EncodeToString(otpBytes)
	// Start a root generation
	err = c.GenerateRootInit(otp, "", GenerateStandardRootTokenStrategy)
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
		result, err = c.GenerateRootUpdate(key, rkconf.Nonce, GenerateStandardRootTokenStrategy)
		if err != nil {
			t.Fatalf("err: %v", err)
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

	tokenBytes, err := xor.XORBase64(encodedToken, otp)
	if err != nil {
		t.Fatal(err)
	}
	token, err := uuid.FormatUUID(tokenBytes)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that the token is a root token
	te, err := c.tokenStore.Lookup(token)
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
	bc, rc := TestSealDefConfigs()
	c, masterKeys, _, _ := TestCoreUnsealedWithConfigs(t, bc, rc)
	testCore_GenerateRoot_Update_PGP_Common(t, c, masterKeys[0:bc.SecretThreshold])
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
		result, err = c.GenerateRootUpdate(key, rkconf.Nonce, GenerateStandardRootTokenStrategy)
		if err != nil {
			t.Fatalf("err: %v", err)
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
	te, err := c.tokenStore.Lookup(token)
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
