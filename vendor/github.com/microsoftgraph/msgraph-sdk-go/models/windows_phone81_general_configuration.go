package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsPhone81GeneralConfiguration this topic provides descriptions of the declared methods, properties and relationships exposed by the windowsPhone81GeneralConfiguration resource.
type WindowsPhone81GeneralConfiguration struct {
    DeviceConfiguration
}
// NewWindowsPhone81GeneralConfiguration instantiates a new WindowsPhone81GeneralConfiguration and sets the default values.
func NewWindowsPhone81GeneralConfiguration()(*WindowsPhone81GeneralConfiguration) {
    m := &WindowsPhone81GeneralConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.windowsPhone81GeneralConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsPhone81GeneralConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsPhone81GeneralConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsPhone81GeneralConfiguration(), nil
}
// GetApplyOnlyToWindowsPhone81 gets the applyOnlyToWindowsPhone81 property value. Value indicating whether this policy only applies to Windows Phone 8.1. This property is read-only.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetApplyOnlyToWindowsPhone81()(*bool) {
    val, err := m.GetBackingStore().Get("applyOnlyToWindowsPhone81")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAppsBlockCopyPaste gets the appsBlockCopyPaste property value. Indicates whether or not to block copy paste.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetAppsBlockCopyPaste()(*bool) {
    val, err := m.GetBackingStore().Get("appsBlockCopyPaste")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBluetoothBlocked gets the bluetoothBlocked property value. Indicates whether or not to block bluetooth.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetBluetoothBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("bluetoothBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCameraBlocked gets the cameraBlocked property value. Indicates whether or not to block camera.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetCameraBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("cameraBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetCellularBlockWifiTethering gets the cellularBlockWifiTethering property value. Indicates whether or not to block Wi-Fi tethering. Has no impact if Wi-Fi is blocked.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetCellularBlockWifiTethering()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockWifiTethering")
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
func (m *WindowsPhone81GeneralConfiguration) GetCompliantAppListType()(*AppListType) {
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
func (m *WindowsPhone81GeneralConfiguration) GetCompliantAppsList()([]AppListItemable) {
    val, err := m.GetBackingStore().Get("compliantAppsList")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AppListItemable)
    }
    return nil
}
// GetDiagnosticDataBlockSubmission gets the diagnosticDataBlockSubmission property value. Indicates whether or not to block diagnostic data submission.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetDiagnosticDataBlockSubmission()(*bool) {
    val, err := m.GetBackingStore().Get("diagnosticDataBlockSubmission")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEmailBlockAddingAccounts gets the emailBlockAddingAccounts property value. Indicates whether or not to block custom email accounts.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetEmailBlockAddingAccounts()(*bool) {
    val, err := m.GetBackingStore().Get("emailBlockAddingAccounts")
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
func (m *WindowsPhone81GeneralConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["applyOnlyToWindowsPhone81"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplyOnlyToWindowsPhone81(val)
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
    res["cellularBlockWifiTethering"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCellularBlockWifiTethering(val)
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
    res["emailBlockAddingAccounts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmailBlockAddingAccounts(val)
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
    res["storageRequireEncryption"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStorageRequireEncryption(val)
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
    res["wifiBlockAutomaticConnectHotspots"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiBlockAutomaticConnectHotspots(val)
        }
        return nil
    }
    res["wifiBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiBlocked(val)
        }
        return nil
    }
    res["wifiBlockHotspotReporting"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWifiBlockHotspotReporting(val)
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
    return res
}
// GetLocationServicesBlocked gets the locationServicesBlocked property value. Indicates whether or not to block location services.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetLocationServicesBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("locationServicesBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMicrosoftAccountBlocked gets the microsoftAccountBlocked property value. Indicates whether or not to block using a Microsoft Account.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetMicrosoftAccountBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("microsoftAccountBlocked")
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
func (m *WindowsPhone81GeneralConfiguration) GetNfcBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("nfcBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordBlockSimple gets the passwordBlockSimple property value. Indicates whether or not to block syncing the calendar.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetPasswordBlockSimple()(*bool) {
    val, err := m.GetBackingStore().Get("passwordBlockSimple")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordExpirationDays gets the passwordExpirationDays property value. Number of days before the password expires.
// returns a *int32 when successful
func (m *WindowsPhone81GeneralConfiguration) GetPasswordExpirationDays()(*int32) {
    val, err := m.GetBackingStore().Get("passwordExpirationDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinimumCharacterSetCount gets the passwordMinimumCharacterSetCount property value. Number of character sets a password must contain.
// returns a *int32 when successful
func (m *WindowsPhone81GeneralConfiguration) GetPasswordMinimumCharacterSetCount()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinimumCharacterSetCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinimumLength gets the passwordMinimumLength property value. Minimum length of passwords.
// returns a *int32 when successful
func (m *WindowsPhone81GeneralConfiguration) GetPasswordMinimumLength()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinimumLength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinutesOfInactivityBeforeScreenTimeout gets the passwordMinutesOfInactivityBeforeScreenTimeout property value. Minutes of inactivity before screen timeout.
// returns a *int32 when successful
func (m *WindowsPhone81GeneralConfiguration) GetPasswordMinutesOfInactivityBeforeScreenTimeout()(*int32) {
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
func (m *WindowsPhone81GeneralConfiguration) GetPasswordPreviousPasswordBlockCount()(*int32) {
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
func (m *WindowsPhone81GeneralConfiguration) GetPasswordRequired()(*bool) {
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
func (m *WindowsPhone81GeneralConfiguration) GetPasswordRequiredType()(*RequiredPasswordType) {
    val, err := m.GetBackingStore().Get("passwordRequiredType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RequiredPasswordType)
    }
    return nil
}
// GetPasswordSignInFailureCountBeforeFactoryReset gets the passwordSignInFailureCountBeforeFactoryReset property value. Number of sign in failures allowed before factory reset.
// returns a *int32 when successful
func (m *WindowsPhone81GeneralConfiguration) GetPasswordSignInFailureCountBeforeFactoryReset()(*int32) {
    val, err := m.GetBackingStore().Get("passwordSignInFailureCountBeforeFactoryReset")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetScreenCaptureBlocked gets the screenCaptureBlocked property value. Indicates whether or not to block screenshots.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetScreenCaptureBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("screenCaptureBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStorageBlockRemovableStorage gets the storageBlockRemovableStorage property value. Indicates whether or not to block removable storage.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetStorageBlockRemovableStorage()(*bool) {
    val, err := m.GetBackingStore().Get("storageBlockRemovableStorage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStorageRequireEncryption gets the storageRequireEncryption property value. Indicates whether or not to require encryption.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetStorageRequireEncryption()(*bool) {
    val, err := m.GetBackingStore().Get("storageRequireEncryption")
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
func (m *WindowsPhone81GeneralConfiguration) GetWebBrowserBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("webBrowserBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWifiBlockAutomaticConnectHotspots gets the wifiBlockAutomaticConnectHotspots property value. Indicates whether or not to block automatically connecting to Wi-Fi hotspots. Has no impact if Wi-Fi is blocked.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetWifiBlockAutomaticConnectHotspots()(*bool) {
    val, err := m.GetBackingStore().Get("wifiBlockAutomaticConnectHotspots")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWifiBlocked gets the wifiBlocked property value. Indicates whether or not to block Wi-Fi.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetWifiBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("wifiBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWifiBlockHotspotReporting gets the wifiBlockHotspotReporting property value. Indicates whether or not to block Wi-Fi hotspot reporting. Has no impact if Wi-Fi is blocked.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetWifiBlockHotspotReporting()(*bool) {
    val, err := m.GetBackingStore().Get("wifiBlockHotspotReporting")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsStoreBlocked gets the windowsStoreBlocked property value. Indicates whether or not to block the Windows Store.
// returns a *bool when successful
func (m *WindowsPhone81GeneralConfiguration) GetWindowsStoreBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("windowsStoreBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsPhone81GeneralConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("appsBlockCopyPaste", m.GetAppsBlockCopyPaste())
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
        err = writer.WriteBoolValue("cellularBlockWifiTethering", m.GetCellularBlockWifiTethering())
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
        err = writer.WriteBoolValue("diagnosticDataBlockSubmission", m.GetDiagnosticDataBlockSubmission())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("emailBlockAddingAccounts", m.GetEmailBlockAddingAccounts())
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
        err = writer.WriteBoolValue("microsoftAccountBlocked", m.GetMicrosoftAccountBlocked())
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
        err = writer.WriteInt32Value("passwordSignInFailureCountBeforeFactoryReset", m.GetPasswordSignInFailureCountBeforeFactoryReset())
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
        err = writer.WriteBoolValue("storageBlockRemovableStorage", m.GetStorageBlockRemovableStorage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("storageRequireEncryption", m.GetStorageRequireEncryption())
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
        err = writer.WriteBoolValue("wifiBlockAutomaticConnectHotspots", m.GetWifiBlockAutomaticConnectHotspots())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wifiBlocked", m.GetWifiBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wifiBlockHotspotReporting", m.GetWifiBlockHotspotReporting())
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
    return nil
}
// SetApplyOnlyToWindowsPhone81 sets the applyOnlyToWindowsPhone81 property value. Value indicating whether this policy only applies to Windows Phone 8.1. This property is read-only.
func (m *WindowsPhone81GeneralConfiguration) SetApplyOnlyToWindowsPhone81(value *bool)() {
    err := m.GetBackingStore().Set("applyOnlyToWindowsPhone81", value)
    if err != nil {
        panic(err)
    }
}
// SetAppsBlockCopyPaste sets the appsBlockCopyPaste property value. Indicates whether or not to block copy paste.
func (m *WindowsPhone81GeneralConfiguration) SetAppsBlockCopyPaste(value *bool)() {
    err := m.GetBackingStore().Set("appsBlockCopyPaste", value)
    if err != nil {
        panic(err)
    }
}
// SetBluetoothBlocked sets the bluetoothBlocked property value. Indicates whether or not to block bluetooth.
func (m *WindowsPhone81GeneralConfiguration) SetBluetoothBlocked(value *bool)() {
    err := m.GetBackingStore().Set("bluetoothBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetCameraBlocked sets the cameraBlocked property value. Indicates whether or not to block camera.
func (m *WindowsPhone81GeneralConfiguration) SetCameraBlocked(value *bool)() {
    err := m.GetBackingStore().Set("cameraBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockWifiTethering sets the cellularBlockWifiTethering property value. Indicates whether or not to block Wi-Fi tethering. Has no impact if Wi-Fi is blocked.
func (m *WindowsPhone81GeneralConfiguration) SetCellularBlockWifiTethering(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockWifiTethering", value)
    if err != nil {
        panic(err)
    }
}
// SetCompliantAppListType sets the compliantAppListType property value. Possible values of the compliance app list.
func (m *WindowsPhone81GeneralConfiguration) SetCompliantAppListType(value *AppListType)() {
    err := m.GetBackingStore().Set("compliantAppListType", value)
    if err != nil {
        panic(err)
    }
}
// SetCompliantAppsList sets the compliantAppsList property value. List of apps in the compliance (either allow list or block list, controlled by CompliantAppListType). This collection can contain a maximum of 10000 elements.
func (m *WindowsPhone81GeneralConfiguration) SetCompliantAppsList(value []AppListItemable)() {
    err := m.GetBackingStore().Set("compliantAppsList", value)
    if err != nil {
        panic(err)
    }
}
// SetDiagnosticDataBlockSubmission sets the diagnosticDataBlockSubmission property value. Indicates whether or not to block diagnostic data submission.
func (m *WindowsPhone81GeneralConfiguration) SetDiagnosticDataBlockSubmission(value *bool)() {
    err := m.GetBackingStore().Set("diagnosticDataBlockSubmission", value)
    if err != nil {
        panic(err)
    }
}
// SetEmailBlockAddingAccounts sets the emailBlockAddingAccounts property value. Indicates whether or not to block custom email accounts.
func (m *WindowsPhone81GeneralConfiguration) SetEmailBlockAddingAccounts(value *bool)() {
    err := m.GetBackingStore().Set("emailBlockAddingAccounts", value)
    if err != nil {
        panic(err)
    }
}
// SetLocationServicesBlocked sets the locationServicesBlocked property value. Indicates whether or not to block location services.
func (m *WindowsPhone81GeneralConfiguration) SetLocationServicesBlocked(value *bool)() {
    err := m.GetBackingStore().Set("locationServicesBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetMicrosoftAccountBlocked sets the microsoftAccountBlocked property value. Indicates whether or not to block using a Microsoft Account.
func (m *WindowsPhone81GeneralConfiguration) SetMicrosoftAccountBlocked(value *bool)() {
    err := m.GetBackingStore().Set("microsoftAccountBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetNfcBlocked sets the nfcBlocked property value. Indicates whether or not to block Near-Field Communication.
func (m *WindowsPhone81GeneralConfiguration) SetNfcBlocked(value *bool)() {
    err := m.GetBackingStore().Set("nfcBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordBlockSimple sets the passwordBlockSimple property value. Indicates whether or not to block syncing the calendar.
func (m *WindowsPhone81GeneralConfiguration) SetPasswordBlockSimple(value *bool)() {
    err := m.GetBackingStore().Set("passwordBlockSimple", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordExpirationDays sets the passwordExpirationDays property value. Number of days before the password expires.
func (m *WindowsPhone81GeneralConfiguration) SetPasswordExpirationDays(value *int32)() {
    err := m.GetBackingStore().Set("passwordExpirationDays", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinimumCharacterSetCount sets the passwordMinimumCharacterSetCount property value. Number of character sets a password must contain.
func (m *WindowsPhone81GeneralConfiguration) SetPasswordMinimumCharacterSetCount(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinimumCharacterSetCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinimumLength sets the passwordMinimumLength property value. Minimum length of passwords.
func (m *WindowsPhone81GeneralConfiguration) SetPasswordMinimumLength(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinimumLength", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinutesOfInactivityBeforeScreenTimeout sets the passwordMinutesOfInactivityBeforeScreenTimeout property value. Minutes of inactivity before screen timeout.
func (m *WindowsPhone81GeneralConfiguration) SetPasswordMinutesOfInactivityBeforeScreenTimeout(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinutesOfInactivityBeforeScreenTimeout", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordPreviousPasswordBlockCount sets the passwordPreviousPasswordBlockCount property value. Number of previous passwords to block. Valid values 0 to 24
func (m *WindowsPhone81GeneralConfiguration) SetPasswordPreviousPasswordBlockCount(value *int32)() {
    err := m.GetBackingStore().Set("passwordPreviousPasswordBlockCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordRequired sets the passwordRequired property value. Indicates whether or not to require a password.
func (m *WindowsPhone81GeneralConfiguration) SetPasswordRequired(value *bool)() {
    err := m.GetBackingStore().Set("passwordRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordRequiredType sets the passwordRequiredType property value. Possible values of required passwords.
func (m *WindowsPhone81GeneralConfiguration) SetPasswordRequiredType(value *RequiredPasswordType)() {
    err := m.GetBackingStore().Set("passwordRequiredType", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordSignInFailureCountBeforeFactoryReset sets the passwordSignInFailureCountBeforeFactoryReset property value. Number of sign in failures allowed before factory reset.
func (m *WindowsPhone81GeneralConfiguration) SetPasswordSignInFailureCountBeforeFactoryReset(value *int32)() {
    err := m.GetBackingStore().Set("passwordSignInFailureCountBeforeFactoryReset", value)
    if err != nil {
        panic(err)
    }
}
// SetScreenCaptureBlocked sets the screenCaptureBlocked property value. Indicates whether or not to block screenshots.
func (m *WindowsPhone81GeneralConfiguration) SetScreenCaptureBlocked(value *bool)() {
    err := m.GetBackingStore().Set("screenCaptureBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageBlockRemovableStorage sets the storageBlockRemovableStorage property value. Indicates whether or not to block removable storage.
func (m *WindowsPhone81GeneralConfiguration) SetStorageBlockRemovableStorage(value *bool)() {
    err := m.GetBackingStore().Set("storageBlockRemovableStorage", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageRequireEncryption sets the storageRequireEncryption property value. Indicates whether or not to require encryption.
func (m *WindowsPhone81GeneralConfiguration) SetStorageRequireEncryption(value *bool)() {
    err := m.GetBackingStore().Set("storageRequireEncryption", value)
    if err != nil {
        panic(err)
    }
}
// SetWebBrowserBlocked sets the webBrowserBlocked property value. Indicates whether or not to block the web browser.
func (m *WindowsPhone81GeneralConfiguration) SetWebBrowserBlocked(value *bool)() {
    err := m.GetBackingStore().Set("webBrowserBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiBlockAutomaticConnectHotspots sets the wifiBlockAutomaticConnectHotspots property value. Indicates whether or not to block automatically connecting to Wi-Fi hotspots. Has no impact if Wi-Fi is blocked.
func (m *WindowsPhone81GeneralConfiguration) SetWifiBlockAutomaticConnectHotspots(value *bool)() {
    err := m.GetBackingStore().Set("wifiBlockAutomaticConnectHotspots", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiBlocked sets the wifiBlocked property value. Indicates whether or not to block Wi-Fi.
func (m *WindowsPhone81GeneralConfiguration) SetWifiBlocked(value *bool)() {
    err := m.GetBackingStore().Set("wifiBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetWifiBlockHotspotReporting sets the wifiBlockHotspotReporting property value. Indicates whether or not to block Wi-Fi hotspot reporting. Has no impact if Wi-Fi is blocked.
func (m *WindowsPhone81GeneralConfiguration) SetWifiBlockHotspotReporting(value *bool)() {
    err := m.GetBackingStore().Set("wifiBlockHotspotReporting", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsStoreBlocked sets the windowsStoreBlocked property value. Indicates whether or not to block the Windows Store.
func (m *WindowsPhone81GeneralConfiguration) SetWindowsStoreBlocked(value *bool)() {
    err := m.GetBackingStore().Set("windowsStoreBlocked", value)
    if err != nil {
        panic(err)
    }
}
type WindowsPhone81GeneralConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplyOnlyToWindowsPhone81()(*bool)
    GetAppsBlockCopyPaste()(*bool)
    GetBluetoothBlocked()(*bool)
    GetCameraBlocked()(*bool)
    GetCellularBlockWifiTethering()(*bool)
    GetCompliantAppListType()(*AppListType)
    GetCompliantAppsList()([]AppListItemable)
    GetDiagnosticDataBlockSubmission()(*bool)
    GetEmailBlockAddingAccounts()(*bool)
    GetLocationServicesBlocked()(*bool)
    GetMicrosoftAccountBlocked()(*bool)
    GetNfcBlocked()(*bool)
    GetPasswordBlockSimple()(*bool)
    GetPasswordExpirationDays()(*int32)
    GetPasswordMinimumCharacterSetCount()(*int32)
    GetPasswordMinimumLength()(*int32)
    GetPasswordMinutesOfInactivityBeforeScreenTimeout()(*int32)
    GetPasswordPreviousPasswordBlockCount()(*int32)
    GetPasswordRequired()(*bool)
    GetPasswordRequiredType()(*RequiredPasswordType)
    GetPasswordSignInFailureCountBeforeFactoryReset()(*int32)
    GetScreenCaptureBlocked()(*bool)
    GetStorageBlockRemovableStorage()(*bool)
    GetStorageRequireEncryption()(*bool)
    GetWebBrowserBlocked()(*bool)
    GetWifiBlockAutomaticConnectHotspots()(*bool)
    GetWifiBlocked()(*bool)
    GetWifiBlockHotspotReporting()(*bool)
    GetWindowsStoreBlocked()(*bool)
    SetApplyOnlyToWindowsPhone81(value *bool)()
    SetAppsBlockCopyPaste(value *bool)()
    SetBluetoothBlocked(value *bool)()
    SetCameraBlocked(value *bool)()
    SetCellularBlockWifiTethering(value *bool)()
    SetCompliantAppListType(value *AppListType)()
    SetCompliantAppsList(value []AppListItemable)()
    SetDiagnosticDataBlockSubmission(value *bool)()
    SetEmailBlockAddingAccounts(value *bool)()
    SetLocationServicesBlocked(value *bool)()
    SetMicrosoftAccountBlocked(value *bool)()
    SetNfcBlocked(value *bool)()
    SetPasswordBlockSimple(value *bool)()
    SetPasswordExpirationDays(value *int32)()
    SetPasswordMinimumCharacterSetCount(value *int32)()
    SetPasswordMinimumLength(value *int32)()
    SetPasswordMinutesOfInactivityBeforeScreenTimeout(value *int32)()
    SetPasswordPreviousPasswordBlockCount(value *int32)()
    SetPasswordRequired(value *bool)()
    SetPasswordRequiredType(value *RequiredPasswordType)()
    SetPasswordSignInFailureCountBeforeFactoryReset(value *int32)()
    SetScreenCaptureBlocked(value *bool)()
    SetStorageBlockRemovableStorage(value *bool)()
    SetStorageRequireEncryption(value *bool)()
    SetWebBrowserBlocked(value *bool)()
    SetWifiBlockAutomaticConnectHotspots(value *bool)()
    SetWifiBlocked(value *bool)()
    SetWifiBlockHotspotReporting(value *bool)()
    SetWindowsStoreBlocked(value *bool)()
}
