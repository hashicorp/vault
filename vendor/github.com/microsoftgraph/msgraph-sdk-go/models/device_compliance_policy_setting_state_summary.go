package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceCompliancePolicySettingStateSummary device Compilance Policy Setting State summary across the account.
type DeviceCompliancePolicySettingStateSummary struct {
    Entity
}
// NewDeviceCompliancePolicySettingStateSummary instantiates a new DeviceCompliancePolicySettingStateSummary and sets the default values.
func NewDeviceCompliancePolicySettingStateSummary()(*DeviceCompliancePolicySettingStateSummary) {
    m := &DeviceCompliancePolicySettingStateSummary{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeviceCompliancePolicySettingStateSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceCompliancePolicySettingStateSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceCompliancePolicySettingStateSummary(), nil
}
// GetCompliantDeviceCount gets the compliantDeviceCount property value. Number of compliant devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicySettingStateSummary) GetCompliantDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("compliantDeviceCount")
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
func (m *DeviceCompliancePolicySettingStateSummary) GetConflictDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("conflictDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDeviceComplianceSettingStates gets the deviceComplianceSettingStates property value. Not yet documented
// returns a []DeviceComplianceSettingStateable when successful
func (m *DeviceCompliancePolicySettingStateSummary) GetDeviceComplianceSettingStates()([]DeviceComplianceSettingStateable) {
    val, err := m.GetBackingStore().Get("deviceComplianceSettingStates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DeviceComplianceSettingStateable)
    }
    return nil
}
// GetErrorDeviceCount gets the errorDeviceCount property value. Number of error devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicySettingStateSummary) GetErrorDeviceCount()(*int32) {
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
func (m *DeviceCompliancePolicySettingStateSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["deviceComplianceSettingStates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDeviceComplianceSettingStateFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DeviceComplianceSettingStateable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DeviceComplianceSettingStateable)
                }
            }
            m.SetDeviceComplianceSettingStates(res)
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
    res["platformType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePolicyPlatformType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPlatformType(val.(*PolicyPlatformType))
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
    res["setting"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSetting(val)
        }
        return nil
    }
    res["settingName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSettingName(val)
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
// GetNonCompliantDeviceCount gets the nonCompliantDeviceCount property value. Number of NonCompliant devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicySettingStateSummary) GetNonCompliantDeviceCount()(*int32) {
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
func (m *DeviceCompliancePolicySettingStateSummary) GetNotApplicableDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("notApplicableDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetPlatformType gets the platformType property value. Supported platform types for policies.
// returns a *PolicyPlatformType when successful
func (m *DeviceCompliancePolicySettingStateSummary) GetPlatformType()(*PolicyPlatformType) {
    val, err := m.GetBackingStore().Get("platformType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PolicyPlatformType)
    }
    return nil
}
// GetRemediatedDeviceCount gets the remediatedDeviceCount property value. Number of remediated devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicySettingStateSummary) GetRemediatedDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("remediatedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetSetting gets the setting property value. The setting class name and property name.
// returns a *string when successful
func (m *DeviceCompliancePolicySettingStateSummary) GetSetting()(*string) {
    val, err := m.GetBackingStore().Get("setting")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSettingName gets the settingName property value. Name of the setting.
// returns a *string when successful
func (m *DeviceCompliancePolicySettingStateSummary) GetSettingName()(*string) {
    val, err := m.GetBackingStore().Get("settingName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUnknownDeviceCount gets the unknownDeviceCount property value. Number of unknown devices
// returns a *int32 when successful
func (m *DeviceCompliancePolicySettingStateSummary) GetUnknownDeviceCount()(*int32) {
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
func (m *DeviceCompliancePolicySettingStateSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteInt32Value("conflictDeviceCount", m.GetConflictDeviceCount())
        if err != nil {
            return err
        }
    }
    if m.GetDeviceComplianceSettingStates() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDeviceComplianceSettingStates()))
        for i, v := range m.GetDeviceComplianceSettingStates() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("deviceComplianceSettingStates", cast)
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
    if m.GetPlatformType() != nil {
        cast := (*m.GetPlatformType()).String()
        err = writer.WriteStringValue("platformType", &cast)
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
        err = writer.WriteStringValue("setting", m.GetSetting())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("settingName", m.GetSettingName())
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
func (m *DeviceCompliancePolicySettingStateSummary) SetCompliantDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("compliantDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetConflictDeviceCount sets the conflictDeviceCount property value. Number of conflict devices
func (m *DeviceCompliancePolicySettingStateSummary) SetConflictDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("conflictDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceComplianceSettingStates sets the deviceComplianceSettingStates property value. Not yet documented
func (m *DeviceCompliancePolicySettingStateSummary) SetDeviceComplianceSettingStates(value []DeviceComplianceSettingStateable)() {
    err := m.GetBackingStore().Set("deviceComplianceSettingStates", value)
    if err != nil {
        panic(err)
    }
}
// SetErrorDeviceCount sets the errorDeviceCount property value. Number of error devices
func (m *DeviceCompliancePolicySettingStateSummary) SetErrorDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("errorDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNonCompliantDeviceCount sets the nonCompliantDeviceCount property value. Number of NonCompliant devices
func (m *DeviceCompliancePolicySettingStateSummary) SetNonCompliantDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("nonCompliantDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetNotApplicableDeviceCount sets the notApplicableDeviceCount property value. Number of not applicable devices
func (m *DeviceCompliancePolicySettingStateSummary) SetNotApplicableDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("notApplicableDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetPlatformType sets the platformType property value. Supported platform types for policies.
func (m *DeviceCompliancePolicySettingStateSummary) SetPlatformType(value *PolicyPlatformType)() {
    err := m.GetBackingStore().Set("platformType", value)
    if err != nil {
        panic(err)
    }
}
// SetRemediatedDeviceCount sets the remediatedDeviceCount property value. Number of remediated devices
func (m *DeviceCompliancePolicySettingStateSummary) SetRemediatedDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("remediatedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetSetting sets the setting property value. The setting class name and property name.
func (m *DeviceCompliancePolicySettingStateSummary) SetSetting(value *string)() {
    err := m.GetBackingStore().Set("setting", value)
    if err != nil {
        panic(err)
    }
}
// SetSettingName sets the settingName property value. Name of the setting.
func (m *DeviceCompliancePolicySettingStateSummary) SetSettingName(value *string)() {
    err := m.GetBackingStore().Set("settingName", value)
    if err != nil {
        panic(err)
    }
}
// SetUnknownDeviceCount sets the unknownDeviceCount property value. Number of unknown devices
func (m *DeviceCompliancePolicySettingStateSummary) SetUnknownDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("unknownDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
type DeviceCompliancePolicySettingStateSummaryable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCompliantDeviceCount()(*int32)
    GetConflictDeviceCount()(*int32)
    GetDeviceComplianceSettingStates()([]DeviceComplianceSettingStateable)
    GetErrorDeviceCount()(*int32)
    GetNonCompliantDeviceCount()(*int32)
    GetNotApplicableDeviceCount()(*int32)
    GetPlatformType()(*PolicyPlatformType)
    GetRemediatedDeviceCount()(*int32)
    GetSetting()(*string)
    GetSettingName()(*string)
    GetUnknownDeviceCount()(*int32)
    SetCompliantDeviceCount(value *int32)()
    SetConflictDeviceCount(value *int32)()
    SetDeviceComplianceSettingStates(value []DeviceComplianceSettingStateable)()
    SetErrorDeviceCount(value *int32)()
    SetNonCompliantDeviceCount(value *int32)()
    SetNotApplicableDeviceCount(value *int32)()
    SetPlatformType(value *PolicyPlatformType)()
    SetRemediatedDeviceCount(value *int32)()
    SetSetting(value *string)()
    SetSettingName(value *string)()
    SetUnknownDeviceCount(value *int32)()
}
