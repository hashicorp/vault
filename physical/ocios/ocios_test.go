package ocios

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
)

func TestOCIOSBackend(t *testing.T) {

	configProvider := common.DefaultConfigProvider()

	if _, err := configProvider.UserOCID(); err != nil {
		t.SkipNow()
	}

	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(configProvider)
	if err != nil {
		t.Fatalf("unable to create test client: %s", err)
	}

	namespaceResponse, err := client.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{})
	if err != nil {
		t.Fatalf("unable to get namespace: %s", err)
	}

	namespace := *namespaceResponse.Value

	namespaceMetadataResponse, err := client.GetNamespaceMetadata(context.Background(), objectstorage.GetNamespaceMetadataRequest{
		NamespaceName: &namespace,
	})
	if err != nil {
		t.Fatalf("unable to get namespace metadata: %s", err)
	}

	var randInt = rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	bucket := fmt.Sprintf("vault-ocios-testacc-%d", randInt)

	_, err = client.CreateBucket(context.Background(), objectstorage.CreateBucketRequest{
		NamespaceName: &namespace,
		CreateBucketDetails: objectstorage.CreateBucketDetails{
			Name:          &bucket,
			CompartmentId: namespaceMetadataResponse.DefaultS3CompartmentId,
		},
	})
	if err != nil {
		t.Fatalf("unable to create test bucket: %s", err)
	}

	defer func() {
		_, err := client.DeleteBucket(context.Background(), objectstorage.DeleteBucketRequest{
			NamespaceName: &namespace,
			BucketName:    &bucket,
		})
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	logger := logging.NewVaultLogger(log.Debug)

	// This uses the same logic to find the OCI OS credentials as we did at the beginning of the test
	b, err := NewOCIOSBackend(map[string]string{
		"bucket": bucket,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
}
