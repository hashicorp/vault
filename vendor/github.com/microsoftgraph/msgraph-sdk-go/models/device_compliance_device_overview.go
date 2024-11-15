package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DeviceComplianceDeviceOverview struct {
    Entity
}
// NewDeviceComplianceDeviceOverview instantiates a new DeviceComplianceDeviceOverview and sets the default values.
func NewDeviceComplianceDeviceOverview()(*DeviceComplianceDeviceOverview) {
    m := &DeviceComplianceDeviceOverview{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceComplianceDeviceOverviewFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceComplianceDeviceOverviewFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceComplianceDeviceOverview(), nil
}
// GetConfigurationVersion gets the configurationVersion property value. Version of the policy for that overview
// returns a *int32 when successful
func (m *DeviceComplianceDeviceOverview) GetConfigurationVersion()(*int32) {
    val, err := m.GetBackingStore().Get("configurationVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetErrorCount gets the errorCount property value. Number of error devices
// returns a *int32 when successful
func (m *DeviceComplianceDeviceOverview) GetErrorCount()(*int32) {
    val, err := m.GetBackingStore().Get("errorCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFailedCount gets the failedCount property value. Number of failed devices
// returns a *int32 when successful
func (m *DeviceComplianceDeviceOverview) GetFailedCount()(*int32) {
    val, err := m.GetBackingStore().Get("failedCount")
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
func (m *DeviceComplianceDeviceOverview) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["configurationVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConfigurationVersion(val)
        }
        return nil
    }
    res["errorCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetErrorCount(val)
        }
        return nil
    }
    res["failedCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFailedCount(val)
        }
        return nil
    }
    res["lastUpdateDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastUpdateDateTime(val)
        }
        return nil
    }
    res["notApplicableCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotApplicableCount(val)
        }
        return nil
    }
    res["pendingCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPendingCount(val)
        }
        return nil
    }
    res["successCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSuccessCount(val)
        }
        return nil
    }
    return res
}
// GetLastUpdateDateTime gets the lastUpdateDateTime property value. Last update time
// returns a *Time when successful
func (m *DeviceComplianceDeviceOverview) GetLastUpdateDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastUpdateDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetNotApplicableCount gets the notApplicableCount property value. Number of not applicable devices
// returns a *int32 when successful
func (m *DeviceComplianceDeviceOverview) GetNotApplicableCount()(*int32) {
    val, err := m.GetBackingStore().Get("notApplicableCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPendingCount gets the pendingCount property value. Number of pending devices
// returns a *int32 when successful
func (m *DeviceComplianceDeviceOverview) GetPendingCount()(*int32) {
    val, err := m.GetBackingStore().Get("pendingCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSuccessCount gets the successCount property value. Number of succeeded devices
// returns a *int32 when successful
func (m *DeviceComplianceDeviceOverview) GetSuccessCount()(*int32) {
    val, err := m.GetBackingStore().Get("successCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceComplianceDeviceOverview) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt32Value("configurationVersion", m.GetConfigurationVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("errorCount", m.GetErrorCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("failedCount", m.GetFailedCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastUpdateDateTime", m.GetLastUpdateDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("notApplicableCount", m.GetNotApplicableCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("pendingCount", m.GetPendingCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("successCount", m.GetSuccessCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetConfigurationVersion sets the configurationVersion property value. Version of the policy for that overview
func (m *DeviceComplianceDeviceOverview) SetConfigurationVersion(value *int32)() {
    err := m.GetBackingStore().Set("configurationVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetErrorCount sets the errorCount property value. Number of error devices
func (m *DeviceComplianceDeviceOverview) SetErrorCount(value *int32)() {
    err := m.GetBackingStore().Set("errorCount", value)
    if err != nil {
        panic(err)
    }
}
// SetFailedCount sets the failedCount property value. Number of failed devices
func (m *DeviceComplianceDeviceOverview) SetFailedCount(value *int32)() {
    err := m.GetBackingStore().Set("failedCount", value)
    if err != nil {
        panic(err)
    }
}
// SetLastUpdateDateTime sets the lastUpdateDateTime property value. Last update time
func (m *DeviceComplianceDeviceOverview) SetLastUpdateDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastUpdateDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetNotApplicableCount sets the notApplicableCount property value. Number of not applicable devices
func (m *DeviceComplianceDeviceOverview) SetNotApplicableCount(value *int32)() {
    err := m.GetBackingStore().Set("notApplicableCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPendingCount sets the pendingCount property value. Number of pending devices
func (m *DeviceComplianceDeviceOverview) SetPendingCount(value *int32)() {
    err := m.GetBackingStore().Set("pendingCount", value)
    if err != nil {
        panic(err)
    }
}
// SetSuccessCount sets the successCount property value. Number of succeeded devices
func (m *DeviceComplianceDeviceOverview) SetSuccessCount(value *int32)() {
    err := m.GetBackingStore().Set("successCount", value)
    if err != nil {
        panic(err)
    }
}
type DeviceComplianceDeviceOverviewable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetConfigurationVersion()(*int32)
    GetErrorCount()(*int32)
    GetFailedCount()(*int32)
    GetLastUpdateDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetNotApplicableCount()(*int32)
    GetPendingCount()(*int32)
    GetSuccessCount()(*int32)
    SetConfigurationVersion(value *int32)()
    SetErrorCount(value *int32)()
    SetFailedCount(value *int32)()
    SetLastUpdateDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetNotApplicableCount(value *int32)()
    SetPendingCount(value *int32)()
    SetSuccessCount(value *int32)()
}
