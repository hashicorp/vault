package transit

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
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
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", false),
			testAccStepReadPolicy(t, "test", false, false),
			testAccStepEncrypt(t, "test", testPlaintext, decryptData),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepDeleteNotDisabledPolicy(t, "test"),
			testAccStepEnableDeletion(t, "test"),
			testAccStepDeletePolicy(t, "test"),
			testAccStepWritePolicy(t, "test", false),
			testAccStepEnableDeletion(t, "test"),
			testAccStepDisableDeletion(t, "test"),
			testAccStepDeleteNotDisabledPolicy(t, "test"),
			testAccStepEnableDeletion(t, "test"),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadPolicy(t, "test", true, false),
		},
	})
}

func TestBackend_datakey(t *testing.T) {
	dataKeyInfo := make(map[string]interface{})
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", false),
			testAccStepReadPolicy(t, "test", false, false),
			testAccStepWriteDatakey(t, "test", false, 256, dataKeyInfo),
			testAccStepDecryptDatakey(t, "test", dataKeyInfo),
			testAccStepWriteDatakey(t, "test", true, 128, dataKeyInfo),
		},
	})
}

func TestBackend_rotation(t *testing.T) {
	decryptData := make(map[string]interface{})
	encryptHistory := make(map[int]map[string]interface{})
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", false),
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 0, encryptHistory),
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 1, encryptHistory),
			testAccStepRotate(t, "test"), // now v2
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 2, encryptHistory),
			testAccStepRotate(t, "test"), // now v3
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 3, encryptHistory),
			testAccStepRotate(t, "test"), // now v4
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 4, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepEncryptVX(t, "test", testPlaintext, decryptData, 99, encryptHistory),
			testAccStepDecryptExpectFailure(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 0, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 1, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 2, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 3, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 99, encryptHistory),
			testAccStepDecryptExpectFailure(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 4, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepDeleteNotDisabledPolicy(t, "test"),
			testAccStepAdjustPolicy(t, "test", 3),
			testAccStepLoadVX(t, "test", decryptData, 0, encryptHistory),
			testAccStepDecryptExpectFailure(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 1, encryptHistory),
			testAccStepDecryptExpectFailure(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 2, encryptHistory),
			testAccStepDecryptExpectFailure(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 3, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 4, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepAdjustPolicy(t, "test", 1),
			testAccStepLoadVX(t, "test", decryptData, 0, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 1, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepLoadVX(t, "test", decryptData, 2, encryptHistory),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepRewrap(t, "test", decryptData, 4),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepEnableDeletion(t, "test"),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadPolicy(t, "test", true, false),
		},
	})
}

func TestBackend_basic_derived(t *testing.T) {
	decryptData := make(map[string]interface{})
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", true),
			testAccStepReadPolicy(t, "test", false, true),
			testAccStepEncryptContext(t, "test", testPlaintext, "my-cool-context", decryptData),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
			testAccStepEnableDeletion(t, "test"),
			testAccStepDeletePolicy(t, "test"),
			testAccStepReadPolicy(t, "test", true, true),
		},
	})
}

func testAccStepWritePolicy(t *testing.T, name string, derived bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name,
		Data: map[string]interface{}{
			"derived": derived,
		},
	}
}

func testAccStepAdjustPolicy(t *testing.T, name string, minVer int) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name + "/config",
		Data: map[string]interface{}{
			"min_decryption_version": minVer,
		},
	}
}

func testAccStepDisableDeletion(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name + "/config",
		Data: map[string]interface{}{
			"deletion_allowed": false,
		},
	}
}

func testAccStepEnableDeletion(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name + "/config",
		Data: map[string]interface{}{
			"deletion_allowed": true,
		},
	}
}

func testAccStepDeletePolicy(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "keys/" + name,
	}
}

func testAccStepDeleteNotDisabledPolicy(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "keys/" + name,
		ErrorOk:   true,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				return fmt.Errorf("Got nil response instead of error")
			}
			if resp.IsError() {
				return nil
			}
			return fmt.Errorf("expected error but did not get one")
		},
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
				Name            string           `mapstructure:"name"`
				Key             []byte           `mapstructure:"key"`
				Keys            map[string]int64 `mapstructure:"keys"`
				CipherMode      string           `mapstructure:"cipher_mode"`
				Derived         bool             `mapstructure:"derived"`
				KDFMode         string           `mapstructure:"kdf_mode"`
				DeletionAllowed bool             `mapstructure:"deletion_allowed"`
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
			if d.Keys == nil {
				return fmt.Errorf("bad: %#v", d)
			}
			if d.DeletionAllowed == true {
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
		Operation: logical.UpdateOperation,
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
		Operation: logical.UpdateOperation,
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
		Operation: logical.UpdateOperation,
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
				return fmt.Errorf("plaintext mismatch: %s expect: %s, decryptData was %#v", plainRaw, plaintext, decryptData)
			}
			return nil
		},
	}
}

