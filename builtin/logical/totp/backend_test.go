package totp

import (
	"fmt"
	"log"
	"path"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
	otplib "github.com/pquerna/otp"
	totplib "github.com/pquerna/otp/totp"
)

func createKey() (string, error) {
	keyUrl, err := totplib.Generate(totplib.GenerateOpts{
		Issuer:      "Vault",
		AccountName: "Test",
	})

	key := keyUrl.Secret()

	return key, err
}

func TestBackend_readCredentialsDefaultValues(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a new shared key
	key, _ := createKey()

	keyData := map[string]interface{}{
		"key": key,
	}

	expected := map[string]interface{}{
		"issuer":       "",
		"account_name": "",
		"digits":       otplib.DigitsSix,
		"period":       30,
		"algorithm":    otplib.AlgorithmSHA1,
		"key":          key,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepReadCreds(t, b, config.StorageView, "test", expected),
		},
	})
}

func TestBackend_readCredentialsEightDigitsThirtySecondPeriod(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a new shared key
	key, _ := createKey()

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
		"digits":       8,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"digits":       otplib.DigitsEight,
		"period":       30,
		"algorithm":    otplib.AlgorithmSHA1,
		"key":          key,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepReadCreds(t, b, config.StorageView, "test", expected),
		},
	})
}

func TestBackend_readCredentialsSixDigitsNinetySecondPeriod(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a new shared key
	key, _ := createKey()

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
		"period":       90,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"digits":       otplib.DigitsSix,
		"period":       90,
		"algorithm":    otplib.AlgorithmSHA1,
		"key":          key,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepReadCreds(t, b, config.StorageView, "test", expected),
		},
	})
}

func TestBackend_readCredentialsSHA256(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a new shared key
	key, _ := createKey()

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
		"algorithm":    "SHA256",
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"digits":       otplib.DigitsSix,
		"period":       30,
		"algorithm":    otplib.AlgorithmSHA256,
		"key":          key,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepReadCreds(t, b, config.StorageView, "test", expected),
		},
	})
}

func TestBackend_readCredentialsSHA512(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a new shared key
	key, _ := createKey()

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
		"algorithm":    "SHA512",
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"digits":       otplib.DigitsSix,
		"period":       30,
		"algorithm":    otplib.AlgorithmSHA512,
		"key":          key,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepReadCreds(t, b, config.StorageView, "test", expected),
		},
	})
}

func TestBackend_keyCrudDefaultValues(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	key, _ := createKey()

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"digits":       otplib.DigitsSix,
		"period":       30,
		"algorithm":    otplib.AlgorithmSHA1,
		"key":          key,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepDeleteKey(t, "test"),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_createKeyMissingKeyValue(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_createKeyInvalidKeyValue(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          "1",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_createKeyInvalidAlgorithm(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a new shared key
	key, _ := createKey()

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
		"algorithm":    "BADALGORITHM",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_createKeyInvalidPeriod(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a new shared key
	key, _ := createKey()

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
		"period":       -1,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_createKeyInvalidDigits(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a new shared key
	key, _ := createKey()

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
		"digits":       20,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func testAccStepCreateKey(t *testing.T, name string, keyData map[string]interface{}, expectFail bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      path.Join("keys", name),
		Data:      keyData,
		ErrorOk:   expectFail,
	}
}

func testAccStepDeleteKey(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      path.Join("keys", name),
	}
}

func testAccStepReadCreds(t *testing.T, b logical.Backend, s logical.Storage, name string, validation map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      path.Join("code", name),
		Check: func(resp *logical.Response) error {
			var d struct {
				Code string `mapstructure:"code"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[TRACE] Generated credentials: %v", d)

			period := validation["period"].(int)
			key := validation["key"].(string)
			algorithm := validation["algorithm"].(otplib.Algorithm)
			digits := validation["digits"].(otplib.Digits)

			valid, _ := totplib.ValidateCustom(d.Code, key, time.Now(), totplib.ValidateOpts{
				Period:    uint(period),
				Skew:      1,
				Digits:    digits,
				Algorithm: algorithm,
			})

			if !valid {
				t.Fatalf("generated code isn't valid")
			}

			return nil
		},
	}
}

func testAccStepReadKey(t *testing.T, name string, expected map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "keys/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if expected == nil {
					return nil
				}
				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				Issuer      string        `mapstructure:"issuer"`
				AccountName string        `mapstructure:"account_name"`
				Period      uint          `mapstructure:"period"`
				Algorithm   string        `mapstructure:"algorithm"`
				Digits      otplib.Digits `mapstructure:"digits"`
			}

			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			var keyAlgorithm otplib.Algorithm
			switch d.Algorithm {
			case "SHA1":
				keyAlgorithm = otplib.AlgorithmSHA1
			case "SHA256":
				keyAlgorithm = otplib.AlgorithmSHA256
			case "SHA512":
				keyAlgorithm = otplib.AlgorithmSHA512
			}

			period := expected["period"].(int)

			switch {
			case d.Issuer != expected["issuer"]:
				return fmt.Errorf("issuer should equal: %s", expected["issuer"])
			case d.AccountName != expected["account_name"]:
				return fmt.Errorf("ccount_Name should equal: %s", expected["account_name"])
			case d.Period != uint(period):
				return fmt.Errorf("period should equal: %i", expected["period"])
			case keyAlgorithm != expected["algorithm"]:
				return fmt.Errorf("algorithm should equal: %s", expected["algorithm"])
			case d.Digits != expected["digits"]:
				return fmt.Errorf("digits should equal: %i", expected["digits"])
			}

			return nil
		},
	}
}
