/**
 * Copyright 2016-2024 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed
 * on an "AS IS" BASIS,WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

// AUTOMATICALLY GENERATED CODE - DO NOT MODIFY

package datatypes

// no documentation yet
type Network struct {
	Entity
}

// The SoftLayer_Network_Application_Delivery_Controller data type models a single instance of an application delivery controller. Local properties are read only, except for a ”notes” property, which can be used to describe your application delivery controller service. The type's relational properties provide more information to the service's function and login information to the controller's backend management if advanced view is enabled.
type Network_Application_Delivery_Controller struct {
	Entity

	// The SoftLayer customer account that owns an application delivery controller record.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The unique identifier of the SoftLayer customer account that owns an application delivery controller record
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The average daily public bandwidth usage for the current billing cycle.
	AverageDailyPublicBandwidthUsage *Float64 `json:"averageDailyPublicBandwidthUsage,omitempty" xmlrpc:"averageDailyPublicBandwidthUsage,omitempty"`

	// The billing item for a Application Delivery Controller.
	BillingItem *Billing_Item_Network_Application_Delivery_Controller `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// Previous configurations for an Application Delivery Controller.
	ConfigurationHistory []Network_Application_Delivery_Controller_Configuration_History `json:"configurationHistory,omitempty" xmlrpc:"configurationHistory,omitempty"`

	// A count of previous configurations for an Application Delivery Controller.
	ConfigurationHistoryCount *uint `json:"configurationHistoryCount,omitempty" xmlrpc:"configurationHistoryCount,omitempty"`

	// The date that an application delivery controller record was created
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The datacenter that the application delivery controller resides in.
	Datacenter *Location `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// A brief description of an application delivery controller record.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// An application delivery controller's unique identifier
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The total public inbound bandwidth for the current billing cycle.
	InboundPublicBandwidthUsage *Float64 `json:"inboundPublicBandwidthUsage,omitempty" xmlrpc:"inboundPublicBandwidthUsage,omitempty"`

	// The date in which the license for this application delivery controller will expire.
	LicenseExpirationDate *Time `json:"licenseExpirationDate,omitempty" xmlrpc:"licenseExpirationDate,omitempty"`

	// A count of the virtual IP address records that belong to an application delivery controller based load balancer.
	LoadBalancerCount *uint `json:"loadBalancerCount,omitempty" xmlrpc:"loadBalancerCount,omitempty"`

	// The virtual IP address records that belong to an application delivery controller based load balancer.
	LoadBalancers []Network_LoadBalancer_VirtualIpAddress `json:"loadBalancers,omitempty" xmlrpc:"loadBalancers,omitempty"`

	// A flag indicating that this Application Delivery Controller is a managed resource.
	ManagedResourceFlag *bool `json:"managedResourceFlag,omitempty" xmlrpc:"managedResourceFlag,omitempty"`

	// An application delivery controller's management ip address.
	ManagementIpAddress *string `json:"managementIpAddress,omitempty" xmlrpc:"managementIpAddress,omitempty"`

	// The date that an application delivery controller record was last modified
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// An application delivery controller's name
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The network VLAN that an application delivery controller resides on.
	NetworkVlan *Network_Vlan `json:"networkVlan,omitempty" xmlrpc:"networkVlan,omitempty"`

	// A count of the network VLANs that an application delivery controller resides on.
	NetworkVlanCount *uint `json:"networkVlanCount,omitempty" xmlrpc:"networkVlanCount,omitempty"`

	// The network VLANs that an application delivery controller resides on.
	NetworkVlans []Network_Vlan `json:"networkVlans,omitempty" xmlrpc:"networkVlans,omitempty"`

	// Editable notes used to describe an application delivery controller's function
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// The total public outbound bandwidth for the current billing cycle.
	OutboundPublicBandwidthUsage *Float64 `json:"outboundPublicBandwidthUsage,omitempty" xmlrpc:"outboundPublicBandwidthUsage,omitempty"`

	// The password used to connect to an application delivery controller's management interface when it is operating in advanced view mode.
	Password *Software_Component_Password `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// An application delivery controller's primary public IP address.
	PrimaryIpAddress *string `json:"primaryIpAddress,omitempty" xmlrpc:"primaryIpAddress,omitempty"`

	// The projected public outbound bandwidth for the current billing cycle.
	ProjectedPublicBandwidthUsage *Float64 `json:"projectedPublicBandwidthUsage,omitempty" xmlrpc:"projectedPublicBandwidthUsage,omitempty"`

	// A count of a network application controller's subnets. A subnet is a group of IP addresses
	SubnetCount *uint `json:"subnetCount,omitempty" xmlrpc:"subnetCount,omitempty"`

	// A network application controller's subnets. A subnet is a group of IP addresses
	Subnets []Network_Subnet `json:"subnets,omitempty" xmlrpc:"subnets,omitempty"`

	// A count of
	TagReferenceCount *uint `json:"tagReferenceCount,omitempty" xmlrpc:"tagReferenceCount,omitempty"`

	// no documentation yet
	TagReferences []Tag_Reference `json:"tagReferences,omitempty" xmlrpc:"tagReferences,omitempty"`

	// no documentation yet
	Type *Network_Application_Delivery_Controller_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// no documentation yet
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`

	// A count of
	VirtualIpAddressCount *uint `json:"virtualIpAddressCount,omitempty" xmlrpc:"virtualIpAddressCount,omitempty"`

	// no documentation yet
	VirtualIpAddresses []Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress `json:"virtualIpAddresses,omitempty" xmlrpc:"virtualIpAddresses,omitempty"`
}

// The SoftLayer_Network_Application_Delivery_Controller_Configuration_History data type models a single instance of a configuration history entry for an application delivery controller. The configuration history entries are used to support creating backups of an application delivery controller's configuration state in order to restore them later if needed.
type Network_Application_Delivery_Controller_Configuration_History struct {
	Entity

	// The application delivery controller that a configuration history record belongs to.
	Controller *Network_Application_Delivery_Controller `json:"controller,omitempty" xmlrpc:"controller,omitempty"`

	// The date a configuration history record was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// An configuration history record's unique identifier
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Editable notes used to describe a configuration history record
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_Health_Attribute struct {
	Entity

	// no documentation yet
	HealthAttributeTypeId *int `json:"healthAttributeTypeId,omitempty" xmlrpc:"healthAttributeTypeId,omitempty"`

	// no documentation yet
	HealthCheck *Network_Application_Delivery_Controller_LoadBalancer_Health_Check `json:"healthCheck,omitempty" xmlrpc:"healthCheck,omitempty"`

	// no documentation yet
	HealthCheckId *int `json:"healthCheckId,omitempty" xmlrpc:"healthCheckId,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Type *Network_Application_Delivery_Controller_LoadBalancer_Health_Attribute_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// no documentation yet
	Value *string `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_Health_Attribute_Type struct {
	Entity

	// no documentation yet
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	ValueExpression *string `json:"valueExpression,omitempty" xmlrpc:"valueExpression,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_Health_Check struct {
	Entity

	// A count of
	AttributeCount *uint `json:"attributeCount,omitempty" xmlrpc:"attributeCount,omitempty"`

	// no documentation yet
	Attributes []Network_Application_Delivery_Controller_LoadBalancer_Health_Attribute `json:"attributes,omitempty" xmlrpc:"attributes,omitempty"`

	// no documentation yet
	HealthCheckTypeId *int `json:"healthCheckTypeId,omitempty" xmlrpc:"healthCheckTypeId,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// A count of
	ServiceCount *uint `json:"serviceCount,omitempty" xmlrpc:"serviceCount,omitempty"`

	// no documentation yet
	Services []Network_Application_Delivery_Controller_LoadBalancer_Service `json:"services,omitempty" xmlrpc:"services,omitempty"`

	// no documentation yet
	Type *Network_Application_Delivery_Controller_LoadBalancer_Health_Check_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_Health_Check_Type struct {
	Entity

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_Routing_Method struct {
	Entity

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_Routing_Type struct {
	Entity

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_Service struct {
	Entity

	// no documentation yet
	Enabled *int `json:"enabled,omitempty" xmlrpc:"enabled,omitempty"`

	// A count of
	GroupCount *uint `json:"groupCount,omitempty" xmlrpc:"groupCount,omitempty"`

	// A count of
	GroupReferenceCount *uint `json:"groupReferenceCount,omitempty" xmlrpc:"groupReferenceCount,omitempty"`

	// no documentation yet
	GroupReferences []Network_Application_Delivery_Controller_LoadBalancer_Service_Group_CrossReference `json:"groupReferences,omitempty" xmlrpc:"groupReferences,omitempty"`

	// no documentation yet
	Groups []Network_Application_Delivery_Controller_LoadBalancer_Service_Group `json:"groups,omitempty" xmlrpc:"groups,omitempty"`

	// no documentation yet
	HealthCheck *Network_Application_Delivery_Controller_LoadBalancer_Health_Check `json:"healthCheck,omitempty" xmlrpc:"healthCheck,omitempty"`

	// A count of
	HealthCheckCount *uint `json:"healthCheckCount,omitempty" xmlrpc:"healthCheckCount,omitempty"`

	// no documentation yet
	HealthChecks []Network_Application_Delivery_Controller_LoadBalancer_Health_Check `json:"healthChecks,omitempty" xmlrpc:"healthChecks,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	IpAddress *Network_Subnet_IpAddress `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// no documentation yet
	IpAddressId *int `json:"ipAddressId,omitempty" xmlrpc:"ipAddressId,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// no documentation yet
	Port *int `json:"port,omitempty" xmlrpc:"port,omitempty"`

	// no documentation yet
	ServiceGroup *Network_Application_Delivery_Controller_LoadBalancer_Service_Group `json:"serviceGroup,omitempty" xmlrpc:"serviceGroup,omitempty"`

	// no documentation yet
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_Service_Group struct {
	Entity

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// no documentation yet
	RoutingMethod *Network_Application_Delivery_Controller_LoadBalancer_Routing_Method `json:"routingMethod,omitempty" xmlrpc:"routingMethod,omitempty"`

	// no documentation yet
	RoutingMethodId *int `json:"routingMethodId,omitempty" xmlrpc:"routingMethodId,omitempty"`

	// no documentation yet
	RoutingType *Network_Application_Delivery_Controller_LoadBalancer_Routing_Type `json:"routingType,omitempty" xmlrpc:"routingType,omitempty"`

	// no documentation yet
	RoutingTypeId *int `json:"routingTypeId,omitempty" xmlrpc:"routingTypeId,omitempty"`

	// A count of
	ServiceCount *uint `json:"serviceCount,omitempty" xmlrpc:"serviceCount,omitempty"`

	// A count of
	ServiceReferenceCount *uint `json:"serviceReferenceCount,omitempty" xmlrpc:"serviceReferenceCount,omitempty"`

	// no documentation yet
	ServiceReferences []Network_Application_Delivery_Controller_LoadBalancer_Service_Group_CrossReference `json:"serviceReferences,omitempty" xmlrpc:"serviceReferences,omitempty"`

	// no documentation yet
	Services []Network_Application_Delivery_Controller_LoadBalancer_Service `json:"services,omitempty" xmlrpc:"services,omitempty"`

	// The timeout value for connections from remote clients to the load balancer. Timeout values are only valid for HTTP service groups.
	Timeout *int `json:"timeout,omitempty" xmlrpc:"timeout,omitempty"`

	// no documentation yet
	VirtualServer *Network_Application_Delivery_Controller_LoadBalancer_VirtualServer `json:"virtualServer,omitempty" xmlrpc:"virtualServer,omitempty"`

	// A count of
	VirtualServerCount *uint `json:"virtualServerCount,omitempty" xmlrpc:"virtualServerCount,omitempty"`

	// no documentation yet
	VirtualServers []Network_Application_Delivery_Controller_LoadBalancer_VirtualServer `json:"virtualServers,omitempty" xmlrpc:"virtualServers,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_Service_Group_CrossReference struct {
	Entity

	// no documentation yet
	Service *Network_Application_Delivery_Controller_LoadBalancer_Service `json:"service,omitempty" xmlrpc:"service,omitempty"`

	// no documentation yet
	ServiceGroup *Network_Application_Delivery_Controller_LoadBalancer_Service_Group `json:"serviceGroup,omitempty" xmlrpc:"serviceGroup,omitempty"`

	// no documentation yet
	ServiceGroupId *int `json:"serviceGroupId,omitempty" xmlrpc:"serviceGroupId,omitempty"`

	// no documentation yet
	ServiceId *int `json:"serviceId,omitempty" xmlrpc:"serviceId,omitempty"`

	// no documentation yet
	Weight *int `json:"weight,omitempty" xmlrpc:"weight,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress struct {
	Entity

	// no documentation yet
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The unique identifier of the SoftLayer customer account that owns the virtual IP address
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// A virtual IP address's associated application delivery controller.
	ApplicationDeliveryController *Network_Application_Delivery_Controller `json:"applicationDeliveryController,omitempty" xmlrpc:"applicationDeliveryController,omitempty"`

	// A count of a virtual IP address's associated application delivery controllers.
	ApplicationDeliveryControllerCount *uint `json:"applicationDeliveryControllerCount,omitempty" xmlrpc:"applicationDeliveryControllerCount,omitempty"`

	// A virtual IP address's associated application delivery controllers.
	ApplicationDeliveryControllers []Network_Application_Delivery_Controller `json:"applicationDeliveryControllers,omitempty" xmlrpc:"applicationDeliveryControllers,omitempty"`

	// The current billing item for the load balancer virtual IP. This is only valid when dedicatedFlag is false. This is an independent virtual IP, and if canceled, will only affect the associated virtual IP.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// The connection limit for this virtual IP address
	ConnectionLimit *int `json:"connectionLimit,omitempty" xmlrpc:"connectionLimit,omitempty"`

	// The units for the connection limit
	ConnectionLimitUnits *string `json:"connectionLimitUnits,omitempty" xmlrpc:"connectionLimitUnits,omitempty"`

	// The current billing item for the load balancing device housing the virtual IP. This billing item represents a device which could contain other virtual IPs. Caution should be taken when canceling. This is only valid when dedicatedFlag is true.
	DedicatedBillingItem *Billing_Item_Network_LoadBalancer `json:"dedicatedBillingItem,omitempty" xmlrpc:"dedicatedBillingItem,omitempty"`

	// A flag that determines if a VIP is dedicated or not. This is used to override the connection limit and use an unlimited value.
	DedicatedFlag *bool `json:"dedicatedFlag,omitempty" xmlrpc:"dedicatedFlag,omitempty"`

	// Denotes whether the virtual IP is configured within a high availability cluster.
	HighAvailabilityFlag *bool `json:"highAvailabilityFlag,omitempty" xmlrpc:"highAvailabilityFlag,omitempty"`

	// The unique identifier of the virtual IP address record
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	IpAddress *Network_Subnet_IpAddress `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// ID of the IP address this virtual IP utilizes
	IpAddressId *int `json:"ipAddressId,omitempty" xmlrpc:"ipAddressId,omitempty"`

	// no documentation yet
	LoadBalancerHardware []Hardware `json:"loadBalancerHardware,omitempty" xmlrpc:"loadBalancerHardware,omitempty"`

	// A count of
	LoadBalancerHardwareCount *uint `json:"loadBalancerHardwareCount,omitempty" xmlrpc:"loadBalancerHardwareCount,omitempty"`

	// A flag indicating that the load balancer is a managed resource.
	ManagedResourceFlag *bool `json:"managedResourceFlag,omitempty" xmlrpc:"managedResourceFlag,omitempty"`

	// User-created notes for this load balancer virtual IP address
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// A count of the list of security ciphers enabled for this virtual IP address
	SecureTransportCipherCount *uint `json:"secureTransportCipherCount,omitempty" xmlrpc:"secureTransportCipherCount,omitempty"`

	// The list of security ciphers enabled for this virtual IP address
	SecureTransportCiphers []Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress_SecureTransportCipher `json:"secureTransportCiphers,omitempty" xmlrpc:"secureTransportCiphers,omitempty"`

	// A count of the list of secure transport protocols enabled for this virtual IP address
	SecureTransportProtocolCount *uint `json:"secureTransportProtocolCount,omitempty" xmlrpc:"secureTransportProtocolCount,omitempty"`

	// The list of secure transport protocols enabled for this virtual IP address
	SecureTransportProtocols []Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress_SecureTransportProtocol `json:"secureTransportProtocols,omitempty" xmlrpc:"secureTransportProtocols,omitempty"`

	// The SSL certificate currently associated with the VIP.
	SecurityCertificate *Security_Certificate `json:"securityCertificate,omitempty" xmlrpc:"securityCertificate,omitempty"`

	// The SSL certificate currently associated with the VIP. Provides chosen certificate visibility to unprivileged users.
	SecurityCertificateEntry *Security_Certificate_Entry `json:"securityCertificateEntry,omitempty" xmlrpc:"securityCertificateEntry,omitempty"`

	// The unique identifier of the Security Certificate to be utilized when SSL support is enabled.
	SecurityCertificateId *int `json:"securityCertificateId,omitempty" xmlrpc:"securityCertificateId,omitempty"`

	// Determines if the VIP currently has SSL acceleration enabled
	SslActiveFlag *bool `json:"sslActiveFlag,omitempty" xmlrpc:"sslActiveFlag,omitempty"`

	// Determines if the VIP is _allowed_ to utilize SSL acceleration
	SslEnabledFlag *bool `json:"sslEnabledFlag,omitempty" xmlrpc:"sslEnabledFlag,omitempty"`

	// A count of
	VirtualServerCount *uint `json:"virtualServerCount,omitempty" xmlrpc:"virtualServerCount,omitempty"`

	// no documentation yet
	VirtualServers []Network_Application_Delivery_Controller_LoadBalancer_VirtualServer `json:"virtualServers,omitempty" xmlrpc:"virtualServers,omitempty"`
}

// A single cipher configured for a load balancer virtual IP address instance. Instances of this class are immutable and should reflect a cipher that is configurable on a load balancer device.
type Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress_SecureTransportCipher struct {
	Entity

	// Unique identifier for the cipher instance
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Identifier for the associated encryption algorithm
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	VirtualIpAddress *Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress `json:"virtualIpAddress,omitempty" xmlrpc:"virtualIpAddress,omitempty"`

	// Identifier for the associated [[SoftLayer_Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress (type)|virtual IP address]] instance
	VirtualIpAddressId *int `json:"virtualIpAddressId,omitempty" xmlrpc:"virtualIpAddressId,omitempty"`
}

// Links a SSL transport protocol with a virtual IP address instance. Instances of this class are immutable and should reflect a protocol that is configurable on a load balancer device.
type Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress_SecureTransportProtocol struct {
	Entity

	// Unique identifier for the protocol instance
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Identifier for the associated communication protocol
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	VirtualIpAddress *Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress `json:"virtualIpAddress,omitempty" xmlrpc:"virtualIpAddress,omitempty"`

	// Identifier for the associated [[SoftLayer_Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress (type)|virtual IP address]] instance
	VirtualIpAddressId *int `json:"virtualIpAddressId,omitempty" xmlrpc:"virtualIpAddressId,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_LoadBalancer_VirtualServer struct {
	Entity

	// no documentation yet
	Allocation *int `json:"allocation,omitempty" xmlrpc:"allocation,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// no documentation yet
	Port *int `json:"port,omitempty" xmlrpc:"port,omitempty"`

	// no documentation yet
	RoutingMethod *Network_Application_Delivery_Controller_LoadBalancer_Routing_Method `json:"routingMethod,omitempty" xmlrpc:"routingMethod,omitempty"`

	// no documentation yet
	RoutingMethodId *int `json:"routingMethodId,omitempty" xmlrpc:"routingMethodId,omitempty"`

	// A count of
	ServiceGroupCount *uint `json:"serviceGroupCount,omitempty" xmlrpc:"serviceGroupCount,omitempty"`

	// no documentation yet
	ServiceGroups []Network_Application_Delivery_Controller_LoadBalancer_Service_Group `json:"serviceGroups,omitempty" xmlrpc:"serviceGroups,omitempty"`

	// no documentation yet
	VirtualIpAddress *Network_Application_Delivery_Controller_LoadBalancer_VirtualIpAddress `json:"virtualIpAddress,omitempty" xmlrpc:"virtualIpAddress,omitempty"`

	// no documentation yet
	VirtualIpAddressId *int `json:"virtualIpAddressId,omitempty" xmlrpc:"virtualIpAddressId,omitempty"`
}

// no documentation yet
type Network_Application_Delivery_Controller_Type struct {
	Entity

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The SoftLayer_Network_Bandwidth_Usage data type contains specific information relating to bandwidth utilization at a specific point in time on a given network interface.
type Network_Bandwidth_Usage struct {
	Entity

	// Incoming bandwidth utilization.
	AmountIn *Float64 `json:"amountIn,omitempty" xmlrpc:"amountIn,omitempty"`

	// Outgoing bandwidth utilization.
	AmountOut *Float64 `json:"amountOut,omitempty" xmlrpc:"amountOut,omitempty"`

	// ID of the bandwidth usage detail type for this record.
	BandwidthUsageDetailTypeId *Float64 `json:"bandwidthUsageDetailTypeId,omitempty" xmlrpc:"bandwidthUsageDetailTypeId,omitempty"`

	// The tracking object this bandwidth usage record describes.
	TrackingObject *Metric_Tracking_Object `json:"trackingObject,omitempty" xmlrpc:"trackingObject,omitempty"`

	// In and out bandwidth utilization for a specified time stamp.
	Type *Network_Bandwidth_Version1_Usage_Detail_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// The SoftLayer_Network_Bandwidth_Version1_Allocation data type contains general information relating to a single bandwidth allocation record.
type Network_Bandwidth_Version1_Allocation struct {
	Entity

	// A bandwidth allotment detail.
	AllotmentDetail *Network_Bandwidth_Version1_Allotment_Detail `json:"allotmentDetail,omitempty" xmlrpc:"allotmentDetail,omitempty"`

	// The amount of bandwidth allocated.
	Amount *Float64 `json:"amount,omitempty" xmlrpc:"amount,omitempty"`

	// Billing item associated with this hardware allocation.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// Internal ID associated with this allocation.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`
}

// The SoftLayer_Network_Bandwidth_Version1_Allotment class provides methods and data structures necessary to work with an array of hardware objects associated with a single Bandwidth Pooling.
type Network_Bandwidth_Version1_Allotment struct {
	Entity

	// The account associated with this virtual rack.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The user account identifier associated with this allotment.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// A count of the bandwidth allotment detail records associated with this virtual rack.
	ActiveDetailCount *uint `json:"activeDetailCount,omitempty" xmlrpc:"activeDetailCount,omitempty"`

	// The bandwidth allotment detail records associated with this virtual rack.
	ActiveDetails []Network_Bandwidth_Version1_Allotment_Detail `json:"activeDetails,omitempty" xmlrpc:"activeDetails,omitempty"`

	// A count of the Application Delivery Controller contained within a virtual rack.
	ApplicationDeliveryControllerCount *uint `json:"applicationDeliveryControllerCount,omitempty" xmlrpc:"applicationDeliveryControllerCount,omitempty"`

	// The Application Delivery Controller contained within a virtual rack.
	ApplicationDeliveryControllers []Network_Application_Delivery_Controller `json:"applicationDeliveryControllers,omitempty" xmlrpc:"applicationDeliveryControllers,omitempty"`

	// The average daily public bandwidth usage for the current billing cycle.
	AverageDailyPublicBandwidthUsage *Float64 `json:"averageDailyPublicBandwidthUsage,omitempty" xmlrpc:"averageDailyPublicBandwidthUsage,omitempty"`

	// The bandwidth allotment type of this virtual rack.
	BandwidthAllotmentType *Network_Bandwidth_Version1_Allotment_Type `json:"bandwidthAllotmentType,omitempty" xmlrpc:"bandwidthAllotmentType,omitempty"`

	// An identifier marking this allotment as a virtual private rack (1) or a bandwidth pooling(2).
	BandwidthAllotmentTypeId *int `json:"bandwidthAllotmentTypeId,omitempty" xmlrpc:"bandwidthAllotmentTypeId,omitempty"`

	// A count of the bare metal server instances contained within a virtual rack.
	BareMetalInstanceCount *uint `json:"bareMetalInstanceCount,omitempty" xmlrpc:"bareMetalInstanceCount,omitempty"`

	// The bare metal server instances contained within a virtual rack.
	BareMetalInstances []Hardware `json:"bareMetalInstances,omitempty" xmlrpc:"bareMetalInstances,omitempty"`

	// A virtual rack's raw bandwidth usage data for an account's current billing cycle. One object is returned for each network this server is attached to.
	BillingCycleBandwidthUsage []Network_Bandwidth_Usage `json:"billingCycleBandwidthUsage,omitempty" xmlrpc:"billingCycleBandwidthUsage,omitempty"`

	// A count of a virtual rack's raw bandwidth usage data for an account's current billing cycle. One object is returned for each network this server is attached to.
	BillingCycleBandwidthUsageCount *uint `json:"billingCycleBandwidthUsageCount,omitempty" xmlrpc:"billingCycleBandwidthUsageCount,omitempty"`

	// A virtual rack's raw private network bandwidth usage data for an account's current billing cycle.
	BillingCyclePrivateBandwidthUsage *Network_Bandwidth_Usage `json:"billingCyclePrivateBandwidthUsage,omitempty" xmlrpc:"billingCyclePrivateBandwidthUsage,omitempty"`

	// A virtual rack's raw public network bandwidth usage data for an account's current billing cycle.
	BillingCyclePublicBandwidthUsage *Network_Bandwidth_Usage `json:"billingCyclePublicBandwidthUsage,omitempty" xmlrpc:"billingCyclePublicBandwidthUsage,omitempty"`

	// The total public bandwidth used in this virtual rack for an account's current billing cycle.
	BillingCyclePublicUsageTotal *uint `json:"billingCyclePublicUsageTotal,omitempty" xmlrpc:"billingCyclePublicUsageTotal,omitempty"`

	// A virtual rack's billing item.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// Creation date for an allotment.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// An object that provides commonly used bandwidth summary components for the current billing cycle.
	CurrentBandwidthSummary *Metric_Tracking_Object_Bandwidth_Summary `json:"currentBandwidthSummary,omitempty" xmlrpc:"currentBandwidthSummary,omitempty"`

	// A count of the bandwidth allotment detail records associated with this virtual rack.
	DetailCount *uint `json:"detailCount,omitempty" xmlrpc:"detailCount,omitempty"`

	// The bandwidth allotment detail records associated with this virtual rack.
	Details []Network_Bandwidth_Version1_Allotment_Detail `json:"details,omitempty" xmlrpc:"details,omitempty"`

	// End date for an allotment.
	EndDate *Time `json:"endDate,omitempty" xmlrpc:"endDate,omitempty"`

	// The hardware contained within a virtual rack.
	Hardware []Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// A count of the hardware contained within a virtual rack.
	HardwareCount *uint `json:"hardwareCount,omitempty" xmlrpc:"hardwareCount,omitempty"`

	// A virtual rack's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The total public inbound bandwidth used in this virtual rack for an account's current billing cycle.
	InboundPublicBandwidthUsage *Float64 `json:"inboundPublicBandwidthUsage,omitempty" xmlrpc:"inboundPublicBandwidthUsage,omitempty"`

	// The location group associated with this virtual rack.
	LocationGroup *Location_Group `json:"locationGroup,omitempty" xmlrpc:"locationGroup,omitempty"`

	// Location Group Id for an allotment
	LocationGroupId *int `json:"locationGroupId,omitempty" xmlrpc:"locationGroupId,omitempty"`

	// A count of the managed bare metal server instances contained within a virtual rack.
	ManagedBareMetalInstanceCount *uint `json:"managedBareMetalInstanceCount,omitempty" xmlrpc:"managedBareMetalInstanceCount,omitempty"`

	// The managed bare metal server instances contained within a virtual rack.
	ManagedBareMetalInstances []Hardware `json:"managedBareMetalInstances,omitempty" xmlrpc:"managedBareMetalInstances,omitempty"`

	// The managed hardware contained within a virtual rack.
	ManagedHardware []Hardware `json:"managedHardware,omitempty" xmlrpc:"managedHardware,omitempty"`

	// A count of the managed hardware contained within a virtual rack.
	ManagedHardwareCount *uint `json:"managedHardwareCount,omitempty" xmlrpc:"managedHardwareCount,omitempty"`

	// A count of the managed Virtual Server contained within a virtual rack.
	ManagedVirtualGuestCount *uint `json:"managedVirtualGuestCount,omitempty" xmlrpc:"managedVirtualGuestCount,omitempty"`

	// The managed Virtual Server contained within a virtual rack.
	ManagedVirtualGuests []Virtual_Guest `json:"managedVirtualGuests,omitempty" xmlrpc:"managedVirtualGuests,omitempty"`

	// A virtual rack's metric tracking object. This object records all periodic polled data available to this rack.
	MetricTrackingObject *Metric_Tracking_Object `json:"metricTrackingObject,omitempty" xmlrpc:"metricTrackingObject,omitempty"`

	// The metric tracking object id for this allotment.
	MetricTrackingObjectId *int `json:"metricTrackingObjectId,omitempty" xmlrpc:"metricTrackingObjectId,omitempty"`

	// Text A virtual rack's name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The total public outbound bandwidth used in this virtual rack for an account's current billing cycle.
	OutboundPublicBandwidthUsage *Float64 `json:"outboundPublicBandwidthUsage,omitempty" xmlrpc:"outboundPublicBandwidthUsage,omitempty"`

	// Whether the bandwidth usage for this bandwidth pool for the current billing cycle exceeds the allocation.
	OverBandwidthAllocationFlag *int `json:"overBandwidthAllocationFlag,omitempty" xmlrpc:"overBandwidthAllocationFlag,omitempty"`

	// The private network only hardware contained within a virtual rack.
	PrivateNetworkOnlyHardware []Hardware `json:"privateNetworkOnlyHardware,omitempty" xmlrpc:"privateNetworkOnlyHardware,omitempty"`

	// A count of the private network only hardware contained within a virtual rack.
	PrivateNetworkOnlyHardwareCount *uint `json:"privateNetworkOnlyHardwareCount,omitempty" xmlrpc:"privateNetworkOnlyHardwareCount,omitempty"`

	// Whether the bandwidth usage for this bandwidth pool for the current billing cycle is projected to exceed the allocation.
	ProjectedOverBandwidthAllocationFlag *int `json:"projectedOverBandwidthAllocationFlag,omitempty" xmlrpc:"projectedOverBandwidthAllocationFlag,omitempty"`

	// The projected public outbound bandwidth for this virtual server for the current billing cycle.
	ProjectedPublicBandwidthUsage *Float64 `json:"projectedPublicBandwidthUsage,omitempty" xmlrpc:"projectedPublicBandwidthUsage,omitempty"`

	// no documentation yet
	ServiceProvider *Service_Provider `json:"serviceProvider,omitempty" xmlrpc:"serviceProvider,omitempty"`

	// Service Provider Id for an allotment
	ServiceProviderId *int `json:"serviceProviderId,omitempty" xmlrpc:"serviceProviderId,omitempty"`

	// The combined allocated bandwidth for all servers in a virtual rack.
	TotalBandwidthAllocated *uint `json:"totalBandwidthAllocated,omitempty" xmlrpc:"totalBandwidthAllocated,omitempty"`

	// A count of the Virtual Server contained within a virtual rack.
	VirtualGuestCount *uint `json:"virtualGuestCount,omitempty" xmlrpc:"virtualGuestCount,omitempty"`

	// The Virtual Server contained within a virtual rack.
	VirtualGuests []Virtual_Guest `json:"virtualGuests,omitempty" xmlrpc:"virtualGuests,omitempty"`
}

// The SoftLayer_Network_Bandwidth_Version1_Allotment_Detail data type contains specific information relating to a single bandwidth allotment record.
type Network_Bandwidth_Version1_Allotment_Detail struct {
	Entity

	// Allocated bandwidth.
	Allocation *Network_Bandwidth_Version1_Allocation `json:"allocation,omitempty" xmlrpc:"allocation,omitempty"`

	// Allocated bandwidth.
	AllocationId *int `json:"allocationId,omitempty" xmlrpc:"allocationId,omitempty"`

	// The parent Bandwidth Pool.
	BandwidthAllotment *Network_Bandwidth_Version1_Allotment `json:"bandwidthAllotment,omitempty" xmlrpc:"bandwidthAllotment,omitempty"`

	// Bandwidth Pool associated with this detail.
	BandwidthAllotmentId *int `json:"bandwidthAllotmentId,omitempty" xmlrpc:"bandwidthAllotmentId,omitempty"`

	// Beginning this date the bandwidth allotment is active.
	EffectiveDate *Time `json:"effectiveDate,omitempty" xmlrpc:"effectiveDate,omitempty"`

	// From this date the bandwidth allotment is no longer active.
	EndEffectiveDate *Time `json:"endEffectiveDate,omitempty" xmlrpc:"endEffectiveDate,omitempty"`

	// Internal ID associated with this allotment detail.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Service Provider Id for an allotment
	ServiceProviderId *int `json:"serviceProviderId,omitempty" xmlrpc:"serviceProviderId,omitempty"`
}

// The SoftLayer_Network_Bandwidth_Version1_Allotment_Type contains a description of the associated SoftLayer_Network_Bandwidth_Version1_Allotment object.
type Network_Bandwidth_Version1_Allotment_Type struct {
	Entity

	// no documentation yet
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	ShortDescription *string `json:"shortDescription,omitempty" xmlrpc:"shortDescription,omitempty"`
}

// The SoftLayer_Network_Bandwidth_Version1_Usage_Detail data type contains specific information relating to bandwidth utilization at a specific point in time on a given network interface.
type Network_Bandwidth_Version1_Usage_Detail struct {
	Entity

	// Incoming bandwidth utilization .
	AmountIn *Float64 `json:"amountIn,omitempty" xmlrpc:"amountIn,omitempty"`

	// Outgoing bandwidth utilization .
	AmountOut *Float64 `json:"amountOut,omitempty" xmlrpc:"amountOut,omitempty"`

	// Describes this bandwidth utilization record as on the public or private network interface.
	BandwidthUsageDetailType *Network_Bandwidth_Version1_Usage_Detail_Type `json:"bandwidthUsageDetailType,omitempty" xmlrpc:"bandwidthUsageDetailType,omitempty"`

	// Day and time this bandwidth utilization event was recorded.
	Day *Time `json:"day,omitempty" xmlrpc:"day,omitempty"`
}

// The SoftLayer_Network_Bandwidth_Version1_Usage_Detail_Type data type contains generic information relating to the types of bandwidth records available, currently just public and private.
type Network_Bandwidth_Version1_Usage_Detail_Type struct {
	Entity

	// Database key associated with this bandwidth detail type.
	Alias *string `json:"alias,omitempty" xmlrpc:"alias,omitempty"`
}

// The SoftLayer_Network_CdnMarketplace_Account data type models an individual CDN account. CDN accounts contain the SoftLayer account ID of the customer, the vendor ID the account belongs to, the customer ID provided by the vendor, and a CDN account's status.
type Network_CdnMarketplace_Account struct {
	Entity

	// SoftLayer account to which the CDN account belongs.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// An associated parent billing item which is active.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`
}

// no documentation yet
type Network_CdnMarketplace_Configuration_Behavior_Geoblocking struct {
	Entity

	// no documentation yet
	AccessType *string `json:"accessType,omitempty" xmlrpc:"accessType,omitempty"`

	// no documentation yet
	RegionType *string `json:"regionType,omitempty" xmlrpc:"regionType,omitempty"`

	// no documentation yet
	Regions []string `json:"regions,omitempty" xmlrpc:"regions,omitempty"`

	// no documentation yet
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`
}

// no documentation yet
type Network_CdnMarketplace_Configuration_Behavior_Geoblocking_Type struct {
	Entity

	// no documentation yet
	AccessType []string `json:"accessType,omitempty" xmlrpc:"accessType,omitempty"`

	// no documentation yet
	Continent []string `json:"continent,omitempty" xmlrpc:"continent,omitempty"`

	// no documentation yet
	CountryOrRegion []string `json:"countryOrRegion,omitempty" xmlrpc:"countryOrRegion,omitempty"`

	// no documentation yet
	RegionType []string `json:"regionType,omitempty" xmlrpc:"regionType,omitempty"`
}

// no documentation yet
type Network_CdnMarketplace_Configuration_Behavior_HotlinkProtection struct {
	Entity

	// no documentation yet
	ProtectionType *string `json:"protectionType,omitempty" xmlrpc:"protectionType,omitempty"`

	// no documentation yet
	RefererValues *string `json:"refererValues,omitempty" xmlrpc:"refererValues,omitempty"`
}

// no documentation yet
type Network_CdnMarketplace_Configuration_Behavior_ModifyResponseHeader struct {
	Entity
}

// no documentation yet
type Network_CdnMarketplace_Configuration_Behavior_TokenAuth struct {
	Entity
}

// This data type models a purge event that occurs in caching server. It contains a reference to a mapping configuration, the path to execute the purge on, the status of the purge, and flag that enables saving the purge information for future use.
type Network_CdnMarketplace_Configuration_Cache_Purge struct {
	Entity
}

// This data type models a purge group event that occurs in caching server. It contains a reference to a mapping configuration and the path to execute the purge on.
type Network_CdnMarketplace_Configuration_Cache_PurgeGroup struct {
	Entity
}

// This data type models a purge history event that occurs in caching server. The purge group history will be deleted after 15 days. The possible purge status of each history can be 'SUCCESS', "FAILED" or "IN_PROGRESS".
type Network_CdnMarketplace_Configuration_Cache_PurgeHistory struct {
	Entity
}

// This data type models a purge event that occurs repetitively and automatically in caching server after a set interval of time. A time to live instance contains a reference to a mapping configuration, the path to execute the purge on, the result of the purge, and the time interval after which the purge will be executed.
type Network_CdnMarketplace_Configuration_Cache_TimeToLive struct {
	Entity

	// date record is created
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// Path where purge will be executed after TTL
	Path *string `json:"path,omitempty" xmlrpc:"path,omitempty"`

	// Time interval after which purge will occur repeatedly
	TimeToLive *int `json:"timeToLive,omitempty" xmlrpc:"timeToLive,omitempty"`
}

// This data type represents the mapping Configuration settings for enabling CDN services. Each instance contains a reference to a CDN account, and CDN configuration properties such as a domain, an origin host and its port, a cname we generate, a cname the vendor generates, and a status. Other properties include the type of content to be cached (static or dynamic), the origin type (a host server or an object storage account), and the protocol to be used for caching.
type Network_CdnMarketplace_Configuration_Mapping struct {
	Entity
}

// no documentation yet
type Network_CdnMarketplace_Configuration_Mapping_Path struct {
	Entity
}

// This Metrics class provides methods to get CDN metrics based on account or mapping unique id.
type Network_CdnMarketplace_Metrics struct {
	Entity
}

// no documentation yet
type Network_CdnMarketplace_Utils_Response struct {
	Entity

	// no documentation yet
	Code *int `json:"code,omitempty" xmlrpc:"code,omitempty"`

	// no documentation yet
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`
}

// The SoftLayer_Network_CdnMarketplace_Vendor contains information regarding a CDN Vendor. This class is associated with SoftLayer_Network_CdnMarketplace_Vendor_Attribute class.
type Network_CdnMarketplace_Vendor struct {
	Entity
}

// Every piece of hardware running in SoftLayer's datacenters connected to the public, private, or management networks (where applicable) have a corresponding network component. These network components are modeled by the SoftLayer_Network_Component data type. These data types reflect the servers' local ethernet and remote management interfaces.
type Network_Component struct {
	Entity

	// Reboot/power (rebootDefault, rebootSoft, rebootHard, powerOn, powerOff and powerCycle) command currently executing by the server's remote management card.
	ActiveCommand *Hardware_Component_RemoteManagement_Command_Request `json:"activeCommand,omitempty" xmlrpc:"activeCommand,omitempty"`

	// The network component linking this object to a child device
	DownlinkComponent *Network_Component `json:"downlinkComponent,omitempty" xmlrpc:"downlinkComponent,omitempty"`

	// The duplex mode of a network component.
	DuplexMode *Network_Component_Duplex_Mode `json:"duplexMode,omitempty" xmlrpc:"duplexMode,omitempty"`

	// A network component's Duplex mode.
	DuplexModeId *string `json:"duplexModeId,omitempty" xmlrpc:"duplexModeId,omitempty"`

	// The hardware that a network component resides in.
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// The internal identifier of the hardware that a network component belongs to.
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// no documentation yet
	HighAvailabilityFirewallFlag *bool `json:"highAvailabilityFirewallFlag,omitempty" xmlrpc:"highAvailabilityFirewallFlag,omitempty"`

	// A network component's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of the records of all IP addresses bound to a network component.
	IpAddressBindingCount *uint `json:"ipAddressBindingCount,omitempty" xmlrpc:"ipAddressBindingCount,omitempty"`

	// The records of all IP addresses bound to a network component.
	IpAddressBindings []Network_Component_IpAddress `json:"ipAddressBindings,omitempty" xmlrpc:"ipAddressBindings,omitempty"`

	// A count of
	IpAddressCount *uint `json:"ipAddressCount,omitempty" xmlrpc:"ipAddressCount,omitempty"`

	// no documentation yet
	IpAddresses []Network_Subnet_IpAddress `json:"ipAddresses,omitempty" xmlrpc:"ipAddresses,omitempty"`

	// The IP address of an IPMI-based management network component.
	IpmiIpAddress *string `json:"ipmiIpAddress,omitempty" xmlrpc:"ipmiIpAddress,omitempty"`

	// The MAC address of an IPMI-based management network component.
	IpmiMacAddress *string `json:"ipmiMacAddress,omitempty" xmlrpc:"ipmiMacAddress,omitempty"`

	// Last reboot/power (rebootDefault, rebootSoft, rebootHard, powerOn, powerOff and powerCycle) command issued to the server's remote management card.
	LastCommand *Hardware_Component_RemoteManagement_Command_Request `json:"lastCommand,omitempty" xmlrpc:"lastCommand,omitempty"`

	// A network component's unique MAC address. IPMI-based management network interfaces may not have a MAC address.
	MacAddress *string `json:"macAddress,omitempty" xmlrpc:"macAddress,omitempty"`

	// A network component's maximum allowed speed, measured in Mbit per second. ''maxSpeed'' is determined by the capabilities of the network interface and the port speed purchased on your SoftLayer server.
	MaxSpeed *int `json:"maxSpeed,omitempty" xmlrpc:"maxSpeed,omitempty"`

	// The metric tracking object for this network component.
	MetricTrackingObject *Metric_Tracking_Object `json:"metricTrackingObject,omitempty" xmlrpc:"metricTrackingObject,omitempty"`

	// The date a network component was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// A network component's short name. For most servers this is the string "eth" for ethernet ports or "mgmt" for remote management ports. Use this in conjunction with the ''port'' property to identify a network component. For instance, the "eth0" interface on a server has the network component name "eth" and port 0.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The upstream network component firewall.
	NetworkComponentFirewall *Network_Component_Firewall `json:"networkComponentFirewall,omitempty" xmlrpc:"networkComponentFirewall,omitempty"`

	// A network component's associated group.
	NetworkComponentGroup *Network_Component_Group `json:"networkComponentGroup,omitempty" xmlrpc:"networkComponentGroup,omitempty"`

	// All network devices in SoftLayer's network hierarchy that this device is connected to.
	NetworkHardware []Hardware `json:"networkHardware,omitempty" xmlrpc:"networkHardware,omitempty"`

	// A count of all network devices in SoftLayer's network hierarchy that this device is connected to.
	NetworkHardwareCount *uint `json:"networkHardwareCount,omitempty" xmlrpc:"networkHardwareCount,omitempty"`

	// The VLAN that a network component's subnet is associated with.
	NetworkVlan *Network_Vlan `json:"networkVlan,omitempty" xmlrpc:"networkVlan,omitempty"`

	// The unique internal id of the network VLAN that the port belongs to.
	NetworkVlanId *int `json:"networkVlanId,omitempty" xmlrpc:"networkVlanId,omitempty"`

	// A count of the VLANs that are trunked to this network component.
	NetworkVlanTrunkCount *uint `json:"networkVlanTrunkCount,omitempty" xmlrpc:"networkVlanTrunkCount,omitempty"`

	// The VLANs that are trunked to this network component.
	NetworkVlanTrunks []Network_Component_Network_Vlan_Trunk `json:"networkVlanTrunks,omitempty" xmlrpc:"networkVlanTrunks,omitempty"`

	// The viable trunking targets of this component. Viable targets include accessible VLANs in the same pod and network as this component, which are not already natively attached nor trunked to this component.
	NetworkVlansTrunkable []Network_Vlan `json:"networkVlansTrunkable,omitempty" xmlrpc:"networkVlansTrunkable,omitempty"`

	// A network component's port number. Most hardware has more than one network interface. The port property separates these interfaces. Use this in conjunction with the ''name'' property to identify a network component. For instance, the "eth0" interface on a server has the network component name "eth" and port 0.
	Port *int `json:"port,omitempty" xmlrpc:"port,omitempty"`

	// A network component's primary IP address. IPMI-based management network interfaces may not have an IP address.
	PrimaryIpAddress *string `json:"primaryIpAddress,omitempty" xmlrpc:"primaryIpAddress,omitempty"`

	// The primary IPv4 Address record for a network component.
	PrimaryIpAddressRecord *Network_Subnet_IpAddress `json:"primaryIpAddressRecord,omitempty" xmlrpc:"primaryIpAddressRecord,omitempty"`

	// The subnet of the primary IP address assigned to this network component.
	PrimarySubnet *Network_Subnet `json:"primarySubnet,omitempty" xmlrpc:"primarySubnet,omitempty"`

	// The primary IPv6 Address record for a network component.
	PrimaryVersion6IpAddressRecord *Network_Subnet_IpAddress `json:"primaryVersion6IpAddressRecord,omitempty" xmlrpc:"primaryVersion6IpAddressRecord,omitempty"`

	// A count of the last five reboot/power (rebootDefault, rebootSoft, rebootHard, powerOn, powerOff and powerCycle) commands issued to the server's remote management card.
	RecentCommandCount *uint `json:"recentCommandCount,omitempty" xmlrpc:"recentCommandCount,omitempty"`

	// The last five reboot/power (rebootDefault, rebootSoft, rebootHard, powerOn, powerOff and powerCycle) commands issued to the server's remote management card.
	RecentCommands []Hardware_Component_RemoteManagement_Command_Request `json:"recentCommands,omitempty" xmlrpc:"recentCommands,omitempty"`

	// Indicates whether the network component is participating in a group of two or more components capable of being operationally redundant, if enabled.
	RedundancyCapableFlag *bool `json:"redundancyCapableFlag,omitempty" xmlrpc:"redundancyCapableFlag,omitempty"`

	// Indicates whether the network component is participating in a group of two or more components which is actively providing link redundancy.
	RedundancyEnabledFlag *bool `json:"redundancyEnabledFlag,omitempty" xmlrpc:"redundancyEnabledFlag,omitempty"`

	// A count of user(s) credentials to issue commands and/or interact with the server's remote management card.
	RemoteManagementUserCount *uint `json:"remoteManagementUserCount,omitempty" xmlrpc:"remoteManagementUserCount,omitempty"`

	// User(s) credentials to issue commands and/or interact with the server's remote management card.
	RemoteManagementUsers []Hardware_Component_RemoteManagement_User `json:"remoteManagementUsers,omitempty" xmlrpc:"remoteManagementUsers,omitempty"`

	// A network component's routers.
	Router *Hardware `json:"router,omitempty" xmlrpc:"router,omitempty"`

	// A network component's speed, measured in Mbit per second.
	Speed *int `json:"speed,omitempty" xmlrpc:"speed,omitempty"`

	// A network component's status. This can take one of four possible values: "ACTIVE", "DISABLE", "USER_OFF", or "MACWAIT". "ACTIVE" network components are enabled and in use on a servers. "DISABLE" status components have been administratively disabled by SoftLayer accounting or abuse. "USER_OFF" components have been administratively disabled by you, the user. "MACWAIT" components only exist on network components that have not been provisioned. You should never see a network interface in MACWAIT state. If you happen to see one please contact SoftLayer support.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// Whether a network component's primary ip address is from a storage network subnet or not. [Deprecated]
	// Deprecated: This function has been marked as deprecated.
	StorageNetworkFlag *bool `json:"storageNetworkFlag,omitempty" xmlrpc:"storageNetworkFlag,omitempty"`

	// A count of a network component's subnets. A subnet is a group of IP addresses
	SubnetCount *uint `json:"subnetCount,omitempty" xmlrpc:"subnetCount,omitempty"`

	// A network component's subnets. A subnet is a group of IP addresses
	Subnets []Network_Subnet `json:"subnets,omitempty" xmlrpc:"subnets,omitempty"`

	// The network component linking this object to parent
	UplinkComponent *Network_Component `json:"uplinkComponent,omitempty" xmlrpc:"uplinkComponent,omitempty"`

	// The duplex mode of the uplink network component linking to this object
	UplinkDuplexMode *Network_Component_Duplex_Mode `json:"uplinkDuplexMode,omitempty" xmlrpc:"uplinkDuplexMode,omitempty"`
}

// Duplex Mode allows finer grained control over networking options and settings.
type Network_Component_Duplex_Mode struct {
	Entity

	// no documentation yet
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The SoftLayer_Network_Component_Firewall data type contains general information relating to a single SoftLayer network component firewall. This is the object which ties the running rules to a specific downstream server. Use the [[SoftLayer Network Firewall Template]] service to pull SoftLayer recommended rule set templates. Use the [[SoftLayer Network Firewall Update Request]] service to submit a firewall update request.
type Network_Component_Firewall struct {
	Entity

	// A count of the additional subnets linked to this network component firewall, that inherit rules from the host that the context slot is attached to.
	ApplyServerRuleSubnetCount *uint `json:"applyServerRuleSubnetCount,omitempty" xmlrpc:"applyServerRuleSubnetCount,omitempty"`

	// The additional subnets linked to this network component firewall, that inherit rules from the host that the context slot is attached to.
	ApplyServerRuleSubnets []Network_Subnet `json:"applyServerRuleSubnets,omitempty" xmlrpc:"applyServerRuleSubnets,omitempty"`

	// The billing item for a Hardware Firewall (Dedicated).
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// The network component of the guest virtual server that this network component firewall belongs to.
	GuestNetworkComponent *Virtual_Guest_Network_Component `json:"guestNetworkComponent,omitempty" xmlrpc:"guestNetworkComponent,omitempty"`

	// Unique ID for the network component of the switch interface that this network component firewall is attached to.
	GuestNetworkComponentId *int `json:"guestNetworkComponentId,omitempty" xmlrpc:"guestNetworkComponentId,omitempty"`

	// Unique ID for the network component firewall.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The network component of the switch interface that this network component firewall belongs to.
	NetworkComponent *Network_Component `json:"networkComponent,omitempty" xmlrpc:"networkComponent,omitempty"`

	// Unique ID for the network component of the switch interface that this network component firewall is attached to.
	NetworkComponentId *int `json:"networkComponentId,omitempty" xmlrpc:"networkComponentId,omitempty"`

	// The update requests made for this firewall.
	NetworkFirewallUpdateRequest []Network_Firewall_Update_Request `json:"networkFirewallUpdateRequest,omitempty" xmlrpc:"networkFirewallUpdateRequest,omitempty"`

	// A count of the update requests made for this firewall.
	NetworkFirewallUpdateRequestCount *uint `json:"networkFirewallUpdateRequestCount,omitempty" xmlrpc:"networkFirewallUpdateRequestCount,omitempty"`

	// A count of the currently running rule set of this network component firewall.
	RuleCount *uint `json:"ruleCount,omitempty" xmlrpc:"ruleCount,omitempty"`

	// The currently running rule set of this network component firewall.
	Rules []Network_Component_Firewall_Rule `json:"rules,omitempty" xmlrpc:"rules,omitempty"`

	// Current status of the network component firewall. Status "no_edit" means this host is not protected by a hardware firewall. Status "allow_edit" means this host is protected by a hardware firewall and processing firewall rules. Status "bypass" means this host is provisioned behind a hardware firewall, but bypassing the firewall rules.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// A count of the additional subnets linked to this network component firewall.
	SubnetCount *uint `json:"subnetCount,omitempty" xmlrpc:"subnetCount,omitempty"`

	// The additional subnets linked to this network component firewall.
	Subnets []Network_Subnet `json:"subnets,omitempty" xmlrpc:"subnets,omitempty"`
}

// A SoftLayer_Network_Component_Firewall_Rule object type represents a currently running firewall rule and contains relative information. Use the [[SoftLayer Network Firewall Update Request]] service to submit a firewall update request. Use the [[SoftLayer Network Firewall Template]] service to pull SoftLayer recommended rule set templates.
type Network_Component_Firewall_Rule struct {
	Entity

	// The action that the rule is to take [permit or deny].
	Action *string `json:"action,omitempty" xmlrpc:"action,omitempty"`

	// The destination IP address considered for determining rule application.
	DestinationIpAddress *string `json:"destinationIpAddress,omitempty" xmlrpc:"destinationIpAddress,omitempty"`

	// The CIDR is used for determining rule application. This value will
	DestinationIpCidr *int `json:"destinationIpCidr,omitempty" xmlrpc:"destinationIpCidr,omitempty"`

	// The destination IP subnet mask considered for determining rule application.
	DestinationIpSubnetMask *string `json:"destinationIpSubnetMask,omitempty" xmlrpc:"destinationIpSubnetMask,omitempty"`

	// The ending (upper end of range) destination port considered for determining rule application.
	DestinationPortRangeEnd *int `json:"destinationPortRangeEnd,omitempty" xmlrpc:"destinationPortRangeEnd,omitempty"`

	// The starting (lower end of range) destination port considered for determining rule application.
	DestinationPortRangeStart *int `json:"destinationPortRangeStart,omitempty" xmlrpc:"destinationPortRangeStart,omitempty"`

	// The rule's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The network component firewall that this rule belongs to.
	NetworkComponentFirewall *Network_Component_Firewall `json:"networkComponentFirewall,omitempty" xmlrpc:"networkComponentFirewall,omitempty"`

	// The notes field for the rule.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// The numeric value describing the order in which the rule should be applied.
	OrderValue *int `json:"orderValue,omitempty" xmlrpc:"orderValue,omitempty"`

	// The protocol considered for determining rule application.
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// The source IP address considered for determining rule application.
	SourceIpAddress *string `json:"sourceIpAddress,omitempty" xmlrpc:"sourceIpAddress,omitempty"`

	// The CIDR is used for determining rule application. This value will
	SourceIpCidr *int `json:"sourceIpCidr,omitempty" xmlrpc:"sourceIpCidr,omitempty"`

	// The source IP subnet mask considered for determining rule application.
	SourceIpSubnetMask *string `json:"sourceIpSubnetMask,omitempty" xmlrpc:"sourceIpSubnetMask,omitempty"`

	// Current status of the network component firewall.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// Whether this rule is an IPv4 rule or an IPv6 rule. If
	Version *int `json:"version,omitempty" xmlrpc:"version,omitempty"`
}

// A SoftLayer_Network_Component_Firewall_Subnets object type represents the current linked subnets and contains relative information. Use the [[SoftLayer Network Firewall Update Request]] service to submit a firewall update request. Use the [[SoftLayer Network Firewall Template]] service to pull SoftLayer recommended rule set templates.
type Network_Component_Firewall_Subnets struct {
	Entity

	// A boolean flag that indicates whether the subnet should receive all the rules intended for the host on this context slot.
	ApplyServerRulesFlag *bool `json:"applyServerRulesFlag,omitempty" xmlrpc:"applyServerRulesFlag,omitempty"`

	// The network component firewall that write rules for this subnet.
	NetworkComponentFirewall *Network_Component_Firewall `json:"networkComponentFirewall,omitempty" xmlrpc:"networkComponentFirewall,omitempty"`

	// The subnet that this link binds to the network component firewall.
	Subnet *Network_Subnet `json:"subnet,omitempty" xmlrpc:"subnet,omitempty"`

	// The unique identifier of the subnet being linked to the network component firewall.
	SubnetId *int `json:"subnetId,omitempty" xmlrpc:"subnetId,omitempty"`
}

// no documentation yet
type Network_Component_Group struct {
	Entity

	// no documentation yet
	GroupTypeId *int `json:"groupTypeId,omitempty" xmlrpc:"groupTypeId,omitempty"`

	// A succinct label describing the members of this grouping.
	MembersDescription *string `json:"membersDescription,omitempty" xmlrpc:"membersDescription,omitempty"`

	// A count of a network component group's associated network components.
	NetworkComponentCount *uint `json:"networkComponentCount,omitempty" xmlrpc:"networkComponentCount,omitempty"`

	// A network component group's associated network components.
	NetworkComponents []Network_Component `json:"networkComponents,omitempty" xmlrpc:"networkComponents,omitempty"`
}

// The SoftLayer_Network_Component_IpAddress data type contains general information relating to the binding of a single network component to a single SoftLayer IP address.
type Network_Component_IpAddress struct {
	Entity

	// The IP address associated with this object's network component.
	IpAddress *Network_Subnet_IpAddress `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// The network component associated with this object's IP address.
	NetworkComponent *Network_Component `json:"networkComponent,omitempty" xmlrpc:"networkComponent,omitempty"`
}

// Represents the association between a Network_Component and Network_Vlan in the manner of a 'trunk'. Trunking a VLAN to a port allows that ports to receive and send packets tagged with the corresponding VLAN number.
type Network_Component_Network_Vlan_Trunk struct {
	Entity

	// A value of '1' indicates the existence of an ongoing request to modify this trunk record.
	IsUpdating *bool `json:"isUpdating,omitempty" xmlrpc:"isUpdating,omitempty"`

	// The network component that the VLAN is being trunked to.
	NetworkComponent *Network_Component `json:"networkComponent,omitempty" xmlrpc:"networkComponent,omitempty"`

	// The network component's identifier.
	NetworkComponentId *int `json:"networkComponentId,omitempty" xmlrpc:"networkComponentId,omitempty"`

	// The VLAN that is being trunked to the network component.
	NetworkVlan *Network_Vlan `json:"networkVlan,omitempty" xmlrpc:"networkVlan,omitempty"`

	// The identifier of the network VLAN that is a trunk on the network component.
	NetworkVlanId *int `json:"networkVlanId,omitempty" xmlrpc:"networkVlanId,omitempty"`
}

// The SoftLayer_Network_Component_RemoteManagement data type contains general information relating to a single SoftLayer remote management network component.
type Network_Component_RemoteManagement struct {
	Network_Component
}

// The SoftLayer_Network_Component_Uplink_Hardware data type abstracts information related to network connections between SoftLayer hardware and SoftLayer network components.
//
// It is populated via triggers on the network_connection table (SoftLayer_Network_Connection), so you shouldn't have to delete or insert records into this table, ever.
type Network_Component_Uplink_Hardware struct {
	Entity

	// A network component uplink's connected [[SoftLayer_Hardware|Hardware]].
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// The [[SoftLayer_Network_Component|Network Component]] that a uplink connection belongs to..
	NetworkComponent *Network_Component `json:"networkComponent,omitempty" xmlrpc:"networkComponent,omitempty"`
}

// The SoftLayer_Network_Customer_Subnet data type contains general information relating to a single customer subnet (remote).
type Network_Customer_Subnet struct {
	Entity

	// The account id a customer subnet belongs to.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// A subnet's Classless Inter-Domain Routing prefix. This is a number between 0 and 32 signifying the number of bits in a subnet's netmask. These bits separate a subnet's network address from it's host addresses. It performs the same function as the ''netmask'' property, but is represented as an integer.
	Cidr *int `json:"cidr,omitempty" xmlrpc:"cidr,omitempty"`

	// A customer subnet's unique identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of all ip addresses associated with a subnet.
	IpAddressCount *uint `json:"ipAddressCount,omitempty" xmlrpc:"ipAddressCount,omitempty"`

	// All ip addresses associated with a subnet.
	IpAddresses []Network_Customer_Subnet_IpAddress `json:"ipAddresses,omitempty" xmlrpc:"ipAddresses,omitempty"`

	// A bitmask in dotted-quad format that is used to separate a subnet's network address from it's host addresses. This performs the same function as the ''cidr'' property, but is expressed in a string format.
	Netmask *string `json:"netmask,omitempty" xmlrpc:"netmask,omitempty"`

	// A subnet's network identifier. This is the first IP address of a subnet.
	NetworkIdentifier *string `json:"networkIdentifier,omitempty" xmlrpc:"networkIdentifier,omitempty"`

	// The total number of ip addresses in a subnet.
	TotalIpAddresses *int `json:"totalIpAddresses,omitempty" xmlrpc:"totalIpAddresses,omitempty"`
}

// The SoftLayer_Network_Customer_Subnet_IpAddress data type contains general information relating to a single Customer Subnet (Remote) IPv4 address.
type Network_Customer_Subnet_IpAddress struct {
	Entity

	// Unique identifier for an ip address.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// An IP address expressed in dotted quad format.
	IpAddress *string `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// An IP address' user defined note.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// The customer subnet (remote) that the ip address belongs to.
	Subnet *Network_Customer_Subnet `json:"subnet,omitempty" xmlrpc:"subnet,omitempty"`

	// The unique identifier for the customer subnet (remote) the ip address belongs to.
	SubnetId *int `json:"subnetId,omitempty" xmlrpc:"subnetId,omitempty"`

	// A count of all the address translations that are tied to an IP address.
	TranslationCount *uint `json:"translationCount,omitempty" xmlrpc:"translationCount,omitempty"`

	// All the address translations that are tied to an IP address.
	Translations []Network_Tunnel_Module_Context_Address_Translation `json:"translations,omitempty" xmlrpc:"translations,omitempty"`
}

// The SoftLayer_Network_DirectLink_Location presents a structure containing attributes of a Direct Link location, and its related object SoftLayer location.
type Network_DirectLink_Location struct {
	Entity

	// The Direct Link specific location owner for POP/DC facilities. Like Equinix, Pacnet, Verizon etc.
	BuildingColocationOwner *string `json:"buildingColocationOwner,omitempty" xmlrpc:"buildingColocationOwner,omitempty"`

	// The unique identifier of a Direct Link location.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Specifies if The Direct Link specific location has Redundancy:secondary XCR availability.
	IsRedundantXcr *bool `json:"isRedundantXcr,omitempty" xmlrpc:"isRedundantXcr,omitempty"`

	// The location of Direct Link facility.
	Location *Location `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// The Direct Link specific location ie. Data Center & Network POP facility. Refer to location object Like Dallas in US, London in England etc.
	LocationId *int `json:"locationId,omitempty" xmlrpc:"locationId,omitempty"`

	// The Direct Link Market location used in Direct Link Order. Like Europe, North America, Asia pacific etc.
	MarketGeography *string `json:"marketGeography,omitempty" xmlrpc:"marketGeography,omitempty"`

	// The Id of Direct Link provider.
	Provider *Network_DirectLink_Provider `json:"provider,omitempty" xmlrpc:"provider,omitempty"`

	// The Id of Direct Link service type.
	ServiceType *Network_DirectLink_ServiceType `json:"serviceType,omitempty" xmlrpc:"serviceType,omitempty"`
}

// The SoftLayer_Network_DirectLink_Provider presents a structure containing attributes of a Direct Link provider.
type Network_DirectLink_Provider struct {
	Entity

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The SoftLayer_Network_DirectLink_ServiceType presents a structure containing attributes of a Direct Link Service Type.
type Network_DirectLink_ServiceType struct {
	Entity

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// The SoftLayer_Network_Firewall_AccessControlList data type contains general information relating to a single SoftLayer firewall access to controll list. This is the object which ties the running rules to a specific context. Use the [[SoftLayer Network Firewall Template]] service to pull SoftLayer recommended rule set templates. Use the [[SoftLayer Network Firewall Update Request]] service to submit a firewall update request.
type Network_Firewall_AccessControlList struct {
	Entity

	// no documentation yet
	Direction *string `json:"direction,omitempty" xmlrpc:"direction,omitempty"`

	// no documentation yet
	FirewallContextInterfaceId *int `json:"firewallContextInterfaceId,omitempty" xmlrpc:"firewallContextInterfaceId,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of the update requests made for this firewall.
	NetworkFirewallUpdateRequestCount *uint `json:"networkFirewallUpdateRequestCount,omitempty" xmlrpc:"networkFirewallUpdateRequestCount,omitempty"`

	// The update requests made for this firewall.
	NetworkFirewallUpdateRequests []Network_Firewall_Update_Request `json:"networkFirewallUpdateRequests,omitempty" xmlrpc:"networkFirewallUpdateRequests,omitempty"`

	// no documentation yet
	NetworkVlan *Network_Vlan `json:"networkVlan,omitempty" xmlrpc:"networkVlan,omitempty"`

	// A count of the currently running rule set of this context access control list firewall.
	RuleCount *uint `json:"ruleCount,omitempty" xmlrpc:"ruleCount,omitempty"`

	// The currently running rule set of this context access control list firewall.
	Rules []Network_Vlan_Firewall_Rule `json:"rules,omitempty" xmlrpc:"rules,omitempty"`
}

// The SoftLayer_Network_Firewall_Interface data type contains general information relating to a single SoftLayer firewall interface. This is the object which ties the firewall context access control list to a firewall. Use the [[SoftLayer Network Firewall Template]] service to pull SoftLayer recommended rule set templates. Use the [[SoftLayer Network Firewall Update Request]] service to submit a firewall update request.
type Network_Firewall_Interface struct {
	Network_Firewall_Module_Context_Interface
}

// no documentation yet
type Network_Firewall_Module_Context_Interface struct {
	Entity

	// A count of
	FirewallContextAccessControlListCount *uint `json:"firewallContextAccessControlListCount,omitempty" xmlrpc:"firewallContextAccessControlListCount,omitempty"`

	// no documentation yet
	FirewallContextAccessControlLists []Network_Firewall_AccessControlList `json:"firewallContextAccessControlLists,omitempty" xmlrpc:"firewallContextAccessControlLists,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	NetworkVlan *Network_Vlan `json:"networkVlan,omitempty" xmlrpc:"networkVlan,omitempty"`
}

// The SoftLayer_Network_Firewall_Template type contains general information for a SoftLayer network firewall template.
//
// Firewall templates are recommend rule sets for use with SoftLayer Hardware Firewall (Dedicated).  These optimized templates are designed to balance security restriction with application availability.  The templates given may be altered to provide custom network security, or may be used as-is for basic security. At least one rule set MUST be applied for the firewall to block traffic. Use the [[SoftLayer Network Component Firewall]] service to view current rules. Use the [[SoftLayer Network Firewall Update Request]] service to submit a firewall update request.
type Network_Firewall_Template struct {
	Entity

	// A Firewall template's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The name of the firewall rules template.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A count of the rule set that belongs to this firewall rules template.
	RuleCount *uint `json:"ruleCount,omitempty" xmlrpc:"ruleCount,omitempty"`

	// The rule set that belongs to this firewall rules template.
	Rules []Network_Firewall_Template_Rule `json:"rules,omitempty" xmlrpc:"rules,omitempty"`
}

// The SoftLayer_Network_Component_Firewall_Rule type contains general information relating to a single SoftLayer firewall template rule. Use the [[SoftLayer Network Component Firewall]] service to view current rules. Use the [[SoftLayer Network Firewall Update Request]] service to submit a firewall update request.
type Network_Firewall_Template_Rule struct {
	Entity

	// The action that this template rule is to take [permit or deny].
	Action *string `json:"action,omitempty" xmlrpc:"action,omitempty"`

	// The destination IP address considered for determining rule application.
	DestinationIpAddress *string `json:"destinationIpAddress,omitempty" xmlrpc:"destinationIpAddress,omitempty"`

	// The destination IP subnet mask considered for determining rule application.
	DestinationIpSubnetMask *string `json:"destinationIpSubnetMask,omitempty" xmlrpc:"destinationIpSubnetMask,omitempty"`

	// The ending (upper end of range) destination port considered for determining rule application.
	DestinationPortRangeEnd *int `json:"destinationPortRangeEnd,omitempty" xmlrpc:"destinationPortRangeEnd,omitempty"`

	// The starting (lower end of range) destination port considered for determining rule application.
	DestinationPortRangeStart *int `json:"destinationPortRangeStart,omitempty" xmlrpc:"destinationPortRangeStart,omitempty"`

	// The firewall template that this rule is attached to.
	FirewallTemplate *Network_Firewall_Template `json:"firewallTemplate,omitempty" xmlrpc:"firewallTemplate,omitempty"`

	// The unique identifier of the firewall template that a firewall template rule is associated with.
	FirewallTemplateId *int `json:"firewallTemplateId,omitempty" xmlrpc:"firewallTemplateId,omitempty"`

	// A Firewall template rule's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The notes field for the firewall template rule.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// The numeric value describing the order in which the rule set should be applied.
	OrderValue *int `json:"orderValue,omitempty" xmlrpc:"orderValue,omitempty"`

	// The protocol considered for determining rule application.
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// The source IP address considered for determining rule application.
	SourceIpAddress *string `json:"sourceIpAddress,omitempty" xmlrpc:"sourceIpAddress,omitempty"`

	// The source IP subnet mask considered for determining rule application.
	SourceIpSubnetMask *string `json:"sourceIpSubnetMask,omitempty" xmlrpc:"sourceIpSubnetMask,omitempty"`
}

// The SoftLayer_Network_Firewall_Update_Request data type contains information relating to a SoftLayer network firewall update request. Use the [[SoftLayer Network Component Firewall]] service to view current rules. Use the [[SoftLayer Network Firewall Template]] service to pull SoftLayer recommended rule set templates.
type Network_Firewall_Update_Request struct {
	Entity

	// Timestamp of when the rules from the update request were applied to the firewall.
	ApplyDate *Time `json:"applyDate,omitempty" xmlrpc:"applyDate,omitempty"`

	// The user that authorized this firewall update request.
	AuthorizingUser *User_Interface `json:"authorizingUser,omitempty" xmlrpc:"authorizingUser,omitempty"`

	// The unique identifier of the user that authorized the update request.
	AuthorizingUserId *int `json:"authorizingUserId,omitempty" xmlrpc:"authorizingUserId,omitempty"`

	// The type of user that authorized the update request [EMP or USR].
	AuthorizingUserType *string `json:"authorizingUserType,omitempty" xmlrpc:"authorizingUserType,omitempty"`

	// Flag indicating whether the request is for a rule bypass configuration [0 or 1].
	BypassFlag *bool `json:"bypassFlag,omitempty" xmlrpc:"bypassFlag,omitempty"`

	// Timestamp of the creation of the record.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The unique identifier of the firewall access control list that the rule set is destined for.
	FirewallContextAccessControlListId *int `json:"firewallContextAccessControlListId,omitempty" xmlrpc:"firewallContextAccessControlListId,omitempty"`

	// The downstream virtual server that the rule set will be applied to.
	Guest *Virtual_Guest `json:"guest,omitempty" xmlrpc:"guest,omitempty"`

	// The downstream server that the rule set will be applied to.
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// The unique identifier of the server that the rule set is destined to protect.
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// The unique identifier of the firewall update request.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The network component firewall that the rule set will be applied to.
	NetworkComponentFirewall *Network_Component_Firewall `json:"networkComponentFirewall,omitempty" xmlrpc:"networkComponentFirewall,omitempty"`

	// The unique identifier of the network component firewall that the rule set is destined for.
	NetworkComponentFirewallId *int `json:"networkComponentFirewallId,omitempty" xmlrpc:"networkComponentFirewallId,omitempty"`

	// A count of the group of rules contained within the update request.
	RuleCount *uint `json:"ruleCount,omitempty" xmlrpc:"ruleCount,omitempty"`

	// The group of rules contained within the update request.
	Rules []Network_Firewall_Update_Request_Rule `json:"rules,omitempty" xmlrpc:"rules,omitempty"`
}

// A SoftLayer_Ticket_Update_Customer is a single update made by a customer to a ticket.
type Network_Firewall_Update_Request_Customer struct {
	Network_Firewall_Update_Request
}

// The SoftLayer_Network_Firewall_Update_Request_Employee data type returns a user object for the SoftLayer employee that created the request.
type Network_Firewall_Update_Request_Employee struct {
	Network_Firewall_Update_Request
}

// The SoftLayer_Network_Firewall_Update_Request_Rule type contains information relating to a SoftLayer network firewall update request rule. This rule is a member of a [[SoftLayer Network Firewall Update Request]]. Use the [[SoftLayer Network Component Firewall]] service to view current rules. Use the [[SoftLayer Network Firewall Template]] service to pull SoftLayer recommended rule set templates.
type Network_Firewall_Update_Request_Rule struct {
	Entity

	// The action that this update request rule is to take [permit or deny].
	Action *string `json:"action,omitempty" xmlrpc:"action,omitempty"`

	// The bypassRuleValidation is used for bypassing the rule validation
	BypassRuleValidation *bool `json:"bypassRuleValidation,omitempty" xmlrpc:"bypassRuleValidation,omitempty"`

	// The destination IP address considered for determining rule application.
	DestinationIpAddress *string `json:"destinationIpAddress,omitempty" xmlrpc:"destinationIpAddress,omitempty"`

	// The CIDR is used for determining rule application. This value will
	DestinationIpCidr *int `json:"destinationIpCidr,omitempty" xmlrpc:"destinationIpCidr,omitempty"`

	// The destination IP subnet mask considered for determining rule application.
	DestinationIpSubnetMask *string `json:"destinationIpSubnetMask,omitempty" xmlrpc:"destinationIpSubnetMask,omitempty"`

	// The ending (upper end of range) destination port considered for determining rule application.
	DestinationPortRangeEnd *int `json:"destinationPortRangeEnd,omitempty" xmlrpc:"destinationPortRangeEnd,omitempty"`

	// The starting (lower end of range) destination port considered for determining rule application.
	DestinationPortRangeStart *int `json:"destinationPortRangeStart,omitempty" xmlrpc:"destinationPortRangeStart,omitempty"`

	// The update request that this rule belongs to.
	FirewallUpdateRequest *Network_Firewall_Update_Request `json:"firewallUpdateRequest,omitempty" xmlrpc:"firewallUpdateRequest,omitempty"`

	// The unique identifier of the firewall update request that a firewall update request rule is associated with.
	FirewallUpdateRequestId *int `json:"firewallUpdateRequestId,omitempty" xmlrpc:"firewallUpdateRequestId,omitempty"`

	// A Firewall update request rule's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The notes field for the firewall update request rule.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// The numeric value describing the order in which the rule should be applied.
	OrderValue *int `json:"orderValue,omitempty" xmlrpc:"orderValue,omitempty"`

	// The protocol considered for determining rule application.
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// The source IP address considered for determining rule application.
	SourceIpAddress *string `json:"sourceIpAddress,omitempty" xmlrpc:"sourceIpAddress,omitempty"`

	// The CIDR is used for determining rule application. This value will
	SourceIpCidr *int `json:"sourceIpCidr,omitempty" xmlrpc:"sourceIpCidr,omitempty"`

	// The source IP subnet mask considered for determining rule application.
	SourceIpSubnetMask *string `json:"sourceIpSubnetMask,omitempty" xmlrpc:"sourceIpSubnetMask,omitempty"`

	// Whether this rule is an IPv4 rule or an IPv6 rule. If
	Version *int `json:"version,omitempty" xmlrpc:"version,omitempty"`
}

// The SoftLayer_Network_Firewall_Update_Request_Rule_Version6 type contains information relating to a SoftLayer network firewall update request rule for IPv6. This rule is a member of a [[SoftLayer Network Firewall Update Request]]. Use the [[SoftLayer Network Component Firewall]] service to view current rules. Use the [[SoftLayer Network Firewall Template]] service to pull SoftLayer recommended rule set templates.
type Network_Firewall_Update_Request_Rule_Version6 struct {
	Network_Firewall_Update_Request_Rule
}

// no documentation yet
type Network_Gateway struct {
	Entity

	// The account for this gateway.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The internal identifier of the account assigned to this gateway.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The VRRP group number for this gateway. This is set internally and cannot be provided on create.
	GroupNumber *int `json:"groupNumber,omitempty" xmlrpc:"groupNumber,omitempty"`

	// A gateway's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of all VLANs trunked to this gateway.
	InsideVlanCount *uint `json:"insideVlanCount,omitempty" xmlrpc:"insideVlanCount,omitempty"`

	// All VLANs trunked to this gateway.
	InsideVlans []Network_Gateway_Vlan `json:"insideVlans,omitempty" xmlrpc:"insideVlans,omitempty"`

	// A count of the members for this gateway.
	MemberCount *uint `json:"memberCount,omitempty" xmlrpc:"memberCount,omitempty"`

	// The members for this gateway.
	Members []Network_Gateway_Member `json:"members,omitempty" xmlrpc:"members,omitempty"`

	// A gateway's name. This is required on create and can be no more than 255 characters.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The firewall associated with this gateway, if any.
	NetworkFirewall *Network_Vlan_Firewall `json:"networkFirewall,omitempty" xmlrpc:"networkFirewall,omitempty"`

	// Whether or not there is a firewall associated with this gateway.
	NetworkFirewallFlag *bool `json:"networkFirewallFlag,omitempty" xmlrpc:"networkFirewallFlag,omitempty"`

	// A gateway's network space. Currently, only 'private'  or 'both' is allowed. When this value is 'private', it is a backend gateway only. Otherwise, it is a gateway for both frontend and backend traffic.
	NetworkSpace *string `json:"networkSpace,omitempty" xmlrpc:"networkSpace,omitempty"`

	// A manufacturer of the gateway os.  This could be different from the manufacturer of the bare metal server os if the gateway is a VM.
	OsManufacturer *string `json:"osManufacturer,omitempty" xmlrpc:"osManufacturer,omitempty"`

	// The private gateway IP address.
	PrivateIpAddress *Network_Subnet_IpAddress `json:"privateIpAddress,omitempty" xmlrpc:"privateIpAddress,omitempty"`

	// The internal identifier of the private IP address for this gateway.
	PrivateIpAddressId *int `json:"privateIpAddressId,omitempty" xmlrpc:"privateIpAddressId,omitempty"`

	// The private VLAN for accessing this gateway.
	PrivateVlan *Network_Vlan `json:"privateVlan,omitempty" xmlrpc:"privateVlan,omitempty"`

	// The internal identifier of the private VLAN for this gateway.
	PrivateVlanId *int `json:"privateVlanId,omitempty" xmlrpc:"privateVlanId,omitempty"`

	// The public gateway IP address.
	PublicIpAddress *Network_Subnet_IpAddress `json:"publicIpAddress,omitempty" xmlrpc:"publicIpAddress,omitempty"`

	// The internal identifier of the public IP address for this gateway.
	PublicIpAddressId *int `json:"publicIpAddressId,omitempty" xmlrpc:"publicIpAddressId,omitempty"`

	// The public gateway IPv6 address.
	PublicIpv6Address *Network_Subnet_IpAddress `json:"publicIpv6Address,omitempty" xmlrpc:"publicIpv6Address,omitempty"`

	// The internal identifier of the public IPv6 address for this gateway.
	PublicIpv6AddressId *int `json:"publicIpv6AddressId,omitempty" xmlrpc:"publicIpv6AddressId,omitempty"`

	// The public VLAN for accessing this gateway.
	PublicVlan *Network_Vlan `json:"publicVlan,omitempty" xmlrpc:"publicVlan,omitempty"`

	// The internal identifier of the public VLAN for this gateway. This is set internally and cannot be provided on create.
	PublicVlanId *int `json:"publicVlanId,omitempty" xmlrpc:"publicVlanId,omitempty"`

	// The current status of the gateway.
	Status *Network_Gateway_Status `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// The current status of this gateway. This is always active unless there is a process running to change the gateway. This can not be set on creation.
	StatusId *int `json:"statusId,omitempty" xmlrpc:"statusId,omitempty"`
}

// no documentation yet
type Network_Gateway_Licenses struct {
	Entity

	// no documentation yet
	Employee *User_Employee `json:"employee,omitempty" xmlrpc:"employee,omitempty"`

	// no documentation yet
	ItemKeyName *string `json:"itemKeyName,omitempty" xmlrpc:"itemKeyName,omitempty"`

	// no documentation yet
	LicenseCategory *string `json:"licenseCategory,omitempty" xmlrpc:"licenseCategory,omitempty"`
}

// no documentation yet
type Network_Gateway_Member struct {
	Entity

	// The attributes for this member.
	Attributes *Network_Gateway_Member_Attribute `json:"attributes,omitempty" xmlrpc:"attributes,omitempty"`

	// The gateway software description for the member.
	GatewaySoftwareDescription *Software_Description `json:"gatewaySoftwareDescription,omitempty" xmlrpc:"gatewaySoftwareDescription,omitempty"`

	// no documentation yet
	GatewaySoftwareId *int `json:"gatewaySoftwareId,omitempty" xmlrpc:"gatewaySoftwareId,omitempty"`

	// The device for this member.
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// The internal identifier of the hardware for this member.
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// A gateway member's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of the gateway licenses for this member.
	LicenseCount *uint `json:"licenseCount,omitempty" xmlrpc:"licenseCount,omitempty"`

	// The gateway licenses for this member.
	Licenses []Network_Gateway_Member_Licenses `json:"licenses,omitempty" xmlrpc:"licenses,omitempty"`

	// The gateway this member belongs to.
	NetworkGateway *Network_Gateway `json:"networkGateway,omitempty" xmlrpc:"networkGateway,omitempty"`

	// The internal identifier of the gateway this member belongs to.
	NetworkGatewayId *int `json:"networkGatewayId,omitempty" xmlrpc:"networkGatewayId,omitempty"`

	// A count of the gateway passwords for this member.
	PasswordCount *uint `json:"passwordCount,omitempty" xmlrpc:"passwordCount,omitempty"`

	// The gateway passwords for this member.
	Passwords []Network_Gateway_Member_Passwords `json:"passwords,omitempty" xmlrpc:"passwords,omitempty"`

	// The priority for this gateway member. This is set internally and cannot be provided on create.
	Priority *int `json:"priority,omitempty" xmlrpc:"priority,omitempty"`

	// The public gateway IP address.
	PublicIpAddress *Network_Subnet_IpAddress `json:"publicIpAddress,omitempty" xmlrpc:"publicIpAddress,omitempty"`
}

// no documentation yet
type Network_Gateway_Member_Attribute struct {
	Entity

	// The gateway member has these attributes.
	GatewayMember *Network_Gateway_Member `json:"gatewayMember,omitempty" xmlrpc:"gatewayMember,omitempty"`

	// A gateway member's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Indicates if the member has been upgraded.
	IsUpgraded *int `json:"isUpgraded,omitempty" xmlrpc:"isUpgraded,omitempty"`

	// The previous version of the gateway software
	LastVersion *string `json:"lastVersion,omitempty" xmlrpc:"lastVersion,omitempty"`

	// Timestamp for the expiration date of the license key
	LicenseExpirationDate *Time `json:"licenseExpirationDate,omitempty" xmlrpc:"licenseExpirationDate,omitempty"`

	// no documentation yet
	LicenseKey *string `json:"licenseKey,omitempty" xmlrpc:"licenseKey,omitempty"`

	// The gateway member for this attribute.
	MemberId *int `json:"memberId,omitempty" xmlrpc:"memberId,omitempty"`

	// Network model of the gateway.
	NetworkModel *string `json:"networkModel,omitempty" xmlrpc:"networkModel,omitempty"`

	// Password of the user name.
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// Timestamp when this gateway member was last upgraded
	UpgradedDate *Time `json:"upgradedDate,omitempty" xmlrpc:"upgradedDate,omitempty"`

	// Username associated with the gateway.
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`

	// The version of the gateway software
	Version *string `json:"version,omitempty" xmlrpc:"version,omitempty"`

	// Precheck Warning code for Version / License Unsupported for member.
	WarningCode *int `json:"warningCode,omitempty" xmlrpc:"warningCode,omitempty"`
}

// no documentation yet
type Network_Gateway_Member_Licenses struct {
	Entity

	// no documentation yet
	ExpirationDate *Time `json:"expirationDate,omitempty" xmlrpc:"expirationDate,omitempty"`

	// The gateway license record.
	GatewayLicense *Network_Gateway_Licenses `json:"gatewayLicense,omitempty" xmlrpc:"gatewayLicense,omitempty"`

	// The gateway member has these licenses.
	GatewayMember *Network_Gateway_Member `json:"gatewayMember,omitempty" xmlrpc:"gatewayMember,omitempty"`

	// no documentation yet
	LicenseKey *string `json:"licenseKey,omitempty" xmlrpc:"licenseKey,omitempty"`
}

// no documentation yet
type Network_Gateway_Member_Passwords struct {
	Entity

	// The gateway member has these password.
	GatewayMember *Network_Gateway_Member `json:"gatewayMember,omitempty" xmlrpc:"gatewayMember,omitempty"`

	// A gateway member passlw internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The gateway member if for this record.
	MemberId *int `json:"memberId,omitempty" xmlrpc:"memberId,omitempty"`

	// Password of the user name.
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// Username associated with the gateway.
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`
}

// no documentation yet
type Network_Gateway_Precheck struct {
	Entity

	// Category name
	Category *string `json:"category,omitempty" xmlrpc:"category,omitempty"`

	// Gateway precheck status
	GatewayReadinessValue *string `json:"gatewayReadinessValue,omitempty" xmlrpc:"gatewayReadinessValue,omitempty"`

	// The gateway member for this precheck.
	MemberId *int `json:"memberId,omitempty" xmlrpc:"memberId,omitempty"`

	// Gateway precheck status
	MemberReadinessValue *string `json:"memberReadinessValue,omitempty" xmlrpc:"memberReadinessValue,omitempty"`

	// The precheck error status of the member
	ReturnCode *int `json:"returnCode,omitempty" xmlrpc:"returnCode,omitempty"`
}

// no documentation yet
type Network_Gateway_Status struct {
	Entity

	// A gateway status's description.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// A gateway status's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A gateway status's programmatic name.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// A gateway status's human-friendly name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Gateway_VersionUpgrade struct {
	Entity

	// Gateway version being upgraded from.
	FromVersion *string `json:"fromVersion,omitempty" xmlrpc:"fromVersion,omitempty"`

	// A gateway status's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Is OS Reload required after version upgrade?.
	OsReloadRequired *int `json:"osReloadRequired,omitempty" xmlrpc:"osReloadRequired,omitempty"`

	// Gateway version available for upgrade.
	ToVersion *string `json:"toVersion,omitempty" xmlrpc:"toVersion,omitempty"`
}

// no documentation yet
type Network_Gateway_Vlan struct {
	Entity

	// If true, this VLAN is bypassed. If false, it is routed through the gateway.
	BypassFlag *bool `json:"bypassFlag,omitempty" xmlrpc:"bypassFlag,omitempty"`

	// A gateway VLAN's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The gateway this VLAN is attached to.
	NetworkGateway *Network_Gateway `json:"networkGateway,omitempty" xmlrpc:"networkGateway,omitempty"`

	// The internal identifier of the gateway this VLAN is attached to.
	NetworkGatewayId *int `json:"networkGatewayId,omitempty" xmlrpc:"networkGatewayId,omitempty"`

	// The network VLAN record.
	NetworkVlan *Network_Vlan `json:"networkVlan,omitempty" xmlrpc:"networkVlan,omitempty"`

	// The internal identifier of the network VLAN.
	NetworkVlanId *int `json:"networkVlanId,omitempty" xmlrpc:"networkVlanId,omitempty"`
}

// no documentation yet
type Network_Interconnect_Tenant struct {
	Entity

	// no documentation yet
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// Specifies ASN used for BGP.
	BgpAsn *int `json:"bgpAsn,omitempty" xmlrpc:"bgpAsn,omitempty"`

	// The active billing item for a network interconnect.
	BillingItem *Billing_Item_Network_Interconnect `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	DatacenterName *string `json:"datacenterName,omitempty" xmlrpc:"datacenterName,omitempty"`

	// no documentation yet
	ErrorMessage *string `json:"errorMessage,omitempty" xmlrpc:"errorMessage,omitempty"`

	// The Direct Link connectivity to all SoftLayer data centers if globalRoutingFlag = 1 and local connectivity if globalRoutingFlag = 0.
	GlobalRoutingFlag *bool `json:"globalRoutingFlag,omitempty" xmlrpc:"globalRoutingFlag,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	InterconnectType *string `json:"interconnectType,omitempty" xmlrpc:"interconnectType,omitempty"`

	// Link speed of a Direct Link connection.
	LinkSpeed *int `json:"linkSpeed,omitempty" xmlrpc:"linkSpeed,omitempty"`

	// IP address (v4 or v6) of "near" router serial interface. No check/update of IP Address table.
	LocalIpAddress *string `json:"localIpAddress,omitempty" xmlrpc:"localIpAddress,omitempty"`

	// no documentation yet
	Location *string `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Specifies the Interconnect connection name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// Direct Link provider can request change to existing routing, Customer can approve the change. newGlobalRoutingFlag = 1 gives connectivity to all IBM data centers, and if newGlobalRoutingFlag = 0, it gives local connectivity.
	NewGlobalRoutingFlag *bool `json:"newGlobalRoutingFlag,omitempty" xmlrpc:"newGlobalRoutingFlag,omitempty"`

	// Updated Link speed of a Direct Link connection.
	NewLinkSpeed *int `json:"newLinkSpeed,omitempty" xmlrpc:"newLinkSpeed,omitempty"`

	// This field will have the ticket id if the tenant workflow fails
	Note *string `json:"note,omitempty" xmlrpc:"note,omitempty"`

	// Link speed of a Direct Link connection on Equinix Side.
	PeerLinkSpeed *int `json:"peerLinkSpeed,omitempty" xmlrpc:"peerLinkSpeed,omitempty"`

	// no documentation yet
	Port *string `json:"port,omitempty" xmlrpc:"port,omitempty"`

	// no documentation yet
	PortLabel *string `json:"portLabel,omitempty" xmlrpc:"portLabel,omitempty"`

	// no documentation yet
	Provider *string `json:"provider,omitempty" xmlrpc:"provider,omitempty"`

	// no documentation yet
	ProviderAccountId *int `json:"providerAccountId,omitempty" xmlrpc:"providerAccountId,omitempty"`

	// Specifies redundant connection is available if 1.
	RedundancyFlag *bool `json:"redundancyFlag,omitempty" xmlrpc:"redundancyFlag,omitempty"`

	// no documentation yet
	RemoteIpAddress *string `json:"remoteIpAddress,omitempty" xmlrpc:"remoteIpAddress,omitempty"`

	// Service key for Interconnect connection.
	ServiceKey *string `json:"serviceKey,omitempty" xmlrpc:"serviceKey,omitempty"`

	// no documentation yet
	ServiceType *Network_DirectLink_ServiceType `json:"serviceType,omitempty" xmlrpc:"serviceType,omitempty"`

	// no documentation yet
	ServiceTypeId *int `json:"serviceTypeId,omitempty" xmlrpc:"serviceTypeId,omitempty"`

	// The direct link connection status. IN_PROGRESS, PROVISIONING, CONNECTION_UP, CONNECTION_DOWN
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// no documentation yet
	VendorName *string `json:"vendorName,omitempty" xmlrpc:"vendorName,omitempty"`

	// no documentation yet
	VlanId *int `json:"vlanId,omitempty" xmlrpc:"vlanId,omitempty"`

	// no documentation yet
	ZoneName *string `json:"zoneName,omitempty" xmlrpc:"zoneName,omitempty"`
}

// The SoftLayer_Network_LBaaS_HealthMonitor type presents a structure containing attributes of a health monitor object associated with load balancer instance. Note that the relationship between backend (pool) and health monitor is N-to-1, especially that the pools object associated with a health monitor must have the same pair of protocol and port. Example: frontend FA: http, 80   - backend BA: tcp, 3456 - healthmonitor HM_tcp3456 frontend FB: https, 443 - backend BB: tcp, 3456 - healthmonitor HM_tcp3456 In above example both backends BA and BB share the same healthmonitor HM_tcp3456
type Network_LBaaS_HealthMonitor struct {
	Entity

	// Create date of the health monitor instance
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// Health monitor's identifier
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Interval in seconds to perform health check
	Interval *int `json:"interval,omitempty" xmlrpc:"interval,omitempty"`

	// Maximum number of health check retries in case of failure
	MaxRetries *int `json:"maxRetries,omitempty" xmlrpc:"maxRetries,omitempty"`

	// Modify date of the health monitor instance
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Type of health check, valid values are "TCP", "HTTP" and "HTTPS"
	MonitorType *string `json:"monitorType,omitempty" xmlrpc:"monitorType,omitempty"`

	// Provisioning status of the health monitor, supported values are "CREATE_PENDING",
	ProvisioningStatus *string `json:"provisioningStatus,omitempty" xmlrpc:"provisioningStatus,omitempty"`

	// Timeout in seconds to wait for health checks response
	Timeout *int `json:"timeout,omitempty" xmlrpc:"timeout,omitempty"`

	// If monitorType is "HTTP" this specifies the whole URL path
	UrlPath *string `json:"urlPath,omitempty" xmlrpc:"urlPath,omitempty"`

	// Health monitor's UUID
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// The SoftLayer_Network_LBaaS_L7HealthMonitor type presents a structure containing attributes of a health monitor object associated with a L7 pool instance. Note that the relationship between backend (L7 pool) and health monitor is 1-to-1, pools object associated with a health monitor must have the same pair of protocol and port. Example: frontend FA: http, 80   - backend BA: http, 3456 - healthmonitor HM_http3456 frontend FB: https, 443 - backend BB: http, 3456 - healthmonitor HM_http3456
type Network_LBaaS_L7HealthMonitor struct {
	Entity

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	Interval *int `json:"interval,omitempty" xmlrpc:"interval,omitempty"`

	// no documentation yet
	MaxRetries *int `json:"maxRetries,omitempty" xmlrpc:"maxRetries,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// no documentation yet
	MonitorType *string `json:"monitorType,omitempty" xmlrpc:"monitorType,omitempty"`

	// no documentation yet
	ProvisioningStatus *string `json:"provisioningStatus,omitempty" xmlrpc:"provisioningStatus,omitempty"`

	// no documentation yet
	Timeout *int `json:"timeout,omitempty" xmlrpc:"timeout,omitempty"`

	// no documentation yet
	UrlPath *string `json:"urlPath,omitempty" xmlrpc:"urlPath,omitempty"`
}

// The SoftLayer_Network_LBaaS_L7Member represents the backend member for a L7 pool. It can be either a virtual server or a bare metal machine.
type Network_LBaaS_L7Member struct {
	Entity

	// The IP address of a L7 pool member.
	Address *string `json:"address,omitempty" xmlrpc:"address,omitempty"`

	// <<< EOT Specifies when a L7 pool member
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The ID of a L7 pool member.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// <<< EOT Specifies when a L7 Pool
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Backends protocol port
	Port *int `json:"port,omitempty" xmlrpc:"port,omitempty"`

	// <<< EOT The provisioning status of a L7 pool member.
	ProvisioningStatus *string `json:"provisioningStatus,omitempty" xmlrpc:"provisioningStatus,omitempty"`

	// The UUID of a L7 pool member.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`

	// The weight of a L7 pool member.
	Weight *int `json:"weight,omitempty" xmlrpc:"weight,omitempty"`
}

// The SoftLayer_Network_LBaaS_L7Policy represents the policy for a listener.
type Network_LBaaS_L7Policy struct {
	Entity

	// The Action to take if the rules belonging to this policy match. It can be set to any of the following values: REDIRECT_URL, REDIRECT_POOL, REDIRECT_HTTPS, REJECT.
	Action *string `json:"action,omitempty" xmlrpc:"action,omitempty"`

	// Specifies when a L7 Policy was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The unique identifier of a policy.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of
	L7RuleCount *uint `json:"l7RuleCount,omitempty" xmlrpc:"l7RuleCount,omitempty"`

	// no documentation yet
	L7Rules []Network_LBaaS_L7Rule `json:"l7Rules,omitempty" xmlrpc:"l7Rules,omitempty"`

	// Specifies when a L7 Policy was updated previously.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Name of a Policy.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The order in which the policy is evaluated. Each policy should have a unique priority
	Priority *int `json:"priority,omitempty" xmlrpc:"priority,omitempty"`

	// The L7 pool id to which traffic is redirected
	RedirectL7PoolId *int `json:"redirectL7PoolId,omitempty" xmlrpc:"redirectL7PoolId,omitempty"`

	// The UUID of the L7 pool object referenced by the policy when the policy action is set to REDIRECT_POOL
	RedirectL7PoolUuid *string `json:"redirectL7PoolUuid,omitempty" xmlrpc:"redirectL7PoolUuid,omitempty"`

	// The URL to which traffic is redirected when the action is set to REDIRECT_URL. Or the port to which listener traffic is redirected to when the action is set to REDIRECT_HTTPS.
	RedirectUrl *string `json:"redirectUrl,omitempty" xmlrpc:"redirectUrl,omitempty"`

	// The UUID of a Policy.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// The SoftLayer_Network_LBaaS_L7Pool type presents a structure containing attributes of a load balancer's L7 pool such as the protocol, and the load balancing algorithm used. L7 pool is used for redirect_pool action of the L7 policy and is different from the default pool
type Network_LBaaS_L7Pool struct {
	Entity

	// Create date of the L7 pool instance
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	L7HealthMonitor *Network_LBaaS_L7HealthMonitor `json:"l7HealthMonitor,omitempty" xmlrpc:"l7HealthMonitor,omitempty"`

	// A count of
	L7MemberCount *uint `json:"l7MemberCount,omitempty" xmlrpc:"l7MemberCount,omitempty"`

	// no documentation yet
	L7Members []Network_LBaaS_L7Member `json:"l7Members,omitempty" xmlrpc:"l7Members,omitempty"`

	// no documentation yet
	L7Policies []Network_LBaaS_L7Policy `json:"l7Policies,omitempty" xmlrpc:"l7Policies,omitempty"`

	// A count of
	L7PolicyCount *uint `json:"l7PolicyCount,omitempty" xmlrpc:"l7PolicyCount,omitempty"`

	// no documentation yet
	L7SessionAffinity *Network_LBaaS_L7SessionAffinity `json:"l7SessionAffinity,omitempty" xmlrpc:"l7SessionAffinity,omitempty"`

	// Load balancing algorithm: "ROUNDROBIN", "WEIGHTED_RR", "LEASTCONNECTION"
	LoadBalancingAlgorithm *string `json:"loadBalancingAlgorithm,omitempty" xmlrpc:"loadBalancingAlgorithm,omitempty"`

	// Last updated date of the L7 pool
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Name of the L7 pool.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// Backends protocol, supported protocol is, "HTTP"
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// Provisioning status of a load balancer's L7 pool.
	ProvisioningStatus *string `json:"provisioningStatus,omitempty" xmlrpc:"provisioningStatus,omitempty"`

	// Instance uuid of the L7 pool
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// SoftLayer_Network_LBaaS_L7PoolMembersHealth provides statistics of members belonging to a particular L7 pool.
type Network_LBaaS_L7PoolMembersHealth struct {
	Entity

	// Instance uuid of the L7 pool
	L7PoolUuid *string `json:"l7PoolUuid,omitempty" xmlrpc:"l7PoolUuid,omitempty"`

	// Members statistics of the L7 pool
	MembersHealth []Network_LBaaS_MemberHealth `json:"membersHealth,omitempty" xmlrpc:"membersHealth,omitempty"`
}

// The SoftLayer_Network_LBaaS_L7Rule represents the Rules that can be attached to a a L7 policy.
type Network_LBaaS_L7Rule struct {
	Entity

	// Comparision type for the Rule, It should any of the following values : REGEX, STARTS_WITH, ENDS_WITH, CONTAINS, EQUAL_TO.
	ComparisonType *string `json:"comparisonType,omitempty" xmlrpc:"comparisonType,omitempty"`

	// Specifies when a Rule was created
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The ID of a Rule.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Inverts the result of the value if set, i.e. True will be inverted to False and vice-versa
	Invert *int `json:"invert,omitempty" xmlrpc:"invert,omitempty"`

	// Key for Rule type HEADER and COOKIE.
	Key *string `json:"key,omitempty" xmlrpc:"key,omitempty"`

	// Specifies when a Rule was updated previously.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Type of the Rule. It  should have any of the following values: HOST_NAME, FILE_TYPE, HEADER, COOKIE, PATH.
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The UUID of a Rule.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`

	// Value for Rule . For type HEADER and COOKIE, this value is compared against the value of the key from HEADER or COOKIE.
	Value *string `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// SoftLayer_Network_LBaaS_L7SessionAffinity represents the session affinity, aka session persistence, configuration for a load balancer backend L7 pool.
type Network_LBaaS_L7SessionAffinity struct {
	Entity

	// no documentation yet
	L7Pool *Network_LBaaS_L7Pool `json:"l7Pool,omitempty" xmlrpc:"l7Pool,omitempty"`

	// Type of the session persistence
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// The SoftLayer_Network_LBaaS_Listener type presents a data structure for a load balancers listener, also called frontend.
type Network_LBaaS_Listener struct {
	Entity

	// maximum idle time in seconds(Range: 1 to 7200), after which the load balancer brings down the
	ClientTimeout *int `json:"clientTimeout,omitempty" xmlrpc:"clientTimeout,omitempty"`

	// Limit of connections a listener can accept
	ConnectionLimit *int `json:"connectionLimit,omitempty" xmlrpc:"connectionLimit,omitempty"`

	// Specifies when the listener was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	DefaultPool *Network_LBaaS_Pool `json:"defaultPool,omitempty" xmlrpc:"defaultPool,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	L7Policies []Network_LBaaS_L7Policy `json:"l7Policies,omitempty" xmlrpc:"l7Policies,omitempty"`

	// A count of
	L7PolicyCount *uint `json:"l7PolicyCount,omitempty" xmlrpc:"l7PolicyCount,omitempty"`

	// Specifies when the listener was updated previously.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Listeners protocol, one of "TCP", "HTTP", "HTTPS".
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// Listeners protocol port number.
	ProtocolPort *int `json:"protocolPort,omitempty" xmlrpc:"protocolPort,omitempty"`

	// The provisioning status of listener.
	ProvisioningStatus *string `json:"provisioningStatus,omitempty" xmlrpc:"provisioningStatus,omitempty"`

	// maximum idle time in seconds(Range: 1 to 7200), after which the load balancer brings down the
	ServerTimeout *int `json:"serverTimeout,omitempty" xmlrpc:"serverTimeout,omitempty"`

	// This references to SSL/TLS certificate (optional) for a listener
	TlsCertificateId *int `json:"tlsCertificateId,omitempty" xmlrpc:"tlsCertificateId,omitempty"`

	// The UUID of a listener.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// The SoftLayer_Network_LBaaS_LoadBalancer type presents a structure containing attributes of a load balancer, and its related objects including listeners, pools and members.
type Network_LBaaS_LoadBalancer struct {
	Entity

	// The account this load balancer belongs to.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// Address (Host name) of a load balancer.
	Address *string `json:"address,omitempty" xmlrpc:"address,omitempty"`

	// Specifies when a load balancer was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// Datacenter, where load balancer is located.
	Datacenter *Location `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// Description of a load balancer.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// A count of health monitors for the backend members.
	HealthMonitorCount *uint `json:"healthMonitorCount,omitempty" xmlrpc:"healthMonitorCount,omitempty"`

	// Health monitors for the backend members.
	HealthMonitors []Network_LBaaS_HealthMonitor `json:"healthMonitors,omitempty" xmlrpc:"healthMonitors,omitempty"`

	// The unique identifier of a load balancer.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Specifies whether the data log is enabled for the load balancer.
	IsDataLogEnabled *int `json:"isDataLogEnabled,omitempty" xmlrpc:"isDataLogEnabled,omitempty"`

	// Specifies whether the load balancer is a public or internal load balancer.
	IsPublic *int `json:"isPublic,omitempty" xmlrpc:"isPublic,omitempty"`

	// A count of l7Pools for load balancer.
	L7PoolCount *uint `json:"l7PoolCount,omitempty" xmlrpc:"l7PoolCount,omitempty"`

	// L7Pools for load balancer.
	L7Pools []Network_LBaaS_L7Pool `json:"l7Pools,omitempty" xmlrpc:"l7Pools,omitempty"`

	// A count of listeners assigned to load balancer.
	ListenerCount *uint `json:"listenerCount,omitempty" xmlrpc:"listenerCount,omitempty"`

	// Listeners assigned to load balancer.
	Listeners []Network_LBaaS_Listener `json:"listeners,omitempty" xmlrpc:"listeners,omitempty"`

	// This references to location with type datacenter
	LocationId *int `json:"locationId,omitempty" xmlrpc:"locationId,omitempty"`

	// A count of members assigned to load balancer.
	MemberCount *uint `json:"memberCount,omitempty" xmlrpc:"memberCount,omitempty"`

	// Members assigned to load balancer.
	Members []Network_LBaaS_Member `json:"members,omitempty" xmlrpc:"members,omitempty"`

	// Specifies when a load balancer was updated last.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The load balancer's name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The operation status "ONLINE" or "OFFLINE" of a load balancer.
	OperatingStatus *string `json:"operatingStatus,omitempty" xmlrpc:"operatingStatus,omitempty"`

	// Error message of previous API call in case of failure
	PreviousErrorText *string `json:"previousErrorText,omitempty" xmlrpc:"previousErrorText,omitempty"`

	// The provisioning status of a load balancer.
	ProvisioningStatus *string `json:"provisioningStatus,omitempty" xmlrpc:"provisioningStatus,omitempty"`

	// A count of list of preferred custom ciphers configured for the load balancer.
	SslCipherCount *uint `json:"sslCipherCount,omitempty" xmlrpc:"sslCipherCount,omitempty"`

	// list of preferred custom ciphers configured for the load balancer.
	SslCiphers []Network_LBaaS_SSLCipher `json:"sslCiphers,omitempty" xmlrpc:"sslCiphers,omitempty"`

	// Specifies the type of load balancer.
	Type *int `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// Applicable for public load balancer only. It specifies whether the public IP addresses are allocated from system public IP pool (1, default) or public subnet (null | 0) from the account ordering the load balancer. For internal load balancer, useSystemPublicIpPool will be ignored, and it always defaults to 1.
	UseSystemPublicIpPool *int `json:"useSystemPublicIpPool,omitempty" xmlrpc:"useSystemPublicIpPool,omitempty"`

	// The UUID of a load balancer.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// This class represents the load balancers appliances, ie virtual servers, on which the actual load balancer service is running. The relationship between load balancer and appliance is 1-to-N with N=2 for beta and very likely N=3 for post beta. Note that this class is for internal use only.
type Network_LBaaS_LoadBalancerAppliance struct {
	Entity

	// no documentation yet
	ComputeId *int `json:"computeId,omitempty" xmlrpc:"computeId,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	LoadBalancerId *int `json:"loadBalancerId,omitempty" xmlrpc:"loadBalancerId,omitempty"`

	// no documentation yet
	OperatingStatus *string `json:"operatingStatus,omitempty" xmlrpc:"operatingStatus,omitempty"`

	// no documentation yet
	PrivateIp *string `json:"privateIp,omitempty" xmlrpc:"privateIp,omitempty"`

	// no documentation yet
	ProvisioningStatus *string `json:"provisioningStatus,omitempty" xmlrpc:"provisioningStatus,omitempty"`

	// no documentation yet
	PublicIp *string `json:"publicIp,omitempty" xmlrpc:"publicIp,omitempty"`

	// no documentation yet
	UnregisteredAt *Time `json:"unregisteredAt,omitempty" xmlrpc:"unregisteredAt,omitempty"`
}

// SoftLayer_Network_LBaaS_LoadBalancerHealthMonitorConfiguration specifies the check method to be used for health monitoring backend members.
type Network_LBaaS_LoadBalancerHealthMonitorConfiguration struct {
	Entity

	// Backends port
	BackendPort *int `json:"backendPort,omitempty" xmlrpc:"backendPort,omitempty"`

	// Backends protocol. Valid values are "TCP", "HTTP"
	BackendProtocol *string `json:"backendProtocol,omitempty" xmlrpc:"backendProtocol,omitempty"`

	// Health Monitor UUID, required for update only
	HealthMonitorUuid *string `json:"healthMonitorUuid,omitempty" xmlrpc:"healthMonitorUuid,omitempty"`

	// <<< EOT Interval in seconds to perform
	Interval *int `json:"interval,omitempty" xmlrpc:"interval,omitempty"`

	// Max number of retries until the member is considered as DOWN
	MaxRetries *int `json:"maxRetries,omitempty" xmlrpc:"maxRetries,omitempty"`

	// Health check methods timeout in
	Timeout *int `json:"timeout,omitempty" xmlrpc:"timeout,omitempty"`

	// If monitor is "HTTP", this specifies URL path
	UrlPath *string `json:"urlPath,omitempty" xmlrpc:"urlPath,omitempty"`
}

// SoftLayer_Network_LBaaS_LoadBalancerMonitoringMetricDataPoint is a collection of datapoints retrieved from a load balancer instance. The available metrics are: <ul> <li>The metric value </li> <li>The timestamp when the metric value was obtained </li> </ul>
type Network_LBaaS_LoadBalancerMonitoringMetricDataPoint struct {
	Entity

	// Epoch Time
	EpochTimestamp *int `json:"epochTimestamp,omitempty" xmlrpc:"epochTimestamp,omitempty"`

	// a value
	Value *Float64 `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// SoftLayer_Network_LBaaS_LoadBalancerProtocolConfiguration specifies the protocol, port, maximum number of allowed connections and session stickiness for load balancer's front- and backend.
type Network_LBaaS_LoadBalancerProtocolConfiguration struct {
	Entity

	// Backends port
	BackendPort *int `json:"backendPort,omitempty" xmlrpc:"backendPort,omitempty"`

	// Backends protocol. Valid values are "TCP", "HTTP"
	BackendProtocol *string `json:"backendProtocol,omitempty" xmlrpc:"backendProtocol,omitempty"`

	// maximum idle time in seconds(Range: 1 to 7200), after which the load balancer brings down the client-side connection
	ClientTimeout *int `json:"clientTimeout,omitempty" xmlrpc:"clientTimeout,omitempty"`

	// Frontends port
	FrontendPort *int `json:"frontendPort,omitempty" xmlrpc:"frontendPort,omitempty"`

	// Frontends protocol. Valid values are "TCP", "HTTP" and "HTTPS"
	FrontendProtocol *string `json:"frontendProtocol,omitempty" xmlrpc:"frontendProtocol,omitempty"`

	// Listeners UUID, required for update only
	ListenerUuid *string `json:"listenerUuid,omitempty" xmlrpc:"listenerUuid,omitempty"`

	// Load balancing method. Valid values are "ROUNDROBIN", "WEIGHTED_RR" and "LEASTCONNECTION"
	LoadBalancingMethod *string `json:"loadBalancingMethod,omitempty" xmlrpc:"loadBalancingMethod,omitempty"`

	// Maximum number of allowed connections
	MaxConn *int `json:"maxConn,omitempty" xmlrpc:"maxConn,omitempty"`

	// maximum idle time in seconds(Range: 1 to 7200), after which the load balancer brings down the server-side connection
	ServerTimeout *int `json:"serverTimeout,omitempty" xmlrpc:"serverTimeout,omitempty"`

	// Sessions cookie name
	SessionCookieName *string `json:"sessionCookieName,omitempty" xmlrpc:"sessionCookieName,omitempty"`

	// Session stickiness type. Valid values are "SOURCE_IP" "HTTP_COOKIE"
	SessionType *string `json:"sessionType,omitempty" xmlrpc:"sessionType,omitempty"`

	// ssl/tls certificate id
	TlsCertificateId *int `json:"tlsCertificateId,omitempty" xmlrpc:"tlsCertificateId,omitempty"`
}

// SoftLayer_Network_LBaaS_LoadBalancerServerInstanceInfo specifies the application server, usually an IBM SoftLayer virtual server or bare metal system, to be assigned to a load balancer.
type Network_LBaaS_LoadBalancerServerInstanceInfo struct {
	Entity

	// Servers private IP address
	PrivateIpAddress *string `json:"privateIpAddress,omitempty" xmlrpc:"privateIpAddress,omitempty"`

	// Servers public IP address
	PublicIpAddress *string `json:"publicIpAddress,omitempty" xmlrpc:"publicIpAddress,omitempty"`

	// Load balancing weight for a server
	Weight *int `json:"weight,omitempty" xmlrpc:"weight,omitempty"`
}

// SoftLayer_Network_LBaaS_LoadBalancerStatistics is a collection of metrics retrieved from a load balancer instance. The available metrics are: <ul> <li>NUmber of members up</li> <li>Number of members down</li> <li>Total number of active connections</li> <li>Throughput</li> <li>Data processed by month</li> <li>Connection rate</li> </ul>
type Network_LBaaS_LoadBalancerStatistics struct {
	Entity

	// Number of connections seen at the
	ConnectionRate *int `json:"connectionRate,omitempty" xmlrpc:"connectionRate,omitempty"`

	// Data processed by month is the total of bin and bout
	DataProcessedByMonth *int `json:"dataProcessedByMonth,omitempty" xmlrpc:"dataProcessedByMonth,omitempty"`

	// Number of members in DOWN health state
	NumberOfMembersDown *int `json:"numberOfMembersDown,omitempty" xmlrpc:"numberOfMembersDown,omitempty"`

	// Number of members in UP health state
	NumberOfMembersUp *int `json:"numberOfMembersUp,omitempty" xmlrpc:"numberOfMembersUp,omitempty"`

	// Throughput measures the total number of bits
	Throughput *Float64 `json:"throughput,omitempty" xmlrpc:"throughput,omitempty"`

	// Number of total active established connections
	TotalConnections *int `json:"totalConnections,omitempty" xmlrpc:"totalConnections,omitempty"`
}

// The SoftLayer_Network_LBaaS_Member represents the backend member for a load balancer. It can be either a virtual server or a bare metal machine.
type Network_LBaaS_Member struct {
	Entity

	// The IP address of a load balancer member.
	Address *string `json:"address,omitempty" xmlrpc:"address,omitempty"`

	// Specifies when a load balancers
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Specifies when a load balancers
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The provisioning status of a load balancer member.
	ProvisioningStatus *string `json:"provisioningStatus,omitempty" xmlrpc:"provisioningStatus,omitempty"`

	// The UUID of a load balancer member.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`

	// The weight of a load balancer member.
	Weight *int `json:"weight,omitempty" xmlrpc:"weight,omitempty"`
}

// SoftLayer_Network_LBaaS_MemberHealth is a collection member metrics retrieved from a LBaaS VSI instance. The available metrics are: <ul> <li>Name of the member</li> <li>Status of the member up or down</li> <li>Uuid of the member</li> </ul>
type Network_LBaaS_MemberHealth struct {
	Entity

	// Members status (UP/DOWN).
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// Members UUID.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// SoftLayer_Network_LBaaS_PolicyRule
//
// This class contains layer 7 policy specifications and an array of associated rules An array of objects of this class must be passed to the API in order to create a policy and its associated rules. <ul> <li>The layer 7 policy object </li> <li>An array of layer 7 rules </li> </ul>
type Network_LBaaS_PolicyRule struct {
	Entity

	// L7 Policy
	L7Policy *Network_LBaaS_L7Policy `json:"l7Policy,omitempty" xmlrpc:"l7Policy,omitempty"`

	// L7 Rules
	L7Rules []Network_LBaaS_L7Rule `json:"l7Rules,omitempty" xmlrpc:"l7Rules,omitempty"`
}

// The SoftLayer_Network_LBaaS_Pool type presents a structure containing attributes of a load balancer pool such as the protocol, protocol port and the load balancing algorithm used.
type Network_LBaaS_Pool struct {
	Entity

	// Create date of the pool instance
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	HealthMonitor *Network_LBaaS_HealthMonitor `json:"healthMonitor,omitempty" xmlrpc:"healthMonitor,omitempty"`

	// Load balancing algorithm: "ROUNDROBIN", "WEIGHTED_RR", "LEASTCONNECTION"
	LoadBalancingAlgorithm *string `json:"loadBalancingAlgorithm,omitempty" xmlrpc:"loadBalancingAlgorithm,omitempty"`

	// A count of
	MemberCount *uint `json:"memberCount,omitempty" xmlrpc:"memberCount,omitempty"`

	// no documentation yet
	Members []Network_LBaaS_Member `json:"members,omitempty" xmlrpc:"members,omitempty"`

	// Last updated date of the pool
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Backends protocol, supported protocols are "TCP", "HTTP" and "HTTPS"
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// Backends protocol port
	ProtocolPort *int `json:"protocolPort,omitempty" xmlrpc:"protocolPort,omitempty"`

	// Provisioning status of a load balancer pool.
	ProvisioningStatus *string `json:"provisioningStatus,omitempty" xmlrpc:"provisioningStatus,omitempty"`

	// no documentation yet
	SessionAffinity *Network_LBaaS_SessionAffinity `json:"sessionAffinity,omitempty" xmlrpc:"sessionAffinity,omitempty"`

	// Instance uuid of the pool
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// SoftLayer_Network_LBaaS_PoolMembersHealth provides statistics of members belonging to a particular pool.
type Network_LBaaS_PoolMembersHealth struct {
	Entity

	// Members statistics of the pool
	MembersHealth []Network_LBaaS_MemberHealth `json:"membersHealth,omitempty" xmlrpc:"membersHealth,omitempty"`

	// Instance uuid of the pool
	PoolUuid *string `json:"poolUuid,omitempty" xmlrpc:"poolUuid,omitempty"`
}

// The SoftLayer_Network_LBaaS_SSLCipher type presents a structure that contains attributes of load balancer cipher suites.
type Network_LBaaS_SSLCipher struct {
	Entity

	// Cipher identifier
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Name of the cipher
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// SoftLayer_Network_LBaaS_SessionAffinity represents the session affinity, aka session persistence, configuration for a load balancer backend pool.
type Network_LBaaS_SessionAffinity struct {
	Entity

	// no documentation yet
	Pool *Network_LBaaS_Pool `json:"pool,omitempty" xmlrpc:"pool,omitempty"`

	// Type of the session persistence
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// The SoftLayer_Network_LoadBalancer_Service data type contains all the information relating to a specific service (destination) on a particular load balancer.
//
// Information retained on the object itself is the the source and destination of the service, routing type, weight, and whether or not the service is currently enabled.
type Network_LoadBalancer_Service struct {
	Entity

	// Connection limit on this service.
	ConnectionLimit *int `json:"connectionLimit,omitempty" xmlrpc:"connectionLimit,omitempty"`

	// Creation Date of this service
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The IP Address of the real server you wish to direct traffic to.  Your account must own this IP
	DestinationIpAddress *string `json:"destinationIpAddress,omitempty" xmlrpc:"destinationIpAddress,omitempty"`

	// The port on the real server to direct the traffic.  This can be different than the source port.  If you wish to obfuscate your HTTP traffic, you can accept requests on port 80 on the load balancer, then redirect them to port 932 on your real server.
	DestinationPort *int `json:"destinationPort,omitempty" xmlrpc:"destinationPort,omitempty"`

	// A flag (either true or false) that determines if this particular service should be enabled on the load balancer.  Set to false to bring the server out of rotation without losing your configuration
	Enabled *bool `json:"enabled,omitempty" xmlrpc:"enabled,omitempty"`

	// The health check type for this service.  If one is supplied, the load balancer will occasionally ping your server to determine if it is still up.  Servers that are down are removed from the queue and will not be used to handle requests until their status returns to "up".  The value of the health check is determined directly by what option you have selected for the routing type.
	//
	// {|
	// |-
	// ! Type
	// ! Valid Health Checks
	// |-
	// | HTTP
	// | HTTP, TCP, ICMP
	// |-
	// | TCP
	// | HTTP, TCP, ICMP
	// |-
	// | FTP
	// | TCP, ICMP
	// |-
	// | DNS
	// | DNS, ICMP
	// |-
	// | UDP
	// | None
	// |}
	//
	//
	HealthCheck *string `json:"healthCheck,omitempty" xmlrpc:"healthCheck,omitempty"`

	// The URL provided here (starting with /) is what the load balancer will request in order to perform a custom HTTP health check.  You must specify either "GET /location/of/file.html" or "HEAD /location/of/file.php"
	HealthCheckURL *string `json:"healthCheckURL,omitempty" xmlrpc:"healthCheckURL,omitempty"`

	// The expected response from the custom HTTP health check.  If the requested page contains this response, the check succeeds.
	HealthResponse *string `json:"healthResponse,omitempty" xmlrpc:"healthResponse,omitempty"`

	// Unique ID for this object, used for the getObject method, and must be set if you are editing this object.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Last modification date of this service
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Name of the load balancer service
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// Holds whether this server is up or down.  Does not affect load balancer configuration at all, just for the customer's informational purposes
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// Peak historical connections since the creation of this service.  Is reset any time you make a configuration change
	PeakConnections *int `json:"peakConnections,omitempty" xmlrpc:"peakConnections,omitempty"`

	// The port on the load balancer that this service maps to.  This is the port for incoming traffic, it needs to be shared with other services to form a group.
	SourcePort *int `json:"sourcePort,omitempty" xmlrpc:"sourcePort,omitempty"`

	// The connection type of this service.  Valid values are HTTP, FTP, TCP, UDP, and DNS.  The value of this variable affects available values of healthCheck, listed in that variable's description
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The load balancer that this service belongs to.
	Vip *Network_LoadBalancer_VirtualIpAddress `json:"vip,omitempty" xmlrpc:"vip,omitempty"`

	// Unique ID for this object's parent.  Probably not useful in the API, as this object will always be a child of a VirtualIpAddress anyway.
	VipId *int `json:"vipId,omitempty" xmlrpc:"vipId,omitempty"`

	// Weight affects the choices the load balancer makes between your services.  The weight of each service is expressed as a percentage of the TOTAL CONNECTION LIMIT on the virtual IP Address.  All services draw from the same pool of connections, so if you expect to have 4 times as much HTTP traffic as HTTPS, your weights for the above example routes would be 40%, 40%, 10%, 10% respectively.  The weights should add up to 100%  If you go over 100%, an exception will be thrown.  Weights must be whole numbers, no fractions or decimals are accepted.
	Weight *int `json:"weight,omitempty" xmlrpc:"weight,omitempty"`
}

// The SoftLayer_Network_LoadBalancer_VirtualIpAddress data type contains all the information relating to a specific load balancer assigned to a customer account.
//
// Information retained on the object itself is the virtual IP address, load balancing method, and any notes that are related to the load balancer.  There is also an array of SoftLayer_Network_LoadBalancer_Service objects, which represent the load balancer services, explained more fully in the SoftLayer_Network_LoadBalancer_Service documentation.
type Network_LoadBalancer_VirtualIpAddress struct {
	Entity

	// The account that owns this load balancer.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The current billing item for the Load Balancer.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// Connection limit on this VIP.  Can be upgraded through the upgradeConnectionLimit() function
	ConnectionLimit *int `json:"connectionLimit,omitempty" xmlrpc:"connectionLimit,omitempty"`

	// If false, this VIP and associated services may be edited via the portal or the API. If true, you must configure this VIP manually on the device.
	CustomerManagedFlag *int `json:"customerManagedFlag,omitempty" xmlrpc:"customerManagedFlag,omitempty"`

	// Unique ID for this object, used for the getObject method, and must be set if you are editing this object.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The load balancing method that determines which server is used "next" by the load balancer.  The method is stored in an abbreviated form, represented in parentheses after the full name. Methods include: Round Robin (Value "rr"):  Each server is used sequentially in a circular queue Shortest Response (Value "sr"):  The server with the lowest ping at the last health check gets the next request Least Connections (Value "lc"):  The server with the least current connections is given the next request Persistent IP - Round Robin (Value "pi"): The same server will be returned to a request during a users session.  Servers are chosen through round robin. Persistent IP - Shortest Response (Value "pi-sr"): The same server will be returned to a request during a users session.  Servers are chosen through shortest response. Persistent IP - Least Connections (Value "pi-lc"): The same server will be returned to a request during a users session.  Servers are chosen through least connections. Insert Cookie - Round Robin (Value "ic"):  Inserts a cookie into the HTTP stream that will tie that client to a particular balanced server. Servers are chosen through round robin. Insert Cookie - Shortest Response (Value "ic-sr"): Inserts a cookie into the HTTP stream that will tie that client to a particular balanced server. Servers are chosen through shortest response. Insert Cookie - Least Connections (Value "ic-lc"): Inserts a cookie into the HTTP stream that will tie that client to a particular balanced server. Servers are chosen through least connections.
	LoadBalancingMethod *string `json:"loadBalancingMethod,omitempty" xmlrpc:"loadBalancingMethod,omitempty"`

	// A human readable version of loadBalancingMethod, intended mainly for API users.
	LoadBalancingMethodFullName *string `json:"loadBalancingMethodFullName,omitempty" xmlrpc:"loadBalancingMethodFullName,omitempty"`

	// A flag indicating that the load balancer is a managed resource.
	ManagedResourceFlag *bool `json:"managedResourceFlag,omitempty" xmlrpc:"managedResourceFlag,omitempty"`

	// Date this load balancer was last modified
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The name of the load balancer instance
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// User-created notes on this load balancer.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// The unique identifier of the Security Certificate to be utilized when SSL support is enabled.
	SecurityCertificateId *int `json:"securityCertificateId,omitempty" xmlrpc:"securityCertificateId,omitempty"`

	// A count of the services on this load balancer.
	ServiceCount *uint `json:"serviceCount,omitempty" xmlrpc:"serviceCount,omitempty"`

	// the services on this load balancer.
	Services []Network_LoadBalancer_Service `json:"services,omitempty" xmlrpc:"services,omitempty"`

	// This is the port for incoming traffic.
	SourcePort *int `json:"sourcePort,omitempty" xmlrpc:"sourcePort,omitempty"`

	// The connection type of this VIP.  Valid values are HTTP, FTP, TCP, UDP, and DNS.
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The virtual, public-facing IP address for your load balancer.  This is the address of all incoming traffic
	VirtualIpAddress *string `json:"virtualIpAddress,omitempty" xmlrpc:"virtualIpAddress,omitempty"`
}

// The Syslog class holds a single line from the Networking Firewall "Syslog" record, for firewall detected and blocked attempts on a server.
type Network_Logging_Syslog struct {
	Entity

	// Timestamp for when the connection was blocked by the firewall
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The Destination IP Address of the blocked connection (your end)
	DestinationIpAddress *string `json:"destinationIpAddress,omitempty" xmlrpc:"destinationIpAddress,omitempty"`

	// The Destination Port of the blocked connection (your end)
	DestinationPort *int `json:"destinationPort,omitempty" xmlrpc:"destinationPort,omitempty"`

	// This tells you what kind of firewall event this log line is for:  accept or deny.
	EventType *string `json:"eventType,omitempty" xmlrpc:"eventType,omitempty"`

	// Raw syslog message for the event
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// Connection protocol used to make the call that was blocked (tcp, udp, etc)
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// The Source IP Address of the call that was blocked (attacker's end)
	SourceIpAddress *string `json:"sourceIpAddress,omitempty" xmlrpc:"sourceIpAddress,omitempty"`

	// The Source Port where the blocked connection was established (attacker's end)
	SourcePort *int `json:"sourcePort,omitempty" xmlrpc:"sourcePort,omitempty"`

	// If this is an aggregation of syslog events, this property shows the total events.
	TotalEvents *int `json:"totalEvents,omitempty" xmlrpc:"totalEvents,omitempty"`
}

// no documentation yet
type Network_Message_Delivery struct {
	Entity

	// The SoftLayer customer account that a network message delivery account belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// no documentation yet
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The billing item for a network message delivery account.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	Guid *string `json:"guid,omitempty" xmlrpc:"guid,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// no documentation yet
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// The message delivery type of a network message delivery account.
	Type *Network_Message_Delivery_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// no documentation yet
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`

	// no documentation yet
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`

	// The vendor for a network message delivery account.
	Vendor *Network_Message_Delivery_Vendor `json:"vendor,omitempty" xmlrpc:"vendor,omitempty"`

	// no documentation yet
	VendorId *int `json:"vendorId,omitempty" xmlrpc:"vendorId,omitempty"`
}

// no documentation yet
type Network_Message_Delivery_Attribute struct {
	Entity

	// no documentation yet
	NetworkMessageDelivery *Network_Message_Delivery `json:"networkMessageDelivery,omitempty" xmlrpc:"networkMessageDelivery,omitempty"`

	// no documentation yet
	Value *string `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// no documentation yet
type Network_Message_Delivery_Email_Sendgrid struct {
	Network_Message_Delivery

	// The contact e-mail address used by SendGrid.
	EmailAddress *string `json:"emailAddress,omitempty" xmlrpc:"emailAddress,omitempty"`

	// A flag that determines if a SendGrid e-mail delivery account has access to send mail through the SendGrid SMTP server.
	SmtpAccess *string `json:"smtpAccess,omitempty" xmlrpc:"smtpAccess,omitempty"`
}

// no documentation yet
type Network_Message_Delivery_Type struct {
	Entity

	// no documentation yet
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Message_Delivery_Vendor struct {
	Entity

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Monitor struct {
	Entity
}

// The SoftLayer_Network_Monitor_Version1_Incident data type models a single virtual server or physical hardware network monitoring event. SoftLayer_Network_Monitor_Version1_Incidents are created when the SoftLayer monitoring system detects a service down on your hardware of virtual server. As the incident is resolved it's status changes from "SERVICE FAILURE" to "COMPLETED".
type Network_Monitor_Version1_Incident struct {
	Entity

	// A network monitoring incident's status, either the string "SERVICE FAILURE" denoting an ongoing incident or "COMPLETE" meaning the incident has been resolved.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`
}

// The Monitoring_Query_Host type represents a monitoring instance.  It consists of a hardware ID to monitor, an IP address attached to that hardware ID, a method of monitoring, and what to do in the instance that the monitor ever fails.
type Network_Monitor_Version1_Query_Host struct {
	Entity

	// The argument to be used for this monitor, if necessary.  The lowest monitoring levels (like ping) ignore this setting, but higher levels like HTTP custom use it.
	Arg1Value *string `json:"arg1Value,omitempty" xmlrpc:"arg1Value,omitempty"`

	// Virtual Guest Identification Number for the guest being monitored.
	GuestId *int `json:"guestId,omitempty" xmlrpc:"guestId,omitempty"`

	// The hardware that is being monitored by this monitoring instance
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// The ID of the hardware being monitored
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// Identification Number for the host being monitored.
	HostId *int `json:"hostId,omitempty" xmlrpc:"hostId,omitempty"`

	// The unique identifier for this object
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The IP address to be monitored.  Must be attached to the hardware on this object
	IpAddress *string `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// The most recent result for this particular monitoring instance.
	LastResult *Network_Monitor_Version1_Query_Result `json:"lastResult,omitempty" xmlrpc:"lastResult,omitempty"`

	// The type of monitoring query that is executed when this hardware is monitored.
	QueryType *Network_Monitor_Version1_Query_Type `json:"queryType,omitempty" xmlrpc:"queryType,omitempty"`

	// The ID of the query type to use.
	QueryTypeId *int `json:"queryTypeId,omitempty" xmlrpc:"queryTypeId,omitempty"`

	// The action taken when a monitor fails.
	ResponseAction *Network_Monitor_Version1_Query_ResponseType `json:"responseAction,omitempty" xmlrpc:"responseAction,omitempty"`

	// The ID of the response action to take when the monitor fails
	ResponseActionId *int `json:"responseActionId,omitempty" xmlrpc:"responseActionId,omitempty"`

	// The status of this monitoring instance.  Anything other than "ON" means that the monitor has been disabled
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// The number of 5-minute cycles to wait before the "responseAction" is taken.  If set to 0, the response action will be taken immediately
	WaitCycles *int `json:"waitCycles,omitempty" xmlrpc:"waitCycles,omitempty"`
}

// The monitoring stratum type stores the maximum level of the various components of the monitoring system that a particular hardware object has access to.  This object cannot be accessed by ID, and cannot be modified. The user can access this object through Hardware_Server->availableMonitoring.
//
// There are two values on this object that are important:
// # monitorLevel determines the highest level of SoftLayer_Network_Monitor_Version1_Query_Type object that can be placed in a monitoring instance on this server
// # responseLevel determines the highest level of SoftLayer_Network_Monitor_Version1_Query_ResponseType object that can be placed in a monitoring instance on this server
//
// Also note that the query type and response types are available through getAllQueryTypes and getAllResponseTypes, respectively.
type Network_Monitor_Version1_Query_Host_Stratum struct {
	Entity

	// The hardware object that these monitoring permissions applies to.
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// The highest level of a monitoring query type allowed on this server
	MonitorLevel *int `json:"monitorLevel,omitempty" xmlrpc:"monitorLevel,omitempty"`

	// The highest level of a monitoring response type allowed on this server
	ResponseLevel *int `json:"responseLevel,omitempty" xmlrpc:"responseLevel,omitempty"`
}

// The ResponseType type stores only an ID and a description of the response type.  The only use for this object is in reference.  The user chooses a response action that would be appropriate for a monitoring instance, and sets the ResponseTypeId to the SoftLayer_Network_Monitor_Version1_Query_Host->responseActionId value.
//
// The user can retrieve all available ResponseTypes with the getAllObjects method on this service.
type Network_Monitor_Version1_Query_ResponseType struct {
	Entity

	// The description of the action the monitoring system will take on failure
	ActionDescription *string `json:"actionDescription,omitempty" xmlrpc:"actionDescription,omitempty"`

	// The unique identifier for this object
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The level of this response.  The level the customer has access to is determined by values in SoftLayer_Network_Monitor_Version1_Query_Host_Stratum
	Level *int `json:"level,omitempty" xmlrpc:"level,omitempty"`
}

// The monitoring result object is used to show the status of the actions taken by the monitoring system.
//
// In general, only the responseStatus variable is needed, as it holds the information on the status of the service.
type Network_Monitor_Version1_Query_Result struct {
	Entity

	// The timestamp of when this monitor was co
	FinishTime *Time `json:"finishTime,omitempty" xmlrpc:"finishTime,omitempty"`

	// References the queryHost that this response relates to.
	QueryHost *Network_Monitor_Version1_Query_Host `json:"queryHost,omitempty" xmlrpc:"queryHost,omitempty"`

	// The response status for this server.  The response status meanings are: 0:  Down/Critical: Server is down and/or has passed the critical response threshold (extremely long ping response, abnormal behavior, etc.) 1:  Warning - Server may be recovering from a previous down state, or may have taken too long to respond 2:  Up 3:  Not used 4:  Unknown - An unknown error has occurred.  If the problem persists, contact support. 5:  Unknown - An unknown error has occurred.  If the problem persists, contact support.
	ResponseStatus *int `json:"responseStatus,omitempty" xmlrpc:"responseStatus,omitempty"`

	// The length of time it took the server to respond
	ResponseTime *Float64 `json:"responseTime,omitempty" xmlrpc:"responseTime,omitempty"`
}

// The MonitorType type stores a name, long description, and default arguments for the monitor types.  The only use for this object is in reference.  The user chooses a monitoring type that would be appropriate for their server, and sets the id of the Query_Type to SoftLayer_Network_Monitor_Version1_Query_Host->queryTypeId
//
// The user can retrieve all available Query Types with the getAllObjects method on this service.
type Network_Monitor_Version1_Query_Type struct {
	Entity

	// The type of parameter sent to the monitoring command.
	ArgumentDescription *string `json:"argumentDescription,omitempty" xmlrpc:"argumentDescription,omitempty"`

	// Long description of the monitoring type.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The unique identifier for this object
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The level of this monitoring type.  The level the customer has access to is determined by values in SoftLayer_Network_Monitor_Version1_Query_Host_Stratum
	MonitorLevel *int `json:"monitorLevel,omitempty" xmlrpc:"monitorLevel,omitempty"`

	// Short name of the monitoring type
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Pod struct {
	Entity

	// Identifier for this Pod's Backend Customer Router (BCR)
	BackendRouterId *int `json:"backendRouterId,omitempty" xmlrpc:"backendRouterId,omitempty"`

	// Host name of Pod's Backend Customer Router (BCR), e.g. bcr01a.dal09
	BackendRouterName *string `json:"backendRouterName,omitempty" xmlrpc:"backendRouterName,omitempty"`

	// Property providing a means to filter Pods based on available capabitilies. See [[SoftLayer_Network_Pod/getAllObjects]] to filter for Pods with specific capabilities. See [[SoftLayer_Network_Pod/getCapabilities]] to retrieve capabilities of a specific Pod.
	Capabilities []string `json:"capabilities,omitempty" xmlrpc:"capabilities,omitempty"`

	// Identifier for the Data Center the Pod resides within
	DatacenterId *int `json:"datacenterId,omitempty" xmlrpc:"datacenterId,omitempty"`

	// Long form name of the data center in which this Pod resides, e.g. Dallas 9
	DatacenterLongName *string `json:"datacenterLongName,omitempty" xmlrpc:"datacenterLongName,omitempty"`

	// Name of data center in which this Pod resides, e.g. dal09
	DatacenterName *string `json:"datacenterName,omitempty" xmlrpc:"datacenterName,omitempty"`

	// (optional) Identifier for this Pod's Frontend Customer Router (FCR)
	FrontendRouterId *int `json:"frontendRouterId,omitempty" xmlrpc:"frontendRouterId,omitempty"`

	// (optional) Host name of Pod's Frontend Customer Router (FCR), e.g. fcr01a.dal09
	FrontendRouterName *string `json:"frontendRouterName,omitempty" xmlrpc:"frontendRouterName,omitempty"`

	// The unique name of the Pod. See [[SoftLayer_Network_Pod (type)]] for details of the name's construction.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Protection_Address struct {
	Entity

	// no documentation yet
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// no documentation yet
	DepartmentId *int `json:"departmentId,omitempty" xmlrpc:"departmentId,omitempty"`

	// no documentation yet
	IpAddress *string `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// no documentation yet
	Location *Location `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// no documentation yet
	ManagementMethodType *string `json:"managementMethodType,omitempty" xmlrpc:"managementMethodType,omitempty"`

	// no documentation yet
	ModifiedUser *User_Employee `json:"modifiedUser,omitempty" xmlrpc:"modifiedUser,omitempty"`

	// no documentation yet
	PrimaryRouter *Hardware_Router `json:"primaryRouter,omitempty" xmlrpc:"primaryRouter,omitempty"`

	// DEPRECATED
	// Deprecated: This function has been marked as deprecated.
	ServiceProvider *Service_Provider `json:"serviceProvider,omitempty" xmlrpc:"serviceProvider,omitempty"`

	// no documentation yet
	Subnet *Network_Subnet `json:"subnet,omitempty" xmlrpc:"subnet,omitempty"`

	// no documentation yet
	SubnetIpAddress *Network_Subnet_IpAddress `json:"subnetIpAddress,omitempty" xmlrpc:"subnetIpAddress,omitempty"`

	// no documentation yet
	TerminatedUser *User_Employee `json:"terminatedUser,omitempty" xmlrpc:"terminatedUser,omitempty"`

	// no documentation yet
	Ticket *Ticket `json:"ticket,omitempty" xmlrpc:"ticket,omitempty"`

	// A count of
	TransactionCount *uint `json:"transactionCount,omitempty" xmlrpc:"transactionCount,omitempty"`

	// no documentation yet
	Transactions []Provisioning_Version1_Transaction `json:"transactions,omitempty" xmlrpc:"transactions,omitempty"`

	// no documentation yet
	UserDepartment *User_Employee_Department `json:"userDepartment,omitempty" xmlrpc:"userDepartment,omitempty"`

	// no documentation yet
	UserRecord *User_Employee `json:"userRecord,omitempty" xmlrpc:"userRecord,omitempty"`
}

// Regional Internet Registries are the organizations who delegate IP address blocks to other groups or organizations around the Internet. The information contained in this data type is used throughout the networking-related services in our systems.
type Network_Regional_Internet_Registry struct {
	Entity

	// Unique ID of the object
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The system-level name of the registry
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The friendly name of the registry
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The SoftLayer_Network_SecurityGroup data type contains general information for a single security group. A security group contains a set of IP filter [[SoftLayer_Network_SecurityGroup_Rule (type)|rules]] that define how to handle incoming (ingress) and outgoing (egress) traffic to both the public and private interfaces of a virtual server instance and a set of [[SoftLayer_Virtual_Network_SecurityGroup_NetworkComponentBinding (type)|bindings]] to associate virtual guest network components with the security group.
type Network_SecurityGroup struct {
	Entity

	// The account this security group belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The date a security group was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The (optional) description for a security group.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The unique ID for a security group.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	Metadata *string `json:"metadata,omitempty" xmlrpc:"metadata,omitempty"`

	// The date a security group was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The (optional) name for a security group.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A count of the network component bindings for this security group.
	NetworkComponentBindingCount *uint `json:"networkComponentBindingCount,omitempty" xmlrpc:"networkComponentBindingCount,omitempty"`

	// The network component bindings for this security group.
	NetworkComponentBindings []Virtual_Network_SecurityGroup_NetworkComponentBinding `json:"networkComponentBindings,omitempty" xmlrpc:"networkComponentBindings,omitempty"`

	// A count of the order bindings for this security group
	OrderBindingCount *uint `json:"orderBindingCount,omitempty" xmlrpc:"orderBindingCount,omitempty"`

	// The order bindings for this security group
	OrderBindings []Network_SecurityGroup_OrderBinding `json:"orderBindings,omitempty" xmlrpc:"orderBindings,omitempty"`

	// A count of the rules for this security group.
	RuleCount *uint `json:"ruleCount,omitempty" xmlrpc:"ruleCount,omitempty"`

	// The rules for this security group.
	Rules []Network_SecurityGroup_Rule `json:"rules,omitempty" xmlrpc:"rules,omitempty"`
}

// The SoftLayer_Network_SecurityGroup_OrderBinding data type contains links between security groups and product orders.
type Network_SecurityGroup_OrderBinding struct {
	Entity

	// The virtual guest associated with the binding
	Guest *Virtual_Guest `json:"guest,omitempty" xmlrpc:"guest,omitempty"`

	// The ID of the Virtual Guest associated with the security group.
	GuestId *int `json:"guestId,omitempty" xmlrpc:"guestId,omitempty"`

	// The unique ID for a security group, order, binding
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The order associated with the binding
	Order *Billing_Order `json:"order,omitempty" xmlrpc:"order,omitempty"`

	// The ID of the order associated with the security group.
	OrderId *int `json:"orderId,omitempty" xmlrpc:"orderId,omitempty"`

	// The security group associated with the order
	SecurityGroup *Network_SecurityGroup `json:"securityGroup,omitempty" xmlrpc:"securityGroup,omitempty"`

	// The ID of the security group that is associated with the order.
	SecurityGroupId *int `json:"securityGroupId,omitempty" xmlrpc:"securityGroupId,omitempty"`
}

// The SoftLayer_Network_SecurityGroup_Request data type contains the ID of a specific request sent to the API. This ID is used to identify specific calls to attach and detach network components, as well as add, edit, and remove security group rules.
type Network_SecurityGroup_Request struct {
	Entity

	// The unique ID for a request.
	RequestId *string `json:"requestId,omitempty" xmlrpc:"requestId,omitempty"`
}

// The SoftLayer_Network_SecurityGroup_RequestRules data type contains the ID of a specific request sent to the API, as well as an associative array of the rules that were created, edited, or removed by the request.
type Network_SecurityGroup_RequestRules struct {
	Network_SecurityGroup_Request

	// Whether the API call was valid or not.
	Rules []Network_SecurityGroup_Rule `json:"rules,omitempty" xmlrpc:"rules,omitempty"`
}

// The SoftLayer_Network_SecurityGroup_Rule data type contains general information for a single rule that belongs to a [[SoftLayer_Network_SecurityGroup|security group]]. By default, all traffic (both inbound and  outbound) to a virtual server instance is blocked. Security group rules are permissive, and define the allowed incoming (ingress) and outgoing (egress) traffic to both the public and private interfaces of a  virtual server instance. The order of rules within a security group does not matter and priority always falls to the least restrictive rule.
type Network_SecurityGroup_Rule struct {
	Entity

	// The createDate field for a rule. It is essentially the date and time that the security group rule was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The direction of traffic (ingress or egress).
	Direction *string `json:"direction,omitempty" xmlrpc:"direction,omitempty"`

	// IPv4 or IPv6. If the remoteIp or ethertype properties are not specified, the default is IPv4. Otherwise ethertype will default based on the format of the specified remoteIp.
	Ethertype *string `json:"ethertype,omitempty" xmlrpc:"ethertype,omitempty"`

	// The unique ID for a rule.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The modifyDate field for a rule. It is essentially the date and time that the security group rule was last changed.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The end of the port range for allowed traffic.  When the protocol is icmp, this value specifies the icmp code to permit.  When icmp code is specified, icmp type is required. When the protocol is vrrp, ports cannot be specified.
	PortRangeMax *int `json:"portRangeMax,omitempty" xmlrpc:"portRangeMax,omitempty"`

	// The start of the port range for allowed traffic.  When the protocol is icmp, this value specifies the icmp type to permit.
	PortRangeMin *int `json:"portRangeMin,omitempty" xmlrpc:"portRangeMin,omitempty"`

	// The protocol of packets (icmp, tcp, udp, or vrrp).
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// The remote security group allowed as part of this rule.
	RemoteGroup *Network_SecurityGroup `json:"remoteGroup,omitempty" xmlrpc:"remoteGroup,omitempty"`

	// The ID of the remote security group allowed as part of the rule. This property is mutually exclusive with the remoteIp property.
	RemoteGroupId *int `json:"remoteGroupId,omitempty" xmlrpc:"remoteGroupId,omitempty"`

	// CIDR or IP address for allowed connections. This property is mutually exclusive with the remoteGroupId property. When the protocol is vrrp, ports cannot be specified.
	RemoteIp *string `json:"remoteIp,omitempty" xmlrpc:"remoteIp,omitempty"`

	// The security group of this rule.
	SecurityGroup *Network_SecurityGroup `json:"securityGroup,omitempty" xmlrpc:"securityGroup,omitempty"`

	// The ID of the security group that owns the rule.
	SecurityGroupId *int `json:"securityGroupId,omitempty" xmlrpc:"securityGroupId,omitempty"`
}

// The SoftLayer_Network_Security_Scanner_Request data type represents a single vulnerability scan request. It provides information on when the scan was created, last updated, and the current status. The status messages are as follows:
// *Scan Pending
// *Scan Processing
// *Scan Complete
// *Scan Cancelled
// *Generating Report.
type Network_Security_Scanner_Request struct {
	Entity

	// The account associated with a security scan request.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// A request's associated customer account identifier.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The date and time that the request is created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The virtual guest a security scan is run against.
	Guest *Virtual_Guest `json:"guest,omitempty" xmlrpc:"guest,omitempty"`

	// Virtual Guest Identification Number for the guest this security scanner request belongs to.
	GuestId *int `json:"guestId,omitempty" xmlrpc:"guestId,omitempty"`

	// The hardware a security scan is run against.
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// The identifier of the hardware item a scan is run on.
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// Identification Number for the host this security scanner request belongs to.
	HostId *int `json:"hostId,omitempty" xmlrpc:"hostId,omitempty"`

	// A security scan request's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The IP address that a scan will be performed on.
	IpAddress *string `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// The date and time that the request was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Flag whether the requestor owns the hardware the scan was run on. This flag will return for hardware servers only, virtual servers will result in a null return even if you have a request out for them.
	RequestorOwnedFlag *bool `json:"requestorOwnedFlag,omitempty" xmlrpc:"requestorOwnedFlag,omitempty"`

	// A security scan request's status.
	Status *Network_Security_Scanner_Request_Status `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// A request status identifier.
	StatusId *int `json:"statusId,omitempty" xmlrpc:"statusId,omitempty"`
}

// The SoftLayer_Network_Security_Scanner_Request_Status data type represents the current status of a vulnerability scan. The status messages are as follows:
// *Scan Pending
// *Scan Processing
// *Scan Complete
// *Scan Cancelled
// *Generating Report.
//
// The status of a vulnerability scan will change over the course of a scan's execution.
type Network_Security_Scanner_Request_Status struct {
	Entity

	// The identifier of a vulnerability scan's status.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The status message of a vulnerability scan.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The SoftLayer_Network_Service_Resource is used to store information related to a service.  It is used for determining the correct resource to connect to for a given service, like NAS, Evault, etc.
type Network_Service_Resource struct {
	Entity

	// no documentation yet
	ApiHost *string `json:"apiHost,omitempty" xmlrpc:"apiHost,omitempty"`

	// no documentation yet
	ApiPassword *string `json:"apiPassword,omitempty" xmlrpc:"apiPassword,omitempty"`

	// no documentation yet
	ApiPath *string `json:"apiPath,omitempty" xmlrpc:"apiPath,omitempty"`

	// no documentation yet
	ApiPort *string `json:"apiPort,omitempty" xmlrpc:"apiPort,omitempty"`

	// no documentation yet
	ApiProtocol *string `json:"apiProtocol,omitempty" xmlrpc:"apiProtocol,omitempty"`

	// no documentation yet
	ApiUsername *string `json:"apiUsername,omitempty" xmlrpc:"apiUsername,omitempty"`

	// no documentation yet
	ApiVersion *string `json:"apiVersion,omitempty" xmlrpc:"apiVersion,omitempty"`

	// A count of
	AttributeCount *uint `json:"attributeCount,omitempty" xmlrpc:"attributeCount,omitempty"`

	// no documentation yet
	Attributes []Network_Service_Resource_Attribute `json:"attributes,omitempty" xmlrpc:"attributes,omitempty"`

	// The backend IP address for this resource
	BackendIpAddress *string `json:"backendIpAddress,omitempty" xmlrpc:"backendIpAddress,omitempty"`

	// no documentation yet
	Datacenter *Location `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// The frontend IP address for this resource
	FrontendIpAddress *string `json:"frontendIpAddress,omitempty" xmlrpc:"frontendIpAddress,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The name associated with this resource
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	NetappVolumeName *string `json:"netappVolumeName,omitempty" xmlrpc:"netappVolumeName,omitempty"`

	// The hardware information associated with this resource.
	NetworkDevice *Hardware `json:"networkDevice,omitempty" xmlrpc:"networkDevice,omitempty"`

	// no documentation yet
	SshUsername *string `json:"sshUsername,omitempty" xmlrpc:"sshUsername,omitempty"`

	// The network information associated with this resource.
	Type *Network_Service_Resource_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// no documentation yet
type Network_Service_Resource_Attribute struct {
	Entity

	// no documentation yet
	AttributeType *Network_Service_Resource_Attribute_Type `json:"attributeType,omitempty" xmlrpc:"attributeType,omitempty"`

	// no documentation yet
	ServiceResource *Network_Service_Resource `json:"serviceResource,omitempty" xmlrpc:"serviceResource,omitempty"`

	// no documentation yet
	Value *string `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// no documentation yet
type Network_Service_Resource_Attribute_Type struct {
	Entity

	// no documentation yet
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`
}

// The SoftLayer_Network_Service_Resource_CosStor is used to store information related to COS service.
type Network_Service_Resource_CosStor struct {
	Network_Service_Resource
}

// no documentation yet
type Network_Service_Resource_Hub struct {
	Network_Service_Resource
}

// no documentation yet
type Network_Service_Resource_Hub_Swift struct {
	Network_Service_Resource_Hub
}

// no documentation yet
type Network_Service_Resource_Type struct {
	Entity

	// A count of
	ServiceResourceCount *uint `json:"serviceResourceCount,omitempty" xmlrpc:"serviceResourceCount,omitempty"`

	// no documentation yet
	ServiceResources []Network_Service_Resource `json:"serviceResources,omitempty" xmlrpc:"serviceResources,omitempty"`

	// no documentation yet
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// The SoftLayer_Network_Service_Vpn_Overrides data type contains information relating user ids to subnet ids when VPN access is manually configured.  It is essentially an entry in a 'white list' of subnets a SoftLayer portal VPN user may access.
type Network_Service_Vpn_Overrides struct {
	Entity

	// The internal identifier of the record.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Subnet components accessible by a SoftLayer VPN portal user.
	Subnet *Network_Subnet `json:"subnet,omitempty" xmlrpc:"subnet,omitempty"`

	// The identifier of a subnet accessible by the SoftLayer portal VPN user.
	SubnetId *int `json:"subnetId,omitempty" xmlrpc:"subnetId,omitempty"`

	// SoftLayer VPN portal user.
	User *User_Customer `json:"user,omitempty" xmlrpc:"user,omitempty"`

	// The identifier of the SoftLayer portal VPN user.
	UserId *int `json:"userId,omitempty" xmlrpc:"userId,omitempty"`
}

// The SoftLayer_Network_Storage data type contains general information regarding a Storage product such as account id, access username and password, the Storage product type, and the server the Storage service is associated with. Currently, only EVault backup storage has an associated server.
type Network_Storage struct {
	Entity

	// The account that a Storage services belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The internal identifier of the SoftLayer customer account that a Storage account belongs to.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// Other usernames and passwords associated with a Storage volume.
	AccountPassword *Account_Password `json:"accountPassword,omitempty" xmlrpc:"accountPassword,omitempty"`

	// A count of the currently active transactions on a network storage volume.
	ActiveTransactionCount *uint `json:"activeTransactionCount,omitempty" xmlrpc:"activeTransactionCount,omitempty"`

	// The currently active transactions on a network storage volume.
	ActiveTransactions []Provisioning_Version1_Transaction `json:"activeTransactions,omitempty" xmlrpc:"activeTransactions,omitempty"`

	// no documentation yet
	AllowDisasterRecoveryFailback *string `json:"allowDisasterRecoveryFailback,omitempty" xmlrpc:"allowDisasterRecoveryFailback,omitempty"`

	// no documentation yet
	AllowDisasterRecoveryFailover *string `json:"allowDisasterRecoveryFailover,omitempty" xmlrpc:"allowDisasterRecoveryFailover,omitempty"`

	// The SoftLayer_Hardware objects which are allowed access to this storage volume.
	AllowedHardware []Hardware `json:"allowedHardware,omitempty" xmlrpc:"allowedHardware,omitempty"`

	// A count of the SoftLayer_Hardware objects which are allowed access to this storage volume.
	AllowedHardwareCount *uint `json:"allowedHardwareCount,omitempty" xmlrpc:"allowedHardwareCount,omitempty"`

	// A count of the SoftLayer_Network_Subnet_IpAddress objects which are allowed access to this storage volume.
	AllowedIpAddressCount *uint `json:"allowedIpAddressCount,omitempty" xmlrpc:"allowedIpAddressCount,omitempty"`

	// The SoftLayer_Network_Subnet_IpAddress objects which are allowed access to this storage volume.
	AllowedIpAddresses []Network_Subnet_IpAddress `json:"allowedIpAddresses,omitempty" xmlrpc:"allowedIpAddresses,omitempty"`

	// The SoftLayer_Hardware objects which are allowed access to this storage volume's Replicant.
	AllowedReplicationHardware []Hardware `json:"allowedReplicationHardware,omitempty" xmlrpc:"allowedReplicationHardware,omitempty"`

	// A count of the SoftLayer_Hardware objects which are allowed access to this storage volume's Replicant.
	AllowedReplicationHardwareCount *uint `json:"allowedReplicationHardwareCount,omitempty" xmlrpc:"allowedReplicationHardwareCount,omitempty"`

	// A count of the SoftLayer_Network_Subnet_IpAddress objects which are allowed access to this storage volume's Replicant.
	AllowedReplicationIpAddressCount *uint `json:"allowedReplicationIpAddressCount,omitempty" xmlrpc:"allowedReplicationIpAddressCount,omitempty"`

	// The SoftLayer_Network_Subnet_IpAddress objects which are allowed access to this storage volume's Replicant.
	AllowedReplicationIpAddresses []Network_Subnet_IpAddress `json:"allowedReplicationIpAddresses,omitempty" xmlrpc:"allowedReplicationIpAddresses,omitempty"`

	// A count of the SoftLayer_Network_Subnet objects which are allowed access to this storage volume's Replicant.
	AllowedReplicationSubnetCount *uint `json:"allowedReplicationSubnetCount,omitempty" xmlrpc:"allowedReplicationSubnetCount,omitempty"`

	// The SoftLayer_Network_Subnet objects which are allowed access to this storage volume's Replicant.
	AllowedReplicationSubnets []Network_Subnet `json:"allowedReplicationSubnets,omitempty" xmlrpc:"allowedReplicationSubnets,omitempty"`

	// A count of the SoftLayer_Hardware objects which are allowed access to this storage volume's Replicant.
	AllowedReplicationVirtualGuestCount *uint `json:"allowedReplicationVirtualGuestCount,omitempty" xmlrpc:"allowedReplicationVirtualGuestCount,omitempty"`

	// The SoftLayer_Hardware objects which are allowed access to this storage volume's Replicant.
	AllowedReplicationVirtualGuests []Virtual_Guest `json:"allowedReplicationVirtualGuests,omitempty" xmlrpc:"allowedReplicationVirtualGuests,omitempty"`

	// A count of the SoftLayer_Network_Subnet objects which are allowed access to this storage volume.
	AllowedSubnetCount *uint `json:"allowedSubnetCount,omitempty" xmlrpc:"allowedSubnetCount,omitempty"`

	// The SoftLayer_Network_Subnet objects which are allowed access to this storage volume.
	AllowedSubnets []Network_Subnet `json:"allowedSubnets,omitempty" xmlrpc:"allowedSubnets,omitempty"`

	// A count of the SoftLayer_Virtual_Guest objects which are allowed access to this storage volume.
	AllowedVirtualGuestCount *uint `json:"allowedVirtualGuestCount,omitempty" xmlrpc:"allowedVirtualGuestCount,omitempty"`

	// The SoftLayer_Virtual_Guest objects which are allowed access to this storage volume.
	AllowedVirtualGuests []Virtual_Guest `json:"allowedVirtualGuests,omitempty" xmlrpc:"allowedVirtualGuests,omitempty"`

	// The current billing item for a Storage volume.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// no documentation yet
	BillingItemCategory *Product_Item_Category `json:"billingItemCategory,omitempty" xmlrpc:"billingItemCategory,omitempty"`

	// The amount of space used by the volume, in bytes.
	BytesUsed *string `json:"bytesUsed,omitempty" xmlrpc:"bytesUsed,omitempty"`

	// A Storage account's capacity, measured in gigabytes.
	CapacityGb *int `json:"capacityGb,omitempty" xmlrpc:"capacityGb,omitempty"`

	// The date a network storage volume was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The schedule id which was executed to create a snapshot.
	CreationScheduleId *string `json:"creationScheduleId,omitempty" xmlrpc:"creationScheduleId,omitempty"`

	// A count of
	CredentialCount *uint `json:"credentialCount,omitempty" xmlrpc:"credentialCount,omitempty"`

	// no documentation yet
	Credentials []Network_Storage_Credential `json:"credentials,omitempty" xmlrpc:"credentials,omitempty"`

	// The Daily Schedule which is associated with this network storage volume.
	DailySchedule *Network_Storage_Schedule `json:"dailySchedule,omitempty" xmlrpc:"dailySchedule,omitempty"`

	// Whether or not a network storage volume is a dependent duplicate.
	DependentDuplicate *string `json:"dependentDuplicate,omitempty" xmlrpc:"dependentDuplicate,omitempty"`

	// A count of the network storage volumes configured to be dependent duplicates of a volume.
	DependentDuplicateCount *uint `json:"dependentDuplicateCount,omitempty" xmlrpc:"dependentDuplicateCount,omitempty"`

	// The network storage volumes configured to be dependent duplicates of a volume.
	DependentDuplicates []Network_Storage `json:"dependentDuplicates,omitempty" xmlrpc:"dependentDuplicates,omitempty"`

	// A count of the events which have taken place on a network storage volume.
	EventCount *uint `json:"eventCount,omitempty" xmlrpc:"eventCount,omitempty"`

	// The events which have taken place on a network storage volume.
	Events []Network_Storage_Event `json:"events,omitempty" xmlrpc:"events,omitempty"`

	// Determines whether the volume is allowed to failback
	FailbackNotAllowed *string `json:"failbackNotAllowed,omitempty" xmlrpc:"failbackNotAllowed,omitempty"`

	// Determines whether the volume is allowed to failover
	FailoverNotAllowed *string `json:"failoverNotAllowed,omitempty" xmlrpc:"failoverNotAllowed,omitempty"`

	// Retrieves the NFS Network Mount Address Name for a given File Storage Volume.
	FileNetworkMountAddress *string `json:"fileNetworkMountAddress,omitempty" xmlrpc:"fileNetworkMountAddress,omitempty"`

	// no documentation yet
	FixReplicationCurrentStatus *string `json:"fixReplicationCurrentStatus,omitempty" xmlrpc:"fixReplicationCurrentStatus,omitempty"`

	// The unique identification number of the guest associated with a Storage volume.
	GuestId *int `json:"guestId,omitempty" xmlrpc:"guestId,omitempty"`

	// When applicable, the hardware associated with a Storage service.
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// The server that is associated with a Storage service.
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// no documentation yet
	HasEncryptionAtRest *bool `json:"hasEncryptionAtRest,omitempty" xmlrpc:"hasEncryptionAtRest,omitempty"`

	// The unique identification number of the host associated with a Storage volume.
	HostId *int `json:"hostId,omitempty" xmlrpc:"hostId,omitempty"`

	// The Hourly Schedule which is associated with this network storage volume.
	HourlySchedule *Network_Storage_Schedule `json:"hourlySchedule,omitempty" xmlrpc:"hourlySchedule,omitempty"`

	// A Storage account's unique identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The Interval Schedule which is associated with this network storage volume.
	IntervalSchedule *Network_Storage_Schedule `json:"intervalSchedule,omitempty" xmlrpc:"intervalSchedule,omitempty"`

	// The maximum number of IOPs selected for this volume.
	Iops *string `json:"iops,omitempty" xmlrpc:"iops,omitempty"`

	// Determines whether network storage volume has an active convert dependent clone to Independent transaction.
	IsConvertToIndependentTransactionInProgress *bool `json:"isConvertToIndependentTransactionInProgress,omitempty" xmlrpc:"isConvertToIndependentTransactionInProgress,omitempty"`

	// Determines whether dependent volume provision is completed on background.
	IsDependentDuplicateProvisionCompleted *bool `json:"isDependentDuplicateProvisionCompleted,omitempty" xmlrpc:"isDependentDuplicateProvisionCompleted,omitempty"`

	// no documentation yet
	IsInDedicatedServiceResource *bool `json:"isInDedicatedServiceResource,omitempty" xmlrpc:"isInDedicatedServiceResource,omitempty"`

	// no documentation yet
	IsMagneticStorage *string `json:"isMagneticStorage,omitempty" xmlrpc:"isMagneticStorage,omitempty"`

	// Determines whether network storage volume has an active provision transaction.
	IsProvisionInProgress *bool `json:"isProvisionInProgress,omitempty" xmlrpc:"isProvisionInProgress,omitempty"`

	// Determines whether a volume is ready to order snapshot space, or, if snapshot space is already available, to assign a snapshot schedule, or to take a manual snapshot.
	IsReadyForSnapshot *bool `json:"isReadyForSnapshot,omitempty" xmlrpc:"isReadyForSnapshot,omitempty"`

	// Determines whether a volume is ready to have Hosts authorized to access it. This does not indicate whether another operation may be blocking, please refer to this volume's volumeStatus property for details.
	IsReadyToMount *bool `json:"isReadyToMount,omitempty" xmlrpc:"isReadyToMount,omitempty"`

	// A count of relationship between a container volume and iSCSI LUNs.
	IscsiLunCount *uint `json:"iscsiLunCount,omitempty" xmlrpc:"iscsiLunCount,omitempty"`

	// Relationship between a container volume and iSCSI LUNs.
	IscsiLuns []Network_Storage `json:"iscsiLuns,omitempty" xmlrpc:"iscsiLuns,omitempty"`

	// The network storage volumes configured to be replicants of this volume.
	IscsiReplicatingVolume *Network_Storage `json:"iscsiReplicatingVolume,omitempty" xmlrpc:"iscsiReplicatingVolume,omitempty"`

	// A count of returns the target IP addresses of an iSCSI volume.
	IscsiTargetIpAddressCount *uint `json:"iscsiTargetIpAddressCount,omitempty" xmlrpc:"iscsiTargetIpAddressCount,omitempty"`

	// Returns the target IP addresses of an iSCSI volume.
	IscsiTargetIpAddresses []string `json:"iscsiTargetIpAddresses,omitempty" xmlrpc:"iscsiTargetIpAddresses,omitempty"`

	// The ID of the LUN volume.
	LunId *string `json:"lunId,omitempty" xmlrpc:"lunId,omitempty"`

	// A count of the manually-created snapshots associated with this SoftLayer_Network_Storage volume. Does not support pagination by result limit and offset.
	ManualSnapshotCount *uint `json:"manualSnapshotCount,omitempty" xmlrpc:"manualSnapshotCount,omitempty"`

	// The manually-created snapshots associated with this SoftLayer_Network_Storage volume. Does not support pagination by result limit and offset.
	ManualSnapshots []Network_Storage `json:"manualSnapshots,omitempty" xmlrpc:"manualSnapshots,omitempty"`

	// A network storage volume's metric tracking object. This object records all periodic polled data available to this volume.
	MetricTrackingObject *Metric_Tracking_Object `json:"metricTrackingObject,omitempty" xmlrpc:"metricTrackingObject,omitempty"`

	// Retrieves the NFS Network Mount Path for a given File Storage Volume.
	MountPath *string `json:"mountPath,omitempty" xmlrpc:"mountPath,omitempty"`

	// Whether or not a network storage volume may be mounted.
	MountableFlag *string `json:"mountableFlag,omitempty" xmlrpc:"mountableFlag,omitempty"`

	// The current status of split or move operation as a part of volume duplication.
	MoveAndSplitStatus *string `json:"moveAndSplitStatus,omitempty" xmlrpc:"moveAndSplitStatus,omitempty"`

	// A Storage account's type. Valid examples are "NAS", "LOCKBOX", "ISCSI", "EVAULT", and "HUB".
	NasType *string `json:"nasType,omitempty" xmlrpc:"nasType,omitempty"`

	// Public notes related to a Storage volume.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// A count of the subscribers that will be notified for usage amount warnings and overages.
	NotificationSubscriberCount *uint `json:"notificationSubscriberCount,omitempty" xmlrpc:"notificationSubscriberCount,omitempty"`

	// The subscribers that will be notified for usage amount warnings and overages.
	NotificationSubscribers []Notification_User_Subscriber `json:"notificationSubscribers,omitempty" xmlrpc:"notificationSubscribers,omitempty"`

	// The name of the snapshot that this volume was duplicated from.
	OriginalSnapshotName *string `json:"originalSnapshotName,omitempty" xmlrpc:"originalSnapshotName,omitempty"`

	// Volume id of the origin volume from which this volume is been cloned.
	OriginalVolumeId *int `json:"originalVolumeId,omitempty" xmlrpc:"originalVolumeId,omitempty"`

	// The name of the volume that this volume was duplicated from.
	OriginalVolumeName *string `json:"originalVolumeName,omitempty" xmlrpc:"originalVolumeName,omitempty"`

	// The size (in GB) of the volume or LUN before any size expansion, or of the volume (before any possible size expansion) from which the duplicate volume or LUN was created.
	OriginalVolumeSize *string `json:"originalVolumeSize,omitempty" xmlrpc:"originalVolumeSize,omitempty"`

	// A volume's configured SoftLayer_Network_Storage_Iscsi_OS_Type.
	OsType *Network_Storage_Iscsi_OS_Type `json:"osType,omitempty" xmlrpc:"osType,omitempty"`

	// A volume's configured SoftLayer_Network_Storage_Iscsi_OS_Type ID.
	OsTypeId *string `json:"osTypeId,omitempty" xmlrpc:"osTypeId,omitempty"`

	// A count of the volumes or snapshots partnered with a network storage volume in a parental role.
	ParentPartnershipCount *uint `json:"parentPartnershipCount,omitempty" xmlrpc:"parentPartnershipCount,omitempty"`

	// The volumes or snapshots partnered with a network storage volume in a parental role.
	ParentPartnerships []Network_Storage_Partnership `json:"parentPartnerships,omitempty" xmlrpc:"parentPartnerships,omitempty"`

	// The parent volume of a volume in a complex storage relationship.
	ParentVolume *Network_Storage `json:"parentVolume,omitempty" xmlrpc:"parentVolume,omitempty"`

	// A count of the volumes or snapshots partnered with a network storage volume.
	PartnershipCount *uint `json:"partnershipCount,omitempty" xmlrpc:"partnershipCount,omitempty"`

	// The volumes or snapshots partnered with a network storage volume.
	Partnerships []Network_Storage_Partnership `json:"partnerships,omitempty" xmlrpc:"partnerships,omitempty"`

	// The password used to access a non-EVault Storage volume. This password is used to register the EVault server agent with the vault backup system.
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// A count of all permissions group(s) this volume is in.
	PermissionsGroupCount *uint `json:"permissionsGroupCount,omitempty" xmlrpc:"permissionsGroupCount,omitempty"`

	// All permissions group(s) this volume is in.
	PermissionsGroups []Network_Storage_Group `json:"permissionsGroups,omitempty" xmlrpc:"permissionsGroups,omitempty"`

	// The properties used to provide additional details about a network storage volume.
	Properties []Network_Storage_Property `json:"properties,omitempty" xmlrpc:"properties,omitempty"`

	// A count of the properties used to provide additional details about a network storage volume.
	PropertyCount *uint `json:"propertyCount,omitempty" xmlrpc:"propertyCount,omitempty"`

	// The number of IOPs provisioned for this volume.
	ProvisionedIops *string `json:"provisionedIops,omitempty" xmlrpc:"provisionedIops,omitempty"`

	// A count of the iSCSI LUN volumes being replicated by this network storage volume.
	ReplicatingLunCount *uint `json:"replicatingLunCount,omitempty" xmlrpc:"replicatingLunCount,omitempty"`

	// The iSCSI LUN volumes being replicated by this network storage volume.
	ReplicatingLuns []Network_Storage `json:"replicatingLuns,omitempty" xmlrpc:"replicatingLuns,omitempty"`

	// The network storage volume being replicated by a volume.
	ReplicatingVolume *Network_Storage `json:"replicatingVolume,omitempty" xmlrpc:"replicatingVolume,omitempty"`

	// A count of the volume replication events.
	ReplicationEventCount *uint `json:"replicationEventCount,omitempty" xmlrpc:"replicationEventCount,omitempty"`

	// The volume replication events.
	ReplicationEvents []Network_Storage_Event `json:"replicationEvents,omitempty" xmlrpc:"replicationEvents,omitempty"`

	// A count of the network storage volumes configured to be replicants of a volume.
	ReplicationPartnerCount *uint `json:"replicationPartnerCount,omitempty" xmlrpc:"replicationPartnerCount,omitempty"`

	// The network storage volumes configured to be replicants of a volume.
	ReplicationPartners []Network_Storage `json:"replicationPartners,omitempty" xmlrpc:"replicationPartners,omitempty"`

	// The Replication Schedule associated with a network storage volume.
	ReplicationSchedule *Network_Storage_Schedule `json:"replicationSchedule,omitempty" xmlrpc:"replicationSchedule,omitempty"`

	// The current replication status of a network storage volume. Indicates Failover or Failback status.
	ReplicationStatus *string `json:"replicationStatus,omitempty" xmlrpc:"replicationStatus,omitempty"`

	// A count of the schedules which are associated with a network storage volume.
	ScheduleCount *uint `json:"scheduleCount,omitempty" xmlrpc:"scheduleCount,omitempty"`

	// The schedules which are associated with a network storage volume.
	Schedules []Network_Storage_Schedule `json:"schedules,omitempty" xmlrpc:"schedules,omitempty"`

	// Service Provider ID
	ServiceProviderId *int `json:"serviceProviderId,omitempty" xmlrpc:"serviceProviderId,omitempty"`

	// The network resource a Storage service is connected to.
	ServiceResource *Network_Service_Resource `json:"serviceResource,omitempty" xmlrpc:"serviceResource,omitempty"`

	// The IP address of a Storage resource.
	ServiceResourceBackendIpAddress *string `json:"serviceResourceBackendIpAddress,omitempty" xmlrpc:"serviceResourceBackendIpAddress,omitempty"`

	// The name of a Storage's network resource.
	ServiceResourceName *string `json:"serviceResourceName,omitempty" xmlrpc:"serviceResourceName,omitempty"`

	// A volume's configured snapshot space size.
	SnapshotCapacityGb *string `json:"snapshotCapacityGb,omitempty" xmlrpc:"snapshotCapacityGb,omitempty"`

	// A count of the snapshots associated with this SoftLayer_Network_Storage volume.
	SnapshotCount *uint `json:"snapshotCount,omitempty" xmlrpc:"snapshotCount,omitempty"`

	// The creation timestamp of the snapshot on the storage platform.
	SnapshotCreationTimestamp *string `json:"snapshotCreationTimestamp,omitempty" xmlrpc:"snapshotCreationTimestamp,omitempty"`

	// The percentage of used snapshot space after which to delete automated snapshots.
	SnapshotDeletionThresholdPercentage *string `json:"snapshotDeletionThresholdPercentage,omitempty" xmlrpc:"snapshotDeletionThresholdPercentage,omitempty"`

	// Whether or not a network storage volume may be mounted.
	SnapshotNotificationStatus *string `json:"snapshotNotificationStatus,omitempty" xmlrpc:"snapshotNotificationStatus,omitempty"`

	// The snapshot size in bytes.
	SnapshotSizeBytes *string `json:"snapshotSizeBytes,omitempty" xmlrpc:"snapshotSizeBytes,omitempty"`

	// A volume's available snapshot reservation space.
	SnapshotSpaceAvailable *string `json:"snapshotSpaceAvailable,omitempty" xmlrpc:"snapshotSpaceAvailable,omitempty"`

	// The snapshots associated with this SoftLayer_Network_Storage volume.
	Snapshots []Network_Storage `json:"snapshots,omitempty" xmlrpc:"snapshots,omitempty"`

	// no documentation yet
	StaasVersion *string `json:"staasVersion,omitempty" xmlrpc:"staasVersion,omitempty"`

	// A count of the network storage groups this volume is attached to.
	StorageGroupCount *uint `json:"storageGroupCount,omitempty" xmlrpc:"storageGroupCount,omitempty"`

	// The network storage groups this volume is attached to.
	StorageGroups []Network_Storage_Group `json:"storageGroups,omitempty" xmlrpc:"storageGroups,omitempty"`

	// no documentation yet
	StorageTierLevel *string `json:"storageTierLevel,omitempty" xmlrpc:"storageTierLevel,omitempty"`

	// A description of the Storage object.
	StorageType *Network_Storage_Type `json:"storageType,omitempty" xmlrpc:"storageType,omitempty"`

	// A storage object's type.
	StorageTypeId *string `json:"storageTypeId,omitempty" xmlrpc:"storageTypeId,omitempty"`

	// The amount of space used by the volume.
	TotalBytesUsed *string `json:"totalBytesUsed,omitempty" xmlrpc:"totalBytesUsed,omitempty"`

	// The total snapshot retention count of all schedules on this network storage volume.
	TotalScheduleSnapshotRetentionCount *uint `json:"totalScheduleSnapshotRetentionCount,omitempty" xmlrpc:"totalScheduleSnapshotRetentionCount,omitempty"`

	// This flag indicates whether this storage type is upgradable or not.
	UpgradableFlag *bool `json:"upgradableFlag,omitempty" xmlrpc:"upgradableFlag,omitempty"`

	// The usage notification for SL Storage services.
	UsageNotification *Notification `json:"usageNotification,omitempty" xmlrpc:"usageNotification,omitempty"`

	// The username used to access a non-EVault Storage volume. This username is used to register the EVault server agent with the vault backup system.
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`

	// The type of network storage service.
	VendorName *string `json:"vendorName,omitempty" xmlrpc:"vendorName,omitempty"`

	// When applicable, the virtual guest associated with a Storage service.
	VirtualGuest *Virtual_Guest `json:"virtualGuest,omitempty" xmlrpc:"virtualGuest,omitempty"`

	// The username and password history for a Storage service.
	VolumeHistory []Network_Storage_History `json:"volumeHistory,omitempty" xmlrpc:"volumeHistory,omitempty"`

	// A count of the username and password history for a Storage service.
	VolumeHistoryCount *uint `json:"volumeHistoryCount,omitempty" xmlrpc:"volumeHistoryCount,omitempty"`

	// The current status of a network storage volume.
	VolumeStatus *string `json:"volumeStatus,omitempty" xmlrpc:"volumeStatus,omitempty"`

	// The account username and password for the EVault webCC interface.
	WebccAccount *Account_Password `json:"webccAccount,omitempty" xmlrpc:"webccAccount,omitempty"`

	// The Weekly Schedule which is associated with this network storage volume.
	WeeklySchedule *Network_Storage_Schedule `json:"weeklySchedule,omitempty" xmlrpc:"weeklySchedule,omitempty"`
}

// no documentation yet
type Network_Storage_Allowed_Host struct {
	Entity

	// The account to which this allowed host belongs to.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// A count of the SoftLayer_Network_Storage_Group objects this SoftLayer_Network_Storage_Allowed_Host is present in.
	AssignedGroupCount *uint `json:"assignedGroupCount,omitempty" xmlrpc:"assignedGroupCount,omitempty"`

	// The SoftLayer_Network_Storage_Group objects this SoftLayer_Network_Storage_Allowed_Host is present in.
	AssignedGroups []Network_Storage_Group `json:"assignedGroups,omitempty" xmlrpc:"assignedGroups,omitempty"`

	// A count of the SoftLayer_Network_Storage volumes to which this SoftLayer_Network_Storage_Allowed_Host is allowed access.
	AssignedIscsiVolumeCount *uint `json:"assignedIscsiVolumeCount,omitempty" xmlrpc:"assignedIscsiVolumeCount,omitempty"`

	// The SoftLayer_Network_Storage volumes to which this SoftLayer_Network_Storage_Allowed_Host is allowed access.
	AssignedIscsiVolumes []Network_Storage `json:"assignedIscsiVolumes,omitempty" xmlrpc:"assignedIscsiVolumes,omitempty"`

	// A count of the SoftLayer_Network_Storage volumes to which this SoftLayer_Network_Storage_Allowed_Host is allowed access.
	AssignedNfsVolumeCount *uint `json:"assignedNfsVolumeCount,omitempty" xmlrpc:"assignedNfsVolumeCount,omitempty"`

	// The SoftLayer_Network_Storage volumes to which this SoftLayer_Network_Storage_Allowed_Host is allowed access.
	AssignedNfsVolumes []Network_Storage `json:"assignedNfsVolumes,omitempty" xmlrpc:"assignedNfsVolumes,omitempty"`

	// A count of the SoftLayer_Network_Storage primary volumes whose replicas are allowed access.
	AssignedReplicationVolumeCount *uint `json:"assignedReplicationVolumeCount,omitempty" xmlrpc:"assignedReplicationVolumeCount,omitempty"`

	// The SoftLayer_Network_Storage primary volumes whose replicas are allowed access.
	AssignedReplicationVolumes []Network_Storage `json:"assignedReplicationVolumes,omitempty" xmlrpc:"assignedReplicationVolumes,omitempty"`

	// A count of the SoftLayer_Network_Storage volumes to which this SoftLayer_Network_Storage_Allowed_Host is allowed access.
	AssignedVolumeCount *uint `json:"assignedVolumeCount,omitempty" xmlrpc:"assignedVolumeCount,omitempty"`

	// The SoftLayer_Network_Storage volumes to which this SoftLayer_Network_Storage_Allowed_Host is allowed access.
	AssignedVolumes []Network_Storage `json:"assignedVolumes,omitempty" xmlrpc:"assignedVolumes,omitempty"`

	// The SoftLayer_Network_Storage_Credential this allowed host uses.
	Credential *Network_Storage_Credential `json:"credential,omitempty" xmlrpc:"credential,omitempty"`

	// The credential this allowed host will use
	CredentialId *int `json:"credentialId,omitempty" xmlrpc:"credentialId,omitempty"`

	// The internal identifier of the igroup
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The name of allowed host, usually an IQN or other identifier
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	ResourceTableId *int `json:"resourceTableId,omitempty" xmlrpc:"resourceTableId,omitempty"`

	// no documentation yet
	ResourceTableName *string `json:"resourceTableName,omitempty" xmlrpc:"resourceTableName,omitempty"`

	// Connections to a target with a source IP in this subnet prefix are allowed.
	SourceSubnet *string `json:"sourceSubnet,omitempty" xmlrpc:"sourceSubnet,omitempty"`

	// The SoftLayer_Network_Subnet records assigned to the ACL for this allowed host.
	SubnetsInAcl []Network_Subnet `json:"subnetsInAcl,omitempty" xmlrpc:"subnetsInAcl,omitempty"`

	// A count of the SoftLayer_Network_Subnet records assigned to the ACL for this allowed host.
	SubnetsInAclCount *uint `json:"subnetsInAclCount,omitempty" xmlrpc:"subnetsInAclCount,omitempty"`
}

// no documentation yet
type Network_Storage_Allowed_Host_Hardware struct {
	Network_Storage_Allowed_Host

	// The SoftLayer_Account object which this SoftLayer_Network_Storage_Allowed_Host belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The SoftLayer_Hardware object which this SoftLayer_Network_Storage_Allowed_Host is referencing.
	Resource *Hardware `json:"resource,omitempty" xmlrpc:"resource,omitempty"`
}

// no documentation yet
type Network_Storage_Allowed_Host_IpAddress struct {
	Network_Storage_Allowed_Host

	// The SoftLayer_Account object which this SoftLayer_Network_Storage_Allowed_Host belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The SoftLayer_Network_Subnet_IpAddress object which this SoftLayer_Network_Storage_Allowed_Host is referencing.
	Resource *Network_Subnet_IpAddress `json:"resource,omitempty" xmlrpc:"resource,omitempty"`
}

// no documentation yet
type Network_Storage_Allowed_Host_Subnet struct {
	Network_Storage_Allowed_Host

	// The SoftLayer_Account object which this SoftLayer_Network_Storage_Allowed_Host belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The SoftLayer_Network_Subnet object which this SoftLayer_Network_Storage_Allowed_Host is referencing.
	Resource *Network_Subnet `json:"resource,omitempty" xmlrpc:"resource,omitempty"`
}

// no documentation yet
type Network_Storage_Allowed_Host_VirtualGuest struct {
	Network_Storage_Allowed_Host

	// The SoftLayer_Account object which this SoftLayer_Network_Storage_Allowed_Host belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The SoftLayer_Virtual_Guest object which this SoftLayer_Network_Storage_Allowed_Host is referencing.
	Resource *Virtual_Guest `json:"resource,omitempty" xmlrpc:"resource,omitempty"`
}

// The SoftLayer_Network_Storage_Backup contains general information regarding a Storage backup service such as account id, username, maximum capacity, password, Storage's product type and the server id.
type Network_Storage_Backup struct {
	Network_Storage

	// Peak number of bytes used in the vault for the current billing cycle.
	CurrentCyclePeakUsage *uint `json:"currentCyclePeakUsage,omitempty" xmlrpc:"currentCyclePeakUsage,omitempty"`

	// Peak number of bytes used in the vault for the previous billing cycle.
	PreviousCyclePeakUsage *uint `json:"previousCyclePeakUsage,omitempty" xmlrpc:"previousCyclePeakUsage,omitempty"`
}

// The SoftLayer_Network_Storage_Backup_Evault contains general information regarding an EVault Storage service such as account id, username, maximum capacity, password, Storage's product type and the server id.
type Network_Storage_Backup_Evault struct {
	Network_Storage_Backup
}

// The SoftLayer_Network_Storage_Backup_Evault_Version6 contains the same properties as the SoftLayer_Network_Storage_Backup_Evault. Additional properties available for the EVault Storage type:  softwareComponent, totalBytesUsed, backupJobDetails, restoreJobDetails and agentStatuses
type Network_Storage_Backup_Evault_Version6 struct {
	Network_Storage_Backup_Evault

	// A count of statuses (most of the time will be one status) for the agent tied to the EVault Storage services.
	AgentStatusCount *uint `json:"agentStatusCount,omitempty" xmlrpc:"agentStatusCount,omitempty"`

	// Statuses (most of the time will be one status) for the agent tied to the EVault Storage services.
	AgentStatuses []Container_Network_Storage_Evault_WebCc_AgentStatus `json:"agentStatuses,omitempty" xmlrpc:"agentStatuses,omitempty"`

	// A count of all the of the backup jobs for the EVault Storage account.
	BackupJobDetailCount *uint `json:"backupJobDetailCount,omitempty" xmlrpc:"backupJobDetailCount,omitempty"`

	// All the of the backup jobs for the EVault Storage account.
	BackupJobDetails []Container_Network_Storage_Evault_WebCc_JobDetails `json:"backupJobDetails,omitempty" xmlrpc:"backupJobDetails,omitempty"`

	// A count of the billing items for plugins tied to the EVault Storage service.
	PluginBillingItemCount *uint `json:"pluginBillingItemCount,omitempty" xmlrpc:"pluginBillingItemCount,omitempty"`

	// The billing items for plugins tied to the EVault Storage service.
	PluginBillingItems []Billing_Item `json:"pluginBillingItems,omitempty" xmlrpc:"pluginBillingItems,omitempty"`

	// A count of all the of the restore jobs for the EVault Storage account.
	RestoreJobDetailCount *uint `json:"restoreJobDetailCount,omitempty" xmlrpc:"restoreJobDetailCount,omitempty"`

	// All the of the restore jobs for the EVault Storage account.
	RestoreJobDetails []Container_Network_Storage_Evault_WebCc_JobDetails `json:"restoreJobDetails,omitempty" xmlrpc:"restoreJobDetails,omitempty"`

	// The software component for the EVault base client.
	SoftwareComponent *Software_Component `json:"softwareComponent,omitempty" xmlrpc:"softwareComponent,omitempty"`

	// A count of retrieve the task information for the EVault Storage service.
	TaskCount *uint `json:"taskCount,omitempty" xmlrpc:"taskCount,omitempty"`

	// Retrieve the task information for the EVault Storage service.
	Tasks []Container_Network_Storage_Evault_Vault_Task `json:"tasks,omitempty" xmlrpc:"tasks,omitempty"`
}

// The SoftLayer_Network_Storage_Credential data type will give you an overview of the usernames that are currently attached to your storage device.
type Network_Storage_Credential struct {
	Entity

	// This is the account that the storage credential is tied to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// This is the account id associated with the volume.
	AccountId *string `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// This is the data that the record was created in the table.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// This is the date that the record was last updated in the table.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// This is the id of the type of credential that this object represents.
	NasCredentialTypeId *int `json:"nasCredentialTypeId,omitempty" xmlrpc:"nasCredentialTypeId,omitempty"`

	// These are the SoftLayer_Network_Storage_Allowed_Host entries that this credential is assigned to.
	NetworkStorageAllowedHosts *Network_Storage_Allowed_Host `json:"networkStorageAllowedHosts,omitempty" xmlrpc:"networkStorageAllowedHosts,omitempty"`

	// This is the password associated with the volume.
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// These are the types of storage that the credential can be assigned to.
	Type *Network_Storage_Credential_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// This is the username associated with the volume.
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`

	// A count of these are the SoftLayer_Network_Storage volumes that this credential is assigned to.
	VolumeCount *uint `json:"volumeCount,omitempty" xmlrpc:"volumeCount,omitempty"`

	// These are the SoftLayer_Network_Storage volumes that this credential is assigned to.
	Volumes []Network_Storage `json:"volumes,omitempty" xmlrpc:"volumes,omitempty"`
}

// <<<
type Network_Storage_Credential_Type struct {
	Entity

	// The date a credential type was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// A short description of the credential type
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The key name of the credential type.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The date a credential was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The human readable name of the credential type.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Storage_Daily_Usage struct {
	Entity

	// no documentation yet
	BytesUsed *uint `json:"bytesUsed,omitempty" xmlrpc:"bytesUsed,omitempty"`

	// no documentation yet
	CdnHttpBandwidth *uint `json:"cdnHttpBandwidth,omitempty" xmlrpc:"cdnHttpBandwidth,omitempty"`

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	NasVolume *Network_Storage `json:"nasVolume,omitempty" xmlrpc:"nasVolume,omitempty"`

	// no documentation yet
	NasVolumeId *int `json:"nasVolumeId,omitempty" xmlrpc:"nasVolumeId,omitempty"`

	// no documentation yet
	PublicBandwidthOut *uint `json:"publicBandwidthOut,omitempty" xmlrpc:"publicBandwidthOut,omitempty"`
}

// no documentation yet
type Network_Storage_DedicatedCluster struct {
	Entity

	// no documentation yet
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The SoftLayer_Account->id of the customer account
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The date when Dedicated service resource entry was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The unique identifier for Dedicated service resource record.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	ServiceResource *Network_Service_Resource `json:"serviceResource,omitempty" xmlrpc:"serviceResource,omitempty"`

	// The cluster Id that is setup as dedicated for the customer.
	ServiceResourceId *int `json:"serviceResourceId,omitempty" xmlrpc:"serviceResourceId,omitempty"`
}

// Storage volumes can create various events to keep track of what has occurred to the volume. Events provide an audit trail that can be used to verify that various tasks have occurred, such as snapshots to be created by a schedule or remote replication synchronization.
type Network_Storage_Event struct {
	Entity

	// The date an event was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The message text for an event.
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// A schedule that is associated with an event. Not all events will have a schedule.
	Schedule *Network_Storage_Schedule `json:"schedule,omitempty" xmlrpc:"schedule,omitempty"`

	// An identifier for the schedule which is associated with an event.
	ScheduleId *int `json:"scheduleId,omitempty" xmlrpc:"scheduleId,omitempty"`

	// A Storage volume's event type. The type provides a standardized definition for an event.
	Type *Network_Storage_Event_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// An identifier for the type of an event.
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`

	// The associated volume for an event.
	Volume *Network_Storage `json:"volume,omitempty" xmlrpc:"volume,omitempty"`

	// The volume id which an event is associated with.
	VolumeId *int `json:"volumeId,omitempty" xmlrpc:"volumeId,omitempty"`
}

// no documentation yet
type Network_Storage_Event_Type struct {
	Entity

	// no documentation yet
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Storage_Group struct {
	Entity

	// The SoftLayer_Account which owns this group.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The account ID which owns this group
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The friendly name of this group
	Alias *string `json:"alias,omitempty" xmlrpc:"alias,omitempty"`

	// A count of the allowed hosts list for this group.
	AllowedHostCount *uint `json:"allowedHostCount,omitempty" xmlrpc:"allowedHostCount,omitempty"`

	// The allowed hosts list for this group.
	AllowedHosts []Network_Storage_Allowed_Host `json:"allowedHosts,omitempty" xmlrpc:"allowedHosts,omitempty"`

	// A count of the network storage volumes this group is attached to.
	AttachedVolumeCount *uint `json:"attachedVolumeCount,omitempty" xmlrpc:"attachedVolumeCount,omitempty"`

	// The network storage volumes this group is attached to.
	AttachedVolumes []Network_Storage `json:"attachedVolumes,omitempty" xmlrpc:"attachedVolumes,omitempty"`

	// The date this group was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The type which defines this group.
	GroupType *Network_Storage_Group_Type `json:"groupType,omitempty" xmlrpc:"groupType,omitempty"`

	// The SoftLayer_Network_Storage_Group_Type which describes this group.
	GroupTypeId *int `json:"groupTypeId,omitempty" xmlrpc:"groupTypeId,omitempty"`

	// The internal identifier of the group
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The OS Type this group is configured for.
	OsType *Network_Storage_Iscsi_OS_Type `json:"osType,omitempty" xmlrpc:"osType,omitempty"`

	// A SoftLayer_Network_Storage_OS_Type Operating System designation that this group was created for.
	OsTypeId *int `json:"osTypeId,omitempty" xmlrpc:"osTypeId,omitempty"`

	// The network resource this group is created on.
	ServiceResource *Network_Service_Resource `json:"serviceResource,omitempty" xmlrpc:"serviceResource,omitempty"`

	// A SoftLayer_Network_Service_Resource that this group was created on.
	ServiceResourceId *int `json:"serviceResourceId,omitempty" xmlrpc:"serviceResourceId,omitempty"`
}

// no documentation yet
type Network_Storage_Group_Iscsi struct {
	Network_Storage_Group
}

// no documentation yet
type Network_Storage_Group_Nfs struct {
	Network_Storage_Group
}

// no documentation yet
type Network_Storage_Group_Type struct {
	Entity

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The SoftLayer_Network_Storage_History contains the username/password past history for Storage services except Evault. Information such as the username, passwords, notes and the date of the password change may be retrieved.
type Network_Storage_History struct {
	Entity

	// The account that the Storage services belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// Date the password was changed.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The Storage service that the password history belongs to.
	NasVolume *Network_Storage `json:"nasVolume,omitempty" xmlrpc:"nasVolume,omitempty"`

	// Past notes for the Storage service.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// Password for the Storage service that was used in the past.
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// Username for the Storage service.
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`
}

// The SoftLayer_Network_Storage_Hub data type models Virtual Server type Storage storage offerings.
type Network_Storage_Hub struct {
	Network_Storage

	// A count of the billing items tied to a Storage service's bandwidth usage.
	BandwidthBillingItemCount *uint `json:"bandwidthBillingItemCount,omitempty" xmlrpc:"bandwidthBillingItemCount,omitempty"`

	// The billing items tied to a Storage service's bandwidth usage.
	BandwidthBillingItems []Billing_Item `json:"bandwidthBillingItems,omitempty" xmlrpc:"bandwidthBillingItems,omitempty"`
}

// no documentation yet
type Network_Storage_Hub_Cleversafe_Account struct {
	Entity

	// SoftLayer account to which an IBM Cloud Object Storage account belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The ID of the SoftLayer_Account which this IBM Cloud Object Storage account is associated with.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// An associated parent billing item which is active. Includes billing items which are scheduled to be cancelled in the future.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// An associated parent billing item which has been cancelled.
	CancelledBillingItem *Billing_Item `json:"cancelledBillingItem,omitempty" xmlrpc:"cancelledBillingItem,omitempty"`

	// A count of credentials used for generating an AWS signature. Max of 2.
	CredentialCount *uint `json:"credentialCount,omitempty" xmlrpc:"credentialCount,omitempty"`

	// Credentials used for generating an AWS signature. Max of 2.
	Credentials []Network_Storage_Credential `json:"credentials,omitempty" xmlrpc:"credentials,omitempty"`

	// The IMS ID of an IBM Cloud Object Storage account.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Provides an interface to various metrics relating to the usage of an IBM Cloud Object Storage account.
	MetricTrackingObject *Metric_Tracking_Object `json:"metricTrackingObject,omitempty" xmlrpc:"metricTrackingObject,omitempty"`

	// A user-defined field of notes.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// Human readable identifier of IBM Cloud Object Storage accounts.
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`

	// Unique identifier for an IBM Cloud Object Storage account.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// no documentation yet
type Network_Storage_Hub_Swift struct {
	Network_Storage_Hub

	// A count of
	StorageNodeCount *uint `json:"storageNodeCount,omitempty" xmlrpc:"storageNodeCount,omitempty"`

	// no documentation yet
	StorageNodes []Network_Service_Resource `json:"storageNodes,omitempty" xmlrpc:"storageNodes,omitempty"`
}

// no documentation yet
type Network_Storage_Hub_Swift_Container struct {
	Network_Storage_Hub_Swift
}

// no documentation yet
type Network_Storage_Hub_Swift_Metrics struct {
	Entity
}

// no documentation yet
type Network_Storage_Hub_Swift_Share struct {
	Entity
}

// no documentation yet
type Network_Storage_Hub_Swift_Version1 struct {
	Network_Storage_Hub_Swift
}

// The iscsi data type provides access to additional information about an iscsi volume such as the snapshot capacity limit and replication partners.
type Network_Storage_Iscsi struct {
	Network_Storage
}

// no documentation yet
type Network_Storage_Iscsi_OS_Type struct {
	Entity

	// The date this OS type record was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The description of this OS type
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The internal identifier of the OS type selection
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The key name of this OS type
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The name of this OS type
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Storage_MassDataMigration_CrossRegion_Country_Xref struct {
	Entity

	// SoftLayer_Locale_Country Id.
	Country *Locale_Country `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// no documentation yet
	CountryId *int `json:"countryId,omitempty" xmlrpc:"countryId,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Location Group ID of CleverSafe cross region.
	LocationGroup *Location_Group `json:"locationGroup,omitempty" xmlrpc:"locationGroup,omitempty"`

	// no documentation yet
	LocationGroupId *int `json:"locationGroupId,omitempty" xmlrpc:"locationGroupId,omitempty"`
}

// The SoftLayer_Network_Storage_MassDataMigration_Request data type contains information on a single Mass Data Migration request. Creation of these requests is limited to SoftLayer customers through the SoftLayer Customer Portal.
type Network_Storage_MassDataMigration_Request struct {
	Entity

	// The account to which the request belongs.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The account id of the request.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// A count of the active tickets that are attached to the MDMS request.
	ActiveTicketCount *uint `json:"activeTicketCount,omitempty" xmlrpc:"activeTicketCount,omitempty"`

	// The active tickets that are attached to the MDMS request.
	ActiveTickets []Ticket `json:"activeTickets,omitempty" xmlrpc:"activeTickets,omitempty"`

	// The customer address where the device is shipped to.
	Address *Account_Address `json:"address,omitempty" xmlrpc:"address,omitempty"`

	// The address id of address assigned to this request.
	AddressId *int `json:"addressId,omitempty" xmlrpc:"addressId,omitempty"`

	// An associated parent billing item which is active. Includes billing items which are scheduled to be cancelled in the future.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// The employee user who created the request.
	CreateEmployee *User_Employee `json:"createEmployee,omitempty" xmlrpc:"createEmployee,omitempty"`

	// The customer user who created the request.
	CreateUser *User_Customer `json:"createUser,omitempty" xmlrpc:"createUser,omitempty"`

	// The create user id of the request.
	CreateUserId *int `json:"createUserId,omitempty" xmlrpc:"createUserId,omitempty"`

	// The device configurations.
	DeviceConfiguration *Network_Storage_MassDataMigration_Request_DeviceConfiguration `json:"deviceConfiguration,omitempty" xmlrpc:"deviceConfiguration,omitempty"`

	// The model of device assigned to this request.
	DeviceModel *string `json:"deviceModel,omitempty" xmlrpc:"deviceModel,omitempty"`

	// The end date of the request.
	EndDate *Time `json:"endDate,omitempty" xmlrpc:"endDate,omitempty"`

	// The unique id of the request.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of the key contacts for this requests.
	KeyContactCount *uint `json:"keyContactCount,omitempty" xmlrpc:"keyContactCount,omitempty"`

	// The key contacts for this requests.
	KeyContacts []Network_Storage_MassDataMigration_Request_KeyContact `json:"keyContacts,omitempty" xmlrpc:"keyContacts,omitempty"`

	// The employee who last modified the request.
	ModifyEmployee *User_Employee `json:"modifyEmployee,omitempty" xmlrpc:"modifyEmployee,omitempty"`

	// The customer user who last modified the request.
	ModifyUser *User_Customer `json:"modifyUser,omitempty" xmlrpc:"modifyUser,omitempty"`

	// The modify user id of the request.
	ModifyUserId *int `json:"modifyUserId,omitempty" xmlrpc:"modifyUserId,omitempty"`

	// The unique id of the request.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A count of the shipments of the request.
	ShipmentCount *uint `json:"shipmentCount,omitempty" xmlrpc:"shipmentCount,omitempty"`

	// The shipments of the request.
	Shipments []Account_Shipment `json:"shipments,omitempty" xmlrpc:"shipments,omitempty"`

	// The start date of the request.
	StartDate *Time `json:"startDate,omitempty" xmlrpc:"startDate,omitempty"`

	// The status of the request.
	Status *Network_Storage_MassDataMigration_Request_Status `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// The status id of the request.
	StatusId *int `json:"statusId,omitempty" xmlrpc:"statusId,omitempty"`

	// Ticket that is attached to this mass data migration request.
	Ticket *Ticket `json:"ticket,omitempty" xmlrpc:"ticket,omitempty"`

	// A count of all tickets that are attached to the mass data migration request.
	TicketCount *uint `json:"ticketCount,omitempty" xmlrpc:"ticketCount,omitempty"`

	// All tickets that are attached to the mass data migration request.
	Tickets []Ticket `json:"tickets,omitempty" xmlrpc:"tickets,omitempty"`
}

// The SoftLayer_Network_Storage_MassDataMigration_Request_DeviceConfiguration data type contains settings such networking, COS account, which needs to be configured on device for a Mass Data Migration Request.
type Network_Storage_MassDataMigration_Request_DeviceConfiguration struct {
	Entity

	// The account id.
	CosAccountId *int `json:"cosAccountId,omitempty" xmlrpc:"cosAccountId,omitempty"`

	// The Cloud Object Storage bucket.
	CosBucket *string `json:"cosBucket,omitempty" xmlrpc:"cosBucket,omitempty"`

	// The eth1 gateway for connecting to private network in datacenter.
	Eth1Gateway *string `json:"eth1Gateway,omitempty" xmlrpc:"eth1Gateway,omitempty"`

	// The eth1 IP address for connecting to private network in datacenter.
	Eth1IpAddress *string `json:"eth1IpAddress,omitempty" xmlrpc:"eth1IpAddress,omitempty"`

	// The eth1 netmask for connecting to private network in datacenter.
	Eth1Netmask *string `json:"eth1Netmask,omitempty" xmlrpc:"eth1Netmask,omitempty"`

	// The eth3 gateway for connecting to private network at customer's location.
	Eth3Gateway *string `json:"eth3Gateway,omitempty" xmlrpc:"eth3Gateway,omitempty"`

	// The eth3 IP address for connecting to private network at customer location.
	Eth3IpAddress *string `json:"eth3IpAddress,omitempty" xmlrpc:"eth3IpAddress,omitempty"`

	// The eth3 netmask for connecting to private network in at customer's location.
	Eth3Netmask *string `json:"eth3Netmask,omitempty" xmlrpc:"eth3Netmask,omitempty"`

	// The unique id of the request status.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The password for configuring network share.
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// The pool lock password for configuring network share.
	PoolLockPassword *string `json:"poolLockPassword,omitempty" xmlrpc:"poolLockPassword,omitempty"`

	// The request this device configurations belongs to.
	Request *Network_Storage_MassDataMigration_Request `json:"request,omitempty" xmlrpc:"request,omitempty"`

	// The request id.
	RequestId *int `json:"requestId,omitempty" xmlrpc:"requestId,omitempty"`

	// The Cloud Object Storage bucket URL.
	S3Url *string `json:"s3Url,omitempty" xmlrpc:"s3Url,omitempty"`

	// The name of network share.
	ShareName *string `json:"shareName,omitempty" xmlrpc:"shareName,omitempty"`

	// The storage account to use for this request.
	StorageAccount *Network_Storage_Hub_Cleversafe_Account `json:"storageAccount,omitempty" xmlrpc:"storageAccount,omitempty"`

	// The username for configuring network share.
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`
}

// The SoftLayer_Network_Storage_MassDataMigration_Request_KeyContact data type contains name, email, and phone for key contact at customer location who will handle Mass Data Migration.
type Network_Storage_MassDataMigration_Request_KeyContact struct {
	Entity

	// The request this key contact belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// An account number that is linked to a KeyContact.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The date a KeyContact was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// KeyContact's Email Id.
	Email *string `json:"email,omitempty" xmlrpc:"email,omitempty"`

	// The unique id of the key contact.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The date a KeyContact was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// KeyContact's Name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A phone number assigned to a KeyContact.
	Phone *string `json:"phone,omitempty" xmlrpc:"phone,omitempty"`

	// The request this key contact belongs to.
	Request *Network_Storage_MassDataMigration_Request `json:"request,omitempty" xmlrpc:"request,omitempty"`

	// A request id that is linked to a KeyContact.
	RequestId *int `json:"requestId,omitempty" xmlrpc:"requestId,omitempty"`
}

// The SoftLayer_Network_Storage_MassDataMigration_Request_Status data type contains general information relating to the statuses to which a Mass Data Migration Request may be set.
type Network_Storage_MassDataMigration_Request_Status struct {
	Entity

	// The description of the request status.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The unique id of the request status.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The unique keyname of the request status.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The name of the request status.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The SoftLayer_Network_Storage_Nas contains general information regarding a NAS Storage service such as account id, username, password, maximum capacity, Storage's product type and capacity.
type Network_Storage_Nas struct {
	Network_Storage

	// no documentation yet
	RecentBytesUsed *Network_Storage_Daily_Usage `json:"recentBytesUsed,omitempty" xmlrpc:"recentBytesUsed,omitempty"`
}

// A network storage partnership is used to link multiple volumes to each other. These partnerships describe replication hierarchies or link volume snapshots to their associated storage volume.
type Network_Storage_Partnership struct {
	Entity

	// The date a partnership was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The date a partnership was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The associated child volume for a partnership.
	PartnerVolume *Network_Storage `json:"partnerVolume,omitempty" xmlrpc:"partnerVolume,omitempty"`

	// The child volume id which a partnership is associated with.
	PartnerVolumeId *int `json:"partnerVolumeId,omitempty" xmlrpc:"partnerVolumeId,omitempty"`

	// The type provides a standardized definition for a partnership.
	Type *Network_Storage_Partnership_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The associated parent volume for a partnership.
	Volume *Network_Storage `json:"volume,omitempty" xmlrpc:"volume,omitempty"`

	// The volume id which a partnership is associated with.
	VolumeId *int `json:"volumeId,omitempty" xmlrpc:"volumeId,omitempty"`
}

// A network storage partnership type is used to define the link between two volumes.
type Network_Storage_Partnership_Type struct {
	Entity

	// A type's description, for example 'ISCSI snapshot partnership'.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// A type's key name, for example 'ISCSI_SNAPSHOT'.
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// A type's name, for example 'ISCSI Snapshot'.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// A property provides additional information about a volume which it is assigned to. This information can range from "Mountable" flags to utilized snapshot space.
type Network_Storage_Property struct {
	Entity

	// The date a property was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The date a property was last modified;
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The type provides a standardized definition for a property.
	Type *Network_Storage_Property_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The value of a property.
	Value *string `json:"value,omitempty" xmlrpc:"value,omitempty"`

	// The associated volume for a property.
	Volume *Network_Storage `json:"volume,omitempty" xmlrpc:"volume,omitempty"`

	// The volume id which a property is associated with.
	VolumeId *int `json:"volumeId,omitempty" xmlrpc:"volumeId,omitempty"`
}

// The storage property types provide standard definitions for properties which can be used with any type for Storage offering.  The properties provide additional information about a volume which they are assigned to.
type Network_Storage_Property_Type struct {
	Entity

	// A type's description, for example 'Determines whether the volume is currently mountable'.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// A type's keyname, for example 'MOUNTABLE'.
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// A type's name, for example 'Mountable'.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Storage_Replicant struct {
	Network_Storage

	// When a replicant is in the process of synchronizing with the parent volume this flag will be true.
	FailbackInProgressFlag *string `json:"failbackInProgressFlag,omitempty" xmlrpc:"failbackInProgressFlag,omitempty"`

	// The volume name for a replicant.
	VolumeName *string `json:"volumeName,omitempty" xmlrpc:"volumeName,omitempty"`
}

// Schedules can be created for select Storage services, such as iscsi. These schedules are used to perform various tasks such as scheduling snapshots or synchronizing replicants.
type Network_Storage_Schedule struct {
	Entity

	// A flag which determines if a schedule is active.
	Active *int `json:"active,omitempty" xmlrpc:"active,omitempty"`

	// The date a schedule was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The hour parameter of this schedule.
	Day *string `json:"day,omitempty" xmlrpc:"day,omitempty"`

	// The day of the month parameter of this schedule.
	DayOfMonth *string `json:"dayOfMonth,omitempty" xmlrpc:"dayOfMonth,omitempty"`

	// The day of the week parameter of this schedule.
	DayOfWeek *string `json:"dayOfWeek,omitempty" xmlrpc:"dayOfWeek,omitempty"`

	// A count of events which have been created as the result of a schedule execution.
	EventCount *uint `json:"eventCount,omitempty" xmlrpc:"eventCount,omitempty"`

	// Events which have been created as the result of a schedule execution.
	Events []Network_Storage_Event `json:"events,omitempty" xmlrpc:"events,omitempty"`

	// The hour parameter of this schedule.
	Hour *string `json:"hour,omitempty" xmlrpc:"hour,omitempty"`

	// A schedule's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The minute parameter of this schedule.
	Minute *string `json:"minute,omitempty" xmlrpc:"minute,omitempty"`

	// The date a schedule was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The month of the year parameter of this schedule.
	MonthOfYear *string `json:"monthOfYear,omitempty" xmlrpc:"monthOfYear,omitempty"`

	// A schedule's name, for example 'Daily'.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The associated partnership for a schedule.
	Partnership *Network_Storage_Partnership `json:"partnership,omitempty" xmlrpc:"partnership,omitempty"`

	// The partnership id which a schedule is associated with.
	PartnershipId *int `json:"partnershipId,omitempty" xmlrpc:"partnershipId,omitempty"`

	// Properties used for configuration of a schedule.
	Properties []Network_Storage_Schedule_Property `json:"properties,omitempty" xmlrpc:"properties,omitempty"`

	// A count of properties used for configuration of a schedule.
	PropertyCount *uint `json:"propertyCount,omitempty" xmlrpc:"propertyCount,omitempty"`

	// A count of replica snapshots which have been created as the result of this schedule's execution.
	ReplicaSnapshotCount *uint `json:"replicaSnapshotCount,omitempty" xmlrpc:"replicaSnapshotCount,omitempty"`

	// Replica snapshots which have been created as the result of this schedule's execution.
	ReplicaSnapshots []Network_Storage `json:"replicaSnapshots,omitempty" xmlrpc:"replicaSnapshots,omitempty"`

	// The number of snapshots this schedule is configured to retain.
	RetentionCount *string `json:"retentionCount,omitempty" xmlrpc:"retentionCount,omitempty"`

	// The minute parameter of this schedule.
	Second *string `json:"second,omitempty" xmlrpc:"second,omitempty"`

	// A count of snapshots which have been created as the result of this schedule's execution.
	SnapshotCount *uint `json:"snapshotCount,omitempty" xmlrpc:"snapshotCount,omitempty"`

	// Snapshots which have been created as the result of this schedule's execution.
	Snapshots []Network_Storage `json:"snapshots,omitempty" xmlrpc:"snapshots,omitempty"`

	// The type provides a standardized definition for a schedule.
	Type *Network_Storage_Schedule_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The type id which a schedule is associated with.
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`

	// The associated volume for a schedule.
	Volume *Network_Storage `json:"volume,omitempty" xmlrpc:"volume,omitempty"`

	// The volume id which a schedule is associated with.
	VolumeId *int `json:"volumeId,omitempty" xmlrpc:"volumeId,omitempty"`
}

// Schedule properties provide attributes such as start date, end date, interval, and other properties to a storage schedule.
type Network_Storage_Schedule_Property struct {
	Entity

	// The date a schedule property was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// A schedule property's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The date a schedule property was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The associated schedule for a property.
	Schedule *Network_Storage_Schedule `json:"schedule,omitempty" xmlrpc:"schedule,omitempty"`

	// The type provides a standardized definition for a property.
	Type *Network_Storage_Schedule_Property_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// An identifier for the type of a property.
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`

	// The value of a property.
	Value *string `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// A schedule property type is used to allow for a standardized method of defining network storage schedules.
type Network_Storage_Schedule_Property_Type struct {
	Entity

	// A type's description, for example 'Date for the schedule to start.'.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// A schedule property type's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A schedule property type's key name, for example 'START_DATE'.
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// A schedule property type's name, for example 'Start Date'.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The type of Storage volume type which a property type may be associated with.
	NasType *string `json:"nasType,omitempty" xmlrpc:"nasType,omitempty"`
}

// A schedule type is used to define what a schedule was created to do. When creating a schedule to take snapshots of a volume, the 'Snapshot' schedule type would be used.
type Network_Storage_Schedule_Type struct {
	Entity

	// A schedule type's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A schedule type's key name, for example 'SNAPSHOT'.
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// A schedule type's name, for example 'Snapshot'.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Network_Storage_Snapshot struct {
	Network_Storage

	// If applicable, the schedule which was executed to create a snapshot.
	CreationSchedule *Network_Storage_Schedule `json:"creationSchedule,omitempty" xmlrpc:"creationSchedule,omitempty"`

	// The volume name for the snapshot.
	VolumeName *string `json:"volumeName,omitempty" xmlrpc:"volumeName,omitempty"`
}

// The SoftLayer_Network_Storage_Type contains a description of the associated SoftLayer_Network_Storage object.
type Network_Storage_Type struct {
	Entity

	// Human readable description for the associated SoftLayer_Network_Storage object.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// ID which corresponds with storageTypeId on storage objects.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Machine readable description code for the associated SoftLayer_Network_Storage object.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// A count of the SoftLayer_Network_Storage object that uses this type.
	VolumeCount *uint `json:"volumeCount,omitempty" xmlrpc:"volumeCount,omitempty"`

	// The SoftLayer_Network_Storage object that uses this type.
	Volumes []Network_Storage `json:"volumes,omitempty" xmlrpc:"volumes,omitempty"`
}

// A subnet represents a continguous range of IP addresses. The range is represented by the networkIdentifer and cidr/netmask properties. The version of a subnet, whether IPv4 or IPv6, is represented by the version property.
//
// When routed, a subnet is associated to a VLAN on your account, which defines its scope on the network. Depending on a subnet's route type, IP addresses may be reserved for network and internal functions, the most common of which is the allocation of network, gateway and broadcast IP addresses.
//
// An unrouted subnet is not active on the network and may generally be routed within the datacenter in which it resides.
//
// [Subnetwork at Wikipedia](http://en.wikipedia.org/wiki/Subnetwork)
//
// [RFC950:Internet Standard Subnetting Procedure](http://datatracker.ietf.org/doc/html/rfc950)
type Network_Subnet struct {
	Entity

	// no documentation yet
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The active regional internet registration for this subnet.
	ActiveRegistration *Network_Subnet_Registration `json:"activeRegistration,omitempty" xmlrpc:"activeRegistration,omitempty"`

	// DEPRECATED
	// Deprecated: This function has been marked as deprecated.
	ActiveSwipTransaction *Network_Subnet_Swip_Transaction `json:"activeSwipTransaction,omitempty" xmlrpc:"activeSwipTransaction,omitempty"`

	// DEPRECATED
	// Deprecated: This function has been marked as deprecated.
	ActiveTransaction *Provisioning_Version1_Transaction `json:"activeTransaction,omitempty" xmlrpc:"activeTransaction,omitempty"`

	// The classifier of IP addresses this subnet represents, generally PUBLIC or PRIVATE. This does not necessarily correlate with the network on which the subnet is used.
	AddressSpace *string `json:"addressSpace,omitempty" xmlrpc:"addressSpace,omitempty"`

	// The link from this subnet to network storage devices supporting access control lists.
	AllowedHost *Network_Storage_Allowed_Host `json:"allowedHost,omitempty" xmlrpc:"allowedHost,omitempty"`

	// The network storage devices this subnet has been granted access to.
	AllowedNetworkStorage []Network_Storage `json:"allowedNetworkStorage,omitempty" xmlrpc:"allowedNetworkStorage,omitempty"`

	// A count of the network storage devices this subnet has been granted access to.
	AllowedNetworkStorageCount *uint `json:"allowedNetworkStorageCount,omitempty" xmlrpc:"allowedNetworkStorageCount,omitempty"`

	// A count of the network storage device replicas this subnet has been granted access to.
	AllowedNetworkStorageReplicaCount *uint `json:"allowedNetworkStorageReplicaCount,omitempty" xmlrpc:"allowedNetworkStorageReplicaCount,omitempty"`

	// The network storage device replicas this subnet has been granted access to.
	AllowedNetworkStorageReplicas []Network_Storage `json:"allowedNetworkStorageReplicas,omitempty" xmlrpc:"allowedNetworkStorageReplicas,omitempty"`

	// The active billing item for this subnet.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// A count of
	BoundDescendantCount *uint `json:"boundDescendantCount,omitempty" xmlrpc:"boundDescendantCount,omitempty"`

	// no documentation yet
	BoundDescendants []Network_Subnet `json:"boundDescendants,omitempty" xmlrpc:"boundDescendants,omitempty"`

	// A count of the list of network routers that this subnet is directly associated with, defining where this subnet may be routed on the network.
	BoundRouterCount *uint `json:"boundRouterCount,omitempty" xmlrpc:"boundRouterCount,omitempty"`

	// Indicates whether this subnet is associated to a network router and is routable on the network.
	BoundRouterFlag *bool `json:"boundRouterFlag,omitempty" xmlrpc:"boundRouterFlag,omitempty"`

	// The list of network routers that this subnet is directly associated with, defining where this subnet may be routed on the network.
	BoundRouters []Hardware `json:"boundRouters,omitempty" xmlrpc:"boundRouters,omitempty"`

	// The IP address of this subnet reserved for use as a broadcast address and which is unavailable for other use. Network traffic targeting this IP address will be broadcast to the entire subnet.
	BroadcastAddress *string `json:"broadcastAddress,omitempty" xmlrpc:"broadcastAddress,omitempty"`

	// The immediate descendants of this subnet.
	Children []Network_Subnet `json:"children,omitempty" xmlrpc:"children,omitempty"`

	// A count of the immediate descendants of this subnet.
	ChildrenCount *uint `json:"childrenCount,omitempty" xmlrpc:"childrenCount,omitempty"`

	// The Classless Inter-Domain Routing prefix of this subnet, which specifies the range of spanned IP addresses.
	//
	// [Classless_Inter-Domain_Routing at Wikipedia](http://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
	Cidr *int `json:"cidr,omitempty" xmlrpc:"cidr,omitempty"`

	// The datacenter this subnet is primarily associated with.
	Datacenter *Location_Datacenter `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// A count of the descendants of this subnet, including all parents and children.
	DescendantCount *uint `json:"descendantCount,omitempty" xmlrpc:"descendantCount,omitempty"`

	// The descendants of this subnet, including all parents and children.
	Descendants []Network_Subnet `json:"descendants,omitempty" xmlrpc:"descendants,omitempty"`

	// [DEPRECATED] The description of this subnet.
	// Deprecated: This function has been marked as deprecated.
	DisplayLabel *string `json:"displayLabel,omitempty" xmlrpc:"displayLabel,omitempty"`

	// The IP address target of this statically routed subnet.
	EndPointIpAddress *Network_Subnet_IpAddress `json:"endPointIpAddress,omitempty" xmlrpc:"endPointIpAddress,omitempty"`

	// The IP address of this subnet reserved for use on the router as a gateway address and which is unavailable for other use.
	Gateway *string `json:"gateway,omitempty" xmlrpc:"gateway,omitempty"`

	// no documentation yet
	GlobalIpRecord *Network_Subnet_IpAddress_Global `json:"globalIpRecord,omitempty" xmlrpc:"globalIpRecord,omitempty"`

	// The Bare Metal devices which have been assigned a primary IP address from this subnet.
	Hardware []Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// A count of the Bare Metal devices which have been assigned a primary IP address from this subnet.
	HardwareCount *uint `json:"hardwareCount,omitempty" xmlrpc:"hardwareCount,omitempty"`

	// The unique identifier of this subnet.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of the IP address records belonging to this subnet.
	IpAddressCount *uint `json:"ipAddressCount,omitempty" xmlrpc:"ipAddressCount,omitempty"`

	// The IP address records belonging to this subnet.
	IpAddresses []Network_Subnet_IpAddress `json:"ipAddresses,omitempty" xmlrpc:"ipAddresses,omitempty"`

	// Indicates whether this subnet is owned by the assigned account.
	IsCustomerOwned *bool `json:"isCustomerOwned,omitempty" xmlrpc:"isCustomerOwned,omitempty"`

	// Indicates whether the route type of this subnet may be altered.
	IsCustomerRoutable *bool `json:"isCustomerRoutable,omitempty" xmlrpc:"isCustomerRoutable,omitempty"`

	// The time this subnet was last modified
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The bitmask in dotted-quad format for this subnet, which specifies the range of spanned IP addresses.
	Netmask *string `json:"netmask,omitempty" xmlrpc:"netmask,omitempty"`

	// The hardware firewall associated to this subnet via access control list.
	NetworkComponentFirewall *Network_Component_Firewall `json:"networkComponentFirewall,omitempty" xmlrpc:"networkComponentFirewall,omitempty"`

	// The first IP address of this subnet.
	NetworkIdentifier *string `json:"networkIdentifier,omitempty" xmlrpc:"networkIdentifier,omitempty"`

	// A count of
	NetworkProtectionAddressCount *uint `json:"networkProtectionAddressCount,omitempty" xmlrpc:"networkProtectionAddressCount,omitempty"`

	// no documentation yet
	NetworkProtectionAddresses []Network_Protection_Address `json:"networkProtectionAddresses,omitempty" xmlrpc:"networkProtectionAddresses,omitempty"`

	// A count of the IPSec VPN tunnels associated to this subnet.
	NetworkTunnelContextCount *uint `json:"networkTunnelContextCount,omitempty" xmlrpc:"networkTunnelContextCount,omitempty"`

	// The IPSec VPN tunnels associated to this subnet.
	NetworkTunnelContexts []Network_Tunnel_Module_Context `json:"networkTunnelContexts,omitempty" xmlrpc:"networkTunnelContexts,omitempty"`

	// The VLAN this subnet is associated with.
	NetworkVlan *Network_Vlan `json:"networkVlan,omitempty" xmlrpc:"networkVlan,omitempty"`

	// The identifier of the VLAN associated to this subnet.
	NetworkVlanId *int `json:"networkVlanId,omitempty" xmlrpc:"networkVlanId,omitempty"`

	// The customer description of this subnet.
	Note *string `json:"note,omitempty" xmlrpc:"note,omitempty"`

	// The pod in which this subnet is currently routed.
	PodName *string `json:"podName,omitempty" xmlrpc:"podName,omitempty"`

	// A count of
	ProtectedIpAddressCount *uint `json:"protectedIpAddressCount,omitempty" xmlrpc:"protectedIpAddressCount,omitempty"`

	// no documentation yet
	ProtectedIpAddresses []Network_Subnet_IpAddress `json:"protectedIpAddresses,omitempty" xmlrpc:"protectedIpAddresses,omitempty"`

	// The RIR which is authoritative over the network in which this subnet resides.
	RegionalInternetRegistry *Network_Regional_Internet_Registry `json:"regionalInternetRegistry,omitempty" xmlrpc:"regionalInternetRegistry,omitempty"`

	// A count of the regional internet registrations that have been created for this subnet.
	RegistrationCount *uint `json:"registrationCount,omitempty" xmlrpc:"registrationCount,omitempty"`

	// The regional internet registrations that have been created for this subnet.
	Registrations []Network_Subnet_Registration `json:"registrations,omitempty" xmlrpc:"registrations,omitempty"`

	// The reverse DNS domain associated with this subnet.
	ReverseDomain *Dns_Domain `json:"reverseDomain,omitempty" xmlrpc:"reverseDomain,omitempty"`

	// The role identifier that this subnet is participating in. Roles dictate how a subnet may be used.
	RoleKeyName *string `json:"roleKeyName,omitempty" xmlrpc:"roleKeyName,omitempty"`

	// The name of the role the subnet is within. Roles dictate how a subnet may be used.
	RoleName *string `json:"roleName,omitempty" xmlrpc:"roleName,omitempty"`

	// The product and route classifier for this routed subnet, with the following values: PRIMARY, SECONDARY, STATIC_TO_IP, GLOBAL_IP, IPSEC_STATIC_NAT.
	RoutingTypeKeyName *string `json:"routingTypeKeyName,omitempty" xmlrpc:"routingTypeKeyName,omitempty"`

	// The description of the product and route classifier for this routed subnet, with the following values: Primary, Portable, Static, Global, IPSec Static NAT.
	RoutingTypeName *string `json:"routingTypeName,omitempty" xmlrpc:"routingTypeName,omitempty"`

	// [DEPRECATED] Used to sort subnets and group subnets of similar type together for use on customer facing portals.
	// Deprecated: This function has been marked as deprecated.
	SortOrder *string `json:"sortOrder,omitempty" xmlrpc:"sortOrder,omitempty"`

	// The product and route classifier for this routed subnet, with the following values:
	// * PRIMARY
	// * ADDITIONAL_PRIMARY
	// * SECONDARY_ON_VLAN
	// * STATIC_IP_ROUTED
	// * PRIMARY_6
	// * SUBNET_ON_VLAN
	// * STATIC_IP_ROUTED_6
	// * GLOBAL_IP
	// * IPSEC_STATIC_NAT
	//
	//
	// "PRIMARY" refers to the principal IPv4 network from which primary IP addresses are assigned to devices.
	//
	// "ADDITIONAL_PRIMARY" refers to extra IPv4 networks from which primary IP addresses are assigned to devices.
	//
	// "SECONDARY_ON_VLAN" refers to a secondary IPv4 subnet routed as portable.
	//
	// "STATIC_IP_ROUTED" refers to a secondary IPv4 subnet routed as static to a single endpoint IPv4 address.
	//
	// "PRIMARY_6" refers to the IPv6 network from which primary IPv6 addresses are assigned to devices.
	//
	// "SUBNET_ON_VLAN" refers to a secondary IPv6 subnet routed as portable.
	//
	// "STATIC_IP_ROUTED_6" refers to a secondary IPv6 subnet routed as static to a single endpoint IPv6 address.
	//
	// "GLOBAL_IP" refers to a global IPv4/IPv6 address routed as static to a single endpoint IP address.
	//
	// "IPSEC_STATIC_NAT" refers to the networks associated to your IPSec VPN tunnels for NAT purposes.
	SubnetType *string `json:"subnetType,omitempty" xmlrpc:"subnetType,omitempty"`

	// DEPRECATED
	// Deprecated: This function has been marked as deprecated.
	SwipTransaction []Network_Subnet_Swip_Transaction `json:"swipTransaction,omitempty" xmlrpc:"swipTransaction,omitempty"`

	// A count of dEPRECATED
	SwipTransactionCount *uint `json:"swipTransactionCount,omitempty" xmlrpc:"swipTransactionCount,omitempty"`

	// A count of the tags associated to this subnet.
	TagReferenceCount *uint `json:"tagReferenceCount,omitempty" xmlrpc:"tagReferenceCount,omitempty"`

	// The tags associated to this subnet.
	TagReferences []Tag_Reference `json:"tagReferences,omitempty" xmlrpc:"tagReferences,omitempty"`

	// The number of IP addresses in this subnet.
	TotalIpAddresses *Float64 `json:"totalIpAddresses,omitempty" xmlrpc:"totalIpAddresses,omitempty"`

	// A count of
	UnboundDescendantCount *uint `json:"unboundDescendantCount,omitempty" xmlrpc:"unboundDescendantCount,omitempty"`

	// no documentation yet
	UnboundDescendants []Network_Subnet `json:"unboundDescendants,omitempty" xmlrpc:"unboundDescendants,omitempty"`

	// The number of IP addresses that can be addressed within this subnet. For IPv4 subnets with a CIDR value of at most 30, a discount of 3 is taken from the total number of IP addresses for the subnet's unusable network, gateway and broadcast IP addresses. For IPv6 subnets with a CIDR value of at most 126, a discount of 2 is taken for the subnet's network and gateway IP addresses.
	UsableIpAddressCount *Float64 `json:"usableIpAddressCount,omitempty" xmlrpc:"usableIpAddressCount,omitempty"`

	// The total number of utilized IP addresses on this subnet. The primary consumer of IP addresses are compute resources, which can consume more than one address. This value is only supported for primary subnets.
	UtilizedIpAddressCount *uint `json:"utilizedIpAddressCount,omitempty" xmlrpc:"utilizedIpAddressCount,omitempty"`

	// The Internet Protocol version of this subnet, either 4 or 6.
	Version *int `json:"version,omitempty" xmlrpc:"version,omitempty"`

	// A count of the Virtual Server devices which have been assigned a primary IP address from this subnet.
	VirtualGuestCount *uint `json:"virtualGuestCount,omitempty" xmlrpc:"virtualGuestCount,omitempty"`

	// The Virtual Server devices which have been assigned a primary IP address from this subnet.
	VirtualGuests []Virtual_Guest `json:"virtualGuests,omitempty" xmlrpc:"virtualGuests,omitempty"`
}

// The SoftLayer_Network_Subnet_IpAddress data type contains general information relating to a single SoftLayer IPv4 address.
type Network_Subnet_IpAddress struct {
	Entity

	// The SoftLayer_Network_Storage_Allowed_Host information to connect this IP Address to Network Storage supporting access control lists.
	AllowedHost *Network_Storage_Allowed_Host `json:"allowedHost,omitempty" xmlrpc:"allowedHost,omitempty"`

	// The SoftLayer_Network_Storage objects that this SoftLayer_Hardware has access to.
	AllowedNetworkStorage []Network_Storage `json:"allowedNetworkStorage,omitempty" xmlrpc:"allowedNetworkStorage,omitempty"`

	// A count of the SoftLayer_Network_Storage objects that this SoftLayer_Hardware has access to.
	AllowedNetworkStorageCount *uint `json:"allowedNetworkStorageCount,omitempty" xmlrpc:"allowedNetworkStorageCount,omitempty"`

	// A count of the SoftLayer_Network_Storage objects whose Replica that this SoftLayer_Hardware has access to.
	AllowedNetworkStorageReplicaCount *uint `json:"allowedNetworkStorageReplicaCount,omitempty" xmlrpc:"allowedNetworkStorageReplicaCount,omitempty"`

	// The SoftLayer_Network_Storage objects whose Replica that this SoftLayer_Hardware has access to.
	AllowedNetworkStorageReplicas []Network_Storage `json:"allowedNetworkStorageReplicas,omitempty" xmlrpc:"allowedNetworkStorageReplicas,omitempty"`

	// The application delivery controller using this address.
	ApplicationDeliveryController *Network_Application_Delivery_Controller `json:"applicationDeliveryController,omitempty" xmlrpc:"applicationDeliveryController,omitempty"`

	// A count of an IPSec network tunnel's address translations. These translations use a SoftLayer ip address from an assigned static NAT subnet to deliver the packets to the remote (customer) destination.
	ContextTunnelTranslationCount *uint `json:"contextTunnelTranslationCount,omitempty" xmlrpc:"contextTunnelTranslationCount,omitempty"`

	// An IPSec network tunnel's address translations. These translations use a SoftLayer ip address from an assigned static NAT subnet to deliver the packets to the remote (customer) destination.
	ContextTunnelTranslations []Network_Tunnel_Module_Context_Address_Translation `json:"contextTunnelTranslations,omitempty" xmlrpc:"contextTunnelTranslations,omitempty"`

	// A count of all the subnets routed to an IP address.
	EndpointSubnetCount *uint `json:"endpointSubnetCount,omitempty" xmlrpc:"endpointSubnetCount,omitempty"`

	// All the subnets routed to an IP address.
	EndpointSubnets []Network_Subnet `json:"endpointSubnets,omitempty" xmlrpc:"endpointSubnets,omitempty"`

	// A network component that is statically routed to an IP address.
	GuestNetworkComponent *Virtual_Guest_Network_Component `json:"guestNetworkComponent,omitempty" xmlrpc:"guestNetworkComponent,omitempty"`

	// A network component that is statically routed to an IP address.
	GuestNetworkComponentBinding *Virtual_Guest_Network_Component_IpAddress `json:"guestNetworkComponentBinding,omitempty" xmlrpc:"guestNetworkComponentBinding,omitempty"`

	// A server that this IP address is routed to.
	Hardware *Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// An IP's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// An IP address expressed in dotted quad format.
	IpAddress *string `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// Indicates if an IP address is reserved to be used as the network broadcast address and cannot be assigned to a network interface
	IsBroadcast *bool `json:"isBroadcast,omitempty" xmlrpc:"isBroadcast,omitempty"`

	// Indicates if an IP address is reserved to a gateway and cannot be assigned to a network interface
	IsGateway *bool `json:"isGateway,omitempty" xmlrpc:"isGateway,omitempty"`

	// Indicates if an IP address is reserved to a network address and cannot be assigned to a network interface
	IsNetwork *bool `json:"isNetwork,omitempty" xmlrpc:"isNetwork,omitempty"`

	// Indicates if an IP address is reserved and cannot be assigned to a network interface
	IsReserved *bool `json:"isReserved,omitempty" xmlrpc:"isReserved,omitempty"`

	// A network component that is statically routed to an IP address.
	NetworkComponent *Network_Component `json:"networkComponent,omitempty" xmlrpc:"networkComponent,omitempty"`

	// An IP address' user defined note.
	Note *string `json:"note,omitempty" xmlrpc:"note,omitempty"`

	// The network gateway appliance using this address as the private IP address.
	PrivateNetworkGateway *Network_Gateway `json:"privateNetworkGateway,omitempty" xmlrpc:"privateNetworkGateway,omitempty"`

	// no documentation yet
	ProtectionAddress []Network_Protection_Address `json:"protectionAddress,omitempty" xmlrpc:"protectionAddress,omitempty"`

	// A count of
	ProtectionAddressCount *uint `json:"protectionAddressCount,omitempty" xmlrpc:"protectionAddressCount,omitempty"`

	// The network gateway appliance using this address as the public IP address.
	PublicNetworkGateway *Network_Gateway `json:"publicNetworkGateway,omitempty" xmlrpc:"publicNetworkGateway,omitempty"`

	// An IPMI-based management network component of the IP address.
	RemoteManagementNetworkComponent *Network_Component `json:"remoteManagementNetworkComponent,omitempty" xmlrpc:"remoteManagementNetworkComponent,omitempty"`

	// An IP address' associated subnet.
	Subnet *Network_Subnet `json:"subnet,omitempty" xmlrpc:"subnet,omitempty"`

	// An IP address' subnet id.
	SubnetId *int `json:"subnetId,omitempty" xmlrpc:"subnetId,omitempty"`

	// All events for this IP address stored in the datacenter syslogs from the last 24 hours
	SyslogEventsOneDay []Network_Logging_Syslog `json:"syslogEventsOneDay,omitempty" xmlrpc:"syslogEventsOneDay,omitempty"`

	// A count of all events for this IP address stored in the datacenter syslogs from the last 24 hours
	SyslogEventsOneDayCount *uint `json:"syslogEventsOneDayCount,omitempty" xmlrpc:"syslogEventsOneDayCount,omitempty"`

	// A count of all events for this IP address stored in the datacenter syslogs from the last 7 days
	SyslogEventsSevenDayCount *uint `json:"syslogEventsSevenDayCount,omitempty" xmlrpc:"syslogEventsSevenDayCount,omitempty"`

	// All events for this IP address stored in the datacenter syslogs from the last 7 days
	SyslogEventsSevenDays []Network_Logging_Syslog `json:"syslogEventsSevenDays,omitempty" xmlrpc:"syslogEventsSevenDays,omitempty"`

	// Top Ten network datacenter syslog events, grouped by destination port, for the last 24 hours
	TopTenSyslogEventsByDestinationPortOneDay []Network_Logging_Syslog `json:"topTenSyslogEventsByDestinationPortOneDay,omitempty" xmlrpc:"topTenSyslogEventsByDestinationPortOneDay,omitempty"`

	// A count of top Ten network datacenter syslog events, grouped by destination port, for the last 24 hours
	TopTenSyslogEventsByDestinationPortOneDayCount *uint `json:"topTenSyslogEventsByDestinationPortOneDayCount,omitempty" xmlrpc:"topTenSyslogEventsByDestinationPortOneDayCount,omitempty"`

	// A count of top Ten network datacenter syslog events, grouped by destination port, for the last 7 days
	TopTenSyslogEventsByDestinationPortSevenDayCount *uint `json:"topTenSyslogEventsByDestinationPortSevenDayCount,omitempty" xmlrpc:"topTenSyslogEventsByDestinationPortSevenDayCount,omitempty"`

	// Top Ten network datacenter syslog events, grouped by destination port, for the last 7 days
	TopTenSyslogEventsByDestinationPortSevenDays []Network_Logging_Syslog `json:"topTenSyslogEventsByDestinationPortSevenDays,omitempty" xmlrpc:"topTenSyslogEventsByDestinationPortSevenDays,omitempty"`

	// Top Ten network datacenter syslog events, grouped by source port, for the last 24 hours
	TopTenSyslogEventsByProtocolsOneDay []Network_Logging_Syslog `json:"topTenSyslogEventsByProtocolsOneDay,omitempty" xmlrpc:"topTenSyslogEventsByProtocolsOneDay,omitempty"`

	// A count of top Ten network datacenter syslog events, grouped by source port, for the last 24 hours
	TopTenSyslogEventsByProtocolsOneDayCount *uint `json:"topTenSyslogEventsByProtocolsOneDayCount,omitempty" xmlrpc:"topTenSyslogEventsByProtocolsOneDayCount,omitempty"`

	// A count of top Ten network datacenter syslog events, grouped by source port, for the last 7 days
	TopTenSyslogEventsByProtocolsSevenDayCount *uint `json:"topTenSyslogEventsByProtocolsSevenDayCount,omitempty" xmlrpc:"topTenSyslogEventsByProtocolsSevenDayCount,omitempty"`

	// Top Ten network datacenter syslog events, grouped by source port, for the last 7 days
	TopTenSyslogEventsByProtocolsSevenDays []Network_Logging_Syslog `json:"topTenSyslogEventsByProtocolsSevenDays,omitempty" xmlrpc:"topTenSyslogEventsByProtocolsSevenDays,omitempty"`

	// Top Ten network datacenter syslog events, grouped by source ip address, for the last 24 hours
	TopTenSyslogEventsBySourceIpOneDay []Network_Logging_Syslog `json:"topTenSyslogEventsBySourceIpOneDay,omitempty" xmlrpc:"topTenSyslogEventsBySourceIpOneDay,omitempty"`

	// A count of top Ten network datacenter syslog events, grouped by source ip address, for the last 24 hours
	TopTenSyslogEventsBySourceIpOneDayCount *uint `json:"topTenSyslogEventsBySourceIpOneDayCount,omitempty" xmlrpc:"topTenSyslogEventsBySourceIpOneDayCount,omitempty"`

	// A count of top Ten network datacenter syslog events, grouped by source ip address, for the last 7 days
	TopTenSyslogEventsBySourceIpSevenDayCount *uint `json:"topTenSyslogEventsBySourceIpSevenDayCount,omitempty" xmlrpc:"topTenSyslogEventsBySourceIpSevenDayCount,omitempty"`

	// Top Ten network datacenter syslog events, grouped by source ip address, for the last 7 days
	TopTenSyslogEventsBySourceIpSevenDays []Network_Logging_Syslog `json:"topTenSyslogEventsBySourceIpSevenDays,omitempty" xmlrpc:"topTenSyslogEventsBySourceIpSevenDays,omitempty"`

	// Top Ten network datacenter syslog events, grouped by source port, for the last 24 hours
	TopTenSyslogEventsBySourcePortOneDay []Network_Logging_Syslog `json:"topTenSyslogEventsBySourcePortOneDay,omitempty" xmlrpc:"topTenSyslogEventsBySourcePortOneDay,omitempty"`

	// A count of top Ten network datacenter syslog events, grouped by source port, for the last 24 hours
	TopTenSyslogEventsBySourcePortOneDayCount *uint `json:"topTenSyslogEventsBySourcePortOneDayCount,omitempty" xmlrpc:"topTenSyslogEventsBySourcePortOneDayCount,omitempty"`

	// A count of top Ten network datacenter syslog events, grouped by source port, for the last 7 days
	TopTenSyslogEventsBySourcePortSevenDayCount *uint `json:"topTenSyslogEventsBySourcePortSevenDayCount,omitempty" xmlrpc:"topTenSyslogEventsBySourcePortSevenDayCount,omitempty"`

	// Top Ten network datacenter syslog events, grouped by source port, for the last 7 days
	TopTenSyslogEventsBySourcePortSevenDays []Network_Logging_Syslog `json:"topTenSyslogEventsBySourcePortSevenDays,omitempty" xmlrpc:"topTenSyslogEventsBySourcePortSevenDays,omitempty"`

	// A virtual guest that this IP address is routed to.
	VirtualGuest *Virtual_Guest `json:"virtualGuest,omitempty" xmlrpc:"virtualGuest,omitempty"`

	// A count of virtual licenses allocated for an IP Address.
	VirtualLicenseCount *uint `json:"virtualLicenseCount,omitempty" xmlrpc:"virtualLicenseCount,omitempty"`

	// Virtual licenses allocated for an IP Address.
	VirtualLicenses []Software_VirtualLicense `json:"virtualLicenses,omitempty" xmlrpc:"virtualLicenses,omitempty"`
}

// no documentation yet
type Network_Subnet_IpAddress_Global struct {
	Entity

	// no documentation yet
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// DEPRECATED
	// Deprecated: This function has been marked as deprecated.
	ActiveTransaction *Provisioning_Version1_Transaction `json:"activeTransaction,omitempty" xmlrpc:"activeTransaction,omitempty"`

	// The billing item for this Global IP.
	BillingItem *Billing_Item_Network_Subnet_IpAddress_Global `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// A Global IP Address' associated description
	Description *int `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	DestinationIpAddress *Network_Subnet_IpAddress `json:"destinationIpAddress,omitempty" xmlrpc:"destinationIpAddress,omitempty"`

	// A Global IP Address' associated [[SoftLayer_Network_Subnet_IpAddress|ipAddress]] ID
	DestinationIpAddressId *int `json:"destinationIpAddressId,omitempty" xmlrpc:"destinationIpAddressId,omitempty"`

	// A Global IP Address' unique identifier
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	IpAddress *Network_Subnet_IpAddress `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// A Global IP Address' associated [[SoftLayer_Account|account]] ID
	IpAddressId *int `json:"ipAddressId,omitempty" xmlrpc:"ipAddressId,omitempty"`

	// A Global IP Address' associated type [[SoftLayer_Network_Subnet_IpAddress_Global_Type|id]] ID
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`
}

// Describes an IP address assigned to a resource on your network.
//
// Details on the associated resource are also provided, described below. Details include the resource's type, unique identifier, name, fully qualified name, and context, the contents of which depends on the resource's type. If the fully qualified name is not included for a resource type below, the resource's name will apply.
//
// The following resource types and associated dependent properties are supported:
//
// * <b>HARDWARE</b>: A [Bare Metal Server](/reference/datatypes/SoftLayer_Hardware_Server)
//
// -- <i>resourceName</i>: The hostname of the server.
//
// -- <i>resourceFullyQualifiedName</i>: The fully qualified domain name of the server.
//
// -- <i>resourceContext</i>: The name of the network component or network component group assigned to the IP address, <i>e.g. eth0/2</i>.
//
// * <b>GUEST</b>: A [Virtual Server Instance](/reference/datatypes/SoftLayer_Virtual_Guest)
//
// -- <i>resourceName</i>: The hostname of the guest.
//
// -- <i>resourceFullyQualifiedName</i>: The fully qualified domain name of the guest.
//
// -- <i>resourceContext</i>: The name of the virtual network component assigned to the IP address, <i>e.g. eth0</i>.
//
// * <b>GATEWAY</b>: A [Network Gateway](/reference/datatypes/SoftLayer_Network_Gateway)
//
// -- <i>resourceName</i>: The name of the gateway.
//
// -- <i>resourceContext</i>: Either the term "virtual" to indicate a gateway IP address, or the name of the network component or network component group assigned to the IP address followed by the id-value of the [Bare Metal Server](/reference/datatypes/SoftLayer_Hardware_Server) gateway member surrounded by '<', '>', <i>e.g. eth1/3<123456></i>.
//
// - <b>FIREWALL_MULTIVLAN</b>: A [Multi-VLAN Firewall](/reference/datatypes/SoftLayer_Network_Vlan_Firewall)
//
// -- <i>resourceName</i>: The name of the firewall.
//
// -- <i>resourceContext</i>: The term "virtual" to indicate a firewall IP address.
//
// - <b>LBAAS</b>: A [Cloud Load Balancer](/reference/datatypes/SoftLayer_Network_LBaaS_LoadBalancer)
//
// -- <i>resourceName</i>: The name of the load balancer.
//
// -- <i>resourceFullyQualifiedName</i>: The full DNS address of the load balancer.
//
// -- <i>resourceContext</i>: The term "ephemeral" to indicate a currently assigned IP address, subject to change. Users are strongly encouraged to access the service by the fully qualified DNS name and not the underlying IP addresses. The UUID of the load balancer is also provided, surrounded by '<' and '>', e.g. ephemeral<84f0affb-0d5e-40f1-ad87-a92d6544936a>
//
// - <b>NETSCALER_VPX</b>: A [Netscaler VPX Load Balancer](/reference/datatypes/SoftLayer_Network_Application_Delivery_Controller)
//
// -- <i>resourceName</i>: The hostname of the load balancer.
//
// -- <i>resourceFullyQualifiedName</i>: The fully qualified domain name of the load balancer.
//
// -- <i>resourceContext</i>: Either the term "nsip" to indicate the management IP address, or the name of the network component assigned to the IP address followed by the id-value of the [Virtual Server Instance](/reference/datatypes/SoftLayer_Virtual_Guest) load balancer host surrounded by '<', '>', <i>e.g. eth1<123456></i>.
//
// - <b>NETSCALER_MPX</b>: A [Netscaler MPX Load Balancer](/reference/datatypes/SoftLayer_Hardware_LoadBalancer)
//
// -- <i>resourceName</i>: The hostname of the load balancer.
//
// -- <i>resourceFullyQualifiedName</i>: The fully qualified domain name of the load balancer.
//
// -- <i>resourceContext</i>: The name of the network component or network component group assigned to the IP address, <i>e.g. eth0/2</i>.
type Network_Subnet_IpAddress_UsageDetail struct {
	Entity

	// The IP address.
	IpAddress *string `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// The unique identifier of the IP address record.
	IpAddressId *int `json:"ipAddressId,omitempty" xmlrpc:"ipAddressId,omitempty"`

	// A description of the resource IP address assignment.
	ResourceContext *string `json:"resourceContext,omitempty" xmlrpc:"resourceContext,omitempty"`

	// The fully qualified name of the assigned resource.
	ResourceFullyQualifiedName *string `json:"resourceFullyQualifiedName,omitempty" xmlrpc:"resourceFullyQualifiedName,omitempty"`

	// The unique identifier of the assigned resource.
	ResourceId *int `json:"resourceId,omitempty" xmlrpc:"resourceId,omitempty"`

	// The name of the assigned resource.
	ResourceName *string `json:"resourceName,omitempty" xmlrpc:"resourceName,omitempty"`

	// The type of the assigned resource.
	ResourceType *string `json:"resourceType,omitempty" xmlrpc:"resourceType,omitempty"`

	// The unique identifier of the subnet the IP address belongs to.
	SubnetId *int `json:"subnetId,omitempty" xmlrpc:"subnetId,omitempty"`
}

// The SoftLayer_Network_Subnet_IpAddress data type contains general information relating to a single SoftLayer IPv6 address.
type Network_Subnet_IpAddress_Version6 struct {
	Network_Subnet_IpAddress

	// The network gateway appliance using this address as the public IPv6 address.
	PublicVersion6NetworkGateway *Network_Gateway `json:"publicVersion6NetworkGateway,omitempty" xmlrpc:"publicVersion6NetworkGateway,omitempty"`
}

// The subnet registration service has been deprecated.
//
// The subnet registration data type contains general information relating to a single subnet registration instance. These registration instances can be updated to reflect changes, and will record the changes in the [[SoftLayer_Network_Subnet_Registration_Event|events]].
type Network_Subnet_Registration struct {
	Entity

	// [Deprecated] The account that this registration belongs to.
	// Deprecated: This function has been marked as deprecated.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The registration object's associated [[SoftLayer_Account|account]] id
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The CIDR prefix for the registered subnet
	Cidr *int `json:"cidr,omitempty" xmlrpc:"cidr,omitempty"`

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// A count of [Deprecated] The cross-reference records that tie the [[SoftLayer_Account_Regional_Registry_Detail]] objects to the registration object.
	DetailReferenceCount *uint `json:"detailReferenceCount,omitempty" xmlrpc:"detailReferenceCount,omitempty"`

	// [Deprecated] The cross-reference records that tie the [[SoftLayer_Account_Regional_Registry_Detail]] objects to the registration object.
	// Deprecated: This function has been marked as deprecated.
	DetailReferences []Network_Subnet_Registration_Details `json:"detailReferences,omitempty" xmlrpc:"detailReferences,omitempty"`

	// A count of [Deprecated] The related registration events.
	EventCount *uint `json:"eventCount,omitempty" xmlrpc:"eventCount,omitempty"`

	// [Deprecated] The related registration events.
	// Deprecated: This function has been marked as deprecated.
	Events []Network_Subnet_Registration_Event `json:"events,omitempty" xmlrpc:"events,omitempty"`

	// Unique ID of the registration object
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// [Deprecated] The "network" detail object.
	// Deprecated: This function has been marked as deprecated.
	NetworkDetail *Account_Regional_Registry_Detail `json:"networkDetail,omitempty" xmlrpc:"networkDetail,omitempty"`

	// The RIR-specific handle or name of the registered subnet. This field is read-only.
	NetworkHandle *string `json:"networkHandle,omitempty" xmlrpc:"networkHandle,omitempty"`

	// The base IP address of the registered subnet
	NetworkIdentifier *string `json:"networkIdentifier,omitempty" xmlrpc:"networkIdentifier,omitempty"`

	// [Deprecated] The "person" detail object.
	// Deprecated: This function has been marked as deprecated.
	PersonDetail *Account_Regional_Registry_Detail `json:"personDetail,omitempty" xmlrpc:"personDetail,omitempty"`

	// [Deprecated] The related Regional Internet Registry.
	// Deprecated: This function has been marked as deprecated.
	RegionalInternetRegistry *Network_Regional_Internet_Registry `json:"regionalInternetRegistry,omitempty" xmlrpc:"regionalInternetRegistry,omitempty"`

	// [Deprecated] The RIR handle that this registration object belongs to. This field may not be populated until the registration is complete.
	// Deprecated: This function has been marked as deprecated.
	RegionalInternetRegistryHandle *Account_Rwhois_Handle `json:"regionalInternetRegistryHandle,omitempty" xmlrpc:"regionalInternetRegistryHandle,omitempty"`

	// The registration object's associated [[SoftLayer_Account_Rwhois_Handle|RIR handle]] id
	RegionalInternetRegistryHandleId *int `json:"regionalInternetRegistryHandleId,omitempty" xmlrpc:"regionalInternetRegistryHandleId,omitempty"`

	// The registration object's associated [[SoftLayer_Network_Regional_Internet_Registry|RIR]] id
	RegionalInternetRegistryId *int `json:"regionalInternetRegistryId,omitempty" xmlrpc:"regionalInternetRegistryId,omitempty"`

	// [Deprecated] The status of this registration.
	// Deprecated: This function has been marked as deprecated.
	Status *Network_Subnet_Registration_Status `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// The registration object's associated [[SoftLayer_Network_Subnet_Registration_Status|status]] id
	StatusId *int `json:"statusId,omitempty" xmlrpc:"statusId,omitempty"`

	// [Deprecated] The subnet that this registration pertains to.
	// Deprecated: This function has been marked as deprecated.
	Subnet *Network_Subnet `json:"subnet,omitempty" xmlrpc:"subnet,omitempty"`
}

// The APNIC subnet registration type has been deprecated.
//
// APNIC-specific registration object. For more detail see [[SoftLayer_Network_Subnet_Registration (type)|SoftLayer_Network_Subnet_Registration]].
type Network_Subnet_Registration_Apnic struct {
	Network_Subnet_Registration
}

// The ARIN subnet registration type has been deprecated.
//
// ARIN-specific registration object. For more detail see [[SoftLayer_Network_Subnet_Registration (type)|SoftLayer_Network_Subnet_Registration]].
type Network_Subnet_Registration_Arin struct {
	Network_Subnet_Registration
}

// The subnet registration details type has been deprecated.
//
// The SoftLayer_Network_Subnet_Registration_Details objects are used to relate [[SoftLayer_Account_Regional_Registry_Detail]] objects to a [[SoftLayer_Network_Subnet_Registration]] object. This allows for easy reuse of registration details. It is important to note that only one detail object per type may be associated to a registration object.
type Network_Subnet_Registration_Details struct {
	Entity

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// [Deprecated] The related [[SoftLayer_Account_Regional_Registry_Detail|detail object]].
	// Deprecated: This function has been marked as deprecated.
	Detail *Account_Regional_Registry_Detail `json:"detail,omitempty" xmlrpc:"detail,omitempty"`

	// Numeric ID of the related [[SoftLayer_Account_Regional_Registry_Detail]] object
	DetailId *int `json:"detailId,omitempty" xmlrpc:"detailId,omitempty"`

	// Unique numeric ID of the object
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// [Deprecated] The related [[SoftLayer_Network_Subnet_Registration|registration object]].
	// Deprecated: This function has been marked as deprecated.
	Registration *Network_Subnet_Registration `json:"registration,omitempty" xmlrpc:"registration,omitempty"`

	// Numeric ID of the related [[SoftLayer_Network_Subnet_Registration]] object
	RegistrationId *int `json:"registrationId,omitempty" xmlrpc:"registrationId,omitempty"`
}

// The subnet registration event type has been deprecated.
//
// Each time a [[SoftLayer_Network_Subnet_Registration|subnet registration]] object is created or modified, the system will generate an event for it. Additional actions that would create an event include RIR responses and error cases. *
type Network_Subnet_Registration_Event struct {
	Entity

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// Unique numeric ID of the event object
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A string message indicating what took place during this event
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// [Deprecated] The registration this event pertains to.
	// Deprecated: This function has been marked as deprecated.
	Registration *Network_Subnet_Registration `json:"registration,omitempty" xmlrpc:"registration,omitempty"`

	// The numeric ID of the related [[SoftLayer_Network_Subnet_Registration]] object
	RegistrationId *int `json:"registrationId,omitempty" xmlrpc:"registrationId,omitempty"`

	// [Deprecated] The type of this event.
	// Deprecated: This function has been marked as deprecated.
	Type *Network_Subnet_Registration_Event_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// The numeric ID of the associated [[SoftLayer_Network_Subnet_Registration_Event_Type|event type]] object
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`
}

// The subnet registration event type type has been deprecated.
//
// Subnet Registration Event Type objects describe the nature of a [[SoftLayer_Network_Subnet_Registration_Event]]
//
// The standard values for these objects are as follows: <ul> <li><strong>REGISTRATION_CREATED</strong> - Indicates that the registration has been created</li> <li><strong>REGISTRATION_UPDATED</strong> - Indicates that the registration has been updated</li> <li><strong>REGISTRATION_CANCELLED</strong> - Indicates that the registration has been cancelled</li> <li><strong>RIR_RESPONSE</strong> - Indicates that an action taken against the RIR has produced a response. More details will be provided in the event message.</li> <li><strong>ERROR</strong> - Indicates that an error has been encountered. More details will be provided in the event message.</li> <li><strong>NOTE</strong> - An employee or other system has entered a note regarding the registration. The note content will be provided in the event message.</li> </ul>
type Network_Subnet_Registration_Event_Type struct {
	Entity

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// Unique numeric ID of the event type object
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Code-friendly string name of the event type
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Human-readable name of the event type
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The RIPE subnet registration type has been deprecated.
//
// RIPE-specific registration object. For more detail see [[SoftLayer_Network_Subnet_Registration (type)|SoftLayer_Network_Subnet_Registration]].
type Network_Subnet_Registration_Ripe struct {
	Network_Subnet_Registration
}

// The subnet registration status type has been deprecated.
//
// Subnet Registration Status objects describe the current status of a subnet registration.
//
// The standard values for these objects are as follows: <ul> <li><strong>OPEN</strong> - Indicates that the registration object is new and has yet to be submitted to the RIR</li> <li><strong>PENDING</strong> - Indicates that the registration object has been submitted to the RIR and is awaiting response</li> <li><strong>COMPLETE</strong> - Indicates that the RIR action has completed</li> <li><strong>DELETED</strong> - Indicates that the registration object has been gracefully removed is no longer valid</li> <li><strong>CANCELLED</strong> - Indicates that the registration object has been abruptly removed is no longer valid</li> </ul>
type Network_Subnet_Registration_Status struct {
	Entity

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// Unique numeric ID of the status object
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Code-friendly string name of the status
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	// Deprecated: This function has been marked as deprecated.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// Human-readable name of the status
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// Every SoftLayer customer account has contact information associated with it for reverse WHOIS purposes. An account's RWHOIS data, modeled by the SoftLayer_Network_Subnet_Rwhois_Data data type, is used by SoftLayer's reverse WHOIS server as well as for SWIP transactions. SoftLayer's reverse WHOIS servers respond to WHOIS queries for IP addresses belonging to a customer's servers, returning this RWHOIS data.
//
// A SoftLayer customer's RWHOIS data may not necessarily match their account or portal users' contact information.
type Network_Subnet_Rwhois_Data struct {
	Entity

	// An email address associated with an account's RWHOIS data that is responsible for responding to network abuse queries about malicious traffic coming from your servers' IP addresses.
	AbuseEmail *string `json:"abuseEmail,omitempty" xmlrpc:"abuseEmail,omitempty"`

	// The SoftLayer customer account associated with this reverse WHOIS data.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// An account's RWHOIS data's associated account identifier.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The first line of the mailing address associated with an account's RWHOIS data.
	Address1 *string `json:"address1,omitempty" xmlrpc:"address1,omitempty"`

	// The second line of the mailing address associated with an account's RWHOIS data.
	Address2 *string `json:"address2,omitempty" xmlrpc:"address2,omitempty"`

	// The city of the mailing address associated with an account's RWHOIS data.
	City *string `json:"city,omitempty" xmlrpc:"city,omitempty"`

	// The company name associated with an account's RWHOIS data.
	CompanyName *string `json:"companyName,omitempty" xmlrpc:"companyName,omitempty"`

	// A two-letter abbreviation of the country of the mailing address associated with an account's RWHOIS data.
	Country *string `json:"country,omitempty" xmlrpc:"country,omitempty"`

	// The date an account's RWHOIS data was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The first name associated with an account's RWHOIS data.
	FirstName *string `json:"firstName,omitempty" xmlrpc:"firstName,omitempty"`

	// An account's RWHOIS data's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The last name associated with an account's RWHOIS data.
	LastName *string `json:"lastName,omitempty" xmlrpc:"lastName,omitempty"`

	// The date an account's RWHOIS data was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The postal code of the mailing address associated with an account's RWHOIS data.
	PostalCode *string `json:"postalCode,omitempty" xmlrpc:"postalCode,omitempty"`

	// Whether an account's RWHOIS data refers to a private residence or not.
	PrivateResidenceFlag *bool `json:"privateResidenceFlag,omitempty" xmlrpc:"privateResidenceFlag,omitempty"`

	// A two-letter abbreviation of the state of the mailing address associated with an account's RWHOIS data. If an account does not reside in a province then this is typically blank.
	State *string `json:"state,omitempty" xmlrpc:"state,omitempty"`
}

// **DEPRECATED**
// The SoftLayer_Network_Subnet_Swip_Transaction data type contains basic information tracked at SoftLayer to allow automation of Swip creation, update, and removal requests.  A specific transaction is attached to an accountId and a subnetId. This also contains a "Status Name" which tells the customer what the transaction is doing:
//
// * REQUEST QUEUED:  Request is queued up to be sent to ARIN
// * REQUEST SENT:  The email request has been sent to ARIN
// * REQUEST CONFIRMED:  ARIN has confirmed that the request is good, and should be available in 24 hours
// * OK:  The subnet has been checked with WHOIS and it the SWIP transaction has completed correctly
// * REMOVE QUEUED:  A subnet is queued to be removed from ARIN's systems
// * REMOVE SENT:  The removal email request has been sent to ARIN
// * REMOVE CONFIRMED:  ARIN has confirmed that the removal request is good, and the subnet should be clear in WHOIS in 24 hours
// * DELETED:  This specific SWIP Transaction has been removed from ARIN and is no longer in effect
// * SOFTLAYER MANUALLY PROCESSING:  Sometimes a request doesn't go through correctly and has to be manually processed by SoftLayer.  This may take some time.
type Network_Subnet_Swip_Transaction struct {
	Entity

	// The Account whose RWHOIS data was used to SWIP this subnet
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// A SWIP transaction's unique identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A Name describing which state a SWIP  transaction is in.
	StatusName *string `json:"statusName,omitempty" xmlrpc:"statusName,omitempty"`

	// The subnet that this SWIP transaction was created for.
	Subnet *Network_Subnet `json:"subnet,omitempty" xmlrpc:"subnet,omitempty"`

	// ID Number of the Subnet for this SWIP transaction.
	SubnetId *int `json:"subnetId,omitempty" xmlrpc:"subnetId,omitempty"`
}

// The SoftLayer_Network_Tunnel_Module_Context data type contains general information relating to a single SoftLayer network tunnel.  The SoftLayer_Network_Tunnel_Module_Context is useful to gather information such as related customer subnets (remote) and internal subnets (local) associated with the network tunnel as well as other information needed to manage the network tunnel.  Account and billing information related to the network tunnel can also be retrieved.
type Network_Tunnel_Module_Context struct {
	Entity

	// The account that a network tunnel belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// A network tunnel's account identifier.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// DEPRECATED
	// Deprecated: This function has been marked as deprecated.
	ActiveTransaction *Provisioning_Version1_Transaction `json:"activeTransaction,omitempty" xmlrpc:"activeTransaction,omitempty"`

	// A count of a network tunnel's address translations.
	AddressTranslationCount *uint `json:"addressTranslationCount,omitempty" xmlrpc:"addressTranslationCount,omitempty"`

	// A network tunnel's address translations.
	AddressTranslations []Network_Tunnel_Module_Context_Address_Translation `json:"addressTranslations,omitempty" xmlrpc:"addressTranslations,omitempty"`

	// A flag used to specify when advanced configurations, complex configurations that require manual setup, are being applied to network devices for a network tunnel. When the flag is set to true (1), a network tunnel cannot be configured through the management portal nor the API.
	AdvancedConfigurationFlag *int `json:"advancedConfigurationFlag,omitempty" xmlrpc:"advancedConfigurationFlag,omitempty"`

	// A count of subnets that provide access to SoftLayer services such as the management portal and the SoftLayer API.
	AllAvailableServiceSubnetCount *uint `json:"allAvailableServiceSubnetCount,omitempty" xmlrpc:"allAvailableServiceSubnetCount,omitempty"`

	// Subnets that provide access to SoftLayer services such as the management portal and the SoftLayer API.
	AllAvailableServiceSubnets []Network_Subnet `json:"allAvailableServiceSubnets,omitempty" xmlrpc:"allAvailableServiceSubnets,omitempty"`

	// The current billing item for network tunnel.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// The date a network tunnel was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The remote end of a network tunnel. This end of the network tunnel resides on an outside network and will be sending and receiving the IPSec packets.
	CustomerPeerIpAddress *string `json:"customerPeerIpAddress,omitempty" xmlrpc:"customerPeerIpAddress,omitempty"`

	// A count of remote subnets that are allowed access through a network tunnel.
	CustomerSubnetCount *uint `json:"customerSubnetCount,omitempty" xmlrpc:"customerSubnetCount,omitempty"`

	// Remote subnets that are allowed access through a network tunnel.
	CustomerSubnets []Network_Customer_Subnet `json:"customerSubnets,omitempty" xmlrpc:"customerSubnets,omitempty"`

	// The datacenter location for one end of the network tunnel that allows access to account's private subnets.
	Datacenter *Location `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// The name giving to a network tunnel by a user.
	FriendlyName *string `json:"friendlyName,omitempty" xmlrpc:"friendlyName,omitempty"`

	// A network tunnel's unique identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The local  end of a network tunnel. This end of the network tunnel resides on the SoftLayer networks and allows access to remote end of the tunnel to subnets on SoftLayer networks.
	InternalPeerIpAddress *string `json:"internalPeerIpAddress,omitempty" xmlrpc:"internalPeerIpAddress,omitempty"`

	// A count of private subnets that can be accessed through the network tunnel.
	InternalSubnetCount *uint `json:"internalSubnetCount,omitempty" xmlrpc:"internalSubnetCount,omitempty"`

	// Private subnets that can be accessed through the network tunnel.
	InternalSubnets []Network_Subnet `json:"internalSubnets,omitempty" xmlrpc:"internalSubnets,omitempty"`

	// The date a network tunnel was last modified.
	//
	// NOTE:  This date should NOT be used to determine when the network tunnel configurations were last applied to the network device.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// A network tunnel's unique name used on the network device.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// Authentication used to generate keys for protecting the negotiations for a network tunnel.
	PhaseOneAuthentication *string `json:"phaseOneAuthentication,omitempty" xmlrpc:"phaseOneAuthentication,omitempty"`

	// Determines the strength of the key used in the key exchange process.  The higher the group number the stronger the key is and the more secure it is.  However, processing time will increase as the strength of the key increases.  Both peers in the must use the Diffie-Hellman Group.
	PhaseOneDiffieHellmanGroup *int `json:"phaseOneDiffieHellmanGroup,omitempty" xmlrpc:"phaseOneDiffieHellmanGroup,omitempty"`

	// Encryption used to generate keys for protecting the negotiations for a network tunnel.
	PhaseOneEncryption *string `json:"phaseOneEncryption,omitempty" xmlrpc:"phaseOneEncryption,omitempty"`

	// Amount of time (in seconds) allowed to pass before the encryption key expires.  A new key is generated without interrupting service. Valid times are from 120 to 172800 seconds.
	PhaseOneKeylife *int `json:"phaseOneKeylife,omitempty" xmlrpc:"phaseOneKeylife,omitempty"`

	// The authentication used in phase 2 proposal negotiation process.
	PhaseTwoAuthentication *string `json:"phaseTwoAuthentication,omitempty" xmlrpc:"phaseTwoAuthentication,omitempty"`

	// Determines the strength of the key used in the key exchange process.  The higher the group number the stronger the key is and the more secure it is.  However, processing time will increase as the strength of the key increases.  Both peers must use the Diffie-Hellman Group.
	PhaseTwoDiffieHellmanGroup *int `json:"phaseTwoDiffieHellmanGroup,omitempty" xmlrpc:"phaseTwoDiffieHellmanGroup,omitempty"`

	// The encryption used in phase 2 proposal negotiation process.
	PhaseTwoEncryption *string `json:"phaseTwoEncryption,omitempty" xmlrpc:"phaseTwoEncryption,omitempty"`

	// Amount of time (in seconds) allowed to pass before the encryption key expires.  A new key is generated without interrupting service. Valid times are from 120 to 172800 seconds.
	PhaseTwoKeylife *int `json:"phaseTwoKeylife,omitempty" xmlrpc:"phaseTwoKeylife,omitempty"`

	// Determines if the generated keys are made from previous keys.  When PFS is specified, a Diffie-Hellman exchange occurs each time a new security association is negotiated.
	PhaseTwoPerfectForwardSecrecy *int `json:"phaseTwoPerfectForwardSecrecy,omitempty" xmlrpc:"phaseTwoPerfectForwardSecrecy,omitempty"`

	// A key used so that peers authenticate each other.  This key is hashed by using the phase one encryption and phase one authentication.
	PresharedKey *string `json:"presharedKey,omitempty" xmlrpc:"presharedKey,omitempty"`

	// A count of service subnets that can be access through the network tunnel.
	ServiceSubnetCount *uint `json:"serviceSubnetCount,omitempty" xmlrpc:"serviceSubnetCount,omitempty"`

	// Service subnets that can be access through the network tunnel.
	ServiceSubnets []Network_Subnet `json:"serviceSubnets,omitempty" xmlrpc:"serviceSubnets,omitempty"`

	// A count of subnets used for a network tunnel's address translations.
	StaticRouteSubnetCount *uint `json:"staticRouteSubnetCount,omitempty" xmlrpc:"staticRouteSubnetCount,omitempty"`

	// Subnets used for a network tunnel's address translations.
	StaticRouteSubnets []Network_Subnet `json:"staticRouteSubnets,omitempty" xmlrpc:"staticRouteSubnets,omitempty"`

	// DEPRECATED
	// Deprecated: This function has been marked as deprecated.
	TransactionHistory []Provisioning_Version1_Transaction `json:"transactionHistory,omitempty" xmlrpc:"transactionHistory,omitempty"`

	// A count of dEPRECATED
	TransactionHistoryCount *uint `json:"transactionHistoryCount,omitempty" xmlrpc:"transactionHistoryCount,omitempty"`
}

// The SoftLayer_Network_Tunnel_Module_Context_Address_Translation data type contains general information relating to a single address translation. Information such as notes, ip addresses, along with record information, and network tunnel data may be retrieved.
type Network_Tunnel_Module_Context_Address_Translation struct {
	Entity

	// The ip address record that will receive the encrypted traffic.
	CustomerIpAddress *string `json:"customerIpAddress,omitempty" xmlrpc:"customerIpAddress,omitempty"`

	// The unique identifier for the ip address record that will receive the encrypted traffic.
	CustomerIpAddressId *int `json:"customerIpAddressId,omitempty" xmlrpc:"customerIpAddressId,omitempty"`

	// The ip address record for the ip that will receive the encrypted traffic from the IPSec network tunnel.
	CustomerIpAddressRecord *Network_Customer_Subnet_IpAddress `json:"customerIpAddressRecord,omitempty" xmlrpc:"customerIpAddressRecord,omitempty"`

	// An address translation's unique identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The ip address record that will deliver the encrypted traffic.
	InternalIpAddress *string `json:"internalIpAddress,omitempty" xmlrpc:"internalIpAddress,omitempty"`

	// The unique identifier for the ip address record that will deliver the encrypted traffic.
	InternalIpAddressId *int `json:"internalIpAddressId,omitempty" xmlrpc:"internalIpAddressId,omitempty"`

	// The ip address record for the ip that will deliver the encrypted traffic from the IPSec network tunnel.
	InternalIpAddressRecord *Network_Subnet_IpAddress `json:"internalIpAddressRecord,omitempty" xmlrpc:"internalIpAddressRecord,omitempty"`

	// The IPSec network tunnel an address translation belongs to.
	NetworkTunnelContext *Network_Tunnel_Module_Context `json:"networkTunnelContext,omitempty" xmlrpc:"networkTunnelContext,omitempty"`

	// An address translation's network tunnel identifier.
	NetworkTunnelContextId *int `json:"networkTunnelContextId,omitempty" xmlrpc:"networkTunnelContextId,omitempty"`

	// A name or description given to an address translation to help identify the address translation.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`
}

// VLANs comprise the fundamental segmentation model on the network, isolating customer networks from one another.
//
// VLANs are scoped to a single network, generally public or private, and a pod. Through association to a single VLAN, assigned subnets are routed on the network to provide IP address connectivity.
//
// Compute devices are associated to a single VLAN per active network, to which the Primary IP Address and containing Primary Subnet belongs. Additional VLANs may be associated to bare metal devices using VLAN trunking.
//
// [VLAN at Wikipedia](https://en.wikipedia.org/wiki/VLAN)
type Network_Vlan struct {
	Entity

	// The account this VLAN is associated with.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The identifier of the account this VLAN is assigned to.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// A count of the primary IPv4 subnets routed on this VLAN, excluding the primarySubnet.
	AdditionalPrimarySubnetCount *uint `json:"additionalPrimarySubnetCount,omitempty" xmlrpc:"additionalPrimarySubnetCount,omitempty"`

	// The primary IPv4 subnets routed on this VLAN, excluding the primarySubnet.
	AdditionalPrimarySubnets []Network_Subnet `json:"additionalPrimarySubnets,omitempty" xmlrpc:"additionalPrimarySubnets,omitempty"`

	// The gateway device this VLAN is associated with for routing purposes.
	AttachedNetworkGateway *Network_Gateway `json:"attachedNetworkGateway,omitempty" xmlrpc:"attachedNetworkGateway,omitempty"`

	// A value of '1' indicates this VLAN is associated with a gateway device for routing purposes.
	AttachedNetworkGatewayFlag *bool `json:"attachedNetworkGatewayFlag,omitempty" xmlrpc:"attachedNetworkGatewayFlag,omitempty"`

	// The gateway device VLAN context this VLAN is associated with for routing purposes.
	AttachedNetworkGatewayVlan *Network_Gateway_Vlan `json:"attachedNetworkGatewayVlan,omitempty" xmlrpc:"attachedNetworkGatewayVlan,omitempty"`

	// The billing item for this VLAN.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// The datacenter this VLAN is associated with.
	Datacenter *Location `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// A value of '1' indicates this VLAN is associated with a firewall device. This does not include Hardware Firewalls.
	DedicatedFirewallFlag *int `json:"dedicatedFirewallFlag,omitempty" xmlrpc:"dedicatedFirewallFlag,omitempty"`

	// [DEPRECATED] The extension router that this VLAN is associated with.
	// Deprecated: This function has been marked as deprecated.
	ExtensionRouter *Hardware_Router `json:"extensionRouter,omitempty" xmlrpc:"extensionRouter,omitempty"`

	// A count of the VSI network interfaces connected to this VLAN and associated with a Hardware Firewall.
	FirewallGuestNetworkComponentCount *uint `json:"firewallGuestNetworkComponentCount,omitempty" xmlrpc:"firewallGuestNetworkComponentCount,omitempty"`

	// The VSI network interfaces connected to this VLAN and associated with a Hardware Firewall.
	FirewallGuestNetworkComponents []Network_Component_Firewall `json:"firewallGuestNetworkComponents,omitempty" xmlrpc:"firewallGuestNetworkComponents,omitempty"`

	// A count of the context for the firewall device associated with this VLAN.
	FirewallInterfaceCount *uint `json:"firewallInterfaceCount,omitempty" xmlrpc:"firewallInterfaceCount,omitempty"`

	// The context for the firewall device associated with this VLAN.
	FirewallInterfaces []Network_Firewall_Module_Context_Interface `json:"firewallInterfaces,omitempty" xmlrpc:"firewallInterfaces,omitempty"`

	// A count of the uplinks of the hardware network interfaces connected natively to this VLAN and associated with a Hardware Firewall.
	FirewallNetworkComponentCount *uint `json:"firewallNetworkComponentCount,omitempty" xmlrpc:"firewallNetworkComponentCount,omitempty"`

	// The uplinks of the hardware network interfaces connected natively to this VLAN and associated with a Hardware Firewall.
	FirewallNetworkComponents []Network_Component_Firewall `json:"firewallNetworkComponents,omitempty" xmlrpc:"firewallNetworkComponents,omitempty"`

	// A count of the access rules for the firewall device associated with this VLAN.
	FirewallRuleCount *uint `json:"firewallRuleCount,omitempty" xmlrpc:"firewallRuleCount,omitempty"`

	// The access rules for the firewall device associated with this VLAN.
	FirewallRules []Network_Vlan_Firewall_Rule `json:"firewallRules,omitempty" xmlrpc:"firewallRules,omitempty"`

	// A human readable, unique identifier for this VLAN.
	FullyQualifiedName *string `json:"fullyQualifiedName,omitempty" xmlrpc:"fullyQualifiedName,omitempty"`

	// A count of the VSI network interfaces connected to this VLAN.
	GuestNetworkComponentCount *uint `json:"guestNetworkComponentCount,omitempty" xmlrpc:"guestNetworkComponentCount,omitempty"`

	// The VSI network interfaces connected to this VLAN.
	GuestNetworkComponents []Virtual_Guest_Network_Component `json:"guestNetworkComponents,omitempty" xmlrpc:"guestNetworkComponents,omitempty"`

	// The hardware with network interfaces connected natively to this VLAN.
	Hardware []Hardware `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// A count of the hardware with network interfaces connected natively to this VLAN.
	HardwareCount *uint `json:"hardwareCount,omitempty" xmlrpc:"hardwareCount,omitempty"`

	// A value of '1' indicates this VLAN is associated with a firewall device in a high availability configuration.
	HighAvailabilityFirewallFlag *bool `json:"highAvailabilityFirewallFlag,omitempty" xmlrpc:"highAvailabilityFirewallFlag,omitempty"`

	// The unique identifier of this VLAN.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A value of '1' indicates this VLAN's pod has VSI local disk storage capability.
	LocalDiskStorageCapabilityFlag *bool `json:"localDiskStorageCapabilityFlag,omitempty" xmlrpc:"localDiskStorageCapabilityFlag,omitempty"`

	// The time this VLAN was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The customer name for this VLAN.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A count of the hardware network interfaces connected natively to this VLAN.
	NetworkComponentCount *uint `json:"networkComponentCount,omitempty" xmlrpc:"networkComponentCount,omitempty"`

	// A count of the hardware network interfaces connected via trunk to this VLAN.
	NetworkComponentTrunkCount *uint `json:"networkComponentTrunkCount,omitempty" xmlrpc:"networkComponentTrunkCount,omitempty"`

	// The hardware network interfaces connected via trunk to this VLAN.
	NetworkComponentTrunks []Network_Component_Network_Vlan_Trunk `json:"networkComponentTrunks,omitempty" xmlrpc:"networkComponentTrunks,omitempty"`

	// The hardware network interfaces connected natively to this VLAN.
	NetworkComponents []Network_Component `json:"networkComponents,omitempty" xmlrpc:"networkComponents,omitempty"`

	// The viable hardware network interface trunking targets of this VLAN. Viable targets include accessible components of assigned hardware in the same pod and network as this VLAN, which are not already connected, either natively or trunked.
	NetworkComponentsTrunkable []Network_Component `json:"networkComponentsTrunkable,omitempty" xmlrpc:"networkComponentsTrunkable,omitempty"`

	// The network that this VLAN is on, either PUBLIC or PRIVATE, if applicable.
	NetworkSpace *string `json:"networkSpace,omitempty" xmlrpc:"networkSpace,omitempty"`

	// The firewall device associated with this VLAN.
	NetworkVlanFirewall *Network_Vlan_Firewall `json:"networkVlanFirewall,omitempty" xmlrpc:"networkVlanFirewall,omitempty"`

	// An internal description of this VLAN, if applicable.
	Note *string `json:"note,omitempty" xmlrpc:"note,omitempty"`

	// The pod this VLAN is associated with.
	PodName *string `json:"podName,omitempty" xmlrpc:"podName,omitempty"`

	// The router device that this VLAN is associated with.
	PrimaryRouter *Hardware_Router `json:"primaryRouter,omitempty" xmlrpc:"primaryRouter,omitempty"`

	// A primary IPv4 subnet routed on this VLAN, if accessible.
	PrimarySubnet *Network_Subnet `json:"primarySubnet,omitempty" xmlrpc:"primarySubnet,omitempty"`

	// A count of all primary subnets routed on this VLAN.
	PrimarySubnetCount *uint `json:"primarySubnetCount,omitempty" xmlrpc:"primarySubnetCount,omitempty"`

	// The identifier of a primary IPv4 subnet routed on this VLAN.
	PrimarySubnetId *int `json:"primarySubnetId,omitempty" xmlrpc:"primarySubnetId,omitempty"`

	// The primary IPv6 subnet routed on this VLAN, if IPv6 is enabled.
	PrimarySubnetVersion6 *Network_Subnet `json:"primarySubnetVersion6,omitempty" xmlrpc:"primarySubnetVersion6,omitempty"`

	// All primary subnets routed on this VLAN.
	PrimarySubnets []Network_Subnet `json:"primarySubnets,omitempty" xmlrpc:"primarySubnets,omitempty"`

	// A count of the gateway devices with connectivity supported by this private VLAN.
	PrivateNetworkGatewayCount *uint `json:"privateNetworkGatewayCount,omitempty" xmlrpc:"privateNetworkGatewayCount,omitempty"`

	// The gateway devices with connectivity supported by this private VLAN.
	PrivateNetworkGateways []Network_Gateway `json:"privateNetworkGateways,omitempty" xmlrpc:"privateNetworkGateways,omitempty"`

	// A count of iP addresses routed on this VLAN which are actively associated with network protections.
	ProtectedIpAddressCount *uint `json:"protectedIpAddressCount,omitempty" xmlrpc:"protectedIpAddressCount,omitempty"`

	// IP addresses routed on this VLAN which are actively associated with network protections.
	ProtectedIpAddresses []Network_Subnet_IpAddress `json:"protectedIpAddresses,omitempty" xmlrpc:"protectedIpAddresses,omitempty"`

	// A count of the gateway devices with connectivity supported by this public VLAN.
	PublicNetworkGatewayCount *uint `json:"publicNetworkGatewayCount,omitempty" xmlrpc:"publicNetworkGatewayCount,omitempty"`

	// The gateway devices with connectivity supported by this public VLAN.
	PublicNetworkGateways []Network_Gateway `json:"publicNetworkGateways,omitempty" xmlrpc:"publicNetworkGateways,omitempty"`

	// A value of '1' indicates this VLAN's pod has VSI SAN disk storage capability.
	SanStorageCapabilityFlag *bool `json:"sanStorageCapabilityFlag,omitempty" xmlrpc:"sanStorageCapabilityFlag,omitempty"`

	// [DEPRECATED] The secondary router device that this VLAN is associated with.
	// Deprecated: This function has been marked as deprecated.
	SecondaryRouter *Hardware `json:"secondaryRouter,omitempty" xmlrpc:"secondaryRouter,omitempty"`

	// A count of all non-primary subnets routed on this VLAN.
	SecondarySubnetCount *uint `json:"secondarySubnetCount,omitempty" xmlrpc:"secondarySubnetCount,omitempty"`

	// All non-primary subnets routed on this VLAN.
	SecondarySubnets []Network_Subnet `json:"secondarySubnets,omitempty" xmlrpc:"secondarySubnets,omitempty"`

	// A count of all subnets routed on this VLAN.
	SubnetCount *uint `json:"subnetCount,omitempty" xmlrpc:"subnetCount,omitempty"`

	// All subnets routed on this VLAN.
	Subnets []Network_Subnet `json:"subnets,omitempty" xmlrpc:"subnets,omitempty"`

	// A count of the tags associated to this VLAN.
	TagReferenceCount *uint `json:"tagReferenceCount,omitempty" xmlrpc:"tagReferenceCount,omitempty"`

	// The tags associated to this VLAN.
	TagReferences []Tag_Reference `json:"tagReferences,omitempty" xmlrpc:"tagReferences,omitempty"`

	// The number of primary IPv4 addresses routed on this VLAN.
	TotalPrimaryIpAddressCount *uint `json:"totalPrimaryIpAddressCount,omitempty" xmlrpc:"totalPrimaryIpAddressCount,omitempty"`

	// The type for this VLAN, with the following values: STANDARD, GATEWAY, INTERCONNECT
	Type *Network_Vlan_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// A count of the VSIs with network interfaces connected to this VLAN.
	VirtualGuestCount *uint `json:"virtualGuestCount,omitempty" xmlrpc:"virtualGuestCount,omitempty"`

	// The VSIs with network interfaces connected to this VLAN.
	VirtualGuests []Virtual_Guest `json:"virtualGuests,omitempty" xmlrpc:"virtualGuests,omitempty"`

	// The number of this VLAN configured on the network.
	VlanNumber *int `json:"vlanNumber,omitempty" xmlrpc:"vlanNumber,omitempty"`
}

// The SoftLayer_Network_Vlan_Firewall data type contains general information relating to a single SoftLayer VLAN firewall. This is the object which ties the running rules to a specific downstream server. Use the [[SoftLayer Network Firewall Template]] service to pull SoftLayer recommended rule set templates. Use the [[SoftLayer Network Firewall Update Request]] service to submit a firewall update request.
type Network_Vlan_Firewall struct {
	Entity

	// no documentation yet
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// A flag to indicate if the firewall is in administrative bypass mode. In other words, no rules are being applied to the traffic coming through.
	AdministrativeBypassFlag *string `json:"administrativeBypassFlag,omitempty" xmlrpc:"administrativeBypassFlag,omitempty"`

	// A firewall's allotted bandwidth (measured in GB).
	BandwidthAllocation *Float64 `json:"bandwidthAllocation,omitempty" xmlrpc:"bandwidthAllocation,omitempty"`

	// The raw bandwidth usage data for the current billing cycle. One object will be returned for each network this firewall is attached to.
	BillingCycleBandwidthUsage []Network_Bandwidth_Usage `json:"billingCycleBandwidthUsage,omitempty" xmlrpc:"billingCycleBandwidthUsage,omitempty"`

	// A count of the raw bandwidth usage data for the current billing cycle. One object will be returned for each network this firewall is attached to.
	BillingCycleBandwidthUsageCount *uint `json:"billingCycleBandwidthUsageCount,omitempty" xmlrpc:"billingCycleBandwidthUsageCount,omitempty"`

	// The raw private bandwidth usage data for the current billing cycle.
	BillingCyclePrivateBandwidthUsage *Network_Bandwidth_Usage `json:"billingCyclePrivateBandwidthUsage,omitempty" xmlrpc:"billingCyclePrivateBandwidthUsage,omitempty"`

	// The raw public bandwidth usage data for the current billing cycle.
	BillingCyclePublicBandwidthUsage *Network_Bandwidth_Usage `json:"billingCyclePublicBandwidthUsage,omitempty" xmlrpc:"billingCyclePublicBandwidthUsage,omitempty"`

	// The billing item for a Hardware Firewall (Dedicated).
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// Administrative bypass request status.
	BypassRequestStatus *string `json:"bypassRequestStatus,omitempty" xmlrpc:"bypassRequestStatus,omitempty"`

	// Whether or not this firewall can be directly logged in to.
	CustomerManagedFlag *bool `json:"customerManagedFlag,omitempty" xmlrpc:"customerManagedFlag,omitempty"`

	// The datacenter that the firewall resides in.
	Datacenter *Location `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// The firewall device type.
	FirewallType *string `json:"firewallType,omitempty" xmlrpc:"firewallType,omitempty"`

	// A name reflecting the hostname and domain of the firewall. This is created from the combined values of the firewall's logical name and vlan number automatically, and thus can not be edited directly.
	FullyQualifiedDomainName *string `json:"fullyQualifiedDomainName,omitempty" xmlrpc:"fullyQualifiedDomainName,omitempty"`

	// A firewall's unique identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The credentials to log in to a firewall device. This is only present for dedicated appliances.
	ManagementCredentials *Software_Component_Password `json:"managementCredentials,omitempty" xmlrpc:"managementCredentials,omitempty"`

	// A firewall's metric tracking object.
	MetricTrackingObject *Metric_Tracking_Object `json:"metricTrackingObject,omitempty" xmlrpc:"metricTrackingObject,omitempty"`

	// The metric tracking object ID for this firewall.
	MetricTrackingObjectId *int `json:"metricTrackingObjectId,omitempty" xmlrpc:"metricTrackingObjectId,omitempty"`

	// A count of the update requests made for this firewall.
	NetworkFirewallUpdateRequestCount *uint `json:"networkFirewallUpdateRequestCount,omitempty" xmlrpc:"networkFirewallUpdateRequestCount,omitempty"`

	// The update requests made for this firewall.
	NetworkFirewallUpdateRequests []Network_Firewall_Update_Request `json:"networkFirewallUpdateRequests,omitempty" xmlrpc:"networkFirewallUpdateRequests,omitempty"`

	// The gateway associated with this firewall, if any.
	NetworkGateway *Network_Gateway `json:"networkGateway,omitempty" xmlrpc:"networkGateway,omitempty"`

	// The VLAN object that a firewall is associated with and protecting.
	NetworkVlan *Network_Vlan `json:"networkVlan,omitempty" xmlrpc:"networkVlan,omitempty"`

	// A count of the VLAN objects that a firewall is associated with and protecting.
	NetworkVlanCount *uint `json:"networkVlanCount,omitempty" xmlrpc:"networkVlanCount,omitempty"`

	// The VLAN objects that a firewall is associated with and protecting.
	NetworkVlans []Network_Vlan `json:"networkVlans,omitempty" xmlrpc:"networkVlans,omitempty"`

	// A firewall's primary IP address. This field will be the IP shown when doing network traces and reverse DNS and is a read-only property.
	PrimaryIpAddress *string `json:"primaryIpAddress,omitempty" xmlrpc:"primaryIpAddress,omitempty"`

	// A count of the currently running rule set of this network component firewall.
	RuleCount *uint `json:"ruleCount,omitempty" xmlrpc:"ruleCount,omitempty"`

	// The currently running rule set of this network component firewall.
	Rules []Network_Vlan_Firewall_Rule `json:"rules,omitempty" xmlrpc:"rules,omitempty"`

	// A count of
	TagReferenceCount *uint `json:"tagReferenceCount,omitempty" xmlrpc:"tagReferenceCount,omitempty"`

	// no documentation yet
	TagReferences []Tag_Reference `json:"tagReferences,omitempty" xmlrpc:"tagReferences,omitempty"`

	// A firewall's associated upgrade request object, if any.
	UpgradeRequest *Product_Upgrade_Request `json:"upgradeRequest,omitempty" xmlrpc:"upgradeRequest,omitempty"`
}

// A SoftLayer_Network_Component_Firewall_Rule object type represents a currently running firewall rule and contains relative information. Use the [[SoftLayer Network Firewall Update Request]] service to submit a firewall update request. Use the [[SoftLayer Network Firewall Template]] service to pull SoftLayer recommended rule set templates.
type Network_Vlan_Firewall_Rule struct {
	Entity

	// The action that the rule is to take [permit or deny].
	Action *string `json:"action,omitempty" xmlrpc:"action,omitempty"`

	// The destination IP address considered for determining rule application.
	DestinationIpAddress *string `json:"destinationIpAddress,omitempty" xmlrpc:"destinationIpAddress,omitempty"`

	// The CIDR is used for determining rule application. This value will
	DestinationIpCidr *int `json:"destinationIpCidr,omitempty" xmlrpc:"destinationIpCidr,omitempty"`

	// The destination IP subnet mask considered for determining rule application.
	DestinationIpSubnetMask *string `json:"destinationIpSubnetMask,omitempty" xmlrpc:"destinationIpSubnetMask,omitempty"`

	// The ending (upper end of range) destination port considered for determining rule application.
	DestinationPortRangeEnd *int `json:"destinationPortRangeEnd,omitempty" xmlrpc:"destinationPortRangeEnd,omitempty"`

	// The starting (lower end of range) destination port considered for determining rule application.
	DestinationPortRangeStart *int `json:"destinationPortRangeStart,omitempty" xmlrpc:"destinationPortRangeStart,omitempty"`

	// The rule's internal identifier.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The network component firewall that this rule belongs to.
	NetworkComponentFirewall *Network_Component_Firewall `json:"networkComponentFirewall,omitempty" xmlrpc:"networkComponentFirewall,omitempty"`

	// The notes field for the rule.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// The numeric value describing the order in which the rule should be applied.
	OrderValue *int `json:"orderValue,omitempty" xmlrpc:"orderValue,omitempty"`

	// The protocol considered for determining rule application.
	Protocol *string `json:"protocol,omitempty" xmlrpc:"protocol,omitempty"`

	// The source IP address considered for determining rule application.
	SourceIpAddress *string `json:"sourceIpAddress,omitempty" xmlrpc:"sourceIpAddress,omitempty"`

	// The CIDR is used for determining rule application. This value will
	SourceIpCidr *int `json:"sourceIpCidr,omitempty" xmlrpc:"sourceIpCidr,omitempty"`

	// The source IP subnet mask considered for determining rule application.
	SourceIpSubnetMask *string `json:"sourceIpSubnetMask,omitempty" xmlrpc:"sourceIpSubnetMask,omitempty"`

	// Current status of the network component firewall.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// Whether this rule is an IPv4 rule or an IPv6 rule. If
	Version *int `json:"version,omitempty" xmlrpc:"version,omitempty"`
}

// no documentation yet
type Network_Vlan_Type struct {
	Entity

	// no documentation yet
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}
