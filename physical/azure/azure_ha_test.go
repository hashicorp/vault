package azure

import (
	"fmt"
	"os"
	"testing"
	"time"

	storage "github.com/Azure/azure-sdk-for-go/storage"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func TestAzureHABackend(t *testing.T) {
	if os.Getenv("AZURE_ACCOUNT_NAME") == "" ||
		os.Getenv("AZURE_ACCOUNT_KEY") == "" {
		t.SkipNow()
	}

	accountName := os.Getenv("AZURE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_ACCOUNT_KEY")
	environmentName := os.Getenv("AZURE_ENVIRONMENT")
	environmentUrl := os.Getenv("AZURE_ARM_ENDPOINT")

	ts := time.Now().UnixNano()
	name := fmt.Sprintf("vault-test-%d", ts)

	cleanupEnvironment, err := environmentForCleanupClient(environmentName, environmentUrl)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	cleanupClient, _ := storage.NewBasicClientOnSovereignCloud(accountName, accountKey, cleanupEnvironment)
	cleanupClient.HTTPClient = cleanhttp.DefaultPooledClient()

	logger := logging.NewVaultLogger(log.Debug)

	defer func() {
		blobService := cleanupClient.GetBlobService()
		container := blobService.GetContainerReference(name)
		container.DeleteIfExists(nil)
	}()

	var b [2]physical.Backend
	for i := 0; i < 2; i++ {
		b[i], err = NewAzureBackend(map[string]string{
			"container":    name,
			"accountName":  accountName,
			"accountKey":   accountKey,
			"environment":  environmentName,
			"arm_endpoint": environmentUrl,
			"ha_enabled":   "true",
		}, logger)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}

	physical.ExerciseHABackend(t, b[0].(physical.HABackend), b[1].(physical.HABackend))
}
