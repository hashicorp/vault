# Release History

## 4.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 4.2.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module
- New value `StorageAccountTypeStandardSSDLRS` added to enum type `StorageAccountType`
- New field `ComputerName` in struct `VirtualMachineScaleSetVMInstanceView`
- New field `HyperVGeneration` in struct `VirtualMachineScaleSetVMInstanceView`
- New field `OSName` in struct `VirtualMachineScaleSetVMInstanceView`
- New field `OSVersion` in struct `VirtualMachineScaleSetVMInstanceView`


## 4.1.0 (2023-01-27)
### Features Added

- New type alias `AlternativeType` with values `AlternativeTypeNone`, `AlternativeTypeOffer`, `AlternativeTypePlan`
- New type alias `ImageState` with values `ImageStateActive`, `ImageStateDeprecated`, `ImageStateScheduledForDeprecation`
- New struct `AlternativeOption`
- New struct `ImageDeprecationStatus`
- New struct `OSImageNotificationProfile`
- New struct `OSProfileProvisioningData`
- New struct `ServiceArtifactReference`
- New field `Zones` in struct `CloudService`
- New field `UserData` in struct `RestorePointSourceMetadata`
- New field `MaxSurge` in struct `RollingUpgradePolicy`
- New field `RollbackFailedInstancesOnPolicyBreach` in struct `RollingUpgradePolicy`
- New field `OSImageNotificationProfile` in struct `ScheduledEventsProfile`
- New field `ImageDeprecationStatus` in struct `VirtualMachineImageProperties`
- New field `ExactVersion` in struct `VirtualMachineReimageParameters`
- New field `OSProfile` in struct `VirtualMachineReimageParameters`
- New field `RequireGuestProvisionSignal` in struct `VirtualMachineScaleSetOSProfile`
- New field `ConstrainedMaximumCapacity` in struct `VirtualMachineScaleSetProperties`
- New field `ExactVersion` in struct `VirtualMachineScaleSetReimageParameters`
- New field `OSProfile` in struct `VirtualMachineScaleSetReimageParameters`
- New field `ServiceArtifactReference` in struct `VirtualMachineScaleSetVMProfile`
- New field `ExactVersion` in struct `VirtualMachineScaleSetVMReimageParameters`
- New field `OSProfile` in struct `VirtualMachineScaleSetVMReimageParameters`


## 4.0.0 (2022-10-04)
### Breaking Changes

- Type of `GalleryImageVersionStorageProfile.Source` has been changed from `*GalleryArtifactVersionSource` to `*GalleryArtifactVersionFullSource`
- Type of `SharingProfile.CommunityGalleryInfo` has been changed from `interface{}` to `*CommunityGalleryInfo`
- Type of `VirtualMachineExtensionUpdateProperties.ProtectedSettingsFromKeyVault` has been changed from `interface{}` to `*KeyVaultSecretReference`
- Type of `GalleryOSDiskImage.Source` has been changed from `*GalleryArtifactVersionSource` to `*GalleryDiskImageSource`
- Type of `GalleryDiskImage.Source` has been changed from `*GalleryArtifactVersionSource` to `*GalleryDiskImageSource`
- Type of `GalleryDataDiskImage.Source` has been changed from `*GalleryArtifactVersionSource` to `*GalleryDiskImageSource`
- Type of `VirtualMachineScaleSetExtensionProperties.ProtectedSettingsFromKeyVault` has been changed from `interface{}` to `*KeyVaultSecretReference`
- Type of `VirtualMachineExtensionProperties.ProtectedSettingsFromKeyVault` has been changed from `interface{}` to `*KeyVaultSecretReference`
- Field `URI` of struct `GalleryArtifactVersionSource` has been removed

### Features Added

