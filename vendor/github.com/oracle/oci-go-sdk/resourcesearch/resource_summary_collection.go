// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Search Service API
//
// Search for resources in your cloud network.
//

package resourcesearch

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ResourceSummaryCollection Summary representation of resources that matched the search criteria.
type ResourceSummaryCollection struct {

	// A list of resources.
	Items []ResourceSummary `mandatory:"false" json:"items"`
}

func (m ResourceSummaryCollection) String() string {
	return common.PointerString(m)
}
