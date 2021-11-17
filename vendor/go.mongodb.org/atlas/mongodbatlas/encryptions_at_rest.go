// Copyright 2021 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/openlyinc/pointy"
)

const (
	// CaCentral1 represents the CA_CENTRAL_1 America region for AWS Configuration.
	CaCentral1 = "CA_CENTRAL_1"
	// UsEast1 represents the US_EAST_1 America region for AWS Configuration.
	UsEast1 = "US_EAST_1"
	// UsEast2 represents the US_EAST_2 America region for AWS Configuration.
	UsEast2 = "US_EAST_2"
	// UsWest1 represents the US_WEST_1 America region for AWS Configuration.
	UsWest1 = "US_WEST_1"
	// UsWest2 represents the US_WEST_2 America region for AWS Configuration.
	UsWest2 = "US_WEST_2"
	// SaEast1 represents the SA_EAST_1 America region for AWS Configuration.
	SaEast1 = "SA_EAST_1"

	// ApNortheast1 represents the AP_NORTHEAST_1 Asia region for AWS Configuration.
	ApNortheast1 = "AP_NORTHEAST_1"
	// ApNortheast2 represents the AP_NORTHEAST_2 Asia region for AWS Configuration.
	ApNortheast2 = "AP_NORTHEAST_2"
	// ApSouth1 represents the AP_SOUTH_1 Asia region for AWS Configuration.
	ApSouth1 = "AP_SOUTH_1"
	// ApSoutheast1 represents the AP_SOUTHEAST_1 Asia region for AWS Configuration.
	ApSoutheast1 = "AP_SOUTHEAST_1"
	// ApSoutheast2 represents the AP_SOUTHEAST_2 Asia region for AWS Configuration.
	ApSoutheast2 = "AP_SOUTHEAST_2"

	// EuCentral1 represents the EU_CENTRAL_1 Europe region for AWS Configuration.
	EuCentral1 = "EU_CENTRAL_1"
	// EuWest1 represents the EU_WEST_1 Europe region for AWS Configuration.
	EuWest1 = "EU_WEST_1"
	// EuWest2 represents the EU_WEST_2 Europe region for AWS Configuration.
	EuWest2 = "EU_WEST_2"
	// EuWest3 represents the EU_WEST_3 Europe region for AWS Configuration.
	EuWest3 = "EU_WEST_3"

	// Azure represents `AZURE` where the Azure account credentials reside.
	Azure = "AZURE"
	// AzureChina represents `AZURE_CHINA` AZURE where the Azure account credentials reside.
	AzureChina = "AZURE_CHINA"
	// AzureGermany represents `AZURE_GERMANY` AZURE where the Azure account credentials reside.
	AzureGermany = "AZURE_GERMANY"

	encryptionsAtRestBasePath = "api/atlas/v1.0/groups/%s/encryptionAtRest"
)

