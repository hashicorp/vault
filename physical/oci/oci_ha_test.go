// Copyright © 2019, Oracle and/or its affiliates.
package oci

import (
	"os"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
)

func TestOCIHABackend(t *testing.T) {
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

	backend := createBackend(bucketName, namespaceName, "true", bucketName, t)
	ha, ok := backend.(physical.HABackend)
	if !ok {
		t.Fatalf("does not implement")
	}

	physical.ExerciseHABackend(t, ha, ha)
}
