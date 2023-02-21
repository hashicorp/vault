package transit

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	stepwise "github.com/hashicorp/vault-testing-stepwise"
	dockerEnvironment "github.com/hashicorp/vault-testing-stepwise/environments/docker"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/mitchellh/mapstructure"
)

// TestBackend_basic_docker is an example test using the Docker Environment
func TestAccBackend_basic_docker(t *testing.T) {
	decryptData := make(map[string]interface{})
	envOptions := stepwise.MountOptions{
		RegistryName:    "updatedtransit",
		PluginType:      api.PluginTypeSecrets,
		PluginName:      "transit",
		MountPathPrefix: "transit_temp",
	}
	stepwise.Run(t, stepwise.Case{
		Environment: dockerEnvironment.NewEnvironment("updatedtransit", &envOptions),
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
		Operation: stepwise.WriteOperation,
		Path:      "keys/" + name,
		Data: map[string]interface{}{
			"derived": derived,
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
		Assert: func(resp *api.Secret, err error) error {
			if (resp == nil || len(resp.Data) == 0) && !expectNone {
				return fmt.Errorf("missing response")
			}
			if expectNone && resp != nil {
				return fmt.Errorf("response data when expecting none")
			}

			if expectNone && resp == nil {
				return nil
			}

			var d struct {
				Keys []string `mapstructure:"keys"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if len(d.Keys) == 0 {
				return fmt.Errorf("missing keys")
			}
			if len(d.Keys) > 1 {
				return fmt.Errorf("only 1 key expected, %d returned", len(d.Keys))
			}
			if d.Keys[0] != name {
				return fmt.Errorf("Actual key: %s\nExpected key: %s", d.Keys[0], name)
			}
			return nil
		},
	}
}

func testAccStepwiseReadPolicy(t *testing.T, name string, expectNone, derived bool) stepwise.Step {
	t.Helper()
	return testAccStepwiseReadPolicyWithVersions(t, name, expectNone, derived, 1, 0)
}

func testAccStepwiseReadPolicyWithVersions(t *testing.T, name string, expectNone, derived bool, minDecryptionVersion int, minEncryptionVersion int) stepwise.Step {
	t.Helper()
	return stepwise.Step{
		Operation: stepwise.ReadOperation,
		Path:      "keys/" + name,
		Assert: func(resp *api.Secret, err error) error {
			t.Helper()
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
				return fmt.Errorf("unexpected key found")
			}
			if d.Keys == nil {
				return fmt.Errorf("no keys found")
			}
			if d.MinDecryptionVersion != minDecryptionVersion {
				return fmt.Errorf("minimum decryption version mismatch, expected (%#v), found (%#v)", minEncryptionVersion, d.MinDecryptionVersion)
			}
			if d.MinEncryptionVersion != minEncryptionVersion {
				return fmt.Errorf("minimum encryption version mismatch, expected (%#v), found (%#v)", minEncryptionVersion, d.MinDecryptionVersion)
			}
			if d.DeletionAllowed {
				return fmt.Errorf("expected DeletionAllowed to be false, but got true")
			}
			if d.Derived != derived {
				return fmt.Errorf("derived mismatch, expected (%t), got (%t)", derived, d.Derived)
			}
			if derived && d.KDF != "hkdf_sha256" {
				return fmt.Errorf("expected KDF to be hkdf_sha256, but got (%s)", d.KDF)
			}
			return nil
		},
	}
}

func testAccStepwiseEncryptContext(
	t *testing.T, name, plaintext, context string, decryptData map[string]interface{},
) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "encrypt/" + name,
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
			"context":   base64.StdEncoding.EncodeToString([]byte(context)),
		},
		Assert: func(resp *api.Secret, err error) error {
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
	t *testing.T, name, plaintext string, decryptData map[string]interface{},
) stepwise.Step {
	return stepwise.Step{
		Operation: stepwise.UpdateOperation,
		Path:      "decrypt/" + name,
		Data:      decryptData,
		Assert: func(resp *api.Secret, err error) error {
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
