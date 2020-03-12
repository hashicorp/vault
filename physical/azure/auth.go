package azure

import (
	"context"
	"strings"

	storagemgmt "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/pkg/errors"
)

// GetStorageAccountKey gets a storage account key using MSI
func GetStorageAccountKey(storageAccountName, subscriptionID, resourceGroupName string) (string, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return "", errors.Wrap(err, "Unable to get token from environment")
	}
	// get storageAccount client
	storageAccountsClient := storagemgmt.NewAccountsClient(subscriptionID)
	storageAccountsClient.Authorizer = authorizer

	// get storage key
	res, err := storageAccountsClient.ListKeys(context.TODO(), resourceGroupName, storageAccountName, storagemgmt.Kerb)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if res.Keys == nil || len(*res.Keys) == 0 {
		return "", errors.New("No storage keys found")
	}
	var storageKey string
	for _, key := range *res.Keys {
		// uppercase both strings for comparison because the ListKeys call returns e.g. "FULL" but
		// the storagemgmt.Full constant in the SDK is defined as "Full".
		if strings.EqualFold(string(key.Permissions), string(storagemgmt.Full)) {
			storageKey = *key.Value
			break
		}
	}

	if storageKey == "" {
		return "", errors.New("No storage key with Full permissions found")
	}

	return storageKey, nil
}
