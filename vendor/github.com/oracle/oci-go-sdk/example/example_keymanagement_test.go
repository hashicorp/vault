// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
//
// Example code for Key Management Service API
//

package example

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/keymanagement"
)

// ExampleKeyManagement_VaultOperations shows how to create, schedule deletion
// and cancel a scheduled deletion of a KMS vault
func ExampleVaultOperations() {
	vaultClient, clientError := keymanagement.NewKmsVaultClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clientError)

	ctx := context.Background()
	vaultName := "KmsVault"
	updatedVaultName := "UpdatedKmsVault"

	vault := createVault(ctx, vaultClient, vaultName)
	defer cleanupResources(ctx, vaultClient, vault.Id)
	// wait for instance lifecycle state becomes active
	waitForStateVaultClient(ctx, vault.Id, vaultClient, keymanagement.VaultLifecycleStateActive)

	updatedVault := updateVault(ctx, vaultClient, &updatedVaultName, vault.Id)
	fmt.Println(fmt.Sprintf("Updated vault display name %s", *updatedVault.DisplayName))

	svdErr := scheduleVaultDeletion(ctx, vaultClient, vault.Id)
	helpers.FatalIfError(svdErr)
	waitForStateVaultClient(ctx, vault.Id, vaultClient, keymanagement.VaultLifecycleStatePendingDeletion)

	cvdErr := cancelVaultDeletion(ctx, vaultClient, vault.Id)
	helpers.FatalIfError(cvdErr)
	waitForStateVaultClient(ctx, vault.Id, vaultClient, keymanagement.VaultLifecycleStateActive)

	// Output:
	// create vault
	// update vault
	// schedule vault deletion
	// cancel vault deletion
	// schedule vault deletion
}

// ExampleKeyManagement_KeyOperations shows how to create, enable and disable a KMS key
func ExampleKeyOperations() {
	vaultClient, clientError := keymanagement.NewKmsVaultClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clientError)

	ctx := context.Background()
	vaultName := "KmsVault"
	keyName := "KmsKey"
	updatedKeyName := "UpdatedKmsKey"

	vault := createVault(ctx, vaultClient, vaultName)
	defer cleanupResources(ctx, vaultClient, vault.Id)
	// wait for instance lifecycle state becomes active
	waitForStateVaultClient(ctx, vault.Id, vaultClient, keymanagement.VaultLifecycleStateActive)

	vaultManagementClient, mgmtClientError := keymanagement.
		NewKmsManagementClientWithConfigurationProvider(common.DefaultConfigProvider(), *vault.ManagementEndpoint)
	helpers.FatalIfError(mgmtClientError)

	// Create Key
	key, _ := createKey(ctx, vaultManagementClient, &keyName)

	// Disable Key
	disableRequest := keymanagement.DisableKeyRequest{
		KeyId: key.Id,
	}

	disableResponse, disableErr := vaultManagementClient.DisableKey(ctx, disableRequest)
	helpers.FatalIfError(disableErr)
	key = disableResponse.Key
	// Wait for key to be in Disabled state
	waitForStateVaultManagementClient(ctx, key.Id, vaultManagementClient, keymanagement.KeyLifecycleStateDisabled)

	fmt.Println("disable key")

	// Enable Key
	enableRequest := keymanagement.EnableKeyRequest{
		KeyId: key.Id,
	}

	enableResponse, enableErr := vaultManagementClient.EnableKey(ctx, enableRequest)
	helpers.FatalIfError(enableErr)
	key = enableResponse.Key
	// Wait for key to be in Enabled state
	waitForStateVaultManagementClient(ctx, key.Id, vaultManagementClient, keymanagement.KeyLifecycleStateEnabled)

	fmt.Println("enable key")

	// Schedule Key Deletion
	scheduleKeyDeletionRequest := keymanagement.ScheduleKeyDeletionRequest{
		KeyId: key.Id,
	}

	scheduleKeyDeletionResponse, scheduleKeyDeletionErr := vaultManagementClient.ScheduleKeyDeletion(ctx, scheduleKeyDeletionRequest)
	helpers.FatalIfError(scheduleKeyDeletionErr)
	key = scheduleKeyDeletionResponse.Key
	// Wait for key to be in PendingDeletion state
	waitForStateVaultManagementClient(ctx, key.Id, vaultManagementClient, keymanagement.KeyLifecycleStatePendingDeletion)

	fmt.Println("schedule key deletion")

	// Cancel Key Deletion
	cancelKeyDeletionRequest := keymanagement.CancelKeyDeletionRequest{
		KeyId: key.Id,
	}

	cancelKeyDeletionResponse, cancelKeyDeletionErr := vaultManagementClient.CancelKeyDeletion(ctx, cancelKeyDeletionRequest)
	helpers.FatalIfError(cancelKeyDeletionErr)
	key = cancelKeyDeletionResponse.Key
	// Wait for key to be in Enabled state
	waitForStateVaultManagementClient(ctx, key.Id, vaultManagementClient, keymanagement.KeyLifecycleStateEnabled)

	fmt.Println("cancel scheduled key deletion")

	// Update Key
	updateKeyDetails := keymanagement.UpdateKeyDetails{
		DisplayName: &updatedKeyName,
	}
	updateKeyRequest := keymanagement.UpdateKeyRequest{
		KeyId:            key.Id,
		UpdateKeyDetails: updateKeyDetails,
	}

	updateResponse, updateErr := vaultManagementClient.UpdateKey(ctx, updateKeyRequest)
	helpers.FatalIfError(updateErr)
	key = updateResponse.Key

	fmt.Println("update key")

	// Output:
	// create vault
	// create key
	// disable key
	// enable key
	// schedule key deletion
	// cancel scheduled key deletion
	// update key
	// schedule vault deletion
}

