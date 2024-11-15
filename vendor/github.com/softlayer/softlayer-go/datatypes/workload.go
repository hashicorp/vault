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
type Workload_Citrix_Deployment struct {
	Entity

	// The [[SoftLayer_Account]] to which the deployment belongs.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The account ID to which the deployment belongs.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// Topology used for the Citrix Virtual Apps And  Desktop deployment.
	ActiveDirectoryTopology *string `json:"activeDirectoryTopology,omitempty" xmlrpc:"activeDirectoryTopology,omitempty"`

	// The date when this record was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// DataCenter of the deployment.
	DataCenter *string `json:"dataCenter,omitempty" xmlrpc:"dataCenter,omitempty"`

	// It is the unique identifier for the deployment.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The date when this record was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Name of the deployment.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A count of it contains a collection of items under the CVAD deployment.
	ResourceCount *uint `json:"resourceCount,omitempty" xmlrpc:"resourceCount,omitempty"`

	// It contains a collection of items under the CVAD deployment.
	Resources []Workload_Citrix_Deployment_Resource `json:"resources,omitempty" xmlrpc:"resources,omitempty"`

	// Current Status of the CVAD deployment.
	Status *Workload_Citrix_Deployment_Status `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// The [[SoftLayer_Workload_Citrix_Deployment_Status]] of the deployment.
	StatusId *int `json:"statusId,omitempty" xmlrpc:"statusId,omitempty"`

	// It shows if the deployment is for Citrix Hypervisor or VMware.
	Type *Workload_Citrix_Deployment_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The [[SoftLayer_Workload_Citrix_Deployment_Type]] of the deployment.
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`

	// It is the [[SoftLayer_User_Customer]] who placed the order for CVAD.
	User *User_Customer `json:"user,omitempty" xmlrpc:"user,omitempty"`

	// The identifier for the customer who placed the CVAD order.
	UserRecordId *int `json:"userRecordId,omitempty" xmlrpc:"userRecordId,omitempty"`

	// It is the VLAN resource for the CVAD deployment.
	Vlan *Network_Vlan `json:"vlan,omitempty" xmlrpc:"vlan,omitempty"`

	// VLAN ID of the deployment.
	VlanId *int `json:"vlanId,omitempty" xmlrpc:"vlanId,omitempty"`

	// It is an internal identifier for the VMware solution. It gets set if the CVAD order is for VMware.
	VmwareOrderId *string `json:"vmwareOrderId,omitempty" xmlrpc:"vmwareOrderId,omitempty"`
}

// The SoftLayer_Workload_Citrix_Deployment_Resource type contains the information of the resource such as the Deployment ID, resource's Billing Item ID, Order ID and Role of the resource in the CVAD deployment.
type Workload_Citrix_Deployment_Resource struct {
	Entity

	// no documentation yet
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// Billing item ID of the resource
	BillingItemId *int `json:"billingItemId,omitempty" xmlrpc:"billingItemId,omitempty"`

	// The point in time at which the resource was ordered.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	Deployment *Workload_Citrix_Deployment `json:"deployment,omitempty" xmlrpc:"deployment,omitempty"`

	// CVAD Deployment ID of the resource
	DeploymentId *int `json:"deploymentId,omitempty" xmlrpc:"deploymentId,omitempty"`

	// Unique Identifier of the CVAD Deployment Resource
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The last time when the resource was modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// no documentation yet
	Order *Billing_Order `json:"order,omitempty" xmlrpc:"order,omitempty"`

	// Billing Order ID of the resource
	OrderId *int `json:"orderId,omitempty" xmlrpc:"orderId,omitempty"`

	// This flag indicates that whether the CVAD APIs have control over this resource. This resource can be cancelled using CVAD cancellation APIs only if this flag is true.
	OrderedByCvad *bool `json:"orderedByCvad,omitempty" xmlrpc:"orderedByCvad,omitempty"`

	// no documentation yet
	Role *Workload_Citrix_Deployment_Resource_Role `json:"role,omitempty" xmlrpc:"role,omitempty"`

	// Role of the resource within the CVAD deployment. For example, a VSI can have different roles such as Proxy Server or DHCP Server.
	RoleId *int `json:"roleId,omitempty" xmlrpc:"roleId,omitempty"`
}

