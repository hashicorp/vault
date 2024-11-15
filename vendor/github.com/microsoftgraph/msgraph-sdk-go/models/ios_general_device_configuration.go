package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosGeneralDeviceConfiguration this topic provides descriptions of the declared methods, properties and relationships exposed by the iosGeneralDeviceConfiguration resource.
type IosGeneralDeviceConfiguration struct {
    DeviceConfiguration
}
// NewIosGeneralDeviceConfiguration instantiates a new IosGeneralDeviceConfiguration and sets the default values.
func NewIosGeneralDeviceConfiguration()(*IosGeneralDeviceConfiguration) {
    m := &IosGeneralDeviceConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.iosGeneralDeviceConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosGeneralDeviceConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosGeneralDeviceConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosGeneralDeviceConfiguration(), nil
}
// GetAccountBlockModification gets the accountBlockModification property value. Indicates whether or not to allow account modification when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAccountBlockModification()(*bool) {
    val, err := m.GetBackingStore().Get("accountBlockModification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetActivationLockAllowWhenSupervised gets the activationLockAllowWhenSupervised property value. Indicates whether or not to allow activation lock when the device is in the supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetActivationLockAllowWhenSupervised()(*bool) {
    val, err := m.GetBackingStore().Get("activationLockAllowWhenSupervised")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAirDropBlocked gets the airDropBlocked property value. Indicates whether or not to allow AirDrop when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAirDropBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("airDropBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAirDropForceUnmanagedDropTarget gets the airDropForceUnmanagedDropTarget property value. Indicates whether or not to cause AirDrop to be considered an unmanaged drop target (iOS 9.0 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAirDropForceUnmanagedDropTarget()(*bool) {
    val, err := m.GetBackingStore().Get("airDropForceUnmanagedDropTarget")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAirPlayForcePairingPasswordForOutgoingRequests gets the airPlayForcePairingPasswordForOutgoingRequests property value. Indicates whether or not to enforce all devices receiving AirPlay requests from this device to use a pairing password.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAirPlayForcePairingPasswordForOutgoingRequests()(*bool) {
    val, err := m.GetBackingStore().Get("airPlayForcePairingPasswordForOutgoingRequests")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppleNewsBlocked gets the appleNewsBlocked property value. Indicates whether or not to block the user from using News when the device is in supervised mode (iOS 9.0 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAppleNewsBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("appleNewsBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppleWatchBlockPairing gets the appleWatchBlockPairing property value. Indicates whether or not to allow Apple Watch pairing when the device is in supervised mode (iOS 9.0 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAppleWatchBlockPairing()(*bool) {
    val, err := m.GetBackingStore().Get("appleWatchBlockPairing")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppleWatchForceWristDetection gets the appleWatchForceWristDetection property value. Indicates whether or not to force a paired Apple Watch to use Wrist Detection (iOS 8.2 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAppleWatchForceWristDetection()(*bool) {
    val, err := m.GetBackingStore().Get("appleWatchForceWristDetection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppsSingleAppModeList gets the appsSingleAppModeList property value. Gets or sets the list of iOS apps allowed to autonomously enter Single App Mode. Supervised only. iOS 7.0 and later. This collection can contain a maximum of 500 elements.
// returns a []AppListItemable when successful
func (m *IosGeneralDeviceConfiguration) GetAppsSingleAppModeList()([]AppListItemable) {
    val, err := m.GetBackingStore().Get("appsSingleAppModeList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppListItemable)
    }
    return nil
}
// GetAppStoreBlockAutomaticDownloads gets the appStoreBlockAutomaticDownloads property value. Indicates whether or not to block the automatic downloading of apps purchased on other devices when the device is in supervised mode (iOS 9.0 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAppStoreBlockAutomaticDownloads()(*bool) {
    val, err := m.GetBackingStore().Get("appStoreBlockAutomaticDownloads")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppStoreBlocked gets the appStoreBlocked property value. Indicates whether or not to block the user from using the App Store. Requires a supervised device for iOS 13 and later.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAppStoreBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("appStoreBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppStoreBlockInAppPurchases gets the appStoreBlockInAppPurchases property value. Indicates whether or not to block the user from making in app purchases.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAppStoreBlockInAppPurchases()(*bool) {
    val, err := m.GetBackingStore().Get("appStoreBlockInAppPurchases")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppStoreBlockUIAppInstallation gets the appStoreBlockUIAppInstallation property value. Indicates whether or not to block the App Store app, not restricting installation through Host apps. Applies to supervised mode only (iOS 9.0 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAppStoreBlockUIAppInstallation()(*bool) {
    val, err := m.GetBackingStore().Get("appStoreBlockUIAppInstallation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppStoreRequirePassword gets the appStoreRequirePassword property value. Indicates whether or not to require a password when using the app store.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetAppStoreRequirePassword()(*bool) {
    val, err := m.GetBackingStore().Get("appStoreRequirePassword")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppsVisibilityList gets the appsVisibilityList property value. List of apps in the visibility list (either visible/launchable apps list or hidden/unlaunchable apps list, controlled by AppsVisibilityListType) (iOS 9.3 and later). This collection can contain a maximum of 10000 elements.
// returns a []AppListItemable when successful
func (m *IosGeneralDeviceConfiguration) GetAppsVisibilityList()([]AppListItemable) {
    val, err := m.GetBackingStore().Get("appsVisibilityList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppListItemable)
    }
    return nil
}
// GetAppsVisibilityListType gets the appsVisibilityListType property value. Possible values of the compliance app list.
// returns a *AppListType when successful
func (m *IosGeneralDeviceConfiguration) GetAppsVisibilityListType()(*AppListType) {
    val, err := m.GetBackingStore().Get("appsVisibilityListType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AppListType)
    }
    return nil
}
// GetBluetoothBlockModification gets the bluetoothBlockModification property value. Indicates whether or not to allow modification of Bluetooth settings when the device is in supervised mode (iOS 10.0 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetBluetoothBlockModification()(*bool) {
    val, err := m.GetBackingStore().Get("bluetoothBlockModification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCameraBlocked gets the cameraBlocked property value. Indicates whether or not to block the user from accessing the camera of the device. Requires a supervised device for iOS 13 and later.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetCameraBlocked()(*bool) {
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
func (m *IosGeneralDeviceConfiguration) GetCellularBlockDataRoaming()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockDataRoaming")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockGlobalBackgroundFetchWhileRoaming gets the cellularBlockGlobalBackgroundFetchWhileRoaming property value. Indicates whether or not to block global background fetch while roaming.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetCellularBlockGlobalBackgroundFetchWhileRoaming()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockGlobalBackgroundFetchWhileRoaming")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockPerAppDataModification gets the cellularBlockPerAppDataModification property value. Indicates whether or not to allow changes to cellular app data usage settings when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetCellularBlockPerAppDataModification()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockPerAppDataModification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockPersonalHotspot gets the cellularBlockPersonalHotspot property value. Indicates whether or not to block Personal Hotspot.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetCellularBlockPersonalHotspot()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockPersonalHotspot")
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
func (m *IosGeneralDeviceConfiguration) GetCellularBlockVoiceRoaming()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockVoiceRoaming")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCertificatesBlockUntrustedTlsCertificates gets the certificatesBlockUntrustedTlsCertificates property value. Indicates whether or not to block untrusted TLS certificates.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetCertificatesBlockUntrustedTlsCertificates()(*bool) {
    val, err := m.GetBackingStore().Get("certificatesBlockUntrustedTlsCertificates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetClassroomAppBlockRemoteScreenObservation gets the classroomAppBlockRemoteScreenObservation property value. Indicates whether or not to allow remote screen observation by Classroom app when the device is in supervised mode (iOS 9.3 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetClassroomAppBlockRemoteScreenObservation()(*bool) {
    val, err := m.GetBackingStore().Get("classroomAppBlockRemoteScreenObservation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetClassroomAppForceUnpromptedScreenObservation gets the classroomAppForceUnpromptedScreenObservation property value. Indicates whether or not to automatically give permission to the teacher of a managed course on the Classroom app to view a student's screen without prompting when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetClassroomAppForceUnpromptedScreenObservation()(*bool) {
    val, err := m.GetBackingStore().Get("classroomAppForceUnpromptedScreenObservation")
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
func (m *IosGeneralDeviceConfiguration) GetCompliantAppListType()(*AppListType) {
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
func (m *IosGeneralDeviceConfiguration) GetCompliantAppsList()([]AppListItemable) {
    val, err := m.GetBackingStore().Get("compliantAppsList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppListItemable)
    }
    return nil
}
// GetConfigurationProfileBlockChanges gets the configurationProfileBlockChanges property value. Indicates whether or not to block the user from installing configuration profiles and certificates interactively when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetConfigurationProfileBlockChanges()(*bool) {
    val, err := m.GetBackingStore().Get("configurationProfileBlockChanges")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDefinitionLookupBlocked gets the definitionLookupBlocked property value. Indicates whether or not to block definition lookup when the device is in supervised mode (iOS 8.1.3 and later ).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetDefinitionLookupBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("definitionLookupBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDeviceBlockEnableRestrictions gets the deviceBlockEnableRestrictions property value. Indicates whether or not to allow the user to enables restrictions in the device settings when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetDeviceBlockEnableRestrictions()(*bool) {
    val, err := m.GetBackingStore().Get("deviceBlockEnableRestrictions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDeviceBlockEraseContentAndSettings gets the deviceBlockEraseContentAndSettings property value. Indicates whether or not to allow the use of the 'Erase all content and settings' option on the device when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetDeviceBlockEraseContentAndSettings()(*bool) {
    val, err := m.GetBackingStore().Get("deviceBlockEraseContentAndSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDeviceBlockNameModification gets the deviceBlockNameModification property value. Indicates whether or not to allow device name modification when the device is in supervised mode (iOS 9.0 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetDeviceBlockNameModification()(*bool) {
    val, err := m.GetBackingStore().Get("deviceBlockNameModification")
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
func (m *IosGeneralDeviceConfiguration) GetDiagnosticDataBlockSubmission()(*bool) {
    val, err := m.GetBackingStore().Get("diagnosticDataBlockSubmission")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDiagnosticDataBlockSubmissionModification gets the diagnosticDataBlockSubmissionModification property value. Indicates whether or not to allow diagnostics submission settings modification when the device is in supervised mode (iOS 9.3.2 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetDiagnosticDataBlockSubmissionModification()(*bool) {
    val, err := m.GetBackingStore().Get("diagnosticDataBlockSubmissionModification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDocumentsBlockManagedDocumentsInUnmanagedApps gets the documentsBlockManagedDocumentsInUnmanagedApps property value. Indicates whether or not to block the user from viewing managed documents in unmanaged apps.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetDocumentsBlockManagedDocumentsInUnmanagedApps()(*bool) {
    val, err := m.GetBackingStore().Get("documentsBlockManagedDocumentsInUnmanagedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDocumentsBlockUnmanagedDocumentsInManagedApps gets the documentsBlockUnmanagedDocumentsInManagedApps property value. Indicates whether or not to block the user from viewing unmanaged documents in managed apps.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetDocumentsBlockUnmanagedDocumentsInManagedApps()(*bool) {
    val, err := m.GetBackingStore().Get("documentsBlockUnmanagedDocumentsInManagedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEmailInDomainSuffixes gets the emailInDomainSuffixes property value. An email address lacking a suffix that matches any of these strings will be considered out-of-domain.
// returns a []string when successful
func (m *IosGeneralDeviceConfiguration) GetEmailInDomainSuffixes()([]string) {
    val, err := m.GetBackingStore().Get("emailInDomainSuffixes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetEnterpriseAppBlockTrust gets the enterpriseAppBlockTrust property value. Indicates whether or not to block the user from trusting an enterprise app.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetEnterpriseAppBlockTrust()(*bool) {
    val, err := m.GetBackingStore().Get("enterpriseAppBlockTrust")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEnterpriseAppBlockTrustModification gets the enterpriseAppBlockTrustModification property value. [Deprecated] Configuring this setting and setting the value to 'true' has no effect on the device.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetEnterpriseAppBlockTrustModification()(*bool) {
    val, err := m.GetBackingStore().Get("enterpriseAppBlockTrustModification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFaceTimeBlocked gets the faceTimeBlocked property value. Indicates whether or not to block the user from using FaceTime. Requires a supervised device for iOS 13 and later.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetFaceTimeBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("faceTimeBlocked")
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
func (m *IosGeneralDeviceConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["accountBlockModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccountBlockModification(val)
        }
        return nil
    }
    res["activationLockAllowWhenSupervised"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivationLockAllowWhenSupervised(val)
        }
        return nil
    }
    res["airDropBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAirDropBlocked(val)
        }
        return nil
    }
    res["airDropForceUnmanagedDropTarget"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAirDropForceUnmanagedDropTarget(val)
        }
        return nil
    }
    res["airPlayForcePairingPasswordForOutgoingRequests"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAirPlayForcePairingPasswordForOutgoingRequests(val)
        }
        return nil
    }
    res["appleNewsBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppleNewsBlocked(val)
        }
        return nil
    }
    res["appleWatchBlockPairing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppleWatchBlockPairing(val)
        }
        return nil
    }
    res["appleWatchForceWristDetection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppleWatchForceWristDetection(val)
        }
        return nil
    }
    res["appsSingleAppModeList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAppsSingleAppModeList(res)
        }
        return nil
    }
    res["appStoreBlockAutomaticDownloads"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppStoreBlockAutomaticDownloads(val)
        }
        return nil
    }
    res["appStoreBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppStoreBlocked(val)
        }
        return nil
    }
    res["appStoreBlockInAppPurchases"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppStoreBlockInAppPurchases(val)
        }
        return nil
    }
    res["appStoreBlockUIAppInstallation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppStoreBlockUIAppInstallation(val)
        }
        return nil
    }
    res["appStoreRequirePassword"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppStoreRequirePassword(val)
        }
        return nil
    }
    res["appsVisibilityList"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetAppsVisibilityList(res)
        }
        return nil
    }
    res["appsVisibilityListType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAppListType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppsVisibilityListType(val.(*AppListType))
        }
        return nil
    }
    res["bluetoothBlockModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBluetoothBlockModification(val)
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
    res["cellularBlockGlobalBackgroundFetchWhileRoaming"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockGlobalBackgroundFetchWhileRoaming(val)
        }
        return nil
    }
    res["cellularBlockPerAppDataModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockPerAppDataModification(val)
        }
        return nil
    }
    res["cellularBlockPersonalHotspot"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockPersonalHotspot(val)
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
    res["certificatesBlockUntrustedTlsCertificates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCertificatesBlockUntrustedTlsCertificates(val)
        }
        return nil
    }
    res["classroomAppBlockRemoteScreenObservation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClassroomAppBlockRemoteScreenObservation(val)
        }
        return nil
    }
    res["classroomAppForceUnpromptedScreenObservation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetClassroomAppForceUnpromptedScreenObservation(val)
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
    res["configurationProfileBlockChanges"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfigurationProfileBlockChanges(val)
        }
        return nil
    }
    res["definitionLookupBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefinitionLookupBlocked(val)
        }
        return nil
    }
    res["deviceBlockEnableRestrictions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceBlockEnableRestrictions(val)
        }
        return nil
    }
    res["deviceBlockEraseContentAndSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceBlockEraseContentAndSettings(val)
        }
        return nil
    }
    res["deviceBlockNameModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceBlockNameModification(val)
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
    res["diagnosticDataBlockSubmissionModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDiagnosticDataBlockSubmissionModification(val)
        }
        return nil
    }
    res["documentsBlockManagedDocumentsInUnmanagedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDocumentsBlockManagedDocumentsInUnmanagedApps(val)
        }
        return nil
    }
    res["documentsBlockUnmanagedDocumentsInManagedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDocumentsBlockUnmanagedDocumentsInManagedApps(val)
        }
        return nil
    }
    res["emailInDomainSuffixes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetEmailInDomainSuffixes(res)
        }
        return nil
    }
    res["enterpriseAppBlockTrust"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseAppBlockTrust(val)
        }
        return nil
    }
    res["enterpriseAppBlockTrustModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnterpriseAppBlockTrustModification(val)
        }
        return nil
    }
    res["faceTimeBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFaceTimeBlocked(val)
        }
        return nil
    }
    res["findMyFriendsBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFindMyFriendsBlocked(val)
        }
        return nil
    }
    res["gameCenterBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGameCenterBlocked(val)
        }
        return nil
    }
    res["gamingBlockGameCenterFriends"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGamingBlockGameCenterFriends(val)
        }
        return nil
    }
    res["gamingBlockMultiplayer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGamingBlockMultiplayer(val)
        }
        return nil
    }
    res["hostPairingBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHostPairingBlocked(val)
        }
        return nil
    }
    res["iBooksStoreBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIBooksStoreBlocked(val)
        }
        return nil
    }
    res["iBooksStoreBlockErotica"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIBooksStoreBlockErotica(val)
        }
        return nil
    }
    res["iCloudBlockActivityContinuation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetICloudBlockActivityContinuation(val)
        }
        return nil
    }
    res["iCloudBlockBackup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetICloudBlockBackup(val)
        }
        return nil
    }
    res["iCloudBlockDocumentSync"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetICloudBlockDocumentSync(val)
        }
        return nil
    }
    res["iCloudBlockManagedAppsSync"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetICloudBlockManagedAppsSync(val)
        }
        return nil
    }
    res["iCloudBlockPhotoLibrary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetICloudBlockPhotoLibrary(val)
        }
        return nil
    }
    res["iCloudBlockPhotoStreamSync"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetICloudBlockPhotoStreamSync(val)
        }
        return nil
    }
    res["iCloudBlockSharedPhotoStream"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetICloudBlockSharedPhotoStream(val)
        }
        return nil
    }
    res["iCloudRequireEncryptedBackup"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetICloudRequireEncryptedBackup(val)
        }
        return nil
    }
    res["iTunesBlockExplicitContent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetITunesBlockExplicitContent(val)
        }
        return nil
    }
    res["iTunesBlockMusicService"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetITunesBlockMusicService(val)
        }
        return nil
    }
    res["iTunesBlockRadio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetITunesBlockRadio(val)
        }
        return nil
    }
    res["keyboardBlockAutoCorrect"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKeyboardBlockAutoCorrect(val)
        }
        return nil
    }
    res["keyboardBlockDictation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKeyboardBlockDictation(val)
        }
        return nil
    }
    res["keyboardBlockPredictive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKeyboardBlockPredictive(val)
        }
        return nil
    }
    res["keyboardBlockShortcuts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKeyboardBlockShortcuts(val)
        }
        return nil
    }
    res["keyboardBlockSpellCheck"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKeyboardBlockSpellCheck(val)
        }
        return nil
    }
    res["kioskModeAllowAssistiveSpeak"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowAssistiveSpeak(val)
        }
        return nil
    }
    res["kioskModeAllowAssistiveTouchSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowAssistiveTouchSettings(val)
        }
        return nil
    }
    res["kioskModeAllowAutoLock"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowAutoLock(val)
        }
        return nil
    }
    res["kioskModeAllowColorInversionSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowColorInversionSettings(val)
        }
        return nil
    }
    res["kioskModeAllowRingerSwitch"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowRingerSwitch(val)
        }
        return nil
    }
    res["kioskModeAllowScreenRotation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowScreenRotation(val)
        }
        return nil
    }
    res["kioskModeAllowSleepButton"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowSleepButton(val)
        }
        return nil
    }
    res["kioskModeAllowTouchscreen"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowTouchscreen(val)
        }
        return nil
    }
    res["kioskModeAllowVoiceOverSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowVoiceOverSettings(val)
        }
        return nil
    }
    res["kioskModeAllowVolumeButtons"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowVolumeButtons(val)
        }
        return nil
    }
    res["kioskModeAllowZoomSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAllowZoomSettings(val)
        }
        return nil
    }
    res["kioskModeAppStoreUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeAppStoreUrl(val)
        }
        return nil
    }
    res["kioskModeBuiltInAppId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeBuiltInAppId(val)
        }
        return nil
    }
    res["kioskModeManagedAppId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeManagedAppId(val)
        }
        return nil
    }
    res["kioskModeRequireAssistiveTouch"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeRequireAssistiveTouch(val)
        }
        return nil
    }
    res["kioskModeRequireColorInversion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeRequireColorInversion(val)
        }
        return nil
    }
    res["kioskModeRequireMonoAudio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeRequireMonoAudio(val)
        }
        return nil
    }
    res["kioskModeRequireVoiceOver"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeRequireVoiceOver(val)
        }
        return nil
    }
    res["kioskModeRequireZoom"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetKioskModeRequireZoom(val)
        }
        return nil
    }
    res["lockScreenBlockControlCenter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLockScreenBlockControlCenter(val)
        }
        return nil
    }
    res["lockScreenBlockNotificationView"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLockScreenBlockNotificationView(val)
        }
        return nil
    }
    res["lockScreenBlockPassbook"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLockScreenBlockPassbook(val)
        }
        return nil
    }
    res["lockScreenBlockTodayView"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLockScreenBlockTodayView(val)
        }
        return nil
    }
    res["mediaContentRatingApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRatingAppsType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaContentRatingApps(val.(*RatingAppsType))
        }
        return nil
    }
    res["mediaContentRatingAustralia"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaContentRatingAustraliaFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaContentRatingAustralia(val.(MediaContentRatingAustraliaable))
        }
        return nil
    }
    res["mediaContentRatingCanada"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaContentRatingCanadaFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaContentRatingCanada(val.(MediaContentRatingCanadaable))
        }
        return nil
    }
    res["mediaContentRatingFrance"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaContentRatingFranceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaContentRatingFrance(val.(MediaContentRatingFranceable))
        }
        return nil
    }
    res["mediaContentRatingGermany"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaContentRatingGermanyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaContentRatingGermany(val.(MediaContentRatingGermanyable))
        }
        return nil
    }
    res["mediaContentRatingIreland"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaContentRatingIrelandFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaContentRatingIreland(val.(MediaContentRatingIrelandable))
        }
        return nil
    }
    res["mediaContentRatingJapan"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaContentRatingJapanFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaContentRatingJapan(val.(MediaContentRatingJapanable))
        }
        return nil
    }
    res["mediaContentRatingNewZealand"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaContentRatingNewZealandFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaContentRatingNewZealand(val.(MediaContentRatingNewZealandable))
        }
        return nil
    }
    res["mediaContentRatingUnitedKingdom"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaContentRatingUnitedKingdomFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaContentRatingUnitedKingdom(val.(MediaContentRatingUnitedKingdomable))
        }
        return nil
    }
    res["mediaContentRatingUnitedStates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMediaContentRatingUnitedStatesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMediaContentRatingUnitedStates(val.(MediaContentRatingUnitedStatesable))
        }
        return nil
    }
    res["messagesBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMessagesBlocked(val)
        }
        return nil
    }
    res["networkUsageRules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIosNetworkUsageRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IosNetworkUsageRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IosNetworkUsageRuleable)
                }
            }
            m.SetNetworkUsageRules(res)
        }
        return nil
    }
    res["notificationsBlockSettingsModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotificationsBlockSettingsModification(val)
        }
        return nil
    }
    res["passcodeBlockFingerprintModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeBlockFingerprintModification(val)
        }
        return nil
    }
    res["passcodeBlockFingerprintUnlock"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeBlockFingerprintUnlock(val)
        }
        return nil
    }
    res["passcodeBlockModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeBlockModification(val)
        }
        return nil
    }
    res["passcodeBlockSimple"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeBlockSimple(val)
        }
        return nil
    }
    res["passcodeExpirationDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeExpirationDays(val)
        }
        return nil
    }
    res["passcodeMinimumCharacterSetCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeMinimumCharacterSetCount(val)
        }
        return nil
    }
    res["passcodeMinimumLength"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeMinimumLength(val)
        }
        return nil
    }
    res["passcodeMinutesOfInactivityBeforeLock"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeMinutesOfInactivityBeforeLock(val)
        }
        return nil
    }
    res["passcodeMinutesOfInactivityBeforeScreenTimeout"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeMinutesOfInactivityBeforeScreenTimeout(val)
        }
        return nil
    }
    res["passcodePreviousPasscodeBlockCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodePreviousPasscodeBlockCount(val)
        }
        return nil
    }
    res["passcodeRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeRequired(val)
        }
        return nil
    }
    res["passcodeRequiredType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRequiredPasswordType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeRequiredType(val.(*RequiredPasswordType))
        }
        return nil
    }
    res["passcodeSignInFailureCountBeforeWipe"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscodeSignInFailureCountBeforeWipe(val)
        }
        return nil
    }
    res["podcastsBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPodcastsBlocked(val)
        }
        return nil
    }
    res["safariBlockAutofill"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSafariBlockAutofill(val)
        }
        return nil
    }
    res["safariBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSafariBlocked(val)
        }
        return nil
    }
    res["safariBlockJavaScript"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSafariBlockJavaScript(val)
        }
        return nil
    }
    res["safariBlockPopups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSafariBlockPopups(val)
        }
        return nil
    }
    res["safariCookieSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWebBrowserCookieSettings)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSafariCookieSettings(val.(*WebBrowserCookieSettings))
        }
        return nil
    }
    res["safariManagedDomains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSafariManagedDomains(res)
        }
        return nil
    }
    res["safariPasswordAutoFillDomains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSafariPasswordAutoFillDomains(res)
        }
        return nil
    }
    res["safariRequireFraudWarning"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSafariRequireFraudWarning(val)
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
    res["siriBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiriBlocked(val)
        }
        return nil
    }
    res["siriBlockedWhenLocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiriBlockedWhenLocked(val)
        }
        return nil
    }
    res["siriBlockUserGeneratedContent"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiriBlockUserGeneratedContent(val)
        }
        return nil
    }
    res["siriRequireProfanityFilter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSiriRequireProfanityFilter(val)
        }
        return nil
    }
    res["spotlightBlockInternetResults"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSpotlightBlockInternetResults(val)
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
    res["wallpaperBlockModification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWallpaperBlockModification(val)
        }
        return nil
    }
    res["wiFiConnectOnlyToConfiguredNetworks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWiFiConnectOnlyToConfiguredNetworks(val)
        }
        return nil
    }
    return res
}
// GetFindMyFriendsBlocked gets the findMyFriendsBlocked property value. Indicates whether or not to block changes to Find My Friends when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetFindMyFriendsBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("findMyFriendsBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetGameCenterBlocked gets the gameCenterBlocked property value. Indicates whether or not to block the user from using Game Center when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetGameCenterBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("gameCenterBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetGamingBlockGameCenterFriends gets the gamingBlockGameCenterFriends property value. Indicates whether or not to block the user from having friends in Game Center. Requires a supervised device for iOS 13 and later.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetGamingBlockGameCenterFriends()(*bool) {
    val, err := m.GetBackingStore().Get("gamingBlockGameCenterFriends")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetGamingBlockMultiplayer gets the gamingBlockMultiplayer property value. Indicates whether or not to block the user from using multiplayer gaming. Requires a supervised device for iOS 13 and later.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetGamingBlockMultiplayer()(*bool) {
    val, err := m.GetBackingStore().Get("gamingBlockMultiplayer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHostPairingBlocked gets the hostPairingBlocked property value. indicates whether or not to allow host pairing to control the devices an iOS device can pair with when the iOS device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetHostPairingBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("hostPairingBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIBooksStoreBlocked gets the iBooksStoreBlocked property value. Indicates whether or not to block the user from using the iBooks Store when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetIBooksStoreBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("iBooksStoreBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIBooksStoreBlockErotica gets the iBooksStoreBlockErotica property value. Indicates whether or not to block the user from downloading media from the iBookstore that has been tagged as erotica.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetIBooksStoreBlockErotica()(*bool) {
    val, err := m.GetBackingStore().Get("iBooksStoreBlockErotica")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetICloudBlockActivityContinuation gets the iCloudBlockActivityContinuation property value. Indicates whether or not to block the user from continuing work they started on iOS device to another iOS or macOS device.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetICloudBlockActivityContinuation()(*bool) {
    val, err := m.GetBackingStore().Get("iCloudBlockActivityContinuation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetICloudBlockBackup gets the iCloudBlockBackup property value. Indicates whether or not to block iCloud backup. Requires a supervised device for iOS 13 and later.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetICloudBlockBackup()(*bool) {
    val, err := m.GetBackingStore().Get("iCloudBlockBackup")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetICloudBlockDocumentSync gets the iCloudBlockDocumentSync property value. Indicates whether or not to block iCloud document sync. Requires a supervised device for iOS 13 and later.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetICloudBlockDocumentSync()(*bool) {
    val, err := m.GetBackingStore().Get("iCloudBlockDocumentSync")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetICloudBlockManagedAppsSync gets the iCloudBlockManagedAppsSync property value. Indicates whether or not to block Managed Apps Cloud Sync.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetICloudBlockManagedAppsSync()(*bool) {
    val, err := m.GetBackingStore().Get("iCloudBlockManagedAppsSync")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetICloudBlockPhotoLibrary gets the iCloudBlockPhotoLibrary property value. Indicates whether or not to block iCloud Photo Library.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetICloudBlockPhotoLibrary()(*bool) {
    val, err := m.GetBackingStore().Get("iCloudBlockPhotoLibrary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetICloudBlockPhotoStreamSync gets the iCloudBlockPhotoStreamSync property value. Indicates whether or not to block iCloud Photo Stream Sync.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetICloudBlockPhotoStreamSync()(*bool) {
    val, err := m.GetBackingStore().Get("iCloudBlockPhotoStreamSync")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetICloudBlockSharedPhotoStream gets the iCloudBlockSharedPhotoStream property value. Indicates whether or not to block Shared Photo Stream.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetICloudBlockSharedPhotoStream()(*bool) {
    val, err := m.GetBackingStore().Get("iCloudBlockSharedPhotoStream")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetICloudRequireEncryptedBackup gets the iCloudRequireEncryptedBackup property value. Indicates whether or not to require backups to iCloud be encrypted.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetICloudRequireEncryptedBackup()(*bool) {
    val, err := m.GetBackingStore().Get("iCloudRequireEncryptedBackup")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetITunesBlockExplicitContent gets the iTunesBlockExplicitContent property value. Indicates whether or not to block the user from accessing explicit content in iTunes and the App Store. Requires a supervised device for iOS 13 and later.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetITunesBlockExplicitContent()(*bool) {
    val, err := m.GetBackingStore().Get("iTunesBlockExplicitContent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetITunesBlockMusicService gets the iTunesBlockMusicService property value. Indicates whether or not to block Music service and revert Music app to classic mode when the device is in supervised mode (iOS 9.3 and later and macOS 10.12 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetITunesBlockMusicService()(*bool) {
    val, err := m.GetBackingStore().Get("iTunesBlockMusicService")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetITunesBlockRadio gets the iTunesBlockRadio property value. Indicates whether or not to block the user from using iTunes Radio when the device is in supervised mode (iOS 9.3 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetITunesBlockRadio()(*bool) {
    val, err := m.GetBackingStore().Get("iTunesBlockRadio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKeyboardBlockAutoCorrect gets the keyboardBlockAutoCorrect property value. Indicates whether or not to block keyboard auto-correction when the device is in supervised mode (iOS 8.1.3 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKeyboardBlockAutoCorrect()(*bool) {
    val, err := m.GetBackingStore().Get("keyboardBlockAutoCorrect")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKeyboardBlockDictation gets the keyboardBlockDictation property value. Indicates whether or not to block the user from using dictation input when the device is in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKeyboardBlockDictation()(*bool) {
    val, err := m.GetBackingStore().Get("keyboardBlockDictation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKeyboardBlockPredictive gets the keyboardBlockPredictive property value. Indicates whether or not to block predictive keyboards when device is in supervised mode (iOS 8.1.3 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKeyboardBlockPredictive()(*bool) {
    val, err := m.GetBackingStore().Get("keyboardBlockPredictive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKeyboardBlockShortcuts gets the keyboardBlockShortcuts property value. Indicates whether or not to block keyboard shortcuts when the device is in supervised mode (iOS 9.0 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKeyboardBlockShortcuts()(*bool) {
    val, err := m.GetBackingStore().Get("keyboardBlockShortcuts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKeyboardBlockSpellCheck gets the keyboardBlockSpellCheck property value. Indicates whether or not to block keyboard spell-checking when the device is in supervised mode (iOS 8.1.3 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKeyboardBlockSpellCheck()(*bool) {
    val, err := m.GetBackingStore().Get("keyboardBlockSpellCheck")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowAssistiveSpeak gets the kioskModeAllowAssistiveSpeak property value. Indicates whether or not to allow assistive speak while in kiosk mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowAssistiveSpeak()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowAssistiveSpeak")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowAssistiveTouchSettings gets the kioskModeAllowAssistiveTouchSettings property value. Indicates whether or not to allow access to the Assistive Touch Settings while in kiosk mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowAssistiveTouchSettings()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowAssistiveTouchSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowAutoLock gets the kioskModeAllowAutoLock property value. Indicates whether or not to allow device auto lock while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockAutoLock instead.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowAutoLock()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowAutoLock")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowColorInversionSettings gets the kioskModeAllowColorInversionSettings property value. Indicates whether or not to allow access to the Color Inversion Settings while in kiosk mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowColorInversionSettings()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowColorInversionSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowRingerSwitch gets the kioskModeAllowRingerSwitch property value. Indicates whether or not to allow use of the ringer switch while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockRingerSwitch instead.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowRingerSwitch()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowRingerSwitch")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowScreenRotation gets the kioskModeAllowScreenRotation property value. Indicates whether or not to allow screen rotation while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockScreenRotation instead.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowScreenRotation()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowScreenRotation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowSleepButton gets the kioskModeAllowSleepButton property value. Indicates whether or not to allow use of the sleep button while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockSleepButton instead.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowSleepButton()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowSleepButton")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowTouchscreen gets the kioskModeAllowTouchscreen property value. Indicates whether or not to allow use of the touchscreen while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockTouchscreen instead.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowTouchscreen()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowTouchscreen")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowVoiceOverSettings gets the kioskModeAllowVoiceOverSettings property value. Indicates whether or not to allow access to the voice over settings while in kiosk mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowVoiceOverSettings()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowVoiceOverSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowVolumeButtons gets the kioskModeAllowVolumeButtons property value. Indicates whether or not to allow use of the volume buttons while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockVolumeButtons instead.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowVolumeButtons()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowVolumeButtons")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAllowZoomSettings gets the kioskModeAllowZoomSettings property value. Indicates whether or not to allow access to the zoom settings while in kiosk mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAllowZoomSettings()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeAllowZoomSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeAppStoreUrl gets the kioskModeAppStoreUrl property value. URL in the app store to the app to use for kiosk mode. Use if KioskModeManagedAppId is not known.
// returns a *string when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeAppStoreUrl()(*string) {
    val, err := m.GetBackingStore().Get("kioskModeAppStoreUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetKioskModeBuiltInAppId gets the kioskModeBuiltInAppId property value. ID for built-in apps to use for kiosk mode. Used when KioskModeManagedAppId and KioskModeAppStoreUrl are not set.
// returns a *string when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeBuiltInAppId()(*string) {
    val, err := m.GetBackingStore().Get("kioskModeBuiltInAppId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetKioskModeManagedAppId gets the kioskModeManagedAppId property value. Managed app id of the app to use for kiosk mode. If KioskModeManagedAppId is specified then KioskModeAppStoreUrl will be ignored.
// returns a *string when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeManagedAppId()(*string) {
    val, err := m.GetBackingStore().Get("kioskModeManagedAppId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetKioskModeRequireAssistiveTouch gets the kioskModeRequireAssistiveTouch property value. Indicates whether or not to require assistive touch while in kiosk mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeRequireAssistiveTouch()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeRequireAssistiveTouch")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeRequireColorInversion gets the kioskModeRequireColorInversion property value. Indicates whether or not to require color inversion while in kiosk mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeRequireColorInversion()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeRequireColorInversion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeRequireMonoAudio gets the kioskModeRequireMonoAudio property value. Indicates whether or not to require mono audio while in kiosk mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeRequireMonoAudio()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeRequireMonoAudio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeRequireVoiceOver gets the kioskModeRequireVoiceOver property value. Indicates whether or not to require voice over while in kiosk mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeRequireVoiceOver()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeRequireVoiceOver")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetKioskModeRequireZoom gets the kioskModeRequireZoom property value. Indicates whether or not to require zoom while in kiosk mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetKioskModeRequireZoom()(*bool) {
    val, err := m.GetBackingStore().Get("kioskModeRequireZoom")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLockScreenBlockControlCenter gets the lockScreenBlockControlCenter property value. Indicates whether or not to block the user from using control center on the lock screen.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetLockScreenBlockControlCenter()(*bool) {
    val, err := m.GetBackingStore().Get("lockScreenBlockControlCenter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLockScreenBlockNotificationView gets the lockScreenBlockNotificationView property value. Indicates whether or not to block the user from using the notification view on the lock screen.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetLockScreenBlockNotificationView()(*bool) {
    val, err := m.GetBackingStore().Get("lockScreenBlockNotificationView")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLockScreenBlockPassbook gets the lockScreenBlockPassbook property value. Indicates whether or not to block the user from using passbook when the device is locked.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetLockScreenBlockPassbook()(*bool) {
    val, err := m.GetBackingStore().Get("lockScreenBlockPassbook")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLockScreenBlockTodayView gets the lockScreenBlockTodayView property value. Indicates whether or not to block the user from using the Today View on the lock screen.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetLockScreenBlockTodayView()(*bool) {
    val, err := m.GetBackingStore().Get("lockScreenBlockTodayView")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMediaContentRatingApps gets the mediaContentRatingApps property value. Apps rating as in media content
// returns a *RatingAppsType when successful
func (m *IosGeneralDeviceConfiguration) GetMediaContentRatingApps()(*RatingAppsType) {
    val, err := m.GetBackingStore().Get("mediaContentRatingApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RatingAppsType)
    }
    return nil
}
// GetMediaContentRatingAustralia gets the mediaContentRatingAustralia property value. Media content rating settings for Australia
// returns a MediaContentRatingAustraliaable when successful
func (m *IosGeneralDeviceConfiguration) GetMediaContentRatingAustralia()(MediaContentRatingAustraliaable) {
    val, err := m.GetBackingStore().Get("mediaContentRatingAustralia")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaContentRatingAustraliaable)
    }
    return nil
}
// GetMediaContentRatingCanada gets the mediaContentRatingCanada property value. Media content rating settings for Canada
// returns a MediaContentRatingCanadaable when successful
func (m *IosGeneralDeviceConfiguration) GetMediaContentRatingCanada()(MediaContentRatingCanadaable) {
    val, err := m.GetBackingStore().Get("mediaContentRatingCanada")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaContentRatingCanadaable)
    }
    return nil
}
// GetMediaContentRatingFrance gets the mediaContentRatingFrance property value. Media content rating settings for France
// returns a MediaContentRatingFranceable when successful
func (m *IosGeneralDeviceConfiguration) GetMediaContentRatingFrance()(MediaContentRatingFranceable) {
    val, err := m.GetBackingStore().Get("mediaContentRatingFrance")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaContentRatingFranceable)
    }
    return nil
}
// GetMediaContentRatingGermany gets the mediaContentRatingGermany property value. Media content rating settings for Germany
// returns a MediaContentRatingGermanyable when successful
func (m *IosGeneralDeviceConfiguration) GetMediaContentRatingGermany()(MediaContentRatingGermanyable) {
    val, err := m.GetBackingStore().Get("mediaContentRatingGermany")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaContentRatingGermanyable)
    }
    return nil
}
// GetMediaContentRatingIreland gets the mediaContentRatingIreland property value. Media content rating settings for Ireland
// returns a MediaContentRatingIrelandable when successful
func (m *IosGeneralDeviceConfiguration) GetMediaContentRatingIreland()(MediaContentRatingIrelandable) {
    val, err := m.GetBackingStore().Get("mediaContentRatingIreland")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaContentRatingIrelandable)
    }
    return nil
}
// GetMediaContentRatingJapan gets the mediaContentRatingJapan property value. Media content rating settings for Japan
// returns a MediaContentRatingJapanable when successful
func (m *IosGeneralDeviceConfiguration) GetMediaContentRatingJapan()(MediaContentRatingJapanable) {
    val, err := m.GetBackingStore().Get("mediaContentRatingJapan")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaContentRatingJapanable)
    }
    return nil
}
// GetMediaContentRatingNewZealand gets the mediaContentRatingNewZealand property value. Media content rating settings for New Zealand
// returns a MediaContentRatingNewZealandable when successful
func (m *IosGeneralDeviceConfiguration) GetMediaContentRatingNewZealand()(MediaContentRatingNewZealandable) {
    val, err := m.GetBackingStore().Get("mediaContentRatingNewZealand")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaContentRatingNewZealandable)
    }
    return nil
}
// GetMediaContentRatingUnitedKingdom gets the mediaContentRatingUnitedKingdom property value. Media content rating settings for United Kingdom
// returns a MediaContentRatingUnitedKingdomable when successful
func (m *IosGeneralDeviceConfiguration) GetMediaContentRatingUnitedKingdom()(MediaContentRatingUnitedKingdomable) {
    val, err := m.GetBackingStore().Get("mediaContentRatingUnitedKingdom")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaContentRatingUnitedKingdomable)
    }
    return nil
}
// GetMediaContentRatingUnitedStates gets the mediaContentRatingUnitedStates property value. Media content rating settings for United States
// returns a MediaContentRatingUnitedStatesable when successful
func (m *IosGeneralDeviceConfiguration) GetMediaContentRatingUnitedStates()(MediaContentRatingUnitedStatesable) {
    val, err := m.GetBackingStore().Get("mediaContentRatingUnitedStates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MediaContentRatingUnitedStatesable)
    }
    return nil
}
// GetMessagesBlocked gets the messagesBlocked property value. Indicates whether or not to block the user from using the Messages app on the supervised device.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetMessagesBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("messagesBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNetworkUsageRules gets the networkUsageRules property value. List of managed apps and the network rules that applies to them. This collection can contain a maximum of 1000 elements.
// returns a []IosNetworkUsageRuleable when successful
func (m *IosGeneralDeviceConfiguration) GetNetworkUsageRules()([]IosNetworkUsageRuleable) {
    val, err := m.GetBackingStore().Get("networkUsageRules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IosNetworkUsageRuleable)
    }
    return nil
}
// GetNotificationsBlockSettingsModification gets the notificationsBlockSettingsModification property value. Indicates whether or not to allow notifications settings modification (iOS 9.3 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetNotificationsBlockSettingsModification()(*bool) {
    val, err := m.GetBackingStore().Get("notificationsBlockSettingsModification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasscodeBlockFingerprintModification gets the passcodeBlockFingerprintModification property value. Block modification of registered Touch ID fingerprints when in supervised mode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeBlockFingerprintModification()(*bool) {
    val, err := m.GetBackingStore().Get("passcodeBlockFingerprintModification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasscodeBlockFingerprintUnlock gets the passcodeBlockFingerprintUnlock property value. Indicates whether or not to block fingerprint unlock.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeBlockFingerprintUnlock()(*bool) {
    val, err := m.GetBackingStore().Get("passcodeBlockFingerprintUnlock")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasscodeBlockModification gets the passcodeBlockModification property value. Indicates whether or not to allow passcode modification on the supervised device (iOS 9.0 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeBlockModification()(*bool) {
    val, err := m.GetBackingStore().Get("passcodeBlockModification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasscodeBlockSimple gets the passcodeBlockSimple property value. Indicates whether or not to block simple passcodes.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeBlockSimple()(*bool) {
    val, err := m.GetBackingStore().Get("passcodeBlockSimple")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasscodeExpirationDays gets the passcodeExpirationDays property value. Number of days before the passcode expires. Valid values 1 to 65535
// returns a *int32 when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeExpirationDays()(*int32) {
    val, err := m.GetBackingStore().Get("passcodeExpirationDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasscodeMinimumCharacterSetCount gets the passcodeMinimumCharacterSetCount property value. Number of character sets a passcode must contain. Valid values 0 to 4
// returns a *int32 when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeMinimumCharacterSetCount()(*int32) {
    val, err := m.GetBackingStore().Get("passcodeMinimumCharacterSetCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasscodeMinimumLength gets the passcodeMinimumLength property value. Minimum length of passcode. Valid values 4 to 14
// returns a *int32 when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeMinimumLength()(*int32) {
    val, err := m.GetBackingStore().Get("passcodeMinimumLength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasscodeMinutesOfInactivityBeforeLock gets the passcodeMinutesOfInactivityBeforeLock property value. Minutes of inactivity before a passcode is required.
// returns a *int32 when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeMinutesOfInactivityBeforeLock()(*int32) {
    val, err := m.GetBackingStore().Get("passcodeMinutesOfInactivityBeforeLock")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasscodeMinutesOfInactivityBeforeScreenTimeout gets the passcodeMinutesOfInactivityBeforeScreenTimeout property value. Minutes of inactivity before the screen times out.
// returns a *int32 when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeMinutesOfInactivityBeforeScreenTimeout()(*int32) {
    val, err := m.GetBackingStore().Get("passcodeMinutesOfInactivityBeforeScreenTimeout")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasscodePreviousPasscodeBlockCount gets the passcodePreviousPasscodeBlockCount property value. Number of previous passcodes to block. Valid values 1 to 24
// returns a *int32 when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodePreviousPasscodeBlockCount()(*int32) {
    val, err := m.GetBackingStore().Get("passcodePreviousPasscodeBlockCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasscodeRequired gets the passcodeRequired property value. Indicates whether or not to require a passcode.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeRequired()(*bool) {
    val, err := m.GetBackingStore().Get("passcodeRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasscodeRequiredType gets the passcodeRequiredType property value. Possible values of required passwords.
// returns a *RequiredPasswordType when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeRequiredType()(*RequiredPasswordType) {
    val, err := m.GetBackingStore().Get("passcodeRequiredType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RequiredPasswordType)
    }
    return nil
}
// GetPasscodeSignInFailureCountBeforeWipe gets the passcodeSignInFailureCountBeforeWipe property value. Number of sign in failures allowed before wiping the device. Valid values 2 to 11
// returns a *int32 when successful
func (m *IosGeneralDeviceConfiguration) GetPasscodeSignInFailureCountBeforeWipe()(*int32) {
    val, err := m.GetBackingStore().Get("passcodeSignInFailureCountBeforeWipe")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPodcastsBlocked gets the podcastsBlocked property value. Indicates whether or not to block the user from using podcasts on the supervised device (iOS 8.0 and later).
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetPodcastsBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("podcastsBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSafariBlockAutofill gets the safariBlockAutofill property value. Indicates whether or not to block the user from using Auto fill in Safari. Requires a supervised device for iOS 13 and later.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetSafariBlockAutofill()(*bool) {
    val, err := m.GetBackingStore().Get("safariBlockAutofill")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSafariBlocked gets the safariBlocked property value. Indicates whether or not to block the user from using Safari. Requires a supervised device for iOS 13 and later.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetSafariBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("safariBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSafariBlockJavaScript gets the safariBlockJavaScript property value. Indicates whether or not to block JavaScript in Safari.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetSafariBlockJavaScript()(*bool) {
    val, err := m.GetBackingStore().Get("safariBlockJavaScript")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSafariBlockPopups gets the safariBlockPopups property value. Indicates whether or not to block popups in Safari.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetSafariBlockPopups()(*bool) {
    val, err := m.GetBackingStore().Get("safariBlockPopups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSafariCookieSettings gets the safariCookieSettings property value. Web Browser Cookie Settings.
// returns a *WebBrowserCookieSettings when successful
func (m *IosGeneralDeviceConfiguration) GetSafariCookieSettings()(*WebBrowserCookieSettings) {
    val, err := m.GetBackingStore().Get("safariCookieSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WebBrowserCookieSettings)
    }
    return nil
}
// GetSafariManagedDomains gets the safariManagedDomains property value. URLs matching the patterns listed here will be considered managed.
// returns a []string when successful
func (m *IosGeneralDeviceConfiguration) GetSafariManagedDomains()([]string) {
    val, err := m.GetBackingStore().Get("safariManagedDomains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSafariPasswordAutoFillDomains gets the safariPasswordAutoFillDomains property value. Users can save passwords in Safari only from URLs matching the patterns listed here. Applies to devices in supervised mode (iOS 9.3 and later).
// returns a []string when successful
func (m *IosGeneralDeviceConfiguration) GetSafariPasswordAutoFillDomains()([]string) {
    val, err := m.GetBackingStore().Get("safariPasswordAutoFillDomains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSafariRequireFraudWarning gets the safariRequireFraudWarning property value. Indicates whether or not to require fraud warning in Safari.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetSafariRequireFraudWarning()(*bool) {
    val, err := m.GetBackingStore().Get("safariRequireFraudWarning")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetScreenCaptureBlocked gets the screenCaptureBlocked property value. Indicates whether or not to block the user from taking Screenshots.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetScreenCaptureBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("screenCaptureBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSiriBlocked gets the siriBlocked property value. Indicates whether or not to block the user from using Siri.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetSiriBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("siriBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSiriBlockedWhenLocked gets the siriBlockedWhenLocked property value. Indicates whether or not to block the user from using Siri when locked.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetSiriBlockedWhenLocked()(*bool) {
    val, err := m.GetBackingStore().Get("siriBlockedWhenLocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSiriBlockUserGeneratedContent gets the siriBlockUserGeneratedContent property value. Indicates whether or not to block Siri from querying user-generated content when used on a supervised device.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetSiriBlockUserGeneratedContent()(*bool) {
    val, err := m.GetBackingStore().Get("siriBlockUserGeneratedContent")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSiriRequireProfanityFilter gets the siriRequireProfanityFilter property value. Indicates whether or not to prevent Siri from dictating, or speaking profane language on supervised device.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetSiriRequireProfanityFilter()(*bool) {
    val, err := m.GetBackingStore().Get("siriRequireProfanityFilter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSpotlightBlockInternetResults gets the spotlightBlockInternetResults property value. Indicates whether or not to block Spotlight search from returning internet results on supervised device.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetSpotlightBlockInternetResults()(*bool) {
    val, err := m.GetBackingStore().Get("spotlightBlockInternetResults")
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
func (m *IosGeneralDeviceConfiguration) GetVoiceDialingBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("voiceDialingBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWallpaperBlockModification gets the wallpaperBlockModification property value. Indicates whether or not to allow wallpaper modification on supervised device (iOS 9.0 and later) .
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetWallpaperBlockModification()(*bool) {
    val, err := m.GetBackingStore().Get("wallpaperBlockModification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWiFiConnectOnlyToConfiguredNetworks gets the wiFiConnectOnlyToConfiguredNetworks property value. Indicates whether or not to force the device to use only Wi-Fi networks from configuration profiles when the device is in supervised mode. Available for devices running iOS and iPadOS versions 14.4 and earlier. Devices running 14.5+ should use the setting, 'WiFiConnectToAllowedNetworksOnlyForced.
// returns a *bool when successful
func (m *IosGeneralDeviceConfiguration) GetWiFiConnectOnlyToConfiguredNetworks()(*bool) {
    val, err := m.GetBackingStore().Get("wiFiConnectOnlyToConfiguredNetworks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IosGeneralDeviceConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("accountBlockModification", m.GetAccountBlockModification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("activationLockAllowWhenSupervised", m.GetActivationLockAllowWhenSupervised())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("airDropBlocked", m.GetAirDropBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("airDropForceUnmanagedDropTarget", m.GetAirDropForceUnmanagedDropTarget())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("airPlayForcePairingPasswordForOutgoingRequests", m.GetAirPlayForcePairingPasswordForOutgoingRequests())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appleNewsBlocked", m.GetAppleNewsBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appleWatchBlockPairing", m.GetAppleWatchBlockPairing())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appleWatchForceWristDetection", m.GetAppleWatchForceWristDetection())
        if err != nil {
            return err
        }
    }
    if m.GetAppsSingleAppModeList() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppsSingleAppModeList()))
        for i, v := range m.GetAppsSingleAppModeList() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appsSingleAppModeList", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appStoreBlockAutomaticDownloads", m.GetAppStoreBlockAutomaticDownloads())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appStoreBlocked", m.GetAppStoreBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appStoreBlockInAppPurchases", m.GetAppStoreBlockInAppPurchases())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appStoreBlockUIAppInstallation", m.GetAppStoreBlockUIAppInstallation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("appStoreRequirePassword", m.GetAppStoreRequirePassword())
        if err != nil {
            return err
        }
    }
    if m.GetAppsVisibilityList() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAppsVisibilityList()))
        for i, v := range m.GetAppsVisibilityList() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("appsVisibilityList", cast)
        if err != nil {
            return err
        }
    }
    if m.GetAppsVisibilityListType() != nil {
        cast := (*m.GetAppsVisibilityListType()).String()
        err = writer.WriteStringValue("appsVisibilityListType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("bluetoothBlockModification", m.GetBluetoothBlockModification())
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
        err = writer.WriteBoolValue("cellularBlockGlobalBackgroundFetchWhileRoaming", m.GetCellularBlockGlobalBackgroundFetchWhileRoaming())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cellularBlockPerAppDataModification", m.GetCellularBlockPerAppDataModification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("cellularBlockPersonalHotspot", m.GetCellularBlockPersonalHotspot())
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
        err = writer.WriteBoolValue("certificatesBlockUntrustedTlsCertificates", m.GetCertificatesBlockUntrustedTlsCertificates())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("classroomAppBlockRemoteScreenObservation", m.GetClassroomAppBlockRemoteScreenObservation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("classroomAppForceUnpromptedScreenObservation", m.GetClassroomAppForceUnpromptedScreenObservation())
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
        err = writer.WriteBoolValue("configurationProfileBlockChanges", m.GetConfigurationProfileBlockChanges())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("definitionLookupBlocked", m.GetDefinitionLookupBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("deviceBlockEnableRestrictions", m.GetDeviceBlockEnableRestrictions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("deviceBlockEraseContentAndSettings", m.GetDeviceBlockEraseContentAndSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("deviceBlockNameModification", m.GetDeviceBlockNameModification())
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
        err = writer.WriteBoolValue("diagnosticDataBlockSubmissionModification", m.GetDiagnosticDataBlockSubmissionModification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("documentsBlockManagedDocumentsInUnmanagedApps", m.GetDocumentsBlockManagedDocumentsInUnmanagedApps())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("documentsBlockUnmanagedDocumentsInManagedApps", m.GetDocumentsBlockUnmanagedDocumentsInManagedApps())
        if err != nil {
            return err
        }
    }
    if m.GetEmailInDomainSuffixes() != nil {
        err = writer.WriteCollectionOfStringValues("emailInDomainSuffixes", m.GetEmailInDomainSuffixes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("enterpriseAppBlockTrust", m.GetEnterpriseAppBlockTrust())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("enterpriseAppBlockTrustModification", m.GetEnterpriseAppBlockTrustModification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("faceTimeBlocked", m.GetFaceTimeBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("findMyFriendsBlocked", m.GetFindMyFriendsBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("gameCenterBlocked", m.GetGameCenterBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("gamingBlockGameCenterFriends", m.GetGamingBlockGameCenterFriends())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("gamingBlockMultiplayer", m.GetGamingBlockMultiplayer())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hostPairingBlocked", m.GetHostPairingBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iBooksStoreBlocked", m.GetIBooksStoreBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iBooksStoreBlockErotica", m.GetIBooksStoreBlockErotica())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iCloudBlockActivityContinuation", m.GetICloudBlockActivityContinuation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iCloudBlockBackup", m.GetICloudBlockBackup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iCloudBlockDocumentSync", m.GetICloudBlockDocumentSync())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iCloudBlockManagedAppsSync", m.GetICloudBlockManagedAppsSync())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iCloudBlockPhotoLibrary", m.GetICloudBlockPhotoLibrary())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iCloudBlockPhotoStreamSync", m.GetICloudBlockPhotoStreamSync())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iCloudBlockSharedPhotoStream", m.GetICloudBlockSharedPhotoStream())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iCloudRequireEncryptedBackup", m.GetICloudRequireEncryptedBackup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iTunesBlockExplicitContent", m.GetITunesBlockExplicitContent())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iTunesBlockMusicService", m.GetITunesBlockMusicService())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iTunesBlockRadio", m.GetITunesBlockRadio())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("keyboardBlockAutoCorrect", m.GetKeyboardBlockAutoCorrect())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("keyboardBlockDictation", m.GetKeyboardBlockDictation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("keyboardBlockPredictive", m.GetKeyboardBlockPredictive())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("keyboardBlockShortcuts", m.GetKeyboardBlockShortcuts())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("keyboardBlockSpellCheck", m.GetKeyboardBlockSpellCheck())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowAssistiveSpeak", m.GetKioskModeAllowAssistiveSpeak())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowAssistiveTouchSettings", m.GetKioskModeAllowAssistiveTouchSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowAutoLock", m.GetKioskModeAllowAutoLock())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowColorInversionSettings", m.GetKioskModeAllowColorInversionSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowRingerSwitch", m.GetKioskModeAllowRingerSwitch())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowScreenRotation", m.GetKioskModeAllowScreenRotation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowSleepButton", m.GetKioskModeAllowSleepButton())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowTouchscreen", m.GetKioskModeAllowTouchscreen())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowVoiceOverSettings", m.GetKioskModeAllowVoiceOverSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowVolumeButtons", m.GetKioskModeAllowVolumeButtons())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeAllowZoomSettings", m.GetKioskModeAllowZoomSettings())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("kioskModeAppStoreUrl", m.GetKioskModeAppStoreUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("kioskModeBuiltInAppId", m.GetKioskModeBuiltInAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("kioskModeManagedAppId", m.GetKioskModeManagedAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeRequireAssistiveTouch", m.GetKioskModeRequireAssistiveTouch())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeRequireColorInversion", m.GetKioskModeRequireColorInversion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeRequireMonoAudio", m.GetKioskModeRequireMonoAudio())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeRequireVoiceOver", m.GetKioskModeRequireVoiceOver())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("kioskModeRequireZoom", m.GetKioskModeRequireZoom())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("lockScreenBlockControlCenter", m.GetLockScreenBlockControlCenter())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("lockScreenBlockNotificationView", m.GetLockScreenBlockNotificationView())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("lockScreenBlockPassbook", m.GetLockScreenBlockPassbook())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("lockScreenBlockTodayView", m.GetLockScreenBlockTodayView())
        if err != nil {
            return err
        }
    }
    if m.GetMediaContentRatingApps() != nil {
        cast := (*m.GetMediaContentRatingApps()).String()
        err = writer.WriteStringValue("mediaContentRatingApps", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaContentRatingAustralia", m.GetMediaContentRatingAustralia())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaContentRatingCanada", m.GetMediaContentRatingCanada())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaContentRatingFrance", m.GetMediaContentRatingFrance())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaContentRatingGermany", m.GetMediaContentRatingGermany())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaContentRatingIreland", m.GetMediaContentRatingIreland())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaContentRatingJapan", m.GetMediaContentRatingJapan())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaContentRatingNewZealand", m.GetMediaContentRatingNewZealand())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaContentRatingUnitedKingdom", m.GetMediaContentRatingUnitedKingdom())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("mediaContentRatingUnitedStates", m.GetMediaContentRatingUnitedStates())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("messagesBlocked", m.GetMessagesBlocked())
        if err != nil {
            return err
        }
    }
    if m.GetNetworkUsageRules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetNetworkUsageRules()))
        for i, v := range m.GetNetworkUsageRules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("networkUsageRules", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("notificationsBlockSettingsModification", m.GetNotificationsBlockSettingsModification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passcodeBlockFingerprintModification", m.GetPasscodeBlockFingerprintModification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passcodeBlockFingerprintUnlock", m.GetPasscodeBlockFingerprintUnlock())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passcodeBlockModification", m.GetPasscodeBlockModification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passcodeBlockSimple", m.GetPasscodeBlockSimple())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passcodeExpirationDays", m.GetPasscodeExpirationDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passcodeMinimumCharacterSetCount", m.GetPasscodeMinimumCharacterSetCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passcodeMinimumLength", m.GetPasscodeMinimumLength())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passcodeMinutesOfInactivityBeforeLock", m.GetPasscodeMinutesOfInactivityBeforeLock())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passcodeMinutesOfInactivityBeforeScreenTimeout", m.GetPasscodeMinutesOfInactivityBeforeScreenTimeout())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passcodePreviousPasscodeBlockCount", m.GetPasscodePreviousPasscodeBlockCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passcodeRequired", m.GetPasscodeRequired())
        if err != nil {
            return err
        }
    }
    if m.GetPasscodeRequiredType() != nil {
        cast := (*m.GetPasscodeRequiredType()).String()
        err = writer.WriteStringValue("passcodeRequiredType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("passcodeSignInFailureCountBeforeWipe", m.GetPasscodeSignInFailureCountBeforeWipe())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("podcastsBlocked", m.GetPodcastsBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("safariBlockAutofill", m.GetSafariBlockAutofill())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("safariBlocked", m.GetSafariBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("safariBlockJavaScript", m.GetSafariBlockJavaScript())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("safariBlockPopups", m.GetSafariBlockPopups())
        if err != nil {
            return err
        }
    }
    if m.GetSafariCookieSettings() != nil {
        cast := (*m.GetSafariCookieSettings()).String()
        err = writer.WriteStringValue("safariCookieSettings", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetSafariManagedDomains() != nil {
        err = writer.WriteCollectionOfStringValues("safariManagedDomains", m.GetSafariManagedDomains())
        if err != nil {
            return err
        }
    }
    if m.GetSafariPasswordAutoFillDomains() != nil {
        err = writer.WriteCollectionOfStringValues("safariPasswordAutoFillDomains", m.GetSafariPasswordAutoFillDomains())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("safariRequireFraudWarning", m.GetSafariRequireFraudWarning())
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
        err = writer.WriteBoolValue("siriBlocked", m.GetSiriBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("siriBlockedWhenLocked", m.GetSiriBlockedWhenLocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("siriBlockUserGeneratedContent", m.GetSiriBlockUserGeneratedContent())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("siriRequireProfanityFilter", m.GetSiriRequireProfanityFilter())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("spotlightBlockInternetResults", m.GetSpotlightBlockInternetResults())
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
        err = writer.WriteBoolValue("wallpaperBlockModification", m.GetWallpaperBlockModification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wiFiConnectOnlyToConfiguredNetworks", m.GetWiFiConnectOnlyToConfiguredNetworks())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccountBlockModification sets the accountBlockModification property value. Indicates whether or not to allow account modification when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetAccountBlockModification(value *bool)() {
    err := m.GetBackingStore().Set("accountBlockModification", value)
    if err != nil {
        panic(err)
    }
}
// SetActivationLockAllowWhenSupervised sets the activationLockAllowWhenSupervised property value. Indicates whether or not to allow activation lock when the device is in the supervised mode.
func (m *IosGeneralDeviceConfiguration) SetActivationLockAllowWhenSupervised(value *bool)() {
    err := m.GetBackingStore().Set("activationLockAllowWhenSupervised", value)
    if err != nil {
        panic(err)
    }
}
// SetAirDropBlocked sets the airDropBlocked property value. Indicates whether or not to allow AirDrop when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetAirDropBlocked(value *bool)() {
    err := m.GetBackingStore().Set("airDropBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetAirDropForceUnmanagedDropTarget sets the airDropForceUnmanagedDropTarget property value. Indicates whether or not to cause AirDrop to be considered an unmanaged drop target (iOS 9.0 and later).
func (m *IosGeneralDeviceConfiguration) SetAirDropForceUnmanagedDropTarget(value *bool)() {
    err := m.GetBackingStore().Set("airDropForceUnmanagedDropTarget", value)
    if err != nil {
        panic(err)
    }
}
// SetAirPlayForcePairingPasswordForOutgoingRequests sets the airPlayForcePairingPasswordForOutgoingRequests property value. Indicates whether or not to enforce all devices receiving AirPlay requests from this device to use a pairing password.
func (m *IosGeneralDeviceConfiguration) SetAirPlayForcePairingPasswordForOutgoingRequests(value *bool)() {
    err := m.GetBackingStore().Set("airPlayForcePairingPasswordForOutgoingRequests", value)
    if err != nil {
        panic(err)
    }
}
// SetAppleNewsBlocked sets the appleNewsBlocked property value. Indicates whether or not to block the user from using News when the device is in supervised mode (iOS 9.0 and later).
func (m *IosGeneralDeviceConfiguration) SetAppleNewsBlocked(value *bool)() {
    err := m.GetBackingStore().Set("appleNewsBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetAppleWatchBlockPairing sets the appleWatchBlockPairing property value. Indicates whether or not to allow Apple Watch pairing when the device is in supervised mode (iOS 9.0 and later).
func (m *IosGeneralDeviceConfiguration) SetAppleWatchBlockPairing(value *bool)() {
    err := m.GetBackingStore().Set("appleWatchBlockPairing", value)
    if err != nil {
        panic(err)
    }
}
// SetAppleWatchForceWristDetection sets the appleWatchForceWristDetection property value. Indicates whether or not to force a paired Apple Watch to use Wrist Detection (iOS 8.2 and later).
func (m *IosGeneralDeviceConfiguration) SetAppleWatchForceWristDetection(value *bool)() {
    err := m.GetBackingStore().Set("appleWatchForceWristDetection", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsSingleAppModeList sets the appsSingleAppModeList property value. Gets or sets the list of iOS apps allowed to autonomously enter Single App Mode. Supervised only. iOS 7.0 and later. This collection can contain a maximum of 500 elements.
func (m *IosGeneralDeviceConfiguration) SetAppsSingleAppModeList(value []AppListItemable)() {
    err := m.GetBackingStore().Set("appsSingleAppModeList", value)
    if err != nil {
        panic(err)
    }
}
// SetAppStoreBlockAutomaticDownloads sets the appStoreBlockAutomaticDownloads property value. Indicates whether or not to block the automatic downloading of apps purchased on other devices when the device is in supervised mode (iOS 9.0 and later).
func (m *IosGeneralDeviceConfiguration) SetAppStoreBlockAutomaticDownloads(value *bool)() {
    err := m.GetBackingStore().Set("appStoreBlockAutomaticDownloads", value)
    if err != nil {
        panic(err)
    }
}
// SetAppStoreBlocked sets the appStoreBlocked property value. Indicates whether or not to block the user from using the App Store. Requires a supervised device for iOS 13 and later.
func (m *IosGeneralDeviceConfiguration) SetAppStoreBlocked(value *bool)() {
    err := m.GetBackingStore().Set("appStoreBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetAppStoreBlockInAppPurchases sets the appStoreBlockInAppPurchases property value. Indicates whether or not to block the user from making in app purchases.
func (m *IosGeneralDeviceConfiguration) SetAppStoreBlockInAppPurchases(value *bool)() {
    err := m.GetBackingStore().Set("appStoreBlockInAppPurchases", value)
    if err != nil {
        panic(err)
    }
}
// SetAppStoreBlockUIAppInstallation sets the appStoreBlockUIAppInstallation property value. Indicates whether or not to block the App Store app, not restricting installation through Host apps. Applies to supervised mode only (iOS 9.0 and later).
func (m *IosGeneralDeviceConfiguration) SetAppStoreBlockUIAppInstallation(value *bool)() {
    err := m.GetBackingStore().Set("appStoreBlockUIAppInstallation", value)
    if err != nil {
        panic(err)
    }
}
// SetAppStoreRequirePassword sets the appStoreRequirePassword property value. Indicates whether or not to require a password when using the app store.
func (m *IosGeneralDeviceConfiguration) SetAppStoreRequirePassword(value *bool)() {
    err := m.GetBackingStore().Set("appStoreRequirePassword", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsVisibilityList sets the appsVisibilityList property value. List of apps in the visibility list (either visible/launchable apps list or hidden/unlaunchable apps list, controlled by AppsVisibilityListType) (iOS 9.3 and later). This collection can contain a maximum of 10000 elements.
func (m *IosGeneralDeviceConfiguration) SetAppsVisibilityList(value []AppListItemable)() {
    err := m.GetBackingStore().Set("appsVisibilityList", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsVisibilityListType sets the appsVisibilityListType property value. Possible values of the compliance app list.
func (m *IosGeneralDeviceConfiguration) SetAppsVisibilityListType(value *AppListType)() {
    err := m.GetBackingStore().Set("appsVisibilityListType", value)
    if err != nil {
        panic(err)
    }
}
// SetBluetoothBlockModification sets the bluetoothBlockModification property value. Indicates whether or not to allow modification of Bluetooth settings when the device is in supervised mode (iOS 10.0 and later).
func (m *IosGeneralDeviceConfiguration) SetBluetoothBlockModification(value *bool)() {
    err := m.GetBackingStore().Set("bluetoothBlockModification", value)
    if err != nil {
        panic(err)
    }
}
// SetCameraBlocked sets the cameraBlocked property value. Indicates whether or not to block the user from accessing the camera of the device. Requires a supervised device for iOS 13 and later.
func (m *IosGeneralDeviceConfiguration) SetCameraBlocked(value *bool)() {
    err := m.GetBackingStore().Set("cameraBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockDataRoaming sets the cellularBlockDataRoaming property value. Indicates whether or not to block data roaming.
func (m *IosGeneralDeviceConfiguration) SetCellularBlockDataRoaming(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockDataRoaming", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockGlobalBackgroundFetchWhileRoaming sets the cellularBlockGlobalBackgroundFetchWhileRoaming property value. Indicates whether or not to block global background fetch while roaming.
func (m *IosGeneralDeviceConfiguration) SetCellularBlockGlobalBackgroundFetchWhileRoaming(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockGlobalBackgroundFetchWhileRoaming", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockPerAppDataModification sets the cellularBlockPerAppDataModification property value. Indicates whether or not to allow changes to cellular app data usage settings when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetCellularBlockPerAppDataModification(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockPerAppDataModification", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockPersonalHotspot sets the cellularBlockPersonalHotspot property value. Indicates whether or not to block Personal Hotspot.
func (m *IosGeneralDeviceConfiguration) SetCellularBlockPersonalHotspot(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockPersonalHotspot", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockVoiceRoaming sets the cellularBlockVoiceRoaming property value. Indicates whether or not to block voice roaming.
func (m *IosGeneralDeviceConfiguration) SetCellularBlockVoiceRoaming(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockVoiceRoaming", value)
    if err != nil {
        panic(err)
    }
}
// SetCertificatesBlockUntrustedTlsCertificates sets the certificatesBlockUntrustedTlsCertificates property value. Indicates whether or not to block untrusted TLS certificates.
func (m *IosGeneralDeviceConfiguration) SetCertificatesBlockUntrustedTlsCertificates(value *bool)() {
    err := m.GetBackingStore().Set("certificatesBlockUntrustedTlsCertificates", value)
    if err != nil {
        panic(err)
    }
}
// SetClassroomAppBlockRemoteScreenObservation sets the classroomAppBlockRemoteScreenObservation property value. Indicates whether or not to allow remote screen observation by Classroom app when the device is in supervised mode (iOS 9.3 and later).
func (m *IosGeneralDeviceConfiguration) SetClassroomAppBlockRemoteScreenObservation(value *bool)() {
    err := m.GetBackingStore().Set("classroomAppBlockRemoteScreenObservation", value)
    if err != nil {
        panic(err)
    }
}
// SetClassroomAppForceUnpromptedScreenObservation sets the classroomAppForceUnpromptedScreenObservation property value. Indicates whether or not to automatically give permission to the teacher of a managed course on the Classroom app to view a student's screen without prompting when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetClassroomAppForceUnpromptedScreenObservation(value *bool)() {
    err := m.GetBackingStore().Set("classroomAppForceUnpromptedScreenObservation", value)
    if err != nil {
        panic(err)
    }
}
// SetCompliantAppListType sets the compliantAppListType property value. Possible values of the compliance app list.
func (m *IosGeneralDeviceConfiguration) SetCompliantAppListType(value *AppListType)() {
    err := m.GetBackingStore().Set("compliantAppListType", value)
    if err != nil {
        panic(err)
    }
}
// SetCompliantAppsList sets the compliantAppsList property value. List of apps in the compliance (either allow list or block list, controlled by CompliantAppListType). This collection can contain a maximum of 10000 elements.
func (m *IosGeneralDeviceConfiguration) SetCompliantAppsList(value []AppListItemable)() {
    err := m.GetBackingStore().Set("compliantAppsList", value)
    if err != nil {
        panic(err)
    }
}
// SetConfigurationProfileBlockChanges sets the configurationProfileBlockChanges property value. Indicates whether or not to block the user from installing configuration profiles and certificates interactively when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetConfigurationProfileBlockChanges(value *bool)() {
    err := m.GetBackingStore().Set("configurationProfileBlockChanges", value)
    if err != nil {
        panic(err)
    }
}
// SetDefinitionLookupBlocked sets the definitionLookupBlocked property value. Indicates whether or not to block definition lookup when the device is in supervised mode (iOS 8.1.3 and later ).
func (m *IosGeneralDeviceConfiguration) SetDefinitionLookupBlocked(value *bool)() {
    err := m.GetBackingStore().Set("definitionLookupBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceBlockEnableRestrictions sets the deviceBlockEnableRestrictions property value. Indicates whether or not to allow the user to enables restrictions in the device settings when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetDeviceBlockEnableRestrictions(value *bool)() {
    err := m.GetBackingStore().Set("deviceBlockEnableRestrictions", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceBlockEraseContentAndSettings sets the deviceBlockEraseContentAndSettings property value. Indicates whether or not to allow the use of the 'Erase all content and settings' option on the device when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetDeviceBlockEraseContentAndSettings(value *bool)() {
    err := m.GetBackingStore().Set("deviceBlockEraseContentAndSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceBlockNameModification sets the deviceBlockNameModification property value. Indicates whether or not to allow device name modification when the device is in supervised mode (iOS 9.0 and later).
func (m *IosGeneralDeviceConfiguration) SetDeviceBlockNameModification(value *bool)() {
    err := m.GetBackingStore().Set("deviceBlockNameModification", value)
    if err != nil {
        panic(err)
    }
}
// SetDiagnosticDataBlockSubmission sets the diagnosticDataBlockSubmission property value. Indicates whether or not to block diagnostic data submission.
func (m *IosGeneralDeviceConfiguration) SetDiagnosticDataBlockSubmission(value *bool)() {
    err := m.GetBackingStore().Set("diagnosticDataBlockSubmission", value)
    if err != nil {
        panic(err)
    }
}
// SetDiagnosticDataBlockSubmissionModification sets the diagnosticDataBlockSubmissionModification property value. Indicates whether or not to allow diagnostics submission settings modification when the device is in supervised mode (iOS 9.3.2 and later).
func (m *IosGeneralDeviceConfiguration) SetDiagnosticDataBlockSubmissionModification(value *bool)() {
    err := m.GetBackingStore().Set("diagnosticDataBlockSubmissionModification", value)
    if err != nil {
        panic(err)
    }
}
// SetDocumentsBlockManagedDocumentsInUnmanagedApps sets the documentsBlockManagedDocumentsInUnmanagedApps property value. Indicates whether or not to block the user from viewing managed documents in unmanaged apps.
func (m *IosGeneralDeviceConfiguration) SetDocumentsBlockManagedDocumentsInUnmanagedApps(value *bool)() {
    err := m.GetBackingStore().Set("documentsBlockManagedDocumentsInUnmanagedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetDocumentsBlockUnmanagedDocumentsInManagedApps sets the documentsBlockUnmanagedDocumentsInManagedApps property value. Indicates whether or not to block the user from viewing unmanaged documents in managed apps.
func (m *IosGeneralDeviceConfiguration) SetDocumentsBlockUnmanagedDocumentsInManagedApps(value *bool)() {
    err := m.GetBackingStore().Set("documentsBlockUnmanagedDocumentsInManagedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetEmailInDomainSuffixes sets the emailInDomainSuffixes property value. An email address lacking a suffix that matches any of these strings will be considered out-of-domain.
func (m *IosGeneralDeviceConfiguration) SetEmailInDomainSuffixes(value []string)() {
    err := m.GetBackingStore().Set("emailInDomainSuffixes", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseAppBlockTrust sets the enterpriseAppBlockTrust property value. Indicates whether or not to block the user from trusting an enterprise app.
func (m *IosGeneralDeviceConfiguration) SetEnterpriseAppBlockTrust(value *bool)() {
    err := m.GetBackingStore().Set("enterpriseAppBlockTrust", value)
    if err != nil {
        panic(err)
    }
}
// SetEnterpriseAppBlockTrustModification sets the enterpriseAppBlockTrustModification property value. [Deprecated] Configuring this setting and setting the value to 'true' has no effect on the device.
func (m *IosGeneralDeviceConfiguration) SetEnterpriseAppBlockTrustModification(value *bool)() {
    err := m.GetBackingStore().Set("enterpriseAppBlockTrustModification", value)
    if err != nil {
        panic(err)
    }
}
// SetFaceTimeBlocked sets the faceTimeBlocked property value. Indicates whether or not to block the user from using FaceTime. Requires a supervised device for iOS 13 and later.
func (m *IosGeneralDeviceConfiguration) SetFaceTimeBlocked(value *bool)() {
    err := m.GetBackingStore().Set("faceTimeBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetFindMyFriendsBlocked sets the findMyFriendsBlocked property value. Indicates whether or not to block changes to Find My Friends when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetFindMyFriendsBlocked(value *bool)() {
    err := m.GetBackingStore().Set("findMyFriendsBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetGameCenterBlocked sets the gameCenterBlocked property value. Indicates whether or not to block the user from using Game Center when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetGameCenterBlocked(value *bool)() {
    err := m.GetBackingStore().Set("gameCenterBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetGamingBlockGameCenterFriends sets the gamingBlockGameCenterFriends property value. Indicates whether or not to block the user from having friends in Game Center. Requires a supervised device for iOS 13 and later.
func (m *IosGeneralDeviceConfiguration) SetGamingBlockGameCenterFriends(value *bool)() {
    err := m.GetBackingStore().Set("gamingBlockGameCenterFriends", value)
    if err != nil {
        panic(err)
    }
}
// SetGamingBlockMultiplayer sets the gamingBlockMultiplayer property value. Indicates whether or not to block the user from using multiplayer gaming. Requires a supervised device for iOS 13 and later.
func (m *IosGeneralDeviceConfiguration) SetGamingBlockMultiplayer(value *bool)() {
    err := m.GetBackingStore().Set("gamingBlockMultiplayer", value)
    if err != nil {
        panic(err)
    }
}
// SetHostPairingBlocked sets the hostPairingBlocked property value. indicates whether or not to allow host pairing to control the devices an iOS device can pair with when the iOS device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetHostPairingBlocked(value *bool)() {
    err := m.GetBackingStore().Set("hostPairingBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetIBooksStoreBlocked sets the iBooksStoreBlocked property value. Indicates whether or not to block the user from using the iBooks Store when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetIBooksStoreBlocked(value *bool)() {
    err := m.GetBackingStore().Set("iBooksStoreBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetIBooksStoreBlockErotica sets the iBooksStoreBlockErotica property value. Indicates whether or not to block the user from downloading media from the iBookstore that has been tagged as erotica.
func (m *IosGeneralDeviceConfiguration) SetIBooksStoreBlockErotica(value *bool)() {
    err := m.GetBackingStore().Set("iBooksStoreBlockErotica", value)
    if err != nil {
        panic(err)
    }
}
// SetICloudBlockActivityContinuation sets the iCloudBlockActivityContinuation property value. Indicates whether or not to block the user from continuing work they started on iOS device to another iOS or macOS device.
func (m *IosGeneralDeviceConfiguration) SetICloudBlockActivityContinuation(value *bool)() {
    err := m.GetBackingStore().Set("iCloudBlockActivityContinuation", value)
    if err != nil {
        panic(err)
    }
}
// SetICloudBlockBackup sets the iCloudBlockBackup property value. Indicates whether or not to block iCloud backup. Requires a supervised device for iOS 13 and later.
func (m *IosGeneralDeviceConfiguration) SetICloudBlockBackup(value *bool)() {
    err := m.GetBackingStore().Set("iCloudBlockBackup", value)
    if err != nil {
        panic(err)
    }
}
// SetICloudBlockDocumentSync sets the iCloudBlockDocumentSync property value. Indicates whether or not to block iCloud document sync. Requires a supervised device for iOS 13 and later.
func (m *IosGeneralDeviceConfiguration) SetICloudBlockDocumentSync(value *bool)() {
    err := m.GetBackingStore().Set("iCloudBlockDocumentSync", value)
    if err != nil {
        panic(err)
    }
}
// SetICloudBlockManagedAppsSync sets the iCloudBlockManagedAppsSync property value. Indicates whether or not to block Managed Apps Cloud Sync.
func (m *IosGeneralDeviceConfiguration) SetICloudBlockManagedAppsSync(value *bool)() {
    err := m.GetBackingStore().Set("iCloudBlockManagedAppsSync", value)
    if err != nil {
        panic(err)
    }
}
// SetICloudBlockPhotoLibrary sets the iCloudBlockPhotoLibrary property value. Indicates whether or not to block iCloud Photo Library.
func (m *IosGeneralDeviceConfiguration) SetICloudBlockPhotoLibrary(value *bool)() {
    err := m.GetBackingStore().Set("iCloudBlockPhotoLibrary", value)
    if err != nil {
        panic(err)
    }
}
// SetICloudBlockPhotoStreamSync sets the iCloudBlockPhotoStreamSync property value. Indicates whether or not to block iCloud Photo Stream Sync.
func (m *IosGeneralDeviceConfiguration) SetICloudBlockPhotoStreamSync(value *bool)() {
    err := m.GetBackingStore().Set("iCloudBlockPhotoStreamSync", value)
    if err != nil {
        panic(err)
    }
}
// SetICloudBlockSharedPhotoStream sets the iCloudBlockSharedPhotoStream property value. Indicates whether or not to block Shared Photo Stream.
func (m *IosGeneralDeviceConfiguration) SetICloudBlockSharedPhotoStream(value *bool)() {
    err := m.GetBackingStore().Set("iCloudBlockSharedPhotoStream", value)
    if err != nil {
        panic(err)
    }
}
// SetICloudRequireEncryptedBackup sets the iCloudRequireEncryptedBackup property value. Indicates whether or not to require backups to iCloud be encrypted.
func (m *IosGeneralDeviceConfiguration) SetICloudRequireEncryptedBackup(value *bool)() {
    err := m.GetBackingStore().Set("iCloudRequireEncryptedBackup", value)
    if err != nil {
        panic(err)
    }
}
// SetITunesBlockExplicitContent sets the iTunesBlockExplicitContent property value. Indicates whether or not to block the user from accessing explicit content in iTunes and the App Store. Requires a supervised device for iOS 13 and later.
func (m *IosGeneralDeviceConfiguration) SetITunesBlockExplicitContent(value *bool)() {
    err := m.GetBackingStore().Set("iTunesBlockExplicitContent", value)
    if err != nil {
        panic(err)
    }
}
// SetITunesBlockMusicService sets the iTunesBlockMusicService property value. Indicates whether or not to block Music service and revert Music app to classic mode when the device is in supervised mode (iOS 9.3 and later and macOS 10.12 and later).
func (m *IosGeneralDeviceConfiguration) SetITunesBlockMusicService(value *bool)() {
    err := m.GetBackingStore().Set("iTunesBlockMusicService", value)
    if err != nil {
        panic(err)
    }
}
// SetITunesBlockRadio sets the iTunesBlockRadio property value. Indicates whether or not to block the user from using iTunes Radio when the device is in supervised mode (iOS 9.3 and later).
func (m *IosGeneralDeviceConfiguration) SetITunesBlockRadio(value *bool)() {
    err := m.GetBackingStore().Set("iTunesBlockRadio", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyboardBlockAutoCorrect sets the keyboardBlockAutoCorrect property value. Indicates whether or not to block keyboard auto-correction when the device is in supervised mode (iOS 8.1.3 and later).
func (m *IosGeneralDeviceConfiguration) SetKeyboardBlockAutoCorrect(value *bool)() {
    err := m.GetBackingStore().Set("keyboardBlockAutoCorrect", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyboardBlockDictation sets the keyboardBlockDictation property value. Indicates whether or not to block the user from using dictation input when the device is in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetKeyboardBlockDictation(value *bool)() {
    err := m.GetBackingStore().Set("keyboardBlockDictation", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyboardBlockPredictive sets the keyboardBlockPredictive property value. Indicates whether or not to block predictive keyboards when device is in supervised mode (iOS 8.1.3 and later).
func (m *IosGeneralDeviceConfiguration) SetKeyboardBlockPredictive(value *bool)() {
    err := m.GetBackingStore().Set("keyboardBlockPredictive", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyboardBlockShortcuts sets the keyboardBlockShortcuts property value. Indicates whether or not to block keyboard shortcuts when the device is in supervised mode (iOS 9.0 and later).
func (m *IosGeneralDeviceConfiguration) SetKeyboardBlockShortcuts(value *bool)() {
    err := m.GetBackingStore().Set("keyboardBlockShortcuts", value)
    if err != nil {
        panic(err)
    }
}
// SetKeyboardBlockSpellCheck sets the keyboardBlockSpellCheck property value. Indicates whether or not to block keyboard spell-checking when the device is in supervised mode (iOS 8.1.3 and later).
func (m *IosGeneralDeviceConfiguration) SetKeyboardBlockSpellCheck(value *bool)() {
    err := m.GetBackingStore().Set("keyboardBlockSpellCheck", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowAssistiveSpeak sets the kioskModeAllowAssistiveSpeak property value. Indicates whether or not to allow assistive speak while in kiosk mode.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowAssistiveSpeak(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowAssistiveSpeak", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowAssistiveTouchSettings sets the kioskModeAllowAssistiveTouchSettings property value. Indicates whether or not to allow access to the Assistive Touch Settings while in kiosk mode.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowAssistiveTouchSettings(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowAssistiveTouchSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowAutoLock sets the kioskModeAllowAutoLock property value. Indicates whether or not to allow device auto lock while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockAutoLock instead.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowAutoLock(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowAutoLock", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowColorInversionSettings sets the kioskModeAllowColorInversionSettings property value. Indicates whether or not to allow access to the Color Inversion Settings while in kiosk mode.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowColorInversionSettings(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowColorInversionSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowRingerSwitch sets the kioskModeAllowRingerSwitch property value. Indicates whether or not to allow use of the ringer switch while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockRingerSwitch instead.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowRingerSwitch(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowRingerSwitch", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowScreenRotation sets the kioskModeAllowScreenRotation property value. Indicates whether or not to allow screen rotation while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockScreenRotation instead.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowScreenRotation(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowScreenRotation", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowSleepButton sets the kioskModeAllowSleepButton property value. Indicates whether or not to allow use of the sleep button while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockSleepButton instead.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowSleepButton(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowSleepButton", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowTouchscreen sets the kioskModeAllowTouchscreen property value. Indicates whether or not to allow use of the touchscreen while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockTouchscreen instead.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowTouchscreen(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowTouchscreen", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowVoiceOverSettings sets the kioskModeAllowVoiceOverSettings property value. Indicates whether or not to allow access to the voice over settings while in kiosk mode.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowVoiceOverSettings(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowVoiceOverSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowVolumeButtons sets the kioskModeAllowVolumeButtons property value. Indicates whether or not to allow use of the volume buttons while in kiosk mode. This property's functionality is redundant with the OS default and is deprecated. Use KioskModeBlockVolumeButtons instead.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowVolumeButtons(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowVolumeButtons", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAllowZoomSettings sets the kioskModeAllowZoomSettings property value. Indicates whether or not to allow access to the zoom settings while in kiosk mode.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAllowZoomSettings(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeAllowZoomSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeAppStoreUrl sets the kioskModeAppStoreUrl property value. URL in the app store to the app to use for kiosk mode. Use if KioskModeManagedAppId is not known.
func (m *IosGeneralDeviceConfiguration) SetKioskModeAppStoreUrl(value *string)() {
    err := m.GetBackingStore().Set("kioskModeAppStoreUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeBuiltInAppId sets the kioskModeBuiltInAppId property value. ID for built-in apps to use for kiosk mode. Used when KioskModeManagedAppId and KioskModeAppStoreUrl are not set.
func (m *IosGeneralDeviceConfiguration) SetKioskModeBuiltInAppId(value *string)() {
    err := m.GetBackingStore().Set("kioskModeBuiltInAppId", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeManagedAppId sets the kioskModeManagedAppId property value. Managed app id of the app to use for kiosk mode. If KioskModeManagedAppId is specified then KioskModeAppStoreUrl will be ignored.
func (m *IosGeneralDeviceConfiguration) SetKioskModeManagedAppId(value *string)() {
    err := m.GetBackingStore().Set("kioskModeManagedAppId", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeRequireAssistiveTouch sets the kioskModeRequireAssistiveTouch property value. Indicates whether or not to require assistive touch while in kiosk mode.
func (m *IosGeneralDeviceConfiguration) SetKioskModeRequireAssistiveTouch(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeRequireAssistiveTouch", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeRequireColorInversion sets the kioskModeRequireColorInversion property value. Indicates whether or not to require color inversion while in kiosk mode.
func (m *IosGeneralDeviceConfiguration) SetKioskModeRequireColorInversion(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeRequireColorInversion", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeRequireMonoAudio sets the kioskModeRequireMonoAudio property value. Indicates whether or not to require mono audio while in kiosk mode.
func (m *IosGeneralDeviceConfiguration) SetKioskModeRequireMonoAudio(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeRequireMonoAudio", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeRequireVoiceOver sets the kioskModeRequireVoiceOver property value. Indicates whether or not to require voice over while in kiosk mode.
func (m *IosGeneralDeviceConfiguration) SetKioskModeRequireVoiceOver(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeRequireVoiceOver", value)
    if err != nil {
        panic(err)
    }
}
// SetKioskModeRequireZoom sets the kioskModeRequireZoom property value. Indicates whether or not to require zoom while in kiosk mode.
func (m *IosGeneralDeviceConfiguration) SetKioskModeRequireZoom(value *bool)() {
    err := m.GetBackingStore().Set("kioskModeRequireZoom", value)
    if err != nil {
        panic(err)
    }
}
// SetLockScreenBlockControlCenter sets the lockScreenBlockControlCenter property value. Indicates whether or not to block the user from using control center on the lock screen.
func (m *IosGeneralDeviceConfiguration) SetLockScreenBlockControlCenter(value *bool)() {
    err := m.GetBackingStore().Set("lockScreenBlockControlCenter", value)
    if err != nil {
        panic(err)
    }
}
// SetLockScreenBlockNotificationView sets the lockScreenBlockNotificationView property value. Indicates whether or not to block the user from using the notification view on the lock screen.
func (m *IosGeneralDeviceConfiguration) SetLockScreenBlockNotificationView(value *bool)() {
    err := m.GetBackingStore().Set("lockScreenBlockNotificationView", value)
    if err != nil {
        panic(err)
    }
}
// SetLockScreenBlockPassbook sets the lockScreenBlockPassbook property value. Indicates whether or not to block the user from using passbook when the device is locked.
func (m *IosGeneralDeviceConfiguration) SetLockScreenBlockPassbook(value *bool)() {
    err := m.GetBackingStore().Set("lockScreenBlockPassbook", value)
    if err != nil {
        panic(err)
    }
}
// SetLockScreenBlockTodayView sets the lockScreenBlockTodayView property value. Indicates whether or not to block the user from using the Today View on the lock screen.
func (m *IosGeneralDeviceConfiguration) SetLockScreenBlockTodayView(value *bool)() {
    err := m.GetBackingStore().Set("lockScreenBlockTodayView", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaContentRatingApps sets the mediaContentRatingApps property value. Apps rating as in media content
func (m *IosGeneralDeviceConfiguration) SetMediaContentRatingApps(value *RatingAppsType)() {
    err := m.GetBackingStore().Set("mediaContentRatingApps", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaContentRatingAustralia sets the mediaContentRatingAustralia property value. Media content rating settings for Australia
func (m *IosGeneralDeviceConfiguration) SetMediaContentRatingAustralia(value MediaContentRatingAustraliaable)() {
    err := m.GetBackingStore().Set("mediaContentRatingAustralia", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaContentRatingCanada sets the mediaContentRatingCanada property value. Media content rating settings for Canada
func (m *IosGeneralDeviceConfiguration) SetMediaContentRatingCanada(value MediaContentRatingCanadaable)() {
    err := m.GetBackingStore().Set("mediaContentRatingCanada", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaContentRatingFrance sets the mediaContentRatingFrance property value. Media content rating settings for France
func (m *IosGeneralDeviceConfiguration) SetMediaContentRatingFrance(value MediaContentRatingFranceable)() {
    err := m.GetBackingStore().Set("mediaContentRatingFrance", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaContentRatingGermany sets the mediaContentRatingGermany property value. Media content rating settings for Germany
func (m *IosGeneralDeviceConfiguration) SetMediaContentRatingGermany(value MediaContentRatingGermanyable)() {
    err := m.GetBackingStore().Set("mediaContentRatingGermany", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaContentRatingIreland sets the mediaContentRatingIreland property value. Media content rating settings for Ireland
func (m *IosGeneralDeviceConfiguration) SetMediaContentRatingIreland(value MediaContentRatingIrelandable)() {
    err := m.GetBackingStore().Set("mediaContentRatingIreland", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaContentRatingJapan sets the mediaContentRatingJapan property value. Media content rating settings for Japan
func (m *IosGeneralDeviceConfiguration) SetMediaContentRatingJapan(value MediaContentRatingJapanable)() {
    err := m.GetBackingStore().Set("mediaContentRatingJapan", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaContentRatingNewZealand sets the mediaContentRatingNewZealand property value. Media content rating settings for New Zealand
func (m *IosGeneralDeviceConfiguration) SetMediaContentRatingNewZealand(value MediaContentRatingNewZealandable)() {
    err := m.GetBackingStore().Set("mediaContentRatingNewZealand", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaContentRatingUnitedKingdom sets the mediaContentRatingUnitedKingdom property value. Media content rating settings for United Kingdom
func (m *IosGeneralDeviceConfiguration) SetMediaContentRatingUnitedKingdom(value MediaContentRatingUnitedKingdomable)() {
    err := m.GetBackingStore().Set("mediaContentRatingUnitedKingdom", value)
    if err != nil {
        panic(err)
    }
}
// SetMediaContentRatingUnitedStates sets the mediaContentRatingUnitedStates property value. Media content rating settings for United States
func (m *IosGeneralDeviceConfiguration) SetMediaContentRatingUnitedStates(value MediaContentRatingUnitedStatesable)() {
    err := m.GetBackingStore().Set("mediaContentRatingUnitedStates", value)
    if err != nil {
        panic(err)
    }
}
// SetMessagesBlocked sets the messagesBlocked property value. Indicates whether or not to block the user from using the Messages app on the supervised device.
func (m *IosGeneralDeviceConfiguration) SetMessagesBlocked(value *bool)() {
    err := m.GetBackingStore().Set("messagesBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetNetworkUsageRules sets the networkUsageRules property value. List of managed apps and the network rules that applies to them. This collection can contain a maximum of 1000 elements.
func (m *IosGeneralDeviceConfiguration) SetNetworkUsageRules(value []IosNetworkUsageRuleable)() {
    err := m.GetBackingStore().Set("networkUsageRules", value)
    if err != nil {
        panic(err)
    }
}
// SetNotificationsBlockSettingsModification sets the notificationsBlockSettingsModification property value. Indicates whether or not to allow notifications settings modification (iOS 9.3 and later).
func (m *IosGeneralDeviceConfiguration) SetNotificationsBlockSettingsModification(value *bool)() {
    err := m.GetBackingStore().Set("notificationsBlockSettingsModification", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeBlockFingerprintModification sets the passcodeBlockFingerprintModification property value. Block modification of registered Touch ID fingerprints when in supervised mode.
func (m *IosGeneralDeviceConfiguration) SetPasscodeBlockFingerprintModification(value *bool)() {
    err := m.GetBackingStore().Set("passcodeBlockFingerprintModification", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeBlockFingerprintUnlock sets the passcodeBlockFingerprintUnlock property value. Indicates whether or not to block fingerprint unlock.
func (m *IosGeneralDeviceConfiguration) SetPasscodeBlockFingerprintUnlock(value *bool)() {
    err := m.GetBackingStore().Set("passcodeBlockFingerprintUnlock", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeBlockModification sets the passcodeBlockModification property value. Indicates whether or not to allow passcode modification on the supervised device (iOS 9.0 and later).
func (m *IosGeneralDeviceConfiguration) SetPasscodeBlockModification(value *bool)() {
    err := m.GetBackingStore().Set("passcodeBlockModification", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeBlockSimple sets the passcodeBlockSimple property value. Indicates whether or not to block simple passcodes.
func (m *IosGeneralDeviceConfiguration) SetPasscodeBlockSimple(value *bool)() {
    err := m.GetBackingStore().Set("passcodeBlockSimple", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeExpirationDays sets the passcodeExpirationDays property value. Number of days before the passcode expires. Valid values 1 to 65535
func (m *IosGeneralDeviceConfiguration) SetPasscodeExpirationDays(value *int32)() {
    err := m.GetBackingStore().Set("passcodeExpirationDays", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeMinimumCharacterSetCount sets the passcodeMinimumCharacterSetCount property value. Number of character sets a passcode must contain. Valid values 0 to 4
func (m *IosGeneralDeviceConfiguration) SetPasscodeMinimumCharacterSetCount(value *int32)() {
    err := m.GetBackingStore().Set("passcodeMinimumCharacterSetCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeMinimumLength sets the passcodeMinimumLength property value. Minimum length of passcode. Valid values 4 to 14
func (m *IosGeneralDeviceConfiguration) SetPasscodeMinimumLength(value *int32)() {
    err := m.GetBackingStore().Set("passcodeMinimumLength", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeMinutesOfInactivityBeforeLock sets the passcodeMinutesOfInactivityBeforeLock property value. Minutes of inactivity before a passcode is required.
func (m *IosGeneralDeviceConfiguration) SetPasscodeMinutesOfInactivityBeforeLock(value *int32)() {
    err := m.GetBackingStore().Set("passcodeMinutesOfInactivityBeforeLock", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeMinutesOfInactivityBeforeScreenTimeout sets the passcodeMinutesOfInactivityBeforeScreenTimeout property value. Minutes of inactivity before the screen times out.
func (m *IosGeneralDeviceConfiguration) SetPasscodeMinutesOfInactivityBeforeScreenTimeout(value *int32)() {
    err := m.GetBackingStore().Set("passcodeMinutesOfInactivityBeforeScreenTimeout", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodePreviousPasscodeBlockCount sets the passcodePreviousPasscodeBlockCount property value. Number of previous passcodes to block. Valid values 1 to 24
func (m *IosGeneralDeviceConfiguration) SetPasscodePreviousPasscodeBlockCount(value *int32)() {
    err := m.GetBackingStore().Set("passcodePreviousPasscodeBlockCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeRequired sets the passcodeRequired property value. Indicates whether or not to require a passcode.
func (m *IosGeneralDeviceConfiguration) SetPasscodeRequired(value *bool)() {
    err := m.GetBackingStore().Set("passcodeRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeRequiredType sets the passcodeRequiredType property value. Possible values of required passwords.
func (m *IosGeneralDeviceConfiguration) SetPasscodeRequiredType(value *RequiredPasswordType)() {
    err := m.GetBackingStore().Set("passcodeRequiredType", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscodeSignInFailureCountBeforeWipe sets the passcodeSignInFailureCountBeforeWipe property value. Number of sign in failures allowed before wiping the device. Valid values 2 to 11
func (m *IosGeneralDeviceConfiguration) SetPasscodeSignInFailureCountBeforeWipe(value *int32)() {
    err := m.GetBackingStore().Set("passcodeSignInFailureCountBeforeWipe", value)
    if err != nil {
        panic(err)
    }
}
// SetPodcastsBlocked sets the podcastsBlocked property value. Indicates whether or not to block the user from using podcasts on the supervised device (iOS 8.0 and later).
func (m *IosGeneralDeviceConfiguration) SetPodcastsBlocked(value *bool)() {
    err := m.GetBackingStore().Set("podcastsBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetSafariBlockAutofill sets the safariBlockAutofill property value. Indicates whether or not to block the user from using Auto fill in Safari. Requires a supervised device for iOS 13 and later.
func (m *IosGeneralDeviceConfiguration) SetSafariBlockAutofill(value *bool)() {
    err := m.GetBackingStore().Set("safariBlockAutofill", value)
    if err != nil {
        panic(err)
    }
}
// SetSafariBlocked sets the safariBlocked property value. Indicates whether or not to block the user from using Safari. Requires a supervised device for iOS 13 and later.
func (m *IosGeneralDeviceConfiguration) SetSafariBlocked(value *bool)() {
    err := m.GetBackingStore().Set("safariBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetSafariBlockJavaScript sets the safariBlockJavaScript property value. Indicates whether or not to block JavaScript in Safari.
func (m *IosGeneralDeviceConfiguration) SetSafariBlockJavaScript(value *bool)() {
    err := m.GetBackingStore().Set("safariBlockJavaScript", value)
    if err != nil {
        panic(err)
    }
}
// SetSafariBlockPopups sets the safariBlockPopups property value. Indicates whether or not to block popups in Safari.
func (m *IosGeneralDeviceConfiguration) SetSafariBlockPopups(value *bool)() {
    err := m.GetBackingStore().Set("safariBlockPopups", value)
    if err != nil {
        panic(err)
    }
}
// SetSafariCookieSettings sets the safariCookieSettings property value. Web Browser Cookie Settings.
func (m *IosGeneralDeviceConfiguration) SetSafariCookieSettings(value *WebBrowserCookieSettings)() {
    err := m.GetBackingStore().Set("safariCookieSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetSafariManagedDomains sets the safariManagedDomains property value. URLs matching the patterns listed here will be considered managed.
func (m *IosGeneralDeviceConfiguration) SetSafariManagedDomains(value []string)() {
    err := m.GetBackingStore().Set("safariManagedDomains", value)
    if err != nil {
        panic(err)
    }
}
// SetSafariPasswordAutoFillDomains sets the safariPasswordAutoFillDomains property value. Users can save passwords in Safari only from URLs matching the patterns listed here. Applies to devices in supervised mode (iOS 9.3 and later).
func (m *IosGeneralDeviceConfiguration) SetSafariPasswordAutoFillDomains(value []string)() {
    err := m.GetBackingStore().Set("safariPasswordAutoFillDomains", value)
    if err != nil {
        panic(err)
    }
}
// SetSafariRequireFraudWarning sets the safariRequireFraudWarning property value. Indicates whether or not to require fraud warning in Safari.
func (m *IosGeneralDeviceConfiguration) SetSafariRequireFraudWarning(value *bool)() {
    err := m.GetBackingStore().Set("safariRequireFraudWarning", value)
    if err != nil {
        panic(err)
    }
}
// SetScreenCaptureBlocked sets the screenCaptureBlocked property value. Indicates whether or not to block the user from taking Screenshots.
func (m *IosGeneralDeviceConfiguration) SetScreenCaptureBlocked(value *bool)() {
    err := m.GetBackingStore().Set("screenCaptureBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetSiriBlocked sets the siriBlocked property value. Indicates whether or not to block the user from using Siri.
func (m *IosGeneralDeviceConfiguration) SetSiriBlocked(value *bool)() {
    err := m.GetBackingStore().Set("siriBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetSiriBlockedWhenLocked sets the siriBlockedWhenLocked property value. Indicates whether or not to block the user from using Siri when locked.
func (m *IosGeneralDeviceConfiguration) SetSiriBlockedWhenLocked(value *bool)() {
    err := m.GetBackingStore().Set("siriBlockedWhenLocked", value)
    if err != nil {
        panic(err)
    }
}
// SetSiriBlockUserGeneratedContent sets the siriBlockUserGeneratedContent property value. Indicates whether or not to block Siri from querying user-generated content when used on a supervised device.
func (m *IosGeneralDeviceConfiguration) SetSiriBlockUserGeneratedContent(value *bool)() {
    err := m.GetBackingStore().Set("siriBlockUserGeneratedContent", value)
    if err != nil {
        panic(err)
    }
}
// SetSiriRequireProfanityFilter sets the siriRequireProfanityFilter property value. Indicates whether or not to prevent Siri from dictating, or speaking profane language on supervised device.
func (m *IosGeneralDeviceConfiguration) SetSiriRequireProfanityFilter(value *bool)() {
    err := m.GetBackingStore().Set("siriRequireProfanityFilter", value)
    if err != nil {
        panic(err)
    }
}
// SetSpotlightBlockInternetResults sets the spotlightBlockInternetResults property value. Indicates whether or not to block Spotlight search from returning internet results on supervised device.
func (m *IosGeneralDeviceConfiguration) SetSpotlightBlockInternetResults(value *bool)() {
    err := m.GetBackingStore().Set("spotlightBlockInternetResults", value)
    if err != nil {
        panic(err)
    }
}
// SetVoiceDialingBlocked sets the voiceDialingBlocked property value. Indicates whether or not to block voice dialing.
func (m *IosGeneralDeviceConfiguration) SetVoiceDialingBlocked(value *bool)() {
    err := m.GetBackingStore().Set("voiceDialingBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetWallpaperBlockModification sets the wallpaperBlockModification property value. Indicates whether or not to allow wallpaper modification on supervised device (iOS 9.0 and later) .
func (m *IosGeneralDeviceConfiguration) SetWallpaperBlockModification(value *bool)() {
    err := m.GetBackingStore().Set("wallpaperBlockModification", value)
    if err != nil {
        panic(err)
    }
}
// SetWiFiConnectOnlyToConfiguredNetworks sets the wiFiConnectOnlyToConfiguredNetworks property value. Indicates whether or not to force the device to use only Wi-Fi networks from configuration profiles when the device is in supervised mode. Available for devices running iOS and iPadOS versions 14.4 and earlier. Devices running 14.5+ should use the setting, 'WiFiConnectToAllowedNetworksOnlyForced.
func (m *IosGeneralDeviceConfiguration) SetWiFiConnectOnlyToConfiguredNetworks(value *bool)() {
    err := m.GetBackingStore().Set("wiFiConnectOnlyToConfiguredNetworks", value)
    if err != nil {
        panic(err)
    }
}
type IosGeneralDeviceConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccountBlockModification()(*bool)
    GetActivationLockAllowWhenSupervised()(*bool)
    GetAirDropBlocked()(*bool)
    GetAirDropForceUnmanagedDropTarget()(*bool)
    GetAirPlayForcePairingPasswordForOutgoingRequests()(*bool)
    GetAppleNewsBlocked()(*bool)
    GetAppleWatchBlockPairing()(*bool)
    GetAppleWatchForceWristDetection()(*bool)
    GetAppsSingleAppModeList()([]AppListItemable)
    GetAppStoreBlockAutomaticDownloads()(*bool)
    GetAppStoreBlocked()(*bool)
    GetAppStoreBlockInAppPurchases()(*bool)
    GetAppStoreBlockUIAppInstallation()(*bool)
    GetAppStoreRequirePassword()(*bool)
    GetAppsVisibilityList()([]AppListItemable)
    GetAppsVisibilityListType()(*AppListType)
    GetBluetoothBlockModification()(*bool)
    GetCameraBlocked()(*bool)
    GetCellularBlockDataRoaming()(*bool)
    GetCellularBlockGlobalBackgroundFetchWhileRoaming()(*bool)
    GetCellularBlockPerAppDataModification()(*bool)
    GetCellularBlockPersonalHotspot()(*bool)
    GetCellularBlockVoiceRoaming()(*bool)
    GetCertificatesBlockUntrustedTlsCertificates()(*bool)
    GetClassroomAppBlockRemoteScreenObservation()(*bool)
    GetClassroomAppForceUnpromptedScreenObservation()(*bool)
    GetCompliantAppListType()(*AppListType)
    GetCompliantAppsList()([]AppListItemable)
    GetConfigurationProfileBlockChanges()(*bool)
    GetDefinitionLookupBlocked()(*bool)
    GetDeviceBlockEnableRestrictions()(*bool)
    GetDeviceBlockEraseContentAndSettings()(*bool)
    GetDeviceBlockNameModification()(*bool)
    GetDiagnosticDataBlockSubmission()(*bool)
    GetDiagnosticDataBlockSubmissionModification()(*bool)
    GetDocumentsBlockManagedDocumentsInUnmanagedApps()(*bool)
    GetDocumentsBlockUnmanagedDocumentsInManagedApps()(*bool)
    GetEmailInDomainSuffixes()([]string)
    GetEnterpriseAppBlockTrust()(*bool)
    GetEnterpriseAppBlockTrustModification()(*bool)
    GetFaceTimeBlocked()(*bool)
    GetFindMyFriendsBlocked()(*bool)
    GetGameCenterBlocked()(*bool)
    GetGamingBlockGameCenterFriends()(*bool)
    GetGamingBlockMultiplayer()(*bool)
    GetHostPairingBlocked()(*bool)
    GetIBooksStoreBlocked()(*bool)
    GetIBooksStoreBlockErotica()(*bool)
    GetICloudBlockActivityContinuation()(*bool)
    GetICloudBlockBackup()(*bool)
    GetICloudBlockDocumentSync()(*bool)
    GetICloudBlockManagedAppsSync()(*bool)
    GetICloudBlockPhotoLibrary()(*bool)
    GetICloudBlockPhotoStreamSync()(*bool)
    GetICloudBlockSharedPhotoStream()(*bool)
    GetICloudRequireEncryptedBackup()(*bool)
    GetITunesBlockExplicitContent()(*bool)
    GetITunesBlockMusicService()(*bool)
    GetITunesBlockRadio()(*bool)
    GetKeyboardBlockAutoCorrect()(*bool)
    GetKeyboardBlockDictation()(*bool)
    GetKeyboardBlockPredictive()(*bool)
    GetKeyboardBlockShortcuts()(*bool)
    GetKeyboardBlockSpellCheck()(*bool)
    GetKioskModeAllowAssistiveSpeak()(*bool)
    GetKioskModeAllowAssistiveTouchSettings()(*bool)
    GetKioskModeAllowAutoLock()(*bool)
    GetKioskModeAllowColorInversionSettings()(*bool)
    GetKioskModeAllowRingerSwitch()(*bool)
    GetKioskModeAllowScreenRotation()(*bool)
    GetKioskModeAllowSleepButton()(*bool)
    GetKioskModeAllowTouchscreen()(*bool)
    GetKioskModeAllowVoiceOverSettings()(*bool)
    GetKioskModeAllowVolumeButtons()(*bool)
    GetKioskModeAllowZoomSettings()(*bool)
    GetKioskModeAppStoreUrl()(*string)
    GetKioskModeBuiltInAppId()(*string)
    GetKioskModeManagedAppId()(*string)
    GetKioskModeRequireAssistiveTouch()(*bool)
    GetKioskModeRequireColorInversion()(*bool)
    GetKioskModeRequireMonoAudio()(*bool)
    GetKioskModeRequireVoiceOver()(*bool)
    GetKioskModeRequireZoom()(*bool)
    GetLockScreenBlockControlCenter()(*bool)
    GetLockScreenBlockNotificationView()(*bool)
    GetLockScreenBlockPassbook()(*bool)
    GetLockScreenBlockTodayView()(*bool)
    GetMediaContentRatingApps()(*RatingAppsType)
    GetMediaContentRatingAustralia()(MediaContentRatingAustraliaable)
    GetMediaContentRatingCanada()(MediaContentRatingCanadaable)
    GetMediaContentRatingFrance()(MediaContentRatingFranceable)
    GetMediaContentRatingGermany()(MediaContentRatingGermanyable)
    GetMediaContentRatingIreland()(MediaContentRatingIrelandable)
    GetMediaContentRatingJapan()(MediaContentRatingJapanable)
    GetMediaContentRatingNewZealand()(MediaContentRatingNewZealandable)
    GetMediaContentRatingUnitedKingdom()(MediaContentRatingUnitedKingdomable)
    GetMediaContentRatingUnitedStates()(MediaContentRatingUnitedStatesable)
    GetMessagesBlocked()(*bool)
    GetNetworkUsageRules()([]IosNetworkUsageRuleable)
    GetNotificationsBlockSettingsModification()(*bool)
    GetPasscodeBlockFingerprintModification()(*bool)
    GetPasscodeBlockFingerprintUnlock()(*bool)
    GetPasscodeBlockModification()(*bool)
    GetPasscodeBlockSimple()(*bool)
    GetPasscodeExpirationDays()(*int32)
    GetPasscodeMinimumCharacterSetCount()(*int32)
    GetPasscodeMinimumLength()(*int32)
    GetPasscodeMinutesOfInactivityBeforeLock()(*int32)
    GetPasscodeMinutesOfInactivityBeforeScreenTimeout()(*int32)
    GetPasscodePreviousPasscodeBlockCount()(*int32)
    GetPasscodeRequired()(*bool)
    GetPasscodeRequiredType()(*RequiredPasswordType)
    GetPasscodeSignInFailureCountBeforeWipe()(*int32)
    GetPodcastsBlocked()(*bool)
    GetSafariBlockAutofill()(*bool)
    GetSafariBlocked()(*bool)
    GetSafariBlockJavaScript()(*bool)
    GetSafariBlockPopups()(*bool)
    GetSafariCookieSettings()(*WebBrowserCookieSettings)
    GetSafariManagedDomains()([]string)
    GetSafariPasswordAutoFillDomains()([]string)
    GetSafariRequireFraudWarning()(*bool)
    GetScreenCaptureBlocked()(*bool)
    GetSiriBlocked()(*bool)
    GetSiriBlockedWhenLocked()(*bool)
    GetSiriBlockUserGeneratedContent()(*bool)
    GetSiriRequireProfanityFilter()(*bool)
    GetSpotlightBlockInternetResults()(*bool)
    GetVoiceDialingBlocked()(*bool)
    GetWallpaperBlockModification()(*bool)
    GetWiFiConnectOnlyToConfiguredNetworks()(*bool)
    SetAccountBlockModification(value *bool)()
    SetActivationLockAllowWhenSupervised(value *bool)()
    SetAirDropBlocked(value *bool)()
    SetAirDropForceUnmanagedDropTarget(value *bool)()
    SetAirPlayForcePairingPasswordForOutgoingRequests(value *bool)()
    SetAppleNewsBlocked(value *bool)()
    SetAppleWatchBlockPairing(value *bool)()
    SetAppleWatchForceWristDetection(value *bool)()
    SetAppsSingleAppModeList(value []AppListItemable)()
    SetAppStoreBlockAutomaticDownloads(value *bool)()
    SetAppStoreBlocked(value *bool)()
    SetAppStoreBlockInAppPurchases(value *bool)()
    SetAppStoreBlockUIAppInstallation(value *bool)()
    SetAppStoreRequirePassword(value *bool)()
    SetAppsVisibilityList(value []AppListItemable)()
    SetAppsVisibilityListType(value *AppListType)()
    SetBluetoothBlockModification(value *bool)()
    SetCameraBlocked(value *bool)()
    SetCellularBlockDataRoaming(value *bool)()
    SetCellularBlockGlobalBackgroundFetchWhileRoaming(value *bool)()
    SetCellularBlockPerAppDataModification(value *bool)()
    SetCellularBlockPersonalHotspot(value *bool)()
    SetCellularBlockVoiceRoaming(value *bool)()
    SetCertificatesBlockUntrustedTlsCertificates(value *bool)()
    SetClassroomAppBlockRemoteScreenObservation(value *bool)()
    SetClassroomAppForceUnpromptedScreenObservation(value *bool)()
    SetCompliantAppListType(value *AppListType)()
    SetCompliantAppsList(value []AppListItemable)()
    SetConfigurationProfileBlockChanges(value *bool)()
    SetDefinitionLookupBlocked(value *bool)()
    SetDeviceBlockEnableRestrictions(value *bool)()
    SetDeviceBlockEraseContentAndSettings(value *bool)()
    SetDeviceBlockNameModification(value *bool)()
    SetDiagnosticDataBlockSubmission(value *bool)()
    SetDiagnosticDataBlockSubmissionModification(value *bool)()
    SetDocumentsBlockManagedDocumentsInUnmanagedApps(value *bool)()
    SetDocumentsBlockUnmanagedDocumentsInManagedApps(value *bool)()
    SetEmailInDomainSuffixes(value []string)()
    SetEnterpriseAppBlockTrust(value *bool)()
    SetEnterpriseAppBlockTrustModification(value *bool)()
    SetFaceTimeBlocked(value *bool)()
    SetFindMyFriendsBlocked(value *bool)()
    SetGameCenterBlocked(value *bool)()
    SetGamingBlockGameCenterFriends(value *bool)()
    SetGamingBlockMultiplayer(value *bool)()
    SetHostPairingBlocked(value *bool)()
    SetIBooksStoreBlocked(value *bool)()
    SetIBooksStoreBlockErotica(value *bool)()
    SetICloudBlockActivityContinuation(value *bool)()
    SetICloudBlockBackup(value *bool)()
    SetICloudBlockDocumentSync(value *bool)()
    SetICloudBlockManagedAppsSync(value *bool)()
    SetICloudBlockPhotoLibrary(value *bool)()
    SetICloudBlockPhotoStreamSync(value *bool)()
    SetICloudBlockSharedPhotoStream(value *bool)()
    SetICloudRequireEncryptedBackup(value *bool)()
    SetITunesBlockExplicitContent(value *bool)()
    SetITunesBlockMusicService(value *bool)()
    SetITunesBlockRadio(value *bool)()
    SetKeyboardBlockAutoCorrect(value *bool)()
    SetKeyboardBlockDictation(value *bool)()
    SetKeyboardBlockPredictive(value *bool)()
    SetKeyboardBlockShortcuts(value *bool)()
    SetKeyboardBlockSpellCheck(value *bool)()
    SetKioskModeAllowAssistiveSpeak(value *bool)()
    SetKioskModeAllowAssistiveTouchSettings(value *bool)()
    SetKioskModeAllowAutoLock(value *bool)()
    SetKioskModeAllowColorInversionSettings(value *bool)()
    SetKioskModeAllowRingerSwitch(value *bool)()
    SetKioskModeAllowScreenRotation(value *bool)()
    SetKioskModeAllowSleepButton(value *bool)()
    SetKioskModeAllowTouchscreen(value *bool)()
    SetKioskModeAllowVoiceOverSettings(value *bool)()
    SetKioskModeAllowVolumeButtons(value *bool)()
    SetKioskModeAllowZoomSettings(value *bool)()
    SetKioskModeAppStoreUrl(value *string)()
    SetKioskModeBuiltInAppId(value *string)()
    SetKioskModeManagedAppId(value *string)()
    SetKioskModeRequireAssistiveTouch(value *bool)()
    SetKioskModeRequireColorInversion(value *bool)()
    SetKioskModeRequireMonoAudio(value *bool)()
    SetKioskModeRequireVoiceOver(value *bool)()
    SetKioskModeRequireZoom(value *bool)()
    SetLockScreenBlockControlCenter(value *bool)()
    SetLockScreenBlockNotificationView(value *bool)()
    SetLockScreenBlockPassbook(value *bool)()
    SetLockScreenBlockTodayView(value *bool)()
    SetMediaContentRatingApps(value *RatingAppsType)()
    SetMediaContentRatingAustralia(value MediaContentRatingAustraliaable)()
    SetMediaContentRatingCanada(value MediaContentRatingCanadaable)()
    SetMediaContentRatingFrance(value MediaContentRatingFranceable)()
    SetMediaContentRatingGermany(value MediaContentRatingGermanyable)()
    SetMediaContentRatingIreland(value MediaContentRatingIrelandable)()
    SetMediaContentRatingJapan(value MediaContentRatingJapanable)()
    SetMediaContentRatingNewZealand(value MediaContentRatingNewZealandable)()
    SetMediaContentRatingUnitedKingdom(value MediaContentRatingUnitedKingdomable)()
    SetMediaContentRatingUnitedStates(value MediaContentRatingUnitedStatesable)()
    SetMessagesBlocked(value *bool)()
    SetNetworkUsageRules(value []IosNetworkUsageRuleable)()
    SetNotificationsBlockSettingsModification(value *bool)()
    SetPasscodeBlockFingerprintModification(value *bool)()
    SetPasscodeBlockFingerprintUnlock(value *bool)()
    SetPasscodeBlockModification(value *bool)()
    SetPasscodeBlockSimple(value *bool)()
    SetPasscodeExpirationDays(value *int32)()
    SetPasscodeMinimumCharacterSetCount(value *int32)()
    SetPasscodeMinimumLength(value *int32)()
    SetPasscodeMinutesOfInactivityBeforeLock(value *int32)()
    SetPasscodeMinutesOfInactivityBeforeScreenTimeout(value *int32)()
    SetPasscodePreviousPasscodeBlockCount(value *int32)()
    SetPasscodeRequired(value *bool)()
    SetPasscodeRequiredType(value *RequiredPasswordType)()
    SetPasscodeSignInFailureCountBeforeWipe(value *int32)()
    SetPodcastsBlocked(value *bool)()
    SetSafariBlockAutofill(value *bool)()
    SetSafariBlocked(value *bool)()
    SetSafariBlockJavaScript(value *bool)()
    SetSafariBlockPopups(value *bool)()
    SetSafariCookieSettings(value *WebBrowserCookieSettings)()
    SetSafariManagedDomains(value []string)()
    SetSafariPasswordAutoFillDomains(value []string)()
    SetSafariRequireFraudWarning(value *bool)()
    SetScreenCaptureBlocked(value *bool)()
    SetSiriBlocked(value *bool)()
    SetSiriBlockedWhenLocked(value *bool)()
    SetSiriBlockUserGeneratedContent(value *bool)()
    SetSiriRequireProfanityFilter(value *bool)()
    SetSpotlightBlockInternetResults(value *bool)()
    SetVoiceDialingBlocked(value *bool)()
    SetWallpaperBlockModification(value *bool)()
    SetWiFiConnectOnlyToConfiguredNetworks(value *bool)()
}
