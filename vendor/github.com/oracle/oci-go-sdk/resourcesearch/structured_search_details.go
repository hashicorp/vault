// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Search Service API
//
// Search for resources in your cloud network.
//

package resourcesearch

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// StructuredSearchDetails A request containing search filters using the structured search query language.
type StructuredSearchDetails struct {

	// The structured query describing which resources to search for.
	Query *string `mandatory:"true" json:"query"`

	// The type of matching context returned in the response. If you specify `HIGHLIGHTS`, then the service will highlight fragments in its response. (See ResourceSummary.searchContext and SearchContext for more information.) The default setting is `NONE`.
	MatchingContextType SearchDetailsMatchingContextTypeEnum `mandatory:"false" json:"matchingContextType,omitempty"`
}

//GetMatchingContextType returns MatchingContextType
func (m StructuredSearchDetails) GetMatchingContextType() SearchDetailsMatchingContextTypeEnum {
	return m.MatchingContextType
}

func (m StructuredSearchDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m StructuredSearchDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeStructuredSearchDetails StructuredSearchDetails
	s := struct {
		DiscriminatorParam string `json:"type"`
		MarshalTypeStructuredSearchDetails
	}{
		"Structured",
		(MarshalTypeStructuredSearchDetails)(m),
	}

	return json.Marshal(&s)
}
