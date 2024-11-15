package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DeviceCompliancePolicyDeviceStateSummary struct {
    Entity
}
// NewDeviceCompliancePolicyDeviceStateSummary instantiates a new DeviceCompliancePolicyDeviceStateSummary and sets the default values.
func NewDeviceCompliancePolicyDeviceStateSummary()(*DeviceCompliancePolicyDeviceStateSummary) {
    m := &DeviceCompliancePolicyDeviceStateSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceCompliancePolicyDeviceStateSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceCompliancePolicyDeviceStateSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceCompliancePolicyDeviceStateSummary(), nil
}
// GetCompliantDeviceCount gets the compliantDeviceCount property value. Number of compliant devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicyDeviceStateSummary) GetCompliantDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("compliantDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetConfigManagerCount gets the configManagerCount property value. Number of devices that have compliance managed by System Center Configuration Manager
// returns a *int32 when successful
func (m *DeviceCompliancePolicyDeviceStateSummary) GetConfigManagerCount()(*int32) {
    val, err := m.GetBackingStore().Get("configManagerCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetConflictDeviceCount gets the conflictDeviceCount property value. Number of conflict devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicyDeviceStateSummary) GetConflictDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("conflictDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetErrorDeviceCount gets the errorDeviceCount property value. Number of error devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicyDeviceStateSummary) GetErrorDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("errorDeviceCount")
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
func (m *DeviceCompliancePolicyDeviceStateSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["configManagerCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfigManagerCount(val)
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
    res["inGracePeriodCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInGracePeriodCount(val)
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
    return res
}
// GetInGracePeriodCount gets the inGracePeriodCount property value. Number of devices that are in grace period
// returns a *int32 when successful
func (m *DeviceCompliancePolicyDeviceStateSummary) GetInGracePeriodCount()(*int32) {
    val, err := m.GetBackingStore().Get("inGracePeriodCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetNonCompliantDeviceCount gets the nonCompliantDeviceCount property value. Number of NonCompliant devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicyDeviceStateSummary) GetNonCompliantDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("nonCompliantDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetNotApplicableDeviceCount gets the notApplicableDeviceCount property value. Number of not applicable devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicyDeviceStateSummary) GetNotApplicableDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("notApplicableDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetRemediatedDeviceCount gets the remediatedDeviceCount property value. Number of remediated devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicyDeviceStateSummary) GetRemediatedDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("remediatedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUnknownDeviceCount gets the unknownDeviceCount property value. Number of unknown devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicyDeviceStateSummary) GetUnknownDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("unknownDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceCompliancePolicyDeviceStateSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteInt32Value("configManagerCount", m.GetConfigManagerCount())
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
        err = writer.WriteInt32Value("errorDeviceCount", m.GetErrorDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("inGracePeriodCount", m.GetInGracePeriodCount())
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
        err = writer.WriteInt32Value("notApplicableDeviceCount", m.GetNotApplicableDeviceCount())
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
        err = writer.WriteInt32Value("unknownDeviceCount", m.GetUnknownDeviceCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCompliantDeviceCount sets the compliantDeviceCount property value. Number of compliant devices
func (m *DeviceCompliancePolicyDeviceStateSummary) SetCompliantDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("compliantDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetConfigManagerCount sets the configManagerCount property value. Number of devices that have compliance managed by System Center Configuration Manager
func (m *DeviceCompliancePolicyDeviceStateSummary) SetConfigManagerCount(value *int32)() {
    err := m.GetBackingStore().Set("configManagerCount", value)
    if err != nil {
        panic(err)
    }
}
// SetConflictDeviceCount sets the conflictDeviceCount property value. Number of conflict devices
func (m *DeviceCompliancePolicyDeviceStateSummary) SetConflictDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("conflictDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetErrorDeviceCount sets the errorDeviceCount property value. Number of error devices
func (m *DeviceCompliancePolicyDeviceStateSummary) SetErrorDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("errorDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetInGracePeriodCount sets the inGracePeriodCount property value. Number of devices that are in grace period
func (m *DeviceCompliancePolicyDeviceStateSummary) SetInGracePeriodCount(value *int32)() {
    err := m.GetBackingStore().Set("inGracePeriodCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNonCompliantDeviceCount sets the nonCompliantDeviceCount property value. Number of NonCompliant devices
func (m *DeviceCompliancePolicyDeviceStateSummary) SetNonCompliantDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("nonCompliantDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNotApplicableDeviceCount sets the notApplicableDeviceCount property value. Number of not applicable devices
func (m *DeviceCompliancePolicyDeviceStateSummary) SetNotApplicableDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("notApplicableDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetRemediatedDeviceCount sets the remediatedDeviceCount property value. Number of remediated devices
func (m *DeviceCompliancePolicyDeviceStateSummary) SetRemediatedDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("remediatedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUnknownDeviceCount sets the unknownDeviceCount property value. Number of unknown devices
func (m *DeviceCompliancePolicyDeviceStateSummary) SetUnknownDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("unknownDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
type DeviceCompliancePolicyDeviceStateSummaryable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCompliantDeviceCount()(*int32)
    GetConfigManagerCount()(*int32)
    GetConflictDeviceCount()(*int32)
    GetErrorDeviceCount()(*int32)
    GetInGracePeriodCount()(*int32)
    GetNonCompliantDeviceCount()(*int32)
    GetNotApplicableDeviceCount()(*int32)
    GetRemediatedDeviceCount()(*int32)
    GetUnknownDeviceCount()(*int32)
    SetCompliantDeviceCount(value *int32)()
    SetConfigManagerCount(value *int32)()
    SetConflictDeviceCount(value *int32)()
    SetErrorDeviceCount(value *int32)()
    SetInGracePeriodCount(value *int32)()
    SetNonCompliantDeviceCount(value *int32)()
    SetNotApplicableDeviceCount(value *int32)()
    SetRemediatedDeviceCount(value *int32)()
    SetUnknownDeviceCount(value *int32)()
}
