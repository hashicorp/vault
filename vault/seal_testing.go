package vault

import (
	"context"

	log "github.com/hashicorp/go-hclog"
	testing "github.com/mitchellh/go-testing-interface"
)

type TestSealOpts struct {
	Logger     log.Logger
	StoredKeys StoredKeysSupport
	Secret     []byte
}

func TestCoreUnsealedWithConfigs(t testing.T, barrierConf, recoveryConf *SealConfig) (*Core, [][]byte, [][]byte, string) {
	t.Helper()
	opts := &TestSealOpts{}
	if recoveryConf == nil {
		opts.StoredKeys = StoredKeysSupportedShamirMaster
	}
	return TestCoreUnsealedWithConfigSealOpts(t, barrierConf, recoveryConf, opts)
}

func TestCoreUnsealedWithConfigSealOpts(t testing.T, barrierConf, recoveryConf *SealConfig, sealOpts *TestSealOpts) (*Core, [][]byte, [][]byte, string) {
	t.Helper()
	seal := NewTestSeal(t, sealOpts)
	core := TestCoreWithSeal(t, seal, false)
	result, err := core.Initialize(context.Background(), &InitParams{
		BarrierConfig:    barrierConf,
		RecoveryConfig:   recoveryConf,
		LegacyShamirSeal: sealOpts.StoredKeys == StoredKeysNotSupported,
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
