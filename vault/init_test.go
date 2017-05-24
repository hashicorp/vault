package vault

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
)

func TestCore_Init(t *testing.T) {
	c, conf := testCore_NewTestCore(t, nil)
	testCore_Init_Common(t, c, conf, &SealConfig{SecretShares: 5, SecretThreshold: 3}, nil)

	c, conf = testCore_NewTestCore(t, newTestSeal(t))
	bc, rc := TestSealDefConfigs()
	rc.SecretShares = 4
	rc.SecretThreshold = 2
	testCore_Init_Common(t, c, conf, bc, rc)

	c, conf = testCore_NewTestCore(t, newTestSeal(t))
	bc, rc = TestSealDefConfigs()
	bc.WrapShares = true
	rc.WrapShares = true
	testCore_Init_Common(t, c, conf, bc, rc)
}

func testCore_NewTestCore(t *testing.T, seal Seal) (*Core, *CoreConfig) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

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

	res, err := c.Initialize(&InitParams{
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

		if recoveryConf.WrapShares {
			testValidateWrappedShare(t, c, string(res.RecoveryShares[0][:]), recoveryConf, "init")
		}
	}

	if barrierConf.WrapShares {
		if res.RootToken != "" {
			t.Fatalf("Bad: %#v", res)
		}

		testValidateWrappedShare(t, c, string(res.SecretShares[0][:]), barrierConf, "init")
	} else {
		if res.RootToken == "" {
			t.Fatalf("Bad: %#v", res)
		}
	}

	_, err = c.Initialize(&InitParams{
		BarrierConfig:  barrierConf,
		RecoveryConfig: recoveryConf,
	})
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

	_, err = c2.Initialize(&InitParams{
		BarrierConfig:  barrierConf,
		RecoveryConfig: recoveryConf,
	})
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

func testValidateWrappedShare(t testing.TB, c *Core, token string, barrierConf *SealConfig, expectedMethod string) string {
	// Make sure tokens are JWT formatted
	if strings.Count(token, ".") != 2 {
		t.Fatalf("Bad: %#v", token)
	}

	req := &logical.Request{
		ClientToken: token,
		Operation:   logical.ReadOperation,
		Path:        "cubbyhole/response",
	}

	ok, err := c.ValidateWrappingToken(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !ok {
		t.Fatalf("bad: %#v", req)
	}

	cubbyResp, err := c.router.Route(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if cubbyResp == nil {
		t.Fatalf("Nothing in cubbyhole")
	}
	if cubbyResp.IsError() {
		t.Fatalf("err: %#v", cubbyResp)
	}
	if cubbyResp.Data == nil {
		t.Fatalf("wrapping information was nil")
	}

	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(cubbyResp.Data["response"].(string)), &m)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	mData := m["data"].(map[string]interface{})

	keyShares := mData["key-shares"].(float64)
	keyThres := mData["key-threshold"].(float64)
	method := mData["method"].(string)
	share := mData["share"].(string)

	if int(keyShares) != barrierConf.SecretShares {
		t.Fatalf("Unexpected number of key shares: got %d, expected %d", keyShares, barrierConf.SecretShares)
	}
	if int(keyThres) != barrierConf.SecretThreshold {
		t.Fatalf("Unexpected threshold: got %d, expected %d", keyShares, barrierConf.SecretThreshold)
	}
	if method != expectedMethod {
		t.Fatalf("Unexpected method: got %d, expected init", expectedMethod)
	}

	return share
}
