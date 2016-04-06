package awsKms

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
)

func createClientToken(keyid *string) (string, error) {
	awsConfig := &aws.Config{
		Region:     aws.String("us-east-1"),
		HTTPClient: cleanhttp.DefaultClient(),
	}
	svc := kms.New(session.New(awsConfig))

	var now = time.Now().UTC()
	notAfter := now.Add(10 * time.Minute)
	notBefore := now.Add(-10 * time.Minute)

	tok := &token{
		NotBefore: notBefore,
		NotAfter:  notAfter,
	}
	tokJson, err := json.Marshal(tok)

	if err != nil {
		return "", fmt.Errorf("err: %v", err)
	}

	params := &kms.EncryptInput{
		KeyId:     keyid,
		Plaintext: tokJson,
	}

	resp, err := svc.Encrypt(params)

	if err != nil {
		return "", fmt.Errorf("err: %v", err)
	}

	tokEncrypted := make([]byte, hex.EncodedLen(len(resp.CiphertextBlob)))
	n := hex.Encode(tokEncrypted, resp.CiphertextBlob)
	if n <= 0 {
		return "", fmt.Errorf("Could not encode to hex, length %v", n)
	}
	return string(tokEncrypted[:]), nil

}

type CLIHandler struct{}

func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
	mount, ok := m["mount"]
	if !ok {
		mount = "aws-kms"
	}
	keyid, ok := m["key"]
	if !ok {
		return "", fmt.Errorf("'key' var must be set")
	}
	tok, err := createClientToken(&keyid)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("auth/%s/login", mount)
	secret, err := c.Logical().Write(path, map[string]interface{}{
		"ciphertext": tok,
	})
	if err != nil {
		return "", err
	}
	if secret == nil {
		return "", fmt.Errorf("empty response from credential provider")
	}

	return secret.Auth.ClientToken, nil
}
func (h *CLIHandler) Help() string {
	help := `
The AWS KMS credential provider allows you to authenticate with AWS KMS.
To use it, specify "key" parameter. The value should be the same that
has been configured in Vault's AWS KMS authentication provider and you
must have kms:Encrypt permission to it.

    Example: vault auth -method=awskms key=12345678-90ab-cdef-1234-567890abcdef


    `

	return strings.TrimSpace(help)
}
