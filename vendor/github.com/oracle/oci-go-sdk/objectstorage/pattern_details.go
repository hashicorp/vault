// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object Storage and Archive Storage APIs for managing buckets, objects, and related resources.
// For more information, see Overview of Object Storage (https://docs.cloud.oracle.com/Content/Object/Concepts/objectstorageoverview.htm) and
// Overview of Archive Storage (https://docs.cloud.oracle.com/Content/Archive/Concepts/archivestorageoverview.htm).
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// PatternDetails Specifying inclusion and exclusion patterns.
type PatternDetails struct {

	// An array of glob patterns to match the object names to include. An empty array includes all objects in the
	// bucket. Exclusion patterns take precedence over inclusion patterns.
	// A Glob pattern is a sequence of characters to match text. Any character that appears in the pattern, other
	// than the special pattern characters described below, matches itself.
	//     Glob patterns must be between 1 and 1024 characters.
	//     The special pattern characters have the following meanings:
	//     \           Escapes the following character
	//     *           Matches any string of characters.
	//     ?           Matches any single character .
	//     [...]       Matches a group of characters. A group of characters can be:
	//                     A set of characters, for example: [Zafg9@]. This matches any character in the brackets.
	//                     A range of characters, for example: [a-z]. This matches any character in the range.
	//                         [a-f] is equivalent to [abcdef].
	//                         For character ranges only the CHARACTER-CHARACTER pattern is supported.
	//                             [ab-yz] is not valid
	//                             [a-mn-z] is not valid
	//                         Character ranges can not start with ^ or :
	//                         To include a '-' in the range, make it the first or last character.
	InclusionPatterns []string `mandatory:"false" json:"inclusionPatterns"`

	// An array of glob patterns to match the object names to exclude. An empty array is ignored. Exclusion
	// patterns take precedence over inclusion patterns.
	// A Glob pattern is a sequence of characters to match text. Any character that appears in the pattern, other
	// than the special pattern characters described below, matches itself.
	//     Glob patterns must be between 1 and 1024 characters.
	//     The special pattern characters have the following meanings:
	//     \           Escapes the following character
	//     *           Matches any string of characters.
	//     ?           Matches any single character .
	//     [...]       Matches a group of characters. A group of characters can be:
	//                     A set of characters, for example: [Zafg9@]. This matches any character in the brackets.
	//                     A range of characters, for example: [a-z]. This matches any character in the range.
	//                         [a-f] is equivalent to [abcdef].
	//                         For character ranges only the CHARACTER-CHARACTER pattern is supported.
	//                             [ab-yz] is not valid
	//                             [a-mn-z] is not valid
	//                         Character ranges can not start with ^ or :
	//                         To include a '-' in the range, make it the first or last character.
	ExclusionPatterns []string `mandatory:"false" json:"exclusionPatterns"`
}

func (m PatternDetails) String() string {
	return common.PointerString(m)
}
