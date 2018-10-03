package vault

import (
	"context"
	"reflect"
	"testing"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical/inmem"
)

func TestCore_Init(t *testing.T) {
	c, conf := testCore_NewTestCore(t, nil)
	testCore_Init_Common(t, c, conf, &SealConfig{SecretShares: 5, SecretThreshold: 3}, nil)
}

func testCore_NewTestCore(t *testing.T, seal Seal) (*Core, *CoreConfig) {
	return testCore_NewTestCoreLicensing(t, seal, nil)
}

func testCore_NewTestCoreLicensing(t *testing.T, seal Seal, licensingConfig *LicensingConfig) (*Core, *CoreConfig) {
	logger := logging.NewVaultLogger(log.Trace)

	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	conf := &CoreConfig{
		Physical:     inm,
		DisableMlock: true,
		LogicalBackends: map[string]logical.Factory{
			"kv": LeasedPassthroughBackendFactory,
		},
		Seal:            seal,
		LicensingConfig: licensingConfig,
	}
	c, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	return c, conf
}

func testCore_Init_Common(t *testing.T, c *Core, conf *CoreConfig, barrierConf, recoveryConf *SealConfig) {
	init, err := c.Initialized(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if init {
		t.Fatalf("should not be init")
	}

	// Check the seal configuration
	outConf, err := c.seal.BarrierConfig(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if outConf != nil {
		t.Fatalf("bad: %v", outConf)
	}
	if recoveryConf != nil {
		outConf, err := c.seal.RecoveryConfig(context.Background())
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if outConf != nil {
			t.Fatalf("bad: %v", outConf)
		}
	}

	res, err := c.Initialize(context.Background(), &InitParams{
		BarrierConfig:  barrierConf,
		RecoveryConfig: recoveryConf,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(res.SecretShares) != (barrierConf.SecretShares - barrierConf.StoredShares) {
		t.Fatalf("Bad: got\n%#v\nexpected conf matching\n%#v\n", *res, *barrierConf)
	}
	if recoveryConf != nil {
		if len(res.RecoveryShares) != recoveryConf.SecretShares {
			t.Fatalf("Bad: got\n%#v\nexpected conf matching\n%#v\n", *res, *recoveryConf)
		}
	}

	if res.RootToken == "" {
		t.Fatalf("Bad: %#v", res)
	}

	_, err = c.Initialize(context.Background(), &InitParams{
		BarrierConfig:  barrierConf,
		RecoveryConfig: recoveryConf,
	})
	if err != ErrAlreadyInit {
		t.Fatalf("err: %v", err)
	}

	init, err = c.Initialized(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !init {
		t.Fatalf("should be init")
	}

	// Check the seal configuration
	outConf, err = c.seal.BarrierConfig(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(outConf, barrierConf) {
		t.Fatalf("bad: %v expect: %v", outConf, barrierConf)
	}
	if recoveryConf != nil {
		outConf, err = c.seal.RecoveryConfig(context.Background())
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if !reflect.DeepEqual(outConf, recoveryConf) {
			t.Fatalf("bad: %v expect: %v", outConf, recoveryConf)
		}
	}

	// New Core, same backend
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	_, err = c2.Initialize(context.Background(), &InitParams{
		BarrierConfig:  barrierConf,
		RecoveryConfig: recoveryConf,
	})
	if err != ErrAlreadyInit {
		t.Fatalf("err: %v", err)
	}

	init, err = c2.Initialized(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !init {
		t.Fatalf("should be init")
	}

	// Check the seal configuration
	outConf, err = c2.seal.BarrierConfig(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(outConf, barrierConf) {
		t.Fatalf("bad: %v expect: %v", outConf, barrierConf)
	}
	if recoveryConf != nil {
		outConf, err = c2.seal.RecoveryConfig(context.Background())
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if !reflect.DeepEqual(outConf, recoveryConf) {
			t.Fatalf("bad: %v expect: %v", outConf, recoveryConf)
		}
	}
}
