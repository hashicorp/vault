package totp

import (
	"fmt"
	"log"
	"path"
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
	totplib "github.com/pquerna/otp/totp"
)

/*
	Test each algorithm type
	Test digits
	Test periods
	Test defaults
	Test invalid period (negative)
	Test invalid key
	Test invalid account_name
	Test invalid issuer
*/

func createKey() (string, error) {
	keyUrl, err := totplib.Generate(totplib.GenerateOpts{
		Issuer:      "Vault",
		AccountName: "Test",
	})

	key := keyUrl.Secret()

	return key, err
}

func TestBackend_basic(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a new shared key
	key := createKey()

	roleData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateRole(t, "test", roleData, false),
			testAccStepReadCreds(t, b, config.StorageView, "test", key),
		},
	})
}

func TestBackend_roleCrud(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	key := createKey()

	roleData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"digits":       6,
		"period":       30,
		"algorithm":    "SHA1",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateRole(t, "test", roleData, false),
			testAccStepReadRole(t, "test", expected),
			testAccStepDeleteRole(t, "test"),
			testAccStepReadRole(t, "test", ""),
		},
	})
}

func testAccStepCreateRole(t *testing.T, name string, data map[string]interface{}, expectFail bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      path.Join("roles", name),
		Data:      data,
		ErrorOk:   expectFail,
	}
}

func testAccStepDeleteRole(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      path.Join("roles", name),
	}
}

func testAccStepReadCreds(t *testing.T, b logical.Backend, s logical.Storage, name string, key string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      path.Join("creds", name),
		Check: func(resp *logical.Response) error {
			var d struct {
				Token string `mapstructure:"token"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[TRACE] Generated credentials: %v", d)

			valid := totplib.Validate(d.Token, key)

			if !valid {
				t.Fatalf("Generated token isn't valid.")
			}

			return nil
		},
	}
}

func testAccStepReadRole(t *testing.T, name string, expected map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				if expected == nil {
					return nil
				}
				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				Issuer       string `mapstructure:"issuer"`
				Account_Name string `mapstructure:"account_name"`
				Period       int    `mapstructure:"period"`
				Algorithm    string `mapstructure:"algorithm"`
				Digits       int    `mapstructure:"digits"`
			}

			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			switch {
			case d.Issuer != expected.Get("issuer"):
				return fmt.Errorf("Issuer should equal: %s", expected.Get("issuer"))
			case d.Account_Name != expected.Get("account_name"):
				return fmt.Errorf("Account_Name should equal: %s", expected.Get("account_name"))
			case d.Period != expected.Get("period"):
				return fmt.Errorf("Period should equal: %i", expected.Get("period"))
			case d.Algorithm != expected.Get("algorithm"):
				return fmt.Errorf("Algorithm should equal: %s", expected.Get("algorithm"))
			case d.Digits != expected.Get("digits"):
				return fmt.Errorf("Digits should equal: %i", expected.Get("digits"))
			}

			return nil
		},
	}
}
