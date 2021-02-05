/**
 * Copyright 2016 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * AUTOMATICALLY GENERATED CODE - DO NOT MODIFY
 */

package datatypes

// no documentation yet
type Workload_Citrix_Client struct {
	Entity
}

// no documentation yet
type Workload_Citrix_Client_Response struct {
	Entity

	// messageId of Citrix account validation response.
	MessageId *string `json:"messageId,omitempty" xmlrpc:"messageId,omitempty"`

	// status of Citrix account validation.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// status message of Citrix account validation.
	StatusMessage *string `json:"statusMessage,omitempty" xmlrpc:"statusMessage,omitempty"`
}

// no documentation yet
type Workload_Citrix_Client_Response_ResourceLocations struct {
	Workload_Citrix_Client_Response

	// no documentation yet
	ResourceLocations []string `json:"resourceLocations,omitempty" xmlrpc:"resourceLocations,omitempty"`
}

// no documentation yet
type Workload_Citrix_Request struct {
	Entity

	// no documentation yet
	ClientId *string `json:"clientId,omitempty" xmlrpc:"clientId,omitempty"`

	// no documentation yet
	ClientSecret *string `json:"clientSecret,omitempty" xmlrpc:"clientSecret,omitempty"`

	// no documentation yet
	CustomerId *string `json:"customerId,omitempty" xmlrpc:"customerId,omitempty"`
}

// no documentation yet
type Workload_Citrix_Request_CreateResourceLocation struct {
	Workload_Citrix_Request

	// no documentation yet
	ResourceLocationName *string `json:"resourceLocationName,omitempty" xmlrpc:"resourceLocationName,omitempty"`
}

// no documentation yet
type Workload_Citrix_Workspace_Order struct {
	Entity
}

