package local

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

func TestLocalSeal_RoundTrip(t *testing.T) {
	l := newLocalSealTest(t)
	defer l.Cleanup()
	l.WriteRandKeys([]string{"test.0.key", "test.1.key", "test.2.key", "should_be_ignored"})

	sealConfig := map[string]string{
		"key_glob": filepath.Join(l.Dir, "*.key"),
	}

	s := NewSeal(logging.NewVaultLogger(log.Trace))
	_, err := s.SetConfig(sealConfig)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	input := []byte("foo")
	enc, err := s.Encrypt(context.Background(), input)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// verify Encrypt() uses the key with the biggest index
	if enc.KeyInfo == nil {
		t.Fatal("KeyInfo not set")
	}
	if enc.KeyInfo.KeyID != "2" {
		t.Fatalf("wrong key used: %s", enc.KeyInfo.KeyID)
	}

	// create another key to verify that decrypting doesn't just use the
	// current key
	l.WriteRandKey("test.3.key")

	dec, err := s.Decrypt(context.Background(), enc)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(input, dec) {
		t.Fatalf("expected %s, got %s", input, dec)
	}
}

func newLocalSealTest(t *testing.T) *localSealTest {
	dir, err := ioutil.TempDir("", "local_seal_test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return &localSealTest{
		t:   t,
		Dir: dir,
	}
}

type localSealTest struct {
	Dir string
	t   *testing.T
}

func (l *localSealTest) Cleanup() {
	os.RemoveAll(l.Dir)
}

func (l *localSealTest) WriteRandKey(file string) {
	key, err := uuid.GenerateRandomBytes(KeyLen)
	if err != nil {
		l.t.Fatalf("err: %s", err)
	}

	fileWithPath := filepath.Join(l.Dir, file)
	if err := ioutil.WriteFile(fileWithPath, key, 0444); err != nil {
		l.t.Fatalf("err: %s", err)
	}
}

func (l *localSealTest) WriteRandKeys(files []string) {
	for _, f := range files {
		l.WriteRandKey(f)
	}
}
