// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package manta

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	tt "github.com/joyent/triton-go/errors"
	"github.com/joyent/triton-go/storage"
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

	logger := logging.NewVaultLogger(log.Debug)
	mb := &MantaBackend{
		client:     client,
		directory:  testHarnessBucket,
		logger:     logger.Named("storage.mantabackend"),
		permitPool: permitpool.New(128),
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
			if !tt.IsResourceNotFoundError(err) {
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