func ExampleCryptoOperations() {
	vaultClient, clientError := keymanagement.NewKmsVaultClientWithConfigurationProvider(common.DefaultConfigProvider())
	helpers.FatalIfError(clientError)

	ctx := context.Background()
	vaultName := "KmsVault"
	keyName := "KmsKey"
	testInput := "CryptoOps Test Input"

	vault := createVault(ctx, vaultClient, vaultName)
	defer cleanupResources(ctx, vaultClient, vault.Id)
	// wait for instance lifecycle state becomes active
	waitForStateVaultClient(ctx, vault.Id, vaultClient, keymanagement.VaultLifecycleStateActive)

	vaultManagementClient, mgmtClientError := keymanagement.
		NewKmsManagementClientWithConfigurationProvider(common.DefaultConfigProvider(), *vault.ManagementEndpoint)
	helpers.FatalIfError(mgmtClientError)

	// Create Key
	key, keyShape := createKey(ctx, vaultManagementClient, &keyName)

	// Create crypto client
	vaultCryptoClient, cryptoClientError := keymanagement.
		NewKmsCryptoClientWithConfigurationProvider(common.DefaultConfigProvider(), *vault.CryptoEndpoint)
	helpers.FatalIfError(cryptoClientError)

	// Generate DEK
	includePlaintextKeyInResponse := true
	generateKeyDetails := keymanagement.GenerateKeyDetails{
		KeyId:               key.Id,
		KeyShape:            &keyShape,
		IncludePlaintextKey: &includePlaintextKeyInResponse,
	}
	generateDekRequest := keymanagement.GenerateDataEncryptionKeyRequest{
		GenerateKeyDetails: generateKeyDetails,
	}

	generateDekResponse, err := vaultCryptoClient.GenerateDataEncryptionKey(ctx, generateDekRequest)
	helpers.FatalIfError(err)
	generatedKey := generateDekResponse.GeneratedKey
	fmt.Println(fmt.Sprintf("Plaintext generated DEK: %s", *generatedKey.Plaintext))

	fmt.Println("generate DEK")

	// Encrypt
	encryptedDataDetails := keymanagement.EncryptDataDetails{
		KeyId:     key.Id,
		Plaintext: &testInput,
	}
	encryptRequest := keymanagement.EncryptRequest{
		EncryptDataDetails: encryptedDataDetails,
	}

	encryptResponse, encryptErr := vaultCryptoClient.Encrypt(ctx, encryptRequest)
	helpers.FatalIfError(encryptErr)

	cipherText := encryptResponse.Ciphertext

	fmt.Print("encrypt data")

	// Decrypt
	decryptDataDetails := keymanagement.DecryptDataDetails{
		KeyId:      key.Id,
		Ciphertext: cipherText,
	}
	decryptRequest := keymanagement.DecryptRequest{
		DecryptDataDetails: decryptDataDetails,
	}
	decryptResponse, decryptErr := vaultCryptoClient.Decrypt(ctx, decryptRequest)
	helpers.FatalIfError(decryptErr)

	plainText := decryptResponse.Plaintext
	fmt.Println(fmt.Sprintf("Decrypted plaintext: %s", *plainText))

	fmt.Print("decrypt data")

	// Output:
	// create vault
	// create key
	// Plaintext generated DEK: <generated key>
	// generate DEK
	// encrypt data
	// Decrypted plaintext: CryptoOps Test Input
	// decrypt data
	// schedule vault deletion
}

