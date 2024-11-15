package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Windows10GeneralConfiguration this topic provides descriptions of the declared methods, properties and relationships exposed by the windows10GeneralConfiguration resource.
type Windows10GeneralConfiguration struct {
    DeviceConfiguration
}
// NewWindows10GeneralConfiguration instantiates a new Windows10GeneralConfiguration and sets the default values.
func NewWindows10GeneralConfiguration()(*Windows10GeneralConfiguration) {
    m := &Windows10GeneralConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.windows10GeneralConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindows10GeneralConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindows10GeneralConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindows10GeneralConfiguration(), nil
}
// GetAccountsBlockAddingNonMicrosoftAccountEmail gets the accountsBlockAddingNonMicrosoftAccountEmail property value. Indicates whether or not to Block the user from adding email accounts to the device that are not associated with a Microsoft account.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetAccountsBlockAddingNonMicrosoftAccountEmail()(*bool) {
    val, err := m.GetBackingStore().Get("accountsBlockAddingNonMicrosoftAccountEmail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAntiTheftModeBlocked gets the antiTheftModeBlocked property value. Indicates whether or not to block the user from selecting an AntiTheft mode preference (Windows 10 Mobile only).
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetAntiTheftModeBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("antiTheftModeBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppsAllowTrustedAppsSideloading gets the appsAllowTrustedAppsSideloading property value. State Management Setting.
// returns a *StateManagementSetting when successful
func (m *Windows10GeneralConfiguration) GetAppsAllowTrustedAppsSideloading()(*StateManagementSetting) {
    val, err := m.GetBackingStore().Get("appsAllowTrustedAppsSideloading")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*StateManagementSetting)
    }
    return nil
}
// GetAppsBlockWindowsStoreOriginatedApps gets the appsBlockWindowsStoreOriginatedApps property value. Indicates whether or not to disable the launch of all apps from Windows Store that came pre-installed or were downloaded.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetAppsBlockWindowsStoreOriginatedApps()(*bool) {
    val, err := m.GetBackingStore().Get("appsBlockWindowsStoreOriginatedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBluetoothAllowedServices gets the bluetoothAllowedServices property value. Specify a list of allowed Bluetooth services and profiles in hex formatted strings.
// returns a []string when successful
func (m *Windows10GeneralConfiguration) GetBluetoothAllowedServices()([]string) {
    val, err := m.GetBackingStore().Get("bluetoothAllowedServices")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetBluetoothBlockAdvertising gets the bluetoothBlockAdvertising property value. Whether or not to Block the user from using bluetooth advertising.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetBluetoothBlockAdvertising()(*bool) {
    val, err := m.GetBackingStore().Get("bluetoothBlockAdvertising")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBluetoothBlockDiscoverableMode gets the bluetoothBlockDiscoverableMode property value. Whether or not to Block the user from using bluetooth discoverable mode.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetBluetoothBlockDiscoverableMode()(*bool) {
    val, err := m.GetBackingStore().Get("bluetoothBlockDiscoverableMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBluetoothBlocked gets the bluetoothBlocked property value. Whether or not to Block the user from using bluetooth.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetBluetoothBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("bluetoothBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBluetoothBlockPrePairing gets the bluetoothBlockPrePairing property value. Whether or not to block specific bundled Bluetooth peripherals to automatically pair with the host device.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetBluetoothBlockPrePairing()(*bool) {
    val, err := m.GetBackingStore().Get("bluetoothBlockPrePairing")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCameraBlocked gets the cameraBlocked property value. Whether or not to Block the user from accessing the camera of the device.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetCameraBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("cameraBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockDataWhenRoaming gets the cellularBlockDataWhenRoaming property value. Whether or not to Block the user from using data over cellular while roaming.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetCellularBlockDataWhenRoaming()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockDataWhenRoaming")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockVpn gets the cellularBlockVpn property value. Whether or not to Block the user from using VPN over cellular.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetCellularBlockVpn()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockVpn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockVpnWhenRoaming gets the cellularBlockVpnWhenRoaming property value. Whether or not to Block the user from using VPN when roaming over cellular.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetCellularBlockVpnWhenRoaming()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockVpnWhenRoaming")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCertificatesBlockManualRootCertificateInstallation gets the certificatesBlockManualRootCertificateInstallation property value. Whether or not to Block the user from doing manual root certificate installation.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetCertificatesBlockManualRootCertificateInstallation()(*bool) {
    val, err := m.GetBackingStore().Get("certificatesBlockManualRootCertificateInstallation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetConnectedDevicesServiceBlocked gets the connectedDevicesServiceBlocked property value. Whether or not to block Connected Devices Service which enables discovery and connection to other devices, remote messaging, remote app sessions and other cross-device experiences.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetConnectedDevicesServiceBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("connectedDevicesServiceBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCopyPasteBlocked gets the copyPasteBlocked property value. Whether or not to Block the user from using copy paste.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetCopyPasteBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("copyPasteBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCortanaBlocked gets the cortanaBlocked property value. Whether or not to Block the user from using Cortana.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetCortanaBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("cortanaBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderBlockEndUserAccess gets the defenderBlockEndUserAccess property value. Whether or not to block end user access to Defender.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderBlockEndUserAccess()(*bool) {
    val, err := m.GetBackingStore().Get("defenderBlockEndUserAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderCloudBlockLevel gets the defenderCloudBlockLevel property value. Possible values of Cloud Block Level
// returns a *DefenderCloudBlockLevelType when successful
func (m *Windows10GeneralConfiguration) GetDefenderCloudBlockLevel()(*DefenderCloudBlockLevelType) {
    val, err := m.GetBackingStore().Get("defenderCloudBlockLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DefenderCloudBlockLevelType)
    }
    return nil
}
// GetDefenderDaysBeforeDeletingQuarantinedMalware gets the defenderDaysBeforeDeletingQuarantinedMalware property value. Number of days before deleting quarantined malware. Valid values 0 to 90
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetDefenderDaysBeforeDeletingQuarantinedMalware()(*int32) {
    val, err := m.GetBackingStore().Get("defenderDaysBeforeDeletingQuarantinedMalware")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDefenderDetectedMalwareActions gets the defenderDetectedMalwareActions property value. Gets or sets Defender’s actions to take on detected Malware per threat level.
// returns a DefenderDetectedMalwareActionsable when successful
func (m *Windows10GeneralConfiguration) GetDefenderDetectedMalwareActions()(DefenderDetectedMalwareActionsable) {
    val, err := m.GetBackingStore().Get("defenderDetectedMalwareActions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DefenderDetectedMalwareActionsable)
    }
    return nil
}
// GetDefenderFileExtensionsToExclude gets the defenderFileExtensionsToExclude property value. File extensions to exclude from scans and real time protection.
// returns a []string when successful
func (m *Windows10GeneralConfiguration) GetDefenderFileExtensionsToExclude()([]string) {
    val, err := m.GetBackingStore().Get("defenderFileExtensionsToExclude")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDefenderFilesAndFoldersToExclude gets the defenderFilesAndFoldersToExclude property value. Files and folder to exclude from scans and real time protection.
// returns a []string when successful
func (m *Windows10GeneralConfiguration) GetDefenderFilesAndFoldersToExclude()([]string) {
    val, err := m.GetBackingStore().Get("defenderFilesAndFoldersToExclude")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDefenderMonitorFileActivity gets the defenderMonitorFileActivity property value. Possible values for monitoring file activity.
// returns a *DefenderMonitorFileActivity when successful
func (m *Windows10GeneralConfiguration) GetDefenderMonitorFileActivity()(*DefenderMonitorFileActivity) {
    val, err := m.GetBackingStore().Get("defenderMonitorFileActivity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DefenderMonitorFileActivity)
    }
    return nil
}
// GetDefenderProcessesToExclude gets the defenderProcessesToExclude property value. Processes to exclude from scans and real time protection.
// returns a []string when successful
func (m *Windows10GeneralConfiguration) GetDefenderProcessesToExclude()([]string) {
    val, err := m.GetBackingStore().Get("defenderProcessesToExclude")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetDefenderPromptForSampleSubmission gets the defenderPromptForSampleSubmission property value. Possible values for prompting user for samples submission.
// returns a *DefenderPromptForSampleSubmission when successful
func (m *Windows10GeneralConfiguration) GetDefenderPromptForSampleSubmission()(*DefenderPromptForSampleSubmission) {
    val, err := m.GetBackingStore().Get("defenderPromptForSampleSubmission")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DefenderPromptForSampleSubmission)
    }
    return nil
}
// GetDefenderRequireBehaviorMonitoring gets the defenderRequireBehaviorMonitoring property value. Indicates whether or not to require behavior monitoring.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderRequireBehaviorMonitoring()(*bool) {
    val, err := m.GetBackingStore().Get("defenderRequireBehaviorMonitoring")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderRequireCloudProtection gets the defenderRequireCloudProtection property value. Indicates whether or not to require cloud protection.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderRequireCloudProtection()(*bool) {
    val, err := m.GetBackingStore().Get("defenderRequireCloudProtection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderRequireNetworkInspectionSystem gets the defenderRequireNetworkInspectionSystem property value. Indicates whether or not to require network inspection system.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderRequireNetworkInspectionSystem()(*bool) {
    val, err := m.GetBackingStore().Get("defenderRequireNetworkInspectionSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderRequireRealTimeMonitoring gets the defenderRequireRealTimeMonitoring property value. Indicates whether or not to require real time monitoring.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderRequireRealTimeMonitoring()(*bool) {
    val, err := m.GetBackingStore().Get("defenderRequireRealTimeMonitoring")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderScanArchiveFiles gets the defenderScanArchiveFiles property value. Indicates whether or not to scan archive files.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderScanArchiveFiles()(*bool) {
    val, err := m.GetBackingStore().Get("defenderScanArchiveFiles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderScanDownloads gets the defenderScanDownloads property value. Indicates whether or not to scan downloads.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderScanDownloads()(*bool) {
    val, err := m.GetBackingStore().Get("defenderScanDownloads")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderScanIncomingMail gets the defenderScanIncomingMail property value. Indicates whether or not to scan incoming mail messages.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderScanIncomingMail()(*bool) {
    val, err := m.GetBackingStore().Get("defenderScanIncomingMail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderScanMappedNetworkDrivesDuringFullScan gets the defenderScanMappedNetworkDrivesDuringFullScan property value. Indicates whether or not to scan mapped network drives during full scan.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderScanMappedNetworkDrivesDuringFullScan()(*bool) {
    val, err := m.GetBackingStore().Get("defenderScanMappedNetworkDrivesDuringFullScan")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderScanMaxCpu gets the defenderScanMaxCpu property value. Max CPU usage percentage during scan. Valid values 0 to 100
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetDefenderScanMaxCpu()(*int32) {
    val, err := m.GetBackingStore().Get("defenderScanMaxCpu")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDefenderScanNetworkFiles gets the defenderScanNetworkFiles property value. Indicates whether or not to scan files opened from a network folder.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderScanNetworkFiles()(*bool) {
    val, err := m.GetBackingStore().Get("defenderScanNetworkFiles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderScanRemovableDrivesDuringFullScan gets the defenderScanRemovableDrivesDuringFullScan property value. Indicates whether or not to scan removable drives during full scan.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderScanRemovableDrivesDuringFullScan()(*bool) {
    val, err := m.GetBackingStore().Get("defenderScanRemovableDrivesDuringFullScan")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderScanScriptsLoadedInInternetExplorer gets the defenderScanScriptsLoadedInInternetExplorer property value. Indicates whether or not to scan scripts loaded in Internet Explorer browser.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDefenderScanScriptsLoadedInInternetExplorer()(*bool) {
    val, err := m.GetBackingStore().Get("defenderScanScriptsLoadedInInternetExplorer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefenderScanType gets the defenderScanType property value. Possible values for system scan type.
// returns a *DefenderScanType when successful
func (m *Windows10GeneralConfiguration) GetDefenderScanType()(*DefenderScanType) {
    val, err := m.GetBackingStore().Get("defenderScanType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DefenderScanType)
    }
    return nil
}
// GetDefenderScheduledQuickScanTime gets the defenderScheduledQuickScanTime property value. The time to perform a daily quick scan.
// returns a *TimeOnly when successful
func (m *Windows10GeneralConfiguration) GetDefenderScheduledQuickScanTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly) {
    val, err := m.GetBackingStore().Get("defenderScheduledQuickScanTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)
    }
    return nil
}
// GetDefenderScheduledScanTime gets the defenderScheduledScanTime property value. The defender time for the system scan.
// returns a *TimeOnly when successful
func (m *Windows10GeneralConfiguration) GetDefenderScheduledScanTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly) {
    val, err := m.GetBackingStore().Get("defenderScheduledScanTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)
    }
    return nil
}
// GetDefenderSignatureUpdateIntervalInHours gets the defenderSignatureUpdateIntervalInHours property value. The signature update interval in hours. Specify 0 not to check. Valid values 0 to 24
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetDefenderSignatureUpdateIntervalInHours()(*int32) {
    val, err := m.GetBackingStore().Get("defenderSignatureUpdateIntervalInHours")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDefenderSystemScanSchedule gets the defenderSystemScanSchedule property value. Possible values for a weekly schedule.
// returns a *WeeklySchedule when successful
func (m *Windows10GeneralConfiguration) GetDefenderSystemScanSchedule()(*WeeklySchedule) {
    val, err := m.GetBackingStore().Get("defenderSystemScanSchedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WeeklySchedule)
    }
    return nil
}
// GetDeveloperUnlockSetting gets the developerUnlockSetting property value. State Management Setting.
// returns a *StateManagementSetting when successful
func (m *Windows10GeneralConfiguration) GetDeveloperUnlockSetting()(*StateManagementSetting) {
    val, err := m.GetBackingStore().Get("developerUnlockSetting")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*StateManagementSetting)
    }
    return nil
}
// GetDeviceManagementBlockFactoryResetOnMobile gets the deviceManagementBlockFactoryResetOnMobile property value. Indicates whether or not to Block the user from resetting their phone.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDeviceManagementBlockFactoryResetOnMobile()(*bool) {
    val, err := m.GetBackingStore().Get("deviceManagementBlockFactoryResetOnMobile")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDeviceManagementBlockManualUnenroll gets the deviceManagementBlockManualUnenroll property value. Indicates whether or not to Block the user from doing manual un-enrollment from device management.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetDeviceManagementBlockManualUnenroll()(*bool) {
    val, err := m.GetBackingStore().Get("deviceManagementBlockManualUnenroll")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDiagnosticsDataSubmissionMode gets the diagnosticsDataSubmissionMode property value. Allow the device to send diagnostic and usage telemetry data, such as Watson.
// returns a *DiagnosticDataSubmissionMode when successful
func (m *Windows10GeneralConfiguration) GetDiagnosticsDataSubmissionMode()(*DiagnosticDataSubmissionMode) {
    val, err := m.GetBackingStore().Get("diagnosticsDataSubmissionMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DiagnosticDataSubmissionMode)
    }
    return nil
}
// GetEdgeAllowStartPagesModification gets the edgeAllowStartPagesModification property value. Allow users to change Start pages on Edge. Use the EdgeHomepageUrls to specify the Start pages that the user would see by default when they open Edge.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeAllowStartPagesModification()(*bool) {
    val, err := m.GetBackingStore().Get("edgeAllowStartPagesModification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockAccessToAboutFlags gets the edgeBlockAccessToAboutFlags property value. Indicates whether or not to prevent access to about flags on Edge browser.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockAccessToAboutFlags()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockAccessToAboutFlags")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockAddressBarDropdown gets the edgeBlockAddressBarDropdown property value. Block the address bar dropdown functionality in Microsoft Edge. Disable this settings to minimize network connections from Microsoft Edge to Microsoft services.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockAddressBarDropdown()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockAddressBarDropdown")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockAutofill gets the edgeBlockAutofill property value. Indicates whether or not to block auto fill.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockAutofill()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockAutofill")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockCompatibilityList gets the edgeBlockCompatibilityList property value. Block Microsoft compatibility list in Microsoft Edge. This list from Microsoft helps Edge properly display sites with known compatibility issues.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockCompatibilityList()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockCompatibilityList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockDeveloperTools gets the edgeBlockDeveloperTools property value. Indicates whether or not to block developer tools in the Edge browser.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockDeveloperTools()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockDeveloperTools")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlocked gets the edgeBlocked property value. Indicates whether or not to Block the user from using the Edge browser.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockExtensions gets the edgeBlockExtensions property value. Indicates whether or not to block extensions in the Edge browser.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockExtensions()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockExtensions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockInPrivateBrowsing gets the edgeBlockInPrivateBrowsing property value. Indicates whether or not to block InPrivate browsing on corporate networks, in the Edge browser.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockInPrivateBrowsing()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockInPrivateBrowsing")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockJavaScript gets the edgeBlockJavaScript property value. Indicates whether or not to Block the user from using JavaScript.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockJavaScript()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockJavaScript")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockLiveTileDataCollection gets the edgeBlockLiveTileDataCollection property value. Block the collection of information by Microsoft for live tile creation when users pin a site to Start from Microsoft Edge.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockLiveTileDataCollection()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockLiveTileDataCollection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockPasswordManager gets the edgeBlockPasswordManager property value. Indicates whether or not to Block password manager.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockPasswordManager()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockPasswordManager")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockPopups gets the edgeBlockPopups property value. Indicates whether or not to block popups.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockPopups()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockPopups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockSearchSuggestions gets the edgeBlockSearchSuggestions property value. Indicates whether or not to block the user from using the search suggestions in the address bar.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockSearchSuggestions()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockSearchSuggestions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockSendingDoNotTrackHeader gets the edgeBlockSendingDoNotTrackHeader property value. Indicates whether or not to Block the user from sending the do not track header.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockSendingDoNotTrackHeader()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockSendingDoNotTrackHeader")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeBlockSendingIntranetTrafficToInternetExplorer gets the edgeBlockSendingIntranetTrafficToInternetExplorer property value. Indicates whether or not to switch the intranet traffic from Edge to Internet Explorer. Note: the name of this property is misleading; the property is obsolete, use EdgeSendIntranetTrafficToInternetExplorer instead.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeBlockSendingIntranetTrafficToInternetExplorer()(*bool) {
    val, err := m.GetBackingStore().Get("edgeBlockSendingIntranetTrafficToInternetExplorer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeClearBrowsingDataOnExit gets the edgeClearBrowsingDataOnExit property value. Clear browsing data on exiting Microsoft Edge.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeClearBrowsingDataOnExit()(*bool) {
    val, err := m.GetBackingStore().Get("edgeClearBrowsingDataOnExit")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeCookiePolicy gets the edgeCookiePolicy property value. Possible values to specify which cookies are allowed in Microsoft Edge.
// returns a *EdgeCookiePolicy when successful
func (m *Windows10GeneralConfiguration) GetEdgeCookiePolicy()(*EdgeCookiePolicy) {
    val, err := m.GetBackingStore().Get("edgeCookiePolicy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EdgeCookiePolicy)
    }
    return nil
}
// GetEdgeDisableFirstRunPage gets the edgeDisableFirstRunPage property value. Block the Microsoft web page that opens on the first use of Microsoft Edge. This policy allows enterprises, like those enrolled in zero emissions configurations, to block this page.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeDisableFirstRunPage()(*bool) {
    val, err := m.GetBackingStore().Get("edgeDisableFirstRunPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeEnterpriseModeSiteListLocation gets the edgeEnterpriseModeSiteListLocation property value. Indicates the enterprise mode site list location. Could be a local file, local network or http location.
// returns a *string when successful
func (m *Windows10GeneralConfiguration) GetEdgeEnterpriseModeSiteListLocation()(*string) {
    val, err := m.GetBackingStore().Get("edgeEnterpriseModeSiteListLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEdgeFirstRunUrl gets the edgeFirstRunUrl property value. The first run URL for when Edge browser is opened for the first time.
// returns a *string when successful
func (m *Windows10GeneralConfiguration) GetEdgeFirstRunUrl()(*string) {
    val, err := m.GetBackingStore().Get("edgeFirstRunUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEdgeHomepageUrls gets the edgeHomepageUrls property value. The list of URLs for homepages shodwn on MDM-enrolled devices on Edge browser.
// returns a []string when successful
func (m *Windows10GeneralConfiguration) GetEdgeHomepageUrls()([]string) {
    val, err := m.GetBackingStore().Get("edgeHomepageUrls")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetEdgeRequireSmartScreen gets the edgeRequireSmartScreen property value. Indicates whether or not to Require the user to use the smart screen filter.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeRequireSmartScreen()(*bool) {
    val, err := m.GetBackingStore().Get("edgeRequireSmartScreen")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeSearchEngine gets the edgeSearchEngine property value. Allows IT admins to set a default search engine for MDM-Controlled devices. Users can override this and change their default search engine provided the AllowSearchEngineCustomization policy is not set.
// returns a EdgeSearchEngineBaseable when successful
func (m *Windows10GeneralConfiguration) GetEdgeSearchEngine()(EdgeSearchEngineBaseable) {
    val, err := m.GetBackingStore().Get("edgeSearchEngine")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EdgeSearchEngineBaseable)
    }
    return nil
}
// GetEdgeSendIntranetTrafficToInternetExplorer gets the edgeSendIntranetTrafficToInternetExplorer property value. Indicates whether or not to switch the intranet traffic from Edge to Internet Explorer.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeSendIntranetTrafficToInternetExplorer()(*bool) {
    val, err := m.GetBackingStore().Get("edgeSendIntranetTrafficToInternetExplorer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEdgeSyncFavoritesWithInternetExplorer gets the edgeSyncFavoritesWithInternetExplorer property value. Enable favorites sync between Internet Explorer and Microsoft Edge. Additions, deletions, modifications and order changes to favorites are shared between browsers.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetEdgeSyncFavoritesWithInternetExplorer()(*bool) {
    val, err := m.GetBackingStore().Get("edgeSyncFavoritesWithInternetExplorer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEnterpriseCloudPrintDiscoveryEndPoint gets the enterpriseCloudPrintDiscoveryEndPoint property value. Endpoint for discovering cloud printers.
// returns a *string when successful
func (m *Windows10GeneralConfiguration) GetEnterpriseCloudPrintDiscoveryEndPoint()(*string) {
    val, err := m.GetBackingStore().Get("enterpriseCloudPrintDiscoveryEndPoint")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnterpriseCloudPrintDiscoveryMaxLimit gets the enterpriseCloudPrintDiscoveryMaxLimit property value. Maximum number of printers that should be queried from a discovery endpoint. This is a mobile only setting. Valid values 1 to 65535
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetEnterpriseCloudPrintDiscoveryMaxLimit()(*int32) {
    val, err := m.GetBackingStore().Get("enterpriseCloudPrintDiscoveryMaxLimit")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetEnterpriseCloudPrintMopriaDiscoveryResourceIdentifier gets the enterpriseCloudPrintMopriaDiscoveryResourceIdentifier property value. OAuth resource URI for printer discovery service as configured in Azure portal.
// returns a *string when successful
func (m *Windows10GeneralConfiguration) GetEnterpriseCloudPrintMopriaDiscoveryResourceIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("enterpriseCloudPrintMopriaDiscoveryResourceIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnterpriseCloudPrintOAuthAuthority gets the enterpriseCloudPrintOAuthAuthority property value. Authentication endpoint for acquiring OAuth tokens.
// returns a *string when successful
func (m *Windows10GeneralConfiguration) GetEnterpriseCloudPrintOAuthAuthority()(*string) {
    val, err := m.GetBackingStore().Get("enterpriseCloudPrintOAuthAuthority")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnterpriseCloudPrintOAuthClientIdentifier gets the enterpriseCloudPrintOAuthClientIdentifier property value. GUID of a client application authorized to retrieve OAuth tokens from the OAuth Authority.
// returns a *string when successful
func (m *Windows10GeneralConfiguration) GetEnterpriseCloudPrintOAuthClientIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("enterpriseCloudPrintOAuthClientIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEnterpriseCloudPrintResourceIdentifier gets the enterpriseCloudPrintResourceIdentifier property value. OAuth resource URI for print service as configured in the Azure portal.
// returns a *string when successful
func (m *Windows10GeneralConfiguration) GetEnterpriseCloudPrintResourceIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("enterpriseCloudPrintResourceIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExperienceBlockDeviceDiscovery gets the experienceBlockDeviceDiscovery property value. Indicates whether or not to enable device discovery UX.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetExperienceBlockDeviceDiscovery()(*bool) {
    val, err := m.GetBackingStore().Get("experienceBlockDeviceDiscovery")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetExperienceBlockErrorDialogWhenNoSIM gets the experienceBlockErrorDialogWhenNoSIM property value. Indicates whether or not to allow the error dialog from displaying if no SIM card is detected.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetExperienceBlockErrorDialogWhenNoSIM()(*bool) {
    val, err := m.GetBackingStore().Get("experienceBlockErrorDialogWhenNoSIM")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetExperienceBlockTaskSwitcher gets the experienceBlockTaskSwitcher property value. Indicates whether or not to enable task switching on the device.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetExperienceBlockTaskSwitcher()(*bool) {
    val, err := m.GetBackingStore().Get("experienceBlockTaskSwitcher")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Windows10GeneralConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["accountsBlockAddingNonMicrosoftAccountEmail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccountsBlockAddingNonMicrosoftAccountEmail(val)
        }
        return nil
    }
    res["antiTheftModeBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAntiTheftModeBlocked(val)
        }
        return nil
    }
    res["appsAllowTrustedAppsSideloading"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseStateManagementSetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppsAllowTrustedAppsSideloading(val.(*StateManagementSetting))
        }
        return nil
    }
    res["appsBlockWindowsStoreOriginatedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppsBlockWindowsStoreOriginatedApps(val)
        }
        return nil
    }
    res["bluetoothAllowedServices"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetBluetoothAllowedServices(res)
        }
        return nil
    }
    res["bluetoothBlockAdvertising"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBluetoothBlockAdvertising(val)
        }
        return nil
    }
    res["bluetoothBlockDiscoverableMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBluetoothBlockDiscoverableMode(val)
        }
        return nil
    }
    res["bluetoothBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBluetoothBlocked(val)
        }
        return nil
    }
    res["bluetoothBlockPrePairing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBluetoothBlockPrePairing(val)
        }
        return nil
    }
    res["cameraBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCameraBlocked(val)
        }
        return nil
    }
    res["cellularBlockDataWhenRoaming"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockDataWhenRoaming(val)
        }
        return nil
    }
    res["cellularBlockVpn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockVpn(val)
        }
        return nil
    }
    res["cellularBlockVpnWhenRoaming"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockVpnWhenRoaming(val)
        }
        return nil
    }
    res["certificatesBlockManualRootCertificateInstallation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCertificatesBlockManualRootCertificateInstallation(val)
        }
        return nil
    }
    res["connectedDevicesServiceBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConnectedDevicesServiceBlocked(val)
        }
        return nil
    }
    res["copyPasteBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCopyPasteBlocked(val)
        }
        return nil
    }
    res["cortanaBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCortanaBlocked(val)
        }
        return nil
    }
    res["defenderBlockEndUserAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderBlockEndUserAccess(val)
        }
        return nil
    }
    res["defenderCloudBlockLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDefenderCloudBlockLevelType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderCloudBlockLevel(val.(*DefenderCloudBlockLevelType))
        }
        return nil
    }
    res["defenderDaysBeforeDeletingQuarantinedMalware"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderDaysBeforeDeletingQuarantinedMalware(val)
        }
        return nil
    }
    res["defenderDetectedMalwareActions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDefenderDetectedMalwareActionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderDetectedMalwareActions(val.(DefenderDetectedMalwareActionsable))
        }
        return nil
    }
    res["defenderFileExtensionsToExclude"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetDefenderFileExtensionsToExclude(res)
        }
        return nil
    }
    res["defenderFilesAndFoldersToExclude"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetDefenderFilesAndFoldersToExclude(res)
        }
        return nil
    }
    res["defenderMonitorFileActivity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDefenderMonitorFileActivity)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderMonitorFileActivity(val.(*DefenderMonitorFileActivity))
        }
        return nil
    }
    res["defenderProcessesToExclude"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetDefenderProcessesToExclude(res)
        }
        return nil
    }
    res["defenderPromptForSampleSubmission"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDefenderPromptForSampleSubmission)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderPromptForSampleSubmission(val.(*DefenderPromptForSampleSubmission))
        }
        return nil
    }
    res["defenderRequireBehaviorMonitoring"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderRequireBehaviorMonitoring(val)
        }
        return nil
    }
    res["defenderRequireCloudProtection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderRequireCloudProtection(val)
        }
        return nil
    }
    res["defenderRequireNetworkInspectionSystem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderRequireNetworkInspectionSystem(val)
        }
        return nil
    }
    res["defenderRequireRealTimeMonitoring"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderRequireRealTimeMonitoring(val)
        }
        return nil
    }
    res["defenderScanArchiveFiles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScanArchiveFiles(val)
        }
        return nil
    }
    res["defenderScanDownloads"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScanDownloads(val)
        }
        return nil
    }
    res["defenderScanIncomingMail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScanIncomingMail(val)
        }
        return nil
    }
    res["defenderScanMappedNetworkDrivesDuringFullScan"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScanMappedNetworkDrivesDuringFullScan(val)
        }
        return nil
    }
    res["defenderScanMaxCpu"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScanMaxCpu(val)
        }
        return nil
    }
    res["defenderScanNetworkFiles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScanNetworkFiles(val)
        }
        return nil
    }
    res["defenderScanRemovableDrivesDuringFullScan"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScanRemovableDrivesDuringFullScan(val)
        }
        return nil
    }
    res["defenderScanScriptsLoadedInInternetExplorer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScanScriptsLoadedInInternetExplorer(val)
        }
        return nil
    }
    res["defenderScanType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDefenderScanType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScanType(val.(*DefenderScanType))
        }
        return nil
    }
    res["defenderScheduledQuickScanTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeOnlyValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScheduledQuickScanTime(val)
        }
        return nil
    }
    res["defenderScheduledScanTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeOnlyValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderScheduledScanTime(val)
        }
        return nil
    }
    res["defenderSignatureUpdateIntervalInHours"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderSignatureUpdateIntervalInHours(val)
        }
        return nil
    }
    res["defenderSystemScanSchedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWeeklySchedule)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefenderSystemScanSchedule(val.(*WeeklySchedule))
        }
        return nil
    }
    res["developerUnlockSetting"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseStateManagementSetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeveloperUnlockSetting(val.(*StateManagementSetting))
        }
        return nil
    }
    res["deviceManagementBlockFactoryResetOnMobile"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceManagementBlockFactoryResetOnMobile(val)
        }
        return nil
    }
    res["deviceManagementBlockManualUnenroll"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceManagementBlockManualUnenroll(val)
        }
        return nil
    }
    res["diagnosticsDataSubmissionMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDiagnosticDataSubmissionMode)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDiagnosticsDataSubmissionMode(val.(*DiagnosticDataSubmissionMode))
        }
        return nil
    }
    res["edgeAllowStartPagesModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeAllowStartPagesModification(val)
        }
        return nil
    }
    res["edgeBlockAccessToAboutFlags"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockAccessToAboutFlags(val)
        }
        return nil
    }
    res["edgeBlockAddressBarDropdown"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockAddressBarDropdown(val)
        }
        return nil
    }
    res["edgeBlockAutofill"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockAutofill(val)
        }
        return nil
    }
    res["edgeBlockCompatibilityList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockCompatibilityList(val)
        }
        return nil
    }
    res["edgeBlockDeveloperTools"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockDeveloperTools(val)
        }
        return nil
    }
    res["edgeBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlocked(val)
        }
        return nil
    }
    res["edgeBlockExtensions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockExtensions(val)
        }
        return nil
    }
    res["edgeBlockInPrivateBrowsing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockInPrivateBrowsing(val)
        }
        return nil
    }
    res["edgeBlockJavaScript"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockJavaScript(val)
        }
        return nil
    }
    res["edgeBlockLiveTileDataCollection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockLiveTileDataCollection(val)
        }
        return nil
    }
    res["edgeBlockPasswordManager"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockPasswordManager(val)
        }
        return nil
    }
    res["edgeBlockPopups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockPopups(val)
        }
        return nil
    }
    res["edgeBlockSearchSuggestions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockSearchSuggestions(val)
        }
        return nil
    }
    res["edgeBlockSendingDoNotTrackHeader"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockSendingDoNotTrackHeader(val)
        }
        return nil
    }
    res["edgeBlockSendingIntranetTrafficToInternetExplorer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeBlockSendingIntranetTrafficToInternetExplorer(val)
        }
        return nil
    }
    res["edgeClearBrowsingDataOnExit"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeClearBrowsingDataOnExit(val)
        }
        return nil
    }
    res["edgeCookiePolicy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEdgeCookiePolicy)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeCookiePolicy(val.(*EdgeCookiePolicy))
        }
        return nil
    }
    res["edgeDisableFirstRunPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeDisableFirstRunPage(val)
        }
        return nil
    }
    res["edgeEnterpriseModeSiteListLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeEnterpriseModeSiteListLocation(val)
        }
        return nil
    }
    res["edgeFirstRunUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeFirstRunUrl(val)
        }
        return nil
    }
    res["edgeHomepageUrls"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetEdgeHomepageUrls(res)
        }
        return nil
    }
    res["edgeRequireSmartScreen"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeRequireSmartScreen(val)
        }
        return nil
    }
    res["edgeSearchEngine"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEdgeSearchEngineBaseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeSearchEngine(val.(EdgeSearchEngineBaseable))
        }
        return nil
    }
    res["edgeSendIntranetTrafficToInternetExplorer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeSendIntranetTrafficToInternetExplorer(val)
        }
        return nil
    }
    res["edgeSyncFavoritesWithInternetExplorer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEdgeSyncFavoritesWithInternetExplorer(val)
        }
        return nil
    }
    res["enterpriseCloudPrintDiscoveryEndPoint"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseCloudPrintDiscoveryEndPoint(val)
        }
        return nil
    }
    res["enterpriseCloudPrintDiscoveryMaxLimit"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseCloudPrintDiscoveryMaxLimit(val)
        }
        return nil
    }
    res["enterpriseCloudPrintMopriaDiscoveryResourceIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseCloudPrintMopriaDiscoveryResourceIdentifier(val)
        }
        return nil
    }
    res["enterpriseCloudPrintOAuthAuthority"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseCloudPrintOAuthAuthority(val)
        }
        return nil
    }
    res["enterpriseCloudPrintOAuthClientIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseCloudPrintOAuthClientIdentifier(val)
        }
        return nil
    }
    res["enterpriseCloudPrintResourceIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseCloudPrintResourceIdentifier(val)
        }
        return nil
    }
    res["experienceBlockDeviceDiscovery"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExperienceBlockDeviceDiscovery(val)
        }
        return nil
    }
    res["experienceBlockErrorDialogWhenNoSIM"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExperienceBlockErrorDialogWhenNoSIM(val)
        }
        return nil
    }
    res["experienceBlockTaskSwitcher"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExperienceBlockTaskSwitcher(val)
        }
        return nil
    }
    res["gameDvrBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGameDvrBlocked(val)
        }
        return nil
    }
    res["internetSharingBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInternetSharingBlocked(val)
        }
        return nil
    }
    res["locationServicesBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocationServicesBlocked(val)
        }
        return nil
    }
    res["lockScreenAllowTimeoutConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLockScreenAllowTimeoutConfiguration(val)
        }
        return nil
    }
    res["lockScreenBlockActionCenterNotifications"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLockScreenBlockActionCenterNotifications(val)
        }
        return nil
    }
    res["lockScreenBlockCortana"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLockScreenBlockCortana(val)
        }
        return nil
    }
    res["lockScreenBlockToastNotifications"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLockScreenBlockToastNotifications(val)
        }
        return nil
    }
    res["lockScreenTimeoutInSeconds"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLockScreenTimeoutInSeconds(val)
        }
        return nil
    }
    res["logonBlockFastUserSwitching"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLogonBlockFastUserSwitching(val)
        }
        return nil
    }
    res["microsoftAccountBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMicrosoftAccountBlocked(val)
        }
        return nil
    }
    res["microsoftAccountBlockSettingsSync"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMicrosoftAccountBlockSettingsSync(val)
        }
        return nil
    }
    res["networkProxyApplySettingsDeviceWide"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNetworkProxyApplySettingsDeviceWide(val)
        }
        return nil
    }
    res["networkProxyAutomaticConfigurationUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNetworkProxyAutomaticConfigurationUrl(val)
        }
        return nil
    }
    res["networkProxyDisableAutoDetect"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNetworkProxyDisableAutoDetect(val)
        }
        return nil
    }
    res["networkProxyServer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWindows10NetworkProxyServerFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNetworkProxyServer(val.(Windows10NetworkProxyServerable))
        }
        return nil
    }
    res["nfcBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNfcBlocked(val)
        }
        return nil
    }
    res["oneDriveDisableFileSync"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOneDriveDisableFileSync(val)
        }
        return nil
    }
    res["passwordBlockSimple"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordBlockSimple(val)
        }
        return nil
    }
    res["passwordExpirationDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordExpirationDays(val)
        }
        return nil
    }
    res["passwordMinimumCharacterSetCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordMinimumCharacterSetCount(val)
        }
        return nil
    }
    res["passwordMinimumLength"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordMinimumLength(val)
        }
        return nil
    }
    res["passwordMinutesOfInactivityBeforeScreenTimeout"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordMinutesOfInactivityBeforeScreenTimeout(val)
        }
        return nil
    }
    res["passwordPreviousPasswordBlockCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordPreviousPasswordBlockCount(val)
        }
        return nil
    }
    res["passwordRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordRequired(val)
        }
        return nil
    }
    res["passwordRequiredType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRequiredPasswordType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordRequiredType(val.(*RequiredPasswordType))
        }
        return nil
    }
    res["passwordRequireWhenResumeFromIdleState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordRequireWhenResumeFromIdleState(val)
        }
        return nil
    }
    res["passwordSignInFailureCountBeforeFactoryReset"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordSignInFailureCountBeforeFactoryReset(val)
        }
        return nil
    }
    res["personalizationDesktopImageUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPersonalizationDesktopImageUrl(val)
        }
        return nil
    }
    res["personalizationLockScreenImageUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPersonalizationLockScreenImageUrl(val)
        }
        return nil
    }
    res["privacyAdvertisingId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseStateManagementSetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrivacyAdvertisingId(val.(*StateManagementSetting))
        }
        return nil
    }
    res["privacyAutoAcceptPairingAndConsentPrompts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrivacyAutoAcceptPairingAndConsentPrompts(val)
        }
        return nil
    }
    res["privacyBlockInputPersonalization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrivacyBlockInputPersonalization(val)
        }
        return nil
    }
    res["resetProtectionModeBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResetProtectionModeBlocked(val)
        }
        return nil
    }
    res["safeSearchFilter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSafeSearchFilterType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSafeSearchFilter(val.(*SafeSearchFilterType))
        }
        return nil
    }
    res["screenCaptureBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScreenCaptureBlocked(val)
        }
        return nil
    }
    res["searchBlockDiacritics"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSearchBlockDiacritics(val)
        }
        return nil
    }
    res["searchDisableAutoLanguageDetection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSearchDisableAutoLanguageDetection(val)
        }
        return nil
    }
    res["searchDisableIndexerBackoff"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSearchDisableIndexerBackoff(val)
        }
        return nil
    }
    res["searchDisableIndexingEncryptedItems"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSearchDisableIndexingEncryptedItems(val)
        }
        return nil
    }
    res["searchDisableIndexingRemovableDrive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSearchDisableIndexingRemovableDrive(val)
        }
        return nil
    }
    res["searchEnableAutomaticIndexSizeManangement"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSearchEnableAutomaticIndexSizeManangement(val)
        }
        return nil
    }
    res["searchEnableRemoteQueries"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSearchEnableRemoteQueries(val)
        }
        return nil
    }
    res["settingsBlockAccountsPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockAccountsPage(val)
        }
        return nil
    }
    res["settingsBlockAddProvisioningPackage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockAddProvisioningPackage(val)
        }
        return nil
    }
    res["settingsBlockAppsPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockAppsPage(val)
        }
        return nil
    }
    res["settingsBlockChangeLanguage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockChangeLanguage(val)
        }
        return nil
    }
    res["settingsBlockChangePowerSleep"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockChangePowerSleep(val)
        }
        return nil
    }
    res["settingsBlockChangeRegion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockChangeRegion(val)
        }
        return nil
    }
    res["settingsBlockChangeSystemTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockChangeSystemTime(val)
        }
        return nil
    }
    res["settingsBlockDevicesPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockDevicesPage(val)
        }
        return nil
    }
    res["settingsBlockEaseOfAccessPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockEaseOfAccessPage(val)
        }
        return nil
    }
    res["settingsBlockEditDeviceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockEditDeviceName(val)
        }
        return nil
    }
    res["settingsBlockGamingPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockGamingPage(val)
        }
        return nil
    }
    res["settingsBlockNetworkInternetPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockNetworkInternetPage(val)
        }
        return nil
    }
    res["settingsBlockPersonalizationPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockPersonalizationPage(val)
        }
        return nil
    }
    res["settingsBlockPrivacyPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockPrivacyPage(val)
        }
        return nil
    }
    res["settingsBlockRemoveProvisioningPackage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockRemoveProvisioningPackage(val)
        }
        return nil
    }
    res["settingsBlockSettingsApp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockSettingsApp(val)
        }
        return nil
    }
    res["settingsBlockSystemPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockSystemPage(val)
        }
        return nil
    }
    res["settingsBlockTimeLanguagePage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockTimeLanguagePage(val)
        }
        return nil
    }
    res["settingsBlockUpdateSecurityPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockUpdateSecurityPage(val)
        }
        return nil
    }
    res["sharedUserAppDataAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharedUserAppDataAllowed(val)
        }
        return nil
    }
    res["smartScreenBlockPromptOverride"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSmartScreenBlockPromptOverride(val)
        }
        return nil
    }
    res["smartScreenBlockPromptOverrideForFiles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSmartScreenBlockPromptOverrideForFiles(val)
        }
        return nil
    }
    res["smartScreenEnableAppInstallControl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSmartScreenEnableAppInstallControl(val)
        }
        return nil
    }
    res["startBlockUnpinningAppsFromTaskbar"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartBlockUnpinningAppsFromTaskbar(val)
        }
        return nil
    }
    res["startMenuAppListVisibility"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsStartMenuAppListVisibilityType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuAppListVisibility(val.(*WindowsStartMenuAppListVisibilityType))
        }
        return nil
    }
    res["startMenuHideChangeAccountSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideChangeAccountSettings(val)
        }
        return nil
    }
    res["startMenuHideFrequentlyUsedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideFrequentlyUsedApps(val)
        }
        return nil
    }
    res["startMenuHideHibernate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideHibernate(val)
        }
        return nil
    }
    res["startMenuHideLock"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideLock(val)
        }
        return nil
    }
    res["startMenuHidePowerButton"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHidePowerButton(val)
        }
        return nil
    }
    res["startMenuHideRecentJumpLists"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideRecentJumpLists(val)
        }
        return nil
    }
    res["startMenuHideRecentlyAddedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideRecentlyAddedApps(val)
        }
        return nil
    }
    res["startMenuHideRestartOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideRestartOptions(val)
        }
        return nil
    }
    res["startMenuHideShutDown"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideShutDown(val)
        }
        return nil
    }
    res["startMenuHideSignOut"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideSignOut(val)
        }
        return nil
    }
    res["startMenuHideSleep"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideSleep(val)
        }
        return nil
    }
    res["startMenuHideSwitchAccount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideSwitchAccount(val)
        }
        return nil
    }
    res["startMenuHideUserTile"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuHideUserTile(val)
        }
        return nil
    }
    res["startMenuLayoutEdgeAssetsXml"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuLayoutEdgeAssetsXml(val)
        }
        return nil
    }
    res["startMenuLayoutXml"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuLayoutXml(val)
        }
        return nil
    }
    res["startMenuMode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsStartMenuModeType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuMode(val.(*WindowsStartMenuModeType))
        }
        return nil
    }
    res["startMenuPinnedFolderDocuments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVisibilitySetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuPinnedFolderDocuments(val.(*VisibilitySetting))
        }
        return nil
    }
    res["startMenuPinnedFolderDownloads"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVisibilitySetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuPinnedFolderDownloads(val.(*VisibilitySetting))
        }
        return nil
    }
    res["startMenuPinnedFolderFileExplorer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVisibilitySetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuPinnedFolderFileExplorer(val.(*VisibilitySetting))
        }
        return nil
    }
    res["startMenuPinnedFolderHomeGroup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVisibilitySetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuPinnedFolderHomeGroup(val.(*VisibilitySetting))
        }
        return nil
    }
    res["startMenuPinnedFolderMusic"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVisibilitySetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuPinnedFolderMusic(val.(*VisibilitySetting))
        }
        return nil
    }
    res["startMenuPinnedFolderNetwork"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVisibilitySetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuPinnedFolderNetwork(val.(*VisibilitySetting))
        }
        return nil
    }
    res["startMenuPinnedFolderPersonalFolder"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVisibilitySetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuPinnedFolderPersonalFolder(val.(*VisibilitySetting))
        }
        return nil
    }
    res["startMenuPinnedFolderPictures"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVisibilitySetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuPinnedFolderPictures(val.(*VisibilitySetting))
        }
        return nil
    }
    res["startMenuPinnedFolderSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVisibilitySetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuPinnedFolderSettings(val.(*VisibilitySetting))
        }
        return nil
    }
    res["startMenuPinnedFolderVideos"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVisibilitySetting)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartMenuPinnedFolderVideos(val.(*VisibilitySetting))
        }
        return nil
    }
    res["storageBlockRemovableStorage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStorageBlockRemovableStorage(val)
        }
        return nil
    }
    res["storageRequireMobileDeviceEncryption"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStorageRequireMobileDeviceEncryption(val)
        }
        return nil
    }
    res["storageRestrictAppDataToSystemVolume"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStorageRestrictAppDataToSystemVolume(val)
        }
        return nil
    }
    res["storageRestrictAppInstallToSystemVolume"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStorageRestrictAppInstallToSystemVolume(val)
        }
        return nil
    }
    res["tenantLockdownRequireNetworkDuringOutOfBoxExperience"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTenantLockdownRequireNetworkDuringOutOfBoxExperience(val)
        }
        return nil
    }
    res["usbBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUsbBlocked(val)
        }
        return nil
    }
    res["voiceRecordingBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVoiceRecordingBlocked(val)
        }
        return nil
    }
    res["webRtcBlockLocalhostIpAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebRtcBlockLocalhostIpAddress(val)
        }
        return nil
    }
    res["wiFiBlockAutomaticConnectHotspots"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWiFiBlockAutomaticConnectHotspots(val)
        }
        return nil
    }
    res["wiFiBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWiFiBlocked(val)
        }
        return nil
    }
    res["wiFiBlockManualConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWiFiBlockManualConfiguration(val)
        }
        return nil
    }
    res["wiFiScanInterval"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWiFiScanInterval(val)
        }
        return nil
    }
    res["windowsSpotlightBlockConsumerSpecificFeatures"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsSpotlightBlockConsumerSpecificFeatures(val)
        }
        return nil
    }
    res["windowsSpotlightBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsSpotlightBlocked(val)
        }
        return nil
    }
    res["windowsSpotlightBlockOnActionCenter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsSpotlightBlockOnActionCenter(val)
        }
        return nil
    }
    res["windowsSpotlightBlockTailoredExperiences"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsSpotlightBlockTailoredExperiences(val)
        }
        return nil
    }
    res["windowsSpotlightBlockThirdPartyNotifications"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsSpotlightBlockThirdPartyNotifications(val)
        }
        return nil
    }
    res["windowsSpotlightBlockWelcomeExperience"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsSpotlightBlockWelcomeExperience(val)
        }
        return nil
    }
    res["windowsSpotlightBlockWindowsTips"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsSpotlightBlockWindowsTips(val)
        }
        return nil
    }
    res["windowsSpotlightConfigureOnLockScreen"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsSpotlightEnablementSettings)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsSpotlightConfigureOnLockScreen(val.(*WindowsSpotlightEnablementSettings))
        }
        return nil
    }
    res["windowsStoreBlockAutoUpdate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsStoreBlockAutoUpdate(val)
        }
        return nil
    }
    res["windowsStoreBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsStoreBlocked(val)
        }
        return nil
    }
    res["windowsStoreEnablePrivateStoreOnly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsStoreEnablePrivateStoreOnly(val)
        }
        return nil
    }
    res["wirelessDisplayBlockProjectionToThisDevice"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWirelessDisplayBlockProjectionToThisDevice(val)
        }
        return nil
    }
    res["wirelessDisplayBlockUserInputFromReceiver"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWirelessDisplayBlockUserInputFromReceiver(val)
        }
        return nil
    }
    res["wirelessDisplayRequirePinForPairing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWirelessDisplayRequirePinForPairing(val)
        }
        return nil
    }
    return res
}
// GetGameDvrBlocked gets the gameDvrBlocked property value. Indicates whether or not to block DVR and broadcasting.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetGameDvrBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("gameDvrBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetInternetSharingBlocked gets the internetSharingBlocked property value. Indicates whether or not to Block the user from using internet sharing.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetInternetSharingBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("internetSharingBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLocationServicesBlocked gets the locationServicesBlocked property value. Indicates whether or not to Block the user from location services.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetLocationServicesBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("locationServicesBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLockScreenAllowTimeoutConfiguration gets the lockScreenAllowTimeoutConfiguration property value. Specify whether to show a user-configurable setting to control the screen timeout while on the lock screen of Windows 10 Mobile devices. If this policy is set to Allow, the value set by lockScreenTimeoutInSeconds is ignored.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetLockScreenAllowTimeoutConfiguration()(*bool) {
    val, err := m.GetBackingStore().Get("lockScreenAllowTimeoutConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLockScreenBlockActionCenterNotifications gets the lockScreenBlockActionCenterNotifications property value. Indicates whether or not to block action center notifications over lock screen.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetLockScreenBlockActionCenterNotifications()(*bool) {
    val, err := m.GetBackingStore().Get("lockScreenBlockActionCenterNotifications")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLockScreenBlockCortana gets the lockScreenBlockCortana property value. Indicates whether or not the user can interact with Cortana using speech while the system is locked.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetLockScreenBlockCortana()(*bool) {
    val, err := m.GetBackingStore().Get("lockScreenBlockCortana")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLockScreenBlockToastNotifications gets the lockScreenBlockToastNotifications property value. Indicates whether to allow toast notifications above the device lock screen.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetLockScreenBlockToastNotifications()(*bool) {
    val, err := m.GetBackingStore().Get("lockScreenBlockToastNotifications")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLockScreenTimeoutInSeconds gets the lockScreenTimeoutInSeconds property value. Set the duration (in seconds) from the screen locking to the screen turning off for Windows 10 Mobile devices. Supported values are 11-1800. Valid values 11 to 1800
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetLockScreenTimeoutInSeconds()(*int32) {
    val, err := m.GetBackingStore().Get("lockScreenTimeoutInSeconds")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetLogonBlockFastUserSwitching gets the logonBlockFastUserSwitching property value. Disables the ability to quickly switch between users that are logged on simultaneously without logging off.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetLogonBlockFastUserSwitching()(*bool) {
    val, err := m.GetBackingStore().Get("logonBlockFastUserSwitching")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMicrosoftAccountBlocked gets the microsoftAccountBlocked property value. Indicates whether or not to Block a Microsoft account.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetMicrosoftAccountBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("microsoftAccountBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMicrosoftAccountBlockSettingsSync gets the microsoftAccountBlockSettingsSync property value. Indicates whether or not to Block Microsoft account settings sync.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetMicrosoftAccountBlockSettingsSync()(*bool) {
    val, err := m.GetBackingStore().Get("microsoftAccountBlockSettingsSync")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNetworkProxyApplySettingsDeviceWide gets the networkProxyApplySettingsDeviceWide property value. If set, proxy settings will be applied to all processes and accounts in the device. Otherwise, it will be applied to the user account that’s enrolled into MDM.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetNetworkProxyApplySettingsDeviceWide()(*bool) {
    val, err := m.GetBackingStore().Get("networkProxyApplySettingsDeviceWide")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNetworkProxyAutomaticConfigurationUrl gets the networkProxyAutomaticConfigurationUrl property value. Address to the proxy auto-config (PAC) script you want to use.
// returns a *string when successful
func (m *Windows10GeneralConfiguration) GetNetworkProxyAutomaticConfigurationUrl()(*string) {
    val, err := m.GetBackingStore().Get("networkProxyAutomaticConfigurationUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNetworkProxyDisableAutoDetect gets the networkProxyDisableAutoDetect property value. Disable automatic detection of settings. If enabled, the system will try to find the path to a proxy auto-config (PAC) script.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetNetworkProxyDisableAutoDetect()(*bool) {
    val, err := m.GetBackingStore().Get("networkProxyDisableAutoDetect")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNetworkProxyServer gets the networkProxyServer property value. Specifies manual proxy server settings.
// returns a Windows10NetworkProxyServerable when successful
func (m *Windows10GeneralConfiguration) GetNetworkProxyServer()(Windows10NetworkProxyServerable) {
    val, err := m.GetBackingStore().Get("networkProxyServer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Windows10NetworkProxyServerable)
    }
    return nil
}
// GetNfcBlocked gets the nfcBlocked property value. Indicates whether or not to Block the user from using near field communication.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetNfcBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("nfcBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOneDriveDisableFileSync gets the oneDriveDisableFileSync property value. Gets or sets a value allowing IT admins to prevent apps and features from working with files on OneDrive.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetOneDriveDisableFileSync()(*bool) {
    val, err := m.GetBackingStore().Get("oneDriveDisableFileSync")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordBlockSimple gets the passwordBlockSimple property value. Specify whether PINs or passwords such as '1111' or '1234' are allowed. For Windows 10 desktops, it also controls the use of picture passwords.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetPasswordBlockSimple()(*bool) {
    val, err := m.GetBackingStore().Get("passwordBlockSimple")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordExpirationDays gets the passwordExpirationDays property value. The password expiration in days. Valid values 0 to 730
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetPasswordExpirationDays()(*int32) {
    val, err := m.GetBackingStore().Get("passwordExpirationDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinimumCharacterSetCount gets the passwordMinimumCharacterSetCount property value. The number of character sets required in the password.
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetPasswordMinimumCharacterSetCount()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinimumCharacterSetCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinimumLength gets the passwordMinimumLength property value. The minimum password length. Valid values 4 to 16
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetPasswordMinimumLength()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinimumLength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinutesOfInactivityBeforeScreenTimeout gets the passwordMinutesOfInactivityBeforeScreenTimeout property value. The minutes of inactivity before the screen times out.
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetPasswordMinutesOfInactivityBeforeScreenTimeout()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinutesOfInactivityBeforeScreenTimeout")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordPreviousPasswordBlockCount gets the passwordPreviousPasswordBlockCount property value. The number of previous passwords to prevent reuse of. Valid values 0 to 50
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetPasswordPreviousPasswordBlockCount()(*int32) {
    val, err := m.GetBackingStore().Get("passwordPreviousPasswordBlockCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordRequired gets the passwordRequired property value. Indicates whether or not to require the user to have a password.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetPasswordRequired()(*bool) {
    val, err := m.GetBackingStore().Get("passwordRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordRequiredType gets the passwordRequiredType property value. Possible values of required passwords.
// returns a *RequiredPasswordType when successful
func (m *Windows10GeneralConfiguration) GetPasswordRequiredType()(*RequiredPasswordType) {
    val, err := m.GetBackingStore().Get("passwordRequiredType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RequiredPasswordType)
    }
    return nil
}
// GetPasswordRequireWhenResumeFromIdleState gets the passwordRequireWhenResumeFromIdleState property value. Indicates whether or not to require a password upon resuming from an idle state.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetPasswordRequireWhenResumeFromIdleState()(*bool) {
    val, err := m.GetBackingStore().Get("passwordRequireWhenResumeFromIdleState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordSignInFailureCountBeforeFactoryReset gets the passwordSignInFailureCountBeforeFactoryReset property value. The number of sign in failures before factory reset. Valid values 0 to 999
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetPasswordSignInFailureCountBeforeFactoryReset()(*int32) {
    val, err := m.GetBackingStore().Get("passwordSignInFailureCountBeforeFactoryReset")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPersonalizationDesktopImageUrl gets the personalizationDesktopImageUrl property value. A http or https Url to a jpg, jpeg or png image that needs to be downloaded and used as the Desktop Image or a file Url to a local image on the file system that needs to used as the Desktop Image.
// returns a *string when successful
func (m *Windows10GeneralConfiguration) GetPersonalizationDesktopImageUrl()(*string) {
    val, err := m.GetBackingStore().Get("personalizationDesktopImageUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPersonalizationLockScreenImageUrl gets the personalizationLockScreenImageUrl property value. A http or https Url to a jpg, jpeg or png image that neeeds to be downloaded and used as the Lock Screen Image or a file Url to a local image on the file system that needs to be used as the Lock Screen Image.
// returns a *string when successful
func (m *Windows10GeneralConfiguration) GetPersonalizationLockScreenImageUrl()(*string) {
    val, err := m.GetBackingStore().Get("personalizationLockScreenImageUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrivacyAdvertisingId gets the privacyAdvertisingId property value. State Management Setting.
// returns a *StateManagementSetting when successful
func (m *Windows10GeneralConfiguration) GetPrivacyAdvertisingId()(*StateManagementSetting) {
    val, err := m.GetBackingStore().Get("privacyAdvertisingId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*StateManagementSetting)
    }
    return nil
}
// GetPrivacyAutoAcceptPairingAndConsentPrompts gets the privacyAutoAcceptPairingAndConsentPrompts property value. Indicates whether or not to allow the automatic acceptance of the pairing and privacy user consent dialog when launching apps.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetPrivacyAutoAcceptPairingAndConsentPrompts()(*bool) {
    val, err := m.GetBackingStore().Get("privacyAutoAcceptPairingAndConsentPrompts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPrivacyBlockInputPersonalization gets the privacyBlockInputPersonalization property value. Indicates whether or not to block the usage of cloud based speech services for Cortana, Dictation, or Store applications.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetPrivacyBlockInputPersonalization()(*bool) {
    val, err := m.GetBackingStore().Get("privacyBlockInputPersonalization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetResetProtectionModeBlocked gets the resetProtectionModeBlocked property value. Indicates whether or not to Block the user from reset protection mode.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetResetProtectionModeBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("resetProtectionModeBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSafeSearchFilter gets the safeSearchFilter property value. Specifies what level of safe search (filtering adult content) is required
// returns a *SafeSearchFilterType when successful
func (m *Windows10GeneralConfiguration) GetSafeSearchFilter()(*SafeSearchFilterType) {
    val, err := m.GetBackingStore().Get("safeSearchFilter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SafeSearchFilterType)
    }
    return nil
}
// GetScreenCaptureBlocked gets the screenCaptureBlocked property value. Indicates whether or not to Block the user from taking Screenshots.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetScreenCaptureBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("screenCaptureBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSearchBlockDiacritics gets the searchBlockDiacritics property value. Specifies if search can use diacritics.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSearchBlockDiacritics()(*bool) {
    val, err := m.GetBackingStore().Get("searchBlockDiacritics")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSearchDisableAutoLanguageDetection gets the searchDisableAutoLanguageDetection property value. Specifies whether to use automatic language detection when indexing content and properties.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSearchDisableAutoLanguageDetection()(*bool) {
    val, err := m.GetBackingStore().Get("searchDisableAutoLanguageDetection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSearchDisableIndexerBackoff gets the searchDisableIndexerBackoff property value. Indicates whether or not to disable the search indexer backoff feature.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSearchDisableIndexerBackoff()(*bool) {
    val, err := m.GetBackingStore().Get("searchDisableIndexerBackoff")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSearchDisableIndexingEncryptedItems gets the searchDisableIndexingEncryptedItems property value. Indicates whether or not to block indexing of WIP-protected items to prevent them from appearing in search results for Cortana or Explorer.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSearchDisableIndexingEncryptedItems()(*bool) {
    val, err := m.GetBackingStore().Get("searchDisableIndexingEncryptedItems")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSearchDisableIndexingRemovableDrive gets the searchDisableIndexingRemovableDrive property value. Indicates whether or not to allow users to add locations on removable drives to libraries and to be indexed.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSearchDisableIndexingRemovableDrive()(*bool) {
    val, err := m.GetBackingStore().Get("searchDisableIndexingRemovableDrive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSearchEnableAutomaticIndexSizeManangement gets the searchEnableAutomaticIndexSizeManangement property value. Specifies minimum amount of hard drive space on the same drive as the index location before indexing stops.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSearchEnableAutomaticIndexSizeManangement()(*bool) {
    val, err := m.GetBackingStore().Get("searchEnableAutomaticIndexSizeManangement")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSearchEnableRemoteQueries gets the searchEnableRemoteQueries property value. Indicates whether or not to block remote queries of this computer’s index.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSearchEnableRemoteQueries()(*bool) {
    val, err := m.GetBackingStore().Get("searchEnableRemoteQueries")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockAccountsPage gets the settingsBlockAccountsPage property value. Indicates whether or not to block access to Accounts in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockAccountsPage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockAccountsPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockAddProvisioningPackage gets the settingsBlockAddProvisioningPackage property value. Indicates whether or not to block the user from installing provisioning packages.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockAddProvisioningPackage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockAddProvisioningPackage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockAppsPage gets the settingsBlockAppsPage property value. Indicates whether or not to block access to Apps in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockAppsPage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockAppsPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockChangeLanguage gets the settingsBlockChangeLanguage property value. Indicates whether or not to block the user from changing the language settings.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockChangeLanguage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockChangeLanguage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockChangePowerSleep gets the settingsBlockChangePowerSleep property value. Indicates whether or not to block the user from changing power and sleep settings.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockChangePowerSleep()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockChangePowerSleep")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockChangeRegion gets the settingsBlockChangeRegion property value. Indicates whether or not to block the user from changing the region settings.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockChangeRegion()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockChangeRegion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockChangeSystemTime gets the settingsBlockChangeSystemTime property value. Indicates whether or not to block the user from changing date and time settings.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockChangeSystemTime()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockChangeSystemTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockDevicesPage gets the settingsBlockDevicesPage property value. Indicates whether or not to block access to Devices in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockDevicesPage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockDevicesPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockEaseOfAccessPage gets the settingsBlockEaseOfAccessPage property value. Indicates whether or not to block access to Ease of Access in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockEaseOfAccessPage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockEaseOfAccessPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockEditDeviceName gets the settingsBlockEditDeviceName property value. Indicates whether or not to block the user from editing the device name.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockEditDeviceName()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockEditDeviceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockGamingPage gets the settingsBlockGamingPage property value. Indicates whether or not to block access to Gaming in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockGamingPage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockGamingPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockNetworkInternetPage gets the settingsBlockNetworkInternetPage property value. Indicates whether or not to block access to Network & Internet in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockNetworkInternetPage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockNetworkInternetPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockPersonalizationPage gets the settingsBlockPersonalizationPage property value. Indicates whether or not to block access to Personalization in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockPersonalizationPage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockPersonalizationPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockPrivacyPage gets the settingsBlockPrivacyPage property value. Indicates whether or not to block access to Privacy in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockPrivacyPage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockPrivacyPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockRemoveProvisioningPackage gets the settingsBlockRemoveProvisioningPackage property value. Indicates whether or not to block the runtime configuration agent from removing provisioning packages.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockRemoveProvisioningPackage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockRemoveProvisioningPackage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockSettingsApp gets the settingsBlockSettingsApp property value. Indicates whether or not to block access to Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockSettingsApp()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockSettingsApp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockSystemPage gets the settingsBlockSystemPage property value. Indicates whether or not to block access to System in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockSystemPage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockSystemPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockTimeLanguagePage gets the settingsBlockTimeLanguagePage property value. Indicates whether or not to block access to Time & Language in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockTimeLanguagePage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockTimeLanguagePage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockUpdateSecurityPage gets the settingsBlockUpdateSecurityPage property value. Indicates whether or not to block access to Update & Security in Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSettingsBlockUpdateSecurityPage()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockUpdateSecurityPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSharedUserAppDataAllowed gets the sharedUserAppDataAllowed property value. Indicates whether or not to block multiple users of the same app to share data.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSharedUserAppDataAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("sharedUserAppDataAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSmartScreenBlockPromptOverride gets the smartScreenBlockPromptOverride property value. Indicates whether or not users can override SmartScreen Filter warnings about potentially malicious websites.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSmartScreenBlockPromptOverride()(*bool) {
    val, err := m.GetBackingStore().Get("smartScreenBlockPromptOverride")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSmartScreenBlockPromptOverrideForFiles gets the smartScreenBlockPromptOverrideForFiles property value. Indicates whether or not users can override the SmartScreen Filter warnings about downloading unverified files
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSmartScreenBlockPromptOverrideForFiles()(*bool) {
    val, err := m.GetBackingStore().Get("smartScreenBlockPromptOverrideForFiles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSmartScreenEnableAppInstallControl gets the smartScreenEnableAppInstallControl property value. This property will be deprecated in July 2019 and will be replaced by property SmartScreenAppInstallControl. Allows IT Admins to control whether users are allowed to install apps from places other than the Store.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetSmartScreenEnableAppInstallControl()(*bool) {
    val, err := m.GetBackingStore().Get("smartScreenEnableAppInstallControl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartBlockUnpinningAppsFromTaskbar gets the startBlockUnpinningAppsFromTaskbar property value. Indicates whether or not to block the user from unpinning apps from taskbar.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartBlockUnpinningAppsFromTaskbar()(*bool) {
    val, err := m.GetBackingStore().Get("startBlockUnpinningAppsFromTaskbar")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuAppListVisibility gets the startMenuAppListVisibility property value. Type of start menu app list visibility.
// returns a *WindowsStartMenuAppListVisibilityType when successful
func (m *Windows10GeneralConfiguration) GetStartMenuAppListVisibility()(*WindowsStartMenuAppListVisibilityType) {
    val, err := m.GetBackingStore().Get("startMenuAppListVisibility")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsStartMenuAppListVisibilityType)
    }
    return nil
}
// GetStartMenuHideChangeAccountSettings gets the startMenuHideChangeAccountSettings property value. Enabling this policy hides the change account setting from appearing in the user tile in the start menu.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideChangeAccountSettings()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideChangeAccountSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideFrequentlyUsedApps gets the startMenuHideFrequentlyUsedApps property value. Enabling this policy hides the most used apps from appearing on the start menu and disables the corresponding toggle in the Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideFrequentlyUsedApps()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideFrequentlyUsedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideHibernate gets the startMenuHideHibernate property value. Enabling this policy hides hibernate from appearing in the power button in the start menu.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideHibernate()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideHibernate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideLock gets the startMenuHideLock property value. Enabling this policy hides lock from appearing in the user tile in the start menu.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideLock()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideLock")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHidePowerButton gets the startMenuHidePowerButton property value. Enabling this policy hides the power button from appearing in the start menu.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHidePowerButton()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHidePowerButton")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideRecentJumpLists gets the startMenuHideRecentJumpLists property value. Enabling this policy hides recent jump lists from appearing on the start menu/taskbar and disables the corresponding toggle in the Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideRecentJumpLists()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideRecentJumpLists")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideRecentlyAddedApps gets the startMenuHideRecentlyAddedApps property value. Enabling this policy hides recently added apps from appearing on the start menu and disables the corresponding toggle in the Settings app.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideRecentlyAddedApps()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideRecentlyAddedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideRestartOptions gets the startMenuHideRestartOptions property value. Enabling this policy hides 'Restart/Update and Restart' from appearing in the power button in the start menu.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideRestartOptions()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideRestartOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideShutDown gets the startMenuHideShutDown property value. Enabling this policy hides shut down/update and shut down from appearing in the power button in the start menu.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideShutDown()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideShutDown")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideSignOut gets the startMenuHideSignOut property value. Enabling this policy hides sign out from appearing in the user tile in the start menu.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideSignOut()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideSignOut")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideSleep gets the startMenuHideSleep property value. Enabling this policy hides sleep from appearing in the power button in the start menu.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideSleep()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideSleep")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideSwitchAccount gets the startMenuHideSwitchAccount property value. Enabling this policy hides switch account from appearing in the user tile in the start menu.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideSwitchAccount()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideSwitchAccount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuHideUserTile gets the startMenuHideUserTile property value. Enabling this policy hides the user tile from appearing in the start menu.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStartMenuHideUserTile()(*bool) {
    val, err := m.GetBackingStore().Get("startMenuHideUserTile")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStartMenuLayoutEdgeAssetsXml gets the startMenuLayoutEdgeAssetsXml property value. This policy setting allows you to import Edge assets to be used with startMenuLayoutXml policy. Start layout can contain secondary tile from Edge app which looks for Edge local asset file. Edge local asset would not exist and cause Edge secondary tile to appear empty in this case. This policy only gets applied when startMenuLayoutXml policy is modified. The value should be a UTF-8 Base64 encoded byte array.
// returns a []byte when successful
func (m *Windows10GeneralConfiguration) GetStartMenuLayoutEdgeAssetsXml()([]byte) {
    val, err := m.GetBackingStore().Get("startMenuLayoutEdgeAssetsXml")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetStartMenuLayoutXml gets the startMenuLayoutXml property value. Allows admins to override the default Start menu layout and prevents the user from changing it. The layout is modified by specifying an XML file based on a layout modification schema. XML needs to be in a UTF8 encoded byte array format.
// returns a []byte when successful
func (m *Windows10GeneralConfiguration) GetStartMenuLayoutXml()([]byte) {
    val, err := m.GetBackingStore().Get("startMenuLayoutXml")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetStartMenuMode gets the startMenuMode property value. Type of display modes for the start menu.
// returns a *WindowsStartMenuModeType when successful
func (m *Windows10GeneralConfiguration) GetStartMenuMode()(*WindowsStartMenuModeType) {
    val, err := m.GetBackingStore().Get("startMenuMode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsStartMenuModeType)
    }
    return nil
}
// GetStartMenuPinnedFolderDocuments gets the startMenuPinnedFolderDocuments property value. Generic visibility state.
// returns a *VisibilitySetting when successful
func (m *Windows10GeneralConfiguration) GetStartMenuPinnedFolderDocuments()(*VisibilitySetting) {
    val, err := m.GetBackingStore().Get("startMenuPinnedFolderDocuments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VisibilitySetting)
    }
    return nil
}
// GetStartMenuPinnedFolderDownloads gets the startMenuPinnedFolderDownloads property value. Generic visibility state.
// returns a *VisibilitySetting when successful
func (m *Windows10GeneralConfiguration) GetStartMenuPinnedFolderDownloads()(*VisibilitySetting) {
    val, err := m.GetBackingStore().Get("startMenuPinnedFolderDownloads")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VisibilitySetting)
    }
    return nil
}
// GetStartMenuPinnedFolderFileExplorer gets the startMenuPinnedFolderFileExplorer property value. Generic visibility state.
// returns a *VisibilitySetting when successful
func (m *Windows10GeneralConfiguration) GetStartMenuPinnedFolderFileExplorer()(*VisibilitySetting) {
    val, err := m.GetBackingStore().Get("startMenuPinnedFolderFileExplorer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VisibilitySetting)
    }
    return nil
}
// GetStartMenuPinnedFolderHomeGroup gets the startMenuPinnedFolderHomeGroup property value. Generic visibility state.
// returns a *VisibilitySetting when successful
func (m *Windows10GeneralConfiguration) GetStartMenuPinnedFolderHomeGroup()(*VisibilitySetting) {
    val, err := m.GetBackingStore().Get("startMenuPinnedFolderHomeGroup")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VisibilitySetting)
    }
    return nil
}
// GetStartMenuPinnedFolderMusic gets the startMenuPinnedFolderMusic property value. Generic visibility state.
// returns a *VisibilitySetting when successful
func (m *Windows10GeneralConfiguration) GetStartMenuPinnedFolderMusic()(*VisibilitySetting) {
    val, err := m.GetBackingStore().Get("startMenuPinnedFolderMusic")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VisibilitySetting)
    }
    return nil
}
// GetStartMenuPinnedFolderNetwork gets the startMenuPinnedFolderNetwork property value. Generic visibility state.
// returns a *VisibilitySetting when successful
func (m *Windows10GeneralConfiguration) GetStartMenuPinnedFolderNetwork()(*VisibilitySetting) {
    val, err := m.GetBackingStore().Get("startMenuPinnedFolderNetwork")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VisibilitySetting)
    }
    return nil
}
// GetStartMenuPinnedFolderPersonalFolder gets the startMenuPinnedFolderPersonalFolder property value. Generic visibility state.
// returns a *VisibilitySetting when successful
func (m *Windows10GeneralConfiguration) GetStartMenuPinnedFolderPersonalFolder()(*VisibilitySetting) {
    val, err := m.GetBackingStore().Get("startMenuPinnedFolderPersonalFolder")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VisibilitySetting)
    }
    return nil
}
// GetStartMenuPinnedFolderPictures gets the startMenuPinnedFolderPictures property value. Generic visibility state.
// returns a *VisibilitySetting when successful
func (m *Windows10GeneralConfiguration) GetStartMenuPinnedFolderPictures()(*VisibilitySetting) {
    val, err := m.GetBackingStore().Get("startMenuPinnedFolderPictures")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VisibilitySetting)
    }
    return nil
}
// GetStartMenuPinnedFolderSettings gets the startMenuPinnedFolderSettings property value. Generic visibility state.
// returns a *VisibilitySetting when successful
func (m *Windows10GeneralConfiguration) GetStartMenuPinnedFolderSettings()(*VisibilitySetting) {
    val, err := m.GetBackingStore().Get("startMenuPinnedFolderSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VisibilitySetting)
    }
    return nil
}
// GetStartMenuPinnedFolderVideos gets the startMenuPinnedFolderVideos property value. Generic visibility state.
// returns a *VisibilitySetting when successful
func (m *Windows10GeneralConfiguration) GetStartMenuPinnedFolderVideos()(*VisibilitySetting) {
    val, err := m.GetBackingStore().Get("startMenuPinnedFolderVideos")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VisibilitySetting)
    }
    return nil
}
// GetStorageBlockRemovableStorage gets the storageBlockRemovableStorage property value. Indicates whether or not to Block the user from using removable storage.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStorageBlockRemovableStorage()(*bool) {
    val, err := m.GetBackingStore().Get("storageBlockRemovableStorage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStorageRequireMobileDeviceEncryption gets the storageRequireMobileDeviceEncryption property value. Indicating whether or not to require encryption on a mobile device.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStorageRequireMobileDeviceEncryption()(*bool) {
    val, err := m.GetBackingStore().Get("storageRequireMobileDeviceEncryption")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStorageRestrictAppDataToSystemVolume gets the storageRestrictAppDataToSystemVolume property value. Indicates whether application data is restricted to the system drive.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStorageRestrictAppDataToSystemVolume()(*bool) {
    val, err := m.GetBackingStore().Get("storageRestrictAppDataToSystemVolume")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStorageRestrictAppInstallToSystemVolume gets the storageRestrictAppInstallToSystemVolume property value. Indicates whether the installation of applications is restricted to the system drive.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetStorageRestrictAppInstallToSystemVolume()(*bool) {
    val, err := m.GetBackingStore().Get("storageRestrictAppInstallToSystemVolume")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetTenantLockdownRequireNetworkDuringOutOfBoxExperience gets the tenantLockdownRequireNetworkDuringOutOfBoxExperience property value. Whether the device is required to connect to the network.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetTenantLockdownRequireNetworkDuringOutOfBoxExperience()(*bool) {
    val, err := m.GetBackingStore().Get("tenantLockdownRequireNetworkDuringOutOfBoxExperience")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUsbBlocked gets the usbBlocked property value. Indicates whether or not to Block the user from USB connection.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetUsbBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("usbBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetVoiceRecordingBlocked gets the voiceRecordingBlocked property value. Indicates whether or not to Block the user from voice recording.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetVoiceRecordingBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("voiceRecordingBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWebRtcBlockLocalhostIpAddress gets the webRtcBlockLocalhostIpAddress property value. Indicates whether or not user's localhost IP address is displayed while making phone calls using the WebRTC
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWebRtcBlockLocalhostIpAddress()(*bool) {
    val, err := m.GetBackingStore().Get("webRtcBlockLocalhostIpAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWiFiBlockAutomaticConnectHotspots gets the wiFiBlockAutomaticConnectHotspots property value. Indicating whether or not to block automatically connecting to Wi-Fi hotspots. Has no impact if Wi-Fi is blocked.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWiFiBlockAutomaticConnectHotspots()(*bool) {
    val, err := m.GetBackingStore().Get("wiFiBlockAutomaticConnectHotspots")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWiFiBlocked gets the wiFiBlocked property value. Indicates whether or not to Block the user from using Wi-Fi.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWiFiBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("wiFiBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWiFiBlockManualConfiguration gets the wiFiBlockManualConfiguration property value. Indicates whether or not to Block the user from using Wi-Fi manual configuration.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWiFiBlockManualConfiguration()(*bool) {
    val, err := m.GetBackingStore().Get("wiFiBlockManualConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWiFiScanInterval gets the wiFiScanInterval property value. Specify how often devices scan for Wi-Fi networks. Supported values are 1-500, where 100 = default, and 500 = low frequency. Valid values 1 to 500
// returns a *int32 when successful
func (m *Windows10GeneralConfiguration) GetWiFiScanInterval()(*int32) {
    val, err := m.GetBackingStore().Get("wiFiScanInterval")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetWindowsSpotlightBlockConsumerSpecificFeatures gets the windowsSpotlightBlockConsumerSpecificFeatures property value. Allows IT admins to block experiences that are typically for consumers only, such as Start suggestions, Membership notifications, Post-OOBE app install and redirect tiles.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWindowsSpotlightBlockConsumerSpecificFeatures()(*bool) {
    val, err := m.GetBackingStore().Get("windowsSpotlightBlockConsumerSpecificFeatures")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsSpotlightBlocked gets the windowsSpotlightBlocked property value. Allows IT admins to turn off all Windows Spotlight features
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWindowsSpotlightBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("windowsSpotlightBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsSpotlightBlockOnActionCenter gets the windowsSpotlightBlockOnActionCenter property value. Block suggestions from Microsoft that show after each OS clean install, upgrade or in an on-going basis to introduce users to what is new or changed
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWindowsSpotlightBlockOnActionCenter()(*bool) {
    val, err := m.GetBackingStore().Get("windowsSpotlightBlockOnActionCenter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsSpotlightBlockTailoredExperiences gets the windowsSpotlightBlockTailoredExperiences property value. Block personalized content in Windows spotlight based on user’s device usage.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWindowsSpotlightBlockTailoredExperiences()(*bool) {
    val, err := m.GetBackingStore().Get("windowsSpotlightBlockTailoredExperiences")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsSpotlightBlockThirdPartyNotifications gets the windowsSpotlightBlockThirdPartyNotifications property value. Block third party content delivered via Windows Spotlight
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWindowsSpotlightBlockThirdPartyNotifications()(*bool) {
    val, err := m.GetBackingStore().Get("windowsSpotlightBlockThirdPartyNotifications")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsSpotlightBlockWelcomeExperience gets the windowsSpotlightBlockWelcomeExperience property value. Block Windows Spotlight Windows welcome experience
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWindowsSpotlightBlockWelcomeExperience()(*bool) {
    val, err := m.GetBackingStore().Get("windowsSpotlightBlockWelcomeExperience")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsSpotlightBlockWindowsTips gets the windowsSpotlightBlockWindowsTips property value. Allows IT admins to turn off the popup of Windows Tips.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWindowsSpotlightBlockWindowsTips()(*bool) {
    val, err := m.GetBackingStore().Get("windowsSpotlightBlockWindowsTips")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsSpotlightConfigureOnLockScreen gets the windowsSpotlightConfigureOnLockScreen property value. Allows IT admind to set a predefined default search engine for MDM-Controlled devices
// returns a *WindowsSpotlightEnablementSettings when successful
func (m *Windows10GeneralConfiguration) GetWindowsSpotlightConfigureOnLockScreen()(*WindowsSpotlightEnablementSettings) {
    val, err := m.GetBackingStore().Get("windowsSpotlightConfigureOnLockScreen")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsSpotlightEnablementSettings)
    }
    return nil
}
// GetWindowsStoreBlockAutoUpdate gets the windowsStoreBlockAutoUpdate property value. Indicates whether or not to block automatic update of apps from Windows Store.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWindowsStoreBlockAutoUpdate()(*bool) {
    val, err := m.GetBackingStore().Get("windowsStoreBlockAutoUpdate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsStoreBlocked gets the windowsStoreBlocked property value. Indicates whether or not to Block the user from using the Windows store.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWindowsStoreBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("windowsStoreBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsStoreEnablePrivateStoreOnly gets the windowsStoreEnablePrivateStoreOnly property value. Indicates whether or not to enable Private Store Only.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWindowsStoreEnablePrivateStoreOnly()(*bool) {
    val, err := m.GetBackingStore().Get("windowsStoreEnablePrivateStoreOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWirelessDisplayBlockProjectionToThisDevice gets the wirelessDisplayBlockProjectionToThisDevice property value. Indicates whether or not to allow other devices from discovering this PC for projection.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWirelessDisplayBlockProjectionToThisDevice()(*bool) {
    val, err := m.GetBackingStore().Get("wirelessDisplayBlockProjectionToThisDevice")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWirelessDisplayBlockUserInputFromReceiver gets the wirelessDisplayBlockUserInputFromReceiver property value. Indicates whether or not to allow user input from wireless display receiver.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWirelessDisplayBlockUserInputFromReceiver()(*bool) {
    val, err := m.GetBackingStore().Get("wirelessDisplayBlockUserInputFromReceiver")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWirelessDisplayRequirePinForPairing gets the wirelessDisplayRequirePinForPairing property value. Indicates whether or not to require a PIN for new devices to initiate pairing.
// returns a *bool when successful
func (m *Windows10GeneralConfiguration) GetWirelessDisplayRequirePinForPairing()(*bool) {
    val, err := m.GetBackingStore().Get("wirelessDisplayRequirePinForPairing")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Windows10GeneralConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("accountsBlockAddingNonMicrosoftAccountEmail", m.GetAccountsBlockAddingNonMicrosoftAccountEmail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("antiTheftModeBlocked", m.GetAntiTheftModeBlocked())
        if err != nil {
            return err
        }
    }
    if m.GetAppsAllowTrustedAppsSideloading() != nil {
        cast := (*m.GetAppsAllowTrustedAppsSideloading()).String()
        err = writer.WriteStringValue("appsAllowTrustedAppsSideloading", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appsBlockWindowsStoreOriginatedApps", m.GetAppsBlockWindowsStoreOriginatedApps())
        if err != nil {
            return err
        }
    }
    if m.GetBluetoothAllowedServices() != nil {
        err = writer.WriteCollectionOfStringValues("bluetoothAllowedServices", m.GetBluetoothAllowedServices())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("bluetoothBlockAdvertising", m.GetBluetoothBlockAdvertising())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("bluetoothBlockDiscoverableMode", m.GetBluetoothBlockDiscoverableMode())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("bluetoothBlocked", m.GetBluetoothBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("bluetoothBlockPrePairing", m.GetBluetoothBlockPrePairing())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cameraBlocked", m.GetCameraBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cellularBlockDataWhenRoaming", m.GetCellularBlockDataWhenRoaming())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cellularBlockVpn", m.GetCellularBlockVpn())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cellularBlockVpnWhenRoaming", m.GetCellularBlockVpnWhenRoaming())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("certificatesBlockManualRootCertificateInstallation", m.GetCertificatesBlockManualRootCertificateInstallation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("connectedDevicesServiceBlocked", m.GetConnectedDevicesServiceBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("copyPasteBlocked", m.GetCopyPasteBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cortanaBlocked", m.GetCortanaBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderBlockEndUserAccess", m.GetDefenderBlockEndUserAccess())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderCloudBlockLevel() != nil {
        cast := (*m.GetDefenderCloudBlockLevel()).String()
        err = writer.WriteStringValue("defenderCloudBlockLevel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("defenderDaysBeforeDeletingQuarantinedMalware", m.GetDefenderDaysBeforeDeletingQuarantinedMalware())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("defenderDetectedMalwareActions", m.GetDefenderDetectedMalwareActions())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderFileExtensionsToExclude() != nil {
        err = writer.WriteCollectionOfStringValues("defenderFileExtensionsToExclude", m.GetDefenderFileExtensionsToExclude())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderFilesAndFoldersToExclude() != nil {
        err = writer.WriteCollectionOfStringValues("defenderFilesAndFoldersToExclude", m.GetDefenderFilesAndFoldersToExclude())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderMonitorFileActivity() != nil {
        cast := (*m.GetDefenderMonitorFileActivity()).String()
        err = writer.WriteStringValue("defenderMonitorFileActivity", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetDefenderProcessesToExclude() != nil {
        err = writer.WriteCollectionOfStringValues("defenderProcessesToExclude", m.GetDefenderProcessesToExclude())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderPromptForSampleSubmission() != nil {
        cast := (*m.GetDefenderPromptForSampleSubmission()).String()
        err = writer.WriteStringValue("defenderPromptForSampleSubmission", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderRequireBehaviorMonitoring", m.GetDefenderRequireBehaviorMonitoring())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderRequireCloudProtection", m.GetDefenderRequireCloudProtection())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderRequireNetworkInspectionSystem", m.GetDefenderRequireNetworkInspectionSystem())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderRequireRealTimeMonitoring", m.GetDefenderRequireRealTimeMonitoring())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderScanArchiveFiles", m.GetDefenderScanArchiveFiles())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderScanDownloads", m.GetDefenderScanDownloads())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderScanIncomingMail", m.GetDefenderScanIncomingMail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderScanMappedNetworkDrivesDuringFullScan", m.GetDefenderScanMappedNetworkDrivesDuringFullScan())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("defenderScanMaxCpu", m.GetDefenderScanMaxCpu())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderScanNetworkFiles", m.GetDefenderScanNetworkFiles())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderScanRemovableDrivesDuringFullScan", m.GetDefenderScanRemovableDrivesDuringFullScan())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("defenderScanScriptsLoadedInInternetExplorer", m.GetDefenderScanScriptsLoadedInInternetExplorer())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderScanType() != nil {
        cast := (*m.GetDefenderScanType()).String()
        err = writer.WriteStringValue("defenderScanType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeOnlyValue("defenderScheduledQuickScanTime", m.GetDefenderScheduledQuickScanTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeOnlyValue("defenderScheduledScanTime", m.GetDefenderScheduledScanTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("defenderSignatureUpdateIntervalInHours", m.GetDefenderSignatureUpdateIntervalInHours())
        if err != nil {
            return err
        }
    }
    if m.GetDefenderSystemScanSchedule() != nil {
        cast := (*m.GetDefenderSystemScanSchedule()).String()
        err = writer.WriteStringValue("defenderSystemScanSchedule", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeveloperUnlockSetting() != nil {
        cast := (*m.GetDeveloperUnlockSetting()).String()
        err = writer.WriteStringValue("developerUnlockSetting", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("deviceManagementBlockFactoryResetOnMobile", m.GetDeviceManagementBlockFactoryResetOnMobile())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("deviceManagementBlockManualUnenroll", m.GetDeviceManagementBlockManualUnenroll())
        if err != nil {
            return err
        }
    }
    if m.GetDiagnosticsDataSubmissionMode() != nil {
        cast := (*m.GetDiagnosticsDataSubmissionMode()).String()
        err = writer.WriteStringValue("diagnosticsDataSubmissionMode", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeAllowStartPagesModification", m.GetEdgeAllowStartPagesModification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockAccessToAboutFlags", m.GetEdgeBlockAccessToAboutFlags())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockAddressBarDropdown", m.GetEdgeBlockAddressBarDropdown())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockAutofill", m.GetEdgeBlockAutofill())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockCompatibilityList", m.GetEdgeBlockCompatibilityList())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockDeveloperTools", m.GetEdgeBlockDeveloperTools())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlocked", m.GetEdgeBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockExtensions", m.GetEdgeBlockExtensions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockInPrivateBrowsing", m.GetEdgeBlockInPrivateBrowsing())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockJavaScript", m.GetEdgeBlockJavaScript())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockLiveTileDataCollection", m.GetEdgeBlockLiveTileDataCollection())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockPasswordManager", m.GetEdgeBlockPasswordManager())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockPopups", m.GetEdgeBlockPopups())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockSearchSuggestions", m.GetEdgeBlockSearchSuggestions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockSendingDoNotTrackHeader", m.GetEdgeBlockSendingDoNotTrackHeader())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeBlockSendingIntranetTrafficToInternetExplorer", m.GetEdgeBlockSendingIntranetTrafficToInternetExplorer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeClearBrowsingDataOnExit", m.GetEdgeClearBrowsingDataOnExit())
        if err != nil {
            return err
        }
    }
    if m.GetEdgeCookiePolicy() != nil {
        cast := (*m.GetEdgeCookiePolicy()).String()
        err = writer.WriteStringValue("edgeCookiePolicy", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeDisableFirstRunPage", m.GetEdgeDisableFirstRunPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("edgeEnterpriseModeSiteListLocation", m.GetEdgeEnterpriseModeSiteListLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("edgeFirstRunUrl", m.GetEdgeFirstRunUrl())
        if err != nil {
            return err
        }
    }
    if m.GetEdgeHomepageUrls() != nil {
        err = writer.WriteCollectionOfStringValues("edgeHomepageUrls", m.GetEdgeHomepageUrls())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeRequireSmartScreen", m.GetEdgeRequireSmartScreen())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("edgeSearchEngine", m.GetEdgeSearchEngine())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeSendIntranetTrafficToInternetExplorer", m.GetEdgeSendIntranetTrafficToInternetExplorer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("edgeSyncFavoritesWithInternetExplorer", m.GetEdgeSyncFavoritesWithInternetExplorer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("enterpriseCloudPrintDiscoveryEndPoint", m.GetEnterpriseCloudPrintDiscoveryEndPoint())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("enterpriseCloudPrintDiscoveryMaxLimit", m.GetEnterpriseCloudPrintDiscoveryMaxLimit())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("enterpriseCloudPrintMopriaDiscoveryResourceIdentifier", m.GetEnterpriseCloudPrintMopriaDiscoveryResourceIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("enterpriseCloudPrintOAuthAuthority", m.GetEnterpriseCloudPrintOAuthAuthority())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("enterpriseCloudPrintOAuthClientIdentifier", m.GetEnterpriseCloudPrintOAuthClientIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("enterpriseCloudPrintResourceIdentifier", m.GetEnterpriseCloudPrintResourceIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("experienceBlockDeviceDiscovery", m.GetExperienceBlockDeviceDiscovery())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("experienceBlockErrorDialogWhenNoSIM", m.GetExperienceBlockErrorDialogWhenNoSIM())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("experienceBlockTaskSwitcher", m.GetExperienceBlockTaskSwitcher())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("gameDvrBlocked", m.GetGameDvrBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("internetSharingBlocked", m.GetInternetSharingBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("locationServicesBlocked", m.GetLocationServicesBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("lockScreenAllowTimeoutConfiguration", m.GetLockScreenAllowTimeoutConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("lockScreenBlockActionCenterNotifications", m.GetLockScreenBlockActionCenterNotifications())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("lockScreenBlockCortana", m.GetLockScreenBlockCortana())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("lockScreenBlockToastNotifications", m.GetLockScreenBlockToastNotifications())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("lockScreenTimeoutInSeconds", m.GetLockScreenTimeoutInSeconds())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("logonBlockFastUserSwitching", m.GetLogonBlockFastUserSwitching())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("microsoftAccountBlocked", m.GetMicrosoftAccountBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("microsoftAccountBlockSettingsSync", m.GetMicrosoftAccountBlockSettingsSync())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("networkProxyApplySettingsDeviceWide", m.GetNetworkProxyApplySettingsDeviceWide())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("networkProxyAutomaticConfigurationUrl", m.GetNetworkProxyAutomaticConfigurationUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("networkProxyDisableAutoDetect", m.GetNetworkProxyDisableAutoDetect())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("networkProxyServer", m.GetNetworkProxyServer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("nfcBlocked", m.GetNfcBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("oneDriveDisableFileSync", m.GetOneDriveDisableFileSync())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordBlockSimple", m.GetPasswordBlockSimple())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordExpirationDays", m.GetPasswordExpirationDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordMinimumCharacterSetCount", m.GetPasswordMinimumCharacterSetCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordMinimumLength", m.GetPasswordMinimumLength())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordMinutesOfInactivityBeforeScreenTimeout", m.GetPasswordMinutesOfInactivityBeforeScreenTimeout())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordPreviousPasswordBlockCount", m.GetPasswordPreviousPasswordBlockCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordRequired", m.GetPasswordRequired())
        if err != nil {
            return err
        }
    }
    if m.GetPasswordRequiredType() != nil {
        cast := (*m.GetPasswordRequiredType()).String()
        err = writer.WriteStringValue("passwordRequiredType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordRequireWhenResumeFromIdleState", m.GetPasswordRequireWhenResumeFromIdleState())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passwordSignInFailureCountBeforeFactoryReset", m.GetPasswordSignInFailureCountBeforeFactoryReset())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("personalizationDesktopImageUrl", m.GetPersonalizationDesktopImageUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("personalizationLockScreenImageUrl", m.GetPersonalizationLockScreenImageUrl())
        if err != nil {
            return err
        }
    }
    if m.GetPrivacyAdvertisingId() != nil {
        cast := (*m.GetPrivacyAdvertisingId()).String()
        err = writer.WriteStringValue("privacyAdvertisingId", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("privacyAutoAcceptPairingAndConsentPrompts", m.GetPrivacyAutoAcceptPairingAndConsentPrompts())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("privacyBlockInputPersonalization", m.GetPrivacyBlockInputPersonalization())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("resetProtectionModeBlocked", m.GetResetProtectionModeBlocked())
        if err != nil {
            return err
        }
    }
    if m.GetSafeSearchFilter() != nil {
        cast := (*m.GetSafeSearchFilter()).String()
        err = writer.WriteStringValue("safeSearchFilter", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("screenCaptureBlocked", m.GetScreenCaptureBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("searchBlockDiacritics", m.GetSearchBlockDiacritics())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("searchDisableAutoLanguageDetection", m.GetSearchDisableAutoLanguageDetection())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("searchDisableIndexerBackoff", m.GetSearchDisableIndexerBackoff())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("searchDisableIndexingEncryptedItems", m.GetSearchDisableIndexingEncryptedItems())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("searchDisableIndexingRemovableDrive", m.GetSearchDisableIndexingRemovableDrive())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("searchEnableAutomaticIndexSizeManangement", m.GetSearchEnableAutomaticIndexSizeManangement())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("searchEnableRemoteQueries", m.GetSearchEnableRemoteQueries())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockAccountsPage", m.GetSettingsBlockAccountsPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockAddProvisioningPackage", m.GetSettingsBlockAddProvisioningPackage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockAppsPage", m.GetSettingsBlockAppsPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockChangeLanguage", m.GetSettingsBlockChangeLanguage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockChangePowerSleep", m.GetSettingsBlockChangePowerSleep())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockChangeRegion", m.GetSettingsBlockChangeRegion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockChangeSystemTime", m.GetSettingsBlockChangeSystemTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockDevicesPage", m.GetSettingsBlockDevicesPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockEaseOfAccessPage", m.GetSettingsBlockEaseOfAccessPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockEditDeviceName", m.GetSettingsBlockEditDeviceName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockGamingPage", m.GetSettingsBlockGamingPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockNetworkInternetPage", m.GetSettingsBlockNetworkInternetPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockPersonalizationPage", m.GetSettingsBlockPersonalizationPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockPrivacyPage", m.GetSettingsBlockPrivacyPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockRemoveProvisioningPackage", m.GetSettingsBlockRemoveProvisioningPackage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockSettingsApp", m.GetSettingsBlockSettingsApp())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockSystemPage", m.GetSettingsBlockSystemPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockTimeLanguagePage", m.GetSettingsBlockTimeLanguagePage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockUpdateSecurityPage", m.GetSettingsBlockUpdateSecurityPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("sharedUserAppDataAllowed", m.GetSharedUserAppDataAllowed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("smartScreenBlockPromptOverride", m.GetSmartScreenBlockPromptOverride())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("smartScreenBlockPromptOverrideForFiles", m.GetSmartScreenBlockPromptOverrideForFiles())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("smartScreenEnableAppInstallControl", m.GetSmartScreenEnableAppInstallControl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startBlockUnpinningAppsFromTaskbar", m.GetStartBlockUnpinningAppsFromTaskbar())
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuAppListVisibility() != nil {
        cast := (*m.GetStartMenuAppListVisibility()).String()
        err = writer.WriteStringValue("startMenuAppListVisibility", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideChangeAccountSettings", m.GetStartMenuHideChangeAccountSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideFrequentlyUsedApps", m.GetStartMenuHideFrequentlyUsedApps())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideHibernate", m.GetStartMenuHideHibernate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideLock", m.GetStartMenuHideLock())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHidePowerButton", m.GetStartMenuHidePowerButton())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideRecentJumpLists", m.GetStartMenuHideRecentJumpLists())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideRecentlyAddedApps", m.GetStartMenuHideRecentlyAddedApps())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideRestartOptions", m.GetStartMenuHideRestartOptions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideShutDown", m.GetStartMenuHideShutDown())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideSignOut", m.GetStartMenuHideSignOut())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideSleep", m.GetStartMenuHideSleep())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideSwitchAccount", m.GetStartMenuHideSwitchAccount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("startMenuHideUserTile", m.GetStartMenuHideUserTile())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("startMenuLayoutEdgeAssetsXml", m.GetStartMenuLayoutEdgeAssetsXml())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("startMenuLayoutXml", m.GetStartMenuLayoutXml())
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuMode() != nil {
        cast := (*m.GetStartMenuMode()).String()
        err = writer.WriteStringValue("startMenuMode", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuPinnedFolderDocuments() != nil {
        cast := (*m.GetStartMenuPinnedFolderDocuments()).String()
        err = writer.WriteStringValue("startMenuPinnedFolderDocuments", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuPinnedFolderDownloads() != nil {
        cast := (*m.GetStartMenuPinnedFolderDownloads()).String()
        err = writer.WriteStringValue("startMenuPinnedFolderDownloads", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuPinnedFolderFileExplorer() != nil {
        cast := (*m.GetStartMenuPinnedFolderFileExplorer()).String()
        err = writer.WriteStringValue("startMenuPinnedFolderFileExplorer", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuPinnedFolderHomeGroup() != nil {
        cast := (*m.GetStartMenuPinnedFolderHomeGroup()).String()
        err = writer.WriteStringValue("startMenuPinnedFolderHomeGroup", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuPinnedFolderMusic() != nil {
        cast := (*m.GetStartMenuPinnedFolderMusic()).String()
        err = writer.WriteStringValue("startMenuPinnedFolderMusic", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuPinnedFolderNetwork() != nil {
        cast := (*m.GetStartMenuPinnedFolderNetwork()).String()
        err = writer.WriteStringValue("startMenuPinnedFolderNetwork", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuPinnedFolderPersonalFolder() != nil {
        cast := (*m.GetStartMenuPinnedFolderPersonalFolder()).String()
        err = writer.WriteStringValue("startMenuPinnedFolderPersonalFolder", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuPinnedFolderPictures() != nil {
        cast := (*m.GetStartMenuPinnedFolderPictures()).String()
        err = writer.WriteStringValue("startMenuPinnedFolderPictures", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuPinnedFolderSettings() != nil {
        cast := (*m.GetStartMenuPinnedFolderSettings()).String()
        err = writer.WriteStringValue("startMenuPinnedFolderSettings", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetStartMenuPinnedFolderVideos() != nil {
        cast := (*m.GetStartMenuPinnedFolderVideos()).String()
        err = writer.WriteStringValue("startMenuPinnedFolderVideos", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("storageBlockRemovableStorage", m.GetStorageBlockRemovableStorage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("storageRequireMobileDeviceEncryption", m.GetStorageRequireMobileDeviceEncryption())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("storageRestrictAppDataToSystemVolume", m.GetStorageRestrictAppDataToSystemVolume())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("storageRestrictAppInstallToSystemVolume", m.GetStorageRestrictAppInstallToSystemVolume())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("tenantLockdownRequireNetworkDuringOutOfBoxExperience", m.GetTenantLockdownRequireNetworkDuringOutOfBoxExperience())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("usbBlocked", m.GetUsbBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("voiceRecordingBlocked", m.GetVoiceRecordingBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("webRtcBlockLocalhostIpAddress", m.GetWebRtcBlockLocalhostIpAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wiFiBlockAutomaticConnectHotspots", m.GetWiFiBlockAutomaticConnectHotspots())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wiFiBlocked", m.GetWiFiBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wiFiBlockManualConfiguration", m.GetWiFiBlockManualConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("wiFiScanInterval", m.GetWiFiScanInterval())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsSpotlightBlockConsumerSpecificFeatures", m.GetWindowsSpotlightBlockConsumerSpecificFeatures())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsSpotlightBlocked", m.GetWindowsSpotlightBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsSpotlightBlockOnActionCenter", m.GetWindowsSpotlightBlockOnActionCenter())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsSpotlightBlockTailoredExperiences", m.GetWindowsSpotlightBlockTailoredExperiences())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsSpotlightBlockThirdPartyNotifications", m.GetWindowsSpotlightBlockThirdPartyNotifications())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsSpotlightBlockWelcomeExperience", m.GetWindowsSpotlightBlockWelcomeExperience())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsSpotlightBlockWindowsTips", m.GetWindowsSpotlightBlockWindowsTips())
        if err != nil {
            return err
        }
    }
    if m.GetWindowsSpotlightConfigureOnLockScreen() != nil {
        cast := (*m.GetWindowsSpotlightConfigureOnLockScreen()).String()
        err = writer.WriteStringValue("windowsSpotlightConfigureOnLockScreen", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsStoreBlockAutoUpdate", m.GetWindowsStoreBlockAutoUpdate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsStoreBlocked", m.GetWindowsStoreBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsStoreEnablePrivateStoreOnly", m.GetWindowsStoreEnablePrivateStoreOnly())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wirelessDisplayBlockProjectionToThisDevice", m.GetWirelessDisplayBlockProjectionToThisDevice())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wirelessDisplayBlockUserInputFromReceiver", m.GetWirelessDisplayBlockUserInputFromReceiver())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wirelessDisplayRequirePinForPairing", m.GetWirelessDisplayRequirePinForPairing())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccountsBlockAddingNonMicrosoftAccountEmail sets the accountsBlockAddingNonMicrosoftAccountEmail property value. Indicates whether or not to Block the user from adding email accounts to the device that are not associated with a Microsoft account.
func (m *Windows10GeneralConfiguration) SetAccountsBlockAddingNonMicrosoftAccountEmail(value *bool)() {
    err := m.GetBackingStore().Set("accountsBlockAddingNonMicrosoftAccountEmail", value)
    if err != nil {
        panic(err)
    }
}
// SetAntiTheftModeBlocked sets the antiTheftModeBlocked property value. Indicates whether or not to block the user from selecting an AntiTheft mode preference (Windows 10 Mobile only).
func (m *Windows10GeneralConfiguration) SetAntiTheftModeBlocked(value *bool)() {
    err := m.GetBackingStore().Set("antiTheftModeBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsAllowTrustedAppsSideloading sets the appsAllowTrustedAppsSideloading property value. State Management Setting.
func (m *Windows10GeneralConfiguration) SetAppsAllowTrustedAppsSideloading(value *StateManagementSetting)() {
    err := m.GetBackingStore().Set("appsAllowTrustedAppsSideloading", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsBlockWindowsStoreOriginatedApps sets the appsBlockWindowsStoreOriginatedApps property value. Indicates whether or not to disable the launch of all apps from Windows Store that came pre-installed or were downloaded.
func (m *Windows10GeneralConfiguration) SetAppsBlockWindowsStoreOriginatedApps(value *bool)() {
    err := m.GetBackingStore().Set("appsBlockWindowsStoreOriginatedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetBluetoothAllowedServices sets the bluetoothAllowedServices property value. Specify a list of allowed Bluetooth services and profiles in hex formatted strings.
func (m *Windows10GeneralConfiguration) SetBluetoothAllowedServices(value []string)() {
    err := m.GetBackingStore().Set("bluetoothAllowedServices", value)
    if err != nil {
        panic(err)
    }
}
// SetBluetoothBlockAdvertising sets the bluetoothBlockAdvertising property value. Whether or not to Block the user from using bluetooth advertising.
func (m *Windows10GeneralConfiguration) SetBluetoothBlockAdvertising(value *bool)() {
    err := m.GetBackingStore().Set("bluetoothBlockAdvertising", value)
    if err != nil {
        panic(err)
    }
}
// SetBluetoothBlockDiscoverableMode sets the bluetoothBlockDiscoverableMode property value. Whether or not to Block the user from using bluetooth discoverable mode.
func (m *Windows10GeneralConfiguration) SetBluetoothBlockDiscoverableMode(value *bool)() {
    err := m.GetBackingStore().Set("bluetoothBlockDiscoverableMode", value)
    if err != nil {
        panic(err)
    }
}
// SetBluetoothBlocked sets the bluetoothBlocked property value. Whether or not to Block the user from using bluetooth.
func (m *Windows10GeneralConfiguration) SetBluetoothBlocked(value *bool)() {
    err := m.GetBackingStore().Set("bluetoothBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetBluetoothBlockPrePairing sets the bluetoothBlockPrePairing property value. Whether or not to block specific bundled Bluetooth peripherals to automatically pair with the host device.
func (m *Windows10GeneralConfiguration) SetBluetoothBlockPrePairing(value *bool)() {
    err := m.GetBackingStore().Set("bluetoothBlockPrePairing", value)
    if err != nil {
        panic(err)
    }
}
// SetCameraBlocked sets the cameraBlocked property value. Whether or not to Block the user from accessing the camera of the device.
func (m *Windows10GeneralConfiguration) SetCameraBlocked(value *bool)() {
    err := m.GetBackingStore().Set("cameraBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockDataWhenRoaming sets the cellularBlockDataWhenRoaming property value. Whether or not to Block the user from using data over cellular while roaming.
func (m *Windows10GeneralConfiguration) SetCellularBlockDataWhenRoaming(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockDataWhenRoaming", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockVpn sets the cellularBlockVpn property value. Whether or not to Block the user from using VPN over cellular.
func (m *Windows10GeneralConfiguration) SetCellularBlockVpn(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockVpn", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockVpnWhenRoaming sets the cellularBlockVpnWhenRoaming property value. Whether or not to Block the user from using VPN when roaming over cellular.
func (m *Windows10GeneralConfiguration) SetCellularBlockVpnWhenRoaming(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockVpnWhenRoaming", value)
    if err != nil {
        panic(err)
    }
}
// SetCertificatesBlockManualRootCertificateInstallation sets the certificatesBlockManualRootCertificateInstallation property value. Whether or not to Block the user from doing manual root certificate installation.
func (m *Windows10GeneralConfiguration) SetCertificatesBlockManualRootCertificateInstallation(value *bool)() {
    err := m.GetBackingStore().Set("certificatesBlockManualRootCertificateInstallation", value)
    if err != nil {
        panic(err)
    }
}
// SetConnectedDevicesServiceBlocked sets the connectedDevicesServiceBlocked property value. Whether or not to block Connected Devices Service which enables discovery and connection to other devices, remote messaging, remote app sessions and other cross-device experiences.
func (m *Windows10GeneralConfiguration) SetConnectedDevicesServiceBlocked(value *bool)() {
    err := m.GetBackingStore().Set("connectedDevicesServiceBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetCopyPasteBlocked sets the copyPasteBlocked property value. Whether or not to Block the user from using copy paste.
func (m *Windows10GeneralConfiguration) SetCopyPasteBlocked(value *bool)() {
    err := m.GetBackingStore().Set("copyPasteBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetCortanaBlocked sets the cortanaBlocked property value. Whether or not to Block the user from using Cortana.
func (m *Windows10GeneralConfiguration) SetCortanaBlocked(value *bool)() {
    err := m.GetBackingStore().Set("cortanaBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderBlockEndUserAccess sets the defenderBlockEndUserAccess property value. Whether or not to block end user access to Defender.
func (m *Windows10GeneralConfiguration) SetDefenderBlockEndUserAccess(value *bool)() {
    err := m.GetBackingStore().Set("defenderBlockEndUserAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderCloudBlockLevel sets the defenderCloudBlockLevel property value. Possible values of Cloud Block Level
func (m *Windows10GeneralConfiguration) SetDefenderCloudBlockLevel(value *DefenderCloudBlockLevelType)() {
    err := m.GetBackingStore().Set("defenderCloudBlockLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderDaysBeforeDeletingQuarantinedMalware sets the defenderDaysBeforeDeletingQuarantinedMalware property value. Number of days before deleting quarantined malware. Valid values 0 to 90
func (m *Windows10GeneralConfiguration) SetDefenderDaysBeforeDeletingQuarantinedMalware(value *int32)() {
    err := m.GetBackingStore().Set("defenderDaysBeforeDeletingQuarantinedMalware", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderDetectedMalwareActions sets the defenderDetectedMalwareActions property value. Gets or sets Defender’s actions to take on detected Malware per threat level.
func (m *Windows10GeneralConfiguration) SetDefenderDetectedMalwareActions(value DefenderDetectedMalwareActionsable)() {
    err := m.GetBackingStore().Set("defenderDetectedMalwareActions", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderFileExtensionsToExclude sets the defenderFileExtensionsToExclude property value. File extensions to exclude from scans and real time protection.
func (m *Windows10GeneralConfiguration) SetDefenderFileExtensionsToExclude(value []string)() {
    err := m.GetBackingStore().Set("defenderFileExtensionsToExclude", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderFilesAndFoldersToExclude sets the defenderFilesAndFoldersToExclude property value. Files and folder to exclude from scans and real time protection.
func (m *Windows10GeneralConfiguration) SetDefenderFilesAndFoldersToExclude(value []string)() {
    err := m.GetBackingStore().Set("defenderFilesAndFoldersToExclude", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderMonitorFileActivity sets the defenderMonitorFileActivity property value. Possible values for monitoring file activity.
func (m *Windows10GeneralConfiguration) SetDefenderMonitorFileActivity(value *DefenderMonitorFileActivity)() {
    err := m.GetBackingStore().Set("defenderMonitorFileActivity", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderProcessesToExclude sets the defenderProcessesToExclude property value. Processes to exclude from scans and real time protection.
func (m *Windows10GeneralConfiguration) SetDefenderProcessesToExclude(value []string)() {
    err := m.GetBackingStore().Set("defenderProcessesToExclude", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderPromptForSampleSubmission sets the defenderPromptForSampleSubmission property value. Possible values for prompting user for samples submission.
func (m *Windows10GeneralConfiguration) SetDefenderPromptForSampleSubmission(value *DefenderPromptForSampleSubmission)() {
    err := m.GetBackingStore().Set("defenderPromptForSampleSubmission", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderRequireBehaviorMonitoring sets the defenderRequireBehaviorMonitoring property value. Indicates whether or not to require behavior monitoring.
func (m *Windows10GeneralConfiguration) SetDefenderRequireBehaviorMonitoring(value *bool)() {
    err := m.GetBackingStore().Set("defenderRequireBehaviorMonitoring", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderRequireCloudProtection sets the defenderRequireCloudProtection property value. Indicates whether or not to require cloud protection.
func (m *Windows10GeneralConfiguration) SetDefenderRequireCloudProtection(value *bool)() {
    err := m.GetBackingStore().Set("defenderRequireCloudProtection", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderRequireNetworkInspectionSystem sets the defenderRequireNetworkInspectionSystem property value. Indicates whether or not to require network inspection system.
func (m *Windows10GeneralConfiguration) SetDefenderRequireNetworkInspectionSystem(value *bool)() {
    err := m.GetBackingStore().Set("defenderRequireNetworkInspectionSystem", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderRequireRealTimeMonitoring sets the defenderRequireRealTimeMonitoring property value. Indicates whether or not to require real time monitoring.
func (m *Windows10GeneralConfiguration) SetDefenderRequireRealTimeMonitoring(value *bool)() {
    err := m.GetBackingStore().Set("defenderRequireRealTimeMonitoring", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScanArchiveFiles sets the defenderScanArchiveFiles property value. Indicates whether or not to scan archive files.
func (m *Windows10GeneralConfiguration) SetDefenderScanArchiveFiles(value *bool)() {
    err := m.GetBackingStore().Set("defenderScanArchiveFiles", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScanDownloads sets the defenderScanDownloads property value. Indicates whether or not to scan downloads.
func (m *Windows10GeneralConfiguration) SetDefenderScanDownloads(value *bool)() {
    err := m.GetBackingStore().Set("defenderScanDownloads", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScanIncomingMail sets the defenderScanIncomingMail property value. Indicates whether or not to scan incoming mail messages.
func (m *Windows10GeneralConfiguration) SetDefenderScanIncomingMail(value *bool)() {
    err := m.GetBackingStore().Set("defenderScanIncomingMail", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScanMappedNetworkDrivesDuringFullScan sets the defenderScanMappedNetworkDrivesDuringFullScan property value. Indicates whether or not to scan mapped network drives during full scan.
func (m *Windows10GeneralConfiguration) SetDefenderScanMappedNetworkDrivesDuringFullScan(value *bool)() {
    err := m.GetBackingStore().Set("defenderScanMappedNetworkDrivesDuringFullScan", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScanMaxCpu sets the defenderScanMaxCpu property value. Max CPU usage percentage during scan. Valid values 0 to 100
func (m *Windows10GeneralConfiguration) SetDefenderScanMaxCpu(value *int32)() {
    err := m.GetBackingStore().Set("defenderScanMaxCpu", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScanNetworkFiles sets the defenderScanNetworkFiles property value. Indicates whether or not to scan files opened from a network folder.
func (m *Windows10GeneralConfiguration) SetDefenderScanNetworkFiles(value *bool)() {
    err := m.GetBackingStore().Set("defenderScanNetworkFiles", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScanRemovableDrivesDuringFullScan sets the defenderScanRemovableDrivesDuringFullScan property value. Indicates whether or not to scan removable drives during full scan.
func (m *Windows10GeneralConfiguration) SetDefenderScanRemovableDrivesDuringFullScan(value *bool)() {
    err := m.GetBackingStore().Set("defenderScanRemovableDrivesDuringFullScan", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScanScriptsLoadedInInternetExplorer sets the defenderScanScriptsLoadedInInternetExplorer property value. Indicates whether or not to scan scripts loaded in Internet Explorer browser.
func (m *Windows10GeneralConfiguration) SetDefenderScanScriptsLoadedInInternetExplorer(value *bool)() {
    err := m.GetBackingStore().Set("defenderScanScriptsLoadedInInternetExplorer", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScanType sets the defenderScanType property value. Possible values for system scan type.
func (m *Windows10GeneralConfiguration) SetDefenderScanType(value *DefenderScanType)() {
    err := m.GetBackingStore().Set("defenderScanType", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScheduledQuickScanTime sets the defenderScheduledQuickScanTime property value. The time to perform a daily quick scan.
func (m *Windows10GeneralConfiguration) SetDefenderScheduledQuickScanTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)() {
    err := m.GetBackingStore().Set("defenderScheduledQuickScanTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderScheduledScanTime sets the defenderScheduledScanTime property value. The defender time for the system scan.
func (m *Windows10GeneralConfiguration) SetDefenderScheduledScanTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)() {
    err := m.GetBackingStore().Set("defenderScheduledScanTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderSignatureUpdateIntervalInHours sets the defenderSignatureUpdateIntervalInHours property value. The signature update interval in hours. Specify 0 not to check. Valid values 0 to 24
func (m *Windows10GeneralConfiguration) SetDefenderSignatureUpdateIntervalInHours(value *int32)() {
    err := m.GetBackingStore().Set("defenderSignatureUpdateIntervalInHours", value)
    if err != nil {
        panic(err)
    }
}
// SetDefenderSystemScanSchedule sets the defenderSystemScanSchedule property value. Possible values for a weekly schedule.
func (m *Windows10GeneralConfiguration) SetDefenderSystemScanSchedule(value *WeeklySchedule)() {
    err := m.GetBackingStore().Set("defenderSystemScanSchedule", value)
    if err != nil {
        panic(err)
    }
}
// SetDeveloperUnlockSetting sets the developerUnlockSetting property value. State Management Setting.
func (m *Windows10GeneralConfiguration) SetDeveloperUnlockSetting(value *StateManagementSetting)() {
    err := m.GetBackingStore().Set("developerUnlockSetting", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceManagementBlockFactoryResetOnMobile sets the deviceManagementBlockFactoryResetOnMobile property value. Indicates whether or not to Block the user from resetting their phone.
func (m *Windows10GeneralConfiguration) SetDeviceManagementBlockFactoryResetOnMobile(value *bool)() {
    err := m.GetBackingStore().Set("deviceManagementBlockFactoryResetOnMobile", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceManagementBlockManualUnenroll sets the deviceManagementBlockManualUnenroll property value. Indicates whether or not to Block the user from doing manual un-enrollment from device management.
func (m *Windows10GeneralConfiguration) SetDeviceManagementBlockManualUnenroll(value *bool)() {
    err := m.GetBackingStore().Set("deviceManagementBlockManualUnenroll", value)
    if err != nil {
        panic(err)
    }
}
// SetDiagnosticsDataSubmissionMode sets the diagnosticsDataSubmissionMode property value. Allow the device to send diagnostic and usage telemetry data, such as Watson.
func (m *Windows10GeneralConfiguration) SetDiagnosticsDataSubmissionMode(value *DiagnosticDataSubmissionMode)() {
    err := m.GetBackingStore().Set("diagnosticsDataSubmissionMode", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeAllowStartPagesModification sets the edgeAllowStartPagesModification property value. Allow users to change Start pages on Edge. Use the EdgeHomepageUrls to specify the Start pages that the user would see by default when they open Edge.
func (m *Windows10GeneralConfiguration) SetEdgeAllowStartPagesModification(value *bool)() {
    err := m.GetBackingStore().Set("edgeAllowStartPagesModification", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockAccessToAboutFlags sets the edgeBlockAccessToAboutFlags property value. Indicates whether or not to prevent access to about flags on Edge browser.
func (m *Windows10GeneralConfiguration) SetEdgeBlockAccessToAboutFlags(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockAccessToAboutFlags", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockAddressBarDropdown sets the edgeBlockAddressBarDropdown property value. Block the address bar dropdown functionality in Microsoft Edge. Disable this settings to minimize network connections from Microsoft Edge to Microsoft services.
func (m *Windows10GeneralConfiguration) SetEdgeBlockAddressBarDropdown(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockAddressBarDropdown", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockAutofill sets the edgeBlockAutofill property value. Indicates whether or not to block auto fill.
func (m *Windows10GeneralConfiguration) SetEdgeBlockAutofill(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockAutofill", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockCompatibilityList sets the edgeBlockCompatibilityList property value. Block Microsoft compatibility list in Microsoft Edge. This list from Microsoft helps Edge properly display sites with known compatibility issues.
func (m *Windows10GeneralConfiguration) SetEdgeBlockCompatibilityList(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockCompatibilityList", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockDeveloperTools sets the edgeBlockDeveloperTools property value. Indicates whether or not to block developer tools in the Edge browser.
func (m *Windows10GeneralConfiguration) SetEdgeBlockDeveloperTools(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockDeveloperTools", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlocked sets the edgeBlocked property value. Indicates whether or not to Block the user from using the Edge browser.
func (m *Windows10GeneralConfiguration) SetEdgeBlocked(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockExtensions sets the edgeBlockExtensions property value. Indicates whether or not to block extensions in the Edge browser.
func (m *Windows10GeneralConfiguration) SetEdgeBlockExtensions(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockExtensions", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockInPrivateBrowsing sets the edgeBlockInPrivateBrowsing property value. Indicates whether or not to block InPrivate browsing on corporate networks, in the Edge browser.
func (m *Windows10GeneralConfiguration) SetEdgeBlockInPrivateBrowsing(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockInPrivateBrowsing", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockJavaScript sets the edgeBlockJavaScript property value. Indicates whether or not to Block the user from using JavaScript.
func (m *Windows10GeneralConfiguration) SetEdgeBlockJavaScript(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockJavaScript", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockLiveTileDataCollection sets the edgeBlockLiveTileDataCollection property value. Block the collection of information by Microsoft for live tile creation when users pin a site to Start from Microsoft Edge.
func (m *Windows10GeneralConfiguration) SetEdgeBlockLiveTileDataCollection(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockLiveTileDataCollection", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockPasswordManager sets the edgeBlockPasswordManager property value. Indicates whether or not to Block password manager.
func (m *Windows10GeneralConfiguration) SetEdgeBlockPasswordManager(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockPasswordManager", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockPopups sets the edgeBlockPopups property value. Indicates whether or not to block popups.
func (m *Windows10GeneralConfiguration) SetEdgeBlockPopups(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockPopups", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockSearchSuggestions sets the edgeBlockSearchSuggestions property value. Indicates whether or not to block the user from using the search suggestions in the address bar.
func (m *Windows10GeneralConfiguration) SetEdgeBlockSearchSuggestions(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockSearchSuggestions", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockSendingDoNotTrackHeader sets the edgeBlockSendingDoNotTrackHeader property value. Indicates whether or not to Block the user from sending the do not track header.
func (m *Windows10GeneralConfiguration) SetEdgeBlockSendingDoNotTrackHeader(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockSendingDoNotTrackHeader", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeBlockSendingIntranetTrafficToInternetExplorer sets the edgeBlockSendingIntranetTrafficToInternetExplorer property value. Indicates whether or not to switch the intranet traffic from Edge to Internet Explorer. Note: the name of this property is misleading; the property is obsolete, use EdgeSendIntranetTrafficToInternetExplorer instead.
func (m *Windows10GeneralConfiguration) SetEdgeBlockSendingIntranetTrafficToInternetExplorer(value *bool)() {
    err := m.GetBackingStore().Set("edgeBlockSendingIntranetTrafficToInternetExplorer", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeClearBrowsingDataOnExit sets the edgeClearBrowsingDataOnExit property value. Clear browsing data on exiting Microsoft Edge.
func (m *Windows10GeneralConfiguration) SetEdgeClearBrowsingDataOnExit(value *bool)() {
    err := m.GetBackingStore().Set("edgeClearBrowsingDataOnExit", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeCookiePolicy sets the edgeCookiePolicy property value. Possible values to specify which cookies are allowed in Microsoft Edge.
func (m *Windows10GeneralConfiguration) SetEdgeCookiePolicy(value *EdgeCookiePolicy)() {
    err := m.GetBackingStore().Set("edgeCookiePolicy", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeDisableFirstRunPage sets the edgeDisableFirstRunPage property value. Block the Microsoft web page that opens on the first use of Microsoft Edge. This policy allows enterprises, like those enrolled in zero emissions configurations, to block this page.
func (m *Windows10GeneralConfiguration) SetEdgeDisableFirstRunPage(value *bool)() {
    err := m.GetBackingStore().Set("edgeDisableFirstRunPage", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeEnterpriseModeSiteListLocation sets the edgeEnterpriseModeSiteListLocation property value. Indicates the enterprise mode site list location. Could be a local file, local network or http location.
func (m *Windows10GeneralConfiguration) SetEdgeEnterpriseModeSiteListLocation(value *string)() {
    err := m.GetBackingStore().Set("edgeEnterpriseModeSiteListLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeFirstRunUrl sets the edgeFirstRunUrl property value. The first run URL for when Edge browser is opened for the first time.
func (m *Windows10GeneralConfiguration) SetEdgeFirstRunUrl(value *string)() {
    err := m.GetBackingStore().Set("edgeFirstRunUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeHomepageUrls sets the edgeHomepageUrls property value. The list of URLs for homepages shodwn on MDM-enrolled devices on Edge browser.
func (m *Windows10GeneralConfiguration) SetEdgeHomepageUrls(value []string)() {
    err := m.GetBackingStore().Set("edgeHomepageUrls", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeRequireSmartScreen sets the edgeRequireSmartScreen property value. Indicates whether or not to Require the user to use the smart screen filter.
func (m *Windows10GeneralConfiguration) SetEdgeRequireSmartScreen(value *bool)() {
    err := m.GetBackingStore().Set("edgeRequireSmartScreen", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeSearchEngine sets the edgeSearchEngine property value. Allows IT admins to set a default search engine for MDM-Controlled devices. Users can override this and change their default search engine provided the AllowSearchEngineCustomization policy is not set.
func (m *Windows10GeneralConfiguration) SetEdgeSearchEngine(value EdgeSearchEngineBaseable)() {
    err := m.GetBackingStore().Set("edgeSearchEngine", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeSendIntranetTrafficToInternetExplorer sets the edgeSendIntranetTrafficToInternetExplorer property value. Indicates whether or not to switch the intranet traffic from Edge to Internet Explorer.
func (m *Windows10GeneralConfiguration) SetEdgeSendIntranetTrafficToInternetExplorer(value *bool)() {
    err := m.GetBackingStore().Set("edgeSendIntranetTrafficToInternetExplorer", value)
    if err != nil {
        panic(err)
    }
}
// SetEdgeSyncFavoritesWithInternetExplorer sets the edgeSyncFavoritesWithInternetExplorer property value. Enable favorites sync between Internet Explorer and Microsoft Edge. Additions, deletions, modifications and order changes to favorites are shared between browsers.
func (m *Windows10GeneralConfiguration) SetEdgeSyncFavoritesWithInternetExplorer(value *bool)() {
    err := m.GetBackingStore().Set("edgeSyncFavoritesWithInternetExplorer", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseCloudPrintDiscoveryEndPoint sets the enterpriseCloudPrintDiscoveryEndPoint property value. Endpoint for discovering cloud printers.
func (m *Windows10GeneralConfiguration) SetEnterpriseCloudPrintDiscoveryEndPoint(value *string)() {
    err := m.GetBackingStore().Set("enterpriseCloudPrintDiscoveryEndPoint", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseCloudPrintDiscoveryMaxLimit sets the enterpriseCloudPrintDiscoveryMaxLimit property value. Maximum number of printers that should be queried from a discovery endpoint. This is a mobile only setting. Valid values 1 to 65535
func (m *Windows10GeneralConfiguration) SetEnterpriseCloudPrintDiscoveryMaxLimit(value *int32)() {
    err := m.GetBackingStore().Set("enterpriseCloudPrintDiscoveryMaxLimit", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseCloudPrintMopriaDiscoveryResourceIdentifier sets the enterpriseCloudPrintMopriaDiscoveryResourceIdentifier property value. OAuth resource URI for printer discovery service as configured in Azure portal.
func (m *Windows10GeneralConfiguration) SetEnterpriseCloudPrintMopriaDiscoveryResourceIdentifier(value *string)() {
    err := m.GetBackingStore().Set("enterpriseCloudPrintMopriaDiscoveryResourceIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseCloudPrintOAuthAuthority sets the enterpriseCloudPrintOAuthAuthority property value. Authentication endpoint for acquiring OAuth tokens.
func (m *Windows10GeneralConfiguration) SetEnterpriseCloudPrintOAuthAuthority(value *string)() {
    err := m.GetBackingStore().Set("enterpriseCloudPrintOAuthAuthority", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseCloudPrintOAuthClientIdentifier sets the enterpriseCloudPrintOAuthClientIdentifier property value. GUID of a client application authorized to retrieve OAuth tokens from the OAuth Authority.
func (m *Windows10GeneralConfiguration) SetEnterpriseCloudPrintOAuthClientIdentifier(value *string)() {
    err := m.GetBackingStore().Set("enterpriseCloudPrintOAuthClientIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseCloudPrintResourceIdentifier sets the enterpriseCloudPrintResourceIdentifier property value. OAuth resource URI for print service as configured in the Azure portal.
func (m *Windows10GeneralConfiguration) SetEnterpriseCloudPrintResourceIdentifier(value *string)() {
    err := m.GetBackingStore().Set("enterpriseCloudPrintResourceIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetExperienceBlockDeviceDiscovery sets the experienceBlockDeviceDiscovery property value. Indicates whether or not to enable device discovery UX.
func (m *Windows10GeneralConfiguration) SetExperienceBlockDeviceDiscovery(value *bool)() {
    err := m.GetBackingStore().Set("experienceBlockDeviceDiscovery", value)
    if err != nil {
        panic(err)
    }
}
// SetExperienceBlockErrorDialogWhenNoSIM sets the experienceBlockErrorDialogWhenNoSIM property value. Indicates whether or not to allow the error dialog from displaying if no SIM card is detected.
func (m *Windows10GeneralConfiguration) SetExperienceBlockErrorDialogWhenNoSIM(value *bool)() {
    err := m.GetBackingStore().Set("experienceBlockErrorDialogWhenNoSIM", value)
    if err != nil {
        panic(err)
    }
}
// SetExperienceBlockTaskSwitcher sets the experienceBlockTaskSwitcher property value. Indicates whether or not to enable task switching on the device.
func (m *Windows10GeneralConfiguration) SetExperienceBlockTaskSwitcher(value *bool)() {
    err := m.GetBackingStore().Set("experienceBlockTaskSwitcher", value)
    if err != nil {
        panic(err)
    }
}
// SetGameDvrBlocked sets the gameDvrBlocked property value. Indicates whether or not to block DVR and broadcasting.
func (m *Windows10GeneralConfiguration) SetGameDvrBlocked(value *bool)() {
    err := m.GetBackingStore().Set("gameDvrBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetInternetSharingBlocked sets the internetSharingBlocked property value. Indicates whether or not to Block the user from using internet sharing.
func (m *Windows10GeneralConfiguration) SetInternetSharingBlocked(value *bool)() {
    err := m.GetBackingStore().Set("internetSharingBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetLocationServicesBlocked sets the locationServicesBlocked property value. Indicates whether or not to Block the user from location services.
func (m *Windows10GeneralConfiguration) SetLocationServicesBlocked(value *bool)() {
    err := m.GetBackingStore().Set("locationServicesBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetLockScreenAllowTimeoutConfiguration sets the lockScreenAllowTimeoutConfiguration property value. Specify whether to show a user-configurable setting to control the screen timeout while on the lock screen of Windows 10 Mobile devices. If this policy is set to Allow, the value set by lockScreenTimeoutInSeconds is ignored.
func (m *Windows10GeneralConfiguration) SetLockScreenAllowTimeoutConfiguration(value *bool)() {
    err := m.GetBackingStore().Set("lockScreenAllowTimeoutConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetLockScreenBlockActionCenterNotifications sets the lockScreenBlockActionCenterNotifications property value. Indicates whether or not to block action center notifications over lock screen.
func (m *Windows10GeneralConfiguration) SetLockScreenBlockActionCenterNotifications(value *bool)() {
    err := m.GetBackingStore().Set("lockScreenBlockActionCenterNotifications", value)
    if err != nil {
        panic(err)
    }
}
// SetLockScreenBlockCortana sets the lockScreenBlockCortana property value. Indicates whether or not the user can interact with Cortana using speech while the system is locked.
func (m *Windows10GeneralConfiguration) SetLockScreenBlockCortana(value *bool)() {
    err := m.GetBackingStore().Set("lockScreenBlockCortana", value)
    if err != nil {
        panic(err)
    }
}
// SetLockScreenBlockToastNotifications sets the lockScreenBlockToastNotifications property value. Indicates whether to allow toast notifications above the device lock screen.
func (m *Windows10GeneralConfiguration) SetLockScreenBlockToastNotifications(value *bool)() {
    err := m.GetBackingStore().Set("lockScreenBlockToastNotifications", value)
    if err != nil {
        panic(err)
    }
}
// SetLockScreenTimeoutInSeconds sets the lockScreenTimeoutInSeconds property value. Set the duration (in seconds) from the screen locking to the screen turning off for Windows 10 Mobile devices. Supported values are 11-1800. Valid values 11 to 1800
func (m *Windows10GeneralConfiguration) SetLockScreenTimeoutInSeconds(value *int32)() {
    err := m.GetBackingStore().Set("lockScreenTimeoutInSeconds", value)
    if err != nil {
        panic(err)
    }
}
// SetLogonBlockFastUserSwitching sets the logonBlockFastUserSwitching property value. Disables the ability to quickly switch between users that are logged on simultaneously without logging off.
func (m *Windows10GeneralConfiguration) SetLogonBlockFastUserSwitching(value *bool)() {
    err := m.GetBackingStore().Set("logonBlockFastUserSwitching", value)
    if err != nil {
        panic(err)
    }
}
// SetMicrosoftAccountBlocked sets the microsoftAccountBlocked property value. Indicates whether or not to Block a Microsoft account.
func (m *Windows10GeneralConfiguration) SetMicrosoftAccountBlocked(value *bool)() {
    err := m.GetBackingStore().Set("microsoftAccountBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetMicrosoftAccountBlockSettingsSync sets the microsoftAccountBlockSettingsSync property value. Indicates whether or not to Block Microsoft account settings sync.
func (m *Windows10GeneralConfiguration) SetMicrosoftAccountBlockSettingsSync(value *bool)() {
    err := m.GetBackingStore().Set("microsoftAccountBlockSettingsSync", value)
    if err != nil {
        panic(err)
    }
}
// SetNetworkProxyApplySettingsDeviceWide sets the networkProxyApplySettingsDeviceWide property value. If set, proxy settings will be applied to all processes and accounts in the device. Otherwise, it will be applied to the user account that’s enrolled into MDM.
func (m *Windows10GeneralConfiguration) SetNetworkProxyApplySettingsDeviceWide(value *bool)() {
    err := m.GetBackingStore().Set("networkProxyApplySettingsDeviceWide", value)
    if err != nil {
        panic(err)
    }
}
// SetNetworkProxyAutomaticConfigurationUrl sets the networkProxyAutomaticConfigurationUrl property value. Address to the proxy auto-config (PAC) script you want to use.
func (m *Windows10GeneralConfiguration) SetNetworkProxyAutomaticConfigurationUrl(value *string)() {
    err := m.GetBackingStore().Set("networkProxyAutomaticConfigurationUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetNetworkProxyDisableAutoDetect sets the networkProxyDisableAutoDetect property value. Disable automatic detection of settings. If enabled, the system will try to find the path to a proxy auto-config (PAC) script.
func (m *Windows10GeneralConfiguration) SetNetworkProxyDisableAutoDetect(value *bool)() {
    err := m.GetBackingStore().Set("networkProxyDisableAutoDetect", value)
    if err != nil {
        panic(err)
    }
}
// SetNetworkProxyServer sets the networkProxyServer property value. Specifies manual proxy server settings.
func (m *Windows10GeneralConfiguration) SetNetworkProxyServer(value Windows10NetworkProxyServerable)() {
    err := m.GetBackingStore().Set("networkProxyServer", value)
    if err != nil {
        panic(err)
    }
}
// SetNfcBlocked sets the nfcBlocked property value. Indicates whether or not to Block the user from using near field communication.
func (m *Windows10GeneralConfiguration) SetNfcBlocked(value *bool)() {
    err := m.GetBackingStore().Set("nfcBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetOneDriveDisableFileSync sets the oneDriveDisableFileSync property value. Gets or sets a value allowing IT admins to prevent apps and features from working with files on OneDrive.
func (m *Windows10GeneralConfiguration) SetOneDriveDisableFileSync(value *bool)() {
    err := m.GetBackingStore().Set("oneDriveDisableFileSync", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordBlockSimple sets the passwordBlockSimple property value. Specify whether PINs or passwords such as '1111' or '1234' are allowed. For Windows 10 desktops, it also controls the use of picture passwords.
func (m *Windows10GeneralConfiguration) SetPasswordBlockSimple(value *bool)() {
    err := m.GetBackingStore().Set("passwordBlockSimple", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordExpirationDays sets the passwordExpirationDays property value. The password expiration in days. Valid values 0 to 730
func (m *Windows10GeneralConfiguration) SetPasswordExpirationDays(value *int32)() {
    err := m.GetBackingStore().Set("passwordExpirationDays", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinimumCharacterSetCount sets the passwordMinimumCharacterSetCount property value. The number of character sets required in the password.
func (m *Windows10GeneralConfiguration) SetPasswordMinimumCharacterSetCount(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinimumCharacterSetCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinimumLength sets the passwordMinimumLength property value. The minimum password length. Valid values 4 to 16
func (m *Windows10GeneralConfiguration) SetPasswordMinimumLength(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinimumLength", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinutesOfInactivityBeforeScreenTimeout sets the passwordMinutesOfInactivityBeforeScreenTimeout property value. The minutes of inactivity before the screen times out.
func (m *Windows10GeneralConfiguration) SetPasswordMinutesOfInactivityBeforeScreenTimeout(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinutesOfInactivityBeforeScreenTimeout", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordPreviousPasswordBlockCount sets the passwordPreviousPasswordBlockCount property value. The number of previous passwords to prevent reuse of. Valid values 0 to 50
func (m *Windows10GeneralConfiguration) SetPasswordPreviousPasswordBlockCount(value *int32)() {
    err := m.GetBackingStore().Set("passwordPreviousPasswordBlockCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordRequired sets the passwordRequired property value. Indicates whether or not to require the user to have a password.
func (m *Windows10GeneralConfiguration) SetPasswordRequired(value *bool)() {
    err := m.GetBackingStore().Set("passwordRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordRequiredType sets the passwordRequiredType property value. Possible values of required passwords.
func (m *Windows10GeneralConfiguration) SetPasswordRequiredType(value *RequiredPasswordType)() {
    err := m.GetBackingStore().Set("passwordRequiredType", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordRequireWhenResumeFromIdleState sets the passwordRequireWhenResumeFromIdleState property value. Indicates whether or not to require a password upon resuming from an idle state.
func (m *Windows10GeneralConfiguration) SetPasswordRequireWhenResumeFromIdleState(value *bool)() {
    err := m.GetBackingStore().Set("passwordRequireWhenResumeFromIdleState", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordSignInFailureCountBeforeFactoryReset sets the passwordSignInFailureCountBeforeFactoryReset property value. The number of sign in failures before factory reset. Valid values 0 to 999
func (m *Windows10GeneralConfiguration) SetPasswordSignInFailureCountBeforeFactoryReset(value *int32)() {
    err := m.GetBackingStore().Set("passwordSignInFailureCountBeforeFactoryReset", value)
    if err != nil {
        panic(err)
    }
}
// SetPersonalizationDesktopImageUrl sets the personalizationDesktopImageUrl property value. A http or https Url to a jpg, jpeg or png image that needs to be downloaded and used as the Desktop Image or a file Url to a local image on the file system that needs to used as the Desktop Image.
func (m *Windows10GeneralConfiguration) SetPersonalizationDesktopImageUrl(value *string)() {
    err := m.GetBackingStore().Set("personalizationDesktopImageUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetPersonalizationLockScreenImageUrl sets the personalizationLockScreenImageUrl property value. A http or https Url to a jpg, jpeg or png image that neeeds to be downloaded and used as the Lock Screen Image or a file Url to a local image on the file system that needs to be used as the Lock Screen Image.
func (m *Windows10GeneralConfiguration) SetPersonalizationLockScreenImageUrl(value *string)() {
    err := m.GetBackingStore().Set("personalizationLockScreenImageUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetPrivacyAdvertisingId sets the privacyAdvertisingId property value. State Management Setting.
func (m *Windows10GeneralConfiguration) SetPrivacyAdvertisingId(value *StateManagementSetting)() {
    err := m.GetBackingStore().Set("privacyAdvertisingId", value)
    if err != nil {
        panic(err)
    }
}
// SetPrivacyAutoAcceptPairingAndConsentPrompts sets the privacyAutoAcceptPairingAndConsentPrompts property value. Indicates whether or not to allow the automatic acceptance of the pairing and privacy user consent dialog when launching apps.
func (m *Windows10GeneralConfiguration) SetPrivacyAutoAcceptPairingAndConsentPrompts(value *bool)() {
    err := m.GetBackingStore().Set("privacyAutoAcceptPairingAndConsentPrompts", value)
    if err != nil {
        panic(err)
    }
}
// SetPrivacyBlockInputPersonalization sets the privacyBlockInputPersonalization property value. Indicates whether or not to block the usage of cloud based speech services for Cortana, Dictation, or Store applications.
func (m *Windows10GeneralConfiguration) SetPrivacyBlockInputPersonalization(value *bool)() {
    err := m.GetBackingStore().Set("privacyBlockInputPersonalization", value)
    if err != nil {
        panic(err)
    }
}
// SetResetProtectionModeBlocked sets the resetProtectionModeBlocked property value. Indicates whether or not to Block the user from reset protection mode.
func (m *Windows10GeneralConfiguration) SetResetProtectionModeBlocked(value *bool)() {
    err := m.GetBackingStore().Set("resetProtectionModeBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetSafeSearchFilter sets the safeSearchFilter property value. Specifies what level of safe search (filtering adult content) is required
func (m *Windows10GeneralConfiguration) SetSafeSearchFilter(value *SafeSearchFilterType)() {
    err := m.GetBackingStore().Set("safeSearchFilter", value)
    if err != nil {
        panic(err)
    }
}
// SetScreenCaptureBlocked sets the screenCaptureBlocked property value. Indicates whether or not to Block the user from taking Screenshots.
func (m *Windows10GeneralConfiguration) SetScreenCaptureBlocked(value *bool)() {
    err := m.GetBackingStore().Set("screenCaptureBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetSearchBlockDiacritics sets the searchBlockDiacritics property value. Specifies if search can use diacritics.
func (m *Windows10GeneralConfiguration) SetSearchBlockDiacritics(value *bool)() {
    err := m.GetBackingStore().Set("searchBlockDiacritics", value)
    if err != nil {
        panic(err)
    }
}
// SetSearchDisableAutoLanguageDetection sets the searchDisableAutoLanguageDetection property value. Specifies whether to use automatic language detection when indexing content and properties.
func (m *Windows10GeneralConfiguration) SetSearchDisableAutoLanguageDetection(value *bool)() {
    err := m.GetBackingStore().Set("searchDisableAutoLanguageDetection", value)
    if err != nil {
        panic(err)
    }
}
// SetSearchDisableIndexerBackoff sets the searchDisableIndexerBackoff property value. Indicates whether or not to disable the search indexer backoff feature.
func (m *Windows10GeneralConfiguration) SetSearchDisableIndexerBackoff(value *bool)() {
    err := m.GetBackingStore().Set("searchDisableIndexerBackoff", value)
    if err != nil {
        panic(err)
    }
}
// SetSearchDisableIndexingEncryptedItems sets the searchDisableIndexingEncryptedItems property value. Indicates whether or not to block indexing of WIP-protected items to prevent them from appearing in search results for Cortana or Explorer.
func (m *Windows10GeneralConfiguration) SetSearchDisableIndexingEncryptedItems(value *bool)() {
    err := m.GetBackingStore().Set("searchDisableIndexingEncryptedItems", value)
    if err != nil {
        panic(err)
    }
}
// SetSearchDisableIndexingRemovableDrive sets the searchDisableIndexingRemovableDrive property value. Indicates whether or not to allow users to add locations on removable drives to libraries and to be indexed.
func (m *Windows10GeneralConfiguration) SetSearchDisableIndexingRemovableDrive(value *bool)() {
    err := m.GetBackingStore().Set("searchDisableIndexingRemovableDrive", value)
    if err != nil {
        panic(err)
    }
}
// SetSearchEnableAutomaticIndexSizeManangement sets the searchEnableAutomaticIndexSizeManangement property value. Specifies minimum amount of hard drive space on the same drive as the index location before indexing stops.
func (m *Windows10GeneralConfiguration) SetSearchEnableAutomaticIndexSizeManangement(value *bool)() {
    err := m.GetBackingStore().Set("searchEnableAutomaticIndexSizeManangement", value)
    if err != nil {
        panic(err)
    }
}
// SetSearchEnableRemoteQueries sets the searchEnableRemoteQueries property value. Indicates whether or not to block remote queries of this computer’s index.
func (m *Windows10GeneralConfiguration) SetSearchEnableRemoteQueries(value *bool)() {
    err := m.GetBackingStore().Set("searchEnableRemoteQueries", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockAccountsPage sets the settingsBlockAccountsPage property value. Indicates whether or not to block access to Accounts in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockAccountsPage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockAccountsPage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockAddProvisioningPackage sets the settingsBlockAddProvisioningPackage property value. Indicates whether or not to block the user from installing provisioning packages.
func (m *Windows10GeneralConfiguration) SetSettingsBlockAddProvisioningPackage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockAddProvisioningPackage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockAppsPage sets the settingsBlockAppsPage property value. Indicates whether or not to block access to Apps in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockAppsPage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockAppsPage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockChangeLanguage sets the settingsBlockChangeLanguage property value. Indicates whether or not to block the user from changing the language settings.
func (m *Windows10GeneralConfiguration) SetSettingsBlockChangeLanguage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockChangeLanguage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockChangePowerSleep sets the settingsBlockChangePowerSleep property value. Indicates whether or not to block the user from changing power and sleep settings.
func (m *Windows10GeneralConfiguration) SetSettingsBlockChangePowerSleep(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockChangePowerSleep", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockChangeRegion sets the settingsBlockChangeRegion property value. Indicates whether or not to block the user from changing the region settings.
func (m *Windows10GeneralConfiguration) SetSettingsBlockChangeRegion(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockChangeRegion", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockChangeSystemTime sets the settingsBlockChangeSystemTime property value. Indicates whether or not to block the user from changing date and time settings.
func (m *Windows10GeneralConfiguration) SetSettingsBlockChangeSystemTime(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockChangeSystemTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockDevicesPage sets the settingsBlockDevicesPage property value. Indicates whether or not to block access to Devices in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockDevicesPage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockDevicesPage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockEaseOfAccessPage sets the settingsBlockEaseOfAccessPage property value. Indicates whether or not to block access to Ease of Access in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockEaseOfAccessPage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockEaseOfAccessPage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockEditDeviceName sets the settingsBlockEditDeviceName property value. Indicates whether or not to block the user from editing the device name.
func (m *Windows10GeneralConfiguration) SetSettingsBlockEditDeviceName(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockEditDeviceName", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockGamingPage sets the settingsBlockGamingPage property value. Indicates whether or not to block access to Gaming in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockGamingPage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockGamingPage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockNetworkInternetPage sets the settingsBlockNetworkInternetPage property value. Indicates whether or not to block access to Network & Internet in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockNetworkInternetPage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockNetworkInternetPage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockPersonalizationPage sets the settingsBlockPersonalizationPage property value. Indicates whether or not to block access to Personalization in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockPersonalizationPage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockPersonalizationPage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockPrivacyPage sets the settingsBlockPrivacyPage property value. Indicates whether or not to block access to Privacy in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockPrivacyPage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockPrivacyPage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockRemoveProvisioningPackage sets the settingsBlockRemoveProvisioningPackage property value. Indicates whether or not to block the runtime configuration agent from removing provisioning packages.
func (m *Windows10GeneralConfiguration) SetSettingsBlockRemoveProvisioningPackage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockRemoveProvisioningPackage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockSettingsApp sets the settingsBlockSettingsApp property value. Indicates whether or not to block access to Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockSettingsApp(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockSettingsApp", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockSystemPage sets the settingsBlockSystemPage property value. Indicates whether or not to block access to System in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockSystemPage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockSystemPage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockTimeLanguagePage sets the settingsBlockTimeLanguagePage property value. Indicates whether or not to block access to Time & Language in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockTimeLanguagePage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockTimeLanguagePage", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockUpdateSecurityPage sets the settingsBlockUpdateSecurityPage property value. Indicates whether or not to block access to Update & Security in Settings app.
func (m *Windows10GeneralConfiguration) SetSettingsBlockUpdateSecurityPage(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockUpdateSecurityPage", value)
    if err != nil {
        panic(err)
    }
}
// SetSharedUserAppDataAllowed sets the sharedUserAppDataAllowed property value. Indicates whether or not to block multiple users of the same app to share data.
func (m *Windows10GeneralConfiguration) SetSharedUserAppDataAllowed(value *bool)() {
    err := m.GetBackingStore().Set("sharedUserAppDataAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetSmartScreenBlockPromptOverride sets the smartScreenBlockPromptOverride property value. Indicates whether or not users can override SmartScreen Filter warnings about potentially malicious websites.
func (m *Windows10GeneralConfiguration) SetSmartScreenBlockPromptOverride(value *bool)() {
    err := m.GetBackingStore().Set("smartScreenBlockPromptOverride", value)
    if err != nil {
        panic(err)
    }
}
// SetSmartScreenBlockPromptOverrideForFiles sets the smartScreenBlockPromptOverrideForFiles property value. Indicates whether or not users can override the SmartScreen Filter warnings about downloading unverified files
func (m *Windows10GeneralConfiguration) SetSmartScreenBlockPromptOverrideForFiles(value *bool)() {
    err := m.GetBackingStore().Set("smartScreenBlockPromptOverrideForFiles", value)
    if err != nil {
        panic(err)
    }
}
// SetSmartScreenEnableAppInstallControl sets the smartScreenEnableAppInstallControl property value. This property will be deprecated in July 2019 and will be replaced by property SmartScreenAppInstallControl. Allows IT Admins to control whether users are allowed to install apps from places other than the Store.
func (m *Windows10GeneralConfiguration) SetSmartScreenEnableAppInstallControl(value *bool)() {
    err := m.GetBackingStore().Set("smartScreenEnableAppInstallControl", value)
    if err != nil {
        panic(err)
    }
}
// SetStartBlockUnpinningAppsFromTaskbar sets the startBlockUnpinningAppsFromTaskbar property value. Indicates whether or not to block the user from unpinning apps from taskbar.
func (m *Windows10GeneralConfiguration) SetStartBlockUnpinningAppsFromTaskbar(value *bool)() {
    err := m.GetBackingStore().Set("startBlockUnpinningAppsFromTaskbar", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuAppListVisibility sets the startMenuAppListVisibility property value. Type of start menu app list visibility.
func (m *Windows10GeneralConfiguration) SetStartMenuAppListVisibility(value *WindowsStartMenuAppListVisibilityType)() {
    err := m.GetBackingStore().Set("startMenuAppListVisibility", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideChangeAccountSettings sets the startMenuHideChangeAccountSettings property value. Enabling this policy hides the change account setting from appearing in the user tile in the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuHideChangeAccountSettings(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideChangeAccountSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideFrequentlyUsedApps sets the startMenuHideFrequentlyUsedApps property value. Enabling this policy hides the most used apps from appearing on the start menu and disables the corresponding toggle in the Settings app.
func (m *Windows10GeneralConfiguration) SetStartMenuHideFrequentlyUsedApps(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideFrequentlyUsedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideHibernate sets the startMenuHideHibernate property value. Enabling this policy hides hibernate from appearing in the power button in the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuHideHibernate(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideHibernate", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideLock sets the startMenuHideLock property value. Enabling this policy hides lock from appearing in the user tile in the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuHideLock(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideLock", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHidePowerButton sets the startMenuHidePowerButton property value. Enabling this policy hides the power button from appearing in the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuHidePowerButton(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHidePowerButton", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideRecentJumpLists sets the startMenuHideRecentJumpLists property value. Enabling this policy hides recent jump lists from appearing on the start menu/taskbar and disables the corresponding toggle in the Settings app.
func (m *Windows10GeneralConfiguration) SetStartMenuHideRecentJumpLists(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideRecentJumpLists", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideRecentlyAddedApps sets the startMenuHideRecentlyAddedApps property value. Enabling this policy hides recently added apps from appearing on the start menu and disables the corresponding toggle in the Settings app.
func (m *Windows10GeneralConfiguration) SetStartMenuHideRecentlyAddedApps(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideRecentlyAddedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideRestartOptions sets the startMenuHideRestartOptions property value. Enabling this policy hides 'Restart/Update and Restart' from appearing in the power button in the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuHideRestartOptions(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideRestartOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideShutDown sets the startMenuHideShutDown property value. Enabling this policy hides shut down/update and shut down from appearing in the power button in the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuHideShutDown(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideShutDown", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideSignOut sets the startMenuHideSignOut property value. Enabling this policy hides sign out from appearing in the user tile in the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuHideSignOut(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideSignOut", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideSleep sets the startMenuHideSleep property value. Enabling this policy hides sleep from appearing in the power button in the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuHideSleep(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideSleep", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideSwitchAccount sets the startMenuHideSwitchAccount property value. Enabling this policy hides switch account from appearing in the user tile in the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuHideSwitchAccount(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideSwitchAccount", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuHideUserTile sets the startMenuHideUserTile property value. Enabling this policy hides the user tile from appearing in the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuHideUserTile(value *bool)() {
    err := m.GetBackingStore().Set("startMenuHideUserTile", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuLayoutEdgeAssetsXml sets the startMenuLayoutEdgeAssetsXml property value. This policy setting allows you to import Edge assets to be used with startMenuLayoutXml policy. Start layout can contain secondary tile from Edge app which looks for Edge local asset file. Edge local asset would not exist and cause Edge secondary tile to appear empty in this case. This policy only gets applied when startMenuLayoutXml policy is modified. The value should be a UTF-8 Base64 encoded byte array.
func (m *Windows10GeneralConfiguration) SetStartMenuLayoutEdgeAssetsXml(value []byte)() {
    err := m.GetBackingStore().Set("startMenuLayoutEdgeAssetsXml", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuLayoutXml sets the startMenuLayoutXml property value. Allows admins to override the default Start menu layout and prevents the user from changing it. The layout is modified by specifying an XML file based on a layout modification schema. XML needs to be in a UTF8 encoded byte array format.
func (m *Windows10GeneralConfiguration) SetStartMenuLayoutXml(value []byte)() {
    err := m.GetBackingStore().Set("startMenuLayoutXml", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuMode sets the startMenuMode property value. Type of display modes for the start menu.
func (m *Windows10GeneralConfiguration) SetStartMenuMode(value *WindowsStartMenuModeType)() {
    err := m.GetBackingStore().Set("startMenuMode", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuPinnedFolderDocuments sets the startMenuPinnedFolderDocuments property value. Generic visibility state.
func (m *Windows10GeneralConfiguration) SetStartMenuPinnedFolderDocuments(value *VisibilitySetting)() {
    err := m.GetBackingStore().Set("startMenuPinnedFolderDocuments", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuPinnedFolderDownloads sets the startMenuPinnedFolderDownloads property value. Generic visibility state.
func (m *Windows10GeneralConfiguration) SetStartMenuPinnedFolderDownloads(value *VisibilitySetting)() {
    err := m.GetBackingStore().Set("startMenuPinnedFolderDownloads", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuPinnedFolderFileExplorer sets the startMenuPinnedFolderFileExplorer property value. Generic visibility state.
func (m *Windows10GeneralConfiguration) SetStartMenuPinnedFolderFileExplorer(value *VisibilitySetting)() {
    err := m.GetBackingStore().Set("startMenuPinnedFolderFileExplorer", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuPinnedFolderHomeGroup sets the startMenuPinnedFolderHomeGroup property value. Generic visibility state.
func (m *Windows10GeneralConfiguration) SetStartMenuPinnedFolderHomeGroup(value *VisibilitySetting)() {
    err := m.GetBackingStore().Set("startMenuPinnedFolderHomeGroup", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuPinnedFolderMusic sets the startMenuPinnedFolderMusic property value. Generic visibility state.
func (m *Windows10GeneralConfiguration) SetStartMenuPinnedFolderMusic(value *VisibilitySetting)() {
    err := m.GetBackingStore().Set("startMenuPinnedFolderMusic", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuPinnedFolderNetwork sets the startMenuPinnedFolderNetwork property value. Generic visibility state.
func (m *Windows10GeneralConfiguration) SetStartMenuPinnedFolderNetwork(value *VisibilitySetting)() {
    err := m.GetBackingStore().Set("startMenuPinnedFolderNetwork", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuPinnedFolderPersonalFolder sets the startMenuPinnedFolderPersonalFolder property value. Generic visibility state.
func (m *Windows10GeneralConfiguration) SetStartMenuPinnedFolderPersonalFolder(value *VisibilitySetting)() {
    err := m.GetBackingStore().Set("startMenuPinnedFolderPersonalFolder", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuPinnedFolderPictures sets the startMenuPinnedFolderPictures property value. Generic visibility state.
func (m *Windows10GeneralConfiguration) SetStartMenuPinnedFolderPictures(value *VisibilitySetting)() {
    err := m.GetBackingStore().Set("startMenuPinnedFolderPictures", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuPinnedFolderSettings sets the startMenuPinnedFolderSettings property value. Generic visibility state.
func (m *Windows10GeneralConfiguration) SetStartMenuPinnedFolderSettings(value *VisibilitySetting)() {
    err := m.GetBackingStore().Set("startMenuPinnedFolderSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetStartMenuPinnedFolderVideos sets the startMenuPinnedFolderVideos property value. Generic visibility state.
func (m *Windows10GeneralConfiguration) SetStartMenuPinnedFolderVideos(value *VisibilitySetting)() {
    err := m.GetBackingStore().Set("startMenuPinnedFolderVideos", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageBlockRemovableStorage sets the storageBlockRemovableStorage property value. Indicates whether or not to Block the user from using removable storage.
func (m *Windows10GeneralConfiguration) SetStorageBlockRemovableStorage(value *bool)() {
    err := m.GetBackingStore().Set("storageBlockRemovableStorage", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageRequireMobileDeviceEncryption sets the storageRequireMobileDeviceEncryption property value. Indicating whether or not to require encryption on a mobile device.
func (m *Windows10GeneralConfiguration) SetStorageRequireMobileDeviceEncryption(value *bool)() {
    err := m.GetBackingStore().Set("storageRequireMobileDeviceEncryption", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageRestrictAppDataToSystemVolume sets the storageRestrictAppDataToSystemVolume property value. Indicates whether application data is restricted to the system drive.
func (m *Windows10GeneralConfiguration) SetStorageRestrictAppDataToSystemVolume(value *bool)() {
    err := m.GetBackingStore().Set("storageRestrictAppDataToSystemVolume", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageRestrictAppInstallToSystemVolume sets the storageRestrictAppInstallToSystemVolume property value. Indicates whether the installation of applications is restricted to the system drive.
func (m *Windows10GeneralConfiguration) SetStorageRestrictAppInstallToSystemVolume(value *bool)() {
    err := m.GetBackingStore().Set("storageRestrictAppInstallToSystemVolume", value)
    if err != nil {
        panic(err)
    }
}
// SetTenantLockdownRequireNetworkDuringOutOfBoxExperience sets the tenantLockdownRequireNetworkDuringOutOfBoxExperience property value. Whether the device is required to connect to the network.
func (m *Windows10GeneralConfiguration) SetTenantLockdownRequireNetworkDuringOutOfBoxExperience(value *bool)() {
    err := m.GetBackingStore().Set("tenantLockdownRequireNetworkDuringOutOfBoxExperience", value)
    if err != nil {
        panic(err)
    }
}
// SetUsbBlocked sets the usbBlocked property value. Indicates whether or not to Block the user from USB connection.
func (m *Windows10GeneralConfiguration) SetUsbBlocked(value *bool)() {
    err := m.GetBackingStore().Set("usbBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetVoiceRecordingBlocked sets the voiceRecordingBlocked property value. Indicates whether or not to Block the user from voice recording.
func (m *Windows10GeneralConfiguration) SetVoiceRecordingBlocked(value *bool)() {
    err := m.GetBackingStore().Set("voiceRecordingBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetWebRtcBlockLocalhostIpAddress sets the webRtcBlockLocalhostIpAddress property value. Indicates whether or not user's localhost IP address is displayed while making phone calls using the WebRTC
func (m *Windows10GeneralConfiguration) SetWebRtcBlockLocalhostIpAddress(value *bool)() {
    err := m.GetBackingStore().Set("webRtcBlockLocalhostIpAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetWiFiBlockAutomaticConnectHotspots sets the wiFiBlockAutomaticConnectHotspots property value. Indicating whether or not to block automatically connecting to Wi-Fi hotspots. Has no impact if Wi-Fi is blocked.
func (m *Windows10GeneralConfiguration) SetWiFiBlockAutomaticConnectHotspots(value *bool)() {
    err := m.GetBackingStore().Set("wiFiBlockAutomaticConnectHotspots", value)
    if err != nil {
        panic(err)
    }
}
// SetWiFiBlocked sets the wiFiBlocked property value. Indicates whether or not to Block the user from using Wi-Fi.
func (m *Windows10GeneralConfiguration) SetWiFiBlocked(value *bool)() {
    err := m.GetBackingStore().Set("wiFiBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetWiFiBlockManualConfiguration sets the wiFiBlockManualConfiguration property value. Indicates whether or not to Block the user from using Wi-Fi manual configuration.
func (m *Windows10GeneralConfiguration) SetWiFiBlockManualConfiguration(value *bool)() {
    err := m.GetBackingStore().Set("wiFiBlockManualConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetWiFiScanInterval sets the wiFiScanInterval property value. Specify how often devices scan for Wi-Fi networks. Supported values are 1-500, where 100 = default, and 500 = low frequency. Valid values 1 to 500
func (m *Windows10GeneralConfiguration) SetWiFiScanInterval(value *int32)() {
    err := m.GetBackingStore().Set("wiFiScanInterval", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsSpotlightBlockConsumerSpecificFeatures sets the windowsSpotlightBlockConsumerSpecificFeatures property value. Allows IT admins to block experiences that are typically for consumers only, such as Start suggestions, Membership notifications, Post-OOBE app install and redirect tiles.
func (m *Windows10GeneralConfiguration) SetWindowsSpotlightBlockConsumerSpecificFeatures(value *bool)() {
    err := m.GetBackingStore().Set("windowsSpotlightBlockConsumerSpecificFeatures", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsSpotlightBlocked sets the windowsSpotlightBlocked property value. Allows IT admins to turn off all Windows Spotlight features
func (m *Windows10GeneralConfiguration) SetWindowsSpotlightBlocked(value *bool)() {
    err := m.GetBackingStore().Set("windowsSpotlightBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsSpotlightBlockOnActionCenter sets the windowsSpotlightBlockOnActionCenter property value. Block suggestions from Microsoft that show after each OS clean install, upgrade or in an on-going basis to introduce users to what is new or changed
func (m *Windows10GeneralConfiguration) SetWindowsSpotlightBlockOnActionCenter(value *bool)() {
    err := m.GetBackingStore().Set("windowsSpotlightBlockOnActionCenter", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsSpotlightBlockTailoredExperiences sets the windowsSpotlightBlockTailoredExperiences property value. Block personalized content in Windows spotlight based on user’s device usage.
func (m *Windows10GeneralConfiguration) SetWindowsSpotlightBlockTailoredExperiences(value *bool)() {
    err := m.GetBackingStore().Set("windowsSpotlightBlockTailoredExperiences", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsSpotlightBlockThirdPartyNotifications sets the windowsSpotlightBlockThirdPartyNotifications property value. Block third party content delivered via Windows Spotlight
func (m *Windows10GeneralConfiguration) SetWindowsSpotlightBlockThirdPartyNotifications(value *bool)() {
    err := m.GetBackingStore().Set("windowsSpotlightBlockThirdPartyNotifications", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsSpotlightBlockWelcomeExperience sets the windowsSpotlightBlockWelcomeExperience property value. Block Windows Spotlight Windows welcome experience
func (m *Windows10GeneralConfiguration) SetWindowsSpotlightBlockWelcomeExperience(value *bool)() {
    err := m.GetBackingStore().Set("windowsSpotlightBlockWelcomeExperience", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsSpotlightBlockWindowsTips sets the windowsSpotlightBlockWindowsTips property value. Allows IT admins to turn off the popup of Windows Tips.
func (m *Windows10GeneralConfiguration) SetWindowsSpotlightBlockWindowsTips(value *bool)() {
    err := m.GetBackingStore().Set("windowsSpotlightBlockWindowsTips", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsSpotlightConfigureOnLockScreen sets the windowsSpotlightConfigureOnLockScreen property value. Allows IT admind to set a predefined default search engine for MDM-Controlled devices
func (m *Windows10GeneralConfiguration) SetWindowsSpotlightConfigureOnLockScreen(value *WindowsSpotlightEnablementSettings)() {
    err := m.GetBackingStore().Set("windowsSpotlightConfigureOnLockScreen", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsStoreBlockAutoUpdate sets the windowsStoreBlockAutoUpdate property value. Indicates whether or not to block automatic update of apps from Windows Store.
func (m *Windows10GeneralConfiguration) SetWindowsStoreBlockAutoUpdate(value *bool)() {
    err := m.GetBackingStore().Set("windowsStoreBlockAutoUpdate", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsStoreBlocked sets the windowsStoreBlocked property value. Indicates whether or not to Block the user from using the Windows store.
func (m *Windows10GeneralConfiguration) SetWindowsStoreBlocked(value *bool)() {
    err := m.GetBackingStore().Set("windowsStoreBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsStoreEnablePrivateStoreOnly sets the windowsStoreEnablePrivateStoreOnly property value. Indicates whether or not to enable Private Store Only.
func (m *Windows10GeneralConfiguration) SetWindowsStoreEnablePrivateStoreOnly(value *bool)() {
    err := m.GetBackingStore().Set("windowsStoreEnablePrivateStoreOnly", value)
    if err != nil {
        panic(err)
    }
}
// SetWirelessDisplayBlockProjectionToThisDevice sets the wirelessDisplayBlockProjectionToThisDevice property value. Indicates whether or not to allow other devices from discovering this PC for projection.
func (m *Windows10GeneralConfiguration) SetWirelessDisplayBlockProjectionToThisDevice(value *bool)() {
    err := m.GetBackingStore().Set("wirelessDisplayBlockProjectionToThisDevice", value)
    if err != nil {
        panic(err)
    }
}
// SetWirelessDisplayBlockUserInputFromReceiver sets the wirelessDisplayBlockUserInputFromReceiver property value. Indicates whether or not to allow user input from wireless display receiver.
func (m *Windows10GeneralConfiguration) SetWirelessDisplayBlockUserInputFromReceiver(value *bool)() {
    err := m.GetBackingStore().Set("wirelessDisplayBlockUserInputFromReceiver", value)
    if err != nil {
        panic(err)
    }
}
// SetWirelessDisplayRequirePinForPairing sets the wirelessDisplayRequirePinForPairing property value. Indicates whether or not to require a PIN for new devices to initiate pairing.
func (m *Windows10GeneralConfiguration) SetWirelessDisplayRequirePinForPairing(value *bool)() {
    err := m.GetBackingStore().Set("wirelessDisplayRequirePinForPairing", value)
    if err != nil {
        panic(err)
    }
}
type Windows10GeneralConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccountsBlockAddingNonMicrosoftAccountEmail()(*bool)
    GetAntiTheftModeBlocked()(*bool)
    GetAppsAllowTrustedAppsSideloading()(*StateManagementSetting)
    GetAppsBlockWindowsStoreOriginatedApps()(*bool)
    GetBluetoothAllowedServices()([]string)
    GetBluetoothBlockAdvertising()(*bool)
    GetBluetoothBlockDiscoverableMode()(*bool)
    GetBluetoothBlocked()(*bool)
    GetBluetoothBlockPrePairing()(*bool)
    GetCameraBlocked()(*bool)
    GetCellularBlockDataWhenRoaming()(*bool)
    GetCellularBlockVpn()(*bool)
    GetCellularBlockVpnWhenRoaming()(*bool)
    GetCertificatesBlockManualRootCertificateInstallation()(*bool)
    GetConnectedDevicesServiceBlocked()(*bool)
    GetCopyPasteBlocked()(*bool)
    GetCortanaBlocked()(*bool)
    GetDefenderBlockEndUserAccess()(*bool)
    GetDefenderCloudBlockLevel()(*DefenderCloudBlockLevelType)
    GetDefenderDaysBeforeDeletingQuarantinedMalware()(*int32)
    GetDefenderDetectedMalwareActions()(DefenderDetectedMalwareActionsable)
    GetDefenderFileExtensionsToExclude()([]string)
    GetDefenderFilesAndFoldersToExclude()([]string)
    GetDefenderMonitorFileActivity()(*DefenderMonitorFileActivity)
    GetDefenderProcessesToExclude()([]string)
    GetDefenderPromptForSampleSubmission()(*DefenderPromptForSampleSubmission)
    GetDefenderRequireBehaviorMonitoring()(*bool)
    GetDefenderRequireCloudProtection()(*bool)
    GetDefenderRequireNetworkInspectionSystem()(*bool)
    GetDefenderRequireRealTimeMonitoring()(*bool)
    GetDefenderScanArchiveFiles()(*bool)
    GetDefenderScanDownloads()(*bool)
    GetDefenderScanIncomingMail()(*bool)
    GetDefenderScanMappedNetworkDrivesDuringFullScan()(*bool)
    GetDefenderScanMaxCpu()(*int32)
    GetDefenderScanNetworkFiles()(*bool)
    GetDefenderScanRemovableDrivesDuringFullScan()(*bool)
    GetDefenderScanScriptsLoadedInInternetExplorer()(*bool)
    GetDefenderScanType()(*DefenderScanType)
    GetDefenderScheduledQuickScanTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)
    GetDefenderScheduledScanTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)
    GetDefenderSignatureUpdateIntervalInHours()(*int32)
    GetDefenderSystemScanSchedule()(*WeeklySchedule)
    GetDeveloperUnlockSetting()(*StateManagementSetting)
    GetDeviceManagementBlockFactoryResetOnMobile()(*bool)
    GetDeviceManagementBlockManualUnenroll()(*bool)
    GetDiagnosticsDataSubmissionMode()(*DiagnosticDataSubmissionMode)
    GetEdgeAllowStartPagesModification()(*bool)
    GetEdgeBlockAccessToAboutFlags()(*bool)
    GetEdgeBlockAddressBarDropdown()(*bool)
    GetEdgeBlockAutofill()(*bool)
    GetEdgeBlockCompatibilityList()(*bool)
    GetEdgeBlockDeveloperTools()(*bool)
    GetEdgeBlocked()(*bool)
    GetEdgeBlockExtensions()(*bool)
    GetEdgeBlockInPrivateBrowsing()(*bool)
    GetEdgeBlockJavaScript()(*bool)
    GetEdgeBlockLiveTileDataCollection()(*bool)
    GetEdgeBlockPasswordManager()(*bool)
    GetEdgeBlockPopups()(*bool)
    GetEdgeBlockSearchSuggestions()(*bool)
    GetEdgeBlockSendingDoNotTrackHeader()(*bool)
    GetEdgeBlockSendingIntranetTrafficToInternetExplorer()(*bool)
    GetEdgeClearBrowsingDataOnExit()(*bool)
    GetEdgeCookiePolicy()(*EdgeCookiePolicy)
    GetEdgeDisableFirstRunPage()(*bool)
    GetEdgeEnterpriseModeSiteListLocation()(*string)
    GetEdgeFirstRunUrl()(*string)
    GetEdgeHomepageUrls()([]string)
    GetEdgeRequireSmartScreen()(*bool)
    GetEdgeSearchEngine()(EdgeSearchEngineBaseable)
    GetEdgeSendIntranetTrafficToInternetExplorer()(*bool)
    GetEdgeSyncFavoritesWithInternetExplorer()(*bool)
    GetEnterpriseCloudPrintDiscoveryEndPoint()(*string)
    GetEnterpriseCloudPrintDiscoveryMaxLimit()(*int32)
    GetEnterpriseCloudPrintMopriaDiscoveryResourceIdentifier()(*string)
    GetEnterpriseCloudPrintOAuthAuthority()(*string)
    GetEnterpriseCloudPrintOAuthClientIdentifier()(*string)
    GetEnterpriseCloudPrintResourceIdentifier()(*string)
    GetExperienceBlockDeviceDiscovery()(*bool)
    GetExperienceBlockErrorDialogWhenNoSIM()(*bool)
    GetExperienceBlockTaskSwitcher()(*bool)
    GetGameDvrBlocked()(*bool)
    GetInternetSharingBlocked()(*bool)
    GetLocationServicesBlocked()(*bool)
    GetLockScreenAllowTimeoutConfiguration()(*bool)
    GetLockScreenBlockActionCenterNotifications()(*bool)
    GetLockScreenBlockCortana()(*bool)
    GetLockScreenBlockToastNotifications()(*bool)
    GetLockScreenTimeoutInSeconds()(*int32)
    GetLogonBlockFastUserSwitching()(*bool)
    GetMicrosoftAccountBlocked()(*bool)
    GetMicrosoftAccountBlockSettingsSync()(*bool)
    GetNetworkProxyApplySettingsDeviceWide()(*bool)
    GetNetworkProxyAutomaticConfigurationUrl()(*string)
    GetNetworkProxyDisableAutoDetect()(*bool)
    GetNetworkProxyServer()(Windows10NetworkProxyServerable)
    GetNfcBlocked()(*bool)
    GetOneDriveDisableFileSync()(*bool)
    GetPasswordBlockSimple()(*bool)
    GetPasswordExpirationDays()(*int32)
    GetPasswordMinimumCharacterSetCount()(*int32)
    GetPasswordMinimumLength()(*int32)
    GetPasswordMinutesOfInactivityBeforeScreenTimeout()(*int32)
    GetPasswordPreviousPasswordBlockCount()(*int32)
    GetPasswordRequired()(*bool)
    GetPasswordRequiredType()(*RequiredPasswordType)
    GetPasswordRequireWhenResumeFromIdleState()(*bool)
    GetPasswordSignInFailureCountBeforeFactoryReset()(*int32)
    GetPersonalizationDesktopImageUrl()(*string)
    GetPersonalizationLockScreenImageUrl()(*string)
    GetPrivacyAdvertisingId()(*StateManagementSetting)
    GetPrivacyAutoAcceptPairingAndConsentPrompts()(*bool)
    GetPrivacyBlockInputPersonalization()(*bool)
    GetResetProtectionModeBlocked()(*bool)
    GetSafeSearchFilter()(*SafeSearchFilterType)
    GetScreenCaptureBlocked()(*bool)
    GetSearchBlockDiacritics()(*bool)
    GetSearchDisableAutoLanguageDetection()(*bool)
    GetSearchDisableIndexerBackoff()(*bool)
    GetSearchDisableIndexingEncryptedItems()(*bool)
    GetSearchDisableIndexingRemovableDrive()(*bool)
    GetSearchEnableAutomaticIndexSizeManangement()(*bool)
    GetSearchEnableRemoteQueries()(*bool)
    GetSettingsBlockAccountsPage()(*bool)
    GetSettingsBlockAddProvisioningPackage()(*bool)
    GetSettingsBlockAppsPage()(*bool)
    GetSettingsBlockChangeLanguage()(*bool)
    GetSettingsBlockChangePowerSleep()(*bool)
    GetSettingsBlockChangeRegion()(*bool)
    GetSettingsBlockChangeSystemTime()(*bool)
    GetSettingsBlockDevicesPage()(*bool)
    GetSettingsBlockEaseOfAccessPage()(*bool)
    GetSettingsBlockEditDeviceName()(*bool)
    GetSettingsBlockGamingPage()(*bool)
    GetSettingsBlockNetworkInternetPage()(*bool)
    GetSettingsBlockPersonalizationPage()(*bool)
    GetSettingsBlockPrivacyPage()(*bool)
    GetSettingsBlockRemoveProvisioningPackage()(*bool)
    GetSettingsBlockSettingsApp()(*bool)
    GetSettingsBlockSystemPage()(*bool)
    GetSettingsBlockTimeLanguagePage()(*bool)
    GetSettingsBlockUpdateSecurityPage()(*bool)
    GetSharedUserAppDataAllowed()(*bool)
    GetSmartScreenBlockPromptOverride()(*bool)
    GetSmartScreenBlockPromptOverrideForFiles()(*bool)
    GetSmartScreenEnableAppInstallControl()(*bool)
    GetStartBlockUnpinningAppsFromTaskbar()(*bool)
    GetStartMenuAppListVisibility()(*WindowsStartMenuAppListVisibilityType)
    GetStartMenuHideChangeAccountSettings()(*bool)
    GetStartMenuHideFrequentlyUsedApps()(*bool)
    GetStartMenuHideHibernate()(*bool)
    GetStartMenuHideLock()(*bool)
    GetStartMenuHidePowerButton()(*bool)
    GetStartMenuHideRecentJumpLists()(*bool)
    GetStartMenuHideRecentlyAddedApps()(*bool)
    GetStartMenuHideRestartOptions()(*bool)
    GetStartMenuHideShutDown()(*bool)
    GetStartMenuHideSignOut()(*bool)
    GetStartMenuHideSleep()(*bool)
    GetStartMenuHideSwitchAccount()(*bool)
    GetStartMenuHideUserTile()(*bool)
    GetStartMenuLayoutEdgeAssetsXml()([]byte)
    GetStartMenuLayoutXml()([]byte)
    GetStartMenuMode()(*WindowsStartMenuModeType)
    GetStartMenuPinnedFolderDocuments()(*VisibilitySetting)
    GetStartMenuPinnedFolderDownloads()(*VisibilitySetting)
    GetStartMenuPinnedFolderFileExplorer()(*VisibilitySetting)
    GetStartMenuPinnedFolderHomeGroup()(*VisibilitySetting)
    GetStartMenuPinnedFolderMusic()(*VisibilitySetting)
    GetStartMenuPinnedFolderNetwork()(*VisibilitySetting)
    GetStartMenuPinnedFolderPersonalFolder()(*VisibilitySetting)
    GetStartMenuPinnedFolderPictures()(*VisibilitySetting)
    GetStartMenuPinnedFolderSettings()(*VisibilitySetting)
    GetStartMenuPinnedFolderVideos()(*VisibilitySetting)
    GetStorageBlockRemovableStorage()(*bool)
    GetStorageRequireMobileDeviceEncryption()(*bool)
    GetStorageRestrictAppDataToSystemVolume()(*bool)
    GetStorageRestrictAppInstallToSystemVolume()(*bool)
    GetTenantLockdownRequireNetworkDuringOutOfBoxExperience()(*bool)
    GetUsbBlocked()(*bool)
    GetVoiceRecordingBlocked()(*bool)
    GetWebRtcBlockLocalhostIpAddress()(*bool)
    GetWiFiBlockAutomaticConnectHotspots()(*bool)
    GetWiFiBlocked()(*bool)
    GetWiFiBlockManualConfiguration()(*bool)
    GetWiFiScanInterval()(*int32)
    GetWindowsSpotlightBlockConsumerSpecificFeatures()(*bool)
    GetWindowsSpotlightBlocked()(*bool)
    GetWindowsSpotlightBlockOnActionCenter()(*bool)
    GetWindowsSpotlightBlockTailoredExperiences()(*bool)
    GetWindowsSpotlightBlockThirdPartyNotifications()(*bool)
    GetWindowsSpotlightBlockWelcomeExperience()(*bool)
    GetWindowsSpotlightBlockWindowsTips()(*bool)
    GetWindowsSpotlightConfigureOnLockScreen()(*WindowsSpotlightEnablementSettings)
    GetWindowsStoreBlockAutoUpdate()(*bool)
    GetWindowsStoreBlocked()(*bool)
    GetWindowsStoreEnablePrivateStoreOnly()(*bool)
    GetWirelessDisplayBlockProjectionToThisDevice()(*bool)
    GetWirelessDisplayBlockUserInputFromReceiver()(*bool)
    GetWirelessDisplayRequirePinForPairing()(*bool)
    SetAccountsBlockAddingNonMicrosoftAccountEmail(value *bool)()
    SetAntiTheftModeBlocked(value *bool)()
    SetAppsAllowTrustedAppsSideloading(value *StateManagementSetting)()
    SetAppsBlockWindowsStoreOriginatedApps(value *bool)()
    SetBluetoothAllowedServices(value []string)()
    SetBluetoothBlockAdvertising(value *bool)()
    SetBluetoothBlockDiscoverableMode(value *bool)()
    SetBluetoothBlocked(value *bool)()
    SetBluetoothBlockPrePairing(value *bool)()
    SetCameraBlocked(value *bool)()
    SetCellularBlockDataWhenRoaming(value *bool)()
    SetCellularBlockVpn(value *bool)()
    SetCellularBlockVpnWhenRoaming(value *bool)()
    SetCertificatesBlockManualRootCertificateInstallation(value *bool)()
    SetConnectedDevicesServiceBlocked(value *bool)()
    SetCopyPasteBlocked(value *bool)()
    SetCortanaBlocked(value *bool)()
    SetDefenderBlockEndUserAccess(value *bool)()
    SetDefenderCloudBlockLevel(value *DefenderCloudBlockLevelType)()
    SetDefenderDaysBeforeDeletingQuarantinedMalware(value *int32)()
    SetDefenderDetectedMalwareActions(value DefenderDetectedMalwareActionsable)()
    SetDefenderFileExtensionsToExclude(value []string)()
    SetDefenderFilesAndFoldersToExclude(value []string)()
    SetDefenderMonitorFileActivity(value *DefenderMonitorFileActivity)()
    SetDefenderProcessesToExclude(value []string)()
    SetDefenderPromptForSampleSubmission(value *DefenderPromptForSampleSubmission)()
    SetDefenderRequireBehaviorMonitoring(value *bool)()
    SetDefenderRequireCloudProtection(value *bool)()
    SetDefenderRequireNetworkInspectionSystem(value *bool)()
    SetDefenderRequireRealTimeMonitoring(value *bool)()
    SetDefenderScanArchiveFiles(value *bool)()
    SetDefenderScanDownloads(value *bool)()
    SetDefenderScanIncomingMail(value *bool)()
    SetDefenderScanMappedNetworkDrivesDuringFullScan(value *bool)()
    SetDefenderScanMaxCpu(value *int32)()
    SetDefenderScanNetworkFiles(value *bool)()
    SetDefenderScanRemovableDrivesDuringFullScan(value *bool)()
    SetDefenderScanScriptsLoadedInInternetExplorer(value *bool)()
    SetDefenderScanType(value *DefenderScanType)()
    SetDefenderScheduledQuickScanTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)()
    SetDefenderScheduledScanTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)()
    SetDefenderSignatureUpdateIntervalInHours(value *int32)()
    SetDefenderSystemScanSchedule(value *WeeklySchedule)()
    SetDeveloperUnlockSetting(value *StateManagementSetting)()
    SetDeviceManagementBlockFactoryResetOnMobile(value *bool)()
    SetDeviceManagementBlockManualUnenroll(value *bool)()
    SetDiagnosticsDataSubmissionMode(value *DiagnosticDataSubmissionMode)()
    SetEdgeAllowStartPagesModification(value *bool)()
    SetEdgeBlockAccessToAboutFlags(value *bool)()
    SetEdgeBlockAddressBarDropdown(value *bool)()
    SetEdgeBlockAutofill(value *bool)()
    SetEdgeBlockCompatibilityList(value *bool)()
    SetEdgeBlockDeveloperTools(value *bool)()
    SetEdgeBlocked(value *bool)()
    SetEdgeBlockExtensions(value *bool)()
    SetEdgeBlockInPrivateBrowsing(value *bool)()
    SetEdgeBlockJavaScript(value *bool)()
    SetEdgeBlockLiveTileDataCollection(value *bool)()
    SetEdgeBlockPasswordManager(value *bool)()
    SetEdgeBlockPopups(value *bool)()
    SetEdgeBlockSearchSuggestions(value *bool)()
    SetEdgeBlockSendingDoNotTrackHeader(value *bool)()
    SetEdgeBlockSendingIntranetTrafficToInternetExplorer(value *bool)()
    SetEdgeClearBrowsingDataOnExit(value *bool)()
    SetEdgeCookiePolicy(value *EdgeCookiePolicy)()
    SetEdgeDisableFirstRunPage(value *bool)()
    SetEdgeEnterpriseModeSiteListLocation(value *string)()
    SetEdgeFirstRunUrl(value *string)()
    SetEdgeHomepageUrls(value []string)()
    SetEdgeRequireSmartScreen(value *bool)()
    SetEdgeSearchEngine(value EdgeSearchEngineBaseable)()
    SetEdgeSendIntranetTrafficToInternetExplorer(value *bool)()
    SetEdgeSyncFavoritesWithInternetExplorer(value *bool)()
    SetEnterpriseCloudPrintDiscoveryEndPoint(value *string)()
    SetEnterpriseCloudPrintDiscoveryMaxLimit(value *int32)()
    SetEnterpriseCloudPrintMopriaDiscoveryResourceIdentifier(value *string)()
    SetEnterpriseCloudPrintOAuthAuthority(value *string)()
    SetEnterpriseCloudPrintOAuthClientIdentifier(value *string)()
    SetEnterpriseCloudPrintResourceIdentifier(value *string)()
    SetExperienceBlockDeviceDiscovery(value *bool)()
    SetExperienceBlockErrorDialogWhenNoSIM(value *bool)()
    SetExperienceBlockTaskSwitcher(value *bool)()
    SetGameDvrBlocked(value *bool)()
    SetInternetSharingBlocked(value *bool)()
    SetLocationServicesBlocked(value *bool)()
    SetLockScreenAllowTimeoutConfiguration(value *bool)()
    SetLockScreenBlockActionCenterNotifications(value *bool)()
    SetLockScreenBlockCortana(value *bool)()
    SetLockScreenBlockToastNotifications(value *bool)()
    SetLockScreenTimeoutInSeconds(value *int32)()
    SetLogonBlockFastUserSwitching(value *bool)()
    SetMicrosoftAccountBlocked(value *bool)()
    SetMicrosoftAccountBlockSettingsSync(value *bool)()
    SetNetworkProxyApplySettingsDeviceWide(value *bool)()
    SetNetworkProxyAutomaticConfigurationUrl(value *string)()
    SetNetworkProxyDisableAutoDetect(value *bool)()
    SetNetworkProxyServer(value Windows10NetworkProxyServerable)()
    SetNfcBlocked(value *bool)()
    SetOneDriveDisableFileSync(value *bool)()
    SetPasswordBlockSimple(value *bool)()
    SetPasswordExpirationDays(value *int32)()
    SetPasswordMinimumCharacterSetCount(value *int32)()
    SetPasswordMinimumLength(value *int32)()
    SetPasswordMinutesOfInactivityBeforeScreenTimeout(value *int32)()
    SetPasswordPreviousPasswordBlockCount(value *int32)()
    SetPasswordRequired(value *bool)()
    SetPasswordRequiredType(value *RequiredPasswordType)()
    SetPasswordRequireWhenResumeFromIdleState(value *bool)()
    SetPasswordSignInFailureCountBeforeFactoryReset(value *int32)()
    SetPersonalizationDesktopImageUrl(value *string)()
    SetPersonalizationLockScreenImageUrl(value *string)()
    SetPrivacyAdvertisingId(value *StateManagementSetting)()
    SetPrivacyAutoAcceptPairingAndConsentPrompts(value *bool)()
    SetPrivacyBlockInputPersonalization(value *bool)()
    SetResetProtectionModeBlocked(value *bool)()
    SetSafeSearchFilter(value *SafeSearchFilterType)()
    SetScreenCaptureBlocked(value *bool)()
    SetSearchBlockDiacritics(value *bool)()
    SetSearchDisableAutoLanguageDetection(value *bool)()
    SetSearchDisableIndexerBackoff(value *bool)()
    SetSearchDisableIndexingEncryptedItems(value *bool)()
    SetSearchDisableIndexingRemovableDrive(value *bool)()
    SetSearchEnableAutomaticIndexSizeManangement(value *bool)()
    SetSearchEnableRemoteQueries(value *bool)()
    SetSettingsBlockAccountsPage(value *bool)()
    SetSettingsBlockAddProvisioningPackage(value *bool)()
    SetSettingsBlockAppsPage(value *bool)()
    SetSettingsBlockChangeLanguage(value *bool)()
    SetSettingsBlockChangePowerSleep(value *bool)()
    SetSettingsBlockChangeRegion(value *bool)()
    SetSettingsBlockChangeSystemTime(value *bool)()
    SetSettingsBlockDevicesPage(value *bool)()
    SetSettingsBlockEaseOfAccessPage(value *bool)()
    SetSettingsBlockEditDeviceName(value *bool)()
    SetSettingsBlockGamingPage(value *bool)()
    SetSettingsBlockNetworkInternetPage(value *bool)()
    SetSettingsBlockPersonalizationPage(value *bool)()
    SetSettingsBlockPrivacyPage(value *bool)()
    SetSettingsBlockRemoveProvisioningPackage(value *bool)()
    SetSettingsBlockSettingsApp(value *bool)()
    SetSettingsBlockSystemPage(value *bool)()
    SetSettingsBlockTimeLanguagePage(value *bool)()
    SetSettingsBlockUpdateSecurityPage(value *bool)()
    SetSharedUserAppDataAllowed(value *bool)()
    SetSmartScreenBlockPromptOverride(value *bool)()
    SetSmartScreenBlockPromptOverrideForFiles(value *bool)()
    SetSmartScreenEnableAppInstallControl(value *bool)()
    SetStartBlockUnpinningAppsFromTaskbar(value *bool)()
    SetStartMenuAppListVisibility(value *WindowsStartMenuAppListVisibilityType)()
    SetStartMenuHideChangeAccountSettings(value *bool)()
    SetStartMenuHideFrequentlyUsedApps(value *bool)()
    SetStartMenuHideHibernate(value *bool)()
    SetStartMenuHideLock(value *bool)()
    SetStartMenuHidePowerButton(value *bool)()
    SetStartMenuHideRecentJumpLists(value *bool)()
    SetStartMenuHideRecentlyAddedApps(value *bool)()
    SetStartMenuHideRestartOptions(value *bool)()
    SetStartMenuHideShutDown(value *bool)()
    SetStartMenuHideSignOut(value *bool)()
    SetStartMenuHideSleep(value *bool)()
    SetStartMenuHideSwitchAccount(value *bool)()
    SetStartMenuHideUserTile(value *bool)()
    SetStartMenuLayoutEdgeAssetsXml(value []byte)()
    SetStartMenuLayoutXml(value []byte)()
    SetStartMenuMode(value *WindowsStartMenuModeType)()
    SetStartMenuPinnedFolderDocuments(value *VisibilitySetting)()
    SetStartMenuPinnedFolderDownloads(value *VisibilitySetting)()
    SetStartMenuPinnedFolderFileExplorer(value *VisibilitySetting)()
    SetStartMenuPinnedFolderHomeGroup(value *VisibilitySetting)()
    SetStartMenuPinnedFolderMusic(value *VisibilitySetting)()
    SetStartMenuPinnedFolderNetwork(value *VisibilitySetting)()
    SetStartMenuPinnedFolderPersonalFolder(value *VisibilitySetting)()
    SetStartMenuPinnedFolderPictures(value *VisibilitySetting)()
    SetStartMenuPinnedFolderSettings(value *VisibilitySetting)()
    SetStartMenuPinnedFolderVideos(value *VisibilitySetting)()
    SetStorageBlockRemovableStorage(value *bool)()
    SetStorageRequireMobileDeviceEncryption(value *bool)()
    SetStorageRestrictAppDataToSystemVolume(value *bool)()
    SetStorageRestrictAppInstallToSystemVolume(value *bool)()
    SetTenantLockdownRequireNetworkDuringOutOfBoxExperience(value *bool)()
    SetUsbBlocked(value *bool)()
    SetVoiceRecordingBlocked(value *bool)()
    SetWebRtcBlockLocalhostIpAddress(value *bool)()
    SetWiFiBlockAutomaticConnectHotspots(value *bool)()
    SetWiFiBlocked(value *bool)()
    SetWiFiBlockManualConfiguration(value *bool)()
    SetWiFiScanInterval(value *int32)()
    SetWindowsSpotlightBlockConsumerSpecificFeatures(value *bool)()
    SetWindowsSpotlightBlocked(value *bool)()
    SetWindowsSpotlightBlockOnActionCenter(value *bool)()
    SetWindowsSpotlightBlockTailoredExperiences(value *bool)()
    SetWindowsSpotlightBlockThirdPartyNotifications(value *bool)()
    SetWindowsSpotlightBlockWelcomeExperience(value *bool)()
    SetWindowsSpotlightBlockWindowsTips(value *bool)()
    SetWindowsSpotlightConfigureOnLockScreen(value *WindowsSpotlightEnablementSettings)()
    SetWindowsStoreBlockAutoUpdate(value *bool)()
    SetWindowsStoreBlocked(value *bool)()
    SetWindowsStoreEnablePrivateStoreOnly(value *bool)()
    SetWirelessDisplayBlockProjectionToThisDevice(value *bool)()
    SetWirelessDisplayBlockUserInputFromReceiver(value *bool)()
    SetWirelessDisplayRequirePinForPairing(value *bool)()
}