func testAccStepRewrap(
	t *testing.T, name string, decryptData map[string]interface{}, expectedVer int) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "rewrap/" + name,
		Data:      decryptData,
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
			splitStrings := strings.Split(d.Ciphertext, ":")
			verString := splitStrings[1][1:]
			ver, err := strconv.Atoi(verString)
			if err != nil {
				return fmt.Errorf("Error pulling out version from verString '%s', ciphertext was %s", verString, d.Ciphertext)
			}
			if ver != expectedVer {
				return fmt.Errorf("Did not get expected version")
			}
			decryptData["ciphertext"] = d.Ciphertext
			return nil
		},
	}
}

func testAccStepEncryptVX(
	t *testing.T, name, plaintext string, decryptData map[string]interface{},
	ver int, encryptHistory map[int]map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
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
			splitStrings := strings.Split(d.Ciphertext, ":")
			splitStrings[1] = "v" + strconv.Itoa(ver)
			ciphertext := strings.Join(splitStrings, ":")
			decryptData["ciphertext"] = ciphertext
			encryptHistory[ver] = map[string]interface{}{
				"ciphertext": ciphertext,
			}
			return nil
		},
	}
}

func testAccStepLoadVX(
	t *testing.T, name string, decryptData map[string]interface{},
	ver int, encryptHistory map[int]map[string]interface{}) logicaltest.TestStep {
	// This is really a no-op to allow us to do data manip in the check function
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "keys/" + name,
		Check: func(resp *logical.Response) error {
			decryptData["ciphertext"] = encryptHistory[ver]["ciphertext"].(string)
			return nil
		},
	}
}

func testAccStepDecryptExpectFailure(
	t *testing.T, name, plaintext string, decryptData map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "decrypt/" + name,
		Data:      decryptData,
		ErrorOk:   true,
		Check: func(resp *logical.Response) error {
			if !resp.IsError() {
				return fmt.Errorf("expected error")
			}
			return nil
		},
	}
}

func testAccStepRotate(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "keys/" + name + "/rotate",
	}
}

func testAccStepWriteDatakey(t *testing.T, name string,
	noPlaintext bool, bits int,
	dataKeyInfo map[string]interface{}) logicaltest.TestStep {
	data := map[string]interface{}{}
	subPath := "plaintext"
	if noPlaintext {
		subPath = "wrapped"
	}
	if bits != 256 {
		data["bits"] = bits
	}
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "datakey/" + subPath + "/" + name,
		Data:      data,
		Check: func(resp *logical.Response) error {
			var d struct {
				Plaintext  string `mapstructure:"plaintext"`
				Ciphertext string `mapstructure:"ciphertext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			if noPlaintext && len(d.Plaintext) != 0 {
				return fmt.Errorf("received plaintxt when we disabled it")
			}
			if !noPlaintext {
				if len(d.Plaintext) == 0 {
					return fmt.Errorf("did not get plaintext when we expected it")
				}
				dataKeyInfo["plaintext"] = d.Plaintext
				plainBytes, err := base64.StdEncoding.DecodeString(d.Plaintext)
				if err != nil {
					return fmt.Errorf("could not base64 decode plaintext string '%s'", d.Plaintext)
				}
				if len(plainBytes)*8 != bits {
					return fmt.Errorf("returned key does not have correct bit length")
				}
			}
			dataKeyInfo["ciphertext"] = d.Ciphertext
			return nil
		},
	}
}

func testAccStepDecryptDatakey(t *testing.T, name string,
	dataKeyInfo map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "decrypt/" + name,
		Data:      dataKeyInfo,
		Check: func(resp *logical.Response) error {
			var d struct {
				Plaintext string `mapstructure:"plaintext"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			if d.Plaintext != dataKeyInfo["plaintext"].(string) {
				return fmt.Errorf("plaintext mismatch: got '%s', expected '%s', decryptData was %#v", d.Plaintext, dataKeyInfo["plaintext"].(string))
			}
			return nil
		},
	}
}

func TestKeyUpgrade(t *testing.T) {
	p := &Policy{
		Name:       "test",
		Key:        []byte(testPlaintext),
		CipherMode: "aes-gcm",
	}

	p.migrateKeyToKeysMap()

	if p.Key != nil ||
		p.Keys == nil ||
		len(p.Keys) != 1 ||
		string(p.Keys[1].Key) != testPlaintext {
		t.Errorf("bad key migration, result is %#v", p.Keys)
	}
}