func getVault(ctx context.Context, client keymanagement.KmsVaultClient, retryPolicy *common.RetryPolicy, vaultId *string) keymanagement.Vault {

	request := keymanagement.GetVaultRequest{
		VaultId: vaultId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	}
	response, err := client.GetVault(ctx, request)
	helpers.FatalIfError(err)
	return response.Vault
}

func updateVault(ctx context.Context, client keymanagement.KmsVaultClient, newName, vaultId *string) keymanagement.Vault {
	updateVaultDetails := keymanagement.UpdateVaultDetails{
		DisplayName: newName,
	}
	request := keymanagement.UpdateVaultRequest{
		VaultId:            vaultId,
		UpdateVaultDetails: updateVaultDetails,
	}
	response, err := client.UpdateVault(ctx, request)
	helpers.FatalIfError(err)

	fmt.Println("update vault")
	return response.Vault
}

func getKey(ctx context.Context, client keymanagement.KmsManagementClient, retryPolicy *common.RetryPolicy, keyId *string) keymanagement.Key {

	request := keymanagement.GetKeyRequest{
		KeyId: keyId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	}
	response, err := client.GetKey(ctx, request)
	helpers.FatalIfError(err)
	return response.Key
}

func createVault(ctx context.Context, c keymanagement.KmsVaultClient, vaultName string) (vault keymanagement.Vault) {
	vaultDetails := keymanagement.CreateVaultDetails{
		CompartmentId: helpers.CompartmentID(),
		DisplayName:   &vaultName,
		VaultType:     keymanagement.CreateVaultDetailsVaultTypePrivate,
	}
	request := keymanagement.CreateVaultRequest{}
	request.CreateVaultDetails = vaultDetails
	response, err := c.CreateVault(ctx, request)
	helpers.FatalIfError(err)

	fmt.Println("create vault")
	return response.Vault
}

func scheduleVaultDeletion(ctx context.Context, c keymanagement.KmsVaultClient, vaultId *string) (err error) {
	request := keymanagement.ScheduleVaultDeletionRequest{
		VaultId: vaultId,
	}
	_, err = c.ScheduleVaultDeletion(ctx, request)
	helpers.FatalIfError(err)

	fmt.Println("schedule vault deletion")
	return
}

func cancelVaultDeletion(ctx context.Context, c keymanagement.KmsVaultClient, vaultId *string) (err error) {
	request := keymanagement.CancelVaultDeletionRequest{
		VaultId: vaultId,
	}
	_, err = c.CancelVaultDeletion(ctx, request)
	helpers.FatalIfError(err)

	fmt.Println("cancel vault deletion")
	return
}

