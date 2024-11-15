package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceCompliancePolicy this is the base class for Compliance policy. Compliance policies are platform specific and individual per-platform compliance policies inherit from here. 
type DeviceCompliancePolicy struct {
    Entity
}
// NewDeviceCompliancePolicy instantiates a new DeviceCompliancePolicy and sets the default values.
func NewDeviceCompliancePolicy()(*DeviceCompliancePolicy) {
    m := &DeviceCompliancePolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceCompliancePolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceCompliancePolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.androidCompliancePolicy":
                        return NewAndroidCompliancePolicy(), nil
                    case "#microsoft.graph.androidWorkProfileCompliancePolicy":
                        return NewAndroidWorkProfileCompliancePolicy(), nil
                    case "#microsoft.graph.iosCompliancePolicy":
                        return NewIosCompliancePolicy(), nil
                    case "#microsoft.graph.macOSCompliancePolicy":
                        return NewMacOSCompliancePolicy(), nil
                    case "#microsoft.graph.windows10CompliancePolicy":
                        return NewWindows10CompliancePolicy(), nil
                    case "#microsoft.graph.windows10MobileCompliancePolicy":
                        return NewWindows10MobileCompliancePolicy(), nil
                    case "#microsoft.graph.windows81CompliancePolicy":
                        return NewWindows81CompliancePolicy(), nil
                    case "#microsoft.graph.windowsPhone81CompliancePolicy":
                        return NewWindowsPhone81CompliancePolicy(), nil
                }
            }
        }
    }
    return NewDeviceCompliancePolicy(), nil
}
// GetAssignments gets the assignments property value. The collection of assignments for this compliance policy.
// returns a []DeviceCompliancePolicyAssignmentable when successful
func (m *DeviceCompliancePolicy) GetAssignments()([]DeviceCompliancePolicyAssignmentable) {
    val, err := m.GetBackingStore().Get("assignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceCompliancePolicyAssignmentable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. DateTime the object was created.
// returns a *Time when successful
func (m *DeviceCompliancePolicy) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. Admin provided description of the Device Configuration.
// returns a *string when successful
func (m *DeviceCompliancePolicy) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceSettingStateSummaries gets the deviceSettingStateSummaries property value. Compliance Setting State Device Summary
// returns a []SettingStateDeviceSummaryable when successful
func (m *DeviceCompliancePolicy) GetDeviceSettingStateSummaries()([]SettingStateDeviceSummaryable) {
    val, err := m.GetBackingStore().Get("deviceSettingStateSummaries")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SettingStateDeviceSummaryable)
    }
    return nil
}
// GetDeviceStatuses gets the deviceStatuses property value. List of DeviceComplianceDeviceStatus.
// returns a []DeviceComplianceDeviceStatusable when successful
func (m *DeviceCompliancePolicy) GetDeviceStatuses()([]DeviceComplianceDeviceStatusable) {
    val, err := m.GetBackingStore().Get("deviceStatuses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceComplianceDeviceStatusable)
    }
    return nil
}
// GetDeviceStatusOverview gets the deviceStatusOverview property value. Device compliance devices status overview
// returns a DeviceComplianceDeviceOverviewable when successful
func (m *DeviceCompliancePolicy) GetDeviceStatusOverview()(DeviceComplianceDeviceOverviewable) {
    val, err := m.GetBackingStore().Get("deviceStatusOverview")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceComplianceDeviceOverviewable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Admin provided name of the device configuration.
// returns a *string when successful
func (m *DeviceCompliancePolicy) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceCompliancePolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["assignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceCompliancePolicyAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceCompliancePolicyAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceCompliancePolicyAssignmentable)
                }
            }
            m.SetAssignments(res)
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
        }
        return nil
    }
    res["deviceSettingStateSummaries"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSettingStateDeviceSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SettingStateDeviceSummaryable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SettingStateDeviceSummaryable)
                }
            }
            m.SetDeviceSettingStateSummaries(res)
        }
        return nil
    }
    res["deviceStatuses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceComplianceDeviceStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceComplianceDeviceStatusable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceComplianceDeviceStatusable)
                }
            }
            m.SetDeviceStatuses(res)
        }
        return nil
    }
    res["deviceStatusOverview"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceComplianceDeviceOverviewFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceStatusOverview(val.(DeviceComplianceDeviceOverviewable))
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["scheduledActionsForRule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceComplianceScheduledActionForRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceComplianceScheduledActionForRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceComplianceScheduledActionForRuleable)
                }
            }
            m.SetScheduledActionsForRule(res)
        }
        return nil
    }
    res["userStatuses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceComplianceUserStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceComplianceUserStatusable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceComplianceUserStatusable)
                }
            }
            m.SetUserStatuses(res)
        }
        return nil
    }
    res["userStatusOverview"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceComplianceUserOverviewFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserStatusOverview(val.(DeviceComplianceUserOverviewable))
        }
        return nil
    }
    res["version"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersion(val)
        }
        return nil
    }
    return res
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. DateTime the object was last modified.
// returns a *Time when successful
func (m *DeviceCompliancePolicy) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetScheduledActionsForRule gets the scheduledActionsForRule property value. The list of scheduled action per rule for this compliance policy. This is a required property when creating any individual per-platform compliance policies.
// returns a []DeviceComplianceScheduledActionForRuleable when successful
func (m *DeviceCompliancePolicy) GetScheduledActionsForRule()([]DeviceComplianceScheduledActionForRuleable) {
    val, err := m.GetBackingStore().Get("scheduledActionsForRule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceComplianceScheduledActionForRuleable)
    }
    return nil
}
// GetUserStatuses gets the userStatuses property value. List of DeviceComplianceUserStatus.
// returns a []DeviceComplianceUserStatusable when successful
func (m *DeviceCompliancePolicy) GetUserStatuses()([]DeviceComplianceUserStatusable) {
    val, err := m.GetBackingStore().Get("userStatuses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceComplianceUserStatusable)
    }
    return nil
}
// GetUserStatusOverview gets the userStatusOverview property value. Device compliance users status overview
// returns a DeviceComplianceUserOverviewable when successful
func (m *DeviceCompliancePolicy) GetUserStatusOverview()(DeviceComplianceUserOverviewable) {
    val, err := m.GetBackingStore().Get("userStatusOverview")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceComplianceUserOverviewable)
    }
    return nil
}
// GetVersion gets the version property value. Version of the device configuration.
// returns a *int32 when successful
func (m *DeviceCompliancePolicy) GetVersion()(*int32) {
    val, err := m.GetBackingStore().Get("version")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceCompliancePolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAssignments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAssignments()))
        for i, v := range m.GetAssignments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("assignments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceSettingStateSummaries() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceSettingStateSummaries()))
        for i, v := range m.GetDeviceSettingStateSummaries() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("deviceSettingStateSummaries", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDeviceStatuses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceStatuses()))
        for i, v := range m.GetDeviceStatuses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("deviceStatuses", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("deviceStatusOverview", m.GetDeviceStatusOverview())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetScheduledActionsForRule() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetScheduledActionsForRule()))
        for i, v := range m.GetScheduledActionsForRule() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("scheduledActionsForRule", cast)
        if err != nil {
            return err
        }
    }
    if m.GetUserStatuses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetUserStatuses()))
        for i, v := range m.GetUserStatuses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("userStatuses", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("userStatusOverview", m.GetUserStatusOverview())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("version", m.GetVersion())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAssignments sets the assignments property value. The collection of assignments for this compliance policy.
func (m *DeviceCompliancePolicy) SetAssignments(value []DeviceCompliancePolicyAssignmentable)() {
    err := m.GetBackingStore().Set("assignments", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. DateTime the object was created.
func (m *DeviceCompliancePolicy) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Admin provided description of the Device Configuration.
func (m *DeviceCompliancePolicy) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceSettingStateSummaries sets the deviceSettingStateSummaries property value. Compliance Setting State Device Summary
func (m *DeviceCompliancePolicy) SetDeviceSettingStateSummaries(value []SettingStateDeviceSummaryable)() {
    err := m.GetBackingStore().Set("deviceSettingStateSummaries", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceStatuses sets the deviceStatuses property value. List of DeviceComplianceDeviceStatus.
func (m *DeviceCompliancePolicy) SetDeviceStatuses(value []DeviceComplianceDeviceStatusable)() {
    err := m.GetBackingStore().Set("deviceStatuses", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceStatusOverview sets the deviceStatusOverview property value. Device compliance devices status overview
func (m *DeviceCompliancePolicy) SetDeviceStatusOverview(value DeviceComplianceDeviceOverviewable)() {
    err := m.GetBackingStore().Set("deviceStatusOverview", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Admin provided name of the device configuration.
func (m *DeviceCompliancePolicy) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. DateTime the object was last modified.
func (m *DeviceCompliancePolicy) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetScheduledActionsForRule sets the scheduledActionsForRule property value. The list of scheduled action per rule for this compliance policy. This is a required property when creating any individual per-platform compliance policies.
func (m *DeviceCompliancePolicy) SetScheduledActionsForRule(value []DeviceComplianceScheduledActionForRuleable)() {
    err := m.GetBackingStore().Set("scheduledActionsForRule", value)
    if err != nil {
        panic(err)
    }
}
// SetUserStatuses sets the userStatuses property value. List of DeviceComplianceUserStatus.
func (m *DeviceCompliancePolicy) SetUserStatuses(value []DeviceComplianceUserStatusable)() {
    err := m.GetBackingStore().Set("userStatuses", value)
    if err != nil {
        panic(err)
    }
}
// SetUserStatusOverview sets the userStatusOverview property value. Device compliance users status overview
func (m *DeviceCompliancePolicy) SetUserStatusOverview(value DeviceComplianceUserOverviewable)() {
    err := m.GetBackingStore().Set("userStatusOverview", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. Version of the device configuration.
func (m *DeviceCompliancePolicy) SetVersion(value *int32)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type DeviceCompliancePolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignments()([]DeviceCompliancePolicyAssignmentable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDeviceSettingStateSummaries()([]SettingStateDeviceSummaryable)
    GetDeviceStatuses()([]DeviceComplianceDeviceStatusable)
    GetDeviceStatusOverview()(DeviceComplianceDeviceOverviewable)
    GetDisplayName()(*string)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetScheduledActionsForRule()([]DeviceComplianceScheduledActionForRuleable)
    GetUserStatuses()([]DeviceComplianceUserStatusable)
    GetUserStatusOverview()(DeviceComplianceUserOverviewable)
    GetVersion()(*int32)
    SetAssignments(value []DeviceCompliancePolicyAssignmentable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDeviceSettingStateSummaries(value []SettingStateDeviceSummaryable)()
    SetDeviceStatuses(value []DeviceComplianceDeviceStatusable)()
    SetDeviceStatusOverview(value DeviceComplianceDeviceOverviewable)()
    SetDisplayName(value *string)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetScheduledActionsForRule(value []DeviceComplianceScheduledActionForRuleable)()
    SetUserStatuses(value []DeviceComplianceUserStatusable)()
    SetUserStatusOverview(value DeviceComplianceUserOverviewable)()
    SetVersion(value *int32)()
}
