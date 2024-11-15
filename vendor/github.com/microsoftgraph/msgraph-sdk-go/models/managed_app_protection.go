package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedAppProtection policy used to configure detailed management settings for a specified set of apps
type ManagedAppProtection struct {
    ManagedAppPolicy
}
// NewManagedAppProtection instantiates a new ManagedAppProtection and sets the default values.
func NewManagedAppProtection()(*ManagedAppProtection) {
    m := &ManagedAppProtection{
        ManagedAppPolicy: *NewManagedAppPolicy(),
    }
    odataTypeValue := "#microsoft.graph.managedAppProtection"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateManagedAppProtectionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedAppProtectionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.androidManagedAppProtection":
                        return NewAndroidManagedAppProtection(), nil
                    case "#microsoft.graph.defaultManagedAppProtection":
                        return NewDefaultManagedAppProtection(), nil
                    case "#microsoft.graph.iosManagedAppProtection":
                        return NewIosManagedAppProtection(), nil
                    case "#microsoft.graph.targetedManagedAppProtection":
                        return NewTargetedManagedAppProtection(), nil
                }
            }
        }
    }
    return NewManagedAppProtection(), nil
}
// GetAllowedDataStorageLocations gets the allowedDataStorageLocations property value. Data storage locations where a user may store managed data.
// returns a []ManagedAppDataStorageLocation when successful
func (m *ManagedAppProtection) GetAllowedDataStorageLocations()([]ManagedAppDataStorageLocation) {
    val, err := m.GetBackingStore().Get("allowedDataStorageLocations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedAppDataStorageLocation)
    }
    return nil
}
// GetAllowedInboundDataTransferSources gets the allowedInboundDataTransferSources property value. Data can be transferred from/to these classes of apps
// returns a *ManagedAppDataTransferLevel when successful
func (m *ManagedAppProtection) GetAllowedInboundDataTransferSources()(*ManagedAppDataTransferLevel) {
    val, err := m.GetBackingStore().Get("allowedInboundDataTransferSources")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ManagedAppDataTransferLevel)
    }
    return nil
}
// GetAllowedOutboundClipboardSharingLevel gets the allowedOutboundClipboardSharingLevel property value. Represents the level to which the device's clipboard may be shared between apps
// returns a *ManagedAppClipboardSharingLevel when successful
func (m *ManagedAppProtection) GetAllowedOutboundClipboardSharingLevel()(*ManagedAppClipboardSharingLevel) {
    val, err := m.GetBackingStore().Get("allowedOutboundClipboardSharingLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ManagedAppClipboardSharingLevel)
    }
    return nil
}
// GetAllowedOutboundDataTransferDestinations gets the allowedOutboundDataTransferDestinations property value. Data can be transferred from/to these classes of apps
// returns a *ManagedAppDataTransferLevel when successful
func (m *ManagedAppProtection) GetAllowedOutboundDataTransferDestinations()(*ManagedAppDataTransferLevel) {
    val, err := m.GetBackingStore().Get("allowedOutboundDataTransferDestinations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ManagedAppDataTransferLevel)
    }
    return nil
}
// GetContactSyncBlocked gets the contactSyncBlocked property value. Indicates whether contacts can be synced to the user's device.
// returns a *bool when successful
func (m *ManagedAppProtection) GetContactSyncBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("contactSyncBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDataBackupBlocked gets the dataBackupBlocked property value. Indicates whether the backup of a managed app's data is blocked.
// returns a *bool when successful
func (m *ManagedAppProtection) GetDataBackupBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("dataBackupBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDeviceComplianceRequired gets the deviceComplianceRequired property value. Indicates whether device compliance is required.
// returns a *bool when successful
func (m *ManagedAppProtection) GetDeviceComplianceRequired()(*bool) {
    val, err := m.GetBackingStore().Get("deviceComplianceRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDisableAppPinIfDevicePinIsSet gets the disableAppPinIfDevicePinIsSet property value. Indicates whether use of the app pin is required if the device pin is set.
// returns a *bool when successful
func (m *ManagedAppProtection) GetDisableAppPinIfDevicePinIsSet()(*bool) {
    val, err := m.GetBackingStore().Get("disableAppPinIfDevicePinIsSet")
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
func (m *ManagedAppProtection) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedAppPolicy.GetFieldDeserializers()
    res["allowedDataStorageLocations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfEnumValues(ParseManagedAppDataStorageLocation)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedAppDataStorageLocation, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*ManagedAppDataStorageLocation))
                }
            }
            m.SetAllowedDataStorageLocations(res)
        }
        return nil
    }
    res["allowedInboundDataTransferSources"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseManagedAppDataTransferLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedInboundDataTransferSources(val.(*ManagedAppDataTransferLevel))
        }
        return nil
    }
    res["allowedOutboundClipboardSharingLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseManagedAppClipboardSharingLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedOutboundClipboardSharingLevel(val.(*ManagedAppClipboardSharingLevel))
        }
        return nil
    }
    res["allowedOutboundDataTransferDestinations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseManagedAppDataTransferLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedOutboundDataTransferDestinations(val.(*ManagedAppDataTransferLevel))
        }
        return nil
    }
    res["contactSyncBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContactSyncBlocked(val)
        }
        return nil
    }
    res["dataBackupBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDataBackupBlocked(val)
        }
        return nil
    }
    res["deviceComplianceRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceComplianceRequired(val)
        }
        return nil
    }
    res["disableAppPinIfDevicePinIsSet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisableAppPinIfDevicePinIsSet(val)
        }
        return nil
    }
    res["fingerprintBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFingerprintBlocked(val)
        }
        return nil
    }
    res["managedBrowser"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseManagedBrowserType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagedBrowser(val.(*ManagedBrowserType))
        }
        return nil
    }
    res["managedBrowserToOpenLinksRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetManagedBrowserToOpenLinksRequired(val)
        }
        return nil
    }
    res["maximumPinRetries"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaximumPinRetries(val)
        }
        return nil
    }
    res["minimumPinLength"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumPinLength(val)
        }
        return nil
    }
    res["minimumRequiredAppVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumRequiredAppVersion(val)
        }
        return nil
    }
    res["minimumRequiredOsVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumRequiredOsVersion(val)
        }
        return nil
    }
    res["minimumWarningAppVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumWarningAppVersion(val)
        }
        return nil
    }
    res["minimumWarningOsVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumWarningOsVersion(val)
        }
        return nil
    }
    res["organizationalCredentialsRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrganizationalCredentialsRequired(val)
        }
        return nil
    }
    res["periodBeforePinReset"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPeriodBeforePinReset(val)
        }
        return nil
    }
    res["periodOfflineBeforeAccessCheck"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPeriodOfflineBeforeAccessCheck(val)
        }
        return nil
    }
    res["periodOfflineBeforeWipeIsEnforced"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPeriodOfflineBeforeWipeIsEnforced(val)
        }
        return nil
    }
    res["periodOnlineBeforeAccessCheck"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPeriodOnlineBeforeAccessCheck(val)
        }
        return nil
    }
    res["pinCharacterSet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseManagedAppPinCharacterSet)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPinCharacterSet(val.(*ManagedAppPinCharacterSet))
        }
        return nil
    }
    res["pinRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPinRequired(val)
        }
        return nil
    }
    res["printBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrintBlocked(val)
        }
        return nil
    }
    res["saveAsBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSaveAsBlocked(val)
        }
        return nil
    }
    res["simplePinBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSimplePinBlocked(val)
        }
        return nil
    }
    return res
}
// GetFingerprintBlocked gets the fingerprintBlocked property value. Indicates whether use of the fingerprint reader is allowed in place of a pin if PinRequired is set to True.
// returns a *bool when successful
func (m *ManagedAppProtection) GetFingerprintBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("fingerprintBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetManagedBrowser gets the managedBrowser property value. Type of managed browser
// returns a *ManagedBrowserType when successful
func (m *ManagedAppProtection) GetManagedBrowser()(*ManagedBrowserType) {
    val, err := m.GetBackingStore().Get("managedBrowser")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ManagedBrowserType)
    }
    return nil
}
// GetManagedBrowserToOpenLinksRequired gets the managedBrowserToOpenLinksRequired property value. Indicates whether internet links should be opened in the managed browser app, or any custom browser specified by CustomBrowserProtocol (for iOS) or CustomBrowserPackageId/CustomBrowserDisplayName (for Android)
// returns a *bool when successful
func (m *ManagedAppProtection) GetManagedBrowserToOpenLinksRequired()(*bool) {
    val, err := m.GetBackingStore().Get("managedBrowserToOpenLinksRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMaximumPinRetries gets the maximumPinRetries property value. Maximum number of incorrect pin retry attempts before the managed app is either blocked or wiped.
// returns a *int32 when successful
func (m *ManagedAppProtection) GetMaximumPinRetries()(*int32) {
    val, err := m.GetBackingStore().Get("maximumPinRetries")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMinimumPinLength gets the minimumPinLength property value. Minimum pin length required for an app-level pin if PinRequired is set to True
// returns a *int32 when successful
func (m *ManagedAppProtection) GetMinimumPinLength()(*int32) {
    val, err := m.GetBackingStore().Get("minimumPinLength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMinimumRequiredAppVersion gets the minimumRequiredAppVersion property value. Versions less than the specified version will block the managed app from accessing company data.
// returns a *string when successful
func (m *ManagedAppProtection) GetMinimumRequiredAppVersion()(*string) {
    val, err := m.GetBackingStore().Get("minimumRequiredAppVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMinimumRequiredOsVersion gets the minimumRequiredOsVersion property value. Versions less than the specified version will block the managed app from accessing company data.
// returns a *string when successful
func (m *ManagedAppProtection) GetMinimumRequiredOsVersion()(*string) {
    val, err := m.GetBackingStore().Get("minimumRequiredOsVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMinimumWarningAppVersion gets the minimumWarningAppVersion property value. Versions less than the specified version will result in warning message on the managed app.
// returns a *string when successful
func (m *ManagedAppProtection) GetMinimumWarningAppVersion()(*string) {
    val, err := m.GetBackingStore().Get("minimumWarningAppVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMinimumWarningOsVersion gets the minimumWarningOsVersion property value. Versions less than the specified version will result in warning message on the managed app from accessing company data.
// returns a *string when successful
func (m *ManagedAppProtection) GetMinimumWarningOsVersion()(*string) {
    val, err := m.GetBackingStore().Get("minimumWarningOsVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOrganizationalCredentialsRequired gets the organizationalCredentialsRequired property value. Indicates whether organizational credentials are required for app use.
// returns a *bool when successful
func (m *ManagedAppProtection) GetOrganizationalCredentialsRequired()(*bool) {
    val, err := m.GetBackingStore().Get("organizationalCredentialsRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPeriodBeforePinReset gets the periodBeforePinReset property value. TimePeriod before the all-level pin must be reset if PinRequired is set to True.
// returns a *ISODuration when successful
func (m *ManagedAppProtection) GetPeriodBeforePinReset()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("periodBeforePinReset")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetPeriodOfflineBeforeAccessCheck gets the periodOfflineBeforeAccessCheck property value. The period after which access is checked when the device is not connected to the internet.
// returns a *ISODuration when successful
func (m *ManagedAppProtection) GetPeriodOfflineBeforeAccessCheck()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("periodOfflineBeforeAccessCheck")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetPeriodOfflineBeforeWipeIsEnforced gets the periodOfflineBeforeWipeIsEnforced property value. The amount of time an app is allowed to remain disconnected from the internet before all managed data it is wiped.
// returns a *ISODuration when successful
func (m *ManagedAppProtection) GetPeriodOfflineBeforeWipeIsEnforced()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("periodOfflineBeforeWipeIsEnforced")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetPeriodOnlineBeforeAccessCheck gets the periodOnlineBeforeAccessCheck property value. The period after which access is checked when the device is connected to the internet.
// returns a *ISODuration when successful
func (m *ManagedAppProtection) GetPeriodOnlineBeforeAccessCheck()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("periodOnlineBeforeAccessCheck")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetPinCharacterSet gets the pinCharacterSet property value. Character set which is to be used for a user's app PIN
// returns a *ManagedAppPinCharacterSet when successful
func (m *ManagedAppProtection) GetPinCharacterSet()(*ManagedAppPinCharacterSet) {
    val, err := m.GetBackingStore().Get("pinCharacterSet")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ManagedAppPinCharacterSet)
    }
    return nil
}
// GetPinRequired gets the pinRequired property value. Indicates whether an app-level pin is required.
// returns a *bool when successful
func (m *ManagedAppProtection) GetPinRequired()(*bool) {
    val, err := m.GetBackingStore().Get("pinRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPrintBlocked gets the printBlocked property value. Indicates whether printing is allowed from managed apps.
// returns a *bool when successful
func (m *ManagedAppProtection) GetPrintBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("printBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSaveAsBlocked gets the saveAsBlocked property value. Indicates whether users may use the 'Save As' menu item to save a copy of protected files.
// returns a *bool when successful
func (m *ManagedAppProtection) GetSaveAsBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("saveAsBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSimplePinBlocked gets the simplePinBlocked property value. Indicates whether simplePin is blocked.
// returns a *bool when successful
func (m *ManagedAppProtection) GetSimplePinBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("simplePinBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ManagedAppProtection) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedAppPolicy.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAllowedDataStorageLocations() != nil {
        err = writer.WriteCollectionOfStringValues("allowedDataStorageLocations", SerializeManagedAppDataStorageLocation(m.GetAllowedDataStorageLocations()))
        if err != nil {
            return err
        }
    }
    if m.GetAllowedInboundDataTransferSources() != nil {
        cast := (*m.GetAllowedInboundDataTransferSources()).String()
        err = writer.WriteStringValue("allowedInboundDataTransferSources", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetAllowedOutboundClipboardSharingLevel() != nil {
        cast := (*m.GetAllowedOutboundClipboardSharingLevel()).String()
        err = writer.WriteStringValue("allowedOutboundClipboardSharingLevel", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetAllowedOutboundDataTransferDestinations() != nil {
        cast := (*m.GetAllowedOutboundDataTransferDestinations()).String()
        err = writer.WriteStringValue("allowedOutboundDataTransferDestinations", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("contactSyncBlocked", m.GetContactSyncBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("dataBackupBlocked", m.GetDataBackupBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("deviceComplianceRequired", m.GetDeviceComplianceRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("disableAppPinIfDevicePinIsSet", m.GetDisableAppPinIfDevicePinIsSet())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("fingerprintBlocked", m.GetFingerprintBlocked())
        if err != nil {
            return err
        }
    }
    if m.GetManagedBrowser() != nil {
        cast := (*m.GetManagedBrowser()).String()
        err = writer.WriteStringValue("managedBrowser", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("managedBrowserToOpenLinksRequired", m.GetManagedBrowserToOpenLinksRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("maximumPinRetries", m.GetMaximumPinRetries())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("minimumPinLength", m.GetMinimumPinLength())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("minimumRequiredAppVersion", m.GetMinimumRequiredAppVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("minimumRequiredOsVersion", m.GetMinimumRequiredOsVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("minimumWarningAppVersion", m.GetMinimumWarningAppVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("minimumWarningOsVersion", m.GetMinimumWarningOsVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("organizationalCredentialsRequired", m.GetOrganizationalCredentialsRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteISODurationValue("periodBeforePinReset", m.GetPeriodBeforePinReset())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteISODurationValue("periodOfflineBeforeAccessCheck", m.GetPeriodOfflineBeforeAccessCheck())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteISODurationValue("periodOfflineBeforeWipeIsEnforced", m.GetPeriodOfflineBeforeWipeIsEnforced())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteISODurationValue("periodOnlineBeforeAccessCheck", m.GetPeriodOnlineBeforeAccessCheck())
        if err != nil {
            return err
        }
    }
    if m.GetPinCharacterSet() != nil {
        cast := (*m.GetPinCharacterSet()).String()
        err = writer.WriteStringValue("pinCharacterSet", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("pinRequired", m.GetPinRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("printBlocked", m.GetPrintBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("saveAsBlocked", m.GetSaveAsBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("simplePinBlocked", m.GetSimplePinBlocked())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowedDataStorageLocations sets the allowedDataStorageLocations property value. Data storage locations where a user may store managed data.
func (m *ManagedAppProtection) SetAllowedDataStorageLocations(value []ManagedAppDataStorageLocation)() {
    err := m.GetBackingStore().Set("allowedDataStorageLocations", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedInboundDataTransferSources sets the allowedInboundDataTransferSources property value. Data can be transferred from/to these classes of apps
func (m *ManagedAppProtection) SetAllowedInboundDataTransferSources(value *ManagedAppDataTransferLevel)() {
    err := m.GetBackingStore().Set("allowedInboundDataTransferSources", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedOutboundClipboardSharingLevel sets the allowedOutboundClipboardSharingLevel property value. Represents the level to which the device's clipboard may be shared between apps
func (m *ManagedAppProtection) SetAllowedOutboundClipboardSharingLevel(value *ManagedAppClipboardSharingLevel)() {
    err := m.GetBackingStore().Set("allowedOutboundClipboardSharingLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedOutboundDataTransferDestinations sets the allowedOutboundDataTransferDestinations property value. Data can be transferred from/to these classes of apps
func (m *ManagedAppProtection) SetAllowedOutboundDataTransferDestinations(value *ManagedAppDataTransferLevel)() {
    err := m.GetBackingStore().Set("allowedOutboundDataTransferDestinations", value)
    if err != nil {
        panic(err)
    }
}
// SetContactSyncBlocked sets the contactSyncBlocked property value. Indicates whether contacts can be synced to the user's device.
func (m *ManagedAppProtection) SetContactSyncBlocked(value *bool)() {
    err := m.GetBackingStore().Set("contactSyncBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetDataBackupBlocked sets the dataBackupBlocked property value. Indicates whether the backup of a managed app's data is blocked.
func (m *ManagedAppProtection) SetDataBackupBlocked(value *bool)() {
    err := m.GetBackingStore().Set("dataBackupBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceComplianceRequired sets the deviceComplianceRequired property value. Indicates whether device compliance is required.
func (m *ManagedAppProtection) SetDeviceComplianceRequired(value *bool)() {
    err := m.GetBackingStore().Set("deviceComplianceRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetDisableAppPinIfDevicePinIsSet sets the disableAppPinIfDevicePinIsSet property value. Indicates whether use of the app pin is required if the device pin is set.
func (m *ManagedAppProtection) SetDisableAppPinIfDevicePinIsSet(value *bool)() {
    err := m.GetBackingStore().Set("disableAppPinIfDevicePinIsSet", value)
    if err != nil {
        panic(err)
    }
}
// SetFingerprintBlocked sets the fingerprintBlocked property value. Indicates whether use of the fingerprint reader is allowed in place of a pin if PinRequired is set to True.
func (m *ManagedAppProtection) SetFingerprintBlocked(value *bool)() {
    err := m.GetBackingStore().Set("fingerprintBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedBrowser sets the managedBrowser property value. Type of managed browser
func (m *ManagedAppProtection) SetManagedBrowser(value *ManagedBrowserType)() {
    err := m.GetBackingStore().Set("managedBrowser", value)
    if err != nil {
        panic(err)
    }
}
// SetManagedBrowserToOpenLinksRequired sets the managedBrowserToOpenLinksRequired property value. Indicates whether internet links should be opened in the managed browser app, or any custom browser specified by CustomBrowserProtocol (for iOS) or CustomBrowserPackageId/CustomBrowserDisplayName (for Android)
func (m *ManagedAppProtection) SetManagedBrowserToOpenLinksRequired(value *bool)() {
    err := m.GetBackingStore().Set("managedBrowserToOpenLinksRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetMaximumPinRetries sets the maximumPinRetries property value. Maximum number of incorrect pin retry attempts before the managed app is either blocked or wiped.
func (m *ManagedAppProtection) SetMaximumPinRetries(value *int32)() {
    err := m.GetBackingStore().Set("maximumPinRetries", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumPinLength sets the minimumPinLength property value. Minimum pin length required for an app-level pin if PinRequired is set to True
func (m *ManagedAppProtection) SetMinimumPinLength(value *int32)() {
    err := m.GetBackingStore().Set("minimumPinLength", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumRequiredAppVersion sets the minimumRequiredAppVersion property value. Versions less than the specified version will block the managed app from accessing company data.
func (m *ManagedAppProtection) SetMinimumRequiredAppVersion(value *string)() {
    err := m.GetBackingStore().Set("minimumRequiredAppVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumRequiredOsVersion sets the minimumRequiredOsVersion property value. Versions less than the specified version will block the managed app from accessing company data.
func (m *ManagedAppProtection) SetMinimumRequiredOsVersion(value *string)() {
    err := m.GetBackingStore().Set("minimumRequiredOsVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumWarningAppVersion sets the minimumWarningAppVersion property value. Versions less than the specified version will result in warning message on the managed app.
func (m *ManagedAppProtection) SetMinimumWarningAppVersion(value *string)() {
    err := m.GetBackingStore().Set("minimumWarningAppVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumWarningOsVersion sets the minimumWarningOsVersion property value. Versions less than the specified version will result in warning message on the managed app from accessing company data.
func (m *ManagedAppProtection) SetMinimumWarningOsVersion(value *string)() {
    err := m.GetBackingStore().Set("minimumWarningOsVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetOrganizationalCredentialsRequired sets the organizationalCredentialsRequired property value. Indicates whether organizational credentials are required for app use.
func (m *ManagedAppProtection) SetOrganizationalCredentialsRequired(value *bool)() {
    err := m.GetBackingStore().Set("organizationalCredentialsRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetPeriodBeforePinReset sets the periodBeforePinReset property value. TimePeriod before the all-level pin must be reset if PinRequired is set to True.
func (m *ManagedAppProtection) SetPeriodBeforePinReset(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("periodBeforePinReset", value)
    if err != nil {
        panic(err)
    }
}
// SetPeriodOfflineBeforeAccessCheck sets the periodOfflineBeforeAccessCheck property value. The period after which access is checked when the device is not connected to the internet.
func (m *ManagedAppProtection) SetPeriodOfflineBeforeAccessCheck(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("periodOfflineBeforeAccessCheck", value)
    if err != nil {
        panic(err)
    }
}
// SetPeriodOfflineBeforeWipeIsEnforced sets the periodOfflineBeforeWipeIsEnforced property value. The amount of time an app is allowed to remain disconnected from the internet before all managed data it is wiped.
func (m *ManagedAppProtection) SetPeriodOfflineBeforeWipeIsEnforced(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("periodOfflineBeforeWipeIsEnforced", value)
    if err != nil {
        panic(err)
    }
}
// SetPeriodOnlineBeforeAccessCheck sets the periodOnlineBeforeAccessCheck property value. The period after which access is checked when the device is connected to the internet.
func (m *ManagedAppProtection) SetPeriodOnlineBeforeAccessCheck(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("periodOnlineBeforeAccessCheck", value)
    if err != nil {
        panic(err)
    }
}
// SetPinCharacterSet sets the pinCharacterSet property value. Character set which is to be used for a user's app PIN
func (m *ManagedAppProtection) SetPinCharacterSet(value *ManagedAppPinCharacterSet)() {
    err := m.GetBackingStore().Set("pinCharacterSet", value)
    if err != nil {
        panic(err)
    }
}
// SetPinRequired sets the pinRequired property value. Indicates whether an app-level pin is required.
func (m *ManagedAppProtection) SetPinRequired(value *bool)() {
    err := m.GetBackingStore().Set("pinRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetPrintBlocked sets the printBlocked property value. Indicates whether printing is allowed from managed apps.
func (m *ManagedAppProtection) SetPrintBlocked(value *bool)() {
    err := m.GetBackingStore().Set("printBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetSaveAsBlocked sets the saveAsBlocked property value. Indicates whether users may use the 'Save As' menu item to save a copy of protected files.
func (m *ManagedAppProtection) SetSaveAsBlocked(value *bool)() {
    err := m.GetBackingStore().Set("saveAsBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetSimplePinBlocked sets the simplePinBlocked property value. Indicates whether simplePin is blocked.
func (m *ManagedAppProtection) SetSimplePinBlocked(value *bool)() {
    err := m.GetBackingStore().Set("simplePinBlocked", value)
    if err != nil {
        panic(err)
    }
}
type ManagedAppProtectionable interface {
    ManagedAppPolicyable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedDataStorageLocations()([]ManagedAppDataStorageLocation)
    GetAllowedInboundDataTransferSources()(*ManagedAppDataTransferLevel)
    GetAllowedOutboundClipboardSharingLevel()(*ManagedAppClipboardSharingLevel)
    GetAllowedOutboundDataTransferDestinations()(*ManagedAppDataTransferLevel)
    GetContactSyncBlocked()(*bool)
    GetDataBackupBlocked()(*bool)
    GetDeviceComplianceRequired()(*bool)
    GetDisableAppPinIfDevicePinIsSet()(*bool)
    GetFingerprintBlocked()(*bool)
    GetManagedBrowser()(*ManagedBrowserType)
    GetManagedBrowserToOpenLinksRequired()(*bool)
    GetMaximumPinRetries()(*int32)
    GetMinimumPinLength()(*int32)
    GetMinimumRequiredAppVersion()(*string)
    GetMinimumRequiredOsVersion()(*string)
    GetMinimumWarningAppVersion()(*string)
    GetMinimumWarningOsVersion()(*string)
    GetOrganizationalCredentialsRequired()(*bool)
    GetPeriodBeforePinReset()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetPeriodOfflineBeforeAccessCheck()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetPeriodOfflineBeforeWipeIsEnforced()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetPeriodOnlineBeforeAccessCheck()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetPinCharacterSet()(*ManagedAppPinCharacterSet)
    GetPinRequired()(*bool)
    GetPrintBlocked()(*bool)
    GetSaveAsBlocked()(*bool)
    GetSimplePinBlocked()(*bool)
    SetAllowedDataStorageLocations(value []ManagedAppDataStorageLocation)()
    SetAllowedInboundDataTransferSources(value *ManagedAppDataTransferLevel)()
    SetAllowedOutboundClipboardSharingLevel(value *ManagedAppClipboardSharingLevel)()
    SetAllowedOutboundDataTransferDestinations(value *ManagedAppDataTransferLevel)()
    SetContactSyncBlocked(value *bool)()
    SetDataBackupBlocked(value *bool)()
    SetDeviceComplianceRequired(value *bool)()
    SetDisableAppPinIfDevicePinIsSet(value *bool)()
    SetFingerprintBlocked(value *bool)()
    SetManagedBrowser(value *ManagedBrowserType)()
    SetManagedBrowserToOpenLinksRequired(value *bool)()
    SetMaximumPinRetries(value *int32)()
    SetMinimumPinLength(value *int32)()
    SetMinimumRequiredAppVersion(value *string)()
    SetMinimumRequiredOsVersion(value *string)()
    SetMinimumWarningAppVersion(value *string)()
    SetMinimumWarningOsVersion(value *string)()
    SetOrganizationalCredentialsRequired(value *bool)()
    SetPeriodBeforePinReset(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetPeriodOfflineBeforeAccessCheck(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetPeriodOfflineBeforeWipeIsEnforced(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetPeriodOnlineBeforeAccessCheck(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetPinCharacterSet(value *ManagedAppPinCharacterSet)()
    SetPinRequired(value *bool)()
    SetPrintBlocked(value *bool)()
    SetSaveAsBlocked(value *bool)()
    SetSimplePinBlocked(value *bool)()
}
