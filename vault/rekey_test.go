package vault

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/inmem"
)

func TestCore_Rekey_Lifecycle(t *testing.T) {
	bc, _ := TestSealDefConfigs()
	bc.SecretShares = 1
	bc.SecretThreshold = 1
	bc.StoredShares = 0
	c, masterKeys, _, _ := TestCoreUnsealedWithConfigs(t, bc, nil)
	if len(masterKeys) != 1 {
		t.Fatalf("expected %d keys, got %d", bc.SecretShares-bc.StoredShares, len(masterKeys))
	}
	testCore_Rekey_Lifecycle_Common(t, c, masterKeys, false)
}

func testCore_Rekey_Lifecycle_Common(t *testing.T, c *Core, masterKeys [][]byte, recovery bool) {
	// Verify update not allowed
	if _, err := c.RekeyUpdate(context.Background(), masterKeys[0], "", recovery); err == nil {
		t.Fatalf("no rekey should be in progress")
	}

	// Should be no progress
	if _, _, err := c.RekeyProgress(recovery, false); err == nil {
		t.Fatal("expected error from RekeyProgress")
	}

	// Should be no config
	conf, err := c.RekeyConfig(recovery)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}

	// Cancel should be idempotent
	err = c.RekeyCancel(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Start a rekey
	newConf := &SealConfig{
		SecretThreshold: 3,
		SecretShares:    5,
	}
	err = c.RekeyInit(newConf, recovery)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should get config
	conf, err = c.RekeyConfig(recovery)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	newConf.Nonce = conf.Nonce
	if !reflect.DeepEqual(conf, newConf) {
		t.Fatalf("bad: %v", conf)
	}

	// Cancel should be clear
	err = c.RekeyCancel(recovery)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be no config
	conf, err = c.RekeyConfig(recovery)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}
}

func TestCore_Rekey_Init(t *testing.T) {
	t.Run("barrier-rekey-init", func(t *testing.T) {
		c, _, _ := TestCoreUnsealed(t)
		testCore_Rekey_Init_Common(t, c, false)
	})
}

func testCore_Rekey_Init_Common(t *testing.T, c *Core, recovery bool) {
	// Try an invalid config
	badConf := &SealConfig{
		SecretThreshold: 5,
		SecretShares:    1,
	}
	err := c.RekeyInit(badConf, recovery)
	if err == nil {
		t.Fatalf("should fail")
	}

	// Start a rekey
	newConf := &SealConfig{
		SecretThreshold: 3,
		SecretShares:    5,
	}

	// If recovery key is supported, set newConf
	// to be a recovery seal config
	if c.seal.RecoveryKeySupported() {
		newConf.Type = c.seal.RecoveryType()
	}

	err = c.RekeyInit(newConf, recovery)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Second should fail
	err = c.RekeyInit(newConf, recovery)
	if err == nil {
		t.Fatalf("should fail")
	}
}

func TestCore_Rekey_Update(t *testing.T) {
	bc, _ := TestSealDefConfigs()
	bc.SecretShares = 1
	bc.SecretThreshold = 1
	bc.StoredShares = 0
	c, masterKeys, _, root := TestCoreUnsealedWithConfigs(t, bc, nil)
	testCore_Rekey_Update_Common(t, c, masterKeys, root, false)
}

