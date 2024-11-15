package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AndroidGeneralDeviceConfiguration this topic provides descriptions of the declared methods, properties and relationships exposed by the androidGeneralDeviceConfiguration resource.
type AndroidGeneralDeviceConfiguration struct {
    DeviceConfiguration
}
// NewAndroidGeneralDeviceConfiguration instantiates a new AndroidGeneralDeviceConfiguration and sets the default values.
func NewAndroidGeneralDeviceConfiguration()(*AndroidGeneralDeviceConfiguration) {
    m := &AndroidGeneralDeviceConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.androidGeneralDeviceConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAndroidGeneralDeviceConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAndroidGeneralDeviceConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAndroidGeneralDeviceConfiguration(), nil
}
// GetAppsBlockClipboardSharing gets the appsBlockClipboardSharing property value. Indicates whether or not to block clipboard sharing to copy and paste between applications.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetAppsBlockClipboardSharing()(*bool) {
    val, err := m.GetBackingStore().Get("appsBlockClipboardSharing")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppsBlockCopyPaste gets the appsBlockCopyPaste property value. Indicates whether or not to block copy and paste within applications.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetAppsBlockCopyPaste()(*bool) {
    val, err := m.GetBackingStore().Get("appsBlockCopyPaste")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppsBlockYouTube gets the appsBlockYouTube property value. Indicates whether or not to block the YouTube app.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetAppsBlockYouTube()(*bool) {
    val, err := m.GetBackingStore().Get("appsBlockYouTube")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppsHideList gets the appsHideList property value. List of apps to be hidden on the KNOX device. This collection can contain a maximum of 500 elements.
// returns a []AppListItemable when successful
func (m *AndroidGeneralDeviceConfiguration) GetAppsHideList()([]AppListItemable) {
    val, err := m.GetBackingStore().Get("appsHideList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppListItemable)
    }
    return nil
}
// GetAppsInstallAllowList gets the appsInstallAllowList property value. List of apps which can be installed on the KNOX device. This collection can contain a maximum of 500 elements.
// returns a []AppListItemable when successful
func (m *AndroidGeneralDeviceConfiguration) GetAppsInstallAllowList()([]AppListItemable) {
    val, err := m.GetBackingStore().Get("appsInstallAllowList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppListItemable)
    }
    return nil
}
// GetAppsLaunchBlockList gets the appsLaunchBlockList property value. List of apps which are blocked from being launched on the KNOX device. This collection can contain a maximum of 500 elements.
// returns a []AppListItemable when successful
func (m *AndroidGeneralDeviceConfiguration) GetAppsLaunchBlockList()([]AppListItemable) {
    val, err := m.GetBackingStore().Get("appsLaunchBlockList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppListItemable)
    }
    return nil
}
// GetBluetoothBlocked gets the bluetoothBlocked property value. Indicates whether or not to block Bluetooth.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetBluetoothBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("bluetoothBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCameraBlocked gets the cameraBlocked property value. Indicates whether or not to block the use of the camera.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetCameraBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("cameraBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockDataRoaming gets the cellularBlockDataRoaming property value. Indicates whether or not to block data roaming.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetCellularBlockDataRoaming()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockDataRoaming")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockMessaging gets the cellularBlockMessaging property value. Indicates whether or not to block SMS/MMS messaging.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetCellularBlockMessaging()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockMessaging")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockVoiceRoaming gets the cellularBlockVoiceRoaming property value. Indicates whether or not to block voice roaming.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetCellularBlockVoiceRoaming()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockVoiceRoaming")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockWiFiTethering gets the cellularBlockWiFiTethering property value. Indicates whether or not to block syncing Wi-Fi tethering.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetCellularBlockWiFiTethering()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockWiFiTethering")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCompliantAppListType gets the compliantAppListType property value. Possible values of the compliance app list.
// returns a *AppListType when successful
func (m *AndroidGeneralDeviceConfiguration) GetCompliantAppListType()(*AppListType) {
    val, err := m.GetBackingStore().Get("compliantAppListType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AppListType)
    }
    return nil
}
// GetCompliantAppsList gets the compliantAppsList property value. List of apps in the compliance (either allow list or block list, controlled by CompliantAppListType). This collection can contain a maximum of 10000 elements.
// returns a []AppListItemable when successful
func (m *AndroidGeneralDeviceConfiguration) GetCompliantAppsList()([]AppListItemable) {
    val, err := m.GetBackingStore().Get("compliantAppsList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppListItemable)
    }
    return nil
}
// GetDeviceSharingAllowed gets the deviceSharingAllowed property value. Indicates whether or not to allow device sharing mode.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetDeviceSharingAllowed()(*bool) {
    val, err := m.GetBackingStore().Get("deviceSharingAllowed")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDiagnosticDataBlockSubmission gets the diagnosticDataBlockSubmission property value. Indicates whether or not to block diagnostic data submission.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetDiagnosticDataBlockSubmission()(*bool) {
    val, err := m.GetBackingStore().Get("diagnosticDataBlockSubmission")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFactoryResetBlocked gets the factoryResetBlocked property value. Indicates whether or not to block user performing a factory reset.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetFactoryResetBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("factoryResetBlocked")
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
func (m *AndroidGeneralDeviceConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["appsBlockClipboardSharing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppsBlockClipboardSharing(val)
        }
        return nil
    }
    res["appsBlockCopyPaste"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppsBlockCopyPaste(val)
        }
        return nil
    }
    res["appsBlockYouTube"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppsBlockYouTube(val)
        }
        return nil
    }
    res["appsHideList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppListItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppListItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppListItemable)
                }
            }
            m.SetAppsHideList(res)
        }
        return nil
    }
    res["appsInstallAllowList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppListItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppListItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppListItemable)
                }
            }
            m.SetAppsInstallAllowList(res)
        }
        return nil
    }
    res["appsLaunchBlockList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppListItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppListItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppListItemable)
                }
            }
            m.SetAppsLaunchBlockList(res)
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
    res["cellularBlockDataRoaming"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockDataRoaming(val)
        }
        return nil
    }
    res["cellularBlockMessaging"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockMessaging(val)
        }
        return nil
    }
    res["cellularBlockVoiceRoaming"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockVoiceRoaming(val)
        }
        return nil
    }
    res["cellularBlockWiFiTethering"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockWiFiTethering(val)
        }
        return nil
    }
    res["compliantAppListType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAppListType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompliantAppListType(val.(*AppListType))
        }
        return nil
    }
    res["compliantAppsList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppListItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppListItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppListItemable)
                }
            }
            m.SetCompliantAppsList(res)
        }
        return nil
    }
    res["deviceSharingAllowed"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceSharingAllowed(val)
        }
        return nil
    }
    res["diagnosticDataBlockSubmission"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDiagnosticDataBlockSubmission(val)
        }
        return nil
    }
    res["factoryResetBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFactoryResetBlocked(val)
        }
        return nil
    }
    res["googleAccountBlockAutoSync"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGoogleAccountBlockAutoSync(val)
        }
        return nil
    }
    res["googlePlayStoreBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGooglePlayStoreBlocked(val)
        }
        return nil
    }
    res["kioskModeApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAppListItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AppListItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AppListItemable)
                }
            }
            m.SetKioskModeApps(res)
        }
        return nil
    }
    res["kioskModeBlockSleepButton"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeBlockSleepButton(val)
        }
        return nil
    }
    res["kioskModeBlockVolumeButtons"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeBlockVolumeButtons(val)
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
    res["passwordBlockFingerprintUnlock"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordBlockFingerprintUnlock(val)
        }
        return nil
    }
    res["passwordBlockTrustAgents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordBlockTrustAgents(val)
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
        val, err := n.GetEnumValue(ParseAndroidRequiredPasswordType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordRequiredType(val.(*AndroidRequiredPasswordType))
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
    res["powerOffBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPowerOffBlocked(val)
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
    res["securityRequireVerifyApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSecurityRequireVerifyApps(val)
        }
        return nil
    }
    res["storageBlockGoogleBackup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStorageBlockGoogleBackup(val)
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
    res["storageRequireDeviceEncryption"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStorageRequireDeviceEncryption(val)
        }
        return nil
    }
    res["storageRequireRemovableStorageEncryption"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStorageRequireRemovableStorageEncryption(val)
        }
        return nil
    }
    res["voiceAssistantBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVoiceAssistantBlocked(val)
        }
        return nil
    }
    res["voiceDialingBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVoiceDialingBlocked(val)
        }
        return nil
    }
    res["webBrowserBlockAutofill"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebBrowserBlockAutofill(val)
        }
        return nil
    }
    res["webBrowserBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebBrowserBlocked(val)
        }
        return nil
    }
    res["webBrowserBlockJavaScript"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebBrowserBlockJavaScript(val)
        }
        return nil
    }
    res["webBrowserBlockPopups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebBrowserBlockPopups(val)
        }
        return nil
    }
    res["webBrowserCookieSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWebBrowserCookieSettings)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebBrowserCookieSettings(val.(*WebBrowserCookieSettings))
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
    return res
}
// GetGoogleAccountBlockAutoSync gets the googleAccountBlockAutoSync property value. Indicates whether or not to block Google account auto sync.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetGoogleAccountBlockAutoSync()(*bool) {
    val, err := m.GetBackingStore().Get("googleAccountBlockAutoSync")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetGooglePlayStoreBlocked gets the googlePlayStoreBlocked property value. Indicates whether or not to block the Google Play store.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetGooglePlayStoreBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("googlePlayStoreBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeApps gets the kioskModeApps property value. A list of apps that will be allowed to run when the device is in Kiosk Mode. This collection can contain a maximum of 500 elements.
// returns a []AppListItemable when successful
func (m *AndroidGeneralDeviceConfiguration) GetKioskModeApps()([]AppListItemable) {
    val, err := m.GetBackingStore().Get("kioskModeApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppListItemable)
    }
    return nil
}
// GetKioskModeBlockSleepButton gets the kioskModeBlockSleepButton property value. Indicates whether or not to block the screen sleep button while in Kiosk Mode.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetKioskModeBlockSleepButton()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeBlockSleepButton")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeBlockVolumeButtons gets the kioskModeBlockVolumeButtons property value. Indicates whether or not to block the volume buttons while in Kiosk Mode.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetKioskModeBlockVolumeButtons()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeBlockVolumeButtons")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLocationServicesBlocked gets the locationServicesBlocked property value. Indicates whether or not to block location services.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetLocationServicesBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("locationServicesBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNfcBlocked gets the nfcBlocked property value. Indicates whether or not to block Near-Field Communication.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetNfcBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("nfcBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordBlockFingerprintUnlock gets the passwordBlockFingerprintUnlock property value. Indicates whether or not to block fingerprint unlock.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetPasswordBlockFingerprintUnlock()(*bool) {
    val, err := m.GetBackingStore().Get("passwordBlockFingerprintUnlock")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordBlockTrustAgents gets the passwordBlockTrustAgents property value. Indicates whether or not to block Smart Lock and other trust agents.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetPasswordBlockTrustAgents()(*bool) {
    val, err := m.GetBackingStore().Get("passwordBlockTrustAgents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordExpirationDays gets the passwordExpirationDays property value. Number of days before the password expires. Valid values 1 to 365
// returns a *int32 when successful
func (m *AndroidGeneralDeviceConfiguration) GetPasswordExpirationDays()(*int32) {
    val, err := m.GetBackingStore().Get("passwordExpirationDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinimumLength gets the passwordMinimumLength property value. Minimum length of passwords. Valid values 4 to 16
// returns a *int32 when successful
func (m *AndroidGeneralDeviceConfiguration) GetPasswordMinimumLength()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinimumLength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinutesOfInactivityBeforeScreenTimeout gets the passwordMinutesOfInactivityBeforeScreenTimeout property value. Minutes of inactivity before the screen times out.
// returns a *int32 when successful
func (m *AndroidGeneralDeviceConfiguration) GetPasswordMinutesOfInactivityBeforeScreenTimeout()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinutesOfInactivityBeforeScreenTimeout")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordPreviousPasswordBlockCount gets the passwordPreviousPasswordBlockCount property value. Number of previous passwords to block. Valid values 0 to 24
// returns a *int32 when successful
func (m *AndroidGeneralDeviceConfiguration) GetPasswordPreviousPasswordBlockCount()(*int32) {
    val, err := m.GetBackingStore().Get("passwordPreviousPasswordBlockCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordRequired gets the passwordRequired property value. Indicates whether or not to require a password.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetPasswordRequired()(*bool) {
    val, err := m.GetBackingStore().Get("passwordRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordRequiredType gets the passwordRequiredType property value. Android required password type.
// returns a *AndroidRequiredPasswordType when successful
func (m *AndroidGeneralDeviceConfiguration) GetPasswordRequiredType()(*AndroidRequiredPasswordType) {
    val, err := m.GetBackingStore().Get("passwordRequiredType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AndroidRequiredPasswordType)
    }
    return nil
}
// GetPasswordSignInFailureCountBeforeFactoryReset gets the passwordSignInFailureCountBeforeFactoryReset property value. Number of sign in failures allowed before factory reset. Valid values 1 to 16
// returns a *int32 when successful
func (m *AndroidGeneralDeviceConfiguration) GetPasswordSignInFailureCountBeforeFactoryReset()(*int32) {
    val, err := m.GetBackingStore().Get("passwordSignInFailureCountBeforeFactoryReset")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPowerOffBlocked gets the powerOffBlocked property value. Indicates whether or not to block powering off the device.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetPowerOffBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("powerOffBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetScreenCaptureBlocked gets the screenCaptureBlocked property value. Indicates whether or not to block screenshots.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetScreenCaptureBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("screenCaptureBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSecurityRequireVerifyApps gets the securityRequireVerifyApps property value. Require the Android Verify apps feature is turned on.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetSecurityRequireVerifyApps()(*bool) {
    val, err := m.GetBackingStore().Get("securityRequireVerifyApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStorageBlockGoogleBackup gets the storageBlockGoogleBackup property value. Indicates whether or not to block Google Backup.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetStorageBlockGoogleBackup()(*bool) {
    val, err := m.GetBackingStore().Get("storageBlockGoogleBackup")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStorageBlockRemovableStorage gets the storageBlockRemovableStorage property value. Indicates whether or not to block removable storage usage.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetStorageBlockRemovableStorage()(*bool) {
    val, err := m.GetBackingStore().Get("storageBlockRemovableStorage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStorageRequireDeviceEncryption gets the storageRequireDeviceEncryption property value. Indicates whether or not to require device encryption.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetStorageRequireDeviceEncryption()(*bool) {
    val, err := m.GetBackingStore().Get("storageRequireDeviceEncryption")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStorageRequireRemovableStorageEncryption gets the storageRequireRemovableStorageEncryption property value. Indicates whether or not to require removable storage encryption.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetStorageRequireRemovableStorageEncryption()(*bool) {
    val, err := m.GetBackingStore().Get("storageRequireRemovableStorageEncryption")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetVoiceAssistantBlocked gets the voiceAssistantBlocked property value. Indicates whether or not to block the use of the Voice Assistant.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetVoiceAssistantBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("voiceAssistantBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetVoiceDialingBlocked gets the voiceDialingBlocked property value. Indicates whether or not to block voice dialing.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetVoiceDialingBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("voiceDialingBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWebBrowserBlockAutofill gets the webBrowserBlockAutofill property value. Indicates whether or not to block the web browser's auto fill feature.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetWebBrowserBlockAutofill()(*bool) {
    val, err := m.GetBackingStore().Get("webBrowserBlockAutofill")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWebBrowserBlocked gets the webBrowserBlocked property value. Indicates whether or not to block the web browser.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetWebBrowserBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("webBrowserBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWebBrowserBlockJavaScript gets the webBrowserBlockJavaScript property value. Indicates whether or not to block JavaScript within the web browser.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetWebBrowserBlockJavaScript()(*bool) {
    val, err := m.GetBackingStore().Get("webBrowserBlockJavaScript")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWebBrowserBlockPopups gets the webBrowserBlockPopups property value. Indicates whether or not to block popups within the web browser.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetWebBrowserBlockPopups()(*bool) {
    val, err := m.GetBackingStore().Get("webBrowserBlockPopups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWebBrowserCookieSettings gets the webBrowserCookieSettings property value. Web Browser Cookie Settings.
// returns a *WebBrowserCookieSettings when successful
func (m *AndroidGeneralDeviceConfiguration) GetWebBrowserCookieSettings()(*WebBrowserCookieSettings) {
    val, err := m.GetBackingStore().Get("webBrowserCookieSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WebBrowserCookieSettings)
    }
    return nil
}
// GetWiFiBlocked gets the wiFiBlocked property value. Indicates whether or not to block syncing Wi-Fi.
// returns a *bool when successful
func (m *AndroidGeneralDeviceConfiguration) GetWiFiBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("wiFiBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AndroidGeneralDeviceConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("appsBlockClipboardSharing", m.GetAppsBlockClipboardSharing())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appsBlockCopyPaste", m.GetAppsBlockCopyPaste())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appsBlockYouTube", m.GetAppsBlockYouTube())
        if err != nil {
            return err
        }
    }
    if m.GetAppsHideList() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppsHideList()))
        for i, v := range m.GetAppsHideList() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appsHideList", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAppsInstallAllowList() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppsInstallAllowList()))
        for i, v := range m.GetAppsInstallAllowList() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appsInstallAllowList", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAppsLaunchBlockList() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppsLaunchBlockList()))
        for i, v := range m.GetAppsLaunchBlockList() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appsLaunchBlockList", cast)
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
        err = writer.WriteBoolValue("cameraBlocked", m.GetCameraBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cellularBlockDataRoaming", m.GetCellularBlockDataRoaming())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cellularBlockMessaging", m.GetCellularBlockMessaging())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cellularBlockVoiceRoaming", m.GetCellularBlockVoiceRoaming())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cellularBlockWiFiTethering", m.GetCellularBlockWiFiTethering())
        if err != nil {
            return err
        }
    }
    if m.GetCompliantAppListType() != nil {
        cast := (*m.GetCompliantAppListType()).String()
        err = writer.WriteStringValue("compliantAppListType", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetCompliantAppsList() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCompliantAppsList()))
        for i, v := range m.GetCompliantAppsList() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("compliantAppsList", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("deviceSharingAllowed", m.GetDeviceSharingAllowed())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("diagnosticDataBlockSubmission", m.GetDiagnosticDataBlockSubmission())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("factoryResetBlocked", m.GetFactoryResetBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("googleAccountBlockAutoSync", m.GetGoogleAccountBlockAutoSync())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("googlePlayStoreBlocked", m.GetGooglePlayStoreBlocked())
        if err != nil {
            return err
        }
    }
    if m.GetKioskModeApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetKioskModeApps()))
        for i, v := range m.GetKioskModeApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("kioskModeApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeBlockSleepButton", m.GetKioskModeBlockSleepButton())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeBlockVolumeButtons", m.GetKioskModeBlockVolumeButtons())
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
        err = writer.WriteBoolValue("nfcBlocked", m.GetNfcBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordBlockFingerprintUnlock", m.GetPasswordBlockFingerprintUnlock())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordBlockTrustAgents", m.GetPasswordBlockTrustAgents())
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
        err = writer.WriteInt32Value("passwordSignInFailureCountBeforeFactoryReset", m.GetPasswordSignInFailureCountBeforeFactoryReset())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("powerOffBlocked", m.GetPowerOffBlocked())
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
        err = writer.WriteBoolValue("securityRequireVerifyApps", m.GetSecurityRequireVerifyApps())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("storageBlockGoogleBackup", m.GetStorageBlockGoogleBackup())
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
        err = writer.WriteBoolValue("storageRequireDeviceEncryption", m.GetStorageRequireDeviceEncryption())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("storageRequireRemovableStorageEncryption", m.GetStorageRequireRemovableStorageEncryption())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("voiceAssistantBlocked", m.GetVoiceAssistantBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("voiceDialingBlocked", m.GetVoiceDialingBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("webBrowserBlockAutofill", m.GetWebBrowserBlockAutofill())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("webBrowserBlocked", m.GetWebBrowserBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("webBrowserBlockJavaScript", m.GetWebBrowserBlockJavaScript())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("webBrowserBlockPopups", m.GetWebBrowserBlockPopups())
        if err != nil {
            return err
        }
    }
    if m.GetWebBrowserCookieSettings() != nil {
        cast := (*m.GetWebBrowserCookieSettings()).String()
        err = writer.WriteStringValue("webBrowserCookieSettings", &cast)
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
    return nil
}
// SetAppsBlockClipboardSharing sets the appsBlockClipboardSharing property value. Indicates whether or not to block clipboard sharing to copy and paste between applications.
func (m *AndroidGeneralDeviceConfiguration) SetAppsBlockClipboardSharing(value *bool)() {
    err := m.GetBackingStore().Set("appsBlockClipboardSharing", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsBlockCopyPaste sets the appsBlockCopyPaste property value. Indicates whether or not to block copy and paste within applications.
func (m *AndroidGeneralDeviceConfiguration) SetAppsBlockCopyPaste(value *bool)() {
    err := m.GetBackingStore().Set("appsBlockCopyPaste", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsBlockYouTube sets the appsBlockYouTube property value. Indicates whether or not to block the YouTube app.
func (m *AndroidGeneralDeviceConfiguration) SetAppsBlockYouTube(value *bool)() {
    err := m.GetBackingStore().Set("appsBlockYouTube", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsHideList sets the appsHideList property value. List of apps to be hidden on the KNOX device. This collection can contain a maximum of 500 elements.
func (m *AndroidGeneralDeviceConfiguration) SetAppsHideList(value []AppListItemable)() {
    err := m.GetBackingStore().Set("appsHideList", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsInstallAllowList sets the appsInstallAllowList property value. List of apps which can be installed on the KNOX device. This collection can contain a maximum of 500 elements.
func (m *AndroidGeneralDeviceConfiguration) SetAppsInstallAllowList(value []AppListItemable)() {
    err := m.GetBackingStore().Set("appsInstallAllowList", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsLaunchBlockList sets the appsLaunchBlockList property value. List of apps which are blocked from being launched on the KNOX device. This collection can contain a maximum of 500 elements.
func (m *AndroidGeneralDeviceConfiguration) SetAppsLaunchBlockList(value []AppListItemable)() {
    err := m.GetBackingStore().Set("appsLaunchBlockList", value)
    if err != nil {
        panic(err)
    }
}
// SetBluetoothBlocked sets the bluetoothBlocked property value. Indicates whether or not to block Bluetooth.
func (m *AndroidGeneralDeviceConfiguration) SetBluetoothBlocked(value *bool)() {
    err := m.GetBackingStore().Set("bluetoothBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetCameraBlocked sets the cameraBlocked property value. Indicates whether or not to block the use of the camera.
func (m *AndroidGeneralDeviceConfiguration) SetCameraBlocked(value *bool)() {
    err := m.GetBackingStore().Set("cameraBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockDataRoaming sets the cellularBlockDataRoaming property value. Indicates whether or not to block data roaming.
func (m *AndroidGeneralDeviceConfiguration) SetCellularBlockDataRoaming(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockDataRoaming", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockMessaging sets the cellularBlockMessaging property value. Indicates whether or not to block SMS/MMS messaging.
func (m *AndroidGeneralDeviceConfiguration) SetCellularBlockMessaging(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockMessaging", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockVoiceRoaming sets the cellularBlockVoiceRoaming property value. Indicates whether or not to block voice roaming.
func (m *AndroidGeneralDeviceConfiguration) SetCellularBlockVoiceRoaming(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockVoiceRoaming", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockWiFiTethering sets the cellularBlockWiFiTethering property value. Indicates whether or not to block syncing Wi-Fi tethering.
func (m *AndroidGeneralDeviceConfiguration) SetCellularBlockWiFiTethering(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockWiFiTethering", value)
    if err != nil {
        panic(err)
    }
}
// SetCompliantAppListType sets the compliantAppListType property value. Possible values of the compliance app list.
func (m *AndroidGeneralDeviceConfiguration) SetCompliantAppListType(value *AppListType)() {
    err := m.GetBackingStore().Set("compliantAppListType", value)
    if err != nil {
        panic(err)
    }
}
// SetCompliantAppsList sets the compliantAppsList property value. List of apps in the compliance (either allow list or block list, controlled by CompliantAppListType). This collection can contain a maximum of 10000 elements.
func (m *AndroidGeneralDeviceConfiguration) SetCompliantAppsList(value []AppListItemable)() {
    err := m.GetBackingStore().Set("compliantAppsList", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceSharingAllowed sets the deviceSharingAllowed property value. Indicates whether or not to allow device sharing mode.
func (m *AndroidGeneralDeviceConfiguration) SetDeviceSharingAllowed(value *bool)() {
    err := m.GetBackingStore().Set("deviceSharingAllowed", value)
    if err != nil {
        panic(err)
    }
}
// SetDiagnosticDataBlockSubmission sets the diagnosticDataBlockSubmission property value. Indicates whether or not to block diagnostic data submission.
func (m *AndroidGeneralDeviceConfiguration) SetDiagnosticDataBlockSubmission(value *bool)() {
    err := m.GetBackingStore().Set("diagnosticDataBlockSubmission", value)
    if err != nil {
        panic(err)
    }
}
// SetFactoryResetBlocked sets the factoryResetBlocked property value. Indicates whether or not to block user performing a factory reset.
func (m *AndroidGeneralDeviceConfiguration) SetFactoryResetBlocked(value *bool)() {
    err := m.GetBackingStore().Set("factoryResetBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetGoogleAccountBlockAutoSync sets the googleAccountBlockAutoSync property value. Indicates whether or not to block Google account auto sync.
func (m *AndroidGeneralDeviceConfiguration) SetGoogleAccountBlockAutoSync(value *bool)() {
    err := m.GetBackingStore().Set("googleAccountBlockAutoSync", value)
    if err != nil {
        panic(err)
    }
}
// SetGooglePlayStoreBlocked sets the googlePlayStoreBlocked property value. Indicates whether or not to block the Google Play store.
func (m *AndroidGeneralDeviceConfiguration) SetGooglePlayStoreBlocked(value *bool)() {
    err := m.GetBackingStore().Set("googlePlayStoreBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeApps sets the kioskModeApps property value. A list of apps that will be allowed to run when the device is in Kiosk Mode. This collection can contain a maximum of 500 elements.
func (m *AndroidGeneralDeviceConfiguration) SetKioskModeApps(value []AppListItemable)() {
    err := m.GetBackingStore().Set("kioskModeApps", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeBlockSleepButton sets the kioskModeBlockSleepButton property value. Indicates whether or not to block the screen sleep button while in Kiosk Mode.
func (m *AndroidGeneralDeviceConfiguration) SetKioskModeBlockSleepButton(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeBlockSleepButton", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeBlockVolumeButtons sets the kioskModeBlockVolumeButtons property value. Indicates whether or not to block the volume buttons while in Kiosk Mode.
func (m *AndroidGeneralDeviceConfiguration) SetKioskModeBlockVolumeButtons(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeBlockVolumeButtons", value)
    if err != nil {
        panic(err)
    }
}
// SetLocationServicesBlocked sets the locationServicesBlocked property value. Indicates whether or not to block location services.
func (m *AndroidGeneralDeviceConfiguration) SetLocationServicesBlocked(value *bool)() {
    err := m.GetBackingStore().Set("locationServicesBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetNfcBlocked sets the nfcBlocked property value. Indicates whether or not to block Near-Field Communication.
func (m *AndroidGeneralDeviceConfiguration) SetNfcBlocked(value *bool)() {
    err := m.GetBackingStore().Set("nfcBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordBlockFingerprintUnlock sets the passwordBlockFingerprintUnlock property value. Indicates whether or not to block fingerprint unlock.
func (m *AndroidGeneralDeviceConfiguration) SetPasswordBlockFingerprintUnlock(value *bool)() {
    err := m.GetBackingStore().Set("passwordBlockFingerprintUnlock", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordBlockTrustAgents sets the passwordBlockTrustAgents property value. Indicates whether or not to block Smart Lock and other trust agents.
func (m *AndroidGeneralDeviceConfiguration) SetPasswordBlockTrustAgents(value *bool)() {
    err := m.GetBackingStore().Set("passwordBlockTrustAgents", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordExpirationDays sets the passwordExpirationDays property value. Number of days before the password expires. Valid values 1 to 365
func (m *AndroidGeneralDeviceConfiguration) SetPasswordExpirationDays(value *int32)() {
    err := m.GetBackingStore().Set("passwordExpirationDays", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinimumLength sets the passwordMinimumLength property value. Minimum length of passwords. Valid values 4 to 16
func (m *AndroidGeneralDeviceConfiguration) SetPasswordMinimumLength(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinimumLength", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinutesOfInactivityBeforeScreenTimeout sets the passwordMinutesOfInactivityBeforeScreenTimeout property value. Minutes of inactivity before the screen times out.
func (m *AndroidGeneralDeviceConfiguration) SetPasswordMinutesOfInactivityBeforeScreenTimeout(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinutesOfInactivityBeforeScreenTimeout", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordPreviousPasswordBlockCount sets the passwordPreviousPasswordBlockCount property value. Number of previous passwords to block. Valid values 0 to 24
func (m *AndroidGeneralDeviceConfiguration) SetPasswordPreviousPasswordBlockCount(value *int32)() {
    err := m.GetBackingStore().Set("passwordPreviousPasswordBlockCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordRequired sets the passwordRequired property value. Indicates whether or not to require a password.
func (m *AndroidGeneralDeviceConfiguration) SetPasswordRequired(value *bool)() {
    err := m.GetBackingStore().Set("passwordRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordRequiredType sets the passwordRequiredType property value. Android required password type.
func (m *AndroidGeneralDeviceConfiguration) SetPasswordRequiredType(value *AndroidRequiredPasswordType)() {
    err := m.GetBackingStore().Set("passwordRequiredType", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordSignInFailureCountBeforeFactoryReset sets the passwordSignInFailureCountBeforeFactoryReset property value. Number of sign in failures allowed before factory reset. Valid values 1 to 16
func (m *AndroidGeneralDeviceConfiguration) SetPasswordSignInFailureCountBeforeFactoryReset(value *int32)() {
    err := m.GetBackingStore().Set("passwordSignInFailureCountBeforeFactoryReset", value)
    if err != nil {
        panic(err)
    }
}
// SetPowerOffBlocked sets the powerOffBlocked property value. Indicates whether or not to block powering off the device.
func (m *AndroidGeneralDeviceConfiguration) SetPowerOffBlocked(value *bool)() {
    err := m.GetBackingStore().Set("powerOffBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetScreenCaptureBlocked sets the screenCaptureBlocked property value. Indicates whether or not to block screenshots.
func (m *AndroidGeneralDeviceConfiguration) SetScreenCaptureBlocked(value *bool)() {
    err := m.GetBackingStore().Set("screenCaptureBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetSecurityRequireVerifyApps sets the securityRequireVerifyApps property value. Require the Android Verify apps feature is turned on.
func (m *AndroidGeneralDeviceConfiguration) SetSecurityRequireVerifyApps(value *bool)() {
    err := m.GetBackingStore().Set("securityRequireVerifyApps", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageBlockGoogleBackup sets the storageBlockGoogleBackup property value. Indicates whether or not to block Google Backup.
func (m *AndroidGeneralDeviceConfiguration) SetStorageBlockGoogleBackup(value *bool)() {
    err := m.GetBackingStore().Set("storageBlockGoogleBackup", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageBlockRemovableStorage sets the storageBlockRemovableStorage property value. Indicates whether or not to block removable storage usage.
func (m *AndroidGeneralDeviceConfiguration) SetStorageBlockRemovableStorage(value *bool)() {
    err := m.GetBackingStore().Set("storageBlockRemovableStorage", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageRequireDeviceEncryption sets the storageRequireDeviceEncryption property value. Indicates whether or not to require device encryption.
func (m *AndroidGeneralDeviceConfiguration) SetStorageRequireDeviceEncryption(value *bool)() {
    err := m.GetBackingStore().Set("storageRequireDeviceEncryption", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageRequireRemovableStorageEncryption sets the storageRequireRemovableStorageEncryption property value. Indicates whether or not to require removable storage encryption.
func (m *AndroidGeneralDeviceConfiguration) SetStorageRequireRemovableStorageEncryption(value *bool)() {
    err := m.GetBackingStore().Set("storageRequireRemovableStorageEncryption", value)
    if err != nil {
        panic(err)
    }
}
// SetVoiceAssistantBlocked sets the voiceAssistantBlocked property value. Indicates whether or not to block the use of the Voice Assistant.
func (m *AndroidGeneralDeviceConfiguration) SetVoiceAssistantBlocked(value *bool)() {
    err := m.GetBackingStore().Set("voiceAssistantBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetVoiceDialingBlocked sets the voiceDialingBlocked property value. Indicates whether or not to block voice dialing.
func (m *AndroidGeneralDeviceConfiguration) SetVoiceDialingBlocked(value *bool)() {
    err := m.GetBackingStore().Set("voiceDialingBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetWebBrowserBlockAutofill sets the webBrowserBlockAutofill property value. Indicates whether or not to block the web browser's auto fill feature.
func (m *AndroidGeneralDeviceConfiguration) SetWebBrowserBlockAutofill(value *bool)() {
    err := m.GetBackingStore().Set("webBrowserBlockAutofill", value)
    if err != nil {
        panic(err)
    }
}
// SetWebBrowserBlocked sets the webBrowserBlocked property value. Indicates whether or not to block the web browser.
func (m *AndroidGeneralDeviceConfiguration) SetWebBrowserBlocked(value *bool)() {
    err := m.GetBackingStore().Set("webBrowserBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetWebBrowserBlockJavaScript sets the webBrowserBlockJavaScript property value. Indicates whether or not to block JavaScript within the web browser.
func (m *AndroidGeneralDeviceConfiguration) SetWebBrowserBlockJavaScript(value *bool)() {
    err := m.GetBackingStore().Set("webBrowserBlockJavaScript", value)
    if err != nil {
        panic(err)
    }
}
// SetWebBrowserBlockPopups sets the webBrowserBlockPopups property value. Indicates whether or not to block popups within the web browser.
func (m *AndroidGeneralDeviceConfiguration) SetWebBrowserBlockPopups(value *bool)() {
    err := m.GetBackingStore().Set("webBrowserBlockPopups", value)
    if err != nil {
        panic(err)
    }
}
// SetWebBrowserCookieSettings sets the webBrowserCookieSettings property value. Web Browser Cookie Settings.
func (m *AndroidGeneralDeviceConfiguration) SetWebBrowserCookieSettings(value *WebBrowserCookieSettings)() {
    err := m.GetBackingStore().Set("webBrowserCookieSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetWiFiBlocked sets the wiFiBlocked property value. Indicates whether or not to block syncing Wi-Fi.
func (m *AndroidGeneralDeviceConfiguration) SetWiFiBlocked(value *bool)() {
    err := m.GetBackingStore().Set("wiFiBlocked", value)
    if err != nil {
        panic(err)
    }
}
type AndroidGeneralDeviceConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppsBlockClipboardSharing()(*bool)
    GetAppsBlockCopyPaste()(*bool)
    GetAppsBlockYouTube()(*bool)
    GetAppsHideList()([]AppListItemable)
    GetAppsInstallAllowList()([]AppListItemable)
    GetAppsLaunchBlockList()([]AppListItemable)
    GetBluetoothBlocked()(*bool)
    GetCameraBlocked()(*bool)
    GetCellularBlockDataRoaming()(*bool)
    GetCellularBlockMessaging()(*bool)
    GetCellularBlockVoiceRoaming()(*bool)
    GetCellularBlockWiFiTethering()(*bool)
    GetCompliantAppListType()(*AppListType)
    GetCompliantAppsList()([]AppListItemable)
    GetDeviceSharingAllowed()(*bool)
    GetDiagnosticDataBlockSubmission()(*bool)
    GetFactoryResetBlocked()(*bool)
    GetGoogleAccountBlockAutoSync()(*bool)
    GetGooglePlayStoreBlocked()(*bool)
    GetKioskModeApps()([]AppListItemable)
    GetKioskModeBlockSleepButton()(*bool)
    GetKioskModeBlockVolumeButtons()(*bool)
    GetLocationServicesBlocked()(*bool)
    GetNfcBlocked()(*bool)
    GetPasswordBlockFingerprintUnlock()(*bool)
    GetPasswordBlockTrustAgents()(*bool)
    GetPasswordExpirationDays()(*int32)
    GetPasswordMinimumLength()(*int32)
    GetPasswordMinutesOfInactivityBeforeScreenTimeout()(*int32)
    GetPasswordPreviousPasswordBlockCount()(*int32)
    GetPasswordRequired()(*bool)
    GetPasswordRequiredType()(*AndroidRequiredPasswordType)
    GetPasswordSignInFailureCountBeforeFactoryReset()(*int32)
    GetPowerOffBlocked()(*bool)
    GetScreenCaptureBlocked()(*bool)
    GetSecurityRequireVerifyApps()(*bool)
    GetStorageBlockGoogleBackup()(*bool)
    GetStorageBlockRemovableStorage()(*bool)
    GetStorageRequireDeviceEncryption()(*bool)
    GetStorageRequireRemovableStorageEncryption()(*bool)
    GetVoiceAssistantBlocked()(*bool)
    GetVoiceDialingBlocked()(*bool)
    GetWebBrowserBlockAutofill()(*bool)
    GetWebBrowserBlocked()(*bool)
    GetWebBrowserBlockJavaScript()(*bool)
    GetWebBrowserBlockPopups()(*bool)
    GetWebBrowserCookieSettings()(*WebBrowserCookieSettings)
    GetWiFiBlocked()(*bool)
    SetAppsBlockClipboardSharing(value *bool)()
    SetAppsBlockCopyPaste(value *bool)()
    SetAppsBlockYouTube(value *bool)()
    SetAppsHideList(value []AppListItemable)()
    SetAppsInstallAllowList(value []AppListItemable)()
    SetAppsLaunchBlockList(value []AppListItemable)()
    SetBluetoothBlocked(value *bool)()
    SetCameraBlocked(value *bool)()
    SetCellularBlockDataRoaming(value *bool)()
    SetCellularBlockMessaging(value *bool)()
    SetCellularBlockVoiceRoaming(value *bool)()
    SetCellularBlockWiFiTethering(value *bool)()
    SetCompliantAppListType(value *AppListType)()
    SetCompliantAppsList(value []AppListItemable)()
    SetDeviceSharingAllowed(value *bool)()
    SetDiagnosticDataBlockSubmission(value *bool)()
    SetFactoryResetBlocked(value *bool)()
    SetGoogleAccountBlockAutoSync(value *bool)()
    SetGooglePlayStoreBlocked(value *bool)()
    SetKioskModeApps(value []AppListItemable)()
    SetKioskModeBlockSleepButton(value *bool)()
    SetKioskModeBlockVolumeButtons(value *bool)()
    SetLocationServicesBlocked(value *bool)()
    SetNfcBlocked(value *bool)()
    SetPasswordBlockFingerprintUnlock(value *bool)()
    SetPasswordBlockTrustAgents(value *bool)()
    SetPasswordExpirationDays(value *int32)()
    SetPasswordMinimumLength(value *int32)()
    SetPasswordMinutesOfInactivityBeforeScreenTimeout(value *int32)()
    SetPasswordPreviousPasswordBlockCount(value *int32)()
    SetPasswordRequired(value *bool)()
    SetPasswordRequiredType(value *AndroidRequiredPasswordType)()
    SetPasswordSignInFailureCountBeforeFactoryReset(value *int32)()
    SetPowerOffBlocked(value *bool)()
    SetScreenCaptureBlocked(value *bool)()
    SetSecurityRequireVerifyApps(value *bool)()
    SetStorageBlockGoogleBackup(value *bool)()
    SetStorageBlockRemovableStorage(value *bool)()
    SetStorageRequireDeviceEncryption(value *bool)()
    SetStorageRequireRemovableStorageEncryption(value *bool)()
    SetVoiceAssistantBlocked(value *bool)()
    SetVoiceDialingBlocked(value *bool)()
    SetWebBrowserBlockAutofill(value *bool)()
    SetWebBrowserBlocked(value *bool)()
    SetWebBrowserBlockJavaScript(value *bool)()
    SetWebBrowserBlockPopups(value *bool)()
    SetWebBrowserCookieSettings(value *WebBrowserCookieSettings)()
    SetWiFiBlocked(value *bool)()
}
