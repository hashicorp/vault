package vault

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
)

func TestCore_Init(t *testing.T) {
	c, conf := testCore_NewTestCore(t, nil)
	testCore_Init_Common(t, c, conf, &SealConfig{SecretShares: 5, SecretThreshold: 3}, nil)

	c, conf = testCore_NewTestCore(t, &TestSeal{})
	bc, rc := TestSealDefConfigs()
	rc.SecretShares = 4
	rc.SecretThreshold = 2
	testCore_Init_Common(t, c, conf, bc, rc)
}

func testCore_NewTestCore(t *testing.T, seal Seal) (*Core, *CoreConfig) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	inm := physical.NewInmem(logger)
	conf := &CoreConfig{
		Physical:     inm,
		DisableMlock: true,
		LogicalBackends: map[string]logical.Factory{
			"generic": LeasedPassthroughBackendFactory,
		},
		Seal: seal,
	}
	c, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	return c, conf
}

func testCore_Init_Common(t *testing.T, c *Core, conf *CoreConfig, barrierConf, recoveryConf *SealConfig) {
	init, err := c.Initialized()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if init {
		t.Fatalf("should not be init")
	}

	// Check the seal configuration
	outConf, err := c.seal.BarrierConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if outConf != nil {
		t.Fatalf("bad: %v", outConf)
	}
	if recoveryConf != nil {
		outConf, err := c.seal.RecoveryConfig()
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if outConf != nil {
			t.Fatalf("bad: %v", outConf)
		}
	}

	res, err := c.Initialize(barrierConf, recoveryConf)
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

	_, err = c.Initialize(barrierConf, recoveryConf)
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
	outConf, err = c.seal.BarrierConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(outConf, barrierConf) {
		t.Fatalf("bad: %v expect: %v", outConf, barrierConf)
	}
	if recoveryConf != nil {
		outConf, err = c.seal.RecoveryConfig()
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

	_, err = c2.Initialize(barrierConf, recoveryConf)
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
	outConf, err = c2.seal.BarrierConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(outConf, barrierConf) {
		t.Fatalf("bad: %v expect: %v", outConf, barrierConf)
	}
	if recoveryConf != nil {
		outConf, err = c2.seal.RecoveryConfig()
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if !reflect.DeepEqual(outConf, recoveryConf) {
			t.Fatalf("bad: %v expect: %v", outConf, recoveryConf)
		}
	}
}
