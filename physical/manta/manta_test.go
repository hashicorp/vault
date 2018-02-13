package manta

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	"github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	tclient "github.com/joyent/triton-go/client"
	"github.com/joyent/triton-go/storage"
	log "github.com/mgutz/logxi/v1"
)

func TestMantaBackend(t *testing.T) {
	user := os.Getenv("MANTA_USER")
	keyId := os.Getenv("MANTA_KEY_ID")
	url := "https://us-east.manta.joyent.com"
	testHarnessBucket := fmt.Sprintf("test-bucket-%d", randInt())

	if user == "" || keyId == "" {
		t.SkipNow()
	}

	input := authentication.SSHAgentSignerInput{
		KeyID:       keyId,
		AccountName: user,
		Username:    "",
	}
	signer, err := authentication.NewSSHAgentSigner(input)
	if err != nil {
		t.Fatalf("Error Creating SSH Agent Signer: %s", err.Error())
	}

	config := &triton.ClientConfig{
		MantaURL:    url,
		AccountName: user,
		Signers:     []authentication.Signer{signer},
	}

	client, err := storage.NewClient(config)
	if err != nil {
		t.Fatalf("failed initialising Storage client: %s", err.Error())
	}

	logger := logformat.NewVaultLogger(log.LevelTrace)
	mb := &MantaBackend{
		client:     client,
		directory:  testHarnessBucket,
		logger:     logger,
		permitPool: physical.NewPermitPool(128),
	}

	err = mb.client.Dir().Put(context.Background(), &storage.PutDirectoryInput{
		DirectoryName: path.Join(mantaDefaultRootStore),
	})
	if err != nil {
		t.Fatal("Error creating test harness directory")
	}

	defer func() {
		err = mb.client.Dir().Delete(context.Background(), &storage.DeleteDirectoryInput{
			DirectoryName: path.Join(mantaDefaultRootStore, testHarnessBucket),
			ForceDelete:   true,
		})
		if err != nil {
			if !tclient.IsResourceNotFoundError(err) {
				t.Fatal("failed to delete test harness directory")
			}
		}
	}()

	physical.ExerciseBackend(t, mb)
	physical.ExerciseBackend_ListPrefix(t, mb)
}

func randInt() int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int()
}
