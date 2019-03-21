package totp

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
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

	return strings.ToLower(key), err
}

func generateCode(key string, period uint, digits otplib.Digits, algorithm otplib.Algorithm) (string, error) {
	// Generate password using totp library
	totpToken, err := totplib.GenerateCodeCustom(key, time.Now(), totplib.ValidateOpts{
		Period:    period,
		Digits:    digits,
		Algorithm: algorithm,
	})

	return totpToken, err
}

func TestBackend_KeyName(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name    string
		KeyName string
		Fail    bool
	}{
		{
			"without @",
			"sample",
			false,
		},
		{
			"with @ in the beginning",
			"@sample.com",
			true,
		},
		{
			"with @ in the end",
			"sample.com@",
			true,
		},
		{
			"with @ in between",
			"sample@sample.com",
			false,
		},
		{
			"with multiple @",
			"sample@sample@@sample.com",
			false,
		},
	}
	var resp *logical.Response
	for _, tc := range tests {
		resp, err = b.HandleRequest(namespace.RootContext(nil), &logical.Request{
			Path:      "keys/" + tc.KeyName,
			Operation: logical.UpdateOperation,
			Storage:   config.StorageView,
			Data: map[string]interface{}{
				"generate":     true,
				"account_name": "vault",
				"issuer":       "hashicorp",
			},
		})
		if tc.Fail {
			if err == nil {
				t.Fatalf("expected an error for test %q", tc.Name)
			}
			continue
		} else if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: test name: %q\nresp: %#v\nerr: %v", tc.Name, resp, err)
		}
		resp, err = b.HandleRequest(namespace.RootContext(nil), &logical.Request{
			Path:      "code/" + tc.KeyName,
			Operation: logical.ReadOperation,
			Storage:   config.StorageView,
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad: test name: %q\nresp: %#v\nerr: %v", tc.Name, resp, err)
		}
		if resp.Data["code"].(string) == "" {
			t.Fatalf("failed to generate code for test %q", tc.Name)
		}
	}
}

func TestBackend_readCredentialsDefaultValues(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a new shared key
	key, _ := createKey()

	keyData := map[string]interface{}{
		"key":      key,
		"generate": false,
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
		LogicalBackend: b,
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
	b, err := Factory(context.Background(), config)
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
		"generate":     false,
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
		LogicalBackend: b,
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
	b, err := Factory(context.Background(), config)
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
		"generate":     false,
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
		LogicalBackend: b,
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
	b, err := Factory(context.Background(), config)
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
		"generate":     false,
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
		LogicalBackend: b,
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
	b, err := Factory(context.Background(), config)
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
		"generate":     false,
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
		LogicalBackend: b,
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
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	key, _ := createKey()

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          key,
		"generate":     false,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"digits":       otplib.DigitsSix,
		"period":       30,
		"algorithm":    otplib.AlgorithmSHA1,
		"key":          key,
	}

	code, _ := generateCode(key, 30, otplib.DigitsSix, otplib.AlgorithmSHA1)
	invalidCode := "12345678"

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepValidateCode(t, "test", code, true, false),
			// Next step should fail because it should be in the used cache
			testAccStepValidateCode(t, "test", code, false, true),
			testAccStepValidateCode(t, "test", invalidCode, false, false),
			testAccStepDeleteKey(t, "test"),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_createKeyMissingKeyValue(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"generate":     false,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_createKeyInvalidKeyValue(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key":          "1",
		"generate":     false,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_createKeyInvalidAlgorithm(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
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
		"generate":     false,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_createKeyInvalidPeriod(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
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
		"generate":     false,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_createKeyInvalidDigits(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
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
		"generate":     false,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_generatedKeyDefaultValues(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"generate":     true,
		"key_size":     20,
		"exported":     true,
		"qr_size":      200,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"digits":       otplib.DigitsSix,
		"period":       30,
		"algorithm":    otplib.AlgorithmSHA1,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
		},
	})
}

func TestBackend_generatedKeyDefaultValuesNoQR(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"generate":     true,
		"key_size":     20,
		"exported":     true,
		"qr_size":      0,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
		},
	})
}

