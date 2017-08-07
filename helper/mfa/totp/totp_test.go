package totp

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	otplib "github.com/pquerna/otp"
	totplib "github.com/pquerna/otp/totp"
)

func setupTestBackend(conf *logical.BackendConfig) (*framework.Backend, error) {
	b := makeTestBackend()
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

func makeTestBackend() *framework.Backend {
	b := framework.Backend{
		Help:        "",
		Paths:       []*framework.Path{},
		Secrets:     []*framework.Secret{},
		BackendType: logical.TypeLogical,
	}
	b.Paths = TotpPaths(&b)
	return &b
}

func generateKey() (string, error) {
	keyUrl, err := totplib.Generate(totplib.GenerateOpts{
		Issuer:      "Vault",
		AccountName: "Test",
	})

	key := keyUrl.Secret()

	return key, err
}

func createKeyPath(t *testing.T, b *framework.Backend, s logical.Storage, secret string) {
	keyData := map[string]interface{}{
		"key":      secret,
		"generate": false,
	}
	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "totp/keys/test",
		Data:      keyData,
		Storage:   s,
	})
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to create role")
	}
	if err != nil {
		t.Fatal(err)
	}
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

func newTestFieldData() *framework.FieldData {
	d := framework.FieldData{
		Raw:    make(map[string]interface{}),
		Schema: make(map[string]*framework.FieldSchema),
	}
	d.Schema["passcode"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "One time passcode (optional)",
	}
	d.Schema["method"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "Multi-factor auth method to use (optional)",
	}
	return &d
}

func TestTotpHandlerSuccess(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := setupTestBackend(config)
	if err != nil {
		t.Fatal(err)
	}

	key, err := generateKey()
	if err != nil {
		t.Fatal(err)
	}
	createKeyPath(t, b, storage, key)

	successResp := &logical.Response{
		Auth: &logical.Auth{Metadata: make(map[string]string)},
	}
	successResp.Auth.Metadata["username"] = "test"

	fieldData := newTestFieldData()
	code, _ := generateCode(key, 30, otplib.DigitsSix, otplib.AlgorithmSHA1)
	fieldData.Raw["passcode"] = code

	handler := GetTotpHandler(b)
	req := logical.Request{
		Storage: storage,
	}
	resp, err := handler(&req, fieldData, successResp)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if resp != successResp {
		t.Fatalf("Testing Totp authentication gave incorrect response (expected success, got: %v)", resp)
	}
}

func TestTotpHandlerReject(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage
	b, err := setupTestBackend(config)
	if err != nil {
		t.Fatal(err)
	}

	key, err := generateKey()
	if err != nil {
		t.Fatal(err)
	}
	createKeyPath(t, b, storage, key)

	successResp := &logical.Response{
		Auth: &logical.Auth{Metadata: make(map[string]string)},
	}
	successResp.Auth.Metadata["username"] = "test"

	fieldData := newTestFieldData()
	code, _ := generateCode(key, 30, otplib.DigitsSix, otplib.AlgorithmSHA1)
	codeI, _ := strconv.Atoi(code)
	fieldData.Raw["passcode"] = strconv.Itoa(codeI + 1)

	handler := GetTotpHandler(b)
	req := logical.Request{
		Storage: storage,
	}
	resp, err := handler(&req, fieldData, successResp)

	if err != nil {
		t.Fatalf(err.Error())
	}
	error, ok := resp.Data["error"].(string)
	if !ok || !strings.Contains(error, "The specified passcode is not valid") {
		t.Fatalf("Testing Duo authentication gave incorrect response (expected deny, got: %v)", error)
	}
}
