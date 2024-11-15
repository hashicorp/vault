//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/resource"
)

// RootResourceID defines the tenant as the root parent of all other ResourceID.
var RootResourceID = resource.RootResourceID

// ResourceID represents a resource ID such as `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg`.
// Don't create this type directly, use ParseResourceID instead.
type ResourceID = resource.ResourceID

// ParseResourceID parses a string to an instance of ResourceID
func ParseResourceID(id string) (*ResourceID, error) {
	return resource.ParseResourceID(id)
}