// This is the datatype that needs to be populated and sent to SoftLayer_Workload_Citrix_Workspace_Order::placeWorkspaceOrder.
type Workload_Citrix_Workspace_Order_Container struct {
	Entity

	// The active directory domain name
	ActiveDirectoryDomainName *string `json:"activeDirectoryDomainName,omitempty" xmlrpc:"activeDirectoryDomainName,omitempty"`

	// The active directory netbios name (optional)
	ActiveDirectoryNetbiosName *string `json:"activeDirectoryNetbiosName,omitempty" xmlrpc:"activeDirectoryNetbiosName,omitempty"`

	// The active directory safe mode password
	ActiveDirectorySafeModePassword *string `json:"activeDirectorySafeModePassword,omitempty" xmlrpc:"activeDirectorySafeModePassword,omitempty"`

	// The active directory topology
	ActiveDirectoryTopology *string `json:"activeDirectoryTopology,omitempty" xmlrpc:"activeDirectoryTopology,omitempty"`

	// The Citrix API Client Id
	CitrixAPIClientId *string `json:"citrixAPIClientId,omitempty" xmlrpc:"citrixAPIClientId,omitempty"`

	// The Citrix API Client Secret
	CitrixAPIClientSecret *string `json:"citrixAPIClientSecret,omitempty" xmlrpc:"citrixAPIClientSecret,omitempty"`

	// The Citrix customer id
	CitrixCustomerId *string `json:"citrixCustomerId,omitempty" xmlrpc:"citrixCustomerId,omitempty"`

	// The Citrix resource location name
	CitrixResourceLocationName *string `json:"citrixResourceLocationName,omitempty" xmlrpc:"citrixResourceLocationName,omitempty"`

	// The default domain to be used for all server orders where the domain is not specified.
	Domain *string `json:"domain,omitempty" xmlrpc:"domain,omitempty"`

	// The specific [[SoftLayer_Location_Datacenter]] id where the order should be provisioned.
	Location *string `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// There should be one child orderContainer for each component ordered.  The containerIdentifier should be set on each and have these exact values: proxy server, bare metal server with hypervisor, dhcp server, citrix connector servers, active directory server, vlan, subnet, storage
	OrderContainers []Container_Product_Order `json:"orderContainers,omitempty" xmlrpc:"orderContainers,omitempty"`

	// Set this value to order IBM Cloud for VMware Solutions servers as part of your Citrix Virtual Apps and Desktops order
	VmwareContainer *Workload_Citrix_Workspace_Order_VMwareContainer `json:"vmwareContainer,omitempty" xmlrpc:"vmwareContainer,omitempty"`
}

// This is the datatype that can be populated by the customer to provide license key information for VMware orders.
type Workload_Citrix_Workspace_Order_LicenseKey struct {
	Entity

	// The license key
	Key *string `json:"key,omitempty" xmlrpc:"key,omitempty"`

	// The name of the product (e.g. vcenter, nsx, vsphere, vsan)
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The license type
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// This is the datatype that can be populated by the customer to provide NFS shared storage information for VMware orders.
type Workload_Citrix_Workspace_Order_SharedStorage struct {
	Entity

	// Which storage tier: e.g. READHEAVY_TIER
	Iops *string `json:"iops,omitempty" xmlrpc:"iops,omitempty"`

	// The number of shared storages to order
	Quantity *int `json:"quantity,omitempty" xmlrpc:"quantity,omitempty"`

	// The size of the storage (e.g. STORAGE_SPACE_FOR_2_IOPS_PER_GB)
	Size *string `json:"size,omitempty" xmlrpc:"size,omitempty"`

	// The volume
	Volume *int `json:"volume,omitempty" xmlrpc:"volume,omitempty"`
}

// This is the datatype that needs to be populated and sent to SoftLayer_Workload_Citrix_Workspace_Order::placeWorkspaceOrder to order and provision one or more VMware server instances to be used with Citrix Virtual Apps and Desktops.
type Workload_Citrix_Workspace_Order_VMwareContainer struct {
	Entity

	// The bare metal disks
	Disks []string `json:"disks,omitempty" xmlrpc:"disks,omitempty"`

	// The domain for the ordered hosts (e.g. example.org)
	Domain *string `json:"domain,omitempty" xmlrpc:"domain,omitempty"`

	// Customer provided license keys (optional)
	LicenseKeys []Workload_Citrix_Workspace_Order_LicenseKey `json:"licenseKeys,omitempty" xmlrpc:"licenseKeys,omitempty"`

	// The datacenter location
	Location *string `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// The name associated with the order
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The nickname for the vSRX service
	Nickname *string `json:"nickname,omitempty" xmlrpc:"nickname,omitempty"`

	// The number of instances to order
	Quantity *int `json:"quantity,omitempty" xmlrpc:"quantity,omitempty"`

	// The bare metal ram type
	Ram *string `json:"ram,omitempty" xmlrpc:"ram,omitempty"`

	// The bare metal server type
	Server *string `json:"server,omitempty" xmlrpc:"server,omitempty"`

	// The bare metal shared nfs storage (optional)
	SharedStorage []Workload_Citrix_Workspace_Order_SharedStorage `json:"sharedStorage,omitempty" xmlrpc:"sharedStorage,omitempty"`

	// The subdomain for the ordered hosts (e.g. corp)
	Subdomain *string `json:"subdomain,omitempty" xmlrpc:"subdomain,omitempty"`

	// The bare metal vsan cache disks (optional)
	VsanCacheDisks []string `json:"vsanCacheDisks,omitempty" xmlrpc:"vsanCacheDisks,omitempty"`
}

// no documentation yet
type Workload_Citrix_Workspace_Response struct {
	Entity

	// messageId associated with any error
	MessageId *string `json:"messageId,omitempty" xmlrpc:"messageId,omitempty"`

	// status of service methods
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// status message
	StatusMessage *string `json:"statusMessage,omitempty" xmlrpc:"statusMessage,omitempty"`
}

// no documentation yet
type Workload_Citrix_Workspace_Response_Item struct {
	Workload_Citrix_Workspace_Response

	// the id of the resource (HARDWARE, GUEST, VLAN, SUBNET, VMWARE)
	Id *string `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// the name associated with the resource (e.g. name, hostname)
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// the type of resource (HARDWARE, GUEST, NETWORK_VLAN, SUBNET)
	TypeName *string `json:"typeName,omitempty" xmlrpc:"typeName,omitempty"`
}

// no documentation yet
type Workload_Citrix_Workspace_Response_Result struct {
	Workload_Citrix_Workspace_Response

	// identification and operation result for each item
	Items []Workload_Citrix_Workspace_Response `json:"items,omitempty" xmlrpc:"items,omitempty"`
}
