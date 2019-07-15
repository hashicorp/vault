// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Identity and Access Management Service API
//
// APIs for managing users, groups, compartments, and policies.
//

package identity

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ApiKey A PEM-format RSA credential for securing requests to the Oracle Cloud Infrastructure REST API. Also known
// as an *API signing key*. Specifically, this is the public key from the key pair. The private key remains with
// the user calling the API. For information about generating a key pair
// in the required PEM format, see Required Keys and OCIDs (https://docs.cloud.oracle.com/Content/API/Concepts/apisigningkey.htm).
// **Important:** This is **not** the SSH key for accessing compute instances.
// Each user can have a maximum of three API signing keys.
// For more information about user credentials, see User Credentials (https://docs.cloud.oracle.com/Content/Identity/Concepts/usercredentials.htm).
type ApiKey struct {

	// An Oracle-assigned identifier for the key, in this format:
	// TENANCY_OCID/USER_OCID/KEY_FINGERPRINT.
	KeyId *string `mandatory:"false" json:"keyId"`

	// The key's value.
	KeyValue *string `mandatory:"false" json:"keyValue"`

	// The key's fingerprint (e.g., 12:34:56:78:90:ab:cd:ef:12:34:56:78:90:ab:cd:ef).
	Fingerprint *string `mandatory:"false" json:"fingerprint"`

	// The OCID of the user the key belongs to.
	UserId *string `mandatory:"false" json:"userId"`

	// Date and time the `ApiKey` object was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The API key's current state. After creating an `ApiKey` object, make sure its `lifecycleState` changes from
	// CREATING to ACTIVE before using it.
	LifecycleState ApiKeyLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`
}

func (m ApiKey) String() string {
	return common.PointerString(m)
}

// ApiKeyLifecycleStateEnum Enum with underlying type: string
type ApiKeyLifecycleStateEnum string

// Set of constants representing the allowable values for ApiKeyLifecycleStateEnum
const (
	ApiKeyLifecycleStateCreating ApiKeyLifecycleStateEnum = "CREATING"
	ApiKeyLifecycleStateActive   ApiKeyLifecycleStateEnum = "ACTIVE"
	ApiKeyLifecycleStateInactive ApiKeyLifecycleStateEnum = "INACTIVE"
	ApiKeyLifecycleStateDeleting ApiKeyLifecycleStateEnum = "DELETING"
	ApiKeyLifecycleStateDeleted  ApiKeyLifecycleStateEnum = "DELETED"
)

var mappingApiKeyLifecycleState = map[string]ApiKeyLifecycleStateEnum{
	"CREATING": ApiKeyLifecycleStateCreating,
	"ACTIVE":   ApiKeyLifecycleStateActive,
	"INACTIVE": ApiKeyLifecycleStateInactive,
	"DELETING": ApiKeyLifecycleStateDeleting,
	"DELETED":  ApiKeyLifecycleStateDeleted,
}

// GetApiKeyLifecycleStateEnumValues Enumerates the set of values for ApiKeyLifecycleStateEnum
func GetApiKeyLifecycleStateEnumValues() []ApiKeyLifecycleStateEnum {
	values := make([]ApiKeyLifecycleStateEnum, 0)
	for _, v := range mappingApiKeyLifecycleState {
		values = append(values, v)
	}
	return values
}
