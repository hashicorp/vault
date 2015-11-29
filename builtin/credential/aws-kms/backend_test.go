package awsKms

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/logical"
	logicaltest "github.com/hashicorp/vault/logical/testing"
)

var keyid string = "FF"
var acctest bool = false

func TestBackend_basic(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) != "" {
		keyid = os.Getenv("AWS_KMS_KEY_ID")
		acctest = true
	}
	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Factory:  Factory,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepMapKeyId(t),
			testAccLogin(t, ""),
			testAccLoginInvalid(t, ""),
			testAccLoginTokenExpired(t, ""),
			testAccLoginTokenInFuture(t, ""),
			testAccLoginInvalidToken(t, ""),
			testAccStepDeleteKeyId(t),
			testAccKeyDeleted(t, ""),
		},
	})
}

func TestBackend_displayName(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) != "" {
		keyid = os.Getenv("AWS_KMS_KEY_ID")
		acctest = true
	}
	logicaltest.Test(t, logicaltest.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Factory:  Factory,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t),
			testAccStepMapAppIdDisplayName(t),
			testAccLogin(t, "tubbin"),
			testAccLoginInvalid(t, ""),
			testAccStepDeleteKeyId(t),
			testAccKeyDeleted(t, ""),
		},
	})
}

func testAccPreCheck(t *testing.T) {

	if v := os.Getenv("AWS_KMS_KEY_ID"); v == "" {
		t.Fatal("AWS_KMS_KEY_ID must be set for acceptance tests")
	}

	if v := os.Getenv("AWS_ACCESS_KEY_ID"); v == "" {
		t.Fatal("AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}

	if v := os.Getenv("AWS_SECRET_ACCESS_KEY"); v == "" {
		t.Fatal("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
}

func testAccStepConfig(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "config/user",
		Data: map[string]interface{}{
			"access_key": os.Getenv("AWS_ACCESS_KEY_ID"),
			"secret_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	}
}

func testAccStepMapKeyId(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "map/key-id/" + keyid,
		Data: map[string]interface{}{
			"value": "foo,bar",
		},
	}
}

func testAccStepMapAppIdDisplayName(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "map/key-id/" + keyid,
		Data: map[string]interface{}{
			"display_name": "tubbin",
			"value":        "foo,bar",
		},
	}
}

func testAccStepDeleteKeyId(t *testing.T) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "map/key-id/" + keyid,
	}
}

type tokeninv struct {
	foo string
}

func createToken(t *testing.T, invalidType bool, invalidCipherText bool, expired bool, infuture bool) string {
	if !acctest {
		return "FF"
	}
	var err error
	if invalidCipherText {
		return "0F01"
	}

	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(os.Getenv("AWS_DEFAULT_REGION")),
		HTTPClient:  cleanhttp.DefaultClient(),
	}
	svc := kms.New(session.New(awsConfig))
	var notAfter time.Time
	var notBefore time.Time
	if expired {
		notAfter = time.Now().UTC().Add(-5 * time.Minute)
	} else {
		notAfter = time.Now().UTC().Add(10 * time.Minute)
	}
	if infuture {
		notBefore = time.Now().UTC().Add(5 * time.Minute)
	} else {
		notBefore = time.Now().UTC().Add(-10 * time.Minute)
	}
	var tokJson []byte
	if invalidType {
		tok := &tokeninv{
			foo: "foobar",
		}
		tokJson, err = json.Marshal(tok)
	} else {
		tok := &token{
			NotBefore: notAfter,
			NotAfter:  notBefore,
		}
		tokJson, err = json.Marshal(tok)
	}

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	params := &kms.EncryptInput{
		KeyId:     &keyid,
		Plaintext: tokJson,
	}

	resp, err := svc.Encrypt(params)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	tokEncrypted := make([]byte, hex.EncodedLen(len(resp.CiphertextBlob)))
	n := hex.Encode(tokEncrypted, resp.CiphertextBlob)
	if n <= 0 {
		t.Fatalf("Could not encode to hex, length %v", n)
	}
	return string(tokEncrypted[:])

}

func testAccLogin(t *testing.T, display string) logicaltest.TestStep {
	tokEncrypted := createToken(t, false, false, false, false)

	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"key_id":     keyid,
			"ciphertext": tokEncrypted,
		},
		Unauthenticated: true,

		Check: logicaltest.TestCheckMulti(
			logicaltest.TestCheckAuth([]string{"bar", "foo"}),
			logicaltest.TestCheckAuthDisplayName(display),
		),
	}
}

func testAccLoginInvalidToken(t *testing.T, display string) logicaltest.TestStep {
	tokEncrypted := createToken(t, true, false, false, false)
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"key_id":     keyid,
			"ciphertext": tokEncrypted,
		},
		ErrorOk:         true,
		Unauthenticated: true,

		Check: logicaltest.TestCheckError(),
	}
}

func testAccLoginInvalid(t *testing.T, display string) logicaltest.TestStep {
	tokEncrypted := createToken(t, false, true, false, false)
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"key_id":     keyid,
			"ciphertext": tokEncrypted,
		},
		ErrorOk:         true,
		Unauthenticated: true,

		Check: logicaltest.TestCheckError(),
	}
}

func testAccKeyDeleted(t *testing.T, display string) logicaltest.TestStep {
	tokEncrypted := createToken(t, false, false, false, false)
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"key_id":     keyid,
			"ciphertext": tokEncrypted,
		},
		ErrorOk:         true,
		Unauthenticated: true,

		Check: logicaltest.TestCheckError(),
	}
}
func testAccLoginTokenExpired(t *testing.T, display string) logicaltest.TestStep {
	tokEncrypted := createToken(t, false, false, true, false)
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"key_id":     keyid,
			"ciphertext": tokEncrypted,
		},
		ErrorOk:         true,
		Unauthenticated: true,

		Check: logicaltest.TestCheckError(),
	}
}
func testAccLoginTokenInFuture(t *testing.T, display string) logicaltest.TestStep {
	tokEncrypted := createToken(t, false, false, false, true)
	return logicaltest.TestStep{
		Operation: logical.WriteOperation,
		Path:      "login",
		Data: map[string]interface{}{
			"key_id":     keyid,
			"ciphertext": tokEncrypted,
		},
		ErrorOk:         true,
		Unauthenticated: true,

		Check: logicaltest.TestCheckError(),
	}
}
