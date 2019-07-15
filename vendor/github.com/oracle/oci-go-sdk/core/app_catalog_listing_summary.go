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

// AppCatalogListingSummary A summary of a listing.
type AppCatalogListingSummary struct {

	// the region free ocid of the listing resource.
	ListingId *string `mandatory:"false" json:"listingId"`

	// The display name of the listing.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The short summary for the listing.
	Summary *string `mandatory:"false" json:"summary"`

	// The name of the publisher who published this listing.
	PublisherName *string `mandatory:"false" json:"publisherName"`
}

func (m AppCatalogListingSummary) String() string {
	return common.PointerString(m)
}
