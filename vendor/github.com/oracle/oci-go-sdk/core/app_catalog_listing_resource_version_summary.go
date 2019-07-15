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

// AppCatalogListingResourceVersionSummary Listing Resource Version summary
type AppCatalogListingResourceVersionSummary struct {

	// The OCID of the listing this resource version belongs to.
	ListingId *string `mandatory:"false" json:"listingId"`

	// Date and time the listing resource version was published, in RFC3339 format.
	// Example: `2018-03-20T12:32:53.532Z`
	TimePublished *common.SDKTime `mandatory:"false" json:"timePublished"`

	// OCID of the listing resource.
	ListingResourceId *string `mandatory:"false" json:"listingResourceId"`

	// Resource Version.
	ListingResourceVersion *string `mandatory:"false" json:"listingResourceVersion"`
}

func (m AppCatalogListingResourceVersionSummary) String() string {
	return common.PointerString(m)
}
