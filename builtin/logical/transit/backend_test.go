package transit

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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

func TestBackend_upsert(t *testing.T) {
	decryptData := make(map[string]interface{})
	logicaltest.Test(t, logicaltest.TestCase{
		Factory: Factory,
		Steps: []logicaltest.TestStep{
			testAccStepReadPolicy(t, "test", true, false),
			testAccStepEncryptUpsert(t, "test", testPlaintext, decryptData),
			testAccStepReadPolicy(t, "test", false, false),
			testAccStepDecrypt(t, "test", testPlaintext, decryptData),
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
				Name                 string           `mapstructure:"name"`
				Key                  []byte           `mapstructure:"key"`
				Keys                 map[string]int64 `mapstructure:"keys"`
				CipherMode           string           `mapstructure:"cipher_mode"`
				Derived              bool             `mapstructure:"derived"`
				KDFMode              string           `mapstructure:"kdf_mode"`
				DeletionAllowed      bool             `mapstructure:"deletion_allowed"`
				ConvergentEncryption bool             `mapstructure:"convergent_encryption"`
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

func testAccStepEncryptUpsert(
	t *testing.T, name, plaintext string, decryptData map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.CreateOperation,
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
				return fmt.Errorf("plaintext mismatch: got '%s', expected '%s', decryptData was %#v", d.Plaintext, dataKeyInfo["plaintext"].(string), resp.Data)
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

func TestConvergentEncryption(t *testing.T) {
	var b *backend
	sysView := logical.TestSystemView()
	storage := &logical.InmemStorage{}

	b = Backend(&logical.BackendConfig{
		StorageView: storage,
		System:      sysView,
	})

	req := &logical.Request{
		Storage:   storage,
		Operation: logical.UpdateOperation,
		Path:      "keys/testkeynonderived",
		Data: map[string]interface{}{
			"derived":               false,
			"convergent_encryption": true,
		},
	}

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if !resp.IsError() {
		t.Fatalf("bad: expected error response, got %#v", *resp)
	}

	req.Path = "keys/testkey"
	req.Data = map[string]interface{}{
		"derived":               true,
		"convergent_encryption": true,
	}

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatalf("bad: got resp %#v", *resp)
	}

	// First, test using an invalid length of nonce
	req.Path = "encrypt/testkey"
	req.Data = map[string]interface{}{
		"plaintext": "emlwIHphcA==", // "zip zap"
		"context":   "Zm9vIGJhcg==", // "foo bar"
	}
	resp, err = b.HandleRequest(req)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if !resp.IsError() {
		t.Fatalf("expected error response, got %#v", *resp)
	}

	// Now test encrypting the same value twice
	req.Data = map[string]interface{}{
		"plaintext": "emlwIHphcA==",     // "zip zap"
		"context":   "b25ldHdvdGhyZWVl", // "onetwothreee"
	}
	resp, err = b.HandleRequest(req)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext1 := resp.Data["ciphertext"].(string)

	resp, err = b.HandleRequest(req)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext2 := resp.Data["ciphertext"].(string)

	if ciphertext1 != ciphertext2 {
		t.Fatalf("expected the same ciphertext but got %s and %s", ciphertext1, ciphertext2)
	}

	// For sanity, also check a different value
	req.Data = map[string]interface{}{
		"plaintext": "emlwIHphcA==",     // "zip zap"
		"context":   "dHdvdGhyZWVmb3Vy", // "twothreefour"
	}
	resp, err = b.HandleRequest(req)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext3 := resp.Data["ciphertext"].(string)

	resp, err = b.HandleRequest(req)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.IsError() {
		t.Fatalf("got error response: %#v", *resp)
	}
	ciphertext4 := resp.Data["ciphertext"].(string)

	if ciphertext3 != ciphertext4 {
		t.Fatalf("expected the same ciphertext but got %s and %s", ciphertext1, ciphertext2)
	}
	if ciphertext1 == ciphertext3 {
		t.Fatalf("expected different ciphertexts")
	}

}

func TestPolicyFuzzing(t *testing.T) {
	var be *backend
	sysView := logical.TestSystemView()

	be = Backend(&logical.BackendConfig{
		System: sysView,
	})
	testPolicyFuzzingCommon(t, be)

	sysView.CachingDisabledVal = true
	be = Backend(&logical.BackendConfig{
		System: sysView,
	})
	testPolicyFuzzingCommon(t, be)
}

func testPolicyFuzzingCommon(t *testing.T, be *backend) {
	storage := &logical.InmemStorage{}
	wg := sync.WaitGroup{}

	funcs := []string{"encrypt", "decrypt", "rotate", "change_min_version"}
	//keys := []string{"test1", "test2", "test3", "test4", "test5"}
	keys := []string{"test1", "test2", "test3"}

	// This is the goroutine loop
	doFuzzy := func(id int) {
		// Check for panics, otherwise notify we're done
		defer func() {
			if err := recover(); err != nil {
				t.Fatalf("got a panic: %v", err)
			}
			wg.Done()
		}()

		// Holds the latest encrypted value for each key
		latestEncryptedText := map[string]string{}

		startTime := time.Now()
		req := &logical.Request{
			Storage: storage,
			Data:    map[string]interface{}{},
		}
		fd := &framework.FieldData{}

		var chosenFunc, chosenKey string

		//t.Errorf("Starting %d", id)
		for {
			// Stop after 10 seconds
			if time.Now().Sub(startTime) > 10*time.Second {
				return
			}

			// Pick a function and a key
			chosenFunc = funcs[rand.Int()%len(funcs)]
			chosenKey = keys[rand.Int()%len(keys)]

			fd.Raw = map[string]interface{}{
				"name": chosenKey,
			}
			fd.Schema = be.pathKeys().Fields

			// Try to write the key to make sure it exists
			_, err := be.pathPolicyWrite(req, fd)
			if err != nil {
				t.Fatalf("got an error: %v", err)
			}

			switch chosenFunc {
			// Encrypt our plaintext and store the result
			case "encrypt":
				//t.Errorf("%s, %s, %d", chosenFunc, chosenKey, id)
				fd.Raw["plaintext"] = base64.StdEncoding.EncodeToString([]byte(testPlaintext))
				fd.Schema = be.pathEncrypt().Fields
				resp, err := be.pathEncryptWrite(req, fd)
				if err != nil {
					t.Fatalf("got an error: %v, resp is %#v", err, *resp)
				}
				latestEncryptedText[chosenKey] = resp.Data["ciphertext"].(string)

			// Rotate to a new key version
			case "rotate":
				//t.Errorf("%s, %s, %d", chosenFunc, chosenKey, id)
				fd.Schema = be.pathRotate().Fields
				resp, err := be.pathRotateWrite(req, fd)
				if err != nil {
					t.Fatalf("got an error: %v, resp is %#v, chosenKey is %s", err, *resp, chosenKey)
				}

			// Decrypt the ciphertext and compare the result
			case "decrypt":
				//t.Errorf("%s, %s, %d", chosenFunc, chosenKey, id)
				ct := latestEncryptedText[chosenKey]
				if ct == "" {
					continue
				}

				fd.Raw["ciphertext"] = ct
				fd.Schema = be.pathDecrypt().Fields
				resp, err := be.pathDecryptWrite(req, fd)
				if err != nil {
					// This could well happen since the min version is jumping around
					if resp.Data["error"].(string) == ErrTooOld {
						continue
					}
					t.Fatalf("got an error: %v, resp is %#v, ciphertext was %s, chosenKey is %s, id is %d", err, *resp, ct, chosenKey, id)
				}
				ptb64 := resp.Data["plaintext"].(string)
				pt, err := base64.StdEncoding.DecodeString(ptb64)
				if err != nil {
					t.Fatalf("got an error decoding base64 plaintext: %v", err)
					return
				}
				if string(pt) != testPlaintext {
					t.Fatalf("got bad plaintext back: %s", pt)
				}

			// Change the min version, which also tests the archive functionality
			case "change_min_version":
				//t.Errorf("%s, %s, %d", chosenFunc, chosenKey, id)
				resp, err := be.pathPolicyRead(req, fd)
				if err != nil {
					t.Fatalf("got an error reading policy %s: %v", chosenKey, err)
				}
				latestVersion := resp.Data["latest_version"].(int)

				// keys start at version 1 so we want [1, latestVersion] not [0, latestVersion)
				setVersion := (rand.Int() % latestVersion) + 1
				fd.Raw["min_decryption_version"] = setVersion
				fd.Schema = be.pathConfig().Fields
				resp, err = be.pathConfigWrite(req, fd)
				if err != nil {
					t.Fatalf("got an error setting min decryption version: %v", err)
				}
			}
		}
	}

	// Spawn 1000 of these workers for 10 seconds
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go doFuzzy(i)
	}

	// Wait for them all to finish
	wg.Wait()
}