- New const `DiskControllerTypesSCSI`
- New const `PolicyViolationCategoryImageFlaggedUnsafe`
- New const `GalleryApplicationCustomActionParameterTypeConfigurationDataBlob`
- New const `PolicyViolationCategoryIPTheft`
- New const `PolicyViolationCategoryCopyrightValidation`
- New const `PolicyViolationCategoryOther`
- New const `GalleryApplicationCustomActionParameterTypeString`
- New const `DiskControllerTypesNVMe`
- New const `GalleryApplicationCustomActionParameterTypeLogOutputBlob`
- New type alias `DiskControllerTypes`
- New type alias `PolicyViolationCategory`
- New type alias `GalleryApplicationCustomActionParameterType`
- New function `PossiblePolicyViolationCategoryValues() []PolicyViolationCategory`
- New function `PossibleGalleryApplicationCustomActionParameterTypeValues() []GalleryApplicationCustomActionParameterType`
- New function `PossibleDiskControllerTypesValues() []DiskControllerTypes`
- New struct `GalleryApplicationCustomAction`
- New struct `GalleryApplicationCustomActionParameter`
- New struct `GalleryApplicationVersionSafetyProfile`
- New struct `GalleryArtifactSafetyProfileBase`
- New struct `GalleryArtifactVersionFullSource`
- New struct `GalleryDiskImageSource`
- New struct `GalleryImageVersionSafetyProfile`
- New struct `LatestGalleryImageVersion`
- New struct `PolicyViolation`
- New struct `PriorityMixPolicy`
- New field `DiskControllerType` in struct `VirtualMachineScaleSetUpdateStorageProfile`
- New field `HardwareProfile` in struct `VirtualMachineScaleSetUpdateVMProfile`
- New field `CustomActions` in struct `GalleryApplicationProperties`
- New field `DisableTCPStateTracking` in struct `VirtualMachineScaleSetNetworkConfigurationProperties`
- New field `DiskControllerType` in struct `StorageProfile`
- New field `OptimizedForFrequentAttach` in struct `DiskProperties`
- New field `BurstingEnabledTime` in struct `DiskProperties`
- New field `DiskControllerTypes` in struct `SupportedCapabilities`
- New field `DisableTCPStateTracking` in struct `VirtualMachineNetworkInterfaceConfigurationProperties`
- New field `EnableVMAgentPlatformUpdates` in struct `WindowsConfiguration`
- New field `PerformancePlus` in struct `CreationData`
- New field `IncrementalSnapshotFamilyID` in struct `SnapshotProperties`
- New field `OptimizedForFrequentAttach` in struct `DiskUpdateProperties`
- New field `DisableTCPStateTracking` in struct `VirtualMachineScaleSetUpdateNetworkConfigurationProperties`
- New field `ExcludeFromLatest` in struct `TargetRegion`
- New field `PrivacyStatementURI` in struct `SharedGalleryImageProperties`
- New field `Eula` in struct `SharedGalleryImageProperties`
- New field `SafetyProfile` in struct `GalleryApplicationVersionProperties`
- New field `SafetyProfile` in struct `GalleryImageVersionProperties`
- New field `EnableVMAgentPlatformUpdates` in struct `LinuxConfiguration`
- New field `CurrentCapacity` in struct `CapacityReservationUtilization`
- New field `PriorityMixPolicy` in struct `VirtualMachineScaleSetProperties`
- New field `CustomActions` in struct `GalleryApplicationVersionPublishingProfile`
- New field `PlatformFaultDomainCount` in struct `CapacityReservationProperties`
- New field `DiskControllerType` in struct `VirtualMachineScaleSetStorageProfile`


## 3.0.1 (2022-07-29)
### Other Changes
- Fix wrong module import for live test

## 3.0.0 (2022-06-24)
### Breaking Changes

- Function `*CloudServicesClient.BeginCreateOrUpdate` parameter(s) have been changed from `(context.Context, string, string, *CloudServicesClientBeginCreateOrUpdateOptions)` to `(context.Context, string, string, CloudService, *CloudServicesClientBeginCreateOrUpdateOptions)`
- Function `*CloudServicesClient.BeginUpdate` parameter(s) have been changed from `(context.Context, string, string, *CloudServicesClientBeginUpdateOptions)` to `(context.Context, string, string, CloudServiceUpdate, *CloudServicesClientBeginUpdateOptions)`
- Function `*CloudServicesUpdateDomainClient.BeginWalkUpdateDomain` parameter(s) have been changed from `(context.Context, string, string, int32, *CloudServicesUpdateDomainClientBeginWalkUpdateDomainOptions)` to `(context.Context, string, string, int32, UpdateDomain, *CloudServicesUpdateDomainClientBeginWalkUpdateDomainOptions)`
- Type of `CloudServiceExtensionProperties.Settings` has been changed from `*string` to `interface{}`
- Type of `CloudServiceExtensionProperties.ProtectedSettings` has been changed from `*string` to `interface{}`
- Field `Parameters` of struct `CloudServicesClientBeginUpdateOptions` has been removed
- Field `Parameters` of struct `CloudServicesClientBeginCreateOrUpdateOptions` has been removed
- Field `Parameters` of struct `CloudServicesUpdateDomainClientBeginWalkUpdateDomainOptions` has been removed

