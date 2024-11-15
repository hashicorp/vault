package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceManagementPartner entity which represents a connection to device management partner.
type DeviceManagementPartner struct {
    Entity
}
// NewDeviceManagementPartner instantiates a new DeviceManagementPartner and sets the default values.
func NewDeviceManagementPartner()(*DeviceManagementPartner) {
    m := &DeviceManagementPartner{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceManagementPartnerFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceManagementPartnerFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceManagementPartner(), nil
}
// GetDisplayName gets the displayName property value. Partner display name
// returns a *string when successful
func (m *DeviceManagementPartner) GetDisplayName()(*string) {
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
func (m *DeviceManagementPartner) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
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
    res["groupsRequiringPartnerEnrollment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceManagementPartnerAssignmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceManagementPartnerAssignmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceManagementPartnerAssignmentable)
                }
            }
            m.SetGroupsRequiringPartnerEnrollment(res)
        }
        return nil
    }
    res["isConfigured"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsConfigured(val)
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
    res["partnerAppType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceManagementPartnerAppType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPartnerAppType(val.(*DeviceManagementPartnerAppType))
        }
        return nil
    }
    res["partnerState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDeviceManagementPartnerTenantState)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPartnerState(val.(*DeviceManagementPartnerTenantState))
        }
        return nil
    }
    res["singleTenantAppId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSingleTenantAppId(val)
        }
        return nil
    }
    res["whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime(val)
        }
        return nil
    }
    res["whenPartnerDevicesWillBeRemovedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWhenPartnerDevicesWillBeRemovedDateTime(val)
        }
        return nil
    }
    return res
}
// GetGroupsRequiringPartnerEnrollment gets the groupsRequiringPartnerEnrollment property value. User groups that specifies whether enrollment is through partner.
// returns a []DeviceManagementPartnerAssignmentable when successful
func (m *DeviceManagementPartner) GetGroupsRequiringPartnerEnrollment()([]DeviceManagementPartnerAssignmentable) {
    val, err := m.GetBackingStore().Get("groupsRequiringPartnerEnrollment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceManagementPartnerAssignmentable)
    }
    return nil
}
// GetIsConfigured gets the isConfigured property value. Whether device management partner is configured or not
// returns a *bool when successful
func (m *DeviceManagementPartner) GetIsConfigured()(*bool) {
    val, err := m.GetBackingStore().Get("isConfigured")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLastHeartbeatDateTime gets the lastHeartbeatDateTime property value. Timestamp of last heartbeat after admin enabled option Connect to Device management Partner
// returns a *Time when successful
func (m *DeviceManagementPartner) GetLastHeartbeatDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastHeartbeatDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetPartnerAppType gets the partnerAppType property value. Partner App Type.
// returns a *DeviceManagementPartnerAppType when successful
func (m *DeviceManagementPartner) GetPartnerAppType()(*DeviceManagementPartnerAppType) {
    val, err := m.GetBackingStore().Get("partnerAppType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceManagementPartnerAppType)
    }
    return nil
}
// GetPartnerState gets the partnerState property value. Partner state of this tenant.
// returns a *DeviceManagementPartnerTenantState when successful
func (m *DeviceManagementPartner) GetPartnerState()(*DeviceManagementPartnerTenantState) {
    val, err := m.GetBackingStore().Get("partnerState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DeviceManagementPartnerTenantState)
    }
    return nil
}
// GetSingleTenantAppId gets the singleTenantAppId property value. Partner Single tenant App id
// returns a *string when successful
func (m *DeviceManagementPartner) GetSingleTenantAppId()(*string) {
    val, err := m.GetBackingStore().Get("singleTenantAppId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime gets the whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime property value. DateTime in UTC when PartnerDevices will be marked as NonCompliant
// returns a *Time when successful
func (m *DeviceManagementPartner) GetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetWhenPartnerDevicesWillBeRemovedDateTime gets the whenPartnerDevicesWillBeRemovedDateTime property value. DateTime in UTC when PartnerDevices will be removed
// returns a *Time when successful
func (m *DeviceManagementPartner) GetWhenPartnerDevicesWillBeRemovedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("whenPartnerDevicesWillBeRemovedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceManagementPartner) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetGroupsRequiringPartnerEnrollment() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetGroupsRequiringPartnerEnrollment()))
        for i, v := range m.GetGroupsRequiringPartnerEnrollment() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("groupsRequiringPartnerEnrollment", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isConfigured", m.GetIsConfigured())
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
    if m.GetPartnerAppType() != nil {
        cast := (*m.GetPartnerAppType()).String()
        err = writer.WriteStringValue("partnerAppType", &cast)
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
        err = writer.WriteStringValue("singleTenantAppId", m.GetSingleTenantAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime", m.GetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("whenPartnerDevicesWillBeRemovedDateTime", m.GetWhenPartnerDevicesWillBeRemovedDateTime())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. Partner display name
func (m *DeviceManagementPartner) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupsRequiringPartnerEnrollment sets the groupsRequiringPartnerEnrollment property value. User groups that specifies whether enrollment is through partner.
func (m *DeviceManagementPartner) SetGroupsRequiringPartnerEnrollment(value []DeviceManagementPartnerAssignmentable)() {
    err := m.GetBackingStore().Set("groupsRequiringPartnerEnrollment", value)
    if err != nil {
        panic(err)
    }
}
// SetIsConfigured sets the isConfigured property value. Whether device management partner is configured or not
func (m *DeviceManagementPartner) SetIsConfigured(value *bool)() {
    err := m.GetBackingStore().Set("isConfigured", value)
    if err != nil {
        panic(err)
    }
}
// SetLastHeartbeatDateTime sets the lastHeartbeatDateTime property value. Timestamp of last heartbeat after admin enabled option Connect to Device management Partner
func (m *DeviceManagementPartner) SetLastHeartbeatDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastHeartbeatDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPartnerAppType sets the partnerAppType property value. Partner App Type.
func (m *DeviceManagementPartner) SetPartnerAppType(value *DeviceManagementPartnerAppType)() {
    err := m.GetBackingStore().Set("partnerAppType", value)
    if err != nil {
        panic(err)
    }
}
// SetPartnerState sets the partnerState property value. Partner state of this tenant.
func (m *DeviceManagementPartner) SetPartnerState(value *DeviceManagementPartnerTenantState)() {
    err := m.GetBackingStore().Set("partnerState", value)
    if err != nil {
        panic(err)
    }
}
// SetSingleTenantAppId sets the singleTenantAppId property value. Partner Single tenant App id
func (m *DeviceManagementPartner) SetSingleTenantAppId(value *string)() {
    err := m.GetBackingStore().Set("singleTenantAppId", value)
    if err != nil {
        panic(err)
    }
}
// SetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime sets the whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime property value. DateTime in UTC when PartnerDevices will be marked as NonCompliant
func (m *DeviceManagementPartner) SetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("whenPartnerDevicesWillBeMarkedAsNonCompliantDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetWhenPartnerDevicesWillBeRemovedDateTime sets the whenPartnerDevicesWillBeRemovedDateTime property value. DateTime in UTC when PartnerDevices will be removed
func (m *DeviceManagementPartner) SetWhenPartnerDevicesWillBeRemovedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("whenPartnerDevicesWillBeRemovedDateTime", value)
    if err != nil {
        panic(err)
    }
}
type DeviceManagementPartnerable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetGroupsRequiringPartnerEnrollment()([]DeviceManagementPartnerAssignmentable)
    GetIsConfigured()(*bool)
    GetLastHeartbeatDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPartnerAppType()(*DeviceManagementPartnerAppType)
    GetPartnerState()(*DeviceManagementPartnerTenantState)
    GetSingleTenantAppId()(*string)
    GetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetWhenPartnerDevicesWillBeRemovedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    SetDisplayName(value *string)()
    SetGroupsRequiringPartnerEnrollment(value []DeviceManagementPartnerAssignmentable)()
    SetIsConfigured(value *bool)()
    SetLastHeartbeatDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPartnerAppType(value *DeviceManagementPartnerAppType)()
    SetPartnerState(value *DeviceManagementPartnerTenantState)()
    SetSingleTenantAppId(value *string)()
    SetWhenPartnerDevicesWillBeMarkedAsNonCompliantDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetWhenPartnerDevicesWillBeRemovedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
}
