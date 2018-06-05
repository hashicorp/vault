package azure

import (
	"fmt"
	"os"
	"testing"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"

	storage "github.com/Azure/azure-sdk-for-go/storage"
)

func TestAzureBackend(t *testing.T) {
	if os.Getenv("AZURE_ACCOUNT_NAME") == "" ||
		os.Getenv("AZURE_ACCOUNT_KEY") == "" {
		t.SkipNow()
	}

	accountName := os.Getenv("AZURE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_ACCOUNT_KEY")

	ts := time.Now().UnixNano()
	name := fmt.Sprintf("vault-test-%d", ts)

	cleanupClient, _ := storage.NewBasicClient(accountName, accountKey)
	cleanupClient.HTTPClient = cleanhttp.DefaultPooledClient()

	logger := logging.NewVaultLogger(log.Debug)

	backend, err := NewAzureBackend(map[string]string{
		"container":   name,
		"accountName": accountName,
		"accountKey":  accountKey,
	}, logger)

	defer func() {
		blobService := cleanupClient.GetBlobService()
		container := blobService.GetContainerReference(name)
		container.DeleteIfExists(nil)
	}()

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, backend)
	physical.ExerciseBackend_ListPrefix(t, backend)
}