func testCore_Rekey_Update_Common(t *testing.T, c *Core, keys [][]byte, root string, recovery bool) {
	var err error
	// Start a rekey
	var expType string
	if recovery {
		expType = c.seal.RecoveryType()
	} else {
		expType = c.seal.BarrierType()
	}

	newConf := &SealConfig{
		Type:            expType,
		SecretThreshold: 3,
		SecretShares:    5,
	}
	hErr := c.RekeyInit(newConf, recovery)
	if hErr != nil {
		t.Fatalf("err: %v", hErr)
	}

	// Fetch new config with generated nonce
	rkconf, hErr := c.RekeyConfig(recovery)
	if hErr != nil {
		t.Fatalf("err: %v", hErr)
	}
	if rkconf == nil {
		t.Fatalf("bad: no rekey config received")
	}

	// Provide the master/recovery keys
	var result *RekeyResult
	for _, key := range keys {
		result, err = c.RekeyUpdate(context.Background(), key, rkconf.Nonce, recovery)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if result != nil {
			break
		}
	}
	if result == nil {
		t.Fatal("nil result after update")
	}
	if newConf.StoredShares > 0 {
		if len(result.SecretShares) > 0 {
			t.Fatal("got secret shares when should have been storing")
		}
	} else if len(result.SecretShares) != newConf.SecretShares {
		t.Fatalf("rekey update error: %#v", result)
	}

	// Should be no progress
	if _, _, err := c.RekeyProgress(recovery, false); err == nil {
		t.Fatal("expected error from RekeyProgress")
	}

	// Should be no config
	conf, hErr := c.RekeyConfig(recovery)
	if hErr != nil {
		t.Fatalf("rekey config error: %v", hErr)
	}
	if conf != nil {
		t.Fatalf("rekey config should be nil, got: %v", conf)
	}

	// SealConfig should update
	var sealConf *SealConfig
	if recovery {
		sealConf, err = c.seal.RecoveryConfig(context.Background())
	} else {
		sealConf, err = c.seal.BarrierConfig(context.Background())
	}
	if err != nil {
		t.Fatalf("seal config retrieval error: %v", err)
	}
	if sealConf == nil {
		t.Fatal("seal configuration is nil")
	}

	newConf.Nonce = rkconf.Nonce
	if !reflect.DeepEqual(sealConf, newConf) {
		t.Fatalf("\nexpected: %#v\nactual: %#v\nexpType: %s\nrecovery: %t", newConf, sealConf, expType, recovery)
	}

	// At this point bail if we are rekeying the barrier key with recovery
	// keys, since a new rekey should still be using the same set of recovery
	// keys and we haven't been returned key shares in this mode.
	if !recovery && c.seal.RecoveryKeySupported() {
		return
	}

	// Attempt unseal if this was not recovery mode
	if !recovery {
		err = c.Seal(root)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		for i := 0; i < newConf.SecretThreshold; i++ {
			_, err = TestCoreUnseal(c, TestKeyCopy(result.SecretShares[i]))
			if err != nil {
				t.Fatalf("err: %v", err)
			}
		}
		if c.Sealed() {
			t.Fatalf("should be unsealed")
		}
	}

	// Start another rekey, this time we require a quorum!

	newConf = &SealConfig{
		Type:            expType,
		SecretThreshold: 1,
		SecretShares:    1,
	}
	err = c.RekeyInit(newConf, recovery)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Fetch new config with generated nonce
	rkconf, err = c.RekeyConfig(recovery)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rkconf == nil {
		t.Fatalf("bad: no rekey config received")
	}

	// Provide the parts master
	oldResult := result
	for i := 0; i < 3; i++ {
		result, err = c.RekeyUpdate(context.Background(), TestKeyCopy(oldResult.SecretShares[i]), rkconf.Nonce, recovery)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Should be progress
		if i < 2 {
			_, num, err := c.RekeyProgress(recovery, false)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if num != i+1 {
				t.Fatalf("bad: %d", num)
			}
		}
	}
	if result == nil || len(result.SecretShares) != 1 {
		t.Fatalf("Bad: %#v", result)
	}

	// Attempt unseal if this was not recovery mode
	if !recovery {
		err = c.Seal(root)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		unseal, err := TestCoreUnseal(c, result.SecretShares[0])
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// SealConfig should update
	if recovery {
		sealConf, err = c.seal.RecoveryConfig(context.Background())
	} else {
		sealConf, err = c.seal.BarrierConfig(context.Background())
	}
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	newConf.Nonce = rkconf.Nonce
	if !reflect.DeepEqual(sealConf, newConf) {
		t.Fatalf("bad: %#v", sealConf)
	}
}

func TestCore_Rekey_Invalid(t *testing.T) {
	bc, _ := TestSealDefConfigs()
	bc.StoredShares = 0
	bc.SecretShares = 1
	bc.SecretThreshold = 1
	c, masterKeys, _, _ := TestCoreUnsealedWithConfigs(t, bc, nil)
	testCore_Rekey_Invalid_Common(t, c, masterKeys, false)
}

func testCore_Rekey_Invalid_Common(t *testing.T, c *Core, keys [][]byte, recovery bool) {
	// Start a rekey
	newConf := &SealConfig{
		SecretThreshold: 3,
		SecretShares:    5,
	}
	err := c.RekeyInit(newConf, recovery)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Fetch new config with generated nonce
	rkconf, err := c.RekeyConfig(recovery)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rkconf == nil {
		t.Fatalf("bad: no rekey config received")
	}

	// Provide the nonce (invalid)
	_, err = c.RekeyUpdate(context.Background(), keys[0], "abcd", recovery)
	if err == nil {
		t.Fatalf("expected error")
	}

	// Provide the key (invalid)
	key := keys[0]
	oldkeystr := fmt.Sprintf("%#v", key)
	key[0]++
	newkeystr := fmt.Sprintf("%#v", key)
	ret, err := c.RekeyUpdate(context.Background(), key, rkconf.Nonce, recovery)
	if err == nil {
		t.Fatalf("expected error, ret is %#v\noldkeystr: %s\nnewkeystr: %s", *ret, oldkeystr, newkeystr)
	}

	// Check progress has been reset
	_, num, err := c.RekeyProgress(recovery, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if num != 0 {
		t.Fatalf("rekey progress should be 0, got: %d", num)
	}
}

func TestCore_Rekey_Standby(t *testing.T) {
	// Create the first core and initialize it
	logger := logging.NewVaultLogger(log.Trace)

	inm, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	redirectOriginal := "http://127.0.0.1:8200"
	core, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal,
		DisableMlock: true,
		DisableCache: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	keys, root := TestCoreInit(t, core)
	for _, key := range keys {
		if _, err := TestCoreUnseal(core, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Wait for core to become active
	TestWaitActive(t, core)

	// Create a second core, attached to same in-memory store
	redirectOriginal2 := "http://127.0.0.1:8500"
	core2, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal2,
		DisableMlock: true,
		DisableCache: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for _, key := range keys {
		if _, err := TestCoreUnseal(core2, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Rekey the master key
	newConf := &SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}
	err = core.RekeyInit(newConf, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Fetch new config with generated nonce
	rkconf, err := core.RekeyConfig(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rkconf == nil {
		t.Fatalf("bad: no rekey config received")
	}
	var rekeyResult *RekeyResult
	for _, key := range keys {
		rekeyResult, err = core.RekeyUpdate(context.Background(), key, rkconf.Nonce, false)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}
	if rekeyResult == nil {
		t.Fatalf("rekey failed")
	}

	// Seal the first core, should step down
	err = core.Seal(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait for core2 to become active
	TestWaitActive(t, core2)

	// Rekey the master key again
	err = core2.RekeyInit(newConf, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Fetch new config with generated nonce
	rkconf, err = core2.RekeyConfig(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rkconf == nil {
		t.Fatalf("bad: no rekey config received")
	}
	var rekeyResult2 *RekeyResult
	for _, key := range rekeyResult.SecretShares {
		rekeyResult2, err = core2.RekeyUpdate(context.Background(), key, rkconf.Nonce, false)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}
	if rekeyResult2 == nil {
		t.Fatalf("rekey failed")
	}

	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rekeyResult2 == nil {
		t.Fatalf("rekey failed")
	}
}