### Features Added

- New const `CloudServiceSlotTypeProduction`
- New const `CloudServiceSlotTypeStaging`
- New function `*VirtualMachineImagesClient.ListByEdgeZone(context.Context, string, string, *VirtualMachineImagesClientListByEdgeZoneOptions) (VirtualMachineImagesClientListByEdgeZoneResponse, error)`
- New function `PossibleCloudServiceSlotTypeValues() []CloudServiceSlotType`
- New struct `SystemData`
- New struct `VMImagesInEdgeZoneListResult`
- New struct `VirtualMachineImagesClientListByEdgeZoneOptions`
- New struct `VirtualMachineImagesClientListByEdgeZoneResponse`
- New field `SystemData` in struct `CloudService`
- New field `SlotType` in struct `CloudServiceNetworkProfile`


## 2.0.0 (2022-06-02)
### Breaking Changes

- Type of `GalleryProperties.ProvisioningState` has been changed from `*GalleryPropertiesProvisioningState` to `*GalleryProvisioningState`
- Type of `GalleryImageVersionProperties.ProvisioningState` has been changed from `*GalleryImageVersionPropertiesProvisioningState` to `*GalleryProvisioningState`
- Type of `GalleryImageProperties.ProvisioningState` has been changed from `*GalleryImagePropertiesProvisioningState` to `*GalleryProvisioningState`
- Type of `GalleryApplicationVersionProperties.ProvisioningState` has been changed from `*GalleryApplicationVersionPropertiesProvisioningState` to `*GalleryProvisioningState`
- Type of `VirtualMachineScaleSetIdentity.UserAssignedIdentities` has been changed from `map[string]*VirtualMachineScaleSetIdentityUserAssignedIdentitiesValue` to `map[string]*UserAssignedIdentitiesValue`
- Const `GalleryImagePropertiesProvisioningStateFailed` has been removed
- Const `GalleryImagePropertiesProvisioningStateMigrating` has been removed
- Const `GalleryImageVersionPropertiesProvisioningStateCreating` has been removed
- Const `GalleryImageVersionPropertiesProvisioningStateMigrating` has been removed
- Const `GalleryApplicationVersionPropertiesProvisioningStateFailed` has been removed
- Const `GalleryPropertiesProvisioningStateMigrating` has been removed
- Const `GalleryApplicationVersionPropertiesProvisioningStateDeleting` has been removed
- Const `GalleryPropertiesProvisioningStateDeleting` has been removed
- Const `GalleryApplicationVersionPropertiesProvisioningStateCreating` has been removed
- Const `GalleryImageVersionPropertiesProvisioningStateSucceeded` has been removed
- Const `GalleryImagePropertiesProvisioningStateCreating` has been removed
- Const `GalleryImagePropertiesProvisioningStateUpdating` has been removed
- Const `GalleryImageVersionPropertiesProvisioningStateDeleting` has been removed
- Const `GalleryPropertiesProvisioningStateFailed` has been removed
- Const `SharingProfileGroupTypesCommunity` has been removed
- Const `GalleryApplicationVersionPropertiesProvisioningStateSucceeded` has been removed
- Const `GalleryApplicationVersionPropertiesProvisioningStateMigrating` has been removed
- Const `GalleryPropertiesProvisioningStateUpdating` has been removed
- Const `GalleryImageVersionPropertiesProvisioningStateFailed` has been removed
- Const `GalleryImagePropertiesProvisioningStateDeleting` has been removed
- Const `GalleryImageVersionPropertiesProvisioningStateUpdating` has been removed
- Const `GalleryPropertiesProvisioningStateCreating` has been removed
- Const `GalleryApplicationVersionPropertiesProvisioningStateUpdating` has been removed
- Const `GalleryImagePropertiesProvisioningStateSucceeded` has been removed
- Const `GalleryPropertiesProvisioningStateSucceeded` has been removed
- Function `PossibleGalleryPropertiesProvisioningStateValues` has been removed
- Function `PossibleGalleryImageVersionPropertiesProvisioningStateValues` has been removed
- Function `PossibleGalleryImagePropertiesProvisioningStateValues` has been removed
- Function `PossibleGalleryApplicationVersionPropertiesProvisioningStateValues` has been removed
- Struct `VirtualMachineScaleSetIdentityUserAssignedIdentitiesValue` has been removed

