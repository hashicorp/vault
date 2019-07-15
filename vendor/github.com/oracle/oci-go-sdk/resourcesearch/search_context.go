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

// SearchContext Contains search context, such as highlighting, for found resources.
type SearchContext struct {

	// Describes what in each field matched the search criteria by showing highlighted values, but only for free text searches or for structured
	// queries that use a MATCHING clause. The list of strings represents fragments of values that matched the query conditions. Highlighted
	// values are wrapped with <hl>..</hl> tags. All values are HTML-encoded (except <hl> tags).
	Highlights map[string][]string `mandatory:"false" json:"highlights"`
}

func (m SearchContext) String() string {
	return common.PointerString(m)
}
