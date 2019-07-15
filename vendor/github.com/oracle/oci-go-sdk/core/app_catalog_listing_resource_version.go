// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"github.com/oracle/oci-go-sdk/common"
)

// AppCatalogListingResourceVersion Listing Resource Version
type AppCatalogListingResourceVersion struct {

	// The OCID of the listing this resource version belongs to.
	ListingId *string `mandatory:"false" json:"listingId"`

	// Date and time the listing resource version was published, in RFC3339 format.
	// Example: `2018-03-20T12:32:53.532Z`
	TimePublished *common.SDKTime `mandatory:"false" json:"timePublished"`

	// OCID of the listing resource.
	ListingResourceId *string `mandatory:"false" json:"listingResourceId"`

	// Resource Version.
	ListingResourceVersion *string `mandatory:"false" json:"listingResourceVersion"`

	// List of regions that this listing resource version is available.
	// For information about Regions, see
	// Regions (https://docs.cloud.oracle.comGeneral/Concepts/regions.htm).
	// Example: `["us-ashburn-1", "us-phoenix-1"]`
	AvailableRegions []string `mandatory:"false" json:"availableRegions"`

	// Array of shapes compatible with this resource.
	// You may enumerate all available shapes by calling listShapes.
	// Example: `["VM.Standard1.1", "VM.Standard1.2"]`
	CompatibleShapes []string `mandatory:"false" json:"compatibleShapes"`

	// List of accessible ports for instances launched with this listing resource version.
	AccessiblePorts []int `mandatory:"false" json:"accessiblePorts"`

	// Allowed actions for the listing resource.
	AllowedActions []AppCatalogListingResourceVersionAllowedActionsEnum `mandatory:"false" json:"allowedActions,omitempty"`
}

func (m AppCatalogListingResourceVersion) String() string {
	return common.PointerString(m)
}

// AppCatalogListingResourceVersionAllowedActionsEnum Enum with underlying type: string
type AppCatalogListingResourceVersionAllowedActionsEnum string

// Set of constants representing the allowable values for AppCatalogListingResourceVersionAllowedActionsEnum
const (
	AppCatalogListingResourceVersionAllowedActionsSnapshot              AppCatalogListingResourceVersionAllowedActionsEnum = "SNAPSHOT"
	AppCatalogListingResourceVersionAllowedActionsBootVolumeDetach      AppCatalogListingResourceVersionAllowedActionsEnum = "BOOT_VOLUME_DETACH"
	AppCatalogListingResourceVersionAllowedActionsPreserveBootVolume    AppCatalogListingResourceVersionAllowedActionsEnum = "PRESERVE_BOOT_VOLUME"
	AppCatalogListingResourceVersionAllowedActionsSerialConsoleAccess   AppCatalogListingResourceVersionAllowedActionsEnum = "SERIAL_CONSOLE_ACCESS"
	AppCatalogListingResourceVersionAllowedActionsBootRecovery          AppCatalogListingResourceVersionAllowedActionsEnum = "BOOT_RECOVERY"
	AppCatalogListingResourceVersionAllowedActionsBackupBootVolume      AppCatalogListingResourceVersionAllowedActionsEnum = "BACKUP_BOOT_VOLUME"
	AppCatalogListingResourceVersionAllowedActionsCaptureConsoleHistory AppCatalogListingResourceVersionAllowedActionsEnum = "CAPTURE_CONSOLE_HISTORY"
)

var mappingAppCatalogListingResourceVersionAllowedActions = map[string]AppCatalogListingResourceVersionAllowedActionsEnum{
	"SNAPSHOT":                AppCatalogListingResourceVersionAllowedActionsSnapshot,
	"BOOT_VOLUME_DETACH":      AppCatalogListingResourceVersionAllowedActionsBootVolumeDetach,
	"PRESERVE_BOOT_VOLUME":    AppCatalogListingResourceVersionAllowedActionsPreserveBootVolume,
	"SERIAL_CONSOLE_ACCESS":   AppCatalogListingResourceVersionAllowedActionsSerialConsoleAccess,
	"BOOT_RECOVERY":           AppCatalogListingResourceVersionAllowedActionsBootRecovery,
	"BACKUP_BOOT_VOLUME":      AppCatalogListingResourceVersionAllowedActionsBackupBootVolume,
	"CAPTURE_CONSOLE_HISTORY": AppCatalogListingResourceVersionAllowedActionsCaptureConsoleHistory,
}

// GetAppCatalogListingResourceVersionAllowedActionsEnumValues Enumerates the set of values for AppCatalogListingResourceVersionAllowedActionsEnum
func GetAppCatalogListingResourceVersionAllowedActionsEnumValues() []AppCatalogListingResourceVersionAllowedActionsEnum {
	values := make([]AppCatalogListingResourceVersionAllowedActionsEnum, 0)
	for _, v := range mappingAppCatalogListingResourceVersionAllowedActions {
		values = append(values, v)
	}
	return values
}
