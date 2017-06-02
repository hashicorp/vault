package vault

import (
	"testing"

	"github.com/hashicorp/vault/helper/pgpkeys"
)

func TestCore_GenerateShare_Lifecycle(t *testing.T) {
	bc, rc := TestSealDefConfigs()
	c, masterKeys, _, _ := TestCoreUnsealedWithConfigs(t, bc, rc)
	c.seal.(*TestSeal).recoveryKeysDisabled = true
	testCore_GenerateShare_Lifecycle_Common(t, c, masterKeys)

	bc, rc = TestSealDefConfigs()
	c, _, recoveryKeys, _ := TestCoreUnsealedWithConfigs(t, bc, rc)
	testCore_GenerateShare_Lifecycle_Common(t, c, recoveryKeys)
}

func testCore_GenerateShare_Lifecycle_Common(t *testing.T, c *Core, keys [][]byte) {
	// Verify update not allowed
	if _, err := c.GenerateShareUpdate(keys[0]); err == nil {
		t.Fatalf("no share generation in progress")
	}

	// Should be no progress
	num, err := c.GenerateShareProgress()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if num != 0 {
		t.Fatalf("bad: %d", num)
	}

	// Should be no config
	conf, err := c.GenerateShareConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}

	// Cancel should be idempotent
	err = c.GenerateShareCancel()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Start a share generation
	err = c.GenerateShareInit(pgpkeys.TestPubKey1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should get config
	conf, err = c.GenerateShareConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Cancel should be clear
	err = c.GenerateShareCancel()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be no config
	conf, err = c.GenerateShareConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}
}

func TestCore_GenerateShare_Init(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	testCore_GenerateShare_Init_Common(t, c)

	bc, rc := TestSealDefConfigs()
	c, _, _, _ = TestCoreUnsealedWithConfigs(t, bc, rc)
	testCore_GenerateShare_Init_Common(t, c)
}

func testCore_GenerateShare_Init_Common(t *testing.T, c *Core) {
	err := c.GenerateShareInit(pgpkeys.TestPubKey1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Second should fail
	err = c.GenerateShareInit(pgpkeys.TestPubKey1)
	if err == nil {
		t.Fatalf("should fail")
	}
}

func TestCore_GenerateShare_Update_PGP(t *testing.T) {
	bc, rc := TestSealDefConfigs()
	c, masterKeys, _, _ := TestCoreUnsealedWithConfigs(t, bc, rc)
	c.seal.(*TestSeal).recoveryKeysDisabled = true
	testCore_GenerateShare_Update_PGP_Common(t, c, masterKeys[0:bc.SecretThreshold])

	bc, rc = TestSealDefConfigs()
	c, _, recoveryKeys, _ := TestCoreUnsealedWithConfigs(t, bc, rc)
	testCore_GenerateShare_Update_PGP_Common(t, c, recoveryKeys[0:rc.SecretThreshold])
}

func testCore_GenerateShare_Update_PGP_Common(t *testing.T, c *Core, keys [][]byte) {
	// Start a share generation
	err := c.GenerateShareInit(pgpkeys.TestPubKey1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Fetch new config with generated nonce
	rkconf, err := c.GenerateShareConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rkconf == nil {
		t.Fatalf("bad: no share generation config received")
	}

	// Provide the keys
	var result *GenerateShareResult
	for _, key := range keys {
		result, err = c.GenerateShareUpdate(key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}
	if result == nil {
		t.Fatalf("Bad, result is nil")
	}

	newShare := result.Key

	// Should be no progress
	num, err := c.GenerateShareProgress()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if num != 0 {
		t.Fatalf("bad: %d", num)
	}

	// Should be no config
	conf, err := c.GenerateShareConfiguration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}

	ptBuf, err := pgpkeys.DecryptBytes(newShare, pgpkeys.TestPrivKey1)
	if err != nil {
		t.Fatal(err)
	}
	if ptBuf == nil {
		t.Fatal("Got nil plaintext key")
	}
}
