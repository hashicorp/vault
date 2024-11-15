package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Windows81GeneralConfiguration this topic provides descriptions of the declared methods, properties and relationships exposed by the windows81GeneralConfiguration resource.
type Windows81GeneralConfiguration struct {
    DeviceConfiguration
}
// NewWindows81GeneralConfiguration instantiates a new Windows81GeneralConfiguration and sets the default values.
func NewWindows81GeneralConfiguration()(*Windows81GeneralConfiguration) {
    m := &Windows81GeneralConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.windows81GeneralConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindows81GeneralConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindows81GeneralConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindows81GeneralConfiguration(), nil
}
// GetAccountsBlockAddingNonMicrosoftAccountEmail gets the accountsBlockAddingNonMicrosoftAccountEmail property value. Indicates whether or not to Block the user from adding email accounts to the device that are not associated with a Microsoft account.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetAccountsBlockAddingNonMicrosoftAccountEmail()(*bool) {
    val, err := m.GetBackingStore().Get("accountsBlockAddingNonMicrosoftAccountEmail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetApplyOnlyToWindows81 gets the applyOnlyToWindows81 property value. Value indicating whether this policy only applies to Windows 8.1. This property is read-only.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetApplyOnlyToWindows81()(*bool) {
    val, err := m.GetBackingStore().Get("applyOnlyToWindows81")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserBlockAutofill gets the browserBlockAutofill property value. Indicates whether or not to block auto fill.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserBlockAutofill()(*bool) {
    val, err := m.GetBackingStore().Get("browserBlockAutofill")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserBlockAutomaticDetectionOfIntranetSites gets the browserBlockAutomaticDetectionOfIntranetSites property value. Indicates whether or not to block automatic detection of Intranet sites.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserBlockAutomaticDetectionOfIntranetSites()(*bool) {
    val, err := m.GetBackingStore().Get("browserBlockAutomaticDetectionOfIntranetSites")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserBlockEnterpriseModeAccess gets the browserBlockEnterpriseModeAccess property value. Indicates whether or not to block enterprise mode access.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserBlockEnterpriseModeAccess()(*bool) {
    val, err := m.GetBackingStore().Get("browserBlockEnterpriseModeAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserBlockJavaScript gets the browserBlockJavaScript property value. Indicates whether or not to Block the user from using JavaScript.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserBlockJavaScript()(*bool) {
    val, err := m.GetBackingStore().Get("browserBlockJavaScript")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserBlockPlugins gets the browserBlockPlugins property value. Indicates whether or not to block plug-ins.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserBlockPlugins()(*bool) {
    val, err := m.GetBackingStore().Get("browserBlockPlugins")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserBlockPopups gets the browserBlockPopups property value. Indicates whether or not to block popups.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserBlockPopups()(*bool) {
    val, err := m.GetBackingStore().Get("browserBlockPopups")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserBlockSendingDoNotTrackHeader gets the browserBlockSendingDoNotTrackHeader property value. Indicates whether or not to Block the user from sending the do not track header.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserBlockSendingDoNotTrackHeader()(*bool) {
    val, err := m.GetBackingStore().Get("browserBlockSendingDoNotTrackHeader")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserBlockSingleWordEntryOnIntranetSites gets the browserBlockSingleWordEntryOnIntranetSites property value. Indicates whether or not to block a single word entry on Intranet sites.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserBlockSingleWordEntryOnIntranetSites()(*bool) {
    val, err := m.GetBackingStore().Get("browserBlockSingleWordEntryOnIntranetSites")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserEnterpriseModeSiteListLocation gets the browserEnterpriseModeSiteListLocation property value. The enterprise mode site list location. Could be a local file, local network or http location.
// returns a *string when successful
func (m *Windows81GeneralConfiguration) GetBrowserEnterpriseModeSiteListLocation()(*string) {
    val, err := m.GetBackingStore().Get("browserEnterpriseModeSiteListLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBrowserInternetSecurityLevel gets the browserInternetSecurityLevel property value. Possible values for internet site security level.
// returns a *InternetSiteSecurityLevel when successful
func (m *Windows81GeneralConfiguration) GetBrowserInternetSecurityLevel()(*InternetSiteSecurityLevel) {
    val, err := m.GetBackingStore().Get("browserInternetSecurityLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*InternetSiteSecurityLevel)
    }
    return nil
}
// GetBrowserIntranetSecurityLevel gets the browserIntranetSecurityLevel property value. Possible values for site security level.
// returns a *SiteSecurityLevel when successful
func (m *Windows81GeneralConfiguration) GetBrowserIntranetSecurityLevel()(*SiteSecurityLevel) {
    val, err := m.GetBackingStore().Get("browserIntranetSecurityLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SiteSecurityLevel)
    }
    return nil
}
// GetBrowserLoggingReportLocation gets the browserLoggingReportLocation property value. The logging report location.
// returns a *string when successful
func (m *Windows81GeneralConfiguration) GetBrowserLoggingReportLocation()(*string) {
    val, err := m.GetBackingStore().Get("browserLoggingReportLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBrowserRequireFirewall gets the browserRequireFirewall property value. Indicates whether or not to require a firewall.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserRequireFirewall()(*bool) {
    val, err := m.GetBackingStore().Get("browserRequireFirewall")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserRequireFraudWarning gets the browserRequireFraudWarning property value. Indicates whether or not to require fraud warning.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserRequireFraudWarning()(*bool) {
    val, err := m.GetBackingStore().Get("browserRequireFraudWarning")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserRequireHighSecurityForRestrictedSites gets the browserRequireHighSecurityForRestrictedSites property value. Indicates whether or not to require high security for restricted sites.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserRequireHighSecurityForRestrictedSites()(*bool) {
    val, err := m.GetBackingStore().Get("browserRequireHighSecurityForRestrictedSites")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserRequireSmartScreen gets the browserRequireSmartScreen property value. Indicates whether or not to require the user to use the smart screen filter.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetBrowserRequireSmartScreen()(*bool) {
    val, err := m.GetBackingStore().Get("browserRequireSmartScreen")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBrowserTrustedSitesSecurityLevel gets the browserTrustedSitesSecurityLevel property value. Possible values for site security level.
// returns a *SiteSecurityLevel when successful
func (m *Windows81GeneralConfiguration) GetBrowserTrustedSitesSecurityLevel()(*SiteSecurityLevel) {
    val, err := m.GetBackingStore().Get("browserTrustedSitesSecurityLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SiteSecurityLevel)
    }
    return nil
}
// GetCellularBlockDataRoaming gets the cellularBlockDataRoaming property value. Indicates whether or not to block data roaming.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetCellularBlockDataRoaming()(*bool) {
    val, err := m.GetBackingStore().Get("cellularBlockDataRoaming")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDiagnosticsBlockDataSubmission gets the diagnosticsBlockDataSubmission property value. Indicates whether or not to block diagnostic data submission.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetDiagnosticsBlockDataSubmission()(*bool) {
    val, err := m.GetBackingStore().Get("diagnosticsBlockDataSubmission")
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
func (m *Windows81GeneralConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["applyOnlyToWindows81"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplyOnlyToWindows81(val)
        }
        return nil
    }
    res["browserBlockAutofill"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserBlockAutofill(val)
        }
        return nil
    }
    res["browserBlockAutomaticDetectionOfIntranetSites"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserBlockAutomaticDetectionOfIntranetSites(val)
        }
        return nil
    }
    res["browserBlockEnterpriseModeAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserBlockEnterpriseModeAccess(val)
        }
        return nil
    }
    res["browserBlockJavaScript"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserBlockJavaScript(val)
        }
        return nil
    }
    res["browserBlockPlugins"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserBlockPlugins(val)
        }
        return nil
    }
    res["browserBlockPopups"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserBlockPopups(val)
        }
        return nil
    }
    res["browserBlockSendingDoNotTrackHeader"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserBlockSendingDoNotTrackHeader(val)
        }
        return nil
    }
    res["browserBlockSingleWordEntryOnIntranetSites"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserBlockSingleWordEntryOnIntranetSites(val)
        }
        return nil
    }
    res["browserEnterpriseModeSiteListLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserEnterpriseModeSiteListLocation(val)
        }
        return nil
    }
    res["browserInternetSecurityLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseInternetSiteSecurityLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserInternetSecurityLevel(val.(*InternetSiteSecurityLevel))
        }
        return nil
    }
    res["browserIntranetSecurityLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSiteSecurityLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserIntranetSecurityLevel(val.(*SiteSecurityLevel))
        }
        return nil
    }
    res["browserLoggingReportLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserLoggingReportLocation(val)
        }
        return nil
    }
    res["browserRequireFirewall"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserRequireFirewall(val)
        }
        return nil
    }
    res["browserRequireFraudWarning"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserRequireFraudWarning(val)
        }
        return nil
    }
    res["browserRequireHighSecurityForRestrictedSites"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserRequireHighSecurityForRestrictedSites(val)
        }
        return nil
    }
    res["browserRequireSmartScreen"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserRequireSmartScreen(val)
        }
        return nil
    }
    res["browserTrustedSitesSecurityLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSiteSecurityLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBrowserTrustedSitesSecurityLevel(val.(*SiteSecurityLevel))
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
    res["diagnosticsBlockDataSubmission"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDiagnosticsBlockDataSubmission(val)
        }
        return nil
    }
    res["passwordBlockPicturePasswordAndPin"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasswordBlockPicturePasswordAndPin(val)
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
    res["updatesRequireAutomaticUpdates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUpdatesRequireAutomaticUpdates(val)
        }
        return nil
    }
    res["userAccountControlSettings"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsUserAccountControlSettings)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserAccountControlSettings(val.(*WindowsUserAccountControlSettings))
        }
        return nil
    }
    res["workFoldersUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorkFoldersUrl(val)
        }
        return nil
    }
    return res
}
// GetPasswordBlockPicturePasswordAndPin gets the passwordBlockPicturePasswordAndPin property value. Indicates whether or not to Block the user from using a pictures password and pin.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetPasswordBlockPicturePasswordAndPin()(*bool) {
    val, err := m.GetBackingStore().Get("passwordBlockPicturePasswordAndPin")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPasswordExpirationDays gets the passwordExpirationDays property value. Password expiration in days.
