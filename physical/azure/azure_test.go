package azure

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/azure"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
)

func environmentForCleanupClient(name string, armURL string) (azure.Environment, error) {
	if armURL != "" {
		return azure.EnvironmentFromURL(armURL)
	}
	if name == "" {
		name = "AzurePublicCloud"
	}
	return azure.EnvironmentFromName(name)
}

func testFixture(t *testing.T) (physical.Backend, func()) {
	t.Helper()
	accountName := os.Getenv("AZURE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_ACCOUNT_KEY")
	environmentName := os.Getenv("AZURE_ENVIRONMENT")
	environmentURL := os.Getenv("AZURE_ARM_ENDPOINT")

	ts := time.Now().UnixNano()
	name := fmt.Sprintf("vault-test-%d", ts)

	cleanupEnvironment, err := environmentForCleanupClient(environmentName, environmentURL)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.%s/%s", accountName, cleanupEnvironment.StorageEndpointSuffix, name))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	containerURL := azblob.NewContainerURL(*URL, p)

	logger := logging.NewVaultLogger(log.Debug)

	backend, err := NewAzureBackend(map[string]string{
		"container":    name,
		"accountName":  accountName,
		"accountKey":   accountKey,
		"environment":  environmentName,
		"arm_endpoint": environmentURL,
	}, logger)

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return backend, func() {
		ctx := context.Background()
		blobService, err := containerURL.GetProperties(ctx, azblob.LeaseAccessConditions{})
		if err != nil {
			t.Logf("failed to retrieve blob container info: %v", err)
			return
		}

		if blobService.StatusCode() == 200 {
			_, err := containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
			if err != nil {
				t.Logf("clean up failed: %v", err)
			}
		}
	}
}

func TestAzureBackend(t *testing.T) {
	if os.Getenv("AZURE_ACCOUNT_NAME") == "" ||
		os.Getenv("AZURE_ACCOUNT_KEY") == "" {
		t.SkipNow()
	}

	backend, cleanup := testFixture(t)
	defer cleanup()

	physical.ExerciseBackend(t, backend)
	physical.ExerciseBackend_ListPrefix(t, backend)
}

func TestAzureBackend_ListPaging(t *testing.T) {
	if os.Getenv("AZURE_ACCOUNT_NAME") == "" ||
		os.Getenv("AZURE_ACCOUNT_KEY") == "" {
		t.SkipNow()
	}

	backend, cleanup := testFixture(t)
	defer cleanup()

	// by default, azure returns 5000 results in a page, load up more than that
	for i := 0; i < MaxListResults+100; i++ {
		if err := backend.Put(context.Background(), &physical.Entry{
			Key:   strconv.Itoa(i),
			Value: []byte(strconv.Itoa(i)),
		}); err != nil {
			t.Fatalf("err: %s", err)
		}
	}

	results, err := backend.List(context.Background(), "")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(results) != MaxListResults+100 {
		t.Fatalf("expected %d, got %d", MaxListResults+100, len(results))
	}
}
