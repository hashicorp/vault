package physical

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/storage"
	"os"
	"testing"
	"time"
)

func TestAzureBackend(t *testing.T) {
	if os.Getenv("AZURE_ACCOUNT_NAME") == "" ||
		os.Getenv("AZURE_ACCOUNT_KEY") == "" {
		t.SkipNow()
	}

	accountName := os.Getenv("AZURE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_ACCOUNT_KEY")

	ts := time.Now().UnixNano()
	container := fmt.Sprintf("vault-test-%d", ts)

	cleanupClient, _ := storage.NewBasicClient(accountName, accountKey)

	backend, err := NewBackend("azure", map[string]string{
		"container":   container,
		"accountName": accountName,
		"accountKey":  accountKey,
	})

	defer func() {
		cleanupClient.GetBlobService().DeleteContainerIfExists(container)
	}()

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, backend)
	testBackend_ListPrefix(t, backend)
}