// returns a *int32 when successful
func (m *Windows81GeneralConfiguration) GetPasswordExpirationDays()(*int32) {
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
func (m *Windows81GeneralConfiguration) GetPasswordMinimumCharacterSetCount()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinimumCharacterSetCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordMinimumLength gets the passwordMinimumLength property value. The minimum password length.
// returns a *int32 when successful
func (m *Windows81GeneralConfiguration) GetPasswordMinimumLength()(*int32) {
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
func (m *Windows81GeneralConfiguration) GetPasswordMinutesOfInactivityBeforeScreenTimeout()(*int32) {
    val, err := m.GetBackingStore().Get("passwordMinutesOfInactivityBeforeScreenTimeout")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordPreviousPasswordBlockCount gets the passwordPreviousPasswordBlockCount property value. The number of previous passwords to prevent re-use of. Valid values 0 to 24
// returns a *int32 when successful
func (m *Windows81GeneralConfiguration) GetPasswordPreviousPasswordBlockCount()(*int32) {
    val, err := m.GetBackingStore().Get("passwordPreviousPasswordBlockCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPasswordRequiredType gets the passwordRequiredType property value. Possible values of required passwords.
// returns a *RequiredPasswordType when successful
func (m *Windows81GeneralConfiguration) GetPasswordRequiredType()(*RequiredPasswordType) {
    val, err := m.GetBackingStore().Get("passwordRequiredType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RequiredPasswordType)
    }
    return nil
}
// GetPasswordSignInFailureCountBeforeFactoryReset gets the passwordSignInFailureCountBeforeFactoryReset property value. The number of sign in failures before factory reset.
// returns a *int32 when successful
func (m *Windows81GeneralConfiguration) GetPasswordSignInFailureCountBeforeFactoryReset()(*int32) {
    val, err := m.GetBackingStore().Get("passwordSignInFailureCountBeforeFactoryReset")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetStorageRequireDeviceEncryption gets the storageRequireDeviceEncryption property value. Indicates whether or not to require encryption on a mobile device.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetStorageRequireDeviceEncryption()(*bool) {
    val, err := m.GetBackingStore().Get("storageRequireDeviceEncryption")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUpdatesRequireAutomaticUpdates gets the updatesRequireAutomaticUpdates property value. Indicates whether or not to require automatic updates.
// returns a *bool when successful
func (m *Windows81GeneralConfiguration) GetUpdatesRequireAutomaticUpdates()(*bool) {
    val, err := m.GetBackingStore().Get("updatesRequireAutomaticUpdates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUserAccountControlSettings gets the userAccountControlSettings property value. Possible values for Windows user account control settings.
// returns a *WindowsUserAccountControlSettings when successful
func (m *Windows81GeneralConfiguration) GetUserAccountControlSettings()(*WindowsUserAccountControlSettings) {
    val, err := m.GetBackingStore().Get("userAccountControlSettings")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsUserAccountControlSettings)
    }
    return nil
}
// GetWorkFoldersUrl gets the workFoldersUrl property value. The work folders url.
// returns a *string when successful
func (m *Windows81GeneralConfiguration) GetWorkFoldersUrl()(*string) {
    val, err := m.GetBackingStore().Get("workFoldersUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Windows81GeneralConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteBoolValue("browserBlockAutofill", m.GetBrowserBlockAutofill())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserBlockAutomaticDetectionOfIntranetSites", m.GetBrowserBlockAutomaticDetectionOfIntranetSites())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserBlockEnterpriseModeAccess", m.GetBrowserBlockEnterpriseModeAccess())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserBlockJavaScript", m.GetBrowserBlockJavaScript())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserBlockPlugins", m.GetBrowserBlockPlugins())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserBlockPopups", m.GetBrowserBlockPopups())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserBlockSendingDoNotTrackHeader", m.GetBrowserBlockSendingDoNotTrackHeader())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserBlockSingleWordEntryOnIntranetSites", m.GetBrowserBlockSingleWordEntryOnIntranetSites())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("browserEnterpriseModeSiteListLocation", m.GetBrowserEnterpriseModeSiteListLocation())
        if err != nil {
            return err
        }
    }
    if m.GetBrowserInternetSecurityLevel() != nil {
        cast := (*m.GetBrowserInternetSecurityLevel()).String()
        err = writer.WriteStringValue("browserInternetSecurityLevel", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetBrowserIntranetSecurityLevel() != nil {
        cast := (*m.GetBrowserIntranetSecurityLevel()).String()
        err = writer.WriteStringValue("browserIntranetSecurityLevel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("browserLoggingReportLocation", m.GetBrowserLoggingReportLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserRequireFirewall", m.GetBrowserRequireFirewall())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserRequireFraudWarning", m.GetBrowserRequireFraudWarning())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserRequireHighSecurityForRestrictedSites", m.GetBrowserRequireHighSecurityForRestrictedSites())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("browserRequireSmartScreen", m.GetBrowserRequireSmartScreen())
        if err != nil {
            return err
        }
    }
    if m.GetBrowserTrustedSitesSecurityLevel() != nil {
        cast := (*m.GetBrowserTrustedSitesSecurityLevel()).String()
        err = writer.WriteStringValue("browserTrustedSitesSecurityLevel", &cast)
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
        err = writer.WriteBoolValue("diagnosticsBlockDataSubmission", m.GetDiagnosticsBlockDataSubmission())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("passwordBlockPicturePasswordAndPin", m.GetPasswordBlockPicturePasswordAndPin())
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
        err = writer.WriteBoolValue("storageRequireDeviceEncryption", m.GetStorageRequireDeviceEncryption())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("updatesRequireAutomaticUpdates", m.GetUpdatesRequireAutomaticUpdates())
        if err != nil {
            return err
        }
    }
    if m.GetUserAccountControlSettings() != nil {
        cast := (*m.GetUserAccountControlSettings()).String()
        err = writer.WriteStringValue("userAccountControlSettings", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("workFoldersUrl", m.GetWorkFoldersUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccountsBlockAddingNonMicrosoftAccountEmail sets the accountsBlockAddingNonMicrosoftAccountEmail property value. Indicates whether or not to Block the user from adding email accounts to the device that are not associated with a Microsoft account.
func (m *Windows81GeneralConfiguration) SetAccountsBlockAddingNonMicrosoftAccountEmail(value *bool)() {
    err := m.GetBackingStore().Set("accountsBlockAddingNonMicrosoftAccountEmail", value)
    if err != nil {
        panic(err)
    }
}
// SetApplyOnlyToWindows81 sets the applyOnlyToWindows81 property value. Value indicating whether this policy only applies to Windows 8.1. This property is read-only.
func (m *Windows81GeneralConfiguration) SetApplyOnlyToWindows81(value *bool)() {
    err := m.GetBackingStore().Set("applyOnlyToWindows81", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserBlockAutofill sets the browserBlockAutofill property value. Indicates whether or not to block auto fill.
func (m *Windows81GeneralConfiguration) SetBrowserBlockAutofill(value *bool)() {
    err := m.GetBackingStore().Set("browserBlockAutofill", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserBlockAutomaticDetectionOfIntranetSites sets the browserBlockAutomaticDetectionOfIntranetSites property value. Indicates whether or not to block automatic detection of Intranet sites.
func (m *Windows81GeneralConfiguration) SetBrowserBlockAutomaticDetectionOfIntranetSites(value *bool)() {
    err := m.GetBackingStore().Set("browserBlockAutomaticDetectionOfIntranetSites", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserBlockEnterpriseModeAccess sets the browserBlockEnterpriseModeAccess property value. Indicates whether or not to block enterprise mode access.
func (m *Windows81GeneralConfiguration) SetBrowserBlockEnterpriseModeAccess(value *bool)() {
    err := m.GetBackingStore().Set("browserBlockEnterpriseModeAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserBlockJavaScript sets the browserBlockJavaScript property value. Indicates whether or not to Block the user from using JavaScript.
func (m *Windows81GeneralConfiguration) SetBrowserBlockJavaScript(value *bool)() {
    err := m.GetBackingStore().Set("browserBlockJavaScript", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserBlockPlugins sets the browserBlockPlugins property value. Indicates whether or not to block plug-ins.
func (m *Windows81GeneralConfiguration) SetBrowserBlockPlugins(value *bool)() {
    err := m.GetBackingStore().Set("browserBlockPlugins", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserBlockPopups sets the browserBlockPopups property value. Indicates whether or not to block popups.
func (m *Windows81GeneralConfiguration) SetBrowserBlockPopups(value *bool)() {
    err := m.GetBackingStore().Set("browserBlockPopups", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserBlockSendingDoNotTrackHeader sets the browserBlockSendingDoNotTrackHeader property value. Indicates whether or not to Block the user from sending the do not track header.
func (m *Windows81GeneralConfiguration) SetBrowserBlockSendingDoNotTrackHeader(value *bool)() {
    err := m.GetBackingStore().Set("browserBlockSendingDoNotTrackHeader", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserBlockSingleWordEntryOnIntranetSites sets the browserBlockSingleWordEntryOnIntranetSites property value. Indicates whether or not to block a single word entry on Intranet sites.
func (m *Windows81GeneralConfiguration) SetBrowserBlockSingleWordEntryOnIntranetSites(value *bool)() {
    err := m.GetBackingStore().Set("browserBlockSingleWordEntryOnIntranetSites", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserEnterpriseModeSiteListLocation sets the browserEnterpriseModeSiteListLocation property value. The enterprise mode site list location. Could be a local file, local network or http location.
func (m *Windows81GeneralConfiguration) SetBrowserEnterpriseModeSiteListLocation(value *string)() {
    err := m.GetBackingStore().Set("browserEnterpriseModeSiteListLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserInternetSecurityLevel sets the browserInternetSecurityLevel property value. Possible values for internet site security level.
func (m *Windows81GeneralConfiguration) SetBrowserInternetSecurityLevel(value *InternetSiteSecurityLevel)() {
    err := m.GetBackingStore().Set("browserInternetSecurityLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserIntranetSecurityLevel sets the browserIntranetSecurityLevel property value. Possible values for site security level.
func (m *Windows81GeneralConfiguration) SetBrowserIntranetSecurityLevel(value *SiteSecurityLevel)() {
    err := m.GetBackingStore().Set("browserIntranetSecurityLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserLoggingReportLocation sets the browserLoggingReportLocation property value. The logging report location.
func (m *Windows81GeneralConfiguration) SetBrowserLoggingReportLocation(value *string)() {
    err := m.GetBackingStore().Set("browserLoggingReportLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserRequireFirewall sets the browserRequireFirewall property value. Indicates whether or not to require a firewall.
func (m *Windows81GeneralConfiguration) SetBrowserRequireFirewall(value *bool)() {
    err := m.GetBackingStore().Set("browserRequireFirewall", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserRequireFraudWarning sets the browserRequireFraudWarning property value. Indicates whether or not to require fraud warning.
func (m *Windows81GeneralConfiguration) SetBrowserRequireFraudWarning(value *bool)() {
    err := m.GetBackingStore().Set("browserRequireFraudWarning", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserRequireHighSecurityForRestrictedSites sets the browserRequireHighSecurityForRestrictedSites property value. Indicates whether or not to require high security for restricted sites.
func (m *Windows81GeneralConfiguration) SetBrowserRequireHighSecurityForRestrictedSites(value *bool)() {
    err := m.GetBackingStore().Set("browserRequireHighSecurityForRestrictedSites", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserRequireSmartScreen sets the browserRequireSmartScreen property value. Indicates whether or not to require the user to use the smart screen filter.
func (m *Windows81GeneralConfiguration) SetBrowserRequireSmartScreen(value *bool)() {
    err := m.GetBackingStore().Set("browserRequireSmartScreen", value)
    if err != nil {
        panic(err)
    }
}
// SetBrowserTrustedSitesSecurityLevel sets the browserTrustedSitesSecurityLevel property value. Possible values for site security level.
func (m *Windows81GeneralConfiguration) SetBrowserTrustedSitesSecurityLevel(value *SiteSecurityLevel)() {
    err := m.GetBackingStore().Set("browserTrustedSitesSecurityLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetCellularBlockDataRoaming sets the cellularBlockDataRoaming property value. Indicates whether or not to block data roaming.
func (m *Windows81GeneralConfiguration) SetCellularBlockDataRoaming(value *bool)() {
    err := m.GetBackingStore().Set("cellularBlockDataRoaming", value)
    if err != nil {
        panic(err)
    }
}
// SetDiagnosticsBlockDataSubmission sets the diagnosticsBlockDataSubmission property value. Indicates whether or not to block diagnostic data submission.
func (m *Windows81GeneralConfiguration) SetDiagnosticsBlockDataSubmission(value *bool)() {
    err := m.GetBackingStore().Set("diagnosticsBlockDataSubmission", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordBlockPicturePasswordAndPin sets the passwordBlockPicturePasswordAndPin property value. Indicates whether or not to Block the user from using a pictures password and pin.
func (m *Windows81GeneralConfiguration) SetPasswordBlockPicturePasswordAndPin(value *bool)() {
    err := m.GetBackingStore().Set("passwordBlockPicturePasswordAndPin", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordExpirationDays sets the passwordExpirationDays property value. Password expiration in days.
func (m *Windows81GeneralConfiguration) SetPasswordExpirationDays(value *int32)() {
    err := m.GetBackingStore().Set("passwordExpirationDays", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinimumCharacterSetCount sets the passwordMinimumCharacterSetCount property value. The number of character sets required in the password.
func (m *Windows81GeneralConfiguration) SetPasswordMinimumCharacterSetCount(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinimumCharacterSetCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinimumLength sets the passwordMinimumLength property value. The minimum password length.
func (m *Windows81GeneralConfiguration) SetPasswordMinimumLength(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinimumLength", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordMinutesOfInactivityBeforeScreenTimeout sets the passwordMinutesOfInactivityBeforeScreenTimeout property value. The minutes of inactivity before the screen times out.
func (m *Windows81GeneralConfiguration) SetPasswordMinutesOfInactivityBeforeScreenTimeout(value *int32)() {
    err := m.GetBackingStore().Set("passwordMinutesOfInactivityBeforeScreenTimeout", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordPreviousPasswordBlockCount sets the passwordPreviousPasswordBlockCount property value. The number of previous passwords to prevent re-use of. Valid values 0 to 24
func (m *Windows81GeneralConfiguration) SetPasswordPreviousPasswordBlockCount(value *int32)() {
    err := m.GetBackingStore().Set("passwordPreviousPasswordBlockCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordRequiredType sets the passwordRequiredType property value. Possible values of required passwords.
func (m *Windows81GeneralConfiguration) SetPasswordRequiredType(value *RequiredPasswordType)() {
    err := m.GetBackingStore().Set("passwordRequiredType", value)
    if err != nil {
        panic(err)
    }
}
// SetPasswordSignInFailureCountBeforeFactoryReset sets the passwordSignInFailureCountBeforeFactoryReset property value. The number of sign in failures before factory reset.
func (m *Windows81GeneralConfiguration) SetPasswordSignInFailureCountBeforeFactoryReset(value *int32)() {
    err := m.GetBackingStore().Set("passwordSignInFailureCountBeforeFactoryReset", value)
    if err != nil {
        panic(err)
    }
}
// SetStorageRequireDeviceEncryption sets the storageRequireDeviceEncryption property value. Indicates whether or not to require encryption on a mobile device.
func (m *Windows81GeneralConfiguration) SetStorageRequireDeviceEncryption(value *bool)() {
    err := m.GetBackingStore().Set("storageRequireDeviceEncryption", value)
    if err != nil {
        panic(err)
    }
}
// SetUpdatesRequireAutomaticUpdates sets the updatesRequireAutomaticUpdates property value. Indicates whether or not to require automatic updates.
func (m *Windows81GeneralConfiguration) SetUpdatesRequireAutomaticUpdates(value *bool)() {
    err := m.GetBackingStore().Set("updatesRequireAutomaticUpdates", value)
    if err != nil {
        panic(err)
    }
}
// SetUserAccountControlSettings sets the userAccountControlSettings property value. Possible values for Windows user account control settings.
func (m *Windows81GeneralConfiguration) SetUserAccountControlSettings(value *WindowsUserAccountControlSettings)() {
    err := m.GetBackingStore().Set("userAccountControlSettings", value)
    if err != nil {
        panic(err)
    }
}
// SetWorkFoldersUrl sets the workFoldersUrl property value. The work folders url.
func (m *Windows81GeneralConfiguration) SetWorkFoldersUrl(value *string)() {
    err := m.GetBackingStore().Set("workFoldersUrl", value)
    if err != nil {
        panic(err)
    }
}
type Windows81GeneralConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccountsBlockAddingNonMicrosoftAccountEmail()(*bool)
    GetApplyOnlyToWindows81()(*bool)
    GetBrowserBlockAutofill()(*bool)
    GetBrowserBlockAutomaticDetectionOfIntranetSites()(*bool)
    GetBrowserBlockEnterpriseModeAccess()(*bool)
    GetBrowserBlockJavaScript()(*bool)
    GetBrowserBlockPlugins()(*bool)
    GetBrowserBlockPopups()(*bool)
    GetBrowserBlockSendingDoNotTrackHeader()(*bool)
    GetBrowserBlockSingleWordEntryOnIntranetSites()(*bool)
    GetBrowserEnterpriseModeSiteListLocation()(*string)
    GetBrowserInternetSecurityLevel()(*InternetSiteSecurityLevel)
    GetBrowserIntranetSecurityLevel()(*SiteSecurityLevel)
    GetBrowserLoggingReportLocation()(*string)
    GetBrowserRequireFirewall()(*bool)
    GetBrowserRequireFraudWarning()(*bool)
    GetBrowserRequireHighSecurityForRestrictedSites()(*bool)
    GetBrowserRequireSmartScreen()(*bool)
    GetBrowserTrustedSitesSecurityLevel()(*SiteSecurityLevel)
    GetCellularBlockDataRoaming()(*bool)
    GetDiagnosticsBlockDataSubmission()(*bool)
    GetPasswordBlockPicturePasswordAndPin()(*bool)
    GetPasswordExpirationDays()(*int32)
    GetPasswordMinimumCharacterSetCount()(*int32)
    GetPasswordMinimumLength()(*int32)
    GetPasswordMinutesOfInactivityBeforeScreenTimeout()(*int32)
    GetPasswordPreviousPasswordBlockCount()(*int32)
    GetPasswordRequiredType()(*RequiredPasswordType)
    GetPasswordSignInFailureCountBeforeFactoryReset()(*int32)
    GetStorageRequireDeviceEncryption()(*bool)
    GetUpdatesRequireAutomaticUpdates()(*bool)
    GetUserAccountControlSettings()(*WindowsUserAccountControlSettings)
    GetWorkFoldersUrl()(*string)
    SetAccountsBlockAddingNonMicrosoftAccountEmail(value *bool)()
    SetApplyOnlyToWindows81(value *bool)()
    SetBrowserBlockAutofill(value *bool)()
    SetBrowserBlockAutomaticDetectionOfIntranetSites(value *bool)()
    SetBrowserBlockEnterpriseModeAccess(value *bool)()
    SetBrowserBlockJavaScript(value *bool)()
    SetBrowserBlockPlugins(value *bool)()
    SetBrowserBlockPopups(value *bool)()
    SetBrowserBlockSendingDoNotTrackHeader(value *bool)()
    SetBrowserBlockSingleWordEntryOnIntranetSites(value *bool)()
    SetBrowserEnterpriseModeSiteListLocation(value *string)()
    SetBrowserInternetSecurityLevel(value *InternetSiteSecurityLevel)()
    SetBrowserIntranetSecurityLevel(value *SiteSecurityLevel)()
    SetBrowserLoggingReportLocation(value *string)()
    SetBrowserRequireFirewall(value *bool)()
    SetBrowserRequireFraudWarning(value *bool)()
    SetBrowserRequireHighSecurityForRestrictedSites(value *bool)()
    SetBrowserRequireSmartScreen(value *bool)()
    SetBrowserTrustedSitesSecurityLevel(value *SiteSecurityLevel)()
    SetCellularBlockDataRoaming(value *bool)()
    SetDiagnosticsBlockDataSubmission(value *bool)()
    SetPasswordBlockPicturePasswordAndPin(value *bool)()
    SetPasswordExpirationDays(value *int32)()
    SetPasswordMinimumCharacterSetCount(value *int32)()
    SetPasswordMinimumLength(value *int32)()
    SetPasswordMinutesOfInactivityBeforeScreenTimeout(value *int32)()
    SetPasswordPreviousPasswordBlockCount(value *int32)()
    SetPasswordRequiredType(value *RequiredPasswordType)()
    SetPasswordSignInFailureCountBeforeFactoryReset(value *int32)()
    SetStorageRequireDeviceEncryption(value *bool)()
    SetUpdatesRequireAutomaticUpdates(value *bool)()
    SetUserAccountControlSettings(value *WindowsUserAccountControlSettings)()
    SetWorkFoldersUrl(value *string)()
}
