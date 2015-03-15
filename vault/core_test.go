package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
)

var (
	// invalidKey is used to test Unseal
	invalidKey = []byte("abcdefghijklmnopqrstuvwxyz")[:17]
)

func TestCore_Init(t *testing.T) {
	inm := physical.NewInmem()
	conf := &CoreConfig{Physical: inm}
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
	c := TestCore(t)
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

func TestCore_Unseal_MultiShare(t *testing.T) {
	c := TestCore(t)

	_, err := c.Unseal(invalidKey)
	if err != ErrNotInit {
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

	sealed, err := c.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !sealed {
		t.Fatalf("should be sealed")
	}

	if prog := c.SecretProgress(); prog != 0 {
		t.Fatalf("bad progress: %d", prog)
	}

	for i := 0; i < 5; i++ {
		unseal, err := c.Unseal(res.SecretShares[i])
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Ignore redundant
		_, err = c.Unseal(res.SecretShares[i])
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i >= 2 {
			if !unseal {
				t.Fatalf("should be unsealed")
			}
			if prog := c.SecretProgress(); prog != 0 {
				t.Fatalf("bad progress: %d", prog)
			}
		} else {
			if unseal {
				t.Fatalf("should not be unsealed")
			}
			if prog := c.SecretProgress(); prog != i+1 {
				t.Fatalf("bad progress: %d", prog)
			}
		}
	}

	sealed, err = c.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if sealed {
		t.Fatalf("should not be sealed")
	}

	err = c.Seal()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ignore redundant
	err = c.Seal()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	sealed, err = c.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !sealed {
		t.Fatalf("should be sealed")
	}
}

func TestCore_Unseal_Single(t *testing.T) {
	c := TestCore(t)

	_, err := c.Unseal(invalidKey)
	if err != ErrNotInit {
		t.Fatalf("err: %v", err)
	}

	sealConf := &SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}
	res, err := c.Initialize(sealConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	sealed, err := c.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !sealed {
		t.Fatalf("should be sealed")
	}

	if prog := c.SecretProgress(); prog != 0 {
		t.Fatalf("bad progress: %d", prog)
	}

	unseal, err := c.Unseal(res.SecretShares[0])
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !unseal {
		t.Fatalf("should be unsealed")
	}
	if prog := c.SecretProgress(); prog != 0 {
		t.Fatalf("bad progress: %d", prog)
	}

	sealed, err = c.Sealed()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if sealed {
		t.Fatalf("should not be sealed")
	}
}

func TestCore_Route_Sealed(t *testing.T) {
	c := TestCore(t)
	sealConf := &SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}

	// Should not route anything
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "sys/mounts",
	}
	_, err := c.HandleRequest(req)
	if err != ErrSealed {
		t.Fatalf("err: %v", err)
	}

	res, err := c.Initialize(sealConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	unseal, err := c.Unseal(res.SecretShares[0])
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Should not error after unseal
	_, err = c.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

// Attempt to unseal after doing a first seal
func TestCore_SealUnseal(t *testing.T) {
	c, key := TestCoreUnsealed(t)
	if err := c.Seal(); err != nil {
		t.Fatalf("err: %v", err)
	}
	if unseal, err := c.Unseal(key); err != nil || !unseal {
		t.Fatalf("err: %v", err)
	}
}
