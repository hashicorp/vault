//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/internal/resource"
)

// SubscriptionResourceType is the ResourceType of a subscription
var SubscriptionResourceType = resource.SubscriptionResourceType

// ResourceGroupResourceType is the ResourceType of a resource group
var ResourceGroupResourceType = resource.ResourceGroupResourceType

// TenantResourceType is the ResourceType of a tenant
var TenantResourceType = resource.TenantResourceType

// ProviderResourceType is the ResourceType of a provider
var ProviderResourceType = resource.ProviderResourceType

// ResourceType represents an Azure resource type, e.g. "Microsoft.Network/virtualNetworks/subnets".
// Don't create this type directly, use ParseResourceType or NewResourceType instead.
type ResourceType = resource.ResourceType

// NewResourceType creates an instance of ResourceType using a provider namespace
// such as "Microsoft.Network" and type such as "virtualNetworks/subnets".
func NewResourceType(providerNamespace, typeName string) ResourceType {
	return resource.NewResourceType(providerNamespace, typeName)
}

// ParseResourceType parses the ResourceType from a resource type string (e.g. Microsoft.Network/virtualNetworks/subsets)
// or a resource identifier string.
// e.g. /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myRg/providers/Microsoft.Network/virtualNetworks/vnet/subnets/mySubnet)
func ParseResourceType(resourceIDOrType string) (ResourceType, error) {
	return resource.ParseResourceType(resourceIDOrType)
}
