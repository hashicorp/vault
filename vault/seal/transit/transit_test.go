package transit

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
)

type testTransitClient struct {
	keyID string
	seal  seal.Access
}

func newTestTransitClient(keyID string) *testTransitClient {
	return &testTransitClient{
		keyID: keyID,
		seal:  seal.NewTestSeal(nil),
	}
}

func (m *testTransitClient) Close() {}

func (m *testTransitClient) Encrypt(plaintext []byte) ([]byte, error) {
	v, err := m.seal.Encrypt(context.Background(), plaintext)
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("v1:%s:%s", m.keyID, string(v.Ciphertext))), nil
}

func (m *testTransitClient) Decrypt(ciphertext []byte) ([]byte, error) {
	splitKey := strings.Split(string(ciphertext), ":")
	if len(splitKey) != 3 {
		return nil, errors.New("invalid ciphertext returned")
	}

	data := &physical.EncryptedBlobInfo{
		Ciphertext: []byte(splitKey[2]),
	}
	v, err := m.seal.Decrypt(context.Background(), data)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func TestTransitSeal_Lifecycle(t *testing.T) {
	s := NewSeal(logging.NewVaultLogger(log.Trace))

	keyID := "test-key"
	s.client = newTestTransitClient(keyID)

	// Test Encrypt and Decrypt calls
	input := []byte("foo")
	swi, err := s.Encrypt(context.Background(), input)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	pt, err := s.Decrypt(context.Background(), swi)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	if !reflect.DeepEqual(input, pt) {
		t.Fatalf("expected %s, got %s", input, pt)
	}

	if s.KeyID() != keyID {
		t.Fatalf("key id does not match: expected %s, got %s", keyID, s.KeyID())
	}
}
