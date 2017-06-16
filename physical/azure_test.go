package physical

import (
	"fmt"
	"os"
	"testing"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"

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

	logger := logformat.NewVaultLogger(log.LevelTrace)

	backend, err := NewBackend("azure", logger, map[string]string{
		"container":   name,
		"accountName": accountName,
		"accountKey":  accountKey,
	})

	defer func() {
		blobService := cleanupClient.GetBlobService()
		container := blobService.GetContainerReference(name)
		container.DeleteIfExists(nil)
	}()

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, backend)
	testBackend_ListPrefix(t, backend)
}