// EncryptionsAtRestService is an interface for interfacing with the Encryption at Rest
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/encryption-at-rest/
type EncryptionsAtRestService interface {
	Create(context.Context, *EncryptionAtRest) (*EncryptionAtRest, *Response, error)
	Get(context.Context, string) (*EncryptionAtRest, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// EncryptionsAtRestServiceOp handles communication with the DatabaseUsers related methods of the
// MongoDB Atlas API.
type EncryptionsAtRestServiceOp service

var _ EncryptionsAtRestService = &EncryptionsAtRestServiceOp{}

// EncryptionAtRest represents a configuration Encryption at Rest for an Atlas project.
type EncryptionAtRest struct {
	GroupID        string                            `json:"groupId,omitempty"` // The unique identifier for the project.
	AwsKms         `json:"awsKms,omitempty"`         // AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
	AzureKeyVault  `json:"azureKeyVault,omitempty"`  // AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
	GoogleCloudKms `json:"googleCloudKms,omitempty"` // Specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
}

// AwsKms specifies AWS KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
type AwsKms struct {
	Enabled             *bool  `json:"enabled,omitempty"`             // Specifies whether Encryption at Rest is enabled for an Atlas project, To disable Encryption at Rest, pass only this parameter with a value of false, When you disable Encryption at Rest, Atlas also removes the configuration details.
	AccessKeyID         string `json:"accessKeyID,omitempty"`         // The IAM access key ID with permissions to access the customer master key specified by customerMasterKeyID.
	SecretAccessKey     string `json:"secretAccessKey,omitempty"`     // The IAM secret access key with permissions to access the customer master key specified by customerMasterKeyID.
	CustomerMasterKeyID string `json:"customerMasterKeyID,omitempty"` // The AWS customer master key used to encrypt and decrypt the MongoDB master keys.
	Region              string `json:"region,omitempty"`              // The AWS region in which the AWS customer master key exists: CA_CENTRAL_1, US_EAST_1, US_EAST_2, US_WEST_1, US_WEST_2, SA_EAST_1
	RoleID              string `json:"roleId,omitempty"`              // ID of an AWS IAM role authorized to manage an AWS customer master key.
	Valid               *bool  `json:"valid,omitempty"`               // Specifies whether the encryption key set for the provider is valid and may be used to encrypt and decrypt data.
}

// AzureKeyVault specifies Azure Key Vault configuration details and whether Encryption at Rest is enabled for an Atlas project.
type AzureKeyVault struct {
	Enabled           *bool  `json:"enabled,omitempty"`           // Specifies whether Encryption at Rest is enabled for an Atlas project. To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	ClientID          string `json:"clientID,omitempty"`          // The Client ID, also known as the application ID, for an Azure application associated with the Azure AD tenant.
	AzureEnvironment  string `json:"azureEnvironment,omitempty"`  // The Azure environment where the Azure account credentials reside. Valid values are the following: AZURE, AZURE_CHINA, AZURE_GERMANY
	SubscriptionID    string `json:"subscriptionID,omitempty"`    // The unique identifier associated with an Azure subscription.
	ResourceGroupName string `json:"resourceGroupName,omitempty"` // The name of the Azure Resource group that contains an Azure Key Vault.
	KeyVaultName      string `json:"keyVaultName,omitempty"`      // The name of an Azure Key Vault containing your key.
	KeyIdentifier     string `json:"keyIdentifier,omitempty"`     // The unique identifier of a key in an Azure Key Vault.
	Secret            string `json:"secret,omitempty"`            // The secret associated with the Azure Key Vault specified by azureKeyVault.tenantID.
	TenantID          string `json:"tenantID,omitempty"`          // The unique identifier for an Azure AD tenant within an Azure subscription.
}

// GoogleCloudKms specifies GCP KMS configuration details and whether Encryption at Rest is enabled for an Atlas project.
type GoogleCloudKms struct {
	Enabled              *bool  `json:"enabled,omitempty"`              // Specifies whether Encryption at Rest is enabled for an Atlas project. To disable Encryption at Rest, pass only this parameter with a value of false. When you disable Encryption at Rest, Atlas also removes the configuration details.
	ServiceAccountKey    string `json:"serviceAccountKey,omitempty"`    // String-formatted JSON object containing GCP KMS credentials from your GCP account.
	KeyVersionResourceID string `json:"keyVersionResourceID,omitempty"` // 	The Key Version Resource ID from your GCP account.
}

// Create takes one on-demand snapshot. Atlas takes on-demand snapshots immediately, unlike scheduled snapshots which occur at regular intervals.
//
// See more: https://docs.atlas.mongodb.com/reference/api/enable-configure-encryptionatrest/
func (s *EncryptionsAtRestServiceOp) Create(ctx context.Context, createRequest *EncryptionAtRest) (*EncryptionAtRest, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}
	if createRequest.GroupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}

	path := fmt.Sprintf(encryptionsAtRestBasePath, createRequest.GroupID)
	createRequest.GroupID = ""

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, createRequest)
	if err != nil {
		return nil, nil, err
	}
	root := new(EncryptionAtRest)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root, resp, err
}

// Get retrieves the current configuration for Encryption at Rest for an Atlas project.
//
// See more: https://docs.atlas.mongodb.com/reference/api/get-configuration-encryptionatrest/
func (s *EncryptionsAtRestServiceOp) Get(ctx context.Context, groupID string) (*EncryptionAtRest, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}

	path := fmt.Sprintf(encryptionsAtRestBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(EncryptionAtRest)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete disable the AWS, Azure and Google Encryption at Rest.
//
// See more: https://docs.atlas.mongodb.com/reference/api/enable-configure-encryptionatrest/
func (s *EncryptionsAtRestServiceOp) Delete(ctx context.Context, groupID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}

	createRequest := EncryptionAtRest{
		AwsKms:         AwsKms{Enabled: pointy.Bool(false)},
		AzureKeyVault:  AzureKeyVault{Enabled: pointy.Bool(false)},
		GoogleCloudKms: GoogleCloudKms{Enabled: pointy.Bool(false)},
	}

	path := fmt.Sprintf(encryptionsAtRestBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, createRequest)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
