package vault

import (
	"context"

	testing "github.com/mitchellh/go-testing-interface"
)

var (
	TestCoreUnsealedWithConfigs = testCoreUnsealedWithConfigs
	TestSealDefConfigs          = testSealDefConfigs
)

type TestSealOpts struct {
	StoredKeysDisabled   bool
	RecoveryKeysDisabled bool
	Secret               []byte
}

func testCoreUnsealedWithConfigs(t testing.T, barrierConf, recoveryConf *SealConfig) (*Core, [][]byte, [][]byte, string) {
	t.Helper()
	var opts *TestSealOpts
	if recoveryConf == nil {
		opts = &TestSealOpts{
			StoredKeysDisabled:   true,
			RecoveryKeysDisabled: true,
		}
	}
	seal := NewTestSeal(t, opts)
	core := TestCoreWithSeal(t, seal, false)
	result, err := core.Initialize(context.Background(), &InitParams{
		BarrierConfig:  barrierConf,
		RecoveryConfig: recoveryConf,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = core.UnsealWithStoredKeys(context.Background())
	if err != nil && IsFatalError(err) {
		t.Fatalf("err: %s", err)
	}
	if core.Sealed() {
		for _, key := range result.SecretShares {
			if _, err := core.Unseal(TestKeyCopy(key)); err != nil {
				t.Fatalf("unseal err: %s", err)
			}
		}

		if core.Sealed() {
			t.Fatal("should not be sealed")
		}
	}

	return core, result.SecretShares, result.RecoveryShares, result.RootToken
}

func testSealDefConfigs() (*SealConfig, *SealConfig) {
	return &SealConfig{
		SecretShares:    5,
		SecretThreshold: 3,
	}, nil
}

func TestCoreUnsealedWithConfigSealOpts(t testing.T, barrierConf, recoveryConf *SealConfig, sealOpts *TestSealOpts) (*Core, [][]byte, [][]byte, string) {
	seal := NewTestSeal(t, sealOpts)
	core := TestCoreWithSeal(t, seal, false)
	result, err := core.Initialize(context.Background(), &InitParams{
		BarrierConfig:  barrierConf,
		RecoveryConfig: recoveryConf,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = core.UnsealWithStoredKeys(context.Background())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if core.Sealed() {
		for _, key := range result.SecretShares {
			if _, err := core.Unseal(TestKeyCopy(key)); err != nil {
				t.Fatalf("unseal err: %s", err)
			}
		}

		if core.Sealed() {
			t.Fatal("should not be sealed")
		}
	}

	return core, result.SecretShares, result.RecoveryShares, result.RootToken
}
