package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MobileThreatDefenseConnector entity which represents a connection to Mobile Threat Defense partner.
type MobileThreatDefenseConnector struct {
    Entity
}
// NewMobileThreatDefenseConnector instantiates a new MobileThreatDefenseConnector and sets the default values.
func NewMobileThreatDefenseConnector()(*MobileThreatDefenseConnector) {
    m := &MobileThreatDefenseConnector{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMobileThreatDefenseConnectorFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMobileThreatDefenseConnectorFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMobileThreatDefenseConnector(), nil
}
// GetAllowPartnerToCollectIOSApplicationMetadata gets the allowPartnerToCollectIOSApplicationMetadata property value. When TRUE, indicates the Mobile Threat Defense partner may collect metadata about installed applications from Intune for IOS devices. When FALSE, indicates the Mobile Threat Defense partner may not collect metadata about installed applications from Intune for IOS devices. Default value is FALSE.
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetAllowPartnerToCollectIOSApplicationMetadata()(*bool) {
    val, err := m.GetBackingStore().Get("allowPartnerToCollectIOSApplicationMetadata")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowPartnerToCollectIOSPersonalApplicationMetadata gets the allowPartnerToCollectIOSPersonalApplicationMetadata property value. When TRUE, indicates the Mobile Threat Defense partner may collect metadata about personally installed applications from Intune for IOS devices. When FALSE, indicates the Mobile Threat Defense partner may not collect metadata about personally installed applications from Intune for IOS devices. Default value is FALSE.
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetAllowPartnerToCollectIOSPersonalApplicationMetadata()(*bool) {
    val, err := m.GetBackingStore().Get("allowPartnerToCollectIOSPersonalApplicationMetadata")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAndroidDeviceBlockedOnMissingPartnerData gets the androidDeviceBlockedOnMissingPartnerData property value. For Android, set whether Intune must receive data from the Mobile Threat Defense partner prior to marking a device compliant
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetAndroidDeviceBlockedOnMissingPartnerData()(*bool) {
    val, err := m.GetBackingStore().Get("androidDeviceBlockedOnMissingPartnerData")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAndroidEnabled gets the androidEnabled property value. For Android, set whether data from the Mobile Threat Defense partner should be used during compliance evaluations
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetAndroidEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("androidEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAndroidMobileApplicationManagementEnabled gets the androidMobileApplicationManagementEnabled property value. When TRUE, inidicates that data from the Mobile Threat Defense partner can be used during Mobile Application Management (MAM) evaluations for Android devices. When FALSE, inidicates that data from the Mobile Threat Defense partner should not be used during Mobile Application Management (MAM) evaluations for Android devices. Only one partner per platform may be enabled for Mobile Application Management (MAM) evaluation. Default value is FALSE.
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetAndroidMobileApplicationManagementEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("androidMobileApplicationManagementEnabled")
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
func (m *MobileThreatDefenseConnector) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["allowPartnerToCollectIOSApplicationMetadata"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowPartnerToCollectIOSApplicationMetadata(val)
        }
        return nil
    }
    res["allowPartnerToCollectIOSPersonalApplicationMetadata"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowPartnerToCollectIOSPersonalApplicationMetadata(val)
        }
        return nil
    }
    res["androidDeviceBlockedOnMissingPartnerData"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidDeviceBlockedOnMissingPartnerData(val)
        }
        return nil
    }
    res["androidEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidEnabled(val)
        }
        return nil
    }
    res["androidMobileApplicationManagementEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAndroidMobileApplicationManagementEnabled(val)
        }
        return nil
    }
    res["iosDeviceBlockedOnMissingPartnerData"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIosDeviceBlockedOnMissingPartnerData(val)
        }
        return nil
    }
    res["iosEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIosEnabled(val)
        }
        return nil
    }
    res["iosMobileApplicationManagementEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIosMobileApplicationManagementEnabled(val)
        }
        return nil
    }
    res["lastHeartbeatDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastHeartbeatDateTime(val)
        }
        return nil
    }
    res["microsoftDefenderForEndpointAttachEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMicrosoftDefenderForEndpointAttachEnabled(val)
        }
        return nil
    }
    res["partnerState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMobileThreatPartnerTenantState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPartnerState(val.(*MobileThreatPartnerTenantState))
        }
        return nil
    }
    res["partnerUnresponsivenessThresholdInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPartnerUnresponsivenessThresholdInDays(val)
        }
        return nil
    }
    res["partnerUnsupportedOsVersionBlocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPartnerUnsupportedOsVersionBlocked(val)
        }
        return nil
    }
    res["windowsDeviceBlockedOnMissingPartnerData"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsDeviceBlockedOnMissingPartnerData(val)
        }
        return nil
    }
    res["windowsEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsEnabled(val)
        }
        return nil
    }
    return res
}
// GetIosDeviceBlockedOnMissingPartnerData gets the iosDeviceBlockedOnMissingPartnerData property value. For IOS, set whether Intune must receive data from the Mobile Threat Defense partner prior to marking a device compliant
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetIosDeviceBlockedOnMissingPartnerData()(*bool) {
    val, err := m.GetBackingStore().Get("iosDeviceBlockedOnMissingPartnerData")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIosEnabled gets the iosEnabled property value. For IOS, get or set whether data from the Mobile Threat Defense partner should be used during compliance evaluations
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetIosEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("iosEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIosMobileApplicationManagementEnabled gets the iosMobileApplicationManagementEnabled property value. When TRUE, inidicates that data from the Mobile Threat Defense partner can be used during Mobile Application Management (MAM) evaluations for IOS devices. When FALSE, inidicates that data from the Mobile Threat Defense partner should not be used during Mobile Application Management (MAM) evaluations for IOS devices. Only one partner per platform may be enabled for Mobile Application Management (MAM) evaluation. Default value is FALSE.
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetIosMobileApplicationManagementEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("iosMobileApplicationManagementEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLastHeartbeatDateTime gets the lastHeartbeatDateTime property value. DateTime of last Heartbeat recieved from the Mobile Threat Defense partner
// returns a *Time when successful
func (m *MobileThreatDefenseConnector) GetLastHeartbeatDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastHeartbeatDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetMicrosoftDefenderForEndpointAttachEnabled gets the microsoftDefenderForEndpointAttachEnabled property value. When TRUE, inidicates that configuration profile management via Microsoft Defender for Endpoint is enabled. When FALSE, inidicates that configuration profile management via Microsoft Defender for Endpoint is disabled. Default value is FALSE.
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetMicrosoftDefenderForEndpointAttachEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("microsoftDefenderForEndpointAttachEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPartnerState gets the partnerState property value. Partner state of this tenant.
// returns a *MobileThreatPartnerTenantState when successful
func (m *MobileThreatDefenseConnector) GetPartnerState()(*MobileThreatPartnerTenantState) {
    val, err := m.GetBackingStore().Get("partnerState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MobileThreatPartnerTenantState)
    }
    return nil
}
// GetPartnerUnresponsivenessThresholdInDays gets the partnerUnresponsivenessThresholdInDays property value. Get or Set days the per tenant tolerance to unresponsiveness for this partner integration
// returns a *int32 when successful
func (m *MobileThreatDefenseConnector) GetPartnerUnresponsivenessThresholdInDays()(*int32) {
    val, err := m.GetBackingStore().Get("partnerUnresponsivenessThresholdInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPartnerUnsupportedOsVersionBlocked gets the partnerUnsupportedOsVersionBlocked property value. Get or set whether to block devices on the enabled platforms that do not meet the minimum version requirements of the Mobile Threat Defense partner
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetPartnerUnsupportedOsVersionBlocked()(*bool) {
    val, err := m.GetBackingStore().Get("partnerUnsupportedOsVersionBlocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsDeviceBlockedOnMissingPartnerData gets the windowsDeviceBlockedOnMissingPartnerData property value. When TRUE, inidicates that Intune must receive data from the Mobile Threat Defense partner prior to marking a device compliant for Windows. When FALSE, inidicates that Intune may make a device compliant without receiving data from the Mobile Threat Defense partner for Windows. Default value is FALSE.
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetWindowsDeviceBlockedOnMissingPartnerData()(*bool) {
    val, err := m.GetBackingStore().Get("windowsDeviceBlockedOnMissingPartnerData")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsEnabled gets the windowsEnabled property value. When TRUE, inidicates that data from the Mobile Threat Defense partner can be used during compliance evaluations for Windows. When FALSE, inidicates that data from the Mobile Threat Defense partner should not be used during compliance evaluations for Windows. Default value is FALSE.
// returns a *bool when successful
func (m *MobileThreatDefenseConnector) GetWindowsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("windowsEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MobileThreatDefenseConnector) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowPartnerToCollectIOSApplicationMetadata", m.GetAllowPartnerToCollectIOSApplicationMetadata())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("allowPartnerToCollectIOSPersonalApplicationMetadata", m.GetAllowPartnerToCollectIOSPersonalApplicationMetadata())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("androidDeviceBlockedOnMissingPartnerData", m.GetAndroidDeviceBlockedOnMissingPartnerData())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("androidEnabled", m.GetAndroidEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("androidMobileApplicationManagementEnabled", m.GetAndroidMobileApplicationManagementEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iosDeviceBlockedOnMissingPartnerData", m.GetIosDeviceBlockedOnMissingPartnerData())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iosEnabled", m.GetIosEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("iosMobileApplicationManagementEnabled", m.GetIosMobileApplicationManagementEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastHeartbeatDateTime", m.GetLastHeartbeatDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("microsoftDefenderForEndpointAttachEnabled", m.GetMicrosoftDefenderForEndpointAttachEnabled())
        if err != nil {
            return err
        }
    }
    if m.GetPartnerState() != nil {
        cast := (*m.GetPartnerState()).String()
        err = writer.WriteStringValue("partnerState", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("partnerUnresponsivenessThresholdInDays", m.GetPartnerUnresponsivenessThresholdInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("partnerUnsupportedOsVersionBlocked", m.GetPartnerUnsupportedOsVersionBlocked())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsDeviceBlockedOnMissingPartnerData", m.GetWindowsDeviceBlockedOnMissingPartnerData())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("windowsEnabled", m.GetWindowsEnabled())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowPartnerToCollectIOSApplicationMetadata sets the allowPartnerToCollectIOSApplicationMetadata property value. When TRUE, indicates the Mobile Threat Defense partner may collect metadata about installed applications from Intune for IOS devices. When FALSE, indicates the Mobile Threat Defense partner may not collect metadata about installed applications from Intune for IOS devices. Default value is FALSE.
func (m *MobileThreatDefenseConnector) SetAllowPartnerToCollectIOSApplicationMetadata(value *bool)() {
    err := m.GetBackingStore().Set("allowPartnerToCollectIOSApplicationMetadata", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowPartnerToCollectIOSPersonalApplicationMetadata sets the allowPartnerToCollectIOSPersonalApplicationMetadata property value. When TRUE, indicates the Mobile Threat Defense partner may collect metadata about personally installed applications from Intune for IOS devices. When FALSE, indicates the Mobile Threat Defense partner may not collect metadata about personally installed applications from Intune for IOS devices. Default value is FALSE.
func (m *MobileThreatDefenseConnector) SetAllowPartnerToCollectIOSPersonalApplicationMetadata(value *bool)() {
    err := m.GetBackingStore().Set("allowPartnerToCollectIOSPersonalApplicationMetadata", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidDeviceBlockedOnMissingPartnerData sets the androidDeviceBlockedOnMissingPartnerData property value. For Android, set whether Intune must receive data from the Mobile Threat Defense partner prior to marking a device compliant
func (m *MobileThreatDefenseConnector) SetAndroidDeviceBlockedOnMissingPartnerData(value *bool)() {
    err := m.GetBackingStore().Set("androidDeviceBlockedOnMissingPartnerData", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidEnabled sets the androidEnabled property value. For Android, set whether data from the Mobile Threat Defense partner should be used during compliance evaluations
func (m *MobileThreatDefenseConnector) SetAndroidEnabled(value *bool)() {
    err := m.GetBackingStore().Set("androidEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetAndroidMobileApplicationManagementEnabled sets the androidMobileApplicationManagementEnabled property value. When TRUE, inidicates that data from the Mobile Threat Defense partner can be used during Mobile Application Management (MAM) evaluations for Android devices. When FALSE, inidicates that data from the Mobile Threat Defense partner should not be used during Mobile Application Management (MAM) evaluations for Android devices. Only one partner per platform may be enabled for Mobile Application Management (MAM) evaluation. Default value is FALSE.
func (m *MobileThreatDefenseConnector) SetAndroidMobileApplicationManagementEnabled(value *bool)() {
    err := m.GetBackingStore().Set("androidMobileApplicationManagementEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIosDeviceBlockedOnMissingPartnerData sets the iosDeviceBlockedOnMissingPartnerData property value. For IOS, set whether Intune must receive data from the Mobile Threat Defense partner prior to marking a device compliant
func (m *MobileThreatDefenseConnector) SetIosDeviceBlockedOnMissingPartnerData(value *bool)() {
    err := m.GetBackingStore().Set("iosDeviceBlockedOnMissingPartnerData", value)
    if err != nil {
        panic(err)
    }
}
// SetIosEnabled sets the iosEnabled property value. For IOS, get or set whether data from the Mobile Threat Defense partner should be used during compliance evaluations
func (m *MobileThreatDefenseConnector) SetIosEnabled(value *bool)() {
    err := m.GetBackingStore().Set("iosEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIosMobileApplicationManagementEnabled sets the iosMobileApplicationManagementEnabled property value. When TRUE, inidicates that data from the Mobile Threat Defense partner can be used during Mobile Application Management (MAM) evaluations for IOS devices. When FALSE, inidicates that data from the Mobile Threat Defense partner should not be used during Mobile Application Management (MAM) evaluations for IOS devices. Only one partner per platform may be enabled for Mobile Application Management (MAM) evaluation. Default value is FALSE.
func (m *MobileThreatDefenseConnector) SetIosMobileApplicationManagementEnabled(value *bool)() {
    err := m.GetBackingStore().Set("iosMobileApplicationManagementEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetLastHeartbeatDateTime sets the lastHeartbeatDateTime property value. DateTime of last Heartbeat recieved from the Mobile Threat Defense partner
func (m *MobileThreatDefenseConnector) SetLastHeartbeatDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastHeartbeatDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMicrosoftDefenderForEndpointAttachEnabled sets the microsoftDefenderForEndpointAttachEnabled property value. When TRUE, inidicates that configuration profile management via Microsoft Defender for Endpoint is enabled. When FALSE, inidicates that configuration profile management via Microsoft Defender for Endpoint is disabled. Default value is FALSE.
func (m *MobileThreatDefenseConnector) SetMicrosoftDefenderForEndpointAttachEnabled(value *bool)() {
    err := m.GetBackingStore().Set("microsoftDefenderForEndpointAttachEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetPartnerState sets the partnerState property value. Partner state of this tenant.
func (m *MobileThreatDefenseConnector) SetPartnerState(value *MobileThreatPartnerTenantState)() {
    err := m.GetBackingStore().Set("partnerState", value)
    if err != nil {
        panic(err)
    }
}
// SetPartnerUnresponsivenessThresholdInDays sets the partnerUnresponsivenessThresholdInDays property value. Get or Set days the per tenant tolerance to unresponsiveness for this partner integration
func (m *MobileThreatDefenseConnector) SetPartnerUnresponsivenessThresholdInDays(value *int32)() {
    err := m.GetBackingStore().Set("partnerUnresponsivenessThresholdInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetPartnerUnsupportedOsVersionBlocked sets the partnerUnsupportedOsVersionBlocked property value. Get or set whether to block devices on the enabled platforms that do not meet the minimum version requirements of the Mobile Threat Defense partner
func (m *MobileThreatDefenseConnector) SetPartnerUnsupportedOsVersionBlocked(value *bool)() {
    err := m.GetBackingStore().Set("partnerUnsupportedOsVersionBlocked", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsDeviceBlockedOnMissingPartnerData sets the windowsDeviceBlockedOnMissingPartnerData property value. When TRUE, inidicates that Intune must receive data from the Mobile Threat Defense partner prior to marking a device compliant for Windows. When FALSE, inidicates that Intune may make a device compliant without receiving data from the Mobile Threat Defense partner for Windows. Default value is FALSE.
func (m *MobileThreatDefenseConnector) SetWindowsDeviceBlockedOnMissingPartnerData(value *bool)() {
    err := m.GetBackingStore().Set("windowsDeviceBlockedOnMissingPartnerData", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsEnabled sets the windowsEnabled property value. When TRUE, inidicates that data from the Mobile Threat Defense partner can be used during compliance evaluations for Windows. When FALSE, inidicates that data from the Mobile Threat Defense partner should not be used during compliance evaluations for Windows. Default value is FALSE.
func (m *MobileThreatDefenseConnector) SetWindowsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("windowsEnabled", value)
    if err != nil {
        panic(err)
    }
}
type MobileThreatDefenseConnectorable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowPartnerToCollectIOSApplicationMetadata()(*bool)
    GetAllowPartnerToCollectIOSPersonalApplicationMetadata()(*bool)
    GetAndroidDeviceBlockedOnMissingPartnerData()(*bool)
    GetAndroidEnabled()(*bool)
    GetAndroidMobileApplicationManagementEnabled()(*bool)
    GetIosDeviceBlockedOnMissingPartnerData()(*bool)
    GetIosEnabled()(*bool)
    GetIosMobileApplicationManagementEnabled()(*bool)
    GetLastHeartbeatDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMicrosoftDefenderForEndpointAttachEnabled()(*bool)
    GetPartnerState()(*MobileThreatPartnerTenantState)
    GetPartnerUnresponsivenessThresholdInDays()(*int32)
    GetPartnerUnsupportedOsVersionBlocked()(*bool)
    GetWindowsDeviceBlockedOnMissingPartnerData()(*bool)
    GetWindowsEnabled()(*bool)
    SetAllowPartnerToCollectIOSApplicationMetadata(value *bool)()
    SetAllowPartnerToCollectIOSPersonalApplicationMetadata(value *bool)()
    SetAndroidDeviceBlockedOnMissingPartnerData(value *bool)()
    SetAndroidEnabled(value *bool)()
    SetAndroidMobileApplicationManagementEnabled(value *bool)()
    SetIosDeviceBlockedOnMissingPartnerData(value *bool)()
    SetIosEnabled(value *bool)()
    SetIosMobileApplicationManagementEnabled(value *bool)()
    SetLastHeartbeatDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMicrosoftDefenderForEndpointAttachEnabled(value *bool)()
    SetPartnerState(value *MobileThreatPartnerTenantState)()
    SetPartnerUnresponsivenessThresholdInDays(value *int32)()
    SetPartnerUnsupportedOsVersionBlocked(value *bool)()
    SetWindowsDeviceBlockedOnMissingPartnerData(value *bool)()
    SetWindowsEnabled(value *bool)()
}