func createKey(ctx context.Context, vaultManagementClient keymanagement.KmsManagementClient, keyName *string) (keymanagement.Key, keymanagement.KeyShape) {
	keyLength := 32

	keyShape := keymanagement.KeyShape{
		Algorithm: keymanagement.KeyShapeAlgorithmAes,
		Length:    &keyLength,
	}
	createKeyDetails := keymanagement.CreateKeyDetails{
		CompartmentId: helpers.CompartmentID(),
		KeyShape:      &keyShape,
		DisplayName:   keyName,
	}
	request := keymanagement.CreateKeyRequest{
		CreateKeyDetails: createKeyDetails,
	}

	response, err := vaultManagementClient.CreateKey(ctx, request)
	helpers.FatalIfError(err)
	key := response.Key

	// Wait for key to be in Enabled state
	waitForStateVaultManagementClient(ctx, key.Id, vaultManagementClient, keymanagement.KeyLifecycleStateEnabled)

	fmt.Println("create key")
	return key, keyShape
}

func waitForStateVaultClient(ctx context.Context, vaultId *string, client keymanagement.KmsVaultClient,
	state keymanagement.VaultLifecycleStateEnum) {
	// maximum times of retry
	attempts := uint(10)

	shouldRetry := func(r common.OCIOperationResponse) bool {
		if _, isServiceError := common.IsServiceError(r.Error); isServiceError {
			// not service error, could be network error or other errors which prevents
			// request send to server, will do retry here
			return true
		}

		if vaultResponse, ok := r.Response.(keymanagement.GetVaultResponse); ok {
			// do the retry until lifecycle state reaches the passed terminal state
			return vaultResponse.Vault.LifecycleState != state
		}

		return true
	}

	nextDuration := func(r common.OCIOperationResponse) time.Duration {
		// you might want wait longer for next retry when your previous one failed
		// this function will return the duration as:
		// 1s, 2s, 4s, 8s, 16s, 32s, 64s etc...
		return time.Duration(math.Pow(float64(2), float64(r.AttemptNumber-1))) * time.Second
	}

	lifecycleStateCheckRetryPolicy := common.NewRetryPolicy(attempts, shouldRetry, nextDuration)

	getVault(ctx, client, &lifecycleStateCheckRetryPolicy, vaultId)
}

func waitForStateVaultManagementClient(ctx context.Context, keyId *string, client keymanagement.KmsManagementClient,
	state keymanagement.KeyLifecycleStateEnum) {
	// maximum times of retry
	attempts := uint(10)

	shouldRetry := func(r common.OCIOperationResponse) bool {
		if _, isServiceError := common.IsServiceError(r.Error); isServiceError {
			// not service error, could be network error or other errors which prevents
			// request send to server, will do retry here
			return true
		}

		if keyResponse, ok := r.Response.(keymanagement.GetKeyResponse); ok {
			// do the retry until lifecycle state reaches the passed terminal state
			return keyResponse.Key.LifecycleState != state
		}

		return true
	}

	nextDuration := func(r common.OCIOperationResponse) time.Duration {
		// you might want wait longer for next retry when your previous one failed
		// this function will return the duration as:
		// 1s, 2s, 4s, 8s, 16s, 32s, 64s etc...
		return time.Duration(math.Pow(float64(2), float64(r.AttemptNumber-1))) * time.Second
	}

	lifecycleStateCheckRetryPolicy := common.NewRetryPolicy(attempts, shouldRetry, nextDuration)

	getKey(ctx, client, &lifecycleStateCheckRetryPolicy, keyId)
}

func cleanupResources(ctx context.Context, client keymanagement.KmsVaultClient, vaultId *string) {
	err := scheduleVaultDeletion(ctx, client, vaultId)
	helpers.FatalIfError(err)
}
