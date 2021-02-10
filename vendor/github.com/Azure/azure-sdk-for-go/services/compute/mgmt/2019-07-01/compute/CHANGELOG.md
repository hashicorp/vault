Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82/specification/compute/resource-manager/readme.md tag: `package-2019-07`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *ContainerServicesCreateOrUpdateFuture.Result(ContainerServicesClient) (ContainerService, error)
1. *ContainerServicesDeleteFuture.Result(ContainerServicesClient) (autorest.Response, error)
1. *DedicatedHostsCreateOrUpdateFuture.Result(DedicatedHostsClient) (DedicatedHost, error)
1. *DedicatedHostsDeleteFuture.Result(DedicatedHostsClient) (autorest.Response, error)
1. *DedicatedHostsUpdateFuture.Result(DedicatedHostsClient) (DedicatedHost, error)
1. *DiskEncryptionSetsCreateOrUpdateFuture.Result(DiskEncryptionSetsClient) (DiskEncryptionSet, error)
1. *DiskEncryptionSetsDeleteFuture.Result(DiskEncryptionSetsClient) (autorest.Response, error)
1. *DiskEncryptionSetsUpdateFuture.Result(DiskEncryptionSetsClient) (DiskEncryptionSet, error)
1. *DisksCreateOrUpdateFuture.Result(DisksClient) (Disk, error)
1. *DisksDeleteFuture.Result(DisksClient) (autorest.Response, error)
1. *DisksGrantAccessFuture.Result(DisksClient) (AccessURI, error)
1. *DisksRevokeAccessFuture.Result(DisksClient) (autorest.Response, error)
1. *DisksUpdateFuture.Result(DisksClient) (Disk, error)
1. *GalleriesCreateOrUpdateFuture.Result(GalleriesClient) (Gallery, error)
1. *GalleriesDeleteFuture.Result(GalleriesClient) (autorest.Response, error)
1. *GalleriesUpdateFuture.Result(GalleriesClient) (Gallery, error)
1. *GalleryApplicationVersionsCreateOrUpdateFuture.Result(GalleryApplicationVersionsClient) (GalleryApplicationVersion, error)
1. *GalleryApplicationVersionsDeleteFuture.Result(GalleryApplicationVersionsClient) (autorest.Response, error)
1. *GalleryApplicationVersionsUpdateFuture.Result(GalleryApplicationVersionsClient) (GalleryApplicationVersion, error)
1. *GalleryApplicationsCreateOrUpdateFuture.Result(GalleryApplicationsClient) (GalleryApplication, error)
1. *GalleryApplicationsDeleteFuture.Result(GalleryApplicationsClient) (autorest.Response, error)
1. *GalleryApplicationsUpdateFuture.Result(GalleryApplicationsClient) (GalleryApplication, error)
1. *GalleryImageVersionsCreateOrUpdateFuture.Result(GalleryImageVersionsClient) (GalleryImageVersion, error)
1. *GalleryImageVersionsDeleteFuture.Result(GalleryImageVersionsClient) (autorest.Response, error)
1. *GalleryImageVersionsUpdateFuture.Result(GalleryImageVersionsClient) (GalleryImageVersion, error)
1. *GalleryImagesCreateOrUpdateFuture.Result(GalleryImagesClient) (GalleryImage, error)
1. *GalleryImagesDeleteFuture.Result(GalleryImagesClient) (autorest.Response, error)
1. *GalleryImagesUpdateFuture.Result(GalleryImagesClient) (GalleryImage, error)
1. *ImagesCreateOrUpdateFuture.Result(ImagesClient) (Image, error)
1. *ImagesDeleteFuture.Result(ImagesClient) (autorest.Response, error)
1. *ImagesUpdateFuture.Result(ImagesClient) (Image, error)
1. *LogAnalyticsExportRequestRateByIntervalFuture.Result(LogAnalyticsClient) (LogAnalyticsOperationResult, error)
1. *LogAnalyticsExportThrottledRequestsFuture.Result(LogAnalyticsClient) (LogAnalyticsOperationResult, error)
1. *SnapshotsCreateOrUpdateFuture.Result(SnapshotsClient) (Snapshot, error)
1. *SnapshotsDeleteFuture.Result(SnapshotsClient) (autorest.Response, error)
1. *SnapshotsGrantAccessFuture.Result(SnapshotsClient) (AccessURI, error)
1. *SnapshotsRevokeAccessFuture.Result(SnapshotsClient) (autorest.Response, error)
1. *SnapshotsUpdateFuture.Result(SnapshotsClient) (Snapshot, error)
1. *VirtualMachineExtensionsCreateOrUpdateFuture.Result(VirtualMachineExtensionsClient) (VirtualMachineExtension, error)
1. *VirtualMachineExtensionsDeleteFuture.Result(VirtualMachineExtensionsClient) (autorest.Response, error)
1. *VirtualMachineExtensionsUpdateFuture.Result(VirtualMachineExtensionsClient) (VirtualMachineExtension, error)
1. *VirtualMachineScaleSetExtensionsCreateOrUpdateFuture.Result(VirtualMachineScaleSetExtensionsClient) (VirtualMachineScaleSetExtension, error)
1. *VirtualMachineScaleSetExtensionsDeleteFuture.Result(VirtualMachineScaleSetExtensionsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetExtensionsUpdateFuture.Result(VirtualMachineScaleSetExtensionsClient) (VirtualMachineScaleSetExtension, error)
1. *VirtualMachineScaleSetRollingUpgradesCancelFuture.Result(VirtualMachineScaleSetRollingUpgradesClient) (autorest.Response, error)
1. *VirtualMachineScaleSetRollingUpgradesStartExtensionUpgradeFuture.Result(VirtualMachineScaleSetRollingUpgradesClient) (autorest.Response, error)
1. *VirtualMachineScaleSetRollingUpgradesStartOSUpgradeFuture.Result(VirtualMachineScaleSetRollingUpgradesClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMExtensionsCreateOrUpdateFuture.Result(VirtualMachineScaleSetVMExtensionsClient) (VirtualMachineExtension, error)
1. *VirtualMachineScaleSetVMExtensionsDeleteFuture.Result(VirtualMachineScaleSetVMExtensionsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMExtensionsUpdateFuture.Result(VirtualMachineScaleSetVMExtensionsClient) (VirtualMachineExtension, error)
1. *VirtualMachineScaleSetVMsDeallocateFuture.Result(VirtualMachineScaleSetVMsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMsDeleteFuture.Result(VirtualMachineScaleSetVMsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMsPerformMaintenanceFuture.Result(VirtualMachineScaleSetVMsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMsPowerOffFuture.Result(VirtualMachineScaleSetVMsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMsRedeployFuture.Result(VirtualMachineScaleSetVMsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMsReimageAllFuture.Result(VirtualMachineScaleSetVMsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMsReimageFuture.Result(VirtualMachineScaleSetVMsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMsRestartFuture.Result(VirtualMachineScaleSetVMsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMsRunCommandFuture.Result(VirtualMachineScaleSetVMsClient) (RunCommandResult, error)
1. *VirtualMachineScaleSetVMsStartFuture.Result(VirtualMachineScaleSetVMsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetVMsUpdateFuture.Result(VirtualMachineScaleSetVMsClient) (VirtualMachineScaleSetVM, error)
1. *VirtualMachineScaleSetsCreateOrUpdateFuture.Result(VirtualMachineScaleSetsClient) (VirtualMachineScaleSet, error)
1. *VirtualMachineScaleSetsDeallocateFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetsDeleteFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetsDeleteInstancesFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetsPerformMaintenanceFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetsPowerOffFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetsRedeployFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetsReimageAllFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetsReimageFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetsRestartFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetsStartFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachineScaleSetsUpdateFuture.Result(VirtualMachineScaleSetsClient) (VirtualMachineScaleSet, error)
1. *VirtualMachineScaleSetsUpdateInstancesFuture.Result(VirtualMachineScaleSetsClient) (autorest.Response, error)
1. *VirtualMachinesCaptureFuture.Result(VirtualMachinesClient) (VirtualMachineCaptureResult, error)
1. *VirtualMachinesConvertToManagedDisksFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesCreateOrUpdateFuture.Result(VirtualMachinesClient) (VirtualMachine, error)
1. *VirtualMachinesDeallocateFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesDeleteFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesPerformMaintenanceFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesPowerOffFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesReapplyFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesRedeployFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesReimageFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesRestartFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesRunCommandFuture.Result(VirtualMachinesClient) (RunCommandResult, error)
1. *VirtualMachinesStartFuture.Result(VirtualMachinesClient) (autorest.Response, error)
1. *VirtualMachinesUpdateFuture.Result(VirtualMachinesClient) (VirtualMachine, error)

## Struct Changes

### Removed Struct Fields

1. ContainerServicesCreateOrUpdateFuture.azure.Future
1. ContainerServicesDeleteFuture.azure.Future
1. DedicatedHostsCreateOrUpdateFuture.azure.Future
1. DedicatedHostsDeleteFuture.azure.Future
1. DedicatedHostsUpdateFuture.azure.Future
1. DiskEncryptionSetsCreateOrUpdateFuture.azure.Future
1. DiskEncryptionSetsDeleteFuture.azure.Future
1. DiskEncryptionSetsUpdateFuture.azure.Future
1. DisksCreateOrUpdateFuture.azure.Future
1. DisksDeleteFuture.azure.Future
1. DisksGrantAccessFuture.azure.Future
1. DisksRevokeAccessFuture.azure.Future
1. DisksUpdateFuture.azure.Future
1. GalleriesCreateOrUpdateFuture.azure.Future
1. GalleriesDeleteFuture.azure.Future
1. GalleriesUpdateFuture.azure.Future
1. GalleryApplicationVersionsCreateOrUpdateFuture.azure.Future
1. GalleryApplicationVersionsDeleteFuture.azure.Future
1. GalleryApplicationVersionsUpdateFuture.azure.Future
1. GalleryApplicationsCreateOrUpdateFuture.azure.Future
1. GalleryApplicationsDeleteFuture.azure.Future
1. GalleryApplicationsUpdateFuture.azure.Future
1. GalleryImageVersionsCreateOrUpdateFuture.azure.Future
1. GalleryImageVersionsDeleteFuture.azure.Future
1. GalleryImageVersionsUpdateFuture.azure.Future
1. GalleryImagesCreateOrUpdateFuture.azure.Future
1. GalleryImagesDeleteFuture.azure.Future
1. GalleryImagesUpdateFuture.azure.Future
1. ImagesCreateOrUpdateFuture.azure.Future
1. ImagesDeleteFuture.azure.Future
1. ImagesUpdateFuture.azure.Future
1. LogAnalyticsExportRequestRateByIntervalFuture.azure.Future
1. LogAnalyticsExportThrottledRequestsFuture.azure.Future
1. SnapshotsCreateOrUpdateFuture.azure.Future
1. SnapshotsDeleteFuture.azure.Future
1. SnapshotsGrantAccessFuture.azure.Future
1. SnapshotsRevokeAccessFuture.azure.Future
1. SnapshotsUpdateFuture.azure.Future
1. VirtualMachineExtensionsCreateOrUpdateFuture.azure.Future
1. VirtualMachineExtensionsDeleteFuture.azure.Future
1. VirtualMachineExtensionsUpdateFuture.azure.Future
1. VirtualMachineScaleSetExtensionsCreateOrUpdateFuture.azure.Future
1. VirtualMachineScaleSetExtensionsDeleteFuture.azure.Future
1. VirtualMachineScaleSetExtensionsUpdateFuture.azure.Future
1. VirtualMachineScaleSetRollingUpgradesCancelFuture.azure.Future
1. VirtualMachineScaleSetRollingUpgradesStartExtensionUpgradeFuture.azure.Future
1. VirtualMachineScaleSetRollingUpgradesStartOSUpgradeFuture.azure.Future
1. VirtualMachineScaleSetVMExtensionsCreateOrUpdateFuture.azure.Future
1. VirtualMachineScaleSetVMExtensionsDeleteFuture.azure.Future
1. VirtualMachineScaleSetVMExtensionsUpdateFuture.azure.Future
1. VirtualMachineScaleSetVMsDeallocateFuture.azure.Future
1. VirtualMachineScaleSetVMsDeleteFuture.azure.Future
1. VirtualMachineScaleSetVMsPerformMaintenanceFuture.azure.Future
1. VirtualMachineScaleSetVMsPowerOffFuture.azure.Future
1. VirtualMachineScaleSetVMsRedeployFuture.azure.Future
1. VirtualMachineScaleSetVMsReimageAllFuture.azure.Future
1. VirtualMachineScaleSetVMsReimageFuture.azure.Future
1. VirtualMachineScaleSetVMsRestartFuture.azure.Future
1. VirtualMachineScaleSetVMsRunCommandFuture.azure.Future
1. VirtualMachineScaleSetVMsStartFuture.azure.Future
1. VirtualMachineScaleSetVMsUpdateFuture.azure.Future
1. VirtualMachineScaleSetsCreateOrUpdateFuture.azure.Future
1. VirtualMachineScaleSetsDeallocateFuture.azure.Future
1. VirtualMachineScaleSetsDeleteFuture.azure.Future
1. VirtualMachineScaleSetsDeleteInstancesFuture.azure.Future
1. VirtualMachineScaleSetsPerformMaintenanceFuture.azure.Future
1. VirtualMachineScaleSetsPowerOffFuture.azure.Future
1. VirtualMachineScaleSetsRedeployFuture.azure.Future
1. VirtualMachineScaleSetsReimageAllFuture.azure.Future
1. VirtualMachineScaleSetsReimageFuture.azure.Future
1. VirtualMachineScaleSetsRestartFuture.azure.Future
1. VirtualMachineScaleSetsStartFuture.azure.Future
1. VirtualMachineScaleSetsUpdateFuture.azure.Future
1. VirtualMachineScaleSetsUpdateInstancesFuture.azure.Future
1. VirtualMachinesCaptureFuture.azure.Future
1. VirtualMachinesConvertToManagedDisksFuture.azure.Future
1. VirtualMachinesCreateOrUpdateFuture.azure.Future
1. VirtualMachinesDeallocateFuture.azure.Future
1. VirtualMachinesDeleteFuture.azure.Future
1. VirtualMachinesPerformMaintenanceFuture.azure.Future
1. VirtualMachinesPowerOffFuture.azure.Future
1. VirtualMachinesReapplyFuture.azure.Future
1. VirtualMachinesRedeployFuture.azure.Future
1. VirtualMachinesReimageFuture.azure.Future
1. VirtualMachinesRestartFuture.azure.Future
1. VirtualMachinesRunCommandFuture.azure.Future
1. VirtualMachinesStartFuture.azure.Future
1. VirtualMachinesUpdateFuture.azure.Future

## Struct Changes

### New Struct Fields

1. ContainerServicesCreateOrUpdateFuture.Result
1. ContainerServicesCreateOrUpdateFuture.azure.FutureAPI
1. ContainerServicesDeleteFuture.Result
1. ContainerServicesDeleteFuture.azure.FutureAPI
1. DedicatedHostsCreateOrUpdateFuture.Result
1. DedicatedHostsCreateOrUpdateFuture.azure.FutureAPI
1. DedicatedHostsDeleteFuture.Result
1. DedicatedHostsDeleteFuture.azure.FutureAPI
1. DedicatedHostsUpdateFuture.Result
1. DedicatedHostsUpdateFuture.azure.FutureAPI
1. DiskEncryptionSetsCreateOrUpdateFuture.Result
1. DiskEncryptionSetsCreateOrUpdateFuture.azure.FutureAPI
1. DiskEncryptionSetsDeleteFuture.Result
1. DiskEncryptionSetsDeleteFuture.azure.FutureAPI
1. DiskEncryptionSetsUpdateFuture.Result
1. DiskEncryptionSetsUpdateFuture.azure.FutureAPI
1. DisksCreateOrUpdateFuture.Result
1. DisksCreateOrUpdateFuture.azure.FutureAPI
1. DisksDeleteFuture.Result
1. DisksDeleteFuture.azure.FutureAPI
1. DisksGrantAccessFuture.Result
1. DisksGrantAccessFuture.azure.FutureAPI
1. DisksRevokeAccessFuture.Result
1. DisksRevokeAccessFuture.azure.FutureAPI
1. DisksUpdateFuture.Result
1. DisksUpdateFuture.azure.FutureAPI
1. GalleriesCreateOrUpdateFuture.Result
1. GalleriesCreateOrUpdateFuture.azure.FutureAPI
1. GalleriesDeleteFuture.Result
1. GalleriesDeleteFuture.azure.FutureAPI
1. GalleriesUpdateFuture.Result
1. GalleriesUpdateFuture.azure.FutureAPI
1. GalleryApplicationVersionsCreateOrUpdateFuture.Result
1. GalleryApplicationVersionsCreateOrUpdateFuture.azure.FutureAPI
1. GalleryApplicationVersionsDeleteFuture.Result
1. GalleryApplicationVersionsDeleteFuture.azure.FutureAPI
1. GalleryApplicationVersionsUpdateFuture.Result
1. GalleryApplicationVersionsUpdateFuture.azure.FutureAPI
1. GalleryApplicationsCreateOrUpdateFuture.Result
1. GalleryApplicationsCreateOrUpdateFuture.azure.FutureAPI
1. GalleryApplicationsDeleteFuture.Result
1. GalleryApplicationsDeleteFuture.azure.FutureAPI
1. GalleryApplicationsUpdateFuture.Result
1. GalleryApplicationsUpdateFuture.azure.FutureAPI
1. GalleryImageVersionsCreateOrUpdateFuture.Result
1. GalleryImageVersionsCreateOrUpdateFuture.azure.FutureAPI
1. GalleryImageVersionsDeleteFuture.Result
1. GalleryImageVersionsDeleteFuture.azure.FutureAPI
1. GalleryImageVersionsUpdateFuture.Result
1. GalleryImageVersionsUpdateFuture.azure.FutureAPI
1. GalleryImagesCreateOrUpdateFuture.Result
1. GalleryImagesCreateOrUpdateFuture.azure.FutureAPI
1. GalleryImagesDeleteFuture.Result
1. GalleryImagesDeleteFuture.azure.FutureAPI
1. GalleryImagesUpdateFuture.Result
1. GalleryImagesUpdateFuture.azure.FutureAPI
1. ImagesCreateOrUpdateFuture.Result
1. ImagesCreateOrUpdateFuture.azure.FutureAPI
1. ImagesDeleteFuture.Result
1. ImagesDeleteFuture.azure.FutureAPI
1. ImagesUpdateFuture.Result
1. ImagesUpdateFuture.azure.FutureAPI
1. LogAnalyticsExportRequestRateByIntervalFuture.Result
1. LogAnalyticsExportRequestRateByIntervalFuture.azure.FutureAPI
1. LogAnalyticsExportThrottledRequestsFuture.Result
1. LogAnalyticsExportThrottledRequestsFuture.azure.FutureAPI
1. SnapshotsCreateOrUpdateFuture.Result
1. SnapshotsCreateOrUpdateFuture.azure.FutureAPI
1. SnapshotsDeleteFuture.Result
1. SnapshotsDeleteFuture.azure.FutureAPI
1. SnapshotsGrantAccessFuture.Result
1. SnapshotsGrantAccessFuture.azure.FutureAPI
1. SnapshotsRevokeAccessFuture.Result
1. SnapshotsRevokeAccessFuture.azure.FutureAPI
1. SnapshotsUpdateFuture.Result
1. SnapshotsUpdateFuture.azure.FutureAPI
1. VirtualMachineExtensionsCreateOrUpdateFuture.Result
1. VirtualMachineExtensionsCreateOrUpdateFuture.azure.FutureAPI
1. VirtualMachineExtensionsDeleteFuture.Result
1. VirtualMachineExtensionsDeleteFuture.azure.FutureAPI
1. VirtualMachineExtensionsUpdateFuture.Result
1. VirtualMachineExtensionsUpdateFuture.azure.FutureAPI
1. VirtualMachineScaleSetExtensionsCreateOrUpdateFuture.Result
1. VirtualMachineScaleSetExtensionsCreateOrUpdateFuture.azure.FutureAPI
1. VirtualMachineScaleSetExtensionsDeleteFuture.Result
1. VirtualMachineScaleSetExtensionsDeleteFuture.azure.FutureAPI
1. VirtualMachineScaleSetExtensionsUpdateFuture.Result
1. VirtualMachineScaleSetExtensionsUpdateFuture.azure.FutureAPI
1. VirtualMachineScaleSetRollingUpgradesCancelFuture.Result
1. VirtualMachineScaleSetRollingUpgradesCancelFuture.azure.FutureAPI
1. VirtualMachineScaleSetRollingUpgradesStartExtensionUpgradeFuture.Result
1. VirtualMachineScaleSetRollingUpgradesStartExtensionUpgradeFuture.azure.FutureAPI
1. VirtualMachineScaleSetRollingUpgradesStartOSUpgradeFuture.Result
1. VirtualMachineScaleSetRollingUpgradesStartOSUpgradeFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMExtensionsCreateOrUpdateFuture.Result
1. VirtualMachineScaleSetVMExtensionsCreateOrUpdateFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMExtensionsDeleteFuture.Result
1. VirtualMachineScaleSetVMExtensionsDeleteFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMExtensionsUpdateFuture.Result
1. VirtualMachineScaleSetVMExtensionsUpdateFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsDeallocateFuture.Result
1. VirtualMachineScaleSetVMsDeallocateFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsDeleteFuture.Result
1. VirtualMachineScaleSetVMsDeleteFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsPerformMaintenanceFuture.Result
1. VirtualMachineScaleSetVMsPerformMaintenanceFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsPowerOffFuture.Result
1. VirtualMachineScaleSetVMsPowerOffFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsRedeployFuture.Result
1. VirtualMachineScaleSetVMsRedeployFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsReimageAllFuture.Result
1. VirtualMachineScaleSetVMsReimageAllFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsReimageFuture.Result
1. VirtualMachineScaleSetVMsReimageFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsRestartFuture.Result
1. VirtualMachineScaleSetVMsRestartFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsRunCommandFuture.Result
1. VirtualMachineScaleSetVMsRunCommandFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsStartFuture.Result
1. VirtualMachineScaleSetVMsStartFuture.azure.FutureAPI
1. VirtualMachineScaleSetVMsUpdateFuture.Result
1. VirtualMachineScaleSetVMsUpdateFuture.azure.FutureAPI
1. VirtualMachineScaleSetsCreateOrUpdateFuture.Result
1. VirtualMachineScaleSetsCreateOrUpdateFuture.azure.FutureAPI
1. VirtualMachineScaleSetsDeallocateFuture.Result
1. VirtualMachineScaleSetsDeallocateFuture.azure.FutureAPI
1. VirtualMachineScaleSetsDeleteFuture.Result
1. VirtualMachineScaleSetsDeleteFuture.azure.FutureAPI
1. VirtualMachineScaleSetsDeleteInstancesFuture.Result
1. VirtualMachineScaleSetsDeleteInstancesFuture.azure.FutureAPI
1. VirtualMachineScaleSetsPerformMaintenanceFuture.Result
1. VirtualMachineScaleSetsPerformMaintenanceFuture.azure.FutureAPI
1. VirtualMachineScaleSetsPowerOffFuture.Result
1. VirtualMachineScaleSetsPowerOffFuture.azure.FutureAPI
1. VirtualMachineScaleSetsRedeployFuture.Result
1. VirtualMachineScaleSetsRedeployFuture.azure.FutureAPI
1. VirtualMachineScaleSetsReimageAllFuture.Result
1. VirtualMachineScaleSetsReimageAllFuture.azure.FutureAPI
1. VirtualMachineScaleSetsReimageFuture.Result
1. VirtualMachineScaleSetsReimageFuture.azure.FutureAPI
1. VirtualMachineScaleSetsRestartFuture.Result
1. VirtualMachineScaleSetsRestartFuture.azure.FutureAPI
1. VirtualMachineScaleSetsStartFuture.Result
1. VirtualMachineScaleSetsStartFuture.azure.FutureAPI
1. VirtualMachineScaleSetsUpdateFuture.Result
1. VirtualMachineScaleSetsUpdateFuture.azure.FutureAPI
1. VirtualMachineScaleSetsUpdateInstancesFuture.Result
1. VirtualMachineScaleSetsUpdateInstancesFuture.azure.FutureAPI
1. VirtualMachinesCaptureFuture.Result
1. VirtualMachinesCaptureFuture.azure.FutureAPI
1. VirtualMachinesConvertToManagedDisksFuture.Result
1. VirtualMachinesConvertToManagedDisksFuture.azure.FutureAPI
1. VirtualMachinesCreateOrUpdateFuture.Result
1. VirtualMachinesCreateOrUpdateFuture.azure.FutureAPI
1. VirtualMachinesDeallocateFuture.Result
1. VirtualMachinesDeallocateFuture.azure.FutureAPI
1. VirtualMachinesDeleteFuture.Result
1. VirtualMachinesDeleteFuture.azure.FutureAPI
1. VirtualMachinesPerformMaintenanceFuture.Result
1. VirtualMachinesPerformMaintenanceFuture.azure.FutureAPI
1. VirtualMachinesPowerOffFuture.Result
1. VirtualMachinesPowerOffFuture.azure.FutureAPI
1. VirtualMachinesReapplyFuture.Result
1. VirtualMachinesReapplyFuture.azure.FutureAPI
1. VirtualMachinesRedeployFuture.Result
1. VirtualMachinesRedeployFuture.azure.FutureAPI
1. VirtualMachinesReimageFuture.Result
1. VirtualMachinesReimageFuture.azure.FutureAPI
1. VirtualMachinesRestartFuture.Result
1. VirtualMachinesRestartFuture.azure.FutureAPI
1. VirtualMachinesRunCommandFuture.Result
1. VirtualMachinesRunCommandFuture.azure.FutureAPI
1. VirtualMachinesStartFuture.Result
1. VirtualMachinesStartFuture.azure.FutureAPI
1. VirtualMachinesUpdateFuture.Result
1. VirtualMachinesUpdateFuture.azure.FutureAPI
