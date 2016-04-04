package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/physical"
)

func TestCore_Rekey_Lifecycle(t *testing.T) {
	c, master, _ := TestCoreUnsealed(t)

	// Verify update not allowed
	if _, err := c.RekeyUpdate(master, "", false); err == nil {
		t.Fatalf("no rekey in progress")
	}

	// Should be no progress
	num, err := c.RekeyProgress(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if num != 0 {
		t.Fatalf("bad: %d", num)
	}

	// Should be no config
	conf, err := c.RekeyConfig(false)
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
	err = c.RekeyInit(newConf, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should get config
	conf, err = c.RekeyConfig(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	newConf.Nonce = conf.Nonce
	if !reflect.DeepEqual(conf, newConf) {
		t.Fatalf("bad: %v", conf)
	}

	// Cancel should be clear
	err = c.RekeyCancel(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should be no config
	conf, err = c.RekeyConfig(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}
}

func TestCore_Rekey_Init(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	// Try an invalid config
	badConf := &SealConfig{
		SecretThreshold: 5,
		SecretShares:    1,
	}
	err := c.RekeyInit(badConf, false)
	if err == nil {
		t.Fatalf("should fail")
	}

	// Start a rekey
	newConf := &SealConfig{
		SecretThreshold: 3,
		SecretShares:    5,
	}
	err = c.RekeyInit(newConf, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Second should fail
	err = c.RekeyInit(newConf, false)
	if err == nil {
		t.Fatalf("should fail")
	}
}

func TestCore_Rekey_Update(t *testing.T) {
	c, master, root := TestCoreUnsealed(t)

	// Start a rekey
	newConf := &SealConfig{
		SecretThreshold: 3,
		SecretShares:    5,
	}
	err := c.RekeyInit(newConf, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Fetch new config with generated nonce
	rkconf, err := c.RekeyConfig(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rkconf == nil {
		t.Fatalf("bad: no rekey config received")
	}

	// Provide the master
	result, err := c.RekeyUpdate(master, rkconf.Nonce, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if result == nil || len(result.SecretShares) != 5 {
		t.Fatalf("Bad: %#v", result)
	}

	// Should be no progress
	num, err := c.RekeyProgress(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if num != 0 {
		t.Fatalf("bad: %d", num)
	}

	// Should be no config
	conf, err := c.RekeyConfig(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if conf != nil {
		t.Fatalf("bad: %v", conf)
	}

	// SealConfig should update
	sealConf, err := c.seal.BarrierConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if sealConf == nil {
		t.Fatal("seal configuration is nil")
	}

	newConf.Nonce = rkconf.Nonce
	if !reflect.DeepEqual(sealConf, newConf) {
		t.Fatalf("\nexpected: %#v\nactual: %#v\n", sealConf, newConf)
	}

	// Attempt unseal
	err = c.Seal(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i := 0; i < 3; i++ {
		_, err = c.Unseal(result.SecretShares[i])
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}
	if sealed, _ := c.Sealed(); sealed {
		t.Fatalf("should be unsealed")
	}

	// Start another rekey, this time we require a quorum!
	newConf = &SealConfig{
		SecretThreshold: 1,
		SecretShares:    1,
	}
	err = c.RekeyInit(newConf, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Fetch new config with generated nonce
	rkconf, err = c.RekeyConfig(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rkconf == nil {
		t.Fatalf("bad: no rekey config received")
	}

	// Provide the parts master
	oldResult := result
	for i := 0; i < 3; i++ {
		result, err = c.RekeyUpdate(oldResult.SecretShares[i], rkconf.Nonce, false)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Should be progress
		num, err := c.RekeyProgress(false)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if (i == 2 && num != 0) || (i != 2 && num != i+1) {
			t.Fatalf("bad: %d", num)
		}
	}
	if result == nil || len(result.SecretShares) != 1 {
		t.Fatalf("Bad: %#v", result)
	}

	// Attempt unseal
	err = c.Seal(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	unseal, err := c.Unseal(result.SecretShares[0])
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// SealConfig should update
	sealConf, err = c.seal.BarrierConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	newConf.Nonce = rkconf.Nonce
	if !reflect.DeepEqual(sealConf, newConf) {
		t.Fatalf("bad: %#v", sealConf)
	}
}

func TestCore_Rekey_InvalidMaster(t *testing.T) {
	c, master, _ := TestCoreUnsealed(t)

	// Start a rekey
	newConf := &SealConfig{
		SecretThreshold: 3,
		SecretShares:    5,
	}
	err := c.RekeyInit(newConf, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Fetch new config with generated nonce
	rkconf, err := c.RekeyConfig(false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if rkconf == nil {
		t.Fatalf("bad: no rekey config received")
	}

	// Provide the master (invalid)
	master[0]++
	_, err = c.RekeyUpdate(master, rkconf.Nonce, false)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCore_Rekey_InvalidNonce(t *testing.T) {
	c, master, _ := TestCoreUnsealed(t)

	// Start a rekey
	newConf := &SealConfig{
		SecretThreshold: 3,
		SecretShares:    5,
	}
	err := c.RekeyInit(newConf, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Provide the nonce (invalid)
	_, err = c.RekeyUpdate(master, "abcd", false)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCore_Standby_Rekey(t *testing.T) {
	// Create the first core and initialize it
	inm := physical.NewInmem()
	inmha := physical.NewInmemHA()
	advertiseOriginal := "http://127.0.0.1:8200"
	core, err := NewCore(&CoreConfig{
		Physical:      inm,
		HAPhysical:    inmha,
		AdvertiseAddr: advertiseOriginal,
		DisableMlock:  true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	key, root := TestCoreInit(t, core)
	if _, err := core.Unseal(TestKeyCopy(key)); err != nil {
		t.Fatalf("unseal err: %s", err)
	}

	// Wait for core to become active
	testWaitActive(t, core)

	// Create a second core, attached to same in-memory store
	advertiseOriginal2 := "http://127.0.0.1:8500"
	core2, err := NewCore(&CoreConfig{
		Physical:      inm,
		HAPhysical:    inmha,
		AdvertiseAddr: advertiseOriginal2,
		DisableMlock:  true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if _, err := core2.Unseal(TestKeyCopy(key)); err != nil {
		t.Fatalf("unseal err: %s", err)
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
	result, err := core.RekeyUpdate(key, rkconf.Nonce, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if result == nil {
		t.Fatalf("rekey failed")
	}

	// Seal the first core, should step down
	err = core.Seal(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait for core2 to become active
	testWaitActive(t, core2)

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
	result, err = core2.RekeyUpdate(result.SecretShares[0], rkconf.Nonce, false)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if result == nil {
		t.Fatalf("rekey failed")
	}
}