func TestBackend_generatedKeyNonDefaultKeySize(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"generate":     true,
		"key_size":     10,
		"exported":     true,
		"qr_size":      200,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"digits":       otplib.DigitsSix,
		"period":       30,
		"algorithm":    otplib.AlgorithmSHA1,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
		},
	})
}

func TestBackend_urlPassedNonGeneratedKeyInvalidPeriod(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	urlString := "otpauth://totp/Vault:test@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&algorithm=SHA512&digits=6&period=AZ"

	keyData := map[string]interface{}{
		"url":      urlString,
		"generate": false,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_urlPassedNonGeneratedKeyInvalidDigits(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	urlString := "otpauth://totp/Vault:test@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&algorithm=SHA512&digits=Q&period=60"

	keyData := map[string]interface{}{
		"url":      urlString,
		"generate": false,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_urlPassedNonGeneratedKeyIssuerInFirstPosition(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	urlString := "otpauth://totp/Vault:test@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&algorithm=SHA512&digits=6&period=60"

	keyData := map[string]interface{}{
		"url":      urlString,
		"generate": false,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "test@email.com",
		"digits":       otplib.DigitsSix,
		"period":       60,
		"algorithm":    otplib.AlgorithmSHA512,
		"key":          "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepReadCreds(t, b, config.StorageView, "test", expected),
		},
	})
}

func TestBackend_urlPassedNonGeneratedKeyIssuerInQueryString(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	urlString := "otpauth://totp/test@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&algorithm=SHA512&digits=6&period=60&issuer=Vault"

	keyData := map[string]interface{}{
		"url":      urlString,
		"generate": false,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "test@email.com",
		"digits":       otplib.DigitsSix,
		"period":       60,
		"algorithm":    otplib.AlgorithmSHA512,
		"key":          "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepReadCreds(t, b, config.StorageView, "test", expected),
		},
	})
}

func TestBackend_urlPassedNonGeneratedKeyMissingIssuer(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	urlString := "otpauth://totp/test@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&algorithm=SHA512&digits=6&period=60"

	keyData := map[string]interface{}{
		"url":      urlString,
		"generate": false,
	}

	expected := map[string]interface{}{
		"issuer":       "",
		"account_name": "test@email.com",
		"digits":       otplib.DigitsSix,
		"period":       60,
		"algorithm":    otplib.AlgorithmSHA512,
		"key":          "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepReadCreds(t, b, config.StorageView, "test", expected),
		},
	})
}

func TestBackend_urlPassedNonGeneratedKeyMissingAccountName(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	urlString := "otpauth://totp/Vault:?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&algorithm=SHA512&digits=6&period=60"

	keyData := map[string]interface{}{
		"url":      urlString,
		"generate": false,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "",
		"digits":       otplib.DigitsSix,
		"period":       60,
		"algorithm":    otplib.AlgorithmSHA512,
		"key":          "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepReadCreds(t, b, config.StorageView, "test", expected),
		},
	})
}

func TestBackend_urlPassedNonGeneratedKeyMissingAccountNameandIssuer(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	urlString := "otpauth://totp/?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&algorithm=SHA512&digits=6&period=60"

	keyData := map[string]interface{}{
		"url":      urlString,
		"generate": false,
	}

	expected := map[string]interface{}{
		"issuer":       "",
		"account_name": "",
		"digits":       otplib.DigitsSix,
		"period":       60,
		"algorithm":    otplib.AlgorithmSHA512,
		"key":          "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
			testAccStepReadCreds(t, b, config.StorageView, "test", expected),
		},
	})
}

