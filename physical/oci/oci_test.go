// Copyright © 2019, Oracle and/or its affiliates.
package oci

import (
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
	"golang.org/x/net/context"
)

func TestOCIBackend(t *testing.T) {
	// Skip tests if we are not running acceptance tests
	if os.Getenv("VAULT_ACC") == "" {
		t.SkipNow()
	}

	if !hasOCICredentials() {
		t.Skip("Skipping because OCI credentials could not be resolved. See https://pkg.go.dev/github.com/oracle/oci-go-sdk/common#DefaultConfigProvider for information on how to set up OCI credentials.")
	}

	bucketName, _ := uuid.GenerateUUID()
	configProvider := common.DefaultConfigProvider()
	objectStorageClient, _ := objectstorage.NewObjectStorageClientWithConfigurationProvider(configProvider)
	namespaceName := getNamespaceName(objectStorageClient, t)

	createBucket(bucketName, getTenancyOcid(configProvider, t), namespaceName, objectStorageClient, t)
	defer deleteBucket(namespaceName, bucketName, objectStorageClient, t)

	backend := createBackend(bucketName, namespaceName, "false", "", t)
	physical.ExerciseBackend(t, backend)
	physical.ExerciseBackend_ListPrefix(t, backend)
}

func createBucket(bucketName string, tenancyOcid string, namespaceName string, objectStorageClient objectstorage.ObjectStorageClient, t *testing.T) {
	createBucketRequest := objectstorage.CreateBucketRequest{
		NamespaceName: &namespaceName,
	}
	createBucketRequest.CompartmentId = &tenancyOcid
	createBucketRequest.Name = &bucketName
	createBucketRequest.Metadata = make(map[string]string)
	createBucketRequest.PublicAccessType = objectstorage.CreateBucketDetailsPublicAccessTypeNopublicaccess
	_, err := objectStorageClient.CreateBucket(context.Background(), createBucketRequest)
	if err != nil {
		t.Fatalf("Failed to create bucket: %v", err)
	}
}

func deleteBucket(nameSpaceName string, bucketName string, objectStorageClient objectstorage.ObjectStorageClient, t *testing.T) {
	request := objectstorage.DeleteBucketRequest{
		NamespaceName: &nameSpaceName,
		BucketName:    &bucketName,
	}
	_, err := objectStorageClient.DeleteBucket(context.Background(), request)
	if err != nil {
		t.Fatalf("Failed to delete bucket: %v", err)
	}
}

func getTenancyOcid(configProvider common.ConfigurationProvider, t *testing.T) (tenancyOcid string) {
	tenancyOcid, err := configProvider.TenancyOCID()
	if err != nil {
		t.Fatalf("Failed to get tenancy ocid: %v", err)
	}
	return tenancyOcid
}

func createBackend(bucketName string, namespaceName string, haEnabledStr string, lockBucketName string, t *testing.T) physical.Backend {
	backend, err := NewBackend(map[string]string{
		"auth_type_api_key": "true",
		"bucket_name":       bucketName,
		"namespace_name":    namespaceName,
		"ha_enabled":        haEnabledStr,
		"lock_bucket_name":  lockBucketName,
	}, logging.NewVaultLogger(log.Trace))
	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}
	return backend
}

func getNamespaceName(objectStorageClient objectstorage.ObjectStorageClient, t *testing.T) string {
	response, err := objectStorageClient.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{})
	if err != nil {
		t.Fatalf("Failed to get namespaceName: %v", err)
	}
	nameSpaceName := *response.Value
	return nameSpaceName
}

func hasOCICredentials() bool {
	configProvider := common.DefaultConfigProvider()

	_, err := configProvider.KeyID()
	if err != nil {
		return false
	}

	return true
}
