package mfa

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func TestTOTP_Generation_Verification(t *testing.T) {
	b, err := MFABackendFactory(logical.TestBackendConfig())
	if err != nil {
		t.Fatal(err)
	}

	storage := &logical.InmemStorage{}
	b.(*MFABackend).SetStorage(storage)

	// Create the TOTP role
	req := logical.TestRequest(t, logical.CreateOperation, "methods/test")
	req.Operation = logical.CreateOperation
	req.Data = map[string]interface{}{
		"type":                "totp",
		"totp_hash_algorithm": "sha256",
	}

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatal("got response, expected a nil response")
	}

	// Create the identifier and get back the TOTP information
	req.Operation = logical.CreateOperation
	req.Path = "methods/test/jeff"
	req.Data = map[string]interface{}{
		"totp_account_name": "jeff@hashicorp.com",
	}

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp == nil {
		t.Fatalf("got a nil response")
	}
	if resp.Data == nil {
		t.Fatalf("got nil data in response")
	}

	totpURL := resp.Data["totp_url"].(string)
	totpSecret := resp.Data["totp_secret"].(string)
	totpImage := resp.Data["totp_qrcode_png_b64"].(string)

	// URL check: verify that the key parses and shows the expected issuer
	key, err := otp.NewKeyFromURL(totpURL)
	if err != nil {
		t.Fatal(err)
	}
	if key == nil {
		t.Fatal("key from URL is nil")
	}
	expIssuer := fmt.Sprintf("Vault MFA: %s/%s (%s)", "test", "jeff", "jeff@hashicorp.com")
	if key.Issuer() != expIssuer {
		t.Fatalf("expected issuer of %s, got %s", expIssuer, key.Issuer())
	}

	// Secret check: generate a code and verify it
	code, err := totp.GenerateCode(totpSecret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if code == "" {
		t.Fatal("got empty code from secret")
	}
	if !totp.Validate(code, totpSecret) {
		t.Fatal("could not validate TOTP secret")
	}

	// Image check: make sure it's a parseable image
	imgBytes, err := base64.StdEncoding.DecodeString(totpImage)
	if err != nil {
		t.Fatal(err)
	}
	imgBuf := bytes.NewBuffer(imgBytes)
	image, err := png.Decode(imgBuf)
	if err != nil {
		t.Fatal(err)
	}
	if image == nil {
		t.Fatal("got a nil image from decoding")
	}
}