### Features Added

- New const `GallerySharingPermissionTypesCommunity`
- New const `GalleryProvisioningStateUpdating`
- New const `SharedGalleryHostCachingReadOnly`
- New const `SharedGalleryHostCachingNone`
- New const `GalleryProvisioningStateSucceeded`
- New const `GalleryProvisioningStateFailed`
- New const `SharedGalleryHostCachingReadWrite`
- New const `GalleryProvisioningStateCreating`
- New const `DiskEncryptionSetIdentityTypeUserAssigned`
- New const `GalleryProvisioningStateMigrating`
- New const `DiskEncryptionSetIdentityTypeSystemAssignedUserAssigned`
- New const `CopyCompletionErrorReasonCopySourceNotFound`
- New const `GalleryProvisioningStateDeleting`
- New const `DiskStorageAccountTypesPremiumV2LRS`
- New function `PossibleCopyCompletionErrorReasonValues() []CopyCompletionErrorReason`
- New function `PossibleSharedGalleryHostCachingValues() []SharedGalleryHostCaching`
- New function `PossibleGalleryProvisioningStateValues() []GalleryProvisioningState`
- New function `EncryptionSetIdentity.MarshalJSON() ([]byte, error)`
- New function `*CommunityGalleryImagesClient.NewListPager(string, string, *CommunityGalleryImagesClientListOptions) *runtime.Pager[CommunityGalleryImagesClientListResponse]`
- New function `*CommunityGalleryImageVersionsClient.NewListPager(string, string, string, *CommunityGalleryImageVersionsClientListOptions) *runtime.Pager[CommunityGalleryImageVersionsClientListResponse]`
- New struct `CommunityGalleryImageList`
- New struct `CommunityGalleryImageVersionList`
- New struct `CommunityGalleryImageVersionsClientListOptions`
- New struct `CommunityGalleryImageVersionsClientListResponse`
- New struct `CommunityGalleryImagesClientListOptions`
- New struct `CommunityGalleryImagesClientListResponse`
- New struct `CopyCompletionError`
- New struct `SharedGalleryDataDiskImage`
- New struct `SharedGalleryDiskImage`
- New struct `SharedGalleryImageVersionStorageProfile`
- New struct `SharedGalleryOSDiskImage`
- New struct `UserArtifactSettings`
- New field `SharedGalleryImageID` in struct `ImageDiskReference`
- New field `CommunityGalleryImageID` in struct `ImageDiskReference`
- New field `AdvancedSettings` in struct `GalleryApplicationVersionPublishingProfile`
- New field `Settings` in struct `GalleryApplicationVersionPublishingProfile`
- New field `CopyCompletionError` in struct `SnapshotProperties`
- New field `ExcludeFromLatest` in struct `SharedGalleryImageVersionProperties`
- New field `StorageProfile` in struct `SharedGalleryImageVersionProperties`
- New field `ExcludeFromLatest` in struct `CommunityGalleryImageVersionProperties`
- New field `StorageProfile` in struct `CommunityGalleryImageVersionProperties`
- New field `Architecture` in struct `SharedGalleryImageProperties`
- New field `UserAssignedIdentities` in struct `EncryptionSetIdentity`
- New field `Eula` in struct `CommunityGalleryImageProperties`
- New field `PrivacyStatementURI` in struct `CommunityGalleryImageProperties`
- New field `Architecture` in struct `CommunityGalleryImageProperties`
- New field `FederatedClientID` in struct `DiskEncryptionSetUpdateProperties`
- New field `FederatedClientID` in struct `EncryptionSetProperties`
- New field `SecurityProfile` in struct `DiskRestorePointProperties`


## 1.0.0 (2022-05-16)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 1.0.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).