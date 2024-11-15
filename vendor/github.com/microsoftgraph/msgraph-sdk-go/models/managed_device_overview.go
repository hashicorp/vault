package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedDeviceOverview summary data for managed devices
type ManagedDeviceOverview struct {
    Entity
}
// NewManagedDeviceOverview instantiates a new ManagedDeviceOverview and sets the default values.
func NewManagedDeviceOverview()(*ManagedDeviceOverview) {
    m := &ManagedDeviceOverview{
        Entity: *NewEntity(),
    }
    return m
}
// CreateManagedDeviceOverviewFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedDeviceOverviewFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewManagedDeviceOverview(), nil
}
// GetDeviceExchangeAccessStateSummary gets the deviceExchangeAccessStateSummary property value. Distribution of Exchange Access State in Intune
// returns a DeviceExchangeAccessStateSummaryable when successful
func (m *ManagedDeviceOverview) GetDeviceExchangeAccessStateSummary()(DeviceExchangeAccessStateSummaryable) {
    val, err := m.GetBackingStore().Get("deviceExchangeAccessStateSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceExchangeAccessStateSummaryable)
    }
    return nil
}
// GetDeviceOperatingSystemSummary gets the deviceOperatingSystemSummary property value. Device operating system summary.
// returns a DeviceOperatingSystemSummaryable when successful
func (m *ManagedDeviceOverview) GetDeviceOperatingSystemSummary()(DeviceOperatingSystemSummaryable) {
    val, err := m.GetBackingStore().Get("deviceOperatingSystemSummary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DeviceOperatingSystemSummaryable)
    }
    return nil
}
// GetDualEnrolledDeviceCount gets the dualEnrolledDeviceCount property value. The number of devices enrolled in both MDM and EAS
// returns a *int32 when successful
func (m *ManagedDeviceOverview) GetDualEnrolledDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("dualEnrolledDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetEnrolledDeviceCount gets the enrolledDeviceCount property value. Total enrolled device count. Does not include PC devices managed via Intune PC Agent
// returns a *int32 when successful
func (m *ManagedDeviceOverview) GetEnrolledDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("enrolledDeviceCount")
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
func (m *ManagedDeviceOverview) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["deviceExchangeAccessStateSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceExchangeAccessStateSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceExchangeAccessStateSummary(val.(DeviceExchangeAccessStateSummaryable))
        }
        return nil
    }
    res["deviceOperatingSystemSummary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDeviceOperatingSystemSummaryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceOperatingSystemSummary(val.(DeviceOperatingSystemSummaryable))
        }
        return nil
    }
    res["dualEnrolledDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDualEnrolledDeviceCount(val)
        }
        return nil
    }
    res["enrolledDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEnrolledDeviceCount(val)
        }
        return nil
    }
    res["mdmEnrolledCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMdmEnrolledCount(val)
        }
        return nil
    }
    return res
}
// GetMdmEnrolledCount gets the mdmEnrolledCount property value. The number of devices enrolled in MDM
// returns a *int32 when successful
func (m *ManagedDeviceOverview) GetMdmEnrolledCount()(*int32) {
    val, err := m.GetBackingStore().Get("mdmEnrolledCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ManagedDeviceOverview) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("deviceExchangeAccessStateSummary", m.GetDeviceExchangeAccessStateSummary())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("deviceOperatingSystemSummary", m.GetDeviceOperatingSystemSummary())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("dualEnrolledDeviceCount", m.GetDualEnrolledDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("enrolledDeviceCount", m.GetEnrolledDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("mdmEnrolledCount", m.GetMdmEnrolledCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeviceExchangeAccessStateSummary sets the deviceExchangeAccessStateSummary property value. Distribution of Exchange Access State in Intune
func (m *ManagedDeviceOverview) SetDeviceExchangeAccessStateSummary(value DeviceExchangeAccessStateSummaryable)() {
    err := m.GetBackingStore().Set("deviceExchangeAccessStateSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceOperatingSystemSummary sets the deviceOperatingSystemSummary property value. Device operating system summary.
func (m *ManagedDeviceOverview) SetDeviceOperatingSystemSummary(value DeviceOperatingSystemSummaryable)() {
    err := m.GetBackingStore().Set("deviceOperatingSystemSummary", value)
    if err != nil {
        panic(err)
    }
}
// SetDualEnrolledDeviceCount sets the dualEnrolledDeviceCount property value. The number of devices enrolled in both MDM and EAS
func (m *ManagedDeviceOverview) SetDualEnrolledDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("dualEnrolledDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetEnrolledDeviceCount sets the enrolledDeviceCount property value. Total enrolled device count. Does not include PC devices managed via Intune PC Agent
func (m *ManagedDeviceOverview) SetEnrolledDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("enrolledDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetMdmEnrolledCount sets the mdmEnrolledCount property value. The number of devices enrolled in MDM
func (m *ManagedDeviceOverview) SetMdmEnrolledCount(value *int32)() {
    err := m.GetBackingStore().Set("mdmEnrolledCount", value)
    if err != nil {
        panic(err)
    }
}
type ManagedDeviceOverviewable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeviceExchangeAccessStateSummary()(DeviceExchangeAccessStateSummaryable)
    GetDeviceOperatingSystemSummary()(DeviceOperatingSystemSummaryable)
    GetDualEnrolledDeviceCount()(*int32)
    GetEnrolledDeviceCount()(*int32)
    GetMdmEnrolledCount()(*int32)
    SetDeviceExchangeAccessStateSummary(value DeviceExchangeAccessStateSummaryable)()
    SetDeviceOperatingSystemSummary(value DeviceOperatingSystemSummaryable)()
    SetDualEnrolledDeviceCount(value *int32)()
    SetEnrolledDeviceCount(value *int32)()
    SetMdmEnrolledCount(value *int32)()
}
