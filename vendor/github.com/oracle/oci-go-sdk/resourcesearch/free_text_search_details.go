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

// FreeTextSearchDetails A request containing arbitrary text that must be present in the resource.
type FreeTextSearchDetails struct {

	// The text to search for.
	Text *string `mandatory:"true" json:"text"`

	// The type of matching context returned in the response. If you specify `HIGHLIGHTS`, then the service will highlight fragments in its response. (See ResourceSummary.searchContext and SearchContext for more information.) The default setting is `NONE`.
	MatchingContextType SearchDetailsMatchingContextTypeEnum `mandatory:"false" json:"matchingContextType,omitempty"`
}

//GetMatchingContextType returns MatchingContextType
func (m FreeTextSearchDetails) GetMatchingContextType() SearchDetailsMatchingContextTypeEnum {
	return m.MatchingContextType
}

func (m FreeTextSearchDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m FreeTextSearchDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeFreeTextSearchDetails FreeTextSearchDetails
	s := struct {
		DiscriminatorParam string `json:"type"`
		MarshalTypeFreeTextSearchDetails
	}{
		"FreeText",
		(MarshalTypeFreeTextSearchDetails)(m),
	}

	return json.Marshal(&s)
}
