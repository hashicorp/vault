package transit

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/testing/stepwise"
	dockerDriver "github.com/hashicorp/vault/sdk/testing/stepwise/drivers/docker"
	"github.com/mitchellh/mapstructure"
	"github.com/y0ssar1an/q"
)

// TestBackend_basic_derived_docker is an example test using the Docker Driver
func TestBackend_basic_derived_docker(t *testing.T) {
	decryptData := make(map[string]interface{})
	stepwise.Run(t, stepwise.Case{
		Driver: dockerDriver.NewDockerDriver("transit"),
		Steps: []stepwise.Step{
			testAccStepwiseListPolicy(t, "test", true),
			testAccStepwiseWritePolicy(t, "test", true),
			testAccStepwiseListPolicy(t, "test", false),
			testAccStepwiseReadPolicy(t, "test", false, true),
			testAccStepwiseEncryptContext(t, "test", testPlaintext, "my-cool-context", decryptData),
			testAccStepwiseDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepwiseEnableDeletion(t, "test"),
			testAccStepwiseDeletePolicy(t, "test"),
			testAccStepwiseReadPolicy(t, "test", true, true),
		},
	})
}

func testAccStepwiseWritePolicy(t *testing.T, name string, derived bool) stepwise.Step {
	ts := stepwise.Step{
		// Operation: logical.UpdateOperation,
		Operation: stepwise.WriteOperation,
		Path:      "keys/" + name,
		Data: map[string]interface{}{
			"derived": derived,
		},
		Check: func(resp *api.Secret) error {
			q.Q("--> stepwise write policy check func")
			return nil
		},
	}
	if os.Getenv("TRANSIT_ACC_KEY_TYPE") == "CHACHA" {
		ts.Data["type"] = "chacha20-poly1305"
	}
	return ts
}

func testAccStepwiseListPolicy(t *testing.T, name string, expectNone bool) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.ListOperation,
		Path:      "keys",
		Check: func(resp *api.Secret) error {
			q.Q("--> stepwise list check func")
			q.Q("resp in check:", resp)
			if (resp == nil || len(resp.Data) == 0) && !expectNone {
				return fmt.Errorf("missing response")
			}
			if expectNone && resp != nil {
				return fmt.Errorf("response data when expecting none")
			}

			if expectNone && resp == nil {
				return nil
			}

			q.Q("--> --> stepwise checking keys list")
			var d struct {
				Keys []string `mapstructure:"keys"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if len(d.Keys) > 0 && d.Keys[0] != name {
				return fmt.Errorf("bad name: %#v", d)
			}
			if len(d.Keys) != 1 {
				return fmt.Errorf("only 1 key expected, %d returned", len(d.Keys))
			}
			q.Q("--> --> stepwise does with check")
			return nil
		},
	}
}

func testAccStepwiseReadPolicy(t *testing.T, name string, expectNone, derived bool) stepwise.Step {
	return testAccStepwiseReadPolicyWithVersions(t, name, expectNone, derived, 1, 0)
}

func testAccStepwiseReadPolicyWithVersions(t *testing.T, name string, expectNone, derived bool, minDecryptionVersion int, minEncryptionVersion int) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.ReadOperation,
		Path:      "keys/" + name,
		Check: func(resp *api.Secret) error {
			q.Q("--> read policy check")
			if resp == nil && !expectNone {
				return fmt.Errorf("missing response")
			} else if expectNone {
				if resp != nil {
					return fmt.Errorf("response when expecting none")
				}
				return nil
			}
			var d struct {
				Name                 string           `mapstructure:"name"`
				Key                  []byte           `mapstructure:"key"`
				Keys                 map[string]int64 `mapstructure:"keys"`
				Type                 string           `mapstructure:"type"`
				Derived              bool             `mapstructure:"derived"`
				KDF                  string           `mapstructure:"kdf"`
				DeletionAllowed      bool             `mapstructure:"deletion_allowed"`
				ConvergentEncryption bool             `mapstructure:"convergent_encryption"`
				MinDecryptionVersion int              `mapstructure:"min_decryption_version"`
				MinEncryptionVersion int              `mapstructure:"min_encryption_version"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			q.Q("d after read:", d)

			if d.Name != name {
				return fmt.Errorf("bad name: %#v", d)
			}
			if os.Getenv("TRANSIT_ACC_KEY_TYPE") == "CHACHA" {
				if d.Type != keysutil.KeyType(keysutil.KeyType_ChaCha20_Poly1305).String() {
					return fmt.Errorf("bad key type: %#v", d)
				}
			} else if d.Type != keysutil.KeyType(keysutil.KeyType_AES256_GCM96).String() {
				return fmt.Errorf("bad key type: %#v", d)
			}
			// Should NOT get a key back
			if d.Key != nil {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.Keys == nil {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.MinDecryptionVersion != minDecryptionVersion {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.MinEncryptionVersion != minEncryptionVersion {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.DeletionAllowed == true {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.Derived != derived {
				return fmt.Errorf("bad: %#v", d)
			}
			if derived && d.KDF != "hkdf_sha256" {
				return fmt.Errorf("bad: %#v", d)
			}
			q.Q("--> read policy check OK")
			return nil
		},
	}
}

func testAccStepwiseEncryptContext(
	t *testing.T, name, plaintext, context string, decryptData map[string]interface{}) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "encrypt/" + name,
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
			"context":   base64.StdEncoding.EncodeToString([]byte(context)),
		},
		Check: func(resp *api.Secret) error {
			var d struct {
				Ciphertext string `mapstructure:"ciphertext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if d.Ciphertext == "" {
				return fmt.Errorf("missing ciphertext")
			}
			decryptData["ciphertext"] = d.Ciphertext
			decryptData["context"] = base64.StdEncoding.EncodeToString([]byte(context))
			return nil
		},
	}
}

func testAccStepwiseDecrypt(
	t *testing.T, name, plaintext string, decryptData map[string]interface{}) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "decrypt/" + name,
		Data:      decryptData,
		Check: func(resp *api.Secret) error {
			var d struct {
				Plaintext string `mapstructure:"plaintext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			// Decode the base64
			plainRaw, err := base64.StdEncoding.DecodeString(d.Plaintext)
			if err != nil {
				return err
			}

			if string(plainRaw) != plaintext {
				return fmt.Errorf("plaintext mismatch: %s expect: %s, decryptData was %#v", plainRaw, plaintext, decryptData)
			}
			return nil
		},
	}
}

func testAccStepwiseEnableDeletion(t *testing.T, name string) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "keys/" + name + "/config",
		Data: map[string]interface{}{
			"deletion_allowed": true,
		},
	}
}

func testAccStepwiseDeletePolicy(t *testing.T, name string) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.DeleteOperation,
		Path:      "keys/" + name,
	}
}
