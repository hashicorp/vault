package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/physical"
)

func TestCore_Init(t *testing.T) {
	inm := physical.NewInmem()
	conf := &CoreConfig{physical: inm}
	c, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	init, err := c.Initialized()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if init {
		t.Fatalf("should not be init")
	}

	// Check the seal configuration
	outConf, err := c.SealConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if outConf != nil {
		t.Fatalf("bad: %v", outConf)
	}

	sealConf := &SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}
	res, err := c.Initialize(sealConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(res.SecretShares) != 1 {
		t.Fatalf("Bad: %v", res)
	}

	_, err = c.Initialize(sealConf)
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
	outConf, err = c.SealConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(outConf, sealConf) {
		t.Fatalf("bad: %v expect: %v", outConf, sealConf)
	}

	// New Core, same backend
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	_, err = c2.Initialize(sealConf)
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
	outConf, err = c2.SealConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(outConf, sealConf) {
		t.Fatalf("bad: %v expect: %v", outConf, sealConf)
	}
}

func TestCore_Init_MultiShare(t *testing.T) {
	inm := physical.NewInmem()
	conf := &CoreConfig{physical: inm}
	c, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	sealConf := &SealConfig{
		SecretShares:    5,
		SecretThreshold: 3,
	}
	res, err := c.Initialize(sealConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(res.SecretShares) != 5 {
		t.Fatalf("Bad: %v", res)
	}

	// Check the seal configuration
	outConf, err := c.SealConfig()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(outConf, sealConf) {
		t.Fatalf("bad: %v expect: %v", outConf, sealConf)
	}
}
