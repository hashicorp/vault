/*
Copyright (c) 2014-2024 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mo

import (
	"reflect"
	"time"

	"github.com/vmware/govmomi/vim25/types"
)

type Alarm struct {
	ExtensibleManagedObject

	Info types.AlarmInfo `json:"info"`
}

func init() {
	t["Alarm"] = reflect.TypeOf((*Alarm)(nil)).Elem()
}

type AlarmManager struct {
	Self types.ManagedObjectReference `json:"self"`

	DefaultExpression []types.BaseAlarmExpression `json:"defaultExpression"`
	Description       types.AlarmDescription      `json:"description"`
}

func (m AlarmManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["AlarmManager"] = reflect.TypeOf((*AlarmManager)(nil)).Elem()
}

type AuthorizationManager struct {
	Self types.ManagedObjectReference `json:"self"`

	PrivilegeList []types.AuthorizationPrivilege `json:"privilegeList"`
	RoleList      []types.AuthorizationRole      `json:"roleList"`
	Description   types.AuthorizationDescription `json:"description"`
}

func (m AuthorizationManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["AuthorizationManager"] = reflect.TypeOf((*AuthorizationManager)(nil)).Elem()
}

type CertificateManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m CertificateManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["CertificateManager"] = reflect.TypeOf((*CertificateManager)(nil)).Elem()
}

type ClusterComputeResource struct {
	ComputeResource

	Configuration     types.ClusterConfigInfo                    `json:"configuration"`
	Recommendation    []types.ClusterRecommendation              `json:"recommendation"`
	DrsRecommendation []types.ClusterDrsRecommendation           `json:"drsRecommendation"`
	HciConfig         *types.ClusterComputeResourceHCIConfigInfo `json:"hciConfig"`
	MigrationHistory  []types.ClusterDrsMigration                `json:"migrationHistory"`
	ActionHistory     []types.ClusterActionHistory               `json:"actionHistory"`
	DrsFault          []types.ClusterDrsFaults                   `json:"drsFault"`
}

func init() {
	t["ClusterComputeResource"] = reflect.TypeOf((*ClusterComputeResource)(nil)).Elem()
}

type ClusterEVCManager struct {
	ExtensibleManagedObject

	ManagedCluster types.ManagedObjectReference    `json:"managedCluster"`
	EvcState       types.ClusterEVCManagerEVCState `json:"evcState"`
}

func init() {
	t["ClusterEVCManager"] = reflect.TypeOf((*ClusterEVCManager)(nil)).Elem()
}

type ClusterProfile struct {
	Profile
}

func init() {
	t["ClusterProfile"] = reflect.TypeOf((*ClusterProfile)(nil)).Elem()
}

type ClusterProfileManager struct {
	ProfileManager
}

func init() {
	t["ClusterProfileManager"] = reflect.TypeOf((*ClusterProfileManager)(nil)).Elem()
}

type ComputeResource struct {
	ManagedEntity

	ResourcePool       *types.ManagedObjectReference       `json:"resourcePool"`
	Host               []types.ManagedObjectReference      `json:"host"`
	Datastore          []types.ManagedObjectReference      `json:"datastore"`
	Network            []types.ManagedObjectReference      `json:"network"`
	Summary            types.BaseComputeResourceSummary    `json:"summary"`
	EnvironmentBrowser *types.ManagedObjectReference       `json:"environmentBrowser"`
	ConfigurationEx    types.BaseComputeResourceConfigInfo `json:"configurationEx"`
	LifecycleManaged   *bool                               `json:"lifecycleManaged"`
}

func (m *ComputeResource) Entity() *ManagedEntity {
	return &m.ManagedEntity
}

func init() {
	t["ComputeResource"] = reflect.TypeOf((*ComputeResource)(nil)).Elem()
}

type ContainerView struct {
	ManagedObjectView

	Container types.ManagedObjectReference `json:"container"`
	Type      []string                     `json:"type"`
	Recursive bool                         `json:"recursive"`
}

func init() {
	t["ContainerView"] = reflect.TypeOf((*ContainerView)(nil)).Elem()
}

type CryptoManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Enabled bool `json:"enabled"`
}

func (m CryptoManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["CryptoManager"] = reflect.TypeOf((*CryptoManager)(nil)).Elem()
}

type CryptoManagerHost struct {
	CryptoManager
}

func init() {
	t["CryptoManagerHost"] = reflect.TypeOf((*CryptoManagerHost)(nil)).Elem()
}

type CryptoManagerHostKMS struct {
	CryptoManagerHost
}

func init() {
	t["CryptoManagerHostKMS"] = reflect.TypeOf((*CryptoManagerHostKMS)(nil)).Elem()
}

type CryptoManagerKmip struct {
	CryptoManager

	KmipServers []types.KmipClusterInfo `json:"kmipServers"`
}

func init() {
	t["CryptoManagerKmip"] = reflect.TypeOf((*CryptoManagerKmip)(nil)).Elem()
}

type CustomFieldsManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Field []types.CustomFieldDef `json:"field"`
}

func (m CustomFieldsManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["CustomFieldsManager"] = reflect.TypeOf((*CustomFieldsManager)(nil)).Elem()
}

type CustomizationSpecManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Info          []types.CustomizationSpecInfo `json:"info"`
	EncryptionKey types.ByteSlice               `json:"encryptionKey"`
}

func (m CustomizationSpecManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["CustomizationSpecManager"] = reflect.TypeOf((*CustomizationSpecManager)(nil)).Elem()
}

type Datacenter struct {
	ManagedEntity

	VmFolder        types.ManagedObjectReference   `json:"vmFolder"`
	HostFolder      types.ManagedObjectReference   `json:"hostFolder"`
	DatastoreFolder types.ManagedObjectReference   `json:"datastoreFolder"`
	NetworkFolder   types.ManagedObjectReference   `json:"networkFolder"`
	Datastore       []types.ManagedObjectReference `json:"datastore"`
	Network         []types.ManagedObjectReference `json:"network"`
	Configuration   types.DatacenterConfigInfo     `json:"configuration"`
}

func (m *Datacenter) Entity() *ManagedEntity {
	return &m.ManagedEntity
}

func init() {
	t["Datacenter"] = reflect.TypeOf((*Datacenter)(nil)).Elem()
}

type Datastore struct {
	ManagedEntity

	Info              types.BaseDatastoreInfo        `json:"info"`
	Summary           types.DatastoreSummary         `json:"summary"`
	Host              []types.DatastoreHostMount     `json:"host"`
	Vm                []types.ManagedObjectReference `json:"vm"`
	Browser           types.ManagedObjectReference   `json:"browser"`
	Capability        types.DatastoreCapability      `json:"capability"`
	IormConfiguration *types.StorageIORMInfo         `json:"iormConfiguration"`
}

func (m *Datastore) Entity() *ManagedEntity {
	return &m.ManagedEntity
}

func init() {
	t["Datastore"] = reflect.TypeOf((*Datastore)(nil)).Elem()
}

type DatastoreNamespaceManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m DatastoreNamespaceManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["DatastoreNamespaceManager"] = reflect.TypeOf((*DatastoreNamespaceManager)(nil)).Elem()
}

type DiagnosticManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m DiagnosticManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["DiagnosticManager"] = reflect.TypeOf((*DiagnosticManager)(nil)).Elem()
}

type DistributedVirtualPortgroup struct {
	Network

	Key      string                      `json:"key"`
	Config   types.DVPortgroupConfigInfo `json:"config"`
	PortKeys []string                    `json:"portKeys"`
}

func init() {
	t["DistributedVirtualPortgroup"] = reflect.TypeOf((*DistributedVirtualPortgroup)(nil)).Elem()
}

type DistributedVirtualSwitch struct {
	ManagedEntity

	Uuid                string                         `json:"uuid"`
	Capability          types.DVSCapability            `json:"capability"`
	Summary             types.DVSSummary               `json:"summary"`
	Config              types.BaseDVSConfigInfo        `json:"config"`
	NetworkResourcePool []types.DVSNetworkResourcePool `json:"networkResourcePool"`
	Portgroup           []types.ManagedObjectReference `json:"portgroup"`
	Runtime             *types.DVSRuntimeInfo          `json:"runtime"`
}

func (m *DistributedVirtualSwitch) Entity() *ManagedEntity {
	return &m.ManagedEntity
}

func init() {
	t["DistributedVirtualSwitch"] = reflect.TypeOf((*DistributedVirtualSwitch)(nil)).Elem()
}

type DistributedVirtualSwitchManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m DistributedVirtualSwitchManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["DistributedVirtualSwitchManager"] = reflect.TypeOf((*DistributedVirtualSwitchManager)(nil)).Elem()
}

type EnvironmentBrowser struct {
	Self types.ManagedObjectReference `json:"self"`

	DatastoreBrowser *types.ManagedObjectReference `json:"datastoreBrowser"`
}

func (m EnvironmentBrowser) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["EnvironmentBrowser"] = reflect.TypeOf((*EnvironmentBrowser)(nil)).Elem()
}

type EventHistoryCollector struct {
	HistoryCollector

	LatestPage []types.BaseEvent `json:"latestPage"`
}

func init() {
	t["EventHistoryCollector"] = reflect.TypeOf((*EventHistoryCollector)(nil)).Elem()
}

type EventManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Description  types.EventDescription `json:"description"`
	LatestEvent  types.BaseEvent        `json:"latestEvent"`
	MaxCollector int32                  `json:"maxCollector"`
}

func (m EventManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["EventManager"] = reflect.TypeOf((*EventManager)(nil)).Elem()
}

type ExtensibleManagedObject struct {
	Self types.ManagedObjectReference `json:"self"`

	Value          []types.BaseCustomFieldValue `json:"value"`
	AvailableField []types.CustomFieldDef       `json:"availableField"`
}

func (m ExtensibleManagedObject) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["ExtensibleManagedObject"] = reflect.TypeOf((*ExtensibleManagedObject)(nil)).Elem()
}

type ExtensionManager struct {
	Self types.ManagedObjectReference `json:"self"`

	ExtensionList []types.Extension `json:"extensionList"`
}

func (m ExtensionManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["ExtensionManager"] = reflect.TypeOf((*ExtensionManager)(nil)).Elem()
}

type FailoverClusterConfigurator struct {
	Self types.ManagedObjectReference `json:"self"`

	DisabledConfigureMethod []string `json:"disabledConfigureMethod"`
}

func (m FailoverClusterConfigurator) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["FailoverClusterConfigurator"] = reflect.TypeOf((*FailoverClusterConfigurator)(nil)).Elem()
}

type FailoverClusterManager struct {
	Self types.ManagedObjectReference `json:"self"`

	DisabledClusterMethod []string `json:"disabledClusterMethod"`
}

func (m FailoverClusterManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["FailoverClusterManager"] = reflect.TypeOf((*FailoverClusterManager)(nil)).Elem()
}

type FileManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m FileManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["FileManager"] = reflect.TypeOf((*FileManager)(nil)).Elem()
}

type Folder struct {
	ManagedEntity

	ChildType   []string                       `json:"childType"`
	ChildEntity []types.ManagedObjectReference `json:"childEntity"`
	Namespace   *string                        `json:"namespace"`
}

func (m *Folder) Entity() *ManagedEntity {
	return &m.ManagedEntity
}

func init() {
	t["Folder"] = reflect.TypeOf((*Folder)(nil)).Elem()
}

type GuestAliasManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m GuestAliasManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["GuestAliasManager"] = reflect.TypeOf((*GuestAliasManager)(nil)).Elem()
}

type GuestAuthManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m GuestAuthManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["GuestAuthManager"] = reflect.TypeOf((*GuestAuthManager)(nil)).Elem()
}

type GuestFileManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m GuestFileManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["GuestFileManager"] = reflect.TypeOf((*GuestFileManager)(nil)).Elem()
}

type GuestOperationsManager struct {
	Self types.ManagedObjectReference `json:"self"`

	AuthManager                 *types.ManagedObjectReference `json:"authManager"`
	FileManager                 *types.ManagedObjectReference `json:"fileManager"`
	ProcessManager              *types.ManagedObjectReference `json:"processManager"`
	GuestWindowsRegistryManager *types.ManagedObjectReference `json:"guestWindowsRegistryManager"`
	AliasManager                *types.ManagedObjectReference `json:"aliasManager"`
}

func (m GuestOperationsManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["GuestOperationsManager"] = reflect.TypeOf((*GuestOperationsManager)(nil)).Elem()
}

type GuestProcessManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m GuestProcessManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["GuestProcessManager"] = reflect.TypeOf((*GuestProcessManager)(nil)).Elem()
}

type GuestWindowsRegistryManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m GuestWindowsRegistryManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["GuestWindowsRegistryManager"] = reflect.TypeOf((*GuestWindowsRegistryManager)(nil)).Elem()
}

type HealthUpdateManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m HealthUpdateManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HealthUpdateManager"] = reflect.TypeOf((*HealthUpdateManager)(nil)).Elem()
}

type HistoryCollector struct {
	Self types.ManagedObjectReference `json:"self"`

	Filter types.AnyType `json:"filter"`
}

func (m HistoryCollector) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HistoryCollector"] = reflect.TypeOf((*HistoryCollector)(nil)).Elem()
}

type HostAccessManager struct {
	Self types.ManagedObjectReference `json:"self"`

	LockdownMode types.HostLockdownMode `json:"lockdownMode"`
}

func (m HostAccessManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostAccessManager"] = reflect.TypeOf((*HostAccessManager)(nil)).Elem()
}

type HostActiveDirectoryAuthentication struct {
	HostDirectoryStore
}

func init() {
	t["HostActiveDirectoryAuthentication"] = reflect.TypeOf((*HostActiveDirectoryAuthentication)(nil)).Elem()
}

type HostAssignableHardwareManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Binding []types.HostAssignableHardwareBinding `json:"binding"`
	Config  types.HostAssignableHardwareConfig    `json:"config"`
}

func (m HostAssignableHardwareManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostAssignableHardwareManager"] = reflect.TypeOf((*HostAssignableHardwareManager)(nil)).Elem()
}

type HostAuthenticationManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Info           types.HostAuthenticationManagerInfo `json:"info"`
	SupportedStore []types.ManagedObjectReference      `json:"supportedStore"`
}

func (m HostAuthenticationManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostAuthenticationManager"] = reflect.TypeOf((*HostAuthenticationManager)(nil)).Elem()
}

type HostAuthenticationStore struct {
	Self types.ManagedObjectReference `json:"self"`

	Info types.BaseHostAuthenticationStoreInfo `json:"info"`
}

func (m HostAuthenticationStore) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostAuthenticationStore"] = reflect.TypeOf((*HostAuthenticationStore)(nil)).Elem()
}

type HostAutoStartManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Config types.HostAutoStartManagerConfig `json:"config"`
}

func (m HostAutoStartManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostAutoStartManager"] = reflect.TypeOf((*HostAutoStartManager)(nil)).Elem()
}

type HostBootDeviceSystem struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m HostBootDeviceSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostBootDeviceSystem"] = reflect.TypeOf((*HostBootDeviceSystem)(nil)).Elem()
}

type HostCacheConfigurationManager struct {
	Self types.ManagedObjectReference `json:"self"`

	CacheConfigurationInfo []types.HostCacheConfigurationInfo `json:"cacheConfigurationInfo"`
}

func (m HostCacheConfigurationManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostCacheConfigurationManager"] = reflect.TypeOf((*HostCacheConfigurationManager)(nil)).Elem()
}

type HostCertificateManager struct {
	Self types.ManagedObjectReference `json:"self"`

	CertificateInfo types.HostCertificateManagerCertificateInfo `json:"certificateInfo"`
}

func (m HostCertificateManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostCertificateManager"] = reflect.TypeOf((*HostCertificateManager)(nil)).Elem()
}

type HostCpuSchedulerSystem struct {
	ExtensibleManagedObject

	HyperthreadInfo *types.HostHyperThreadScheduleInfo `json:"hyperthreadInfo"`
}

func init() {
	t["HostCpuSchedulerSystem"] = reflect.TypeOf((*HostCpuSchedulerSystem)(nil)).Elem()
}

type HostDatastoreBrowser struct {
	Self types.ManagedObjectReference `json:"self"`

	Datastore     []types.ManagedObjectReference `json:"datastore"`
	SupportedType []types.BaseFileQuery          `json:"supportedType"`
}

func (m HostDatastoreBrowser) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostDatastoreBrowser"] = reflect.TypeOf((*HostDatastoreBrowser)(nil)).Elem()
}

type HostDatastoreSystem struct {
	Self types.ManagedObjectReference `json:"self"`

	Datastore    []types.ManagedObjectReference        `json:"datastore"`
	Capabilities types.HostDatastoreSystemCapabilities `json:"capabilities"`
}

func (m HostDatastoreSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostDatastoreSystem"] = reflect.TypeOf((*HostDatastoreSystem)(nil)).Elem()
}

type HostDateTimeSystem struct {
	Self types.ManagedObjectReference `json:"self"`

	DateTimeInfo types.HostDateTimeInfo `json:"dateTimeInfo"`
}

func (m HostDateTimeSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostDateTimeSystem"] = reflect.TypeOf((*HostDateTimeSystem)(nil)).Elem()
}

type HostDiagnosticSystem struct {
	Self types.ManagedObjectReference `json:"self"`

	ActivePartition *types.HostDiagnosticPartition `json:"activePartition"`
}

func (m HostDiagnosticSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostDiagnosticSystem"] = reflect.TypeOf((*HostDiagnosticSystem)(nil)).Elem()
}

type HostDirectoryStore struct {
	HostAuthenticationStore
}

func init() {
	t["HostDirectoryStore"] = reflect.TypeOf((*HostDirectoryStore)(nil)).Elem()
}

type HostEsxAgentHostManager struct {
	Self types.ManagedObjectReference `json:"self"`

	ConfigInfo types.HostEsxAgentHostManagerConfigInfo `json:"configInfo"`
}

func (m HostEsxAgentHostManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostEsxAgentHostManager"] = reflect.TypeOf((*HostEsxAgentHostManager)(nil)).Elem()
}

type HostFirewallSystem struct {
	ExtensibleManagedObject

	FirewallInfo *types.HostFirewallInfo `json:"firewallInfo"`
}

func init() {
	t["HostFirewallSystem"] = reflect.TypeOf((*HostFirewallSystem)(nil)).Elem()
}

type HostFirmwareSystem struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m HostFirmwareSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostFirmwareSystem"] = reflect.TypeOf((*HostFirmwareSystem)(nil)).Elem()
}

type HostGraphicsManager struct {
	ExtensibleManagedObject

	GraphicsInfo           []types.HostGraphicsInfo          `json:"graphicsInfo"`
	GraphicsConfig         *types.HostGraphicsConfig         `json:"graphicsConfig"`
	SharedPassthruGpuTypes []string                          `json:"sharedPassthruGpuTypes"`
	SharedGpuCapabilities  []types.HostSharedGpuCapabilities `json:"sharedGpuCapabilities"`
}

func init() {
	t["HostGraphicsManager"] = reflect.TypeOf((*HostGraphicsManager)(nil)).Elem()
}

type HostHealthStatusSystem struct {
	Self types.ManagedObjectReference `json:"self"`

	Runtime types.HealthSystemRuntime `json:"runtime"`
}

func (m HostHealthStatusSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostHealthStatusSystem"] = reflect.TypeOf((*HostHealthStatusSystem)(nil)).Elem()
}

type HostImageConfigManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m HostImageConfigManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostImageConfigManager"] = reflect.TypeOf((*HostImageConfigManager)(nil)).Elem()
}

type HostKernelModuleSystem struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m HostKernelModuleSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostKernelModuleSystem"] = reflect.TypeOf((*HostKernelModuleSystem)(nil)).Elem()
}

type HostLocalAccountManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m HostLocalAccountManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostLocalAccountManager"] = reflect.TypeOf((*HostLocalAccountManager)(nil)).Elem()
}

type HostLocalAuthentication struct {
	HostAuthenticationStore
}

func init() {
	t["HostLocalAuthentication"] = reflect.TypeOf((*HostLocalAuthentication)(nil)).Elem()
}

type HostMemorySystem struct {
	ExtensibleManagedObject

	ConsoleReservationInfo        *types.ServiceConsoleReservationInfo       `json:"consoleReservationInfo"`
	VirtualMachineReservationInfo *types.VirtualMachineMemoryReservationInfo `json:"virtualMachineReservationInfo"`
}

func init() {
	t["HostMemorySystem"] = reflect.TypeOf((*HostMemorySystem)(nil)).Elem()
}

type HostNetworkSystem struct {
	ExtensibleManagedObject

	Capabilities         *types.HostNetCapabilities        `json:"capabilities"`
	NetworkInfo          *types.HostNetworkInfo            `json:"networkInfo"`
	OffloadCapabilities  *types.HostNetOffloadCapabilities `json:"offloadCapabilities"`
	NetworkConfig        *types.HostNetworkConfig          `json:"networkConfig"`
	DnsConfig            types.BaseHostDnsConfig           `json:"dnsConfig"`
	IpRouteConfig        types.BaseHostIpRouteConfig       `json:"ipRouteConfig"`
	ConsoleIpRouteConfig types.BaseHostIpRouteConfig       `json:"consoleIpRouteConfig"`
}

func init() {
	t["HostNetworkSystem"] = reflect.TypeOf((*HostNetworkSystem)(nil)).Elem()
}

type HostNvdimmSystem struct {
	Self types.ManagedObjectReference `json:"self"`

	NvdimmSystemInfo types.NvdimmSystemInfo `json:"nvdimmSystemInfo"`
}

func (m HostNvdimmSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostNvdimmSystem"] = reflect.TypeOf((*HostNvdimmSystem)(nil)).Elem()
}

type HostPatchManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m HostPatchManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostPatchManager"] = reflect.TypeOf((*HostPatchManager)(nil)).Elem()
}

type HostPciPassthruSystem struct {
	ExtensibleManagedObject

	PciPassthruInfo     []types.BaseHostPciPassthruInfo     `json:"pciPassthruInfo"`
	SriovDevicePoolInfo []types.BaseHostSriovDevicePoolInfo `json:"sriovDevicePoolInfo"`
}

func init() {
	t["HostPciPassthruSystem"] = reflect.TypeOf((*HostPciPassthruSystem)(nil)).Elem()
}

type HostPowerSystem struct {
	Self types.ManagedObjectReference `json:"self"`

	Capability types.PowerSystemCapability `json:"capability"`
	Info       types.PowerSystemInfo       `json:"info"`
}

func (m HostPowerSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostPowerSystem"] = reflect.TypeOf((*HostPowerSystem)(nil)).Elem()
}

type HostProfile struct {
	Profile

	ValidationState           *string                                 `json:"validationState"`
	ValidationStateUpdateTime *time.Time                              `json:"validationStateUpdateTime"`
	ValidationFailureInfo     *types.HostProfileValidationFailureInfo `json:"validationFailureInfo"`
	ReferenceHost             *types.ManagedObjectReference           `json:"referenceHost"`
}

func init() {
	t["HostProfile"] = reflect.TypeOf((*HostProfile)(nil)).Elem()
}

type HostProfileManager struct {
	ProfileManager
}

func init() {
	t["HostProfileManager"] = reflect.TypeOf((*HostProfileManager)(nil)).Elem()
}

type HostServiceSystem struct {
	ExtensibleManagedObject

	ServiceInfo types.HostServiceInfo `json:"serviceInfo"`
}

func init() {
	t["HostServiceSystem"] = reflect.TypeOf((*HostServiceSystem)(nil)).Elem()
}

type HostSnmpSystem struct {
	Self types.ManagedObjectReference `json:"self"`

	Configuration types.HostSnmpConfigSpec        `json:"configuration"`
	Limits        types.HostSnmpSystemAgentLimits `json:"limits"`
}

func (m HostSnmpSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostSnmpSystem"] = reflect.TypeOf((*HostSnmpSystem)(nil)).Elem()
}

type HostSpecificationManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m HostSpecificationManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostSpecificationManager"] = reflect.TypeOf((*HostSpecificationManager)(nil)).Elem()
}

type HostStorageSystem struct {
	ExtensibleManagedObject

	StorageDeviceInfo    *types.HostStorageDeviceInfo   `json:"storageDeviceInfo"`
	FileSystemVolumeInfo types.HostFileSystemVolumeInfo `json:"fileSystemVolumeInfo"`
	SystemFile           []string                       `json:"systemFile"`
	MultipathStateInfo   *types.HostMultipathStateInfo  `json:"multipathStateInfo"`
}

func init() {
	t["HostStorageSystem"] = reflect.TypeOf((*HostStorageSystem)(nil)).Elem()
}

type HostSystem struct {
	ManagedEntity

	Runtime                    types.HostRuntimeInfo                      `json:"runtime"`
	Summary                    types.HostListSummary                      `json:"summary"`
	Hardware                   *types.HostHardwareInfo                    `json:"hardware"`
	Capability                 *types.HostCapability                      `json:"capability"`
	LicensableResource         types.HostLicensableResourceInfo           `json:"licensableResource"`
	RemediationState           *types.HostSystemRemediationState          `json:"remediationState"`
	PrecheckRemediationResult  *types.ApplyHostProfileConfigurationSpec   `json:"precheckRemediationResult"`
	RemediationResult          *types.ApplyHostProfileConfigurationResult `json:"remediationResult"`
	ComplianceCheckState       *types.HostSystemComplianceCheckState      `json:"complianceCheckState"`
	ComplianceCheckResult      *types.ComplianceResult                    `json:"complianceCheckResult"`
	ConfigManager              types.HostConfigManager                    `json:"configManager"`
	Config                     *types.HostConfigInfo                      `json:"config"`
	Vm                         []types.ManagedObjectReference             `json:"vm"`
	Datastore                  []types.ManagedObjectReference             `json:"datastore"`
	Network                    []types.ManagedObjectReference             `json:"network"`
	DatastoreBrowser           types.ManagedObjectReference               `json:"datastoreBrowser"`
	SystemResources            *types.HostSystemResourceInfo              `json:"systemResources"`
	AnswerFileValidationState  *types.AnswerFileStatusResult              `json:"answerFileValidationState"`
	AnswerFileValidationResult *types.AnswerFileStatusResult              `json:"answerFileValidationResult"`
}

func (m *HostSystem) Entity() *ManagedEntity {
	return &m.ManagedEntity
}

func init() {
	t["HostSystem"] = reflect.TypeOf((*HostSystem)(nil)).Elem()
}

type HostVFlashManager struct {
	Self types.ManagedObjectReference `json:"self"`

	VFlashConfigInfo *types.HostVFlashManagerVFlashConfigInfo `json:"vFlashConfigInfo"`
}

func (m HostVFlashManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostVFlashManager"] = reflect.TypeOf((*HostVFlashManager)(nil)).Elem()
}

type HostVMotionSystem struct {
	ExtensibleManagedObject

	NetConfig *types.HostVMotionNetConfig `json:"netConfig"`
	IpConfig  *types.HostIpConfig         `json:"ipConfig"`
}

func init() {
	t["HostVMotionSystem"] = reflect.TypeOf((*HostVMotionSystem)(nil)).Elem()
}

type HostVStorageObjectManager struct {
	VStorageObjectManagerBase
}

func init() {
	t["HostVStorageObjectManager"] = reflect.TypeOf((*HostVStorageObjectManager)(nil)).Elem()
}

type HostVirtualNicManager struct {
	ExtensibleManagedObject

	Info types.HostVirtualNicManagerInfo `json:"info"`
}

func init() {
	t["HostVirtualNicManager"] = reflect.TypeOf((*HostVirtualNicManager)(nil)).Elem()
}

type HostVsanInternalSystem struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m HostVsanInternalSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostVsanInternalSystem"] = reflect.TypeOf((*HostVsanInternalSystem)(nil)).Elem()
}

type HostVsanSystem struct {
	Self types.ManagedObjectReference `json:"self"`

	Config types.VsanHostConfigInfo `json:"config"`
}

func (m HostVsanSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HostVsanSystem"] = reflect.TypeOf((*HostVsanSystem)(nil)).Elem()
}

type HttpNfcLease struct {
	Self types.ManagedObjectReference `json:"self"`

	InitializeProgress int32                          `json:"initializeProgress"`
	TransferProgress   int32                          `json:"transferProgress"`
	Mode               string                         `json:"mode"`
	Capabilities       types.HttpNfcLeaseCapabilities `json:"capabilities"`
	Info               *types.HttpNfcLeaseInfo        `json:"info"`
	State              types.HttpNfcLeaseState        `json:"state"`
	Error              *types.LocalizedMethodFault    `json:"error"`
}

func (m HttpNfcLease) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["HttpNfcLease"] = reflect.TypeOf((*HttpNfcLease)(nil)).Elem()
}

type InventoryView struct {
	ManagedObjectView
}

func init() {
	t["InventoryView"] = reflect.TypeOf((*InventoryView)(nil)).Elem()
}

type IoFilterManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m IoFilterManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["IoFilterManager"] = reflect.TypeOf((*IoFilterManager)(nil)).Elem()
}

type IpPoolManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m IpPoolManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["IpPoolManager"] = reflect.TypeOf((*IpPoolManager)(nil)).Elem()
}

type IscsiManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m IscsiManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["IscsiManager"] = reflect.TypeOf((*IscsiManager)(nil)).Elem()
}

type LicenseAssignmentManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m LicenseAssignmentManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["LicenseAssignmentManager"] = reflect.TypeOf((*LicenseAssignmentManager)(nil)).Elem()
}

type LicenseManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Source                   types.BaseLicenseSource            `json:"source"`
	SourceAvailable          bool                               `json:"sourceAvailable"`
	Diagnostics              *types.LicenseDiagnostics          `json:"diagnostics"`
	FeatureInfo              []types.LicenseFeatureInfo         `json:"featureInfo"`
	LicensedEdition          string                             `json:"licensedEdition"`
	Licenses                 []types.LicenseManagerLicenseInfo  `json:"licenses"`
	LicenseAssignmentManager *types.ManagedObjectReference      `json:"licenseAssignmentManager"`
	Evaluation               types.LicenseManagerEvaluationInfo `json:"evaluation"`
}

func (m LicenseManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["LicenseManager"] = reflect.TypeOf((*LicenseManager)(nil)).Elem()
}

type ListView struct {
	ManagedObjectView
}

func init() {
	t["ListView"] = reflect.TypeOf((*ListView)(nil)).Elem()
}

type LocalizationManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Catalog []types.LocalizationManagerMessageCatalog `json:"catalog"`
}

func (m LocalizationManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["LocalizationManager"] = reflect.TypeOf((*LocalizationManager)(nil)).Elem()
}

type ManagedEntity struct {
	ExtensibleManagedObject

	Parent              *types.ManagedObjectReference  `json:"parent"`
	CustomValue         []types.BaseCustomFieldValue   `json:"customValue"`
	OverallStatus       types.ManagedEntityStatus      `json:"overallStatus"`
	ConfigStatus        types.ManagedEntityStatus      `json:"configStatus"`
	ConfigIssue         []types.BaseEvent              `json:"configIssue"`
	EffectiveRole       []int32                        `json:"effectiveRole"`
	Permission          []types.Permission             `json:"permission"`
	Name                string                         `json:"name"`
	DisabledMethod      []string                       `json:"disabledMethod"`
	RecentTask          []types.ManagedObjectReference `json:"recentTask"`
	DeclaredAlarmState  []types.AlarmState             `json:"declaredAlarmState"`
	TriggeredAlarmState []types.AlarmState             `json:"triggeredAlarmState"`
	AlarmActionsEnabled *bool                          `json:"alarmActionsEnabled"`
	Tag                 []types.Tag                    `json:"tag"`
}

func init() {
	t["ManagedEntity"] = reflect.TypeOf((*ManagedEntity)(nil)).Elem()
}

type ManagedObjectView struct {
	Self types.ManagedObjectReference `json:"self"`

	View []types.ManagedObjectReference `json:"view"`
}

func (m ManagedObjectView) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["ManagedObjectView"] = reflect.TypeOf((*ManagedObjectView)(nil)).Elem()
}

type MessageBusProxy struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m MessageBusProxy) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["MessageBusProxy"] = reflect.TypeOf((*MessageBusProxy)(nil)).Elem()
}

type Network struct {
	ManagedEntity

	Summary types.BaseNetworkSummary       `json:"summary"`
	Host    []types.ManagedObjectReference `json:"host"`
	Vm      []types.ManagedObjectReference `json:"vm"`
	Name    string                         `json:"name"`
}

func (m *Network) Entity() *ManagedEntity {
	return &m.ManagedEntity
}

func init() {
	t["Network"] = reflect.TypeOf((*Network)(nil)).Elem()
}

type OpaqueNetwork struct {
	Network

	Capability  *types.OpaqueNetworkCapability `json:"capability"`
	ExtraConfig []types.BaseOptionValue        `json:"extraConfig"`
}

func init() {
	t["OpaqueNetwork"] = reflect.TypeOf((*OpaqueNetwork)(nil)).Elem()
}

type OptionManager struct {
	Self types.ManagedObjectReference `json:"self"`

	SupportedOption []types.OptionDef       `json:"supportedOption"`
	Setting         []types.BaseOptionValue `json:"setting"`
}

func (m OptionManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["OptionManager"] = reflect.TypeOf((*OptionManager)(nil)).Elem()
}

type OverheadMemoryManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m OverheadMemoryManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["OverheadMemoryManager"] = reflect.TypeOf((*OverheadMemoryManager)(nil)).Elem()
}

type OvfManager struct {
	Self types.ManagedObjectReference `json:"self"`

	OvfImportOption []types.OvfOptionInfo `json:"ovfImportOption"`
	OvfExportOption []types.OvfOptionInfo `json:"ovfExportOption"`
}

func (m OvfManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["OvfManager"] = reflect.TypeOf((*OvfManager)(nil)).Elem()
}

type PerformanceManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Description        types.PerformanceDescription `json:"description"`
	HistoricalInterval []types.PerfInterval         `json:"historicalInterval"`
	PerfCounter        []types.PerfCounterInfo      `json:"perfCounter"`
}

func (m PerformanceManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["PerformanceManager"] = reflect.TypeOf((*PerformanceManager)(nil)).Elem()
}

type Profile struct {
	Self types.ManagedObjectReference `json:"self"`

	Config           types.BaseProfileConfigInfo    `json:"config"`
	Description      *types.ProfileDescription      `json:"description"`
	Name             string                         `json:"name"`
	CreatedTime      time.Time                      `json:"createdTime"`
	ModifiedTime     time.Time                      `json:"modifiedTime"`
	Entity           []types.ManagedObjectReference `json:"entity"`
	ComplianceStatus string                         `json:"complianceStatus"`
}

func (m Profile) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["Profile"] = reflect.TypeOf((*Profile)(nil)).Elem()
}

type ProfileComplianceManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m ProfileComplianceManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["ProfileComplianceManager"] = reflect.TypeOf((*ProfileComplianceManager)(nil)).Elem()
}

type ProfileManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Profile []types.ManagedObjectReference `json:"profile"`
}

func (m ProfileManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["ProfileManager"] = reflect.TypeOf((*ProfileManager)(nil)).Elem()
}

type PropertyCollector struct {
	Self types.ManagedObjectReference `json:"self"`

	Filter []types.ManagedObjectReference `json:"filter"`
}

func (m PropertyCollector) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["PropertyCollector"] = reflect.TypeOf((*PropertyCollector)(nil)).Elem()
}

type PropertyFilter struct {
	Self types.ManagedObjectReference `json:"self"`

	Spec           types.PropertyFilterSpec `json:"spec"`
	PartialUpdates bool                     `json:"partialUpdates"`
}

func (m PropertyFilter) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["PropertyFilter"] = reflect.TypeOf((*PropertyFilter)(nil)).Elem()
}

type ResourcePlanningManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m ResourcePlanningManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["ResourcePlanningManager"] = reflect.TypeOf((*ResourcePlanningManager)(nil)).Elem()
}

type ResourcePool struct {
	ManagedEntity

	Summary            types.BaseResourcePoolSummary  `json:"summary"`
	Runtime            types.ResourcePoolRuntimeInfo  `json:"runtime"`
	Owner              types.ManagedObjectReference   `json:"owner"`
	ResourcePool       []types.ManagedObjectReference `json:"resourcePool"`
	Vm                 []types.ManagedObjectReference `json:"vm"`
	Config             types.ResourceConfigSpec       `json:"config"`
	Namespace          *string                        `json:"namespace"`
	ChildConfiguration []types.ResourceConfigSpec     `json:"childConfiguration"`
}

func (m *ResourcePool) Entity() *ManagedEntity {
	return &m.ManagedEntity
}

func init() {
	t["ResourcePool"] = reflect.TypeOf((*ResourcePool)(nil)).Elem()
}

type ScheduledTask struct {
	ExtensibleManagedObject

	Info types.ScheduledTaskInfo `json:"info"`
}

func init() {
	t["ScheduledTask"] = reflect.TypeOf((*ScheduledTask)(nil)).Elem()
}

type ScheduledTaskManager struct {
	Self types.ManagedObjectReference `json:"self"`

	ScheduledTask []types.ManagedObjectReference `json:"scheduledTask"`
	Description   types.ScheduledTaskDescription `json:"description"`
}

func (m ScheduledTaskManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["ScheduledTaskManager"] = reflect.TypeOf((*ScheduledTaskManager)(nil)).Elem()
}

type SearchIndex struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m SearchIndex) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["SearchIndex"] = reflect.TypeOf((*SearchIndex)(nil)).Elem()
}

type ServiceInstance struct {
	Self types.ManagedObjectReference `json:"self"`

	ServerClock time.Time            `json:"serverClock"`
	Capability  types.Capability     `json:"capability"`
	Content     types.ServiceContent `json:"content"`
}

func (m ServiceInstance) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["ServiceInstance"] = reflect.TypeOf((*ServiceInstance)(nil)).Elem()
}

type ServiceManager struct {
	Self types.ManagedObjectReference `json:"self"`

	Service []types.ServiceManagerServiceInfo `json:"service"`
}

func (m ServiceManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["ServiceManager"] = reflect.TypeOf((*ServiceManager)(nil)).Elem()
}

type SessionManager struct {
	Self types.ManagedObjectReference `json:"self"`

	SessionList         []types.UserSession `json:"sessionList"`
	CurrentSession      *types.UserSession  `json:"currentSession"`
	Message             *string             `json:"message"`
	MessageLocaleList   []string            `json:"messageLocaleList"`
	SupportedLocaleList []string            `json:"supportedLocaleList"`
	DefaultLocale       string              `json:"defaultLocale"`
}

func (m SessionManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["SessionManager"] = reflect.TypeOf((*SessionManager)(nil)).Elem()
}

type SimpleCommand struct {
	Self types.ManagedObjectReference `json:"self"`

	EncodingType types.SimpleCommandEncoding     `json:"encodingType"`
	Entity       types.ServiceManagerServiceInfo `json:"entity"`
}

func (m SimpleCommand) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["SimpleCommand"] = reflect.TypeOf((*SimpleCommand)(nil)).Elem()
}

type SiteInfoManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m SiteInfoManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["SiteInfoManager"] = reflect.TypeOf((*SiteInfoManager)(nil)).Elem()
}

type StoragePod struct {
	Folder

	Summary            *types.StoragePodSummary  `json:"summary"`
	PodStorageDrsEntry *types.PodStorageDrsEntry `json:"podStorageDrsEntry"`
}

func init() {
	t["StoragePod"] = reflect.TypeOf((*StoragePod)(nil)).Elem()
}

type StorageQueryManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m StorageQueryManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["StorageQueryManager"] = reflect.TypeOf((*StorageQueryManager)(nil)).Elem()
}

type StorageResourceManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m StorageResourceManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["StorageResourceManager"] = reflect.TypeOf((*StorageResourceManager)(nil)).Elem()
}

type Task struct {
	ExtensibleManagedObject

	Info types.TaskInfo `json:"info"`
}

func init() {
	t["Task"] = reflect.TypeOf((*Task)(nil)).Elem()
}

type TaskHistoryCollector struct {
	HistoryCollector

	LatestPage []types.TaskInfo `json:"latestPage"`
}

func init() {
	t["TaskHistoryCollector"] = reflect.TypeOf((*TaskHistoryCollector)(nil)).Elem()
}

type TaskManager struct {
	Self types.ManagedObjectReference `json:"self"`

	RecentTask   []types.ManagedObjectReference `json:"recentTask"`
	Description  types.TaskDescription          `json:"description"`
	MaxCollector int32                          `json:"maxCollector"`
}

func (m TaskManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["TaskManager"] = reflect.TypeOf((*TaskManager)(nil)).Elem()
}

type TenantTenantManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m TenantTenantManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["TenantTenantManager"] = reflect.TypeOf((*TenantTenantManager)(nil)).Elem()
}

type UserDirectory struct {
	Self types.ManagedObjectReference `json:"self"`

	DomainList []string `json:"domainList"`
}

func (m UserDirectory) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["UserDirectory"] = reflect.TypeOf((*UserDirectory)(nil)).Elem()
}

type VStorageObjectManagerBase struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m VStorageObjectManagerBase) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["VStorageObjectManagerBase"] = reflect.TypeOf((*VStorageObjectManagerBase)(nil)).Elem()
}

type VcenterVStorageObjectManager struct {
	VStorageObjectManagerBase
}

func init() {
	t["VcenterVStorageObjectManager"] = reflect.TypeOf((*VcenterVStorageObjectManager)(nil)).Elem()
}

type View struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m View) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["View"] = reflect.TypeOf((*View)(nil)).Elem()
}

type ViewManager struct {
	Self types.ManagedObjectReference `json:"self"`

	ViewList []types.ManagedObjectReference `json:"viewList"`
}

func (m ViewManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["ViewManager"] = reflect.TypeOf((*ViewManager)(nil)).Elem()
}

type VirtualApp struct {
	ResourcePool

	ParentFolder *types.ManagedObjectReference  `json:"parentFolder"`
	Datastore    []types.ManagedObjectReference `json:"datastore"`
	Network      []types.ManagedObjectReference `json:"network"`
	VAppConfig   *types.VAppConfigInfo          `json:"vAppConfig"`
	ParentVApp   *types.ManagedObjectReference  `json:"parentVApp"`
	ChildLink    []types.VirtualAppLinkInfo     `json:"childLink"`
}

func init() {
	t["VirtualApp"] = reflect.TypeOf((*VirtualApp)(nil)).Elem()
}

type VirtualDiskManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m VirtualDiskManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["VirtualDiskManager"] = reflect.TypeOf((*VirtualDiskManager)(nil)).Elem()
}

type VirtualMachine struct {
	ManagedEntity

	Capability           types.VirtualMachineCapability    `json:"capability"`
	Config               *types.VirtualMachineConfigInfo   `json:"config"`
	Layout               *types.VirtualMachineFileLayout   `json:"layout"`
	LayoutEx             *types.VirtualMachineFileLayoutEx `json:"layoutEx"`
	Storage              *types.VirtualMachineStorageInfo  `json:"storage"`
	EnvironmentBrowser   types.ManagedObjectReference      `json:"environmentBrowser"`
	ResourcePool         *types.ManagedObjectReference     `json:"resourcePool"`
	ParentVApp           *types.ManagedObjectReference     `json:"parentVApp"`
	ResourceConfig       *types.ResourceConfigSpec         `json:"resourceConfig"`
	Runtime              types.VirtualMachineRuntimeInfo   `json:"runtime"`
	Guest                *types.GuestInfo                  `json:"guest"`
	Summary              types.VirtualMachineSummary       `json:"summary"`
	Datastore            []types.ManagedObjectReference    `json:"datastore"`
	Network              []types.ManagedObjectReference    `json:"network"`
	Snapshot             *types.VirtualMachineSnapshotInfo `json:"snapshot"`
	RootSnapshot         []types.ManagedObjectReference    `json:"rootSnapshot"`
	GuestHeartbeatStatus types.ManagedEntityStatus         `json:"guestHeartbeatStatus"`
}

func (m *VirtualMachine) Entity() *ManagedEntity {
	return &m.ManagedEntity
}

func init() {
	t["VirtualMachine"] = reflect.TypeOf((*VirtualMachine)(nil)).Elem()
}

type VirtualMachineCompatibilityChecker struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m VirtualMachineCompatibilityChecker) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["VirtualMachineCompatibilityChecker"] = reflect.TypeOf((*VirtualMachineCompatibilityChecker)(nil)).Elem()
}

type VirtualMachineGuestCustomizationManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m VirtualMachineGuestCustomizationManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["VirtualMachineGuestCustomizationManager"] = reflect.TypeOf((*VirtualMachineGuestCustomizationManager)(nil)).Elem()
}

type VirtualMachineProvisioningChecker struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m VirtualMachineProvisioningChecker) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["VirtualMachineProvisioningChecker"] = reflect.TypeOf((*VirtualMachineProvisioningChecker)(nil)).Elem()
}

type VirtualMachineSnapshot struct {
	ExtensibleManagedObject

	Config        types.VirtualMachineConfigInfo `json:"config"`
	ChildSnapshot []types.ManagedObjectReference `json:"childSnapshot"`
	Vm            types.ManagedObjectReference   `json:"vm"`
}

func init() {
	t["VirtualMachineSnapshot"] = reflect.TypeOf((*VirtualMachineSnapshot)(nil)).Elem()
}

type VirtualizationManager struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m VirtualizationManager) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["VirtualizationManager"] = reflect.TypeOf((*VirtualizationManager)(nil)).Elem()
}

type VmwareDistributedVirtualSwitch struct {
	DistributedVirtualSwitch
}

func init() {
	t["VmwareDistributedVirtualSwitch"] = reflect.TypeOf((*VmwareDistributedVirtualSwitch)(nil)).Elem()
}

type VsanUpgradeSystem struct {
	Self types.ManagedObjectReference `json:"self"`
}

func (m VsanUpgradeSystem) Reference() types.ManagedObjectReference {
	return m.Self
}

func init() {
	t["VsanUpgradeSystem"] = reflect.TypeOf((*VsanUpgradeSystem)(nil)).Elem()
}
