package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedDeviceMobileAppConfiguration an abstract class for Mobile app configuration for enrolled devices.
type ManagedDeviceMobileAppConfiguration struct {
    Entity
}
// NewManagedDeviceMobileAppConfiguration instantiates a new ManagedDeviceMobileAppConfiguration and sets the default values.
func NewManagedDeviceMobileAppConfiguration()(*ManagedDeviceMobileAppConfiguration) {
    m := &ManagedDeviceMobileAppConfiguration{
        Entity: *NewEntity(),
    }
    return m
}
// CreateManagedDeviceMobileAppConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedDeviceMobileAppConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
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
                    case "#microsoft.graph.iosMobileAppConfiguration":
                        return NewIosMobileAppConfiguration(), nil
                }
            }
        }
    }
    return NewManagedDeviceMobileAppConfiguration(), nil
}
// GetAssignments gets the assignments property value. The list of group assignemenets for app configration.
// returns a []ManagedDeviceMobileAppConfigurationAssignmentable when successful
func (m *ManagedDeviceMobileAppConfiguration) GetAssignments()([]ManagedDeviceMobileAppConfigurationAssignmentable) {
    val, err := m.GetBackingStore().Get("assignments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedDeviceMobileAppConfigurationAssignmentable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. DateTime the object was created.
// returns a *Time when successful
func (m *ManagedDeviceMobileAppConfiguration) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
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
func (m *ManagedDeviceMobileAppConfiguration) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeviceStatuses gets the deviceStatuses property value. List of ManagedDeviceMobileAppConfigurationDeviceStatus.
// returns a []ManagedDeviceMobileAppConfigurationDeviceStatusable when successful
func (m *ManagedDeviceMobileAppConfiguration) GetDeviceStatuses()([]ManagedDeviceMobileAppConfigurationDeviceStatusable) {
    val, err := m.GetBackingStore().Get("deviceStatuses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedDeviceMobileAppConfigurationDeviceStatusable)
    }
    return nil
}
// GetDeviceStatusSummary gets the deviceStatusSummary property value. App configuration device status summary.
// returns a ManagedDeviceMobileAppConfigurationDeviceSummaryable when successful
func (m *ManagedDeviceMobileAppConfiguration) GetDeviceStatusSummary()(ManagedDeviceMobileAppConfigurationDeviceSummaryable) {
    val, err := m.GetBackingStore().Get("deviceStatusSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ManagedDeviceMobileAppConfigurationDeviceSummaryable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Admin provided name of the device configuration.
// returns a *string when successful
func (m *ManagedDeviceMobileAppConfiguration) GetDisplayName()(*string) {
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
func (m *ManagedDeviceMobileAppConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["assignments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedDeviceMobileAppConfigurationAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedDeviceMobileAppConfigurationAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedDeviceMobileAppConfigurationAssignmentable)
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
    res["deviceStatuses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedDeviceMobileAppConfigurationDeviceStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedDeviceMobileAppConfigurationDeviceStatusable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedDeviceMobileAppConfigurationDeviceStatusable)
                }
            }
            m.SetDeviceStatuses(res)
        }
        return nil
    }
    res["deviceStatusSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateManagedDeviceMobileAppConfigurationDeviceSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceStatusSummary(val.(ManagedDeviceMobileAppConfigurationDeviceSummaryable))
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
    res["targetedMobileApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetTargetedMobileApps(res)
        }
        return nil
    }
    res["userStatuses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateManagedDeviceMobileAppConfigurationUserStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ManagedDeviceMobileAppConfigurationUserStatusable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ManagedDeviceMobileAppConfigurationUserStatusable)
                }
            }
            m.SetUserStatuses(res)
        }
        return nil
    }
    res["userStatusSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateManagedDeviceMobileAppConfigurationUserSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserStatusSummary(val.(ManagedDeviceMobileAppConfigurationUserSummaryable))
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
func (m *ManagedDeviceMobileAppConfiguration) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetTargetedMobileApps gets the targetedMobileApps property value. the associated app.
// returns a []string when successful
func (m *ManagedDeviceMobileAppConfiguration) GetTargetedMobileApps()([]string) {
    val, err := m.GetBackingStore().Get("targetedMobileApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetUserStatuses gets the userStatuses property value. List of ManagedDeviceMobileAppConfigurationUserStatus.
// returns a []ManagedDeviceMobileAppConfigurationUserStatusable when successful
func (m *ManagedDeviceMobileAppConfiguration) GetUserStatuses()([]ManagedDeviceMobileAppConfigurationUserStatusable) {
    val, err := m.GetBackingStore().Get("userStatuses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ManagedDeviceMobileAppConfigurationUserStatusable)
    }
    return nil
}
// GetUserStatusSummary gets the userStatusSummary property value. App configuration user status summary.
// returns a ManagedDeviceMobileAppConfigurationUserSummaryable when successful
func (m *ManagedDeviceMobileAppConfiguration) GetUserStatusSummary()(ManagedDeviceMobileAppConfigurationUserSummaryable) {
    val, err := m.GetBackingStore().Get("userStatusSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ManagedDeviceMobileAppConfigurationUserSummaryable)
    }
    return nil
}
// GetVersion gets the version property value. Version of the device configuration.
// returns a *int32 when successful
func (m *ManagedDeviceMobileAppConfiguration) GetVersion()(*int32) {
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
func (m *ManagedDeviceMobileAppConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteObjectValue("deviceStatusSummary", m.GetDeviceStatusSummary())
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
    if m.GetTargetedMobileApps() != nil {
        err = writer.WriteCollectionOfStringValues("targetedMobileApps", m.GetTargetedMobileApps())
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
        err = writer.WriteObjectValue("userStatusSummary", m.GetUserStatusSummary())
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
// SetAssignments sets the assignments property value. The list of group assignemenets for app configration.
func (m *ManagedDeviceMobileAppConfiguration) SetAssignments(value []ManagedDeviceMobileAppConfigurationAssignmentable)() {
    err := m.GetBackingStore().Set("assignments", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. DateTime the object was created.
func (m *ManagedDeviceMobileAppConfiguration) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Admin provided description of the Device Configuration.
func (m *ManagedDeviceMobileAppConfiguration) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceStatuses sets the deviceStatuses property value. List of ManagedDeviceMobileAppConfigurationDeviceStatus.
func (m *ManagedDeviceMobileAppConfiguration) SetDeviceStatuses(value []ManagedDeviceMobileAppConfigurationDeviceStatusable)() {
    err := m.GetBackingStore().Set("deviceStatuses", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceStatusSummary sets the deviceStatusSummary property value. App configuration device status summary.
func (m *ManagedDeviceMobileAppConfiguration) SetDeviceStatusSummary(value ManagedDeviceMobileAppConfigurationDeviceSummaryable)() {
    err := m.GetBackingStore().Set("deviceStatusSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Admin provided name of the device configuration.
func (m *ManagedDeviceMobileAppConfiguration) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. DateTime the object was last modified.
func (m *ManagedDeviceMobileAppConfiguration) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetedMobileApps sets the targetedMobileApps property value. the associated app.
func (m *ManagedDeviceMobileAppConfiguration) SetTargetedMobileApps(value []string)() {
    err := m.GetBackingStore().Set("targetedMobileApps", value)
    if err != nil {
        panic(err)
    }
}
// SetUserStatuses sets the userStatuses property value. List of ManagedDeviceMobileAppConfigurationUserStatus.
func (m *ManagedDeviceMobileAppConfiguration) SetUserStatuses(value []ManagedDeviceMobileAppConfigurationUserStatusable)() {
    err := m.GetBackingStore().Set("userStatuses", value)
    if err != nil {
        panic(err)
    }
}
// SetUserStatusSummary sets the userStatusSummary property value. App configuration user status summary.
func (m *ManagedDeviceMobileAppConfiguration) SetUserStatusSummary(value ManagedDeviceMobileAppConfigurationUserSummaryable)() {
    err := m.GetBackingStore().Set("userStatusSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetVersion sets the version property value. Version of the device configuration.
func (m *ManagedDeviceMobileAppConfiguration) SetVersion(value *int32)() {
    err := m.GetBackingStore().Set("version", value)
    if err != nil {
        panic(err)
    }
}
type ManagedDeviceMobileAppConfigurationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAssignments()([]ManagedDeviceMobileAppConfigurationAssignmentable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDeviceStatuses()([]ManagedDeviceMobileAppConfigurationDeviceStatusable)
    GetDeviceStatusSummary()(ManagedDeviceMobileAppConfigurationDeviceSummaryable)
    GetDisplayName()(*string)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetTargetedMobileApps()([]string)
    GetUserStatuses()([]ManagedDeviceMobileAppConfigurationUserStatusable)
    GetUserStatusSummary()(ManagedDeviceMobileAppConfigurationUserSummaryable)
    GetVersion()(*int32)
    SetAssignments(value []ManagedDeviceMobileAppConfigurationAssignmentable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDeviceStatuses(value []ManagedDeviceMobileAppConfigurationDeviceStatusable)()
    SetDeviceStatusSummary(value ManagedDeviceMobileAppConfigurationDeviceSummaryable)()
    SetDisplayName(value *string)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetTargetedMobileApps(value []string)()
    SetUserStatuses(value []ManagedDeviceMobileAppConfigurationUserStatusable)()
    SetUserStatusSummary(value ManagedDeviceMobileAppConfigurationUserSummaryable)()
    SetVersion(value *int32)()
}
