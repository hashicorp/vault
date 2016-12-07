// +build vault

package vault

import (
	"bytes"
	"fmt"
	"testing"
)

type TestSeal struct {
	defseal        *DefaultSeal
	barrierKeys    [][]byte
	recoveryKey    []byte
	recoveryConfig *SealConfig
}

func newTestSeal(t *testing.T) Seal {
	return &TestSeal{}
}

func (d *TestSeal) checkCore() error {
	if d.defseal.core == nil {
		return fmt.Errorf("seal does not have a core set")
	}
	return nil
}

func (d *TestSeal) SetCore(core *Core) {
	d.defseal = &DefaultSeal{}
	d.defseal.core = core
}

func (d *TestSeal) Init() error {
	d.barrierKeys = [][]byte{}
	return d.defseal.Init()
}

func (d *TestSeal) Finalize() error {
	return d.defseal.Finalize()
}

func (d *TestSeal) BarrierType() string {
	return "shamir"
}

func (d *TestSeal) StoredKeysSupported() bool {
	return true
}

func (d *TestSeal) RecoveryKeySupported() bool {
	return true
}

func (d *TestSeal) SetStoredKeys(keys [][]byte) error {
	d.barrierKeys = keys
	return nil
}

func (d *TestSeal) GetStoredKeys() ([][]byte, error) {
	return d.barrierKeys, nil
}

func (d *TestSeal) BarrierConfig() (*SealConfig, error) {
	return d.defseal.BarrierConfig()
}

func (d *TestSeal) SetBarrierConfig(config *SealConfig) error {
	return d.defseal.SetBarrierConfig(config)
}

func (d *TestSeal) RecoveryType() string {
	return "shamir"
}

func (d *TestSeal) RecoveryConfig() (*SealConfig, error) {
	return d.recoveryConfig, nil
}

func (d *TestSeal) SetRecoveryConfig(config *SealConfig) error {
	d.recoveryConfig = config
	return nil
}

func (d *TestSeal) VerifyRecoveryKey(key []byte) error {
	if bytes.Equal(d.recoveryKey, key) {
		return nil
	}
	return fmt.Errorf("not equivalent")
}

func (d *TestSeal) SetRecoveryKey(key []byte) error {
	newbuf := bytes.NewBuffer(nil)
	newbuf.Write(key)
	d.recoveryKey = newbuf.Bytes()
	return nil
}

func TestCoreUnsealedWithConfigs(t *testing.T, barrierConf, recoveryConf *SealConfig) (*Core, [][]byte, [][]byte, string) {
	seal := &TestSeal{}
	core := TestCoreWithSeal(t, seal)
	result, err := core.Initialize(&InitParams{
		BarrierConfig:  barrierConf,
		RecoveryConfig: recoveryConf,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = core.UnsealWithStoredKeys()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if sealed, _ := core.Sealed(); sealed {
		for _, key := range result.SecretShares {
			if _, err := core.Unseal(key); err != nil {

				t.Fatalf("unseal err: %s", err)
			}
		}

		sealed, err = core.Sealed()
		if err != nil {
			t.Fatalf("err checking seal status: %s", err)
		}
		if sealed {
			t.Fatal("should not be sealed")
		}
	}

	return core, result.SecretShares, result.RecoveryShares, result.RootToken
}

func TestSealDefConfigs() (*SealConfig, *SealConfig) {
	return &SealConfig{
			SecretShares:    5,
			SecretThreshold: 3,
			StoredShares:    2,
		}, &SealConfig{
			SecretShares:    5,
			SecretThreshold: 3,
		}
}
