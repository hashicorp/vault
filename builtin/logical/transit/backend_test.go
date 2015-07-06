package transit

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
	"github.com/mitchellh/mapstructure"
)

const (
	testPlaintext = "the quick brown fox"
)

func TestBackend_basic(t *testing.T) {
	decryptData := make(map[string]interface{})
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", false),
			testAccStepReadPolicy(t, "test", false, false),
			testAccStepReadRaw(t, "test", false, false),
			testAccStepEncrypt(t, "test", testPlaintext, decryptData),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadPolicy(t, "test", true, false),
			testAccStepReadRaw(t, "test", true, false),
		},
	})
}

func TestBackend_upsert(t *testing.T) {
	decryptData := make(map[string]interface{})
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepReadPolicy(t, "test", true, false),
			testAccStepEncrypt(t, "test", testPlaintext, decryptData),
			testAccStepReadPolicy(t, "test", false, false),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadPolicy(t, "test", true, false),
		},
	})
}

func TestBackend_basic_derived(t *testing.T) {
	decryptData := make(map[string]interface{})
	logicaltest.Test(t, logicaltest.TestCase{
		Backend: Backend(),
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", true),
			testAccStepReadPolicy(t, "test", false, true),
			testAccStepReadRaw(t, "test", false, true),
			testAccStepEncryptContext(t, "test", testPlaintext, "my-cool-context", decryptData),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadPolicy(t, "test", true, true),
			testAccStepReadRaw(t, "test", true, true),
		},
	})
}

func testAccStepWritePolicy(t *testing.T, name string, derived bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "keys/" + name,
		Data: map[string]interface{}{
			"derived": derived,
		},
	}
}

func testAccStepDeletePolicy(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "keys/" + name,
	}
}

func testAccStepReadPolicy(t *testing.T, name string, expectNone, derived bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "keys/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil && !expectNone {
				return fmt.Errorf("missing response")
			} else if expectNone {
				if resp != nil {
					return fmt.Errorf("response when expecting none")
				}
				return nil
			}
			var d struct {
				Name       string `mapstructure:"name"`
				Key        []byte `mapstructure:"key"`
				CipherMode string `mapstructure:"cipher_mode"`
				Derived    bool   `mapstructure:"derived"`
				KDFMode    string `mapstructure:"kdf_mode"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Name != name {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.CipherMode != "aes-gcm" {
				return fmt.Errorf("bad: %#v", d)
			}
			// Should NOT get a key back
			if d.Key != nil {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.Derived != derived {
				return fmt.Errorf("bad: %#v", d)
			}
			if derived && d.KDFMode != kdfMode {
				return fmt.Errorf("bad: %#v", d)
			}
			return nil
		},
	}
}

func testAccStepReadRaw(t *testing.T, name string, expectNone, derived bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "raw/" + name,
		Check: func(resp *logical.Response) error {
			if resp == nil && !expectNone {
				return fmt.Errorf("missing response")
			} else if expectNone {
				if resp != nil {
					return fmt.Errorf("response when expecting none")
				}
				return nil
			}
			var d struct {
				Name       string `mapstructure:"name"`
				Key        []byte `mapstructure:"key"`
				CipherMode string `mapstructure:"cipher_mode"`
				Derived    bool   `mapstructure:"derived"`
				KDFMode    string `mapstructure:"kdf_mode"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Name != name {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.CipherMode != "aes-gcm" {
				return fmt.Errorf("bad: %#v", d)
			}
			if len(d.Key) != 32 {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.Derived != derived {
				return fmt.Errorf("bad: %#v", d)
			}
			if derived && d.KDFMode != kdfMode {
				return fmt.Errorf("bad: %#v", d)
			}
			return nil
		},
	}
}

func testAccStepEncrypt(
	t *testing.T, name, plaintext string, decryptData map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "encrypt/" + name,
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
		},
		Check: func(resp *logical.Response) error {
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
			return nil
		},
	}
}

func testAccStepEncryptContext(
	t *testing.T, name, plaintext, context string, decryptData map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "encrypt/" + name,
		Data: map[string]interface{}{
			"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
			"context":   base64.StdEncoding.EncodeToString([]byte(context)),
		},
		Check: func(resp *logical.Response) error {
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

func testAccStepDecrypt(
	t *testing.T, name, plaintext string, decryptData map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "decrypt/" + name,
		Data:      decryptData,
		Check: func(resp *logical.Response) error {
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
				return fmt.Errorf("plaintext mismatch: %s expect: %s", plainRaw, plaintext)
			}
			return nil
		},
	}
}