// The SoftLayer_Workload_Citrix_Deployment_Resource_Response constructs a response object for [[SoftLayer_Workload_Citrix_Deployment_Resource_Response]] for the CVAD resource.
type Workload_Citrix_Deployment_Resource_Response struct {
	Entity

	// Represents the hardware resource of the CVAD deployment.
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// It is a flag for internal usage that represents if the underlying resource is ordered by another system of the same infrastructure provider.
	IsDeploymentOwned *bool `json:"isDeploymentOwned,omitempty" xmlrpc:"isDeploymentOwned,omitempty"`

	// It represents the role of a VSI resource in the CVAD deployment, e.g., a proxy server, DHCP server, cloud connector.
	Role *Workload_Citrix_Deployment_Resource_Role `json:"role,omitempty" xmlrpc:"role,omitempty"`

	// Storage resource for the CVAD deployment.
	Storage *Network_Storage `json:"storage,omitempty" xmlrpc:"storage,omitempty"`

	// Represents the subnet resource of the CVAD deployment.
	Subnet *Network_Subnet `json:"subnet,omitempty" xmlrpc:"subnet,omitempty"`

	// It contains the category of the item which is set for the current response.
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// VSI resource for the CVAD deployment.
	VirtualGuest *Virtual_Guest `json:"virtualGuest,omitempty" xmlrpc:"virtualGuest,omitempty"`

	// Represents the VLAN resource of the CVAD deployment.
	Vlan *Network_Vlan `json:"vlan,omitempty" xmlrpc:"vlan,omitempty"`
}

// SoftLayer_Workload_Citrix_Deployment_Resource_Role contains the role and its description of any resource of Citrix Virtual Apps & Desktops deployment.
type Workload_Citrix_Deployment_Resource_Role struct {
	Entity

	// Description of the resource role
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// ID of the role
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Unique keyName of the role
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// Name of the role
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The SoftLayer_Workload_Citrix_Deployment_Response constructs a response object for the [[SoftLayer_Workload_Citrix_Deployment]] that includes all resources, i.e., [[SoftLayer_Workload_Citrix_Deployment_Resource]].
type Workload_Citrix_Deployment_Response struct {
	Entity

	// The account ID to which the deployment belongs.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// Topology used for the CVAD deployment
	ActiveDirectoryTopology *string `json:"activeDirectoryTopology,omitempty" xmlrpc:"activeDirectoryTopology,omitempty"`

	// The date when this deployment was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// Location name of the deployment.
	DataCenter *string `json:"dataCenter,omitempty" xmlrpc:"dataCenter,omitempty"`

	// ID of the CVAD deployment.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The date when this deployment was modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Name of the deployment.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// It is a collection of objects representing deployment resources such as VLAN, subnet, bare metal, proxy, DHCP, cloud connectors.
	Resources []Workload_Citrix_Deployment_Resource_Response `json:"resources,omitempty" xmlrpc:"resources,omitempty"`

	// Status of the deployment.
	Status *Workload_Citrix_Deployment_Status `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// Represents if the deployment is for Citrix Hypervisor or VMware
	Type *Workload_Citrix_Deployment_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The identifier for the customer who placed the CVAD order.
	UserRecordId *int `json:"userRecordId,omitempty" xmlrpc:"userRecordId,omitempty"`

	// VLAN ID of the deployment.
	VlanId *int `json:"vlanId,omitempty" xmlrpc:"vlanId,omitempty"`

	// It is an internal identifier for the VMware solution. It gets set if the CVAD order is for VMware.
	VmwareOrderId *string `json:"vmwareOrderId,omitempty" xmlrpc:"vmwareOrderId,omitempty"`
}

// The SoftLayer_Workload_Citrix_Deployment_Status shows the status of Citrix Virtual Apps and Desktop deployment. The deployment can be in one of the following statuses at a given point in time: - PROVISIONING: The resources are being provisioned for the deployment. - ACTIVE: All the resources for the deployment are ready. - CANCELLING: Resources of the deployment are being cancelled. - CANCELLED: All the resources of the deployment are cancelled.
type Workload_Citrix_Deployment_Status struct {
	Entity

	// The description of the deployment status.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The ID of the deployment status.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The keyName of the deployment status.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The title of the deployment status.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Workload_Citrix_Deployment_Type struct {
	Entity

	// Description of the deployment type.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The identifer of the deployment type.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// KeyName of the deployment type.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// Name of the deployment type.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
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

	// The subdomain for the ordered hosts (e.g. corp).
	Subdomain *string `json:"subdomain,omitempty" xmlrpc:"subdomain,omitempty"`

	// The vSphere version. Valid values are: "6.7" and "7.0"
	VSphereVersion *string `json:"vSphereVersion,omitempty" xmlrpc:"vSphereVersion,omitempty"`

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
