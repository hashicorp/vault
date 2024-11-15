package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Windows10TeamGeneralConfiguration this topic provides descriptions of the declared methods, properties and relationships exposed by the windows10TeamGeneralConfiguration resource.
type Windows10TeamGeneralConfiguration struct {
    DeviceConfiguration
}
// NewWindows10TeamGeneralConfiguration instantiates a new Windows10TeamGeneralConfiguration and sets the default values.
func NewWindows10TeamGeneralConfiguration()(*Windows10TeamGeneralConfiguration) {
    m := &Windows10TeamGeneralConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.windows10TeamGeneralConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindows10TeamGeneralConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindows10TeamGeneralConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindows10TeamGeneralConfiguration(), nil
}
// GetAzureOperationalInsightsBlockTelemetry gets the azureOperationalInsightsBlockTelemetry property value. Indicates whether or not to Block Azure Operational Insights.
// returns a *bool when successful
func (m *Windows10TeamGeneralConfiguration) GetAzureOperationalInsightsBlockTelemetry()(*bool) {
    val, err := m.GetBackingStore().Get("azureOperationalInsightsBlockTelemetry")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAzureOperationalInsightsWorkspaceId gets the azureOperationalInsightsWorkspaceId property value. The Azure Operational Insights workspace id.
// returns a *string when successful
func (m *Windows10TeamGeneralConfiguration) GetAzureOperationalInsightsWorkspaceId()(*string) {
    val, err := m.GetBackingStore().Get("azureOperationalInsightsWorkspaceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAzureOperationalInsightsWorkspaceKey gets the azureOperationalInsightsWorkspaceKey property value. The Azure Operational Insights Workspace key.
// returns a *string when successful
func (m *Windows10TeamGeneralConfiguration) GetAzureOperationalInsightsWorkspaceKey()(*string) {
    val, err := m.GetBackingStore().Get("azureOperationalInsightsWorkspaceKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetConnectAppBlockAutoLaunch gets the connectAppBlockAutoLaunch property value. Specifies whether to automatically launch the Connect app whenever a projection is initiated.
// returns a *bool when successful
func (m *Windows10TeamGeneralConfiguration) GetConnectAppBlockAutoLaunch()(*bool) {
    val, err := m.GetBackingStore().Get("connectAppBlockAutoLaunch")
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
func (m *Windows10TeamGeneralConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["azureOperationalInsightsBlockTelemetry"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureOperationalInsightsBlockTelemetry(val)
        }
        return nil
    }
    res["azureOperationalInsightsWorkspaceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureOperationalInsightsWorkspaceId(val)
        }
        return nil
    }
    res["azureOperationalInsightsWorkspaceKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAzureOperationalInsightsWorkspaceKey(val)
        }
        return nil
    }
    res["connectAppBlockAutoLaunch"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConnectAppBlockAutoLaunch(val)
        }
        return nil
    }
    res["maintenanceWindowBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaintenanceWindowBlocked(val)
        }
        return nil
    }
    res["maintenanceWindowDurationInHours"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaintenanceWindowDurationInHours(val)
        }
        return nil
    }
    res["maintenanceWindowStartTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeOnlyValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaintenanceWindowStartTime(val)
        }
        return nil
    }
    res["miracastBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMiracastBlocked(val)
        }
        return nil
    }
    res["miracastChannel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMiracastChannel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMiracastChannel(val.(*MiracastChannel))
        }
        return nil
    }
    res["miracastRequirePin"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMiracastRequirePin(val)
        }
        return nil
    }
    res["settingsBlockMyMeetingsAndFiles"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockMyMeetingsAndFiles(val)
        }
        return nil
    }
    res["settingsBlockSessionResume"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockSessionResume(val)
        }
        return nil
    }
    res["settingsBlockSigninSuggestions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsBlockSigninSuggestions(val)
        }
        return nil
    }
    res["settingsDefaultVolume"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsDefaultVolume(val)
        }
        return nil
    }
    res["settingsScreenTimeoutInMinutes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsScreenTimeoutInMinutes(val)
        }
        return nil
    }
    res["settingsSessionTimeoutInMinutes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsSessionTimeoutInMinutes(val)
        }
        return nil
    }
    res["settingsSleepTimeoutInMinutes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingsSleepTimeoutInMinutes(val)
        }
        return nil
    }
    res["welcomeScreenBackgroundImageUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWelcomeScreenBackgroundImageUrl(val)
        }
        return nil
    }
    res["welcomeScreenBlockAutomaticWakeUp"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWelcomeScreenBlockAutomaticWakeUp(val)
        }
        return nil
    }
    res["welcomeScreenMeetingInformation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWelcomeScreenMeetingInformation)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWelcomeScreenMeetingInformation(val.(*WelcomeScreenMeetingInformation))
        }
        return nil
    }
    return res
}
// GetMaintenanceWindowBlocked gets the maintenanceWindowBlocked property value. Indicates whether or not to Block setting a maintenance window for device updates.
// returns a *bool when successful
func (m *Windows10TeamGeneralConfiguration) GetMaintenanceWindowBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("maintenanceWindowBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMaintenanceWindowDurationInHours gets the maintenanceWindowDurationInHours property value. Maintenance window duration for device updates. Valid values 0 to 5
// returns a *int32 when successful
func (m *Windows10TeamGeneralConfiguration) GetMaintenanceWindowDurationInHours()(*int32) {
    val, err := m.GetBackingStore().Get("maintenanceWindowDurationInHours")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMaintenanceWindowStartTime gets the maintenanceWindowStartTime property value. Maintenance window start time for device updates.
// returns a *TimeOnly when successful
func (m *Windows10TeamGeneralConfiguration) GetMaintenanceWindowStartTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly) {
    val, err := m.GetBackingStore().Get("maintenanceWindowStartTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)
    }
    return nil
}
// GetMiracastBlocked gets the miracastBlocked property value. Indicates whether or not to Block wireless projection.
// returns a *bool when successful
func (m *Windows10TeamGeneralConfiguration) GetMiracastBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("miracastBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMiracastChannel gets the miracastChannel property value. Possible values for Miracast channel.
// returns a *MiracastChannel when successful
func (m *Windows10TeamGeneralConfiguration) GetMiracastChannel()(*MiracastChannel) {
    val, err := m.GetBackingStore().Get("miracastChannel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MiracastChannel)
    }
    return nil
}
// GetMiracastRequirePin gets the miracastRequirePin property value. Indicates whether or not to require a pin for wireless projection.
// returns a *bool when successful
func (m *Windows10TeamGeneralConfiguration) GetMiracastRequirePin()(*bool) {
    val, err := m.GetBackingStore().Get("miracastRequirePin")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockMyMeetingsAndFiles gets the settingsBlockMyMeetingsAndFiles property value. Specifies whether to disable the 'My meetings and files' feature in the Start menu, which shows the signed-in user's meetings and files from Office 365.
// returns a *bool when successful
func (m *Windows10TeamGeneralConfiguration) GetSettingsBlockMyMeetingsAndFiles()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockMyMeetingsAndFiles")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockSessionResume gets the settingsBlockSessionResume property value. Specifies whether to allow the ability to resume a session when the session times out.
// returns a *bool when successful
func (m *Windows10TeamGeneralConfiguration) GetSettingsBlockSessionResume()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockSessionResume")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsBlockSigninSuggestions gets the settingsBlockSigninSuggestions property value. Specifies whether to disable auto-populating of the sign-in dialog with invitees from scheduled meetings.
// returns a *bool when successful
func (m *Windows10TeamGeneralConfiguration) GetSettingsBlockSigninSuggestions()(*bool) {
    val, err := m.GetBackingStore().Get("settingsBlockSigninSuggestions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSettingsDefaultVolume gets the settingsDefaultVolume property value. Specifies the default volume value for a new session. Permitted values are 0-100. The default is 45. Valid values 0 to 100
// returns a *int32 when successful
func (m *Windows10TeamGeneralConfiguration) GetSettingsDefaultVolume()(*int32) {
    val, err := m.GetBackingStore().Get("settingsDefaultVolume")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSettingsScreenTimeoutInMinutes gets the settingsScreenTimeoutInMinutes property value. Specifies the number of minutes until the Hub screen turns off.
// returns a *int32 when successful
func (m *Windows10TeamGeneralConfiguration) GetSettingsScreenTimeoutInMinutes()(*int32) {
    val, err := m.GetBackingStore().Get("settingsScreenTimeoutInMinutes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSettingsSessionTimeoutInMinutes gets the settingsSessionTimeoutInMinutes property value. Specifies the number of minutes until the session times out.
// returns a *int32 when successful
func (m *Windows10TeamGeneralConfiguration) GetSettingsSessionTimeoutInMinutes()(*int32) {
    val, err := m.GetBackingStore().Get("settingsSessionTimeoutInMinutes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSettingsSleepTimeoutInMinutes gets the settingsSleepTimeoutInMinutes property value. Specifies the number of minutes until the Hub enters sleep mode.
// returns a *int32 when successful
func (m *Windows10TeamGeneralConfiguration) GetSettingsSleepTimeoutInMinutes()(*int32) {
    val, err := m.GetBackingStore().Get("settingsSleepTimeoutInMinutes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetWelcomeScreenBackgroundImageUrl gets the welcomeScreenBackgroundImageUrl property value. The welcome screen background image URL. The URL must use the HTTPS protocol and return a PNG image.
// returns a *string when successful
func (m *Windows10TeamGeneralConfiguration) GetWelcomeScreenBackgroundImageUrl()(*string) {
    val, err := m.GetBackingStore().Get("welcomeScreenBackgroundImageUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWelcomeScreenBlockAutomaticWakeUp gets the welcomeScreenBlockAutomaticWakeUp property value. Indicates whether or not to Block the welcome screen from waking up automatically when someone enters the room.
// returns a *bool when successful
func (m *Windows10TeamGeneralConfiguration) GetWelcomeScreenBlockAutomaticWakeUp()(*bool) {
    val, err := m.GetBackingStore().Get("welcomeScreenBlockAutomaticWakeUp")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWelcomeScreenMeetingInformation gets the welcomeScreenMeetingInformation property value. Possible values for welcome screen meeting information.
// returns a *WelcomeScreenMeetingInformation when successful
func (m *Windows10TeamGeneralConfiguration) GetWelcomeScreenMeetingInformation()(*WelcomeScreenMeetingInformation) {
    val, err := m.GetBackingStore().Get("welcomeScreenMeetingInformation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WelcomeScreenMeetingInformation)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Windows10TeamGeneralConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("azureOperationalInsightsBlockTelemetry", m.GetAzureOperationalInsightsBlockTelemetry())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("azureOperationalInsightsWorkspaceId", m.GetAzureOperationalInsightsWorkspaceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("azureOperationalInsightsWorkspaceKey", m.GetAzureOperationalInsightsWorkspaceKey())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("connectAppBlockAutoLaunch", m.GetConnectAppBlockAutoLaunch())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("maintenanceWindowBlocked", m.GetMaintenanceWindowBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("maintenanceWindowDurationInHours", m.GetMaintenanceWindowDurationInHours())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeOnlyValue("maintenanceWindowStartTime", m.GetMaintenanceWindowStartTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("miracastBlocked", m.GetMiracastBlocked())
        if err != nil {
            return err
        }
    }
    if m.GetMiracastChannel() != nil {
        cast := (*m.GetMiracastChannel()).String()
        err = writer.WriteStringValue("miracastChannel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("miracastRequirePin", m.GetMiracastRequirePin())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockMyMeetingsAndFiles", m.GetSettingsBlockMyMeetingsAndFiles())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockSessionResume", m.GetSettingsBlockSessionResume())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("settingsBlockSigninSuggestions", m.GetSettingsBlockSigninSuggestions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("settingsDefaultVolume", m.GetSettingsDefaultVolume())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("settingsScreenTimeoutInMinutes", m.GetSettingsScreenTimeoutInMinutes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("settingsSessionTimeoutInMinutes", m.GetSettingsSessionTimeoutInMinutes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("settingsSleepTimeoutInMinutes", m.GetSettingsSleepTimeoutInMinutes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("welcomeScreenBackgroundImageUrl", m.GetWelcomeScreenBackgroundImageUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("welcomeScreenBlockAutomaticWakeUp", m.GetWelcomeScreenBlockAutomaticWakeUp())
        if err != nil {
            return err
        }
    }
    if m.GetWelcomeScreenMeetingInformation() != nil {
        cast := (*m.GetWelcomeScreenMeetingInformation()).String()
        err = writer.WriteStringValue("welcomeScreenMeetingInformation", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAzureOperationalInsightsBlockTelemetry sets the azureOperationalInsightsBlockTelemetry property value. Indicates whether or not to Block Azure Operational Insights.
func (m *Windows10TeamGeneralConfiguration) SetAzureOperationalInsightsBlockTelemetry(value *bool)() {
    err := m.GetBackingStore().Set("azureOperationalInsightsBlockTelemetry", value)
    if err != nil {
        panic(err)
    }
}
// SetAzureOperationalInsightsWorkspaceId sets the azureOperationalInsightsWorkspaceId property value. The Azure Operational Insights workspace id.
func (m *Windows10TeamGeneralConfiguration) SetAzureOperationalInsightsWorkspaceId(value *string)() {
    err := m.GetBackingStore().Set("azureOperationalInsightsWorkspaceId", value)
    if err != nil {
        panic(err)
    }
}
// SetAzureOperationalInsightsWorkspaceKey sets the azureOperationalInsightsWorkspaceKey property value. The Azure Operational Insights Workspace key.
func (m *Windows10TeamGeneralConfiguration) SetAzureOperationalInsightsWorkspaceKey(value *string)() {
    err := m.GetBackingStore().Set("azureOperationalInsightsWorkspaceKey", value)
    if err != nil {
        panic(err)
    }
}
// SetConnectAppBlockAutoLaunch sets the connectAppBlockAutoLaunch property value. Specifies whether to automatically launch the Connect app whenever a projection is initiated.
func (m *Windows10TeamGeneralConfiguration) SetConnectAppBlockAutoLaunch(value *bool)() {
    err := m.GetBackingStore().Set("connectAppBlockAutoLaunch", value)
    if err != nil {
        panic(err)
    }
}
// SetMaintenanceWindowBlocked sets the maintenanceWindowBlocked property value. Indicates whether or not to Block setting a maintenance window for device updates.
func (m *Windows10TeamGeneralConfiguration) SetMaintenanceWindowBlocked(value *bool)() {
    err := m.GetBackingStore().Set("maintenanceWindowBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetMaintenanceWindowDurationInHours sets the maintenanceWindowDurationInHours property value. Maintenance window duration for device updates. Valid values 0 to 5
func (m *Windows10TeamGeneralConfiguration) SetMaintenanceWindowDurationInHours(value *int32)() {
    err := m.GetBackingStore().Set("maintenanceWindowDurationInHours", value)
    if err != nil {
        panic(err)
    }
}
// SetMaintenanceWindowStartTime sets the maintenanceWindowStartTime property value. Maintenance window start time for device updates.
func (m *Windows10TeamGeneralConfiguration) SetMaintenanceWindowStartTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)() {
    err := m.GetBackingStore().Set("maintenanceWindowStartTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMiracastBlocked sets the miracastBlocked property value. Indicates whether or not to Block wireless projection.
func (m *Windows10TeamGeneralConfiguration) SetMiracastBlocked(value *bool)() {
    err := m.GetBackingStore().Set("miracastBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetMiracastChannel sets the miracastChannel property value. Possible values for Miracast channel.
func (m *Windows10TeamGeneralConfiguration) SetMiracastChannel(value *MiracastChannel)() {
    err := m.GetBackingStore().Set("miracastChannel", value)
    if err != nil {
        panic(err)
    }
}
// SetMiracastRequirePin sets the miracastRequirePin property value. Indicates whether or not to require a pin for wireless projection.
func (m *Windows10TeamGeneralConfiguration) SetMiracastRequirePin(value *bool)() {
    err := m.GetBackingStore().Set("miracastRequirePin", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockMyMeetingsAndFiles sets the settingsBlockMyMeetingsAndFiles property value. Specifies whether to disable the 'My meetings and files' feature in the Start menu, which shows the signed-in user's meetings and files from Office 365.
func (m *Windows10TeamGeneralConfiguration) SetSettingsBlockMyMeetingsAndFiles(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockMyMeetingsAndFiles", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockSessionResume sets the settingsBlockSessionResume property value. Specifies whether to allow the ability to resume a session when the session times out.
func (m *Windows10TeamGeneralConfiguration) SetSettingsBlockSessionResume(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockSessionResume", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsBlockSigninSuggestions sets the settingsBlockSigninSuggestions property value. Specifies whether to disable auto-populating of the sign-in dialog with invitees from scheduled meetings.
func (m *Windows10TeamGeneralConfiguration) SetSettingsBlockSigninSuggestions(value *bool)() {
    err := m.GetBackingStore().Set("settingsBlockSigninSuggestions", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsDefaultVolume sets the settingsDefaultVolume property value. Specifies the default volume value for a new session. Permitted values are 0-100. The default is 45. Valid values 0 to 100
func (m *Windows10TeamGeneralConfiguration) SetSettingsDefaultVolume(value *int32)() {
    err := m.GetBackingStore().Set("settingsDefaultVolume", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsScreenTimeoutInMinutes sets the settingsScreenTimeoutInMinutes property value. Specifies the number of minutes until the Hub screen turns off.
func (m *Windows10TeamGeneralConfiguration) SetSettingsScreenTimeoutInMinutes(value *int32)() {
    err := m.GetBackingStore().Set("settingsScreenTimeoutInMinutes", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsSessionTimeoutInMinutes sets the settingsSessionTimeoutInMinutes property value. Specifies the number of minutes until the session times out.
func (m *Windows10TeamGeneralConfiguration) SetSettingsSessionTimeoutInMinutes(value *int32)() {
    err := m.GetBackingStore().Set("settingsSessionTimeoutInMinutes", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingsSleepTimeoutInMinutes sets the settingsSleepTimeoutInMinutes property value. Specifies the number of minutes until the Hub enters sleep mode.
func (m *Windows10TeamGeneralConfiguration) SetSettingsSleepTimeoutInMinutes(value *int32)() {
    err := m.GetBackingStore().Set("settingsSleepTimeoutInMinutes", value)
    if err != nil {
        panic(err)
    }
}
// SetWelcomeScreenBackgroundImageUrl sets the welcomeScreenBackgroundImageUrl property value. The welcome screen background image URL. The URL must use the HTTPS protocol and return a PNG image.
func (m *Windows10TeamGeneralConfiguration) SetWelcomeScreenBackgroundImageUrl(value *string)() {
    err := m.GetBackingStore().Set("welcomeScreenBackgroundImageUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetWelcomeScreenBlockAutomaticWakeUp sets the welcomeScreenBlockAutomaticWakeUp property value. Indicates whether or not to Block the welcome screen from waking up automatically when someone enters the room.
func (m *Windows10TeamGeneralConfiguration) SetWelcomeScreenBlockAutomaticWakeUp(value *bool)() {
    err := m.GetBackingStore().Set("welcomeScreenBlockAutomaticWakeUp", value)
    if err != nil {
        panic(err)
    }
}
// SetWelcomeScreenMeetingInformation sets the welcomeScreenMeetingInformation property value. Possible values for welcome screen meeting information.
func (m *Windows10TeamGeneralConfiguration) SetWelcomeScreenMeetingInformation(value *WelcomeScreenMeetingInformation)() {
    err := m.GetBackingStore().Set("welcomeScreenMeetingInformation", value)
    if err != nil {
        panic(err)
    }
}
type Windows10TeamGeneralConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAzureOperationalInsightsBlockTelemetry()(*bool)
    GetAzureOperationalInsightsWorkspaceId()(*string)
    GetAzureOperationalInsightsWorkspaceKey()(*string)
    GetConnectAppBlockAutoLaunch()(*bool)
    GetMaintenanceWindowBlocked()(*bool)
    GetMaintenanceWindowDurationInHours()(*int32)
    GetMaintenanceWindowStartTime()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)
    GetMiracastBlocked()(*bool)
    GetMiracastChannel()(*MiracastChannel)
    GetMiracastRequirePin()(*bool)
    GetSettingsBlockMyMeetingsAndFiles()(*bool)
    GetSettingsBlockSessionResume()(*bool)
    GetSettingsBlockSigninSuggestions()(*bool)
    GetSettingsDefaultVolume()(*int32)
    GetSettingsScreenTimeoutInMinutes()(*int32)
    GetSettingsSessionTimeoutInMinutes()(*int32)
    GetSettingsSleepTimeoutInMinutes()(*int32)
    GetWelcomeScreenBackgroundImageUrl()(*string)
    GetWelcomeScreenBlockAutomaticWakeUp()(*bool)
    GetWelcomeScreenMeetingInformation()(*WelcomeScreenMeetingInformation)
    SetAzureOperationalInsightsBlockTelemetry(value *bool)()
    SetAzureOperationalInsightsWorkspaceId(value *string)()
    SetAzureOperationalInsightsWorkspaceKey(value *string)()
    SetConnectAppBlockAutoLaunch(value *bool)()
    SetMaintenanceWindowBlocked(value *bool)()
    SetMaintenanceWindowDurationInHours(value *int32)()
    SetMaintenanceWindowStartTime(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.TimeOnly)()
    SetMiracastBlocked(value *bool)()
    SetMiracastChannel(value *MiracastChannel)()
    SetMiracastRequirePin(value *bool)()
    SetSettingsBlockMyMeetingsAndFiles(value *bool)()
    SetSettingsBlockSessionResume(value *bool)()
    SetSettingsBlockSigninSuggestions(value *bool)()
    SetSettingsDefaultVolume(value *int32)()
    SetSettingsScreenTimeoutInMinutes(value *int32)()
    SetSettingsSessionTimeoutInMinutes(value *int32)()
    SetSettingsSleepTimeoutInMinutes(value *int32)()
    SetWelcomeScreenBackgroundImageUrl(value *string)()
    SetWelcomeScreenBlockAutomaticWakeUp(value *bool)()
    SetWelcomeScreenMeetingInformation(value *WelcomeScreenMeetingInformation)()
}