func TestBackend_generatedKeyInvalidSkew(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"skew":         "2",
		"generate":     true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_generatedKeyInvalidQRSize(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"qr_size":      "-100",
		"generate":     true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_generatedKeyInvalidKeySize(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "Test",
		"key_size":     "-100",
		"generate":     true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_generatedKeyMissingAccountName(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":   "Vault",
		"generate": true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_generatedKeyMissingIssuer(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"account_name": "test@email.com",
		"generate":     true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_invalidURLValue(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"url":      "notaurl",
		"generate": false,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_urlAndGenerateTrue(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"url":      "otpauth://totp/Vault:test@email.com?secret=HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ&algorithm=SHA512&digits=6&period=60",
		"generate": true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_keyAndGenerateTrue(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"key":      "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ",
		"generate": true,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, true),
			testAccStepReadKey(t, "test", nil),
		},
	})
}

func TestBackend_generatedKeyExportedFalse(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	keyData := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "test@email.com",
		"generate":     true,
		"exported":     false,
	}

	expected := map[string]interface{}{
		"issuer":       "Vault",
		"account_name": "test@email.com",
		"digits":       otplib.DigitsSix,
		"period":       30,
		"algorithm":    otplib.AlgorithmSHA1,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepCreateKey(t, "test", keyData, false),
			testAccStepReadKey(t, "test", expected),
		},
	})
}

func testAccStepCreateKey(t *testing.T, name string, keyData map[string]interface{}, expectFail bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      path.Join("keys", name),
		Data:      keyData,
		ErrorOk:   expectFail,
		Check: func(resp *logical.Response) error {
			//Skip this if the key is not generated by vault or if the test is expected to fail
			if !keyData["generate"].(bool) || expectFail {
				return nil
			}

			// Check to see if barcode and url were returned if exported is false
			if !keyData["exported"].(bool) {
				if resp != nil {
					t.Fatalf("data was returned when exported was set to false")
				}
				return nil
			}

			// Check to see if a barcode was returned when qr_size is zero
			if keyData["qr_size"].(int) == 0 {
				if _, exists := resp.Data["barcode"]; exists {
					t.Fatalf("a barcode was returned when qr_size was set to zero")
				}
				return nil
			}

			var d struct {
				Url     string `mapstructure:"url"`
				Barcode string `mapstructure:"barcode"`
			}

			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			//Check to see if barcode and url are returned
			if d.Barcode == "" {
				t.Fatalf("a barcode was not returned for a generated key")
			}

			if d.Url == "" {
				t.Fatalf("a url was not returned for a generated key")
			}

			//Parse url
			urlObject, err := url.Parse(d.Url)

			if err != nil {
				t.Fatal("an error occurred while parsing url string")
			}

			//Set up query object
			urlQuery := urlObject.Query()

			//Read secret
			urlSecret := urlQuery.Get("secret")

			//Check key length
			keySize := keyData["key_size"].(int)
			correctSecretStringSize := (keySize / 5) * 8
			actualSecretStringSize := len(urlSecret)

			if actualSecretStringSize != correctSecretStringSize {
				t.Fatal("incorrect key string length")
			}

			return nil
		},
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
				return fmt.Errorf("account_name should equal: %s", expected["account_name"])
			case d.Period != uint(period):
				return fmt.Errorf("period should equal: %d", expected["period"])
			case keyAlgorithm != expected["algorithm"]:
				return fmt.Errorf("algorithm should equal: %s", expected["algorithm"])
			case d.Digits != expected["digits"]:
				return fmt.Errorf("digits should equal: %d", expected["digits"])
			}
			return nil
		},
	}
}

func testAccStepValidateCode(t *testing.T, name string, code string, valid, expectError bool) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "code/" + name,
		Data: map[string]interface{}{
			"code": code,
		},
		ErrorOk: expectError,
		Check: func(resp *logical.Response) error {
			if resp == nil {
				return fmt.Errorf("bad: %#v", resp)
			}

			var d struct {
				Valid bool `mapstructure:"valid"`
			}

			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}

			switch valid {
			case true:
				if d.Valid != true {
					return fmt.Errorf("code was not valid: %s", code)
				}

			default:
				if d.Valid != false {
					return fmt.Errorf("code was incorrectly validated: %s", code)
				}
			}
			return nil
		},
	}
}
