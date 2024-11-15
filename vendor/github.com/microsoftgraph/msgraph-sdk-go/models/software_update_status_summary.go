package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SoftwareUpdateStatusSummary struct {
    Entity
}
// NewSoftwareUpdateStatusSummary instantiates a new SoftwareUpdateStatusSummary and sets the default values.
func NewSoftwareUpdateStatusSummary()(*SoftwareUpdateStatusSummary) {
    m := &SoftwareUpdateStatusSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSoftwareUpdateStatusSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSoftwareUpdateStatusSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSoftwareUpdateStatusSummary(), nil
}
// GetCompliantDeviceCount gets the compliantDeviceCount property value. Number of compliant devices.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetCompliantDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("compliantDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetCompliantUserCount gets the compliantUserCount property value. Number of compliant users.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetCompliantUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("compliantUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetConflictDeviceCount gets the conflictDeviceCount property value. Number of conflict devices.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetConflictDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("conflictDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetConflictUserCount gets the conflictUserCount property value. Number of conflict users.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetConflictUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("conflictUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The name of the policy.
// returns a *string when successful
func (m *SoftwareUpdateStatusSummary) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetErrorDeviceCount gets the errorDeviceCount property value. Number of devices had error.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetErrorDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("errorDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetErrorUserCount gets the errorUserCount property value. Number of users had error.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetErrorUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("errorUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SoftwareUpdateStatusSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["compliantDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompliantDeviceCount(val)
        }
        return nil
    }
    res["compliantUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompliantUserCount(val)
        }
        return nil
    }
    res["conflictDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConflictDeviceCount(val)
        }
        return nil
    }
    res["conflictUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConflictUserCount(val)
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
    res["errorDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetErrorDeviceCount(val)
        }
        return nil
    }
    res["errorUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetErrorUserCount(val)
        }
        return nil
    }
    res["nonCompliantDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNonCompliantDeviceCount(val)
        }
        return nil
    }
    res["nonCompliantUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNonCompliantUserCount(val)
        }
        return nil
    }
    res["notApplicableDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotApplicableDeviceCount(val)
        }
        return nil
    }
    res["notApplicableUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotApplicableUserCount(val)
        }
        return nil
    }
    res["remediatedDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemediatedDeviceCount(val)
        }
        return nil
    }
    res["remediatedUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemediatedUserCount(val)
        }
        return nil
    }
    res["unknownDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnknownDeviceCount(val)
        }
        return nil
    }
    res["unknownUserCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnknownUserCount(val)
        }
        return nil
    }
    return res
}
// GetNonCompliantDeviceCount gets the nonCompliantDeviceCount property value. Number of non compliant devices.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetNonCompliantDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("nonCompliantDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetNonCompliantUserCount gets the nonCompliantUserCount property value. Number of non compliant users.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetNonCompliantUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("nonCompliantUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetNotApplicableDeviceCount gets the notApplicableDeviceCount property value. Number of not applicable devices.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetNotApplicableDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("notApplicableDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetNotApplicableUserCount gets the notApplicableUserCount property value. Number of not applicable users.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetNotApplicableUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("notApplicableUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetRemediatedDeviceCount gets the remediatedDeviceCount property value. Number of remediated devices.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetRemediatedDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("remediatedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetRemediatedUserCount gets the remediatedUserCount property value. Number of remediated users.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetRemediatedUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("remediatedUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUnknownDeviceCount gets the unknownDeviceCount property value. Number of unknown devices.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetUnknownDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("unknownDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUnknownUserCount gets the unknownUserCount property value. Number of unknown users.
// returns a *int32 when successful
func (m *SoftwareUpdateStatusSummary) GetUnknownUserCount()(*int32) {
    val, err := m.GetBackingStore().Get("unknownUserCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SoftwareUpdateStatusSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("compliantDeviceCount", m.GetCompliantDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("compliantUserCount", m.GetCompliantUserCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("conflictDeviceCount", m.GetConflictDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("conflictUserCount", m.GetConflictUserCount())
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
        err = writer.WriteInt32Value("errorDeviceCount", m.GetErrorDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("errorUserCount", m.GetErrorUserCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("nonCompliantDeviceCount", m.GetNonCompliantDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("nonCompliantUserCount", m.GetNonCompliantUserCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("notApplicableDeviceCount", m.GetNotApplicableDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("notApplicableUserCount", m.GetNotApplicableUserCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("remediatedDeviceCount", m.GetRemediatedDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("remediatedUserCount", m.GetRemediatedUserCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("unknownDeviceCount", m.GetUnknownDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("unknownUserCount", m.GetUnknownUserCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCompliantDeviceCount sets the compliantDeviceCount property value. Number of compliant devices.
func (m *SoftwareUpdateStatusSummary) SetCompliantDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("compliantDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetCompliantUserCount sets the compliantUserCount property value. Number of compliant users.
func (m *SoftwareUpdateStatusSummary) SetCompliantUserCount(value *int32)() {
    err := m.GetBackingStore().Set("compliantUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetConflictDeviceCount sets the conflictDeviceCount property value. Number of conflict devices.
func (m *SoftwareUpdateStatusSummary) SetConflictDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("conflictDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetConflictUserCount sets the conflictUserCount property value. Number of conflict users.
func (m *SoftwareUpdateStatusSummary) SetConflictUserCount(value *int32)() {
    err := m.GetBackingStore().Set("conflictUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The name of the policy.
func (m *SoftwareUpdateStatusSummary) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetErrorDeviceCount sets the errorDeviceCount property value. Number of devices had error.
func (m *SoftwareUpdateStatusSummary) SetErrorDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("errorDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetErrorUserCount sets the errorUserCount property value. Number of users had error.
func (m *SoftwareUpdateStatusSummary) SetErrorUserCount(value *int32)() {
    err := m.GetBackingStore().Set("errorUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNonCompliantDeviceCount sets the nonCompliantDeviceCount property value. Number of non compliant devices.
func (m *SoftwareUpdateStatusSummary) SetNonCompliantDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("nonCompliantDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNonCompliantUserCount sets the nonCompliantUserCount property value. Number of non compliant users.
func (m *SoftwareUpdateStatusSummary) SetNonCompliantUserCount(value *int32)() {
    err := m.GetBackingStore().Set("nonCompliantUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNotApplicableDeviceCount sets the notApplicableDeviceCount property value. Number of not applicable devices.
func (m *SoftwareUpdateStatusSummary) SetNotApplicableDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("notApplicableDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNotApplicableUserCount sets the notApplicableUserCount property value. Number of not applicable users.
func (m *SoftwareUpdateStatusSummary) SetNotApplicableUserCount(value *int32)() {
    err := m.GetBackingStore().Set("notApplicableUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetRemediatedDeviceCount sets the remediatedDeviceCount property value. Number of remediated devices.
func (m *SoftwareUpdateStatusSummary) SetRemediatedDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("remediatedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetRemediatedUserCount sets the remediatedUserCount property value. Number of remediated users.
func (m *SoftwareUpdateStatusSummary) SetRemediatedUserCount(value *int32)() {
    err := m.GetBackingStore().Set("remediatedUserCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUnknownDeviceCount sets the unknownDeviceCount property value. Number of unknown devices.
func (m *SoftwareUpdateStatusSummary) SetUnknownDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("unknownDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUnknownUserCount sets the unknownUserCount property value. Number of unknown users.
func (m *SoftwareUpdateStatusSummary) SetUnknownUserCount(value *int32)() {
    err := m.GetBackingStore().Set("unknownUserCount", value)
    if err != nil {
        panic(err)
    }
}
type SoftwareUpdateStatusSummaryable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCompliantDeviceCount()(*int32)
    GetCompliantUserCount()(*int32)
    GetConflictDeviceCount()(*int32)
    GetConflictUserCount()(*int32)
    GetDisplayName()(*string)
    GetErrorDeviceCount()(*int32)
    GetErrorUserCount()(*int32)
    GetNonCompliantDeviceCount()(*int32)
    GetNonCompliantUserCount()(*int32)
    GetNotApplicableDeviceCount()(*int32)
    GetNotApplicableUserCount()(*int32)
    GetRemediatedDeviceCount()(*int32)
    GetRemediatedUserCount()(*int32)
    GetUnknownDeviceCount()(*int32)
    GetUnknownUserCount()(*int32)
    SetCompliantDeviceCount(value *int32)()
    SetCompliantUserCount(value *int32)()
    SetConflictDeviceCount(value *int32)()
    SetConflictUserCount(value *int32)()
    SetDisplayName(value *string)()
    SetErrorDeviceCount(value *int32)()
    SetErrorUserCount(value *int32)()
    SetNonCompliantDeviceCount(value *int32)()
    SetNonCompliantUserCount(value *int32)()
    SetNotApplicableDeviceCount(value *int32)()
    SetNotApplicableUserCount(value *int32)()
    SetRemediatedDeviceCount(value *int32)()
    SetRemediatedUserCount(value *int32)()
    SetUnknownDeviceCount(value *int32)()
    SetUnknownUserCount(value *int32)()
}
