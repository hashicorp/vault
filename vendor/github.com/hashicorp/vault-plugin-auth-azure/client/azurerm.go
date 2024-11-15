// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type ComputeClient interface {
	Get(ctx context.Context, resourceGroupName string, vmName string, options *armcompute.VirtualMachinesClientGetOptions) (armcompute.VirtualMachinesClientGetResponse, error)
}

type VMSSClient interface {
	Get(ctx context.Context, resourceGroupName string, vmScaleSetName string, options *armcompute.VirtualMachineScaleSetsClientGetOptions) (armcompute.VirtualMachineScaleSetsClientGetResponse, error)
}

type MSIClient interface {
	Get(ctx context.Context, resourceGroupName string, resourceName string, options *armmsi.UserAssignedIdentitiesClientGetOptions) (armmsi.UserAssignedIdentitiesClientGetResponse, error)
	NewListByResourceGroupPager(resourceGroupName string, options *armmsi.UserAssignedIdentitiesClientListByResourceGroupOptions) *runtime.Pager[armmsi.UserAssignedIdentitiesClientListByResourceGroupResponse]
}

type ResourceClient interface {
	GetByID(ctx context.Context, resourceID, apiVersion string, options *armresources.ClientGetByIDOptions) (armresources.ClientGetByIDResponse, error)
}

type ProvidersClient interface {
	Get(ctx context.Context, resourceProviderNamespace string, options *armresources.ProvidersClientGetOptions) (armresources.ProvidersClientGetResponse, error)
}
