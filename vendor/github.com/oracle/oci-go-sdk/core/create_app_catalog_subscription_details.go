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

// CreateAppCatalogSubscriptionDetails details for creating a subscription for a listing resource version.
type CreateAppCatalogSubscriptionDetails struct {

	// The compartmentID for the subscription.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The OCID of the listing.
	ListingId *string `mandatory:"false" json:"listingId"`

	// Listing resource version.
	ListingResourceVersion *string `mandatory:"false" json:"listingResourceVersion"`

	// Oracle TOU link
	OracleTermsOfUseLink *string `mandatory:"false" json:"oracleTermsOfUseLink"`

	// EULA link
	EulaLink *string `mandatory:"false" json:"eulaLink"`

	// Date and time the agreements were retrieved, in RFC3339 format.
	// Example: `2018-03-20T12:32:53.532Z`
	TimeRetrieved *common.SDKTime `mandatory:"false" json:"timeRetrieved"`

	// A generated signature for this listing resource version retrieved the agreements API.
	Signature *string `mandatory:"false" json:"signature"`
}

func (m CreateAppCatalogSubscriptionDetails) String() string {
	return common.PointerString(m)
}
