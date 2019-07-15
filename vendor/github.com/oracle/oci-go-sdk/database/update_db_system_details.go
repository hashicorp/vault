// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"github.com/oracle/oci-go-sdk/common"
)

// UpdateDbSystemDetails Describes the parameters for updating the DB system.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type UpdateDbSystemDetails struct {

	// The new number of CPU cores to set for the DB system. Not applicable for virtual machine DB systems.
	CpuCoreCount *int `mandatory:"false" json:"cpuCoreCount"`

	Version *PatchDetails `mandatory:"false" json:"version"`

	// The public key portion of the key pair to use for SSH access to the DB system. Multiple public keys can be provided. The length of the combined keys cannot exceed 40,000 characters.
	SshPublicKeys []string `mandatory:"false" json:"sshPublicKeys"`

	// The size, in gigabytes, to scale the attached storage up to for this virtual machine DB system. This value must be greater than current storage size. Note that the resulting total storage size attached will be greater than the amount requested to allow for REDO/RECO space and software volume. Applies only to virtual machine DB systems.
	DataStorageSizeInGBs *int `mandatory:"false" json:"dataStorageSizeInGBs"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The list of Network Security Group OCIDs (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) associated with this DB system.
	// A maximum of 5 allowed.
	NsgIds []string `mandatory:"false" json:"nsgIds"`

	// The list of Network Security Group OCIDs (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) associated with the backup network of this DB system.
	// Applicable only to Exadata DB systems.
	// A maximum of 5 allowed.
	BackupNetworkNsgIds []string `mandatory:"false" json:"backupNetworkNsgIds"`
}

func (m UpdateDbSystemDetails) String() string {
	return common.PointerString(m)
}
