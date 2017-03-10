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
	//otplib "github.com/pquerna/otp/"
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
	key, _ := createKey()

	roleData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		Backend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateRole(t, "test", roleData, false),
			testAccStepReadCreds(t, b, config.StorageView, "test"),
		},
	})
}

func TestBackend_roleCrudDefaultValues(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(config)
	if err != nil {
		t.Fatal(err)
	}

	key, _ := createKey()

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
			testAccStepReadRole(t, "test", nil),
		},
	})
}

func testAccStepCreateRole(t *testing.T, name string, roleData map[string]interface{}, expectFail bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      path.Join("roles", name),
		Data:      roleData,
		ErrorOk:   expectFail,
	}
}

func testAccStepDeleteRole(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      path.Join("roles", name),
	}
}

func testAccStepReadCreds(t *testing.T, b logical.Backend, s logical.Storage, name string) logicaltest.TestStep {
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

			role, err := (backend) b.Role(s, name)

			if err != nil {
				t.Fatalf("Error retrieving role.")
			}

			if role == nil {
				t.Fatalf("Retrieved role is nil.")
			}

			// Translate digits and algorithm to a format the totp library understands
			var digits otplib.Digits
			switch role.Digits {
			case 6:
				digits = otplib.DigitsSix
			case 8:
				digits = otplib.DigitsEight
			}

			var algorithm otplib.Algorithm
			switch role.Algorithm {
			case "SHA1":
				algorithm = otplib.AlgorithmSHA1
			case "SHA256":
				algorithm = otplib.AlgorithmSHA256
			case "SHA512":
				algorithm = otplib.AlgorithmSHA512
			case "MD5":
				algorithm = otplib.AlgorithmMD5
			default:
				algorithm = otplib.AlgorithmSHA1
			}

			period := uint(role.Period)

			valid := totplib.ValidateCustom(d.Token, role.Key, time.Now().UTC(), totplib.ValidateOpts{
				Period:    period,
				Skew:      1,
				Digits:    digits,
				Algorithm: algorithm,
			})

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
			case d.Issuer != expected["issuer"]:
				return fmt.Errorf("Issuer should equal: %s", expected["issuer"])
			case d.Account_Name != expected["account_name"]:
				return fmt.Errorf("Account_Name should equal: %s", expected["account_name"])
			case d.Period != expected["period"]:
				return fmt.Errorf("Period should equal: %i", expected["period"])
			case d.Algorithm != expected["algorithm"]:
				return fmt.Errorf("Algorithm should equal: %s", expected["algorithm"])
			case d.Digits != expected["digits"]:
				return fmt.Errorf("Digits should equal: %i", expected["digits"])
			}

			return nil
		},
	}
}
